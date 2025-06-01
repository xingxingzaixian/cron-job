package schemas

import "github.com/gin-gonic/gin"

type LoginInput struct {
	UserName string `form:"userName" json:"userName" binding:"required" comment:"用户名必须填写"`
	Password string `form:"password" json:"password" binding:"required" comment:"密码必须填写"`
}

func (s *LoginInput) BindValidParam(ctx *gin.Context) error {
	return DefaultGetValidParams(ctx, s)
}

type LoginOutput struct {
	User  *UserEditInput `json:"user"`
	Token string         `json:"token"`
}

type CheckAuthOutput struct {
	AuthEnabled bool `json:"authEnabled" comment:"是否启用认证"`
}
