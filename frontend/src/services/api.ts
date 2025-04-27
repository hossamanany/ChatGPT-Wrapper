import type {
  ChatCompletionRequest,
  ChatCompletionResponse,
} from "@/models/api.model";

const API_BASE_URL =
  import.meta.env.VITE_API_URL || "http://localhost:8080/api";

export class ApiService {
  private static async request<T>(
    endpoint: string,
    options: RequestInit = {}
  ): Promise<T> {
    const response = await fetch(`${API_BASE_URL}${endpoint}`, {
      ...options,
      headers: {
        "Content-Type": "application/json",
        ...options.headers,
      },
    });

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error || "An error occurred");
    }

    return response.json();
  }

  static async createChatCompletion(
    request: ChatCompletionRequest
  ): Promise<ChatCompletionResponse> {
    return this.request<ChatCompletionResponse>("/chat", {
      method: "POST",
      body: JSON.stringify(request),
    });
  }

  static async createChatCompletionStream(
    request: ChatCompletionRequest,
    onMessage: (content: string) => void,
    onError: (error: Error) => void,
    onDone: () => void
  ): Promise<void> {
    const response = await fetch(`${API_BASE_URL}/stream`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(request),
    });

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error || "An error occurred");
    }

    const reader = response.body?.getReader();
    if (!reader) {
      throw new Error("Failed to get stream reader");
    }

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
          if (line.startsWith("data: ")) {
            const data = line.slice(6);
            if (data === "[DONE]") {
              onDone();
              return;
            }
            try {
              const parsed = JSON.parse(data);
              if (parsed.choices?.[0]?.delta?.content) {
                onMessage(parsed.choices[0].delta.content);
              }
            } catch (e) {
              console.error("Error parsing stream data:", e);
            }
          }
        }
      }
    } catch (e) {
      onError(e instanceof Error ? e : new Error("Stream error"));
    } finally {
      reader.releaseLock();
    }
  }

  static async generateTitle(
    request: ChatCompletionRequest
  ): Promise<ChatCompletionResponse> {
    return this.request<ChatCompletionResponse>("/title", {
      method: "POST",
      body: JSON.stringify(request),
    });
  }
}
