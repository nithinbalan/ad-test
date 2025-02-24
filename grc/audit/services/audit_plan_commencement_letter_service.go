package services

import (
	"errors"
	"fmt"
	"html/template"
	"strings"
	"time"

	auditDTOs "jethur/grc/audit/dtos"
	"jethur/grc/audit/repositories"
	commonDTOs "jethur/grc/common/dtos"
	commonServices "jethur/grc/common/services"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type AuditPlanCommencementLetterServiceInterface interface {
	GetAll(auditPlanID uint, keyword string) ([]*auditDTOs.AuditPlanCommencementLetterListDTO, error)
	GetByID(auditPlanID uint, letterID uint) (*models.AuditPlanCommencementLetter, error)
	GetFullDetailsByID(auditPlanID uint, letterID uint, tx ...utility.DatabaseTransaction) (*auditDTOs.AuditPlanCommencementLetterDetailDTO, error)
	Create(auditPlanId uint, currentUserID uint) error
	Update(auditPlanID uint, letterID uint, currentUserID uint, commencementLetterFormDTO *auditDTOs.AuditPlanCommencementLetterFormDTO) error
	Delete(auditPlanID uint, letterID uint, currentUserID uint) (uint, error)
	Send(auditPlanID uint, letterID uint, currentUserID uint, commencementLetterSendFormDTO *auditDTOs.AuditPlanCommencementLetterSendFormDTO) error
}

type AuditPlanCommencementLetterService struct {
	repo                                 repositories.AuditPlanCommencementLetterRepositoryInterface
	activityLogSvc                       logServices.ActivityLogServiceInterface
	auditPlanSvc                         AuditPlanServiceInterface
	userSvc                              organizationServices.UserServiceInterface
	notificationSvc                      commonServices.NotificationServiceInterface
	auditPlanCommencementLetterStatusSvc AuditPlanCommencementLetterStatusServiceInterface
}

func NewAuditPlanCommencementLetterService(
	repo repositories.AuditPlanCommencementLetterRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	auditPlanSvc AuditPlanServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	notificationSvc commonServices.NotificationServiceInterface,
	auditPlanCommencementLetterStatusSvc AuditPlanCommencementLetterStatusServiceInterface,
) AuditPlanCommencementLetterServiceInterface {
	return &AuditPlanCommencementLetterService{
		repo:                                 repo,
		activityLogSvc:                       activityLogSvc,
		auditPlanSvc:                         auditPlanSvc,
		userSvc:                              userSvc,
		notificationSvc:                      notificationSvc,
		auditPlanCommencementLetterStatusSvc: auditPlanCommencementLetterStatusSvc,
	}
}

func (apcls *AuditPlanCommencementLetterService) GetAll(auditPlanID uint, keyword string) ([]*auditDTOs.AuditPlanCommencementLetterListDTO, error) {
	return apcls.repo.GetAll(auditPlanID, keyword)
}

func (apcls *AuditPlanCommencementLetterService) GetFullDetailsByID(auditPlanID uint, letterID uint, tx ...utility.DatabaseTransaction) (*auditDTOs.AuditPlanCommencementLetterDetailDTO, error) {
	var commencementLetterDTO auditDTOs.AuditPlanCommencementLetterDetailDTO

	// Fetch full details of the commencement letter by ID
	commencementLetter, err := apcls.repo.GetFullDetailsByID(auditPlanID, letterID, tx...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, letterID)
		}
		return nil, err
	}

	// Decode the model into DTO
	err = mapstructure.Decode(commencementLetter, &commencementLetterDTO)
	if err != nil {
		return nil, err
	}

	return &commencementLetterDTO, nil
}

func (aps *AuditPlanCommencementLetterService) GetByID(auditPlanID uint, letterID uint) (*models.AuditPlanCommencementLetter, error) {
	if _, err := aps.auditPlanSvc.GetByID(auditPlanID); err != nil {
		return nil, err
	}

	commencementLetter, err := aps.repo.GetByID(letterID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, letterID)
		}
		return nil, err
	}

	// Return the found commencement letter
	return commencementLetter, nil
}

