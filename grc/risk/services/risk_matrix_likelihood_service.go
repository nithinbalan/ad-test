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

func (rmis *RiskMatrixLikelihoodService) GetAll(sort string, sortBy string, keyword string, status string) ([]*riskDTOs.RiskMatrixLikelihoodDTO, error) {
	return rmis.repo.GetAll(sort, sortBy, keyword, status)
}

func (rmis RiskMatrixLikelihoodService) GetByID(id uint) (*models.RiskMatrixLikelihood, error) {
	return rmis.repo.GetByID(id)
}

func (rmis *RiskMatrixLikelihoodService) Create(currentUserID uint, riskMatrixLikelihoodFormDTO *riskDTOs.RiskMatrixLikelihoodFormDTO) (uint, error) {

	riskMatrix, err := rmis.riskMatrixSvc.GetByType("erm")
	if err != nil {
		return 0, err
	}

	if err := rmis.validateRiskMatrixLikelihoodFormDTO(riskMatrixLikelihoodFormDTO, riskMatrix.ID, nil); err != nil {
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

	id, err := rmis.repo.Create(likelihood)
	if err != nil {
		return 0, err
	}

	go func() {
		err := rmis.activityLogSvc.LogActivity(
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

func (rmis *RiskMatrixLikelihoodService) Update(id uint, currentUserID uint, riskMatrixLikelihoodFormDTO *riskDTOs.RiskMatrixLikelihoodFormDTO) error {

	existingLikelihood, err := rmis.GetByID(id)

	if err != nil {
		return err
	}

	if err := rmis.validateRiskMatrixLikelihoodFormDTO(riskMatrixLikelihoodFormDTO, existingLikelihood.RiskMatrixID, &id); err != nil {
		return err
	}

	existingLikelihood.TitleEn = riskMatrixLikelihoodFormDTO.TitleEn
	existingLikelihood.TitleAr = riskMatrixLikelihoodFormDTO.TitleAr
	existingLikelihood.DescriptionEn = riskMatrixLikelihoodFormDTO.DescriptionEn
	existingLikelihood.DescriptionAr = riskMatrixLikelihoodFormDTO.DescriptionAr
	existingLikelihood.Status = riskMatrixLikelihoodFormDTO.Status

	err = rmis.repo.Update(existingLikelihood)

	if err != nil {
		return err
	}

	return nil

}

func (rmis *RiskMatrixLikelihoodService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
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

func (rmis *RiskMatrixLikelihoodService) Activate(ids []uint, currentUserID uint) error {
	successfulIds, err := rmis.repo.Activate(ids)

	if len(successfulIds) > 0 {

		go func() {
			err := rmis.activityLogSvc.LogBulkActivities(
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

func (rmis *RiskMatrixLikelihoodService) Deactivate(ids []uint, currentUserID uint) error {
	successfulIds, err := rmis.repo.Deactivate(ids)

	if len(successfulIds) > 0 {

		go func() {
			err := rmis.activityLogSvc.LogBulkActivities(
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

func (rmis *RiskMatrixLikelihoodService) validateRiskMatrixLikelihoodFormDTO(riskMatrixLikelihoodFormDTO *riskDTOs.RiskMatrixLikelihoodFormDTO, riskMatrixID uint, id *uint) error {

	if id == nil {
		if exists, err := rmis.repo.ExistsByTitle(riskMatrixLikelihoodFormDTO.TitleEn, riskMatrixLikelihoodFormDTO.TitleAr, riskMatrixID); err != nil {
			return err
		} else if exists {
			return utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
		}

	} else {
		if exists, err := rmis.repo.ExistsByTitleExcludingID(*id, riskMatrixLikelihoodFormDTO.TitleEn, riskMatrixLikelihoodFormDTO.TitleAr, riskMatrixID); err != nil {
			return err
		} else if exists {
			return utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
		}
	}

	return nil
}
