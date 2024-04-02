<template>
  <DarkModeContainer class="h-full flex items-center">
    <LayoutLogo class="h-full" style="width: 240px" />
    <div class="h-full flex items-center flex-1-hidden">
      <Breadcrumb class="ml-12px" />
    </div>
    <div class="h-full flex items-center justify-end">
      <FullScreen :full="isFullscreen" @toggle="toggle" />
      <LangSwitch :lang="appStore.locale" :lang-options="appStore.localeOptions" @change-lang="appStore.changeLocale" />
      <ThemeSwitch :is-dark="themeStore.darkMode" @switch="themeStore.toggleThemeScheme" />
      <!-- <UserAvatar /> -->
    </div>
  </DarkModeContainer>
</template>

<script setup lang="ts">
import { useFullscreen } from '@vueuse/core';
import useAppStore from '@/store/modules/app';
import useThemeStore from '@/store/modules/theme';
import ThemeSwitch from './components/theme-switch.vue';
import DarkModeContainer from './components/dark-mode-container.vue';
import LayoutLogo from '../logo/index.vue';
import FullScreen from './components/full-screen.vue';
import Breadcrumb from '../breadcrumb/index.vue';
import LangSwitch from './components/lang-switch.vue';

defineOptions({
  name: 'LayoutHeader'
});

const appStore = useAppStore();
const themeStore = useThemeStore();
const { isFullscreen, toggle } = useFullscreen();
</script>

<style scoped>
.shadow-header {
  border-bottom: 1px solid rgb(0, 21, 41, 0.08);
}

.flex-1-hidden {
  flex: 1 1 0%;
  overflow: hidden;
}
</style>
