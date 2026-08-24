export enum NotificationType {
  Email = 'email',
  Webhook = 'webhook',
}

export enum NotificationTrigger {
  Success = 'success',
  Failure = 'failure',
  All = 'all',
  Timeout = 'timeout',
  Cancel = 'cancel',
}

export enum NotificationStatus {
  Disabled = 0,
  Enabled = 1,
}
