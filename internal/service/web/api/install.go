package api

import (
	"context"
	"cronJob/internal/global"
	"cronJob/internal/models"
	"cronJob/internal/schemas"
	taskManager "cronJob/internal/service/cron/task_manager"
	"cronJob/lib/config"
	"cronJob/lib/database"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type InstallApi struct{}

func InstallRegister(group *gin.RouterGroup) {
	service := &InstallApi{}
	group.GET("/check", service.CheckInstall)
	group.POST("/install", service.Install)
	group.POST("/test-db", service.TestDatabase)
}

// CheckInstall godoc
// @Summary 检查安装状态
// @Description 检查系统是否已安装
// @Tags 安装
// @ID /api/install/check
// @Accept json
// @Produce json
// @Success 200 {object} schemas.Response{data=schemas.InstallCheckOutput} "success"
// @Router /api/install/check [get]
func (s *InstallApi) CheckInstall(ctx *gin.Context) {
	// 检查是否已配置数据库
	dbConfigured := s.isDatabaseConfigured()

	// 检查配置文件是否存在
	configExists := s.isConfigFileExists()

	// 检查系统是否真正安装完成（数据库表已创建）
	installed := s.isSystemInstalled()

	schemas.ResponseSuccess(ctx, schemas.InstallCheckOutput{
		Installed:    installed,
		ConfigExists: configExists,
		DbConfigured: dbConfigured,
	})
}

// TestDatabase godoc
// @Summary 测试数据库连接
// @Description 测试数据库连接是否正常
// @Tags 安装
// @ID /api/install/test-db
// @Accept json
// @Produce json
// @Param data body schemas.DatabaseTestInput true "数据库配置"
// @Success 200 {object} schemas.Response{data=schemas.DatabaseTestOutput} "success"
// @Router /api/install/test-db [post]
func (s *InstallApi) TestDatabase(ctx *gin.Context) {
	params := &schemas.DatabaseTestInput{}
	if err := params.BindValidParam(ctx); err != nil {
		schemas.ResponseError(ctx, schemas.InstallParamInvalid, err)
		return
	}

	// 临时设置数据库配置用于测试
	s.setTempDatabaseConfig(params)

	// 创建数据库工厂并测试连接
	factory := &database.DatabaseFactory{}

	// 验证配置
	if err := factory.ValidateConfig(); err != nil {
		schemas.ResponseError(ctx, schemas.DatabaseTestFailed, err)
		return
	}

	// 尝试创建数据库连接
	db, err := factory.CreateDatabase()
	if err != nil {
		schemas.ResponseError(ctx, schemas.DatabaseTestFailed, err)
		return
	}

	// 测试连接
	gormDB, err := db.Create(&gorm.Config{})
	if err != nil {
		schemas.ResponseError(ctx, schemas.DatabaseTestFailed, err)
		return
	}

	// 测试ping
	sqlDB, err := gormDB.DB()
	if err != nil {
		schemas.ResponseError(ctx, schemas.DatabaseTestFailed, err)
		return
	}

	if err := sqlDB.Ping(); err != nil {
		schemas.ResponseError(ctx, schemas.DatabaseTestFailed, err)
		return
	}

	sqlDB.Close()

	schemas.ResponseSuccess(ctx, schemas.DatabaseTestOutput{
		Success: true,
		Message: "数据库连接测试成功",
	})
}

// Install godoc
// @Summary 执行安装
// @Description 执行系统安装，配置数据库和管理员账户
// @Tags 安装
// @ID /api/install/install
// @Accept json
// @Produce json
// @Param data body schemas.InstallInput true "安装配置"
// @Success 200 {object} schemas.Response{data=schemas.InstallOutput} "success"
// @Router /api/install/install [post]
func (s *InstallApi) Install(ctx *gin.Context) {
	params := &schemas.InstallInput{}
	if err := params.BindValidParam(ctx); err != nil {
		schemas.ResponseError(ctx, schemas.InstallParamInvalid, err)
		return
	}

	// 检查是否已安装（通过检查数据库是否已初始化）
	if s.isSystemInstalled() {
		schemas.ResponseError(ctx, schemas.InstallAlreadyInstalled, fmt.Errorf("系统已安装，请勿重复安装"))
		return
	}

	// 1. 写入配置文件
	if err := s.writeConfigFile(params); err != nil {
		schemas.ResponseError(ctx, schemas.InstallConfigWriteFailed, err)
		return
	}

	// 2. 重新加载配置
	config.InitConfig("config.yaml")

	// 3. 初始化数据库
	if err := s.initDatabase(); err != nil {
		schemas.ResponseError(ctx, schemas.InstallDatabaseInitFailed, err)
		return
	}

	// 4. 创建管理员账户（如果启用认证）
	if params.AuthEnabled {
		if err := s.createAdminUser(params.AdminUser); err != nil {
			schemas.ResponseError(ctx, schemas.InstallCreateAdminFailed, err)
			return
		}
	}

	// 5. 安装完成，准备重启服务
	schemas.ResponseSuccess(ctx, schemas.InstallOutput{
		Success: true,
		Message: "系统安装成功，服务器将在3秒后重启以应用新配置",
	})

	// 异步执行重启操作，给前端时间接收响应
	go func() {
		time.Sleep(3 * time.Second)
		s.restartApplication()
	}()
}

// 检查数据库是否已配置
func (s *InstallApi) isDatabaseConfigured() bool {
	engine := viper.GetString("db.engine")
	zap.S().Infof(" [INFO] isDatabaseConfigured engine:%s\n", engine)
	if engine == "" {
		return false
	}

	switch engine {
	case "mysql", "postgresql", "postgres":
		host := viper.GetString("db.host")
		name := viper.GetString("db.name")
		user := viper.GetString("db.user")
		return host != "" && name != "" && user != ""
	case "sqlite", "sqlite3":
		name := viper.GetString("db.name")
		path := viper.GetString("db.path")
		return name != "" || path != ""
	}

	return false
}

// 检查系统是否已安装（通过检查数据库表是否存在）
func (s *InstallApi) isSystemInstalled() bool {
	// 如果数据库配置不存在，肯定没有安装
	if !s.isDatabaseConfigured() {
		return false
	}

	// 尝试连接数据库并检查表是否存在
	factory := &database.DatabaseFactory{}
	if err := factory.ValidateConfig(); err != nil {
		return false
	}

	db, err := factory.CreateDatabase()
	if err != nil {
		return false
	}

	gormDB, err := db.Create(&gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   viper.GetString("db.prefix"),
			SingularTable: true,
		},
	})
	if err != nil {
		return false
	}

	defer func() {
		if sqlDB, err := gormDB.DB(); err == nil {
			sqlDB.Close()
		}
	}()

	// 检查用户表是否存在
	return gormDB.Migrator().HasTable(&models.User{})
}

