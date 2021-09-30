<template>
  <PageHeader title="用户详情" :subtitle="model.user.name">
    <template #action>
      <n-button size="small" @click="$router.push('/account/users')">
        <template #icon>
          <n-icon>
            <back-icon />
          </n-icon>
        </template>返回
      </n-button>
      <n-button size="small" @click="$router.push(`/account/users/${model.user.id}/edit`)">编辑</n-button>
    </template>
  </PageHeader>
  <n-space class="page-body" vertical :size="16">
    <Description cols="1 640:2" label-position="left" label-align="right" :label-width="75">
      <DescriptionItem label="ID">{{ model.user.id }}</DescriptionItem>
      <DescriptionItem label="状态">
        <n-tag
          round
          size="small"
          :type="model.user.status ? 'primary' : 'warning'"
        >{{ model.user.status ? '正常' : '禁用' }}</n-tag>
      </DescriptionItem>
      <DescriptionItem label="用户名">{{ model.user.name }}</DescriptionItem>
      <DescriptionItem label="登录名">{{ model.user.login_name }}</DescriptionItem>
      <DescriptionItem label="邮箱">{{ model.user.email }}</DescriptionItem>
      <DescriptionItem label="手机">{{ model.user.phone }}</DescriptionItem>
      <DescriptionItem label="企业微信">{{ model.user.wecom }}</DescriptionItem>
      <DescriptionItem label="管理员">
        <n-tag
          size="small"
          round
          :type="model.user.admin ? 'success' : 'default'"
        >{{ model.user.admin ? "是" : "否" }}</n-tag>
      </DescriptionItem>
      <DescriptionItem label="角色">
        <n-space :size="6">
          <n-tag
            round
            size="small"
            type="info"
            v-for="r in model.user.roles"
          >{{ model.roles.get(r) }}</n-tag>
        </n-space>
      </DescriptionItem>
    </Description>
  </n-space>
</template>

<script setup lang="ts">
import { onMounted, reactive } from "vue";
import {
  NButton,
  NTag,
  NSpace,
  NIcon,
} from "naive-ui";
import { useRoute } from "vue-router";
import { ArrowBackCircleOutline as BackIcon } from "@vicons/ionicons5";
import PageHeader from "@/components/PageHeader.vue";
import { Description, DescriptionItem } from "@/components/description";
import userApi from "@/api/user";
import roleApi from "@/api/role";
import type { User } from "@/api/user";

const route = useRoute();
const model = reactive({
  user: {} as User,
  roles: new Map<string, string>(),
});

async function fetchData() {
  let user = await userApi.find(route.params.id as string);
  let roles = await roleApi.search()
  model.user = user.data as User
  roles.data?.forEach(r => model.roles.set(r.id, r.name))
}

onMounted(fetchData);
</script>