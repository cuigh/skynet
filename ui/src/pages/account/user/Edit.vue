<template>
  <PageHeader :title="$route.meta.title" :subtitle="user.id">
    <template #action>
      <n-button size="small" @click="$router.push('/account/users')">
        <template #icon>
          <n-icon>
            <back-icon />
          </n-icon>
        </template>返回
      </n-button>
    </template>
  </PageHeader>
  <n-space class="page-body" vertical :size="12">
    <n-alert type="info" v-if="!$route.params.id">初始密码为：754321，创建后请让用户尽快修改。</n-alert>
    <n-form :model="user" :rules="rules" ref="form" label-placement="top">
      <n-grid cols="1 640:2" :x-gap="24">
        <n-form-item-gi label="用户名" path="name">
          <n-input placeholder="用户名" v-model:value="user.name" />
        </n-form-item-gi>
        <n-form-item-gi label="登录名" path="login_name">
          <n-input placeholder="登录名" v-model:value="user.login_name" />
        </n-form-item-gi>
        <n-form-item-gi label="邮箱" path="email">
          <n-input placeholder="电子邮箱" v-model:value="user.email" />
        </n-form-item-gi>
        <n-form-item-gi label="手机" path="phone">
          <n-input placeholder="手机号码" v-model:value="user.phone" />
        </n-form-item-gi>
        <n-form-item-gi label="企业微信" path="wecom">
          <n-input placeholder="企业微信用户 ID" v-model:value="user.wecom" />
        </n-form-item-gi>
        <n-form-item-gi label="管理员" span="2" path="admin" label-placement="left">
          <n-switch v-model:value="user.admin">
            <template #checked>是</template>
            <template #unchecked>否</template>
          </n-switch>
        </n-form-item-gi>
        <n-form-item-gi label="角色" span="2" path="roles">
          <n-checkbox-group v-model:value="user.roles">
            <n-space item-style="display: flex;">
              <n-checkbox :value="r.id" :label="r.name" v-for="r of roles" />
            </n-space>
          </n-checkbox-group>
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
  </n-space>
</template>

<script setup lang="ts">
import { onMounted, ref } from "vue";
import {
  NButton,
  NSpace,
  NInput,
  NIcon,
  NForm,
  NGrid,
  NGi,
  NFormItemGi,
  NSwitch,
  NCheckboxGroup,
  NCheckbox,
  NAlert,
} from "naive-ui";
import {
  ArrowBackCircleOutline as BackIcon,
  SaveOutline as SaveIcon,
} from "@vicons/ionicons5";
import PageHeader from "@/components/PageHeader.vue";
import { useRoute } from "vue-router";
import { router } from "@/router/router";
import userApi from "@/api/user";
import roleApi from "@/api/role";
import type { User } from "@/api/user";
import type { Role } from "@/api/role";
import { useForm, emailRule, phoneRule, requiredRule } from "@/utils/form";

const route = useRoute();
const user = ref({ status: 1, admin: false } as User)
const roles = ref([] as Role[]);
const rules: any = {
  name: requiredRule(),
  login_name: requiredRule(),
  email: emailRule(),
  phone: phoneRule(),
};
const form = ref();
const { submit, submiting } = useForm(form, () => userApi.save(user.value), () => {
  window.message.info("操作成功");
  router.push("/account/users")
})

async function fetchData() {
  const id = route.params.id as string || ''
  if (id) {
    let r = await userApi.find(id);
    user.value = r.data as User;
  }
  let r = await roleApi.search()
  roles.value = r.data as Role[]
}

onMounted(fetchData);
</script>
