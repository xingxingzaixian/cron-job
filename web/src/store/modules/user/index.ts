import { ref, computed } from 'vue';
import { defineStore } from 'pinia';
import store from '../../index';

const userStore = defineStore('user-store', () => {
  const id = ref<number>(-1);
  const name = ref<string>('');
  const token = ref<string>('');

  const isLoggedIn = computed(() => {
    return token.value !== '';
  });

  const userToken = computed(() => {
    return token.value;
  });

  const userInfo = computed(() => {
    return {};
  });
  return {
    id,
    name,
    token,
    isLoggedIn,
    userToken,
    userInfo
  };
});

const useUserStore = () => {
  return userStore(store);
};

export default useUserStore;
