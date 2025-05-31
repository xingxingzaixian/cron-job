<template>
  <div class="login-container">
    <!-- 背景装饰 -->
    <div class="login-bg">
      <div class="bg-shape shape-1"></div>
      <div class="bg-shape shape-2"></div>
      <div class="bg-shape shape-3"></div>
      <div class="particles">
        <div v-for="i in 20" :key="i" class="particle"></div>
      </div>
    </div>

    <!-- 登录表单 -->
    <div class="login-wrapper">
      <div class="login-form">
        <!-- Logo和标题 -->
        <div class="login-header">
          <div class="logo">
            <div class="logo-icon">
              <svg viewBox="0 0 24 24" width="48" height="48">
                <path fill="currentColor" d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-2 15l-5-5 1.41-1.41L10 14.17l7.59-7.59L19 8l-9 9z"/>
              </svg>
            </div>
          </div>
          <h1 class="title">定时任务管理系统</h1>
          <p class="subtitle">欢迎回来，请登录您的账户</p>
        </div>

        <!-- 登录表单 -->
        <form class="form" @submit.prevent="onSubmit">
          <div class="form-group">
            <div class="input-wrapper">
              <div class="input-icon">
                <svg viewBox="0 0 24 24" width="20" height="20">
                  <path fill="currentColor" d="M12 12c2.21 0 4-1.79 4-4s-1.79-4-4-4-4 1.79-4 4 1.79 4 4 4zm0 2c-2.67 0-8 1.34-8 4v2h16v-2c0-2.66-5.33-4-8-4z"/>
                </svg>
              </div>
              <input
                v-model="form.userName"
                type="text"
                name="username"
                :placeholder="$t('page.login.accout')"
                autocomplete="username"
                required
              />
            </div>
          </div>

          <div class="form-group">
            <div class="input-wrapper">
              <div class="input-icon">
                <svg viewBox="0 0 24 24" width="20" height="20">
                  <path fill="currentColor" d="M18 8h-1V6c0-2.76-2.24-5-5-5S7 3.24 7 6v2H6c-1.1 0-2 .9-2 2v10c0 1.1.9 2 2 2h12c1.1 0 2-.9 2-2V10c0-1.1-.9-2-2-2zM12 17c-1.1 0-2-.9-2-2s.9-2 2-2 2 .9 2 2-.9 2-2 2zM15.1 8H8.9V6c0-1.71 1.39-3.1 3.1-3.1 1.71 0 3.1 1.39 3.1 3.1v2z"/>
                </svg>
              </div>
              <input
                v-model="form.password"
                type="password"
                name="password"
                :placeholder="$t('page.login.pass')"
                autocomplete="current-password"
                required
              />
            </div>
          </div>

          <button type="submit" class="login-btn" :disabled="loading">
            <span v-if="loading" class="loading-spinner"></span>
            <span>{{ loading ? '登录中...' : $t('page.login.submit') }}</span>
          </button>
        </form>

        <!-- 默认账户提示 -->
        <div class="login-tips">
          <div class="tip-card">
            <div class="tip-icon">
              <svg viewBox="0 0 24 24" width="16" height="16">
                <path fill="currentColor" d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm1 15h-2v-6h2v6zm0-8h-2V7h2v2z"/>
              </svg>
            </div>
            <div class="tip-content">
              <div class="tip-title">默认账户</div>
              <div class="tip-text">用户名：<code>admin</code> 密码：<code>admin</code></div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue';
import type { LoginInput } from '@/api/user/model';
import { $t } from '@/locales';
import useUserStore from '@/store/modules/user';
import { useRouter } from 'vue-router';
import { message } from '@/utils/message';

const userStore = useUserStore();
const router = useRouter();
const loading = ref(false);

const form = reactive<LoginInput>({
  userName: 'admin',
  password: 'admin'
});

const onSubmit = async () => {
  if (loading.value) return;

  if (!form.userName.trim() || !form.password.trim()) {
    message.error('请输入用户名和密码');
    return;
  }

  loading.value = true;

  try {
    const res = await userStore.Login(form);
    if (res.code === 200) {
      message.success('登录成功');
      const redirect = router.currentRoute.value.query.redirect as string;
      router.push(redirect || { name: 'Task' });
    } else {
      message.error(res.message || $t('page.login.failed'));
    }
  } catch (error) {
    message.error('登录失败，请检查网络连接');
  } finally {
    loading.value = false;
  }
};
</script>

<style scoped>
/* 登录容器 */
.login-container {
  position: relative;
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  overflow: hidden;
}

/* 背景装饰 */
.login-bg {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  pointer-events: none;
}

.bg-shape {
  position: absolute;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.1);
  animation: float 6s ease-in-out infinite;
}

.shape-1 {
  width: 200px;
  height: 200px;
  top: 10%;
  left: 10%;
  animation-delay: 0s;
}

.shape-2 {
  width: 150px;
  height: 150px;
  top: 60%;
  right: 10%;
  animation-delay: 2s;
}

