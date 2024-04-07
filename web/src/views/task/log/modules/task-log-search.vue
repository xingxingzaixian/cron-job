<template>
  <NCard :title="$t('common.search')">
    <NForm ref="formRef" :model="model" label-placement="left">
      <NGrid responsive="screen" item-responsive>
        <NFormItemGi span="24 s:12 m:6" :label="$t('page.task.list.name')" path="name" class="pr-24px">
          <NInput v-model:value="model.taskName" :placeholder="$t('page.task.list.name')" />
        </NFormItemGi>
        <NFormItemGi span="24 s:12 m:6" :label="$t('page.task.log.taskId')" path="taskId" class="pr-24px">
          <NInput v-model:value="model.taskId" :placeholder="$t('page.task.log.taskId')" />
        </NFormItemGi>
        <NFormItemGi span="24 s:12 m:6" :label="$t('page.task.list.status')" path="status" class="pr-24px">
          <NSelect
            v-model:value="model.status"
            :placeholder="$t('page.task.list.status')"
            :options="statusOptions"
            clearable
          />
        </NFormItemGi>
        <NFormItemGi span="24 s:12 m:6" class="pr-24px">
          <NSpace class="w-full" justify="end">
            <NButton @click="reset">
              <template #icon>
                <icon-ic-round-refresh />
              </template>
              {{ $t('common.reset') }}
            </NButton>
            <NButton type="primary" ghost @click="search">
              <template #icon>
                <icon-ic-round-search />
              </template>
              {{ $t('common.search') }}
            </NButton>
          </NSpace>
        </NFormItemGi>
      </NGrid>
    </NForm>
  </NCard>
</template>

<script lang="tsx" setup>
import { $t } from '@/locales';
import { reactive } from 'vue';
import { useForm } from '@/hooks';
import type { QueryTaskLog } from '@/api/task/types';
import type { SelectOption } from 'naive-ui';
import type { TaskStatus } from '@/enum/task';

type SearchTask = Omit<QueryTaskLog, 'pageNo' | 'pageSize'>;
const emit = defineEmits<{
  (e: 'reset'): void;
  (e: 'search', searchParams: SearchTask): void;
}>();

const { formRef, validate, restoreValidation } = useForm();
const model = reactive<SearchTask>({
  taskId: '',
  taskName: '',
  status: -1
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
  emit('reset');
}

async function search() {
  await validate();
  emit('search', model);
}
</script>
