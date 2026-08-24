import apiHttp from '@/request';
import type { HttpResult } from '@/types/api';
import type {
  SearchNotificationParams,
  SearchNotificationResponse,
  NotificationConfigInput,
  NotificationConfigOutput,
  NotificationLogInput,
  NotificationLogResponse,
  NotificationTestInput,
  NotificationTestOutput,
} from './types';

export const fetchNotificationList = (params: SearchNotificationParams) => {
  return apiHttp.get<HttpResult<SearchNotificationResponse>>({
    url: '/api/notification/list',
    params,
  });
};

export const fetchNotificationView = (id: number) => {
  return apiHttp.get<HttpResult<NotificationConfigOutput>>({
    url: `/api/notification/view?id=${id}`,
  });
};

export const fetchNotificationCreate = (data: NotificationConfigInput) => {
  return apiHttp.post<HttpResult<string>>({
    url: '/api/notification/create',
    data,
  });
};

export const fetchNotificationUpdate = (data: NotificationConfigInput) => {
  return apiHttp.post<HttpResult<string>>({
    url: '/api/notification/update',
    data,
  });
};

export const fetchNotificationDelete = (id: number) => {
  return apiHttp.post<HttpResult<boolean>>({
    url: `/api/notification/delete?id=${id}`,
  });
};

export const fetchNotificationTest = (data: NotificationTestInput) => {
  return apiHttp.post<HttpResult<NotificationTestOutput>>({
    url: '/api/notification/test',
    data,
  });
};

export const fetchNotificationLogs = (params: NotificationLogInput) => {
  return apiHttp.get<HttpResult<NotificationLogResponse>>({
    url: '/api/notification/logs',
    params,
  });
};
