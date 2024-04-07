<template>
  <NModal
    :show="show"
    :on-update:show="updateValue"
    :title="$t('page.task.log.title')"
    preset="dialog"
    style="width: 80%; min-width: 1000px"
  >
    <NCard :bordered="true" :embedded="true" :title="$t('page.task.log.command')" class="mb-4">
      <VueJsonPretty
        :data="request"
        :deep="1"
        :show-key-value-space="true"
        :show-line="true"
        :show-double-quotes="true"
        :show-line-numbers="true"
      />
    </NCard>

    <NCard :bordered="true" :embedded="true" :title="$t('page.task.log.result')">
      <VueJsonPretty
        :data="result"
        :deep="1"
        :show-key-value-space="true"
        :show-line="true"
        :show-double-quotes="true"
      />
    </NCard>
  </NModal>
</template>

<script lang="tsx" setup>
import { $t } from '@/locales';
import { watch, ref } from 'vue';
import { fetchTaskView } from '@/api/task';
import { TaskProtocol } from '@/enum/task';
import type { TaskLogItemOutput } from '@/api/task/types';
import VueJsonPretty from 'vue-json-pretty';
import 'vue-json-pretty/lib/styles.css';

const props = defineProps<{
  show: boolean;
  taskItem: TaskLogItemOutput | null;
}>();

const emits = defineEmits<{
  (e: 'update:show', value: boolean): void;
}>();

const updateValue = (value: boolean) => {
  emits('update:show', value);
};

const result = ref({});
const request = ref({});

watch(
  () => props.taskItem,
  async () => {
    if (props.taskItem) {
      try {
        result.value = JSON.parse(props.taskItem.result);
      } catch (error) {
        result.value = props.taskItem.result;
      }
      const res = await fetchTaskView(props.taskItem.taskId);
      if (res.code === 200) {
        if (res.data.protocol === TaskProtocol.HTTP) {
          const command = JSON.parse(res.data.command);
          const params = JSON.parse(res.data.params);
          request.value = { ...command, ...params };
        } else {
          request.value = res.data.command;
        }
      }
    }
  },
  {
    immediate: true
  }
);
</script>
