package api

import (
	"cronJob/internal/global"
	"cronJob/internal/models"
	"cronJob/internal/schemas"
	"github.com/gin-gonic/gin"
)

type RoleApi struct{}

func RoleRegister(group *gin.RouterGroup) {
	service := &RoleApi{}
	group.GET("/list", service.GetList)
}

// GetList godoc
// @Summary 查询角色列表
// @Description 查询角色列表
// @Tags 用户管理
// @Security ApiKeyAuth
// @ID /api/role/list
// @Accept json
// @Produce json
// @Success 200 {object} schemas.Response{data=schemas.RoleList} "success"
// @Router /api/role/list [get]
func (s *RoleApi) GetList(c *gin.Context) {
	role := &models.Role{}
	roles, count, err := role.PageList(global.GormDB)
	if err != nil {
		schemas.ResponseError(c, schemas.RoleListFailed, err)
		return
	}

	result := schemas.RoleList{}
	result.Total = count
	for _, item := range roles {
		result.List = append(result.List, &schemas.RoleItem{
			ID:      item.ID,
			Name:    item.Name,
			Key:     item.Key,
			Status:  item.Status,
			IsSuper: item.IsSuper,
		})
	}
	schemas.ResponseSuccess(c, result)
}
