<template>
  <n-popover trigger="hover">
    <template #trigger>
      <n-button text class="h-36px text-icon" @click="$emit('toggle')">
        <SvgIcon :icon="full ? 'gridicons:fullscreen-exit' : 'gridicons:fullscreen'" />
      </n-button>
    </template>
    <span>{{ full ? $t('icon.fullscreenExit') : $t('icon.fullscreen') }}</span>
  </n-popover>
</template>

<script setup lang="ts">
import { $t } from '@/locales';
import SvgIcon from '@/components/custom/SvgIcon.vue';

defineOptions({
  name: 'FullScreen'
});

interface Props {
  full?: boolean;
}

defineProps<Props>();
</script>

<style scoped></style>
