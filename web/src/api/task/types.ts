import type { TaskPolicy, TaskProtocol, TaskStatus } from '@/enum/task';

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
  status: TaskStatus;
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
  op: string;
}

export interface QueryTaskLog {
  pageNo: number;
  pageSize: number;
  taskId: number;
}

export interface TaskLogItemOutput {
  id: number;
  result: string;
  retry_times: number;
  status: number;
  task_id: number;
  total_time: number;
}

export interface TaskLogOutput {
  list: TaskLogItemOutput[];
  total: number;
}
