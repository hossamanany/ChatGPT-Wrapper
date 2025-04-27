import { defineStore } from "pinia";
import { computed, ref } from "vue";
import type { Chat } from "@/models/chat.model";
import { db } from "@/db";
import Dexie from "dexie";
import type { Message } from "@/models/message.model";
import { Role } from "@/models/role.model";
import { ApiService } from "@/services/api";
import type { ChatCompletionRequest } from "@/models/api.model";

export const useChatStore = defineStore("chat", () => {
  const chats = ref<Chat[]>([]);
  const currentChatId = ref<number | null>(null);
  const isStreaming = ref(false);

  const currentChat = computed(() =>
    chats.value.find((item) => item.id === currentChatId.value)
  );

  function setCurrentChatId(id: number | null) {
    currentChatId.value = id;
  }

  async function reloadChats() {
    try {
      chats.value = await db.chats.reverse().toArray();
    } catch (e) {
      console.error(e);
    }
  }

  async function addMessage(message: Message) {
    if (currentChat.value) {
      const messages = [
        ...Dexie.deepClone(currentChat.value.messages),
        message,
      ];
      try {
        await db.chats.update(currentChatId.value, { messages });
      } catch (e) {
        console.error(e);
      }
    } else {
      try {
        setCurrentChatId(
          await db.chats.add({
            title: null,
            messages: [
              { role: Role.system, content: "Always answer in markdown" },
              message,
            ],
          })
        );
      } catch (e) {
        console.error(e);
      }
    }
    await reloadChats();
  }

  async function sendMessage(content: string) {
    if (!currentChat.value) return;

    const userMessage: Message = { role: Role.user, content };
    await addMessage(userMessage);

    const request: ChatCompletionRequest = {
      model: "gpt-3.5-turbo",
      messages: currentChat.value.messages.map((msg) => ({
        role: msg.role as "system" | "user" | "assistant",
        content: msg.content || "",
      })),
      stream: true,
    };

    isStreaming.value = true;
    try {
      await ApiService.createChatCompletionStream(
        request,
        (chunk) => updateLastMessageStream(chunk),
        (error) => {
          console.error("Stream error:", error);
          isStreaming.value = false;
        },
        () => {
          isStreaming.value = false;
          generateTitle();
        }
      );
    } catch (error) {
      console.error("Error sending message:", error);
      isStreaming.value = false;
    }
  }

  async function generateTitle() {
    if (!currentChat.value || currentChat.value.title) return;

    const request: ChatCompletionRequest = {
      model: "gpt-3.5-turbo",
      messages: [
        {
          role: Role.system,
          content:
            "Generate a short title (max 5 words) for this conversation based on the first message.",
        },
        ...currentChat.value.messages.slice(0, 2).map((msg) => ({
          role: msg.role as "system" | "user" | "assistant",
          content: msg.content || "",
        })),
      ],
    };

    try {
      const response = await ApiService.generateTitle(request);
      if (response.choices[0]?.message?.content) {
        await setCurrentChatTitle(response.choices[0].message.content.trim());
      }
    } catch (error) {
      console.error("Error generating title:", error);
    }
  }

  async function setCurrentChatTitle(title: string | null) {
    if (currentChat.value) {
      try {
        await db.chats.update(currentChatId.value, { title });
      } catch (e) {
        console.error(e);
      }
      await reloadChats();
    }
  }

  async function updateLastMessageStream(messageChunk: string) {
    if (currentChat.value) {
      const messages = Dexie.deepClone(currentChat.value.messages);
      const lastMessage = messages[messages.length - 1];
      if (lastMessage.role === Role.assistant) {
        lastMessage.content = (lastMessage.content || "") + messageChunk;
      } else {
        messages.push({
          role: Role.assistant,
          content: messageChunk,
        });
      }
      try {
        await db.chats.update(currentChatId.value, { messages });
      } catch (e) {
        console.error(e);
      }
      await reloadChats();
    }
  }

  async function deleteChat(chatId: number | null | undefined) {
    try {
      await db.chats.delete(chatId);
    } catch (e) {
      console.error(e);
    }
    await reloadChats();
  }

  return {
    chats,
    currentChatId,
    currentChat,
    isStreaming,
    setCurrentChatId,
    reloadChats,
    sendMessage,
    setCurrentChatTitle,
    updateLastMessageStream,
    deleteChat,
    addMessage,
  };
});
