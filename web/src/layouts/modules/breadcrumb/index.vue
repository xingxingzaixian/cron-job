<template>
  <NBreadcrumb>
    <!-- define component: BreadcrumbContent -->
    <DefineBreadcrumbContent v-slot="{ breadcrumb }">
      <div class="flex items-center align-middle">
        <component :is="breadcrumb.icon" class="mr-4px" />
        {{ breadcrumb.label }}
      </div>
    </DefineBreadcrumbContent>
    <!-- define component: BreadcrumbContent -->

    <NBreadcrumbItem v-for="item in routeStore.breadcrumbs" :key="item.key">
      <NDropdown v-if="item.options?.length" :options="item.options" @select="handleClickMenu">
        <BreadcrumbContent :breadcrumb="item" />
      </NDropdown>
      <BreadcrumbContent v-else :breadcrumb="item" />
    </NBreadcrumbItem>
  </NBreadcrumb>
</template>

<script setup lang="ts">
import { createReusableTemplate } from '@vueuse/core';
import useRouteStore from '@/store/modules/route';
import { useRouter } from 'vue-router';

defineOptions({
  name: 'GlobalBreadcrumb'
});

const routeStore = useRouteStore();
const router = useRouter();

interface BreadcrumbContentProps {
  breadcrumb: Global.Menu;
}

const [DefineBreadcrumbContent, BreadcrumbContent] = createReusableTemplate<BreadcrumbContentProps>();

function handleClickMenu(key: string) {
  router.push({ name: key });
}
</script>

<style scoped></style>
