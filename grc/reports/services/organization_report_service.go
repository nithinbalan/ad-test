package services

import (
	"jethur/grc/reports/dtos"
	"jethur/grc/reports/repositories"
)

type OrganizationReportServiceInterface interface {
	GetUsersByStatus(params *dtos.ReportListQueryParamsDTO) ([]*dtos.UsersByStatusDTO, error)
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
