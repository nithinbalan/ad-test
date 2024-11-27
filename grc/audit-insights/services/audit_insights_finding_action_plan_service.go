package services

import (
	"errors"
	"jethur/grc/audit-insights/repositories"
	auditDTOs "jethur/grc/audit/dtos"
	auditServices "jethur/grc/audit/services"
	logServices "jethur/grc/logs/services"
	organizationServices "jethur/grc/organization/services"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type AuditInsightsFindingActionPlanServiceInterface interface {
	GetAll(params *auditDTOs.AuditFindingActionPlanListQueryParamsDTO, currentUserID uint) (*auditDTOs.AuditFindingActionPlanListResponseDTO, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*auditDTOs.AuditFindingActionPlanDetailDTO, error)
}

type AuditInsightsFindingActionPlanService struct {
	repo            repositories.AuditInsightsFindingActionPlanRepositoryInterface
	activityLogSvc  logServices.ActivityLogServiceInterface
	userSvc         organizationServices.UserServiceInterface
	auditFindingActionPlanSvc auditServices.AuditFindingActionPlanServiceInterface
}

func NewAuditInsightsFindingActionPlanService(
	repo repositories.AuditInsightsFindingActionPlanRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	auditFindingActionPlanSvc auditServices.AuditFindingActionPlanServiceInterface,
) AuditInsightsFindingActionPlanServiceInterface {
	return &AuditInsightsFindingActionPlanService{
		repo:            repo,
		activityLogSvc:  activityLogSvc,
		userSvc:         userSvc,
		auditFindingActionPlanSvc: auditFindingActionPlanSvc,
	}
}

func (aifaps *AuditInsightsFindingActionPlanService) GetAll(params *auditDTOs.AuditFindingActionPlanListQueryParamsDTO, currentUserID uint) (*auditDTOs.AuditFindingActionPlanListResponseDTO, error) {
	orgUnitID, err := aifaps.userSvc.GetUserPrimaryUnitID(currentUserID)
	if err != nil {
		return nil, err
	}
	return aifaps.repo.GetAll(params, orgUnitID)
}

func (aifaps *AuditInsightsFindingActionPlanService) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*auditDTOs.AuditFindingActionPlanDetailDTO, error) {
	var auditFindingActionPlanDTO auditDTOs.AuditFindingActionPlanDetailDTO
	auditFindingActionPlan, err := aifaps.repo.GetFullDetailsByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	err = mapstructure.Decode(auditFindingActionPlan, &auditFindingActionPlanDTO)
	if err != nil {
		return nil, err
	}

	// todo :- also add extra fields which required in details screen
	if auditFindingActionPlan.TargetDate != nil {
		targetDate := auditFindingActionPlan.TargetDate.Format("2006-01-02")
		auditFindingActionPlanDTO.TargetDateString = &targetDate
	}

	auditFindingActionPlanDTO.CreatedAt = auditFindingActionPlan.CreatedAt
	auditFindingActionPlanDTO.UpdatedAt = auditFindingActionPlan.UpdatedAt

	return &auditFindingActionPlanDTO, nil
}
