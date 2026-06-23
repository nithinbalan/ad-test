package services

import (
	"errors"
	"fmt"
	"log"

	assetDTOs "jethur/grc/asset/dtos"
	"jethur/grc/asset/repositories"
	activityLogServices "jethur/grc/logs/services"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type SystemRegisterDisruptionLevelServiceInterface interface {
	GetAll(params *assetDTOs.SystemRegisterDisruptionLevelListQueryParamsDTO) ([]*assetDTOs.SystemRegisterDisruptionLevelDTO, error)
	GetByID(id uint) (*models.SystemRegisterDisruptionLevel, error)
	Create(currentUserID uint, levelFormDTO *assetDTOs.SystemRegisterDisruptionLevelFormDTO) (uint, error)
	Update(id uint, currentUserID uint, levelFormDTO *assetDTOs.SystemRegisterDisruptionLevelFormDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	Activate(ids []uint, currentUserID uint) error
	Deactivate(ids []uint, currentUserID uint) error
}

// Service and implementations
type SystemRegisterDisruptionLevelService struct {
	repo           repositories.SystemRegisterDisruptionLevelRepositoryInterface
	activityLogSvc activityLogServices.ActivityLogServiceInterface
}

func NewSystemRegisterDisruptionLevelService(repo repositories.SystemRegisterDisruptionLevelRepositoryInterface, activityLogSvc activityLogServices.ActivityLogServiceInterface) SystemRegisterDisruptionLevelServiceInterface {
	return &SystemRegisterDisruptionLevelService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
	}
}

func (srdls *SystemRegisterDisruptionLevelService) GetAll(params *assetDTOs.SystemRegisterDisruptionLevelListQueryParamsDTO) ([]*assetDTOs.SystemRegisterDisruptionLevelDTO, error) {
	return srdls.repo.GetAll(params)
}

func (srdls SystemRegisterDisruptionLevelService) GetByID(id uint) (*models.SystemRegisterDisruptionLevel, error) {
	level, err := srdls.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return level, nil
}

func (srdls *SystemRegisterDisruptionLevelService) Create(currentUserID uint, levelFormDTO *assetDTOs.SystemRegisterDisruptionLevelFormDTO) (uint, error) {
	level := &models.SystemRegisterDisruptionLevel{
		TitleEn:     levelFormDTO.TitleEn,
		TitleAr:     levelFormDTO.TitleAr,
		Status:      levelFormDTO.Status,
		CreatedBYID: currentUserID,
	}

	id, err := srdls.repo.Create(level)
	if err != nil {
		return 0, err
	}

	go func() {
		createdLevel, err := srdls.GetByID(id)
		if err != nil {
			return
		}

		err = srdls.activityLogSvc.LogActivity(
			models.CreateAction,
			models.SystemRegisterDisruptionLevel{}.TableName(),
			id,
			createdLevel,
			currentUserID,
			"",
			"",
			"",
			nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.SystemRegisterDisruptionLevel{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return id, nil
}

func (srdls *SystemRegisterDisruptionLevelService) Update(id uint, currentUserID uint, levelFormDTO *assetDTOs.SystemRegisterDisruptionLevelFormDTO) error {
	oldLevel, err := srdls.GetByID(id)
	if err != nil {
		return err
	}

	existingLevel, err := srdls.GetByID(id)
	if err != nil {
		return err
	}

	existingLevel.TitleEn = levelFormDTO.TitleEn
	existingLevel.TitleAr = levelFormDTO.TitleAr
	existingLevel.Status = levelFormDTO.Status
	existingLevel.UpdatedBYID = &currentUserID

	err = srdls.repo.Update(existingLevel)
	if err != nil {
		return err
	}

	fieldsToCompare := []string{"TitleEn", "TitleAr", "Status"}
	changes, err := utility.CompareFields(oldLevel, existingLevel, fieldsToCompare)
	if err != nil {
		log.Println("error :", err.Error())
	}

	changedDataJSON, err := utility.ChangesToJSON(changes)
	if err != nil {
		log.Println("error :", err.Error())
	}

	go func() {
		err := srdls.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.SystemRegisterDisruptionLevel{}.TableName(),
			id,
			changedDataJSON,
			currentUserID,
			"",
			"",
			"",
			nil,
		)
		if err != nil {
			// @todo - save to log file
			fmt.Printf("Failed to log bulk delete activity: %v\n", err)
		}
	}()

	return nil
}

func (srdls *SystemRegisterDisruptionLevelService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIds []uint
	var successfulIds []uint
	var deletedLevel []map[string]interface{}
	for _, id := range ids {
		level, err := srdls.repo.GetByID(id)
		if err != nil {
			fmt.Print(err)
			continue
		}

		err = srdls.repo.Delete(id)
		if err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}
		deletedLevel = append(deletedLevel, map[string]interface{}{
			"id":       id,
			"title_en": level.TitleEn,
			"title_ar": level.TitleAr,
		})
	}

	if len(successfulIds) > 0 {
		go func() {
			err := srdls.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.SystemRegisterDisruptionLevel{}.TableName(),
				successfulIds,
				deletedLevel,
				currentUserID,
				nil,
			)
			if err != nil {
				fmt.Printf("Failed to log bulk delete activity: %v\n", err)
			}
		}()
	}

	return failedIds, nil
}

func (srdls *SystemRegisterDisruptionLevelService) Activate(ids []uint, currentUserID uint) error {
	successfulIds, err := srdls.repo.Activate(ids)

	if len(successfulIds) > 0 {
		go func() {
			err := srdls.activityLogSvc.LogBulkActivities(
				models.ActivateAction,
				models.SystemRegisterDisruptionLevel{}.TableName(),
				successfulIds,
				nil,
				currentUserID,
				nil,
			)
			if err != nil {
				fmt.Printf("Failed to log bulk delete activity: %v\n", err)
			}
		}()
	}

	return err
}

func (srdls *SystemRegisterDisruptionLevelService) Deactivate(ids []uint, currentUserID uint) error {
	successfulIds, err := srdls.repo.Deactivate(ids)

	if len(successfulIds) > 0 {
		go func() {
			err := srdls.activityLogSvc.LogBulkActivities(
				models.DeactivateAction,
				models.SystemRegisterDisruptionLevel{}.TableName(),
				successfulIds,
				nil,
				currentUserID,
				nil,
			)
			if err != nil {
				fmt.Printf("Failed to log bulk delete activity: %v\n", err)
			}
		}()
	}

	return err
}
