<template>
  <NModal
    v-model:show="modalVisible"
    :mask-closable="false"
    preset="card"
    :title="title"
    class="task-operate-modal"
    style="width: 800px"
    :segmented="segmented"
  >
    <NForm ref="formRef" :model="model" :rules="rules" label-placement="left" :label-width="100">
      <NGrid x-gap="18" :cols="2">
        <NGridItem>
          <NFormItem :label="$t('page.task.list.form.name')" path="name">
            <NInput v-model:value="model.name" :placeholder="$t('page.task.list.form.name')" />
          </NFormItem>
        </NGridItem>
        <NGridItem>
          <NFormItem :label="$t('page.task.list.form.tag')" path="tag">
            <NInput v-model:value="model.tag" :placeholder="$t('page.task.list.form.tag')" />
          </NFormItem>
        </NGridItem>
      </NGrid>

      <NGrid x-gap="18" :cols="2">
        <NGridItem>
          <NFormItem :label="$t('page.task.list.form.protocol')" path="protocol">
            <NSelect v-model:value="model.protocol" :options="protocolOptions" @update:value="onProtocolChange" />
          </NFormItem>
        </NGridItem>
        <NGridItem>
          <NFormItem :label="$t('page.task.list.form.spec')" path="spec">
            <NInput v-model:value="model.spec" :placeholder="$t('page.task.list.form.spec')" />
          </NFormItem>
        </NGridItem>
      </NGrid>

      <NGrid x-gap="18" :cols="2">
        <NGridItem>
          <NFormItem :label="$t('page.task.list.form.policy')" path="policy">
            <NSelect v-model:value="model.policy" :options="policyOptions" />
          </NFormItem>
        </NGridItem>
        <NGridItem>
          <NFormItem :label="$t('page.task.list.form.count')" path="count">
            <NInputNumber v-model:value="model.count" :min="1" :max="999" />
          </NFormItem>
        </NGridItem>
      </NGrid>

      <NGrid x-gap="18" :cols="2">
        <NGridItem>
          <NFormItem :label="$t('page.task.list.form.timeout')" path="timeout">
            <NInputNumber v-model:value="model.timeout" :min="1" :max="86400" />
            <template #feedback>
              <span class="text-xs text-gray-500">{{ $t('page.task.list.form.timeoutUnit') }}</span>
            </template>
          </NFormItem>
        </NGridItem>
        <NGridItem>
          <NFormItem :label="$t('page.task.list.form.status')" path="status">
            <NSwitch v-model:value="model.status" :checked-value="1" :unchecked-value="0">
              <template #checked>{{ $t('common.enable') }}</template>
              <template #unchecked>{{ $t('common.disable') }}</template>
            </NSwitch>
          </NFormItem>
        </NGridItem>
      </NGrid>

      <NFormItem :label="$t('page.task.list.form.remark')" path="remark" :show-feedback="false">
        <NInput
          v-model:value="model.remark"
          type="textarea"
          :rows="2"
          :placeholder="$t('page.task.list.form.remark')"
        />
      </NFormItem>

      <!-- 协议配置区域 -->
      <NFormItem :label="$t('page.task.list.form.config')" path="config" :show-feedback="false">
        <div class="w-full">
          <!-- HTTP配置 -->
          <Http
            v-if="model.protocol === TaskProtocol.HTTP"
            :command="model.command"
            :params="model.params"
            @update="onConfigUpdate"
          />
          
          <!-- Shell配置 -->
          <NInput
            v-else-if="model.protocol === TaskProtocol.Shell"
            v-model:value="model.command"
            type="textarea"
            :rows="6"
            :placeholder="$t('page.task.list.form.shellPlaceholder')"
          />
          
          <!-- SSH配置 -->
          <Ssh
            v-else-if="model.protocol === TaskProtocol.SSH"
            :command="model.command"
            :params="model.params"
            @update="onConfigUpdate"
          />
          
        </div>
      </NFormItem>
    </NForm>

    <template #footer>
      <NSpace justify="end">
        <NButton @click="closeModal">{{ $t('common.cancel') }}</NButton>
        <NButton type="primary" :loading="submitLoading" @click="handleSubmit">
          {{ $t('common.confirm') }}
        </NButton>
      </NSpace>
    </template>
  </NModal>
</template>

<script setup lang="ts">
import { computed, reactive, ref, watch } from 'vue';
import type { FormInst, FormRules } from 'naive-ui';
import { $t } from '@/locales';
import { TaskProtocol, TaskPolicy, TaskStatus } from '@/enum/task';
import { fetchTaskEdit, fetchTaskView } from '@/api/task';
import { message } from '@/utils/message';
import Http from './http.vue';
import Ssh from './ssh.vue';

