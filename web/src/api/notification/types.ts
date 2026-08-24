import type { NotificationType, NotificationTrigger } from '@/enum/notification';

export interface NotificationConfigInput {
  id?: number | null;
  task_id: number;
  name: string;
  type: NotificationType;
  target: string;
  trigger: NotificationTrigger;
  enabled: boolean;

  email_recipients?: string;
  email_subject?: string;
  email_body?: string;

  webhook_url?: string;
  webhook_method?: string;
  webhook_headers?: string;
  webhook_body?: string;

  retry_times?: number;
  retry_interval?: number;
}

export interface NotificationConfigOutput {
  id: number;
  task_id: number;
  name: string;
  type: NotificationType;
  target: string;
  trigger: NotificationTrigger;
  enabled: boolean;
  created_at: string;
  updated_at: string;
}

export interface SearchNotificationParams {
  pageNo: number;
  pageSize: number;
  task_id?: number;
}

export interface SearchNotificationResponse {
  list: NotificationConfigOutput[];
  total: number;
}

export interface NotificationLogInput {
  pageNo: number;
  pageSize: number;
  task_id: number;
  notification_id?: number;
}

export interface NotificationLogOutput {
  id: number;
  notification_id: number;
  task_id: number;
  task_log_id: number;
  status: number;
  result: string;
  retry_count: number;
  start_time: string;
  end_time: string;
}

export interface NotificationLogResponse {
  list: NotificationLogOutput[];
  total: number;
}

export interface NotificationTestInput {
  type: NotificationType;
  target: string;
  content?: string;
}

export interface NotificationTestOutput {
  success: boolean;
  message: string;
}
