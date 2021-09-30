<template>
  <PageHeader title="作业列表"></PageHeader>
  <n-space class="page-body" vertical :size="12">
    <n-space :size="12">
      <n-input size="small" v-model:value="filter.task" placeholder="任务名称" clearable />
      <n-select
        size="small"
        placeholder="执行方式"
        v-model:value="filter.mode"
        :options="modeOptions"
        style="width: 120px"
        clearable
      />
      <n-select
        size="small"
        placeholder="调度状态"
        v-model:value="filter.dispatch_status"
        :options="statusOptions"
        style="width: 120px"
        clearable
      />
      <n-select
        size="small"
        placeholder="执行状态"
        v-model:value="filter.execute_status"
        :options="statusOptions"
        style="width: 120px"
        clearable
      />
      <n-button size="small" type="primary" @click="() => fetchData()">查询</n-button>
    </n-space>
    <n-data-table
      remote
      size="small"
      :columns="columns"
      :data="state.data"
      :pagination="pagination"
      :loading="state.loading"
      :row-key="key"
      @update:page="fetchData"
      scroll-x="max-content"
    />
  </n-space>
</template>

<script setup lang="ts">
import { reactive } from "vue";
import {
  NButton,
  NSpace,
  NDataTable,
  NInput,
  NSelect,
} from "naive-ui";
import PageHeader from "@/components/PageHeader.vue";
import jobApi from "@/api/job";
import type { Job } from "@/api/job";
import { renderLink, renderTag, renderTime, formatDuration } from "@/utils/render";
import { statusType, statusText } from "./job";
import { useDataTable } from "@/utils/data-table";

const modeOptions = [{ label: '自动', value: 0 }, { label: '手动', value: 1 }];
const statusOptions = [
  { value: 0, label: statusText(0) },
  { value: 1, label: statusText(1) },
  { value: 2, label: statusText(2) },
];
const filter = reactive({
  task: "",
  mode: undefined,
  dispatch_status: undefined,
  execute_status: undefined,
});
const columns = [
  {
    title: "ID",
    key: "id",
    width: 100,
    fixed: 'left' as const,
    render: (row: Job) => renderLink(`/jobs/${row.id}`, row.id.substr(0, 8)),
  },
  {
    title: "任务",
    key: "task",
  },
  {
    title: "执行方式",
    key: "mode",
    render: (row: Job) =>
      renderTag(
        row.mode ? "手动" : "自动",
        row.mode ? "warning" : "success"
      ),
  },
  {
    title: "调度状态",
    key: "dispatch_status",
    render: (row: Job) => renderTag(statusText(row.dispatch.status), statusType(row.dispatch.status)),
  },
  {
    title: "执行状态",
    key: "execute_status",
    render: (row: Job) => renderTag(statusText(row.execute.status), statusType(row.execute.status)),
  },
  {
    title: "触发时间",
    key: "fire_time",
    width: 160,
    render: (row: Job) => renderTime(row.fire_time),
  },
  {
    title: "耗时",
    key: "duration",
    render: (row: Job) => row.execute.end_time ? formatDuration(row.execute.end_time - row.execute.start_time) : '',
  },
];
const key = (row: Job) => row.id
const { state, pagination, fetchData } = useDataTable(jobApi.search, filter)
</script>