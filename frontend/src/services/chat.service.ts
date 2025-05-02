import type { ChatCompletionMessageParam } from "openai/resources/chat/completions"; // TODO: We should use our own type
import { useChatStore } from "@/stores/chat.store";

/**
 * Fetches the chat stream response from the API.
 */
async function fetchChatStream(
  messages: ChatCompletionMessageParam[]
): Promise<Response> {
  const response = await fetch("/api/chat/stream", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ messages }),
  });

  if (!response.ok) {
    const errorData = await response.json();
    throw new Error(errorData.error || "Failed to get streaming response");
  }

  return response;
}

/**
 * Initializes the stream reader from the response.
 */
function getStreamReader(
  response: Response
): ReadableStreamDefaultReader<Uint8Array> {
  const reader = response.body?.getReader();
  if (!reader) {
    throw new Error("Failed to initialize stream reader");
  }
  return reader;
}

/**
 * Processes the streaming response line by line.
 */
async function processStream(
  reader: ReadableStreamDefaultReader<Uint8Array>,
  onContent: (content: string) => Promise<void>
): Promise<void> {
  const decoder = new TextDecoder();
  let buffer = "";

  try {
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
          await onContent(content);
        }
      }
    }
  } finally {
    reader.releaseLock();
  }
}

/**
 * Main function to stream chat response and update chat store.
 */
export async function streamChatResponse(): Promise<void> {
  const chatStore = useChatStore();

  if (!chatStore.currentChat) {
    return;
  }

  const response = await fetchChatStream(
    chatStore.currentChat.messages as ChatCompletionMessageParam[]
  );
  const reader = getStreamReader(response);

  await processStream(reader, async (content) => {
    await chatStore.updateLastMessageStream(content);
  });
}
