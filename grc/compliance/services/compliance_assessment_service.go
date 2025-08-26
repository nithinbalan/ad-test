package services

import (
	"errors"
	"fmt"
	"math"
	"time"

	commonDTOs "jethur/grc/common/dtos"
	commonServices "jethur/grc/common/services"
	complianceDTOs "jethur/grc/compliance/dtos"
	complianceRepositories "jethur/grc/compliance/repositories"
	logServices "jethur/grc/logs/services"
	masterServices "jethur/grc/masters/services"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type ComplianceAssessmentServiceInterface interface {
	GetAll(params *complianceDTOs.ComplianceAssessmentListQueryParamsDTO) (*complianceDTOs.ComplianceAssessmentListResponseDTO, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*complianceDTOs.ComplianceAssessmentDetailDTO, error)
	GetByID(id uint) (*models.ComplianceAssessment, error)
	Create(currentUserID uint, complianceAssessmentFormDTO *complianceDTOs.ComplianceAssessmentFormDTO) error
	Update(id uint, currentUserID uint, compliance *complianceDTOs.ComplianceAssessmentUpdateFormDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)

	AssessmentCompletionRollup(assessmentID uint) (*complianceDTOs.AssessmentCompletionRollupDTO, error)
	GetComplianceStatusSummaryByRequirement(assessmentID uint) ([]complianceDTOs.ComplianceStatusCountDTO, error)
	GetComplianceScoreSummary(assessmentID uint) (*complianceDTOs.ComplianceScoreSummaryDTO, error)

	SubmitComplianceAssessment(assessmentID uint, currentUserID uint) error
	ApproveComplianceAssessment(assessmentID uint, currentUserID uint) error

	GetApprovalStatus(id uint, tx ...utility.DatabaseTransaction) (*complianceDTOs.ComplianceAssessmentApprovalStatusDTO, error)
}

// Service and implementations
type ComplianceAssessmentService struct {
	repo                          complianceRepositories.ComplianceAssessmentRepositoryInterface
	activityLogSvc                logServices.ActivityLogServiceInterface
	userSvc                       organizationServices.UserServiceInterface
	complianceRegisterSvc         ComplianceRegisterServiceInterface
	applicablityStatusSvc         masterServices.ApplicabilityStatusServiceInterface
	complianceAssessmentStatusSvc ComplianceAssessmentStatusServiceInterface
	notificationSvc               commonServices.NotificationServiceInterface
	complianceClauseSvc           ComplianceClauseServiceInterface
	complianceStatusSvc           ComplianceStatusServiceInterface
}

func NewComplianceAssessmentService(
	repo complianceRepositories.ComplianceAssessmentRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	complianceRegisterSvc ComplianceRegisterServiceInterface,
	applicablityStatusSvc masterServices.ApplicabilityStatusServiceInterface,
	complianceAssessmentStatusSvc ComplianceAssessmentStatusServiceInterface,
	notificationSvc commonServices.NotificationServiceInterface,
	complianceClauseSvc ComplianceClauseServiceInterface,
	complianceStatusSvc ComplianceStatusServiceInterface,
) ComplianceAssessmentServiceInterface {
	return &ComplianceAssessmentService{
		repo:                          repo,
		activityLogSvc:                activityLogSvc,
		userSvc:                       userSvc,
		complianceRegisterSvc:         complianceRegisterSvc,
		applicablityStatusSvc:         applicablityStatusSvc,
		complianceAssessmentStatusSvc: complianceAssessmentStatusSvc,
		notificationSvc:               notificationSvc,
		complianceClauseSvc:           complianceClauseSvc,
		complianceStatusSvc:           complianceStatusSvc,
	}
}

func (cras *ComplianceAssessmentService) GetAll(params *complianceDTOs.ComplianceAssessmentListQueryParamsDTO) (*complianceDTOs.ComplianceAssessmentListResponseDTO, error) {
	return cras.repo.GetAll(params)
}

func (cras ComplianceAssessmentService) GetByID(id uint) (*models.ComplianceAssessment, error) {
	assessment, err := cras.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return assessment, nil
}

func (cras ComplianceAssessmentService) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*complianceDTOs.ComplianceAssessmentDetailDTO, error) {
	var dto complianceDTOs.ComplianceAssessmentDetailDTO

	// Fetch full model with relations
	model, err := cras.repo.GetFullDetailsByID(id, tx...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	// Map model to DTO
	if err := mapstructure.Decode(model, &dto); err != nil {
		return nil, err
	}

	allCompleted := true
	allApproved := true

	// Loop over requirements in the model directly
	for _, reqModel := range model.AssessmentRequirements {
		reqCompleted := true
		reqApproved := true

		if len(reqModel.ResponsibleUnits) == 0 {
			reqCompleted = false
			reqApproved = false
		} else {
			for _, ru := range reqModel.ResponsibleUnits {
				if !ru.IsCompleted {
					reqCompleted = false
				}
				if ru.ComplianceAssessmentRequirementStatus == nil || ru.ComplianceAssessmentRequirementStatus.Type != "approved" {
					reqApproved = false
				}
			}
		}

		if !reqCompleted {
			allCompleted = false
		}
		if !reqApproved {
			allApproved = false
		}
	}

	// Set assessment-level flags
	dto.IsCompleted = allCompleted
	dto.IsApproved = allApproved

	dto.CreatedAt = model.CreatedAt
	dto.UpdatedAt = model.UpdatedAt

	return &dto, nil
}

func (cras *ComplianceAssessmentService) Create(currentUserID uint, complianceAssessmentFormDTO *complianceDTOs.ComplianceAssessmentFormDTO) error {
	tx := cras.repo.BeginTransaction()
	var createdIds []uint
	var addedComplianceAssessments []map[string]interface{}

	if err := cras.validateComplianceAssessmentFormDTO(complianceAssessmentFormDTO); err != nil {
		return err
	}

	compliance, err := cras.complianceRegisterSvc.GetFullDetailsByID(complianceAssessmentFormDTO.ComplianceRegisterID)
	if err != nil {
		return err
	}

	complianceAssessmentStatus, err := cras.complianceAssessmentStatusSvc.GetByType("draft")
	if err != nil {
		return err
	}

	complianceAssessment := &models.ComplianceAssessment{
		ComplianceRegisterID:         complianceAssessmentFormDTO.ComplianceRegisterID,
		ReviewUserID:                 complianceAssessmentFormDTO.ReviewUserID,
		Description:                  complianceAssessmentFormDTO.Description,
		CreatedBYID:                  currentUserID,
		ComplianceAssessmentStatusID: complianceAssessmentStatus.ID,
	}
	id, err := cras.repo.Create(complianceAssessment, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	for _, clause := range compliance.ComplianceClauses {
		clauseDetails, err := cras.complianceClauseSvc.GetFullDetailsByID(clause.ID)
		if err != nil {
			tx.Rollback()
			return err
		}
		for _, complianceRequirement := range clauseDetails.ComplianceRequirements {
			requirement := &models.ComplianceAssessmentRequirement{
				ComplianceAssessmentID:  id,
				ComplianceRequirementID: uint(complianceRequirement.ID),
				ComplianceStatusID:      nil,
				CreatedBYID:             currentUserID,
			}
			requirementId, err := cras.repo.CreateComplianceAssessmentRequirement(requirement, tx)
			if err != nil {
				tx.Rollback()
				return err
			}

			var deadline *time.Time
			if complianceRequirement.DeadlineString != nil {
				parsedDeadline, err := time.Parse("2006-01-02", *complianceRequirement.DeadlineString)
				if err != nil {
					tx.Rollback()
					return err
				}
				deadline = &parsedDeadline
			}

			applicabilityStatus, err := cras.applicablityStatusSvc.GetByType("applicable")
			if err != nil {
				return err
			}

			for _, ru := range complianceRequirement.ResponsibleUnits {
				// Step 1: Create ComplianceAssessmentRequirementResponsibleUnit one-by-one
				unit := &models.ComplianceAssessmentRequirementResponsibleUnit{
					ComplianceAssessmentRequirementID: requirementId,
					ResponsibleUnitID:                 uint(ru.ResponsibleUnit.ID),
					Deadline:                          deadline,
					CreatedBYID:                       currentUserID,
				}
				if err := cras.repo.CreateComplianceAssessmentRequirementResponsibleUnit(unit, tx); err != nil {
					tx.Rollback()
					return err
				}

				// Step 2: Create associated DocumentTitles using the created unit's ID
				var requirementDocumentTitles []*models.ComplianceAssessmentRequirementDocumentTitle

				for _, documentTitle := range complianceRequirement.RequiredDocumentTitles {
					doc := &models.ComplianceAssessmentRequirementDocumentTitle{
						ComplianceAssessmentRequirementResponsibleUnitID: unit.ID,
						RequiredDocumentTitleID:                          &documentTitle.RequiredDocumentTitleID,
						ApplicabilityStatusID:                            applicabilityStatus.ID,
						CreatedBYID:                                      currentUserID,
					}
					requirementDocumentTitles = append(requirementDocumentTitles, doc)
				}

				// Add an empty document title (with nil RequiredDocumentTitleID)
				emptyDoc := &models.ComplianceAssessmentRequirementDocumentTitle{
					ComplianceAssessmentRequirementResponsibleUnitID: unit.ID,
					RequiredDocumentTitleID:                          nil,
					ApplicabilityStatusID:                            applicabilityStatus.ID,
					CreatedBYID:                                      currentUserID,
				}
				requirementDocumentTitles = append(requirementDocumentTitles, emptyDoc)

				// Bulk insert documents for this unit
				if len(requirementDocumentTitles) > 0 {
					if err := cras.repo.CreateComplianceAssessmentRequirementDocumentTitle(requirementDocumentTitles, tx); err != nil {
						tx.Rollback()
						return err
					}
				}
			}

		}
	}

	addedComplianceAssessments = append(addedComplianceAssessments, map[string]interface{}{
		"id":                    id,
		"compliance_assessment": complianceAssessment,
	})

	createdIds = append(createdIds, id)

	tx.Commit()

	go func() {
		if len(createdIds) > 0 {
			err := cras.activityLogSvc.LogBulkActivities(
				models.CreateAction,
				models.ComplianceAssessment{}.TableName(),
				createdIds,
				addedComplianceAssessments,
				currentUserID,
				nil,
			)
			if err != nil {
				fields := map[string]interface{}{
					"action":    models.CreateAction,
					"tableName": models.ComplianceAssessment{}.TableName(),
					"recordID":  createdIds,
					"userID":    currentUserID,
					"error":     err.Error(),
				}

				utility.LogError("Failed to log activity", fields)
			}
		}
	}()

	return nil
}

func (cras *ComplianceAssessmentService) Update(id uint, currentUserID uint, complianceAssessmentUpdateFormDTO *complianceDTOs.ComplianceAssessmentUpdateFormDTO) error {
	tx := cras.repo.BeginTransaction()

	existingComplianceAssessment, err := cras.GetByID(id)
	if err != nil {
		return err
	}

	oldComplianceAssessment, err := cras.GetFullDetailsByID(id, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	existingComplianceAssessment.ComplianceRegisterID = complianceAssessmentUpdateFormDTO.ComplianceRegisterID
	existingComplianceAssessment.ReviewUserID = complianceAssessmentUpdateFormDTO.ReviewUserID
	// existingComplianceAssessment.PrimaryUnitID = complianceAssessmentUpdateFormDTO.OrganizationUnitIDs
	existingComplianceAssessment.Description = complianceAssessmentUpdateFormDTO.Description
	existingComplianceAssessment.UpdatedBYID = &currentUserID

	err = cras.repo.Update(existingComplianceAssessment, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	go func() {
		existingAssessment, err := cras.GetFullDetailsByID(id)
		if err != nil {
			return
		}
		// Prepare log data
		changedData := make(map[string]interface{})
		// log changes
		oldData := map[string]interface{}{
			"compliance_register": oldComplianceAssessment.ComplianceRegister.TitleEn,
			// "organization_unit":   oldComplianceAssessment.OrganizationUnit.TitleEn,
			"description": oldComplianceAssessment.Description,
			"review_user": fmt.Sprintf("%s %s", oldComplianceAssessment.ReviewUser.FirstNameEN, oldComplianceAssessment.ReviewUser.LastNameEN),
		}

		newData := map[string]interface{}{
			"compliance_register": existingAssessment.ComplianceRegister.TitleEn,
			// "organization_unit":   existingAssessment.OrganizationUnit.TitleEn,
			"description": existingAssessment.Description,
			"review_user": fmt.Sprintf("%s %s", existingAssessment.ReviewUser.FirstNameEN, existingAssessment.ReviewUser.LastNameEN),
		}
		// Compare other fields and add them to changedData
		fieldChanges := utility.CompareFieldChanges(oldData, newData)
		for k, v := range fieldChanges {
			changedData[k] = v
		}
		// Log activity if there are any changes
		if len(changedData) > 0 {
			// Log the update activity
			err := cras.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.ComplianceAssessment{}.TableName(),
				id,
				changedData,
				currentUserID,
				nil,
			)
			if err != nil {
				utility.LogError("Failed to log activity", map[string]interface{}{
					"action":    models.UpdateAction,
					"tableName": models.ComplianceAssessment{}.TableName(),
					"recordID":  id,
					"userID":    currentUserID,
					"error":     err.Error(),
				})
			}
		}
	}()

	return nil
}

func (cras *ComplianceAssessmentService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	tx := cras.repo.BeginTransaction()

	var failedIds []uint
	var successfulIds []uint
	var deletedComplianceAssessment []map[string]interface{}

	for _, id := range ids {
		complianceAssessment, err := cras.GetFullDetailsByID(id, tx)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		if len(complianceAssessment.AssessmentRequirements) > 0 {
			var deleteRequirementIds []*uint
			var deletedResponsibleUnitIds []*uint

			for _, assmntRequirement := range complianceAssessment.AssessmentRequirements {
				deleteRequirementIds = append(deleteRequirementIds, assmntRequirement.ID)

				// Delete all document titles under each Responsible Unit
				for _, ru := range assmntRequirement.ResponsibleUnits {
					// Delete document titles for this responsible unit
					if len(ru.RequiredDocumentTitles) > 0 {
						var docTitleIDs []*uint
						for _, doc := range ru.RequiredDocumentTitles {
							docTitleIDs = append(docTitleIDs, doc.ID)
						}
						if err := cras.repo.DeleteComplianceAssessmentRequirementDocumentTitles(docTitleIDs); err != nil {
							failedIds = append(failedIds, id)
							continue
						}
					}
					deletedResponsibleUnitIds = append(deletedResponsibleUnitIds, &ru.ID)
				}
			}
			if len(deletedResponsibleUnitIds) > 0 {
				if err := cras.repo.DeleteComplianceAssessmentRequirementResponsibleUnits(deletedResponsibleUnitIds); err != nil {
					failedIds = append(failedIds, id)
					continue
				}
			}
			// Delete Assessment Requirements
			if len(deleteRequirementIds) > 0 {
				if err := cras.repo.DeleteComplianceAssessmentRequirements(deleteRequirementIds); err != nil {
					failedIds = append(failedIds, id)
					continue
				}
			}
		}

		// Delete the Compliance Assessment itself
		if err := cras.repo.Delete(id, tx); err != nil {
			failedIds = append(failedIds, id)
			continue
		} else {
			successfulIds = append(successfulIds, id)
		}

		deletedComplianceAssessment = append(deletedComplianceAssessment, map[string]interface{}{
			"id":       id,
			"title_en": complianceAssessment.ComplianceRegister.TitleEn,
		})
	}

	if len(failedIds) > 0 {
		tx.Rollback()
		return failedIds, nil
	}

	tx.Commit()

	if len(successfulIds) > 0 {
		go func() {
			err := cras.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.ComplianceAssessment{}.TableName(),
				successfulIds,
				deletedComplianceAssessment,
				currentUserID,
				nil,
			)
			if err != nil {
				fmt.Printf("Failed to log bulk delete activity: %v\n", err)
			}
		}()
	}

	return failedIds, nil
}

func (cras *ComplianceAssessmentService) validateComplianceAssessmentFormDTO(complianceAssessmentFormDTO *complianceDTOs.ComplianceAssessmentFormDTO) error {
	if _, err := cras.complianceRegisterSvc.GetByID(complianceAssessmentFormDTO.ComplianceRegisterID); err != nil {
		return utility.NewValidationError(fmt.Errorf("compliance register does not exist"), "not_exist", "compliance register")
	}

	if _, err := cras.userSvc.GetByID(complianceAssessmentFormDTO.ReviewUserID); err != nil {
		return utility.NewValidationError(fmt.Errorf("review user does not exist"), "not_exist", "review user")
	}

	statusIds := []uint{}

	complianceAssessmentDraftStatus, err := cras.complianceAssessmentStatusSvc.GetByType("draft")
	if err != nil {
		return err
	}

	statusIds = append(statusIds, uint(complianceAssessmentDraftStatus.ID))

	complianceAssessmentInreviewStatus, err := cras.complianceAssessmentStatusSvc.GetByType("in_review")
	if err != nil {
		return err
	}

	statusIds = append(statusIds, uint(complianceAssessmentInreviewStatus.ID))

	complianceAssessmentInProgressStatus, err := cras.complianceAssessmentStatusSvc.GetByType("in_progress")
	if err != nil {
		return err
	}

	statusIds = append(statusIds, uint(complianceAssessmentInProgressStatus.ID))

	exists, err := cras.repo.IsAssessmentExists(0, complianceAssessmentFormDTO.ComplianceRegisterID, statusIds)
	if err != nil {
		return fmt.Errorf("failed to validate assessment %w", err)
	}

	if exists {
		return utility.NewForbiddenError("assessment_already_in_progress_for_compliance_register")
	}
	hasRequirement, err := cras.repo.HasRequirementForRegister(complianceAssessmentFormDTO.ComplianceRegisterID)
	if err != nil {
		return fmt.Errorf("failed to validate requirements: %w", err)
	}

	if !hasRequirement {
		return utility.NewValidationError(
			fmt.Errorf("no compliance requirements found for this compliance register"),
			"reqirement_not_exist",
			"compliance requirement",
		)
	}

	return nil
}

func (cras *ComplianceAssessmentService) validateComplianceAssessmentOrganizationUnit(organizationUnitID uint, complianceRegisterID uint, id *uint) error {
	statusIds := []uint{}

	complianceAssessmentDraftStatus, err := cras.complianceAssessmentStatusSvc.GetByType("draft")
	if err != nil {
		return err
	}

	statusIds = append(statusIds, uint(complianceAssessmentDraftStatus.ID))

	complianceAssessmentInreviewStatus, err := cras.complianceAssessmentStatusSvc.GetByType("in_review")
	if err != nil {
		return err
	}

	statusIds = append(statusIds, uint(complianceAssessmentInreviewStatus.ID))
	if id == nil {
		exists, err := cras.repo.IsAssessmentExists(organizationUnitID, complianceRegisterID, statusIds)
		if err != nil {
			return fmt.Errorf("failed to validate organization unit ID %d: %w", organizationUnitID, err)
		}

		if exists {
			return utility.NewValidationError(fmt.Errorf("a assessment already exists for this organization unit"), "unique", "organization unit")
		}
	} else {
		exists, err := cras.repo.IsAssessmentExistExcludingID(*id, organizationUnitID, complianceRegisterID, statusIds)
		if err != nil {
			return fmt.Errorf("failed to validate organization unit ID %d: %w", organizationUnitID, err)
		}

		if exists {
			return utility.NewValidationError(fmt.Errorf("a assessment already exists for this organization unit"), "unique", "organization unit")
		}
	}

	return nil
}

func (cras *ComplianceAssessmentService) AssessmentCompletionRollup(assessmentID uint) (*complianceDTOs.AssessmentCompletionRollupDTO, error) {
	if _, err := cras.GetByID(assessmentID); err != nil {
		return nil, err
	}
	return cras.repo.AssessmentCompletionRollup(assessmentID)
}

func (cras *ComplianceAssessmentService) GetComplianceStatusSummaryByRequirement(assessmentID uint) ([]complianceDTOs.ComplianceStatusCountDTO, error) {
	if _, err := cras.GetByID(assessmentID); err != nil {
		return nil, err
	}
	return cras.repo.GetComplianceStatusSummaryByRequirement(assessmentID)
}

func (cras *ComplianceAssessmentService) GetComplianceScoreSummary(assessmentID uint) (*complianceDTOs.ComplianceScoreSummaryDTO, error) {
	// Optional: ensure the assessment exists (keeps behavior consistent with your other methods)
	if _, err := cras.GetByID(assessmentID); err != nil {
		return nil, err
	}

	return cras.repo.GetComplianceScoreSummary(assessmentID)
}

func (cras *ComplianceAssessmentService) SubmitComplianceAssessment(assessmentID uint, currentUserID uint) error {
	complianceAssessmentInProgressStatus, err := cras.complianceAssessmentStatusSvc.GetByType("in_progress")
	if err != nil {
		return err
	}

	assessment, err := cras.GetByID(assessmentID)
	if err != nil {
		return err
	}

	assessmentDetails, err := cras.GetFullDetailsByID(assessmentID)
	if err != nil {
		return err
	}

	// Validate if the assessment is already submitted
	if assessment.ComplianceAssessmentStatusID != complianceAssessmentInProgressStatus.ID {
		return utility.NewForbiddenError("Assessment not in progress")
	}

	complianceAssessmentInreviewStatus, err := cras.complianceAssessmentStatusSvc.GetByType("in_review")
	if err != nil {
		return err
	}

	submittedAt := time.Now()
	assessment.ComplianceAssessmentStatusID = complianceAssessmentInreviewStatus.ID
	assessment.SubmittedBy = &currentUserID
	assessment.SubmittedAt = &submittedAt

	err = cras.repo.Update(assessment)
	if err != nil {
		return err
	}

	go func() {
		err := cras.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.ComplianceAssessment{}.TableName(),
			assessmentID,
			map[string]interface{}{
				"status":       "submitted",
				"submitted_by": currentUserID,
				"submitted_at": submittedAt,
			},
			currentUserID,
			nil,
		)
		if err != nil {
			utility.LogError("Failed to log activity", map[string]interface{}{
				"action":    models.UpdateAction,
				"tableName": models.ComplianceAssessment{}.TableName(),
				"recordID":  assessmentID,
				"userID":    currentUserID,
				"error":     err.Error(),
			})
			fmt.Println("Failed to log activity:", err)
			return
		}

		// Fetch Review User
		reviewUser, err := cras.userSvc.GetByID(assessment.ReviewUserID)
		if err != nil {
			fmt.Println("Error: Failed to fetch review user", err)
			return
		}
		if reviewUser == nil {
			fmt.Println("Error: reviewUser is nil")
			return
		}

		createdUser, err := cras.userSvc.GetByID(currentUserID)
		if err != nil {
			fmt.Println("Error: Failed to fetch created user", err)
			return
		}
		if createdUser == nil {
			fmt.Println("Error: createdUser is nil")
			return
		}

		// Prepare DB Notification
		urlParam := fmt.Sprintf(`{"id":%d}`, assessmentID)
		entity := "compliance-assessment"
		action := "submitted"

		dbNotificationDTO := &commonDTOs.DBNotificationDTO{
			UserIDs:              []uint{assessment.ReviewUserID},
			ModuleIdentifier:     "compliance",
			NotificationType:     string(commonDTOs.AssignmentType),
			HeadingEn:            "New Compliance Assessment Submitted",
			HeadingAr:            "تم تقديم تقييم الامتثال الجديد",
			TitleEn:              assessmentDetails.ComplianceRegister.TitleEn,
			TitleAr:              assessmentDetails.ComplianceRegister.TitleAr,
			ReferenceNumber:      assessmentDetails.ComplianceRegister.ReferenceNumber,
			Entity:               &entity,
			Action:               &action,
			URLParam:             urlParam,
			IsSystemNotification: false,
			CreatedByID:          &currentUserID,
		}

		// Prepare Email Notification
		emailTo := []string{reviewUser.Email}
		subject := fmt.Sprintf("Compliance Assessment/%s: Review Assigned", assessmentDetails.ComplianceRegister.TitleEn)
		appURL := utility.GetAppURL()
		emailData := map[string]interface{}{
			"ReviewerFirstName": reviewUser.FirstNameEn,
			"ReviewerLastName":  reviewUser.LastNameEn,
			// "PrimaryUnit":       assessmentDetails.OrganizationUnit.TitleEn,
			"ReferenceNumber": assessmentDetails.ComplianceRegister.ReferenceNumber,
			"TitleEn":         assessmentDetails.ComplianceRegister.TitleEn,
			"TitleAr":         assessmentDetails.ComplianceRegister.TitleAr,
			"Description":     assessment.Description,
			"CreatedBy":       fmt.Sprintf("%s %s", createdUser.FirstNameEn, createdUser.LastNameEn),
			"SubmittedAt":     submittedAt.Format("2006-01-02"),
			"RedirectURL":     fmt.Sprintf("%s/compliance/assessments/%d/overview", appURL, assessmentID),
		}

		emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
			EmailTo:              emailTo,
			Subject:              subject,
			TemplateRelativePath: "compliance/templates/email/compliance_assessment_review_mail.html",
			EmailData:            emailData,
		}

		// Send Notification
		err = cras.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "compliance_assessment")
		if err != nil {
			fmt.Println("Failed to send notifications:", err)
		}
	}()

	return nil
}

