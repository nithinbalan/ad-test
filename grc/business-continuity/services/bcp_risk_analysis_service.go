package services

import (
	"jethur/grc/business-continuity/repositories"
)

type BcpRiskAnalysisService struct {
	repo repositories.BcpRiskAnalysisRepositoryInterface
}

func NewBcpRiskAnalysisService(repo repositories.BcpRiskAnalysisRepositoryInterface) *BcpRiskAnalysisService {
	return &BcpRiskAnalysisService{
		repo: repo,
	}
}
