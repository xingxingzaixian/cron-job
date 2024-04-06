export enum TaskProtocol {
  HTTP = 1,
  Shell = 2,
  Grpc = 3
}

export enum TaskPolicy {
  Multi = 1, // 并行策略
  Once = 2, // 单词策略
  Single = 3, // 单利策略
  Times = 4 // 多次策略
}

export enum TaskStatus {
  Disabled = 0,
  Enabled = 1,
  Failure = 2,
  Running = 3,
  Finish = 4
}
