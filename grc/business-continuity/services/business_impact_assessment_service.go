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
	organizationDTOs "jethur/grc/organization/dtos"
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
	ExportBia(data []*organizationDTOs.ProcessListDTO, currentUserID uint) (*os.File, error)
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
	biaProcessSvc             BiaProcessServiceInterface
	biaProcessWorkflowSvc     BiaProcessWorkflowServiceInterface
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
	biaProcessSvc BiaProcessServiceInterface,
	biaProcessWorkflowSvc BiaProcessWorkflowServiceInterface,
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
		biaProcessSvc:             biaProcessSvc,
		biaProcessWorkflowSvc:     biaProcessWorkflowSvc,
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
			existingProcess.RtoFromValue = bcmBiaImpactLevel.RtoFromValue
			existingProcess.RtoToValue = bcmBiaImpactLevel.RtoToValue
			existingProcess.RtoUnitId = bcmBiaImpactLevel.RtoUnitID
			existingProcess.RtoCondition = bcmBiaImpactLevel.RtoCondition
			existingProcess.MtpdFromValue = bcmBiaImpactLevel.MtpdFromValue
			existingProcess.MtpdToValue = bcmBiaImpactLevel.MtpdToValue
			existingProcess.MtpdUnitId = bcmBiaImpactLevel.MtpdUnitID
			existingProcess.MtpdCondition = bcmBiaImpactLevel.MtpdCondition
			existingProcess.RpoFromValue = bcmBiaImpactLevel.RpoFromValue
			existingProcess.RpoToValue = bcmBiaImpactLevel.RpoToValue
			existingProcess.RpoUnitId = bcmBiaImpactLevel.RpoUnitID
			existingProcess.RpoCondition = bcmBiaImpactLevel.RpoCondition
			existingProcess.RspFromValue = bcmBiaImpactLevel.RspFromValue
			existingProcess.RspToValue = bcmBiaImpactLevel.RspToValue
			existingProcess.RspUnitID = bcmBiaImpactLevel.RspUnitID
			existingProcess.RspCondition = bcmBiaImpactLevel.RspCondition
			existingProcess.UpdatedBYID = &currentUserID
			if err := bias.processSvc.UpdateProcessRecoveryObjective(existingProcess, tx); err != nil {
				tx.Rollback()
				return err
			}
		}

	}

	tx.Commit()

	// Create BiaProcess record if not exists
	if businessImpactAssessmentFormDTO.ProcessID != nil {
		if err := bias.biaProcessSvc.CreateIfNotExists(*businessImpactAssessmentFormDTO.ProcessID, currentUserID); err != nil {
			return err
		}
	}

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

	oldData := map[string]interface{}{
		"rto_from_value":  existing.RtoFromValue,
		"override_reason": existing.OverrideReason,
	}

	existing.RtoFromValue = dto.RtoFromValue
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
	process.RtoFromValue = existing.RtoFromValue
	process.RtoToValue = existing.RtoToValue
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

		existingProcess.RtoFromValue = businessImpactAssessmentRecoveryObjectiveFormDTO.Rto
		existingProcess.RtoUnitId = businessImpactAssessmentRecoveryObjectiveFormDTO.RtoUnitId
		existingProcess.MtpdFromValue = businessImpactAssessmentRecoveryObjectiveFormDTO.Mtpd
		existingProcess.MtpdUnitId = businessImpactAssessmentRecoveryObjectiveFormDTO.MtpdUnitId
		existingProcess.RpoFromValue = businessImpactAssessmentRecoveryObjectiveFormDTO.Rpo
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

		// Delete BiaProcess workflows, histories, review comments, then BiaProcess itself
		biaProcess, biaProcessErr := bias.biaProcessSvc.GetByProcessID(id, tx)
		if biaProcessErr == nil && biaProcess != nil {
			if err := bias.biaProcessWorkflowSvc.DeleteAllWorkflowsForBiaProcess(biaProcess.ID, tx); err != nil {
				tx.Rollback()
				return nil, err
			}
			if err := bias.biaProcessSvc.DeleteByProcessID(id, tx); err != nil {
				tx.Rollback()
				return nil, err
			}
		}

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
		existingProcess.RtoFromValue = nil
		existingProcess.RtoToValue = nil
		existingProcess.RtoUnitId = nil
		existingProcess.RtoCondition = nil
		existingProcess.MtpdFromValue = nil
		existingProcess.MtpdToValue = nil
		existingProcess.MtpdUnitId = nil
		existingProcess.MtpdCondition = nil
		existingProcess.RpoFromValue = nil
		existingProcess.RpoToValue = nil
		existingProcess.RpoUnitId = nil
		existingProcess.RpoCondition = nil
		existingProcess.RspFromValue = nil
		existingProcess.RspToValue = nil
		existingProcess.RspUnitID = nil
		existingProcess.RspCondition = nil
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
			"Impact Analysis Level",
			"Risk Level",
			"RTO",
			"MTPD",
			"RPO",
			"RSP",
		}
	} else {
		headers = []string{
			"SL No.",
			"Business Function",
			"Process",
			"Impact Analysis Level",
			"Risk Level",
			"RTO",
			"MTPD",
			"RPO",
			"RSP",
		}
	}

	rows := make([][]interface{}, len(data))
	for i, bia := range data {
		// Helper function to format from/to values
		formatRange := func(fromValue *int, toValue *int, condition *string, unitTitle *string) string {
			if fromValue == nil && toValue == nil {
				return ""
			}
			if fromValue != nil && toValue != nil && *fromValue == *toValue {
				return utility.GetIntValueAsString(*fromValue) + " " + utility.GetStringValue(unitTitle)
			}
			if condition != nil && *condition != "" {
				if *condition == "less_than" {
					return "< " + utility.GetIntValueAsString(*toValue) + " " + utility.GetStringValue(unitTitle)
				} else if *condition == "greater_than" {
					return "> " + utility.GetIntValueAsString(*fromValue) + " " + utility.GetStringValue(unitTitle)
				}
			}
			return utility.GetIntValueAsString(*fromValue) + " - " + utility.GetIntValueAsString(*toValue) + " " + utility.GetStringValue(unitTitle)
		}

		if userSetting.Language == "ar" {
			rows[i] = []interface{}{
				i + 1,
				bia.BusinessFunctionTitle,
				bia.ProcessTitleAr,
				utility.GetStringValue(bia.ImpactAnalysisLevelTitleAr),
				utility.GetStringValue(bia.RiskLevelTitleAr),
				formatRange(bia.ProcessRTOFromValue, bia.ProcessRTOToValue, bia.ProcessRTOCondition, bia.ProcessRTOUnitTitleAr),
				formatRange(bia.ProcessMTPDFromValue, bia.ProcessMTPDToValue, bia.ProcessMTPDCondition, bia.ProcessMTPDUnitTitleAr),
				formatRange(bia.ProcessRPOFromValue, bia.ProcessRPOToValue, bia.ProcessRPOCondition, bia.ProcessRPOUnitTitleAr),
				formatRange(bia.ProcessRSPFromValue, bia.ProcessRSPToValue, bia.ProcessRSPCondition, bia.ProcessRSPUnitTitleAr),
			}
		} else {
			rows[i] = []interface{}{
				i + 1,
				bia.BusinessFunctionTitle,
				bia.ProcessTitleEn,
				utility.GetStringValue(bia.ImpactAnalysisLevelTitleEn),
				utility.GetStringValue(bia.RiskLevelTitleEn),
				formatRange(bia.ProcessRTOFromValue, bia.ProcessRTOToValue, bia.ProcessRTOCondition, bia.ProcessRTOUnitTitleEn),
				formatRange(bia.ProcessMTPDFromValue, bia.ProcessMTPDToValue, bia.ProcessMTPDCondition, bia.ProcessMTPDUnitTitleEn),
				formatRange(bia.ProcessRPOFromValue, bia.ProcessRPOToValue, bia.ProcessRPOCondition, bia.ProcessRPOUnitTitleEn),
				formatRange(bia.ProcessRSPFromValue, bia.ProcessRSPToValue, bia.ProcessRSPCondition, bia.ProcessRSPUnitTitleEn),
			}
		}
	}
	return utility.GenerateExcelFile(headers, rows)
}

