package schemas

import "github.com/gin-gonic/gin"

type LoginInput struct {
	UserName string `form:"userName" json:"userName" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

func (s *LoginInput) BindValidParam(ctx *gin.Context) error {
	return DefaultGetValidParams(ctx, s)
}

type LoginOutput struct {
	User  *UserEditInput `json:"user"`
	Token string         `json:"token"`
}
