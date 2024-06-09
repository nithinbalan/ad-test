package services

import (
	"errors"
	"fmt"
	"jethur/grc/logs/services"
	"jethur/grc/models"
	"jethur/grc/risk/repositories"
	"jethur/utility"
	"log"

	"gorm.io/gorm"
)

type RiskResponseStrategyServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.RiskResponseStrategy, error)
	GetByID(id uint) (*models.RiskResponseStrategy, error)
	Create(currentUserID uint, strategy *models.RiskResponseStrategy) (*models.RiskResponseStrategy, error)
	Update(id uint, currentUserID uint, strategy *models.RiskResponseStrategy) (*models.RiskResponseStrategy, error)
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	Activate(ids []uint, currentUserID uint) error
	Deactivate(ids []uint, currentUserID uint) error
}

// Service and implementations
type RiskResponseStrategyService struct {
	repo           repositories.RiskResponseStrategyRepositoryInterface
	activityLogSvc services.ActivityLogServiceInterface
}

func NewRiskResponseStrategyService(repo repositories.RiskResponseStrategyRepositoryInterface, activityLogSvc services.ActivityLogServiceInterface) RiskResponseStrategyServiceInterface {
	return &RiskResponseStrategyService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
	}
}

func (rrss *RiskResponseStrategyService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.RiskResponseStrategy, error) {
	return rrss.repo.GetAll(sort, sortBy, keyword, status)
}

func (rrss RiskResponseStrategyService) GetByID(id uint) (*models.RiskResponseStrategy, error) {
	riskResponseStrategy, err := rrss.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return riskResponseStrategy, nil
}

func (rrss *RiskResponseStrategyService) Create(currentUserID uint, strategy *models.RiskResponseStrategy) (*models.RiskResponseStrategy, error) {

	if exists, err := rrss.repo.ExistsByTitle(strategy.TitleEn, strategy.TitleAr); err != nil {
		return nil, err
	} else if exists {
		return nil, utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
	}

	createdStrategy, err := rrss.repo.Create(strategy)
	if err != nil {
		return nil, err
	}

	go func() {
		err := rrss.activityLogSvc.LogActivity(
			models.CreateAction,
			models.RiskResponseStrategy{}.TableName(),
			createdStrategy.ID,
			createdStrategy,
			currentUserID,
			nil,
		)

		if err != nil {
			// @todo - save to log file
			log.Printf("Failed to log activity: %v", err)
		}
	}()

	return createdStrategy, nil
}

func (rrss *RiskResponseStrategyService) Update(id uint, currentUserID uint, updatedStrategy *models.RiskResponseStrategy) (*models.RiskResponseStrategy, error) {
	existingStrategy, err := rrss.repo.GetByID(id)

	if err != nil {
		return nil, err
	}

	if updatedStrategy.TitleEn != existingStrategy.TitleEn || updatedStrategy.TitleAr != existingStrategy.TitleAr {
		if exists, err := rrss.repo.ExistsByTitleExcludingID(id, updatedStrategy.TitleEn, updatedStrategy.TitleAr); err != nil {
			return nil, err
		} else if exists {
			return nil, utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
		}

	}

	result, err := rrss.repo.Update(id, updatedStrategy)

	if err != nil {
		return nil, err
	}

	fieldsToCompare := []string{"TitleEn", "TitleAr", "DescriptionEn", "DescriptionAr", "Status"}
	changes, err := utility.CompareFields(existingStrategy, updatedStrategy, fieldsToCompare)

	if err != nil {
		log.Println("error :", err.Error())
	}

	changedDataJSON, err := utility.ChangesToJSON(changes)
	if err != nil {
		log.Println("error :", err.Error())
	}

	go func() {
		err := rrss.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.RiskResponseStrategy{}.TableName(),
			result.ID,
			changedDataJSON,
			currentUserID,
			nil,
		)
		if err != nil {
			// @todo - save to log file
			fmt.Printf("Failed to log bulk delete activity: %v\n", err)
		}
	}()
	return result, err
}

func (rrss *RiskResponseStrategyService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIds []uint
	var successfulIds []uint
	for _, id := range ids {

		_, err := rrss.repo.GetByID(id)

		if err != nil {
			fmt.Print(err)
			continue
		}

		err = rrss.repo.Delete(id)
		if err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}
	}

	if len(successfulIds) > 0 {

		go func() {
			err := rrss.activityLogSvc.LogBulkActivities(
				models.DeleteAction, // Make sure this action constant is defined appropriately
				models.RiskResponseStrategy{}.TableName(),
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

func (rrss *RiskResponseStrategyService) Activate(ids []uint, currentUserID uint) error {
	successfulIds, err := rrss.repo.Activate(ids)

	if len(successfulIds) > 0 {

		go func() {
			err := rrss.activityLogSvc.LogBulkActivities(
				models.ActivateAction,
				models.RiskResponseStrategy{}.TableName(),
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

	return err
}

func (rrss *RiskResponseStrategyService) Deactivate(ids []uint, currentUserID uint) error {
	successfulIds, err := rrss.repo.Deactivate(ids)

	if len(successfulIds) > 0 {

		go func() {
			err := rrss.activityLogSvc.LogBulkActivities(
				models.DeactivateAction,
				models.RiskResponseStrategy{}.TableName(),
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

	return err
}
