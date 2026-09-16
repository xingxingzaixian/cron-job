<template>
  <div class="task-edit-page">
    <!-- 页面头部 -->
    <div class="flex items-center justify-between mb-16px">
      <div class="flex items-center gap-8px">
        <NButton quaternary circle @click="goBack">
          <template #icon>
            <icon-ic-round-keyboard-backspace />
          </template>
        </NButton>
        <h1 class="page-title">{{ isEdit ? $t('page.task.list.editTask') : $t('page.task.list.addTask') }}</h1>
      </div>
    </div>

    <!-- 左右分栏：左侧基本信息+通知配置，右侧执行方式配置 -->
    <div class="edit-grid">
      <!-- 左侧：基本信息 + 通知配置 -->
      <div class="left-col">
        <!-- 基本信息 -->
        <NCard :bordered="false" size="small" class="glass-card">
          <template #header>
            <span class="text-base font-medium">{{ $t('page.task.list.basicInfoTitle') }}</span>
          </template>
          <NForm ref="formRef" :model="model" :rules="rules" label-placement="top">
            <NFormItem :label="$t('page.task.list.form.name')" path="name">
              <NInput v-model:value="model.name" :placeholder="$t('page.task.list.form.name')" />
            </NFormItem>
            <NFormItem :label="$t('page.task.list.form.tag')" path="tag">
              <NInput v-model:value="model.tag" :placeholder="$t('page.task.list.form.tag')" />
            </NFormItem>
            <NFormItem :label="$t('page.task.list.form.spec')" path="spec">
              <NSpace vertical class="w-full">
                <NSpace>
                  <NSelect
                    v-model:value="cronPreset"
                    :options="cronPresetOptions"
                    :placeholder="$t('page.task.list.form.cronPreset')"
                    size="small"
                    class="cron-preset-select"
                    @update:value="applyCronPreset"
                  />
                  <NInput
                    v-model:value="model.spec"
                    :placeholder="$t('page.task.list.form.specPlaceholder')"
                    @update:value="onSpecChange"
                  />
                </NSpace>
                <NText v-if="specError" depth="3" style="color: #e88080; font-size: 12px">
                  {{ specError }}
                </NText>
                <NText v-else-if="specBreakdown" depth="3" style="font-size: 12px">
                  {{ specBreakdown }}
                </NText>
              </NSpace>
            </NFormItem>
            <NGrid :cols="2" :x-gap="12">
              <NGridItem>
                <NFormItem :label="$t('page.task.list.form.protocol')" path="protocol">
                  <NSelect v-model:value="model.protocol" :options="protocolOptions" @update:value="onProtocolChange" />
                </NFormItem>
              </NGridItem>
              <NGridItem>
                <NFormItem :label="$t('page.task.list.form.policy')" path="policy">
                  <NSelect v-model:value="model.policy" :options="policyOptions" />
                </NFormItem>
              </NGridItem>
            </NGrid>
            <NGrid :cols="2" :x-gap="12">
              <NGridItem>
                <NFormItem :label="$t('page.task.list.form.count')" path="count">
                  <NInputNumber v-model:value="model.count" :min="1" :max="999" class="w-full" />
                </NFormItem>
              </NGridItem>
              <NGridItem>
                <NFormItem :label="$t('page.task.list.form.timeout')" path="timeout">
                  <NInputNumber v-model:value="model.timeout" :min="1" :max="86400" class="w-full" />
                  <template #feedback>
                    <span class="text-xs text-gray-500">{{ $t('page.task.list.form.timeoutUnit') }}</span>
                  </template>
                </NFormItem>
              </NGridItem>
            </NGrid>
            <NGrid :cols="2" :x-gap="12">
              <NGridItem>
                <NFormItem :label="$t('page.task.list.form.status')" path="status">
                  <NSwitch v-model:value="statusEnabled">
                    <template #checked>{{ $t('common.enable') }}</template>
                    <template #unchecked>{{ $t('common.disable') }}</template>
                  </NSwitch>
                </NFormItem>
              </NGridItem>
              <NGridItem>
                <NFormItem label="通知配置">
                  <div class="flex items-center gap-8px">
                    <NSwitch v-model:value="notificationEnabled" @update:value="onNotificationSwitchChange">
                      <template #checked>启用</template>
                      <template #unchecked>禁用</template>
                    </NSwitch>
                  </div>
                </NFormItem>
              </NGridItem>
            </NGrid>
            <NFormItem :label="$t('page.task.list.form.remark')" path="remark" :show-feedback="false">
              <NInput v-model:value="model.remark" type="textarea" :rows="2" :placeholder="$t('page.task.list.form.remark')" />
            </NFormItem>
          </NForm>
        </NCard>
      </div>

      <!-- 右侧：执行方式配置 + 测试响应 -->
      <div class="right-col">
        <NCard :bordered="false" size="small" class="glass-card">
          <template #header>
            <div class="flex items-center justify-between w-full">
              <span class="text-base font-medium">{{ $t('page.task.list.configTitle') }}</span>
              <NButton
                v-if="model.protocol !== TaskProtocol.HTTP"
                type="primary"
                :loading="testLoading"
                @click="handleTest"
              >
                {{ $t('page.task.list.testButton') }}
              </NButton>
            </div>
          </template>

          <!-- HTTP：Hoppscotch 式请求构建器 -->
          <HttpBuilder
            v-if="model.protocol === TaskProtocol.HTTP"
            :command="model.command"
            :params="model.params"
            :testing="testLoading"
            @update="onConfigUpdate"
            @test="handleTest"
          />

          <!-- Shell：命令文本 -->
          <NInput
            v-else-if="model.protocol === TaskProtocol.Shell"
            v-model:value="model.command"
            type="textarea"
            :rows="10"
            :placeholder="$t('page.task.list.form.shellPlaceholder')"
          />

          <!-- SSH：连接与命令配置 -->
          <Ssh
            v-else-if="model.protocol === TaskProtocol.SSH"
            :command="model.command"
            :params="model.params"
            @update="onConfigUpdate"
          />
        </NCard>

        <!-- 测试响应 -->
        <ResponsePanel :result="testResult" />
      </div>
    </div>

    <!-- 底部固定操作栏 -->
    <div class="edit-footer">
      <NSpace>
        <NButton @click="goBack">{{ $t('common.cancel') }}</NButton>
        <NButton type="primary" :loading="submitLoading" @click="handleSave">
          {{ $t('page.task.list.saveButton') }}
        </NButton>
      </NSpace>
    </div>

    <!-- 通知配置弹窗 -->
    <NModal v-model:show="showNotificationModal" preset="card" title="通知配置" style="width: 680px">
      <NotificationConfig 
        ref="notificationConfigRef"
        :task-id="model.id"
        @update="onNotificationUpdate"
      />
      <template #footer>
        <NSpace justify="end">
          <NButton @click="onNotificationModalClose">{{ $t('common.close') }}</NButton>
        </NSpace>
      </template>
    </NModal>
  </div>
