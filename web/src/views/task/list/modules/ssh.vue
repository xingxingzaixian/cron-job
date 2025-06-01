<template>
  <NSpace vertical class="w-full">
    <!-- SSH连接配置 -->
    <NFormItem path="connection" :show-feedback="false">
      <NCard size="small" class="ssh-config-card">
        <template #header>
          <span class="config-header">{{ $t('ssh.connectionConfig') }}</span>
        </template>
        
        <NGrid x-gap="12" :y-gap="8" :cols="2">
          <NGridItem>
            <NFormItem :label="$t('ssh.host')" path="host" :show-feedback="false">
              <NInput 
                v-model:value="sshConfig.host" 
                :placeholder="$t('ssh.hostPlaceholder')"
                @input="onConfigChange"
              />
            </NFormItem>
          </NGridItem>
          <NGridItem>
            <NFormItem :label="$t('ssh.port')" path="port" :show-feedback="false">
              <NInputNumber 
                v-model:value="sshConfig.port" 
                :placeholder="22"
                :min="1"
                :max="65535"
                @update:value="onConfigChange"
              />
            </NFormItem>
          </NGridItem>
          <NGridItem>
            <NFormItem :label="$t('ssh.username')" path="username" :show-feedback="false">
              <NInput 
                v-model:value="sshConfig.username" 
                :placeholder="$t('ssh.usernamePlaceholder')"
                @input="onConfigChange"
              />
            </NFormItem>
          </NGridItem>
          <NGridItem>
            <NFormItem :label="$t('ssh.password')" path="password" :show-feedback="false">
              <NInput 
                v-model:value="sshConfig.password" 
                type="password"
                show-password-on="click"
                :placeholder="$t('ssh.passwordPlaceholder')"
                @input="onConfigChange"
              />
            </NFormItem>
          </NGridItem>
          <NGridItem>
            <NFormItem :label="$t('ssh.mode')" path="mode" :show-feedback="false">
              <NSelect
                v-model:value="sshConfig.mode"
                :options="modeOptions"
                @update:value="onConfigChange"
              />
            </NFormItem>
          </NGridItem>
          <NGridItem>
            <NSpace>
              <NButton :loading="testing" @click="testConnection">
                {{ $t('ssh.testConnection') }}
              </NButton>
              <NButton @click="showExample">
                {{ $t('ssh.showExample') }}
              </NButton>
            </NSpace>
          </NGridItem>
        </NGrid>
      </NCard>
    </NFormItem>

    <!-- 执行命令 -->
    <NFormItem path="commands" :show-feedback="false">
      <NCard size="small" class="ssh-commands-card">
        <template #header>
          <span class="config-header">{{ $t('ssh.commandsConfig') }}</span>
        </template>
        
        <NInput 
          v-model:value="commands" 
          type="textarea"
          :rows="8"
          :placeholder="$t('ssh.commandsPlaceholder')"
          @input="onCommandsChange"
        />
        
        <template #footer>
          <NText depth="3" style="font-size: 12px;color: antiquewhite;">
            {{ $t('ssh.commandsHelp') }}
          </NText>
        </template>
      </NCard>
    </NFormItem>

    <!-- 配置预览 -->
    <NCollapse v-if="showPreview">
      <NCollapseItem :title="$t('ssh.configPreview')" name="preview">
        <NCode :code="configPreview" language="json" />
      </NCollapseItem>
    </NCollapse>
  </NSpace>
</template>

<script setup lang="ts">
import { $t } from '@/locales';
import { reactive, watch, ref, computed } from 'vue';
import type { SelectOption } from 'naive-ui';
import { message } from '@/utils/message';

const props = defineProps<{
  command: string;
  params: string;
}>();

const emit = defineEmits<(e: 'update', command: string, params: string) => void>();

// SSH配置模型
const sshConfig = reactive({
  host: '',
  port: 22,
  username: '',
  password: '',
  mode: 'sequential'
});

// 执行命令
const commands = ref('');

// 测试连接状态
const testing = ref(false);

// 显示预览
const showPreview = ref(false);

// 执行模式选项
const modeOptions: SelectOption[] = [
  {
    label: $t('ssh.sequentialMode'),
    value: 'sequential'
  },
  {
    label: $t('ssh.scriptMode'),
    value: 'script'
  }
];

// 配置预览
const configPreview = computed(() => {
  return JSON.stringify({
    params: sshConfig,
    command: commands.value
  }, null, 2);
});

// 配置变化
const onConfigChange = () => {
  updateData();
};

// 命令变化
const onCommandsChange = () => {
  updateData();
};

// 更新数据
const updateData = () => {
  const params = JSON.stringify(sshConfig);
  const command = commands.value;
  
  emit('update', command, params);
};

// 测试连接
const testConnection = async () => {
  if (!sshConfig.host || !sshConfig.username) {
    message.warning($t('ssh.fillRequiredFields'));
    return;
  }
  
  if (!sshConfig.password) {
    message.warning($t('ssh.fillAuthInfo'));
    return;
  }
  
  testing.value = true;
  
  // 模拟测试连接
  setTimeout(() => {
    testing.value = false;
    message.success($t('ssh.connectionSuccess'));
  }, 2000);
};

// 显示示例
const showExample = () => {
  // 填充示例数据
  sshConfig.host = '192.168.1.100';
  sshConfig.port = 22;
  sshConfig.username = 'root';
  sshConfig.password = 'your_password';
  sshConfig.mode = 'sequential';
  
  commands.value = `whoami
pwd
ls -la
df -h
# 这是注释，会被忽略
uptime`;

  showPreview.value = true;
  
  updateData();
  message.info($t('ssh.exampleLoaded'));
};

// 监听props变化
watch(
  () => props.params,
  () => {
    try {
      if (props.params) {
        const config = JSON.parse(props.params);
        Object.assign(sshConfig, config);
      }
    } catch (err) {
      console.log('解析SSH配置失败:', err);
    }
  },
  {
    immediate: true
  }
);

watch(
  () => props.command,
  () => {
    commands.value = props.command || '';
  },
  {
    immediate: true
  }
);
</script>

<style scoped>
.ssh-config-card {
  border: 1px solid var(--n-border-color);
}

.ssh-commands-card {
  border: 1px solid var(--n-border-color);
}

.config-header {
  font-weight: 600;
  color: var(--n-text-color);
}

:deep(.n-card-header) {
  text-align: center;
  padding: 2px 10px;
  border-bottom: 1px solid var(--n-divider-color);
}

:deep(.n-card__content) {
  padding: 6px 16px;
}

:deep(.n-card__footer) {
  padding: 2px 16px;
  border-top: 1px solid var(--n-divider-color);
  background-color: var(--n-color-target);
}

:deep(.n-form-item) {
  margin-bottom: 16px;
}

:deep(.n-form-item:last-child) {
  margin-bottom: 0;
}
</style>
