package services

import (
	"errors"
	"fmt"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	strategyDTOs "jethur/grc/strategy/dtos"
	strategyProfileProjectRepositories "jethur/grc/strategy/repositories"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type StrategyProfileProjectServiceInterface interface {
	GetAll(strategyProfileID uint, params *strategyDTOs.StrategyProfileProjectListQueryParamsDTO) ([]*strategyDTOs.StrategyProfileProjectListDTO, error)
	GetFullDetailsByID(strategyProfileID uint, id uint, tx ...utility.DatabaseTransaction) (*strategyDTOs.StrategyProfileProjectDetailDTO, error)
	GetByID(id uint) (*models.StrategyProfileProject, error)
	Create(strategyProfileID uint, currentUserID uint, strategyDTOs *strategyDTOs.StrategyProfileProjectFormDTO) (uint, error)
	Update(strategyProfileID uint, id uint, currentUserID uint, strategyProfileProject *strategyDTOs.StrategyProfileProjectFormDTO) error
	Delete(strategyProfileID uint, ids []uint, currentUserID uint) ([]uint, error)
}

// Service and implementations
type StrategyProfileProjectService struct {
	repo               strategyProfileProjectRepositories.StrategyProfileProjectRepositoryInterface
	activityLogSvc     logServices.ActivityLogServiceInterface
	userSvc            organizationServices.UserServiceInterface
	userSettingSvc     organizationServices.UserSettingServiceInterface
	strategyProfileSvc StrategyProfileServiceInterface
}

func NewStrategyProfileProjectService(
	repo strategyProfileProjectRepositories.StrategyProfileProjectRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	userSettingSvc organizationServices.UserSettingServiceInterface,
	strategyProfileSvc StrategyProfileServiceInterface,
) StrategyProfileProjectServiceInterface {
	return &StrategyProfileProjectService{
		repo:               repo,
		activityLogSvc:     activityLogSvc,
		userSvc:            userSvc,
		userSettingSvc:     userSettingSvc,
		strategyProfileSvc: strategyProfileSvc,
	}
}

func (sps *StrategyProfileProjectService) GetAll(strategyProfileID uint, params *strategyDTOs.StrategyProfileProjectListQueryParamsDTO) ([]*strategyDTOs.StrategyProfileProjectListDTO, error) {
	return sps.repo.GetAll(strategyProfileID, params)
}

func (sps StrategyProfileProjectService) GetByID(id uint) (*models.StrategyProfileProject, error) {
	strategyProfileProject, err := sps.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return strategyProfileProject, nil
}

func (sps StrategyProfileProjectService) GetFullDetailsByID(strategyProfileID uint, id uint, tx ...utility.DatabaseTransaction) (*strategyDTOs.StrategyProfileProjectDetailDTO, error) {
	// Initialize your DTO
	var strategyProfileProjectDTO strategyDTOs.StrategyProfileProjectDetailDTO

	if _, err := sps.strategyProfileSvc.GetByID(strategyProfileID); err != nil {
		return nil, err
	}

	strategyProfileProject, err := sps.repo.GetFullDetailsByID(id, tx...)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	// Map the result to your DTO
	err = mapstructure.Decode(strategyProfileProject, &strategyProfileProjectDTO)
	if err != nil {
		return nil, err
	}

	if strategyProfileProject.StartDate != nil {
		strategyProfileProjectDTO.StartDateString = strategyProfileProject.StartDate.Format("2006-01-02")
	}

	if strategyProfileProject.EndDate != nil {
		strategyProfileProjectDTO.EndDateString = strategyProfileProject.EndDate.Format("2006-01-02")
	}

	strategyProfileProjectDTO.StrategyProfile.StartDateString = strategyProfileProject.StrategyProfile.StartDate.Format("2006-01-02")
	strategyProfileProjectDTO.StrategyProfile.EndDateString = strategyProfileProject.StrategyProfile.EndDate.Format("2006-01-02")

	strategyProfileProjectDTO.CreatedAt = strategyProfileProject.CreatedAt
	strategyProfileProjectDTO.UpdatedAt = strategyProfileProject.UpdatedAt

	return &strategyProfileProjectDTO, nil
}

func (sps *StrategyProfileProjectService) Create(strategyProfileID uint, currentUserID uint, strategyProfileProjectFormDTO *strategyDTOs.StrategyProfileProjectFormDTO) (uint, error) {

	if _, err := sps.strategyProfileSvc.GetByID(strategyProfileID); err != nil {
		return 0, err
	}

	tx := sps.repo.BeginTransaction()

	dates, err := sps.parseStrategyProfileProjectDates(strategyProfileProjectFormDTO)
	if err != nil {
		return 0, err
	}

	strategyProfileProject := &models.StrategyProfileProject{
		StrategyProfileID:           strategyProfileID,
		StrategyProfileFocusAreaID:  strategyProfileProjectFormDTO.StrategyProfileFocusAreaID,
		StrategyProfileGoalID:       strategyProfileProjectFormDTO.StrategyProfileGoalID,
		StrategyProfileObjectiveID:  strategyProfileProjectFormDTO.StrategyProfileObjectiveID,
		StrategyProfileInitiativeID: strategyProfileProjectFormDTO.StrategyProfileInitiativeID,
		Title:                       strategyProfileProjectFormDTO.Title,
		Description:                 strategyProfileProjectFormDTO.Description,
		StartDate:                   dates.StartDate,
		EndDate:                     dates.EndDate,
		ResponsibleUserID:           strategyProfileProjectFormDTO.ResponsibleUserID,
		CreatedBYID:                 currentUserID,
	}
	id, err := sps.repo.Create(strategyProfileProject, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	tx.Commit()
	go func() {
		existingStrategyProfileProject, err := sps.GetFullDetailsByID(strategyProfileID, id)
		if err != nil {
			utility.LogError("Failed to fetch strategy profile project", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}
		err = sps.activityLogSvc.LogActivity(
			models.CreateAction,
			models.StrategyProfileProject{}.TableName(),
			id,
			existingStrategyProfileProject,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.StrategyProfileProject{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()
	return id, nil
}
func (sps *StrategyProfileProjectService) Update(strategyProfileID uint, id uint, currentUserID uint, strategyProfileProjectFormDTO *strategyDTOs.StrategyProfileProjectFormDTO) error {
	if _, err := sps.strategyProfileSvc.GetByID(strategyProfileID); err != nil {
		return err
	}

	tx := sps.repo.BeginTransaction()

	dates, err := sps.parseStrategyProfileProjectDates(strategyProfileProjectFormDTO)
	if err != nil {
		return err
	}

	existingStrategyProfileProject, err := sps.repo.GetByID(id)
	if err != nil {
		return utility.NewResourceNotFoundError(err, id)
	}

	oldStrategyProfileProject, err := sps.GetFullDetailsByID(strategyProfileID, id, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	existingStrategyProfileProject.Title = strategyProfileProjectFormDTO.Title
	existingStrategyProfileProject.Description = strategyProfileProjectFormDTO.Description
	existingStrategyProfileProject.StartDate = dates.StartDate
	existingStrategyProfileProject.EndDate = dates.EndDate
	existingStrategyProfileProject.ResponsibleUserID = strategyProfileProjectFormDTO.ResponsibleUserID
	existingStrategyProfileProject.UpdatedBYID = &currentUserID

	err = sps.repo.Update(existingStrategyProfileProject, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	go func() {
		existingStrategyProfileProject, err := sps.GetFullDetailsByID(strategyProfileID, id)
		if err != nil {
			return
		}

		// Prepare log data
		changedData := make(map[string]interface{})

		//log changes
		oldData := map[string]interface{}{
			"title":            oldStrategyProfileProject.Title,
			"description":      oldStrategyProfileProject.Description,
			"start_date":       oldStrategyProfileProject.StartDateString,
			"end_date":         oldStrategyProfileProject.EndDateString,
			"responsible_user": utility.GetFieldValue(oldStrategyProfileProject.ResponsibleUser),
		}

		newData := map[string]interface{}{
			"title":            existingStrategyProfileProject.Title,
			"description":      existingStrategyProfileProject.Description,
			"start_date":       existingStrategyProfileProject.StartDateString,
			"end_date":         existingStrategyProfileProject.EndDateString,
			"responsible_user": utility.GetFieldValue(existingStrategyProfileProject.ResponsibleUser),
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
				models.StrategyProfileProject{}.TableName(),
				id,
				changedData,
				currentUserID,
				nil,
			)

			if err != nil {
				fields := map[string]interface{}{
					"action":    models.UpdateAction,
					"tableName": models.StrategyProfileProject{}.TableName(),
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

func (sps *StrategyProfileProjectService) Delete(strategyProfileID uint, ids []uint, currentUserID uint) ([]uint, error) {
	if _, err := sps.strategyProfileSvc.GetByID(strategyProfileID); err != nil {
		return nil, err
	}

	tx := sps.repo.BeginTransaction()
	var failedIds []uint
	var successfulIds []uint
	var deletedStrategyProfileProject []map[string]interface{}

	for _, id := range ids {
		strategyProfileProject, err := sps.GetFullDetailsByID(strategyProfileID, id, tx)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		deletedStrategyProfileProject = append(deletedStrategyProfileProject, map[string]interface{}{
			"id":       id,
			"title_en": strategyProfileProject.Title,
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
				models.StrategyProfileProject{}.TableName(),
				successfulIds,
				deletedStrategyProfileProject,
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

func (sps *StrategyProfileProjectService) parseStrategyProfileProjectDates(dto *strategyDTOs.StrategyProfileProjectFormDTO) (*strategyDTOs.StrategyProfileProjectParsedDates, error) {
	parsed := &strategyDTOs.StrategyProfileProjectParsedDates{}

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
