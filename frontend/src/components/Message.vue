<script setup lang="ts">
import { Role } from "@/models/role.model";
import MarkdownIt from "markdown-it";
import hljs from "highlight.js";

const props = defineProps<{
  content: string;
  role: Role;
}>();

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
</script>

<template>
  <div class="flex">
    <div
      :class="[
        'py-2 px-3 rounded mb-4 message-content',
        role === Role.user ? 'border-green-600 border-2 border-solid' : 'ml-5',
      ]"
      v-html="md.render(content)"
    />
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
