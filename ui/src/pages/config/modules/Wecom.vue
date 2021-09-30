<template>
  <n-form :model="model" :rules="rules" ref="form" label-placement="top">
    <n-grid cols="1 640:6" :x-gap="24">
      <n-form-item-gi label="启用" span="3" label-placement="left">
        <n-switch v-model:value="enabled">
          <template #checked>是</template>
          <template #unchecked>否</template>
        </n-switch>
      </n-form-item-gi>
      <n-form-item-gi label="报警方式" path="mode" span="3" label-placement="left">
        <n-radio-group v-model:value="model.mode" name="mode" size="small">
          <n-space>
            <n-radio value="robot">群聊机器人</n-radio>
            <n-radio value="app" disabled>自定义应用</n-radio>
          </n-space>
        </n-radio-group>
      </n-form-item-gi>
      <n-form-item-gi label="机器人密钥" path="robot_key" span="6" v-if="model.mode === 'robot'">
        <n-input
          placeholder="机器人密钥（GUID），格式为：xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxxx"
          v-model:value="model.robot_key"
        />
      </n-form-item-gi>
      <n-form-item-gi label="企业ID" path="corp_id" span="2" v-if="model.mode === 'app'">
        <n-input placeholder="企业ID(corp_id)" v-model:value="model.corp_id" />
      </n-form-item-gi>
      <n-form-item-gi label="应用ID" path="app_id" span="2" v-if="model.mode === 'app'">
        <n-input placeholder="应用ID(agent_id)" v-model:value="model.app_id" />
      </n-form-item-gi>
      <n-form-item-gi label="应用密钥" path="app_secret" span="2" v-if="model.mode === 'app'">
        <n-input placeholder="应用密钥(secret)" v-model:value="model.app_secret" />
      </n-form-item-gi>
      <n-form-item-gi label="用户" path="users" span="3" v-if="model.mode === 'app'">
        <n-input placeholder="默认接收消息的用户 ID 列表（多个 ID 用 | 分割）" v-model:value="model.users" />
      </n-form-item-gi>
      <n-form-item-gi label="部门" path="parties" span="3" v-if="model.mode === 'app'">
        <n-input placeholder="默认接收消息的部门 ID 列表（多个 ID 用 | 分割）" v-model:value="model.parties" />
      </n-form-item-gi>
      <n-form-item-gi label="标签" path="tags" span="3" v-if="model.mode === 'app'">
        <n-input placeholder="默认接收消息的标签 ID 列表（多个 ID 用 | 分割）" v-model:value="model.tags" />
      </n-form-item-gi>
      <n-form-item-gi label="群聊" path="chats" span="3" v-if="model.mode === 'app'">
        <n-input placeholder="默认接收消息的群聊 ID 列表（多个 ID 用 | 分割）" v-model:value="model.chats" />
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
      <n-form-item-gi label="消息类型" path="msg_type" span="3" label-placement="left">
        <n-radio-group v-model:value="model.msg_type" name="mode" size="small">
          <n-space>
            <n-radio value="markdown">Markdown</n-radio>
            <n-radio value="text">文本</n-radio>
            <!-- <n-radio value="textcard" disabled>文本卡片</n-radio> -->
          </n-space>
        </n-radio-group>
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
import type { WecomOptions } from "@/api/config";
import { useForm } from "@/utils/form";

const key = 'alert.wecom'
const form = ref();
const rules: any = {
};
const model = ref({
  enabled: '',
  mode: '',
  robot_key: '',
  corp_id: '',
  app_id: '',
  app_secret: '',
  receiver: '',
  title: '',
  body: '',
  msg_type: '',
} as WecomOptions)
const enabled = computed({
  get() { return model.value.enabled === 'true' },
  set(v: boolean) { model.value.enabled = v.toString() },
});
const { submit, submiting } = useForm(form, () => configApi.save(key, model.value))

async function fetchData() {
  let r = await configApi.find(key);
  model.value = r.data as WecomOptions;
  model.value.mode = model.value.mode ?? 'robot';
  model.value.msg_type = model.value.msg_type ?? 'markdown';
}

onMounted(fetchData);
</script>
