<template>
  <div class="install-container">
    <div class="install-card">
      <div class="install-header">
        <h1 class="install-title">
          <n-icon size="32" class="title-icon">
            <svg viewBox="0 0 24 24">
              <path fill="currentColor" d="M12,2A10,10 0 0,0 2,12A10,10 0 0,0 12,22A10,10 0 0,0 22,12A10,10 0 0,0 12,2M11,16.5L6.5,12L7.91,10.59L11,13.67L16.59,8.09L18,9.5L11,16.5Z"/>
            </svg>
          </n-icon>
          系统安装向导
        </h1>
        <p class="install-description">欢迎使用定时任务管理系统，请按照向导完成系统初始化配置</p>
      </div>

      <n-steps :current="currentStep" :status="stepStatus">
        <n-step title="基础配置" description="配置HTTP端口和数据库连接" />
        <n-step title="认证设置" description="配置用户认证和管理员账户" />
        <n-step title="完成安装" description="确认配置并完成安装" />
      </n-steps>

      <div class="install-content">
        <!-- 步骤1: 基础配置 -->
        <div v-if="currentStep === 1" class="step-content">
          <n-form ref="basicFormRef" :model="basicForm" :rules="basicRules" label-placement="top">
            <!-- HTTP端口配置 -->
            <n-form-item label="HTTP服务端口" path="httpPort">
              <n-input-number 
                v-model:value="basicForm.httpPort" 
                :min="1" 
                :max="65535" 
                placeholder="8210" 
                class="w-full"
                :precision="0"
              />
              <template #suffix>
                <n-text depth="3">默认端口：8210</n-text>
              </template>
            </n-form-item>
            <n-alert type="info" class="mb-4">
              修改端口后，您需要使用新端口访问系统。确保端口未被其他服务占用。
            </n-alert>
            
            <!-- 数据库配置 -->
            <n-divider title-placement="left">数据库配置</n-divider>
            
            <n-form-item label="数据库类型" path="database.engine">
              <n-select
                v-model:value="basicForm.database.engine"
                :options="dbEngineOptions"
                placeholder="请选择数据库类型"
                @update:value="handleDbEngineChange"
              />
            </n-form-item>

            <!-- MySQL/PostgreSQL 配置 -->
            <template v-if="['mysql', 'postgresql'].includes(basicForm.database.engine)">
              <n-grid :cols="2" :x-gap="16">
                <n-grid-item>
                  <n-form-item label="主机地址" path="database.host">
                    <n-input v-model:value="basicForm.database.host" placeholder="127.0.0.1" />
                  </n-form-item>
                </n-grid-item>
                <n-grid-item>
                  <n-form-item label="端口" path="database.port">
                    <n-input-number v-model:value="basicForm.database.port" :min="1" :max="65535" class="w-full" />
                  </n-form-item>
                </n-grid-item>
              </n-grid>
              
              <n-form-item label="数据库名" path="database.name">
                <n-input v-model:value="basicForm.database.name" placeholder="cronJob" />
              </n-form-item>
              
              <n-grid :cols="2" :x-gap="16">
                <n-grid-item>
                  <n-form-item label="用户名" path="database.user">
                    <n-input v-model:value="basicForm.database.user" placeholder="root" />
                  </n-form-item>
                </n-grid-item>
                <n-grid-item>
                  <n-form-item label="密码" path="database.password">
                    <n-input v-model:value="basicForm.database.password" type="password" show-password-on="click" placeholder="数据库密码" />
                  </n-form-item>
                </n-grid-item>
              </n-grid>
            </template>

            <!-- SQLite 配置 -->
            <template v-if="basicForm.database.engine === 'sqlite'">
              <n-form-item label="数据库文件" path="database.name">
                <n-input v-model:value="basicForm.database.name" placeholder="cronJob" />
                <template #suffix>.db</template>
              </n-form-item>
              <n-alert type="info" class="mt-2">
                SQLite 数据库文件将保存在 data 目录下
              </n-alert>
            </template>
          </n-form>

          <div class="step-actions">
            <n-space>
              <n-button @click="testDbConnection" :loading="testing" type="primary" ghost>
                测试数据库连接
              </n-button>
              <n-button @click="nextStep" :disabled="!dbTestPassed" type="primary">
                下一步
              </n-button>
            </n-space>
          </div>
        </div>

        <!-- 步骤2: 认证设置 -->
        <div v-if="currentStep === 2" class="step-content">
          <n-form ref="authFormRef" :model="authForm" :rules="authRules" label-placement="top">
            <n-form-item label="用户认证">
              <n-switch v-model:value="authForm.authEnabled" @update:value="handleAuthToggle">
                <template #checked>启用认证</template>
                <template #unchecked>禁用认证</template>
              </n-switch>
              <n-text depth="3" class="ml-2">
                {{ authForm.authEnabled ? '启用用户登录认证功能' : '禁用认证，直接访问系统功能' }}
              </n-text>
            </n-form-item>

            <!-- 管理员账户配置 -->
            <template v-if="authForm.authEnabled">
              <n-divider title-placement="left">管理员账户</n-divider>
              
              <n-grid :cols="2" :x-gap="16">
                <n-grid-item>
                  <n-form-item label="用户名" path="adminUser.username">
                    <n-input v-model:value="authForm.adminUser.username" placeholder="admin" />
                  </n-form-item>
                </n-grid-item>
                <n-grid-item>
                  <n-form-item label="昵称" path="adminUser.nickname">
                    <n-input v-model:value="authForm.adminUser.nickname" placeholder="系统管理员" />
                  </n-form-item>
                </n-grid-item>
              </n-grid>
              
              <n-form-item label="密码" path="adminUser.password">
                <n-input v-model:value="authForm.adminUser.password" type="password" show-password-on="click" placeholder="管理员密码" />
              </n-form-item>
              
              <n-form-item label="确认密码" path="adminUser.confirmPassword">
                <n-input v-model:value="authForm.adminUser.confirmPassword" type="password" show-password-on="click" placeholder="再次输入密码" />
              </n-form-item>
              
              <n-form-item label="邮箱" path="adminUser.email">
                <n-input v-model:value="authForm.adminUser.email" placeholder="admin@example.com（可选）" />
              </n-form-item>
            </template>
          </n-form>

          <div class="step-actions">
            <n-space>
              <n-button @click="prevStep">上一步</n-button>
              <n-button @click="nextStep" type="primary">下一步</n-button>
            </n-space>
          </div>
        </div>

        <!-- 步骤3: 完成安装 -->
        <div v-if="currentStep === 3" class="step-content">
          <n-alert type="info" title="安装确认" class="mb-4">
            请确认以下配置信息，点击"开始安装"完成系统初始化
          </n-alert>

          <n-descriptions bordered :column="1">
            <n-descriptions-item label="HTTP端口">
              {{ basicForm.httpPort }}
            </n-descriptions-item>
            <n-descriptions-item label="数据库类型">
              {{ getDbEngineLabel(basicForm.database.engine) }}
            </n-descriptions-item>
            <n-descriptions-item v-if="['mysql', 'postgresql'].includes(basicForm.database.engine)" label="连接地址">
              {{ basicForm.database.host }}:{{ basicForm.database.port }}
            </n-descriptions-item>
            <n-descriptions-item label="数据库名">
              {{ basicForm.database.name }}
            </n-descriptions-item>
            <n-descriptions-item v-if="['mysql', 'postgresql'].includes(basicForm.database.engine)" label="用户名">
              {{ basicForm.database.user }}
            </n-descriptions-item>
            <n-descriptions-item label="用户认证">
              {{ authForm.authEnabled ? '已启用' : '已禁用' }}
            </n-descriptions-item>
            <n-descriptions-item v-if="authForm.authEnabled" label="管理员账户">
              {{ authForm.adminUser.username }} ({{ authForm.adminUser.nickname }})
            </n-descriptions-item>
          </n-descriptions>

          <div class="step-actions">
            <n-space>
              <n-button @click="prevStep">上一步</n-button>
              <n-button @click="startInstall" :loading="installing" type="primary">
                开始安装
              </n-button>
            </n-space>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { useMessage } from 'naive-ui';
