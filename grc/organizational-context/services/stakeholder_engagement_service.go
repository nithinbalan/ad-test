package services

import (
	"errors"
	"fmt"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	organizationalContextDTOs "jethur/grc/organizational-context/dtos"
	stakeholderEngagementRepositories "jethur/grc/organizational-context/repositories"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type StakeholderEngagementServiceInterface interface {
	GetAll(stakeholderID uint, params *organizationalContextDTOs.StakeholderEngagementListQueryParamsDTO) ([]*organizationalContextDTOs.StakeholderEngagementListDTO, error)
	GetFullDetailsByID(stakeholderID uint, id uint, tx ...utility.DatabaseTransaction) (*organizationalContextDTOs.StakeholderEngagementDetailDTO, error)
	GetByID(stakeholderID uint, id uint) (*models.StakeholderEngagement, error)
	Create(stakeholderID uint, currentUserID uint, organizationalContextDTOs *organizationalContextDTOs.StakeholderEngagementFormDTO) (uint, error)
	Update(stakeholderID uint, id uint, currentUserID uint, stakeholderEngagement *organizationalContextDTOs.StakeholderEngagementFormDTO) error
	Delete(stakeholderID uint, ids []uint, currentUserID uint) ([]uint, error)

	StatusUpdate(stakeholderID uint, id uint, currentUserID uint, stakeholderEngagementEvaluationFormDTO *organizationalContextDTOs.StakeholderEngagementEvaluationFormDTO, tx utility.DatabaseTransaction) error
}

// Service and implementations
type StakeholderEngagementService struct {
	repo           stakeholderEngagementRepositories.StakeholderEngagementRepositoryInterface
	activityLogSvc logServices.ActivityLogServiceInterface
	userSvc        organizationServices.UserServiceInterface
	userSettingSvc organizationServices.UserSettingServiceInterface
}

func NewStakeholderEngagementService(
	repo stakeholderEngagementRepositories.StakeholderEngagementRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	userSettingSvc organizationServices.UserSettingServiceInterface,
) StakeholderEngagementServiceInterface {
	return &StakeholderEngagementService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
		userSvc:        userSvc,
		userSettingSvc: userSettingSvc,
	}
}

func (ses *StakeholderEngagementService) GetAll(stakeholderID uint, params *organizationalContextDTOs.StakeholderEngagementListQueryParamsDTO) ([]*organizationalContextDTOs.StakeholderEngagementListDTO, error) {
	return ses.repo.GetAll(stakeholderID, params)
}

func (ses StakeholderEngagementService) GetByID(stakeholderID uint, id uint) (*models.StakeholderEngagement, error) {
	stakeholderEngagement, err := ses.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return stakeholderEngagement, nil
}

func (ses StakeholderEngagementService) GetFullDetailsByID(stakeholderID uint, id uint, tx ...utility.DatabaseTransaction) (*organizationalContextDTOs.StakeholderEngagementDetailDTO, error) {
	// Initialize your DTO
	var stakeholderEngagementDTO organizationalContextDTOs.StakeholderEngagementDetailDTO

	stakeholderEngagement, err := ses.repo.GetFullDetailsByID(id, tx...)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	// Map the result to your DTO
	err = mapstructure.Decode(stakeholderEngagement, &stakeholderEngagementDTO)
	if err != nil {
		return nil, err
	}

	stakeholderEngagementDTO.CreatedAt = stakeholderEngagement.CreatedAt
	stakeholderEngagementDTO.UpdatedAt = stakeholderEngagement.UpdatedAt

	return &stakeholderEngagementDTO, nil
}

