package schemas

import "github.com/gin-gonic/gin"

// InstallCheckOutput 安装检查输出
type InstallCheckOutput struct {
	Installed    bool `json:"installed" comment:"是否已安装"`
	ConfigExists bool `json:"config_exists" comment:"配置文件是否存在"`
	DbConfigured bool `json:"db_configured" comment:"数据库是否已配置"`
}

// DatabaseTestInput 数据库测试输入
type DatabaseTestInput struct {
	Engine   string `json:"engine" form:"engine" comment:"数据库类型" example:"mysql" validate:"required"`
	Host     string `json:"host" form:"host" comment:"数据库主机" example:"127.0.0.1"`
	Port     int    `json:"port" form:"port" comment:"数据库端口" example:"3306"`
	Name     string `json:"name" form:"name" comment:"数据库名称" example:"cronJob" validate:"required"`
	User     string `json:"user" form:"user" comment:"数据库用户名" example:"root"`
	Password string `json:"password" form:"password" comment:"数据库密码" example:""`
	Path     string `json:"path" form:"path" comment:"SQLite数据库路径" example:""`
}

func (param *DatabaseTestInput) BindValidParam(c *gin.Context) error {
	return DefaultGetValidParams(c, param)
}

// DatabaseTestOutput 数据库测试输出
type DatabaseTestOutput struct {
	Success bool   `json:"success" comment:"测试是否成功"`
	Message string `json:"message" comment:"测试结果消息"`
}

// AdminUserInput 管理员用户输入
type AdminUserInput struct {
	Username string `json:"username" form:"username" comment:"用户名" example:"admin" validate:"required"`
	Nickname string `json:"nickname" form:"nickname" comment:"昵称" example:"系统管理员" validate:"required"`
	Password string `json:"password" form:"password" comment:"密码" example:"" validate:"required"`
	Email    string `json:"email" form:"email" comment:"邮箱" example:"admin@example.com"`
}

// InstallInput 安装输入
type InstallInput struct {
	Database    DatabaseTestInput `json:"database" comment:"数据库配置" validate:"required"`
	AuthEnabled bool              `json:"auth_enabled" comment:"是否启用认证" example:"true"`
	AdminUser   AdminUserInput    `json:"admin_user" comment:"管理员用户配置"`
}

func (param *InstallInput) BindValidParam(c *gin.Context) error {
	return DefaultGetValidParams(c, param)
}

// InstallOutput 安装输出
type InstallOutput struct {
	Success bool   `json:"success" comment:"安装是否成功"`
	Message string `json:"message" comment:"安装结果消息"`
}