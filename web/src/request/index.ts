/* eslint-disable @typescript-eslint/no-explicit-any */
import Axios from 'axios';
import type { InternalAxiosRequestConfig, AxiosRequestConfig, AxiosInstance, AxiosResponse, AxiosError } from 'axios';
import { defaultConfig } from './types';
import useUserStore from '@/store/modules/user';

class AxiosHttp {
  private axiosInstance: AxiosInstance;

  constructor() {
    this.axiosInstance = Axios.create(defaultConfig);
    this.httpHookRequest();
    this.httpHookResponse();
  }

  // 请求拦截
  private httpHookRequest(): void {
    const userStore = useUserStore();
    this.axiosInstance.interceptors.request.use(
      (config: InternalAxiosRequestConfig) => {
        // 将 Token 添加到 header 中
        const token: string | undefined = userStore.userToken;
        if (token && config.headers) {
          config.headers['authorization'] = token;
        }
        return config;
      },
      (error) => {
        return Promise.reject(error);
      }
    );
  }

  // 响应拦截
  private httpHookResponse(): void {
    this.axiosInstance.interceptors.response.use(
      (response: AxiosResponse) => {
        return response.data;
      },
      (error: AxiosError) => {
        const { response } = error;
        if (response) {
          AxiosHttp.errorHandler(response.status, (response.data as any).message);
        }

        return Promise.reject(error);
      }
    );
  }

  // 异常请求处理
  static errorHandler(status: number, message?: string): void {
    switch (status) {
      case 401: {
        // 未登录
        break;
      }
      case 403: {
        // 没有权限
        break;
      }
      default: {
        console.log(message);
      }
    }
  }

  public get<T = any>(config: AxiosRequestConfig): Promise<T> {
    return this.request({ ...config, method: 'GET' });
  }

  public post<T = any>(config: AxiosRequestConfig): Promise<T> {
    return this.request({ ...config, method: 'POST' });
  }

  public put<T = any>(config: AxiosRequestConfig): Promise<T> {
    return this.request({ ...config, method: 'PUT' });
  }

  public patch<T = any>(config: AxiosRequestConfig): Promise<T> {
    return this.request({ ...config, method: 'PATCH' });
  }

  public delete<T = any>(config: AxiosRequestConfig): Promise<T> {
    return this.request({ ...config, method: 'DELETE' });
  }

  private request<T = any>(config: AxiosRequestConfig): Promise<T> {
    return new Promise<T>((resolve, reject) => {
      this.axiosInstance
        .request(config)
        .then((resp: AxiosResponse) => {
          resolve(resp as any);
        })
        .catch((err) => {
          reject(err);
        });
    });
  }
}

const http = new AxiosHttp();
export default http;
