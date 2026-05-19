package services

import (
	"errors"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
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
	GetFullDetailsByIDOptimized(id uint, tx ...utility.DatabaseTransaction) (*complianceDTOs.ComplianceAssessmentDetailDTO, error)
	GetByID(id uint) (*models.ComplianceAssessment, error)
	Create(currentUserID uint, complianceAssessmentFormDTO *complianceDTOs.ComplianceAssessmentFormDTO) error
	Update(id uint, currentUserID uint, compliance *complianceDTOs.ComplianceAssessmentUpdateFormDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	ForceDelete(ids []uint, currentUserID uint) ([]uint, error)

	AssessmentCompletionRollup(assessmentID uint) (*complianceDTOs.AssessmentCompletionRollupDTO, error)
	GetComplianceStatusSummary(assessmentID uint) ([]complianceDTOs.ComplianceStatusCountDTO, error)
	GetComplianceScoreSummary(assessmentID uint, params *complianceDTOs.ComplianceAssessmentScoreSummaryQueryParamsDTO) (*complianceDTOs.ComplianceScoreSummaryDTO, error)

	SubmitComplianceAssessment(assessmentID uint, currentUserID uint) error
	ApproveComplianceAssessment(assessmentID uint, currentUserID uint) error

	GetApprovalStatus(id uint, tx ...utility.DatabaseTransaction) (*complianceDTOs.ComplianceAssessmentApprovalStatusDTO, error)
	GetControlApprovalStatus(id uint, tx ...utility.DatabaseTransaction) (*complianceDTOs.ComplianceAssessmentApprovalStatusDTO, error)

	ExportComplianceAssessments(data []*complianceDTOs.ComplianceAssessmentListDTO, currentUserID uint) (*os.File, error)
}

// Service and implementations
type ComplianceAssessmentService struct {
	repo                          complianceRepositories.ComplianceAssessmentRepositoryInterface
	activityLogSvc                logServices.ActivityLogServiceInterface
	userSvc                       organizationServices.UserServiceInterface
	userSettingSvc                organizationServices.UserSettingServiceInterface
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
	userSettingSvc organizationServices.UserSettingServiceInterface,
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
		userSettingSvc:                userSettingSvc,
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

	allControlsCompleted := true
	allControlsApproved := true

	for _, ctrlModel := range model.AssessmentControls {
		if !ctrlModel.IsCompleted {
			allControlsCompleted = false
		}
		if ctrlModel.ComplianceAssessmentRequirementStatus == nil || ctrlModel.ComplianceAssessmentRequirementStatus.Type != "approved" {
			allControlsApproved = false
		}
	}

	// Merge requirement + control flags
	dto.IsCompleted = allCompleted && allControlsCompleted
	dto.IsApproved = allApproved && allControlsApproved

	dto.CreatedAt = model.CreatedAt
	dto.UpdatedAt = model.UpdatedAt

	return &dto, nil
}

// service/compliance_assessment_service.go
func (cras *ComplianceAssessmentService) GetFullDetailsByIDOptimized(id uint, tx ...utility.DatabaseTransaction) (*complianceDTOs.ComplianceAssessmentDetailDTO, error) {
	model, err := cras.repo.GetFullDetailsByIDOptimized(id, tx...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	isCompleted, isApproved, err := cras.repo.ComputeFlagsFast(id, tx...)
	if err != nil {
		return nil, err
	}
	var dto complianceDTOs.ComplianceAssessmentDetailDTO
	if err := mapstructure.Decode(model, &dto); err != nil {
		return nil, err
	}

	dto.IsCompleted = isCompleted
	dto.IsApproved = isApproved
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
	// fmt.Println(compliance.IsRiskBasedAssessment)
	//
	// Only create compliance requirements and controls if it's a risk-based assessment
	// if compliance.IsRiskBasedAssessment {
	for _, clause := range compliance.ComplianceClauses {
		clauseDetails, err := cras.complianceClauseSvc.GetFullDetailsByID(clause.ID)
		if err != nil {
			tx.Rollback()
			return err
		}
		// ====================================
		// CASE 1: Requirement-based assessment
		// ====================================
		if compliance.IsRequirementAssessment {
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
		// =================================
		// CASE 2: Control-based assessment
		// =================================
		// if compliance.IsControlAssessment {
		// 	for _, control := range clauseDetails.ComplianceControls {
		// 		assessmentControl := &models.ComplianceAssessmentControl{
		// 			ComplianceAssessmentID: id,
		// 			ComplianceControlID:    control.ID,
		// 			IsCompleted:            false,
		// 			CreatedBYID:            currentUserID,
		// 		}
		// 		if err := cras.repo.CreateComplianceAssessmentControl(assessmentControl, tx); err != nil {
		// 			tx.Rollback()
		// 			return err
		// 		}
		// 	}
		// }
	}
	// }

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

		// Send notification and email to review user if assigned
		if complianceAssessmentFormDTO.ReviewUserID != 0 {
			// Prepare notification
			urlParam := fmt.Sprintf(`{"id":%d}`, id)
			queryParam := ""
			entity := "compliance-assessment"
			action := "review-user-assign"

			dbNotificationDTO := &commonDTOs.DBNotificationDTO{
				UserIDs:              []uint{complianceAssessmentFormDTO.ReviewUserID},
				ModuleIdentifier:     "compliance",
				NotificationType:     string(commonDTOs.AssignmentType),
				HeadingEn:            "Compliance Assessment Review Assigned",
				HeadingAr:            "تم تعيين مراجعة لتقييم الامتثال",
				TitleEn:              compliance.TitleEn,
				TitleAr:              compliance.TitleAr,
				ReferenceNumber:      compliance.ReferenceNumber,
				Entity:               &entity,
				Action:               &action,
				URLParam:             urlParam,
				QueryParam:           queryParam,
				IsSystemNotification: true,
				CreatedByID:          &currentUserID,
			}

			// Prepare email data
			reviewUser, err := cras.userSvc.GetByID(complianceAssessmentFormDTO.ReviewUserID)
			if err != nil {
				utility.LogError("Failed to fetch review user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}

			assignedBy, err := cras.userSvc.GetByID(currentUserID)
			if err != nil {
				utility.LogError("Failed to fetch assigning user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}

			emailTo := []string{reviewUser.Email}
			subject := fmt.Sprintf("Action Required: Review Compliance Assessment [%s] - %s", compliance.ReferenceNumber, compliance.TitleEn)
			appURL := utility.GetAppURL()
			emailData := map[string]interface{}{
				"FirstNameEn":     reviewUser.FirstNameEn,
				"LastNameEn":      reviewUser.LastNameEn,
				"TitleEn":         compliance.TitleEn,
				"TitleAr":         compliance.TitleAr,
				"ReferenceNumber": compliance.ReferenceNumber,
				"AssignedBy":      fmt.Sprintf("%s %s", assignedBy.FirstNameEn, assignedBy.LastNameEn),
				"AssignedOn":      time.Now().Format("2006-01-02"),
				"Description":     "N/A",
				"RedirectURL":     fmt.Sprintf("%s/compliance/assessments/%d/overview", appURL, id),
			}

			// Set description if available
			if complianceAssessmentFormDTO.Description != nil {
				emailData["Description"] = *complianceAssessmentFormDTO.Description
			}

			emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
				EmailTo:              emailTo,
				Subject:              subject,
				TemplateRelativePath: "compliance/templates/email/compliance_assessment_review_user_assign_mail.html",
				EmailData:            emailData,
			}

			// Send notification and email
			if err := cras.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "compliance_assessment"); err != nil {
				utility.LogError("Failed to send notifications", map[string]interface{}{
					"error": err.Error(),
				})
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
				"compliance", "compliance_assessment", "", nil,
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

		// Get the old review user ID for comparison
		oldReviewUserID := oldComplianceAssessment.ReviewUserID

		// sends notification if review user has changed
		if oldReviewUserID != complianceAssessmentUpdateFormDTO.ReviewUserID {
			urlParam := fmt.Sprintf(`{"id":%d}`, id)
			queryParam := ""
			entity := "compliance-assessment"

			if oldReviewUserID != 0 {
				oldReviewUserAction := "review-user-remove"
				oldReviewUserNotificationDTO := &commonDTOs.DBNotificationDTO{
					UserIDs:              []uint{oldReviewUserID},
					ModuleIdentifier:     "compliance",
					NotificationType:     string(commonDTOs.InformationType),
					HeadingEn:            "Compliance Assessment Review Responsibility Removed",
					HeadingAr:            "تمت إزالة مسؤولية مراجعة تقييم الامتثال",
					TitleEn:              oldComplianceAssessment.ComplianceRegister.TitleEn,
					TitleAr:              oldComplianceAssessment.ComplianceRegister.TitleAr,
					ReferenceNumber:      oldComplianceAssessment.ComplianceRegister.ReferenceNumber,
					Entity:               &entity,
					Action:               &oldReviewUserAction,
					URLParam:             urlParam,
					QueryParam:           queryParam,
					IsSystemNotification: false,
					CreatedByID:          &currentUserID,
				}

				oldReviewUser, err := cras.userSvc.GetByID(oldReviewUserID)
				if err != nil {
					utility.LogError("Failed to fetch old review user", map[string]interface{}{
						"error": err.Error(),
					})
					return
				}

				appURL := utility.GetAppURL()
				oldReviewUserEmailTo := []string{oldReviewUser.Email}
				oldReviewUserSubject := fmt.Sprintf("Compliance Assessment/%s: Review Responsibility Removed", oldComplianceAssessment.ComplianceRegister.TitleEn)
				oldReviewUserTemplateRelativePath := "compliance/templates/email/compliance_assessment_remove_review_user_mail.html"
				oldReviewUserEmailData := map[string]interface{}{
					"FirstNameEn":     oldReviewUser.FirstNameEn,
					"LastNameEn":      oldReviewUser.LastNameEn,
					"ReferenceNumber": oldComplianceAssessment.ComplianceRegister.ReferenceNumber,
					"TitleEn":         oldComplianceAssessment.ComplianceRegister.TitleEn,
					"TitleAr":         oldComplianceAssessment.ComplianceRegister.TitleAr,
					"Description":     "N/A",
					"RedirectURL":     fmt.Sprintf("%s/compliance/assessments/%d/overview", appURL, id),
				}

				if complianceAssessmentUpdateFormDTO.Description != nil {
					oldReviewUserEmailData["Description"] = *complianceAssessmentUpdateFormDTO.Description
				}

				oldReviewUserEmailNotificationDTO := &commonDTOs.EmailNotificationDTO{
					EmailTo:              oldReviewUserEmailTo,
					Subject:              oldReviewUserSubject,
					TemplateRelativePath: oldReviewUserTemplateRelativePath,
					EmailData:            oldReviewUserEmailData,
				}

				// send old review user notification + email
				err = cras.notificationSvc.SendNotification(
					oldReviewUserNotificationDTO, oldReviewUserEmailNotificationDTO, "compliance_assessment",
				)
				if err != nil {
					utility.LogError("Failed to send old review user notifications", map[string]interface{}{
						"error": err.Error(),
					})
				}
			}

			if complianceAssessmentUpdateFormDTO.ReviewUserID != 0 {
				newReviewUserAction := "review-user-assign"
				newReviewUserNotificationDTO := &commonDTOs.DBNotificationDTO{
					UserIDs:              []uint{complianceAssessmentUpdateFormDTO.ReviewUserID},
					ModuleIdentifier:     "compliance",
					NotificationType:     string(commonDTOs.AssignmentType),
					HeadingEn:            "Compliance Assessment Review Assigned",
					HeadingAr:            "تم تعيين مراجعة لتقييم الامتثال",
					TitleEn:              existingAssessment.ComplianceRegister.TitleEn,
					TitleAr:              existingAssessment.ComplianceRegister.TitleAr,
					ReferenceNumber:      existingAssessment.ComplianceRegister.ReferenceNumber,
					Entity:               &entity,
					Action:               &newReviewUserAction,
					URLParam:             urlParam,
					QueryParam:           queryParam,
					IsSystemNotification: false,
					CreatedByID:          &currentUserID,
				}

				newReviewUser, err := cras.userSvc.GetByID(complianceAssessmentUpdateFormDTO.ReviewUserID)
				if err != nil {
					utility.LogError("Failed to fetch new review user", map[string]interface{}{
						"error": err.Error(),
					})
					return
				}

				assignedBy, err := cras.userSvc.GetByID(currentUserID)
				if err != nil {
					utility.LogError("Failed to fetch assigning user", map[string]interface{}{
						"error": err.Error(),
					})
					return
				}

				appURL := utility.GetAppURL()
				newReviewUserEmailTo := []string{newReviewUser.Email}
				newReviewUserSubject := fmt.Sprintf("Action Required: Review Compliance Assessment [%s] - %s", existingAssessment.ComplianceRegister.ReferenceNumber, existingAssessment.ComplianceRegister.TitleEn)
				newReviewUserTemplateRelativePath := "compliance/templates/email/compliance_assessment_review_user_assign_mail.html"
				newReviewUserEmailData := map[string]interface{}{
					"FirstNameEn":     newReviewUser.FirstNameEn,
					"LastNameEn":      newReviewUser.LastNameEn,
					"TitleEn":         existingAssessment.ComplianceRegister.TitleEn,
					"TitleAr":         existingAssessment.ComplianceRegister.TitleAr,
					"ReferenceNumber": existingAssessment.ComplianceRegister.ReferenceNumber,
					"AssignedBy":      fmt.Sprintf("%s %s", assignedBy.FirstNameEn, assignedBy.LastNameEn),
					"AssignedOn":      time.Now().Format("2006-01-02"),
					"Description":     "N/A",
					"RedirectURL":     fmt.Sprintf("%s/compliance/assessments/%d/overview", appURL, id),
				}

				if complianceAssessmentUpdateFormDTO.Description != nil {
					newReviewUserEmailData["Description"] = *complianceAssessmentUpdateFormDTO.Description
				}

				newReviewUserEmailNotificationDTO := &commonDTOs.EmailNotificationDTO{
					EmailTo:              newReviewUserEmailTo,
					Subject:              newReviewUserSubject,
					TemplateRelativePath: newReviewUserTemplateRelativePath,
					EmailData:            newReviewUserEmailData,
				}

				// send new review user notification + email
				err = cras.notificationSvc.SendNotification(
					newReviewUserNotificationDTO, newReviewUserEmailNotificationDTO, "compliance_assessment",
				)
				if err != nil {
					utility.LogError("Failed to send new review user notifications", map[string]interface{}{
						"error": err.Error(),
					})
				}
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
		if len(complianceAssessment.AssessmentControls) > 0 {
			var deleteControlIds []*uint
			for _, assmntControl := range complianceAssessment.AssessmentControls {
				deleteControlIds = append(deleteControlIds, assmntControl.ID)
			}

			if len(deleteControlIds) > 0 {
				if err := cras.repo.DeleteComplianceAssessmentControls(deleteControlIds); err != nil {
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

func (cras *ComplianceAssessmentService) ForceDelete(ids []uint, currentUserID uint) ([]uint, error) {

	tx := cras.repo.BeginTransaction()
	var failedIds []uint
	var deletedData []map[string]interface{}

	for _, id := range ids {
		if err := cras.repo.ForceDelete(id, tx); err != nil {
			failedIds = append(failedIds, id)
			continue
		}

		deletedData = append(deletedData, map[string]interface{}{
			"id": id,
		})
	}

	if len(failedIds) > 0 {
		tx.Rollback()
		return failedIds, nil
	}

	tx.Commit()

	go func() {
		err := cras.activityLogSvc.LogBulkActivities(
			models.DeleteAction,
			models.ComplianceAssessment{}.TableName(),
			ids,
			deletedData,
			currentUserID,
			nil,
		)
		if err != nil {
			// @todo - save to log file
			fmt.Printf("Failed to log bulk delete activity: %v\n", err)
		}
	}()

	return nil, nil
}

func (cras *ComplianceAssessmentService) validateComplianceAssessmentFormDTO(complianceAssessmentFormDTO *complianceDTOs.ComplianceAssessmentFormDTO) error {
	complianceRegister, err := cras.complianceRegisterSvc.GetByID(complianceAssessmentFormDTO.ComplianceRegisterID)
	if err != nil {
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

	// validate based on register type
	// if complianceRegister.IsRiskBasedAssessment {
	if complianceRegister.IsRequirementAssessment {
		requirements, err := cras.repo.GetRequirementsForRegister(complianceAssessmentFormDTO.ComplianceRegisterID)
		if err != nil {
			return fmt.Errorf("failed to validate requirements: %w", err)
		}

		if len(requirements) == 0 {
			return utility.NewValidationError(
				fmt.Errorf("no compliance requirements found for this compliance register"),
				"requirement_not_exist",
				"compliance requirement",
			)
		}

		hasRiskBasedRequirement := false
		for _, requirement := range requirements {
			if requirement == nil {
				continue
			}

			for _, responsibleUnit := range requirement.ResponsibleUnits {
				if responsibleUnit != nil && responsibleUnit.IsRiskBasedAssessment != nil && *responsibleUnit.IsRiskBasedAssessment {
					hasRiskBasedRequirement = true
					break
				}
			}

			if hasRiskBasedRequirement {
				break
			}
		}

		if hasRiskBasedRequirement {
			allRiskBasedHaveInherentScore, err := cras.repo.CheckAllRiskBasedRequirementsHaveInherentScore(complianceAssessmentFormDTO.ComplianceRegisterID)
			if err != nil {
				return fmt.Errorf("failed to validate inherent risk scores: %w", err)
			}
			if !allRiskBasedHaveInherentScore {
				return utility.NewValidationError(
					fmt.Errorf("all risk-based requirements must have inherent risk analysis completed before creating assessment"),
					"inherent_risk_analysis_incomplete",
					"compliance requirement",
				)
			}
		}
	}
	// }

	// if complianceRegister.IsRiskBasedAssessment {
	// 	if complianceRegister.IsControlAssessment {
	// 		hasControl, err := cras.repo.HasControlForRegister(complianceAssessmentFormDTO.ComplianceRegisterID)
	// 		if err != nil {
	// 			return fmt.Errorf("failed to validate controls: %w", err)
	// 		}
	// 		if !hasControl {
	// 			return utility.NewValidationError(
	// 				fmt.Errorf("no compliance controls found for this compliance register"),
	// 				"control_not_exist",
	// 				"compliance control",
	// 			)
	// 		}
	// 	}
	// }

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
	assessment, err := cras.repo.GetFullDetailsByID(assessmentID)
	if err != nil {
		return nil, err
	}
	if assessment.ComplianceRegister == nil {
		return nil, fmt.Errorf("assessment %d has no compliance register loaded", assessmentID)
	}

	switch {
	case assessment.ComplianceRegister.IsControlAssessment:
		return cras.repo.AssessmentControlCompletionRollup(assessmentID)
	case assessment.ComplianceRegister.IsRequirementAssessment:
		return cras.repo.AssessmentRequirementCompletionRollup(assessmentID)
	default:
		// fallback: zero summary
		return &complianceDTOs.AssessmentCompletionRollupDTO{
			Total:                0,
			Completed:            0,
			NotCompleted:         0,
			CompletionPercentage: 0,
		}, nil
	}
}

func (cras *ComplianceAssessmentService) GetComplianceStatusSummary(assessmentID uint) ([]complianceDTOs.ComplianceStatusCountDTO, error) {
	assessment, err := cras.repo.GetFullDetailsByID(assessmentID)
	if err != nil {
		return nil, err
	}
	if assessment.ComplianceRegister == nil {
		return nil, fmt.Errorf("assessment %d has no compliance register loaded", assessmentID)
	}

	switch {
	case assessment.ComplianceRegister.IsControlAssessment:
		return cras.repo.GetComplianceStatusSummaryByControl(assessmentID)
	case assessment.ComplianceRegister.IsRequirementAssessment:
		return cras.repo.GetComplianceStatusSummaryByRequirement(assessmentID)
	default:
		return []complianceDTOs.ComplianceStatusCountDTO{}, nil
	}
}

func (cras *ComplianceAssessmentService) GetComplianceScoreSummary(assessmentID uint, params *complianceDTOs.ComplianceAssessmentScoreSummaryQueryParamsDTO) (*complianceDTOs.ComplianceScoreSummaryDTO, error) {
	assessment, err := cras.repo.GetFullDetailsByID(assessmentID)
	if err != nil {
		return nil, err
	}
	if assessment.ComplianceRegister == nil {
		return nil, fmt.Errorf("assessment %d has no compliance register loaded", assessmentID)
	}
	var ruIDs []uint
	if params != nil {
		if raw := params.Filters["responsible_unit_id"]; len(raw) > 0 {
			ruIDs = parseUintList(raw) // helper below
		}
	}

	var total int
	var avg float64

	switch {
	case assessment.ComplianceRegister.IsControlAssessment:
		total, avg, err = cras.computeControlSnapshot(assessmentID)
	case assessment.ComplianceRegister.IsRequirementAssessment:
		total, avg, err = cras.computeRequirementSnapshot(assessmentID, ruIDs...)
	default:
		return &complianceDTOs.ComplianceScoreSummaryDTO{TotalRequirements: 0, AverageScore: 0}, nil
	}
	if err != nil {
		return nil, err
	}

	return &complianceDTOs.ComplianceScoreSummaryDTO{
		TotalRequirements: total,
		AverageScore:      avg, // already rounded and 99.99-capped when needed
	}, nil
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
			"compliance", "compliance_assessment", "", nil,
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

	assessment, err := cras.repo.GetByID(assessmentID)
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

	// Get assessment details to calculate average residual score and rating
	assessmentDetails, err := cras.repo.GetFullDetailsByID(assessmentID, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	// STEP 1: First, update all requirement averages from their responsible units
	for _, req := range assessmentDetails.AssessmentRequirements {
		if req == nil {
			continue
		}

		// Calculate average residual score from all responsible units for this requirement
		var totalResidualScore float64
		var countWithScore int

		for _, ru := range req.ResponsibleUnits {
			if ru != nil && ru.ResidualRiskScore != nil {
				totalResidualScore += *ru.ResidualRiskScore
				countWithScore++

			}
		}

		// If at least one responsible unit has a residual score, calculate and update requirement average
		if countWithScore > 0 {
			requirementAverageScore := totalResidualScore / float64(countWithScore)

			// Get the rating based on the average score
			requirementRatingID, err := cras.getResidualRiskRatingByScore(requirementAverageScore)
			if err != nil {
				tx.Rollback()
				return err
			}

			// Update the requirement with average residual score and rating
			if err := cras.repo.UpdateRequirementAverageResidualScoreAndRating(req.ID, requirementAverageScore, requirementRatingID, tx); err != nil {
				tx.Rollback()
				return err
			}

			// Update the in-memory object so we can use it for assessment calculation
			req.AverageResidualRiskScore = &requirementAverageScore
			req.AverageResidualRiskRatingID = requirementRatingID
		}
	}

	// STEP 2: Now calculate assessment average from all requirement averages
	var totalResidualScore float64
	var countWithScore int

	for _, req := range assessmentDetails.AssessmentRequirements {
		if req != nil && req.AverageResidualRiskScore != nil {
			totalResidualScore += *req.AverageResidualRiskScore
			countWithScore++
		}
	}

	// If at least one requirement has a residual score, calculate average and cascade
	if countWithScore > 0 {
		assessmentAverageScore := totalResidualScore / float64(countWithScore)

		// Get the rating based on the average score
		assessmentRatingID, err := cras.getResidualRiskRatingByScore(assessmentAverageScore)
		if err != nil {
			tx.Rollback()
			return err
		}

		// Update the assessment with average residual score and rating
		if err := cras.repo.UpdateAverageResidualScoreAndRating(assessmentID, assessmentAverageScore, assessmentRatingID, tx); err != nil {
			tx.Rollback()
			return err
		}

		// STEP 3: Finally cascade to compliance register - use the same assessment average
		registerID := assessmentDetails.ComplianceRegisterID

		if err := cras.repo.UpdateRegisterAverageResidualScoreAndRating(registerID, assessmentAverageScore, assessmentRatingID, tx); err != nil {
			tx.Rollback()
			return err
		}
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
			"compliance", "compliance_assessment", "", nil,
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

func (cras *ComplianceAssessmentService) UpdateComplianceRequirementValuesFromModel(assessmentID uint, tx ...utility.DatabaseTransaction) error {
	// Load assessment with requirements + RUs (your existing GetFullDetailsByID already preloads these)
	assessment, err := cras.repo.GetFullDetailsByIDOptimized(assessmentID, tx...)
	if err != nil {
		return err
	}
	if assessment == nil || len(assessment.AssessmentRequirements) == 0 {
		return nil
	}

	for _, ar := range assessment.AssessmentRequirements {
		// ---------- 1) Push RU selections/audit up to ComplianceRequirementResponsibleUnit ----------
		for _, ru := range ar.ResponsibleUnits {
			// Skip if there is literally nothing to persist for this RU
			hasFV := ru.ComplianceAssessmentFrameworkValueID != nil
			if !hasFV && ru.Comments == nil && ru.CompletedByUserID == nil && ru.CompletedAt == nil {
				continue
			}

			var fvID uint
			if ru.ComplianceAssessmentFrameworkValueID != nil {
				fvID = *ru.ComplianceAssessmentFrameworkValueID
			}

			dtoRU := complianceDTOs.ComplianceRequirementValueUpdateDTO{
				ComplianceRequirementID:              ar.ComplianceRequirementID,
				ResponsibleUnitID:                    ru.ResponsibleUnitID,
				ComplianceAssessmentFrameworkValueID: fvID,        // 0 means "leave unchanged" per repo logic
				Comments:                             ru.Comments, // pass through
				CompletedByUserID:                    ru.CompletedByUserID,
				CompletedAt:                          ru.CompletedAt,
			}

			if err := cras.repo.UpdateComplianceRequirementResponsibleUintValue(dtoRU, tx...); err != nil {
				return err
			}
		}

		// ---------- 2) Push requirement status/score up to ComplianceRequirement ----------
		// Take status from the AR row (prefer explicit ID; fallback to loaded relation)
		var statusID uint
		if ar.ComplianceStatusID != nil {
			statusID = *ar.ComplianceStatusID
		} else if ar.ComplianceStatus != nil {
			statusID = ar.ComplianceStatus.ID
		}

		// Score on AR is *int in your DTO/model; convert to float64 for parent table
		var scoreF float64
		if ar.Score != nil {
			scoreF = float64(*ar.Score)
		} else {
			// keep 0.0 if not present; adjust if you'd rather leave it unchanged
			scoreF = 0.0
		}

		dtoReq := complianceDTOs.ComplianceRequirementValueUpdateDTO{
			ComplianceRequirementID: ar.ComplianceRequirementID,
			ComplianceStatusID:      &statusID,
			Score:                   &scoreF,
		}

		if err := cras.repo.UpdateComplianceRequirementStatusUpdate(dtoReq, tx...); err != nil {
			return err
		}
	}

	return nil
}

// UpdateComplianceControlValuesFromModel syncs the *current* values from the
// assessment graph (DTO) into your control tables. This is intended to be called
// during final approve, not during calculation.
func (cras *ComplianceAssessmentService) UpdateComplianceControlValuesFromModel(assessmentID uint, tx ...utility.DatabaseTransaction) error {
	// Load the assessment with controls (optimized preload)
	assessment, err := cras.repo.GetFullDetailsByIDOptimized(assessmentID, tx...)
	if err != nil {
		return err
	}
	if assessment == nil || len(assessment.AssessmentControls) == 0 {
		return nil
	}

	for _, ac := range assessment.AssessmentControls {
		// Strong null-safety: skip anything we can’t reliably map
		if ac == nil || ac.ComplianceControl == nil {
			continue
		}

		// Map status from the selected framework value if present
		var statusID *uint
		if ac.ComplianceAssessmentFrameworkValue != nil {
			sid := ac.ComplianceAssessmentFrameworkValue.ComplianceStatusID
			statusID = &sid
		}

		// Prepare the sync DTO from the DTO node (no extra logic here)
		dto := complianceDTOs.ControlSyncUpdateDTO{
			ControlID:                            ac.ComplianceControl.ControlID, // base control id
			ComplianceControlID:                  ac.ComplianceControlID,         // linking table id
			ComplianceAssessmentFrameworkValueID: ac.ComplianceAssessmentFrameworkValueID,
			CompletedByUserID:                    ac.CompletedByUserID,
			CompletedAt:                          ac.CompletedAt,
			ComplianceStatusID:                   statusID,
		}

		// 1) Push status (band) onto the control’s compliance status record
		if err := cras.repo.UpdateControlComplianceStatus(dto, tx...); err != nil {
			return err
		}

		// 2) Push metadata (completed by / completed at) to the control meta table
		if err := cras.repo.UpdateComplianceControlMeta(dto, tx...); err != nil {
			return err
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

func (cras *ComplianceAssessmentService) CalculateComplianceRequirementScores(assessmentID uint, tx ...utility.DatabaseTransaction) ([]*complianceDTOs.ResponsibleUnitComplianceScoreDTO, error) {
	assessment, err := cras.repo.GetFullDetailsByID(assessmentID, tx...)
	if err != nil {
		return nil, err
	}
	notApplicable, err := cras.complianceStatusSvc.GetByType("not_applicable")
	if err != nil {
		return nil, err
	}

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

	// Map to track responsibleUnit scores
	responsibleUnitScores := make(map[uint]*complianceDTOs.ResponsibleUnitComplianceScoreDTO)
	reqClause := make(map[uint]uint)

	// Walk all requirement-responsibleUnit rows and accumulate globally
	for _, req := range assessment.AssessmentRequirements {
		if req.ComplianceRequirement != nil {
			reqClause[req.ID] = req.ComplianceRequirement.ComplianceClauseID
		}
		ComplianceClauseID := req.ComplianceRequirement.ComplianceClauseID

		applicableFound := false
		notApplicableCount := 0

		for _, ru := range req.ResponsibleUnits {
			complianceAssessmentRequirementID := ru.ComplianceAssessmentRequirementID

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
			if ru.ComplianceAssessmentFrameworkValueID != nil {
				isApplicable = ru.ComplianceAssessmentFrameworkValue.IsApplicable
			}
			if !isApplicable {
				notApplicableCount++
				continue
			}

			applicableFound = true
			responsibleUnitScores[complianceAssessmentRequirementID].TotalRequirements++

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

	// Clause aggregation structures
	type clauseScoreAgg struct {
		sum   float64
		count int
	}
	clauseScores := make(map[uint]*clauseScoreAgg)

	// Calculate final percentages and statuses for each responsibleUnit
	var results []*complianceDTOs.ResponsibleUnitComplianceScoreDTO
	for _, responsibleUnitScore := range responsibleUnitScores {
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

		// Map requirement score % → status
		if assessment.ComplianceRegister != nil && assessment.ComplianceRegister.ComplianceAssessmentFramework != nil {
			for _, fv := range assessment.ComplianceRegister.ComplianceAssessmentFramework.Values {
				if fv != nil && fv.ScoreMin != nil && fv.ScoreMax != nil {
					min := float64(*fv.ScoreMin)
					max := float64(*fv.ScoreMax)
					score := responsibleUnitScore.CompliancePercentage
					if score >= min && (score < max || (score == max && max == 100)) {
						responsibleUnitScore.MappedFrameworkValueID = &fv.ID
						responsibleUnitScore.MappedComplianceStatusID = &fv.ComplianceStatusID
						break
					}
				}
			}
		}
		results = append(results, responsibleUnitScore)

		if responsibleUnitScore.MappedComplianceStatusID != nil {
			if err := cras.repo.UpdateComplianceAssessmentRequirementStatus(
				responsibleUnitScore.ComplianceAssessmentRequirementID,
				*responsibleUnitScore.MappedComplianceStatusID,
				responsibleUnitScore.Score,
				tx...,
			); err != nil {
				return nil, err
			}
		}

		// Collect clause averages
		if clauseID, ok := reqClause[responsibleUnitScore.ComplianceAssessmentRequirementID]; ok && clauseID != 0 {
			if clauseScores[clauseID] == nil {
				clauseScores[clauseID] = &clauseScoreAgg{}
			}
			cs := clauseScores[clauseID]
			cs.count++
			cs.sum += responsibleUnitScore.Score
		}
	}

	// Persist clause status + average score
	clauseAverages := make([]float64, 0)
	for clauseID, cs := range clauseScores {
		avg := 0.0
		if cs != nil && cs.count > 0 {
			avg = cs.sum / float64(cs.count)
			avg = math.Round(avg*100) / 100
		}

		// Map clause % → status
		var clauseStatusID uint
		if assessment.ComplianceRegister != nil && assessment.ComplianceRegister.ComplianceAssessmentFramework != nil {
			for _, fv := range assessment.ComplianceRegister.ComplianceAssessmentFramework.Values {
				if fv != nil && fv.ScoreMin != nil && fv.ScoreMax != nil {
					min := float64(*fv.ScoreMin)
					max := float64(*fv.ScoreMax)
					if avg >= min && (avg < max || (avg == max && max == 100)) {
						clauseStatusID = fv.ComplianceStatusID
						break
					}
				}
			}
		}

		if err := cras.repo.UpdateClauseComplianceStatus(clauseID, clauseStatusID, avg, tx...); err != nil {
			return nil, err
		}

		clauseAverages = append(clauseAverages, avg)
	}

	// ========= OVERALL COMPLIANCE SCORE (based on CLAUSE AVERAGES) =========
	if len(clauseAverages) > 0 {
		var sum float64
		for _, v := range clauseAverages {
			sum += v
		}
		overallPct := sum / float64(len(clauseAverages))
		if overallPct == 100 {
			overallPct = 99.99
		}
		overallScore := float64(int(overallPct*100)) / 100

		// Map overall % → status
		var overallStatusID uint
		if assessment.ComplianceRegister != nil && assessment.ComplianceRegister.ComplianceAssessmentFramework != nil {
			for _, fv := range assessment.ComplianceRegister.ComplianceAssessmentFramework.Values {
				if fv != nil && fv.ScoreMin != nil && fv.ScoreMax != nil {
					min := float64(*fv.ScoreMin)
					max := float64(*fv.ScoreMax)
					if overallPct >= min && (overallPct < max || (overallPct == max && max == 100)) {
						overallStatusID = fv.ComplianceStatusID
						break
					}
				}
			}
		}

		if err := cras.repo.UpdateComplianceStatus(assessment.ComplianceRegisterID, overallStatusID, overallScore, tx...); err != nil {
			return nil, err
		}
	}

	return results, nil
}

func (cras *ComplianceAssessmentService) CalculateComplianceControlScores(assessmentID uint, tx ...utility.DatabaseTransaction) ([]*complianceDTOs.ControlComplianceScoreDTO, error) {
	assessment, err := cras.repo.GetFullDetailsByID(assessmentID, tx...)
	if err != nil {
		return nil, err
	}

	// ---------------- Framework setup ----------------
	frameworkType := ""
	var frameworkValues []*models.ComplianceAssessmentFrameworkValue
	frameworkValuesMap := make(map[uint]int) // framework_value_id -> numeric raw value (points/binary)
	maxValueInFramework := 0

	if reg := assessment.ComplianceRegister; reg != nil && reg.ComplianceAssessmentFramework != nil {
		if ft := reg.ComplianceAssessmentFramework.AssessmentFrameworkType; ft != nil {
			frameworkType = ft.Type
		}
		frameworkValues = assessment.ComplianceRegister.ComplianceAssessmentFramework.Values
		for _, v := range frameworkValues {
			if v != nil && v.ID > 0 && v.Value != nil {
				frameworkValuesMap[v.ID] = *v.Value
				if *v.Value > maxValueInFramework {
					maxValueInFramework = *v.Value
				}
			}
		}
	}

	// Per-control max for denominator
	perControlMax := func() float64 {
		switch frameworkType {
		case "percentage":
			return 100.0
		case "binary":
			return 1.0
		case "point", "maturity":
			if maxValueInFramework > 0 {
				return float64(maxValueInFramework)
			}
			// sensible fallbacks
			if frameworkType == "point" {
				return 2.0
			}
			return 5.0
		default:
			return 100.0
		}
	}()
	// ---- DEBUG toggle ----
	// const DEBUG = true
	// debugf := func(format string, args ...any) {
	// 	if DEBUG {
	// 		fmt.Printf(format, args...)
	// 	}
	// }
	// Map % (0–100) -> framework band/status (from ScoreMin/ScoreMax ranges).
	// IMPORTANT: For SAMA point/maturity, ScoreMin/ScoreMax are already percentages.
	mapPctToBand := func(pct float64) (fvID *uint, statusID *uint) {
		if assessment.ComplianceRegister == nil || assessment.ComplianceRegister.ComplianceAssessmentFramework == nil {
			return nil, nil
		}
		for _, fv := range frameworkValues {
			if fv == nil || fv.ScoreMin == nil || fv.ScoreMax == nil {
				continue
			}
			min := float64(*fv.ScoreMin) // already %
			max := float64(*fv.ScoreMax) // already %

			// debugf("[mapPctToBand] pct=%.2f  fv=%d  band=[%.2f, %.2f]  status=%d\n",
			// 	pct, fv.ID, min, max, fv.ComplianceStatusID)

			if pct >= min && (pct < max || (math.Abs(pct-max) < 1e-9 && math.Abs(max-100.0) < 1e-9)) {
				id := fv.ID
				st := fv.ComplianceStatusID
				return &id, &st
			}
		}
		// debugf("[mapPctToBand] NO MATCH for pct=%.2f (framework=%s)\n", pct, frameworkType)
		return nil, nil
	}

	// ---------------- Applicability + group controls by clause ----------------
	notApplicable, err := cras.complianceStatusSvc.GetByType("not_applicable")
	if err != nil {
		return nil, err
	}

	type clauseStat struct {
		total         int // total controls under clause
		applicable    int // CAFV.IsApplicable = true
		notApplicable int // CAFV.IsApplicable = false (explicit)
	}

	clauseControls := make(map[uint][]*models.ComplianceAssessmentControl) // applicable only (for scoring)
	clauseStats := make(map[uint]*clauseStat)                              // applicability counters

	for _, ac := range assessment.AssessmentControls {
		if ac == nil || ac.ComplianceControl == nil || ac.ComplianceControl.ComplianceClauseID == nil {
			continue
		}
		clauseID := *ac.ComplianceControl.ComplianceClauseID

		if clauseStats[clauseID] == nil {
			clauseStats[clauseID] = &clauseStat{}
		}
		cs := clauseStats[clauseID]
		cs.total++

		// Determine applicability
		isApplicable := false
		if ac.ComplianceAssessmentFrameworkValue != nil {
			isApplicable = ac.ComplianceAssessmentFrameworkValue.IsApplicable
		}
		if isApplicable {
			cs.applicable++
			clauseControls[clauseID] = append(clauseControls[clauseID], ac)
		} else {
			if ac.ComplianceAssessmentFrameworkValue != nil {
				cs.notApplicable++
			}
		}

		// Sync control metadata + control's control status
		var statusID *uint
		if ac.ComplianceAssessmentFrameworkValue != nil {
			sid := ac.ComplianceAssessmentFrameworkValue.ComplianceStatusID
			statusID = &sid
		}
		dto := complianceDTOs.ControlSyncUpdateDTO{
			ControlID:                            ac.ComplianceControl.ControlID,
			ComplianceControlID:                  ac.ComplianceControlID,
			ComplianceAssessmentFrameworkValueID: ac.ComplianceAssessmentFrameworkValueID,
			CompletedByUserID:                    ac.CompletedByUserID,
			CompletedAt:                          ac.CompletedAt,
			ComplianceStatusID:                   statusID,
		}
		if err := cras.repo.UpdateControlComplianceStatus(dto, tx...); err != nil {
			return nil, err
		}
		if err := cras.repo.UpdateComplianceControlMeta(dto, tx...); err != nil {
			return nil, err
		}
	}

	// return signature satisfaction (you don't need per-control DTOs)
	var results []*complianceDTOs.ControlComplianceScoreDTO

	// ---------------- Compute per-clause score & update ----------------
	var (
		clauseCount   int     // number of clauses that produced a numeric score
		sumClauseAvgs float64 // sum of clause scores used in overall
		totalClauses  int
		naClauses     int
	)

	for clauseID, stat := range clauseStats {
		totalClauses++

		// All controls in clause are NA → mark clause NA and skip from overall
		if stat.total > 0 && stat.applicable == 0 && stat.notApplicable == stat.total {
			// debugf("[CLAUSE %d] all controls NA → set status=not_applicable, score=0\n", clauseID)
			if err := cras.repo.UpdateClauseComplianceStatus(clauseID, notApplicable.ID, 0, tx...); err != nil {
				return nil, err
			}
			naClauses++
			continue
		}

		// Only average applicable controls
		ctrls := clauseControls[clauseID]
		if len(ctrls) == 0 {
			// debugf("[CLAUSE %d] no applicable controls → skip averaging\n", clauseID)
			continue // neither score nor status change here
		}

		var totalScore, maxPossible float64
		for _, ac := range ctrls {
			switch frameworkType {
			case "point", "maturity":
				maxPossible += perControlMax
				if ac.ComplianceAssessmentFrameworkValueID != nil {
					if val, ok := frameworkValuesMap[*ac.ComplianceAssessmentFrameworkValueID]; ok {
						totalScore += float64(val)
						// vStr = fmt.Sprintf("%d", val)
					}
				}
				// debugf("  [ctrl #%d] point val=%s  maxPossible+=%.0f  totals: score=%.2f max=%.2f\n",
				// 	ac.ID, vStr, perControlMax, totalScore, maxPossible)

			case "percentage":
				if ac.Score != nil {
					totalScore += float64(*ac.Score)
					maxPossible += 100.0
				}
			case "binary":
				maxPossible += 1.0
				if ac.ComplianceAssessmentFrameworkValueID != nil {
					if raw, ok := frameworkValuesMap[*ac.ComplianceAssessmentFrameworkValueID]; ok && raw > 0 {
						totalScore += 1.0
					}
				}
			default:
				if ac.Score != nil {
					totalScore += float64(*ac.Score)
					maxPossible += 100.0
				}
			}
		}

		if maxPossible <= 0 {
			// debugf("[CLAUSE %d] maxPossible=0 → skip\n", clauseID)
			continue
		}

		clausePct := (totalScore / maxPossible) * 100.0
		clauseScore := math.Round(clausePct*100) / 100
		// debugf("[CLAUSE %d] totalScore=%.2f max=%.2f → pct=%.2f score=%.2f\n",
		// 	clauseID, totalScore, maxPossible, clausePct, clauseScore)

		if _, stID := mapPctToBand(clausePct); stID != nil {
			// debugf("[CLAUSE %d] mapped statusID=%d → updating DB\n", clauseID, *stID)
			if err := cras.repo.UpdateClauseComplianceStatus(clauseID, *stID, clauseScore, tx...); err != nil {
				// debugf("[CLAUSE %d] update ERROR: %v\n", clauseID, err)
				return nil, err
			}
		}

		sumClauseAvgs += clauseScore
		clauseCount++
	}

	// ---------------- Overall (average of clause averages) ----------------
	if clauseCount > 0 {
		overall := math.Round((sumClauseAvgs/float64(clauseCount))*100) / 100
		if _, stID := mapPctToBand(overall); stID != nil {
			// debugf("[OVERALL] score=%.2f statusID=%d → updating register\n", overall, *stID)
			if err := cras.repo.UpdateComplianceStatus(assessment.ComplianceRegisterID, *stID, overall, tx...); err != nil {
				// debugf("[OVERALL] update ERROR: %v\n", err)
				return nil, err
			}
		}
	} else if totalClauses > 0 && naClauses == totalClauses {
		// debugf("[OVERALL] all clauses NA → set register NA, score=0\n")
		if err := cras.repo.UpdateComplianceStatus(assessment.ComplianceRegisterID, notApplicable.ID, 0, tx...); err != nil {
			return nil, err
		}
	}

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

func (cras ComplianceAssessmentService) GetControlApprovalStatus(id uint, tx ...utility.DatabaseTransaction) (*complianceDTOs.ComplianceAssessmentApprovalStatusDTO, error) {
	model, err := cras.repo.GetFullDetailsByID(id, tx...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	isApproved := true
	if isApproved { // only check if requirements already passed
		for _, ctrl := range model.AssessmentControls {
			if ctrl.ComplianceAssessmentRequirementStatus == nil ||
				ctrl.ComplianceAssessmentRequirementStatus.Type != "approved" {
				isApproved = false
				break
			}
		}
	}

	return &complianceDTOs.ComplianceAssessmentApprovalStatusDTO{
		AssessmentID: id,
		IsApproved:   isApproved,
	}, nil
}

func (cras *ComplianceAssessmentService) computeRequirementSnapshot(assessmentID uint, responsibleUnitIDs ...uint) (int, float64, error) {
	assessment, err := cras.repo.GetFullDetailsByIDOptimized(assessmentID)
	if err != nil {
		return 0, 0, err
	}
	if assessment == nil || assessment.ComplianceRegister == nil {
		return 0, 0, fmt.Errorf("assessment/register not found")
	}
	register := assessment.ComplianceRegister

	// ---------- RU filter (unchanged) ----------
	ruFilter := make(map[uint]struct{}, len(responsibleUnitIDs))
	for _, id := range responsibleUnitIDs {
		ruFilter[id] = struct{}{}
	}
	hasRUFilter := len(ruFilter) > 0

	// ---------- Framework setup (same rules) ----------
	frameworkType := ""
	fvMap := make(map[uint]int) // frameworkValueID -> numeric value
	maxFV := 0

	if f := register.ComplianceAssessmentFramework; f != nil {
		if f.AssessmentFrameworkType != nil {
			frameworkType = f.AssessmentFrameworkType.Type
		}
		for _, v := range f.Values {
			if v == nil || v.Value == nil {
				continue
			}
			fvMap[uint(v.ID)] = *v.Value
			if *v.Value > maxFV {
				maxFV = *v.Value
			}
		}
	}

	perReqMax := func() float64 {
		switch frameworkType {
		case "percentage":
			return 100.0
		case "binary":
			return 1.0
		case "point", "maturity":
			if maxFV > 0 {
				return float64(maxFV)
			}
			if frameworkType == "point" {
				return 2.0
			}
			return 5.0
		default:
			return 100.0
		}
	}()

	// only count MaxPossible when RU has an actual value/score (same rule)
	ruContributesMax := func(ru *models.ComplianceAssessmentRequirementResponsibleUnit) bool {
		isApplicable := true
		if ru.ComplianceAssessmentFrameworkValue != nil {
			isApplicable = ru.ComplianceAssessmentFrameworkValue.IsApplicable
		}
		if !isApplicable {
			return false
		}
		switch frameworkType {
		case "point", "maturity", "binary":
			return ru.ComplianceAssessmentFrameworkValueID != nil
		case "percentage":
			return ru.Score != nil
		default:
			return ru.Score != nil
		}
	}

	round2 := func(x float64) float64 { return math.Round(x*100) / 100 }

	// ---------- (A) per-requirement snapshot with RU filter ----------
	type reqScore struct {
		score   float64 // 0..100 (rounded 2dp)
		counted bool    // counted in clause averages
		clause  uint
	}
	reqScores := make(map[uint]reqScore) // requirementID -> reqScore

	for _, r := range assessment.AssessmentRequirements {
		if r == nil || r.ID == 0 || r.ComplianceRequirement == nil {
			continue
		}

		var s, m float64
		var anyApplicable bool
		na := 0
		seenAfterFilter := 0

		for _, ru := range r.ResponsibleUnits {
			if ru == nil {
				continue
			}
			// RU filter (unchanged)
			if hasRUFilter {
				if _, ok := ruFilter[ru.ResponsibleUnitID]; !ok {
					continue
				}
			}
			seenAfterFilter++

			isApplicable := true
			if ru.ComplianceAssessmentFrameworkValue != nil {
				isApplicable = ru.ComplianceAssessmentFrameworkValue.IsApplicable
			}
			if !isApplicable {
				na++
				continue
			}
			anyApplicable = true

			if ruContributesMax(ru) {
				m += perReqMax
			}

			switch frameworkType {
			case "point", "maturity":
				if ru.ComplianceAssessmentFrameworkValueID != nil {
					if val, ok := fvMap[*ru.ComplianceAssessmentFrameworkValueID]; ok {
						s += float64(val)
					}
				}
			case "percentage":
				if ru.Score != nil {
					s += float64(*ru.Score)
				}
			case "binary":
				if ru.ComplianceAssessmentFrameworkValueID != nil {
					if raw, ok := fvMap[*ru.ComplianceAssessmentFrameworkValueID]; ok && raw > 0 {
						s += 1.0
					}
				}
			default:
				if ru.Score != nil {
					s += float64(*ru.Score)
				}
			}
		}

		if seenAfterFilter == 0 {
			continue
		}

		out := reqScore{
			clause: r.ComplianceRequirement.ComplianceClause.ID,
		}

		if !anyApplicable && na > 0 {
			out.counted = false
			reqScores[r.ID] = out
			continue
		}
		if m == 0 {
			out.counted = false
			reqScores[r.ID] = out
			continue
		}

		pct := (s / m) * 100.0
		out.score = round2(pct)
		out.counted = true
		reqScores[r.ID] = out
	}

	// ---------- (B) build clause tree ----------
	type cAgg struct {
		selfSum  float64 // sum of direct req scores under this clause
		selfCnt  int     // number of counted direct reqs
		final    float64 // final score after roll-up
		hasFinal bool
	}

	agg := map[uint]*cAgg{}           // clauseID -> aggregates
	parentOf := map[uint]*uint{}      // childID -> parentID
	children := map[uint][]uint{}     // parentID -> []childIDs
	seenClause := map[uint]struct{}{} // all clauses present

	for _, r := range assessment.AssessmentRequirements {
		if r == nil || r.ComplianceRequirement == nil {
			continue
		}
		cl := r.ComplianceRequirement.ComplianceClause
		if cl.ID == 0 {
			continue
		}
		cid := cl.ID
		seenClause[cid] = struct{}{}
		if agg[cid] == nil {
			agg[cid] = &cAgg{}
		}
		if cl.ParentClauseID != nil {
			parentOf[cid] = cl.ParentClauseID
			children[*cl.ParentClauseID] = append(children[*cl.ParentClauseID], cid)
			seenClause[*cl.ParentClauseID] = struct{}{}
			if agg[*cl.ParentClauseID] == nil {
				agg[*cl.ParentClauseID] = &cAgg{}
			}
		} else {
			parentOf[cid] = nil
		}
	}

	// seed SELF from reqScores (direct reqs only)
	for _, rs := range reqScores {
		if !rs.counted || rs.clause == 0 {
			continue
		}
		a := agg[rs.clause]
		a.selfSum += rs.score
		a.selfCnt++
	}

	// ---------- (C) post-order roll-up with SIMPLE child avg + self avg ----------
	roots := make([]uint, 0)
	for cid := range seenClause {
		if parentOf[cid] == nil {
			roots = append(roots, cid)
		}
	}

	visited := map[uint]bool{}
	var postOrder func(uint)
	postOrder = func(cid uint) {
		if visited[cid] {
			return
		}
		visited[cid] = true

		// compute children first
		childSum, childCnt := 0.0, 0
		for _, ch := range children[cid] {
			postOrder(ch)
			if agg[ch] != nil && agg[ch].hasFinal {
				childSum += agg[ch].final
				childCnt++
			}
		}

		var (
			hasSelf, hasChild bool
			selfAvg, childAvg float64
		)
		if a := agg[cid]; a != nil && a.selfCnt > 0 {
			selfAvg = round2(a.selfSum / float64(a.selfCnt))
			hasSelf = true
		}
		if childCnt > 0 {
			childAvg = round2(childSum / float64(childCnt)) // simple average of child finals
			hasChild = true
		}

		switch {
		case hasSelf && hasChild:
			agg[cid].final = round2((selfAvg + childAvg) / 2.0)
			agg[cid].hasFinal = true
		case hasSelf:
			agg[cid].final = selfAvg
			agg[cid].hasFinal = true
		case hasChild:
			agg[cid].final = childAvg
			agg[cid].hasFinal = true
		default:
			agg[cid].hasFinal = false
		}
	}

	for _, root := range roots {
		postOrder(root)
	}

	// ---------- (D) overall = SIMPLE average of root finals ----------
	var overallSum float64
	var overallCnt int
	for _, root := range roots {
		if a := agg[root]; a != nil && a.hasFinal {
			overallSum += a.final
			overallCnt++
		}
	}
	overall := 0.0
	if overallCnt > 0 {
		overall = round2(overallSum / float64(overallCnt))
	}

	// ---------- total applicable requirements (post-filter) ----------
	totalApplicable := 0
	for _, rs := range reqScores {
		if rs.counted {
			totalApplicable++
		}
	}

	return totalApplicable, overall, nil
}

func (cras *ComplianceAssessmentService) computeControlSnapshot(assessmentID uint) (total int, avg float64, err error) {
	// Prefer the optimized load so we can use the same DTO shape you used in CalculateComplianceControlScores
	assessment, err := cras.GetFullDetailsByIDOptimized(assessmentID)
	if err != nil {
		return 0, 0, err
	}
	if assessment == nil || assessment.ComplianceRegister == nil {
		return 0, 0, fmt.Errorf("assessment/register not found")
	}

	// ---------- Framework setup (same as scoring logic) ----------
	frameworkType := ""
	fvMap := make(map[uint]int)
	maxFV := 0

	var frameworkValues []*complianceDTOs.ComplianceAssessmentFrameworkValueDetailDTO
	if reg := assessment.ComplianceRegister; reg != nil && reg.ComplianceAssessmentFramework != nil {
		if t := reg.ComplianceAssessmentFramework.AssessmentFrameworkType; t != nil {
			frameworkType = t.Type
		}
		frameworkValues = reg.ComplianceAssessmentFramework.Values
		for _, v := range frameworkValues {
			if v == nil || v.Value == nil {
				continue
			}
			fvMap[uint(v.ID)] = *v.Value
			if *v.Value > maxFV {
				maxFV = *v.Value
			}
		}
	}

	perControlMax := func() float64 {
		switch frameworkType {
		case "percentage":
			return 100.0
		case "binary":
			return 1.0
		case "point", "maturity":
			if maxFV > 0 {
				return float64(maxFV)
			}
			if frameworkType == "point" {
				return 2.0
			}
			return 5.0
		default:
			return 100.0
		}
	}()

	round2 := func(x float64) float64 { return math.Round(x*100) / 100 }

	// Count MaxPossible only when control is applicable AND has a value/score
	ctrlContributesMax := func(ac *complianceDTOs.ComplianceAssessmentControlDTO) bool {
		if ac == nil {
			return false
		}
		isApplicable := true
		if ac.ComplianceAssessmentFrameworkValue != nil && ac.ComplianceAssessmentFrameworkValue.IsApplicable != nil {
			isApplicable = *ac.ComplianceAssessmentFrameworkValue.IsApplicable
		}
		if !isApplicable {
			return false
		}
		switch frameworkType {
		case "point", "maturity", "binary":
			return ac.ComplianceAssessmentFrameworkValueID != nil
		case "percentage":
			return ac.Score != nil
		default:
			return ac.Score != nil
		}
	}

	// ---------- (A) Self (direct) clause averages ----------
	type selfAgg struct {
		scoreSum       float64
		maxPossible    float64
		seenApplicable bool
		anySeen        bool
	}
	selfByClause := map[uint]*selfAgg{}
	totalApplicable := 0 // keep your old meaning: count applicable controls (regardless of whether they carry a value)

	for _, ac := range assessment.AssessmentControls {
		if ac == nil {
			continue
		}
		cid := ac.ComplianceControl.ComplianceClauseID
		if cid == 0 {
			continue
		}
		if selfByClause[cid] == nil {
			selfByClause[cid] = &selfAgg{}
		}
		a := selfByClause[cid]
		a.anySeen = true

		// applicability
		isApplicable := false
		if ac.ComplianceAssessmentFrameworkValue != nil && ac.ComplianceAssessmentFrameworkValue.IsApplicable != nil {
			isApplicable = *ac.ComplianceAssessmentFrameworkValue.IsApplicable
		}
		if !isApplicable {
			continue
		}
		a.seenApplicable = true
		totalApplicable++ // unchanged semantics

		// denominator only if the control carries a value/score
		if ctrlContributesMax(ac) {
			a.maxPossible += perControlMax
		}

		// numerator based on framework type
		switch frameworkType {
		case "point", "maturity":
			if ac.ComplianceAssessmentFrameworkValueID != nil {
				if val, ok := fvMap[*ac.ComplianceAssessmentFrameworkValueID]; ok {
					a.scoreSum += float64(val)
				}
			}
		case "percentage":
			if ac.Score != nil {
				a.scoreSum += float64(*ac.Score)
			}
		case "binary":
			if ac.ComplianceAssessmentFrameworkValueID != nil {
				if raw, ok := fvMap[*ac.ComplianceAssessmentFrameworkValueID]; ok && raw > 0 {
					a.scoreSum += 1.0
				}
			}
		default:
			if ac.Score != nil {
				a.scoreSum += float64(*ac.Score)
			}
		}
	}

	// ---------- (B) Build clause tree & seed self averages ----------
	type cAgg struct {
		selfAvg *float64
		final   *float64
	}
	agg := map[uint]*cAgg{}
	parentOf := map[uint]*uint{}
	children := map[uint][]uint{}
	seen := map[uint]struct{}{}

	for _, ac := range assessment.AssessmentControls {
		if ac == nil {
			continue
		}
		cl := ac.ComplianceControl.ComplianceClause
		cid := cl.ID
		if cid == 0 {
			continue
		}
		seen[cid] = struct{}{}
		if agg[cid] == nil {
			agg[cid] = &cAgg{}
		}
		if cl.ParentClauseID != nil {
			parentOf[cid] = cl.ParentClauseID
			children[*cl.ParentClauseID] = append(children[*cl.ParentClauseID], cid)
			seen[*cl.ParentClauseID] = struct{}{}
			if agg[*cl.ParentClauseID] == nil {
				agg[*cl.ParentClauseID] = &cAgg{}
			}
		}
	}

	for cid, s := range selfByClause {
		if s == nil {
			continue
		}
		if s.seenApplicable && s.maxPossible > 0 {
			avg := round2((s.scoreSum / s.maxPossible) * 100.0)
			if agg[cid] == nil {
				agg[cid] = &cAgg{}
			}
			agg[cid].selfAvg = &avg
		}
	}

	// ---------- (C) Bottom-up roll-up: final = avg(selfAvg, childAvg) ----------
	roots := make([]uint, 0)
	for cid := range seen {
		if parentOf[cid] == nil {
			roots = append(roots, cid)
		}
	}

	visited := map[uint]bool{}
	var postOrder func(uint)
	postOrder = func(cid uint) {
		if visited[cid] {
			return
		}
		visited[cid] = true

		childSum := 0.0
		childCnt := 0
		for _, ch := range children[cid] {
			postOrder(ch)
			if agg[ch] != nil && agg[ch].final != nil {
				childSum += *agg[ch].final
				childCnt++
			}
		}

		var (
			hasSelf  bool
			selfAvg  float64
			hasChild bool
			childAvg float64
		)
		if a := agg[cid]; a != nil && a.selfAvg != nil {
			selfAvg = *a.selfAvg
			hasSelf = true
		}
		if childCnt > 0 {
			childAvg = round2(childSum / float64(childCnt))
			hasChild = true
		}

		switch {
		case hasSelf && hasChild:
			v := round2((selfAvg + childAvg) / 2.0)
			agg[cid].final = &v
		case hasSelf:
			v := selfAvg
			agg[cid].final = &v
		case hasChild:
			v := childAvg
			agg[cid].final = &v
		default:
			agg[cid].final = nil
		}
	}
	for _, r := range roots {
		postOrder(r)
	}

	// ---------- (D) Overall = simple average of root finals ----------
	var overallSum float64
	var overallCnt int
	for _, root := range roots {
		if a := agg[root]; a != nil && a.final != nil {
			overallSum += *a.final
			overallCnt++
		}
	}
	if overallCnt == 0 {
		return totalApplicable, 0, nil
	}
	overall := round2(overallSum / float64(overallCnt))
	return totalApplicable, overall, nil
}

func parseUintList(ss []string) []uint {
	out := make([]uint, 0, len(ss))
	for _, s := range ss {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		if v, err := strconv.ParseUint(s, 10, 32); err == nil {
			out = append(out, uint(v))
		}
	}
	return out
}

func (cas *ComplianceAssessmentService) ExportComplianceAssessments(data []*complianceDTOs.ComplianceAssessmentListDTO, currentUserID uint) (*os.File, error) {
	userSetting, err := cas.userSettingSvc.GetUserSettingsByUserID(currentUserID)
	if err != nil {
		return nil, err
	}

	var headers []string
	if userSetting.Language == "ar" {
		headers = []string{
			"SL No.",
			"Compliance Register",
			"Review User",
			"Description",
			"Status",
			"Created By",
		}
	} else {
		headers = []string{
			"SL No.",
			"Compliance Register",
			"Review User",
			"Description",
			"Status",
			"Created By",
		}
	}

	rows := make([][]interface{}, len(data))
	for i, assessment := range data {
		var complianceRegister, reviewUser, description, status, createdBy string

		if userSetting.Language == "ar" {
			complianceRegister = assessment.ComplianceRegisterTitleAr
			reviewUser = assessment.ReviewUserFirstNameAr + " " + assessment.ReviewUserLastNameAr
			status = assessment.ComplianceAssessmentStatusTitleAr
			createdBy = assessment.CreatedUserFirstNameAr + " " + assessment.CreatedUserLastNameAr
			description = utility.GetStringValue(assessment.Description)
		} else {
			complianceRegister = assessment.ComplianceRegisterTitleEn
			reviewUser = assessment.ReviewUserFirstNameEn + " " + assessment.ReviewUserLastNameEn
			status = assessment.ComplianceAssessmentStatusTitleEn
			createdBy = assessment.CreatedUserFirstNameEn + " " + assessment.CreatedUserLastNameEn
			description = utility.GetStringValue(assessment.Description)

		}

		rows[i] = []interface{}{
			i + 1,
			complianceRegister,
			reviewUser,
			description,
			status,
			createdBy,
		}
	}

	return utility.GenerateExcelFile(headers, rows)
}

func derefUintPtr(p *uint) uint {
	if p == nil {
		return 0
	}
	return *p
}

func (cras *ComplianceAssessmentService) getResidualRiskRatingByScore(score float64) (*uint, error) {
	// Get all compliance assessment ratings
	ratings, err := cras.repo.GetComplianceAssessmentRatings()
	if err != nil {
		return nil, err
	}
	fmt.Printf("[SERVICE] Total ratings found: %d\n", len(ratings))

	// Find the appropriate rating based on score range
	for _, rating := range ratings {
		if rating.ScoreFrom != nil && rating.ScoreTo != nil {
			fmt.Printf("[SERVICE] Checking rating ID %d: ScoreFrom=%v, ScoreTo=%v, Score=%v\n", rating.ID, *rating.ScoreFrom, *rating.ScoreTo, score)
			if score >= *rating.ScoreFrom && score <= *rating.ScoreTo {
				fmt.Printf("[SERVICE] Match found! Rating ID: %d\n", rating.ID)
				return &rating.ID, nil
			}
		}
	}

	fmt.Printf("[SERVICE] No rating found for score: %v\n", score)
	return nil, nil // No rating found
}
