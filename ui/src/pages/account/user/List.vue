<template>
  <PageHeader title="用户列表">
    <template #action>
      <n-button size="small" @click="$router.push('/account/users/new')">
        <template #icon>
          <n-icon>
            <add-icon />
          </n-icon>
        </template>添加
      </n-button>
    </template>
  </PageHeader>
  <n-space class="page-body" vertical :size="12">
    <tab>
      <tab-pane href="/account/users" :active="!$route.query.filter">All</tab-pane>
      <tab-pane
        href="/account/users?filter=admins"
        :active="$route.query.filter === 'admins'"
      >Admins</tab-pane>
      <tab-pane
        href="/account/users?filter=active"
        :active="$route.query.filter === 'active'"
      >Active</tab-pane>
      <tab-pane
        href="/account/users?filter=blocked"
        :active="$route.query.filter === 'blocked'"
      >Blocked</tab-pane>
    </tab>
    <n-space :size="12">
      <n-input size="small" v-model:value="args.name" placeholder="用户名" clearable />
      <n-input size="small" v-model:value="args.login_name" placeholder="登录名" clearable />
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
</template>

<script setup lang="ts">
import { reactive, watch } from "vue";
import {
  NSpace,
  NInput,
  NButton,
  NIcon,
  NDataTable,
} from "naive-ui";
import {
  AddOutline as AddIcon,
} from "@vicons/ionicons5";
import { useRoute, useRouter } from "vue-router";
import PageHeader from "@/components/PageHeader.vue";
import { Tab, TabPane } from "@/components/tab";
import userApi from "@/api/user";
import type { User } from "@/api/user";
import { useDataTable } from "@/utils/data-table";
import { renderButtons, renderLink, renderTag } from "@/utils/render";

const route = useRoute();
const router = useRouter();
const args = reactive({
  name: "",
  login_name: "",
});
const columns = [
  {
    title: "ID",
    key: "id",
    render: (row: User) => renderLink(`/account/users/${row.id}`, row.id),
  },
  {
    title: "用户名",
    key: "name",
  },
  {
    title: "登录名",
    key: "login_name",
  },
  {
    title: "邮箱",
    key: "email",
  },
  {
    title: "手机",
    key: "phone",
  },
  {
    title: "管理员",
    key: "admin",
    render: (row: User) => row.admin ? '是' : '否',
  },
  {
    title: "状态",
    key: "status",
    render: (row: User) => renderTag(
      row.status ? "正常" : "禁用",
      row.status ? "success" : "warning"
    ),
  },
  {
    title: "操作",
    key: "actions",
    render(row: User, index: number) {
      return renderButtons([
        row.status ?
          { type: 'warning', text: '禁用', action: () => setStatus(row, 0), prompt: '确定禁用此用户？' } :
          { type: 'success', text: '启用', action: () => setStatus(row, 1) },
        { type: 'warning', text: '编辑', action: () => router.push(`/account/users/${row.id}/edit`) }
      ])
    },
  },
];
const { state, pagination, fetchData } = useDataTable(userApi.search, () => {
  return { ...args, filter: route.query.filter }
})

async function setStatus(u: User, status: number) {
  await userApi.setStatus({ id: u.id, status });
  u.status = status
}

watch(() => route.query.filter, (newValue: any, oldValue: any) => {
  fetchData()
})
</script>