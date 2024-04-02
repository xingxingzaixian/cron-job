import { locale } from 'dayjs';
import 'dayjs/locale/zh-cn';
import 'dayjs/locale/en';
import { getStorageItem } from '@/utils/storage';

export const setDayjsLocale = (lang: I18n.LangType = 'zh-CN') => {
  const localMap = {
    'zh-CN': 'zh-cn',
    'en-US': 'en'
  };

  const l = lang || getStorageItem('lang') || 'zh-CN';
  locale(localMap[l]);
};
