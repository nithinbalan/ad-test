package services

import (
	auditDTOs "jethur/grc/audit/dtos"
	"jethur/grc/audit/repositories"
	logServices "jethur/grc/logs/services"
	organizationServices "jethur/grc/organization/services"
	"jethur/utility"
)

type AuditInsightsAuditServiceInterface interface {
	GetAll(params *auditDTOs.AuditListQueryParamsDTO, currentUserID uint) (*auditDTOs.AuditListResponseDTO, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*auditDTOs.AuditDetailDTO, error)
}

type AuditInsightsAuditService struct {
	repo           repositories.AuditInsightsAuditRepositoryInterface
	activityLogSvc logServices.ActivityLogServiceInterface
	userSvc        organizationServices.UserServiceInterface
	auditSvc       AuditServiceInterface
}

func NewAuditInsightsAuditService(
	repo repositories.AuditInsightsAuditRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	auditSvc AuditServiceInterface,
) AuditInsightsAuditServiceInterface {
	return &AuditInsightsAuditService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
		userSvc:        userSvc,
		auditSvc:       auditSvc,
	}
}

func (aias *AuditInsightsAuditService) GetAll(params *auditDTOs.AuditListQueryParamsDTO, currentUserID uint) (*auditDTOs.AuditListResponseDTO, error) {
	orgUnitID, err := aias.userSvc.GetUserPrimaryUnitID(currentUserID)
	if err != nil {
		return nil, err
	}
	return aias.repo.GetAll(params, orgUnitID)
}

func (aias *AuditInsightsAuditService) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*auditDTOs.AuditDetailDTO, error) {
	audit, err := aias.auditSvc.GetFullDetailsByID(id, tx...)
	if err != nil {
		return nil, err
	}
	return audit, nil
}
