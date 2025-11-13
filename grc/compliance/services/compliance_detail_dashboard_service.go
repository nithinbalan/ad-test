package services

import (
	complianceDTOs "jethur/grc/compliance/dtos"
	"jethur/grc/compliance/repositories"
)

type ComplianceDetailDashboardServiceInterface interface {
	GetComplianceDetailSummary(complianceRegisterID uint) (*complianceDTOs.ComplianceDetailSummaryDTO, error)
	GetComplianceClauseSummary(complianceRegisterID uint) (*complianceDTOs.ComplianceClauseSummaryDTO, error)
	GetTop5ClausePerformance(complianceRegisterID uint) ([]complianceDTOs.ComplianceClausePerformanceDTO, error)
	GetBottom5ClausePerformance(complianceRegisterID uint) ([]complianceDTOs.ComplianceClausePerformanceDTO, error)
	GetComplianceDomainWiseOverview(complianceRegisterID uint) ([]complianceDTOs.ComplianceDomainWiseOverviewDTO, error)
	GetPolicyByOrganizationUnit(complianceRegisterID uint) ([]complianceDTOs.CompliancePolicyByOrganizationUnitDTO, error)
	GetControlEffectivenessSummary(complianceRegisterID uint) (*complianceDTOs.ComplianceControlEffectivenessDTO, error)
	GetControlImplementationSummary(complianceRegisterID uint) (*complianceDTOs.ComplianceControlImplementationSummaryDTO, error)
	GetControlSummary(complianceRegisterID uint) (*complianceDTOs.ComplianceControlSummaryDTO, error)
	GetControlEfficiencySummary(complianceRegisterID uint) ([]complianceDTOs.ComplianceSummaryItemDTO, error)
	GetControlCoverageByDomain(complianceRegisterID uint) ([]complianceDTOs.ComplianceControlCoverageByDomainDTO, error)
	GetCompliancePolicySummary(complianceRegisterID uint) (*complianceDTOs.CompliancePolicySummaryDTO, error)
	GetPolicyAcknowledgementByOrganizationUnit(complianceRegisterID uint) ([]complianceDTOs.CompliancePolicyAcknowledgementOrganizationUnitDTO, error)
	GetActionPlanOverview(complianceRegisterID uint) (*complianceDTOs.ComplianceActionPlanOverviewDTO, error)
	GetActionPlanStatusSummary(complianceRegisterID uint) (*complianceDTOs.ComplianceActionPlanStatusSummaryDTO, error)
	GetActionPlansByOrganizationUnit(complianceRegisterID uint) ([]complianceDTOs.ComplianceActionPlanByOrganizationUnitDTO, error)
}

type ComplianceDetailDashboardService struct {
	repo repositories.ComplianceDetailDashboardRepositoryInterface
}

func NewComplianceDetailDashboardService(repo repositories.ComplianceDetailDashboardRepositoryInterface) ComplianceDetailDashboardServiceInterface {
	return &ComplianceDetailDashboardService{
		repo: repo,
	}
}

func (cdds *ComplianceDetailDashboardService) GetComplianceDetailSummary(complianceRegisterID uint) (*complianceDTOs.ComplianceDetailSummaryDTO, error) {
	return cdds.repo.GetComplianceDetailSummary(complianceRegisterID)
}

func (cdds *ComplianceDetailDashboardService) GetComplianceClauseSummary(complianceRegisterID uint) (*complianceDTOs.ComplianceClauseSummaryDTO, error) {
	return cdds.repo.GetComplianceClauseSummary(complianceRegisterID)
}

func (cdds *ComplianceDetailDashboardService) GetTop5ClausePerformance(complianceRegisterID uint) ([]complianceDTOs.ComplianceClausePerformanceDTO, error) {
	return cdds.repo.GetTop5ClausePerformance(complianceRegisterID)
}

func (cdds *ComplianceDetailDashboardService) GetBottom5ClausePerformance(complianceRegisterID uint) ([]complianceDTOs.ComplianceClausePerformanceDTO, error) {
	return cdds.repo.GetBottom5ClausePerformance(complianceRegisterID)
}

