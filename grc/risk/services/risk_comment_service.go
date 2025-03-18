package services

import (
	"errors"
	"fmt"
	"time"

	commonDTOs "jethur/grc/common/dtos"
	commonServices "jethur/grc/common/services"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	"jethur/grc/risk/dtos"
	"jethur/grc/risk/repositories"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type RiskCommentServiceInterface interface {
	GetAll(riskID uint) ([]*dtos.RiskCommentListDTO, error)
	GetByID(riskID uint, id uint) (*models.RiskComment, error)
	GetFullDetailsByID(riskID uint, id uint, tx ...utility.DatabaseTransaction) (*dtos.RiskCommentDetailDTO, error)
	Create(riskID uint, currentUserID uint, dto *dtos.RiskCommentFormDTO) (uint, error)
	Update(riskID uint, id uint, currentUserID uint, dto *dtos.RiskCommentFormDTO) error
	Delete(riskID uint, ids []uint, currentUserID uint) ([]uint, error)

	SubmitRisk(riskID uint, currentUserID uint) error
	RevertRisk(riskID uint, currentUserID uint, commentForm *dtos.RiskCommentFormDTO) error
	RejectRisk(riskID uint, currentUserID uint, commentForm *dtos.RiskCommentFormDTO) error
	ApproveRisk(riskID uint, currentUserID uint, commentForm *dtos.RiskCommentFormDTO) error
	ResubmitRisk(riskID uint, currentUserID uint, commentForm *dtos.RiskCommentFormDTO) error
}

type RiskCommentService struct {
	repo            repositories.RiskCommentRepositoryInterface
	riskRegisterSvc RiskRegisterServiceInterface
	riskStatusSvc   RiskStatusServiceInterface
	activityLogSvc  logServices.ActivityLogServiceInterface
	userSvc         organizationServices.UserServiceInterface
	riskSettingsSvc RiskSettingsServiceInterface
	notificationSvc commonServices.NotificationServiceInterface
}

func NewRiskCommentService(repo repositories.RiskCommentRepositoryInterface,
	riskRegisterSvc RiskRegisterServiceInterface,
	riskStatusSvc RiskStatusServiceInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	riskSettingsSvc RiskSettingsServiceInterface,
	notificationSvc commonServices.NotificationServiceInterface,
) *RiskCommentService {
	return &RiskCommentService{
		repo:            repo,
		riskRegisterSvc: riskRegisterSvc,
		riskStatusSvc:   riskStatusSvc,
		activityLogSvc:  activityLogSvc,
		userSvc:         userSvc,
		riskSettingsSvc: riskSettingsSvc,
		notificationSvc: notificationSvc,
	}
}

func (rcs *RiskCommentService) GetAll(riskID uint) ([]*dtos.RiskCommentListDTO, error) {
	return rcs.repo.GetAll(riskID)
}

func (rcs *RiskCommentService) GetByID(riskID uint, id uint) (*models.RiskComment, error) {
	comment, err := rcs.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return comment, nil
}

func (rcs *RiskCommentService) GetFullDetailsByID(riskID uint, id uint, tx ...utility.DatabaseTransaction) (*dtos.RiskCommentDetailDTO, error) {
	var commentDetail dtos.RiskCommentDetailDTO

	comment, err := rcs.repo.GetFullDetailsByID(id, tx...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	err = mapstructure.Decode(comment, &commentDetail)
	if err != nil {
		return nil, err
	}

	commentDetail.CreatedAt = comment.CreatedAt
	return &commentDetail, nil
}

func (rcs *RiskCommentService) Create(riskID uint, currentUserID uint, riskCommentFormDTO *dtos.RiskCommentFormDTO) (uint, error) {
	tx := rcs.repo.BeginTransaction()

	comment := &models.RiskComment{
		RiskRegisterID:     riskID,
		RiskTimelineTypeID: riskCommentFormDTO.RiskTimelineTypeID,
		Comments:           riskCommentFormDTO.Comment,
		CreatedBYID:        currentUserID,
	}

	id, err := rcs.repo.Create(comment)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	tx.Commit()
	go func() {
		existingComment, err := rcs.GetFullDetailsByID(riskID, id)
		if err != nil {
			utility.LogError("Failed to fetch risk comment", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		err = rcs.activityLogSvc.LogActivity(
			models.CreateAction,
			models.RiskComment{}.TableName(),
			id,
			existingComment,
			currentUserID,
			nil,
		)
		if err != nil {
			utility.LogError("Failed to log activity", map[string]interface{}{
				"error": err.Error(),
			})
		}
	}()

	return id, nil
}

func (rcs *RiskCommentService) Update(riskID uint, id uint, currentUserID uint, riskCommentFormDTO *dtos.RiskCommentFormDTO) error {
	tx := rcs.repo.BeginTransaction()

	existingComment, err := rcs.GetByID(riskID, id)
	if err != nil {
		return err
	}

	oldComment, err := rcs.GetFullDetailsByID(riskID, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	existingComment.Comments = riskCommentFormDTO.Comment
	existingComment.RiskTimelineTypeID = riskCommentFormDTO.RiskTimelineTypeID
	err = rcs.repo.Update(existingComment, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	go func() {
		existingComment, err := rcs.GetFullDetailsByID(riskID, id)
		if err != nil {
			return
		}

		changedData := utility.CompareFieldChanges(
			map[string]interface{}{"comment": oldComment.Comment},
			map[string]interface{}{"comment": existingComment.Comment},
		)

		if len(changedData) > 0 {
			err := rcs.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.RiskComment{}.TableName(),
				id,
				changedData,
				currentUserID,
				nil,
			)
			if err != nil {
				utility.LogError("Failed to log activity", map[string]interface{}{
					"error": err.Error(),
				})
			}
		}
	}()
	return nil
}

func (rcs *RiskCommentService) Delete(riskID uint, ids []uint, currentUserID uint) ([]uint, error) {
	tx := rcs.repo.BeginTransaction()

	var failedIDs []uint
	var successfulIDs []uint
	var deletedComments []map[string]interface{}

	for _, id := range ids {
		if _, err := rcs.repo.GetFullDetailsByID(id); err != nil {
			failedIDs = append(failedIDs, id)
			continue
		}

		deletedComments = append(deletedComments, map[string]interface{}{
			"id":      id,
			"risk_id": riskID,
		})

		err := rcs.repo.Delete(id, tx)
		if err != nil {
			failedIDs = append(failedIDs, id)
		} else {
			successfulIDs = append(successfulIDs, id)
		}
	}

	if len(failedIDs) > 0 {
		tx.Rollback()
		return failedIDs, nil
	}

	tx.Commit()
	if len(successfulIDs) > 0 {
		go func() {
			err := rcs.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.RiskComment{}.TableName(),
				successfulIDs,
				deletedComments,
				currentUserID,
				nil,
			)
			if err != nil {
				fmt.Printf("Failed to log bulk delete activity: %v\n", err)
			}
		}()
	}

	return failedIDs, nil
}

func (rcs *RiskCommentService) SubmitRisk(riskID uint, currentUserID uint) error {
	tx := rcs.repo.BeginTransaction()

	submittedStatus, err := rcs.riskStatusSvc.GetByType("submitted")
	if err != nil {
		tx.Rollback()
		return err
	}

	risk, err := rcs.riskRegisterSvc.GetByID(riskID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Ensure only the responsible user can submit
	if risk.ResponsibleUserID != currentUserID {
		tx.Rollback()
		return utility.NewForbiddenError("Only the responsible user can submit this risk.")
	}

	// Update risk status to "submitted"
	risk.RiskStatusID = submittedStatus.ID
	risk.SubmittedBy = &currentUserID
	now := time.Now()
	risk.SubmittedAt = &now
	err = rcs.repo.UpdateRiskStatus(risk, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	// Log the action
	go func() {
		err := rcs.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.RiskRegister{}.TableName(),
			riskID,
			map[string]interface{}{
				"status": "submitted",
				"user":   currentUserID,
			},
			currentUserID,
			nil,
		)
		if err != nil {
			utility.LogError("Failed to log activity", map[string]interface{}{
				"action":    models.UpdateAction,
				"tableName": models.RiskRegister{}.TableName(),
				"recordID":  riskID,
				"userID":    currentUserID,
				"error":     err.Error(),
			})
		}

		// Fetch risk settings
		riskSetting, err := rcs.riskSettingsSvc.Get()
		if err != nil {
			utility.LogError("Failed to fetch risk setting", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		// Get Review User IDs
		reviewUserIDs, err := rcs.userSvc.GetUserIdsByRole(riskSetting.ReviewRoleID)
		if err != nil {
			utility.LogError("Failed to fetch review user IDs", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		// Ensure we have valid reviewers
		if len(reviewUserIDs) == 0 {
			fmt.Println("No review users found, skipping notifications")
			return
		}

		// Prepare Review User Emails & Names
		reviewEmailTo := make([]string, len(reviewUserIDs))
		reviewers := make(map[uint]string) // Map of userID -> full name

		for i, userID := range reviewUserIDs {
			user, err := rcs.userSvc.GetByID(userID)
			if err != nil {
				utility.LogError("Failed to fetch reviewer details", map[string]interface{}{
					"error":  err.Error(),
					"userID": userID,
				})
				continue
			}
			reviewEmailTo[i] = user.Email
			reviewers[userID] = user.FirstNameEn + " " + user.LastNameEn
		}

		createdBy, err := rcs.userSvc.GetByID(currentUserID)
		if err != nil {
			utility.LogError("Failed to fetch created user", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		// Prepare notification and email
		urlParam := fmt.Sprintf(`{"riskID":%d}`, riskID)
		entity := "risk-register"
		action := "submitted"

		dbNotification := &commonDTOs.DBNotificationDTO{
			UserIDs:              reviewUserIDs,
			ModuleIdentifier:     "risk_register",
			NotificationType:     string(commonDTOs.AssignmentType),
			HeadingEn:            "Risk Submitted",
			HeadingAr:            "تم تقديم المخاطر",
			TitleEn:              risk.Title,
			TitleAr:              risk.Title,
			ReferenceNumber:      risk.ReferenceNumber,
			Entity:               &entity,
			Action:               &action,
			URLParam:             urlParam,
			IsSystemNotification: false,
			CreatedByID:          &currentUserID,
		}

		// **Send DB Notification only once**
		err = rcs.notificationSvc.SendNotification(dbNotification, nil, "risk")
		if err != nil {
			utility.LogError("Failed to send DB notification", map[string]interface{}{
				"error": err.Error(),
			})
		}

		for userID, fullName := range reviewers {
			emailNotification := &commonDTOs.EmailNotificationDTO{
				EmailTo:              reviewEmailTo,
				Subject:              fmt.Sprintf("Risk #%d Submitted", riskID),
				TemplateRelativePath: "risk/templates/email/risk_submitted.html",
				EmailData: map[string]interface{}{
					"FirstNameEn":          fullName,
					"TitleEn":              risk.Title,
					"Description":          risk.Description,
					"ReferenceNumber":      risk.ReferenceNumber,
					"SubmittedByFirstName": createdBy.FirstNameEn,
					"SubmittedByLastName":  createdBy.LastNameEn,
					"SubmissionDate":       time.Now().Format("2006-01-02"),
					"RedirectURL":          fmt.Sprintf("%srisk/risk-register/%d/overview", utility.GetAppURL(), riskID),
				},
			}

			err = rcs.notificationSvc.SendNotification(nil, emailNotification, "risk")
			if err != nil {
				utility.LogError("Failed to send review user email", map[string]interface{}{
					"error":  err.Error(),
					"userID": userID,
				})
			}
		}
	}()

	return nil
}

func (rcs *RiskCommentService) RevertRisk(riskID uint, currentUserID uint, commentForm *dtos.RiskCommentFormDTO) error {
	tx := rcs.repo.BeginTransaction()

	riskSetting, err := rcs.riskSettingsSvc.Get()
	if err != nil {
		return err
	}

	// Get Review User IDs
	reviewUserIDs, err := rcs.userSvc.GetUserIdsByRole(riskSetting.ReviewRoleID)
	if err != nil {
		return err
	}

	if !contains(reviewUserIDs, currentUserID) {
		return utility.NewForbiddenError("Only the review user can revert this risk.")
	}

	revertedStatus, err := rcs.riskStatusSvc.GetByType("reverted")
	if err != nil {
		tx.Rollback()
		return err
	}

	risk, err := rcs.riskRegisterSvc.GetByID(riskID)
	if err != nil {
		tx.Rollback()
		return err
	}

	comment := &models.RiskComment{
		RiskRegisterID:     riskID,
		RiskTimelineTypeID: commentForm.RiskTimelineTypeID,
		Comments:           commentForm.Comment,
		RiskStatusID:       &revertedStatus.ID,
		CreatedBYID:        currentUserID,
	}

	_, err = rcs.repo.Create(comment, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Update risk status to "reverted"
	risk.RiskStatusID = revertedStatus.ID
	risk.ReviewedBy = &currentUserID
	err = rcs.repo.UpdateRiskStatus(risk, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	// Log the action
	go func() {
		err := rcs.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.RiskRegister{}.TableName(),
			riskID,
			map[string]interface{}{
				"status":  "reverted",
				"comment": commentForm.Comment,
				"user":    currentUserID,
			},
			currentUserID,
			nil,
		)
		if err != nil {
			utility.LogError("Failed to log activity", map[string]interface{}{
				"action":    models.UpdateAction,
				"tableName": models.RiskRegister{}.TableName(),
				"recordID":  riskID,
				"userID":    currentUserID,
				"error":     err.Error(),
			})
		}

		responsibleUser, err := rcs.userSvc.GetByID(risk.ResponsibleUserID)
		if err != nil {
			utility.LogError("Failed to fetch responsible user", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		createdBy, err := rcs.userSvc.GetByID(currentUserID)
		if err != nil {
			utility.LogError("Failed to fetch created user", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		// Prepare notification and email
		urlParam := fmt.Sprintf(`{"riskID":%d}`, riskID)
		entity := "risk-register"
		action := "reverted"

		dbNotification := &commonDTOs.DBNotificationDTO{
			UserIDs:              []uint{responsibleUser.ID},
			ModuleIdentifier:     "risk_register",
			NotificationType:     string(commonDTOs.AssignmentType),
			HeadingEn:            "Risk Reverted",
			HeadingAr:            "عادت المخاطر",
			TitleEn:              risk.Title,
			TitleAr:              risk.Title,
			ReferenceNumber:      risk.ReferenceNumber,
			Entity:               &entity,
			Action:               &action,
			URLParam:             urlParam,
			IsSystemNotification: false,
			CreatedByID:          &currentUserID,
		}

		emailNotification := &commonDTOs.EmailNotificationDTO{
			EmailTo:              []string{responsibleUser.Email},
			Subject:              fmt.Sprintf("Risk #%d Reverted", riskID),
			TemplateRelativePath: "risk/templates/email/risk_reverted.html",
			EmailData: map[string]interface{}{
				"FirstNameEn":         responsibleUser.FirstNameEn,
				"LastNameEn":          responsibleUser.LastNameEn,
				"TitleEn":             risk.Title,
				"Description":         risk.Description,
				"ReferenceNumber":     risk.ReferenceNumber,
				"Comment":             commentForm.Comment,
				"RevertedByFirstName": createdBy.FirstNameEn,
				"RevertedByLastName":  createdBy.LastNameEn,
				"ReversionDate":       time.Now().Format("2006-01-02"),
				"RedirectURL":         fmt.Sprintf("%srisk/risk-register/%d/overview", utility.GetAppURL(), riskID),
			},
		}

		// **Send Both Notifications at ONCE**
		err = rcs.notificationSvc.SendNotification(dbNotification, emailNotification, "risk")
		if err != nil {
			utility.LogError("Failed to send notifications", map[string]interface{}{
				"error": err.Error(),
			})
		}
	}()

	return nil
}

func (rcs *RiskCommentService) RejectRisk(riskID uint, currentUserID uint, commentForm *dtos.RiskCommentFormDTO) error {
	tx := rcs.repo.BeginTransaction()

	riskSetting, err := rcs.riskSettingsSvc.Get()
	if err != nil {
		return err
	}

	// Get Review User IDs
	reviewUserIDs, err := rcs.userSvc.GetUserIdsByRole(riskSetting.ReviewRoleID)
	if err != nil {
		return err
	}

	if !contains(reviewUserIDs, currentUserID) {
		return utility.NewForbiddenError("Only the review user can reject this risk.")
	}

	rejectedStatus, err := rcs.riskStatusSvc.GetByType("rejected")
	if err != nil {
		tx.Rollback()
		return err
	}

	risk, err := rcs.riskRegisterSvc.GetByID(riskID)
	if err != nil {
		tx.Rollback()
		return err
	}

	comment := &models.RiskComment{
		RiskRegisterID:     riskID,
		RiskTimelineTypeID: commentForm.RiskTimelineTypeID,
		Comments:           commentForm.Comment,
		RiskStatusID:       &rejectedStatus.ID,
		CreatedBYID:        currentUserID,
	}

	_, err = rcs.repo.Create(comment, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Update risk status to "rejected"
	risk.RiskStatusID = rejectedStatus.ID
	risk.ReviewedBy = &currentUserID
	err = rcs.repo.UpdateRiskStatus(risk, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	// Log the action
	go func() {
		err := rcs.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.RiskRegister{}.TableName(),
			riskID,
			map[string]interface{}{
				"status":  "rejected",
				"comment": commentForm.Comment,
				"user":    currentUserID,
			},
			currentUserID,
			nil,
		)
		if err != nil {
			utility.LogError("Failed to log activity", map[string]interface{}{
				"action":    models.UpdateAction,
				"tableName": models.RiskRegister{}.TableName(),
				"recordID":  riskID,
				"userID":    currentUserID,
				"error":     err.Error(),
			})
		}

		responsibleUser, err := rcs.userSvc.GetByID(risk.ResponsibleUserID)
		if err != nil {
			utility.LogError("Failed to fetch responsible user", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		createdBy, err := rcs.userSvc.GetByID(currentUserID)
		if err != nil {
			utility.LogError("Failed to fetch created user", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		// Prepare notification and email
		urlParam := fmt.Sprintf(`{"riskID":%d}`, riskID)
		entity := "risk-register"
		action := "rejected"

		dbNotification := &commonDTOs.DBNotificationDTO{
			UserIDs:              []uint{responsibleUser.ID},
			ModuleIdentifier:     "risk_register",
			NotificationType:     string(commonDTOs.InformationType),
			HeadingEn:            "Risk Rejected",
			HeadingAr:            "تم رفض المخاطرة",
			TitleEn:              risk.Title,
			TitleAr:              risk.Title,
			ReferenceNumber:      risk.ReferenceNumber,
			Entity:               &entity,
			Action:               &action,
			URLParam:             urlParam,
			IsSystemNotification: false,
			CreatedByID:          &currentUserID,
		}

		emailNotification := &commonDTOs.EmailNotificationDTO{
			EmailTo:              []string{responsibleUser.Email},
			Subject:              fmt.Sprintf("Risk #%d Rejected", riskID),
			TemplateRelativePath: "risk/templates/email/risk_rejected.html",
			EmailData: map[string]interface{}{
				"FirstNameEn":         responsibleUser.FirstNameEn,
				"LastNameEn":          responsibleUser.LastNameEn,
				"TitleEn":             risk.Title,
				"Description":         risk.Description,
				"ReferenceNumber":     risk.ReferenceNumber,
				"Comment":             commentForm.Comment,
				"RejectedByFirstName": createdBy.FirstNameEn,
				"RejectedByLastName":  createdBy.LastNameEn,
				"RejectionDate":       time.Now().Format("2006-01-02"),
				"RedirectURL":         fmt.Sprintf("%srisk/risk-register/%d/overview", utility.GetAppURL(), riskID),
			},
		}

		// **Send Both Notifications at ONCE**
		err = rcs.notificationSvc.SendNotification(dbNotification, emailNotification, "risk")
		if err != nil {
			utility.LogError("Failed to send notifications", map[string]interface{}{
				"error": err.Error(),
			})
		}
	}()

	return nil
}

func (rcs *RiskCommentService) ApproveRisk(riskID uint, currentUserID uint, commentForm *dtos.RiskCommentFormDTO) error {
	tx := rcs.repo.BeginTransaction()

	riskSetting, err := rcs.riskSettingsSvc.Get()
	if err != nil {
		return err
	}

	// Get Review User IDs
	reviewUserIDs, err := rcs.userSvc.GetUserIdsByRole(riskSetting.ReviewRoleID)
	if err != nil {
		return err
	}

	if !contains(reviewUserIDs, currentUserID) {
		return utility.NewForbiddenError("Only the review user can approve this risk.")
	}

	approvedStatus, err := rcs.riskStatusSvc.GetByType("approved")
	if err != nil {
		tx.Rollback()
		return err
	}

	risk, err := rcs.riskRegisterSvc.GetByID(riskID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Validate if the current status allows approval
	if risk.RiskStatusID == approvedStatus.ID {
		tx.Rollback()
		return utility.NewValidationError(fmt.Errorf("risk is already approved"), "invalid_state", "risk")
	}

	comment := &models.RiskComment{
		RiskRegisterID:     riskID,
		RiskTimelineTypeID: commentForm.RiskTimelineTypeID,
		Comments:           commentForm.Comment,
		RiskStatusID:       &approvedStatus.ID,
		CreatedBYID:        currentUserID,
	}

	_, err = rcs.repo.Create(comment, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Update risk status to "approved"
	risk.RiskStatusID = approvedStatus.ID
	risk.ReviewedBy = &currentUserID
	err = rcs.repo.UpdateRiskStatus(risk, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	// Log the action
	go func() {
		err := rcs.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.RiskRegister{}.TableName(),
			riskID,
			map[string]interface{}{
				"status":  "approved",
				"comment": commentForm.Comment,
				"user":    currentUserID,
			},
			currentUserID,
			nil,
		)
		if err != nil {
			utility.LogError("Failed to log activity", map[string]interface{}{
				"action":    models.UpdateAction,
				"tableName": models.RiskRegister{}.TableName(),
				"recordID":  riskID,
				"userID":    currentUserID,
				"error":     err.Error(),
			})
		}

		responsibleUser, err := rcs.userSvc.GetByID(risk.ResponsibleUserID)
		if err != nil {
			utility.LogError("Failed to fetch responsible user", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		createdBy, err := rcs.userSvc.GetByID(currentUserID)
		if err != nil {
			utility.LogError("Failed to fetch created user", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		// Prepare notification and email
		urlParam := fmt.Sprintf(`{"riskID":%d}`, riskID)
		entity := "risk-register"
		action := "approved"

		dbNotificationDTO := &commonDTOs.DBNotificationDTO{
			UserIDs:              []uint{responsibleUser.ID},
			ModuleIdentifier:     "risk_register",
			NotificationType:     string(commonDTOs.InformationType),
			HeadingEn:            "Risk Approved",
			HeadingAr:            "تمت الموافقة على المخاطر",
			TitleEn:              risk.Title,
			TitleAr:              risk.Title,
			ReferenceNumber:      risk.ReferenceNumber,
			Entity:               &entity,
			Action:               &action,
			URLParam:             urlParam,
			IsSystemNotification: false,
			CreatedByID:          &currentUserID,
		}

		emailNotification := &commonDTOs.EmailNotificationDTO{
			EmailTo:              []string{responsibleUser.Email},
			Subject:              fmt.Sprintf("Risk #%d Approved", riskID),
			TemplateRelativePath: "risk/templates/email/risk_approved.html",
			EmailData: map[string]interface{}{
				"FirstNameEn":         responsibleUser.FirstNameEn,
				"LastNameEn":          responsibleUser.LastNameEn,
				"TitleEn":             risk.Title,
				"Description":         risk.Description,
				"ReferenceNumber":     risk.ReferenceNumber,
				"Comment":             commentForm.Comment,
				"ApprovedByFirstName": createdBy.FirstNameEn,
				"ApprovedByLastName":  createdBy.LastNameEn,
				"ApprovalDate":        time.Now().Format("2006-01-02"),
				"RedirectURL":         fmt.Sprintf("%srisk/risk-register/%d/overview", utility.GetAppURL(), riskID),
			},
		}

		// **Send Both Notifications at ONCE**
		err = rcs.notificationSvc.SendNotification(dbNotificationDTO, emailNotification, "risk")
		if err != nil {
			utility.LogError("Failed to send notifications", map[string]interface{}{
				"error": err.Error(),
			})
		}
	}()

	return nil
}

func (rcs *RiskCommentService) ResubmitRisk(riskID uint, currentUserID uint, commentForm *dtos.RiskCommentFormDTO) error {
	tx := rcs.repo.BeginTransaction()

	resubmittedStatus, err := rcs.riskStatusSvc.GetByType("resubmitted")
	if err != nil {
		tx.Rollback()
		return err
	}

	risk, err := rcs.riskRegisterSvc.GetByID(riskID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Ensure only the responsible user can submit
	if risk.ResponsibleUserID != currentUserID {
		tx.Rollback()
		return utility.NewForbiddenError("Only the responsible user can re-submit this risk.")
	}

	// Validate if the current status allows resubmission
	if risk.RiskStatusID == resubmittedStatus.ID {
		tx.Rollback()
		return utility.NewValidationError(fmt.Errorf("risk is already resubmitted"), "invalid_state", "risk")
	}

	comment := &models.RiskComment{
		RiskRegisterID:     riskID,
		RiskTimelineTypeID: commentForm.RiskTimelineTypeID,
		Comments:           commentForm.Comment,
		RiskStatusID:       &resubmittedStatus.ID,
		CreatedBYID:        currentUserID,
	}

	_, err = rcs.repo.Create(comment, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Update risk status to "resubmitted"
	risk.RiskStatusID = resubmittedStatus.ID
	err = rcs.repo.UpdateRiskStatus(risk, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	// Log the action
	go func() {
		err := rcs.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.RiskRegister{}.TableName(),
			riskID,
			map[string]interface{}{
				"status":  "resubmitted",
				"comment": commentForm.Comment,
				"user":    currentUserID,
			},
			currentUserID,
			nil,
		)
		if err != nil {
			utility.LogError("Failed to log activity", map[string]interface{}{
				"action":    models.UpdateAction,
				"tableName": models.RiskRegister{}.TableName(),
				"recordID":  riskID,
				"userID":    currentUserID,
				"error":     err.Error(),
			})
		}

		// Fetch risk settings
		riskSetting, err := rcs.riskSettingsSvc.Get()
		if err != nil {
			utility.LogError("Failed to fetch risk setting", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		// Get Review User IDs
		reviewUserIDs, err := rcs.userSvc.GetUserIdsByRole(riskSetting.ReviewRoleID)
		if err != nil {
			utility.LogError("Failed to fetch review user IDs", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		// Ensure we have valid reviewers
		if len(reviewUserIDs) == 0 {
			fmt.Println("No review users found, skipping notifications")
			return
		}

		// Prepare Review User Emails & Names
		reviewEmailTo := make([]string, len(reviewUserIDs))
		reviewers := make(map[uint]string)

		for i, userID := range reviewUserIDs {
			user, err := rcs.userSvc.GetByID(userID)
			if err != nil {
				utility.LogError("Failed to fetch reviewer details", map[string]interface{}{
					"error":  err.Error(),
					"userID": userID,
				})
				continue
			}
			reviewEmailTo[i] = user.Email
			reviewers[userID] = user.FirstNameEn + " " + user.LastNameEn
		}

		createdBy, err := rcs.userSvc.GetByID(currentUserID)
		if err != nil {
			utility.LogError("Failed to fetch created user", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		// Prepare notification and email
		urlParam := fmt.Sprintf(`{"riskID":%d}`, riskID)
		entity := "risk-register"
		action := "risk-resubmitted"

		dbNotification := &commonDTOs.DBNotificationDTO{
			UserIDs:              reviewUserIDs,
			ModuleIdentifier:     "risk_register",
			NotificationType:     string(commonDTOs.AssignmentType),
			HeadingEn:            "Risk Resubmitted",
			HeadingAr:            "تم إعادة تقديم المخاطر",
			TitleEn:              risk.Title,
			TitleAr:              risk.Title,
			ReferenceNumber:      risk.ReferenceNumber,
			Entity:               &entity,
			Action:               &action,
			URLParam:             urlParam,
			IsSystemNotification: false,
			CreatedByID:          &currentUserID,
		}

		err = rcs.notificationSvc.SendNotification(dbNotification, nil, "risk")
		if err != nil {
			utility.LogError("Failed to send DB notification", map[string]interface{}{
				"error": err.Error(),
			})
		}

		for userID, fullName := range reviewers {
			emailNotification := &commonDTOs.EmailNotificationDTO{
				EmailTo:              reviewEmailTo,
				Subject:              fmt.Sprintf("Risk #%d Resubmitted", riskID),
				TemplateRelativePath: "risk/templates/email/risk_resubmitted.html",
				EmailData: map[string]interface{}{
					"FirstNameEn":            fullName,
					"TitleEn":                risk.Title,
					"Description":            risk.Description,
					"ReferenceNumber":        risk.ReferenceNumber,
					"Comment":                commentForm.Comment,
					"ResubmittedByFirstName": createdBy.FirstNameEn,
					"ResubmittedByLastName":  createdBy.LastNameEn,
					"ResubmissionDate":       time.Now().Format("2006-01-02"),
					"RedirectURL":            fmt.Sprintf("%srisk/risk-register/%d/overview", utility.GetAppURL(), riskID),
				},
			}

			err = rcs.notificationSvc.SendNotification(nil, emailNotification, "risk")
			if err != nil {
				utility.LogError("Failed to send review user email", map[string]interface{}{
					"error":  err.Error(),
					"userID": userID,
				})
			}
		}
	}()

	return nil
}

func contains(slice []uint, value uint) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}
