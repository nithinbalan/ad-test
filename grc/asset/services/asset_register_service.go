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

type AssetRegisterServiceInterface interface {
	GetAll(params *dtos.AssetRegisterListQueryParamsDTO) (*dtos.AssetRegisterListResponseDTO, error)
	GetByID(id uint) (*models.AssetRegister, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*dtos.AssetRegisterDetailDTO, error)
	Create(assetForm *dtos.AssetRegisterFormDTO, currentUserID uint) (uint, error)
	Update(id uint, assetForm *dtos.AssetRegisterFormDTO, currentUserID uint) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)
}

type AssetRegisterService struct {
	repo             repositories.AssetRegisterRepositoryInterface
	activityLogSvc   logServices.ActivityLogServiceInterface
	assetTypeSvc     AssetTypeServiceInterface
	assetCategorySvc AssetCategoryServiceInterface
	assetStatusSvc   AssetStatusServiceInterface
	userSvc          organizationServices.UserServiceInterface
}

func NewAssetRegisterService(
	repo repositories.AssetRegisterRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	assetTypeSvc AssetTypeServiceInterface,
	assetCategorySvc AssetCategoryServiceInterface,
	assetStatusSvc AssetStatusServiceInterface,
	userSvc organizationServices.UserServiceInterface,
) AssetRegisterServiceInterface {
	return &AssetRegisterService{
		repo:             repo,
		activityLogSvc:   activityLogSvc,
		assetTypeSvc:     assetTypeSvc,
		assetCategorySvc: assetCategorySvc,
		assetStatusSvc:   assetStatusSvc,
		userSvc:          userSvc,
	}
}

func (ars *AssetRegisterService) GetAll(params *dtos.AssetRegisterListQueryParamsDTO) (*dtos.AssetRegisterListResponseDTO, error) {
	return ars.repo.GetAll(params)
}

func (ars *AssetRegisterService) GetByID(id uint) (*models.AssetRegister, error) {
	asset, err := ars.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return asset, nil
}

