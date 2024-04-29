package services

import (
	"errors"
	"fmt"
	logModels "jethur/grc/logs/models"
	activityLogServices "jethur/grc/logs/services"
	riskDTOs "jethur/grc/risk/dtos"
	"jethur/grc/risk/models"
	"jethur/grc/risk/repositories"
	"jethur/utility"

	"gorm.io/gorm"
)

type RiskMatrixLikelihoodServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*riskDTOs.RiskMatrixLikelihoodDTO, error)
	GetByID(id uint) (*models.RiskMatrixLikelihood, error)
	Create(currentUserID uint, riskMatrixLikelihoodFormDTO *riskDTOs.RiskMatrixLikelihoodFormDTO) (uint, error)
	Update(id uint, currentUserID uint, riskMatrixLikelihoodFormDTO *riskDTOs.RiskMatrixLikelihoodFormDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	Activate(ids []uint, currentUserID uint) error
	Deactivate(ids []uint, currentUserID uint) error
}

// Service and implementations
type RiskMatrixLikelihoodService struct {
	repo           repositories.RiskMatrixLikelihoodRepositoryInterface
	activityLogSvc activityLogServices.ActivityLogServiceInterface
	riskMatrixSvc  RiskMatrixServiceInterface
}

func NewRiskMatrixLikelihoodService(repo repositories.RiskMatrixLikelihoodRepositoryInterface, activityLogSvc activityLogServices.ActivityLogServiceInterface, riskMatrixService RiskMatrixServiceInterface) RiskMatrixLikelihoodServiceInterface {
	return &RiskMatrixLikelihoodService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
		riskMatrixSvc:  riskMatrixService,
	}
}

func (rmls *RiskMatrixLikelihoodService) GetAll(sort string, sortBy string, keyword string, status string) ([]*riskDTOs.RiskMatrixLikelihoodDTO, error) {
	return rmls.repo.GetAll(sort, sortBy, keyword, status)
}

func (rmls RiskMatrixLikelihoodService) GetByID(id uint) (*models.RiskMatrixLikelihood, error) {
	riskMatrixLikelihood, err := rmls.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return riskMatrixLikelihood, nil
}

func (rmls *RiskMatrixLikelihoodService) Create(currentUserID uint, riskMatrixLikelihoodFormDTO *riskDTOs.RiskMatrixLikelihoodFormDTO) (uint, error) {
	riskMatrix, err := rmls.riskMatrixSvc.GetByType("erm")
	if err != nil {
		return 0, err
	}

	if err := rmls.validateRiskMatrixLikelihoodFormDTO(riskMatrixLikelihoodFormDTO, riskMatrix.ID, nil); err != nil {
		return 0, err
	}

	likelihood := &models.RiskMatrixLikelihood{
		RiskMatrixID:  riskMatrix.ID,
		TitleEn:       riskMatrixLikelihoodFormDTO.TitleEn,
		TitleAr:       riskMatrixLikelihoodFormDTO.TitleAr,
		DescriptionEn: riskMatrixLikelihoodFormDTO.DescriptionEn,
		DescriptionAr: riskMatrixLikelihoodFormDTO.DescriptionAr,
		Status:        riskMatrixLikelihoodFormDTO.Status,
	}

	id, err := rmls.repo.Create(likelihood)
	if err != nil {
		return 0, err
	}

	go func() {
		err := rmls.activityLogSvc.LogActivity(
			logModels.CreateAction,
			models.RiskMatrixLikelihood{}.TableName(),
			id,
			riskMatrixLikelihoodFormDTO,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    logModels.CreateAction,
				"tableName": models.RiskMatrixLikelihood{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return id, nil

}

func (rmls *RiskMatrixLikelihoodService) Update(id uint, currentUserID uint, riskMatrixLikelihoodFormDTO *riskDTOs.RiskMatrixLikelihoodFormDTO) error {
	existingLikelihood, err := rmls.GetByID(id)

	if err != nil {
		return err
	}

	if err := rmls.validateRiskMatrixLikelihoodFormDTO(riskMatrixLikelihoodFormDTO, existingLikelihood.RiskMatrixID, &id); err != nil {
		return err
	}

	existingLikelihood.TitleEn = riskMatrixLikelihoodFormDTO.TitleEn
	existingLikelihood.TitleAr = riskMatrixLikelihoodFormDTO.TitleAr
	existingLikelihood.DescriptionEn = riskMatrixLikelihoodFormDTO.DescriptionEn
	existingLikelihood.DescriptionAr = riskMatrixLikelihoodFormDTO.DescriptionAr
	existingLikelihood.Status = riskMatrixLikelihoodFormDTO.Status

	err = rmls.repo.Update(existingLikelihood)

	if err != nil {
		return err
	}

	return nil

}

func (rmls *RiskMatrixLikelihoodService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIds []uint
	var successfulIds []uint
	for _, id := range ids {

		_, err := rmls.repo.GetByID(id)

		if err != nil {
			fmt.Print(err)
			continue
		}

		err = rmls.repo.Delete(id)
		if err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}
	}

	if len(successfulIds) > 0 {

		go func() {
			err := rmls.activityLogSvc.LogBulkActivities(
				logModels.DeleteAction, // Make sure this action constant is defined appropriately
				models.RiskMatrixLikelihood{}.TableName(),
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

func (rmls *RiskMatrixLikelihoodService) Activate(ids []uint, currentUserID uint) error {
	successfulIds, err := rmls.repo.Activate(ids)

	if len(successfulIds) > 0 {

		go func() {
			err := rmls.activityLogSvc.LogBulkActivities(
				logModels.ActivateAction,
				models.RiskMatrixLikelihood{}.TableName(),
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

func (rmls *RiskMatrixLikelihoodService) Deactivate(ids []uint, currentUserID uint) error {
	successfulIds, err := rmls.repo.Deactivate(ids)

	if len(successfulIds) > 0 {

		go func() {
			err := rmls.activityLogSvc.LogBulkActivities(
				logModels.DeactivateAction,
				models.RiskMatrixLikelihood{}.TableName(),
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

func (rmls *RiskMatrixLikelihoodService) validateRiskMatrixLikelihoodFormDTO(riskMatrixLikelihoodFormDTO *riskDTOs.RiskMatrixLikelihoodFormDTO, riskMatrixID uint, id *uint) error {
	if id == nil {
		if exists, err := rmls.repo.ExistsByTitle(riskMatrixLikelihoodFormDTO.TitleEn, riskMatrixLikelihoodFormDTO.TitleAr, riskMatrixID); err != nil {
			return err
		} else if exists {
			return utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
		}

	} else {
		if exists, err := rmls.repo.ExistsByTitleExcludingID(*id, riskMatrixLikelihoodFormDTO.TitleEn, riskMatrixLikelihoodFormDTO.TitleAr, riskMatrixID); err != nil {
			return err
		} else if exists {
			return utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
		}
	}

	return nil
}
