package schemas

type RoleItem struct {
	ID      uint   `json:"id" form:"id" comment:"角色ID" example:"1"`
	Name    string `json:"name" form:"name" comment:"角色名称" example:""`
	Key     string `json:"key" form:"key" comment:"角色Key" example:""`
	Status  int    `json:"status" form:"status" comment:"状态"`
	IsSuper bool   `json:"is_super" form:"is_super" comment:"是否超级管理员"`
}

type RoleList struct {
	Total int64       `json:"total" comment:"总数"`
	List  []*RoleItem `json:"list" comment:"列表"`
}
