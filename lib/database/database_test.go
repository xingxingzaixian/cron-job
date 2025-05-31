package database

import (
	"log"
	"os"
	"testing"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// TestDatabaseConnection 测试数据库连接
func TestDatabaseConnection(t *testing.T) {
	// 初始化配置
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("../")
	viper.AddConfigPath("../../")

	if err := viper.ReadInConfig(); err != nil {
		t.Skipf("跳过数据库连接测试，配置文件未找到: %v", err)
		return
	}

	// 初始化日志
	config := zap.NewDevelopmentConfig()
	zapLogger, _ := config.Build()
	zap.ReplaceGlobals(zapLogger)

	t.Log("=== 数据库连接测试 ===")

	// 显示当前配置
	engine := viper.GetString("db.engine")
	host := viper.GetString("db.host")
	port := viper.GetInt("db.port")
	name := viper.GetString("db.name")
	user := viper.GetString("db.user")

	t.Logf("数据库类型: %s", engine)
	t.Logf("主机地址: %s:%d", host, port)
	t.Logf("数据库名: %s", name)
	t.Logf("用户名: %s", user)

	// 创建数据库工厂
	factory := &DatabaseFactory{}

	// 验证配置
	t.Log("1. 验证数据库配置...")
	if err := factory.ValidateConfig(); err != nil {
		t.Fatalf("配置验证失败: %v", err)
	}
	t.Log("✓ 配置验证通过")

	// 创建数据库实例
	t.Log("2. 创建数据库实例...")
	db, err := factory.CreateDatabase()
	if err != nil {
		t.Fatalf("创建数据库实例失败: %v", err)
	}
	t.Log("✓ 数据库实例创建成功")

	// 配置GORM
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		})

	option := gorm.Config{
		Logger:                 newLogger,
		PrepareStmt:            true,
		SkipDefaultTransaction: true,
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   viper.GetString("db.prefix"),
			SingularTable: true,
		},
	}

	// 测试连接
	t.Log("3. 测试数据库连接...")
	gormDB, err := db.Create(&option)
	if err != nil {
		t.Fatalf("数据库连接失败: %v", err)
	}
	t.Log("✓ 数据库连接成功")

	// 获取底层数据库连接
	t.Log("4. 测试连接池...")
	sqlDB, err := gormDB.DB()
	if err != nil {
		t.Fatalf("获取数据库连接失败: %v", err)
	}

	// 测试ping
	if err := sqlDB.Ping(); err != nil {
		t.Fatalf("数据库ping失败: %v", err)
	}
	t.Log("✓ 数据库ping成功")

	// 显示连接池状态
	stats := sqlDB.Stats()
	t.Logf("连接池状态:")
	t.Logf("  - 打开连接数: %d", stats.OpenConnections)
	t.Logf("  - 使用中连接数: %d", stats.InUse)
	t.Logf("  - 空闲连接数: %d", stats.Idle)

	// 测试基本查询
	t.Log("5. 测试基本查询...")
	var result struct {
		Count int
	}

	switch engine {
	case "mysql":
		err = gormDB.Raw("SELECT 1 as count").Scan(&result).Error
	case "postgresql", "postgres":
		err = gormDB.Raw("SELECT 1 as count").Scan(&result).Error
	}

	if err != nil {
		t.Fatalf("基本查询失败: %v", err)
	}
	t.Log("✓ 基本查询成功")

	// 关闭连接
	t.Log("6. 关闭数据库连接...")
	if err := sqlDB.Close(); err != nil {
		t.Logf("关闭数据库连接失败: %v", err)
	} else {
		t.Log("✓ 数据库连接已关闭")
	}

	t.Log("=== 数据库连接测试完成 ===")
	t.Log("所有测试通过！数据库配置正确。")
}

// TestDatabaseFactory 测试数据库工厂
func TestDatabaseFactory(t *testing.T) {
	factory := &DatabaseFactory{}

	// 测试支持的数据库引擎
	engines := factory.GetSupportedEngines()
	expectedEngines := []string{"mysql", "postgresql", "postgres"}

	if len(engines) != len(expectedEngines) {
		t.Errorf("期望支持 %d 种数据库引擎，实际支持 %d 种", len(expectedEngines), len(engines))
	}

	for _, expected := range expectedEngines {
		found := false
		for _, engine := range engines {
			if engine == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("期望支持数据库引擎 %s，但未找到", expected)
		}
	}

	// 测试默认端口
	mysqlPort := factory.GetDefaultPort("mysql")
	if mysqlPort != 3306 {
		t.Errorf("MySQL默认端口期望为3306，实际为%d", mysqlPort)
	}

	postgresPort := factory.GetDefaultPort("postgresql")
	if postgresPort != 5432 {
		t.Errorf("PostgreSQL默认端口期望为5432，实际为%d", postgresPort)
	}

	unknownPort := factory.GetDefaultPort("unknown")
	if unknownPort != 0 {
		t.Errorf("未知数据库类型默认端口期望为0，实际为%d", unknownPort)
	}
}

// TestDatabaseConfigValidation 测试数据库配置验证
func TestDatabaseConfigValidation(t *testing.T) {
	factory := &DatabaseFactory{}

	// 保存原始配置
	originalEngine := viper.GetString("db.engine")
	originalHost := viper.GetString("db.host")
	originalPort := viper.GetInt("db.port")
	originalName := viper.GetString("db.name")
	originalUser := viper.GetString("db.user")

	// 恢复原始配置
	defer func() {
		viper.Set("db.engine", originalEngine)
		viper.Set("db.host", originalHost)
		viper.Set("db.port", originalPort)
		viper.Set("db.name", originalName)
		viper.Set("db.user", originalUser)
	}()

	// 测试空引擎
	viper.Set("db.engine", "")
	if err := factory.ValidateConfig(); err == nil {
		t.Error("期望空引擎配置验证失败，但验证通过了")
	}

	// 测试不支持的引擎
	viper.Set("db.engine", "oracle")
	if err := factory.ValidateConfig(); err == nil {
		t.Error("期望不支持的引擎配置验证失败，但验证通过了")
	}

	// 测试空主机
	viper.Set("db.engine", "mysql")
	viper.Set("db.host", "")
	if err := factory.ValidateConfig(); err == nil {
		t.Error("期望空主机配置验证失败，但验证通过了")
	}

	// 测试无效端口
	viper.Set("db.host", "localhost")
	viper.Set("db.port", 0)
	if err := factory.ValidateConfig(); err == nil {
		t.Error("期望无效端口配置验证失败，但验证通过了")
	}

	// 测试空数据库名
	viper.Set("db.port", 3306)
	viper.Set("db.name", "")
	if err := factory.ValidateConfig(); err == nil {
		t.Error("期望空数据库名配置验证失败，但验证通过了")
	}

	// 测试空用户名
	viper.Set("db.name", "test")
	viper.Set("db.user", "")
	if err := factory.ValidateConfig(); err == nil {
		t.Error("期望空用户名配置验证失败，但验证通过了")
	}

	// 测试有效配置
	viper.Set("db.user", "root")
	if err := factory.ValidateConfig(); err != nil {
		t.Errorf("期望有效配置验证通过，但验证失败: %v", err)
	}
}
