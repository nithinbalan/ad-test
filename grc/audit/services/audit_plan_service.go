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
	organizationServices "jethur/grc/organization/services"

	"jethur/grc/models"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type AuditPlanServiceInterface interface {
	GetAll(params *auditDTOs.AuditPlanListQueryParamsDTO) (*auditDTOs.AuditPlanListResponseDTO, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*auditDTOs.AuditPlanDetailDTO, error)
	GetByID(id uint) (*models.AuditPlan, error)
	Create(currentUserID uint, auditPlanFormDTO *auditDTOs.AuditPlanFormDTO) (uint, error)
	Update(id uint, currentUserID uint, auditPlanFormDTO *auditDTOs.AuditPlanFormDTO) error
	Publish(id uint, currentUserID uint) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	StatusUpdate(id uint, currentUserID uint, auditStatusID uint, tx utility.DatabaseTransaction) error
}

type AuditPlanService struct {
	repo             repositories.AuditPlanRepositoryInterface
	activityLogSvc   logServices.ActivityLogServiceInterface
	auditCategorySvc AuditCategoryServiceInterface
	auditStatusSvc   AuditStatusServiceInterface
	auditTeamSvc     AuditTeamServiceInterface
	userSvc          organizationServices.UserServiceInterface
	auditableItemSvc AuditableItemServiceInterface
	notificationSvc  commonServices.NotificationServiceInterface
}

func NewAuditPlanService(
	repo repositories.AuditPlanRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	auditCategorySvc AuditCategoryServiceInterface,
	auditStatusSvc AuditStatusServiceInterface,
	auditTeamSvc AuditTeamServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	auditableItemSvc AuditableItemServiceInterface,
	notificationSvc commonServices.NotificationServiceInterface,
) AuditPlanServiceInterface {
	return &AuditPlanService{
		repo:             repo,
		activityLogSvc:   activityLogSvc,
		auditCategorySvc: auditCategorySvc,
		auditStatusSvc:   auditStatusSvc,
		auditTeamSvc:     auditTeamSvc,
		userSvc:          userSvc,
		auditableItemSvc: auditableItemSvc,
		notificationSvc:  notificationSvc,
	}
}

func (aps *AuditPlanService) GetAll(params *auditDTOs.AuditPlanListQueryParamsDTO) (*auditDTOs.AuditPlanListResponseDTO, error) {
	return aps.repo.GetAll(params)
}

func (aps *AuditPlanService) GetByID(id uint) (*models.AuditPlan, error) {
	auditPlan, err := aps.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return auditPlan, nil
}

