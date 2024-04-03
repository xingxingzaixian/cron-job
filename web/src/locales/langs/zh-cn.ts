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
        addTask: '新增任务',
        editTask: '编辑任务',
        form: {
          serviceName: '请输入服务名称',
          serviceDesc: '请输入服务描述',
          rule: '请输入代理规则',
          urlRewrite: '请输入代理地址'
        }
      },
      log: {
        originUrl: '源地址',
        proxyUrl: '代理地址',
        serviceName: '服务名称'
      }
    }
  }
};

export default local;
