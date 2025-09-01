package services

import (
	"errors"
	"fmt"
	"time"

	commonDTOs "jethur/grc/common/dtos"
	commonServices "jethur/grc/common/services"
	complianceDTOs "jethur/grc/compliance/dtos"
	"jethur/grc/compliance/repositories"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type ComplianceAssessmentControlCommentServiceInterface interface {
	GetAll(assessmentID uint, controlID uint) ([]*complianceDTOs.ComplianceAssessmentControlCommentListDTO, error)
	GetByID(assessmentID uint, controlID uint, id uint) (*models.ComplianceAssessmentControlComment, error)
	GetFullDetailsByID(assessmentID uint, controlID uint, id uint, tx ...utility.DatabaseTransaction) (*complianceDTOs.ComplianceAssessmentControlCommentDetailDTO, error)
	Create(assessmentID uint, controlID uint, currentUserID uint, dto *complianceDTOs.ComplianceAssessmentControlCommentFormDTO) (uint, error)
	Update(assessmentID uint, controlID uint, id uint, currentUserID uint, dto *complianceDTOs.ComplianceAssessmentControlCommentFormDTO) error
	Delete(assessmentID uint, controlID uint, ids []uint, currentUserID uint) ([]uint, error)
}

type ComplianceAssessmentControlCommentService struct {
	repo                           repositories.ComplianceAssessmentControlCommentRepositoryInterface
	userSvc                        organizationServices.UserServiceInterface
	activityLogSvc                 logServices.ActivityLogServiceInterface
	complianceAssessmentControlSvc ComplianceAssessmentControlServiceInterface
	assessmentSvc                  ComplianceAssessmentServiceInterface
	notificationSvc                commonServices.NotificationServiceInterface
	statusSvc                      ComplianceAssessmentStatusServiceInterface
}

func NewComplianceAssessmentControlCommentService(
	repo repositories.ComplianceAssessmentControlCommentRepositoryInterface,
	userSvc organizationServices.UserServiceInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	complianceAssessmentControlSvc ComplianceAssessmentControlServiceInterface,
	assessmentSvc ComplianceAssessmentServiceInterface,
	notificationSvc commonServices.NotificationServiceInterface,
	statusSvc ComplianceAssessmentStatusServiceInterface,
) ComplianceAssessmentControlCommentServiceInterface {
	return &ComplianceAssessmentControlCommentService{
		repo:                           repo,
		userSvc:                        userSvc,
		activityLogSvc:                 activityLogSvc,
		complianceAssessmentControlSvc: complianceAssessmentControlSvc,
		assessmentSvc:                  assessmentSvc,
		notificationSvc:                notificationSvc,
		statusSvc:                      statusSvc,
	}
}

func (s *ComplianceAssessmentControlCommentService) GetAll(assessmentID uint, controlID uint) ([]*complianceDTOs.ComplianceAssessmentControlCommentListDTO, error) {
	if _, err := s.complianceAssessmentControlSvc.GetByID(controlID); err != nil {
		return nil, err
	}
	return s.repo.GetAll(controlID)
}

func (s *ComplianceAssessmentControlCommentService) GetByID(assessmentID uint, controlID uint, id uint) (*models.ComplianceAssessmentControlComment, error) {
	if _, err := s.complianceAssessmentControlSvc.GetByID(controlID); err != nil {
		return nil, err
	}
	comment, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return comment, nil
}

func (s *ComplianceAssessmentControlCommentService) GetFullDetailsByID(assessmentID uint, controlID uint, id uint, tx ...utility.DatabaseTransaction) (*complianceDTOs.ComplianceAssessmentControlCommentDetailDTO, error) {
	if _, err := s.complianceAssessmentControlSvc.GetByID(controlID); err != nil {
		return nil, err
	}

	var dto complianceDTOs.ComplianceAssessmentControlCommentDetailDTO
	comment, err := s.repo.GetFullDetailsByID(id, tx...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	err = mapstructure.Decode(comment, &dto)
	if err != nil {
		return nil, err
	}
	dto.CreatedAt = comment.CreatedAt
	return &dto, nil
}

func (s *ComplianceAssessmentControlCommentService) Create(assessmentID uint, controlID uint, currentUserID uint, dto *complianceDTOs.ComplianceAssessmentControlCommentFormDTO) (uint, error) {

	control, err := s.complianceAssessmentControlSvc.GetFullDetailsByID(assessmentID, controlID)
	if err != nil {
		return 0, err
	}

	assessment, err := s.assessmentSvc.GetFullDetailsByID(assessmentID)
	if err != nil {
		return 0, err
	}
	tx := s.repo.BeginTransaction()

	comment := &models.ComplianceAssessmentControlComment{
		ComplianceAssessmentControlID: controlID,
		Comment:                       dto.Comment,
		CreatedBYID:                   currentUserID,
		CreatedAt:                     time.Now(),
	}
	id, err := s.repo.Create(comment, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	tx.Commit()

	go func() {
		detail, _ := s.GetFullDetailsByID(assessmentID, controlID, id)
		_ = s.activityLogSvc.LogActivity(
			models.CreateAction,
			models.ComplianceAssessmentControlComment{}.TableName(),
			id,
			detail,
			currentUserID,
			nil,
		)

		// ---- Notification logic ----
		createdBy, _ := s.userSvc.GetByID(currentUserID)

		// decide receiver: if reviewer commented → notify submitter, if submitter commented → notify reviewer
		var notifyUserID uint
		// if assessment.ReviewUserID == currentUserID {
		// 	notifyUserID = uint(assessment.SubmittedByUser.ID)
		// } else {
		// 	notifyUserID = uint(assessment.ReviewUser.ID)
		// }
		notifyUserID = uint(assessment.ReviewUser.ID)

		notifyUser, err := s.userSvc.GetByID(notifyUserID)
		if err != nil {
			return
		}

		urlParam := fmt.Sprintf(`{"assessment_id":%d,"control_id":%d}`, assessmentID, controlID)
		entity := "compliance-assessment-control"
		action := "comment-added"

		dbNotification := &commonDTOs.DBNotificationDTO{
			UserIDs:              []uint{notifyUserID},
			ModuleIdentifier:     "compliance_register",
			NotificationType:     string(commonDTOs.AssignmentType),
			HeadingEn:            "Comment Added For Compliance Control",
			HeadingAr:            "تمت إضافة تعليق لضبط الامتثال",
			TitleEn:              assessment.ComplianceRegister.TitleEn,
			TitleAr:              assessment.ComplianceRegister.TitleAr,
			ReferenceNumber:      assessment.ComplianceRegister.ReferenceNumber,
			Entity:               &entity,
			Action:               &action,
			URLParam:             urlParam,
			IsSystemNotification: false,
			CreatedByID:          &currentUserID,
		}

		emailNotification := &commonDTOs.EmailNotificationDTO{
			EmailTo:              []string{notifyUser.Email},
			Subject:              fmt.Sprintf("New Comment on Compliance Control #%d", controlID),
			TemplateRelativePath: "compliance/templates/email/compliance_assessment_control_comment.html",
			EmailData: map[string]interface{}{
				"FirstNameEn":    notifyUser.FirstNameEn,
				"LastNameEn":     notifyUser.LastNameEn,
				"ControlTitle":   control.ComplianceControl.Control.TitleEn,
				"Comment":        dto.Comment,
				"CommentedBy":    createdBy.FirstNameEn,
				"CommentedOn":    time.Now().Format("2006-01-02"),
				"ReferenceNo":    assessment.ComplianceRegister.ReferenceNumber,
				"AssessmentName": assessment.ComplianceRegister.TitleEn,
				"RedirectURL":    fmt.Sprintf("%s/compliance/assessments/%d/controls/%d", utility.GetAppURL(), assessmentID, controlID),
			},
		}

		_ = s.notificationSvc.SendNotification(dbNotification, emailNotification, "compliance_assessment_control")
	}()

	return id, nil
}

func (s *ComplianceAssessmentControlCommentService) Update(assessmentID uint, controlID uint, id uint, currentUserID uint, dto *complianceDTOs.ComplianceAssessmentControlCommentFormDTO) error {
	tx := s.repo.BeginTransaction()

	comment, err := s.GetByID(assessmentID, controlID, id)
	if err != nil {
		return err
	}

	oldDetail, _ := s.GetFullDetailsByID(assessmentID, controlID, id)

	comment.Comment = dto.Comment
	if err := s.repo.Update(comment, tx); err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	go func() {
		newDetail, _ := s.GetFullDetailsByID(assessmentID, controlID, id)
		changes := utility.CompareFieldChanges(
			map[string]interface{}{"comment": oldDetail.Comment},
			map[string]interface{}{"comment": newDetail.Comment},
		)
		if len(changes) > 0 {
			_ = s.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.ComplianceAssessmentControlComment{}.TableName(),
				id,
				changes,
				currentUserID,
				nil,
			)
		}
	}()
	return nil
}

func (s *ComplianceAssessmentControlCommentService) Delete(assessmentID uint, controlID uint, ids []uint, currentUserID uint) ([]uint, error) {
	tx := s.repo.BeginTransaction()
	if _, err := s.complianceAssessmentControlSvc.GetByID(controlID); err != nil {
		return nil, err
	}

	var failedIDs, deletedIDs []uint
	var deletedComments []map[string]interface{}

	for _, id := range ids {
		if _, err := s.repo.GetFullDetailsByID(id); err != nil {
			failedIDs = append(failedIDs, id)
			continue
		}
		if err := s.repo.Delete(id, tx); err != nil {
			failedIDs = append(failedIDs, id)
			continue
		}
		deletedIDs = append(deletedIDs, id)
		deletedComments = append(deletedComments, map[string]interface{}{
			"id":         id,
			"control_id": controlID,
		})
	}

	if len(failedIDs) > 0 {
		tx.Rollback()
		return failedIDs, nil
	}

	tx.Commit()
	if len(deletedIDs) > 0 {
		go func() {
			_ = s.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.ComplianceAssessmentControlComment{}.TableName(),
				deletedIDs,
				deletedComments,
				currentUserID,
				nil,
			)
		}()
	}
	return failedIDs, nil
}
