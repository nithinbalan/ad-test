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

type RiskMatrixRiskRatingServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*riskDTOs.RiskMatrixRiskRatingDTO, error)
	GetByID(id uint) (*models.RiskMatrixRiskRating, error)
	Create(currentUserID uint, riskMatrixRiskRatingFormDTO *riskDTOs.RiskMatrixRiskRatingFormDTO) (uint, error)
	Update(id uint, currentUserID uint, riskMatrixRiskRatingFormDTO *riskDTOs.RiskMatrixRiskRatingFormDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	Activate(ids []uint, currentUserID uint) error
	Deactivate(ids []uint, currentUserID uint) error
}

// Service and implementations
type RiskMatrixRiskRatingService struct {
	repo           repositories.RiskMatrixRiskRatingRepositoryInterface
	activityLogSvc activityLogServices.ActivityLogServiceInterface
	riskMatrixSvc  RiskMatrixServiceInterface
}

func NewRiskMatrixRiskRatingService(repo repositories.RiskMatrixRiskRatingRepositoryInterface, activityLogSvc activityLogServices.ActivityLogServiceInterface, riskMatrixService RiskMatrixServiceInterface) RiskMatrixRiskRatingServiceInterface {
	return &RiskMatrixRiskRatingService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
		riskMatrixSvc:  riskMatrixService,
	}
}

func (rmrrs *RiskMatrixRiskRatingService) GetAll(sort string, sortBy string, keyword string, status string) ([]*riskDTOs.RiskMatrixRiskRatingDTO, error) {
	return rmrrs.repo.GetAll(sort, sortBy, keyword, status)
}

func (rmrrs RiskMatrixRiskRatingService) GetByID(id uint) (*models.RiskMatrixRiskRating, error) {
	return rmrrs.repo.GetByID(id)
}

func (rmrrs *RiskMatrixRiskRatingService) Create(currentUserID uint, riskMatrixRiskRatingFormDTO *riskDTOs.RiskMatrixRiskRatingFormDTO) (uint, error) {
	riskMatrix, err := rmrrs.riskMatrixSvc.GetByType("erm")
	if err != nil {
		return 0, err
	}

	if err := rmrrs.validateRiskMatrixRiskRatingFormDTO(riskMatrixRiskRatingFormDTO, riskMatrix.ID, nil); err != nil {
		return 0, err
	}

	riskRating := &models.RiskMatrixRiskRating{
		RiskMatrixID: riskMatrix.ID,
		TitleEn:      riskMatrixRiskRatingFormDTO.TitleEn,
		TitleAr:      riskMatrixRiskRatingFormDTO.TitleAr,
		Type:         riskMatrixRiskRatingFormDTO.Type,
		ColorCode:    riskMatrixRiskRatingFormDTO.ColorCode,
		Status:       riskMatrixRiskRatingFormDTO.Status,
	}

	id, err := rmrrs.repo.Create(riskRating)
	if err != nil {
		return 0, err
	}

	go func() {
		err := rmrrs.activityLogSvc.LogActivity(
			logModels.CreateAction,
			models.RiskMatrixRiskRating{}.TableName(),
			id,
			riskMatrixRiskRatingFormDTO,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    logModels.CreateAction,
				"tableName": models.RiskMatrixRiskRating{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return id, nil

}

func (rmrrs *RiskMatrixRiskRatingService) Update(id uint, currentUserID uint, riskMatrixRiskRatingFormDTO *riskDTOs.RiskMatrixRiskRatingFormDTO) error {
	existingRiskRating, err := rmrrs.GetByID(id)

	if err != nil {
		return err
	}

	if err := rmrrs.validateRiskMatrixRiskRatingFormDTO(riskMatrixRiskRatingFormDTO, existingRiskRating.RiskMatrixID, &id); err != nil {
		return err
	}

	existingRiskRating.TitleEn = riskMatrixRiskRatingFormDTO.TitleEn
	existingRiskRating.TitleAr = riskMatrixRiskRatingFormDTO.TitleAr
	existingRiskRating.Type = riskMatrixRiskRatingFormDTO.Type
	existingRiskRating.ColorCode = riskMatrixRiskRatingFormDTO.ColorCode
	existingRiskRating.Status = riskMatrixRiskRatingFormDTO.Status

	err = rmrrs.repo.Update(existingRiskRating)

	if err != nil {
		return err
	}

	return nil

}

func (rmrrs *RiskMatrixRiskRatingService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIds []uint
	var successfulIds []uint
	for _, id := range ids {

		_, err := rmrrs.repo.GetByID(id)

		if err != nil {
			fmt.Print(err)
			continue
		}

		err = rmrrs.repo.Delete(id)
		if err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}
	}

	if len(successfulIds) > 0 {

		go func() {
			err := rmrrs.activityLogSvc.LogBulkActivities(
				logModels.DeleteAction, // Make sure this action constant is defined appropriately
				models.RiskMatrixRiskRating{}.TableName(),
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

func (rmrrs *RiskMatrixRiskRatingService) Activate(ids []uint, currentUserID uint) error {
	successfulIds, err := rmrrs.repo.Activate(ids)

	if len(successfulIds) > 0 {

		go func() {
			err := rmrrs.activityLogSvc.LogBulkActivities(
				logModels.ActivateAction,
				models.RiskMatrixRiskRating{}.TableName(),
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

func (rmrrs *RiskMatrixRiskRatingService) Deactivate(ids []uint, currentUserID uint) error {
	successfulIds, err := rmrrs.repo.Deactivate(ids)

	if len(successfulIds) > 0 {

		go func() {
			err := rmrrs.activityLogSvc.LogBulkActivities(
				logModels.DeactivateAction,
				models.RiskMatrixRiskRating{}.TableName(),
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

func (rmrrs *RiskMatrixRiskRatingService) validateRiskMatrixRiskRatingFormDTO(riskMatrixRiskRatingFormDTO *riskDTOs.RiskMatrixRiskRatingFormDTO, riskMatrixID uint, id *uint) error {
	if id == nil {
		if exists, err := rmrrs.repo.ExistsByTitle(riskMatrixRiskRatingFormDTO.TitleEn, riskMatrixRiskRatingFormDTO.TitleAr, riskMatrixID); err != nil {
			return err
		} else if exists {
			return utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
		}

		if exists, err := rmrrs.repo.ExistsByType(riskMatrixRiskRatingFormDTO.Type, riskMatrixID); err != nil {
			return err
		} else if exists {
			return utility.NewValidationError(fmt.Errorf("the type must be unique"), "unique", "type")
		}

	} else {
		if exists, err := rmrrs.repo.ExistsByTitleExcludingID(*id, riskMatrixRiskRatingFormDTO.TitleEn, riskMatrixRiskRatingFormDTO.TitleAr, riskMatrixID); err != nil {
			return err
		} else if exists {
			return utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
		}

		if exists, err := rmrrs.repo.ExistsByTypeExcludingID(*id, riskMatrixRiskRatingFormDTO.Type, riskMatrixID); err != nil {
			return err
		} else if exists {
			return utility.NewValidationError(fmt.Errorf("the type must be unique"), "unique", "type")
		}
	}

	return nil
}
