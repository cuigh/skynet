<template>
  <n-space class="page-body" vertical :size="12">
    <n-alert type="info">Welcome to Skynet.</n-alert>
    <n-grid cols="1 s:2 m:3 l:4 xl:5 2xl:6" x-gap="12" y-gap="12" responsive="screen">
      <n-gi>
        <n-card size="small" title="任务">
          <n-button text type="primary" @click="$router.push('/tasks')">{{ summary.taskCount }}</n-button>
        </n-card>
      </n-gi>
      <n-gi>
        <n-card size="small" title="作业">
          <n-button text type="primary" @click="$router.push('/jobs')">{{ summary.jobCount }}</n-button>
        </n-card>
      </n-gi>
      <n-gi>
        <n-card size="small" title="用户">
          <n-button
            text
            type="primary"
            @click="$router.push('/account/users')"
          >{{ summary.userCount }}</n-button>
        </n-card>
      </n-gi>
      <n-gi>
        <n-card size="small" title="版本">v{{ summary.version }}({{ summary.goVersion }})</n-card>
      </n-gi>
      <!-- <n-gi>
        <n-statistic label="任务" value="125" />
      </n-gi>-->
    </n-grid>
  </n-space>
</template>

<script setup lang="ts">
import { onMounted, ref } from "vue";
import { NSpace, NAlert, NGrid, NGi, NButton, NStatistic, NCard } from "naive-ui";
import systemApi from "@/api/system";
import type { Summary } from "@/api/system";

const summary = ref({
  taskCount: 0,
  jobCount: 0,
  userCount: 0,
} as Summary)

onMounted(async () => {
  const r = await systemApi.summarize();
  summary.value = r.data as Summary;
});
</script>