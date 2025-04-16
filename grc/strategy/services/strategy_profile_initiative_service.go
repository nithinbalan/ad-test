package services

import (
	"errors"
	"fmt"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	strategyDTOs "jethur/grc/strategy/dtos"
	strategyProfileInitiativeRepositories "jethur/grc/strategy/repositories"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type StrategyProfileInitiativeServiceInterface interface {
	GetAll(strategyProfileID uint, params *strategyDTOs.StrategyProfileInitiativeListQueryParamsDTO) ([]*strategyDTOs.StrategyProfileInitiativeListDTO, error)
	GetFullDetailsByID(strategyProfileID uint, id uint, tx ...utility.DatabaseTransaction) (*strategyDTOs.StrategyProfileInitiativeDetailDTO, error)
	GetByID(id uint) (*models.StrategyProfileInitiative, error)
	Create(strategyProfileID uint, currentUserID uint, strategyDTOs *strategyDTOs.StrategyProfileInitiativeFormDTO) (uint, error)
	Update(strategyProfileID uint, id uint, currentUserID uint, strategyProfileInitiative *strategyDTOs.StrategyProfileInitiativeFormDTO) error
	Delete(strategyProfileID uint, ids []uint, currentUserID uint) ([]uint, error)
}

// Service and implementations
type StrategyProfileInitiativeService struct {
	repo               strategyProfileInitiativeRepositories.StrategyProfileInitiativeRepositoryInterface
	activityLogSvc     logServices.ActivityLogServiceInterface
	userSvc            organizationServices.UserServiceInterface
	userSettingSvc     organizationServices.UserSettingServiceInterface
	strategyProfileSvc StrategyProfileServiceInterface
}

func NewStrategyProfileInitiativeService(
	repo strategyProfileInitiativeRepositories.StrategyProfileInitiativeRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	userSettingSvc organizationServices.UserSettingServiceInterface,
	strategyProfileSvc StrategyProfileServiceInterface,
) StrategyProfileInitiativeServiceInterface {
	return &StrategyProfileInitiativeService{
		repo:               repo,
		activityLogSvc:     activityLogSvc,
		userSvc:            userSvc,
		userSettingSvc:     userSettingSvc,
		strategyProfileSvc: strategyProfileSvc,
	}
}

func (sps *StrategyProfileInitiativeService) GetAll(strategyProfileID uint, params *strategyDTOs.StrategyProfileInitiativeListQueryParamsDTO) ([]*strategyDTOs.StrategyProfileInitiativeListDTO, error) {
	return sps.repo.GetAll(strategyProfileID, params)
}

func (sps StrategyProfileInitiativeService) GetByID(id uint) (*models.StrategyProfileInitiative, error) {
	strategyProfileInitiative, err := sps.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return strategyProfileInitiative, nil
}

func (sps StrategyProfileInitiativeService) GetFullDetailsByID(strategyProfileID uint, id uint, tx ...utility.DatabaseTransaction) (*strategyDTOs.StrategyProfileInitiativeDetailDTO, error) {
	// Initialize your DTO
	var strategyProfileInitiativeDTO strategyDTOs.StrategyProfileInitiativeDetailDTO

	if _, err := sps.strategyProfileSvc.GetByID(strategyProfileID); err != nil {
		return nil, err
	}

	strategyProfileInitiative, err := sps.repo.GetFullDetailsByID(id, tx...)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	// Map the result to your DTO
	err = mapstructure.Decode(strategyProfileInitiative, &strategyProfileInitiativeDTO)
	if err != nil {
		return nil, err
	}

	if strategyProfileInitiative.StartDate != nil {
		strategyProfileInitiativeDTO.StartDateString = strategyProfileInitiative.StartDate.Format("2006-01-02")
	}

	if strategyProfileInitiative.EndDate != nil {
		strategyProfileInitiativeDTO.EndDateString = strategyProfileInitiative.EndDate.Format("2006-01-02")
	}

	strategyProfileInitiativeDTO.StrategyProfile.StartDateString = strategyProfileInitiative.StrategyProfile.StartDate.Format("2006-01-02")
	strategyProfileInitiativeDTO.StrategyProfile.EndDateString = strategyProfileInitiative.StrategyProfile.EndDate.Format("2006-01-02")

	strategyProfileInitiativeDTO.CreatedAt = strategyProfileInitiative.CreatedAt
	strategyProfileInitiativeDTO.UpdatedAt = strategyProfileInitiative.UpdatedAt

	return &strategyProfileInitiativeDTO, nil
}

func (sps *StrategyProfileInitiativeService) Create(strategyProfileID uint, currentUserID uint, strategyProfileInitiativeFormDTO *strategyDTOs.StrategyProfileInitiativeFormDTO) (uint, error) {

	if _, err := sps.strategyProfileSvc.GetByID(strategyProfileID); err != nil {
		return 0, err
	}

	tx := sps.repo.BeginTransaction()

	dates, err := sps.parseStrategyProfileInitiativeDates(strategyProfileInitiativeFormDTO)
	if err != nil {
		return 0, err
	}

	strategyProfileInitiative := &models.StrategyProfileInitiative{
		StrategyProfileID:          strategyProfileID,
		StrategyProfileFocusAreaID: strategyProfileInitiativeFormDTO.StrategyProfileFocusAreaID,
		StrategyProfileGoalID:      strategyProfileInitiativeFormDTO.StrategyProfileGoalID,
		StrategyProfileObjectiveID: strategyProfileInitiativeFormDTO.StrategyProfileObjectiveID,
		Title:                      strategyProfileInitiativeFormDTO.Title,
		Description:                strategyProfileInitiativeFormDTO.Description,
		StartDate:                  dates.StartDate,
		EndDate:                    dates.EndDate,
		ResponsibleUserID:          strategyProfileInitiativeFormDTO.ResponsibleUserID,
		CreatedBYID:                currentUserID,
	}
	id, err := sps.repo.Create(strategyProfileInitiative, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	tx.Commit()
	go func() {
		existingStrategyProfileInitiative, err := sps.GetFullDetailsByID(strategyProfileID, id)
		if err != nil {
			utility.LogError("Failed to fetch strategy profile initiative", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}
		err = sps.activityLogSvc.LogActivity(
			models.CreateAction,
			models.StrategyProfileInitiative{}.TableName(),
			id,
			existingStrategyProfileInitiative,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.StrategyProfileInitiative{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()
	return id, nil
}
func (sps *StrategyProfileInitiativeService) Update(strategyProfileID uint, id uint, currentUserID uint, strategyProfileInitiativeFormDTO *strategyDTOs.StrategyProfileInitiativeFormDTO) error {
	if _, err := sps.strategyProfileSvc.GetByID(strategyProfileID); err != nil {
		return err
	}

	tx := sps.repo.BeginTransaction()

	dates, err := sps.parseStrategyProfileInitiativeDates(strategyProfileInitiativeFormDTO)
	if err != nil {
		return err
	}

	existingStrategyProfileInitiative, err := sps.repo.GetByID(id)
	if err != nil {
		return utility.NewResourceNotFoundError(err, id)
	}

	oldStrategyProfileInitiative, err := sps.GetFullDetailsByID(strategyProfileID, id, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	existingStrategyProfileInitiative.Title = strategyProfileInitiativeFormDTO.Title
	existingStrategyProfileInitiative.Description = strategyProfileInitiativeFormDTO.Description
	existingStrategyProfileInitiative.StartDate = dates.StartDate
	existingStrategyProfileInitiative.EndDate = dates.EndDate
	existingStrategyProfileInitiative.ResponsibleUserID = strategyProfileInitiativeFormDTO.ResponsibleUserID
	existingStrategyProfileInitiative.UpdatedBYID = &currentUserID

	err = sps.repo.Update(existingStrategyProfileInitiative, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	go func() {
		existingStrategyProfileInitiative, err := sps.GetFullDetailsByID(strategyProfileID, id)
		if err != nil {
			return
		}

		// Prepare log data
		changedData := make(map[string]interface{})

		//log changes
		oldData := map[string]interface{}{
			"title":            oldStrategyProfileInitiative.Title,
			"description":      oldStrategyProfileInitiative.Description,
			"start_date":       oldStrategyProfileInitiative.StartDateString,
			"end_date":         oldStrategyProfileInitiative.EndDateString,
			"responsible_user": utility.GetFieldValue(oldStrategyProfileInitiative.ResponsibleUser),
		}

		newData := map[string]interface{}{
			"title":            existingStrategyProfileInitiative.Title,
			"description":      existingStrategyProfileInitiative.Description,
			"start_date":       existingStrategyProfileInitiative.StartDateString,
			"end_date":         existingStrategyProfileInitiative.EndDateString,
			"responsible_user": utility.GetFieldValue(existingStrategyProfileInitiative.ResponsibleUser),
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
				models.StrategyProfileInitiative{}.TableName(),
				id,
				changedData,
				currentUserID,
				nil,
			)

			if err != nil {
				fields := map[string]interface{}{
					"action":    models.UpdateAction,
					"tableName": models.StrategyProfileInitiative{}.TableName(),
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

func (sps *StrategyProfileInitiativeService) Delete(strategyProfileID uint, ids []uint, currentUserID uint) ([]uint, error) {
	if _, err := sps.strategyProfileSvc.GetByID(strategyProfileID); err != nil {
		return nil, err
	}

	tx := sps.repo.BeginTransaction()
	var failedIds []uint
	var successfulIds []uint
	var deletedStrategyProfileInitiative []map[string]interface{}

	for _, id := range ids {
		strategyProfileInitiative, err := sps.GetFullDetailsByID(strategyProfileID, id, tx)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		deletedStrategyProfileInitiative = append(deletedStrategyProfileInitiative, map[string]interface{}{
			"id":       id,
			"title_en": strategyProfileInitiative.Title,
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
				models.StrategyProfileInitiative{}.TableName(),
				successfulIds,
				deletedStrategyProfileInitiative,
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

func (sps *StrategyProfileInitiativeService) parseStrategyProfileInitiativeDates(dto *strategyDTOs.StrategyProfileInitiativeFormDTO) (*strategyDTOs.StrategyProfileInitiativeParsedDates, error) {
	parsed := &strategyDTOs.StrategyProfileInitiativeParsedDates{}

	var err error
	if dto.StartDate != "" {
		parsed.StartDate, err = utility.ParseDate(dto.StartDate)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid start date"), "invalid_date", "start date")
		}
	}

	if dto.EndDate != "" {
		parsed.EndDate, err = utility.ParseDate(dto.EndDate)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid end date"), "invalid_date", "end date")
		}
	}

	return parsed, nil
}
