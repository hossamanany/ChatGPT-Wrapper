<script setup lang="ts">
import { FwbButton, FwbModal } from "flowbite-vue";
import { useAppStore } from "@/stores/app.store";

const appStore = useAppStore();
</script>

<template>
  <fwb-modal
    v-if="appStore.isAboutVisible"
    @close="appStore.hideAbout()"
    not-escapable
    persistent
  >
    <template #header>
      <div class="text-lg">About</div>
    </template>
    <template #body>
      <div class="flex">
        <div><img src="/icon-512.png" alt="Logo" /></div>
        <div class="ml-5">
          <pre class="whitespace-pre-wrap text-xs">
<b>ChatGPT-Wrapper 1.0</b>

Source, info and contact:
<a href="https://github.com/hossamanany/ChatGPT-Wrapper" class="hover-blue">https://github.com/hossamanany/ChatGPT-Wrapper</a>

License:
MIT
        </pre>
        </div>
      </div>
      <pre class="whitespace-pre-wrap text-xs">

<b>Disclaimer</b>

The ChatGPT-Wrapper ("Application") and all associated content are provided on an "as is" and "as available" basis without any warranties of any kind, either express or implied, including, but not limited to, the implied warranties of merchantability, fitness for a particular purpose, or non-infringement.

The developers, contributors, and owners of ChatGPT-Wrapper take no responsibility and assume no liability for any content that users create, upload, post, send, receive, or store on or through our Application. You understand and agree that you may be exposed to content that might be offensive, illegal, misleading, or otherwise inappropriate, none of which ChatGPT-Wrapper's developers, contributors, and owners will be responsible for.

<b>Limitation of Liability</b>

To the fullest extent permitted by applicable law, in no event will ChatGPT-Wrapper, its affiliates, officers, directors, employees, agents, suppliers, or licensors be liable for (a) any indirect, special, incidental, punitive, exemplary, or consequential damages; (b) any loss of use, data, business, or profits; regardless of legal theory, whether or not ChatGPT-Wrapper has been warned of the possibility of such damages, and even if a remedy fails of its essential purpose.

<b>User Responsibility</b>

By using ChatGPT-Wrapper, you agree and acknowledge that you are solely responsible for your conduct while using our Application and for any consequences thereof. You agree to use ChatGPT-Wrapper only for purposes that are legal, proper, and in accordance with these Terms and any applicable laws or regulations.

ChatGPT-Wrapper does not endorse, support, represent or guarantee the completeness, truthfulness, accuracy, or reliability of any content or communications posted using our Application. You understand that by using ChatGPT-Wrapper, you may be exposed to content that might be inaccurate, offensive, indecent, or objectionable. You agree to waive, and hereby do waive, any legal or equitable rights or remedies you have or may have against ChatGPT-Wrapper with respect thereto.

<b>Amendments to this Disclaimer</b>

ChatGPT-Wrapper reserves the right to modify this disclaimer at any time. You should review this disclaimer periodically for changes. Your continued use of the Application after the posting of changes constitutes your binding acceptance of such changes.

This disclaimer constitutes the entire agreement between you and ChatGPT-Wrapper regarding the use of the Application, superseding any prior agreements between you and ChatGPT-Wrapper regarding its use.</pre>
    </template>
    <template #footer>
      <div class="modal-footer">
        <fwb-button
          @click="appStore.hideAbout()"
          color="alternative"
          class="modal-button"
          >Close</fwb-button
        >
      </div>
    </template>
  </fwb-modal>
</template>

<style scoped>
img {
  width: 110px;
}
</style>
