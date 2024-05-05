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
	group.GET("/init", service.Init)
}

// Login godoc
// @Summary 登陆
// @Description 登陆
// @Tags 登陆
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
			Role:     user.Role.ID,
		},
	})
}

func (service *LoginApi) Init(ctx *gin.Context) {
	tx := global.GormDB.Begin()
	role1 := &models.Role{}
	role1.Name = "管理员"
	role1.Key = "admin"
	role1.Status = 1
	role1.IsSuper = true
	if _, err := role1.Create(tx); err != nil {
		tx.Rollback()
		schemas.ResponseError(ctx, schemas.RoleCreateFailed, err)
		return
	}

	role2 := &models.Role{}
	role2.Name = "普通成员"
	role2.Key = "default"
	role2.Status = 1
	role2.IsSuper = false
	if _, err := role2.Create(tx); err != nil {
		tx.Rollback()
		schemas.ResponseError(ctx, schemas.RoleCreateFailed, err)
		return
	}

	user := &models.User{}
	user.NickName = "admin"
	user.UserName = "admin"
	user.Password = "admin"
	user.Role = *role1
	user.Email = "admin@qq.com"
	if _, err := user.Create(tx); err != nil {
		tx.Rollback()
	}
	tx.Commit()
	schemas.ResponseSuccess(ctx, nil)
}