func (cras *ComplianceAssessmentService) ApproveComplianceAssessment(assessmentID uint, currentUserID uint) error {
	complianceAssessmentInReviewStatus, err := cras.complianceAssessmentStatusSvc.GetByType("in_review")
	if err != nil {
		return err
	}

	assessment, err := cras.GetByID(assessmentID)
	if err != nil {
		return err
	}

	// Validate if the assessment is in draft
	if assessment.ComplianceAssessmentStatusID != complianceAssessmentInReviewStatus.ID {
		return utility.NewForbiddenError("Assessment is still in draft")
	}

	if currentUserID != assessment.ReviewUserID {
		return utility.NewForbiddenError("User is not authorized to submit this assessment")
	}

	complianceAssessmentApprovedStatus, err := cras.complianceAssessmentStatusSvc.GetByType("approved")
	if err != nil {
		return err
	}
	tx := cras.repo.BeginTransaction()

	assessment.ComplianceAssessmentStatusID = complianceAssessmentApprovedStatus.ID
	assessment.UpdatedBYID = &currentUserID
	err = cras.repo.Update(assessment, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	scores, err := cras.CalculateComplianceScores(assessmentID, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	reqEval := make(map[uint]complianceDTOs.ComplianceRequirementScoreStatusDTO)
	for _, s := range scores {
		if s == nil || s.ComplianceAssessmentRequirementID == 0 || s.MappedComplianceStatusID == nil {
			continue
		}
		pct := s.Score
		pctCopy := pct
		reqEval[s.ComplianceAssessmentRequirementID] = complianceDTOs.ComplianceRequirementScoreStatusDTO{
			ComplianceStatusID: s.MappedComplianceStatusID,
			Score:              &pctCopy,
		}
	}

	if err := cras.UpdateComplianceRequirementValueFromScores(assessmentID, reqEval, tx); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	go func() {
		err := cras.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.ComplianceAssessment{}.TableName(),
			assessmentID,
			map[string]interface{}{
				"status":      "approved",
				"approved_by": currentUserID,
			},
			currentUserID,
			nil,
		)
		if err != nil {
			utility.LogError("Failed to log activity", map[string]interface{}{
				"action":    models.UpdateAction,
				"tableName": models.ComplianceAssessment{}.TableName(),
				"recordID":  assessmentID,
				"userID":    currentUserID,
				"error":     err.Error(),
			})
		}

		assessmentDetails, err := cras.GetFullDetailsByID(assessmentID)
		if err != nil {
			return
		}

		submittedUser, err := cras.userSvc.GetByID(*assessment.SubmittedBy)
		if err != nil {
			fmt.Println("Error: Failed to fetch submitted user", err)
			return
		}

		createdUser, err := cras.userSvc.GetByID(currentUserID)
		if err != nil {
			fmt.Println("Error: Failed to fetch created user", err)
			return
		}

		// Prepare DB Notification
		urlParam := fmt.Sprintf(`{"id":%d}`, assessmentID)
		entity := "compliance-assessment"
		action := "approved"

		dbNotificationDTO := &commonDTOs.DBNotificationDTO{
			UserIDs:              []uint{submittedUser.ID},
			ModuleIdentifier:     "compliance",
			NotificationType:     string(commonDTOs.InformationType),
			HeadingEn:            "New Compliance Assessment Submitted",
			HeadingAr:            "تم تقديم تقييم الامتثال الجديد",
			TitleEn:              assessmentDetails.ComplianceRegister.TitleEn,
			TitleAr:              assessmentDetails.ComplianceRegister.TitleAr,
			ReferenceNumber:      assessmentDetails.ComplianceRegister.ReferenceNumber,
			Entity:               &entity,
			Action:               &action,
			URLParam:             urlParam,
			IsSystemNotification: false,
			CreatedByID:          &currentUserID,
		}

		// Prepare Email Notification
		emailTo := []string{submittedUser.Email}
		subject := fmt.Sprintf("Compliance Assessment/%s: Approved", assessmentDetails.ComplianceRegister.TitleEn)
		appURL := utility.GetAppURL()
		emailData := map[string]interface{}{
			"FirstName": submittedUser.FirstNameEn,
			"LastName":  submittedUser.LastNameEn,
			// "PrimaryUnit":     assessmentDetails.OrganizationUnit.TitleEn,
			"ReferenceNumber": assessmentDetails.ComplianceRegister.ReferenceNumber,
			"TitleEn":         assessmentDetails.ComplianceRegister.TitleEn,
			"TitleAr":         assessmentDetails.ComplianceRegister.TitleAr,
			"Description":     assessment.Description,
			"CreatedBy":       fmt.Sprintf("%s %s", createdUser.FirstNameEn, createdUser.LastNameEn),
			"ApprovedAt":      time.Now().Format("2006-01-02"),
			"RedirectURL":     fmt.Sprintf("%s/compliance/assessments/%d/overview", appURL, assessmentID),
		}

		emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
			EmailTo:              emailTo,
			Subject:              subject,
			TemplateRelativePath: "compliance/templates/email/compliance_assessment_approve_mail.html",
			EmailData:            emailData,
		}

		// Send Notification
		err = cras.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "compliance_assessment")
		if err != nil {
			fmt.Println("Failed to send notifications:", err)
		}
	}()

	return nil
}

