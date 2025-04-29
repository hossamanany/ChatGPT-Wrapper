import { Role } from "@/models/role.model";
import { useChatStore } from "@/stores/chat.store";

export async function generateTitle(message: string): Promise<void> {
  const chatStore = useChatStore();
  const openaiModel = import.meta.env.VITE_OPENAI_MODEL || "gpt-3.5-turbo";
  const openaiTemperature =
    Number(import.meta.env.VITE_OPENAI_TEMPERATURE) || 0.5;
  const openaiMaxTokens =
    Number(import.meta.env.VITE_OPENAI_MAX_TOKENS) || 1000;

  if (!chatStore.currentChat || chatStore.currentChat.title) {
    return;
  }

  try {
    console.log("Generating title for message:", message);
    const response = await fetch("/api/chat/title", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        messages: [
          {
            role: Role.user,
            content:
              "Summarize the input as title of no more than 5 words. Output only the summarized title. The input is: " +
              message,
          },
        ],
        model: openaiModel,
        temperature: openaiTemperature,
        max_tokens: openaiMaxTokens,
      }),
    });

    if (!response.ok) {
      const errorData = await response.json();
      console.error("Title generation failed:", errorData);
      throw new Error(errorData.error || "Failed to generate title");
    }

    const data = await response.json();
    console.log(
      "Title generated successfully:",
      data.choices[0].message.content
    );
    await chatStore.setCurrentChatTitle(data.choices[0].message.content);
  } catch (error) {
    console.error("Error generating title:", error);
    throw error;
  }
}