func (aps *AuditPlanService) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*auditDTOs.AuditPlanDetailDTO, error) {
	var auditPlanDTO auditDTOs.AuditPlanDetailDTO

	auditPlan, err := aps.repo.GetFullDetailsByID(id, tx...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	err = mapstructure.Decode(auditPlan, &auditPlanDTO)
	if err != nil {
		return nil, err
	}

	if auditPlan.StartDate != nil {
		auditPlanDTO.StartDateString = auditPlan.StartDate.Format("2006-01-02")
	}
	if auditPlan.TargetDate != nil {
		auditPlanDTO.TargetDateString = auditPlan.TargetDate.Format("2006-01-02")
	}
	if auditPlan.PublishedAt != nil {
		auditPlanDTO.PublishedDateString = auditPlan.PublishedAt.Format("2006-01-02")
	}
	auditPlanDTO.CreatedAt = auditPlan.CreatedAt
	auditPlanDTO.UpdatedAt = auditPlan.UpdatedAt

	return &auditPlanDTO, nil
}

func (aps *AuditPlanService) Create(currentUserID uint, auditPlanFormDTO *auditDTOs.AuditPlanFormDTO) (uint, error) {
	tx := aps.repo.BeginTransaction()
	dates, err := aps.parseAuditPlanDates(auditPlanFormDTO)
	if err != nil {
		return 0, err
	}

	if err := aps.validateAuditPlanFormDTO(auditPlanFormDTO); err != nil {
		tx.Rollback()
		return 0, err
	}
	referenceNumber, err := aps.repo.GenerateReferenceNumber()
	if err != nil {
		return 0, err
	}
	auditStatus, err := aps.auditStatusSvc.GetByType("draft")
	if err != nil {
		return 0, err
	}

	auditPlan := &models.AuditPlan{
		ReferenceNumber: referenceNumber,
		AuditCategoryID: auditPlanFormDTO.AuditCategoryID,
		AuditTeamID:     auditPlanFormDTO.AuditTeamID,
		TeamLeadID:      auditPlanFormDTO.TeamLeadID,
		Title:           auditPlanFormDTO.Title,
		Description:     auditPlanFormDTO.Description,
		StartDate:       dates.StartDate,
		TargetDate:      dates.TargetDate,
		Purpose:         auditPlanFormDTO.Purpose,
		Objective:       auditPlanFormDTO.Objective,
		AuditStatusID:   uint(auditStatus.ID),
		CreatedBYID:     currentUserID,
	}

	id, err := aps.repo.Create(auditPlan, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	var auditMembers []*models.AuditPlanAuditMember
	for _, userID := range auditPlanFormDTO.UserIds {
		auditMember := &models.AuditPlanAuditMember{
			AuditPlanID: id,
			UserID:      *userID,
			CreatedBYID: currentUserID,
		}
		auditMembers = append(auditMembers, auditMember)
	}

	if len(auditMembers) > 0 {
		if err := aps.repo.CreateAuditPlanMembers(auditMembers, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	var auditableItems []*models.AuditPlanAuditableItem
	for _, itemID := range auditPlanFormDTO.AuditableItemIds {
		auditableItem := &models.AuditPlanAuditableItem{
			AuditPlanID:     id,
			AuditableItemID: *itemID,
			CreatedBYID:     currentUserID,
		}
		auditableItems = append(auditableItems, auditableItem)
	}

	if len(auditableItems) > 0 {
		if err := aps.repo.CreateAuditPlanAuditableItems(auditableItems, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	tx.Commit()

	go func() {
		existingAuditPlan, err := aps.GetFullDetailsByID(id)
		if err != nil {
			utility.LogError("Failed to fetch audit plan", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		err = aps.activityLogSvc.LogActivity(
			models.CreateAction,
			models.AuditPlan{}.TableName(),
			id,
			existingAuditPlan,
			currentUserID,
			nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.AuditPlan{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return id, nil
}

func (aps *AuditPlanService) Update(id uint, currentUserID uint, auditPlanFormDTO *auditDTOs.AuditPlanFormDTO) error {
	tx := aps.repo.BeginTransaction()
	dates, err := aps.parseAuditPlanDates(auditPlanFormDTO)
	if err != nil {
		return err
	}

	oldAuditPlan, err := aps.repo.GetFullDetailsByID(id)
	if err != nil {
		tx.Rollback()
		return utility.NewResourceNotFoundError(err, id)
	}

	existingAuditPlan, err := aps.repo.GetByID(id)
	if err != nil {
		tx.Rollback()
		return utility.NewResourceNotFoundError(err, id)
	}

	if err := aps.validateAuditPlanFormDTO(auditPlanFormDTO); err != nil {
		tx.Rollback()
		return err
	}

	existingAuditPlan.AuditCategoryID = auditPlanFormDTO.AuditCategoryID
	existingAuditPlan.AuditTeamID = auditPlanFormDTO.AuditTeamID
	existingAuditPlan.TeamLeadID = auditPlanFormDTO.TeamLeadID
	existingAuditPlan.Title = auditPlanFormDTO.Title
	existingAuditPlan.Description = auditPlanFormDTO.Description
	existingAuditPlan.StartDate = dates.StartDate
	existingAuditPlan.TargetDate = dates.TargetDate
	existingAuditPlan.Purpose = auditPlanFormDTO.Purpose
	existingAuditPlan.Objective = auditPlanFormDTO.Objective
	existingAuditPlan.UpdatedBYID = &currentUserID

	err = aps.repo.Update(existingAuditPlan, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	existingAuditPlanDetails, err := aps.GetFullDetailsByID(id, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	existingUserIDs := make(map[uint]bool)
	for _, member := range existingAuditPlanDetails.Members {
		existingUserIDs[member.UserID] = true
	}

	newUserIDs := make(map[uint]bool)
	for _, userID := range auditPlanFormDTO.UserIds {
		newUserIDs[*userID] = true
	}

	var usersToAdd []*models.AuditPlanAuditMember
	for _, userID := range auditPlanFormDTO.UserIds {
		if _, exists := existingUserIDs[*userID]; !exists {
			user := &models.AuditPlanAuditMember{
				AuditPlanID: id,
				UserID:      *userID,
				CreatedBYID: currentUserID,
			}
			usersToAdd = append(usersToAdd, user)
		}
	}

	var usersToRemove []*models.AuditPlanAuditMember
	for userID := range existingUserIDs {
		if _, exists := newUserIDs[userID]; !exists {
			usersToRemove = append(usersToRemove, &models.AuditPlanAuditMember{
				AuditPlanID: id,
				UserID:      userID,
			})
		}
	}

	// Extract user IDs to remove
	var userIDsToRemove []*uint
	for _, user := range usersToRemove {
		userID := user.UserID
		userIDsToRemove = append(userIDsToRemove, &userID)
	}

	// Add new members
	if len(usersToAdd) > 0 {
		if err := aps.repo.CreateAuditPlanMembers(usersToAdd, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	// Remove old members
	if len(userIDsToRemove) > 0 {
		if err := aps.repo.DeleteAuditPlanMembers(userIDsToRemove, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	// Handle auditable items
	existingAuditableItemIDs := make(map[uint]bool)
	for _, item := range existingAuditPlanDetails.AuditableItems {
		existingAuditableItemIDs[item.AuditableItemID] = true
	}

	newAuditableItemIDs := make(map[uint]bool)
	for _, itemID := range auditPlanFormDTO.AuditableItemIds {
		newAuditableItemIDs[*itemID] = true
	}

	var itemsToAdd []*models.AuditPlanAuditableItem
	for _, itemID := range auditPlanFormDTO.AuditableItemIds {
		if _, exists := existingAuditableItemIDs[*itemID]; !exists {
			item := &models.AuditPlanAuditableItem{
				AuditPlanID:     id,
				AuditableItemID: *itemID,
				CreatedBYID:     currentUserID,
			}
			itemsToAdd = append(itemsToAdd, item)
		}
	}

	var auditableItemIDsToRemove []*uint
	for itemID := range existingAuditableItemIDs {
		if _, exists := newAuditableItemIDs[itemID]; !exists {
			auditableItemIDsToRemove = append(auditableItemIDsToRemove, &itemID)
		}
	}

	if len(itemsToAdd) > 0 {
		if err := aps.repo.CreateAuditPlanAuditableItems(itemsToAdd, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	if len(auditableItemIDsToRemove) > 0 {
		if err := aps.repo.DeleteAuditPlanAuditableItems(auditableItemIDsToRemove, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()

	go func() {
		existingPlan, err := aps.repo.GetFullDetailsByID(id)
		if err != nil {
			return
		}
		// Determine auditable items and users changes
		auditableItemsAdded := []uint{}
		auditableItemsRemoved := []uint{}
		for id := range newAuditableItemIDs {
			if !existingAuditableItemIDs[id] {
				auditableItemsAdded = append(auditableItemsAdded, id)
			}
		}
		for id := range existingAuditableItemIDs {
			if !newAuditableItemIDs[id] {
				auditableItemsRemoved = append(auditableItemsRemoved, id)
			}
		}

		usersAdded := []uint{}
		usersRemoved := []uint{}
		for id := range newUserIDs {
			if !existingUserIDs[id] {
				usersAdded = append(usersAdded, id)
			}
		}
		for id := range existingUserIDs {
			if !newUserIDs[id] {
				usersRemoved = append(usersRemoved, id)
			}
		}

		// Prepare log data
		changedData := make(map[string]interface{})
		associatedItems := make(map[string]interface{})

		if len(auditableItemsAdded) > 0 || len(auditableItemsRemoved) > 0 {
			associatedItems["auditable_item_ids"] = map[string]interface{}{
				"added":   auditableItemsAdded,
				"removed": auditableItemsRemoved,
			}
		}

		if len(usersAdded) > 0 || len(usersRemoved) > 0 {
			associatedItems["user_ids"] = map[string]interface{}{
				"added":   usersAdded,
				"removed": usersRemoved,
			}
		}

		if len(associatedItems) > 0 {
			changedData["associated_items"] = associatedItems
		}
		//log changes
		oldData := map[string]interface{}{
			"audit_category": oldAuditPlan.AuditCategory.TitleEn,
			"audit_team":     oldAuditPlan.AuditTeam.TitleEn,
			"team_lead":      fmt.Sprintf("%s %s", oldAuditPlan.TeamLead.FirstNameEn, oldAuditPlan.TeamLead.LastNameEn),
			"title":          oldAuditPlan.Title,
			"description":    oldAuditPlan.Description,
			"start_date":     oldAuditPlan.StartDate,
			"target_date":    oldAuditPlan.TargetDate,
			"purpose":        oldAuditPlan.Purpose,
			"objective":      oldAuditPlan.Objective,
		}

		newData := map[string]interface{}{
			"audit_category": existingPlan.AuditCategory.TitleEn,
			"audit_team":     existingPlan.AuditTeam.TitleEn,
			"team_lead":      fmt.Sprintf("%s %s", existingPlan.TeamLead.FirstNameEn, existingPlan.TeamLead.LastNameEn),
			"title":          existingPlan.Title,
			"description":    existingPlan.Description,
			"start_date":     existingPlan.StartDate,
			"target_date":    existingPlan.TargetDate,
			"purpose":        existingPlan.Purpose,
			"objective":      existingPlan.Objective,
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
				models.AuditPlan{}.TableName(),
				id,
				changedData,
				currentUserID,
				nil,
			)
			if err != nil {
				fields := map[string]interface{}{
					"action":    models.UpdateAction,
					"tableName": models.AuditPlan{}.TableName(),
					"recordID":  id,
					"userID":    currentUserID,
					"error":     err.Error(),
				}

				utility.LogError("Failed to log activity", fields)
			}
		}
	}()

	return nil
}

func (aps *AuditPlanService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	tx := aps.repo.BeginTransaction()
	var failedIDs []uint
	var successfulIDs []uint
	var deletedAuditPlan []map[string]interface{}
	for _, id := range ids {
		// Get existing audit plan details
		existingAuditPlan, err := aps.GetFullDetailsByID(id, tx)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				failedIDs = append(failedIDs, id)
				continue
			}
			tx.Rollback()
			return nil, err
		}

		if len(existingAuditPlan.AuditableItems) > 0 {
			var deleteAuditableItemIDs []*uint
			for _, item := range existingAuditPlan.AuditableItems {
				deleteAuditableItemIDs = append(deleteAuditableItemIDs, item.ID)
			}
			if err := aps.repo.DeleteAuditPlanAuditableItems(deleteAuditableItemIDs, tx); err != nil {
				tx.Rollback()
				return nil, err
			}
		}

		if len(existingAuditPlan.Members) > 0 {
			var deleteMemberIDs []*uint
			for _, member := range existingAuditPlan.Members {
				deleteMemberIDs = append(deleteMemberIDs, member.ID)
			}
			if err := aps.repo.DeleteAuditPlanMembers(deleteMemberIDs, tx); err != nil {
				tx.Rollback()
				return nil, err
			}
		}

		// Delete the audit plan itself
		if err := aps.repo.Delete(id, tx); err != nil {
			failedIDs = append(failedIDs, id)
			continue
		}

		successfulIDs = append(successfulIDs, id)

		deletedAuditPlan = append(deletedAuditPlan, map[string]interface{}{
			"id":    id,
			"title": existingAuditPlan.Title,
		})
	}

	if len(successfulIDs) > 0 {
		tx.Commit()

		go func() {
			for _, id := range successfulIDs {
				err := aps.activityLogSvc.LogActivity(
					models.DeleteAction,
					models.AuditPlan{}.TableName(),
					id,
					deletedAuditPlan,
					currentUserID,
					nil,
				)
				if err != nil {
					fields := map[string]interface{}{
						"action":    models.DeleteAction,
						"tableName": models.AuditPlan{}.TableName(),
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

func (aps *AuditPlanService) StatusUpdate(id uint, currentUserID uint, auditStatusID uint, tx utility.DatabaseTransaction) error {
	existingAuditPlan, err := aps.GetByID(id)
	if err != nil {
		return err
	}
	now := time.Now()

	existingAuditPlan.UpdatedAt = &now
	existingAuditPlan.AuditStatusID = auditStatusID
	existingAuditPlan.UpdatedBYID = &currentUserID

	err = aps.repo.Update(existingAuditPlan, tx)
	if err != nil {
		return err
	}

	return nil
}

func (aps *AuditPlanService) Publish(id uint, currentUserID uint) error {
	tx := aps.repo.BeginTransaction()
	existingAuditPlanById, err := aps.repo.GetByID(id)
	if err != nil {
		tx.Rollback()
		return utility.NewResourceNotFoundError(err, id)
	}
	auditStatus, err := aps.auditStatusSvc.GetByType("published")
	if err != nil {
		tx.Rollback()
		return err
	}
	if existingAuditPlanById.AuditStatusID == uint(auditStatus.ID) {
		tx.Rollback()
		return utility.NewForbiddenError("The audit draft report has already been published.")
	}
	now := time.Now()

	existingAuditPlanById.PublishedAt = &now
	existingAuditPlanById.AuditStatusID = uint(auditStatus.ID)
	existingAuditPlanById.PublishedBYID = &currentUserID

	err = aps.repo.Update(existingAuditPlanById, tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	existingAuditPlan, err := aps.repo.GetFullDetailsByID(id, tx)
	if err != nil {
		tx.Rollback()
		return utility.NewResourceNotFoundError(err, id)
	}

	tx.Commit()

	go func() {
		err := aps.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.AuditPlan{}.TableName(),
			id,
			nil,
			currentUserID,
			nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.UpdateAction,
				"tableName": models.AuditPlan{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}

		uniqueUserIDs := make(map[uint]bool)

		teamLeadDetails := "None"
		if existingAuditPlan.TeamLead != nil {
			uniqueUserIDs[existingAuditPlan.TeamLead.ID] = true
			teamLeadDetails = fmt.Sprintf("%s %s", existingAuditPlan.TeamLead.FirstNameEn, existingAuditPlan.TeamLead.LastNameEn)
		} else {
			utility.LogError("TeamLead is nil", nil)
		}

		auditMembers := []string{}
		for _, member := range existingAuditPlan.Members {
			if member.User != nil {
				uniqueUserIDs[member.User.ID] = true
				auditMembers = append(auditMembers, fmt.Sprintf("%s %s", member.User.FirstNameEn, member.User.LastNameEn))
			} else {
				utility.LogError("Member User is nil", map[string]interface{}{
					"memberID": member.ID,
				})
			}
		}

		assignedBy, err := aps.userSvc.GetByID(currentUserID)
		if err != nil {
			utility.LogError("Failed to fetch assigning user", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		auditMembersDetails := "None"
		if len(auditMembers) > 0 {
			auditMembersDetails = strings.Join(auditMembers, "\n- ")
		}

		// Send notifications
		for userID := range uniqueUserIDs {
			user, err := aps.userSvc.GetByID(userID)
			if err != nil {
				utility.LogError("Failed to fetch user details for notification", map[string]interface{}{
					"userID": userID,
					"error":  err.Error(),
				})
				continue
			}

			urlParam := fmt.Sprintf(`{"audit_plan_id":%d}`, id)
			entity := "audit-plan-publish"
			action := "responsible-user-assign"
			appURL := utility.GetAppURL()

			// DB Notification
			dbNotificationDTO := &commonDTOs.DBNotificationDTO{
				UserID:               userID,
				ModuleIdentifier:     "audit",
				NotificationType:     string(commonDTOs.InformationType),
				ReferenceNumber:      existingAuditPlan.ReferenceNumber,
				HeadingEn:            "Audit Plan Published",
				HeadingAr:            "تم نشر خطة التدقيق",
				TitleEn:              existingAuditPlan.Title,
				TitleAr:              existingAuditPlan.Title,
				Entity:               &entity,
				Action:               &action,
				URLParam:             urlParam,
				IsSystemNotification: false,
				CreatedByID:          &currentUserID,
			}

			// Email Notification
			emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
				EmailTo:              []string{user.Email},
				Subject:              fmt.Sprintf("Audit Plan Published: %s", existingAuditPlan.Title),
				TemplateRelativePath: "audit/templates/email/audit_plan_publish_create_responsible_user_mail.html",
				EmailData: map[string]interface{}{
					"FirstNameEn":           user.FirstNameEn,
					"LastNameEn":            user.LastNameEn,
					"ReferenceNumber":       existingAuditPlan.ReferenceNumber,
					"Title":                 existingAuditPlan.Title,
					"Description":           existingAuditPlan.Description,
					"Category":              existingAuditPlan.AuditCategory.TitleEn,
					"StartDate":             existingAuditPlan.StartDate.Format("2006-01-02"),
					"TargetDate":            existingAuditPlan.TargetDate.Format("2006-01-02"),
					"AssignedByFirstNameEn": assignedBy.FirstNameEn,
					"AssignedByLastNameEn":  assignedBy.LastNameEn,
					"AssignedOn":            time.Now().Format("2006-01-02"),
					"TeamLead":              teamLeadDetails,
					"AuditMembers":          auditMembersDetails,
					"RedirectURL":           fmt.Sprintf("%s/audit/audit-plans/%d/overview", appURL, id),
				},
			}

			// Send Notifications
			err = aps.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "audit")
			if err != nil {
				utility.LogError("Failed to send audit plan notifications", map[string]interface{}{
					"userID": userID,
					"error":  err.Error(),
				})
			}
		}
	}()

	return nil
}

func (aps *AuditPlanService) validateAuditPlanFormDTO(dto *auditDTOs.AuditPlanFormDTO) error {
	if _, err := aps.auditCategorySvc.GetByID(dto.AuditCategoryID); err != nil {
		return utility.NewValidationError(fmt.Errorf("audit category does not exist"), "not_exist", "audit category")
	}

	if dto.AuditTeamID != nil {
		if _, err := aps.auditTeamSvc.GetByID(*dto.AuditTeamID); err != nil {
			return utility.NewValidationError(fmt.Errorf("audit team does not exist"), "not_exist", "audit team")
		}
	}

	if dto.TeamLeadID != nil {
		if _, err := aps.userSvc.GetByID(*dto.TeamLeadID); err != nil {
			return utility.NewValidationError(fmt.Errorf("team lead does not exist"), "not_exist", "team lead")
		}
	}

	// Validate UserIds
	for _, userID := range dto.UserIds {
		if userID != nil {
			if _, err := aps.userSvc.GetByID(*userID); err != nil {
				return utility.NewValidationError(fmt.Errorf("user does not exist"), "not_exist", "user")
			}
		}
	}

	// Validate AuditableItemIds
	for _, itemID := range dto.AuditableItemIds {
		if itemID != nil {
			if _, err := aps.auditableItemSvc.GetByID(*itemID); err != nil {
				return utility.NewValidationError(fmt.Errorf("auditable item does not exist"), "not_exist", "auditable item")
			}
		}
	}

	return nil
}

func (aps *AuditPlanService) parseAuditPlanDates(dto *auditDTOs.AuditPlanFormDTO) (*auditDTOs.AuditPlanParsedDates, error) {
	parsed := &auditDTOs.AuditPlanParsedDates{}

	var err error
	if dto.StartDate != nil && *dto.StartDate != "" {
		parsed.StartDate, err = utility.ParseDate(*dto.StartDate)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid start date"), "invalid_date", "start date")
		}
	}

	if dto.TargetDate != nil && *dto.TargetDate != "" {
		parsed.TargetDate, err = utility.ParseDate(*dto.TargetDate)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid target date"), "invalid_date", "target date")
		}
	}

	return parsed, nil
}
