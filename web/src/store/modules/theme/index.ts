import { ref, computed } from 'vue';
import { defineStore } from 'pinia';
import store from '../../index';
import { darkTheme, lightTheme } from 'naive-ui';

const themeStore = defineStore('theme-store', () => {
  const darkMode = ref<boolean>(false);
  const animateMode = ref<string>('fade-slide');
  const collapsedWidth = ref<number>(64);
  const siderWidth = ref<number>(240);
  const naiveTheme = computed(() => {
    return darkMode.value ? darkTheme : lightTheme;
  });

  const themeOverrides = computed(() => {
    return darkMode.value
      ? {
          common: {
            primaryColor: '#165DFF'
          }
        }
      : {
          common: {
            primaryColor: '#722ED1'
          }
        };
  });

  function toggleThemeScheme() {
    darkMode.value = !darkMode.value;
  }

  return {
    darkMode,
    animateMode,
    naiveTheme,
    themeOverrides,
    toggleThemeScheme,
    collapsedWidth,
    siderWidth
  };
});

const useThemeStore = () => {
  return themeStore(store);
};

export default useThemeStore;
