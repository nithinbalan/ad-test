package services

import (
	"fmt"
	logModels "jethur/grc/logs/models"
	activityLogServices "jethur/grc/logs/services"
	riskDTOs "jethur/grc/risk/dtos"
	"jethur/grc/risk/models"
	"jethur/grc/risk/repositories"
	"jethur/utility"
)

type RiskMatrixImpactServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*riskDTOs.RiskMatrixImpactDTO, error)
	GetByID(id uint) (*models.RiskMatrixImpact, error)
	Create(currentUserID uint, riskMatrixImpactFormDTO *riskDTOs.RiskMatrixImpactFormDTO) (uint, error)
	Update(id uint, currentUserID uint, riskMatrixImpactFormDTO *riskDTOs.RiskMatrixImpactFormDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	Activate(ids []uint, currentUserID uint) error
	Deactivate(ids []uint, currentUserID uint) error
}

// Service and implementations
type RiskMatrixImpactService struct {
	repo           repositories.RiskMatrixImpactRepositoryInterface
	activityLogSvc activityLogServices.ActivityLogServiceInterface
	riskMatrixSvc  RiskMatrixServiceInterface
}

func NewRiskMatrixImpactService(repo repositories.RiskMatrixImpactRepositoryInterface, activityLogSvc activityLogServices.ActivityLogServiceInterface, riskMatrixService RiskMatrixServiceInterface) RiskMatrixImpactServiceInterface {
	return &RiskMatrixImpactService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
		riskMatrixSvc:  riskMatrixService,
	}
}

func (rmis *RiskMatrixImpactService) GetAll(sort string, sortBy string, keyword string, status string) ([]*riskDTOs.RiskMatrixImpactDTO, error) {
	return rmis.repo.GetAll(sort, sortBy, keyword, status)
}

func (rmis RiskMatrixImpactService) GetByID(id uint) (*models.RiskMatrixImpact, error) {
	return rmis.repo.GetByID(id)
}

func (rmis *RiskMatrixImpactService) Create(currentUserID uint, riskMatrixImpactFormDTO *riskDTOs.RiskMatrixImpactFormDTO) (uint, error) {
	riskMatrix, err := rmis.riskMatrixSvc.GetByType("erm")
	if err != nil {
		return 0, err
	}

	if err := rmis.validateRiskMatrixImpactFormDTO(riskMatrixImpactFormDTO, riskMatrix.ID, nil); err != nil {
		return 0, err
	}

	impact := &models.RiskMatrixImpact{
		RiskMatrixID:  riskMatrix.ID,
		TitleEn:       riskMatrixImpactFormDTO.TitleEn,
		TitleAr:       riskMatrixImpactFormDTO.TitleAr,
		DescriptionEn: riskMatrixImpactFormDTO.DescriptionEn,
		DescriptionAr: riskMatrixImpactFormDTO.DescriptionAr,
		Status:        riskMatrixImpactFormDTO.Status,
	}

	id, err := rmis.repo.Create(impact)
	if err != nil {
		return 0, err
	}

	go func() {
		err := rmis.activityLogSvc.LogActivity(
			logModels.CreateAction,
			models.RiskMatrixImpact{}.TableName(),
			id,
			riskMatrixImpactFormDTO,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    logModels.CreateAction,
				"tableName": models.RiskMatrixImpact{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return id, nil

}

func (rmis *RiskMatrixImpactService) Update(id uint, currentUserID uint, riskMatrixImpactFormDTO *riskDTOs.RiskMatrixImpactFormDTO) error {
	existingImpact, err := rmis.GetByID(id)

	if err != nil {
		return err
	}

	if err := rmis.validateRiskMatrixImpactFormDTO(riskMatrixImpactFormDTO, existingImpact.RiskMatrixID, &id); err != nil {
		return err
	}

	existingImpact.TitleEn = riskMatrixImpactFormDTO.TitleEn
	existingImpact.TitleAr = riskMatrixImpactFormDTO.TitleAr
	existingImpact.DescriptionEn = riskMatrixImpactFormDTO.DescriptionEn
	existingImpact.DescriptionAr = riskMatrixImpactFormDTO.DescriptionAr
	existingImpact.Status = riskMatrixImpactFormDTO.Status

	err = rmis.repo.Update(existingImpact)

	if err != nil {
		return err
	}

	return nil

}

func (rmis *RiskMatrixImpactService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIds []uint
	var successfulIds []uint
	for _, id := range ids {

		_, err := rmis.repo.GetByID(id)

		if err != nil {
			fmt.Print(err)
			continue
		}

		err = rmis.repo.Delete(id)
		if err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}
	}

	if len(successfulIds) > 0 {

		go func() {
			err := rmis.activityLogSvc.LogBulkActivities(
				logModels.DeleteAction, // Make sure this action constant is defined appropriately
				models.RiskMatrixImpact{}.TableName(),
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

func (rmis *RiskMatrixImpactService) Activate(ids []uint, currentUserID uint) error {
	successfulIds, err := rmis.repo.Activate(ids)

	if len(successfulIds) > 0 {

		go func() {
			err := rmis.activityLogSvc.LogBulkActivities(
				logModels.ActivateAction,
				models.RiskMatrixImpact{}.TableName(),
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

	return err
}

func (rmis *RiskMatrixImpactService) Deactivate(ids []uint, currentUserID uint) error {
	successfulIds, err := rmis.repo.Deactivate(ids)

	if len(successfulIds) > 0 {

		go func() {
			err := rmis.activityLogSvc.LogBulkActivities(
				logModels.DeactivateAction,
				models.RiskMatrixImpact{}.TableName(),
				successfulIds,
				nil,
				currentUserID,
				nil,
			)
			if err != nil {
				// @todo - save to log file
				// You might want to log this error or handle it appropriately
				fmt.Printf("Failed to log bulk delete activity: %v\n", err)
			}
		}()
	}

	return err
}

func (rmis *RiskMatrixImpactService) validateRiskMatrixImpactFormDTO(riskMatrixImpactFormDTO *riskDTOs.RiskMatrixImpactFormDTO, riskMatrixID uint, id *uint) error {
	if id == nil {
		if exists, err := rmis.repo.ExistsByTitle(riskMatrixImpactFormDTO.TitleEn, riskMatrixImpactFormDTO.TitleAr, riskMatrixID); err != nil {
			return err
		} else if exists {
			return utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
		}

	} else {
		if exists, err := rmis.repo.ExistsByTitleExcludingID(*id, riskMatrixImpactFormDTO.TitleEn, riskMatrixImpactFormDTO.TitleAr, riskMatrixID); err != nil {
			return err
		} else if exists {
			return utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
		}
	}

	return nil
}
