package global

import "gorm.io/gorm"

const (
	ValidatorKey  = "ValidatorKey"
	TranslatorKey = "TranslatorKey"
)

var (
	GormDB *gorm.DB
)
