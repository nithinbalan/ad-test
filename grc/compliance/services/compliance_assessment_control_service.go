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
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type ComplianceAssessmentControlServiceInterface interface {
	GetAllClauses(assessmentID uint, currentUserID uint) ([]*complianceDTOs.ComplianceAssessmentControlClausesListDTO, error)
	GetAll(assessmentID uint, params *complianceDTOs.ComplianceAssessmentControlListQueryParamsDTO) (*complianceDTOs.ComplianceAssessmentControlListResponseDTO, error)
	GetByID(controlID uint) (*models.ComplianceAssessmentControl, error)
	GetFullDetailsByID(assessmentId uint, controlId uint, tx ...utility.DatabaseTransaction) (*complianceDTOs.ComplianceAssessmentControlDetailDTO, error)
	Update(complianceAssessmentId uint, controlId uint, currentUserID uint, updateForm *complianceDTOs.ComplianceAssessmentControlUpdateFormDTO) error

	MarkAsCompleted(complianceAssessmentId uint, controlId uint, currentUserID uint) error
	MarkAsOpen(complianceAssessmentId uint, controlId uint, currentUserID uint) error

	ApproveComplianceControl(assessmentID uint, controlID uint, currentUserID uint, commentForm *complianceDTOs.ComplianceAssessmentControlCommentFormDTO) error
	RevertComplianceControl(assessmentID uint, controlID uint, currentUserID uint, commentForm *complianceDTOs.ComplianceAssessmentControlCommentFormDTO) error
	ResubmitComplianceControl(assessmentID uint, controlID uint, currentUserID uint, commentForm *complianceDTOs.ComplianceAssessmentControlCommentFormDTO) error
}

type ComplianceAssessmentControlService struct {
	repo                             complianceRepositories.ComplianceAssessmentControlRepositoryInterface
	activityLogSvc                   logServices.ActivityLogServiceInterface
	userSvc                          organizationServices.UserServiceInterface
	complianceAssessmentSvc          ComplianceAssessmentServiceInterface
	controlStatusSvc                 ComplianceAssessmentRequirementStatusServiceInterface
	notificationSvc                  commonServices.NotificationServiceInterface
	complianceAssessmentFrameworkSvc ComplianceAssessmentFrameworkServiceInterface
	complianceAssessmentStatusSvc    ComplianceAssessmentStatusServiceInterface
}

func NewComplianceAssessmentControlService(
	repo complianceRepositories.ComplianceAssessmentControlRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	complianceAssessmentSvc ComplianceAssessmentServiceInterface,
	controlStatusSvc ComplianceAssessmentRequirementStatusServiceInterface,
	notificationSvc commonServices.NotificationServiceInterface,
	complianceAssessmentFrameworkSvc ComplianceAssessmentFrameworkServiceInterface,
	complianceAssessmentStatusSvc ComplianceAssessmentStatusServiceInterface,
) ComplianceAssessmentControlServiceInterface {
	return &ComplianceAssessmentControlService{
		repo:                             repo,
		activityLogSvc:                   activityLogSvc,
		userSvc:                          userSvc,
		complianceAssessmentSvc:          complianceAssessmentSvc,
		controlStatusSvc:                 controlStatusSvc,
		notificationSvc:                  notificationSvc,
		complianceAssessmentFrameworkSvc: complianceAssessmentFrameworkSvc,
		complianceAssessmentStatusSvc:    complianceAssessmentStatusSvc,
	}
}

func (s *ComplianceAssessmentControlService) GetAllClauses(assessmentID uint, currentUserID uint) ([]*complianceDTOs.ComplianceAssessmentControlClausesListDTO, error) {
	clauses, err := s.repo.GetAllClauses(assessmentID, currentUserID)
	if err != nil {
		return nil, err
	}
	if len(clauses) == 0 {
		return []*complianceDTOs.ComplianceAssessmentControlClausesListDTO{}, nil
	}
	countsMap, err := s.repo.GetClauseControlCounts(assessmentID)
	if err != nil {
		return nil, err
	}
	return s.buildHierarchicalClauseStructure(clauses, countsMap)
}

