package services

import (
	organizationServices "jethur/grc/organization/services"
	reportDTOs "jethur/grc/reports/dtos"
	"jethur/grc/reports/repositories"
)

type ComplianceReportServiceInterface interface {
	GetComplianceByComplianceCategory(params *reportDTOs.ReportListQueryParamsDTO, currentUserID uint) ([]*reportDTOs.ComplianceByCategoryDTO, error)
	GetComplianceBySeverityLevel(params *reportDTOs.ReportListQueryParamsDTO, currentUserID uint) ([]*reportDTOs.ComplianceBySeverityLevelDTO, error)
	GetComplianceByComplianceStatus(params *reportDTOs.ReportListQueryParamsDTO, currentUserID uint) ([]*reportDTOs.ComplianceByStatusDTO, error)
	GetComplianceByResponsibleUnit(params *reportDTOs.ReportListQueryParamsDTO, currentUserID uint) ([]*reportDTOs.ComplianceByResponsibleUnitDTO, error)
	GetComplianceByResponsibleUser(params *reportDTOs.ReportListQueryParamsDTO, currentUserID uint) ([]*reportDTOs.ComplianceByResponsibleUserDTO, error)

	GetComplianceRequirementByComplianceStatus(params *reportDTOs.ReportListQueryParamsDTO, currentUserID uint) ([]*reportDTOs.ComplianceRequirementByComplianceStatusDTO, error)
	GetComplianceRequirementBySeverity(params *reportDTOs.ReportListQueryParamsDTO, currentUserID uint) ([]*reportDTOs.ComplianceRequirementBySeverityLevelDTO, error)
	GetComplianceRequirementByCategory(params *reportDTOs.ReportListQueryParamsDTO, currentUserID uint) ([]*reportDTOs.ComplianceRequirementByCategoryDTO, error)

	GetComplianceRequirementByResponsibleUnit(params *reportDTOs.ReportListQueryParamsDTO, currentUserID uint) ([]*reportDTOs.ComplianceRequirementByResponsibleUnitDTO, error)
	GetComplianceRequirementByResponsibleUser(params *reportDTOs.ReportListQueryParamsDTO, currentUserID uint) ([]*reportDTOs.ComplianceRequirementByResponsibleUserDTO, error)

	GetComplianceActionPlanByStatus(params *reportDTOs.ReportListQueryParamsDTO, currentUserID uint) ([]*reportDTOs.ComplianceActionPlanByStatusDTO, error)
	GetComplianceActionPlanByResponsibleUser(params *reportDTOs.ReportListQueryParamsDTO, currentUserID uint) ([]*reportDTOs.ComplianceActionPlanByResponsibleUserDTO, error)

	GetSLAByCategory(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.SLAByCategoryDTO, error)
	GetSLABySeverityLevel(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.SLABySeverityLevelDTO, error)
	GetSLAByComplianceStatus(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.SLAByComplianceStatusDTO, error)
}

type ComplianceReportService struct {
	repo    repositories.ComplianceReportRepositoryInterface
	userSvc organizationServices.UserServiceInterface
}

func NewComplianceReportService(repo repositories.ComplianceReportRepositoryInterface, userSvc organizationServices.UserServiceInterface) ComplianceReportServiceInterface {
	return &ComplianceReportService{repo: repo, userSvc: userSvc}
}

func (crs *ComplianceReportService) setCyberFilter(params *reportDTOs.ReportListQueryParamsDTO, currentUserID uint) error {
	isCyberUser, err := crs.userSvc.IsCyberUser(currentUserID)
	if err != nil {
		return err
	}
	params.IsCyberUser = isCyberUser
	return nil
}

func (crs *ComplianceReportService) GetComplianceByComplianceCategory(params *reportDTOs.ReportListQueryParamsDTO, currentUserID uint) ([]*reportDTOs.ComplianceByCategoryDTO, error) {
	if err := crs.setCyberFilter(params, currentUserID); err != nil {
		return nil, err
	}
	return crs.repo.GetComplianceByComplianceCategory(params)
}

func (crs *ComplianceReportService) GetComplianceBySeverityLevel(params *reportDTOs.ReportListQueryParamsDTO, currentUserID uint) ([]*reportDTOs.ComplianceBySeverityLevelDTO, error) {
	if err := crs.setCyberFilter(params, currentUserID); err != nil {
		return nil, err
	}
	return crs.repo.GetComplianceBySeverityLevel(params)
}

