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
    welcome: '欢迎访问 GateKeeper 管理系统',
    exception: '异常',
    exception_403: '没有权限',
    exception_404: '请求错误',
    exception_500: '服务器错误',
    service: '服务列表',
    service_list: '列表',
    service_log: '日志',
    task: '任务列表'
  },
  icon: {
    fullscreenExit: '退出全屏',
    fullscreen: '全屏'
  },
  system: {
    title: 'GateKeeper 管理系统'
  },
  page: {
    home: {
      dayTitle: '今日访问量',
      weekTitle: '本周访问量',
      monthTitle: '本月访问量',
      interfaceTitle: '接口调用量',
      visits: '访问量',
      hour: '小时',
      day: '天',
      interface: '接口'
    },
    service: {
      list: {
        title: '服务列表',
        serviceName: '服务名称',
        rule: '代理规则',
        serviceDesc: '服务描述',
        urlRewrite: '代理地址',
        needHttps: '是否支持https',
        addService: '新增服务',
        editService: '编辑服务',
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
