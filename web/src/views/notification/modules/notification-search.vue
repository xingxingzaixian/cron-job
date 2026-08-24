<template>
  <NForm ref="formRef" :model="model" inline class="search-form" label-placement="left" :show-feedback="false">
    <NFormItem label="任务ID" path="task_id" class="search-field">
      <NInputNumber
        v-model:value="model.task_id"
        placeholder="任务ID"
        size="small"
        class="search-input"
        :min="0"
        clearable
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
import { useForm } from '@/hooks';
import type { SearchNotificationParams } from '@/api/notification/types';

type SearchTask = Omit<SearchNotificationParams, 'pageNo' | 'pageSize'>;
const emit = defineEmits<{
  (e: 'reset'): void;
  (e: 'search', searchParams: SearchTask): void;
}>();

const { formRef, validate, restoreValidation } = useForm();
const model = reactive<SearchTask>({
  task_id: 0,
});

async function reset() {
  await restoreValidation();
  emit('reset');
}

async function search() {
  await validate();
  emit('search', model);
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

:deep(.n-input-number) {
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

  .search-input {
    width: 100%;
  }

  .search-actions {
    display: flex;
    justify-content: center;
  }
}
</style>