// 检查配置文件是否存在
func (s *InstallApi) isConfigFileExists() bool {
	_, err := os.Stat("config.yaml")
	return !os.IsNotExist(err)
}

// 设置临时数据库配置用于测试
func (s *InstallApi) setTempDatabaseConfig(params *schemas.DatabaseTestInput) {
	viper.Set("db.engine", params.Engine)

	switch params.Engine {
	case "mysql", "postgresql", "postgres":
		viper.Set("db.host", params.Host)
		viper.Set("db.port", params.Port)
		viper.Set("db.name", params.Name)
		viper.Set("db.user", params.User)
		viper.Set("db.password", params.Password)
	case "sqlite", "sqlite3":
		if params.Path != "" {
			viper.Set("db.path", params.Path)
		} else {
			viper.Set("db.name", params.Name)
			viper.Set("db.data_dir", "data")
		}
	}
}

// 写入配置文件
func (s *InstallApi) writeConfigFile(params *schemas.InstallInput) error {
	// 读取示例配置文件
	configData := make(map[string]interface{})

	// 基础配置
	configData["debug"] = "release"
	configData["http"] = map[string]interface{}{
		"addr":             ":8210",
		"read_timeout":     10,
		"write_timeout":    10,
		"max_header_bytes": 20,
	}

	configData["swagger"] = map[string]interface{}{
		"title":     "定时任务服务swagger API",
		"desc":      "这是一个简单的定时任务执行系统",
		"host":      "127.0.0.1:8210",
		"base_path": "",
	}

	configData["cors"] = map[string]interface{}{
		"allowed_origins": "http://localhost:3000,http://127.0.0.1:3000,http://localhost:8210,http://127.0.0.1:8210",
	}

	// 数据库配置
	dbConfig := map[string]interface{}{
		"engine": params.Database.Engine,
		"prefix": "sched_",
	}

	switch params.Database.Engine {
	case "mysql", "postgresql", "postgres":
		dbConfig["host"] = params.Database.Host
		dbConfig["port"] = params.Database.Port
		dbConfig["name"] = params.Database.Name
		dbConfig["user"] = params.Database.User
		dbConfig["password"] = params.Database.Password
		dbConfig["max_idle_conns"] = 10
		dbConfig["max_open_conns"] = 100
		dbConfig["conn_max_lifetime"] = 3600
	case "sqlite", "sqlite3":
		if params.Database.Path != "" {
			dbConfig["path"] = params.Database.Path
		} else {
			dbConfig["name"] = params.Database.Name
			dbConfig["data_dir"] = "data"
		}
		dbConfig["sqlite"] = map[string]interface{}{
			"synchronous":  "NORMAL",
			"cache_size":   -64000,
			"temp_store":   "MEMORY",
			"busy_timeout": 30000,
		}
	}

	configData["db"] = dbConfig

	// 认证配置
	configData["auth"] = map[string]interface{}{
		"enable": params.AuthEnabled,
	}

	// JWT配置
	configData["jwt"] = map[string]interface{}{
		"secret":  "your_jwt_secret_here_" + fmt.Sprintf("%d", g.NewVar(nil).Int64()),
		"expires": 7200,
	}

	// 写入配置文件
	yamlData, err := yaml.Marshal(configData)
	if err != nil {
		return fmt.Errorf("序列化配置失败: %v", err)
	}

	// 确保目录存在
	configDir := filepath.Dir("config.yaml")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("创建配置目录失败: %v", err)
	}

	// 写入文件
	if err := os.WriteFile("config.yaml", yamlData, 0644); err != nil {
		return fmt.Errorf("写入配置文件失败: %v", err)
	}

	return nil
}

