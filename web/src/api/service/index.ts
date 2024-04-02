import apiHttp from '@/request';
import type {
  ServiceListOutput,
  ServiceSearchParam,
  ServiceAddHTTPInput,
  ServiceUpdateHTTPInput,
  ServiceDetail,
  ServiceRequestLogList,
  ServiceStatisticsRes
} from './types';
import type { HttpResult } from '@/types/api';

export const fetchServiceList = (params: ServiceSearchParam) => {
  return apiHttp.get<HttpResult<ServiceListOutput>>({
    url: '/api/service/list',
    params
  });
};

export const fetchServiceDel = (id: number) => {
  return apiHttp.delete({
    url: `/api/service/${id}`
  });
};

export const fetchServiceAdd = (data: ServiceAddHTTPInput) => {
  return apiHttp.post({
    url: `/api/service/service_add_http`,
    data
  });
};

export const fetchServiceUpdate = (data: ServiceUpdateHTTPInput) => {
  return apiHttp.post({
    url: `/api/service/service_update_http`,
    data
  });
};

export const fetchServiceGet = (id: number) => {
  return apiHttp.get<HttpResult<ServiceDetail>>({
    url: `/api/service/${id}`
  });
};

export const fetchServiceLogList = (params: ServiceSearchParam) => {
  return apiHttp.get<HttpResult<ServiceRequestLogList>>({
    url: '/api/service/log',
    params
  });
};

export const fetchServiceStatistics = (classify: 'day' | 'week' | 'month' | 'interface') => {
  return apiHttp.get<HttpResult<ServiceStatisticsRes[]>>({
    url: '/api/service/statistics',
    params: {
      classify
    }
  });
};
