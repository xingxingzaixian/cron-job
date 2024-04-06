import apiHttp from '@/request';
import type { HttpResult } from '@/types/api';
import type {
  QueryTask,
  SearchTaskResponse,
  TaskEditHTTPInput,
  TaskOptionInput,
  TaskLogOutput,
  TaskLogItemOutput,
  QueryTaskLog,
  TaskItemOutput
} from './types';

export const fetchTaskList = (params: QueryTask) => {
  return apiHttp.get<HttpResult<SearchTaskResponse>>({
    url: '/api/task/list',
    params
  });
};

export const fetchTaskEdit = (data: TaskEditHTTPInput) => {
  return apiHttp.get<HttpResult<string>>({
    url: '/api/task/edit',
    data
  });
};

export const fetchTaskOp = (data: TaskOptionInput) => {
  return apiHttp.post<HttpResult<boolean>>({
    url: '/api/task/op',
    data
  });
};

export const fetchTaskView = (id: number) => {
  return apiHttp.get<HttpResult<TaskItemOutput>>({
    url: `/api/task/view?id=${id}`
  });
};

export const fetchTaskLogList = (params: QueryTaskLog) => {
  return apiHttp.get<HttpResult<TaskLogOutput>>({
    url: '/api/taskLog/list',
    params
  });
};

export const fetchTaskLogQuery = (id: number) => {
  return apiHttp.get<HttpResult<TaskLogItemOutput>>({
    url: `/api/taskLog/query?id=${id}`
  });
};
