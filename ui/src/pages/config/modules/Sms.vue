<template>
  <n-form :model="model" :rules="rules" ref="form" label-placement="top">
    <n-grid cols="1 640:6" :x-gap="24">
      <n-form-item-gi label="启用" span="3" label-placement="left">
        <n-switch v-model:value="enabled">
          <template #checked>是</template>
          <template #unchecked>否</template>
        </n-switch>
      </n-form-item-gi>
      <n-form-item-gi label="服务商" path="provider" span="3" label-placement="left">
        <n-radio-group v-model:value="model.provider" name="mode" size="small">
          <n-space>
            <n-radio value="aliyun">阿里云</n-radio>
          </n-space>
        </n-radio-group>
      </n-form-item-gi>
      <n-form-item-gi label="访问凭证" path="key" span="3">
        <n-input placeholder="访问凭证(Key)" v-model:value="model.key" />
      </n-form-item-gi>
      <n-form-item-gi label="访问密钥" path="secret" span="3">
        <n-input placeholder="访问密钥(Secret)" v-model:value="model.secret" />
      </n-form-item-gi>
      <n-form-item-gi label="短信签名" path="sign" span="3">
        <n-input placeholder="已审核通过的短信签名名称" v-model:value="model.sign" />
      </n-form-item-gi>
      <n-form-item-gi label="模版ID" path="template" span="3">
        <n-input placeholder="服务商平台上配置的短信模版ID" v-model:value="model.template" />
      </n-form-item-gi>
      <n-form-item-gi label="接收手机" path="receiver" span="6">
        <n-input placeholder="默认报警手机，当任务没有配置维护者时会发送到此手机号" v-model:value="model.receiver" />
      </n-form-item-gi>
      <n-form-item-gi label="内容模版" path="body" span="6">
        <n-input
          type="textarea"
          rows="5"
          placeholder="内容模版，支持用 Go 模版语法插入变量，转换后的内容会用 content 模版变量传递给短信平台"
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
  NSpace,
  NIcon,
  NGrid,
  NGi,
  NInput,
  NForm,
  NFormItemGi,
  NRadioGroup,
  NRadio,
  NSwitch,
} from "naive-ui";
import { SaveOutline as SaveIcon, } from "@vicons/ionicons5";
import configApi from "@/api/config";
import type { SmsOptions } from "@/api/config";
import { requiredRule, useForm } from "@/utils/form";

const key = 'alert.sms'
const form = ref();
const rules: any = {
  provider: requiredRule(),
  key: requiredRule(),
  secret: requiredRule(),
  sign: requiredRule(),
  template: requiredRule(),
};
const model = ref({
  provider: 'aliyun',
} as SmsOptions)
const enabled = computed({
  get() { return model.value.enabled === 'true' },
  set(v: boolean) { model.value.enabled = v.toString() },
});
const { submit, submiting } = useForm(form, () => configApi.save(key, model.value))

async function fetchData() {
  let r = await configApi.find(key);
  model.value = r.data as SmsOptions;
}

onMounted(fetchData);
</script>
