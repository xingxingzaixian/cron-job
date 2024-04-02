import { ref } from 'vue';
import { defineStore } from 'pinia';
import store from '../../index';
import { getStorageItem, setStorageItem } from '@/utils/storage';
import { setDayjsLocale, setLocale } from '@/locales';
import { useBoolean } from '@/hooks';
import { useTitle } from '@vueuse/core';
import useRouteStore from '../route';
import { $t } from '@/locales';
import router from '@/router';

const appStore = defineStore('app-store', () => {
  const routeStore = useRouteStore();

  const locale = ref<I18n.LangType>((getStorageItem('lang') as I18n.LangType) || 'zh-CN');
  const localeOptions: I18n.LangOption[] = [
    {
      label: '中文',
      key: 'zh-CN'
    },
    {
      label: 'English',
      key: 'en-US'
    }
  ];

  const { bool: reloadFlag, setBool: setReloadFlag } = useBoolean(true);
  const { bool: contentXScrollable, setBool: setContentXScrollable } = useBoolean();
  const { bool: siderCollapse, setBool: setSiderCollapse } = useBoolean();

  function changeLocale(lang: I18n.LangType) {
    locale.value = lang;
    setLocale(lang);
    setDayjsLocale(lang);
    setStorageItem('lang', lang);
    updateDocumentTitleByLocale();
    routeStore.updateGlobalMenusByLocale();
    reloadPage(1000);
  }

  async function reloadPage(duration = 0) {
    setReloadFlag(false);

    if (duration > 0) {
      await new Promise((resolve) => {
        setTimeout(resolve, duration);
      });
    }

    setReloadFlag(true);
  }

  function updateDocumentTitleByLocale() {
    const { title } = router.currentRoute.value.meta;

    const documentTitle = $t(title as string);

    useTitle(documentTitle);
  }

  return {
    locale,
    localeOptions,
    reloadFlag,
    reloadPage,
    setReloadFlag,
    changeLocale,
    siderCollapse,
    setSiderCollapse,
    contentXScrollable,
    setContentXScrollable
  };
});

const useAppStore = () => {
  return appStore(store);
};

export default useAppStore;
