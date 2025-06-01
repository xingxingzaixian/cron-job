import type { HttpResult } from '@/types/api';
import type {
  LoginInput,
  LoginOutput,
  UserEditInput,
  UserList,
  SearchUserParams,
  DelUserParams,
  UpdatePasswordInput,
  CheckAuthOutput
} from './model';
import apiHttp from '@/request';

// 登录
export const fetchLoginApi = async (data: LoginInput): Promise<HttpResult<LoginOutput>> => {
  return apiHttp.post<HttpResult<LoginOutput>>({
    url: '/api/login',
    data
  });
};

// 获取当前用户信息
export const fetchUserInfo = async (): Promise<HttpResult<UserEditInput>> => {
  return apiHttp.get<HttpResult<UserEditInput>>({
    url: '/api/user/info'
  });
};

// 获取用户列表
export const fetchUserList = async (params: SearchUserParams): Promise<HttpResult<UserList>> => {
  return apiHttp.get<HttpResult<UserList>>({
    url: '/api/user/list',
    params
  });
};

// 获取用户详情
export const fetchUserDetail = async (id: number): Promise<HttpResult<UserEditInput>> => {
  return apiHttp.get<HttpResult<UserEditInput>>({
    url: '/api/user/view',
    params: { id }
  });
};

// 新增/编辑用户
export const fetchUserEdit = async (data: UserEditInput): Promise<HttpResult<string>> => {
  return apiHttp.post<HttpResult<string>>({
    url: '/api/user/edit',
    data
  });
};

// 删除用户
export const fetchUserDelete = async (data: DelUserParams): Promise<HttpResult<string>> => {
  return apiHttp.post<HttpResult<string>>({
    url: '/api/user/del',
    data
  });
};

// 更新密码
export const fetchUpdatePassword = async (data: UpdatePasswordInput): Promise<HttpResult<string>> => {
  return apiHttp.post<HttpResult<string>>({
    url: '/api/user/update-password',
    data
  });
};

// 获取是否启用认证
export const fetchCheckAuth = async (): Promise<HttpResult<CheckAuthOutput>> => {
  return apiHttp.get<HttpResult<CheckAuthOutput>>({
    url: '/api/check-auth'
  });
};
