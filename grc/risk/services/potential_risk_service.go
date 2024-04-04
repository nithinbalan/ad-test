package services

import (
	"errors"
	"fmt"
	logModels "jethur/grc/logs/models"
	logServices "jethur/grc/logs/services"
	masterServices "jethur/grc/masters/services"
	riskDTOs "jethur/grc/risk/dtos"
	"jethur/grc/risk/models"
	riskRepositories "jethur/grc/risk/repositories"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type PotentialRiskServiceInterface interface {
	GetAll(params *riskDTOs.PotentialRiskListQueryParamsDTO) (*riskDTOs.PotentialRiskListResponseDTO, error)
	GetFullDetailsByID(id uint) (*riskDTOs.PotentialRiskDetailsDTO, error)
	GetByID(id uint) (*models.PotentialRisk, error)
	Create(currentUserID uint, potentialRiskFormDTO *riskDTOs.PotentialRiskFormDTO) (uint, error)
	Update(id uint, currentUserID uint, potentialRiskFormDTO *riskDTOs.PotentialRiskFormDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)
}

// Service and implementations
type PotentialRiskService struct {
	repo                   riskRepositories.PotentialRiskRepositoryInterface
	activityLogSvc         logServices.ActivityLogServiceInterface
	severityLevelSvc       masterServices.SeverityLevelServiceInterface
	riskCategorySvc        RiskCategoryServiceInterface
	potentialRiskStatusSvc PotentialRiskStatusServiceInterface
}

func NewPotentialRiskService(
	repo riskRepositories.PotentialRiskRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	severityLevelSvc masterServices.SeverityLevelServiceInterface,
	riskCategorySvc RiskCategoryServiceInterface,
	potentialRiskStatusSvc PotentialRiskStatusServiceInterface) PotentialRiskServiceInterface {
	return &PotentialRiskService{
		repo:                   repo,
		activityLogSvc:         activityLogSvc,
		severityLevelSvc:       severityLevelSvc,
		riskCategorySvc:        riskCategorySvc,
		potentialRiskStatusSvc: potentialRiskStatusSvc,
	}
}

func (prs *PotentialRiskService) GetAll(params *riskDTOs.PotentialRiskListQueryParamsDTO) (*riskDTOs.PotentialRiskListResponseDTO, error) {
	return prs.repo.GetAll(params)
}

func (prs PotentialRiskService) GetByID(id uint) (*models.PotentialRisk, error) {
	potentialRisk, err := prs.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return potentialRisk, nil
}

func (prs PotentialRiskService) GetFullDetailsByID(id uint) (*riskDTOs.PotentialRiskDetailsDTO, error) {

	var potentialRiskDTO riskDTOs.PotentialRiskDetailsDTO

	potentialRisk, err := prs.repo.GetFullDetailsByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	err = mapstructure.Decode(potentialRisk, &potentialRiskDTO)
	if err != nil {
		return nil, err
	}

	//todo :- also add extra fields which required in details screen

	return &potentialRiskDTO, nil
}

func (prs *PotentialRiskService) Create(currentUserID uint, potentialRiskFormDTO *riskDTOs.PotentialRiskFormDTO) (uint, error) {

	if err := prs.validatePotentialRiskFormDTO(potentialRiskFormDTO, nil); err != nil {
		return 0, err
	}

	potentialRiskStatus, err := prs.potentialRiskStatusSvc.GetByStatus("new")
	if err != nil {
		return 0, err
	}

	potentialRisk := &models.PotentialRisk{
		RiskCategoryID:        potentialRiskFormDTO.RiskCategoryID,
		SeverityLevelID:       potentialRiskFormDTO.SeverityLevelID,
		PotentialRiskStatusID: potentialRiskStatus.ID,
		Title:                 potentialRiskFormDTO.Title,
		Description:           potentialRiskFormDTO.Description,
		RiskSource:            "manual",
		CreatedBy:             currentUserID,
		UpdatedAt:             nil,
	}

	id, err := prs.repo.Create(potentialRisk)
	if err != nil {
		return 0, err
	}

	go func() {
		err := prs.activityLogSvc.LogActivity(
			logModels.CreateAction,
			models.PotentialRisk{}.TableName(),
			id,
			potentialRiskFormDTO,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    logModels.CreateAction,
				"tableName": models.PotentialRisk{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return id, nil
}

func (prs *PotentialRiskService) Update(id uint, currentUserID uint, potentialRiskFormDTO *riskDTOs.PotentialRiskFormDTO) error {

	existingPotentialRisk, err := prs.GetByID(id)

	if err != nil {
		return err
	}

	if err := prs.validatePotentialRiskFormDTO(potentialRiskFormDTO, &id); err != nil {
		return err
	}

	existingPotentialRisk.RiskCategoryID = potentialRiskFormDTO.RiskCategoryID
	existingPotentialRisk.SeverityLevelID = potentialRiskFormDTO.SeverityLevelID
	existingPotentialRisk.Title = potentialRiskFormDTO.Title
	existingPotentialRisk.Description = potentialRiskFormDTO.Description
	existingPotentialRisk.UpdatedBy = &currentUserID

	err = prs.repo.Update(existingPotentialRisk)

	if err != nil {
		return err
	}

	return nil
}

func (prs *PotentialRiskService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIds []uint
	var successfulIds []uint
	for _, id := range ids {

		_, err := prs.repo.GetByID(id)

		if err != nil {
			continue
		}

		err = prs.repo.Delete(id)
		if err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}
	}

	if len(successfulIds) > 0 {

		go func() {
			err := prs.activityLogSvc.LogBulkActivities(
				logModels.DeleteAction, // Make sure this action constant is defined appropriately
				models.PotentialRisk{}.TableName(),
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

func (prs *PotentialRiskService) validatePotentialRiskFormDTO(potentialRiskFormDTO *riskDTOs.PotentialRiskFormDTO, id *uint) error {
	if _, err := prs.riskCategorySvc.GetByID(potentialRiskFormDTO.RiskCategoryID); err != nil {
		return utility.NewValidationError(fmt.Errorf("risk category does not exist"), "not_exist", "risk category")
	}

	if _, err := prs.severityLevelSvc.GetByID(potentialRiskFormDTO.SeverityLevelID); err != nil {
		return utility.NewValidationError(fmt.Errorf("severity level does not exist"), "not_exist", "severity level")
	}

	if id == nil {
		if exists, err := prs.repo.ExistsByTitle(potentialRiskFormDTO.Title, potentialRiskFormDTO.RiskCategoryID); err != nil {
			return err
		} else if exists {
			return utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
		}

	} else {
		if exists, err := prs.repo.ExistsByTitleExcludingID(*id, potentialRiskFormDTO.Title, potentialRiskFormDTO.RiskCategoryID); err != nil {
			return err
		} else if exists {
			return utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
		}
	}

	return nil
}
