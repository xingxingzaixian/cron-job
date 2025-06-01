<template>
  <NModal
    :show="show"
    :on-update:show="updateValue"
    :title="$t('page.task.log.title')"
    preset="dialog"
    style="width: 90%; min-width: 1200px; max-width: 1600px"
  >
    <!-- 任务命令显示 -->
    <NCard bordered embedded class="mb-4 command-card">
      <div class="command-header">
        <span class="command-type">{{ getProtocolName(props.taskItem?.protocol) }}</span>
        <NButton size="small" quaternary @click="copyToClipboard(request, 'command')">
          <template #icon>
            <SvgIcon icon="material-symbols:content-copy" />
          </template>
          复制
        </NButton>
      </div>
      <VueJsonPretty
        v-if="props.taskItem?.protocol === TaskProtocol.HTTP"
        :data="request"
        :deep="1"
        show-key-value-space
        show-line
        show-double-quotes
        show-line-numbers
      />
      <div v-else class="command-display">
        <pre class="command-content"><code>{{ request }}</code></pre>
      </div>
    </NCard>

    <!-- 执行结果显示 -->
    <NCard bordered embedded class="result-card">
      <div class="result-header">
        <div class="result-info">
          <NBadge :type="getResultType()" :value="getResultStatus()" class="status-badge" />
          <span class="result-lines">{{ getLineCount() }} 行</span>
        </div>
        <NButton size="small" quaternary @click="copyToClipboard(result, 'result')">
          <template #icon>
            <SvgIcon icon="material-symbols:content-copy" />
          </template>
          复制
        </NButton>
      </div>
      <VueJsonPretty
        v-if="props.taskItem?.protocol === TaskProtocol.HTTP"
        :data="result"
        :deep="1"
        show-key-value-space
        show-line
        show-double-quotes
        show-line-numbers
      />
      <div v-else class="result-display">
        <div class="result-content-wrapper">
          <div class="line-numbers">
            <div v-for="(line, index) in getResultLines()" :key="index" class="line-number">
              {{ index + 1 }}
            </div>
          </div>
          <pre class="result-content"><code v-for="(line, index) in getResultLines()"
            :key="index"
            :class="['result-line', getLineClass(line)]">{{ line || ' ' }}</code></pre>
        </div>
      </div>
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
import { message } from '@/utils/message';
import SvgIcon from '@/components/custom/SvgIcon.vue';

const props = defineProps<{
  show: boolean;
  taskItem: TaskLogItemOutput | null;
}>();

const emits = defineEmits<(e: 'update:show', value: boolean) => void>();

const updateValue = (value: boolean) => {
  emits('update:show', value);
};

const result = ref({});
const request = ref({});

// 获取协议名称
const getProtocolName = (protocol?: number) => {
  switch (protocol) {
    case TaskProtocol.HTTP:
      return 'HTTP';
    case TaskProtocol.Shell:
      return 'Shell';
    case TaskProtocol.SSH:
      return 'SSH';
    default:
      return 'Unknown';
  }
};

// 复制到剪贴板
const copyToClipboard = async (content: any, type: string) => {
  try {
    const text = typeof content === 'string' ? content : JSON.stringify(content, null, 2);
    await navigator.clipboard.writeText(text);
    message.success(`${type === 'command' ? '命令' : '结果'}已复制到剪贴板`);
  } catch (err) {
    message.error('复制失败');
  }
};

// 获取结果行数组
const getResultLines = () => {
  if (typeof result.value === 'string') {
    return result.value.split('\n');
  }
  return [JSON.stringify(result.value, null, 2)];
};

// 获取行数
const getLineCount = () => {
  return getResultLines().length;
};

// 获取结果类型（用于状态徽章）
const getResultType = () => {
  const resultStr = typeof result.value === 'string' ? result.value : JSON.stringify(result.value);
  if (resultStr.toLowerCase().includes('error') || resultStr.toLowerCase().includes('failed')) {
    return 'error';
  }
  if (resultStr.toLowerCase().includes('warning') || resultStr.toLowerCase().includes('warn')) {
    return 'warning';
  }
  return 'success';
};

// 获取结果状态文本
const getResultStatus = () => {
  const type = getResultType();
  switch (type) {
    case 'error':
      return '执行失败';
    case 'warning':
      return '有警告';
    default:
      return '执行成功';
  }
};

