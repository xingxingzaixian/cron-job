<template>
  <NCard :title="$t('common.search')">
    <NForm ref="formRef" :model="model" label-placement="left">
      <NGrid responsive="screen" item-responsive>
        <NFormItemGi span="24 s:12 m:6" :label="$t('page.task.list.name')" path="name" class="pr-24px">
          <NInput v-model:value="model.name" :placeholder="$t('page.task.list.name')" />
        </NFormItemGi>
        <NFormItemGi span="24 s:12 m:6" :label="$t('page.task.list.tag')" path="tag" class="pr-24px">
          <NInput v-model:value="model.tag" :placeholder="$t('page.task.list.tag')" />
        </NFormItemGi>
        <NFormItemGi span="24 s:12 m:6" :label="$t('page.task.list.protocol')" path="protocol" class="pr-24px">
          <NSelect
            v-model:value="model.protocol"
            :placeholder="$t('page.task.list.protocol')"
            :options="protocolOptions"
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
import type { QueryTask } from '@/api/task/types';
import type { SelectOption } from 'naive-ui';

type SearchTask = Omit<QueryTask, 'pageNo' | 'pageSize'>;
const emit = defineEmits<{
  (e: 'reset'): void;
  (e: 'search', searchParams: SearchTask): void;
}>();

const { formRef, validate, restoreValidation } = useForm();
const model = reactive<SearchTask>({
  name: '',
  tag: '',
  protocol: 0
});

const protocolOptions: SelectOption[] = [
  {
    label: $t('task.protocol.all'),
    value: 0
  },
  {
    label: $t('task.protocol.http'),
    value: 1
  },
  {
    label: $t('task.protocol.shell'),
    value: 2
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