func (s *ComplianceAssessmentControlService) buildHierarchicalClauseStructure(clauses []*models.ComplianceClause, countMap map[uint]complianceDTOs.ClauseControlCountDTO) ([]*complianceDTOs.ComplianceAssessmentControlClausesListDTO, error) {
	// Build a proper hierarchy by filtering out child clauses from the root level
	var rootClauses []*complianceDTOs.ComplianceAssessmentControlClausesListDTO
	clauseMap := make(map[uint]*complianceDTOs.ComplianceAssessmentControlClausesListDTO)

	// Convert all clauses into DTOs and map them by ID
	for _, clause := range clauses {
		dto := &complianceDTOs.ComplianceAssessmentControlClausesListDTO{
			ID:                clause.ID,
			ClauseNo:          clause.ClauseNo,
			TitleEn:           clause.TitleEn,
			TitleAr:           clause.TitleAr,
			ChildClauses:      []*complianceDTOs.ComplianceAssessmentControlClausesListDTO{},
			TotalControls:     0,
			CompletedControls: 0,
		}
		if counts, ok := countMap[clause.ID]; ok {
			dto.TotalControls = counts.TotalControls
			dto.CompletedControls = counts.CompletedControls
		}
		clauseMap[clause.ID] = dto
	}

	// Build hierarchy by linking children to their parent
	for _, clause := range clauses {
		if clause.ParentClauseID != nil {
			parent := clauseMap[*clause.ParentClauseID]
			parent.ChildClauses = append(parent.ChildClauses, clauseMap[clause.ID])
		} else {
			rootClauses = append(rootClauses, clauseMap[clause.ID])
		}
	}

	return rootClauses, nil
}
func (s *ComplianceAssessmentControlService) GetAll(assessmentID uint, params *complianceDTOs.ComplianceAssessmentControlListQueryParamsDTO) (*complianceDTOs.ComplianceAssessmentControlListResponseDTO, error) {
	if _, err := s.complianceAssessmentSvc.GetByID(assessmentID); err != nil {
		return nil, err
	}
	return s.repo.GetAll(assessmentID, params)
}
func (s *ComplianceAssessmentControlService) GetByID(id uint) (*models.ComplianceAssessmentControl, error) {
	control, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return control, nil
}
func (s *ComplianceAssessmentControlService) GetFullDetailsByID(assessmentId uint, controlId uint, tx ...utility.DatabaseTransaction) (*complianceDTOs.ComplianceAssessmentControlDetailDTO, error) {
	var dto complianceDTOs.ComplianceAssessmentControlDetailDTO

	control, err := s.repo.GetFullDetailsByID(controlId, tx...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, controlId)
		}
		return nil, err
	}

	err = mapstructure.Decode(control, &dto)
	if err != nil {
		return nil, err
	}
	// ❗ Prevent duplicates: mapstructure already filled ActionPlans.
	// Clear it, then repopulate with TargetDateString set.
	dto.ActionPlans = nil

	seen := make(map[uint]struct{}, len(control.ActionPlans))
	for _, ap := range control.ActionPlans {
		if ap == nil {
			continue
		}
		if _, ok := seen[ap.ID]; ok {
			continue
		}
		seen[ap.ID] = struct{}{}

		var apDTO complianceDTOs.ComplianceActionPlanDetailDTO
		if err := mapstructure.Decode(ap, &apDTO); err != nil {
			return nil, err
		}

		if ap.TargetDate != nil {
			td := ap.TargetDate.Format("2006-01-02")
			apDTO.TargetDateString = &td
		}

		dto.ActionPlans = append(dto.ActionPlans, &apDTO)
	}
	// --- ControlComments: ensure CreatedAt & Comment are set; prevent duplicates ---
	dto.ControlComments = nil
	seenCmt := make(map[uint]struct{}, len(control.ControlComments))
	for _, c := range control.ControlComments {
		if c == nil {
			continue
		}
		if _, ok := seenCmt[c.ID]; ok {
			continue
		}
		seenCmt[c.ID] = struct{}{}

		var cDTO complianceDTOs.ComplianceAssessmentControlCommentDetailDTO
		if err := mapstructure.Decode(c, &cDTO); err != nil {
			return nil, err
		}
		cDTO.CreatedAt = c.CreatedAt
		cDTO.Comment = c.Comment

		dto.ControlComments = append(dto.ControlComments, &cDTO)
	}

	dto.CreatedAt = control.CreatedAt
	dto.UpdatedAt = control.UpdatedAt
	return &dto, nil
}
func (s *ComplianceAssessmentControlService) Update(complianceAssessmentId uint, controlId uint, currentUserID uint, updateForm *complianceDTOs.ComplianceAssessmentControlUpdateFormDTO) error {
	control, err := s.GetByID(controlId)
	if err != nil {
		return err
	}
	complianceAssessment, err := s.complianceAssessmentSvc.GetByID(complianceAssessmentId)
	if err != nil {
		return err
	}
	assessmentDetails, err := s.complianceAssessmentSvc.GetFullDetailsByID(complianceAssessmentId)
	if err != nil {
		return err
	}
	assessmentInProgressStatus, err := s.complianceAssessmentStatusSvc.GetByType("in_progress")
	if err != nil {
		return err
	}

	assessmentInitiatedStatus, err := s.complianceAssessmentStatusSvc.GetByType("draft")
	if err != nil {
		return err
	}
	complianceAssessmentFrameworkValue, err := s.complianceAssessmentFrameworkSvc.GetByFrameworkValueID(updateForm.AssessmentFrameworkValueID)
	if err != nil {
		return err
	}

	tx := s.repo.BeginTransaction()

	if complianceAssessment.ComplianceAssessmentStatusID == assessmentInitiatedStatus.ID {
		complianceAssessment.ComplianceAssessmentStatusID = assessmentInProgressStatus.ID
		err = s.repo.ComplianceAssessmentUpdateStatus(complianceAssessment, tx)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	if complianceAssessmentFrameworkValue.IsApplicable && assessmentDetails.ComplianceRegister.ComplianceAssessmentFramework.AssessmentFrameworkType.Type == "percentage" {
		control.Score = updateForm.Score
	} else {
		control.Score = nil
	}
	control.Comments = updateForm.Comments
	control.ComplianceAssessmentFrameworkValueID = &updateForm.AssessmentFrameworkValueID
	control.UpdatedBYID = &currentUserID

	err = s.repo.Update(control, tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	// async logging
	go func() {
		changedData := map[string]interface{}{
			"score":   updateForm.Score,
			"comment": updateForm.Comments,
		}
		_ = s.activityLogSvc.LogActivity(models.UpdateAction, models.ComplianceAssessmentControl{}.TableName(), controlId, changedData, currentUserID, nil)
	}()

	return nil
}
func (s *ComplianceAssessmentControlService) MarkAsCompleted(complianceAssessmentId uint, controlId uint, currentUserID uint) error {
	tx := s.repo.BeginTransaction()

	control, err := s.repo.GetByID(controlId)
	if err != nil {
		return err
	}
	if control.ComplianceAssessmentFrameworkValueID == nil {
		return utility.NewForbiddenError("score_not_updated")
	}
	now := time.Now()
	control.IsCompleted = true
	control.CompletedByUserID = &currentUserID
	control.CompletedAt = &now
	control.UpdatedBYID = &currentUserID

	if err := s.repo.Update(control, tx); err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}
func (s *ComplianceAssessmentControlService) MarkAsOpen(complianceAssessmentId uint, controlId uint, currentUserID uint) error {
	tx := s.repo.BeginTransaction()

	control, err := s.repo.GetByID(controlId)
	if err != nil {
		return err
	}

	control.IsCompleted = false
	control.UpdatedBYID = &currentUserID

	if err := s.repo.Update(control, tx); err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}
func (s *ComplianceAssessmentControlService) ApproveComplianceControl(assessmentID uint, controlID uint, currentUserID uint, commentForm *complianceDTOs.ComplianceAssessmentControlCommentFormDTO) error {
	tx := s.repo.BeginTransaction()

	approvedStatus, err := s.controlStatusSvc.GetByType("approved")
	if err != nil {
		tx.Rollback()
		return err
	}

	control, err := s.repo.GetByID(controlID)
	if err != nil {
		tx.Rollback()
		return err
	}

	if control.ComplianceAssessmentRequirementStatusID != nil && *control.ComplianceAssessmentRequirementStatusID == approvedStatus.ID {
		tx.Rollback()
		return utility.NewValidationError(fmt.Errorf("control already approved"), "invalid_state", "control")
	}
	comment := &models.ComplianceAssessmentControlComment{
		ComplianceAssessmentControlID:           controlID,
		Comment:                                 commentForm.Comment,
		CreatedBYID:                             currentUserID,
		ComplianceAssessmentRequirementStatusID: &approvedStatus.ID,
	}

	err = s.repo.CreateComment(comment, tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	control.ComplianceAssessmentRequirementStatusID = &approvedStatus.ID
	control.UpdatedBYID = &currentUserID

	if err := s.repo.Update(control, tx); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	go func() {
		// log activity
		_ = s.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.ComplianceAssessmentControl{}.TableName(),
			controlID,
			map[string]interface{}{"status": "approved", "comment": commentForm.Comment},
			currentUserID,
			nil,
		)

		assessment, _ := s.complianceAssessmentSvc.GetFullDetailsByID(assessmentID)
		createdBy, _ := s.userSvc.GetByID(currentUserID)

		if control.CompletedByUserID != nil {
			completedUser, _ := s.userSvc.GetByID(*control.CompletedByUserID)

			urlParam := fmt.Sprintf(`{"assessment_id":%d, "control_id":%d}`, assessmentID, controlID)
			entity := "compliance-assessment-control"
			action := "approved"

			// DB Notification
			dbNotification := &commonDTOs.DBNotificationDTO{
				UserIDs:              []uint{completedUser.ID},
				ModuleIdentifier:     "compliance_register",
				NotificationType:     string(commonDTOs.AssignmentType),
				HeadingEn:            "Compliance Control Approved",
				HeadingAr:            "تمت الموافقة على الضبط",
				TitleEn:              assessment.ComplianceRegister.TitleEn,
				TitleAr:              assessment.ComplianceRegister.TitleAr,
				ReferenceNumber:      assessment.ComplianceRegister.ReferenceNumber,
				Entity:               &entity,
				Action:               &action,
				URLParam:             urlParam,
				IsSystemNotification: false,
				CreatedByID:          &currentUserID,
			}
			appURL := utility.GetAppURL()

			// Email Notification
			emailNotification := &commonDTOs.EmailNotificationDTO{
				EmailTo:              []string{completedUser.Email},
				Subject:              fmt.Sprintf("Compliance Control #%d Approved", controlID),
				TemplateRelativePath: "compliance/templates/email/compliance_control_approved.html",
				EmailData: map[string]interface{}{
					"FirstNameEn":     completedUser.FirstNameEn,
					"LastNameEn":      completedUser.LastNameEn,
					"TitleEn":         assessment.ComplianceRegister.TitleEn,
					"TitleAr":         assessment.ComplianceRegister.TitleAr,
					"ReferenceNumber": assessment.ComplianceRegister.ReferenceNumber,
					"Comment":         commentForm.Comment,
					"ApprovedBy":      createdBy.FirstNameEn,
					"ApprovalDate":    time.Now().Format("2006-01-02"),
					"RedirectURL":     fmt.Sprintf("%s/compliance/assessments/%d/overview", appURL, assessmentID),
				},
			}

			_ = s.notificationSvc.SendNotification(dbNotification, emailNotification, "compliance_assessment")
		}
	}()

	return nil
}
func (s *ComplianceAssessmentControlService) RevertComplianceControl(assessmentID uint, controlID uint, currentUserID uint, commentForm *complianceDTOs.ComplianceAssessmentControlCommentFormDTO) error {
	tx := s.repo.BeginTransaction()

	revertedStatus, err := s.controlStatusSvc.GetByType("reverted")
	if err != nil {
		return err
	}

	control, err := s.repo.GetByID(controlID)
	if err != nil {
		return err
	}

	// prevent duplicate revert
	if control.ComplianceAssessmentRequirementStatusID != nil &&
		*control.ComplianceAssessmentRequirementStatusID == revertedStatus.ID {
		return utility.NewValidationError(fmt.Errorf("control is already reverted"), "invalid_state", "control")
	}

	// authorization: only compliance register responsible user can revert
	assessment, err := s.complianceAssessmentSvc.GetFullDetailsByID(assessmentID)
	if err != nil {
		return err
	}
	if currentUserID != assessment.ComplianceRegister.ResponsibleUserID {
		return utility.NewForbiddenError("User is not authorized to revert this control")
	}

	// add a comment
	comment := &models.ComplianceAssessmentControlComment{
		ComplianceAssessmentControlID:           control.ID,
		Comment:                                 commentForm.Comment,
		CreatedBYID:                             currentUserID,
		ComplianceAssessmentRequirementStatusID: &revertedStatus.ID,
		CreatedAt:                               time.Now(),
	}
	if err := s.repo.CreateComment(comment, tx); err != nil {
		tx.Rollback()
		return err
	}

	// update control status
	control.ComplianceAssessmentRequirementStatusID = &revertedStatus.ID
	control.UpdatedBYID = &currentUserID
	if err := s.repo.Update(control, tx); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	// async log + notification
	go func() {
		_ = s.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.ComplianceAssessmentControl{}.TableName(),
			controlID,
			map[string]interface{}{"status": "reverted", "comment": commentForm.Comment},
			currentUserID,
			nil,
		)

		createdBy, _ := s.userSvc.GetByID(currentUserID)
		if control.CompletedByUserID != nil {
			completedUser, _ := s.userSvc.GetByID(*control.CompletedByUserID)

			urlParam := fmt.Sprintf(`{"assessment_id":%d, "control_id":%d}`, assessmentID, controlID)
			entity := "compliance-assessment-control"
			action := "reverted"

			dbNotification := &commonDTOs.DBNotificationDTO{
				UserIDs:              []uint{completedUser.ID},
				ModuleIdentifier:     "compliance_register",
				NotificationType:     string(commonDTOs.AssignmentType),
				HeadingEn:            "Compliance Control Reverted",
				HeadingAr:            "تم إرجاع الضبط",
				TitleEn:              assessment.ComplianceRegister.TitleEn,
				TitleAr:              assessment.ComplianceRegister.TitleAr,
				ReferenceNumber:      assessment.ComplianceRegister.ReferenceNumber,
				Entity:               &entity,
				Action:               &action,
				URLParam:             urlParam,
				IsSystemNotification: false,
				CreatedByID:          &currentUserID,
			}
			appURL := utility.GetAppURL()

			emailNotification := &commonDTOs.EmailNotificationDTO{
				EmailTo:              []string{completedUser.Email},
				Subject:              fmt.Sprintf("Compliance Control #%d Reverted", controlID),
				TemplateRelativePath: "compliance/templates/email/compliance_control_reverted.html",

				EmailData: map[string]interface{}{
					"FirstNameEn":     completedUser.FirstNameEn,
					"LastNameEn":      completedUser.LastNameEn,
					"TitleEn":         assessment.ComplianceRegister.TitleEn,
					"TitleAr":         assessment.ComplianceRegister.TitleAr,
					"ReferenceNumber": assessment.ComplianceRegister.ReferenceNumber,
					"Comment":         commentForm.Comment,
					"RevertedBy":      createdBy.FirstNameEn,
					"ReversionDate":   time.Now().Format("2006-01-02"),
					"RedirectURL":     fmt.Sprintf("%s/compliance/assessments/%d/overview", appURL, assessmentID),
				},
			}
			_ = s.notificationSvc.SendNotification(dbNotification, emailNotification, "compliance_assessment")
		}
	}()
	return nil
}
func (s *ComplianceAssessmentControlService) ResubmitComplianceControl(assessmentID uint, controlID uint, currentUserID uint, commentForm *complianceDTOs.ComplianceAssessmentControlCommentFormDTO) error {
	tx := s.repo.BeginTransaction()

	resubmittedStatus, err := s.controlStatusSvc.GetByType("resubmitted")
	if err != nil {
		tx.Rollback()
		return err
	}

	control, err := s.repo.GetByID(controlID)
	if err != nil {
		tx.Rollback()
		return err
	}

	if control.ComplianceAssessmentRequirementStatusID != nil &&
		*control.ComplianceAssessmentRequirementStatusID == resubmittedStatus.ID {
		tx.Rollback()
		return utility.NewValidationError(fmt.Errorf("control is already resubmitted"), "invalid_state", "control")
	}

	// add comment
	comment := &models.ComplianceAssessmentControlComment{
		ComplianceAssessmentControlID:           control.ID,
		Comment:                                 commentForm.Comment,
		CreatedBYID:                             currentUserID,
		ComplianceAssessmentRequirementStatusID: &resubmittedStatus.ID,
	}
	if err := s.repo.CreateComment(comment, tx); err != nil {
		tx.Rollback()
		return err
	}

	// update control status
	control.ComplianceAssessmentRequirementStatusID = &resubmittedStatus.ID
	control.UpdatedBYID = &currentUserID
	if err := s.repo.Update(control, tx); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	// async log + notify responsible user
	go func() {
		_ = s.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.ComplianceAssessmentControl{}.TableName(),
			controlID,
			map[string]interface{}{"status": "resubmitted", "comment": commentForm.Comment},
			currentUserID,
			nil,
		)

		assessmentDetails, _ := s.complianceAssessmentSvc.GetFullDetailsByID(assessmentID)
		createdBy, _ := s.userSvc.GetByID(currentUserID)
		responsibleUser, _ := s.userSvc.GetByID(assessmentDetails.ComplianceRegister.ResponsibleUserID)

		urlParam := fmt.Sprintf(`{"assessment_id":%d, "control_id":%d}`, assessmentID, controlID)
		entity := "compliance-assessment-control"
		action := "resubmitted"

		dbNotification := &commonDTOs.DBNotificationDTO{
			UserIDs:              []uint{responsibleUser.ID},
			ModuleIdentifier:     "compliance_register",
			NotificationType:     string(commonDTOs.AssignmentType),
			HeadingEn:            "Compliance Control Resubmitted",
			HeadingAr:            "إعادة تقديم الضبط",
			TitleEn:              assessmentDetails.ComplianceRegister.TitleEn,
			TitleAr:              assessmentDetails.ComplianceRegister.TitleAr,
			ReferenceNumber:      assessmentDetails.ComplianceRegister.ReferenceNumber,
			Entity:               &entity,
			Action:               &action,
			URLParam:             urlParam,
			IsSystemNotification: false,
			CreatedByID:          &currentUserID,
		}
		appURL := utility.GetAppURL()

		emailNotification := &commonDTOs.EmailNotificationDTO{
			EmailTo:              []string{responsibleUser.Email},
			Subject:              fmt.Sprintf("Compliance Control #%d Resubmitted", controlID),
			TemplateRelativePath: "compliance/templates/email/compliance_control_resubmitted.html",
			EmailData: map[string]interface{}{
				"FirstNameEn":      responsibleUser.FirstNameEn,
				"LastNameEn":       responsibleUser.LastNameEn,
				"TitleEn":          assessmentDetails.ComplianceRegister.TitleEn,
				"TitleAr":          assessmentDetails.ComplianceRegister.TitleAr,
				"ReferenceNumber":  assessmentDetails.ComplianceRegister.ReferenceNumber,
				"Comment":          commentForm.Comment,
				"ResubmittedBy":    createdBy.FirstNameEn,
				"ResubmissionDate": time.Now().Format("2006-01-02"),
				"RedirectURL":      fmt.Sprintf("%s/compliance/assessments/%d/overview", appURL, assessmentID),
			},
		}
		_ = s.notificationSvc.SendNotification(dbNotification, emailNotification, "compliance_assessment")
	}()
	return nil
}
