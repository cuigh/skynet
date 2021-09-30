<template>
  <PageHeader title="通知设置" subtitle="任务执行出错时会触发报警" />
  <n-space class="page-body" vertical :size="12">
    <n-tabs type="line" tab-style="padding-top: 0">
      <n-tab-pane name="email" tab="邮件">
        <Email />
      </n-tab-pane>
      <n-tab-pane name="wecom" tab="企业微信">
        <Wecom />
      </n-tab-pane>
    </n-tabs>
    <n-alert
      type="info"
      :show-icon="false"
    >可以在模版中使用变量，有效变量如下：error, task, handler, runner, job, mode, fire, args, scheduler, duration.</n-alert>
  </n-space>
</template>

<script setup lang="ts">
import {
  NSpace,
  NAlert,
  NTabs,
  NTabPane,
} from "naive-ui";
import PageHeader from "@/components/PageHeader.vue";
import Email from "./modules/Email.vue";
import Wecom from "./modules/Wecom.vue";
</script>
