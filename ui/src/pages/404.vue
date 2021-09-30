<template>
  <n-result size="huge" status="404" title="404 资源不存在" description="您要请求的页面已被移动、删除、重命名或可能不存在">
    <template #footer>
      <n-button type="primary" @click="$router.push('/')">返回首页</n-button>
    </template>
  </n-result>
</template>

<script setup lang="ts">
import { NResult, NButton } from "naive-ui";
</script>
  