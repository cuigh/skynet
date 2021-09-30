<template>
  <page-header title="任务列表">
    <template #action>
      <n-button size="small" @click="$router.push('/tasks/new')">
        <template #icon>
          <n-icon>
            <add-icon />
          </n-icon>
        </template>新建
      </n-button>
    </template>
  </page-header>
  <n-space class="page-body" vertical :size="12">
    <n-space :size="12">
      <n-input size="small" v-model:value="filter.name" placeholder="名称" clearable />
      <n-input size="small" v-model:value="filter.runner" placeholder="执行器" clearable />
      <n-button size="small" type="primary" @click="() => fetchData()">查询</n-button>
    </n-space>
    <n-data-table
      remote
      :row-key="row => row.name"
      size="small"
      :columns="columns"
      :data="state.data"
      :pagination="pagination"
      :loading="state.loading"
      @update:page="fetchData"
      scroll-x="max-content"
    />
  </n-space>
  <n-modal
    preset="card"
    size="small"
    :title="`执行任务: ${execModel.name}`"
    style="width: 500px"
    v-model:show="showModal"
  >
    <n-form :model="execModel" ref="form">
      <n-form-item path="name" label="名称">
        <n-input v-model:value="execModel.name" disabled />
      </n-form-item>
      <n-form-item path="args" label="参数">
        <n-dynamic-input v-model:value="execModel.args" #="{ index, value }" :on-create="newArg">
          <n-input placeholder="参数名" v-model:value="value.name" />
          <div style="height: 34px; line-height: 34px; margin: 0 8px">=</div>
          <n-input placeholder="参数值" v-model:value="value.value" />
        </n-dynamic-input>
      </n-form-item>
    </n-form>
    <template #footer>
      <n-button round type="primary" @click.prevent="submit" :disabled="submiting">确定</n-button>
    </template>
  </n-modal>
</template>

<script setup lang="ts">
import { reactive, ref, h } from "vue";
import {
  NButton,
  NSpace,
  NDataTable,
  NInput,
  NIcon,
  NModal,
  NDynamicInput,
  NForm,
  NFormItem,
} from "naive-ui";
import { AddOutline as AddIcon } from "@vicons/ionicons5";
import PageHeader from "@/components/PageHeader.vue";
import { renderButtons, renderLink, renderTag } from "@/utils/render";
import { useRouter } from "vue-router";
import taskApi from "@/api/task";
import type { Task, ExecuteArgs } from "@/api/task";
import { useDataTable } from "@/utils/data-table";
import { useForm } from "@/utils/form";

const router = useRouter();
const filter = reactive({
  name: "",
  runner: "",
});
const showModal = ref(false)
const execModel = reactive({} as ExecuteArgs);
const columns = [
  {
    title: "名称",
    key: "name",
    fixed: "left" as const,
    render: (t: Task) => renderLink(`/tasks/${t.name}`, t.name),
  },
  {
    title: "执行器",
    key: "runner"
  },
  {
    title: "处理器",
    key: "handler"
  },
  {
    title: "触发器",
    key: "triggers",
    render: (t: Task) => h(NSpace, { vertical: true }, { default: () => t.triggers.map(c => renderTag(c)) }),
  },
  {
    title: "描述",
    key: "desc"
  },
  {
    title: "状态",
    key: "enabled",
    render: (t: Task) => renderTag(t.enabled ? "启用" : "禁用", t.enabled ? "success" : "error"),
  },
  {
    title: "操作",
    key: "actions",
    render(t: Task, index: number) {
      return renderButtons([
        {
          type: 'info',
          text: '执行',
          action: () => {
            execModel.name = t.name
            execModel.args = []
            showModal.value = true
          }
        },
        {
          type: 'error',
          text: '删除',
          action: () => deleteTask(t, index),
          prompt: '你确定要删除此任务？'
        },
        {
          type: 'warning',
          text: '编辑',
          action: () => router.push(`/tasks/${t.name}/edit`),
        },
      ])
    },
  },
];
const { state, pagination, fetchData } = useDataTable(taskApi.search, filter)
const form = ref();
const { submit, submiting } = useForm(form, () => taskApi.execute(execModel), () => {
  window.message.info("任务执行成功");
  showModal.value = false;
})

function newArg() {
  return {
    name: '',
    value: ''
  }
}

async function deleteTask(row: Task, index: number) {
  await taskApi.delete(row.name)
  state.data.splice(index, 1)
}
</script>