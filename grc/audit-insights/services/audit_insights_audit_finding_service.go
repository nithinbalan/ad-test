package services

import (
	"jethur/grc/audit-insights/repositories"
	auditDTOs "jethur/grc/audit/dtos"
	auditServices "jethur/grc/audit/services"
	logServices "jethur/grc/logs/services"
	organizationServices "jethur/grc/organization/services"
	"jethur/utility"
)

type AuditInsightsAuditFindingServiceInterface interface {
	GetAll(params *auditDTOs.AuditFindingListQueryParamsDTO, currentUserID uint) (*auditDTOs.AuditFindingListResponseDTO, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*auditDTOs.AuditFindingDetailDTO, error)
	GetAuditInsightFindingCountByCategory(auditId uint) ([]auditDTOs.AuditFindingCategoryCountResponseDTO, error)
}

type AuditInsightsAuditFindingService struct {
	repo            repositories.AuditInsightsAuditFindingRepositoryInterface
	activityLogSvc  logServices.ActivityLogServiceInterface
	userSvc         organizationServices.UserServiceInterface
	auditFindingSvc auditServices.AuditFindingServiceInterface
}

func NewAuditInsightsAuditFindingService(
	repo repositories.AuditInsightsAuditFindingRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	auditFindingSvc auditServices.AuditFindingServiceInterface,
) AuditInsightsAuditFindingServiceInterface {
	return &AuditInsightsAuditFindingService{
		repo:            repo,
		activityLogSvc:  activityLogSvc,
		userSvc:         userSvc,
		auditFindingSvc: auditFindingSvc,
	}
}

func (aiafs *AuditInsightsAuditFindingService) GetAll(params *auditDTOs.AuditFindingListQueryParamsDTO, currentUserID uint) (*auditDTOs.AuditFindingListResponseDTO, error) {
	orgUnitID, err := aiafs.userSvc.GetUserPrimaryUnitID(currentUserID)
	if err != nil {
		return nil, err
	}
	return aiafs.repo.GetAll(params, orgUnitID)
}

func (aiafs *AuditInsightsAuditFindingService) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*auditDTOs.AuditFindingDetailDTO, error) {
	auditFinding, err := aiafs.auditFindingSvc.GetFullDetailsByID(id, tx...)
	if err != nil {
		return nil, err
	}
	return auditFinding, nil
}

func (aiafs *AuditInsightsAuditFindingService) GetAuditInsightFindingCountByCategory(auditId uint) ([]auditDTOs.AuditFindingCategoryCountResponseDTO, error) {
	response, err := aiafs.auditFindingSvc.GetAuditFindingCountByCategory(auditId)
	if err != nil {
		return nil, err
	}

	return response, nil
}