func (aps *AuditPlanCommencementLetterService) Create(auditPlanId uint, currentUserID uint) error {
	tx := aps.repo.BeginTransaction()

	var createdIds []uint
	var addedAuditPlanCommencementLetter []map[string]interface{}

	auditPlanDetails, err := aps.auditPlanSvc.GetFullDetailsByID(auditPlanId)
	if err != nil {
		return err
	}

	// Group Auditable Items by ResponsibleUnitID
	unitToItemsMap := make(map[uint][]*auditDTOs.AuditPlanAuditableItemDTO)
	for _, item := range auditPlanDetails.AuditableItems {
		unitToItemsMap[item.AuditableItem.ResponsibleUnitID] = append(unitToItemsMap[item.AuditableItem.ResponsibleUnitID], item)
	}

	auditPlanCommencementLetterStatus, err := aps.auditPlanCommencementLetterStatusSvc.GetByType("draft")
	if err != nil {
		return err
	}

	for unitID, auditableItems := range unitToItemsMap {
		subjectEn := fmt.Sprintf("Audit Commencement Letter - %s", auditPlanDetails.Title)
		subjectAr := fmt.Sprintf("خطاب بدء التدقيق - %s", auditPlanDetails.Title)
		bodyEn := generateLetterBodyEn(auditableItems, auditPlanDetails)
		bodyAr := generateLetterBodyAr(auditableItems, auditPlanDetails)

		commencementLetter := &models.AuditPlanCommencementLetter{
			AuditPlanID:                         auditPlanId,
			ResponsibleUnitID:                   unitID,
			SubjectEn:                           subjectEn,
			SubjectAr:                           subjectAr,
			BodyEn:                              bodyEn,
			BodyAr:                              bodyAr,
			FromUserID:                          *auditPlanDetails.TeamLeadID,
			AuditPlanCommencementLetterStatusID: auditPlanCommencementLetterStatus.ID,
			CreatedBYID:                         currentUserID,
			CreatedAt:                           time.Now(),
		}

		id, err := aps.repo.Create(commencementLetter, tx)
		if err != nil {
			tx.Rollback()
			return err
		}
		addedAuditPlanCommencementLetter = append(addedAuditPlanCommencementLetter, map[string]interface{}{
			"id":                  id,
			"commencement_letter": commencementLetter,
		})

		createdIds = append(createdIds, id)
	}

	tx.Commit()
	go func() {
		err := aps.activityLogSvc.LogBulkActivities(
			models.CreateAction,
			models.AuditPlanCommencementLetter{}.TableName(),
			createdIds,
			addedAuditPlanCommencementLetter,
			currentUserID,
			nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.AuditPlanCommencementLetter{}.TableName(),
				"recordID":  auditPlanId,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log control activity", fields)
		}
	}()
	return nil
}

