<template>
  <n-layout :position="isMobile ? 'static' : 'absolute'">
    <n-layout-header bordered>
      <div class="header-left" align="center">
        <n-popover
          v-if="isMobile || isTablet"
          style="padding: 0; width: 250px"
          placement="bottom-end"
          display-directive="show"
          trigger="click"
          ref="menuPopover"
        >
          <template #trigger>
            <n-button size="small" style="margin-right: 8px">
              <template #icon>
                <n-icon>
                  <menu-outline />
                </n-icon>
              </template>
            </n-button>
          </template>
          <div style="overflow: auto; max-height: 79vh">
            <n-menu
              :value="menuValue"
              :options="menuOptions"
              :indent="18"
              @update:value="menuPopover.setShow(false)"
              :render-label="renderMenuLabel"
            />
          </div>
        </n-popover>
        <n-text tag="div" class="logo" :depth="1" @click="$router.push('/')">
          <img src="../assets/skynet.svg" v-if="!isMobile" />
          Skynet
        </n-text>
        <div style="margin: 4px 0 0 8px">
          <n-text depth="3">v0.2</n-text>
        </div>
      </div>
      <n-space justify="end" align="center" class="header-right" :size="0">
        <n-button
          type="default"
          size="large"
          :bordered="false"
          style="margin-right: 4px"
          tag="a"
          href="https://github.com/cuigh/skynet"
          target="_blank"
        >
          <template #icon>
            <n-icon>
              <LogoGithub />
            </n-icon>
          </template>
        </n-button>
        <n-switch :value="darkTheme" @update:value="changeTheme">
          <template #checked>深色</template>
          <template #unchecked>浅色</template>
        </n-switch>
        <n-dropdown @select="selectOption" trigger="hover" :options="dropdownOptions" show-arrow>
          <n-button type="default" size="small" :bordered="false" style="margin-left: 12px;">
            <template #icon>
              <n-icon>
                <PersonOutline />
              </n-icon>
            </template>
            {{ store.state.name }}
          </n-button>
        </n-dropdown>
        <n-tooltip trigger="hover">
          <template #trigger>
            <n-button type="default" size="small" :bordered="false" @click="logout">
              <template #icon>
                <n-icon>
                  <LogOutOutline />
                </n-icon>
              </template>
            </n-button>
          </template>
          登出
        </n-tooltip>
      </n-space>
    </n-layout-header>
    <n-layout
      has-sider
      :position="isMobile ? 'static' : 'absolute'"
      :style="isMobile ? '' : 'top: 56px; bottom: 64px'"
    >
      <n-layout-sider
        v-if="!isMobile && !isTablet"
        bordered
        width="240"
        :collapsed-width="64"
        :collapsed="collapsed"
        collapse-mode="width"
        show-trigger="bar"
        trigger-style="right: -25px"
        collapsed-trigger-style="right: -25px"
        @collapse="collapsed = true"
        @expand="collapsed = false"
      >
        <n-menu
          :value="menuValue"
          :options="menuOptions"
          :collapsed="collapsed"
          :collapsed-width="64"
          :collapsed-icon-size="22"
          :root-indent="20"
          :indent="24"
          :render-label="renderMenuLabel"
          :expanded-keys="expandedKeys"
          @update:expanded-keys="updateExpandedKeys"
        />
      </n-layout-sider>
      <n-layout-content>
        <router-view></router-view>
        <n-back-top :right="16" :bottom="80" />
      </n-layout-content>
    </n-layout>
    <n-layout-footer bordered :position="isMobile ? 'static' : 'absolute'">
      <span>© 2021 cuigh. 保留所有权利。</span>
    </n-layout-footer>
  </n-layout>
  <n-modal
    v-model:show="modelPwd.showDlg"
    preset="card"
    title="修改密码"
    size="small"
    style="width: 500px;"
  >
    <n-form :model="modelPwd" ref="form" :rules="rules">
      <n-form-item path="old_pwd" label="旧密码">
        <n-input v-model:value="modelPwd.old_pwd" type="password" placeholder="请输入旧密码" />
      </n-form-item>
      <n-form-item first path="new_pwd" label="新密码">
        <n-input v-model:value="modelPwd.new_pwd" type="password" placeholder="请输入新密码" />
      </n-form-item>
      <n-form-item first path="confirm_pwd" label="新密码确认">
        <n-input
          :disabled="!modelPwd.new_pwd"
          v-model:value="modelPwd.confirm_pwd"
          type="password"
          placeholder="请再次输入新密码"
        />
      </n-form-item>
      <div style="display: flex; justify-content: flex-end;">
        <n-button
          type="primary"
          :disabled="submiting"
          :loading="submiting"
          @click.prevent="modifyPwd"
        >确定</n-button>
      </div>
    </n-form>
  </n-modal>
