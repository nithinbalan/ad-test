package services

import (
	"errors"
	"fmt"
	"jethur/grc/asset/dtos"
	"jethur/grc/asset/repositories"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type AssetMaintenanceServiceInterface interface {
	GetAll(params *dtos.AssetMaintenanceListQueryParamsDTO, assetId uint) (*dtos.AssetMaintenanceListResponseDTO, error)
	GetByID(id uint) (*models.AssetMaintenance, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*dtos.AssetMaintenanceDetailDTO, error)
	Create(maintenanceForm *dtos.AssetMaintenanceFormDTO, currentUserID uint, assetId uint) (uint, error)
	Update(id uint, maintenanceForm *dtos.AssetMaintenanceFormDTO, currentUserID uint) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)
}

type AssetMaintenanceService struct {
	repo             repositories.AssetMaintenanceRepositoryInterface
	activityLogSvc   logServices.ActivityLogServiceInterface
	assetRegisterSvc AssetRegisterServiceInterface
	userSvc          organizationServices.UserServiceInterface
}

func NewAssetMaintenanceService(
	repo repositories.AssetMaintenanceRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	assetRegisterSvc AssetRegisterServiceInterface,
	userSvc organizationServices.UserServiceInterface,
) AssetMaintenanceServiceInterface {
	return &AssetMaintenanceService{
		repo:             repo,
		activityLogSvc:   activityLogSvc,
		assetRegisterSvc: assetRegisterSvc,
		userSvc:          userSvc,
	}
}

func (ams *AssetMaintenanceService) GetAll(params *dtos.AssetMaintenanceListQueryParamsDTO, assetId uint) (*dtos.AssetMaintenanceListResponseDTO, error) {
	return ams.repo.GetAll(params, assetId)
}

func (ams *AssetMaintenanceService) GetByID(id uint) (*models.AssetMaintenance, error) {
	maintenance, err := ams.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return maintenance, nil
}

func (ams *AssetMaintenanceService) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*dtos.AssetMaintenanceDetailDTO, error) {
	var maintenanceDetailDTO dtos.AssetMaintenanceDetailDTO

	maintenance, err := ams.repo.GetFullDetailsByID(id, tx...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	err = mapstructure.Decode(maintenance, &maintenanceDetailDTO)
	if err != nil {
		fmt.Printf("Decode error: %v\n", err)
		return nil, err
	}

	performedDateString := maintenance.PerformedDate.Format("2006-01-02")
	maintenanceDetailDTO.PerformedDateString = performedDateString
	return &maintenanceDetailDTO, nil
}

func (ams *AssetMaintenanceService) Create(maintenanceForm *dtos.AssetMaintenanceFormDTO, currentUserID uint, assetId uint) (uint, error) {
	tx := ams.repo.BeginTransaction()

	if err := ams.validateAssetMaintenanceFormDTO(maintenanceForm); err != nil {
		return 0, err
	}

	dates, err := ams.parseMaintenanceDates(maintenanceForm)
	if err != nil {
		return 0, err
	}

	maintenance := &models.AssetMaintenance{
		AssetRegisterID: assetId,
		Type:            maintenanceForm.Type,
		PerformedDate:   *dates.PerformedDate,
		Title:           maintenanceForm.Title,
		Description:     maintenanceForm.Description,
		PerformedBYID:   maintenanceForm.PerformedByID,
		Outcome:         maintenanceForm.Outcome,
		CreatedBYID:     currentUserID,
	}

	id, err := ams.repo.Create(maintenance, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	tx.Commit()

	go func() {
		err := ams.activityLogSvc.LogActivity(
			models.CreateAction,
			models.AssetMaintenance{}.TableName(),
			id,
			maintenanceForm,
			currentUserID,
			nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.AssetMaintenance{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}
			utility.LogError("Failed to log maintenance activity", fields)
		}
	}()

	return id, nil
}

func (ams *AssetMaintenanceService) Update(id uint, maintenanceForm *dtos.AssetMaintenanceFormDTO, currentUserID uint) error {
	tx := ams.repo.BeginTransaction()

	existingMaintenance, err := ams.repo.GetByID(id)
	if err != nil {
		return err
	}

	if err := ams.validateAssetMaintenanceFormDTO(maintenanceForm); err != nil {
		return err
	}

	dates, err := ams.parseMaintenanceDates(maintenanceForm)
	if err != nil {
		return err
	}

	existingMaintenance.Type = maintenanceForm.Type
	existingMaintenance.PerformedDate = *dates.PerformedDate
	existingMaintenance.Title = maintenanceForm.Title
	existingMaintenance.Description = maintenanceForm.Description
	existingMaintenance.PerformedBYID = maintenanceForm.PerformedByID
	existingMaintenance.Outcome = maintenanceForm.Outcome
	existingMaintenance.UpdatedBYID = &currentUserID

	err = ams.repo.Update(existingMaintenance, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (ams *AssetMaintenanceService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIds []uint
	var successfulIds []uint

	for _, id := range ids {
		err := ams.repo.Delete(id)
		if err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}
	}

	if len(successfulIds) > 0 {
		go func() {
			err := ams.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.AssetMaintenance{}.TableName(),
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

	return failedIds, nil
}

func (ams *AssetMaintenanceService) validateAssetMaintenanceFormDTO(maintenanceForm *dtos.AssetMaintenanceFormDTO) error {
	if _, err := ams.userSvc.GetByID(uint(maintenanceForm.PerformedByID)); err != nil {
		return utility.NewValidationError(fmt.Errorf("performed by user does not exist"), "not_exist", "performed by user")
	}

	return nil
}

func (ams *AssetMaintenanceService) parseMaintenanceDates(maintenanceForm *dtos.AssetMaintenanceFormDTO) (*dtos.AssetMaintenanceParsedDates, error) {
	parsed := &dtos.AssetMaintenanceParsedDates{}
	var err error

	if maintenanceForm.PerformedDate != "" {
		parsed.PerformedDate, err = utility.ParseDate(maintenanceForm.PerformedDate)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid performed date"), "invalid_date", "performed date")
		}
	}

	return parsed, nil
}
