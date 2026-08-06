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
  TaskItemOutput,
  TaskTestInput,
  TaskTestOutput
} from './types';

export const fetchTaskList = (params: QueryTask) => {
  return apiHttp.get<HttpResult<SearchTaskResponse>>({
    url: '/api/task/list',
    params
  });
};

export const fetchTaskCreate = (data: TaskEditHTTPInput) => {
  return apiHttp.post<HttpResult<string>>({
    url: '/api/task/create',
    data
  });
};

export const fetchTaskUpdate = (data: TaskEditHTTPInput) => {
  return apiHttp.post<HttpResult<string>>({
    url: '/api/task/update',
    data
  });
};

export const fetchTaskOp = (data: TaskOptionInput) => {
  const opUrls: Record<TaskOptionInput['op'], string> = {
    start: '/api/task/start',
    stop: '/api/task/stop',
    run: '/api/task/execute',
    delete: '/api/task/delete'
  };

  return apiHttp.post<HttpResult<boolean>>({
    url: opUrls[data.op],
    data
  });
};

export const fetchTaskTest = (data: TaskTestInput) => {
  return apiHttp.post<HttpResult<TaskTestOutput>>({
    url: '/api/task/test',
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

export const fetchTaskLogDelete = (ids: number[]) => {
  return apiHttp.delete<HttpResult<boolean>>({
    url: '/api/taskLog/delete',
    data: {
      ids
    }
  });
};

export const fetchTaskLogClean = () => {
  return apiHttp.post<HttpResult<number>>({
    url: '/api/taskLog/clean'
  });
};
