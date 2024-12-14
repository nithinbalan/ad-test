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
	Create(findingID uint, currentUserID uint, AuditDraftFindingActionPlanFormDTO *auditFindingDTOs.AuditDraftFindingActionPlanFormDTO, tx ...utility.DatabaseTransaction) (uint, error)
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

func (adfaps *AuditDraftFindingActionPlanService) Create(findingID uint, currentUserID uint, AuditDraftFindingActionPlanFormDTO *auditFindingDTOs.AuditDraftFindingActionPlanFormDTO, tx ...utility.DatabaseTransaction) (uint, error) {
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

	id, err := adfaps.repo.Create(auditDraftFindingActionPlan, tx...)
	if err != nil {
		return 0, err
	}

	go func() {

		existingAuditDraftFindingActionPlan, err := adfaps.GetFullDetailsByID(findingID, id)
		if err != nil {
			utility.LogError("Failed to fetch audit draft finding", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		err = adfaps.activityLogSvc.LogActivity(
			models.CreateAction,
			models.AuditDraftFindingActionPlan{}.TableName(),
			id,
			existingAuditDraftFindingActionPlan,
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

	oldAuditDraftFindingActionPlan, err := adfaps.repo.GetFullDetailsByID(id)
	if err != nil {
		tx.Rollback()
		return utility.NewResourceNotFoundError(err, id)
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
	go func() {
		existingAuditDraftFindingActionPlan, err := adfaps.repo.GetFullDetailsByID(id)
		if err != nil {
			return
		}
		// Prepare log data
		changedData := make(map[string]interface{})
		//log changes
		oldData := map[string]interface{}{
			"title":            oldAuditDraftFindingActionPlan.Title,
			"description":      oldAuditDraftFindingActionPlan.Description,
			"responsible_user": fmt.Sprintf("%s %s", oldAuditDraftFindingActionPlan.ResponsibleUser.FirstNameEn, oldAuditDraftFindingActionPlan.ResponsibleUser.LastNameEn),
			"target_date":      oldAuditDraftFindingActionPlan.TargetDate,
		}

		newData := map[string]interface{}{
			"title":            existingAuditDraftFindingActionPlan.Title,
			"description":      existingAuditDraftFindingActionPlan.Description,
			"responsible_user": fmt.Sprintf("%s %s", existingAuditDraftFindingActionPlan.ResponsibleUser.FirstNameEn, existingAuditDraftFindingActionPlan.ResponsibleUser.LastNameEn),
			"target_date":      existingAuditDraftFindingActionPlan.TargetDate,
		}
		// Compare other fields and add them to changedData
		fieldChanges := utility.CompareFieldChanges(oldData, newData)
		for k, v := range fieldChanges {
			changedData[k] = v
		}
		// Log activity if there are any changes
		if len(changedData) > 0 {
			err := adfaps.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.AuditDraftFindingActionPlan{}.TableName(),
				id,
				changedData,
				currentUserID,
				nil,
			)

			if err != nil {
				fields := map[string]interface{}{
					"action":    models.UpdateAction,
					"tableName": models.AuditDraftFindingActionPlan{}.TableName(),
					"recordID":  id,
					"userID":    currentUserID,
					"error":     err.Error(),
				}

				utility.LogError("Failed to log activity", fields)
			}
		}
	}()

	return nil
}

func (adfaps *AuditDraftFindingActionPlanService) Delete(findingID uint, ids []uint, currentUserID uint) ([]uint, error) {
	tx := adfaps.repo.BeginTransaction()

	if _, err := adfaps.auditDraftFindingSvc.GetByID(findingID); err != nil {
		return nil, err
	}

	var failedIds []uint
	var successfulIds []uint
	var deletedAuditDraftFindingActionPlan []map[string]interface{}
	for _, id := range ids {

		auditDraftFindingActionPlan, err := adfaps.GetFullDetailsByID(findingID, id)

		if err != nil {
			continue
		}

		err = adfaps.repo.Delete(id, tx)
		if err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}
		deletedAuditDraftFindingActionPlan = append(deletedAuditDraftFindingActionPlan, map[string]interface{}{
			"id":                  id,
			"audit_draft_finding": auditDraftFindingActionPlan.AuditDraftFinding.Title,
			"auditable_item":      auditDraftFindingActionPlan.Title,
		})
	}
	tx.Commit()

	if len(successfulIds) > 0 {

		go func() {
			err := adfaps.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.AuditFinding{}.TableName(),
				successfulIds,
				deletedAuditDraftFindingActionPlan,
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
