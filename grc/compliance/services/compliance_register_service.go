package services

import (
	"fmt"
	complianceDTOs "jethur/grc/compliance/dtos"
	complianceModels "jethur/grc/compliance/models"
	complianceRepositories "jethur/grc/compliance/repositories"
	logModels "jethur/grc/logs/models"
	logServices "jethur/grc/logs/services"
	masterRepositories "jethur/grc/masters/repositories"
	masterServices "jethur/grc/masters/services"
	"jethur/utility"
	"time"
)

type ComplianceRegisterServiceInterface interface {
	GetAll(params *complianceDTOs.ComplianceRegisterListQueryParamsDTO) (*complianceDTOs.ComplianceRegisterListResponseDTO, error)
	GetByID(id uint) (*complianceModels.ComplianceRegister, error)
	Create(currentUserID uint, complianceRegisterFormDTO *complianceDTOs.ComplianceRegisterFormDTO) (uint, error)
	Update(id uint, currentUserID uint, compliance *complianceDTOs.ComplianceRegisterFormDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	Archive(ids []uint, currentUserID uint) ([]uint, error)
	Restore(ids []uint, currentUserID uint) ([]uint, error)
}

// Service and implementations
type ComplianceRegisterService struct {
	repo                  complianceRepositories.ComplianceRegisterRepositoryInterface
	activityLogSvc        logServices.ActivityLogServiceInterface
	lifecycleStatusSvc    masterServices.LifecycleStatusServiceInterface
	complianceCategorySvc ComplianceCategoryServiceInterface
	complianceStatusSvc   ComplianceStatusServiceInterface
	reviewFrequencySvc    masterServices.ReviewFrequencyServiceInterface
	severityLevelSvc      masterServices.SeverityLevelServiceInterface
}

func NewComplianceRegisterService(
	repo complianceRepositories.ComplianceRegisterRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	lifecycleStatusSvc masterServices.LifecycleStatusServiceInterface,
	complianceCategorySvc ComplianceCategoryServiceInterface,
	complianceStatusSvc ComplianceStatusServiceInterface,
	reviewFrequencySvc masterServices.ReviewFrequencyServiceInterface,
	severityLevelSvc masterRepositories.SeverityLevelRepositoryInterface) ComplianceRegisterServiceInterface {
	return &ComplianceRegisterService{
		repo:                  repo,
		activityLogSvc:        activityLogSvc,
		lifecycleStatusSvc:    lifecycleStatusSvc,
		complianceCategorySvc: complianceCategorySvc,
		complianceStatusSvc:   complianceStatusSvc,
		reviewFrequencySvc:    reviewFrequencySvc,
		severityLevelSvc:      severityLevelSvc,
	}
}

func (crs *ComplianceRegisterService) GetAll(params *complianceDTOs.ComplianceRegisterListQueryParamsDTO) (*complianceDTOs.ComplianceRegisterListResponseDTO, error) {
	return crs.repo.GetAll(params)
}

func (crs ComplianceRegisterService) GetByID(id uint) (*complianceModels.ComplianceRegister, error) {
	return crs.repo.GetByID(id)
}

func (crs *ComplianceRegisterService) Create(currentUserID uint, complianceRegisterFormDTO *complianceDTOs.ComplianceRegisterFormDTO) (uint, error) {

	if err := crs.validateComplianceRegisterFormDTO(complianceRegisterFormDTO); err != nil {
		return 0, err
	}

	dates, err := crs.parseComplianceRegisterDates(complianceRegisterFormDTO)
	if err != nil {
		return 0, err
	}

	referenceNumber, err := crs.repo.GenerateReferenceNumber()
	if err != nil {
		return 0, err
	}

	lifecycleStatus, err := crs.lifecycleStatusSvc.GetByStatus("active")
	if err != nil {
		return 0, err
	}

	currentTime := time.Now().UTC()

	complianceRegister := &complianceModels.ComplianceRegister{
		ComplianceCategoryID:    complianceRegisterFormDTO.ComplianceCategoryID,
		ComplianceStatusID:      complianceRegisterFormDTO.ComplianceStatusID,
		ResponsibleUnitID:       complianceRegisterFormDTO.ResponsibleUnitID,
		ResponsibleUserID:       complianceRegisterFormDTO.ResponsibleUserID,
		ReviewFrequencyID:       complianceRegisterFormDTO.ReviewFrequencyID,
		ReviewResponsibleUserID: complianceRegisterFormDTO.ReviewResponsibleUserID,
		SeverityLevelID:         complianceRegisterFormDTO.SeverityLevelID,
		ReferenceNumber:         referenceNumber,
		StandardCode:            complianceRegisterFormDTO.StandardCode,
		TitleEn:                 complianceRegisterFormDTO.TitleEn,
		TitleAr:                 complianceRegisterFormDTO.TitleAr,
		DescriptionEn:           complianceRegisterFormDTO.DescriptionEn,
		DescriptionAr:           complianceRegisterFormDTO.DescriptionAr,
		Version:                 complianceRegisterFormDTO.Version,
		Deadline:                dates.Deadline,
		IssueDate:               dates.IssueDate,
		ExpiryDate:              dates.ExpiryDate,
		NextReviewDate:          dates.NextReviewDate,
		CreatedBy:               currentUserID,
		CreatedAt:               currentTime,
		UpdatedAt:               currentTime,
		LifecycleStatusID:       lifecycleStatus.ID,
	}

	id, err := crs.repo.Create(complianceRegister)
	if err != nil {
		return 0, err
	}

	go func() {
		err := crs.activityLogSvc.LogActivity(
			logModels.CreateAction,
			complianceModels.ComplianceRegister{}.TableName(),
			id,
			complianceRegisterFormDTO,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    logModels.CreateAction,
				"tableName": complianceModels.ComplianceRegister{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return id, nil
}

func (crs *ComplianceRegisterService) Update(id uint, currentUserID uint, complianceRegisterFormDTO *complianceDTOs.ComplianceRegisterFormDTO) error {

	existingComplianceRegister, err := crs.repo.GetByID(id)

	if err != nil {
		return err
	}

	if err := crs.validateComplianceRegisterFormDTO(complianceRegisterFormDTO); err != nil {
		return err
	}

	dates, err := crs.parseComplianceRegisterDates(complianceRegisterFormDTO)
	if err != nil {
		return err
	}

	existingComplianceRegister.ComplianceCategoryID = complianceRegisterFormDTO.ComplianceCategoryID
	existingComplianceRegister.ComplianceStatusID = complianceRegisterFormDTO.ComplianceStatusID
	existingComplianceRegister.ResponsibleUnitID = complianceRegisterFormDTO.ResponsibleUnitID
	existingComplianceRegister.ResponsibleUserID = complianceRegisterFormDTO.ResponsibleUserID
	existingComplianceRegister.ReviewFrequencyID = complianceRegisterFormDTO.ReviewFrequencyID
	existingComplianceRegister.ReviewResponsibleUserID = complianceRegisterFormDTO.ReviewResponsibleUserID
	existingComplianceRegister.SeverityLevelID = complianceRegisterFormDTO.SeverityLevelID
	existingComplianceRegister.StandardCode = complianceRegisterFormDTO.StandardCode
	existingComplianceRegister.TitleEn = complianceRegisterFormDTO.TitleEn
	existingComplianceRegister.TitleAr = complianceRegisterFormDTO.TitleAr
	existingComplianceRegister.DescriptionEn = complianceRegisterFormDTO.DescriptionEn
	existingComplianceRegister.DescriptionAr = complianceRegisterFormDTO.DescriptionAr
	existingComplianceRegister.Version = complianceRegisterFormDTO.Version
	existingComplianceRegister.Deadline = dates.Deadline
	existingComplianceRegister.IssueDate = dates.IssueDate
	existingComplianceRegister.ExpiryDate = dates.ExpiryDate
	existingComplianceRegister.NextReviewDate = dates.NextReviewDate
	existingComplianceRegister.UpdatedBy = &currentUserID

	err = crs.repo.Update(existingComplianceRegister)

	if err != nil {
		return err
	}

	return nil
}

func (crs *ComplianceRegisterService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIds []uint
	var successfulIds []uint
	for _, id := range ids {

		_, err := crs.repo.GetByID(id)

		if err != nil {
			continue
		}

		err = crs.repo.Delete(id)
		if err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}
	}

	if len(successfulIds) > 0 {

		go func() {
			err := crs.activityLogSvc.LogBulkActivities(
				logModels.DeleteAction, // Make sure this action constant is defined appropriately
				complianceModels.ComplianceRegister{}.TableName(),
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

func (crs *ComplianceRegisterService) Archive(ids []uint, currentUserID uint) ([]uint, error) {
	successfulIds, err := crs.repo.Archive(ids)

	if err != nil {
		return nil, err
	}

	if len(successfulIds) > 0 {
		go func() {
			err := crs.activityLogSvc.LogBulkActivities(
				logModels.ArchiveAction,
				complianceModels.ComplianceRegister{}.TableName(),
				successfulIds,
				nil,
				currentUserID,
				nil,
			)
			if err != nil {
				// @todo - save to log file
				fmt.Printf("Failed to log bulk archive activity: %v\n", err)
			}
		}()
	}

	return successfulIds, nil
}

func (crs *ComplianceRegisterService) Restore(ids []uint, currentUserID uint) ([]uint, error) {
	successfulIds, err := crs.repo.Restore(ids)

	if err != nil {
		return nil, err
	}

	if len(successfulIds) > 0 {
		go func() {
			err := crs.activityLogSvc.LogBulkActivities(
				logModels.RestoreAction,
				complianceModels.ComplianceRegister{}.TableName(),
				successfulIds,
				nil,
				currentUserID,
				nil,
			)
			if err != nil {
				// @todo - save to log file
				fmt.Printf("Failed to log bulk restore activity: %v\n", err)
			}
		}()
	}

	return successfulIds, nil
}

func (crs *ComplianceRegisterService) validateComplianceRegisterFormDTO(complianceRegisterFormDTO *complianceDTOs.ComplianceRegisterFormDTO) error {
	if _, err := crs.complianceCategorySvc.GetByID(complianceRegisterFormDTO.ComplianceCategoryID); err != nil {
		return utility.NewValidationError(fmt.Errorf("compliance category does not exist"), "not_exist", "compliance category")
	}

	if _, err := crs.complianceStatusSvc.GetByID(complianceRegisterFormDTO.ComplianceStatusID); err != nil {
		return utility.NewValidationError(fmt.Errorf("compliance status does not exist"), "not_exist", "compliance status")
	}

	if complianceRegisterFormDTO.ReviewFrequencyID != nil {
		if _, err := crs.reviewFrequencySvc.GetByID(*complianceRegisterFormDTO.ReviewFrequencyID); err != nil {
			return utility.NewValidationError(fmt.Errorf("review frequency does not exist"), "not_exist", "review frequency")
		}
	}
	if complianceRegisterFormDTO.SeverityLevelID != nil {
		if _, err := crs.severityLevelSvc.GetByID(*complianceRegisterFormDTO.SeverityLevelID); err != nil {
			return utility.NewValidationError(fmt.Errorf("severity level does not exist"), "not_exist", "severity level")
		}
	}

	return nil
}

func (crs *ComplianceRegisterService) parseComplianceRegisterDates(dto *complianceDTOs.ComplianceRegisterFormDTO) (*complianceDTOs.ParsedDates, error) {
	parsed := &complianceDTOs.ParsedDates{}

	var err error
	if dto.Deadline != "" {
		parsed.Deadline, err = utility.ParseDate(dto.Deadline)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid deadline"), "invalid_date", "deadline")
		}
	}

	if dto.IssueDate != "" {
		parsed.IssueDate, err = utility.ParseDate(dto.IssueDate)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid issue date"), "invalid_date", "issue date")
		}
	}

	if dto.ExpiryDate != "" {
		parsed.ExpiryDate, err = utility.ParseDate(dto.ExpiryDate)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid expiry date"), "invalid_date", "expiry date")
		}
	}

	if dto.NextReviewDate != "" {
		parsed.NextReviewDate, err = utility.ParseDate(dto.NextReviewDate)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid next review date"), "invalid_date", "next review date")
		}
	}

	return parsed, nil
}
