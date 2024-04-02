export interface HttpResult<T> {
  code: number;
  data: T;
  msg: string;
}
