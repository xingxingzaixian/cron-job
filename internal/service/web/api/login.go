package api

import (
	"cronJob/internal/global"
	"cronJob/internal/models"
	"cronJob/internal/schemas"
	jwt2 "cronJob/lib/jwt"
	"github.com/gin-gonic/gin"
	"github.com/gogf/gf/v2/frame/g"
	"golang.org/x/crypto/bcrypt"
)

type LoginApi struct{}

func LoginRegister(group *gin.RouterGroup) {
	service := &LoginApi{}
	group.POST("/login", service.Login)
}

// Login godoc
// @Summary 登陆
// @Description 登陆
// @Tags 登陆
// @Security ApiKeyAuth
// @ID /api/login
// @Accept json
// @Produce json
// @Param data body schemas.LoginInput true "body"
// @Success 200 {object} schemas.Response{data=schemas.LoginOutput} "success"
// @Router /api/login [post]
func (service *LoginApi) Login(ctx *gin.Context) {
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

func (service *LoginApi) Logout(ctx *gin.Context) {

}
