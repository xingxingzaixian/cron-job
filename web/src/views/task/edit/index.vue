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

    <!-- 左右分栏：左侧基本信息，右侧执行方式配置 -->
    <div class="edit-grid">
      <!-- 基本信息 -->
      <NCard :bordered="false" size="small" class="glass-card left-col">
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
            <NInput v-model:value="model.spec" :placeholder="$t('page.task.list.form.spec')" />
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
          <NFormItem :label="$t('page.task.list.form.status')" path="status">
            <NSwitch v-model:value="model.status" :checked-value="1" :unchecked-value="0">
              <template #checked>{{ $t('common.enable') }}</template>
              <template #unchecked>{{ $t('common.disable') }}</template>
            </NSwitch>
          </NFormItem>
          <NFormItem :label="$t('page.task.list.form.remark')" path="remark" :show-feedback="false">
            <NInput v-model:value="model.remark" type="textarea" :rows="2" :placeholder="$t('page.task.list.form.remark')" />
          </NFormItem>
        </NForm>
      </NCard>

      <!-- 执行方式配置 + 测试响应 -->
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
import Ssh from '@/views/task/list/modules/ssh.vue';
import HttpBuilder from './modules/http-builder.vue';
import ResponsePanel from './modules/response-panel.vue';

defineOptions({ name: 'TaskEdit' });

const route = useRoute();
const router = useRouter();

const isEdit = computed(() => Boolean(route.query.id));

const formRef = ref<HTMLElement & FormInst>();
const submitLoading = ref(false);
const testLoading = ref(false);
const testResult = ref<TaskTestOutput | null>(null);

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
  spec: [{ required: true, message: $t('page.task.list.form.spec') }]
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
    if (isEdit.value) {
      await fetchTaskUpdate(model);
      message.success($t('task.message.editSuccess'));
    } else {
      await fetchTaskCreate(model);
      message.success($t('task.message.addSuccess'));
    }
    goBack();
  } catch (error) {
    console.error('保存失败:', error);
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
  max-width: 1280px;
  margin: 0 auto;
}

.edit-grid {
  display: grid;
  grid-template-columns: 340px 1fr;
  gap: 16px;
  align-items: start;
}

.left-col,
.right-col {
  min-width: 0;
}

@media (max-width: 1024px) {
  .edit-grid {
    grid-template-columns: 1fr;
  }
}

.page-title {
  font-size: 20px;
  font-weight: 600;
  margin: 0;
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
