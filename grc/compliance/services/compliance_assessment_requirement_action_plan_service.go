package services

import (
	"errors"
	"fmt"

	complianceDTOs "jethur/grc/compliance/dtos"
	complianceRepositories "jethur/grc/compliance/repositories"
	logServices "jethur/grc/logs/services"
	masterServices "jethur/grc/masters/services"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	settingService "jethur/grc/settings/services"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type ComplianceAssessmentRequirementActionPlanServiceInterface interface {
	GetAll(complianceAssessmentID uint, requirementID uint, params *complianceDTOs.ComplianceActionPlanListQueryParamsDTO) (*complianceDTOs.ComplianceActionPlanListResponseDTO, error)
	GetByID(complianceAssessmentID uint, requirementID uint, id uint) (*models.ComplianceActionPlan, error)
	GetFullDetailsByID(complianceAssessmentID uint, requirementID uint, id uint, tx ...utility.DatabaseTransaction) (*complianceDTOs.ComplianceAssessmentRequirementActionPlanDetailDTO, error)
	Create(complianceAssessmentID uint, requirementID uint, currentUserID uint, ComplianceActionPlanFormDTO *complianceDTOs.ComplianceActionPlanFormDTO) (uint, error)
	Update(complianceAssessmentID uint, requirementID uint, id uint, currentUserID uint, compliance *complianceDTOs.ComplianceActionPlanFormDTO) error
	Delete(complianceAssessmentID uint, requirementID uint, ids []uint, currentUserID uint) ([]uint, error)
}

// Service and implementations
type ComplianceAssessmentRequirementActionPlanService struct {
	repo                               complianceRepositories.ComplianceActionPlanRepositoryInterface
	activityLogSvc                     logServices.ActivityLogServiceInterface
	complianceAssessmentRequirementSvc ComplianceAssessmentRequirementServiceInterface
	actionPlanStatusSvc                masterServices.ActionPlanStatusServiceInterface
	userSvc                            organizationServices.UserServiceInterface
	organizationUnitSvc                settingService.OrganizationUnitServiceInterface
}

func NewComplianceAssessmentRequirementActionPlanService(
	repo complianceRepositories.ComplianceActionPlanRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	complianceAssessmentRequirementSvc ComplianceAssessmentRequirementServiceInterface,
	actionPlanStatusSvc masterServices.ActionPlanStatusServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	organizationUnitSvc settingService.OrganizationUnitServiceInterface,
) ComplianceAssessmentRequirementActionPlanServiceInterface {
	return &ComplianceAssessmentRequirementActionPlanService{
		repo:                               repo,
		activityLogSvc:                     activityLogSvc,
		complianceAssessmentRequirementSvc: complianceAssessmentRequirementSvc,
		actionPlanStatusSvc:                actionPlanStatusSvc,
		userSvc:                            userSvc,
		organizationUnitSvc:                organizationUnitSvc,
	}
}

func (caraps *ComplianceAssessmentRequirementActionPlanService) GetAll(complianceAssessmentID uint, requirementID uint, params *complianceDTOs.ComplianceActionPlanListQueryParamsDTO) (*complianceDTOs.ComplianceActionPlanListResponseDTO, error) {
	if requirementID != 0 {
		if _, err := caraps.complianceAssessmentRequirementSvc.GetByID(requirementID); err != nil {
			return nil, err
		}
	}
	return caraps.repo.GetAll(uint(0), params, requirementID)
}

func (caraps ComplianceAssessmentRequirementActionPlanService) GetByID(complianceAssessmentID uint, requirementID uint, id uint) (*models.ComplianceActionPlan, error) {
	if _, err := caraps.complianceAssessmentRequirementSvc.GetByID(requirementID); err != nil {
		return nil, err
	}

	complianceAssessmentRequirementActionPlan, err := caraps.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return complianceAssessmentRequirementActionPlan, nil
}

func (caraps ComplianceAssessmentRequirementActionPlanService) GetFullDetailsByID(complianceAssessmentID uint, requirementID uint, id uint, tx ...utility.DatabaseTransaction) (*complianceDTOs.ComplianceAssessmentRequirementActionPlanDetailDTO, error) {
	var complianceAssessmentRequirementActionPlanDTO complianceDTOs.ComplianceAssessmentRequirementActionPlanDetailDTO

	// Check if the compliance register exists
	if _, err := caraps.complianceAssessmentRequirementSvc.GetByID(requirementID); err != nil {
		return nil, err
	}

	// Get full details of the compliance action plan
	complianceAssessmentRequirementActionPlan, err := caraps.repo.GetFullDetailsByID(id, tx...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	// Decode the compliance action plan into the DTO
	err = mapstructure.Decode(complianceAssessmentRequirementActionPlan, &complianceAssessmentRequirementActionPlanDTO)
	if err != nil {
		return nil, err
	}
	if complianceAssessmentRequirementActionPlan.TargetDate != nil {
		targetDate := complianceAssessmentRequirementActionPlan.TargetDate.Format("2006-01-02")
		complianceAssessmentRequirementActionPlanDTO.TargetDateString = &targetDate
	}

	complianceAssessmentRequirementActionPlanDTO.CreatedAt = complianceAssessmentRequirementActionPlan.CreatedAt
	complianceAssessmentRequirementActionPlanDTO.UpdatedAt = complianceAssessmentRequirementActionPlan.UpdatedAt

	// TODO: Add any additional fields required for the details screen

	return &complianceAssessmentRequirementActionPlanDTO, nil
}

func (caraps *ComplianceAssessmentRequirementActionPlanService) Create(complianceAssessmentID uint, requirementID uint, currentUserID uint, ComplianceActionPlanFormDTO *complianceDTOs.ComplianceActionPlanFormDTO) (uint, error) {
	tx := caraps.repo.BeginTransaction()

	if err := caraps.validateComplianceActionPlanFormDTO(ComplianceActionPlanFormDTO); err != nil {
		return 0, err
	}

	dates, err := caraps.parseComplianceAssessmentRequirementActionPlanDates(ComplianceActionPlanFormDTO)
	if err != nil {
		return 0, err
	}

	complianceAssessmentRequirementActionPlan := &models.ComplianceActionPlan{
		ComplianceRegisterID:              nil,
		ComplianceAssessmentRequirementID: &requirementID,
		Title:                             ComplianceActionPlanFormDTO.Title,
		Description:                       ComplianceActionPlanFormDTO.Description,
		ResponsibleUserID:                 ComplianceActionPlanFormDTO.ResponsibleUserID,
		TargetDate:                        dates.TargetDate,
		Percentage:                        ComplianceActionPlanFormDTO.Percentage,
		ActionPlanStatusID:                ComplianceActionPlanFormDTO.ActionPlanStatusID,
		CreatedByID:                       currentUserID,
	}

	id, err := caraps.repo.Create(complianceAssessmentRequirementActionPlan, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	tx.Commit()

	go func() {
		existingComplianceAssessmentRequirementActionPlan, err := caraps.GetFullDetailsByID(complianceAssessmentID, requirementID, id)
		if err != nil {
			utility.LogError("Failed to fetch compliance action plan", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		err = caraps.activityLogSvc.LogActivity(
			models.CreateAction,
			models.ComplianceActionPlan{}.TableName(),
			id,
			existingComplianceAssessmentRequirementActionPlan,
			currentUserID,
			nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.ComplianceActionPlan{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return id, nil
}

func (caraps *ComplianceAssessmentRequirementActionPlanService) Update(complianceAssessmentID uint, requirementID uint, id uint, currentUserID uint, ComplianceActionPlanFormDTO *complianceDTOs.ComplianceActionPlanFormDTO) error {
	tx := caraps.repo.BeginTransaction()

	existingComplianceAssessmentRequirementActionPlan, err := caraps.GetByID(complianceAssessmentID, requirementID, id)
	if err != nil {
		return err
	}

	oldComplianceAssessmentRequirementActionPlan, err := caraps.GetFullDetailsByID(complianceAssessmentID, requirementID, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := caraps.validateComplianceActionPlanFormDTO(ComplianceActionPlanFormDTO); err != nil {
		return err
	}

	dates, err := caraps.parseComplianceAssessmentRequirementActionPlanDates(ComplianceActionPlanFormDTO)
	if err != nil {
		return err
	}

	// Update Compliance ActionPlan
	existingComplianceAssessmentRequirementActionPlan.Title = ComplianceActionPlanFormDTO.Title
	existingComplianceAssessmentRequirementActionPlan.Description = ComplianceActionPlanFormDTO.Description
	existingComplianceAssessmentRequirementActionPlan.ResponsibleUserID = ComplianceActionPlanFormDTO.ResponsibleUserID
	existingComplianceAssessmentRequirementActionPlan.TargetDate = dates.TargetDate
	existingComplianceAssessmentRequirementActionPlan.Percentage = ComplianceActionPlanFormDTO.Percentage
	existingComplianceAssessmentRequirementActionPlan.ActionPlanStatusID = ComplianceActionPlanFormDTO.ActionPlanStatusID
	existingComplianceAssessmentRequirementActionPlan.UpdatedByID = &currentUserID

	err = caraps.repo.Update(existingComplianceAssessmentRequirementActionPlan, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	go func() {
		existingActionPlan, err := caraps.GetFullDetailsByID(complianceAssessmentID, requirementID, id)
		if err != nil {
			return
		}
		//log changes
		changedData := make(map[string]interface{})
		oldData := map[string]interface{}{
			"title":              oldComplianceAssessmentRequirementActionPlan.Title,
			"description":        oldComplianceAssessmentRequirementActionPlan.Description,
			"responsible_user":   fmt.Sprintf("%s %s", oldComplianceAssessmentRequirementActionPlan.ResponsibleUser.FirstNameEN, oldComplianceAssessmentRequirementActionPlan.ResponsibleUser.LastNameEN),
			"target_date":        oldComplianceAssessmentRequirementActionPlan.TargetDateString,
			"percentage":         oldComplianceAssessmentRequirementActionPlan.Percentage,
			"action_plan_status": oldComplianceAssessmentRequirementActionPlan.ActionPlanStatus.TitleEn,
		}

		newData := map[string]interface{}{
			"title":              existingActionPlan.Title,
			"description":        existingActionPlan.Description,
			"responsible_user":   fmt.Sprintf("%s %s", existingActionPlan.ResponsibleUser.FirstNameEN, existingActionPlan.ResponsibleUser.LastNameEN),
			"target_date":        existingActionPlan.TargetDateString,
			"percentage":         existingActionPlan.Percentage,
			"action_plan_status": existingActionPlan.ActionPlanStatus.TitleEn,
		}
		// Compare other fields and add them to changedData
		fieldChanges := utility.CompareFieldChanges(oldData, newData)
		for k, v := range fieldChanges {
			changedData[k] = v
		}
		// Log activity if there are any changes
		if len(changedData) > 0 {
			err := caraps.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.ComplianceActionPlan{}.TableName(),
				id,
				changedData,
				currentUserID,
				nil,
			)

			if err != nil {
				fields := map[string]interface{}{
					"action":    models.UpdateAction,
					"tableName": models.ComplianceActionPlan{}.TableName(),
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

func (caraps *ComplianceAssessmentRequirementActionPlanService) Delete(complianceAssessmentID uint, requirementID uint, ids []uint, currentUserID uint) ([]uint, error) {
	if _, err := caraps.complianceAssessmentRequirementSvc.GetByID(requirementID); err != nil {
		return nil, err
	}

	var failedIds []uint
	var successfulIds []uint
	var deletedComplianceAssessmentRequirementActionPlan []map[string]interface{}
	for _, id := range ids {

		complianceAssessmentRequirementActionPlan, err := caraps.GetFullDetailsByID(complianceAssessmentID, requirementID, id)
		if err != nil {
			continue
		}

		err = caraps.repo.Delete(id)
		if err != nil {
			failedIds = append(failedIds, id)
			continue
		} else {
			successfulIds = append(successfulIds, id)
		}

		deletedComplianceAssessmentRequirementActionPlan = append(deletedComplianceAssessmentRequirementActionPlan, map[string]interface{}{
			"id":                     id,
			"compliance_register_id": requirementID,
			"title":                  complianceAssessmentRequirementActionPlan.Title,
		})

	}

	if len(failedIds) > 0 {
		return failedIds, nil
	}

	if len(successfulIds) > 0 {
		go func() {
			err := caraps.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.ComplianceActionPlan{}.TableName(),
				successfulIds,
				deletedComplianceAssessmentRequirementActionPlan,
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

func (caraps *ComplianceAssessmentRequirementActionPlanService) validateComplianceActionPlanFormDTO(ComplianceActionPlanFormDTO *complianceDTOs.ComplianceActionPlanFormDTO) error {

	if _, err := caraps.actionPlanStatusSvc.GetByID(ComplianceActionPlanFormDTO.ActionPlanStatusID); err != nil {
		return utility.NewValidationError(fmt.Errorf("action plan status does not exist"), "not_exist", "action plan status")
	}

	if ComplianceActionPlanFormDTO.Percentage == nil {
		defaultValue := int(0)
		ComplianceActionPlanFormDTO.Percentage = &defaultValue
	}

	return nil
}

func (caraps *ComplianceAssessmentRequirementActionPlanService) parseComplianceAssessmentRequirementActionPlanDates(dto *complianceDTOs.ComplianceActionPlanFormDTO) (*complianceDTOs.ComplianceAssessmentRequirementActionPlanParsedDates, error) {
	parsed := &complianceDTOs.ComplianceAssessmentRequirementActionPlanParsedDates{}

	var err error
	if *dto.TargetDate != "" {
		parsed.TargetDate, err = utility.ParseDate(*dto.TargetDate)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid target date"), "invalid_date", "target date")
		}
	}

	return parsed, nil
}
