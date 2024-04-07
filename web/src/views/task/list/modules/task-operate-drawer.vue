<template>
  <NDrawer v-model:show="visible" :title="title" display-directive="show" :width="500">
    <NDrawerContent :title="title" :native-scrollbar="false" closable>
      <NForm ref="formRef" :model="model" :rules="rules">
        <NFormItem :label="$t('page.task.list.name')" path="name">
          <NInput v-model:value="model.name" :placeholder="$t('page.task.list.form.name')" />
        </NFormItem>
        <NFormItem :label="$t('page.task.list.tag')" path="tag">
          <NInput v-model:value="model.tag" :placeholder="$t('page.task.list.form.tag')" />
        </NFormItem>
        <NFormItem :label="$t('page.task.list.spec')" path="spec">
          <NInput v-model:value="model.spec" :placeholder="$t('page.task.list.form.specDesc')" />
        </NFormItem>
        <NFormItem :label="$t('page.task.list.protocol')" path="protocol">
          <NSelect
            v-model:value="model.protocol"
            :options="protocolOptions"
            :placeholder="$t('page.task.list.form.protocol')"
          />
        </NFormItem>
        <NFormItem v-if="model.protocol === 2" :label="$t('page.task.list.command')" path="command">
          <NInput v-model:value="model.command" :placeholder="$t('page.task.list.form.command')" type="textarea" />
        </NFormItem>
        <HttpVue v-else :command="model.command" :params="model.params" @update="updateValue" />
        <NFormItem :label="$t('page.task.list.remark')" path="remark">
          <NInput v-model:value="model.remark" :placeholder="$t('page.task.list.form.remark')" type="textarea" />
        </NFormItem>

        <NDivider>{{ $t('task.runPolicy') }}</NDivider>
        <NGrid x-gap="12" :cols="2">
          <NGridItem>
            <NFormItem :label="$t('page.task.list.policy')" path="policy">
              <NSelect
                v-model:value="model.policy"
                :options="policyOptions"
                :placeholder="$t('page.task.list.form.policy')"
              />
            </NFormItem>
          </NGridItem>
          <NGridItem>
            <NFormItem v-show="model.policy === TaskPolicy.Times" :label="$t('page.task.list.count')" path="count">
              <NInputNumber v-model:value="model.count" />
            </NFormItem>
          </NGridItem>
        </NGrid>
        <NGrid x-gap="12" :cols="2">
          <NGridItem>
            <NFormItem :label="$t('page.task.list.delay')" path="delay">
              <NInputNumber v-model:value="model.delay" />
            </NFormItem>
          </NGridItem>
          <NGridItem>
            <NFormItem :label="$t('page.task.list.timeout')" path="timeout">
              <NInputNumber v-model:value="model.timeout" />
            </NFormItem>
          </NGridItem>
        </NGrid>
        <NGrid x-gap="12" :cols="2">
          <NGridItem>
            <NFormItem :label="$t('page.task.list.retry_times')" path="retry_times">
              <NInputNumber v-model:value="model.retry_times" />
            </NFormItem>
          </NGridItem>
          <NGridItem>
            <NFormItem :label="$t('page.task.list.retry_interval')" path="retry_interval">
              <NInputNumber v-model:value="model.retry_interval" />
            </NFormItem>
          </NGridItem>
        </NGrid>
      </NForm>
      <template #footer>
        <NSpace :size="16">
          <NButton @click="closeDrawer">{{ $t('common.cancel') }}</NButton>
          <NButton type="primary" :loading="loading" @click="handleSubmit">{{ $t('common.confirm') }}</NButton>
        </NSpace>
      </template>
    </NDrawerContent>
  </NDrawer>
</template>

<script setup lang="ts">
import { $t } from '@/locales';
import type { FormRules, SelectOption } from 'naive-ui';
import { message } from '@/utils/message';
import { computed, reactive, watch, ref } from 'vue';
import type { TaskEditHTTPInput } from '@/api/task/types';
import { fetchTaskView, fetchTaskEdit } from '@/api/task';
import HttpVue from './http.vue';
import { TaskPolicy, TaskProtocol } from '@/enum/task';
import { useForm } from '@/hooks';

defineOptions({
  name: 'TaskOperateDrawer'
});

interface Props {
  /** the type of operation */
  operateType: OperateType;
  /** the edit row data */
  dataId?: number | null;
}
const props = defineProps<Props>();
interface Emits {
  (e: 'submitted'): void;
}
const emit = defineEmits<Emits>();
export type OperateType = 'add' | 'edit';
const loading = ref<boolean>(false);
const visible = defineModel<boolean>('visible', {
  default: false
});

const { formRef, validate, restoreValidation } = useForm();

const protocolOptions: SelectOption[] = [
  {
    label: $t('task.protocol.http'),
    value: 1
  },
  {
    label: $t('task.protocol.shell'),
    value: 2
  }
  // {
  //   label: $t('task.protocol.grpc'),
  //   value: 3
  // }
];

const policyOptions: SelectOption[] = [
  {
    label: $t('task.policy.multi'),
    value: 1
  },
  {
    label: $t('task.policy.once'),
    value: 2
  },
  {
    label: $t('task.policy.single'),
    value: 3
  },
  {
    label: $t('task.policy.times'),
    value: 4
  }
];
const rules: FormRules = {};
const title = computed(() => {
  const titles: Record<OperateType, string> = {
    add: $t('page.task.list.addTask'),
    edit: $t('page.task.list.editTask')
  };
  return titles[props.operateType];
});

async function handleSubmit() {
  await validate();
  // request
  loading.value = true;
  const res = await fetchTaskEdit(model);
  loading.value = false;
  if (res.code !== 200) {
    res.msg && message.error(res.msg);
    return;
  }

  message.success($t('common.updateSuccess'));
  closeDrawer();
  emit('submitted');
}

const updateValue = (command: string, params: string) => {
  model.command = command;
  model.params = params;
};

const model: TaskEditHTTPInput = reactive(createDefaultModel());
function createDefaultModel(): TaskEditHTTPInput {
  return {
    command: '',
    count: 0,
    delay: 0,
    id: null,
    name: '',
    params: '',
    policy: TaskPolicy.Single,
    protocol: TaskProtocol.Shell,
    remark: '',
    retry_interval: 0,
    retry_times: 0,
    spec: '',
    tag: '',
    timeout: 0
  };
}

async function handleUpdateModelWhenEdit() {
  if (props.operateType === 'add') {
    Object.assign(model, createDefaultModel());
    return;
  }

  if (props.operateType === 'edit' && props.dataId) {
    fetchTaskView(props.dataId).then((res) => {
      if (res.code === 200) {
        Object.assign(model, res.data);
      }
    });
  }
}

function closeDrawer() {
  visible.value = false;
}

watch(visible, () => {
  if (visible.value) {
    handleUpdateModelWhenEdit();
    restoreValidation();
  }
});
</script>

<style scoped></style>
