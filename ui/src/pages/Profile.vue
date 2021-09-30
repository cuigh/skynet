<template>
  <PageHeader title="个人资料" />
  <div class="page-body">
    <n-form :model="model" :rules="rules" ref="form" label-placement="top">
      <n-grid cols="1 640:2" :x-gap="24">
        <n-form-item-gi label="用户名" path="name">
          <n-input placeholder="用户名" v-model:value="model.name" />
        </n-form-item-gi>
        <n-form-item-gi label="登录名" path="login_name">
          <n-input placeholder="登录名" v-model:value="model.login_name" />
        </n-form-item-gi>
        <n-form-item-gi label="邮箱" path="email">
          <n-input placeholder="电子邮箱" v-model:value="model.email" />
        </n-form-item-gi>
        <n-form-item-gi label="手机" path="phone">
          <n-input placeholder="手机号码" v-model:value="model.phone" />
        </n-form-item-gi>
        <n-gi :span="2">
          <n-button
            :disabled="submiting"
            :loading="submiting"
            @click.prevent="submit"
            type="primary"
          >
            <template #icon>
              <n-icon>
                <save-icon />
              </n-icon>
            </template>
            保存
          </n-button>
        </n-gi>
      </n-grid>
    </n-form>
  </div>
  <n-space class="page-body" vertical :size="12"></n-space>
</template>

<script setup lang="ts">
import { h, onMounted, reactive, ref } from "vue";
import {
  NButton,
  NSpace,
  NIcon,
  NGrid,
  NGi,
  NForm,
  NFormItemGi,
  NInput,
} from "naive-ui";
import {
  SaveOutline as SaveIcon,
} from "@vicons/ionicons5";
import userApi from "@/api/user";
import type { User } from "@/api/user";
import PageHeader from "@/components/PageHeader.vue";
import { useForm, requiredRule } from "@/utils/form";

const form = ref();
const model = ref({} as User);
const rules: any = {
  name: requiredRule(),
  login_name: requiredRule(),
};
const { submit, submiting } = useForm(form, () => userApi.modifyProfile(model.value))

async function fetchData() {
  let user = await userApi.find("");
  model.value = user.data as User || {};
}

onMounted(fetchData);
</script>
