package services

import (
	"errors"
	"fmt"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	settingsDTOs "jethur/grc/settings/dtos"
	settingsRepositories "jethur/grc/settings/repositories"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type OrganizationUnitServiceInterface interface {
	GetAll(params *settingsDTOs.OrganizationUnitListQueryParamsDTO) ([]*settingsDTOs.OrganizationUnitListDTO, error)
	GetByID(id uint) (*models.OrganizationUnit, error)
	GetFullDetailsByID(id uint) (*settingsDTOs.OrganizationUnitDTO, error)
	Create(currentUserID uint, orgUnitFormDTO *settingsDTOs.OrganizationUnitFormDTO) (uint, error)
	Update(id uint, currentUserID uint, orgUnit *settingsDTOs.OrganizationUnitFormDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	Activate(ids []uint, currentUserID uint) error
	Deactivate(ids []uint, currentUserID uint) error
}

type OrganizationUnitService struct {
	repo                 settingsRepositories.OrganizationUnitRepositoryInterface
	activityLogSvc       logServices.ActivityLogServiceInterface
	organizationLevelSvc OrganizationLevelServiceInterface
}

func NewOrganizationUnitService(
	repo settingsRepositories.OrganizationUnitRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	organizationLevelSvc OrganizationLevelServiceInterface) OrganizationUnitServiceInterface {
	return &OrganizationUnitService{
		repo:                 repo,
		activityLogSvc:       activityLogSvc,
		organizationLevelSvc: organizationLevelSvc,
	}
}

func (ous *OrganizationUnitService) GetAll(params *settingsDTOs.OrganizationUnitListQueryParamsDTO) ([]*settingsDTOs.OrganizationUnitListDTO, error) {
	return ous.repo.GetAll(params)
}

func (ous OrganizationUnitService) GetByID(id uint) (*models.OrganizationUnit, error) {
	orgUnit, err := ous.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return orgUnit, nil
}

func (ous OrganizationUnitService) GetFullDetailsByID(unitId uint) (*settingsDTOs.OrganizationUnitDTO, error) {
	// Initialize your DTO
	var unitDTO settingsDTOs.OrganizationUnitDTO

	unit, err := ous.repo.GetFullDetailsByID(unitId)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, unitId)
		}
		return nil, err
	}

	// Map the result to your DTO
	err = mapstructure.Decode(unit, &unitDTO)
	if err != nil {
		return nil, err
	}

	return &unitDTO, nil
}

func (ous *OrganizationUnitService) Create(currentUserID uint, orgUnitFormDTO *settingsDTOs.OrganizationUnitFormDTO) (uint, error) {

	if err := ous.validateOrganizationUnitFormDTO(orgUnitFormDTO); err != nil {
		return 0, err
	}

	colorCode := utility.GenerateColorCode("dept")

	organizationUnit := &models.OrganizationUnit{
		OrganizationLevelID: orgUnitFormDTO.OrganizationLevelID,
		OrganizationUnitID:  orgUnitFormDTO.OrganizationUnitID,
		TitleEn:             orgUnitFormDTO.TitleEn,
		TitleAr:             orgUnitFormDTO.TitleAr,
		Status:              orgUnitFormDTO.Status,
		ColorCode:           colorCode,
	}

	id, err := ous.repo.Create(organizationUnit)
	if err != nil {
		return 0, err
	}

	go func() {
		err := ous.activityLogSvc.LogActivity(
			models.CreateAction,
			models.OrganizationUnit{}.TableName(),
			id,
			orgUnitFormDTO,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.OrganizationUnit{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return id, nil
}

func (ous *OrganizationUnitService) Update(id uint, currentUserID uint, orgUnitFormDTO *settingsDTOs.OrganizationUnitFormDTO) error {
	existingOrgUnit, err := ous.GetByID(id)
	if err != nil {
		return err
	}

	existingOrgUnit.OrganizationLevelID = orgUnitFormDTO.OrganizationLevelID
	existingOrgUnit.OrganizationUnitID = orgUnitFormDTO.OrganizationUnitID
	existingOrgUnit.TitleEn = orgUnitFormDTO.TitleEn
	existingOrgUnit.TitleAr = orgUnitFormDTO.TitleAr

	err = ous.repo.Update(existingOrgUnit)
	if err != nil {
		return err
	}

	// Log the update activity
	go func() {
		err := ous.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.OrganizationUnit{}.TableName(),
			id,
			orgUnitFormDTO,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.UpdateAction,
				"tableName": models.OrganizationUnit{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return nil
}

func (ous *OrganizationUnitService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIds []uint
	var successfulIds []uint

	for _, id := range ids {
		_, err := ous.GetByID(id)
		if err != nil {
			failedIds = append(failedIds, id)
			continue
		}
		// Delete the organization unit
		err = ous.repo.Delete(id)
		if err != nil {
			failedIds = append(failedIds, id)
			continue
		}
		successfulIds = append(successfulIds, id)
	}

	// Log the delete activity for successful deletions
	go func() {
		err := ous.activityLogSvc.LogBulkActivities(
			models.DeleteAction,
			models.OrganizationUnit{}.TableName(),
			successfulIds,
			nil,
			currentUserID,
			nil,
		)
		if err != nil {
			fmt.Printf("Failed to log bulk delete activity: %v\n", err)
		}
	}()

	return failedIds, nil
}

func (ous *OrganizationUnitService) Activate(ids []uint, currentUserID uint) error {
	successfulIds, err := ous.repo.Activate(ids)

	if len(successfulIds) > 0 {

		go func() {
			err := ous.activityLogSvc.LogBulkActivities(
				models.ActivateAction,
				models.Role{}.TableName(),
				successfulIds,
				nil,
				currentUserID,
				nil,
			)
			if err != nil {
				// @todo - save to log file
				fmt.Printf("Failed to log activity: %v\n", err)
			}
		}()
	}

	return err
}

func (ous *OrganizationUnitService) Deactivate(ids []uint, currentUserID uint) error {
	successfulIds, err := ous.repo.Deactivate(ids)

	if len(successfulIds) > 0 {

		go func() {
			err := ous.activityLogSvc.LogBulkActivities(
				models.DeactivateAction,
				models.Role{}.TableName(),
				successfulIds,
				nil,
				currentUserID,
				nil,
			)
			if err != nil {
				// @todo - save to log file
				// You might want to log this error or handle it appropriately
				fmt.Printf("Failed to log activity: %v\n", err)
			}
		}()
	}

	return err
}

func (ous *OrganizationUnitService) validateOrganizationUnitFormDTO(orgUnitFormDTO *settingsDTOs.OrganizationUnitFormDTO) error {
	if _, err := ous.organizationLevelSvc.GetByID(orgUnitFormDTO.OrganizationLevelID); err != nil {
		return utility.NewValidationError(fmt.Errorf("organization level does not exist"), "not_exist", "organization level")
	}

	if orgUnitFormDTO.OrganizationUnitID != nil {
		if _, err := ous.GetByID(*orgUnitFormDTO.OrganizationUnitID); err != nil {
			return utility.NewValidationError(fmt.Errorf("organization unit does not exist"), "not_exist", "organization unit")
		}
	}

	return nil
}
