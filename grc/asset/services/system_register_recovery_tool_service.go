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

type SystemRegisterRecoveryToolServiceInterface interface {
	GetAll(params *assetDTOs.SystemRegisterRecoveryToolListQueryParamsDTO) ([]*assetDTOs.SystemRegisterRecoveryToolDTO, error)
	GetByID(id uint) (*models.SystemRegisterRecoveryTool, error)
	Create(currentUserID uint, toolFormDTO *assetDTOs.SystemRegisterRecoveryToolFormDTO) (uint, error)
	Update(id uint, currentUserID uint, toolFormDTO *assetDTOs.SystemRegisterRecoveryToolFormDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	Activate(ids []uint, currentUserID uint) error
	Deactivate(ids []uint, currentUserID uint) error
}

// Service and implementations
type SystemRegisterRecoveryToolService struct {
	repo           repositories.SystemRegisterRecoveryToolRepositoryInterface
	activityLogSvc activityLogServices.ActivityLogServiceInterface
}

func NewSystemRegisterRecoveryToolService(repo repositories.SystemRegisterRecoveryToolRepositoryInterface, activityLogSvc activityLogServices.ActivityLogServiceInterface) SystemRegisterRecoveryToolServiceInterface {
	return &SystemRegisterRecoveryToolService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
	}
}

func (srrts *SystemRegisterRecoveryToolService) GetAll(params *assetDTOs.SystemRegisterRecoveryToolListQueryParamsDTO) ([]*assetDTOs.SystemRegisterRecoveryToolDTO, error) {
	return srrts.repo.GetAll(params)
}

func (srrts SystemRegisterRecoveryToolService) GetByID(id uint) (*models.SystemRegisterRecoveryTool, error) {
	tool, err := srrts.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return tool, nil
}

func (srrts *SystemRegisterRecoveryToolService) Create(currentUserID uint, toolFormDTO *assetDTOs.SystemRegisterRecoveryToolFormDTO) (uint, error) {
	tool := &models.SystemRegisterRecoveryTool{
		TitleEn:     toolFormDTO.TitleEn,
		TitleAr:     toolFormDTO.TitleAr,
		Status:      toolFormDTO.Status,
		CreatedBYID: currentUserID,
	}

	id, err := srrts.repo.Create(tool)
	if err != nil {
		return 0, err
	}

	go func() {
		createdTool, err := srrts.GetByID(id)
		if err != nil {
			return
		}

		err = srrts.activityLogSvc.LogActivity(
			models.CreateAction,
			models.SystemRegisterRecoveryTool{}.TableName(),
			id,
			createdTool,
			currentUserID,
			"",
			"",
			"",
			nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.SystemRegisterRecoveryTool{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return id, nil
}

func (srrts *SystemRegisterRecoveryToolService) Update(id uint, currentUserID uint, toolFormDTO *assetDTOs.SystemRegisterRecoveryToolFormDTO) error {
	oldTool, err := srrts.GetByID(id)
	if err != nil {
		return err
	}

	existingTool, err := srrts.GetByID(id)
	if err != nil {
		return err
	}

	existingTool.TitleEn = toolFormDTO.TitleEn
	existingTool.TitleAr = toolFormDTO.TitleAr
	existingTool.Status = toolFormDTO.Status
	existingTool.UpdatedBYID = &currentUserID

	err = srrts.repo.Update(existingTool)
	if err != nil {
		return err
	}

	fieldsToCompare := []string{"TitleEn", "TitleAr", "Status"}
	changes, err := utility.CompareFields(oldTool, existingTool, fieldsToCompare)
	if err != nil {
		log.Println("error :", err.Error())
	}

	changedDataJSON, err := utility.ChangesToJSON(changes)
	if err != nil {
		log.Println("error :", err.Error())
	}

	go func() {
		err := srrts.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.SystemRegisterRecoveryTool{}.TableName(),
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

func (srrts *SystemRegisterRecoveryToolService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIds []uint
	var successfulIds []uint
	var deletedTool []map[string]interface{}
	for _, id := range ids {
		tool, err := srrts.repo.GetByID(id)
		if err != nil {
			fmt.Print(err)
			continue
		}

		err = srrts.repo.Delete(id)
		if err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}
		deletedTool = append(deletedTool, map[string]interface{}{
			"id":       id,
			"title_en": tool.TitleEn,
			"title_ar": tool.TitleAr,
		})
	}

	if len(successfulIds) > 0 {
		go func() {
			err := srrts.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.SystemRegisterRecoveryTool{}.TableName(),
				successfulIds,
				deletedTool,
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

func (srrts *SystemRegisterRecoveryToolService) Activate(ids []uint, currentUserID uint) error {
	successfulIds, err := srrts.repo.Activate(ids)

	if len(successfulIds) > 0 {
		go func() {
			err := srrts.activityLogSvc.LogBulkActivities(
				models.ActivateAction,
				models.SystemRegisterRecoveryTool{}.TableName(),
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

func (srrts *SystemRegisterRecoveryToolService) Deactivate(ids []uint, currentUserID uint) error {
	successfulIds, err := srrts.repo.Deactivate(ids)

	if len(successfulIds) > 0 {
		go func() {
			err := srrts.activityLogSvc.LogBulkActivities(
				models.DeactivateAction,
				models.SystemRegisterRecoveryTool{}.TableName(),
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
