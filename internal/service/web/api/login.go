package api

import (
	"cronJob/internal/global"
	"cronJob/internal/models"
	"cronJob/internal/schemas"
	jwt2 "cronJob/lib/jwt"
	"errors"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

// 登录失败限流：同一IP在时间窗口内最多允许失败次数
const (
	maxLoginFailures    = 5
	loginFailureWindow  = 10 * time.Minute
	maxLoginAttemptKeys = 10000
)

type loginAttempt struct {
	failures  int
	windowEnd time.Time
}

var (
	loginAttempts   = make(map[string]*loginAttempt)
	loginAttemptsMu sync.Mutex
)

// checkLoginRateLimit 检查是否允许继续尝试登录
func checkLoginRateLimit(key string) bool {
	loginAttemptsMu.Lock()
	defer loginAttemptsMu.Unlock()

	now := time.Now()
	entry, ok := loginAttempts[key]
	if !ok || now.After(entry.windowEnd) {
		entry = &loginAttempt{windowEnd: now.Add(loginFailureWindow)}
		loginAttempts[key] = entry
	}
	return entry.failures < maxLoginFailures
}

// recordLoginFailure 记录一次登录失败
func recordLoginFailure(key string) {
	loginAttemptsMu.Lock()
	defer loginAttemptsMu.Unlock()

	now := time.Now()
	entry, ok := loginAttempts[key]
	if !ok || now.After(entry.windowEnd) {
		entry = &loginAttempt{windowEnd: now.Add(loginFailureWindow)}
		loginAttempts[key] = entry
	}
	entry.failures++

	// 防止map无限增长：条目过多时清理过期记录
	if len(loginAttempts) > maxLoginAttemptKeys {
		for k, e := range loginAttempts {
			if now.After(e.windowEnd) {
				delete(loginAttempts, k)
			}
		}
	}
}

// resetLoginFailures 登录成功后清零失败记录
func resetLoginFailures(key string) {
	loginAttemptsMu.Lock()
	defer loginAttemptsMu.Unlock()
	delete(loginAttempts, key)
}

type LoginApi struct{}

func LoginRegister(group *gin.RouterGroup) {
	service := &LoginApi{}
	group.POST("/login", service.Login)
	group.GET("/check-auth", service.CheckAuth)
}

// Login godoc
// @Summary 登录
// @Description 登录
// @Tags 登录
// @ID /api/login
// @Accept json
// @Produce json
// @Param data body schemas.LoginInput true "body"
// @Success 200 {object} schemas.Response{data=schemas.LoginOutput} "success"
// @Router /api/login [post]
func (service *LoginApi) Login(ctx *gin.Context) {
	// 检查是否启用认证，如果未启用则返回默认登录信息
	if !viper.GetBool("auth.enable") {
		schemas.ResponseSuccess(ctx, schemas.LoginOutput{
			Token: "no-auth-token",
			User: &schemas.UserEditInput{
				UserName: "admin",
				NickName: "系统管理员",
				Email:    "admin@example.com",
				ID:       1,
			},
		})
		return
	}

	params := &schemas.LoginInput{}
	if err := params.BindValidParam(ctx); err != nil {
		schemas.ResponseError(ctx, schemas.LoginParamInvalid, err)
		return
	}

	// 登录失败限流，防止暴力破解
	rateKey := ctx.ClientIP()
	if !checkLoginRateLimit(rateKey) {
		schemas.ResponseError(ctx, schemas.LoginTooManyAttempts, errors.New("登录尝试过于频繁，请稍后再试"))
		return
	}

	// 判断用户是否存在
	user := &models.User{}
	if err := user.FindOne(global.GormDB, g.Map{
		"username": params.UserName,
	}); err != nil {
		recordLoginFailure(rateKey)
		// 统一错误提示，避免用户名枚举
		schemas.ResponseError(ctx, schemas.UserPasswordError, errors.New("用户名或密码错误"))
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(params.Password)); err != nil {
		recordLoginFailure(rateKey)
		schemas.ResponseError(ctx, schemas.UserPasswordError, errors.New("用户名或密码错误"))
		return
	}

	resetLoginFailures(rateKey)

	token, err := jwt2.GenToken(user.UserName)
	if err != nil {
		schemas.ResponseError(ctx, schemas.UserPasswordError, err)
		return
	}

	schemas.ResponseSuccess(ctx, schemas.LoginOutput{
		Token: token,
		User: &schemas.UserEditInput{
			UserName: user.UserName,
			NickName: user.NickName,
			Email:    user.Email,
			ID:       user.ID,
		},
	})
}

// CheckAuth godoc
// @Summary 获取是否启用认证
// @Description 获取是否启用认证
// @Tags 登录
// @ID /api/check-auth
// @Accept json
// @Produce json
// @Success 200 {object} schemas.Response{data=schemas.CheckAuthOutput} "success"
// @Router /api/check-auth [get]
func (service *LoginApi) CheckAuth(ctx *gin.Context) {
	schemas.ResponseSuccess(ctx, schemas.CheckAuthOutput{
		AuthEnabled: viper.GetBool("auth.enable"),
	})
}