import { checkInstall, testDatabase, install } from '@/api/install';
import type { DatabaseTestParams, InstallParams } from '@/api/install/types';

const router = useRouter();
const message = useMessage();

// 当前步骤
const currentStep = ref(1);
const stepStatus = ref<'process' | 'finish' | 'error' | 'wait'>('process');

// 表单引用
const basicFormRef = ref<any>();
const authFormRef = ref<any>();

// 状态
const testing = ref(false);
const installing = ref(false);
const dbTestPassed = ref(false);

// 基础配置表单
const basicForm = reactive({
  httpPort: 8210,
  database: {
    engine: 'mysql',
    host: '127.0.0.1',
    port: 3306,
    name: 'cronJob',
    user: 'root',
    password: '',
    path: ''
  }
});

// 认证表单
const authForm = reactive({
  authEnabled: true,
  adminUser: {
    username: 'admin',
    nickname: '系统管理员',
    password: '',
    confirmPassword: '',
    email: ''
  }
});

// 数据库引擎选项
const dbEngineOptions = [
  { label: 'MySQL', value: 'mysql' },
  { label: 'PostgreSQL', value: 'postgresql' },
  { label: 'SQLite', value: 'sqlite' }
];

// 表单验证规则
const basicRules = {
  httpPort: [
    { required: true, type: 'number' as const, message: '请输入HTTP端口' },
    { type: 'number' as const, min: 1, max: 65535, message: '端口范围应在1-65535之间' }
  ],
  'database.engine': { required: true, message: '请选择数据库类型' },
  'database.host': { required: true, message: '请输入主机地址' },
  'database.port': { required: true, type: 'number' as const, message: '请输入端口号' },
  'database.name': { required: true, message: '请输入数据库名' },
  'database.user': { required: true, message: '请输入用户名' }
};

