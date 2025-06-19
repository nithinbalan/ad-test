package controllers

import (
	"net/http"

	libraryDTOs "jethur/grc/library/dtos"
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
	keyword := c.Query("keyword")

	queryParams := libraryDTOs.KnowledgeCenterResourceListQueryParamsDTO{
		Keyword: keyword,
	}

	resources, err := kcrc.service.GetAll(&queryParams)
	if err != nil {
		utility.SendInternalErrResp(kcrc.db, c, err.Error())
		return
	}
	c.JSON(http.StatusOK, resources)
}
