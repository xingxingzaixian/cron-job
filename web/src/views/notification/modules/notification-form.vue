<template>
  <NModal v-model:show="showModal" preset="card" :title="isEdit ? '编辑通知配置' : '新增通知配置'" style="width: 720px">
    <NForm ref="formRef" :model="model" :rules="rules" label-placement="left" label-width="100">
      <NGrid :cols="2" :x-gap="16">
        <NGridItem>
          <NFormItem label="任务" path="task_id">
            <NSelect
              v-model:value="model.task_id"
              :options="taskOptions"
              placeholder="请选择任务"
              :loading="taskListLoading"
              filterable
              @update:value="onTaskChange"
            />
          </NFormItem>
        </NGridItem>
        <NGridItem>
          <NFormItem label="通知名称" path="name">
            <NInput v-model:value="model.name" placeholder="请输入通知名称" />
          </NFormItem>
        </NGridItem>
      </NGrid>

      <!-- 全局配置说明 -->
      <NAlert v-if="model.task_id === 0" type="info" :show-icon="true" class="mb-4">
        <template #header>全局配置说明</template>
        当前为全局通知配置，该通知将应用于<strong>所有任务</strong>。当任意任务触发通知条件时，都会发送此通知。如果某个任务同时配置了全局通知和专属通知，两种通知都会被触发。
      </NAlert>

      <NGrid :cols="2" :x-gap="16">
        <NGridItem>
          <NFormItem label="通知类型" path="type">
            <NSelect
              v-model:value="model.type"
              :options="typeOptions"
              placeholder="请选择通知类型"
              @update:value="onTypeChange"
            />
          </NFormItem>
        </NGridItem>
        <NGridItem>
          <NFormItem label="触发条件" path="trigger">
            <NSelect v-model:value="model.trigger" :options="triggerOptions" placeholder="请选择触发条件" />
          </NFormItem>
        </NGridItem>
      </NGrid>

      <NFormItem label="状态" path="enabled">
        <NSwitch v-model:value="model.enabled">
          <template #checked>{{ $t('common.enable') }}</template>
          <template #unchecked>{{ $t('common.disable') }}</template>
        </NSwitch>
      </NFormItem>

      <!-- 邮件配置 -->
      <template v-if="model.type === NotificationType.Email">
        <NDivider>
          邮件配置
          <NPopover trigger="hover" :show-arrow="true" :delay="300" :duration="200">
            <template #trigger>
              <NButton text type="info" class="ml-2">
                <template #icon>
                  <NIcon>
                    <svg viewBox="0 0 24 24">
                      <path fill="currentColor" d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm1 15h-2v-6h2v6zm0-8h-2V7h2v2z"/>
                    </svg>
                  </NIcon>
                </template>
              </NButton>
            </template>
            <div style="max-width: 300px;">
              <p class="mb-2 font-medium">模板变量说明：</p>
              <p class="text-sm">在邮件主题和正文中，可以使用以下变量：</p>
              <ul class="mt-1 text-sm list-disc list-inside">
                <li><code>.TaskName</code> - 任务名称</li>
                <li><code>.Status</code> - 执行状态</li>
                <li><code>.Command</code> - 执行命令</li>
                <li><code>.Result</code> - 执行结果</li>
                <li><code>.StartTime</code> - 开始时间</li>
                <li><code>.EndTime</code> - 结束时间</li>
                <li><code>.Duration</code> - 执行时长</li>
              </ul>
              <p class="mt-2 text-xs text-gray-400">示例：任务 .TaskName 执行 .Status</p>
            </div>
          </NPopover>
        </NDivider>
        <NFormItem label="收件人" path="target">
          <NInput v-model:value="model.target" placeholder="主要收件人邮箱" />
        </NFormItem>
        <NFormItem label="抄送列表" path="email_recipients">
          <NInput v-model:value="model.email_recipients" placeholder="多个邮箱用逗号分隔" />
        </NFormItem>
        <NFormItem label="邮件主题" path="email_subject">
          <NInput v-model:value="model.email_subject" placeholder="支持模板变量：TaskName Status" />
        </NFormItem>
        <NFormItem label="邮件正文" path="email_body">
          <NInput
            v-model:value="model.email_body"
            type="textarea"
            :rows="4"
            placeholder="支持HTML和模板变量"
          />
        </NFormItem>
      </template>

      <!-- Webhook配置 -->
      <template v-if="model.type === NotificationType.Webhook">
        <NDivider>
          Webhook配置
          <NPopover trigger="hover" :show-arrow="true" :delay="300" :duration="200">
            <template #trigger>
              <NButton text type="info" class="ml-2">
                <template #icon>
                  <NIcon>
                    <svg viewBox="0 0 24 24">
                      <path fill="currentColor" d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm1 15h-2v-6h2v6zm0-8h-2V7h2v2z"/>
                    </svg>
                  </NIcon>
                </template>
              </NButton>
            </template>
            <div style="max-width: 300px;">
              <p class="mb-2 font-medium">Webhook 说明：</p>
              <p class="text-sm">Webhook 是一种 HTTP 回调机制，当任务执行完成时，系统会向指定的 URL 发送 POST 请求。</p>
              <p class="mt-2 text-sm">常用场景：</p>
              <ul class="mt-1 text-sm list-disc list-inside">
                <li>钉钉机器人</li>
                <li>企业微信</li>
                <li>飞书机器人</li>
                <li>自定义 API</li>
              </ul>
              <p class="mt-2 text-xs text-gray-400">请求体支持模板变量</p>
            </div>
          </NPopover>
        </NDivider>
        <NFormItem label="Webhook URL" path="webhook_url">
          <NInput v-model:value="model.webhook_url" placeholder="https://example.com/webhook" />
        </NFormItem>
        <NFormItem label="请求方法" path="webhook_method">
          <NSelect v-model:value="model.webhook_method" :options="methodOptions" />
        </NFormItem>
        <NFormItem label="请求头" path="webhook_headers">
          <NInput
            v-model:value="model.webhook_headers"
            type="textarea"
            :rows="3"
            placeholder='JSON格式，如：{"Content-Type": "application/json"}'
          />
        </NFormItem>
        <NFormItem label="请求体" path="webhook_body">
          <NInput
            v-model:value="model.webhook_body"
            type="textarea"
            :rows="4"
            placeholder="支持模板变量"
          />
        </NFormItem>
      </template>

      <NDivider>重试配置</NDivider>
      <NGrid :cols="2" :x-gap="16">
        <NGridItem>
          <NFormItem label="重试次数" path="retry_times">
            <NInputNumber v-model:value="model.retry_times" :min="0" :max="10" class="w-full" />
          </NFormItem>
        </NGridItem>
        <NGridItem>
          <NFormItem label="重试间隔" path="retry_interval">
            <NInputNumber v-model:value="model.retry_interval" :min="1" :max="3600" class="w-full" />
            <template #feedback>
              <span class="text-xs text-gray-500">单位：秒</span>
            </template>
          </NFormItem>
        </NGridItem>
      </NGrid>

      <!-- 模板变量说明 -->
      <NDivider>
        模板变量说明
        <NPopover trigger="hover" :show-arrow="true" :delay="300" :duration="200">
          <template #trigger>
            <NButton text type="info" class="ml-2">
              <template #icon>
                <NIcon>
                  <svg viewBox="0 0 24 24">
                    <path fill="currentColor" d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm1 15h-2v-6h2v6zm0-8h-2V7h2v2z"/>
                  </svg>
                </NIcon>
              </template>
            </NButton>
          </template>
          <div style="max-width: 400px;">
            <p class="mb-2 font-medium">可用的模板变量：</p>
            <NGrid :cols="2" :x-gap="8" :y-gap="4">
              <NGridItem v-for="item in templateVars" :key="item.var">
                <div class="flex items-center">
                  <code class="px-2 py-1 bg-gray-100 rounded text-xs mr-1">.{{ item.var }}</code>
                  <span class="text-xs text-gray-500">{{ item.desc }}</span>
                </div>
              </NGridItem>
            </NGrid>
            <p class="mt-2 text-xs text-gray-400">在邮件主题、正文或 Webhook 请求体中使用</p>
          </div>
        </NPopover>
      </NDivider>
      <NAlert type="info" :show-icon="true" :bordered="false" class="template-vars-alert">
        <template #header>可用的模板变量</template>
        <p class="mb-2 text-sm text-gray-600">在邮件主题、正文或 Webhook 请求体中，您可以使用以下模板变量：</p>
        <NGrid :cols="2" :x-gap="12" :y-gap="4">
          <NGridItem v-for="item in templateVars" :key="item.var">
            <div class="template-var-item">
              <code class="mr-1 px-2 py-1 bg-gray-100 rounded text-sm">.{{ item.var }}</code>
              <span class="text-xs text-gray-500">{{ item.desc }}</span>
            </div>
          </NGridItem>
        </NGrid>
        <p class="mt-3 text-xs text-gray-400">使用示例：在模板中输入 <code>.TaskName</code> 即可获取任务名称</p>
      </NAlert>
    </NForm>

    <template #footer>
      <NSpace justify="end">
        <NButton @click="handleCancel">{{ $t('common.cancel') }}</NButton>
        <NButton type="primary" :loading="testLoading" @click="handleTest">
          测试发送
        </NButton>
        <NButton type="primary" :loading="submitLoading" @click="handleSave">
          {{ $t('common.save') }}
        </NButton>
      </NSpace>
    </template>
  </NModal>
