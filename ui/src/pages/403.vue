<template>
  <n-result size="huge" status="403" title="403 禁止访问" description="很抱歉，你没有权限访问此页面">
    <template #footer>
      <n-button type="primary" @click="$router.push('/')">返回首页</n-button>
    </template>
  </n-result>
</template>

<script setup lang="ts">
import { NResult, NButton } from "naive-ui";
</script>
