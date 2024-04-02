import { LoginData } from './types';
import apiHttp from '@/request';

// 获取页面数据
export const fetchLoginApi = async (data: LoginData): Promise<string> => {
  return apiHttp.post<string>({
    url: '/api/user/login',
    data
  });
};
