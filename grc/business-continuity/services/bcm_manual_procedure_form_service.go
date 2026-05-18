package services

import (
	"errors"
	"fmt"
	"log"

	bcmDTOs "jethur/grc/business-continuity/dtos"
	"jethur/grc/business-continuity/repositories"
	activityLogServices "jethur/grc/logs/services"
	"jethur/grc/models"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type BCMManualProcedureFormServiceInterface interface {
	GetAll(params *bcmDTOs.BCMManualProcedureFormListQueryParamsDTO) ([]*bcmDTOs.BCMManualProcedureFormDTO, error)
	GetByID(id uint) (*models.BCMManualProcedureForm, error)
	GetFullDetailsByID(id uint) (*bcmDTOs.BCMManualProcedureFormDetailDTO, error)
	Create(currentUserID uint, formDTO *bcmDTOs.BCMManualProcedureFormFormDTO) (uint, error)
	Update(id uint, currentUserID uint, formDTO *bcmDTOs.BCMManualProcedureFormFormDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)
}

// Service and implementations
type BCMManualProcedureFormService struct {
	repo           repositories.BCMManualProcedureFormRepositoryInterface
	activityLogSvc activityLogServices.ActivityLogServiceInterface
}

func NewBCMManualProcedureFormService(repo repositories.BCMManualProcedureFormRepositoryInterface, activityLogSvc activityLogServices.ActivityLogServiceInterface) BCMManualProcedureFormServiceInterface {
	return &BCMManualProcedureFormService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
	}
}

func (service *BCMManualProcedureFormService) GetAll(params *bcmDTOs.BCMManualProcedureFormListQueryParamsDTO) ([]*bcmDTOs.BCMManualProcedureFormDTO, error) {
	return service.repo.GetAll(params)
}

func (service BCMManualProcedureFormService) GetByID(id uint) (*models.BCMManualProcedureForm, error) {
	form, err := service.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return form, nil
}

func (service *BCMManualProcedureFormService) GetFullDetailsByID(id uint) (*bcmDTOs.BCMManualProcedureFormDetailDTO, error) {
	var detailDTO bcmDTOs.BCMManualProcedureFormDetailDTO

	form, err := service.repo.GetFullDetailsByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	if err := mapstructure.Decode(form, &detailDTO); err != nil {
		return nil, err
	}

	detailDTO.CreatedAt = form.CreatedAt
	detailDTO.UpdatedAt = form.UpdatedAt

	return &detailDTO, nil
}

func (service *BCMManualProcedureFormService) Create(currentUserID uint, formDTO *bcmDTOs.BCMManualProcedureFormFormDTO) (uint, error) {
	form := &models.BCMManualProcedureForm{
		TitleEn:       formDTO.TitleEn,
		TitleAr:       formDTO.TitleAr,
		Code:          formDTO.Code,
		OwnerID:       formDTO.OwnerID,
		LocationID:    formDTO.LocationID,
		StakeholderID: formDTO.StakeholderID,
		Notes:         formDTO.Notes,
		CreatedBYID:   currentUserID,
	}

	id, err := service.repo.Create(form)
	if err != nil {
		return 0, err
	}

	go func() {
		createdForm, err := service.GetByID(id)
		if err != nil {
			return
		}

		err = service.activityLogSvc.LogActivity(
			models.CreateAction,
			models.BCMManualProcedureForm{}.TableName(),
			id,
			createdForm,
			currentUserID,
			"",
			"",
			"",
			nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.BCMManualProcedureForm{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return id, nil
}

func (service *BCMManualProcedureFormService) Update(id uint, currentUserID uint, formDTO *bcmDTOs.BCMManualProcedureFormFormDTO) error {
	existingForm, err := service.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utility.NewResourceNotFoundError(err, id)
		}
		return err
	}

	oldForm := &models.BCMManualProcedureForm{
		TitleEn:       existingForm.TitleEn,
		TitleAr:       existingForm.TitleAr,
		Code:          existingForm.Code,
		OwnerID:       existingForm.OwnerID,
		LocationID:    existingForm.LocationID,
		StakeholderID: existingForm.StakeholderID,
		Notes:         existingForm.Notes,
	}

	existingForm.TitleEn = formDTO.TitleEn
	existingForm.TitleAr = formDTO.TitleAr
	existingForm.Code = formDTO.Code
	existingForm.OwnerID = formDTO.OwnerID
	existingForm.LocationID = formDTO.LocationID
	existingForm.StakeholderID = formDTO.StakeholderID
	existingForm.Notes = formDTO.Notes
	existingForm.UpdatedBYID = &currentUserID

	if err := service.repo.Update(existingForm); err != nil {
		return err
	}

	go func() {
		fieldsToCompare := []string{"TitleEn", "TitleAr", "Code", "OwnerID", "LocationID", "StakeholderID", "Notes"}
		changes, err := utility.CompareFields(oldForm, existingForm, fieldsToCompare)
		if err != nil {
			log.Println("error :", err.Error())
			return
		}

		changedDataJSON, err := utility.ChangesToJSON(changes)
		if err != nil {
			log.Println("error :", err.Error())
			return
		}

		err = service.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.BCMManualProcedureForm{}.TableName(),
			id,
			changedDataJSON,
			currentUserID,
			"",
			"",
			"",
			nil,
		)
		if err != nil {
			fmt.Printf("Failed to log activity: %v\n", err)
		}
	}()

	return nil
}

func (service *BCMManualProcedureFormService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIDs []uint
	var successfulIDs []uint
	var deletedForms []map[string]interface{}

	for _, id := range ids {
		form, err := service.repo.GetByID(id)
		if err != nil {
			failedIDs = append(failedIDs, id)
			continue
		}

		if err := service.repo.Delete(id); err != nil {
			failedIDs = append(failedIDs, id)
			continue
		}

		successfulIDs = append(successfulIDs, id)
		deletedForms = append(deletedForms, map[string]interface{}{
			"id":      id,
			"titleEn": form.TitleEn,
			"titleAr": form.TitleAr,
		})
	}

	if len(successfulIDs) > 0 {
		go func() {
			err := service.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.BCMManualProcedureForm{}.TableName(),
				successfulIDs,
				deletedForms,
				currentUserID,
				nil,
			)
			if err != nil {
				fmt.Printf("Failed to log bulk delete activity: %v\n", err)
			}
		}()
	}

	return failedIDs, nil
}
