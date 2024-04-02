<template>
  <NCard :bordered="false" size="small" class="rd-8 shadow-sm">
    <div class="rd-8px px-16px pb-4px pt-8px text-white bg-gradient-to-br from-pink-500 to-red-500">
      <h3 class="text-16px">{{ title }}</h3>
      <div class="flex justify-between items-center pt-12px">
        <SvgIcon :icon="icon" class="text-32px" />
        <span class="text-32px">{{ count }}</span>
      </div>
    </div>
  </NCard>
</template>

<script lang="ts" setup>
defineProps<{
  title: string;
  icon: string;
  count: number;
}>();
</script>
