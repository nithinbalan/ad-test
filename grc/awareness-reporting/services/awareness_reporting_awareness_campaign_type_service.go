package services

import (
	"errors"
	"fmt"
	"log"

	"jethur/grc/awareness-reporting/dtos"
	"jethur/grc/awareness-reporting/repositories"
	"jethur/grc/logs/services"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type AwarenessCampaignTypeServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.AwarenessCampaignType, error)
	GetByID(id uint) (*models.AwarenessCampaignType, error)
	Create(currentUserID uint, awarenessCampaignTypeForm *dtos.AwarenessCampaignTypeFormDTO) (*models.AwarenessCampaignType, error)
	Update(id uint, currentUserID uint, awarenessCampaignTypeForm *dtos.AwarenessCampaignTypeFormDTO) (*models.AwarenessCampaignType, error)
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	Activate(ids []uint, currentUserID uint) error
	Deactivate(ids []uint, currentUserID uint) error
}

// Service and implementations
type AwarenessCampaignTypeService struct {
	repo           repositories.AwarenessCampaignTypeRepositoryInterface
	activityLogSvc services.ActivityLogServiceInterface
}

func NewAwarenessCampaignTypeService(repo repositories.AwarenessCampaignTypeRepositoryInterface, activityLogSvc services.ActivityLogServiceInterface) AwarenessCampaignTypeServiceInterface {
	return &AwarenessCampaignTypeService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
	}
}

func (acts *AwarenessCampaignTypeService) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.AwarenessCampaignType, error) {
	return acts.repo.GetAll(sort, sortBy, keyword, status)
}

func (acts AwarenessCampaignTypeService) GetByID(id uint) (*models.AwarenessCampaignType, error) {
	awarenessCampaignType, err := acts.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return awarenessCampaignType, nil
}

func (acts *AwarenessCampaignTypeService) Create(currentUserID uint, awarenessCampaignTypeForm *dtos.AwarenessCampaignTypeFormDTO) (*models.AwarenessCampaignType, error) {
	
	awarenessCampaignType := &models.AwarenessCampaignType{
		TitleEn: awarenessCampaignTypeForm.TitleEn,
		TitleAr: awarenessCampaignTypeForm.TitleAr,
		Type:    "others",
		Status:  awarenessCampaignTypeForm.Status,
	}

	createdType, err := acts.repo.Create(awarenessCampaignType)
	if err != nil {
		return nil, err
	}

	go func() {
		err := acts.activityLogSvc.LogActivity(
			models.CreateAction,
			models.AwarenessCampaignType{}.TableName(),
			createdType.ID,
			createdType,
			currentUserID,
			"", "", "", nil,
		)
		if err != nil {
			// @todo - save to log file
			log.Printf("Failed to log activity: %v", err)
		}
	}()

	return createdType, nil
}

func (acts *AwarenessCampaignTypeService) Update(id uint, currentUserID uint, awarenessCampaignTypeForm *dtos.AwarenessCampaignTypeFormDTO) (*models.AwarenessCampaignType, error) {
	existingType, err := acts.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	updatedType := &models.AwarenessCampaignType{
		TitleEn: awarenessCampaignTypeForm.TitleEn,
		TitleAr: awarenessCampaignTypeForm.TitleAr,
		Type:    "others",
		Status:  awarenessCampaignTypeForm.Status,
	}

	result, err := acts.repo.Update(id, updatedType)
	if err != nil {
		return nil, err
	}

	fieldsToCompare := []string{"TitleEn", "TitleAr", "Type", "Status"}
	changes, err := utility.CompareFields(existingType, updatedType, fieldsToCompare)
	if err != nil {
		log.Println("error :", err.Error())
	}

	changedDataJSON, err := utility.ChangesToJSON(changes)
	if err != nil {
		log.Println("error :", err.Error())
	}

	go func() {
		err := acts.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.AwarenessCampaignType{}.TableName(),
			result.ID,
			changedDataJSON,
			currentUserID,
			"", "", "", nil,
		)
		if err != nil {
			// @todo - save to log file
			fmt.Printf("Failed to log bulk delete activity: %v\n", err)
		}
	}()
	return result, err
}

func (acts *AwarenessCampaignTypeService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIds []uint
	var successfulIds []uint
	var deletedAwarenessCampaignType []map[string]interface{}

	for _, id := range ids {

		awarenessCampaignType, err := acts.repo.GetByID(id)
		if err != nil {
			fmt.Print(err)
			continue
		}

		err = acts.repo.Delete(id)
		if err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}

		deletedAwarenessCampaignType = append(deletedAwarenessCampaignType, map[string]interface{}{
			"id":       id,
			"title_en": awarenessCampaignType.TitleEn,
			"title_ar": awarenessCampaignType.TitleAr,
			"type":     awarenessCampaignType.Type,
		})
	}

	if len(successfulIds) > 0 {
		go func() {
			err := acts.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.AwarenessCampaignType{}.TableName(),
				successfulIds,
				deletedAwarenessCampaignType,
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

func (acts *AwarenessCampaignTypeService) Activate(ids []uint, currentUserID uint) error {
	successfulIds, err := acts.repo.Activate(ids)

	if len(successfulIds) > 0 {
		go func() {
			err := acts.activityLogSvc.LogBulkActivities(
				models.ActivateAction,
				models.AwarenessCampaignType{}.TableName(),
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

func (acts *AwarenessCampaignTypeService) Deactivate(ids []uint, currentUserID uint) error {
	successfulIds, err := acts.repo.Deactivate(ids)

	if len(successfulIds) > 0 {
		go func() {
			err := acts.activityLogSvc.LogBulkActivities(
				models.DeactivateAction,
				models.AwarenessCampaignType{}.TableName(),
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