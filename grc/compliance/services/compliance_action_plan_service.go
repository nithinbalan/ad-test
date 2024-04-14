package services

import (
	"errors"
	"fmt"
	complianceDTOs "jethur/grc/compliance/dtos"
	complianceModels "jethur/grc/compliance/models"
	complianceRepositories "jethur/grc/compliance/repositories"
	logModels "jethur/grc/logs/models"
	logServices "jethur/grc/logs/services"
	masterRepositories "jethur/grc/masters/repositories"
	masterServices "jethur/grc/masters/services"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type ComplianceActionPlanServiceInterface interface {
	GetAll(complianceID uint, params *complianceDTOs.ComplianceActionPlanListQueryParamsDTO) (*complianceDTOs.ComplianceActionPlanListResponseDTO, error)
	GetByID(complianceID uint, id uint) (*complianceModels.ComplianceActionPlan, error)
	GetFullDetailsByID(complianceID uint, id uint) (*complianceDTOs.ComplianceActionPlanDetailDTO, error)
	Create(complianceID uint, currentUserID uint, ComplianceActionPlanFormDTO *complianceDTOs.ComplianceActionPlanFormDTO) (uint, error)
	Update(complianceID uint, id uint, currentUserID uint, compliance *complianceDTOs.ComplianceActionPlanFormDTO) error
	Delete(complianceID uint, ids []uint, currentUserID uint) ([]uint, error)
	StatusUpdate(complianceID uint, id uint, currentUserID uint, ActionPlanStatusID uint) error
}

// Service and implementations
type ComplianceActionPlanService struct {
	repo                     complianceRepositories.ComplianceActionPlanRepositoryInterface
	activityLogSvc           logServices.ActivityLogServiceInterface
	complianceRegisterSvc    ComplianceRegisterServiceInterface
	complianceRequirementSvc ComplianceRequirementServiceInterface
	actionPlanStatusSvc      masterServices.ActionPlanStatusServiceInterface
	severityLevelSvc         masterServices.SeverityLevelServiceInterface
}

func NewComplianceActionPlanService(
	repo complianceRepositories.ComplianceActionPlanRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	complianceRegisterSvc ComplianceRegisterServiceInterface,
	complianceRequirementSvc ComplianceRequirementServiceInterface,
	actionPlanStatusSvc masterServices.ActionPlanStatusServiceInterface,
	severityLevelSvc masterRepositories.SeverityLevelRepositoryInterface) ComplianceActionPlanServiceInterface {
	return &ComplianceActionPlanService{
		repo:                     repo,
		activityLogSvc:           activityLogSvc,
		complianceRegisterSvc:    complianceRegisterSvc,
		complianceRequirementSvc: complianceRequirementSvc,
		actionPlanStatusSvc:      actionPlanStatusSvc,
		severityLevelSvc:         severityLevelSvc,
	}
}

func (caps *ComplianceActionPlanService) GetAll(complianceID uint, params *complianceDTOs.ComplianceActionPlanListQueryParamsDTO) (*complianceDTOs.ComplianceActionPlanListResponseDTO, error) {
	if _, err := caps.complianceRegisterSvc.GetByID(complianceID); err != nil {
		return nil, err
	}
	return caps.repo.GetAll(complianceID, params)
}

func (caps ComplianceActionPlanService) GetByID(complianceID uint, id uint) (*complianceModels.ComplianceActionPlan, error) {
	if _, err := caps.complianceRegisterSvc.GetByID(complianceID); err != nil {
		return nil, err
	}

	complianceActionPlan, err := caps.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return complianceActionPlan, nil
}

func (caps ComplianceActionPlanService) GetFullDetailsByID(complianceID uint, id uint) (*complianceDTOs.ComplianceActionPlanDetailDTO, error) {
	if _, err := caps.complianceRegisterSvc.GetByID(complianceID); err != nil {
		return nil, err
	}

	var complianceActionPlanDTO complianceDTOs.ComplianceActionPlanDetailDTO

	complianceActionPlan, err := caps.repo.GetFullDetailsByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	err = mapstructure.Decode(complianceActionPlan, &complianceActionPlanDTO)
	if err != nil {
		return nil, err
	}

	//todo :- also add extra fields which required in details screen

	return &complianceActionPlanDTO, nil
}

func (caps *ComplianceActionPlanService) Create(complianceID uint, currentUserID uint, ComplianceActionPlanFormDTO *complianceDTOs.ComplianceActionPlanFormDTO) (uint, error) {

	if _, err := caps.complianceRegisterSvc.GetByID(complianceID); err != nil {
		return 0, err
	}

	if err := caps.validateComplianceActionPlanFormDTO(complianceID, ComplianceActionPlanFormDTO); err != nil {
		return 0, err
	}

	dates, err := caps.parseComplianceActionPlanDates(ComplianceActionPlanFormDTO)
	if err != nil {
		return 0, err
	}

	complianceActionPlan := &complianceModels.ComplianceActionPlan{
		ComplianceRegisterID:    complianceID,
		ComplianceRequirementID: ComplianceActionPlanFormDTO.ComplianceRequirementID,
		Title:                   ComplianceActionPlanFormDTO.Title,
		Description:             ComplianceActionPlanFormDTO.Description,
		SeverityLevelID:         ComplianceActionPlanFormDTO.SeverityLevelID,
		ResponsibleUserID:       ComplianceActionPlanFormDTO.ResponsibleUserID,
		TargetDate:              dates.TargetDate,
		ActionPlanStatusID:      ComplianceActionPlanFormDTO.ActionPlanStatusID,
		CreatedBy:               currentUserID,
	}

	id, err := caps.repo.Create(complianceActionPlan)
	if err != nil {
		return 0, err
	}

	go func() {
		err := caps.activityLogSvc.LogActivity(
			logModels.CreateAction,
			complianceModels.ComplianceActionPlan{}.TableName(),
			id,
			ComplianceActionPlanFormDTO,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    logModels.CreateAction,
				"tableName": complianceModels.ComplianceActionPlan{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return id, nil
}

func (caps *ComplianceActionPlanService) Update(complianceID uint, id uint, currentUserID uint, ComplianceActionPlanFormDTO *complianceDTOs.ComplianceActionPlanFormDTO) error {

	if _, err := caps.complianceRegisterSvc.GetByID(complianceID); err != nil {
		return err
	}

	existingComplianceActionPlan, err := caps.GetByID(complianceID, id)

	if err != nil {
		return err
	}

	if err := caps.validateComplianceActionPlanFormDTO(complianceID, ComplianceActionPlanFormDTO); err != nil {
		return err
	}

	dates, err := caps.parseComplianceActionPlanDates(ComplianceActionPlanFormDTO)
	if err != nil {
		return err
	}

	// Update Compliance ActionPlan
	existingComplianceActionPlan.Title = ComplianceActionPlanFormDTO.Title
	existingComplianceActionPlan.Description = ComplianceActionPlanFormDTO.Description
	existingComplianceActionPlan.SeverityLevelID = ComplianceActionPlanFormDTO.SeverityLevelID
	existingComplianceActionPlan.ResponsibleUserID = ComplianceActionPlanFormDTO.ResponsibleUserID
	existingComplianceActionPlan.TargetDate = dates.TargetDate
	existingComplianceActionPlan.Percentage = ComplianceActionPlanFormDTO.Percentage
	existingComplianceActionPlan.Comments = ComplianceActionPlanFormDTO.Comments
	existingComplianceActionPlan.ActionPlanStatusID = ComplianceActionPlanFormDTO.ActionPlanStatusID
	existingComplianceActionPlan.UpdatedBy = &currentUserID

	err = caps.repo.Update(existingComplianceActionPlan)

	if err != nil {
		return err
	}

	return nil
}

func (caps *ComplianceActionPlanService) Delete(complianceID uint, ids []uint, currentUserID uint) ([]uint, error) {
	if _, err := caps.complianceRegisterSvc.GetByID(complianceID); err != nil {
		return nil, err
	}

	var failedIds []uint
	var successfulIds []uint
	for _, id := range ids {

		_, err := caps.repo.GetByID(id)

		if err != nil {
			continue
		}

		err = caps.repo.Delete(id)
		if err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}
	}

	if len(successfulIds) > 0 {

		go func() {
			err := caps.activityLogSvc.LogBulkActivities(
				logModels.DeleteAction,
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

func (caps *ComplianceActionPlanService) StatusUpdate(complianceID uint, id uint, currentUserID uint, ActionPlanStatusID uint) error {
	existingComplianceActionPlan, err := caps.GetByID(complianceID, id)

	if err != nil {
		return err
	}

	existingComplianceActionPlan.ActionPlanStatusID = ActionPlanStatusID
	existingComplianceActionPlan.UpdatedBy = &currentUserID

	err = caps.repo.Update(existingComplianceActionPlan)

	if err != nil {
		return err
	}

	return nil
}

func (caps *ComplianceActionPlanService) validateComplianceActionPlanFormDTO(complianceID uint, ComplianceActionPlanFormDTO *complianceDTOs.ComplianceActionPlanFormDTO) error {

	if ComplianceActionPlanFormDTO.SeverityLevelID != nil {
		if _, err := caps.severityLevelSvc.GetByID(*ComplianceActionPlanFormDTO.SeverityLevelID); err != nil {
			return utility.NewValidationError(fmt.Errorf("severity level does not exist"), "not_exist", "severity level")
		}
	}

	if ComplianceActionPlanFormDTO.ComplianceRequirementID != nil {
		if _, err := caps.complianceRequirementSvc.GetByID(complianceID, *ComplianceActionPlanFormDTO.ComplianceRequirementID); err != nil {
			return utility.NewValidationError(fmt.Errorf("compliance requirement does not exist"), "not_exist", "compliance requirement")
		}
	}

	if _, err := caps.actionPlanStatusSvc.GetByID(ComplianceActionPlanFormDTO.ActionPlanStatusID); err != nil {
		return utility.NewValidationError(fmt.Errorf("action plan status does not exist"), "not_exist", "action plan status")
	}

	return nil
}

func (caps *ComplianceActionPlanService) parseComplianceActionPlanDates(dto *complianceDTOs.ComplianceActionPlanFormDTO) (*complianceDTOs.ComplianceActionPlanParsedDates, error) {
	parsed := &complianceDTOs.ComplianceActionPlanParsedDates{}

	var err error
	if *dto.TargetDate != "" {
		parsed.TargetDate, err = utility.ParseDate(*dto.TargetDate)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid target date"), "invalid_date", "target date")
		}
	}

	return parsed, nil
}
