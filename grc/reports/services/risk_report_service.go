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

	GetRiskByCategory(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.RiskByCategoryDTO, error)
	GetRiskByResponsibleUnit(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.RiskByResponsibleUnitDTO, error)
	GetRiskByResponsibleUser(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.RiskByResponsibleUserDTO, error)
	GetRiskByRiskSource(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.RiskByRiskSourceDTO, error)
	GetRiskByInherentRiskRating(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.RiskByRiskRatingDTO, error)
	GetRiskByResidualRiskRating(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.RiskByRiskRatingDTO, error)
	GetRiskByRiskStatus(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.RiskByStatusDTO, error)
	GetRiskByRiskResponseStrategy(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.RiskByResponseStrategyDTO, error)
	GetRiskByRiskAppetiteToleranceLevel(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.RiskByRiskAppetiteToleranceLevelDTO, error)
	GetRiskByKRIEnabledStatus(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.RiskByKRIEnabledStatusDTO, error)
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

func (r *RiskReportService) GetRiskByCategory(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.RiskByCategoryDTO, error) {
	return r.repo.GetRiskByCategory(params)
}

func (r *RiskReportService) GetRiskByResponsibleUnit(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.RiskByResponsibleUnitDTO, error) {
	return r.repo.GetRiskByResponsibleUnit(params)
}

func (r *RiskReportService) GetRiskByResponsibleUser(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.RiskByResponsibleUserDTO, error) {
	return r.repo.GetRiskByResponsibleUser(params)
}

func (r *RiskReportService) GetRiskByRiskSource(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.RiskByRiskSourceDTO, error) {
	return r.repo.GetRiskByRiskSource(params)
}

func (r *RiskReportService) GetRiskByInherentRiskRating(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.RiskByRiskRatingDTO, error) {
	return r.repo.GetRiskByInherentRiskRating(params)
}

func (r *RiskReportService) GetRiskByResidualRiskRating(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.RiskByRiskRatingDTO, error) {
	return r.repo.GetRiskByResidualRiskRating(params)
}

func (r *RiskReportService) GetRiskByRiskStatus(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.RiskByStatusDTO, error) {
	return r.repo.GetRiskByRiskStatus(params)
}

func (r *RiskReportService) GetRiskByRiskResponseStrategy(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.RiskByResponseStrategyDTO, error) {
	return r.repo.GetRiskByRiskResponseStrategy(params)
}

func (r *RiskReportService) GetRiskByRiskAppetiteToleranceLevel(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.RiskByRiskAppetiteToleranceLevelDTO, error) {
	return r.repo.GetRiskByRiskAppetiteToleranceLevel(params)
}

func (r *RiskReportService) GetRiskByKRIEnabledStatus(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.RiskByKRIEnabledStatusDTO, error) {
	return r.repo.GetRiskByKRIEnabledStatus(params)
}
