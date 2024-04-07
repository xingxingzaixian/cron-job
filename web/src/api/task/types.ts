import type { TaskPolicy, TaskProtocol } from '@/enum/task';

export interface TaskEditHTTPInput {
  command: string;
  count: number;
  delay: number;
  id?: number | null;
  name: string;
  params: string;
  policy: TaskPolicy;
  protocol: TaskProtocol;
  remark: string;
  retry_interval: number;
  retry_times: number;
  spec: string;
  tag: string;
  timeout: number;
}

export interface TaskItemOutput {
  command: string;
  count: number;
  delay: number;
  id: number;
  name: string;
  params: string;
  policy: number;
  protocol: number;
  remark: string;
  retry_interval: number;
  retry_times: number;
  spec: string;
  status: number;
  tag: string;
  timeout: number;
}

export interface SearchTaskResponse {
  list: TaskItemOutput[];
  total: number;
}

export interface QueryTask {
  pageNo: number;
  pageSize: number;
  name?: string;
  protocol?: number;
  tag?: string;
}

export interface TaskOptionInput {
  id: number;
  op: 'start' | 'stop' | 'run' | 'delete';
}

export interface QueryTaskLog {
  pageNo: number;
  pageSize: number;
  taskId?: number;
  taskName?: string;
  status?: number;
}

export interface TaskLogItemOutput {
  id: number;
  taskId: number;
  taskName: string;
  result: string;
  protocol: number;
  retryTimes: number;
  status: number;
  startTime: string;
  endTime: string;
  totalTime: number;
}

export interface TaskLogOutput {
  list: TaskLogItemOutput[];
  total: number;
}
