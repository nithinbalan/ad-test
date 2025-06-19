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
	UpdateRecoveryObjectives(businessImpactAssessmentRecoveryObjectiveFormDTO businessContinuityDTOs.BusinessImpactAssessmentRecoveryObjectiveFormDTO, currentUserID uint) error
	UpdateProcessCriticality(businessImpactAssessmentProcessCriticalityFormDTO businessContinuityDTOs.BusinessImpactAssessmentProcessCriticalityFormDTO, currentUserID uint) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	GetBIAMatrix(params *businessContinuityDTOs.BIAMatrixQueryParamsDTO) ([]businessContinuityDTOs.BIAMatrixResponseDTO, error)

	ExportBiaMatrix(data []businessContinuityDTOs.BIAMatrixResponseDTO, currentUserID uint) (*os.File, error)
}

// Service and implementations
type BusinessImpactAssessmentService struct {
	repo                 businessImpactAssessmentRepositories.BusinessImpactAssessmentRepositoryInterface
	activityLogSvc       logServices.ActivityLogServiceInterface
	userSvc              organizationServices.UserServiceInterface
	impactCategorySvc    masterServices.ImpactCategoryServiceInterface
	processSvc           organizationServices.ProcessServiceInterface
	assetRegisterSvc     assetServices.AssetRegisterServiceInterface
	documentSvc          documentServices.DocumentServiceInterface
	processDependencySvc organizationServices.ProcessDependencyServiceInterface
	riskLevelSvc         riskServices.RiskLevelServiceInterface
	userSettingSvc       organizationServices.UserSettingServiceInterface
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
) BusinessImpactAssessmentServiceInterface {
	return &BusinessImpactAssessmentService{
		repo:                 repo,
		activityLogSvc:       activityLogSvc,
		userSvc:              userSvc,
		impactCategorySvc:    impactCategorySvc,
		processSvc:           processSvc,
		assetRegisterSvc:     assetRegisterSvc,
		documentSvc:          documentSvc,
		processDependencySvc: processDependencySvc,
		riskLevelSvc:         riskLevelSvc,
		userSettingSvc:       userSettingSvc,
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
			RiskLevelID:      nil,
			Cost:             nil,
			Comments:         nil,
			IsApplicable:     false,
			CreatedByID:      currentUserID,
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
	tx := bias.repo.BeginTransaction()

	var oldBusinessImpactAssessment []map[string]interface{}
	var existingBusinessImpactAssessment []*models.BusinessImpactAssessment
	var id uint

	if businessImpactAssessmentFormDTO.ProcessID != nil {
		id = *businessImpactAssessmentFormDTO.ProcessID
		existingImpactAssessments, err := bias.repo.GetImpactAssessmentsByProcessID(*businessImpactAssessmentFormDTO.ProcessID)
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
				"is_applicable":      impactAssessment.IsApplicable,
			})
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
				RiskLevelID:      impactAssessment.RiskLevelID,
				Cost:             impactAssessment.Cost,
				Comments:         impactAssessment.Comments,
				IsApplicable:     impactAssessment.IsApplicable,
				CreatedByID:      currentUserID,
			}
			existingBusinessImpactAssessment = append(existingBusinessImpactAssessment, businessImpactAssessment)
			err = bias.repo.CreateBusinessImpactAssessment(businessImpactAssessment, tx)
			if err != nil {
				tx.Rollback()
				return err
			}
		}

		// Update the process risk level
		existingProcess, err := bias.processSvc.GetByID(*businessImpactAssessmentFormDTO.ProcessID, tx)
		if err != nil {
			tx.Rollback()
			return err
		}

		existingProcess.RiskLevelId = businessImpactAssessmentFormDTO.RiskLevelID
		existingProcess.ImpactJustification = businessImpactAssessmentFormDTO.ImpactJustification
		if err := bias.processSvc.UpdateProcessRecoveryObjective(existingProcess, tx); err != nil {
			tx.Rollback()
			return err
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

		// Delete existing BusinessImpactAssessments
		if err := bias.repo.DeleteImpactAssessmentsByProcessID(id, tx); err != nil {
			tx.Rollback()
			return nil, err
		}

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
