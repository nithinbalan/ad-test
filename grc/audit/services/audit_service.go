package services

import (
	"errors"
	"fmt"
	"strings"
	"time"

	auditDTOs "jethur/grc/audit/dtos"
	"jethur/grc/audit/repositories"
	commonDTOs "jethur/grc/common/dtos"
	commonServices "jethur/grc/common/services"
	logServices "jethur/grc/logs/services"
	masterServices "jethur/grc/masters/services"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type AuditServiceInterface interface {
	GetAll(params *auditDTOs.AuditListQueryParamsDTO) (*auditDTOs.AuditListResponseDTO, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*auditDTOs.AuditDetailDTO, error)
	GetByID(id uint) (*models.Audit, error)
	Create(currentUserID uint, auditFormDTO *auditDTOs.AuditFormDTO) (uint, error)
	Update(id uint, currentUserID uint, auditFormDTO *auditDTOs.AuditFormDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)

	GetAllAuditableItems(auditID uint, params *auditDTOs.AuditAuditableItemListQueryParamsDTO) (*auditDTOs.AuditAuditableItemListResponseDTO, error)
}

type AuditService struct {
	repo                 repositories.AuditRepositoryInterface
	activityLogSvc       logServices.ActivityLogServiceInterface
	auditPlanSvc         AuditPlanServiceInterface
	auditTimelineSvc     AuditTimelineServiceInterface
	auditPlanScheduleSvc AuditPlanScheduleServiceInterface
	statusSvc            masterServices.StatusServiceInterface
	auditStatusSvc       AuditStatusServiceInterface
	userSvc              organizationServices.UserServiceInterface
	notificationSvc      commonServices.NotificationServiceInterface
}

func NewAuditService(
	repo repositories.AuditRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	AuditPlanSvc AuditPlanServiceInterface,
	auditTimelineSvc AuditTimelineServiceInterface,
	auditPlanScheduleSvc AuditPlanScheduleServiceInterface,
	statusSvc masterServices.StatusServiceInterface,
	auditStatusSvc AuditStatusServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	notificationSvc commonServices.NotificationServiceInterface,
) AuditServiceInterface {
	return &AuditService{
		repo:                 repo,
		activityLogSvc:       activityLogSvc,
		auditPlanSvc:         AuditPlanSvc,
		auditTimelineSvc:     auditTimelineSvc,
		auditPlanScheduleSvc: auditPlanScheduleSvc,
		statusSvc:            statusSvc,
		auditStatusSvc:       auditStatusSvc,
		userSvc:              userSvc,
		notificationSvc:      notificationSvc,
	}
}

func (as *AuditService) GetAll(params *auditDTOs.AuditListQueryParamsDTO) (*auditDTOs.AuditListResponseDTO, error) {
	return as.repo.GetAll(params)
}

func (as *AuditService) GetByID(id uint) (*models.Audit, error) {
	audit, err := as.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return audit, nil
}

func (as *AuditService) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*auditDTOs.AuditDetailDTO, error) {
	var auditDTO auditDTOs.AuditDetailDTO

	audit, err := as.repo.GetFullDetailsByID(id, tx...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	err = mapstructure.Decode(audit, &auditDTO)
	if err != nil {
		return nil, err
	}

	auditDTO.StartDateString = audit.StartDate.Format("2006-01-02")
	auditDTO.EndDateString = audit.EndDate.Format("2006-01-02")

	auditDTO.CreatedAt = audit.CreatedAt
	auditDTO.UpdatedAt = audit.UpdatedAt

	return &auditDTO, nil
}

func (as *AuditService) Create(currentUserID uint, auditFormDTO *auditDTOs.AuditFormDTO) (uint, error) {
	// First, check if an audit already exists for the given AuditPlanID
	existingAudit, err := as.repo.GetAuditByAuditPlanID(auditFormDTO.AuditPlanID)
	if err != nil {
		return 0, err // Handle DB query error
	}
	// If an audit already exists, return an error without starting a transaction
	if existingAudit != nil {
		return 0, utility.NewForbiddenError("Audit has already been created for this audit plan")
	}
	tx := as.repo.BeginTransaction()

	dates, err := as.parseAuditPlanDates(auditFormDTO)
	if err != nil {
		return 0, err
	}

	if err := as.validateAuditFormDTO(auditFormDTO); err != nil {
		tx.Rollback()
		return 0, err
	}
	auditStatus, err := as.auditStatusSvc.GetByType("in-progress")
	if err != nil {
		return 0, err
	}
	audit := &models.Audit{
		AuditPlanID:   auditFormDTO.AuditPlanID,
		StartDate:     dates.StartDate,
		EndDate:       dates.EndDate,
		AuditStatusID: uint(auditStatus.ID),
		CreatedBYID:   currentUserID,
	}

	id, err := as.repo.Create(audit, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	if err := as.copyAuditPlanSchedulesToAuditSchedule(id, auditFormDTO.AuditPlanID, currentUserID, tx); err != nil {
		tx.Rollback()
		return 0, err
	}

	timelineTypes, err := as.auditTimelineSvc.GetAllTypes()
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	// Get timeline statuses
	completedStatus, err := as.auditTimelineSvc.GetByStatusType("completed")
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	pendingStatus, err := as.auditTimelineSvc.GetByStatusType("pending")
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	var timelineFormDTOs []*auditDTOs.AuditTimelineFormDTO

	for _, timelineType := range timelineTypes {
		var (
			statusID      uint
			performedByID *uint
			currentTime   *time.Time
		)

		if timelineType.Type == "initiated" {
			statusID = completedStatus.ID
			performedByID = &currentUserID
			now := time.Now()
			currentTime = &now
		} else {
			statusID = pendingStatus.ID
			performedByID = nil
			currentTime = nil
		}

		timelineFormDTO := &auditDTOs.AuditTimelineFormDTO{
			AuditID:             id,
			AuditTimelineTypeID: timelineType.ID,
			TimelineStatusID:    statusID,
			PerformedByID:       performedByID,
			PerformedAt:         currentTime,
		}
		timelineFormDTOs = append(timelineFormDTOs, timelineFormDTO)
	}
	if err := as.auditTimelineSvc.Create(timelineFormDTOs, tx); err != nil {
		tx.Rollback()
		return 0, err
	}

	if err := as.CopyChecklistsToAudit(auditFormDTO.AuditPlanID, id, currentUserID, tx); err != nil {
		tx.Rollback()
		return 0, err
	}
	// Update Audit plan status in to in-progress
	auditStatusID := auditStatus.ID
	if err := as.auditPlanSvc.StatusUpdate(auditFormDTO.AuditPlanID, currentUserID, uint(auditStatusID), tx); err != nil {
		tx.Rollback()
		return 0, err
	}
	tx.Commit()

	go func() {
		err := as.activityLogSvc.LogActivity(
			models.CreateAction,
			models.Audit{}.TableName(),
			id,
			auditFormDTO,
			currentUserID,
			nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.Audit{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}

		auditPlan, err := as.auditPlanSvc.GetFullDetailsByID(auditFormDTO.AuditPlanID)
		if err != nil {
			utility.LogError("Failed to fetch audit plan details", map[string]interface{}{
				"auditPlanID": auditFormDTO.AuditPlanID,
				"error":       err.Error(),
			})
			return
		}

		if auditPlan == nil {
			utility.LogError("AuditPlan is nil", map[string]interface{}{
				"auditPlanID": auditFormDTO.AuditPlanID,
			})
			return
		}

		uniqueUserIDs := make(map[uint]bool)

		// Check TeamLead
		if auditPlan.TeamLead != nil && auditPlan.TeamLead.ID != 0 {
			uniqueUserIDs[uint(auditPlan.TeamLead.ID)] = true
		} else {
			utility.LogError("TeamLead is missing or invalid", nil)
		}

		auditMembers := []string{}
		if auditPlan.Members != nil {
			for _, member := range auditPlan.Members {
				if member.User.ID != 0 {
					uniqueUserIDs[uint(member.User.ID)] = true
					auditMembers = append(auditMembers, fmt.Sprintf("%s %s", member.User.FirstNameEN, member.User.LastNameEN))
				} else {
					utility.LogError("Audit Member User is invalid", map[string]interface{}{
						"memberID": member.ID,
					})
				}
			}
		} else {
			utility.LogError("Audit Members are nil", nil)
		}

		auditMembersDetails := "None"
		if len(auditMembers) > 0 {
			auditMembersDetails = strings.Join(auditMembers, "\n- ")
		}

		assignedBy, err := as.userSvc.GetByID(currentUserID)
		if err != nil {
			utility.LogError("Failed to fetch assigning user", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		for userID := range uniqueUserIDs {
			user, err := as.userSvc.GetByID(userID)
			if err != nil {
				utility.LogError("Failed to fetch user details for notification", map[string]interface{}{
					"userID": userID,
					"error":  err.Error(),
				})
				continue
			}

			urlParam := fmt.Sprintf(`{"audit_plan_id":%d}`, auditPlan.ID)
			entity := "audit-plan-publish"
			action := "responsible-user-assign"
			appURL := utility.GetAppURL()

			dbNotificationDTO := &commonDTOs.DBNotificationDTO{
				UserID:               userID,
				ModuleIdentifier:     "audit",
				NotificationType:     string(commonDTOs.InformationType),
				ReferenceNumber:      auditPlan.ReferenceNumber,
				HeadingEn:            "Audit Created",
				HeadingAr:            "تم إنشاء التدقيق",
				TitleEn:              auditPlan.Title,
				TitleAr:              auditPlan.Title,
				Entity:               &entity,
				Action:               &action,
				URLParam:             urlParam,
				IsSystemNotification: false,
				CreatedByID:          &currentUserID,
			}

			emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
				EmailTo:              []string{user.Email},
				Subject:              fmt.Sprintf("New Audit Created: %s", auditPlan.Title),
				TemplateRelativePath: "audit/templates/email/audit_plan_publish_create_responsible_user_mail copy.html",
				EmailData: map[string]interface{}{
					"FirstNameEn":           user.FirstNameEn,
					"LastNameEn":            user.LastNameEn,
					"AuditTitle":            auditPlan.Title,
					"Description":           auditPlan.Description,
					"TeamLead":              fmt.Sprintf("%s %s", auditPlan.TeamLead.FirstNameEN, auditPlan.TeamLead.LastNameEN),
					"AuditMembers":          auditMembersDetails,
					"ReferenceNumber":       auditPlan.ReferenceNumber,
					"AuditPlanStartDate":    auditPlan.StartDateString,
					"AuditPlanEndDate":      auditPlan.TargetDateString,
					"AssignedByFirstNameEn": assignedBy.FirstNameEn,
					"AssignedByLastNameEn":  assignedBy.LastNameEn,
					"AssignedOn":            time.Now().Format("2006-01-02"),
					"RedirectURL":           fmt.Sprintf("%s/audit/audits", appURL),
				},
			}

			// Send Notifications
			err = as.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "audit")
			if err != nil {
				utility.LogError("Failed to send audit plan notifications", map[string]interface{}{
					"userID": userID,
					"error":  err.Error(),
				})
			}
		}
	}()

	return id, nil
}

func (as *AuditService) Update(id uint, currentUserID uint, auditFormDTO *auditDTOs.AuditFormDTO) error {
	tx := as.repo.BeginTransaction()

	existingAudit, err := as.repo.GetByID(id)
	if err != nil {
		tx.Rollback()
		return utility.NewResourceNotFoundError(err, id)
	}

	if err := as.validateAuditFormDTO(auditFormDTO); err != nil {
		tx.Rollback()
		return err
	}

	dates, err := as.parseAuditPlanDates(auditFormDTO)
	if err != nil {
		return err
	}

	existingAudit.AuditPlanID = auditFormDTO.AuditPlanID
	existingAudit.StartDate = dates.StartDate
	existingAudit.EndDate = dates.EndDate
	existingAudit.UpdatedBYID = &currentUserID

	err = as.repo.Update(existingAudit, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	go func() {
		err := as.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.Audit{}.TableName(),
			id,
			auditFormDTO,
			currentUserID,
			nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.UpdateAction,
				"tableName": models.Audit{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return nil
}

func (as *AuditService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	tx := as.repo.BeginTransaction()
	var failedIDs []uint
	var successfulIDs []uint

	for _, id := range ids {
		// Get existing audit details
		_, err := as.GetFullDetailsByID(id, tx)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				failedIDs = append(failedIDs, id)
				continue
			}
			tx.Rollback()
			return nil, err
		}

		if err := as.auditTimelineSvc.Delete(id, tx); err != nil {
			tx.Rollback()
			return nil, err
		}

		// Delete the audit itself
		if err := as.repo.Delete(id, tx); err != nil {
			failedIDs = append(failedIDs, id)
			continue
		}

		successfulIDs = append(successfulIDs, id)
	}

	if len(successfulIDs) > 0 {
		tx.Commit()

		go func() {
			for _, id := range successfulIDs {
				err := as.activityLogSvc.LogActivity(
					models.DeleteAction,
					models.Audit{}.TableName(),
					id,
					nil,
					currentUserID,
					nil,
				)
				if err != nil {
					fields := map[string]interface{}{
						"action":    models.DeleteAction,
						"tableName": models.Audit{}.TableName(),
						"recordID":  id,
						"userID":    currentUserID,
						"error":     err.Error(),
					}

					utility.LogError("Failed to log activity", fields)
				}
			}
		}()
	} else {
		tx.Rollback()
	}

	return failedIDs, nil
}

func (as *AuditService) validateAuditFormDTO(dto *auditDTOs.AuditFormDTO) error {
	if _, err := as.auditPlanSvc.GetByID(dto.AuditPlanID); err != nil {
		return utility.NewValidationError(fmt.Errorf("audit plan does not exist"), "not_exist", "audit plan")
	}

	return nil
}

func (as *AuditService) parseAuditPlanDates(dto *auditDTOs.AuditFormDTO) (*auditDTOs.AuditParsedDates, error) {
	parsed := &auditDTOs.AuditParsedDates{}

	var err error

	parsed.StartDate, err = utility.ParseMandatoryDate(dto.StartDate)
	if err != nil {
		return nil, utility.NewValidationError(fmt.Errorf("invalid start date"), "invalid_date", "start date")
	}

	parsed.EndDate, err = utility.ParseMandatoryDate(dto.EndDate)
	if err != nil {
		return nil, utility.NewValidationError(fmt.Errorf("invalid end date"), "invalid_date", "end date")
	}

	return parsed, nil
}

func (as *AuditService) copyAuditPlanSchedulesToAuditSchedule(auditID uint, auditPlanID uint, currentUserID uint, tx ...utility.DatabaseTransaction) error {
	auditPlanSchedules, err := as.auditPlanScheduleSvc.GetAllSchedulesByAuditPlanID(auditPlanID, tx...)
	if err != nil {
		return err
	}

	for _, schedule := range auditPlanSchedules {
		auditSchedule := &models.AuditSchedule{
			AuditPlanID:           schedule.AuditPlanID,
			AuditPlanScheduleID:   &schedule.ID,
			AuditActivityTypeID:   schedule.AuditActivityTypeID,
			AuditID:               auditID,
			Title:                 schedule.Title,
			Description:           schedule.Description,
			StartDate:             schedule.StartDate,
			TargetDate:            schedule.TargetDate,
			ResponsibleUserID:     schedule.ResponsibleUserID,
			AuditScheduleStatusID: schedule.AuditScheduleStatusID,
			CreatedBYID:           currentUserID,
		}

		createdAuditSchedule, err := as.repo.CreateAuditSchedule(auditSchedule, tx...)
		if err != nil {
			return err
		}

		statusUpdate := &models.AuditScheduleStatusUpdate{
			AuditScheduleID:       createdAuditSchedule.ID,
			AuditScheduleStatusID: auditSchedule.AuditScheduleStatusID,
			Comment:               "Initial audit schedule status",
			CreatedBYID:           currentUserID,
		}
		if err := as.repo.CreateAuditScheduleStatusUpdate(statusUpdate, tx...); err != nil {
			return err
		}

		for _, auditableItem := range schedule.AuditableItems {
			auditScheduleAuditableItem := &models.AuditScheduleAuditableItem{
				AuditScheduleID: createdAuditSchedule.ID,
				AuditableItemID: auditableItem.AuditableItemID,
				CreatedBYID:     currentUserID,
			}

			if err := as.repo.CreateAuditScheduleAuditableItem(auditScheduleAuditableItem, tx...); err != nil {
				return err
			}
		}
	}

	return nil
}

func (as *AuditService) CopyChecklistsToAudit(auditPlanID uint, auditID uint, currentUserID uint, tx ...utility.DatabaseTransaction) error {
	// Retrieve the full details of the audit plan, including the checklists
	auditPlan, err := as.auditPlanSvc.GetFullDetailsByID(auditPlanID, tx...)
	if err != nil {
		return err
	}

	status, err := as.statusSvc.GetByType("pending")
	if err != nil {
		return err
	}

	var auditChecklists []*models.AuditChecklist
	if len(auditPlan.AuditPlanChecklists) > 0 {
		// Iterate through each checklist in the audit plan and create a corresponding audit checklist
		for _, checklist := range auditPlan.AuditPlanChecklists {
			auditChecklist := &models.AuditChecklist{
				AuditID:              auditID,
				AuditPlanChecklistID: &checklist.ID,
				AuditableItemID:      &checklist.AuditableItemID,
				Question:             checklist.Question,
				CreatedBYID:          uint(auditPlan.CreatedByUser.ID),
				CreatedAt:            auditPlan.CreatedAt,
				StatusID:             status.ID,
			}
			auditChecklists = append(auditChecklists, auditChecklist)
		}
	}

	// Save the new audit checklists to the database
	if len(auditChecklists) > 0 {
		if err := as.repo.CreateAuditChecklists(auditChecklists, tx...); err != nil {
			return err
		}
	}

	return nil
}

func (as *AuditService) GetAllAuditableItems(auditID uint, params *auditDTOs.AuditAuditableItemListQueryParamsDTO) (*auditDTOs.AuditAuditableItemListResponseDTO, error) {
	if _, err := as.GetByID(auditID); err != nil {
		return nil, err
	}
	return as.repo.GetAllAuditableItems(auditID, params)
}