func (bias *BusinessImpactAssessmentService) ExportBia(data []*organizationDTOs.ProcessListDTO, currentUserID uint) (*os.File, error) {
	userSetting, err := bias.userSettingSvc.GetUserSetting(currentUserID)
	if err != nil {
		return nil, err
	}

	var headers []string
	if userSetting.Language == "ar" {
		headers = []string{
			"SL No.",
			"Title",
			"Impact Analysis Level",
			"RTO",
			"MTPD",
			"RPO",
			"RSP",
			"Last Assessed At",
			"Last Assessed By",
			"Created At",
			"Created By",
		}
	} else {
		headers = []string{
			"SL No.",
			"Title",
			"Impact Analysis Level",
			"RTO",
			"MTPD",
			"RPO",
			"RSP",
			"Last Assessed At",
			"Last Assessed By",
			"Created At",
			"Created By",
		}
	}

	formatRange := func(fromValue *int, toValue *int, condition *string, unitTitle *string) string {
		if fromValue == nil && toValue == nil {
			return ""
		}
		unit := utility.GetStringValue(unitTitle)
		cond := ""
		if condition != nil {
			cond = *condition
		}
		switch cond {
		case "exact":
			if fromValue != nil {
				return utility.GetIntValueAsString(*fromValue) + " " + unit
			}
		case "range":
			if fromValue != nil && toValue != nil {
				return utility.GetIntValueAsString(*fromValue) + "-" + utility.GetIntValueAsString(*toValue) + " " + unit
			}
			if fromValue != nil {
				return utility.GetIntValueAsString(*fromValue) + " " + unit
			}
		case "less_than":
			if fromValue != nil {
				return utility.GetIntValueAsString(*fromValue) + " " + unit + " or less"
			}
		case "greater_than":
			if fromValue != nil {
				return utility.GetIntValueAsString(*fromValue) + "+ " + unit
			}
		}
		if fromValue != nil {
			return utility.GetIntValueAsString(*fromValue) + " " + unit
		}
		return ""
	}

	rows := make([][]interface{}, len(data))
	for i, process := range data {
		lastAssessedAt := ""
		if process.LastAssessedAt != nil {
			lastAssessedAt = process.LastAssessedAt.Format("2006-01-02")
		}
		createdAt := process.CreatedAt.Format("2006-01-02")
		if userSetting.Language == "ar" {
			rows[i] = []interface{}{
				i + 1,
				process.TitleAr,
				utility.GetStringValue(process.ImpactAnalysisLevelTitleAr),
				formatRange(process.RtoFromValue, process.RtoToValue, process.RtoCondition, process.RtoUnitTitleAr),
				formatRange(process.MtpdFromValue, process.MtpdToValue, process.MtpdCondition, process.MtpdUnitTitleAr),
				formatRange(process.RpoFromValue, process.RpoToValue, process.RpoCondition, process.RpoUnitTitleAr),
				formatRange(process.RspFromValue, process.RspToValue, process.RspCondition, process.RspUnitTitleAr),
				lastAssessedAt,
				fmt.Sprintf("%s %s", utility.GetStringValue(process.LastAssessedUserFirstNameAr), utility.GetStringValue(process.LastAssessedUserLastNameAr)),
				createdAt,
				fmt.Sprintf("%s %s", process.CreatedUserFirstNameAr, process.CreatedUserLastNameAr),
			}
		} else {
			rows[i] = []interface{}{
				i + 1,
				process.TitleEn,
				utility.GetStringValue(process.ImpactAnalysisLevelTitleEn),
				formatRange(process.RtoFromValue, process.RtoToValue, process.RtoCondition, process.RtoUnitTitleEn),
				formatRange(process.MtpdFromValue, process.MtpdToValue, process.MtpdCondition, process.MtpdUnitTitleEn),
				formatRange(process.RpoFromValue, process.RpoToValue, process.RpoCondition, process.RpoUnitTitleEn),
				formatRange(process.RspFromValue, process.RspToValue, process.RspCondition, process.RspUnitTitleEn),
				lastAssessedAt,
				fmt.Sprintf("%s %s", utility.GetStringValue(process.LastAssessedUserFirstNameEn), utility.GetStringValue(process.LastAssessedUserLastNameEn)),
				createdAt,
				fmt.Sprintf("%s %s", process.CreatedUserFirstNameEn, process.CreatedUserLastNameEn),
			}
		}
	}

	return utility.GenerateExcelFile(headers, rows)
}

func (bias *BusinessImpactAssessmentService) GetCriticalSystems() ([]*businessContinuityDTOs.CriticalSystemListDTO, error) {
	return bias.repo.GetCriticalSystems()
}
