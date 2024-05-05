package models

import (
	"errors"
	"github.com/gogf/gf/v2/frame/g"
	"gorm.io/gorm"
)

type Role struct {
	gorm.Model
	Name    string `gorm:"size:64;not null;unique;column:name" json:"name" comment:"角色名"`
	Key     string `gorm:"size:64;not null;unique;column:key" json:"key" comment:"角色Key"`
	Status  int    `gorm:"type:tinyint;not null;default:1" json:"status" comment:"状态（0:禁用，1:启用）"`
	IsSuper bool   `gorm:"type:tinyint;not null;default:0" json:"is_super" comment:"是否超级管理员"`
}

func (r *Role) Create(tx *gorm.DB) (uint, error) {
	result := tx.Create(r)
	if result.Error != nil {
		return 0, result.Error
	}
	return r.ID, nil
}

func (r *Role) Update(tx *gorm.DB, id uint, data g.Map) (int64, error) {
	result := tx.Model(r).Where("id = ?", id).Updates(data)
	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}

func (r *Role) Delete(tx *gorm.DB, id uint) (int64, error) {
	result := tx.Model(r).Where("id = ?", id).Delete(r)
	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}

func (r *Role) FindOne(tx *gorm.DB, roleModel g.Map) error {
	result := tx.Where(roleModel).First(r)
	if result.RowsAffected == 0 {
		return errors.New("角色不存在")
	}
	return nil
}

func (r *Role) PageList(tx *gorm.DB) (roles []*Role, count int64, err error) {
	query := tx.Model(r)

	query.Count(&count)
	result := query.Order("id desc").Find(&roles)
	if result.Error != nil {
		return
	}

	return
}