func (crs *ComplianceReportService) GetComplianceByComplianceStatus(params *reportDTOs.ReportListQueryParamsDTO, currentUserID uint) ([]*reportDTOs.ComplianceByStatusDTO, error) {
	if err := crs.setCyberFilter(params, currentUserID); err != nil {
		return nil, err
	}
	return crs.repo.GetComplianceByComplianceStatus(params)
}

func (crs *ComplianceReportService) GetComplianceByResponsibleUnit(params *reportDTOs.ReportListQueryParamsDTO, currentUserID uint) ([]*reportDTOs.ComplianceByResponsibleUnitDTO, error) {
	if err := crs.setCyberFilter(params, currentUserID); err != nil {
		return nil, err
	}
	return crs.repo.GetComplianceByResponsibleUnit(params)
}

func (crs *ComplianceReportService) GetComplianceByResponsibleUser(params *reportDTOs.ReportListQueryParamsDTO, currentUserID uint) ([]*reportDTOs.ComplianceByResponsibleUserDTO, error) {
	if err := crs.setCyberFilter(params, currentUserID); err != nil {
		return nil, err
	}
	return crs.repo.GetComplianceByResponsibleUser(params)
}

func (crs *ComplianceReportService) GetComplianceRequirementByComplianceStatus(params *reportDTOs.ReportListQueryParamsDTO, currentUserID uint) ([]*reportDTOs.ComplianceRequirementByComplianceStatusDTO, error) {
	if err := crs.setCyberFilter(params, currentUserID); err != nil {
		return nil, err
	}
	return crs.repo.GetComplianceRequirementByComplianceStatus(params)
}

func (crs *ComplianceReportService) GetComplianceRequirementBySeverity(params *reportDTOs.ReportListQueryParamsDTO, currentUserID uint) ([]*reportDTOs.ComplianceRequirementBySeverityLevelDTO, error) {
	if err := crs.setCyberFilter(params, currentUserID); err != nil {
		return nil, err
	}
	return crs.repo.GetComplianceRequirementBySeverity(params)
}

func (crs *ComplianceReportService) GetComplianceRequirementByCategory(params *reportDTOs.ReportListQueryParamsDTO, currentUserID uint) ([]*reportDTOs.ComplianceRequirementByCategoryDTO, error) {
	if err := crs.setCyberFilter(params, currentUserID); err != nil {
		return nil, err
	}
	return crs.repo.GetComplianceRequirementByCategory(params)
}

func (crs *ComplianceReportService) GetComplianceRequirementByResponsibleUnit(params *reportDTOs.ReportListQueryParamsDTO, currentUserID uint) ([]*reportDTOs.ComplianceRequirementByResponsibleUnitDTO, error) {
	if err := crs.setCyberFilter(params, currentUserID); err != nil {
		return nil, err
	}
	return crs.repo.GetComplianceRequirementByResponsibleUnit(params)
}

func (crs *ComplianceReportService) GetComplianceRequirementByResponsibleUser(params *reportDTOs.ReportListQueryParamsDTO, currentUserID uint) ([]*reportDTOs.ComplianceRequirementByResponsibleUserDTO, error) {
	if err := crs.setCyberFilter(params, currentUserID); err != nil {
		return nil, err
	}
	return crs.repo.GetComplianceRequirementByResponsibleUser(params)
}

func (crs *ComplianceReportService) GetComplianceActionPlanByStatus(params *reportDTOs.ReportListQueryParamsDTO, currentUserID uint) ([]*reportDTOs.ComplianceActionPlanByStatusDTO, error) {
	if err := crs.setCyberFilter(params, currentUserID); err != nil {
		return nil, err
	}
	return crs.repo.GetComplianceActionPlanByStatus(params)
}

func (crs *ComplianceReportService) GetComplianceActionPlanByResponsibleUser(params *reportDTOs.ReportListQueryParamsDTO, currentUserID uint) ([]*reportDTOs.ComplianceActionPlanByResponsibleUserDTO, error) {
	if err := crs.setCyberFilter(params, currentUserID); err != nil {
		return nil, err
	}
	return crs.repo.GetComplianceActionPlanByResponsibleUser(params)
}

func (crs *ComplianceReportService) GetSLAByCategory(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.SLAByCategoryDTO, error) {
	return crs.repo.GetSLAByCategory(params)
}

func (crs *ComplianceReportService) GetSLABySeverityLevel(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.SLABySeverityLevelDTO, error) {
	return crs.repo.GetSLABySeverityLevel(params)
}

func (crs *ComplianceReportService) GetSLAByComplianceStatus(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.SLAByComplianceStatusDTO, error) {
	return crs.repo.GetSLAByComplianceStatus(params)
}
