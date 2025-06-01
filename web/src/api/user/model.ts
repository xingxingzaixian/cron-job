export interface LoginInput {
  userName: string;
  password: string;
}

export interface UserEditInput {
  id?: number;
  username: string;
  nickname: string;
  password?: string;
  email?: string;
  created_at?: string;
  updated_at?: string;
}

export interface LoginOutput {
  token: string;
  user: UserEditInput;
}

export interface UserList {
  total: number;
  list: UserEditInput[];
}

export interface SearchUserParams {
  pageNo: number;
  pageSize: number;
  name?: string;
  username?: string;
  email?: string;
  id?: number;
}

export interface DelUserParams {
  id: number;
}

export interface UpdatePasswordInput {
  password: string;
  newPass: string;
  confirmPassword: string;
  username: string;
}

export interface CheckAuthOutput {
  authEnabled: boolean;
}
