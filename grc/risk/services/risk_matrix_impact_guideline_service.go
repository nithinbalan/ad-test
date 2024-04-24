package services

import (
	"errors"
	"fmt"
	logModels "jethur/grc/logs/models"
	logServices "jethur/grc/logs/services"
	riskDTOs "jethur/grc/risk/dtos"
	"jethur/grc/risk/models"
	riskRepositories "jethur/grc/risk/repositories"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type RiskMatrixImpactGuidelineServiceInterface interface {
	GetAll(params *riskDTOs.RiskMatrixImpactGuidelineListQueryParamsDTO) (*riskDTOs.RiskMatrixImpactGuidelineListResponseDTO, error)
	GetFullDetailsByID(id uint) (*riskDTOs.RiskMatrixImpactGuidelineDetailsDTO, error)
	GetByID(id uint) (*models.RiskMatrixImpactGuideline, error)
	Create(currentUserID uint, riskMatrixImpactGuidelineFormDTO *riskDTOs.RiskMatrixImpactGuidelineFormDTO) ([]uint, error)
	Update(id uint, currentUserID uint, riskMatrixImpactGuidelineUpdateFormDTO *riskDTOs.RiskMatrixImpactGuidelineUpdateFormDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)
}

// Service and implementations
type RiskMatrixImpactGuidelineService struct {
	repo                  riskRepositories.RiskMatrixImpactGuidelineRepositoryInterface
	activityLogSvc        logServices.ActivityLogServiceInterface
	riskImpactCategorySvc RiskImpactCategoryServiceInterface
	// riskMatrixImpactSvc     RiskMatrixImpactServiceInterface
}

func NewRiskMatrixImpactGuidelineService(
	repo riskRepositories.RiskMatrixImpactGuidelineRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	riskImpactCategorySvc RiskImpactCategoryServiceInterface,
	// riskMatrixImpactSvc     RiskMatrixImpactServiceInterface,
) RiskMatrixImpactGuidelineServiceInterface {
	return &RiskMatrixImpactGuidelineService{
		repo:                  repo,
		activityLogSvc:        activityLogSvc,
		riskImpactCategorySvc: riskImpactCategorySvc,
		// riskMatrixImpactSvc:     riskMatrixImpactSvc,
	}
}

func (rmigs *RiskMatrixImpactGuidelineService) GetAll(params *riskDTOs.RiskMatrixImpactGuidelineListQueryParamsDTO) (*riskDTOs.RiskMatrixImpactGuidelineListResponseDTO, error) {
	return rmigs.repo.GetAll(params)
}

func (rmigs RiskMatrixImpactGuidelineService) GetByID(id uint) (*models.RiskMatrixImpactGuideline, error) {
	riskMatrixImpactGuideline, err := rmigs.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return riskMatrixImpactGuideline, nil
}

func (rmigs RiskMatrixImpactGuidelineService) GetFullDetailsByID(id uint) (*riskDTOs.RiskMatrixImpactGuidelineDetailsDTO, error) {

	var riskMatrixImpactGuidelineDTO riskDTOs.RiskMatrixImpactGuidelineDetailsDTO

	riskMatrixImpactGuideline, err := rmigs.repo.GetFullDetailsByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	err = mapstructure.Decode(riskMatrixImpactGuideline, &riskMatrixImpactGuidelineDTO)
	if err != nil {
		return nil, err
	}

	//todo :- also add extra fields which required in details screen

	return &riskMatrixImpactGuidelineDTO, nil
}

func (rmigs *RiskMatrixImpactGuidelineService) Create(currentUserID uint, riskMatrixImpactGuidelineFormDTO *riskDTOs.RiskMatrixImpactGuidelineFormDTO) ([]uint, error) {
	tx := rmigs.repo.BeginTransaction()

	if err := rmigs.validateRiskMatrixImpactGuidelineFormDTO(riskMatrixImpactGuidelineFormDTO, nil); err != nil {
		return nil, err
	}

	// Create a slice to store the data to be created
	var guidelinesToCreate []*models.RiskMatrixImpactGuideline
	// Iterate over each guideline in the form DTO and prepare the data
	for _, guideline := range riskMatrixImpactGuidelineFormDTO.Guidelines {
		guideline := &models.RiskMatrixImpactGuideline{
			RiskImpactCategoryID: riskMatrixImpactGuidelineFormDTO.RiskImpactCategoryID,
			RiskMatrixImpactID:   riskMatrixImpactGuidelineFormDTO.RiskMatrixImpactID,
			Guideline:            guideline,
			Status:               riskMatrixImpactGuidelineFormDTO.Status,
			CreatedBy:            currentUserID,
		}
		guidelinesToCreate = append(guidelinesToCreate, guideline)
	}

	// Call the repository method to create the records in batch
	ids, err := rmigs.repo.CreateBatch(guidelinesToCreate, tx)

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	go func() {
		err := rmigs.activityLogSvc.LogBulkActivities(
			logModels.CreateAction,
			models.RiskMatrixImpactGuideline{}.TableName(),
			ids,
			riskMatrixImpactGuidelineFormDTO,
			currentUserID,
			nil,
		)
		if err != nil {
			// @todo - save to log file
			fmt.Printf("Failed to log bulk activity: %v\n", err)
		}
	}()
	tx.Commit()

	return ids, nil
}

func (rmigs *RiskMatrixImpactGuidelineService) Update(id uint, currentUserID uint, riskMatrixImpactGuidelineUpdateFormDTO *riskDTOs.RiskMatrixImpactGuidelineUpdateFormDTO) error {

	// Retrieve existing risk matrix impact guideline by ID
	existingRiskMatrixImpactGuideline, err := rmigs.GetByID(id)
	if err != nil {
		return err
	}

	// Validate the form DTO
	if err := rmigs.validateRiskMatrixImpactGuidelineUpdateFormDTO(riskMatrixImpactGuidelineUpdateFormDTO); err != nil {
		return err
	}

	// Update the specific guideline
	existingRiskMatrixImpactGuideline.Guideline = riskMatrixImpactGuidelineUpdateFormDTO.Guideline

	// Update other fields as needed
	existingRiskMatrixImpactGuideline.RiskImpactCategoryID = riskMatrixImpactGuidelineUpdateFormDTO.RiskImpactCategoryID
	existingRiskMatrixImpactGuideline.RiskMatrixImpactID = riskMatrixImpactGuidelineUpdateFormDTO.RiskMatrixImpactID
	// Update other fields similarly...

	// Update the existing risk matrix impact guideline in the repository
	err = rmigs.repo.Update(existingRiskMatrixImpactGuideline)
	if err != nil {
		return err
	}

	return nil
}

func (rmigs *RiskMatrixImpactGuidelineService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIds []uint
	var successfulIds []uint
	for _, id := range ids {

		_, err := rmigs.repo.GetByID(id)

		if err != nil {
			continue
		}

		err = rmigs.repo.Delete(id)
		if err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}
	}

	if len(successfulIds) > 0 {

		go func() {
			err := rmigs.activityLogSvc.LogBulkActivities(
				logModels.DeleteAction, // Make sure this action constant is defined appropriately
				models.RiskMatrixImpactGuideline{}.TableName(),
				successfulIds,
				nil,
				currentUserID,
				nil,
			)
			if err != nil {
				// @todo - save to log file
				fmt.Printf("Failed to log bulk delete activity: %v\n", err)
			}
		}()
	}

	return failedIds, nil
}

