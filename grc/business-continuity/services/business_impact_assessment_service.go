package services

import (
	"fmt"
	"time"

	assetServices "jethur/grc/asset/services"
	businessContinuityDTOs "jethur/grc/business-continuity/dtos"
	businessImpactAssessmentRepositories "jethur/grc/business-continuity/repositories"
	documentServices "jethur/grc/documents/services"
	logServices "jethur/grc/logs/services"
	masterServices "jethur/grc/masters/services"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	"jethur/utility"
)

type BusinessImpactAssessmentServiceInterface interface {
	GetAllImpactAnalysis(params *businessContinuityDTOs.BusinessImpactAssessmentListQueryParamsDTO) ([]*businessContinuityDTOs.BusinessImpactAssessmentListDTO, error)
	GetAllRecoveryObjectives(params *businessContinuityDTOs.BusinessImpactAssessmentRecoveryObjectiveListQueryParamsDTO) ([]*businessContinuityDTOs.BusinessImpactAssessmentRecoveryObjectiveListDTO, error)
	Create(currentUserID uint, businessContinuityDTOs *businessContinuityDTOs.BusinessImpactFormDTO) (uint, error)
	UpdateImpactAnalysis(businessImpactAssessmentFormDTO businessContinuityDTOs.BusinessImpactAssessmentFormDTO, currentUserID uint) error
	UpdateRecoveryObjectives(businessImpactAssessmentRecoveryObjectiveFormDTO businessContinuityDTOs.BusinessImpactAssessmentRecoveryObjectiveFormDTO, currentUserID uint) error
	GetBIAMatrix(params *businessContinuityDTOs.BIAMatrixQueryParamsDTO) (*businessContinuityDTOs.BIAMatrixResponsePaginatedDTO, error)
}

// Service and implementations
type BusinessImpactAssessmentService struct {
	repo              businessImpactAssessmentRepositories.BusinessImpactAssessmentRepositoryInterface
	activityLogSvc    logServices.ActivityLogServiceInterface
	userSvc           organizationServices.UserServiceInterface
	impactCategorySvc masterServices.ImpactCategoryServiceInterface
	processSvc        organizationServices.ProcessServiceInterface
	assetRegisterSvc  assetServices.AssetRegisterServiceInterface
	documentSvc       documentServices.DocumentServiceInterface
}

func NewBusinessImpactAssessmentService(
	repo businessImpactAssessmentRepositories.BusinessImpactAssessmentRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	impactCategorySvc masterServices.ImpactCategoryServiceInterface,
	processSvc organizationServices.ProcessServiceInterface,
	assetRegisterSvc assetServices.AssetRegisterServiceInterface,
	documentSvc documentServices.DocumentServiceInterface,
) BusinessImpactAssessmentServiceInterface {
	return &BusinessImpactAssessmentService{
		repo:              repo,
		activityLogSvc:    activityLogSvc,
		userSvc:           userSvc,
		impactCategorySvc: impactCategorySvc,
		processSvc:        processSvc,
		assetRegisterSvc:  assetRegisterSvc,
		documentSvc:       documentSvc,
	}
}

func (bias *BusinessImpactAssessmentService) GetAllImpactAnalysis(params *businessContinuityDTOs.BusinessImpactAssessmentListQueryParamsDTO) ([]*businessContinuityDTOs.BusinessImpactAssessmentListDTO, error) {
	return bias.repo.GetAllImpactAnalysis(params)
}

func (bias *BusinessImpactAssessmentService) GetAllRecoveryObjectives(params *businessContinuityDTOs.BusinessImpactAssessmentRecoveryObjectiveListQueryParamsDTO) ([]*businessContinuityDTOs.BusinessImpactAssessmentRecoveryObjectiveListDTO, error) {
	return bias.repo.GetAllRecoveryObjectives(params)
}

