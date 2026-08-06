<template>
  <div class="w-full">
    <div v-for="(row, index) in rows" :key="index" class="flex items-center gap-8px mb-8px">
      <NCheckbox v-model:checked="row.enabled" />
      <NInput v-model:value="row.key" size="small" :placeholder="$t('http.kvKey')" class="flex-1" @update:value="onChange" />
      <NInput v-model:value="row.value" size="small" :placeholder="$t('http.kvValue')" class="flex-1" @update:value="onChange" />
      <NInput v-model:value="row.desc" size="small" :placeholder="$t('http.kvDesc')" class="flex-1" @update:value="onChange" />
      <NButton type="error" quaternary circle size="small" @click="removeRow(index)">
        <template #icon>
          <icon-ic-round-delete color="red" />
        </template>
      </NButton>
    </div>
    <NButton dashed block size="small" @click="addRow">
      <template #icon>
        <icon-ic-round-plus />
      </template>
      {{ $t('http.addRow') }}
    </NButton>
  </div>
</template>

<script setup lang="ts">
import { $t } from '@/locales';
import { newKvRow, type KvRow } from './types';

const props = defineProps<{
  rows: KvRow[];
}>();

const emit = defineEmits<{
  (e: 'change'): void;
}>();

const onChange = () => {
  emit('change');
};

const addRow = () => {
  props.rows.push(newKvRow());
  emit('change');
};

const removeRow = (index: number) => {
  props.rows.splice(index, 1);
  emit('change');
};
</script>
