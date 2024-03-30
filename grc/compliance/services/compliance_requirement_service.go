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
)

type ComplianceRequirementServiceInterface interface {
	GetAll(params *complianceDTOs.ComplianceRequirementListQueryParamsDTO) (*complianceDTOs.ComplianceRequirementListResponseDTO, error)
	GetByID(complianceID uint, id uint) (*complianceModels.ComplianceRequirement, error)
	Create(complianceID uint, currentUserID uint, ComplianceRequirementFormDTO *complianceDTOs.ComplianceRequirementFormDTO) (uint, error)
	Update(complianceID uint, id uint, currentUserID uint, compliance *complianceDTOs.ComplianceRequirementFormDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)
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

func (crs ComplianceRequirementService) GetByID(complianceID uint, id uint) (*complianceModels.ComplianceRequirement, error) {
	return crs.repo.GetByID(complianceID, id)
}

func (crs *ComplianceRequirementService) Create(complianceID uint, currentUserID uint, ComplianceRequirementFormDTO *complianceDTOs.ComplianceRequirementFormDTO) (uint, error) {

	if err := crs.complianceRegisterSvc.IsValidID(complianceID); err != nil {
		return 0, err
	}

	if err := crs.validateComplianceRequirementFormDTO(ComplianceRequirementFormDTO); err != nil {
		return 0, err
	}

	dates, err := crs.parseComplianceRequirementDates(ComplianceRequirementFormDTO)
	if err != nil {
		return 0, err
	}

	complianceClause := &complianceModels.ComplianceClause{
		ComplianceRegisterID: complianceID,
		ClauseNo:             ComplianceRequirementFormDTO.ComplianceClauseNo,
		Title:                ComplianceRequirementFormDTO.ComplianceClauseTitle,
	}

	complianceRequirement := &complianceModels.ComplianceRequirement{
		ComplianceRequirementCategoryID: ComplianceRequirementFormDTO.ComplianceRequirementCategoryID,
		ComplianceStatusID:              ComplianceRequirementFormDTO.ComplianceStatusID,
		ResponsibleUserID:               ComplianceRequirementFormDTO.ResponsibleUserID,
		SeverityLevelID:                 ComplianceRequirementFormDTO.SeverityLevelID,
		Title:                           ComplianceRequirementFormDTO.Title,
		Deadline:                        dates.Deadline,
		CreatedBy:                       currentUserID,
	}

	id, err := crs.repo.Create(complianceClause, complianceRequirement)
	if err != nil {
		return 0, err
	}

	go func() {
		err := crs.activityLogSvc.LogActivity(
			logModels.CreateAction,
			complianceModels.ComplianceRequirement{}.TableName(),
			id,
			ComplianceRequirementFormDTO,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    logModels.CreateAction,
				"tableName": complianceModels.ComplianceRequirement{}.TableName(),
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
	existingComplianceRequirement, err := crs.repo.GetByID(complianceID, id)

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

	existingComplianceRequirement.ComplianceRequirementCategoryID = ComplianceRequirementFormDTO.ComplianceRequirementCategoryID
	existingComplianceRequirement.ComplianceStatusID = ComplianceRequirementFormDTO.ComplianceStatusID
	existingComplianceRequirement.ResponsibleUserID = ComplianceRequirementFormDTO.ResponsibleUserID
	existingComplianceRequirement.SeverityLevelID = ComplianceRequirementFormDTO.SeverityLevelID
	existingComplianceRequirement.Title = ComplianceRequirementFormDTO.Title
	existingComplianceRequirement.Deadline = dates.Deadline
	existingComplianceRequirement.UpdatedBy = &currentUserID

	err = crs.repo.Update(existingComplianceRequirement)

	if err != nil {
		return err
	}

	return nil
}

func (crs *ComplianceRequirementService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIds []uint
	var successfulIds []uint
	for _, id := range ids {

		_, err := crs.repo.GetByID(id, id)

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
				complianceModels.ComplianceRequirement{}.TableName(),
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

func (crs *ComplianceRequirementService) parseComplianceRequirementDates(dto *complianceDTOs.ComplianceRequirementFormDTO) (*complianceDTOs.ParsedDates, error) {
	parsed := &complianceDTOs.ParsedDates{}

	var err error
	if dto.Deadline != "" {
		parsed.Deadline, err = utility.ParseDate(dto.Deadline)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid deadline"), "invalid_date", "deadline")
		}
	}

	return parsed, nil
}
