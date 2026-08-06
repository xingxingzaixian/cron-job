<template>
  <NCard v-if="result" size="small" class="mt-16px">
    <template #header>
      <div class="flex items-center justify-between w-full">
        <div class="flex items-center gap-16px">
          <NTag :type="result.success ? 'success' : 'error'" size="small" round>
            {{ result.success ? $t('page.task.list.testSuccess') : $t('page.task.list.testFailed') }}
          </NTag>
          <span v-if="result.status_code" class="text-sm">{{ $t('page.task.list.testStatusCode') }}: {{ result.status_code }}</span>
          <span class="text-sm">{{ $t('page.task.list.testDuration') }}: {{ result.duration_ms }} ms</span>
          <span v-if="result.size" class="text-sm">{{ $t('page.task.list.testSize') }}: {{ result.size }} B</span>
        </div>
        <NRadioGroup v-if="isJson" v-model:value="displayMode" size="small">
          <NRadioButton value="json">{{ $t('page.task.list.responseJson') }}</NRadioButton>
          <NRadioButton value="raw">{{ $t('page.task.list.responseRaw') }}</NRadioButton>
        </NRadioGroup>
      </div>
    </template>
    <p v-if="result.error" class="mb-8px" style="color: #d03050; word-break: break-all;">
      {{ result.error }}
    </p>
    <NInput :value="displayedOutput" type="textarea" readonly :rows="10" class="font-mono" :placeholder="$t('page.task.list.testOutput')" />
  </NCard>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue';
import { $t } from '@/locales';
import type { TaskTestOutput } from '@/api/task/types';

const props = defineProps<{
  result: TaskTestOutput | null;
}>();

const displayMode = ref<'json' | 'raw'>('json');

const isJson = computed(() => {
  if (!props.result?.output) {
    return false;
  }
  try {
    JSON.parse(props.result.output);
    return true;
  } catch (error) {
    return false;
  }
});

const displayedOutput = computed(() => {
  if (!props.result) {
    return '';
  }
  if (displayMode.value === 'json' && isJson.value) {
    try {
      return JSON.stringify(JSON.parse(props.result.output), null, 2);
    } catch (error) {
      return props.result.output;
    }
  }
  return props.result.output;
});

watch(
  () => props.result,
  () => {
    displayMode.value = isJson.value ? 'json' : 'raw';
  }
);
</script>
