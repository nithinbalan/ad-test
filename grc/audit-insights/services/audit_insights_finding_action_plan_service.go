package services

import (
	"errors"
	"jethur/grc/audit-insights/repositories"
	auditDTOs "jethur/grc/audit/dtos"
	auditFindingDTOs "jethur/grc/audit/dtos"
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
	Create(findingID uint, currentUserID uint, AuditInsightsFindingActionPlanFormDTO *auditDTOs.AuditFindingActionPlanFormDTO) (uint, error)
	Update(findingID uint, id uint, currentUserID uint, auditFinding *auditDTOs.AuditFindingActionPlanFormDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	GetAuditInsightsFindingActionPlanCountByStatus(findingId uint) ([]auditFindingDTOs.AuditFindingActionPlanStatusCountResponseDTO, error)
}

type AuditInsightsFindingActionPlanService struct {
	repo                      repositories.AuditInsightsFindingActionPlanRepositoryInterface
	activityLogSvc            logServices.ActivityLogServiceInterface
	userSvc                   organizationServices.UserServiceInterface
	auditFindingSvc           auditServices.AuditFindingServiceInterface
	auditFindingActionPlanSvc auditServices.AuditFindingActionPlanServiceInterface
}

func NewAuditInsightsFindingActionPlanService(
	repo repositories.AuditInsightsFindingActionPlanRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	auditFindingSvc auditServices.AuditFindingServiceInterface,
	auditFindingActionPlanSvc auditServices.AuditFindingActionPlanServiceInterface,
) AuditInsightsFindingActionPlanServiceInterface {
	return &AuditInsightsFindingActionPlanService{
		repo:                      repo,
		activityLogSvc:            activityLogSvc,
		userSvc:                   userSvc,
		auditFindingSvc:           auditFindingSvc,
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

func (aifaps *AuditInsightsFindingActionPlanService) Create(findingID uint, currentUserID uint, auditInsightsFindingActionPlanFormDTO *auditFindingDTOs.AuditFindingActionPlanFormDTO) (uint, error) {
	if err := aifaps.validateAuditInsightsFindingActionPlanDTO(findingID, currentUserID); err != nil {
		return 0, err
	}
	id, err := aifaps.auditFindingActionPlanSvc.Create(findingID, currentUserID, auditInsightsFindingActionPlanFormDTO)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (aifaps *AuditInsightsFindingActionPlanService) Update(findingID uint, id uint, currentUserID uint, auditInsightsFindingActionPlanFormDTO *auditFindingDTOs.AuditFindingActionPlanFormDTO) error {
	if err := aifaps.validateAuditInsightsFindingActionPlanDTO(findingID, currentUserID); err != nil {
		return err
	}
	err := aifaps.auditFindingActionPlanSvc.Update(findingID, id, currentUserID, auditInsightsFindingActionPlanFormDTO)
	if err != nil {
		return err
	}
	return nil
}

func (aifaps *AuditInsightsFindingActionPlanService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	failedIds, err := aifaps.auditFindingActionPlanSvc.Delete(ids, currentUserID)
	if err != nil {
		return nil, err
	}
	return failedIds, nil
}

func (aifaps *AuditInsightsFindingActionPlanService) validateAuditInsightsFindingActionPlanDTO(findingID uint, currentUserID uint) error {
	existingAuditFinding, err := aifaps.auditFindingSvc.GetFullDetailsByID(findingID)
	if err != nil {
		return utility.NewResourceNotFoundError(err, findingID)
	}
	orgUnitID, err := aifaps.userSvc.GetUserPrimaryUnitID(currentUserID)
	if err != nil {
		return err
	}

	if existingAuditFinding.AuditableItem.ResponsibleUnitID != orgUnitID {
		return utility.NewForbiddenError("Action not allowed:  Only users from the responsible unit can perform this action.")
	}
	return nil
}

func (aifaps *AuditInsightsFindingActionPlanService) GetAuditInsightsFindingActionPlanCountByStatus(findingId uint) ([]auditFindingDTOs.AuditFindingActionPlanStatusCountResponseDTO, error) {
	response, err := aifaps.auditFindingActionPlanSvc.GetAuditFindingActionPlanCountByStatus(findingId)
	if err != nil {
		return nil, err
	}

	return response, nil
}
