<template>
  <div class="tab-pane" :class="{active: active}" :style="{borderColor: themeVars.primaryColor}">
    <router-link
      :to="href"
      :style="{color: active ? themeVars.primaryColor : themeVars.textColorBase}"
    >
      <slot />
    </router-link>
  </div>
</template>

<script setup lang="ts">
import { NA, useThemeVars } from "naive-ui";
import { RouterLink } from "vue-router";

const props = defineProps({
  active: Boolean,
  title: String,
  href: {
    type: String,
    required: true,
  },
});
const themeVars = useThemeVars();
</script>

<style scoped>
.tab-pane {
  padding-bottom: 4px;
  margin-bottom: -1px;
}
.active {
  border-bottom: 2px solid;
}
a {
  text-decoration: none;
}
</style>