// 初始化数据库
func (s *InstallApi) initDatabase() error {
	defer func() {
		if r := recover(); r != nil {
			zap.S().Errorf("数据库初始化失败: %v", r)
		}
	}()

	// 初始化数据库连接
	database.InitDB(viper.GetString("db.prefix"))

	return nil
}

// 创建管理员用户
func (s *InstallApi) createAdminUser(adminUser schemas.AdminUserInput) error {
	user := &models.User{
		UserName: adminUser.Username,
		NickName: adminUser.Nickname,
		Password: adminUser.Password,
		Email:    adminUser.Email,
	}

	tx := global.GormDB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	_, err := user.Create(tx)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("创建管理员用户失败: %v", err)
	}

	tx.Commit()
	return nil
}

// 重启应用程序
func (s *InstallApi) restartApplication() {
	zap.S().Info("开始重启应用程序...")

	// 尝试进程重启
	if s.tryProcessRestart() {
		return
	}

	// 如果进程重启失败，则进行服务重新初始化
	zap.S().Warn("进程重启失败，尝试重新初始化服务...")
	s.reinitializeServices()
}

// 尝试进程重启
func (s *InstallApi) tryProcessRestart() bool {
	defer func() {
		if r := recover(); r != nil {
			zap.S().Errorf("进程重启失败: %v", r)
		}
	}()

	// 获取当前执行文件路径
	executable, err := os.Executable()
	if err != nil {
		zap.S().Errorf("获取执行文件路径失败: %v", err)
		return false
	}

	// 获取当前进程的参数
	args := os.Args[1:] // 排除程序名

	zap.S().Infof("准备重启进程: %s %v", executable, args)

	// 根据操作系统选择重启方式
	switch runtime.GOOS {
	case "windows":
		return s.restartOnWindows(executable, args)
	case "linux", "darwin":
		return s.restartOnUnix(executable, args)
	default:
		zap.S().Errorf("不支持的操作系统: %s", runtime.GOOS)
		return false
	}
}

// Windows 系统重启
func (s *InstallApi) restartOnWindows(executable string, args []string) bool {
	zap.S().Info("Windows 重启：先停止HTTP服务器避免端口冲突")

	// 延迟执行重启操作，避免端口冲突
	go func() {
		// 等待当前响应发送完成
		time.Sleep(2 * time.Second)

		// 先停止HTTP服务器释放端口
		zap.S().Info("停止HTTP服务器...")
		if global.HttpSrvHandler != nil {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if err := global.HttpSrvHandler.Shutdown(ctx); err != nil {
				zap.S().Errorf("停止HTTP服务器失败: %v", err)
			} else {
				zap.S().Info("HTTP服务器已停止")
			}
		}

		// 再启动新进程
		zap.S().Info("启动新进程...")
		cmd := exec.Command(executable, args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		err := cmd.Start()
		if err != nil {
			zap.S().Errorf("启动新进程失败: %v", err)
			return
		}

		zap.S().Info("新进程已启动，退出当前进程")
		os.Exit(0)
	}()

	return true
}

// Unix 系统重启 (Linux/macOS)
func (s *InstallApi) restartOnUnix(executable string, args []string) bool {
	// 使用 execve 系统调用替换当前进程
	env := os.Environ()

	// 准备参数数组
	argv := make([]string, len(args)+1)
	argv[0] = executable
	copy(argv[1:], args)

	zap.S().Info("使用 execve 重启进程")

	// 延迟执行 execve，给日志时间输出
	go func() {
		time.Sleep(1 * time.Second)
		err := syscall.Exec(executable, argv, env)
		if err != nil {
			zap.S().Errorf("execve 重启失败: %v", err)
		}
	}()

	return true
}

// 重新初始化服务（如果进程重启失败的备用方案）
func (s *InstallApi) reinitializeServices() {
	zap.S().Info("开始重新初始化服务...")

	defer func() {
		if r := recover(); r != nil {
			zap.S().Errorf("服务重新初始化失败: %v", r)
		}
	}()

	// 1. 重新加载配置
	zap.S().Info("重新加载配置...")
	config.InitConfig("config.yaml")

	// 2. 重新初始化数据库
	zap.S().Info("重新初始化数据库...")
	database.InitDB(viper.GetString("db.prefix"))

	// 3. 启动定时任务调度器
	zap.S().Info("启动定时任务调度器...")
	go func() {
		taskManager.CronServerRun()
	}()

	// 4. 设置全局标志，表示系统已完成安装
	global.IsInstalled = true
	global.InstallMode = false

	zap.S().Info("服务重新初始化完成，系统已切换到正常运行模式")
}