</template>

<script setup lang="ts">
import { computed, reactive, ref } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import type { FormInst, FormRules } from 'naive-ui';
import { $t } from '@/locales';
import { message } from '@/utils/message';
import { TaskProtocol, TaskPolicy, TaskStatus } from '@/enum/task';
import { fetchTaskCreate, fetchTaskUpdate, fetchTaskView, fetchTaskTest } from '@/api/task';
import type { TaskTestOutput } from '@/api/task/types';
import { CRON_PRESETS, validateCronSpec, cronSpecBreakdown } from '@/utils/cron';
import Ssh from '@/views/task/list/modules/ssh.vue';
import HttpBuilder from './modules/http-builder.vue';
import ResponsePanel from './modules/response-panel.vue';
import NotificationConfig from './modules/notification-config.vue';
import type { NotificationConfigInput } from '@/api/notification/types';

defineOptions({ name: 'TaskEdit' });

const route = useRoute();
const router = useRouter();

const isEdit = computed(() => Boolean(route.query.id));

const formRef = ref<HTMLElement & FormInst>();
const notificationConfigRef = ref<InstanceType<typeof NotificationConfig>>();
const submitLoading = ref(false);
const testLoading = ref(false);
const showNotificationModal = ref(false);
const notificationEnabled = ref(false);
const testResult = ref<TaskTestOutput | null>(null);
const cronPreset = ref<string>('custom');
const notificationData = ref<NotificationConfigInput | null>(null);

const cronPresetOptions = [
  { label: $t('page.task.list.form.cronPresetCustom'), value: 'custom' },
  ...CRON_PRESETS.map((preset) => ({ label: preset.label, value: preset.value }))
];

const specError = computed(() => validateCronSpec(model.spec));
const specBreakdown = computed(() => cronSpecBreakdown(model.spec));

// 状态开关：与列表页保持一致，凡非"禁用"状态均视为启用。
// 任务执行后会变为 运行中/成功/失败 等状态，因此不能只用 1 判定。
const statusEnabled = computed<boolean>({
  get: () => model.status !== TaskStatus.Disabled,
  set: (val: boolean) => {
    model.status = val ? TaskStatus.Enabled : TaskStatus.Disabled;
  }
});

function applyCronPreset(value: string) {
  if (value !== 'custom') {
    model.spec = value;
  }
}

function onSpecChange() {
  // 手动编辑时切回"自定义"，避免预设选中态误导
  cronPreset.value = 'custom';
}

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
  name: [{ required: true, message: $t('page.task.list.form.name') }],
  protocol: [{ required: true, type: 'number', message: $t('page.task.list.form.protocol') }],
  spec: [
    { required: true, message: $t('page.task.list.form.spec') },
    {
      validator: (_rule: unknown, value: string) => {
        const err = validateCronSpec(value);
        return err ? new Error(err) : true;
      }
    }
  ]
};

