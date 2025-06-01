package api

import (
	"cronJob/internal/global"
	"cronJob/internal/models"
	"cronJob/internal/schemas"
	jwt2 "cronJob/lib/jwt"

	"github.com/gin-gonic/gin"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

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

	// 判断用户是否存在
	user := &models.User{}
	if err := user.FindOne(global.GormDB, g.Map{
		"username": params.UserName,
	}); err != nil {
		schemas.ResponseError(ctx, schemas.UserNotExist, err)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(params.Password)); err != nil {
		schemas.ResponseError(ctx, schemas.UserPasswordError, err)
		return
	}

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
