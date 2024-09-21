package services

import (
	"jethur/grc/audit-insights/repositories"
	auditDTOs "jethur/grc/audit/dtos"
	auditServices "jethur/grc/audit/services"
	logServices "jethur/grc/logs/services"
	organizationServices "jethur/grc/organization/services"

	"jethur/utility"
)

type AuditInsightsAuditPlanServiceInterface interface {
	GetAll(params *auditDTOs.AuditPlanListQueryParamsDTO, currentUserID uint) (*auditDTOs.AuditPlanListResponseDTO, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*auditDTOs.AuditPlanDetailDTO, error)
}

type AuditInsightsAuditPlanService struct {
	repo           repositories.AuditInsightsAuditPlanRepositoryInterface
	activityLogSvc logServices.ActivityLogServiceInterface
	userSvc        organizationServices.UserServiceInterface
	auditPlanSvc   auditServices.AuditPlanServiceInterface
}

func NewAuditInsightsAuditPlanService(
	repo repositories.AuditInsightsAuditPlanRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	auditPlanSvc auditServices.AuditPlanServiceInterface,
) AuditInsightsAuditPlanServiceInterface {
	return &AuditInsightsAuditPlanService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
		userSvc:        userSvc,
		auditPlanSvc:   auditPlanSvc,
	}
}

func (aiaps *AuditInsightsAuditPlanService) GetAll(params *auditDTOs.AuditPlanListQueryParamsDTO, currentUserID uint) (*auditDTOs.AuditPlanListResponseDTO, error) {
	orgUnitID, err := aiaps.userSvc.GetUserPrimaryUnitID(currentUserID)
	if err != nil {
		return nil, err
	}
	return aiaps.repo.GetAll(params, orgUnitID)
}

func (aiaps *AuditInsightsAuditPlanService) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*auditDTOs.AuditPlanDetailDTO, error) {
	auditPlan, err := aiaps.auditPlanSvc.GetFullDetailsByID(id, tx...)
	if err != nil {
		return nil, err
	}
	return auditPlan, nil
}
