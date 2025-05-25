package controllers

import (
	"net/http"

	"jethur/grc/library/services"
	"jethur/utility"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type KnowledgeCenterResourceController struct {
	service services.KnowledgeCenterResourceServiceInterface
	db      *gorm.DB
}

func NewKnowledgeCenterResourceController(service services.KnowledgeCenterResourceServiceInterface, db *gorm.DB) *KnowledgeCenterResourceController {
	return &KnowledgeCenterResourceController{
		service: service,
		db:      db,
	}
}

func (kcrc *KnowledgeCenterResourceController) GetAllResources(c *gin.Context) {
	resources, err := kcrc.service.GetAll()
	if err != nil {
		utility.SendInternalErrResp(kcrc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, resources)
}
