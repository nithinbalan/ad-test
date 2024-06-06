package services

import (
	"errors"
	"fmt"
	complianceDTOs "jethur/grc/compliance/dtos"
	complianceRepositories "jethur/grc/compliance/repositories"
	logServices "jethur/grc/logs/services"
	masterRepositories "jethur/grc/masters/repositories"
	masterServices "jethur/grc/masters/services"
	"jethur/grc/models"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type ComplianceRequirementServiceInterface interface {
	GetAll(params *complianceDTOs.ComplianceRequirementListQueryParamsDTO) (*complianceDTOs.ComplianceRequirementListResponseDTO, error)
	GetByID(complianceID uint, id uint) (*models.ComplianceRequirement, error)
	GetRequirementsByComplianceID(complianceID uint, params *complianceDTOs.ComplianceRequirementListQueryParamsDTO) ([]*complianceDTOs.ComplianceRequirementByComplianceIDListResponseDTO, error)
	GetByComplianceRequirementID(id uint) (*models.ComplianceRequirement, error)
	GetFullDetailsByID(complianceID uint, id uint) (*complianceDTOs.ComplianceRequirementDetailDTO, error)
	Create(complianceID uint, currentUserID uint, ComplianceRequirementFormDTO *complianceDTOs.ComplianceRequirementFormDTO) (uint, error)
	Update(complianceID uint, id uint, currentUserID uint, compliance *complianceDTOs.ComplianceRequirementFormDTO) error
	Delete(complianceID uint, ids []uint, currentUserID uint) ([]uint, error)
	StatusUpdate(complianceID uint, id uint, currentUserID uint, complianceStatusID uint, tx utility.DatabaseTransaction) error
}

// Service and implementations
type ComplianceRequirementService struct {
	repo                             complianceRepositories.ComplianceRequirementRepositoryInterface
	activityLogSvc                   logServices.ActivityLogServiceInterface
	complianceRegisterSvc            ComplianceRegisterServiceInterface
	complianceRequirementCategorySvc ComplianceRequirementCategoryServiceInterface
	complianceStatusSvc              ComplianceStatusServiceInterface
	severityLevelSvc                 masterServices.SeverityLevelServiceInterface
}

func NewComplianceRequirementService(
	repo complianceRepositories.ComplianceRequirementRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	complianceRegisterSvc ComplianceRegisterServiceInterface,
	complianceRequirementCategorySvc ComplianceRequirementCategoryServiceInterface,
	complianceStatusSvc ComplianceStatusServiceInterface,
	severityLevelSvc masterRepositories.SeverityLevelRepositoryInterface) ComplianceRequirementServiceInterface {
	return &ComplianceRequirementService{
		repo:                             repo,
		activityLogSvc:                   activityLogSvc,
		complianceRegisterSvc:            complianceRegisterSvc,
		complianceRequirementCategorySvc: complianceRequirementCategorySvc,
		complianceStatusSvc:              complianceStatusSvc,
		severityLevelSvc:                 severityLevelSvc,
	}
}

func (crs *ComplianceRequirementService) GetAll(params *complianceDTOs.ComplianceRequirementListQueryParamsDTO) (*complianceDTOs.ComplianceRequirementListResponseDTO, error) {
	return crs.repo.GetAll(params)
}

func (crs *ComplianceRequirementService) GetRequirementsByComplianceID(complianceID uint, params *complianceDTOs.ComplianceRequirementListQueryParamsDTO) ([]*complianceDTOs.ComplianceRequirementByComplianceIDListResponseDTO, error) {
	return crs.repo.GetRequirementsByComplianceID(complianceID, params)
}

func (crs ComplianceRequirementService) GetByID(complianceID uint, id uint) (*models.ComplianceRequirement, error) {
	if _, err := crs.complianceRegisterSvc.GetByID(complianceID); err != nil {
		return nil, err
	}

	complianceRequirement, err := crs.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return complianceRequirement, nil
}

func (crs ComplianceRequirementService) GetByComplianceRequirementID(id uint) (*models.ComplianceRequirement, error) {

	complianceRequirement, err := crs.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return complianceRequirement, nil
}

func (crs ComplianceRequirementService) GetFullDetailsByID(complianceID uint, id uint) (*complianceDTOs.ComplianceRequirementDetailDTO, error) {
	if _, err := crs.complianceRegisterSvc.GetByID(complianceID); err != nil {
		return nil, err
	}

	var complianceRequirementDTO complianceDTOs.ComplianceRequirementDetailDTO

	complianceRequirement, err := crs.repo.GetFullDetailsByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	// Map the result to DTO
	err = mapstructure.Decode(complianceRequirement, &complianceRequirementDTO)
	if err != nil {
		return nil, err
	}

	//todo :- also add extra fields which required in details screen

	return &complianceRequirementDTO, nil
}

func (crs *ComplianceRequirementService) Create(complianceID uint, currentUserID uint, ComplianceRequirementFormDTO *complianceDTOs.ComplianceRequirementFormDTO) (uint, error) {
	tx := crs.repo.BeginTransaction()

	if _, err := crs.complianceRegisterSvc.GetByID(complianceID); err != nil {
		return 0, err
	}

	if err := crs.validateComplianceRequirementFormDTO(ComplianceRequirementFormDTO); err != nil {
		return 0, err
	}

	dates, err := crs.parseComplianceRequirementDates(ComplianceRequirementFormDTO)
	if err != nil {
		return 0, err
	}

	// Create or retrieve the compliance clause
	complianceClause := &models.ComplianceClause{
		ComplianceRegisterID: complianceID,
		ClauseNo:             ComplianceRequirementFormDTO.ComplianceClauseNo,
		Title:                ComplianceRequirementFormDTO.ComplianceClauseTitle,
	}

	clause, err := crs.repo.GetOrCreateComplianceClause(complianceClause, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	complianceRequirement := &models.ComplianceRequirement{
		ComplianceRequirementCategoryID: ComplianceRequirementFormDTO.ComplianceRequirementCategoryID,
		ComplianceStatusID:              ComplianceRequirementFormDTO.ComplianceStatusID,
		ResponsibleUserID:               ComplianceRequirementFormDTO.ResponsibleUserID,
		SeverityLevelID:                 ComplianceRequirementFormDTO.SeverityLevelID,
		Title:                           ComplianceRequirementFormDTO.Title,
		Deadline:                        dates.Deadline,
		CreatedBy:                       currentUserID,
		ComplianceClauseID:              clause.ID,
	}

	id, err := crs.repo.Create(complianceRequirement, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	tx.Commit()
	go func() {
		err := crs.activityLogSvc.LogActivity(
			models.CreateAction,
			models.ComplianceRequirement{}.TableName(),
			id,
			ComplianceRequirementFormDTO,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.ComplianceRequirement{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return id, nil
}

func (crs *ComplianceRequirementService) Update(complianceID uint, id uint, currentUserID uint, ComplianceRequirementFormDTO *complianceDTOs.ComplianceRequirementFormDTO) error {
	tx := crs.repo.BeginTransaction()

	if _, err := crs.complianceRegisterSvc.GetByID(complianceID); err != nil {
		return err
	}

	existingComplianceRequirement, err := crs.GetByID(complianceID, id)

	if err != nil {
		return err
	}

	if err := crs.validateComplianceRequirementFormDTO(ComplianceRequirementFormDTO); err != nil {
		return err
	}

	dates, err := crs.parseComplianceRequirementDates(ComplianceRequirementFormDTO)
	if err != nil {
		return err
	}

	complianceClause := &models.ComplianceClause{
		ComplianceRegisterID: complianceID,
		ClauseNo:             ComplianceRequirementFormDTO.ComplianceClauseNo,
		Title:                ComplianceRequirementFormDTO.ComplianceClauseTitle,
	}

	clause, err := crs.repo.GetOrCreateComplianceClause(complianceClause)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Update Compliance Requirement
	existingComplianceRequirement.ComplianceRequirementCategoryID = ComplianceRequirementFormDTO.ComplianceRequirementCategoryID
	existingComplianceRequirement.ComplianceStatusID = ComplianceRequirementFormDTO.ComplianceStatusID
	existingComplianceRequirement.ResponsibleUserID = ComplianceRequirementFormDTO.ResponsibleUserID
	existingComplianceRequirement.SeverityLevelID = ComplianceRequirementFormDTO.SeverityLevelID
	existingComplianceRequirement.Title = ComplianceRequirementFormDTO.Title
	existingComplianceRequirement.Deadline = dates.Deadline
	existingComplianceRequirement.UpdatedBy = &currentUserID
	existingComplianceRequirement.ComplianceClauseID = clause.ID

	err = crs.repo.Update(existingComplianceRequirement)

	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (crs *ComplianceRequirementService) Delete(complianceID uint, ids []uint, currentUserID uint) ([]uint, error) {
	tx := crs.repo.BeginTransaction()
	if _, err := crs.complianceRegisterSvc.GetByID(complianceID); err != nil {
		return nil, err
	}

	var failedIds []uint
	var successfulIds []uint

	// Track clause IDs for deletion
	clauseIDsToDelete := make(map[uint]bool)

	for _, id := range ids {
		requirement, err := crs.repo.GetByID(id)
		if err != nil {
			continue
		}

		if err := crs.repo.Delete(id); err != nil {
			tx.Rollback()
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
			// Track the clause IDs
			clauseIDsToDelete[requirement.ComplianceClauseID] = true
		}
	}

	// Delete clauses if all their requirements are deleted
	for clauseID := range clauseIDsToDelete {
		count, err := crs.repo.CountRequirementsByClauseID(clauseID)
		if err != nil {
			// Log error, but proceed with deletion as it might be recoverable later
			fmt.Printf("Error counting requirements for clause %d: %v\n", clauseID, err)
		} else if count == 0 {
			if err := crs.repo.DeleteClause(clauseID); err != nil {
				tx.Rollback()
				// Log error, but proceed with other deletions
				fmt.Printf("Error deleting clause %d: %v\n", clauseID, err)
			}
		}
	}

	// Log successful deletions
	if len(successfulIds) > 0 {
		go func() {
			err := crs.activityLogSvc.LogBulkActivities(
				models.DeleteAction, // Make sure this action constant is defined appropriately
				models.ComplianceRequirement{}.TableName(),
				successfulIds,
				nil,
				currentUserID,
				nil,
			)
			if err != nil {
				// Log error, but proceed as it might be recoverable later
				fmt.Printf("Failed to log bulk delete activity: %v\n", err)
			}
		}()
	}

	tx.Commit()

	return failedIds, nil
}

func (crs *ComplianceRequirementService) StatusUpdate(complianceID uint, id uint, currentUserID uint, complianceStatusID uint, tx utility.DatabaseTransaction) error {
	existingComplianceRequirement, err := crs.GetByID(complianceID, id)

	if err != nil {
		return err
	}

	existingComplianceRequirement.ComplianceStatusID = complianceStatusID
	existingComplianceRequirement.UpdatedBy = &currentUserID

	err = crs.repo.Update(existingComplianceRequirement, tx)

	if err != nil {
		return err
	}

	return nil
}

func (crs *ComplianceRequirementService) validateComplianceRequirementFormDTO(ComplianceRequirementFormDTO *complianceDTOs.ComplianceRequirementFormDTO) error {
	if _, err := crs.complianceRequirementCategorySvc.GetByID(ComplianceRequirementFormDTO.ComplianceRequirementCategoryID); err != nil {
		return utility.NewValidationError(fmt.Errorf("compliance requirement category does not exist"), "not_exist", "compliance requirement category")
	}

	if _, err := crs.complianceStatusSvc.GetByID(ComplianceRequirementFormDTO.ComplianceStatusID); err != nil {
		return utility.NewValidationError(fmt.Errorf("compliance status does not exist"), "not_exist", "compliance status")
	}

	if _, err := crs.severityLevelSvc.GetByID(ComplianceRequirementFormDTO.SeverityLevelID); err != nil {
		return utility.NewValidationError(fmt.Errorf("severity level does not exist"), "not_exist", "severity level")
	}

	return nil
}

func (crs *ComplianceRequirementService) parseComplianceRequirementDates(dto *complianceDTOs.ComplianceRequirementFormDTO) (*complianceDTOs.ComplianceRequirementParsedDates, error) {
	parsed := &complianceDTOs.ComplianceRequirementParsedDates{}

	var err error
	if dto.Deadline != "" {
		parsed.Deadline, err = utility.ParseDate(dto.Deadline)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid deadline"), "invalid_date", "deadline")
		}
	}

	return parsed, nil
}
