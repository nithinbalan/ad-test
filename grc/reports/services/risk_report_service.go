package services

import (
	reportDTOs "jethur/grc/reports/dtos"
	"jethur/grc/reports/repositories"
)

type RiskReportServiceInterface interface {
	GetKRIByCategory(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.KRIByCategoryDTO, error)
	GetKRIByResponsibleUnit(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.KRIByResponsibleUnitDTO, error)
	GetKRIByResponsibleUser(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.KRIByResponsibleUserDTO, error)
	GetKRIByRiskLevel(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.KRIByRiskLevelDTO, error)
	GetKRIByReviewFrequency(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.KRIByReviewFrequencyDTO, error)
}

type RiskReportService struct {
	repo repositories.RiskReportRepositoryInterface
}

func NewRiskReportService(repo repositories.RiskReportRepositoryInterface) RiskReportServiceInterface {
	return &RiskReportService{repo: repo}
}

func (r *RiskReportService) GetKRIByCategory(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.KRIByCategoryDTO, error) {
	return r.repo.GetKRIByCategory(params)
}

func (r *RiskReportService) GetKRIByResponsibleUnit(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.KRIByResponsibleUnitDTO, error) {
	return r.repo.GetKRIByResponsibleUnit(params)
}

func (r *RiskReportService) GetKRIByResponsibleUser(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.KRIByResponsibleUserDTO, error) {
	return r.repo.GetKRIByResponsibleUser(params)
}

func (r *RiskReportService) GetKRIByRiskLevel(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.KRIByRiskLevelDTO, error) {
	return r.repo.GetKRIByRiskLevel(params)
}

func (r *RiskReportService) GetKRIByReviewFrequency(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.KRIByReviewFrequencyDTO, error) {
	return r.repo.GetKRIByReviewFrequency(params)
}
