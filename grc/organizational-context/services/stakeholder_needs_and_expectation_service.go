package services

import (
	"errors"
	"fmt"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	organizationalContextDTOs "jethur/grc/organizational-context/dtos"
	stakeholderNeedsAndExpectationRepositories "jethur/grc/organizational-context/repositories"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type StakeholderNeedsAndExpectationServiceInterface interface {
	GetAll(stakeholderID uint, params *organizationalContextDTOs.StakeholderNeedsAndExpectationListQueryParamsDTO) ([]*organizationalContextDTOs.StakeholderNeedsAndExpectationListDTO, error)
	GetFullDetailsByID(stakeholderID uint, id uint, tx ...utility.DatabaseTransaction) (*organizationalContextDTOs.StakeholderNeedsAndExpectationDetailDTO, error)
	GetByID(stakeholderID uint, id uint) (*models.StakeholderNeedsAndExpectation, error)
	Create(stakeholderID uint, currentUserID uint, organizationalContextDTOs *organizationalContextDTOs.StakeholderNeedsAndExpectationFormDTO) (uint, error)
	Update(stakeholderID uint, id uint, currentUserID uint, stakeholderNeedsAndExpectation *organizationalContextDTOs.StakeholderNeedsAndExpectationFormDTO) error
	Delete(stakeholderID uint, ids []uint, currentUserID uint) ([]uint, error)

	StatusUpdate(stakeholderID uint, id uint, currentUserID uint, stakeholderNeedsAndExpectationEvaluationFormDTO *organizationalContextDTOs.StakeholderNeedsAndExpectationEvaluationFormDTO, tx utility.DatabaseTransaction) error
}

// Service and implementations
type StakeholderNeedsAndExpectationService struct {
	repo           stakeholderNeedsAndExpectationRepositories.StakeholderNeedsAndExpectationRepositoryInterface
	activityLogSvc logServices.ActivityLogServiceInterface
	userSvc        organizationServices.UserServiceInterface
	userSettingSvc organizationServices.UserSettingServiceInterface
}

func NewStakeholderNeedsAndExpectationService(
	repo stakeholderNeedsAndExpectationRepositories.StakeholderNeedsAndExpectationRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	userSettingSvc organizationServices.UserSettingServiceInterface,
) StakeholderNeedsAndExpectationServiceInterface {
	return &StakeholderNeedsAndExpectationService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
		userSvc:        userSvc,
		userSettingSvc: userSettingSvc,
	}
}

func (snes *StakeholderNeedsAndExpectationService) GetAll(stakeholderID uint, params *organizationalContextDTOs.StakeholderNeedsAndExpectationListQueryParamsDTO) ([]*organizationalContextDTOs.StakeholderNeedsAndExpectationListDTO, error) {
	return snes.repo.GetAll(stakeholderID, params)
}

func (snes StakeholderNeedsAndExpectationService) GetByID(stakeholderID uint, id uint) (*models.StakeholderNeedsAndExpectation, error) {
	stakeholderNeedsAndExpectation, err := snes.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return stakeholderNeedsAndExpectation, nil
}

func (snes StakeholderNeedsAndExpectationService) GetFullDetailsByID(stakeholderID uint, id uint, tx ...utility.DatabaseTransaction) (*organizationalContextDTOs.StakeholderNeedsAndExpectationDetailDTO, error) {
	// Initialize your DTO
	var stakeholderNeedsAndExpectationDTO organizationalContextDTOs.StakeholderNeedsAndExpectationDetailDTO

	stakeholderNeedsAndExpectation, err := snes.repo.GetFullDetailsByID(id, tx...)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	// Map the result to your DTO
	err = mapstructure.Decode(stakeholderNeedsAndExpectation, &stakeholderNeedsAndExpectationDTO)
	if err != nil {
		return nil, err
	}

	stakeholderNeedsAndExpectationDTO.CreatedAt = stakeholderNeedsAndExpectation.CreatedAt
	stakeholderNeedsAndExpectationDTO.UpdatedAt = stakeholderNeedsAndExpectation.UpdatedAt

	return &stakeholderNeedsAndExpectationDTO, nil
}

