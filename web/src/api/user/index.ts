import { HttpResult } from '@/types/api';
import type { LoginInput, LoginOutput, UserEditInput, RoleList } from './model';
import apiHttp from '@/request';

// 获取页面数据
export const fetchLoginApi = async (data: LoginInput): Promise<HttpResult<LoginOutput>> => {
  return apiHttp.post<HttpResult<LoginOutput>>({
    url: '/api/user/login',
    data
  });
};

export const fetchUserInfo = async (): Promise<HttpResult<UserEditInput>> => {
  return apiHttp.get<HttpResult<UserEditInput>>({
    url: '/api/user/info'
  });
};

export const fetchRoleList = async (): Promise<HttpResult<RoleList>> => {
  return apiHttp.get<HttpResult<RoleList>>({
    url: '/api/role/list'
  });
};
