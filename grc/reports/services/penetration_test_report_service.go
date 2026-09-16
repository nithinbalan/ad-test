package services

import (
	reportDTOs "jethur/grc/reports/dtos"
	"jethur/grc/reports/repositories"
)

type PenetrationTestReportServiceInterface interface {
	GetPenetrationTestFindingBySeverityLevel(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.PenetrationTestFindingBySeverityLevelDTO, error)
	GetPenetrationTestFindingByResponsibleUser(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.PenetrationTestFindingByResponsibleUserDTO, error)
	GetPenetrationTestFindingByResponsibleUnit(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.PenetrationTestFindingByResponsibleUnitDTO, error)
	GetPenetrationTestFindingByStatus(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.PenetrationTestFindingByStatusDTO, error)
	GetPenetrationTestByType(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.PenetrationTestByTypeDTO, error)
	GetPenetrationTestByStatus(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.PenetrationTestByStatusDTO, error)
	GetPenetrationTestByClassification(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.PenetrationTestByClassificationDTO, error)
	GetPenetrationTestByPlan(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.PenetrationTestByPlanDTO, error)
	GetPenetrationTestByExecutor(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.PenetrationTestByExecutorDTO, error)
	GetPenetrationTestByResponsibleUser(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.PenetrationTestByResponsibleUserDTO, error)
	GetPenetrationTestByCreatedUser(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.PenetrationTestByCreatedUserDTO, error)
	GetPenetrationTestByResponsibleUnit(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.PenetrationTestByResponsibleUnitDTO, error)
	GetPenetrationTestByYear(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.PenetrationTestByYearDTO, error)
	GetPenetrationTestByPeriod(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.PenetrationTestByPeriodDTO, error)
}

type PenetrationTestReportService struct {
	repo repositories.PenetrationTestReportRepositoryInterface
}

func NewPenetrationTestReportService(repo repositories.PenetrationTestReportRepositoryInterface) PenetrationTestReportServiceInterface {
	return &PenetrationTestReportService{repo: repo}
}

func (ptrs *PenetrationTestReportService) GetPenetrationTestByType(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.PenetrationTestByTypeDTO, error) {
	return ptrs.repo.GetPenetrationTestByType(params)
}

func (ptrs *PenetrationTestReportService) GetPenetrationTestByStatus(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.PenetrationTestByStatusDTO, error) {
	return ptrs.repo.GetPenetrationTestByStatus(params)
}

func (ptrs *PenetrationTestReportService) GetPenetrationTestByClassification(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.PenetrationTestByClassificationDTO, error) {
	return ptrs.repo.GetPenetrationTestByClassification(params)
}

func (ptrs *PenetrationTestReportService) GetPenetrationTestByPlan(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.PenetrationTestByPlanDTO, error) {
	return ptrs.repo.GetPenetrationTestByPlan(params)
}

func (ptrs *PenetrationTestReportService) GetPenetrationTestByExecutor(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.PenetrationTestByExecutorDTO, error) {
	return ptrs.repo.GetPenetrationTestByExecutor(params)
}

func (ptrs *PenetrationTestReportService) GetPenetrationTestByResponsibleUser(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.PenetrationTestByResponsibleUserDTO, error) {
	return ptrs.repo.GetPenetrationTestByResponsibleUser(params)
}

func (ptrs *PenetrationTestReportService) GetPenetrationTestByCreatedUser(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.PenetrationTestByCreatedUserDTO, error) {
	return ptrs.repo.GetPenetrationTestByCreatedUser(params)
}

func (ptrs *PenetrationTestReportService) GetPenetrationTestByResponsibleUnit(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.PenetrationTestByResponsibleUnitDTO, error) {
	return ptrs.repo.GetPenetrationTestByResponsibleUnit(params)
}

func (ptrs *PenetrationTestReportService) GetPenetrationTestByYear(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.PenetrationTestByYearDTO, error) {
	return ptrs.repo.GetPenetrationTestByYear(params)
}

func (ptrs *PenetrationTestReportService) GetPenetrationTestByPeriod(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.PenetrationTestByPeriodDTO, error) {
	return ptrs.repo.GetPenetrationTestByPeriod(params)
}

func (ptrs *PenetrationTestReportService) GetPenetrationTestFindingByStatus(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.PenetrationTestFindingByStatusDTO, error) {
	return ptrs.repo.GetPenetrationTestFindingByStatus(params)
}

func (ptrs *PenetrationTestReportService) GetPenetrationTestFindingByResponsibleUnit(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.PenetrationTestFindingByResponsibleUnitDTO, error) {
	return ptrs.repo.GetPenetrationTestFindingByResponsibleUnit(params)
}

func (ptrs *PenetrationTestReportService) GetPenetrationTestFindingByResponsibleUser(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.PenetrationTestFindingByResponsibleUserDTO, error) {
	return ptrs.repo.GetPenetrationTestFindingByResponsibleUser(params)
}

func (ptrs *PenetrationTestReportService) GetPenetrationTestFindingBySeverityLevel(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.PenetrationTestFindingBySeverityLevelDTO, error) {
	return ptrs.repo.GetPenetrationTestFindingBySeverityLevel(params)
}