</template>

<script setup lang="ts">
import { ref, computed, reactive, watch } from "vue";
import {
  NButton,
  NIcon,
  NMenu,
  NText,
  NSpace,
  NLayout,
  NLayoutHeader,
  NLayoutSider,
  NLayoutContent,
  NLayoutFooter,
  NPopover,
  NTooltip,
  NDropdown,
  NSwitch,
  NBackTop,
  NModal,
  NForm,
  NFormItem,
  NInput,
} from "naive-ui";
import { MenuOutline, PersonOutline, LogOutOutline, LogoGithub } from "@vicons/ionicons5";
import { RouterView, useRouter, useRoute } from "vue-router";
import { useStore } from "vuex";
import { useIsMobile, useIsTablet } from "@/utils";
import { findMenuValue, renderMenuLabel, menuOptions, findActiveOptions } from "@/router/menu";
import userApi from "@/api/user";
import { Mutations } from "@/store/mutations";
import { useForm, requiredRule, customRule } from "@/utils/form";

// user dropdown options
const dropdownOptions = [
  {
    label: "个人信息",
    key: "profile",
  },
  {
    label: "修改密码",
    key: "password",
  },
];
const store = useStore();
const router = useRouter();
const route = useRoute();
const menuPopover = ref();
const collapsed = ref(false)
const expandedKeys = ref([] as string[]);
const isMobile = useIsMobile()
const isTablet = useIsTablet()
const darkTheme = computed(() => store.state.theme === "dark")
const menuValue = computed(() => findMenuValue(route))

// modify password
const form = ref()
const modelPwd = reactive({
  showDlg: false,
  old_pwd: "",
  new_pwd: "",
  confirm_pwd: "",
})
const rules = {
  old_pwd: requiredRule(),
  new_pwd: [
    requiredRule(),
    customRule((rule: any, value: string) => value !== modelPwd.old_pwd, '新密码不能跟旧密码一样'),
  ],
  confirm_pwd: [
    requiredRule("", "请再次输入新密码"),
    customRule((rule: any, value: string) => value === modelPwd.new_pwd, '两次输入的密码不一致'),
  ],
};
const { submit: modifyPwd, submiting } = useForm(
  form,
  () => userApi.modifyPassword({ old_pwd: modelPwd.old_pwd, new_pwd: modelPwd.new_pwd }),
  () => window.message.info("密码修改成功")
);

function updateExpandedKeys(data: any) {
  expandedKeys.value = data
}

function selectOption(key: any) {
  switch (key as string) {
    case "profile":
      router.push("/profile")
      return
    case "password":
      modelPwd.showDlg = true;
      return
    default:
      console.info(key)
  }
}

function logout() {
  store.commit(Mutations.Logout);
  router.push("/login");
}

function changeTheme(value: boolean | undefined) {
  store.commit(Mutations.SetTheme, value ? "dark" : "light");
}

watch(() => route.path, (path: string) => {
  let keys = findActiveOptions(route).map((opt: any) => opt.key) as string[]
  expandedKeys.value = keys;
})
</script>

<style scoped>
::v-deep(.header-right .n-button__content) {
  margin-top: 4px;
}
.header-left {
  flex-grow: 1;
  width: 250px;
  display: flex;
  align-items: center;
}
.header-right {
  width: 320px;
}
/* .n-layout-header {
  background-color: #363636;
} */
.n-layout-sider {
  box-shadow: 0 1px 2px rgb(10 10 10 / 10%);
}
.n-layout-footer {
  box-shadow: 1px 0px 2px rgb(10 10 10 / 10%);
  /* background-image: radial-gradient(circle at 1% 1%,#328bf2,#1644ad); */
}
/* .n-layout-header {
  background-image: linear-gradient(to right, rgb(91, 121, 162) 0%, rgb(46, 68, 105) 100%);
}
.logo {
  color: white;
}
.n-layout-header .n-icon {
  color: white;
}
::v-deep(.n-layout-header .n-button__content) {
  color: white;
}
::v-deep(.n-layout-header .n-button__content:hover) {
  color: green;
} */
</style>
