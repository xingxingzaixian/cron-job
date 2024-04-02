import { ContentTypeEnum } from '@/enum';
import { AxiosRequestConfig } from 'axios';

/**
 * 默认配置
 */
export const defaultConfig: AxiosRequestConfig = {
  baseURL: import.meta.env.VITE_APP_BASE_URL ? import.meta.env.VITE_APP_BASE_URL : '',
  //10秒超时
  timeout: 100000,
  headers: { 'Content-Type': ContentTypeEnum.JSON }
};
