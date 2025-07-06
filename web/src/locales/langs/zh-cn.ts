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
    close: '关闭',
    enable: '启用',
    disable: '禁用',
    addSuccess: '新增成功',
    addFailed: '新增失败'
  },
  route: {
    login: '登陆',
    install: '系统安装',
    home: '首页',
    welcome: '欢迎访问 CronJob 任务管理系统',
    exception: '异常',
    exception_403: '没有权限',
    exception_404: '请求错误',
    exception_500: '服务器错误',
    task_log: '日志',
    task_list: '列表',
    task: '任务列表',
    user: '用户管理',
    user_list: '用户列表'
  },
  icon: {
    fullscreenExit: '退出全屏',
    fullscreen: '全屏'
  },
  system: {
    title: 'CronJob 任务管理系统'
  },
  page: {
    login: {
      accout: '账号',
      pass: '密码',
      submit: '登录',
      failed: '用户或密码错误'
    },
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
          timeoutUnit: '单位：秒',
          command: '执行命令',
          config: '协议配置',
          retry_times: '重试次数',
          retry_interval: '重试间隔时间',
          delay: '延时时间(秒)',
          remark: '备注',
          request: '请求地址',
          shellPlaceholder: '请输入Shell命令...',
          policyMulti: '并行策略',
          policyOnce: '单次策略',
          policySingle: '单例策略',
          policyTimes: '多次策略',
          count: '执行次数'
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
    },
    user: {
      list: {
        title: '用户管理',
        username: '用户名',
        nickname: '昵称',
        email: '邮箱',
        role: '角色',
        status: '状态',
        createTime: '创建时间',
        updateTime: '更新时间',
        addUser: '新增用户',
        editUser: '编辑用户',
        changePassword: '修改密码',
        searchPlaceholder: '搜索用户昵称',
        form: {
          username: '用户名',
          nickname: '昵称',
          email: '邮箱',
          password: '密码',
          role: '角色',
          oldPassword: '原密码',
          newPassword: '新密码',
          confirmPassword: '确认密码'
        }
      },
      message: {
        addSuccess: '用户新增成功',
        addFailed: '用户新增失败',
        editSuccess: '用户编辑成功',
        editFailed: '用户编辑失败',
        deleteSuccess: '用户删除成功',
        deleteFailed: '用户删除失败',
        passwordChangeSuccess: '密码修改成功',
        passwordChangeFailed: '密码修改失败',
        confirmDelete: '确认删除该用户吗？'
      }
    }
  },
  task: {
    protocol: {
      all: '全部',
      http: 'HTTP',
      shell: 'Shell',
      ssh: 'SSH'
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
  },
  ssh: {
    connection: 'SSH连接',
    connectionConfig: 'SSH连接配置',
    host: '主机地址',
    hostPlaceholder: '请输入服务器IP或域名',
    port: '端口',
    username: '用户名',
    usernamePlaceholder: '请输入SSH用户名',
    authType: '认证方式',
    passwordAuth: '密码认证',
    password: '密码',
    passwordPlaceholder: '请输入SSH密码',
    mode: '执行模式',
    sequentialMode: '顺序执行',
    scriptMode: '脚本执行',
    commands: '执行命令',
    commandsConfig: '命令配置',
    commandsPlaceholder: '请输入要执行的命令，每行一个命令\n支持注释（以#开头）\n\n示例：\nwhoami\npwd\nls -la\n# 这是注释\ndf -h',
    commandsHelp: '提示：每行一个命令，支持注释（#开头），空行会被忽略',
    testConnection: '测试连接',
    showExample: '加载示例',
    configPreview: '配置预览',
    fillRequiredFields: '请填写必填字段',
    fillAuthInfo: '请填写认证信息（密码或私钥）',
    connectionSuccess: 'SSH连接测试成功',
    connectionFailed: 'SSH连接测试失败',
    exampleLoaded: '示例配置已加载'
  }
};

export default local;
