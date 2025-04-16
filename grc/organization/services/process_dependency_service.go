package services

import (
	"errors"
	"fmt"
	activityLogServices "jethur/grc/logs/services"
	"jethur/grc/models"
	orgDTOs "jethur/grc/organization/dtos"
	"jethur/grc/organization/repositories"
	settingServices "jethur/grc/settings/services"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type ProcessDependencyServiceInterface interface {
	GetAll(processID uint, params *orgDTOs.ProcessDependencyListQueryParamsDTO) ([]*orgDTOs.ProcessDependencyListDTO, error)
	GetByID(processID uint, id uint) (*models.ProcessDependency, error)
	Create(processID uint, currentUserID uint, processDependencyFormDTO *orgDTOs.ProcessDependencyFormDTO) (uint, error)
	Update(processID uint, id uint, currentUserID uint, processDependencyFormDTO *orgDTOs.ProcessDependencyUpdateFormDTO) error
	Delete(processID uint, ids []uint, currentUserID uint) ([]uint, error)

	UpdateProcessDependencyRecoveryObjective(processDependency *models.ProcessDependency, tx utility.DatabaseTransaction) error
}

type ProcessDependencyService struct {
	repo                        repositories.ProcessDependencyRepositoryInterface
	processDependencySubTypeSvc ProcessDependencySubTypeServiceInterface
	activityLogSvc              activityLogServices.ActivityLogServiceInterface
	processSvc                  ProcessServiceInterface
	organizationUnitSvc         settingServices.OrganizationUnitServiceInterface
}

func NewProcessDependencyService(
	repo repositories.ProcessDependencyRepositoryInterface,
	processDependencySubTypeSvc ProcessDependencySubTypeServiceInterface,
	activityLogSvc activityLogServices.ActivityLogServiceInterface,
	processSvc ProcessServiceInterface,
	organizationUnitSvc settingServices.OrganizationUnitServiceInterface,
) ProcessDependencyServiceInterface {
	return &ProcessDependencyService{
		repo:                        repo,
		processDependencySubTypeSvc: processDependencySubTypeSvc,
		activityLogSvc:              activityLogSvc,
		processSvc:                  processSvc,
		organizationUnitSvc:         organizationUnitSvc,
	}
}

func (pds *ProcessDependencyService) GetAll(processID uint, params *orgDTOs.ProcessDependencyListQueryParamsDTO) ([]*orgDTOs.ProcessDependencyListDTO, error) {
	if _, err := pds.processSvc.GetByID(processID); err != nil {
		return nil, err
	}
	return pds.repo.GetAll(processID, params)
}

func (pds *ProcessDependencyService) GetByID(processID uint, id uint) (*models.ProcessDependency, error) {
	if _, err := pds.processSvc.GetByID(processID); err != nil {
		return nil, err
	}
	processDependency, err := pds.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return processDependency, nil
}

func (pds *ProcessDependencyService) GetFullDetailsByID(processID uint, id uint) (*orgDTOs.ProcessDependencyDetailDTO, error) {
	var processDependencyDetailDTO orgDTOs.ProcessDependencyDetailDTO

	processDependency, err := pds.repo.GetFullDetailsByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	err = mapstructure.Decode(processDependency, &processDependencyDetailDTO)
	if err != nil {
		return nil, err
	}

	return &processDependencyDetailDTO, nil
}

func (pds *ProcessDependencyService) Create(processID uint, currentUserID uint, dto *orgDTOs.ProcessDependencyFormDTO) (uint, error) {
	if _, err := pds.processSvc.GetByID(processID); err != nil {
		return 0, err
	}

	if err := pds.validateProcessDependencyFormDTO(dto, nil); err != nil {
		return 0, err
	}

	var createdIDs []uint

	// Case 1: Dependent Process
	if len(dto.DependentProcessIDs) > 0 {
		for _, dependentProcessID := range dto.DependentProcessIDs {
			processDependency := &models.ProcessDependency{
				ProcessID:                  processID,
				ProcessDependencySubTypeID: dto.ProcessDependencySubTypeID,
				DependentProcessID:         dependentProcessID,
				Title:                      dto.Title,
				Description:                dto.Description,
				CreatedBYID:                currentUserID,
			}

			id, err := pds.repo.Create(processDependency)
			if err != nil {
				return 0, err
			}
			createdIDs = append(createdIDs, id)

			go pds.logDependencyCreation(processID, id, currentUserID)
		}
	}

	// Case 2: Organization Unit
	if len(dto.OrganizationUnitIDs) > 0 {
		for _, unitID := range dto.OrganizationUnitIDs {
			processDependency := &models.ProcessDependency{
				ProcessID:                  processID,
				ProcessDependencySubTypeID: dto.ProcessDependencySubTypeID,
				OrganizationUnitID:         unitID,
				Title:                      dto.Title,
				Description:                dto.Description,
				CreatedBYID:                currentUserID,
			}

			id, err := pds.repo.Create(processDependency)
			if err != nil {
				return 0, err
			}
			createdIDs = append(createdIDs, id)

			go pds.logDependencyCreation(processID, id, currentUserID)
		}
	}

	// Case 3: Asset
	if len(dto.AssetIDs) > 0 {
		for _, assetID := range dto.AssetIDs {
			processDependency := &models.ProcessDependency{
				ProcessID:                  processID,
				ProcessDependencySubTypeID: dto.ProcessDependencySubTypeID,
				AssetID:                    assetID,
				Title:                      dto.Title,
				Description:                dto.Description,
				CreatedBYID:                currentUserID,
			}

			id, err := pds.repo.Create(processDependency)
			if err != nil {
				return 0, err
			}
			createdIDs = append(createdIDs, id)

			go pds.logDependencyCreation(processID, id, currentUserID)
		}
	}

	// Case 4: Only Title & Description (no IDs)
	if len(dto.DependentProcessIDs) == 0 && len(dto.OrganizationUnitIDs) == 0 && len(dto.AssetIDs) == 0 {
		processDependency := &models.ProcessDependency{
			ProcessID:                  processID,
			ProcessDependencySubTypeID: dto.ProcessDependencySubTypeID,
			Title:                      dto.Title,
			Description:                dto.Description,
			CreatedBYID:                currentUserID,
		}

		id, err := pds.repo.Create(processDependency)
		if err != nil {
			return 0, err
		}
		createdIDs = append(createdIDs, id)

		go pds.logDependencyCreation(processID, id, currentUserID)
	}

	if len(createdIDs) > 0 {
		return createdIDs[0], nil
	}
	return 0, nil
}

func (pds *ProcessDependencyService) logDependencyCreation(processID, id, userID uint) {
	created, _ := pds.GetByID(processID, id)
	_ = pds.activityLogSvc.LogActivity(
		models.CreateAction,
		models.ProcessDependency{}.TableName(),
		id,
		created,
		userID,
		nil,
	)
}

func (pds *ProcessDependencyService) Update(processID uint, id uint, currentUserID uint, dto *orgDTOs.ProcessDependencyUpdateFormDTO) error {
	if _, err := pds.processSvc.GetByID(processID); err != nil {
		return err
	}

	oldData, err := pds.GetFullDetailsByID(processID, id)
	if err != nil {
		return err
	}

	existing, err := pds.GetByID(processID, id)
	if err != nil {
		return err
	}

	if err := pds.validateProcessDependencyUpdateFormDTO(dto, &id); err != nil {
		return err
	}

	existing.ProcessDependencySubTypeID = dto.ProcessDependencySubTypeID
	existing.DependentProcessID = dto.DependentProcessID
	existing.OrganizationUnitID = dto.OrganizationUnitID
	existing.AssetID = dto.AssetID
	existing.Title = dto.Title
	existing.Description = dto.Description
	existing.UpdatedBYID = &currentUserID

	err = pds.repo.Update(existing)
	if err != nil {
		return err
	}

	go func() {
		// Fetch updated full details after update
		newData, err := pds.GetFullDetailsByID(processID, id)
		if err != nil {
			return
		}

		oldFields := map[string]interface{}{
			"sub_type_title_en":          oldData.ProcessDependencySubType.TitleEn,
			"dependent_process_en":       utility.GetFieldValue(oldData.DependentProcess).(string),
			"organization_unit_title_en": utility.GetFieldValue(oldData.OrganizationUnit).(string),
			"asset_title_en":             utility.GetFieldValue(oldData.Asset).(string),
			"title":                      oldData.Title,
			"description":                oldData.Description,
		}

		newFields := map[string]interface{}{
			"sub_type_title_en":          newData.ProcessDependencySubType.TitleEn,
			"dependent_process_en":       utility.GetFieldValue(newData.DependentProcess).(string),
			"organization_unit_title_en": utility.GetFieldValue(newData.OrganizationUnit).(string),
			"asset_title_en":             utility.GetFieldValue(newData.Asset).(string),
			"title":                      newData.Title,
			"description":                newData.Description,
		}

		changedData := utility.CompareFieldChanges(oldFields, newFields)

		if len(changedData) > 0 {
			err := pds.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.ProcessDependency{}.TableName(),
				id,
				changedData,
				currentUserID,
				nil,
			)
			if err != nil {
				utility.LogError("Failed to log process dependency update", map[string]interface{}{
					"action":    models.UpdateAction,
					"table":     models.ProcessDependency{}.TableName(),
					"record_id": id,
					"user_id":   currentUserID,
					"error":     err.Error(),
				})
			}
		}
	}()

	return nil
}