</template>

<script lang="tsx" setup>
import { ref, reactive, watch, onMounted } from 'vue';
import type { FormInst, FormRules } from 'naive-ui';
import { $t } from '@/locales';
import { message } from '@/utils/message';
import { useForm } from '@/hooks';
import { NotificationType, NotificationTrigger } from '@/enum/notification';
import {
  fetchNotificationView,
  fetchNotificationCreate,
  fetchNotificationUpdate,
  fetchNotificationTest,
} from '@/api/notification';
import { fetchTaskList } from '@/api/task';
import type { NotificationConfigInput } from '@/api/notification/types';
import type { TaskItemOutput } from '@/api/task/types';

const props = defineProps<{
  show: boolean;
  notificationId?: number | null;
  taskId?: number;
}>();

const emit = defineEmits<{
  (e: 'update:show', value: boolean): void;
  (e: 'success'): void;
}>();

const showModal = ref(props.show);
const submitLoading = ref(false);
const testLoading = ref(false);
const isEdit = ref(false);

const { formRef, validate, restoreValidation } = useForm();

// ==================== 任务列表 ====================
const taskList = ref<TaskItemOutput[]>([]);
const taskListLoading = ref(false);
const taskOptions = ref<{ label: string; value: number }[]>([
  { label: '全局配置（所有任务）', value: 0 },
]);

