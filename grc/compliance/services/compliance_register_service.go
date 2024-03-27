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
	"log"
)

type ComplianceRegisterServiceInterface interface {
	GetAll(params *complianceDTOs.ComplianceRegisterListQueryParamsDTO) (*complianceDTOs.ComplianceRegisterListResponseDTO, error)
	GetByID(id uint) (*complianceModels.ComplianceRegister, error)
	Create(currentUserID uint, complianceRegFormDTO *complianceDTOs.ComplianceRegisterFormDTO) (uint, error)
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

func (crs *ComplianceRegisterService) Create(currentUserID uint, complianceRegFormDTO *complianceDTOs.ComplianceRegisterFormDTO) (uint, error) {

	if err := crs.validateComplianceRegisterFormDTO(complianceRegFormDTO); err != nil {
		return 0, err
	}

	dates, err := crs.parseComplianceRegisterDates(complianceRegFormDTO)
	if err != nil {
		return 0, err
	}

	referenceNumber, err := crs.repo.GenerateReferenceNumber()
	if err != nil {
		return 0, err
	}

	lifecycleStatus, err := crs.lifecycleStatusSvc.GetByStatus("active")
	if err != nil {
		log.Println("Error parsing next review date:", err)
		return 0, err
	}

	complianceRegister := &complianceModels.ComplianceRegister{
		ComplianceCategoryID:    complianceRegFormDTO.ComplianceCategoryID,
		ComplianceStatusID:      complianceRegFormDTO.ComplianceStatusID,
		ResponsibleUnitID:       complianceRegFormDTO.ResponsibleUnitID,
		ResponsibleUserID:       complianceRegFormDTO.ResponsibleUserID,
		ReviewFrequencyID:       complianceRegFormDTO.ReviewFrequencyID,
		ReviewResponsibleUserID: complianceRegFormDTO.ReviewResponsibleUserID,
		SeverityLevelID:         complianceRegFormDTO.SeverityLevelID,
		ReferenceNumber:         referenceNumber,
		StandardCode:            complianceRegFormDTO.StandardCode,
		TitleEn:                 complianceRegFormDTO.TitleEn,
		TitleAr:                 complianceRegFormDTO.TitleAr,
		DescriptionEn:           complianceRegFormDTO.DescriptionEn,
		DescriptionAr:           complianceRegFormDTO.DescriptionAr,
		Version:                 complianceRegFormDTO.Version,
		Deadline:                dates.Deadline,
		IssueDate:               dates.IssueDate,
		ExpiryDate:              dates.ExpiryDate,
		NextReviewDate:          dates.NextReviewDate,
		CreatedBy:               currentUserID,
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
			complianceRegFormDTO,
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

func (crs *ComplianceRegisterService) Update(id uint, currentUserID uint, complianceRegFormDTO *complianceDTOs.ComplianceRegisterFormDTO) error {

	existingComplianceRegister, err := crs.repo.GetByID(id)

	if err != nil {
		return err
	}

	if err := crs.validateComplianceRegisterFormDTO(complianceRegFormDTO); err != nil {
		return err
	}

	dates, err := crs.parseComplianceRegisterDates(complianceRegFormDTO)
	if err != nil {
		return err
	}

	existingComplianceRegister.ComplianceCategoryID = complianceRegFormDTO.ComplianceCategoryID
	existingComplianceRegister.ComplianceStatusID = complianceRegFormDTO.ComplianceStatusID
	existingComplianceRegister.ResponsibleUnitID = complianceRegFormDTO.ResponsibleUnitID
	existingComplianceRegister.ResponsibleUserID = complianceRegFormDTO.ResponsibleUserID
	existingComplianceRegister.ReviewFrequencyID = complianceRegFormDTO.ReviewFrequencyID
	existingComplianceRegister.ReviewResponsibleUserID = complianceRegFormDTO.ReviewResponsibleUserID
	existingComplianceRegister.SeverityLevelID = complianceRegFormDTO.SeverityLevelID
	existingComplianceRegister.StandardCode = complianceRegFormDTO.StandardCode
	existingComplianceRegister.TitleEn = complianceRegFormDTO.TitleEn
	existingComplianceRegister.TitleAr = complianceRegFormDTO.TitleAr
	existingComplianceRegister.DescriptionEn = complianceRegFormDTO.DescriptionEn
	existingComplianceRegister.DescriptionAr = complianceRegFormDTO.DescriptionAr
	existingComplianceRegister.Version = complianceRegFormDTO.Version
	existingComplianceRegister.Deadline = dates.Deadline
	existingComplianceRegister.IssueDate = dates.IssueDate
	existingComplianceRegister.ExpiryDate = dates.ExpiryDate
	existingComplianceRegister.NextReviewDate = dates.NextReviewDate
	existingComplianceRegister.UpdatedBy = currentUserID

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

func (crs *ComplianceRegisterService) validateComplianceRegisterFormDTO(complianceRegFormDTO *complianceDTOs.ComplianceRegisterFormDTO) error {
	if _, err := crs.complianceCategorySvc.GetByID(complianceRegFormDTO.ComplianceCategoryID); err != nil {
		return utility.NewValidationError(fmt.Errorf("compliance category does not exist"), "not_exist", "compliance category")
	}

	if _, err := crs.complianceStatusSvc.GetByID(complianceRegFormDTO.ComplianceStatusID); err != nil {
		return utility.NewValidationError(fmt.Errorf("compliance status does not exist"), "not_exist", "compliance status")
	}

	if complianceRegFormDTO.ReviewFrequencyID != nil {
		if _, err := crs.reviewFrequencySvc.GetByID(*complianceRegFormDTO.ReviewFrequencyID); err != nil {
			return utility.NewValidationError(fmt.Errorf("review frequency does not exist"), "not_exist", "review frequency")
		}
	}
	if complianceRegFormDTO.SeverityLevelID != nil {
		if _, err := crs.severityLevelSvc.GetByID(*complianceRegFormDTO.SeverityLevelID); err != nil {
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
			return nil, fmt.Errorf("invalid deadline: %w", err)
		}
	}

	if dto.IssueDate != "" {
		parsed.IssueDate, err = utility.ParseDate(dto.IssueDate)
		if err != nil {
			return nil, fmt.Errorf("invalid issue date: %w", err)
		}
	}

	if dto.ExpiryDate != "" {
		parsed.ExpiryDate, err = utility.ParseDate(dto.ExpiryDate)
		if err != nil {
			return nil, fmt.Errorf("invalid expiry date: %w", err)
		}
	}

	if dto.NextReviewDate != "" {
		parsed.NextReviewDate, err = utility.ParseDate(dto.NextReviewDate)
		if err != nil {
			return nil, fmt.Errorf("invalid next review date: %w", err)
		}
	}

	return parsed, nil
}