const authRules = {
  'adminUser.username': { required: true, message: '请输入用户名' },
  'adminUser.nickname': { required: true, message: '请输入昵称' },
  'adminUser.password': { required: true, message: '请输入密码' },
  'adminUser.confirmPassword': [
    { required: true, message: '请确认密码' },
    {
      validator: (_rule: any, value: string) => {
        return value === authForm.adminUser.password;
      },
      message: '两次输入的密码不一致',
      trigger: 'blur' as const
    }
  ]
};

// 获取数据库引擎标签
const getDbEngineLabel = (engine: string) => {
  const option = dbEngineOptions.find(opt => opt.value === engine);
  return option?.label || engine;
};

// 处理数据库引擎变化
const handleDbEngineChange = (value: string) => {
  if (value === 'mysql') {
    basicForm.database.port = 3306;
  } else if (value === 'postgresql') {
    basicForm.database.port = 5432;
  }
  dbTestPassed.value = false;
};

// 处理认证开关
const handleAuthToggle = (value: boolean) => {
  if (!value) {
    // 禁用认证时清空管理员信息
    Object.assign(authForm.adminUser, {
      username: '',
      nickname: '',
      password: '',
      confirmPassword: '',
      email: ''
    });
  } else {
    // 启用认证时设置默认值
    Object.assign(authForm.adminUser, {
      username: 'admin',
      nickname: '系统管理员',
      password: '',
      confirmPassword: '',
      email: ''
    });
  }
};

// 测试数据库连接
const testDbConnection = async () => {
  try {
    await basicFormRef.value?.validate();
  } catch (error) {
    console.error('表单验证失败:', error);
    return;
  }

  testing.value = true;
  try {
    const params: DatabaseTestParams = {
      engine: basicForm.database.engine,
      name: basicForm.database.name
    };

    if (['mysql', 'postgresql'].includes(basicForm.database.engine)) {
      params.host = basicForm.database.host;
      params.port = basicForm.database.port;
      params.user = basicForm.database.user;
      params.password = basicForm.database.password;
    }

    const response = await testDatabase(params);
    const { data } = response;
    if (data.success) {
      message.success('数据库连接测试成功');
      dbTestPassed.value = true;
    } else {
      message.error(data.message || '数据库连接测试失败');
      dbTestPassed.value = false;
    }
  } catch (error: any) {
    message.error(error?.message || '数据库连接测试失败');
    dbTestPassed.value = false;
  } finally {
    testing.value = false;
  }
};

