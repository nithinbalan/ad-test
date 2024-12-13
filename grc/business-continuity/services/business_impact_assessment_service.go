package services

import (
	businessContinuityDTOs "jethur/grc/business-continuity/dtos"
	businessImpactAssessmentRepositories "jethur/grc/business-continuity/repositories"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	"jethur/utility"
)

type BusinessImpactAssessmentServiceInterface interface {
	GetAll(params *businessContinuityDTOs.BusinessImpactAssessmentListQueryParamsDTO) (*businessContinuityDTOs.BusinessImpactAssessmentListResponseDTO, error)
	Update(businessImpactAssessmentFormDTO businessContinuityDTOs.BusinessImpactAssessmentFormDTO, currentUserID uint) error
}

// Service and implementations
type BusinessImpactAssessmentService struct {
	repo           businessImpactAssessmentRepositories.BusinessImpactAssessmentRepositoryInterface
	activityLogSvc logServices.ActivityLogServiceInterface
	userSvc        organizationServices.UserServiceInterface
}

func NewBusinessImpactAssessmentService(
	repo businessImpactAssessmentRepositories.BusinessImpactAssessmentRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	userSvc organizationServices.UserServiceInterface,
) BusinessImpactAssessmentServiceInterface {
	return &BusinessImpactAssessmentService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
		userSvc:        userSvc,
	}
}

func (bias *BusinessImpactAssessmentService) GetAll(params *businessContinuityDTOs.BusinessImpactAssessmentListQueryParamsDTO) (*businessContinuityDTOs.BusinessImpactAssessmentListResponseDTO, error) {
	return bias.repo.GetAll(params)
}

func (bias *BusinessImpactAssessmentService) Update(businessImpactAssessmentFormDTO businessContinuityDTOs.BusinessImpactAssessmentFormDTO, currentUserID uint) error {
	tx := bias.repo.BeginTransaction()
	var oldBusinessImpactAssessment []map[string]interface{}
	var existingBusinessImpactAssessment []*models.BusinessImpactAssessment
	var id uint
	if businessImpactAssessmentFormDTO.ProcessID != nil {
		id = *businessImpactAssessmentFormDTO.ProcessID
		existingImpactAssessments, err := bias.repo.GetBusinessImpactAssessmentsByProcessID(*businessImpactAssessmentFormDTO.ProcessID)
		if err != nil {
			return err
		}

		for _, impactAssessment := range existingImpactAssessments {
			oldBusinessImpactAssessment = append(oldBusinessImpactAssessment, map[string]interface{}{
				"process_id":         impactAssessment.ProcessID,
				"impact_category_id": impactAssessment.ImpactCategoryID,
				"risk_level_id":      impactAssessment.RiskLevelID,
				"cost":               impactAssessment.Cost,
				"comments":           impactAssessment.Comments,
			})
		}

		// Delete existing BusinessImpactAssessments
		if err := bias.repo.DeleteBusinessImpactAssessmentsByProcessID(*businessImpactAssessmentFormDTO.ProcessID, tx); err != nil {
			tx.Rollback()
			return err
		}

		for _, impactAssessment := range businessImpactAssessmentFormDTO.BusinessImpactAssessments {
			businessImpactAssessment := &models.BusinessImpactAssessment{
				ProcessID:        businessImpactAssessmentFormDTO.ProcessID,
				ImpactCategoryID: impactAssessment.ImpactCategoryID,
				RiskLevelID:      impactAssessment.RiskLevelID,
				Cost:             impactAssessment.Cost,
				Comments:         impactAssessment.Comments,
				CreatedByID:      currentUserID,
			}
			existingBusinessImpactAssessment = append(existingBusinessImpactAssessment, businessImpactAssessment)
			err := bias.repo.Create(businessImpactAssessment, tx)
			if err != nil {
				tx.Rollback()
				return err
			}

		}
	}

	if businessImpactAssessmentFormDTO.AssetRegisterID != nil {
		id = *businessImpactAssessmentFormDTO.AssetRegisterID
		existingImpactAssessments, err := bias.repo.GetBusinessImpactAssessmentsByAssetRegisterID(*businessImpactAssessmentFormDTO.AssetRegisterID)
		if err != nil {
			return err
		}
		for _, impactAssessment := range existingImpactAssessments {
			oldBusinessImpactAssessment = append(oldBusinessImpactAssessment, map[string]interface{}{
				"asset_register_id":  impactAssessment.AssetRegisterID,
				"impact_category_id": impactAssessment.ImpactCategory,
				"risk_level_id":      impactAssessment.RiskLevelID,
				"cost":               impactAssessment.Cost,
				"comments":           impactAssessment.Comments,
			})
		}
		// Delete existing BusinessImpactAssessments
		if err := bias.repo.DeleteBusinessImpactAssessmentsByAssetRegisterID(*businessImpactAssessmentFormDTO.AssetRegisterID, tx); err != nil {
			tx.Rollback()
			return err
		}
		for _, impactAssessment := range businessImpactAssessmentFormDTO.BusinessImpactAssessments {
			businessImpactAssessment := &models.BusinessImpactAssessment{
				AssetRegisterID:  businessImpactAssessmentFormDTO.AssetRegisterID,
				ImpactCategoryID: impactAssessment.ImpactCategoryID,
				RiskLevelID:      impactAssessment.RiskLevelID,
				Cost:             impactAssessment.Cost,
				Comments:         impactAssessment.Comments,
				CreatedByID:      currentUserID,
			}
			existingBusinessImpactAssessment = append(existingBusinessImpactAssessment, businessImpactAssessment)
			err := bias.repo.Create(businessImpactAssessment, tx)
			if err != nil {
				tx.Rollback()
				return err
			}

		}
	}
	tx.Commit()
	go func() {
		changedData := make(map[string]interface{})
		if len(oldBusinessImpactAssessment) > 0 || len(existingBusinessImpactAssessment) > 0 {
			changedData["impact_assessments"] = map[string]interface{}{
				"added":   existingBusinessImpactAssessment,
				"removed": oldBusinessImpactAssessment,
			}
		}
		// Log activity if there are any changes
		if len(changedData) > 0 {
			err := bias.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.BusinessImpactAssessment{}.TableName(),
				id,
				changedData,
				currentUserID,
				nil,
			)

			if err != nil {
				fields := map[string]interface{}{
					"action":    models.UpdateAction,
					"tableName": models.BusinessImpactAssessment{}.TableName(),
					"recordID":  id,
					"userID":    currentUserID,
					"error":     err.Error(),
				}

				utility.LogError("Failed to log activity", fields)
			}
		}
	}()
	return nil
}