func (bias *BusinessImpactAssessmentService) Create(currentUserID uint, businessImpactAssessmentFormDTO *businessContinuityDTOs.BusinessImpactFormDTO) (uint, error) {
	tx := bias.repo.BeginTransaction()

	impactCategories, err := bias.impactCategorySvc.GetAll("", "", "", "")
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
			RiskLevelID:      1,
			Cost:             nil,
			Comments:         nil,
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

	if process.Assets != nil {
		for _, asset := range process.Assets {
			recoveryObjective := &models.BusinessImpactAssessmentRecoveryObjective{
				ProcessID:       businessImpactAssessmentFormDTO.ProcessID,
				AssetRegisterID: &asset.Asset.ID,
				Rto:             asset.Asset.Rto,
				RtoUnitId:       asset.Asset.RtoUnitId,
				Mtpd:            asset.Asset.Mtpd,
				MtpdUnitId:      asset.Asset.MtpdUnitId,
				Rpo:             asset.Asset.Rpo,
				RpoUnitId:       asset.Asset.RpoUnitId,
				CreatedByID:     currentUserID,
			}
			err := bias.repo.CreateRecoveryObjective(recoveryObjective, tx)
			if err != nil {
				tx.Rollback()
				return 0, err
			}
		}
	}
	if process.Documents != nil {
		for _, document := range process.Documents {
			recoveryObjective := &models.BusinessImpactAssessmentRecoveryObjective{
				ProcessID:   businessImpactAssessmentFormDTO.ProcessID,
				DocumentID:  &document.DocumentID,
				Rto:         document.Document.Rto,
				RtoUnitId:   document.Document.RtoUnitId,
				Mtpd:        document.Document.Mtpd,
				MtpdUnitId:  document.Document.MtpdUnitId,
				Rpo:         document.Document.Rpo,
				RpoUnitId:   document.Document.RpoUnitId,
				CreatedByID: currentUserID,
			}
			err := bias.repo.CreateRecoveryObjective(recoveryObjective, tx)
			if err != nil {
				tx.Rollback()
				return 0, err
			}
		}
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
				CreatedByID:      currentUserID,
			}
			existingBusinessImpactAssessment = append(existingBusinessImpactAssessment, businessImpactAssessment)
			err = bias.repo.CreateBusinessImpactAssessment(businessImpactAssessment, tx)
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

			ImpactRecoveryObjective := &models.BusinessImpactAssessmentRecoveryObjective{
				ProcessID:         businessImpactAssessmentRecoveryObjectiveFormDTO.ProcessID,
				AssetRegisterID:   recoveryObjective.AssetRegisterID,
				DocumentID:        recoveryObjective.DocumentID,
				RoleInProcess:     recoveryObjective.RoleInProcess,
				Rto:               recoveryObjective.Rto,
				RtoUnitId:         recoveryObjective.RtoUnitId,
				Mtpd:              recoveryObjective.Mtpd,
				MtpdUnitId:        recoveryObjective.MtpdUnitId,
				Rpo:               recoveryObjective.Rpo,
				RpoUnitId:         recoveryObjective.RpoUnitId,
				ImpactDescription: recoveryObjective.ImpactDescription,
				CreatedByID:       currentUserID,
			}

			newBusinessImpactAssessment = append(newBusinessImpactAssessment, ImpactRecoveryObjective)

			err := bias.repo.CreateRecoveryObjective(ImpactRecoveryObjective, tx)
			if err != nil {
				tx.Rollback()
				return err
			}

			if recoveryObjective.AssetRegisterID != nil {
				existingAsset, err := bias.assetRegisterSvc.GetByID(*recoveryObjective.AssetRegisterID)
				if err != nil {
					return err
				}

				existingAsset.Rto = recoveryObjective.Rto
				existingAsset.RtoUnitId = recoveryObjective.RtoUnitId
				existingAsset.Mtpd = recoveryObjective.Mtpd
				existingAsset.MtpdUnitId = recoveryObjective.MtpdUnitId
				existingAsset.Rpo = recoveryObjective.Rpo
				existingAsset.RpoUnitId = recoveryObjective.RpoUnitId
				existingAsset.UpdatedBYID = &currentUserID

				err = bias.assetRegisterSvc.UpdateAssetRecoveryObjective(existingAsset, tx)
				if err != nil {
					tx.Rollback()
					return err
				}
			}

			if recoveryObjective.DocumentID != nil {
				existingDocument, err := bias.documentSvc.GetByID(*recoveryObjective.DocumentID)
				if err != nil {
					return err
				}

				existingDocument.Rto = recoveryObjective.Rto
				existingDocument.RtoUnitId = recoveryObjective.RtoUnitId
				existingDocument.Mtpd = recoveryObjective.Mtpd
				existingDocument.MtpdUnitId = recoveryObjective.MtpdUnitId
				existingDocument.Rpo = recoveryObjective.Rpo
				existingDocument.RpoUnitId = recoveryObjective.RpoUnitId
				existingDocument.UpdatedBYID = &currentUserID

				err = bias.documentSvc.UpdateDocumentRecoveryObjective(existingDocument, tx)
				if err != nil {
					tx.Rollback()
					return err
				}
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

func (bias *BusinessImpactAssessmentService) GetBIAMatrix(params *businessContinuityDTOs.BIAMatrixQueryParamsDTO) (*businessContinuityDTOs.BIAMatrixResponsePaginatedDTO, error) {
	return bias.repo.GetBIAMatrix(params)
}
