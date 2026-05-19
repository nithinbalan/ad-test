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
	repo                          complianceRepositories.ComplianceAssessmentAcknowledgeDepartmentRepositoryInterface
	activityLogSvc                logServices.ActivityLogServiceInterface
	assessmentSvc                 ComplianceAssessmentServiceInterface
	userSvc                       organizationServices.UserServiceInterface
	organizationLevelSvc          settingServices.OrganizationLevelServiceInterface
	complianceAssessmentStatusSvc ComplianceAssessmentStatusServiceInterface
}

func NewComplianceAssessmentAcknowledgeDepartmentService(
	repo complianceRepositories.ComplianceAssessmentAcknowledgeDepartmentRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	assessmentSvc ComplianceAssessmentServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	organizationLevelSvc settingServices.OrganizationLevelServiceInterface,
	complianceAssessmentStatusSvc ComplianceAssessmentStatusServiceInterface,
) ComplianceAssessmentAcknowledgeDepartmentServiceInterface {
	return &ComplianceAssessmentAcknowledgeDepartmentService{
		repo:                          repo,
		activityLogSvc:                activityLogSvc,
		assessmentSvc:                 assessmentSvc,
		userSvc:                       userSvc,
		organizationLevelSvc:          organizationLevelSvc,
		complianceAssessmentStatusSvc: complianceAssessmentStatusSvc,
	}
}

func (caads *ComplianceAssessmentAcknowledgeDepartmentService) GetAll(assessmentID uint, currentUserID uint, params *complianceDTOs.ComplianceAssessmentAcknowledgeListQueryParamsDTO) (*complianceDTOs.ComplianceAssessmentAcknowledgementListResponseDTO, error) {
	if err := caads.ensureRegulationAcknowledgementApplicable(assessmentID); err != nil {
		return nil, err
	}

	primaryOrgUnitField, err := caads.getPrimaryOrgUnitField()
	if err != nil {
		return nil, err
	}

	// Keep existing visibility behavior by returning the full eligible department set.
	_ = currentUserID
	return caads.repo.GetAllForRegulationAssessment(assessmentID, nil, primaryOrgUnitField, params)
}

func (caads *ComplianceAssessmentAcknowledgeDepartmentService) Create(assessmentID uint, currentUserID uint, acknowledgeFormDTO *complianceDTOs.ComplianceAssessmentAcknowledgeFormDTO) (uint, error) {
	tx := caads.repo.BeginTransaction()

	if err := caads.ensureRegulationAcknowledgementApplicable(assessmentID, tx); err != nil {
		tx.Rollback()
		return 0, err
	}

	// Get user's primary organization unit ID
	organizationUnitID, err := caads.userSvc.GetUserPrimaryUnitID(currentUserID)
	if err != nil {
		tx.Rollback()
		return 0, utility.NewForbiddenError("organization_unit_not_set_for_user")
	}

	isEligibleDepartment, err := caads.repo.IsDepartmentEligibleForRegulationAcknowledgement(assessmentID, organizationUnitID, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	if !isEligibleDepartment {
		tx.Rollback()
		return 0, utility.NewForbiddenError("user_not_part_of_regulation_requirement_responsible_unit")
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

	// Check if all departments have acknowledged within the transaction
	allAcknowledged, err := caads.repo.CheckAllEligibleDepartmentsAcknowledged(assessmentID, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	// If all departments have acknowledged, update assessment status to "Approved"
	if allAcknowledged {
		// Get approved status using GetByType following existing pattern
		approvedStatus, err := caads.complianceAssessmentStatusSvc.GetByType("approved")
		if err != nil {
			tx.Rollback()
			return 0, err
		}

		// Get the assessment model to update
		complianceAssessment, err := caads.assessmentSvc.GetByID(assessmentID)
		if err != nil {
			tx.Rollback()
			return 0, err
		}

		// Update status to Approved within the same transaction
		complianceAssessment.ComplianceAssessmentStatusID = approvedStatus.ID
		err = caads.repo.ComplianceAssessmentUpdateStatus(complianceAssessment, tx)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
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

	if err := caads.ensureRegulationAcknowledgementApplicable(assessmentID, tx); err != nil {
		tx.Rollback()
		return 0, err
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
	if err := caads.ensureRegulationAcknowledgementApplicable(assessmentID); err != nil {
		return nil, err
	}

	return caads.repo.GetRegulationAssessmentAcknowledgmentSummary(assessmentID)
}

func (caads *ComplianceAssessmentAcknowledgeDepartmentService) GetAssessmentAcknowledgmentSummaryByOrgUnit(assessmentID uint) ([]complianceDTOs.ComplianceAssessmentAcknowledgmentOrgUnitSummaryDTO, error) {
	if err := caads.ensureRegulationAcknowledgementApplicable(assessmentID); err != nil {
		return nil, err
	}

	primaryOrgUnitField, err := caads.getPrimaryOrgUnitField()
	if err != nil {
		return nil, err
	}

	return caads.repo.GetRegulationOrgUnitAcknowledgmentSummary(assessmentID, primaryOrgUnitField)
}

func (caads *ComplianceAssessmentAcknowledgeDepartmentService) ensureRegulationAcknowledgementApplicable(assessmentID uint, tx ...utility.DatabaseTransaction) error {
	if _, err := caads.assessmentSvc.GetByID(assessmentID); err != nil {
		return err
	}

	hasRegulationRequirements, err := caads.repo.HasRegulationNotificationRequirements(assessmentID, tx...)
	if err != nil {
		return err
	}
	if !hasRegulationRequirements {
		return utility.NewForbiddenError("assessment_is_not_a_regulation_notification_assessment")
	}

	return nil
}

func (caads *ComplianceAssessmentAcknowledgeDepartmentService) getPrimaryOrgUnitField() (string, error) {
	orgLevel, err := caads.organizationLevelSvc.GetPrimaryLevel()
	if err != nil {
		return "", err
	}

	switch orgLevel.Level {
	case 1:
		return "organization_unit_1_id", nil
	case 2:
		return "organization_unit_2_id", nil
	case 3:
		return "organization_unit_3_id", nil
	case 4:
		return "organization_unit_4_id", nil
	case 5:
		return "organization_unit_5_id", nil
	default:
		return "", fmt.Errorf("unsupported primary level")
	}
}
