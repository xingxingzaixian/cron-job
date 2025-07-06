package global

import (
	"net/http"
	"gorm.io/gorm"
)

const (
	ValidatorKey  = "ValidatorKey"
	TranslatorKey = "TranslatorKey"
)

var (
	GormDB         *gorm.DB
	HttpSrvHandler *http.Server
	// Install status flags
	IsInstalled bool = false
	InstallMode bool = false
)