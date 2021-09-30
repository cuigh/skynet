<template>
  <n-config-provider
    :theme="theme"
    :locale="zhCN"
    :date-locale="dateZhCN"
    :theme-overrides="themeFixed"
  >
    <n-loading-bar-provider>
      <n-message-provider>
        <n-notification-provider>
          <n-dialog-provider>
            <root />
          </n-dialog-provider>
        </n-notification-provider>
      </n-message-provider>
    </n-loading-bar-provider>
    <n-global-style />
  </n-config-provider>
</template>

<script lang="ts">
import { computed, defineComponent, onMounted } from "vue";
import {
  zhCN,
  dateZhCN,
  NConfigProvider,
  NDialogProvider,
  NNotificationProvider,
  NMessageProvider,
  NLoadingBarProvider,
  NGlobalStyle,
  useMessage,
  useLoadingBar,
useDialog,
useNotification,
} from "naive-ui";
import { darkTheme } from "naive-ui";
import { useRoute } from "vue-router";
import { useStore } from "vuex";
import { initLoadingBar } from "@/router/router";
import DefaultLayout from "./layouts/Default.vue";
import SimpleLayout from "./layouts/Simple.vue";
import EmptyLayout from "./layouts/Empty.vue";

const Root = defineComponent({
  name: "App",
  components: {
    DefaultLayout,
    SimpleLayout,
    EmptyLayout,
  },
  template: '<component :is="layout"></component>',
  setup() {
    window.message = useMessage();
    window.dialog = useDialog();
    window.notification = useNotification();
    initLoadingBar(useLoadingBar());

    const route = useRoute();
    return {
      layout: computed(() => (route.meta.layout || "default") + "-layout"),
    };
  },
})

export default defineComponent({
  name: "App",
  components: {
    NConfigProvider,
    NDialogProvider,
    NNotificationProvider,
    NMessageProvider,
    NLoadingBarProvider,
    NGlobalStyle,
    Root,
  },
  setup() {
    const store = useStore();
    const theme = computed(() => store.state.theme === "dark" ? darkTheme : null);
    const themeFixed = {
      "Form": {
        "feedbackHeightMedium": "20px",
        "feedbackFontSizeMedium": "12px",
        // "blankHeightMedium": "30px",
      }
    }
    return {
      zhCN,
      dateZhCN,
      theme,
      themeFixed,
    };
  },
});
</script>

<style>
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}
.n-layout-header {
  padding: 0 16px;
  height: 56px;
  display: flex;
  display: -webkit-flex;
  align-items: center;
}
.n-layout-footer {
  text-align: center;
  line-height: 40px;
  height: 64px;
  padding: 12px;
}
.logo {
  cursor: pointer;
  display: flex;
  align-items: center;
  font-size: 18px;
}
.logo > img {
  margin-right: 6px;
  height: 32px;
  width: 32px;
}
.page-body {
  padding: 16px;
}
</style>