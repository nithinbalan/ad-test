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

type ComplianceServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string, page int, pageSize int) (*complianceDTOs.ComplianceRegisterListResponseDTO, error)
	GetByID(id uint) (*models.ComplianceRegister, error)
	Create(currentUserID uint, complianceRegFormDTO *complianceDTOs.ComplianceRegisterFormDTO) (uint, error)
	Update(id uint, currentUserID uint, compliance *complianceDTOs.ComplianceRegisterFormDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	Activate(ids []uint, currentUserID uint) ([]uint, error)
	Archive(ids []uint, currentUserID uint) ([]uint, error)
}

// Service and implementations
type ComplianceService struct {
	repo                   repositories.ComplianceRepositoryInterface
	activityLogSvc         services.ActivityLogServiceInterface
	complianceCategoryRepo repositories.ComplianceCategoryRepositoryInterface
	complianceStatusRepo   repositories.ComplianceStatusRepositoryInterface
	reviewFrequencyRepo    masterRepositories.ReviewFrequencyRepositoryInterface
	severityLevelRepo      masterRepositories.SeverityLevelRepositoryInterface
}

func NewComplianceService(
	repo repositories.ComplianceRepositoryInterface,
	activityLogSvc services.ActivityLogServiceInterface,
	complianceCategoryRepo repositories.ComplianceCategoryRepositoryInterface,
	complianceStatusRepo repositories.ComplianceStatusRepositoryInterface,
	reviewFrequencyRepo masterRepositories.ReviewFrequencyRepositoryInterface,
	severityLevelRepo masterRepositories.SeverityLevelRepositoryInterface) ComplianceServiceInterface {
	return &ComplianceService{
		repo:                   repo,
		activityLogSvc:         activityLogSvc,
		complianceCategoryRepo: complianceCategoryRepo,
		complianceStatusRepo:   complianceStatusRepo,
		reviewFrequencyRepo:    reviewFrequencyRepo,
		severityLevelRepo:      severityLevelRepo,
	}
}

func (ccs *ComplianceService) GetAll(sort string, sortBy string, keyword string, status string, page int, pageSize int) (*complianceDTOs.ComplianceRegisterListResponseDTO, error) {
	return ccs.repo.GetAll(sort, sortBy, keyword, status, page, pageSize)
}

func (ccs ComplianceService) GetByID(id uint) (*models.ComplianceRegister, error) {
	return ccs.repo.GetByID(id)
}

func (ccs *ComplianceService) Create(currentUserID uint, complianceRegFormDTO *complianceDTOs.ComplianceRegisterFormDTO) (uint, error) {

	complianceCategory, err := ccs.complianceCategoryRepo.GetByID(complianceRegFormDTO.ComplianceCategoryID)
	if err != nil {
		return 0, utility.NewValidationError(fmt.Errorf("compliance category does not exist"), "not_exist", "compliance category")
	}

	complianceStatus, err := ccs.complianceStatusRepo.GetByID(complianceRegFormDTO.ComplianceStatusID)
	if err != nil {
		return 0, utility.NewValidationError(fmt.Errorf("compliance status does not exist"), "not_exist", "compliance status")
	}

	reviewFrequency, err := ccs.reviewFrequencyRepo.GetByID(complianceRegFormDTO.ReviewFrequencyID)
	if err != nil {
		return 0, utility.NewValidationError(fmt.Errorf("review frequency does not exist"), "not_exist", "review frequency")
	}

	severityLevel, err := ccs.severityLevelRepo.GetByID(complianceRegFormDTO.SeverityLevelID)
	if err != nil {
		return 0, utility.NewValidationError(fmt.Errorf("severity level does not exist"), "not_exist", "severity level")
	}

	_, _, _, _ = complianceCategory, complianceStatus, reviewFrequency, severityLevel

	// Proceed with creating the compliance record
	id, err := ccs.repo.Create(currentUserID, complianceRegFormDTO)
	if err != nil {
		return 0, err
	}

	go func() {
		err := ccs.activityLogSvc.LogActivity(
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

func (ccs *ComplianceService) Update(id uint, currentUserID uint, complianceRegFormDTO *complianceDTOs.ComplianceRegisterFormDTO) error {

	existingComplianceRegister, err := ccs.repo.GetByID(id)

	if err != nil {
		return err
	}

	complianceCategory, err := ccs.complianceCategoryRepo.GetByID(complianceRegFormDTO.ComplianceCategoryID)
	if err != nil {
		return utility.NewValidationError(fmt.Errorf("compliance category does not exist"), "not_exist", "compliance category")
	}

	complianceStatus, err := ccs.complianceStatusRepo.GetByID(complianceRegFormDTO.ComplianceStatusID)
	if err != nil {
		return utility.NewValidationError(fmt.Errorf("compliance status does not exist"), "not_exist", "compliance status")
	}

	reviewFrequency, err := ccs.reviewFrequencyRepo.GetByID(complianceRegFormDTO.ReviewFrequencyID)
	if err != nil {
		return utility.NewValidationError(fmt.Errorf("review frequency does not exist"), "not_exist", "review frequency")
	}

	severityLevel, err := ccs.severityLevelRepo.GetByID(complianceRegFormDTO.SeverityLevelID)
	if err != nil {
		return utility.NewValidationError(fmt.Errorf("severity level does not exist"), "not_exist", "severity level")
	}

	_, _, _, _ = complianceCategory, complianceStatus, reviewFrequency, severityLevel

	err = ccs.repo.Update(id, currentUserID, complianceRegFormDTO)

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
		err := ccs.activityLogSvc.LogActivity(
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

func (ccs *ComplianceService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIds []uint
	var successfulIds []uint
	for _, id := range ids {

		_, err := ccs.repo.GetByID(id)

		if err != nil {
			continue
		}

		err = ccs.repo.Delete(id)
		if err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}
	}

	if len(successfulIds) > 0 {

		go func() {
			err := ccs.activityLogSvc.LogBulkActivities(
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

func (ccs *ComplianceService) Activate(ids []uint, currentUserID uint) ([]uint, error) {
	successfulIds, err := ccs.repo.Activate(ids)

	if err != nil {
		return nil, err
	}

	if len(successfulIds) > 0 {
		go func() {
			err := ccs.activityLogSvc.LogBulkActivities(
				logModels.ActivateAction,
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

	return successfulIds, nil
}

func (ccs *ComplianceService) Archive(ids []uint, currentUserID uint) ([]uint, error) {
	successfulIds, err := ccs.repo.Archive(ids)

	if err != nil {
		return nil, err
	}

	if len(successfulIds) > 0 {
		go func() {
			err := ccs.activityLogSvc.LogBulkActivities(
				logModels.DeactivateAction,
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

	return successfulIds, nil
}
