<template>
  <PageHeader title="角色列表">
    <template #action>
      <n-button size="small" @click="$router.push({ name: 'role.new' })">
        <template #icon>
          <n-icon>
            <add-icon />
          </n-icon>
        </template>添加
      </n-button>
    </template>
  </PageHeader>
  <n-space class="page-body" vertical :size="12">
    <n-space :size="12">
      <n-input size="small" v-model:value="model.name" placeholder="角色名" clearable />
      <n-button size="small" type="primary" @click="fetchData">查询</n-button>
    </n-space>
    <n-table size="small" :bordered="true" :single-line="false">
      <thead>
        <tr>
          <th>ID</th>
          <th>名称</th>
          <th>说明</th>
          <th>操作</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="(r, index) of model.roles" :key="r.id">
          <td>
            <anchor :href="`/account/roles/${r.id}`">{{ r.id }}</anchor>
          </td>
          <td>{{ r.name }}</td>
          <td>{{ r.desc }}</td>
          <td>
            <n-space :size="4">
              <n-popconfirm :show-icon="false" @positive-click="deleteRole(r.id, index)">
                <template #trigger>
                  <n-button size="tiny" ghost type="error">删除</n-button>
                </template>
                确定删除此角色？
              </n-popconfirm>
              <n-button
                size="tiny"
                ghost
                type="warning"
                @click="$router.push({ name: 'role.edit', params: { id: r.id } })"
              >编辑</n-button>
            </n-space>
          </td>
        </tr>
      </tbody>
    </n-table>
  </n-space>
</template>

<script setup lang="ts">
import { onMounted, reactive } from "vue";
import {
  NSpace,
  NInput,
  NButton,
  NIcon,
  NTable,
  NPopconfirm,
} from "naive-ui";
import {
  AddOutline as AddIcon,
} from "@vicons/ionicons5";
import Anchor from "@/components/Anchor.vue";
import PageHeader from "@/components/PageHeader.vue";
import roleApi from "@/api/role";
import type { Role } from "@/api/role";

const model = reactive({
  name: "",
  roles: [] as Role[],
});

async function deleteRole(id: string, index: number) {
  await roleApi.delete(id);
  model.roles.splice(index, 1)
}

async function fetchData() {
  let r = await roleApi.search(model.name);
  model.roles = r.data as Role[];
}

onMounted(fetchData);
</script>