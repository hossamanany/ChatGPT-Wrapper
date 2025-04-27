import "./assets/main.css";
import "../node_modules/flowbite-vue/dist/index.css";
import "./assets/styles/common.css";

import { createApp } from "vue";
import { createPinia } from "pinia";
import "flowbite-vue";

import App from "./App.vue";

const app = createApp(App);

app.use(createPinia());

app.mount("#app");
