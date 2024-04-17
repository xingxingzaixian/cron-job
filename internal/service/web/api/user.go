package api

import (
	"cronJob/internal/global"
	"cronJob/internal/models"
	"cronJob/internal/schemas"
	"github.com/gin-gonic/gin"
	"github.com/gogf/gf/v2/frame/g"
	"golang.org/x/crypto/bcrypt"
)

type UserApi struct{}

func UserRegister(group *gin.RouterGroup) {
	service := &UserApi{}
	group.GET("/list", service.GetList)
	group.POST("/edit", service.EditUser)
	group.GET("/view", service.GetUser)
}

// GetList godoc
// @Summary 查询用户列表
// @Description 查询用户列表
// @Tags 用户管理
// @Security ApiKeyAuth
// @ID /api/user/list
// @Accept json
// @Produce json
// @Param pageNo query int true "页数"
// @Param pageSize query int true "每页条数"
// @Param name query string false "用户昵称"
// @Success 200 {object} schemas.Response{data=schemas.UserList} "success"
// @Router /api/user/list [get]
func (u *UserApi) GetList(c *gin.Context) {
	params := &schemas.SearchUserParams{}
	if err := params.BindValidParam(c); err != nil {
		schemas.ResponseError(c, schemas.UserSearchParamInvalid, err)
		return
	}

	user := &models.User{}
	users, count, err := user.PageList(global.GormDB, params)
	if err != nil {
		schemas.ResponseError(c, schemas.UserSearchListInvalid, err)
		return
	}

	result := schemas.UserList{}
	result.Total = count
	for _, item := range users {
		result.List = append(result.List, &schemas.UserEditInput{
			ID:       item.ID,
			UserName: item.UserName,
			NickName: item.NickName,
			Email:    item.Email,
		})
	}

	schemas.ResponseSuccess(c, result)
}

// EditUser godoc
// @Summary 新增/修改用户
// @Description 新增/修改用户
// @Tags 用户管理
// @Security ApiKeyAuth
// @ID /api/user/edit
// @Accept json
// @Produce json
// @Param data body schemas.UserEditInput true "body"
// @Success 200 {object} schemas.Response{data=string} "success"
// @Router /api/user/edit [post]
func (u *UserApi) EditUser(c *gin.Context) {
	params := &schemas.UserEditInput{}
	if err := params.BindValidParam(c); err != nil {
		schemas.ResponseError(c, schemas.UserEditParamInvalid, err)
		return
	}

	user := &models.User{}
	tx := global.GormDB.Begin()
	if params.ID > 0 {
		err := user.FindOne(tx, g.Map{
			"id": params.ID,
		})
		if err != nil {
			tx.Rollback()
			schemas.ResponseError(c, schemas.UserNotExist, err)
			return
		}
		_, err = user.Update(tx, params.ID, g.Map{
			"nickname": params.NickName,
		})
		if err != nil {
			tx.Rollback()
			schemas.ResponseError(c, schemas.UserUpdateFailed, err)
		}
	} else {
		user.UserName = params.UserName
		user.NickName = params.NickName
		user.Email = params.Email
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
		if err != nil {
			tx.Rollback()
			schemas.ResponseError(c, schemas.UserCreateFailed, err)
			return
		}
		user.Password = string(hashedPassword)
		//user.RoleId = 2
		_, err = user.Create(tx)
		if err != nil {
			tx.Rollback()
			schemas.ResponseError(c, schemas.UserCreateFailed, err)
		}
	}
	tx.Commit()
	schemas.ResponseSuccess(c, nil)
}

// GetUser godoc
// @Summary 用户信息
// @Description 用户信息
// @Tags 用户管理
// @Security ApiKeyAuth
// @ID /api/user/view
// @Accept json
// @Produce json
// @Param id query uint true "用户ID"
// @Success 200 {object} schemas.Response{data=schemas.UserEditInput} "success"
// @Router /api/user/view [get]
func (u *UserApi) GetUser(c *gin.Context) {
	params := &schemas.SearchUserParams{}
	if err := params.BindValidParam(c); err != nil {
		schemas.ResponseError(c, schemas.UserSearchParamInvalid, err)
		return
	}

	user := &models.User{}
	if err := user.FindOne(global.GormDB, g.Map{
		"id": params.ID,
	}); err != nil {
		schemas.ResponseError(c, schemas.UserNotExist, err)
		return
	}

	schemas.ResponseSuccess(c, user)
}
