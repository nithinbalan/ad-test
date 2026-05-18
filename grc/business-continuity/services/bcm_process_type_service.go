package services

import (
	"errors"

	bcmDTOs "jethur/grc/business-continuity/dtos"
	"jethur/grc/business-continuity/repositories"
	activityLogServices "jethur/grc/logs/services"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type BCMProcessTypeServiceInterface interface {
	GetAll(params *bcmDTOs.BCMProcessTypeListQueryParamsDTO) ([]*bcmDTOs.BCMProcessTypeDTO, error)
	GetByID(id uint) (*models.BCMProcessType, error)
	Create(currentUserID uint, formDTO *bcmDTOs.BCMProcessTypeFormDTO) (uint, error)
	Update(id uint, currentUserID uint, formDTO *bcmDTOs.BCMProcessTypeFormDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	Activate(ids []uint, currentUserID uint) error
	Deactivate(ids []uint, currentUserID uint) error
}

// Service and implementations
type BCMProcessTypeService struct {
	repo           repositories.BCMProcessTypeRepositoryInterface
	activityLogSvc activityLogServices.ActivityLogServiceInterface
}

func NewBCMProcessTypeService(repo repositories.BCMProcessTypeRepositoryInterface, activityLogSvc activityLogServices.ActivityLogServiceInterface) BCMProcessTypeServiceInterface {
	return &BCMProcessTypeService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
	}
}

func (service *BCMProcessTypeService) GetAll(params *bcmDTOs.BCMProcessTypeListQueryParamsDTO) ([]*bcmDTOs.BCMProcessTypeDTO, error) {
	return service.repo.GetAll(params)
}

func (service *BCMProcessTypeService) GetByID(id uint) (*models.BCMProcessType, error) {
	tool, err := service.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return tool, nil
}

func (service *BCMProcessTypeService) Create(currentUserID uint, formDTO *bcmDTOs.BCMProcessTypeFormDTO) (uint, error) {
	tool := &models.BCMProcessType{
		TitleEn:     formDTO.TitleEn,
		TitleAr:     formDTO.TitleAr,
		Status:      formDTO.Status,
		CreatedBYID: currentUserID,
	}

	id, err := service.repo.Create(tool)
	if err != nil {
		return 0, err
	}

	go func() {
		createdTool, err := service.GetByID(id)
		if err != nil {
			return
		}

		err = service.activityLogSvc.LogActivity(
			models.CreateAction,
			models.BCMProcessType{}.TableName(),
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
				"tableName": models.BCMProcessType{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return id, nil
}

func (service *BCMProcessTypeService) Update(id uint, currentUserID uint, formDTO *bcmDTOs.BCMProcessTypeFormDTO) error {
	existingTool, err := service.GetByID(id)
	if err != nil {
		return err
	}

	oldTool := &models.BCMProcessType{
		TitleEn: existingTool.TitleEn,
		TitleAr: existingTool.TitleAr,
		Status:  existingTool.Status,
	}

	existingTool.TitleEn = formDTO.TitleEn
	existingTool.TitleAr = formDTO.TitleAr
	existingTool.Status = formDTO.Status
	existingTool.UpdatedBYID = &currentUserID

	err = service.repo.Update(existingTool)
	if err != nil {
		return err
	}

	fieldsToCompare := []string{"TitleEn", "TitleAr", "Status"}
	changes, err := utility.CompareFields(oldTool, existingTool, fieldsToCompare)
	if err != nil {
		fields := map[string]interface{}{
			"error": err.Error(),
		}
		utility.LogError("Failed to compare fields", fields)
	}

	changedDataJSON, err := utility.ChangesToJSON(changes)
	if err != nil {
		fields := map[string]interface{}{
			"error": err.Error(),
		}
		utility.LogError("Failed to convert changes to JSON", fields)
	}

	go func() {
		err := service.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.BCMProcessType{}.TableName(),
			id,
			changedDataJSON,
			currentUserID,
			"",
			"",
			"",
			nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"error": err.Error(),
			}
			utility.LogError("Failed to log activity", fields)
		}
	}()

	return nil
}

func (service *BCMProcessTypeService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIds []uint
	var successfulIds []uint
	var deletedTool []map[string]interface{}
	for _, id := range ids {
		tool, err := service.repo.GetByID(id)
		if err != nil {
			failedIds = append(failedIds, id)
			continue
		}

		err = service.repo.Delete(id)
		if err != nil {
			failedIds = append(failedIds, id)
			continue
		}

		successfulIds = append(successfulIds, id)
		deletedTool = append(deletedTool, map[string]interface{}{
			"id":       id,
			"title_en": tool.TitleEn,
			"title_ar": tool.TitleAr,
		})
	}

	if len(successfulIds) > 0 {
		go func() {
			err := service.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.BCMProcessType{}.TableName(),
				successfulIds,
				deletedTool,
				currentUserID,
				nil,
			)
			if err != nil {
				fields := map[string]interface{}{
					"error": err.Error(),
				}
				utility.LogError("Failed to log bulk delete activity", fields)
			}
		}()
	}

	return failedIds, nil
}

func (service *BCMProcessTypeService) Activate(ids []uint, currentUserID uint) error {
	successfulIds, err := service.repo.Activate(ids)

	if len(successfulIds) > 0 {
		go func() {
			err := service.activityLogSvc.LogBulkActivities(
				models.ActivateAction,
				models.BCMProcessType{}.TableName(),
				successfulIds,
				nil,
				currentUserID,
				nil,
			)
			if err != nil {
				fields := map[string]interface{}{
					"error": err.Error(),
				}
				utility.LogError("Failed to log bulk activate activity", fields)
			}
		}()
	}

	return err
}

func (service *BCMProcessTypeService) Deactivate(ids []uint, currentUserID uint) error {
	successfulIds, err := service.repo.Deactivate(ids)

	if len(successfulIds) > 0 {
		go func() {
			err := service.activityLogSvc.LogBulkActivities(
				models.DeactivateAction,
				models.BCMProcessType{}.TableName(),
				successfulIds,
				nil,
				currentUserID,
				nil,
			)
			if err != nil {
				fields := map[string]interface{}{
					"error": err.Error(),
				}
				utility.LogError("Failed to log bulk deactivate activity", fields)
			}
		}()
	}

	return err
}
