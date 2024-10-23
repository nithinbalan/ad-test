package services

import (
	"jethur/grc/reports/dtos"
	"jethur/grc/reports/repositories"
)

type OrganizationReportServiceInterface interface {
	GetUsersByStatus(params *dtos.ReportListQueryParamsDTO) ([]*dtos.UsersByStatusDTO, error)
	GetUsersByRole(params *dtos.ReportListQueryParamsDTO) ([]*dtos.UsersByRoleDTO, error)
	GetUsersByDesignation(params *dtos.ReportListQueryParamsDTO) ([]*dtos.UsersByDesignationDTO, error)
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
