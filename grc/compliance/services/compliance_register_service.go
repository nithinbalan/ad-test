package services

import (
	"fmt"
	complianceDTOs "jethur/grc/compliance/dtos"
	"jethur/grc/compliance/models"
	"jethur/grc/compliance/repositories"
	logModels "jethur/grc/logs/models"
	"jethur/grc/logs/services"
	masterRepositories "jethur/grc/masters/repositories"
	"jethur/utility"
	"log"
)

type ComplianceRegisterServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string, page int, pageSize int) (*complianceDTOs.ComplianceRegisterListResponseDTO, error)
	GetByID(id uint) (*models.ComplianceRegister, error)
	Create(currentUserID uint, complianceRegFormDTO *complianceDTOs.ComplianceRegisterFormDTO) (uint, error)
	Update(id uint, currentUserID uint, compliance *complianceDTOs.ComplianceRegisterFormDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	Archive(ids []uint, currentUserID uint) ([]uint, error)
	Restore(ids []uint, currentUserID uint) ([]uint, error)
}

// Service and implementations
type ComplianceRegisterService struct {
	repo                   repositories.ComplianceRegisterRepositoryInterface
	activityLogSvc         services.ActivityLogServiceInterface
	complianceCategoryRepo repositories.ComplianceCategoryRepositoryInterface
	complianceStatusRepo   repositories.ComplianceStatusRepositoryInterface
	reviewFrequencyRepo    masterRepositories.ReviewFrequencyRepositoryInterface
	severityLevelRepo      masterRepositories.SeverityLevelRepositoryInterface
}

func NewComplianceRegisterService(
	repo repositories.ComplianceRegisterRepositoryInterface,
	activityLogSvc services.ActivityLogServiceInterface,
	complianceCategoryRepo repositories.ComplianceCategoryRepositoryInterface,
	complianceStatusRepo repositories.ComplianceStatusRepositoryInterface,
	reviewFrequencyRepo masterRepositories.ReviewFrequencyRepositoryInterface,
	severityLevelRepo masterRepositories.SeverityLevelRepositoryInterface) ComplianceRegisterServiceInterface {
	return &ComplianceRegisterService{
		repo:                   repo,
		activityLogSvc:         activityLogSvc,
		complianceCategoryRepo: complianceCategoryRepo,
		complianceStatusRepo:   complianceStatusRepo,
		reviewFrequencyRepo:    reviewFrequencyRepo,
		severityLevelRepo:      severityLevelRepo,
	}
}

func (crs *ComplianceRegisterService) GetAll(sort string, sortBy string, keyword string, status string, page int, pageSize int) (*complianceDTOs.ComplianceRegisterListResponseDTO, error) {
	return crs.repo.GetAll(sort, sortBy, keyword, status, page, pageSize)
}

func (crs ComplianceRegisterService) GetByID(id uint) (*models.ComplianceRegister, error) {
	return crs.repo.GetByID(id)
}