func (pds *ProcessDependencyService) Delete(processID uint, ids []uint, currentUserID uint) ([]uint, error) {
	if _, err := pds.processSvc.GetByID(processID); err != nil {
		return nil, err
	}

	var failedIds []uint
	var successfulIds []uint
	var deletedItems []map[string]interface{}

	for _, id := range ids {
		item, err := pds.repo.GetByID(id)
		if err != nil {
			failedIds = append(failedIds, id)
			continue
		}
		if err := pds.repo.Delete(id); err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
			deletedItems = append(deletedItems, map[string]interface{}{
				"id":    id,
				"title": item.Title,
			})
		}
	}

	if len(successfulIds) > 0 {
		go pds.activityLogSvc.LogBulkActivities(
			models.DeleteAction,
			models.ProcessDependency{}.TableName(),
			successfulIds,
			deletedItems,
			currentUserID,
			nil,
		)
	}

	return failedIds, nil
}

func (pds *ProcessDependencyService) validateProcessDependencyFormDTO(processDependencyFormDTO *orgDTOs.ProcessDependencyFormDTO, id *uint) error {
	// Validate sub type
	if _, err := pds.processDependencySubTypeSvc.GetByID(processDependencyFormDTO.ProcessDependencySubTypeID); err != nil {
		return utility.NewValidationError(fmt.Errorf("dependency sub type does not exist"), "not_exist", "process dependency sub type")
	}

	// Validate each Dependent Process ID
	for _, processID := range processDependencyFormDTO.DependentProcessIDs {
		if processID != nil {
			if _, err := pds.processSvc.GetByID(*processID); err != nil {
				return utility.NewValidationError(fmt.Errorf("dependent process does not exist"), "not_exist", "dependent process")
			}
		}
	}

	// Validate each Organization Unit ID
	for _, unitID := range processDependencyFormDTO.OrganizationUnitIDs {
		if unitID != nil {
			if _, err := pds.organizationUnitSvc.GetByID(*unitID); err != nil {
				return utility.NewValidationError(fmt.Errorf("organization unit does not exist"), "not_exist", "organization unit")
			}
		}
	}

	// Validate each Asset ID
	for _, assetID := range processDependencyFormDTO.AssetIDs {
		if assetID != nil {
			if _, err := pds.GetAssetByID(*assetID); err != nil {
				return utility.NewValidationError(fmt.Errorf("asset register does not exist"), "not_exist", "asset register")
			}
		}
	}

	return nil
}