// 下一步
const nextStep = async () => {
  if (currentStep.value === 1) {
    if (!dbTestPassed.value) {
      message.warning('请先测试数据库连接');
      return;
    }
  } else if (currentStep.value === 2) {
    if (authForm.authEnabled) {
      try {
        await authFormRef.value?.validate();
      } catch (error) {
        console.error('认证表单验证失败:', error);
        return;
      }
    }
  }
  
  currentStep.value++;
};

// 上一步
const prevStep = () => {
  currentStep.value--;
};

// 开始安装
const startInstall = async () => {
  installing.value = true;
  try {
    const params: InstallParams = {
      database: {
        engine: basicForm.database.engine,
        name: basicForm.database.name
      },
      auth_enabled: authForm.authEnabled,
      http_port: basicForm.httpPort
    };

    if (['mysql', 'postgresql'].includes(basicForm.database.engine)) {
      params.database.host = basicForm.database.host;
      params.database.port = basicForm.database.port;
      params.database.user = basicForm.database.user;
      params.database.password = basicForm.database.password;
    }

    if (authForm.authEnabled) {
      params.admin_user = {
        username: authForm.adminUser.username,
        nickname: authForm.adminUser.nickname,
        password: authForm.adminUser.password,
        email: authForm.adminUser.email
      };
    }

    const response = await install(params);
    const { data } = response;
    if (data.success) {
      message.success(data.message || '系统安装成功');
      stepStatus.value = 'finish';
      
      // 显示重启提示
      message.info('服务器正在重启，请稍候...', { duration: 5000 });
      
      // 等待服务器重启完成后刷新页面
      setTimeout(() => {
        // 轮询检查服务器是否重启完成
        checkServerRestart();
      }, 5000);
    } else {
      message.error(data.message || '安装失败');
      stepStatus.value = 'error';
    }
  } catch (error: any) {
    message.error(error?.message || '安装失败');
    stepStatus.value = 'error';
  } finally {
    installing.value = false;
  }
};

// 检查安装状态
const checkInstallStatus = async () => {
  try {
    const response = await checkInstall();
    const { data } = response;
    if (data.installed) {
      // 已安装，跳转到主页
      router.push('/');
    }
  } catch (error) {
    console.error('检查安装状态失败:', error);
  }
};

// 检查服务器重启状态
const checkServerRestart = async () => {
  let attempts = 0;
  const maxAttempts = 30; // 最多尝试30次，每次间隔2秒
  
  const poll = async () => {
    attempts++;
    try {
      const response = await checkInstall();
      const { data } = response;
      if (data.installed) {
        // 服务器已重启并完成安装
        message.success('服务器重启完成，即将跳转到系统');
        setTimeout(() => {
          window.location.href = '/'; // 跳转到首页
        }, 1000);
        return;
      }
    } catch (_error) {
      // 服务器可能还在重启中，继续等待
    }
    
    if (attempts < maxAttempts) {
      setTimeout(poll, 2000); // 2秒后再次检查
    } else {
      message.warning('服务器重启时间较长，请手动刷新页面');
    }
  };
  
  poll();
};

onMounted(() => {
  checkInstallStatus();
});
</script>

<style scoped>
.install-container {
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 20px;
}

.install-card {
  background: white;
  border-radius: 12px;
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.1);
  width: 100%;
  max-width: 800px;
  padding: 40px;
}

.install-header {
  text-align: center;
  margin-bottom: 40px;
}

.install-title {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 12px;
  font-size: 28px;
  font-weight: 600;
  color: #333;
  margin: 0 0 12px 0;
}

.title-icon {
  color: #18a058;
}

.install-description {
  color: #666;
  font-size: 16px;
  margin: 0;
}

.install-content {
  margin-top: 40px;
}

.step-content {
  min-height: 400px;
}

.step-actions {
  margin-top: 40px;
  padding-top: 20px;
  border-top: 1px solid #f0f0f0;
  display: flex;
  justify-content: center;
}

.w-full {
  width: 100%;
}

.mt-2 {
  margin-top: 8px;
}

.mb-4 {
  margin-bottom: 16px;
}

.ml-2 {
  margin-left: 8px;
}
</style>