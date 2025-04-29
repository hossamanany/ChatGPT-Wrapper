import { defineStore } from "pinia";
import { computed, ref } from "vue";
import type { Chat } from "@/models/chat.model";
import { db } from "@/db";
import Dexie from "dexie";
import type { Message } from "@/models/message.model";
import { Role } from "@/models/role.model";
import { streamChatResponse } from "@/services/chat.service";
import { generateTitle as generateChatTitle } from "@/services/title.service";

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

    isStreaming.value = true;
    try {
      await streamChatResponse();
    } catch (error) {
      console.error("Error sending message:", error);
      isStreaming.value = false;
    }
  }

  async function generateTitle() {
    if (!currentChat.value || currentChat.value.title) return;

    const firstUserMessage = currentChat.value.messages[1]?.content;
    if (!firstUserMessage) return;

    try {
      await generateChatTitle(firstUserMessage);
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