func (cras *ComplianceAssessmentService) UpdateComplianceRequirementValueFromScores(assessmentID uint, reqEval map[uint]complianceDTOs.ComplianceRequirementScoreStatusDTO, tx ...utility.DatabaseTransaction) error {
	// We still need RU rows to write values -> fetch with tx
	assessment, err := cras.repo.GetFullDetailsByID(assessmentID, tx...)
	if err != nil {
		return err
	}

	for _, ar := range assessment.AssessmentRequirements {
		eval, ok := reqEval[ar.ID]
		if !ok {
			continue
		}

		for _, ru := range ar.ResponsibleUnits {
			if ru.ComplianceAssessmentFrameworkValueID == nil {
				// nothing selected at RU level; skip or set a default if business rule says so
				continue
			}

			dto := complianceDTOs.ComplianceRequirementValueUpdateDTO{
				ComplianceRequirementID:              ar.ComplianceRequirementID,
				ComplianceStatusID:                   eval.ComplianceStatusID,
				ComplianceAssessmentFrameworkValueID: *ru.ComplianceAssessmentFrameworkValueID,
				ResponsibleUnitID:                    ru.ResponsibleUnitID,
				Score:                                eval.Score,
				Comments:                             ru.Comments,          // new
				CompletedByUserID:                    ru.CompletedByUserID, // new
				CompletedAt:                          ru.CompletedAt,
			}

			if err := cras.repo.UpdateComplianceRequirementResponsibleUintValue(dto, tx...); err != nil {
				return err
			}
			if err := cras.repo.UpdateComplianceRequirementStatusUpdate(dto, tx...); err != nil {
				return err
			}
		}
	}
	return nil
}

