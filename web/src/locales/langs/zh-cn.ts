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
    updateFailed: '更新失败',
    log: '日志',
    search: '搜索',
    reset: '重置',
    save: '保存',
    close: '关闭'
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
        runTask: '运行',
        startTask: '启动',
        startDesc: '启动任务运行',
        stopTask: '停止',
        runDesc: '运行一次任务',
        params: '请求参数',
        count: '执行次数',
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
        taskId: '任务ID',
        taskName: '任务名称',
        protocol: '执行方式',
        retryTimes: '重试次数',
        runTime: '运行时长(秒)',
        startTime: '开始时间',
        runStatus: '运行状态',
        runResult: '运行结果',

        failure: '失败',
        success: '成功',
        running: '运行中',
        timeout: '超时',

        title: '任务执行结果',
        command: '任务命令',
        result: '执行结果'
      }
    }
  },
  task: {
    protocol: {
      all: '全部',
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
      once: '单次策略',
      single: '单例策略',
      times: '多次策略'
    },
    runPolicy: '执行策略',
    message: {
      runSuccess: '任务运行成功',
      runFailed: '任务运行失败',
      startSuccess: '任务启动成功',
      startFailed: '任务启动失败',
      stopSuccess: '任务停止成功',
      stopFailed: '任务停止失败',
      addSuccess: '任务新增成功',
      addFailed: '任务新增失败',
      editSuccess: '任务编辑成功',
      editFailed: '任务编辑失败',
      deleteSuccess: '任务删除成功',
      deleteFailed: '任务删除失败',
      confirmDelete: '确认删除该任务吗？'
    }
  },
  http: {
    header: '请求头',
    query: '参数',
    data: '请求体'
  }
};

export default local;
