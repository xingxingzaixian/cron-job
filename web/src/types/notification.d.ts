declare namespace Notification {
  interface ConfigInput {
    id?: number | null;
    task_id: number;
    name: string;
    type: 'email' | 'webhook';
    target: string;
    trigger: 'success' | 'failure' | 'all' | 'timeout' | 'cancel';
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

  interface ConfigOutput {
    id: number;
    task_id: number;
    name: string;
    type: 'email' | 'webhook';
    target: string;
    trigger: 'success' | 'failure' | 'all' | 'timeout' | 'cancel';
    enabled: boolean;
    created_at: string;
    updated_at: string;
  }

  interface LogOutput {
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
}