func (cras *ComplianceAssessmentService) UpdateClauseStatusesFromScores(assessmentID uint, reqEval map[uint]complianceDTOs.ComplianceRequirementScoreStatusDTO, tx ...utility.DatabaseTransaction) error {
	// canonical status ids
	compliant, err := cras.complianceStatusSvc.GetByType("compliant")
	if err != nil {
		return err
	}
	partial, err := cras.complianceStatusSvc.GetByType("partially_compliant")
	if err != nil {
		return err
	}
	notImpl, err := cras.complianceStatusSvc.GetByType("not_implemented")
	if err != nil {
		return err
	}

	// need AR -> Requirement -> ClauseID
	ass, err := cras.repo.GetFullDetailsByID(assessmentID, tx...)
	if err != nil {
		return err
	}
	// Ensure repo has: Preload("AssessmentRequirements.ComplianceRequirement")

	// clauseID -> aggregate DTO
	clauseAgg := make(map[uint]*complianceDTOs.ClauseStatusAggregateDTO)

	// Also track sum of scores to compute clause score
	type clauseScoreAgg struct {
		sum   float64
		count int
	}
	clauseScore := make(map[uint]*clauseScoreAgg)

	for _, ar := range ass.AssessmentRequirements {
		if ar.ComplianceRequirement == nil {
			continue // must be preloaded
		}
		clauseID := ar.ComplianceRequirement.ComplianceClauseID

		if clauseAgg[clauseID] == nil {
			clauseAgg[clauseID] = &complianceDTOs.ClauseStatusAggregateDTO{}
		}
		if clauseScore[clauseID] == nil {
			clauseScore[clauseID] = &clauseScoreAgg{}
		}

		a := clauseAgg[clauseID]
		a.TotalCount++

		// Use in-memory computed status/score from reqEval (keyed by AR.ID)
		eval, ok := reqEval[ar.ID]
		if !ok || eval.ComplianceStatusID == nil {
			// treat missing as Not Implemented (and 0 score)
			a.NotImplementedCount++
			// score contribution = 0; still counts in denominator
			clauseScore[clauseID].count++
			continue
		}

		switch *eval.ComplianceStatusID {
		case compliant.ID:
			a.CompliantCount++
		case partial.ID:
			a.PartialCount++
		case notImpl.ID:
			a.NotImplementedCount++
		default:
			a.PartialCount++
		}

		// Score: if provided, use it; else treat as 0
		if eval.Score != nil {
			clauseScore[clauseID].sum += *eval.Score
		}
		clauseScore[clauseID].count++
	}

	// Decide status + compute/round score, then persist
	for clauseID, a := range clauseAgg {
		var setStatusID uint
		switch {
		case a.PartialCount > 0:
			setStatusID = partial.ID
		case a.TotalCount > 0 && a.CompliantCount == a.TotalCount:
			setStatusID = compliant.ID
		case a.TotalCount > 0 && a.NotImplementedCount == a.TotalCount:
			setStatusID = notImpl.ID
		default:
			setStatusID = partial.ID
		}

		// Average score across its requirements (0–100), round to 2 dp
		cs := clauseScore[clauseID]
		avg := 0.0
		if cs != nil && cs.count > 0 {
			avg = cs.sum / float64(cs.count)
			avg = math.Round(avg*100) / 100
		}

		// Update both fields (use one method or two, your choice)
		if err := cras.repo.UpdateClauseComplianceStatus(clauseID, setStatusID, avg, tx...); err != nil {
			return err
		}
	}

	return nil
}

