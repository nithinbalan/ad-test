package services

import (
	reportDTOs "jethur/grc/reports/dtos"
	"jethur/grc/reports/repositories"
)

type DocumentReportServiceInterface interface {
	GetDocumentsByStatus(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.DocumentByStatusDTO, error)
	GetDocumentsByType(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.DocumentByTypeDTO, error)
	GetDocumentsByResponsibleUnit(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.DocumentByResponsibleUnitDTO, error)
	GetDocumentsByCategory(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.DocumentByCategoryDTO, error)
	GetDocumentsByResponsibleUser(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.DocumentByResponsibleUserDTO, error)
	GetDocumentsByAccess(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.DocumentByAccessDTO, error)
	GetDocumentsByPolicyAndProcedure(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.DocumentByTypeDTO, error)
	GetDocumentsByPublishedAndArchived(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.DocumentByStatusDTO, error)
	GetDocumentsPublishedByYear(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.DocumentByYearDTO, error)
	GetDocumentsPublishedByMonth(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.DocumentByMonthDTO, error)
	GetDocumentChangeRequestsByType(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.DocumentChangeRequestByTypeDTO, error)
	GetDocumentChangeRequestsByStatus(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.DocumentChangeRequestByStatusDTO, error)
	GetDocumentsUpdatedByChangeType(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.DocumentChangeRequestByTypeDTO, error)
	GetDocumentsArchivedByQuarter(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.DocumentsArchivedByQuarterDTO, error)
}

type DocumentReportService struct {
	repo repositories.DocumentReportRepositoryInterface
}

func NewDocumentReportService(repo repositories.DocumentReportRepositoryInterface) DocumentReportServiceInterface {
	return &DocumentReportService{repo: repo}
}

func (drs *DocumentReportService) GetDocumentsByStatus(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.DocumentByStatusDTO, error) {
	return drs.repo.GetDocumentsByStatus(params)
}

func (drs *DocumentReportService) GetDocumentsByType(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.DocumentByTypeDTO, error) {
	return drs.repo.GetDocumentsByType(params)
}

func (drs *DocumentReportService) GetDocumentsByResponsibleUnit(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.DocumentByResponsibleUnitDTO, error) {
	return drs.repo.GetDocumentsByResponsibleUnit(params)
}

func (drs *DocumentReportService) GetDocumentsByCategory(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.DocumentByCategoryDTO, error) {
	return drs.repo.GetDocumentsByCategory(params)
}

func (drs *DocumentReportService) GetDocumentsByResponsibleUser(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.DocumentByResponsibleUserDTO, error) {
	return drs.repo.GetDocumentsByResponsibleUser(params)
}

func (drs *DocumentReportService) GetDocumentsByAccess(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.DocumentByAccessDTO, error) {
	return drs.repo.GetDocumentsByAccess(params)
}

func (drs *DocumentReportService) GetDocumentsByPolicyAndProcedure(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.DocumentByTypeDTO, error) {
	return drs.repo.GetDocumentsByPolicyAndProcedure(params)
}

func (drs *DocumentReportService) GetDocumentsByPublishedAndArchived(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.DocumentByStatusDTO, error) {
	return drs.repo.GetDocumentsByPublishedAndArchived(params)
}

func (drs *DocumentReportService) GetDocumentsPublishedByYear(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.DocumentByYearDTO, error) {
	return drs.repo.GetDocumentsPublishedByYear(params)
}

func (drs *DocumentReportService) GetDocumentsPublishedByMonth(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.DocumentByMonthDTO, error) {
	return drs.repo.GetDocumentsPublishedByMonth(params)
}

func (drs *DocumentReportService) GetDocumentChangeRequestsByType(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.DocumentChangeRequestByTypeDTO, error) {
	return drs.repo.GetDocumentChangeRequestsByType(params)
}

func (drs *DocumentReportService) GetDocumentChangeRequestsByStatus(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.DocumentChangeRequestByStatusDTO, error) {
	return drs.repo.GetDocumentChangeRequestsByStatus(params)
}

func (drs *DocumentReportService) GetDocumentsUpdatedByChangeType(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.DocumentChangeRequestByTypeDTO, error) {
	return drs.repo.GetDocumentsUpdatedByChangeType(params)
}

func (drs *DocumentReportService) GetDocumentsArchivedByQuarter(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.DocumentsArchivedByQuarterDTO, error) {
	return drs.repo.GetDocumentsArchivedByQuarter(params)
}