.shape-3 {
  width: 100px;
  height: 100px;
  bottom: 20%;
  left: 20%;
  animation-delay: 4s;
}

@keyframes float {
  0%, 100% { transform: translateY(0px) rotate(0deg); }
  50% { transform: translateY(-20px) rotate(180deg); }
}

/* 粒子效果 */
.particles {
  position: absolute;
  width: 100%;
  height: 100%;
}

.particle {
  position: absolute;
  width: 4px;
  height: 4px;
  background: rgba(255, 255, 255, 0.6);
  border-radius: 50%;
  animation: particle-float 8s linear infinite;
}

.particle:nth-child(odd) {
  animation-duration: 10s;
  background: rgba(255, 255, 255, 0.3);
}

.particle:nth-child(1) { left: 10%; animation-delay: 0s; }
.particle:nth-child(2) { left: 20%; animation-delay: 1s; }
.particle:nth-child(3) { left: 30%; animation-delay: 2s; }
.particle:nth-child(4) { left: 40%; animation-delay: 3s; }
.particle:nth-child(5) { left: 50%; animation-delay: 4s; }
.particle:nth-child(6) { left: 60%; animation-delay: 5s; }
.particle:nth-child(7) { left: 70%; animation-delay: 6s; }
.particle:nth-child(8) { left: 80%; animation-delay: 7s; }
.particle:nth-child(9) { left: 90%; animation-delay: 8s; }
.particle:nth-child(10) { left: 15%; animation-delay: 9s; }

@keyframes particle-float {
  0% { transform: translateY(100vh) rotate(0deg); opacity: 0; }
  10% { opacity: 1; }
  90% { opacity: 1; }
  100% { transform: translateY(-100px) rotate(360deg); opacity: 0; }
}

/* 登录表单容器 */
.login-wrapper {
  position: relative;
  z-index: 1;
}

.login-form {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-radius: 20px;
  padding: 40px;
  width: 420px;
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
}

/* 登录头部 */
.login-header {
  text-align: center;
  margin-bottom: 40px;
}

.logo {
  margin-bottom: 20px;
}

.logo-icon {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 80px;
  height: 80px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 20px;
  color: white;
  box-shadow: 0 10px 30px rgba(102, 126, 234, 0.3);
}

.title {
  font-size: 28px;
  font-weight: 700;
  color: #2d3748;
  margin: 0 0 8px 0;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.subtitle {
  font-size: 16px;
  color: #718096;
  margin: 0;
}

/* 表单样式 */
.form {
  margin-bottom: 30px;
}

.form-group {
  margin-bottom: 20px;
}

.input-wrapper {
  position: relative;
  display: flex;
  align-items: center;
}

.input-icon {
  position: absolute;
  left: 16px;
  z-index: 2;
  color: #a0aec0;
  transition: color 0.3s ease;
}

.input-wrapper:focus-within .input-icon {
  color: #667eea;
}

input {
  width: 100%;
  height: 56px;
  padding: 0 16px 0 50px;
  border: 2px solid #e2e8f0;
  border-radius: 12px;
  font-size: 16px;
  background: #f7fafc;
  transition: all 0.3s ease;
  outline: none;
}

input:focus {
  border-color: #667eea;
  background: white;
  box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.1);
}

input::placeholder {
  color: #a0aec0;
}

/* 登录按钮 */
.login-btn {
  width: 100%;
  height: 56px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border: none;
  border-radius: 12px;
  color: white;
  font-size: 16px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  margin-top: 10px;
}

.login-btn:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 10px 30px rgba(102, 126, 234, 0.3);
}

.login-btn:active {
  transform: translateY(0);
}

.login-btn:disabled {
  opacity: 0.7;
  cursor: not-allowed;
  transform: none;
}

/* 加载动画 */
.loading-spinner {
  width: 20px;
  height: 20px;
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-top: 2px solid white;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

/* 提示信息 */
.login-tips {
  margin-top: 20px;
}

.tip-card {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  padding: 16px;
  background: #ebf8ff;
  border: 1px solid #bee3f8;
  border-radius: 12px;
}

.tip-icon {
  color: #3182ce;
  margin-top: 2px;
}

.tip-content {
  flex: 1;
}

.tip-title {
  font-size: 14px;
  font-weight: 600;
  color: #2b6cb0;
  margin-bottom: 4px;
}

.tip-text {
  font-size: 13px;
  color: #2c5282;
  line-height: 1.4;
}

.tip-text code {
  background: #2d3748;
  color: #e2e8f0;
  padding: 2px 6px;
  border-radius: 4px;
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  font-size: 12px;
}

/* 响应式设计 */
@media (max-width: 480px) {
  .login-form {
    width: 90vw;
    padding: 30px 20px;
  }

  .title {
    font-size: 24px;
  }

  .logo-icon {
    width: 60px;
    height: 60px;
  }
}
</style>