/** 从接口加载任务列表 */
async function loadTaskList() {
  taskListLoading.value = true;
  try {
    const res = await fetchTaskList({ pageNo: 1, pageSize: 999 });
    if (res.code === 200 && res.data?.list) {
      taskList.value = res.data.list;
      taskOptions.value = [
        { label: '🌍 全局配置（所有任务）', value: 0 },
        ...res.data.list.map((t) => ({
          label: `${t.name}（ID: ${t.id}）`,
          value: t.id,
        })),
      ];
    }
  } catch (error) {
    console.error('获取任务列表失败:', error);
  } finally {
    taskListLoading.value = false;
  }
}

function onTaskChange(value: number) {
  model.task_id = value;
}

// ==================== 模板变量 ====================
const templateVars = [
  { var: 'TaskID', desc: '任务ID' },
  { var: 'TaskName', desc: '任务名称' },
  { var: 'Status', desc: '执行状态（成功/失败）' },
  { var: 'Command', desc: '执行命令' },
  { var: 'Result', desc: '执行结果' },
  { var: 'RetryTimes', desc: '重试次数' },
  { var: 'StartTime', desc: '开始时间' },
  { var: 'EndTime', desc: '结束时间' },
  { var: 'Duration', desc: '执行时长' },
];

// ==================== 表单模型 ====================
const model = reactive<NotificationConfigInput>({
  task_id: 0,
  name: '',
  type: NotificationType.Email,
  target: '',
  trigger: NotificationTrigger.All,
  enabled: true,
  email_recipients: '',
  email_subject: '',
  email_body: '',
  webhook_url: '',
  webhook_method: 'POST',
  webhook_headers: '',
  webhook_body: '',
  retry_times: 0,
  retry_interval: 5,
});

const rules: FormRules = {
  name: [{ required: true, message: '请输入通知名称' }],
  type: [{ required: true, message: '请选择通知类型' }],
  target: [{ required: true, message: '请输入目标地址' }],
  trigger: [{ required: true, message: '请选择触发条件' }],
};