func (crs *ComplianceRegisterService) Create(currentUserID uint, complianceRegFormDTO *complianceDTOs.ComplianceRegisterFormDTO) (uint, error) {

	if complianceRegFormDTO.ComplianceCategoryID != 0 {
		_, err := ccs.complianceCategoryRepo.GetByID(complianceRegFormDTO.ComplianceCategoryID)
		if err != nil {
			return 0, utility.NewValidationError(fmt.Errorf("compliance category does not exist"), "not_exist", "compliance category")
		}
	}

	complianceStatus, err := crs.complianceStatusRepo.GetByID(complianceRegFormDTO.ComplianceStatusID)
	if err != nil {
		return 0, utility.NewValidationError(fmt.Errorf("compliance status does not exist"), "not_exist", "compliance status")
	}

	reviewFrequency, err := crs.reviewFrequencyRepo.GetByID(complianceRegFormDTO.ReviewFrequencyID)
	if err != nil {
		return 0, utility.NewValidationError(fmt.Errorf("review frequency does not exist"), "not_exist", "review frequency")
	}

	severityLevel, err := crs.severityLevelRepo.GetByID(complianceRegFormDTO.SeverityLevelID)
	if err != nil {
		return 0, utility.NewValidationError(fmt.Errorf("severity level does not exist"), "not_exist", "severity level")
	}

	_, _, _ = complianceStatus, reviewFrequency, severityLevel

	// Proceed with creating the compliance record
	id, err := crs.repo.Create(currentUserID, complianceRegFormDTO)
	if err != nil {
		return 0, err
	}

	go func() {
		err := crs.activityLogSvc.LogActivity(
			logModels.CreateAction,
			models.ComplianceRegister{}.TableName(),
			id,
			complianceRegFormDTO,
			currentUserID,
			nil,
		)

		if err != nil {
			// @todo - save to log file
			log.Printf("Failed to log activity: %v", err)
		}
	}()

	return id, nil
}

func (crs *ComplianceRegisterService) Update(id uint, currentUserID uint, complianceRegFormDTO *complianceDTOs.ComplianceRegisterFormDTO) error {

	existingComplianceRegister, err := crs.repo.GetByID(id)

	if err != nil {
		return err
	}

	complianceCategory, err := crs.complianceCategoryRepo.GetByID(complianceRegFormDTO.ComplianceCategoryID)
	if err != nil {
		return utility.NewValidationError(fmt.Errorf("compliance category does not exist"), "not_exist", "compliance category")
	}

	complianceStatus, err := crs.complianceStatusRepo.GetByID(complianceRegFormDTO.ComplianceStatusID)
	if err != nil {
		return utility.NewValidationError(fmt.Errorf("compliance status does not exist"), "not_exist", "compliance status")
	}

	reviewFrequency, err := crs.reviewFrequencyRepo.GetByID(complianceRegFormDTO.ReviewFrequencyID)
	if err != nil {
		return utility.NewValidationError(fmt.Errorf("review frequency does not exist"), "not_exist", "review frequency")
	}

	severityLevel, err := crs.severityLevelRepo.GetByID(complianceRegFormDTO.SeverityLevelID)
	if err != nil {
		return utility.NewValidationError(fmt.Errorf("severity level does not exist"), "not_exist", "severity level")
	}

	_, _, _, _ = complianceCategory, complianceStatus, reviewFrequency, severityLevel

	err = crs.repo.Update(id, currentUserID, complianceRegFormDTO)

	if err != nil {
		return err
	}

	fieldsToCompare := []string{"ComplianceCategoryID",
		"ComplianceStatusID",
		"ResponsibleUnitID",
		"ResponsibleUserID",
		"ReviewFrequencyID",
		"ReviewResponsibleUserID",
		"SeverityLevelID",
		"StandardCode",
		"TitleEn",
		"TitleAr",
		"DescriptionEn",
		"DescriptionAr",
		"Version",
		"Deadline",
		"IssueDate",
		"ExpiryDate",
		"NextReviewDate"}
	changes, err := utility.CompareFields(existingComplianceRegister, complianceRegFormDTO, fieldsToCompare)

	if err != nil {
		log.Println("error :", err.Error())
	}

	changedDataJSON, err := utility.ChangesToJSON(changes)
	if err != nil {
		log.Println("error :", err.Error())
	}

	go func() {
		err := crs.activityLogSvc.LogActivity(
			logModels.UpdateAction,
			models.ComplianceRegister{}.TableName(),
			id,
			changedDataJSON,
			currentUserID,
			nil,
		)
		if err != nil {
			// @todo - save to log file
			fmt.Printf("Failed to log bulk delete activity: %v\n", err)
		}
	}()
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
				models.ComplianceRegister{}.TableName(),
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
				models.ComplianceRegister{}.TableName(),
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
				models.ComplianceRegister{}.TableName(),
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
