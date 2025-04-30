import type { ChatCompletionMessageParam } from "openai/resources/chat/completions";
import { useChatStore } from "@/stores/chat.store";

export async function streamChatResponse(): Promise<void> {
  const chatStore = useChatStore();
  const openaiModel = import.meta.env.VITE_OPENAI_MODEL || "gpt-3.5-turbo";
  const openaiTemperature =
    Number(import.meta.env.VITE_OPENAI_TEMPERATURE) || 0.5;
  const openaiMaxTokens =
    Number(import.meta.env.VITE_OPENAI_MAX_TOKENS) || 1000;

  if (!chatStore.currentChat) {
    return;
  }

  const response = await fetch("/api/chat/stream", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({
      messages: chatStore.currentChat.messages as ChatCompletionMessageParam[],
      model: openaiModel,
      temperature: openaiTemperature,
      max_tokens: openaiMaxTokens,
      stream: true,
    }),
  });

  if (!response.ok) {
    const errorData = await response.json();
    throw new Error(errorData.error || "Failed to get streaming response");
  }

  const reader = response.body?.getReader();
  const decoder = new TextDecoder();

  if (!reader) {
    throw new Error("Failed to initialize stream reader");
  }

  try {
    let buffer = "";
    while (true) {
      const { done, value } = await reader.read();
      if (done) break;

      buffer += decoder.decode(value, { stream: true });
      const lines = buffer.split("\n");
      buffer = lines.pop() || "";

      for (const line of lines) {
        if (!line.trim()) continue;
        const parsed = JSON.parse(line);
        const content = parsed.choices[0]?.delta?.content;
        if (content) {
          await chatStore.updateLastMessageStream(content);
        }
      }
    }
  } finally {
    reader.releaseLock();
  }
}