const typeOptions = [
  { label: '邮件通知', value: NotificationType.Email },
  { label: 'Webhook通知', value: NotificationType.Webhook },
];

const triggerOptions = [
  { label: '所有情况', value: NotificationTrigger.All },
  { label: '成功时', value: NotificationTrigger.Success },
  { label: '失败时', value: NotificationTrigger.Failure },
  { label: '超时时', value: NotificationTrigger.Timeout },
  { label: '取消时', value: NotificationTrigger.Cancel },
];

const methodOptions = [
  { label: 'POST', value: 'POST' },
  { label: 'PUT', value: 'PUT' },
  { label: 'PATCH', value: 'PATCH' },
];

function onTypeChange() {
  // 重置相关字段
  model.target = '';
  model.email_recipients = '';
  model.email_subject = '';
  model.email_body = '';
  model.webhook_url = '';
  model.webhook_method = 'POST';
  model.webhook_headers = '';
  model.webhook_body = '';
}

function handleCancel() {
  emit('update:show', false);
  restoreValidation();
}

async function handleTest() {
  try {
    await validate();
    testLoading.value = true;
    const res = await fetchNotificationTest({
      type: model.type,
      target: model.target,
      content: '这是一条测试消息',
    });
    if (res.code === 200) {
      message.success('测试发送成功');
    } else {
      message.error(res.message || '测试发送失败');
    }
  } catch (error) {
    console.error('测试发送失败:', error);
    message.error('测试发送失败');
  } finally {
    testLoading.value = false;
  }
}

async function handleSave() {
  try {
    await validate();
    submitLoading.value = true;

    if (isEdit.value && props.notificationId) {
      const res = await fetchNotificationUpdate({
        ...model,
        id: props.notificationId,
      });
      if (res.code === 200) {
        message.success($t('common.updateSuccess'));
        emit('update:show', false);
        emit('success');
      } else {
        message.error(res.message || $t('common.updateFailed'));
      }
    } else {
      const res = await fetchNotificationCreate(model);
      if (res.code === 200) {
        message.success($t('common.addSuccess'));
        emit('update:show', false);
        emit('success');
      } else {
        message.error(res.message || $t('common.addFailed'));
      }
    }
  } catch (error) {
    console.error('保存失败:', error);
    message.error(isEdit.value ? $t('common.updateFailed') : $t('common.addFailed'));
  } finally {
    submitLoading.value = false;
  }
}

async function getNotificationData() {
  if (props.notificationId) {
    isEdit.value = true;
    try {
      const res = await fetchNotificationView(props.notificationId);
      if (res.code === 200 && res.data) {
        Object.assign(model, res.data);
      }
    } catch (error) {
      console.error('获取通知配置失败:', error);
    }
  } else {
    isEdit.value = false;
    // 重置表单
    Object.assign(model, {
      task_id: props.taskId || 0,
      name: '',
      type: NotificationType.Email,
      target: '',
      trigger: NotificationTrigger.All,
      enabled: true,
      email_recipients: '',
      email_subject: '',
      email_body: '',
      webhook_url: '',
      webhook_method: 'POST',
      webhook_headers: '',
      webhook_body: '',
      retry_times: 0,
      retry_interval: 5,
    });
  }
}

watch(
  () => props.show,
  (val) => {
    showModal.value = val;
    if (val) {
      getNotificationData();
    }
  }
);

watch(
  () => showModal.value,
  (val) => {
    emit('update:show', val);
  }
);

// 组件挂载时加载任务列表
onMounted(() => {
  loadTaskList();
});
</script>

<style scoped>
.w-full {
  width: 100%;
}

.text-xs {
  font-size: 12px;
}

.text-gray-500 {
  color: #64748b;
}

.text-gray-600 {
  color: #4b5563;
}

.text-gray-400 {
  color: #9ca3af;
}

.text-sm {
  font-size: 14px;
}

.mb-2 {
  margin-bottom: 8px;
}

.mb-4 {
  margin-bottom: 16px;
}

.mt-1 {
  margin-top: 4px;
}

.mt-3 {
  margin-top: 12px;
}

.mr-1 {
  margin-right: 4px;
}

.template-var-item {
  display: flex;
  align-items: center;
  padding: 4px 0;
}

.template-vars-alert {
  background-color: rgba(32, 128, 240, 0.04);
}
</style>
