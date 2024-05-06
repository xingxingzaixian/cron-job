export interface HttpResult<T> {
  code: number;
  data: T;
  message: string;
}
