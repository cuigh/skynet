<template>
  <PageHeader :title="$route.meta.title" :subtitle="model.name">
    <template #action>
      <n-button size="small" @click="$router.push('/tasks')">
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
        <n-form-item-gi label="名称" path="name">
          <n-input placeholder="任务名称" v-model:value="model.name" :disabled="Boolean(name)" />
        </n-form-item-gi>
        <n-form-item-gi label="描述" path="desc">
          <n-input placeholder="任务描述" v-model:value="model.desc" />
        </n-form-item-gi>
        <n-form-item-gi label="执行器" path="runner">
          <n-input
            placeholder="任务执行器，格式：schema://[name or address]，示例：http://task.test.com"
            v-model:value="model.runner"
          />
        </n-form-item-gi>
        <n-form-item-gi label="处理器" path="handler">
          <n-input placeholder="在执行器中注册的任务处理器" v-model:value="model.handler" />
        </n-form-item-gi>
        <n-form-item-gi label="是否启用" path="enabled">
          <n-switch v-model:value="model.enabled" />
        </n-form-item-gi>
        <n-form-item-gi label="报警方式" path="alert">
          <n-checkbox-group v-model:value="model.alerts">
            <n-space item-style="display: flex;">
              <n-checkbox :value="a.value" :label="a.text" v-for="a of alerts" />
            </n-space>
          </n-checkbox-group>
        </n-form-item-gi>
        <n-form-item-gi label="维护者" path="maintainers" span="2">
          <n-select
            placeholder="任务维护者"
            v-model:value="model.maintainers"
            multiple
            clearable
            filterable
            :options="users"
          />
        </n-form-item-gi>
        <n-form-item-gi span="2" label="触发器" path="triggers">
          <n-dynamic-input v-model:value="model.triggers" #="{ index, value }" :min="1" :max="5">
            <n-input-group>
              <n-input
                placeholder="Cron表达式：[秒] [分] [时] [日] [月] [周]，支持预定义宏：@yearly, @monthly, @weekly, @daily, @hourly"
                v-model:value="model.triggers[index]"
              />
              <n-button type="default" ghost @click="testCron(value)" :disabled="!value">测试</n-button>
            </n-input-group>
          </n-dynamic-input>
        </n-form-item-gi>
        <n-form-item-gi span="2" label="参数" path="args">
          <n-dynamic-input v-model:value="model.args" #="{ index, value }" :on-create="newArg">
            <n-input placeholder="参数名" v-model:value="value.name" />
            <div style="height: 34px; line-height: 34px; margin: 0 8px">=</div>
            <n-input placeholder="参数值" v-model:value="value.value" />
          </n-dynamic-input>
        </n-form-item-gi>
        <n-gi :span="2">
          <n-button
            @click.prevent="submit"
            type="primary"
            :disabled="submiting"
            :loading="submiting"
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
import { h, onMounted, reactive, ref } from "vue";
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
  NDynamicInput,
  NSelect,
  NCheckboxGroup,
  NCheckbox,
  NInputGroup,
} from "naive-ui";
import type { FormItemRule } from "naive-ui";
import {
  ArrowBackCircleOutline as BackIcon,
  SaveOutline as SaveIcon,
} from "@vicons/ionicons5";
import { parseExpression } from "cron-parser";
import PageHeader from "@/components/PageHeader.vue";
import taskApi from "@/api/task";
import userApi from "@/api/user";
import type { Task } from "@/api/task";
import { useRoute } from "vue-router";
import { router } from "@/router/router";
import { useForm, requiredRule, customRule } from "@/utils/form";
import { alerts } from "./task";
import { renderTime } from "@/utils/render";

const route = useRoute();
const name = route.params.name as string || ''
const model = ref({} as Task);
const rules: any = {
  name: requiredRule(),
  runner: requiredRule(),
  handler: requiredRule(),
  maintainers: customRule((rule: any, value: any) => value != null && value.length > 0, '不能为空', '', true),
  triggers: {
    required: true,
    trigger: ["blur", "input"],
    validator(rule: FormItemRule, values: string[]) {
      var empty = true
      if (values) {
        for (let v of values) {
          if (v) {
            empty = false
            try {
              parseExpression(v)
            } catch {
              return new Error(`'${v}' 不是一个有效的 Cron 表达式`)
            }
          }
        }
      }
      return empty ? new Error('请输入触发器') : true
    },
  },
};
const form = ref();
const { submit, submiting } = useForm(form, () => taskApi.save(model.value), () => {
  window.message.info("操作成功");
  router.push("/tasks")
})
const users = ref([] as any)

function newArg() {
  return {
    name: '',
    value: ''
  }
}

async function testCron(cron: string) {
  try {
    const exp = parseExpression(cron)
    const times: Date[] = []
    for (let i = 0; i < 10; i++) {
      times.push(exp.next().toDate())
    }
    window.dialog.success({
      iconPlacement: "top",
      title: `未来 ${times.length} 次触发时间`,
      content: () => h(NSpace, { vertical: true, size: 0 }, {
        default: () => times.map(t => renderTime(t.getTime()))
      }),
    })
  } catch (err: any) {
    window.dialog.error({
      iconPlacement: "top",
      content: `'${cron}' 不是一个有效的 Cron 表达式`,
    })
  }
}

async function fetchData() {
  if (name) {
    let tr = await taskApi.find(name);
    model.value = tr.data as Task;
  }

  let ur = await userApi.search({ page_index: 1, page_size: 1000 })
  users.value = ur.data?.items.map(u => {
    return {
      label: u.name,
      value: u.id,
    }
  })
}

onMounted(fetchData);
</script>