export type OperateType = 'add' | 'edit';

export interface Props {
  /** 弹框可见性 */
  visible: boolean;
  /** 操作类型 */
  operateType: OperateType;
  /** 编辑的数据id */
  dataId?: number;
}

export interface Emits {
  (e: 'update:visible', visible: boolean): void;
  (e: 'submitted'): void;
}

const props = withDefaults(defineProps<Props>(), {
  dataId: 0
});

const emit = defineEmits<Emits>();

const modalVisible = computed({
  get() {
    return props.visible;
  },
  set(visible) {
    emit('update:visible', visible);
  }
});

const title = computed(() => {
  if (props.operateType === 'add') {
    return $t('page.task.list.addTask');
  }
  return $t('page.task.list.editTask');
});

const segmented = {
  content: 'soft',
  footer: 'soft'
} as const;

const formRef = ref<HTMLElement & FormInst>();
const submitLoading = ref(false);

const model = reactive({
  id: 0,
  name: '',
  tag: '',
  protocol: TaskProtocol.HTTP,
  spec: '',
  policy: TaskPolicy.Once,
  count: 1,
  timeout: 30,
  command: '',
  params: '',
  remark: '',
  status: TaskStatus.Enabled,
  delay: 0,
  retry_times: 0,
  retry_interval: 0
});

const rules: FormRules = {
  name: [
    {
      required: true,
      message: $t('page.task.list.form.name')
    }
  ],
  protocol: [
    {
      required: true,
      type: 'number',
      message: $t('page.task.list.form.protocol')
    }
  ],
  spec: [
    {
      required: true,
      message: $t('page.task.list.form.spec')
    }
  ]
};

const protocolOptions = [
  {
    label: 'HTTP',
    value: TaskProtocol.HTTP
  },
  {
    label: 'Shell',
    value: TaskProtocol.Shell
  },
  {
    label: 'SSH',
    value: TaskProtocol.SSH
  }
];

const policyOptions = [
  {
    label: $t('page.task.list.form.policyMulti'),
    value: TaskPolicy.Multi
  },
  {
    label: $t('page.task.list.form.policyOnce'),
    value: TaskPolicy.Once
  },
  {
    label: $t('page.task.list.form.policySingle'),
    value: TaskPolicy.Single
  },
  {
    label: $t('page.task.list.form.policyTimes'),
    value: TaskPolicy.Times
  }
];

function onProtocolChange() {
  // 切换协议时清空配置
  model.command = '';
  model.params = '';
}

function onConfigUpdate(command: string, params: string) {
  model.command = command;
  model.params = params;
}

function closeModal() {
  modalVisible.value = false;
}

async function handleSubmit() {
  await formRef.value?.validate();
  
  submitLoading.value = true;
  
  try {
    if (props.operateType === 'add') {
      await fetchTaskEdit(model);
      message.success($t('common.addSuccess'));
    } else {
      await fetchTaskEdit(model);
      message.success($t('common.updateSuccess'));
    }
    
    closeModal();
    emit('submitted');
  } catch (error) {
    console.error('提交失败:', error);
  } finally {
    submitLoading.value = false;
  }
}

async function getTaskData() {
  if (props.operateType === 'edit' && props.dataId) {
    try {
      const res = await fetchTaskView(props.dataId);
      Object.assign(model, res.data);
    } catch (error) {
      console.error('获取任务数据失败:', error);
    }
  }
}

function resetForm() {
  Object.assign(model, {
    id: 0,
    name: '',
    tag: '',
    protocol: TaskProtocol.HTTP,
    spec: '',
    policy: TaskPolicy.Once,
    count: 1,
    timeout: 30,
    command: '',
    params: '',
    remark: '',
    status: TaskStatus.Enabled,
    delay: 0,
    retry_times: 0,
    retry_interval: 0
  });
}

watch(
  () => props.visible,
  (newValue) => {
    if (newValue) {
      if (props.operateType === 'add') {
        resetForm();
      } else {
        getTaskData();
      }
    }
  }
);
</script>

<style scoped>
.task-operate-modal {
  max-height: 90vh;
  overflow-y: auto;
}

:deep(.n-card__content) {
  max-height: 70vh;
  overflow-y: auto;
}

:deep(.n-form-item) {
  margin-bottom: 18px;
}

:deep(.n-form-item:last-child) {
  margin-bottom: 0;
}
</style>
