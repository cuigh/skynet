<template>
  <PageHeader title="角色详情" :subtitle="model.name">
    <template #action>
      <n-button size="small" @click="$router.push({ name: 'role.list' })">
        <template #icon>
          <n-icon>
            <back-icon />
          </n-icon>
        </template>返回
      </n-button>
      <n-button
        size="small"
        @click="$router.push({ name: 'role.edit', params: { id: model.id } })"
      >编辑</n-button>
    </template>
  </PageHeader>
  <n-space class="page-body" vertical :size="16">
    <Description cols="1 640:2" label-position="left" label-align="right" :label-width="60">
      <DescriptionItem label="ID">{{ model.id }}</DescriptionItem>
      <DescriptionItem label="名称">{{ model.name }}</DescriptionItem>
      <DescriptionItem :span="2" label="说明">{{ model.desc }}</DescriptionItem>
      <DescriptionItem :span="2" label="权限">
        <n-space :size="6">
          <n-tag round size="small" type="info" v-for="r in model.perms">{{ r }}</n-tag>
        </n-space>
      </DescriptionItem>
    </Description>
  </n-space>
</template>

<script setup lang="ts">
import { onMounted, ref } from "vue";
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
import roleApi from "@/api/role";
import type { Role } from "@/api/role";

const route = useRoute();
const model = ref({} as Role);

async function fetchData() {
  let r = await roleApi.find(route.params.id as string);
  model.value = r.data as Role;
}

onMounted(fetchData);
</script>