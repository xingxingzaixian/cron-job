<template>
  <SvgIcon :fill="fill" :icon="icon" />
</template>

<script lang="ts" setup>
import { computed } from 'vue';
defineOptions({ name: 'SystemLogo' });

interface Props {
  /** logo是否填充 */
  fill?: boolean;
}

const props = withDefaults(defineProps<Props>(), {
  fill: false
});

const icon = computed(() => (props.fill ? 'eos-icons:cronjob' : 'eos-icons:cronjob'));
</script>

<style scoped></style>
