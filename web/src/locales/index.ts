import { createI18n } from 'vue-i18n';
import messages from './locale';
import { getStorageItem } from '@/utils/storage';
import { setDayjsLocale } from './dayjs';
import { naiveLocales, naiveDateLocales } from './naive';

const i18n = createI18n({
  locale: getStorageItem('lang') || 'zh-CN',
  fallbackLocale: 'en',
  messages,
  legacy: false
});

const $t = i18n.global.t;

const setLocale = (locale: I18n.LangType) => {
  i18n.global.locale.value = locale;
};

export { i18n, setLocale, setDayjsLocale, naiveLocales, naiveDateLocales, $t };