func (as *AssetRegisterService) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*dtos.AssetRegisterDetailDTO, error) {
	var assetDetailDTO dtos.AssetRegisterDetailDTO

	asset, err := as.repo.GetFullDetailsByID(id, tx...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	err = mapstructure.Decode(asset, &assetDetailDTO)
	if err != nil {
		return nil, err
	}

	if asset.AcquisitionDate != nil {
		acquisitionDateString := asset.AcquisitionDate.Format("2006-01-02")
		assetDetailDTO.AcquisitionDateString = &acquisitionDateString
	}

	// Add additional fields as necessary for the details screen
	// For example, setting relationship data or additional fields as required by the DTO

	return &assetDetailDTO, nil
}

func (ars *AssetRegisterService) Create(assetRegisterFormDTO *dtos.AssetRegisterFormDTO, currentUserID uint) (uint, error) {
	if err := ars.validateAssetRegisterFormDTO(assetRegisterFormDTO); err != nil {
		return 0, err
	}
	dates, err := ars.parseAssetRegisterDates(assetRegisterFormDTO)
	if err != nil {
		return 0, err
	}

	referenceNumber, err := ars.repo.GenerateReferenceNumber()
	if err != nil {
		return 0, err
	}

	organizationUnitID, err := ars.userSvc.GetUserPrimaryUnitID(assetRegisterFormDTO.ResponsibleUserID)
	if err != nil {
		return 0, utility.NewValidationError(fmt.Errorf("organization unit not set for responsible user"), "not_exist", "responsible user unit")
	}

	asset := &models.AssetRegister{
		ReferenceNumber:   referenceNumber,
		TitleEn:           assetRegisterFormDTO.TitleEn,
		TitleAr:           assetRegisterFormDTO.TitleAr,
		AssetTypeID:       assetRegisterFormDTO.AssetTypeID,
		AssetCategoryID:   assetRegisterFormDTO.AssetCategoryID,
		AssetStatusID:     assetRegisterFormDTO.AssetStatusID,
		ResponsibleUserID: assetRegisterFormDTO.ResponsibleUserID,
		ResponsibleUnitID: organizationUnitID,
		Location:          assetRegisterFormDTO.Location,
		Description:       assetRegisterFormDTO.Description,
		AcquisitionDate:   dates.AcquisitionDate,
		CreatedBYID:       currentUserID,
	}

	assetID, err := ars.repo.Create(asset)
	if err != nil {
		return 0, err
	}

	go func() {
		err := ars.activityLogSvc.LogActivity(
			models.CreateAction,
			models.AssetRegister{}.TableName(),
			assetID,
			assetRegisterFormDTO,
			currentUserID,
			nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.AssetRegister{}.TableName(),
				"recordID":  assetID,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log control activity", fields)
		}
	}()

	return assetID, nil
}

func (ars *AssetRegisterService) Update(id uint, assetRegisterFormDTO *dtos.AssetRegisterFormDTO, currentUserID uint) error {
	existingAsset, err := ars.repo.GetByID(id)
	if err != nil {
		return err
	}

	//oldResponsibleUserID := existingAsset.ResponsibleUserID

	if err := ars.validateAssetRegisterFormDTO(assetRegisterFormDTO); err != nil {
		return err
	}
	dates, err := ars.parseAssetRegisterDates(assetRegisterFormDTO)
	if err != nil {
		return err
	}

	if existingAsset.ResponsibleUserID != assetRegisterFormDTO.ResponsibleUserID {
		organizationUnitID, err := ars.userSvc.GetUserPrimaryUnitID(assetRegisterFormDTO.ResponsibleUserID)
		if err != nil {
			return utility.NewValidationError(fmt.Errorf("organization unit not set for owner"), "not_exist", "control owner unit")
		}
		existingAsset.ResponsibleUnitID = organizationUnitID
	}

	// Update fields
	existingAsset.TitleEn = assetRegisterFormDTO.TitleEn
	existingAsset.TitleAr = assetRegisterFormDTO.TitleAr
	existingAsset.AssetTypeID = assetRegisterFormDTO.AssetTypeID
	existingAsset.AssetCategoryID = assetRegisterFormDTO.AssetCategoryID
	existingAsset.AssetStatusID = assetRegisterFormDTO.AssetStatusID
	existingAsset.ResponsibleUserID = assetRegisterFormDTO.ResponsibleUserID
	existingAsset.ResponsibleUnitID = assetRegisterFormDTO.ResponsibleUserID
	existingAsset.Location = assetRegisterFormDTO.Location
	existingAsset.Description = assetRegisterFormDTO.Description
	existingAsset.AcquisitionDate = dates.AcquisitionDate
	existingAsset.UpdatedBYID = &currentUserID

	return ars.repo.Update(existingAsset)
}

func (as *AssetRegisterService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIds []uint
	var successfulIds []uint

	for _, id := range ids {
		err := as.repo.Delete(id)
		if err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}
	}

	if len(successfulIds) > 0 {
		go func() {
			err := as.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.AssetRegister{}.TableName(),
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

func (ars *AssetRegisterService) validateAssetRegisterFormDTO(assetRegisterFormDTO *dtos.AssetRegisterFormDTO) error {
	if _, err := ars.assetTypeSvc.GetByID(uint(assetRegisterFormDTO.AssetTypeID)); err != nil {
		return utility.NewValidationError(fmt.Errorf("asset type does not exist"), "not_exist", "asset type")
	}

	if _, err := ars.assetCategorySvc.GetByID(uint(assetRegisterFormDTO.AssetCategoryID)); err != nil {
		return utility.NewValidationError(fmt.Errorf("asset category does not exist"), "not_exist", "asset category")
	}

	if _, err := ars.userSvc.GetByID(uint(assetRegisterFormDTO.ResponsibleUserID)); err != nil {
		return utility.NewValidationError(fmt.Errorf("responsible user does not exist"), "not_exist", "responsible user")
	}

	if _, err := ars.assetStatusSvc.GetByID(uint(assetRegisterFormDTO.AssetStatusID)); err != nil {
		return utility.NewValidationError(fmt.Errorf("asset status does not exist"), "not_exist", "asset status")
	}

	return nil
}

func (ars *AssetRegisterService) parseAssetRegisterDates(assetRegisterFormDTO *dtos.AssetRegisterFormDTO) (*dtos.AssetRegisterParsedDates, error) {
	parsed := &dtos.AssetRegisterParsedDates{}

	var err error
	if assetRegisterFormDTO.AcquisitionDate != nil {
		parsed.AcquisitionDate, err = utility.ParseDate(*assetRegisterFormDTO.AcquisitionDate)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid acquisition date"), "invalid_date", "acquisition date")
		}
	}

	return parsed, nil
}
