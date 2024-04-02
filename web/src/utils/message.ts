import type { ConfigProviderProps } from 'naive-ui';
import { createDiscreteApi } from 'naive-ui';
import { computed } from 'vue';

import useThemeStore from '@/store/modules/theme';

const configProviderPropsRef = computed<ConfigProviderProps>(() => {
  const themeStore = useThemeStore();
  return {
    theme: themeStore.naiveTheme
  };
});

const { message, dialog, notification, loadingBar } = createDiscreteApi(
  ['message', 'dialog', 'notification', 'loadingBar'],
  {
    configProviderProps: configProviderPropsRef
  }
);
export { dialog, loadingBar, message, notification };
