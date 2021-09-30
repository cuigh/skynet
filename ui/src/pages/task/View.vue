<template>
  <PageHeader title="任务详情" :subtitle="model.name">
    <template #action>
      <n-button size="small" @click="$router.push('/tasks')">
        <template #icon>
          <n-icon>
            <back-icon />
          </n-icon>
        </template>返回
      </n-button>
      <n-button size="small" @click="$router.push(`/tasks/${model.name}/edit`)">编辑</n-button>
    </template>
  </PageHeader>
  <n-space class="page-body" vertical :size="16">
    <Description cols="1 640:2" label-position="left" label-align="right" :label-width="75">
      <DescriptionItem label="名称">{{ model.name }}</DescriptionItem>
      <DescriptionItem label="描述">{{ model.desc }}</DescriptionItem>
      <DescriptionItem label="执行器">{{ model.runner }}</DescriptionItem>
      <DescriptionItem label="处理器">{{ model.handler }}</DescriptionItem>
      <DescriptionItem label="状态">
        <n-space :size="6">
          <n-tag
            size="small"
            round
            :type="model.enabled ? 'success' : 'error'"
          >{{ model.enabled ? "启用" : "禁用" }}</n-tag>
        </n-space>
      </DescriptionItem>
      <DescriptionItem label="报警方式">
        <n-space :size="6">
          <n-tag size="small" round type="info" v-for="a in model.alerts">{{ alertText(a) }}</n-tag>
        </n-space>
      </DescriptionItem>
      <DescriptionItem label="维护者" :span="2" v-if="model.maintainers && model.maintainers.length">
        <n-space :size="6">
          <n-button
            text
            size="small"
            type="info"
            @click="$router.push(`/account/users/${u.id}`)"
            v-for="u in maintainers"
          >{{ u.name }}</n-button>
        </n-space>
      </DescriptionItem>
    </Description>
    <Panel title="触发器">
      <n-space :size="6">
        <n-tag round v-for="t in model.triggers">{{ t }}</n-tag>
      </n-space>
    </Panel>
    <Panel title="参数" v-if="model.args">
      <n-table size="small" :bordered="true" :single-line="true">
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
  </n-space>
</template>

<script setup lang="ts">
import { onMounted, ref } from "vue";
import {
  NButton,
  NTag,
  NSpace,
  NIcon,
  NTable,
} from "naive-ui";
import { ArrowBackCircleOutline as BackIcon } from "@vicons/ionicons5";
import PageHeader from "@/components/PageHeader.vue";
import taskApi from "@/api/task";
import userApi from "@/api/user";
import type { Task } from "@/api/task";
import { useRoute } from "vue-router";
import Panel from "@/components/Panel.vue";
import { Description, DescriptionItem } from "@/components/description";
import { alertText } from "./task";

const route = useRoute();
const model = ref({} as Task);
const maintainers = ref();

async function fetchData() {
  let tr = await taskApi.find(route.params.name as string);
  model.value = tr.data as Task;
  if (model.value.maintainers && model.value.maintainers.length) {
    let ur = await userApi.fetch(model.value.maintainers as string[]);
    maintainers.value = ur.data;
  }
}

onMounted(fetchData);
</script>