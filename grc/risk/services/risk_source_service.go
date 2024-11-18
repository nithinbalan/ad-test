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

type RiskSourceServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.RiskSource, error)
	GetByID(id uint) (*models.RiskSource, error)
	Create(currentUserID uint, source *models.RiskSource) (*models.RiskSource, error)
	Update(id uint, currentUserID uint, source *models.RiskSource) (*models.RiskSource, error)
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	Activate(ids []uint, currentUserID uint) error
	Deactivate(ids []uint, currentUserID uint) error
}

// Service and implementations
type RiskSourceService struct {
	repo           repositories.RiskSourceRepositoryInterface
	activityLogSvc services.ActivityLogServiceInterface
}

func NewRiskSourceService(repo repositories.RiskSourceRepositoryInterface, activityLogSvc services.ActivityLogServiceInterface) RiskSourceServiceInterface {
	return &RiskSourceService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
	}
}

func (rss *RiskSourceService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.RiskSource, error) {
	return rss.repo.GetAll(sort, sortBy, keyword, status)
}

func (rss RiskSourceService) GetByID(id uint) (*models.RiskSource, error) {
	riskSource, err := rss.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return riskSource, nil
}

func (rss *RiskSourceService) Create(currentUserID uint, source *models.RiskSource) (*models.RiskSource, error) {

	if exists, err := rss.repo.ExistsByTitle(source.TitleEn, source.TitleAr); err != nil {
		return nil, err
	} else if exists {
		return nil, utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
	}

	createdSource, err := rss.repo.Create(source)
	if err != nil {
		return nil, err
	}

	go func() {
		err := rss.activityLogSvc.LogActivity(
			models.CreateAction,
			models.RiskSource{}.TableName(),
			createdSource.ID,
			createdSource,
			currentUserID,
			nil,
		)

		if err != nil {
			// @todo - save to log file
			log.Printf("Failed to log activity: %v", err)
		}
	}()

	return createdSource, nil
}

func (rss *RiskSourceService) Update(id uint, currentUserID uint, updatedSource *models.RiskSource) (*models.RiskSource, error) {
	existingSource, err := rss.repo.GetByID(id)

	if err != nil {
		return nil, err
	}

	if updatedSource.TitleEn != existingSource.TitleEn || updatedSource.TitleAr != existingSource.TitleAr {
		if exists, err := rss.repo.ExistsByTitleExcludingID(id, updatedSource.TitleEn, updatedSource.TitleAr); err != nil {
			return nil, err
		} else if exists {
			return nil, utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
		}

	}

	result, err := rss.repo.Update(id, updatedSource)

	if err != nil {
		return nil, err
	}

	fieldsToCompare := []string{"TitleEn", "TitleAr", "Status"}
	changes, err := utility.CompareFields(existingSource, updatedSource, fieldsToCompare)

	if err != nil {
		log.Println("error :", err.Error())
	}

	changedDataJSON, err := utility.ChangesToJSON(changes)
	if err != nil {
		log.Println("error :", err.Error())
	}

	go func() {
		err := rss.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.RiskSource{}.TableName(),
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

func (rss *RiskSourceService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIds []uint
	var successfulIds []uint
	var deletedRiskSource []map[string]interface{}
	for _, id := range ids {

		riskSource, err := rss.repo.GetByID(id)

		if err != nil {
			fmt.Print(err)
			continue
		}

		err = rss.repo.Delete(id)
		if err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}
		deletedRiskSource = append(deletedRiskSource, map[string]interface{}{
			"ID":    id,
			"Title": riskSource.TitleEn,
		})
	}

	if len(successfulIds) > 0 {

		go func() {
			err := rss.activityLogSvc.LogBulkActivities(
				models.DeleteAction, // Make sure this action constant is defined appropriately
				models.RiskSource{}.TableName(),
				successfulIds,
				deletedRiskSource,
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

func (rss *RiskSourceService) Activate(ids []uint, currentUserID uint) error {
	successfulIds, err := rss.repo.Activate(ids)

	if len(successfulIds) > 0 {

		go func() {
			err := rss.activityLogSvc.LogBulkActivities(
				models.ActivateAction,
				models.RiskSource{}.TableName(),
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

func (rss *RiskSourceService) Deactivate(ids []uint, currentUserID uint) error {
	successfulIds, err := rss.repo.Deactivate(ids)

	if len(successfulIds) > 0 {

		go func() {
			err := rss.activityLogSvc.LogBulkActivities(
				models.DeactivateAction,
				models.RiskSource{}.TableName(),
				successfulIds,
				nil,
				currentUserID,
				nil,
			)
			if err != nil {
				// @todo - save to log file
				// You might want to log this error or handle it appropriately
				fmt.Printf("Failed to log bulk delete activity: %v\n", err)
			}
		}()
	}

	return err
}
