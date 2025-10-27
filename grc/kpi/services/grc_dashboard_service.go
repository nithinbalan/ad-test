package services

import (
	"fmt"
	documentServices "jethur/grc/documents/services"
	kpiDTOs "jethur/grc/kpi/dtos"
	"jethur/grc/kpi/repositories"
	organizationServices "jethur/grc/organization/services"
)

type GRCDashboardServiceInterface interface {
	GetGRCComplianceSummary(params *kpiDTOs.GRCDashboardListQueryParamsDTO) (*kpiDTOs.GRCComplianceSummaryDTO, error)
	GetGRCRiskSummary(params *kpiDTOs.GRCDashboardListQueryParamsDTO) (*kpiDTOs.GRCRiskSummaryDTO, error)
	GetGRCAuditSummary(params *kpiDTOs.GRCDashboardListQueryParamsDTO) (*kpiDTOs.GRCAuditSummaryDTO, error)
	GetGRCDocumentSummary(params *kpiDTOs.GRCDashboardListQueryParamsDTO) (*kpiDTOs.GRCDocumentSummaryDTO, error)
	GetGRCOrganizationSummary(params *kpiDTOs.GRCDashboardListQueryParamsDTO) (*kpiDTOs.GRCOrganizationSummaryDTO, error)
	GetGRCIncidentSummary(params *kpiDTOs.GRCDashboardListQueryParamsDTO) (*kpiDTOs.GRCIncidentSummaryDTO, error)

	GetGRCComplianceOverview(params *kpiDTOs.GRCDashboardListQueryParamsDTO) (*kpiDTOs.GRCComplianceOverviewDTO, error)
	GetGRCRiskOverview(params *kpiDTOs.GRCDashboardListQueryParamsDTO) (*kpiDTOs.GRCRiskOverviewDTO, error)
	GetGRCDocumentOverview(params *kpiDTOs.GRCDashboardListQueryParamsDTO) (*kpiDTOs.GRCDocumentOverviewDTO, error)

	GetGRCTopCompliances(params *kpiDTOs.GRCDashboardListQueryParamsDTO) ([]*kpiDTOs.GRCTopComplianceRegistersDTO, error)
	GetGRCTopRisks(params *kpiDTOs.GRCDashboardListQueryParamsDTO) ([]*kpiDTOs.GRCTopRiskRegistersDTO, error)
	GetGRCTopPublishedDocuments(params *kpiDTOs.GRCDashboardListQueryParamsDTO) ([]*kpiDTOs.GRCTopDocumentsDTO, error)

	GetGRCDocumentTypeChart(params *kpiDTOs.GRCDashboardListQueryParamsDTO) ([]*kpiDTOs.GRCDocumentTypeChartDTO, error)
	GetPolicyDocuments(params *kpiDTOs.GRCDashboardDocumentListQueryParamsDTO) (*kpiDTOs.GRCDashboardDocumentListResponseDTO, error)
}

type GRCDashboardService struct {
	repo         repositories.GRCDashboardRepositoryInterface
	grcReportSvc organizationServices.GovernanceReportServiceInterface
	documentSvc  documentServices.DocumentServiceInterface
}

func NewGRCDashboardService(repo repositories.GRCDashboardRepositoryInterface, grcReportSvc organizationServices.GovernanceReportServiceInterface, documentSvc documentServices.DocumentServiceInterface) GRCDashboardServiceInterface {
	return &GRCDashboardService{
		repo:         repo,
		grcReportSvc: grcReportSvc,
		documentSvc:  documentSvc,
	}
}

func (gds *GRCDashboardService) GetGRCComplianceSummary(params *kpiDTOs.GRCDashboardListQueryParamsDTO) (*kpiDTOs.GRCComplianceSummaryDTO, error) {
	return gds.repo.GetGRCComplianceSummary(params)
}

func (gds *GRCDashboardService) GetGRCRiskSummary(params *kpiDTOs.GRCDashboardListQueryParamsDTO) (*kpiDTOs.GRCRiskSummaryDTO, error) {
	return gds.repo.GetGRCRiskSummary(params)
}

func (gds *GRCDashboardService) GetGRCAuditSummary(params *kpiDTOs.GRCDashboardListQueryParamsDTO) (*kpiDTOs.GRCAuditSummaryDTO, error) {
	return gds.repo.GetGRCAuditSummary(params)
}

func (gds *GRCDashboardService) GetGRCDocumentSummary(params *kpiDTOs.GRCDashboardListQueryParamsDTO) (*kpiDTOs.GRCDocumentSummaryDTO, error) {
	return gds.repo.GetGRCDocumentSummary(params)
}

