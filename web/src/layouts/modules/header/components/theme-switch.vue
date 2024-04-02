<template>
  <n-button text @click="handleSwitch" class="h-36px text-icon">
    <SvgIcon :icon="icon" />
  </n-button>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import SvgIcon from '@/components/custom/SvgIcon.vue';

defineOptions({ name: 'ThemeSwitch' });

interface Props {
  isDark?: boolean;
}

const props = defineProps<Props>();

interface Emits {
  (e: 'switch'): void;
}

const emit = defineEmits<Emits>();

function handleSwitch() {
  emit('switch');
}

const icon = computed(() => (props.isDark ? 'material-symbols:nightlight-rounded' : 'material-symbols:sunny'));
</script>
<style scoped></style>
