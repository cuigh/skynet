<template>
  <n-form :model="model" :rules="rules" ref="form" label-placement="top">
    <n-grid cols="1 640:6" :x-gap="24">
      <n-form-item-gi label="启用" span="6" label-placement="left">
        <n-switch v-model:value="enabled">
          <template #checked>是</template>
          <template #unchecked>否</template>
        </n-switch>
      </n-form-item-gi>
      <n-form-item-gi label="邮件服务器" path="smtp_address" span="2">
        <n-input placeholder="邮件服务器地址" v-model:value="model.smtp_address" />
      </n-form-item-gi>
      <n-form-item-gi label="用户名" path="smtp_username" span="2">
        <n-input placeholder="用户名" v-model:value="model.smtp_username" />
      </n-form-item-gi>
      <n-form-item-gi label="密码" path="smtp_password" span="2">
        <n-input type="password" placeholder="密码" v-model:value="model.smtp_password" />
      </n-form-item-gi>
      <n-form-item-gi label="发送邮箱" path="sender" span="2">
        <n-input placeholder="发送邮箱" v-model:value="model.sender" />
      </n-form-item-gi>
      <n-form-item-gi label="接收邮箱" path="receiver" span="4">
        <n-input placeholder="默认报警邮箱，当任务没有配置维护者时会发送到此邮箱" v-model:value="model.receiver" />
      </n-form-item-gi>
      <n-form-item-gi label="标题模版" path="title" span="6">
        <n-input placeholder="标题模版，支持用 Go 模版语法插入变量" v-model:value="model.title" />
      </n-form-item-gi>
      <n-form-item-gi label="内容模版" path="body" span="6">
        <n-input
          type="textarea"
          rows="5"
          placeholder="内容模版，支持用 Go 模版语法插入变量"
          v-model:value="model.body"
        />
      </n-form-item-gi>
      <n-gi span="6">
        <n-button @click.prevent="submit" type="primary" :disabled="submiting" :loading="submiting">
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
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from "vue";
import {
  NButton,
  NIcon,
  NGrid,
  NGi,
  NInput,
  NForm,
  NFormItemGi,
  NSwitch,
} from "naive-ui";
import { SaveOutline as SaveIcon, } from "@vicons/ionicons5";
import configApi from "@/api/config";
import type { EmailOptions } from "@/api/config";
import { useForm, requiredRule } from "@/utils/form";

const form = ref();
const model = ref({
  enabled: "",
  sender: "",
  smtp_address: "",
  smtp_username: "",
  smtp_password: "",
  title: "",
  body: "",
} as EmailOptions)
const enabled = computed({
  get() { return model.value.enabled === 'true' },
  set(v: boolean) { model.value.enabled = v.toString() },
});
const rules: any = {
  sender: requiredRule(),
  smtp_address: requiredRule(),
  smtp_username: requiredRule(),
  smtp_password: requiredRule(),
};
const { submit, submiting } = useForm(form, () => configApi.save("alert.email", model.value))

async function fetchData() {
  let r = await configApi.find("alert.email");
  model.value = r.data as EmailOptions;
}

onMounted(fetchData);
</script>