const protocolOptions = [
  { label: 'HTTP', value: TaskProtocol.HTTP },
  { label: 'Shell', value: TaskProtocol.Shell },
  { label: 'SSH', value: TaskProtocol.SSH }
];

const policyOptions = [
  { label: $t('page.task.list.form.policyMulti'), value: TaskPolicy.Multi },
  { label: $t('page.task.list.form.policyOnce'), value: TaskPolicy.Once },
  { label: $t('page.task.list.form.policySingle'), value: TaskPolicy.Single },
  { label: $t('page.task.list.form.policyTimes'), value: TaskPolicy.Times }
];

const onProtocolChange = () => {
  model.command = '';
  model.params = '';
  testResult.value = null;
};

const onConfigUpdate = (command: string, params: string) => {
  model.command = command;
  model.params = params;
};

const onNotificationUpdate = (data: NotificationConfigInput | null) => {
  notificationData.value = data;
  notificationEnabled.value = data !== null && data.enabled;
};

const onNotificationSwitchChange = (enabled: boolean) => {
  if (enabled) {
    showNotificationModal.value = true;
  } else {
    notificationData.value = null;
  }
};

const onNotificationModalClose = () => {
  showNotificationModal.value = false;
};

const goBack = () => {
  router.push({ name: 'TaskList' });
};

async function handleTest() {
  if (model.protocol === TaskProtocol.HTTP) {
    try {
      const cmd = JSON.parse(model.command || '{}');
      if (!/^https?:\/\//i.test(cmd.url || '')) {
        message.warning($t('http.urlInvalid'));
        return;
      }
    } catch (error) {
      message.warning($t('http.urlInvalid'));
      return;
    }
  } else if (model.protocol === TaskProtocol.Shell) {
    if (!model.command.trim()) {
      message.warning($t('page.task.list.testShellCommandEmpty'));
      return;
    }
  } else if (model.protocol === TaskProtocol.SSH) {
    try {
      const cfg = JSON.parse(model.params || '{}');
      if (!cfg.host || !cfg.username) {
        message.warning($t('page.task.list.testSshConfigEmpty'));
        return;
      }
    } catch (error) {
      message.warning($t('page.task.list.testSshConfigEmpty'));
      return;
    }
  }

  testLoading.value = true;
  testResult.value = null;
  try {
    const res = await fetchTaskTest({
      protocol: model.protocol,
      command: model.command,
      params: model.params,
      timeout: model.timeout
    });
    testResult.value = res.data;
  } catch (error) {
    console.error('测试执行失败:', error);
  } finally {
    testLoading.value = false;
  }
}

async function handleSave() {
  await formRef.value?.validate();

  submitLoading.value = true;
  try {
    let taskId = model.id;
    
    if (isEdit.value) {
      const res = await fetchTaskUpdate(model);
      if (res.code !== 200) {
        message.error(res.message || $t('task.message.editFailed'));
        return;
      }
      message.success($t('task.message.editSuccess'));
    } else {
      const res = await fetchTaskCreate(model);
      if (res.code !== 200) {
        message.error(res.message || $t('task.message.addFailed'));
        return;
      }
      // 创建成功后从返回数据中获取任务ID，用于关联通知配置
      taskId = Number(res.data?.id) || model.id;
      message.success($t('task.message.addSuccess'));
    }
    
    // 保存通知配置
    if (notificationConfigRef.value && taskId) {
      await notificationConfigRef.value.saveNotification(taskId);
    }
    
    goBack();
  } catch (error) {
    console.error('保存失败:', error);
    message.error(isEdit.value ? $t('task.message.editFailed') : $t('task.message.addFailed'));
  } finally {
    submitLoading.value = false;
  }
}

async function getTaskData() {
  if (isEdit.value && route.query.id) {
    try {
      const res = await fetchTaskView(Number(route.query.id));
      Object.assign(model, res.data);
    } catch (error) {
      console.error('获取任务数据失败:', error);
    }
  }
}

getTaskData();
</script>

<style scoped>
.task-edit-page {
  padding: 16px 24px 88px;
  max-width: 1400px;
  margin: 0 auto;
}

.edit-grid {
  display: grid;
  grid-template-columns: 380px 1fr;
  gap: 16px;
  align-items: start;
}

.left-col,
.right-col {
  min-width: 0;
}

.left-col {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

@media (max-width: 1024px) {
  .edit-grid {
    grid-template-columns: 1fr;
  }
  
  .left-col {
    grid-column: 1 / -1;
  }
}

.page-title {
  font-size: 20px;
  font-weight: 600;
  margin: 0;
}

.cron-preset-select {
  width: 130px;
  flex-shrink: 0;
}

.mt-4 {
  margin-top: 16px;
}

.edit-footer {
  position: sticky;
  bottom: 0;
  z-index: 10;
  display: flex;
  justify-content: flex-end;
  padding: 12px 0;
  margin-top: 16px;
  background: var(--n-color, rgba(255, 255, 255, 0.9));
  backdrop-filter: blur(8px);
  border-top: 1px solid rgba(0, 0, 0, 0.06);
}
</style>