func (cdds *ComplianceDetailDashboardService) GetComplianceDomainWiseOverview(complianceRegisterID uint) ([]complianceDTOs.ComplianceDomainWiseOverviewDTO, error) {

	rawItems, err := cdds.repo.GetComplianceDomainWiseOverview(complianceRegisterID)
	if err != nil {
		return nil, err
	}

	domainMap := map[uint]*complianceDTOs.ComplianceDomainWiseOverviewDTO{}

	for _, row := range rawItems {

		if _, exists := domainMap[row.DomainID]; !exists {
			domainMap[row.DomainID] = &complianceDTOs.ComplianceDomainWiseOverviewDTO{
				DomainTitleEn:      row.DomainTitleEn,
				DomainTitleAr:      row.DomainTitleAr,
				ComplianceStatuses: []complianceDTOs.ComplianceStatusResponseDTO{},
			}
		}

		domainMap[row.DomainID].ComplianceStatuses = append(
			domainMap[row.DomainID].ComplianceStatuses,
			complianceDTOs.ComplianceStatusResponseDTO{
				TitleEn:   row.StatusTitleEn,
				TitleAr:   row.StatusTitleAr,
				ColorCode: row.ColorCode,
				Type:      row.StatusType,
				Count:     row.Count,
			},
		)
	}

	response := []complianceDTOs.ComplianceDomainWiseOverviewDTO{}
	for _, v := range domainMap {
		response = append(response, *v)
	}

	return response, nil
}

func (cdds *ComplianceDetailDashboardService) GetPolicyByOrganizationUnit(complianceRegisterID uint) ([]complianceDTOs.CompliancePolicyByOrganizationUnitDTO, error) {
	return cdds.repo.GetPolicyByOrganizationUnit(complianceRegisterID)
}

func (cdds *ComplianceDetailDashboardService) GetControlEffectivenessSummary(complianceRegisterID uint) (*complianceDTOs.ComplianceControlEffectivenessDTO, error) {
	return cdds.repo.GetControlEffectivenessSummary(complianceRegisterID)
}

func (cdds *ComplianceDetailDashboardService) GetControlImplementationSummary(complianceRegisterID uint) (*complianceDTOs.ComplianceControlImplementationSummaryDTO, error) {
	return cdds.repo.GetControlImplementationSummary(complianceRegisterID)
}

func (cdds *ComplianceDetailDashboardService) GetControlSummary(complianceRegisterID uint) (*complianceDTOs.ComplianceControlSummaryDTO, error) {
	return cdds.repo.GetControlSummary(complianceRegisterID)
}

func (cdds *ComplianceDetailDashboardService) GetControlEfficiencySummary(complianceRegisterID uint) ([]complianceDTOs.ComplianceSummaryItemDTO, error) {
	return cdds.repo.GetControlEfficiencySummary(complianceRegisterID)
}

func (cdds *ComplianceDetailDashboardService) GetControlCoverageByDomain(complianceRegisterID uint) ([]complianceDTOs.ComplianceControlCoverageByDomainDTO, error) {
	return cdds.repo.GetControlCoverageByDomain(complianceRegisterID)
}

func (cdds *ComplianceDetailDashboardService) GetCompliancePolicySummary(complianceRegisterID uint) (*complianceDTOs.CompliancePolicySummaryDTO, error) {
	return cdds.repo.GetCompliancePolicySummary(complianceRegisterID)
}

func (cdds *ComplianceDetailDashboardService) GetPolicyAcknowledgementByOrganizationUnit(complianceRegisterID uint) ([]complianceDTOs.CompliancePolicyAcknowledgementOrganizationUnitDTO, error) {
	return cdds.repo.GetPolicyAcknowledgementByOrganizationUnit(complianceRegisterID)
}

func (cdds *ComplianceDetailDashboardService) GetActionPlanOverview(complianceRegisterID uint) (*complianceDTOs.ComplianceActionPlanOverviewDTO, error) {
	return cdds.repo.GetActionPlanOverview(complianceRegisterID)
}

func (cdds *ComplianceDetailDashboardService) GetActionPlanStatusSummary(complianceRegisterID uint) (*complianceDTOs.ComplianceActionPlanStatusSummaryDTO, error) {
	return cdds.repo.GetActionPlanStatusSummary(complianceRegisterID)
}

func (cdds *ComplianceDetailDashboardService) GetActionPlansByOrganizationUnit(complianceRegisterID uint) ([]complianceDTOs.ComplianceActionPlanByOrganizationUnitDTO, error) {
	return cdds.repo.GetActionPlansByOrganizationUnit(complianceRegisterID)
}
