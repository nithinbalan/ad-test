package services

import (
	"errors"
	"fmt"
	auditFindingDTOs "jethur/grc/audit/dtos"
	auditFindingRepositories "jethur/grc/audit/repositories"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type AuditDraftFindingActionPlanServiceInterface interface {
	GetAll(findingID uint, params *auditFindingDTOs.AuditDraftFindingActionPlanListQueryParamsDTO) (*auditFindingDTOs.AuditDraftFindingActionPlanListResponseDTO, error)
	GetByID(findingID uint, id uint) (*models.AuditDraftFindingActionPlan, error)
	GetFullDetailsByID(findingID uint, id uint) (*auditFindingDTOs.AuditDraftFindingActionPlanDetailDTO, error)
	Create(findingID uint, currentUserID uint, AuditDraftFindingActionPlanFormDTO *auditFindingDTOs.AuditDraftFindingActionPlanFormDTO) (uint, error)
	Update(findingID uint, id uint, currentUserID uint, auditFinding *auditFindingDTOs.AuditDraftFindingActionPlanFormDTO) error
	Delete(findingID uint, ids []uint, currentUserID uint) ([]uint, error)
}

// Service and implementations
type AuditDraftFindingActionPlanService struct {
	repo                       auditFindingRepositories.AuditDraftFindingActionPlanRepositoryInterface
	activityLogSvc             logServices.ActivityLogServiceInterface
	auditDraftFindingSvc       AuditDraftFindingServiceInterface
	auditDraftFindingStatusSvc AuditDraftFindingStatusServiceInterface
}

func NewAuditDraftFindingActionPlanService(
	repo auditFindingRepositories.AuditDraftFindingActionPlanRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	auditDraftFindingSvc AuditDraftFindingServiceInterface,
	auditDraftFindingStatusSvc AuditDraftFindingStatusServiceInterface,
) AuditDraftFindingActionPlanServiceInterface {
	return &AuditDraftFindingActionPlanService{
		repo:                       repo,
		activityLogSvc:             activityLogSvc,
		auditDraftFindingSvc:       auditDraftFindingSvc,
		auditDraftFindingStatusSvc: auditDraftFindingStatusSvc,
	}
}

func (adfaps *AuditDraftFindingActionPlanService) GetAll(findingID uint, params *auditFindingDTOs.AuditDraftFindingActionPlanListQueryParamsDTO) (*auditFindingDTOs.AuditDraftFindingActionPlanListResponseDTO, error) {
	return adfaps.repo.GetAll(findingID, params)
}

func (adfaps AuditDraftFindingActionPlanService) GetByID(findingID uint, id uint) (*models.AuditDraftFindingActionPlan, error) {
	if _, err := adfaps.auditDraftFindingSvc.GetByID(findingID); err != nil {
		return nil, err
	}

	auditDraftFindingActionPlan, err := adfaps.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return auditDraftFindingActionPlan, nil
}

func (adfaps AuditDraftFindingActionPlanService) GetFullDetailsByID(findingID uint, id uint) (*auditFindingDTOs.AuditDraftFindingActionPlanDetailDTO, error) {
	var auditDraftFindingActionPlanDTO auditFindingDTOs.AuditDraftFindingActionPlanDetailDTO

	if _, err := adfaps.auditDraftFindingSvc.GetByID(findingID); err != nil {
		return nil, err
	}

	auditDraftFindingActionPlan, err := adfaps.repo.GetFullDetailsByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	err = mapstructure.Decode(auditDraftFindingActionPlan, &auditDraftFindingActionPlanDTO)
	if err != nil {
		return nil, err
	}

	// todo :- also add extra fields which required in details screen
	if auditDraftFindingActionPlan.TargetDate != nil {
		targetDate := auditDraftFindingActionPlan.TargetDate.Format("2006-01-02")
		auditDraftFindingActionPlanDTO.TargetDateString = &targetDate
	}

	auditDraftFindingActionPlanDTO.CreatedAt = auditDraftFindingActionPlan.CreatedAt
	auditDraftFindingActionPlanDTO.UpdatedAt = auditDraftFindingActionPlan.UpdatedAt

	return &auditDraftFindingActionPlanDTO, nil
}

func (adfaps *AuditDraftFindingActionPlanService) Create(findingID uint, currentUserID uint, AuditDraftFindingActionPlanFormDTO *auditFindingDTOs.AuditDraftFindingActionPlanFormDTO) (uint, error) {
	tx := adfaps.repo.BeginTransaction()
	if _, err := adfaps.auditDraftFindingSvc.GetByID(findingID); err != nil {
		return 0, err
	}

	dates, err := adfaps.parseAuditDraftFindingActionPlanDates(AuditDraftFindingActionPlanFormDTO)
	if err != nil {
		return 0, err
	}

	auditDraftFindingActionPlan := &models.AuditDraftFindingActionPlan{
		AuditDraftFindingID: findingID,
		Title:               AuditDraftFindingActionPlanFormDTO.Title,
		Description:         *AuditDraftFindingActionPlanFormDTO.Description,
		ResponsibleUserID:   &AuditDraftFindingActionPlanFormDTO.ResponsibleUserID,
		TargetDate:          dates.TargetDate,
		CreatedBYID:         currentUserID,
	}

	id, err := adfaps.repo.Create(auditDraftFindingActionPlan, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	tx.Commit()

	go func() {
		err := adfaps.activityLogSvc.LogActivity(
			models.CreateAction,
			models.AuditDraftFindingActionPlan{}.TableName(),
			id,
			AuditDraftFindingActionPlanFormDTO,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.AuditDraftFindingActionPlan{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return id, nil
}

func (adfaps *AuditDraftFindingActionPlanService) Update(findingID uint, id uint, currentUserID uint, AuditDraftFindingActionPlanFormDTO *auditFindingDTOs.AuditDraftFindingActionPlanFormDTO) error {
	tx := adfaps.repo.BeginTransaction()

	if _, err := adfaps.auditDraftFindingSvc.GetByID(findingID); err != nil {
		return err
	}

	existingAuditDraftFindingActionPlan, err := adfaps.GetByID(findingID, id)

	if err != nil {
		return err
	}

	dates, err := adfaps.parseAuditDraftFindingActionPlanDates(AuditDraftFindingActionPlanFormDTO)
	if err != nil {
		return err
	}

	// Update Audit Finding ActionPlan
	existingAuditDraftFindingActionPlan.Title = AuditDraftFindingActionPlanFormDTO.Title
	existingAuditDraftFindingActionPlan.Description = *AuditDraftFindingActionPlanFormDTO.Description
	existingAuditDraftFindingActionPlan.ResponsibleUserID = &AuditDraftFindingActionPlanFormDTO.ResponsibleUserID
	existingAuditDraftFindingActionPlan.TargetDate = dates.TargetDate
	existingAuditDraftFindingActionPlan.UpdatedBYID = &currentUserID

	err = adfaps.repo.Update(existingAuditDraftFindingActionPlan, tx)

	if err != nil {
		return err
	}

	tx.Commit()

	return nil
}

func (adfaps *AuditDraftFindingActionPlanService) Delete(findingID uint, ids []uint, currentUserID uint) ([]uint, error) {
	tx := adfaps.repo.BeginTransaction()

	if _, err := adfaps.auditDraftFindingSvc.GetByID(findingID); err != nil {
		return nil, err
	}

	var failedIds []uint
	var successfulIds []uint
	for _, id := range ids {

		_, err := adfaps.GetFullDetailsByID(findingID, id)

		if err != nil {
			continue
		}

		err = adfaps.repo.Delete(id, tx)
		if err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}
	}
	tx.Commit()

	if len(successfulIds) > 0 {

		go func() {
			err := adfaps.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.AuditFinding{}.TableName(),
				successfulIds,
				nil,
				currentUserID,
				nil,
			)
			if err != nil {
				// @todo - save to log file
				fmt.Printf("Failed to log bulk delete activity: %v\n", err)
			}
		}()
	}

	return failedIds, nil
}

func (adfaps *AuditDraftFindingActionPlanService) parseAuditDraftFindingActionPlanDates(dto *auditFindingDTOs.AuditDraftFindingActionPlanFormDTO) (*auditFindingDTOs.AuditDraftFindingActionPlanParsedDates, error) {
	parsed := &auditFindingDTOs.AuditDraftFindingActionPlanParsedDates{}

	var err error
	if *dto.TargetDate != "" {
		parsed.TargetDate, err = utility.ParseDate(*dto.TargetDate)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid target date"), "invalid_date", "target date")
		}
	}

	return parsed, nil
}
