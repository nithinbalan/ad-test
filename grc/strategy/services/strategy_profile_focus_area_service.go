package services

import (
	"errors"
	"fmt"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	strategyDTOs "jethur/grc/strategy/dtos"
	strategyProfileFocusAreaRepositories "jethur/grc/strategy/repositories"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type StrategyProfileFocusAreaServiceInterface interface {
	GetAll(strategyProfileID uint, params *strategyDTOs.StrategyProfileFocusAreaListQueryParamsDTO) ([]*strategyDTOs.StrategyProfileFocusAreaListDTO, error)
	GetFullDetailsByID(strategyProfileID uint, id uint, tx ...utility.DatabaseTransaction) (*strategyDTOs.StrategyProfileFocusAreaDetailDTO, error)
	GetByID(id uint) (*models.StrategyProfileFocusArea, error)
	Create(strategyProfileID uint, currentUserID uint, strategyDTOs *strategyDTOs.StrategyProfileFocusAreaFormDTO) (uint, error)
	Update(strategyProfileID uint, id uint, currentUserID uint, strategyProfileFocusArea *strategyDTOs.StrategyProfileFocusAreaFormDTO) error
	Delete(strategyProfileID uint, ids []uint, currentUserID uint) ([]uint, error)
}

// Service and implementations
type StrategyProfileFocusAreaService struct {
	repo               strategyProfileFocusAreaRepositories.StrategyProfileFocusAreaRepositoryInterface
	activityLogSvc     logServices.ActivityLogServiceInterface
	userSvc            organizationServices.UserServiceInterface
	userSettingSvc     organizationServices.UserSettingServiceInterface
	strategyProfileSvc StrategyProfileServiceInterface
}

func NewStrategyProfileFocusAreaService(
	repo strategyProfileFocusAreaRepositories.StrategyProfileFocusAreaRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	userSettingSvc organizationServices.UserSettingServiceInterface,
	strategyProfileSvc StrategyProfileServiceInterface,
) StrategyProfileFocusAreaServiceInterface {
	return &StrategyProfileFocusAreaService{
		repo:               repo,
		activityLogSvc:     activityLogSvc,
		userSvc:            userSvc,
		userSettingSvc:     userSettingSvc,
		strategyProfileSvc: strategyProfileSvc,
	}
}

func (sps *StrategyProfileFocusAreaService) GetAll(strategyProfileID uint, params *strategyDTOs.StrategyProfileFocusAreaListQueryParamsDTO) ([]*strategyDTOs.StrategyProfileFocusAreaListDTO, error) {
	return sps.repo.GetAll(strategyProfileID, params)
}

func (sps StrategyProfileFocusAreaService) GetByID(id uint) (*models.StrategyProfileFocusArea, error) {
	strategyProfileFocusArea, err := sps.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return strategyProfileFocusArea, nil
}

func (sps StrategyProfileFocusAreaService) GetFullDetailsByID(strategyProfileID uint, id uint, tx ...utility.DatabaseTransaction) (*strategyDTOs.StrategyProfileFocusAreaDetailDTO, error) {
	// Initialize your DTO
	var strategyProfileFocusAreaDTO strategyDTOs.StrategyProfileFocusAreaDetailDTO

	if _, err := sps.strategyProfileSvc.GetByID(strategyProfileID); err != nil {
		return nil, err
	}

	strategyProfileFocusArea, err := sps.repo.GetFullDetailsByID(id, tx...)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	// Map the result to your DTO
	err = mapstructure.Decode(strategyProfileFocusArea, &strategyProfileFocusAreaDTO)
	if err != nil {
		return nil, err
	}

	strategyProfileFocusAreaDTO.StrategyProfile.StartDateString = strategyProfileFocusArea.StrategyProfile.StartDate.Format("2006-01-02")
	strategyProfileFocusAreaDTO.StrategyProfile.EndDateString = strategyProfileFocusArea.StrategyProfile.EndDate.Format("2006-01-02")

	strategyProfileFocusAreaDTO.CreatedAt = strategyProfileFocusArea.CreatedAt
	strategyProfileFocusAreaDTO.UpdatedAt = strategyProfileFocusArea.UpdatedAt

	return &strategyProfileFocusAreaDTO, nil
}

func (sps *StrategyProfileFocusAreaService) Create(strategyProfileID uint, currentUserID uint, strategyProfileFocusAreaFormDTO *strategyDTOs.StrategyProfileFocusAreaFormDTO) (uint, error) {

	if _, err := sps.strategyProfileSvc.GetByID(strategyProfileID); err != nil {
		return 0, err
	}

	tx := sps.repo.BeginTransaction()

	strategyProfileFocusArea := &models.StrategyProfileFocusArea{
		StrategyProfileID: strategyProfileID,
		Title:             strategyProfileFocusAreaFormDTO.Title,
		Description:       strategyProfileFocusAreaFormDTO.Description,
		CreatedBYID:       currentUserID,
	}
	id, err := sps.repo.Create(strategyProfileFocusArea, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	tx.Commit()
	go func() {
		existingStrategyProfileFocusArea, err := sps.GetFullDetailsByID(strategyProfileID, id)
		if err != nil {
			utility.LogError("Failed to fetch strategy profile", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}
		err = sps.activityLogSvc.LogActivity(
			models.CreateAction,
			models.StrategyProfileFocusArea{}.TableName(),
			id,
			existingStrategyProfileFocusArea,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.StrategyProfileFocusArea{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()
	return id, nil
}
func (sps *StrategyProfileFocusAreaService) Update(strategyProfileID uint, id uint, currentUserID uint, strategyProfileFocusAreaFormDTO *strategyDTOs.StrategyProfileFocusAreaFormDTO) error {
	if _, err := sps.strategyProfileSvc.GetByID(strategyProfileID); err != nil {
		return err
	}

	tx := sps.repo.BeginTransaction()

	existingStrategyProfileFocusArea, err := sps.repo.GetByID(id)
	if err != nil {
		return utility.NewResourceNotFoundError(err, id)
	}

	oldStrategyProfileFocusArea, err := sps.GetFullDetailsByID(strategyProfileID, id, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	existingStrategyProfileFocusArea.Title = strategyProfileFocusAreaFormDTO.Title
	existingStrategyProfileFocusArea.Description = strategyProfileFocusAreaFormDTO.Description
	existingStrategyProfileFocusArea.UpdatedBYID = &currentUserID

	err = sps.repo.Update(existingStrategyProfileFocusArea, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	go func() {
		existingStrategyProfileFocusArea, err := sps.GetFullDetailsByID(strategyProfileID, id)
		if err != nil {
			return
		}

		// Prepare log data
		changedData := make(map[string]interface{})

		//log changes
		oldData := map[string]interface{}{
			"title":       oldStrategyProfileFocusArea.Title,
			"description": oldStrategyProfileFocusArea.Description,
		}

		newData := map[string]interface{}{
			"title":       existingStrategyProfileFocusArea.Title,
			"description": existingStrategyProfileFocusArea.Description,
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
				models.StrategyProfileFocusArea{}.TableName(),
				id,
				changedData,
				currentUserID,
				nil,
			)

			if err != nil {
				fields := map[string]interface{}{
					"action":    models.UpdateAction,
					"tableName": models.StrategyProfileFocusArea{}.TableName(),
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

func (sps *StrategyProfileFocusAreaService) Delete(strategyProfileID uint, ids []uint, currentUserID uint) ([]uint, error) {
	if _, err := sps.strategyProfileSvc.GetByID(strategyProfileID); err != nil {
		return nil, err
	}

	tx := sps.repo.BeginTransaction()
	var failedIds []uint
	var successfulIds []uint
	var deletedStrategyProfileFocusArea []map[string]interface{}

	for _, id := range ids {
		strategyProfileFocusArea, err := sps.GetFullDetailsByID(strategyProfileID, id, tx)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		deletedStrategyProfileFocusArea = append(deletedStrategyProfileFocusArea, map[string]interface{}{
			"id":       id,
			"title_en": strategyProfileFocusArea.Title,
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
				models.StrategyProfileFocusArea{}.TableName(),
				successfulIds,
				deletedStrategyProfileFocusArea,
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
