<template>
  <div class="notification-config-container">
    <NForm label-placement="left" :label-width="100">
      <!-- 通知名称 -->
      <NFormItem label="通知名称" required>
        <NInput v-model:value="model.name" placeholder="输入通知配置名称" />
      </NFormItem>

      <!-- 通知类型和触发条件 -->
      <NGrid :cols="2" :x-gap="16">
        <NGridItem>
          <NFormItem label="通知类型" required>
            <NSelect
              v-model:value="model.type"
              :options="typeOptions"
              placeholder="选择通知类型"
              @update:value="onTypeChange"
            />
          </NFormItem>
        </NGridItem>
        <NGridItem>
          <NFormItem label="触发条件" required>
            <NSelect v-model:value="model.trigger" :options="triggerOptions" placeholder="选择触发条件" />
          </NFormItem>
        </NGridItem>
      </NGrid>

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
            <div style="max-width: 450px;">
              <p class="mb-2 font-medium">可用的模板变量：</p>
              <NGrid :cols="2" :x-gap="8" :y-gap="4">
                <NGridItem v-for="item in templateVars" :key="item.var">
                  <div class="flex items-center">
                    <code class="px-2 py-1 bg-gray-100 rounded text-xs mr-1">.{{ item.var }}</code>
                    <span class="text-xs text-gray-500">{{ item.desc }}</span>
                  </div>
                </NGridItem>
              </NGrid>
              <p class="mt-2 text-xs text-gray-400">在邮件主题和正文中使用</p>
            </div>
          </NPopover>
        </NDivider>
        <NFormItem label="收件人" required>
          <NInput v-model:value="model.target" placeholder="主要收件人邮箱" />
        </NFormItem>
        <NFormItem label="抄送列表">
          <NInput v-model:value="model.email_recipients" placeholder="多个邮箱用逗号分隔" />
        </NFormItem>
        <NFormItem label="邮件主题">
          <NInput v-model:value="model.email_subject" placeholder="支持模板变量：.TaskName .Status" />
        </NFormItem>
        <NFormItem label="邮件正文">
          <NInput
            v-model:value="model.email_body"
            type="textarea"
            :rows="3"
            placeholder="支持HTML和模板变量，如：.TaskName .Status .Result"
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
            <div style="max-width: 350px;">
              <p class="mb-2 font-medium">可用的模板变量：</p>
              <NGrid :cols="2" :x-gap="8" :y-gap="4">
                <NGridItem v-for="item in templateVars" :key="item.var">
                  <div class="flex items-center">
                    <code class="px-2 py-1 bg-gray-100 rounded text-xs mr-1">.{{ item.var }}</code>
                    <span class="text-xs text-gray-500">{{ item.desc }}</span>
                  </div>
                </NGridItem>
              </NGrid>
              <p class="mt-2 text-xs text-gray-400">在请求体中使用</p>
            </div>
          </NPopover>
        </NDivider>
        <NFormItem label="Webhook URL" required>
          <NInput v-model:value="model.target" placeholder="https://example.com/webhook" />
        </NFormItem>
        <NFormItem label="请求方法">
          <NSelect v-model:value="model.webhook_method" :options="methodOptions" />
        </NFormItem>
        <NFormItem label="请求头">
          <NInput
            v-model:value="model.webhook_headers"
            type="textarea"
            :rows="2"
            placeholder='JSON格式，如：{"Content-Type": "application/json"}'
          />
        </NFormItem>
        <NFormItem label="请求体">
          <NInput
            v-model:value="model.webhook_body"
            type="textarea"
            :rows="3"
            placeholder="支持模板变量，如：.TaskName .Status .Result"
          />
        </NFormItem>
      </template>

      <!-- 重试配置 -->
      <NDivider>重试配置</NDivider>
      <NGrid :cols="2" :x-gap="16">
        <NGridItem>
          <NFormItem label="重试次数">
            <NInputNumber v-model:value="model.retry_times" :min="0" :max="10" class="w-full" />
          </NFormItem>
        </NGridItem>
        <NGridItem>
          <NFormItem label="重试间隔">
            <NInputNumber v-model:value="model.retry_interval" :min="1" :max="3600" class="w-full" />
            <template #feedback>
              <span class="text-xs text-gray-500">单位：秒</span>
            </template>
          </NFormItem>
        </NGridItem>
      </NGrid>
    </NForm>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, watch, onMounted } from 'vue';
import { NotificationType, NotificationTrigger } from '@/enum/notification';
import {
  fetchNotificationList,
  fetchNotificationCreate,
  fetchNotificationUpdate,
} from '@/api/notification';
import type { NotificationConfigInput } from '@/api/notification/types';

const props = defineProps<{
  taskId: number;
}>();

const emit = defineEmits<{
  (e: 'update', data: NotificationConfigInput | null): void;
}>();

const notificationId = ref<number | null>(null);

// 模板变量
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

// 表单模型
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

// 保存状态
const saving = ref(false);

// 选项配置
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

// 处理类型变化
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

// 加载现有通知配置
async function loadNotificationConfig() {
  if (!props.taskId) return;

  try {
    const res = await fetchNotificationList({
      pageNo: 1,
      pageSize: 1,
      task_id: props.taskId
    });

    if (res.code === 200 && res.data?.list && res.data.list.length > 0) {
      const config = res.data.list[0];
      notificationId.value = config.id;

      // 填充表单数据（使用类型断言处理可能缺少的字段）
      Object.assign(model, {
        task_id: config.task_id,
        name: config.name,
        type: config.type,
        target: config.target,
        trigger: config.trigger,
        enabled: config.enabled,
        // 这些字段在NotificationConfigOutput中可能不存在，使用默认值
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
  } catch (error) {
    console.error('加载通知配置失败:', error);
  }
}

// 监听taskId变化
watch(() => props.taskId, (newTaskId) => {
  if (newTaskId) {
    model.task_id = newTaskId;
    loadNotificationConfig();
  }
}, { immediate: true });

// 监听model变化，实时更新父组件
watch(model, () => {
  emit('update', { ...model, enabled: true });
}, { deep: true });

// 组件挂载时加载配置
onMounted(() => {
  if (props.taskId) {
    model.task_id = props.taskId;
    loadNotificationConfig();
  }
});

// 暴露保存方法
defineExpose({
  saveNotification: async (taskId: number) => {
    // 防止重复提交
    if (saving.value) return null;
    saving.value = true;

    const data = { ...model, task_id: taskId, enabled: true };

    try {
      if (notificationId.value) {
        // 更新现有配置
        const res = await fetchNotificationUpdate({ ...data, id: notificationId.value });
        if (res.code === 200) {
          return res.data;
        }
      } else {
        // 创建新配置
        const res = await fetchNotificationCreate(data);
        if (res.code === 200) {
          return res.data;
        }
      }
    } catch (error) {
      console.error('保存通知配置失败:', error);
    } finally {
      saving.value = false;
    }
    return null;
  }
});
</script>

<style scoped>
.notification-config-container {
  padding: 12px 0;
}

.w-full {
  width: 100%;
}

.text-xs {
  font-size: 12px;
}

.text-gray-500 {
  color: #64748b;
}

.text-gray-400 {
  color: #9ca3af;
}

.mb-2 {
  margin-bottom: 8px;
}

.mt-2 {
  margin-top: 8px;
}

.ml-2 {
  margin-left: 8px;
}
</style>
