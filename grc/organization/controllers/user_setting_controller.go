package controllers

import (
	"jethur/grc/organization/services"

	"gorm.io/gorm"
)

type UserSettingController struct {
	service services.UserServiceInterface
	db      *gorm.DB
}

func NewUSerSettingController(service services.UserServiceInterface) *UserController {
	return &UserController{
		service: service,
	}
}
