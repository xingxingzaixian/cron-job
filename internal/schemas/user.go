package schemas

import "github.com/gin-gonic/gin"

type UserEditInput struct {
	ID       uint   `json:"id" form:"id" comment:"用户ID" example:"1"`
	UserName string `json:"username" form:"username" comment:"账号名" example:""`
	NickName string `json:"nickname" form:"nickname" comment:"用户名" example:""`
	Password string `json:"password" form:"password" comment:"密码" example:""`
	Email    string `json:"email" form:"email" comment:"邮箱" example:""`
	Role     uint   `json:"role" form:"role" comment:"角色ID"`
}

func (param *UserEditInput) BindValidParam(c *gin.Context) error {
	return DefaultGetValidParams(c, param)
}

type UserList struct {
	Total int64            `json:"total" comment:"总数"`
	List  []*UserEditInput `json:"list" comment:"列表"`
}

type SearchUserParams struct {
	FormPage
	Name string `json:"name" form:"name" comment:"用户名" example:""`
	ID   uint   `json:"id" form:"id" comment:"用户ID" example:""`
}

func (param *SearchUserParams) BindValidParam(c *gin.Context) error {
	return DefaultGetValidParams(c, param)
}

type DelUserParams struct {
	ID uint `json:"id" form:"id" comment:"用户ID"`
}

func (param *DelUserParams) BindValidParam(c *gin.Context) error {
	return DefaultGetValidParams(c, param)
}

type UpdatePasswordInput struct {
	Password        string `json:"password" form:"password" comment:"密码"`
	NewPass         string `json:"newPass" form:"newPass" comment:"新密码"`
	ConfirmPassword string `json:"confirmPassword" form:"confirmPassword" comment:"确认密码"`
	UserName        string `json:"username" form:"username" comment:"账号名"`
}

func (param *UpdatePasswordInput) BindValidParam(c *gin.Context) error {
	return DefaultGetValidParams(c, param)
}
