package repositories

import (
	"gorm.io/gorm"
)

type BcpRiskAnalysisRepositoryInterface interface {
}

type BcpRiskAnalysisRepository struct {
	db *gorm.DB
}

func NewBcpRiskAnalysisRepository(db *gorm.DB) BcpRiskAnalysisRepositoryInterface {
	return &BcpRiskAnalysisRepository{
		db: db,
	}
}
