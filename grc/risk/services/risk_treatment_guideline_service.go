package services

import (
	"errors"
	"fmt"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	riskDTOs "jethur/grc/risk/dtos"
	riskRepositories "jethur/grc/risk/repositories"
	"jethur/utility"

	"gorm.io/gorm"
)

type RiskTreatmentGuidelineServiceInterface interface {
	GetAll(params *riskDTOs.RiskTreatmentGuidelineListQueryParamsDTO) ([]*riskDTOs.RiskTreatmentGuidelineListDTO, error)
	GetByID(id uint) (*models.RiskTreatmentGuideline, error)
	Create(currentUserID uint, guidelineDTO *riskDTOs.RiskTreatmentGuidelineFormDTO) (uint, error)
	Update(id uint, currentUserID uint, guidelineDTO *riskDTOs.RiskTreatmentGuidelineFormDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	Activate(ids []uint, currentUserID uint) error
	Deactivate(ids []uint, currentUserID uint) error
}

// Service and implementations
type RiskTreatmentGuidelineService struct {
	repo                    riskRepositories.RiskTreatmentGuidelineRepositoryInterface
	activityLogSvc          logServices.ActivityLogServiceInterface
	riskMatrixRiskRatingSvc RiskMatrixRiskRatingServiceInterface
}

func NewRiskTreatmentGuidelineService(
	repo riskRepositories.RiskTreatmentGuidelineRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	riskMatrixRiskRatingSvc RiskMatrixRiskRatingServiceInterface,
) RiskTreatmentGuidelineServiceInterface {
	return &RiskTreatmentGuidelineService{
		repo:                    repo,
		activityLogSvc:          activityLogSvc,
		riskMatrixRiskRatingSvc: riskMatrixRiskRatingSvc,
	}
}

func (rgs *RiskTreatmentGuidelineService) GetAll(params *riskDTOs.RiskTreatmentGuidelineListQueryParamsDTO) ([]*riskDTOs.RiskTreatmentGuidelineListDTO, error) {
	return rgs.repo.GetAll(params)
}

func (rgs *RiskTreatmentGuidelineService) GetByID(id uint) (*models.RiskTreatmentGuideline, error) {
	guideline, err := rgs.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return guideline, nil
}

func (rgs *RiskTreatmentGuidelineService) Create(currentUserID uint, guidelineDTO *riskDTOs.RiskTreatmentGuidelineFormDTO) (uint, error) {
	if err := rgs.validateRiskTreatmentGuidelineFormDTO(guidelineDTO, nil); err != nil {
		return 0, err
	}

	guideline := &models.RiskTreatmentGuideline{
		RiskMatrixRiskRatingID: guidelineDTO.RiskMatrixRiskRatingID,
		DescriptionEn:          guidelineDTO.DescriptionEn,
		DescriptionAr:          guidelineDTO.DescriptionAr,
		Status:                 guidelineDTO.Status,
	}

	id, err := rgs.repo.Create(guideline)
	if err != nil {
		return 0, err
	}

	go func() {
		existingRiskTreatmentGuideline, err := rgs.GetByID(id)
		if err != nil {
			utility.LogError("Failed to fetch risk treatment guideline", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		err = rgs.activityLogSvc.LogActivity(
			models.CreateAction,
			models.RiskTreatmentGuideline{}.TableName(),
			id,
			existingRiskTreatmentGuideline,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.RiskTreatmentGuideline{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return id, nil
}

func (rgs *RiskTreatmentGuidelineService) Update(id uint, currentUserID uint, guidelineDTO *riskDTOs.RiskTreatmentGuidelineFormDTO) error {
	oldGuideline, err := rgs.GetByID(id)
	if err != nil {
		return err
	}

	existingGuideline, err := rgs.GetByID(id)

	if err != nil {
		return err
	}

	if err := rgs.validateRiskTreatmentGuidelineFormDTO(guidelineDTO, &id); err != nil {
		return err
	}

	existingGuideline.RiskMatrixRiskRatingID = guidelineDTO.RiskMatrixRiskRatingID
	existingGuideline.DescriptionEn = guidelineDTO.DescriptionEn
	existingGuideline.DescriptionAr = guidelineDTO.DescriptionAr
	existingGuideline.Status = guidelineDTO.Status

	err = rgs.repo.Update(existingGuideline)
	if err != nil {
		return err
	}

	go func() {
		existingTreatmentGuideline, err := rgs.GetByID(id)
		if err != nil {
			return
		}
		//log changes
		changedData := make(map[string]interface{})
		oldData := map[string]interface{}{
			"risk_matrix_risk_rating": oldGuideline.RiskMatrixRiskRating.ID,
			"description_en":          oldGuideline.DescriptionEn,
			"description_ar":          oldGuideline.DescriptionAr,
			"status":                  oldGuideline.Status,
		}

		newData := map[string]interface{}{
			"risk_matrix_risk_rating": existingTreatmentGuideline.RiskMatrixRiskRating.ID,
			"description_en":          existingTreatmentGuideline.DescriptionEn,
			"description_ar":          existingTreatmentGuideline.DescriptionAr,
			"status":                  existingTreatmentGuideline.Status,
		}
		// Compare other fields and add them to changedData
		fieldChanges := utility.CompareFieldChanges(oldData, newData)
		for k, v := range fieldChanges {
			changedData[k] = v
		}

		// Log activity if there are any changes
		if len(changedData) > 0 {
			err := rgs.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.RiskTreatmentGuideline{}.TableName(),
				id,
				changedData,
				currentUserID,
				nil,
			)
			if err != nil {
				utility.LogError("Failed to log activity", map[string]interface{}{
					"action":   models.UpdateAction,
					"error":    err.Error(),
					"recordID": id,
				})
			}
		}

	}()

	return nil
}

func (rgs *RiskTreatmentGuidelineService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIds []uint
	var successfulIds []uint
	var deletedRiskTreatmentGuideline []map[string]interface{}
	for _, id := range ids {

		_, err := rgs.repo.GetByID(id)

		if err != nil {
			continue
		}

		err = rgs.repo.Delete(id)
		if err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}

		deletedRiskTreatmentGuideline = append(deletedRiskTreatmentGuideline, map[string]interface{}{
			"id": id,
		})
	}

	if len(successfulIds) > 0 {

		go func() {
			err := rgs.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.RiskTreatmentGuideline{}.TableName(),
				successfulIds,
				deletedRiskTreatmentGuideline,
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

func (rgs *RiskTreatmentGuidelineService) Activate(ids []uint, currentUserID uint) error {
	successfulIds, err := rgs.repo.Activate(ids)

	if len(successfulIds) > 0 {

		go func() {
			err := rgs.activityLogSvc.LogBulkActivities(
				models.ActivateAction,
				models.RiskTreatmentGuideline{}.TableName(),
				successfulIds,
				nil,
				currentUserID,
				nil,
			)
			if err != nil {
				fmt.Printf("Failed to log bulk activate activity: %v\n", err)
			}
		}()
	}

	return err
}

func (rgs *RiskTreatmentGuidelineService) Deactivate(ids []uint, currentUserID uint) error {
	successfulIds, err := rgs.repo.Deactivate(ids)

	if len(successfulIds) > 0 {

		go func() {
			err := rgs.activityLogSvc.LogBulkActivities(
				models.DeactivateAction,
				models.RiskTreatmentGuideline{}.TableName(),
				successfulIds,
				nil,
				currentUserID,
				nil,
			)
			if err != nil {
				fmt.Printf("Failed to log bulk deactivate activity: %v\n", err)
			}
		}()
	}

	return err
}

func (rgs *RiskTreatmentGuidelineService) validateRiskTreatmentGuidelineFormDTO(guidelineDTO *riskDTOs.RiskTreatmentGuidelineFormDTO, id *uint) error {
	if _, err := rgs.riskMatrixRiskRatingSvc.GetByID(guidelineDTO.RiskMatrixRiskRatingID); err != nil {
		return utility.NewValidationError(fmt.Errorf("risk rating does not exist"), "not_exist", "risk rating")
	}

	if exists, err := rgs.repo.ExistsByRiskMatrixRiskRatingID(guidelineDTO.RiskMatrixRiskRatingID, id); err != nil {
		return err
	} else if exists {
		return utility.NewValidationError(fmt.Errorf("risk matrix risk rating id must be unique"), "unique", "risk matrix risk rating id")
	}

	if id == nil {
		if exists, err := rgs.repo.ExistsByDescription(guidelineDTO.DescriptionEn, guidelineDTO.DescriptionAr); err != nil {
			return err
		} else if exists {
			return utility.NewValidationError(fmt.Errorf("the description must be unique"), "unique", "description")
		}
	} else {
		if exists, err := rgs.repo.ExistsByDescriptionExcludingID(*id, guidelineDTO.DescriptionEn, guidelineDTO.DescriptionAr); err != nil {
			return err
		} else if exists {
			return utility.NewValidationError(fmt.Errorf("the description must be unique"), "unique", "description")
		}
	}

	return nil
}
