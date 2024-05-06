import { ref, computed } from 'vue';
import { defineStore } from 'pinia';
import store from '../../index';
import type { LoginInput, LoginOutput, UserEditInput } from '@/api/user/model';
import { fetchLoginApi, fetchUserInfo } from '@/api/user';
import { HttpResult } from '@/types/api';

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
    return {};
  });

  const Login = async (form: LoginInput): Promise<HttpResult<LoginOutput | UserEditInput>> => {
    const res = await fetchLoginApi(form);
    if (res.code === 200) {
      // TODO: 获取用户信息
      user.value = res.data.user;
      token.value = res.data.token;
    }
    return res;
  };

  const getUserInfo = async (): Promise<void> => {
    const res = await fetchUserInfo();
    if (res.code === 200) {
      user.value = res.data;
    }
  };

  return {
    id,
    name,
    token,
    isLoggedIn,
    userToken,
    userInfo,
    Login,
    getUserInfo
  };
});

const useUserStore = () => {
  return userStore(store);
};

export default useUserStore;
