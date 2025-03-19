package controllers

import (
	"jethur/grc/business-continuity/services"

	"gorm.io/gorm"
)

type BcpRiskAnalysisController struct {
	service services.BcpRiskAnalysisService
	db      *gorm.DB
}

func NewBcpRiskAnalysisController(service services.BcpRiskAnalysisService) *BcpRiskAnalysisController {
	return &BcpRiskAnalysisController{
		service: service,
	}
}
