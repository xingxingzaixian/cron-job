export interface LoginInput {
  userName: string;
  password: string;
}

export interface UserEditInput {
  id?: number;
  username: string;
  nickname: string;
  password: string;
  email?: string;
  role: number;
}

export interface LoginOutput {
  token: string;
  user: UserEditInput;
}

export interface RoleItem {
  id: number;
  is_super: number;
  key: string;
  name: string;
  status: string;
}

export interface RoleList {
  total: number;
  list: RoleItem[];
}
