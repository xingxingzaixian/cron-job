<script setup lang="ts">
import { computed } from 'vue';
import { NConfigProvider, NGlobalStyle } from 'naive-ui';
import useAppStore from './store/modules/app';
import useThemeStore from './store/modules/theme';
import { naiveDateLocales, naiveLocales } from './locales';
import useRouteStore from '@/store/modules/route';

defineOptions({
  name: 'App'
});

const appStore = useAppStore();
const themeStore = useThemeStore();
const routeStore = useRouteStore();
routeStore.initAuthRoute();

const naiveLocale = computed(() => {
  return naiveLocales[appStore.locale];
});

const naiveDateLocale = computed(() => {
  return naiveDateLocales[appStore.locale];
});
</script>

<template>
  <n-config-provider
    :theme="themeStore.naiveTheme"
    :theme-overrides="themeStore.themeOverrides"
    :locale="naiveLocale"
    :date-locale="naiveDateLocale"
    class="h-full"
  >
    <n-global-style />
    <RouterView class="bg-layout" />
  </n-config-provider>
</template>

<style scoped></style>
