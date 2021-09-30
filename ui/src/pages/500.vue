<template>
  <n-result size="huge" status="500" title="500 服务器错误" description="">
    <template #footer>
      <n-button type="primary" @click="$router.push('/')">返回首页</n-button>
    </template>
  </n-result>
</template>

<script setup lang="ts">
import { NResult, NButton } from "naive-ui";
</script>  
