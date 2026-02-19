package services

import (
	"fmt"
	"time"

	complianceDTOs "jethur/grc/compliance/dtos"
	complianceRepositories "jethur/grc/compliance/repositories"
	logServices "jethur/grc/logs/services"
	organizationServices "jethur/grc/organization/services"
	settingServices "jethur/grc/settings/services"
	"jethur/utility"

	"jethur/grc/models"
)

type ComplianceAssessmentAcknowledgeDepartmentServiceInterface interface {
	GetAll(assessmentID uint, currentUserID uint, params *complianceDTOs.ComplianceAssessmentAcknowledgeListQueryParamsDTO) (*complianceDTOs.ComplianceAssessmentAcknowledgementListResponseDTO, error)
	Create(assessmentID uint, currentUserID uint, acknowledgeFormDTO *complianceDTOs.ComplianceAssessmentAcknowledgeFormDTO) (uint, error)
	Delete(assessmentID uint, currentUserID uint) (uint, error)
	GetAssessmentAcknowledgmentSummary(assessmentID uint) (*complianceDTOs.ComplianceAssessmentAcknowledgmentSummaryDTO, error)
	GetAssessmentAcknowledgmentSummaryByOrgUnit(assessmentID uint) ([]complianceDTOs.ComplianceAssessmentAcknowledgmentOrgUnitSummaryDTO, error)
}

type ComplianceAssessmentAcknowledgeDepartmentService struct {
	repo                 complianceRepositories.ComplianceAssessmentAcknowledgeDepartmentRepositoryInterface
	activityLogSvc       logServices.ActivityLogServiceInterface
	assessmentSvc        ComplianceAssessmentServiceInterface
	userSvc              organizationServices.UserServiceInterface
	organizationLevelSvc settingServices.OrganizationLevelServiceInterface
}

func NewComplianceAssessmentAcknowledgeDepartmentService(
	repo complianceRepositories.ComplianceAssessmentAcknowledgeDepartmentRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	assessmentSvc ComplianceAssessmentServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	organizationLevelSvc settingServices.OrganizationLevelServiceInterface,
) ComplianceAssessmentAcknowledgeDepartmentServiceInterface {
	return &ComplianceAssessmentAcknowledgeDepartmentService{
		repo:                 repo,
		activityLogSvc:       activityLogSvc,
		assessmentSvc:        assessmentSvc,
		userSvc:              userSvc,
		organizationLevelSvc: organizationLevelSvc,
	}
}

func (caads *ComplianceAssessmentAcknowledgeDepartmentService) GetAll(assessmentID uint, currentUserID uint, params *complianceDTOs.ComplianceAssessmentAcknowledgeListQueryParamsDTO) (*complianceDTOs.ComplianceAssessmentAcknowledgementListResponseDTO, error) {
	// Get assessment details to check if it's public or shared
	assessment, err := caads.assessmentSvc.GetFullDetailsByID(assessmentID, nil)
	if err != nil {
		return nil, err
	}

	// Check if the compliance register has IsRegulationNotificationAssessment = true
	if !assessment.ComplianceRegister.IsRegulationNotificationAssessment {
		return nil, fmt.Errorf("assessment is not a regulation notification assessment")
	}

	// Get organization level settings to determine primary org unit field
	orgLevel, err := caads.organizationLevelSvc.GetPrimaryLevel()
	if err != nil {
		return nil, err
	}

	primaryOrgUnitField := "organization_unit_1_id"
	if orgLevel.Level == 2 {
		primaryOrgUnitField = "organization_unit_2_id"
	} else if orgLevel.Level == 3 {
		primaryOrgUnitField = "organization_unit_3_id"
	}

	// Check if it's public or shared
	if assessment.ComplianceRegister.IsPublic {
		// For public assessments, get all active departments
		return caads.repo.GetAllForPublicAssessment(assessmentID, nil, primaryOrgUnitField, params)
	} else if assessment.ComplianceRegister.IsShared {
		// For shared assessments, try to get shared departments first
		result, err := caads.repo.GetAllForSharedAssessment(assessmentID, uint(assessment.ComplianceRegister.ID), nil, primaryOrgUnitField, params)
		if err != nil {
			return nil, err
		}
		return result, nil
	} else {
		// For non-public/non-shared assessments, use the same active filter as public
		return caads.repo.GetAllForPublicAssessment(assessmentID, nil, primaryOrgUnitField, params)
	}
}

