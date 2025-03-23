package services

import (
	"errors"
	"fmt"
	assetDTOs "jethur/grc/asset/dtos"
	assetRepositories "jethur/grc/asset/repositories"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type AssetKriServiceInterface interface {
	GetAll(assetID uint, params *assetDTOs.AssetKriListQueryParamsDTO) (*assetDTOs.AssetKriListResponseDTO, error)
	GetByID(assetID uint, id uint) (*models.AssetKri, error)
	GetFullDetailsByID(assetID uint, id uint) (*assetDTOs.AssetKriDetailDTO, error)
	Create(assetID uint, currentUserID uint, form *assetDTOs.AssetKriFormDTO) ([]uint, error)
	Update(assetID uint, currentUserID uint, form *assetDTOs.AssetKriFormDTO) error
	Delete(assetID uint, ids []uint, currentUserID uint) ([]uint, error)
}

type AssetKriService struct {
	repo             assetRepositories.AssetKriRepositoryInterface
	activityLogSvc   logServices.ActivityLogServiceInterface
	assetRegisterSvc AssetRegisterServiceInterface
}

func NewAssetKriService(
	repo assetRepositories.AssetKriRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	assetRegisterSvc AssetRegisterServiceInterface,
) AssetKriServiceInterface {
	return &AssetKriService{
		repo:             repo,
		activityLogSvc:   activityLogSvc,
		assetRegisterSvc: assetRegisterSvc,
	}
}

func (aks *AssetKriService) GetAll(assetID uint, params *assetDTOs.AssetKriListQueryParamsDTO) (*assetDTOs.AssetKriListResponseDTO, error) {
	if assetID != 0 {
		if _, err := aks.assetRegisterSvc.GetByID(assetID); err != nil {
			return nil, err
		}
	}
	return aks.repo.GetAll(assetID, params)
}

func (aks *AssetKriService) GetByID(assetID uint, id uint) (*models.AssetKri, error) {
	if _, err := aks.assetRegisterSvc.GetByID(assetID); err != nil {
		return nil, err
	}

	kri, err := aks.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return kri, nil
}

func (aks *AssetKriService) GetFullDetailsByID(assetID uint, id uint) (*assetDTOs.AssetKriDetailDTO, error) {
	var dto assetDTOs.AssetKriDetailDTO

	if _, err := aks.assetRegisterSvc.GetByID(assetID); err != nil {
		return nil, err
	}

	kri, err := aks.repo.GetFullDetailsByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	err = mapstructure.Decode(kri, &dto)
	if err != nil {
		return nil, err
	}

	return &dto, nil
}

func (aks *AssetKriService) Create(assetID uint, currentUserID uint, form *assetDTOs.AssetKriFormDTO) ([]uint, error) {
	tx := aks.repo.BeginTransaction()

	if _, err := aks.assetRegisterSvc.GetByID(assetID); err != nil {
		return nil, err
	}

	if err := aks.validateAssetKriFormDTO(assetID, form); err != nil {
		return nil, err
	}

	var kris []*models.AssetKri
	for _, kriID := range form.KriIDs {
		kris = append(kris, &models.AssetKri{
			AssetRegisterID: assetID,
			KriRegisterID:   *kriID,
			CreatedBYID:     currentUserID,
		})
	}

	ids, err := aks.repo.Create(kris, tx)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()

	go func() {
		_ = aks.activityLogSvc.LogBulkActivities(
			models.CreateAction,
			models.AssetKri{}.TableName(),
			ids,
			kris,
			currentUserID,
			nil,
		)
	}()

	return ids, nil
}

func (aks *AssetKriService) Update(assetID uint, currentUserID uint, form *assetDTOs.AssetKriFormDTO) error {
	tx := aks.repo.BeginTransaction()

	existingAsset, err := aks.assetRegisterSvc.GetFullDetailsByID(assetID)
	if err != nil {
		return err
	}

	if err := aks.validateAssetKriFormDTO(assetID, form); err != nil {
		return err
	}

	existingKRIIDs := make(map[uint]bool)
	for _, kri := range existingAsset.Kris {
		existingKRIIDs[kri.KriID] = true
	}

	newKRIIDs := make(map[uint]bool)
	for _, kri := range form.KriIDs {
		newKRIIDs[*kri] = true
	}

	var toAdd []*models.AssetKri
	for _, kriID := range form.KriIDs {
		if !existingKRIIDs[*kriID] {
			toAdd = append(toAdd, &models.AssetKri{
				AssetRegisterID: assetID,
				KriRegisterID:   *kriID,
				CreatedBYID:     currentUserID,
			})
		}
	}

	if len(toAdd) > 0 {
		_, err := aks.repo.Create(toAdd, tx)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	var toDelete []*uint
	for _, existing := range existingAsset.Kris {
		if !newKRIIDs[existing.KriID] {
			toDelete = append(toDelete, existing.ID)
		}
	}
	if len(toDelete) > 0 {
		if err := aks.repo.DeleteAssetKRIs(toDelete); err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()

	go func() {
		// Build log diff
		added := []uint{}
		removed := []uint{}
		for id := range newKRIIDs {
			if !existingKRIIDs[id] {
				added = append(added, id)
			}
		}
		for id := range existingKRIIDs {
			if !newKRIIDs[id] {
				removed = append(removed, id)
			}
		}

		if len(added) > 0 || len(removed) > 0 {
			changedData := map[string]interface{}{
				"associated_items": map[string]interface{}{
					"kri_ids": map[string]interface{}{
						"added":   added,
						"removed": removed,
					},
				},
			}

			_ = aks.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.AssetKri{}.TableName(),
				assetID,
				changedData,
				currentUserID,
				nil,
			)
		}
	}()

	return nil
}

func (aks *AssetKriService) Delete(assetID uint, ids []uint, currentUserID uint) ([]uint, error) {
	tx := aks.repo.BeginTransaction()

	if _, err := aks.assetRegisterSvc.GetByID(assetID); err != nil {
		return nil, err
	}

	var failed []uint
	var successful []uint
	var deleted []map[string]interface{}

	for _, id := range ids {
		_, err := aks.repo.GetByID(id)
		if err != nil {
			failed = append(failed, id)
			continue
		}
		if err := aks.repo.Delete(&id, tx); err != nil {
			tx.Rollback()
			failed = append(failed, id)
			continue
		}
		successful = append(successful, id)
		deleted = append(deleted, map[string]interface{}{
			"id":                id,
			"asset_register_id": assetID,
		})
	}

	tx.Commit()

	if len(successful) > 0 {
		go func() {
			_ = aks.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.AssetKri{}.TableName(),
				successful,
				deleted,
				currentUserID,
				nil,
			)
		}()
	}

	return failed, nil
}

func (aks *AssetKriService) validateAssetKriFormDTO(assetID uint, form *assetDTOs.AssetKriFormDTO) error {
	if len(form.KriIDs) == 0 {
		return utility.NewValidationError(fmt.Errorf("no KRI IDs provided"), "empty", "kri")
	}

	kriIDs := map[uint]bool{}
	for _, kriID := range form.KriIDs {
		if kriIDs[*kriID] {
			return utility.NewValidationError(fmt.Errorf("duplicate KRI ID: %d", *kriID), "duplicate", "kri")
		}
		if _, err := aks.repo.GetKRIByID(*kriID); err != nil {
			return utility.NewValidationError(fmt.Errorf("KRI does not exist for ID: %d", *kriID), "not_exist", "kri")
		}

		// Check if already associated
		existing, err := aks.repo.GetKrisByAssetID(assetID)
		if err != nil {
			return utility.NewValidationError(fmt.Errorf("failed to fetch existing KRIs"), "fetch_error", "kri")
		}
		for _, e := range existing {
			if e.KriRegisterID == *kriID {
				return utility.NewValidationError(fmt.Errorf("KRI already associated"), "already_associated", "kri")
			}
		}

		kriIDs[*kriID] = true
	}
	return nil
}
