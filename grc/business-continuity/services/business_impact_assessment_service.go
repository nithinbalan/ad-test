package services

import (
	"errors"
	"fmt"
	"os"
	"time"

	assetServices "jethur/grc/asset/services"
	businessContinuityDTOs "jethur/grc/business-continuity/dtos"
	businessImpactAssessmentRepositories "jethur/grc/business-continuity/repositories"
	documentServices "jethur/grc/documents/services"
	logServices "jethur/grc/logs/services"
	masterServices "jethur/grc/masters/services"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	riskServices "jethur/grc/risk/services"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type BusinessImpactAssessmentServiceInterface interface {
	GetAllImpactAnalysis(params *businessContinuityDTOs.BusinessImpactAssessmentListQueryParamsDTO) ([]*businessContinuityDTOs.BusinessImpactAssessmentListDTO, error)
	GetAllRecoveryObjectives(params *businessContinuityDTOs.BusinessImpactAssessmentRecoveryObjectiveListQueryParamsDTO) ([]*businessContinuityDTOs.BusinessImpactAssessmentRecoveryObjectiveListDTO, error)
	GetProcessDetails(id uint, tx ...utility.DatabaseTransaction) (*businessContinuityDTOs.BIAProcessDetailDTO, error)
	Create(currentUserID uint, businessContinuityDTOs *businessContinuityDTOs.BusinessImpactFormDTO) (uint, error)
	UpdateImpactAnalysis(businessImpactAssessmentFormDTO businessContinuityDTOs.BusinessImpactAssessmentFormDTO, currentUserID uint) error
	UpdateBIARating(id uint, currentUserID uint, dto businessContinuityDTOs.BusinessImpactAssessmentRatingFormDTO) error
	UpdateRecoveryObjectives(businessImpactAssessmentRecoveryObjectiveFormDTO businessContinuityDTOs.BusinessImpactAssessmentRecoveryObjectiveFormDTO, currentUserID uint) error
	UpdateProcessCriticality(businessImpactAssessmentProcessCriticalityFormDTO businessContinuityDTOs.BusinessImpactAssessmentProcessCriticalityFormDTO, currentUserID uint) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	GetBIAMatrix(params *businessContinuityDTOs.BIAMatrixQueryParamsDTO) ([]businessContinuityDTOs.BIAMatrixResponseDTO, error)

	ExportBiaMatrix(data []businessContinuityDTOs.BIAMatrixResponseDTO, currentUserID uint) (*os.File, error)
	GetCriticalSystems() ([]*businessContinuityDTOs.CriticalSystemListDTO, error)
}

// Service and implementations
type BusinessImpactAssessmentService struct {
	repo                      businessImpactAssessmentRepositories.BusinessImpactAssessmentRepositoryInterface
	activityLogSvc            logServices.ActivityLogServiceInterface
	userSvc                   organizationServices.UserServiceInterface
	impactCategorySvc         masterServices.ImpactCategoryServiceInterface
	processSvc                organizationServices.ProcessServiceInterface
	assetRegisterSvc          assetServices.AssetRegisterServiceInterface
	documentSvc               documentServices.DocumentServiceInterface
	processDependencySvc      organizationServices.ProcessDependencyServiceInterface
	riskLevelSvc              riskServices.RiskLevelServiceInterface
	userSettingSvc            organizationServices.UserSettingServiceInterface
	bcmBiaImpactLevelSvc      BCMBIAImpactLevelServiceInterface
	bcmBiaOutageBucketSvc     BCMBIAOutageBucketServiceInterface
	bcmImpactAnalysisLevelSvc BCMImpactAnalysisLevelServiceInterface
	bcmStaffAllocationSvc     BCMStaffAllocationServiceInterface
}

func NewBusinessImpactAssessmentService(
	repo businessImpactAssessmentRepositories.BusinessImpactAssessmentRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	impactCategorySvc masterServices.ImpactCategoryServiceInterface,
	processSvc organizationServices.ProcessServiceInterface,
	assetRegisterSvc assetServices.AssetRegisterServiceInterface,
	documentSvc documentServices.DocumentServiceInterface,
	processDependencySvc organizationServices.ProcessDependencyServiceInterface,
	riskLevelSvc riskServices.RiskLevelServiceInterface,
	userSettingSvc organizationServices.UserSettingServiceInterface,
	bcmBiaImpactLevelSvc BCMBIAImpactLevelServiceInterface,
	bcmBiaOutageBucketSvc BCMBIAOutageBucketServiceInterface,
	bcmImpactAnalysisLevelSvc BCMImpactAnalysisLevelServiceInterface,
	bcmStaffAllocationSvc BCMStaffAllocationServiceInterface,
) BusinessImpactAssessmentServiceInterface {
	return &BusinessImpactAssessmentService{
		repo:                      repo,
		activityLogSvc:            activityLogSvc,
		userSvc:                   userSvc,
		impactCategorySvc:         impactCategorySvc,
		processSvc:                processSvc,
		assetRegisterSvc:          assetRegisterSvc,
		documentSvc:               documentSvc,
		processDependencySvc:      processDependencySvc,
		riskLevelSvc:              riskLevelSvc,
		userSettingSvc:            userSettingSvc,
		bcmBiaImpactLevelSvc:      bcmBiaImpactLevelSvc,
		bcmBiaOutageBucketSvc:     bcmBiaOutageBucketSvc,
		bcmImpactAnalysisLevelSvc: bcmImpactAnalysisLevelSvc,
		bcmStaffAllocationSvc:     bcmStaffAllocationSvc,
	}
}

func (bias *BusinessImpactAssessmentService) GetAllImpactAnalysis(params *businessContinuityDTOs.BusinessImpactAssessmentListQueryParamsDTO) ([]*businessContinuityDTOs.BusinessImpactAssessmentListDTO, error) {
	return bias.repo.GetAllImpactAnalysis(params)
}

func (bias *BusinessImpactAssessmentService) GetAllRecoveryObjectives(params *businessContinuityDTOs.BusinessImpactAssessmentRecoveryObjectiveListQueryParamsDTO) ([]*businessContinuityDTOs.BusinessImpactAssessmentRecoveryObjectiveListDTO, error) {
	return bias.repo.GetAllRecoveryObjectives(params)
}

func (bias *BusinessImpactAssessmentService) GetProcessDetails(id uint, tx ...utility.DatabaseTransaction) (*businessContinuityDTOs.BIAProcessDetailDTO, error) {
	var processDetailDTO businessContinuityDTOs.BIAProcessDetailDTO

	process, err := bias.repo.GetProcessDetails(id, tx...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	err = mapstructure.Decode(process, &processDetailDTO)
	if err != nil {
		return nil, err
	}

	processDetailDTO.CreatedAt = process.CreatedAt
	processDetailDTO.UpdatedAt = process.UpdatedAt

	return &processDetailDTO, nil
}

func (bias *BusinessImpactAssessmentService) Create(currentUserID uint, businessImpactAssessmentFormDTO *businessContinuityDTOs.BusinessImpactFormDTO) (uint, error) {
	tx := bias.repo.BeginTransaction()

	impactCategories, err := bias.impactCategorySvc.GetAll("", "", "", "active")
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	if len(impactCategories) == 0 {
		return 0, fmt.Errorf("no impact categories found")
	}

	for _, impactCategory := range impactCategories {
		businessImpactAssessment := &models.BusinessImpactAssessment{
			ProcessID:        businessImpactAssessmentFormDTO.ProcessID,
			ImpactCategoryID: impactCategory.ID,
			// New fields
			IsApplicable: false,
			CreatedByID:  currentUserID,
		}
		err := bias.repo.CreateBusinessImpactAssessment(businessImpactAssessment, tx)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	process, err := bias.processSvc.GetFullDetailsByID(*businessImpactAssessmentFormDTO.ProcessID, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	existingProcess, err := bias.processSvc.GetByID(*businessImpactAssessmentFormDTO.ProcessID, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	now := time.Now()
	existingProcess.IsBiaPerformed = true
	existingProcess.LastAssessedByID = &currentUserID
	existingProcess.LastAssessedAt = &now

	err = bias.processSvc.UpdateProcessRecoveryObjective(existingProcess, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	tx.Commit()

	return process.ID, nil
}

func (bias *BusinessImpactAssessmentService) UpdateImpactAnalysis(businessImpactAssessmentFormDTO businessContinuityDTOs.BusinessImpactAssessmentFormDTO, currentUserID uint) error {
	if err := bias.validateNonDecreasingImpact(businessImpactAssessmentFormDTO.BusinessImpactAssessments); err != nil {
		return err
	}

	tx := bias.repo.BeginTransaction()

	var oldBusinessImpactAssessment []map[string]interface{}
	var existingBusinessImpactAssessment []*models.BusinessImpactAssessment
	var id uint

	if businessImpactAssessmentFormDTO.ProcessID != nil {
		id = *businessImpactAssessmentFormDTO.ProcessID
		existingImpactAssessments, err := bias.repo.GetImpactAssessmentsWithAnalysisByProcessID(*businessImpactAssessmentFormDTO.ProcessID)
		if err != nil {
			tx.Rollback()
			return err
		}

		for _, impactAssessment := range existingImpactAssessments {
			oldBusinessImpactAssessment = append(oldBusinessImpactAssessment, map[string]interface{}{
				"process_id":         impactAssessment.ProcessID,
				"impact_category_id": impactAssessment.ImpactCategoryID,
				"is_applicable":      impactAssessment.IsApplicable,
			})
		}

		// Delete existing analysis records explicitly before deleting parent assessments
		if err := bias.repo.DeleteBusinessImpactAssessmentAnalysisByProcessID(*businessImpactAssessmentFormDTO.ProcessID, tx); err != nil {
			tx.Rollback()
			return err
		}

		// Delete existing staff allocations for this process
		if err := bias.bcmStaffAllocationSvc.DeleteByProcessID(*businessImpactAssessmentFormDTO.ProcessID, tx); err != nil {
			tx.Rollback()
			return err
		}

		// Delete existing BusinessImpactAssessments
		if err := bias.repo.DeleteImpactAssessmentsByProcessID(*businessImpactAssessmentFormDTO.ProcessID, tx); err != nil {
			tx.Rollback()
			return err
		}

		for _, impactAssessment := range businessImpactAssessmentFormDTO.BusinessImpactAssessments {
			businessImpactAssessment := &models.BusinessImpactAssessment{
				ProcessID:        businessImpactAssessmentFormDTO.ProcessID,
				ImpactCategoryID: impactAssessment.ImpactCategoryID,
				IsApplicable:     impactAssessment.IsApplicable,
				CreatedByID:      currentUserID,
			}
			err = bias.repo.CreateBusinessImpactAssessment(businessImpactAssessment, tx)
			if err != nil {
				tx.Rollback()
				return err
			}

			// Create analysis records for each outage bucket
			var analysisRecords []*models.BusinessImpactAssessmentAnalysis
			for _, analysis := range impactAssessment.BusinessImpactAssessmentAnalysis {
				analysisRecord := &models.BusinessImpactAssessmentAnalysis{
					BusinessImpactAssessmentID: businessImpactAssessment.ID,
					OutageBucketID:             analysis.OutageBucketID,
					ImpactAnalysisLevelID:      analysis.ImpactAnalysisLevelID,
					CreatedBYID:                currentUserID,
					CreatedAt:                  time.Now(),
				}
				analysisRecords = append(analysisRecords, analysisRecord)
			}

			if len(analysisRecords) > 0 {
				err = bias.repo.CreateBusinessImpactAssessmentAnalysis(analysisRecords, tx)
				if err != nil {
					tx.Rollback()
					return err
				}
			}
		}

		// Create BusinessImpactAssessmentRating if ImpactAnalysisLevelID and OutageBucketID are provided
		// This is process-level, not category-level, so it should be outside the loop
		if businessImpactAssessmentFormDTO.OverallImpactLevelID != nil && businessImpactAssessmentFormDTO.ProcessID != nil {
			// Find BCMBIAImpactLevel using ImpactAnalysisLevelID
			bcmBiaImpactLevel, err := bias.bcmBiaImpactLevelSvc.GetByImpactAnalysisLevelID(*businessImpactAssessmentFormDTO.OverallImpactLevelID)
			if err != nil {
				tx.Rollback()
				return err
			}
			if bcmBiaImpactLevel == nil {
				tx.Rollback()
				return fmt.Errorf("bcm bia impact level not found for impact analysis level ID: %d", *businessImpactAssessmentFormDTO.OverallImpactLevelID)
			}

			// Delete existing rating for this process to ensure single record
			err = bias.repo.DeleteBusinessImpactAssessmentRatingByProcessID(*businessImpactAssessmentFormDTO.ProcessID, tx)
			if err != nil {
				tx.Rollback()
				return err
			}

			// Create rating record with data copied from BCMBIAImpactLevel
			rating := &models.BusinessImpactAssessmentRating{
				ProcessID:            *businessImpactAssessmentFormDTO.ProcessID,
				OverallImpactLevelID: *businessImpactAssessmentFormDTO.OverallImpactLevelID,
				BCMBIAImpactLevelID:  bcmBiaImpactLevel.ID,
				IsEdited:             false,
				RtoFromValue:         bcmBiaImpactLevel.RtoFromValue,
				RtoToValue:           bcmBiaImpactLevel.RtoToValue,
				RtoCondition:         bcmBiaImpactLevel.RtoCondition,
				RtoUnitID:            bcmBiaImpactLevel.RtoUnitID,
				MtpdFromValue:        bcmBiaImpactLevel.MtpdFromValue,
				MtpdToValue:          bcmBiaImpactLevel.MtpdToValue,
				MtpdCondition:        bcmBiaImpactLevel.MtpdCondition,
				MtpdUnitID:           bcmBiaImpactLevel.MtpdUnitID,
				RpoFromValue:         bcmBiaImpactLevel.RpoFromValue,
				RpoToValue:           bcmBiaImpactLevel.RpoToValue,
				RpoCondition:         bcmBiaImpactLevel.RpoCondition,
				RpoUnitID:            bcmBiaImpactLevel.RpoUnitID,
				RspFromValue:         bcmBiaImpactLevel.RspFromValue,
				RspToValue:           bcmBiaImpactLevel.RspToValue,
				RspCondition:         bcmBiaImpactLevel.RspCondition,
				RspUnitID:            bcmBiaImpactLevel.RspUnitID,
				CreatedBYID:          &currentUserID,
			}
			err = bias.repo.CreateBusinessImpactAssessmentRating(rating, tx)
			if err != nil {
				tx.Rollback()
				return err
			}

			// Update the process with impact analysis level, justification, and all RTO/MTPD/RPO fields
			existingProcess, err := bias.processSvc.GetByID(*businessImpactAssessmentFormDTO.ProcessID, tx)
			if err != nil {
				tx.Rollback()
				return err
			}
			existingProcess.ImpactAnalysisLevelID = businessImpactAssessmentFormDTO.OverallImpactLevelID
			existingProcess.ImpactJustification = businessImpactAssessmentFormDTO.ImpactJustification
			existingProcess.Rto = bcmBiaImpactLevel.RtoFromValue
			existingProcess.RtoUnitId = bcmBiaImpactLevel.RtoUnitID
			existingProcess.RtoCondition = bcmBiaImpactLevel.RtoCondition
			existingProcess.Mtpd = bcmBiaImpactLevel.MtpdFromValue
			existingProcess.MtpdUnitId = bcmBiaImpactLevel.MtpdUnitID
			existingProcess.MtpdCondition = bcmBiaImpactLevel.MtpdCondition
			existingProcess.Rpo = bcmBiaImpactLevel.RpoFromValue
			existingProcess.RpoUnitId = bcmBiaImpactLevel.RpoUnitID
			existingProcess.RpoCondition = bcmBiaImpactLevel.RpoCondition
			existingProcess.UpdatedBYID = &currentUserID
			if err := bias.processSvc.UpdateProcessRecoveryObjective(existingProcess, tx); err != nil {
				tx.Rollback()
				return err
			}
		}

	}

	tx.Commit()

	// Auto-generate staff allocations when process is critical (OverallImpactLevelID is not non_applicable)
	if businessImpactAssessmentFormDTO.ProcessID != nil && businessImpactAssessmentFormDTO.OverallImpactLevelID != nil {
		nonApplicableLevel, naErr := bias.bcmImpactAnalysisLevelSvc.GetByType("non_applicable")
		if naErr == nil && nonApplicableLevel != nil && *businessImpactAssessmentFormDTO.OverallImpactLevelID != nonApplicableLevel.ID {
			if err := bias.bcmStaffAllocationSvc.GenerateAllocations(*businessImpactAssessmentFormDTO.ProcessID, currentUserID); err != nil {
				return err
			}
		}
	}

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
				"",
				"",
				"",
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

func (bias *BusinessImpactAssessmentService) validateNonDecreasingImpact(assessments []businessContinuityDTOs.BusinessImpactAssessmentDTO) error {
	for _, impactAssessment := range assessments {
		type bucketLevel struct {
			order int
			score int
		}
		var items []bucketLevel

		for _, analysis := range impactAssessment.BusinessImpactAssessmentAnalysis {
			if analysis.OutageBucketID == nil || analysis.ImpactAnalysisLevelID == nil {
				continue
			}
			bucket, err := bias.bcmBiaOutageBucketSvc.GetByID(*analysis.OutageBucketID)
			if err != nil {
				return utility.NewValidationError(fmt.Errorf("invalid outage_bucket_id: %d", *analysis.OutageBucketID), "not_exist", "outage_bucket_id")
			}
			level, err := bias.bcmImpactAnalysisLevelSvc.GetByID(*analysis.ImpactAnalysisLevelID)
			if err != nil {
				return utility.NewValidationError(fmt.Errorf("invalid impact_analysis_level_id: %d", *analysis.ImpactAnalysisLevelID), "not_exist", "impact_analysis_level_id")
			}
			items = append(items, bucketLevel{order: bucket.Order, score: level.Score})
		}

		// Sort by outage bucket order
		for i := 0; i < len(items); i++ {
			for j := i + 1; j < len(items); j++ {
				if items[i].order > items[j].order {
					items[i], items[j] = items[j], items[i]
				}
			}
		}

		// Validate non-decreasing scores
		for i := 1; i < len(items); i++ {
			if items[i].score < items[i-1].score {
				return utility.NewValidationError(
					fmt.Errorf("impact level must be non-decreasing across outage bucket intervals"),
					"bucket_range",
					"impact_analysis_level_id",
				)
			}
		}
	}
	return nil
}

func (bias *BusinessImpactAssessmentService) UpdateBIARating(id uint, currentUserID uint, dto businessContinuityDTOs.BusinessImpactAssessmentRatingFormDTO) error {
	existing, err := bias.repo.GetRatingByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utility.NewResourceNotFoundError(err, id)
		}
		return err
	}

	if dto.RtoFromValue != nil && dto.RtoToValue != nil && *dto.RtoFromValue > *dto.RtoToValue {
		return utility.NewValidationError(fmt.Errorf("rto_from_value must be less than or equal to rto_to_value"), "invalid_range", "rto_from_value")
	}
	if dto.MtpdFromValue != nil && dto.MtpdToValue != nil && *dto.MtpdFromValue > *dto.MtpdToValue {
		return utility.NewValidationError(fmt.Errorf("mtpd_from_value must be less than or equal to mtpd_to_value"), "invalid_range", "mtpd_from_value")
	}
	if dto.RpoFromValue != nil && dto.RpoToValue != nil && *dto.RpoFromValue > *dto.RpoToValue {
		return utility.NewValidationError(fmt.Errorf("rpo_from_value must be less than or equal to rpo_to_value"), "invalid_range", "rpo_from_value")
	}
	if dto.RspFromValue != nil && dto.RspToValue != nil && *dto.RspFromValue > *dto.RspToValue {
		return utility.NewValidationError(fmt.Errorf("rsp_from_value must be less than or equal to rsp_to_value"), "invalid_range", "rsp_from_value")
	}

	oldData := map[string]interface{}{
		"rto_from_value":  existing.RtoFromValue,
		"rto_to_value":    existing.RtoToValue,
		"mtpd_from_value": existing.MtpdFromValue,
		"mtpd_to_value":   existing.MtpdToValue,
		"rpo_from_value":  existing.RpoFromValue,
		"rpo_to_value":    existing.RpoToValue,
		"rsp_from_value":  existing.RspFromValue,
		"rsp_to_value":    existing.RspToValue,
		"override_reason": existing.OverrideReason,
	}

	existing.RtoFromValue = dto.RtoFromValue
	existing.RtoToValue = dto.RtoToValue
	existing.MtpdFromValue = dto.MtpdFromValue
	existing.MtpdToValue = dto.MtpdToValue
	existing.RpoFromValue = dto.RpoFromValue
	existing.RpoToValue = dto.RpoToValue
	existing.RspFromValue = dto.RspFromValue
	existing.RspToValue = dto.RspToValue
	existing.OverrideReason = dto.OverrideReason
	existing.IsEdited = true
	existing.UpdatedBYID = &currentUserID

	if err := bias.repo.UpdateBusinessImpactAssessmentRating(existing); err != nil {
		return err
	}

	process, err := bias.processSvc.GetByID(existing.ProcessID)
	if err != nil {
		return err
	}
	process.Rto = existing.RtoFromValue
	process.UpdatedBYID = &currentUserID
	tx := bias.repo.BeginTransaction()
	if err := bias.processSvc.UpdateProcessRecoveryObjective(process, tx); err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	go func() {
		changedData := map[string]interface{}{
			"rating_values": map[string]interface{}{
				"old": oldData,
				"new": map[string]interface{}{
					"rto_from_value":  dto.RtoFromValue,
					"rto_to_value":    dto.RtoToValue,
					"mtpd_from_value": dto.MtpdFromValue,
					"mtpd_to_value":   dto.MtpdToValue,
					"rpo_from_value":  dto.RpoFromValue,
					"rpo_to_value":    dto.RpoToValue,
					"rsp_from_value":  dto.RspFromValue,
					"rsp_to_value":    dto.RspToValue,
				},
			},
		}
		err := bias.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.BusinessImpactAssessmentRating{}.TableName(),
			id,
			changedData,
			currentUserID,
			"",
			"",
			"",
			nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.UpdateAction,
				"tableName": models.BusinessImpactAssessmentRating{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}
			utility.LogError("Failed to log activity", fields)
		}
	}()

	return nil
}

func (bias *BusinessImpactAssessmentService) UpdateRecoveryObjectives(businessImpactAssessmentRecoveryObjectiveFormDTO businessContinuityDTOs.BusinessImpactAssessmentRecoveryObjectiveFormDTO, currentUserID uint) error {
	tx := bias.repo.BeginTransaction()

	var oldBusinessImpactAssessment []map[string]interface{}
	var newBusinessImpactAssessment []*models.BusinessImpactAssessmentRecoveryObjective
	var id uint

	if businessImpactAssessmentRecoveryObjectiveFormDTO.ProcessID != nil {
		id = *businessImpactAssessmentRecoveryObjectiveFormDTO.ProcessID

		if err := bias.repo.DeleteRecoveryObjectivesByProcessID(*businessImpactAssessmentRecoveryObjectiveFormDTO.ProcessID, tx); err != nil {
			tx.Rollback()
			return err
		}

		for _, recoveryObjective := range businessImpactAssessmentRecoveryObjectiveFormDTO.RecoveryObjectives {
			if recoveryObjective.ProcessDependencyID != nil {
				existingDependency, err := bias.processDependencySvc.GetByID(*businessImpactAssessmentRecoveryObjectiveFormDTO.ProcessID, *recoveryObjective.ProcessDependencyID)
				if err != nil {
					return err
				}

				existingDependency.Rto = recoveryObjective.Rto
				existingDependency.RtoUnitId = recoveryObjective.RtoUnitId
				existingDependency.Mtpd = recoveryObjective.Mtpd
				existingDependency.MtpdUnitId = recoveryObjective.MtpdUnitId
				existingDependency.Rpo = recoveryObjective.Rpo
				existingDependency.RpoUnitId = recoveryObjective.RpoUnitId
				existingDependency.UpdatedBYID = &currentUserID

				err = bias.processDependencySvc.UpdateProcessDependencyRecoveryObjective(existingDependency, tx)
				if err != nil {
					tx.Rollback()
					return err
				}

				// if existingDependency.AssetID != nil {
				// 	existingAsset, err := bias.assetRegisterSvc.GetByID(*existingDependency.AssetID)
				// 	if err != nil {
				// 		return err
				// 	}

				// 	existingAsset.Rto = recoveryObjective.Rto
				// 	existingAsset.RtoUnitId = recoveryObjective.RtoUnitId
				// 	existingAsset.Mtpd = recoveryObjective.Mtpd
				// 	existingAsset.MtpdUnitId = recoveryObjective.MtpdUnitId
				// 	existingAsset.Rpo = recoveryObjective.Rpo
				// 	existingAsset.RpoUnitId = recoveryObjective.RpoUnitId
				// 	existingAsset.UpdatedBYID = &currentUserID

				// 	err = bias.assetRegisterSvc.UpdateAssetRecoveryObjective(existingAsset, tx)
				// 	if err != nil {
				// 		tx.Rollback()
				// 		return err
				// 	}
				// }
			}
		}

		existingProcess, err := bias.processSvc.GetByID(*businessImpactAssessmentRecoveryObjectiveFormDTO.ProcessID, tx)
		if err != nil {
			tx.Rollback()
			return err
		}

		existingProcess.Rto = businessImpactAssessmentRecoveryObjectiveFormDTO.Rto
		existingProcess.RtoUnitId = businessImpactAssessmentRecoveryObjectiveFormDTO.RtoUnitId
		existingProcess.Mtpd = businessImpactAssessmentRecoveryObjectiveFormDTO.Mtpd
		existingProcess.MtpdUnitId = businessImpactAssessmentRecoveryObjectiveFormDTO.MtpdUnitId
		existingProcess.Rpo = businessImpactAssessmentRecoveryObjectiveFormDTO.Rpo
		existingProcess.RpoUnitId = businessImpactAssessmentRecoveryObjectiveFormDTO.RpoUnitId
		existingProcess.UpdatedBYID = &currentUserID

		err = bias.processSvc.UpdateProcessRecoveryObjective(existingProcess, tx)
		if err != nil {
			tx.Rollback()
			return err
		}

	}

	tx.Commit()

	go func() {
		changedData := make(map[string]interface{})
		if len(oldBusinessImpactAssessment) > 0 || len(newBusinessImpactAssessment) > 0 {
			changedData["impact_assessments"] = map[string]interface{}{
				"added":   newBusinessImpactAssessment,
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
				"",
				"",
				"",
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

func (bias *BusinessImpactAssessmentService) UpdateProcessCriticality(businessImpactAssessmenProcessCriticalityFormDTO businessContinuityDTOs.BusinessImpactAssessmentProcessCriticalityFormDTO, currentUserID uint) error {
	tx := bias.repo.BeginTransaction()

	existingProcess, err := bias.processSvc.GetByID(*businessImpactAssessmenProcessCriticalityFormDTO.ProcessID, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	existingProcess.ProcessCriticalityLevelId = businessImpactAssessmenProcessCriticalityFormDTO.ProcessCriticalityLevelID
	existingProcess.CriticalityJustification = businessImpactAssessmenProcessCriticalityFormDTO.CriticalityJustification
	if err := bias.processSvc.UpdateProcessRecoveryObjective(existingProcess, tx); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil
}

func (bias *BusinessImpactAssessmentService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	tx := bias.repo.BeginTransaction()
	var failedIds []uint
	var successfulIds []uint
	var deletedProcess []map[string]interface{}

	for _, id := range ids {
		process, err := bias.processSvc.GetFullDetailsByID(id, tx)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		deletedProcess = append(deletedProcess, map[string]interface{}{
			"id":    id,
			"title": process.TitleEn,
		})

		// Delete BusinessImpactAssessmentRating first
		if err := bias.repo.DeleteBusinessImpactAssessmentRatingByProcessID(id, tx); err != nil {
			tx.Rollback()
			return nil, err
		}

		// Delete staff allocations
		if err := bias.bcmStaffAllocationSvc.DeleteByProcessID(id, tx); err != nil {
			tx.Rollback()
			return nil, err
		}

		// Delete BusinessImpactAssessmentAnalysis explicitly before parent
		if err := bias.repo.DeleteBusinessImpactAssessmentAnalysisByProcessID(id, tx); err != nil {
			tx.Rollback()
			return nil, err
		}

		// Delete BusinessImpactAssessments
		if err := bias.repo.DeleteImpactAssessmentsByProcessID(id, tx); err != nil {
			tx.Rollback()
			return nil, err
		}

		// Delete RecoveryObjectives
		if err := bias.repo.DeleteRecoveryObjectivesByProcessID(id, tx); err != nil {
			tx.Rollback()
			return nil, err
		}

		existingProcess, err := bias.processSvc.GetByID(id, tx)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		existingProcess.Rto = nil
		existingProcess.RtoUnitId = nil
		existingProcess.Mtpd = nil
		existingProcess.MtpdUnitId = nil
		existingProcess.Rpo = nil
		existingProcess.RpoUnitId = nil
		existingProcess.IsBiaPerformed = false
		existingProcess.RiskLevelId = nil
		existingProcess.ImpactJustification = nil
		existingProcess.ProcessCriticalityLevelId = nil
		existingProcess.CriticalityJustification = nil
		existingProcess.LastAssessedByID = nil

		err = bias.processSvc.UpdateProcessRecoveryObjective(existingProcess, tx)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}
	if len(failedIds) > 0 {
		tx.Rollback()
		return failedIds, nil
	}

	tx.Commit()

	// Log bulk delete activities for successful deletions
	if len(successfulIds) > 0 {
		go func() {
			err := bias.activityLogSvc.LogBulkActivities(
				models.DeleteAction, // Make sure this action constant is defined appropriately
				models.BusinessImpactAssessment{}.TableName(),
				successfulIds,
				deletedProcess,
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

func (bias *BusinessImpactAssessmentService) GetBIAMatrix(params *businessContinuityDTOs.BIAMatrixQueryParamsDTO) ([]businessContinuityDTOs.BIAMatrixResponseDTO, error) {
	return bias.repo.GetBIAMatrix(params)
}

func (bias *BusinessImpactAssessmentService) ExportBiaMatrix(data []businessContinuityDTOs.BIAMatrixResponseDTO, currentUserID uint) (*os.File, error) {
	userSetting, err := bias.userSettingSvc.GetUserSetting(currentUserID)
	if err != nil {
		return nil, err
	}

	var headers []string
	if userSetting.Language == "ar" {
		headers = []string{
			"SL No.",
			"Business Function",
			"Process",
			"Risk Level",
			"RTO",
			"MTPD",
			"RPO",
		}
	} else {
		headers = []string{
			"SL No.",
			"Business Function",
			"Process",
			"Risk Level",
			"RTO",
			"MTPD",
			"RPO",
		}
	}

	rows := make([][]interface{}, len(data))
	for i, bia := range data {
		if userSetting.Language == "ar" {
			rows[i] = []interface{}{
				i + 1,
				bia.BusinessFunctionTitle,
				bia.ProcessTitleAr,
				*bia.RiskLevelTitleAr,
				utility.GetIntValueAsString(bia.ProcessRTO) + " " + utility.GetStringValue(bia.ProcessRTOUnitTitleAr),
				utility.GetIntValueAsString(bia.ProcessMTPD) + " " + utility.GetStringValue(bia.ProcessMTPDUnitTitleAr),
				utility.GetIntValueAsString(bia.ProcessRPO) + " " + utility.GetStringValue(bia.ProcessRPOUnitTitleAr),
			}
		} else {
			rows[i] = []interface{}{
				i + 1,
				bia.BusinessFunctionTitle,
				bia.ProcessTitleEn,
				*bia.RiskLevelTitleEn,
				utility.GetIntValueAsString(bia.ProcessRTO) + " " + utility.GetStringValue(bia.ProcessRTOUnitTitleEn),
				utility.GetIntValueAsString(bia.ProcessMTPD) + " " + utility.GetStringValue(bia.ProcessMTPDUnitTitleEn),
				utility.GetIntValueAsString(bia.ProcessRPO) + " " + utility.GetStringValue(bia.ProcessRPOUnitTitleEn),
			}
		}
	}
	return utility.GenerateExcelFile(headers, rows)
}

func (bias *BusinessImpactAssessmentService) GetCriticalSystems() ([]*businessContinuityDTOs.CriticalSystemListDTO, error) {
	return bias.repo.GetCriticalSystems()
}
