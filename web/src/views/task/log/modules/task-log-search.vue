<template>
  <NForm ref="formRef" :model="model" class="search-form" label-placement="left" :show-feedback="false">
    <NFormItem :label="$t('page.task.list.name')" path="name" class="search-field">
      <NInput
        v-model:value="model.taskName"
        :placeholder="$t('page.task.list.name')"
        size="small"
        class="search-input"
        clearable
      />
    </NFormItem>
    <NFormItem :label="$t('page.task.list.status')" path="status" class="search-field">
      <NSelect
        v-model:value="model.status"
        :placeholder="$t('page.task.list.status')"
        :options="statusOptions"
        clearable
        size="small"
        class="search-select"
      />
    </NFormItem>
    <NFormItem :label="$t('page.task.log.dateRange')" path="dateRange" class="search-field">
      <NDatePicker
        v-model:value="model.dateRange"
        type="daterange"
        clearable
        size="small"
        class="search-date"
        :start-placeholder="$t('page.task.log.startDate')"
        :end-placeholder="$t('page.task.log.endDate')"
      />
    </NFormItem>
    <NFormItem class="search-actions">
      <NSpace size="small">
        <NButton size="small" @click="reset">
          <template #icon>
            <icon-ic-round-refresh />
          </template>
          {{ $t('common.reset') }}
        </NButton>
        <NButton type="primary" size="small" @click="search">
          <template #icon>
            <icon-ic-round-search />
          </template>
          {{ $t('common.search') }}
        </NButton>
      </NSpace>
    </NFormItem>
  </NForm>
</template>

<script lang="tsx" setup>
import { $t } from '@/locales';
import { reactive } from 'vue';
import dayjs from 'dayjs';
import { useForm } from '@/hooks';
import type { QueryTaskLog } from '@/api/task/types';
import type { SelectOption } from 'naive-ui';

type SearchTask = Omit<QueryTaskLog, 'taskId' | 'pageNo' | 'pageSize'>;
const emit = defineEmits<{
  (e: 'reset'): void;
  (e: 'search', searchParams: SearchTask): void;
}>();

const { formRef, validate, restoreValidation } = useForm();
const model = reactive<SearchTask & { dateRange: [number, number] | null }>({
  taskName: '',
  status: -1,
  dateRange: null
});

const statusOptions: SelectOption[] = [
  {
    label: $t('task.protocol.all'),
    value: -1
  },
  {
    label: $t('page.task.log.failure'),
    value: 2
  },
  {
    label: $t('page.task.log.running'),
    value: 3
  },
  {
    label: $t('page.task.log.success'),
    value: 4
  }
];
async function reset() {
  await restoreValidation();
  model.taskName = '';
  model.status = -1;
  model.dateRange = null;
  emit('reset');
}

async function search() {
  await validate();
  const { dateRange, ...params } = model;
  const searchParams: SearchTask = { ...params };
  if (dateRange && dateRange.length === 2) {
    searchParams.startTime = dayjs(dateRange[0]).format('YYYY-MM-DD HH:mm:ss');
    searchParams.endTime = dayjs(dateRange[1]).format('YYYY-MM-DD HH:mm:ss');
  }
  emit('search', searchParams);
}
</script>

<style scoped>
.search-form {
  display: flex;
  align-items: center;
  gap: 16px;
  margin: 0;
}

.search-field {
  flex-shrink: 0;
}

.search-input {
  width: 140px;
}

.search-select {
  width: 120px;
}

.search-date {
  width: 260px;
}

.search-actions {
  flex-shrink: 0;
}

/* 表单项样式优化 */
:deep(.n-form-item) {
  margin-bottom: 0;
}

:deep(.n-form-item-label) {
  font-size: 13px;
  color: #64748b;
  font-weight: 500;
  margin-right: 8px;
}

:deep(.n-input) {
  --n-border-radius: 6px;
  --n-border-color: rgba(209, 213, 219, 0.6);
  --n-border-color-hover: rgba(102, 126, 234, 0.4);
  --n-border-color-focus: rgba(102, 126, 234, 0.6);
}

:deep(.n-select) {
  --n-border-radius: 6px;
  --n-border-color: rgba(209, 213, 219, 0.6);
  --n-border-color-hover: rgba(102, 126, 234, 0.4);
  --n-border-color-focus: rgba(102, 126, 234, 0.6);
}

:deep(.n-button) {
  --n-border-radius: 6px;
  font-size: 13px;
  height: 28px;
  padding: 0 12px;
}

/* 响应式设计 */
@media (max-width: 1200px) {
  .search-form {
    flex-wrap: wrap;
    justify-content: center;
    gap: 12px;
  }
}

@media (max-width: 768px) {
  .search-form {
    flex-direction: column;
    align-items: stretch;
    gap: 8px;
  }

  .search-input,
  .search-select {
    width: 100%;
  }

  .search-actions {
    display: flex;
    justify-content: center;
  }
}
</style>
