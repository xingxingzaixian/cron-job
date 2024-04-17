package models

import (
	"cronJob/internal/schemas"
	"errors"
	"github.com/gogf/gf/v2/frame/g"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UserName string `gorm:"size:64;not null;unique" json:"username" comment:"账号名"`
	NickName string `gorm:"size:64;not null" json:"nickname" comment:"用户名"`
	Password string `gorm:"size:128;not null" json:"password" comment:"密码"`
	Email    string `gorm:"size:64" json:"email"`
}

func (u *User) Create(tx *gorm.DB) (uint, error) {
	result := tx.Create(u)
	if result.Error != nil {
		return 0, result.Error
	}
	return u.ID, nil
}

func (u *User) Update(tx *gorm.DB, id uint, data g.Map) (int64, error) {
	result := tx.Model(u).Where("id = ?", id).Updates(data)
	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}

func (u *User) PageList(tx *gorm.DB, params *schemas.SearchUserParams) (users []User, count int64, err error) {
	query := tx.Model(u)
	if params.Name != "" {
		query = query.Where("nickname like ?", "%"+params.Name+"%")
	}

	query.Count(&count)
	offset := (params.PageNo - 1) * params.PageSize
	result := query.Limit(params.PageSize).Offset(offset).Order("id desc").Find(&users)
	if result.Error != nil {
		return
	}

	return
}

func (u *User) Find(tx *gorm.DB, userModel g.Map) (list []User, err error) {
	result := tx.Where(userModel).Find(&list)
	if result.RowsAffected == 0 {
		return nil, errors.New("用户不存在")
	}

	return
}

func (u *User) FindOne(tx *gorm.DB, userModel g.Map) error {
	result := tx.Where(userModel).First(u)
	if result.RowsAffected == 0 {
		return errors.New("用户不存在")
	}
	return nil
}