func (aps *AuditPlanCommencementLetterService) Update(auditPlanID uint, letterID uint, currentUserID uint, commencementLetterFormDTO *auditDTOs.AuditPlanCommencementLetterFormDTO) error {
	if _, err := aps.auditPlanSvc.GetByID(auditPlanID); err != nil {
		return err
	}

	oldCommencementLetter, err := aps.GetFullDetailsByID(auditPlanID, letterID)
	if err != nil {
		return err
	}

	existingCommencementLetter, err := aps.GetByID(auditPlanID, letterID)
	if err != nil {
		return err
	}

	// Update the fields of the commencement letter
	existingCommencementLetter.SubjectEn = commencementLetterFormDTO.SubjectEn
	existingCommencementLetter.SubjectAr = commencementLetterFormDTO.SubjectAr
	existingCommencementLetter.BodyEn = commencementLetterFormDTO.BodyEn
	existingCommencementLetter.BodyAr = commencementLetterFormDTO.BodyAr
	existingCommencementLetter.FromUserID = commencementLetterFormDTO.FromUserID
	existingCommencementLetter.UpdatedBYID = &currentUserID

	// Call the repository to update the commencement letter in the database
	err = aps.repo.Update(existingCommencementLetter)
	if err != nil {
		return err
	}
	go func() {
		existingCommencementLetter, err := aps.GetFullDetailsByID(auditPlanID, letterID)
		if err != nil {
			return
		}
		// Prepare log data
		changedData := make(map[string]interface{})
		// log changes
		oldData := map[string]interface{}{
			"subject_en": oldCommencementLetter.SubjectEn,
			"subject_ar": oldCommencementLetter.SubjectAr,
			"body_en":    oldCommencementLetter.BodyEn,
			"body_ar":    oldCommencementLetter.BodyAr,
			"from_user":  fmt.Sprintf("%s %s", oldCommencementLetter.FromUser.FirstNameEn, oldCommencementLetter.FromUser.LastNameEn),
		}

		newData := map[string]interface{}{
			"subject_en": existingCommencementLetter.SubjectEn,
			"subject_ar": existingCommencementLetter.SubjectAr,
			"body_en":    existingCommencementLetter.BodyEn,
			"body_ar":    existingCommencementLetter.BodyAr,
			"from_user":  fmt.Sprintf("%s %s", existingCommencementLetter.FromUser.FirstNameEn, existingCommencementLetter.FromUser.LastNameEn),
		}
		// Compare other fields and add them to changedData
		fieldChanges := utility.CompareFieldChanges(oldData, newData)
		for k, v := range fieldChanges {
			changedData[k] = v
		}
		// Log activity if there are any changes
		if len(changedData) > 0 {
			err := aps.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.AuditPlanCommencementLetter{}.TableName(),
				letterID,
				changedData,
				currentUserID,
				nil,
			)
			if err != nil {
				fields := map[string]interface{}{
					"action":    models.UpdateAction,
					"tableName": models.AuditPlanCommencementLetter{}.TableName(),
					"recordID":  letterID,
					"userID":    currentUserID,
					"error":     err.Error(),
				}
				fmt.Printf("Failed to log activity: %v\n", fields)
			}
		}
	}()
	return nil
}

func (aps *AuditPlanCommencementLetterService) Delete(auditPlanID uint, letterID uint, currentUserID uint) (uint, error) {
	if _, err := aps.auditPlanSvc.GetByID(auditPlanID); err != nil {
		return 0, err
	}

	var failedID uint
	var deletedAuditPlanCommecementLetter []map[string]interface{}

	_, err := aps.GetByID(auditPlanID, letterID)
	if err != nil {
		return 0, err
	}

	auditPlanCommencementLetter, err := aps.GetFullDetailsByID(auditPlanID, letterID)
	if err != nil {
		return 0, err
	}

	if err := aps.repo.Delete(letterID); err != nil {
		failedID = letterID
		return failedID, err
	}

	deletedAuditPlanCommecementLetter = append(deletedAuditPlanCommecementLetter, map[string]interface{}{
		"id":         letterID,
		"audit_plan": auditPlanCommencementLetter.AuditPlan.Title,
		"subject_en": auditPlanCommencementLetter.SubjectEn,
		"subject_ar": auditPlanCommencementLetter.SubjectAr,
	})

	go func() {
		err := aps.activityLogSvc.LogActivity(
			models.DeleteAction,
			models.AuditPlanCommencementLetter{}.TableName(),
			letterID,
			deletedAuditPlanCommecementLetter,
			currentUserID,
			nil,
		)
		if err != nil {
			// Handle error if logging fails
			fmt.Printf("Failed to log delete activity: %v\n", err)
		}
	}()

	// Return nil for no failure
	return failedID, nil
}

