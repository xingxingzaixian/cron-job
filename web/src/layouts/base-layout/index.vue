<template>
  <div class="h-full relative">
    <n-layout position="absolute">
      <n-layout-header class="layout-header" bordered>
        <LayoutHeader />
      </n-layout-header>
      <n-layout has-sider>
        <n-layout-sider
          bordered
          content-style="padding: 10px;"
          class="layout-sider"
          collapse-mode="width"
          :collapsed-width="themeStore.collapsedWidth"
          :width="themeStore.siderWidth"
        >
          <LayoutMenu />
        </n-layout-sider>
        <n-layout-content embedded content-style="padding: 24px; height: 100%" :native-scrollbar="false">
          <LayoutContent :show-padding="false" />
        </n-layout-content>
      </n-layout>
    </n-layout>
  </div>
</template>

<script setup lang="ts">
import LayoutHeader from '../modules/header/index.vue';
import LayoutContent from '../modules/content/index.vue';
import LayoutMenu from '../modules/menu/index.vue';
import useThemeStore from '@/store/modules/theme';

defineOptions({
  name: 'BaseLayout'
});

const themeStore = useThemeStore();
</script>

<style scoped>
.layout-header {
  height: 56px;
}

.layout-sider {
  height: calc(100vh - 56px);
}
</style>