func (rmigs *RiskMatrixImpactGuidelineService) validateRiskMatrixImpactGuidelineFormDTO(riskMatrixImpactGuidelineFormDTO *riskDTOs.RiskMatrixImpactGuidelineFormDTO, id *uint) error {
	// Check if Guidelines array exists and is not null
	if riskMatrixImpactGuidelineFormDTO.Guidelines == nil || len(riskMatrixImpactGuidelineFormDTO.Guidelines) == 0 {
		return utility.NewValidationError(fmt.Errorf("guidelines array is required and cannot be empty"), "required", "guidelines")
	}

	// Validate Guidelines array
	for _, guideline := range riskMatrixImpactGuidelineFormDTO.Guidelines {
		if len(guideline) > 500 {
			return utility.NewValidationError(fmt.Errorf("guideline exceeds maximum length of 500 characters"), "max_length_exceeded", "guidelines")
		}
		// You can add more validation rules for each guideline here if needed
	}

	// Validate RiskImpactCategoryID
	if _, err := rmigs.riskImpactCategorySvc.GetByID(riskMatrixImpactGuidelineFormDTO.RiskImpactCategoryID); err != nil {
		return utility.NewValidationError(fmt.Errorf("risk impact category does not exist"), "not_exist", "risk impact category")
	}

	// Validate RiskMatrixImpactID
	// if _, err := rmigs.riskMatrixImpactSvc.GetByID(riskMatrixImpactGuidelineFormDTO.RiskMatrixImpactID); err != nil {
	//     return utility.NewValidationError(fmt.Errorf("risk matrix impact does not exist"), "not_exist", "risk matrix impact")
	// }

	return nil
}

func (rmigs *RiskMatrixImpactGuidelineService) validateRiskMatrixImpactGuidelineUpdateFormDTO(riskMatrixImpactGuidelineUpdateFormDTO *riskDTOs.RiskMatrixImpactGuidelineUpdateFormDTO) error {
	// Validate RiskImpactCategoryID
	if _, err := rmigs.riskImpactCategorySvc.GetByID(riskMatrixImpactGuidelineUpdateFormDTO.RiskImpactCategoryID); err != nil {
		return utility.NewValidationError(fmt.Errorf("risk impact category does not exist"), "not_exist", "risk impact category")
	}

	// Validate RiskMatrixImpactID
	// if _, err := rmigs.riskMatrixImpactSvc.GetByID(riskMatrixImpactGuidelineUpdateFormDTO.RiskMatrixImpactID); err != nil {
	//     return utility.NewValidationError(fmt.Errorf("risk matrix impact does not exist"), "not_exist", "risk matrix impact")
	// }

	// Validate Guideline
	if len(riskMatrixImpactGuidelineUpdateFormDTO.Guideline) == 0 {
		return utility.NewValidationError(fmt.Errorf("guideline is required"), "required", "guideline")
	}
	if len(riskMatrixImpactGuidelineUpdateFormDTO.Guideline) > 500 {
		return utility.NewValidationError(fmt.Errorf("guideline exceeds maximum length of 500 characters"), "max_length_exceeded", "guideline")
	}

	return nil
}