// 获取行的CSS类
const getLineClass = (line: string) => {
  const lowerLine = line.toLowerCase();
  if (lowerLine.includes('error') || lowerLine.includes('failed') || lowerLine.includes('exception')) {
    return 'error-line';
  }
  if (lowerLine.includes('warning') || lowerLine.includes('warn')) {
    return 'warning-line';
  }
  if (lowerLine.includes('success') || lowerLine.includes('ok') || lowerLine.includes('done')) {
    return 'success-line';
  }
  return '';
};

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

<style scoped>
/* 命令显示区域 */
.command-card {
  border: 1px solid #e5e7eb;
  border-radius: 12px;
  overflow: hidden;
}

.command-display {
  background: #f8fafc;
  border-radius: 8px;
  overflow: hidden;
}

.command-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px;
  background: #f1f5f9;
  border-bottom: 1px solid #e2e8f0;
}

.command-type {
  font-size: 12px;
  font-weight: 600;
  color: #475569;
  background: #e2e8f0;
  padding: 4px 8px;
  border-radius: 4px;
}

.command-content {
  margin: 0;
  padding: 16px;
  background: #ffffff;
  font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
  font-size: 13px;
  line-height: 1.5;
  color: #374151;
  white-space: pre-wrap;
  word-break: break-all;
  border: none;
  overflow-x: auto;
}

/* 结果显示区域 */
.result-card {
  border: 1px solid #e5e7eb;
  border-radius: 12px;
  overflow: hidden;
}

.result-display {
  background: #f8fafc;
  border-radius: 8px;
  overflow: hidden;
}

.result-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px;
  background: #f1f5f9;
  border-bottom: 1px solid #e2e8f0;
}

.result-info {
  display: flex;
  align-items: center;
  gap: 12px;
}

.status-badge {
  font-size: 11px;
}

.result-lines {
  font-size: 12px;
  color: #64748b;
  font-weight: 500;
}

.result-content-wrapper {
  display: flex;
  background: #ffffff;
  max-height: 500px;
  overflow: auto;
}

.line-numbers {
  background: #f8fafc;
  border-right: 1px solid #e2e8f0;
  padding: 16px 8px;
  min-width: 50px;
  text-align: right;
  user-select: none;
  flex-shrink: 0;
}

.line-number {
  font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
  font-size: 12px;
  line-height: 1.5;
  color: #9ca3af;
  height: 19.5px;
  display: flex;
  align-items: center;
  justify-content: flex-end;
}

.result-content {
  margin: 0;
  padding: 16px;
  flex: 1;
  font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
  font-size: 13px;
  line-height: 1.5;
  color: #374151;
  white-space: pre;
  border: none;
  overflow: visible;
}

.result-line {
  display: block;
  min-height: 19.5px;
  margin: 0;
  padding: 0;
  border: none;
  background: transparent;
}

/* 不同类型行的颜色 */
.error-line {
  color: #dc2626;
  background-color: rgba(220, 38, 38, 0.05);
  padding: 0 4px;
  margin: 0 -4px;
  border-radius: 2px;
}

.warning-line {
  color: #d97706;
  background-color: rgba(217, 119, 6, 0.05);
  padding: 0 4px;
  margin: 0 -4px;
  border-radius: 2px;
}

.success-line {
  color: #059669;
  background-color: rgba(5, 150, 105, 0.05);
  padding: 0 4px;
  margin: 0 -4px;
  border-radius: 2px;
}

/* 滚动条样式 */
.result-content-wrapper::-webkit-scrollbar {
  width: 8px;
  height: 8px;
}

.result-content-wrapper::-webkit-scrollbar-track {
  background: #f1f5f9;
}

.result-content-wrapper::-webkit-scrollbar-thumb {
  background: #cbd5e1;
  border-radius: 4px;
}

.result-content-wrapper::-webkit-scrollbar-thumb:hover {
  background: #94a3b8;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .command-header,
  .result-header {
    flex-direction: column;
    gap: 8px;
    align-items: flex-start;
  }

  .result-info {
    flex-direction: column;
    align-items: flex-start;
    gap: 4px;
  }

  .line-numbers {
    min-width: 40px;
    padding: 16px 4px;
  }

  .command-content,
  .result-content {
    font-size: 12px;
    padding: 12px;
  }
}
</style>
