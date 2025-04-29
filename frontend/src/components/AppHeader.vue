<script setup lang="ts">
import { useAppStore } from "@/stores/app.store";
import { Bars3Icon } from "@heroicons/vue/24/outline";
import { QuestionMarkCircleIcon } from "@heroicons/vue/24/outline";

const appStore = useAppStore();
</script>

<template>
  <div class="p-2 bg-gray-200 flex justify-between">
    <div title="Toggle sidebar">
      <Bars3Icon
        class="h-6 w-6 text-gray-500 cursor-pointer hover-blue"
        @click="appStore.toggleSidebar()"
      ></Bars3Icon>
    </div>
    <div class="flex">
      <div title="About">
        <QuestionMarkCircleIcon
          class="h-6 w-6 text-gray-500 cursor-pointer hover-blue"
          @click="appStore.showAbout()"
        ></QuestionMarkCircleIcon>
      </div>
    </div>
  </div>
</template>
