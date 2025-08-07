package services

import (
	"errors"
	"fmt"
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

	AssessmentCompletionSummary(assessmentID uint) (*complianceDTOs.AssessmentCompletionSummaryDTO, error)
	GetComplianceStatusSummary(assessmentID uint) ([]*complianceDTOs.ComplianceStatusCountDTO, error)
	GetComplianceScoreSummary(assessmentID uint) (*complianceDTOs.ComplianceScoreSummaryDTO, error)

	SubmitComplianceAssessment(assessmentID uint, currentUserID uint) error
	ApproveComplianceAssessment(assessmentID uint, currentUserID uint) error
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
	var complianceAssessmentDTO complianceDTOs.ComplianceAssessmentDetailDTO

	complianceAssessment, err := cras.repo.GetFullDetailsByID(id, tx...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	err = mapstructure.Decode(complianceAssessment, &complianceAssessmentDTO)
	if err != nil {
		return nil, err
	}

	// complianceAssessmentDTO.OrganizationUnitID = complianceAssessment.PrimaryUnitID
	complianceAssessmentDTO.CreatedAt = complianceAssessment.CreatedAt
	complianceAssessmentDTO.UpdatedAt = complianceAssessment.UpdatedAt
	// todo :- also add extra fields which required in details screen

	return &complianceAssessmentDTO, nil
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

	exists, err := cras.repo.IsAssessmentExists(0, complianceAssessmentFormDTO.ComplianceRegisterID, statusIds)
	if err != nil {
		return fmt.Errorf("failed to validate assessment %w", err)
	}

	if exists {
		return utility.NewValidationError(fmt.Errorf("an assessment already exists for this compliance register"), "unique", "compliance")
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

func (cras *ComplianceAssessmentService) AssessmentCompletionSummary(assessmentID uint) (*complianceDTOs.AssessmentCompletionSummaryDTO, error) {
	if _, err := cras.GetByID(assessmentID); err != nil {
		return nil, err
	}
	response, err := cras.repo.AssessmentCompletionSummary(assessmentID)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (cras *ComplianceAssessmentService) GetComplianceStatusSummary(assessmentID uint) ([]*complianceDTOs.ComplianceStatusCountDTO, error) {
	if _, err := cras.GetByID(assessmentID); err != nil {
		return nil, err
	}
	return cras.repo.GetComplianceStatusSummary(assessmentID)
}

func (cras *ComplianceAssessmentService) GetComplianceScoreSummary(assessmentID uint) (*complianceDTOs.ComplianceScoreSummaryDTO, error) {
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
	err = cras.UpdateComplianceRequirementValue(assessmentID, tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = cras.CalculateComplianceScores(assessmentID, tx)
	if err != nil {
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

// function for update compliance requiremnt status based on assessment status
func (cras *ComplianceAssessmentService) UpdateComplianceRequirementValue(assessmentID uint, tx ...utility.DatabaseTransaction) error {
	assessment, err := cras.repo.GetFullDetailsByID(assessmentID, tx...)
	if err != nil {
		return err
	}

	// Update compliance requirement status for each responsibleUnit
	for _, assessmentRequirement := range assessment.AssessmentRequirements {
		for _, assessmentRequirementResponsibleUnit := range assessmentRequirement.ResponsibleUnits {
			complianceRequirementValueUpdateDTO := complianceDTOs.ComplianceRequirementValueUpdateDTO{
				ComplianceRequirementID:              assessmentRequirement.ComplianceRequirementID,
				ComplianceAssessmentFrameworkValueID: *assessmentRequirementResponsibleUnit.ComplianceAssessmentFrameworkValueID,
				ResponsibleUnitID:                    assessmentRequirementResponsibleUnit.ResponsibleUnitID,
			}

			err = cras.repo.UpdateComplianceRequirementValue(complianceRequirementValueUpdateDTO, tx...)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (cras *ComplianceAssessmentService) CalculateComplianceScores(assessmentID uint, tx ...utility.DatabaseTransaction) ([]*complianceDTOs.ResponsibleUnitComplianceScoreDTO, error) {
	assessment, err := cras.repo.GetFullDetailsByID(assessmentID, tx...)
	if err != nil {
		return nil, err
	}

	// Get framework type and values map
	frameworkType := ""
	frameworkValuesMap := make(map[uint]int)
	if assessment.ComplianceRegister != nil && assessment.ComplianceRegister.ComplianceAssessmentFramework != nil {
		frameworkType = assessment.ComplianceRegister.ComplianceAssessmentFramework.AssessmentFrameworkType.Type
		for _, v := range assessment.ComplianceRegister.ComplianceAssessmentFramework.Values {
			if v != nil && v.ID > 0 && v.Value != nil {
				frameworkValuesMap[v.ID] = *v.Value
			}
		}
	}

	// Map to track responsibleUnit scores
	responsibleUnitScores := make(map[uint]*complianceDTOs.ResponsibleUnitComplianceScoreDTO)

	// Calculate scores for each responsibleUnit
	for _, assessmentRequirement := range assessment.AssessmentRequirements {
		for _, assessmentRequirementResponsibleUnit := range assessmentRequirement.ResponsibleUnits {
			responsibleUnitID := assessmentRequirementResponsibleUnit.ResponsibleUnitID

			// Check if requirement is applicable
			isApplicable := true
			if assessmentRequirementResponsibleUnit.ComplianceAssessmentFrameworkValueID != nil {
				isApplicable = assessmentRequirementResponsibleUnit.ComplianceAssessmentFrameworkValue.IsApplicable
			}

			// Only count applicable requirements
			if isApplicable {
				// Initialize responsibleUnit score if not exists
				if responsibleUnitScores[responsibleUnitID] == nil {
					responsibleUnitScores[responsibleUnitID] = &complianceDTOs.ResponsibleUnitComplianceScoreDTO{
						ComplianceRegisterID: assessment.ComplianceRegisterID,
						ResponsibleUnitID:    responsibleUnitID,
						TotalRequirements:    0,
						ScoredRequirements:   0,
						TotalScore:           0,
						MaxPossibleScore:     0,
						CompliancePercentage: 0,
						Score:                0,
						FrameworkType:        frameworkType,
					}
				}
				responsibleUnitScores[responsibleUnitID].TotalRequirements++

				// Calculate score based on framework type
				// if frameworkType == "point" {
				// 	// Point scoring: use framework value (0-2 scale)
				// 	if assessmentRequirementResponsibleUnit.ComplianceAssessmentFrameworkValueID != nil {
				// 		if val, ok := frameworkValuesMap[*assessmentRequirementResponsibleUnit.ComplianceAssessmentFrameworkValueID]; ok {
				// 			responsibleUnitScores[responsibleUnitID].ScoredRequirements++
				// 			responsibleUnitScores[responsibleUnitID].TotalScore += float64(val)
				// 			responsibleUnitScores[responsibleUnitID].MaxPossibleScore += 2.0 // Max score is 2 for point scoring
				// 		}
				// 	}
				// } else if frameworkType == "percentage" {
				// 	// Percentage scoring: use score field directly (percentage scale)
				// 	if assessmentRequirementResponsibleUnit.Score != nil {
				// 		responsibleUnitScores[responsibleUnitID].ScoredRequirements++
				// 		responsibleUnitScores[responsibleUnitID].TotalScore += float64(*assessmentRequirementResponsibleUnit.Score)
				// 		responsibleUnitScores[responsibleUnitID].MaxPossibleScore += 100.0 // Max score is 100 for percentage
				// 	}
				// }
				switch frameworkType {
				case "point":
					// Point scoring: use framework value (0-2 scale)
					if assessmentRequirementResponsibleUnit.ComplianceAssessmentFrameworkValueID != nil {
						if val, ok := frameworkValuesMap[*assessmentRequirementResponsibleUnit.ComplianceAssessmentFrameworkValueID]; ok {
							responsibleUnitScores[responsibleUnitID].ScoredRequirements++
							responsibleUnitScores[responsibleUnitID].TotalScore += float64(val)
							responsibleUnitScores[responsibleUnitID].MaxPossibleScore += 2.0
						}
					}
				case "percentage":
					// Percentage scoring: use score field directly (percentage scale)
					if assessmentRequirementResponsibleUnit.Score != nil {
						responsibleUnitScores[responsibleUnitID].ScoredRequirements++
						responsibleUnitScores[responsibleUnitID].TotalScore += float64(*assessmentRequirementResponsibleUnit.Score)
						responsibleUnitScores[responsibleUnitID].MaxPossibleScore += 100.0
					}
				case "binary":
					// Binary scoring: 1 = Yes, 0 = No, -1 = Not Applicable (already filtered out)
					if assessmentRequirementResponsibleUnit.ComplianceAssessmentFrameworkValueID != nil {
						if val, ok := frameworkValuesMap[*assessmentRequirementResponsibleUnit.ComplianceAssessmentFrameworkValueID]; ok {
							responsibleUnitScores[responsibleUnitID].ScoredRequirements++
							responsibleUnitScores[responsibleUnitID].TotalScore += float64(val)
							responsibleUnitScores[responsibleUnitID].MaxPossibleScore += 1.0 // Max score is 1 for binary
						}
					}
				case "maturity":
					// Maturity scoring: use framework value as maturity level (0–5)
					if assessmentRequirementResponsibleUnit.ComplianceAssessmentFrameworkValueID != nil {
						if val, ok := frameworkValuesMap[*assessmentRequirementResponsibleUnit.ComplianceAssessmentFrameworkValueID]; ok {
							responsibleUnitScores[responsibleUnitID].ScoredRequirements++
							responsibleUnitScores[responsibleUnitID].TotalScore += float64(val)
							responsibleUnitScores[responsibleUnitID].MaxPossibleScore += 5.0 // Max score is 5 for maturity
						}
					}
				}

			}
		}
	}

	// Calculate final percentages and statuses for each responsibleUnit
	var results []*complianceDTOs.ResponsibleUnitComplianceScoreDTO
	for _, responsibleUnitScore := range responsibleUnitScores {
		if responsibleUnitScore.MaxPossibleScore > 0 {
			responsibleUnitScore.CompliancePercentage = (responsibleUnitScore.TotalScore / responsibleUnitScore.MaxPossibleScore) * 100

			//score column is decimal(10,2)
			responsibleUnitScore.Score = float64(int(responsibleUnitScore.CompliancePercentage*100)) / 100
			//if i get 100% make it 99.99% because we checking max score less than 100
			if responsibleUnitScore.CompliancePercentage == 100 {
				responsibleUnitScore.CompliancePercentage = 99.99
			}

		}

		// Map score to framework value if framework exists
		if assessment.ComplianceRegister != nil && assessment.ComplianceRegister.ComplianceAssessmentFramework != nil {
			for _, frameworkValue := range assessment.ComplianceRegister.ComplianceAssessmentFramework.Values {
				if frameworkValue != nil && frameworkValue.ScoreMin != nil && frameworkValue.ScoreMax != nil {
					if responsibleUnitScore.CompliancePercentage >= float64(*frameworkValue.ScoreMin) &&
						responsibleUnitScore.CompliancePercentage < float64(*frameworkValue.ScoreMax) {
						responsibleUnitScore.MappedFrameworkValueID = &frameworkValue.ID
						responsibleUnitScore.MappedComplianceStatusID = &frameworkValue.ComplianceStatusID
						break
					}
				}
			}
		}

		results = append(results, responsibleUnitScore)
		//insted append please updaate compliance status in compliance register responsibleUnit table
		err = cras.repo.UpdateComplianceStatus(responsibleUnitScore.ComplianceRegisterID, *responsibleUnitScore.MappedComplianceStatusID, tx...)
		if err != nil {
			return nil, err
		}

	}

	return results, nil
}