func (pds *ProcessDependencyService) validateProcessDependencyUpdateFormDTO(processDependencyFormDTO *orgDTOs.ProcessDependencyUpdateFormDTO, id *uint) error {
	if _, err := pds.processDependencySubTypeSvc.GetByID(processDependencyFormDTO.ProcessDependencySubTypeID); err != nil {
		return utility.NewValidationError(fmt.Errorf("dependency sub type does not exist"), "not_exist", "process dependency sub type")
	}

	if processDependencyFormDTO.DependentProcessID != nil {
		if _, err := pds.processSvc.GetByID(*processDependencyFormDTO.DependentProcessID); err != nil {
			return utility.NewValidationError(fmt.Errorf("process does not exist"), "not_exist", "process")
		}
	}

	if processDependencyFormDTO.OrganizationUnitID != nil {
		if _, err := pds.organizationUnitSvc.GetByID(*processDependencyFormDTO.OrganizationUnitID); err != nil {
			return utility.NewValidationError(fmt.Errorf("organization unit does not exist"), "not_exist", "organization unit")
		}
	}

	if processDependencyFormDTO.AssetID != nil {
		if _, err := pds.GetAssetByID(*processDependencyFormDTO.AssetID); err != nil {
			return utility.NewValidationError(fmt.Errorf("asset register does not exist"), "not_exist", "asset register")
		}
	}
	return nil
}

func (pds *ProcessDependencyService) GetAssetByID(id uint) (*models.AssetRegister, error) {
	asset, err := pds.repo.GetAssetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return asset, nil
}

func (pds *ProcessDependencyService) UpdateProcessDependencyRecoveryObjective(processDependency *models.ProcessDependency, tx utility.DatabaseTransaction) error {

	err := pds.repo.UpdateProcessDependencyRecoveryObjective(processDependency, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}
