import { ref, computed, watch } from 'vue';
import { defineStore } from 'pinia';
import store from '../../index';
import type { LoginInput, LoginOutput, UserEditInput } from '@/api/user/model';
import { fetchLoginApi, fetchUserInfo } from '@/api/user';
import type { HttpResult } from '@/types/api';

const userStore = defineStore('user-store', () => {
  const id = ref<number>(-1);
  const name = ref<string>('');
  const token = ref<string>('');
  const user = ref<UserEditInput | null>(null);

  const isLoggedIn = computed(() => {
    return token.value !== '';
  });

  const userToken = computed(() => {
    return token.value;
  });

  const userInfo = computed(() => {
    return user.value || {};
  });

  // 初始化从本地存储加载数据
  const initFromStorage = () => {
    const storedToken = localStorage.getItem('token');
    const storedUser = localStorage.getItem('user');

    if (storedToken) {
      token.value = storedToken;
    }

    if (storedUser) {
      try {
        const userData = JSON.parse(storedUser);
        user.value = userData;
        name.value = userData.nickname || '';
        id.value = userData.id || -1;
      } catch (error) {
        console.error('Failed to parse stored user data:', error);
      }
    }
  };

  const Login = async (form: LoginInput): Promise<HttpResult<LoginOutput | UserEditInput>> => {
    const res = await fetchLoginApi(form);
    if (res.code === 200) {
      user.value = res.data.user;
      token.value = res.data.token;
      name.value = res.data.user.nickname || '';
      id.value = res.data.user.id || -1;
    }
    return res;
  };

  const getUserInfo = () => {
    fetchUserInfo().then((res) => {
      if (res.code === 200) {
        user.value = res.data;
        name.value = res.data.nickname || '';
        id.value = res.data.id || -1;
      }
    });
  };

  const logout = () => {
    token.value = '';
    user.value = null;
    name.value = '';
    id.value = -1;
  };

  // 监听token变化，自动保存到本地存储
  watch(token, (newToken) => {
    if (newToken) {
      localStorage.setItem('token', newToken);
    } else {
      localStorage.removeItem('token');
    }
  });

  // 监听用户信息变化，自动保存到本地存储
  watch(user, (newUser) => {
    if (newUser) {
      localStorage.setItem('user', JSON.stringify(newUser));
    } else {
      localStorage.removeItem('user');
    }
  }, { deep: true });

  // 初始化
  initFromStorage();

  return {
    id,
    name,
    token,
    user,
    isLoggedIn,
    userToken,
    userInfo,
    Login,
    getUserInfo,
    logout,
    initFromStorage
  };
});

const useUserStore = () => {
  return userStore(store);
};

export default useUserStore;
