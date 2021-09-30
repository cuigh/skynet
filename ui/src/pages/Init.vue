<template>
  <div :class="['container', isMobile ? '' : 'pc']">
    <h1 class="title">系统初始化</h1>
    <n-alert type="info" :show-icon="false">你正在首次使用 Skynet，请按照下面的流程初始化系统。</n-alert>
    <!-- <n-steps :current="current" :status="currentStatus">
      <n-step title></n-step>
      <n-step title></n-step>
      <n-step title></n-step>
    </n-steps>-->
    <n-steps :current="step.current" :status="step.status" vertical>
      <n-step title="初始化数据库">
        <div class="n-step-description" v-if="step.current === 1">
          <n-text depth="3" tag="p">创建 Skynet 数据库，并添加必需的索引。</n-text>
          <n-button type="primary" @click="initDB">开始</n-button>
        </div>
      </n-step>
      <n-step title="创建管理员账号">
        <div class="n-step-description" v-if="step.current === 2">
          <n-space vertical>
            <n-alert id="error" type="error" :show-icon="false" v-if="errors && errors.length > 0">
              <n-ul>
                <n-li v-for="e in errors">{{ e[0].message }}</n-li>
              </n-ul>
            </n-alert>
            <n-form
              :model="model"
              ref="form"
              :rules="rules"
              label-placement="top"
              :show-feedback="false"
            >
              <n-grid cols="1 600:2" :x-gap="12" :y-gap="12">
                <n-form-item-gi path="login_name" label="登录名">
                  <n-input v-model:value="model.login_name" placeholder="请输入登录名" clearable />
                </n-form-item-gi>
                <n-form-item-gi path="name" label="用户名">
                  <n-input placeholder="请输入用户名" v-model:value="model.name" />
                </n-form-item-gi>
                <n-form-item-gi path="password" label="密码">
                  <n-input
                    v-model:value="model.password"
                    type="password"
                    placeholder="请输入密码"
                    clearable
                  />
                </n-form-item-gi>
                <n-form-item-gi path="password_confirm" label="密码确认">
                  <n-input
                    v-model:value="model.password_confirm"
                    type="password"
                    placeholder="请再次输入密码"
                    clearable
                  />
                </n-form-item-gi>
                <n-form-item-gi path="email" label="电子邮箱">
                  <n-input placeholder="电子邮箱" v-model:value="model.email" />
                </n-form-item-gi>
                <n-form-item-gi path="phone" label="手机号码">
                  <n-input placeholder="手机号码" v-model:value="model.phone" />
                </n-form-item-gi>
              </n-grid>
              <n-button
                type="primary"
                :disabled="submiting"
                :loading="submiting"
                @click.prevent="createAdmin"
                style="margin-top: 12px"
              >提交</n-button>
            </n-form>
          </n-space>
        </div>
      </n-step>
      <n-step title="开始使用">
        <div class="n-step-description" v-if="step.current === 3">
          <n-text depth="3" tag="p">恭喜你，你已经完成了初始化。</n-text>
          <n-button type="primary" @click="router.push('/login')">登录系统</n-button>
        </div>
      </n-step>
    </n-steps>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from "vue";
import { useRouter } from "vue-router";
import {
  NSpace,
  NForm,
  NFormItemGi,
  NInput,
  NButton,
  NSteps,
  NStep,
  NAlert,
  NText,
  NGrid,
  NUl,
  NLi,
} from "naive-ui";
import systemApi from "@/api/system";
import { useIsMobile } from "@/utils";
import { requiredRule, lengthRule, passwordRule, emailRule, phoneRule, customRule, regexRule } from "@/utils/form";

const form = ref();
const submiting = ref(false);
const model = reactive({
  login_name: "",
  name: "",
  password: "",
  password_confirm: "",
  email: "",
  phone: "",
});
const rules = {
  login_name: [requiredRule('登录名'), regexRule(/^[a-z]+[a-z0-9_.]*$/i, '只能包含英文字母、数字、下划线和点，且必须以英文字母开头', '登录名')],
  name: requiredRule('用户名'),
  password: [requiredRule('密码'), lengthRule(6, 15, "密码"), passwordRule("密码")],
  password_confirm: [
    requiredRule('密码确认'),
    customRule((_, value: string): boolean => {
      return value === model.password;
    }, '两次输入的密码不一致', '密码确认')
  ],
  email: emailRule("电子邮箱"),
  phone: phoneRule("手机号码"),
};
const step = reactive({
  current: 1,
  status: 'process' as 'process' | 'finish' | 'error',
});
const errors = ref([] as any);
const router = useRouter();
const isMobile = useIsMobile()

function setStep(current: number) {
  step.current = current
  step.status = current === 3 ? 'finish' : 'process'
}

async function initDB(e: Event) {
  submiting.value = true;
  try {
    const r = await systemApi.initDB();
    setStep(2)
  } finally {
    submiting.value = false;
  }
}

function createAdmin(e: Event) {
  form.value.validate(async (errs: any) => {
    if (errs) {
      errors.value = errs
      step.status = 'error'
      return;
    }

    submiting.value = true;
    try {
      const r = await systemApi.initUser(model);
      setStep(3)
    } finally {
      submiting.value = false;
    }
  });
}

async function checkState() {
  const r = await systemApi.checkState();
  if (!r.data?.fresh) {
    router.push("/login")
  }
}

onMounted(checkState);
</script>

<style lang="scss" scoped>
.container {
  padding: 20px;
  border-radius: 5px;
  box-shadow: 1px 1px 10px #ddd;
  .title {
    text-align: center;
    margin-top: 0px;
  }
  .n-steps {
    margin-top: 20px;
  }
}
.pc {
  width: 640px;
  margin: 20px auto;
}
::v-deep(#error .n-alert-body) {
  padding: 6px 6px 6px 0;
}
</style>