func (gds *GRCDashboardService) GetGRCOrganizationSummary(params *kpiDTOs.GRCDashboardListQueryParamsDTO) (*kpiDTOs.GRCOrganizationSummaryDTO, error) {
	return gds.repo.GetGRCOrganizationSummary(params)
}

func (gds *GRCDashboardService) GetGRCIncidentSummary(params *kpiDTOs.GRCDashboardListQueryParamsDTO) (*kpiDTOs.GRCIncidentSummaryDTO, error) {
	return gds.repo.GetGRCIncidentSummary(params)
}

func (gds *GRCDashboardService) GetGRCComplianceOverview(params *kpiDTOs.GRCDashboardListQueryParamsDTO) (*kpiDTOs.GRCComplianceOverviewDTO, error) {
	return gds.repo.GetGRCComplianceOverview(params)
}

func (gds *GRCDashboardService) GetGRCRiskOverview(params *kpiDTOs.GRCDashboardListQueryParamsDTO) (*kpiDTOs.GRCRiskOverviewDTO, error) {
	return gds.repo.GetGRCRiskOverview(params)
}

func (gds *GRCDashboardService) GetGRCDocumentOverview(params *kpiDTOs.GRCDashboardListQueryParamsDTO) (*kpiDTOs.GRCDocumentOverviewDTO, error) {
	return gds.repo.GetGRCDocumentOverview(params)
}

func (gds *GRCDashboardService) GetGRCTopCompliances(params *kpiDTOs.GRCDashboardListQueryParamsDTO) ([]*kpiDTOs.GRCTopComplianceRegistersDTO, error) {
	return gds.repo.GetTopComplianceRegisters(params)
}

func (gds *GRCDashboardService) GetGRCTopRisks(params *kpiDTOs.GRCDashboardListQueryParamsDTO) ([]*kpiDTOs.GRCTopRiskRegistersDTO, error) {
	return gds.repo.GetTopRiskRegisters(params)
}

func (gds *GRCDashboardService) GetGRCTopPublishedDocuments(params *kpiDTOs.GRCDashboardListQueryParamsDTO) ([]*kpiDTOs.GRCTopDocumentsDTO, error) {
	return gds.repo.GetTopPublishedDocuments(params)
}

func (gds *GRCDashboardService) GetGRCDocumentTypeChart(params *kpiDTOs.GRCDashboardListQueryParamsDTO) ([]*kpiDTOs.GRCDocumentTypeChartDTO, error) {
	return gds.repo.GetGRCDocumentTypeChart(params)
}

func (gds *GRCDashboardService) GetPolicyDocuments(params *kpiDTOs.GRCDashboardDocumentListQueryParamsDTO) (*kpiDTOs.GRCDashboardDocumentListResponseDTO, error) {
	// return gds.repo.GetPolicyDocuments(params)
	var documentDTO []*kpiDTOs.GRCDashboardDocumentResponseDTO
	documents, err := gds.repo.GetPolicyDocuments(params)
	if err != nil {
		return nil, err
	}

	for _, document := range documents.Items {
		fmt.Println(document)
		policy, err := gds.repo.GetDocumentByID(document.ID)
		if err != nil {
			return nil, err
		}
		summary, err := gds.grcReportSvc.GetDocumentAcknowledgmentSummary(policy)
		if err != nil {
			return nil, fmt.Errorf("failed to get acknowledgment summary for document ID %d: %w", document.ID, err)
		}

		docResponse := kpiDTOs.GRCDashboardDocumentResponseDTO{
			ID:                         document.ID,
			ReferenceNumber:            document.ReferenceNumber,
			DocumentTitleEn:            policy.TitleEn,
			DocumentTitleAr:            policy.TitleAr,
			ResponsibleUserFirstNameEn: document.ResponsibleUserFirstNameEn,
			ResponsibleUserFirstNameAr: document.ResponsibleUserFirstNameAr,
			ResponsibleUserLastNameEn:  document.ResponsibleUserLastNameEn,
			ResponsibleUserLastNameAr:  document.ResponsibleUserLastNameAr,
			ResponsibleUnitTitleEn:     document.ResponsibleUnitTitleEn,
			ResponsibleUnitTitleAr:     document.ResponsibleUnitTitleAr,
			TotalEmployees:             int(summary.TotalUsers),
			AcknowledgedCount:          int(summary.AcknowledgedUsers),
			PendingCount:               int(summary.NotAcknowledgedUsers),
			CreatedByID:                document.CreatedByID,
		}
		documentDTO = append(documentDTO, &docResponse)

	}
	return documents, nil
}
