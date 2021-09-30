<template>
  <PageHeader :title="$route.meta.title" :subtitle="model.name">
    <template #action>
      <n-button size="small" @click="$router.push({ name: 'role.list' })">
        <template #icon>
          <n-icon>
            <back-icon />
          </n-icon>
        </template>返回
      </n-button>
    </template>
  </PageHeader>
  <n-space class="page-body" vertical :size="12">
    <n-form :model="model" :rules="rules" ref="form" label-placement="top">
      <n-grid cols="1 640:2" :x-gap="24">
        <n-form-item-gi label="ID" path="id">
          <n-input
            placeholder="角色 ID 只能是小写字母和中划线，如：task-editor"
            v-model:value="model.id"
            :disabled="Boolean(route.params.id)"
          />
        </n-form-item-gi>
        <n-form-item-gi label="名称" path="name">
          <n-input placeholder="角色名" v-model:value="model.name" />
        </n-form-item-gi>
        <n-form-item-gi label="说明" path="desc" span="2">
          <n-input placeholder="角色说明" v-model:value="model.desc" />
        </n-form-item-gi>
        <n-form-item-gi label="角色" span="2" path="perms">
          <n-checkbox-group v-model:value="model.perms">
            <n-space item-style="display: flex;">
              <n-checkbox :value="p.value" :label="p.text" v-for="p of perms" />
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
import { onMounted, reactive, ref } from "vue";
import {
  NButton,
  NSpace,
  NInput,
  NIcon,
  NForm,
  NGrid,
  NGi,
  NFormItemGi,
  NCheckboxGroup,
  NCheckbox,
} from "naive-ui";
import {
  ArrowBackCircleOutline as BackIcon,
  SaveOutline as SaveIcon,
} from "@vicons/ionicons5";
import { useRoute } from "vue-router";
import PageHeader from "@/components/PageHeader.vue";
import { router } from "@/router/router";
import roleApi from "@/api/role";
import type { Role } from "@/api/role";
import { perms } from "@/utils/perm";
import { useForm, regexRule, requiredRule } from "@/utils/form";

const route = useRoute();
const model = ref({} as Role);
const rules: any = {
  id: [requiredRule(), regexRule(/^[a-z]+[a-z-]*[a-z]+$/, '只能包含小写字母和中划线')],
  name: requiredRule(),
};
const form = ref();
const { submit, submiting } = useForm(form, () => roleApi.save(model.value), () => {
  window.message.info("操作成功");
  router.push("/account/roles")
})

async function fetchData() {
  let id = route.params.id as string
  if (id) {
    let r = await roleApi.find(id);
    model.value = r.data as Role
  }
}

onMounted(fetchData);
</script>
