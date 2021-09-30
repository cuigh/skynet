<template>
  <n-layout :position="isMobile ? 'static' : 'absolute'">
    <n-layout-header bordered>
      <n-text tag="div" class="logo" :depth="1" @click="$router.push('/')">
        <img src="../assets/skynet.svg" />
        <span>Skynet</span>
      </n-text>
    </n-layout-header>
    <n-layout-content
      :position="isMobile ? 'static' : 'absolute'"
      :style="isMobile ? '' : 'top: 64px; bottom: 64px'"
    >
      <div class="content">
        <router-view></router-view>
      </div>
    </n-layout-content>
    <n-layout-footer bordered :position="isMobile ? 'static' : 'absolute'">
      <span>© 2021 cuigh. 保留所有权利。</span>
    </n-layout-footer>
  </n-layout>
</template>

<script setup lang="ts">
import {
  NText,
  NLayout,
  NLayoutHeader,
  NLayoutContent,
  NLayoutFooter,
} from "naive-ui";
import { RouterView } from "vue-router";
import { useIsMobile, useIsTablet } from "@/utils";

const isMobile = useIsMobile()
</script>

<style scoped>
.content {
  padding: 48px 24px;
  overflow-x: scroll;
}
</style>
