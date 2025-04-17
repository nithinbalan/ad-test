package services

import (
	"errors"
	"fmt"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	strategyDTOs "jethur/grc/strategy/dtos"
	strategyProfileActionPlanRepositories "jethur/grc/strategy/repositories"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type StrategyProfileActionPlanServiceInterface interface {
	GetAll(strategyProfileID uint, params *strategyDTOs.StrategyProfileActionPlanListQueryParamsDTO) ([]*strategyDTOs.StrategyProfileActionPlanListDTO, error)
	GetFullDetailsByID(strategyProfileID uint, id uint, tx ...utility.DatabaseTransaction) (*strategyDTOs.StrategyProfileActionPlanDetailDTO, error)
	GetByID(id uint) (*models.StrategyProfileActionPlan, error)
	Create(strategyProfileID uint, currentUserID uint, strategyDTOs *strategyDTOs.StrategyProfileActionPlanFormDTO) (uint, error)
	Update(strategyProfileID uint, id uint, currentUserID uint, strategyProfileActionPlan *strategyDTOs.StrategyProfileActionPlanFormDTO) error
	Delete(strategyProfileID uint, ids []uint, currentUserID uint) ([]uint, error)
}

// Service and implementations
type StrategyProfileActionPlanService struct {
	repo               strategyProfileActionPlanRepositories.StrategyProfileActionPlanRepositoryInterface
	activityLogSvc     logServices.ActivityLogServiceInterface
	userSvc            organizationServices.UserServiceInterface
	userSettingSvc     organizationServices.UserSettingServiceInterface
	strategyProfileSvc StrategyProfileServiceInterface
}

func NewStrategyProfileActionPlanService(
	repo strategyProfileActionPlanRepositories.StrategyProfileActionPlanRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	userSettingSvc organizationServices.UserSettingServiceInterface,
	strategyProfileSvc StrategyProfileServiceInterface,
) StrategyProfileActionPlanServiceInterface {
	return &StrategyProfileActionPlanService{
		repo:               repo,
		activityLogSvc:     activityLogSvc,
		userSvc:            userSvc,
		userSettingSvc:     userSettingSvc,
		strategyProfileSvc: strategyProfileSvc,
	}
}

func (sps *StrategyProfileActionPlanService) GetAll(strategyProfileID uint, params *strategyDTOs.StrategyProfileActionPlanListQueryParamsDTO) ([]*strategyDTOs.StrategyProfileActionPlanListDTO, error) {
	return sps.repo.GetAll(strategyProfileID, params)
}

func (sps StrategyProfileActionPlanService) GetByID(id uint) (*models.StrategyProfileActionPlan, error) {
	strategyProfileActionPlan, err := sps.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return strategyProfileActionPlan, nil
}

func (sps StrategyProfileActionPlanService) GetFullDetailsByID(strategyProfileID uint, id uint, tx ...utility.DatabaseTransaction) (*strategyDTOs.StrategyProfileActionPlanDetailDTO, error) {
	// Initialize your DTO
	var strategyProfileActionPlanDTO strategyDTOs.StrategyProfileActionPlanDetailDTO

	if _, err := sps.strategyProfileSvc.GetByID(strategyProfileID); err != nil {
		return nil, err
	}

	strategyProfileActionPlan, err := sps.repo.GetFullDetailsByID(id, tx...)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	// Map the result to your DTO
	err = mapstructure.Decode(strategyProfileActionPlan, &strategyProfileActionPlanDTO)
	if err != nil {
		return nil, err
	}

	if strategyProfileActionPlan.TargetDate != nil {
		strategyProfileActionPlanDTO.TargetDateString = strategyProfileActionPlan.TargetDate.Format("2006-01-02")
	}

	strategyProfileActionPlanDTO.StrategyProfile.StartDateString = strategyProfileActionPlan.StrategyProfile.StartDate.Format("2006-01-02")
	strategyProfileActionPlanDTO.StrategyProfile.EndDateString = strategyProfileActionPlan.StrategyProfile.EndDate.Format("2006-01-02")

	strategyProfileActionPlanDTO.CreatedAt = strategyProfileActionPlan.CreatedAt
	strategyProfileActionPlanDTO.UpdatedAt = strategyProfileActionPlan.UpdatedAt

	return &strategyProfileActionPlanDTO, nil
}

func (sps *StrategyProfileActionPlanService) Create(strategyProfileID uint, currentUserID uint, strategyProfileActionPlanFormDTO *strategyDTOs.StrategyProfileActionPlanFormDTO) (uint, error) {

	if _, err := sps.strategyProfileSvc.GetByID(strategyProfileID); err != nil {
		return 0, err
	}

	tx := sps.repo.BeginTransaction()

	dates, err := sps.parseStrategyProfileActionPlanDates(strategyProfileActionPlanFormDTO)
	if err != nil {
		return 0, err
	}

	strategyProfileActionPlan := &models.StrategyProfileActionPlan{
		StrategyProfileID:           strategyProfileID,
		StrategyProfileFocusAreaID:  strategyProfileActionPlanFormDTO.StrategyProfileFocusAreaID,
		StrategyProfileGoalID:       strategyProfileActionPlanFormDTO.StrategyProfileGoalID,
		StrategyProfileObjectiveID:  strategyProfileActionPlanFormDTO.StrategyProfileObjectiveID,
		StrategyProfileInitiativeID: strategyProfileActionPlanFormDTO.StrategyProfileInitiativeID,
		StrategyProfileProjectID:    strategyProfileActionPlanFormDTO.StrategyProfileProjectID,
		StrategyProfileMilestoneID:  strategyProfileActionPlanFormDTO.StrategyProfileMilestoneID,
		Title:                       strategyProfileActionPlanFormDTO.Title,
		TargetDate:                  dates.TargetDate,
		ResponsibleUserID:           strategyProfileActionPlanFormDTO.ResponsibleUserID,
		CreatedBYID:                 currentUserID,
	}
	id, err := sps.repo.Create(strategyProfileActionPlan, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	tx.Commit()
	go func() {
		existingStrategyProfileActionPlan, err := sps.GetFullDetailsByID(strategyProfileID, id)
		if err != nil {
			utility.LogError("Failed to fetch strategy profile action plan", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}
		err = sps.activityLogSvc.LogActivity(
			models.CreateAction,
			models.StrategyProfileActionPlan{}.TableName(),
			id,
			existingStrategyProfileActionPlan,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.StrategyProfileActionPlan{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()
	return id, nil
}
func (sps *StrategyProfileActionPlanService) Update(strategyProfileID uint, id uint, currentUserID uint, strategyProfileActionPlanFormDTO *strategyDTOs.StrategyProfileActionPlanFormDTO) error {
	if _, err := sps.strategyProfileSvc.GetByID(strategyProfileID); err != nil {
		return err
	}

	tx := sps.repo.BeginTransaction()

	dates, err := sps.parseStrategyProfileActionPlanDates(strategyProfileActionPlanFormDTO)
	if err != nil {
		return err
	}

	existingStrategyProfileActionPlan, err := sps.repo.GetByID(id)
	if err != nil {
		return utility.NewResourceNotFoundError(err, id)
	}

	oldStrategyProfileActionPlan, err := sps.GetFullDetailsByID(strategyProfileID, id, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	existingStrategyProfileActionPlan.Title = strategyProfileActionPlanFormDTO.Title
	existingStrategyProfileActionPlan.TargetDate = dates.TargetDate
	existingStrategyProfileActionPlan.ResponsibleUserID = strategyProfileActionPlanFormDTO.ResponsibleUserID
	existingStrategyProfileActionPlan.UpdatedBYID = &currentUserID

	err = sps.repo.Update(existingStrategyProfileActionPlan, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	go func() {
		existingStrategyProfileActionPlan, err := sps.GetFullDetailsByID(strategyProfileID, id)
		if err != nil {
			return
		}

		// Prepare log data
		changedData := make(map[string]interface{})

		//log changes
		oldData := map[string]interface{}{
			"title":            oldStrategyProfileActionPlan.Title,
			"target_date":      oldStrategyProfileActionPlan.TargetDateString,
			"responsible_user": utility.GetFieldValue(oldStrategyProfileActionPlan.ResponsibleUser),
		}

		newData := map[string]interface{}{
			"title":            existingStrategyProfileActionPlan.Title,
			"target_date":      existingStrategyProfileActionPlan.TargetDateString,
			"responsible_user": utility.GetFieldValue(existingStrategyProfileActionPlan.ResponsibleUser),
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
				models.StrategyProfileActionPlan{}.TableName(),
				id,
				changedData,
				currentUserID,
				nil,
			)

			if err != nil {
				fields := map[string]interface{}{
					"action":    models.UpdateAction,
					"tableName": models.StrategyProfileActionPlan{}.TableName(),
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

func (sps *StrategyProfileActionPlanService) Delete(strategyProfileID uint, ids []uint, currentUserID uint) ([]uint, error) {
	if _, err := sps.strategyProfileSvc.GetByID(strategyProfileID); err != nil {
		return nil, err
	}

	tx := sps.repo.BeginTransaction()
	var failedIds []uint
	var successfulIds []uint
	var deletedStrategyProfileActionPlan []map[string]interface{}

	for _, id := range ids {
		strategyProfileActionPlan, err := sps.GetFullDetailsByID(strategyProfileID, id, tx)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		deletedStrategyProfileActionPlan = append(deletedStrategyProfileActionPlan, map[string]interface{}{
			"id":       id,
			"title_en": strategyProfileActionPlan.Title,
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
				models.StrategyProfileActionPlan{}.TableName(),
				successfulIds,
				deletedStrategyProfileActionPlan,
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

func (sps *StrategyProfileActionPlanService) parseStrategyProfileActionPlanDates(dto *strategyDTOs.StrategyProfileActionPlanFormDTO) (*strategyDTOs.StrategyProfileActionPlanParsedDates, error) {
	parsed := &strategyDTOs.StrategyProfileActionPlanParsedDates{}

	var err error
	if dto.TargetDate != "" {
		parsed.TargetDate, err = utility.ParseDate(dto.TargetDate)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid start date"), "invalid_date", "start date")
		}
	}

	return parsed, nil
}
