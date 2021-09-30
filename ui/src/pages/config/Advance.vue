<template>
  <PageHeader title="高级设置" subtitle="测试页面，功能暂未实现……"/>
  <n-space class="page-body" vertical :size="12">
    <Panel title="外观" subtitle="网站风格与配色" divider="bottom" :collapsed="panel !== 'appearance'">
      <template #action>
        <n-button size="small" @click="togglePanel('appearance')">{{ panel === 'appearance' ? 'Collapse' : 'Expand' }}</n-button>
      </template>
      <div>
        <n-alert :show-icon="false">Test</n-alert>
        <n-button type="primary" size="small" style="margin-top: 12px">保存</n-button>
      </div>
    </Panel>
    <Panel title="语言" subtitle="系统使用的语言" divider="bottom" :collapsed="panel !== 'language'">
      <template #action>
        <n-button size="small" @click="togglePanel('language')">{{ panel === 'language' ? 'Collapse' : 'Expand' }}</n-button>
      </template>
      <div>
        <n-form :model="modelLang" ref="formLang" label-placement="left" :show-feedback="false">
          <n-form-item label="语言" path="lang">
            <n-radio-group v-model:value="modelLang.lang" size="small">
              <n-radio-button value="zh">中文</n-radio-button>
              <n-radio-button value="en">English</n-radio-button>
            </n-radio-group>
          </n-form-item>
        </n-form>
        <n-button type="primary" size="small" style="margin-top: 12px">保存</n-button>
      </div>
    </Panel>
    <Panel title="LDAP" subtitle="LDAP 登录设置" :collapsed="panel !== 'ldap'">
      <template #action>
        <n-button size="small" @click="togglePanel('ldap')">{{ panel === 'ldap' ? 'Collapse' : 'Expand' }}</n-button>
      </template>
      <n-alert :show-icon="false">Test</n-alert>
    </Panel>
  </n-space>
</template>

<script setup lang="ts">
import { h, onMounted, reactive, ref } from "vue";
import {
  NButton,
  NSpace,
  NIcon,
  NInput,
  NForm,
  NFormItem,
  NRadioGroup,
  NRadio,
  NRadioButton,
  NAlert,
} from "naive-ui";
import { SaveOutline as SaveIcon } from "@vicons/ionicons5";
import PageHeader from "@/components/PageHeader.vue";
import Panel from "@/components/Panel.vue";
import configApi from "@/api/config";

const modelLang = reactive({
  show: false,
  lang: 'zh',
});
const panel = ref('')

function togglePanel(name: string) {
  if (panel.value === name) {
    panel.value = ''
  } else {
    panel.value = name
  }
}

// async function fetchData() {
//   model.task = {} as Task;
//   try {
//     let r = await task.find(route.params.name as string);
//     model.task = r.data as Task;
//   } catch (e) {
//     message.error(e.toString());
//   }
// }

// onMounted(fetchData);
</script>

<style scoped>
.n-button {
  width: 75px;
}
</style>
