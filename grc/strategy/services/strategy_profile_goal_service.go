package services

import (
	"errors"
	"fmt"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	strategyDTOs "jethur/grc/strategy/dtos"
	strategyProfileGoalRepositories "jethur/grc/strategy/repositories"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type StrategyProfileGoalServiceInterface interface {
	GetAll(strategyProfileID uint, params *strategyDTOs.StrategyProfileGoalListQueryParamsDTO) ([]*strategyDTOs.StrategyProfileGoalListDTO, error)
	GetFullDetailsByID(strategyProfileID uint, id uint, tx ...utility.DatabaseTransaction) (*strategyDTOs.StrategyProfileGoalDetailDTO, error)
	GetByID(id uint) (*models.StrategyProfileGoal, error)
	Create(strategyProfileID uint, currentUserID uint, strategyDTOs *strategyDTOs.StrategyProfileGoalFormDTO) (uint, error)
	Update(strategyProfileID uint, id uint, currentUserID uint, strategyProfileGoal *strategyDTOs.StrategyProfileGoalFormDTO) error
	Delete(strategyProfileID uint, ids []uint, currentUserID uint) ([]uint, error)
}

// Service and implementations
type StrategyProfileGoalService struct {
	repo               strategyProfileGoalRepositories.StrategyProfileGoalRepositoryInterface
	activityLogSvc     logServices.ActivityLogServiceInterface
	userSvc            organizationServices.UserServiceInterface
	userSettingSvc     organizationServices.UserSettingServiceInterface
	strategyProfileSvc StrategyProfileServiceInterface
}

func NewStrategyProfileGoalService(
	repo strategyProfileGoalRepositories.StrategyProfileGoalRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	userSettingSvc organizationServices.UserSettingServiceInterface,
	strategyProfileSvc StrategyProfileServiceInterface,
) StrategyProfileGoalServiceInterface {
	return &StrategyProfileGoalService{
		repo:               repo,
		activityLogSvc:     activityLogSvc,
		userSvc:            userSvc,
		userSettingSvc:     userSettingSvc,
		strategyProfileSvc: strategyProfileSvc,
	}
}

func (sps *StrategyProfileGoalService) GetAll(strategyProfileID uint, params *strategyDTOs.StrategyProfileGoalListQueryParamsDTO) ([]*strategyDTOs.StrategyProfileGoalListDTO, error) {
	return sps.repo.GetAll(strategyProfileID, params)
}

func (sps StrategyProfileGoalService) GetByID(id uint) (*models.StrategyProfileGoal, error) {
	strategyProfileGoal, err := sps.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return strategyProfileGoal, nil
}

func (sps StrategyProfileGoalService) GetFullDetailsByID(strategyProfileID uint, id uint, tx ...utility.DatabaseTransaction) (*strategyDTOs.StrategyProfileGoalDetailDTO, error) {
	// Initialize your DTO
	var strategyProfileGoalDTO strategyDTOs.StrategyProfileGoalDetailDTO

	if _, err := sps.strategyProfileSvc.GetByID(strategyProfileID); err != nil {
		return nil, err
	}

	strategyProfileGoal, err := sps.repo.GetFullDetailsByID(id, tx...)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	// Map the result to your DTO
	err = mapstructure.Decode(strategyProfileGoal, &strategyProfileGoalDTO)
	if err != nil {
		return nil, err
	}
	if strategyProfileGoal.TargetDate != nil {
		strategyProfileGoalDTO.TargetDateString = strategyProfileGoal.TargetDate.Format("2006-01-02")
	}

	strategyProfileGoalDTO.StrategyProfile.StartDateString = strategyProfileGoal.StrategyProfile.StartDate.Format("2006-01-02")
	strategyProfileGoalDTO.StrategyProfile.EndDateString = strategyProfileGoal.StrategyProfile.EndDate.Format("2006-01-02")

	strategyProfileGoalDTO.CreatedAt = strategyProfileGoal.CreatedAt
	strategyProfileGoalDTO.UpdatedAt = strategyProfileGoal.UpdatedAt

	return &strategyProfileGoalDTO, nil
}

func (sps *StrategyProfileGoalService) Create(strategyProfileID uint, currentUserID uint, strategyProfileGoalFormDTO *strategyDTOs.StrategyProfileGoalFormDTO) (uint, error) {

	if _, err := sps.strategyProfileSvc.GetByID(strategyProfileID); err != nil {
		return 0, err
	}

	tx := sps.repo.BeginTransaction()

	dates, err := sps.parseStrategyProfileDates(strategyProfileGoalFormDTO)
	if err != nil {
		return 0, err
	}

	strategyProfileGoal := &models.StrategyProfileGoal{
		StrategyProfileID:          strategyProfileID,
		StrategyProfileFocusAreaID: strategyProfileGoalFormDTO.StrategyProfileFocusAreaID,
		Title:                      strategyProfileGoalFormDTO.Title,
		Description:                strategyProfileGoalFormDTO.Description,
		TargetDate:                 dates.TargetDate,
		ResponsibleUserID:          strategyProfileGoalFormDTO.ResponsibleUserID,
		CreatedBYID:                currentUserID,
	}
	id, err := sps.repo.Create(strategyProfileGoal, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	tx.Commit()
	go func() {
		existingStrategyProfileGoal, err := sps.GetFullDetailsByID(strategyProfileID, id)
		if err != nil {
			utility.LogError("Failed to fetch strategy profile goal", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}
		err = sps.activityLogSvc.LogActivity(
			models.CreateAction,
			models.StrategyProfileGoal{}.TableName(),
			id,
			existingStrategyProfileGoal,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.StrategyProfileGoal{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()
	return id, nil
}
func (sps *StrategyProfileGoalService) Update(strategyProfileID uint, id uint, currentUserID uint, strategyProfileGoalFormDTO *strategyDTOs.StrategyProfileGoalFormDTO) error {
	if _, err := sps.strategyProfileSvc.GetByID(strategyProfileID); err != nil {
		return err
	}

	tx := sps.repo.BeginTransaction()

	dates, err := sps.parseStrategyProfileDates(strategyProfileGoalFormDTO)
	if err != nil {
		return err
	}

	existingStrategyProfileGoal, err := sps.repo.GetByID(id)
	if err != nil {
		return utility.NewResourceNotFoundError(err, id)
	}

	oldStrategyProfileGoal, err := sps.GetFullDetailsByID(strategyProfileID, id, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	existingStrategyProfileGoal.Title = strategyProfileGoalFormDTO.Title
	existingStrategyProfileGoal.Description = strategyProfileGoalFormDTO.Description
	existingStrategyProfileGoal.TargetDate = dates.TargetDate
	existingStrategyProfileGoal.ResponsibleUserID = strategyProfileGoalFormDTO.ResponsibleUserID
	existingStrategyProfileGoal.UpdatedBYID = &currentUserID

	err = sps.repo.Update(existingStrategyProfileGoal, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	go func() {
		existingStrategyProfileGoal, err := sps.GetFullDetailsByID(strategyProfileID, id)
		if err != nil {
			return
		}

		// Prepare log data
		changedData := make(map[string]interface{})

		//log changes
		oldData := map[string]interface{}{
			"title":            oldStrategyProfileGoal.Title,
			"description":      oldStrategyProfileGoal.Description,
			"target_date":      oldStrategyProfileGoal.TargetDateString,
			"responsible_user": utility.GetFieldValue(oldStrategyProfileGoal.ResponsibleUser),
		}

		newData := map[string]interface{}{
			"title":            existingStrategyProfileGoal.Title,
			"description":      existingStrategyProfileGoal.Description,
			"target_date":      existingStrategyProfileGoal.TargetDateString,
			"responsible_user": utility.GetFieldValue(existingStrategyProfileGoal.ResponsibleUser),
		}
		// Compare other fields and add them to changedData
		fieldChanges := utility.CompareFieldChanges(oldData, newData)
		for k, v := range fieldChanges {
			changedData[k] = v
		}
		// Log activity if there are any changes
		if len(changedData) > 0 {
			err := sps.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.StrategyProfileGoal{}.TableName(),
				id,
				changedData,
				currentUserID,
				nil,
			)

			if err != nil {
				fields := map[string]interface{}{
					"action":    models.UpdateAction,
					"tableName": models.StrategyProfileGoal{}.TableName(),
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

func (sps *StrategyProfileGoalService) Delete(strategyProfileID uint, ids []uint, currentUserID uint) ([]uint, error) {
	if _, err := sps.strategyProfileSvc.GetByID(strategyProfileID); err != nil {
		return nil, err
	}

	tx := sps.repo.BeginTransaction()
	var failedIds []uint
	var successfulIds []uint
	var deletedStrategyProfileGoal []map[string]interface{}

	for _, id := range ids {
		strategyProfileGoal, err := sps.GetFullDetailsByID(strategyProfileID, id, tx)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		deletedStrategyProfileGoal = append(deletedStrategyProfileGoal, map[string]interface{}{
			"id":       id,
			"title_en": strategyProfileGoal.Title,
		})

		// Delete the business
		if err := sps.repo.Delete(id, tx); err != nil {
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
			err := sps.activityLogSvc.LogBulkActivities(
				models.DeleteAction, // Make sure this action constant is defined appropriately
				models.StrategyProfileGoal{}.TableName(),
				successfulIds,
				deletedStrategyProfileGoal,
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

func (sps *StrategyProfileGoalService) parseStrategyProfileDates(dto *strategyDTOs.StrategyProfileGoalFormDTO) (*strategyDTOs.StrategyProfileGoalParsedDates, error) {
	parsed := &strategyDTOs.StrategyProfileGoalParsedDates{}

	var err error
	if dto.TargetDate != "" {
		parsed.TargetDate, err = utility.ParseDate(dto.TargetDate)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid start date"), "invalid_date", "start date")
		}
	}

	return parsed, nil
}
