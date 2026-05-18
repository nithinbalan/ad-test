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

type BCMExecutionRecordServiceInterface interface {
	GetAll(params *bcmDTOs.BCMExecutionRecordListQueryParamsDTO) ([]*bcmDTOs.BCMExecutionRecordDTO, error)
	GetByID(id uint) (*models.BCMExecutionRecord, error)
	GetFullDetailsByID(id uint) (*bcmDTOs.BCMExecutionRecordDetailDTO, error)
	Create(currentUserID uint, formDTO *bcmDTOs.BCMExecutionRecordFormDTO) (uint, error)
	Update(id uint, currentUserID uint, formDTO *bcmDTOs.BCMExecutionRecordFormDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)
}

// Service and implementations
type BCMExecutionRecordService struct {
	repo           repositories.BCMExecutionRecordRepositoryInterface
	activityLogSvc activityLogServices.ActivityLogServiceInterface
}

func NewBCMExecutionRecordService(repo repositories.BCMExecutionRecordRepositoryInterface, activityLogSvc activityLogServices.ActivityLogServiceInterface) BCMExecutionRecordServiceInterface {
	return &BCMExecutionRecordService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
	}
}

func (service *BCMExecutionRecordService) GetAll(params *bcmDTOs.BCMExecutionRecordListQueryParamsDTO) ([]*bcmDTOs.BCMExecutionRecordDTO, error) {
	return service.repo.GetAll(params)
}

func (service BCMExecutionRecordService) GetByID(id uint) (*models.BCMExecutionRecord, error) {
	record, err := service.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return record, nil
}

func (service *BCMExecutionRecordService) GetFullDetailsByID(id uint) (*bcmDTOs.BCMExecutionRecordDetailDTO, error) {
	var detailDTO bcmDTOs.BCMExecutionRecordDetailDTO

	record, err := service.repo.GetFullDetailsByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	if err := mapstructure.Decode(record, &detailDTO); err != nil {
		return nil, err
	}

	detailDTO.CreatedAt = record.CreatedAt
	detailDTO.UpdatedAt = record.UpdatedAt

	return &detailDTO, nil
}

func (service *BCMExecutionRecordService) Create(currentUserID uint, formDTO *bcmDTOs.BCMExecutionRecordFormDTO) (uint, error) {
	record := &models.BCMExecutionRecord{
		TitleEn:           formDTO.TitleEn,
		TitleAr:           formDTO.TitleAr,
		OwnerEntityID:     formDTO.OwnerEntityID,
		RecordTypeID:      formDTO.RecordTypeID,
		IsBackUpAvailable: formDTO.IsBackUpAvailable,
		LocationID:        formDTO.LocationID,
		RecoveryPlan:      formDTO.RecoveryPlan,
		CreatedBYID:       currentUserID,
	}

	id, err := service.repo.Create(record)
	if err != nil {
		return 0, err
	}

	go func() {
		createdRecord, err := service.GetByID(id)
		if err != nil {
			return
		}

		err = service.activityLogSvc.LogActivity(
			models.CreateAction,
			models.BCMExecutionRecord{}.TableName(),
			id,
			createdRecord,
			currentUserID,
			"",
			"",
			"",
			nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.BCMExecutionRecord{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return id, nil
}

func (service *BCMExecutionRecordService) Update(id uint, currentUserID uint, formDTO *bcmDTOs.BCMExecutionRecordFormDTO) error {
	existingRecord, err := service.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utility.NewResourceNotFoundError(err, id)
		}
		return err
	}

	oldRecord := &models.BCMExecutionRecord{
		TitleEn:           existingRecord.TitleEn,
		TitleAr:           existingRecord.TitleAr,
		OwnerEntityID:     existingRecord.OwnerEntityID,
		RecordTypeID:      existingRecord.RecordTypeID,
		IsBackUpAvailable: existingRecord.IsBackUpAvailable,
		LocationID:        existingRecord.LocationID,
		RecoveryPlan:      existingRecord.RecoveryPlan,
	}

	existingRecord.TitleEn = formDTO.TitleEn
	existingRecord.TitleAr = formDTO.TitleAr
	existingRecord.OwnerEntityID = formDTO.OwnerEntityID
	existingRecord.RecordTypeID = formDTO.RecordTypeID
	existingRecord.IsBackUpAvailable = formDTO.IsBackUpAvailable
	existingRecord.LocationID = formDTO.LocationID
	existingRecord.RecoveryPlan = formDTO.RecoveryPlan
	existingRecord.UpdatedBYID = &currentUserID

	if err := service.repo.Update(existingRecord); err != nil {
		return err
	}

	go func() {
		fieldsToCompare := []string{"TitleEn", "TitleAr", "OwnerEntityID", "RecordTypeID", "IsBackUpAvailable", "LocationID", "RecoveryPlan"}
		changes, err := utility.CompareFields(oldRecord, existingRecord, fieldsToCompare)
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
			models.BCMExecutionRecord{}.TableName(),
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

func (service *BCMExecutionRecordService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIDs []uint
	var successfulIDs []uint
	var deletedRecords []map[string]interface{}

	for _, id := range ids {
		record, err := service.repo.GetByID(id)
		if err != nil {
			failedIDs = append(failedIDs, id)
			continue
		}

		if err := service.repo.Delete(id); err != nil {
			failedIDs = append(failedIDs, id)
			continue
		}

		successfulIDs = append(successfulIDs, id)
		deletedRecords = append(deletedRecords, map[string]interface{}{
			"id":      id,
			"titleEn": record.TitleEn,
			"titleAr": record.TitleAr,
		})
	}

	if len(successfulIDs) > 0 {
		go func() {
			err := service.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.BCMExecutionRecord{}.TableName(),
				successfulIDs,
				deletedRecords,
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