func (caads *ComplianceAssessmentAcknowledgeDepartmentService) Create(assessmentID uint, currentUserID uint, acknowledgeFormDTO *complianceDTOs.ComplianceAssessmentAcknowledgeFormDTO) (uint, error) {
	tx := caads.repo.BeginTransaction()

	// Get assessment details
	assessment, err := caads.assessmentSvc.GetFullDetailsByID(assessmentID, nil)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	// Check if the compliance register has IsRegulationNotificationAssessment = true
	if !assessment.ComplianceRegister.IsRegulationNotificationAssessment {
		tx.Rollback()
		return 0, utility.NewForbiddenError("assessment_is_not_a_regulation_notification_assessment")
	}

	// Get user's primary organization unit ID
	organizationUnitID, err := caads.userSvc.GetUserPrimaryUnitID(currentUserID)
	if err != nil {
		tx.Rollback()
		return 0, utility.NewForbiddenError("organization_unit_not_set_for_user")
	}

	// Check if department already has an acknowledgment
	departmentExists, err := caads.repo.DepartmentExists(assessmentID, organizationUnitID, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	if departmentExists {
		tx.Rollback()
		return 0, utility.NewValidationError(fmt.Errorf("department has already been acknowledged by another user"), "already_acknowledged", "department")
	}

	// Create acknowledge Department
	acknowledgeDepartment := &models.ComplianceAssessmentAcknowledgeDepartment{
		ComplianceAssessmentID: assessmentID,
		UserID:                 currentUserID,
		DepartmentID:           organizationUnitID,
		Comment:                acknowledgeFormDTO.Comment,
		CreatedBYID:            currentUserID,
		CreatedAt:              time.Now(),
	}

	// Create the acknowledge Department
	createdDepartment, err := caads.repo.Create(acknowledgeDepartment, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	// Note: CreateDepartmentAuditLog is not required for compliance as per requirements

	// Commit transaction
	tx.Commit()

	// Add activity log asynchronously
	go func() {
		err := caads.activityLogSvc.LogActivity(
			models.CreateAction,
			models.ComplianceAssessmentAcknowledgeDepartment{}.TableName(),
			createdDepartment.ID,
			createdDepartment,
			currentUserID,
			"", "", "", nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.ComplianceAssessmentAcknowledgeDepartment{}.TableName(),
				"recordID":  createdDepartment.ID,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return createdDepartment.ID, nil
}

func (caads *ComplianceAssessmentAcknowledgeDepartmentService) Delete(assessmentID uint, currentUserID uint) (uint, error) {
	tx := caads.repo.BeginTransaction()
	var deletedUser []map[string]interface{}

	// Get assessment details
	assessment, err := caads.assessmentSvc.GetFullDetailsByID(assessmentID, nil)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	// Check if the compliance register has IsRegulationNotificationAssessment = true
	if !assessment.ComplianceRegister.IsRegulationNotificationAssessment {
		tx.Rollback()
		return 0, utility.NewForbiddenError("assessment_is_not_a_regulation_notification_assessment")
	}

	// Get user's primary organization unit ID
	organizationUnitID, err := caads.userSvc.GetUserPrimaryUnitID(currentUserID)
	if err != nil {
		tx.Rollback()
		return 0, utility.NewForbiddenError("organization_unit_not_set_for_user")
	}

	// Check if department already has an acknowledgment
	departmentExists, err := caads.repo.DepartmentExists(assessmentID, organizationUnitID, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	if !departmentExists {
		tx.Rollback()
		return 0, utility.NewForbiddenError("assessment_has_not_been_acknowledged")
	}

	// Delete the acknowledge Department
	if err := caads.repo.Delete(assessmentID, currentUserID, tx); err != nil {
		tx.Rollback()
		return 0, err
	}

	deletedUser = append(deletedUser, map[string]interface{}{
		"assessment_id": assessmentID,
		"user_id":       currentUserID,
	})

	// Commit transaction
	tx.Commit()

	// Log successful delete activity
	go func() {
		err := caads.activityLogSvc.LogActivity(
			models.DeleteAction,
			models.ComplianceAssessmentAcknowledgeDepartment{}.TableName(),
			assessmentID,
			deletedUser,
			currentUserID,
			"", "", "", nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.DeleteAction,
				"tableName": models.ComplianceAssessmentAcknowledgeDepartment{}.TableName(),
				"recordID":  assessmentID,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log delete activity", fields)
		}
	}()

	return currentUserID, nil
}

func (caads *ComplianceAssessmentAcknowledgeDepartmentService) GetAssessmentAcknowledgmentSummary(assessmentID uint) (*complianceDTOs.ComplianceAssessmentAcknowledgmentSummaryDTO, error) {
	// Get assessment details
	assessment, err := caads.assessmentSvc.GetFullDetailsByID(assessmentID, nil)
	if err != nil {
		return nil, err
	}

	// Check if the compliance register has IsRegulationNotificationAssessment = true
	if !assessment.ComplianceRegister.IsRegulationNotificationAssessment {
		return nil, utility.NewForbiddenError("assessment_is_not_a_regulation_notification_assessment")
	}

	// Get organization level settings
	orgLevel, err := caads.organizationLevelSvc.GetPrimaryLevel()
	if err != nil {
		return nil, err
	}

	primaryOrgUnitField := "organization_unit_1_id"
	if orgLevel.Level == 2 {
		primaryOrgUnitField = "organization_unit_2_id"
	} else if orgLevel.Level == 3 {
		primaryOrgUnitField = "organization_unit_3_id"
	}

	// Check if it's public or shared
	if assessment.ComplianceRegister.IsPublic {
		// For public assessments, get summary for all active departments
		return caads.repo.GetPublicAssessmentAcknowledgmentSummary(assessmentID)
	} else if assessment.ComplianceRegister.IsShared {
		// For shared assessments, get summary only for shared departments
		return caads.repo.GetSharedAssessmentAcknowledgmentSummary(assessmentID, uint(assessment.ComplianceRegister.ID))
	} else {
		// For non-public/non-shared assessments, get summary based on organization unit
		return caads.repo.GetAssessmentAcknowledgmentSummary(assessmentID, primaryOrgUnitField)
	}
}

func (caads *ComplianceAssessmentAcknowledgeDepartmentService) GetAssessmentAcknowledgmentSummaryByOrgUnit(assessmentID uint) ([]complianceDTOs.ComplianceAssessmentAcknowledgmentOrgUnitSummaryDTO, error) {
	// Get assessment details
	assessment, err := caads.assessmentSvc.GetFullDetailsByID(assessmentID, nil)
	if err != nil {
		return nil, err
	}

	// Check if the compliance register has IsRegulationNotificationAssessment = true
	if !assessment.ComplianceRegister.IsRegulationNotificationAssessment {
		return nil, utility.NewForbiddenError("assessment_is_not_a_regulation_notification_assessment")
	}

	// Get organization level settings
	orgLevel, err := caads.organizationLevelSvc.GetPrimaryLevel()
	if err != nil {
		return nil, err
	}

	primaryOrgUnitField := "organization_unit_1_id"
	if orgLevel.Level == 2 {
		primaryOrgUnitField = "organization_unit_2_id"
	} else if orgLevel.Level == 3 {
		primaryOrgUnitField = "organization_unit_3_id"
	}

	// Check if it's public or shared
	if assessment.ComplianceRegister.IsPublic {
		// For public assessments, get summary for all organization units
		return caads.repo.GetPublicAcknowledgmentSummaryByOrgUnit(assessmentID)
	} else {
		// For non-public assessments, get summary based on organization unit
		return caads.repo.GetOrgUnitAcknowledgmentSummary(assessmentID, primaryOrgUnitField)
	}
}