func (snes *StakeholderNeedsAndExpectationService) Create(stakeholderID, currentUserID uint, stakeholderNeedsAndExpectationFormDTO *organizationalContextDTOs.StakeholderNeedsAndExpectationFormDTO) (uint, error) {

	tx := snes.repo.BeginTransaction()

	stakeholderNeedsAndExpectation := &models.StakeholderNeedsAndExpectation{
		StakeholderID:   stakeholderID,
		Title:           stakeholderNeedsAndExpectationFormDTO.Title,
		SeverityLevelID: stakeholderNeedsAndExpectationFormDTO.SeverityLevelID,
		SuccessCriteria: stakeholderNeedsAndExpectationFormDTO.SuccessCriteria,
		Challenges:      stakeholderNeedsAndExpectationFormDTO.Challenges,
		Description:     stakeholderNeedsAndExpectationFormDTO.Description,
		CreatedBYID:     currentUserID,
	}
	id, err := snes.repo.Create(stakeholderNeedsAndExpectation, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	tx.Commit()
	go func() {
		existingStakeholderNeedsAndExpectation, err := snes.GetFullDetailsByID(stakeholderID, id)
		if err != nil {
			utility.LogError("Failed to fetch swot pestel question factor", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}
		err = snes.activityLogSvc.LogActivity(
			models.CreateAction,
			models.StakeholderNeedsAndExpectation{}.TableName(),
			id,
			existingStakeholderNeedsAndExpectation,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.StakeholderNeedsAndExpectation{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()
	return id, nil
}
func (snes *StakeholderNeedsAndExpectationService) Update(stakeholderID uint, id uint, currentUserID uint, stakeholderNeedsAndExpectationFormDTO *organizationalContextDTOs.StakeholderNeedsAndExpectationFormDTO) error {
	tx := snes.repo.BeginTransaction()

	existingStakeholderNeedsAndExpectation, err := snes.repo.GetByID(id)
	if err != nil {
		return utility.NewResourceNotFoundError(err, id)
	}

	oldStakeholderNeedsAndExpectation, err := snes.GetFullDetailsByID(stakeholderID, id, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	existingStakeholderNeedsAndExpectation.Title = stakeholderNeedsAndExpectationFormDTO.Title
	existingStakeholderNeedsAndExpectation.SeverityLevelID = stakeholderNeedsAndExpectationFormDTO.SeverityLevelID
	existingStakeholderNeedsAndExpectation.SuccessCriteria = stakeholderNeedsAndExpectationFormDTO.SuccessCriteria
	existingStakeholderNeedsAndExpectation.Challenges = stakeholderNeedsAndExpectationFormDTO.Challenges
	existingStakeholderNeedsAndExpectation.Description = stakeholderNeedsAndExpectationFormDTO.Description
	existingStakeholderNeedsAndExpectation.UpdatedBYID = &currentUserID

	err = snes.repo.Update(existingStakeholderNeedsAndExpectation, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	go func() {
		existingStakeholderNeedsAndExpectation, err := snes.GetFullDetailsByID(stakeholderID, id)
		if err != nil {
			return
		}

		// Prepare log data
		changedData := make(map[string]interface{})

		//log changes
		oldData := map[string]interface{}{
			"title":             oldStakeholderNeedsAndExpectation.Title,
			"severity_level_id": oldStakeholderNeedsAndExpectation.SeverityLevel,
			"success_criteria":  oldStakeholderNeedsAndExpectation.SuccessCriteria,
			"challenges":        oldStakeholderNeedsAndExpectation.Challenges,
			"description":       oldStakeholderNeedsAndExpectation.Description,
		}

		newData := map[string]interface{}{
			"title":             existingStakeholderNeedsAndExpectation.Title,
			"severity_level_id": existingStakeholderNeedsAndExpectation.SeverityLevel,
			"success_criteria":  existingStakeholderNeedsAndExpectation.SuccessCriteria,
			"challenges":        existingStakeholderNeedsAndExpectation.Challenges,
			"description":       existingStakeholderNeedsAndExpectation.Description,
		}
		// Compare other fields and add them to changedData
		fieldChanges := utility.CompareFieldChanges(oldData, newData)
		for k, v := range fieldChanges {
			changedData[k] = v
		}
		// Log activity if there are any changes
		if len(changedData) > 0 {
			err := snes.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.StakeholderNeedsAndExpectation{}.TableName(),
				id,
				changedData,
				currentUserID,
				nil,
			)

			if err != nil {
				fields := map[string]interface{}{
					"action":    models.UpdateAction,
					"tableName": models.StakeholderNeedsAndExpectation{}.TableName(),
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

func (snes *StakeholderNeedsAndExpectationService) Delete(stakeholderID uint, ids []uint, currentUserID uint) ([]uint, error) {
	tx := snes.repo.BeginTransaction()
	var failedIds []uint
	var successfulIds []uint
	var deletedStakeholderNeedsAndExpectation []map[string]interface{}

	for _, id := range ids {
		stakeholderNeedsAndExpectation, err := snes.GetFullDetailsByID(stakeholderID, id, tx)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		deletedStakeholderNeedsAndExpectation = append(deletedStakeholderNeedsAndExpectation, map[string]interface{}{
			"id":    id,
			"title": stakeholderNeedsAndExpectation.Title,
		})

		// Delete the business
		if err := snes.repo.Delete(id, tx); err != nil {
			failedIds = append(failedIds, id)
			continue
		} else {
			successfulIds = append(successfulIds, id)
		}
	}
	if len(failedIds) > 0 {
		tx.Rollback()
		return failedIds, nil
	}

	tx.Commit()

	// Log bulk delete activities for successful deletions
	if len(successfulIds) > 0 {
		go func() {
			err := snes.activityLogSvc.LogBulkActivities(
				models.DeleteAction, // Make sure this action constant is defined appropriately
				models.StakeholderNeedsAndExpectation{}.TableName(),
				successfulIds,
				deletedStakeholderNeedsAndExpectation,
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

func (snes *StakeholderNeedsAndExpectationService) StatusUpdate(stakeholderID uint, id uint, currentUserID uint, stakeholderNeedsAndExpectationEvaluationFormDTO *organizationalContextDTOs.StakeholderNeedsAndExpectationEvaluationFormDTO, tx utility.DatabaseTransaction) error {
	existingNeedsAndExpectation, err := snes.GetByID(stakeholderID, id)
	if err != nil {
		return err
	}

	dates, err := snes.parseEvaluationDates(stakeholderNeedsAndExpectationEvaluationFormDTO)
	if err != nil {
		return err
	}

	existingNeedsAndExpectation.EvaluationStatusID = &stakeholderNeedsAndExpectationEvaluationFormDTO.EvaluationStatusID
	existingNeedsAndExpectation.Score = stakeholderNeedsAndExpectationEvaluationFormDTO.Score
	existingNeedsAndExpectation.EvaluatedBYID = &stakeholderNeedsAndExpectationEvaluationFormDTO.EvaluatedByID
	existingNeedsAndExpectation.EvaluatedAt = dates.EvaluatedAt
	existingNeedsAndExpectation.UpdatedBYID = &currentUserID

	err = snes.repo.Update(existingNeedsAndExpectation, tx)
	if err != nil {
		return err
	}

	return nil
}

func (snes *StakeholderNeedsAndExpectationService) parseEvaluationDates(dto *organizationalContextDTOs.StakeholderNeedsAndExpectationEvaluationFormDTO) (*organizationalContextDTOs.StakeholderNeedsAndExpectationEvaluationParsedDates, error) {
	parsed := &organizationalContextDTOs.StakeholderNeedsAndExpectationEvaluationParsedDates{}

	var err error
	if *dto.EvaluatedAt != "" {
		parsed.EvaluatedAt, err = utility.ParseDate(*dto.EvaluatedAt)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid evaluated date"), "invalid_date", "evaluated at")
		}
	}

	return parsed, nil
}
