<template>
  <PageHeader title="作业详情" :subtitle="model.id">
    <template #action>
      <n-popconfirm @positive-click="retry(model.id)" v-if="!model.execute.status">
        <template #trigger>
          <n-button size="small" type="info">重试</n-button>
        </template>
        你确定要重试此作业？
      </n-popconfirm>
      <n-button size="small" @click="$router.push('/jobs')">
        <template #icon>
          <n-icon>
            <back-icon />
          </n-icon>
        </template>返回
      </n-button>
    </template>
  </PageHeader>
  <n-space class="page-body" vertical :size="16">
    <Description cols="1 640:2" label-position="left" label-align="right" :label-width="80">
      <DescriptionItem label="任务">{{ model.task }}</DescriptionItem>
      <DescriptionItem label="处理器">{{ model.handler }}</DescriptionItem>
      <DescriptionItem label="调度器">{{ model.scheduler }}</DescriptionItem>
      <DescriptionItem label="触发时间">
        <n-time :time="model.fire_time" format="yyyy-MM-dd HH:mm:ss" />
      </DescriptionItem>
      <DescriptionItem label="执行方式">
        <n-tag
          size="small"
          round
          :type="model.mode ? 'warning' : 'success'"
        >{{ model.mode ? "手动" : "自动" }}</n-tag>
      </DescriptionItem>
    </Description>
    <Panel title="参数" v-if="model.args && model.args.length">
      <n-table size="small" :single-line="false">
        <thead>
          <tr>
            <th>名称</th>
            <th>值</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="arg in model.args">
            <td>{{ arg.name }}</td>
            <td>{{ arg.value }}</td>
          </tr>
        </tbody>
      </n-table>
    </Panel>
    <Panel title="调度信息" key="dispatch">
      <Description>
        <DescriptionItem label="状态">
          <n-tag
            size="small"
            round
            :type="statusType(model.dispatch.status)"
          >{{ statusText(model.dispatch.status) }}</n-tag>
        </DescriptionItem>
        <DescriptionItem label="时间">
          <n-time :time="model.dispatch.time" format="yyyy-MM-dd HH:mm:ss" />
        </DescriptionItem>
        <DescriptionItem :span="2" label="错误信息" v-if="model.dispatch.error">
          <n-text type="error">{{ model.dispatch.error }}</n-text>
        </DescriptionItem>
      </Description>
    </Panel>
    <Panel title="执行信息" key="execute">
      <Description>
        <DescriptionItem label="状态">
          <n-tag
            size="small"
            round
            :type="statusType(model.execute.status)"
          >{{ statusText(model.execute.status) }}</n-tag>
        </DescriptionItem>
        <DescriptionItem
          label="耗时"
          v-if="model.execute.status"
        >{{ formatDuration(model.execute.end_time - model.execute.start_time) }}</DescriptionItem>
        <DescriptionItem label="开始时间" v-if="model.execute.status">
          <n-time :time="model.execute.start_time" format="y-MM-dd HH:mm:ss.SSS" />
        </DescriptionItem>
        <DescriptionItem label="结束时间" v-if="model.execute.status">
          <n-time :time="model.execute.end_time" format="y-MM-dd HH:mm:ss.SSS" />
        </DescriptionItem>
        <DescriptionItem :span="2" label="错误信息" v-if="model.execute.error">
          <n-text type="error">{{ model.execute.error }}</n-text>
        </DescriptionItem>
      </Description>
    </Panel>
  </n-space>
</template>

<script setup lang="ts">
import { onMounted, ref } from "vue";
import {
  NButton,
  NTag,
  NSpace,
  NIcon,
  NTime,
  NText,
  NTable,
  NPopconfirm,
} from "naive-ui";
import { ArrowBackCircleOutline as BackIcon } from "@vicons/ionicons5";
import PageHeader from "@/components/PageHeader.vue";
import Panel from "@/components/Panel.vue";
import jobApi from "@/api/job";
import type { Job } from "@/api/job";
import { useRoute } from "vue-router";
import { Description, DescriptionItem } from "@/components/description";
import { statusType, statusText } from "./job";
import { formatDuration } from "@/utils/render";

const route = useRoute();
const model = ref({
  args: [] as any,
  dispatch: {},
  execute: {},
} as Job);

async function retry(id: string) {
  await jobApi.retry(id)
  window.message.info("操作成功");
}

async function fetchData() {
  let r = await jobApi.find(route.params.id as string);
  model.value = r.data as Job;
}

onMounted(fetchData);
</script>