import type {
  ChatCompletionRequest,
  ChatCompletionResponse,
  ChatCompletionStreamResponse,
} from "@/models/api.model";

export const createChatCompletion = async (
  request: ChatCompletionRequest
): Promise<ChatCompletionResponse> => {
  const response = await fetch("/api/chat/completions", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(request),
  });

  if (!response.ok) {
    throw new Error("Failed to create chat completion");
  }

  return response.json();
};

export const createChatCompletionStream = async (
  request: ChatCompletionRequest
): Promise<ReadableStream<ChatCompletionStreamResponse>> => {
  const response = await fetch("/api/chat/completions/stream", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(request),
  });

  if (!response.ok) {
    throw new Error("Failed to create chat completion stream");
  }

  if (!response.body) {
    throw new Error("Response body is null");
  }

  return response.body.pipeThrough(
    new TransformStream({
      transform(chunk, controller) {
        const text = new TextDecoder().decode(chunk);
        const lines = text.split("\n").filter((line) => line.trim() !== "");
        for (const line of lines) {
          if (line.startsWith("data: ")) {
            const data = line.slice(6);
            if (data === "[DONE]") {
              controller.terminate();
              return;
            }
            try {
              const parsed = JSON.parse(data);
              controller.enqueue(parsed);
            } catch (e) {
              console.error("Error parsing stream data:", e);
            }
          }
        }
      },
    })
  ) as ReadableStream<ChatCompletionStreamResponse>;
};