func (aps *AuditPlanCommencementLetterService) Send(auditPlanID uint, letterID uint, currentUserID uint, commencementLetterSendFormDTO *auditDTOs.AuditPlanCommencementLetterSendFormDTO) error {
	// Validate if the audit plan exists
	auditPlan, err := aps.auditPlanSvc.GetByID(auditPlanID)
	if err != nil {
		return err
	}

	// Fetch the commencement letter details
	existingCommencementLetter, err := aps.GetByID(auditPlanID, letterID)
	if err != nil {
		return err
	}

	tx := aps.repo.BeginTransaction()

	auditPlanCommencementLetterStatus, err := aps.auditPlanCommencementLetterStatusSvc.GetByType("sent")
	if err != nil {
		return err
	}

	currentTime := time.Now()
	existingCommencementLetter.AuditPlanCommencementLetterStatusID = auditPlanCommencementLetterStatus.ID
	existingCommencementLetter.SentAt = &currentTime
	existingCommencementLetter.UpdatedBYID = &currentUserID

	err = aps.repo.Update(existingCommencementLetter, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	toUser := &models.AuditPlanCommencementLetterToUser{
		AuditPlanCommencementLetterID: letterID,
		UserID:                        commencementLetterSendFormDTO.ToUserID,
		CreatedByID:                   currentUserID,
	}

	err = aps.repo.CreateToUser(toUser, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	for _, ccUserID := range commencementLetterSendFormDTO.CCUserIDs {
		ccUser := &models.AuditPlanCommencementLetterCCUser{
			AuditPlanCommencementLetterID: letterID,
			UserID:                        *ccUserID,
			CreatedByID:                   currentUserID,
		}

		err := aps.repo.CreateCCUser(ccUser, tx)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()

	// Start sending notifications and emails asynchronously
	go func() {
		toUserDetails, err := aps.userSvc.GetByID(commencementLetterSendFormDTO.ToUserID)
		if err != nil {
			utility.LogError("Failed to fetch user details", map[string]interface{}{
				"error":  err.Error(),
				"userID": commencementLetterSendFormDTO.ToUserID,
			})
			return
		}

		fromUserDetails, err := aps.userSvc.GetFullDetailsByID(existingCommencementLetter.FromUserID)
		if err != nil {
			utility.LogError("Failed to fetch user details", map[string]interface{}{
				"error":  err.Error(),
				"userID": existingCommencementLetter.FromUserID,
			})
			return
		}

		// Fetch CC User Details
		var ccEmails []string
		for _, ccUserID := range commencementLetterSendFormDTO.CCUserIDs {
			ccUserDetails, err := aps.userSvc.GetByID(*ccUserID) // Ensure *ccUserID is dereferenced safely
			if err != nil {
				utility.LogError("Failed to fetch CC user details", map[string]interface{}{
					"error":  err.Error(),
					"userID": *ccUserID,
				})
				continue
			}
			ccEmails = append(ccEmails, ccUserDetails.Email)
		}

		// Prepare DB Notification (unchanged)
		entity := "audit-plan"
		urlParam := fmt.Sprintf(`{"id":%d}`, auditPlanID)
		action := "send-commencement-letter"
		dbNotificationDTO := &commonDTOs.DBNotificationDTO{
			UserIDs:              []uint{commencementLetterSendFormDTO.ToUserID},
			ModuleIdentifier:     "audit",
			NotificationType:     string(commonDTOs.AssignmentType),
			HeadingEn:            "Audit Commencement Letter Sent",
			HeadingAr:            "تم إرسال خطاب بدء التدقيق",
			TitleEn:              auditPlan.Title,
			TitleAr:              auditPlan.Title,
			Entity:               &entity,
			Action:               &action,
			URLParam:             urlParam,
			IsSystemNotification: false,
			CreatedByID:          &currentUserID,
		}

		// Prepare Email Notification
		emailTo := []string{toUserDetails.Email}
		emailData := map[string]interface{}{
			"FirstNameEn":          toUserDetails.FirstNameEn,
			"LastNameEn":           toUserDetails.LastNameEn,
			"SubjectEn":            existingCommencementLetter.SubjectEn,
			"BodyEn":               template.HTML(strings.ReplaceAll(existingCommencementLetter.BodyEn, "\n", "<br>")),
			"FromUserFirstNameEn":  fromUserDetails.FirstNameEn,
			"FromUserLastNameEn":   fromUserDetails.LastNameEn,
			"FromUserDesignation":  fromUserDetails.Designation.TitleEn,
			"FromUserMobileNumber": fromUserDetails.MobileNumber,
			"FromUserEmail":        fromUserDetails.Email,
		}
		emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
			EmailTo:              emailTo,
			Subject:              existingCommencementLetter.SubjectEn,
			TemplateRelativePath: "audit/templates/email/audit_commencement_letter_send.html",
			EmailData:            emailData,
		}

		// Include CC Emails if they exist
		if len(ccEmails) > 0 {
			emailNotificationDTO.CC = ccEmails
		}

		// Send Notifications and Email
		err = aps.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "audit")
		if err != nil {
			utility.LogError("Failed to send notifications", map[string]interface{}{
				"error": err.Error(),
				"user":  toUserDetails.Email,
			})
		}
	}()

	return nil
}

func generateLetterBodyEn(auditableItems []*auditDTOs.AuditPlanAuditableItemDTO, auditPlan *auditDTOs.AuditPlanDetailDTO) string {
	auditScope := ""
	var responsibleUnitName string

	for _, item := range auditableItems {
		// Generate the audit scope
		if item.AuditableItem.TitleEn != nil {
			auditScope += fmt.Sprintf("- %s\n", *item.AuditableItem.TitleEn) // Use English title if available
		}

		if item.AuditableItem.ResponsibleUnit.TitleEn != "" {
			responsibleUnitName = item.AuditableItem.ResponsibleUnit.TitleEn
		}
	}

	teamLeadName := fmt.Sprintf("%s %s", auditPlan.TeamLead.FirstNameEN, auditPlan.TeamLead.LastNameEN)

	body := fmt.Sprintf(`
Dear %s,

This is to inform you that JethurGRC will commence an audit of your department for the period %s to %s.

The audit scope will include:
%s

Please prepare relevant documents for review. Our team, led by %s, will contact you to schedule an initial meeting.

Your cooperation is appreciated. For any questions, please contact %s.

Thank you,
`, responsibleUnitName, auditPlan.StartDateString, auditPlan.TargetDateString, auditScope, teamLeadName, teamLeadName)

	return body
}

func generateLetterBodyAr(auditableItems []*auditDTOs.AuditPlanAuditableItemDTO, auditPlan *auditDTOs.AuditPlanDetailDTO) string {
	auditScope := ""
	var responsibleUnitName string

	for _, item := range auditableItems {
		// Generate the audit scope
		if item.AuditableItem.TitleAr != nil {
			auditScope += fmt.Sprintf("- %s\n", *item.AuditableItem.TitleAr) // Use Arabic title if available
		}

		if item.AuditableItem.ResponsibleUnit.TitleAr != "" {
			responsibleUnitName = item.AuditableItem.ResponsibleUnit.TitleAr
		}
	}

	teamLeadName := fmt.Sprintf("%s %s", auditPlan.TeamLead.FirstNameAR, auditPlan.TeamLead.LastNameAR)

	body := fmt.Sprintf(`
عزيزي %s،

نود إبلاغكم أن JethurGRC سيبدأ عملية التدقيق على إدارتكم خلال الفترة من %s إلى %s.

سيشمل نطاق التدقيق:
%s

يرجى إعداد المستندات ذات الصلة للمراجعة. سيقوم فريقنا، بقيادة %s، بالتواصل معكم لتحديد موعد الاجتماع الأولي.

نقدر تعاونكم. لأي استفسارات، يرجى التواصل مع %s.

شكرًا لكم،

`, responsibleUnitName, auditPlan.StartDateString, auditPlan.TargetDateString, auditScope, teamLeadName, teamLeadName)

	return body
}
