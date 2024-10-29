package services

import (
	"jethur/grc/reports/dtos"
	"jethur/grc/reports/repositories"
)

type OrganizationReportServiceInterface interface {
	GetUsersByStatus(params *dtos.ReportListQueryParamsDTO) ([]*dtos.UsersByStatusDTO, error)
	GetUsersByRole(params *dtos.ReportListQueryParamsDTO) ([]*dtos.UsersByRoleDTO, error)
	GetUsersByDesignation(params *dtos.ReportListQueryParamsDTO) ([]*dtos.UsersByDesignationDTO, error)

	GetProcessesByResponsibleUnit(params *dtos.ReportListQueryParamsDTO) ([]*dtos.ProcessByUnitDTO, error)
	GetProcessesByRiskRating(params *dtos.ReportListQueryParamsDTO) ([]*dtos.ProcessByRiskRatingDTO, error)
	GetProcessesByResponsibleUser(params *dtos.ReportListQueryParamsDTO) ([]*dtos.ProcessByResponsibleUserDTO, error)
	GetProcessesByProcessGroup(params *dtos.ReportListQueryParamsDTO) ([]*dtos.ProcessByProcessGroupDTO, error)

	GetControlsByControlType(params *dtos.ReportListQueryParamsDTO) ([]*dtos.ControlsByControlTypeDTO, error)
	GetControlsByCategory(params *dtos.ReportListQueryParamsDTO) ([]*dtos.ControlsByControlCategoryDTO, error)
	GetControlsByComplianceRegister(params *dtos.ReportListQueryParamsDTO) ([]*dtos.ControlsByComplianceRegisterDTO, error)
	GetControlsByRiskRating(params *dtos.ReportListQueryParamsDTO) ([]*dtos.ControlsByRiskRatingDTO, error)
	GetControlsByImplementationStatus(params *dtos.ReportListQueryParamsDTO) ([]*dtos.ControlsByImplementationStatusDTO, error)
	GetControlsByEfficiencyStatus(params *dtos.ReportListQueryParamsDTO) ([]*dtos.ControlsByEfficiencyDTO, error)
}

type OrganizationReportService struct {
	repo repositories.OrganizationReportRepositoryInterface
}

func NewOrganizationReportService(repo repositories.OrganizationReportRepositoryInterface) OrganizationReportServiceInterface {
	return &OrganizationReportService{repo: repo}
}

func (ors *OrganizationReportService) GetUsersByStatus(params *dtos.ReportListQueryParamsDTO) ([]*dtos.UsersByStatusDTO, error) {
	return ors.repo.GetUsersByStatus(params)
}

func (ors *OrganizationReportService) GetUsersByRole(params *dtos.ReportListQueryParamsDTO) ([]*dtos.UsersByRoleDTO, error) {
	return ors.repo.GetUsersByRole(params)
}

func (ors *OrganizationReportService) GetUsersByDesignation(params *dtos.ReportListQueryParamsDTO) ([]*dtos.UsersByDesignationDTO, error) {
	return ors.repo.GetUsersByDesignation(params)
}

func (ors *OrganizationReportService) GetProcessesByResponsibleUnit(params *dtos.ReportListQueryParamsDTO) ([]*dtos.ProcessByUnitDTO, error) {
	return ors.repo.GetProcessesByResponsibleUnit(params)
}

func (ors *OrganizationReportService) GetProcessesByRiskRating(params *dtos.ReportListQueryParamsDTO) ([]*dtos.ProcessByRiskRatingDTO, error) {
	return ors.repo.GetProcessesByRiskRating(params)
}

func (ors *OrganizationReportService) GetProcessesByResponsibleUser(params *dtos.ReportListQueryParamsDTO) ([]*dtos.ProcessByResponsibleUserDTO, error) {
	return ors.repo.GetProcessesByResponsibleUser(params)
}

func (ors *OrganizationReportService) GetProcessesByProcessGroup(params *dtos.ReportListQueryParamsDTO) ([]*dtos.ProcessByProcessGroupDTO, error) {
	return ors.repo.GetProcessesByProcessGroup(params)
}

func (ors *OrganizationReportService) GetControlsByControlType(params *dtos.ReportListQueryParamsDTO) ([]*dtos.ControlsByControlTypeDTO, error) {
	return ors.repo.GetControlsByControlType(params)
}

func (ors *OrganizationReportService) GetControlsByCategory(queryParams *dtos.ReportListQueryParamsDTO) ([]*dtos.ControlsByControlCategoryDTO, error) {
	return ors.repo.GetControlsByCategory(queryParams)
}

func (ors *OrganizationReportService) GetControlsByComplianceRegister(queryParams *dtos.ReportListQueryParamsDTO) ([]*dtos.ControlsByComplianceRegisterDTO, error) {
	return ors.repo.GetControlsByComplianceRegister(queryParams)
}

func (ors *OrganizationReportService) GetControlsByRiskRating(queryParams *dtos.ReportListQueryParamsDTO) ([]*dtos.ControlsByRiskRatingDTO, error) {
	return ors.repo.GetControlsByRiskRating(queryParams)
}

func (ors *OrganizationReportService) GetControlsByImplementationStatus(queryParams *dtos.ReportListQueryParamsDTO) ([]*dtos.ControlsByImplementationStatusDTO, error) {
	return ors.repo.GetControlsByImplementationStatus(queryParams)
}

func (ors *OrganizationReportService) GetControlsByEfficiencyStatus(queryParams *dtos.ReportListQueryParamsDTO) ([]*dtos.ControlsByEfficiencyDTO, error) {
	return ors.repo.GetControlsByEfficiencyStatus(queryParams)
}