func (cras *ComplianceAssessmentService) CalculateComplianceScores(assessmentID uint, tx ...utility.DatabaseTransaction) ([]*complianceDTOs.ResponsibleUnitComplianceScoreDTO, error) {
	assessment, err := cras.repo.GetFullDetailsByID(assessmentID, tx...)
	if err != nil {
		return nil, err
	}
	// Canonical status ids for clause roll-up
	compliant, err := cras.complianceStatusSvc.GetByType("compliant")
	if err != nil {
		return nil, err
	}
	partial, err := cras.complianceStatusSvc.GetByType("partially_compliant")
	if err != nil {
		return nil, err
	}
	notImpl, err := cras.complianceStatusSvc.GetByType("not_implemented")
	if err != nil {
		return nil, err
	}
	notApplicable, err := cras.complianceStatusSvc.GetByType("not_applicable")
	if err != nil {
		return nil, err
	}
	// fmt.Println("=== CalculateComplianceScores Debug ===")

	frameworkType := ""
	frameworkValuesMap := make(map[uint]int)
	var maxValueInFramework int

	// Framework setup
	if reg := assessment.ComplianceRegister; reg != nil && reg.ComplianceAssessmentFramework != nil {
		if ft := reg.ComplianceAssessmentFramework.AssessmentFrameworkType; ft != nil {
			frameworkType = ft.Type
		}
		for _, v := range reg.ComplianceAssessmentFramework.Values {
			if v != nil && v.ID > 0 && v.Value != nil {
				frameworkValuesMap[v.ID] = *v.Value
				if *v.Value > maxValueInFramework {
					maxValueInFramework = *v.Value
				}
			}
		}
	}

	// fmt.Println("Framework Type:", frameworkType)
	// fmt.Println("Framework Values Map:", frameworkValuesMap)
	// fmt.Println("Max Value In Framework:", maxValueInFramework)

	// per-requirement max (dynamic)
	perReqMax := func() float64 {
		switch frameworkType {
		case "percentage":
			return 100.0
		case "binary":
			return 1.0
		case "point", "maturity":
			if maxValueInFramework > 0 {
				return float64(maxValueInFramework)
			}
			if frameworkType == "point" {
				return 2.0
			}
			return 5.0
		default:
			return 100.0
		}
	}()
	// fmt.Println("Per Requirement Max:", perReqMax)

	// Map to track responsibleUnit scores
	responsibleUnitScores := make(map[uint]*complianceDTOs.ResponsibleUnitComplianceScoreDTO)
	reqClause := make(map[uint]uint)

	// Walk all requirement-responsibleUnit rows and accumulate globally
	for _, req := range assessment.AssessmentRequirements {
		if req.ComplianceRequirement != nil {
			reqClause[req.ID] = req.ComplianceRequirement.ComplianceClauseID
		}
		ComplianceClauseID := req.ComplianceRequirement.ComplianceClauseID
		// fmt.Println("ComplianceClauseID ID----:", ComplianceClauseID)
		applicableFound := false
		notApplicableCount := 0
		for _, ru := range req.ResponsibleUnits {

			complianceAssessmentRequirementID := ru.ComplianceAssessmentRequirementID

			// fmt.Println("complianceAssessmentRequirement ID:", complianceAssessmentRequirementID)
			if responsibleUnitScores[complianceAssessmentRequirementID] == nil {
				responsibleUnitScores[complianceAssessmentRequirementID] = &complianceDTOs.ResponsibleUnitComplianceScoreDTO{
					ComplianceRegisterID:              assessment.ComplianceRegisterID,
					ComplianceAssessmentRequirementID: complianceAssessmentRequirementID,
					FrameworkType:                     frameworkType,
				}
			}
			if responsibleUnitScores[complianceAssessmentRequirementID].ComplianceClauseID == 0 {
				responsibleUnitScores[complianceAssessmentRequirementID].ComplianceClauseID = ComplianceClauseID
			}
			isApplicable := true
			// fmt.Println("isApplicable:", ru.ComplianceAssessmentFrameworkValue.IsApplicable)

			if ru.ComplianceAssessmentFrameworkValueID != nil {
				isApplicable = ru.ComplianceAssessmentFrameworkValue.IsApplicable
			}
			// Only count applicable requirements
			if !isApplicable {
				notApplicableCount++
				continue
			}
			applicableFound = true
			responsibleUnitScores[complianceAssessmentRequirementID].TotalRequirements++
			// totalRequirements++

			switch frameworkType {
			case "point", "maturity":
				if ru.ComplianceAssessmentFrameworkValueID != nil {
					if val, ok := frameworkValuesMap[*ru.ComplianceAssessmentFrameworkValueID]; ok {
						responsibleUnitScores[complianceAssessmentRequirementID].ScoredRequirements++
						responsibleUnitScores[complianceAssessmentRequirementID].TotalScore += float64(val)
						responsibleUnitScores[complianceAssessmentRequirementID].MaxPossibleScore += perReqMax
					}
				}
			case "percentage":
				if ru.Score != nil {
					responsibleUnitScores[complianceAssessmentRequirementID].ScoredRequirements++
					responsibleUnitScores[complianceAssessmentRequirementID].TotalScore += float64(*ru.Score)
					responsibleUnitScores[complianceAssessmentRequirementID].MaxPossibleScore += perReqMax
				}
			case "binary":
				responsibleUnitScores[complianceAssessmentRequirementID].ScoredRequirements++
				responsibleUnitScores[complianceAssessmentRequirementID].MaxPossibleScore += perReqMax
				if ru.ComplianceAssessmentFrameworkValueID != nil {
					if raw, ok := frameworkValuesMap[*ru.ComplianceAssessmentFrameworkValueID]; ok && raw > 0 {
						responsibleUnitScores[complianceAssessmentRequirementID].TotalScore += 1.0
					}
				}
			}

		}
		// If NO applicable RU and at least one NA → mark whole requirement NA
		if !applicableFound && notApplicableCount > 0 {
			if err := cras.repo.UpdateComplianceAssessmentRequirementStatus(
				req.ID, notApplicable.ID, 0, tx...,
			); err != nil {
				return nil, err
			}
			responsibleUnitScores[req.ID].CompliancePercentage = 0
			responsibleUnitScores[req.ID].Score = 0
			responsibleUnitScores[req.ID].MappedComplianceStatusID = &notApplicable.ID
			continue
		}
	}

	type clauseAgg struct{ total, compliant, partial, notImpl int }
	type clauseScoreAgg struct {
		sum   float64
		count int
	}
	clauseCounts := make(map[uint]*clauseAgg)      // clauseID -> counts
	clauseScores := make(map[uint]*clauseScoreAgg) // clauseID -> score sum/count

	// Calculate final percentages and statuses for each responsibleUnit
	var results []*complianceDTOs.ResponsibleUnitComplianceScoreDTO
	for _, responsibleUnitScore := range responsibleUnitScores {
		// fmt.Println("responsibleUnitScore ID:", responsibleUnitScore.ComplianceAssessmentRequirementID)
		if responsibleUnitScore.MappedComplianceStatusID != nil &&
			*responsibleUnitScore.MappedComplianceStatusID == notApplicable.ID {
			results = append(results, responsibleUnitScore)
			continue
		}
		if responsibleUnitScore.MaxPossibleScore > 0 {
			responsibleUnitScore.CompliancePercentage = (responsibleUnitScore.TotalScore / responsibleUnitScore.MaxPossibleScore) * 100
			responsibleUnitScore.Score = float64(int(responsibleUnitScore.CompliancePercentage*100)) / 100
			if responsibleUnitScore.CompliancePercentage == 100 {
				responsibleUnitScore.CompliancePercentage = 99.99
			}
		}

		// fmt.Println("Compliance Percentage:", responsibleUnitScore.CompliancePercentage)
		// fmt.Println("Score (Rounded):", responsibleUnitScore.Score)

		// Map to framework range → status (once)
		if assessment.ComplianceRegister != nil && assessment.ComplianceRegister.ComplianceAssessmentFramework != nil {
			for _, fv := range assessment.ComplianceRegister.ComplianceAssessmentFramework.Values {
				if fv != nil && fv.ScoreMin != nil && fv.ScoreMax != nil {
					min := float64(*fv.ScoreMin)
					max := float64(*fv.ScoreMax)
					score := responsibleUnitScore.CompliancePercentage

					// For the top band (highest max), include equality
					if score >= min && (score < max || (score == max && max == 100)) {
						responsibleUnitScore.MappedFrameworkValueID = &fv.ID
						responsibleUnitScore.MappedComplianceStatusID = &fv.ComplianceStatusID
						break
					}
				}

			}
		}
		results = append(results, responsibleUnitScore)
		// Update register status ONCE based on responsibleUnitScore mapping
		if responsibleUnitScore.MappedComplianceStatusID != nil {
			// fmt.Printf("Updating status: ComplianceAssessmentRequirementID=%v, StatusID=%v\n", responsibleUnitScore.ComplianceAssessmentRequirementID, *responsibleUnitScore.MappedComplianceStatusID)
			if err := cras.repo.UpdateComplianceAssessmentRequirementStatus(responsibleUnitScore.ComplianceAssessmentRequirementID, *responsibleUnitScore.MappedComplianceStatusID, responsibleUnitScore.Score, tx...); err != nil {
				return nil, err
			}
		}
		// Collect for clause roll-up
		if clauseID, ok := reqClause[responsibleUnitScore.ComplianceAssessmentRequirementID]; ok && clauseID != 0 {
			if clauseCounts[clauseID] == nil {
				clauseCounts[clauseID] = &clauseAgg{}
			}
			if clauseScores[clauseID] == nil {
				clauseScores[clauseID] = &clauseScoreAgg{}
			}

			cc := clauseCounts[clauseID]
			cs := clauseScores[clauseID]

			cc.total++
			switch {
			case responsibleUnitScore.MappedComplianceStatusID == nil:
				cc.notImpl++
			case *responsibleUnitScore.MappedComplianceStatusID == compliant.ID:
				cc.compliant++
			case *responsibleUnitScore.MappedComplianceStatusID == partial.ID:
				cc.partial++
			case *responsibleUnitScore.MappedComplianceStatusID == notImpl.ID:
				cc.notImpl++
			default:
				cc.partial++
			}

			// Score: use AR score (0–100). Treat missing as 0 and still count.
			cs.count++
			cs.sum += responsibleUnitScore.Score
		}
	}

	// Persist clause status + average score (one update per clause)
	// --- before clause loop: keep as-is ---
	// clauseCounts, clauseScores already exist

	// NEW: collect clause averages
	clauseAverages := make([]float64, 0)
	var (
		clausesTotal     int
		clausesCompliant int
		clausesPartial   int
		clausesNotImpl   int
	)
	// ========= OVERALL CLAUSE SCORE(based on REQ AVERAGES) =========
	// ---- persist clause status + average score (one update per clause) ----
	for clauseID, cc := range clauseCounts {
		// decide clause status (same logic)
		var setStatusID uint
		switch {
		case cc.partial > 0:
			setStatusID = partial.ID
			clausesPartial++
		case cc.total > 0 && cc.compliant == cc.total:
			setStatusID = compliant.ID
			clausesCompliant++
		case cc.total > 0 && cc.notImpl == cc.total:
			setStatusID = notImpl.ID
			clausesNotImpl++
		default:
			setStatusID = partial.ID
			clausesPartial++
		}
		clausesTotal++
		// average score (0–100), 2 dp
		cs := clauseScores[clauseID]
		avg := 0.0
		if cs != nil && cs.count > 0 {
			avg = cs.sum / float64(cs.count)
			avg = math.Round(avg*100) / 100
		}

		// persist clause status + score
		if err := cras.repo.UpdateClauseComplianceStatus(clauseID, setStatusID, avg, tx...); err != nil {
			return nil, err
		}

		// collect avg for overall
		clauseAverages = append(clauseAverages, avg)
	}

	// ========= OVERALL COMPLIANCE SCORE(based on CLAUSE AVERAGES) =========
	if clausesTotal > 0 {
		// overall score
		var sum float64
		for _, v := range clauseAverages {
			sum += v
		}
		overallPct := sum / float64(len(clauseAverages)) // already 0–100
		if overallPct == 100 {
			overallPct = 99.99
		}
		overallScore := float64(int(overallPct*100)) / 100

		// overall status (mirror clause logic)
		var overallStatusID uint
		switch {
		case clausesPartial > 0:
			overallStatusID = partial.ID
		case clausesCompliant == clausesTotal:
			overallStatusID = compliant.ID
		case clausesNotImpl == clausesTotal:
			overallStatusID = notImpl.ID
		default:
			overallStatusID = partial.ID
		}

		if err := cras.repo.UpdateComplianceStatus(assessment.ComplianceRegisterID, overallStatusID, overallScore, tx...); err != nil {
			return nil, err
		}
	}

	// fmt.Println("=== End CalculateComplianceScores ===")
	return results, nil
}

func (cras ComplianceAssessmentService) GetApprovalStatus(id uint, tx ...utility.DatabaseTransaction) (*complianceDTOs.ComplianceAssessmentApprovalStatusDTO, error) {

	model, err := cras.repo.GetFullDetailsByID(id, tx...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	isApproved := true
	for _, req := range model.AssessmentRequirements {
		// No responsible units means not approved
		if len(req.ResponsibleUnits) == 0 {
			isApproved = false
			break
		}

		for _, ru := range req.ResponsibleUnits {
			if ru.ComplianceAssessmentRequirementStatus == nil ||
				ru.ComplianceAssessmentRequirementStatus.Type != "approved" {
				isApproved = false
				break
			}
		}

		if !isApproved {
			break
		}
	}

	return &complianceDTOs.ComplianceAssessmentApprovalStatusDTO{
		AssessmentID: id,
		IsApproved:   isApproved,
	}, nil
}
