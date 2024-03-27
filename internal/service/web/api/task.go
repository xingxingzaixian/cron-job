package api

import "github.com/gin-gonic/gin"

type TaskApi struct{}

func ServiceRegister(group *gin.RouterGroup) {
	service := &TaskApi{}
	group.GET("/list", service.GetTaskList)
}

func (s *TaskApi) GetTaskList(ctx *gin.Context) {

}
