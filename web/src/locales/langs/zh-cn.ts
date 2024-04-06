const local = {
  common: {
    tip: '提示',
    add: '新增',
    edit: '编辑',
    batchDelete: '批量删除',
    deleteSuccess: '删除成功',
    confirmDelete: '确认删除吗？',
    logoutConfirm: '确认退出登录吗？',
    confirm: '确认',
    delete: '删除',
    cancel: '取消',
    refresh: '刷新',
    check: '勾选',
    columnSetting: '列设置',
    backToHome: '返回首页',
    operate: '操作',
    update: '更新',
    updateSuccess: '更新成功',
    updateFailed: '更新失败'
  },
  route: {
    login: '登陆',
    home: '首页',
    welcome: '欢迎访问 CronJob 任务管理系统',
    exception: '异常',
    exception_403: '没有权限',
    exception_404: '请求错误',
    exception_500: '服务器错误',
    task_log: '日志',
    task_list: '列表',
    task: '任务列表'
  },
  icon: {
    fullscreenExit: '退出全屏',
    fullscreen: '全屏'
  },
  system: {
    title: 'CronJob 任务管理系统'
  },
  page: {
    task: {
      list: {
        title: '任务列表',
        name: '任务名称',
        tag: '任务标签',
        spec: 'cron表达式',
        protocol: '执行方式',
        status: '状态',
        policy: '策略',
        command: '命令',
        timeout: '超时时间',
        retry_times: '重试次数',
        retry_interval: '重试间隔时间',
        delay: '延时启动(秒)',
        remark: '备注',
        params: '请求参数',
        addTask: '新增任务',
        editTask: '编辑任务',
        form: {
          name: '任务名称',
          tag: '任务标签',
          spec: 'cron表达式',
          specDesc: '秒 分 时 日 月 周',
          protocol: '执行方式',
          status: '状态',
          policy: '策略',
          timeout: '超时时间',
          command: '执行命令',
          retry_times: '重试次数',
          retry_interval: '重试间隔时间',
          delay: '延时时间(秒)',
          remark: '备注',
          request: '请求地址'
        }
      },
      log: {
        originUrl: '源地址',
        proxyUrl: '代理地址',
        serviceName: '服务名称'
      }
    }
  },
  task: {
    protocol: {
      http: 'HTTP',
      shell: 'Shell',
      grpc: 'Grpc'
    },
    status: {
      disabled: '禁用',
      enabled: '启用',
      failure: '失败',
      running: '运行中',
      finish: '完成'
    },
    policy: {
      multi: '并行策略',
      once: '单词策略',
      single: '单利策略',
      times: '多次策略'
    },
    runPolicy: '执行策略'
  },
  http: {
    header: '请求头',
    query: '参数',
    data: '请求体'
  }
};

export default local;
