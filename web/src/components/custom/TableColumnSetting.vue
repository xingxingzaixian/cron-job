<template>
  <NPopover placement="bottom-end" trigger="click">
    <template #trigger>
      <NButton size="small">
        <template #icon>
          <icon-ant-design-setting-outlined />
        </template>
        {{ $t('common.columnSetting') }}
      </NButton>
    </template>
    <VueDraggable v-model="columns">
      <div
        v-for="item in columns"
        :key="item.key"
        class="h-36px flex items-center flex items-center rd-4px hover:(bg-primary bg-opacity-20)"
      >
        <icon-mdi:drag class="cursor-move" />
        <NCheckbox v-model:checked="item.checked">
          {{ item.title }}
        </NCheckbox>
      </div>
    </VueDraggable>
  </NPopover>
</template>

<script setup lang="ts" generic="T extends Record<string, unknown>, K = never">
import { VueDraggable } from 'vue-draggable-plus';
import { $t } from '@/locales';

defineOptions({
  name: 'TableColumnSetting'
});

type FilteredColumn = {
  key: string;
  title: string;
  checked: boolean;
};

const columns = defineModel<FilteredColumn[]>('columns', {
  required: true
});
</script>

<style scoped></style>
