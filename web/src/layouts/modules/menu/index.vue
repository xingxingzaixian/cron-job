<template>
  <NMenu
    v-model:expanded-keys="expandedKeys"
    :value="selectedKey"
    :collapsed="siderCollapse"
    :collapsed-width="themeStore.collapsedWidth"
    :collapsed-icon-size="22"
    :options="naiveMenus"
    :inverted="themeStore.darkMode"
    :indent="18"
    responsive
    @update:value="handleClickMenu"
  />
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import type { MentionOption } from 'naive-ui';
import useAppStore from '@/store/modules/app';
import useThemeStore from '@/store/modules/theme';
import useRouteStore from '@/store/modules/route';

defineOptions({
  name: 'LayoutMenu'
});

const route = useRoute();
const appStore = useAppStore();
const themeStore = useThemeStore();
const routeStore = useRouteStore();
const router = useRouter();

const naiveMenus = computed(() => routeStore.menus as unknown as MentionOption[]);

const siderCollapse = computed(() => appStore.siderCollapse);

const selectedKey = computed<string>(() => {
  const { hideInMenu, activeMenu } = route.meta;
  const name = route.name as string;

  const routeName = ((hideInMenu ? activeMenu : name) as string) || name;

  console.log(routeName);
  return routeName;
});

const expandedKeys = ref<string[]>([]);

function handleClickMenu(key: string) {
  router.push({ name: key });
}

function updateExpandedKeys() {
  if (siderCollapse.value || !selectedKey.value) {
    expandedKeys.value = [];
    return;
  }
  expandedKeys.value = routeStore.getSelectedMenuKeyPath(selectedKey.value);
}

watch(
  () => route.name,
  () => {
    updateExpandedKeys();
  },
  { immediate: true }
);
</script>

<style scoped></style>