func (ses *StakeholderEngagementService) Create(stakeholderID, currentUserID uint, stakeholderEngagementFormDTO *organizationalContextDTOs.StakeholderEngagementFormDTO) (uint, error) {

	tx := ses.repo.BeginTransaction()

	stakeholderEngagement := &models.StakeholderEngagement{
		StakeholderID:      stakeholderID,
		Title:              stakeholderEngagementFormDTO.Title,
		EngagementStrategy: stakeholderEngagementFormDTO.EngagementStrategy,
		AssignedOwnerID:    stakeholderEngagementFormDTO.AssignedOwnerID,
		ReviewFrequencyID:  stakeholderEngagementFormDTO.ReviewFrequencyID,
		Description:        stakeholderEngagementFormDTO.Description,
		CreatedBYID:        currentUserID,
	}
	id, err := ses.repo.Create(stakeholderEngagement, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	tx.Commit()
	go func() {
		existingStakeholderEngagement, err := ses.GetFullDetailsByID(stakeholderID, id)
		if err != nil {
			utility.LogError("Failed to fetch swot pestel question factor", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}
		err = ses.activityLogSvc.LogActivity(
			models.CreateAction,
			models.StakeholderEngagement{}.TableName(),
			id,
			existingStakeholderEngagement,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.StakeholderEngagement{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()
	return id, nil
}
func (ses *StakeholderEngagementService) Update(stakeholderID uint, id uint, currentUserID uint, stakeholderEngagementFormDTO *organizationalContextDTOs.StakeholderEngagementFormDTO) error {
	tx := ses.repo.BeginTransaction()

	existingStakeholderEngagement, err := ses.repo.GetByID(id)
	if err != nil {
		return utility.NewResourceNotFoundError(err, id)
	}

	oldStakeholderEngagement, err := ses.GetFullDetailsByID(stakeholderID, id, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	existingStakeholderEngagement.Title = stakeholderEngagementFormDTO.Title
	existingStakeholderEngagement.EngagementStrategy = stakeholderEngagementFormDTO.EngagementStrategy
	existingStakeholderEngagement.AssignedOwnerID = stakeholderEngagementFormDTO.AssignedOwnerID
	existingStakeholderEngagement.ReviewFrequencyID = stakeholderEngagementFormDTO.ReviewFrequencyID
	existingStakeholderEngagement.Description = stakeholderEngagementFormDTO.Description
	existingStakeholderEngagement.UpdatedBYID = &currentUserID

	err = ses.repo.Update(existingStakeholderEngagement, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	go func() {
		existingStakeholderEngagement, err := ses.GetFullDetailsByID(stakeholderID, id)
		if err != nil {
			return
		}

		// Prepare log data
		changedData := make(map[string]interface{})

		//log changes
		oldData := map[string]interface{}{
			"title":               oldStakeholderEngagement.Title,
			"engagement_strategy": oldStakeholderEngagement.EngagementStrategy,
			"assigned_owner":      utility.GetFieldValue(oldStakeholderEngagement.AssignedOwner),
			"review_frequency":    utility.GetFieldValue(oldStakeholderEngagement.ReviewFrequency),
			"description":         oldStakeholderEngagement.Description,
		}

		newData := map[string]interface{}{
			"title":               existingStakeholderEngagement.Title,
			"engagement_strategy": existingStakeholderEngagement.EngagementStrategy,
			"assigned_owner":      utility.GetFieldValue(existingStakeholderEngagement.AssignedOwner),
			"review_frequency":    utility.GetFieldValue(existingStakeholderEngagement.ReviewFrequency),
			"description":         existingStakeholderEngagement.Description,
		}
		// Compare other fields and add them to changedData
		fieldChanges := utility.CompareFieldChanges(oldData, newData)
		for k, v := range fieldChanges {
			changedData[k] = v
		}
		// Log activity if there are any changes
		if len(changedData) > 0 {
			err := ses.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.StakeholderEngagement{}.TableName(),
				id,
				changedData,
				currentUserID,
				nil,
			)

			if err != nil {
				fields := map[string]interface{}{
					"action":    models.UpdateAction,
					"tableName": models.StakeholderEngagement{}.TableName(),
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

func (ses *StakeholderEngagementService) Delete(stakeholderID uint, ids []uint, currentUserID uint) ([]uint, error) {
	tx := ses.repo.BeginTransaction()
	var failedIds []uint
	var successfulIds []uint
	var deletedStakeholderEngagement []map[string]interface{}

	for _, id := range ids {
		stakeholderEngagement, err := ses.GetFullDetailsByID(stakeholderID, id, tx)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		deletedStakeholderEngagement = append(deletedStakeholderEngagement, map[string]interface{}{
			"id":    id,
			"title": stakeholderEngagement.Title,
		})

		// Delete the business
		if err := ses.repo.Delete(id, tx); err != nil {
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
			err := ses.activityLogSvc.LogBulkActivities(
				models.DeleteAction, // Make sure this action constant is defined appropriately
				models.StakeholderEngagement{}.TableName(),
				successfulIds,
				deletedStakeholderEngagement,
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

func (ses *StakeholderEngagementService) StatusUpdate(stakeholderID uint, id uint, currentUserID uint, stakeholderEngagementEvaluationFormDTO *organizationalContextDTOs.StakeholderEngagementEvaluationFormDTO, tx utility.DatabaseTransaction) error {
	existingEngagement, err := ses.GetByID(stakeholderID, id)
	if err != nil {
		return err
	}

	dates, err := ses.parseEvaluationDates(stakeholderEngagementEvaluationFormDTO)
	if err != nil {
		return err
	}

	existingEngagement.Score = stakeholderEngagementEvaluationFormDTO.Score
	existingEngagement.EvaluatedBYID = &stakeholderEngagementEvaluationFormDTO.EvaluatedByID
	existingEngagement.EvaluatedAt = dates.EvaluatedAt
	existingEngagement.UpdatedBYID = &currentUserID

	err = ses.repo.Update(existingEngagement, tx)
	if err != nil {
		return err
	}

	return nil
}

func (ses *StakeholderEngagementService) parseEvaluationDates(dto *organizationalContextDTOs.StakeholderEngagementEvaluationFormDTO) (*organizationalContextDTOs.StakeholderEngagementEvaluationParsedDates, error) {
	parsed := &organizationalContextDTOs.StakeholderEngagementEvaluationParsedDates{}

	var err error
	if *dto.EvaluatedAt != "" {
		parsed.EvaluatedAt, err = utility.ParseDate(*dto.EvaluatedAt)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid evaluated date"), "invalid_date", "evaluated at")
		}
	}

	return parsed, nil
}
