package services

import (
	"errors"
	auditDTOs "jethur/grc/audit/dtos"
	"jethur/grc/audit/repositories"
	logServices "jethur/grc/logs/services"
	organizationServices "jethur/grc/organization/services"

	"jethur/grc/models"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type AuditInsightsAuditPlanServiceInterface interface {
	GetAll(params *auditDTOs.AuditPlanListQueryParamsDTO, currentUserID uint) (*auditDTOs.AuditPlanListResponseDTO, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*auditDTOs.AuditPlanDetailDTO, error)
	GetByID(id uint) (*models.AuditPlan, error)
}

type AuditInsightsAuditPlanService struct {
	repo             repositories.AuditInsightsAuditPlanRepositoryInterface
	activityLogSvc   logServices.ActivityLogServiceInterface
	auditCategorySvc AuditCategoryServiceInterface
	auditStatusSvc   AuditStatusServiceInterface
	auditTeamSvc     AuditTeamServiceInterface
	userSvc          organizationServices.UserServiceInterface
	auditableItemSvc AuditableItemServiceInterface
}

func NewAuditInsightsAuditPlanService(
	repo repositories.AuditInsightsAuditPlanRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	auditCategorySvc AuditCategoryServiceInterface,
	auditStatusSvc AuditStatusServiceInterface,
	auditTeamSvc AuditTeamServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	auditableItemSvc AuditableItemServiceInterface,
) AuditInsightsAuditPlanServiceInterface {
	return &AuditInsightsAuditPlanService{
		repo:             repo,
		activityLogSvc:   activityLogSvc,
		auditCategorySvc: auditCategorySvc,
		auditStatusSvc:   auditStatusSvc,
		auditTeamSvc:     auditTeamSvc,
		userSvc:          userSvc,
		auditableItemSvc: auditableItemSvc,
	}
}

func (aiaps *AuditInsightsAuditPlanService) GetAll(params *auditDTOs.AuditPlanListQueryParamsDTO, currentUserID uint) (*auditDTOs.AuditPlanListResponseDTO, error) {
	orgUnitID, err := aiaps.userSvc.GetUserPrimaryUnitID(currentUserID)
	if err != nil {
		return nil, err
	}
	return aiaps.repo.GetAll(params, currentUserID, orgUnitID)
}

func (aiaps *AuditInsightsAuditPlanService) GetByID(id uint) (*models.AuditPlan, error) {
	auditPlan, err := aiaps.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return auditPlan, nil
}

func (aiaps *AuditInsightsAuditPlanService) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*auditDTOs.AuditPlanDetailDTO, error) {
	var auditPlanDTO auditDTOs.AuditPlanDetailDTO

	auditPlan, err := aiaps.repo.GetFullDetailsByID(id, tx...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	err = mapstructure.Decode(auditPlan, &auditPlanDTO)
	if err != nil {
		return nil, err
	}

	if auditPlan.StartDate != nil {
		auditPlanDTO.StartDateString = auditPlan.StartDate.Format("2006-01-02")
	}
	if auditPlan.TargetDate != nil {
		auditPlanDTO.TargetDateString = auditPlan.TargetDate.Format("2006-01-02")
	}

	auditPlanDTO.CreatedAt = auditPlan.CreatedAt
	auditPlanDTO.UpdatedAt = auditPlan.UpdatedAt

	return &auditPlanDTO, nil
}

