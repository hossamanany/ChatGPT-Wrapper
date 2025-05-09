<script setup lang="ts">
import { computed, onMounted, ref } from "vue";
import { PlayIcon } from "@heroicons/vue/24/outline";
import { Role } from "@/models/role.model";
import { useChatStore } from "@/stores/chat.store";
import MarkdownIt from "markdown-it";
import hljs from "highlight.js";
import { FwbAlert, FwbButton, FwbSpinner } from "flowbite-vue";
import { useAppStore } from "@/stores/app.store";
import { generateTitle } from "@/services/title.service";
import { streamChatResponse } from "@/services/chat.service";
import { watch } from "vue";

const input = ref("");
const inputTextarea = ref<HTMLTextAreaElement | null>(null);
const scrollingDiv = ref<HTMLElement | null>(null);
const userScrolled = ref(false);
const pending = ref(false);

const appStore = useAppStore();
const chatStore = useChatStore();

const md = new MarkdownIt({
  breaks: true,
  linkify: true,
  highlight: function (str, lang) {
    if (lang && hljs.getLanguage(lang)) {
      try {
        return hljs.highlight(str, { language: lang }).value;
      } catch (e) {
        console.log(e);
      }
    }
    return "";
  },
});

const isInputEnabled = computed(() => !pending.value);
const isSendBtnEnabled = computed(() => input.value?.trim().length > 0);

onMounted(() => {
  setTimeout(() => inputTextarea.value?.focus(), 100);
});

watch(input, () => {
  if (inputTextarea.value) {
    inputTextarea.value.style.height = "auto";
    const borderOffset = 2;
    inputTextarea.value.style.height =
      inputTextarea.value.scrollHeight + borderOffset + "px";
  }
});

async function onSend() {
  pending.value = true;
  try {
    userScrolled.value = false;
    inputTextarea.value?.blur();
    await chatStore.addMessage({ role: Role.user, content: input.value });
    autoScrollDown();

    await sendRequestForTitle(input.value);

    input.value = "";
    await sendRequestForResponse();
  } catch (e) {
    if (e instanceof Error) {
      appStore.addError(e.message);
    }
  }
  pending.value = false;
}

async function sendRequestForTitle(message: string) {
  try {
    await generateTitle(message);
  } catch (e) {
    if (e instanceof Error) {
      appStore.addError(e.message);
    }
  }
}

async function sendRequestForResponse() {
  try {
    await streamChatResponse();
    autoScrollDown();
  } catch (e) {
    if (e instanceof Error) {
      appStore.addError(e.message);
    }
  }
}

function autoScrollDown() {
  if (scrollingDiv.value && !userScrolled.value) {
    scrollingDiv.value.scrollTop = scrollingDiv.value.scrollHeight;
  }
}

function checkIfUserScrolled() {
  if (scrollingDiv.value) {
    userScrolled.value =
      scrollingDiv.value.scrollTop + scrollingDiv.value.clientHeight !==
      scrollingDiv.value.scrollHeight;
  }
}
</script>

<template>
  <div class="flex flex-1 flex-col overflow-auto">
    <fwb-alert
      closable
      type="danger"
      class="mt-4 ml-4 mr-4 gap-0"
      v-for="error in appStore.errors"
      :key="error.id"
      @close="appStore.removeError(error.id)"
    >
      {{ error.message }}
    </fwb-alert>
    <main
      class="flex-1 p-4 overflow-auto"
      ref="scrollingDiv"
      @scroll="checkIfUserScrolled()"
    >
      <template v-if="chatStore.currentChat">
        <template
          v-for="(message, index) in chatStore.currentChat.messages"
          :key="index"
        >
          <template v-if="message.content && message.role === Role.user">
            <div class="flex">
              <div
                class="border-green-600 border-2 border-solid py-2 px-3 rounded mb-4 message-content"
                v-html="md.render(message.content)"
              />
            </div>
          </template>
          <template v-if="message.content && message.role === Role.assistant">
            <div class="flex">
              <div
                class="py-2 px-3 rounded mb-4 ml-5 message-content"
                v-html="md.render(message.content)"
              />
            </div>
          </template>
        </template>
      </template>
    </main>
    <div class="flex w-full p-4">
      <textarea
        class="p-2 overflow-x-hidden w-full text-gray-900 bg-gray-50 rounded border border-gray-300 focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500 resize-none max-h-[300px] overflow-y-auto box-border"
        :rows="1"
        :placeholder="pending ? 'Answering...' : 'Chat with AI...'"
        ref="inputTextarea"
        v-model="input"
        @keydown.enter="
          !$event.shiftKey && isInputEnabled && isSendBtnEnabled
            ? ($event.preventDefault(), onSend())
            : null
        "
        :disabled="!isInputEnabled"
      />
      <fwb-button
        color="default"
        @click="onSend"
        :disabled="!isSendBtnEnabled"
        class="ml-2 p-2 rounded"
      >
        <PlayIcon class="h-6 w-6" v-if="!pending"></PlayIcon>
        <fwb-spinner size="6" v-if="pending" />
      </fwb-button>
    </div>
  </div>
</template>

<style>
@import "../../node_modules/highlight.js/styles/github.css";

.message-content {
  pre:not(:last-child),
  p:not(:last-child),
  ol:not(:last-child),
  ul:not(:last-child),
  li:not(:last-child),
  table:not(:last-child),
  blockquote:not(:last-child),
  hr:not(:last-child),
  h1,
  h2,
  h3,
  h4,
  h5,
  h6 {
    margin-bottom: 0.5rem;
  }

  blockquote {
    margin-left: 1rem;
    font-style: italic;
  }

  h1 {
    font-size: 1.5rem;
  }

  h2 {
    font-size: 1.25rem;
  }

  h3 {
    font-size: 1.125rem;
  }

  h1,
  h2,
  h3,
  h4,
  h5,
  h6 {
    font-weight: bold;
    margin-top: 1rem;
  }

  pre {
    margin-left: 1rem;
    background-color: rgb(249 250 251);
    display: table;
    border-radius: 5px;
    padding: 0 5px;
    white-space: pre-wrap;
  }

  code:not(pre code) {
    background-color: rgb(249 250 251);
    border-radius: 5px;
    padding: 0 1px;
  }

  a {
    color: rgb(37 99 235);
  }

  ul {
    list-style-type: disc;
    margin-left: 2rem;
  }

  ol {
    list-style-type: decimal;
    margin-left: 2rem;
  }

  td,
  th {
    border: 1px solid black;
    padding: 5px;
  }
}
</style>
