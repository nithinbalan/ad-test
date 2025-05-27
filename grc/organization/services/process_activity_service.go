package services

import (
	"errors"
	"fmt"
	"strings"
	"time"

	commonDTOs "jethur/grc/common/dtos"
	commonServices "jethur/grc/common/services"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	processDTOs "jethur/grc/organization/dtos"
	processRepositories "jethur/grc/organization/repositories"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type ProcessActivityServiceInterface interface {
	GetAll(processID uint, params *processDTOs.ProcessActivityListQueryParamsDTO) (*processDTOs.ProcessActivityListResponseDTO, error)
	GetByID(processID uint, id uint) (*models.ProcessActivity, error)
	GetFullDetailsByID(processID uint, id uint) (*processDTOs.ProcessActivityDetailDTO, error)
	Create(processID uint, currentUserID uint, form *processDTOs.ProcessActivityFormDTO) (uint, error)
	Update(processID uint, id uint, currentUserID uint, form *processDTOs.ProcessActivityFormDTO) error
	Delete(processID uint, ids []uint, currentUserID uint) ([]uint, error)

	UpdateProcessActivityUsers(processID, activityID, currentUserID uint, form *processDTOs.ProcessActivityUserRACIDTO) error
	DeleteActivityUser(activityUserID, processID, activityID, currentUserID uint) error
}

type ProcessActivityService struct {
	repo            processRepositories.ProcessActivityRepositoryInterface
	activityLogSvc  logServices.ActivityLogServiceInterface
	processSvc      ProcessServiceInterface
	userSvc         UserServiceInterface
	notificationSvc commonServices.NotificationServiceInterface
}

func NewProcessActivityService(
	repo processRepositories.ProcessActivityRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	processSvc ProcessServiceInterface,
	userSvc UserServiceInterface,
	notificationSvc commonServices.NotificationServiceInterface,
) ProcessActivityServiceInterface {
	return &ProcessActivityService{
		repo:            repo,
		activityLogSvc:  activityLogSvc,
		processSvc:      processSvc,
		userSvc:         userSvc,
		notificationSvc: notificationSvc,
	}
}

func (pas *ProcessActivityService) GetAll(processID uint, params *processDTOs.ProcessActivityListQueryParamsDTO) (*processDTOs.ProcessActivityListResponseDTO, error) {
	if _, err := pas.processSvc.GetByID(processID); err != nil {
		return nil, err
	}
	return pas.repo.GetAll(processID, params)
}

func (pas *ProcessActivityService) GetByID(processID uint, id uint) (*models.ProcessActivity, error) {
	if _, err := pas.processSvc.GetByID(processID); err != nil {
		return nil, err
	}
	activity, err := pas.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return activity, nil
}

func (pas *ProcessActivityService) GetFullDetailsByID(processID uint, id uint) (*processDTOs.ProcessActivityDetailDTO, error) {
	var dto processDTOs.ProcessActivityDetailDTO
	if _, err := pas.processSvc.GetByID(processID); err != nil {
		return nil, err
	}
	activity, err := pas.repo.GetFullDetailsByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	err = mapstructure.Decode(activity, &dto)
	if err != nil {
		return nil, err
	}
	dto.CreatedAt = activity.CreatedAt
	dto.UpdatedAt = activity.UpdatedAt
	return &dto, nil
}

func (pas *ProcessActivityService) Create(processID uint, currentUserID uint, form *processDTOs.ProcessActivityFormDTO) (uint, error) {
	if _, err := pas.processSvc.GetByID(processID); err != nil {
		return 0, err
	}

	activity := &models.ProcessActivity{
		ProcessID:         processID,
		TitleEn:           form.TitleEn,
		TitleAr:           form.TitleAr,
		DescriptionEn:     form.DescriptionEn,
		DescriptionAr:     form.DescriptionAr,
		Sequence:          form.Sequence,
		EstimatedDuration: form.EstimatedDuration,
		ExpectedOutput:    form.ExpectedOutput,
		CreatedBYID:       currentUserID,
	}
	tx := pas.repo.BeginTransaction()

	id, err := pas.repo.Create(activity, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	tx.Commit()
	go func() {
		existingActivity, err := pas.GetFullDetailsByID(processID, id)
		if err != nil {
			utility.LogError("Failed to fetch created activity for logging", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		err = pas.activityLogSvc.LogActivity(
			models.CreateAction,
			models.ProcessActivity{}.TableName(),
			id,
			existingActivity,
			currentUserID,
			nil,
		)
		if err != nil {
			utility.LogError("Failed to log activity", map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.ProcessActivity{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			})
		}
	}()

	return id, nil
}

func (pas *ProcessActivityService) Update(processID uint, id uint, currentUserID uint, form *processDTOs.ProcessActivityFormDTO) error {
	oldActivity, err := pas.GetFullDetailsByID(processID, id)
	if err != nil {
		return err
	}
	existing, err := pas.GetByID(processID, id)
	if err != nil {
		return err
	}

	existing.TitleEn = form.TitleEn
	existing.TitleAr = form.TitleAr
	existing.DescriptionEn = form.DescriptionEn
	existing.DescriptionAr = form.DescriptionAr
	existing.Sequence = form.Sequence
	existing.EstimatedDuration = form.EstimatedDuration
	existing.ExpectedOutput = form.ExpectedOutput
	existing.UpdatedBYID = &currentUserID

	tx := pas.repo.BeginTransaction()

	if err := pas.repo.Update(existing, tx); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	go func() {
		existingActivity, err := pas.GetFullDetailsByID(processID, id)
		if err != nil {
			utility.LogError("Failed to fetch updated activity for logging", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}
		// Build oldData and newData for comparison
		oldData := map[string]interface{}{
			"title_en":        oldActivity.TitleEn,
			"title_ar":        oldActivity.TitleAr,
			"description_en":  oldActivity.DescriptionEn,
			"description_ar":  oldActivity.DescriptionAr,
			"sequence":        oldActivity.Sequence,
			"expected_output": oldActivity.ExpectedOutput,
			"estimated_duration": func() string {
				if oldActivity.EstimatedDuration != nil {
					return *oldActivity.EstimatedDuration
				}
				return ""
			}(),
		}

		newData := map[string]interface{}{
			"title_en":        existingActivity.TitleEn,
			"title_ar":        existingActivity.TitleAr,
			"description_en":  existingActivity.DescriptionEn,
			"description_ar":  existingActivity.DescriptionAr,
			"sequence":        existingActivity.Sequence,
			"expected_output": existingActivity.ExpectedOutput,
			"estimated_duration": func() string {
				if existingActivity.EstimatedDuration != nil {
					return *existingActivity.EstimatedDuration
				}
				return ""
			}(),
		}

		// Compare and collect changed fields
		changedData := utility.CompareFieldChanges(oldData, newData)

		// Log only if any changes exist
		if len(changedData) > 0 {
			err := pas.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.ProcessActivity{}.TableName(),
				id,
				changedData,
				currentUserID,
				nil,
			)
			if err != nil {
				utility.LogError("Failed to log update activity", map[string]interface{}{
					"action":   models.UpdateAction,
					"table":    models.ProcessActivity{}.TableName(),
					"recordID": id,
					"error":    err.Error(),
				})
			}
		}
	}()

	return nil
}

func (pas *ProcessActivityService) Delete(processID uint, ids []uint, currentUserID uint) ([]uint, error) {
	if _, err := pas.processSvc.GetByID(processID); err != nil {
		return nil, err
	}

	var failed []uint
	var successful []uint
	var deletedActivities []*models.ProcessActivity

	for _, id := range ids {
		// Fetch full detail for log before delete
		activity, err := pas.repo.GetFullDetailsByID(id)
		if err != nil {
			failed = append(failed, id)
			continue
		}
		deletedActivities = append(deletedActivities, activity)

		if err := pas.repo.DeleteProcessActivityUsersByActivityID(id); err != nil {
			failed = append(failed, id)
			continue
		}
		if err := pas.repo.Delete(id); err != nil {
			failed = append(failed, id)
		} else {
			successful = append(successful, id)
		}
	}

	if len(successful) > 0 {
		go func() {
			err := pas.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.ProcessActivity{}.TableName(),
				successful,
				deletedActivities,
				currentUserID,
				nil,
			)
			if err != nil {
				utility.LogError("Failed to log bulk delete activity for process activities", map[string]interface{}{
					"action":    models.DeleteAction,
					"recordIDs": successful,
					"userID":    currentUserID,
					"error":     err.Error(),
				})
			}
		}()
	}

	return failed, nil
}
func (pas *ProcessActivityService) UpdateProcessActivityUsers(processID, activityID, currentUserID uint, form *processDTOs.ProcessActivityUserRACIDTO) error {
	if err := pas.validateProcessActivityUserRACIDTO(form); err != nil {
		return err
	}

	// Validate process and activity
	if _, err := pas.processSvc.GetByID(processID); err != nil {
		return err
	}
	oldActivity, err := pas.GetFullDetailsByID(processID, activityID)
	if err != nil {
		return err
	}
	oldUserMap := groupUserDTOsByRole(oldActivity.Users)

	// Fetch existing mappings
	existingMappings, err := pas.repo.GetActivityUsers(activityID)
	if err != nil {
		return err
	}
	tx := pas.repo.BeginTransaction()

	// Build maps
	newUserMap := map[string]map[uint]bool{
		"responsible": make(map[uint]bool),
		"accountable": make(map[uint]bool),
		"consulted":   make(map[uint]bool),
		"informed":    make(map[uint]bool),
	}
	for _, id := range form.ResponsibleUserIDs {
		newUserMap["responsible"][id] = true
	}
	for _, id := range form.AccountableUserIDs {
		newUserMap["accountable"][id] = true
	}
	for _, id := range form.ConsultedUserIDs {
		newUserMap["consulted"][id] = true
	}
	for _, id := range form.InformedUserIDs {
		newUserMap["informed"][id] = true
	}

	existingUserMap := make(map[string]map[uint]*uint)
	for _, u := range existingMappings {
		if existingUserMap[u.RoleType] == nil {
			existingUserMap[u.RoleType] = make(map[uint]*uint)
		}
		existingUserMap[u.RoleType][u.UserID] = &u.ID
	}

	// Calculate differences
	var toInsert []*models.ProcessActivityUser
	var toDelete []*uint

	for role, newIDs := range newUserMap {
		existingIDs := existingUserMap[role]

		// Insert
		for newID := range newIDs {
			if _, exists := existingIDs[newID]; !exists {
				newEntry := &models.ProcessActivityUser{
					ProcessActivityID: activityID,
					UserID:            newID,
					RoleType:          role,
				}
				if role == "responsible" {
					orgID, err := pas.userSvc.GetUserPrimaryUnitID(newID)
					if err != nil {
						tx.Rollback()
						return utility.NewValidationError(fmt.Errorf("organization unit not set for user ID %d", newID), "not_exist", "responsible unit")
					}
					newEntry.ResponsibleUnitID = &orgID
				}
				toInsert = append(toInsert, newEntry)
			}
		}

		// Delete
		for existingID, recID := range existingIDs {
			if _, still := newIDs[existingID]; !still {
				toDelete = append(toDelete, recID)
			}
		}
	}

	// Apply Deletes
	var deleteIDs []uint
	for _, ptr := range toDelete {
		if ptr != nil {
			deleteIDs = append(deleteIDs, *ptr)
		}
	}
	if len(deleteIDs) > 0 {
		if err := pas.repo.DeleteProcessActivityUsers(deleteIDs, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	// Apply Inserts
	if len(toInsert) > 0 {
		if err := pas.repo.CreateProcessActivityUsers(toInsert, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()
	go func() {
		existingActivity, err := pas.GetFullDetailsByID(processID, activityID)
		if err != nil {
			utility.LogError("Failed to fetch updated activity for logging", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}
		newUserMap := groupUserDTOsByRole(existingActivity.Users)

		changedData := make(map[string]interface{})
		roleTypes := []string{"responsible", "accountable", "consulted", "informed"}

		for _, role := range roleTypes {
			added := difference(newUserMap[role], oldUserMap[role])
			removed := difference(oldUserMap[role], newUserMap[role])

			if len(added) > 0 || len(removed) > 0 {
				changedData[role+"_user_ids"] = map[string]interface{}{
					"added":   added,
					"removed": removed,
				}
			}

			// Send notification only for responsible role
			if role == "responsible" && (len(added) > 0 || len(removed) > 0) {
				existingActivity, _ := pas.GetFullDetailsByID(processID, activityID)
				urlParam := fmt.Sprintf(`{"id":%d}`, activityID)
				entity := "process-activity"
				action := "responsible-user-assign"

				for _, uid := range added {
					user, err := pas.userSvc.GetByID(uid)
					if err != nil {
						utility.LogError("Failed to fetch responsible user", map[string]interface{}{
							"userID": uid, "error": err.Error(),
						})
						continue
					}
					emailTo := []string{user.Email}
					subject := fmt.Sprintf("Process Activity/%s: Responsibility Assigned", existingActivity.TitleEn)
					appURL := utility.GetAppURL()
					emailData := map[string]interface{}{
						"FirstNameEn":           user.FirstNameEn,
						"LastNameEn":            user.LastNameEn,
						"ProcessActivityTitle":  existingActivity.TitleEn,
						"AssignedByFirstNameEn": existingActivity.CreatedByUser.FirstNameEN,
						"AssignedByLastNameEn":  existingActivity.CreatedByUser.LastNameEN,
						"AssignedOn":            existingActivity.CreatedAt.Format("2006-01-02"),
						"RedirectURL":           fmt.Sprintf("%sorganization/process/%d/activities", appURL, processID),
					}
					emailDTO := &commonDTOs.EmailNotificationDTO{
						EmailTo:              emailTo,
						Subject:              subject,
						TemplateRelativePath: "organization/templates/email/process_activity_assign.html",
						EmailData:            emailData,
					}
					dbNotification := &commonDTOs.DBNotificationDTO{
						UserIDs:              []uint{uid},
						ModuleIdentifier:     "organization",
						NotificationType:     string(commonDTOs.AssignmentType),
						HeadingEn:            "Activity Responsibility Assigned",
						HeadingAr:            "تعيين مسؤولية النشاط",
						TitleEn:              existingActivity.TitleEn,
						TitleAr:              existingActivity.TitleAr,
						Entity:               &entity,
						Action:               &action,
						URLParam:             urlParam,
						IsSystemNotification: false,
						CreatedByID:          &currentUserID,
					}

					// Send and log
					_ = pas.notificationSvc.SendNotification(dbNotification, emailDTO, "process-activity")
				}
				// Notify removed responsible users
				for _, uid := range removed {
					user, err := pas.userSvc.GetByID(uid)
					if err != nil {
						utility.LogError("Failed to fetch removed responsible user", map[string]interface{}{
							"userID": uid, "error": err.Error(),
						})
						continue
					}

					emailTo := []string{user.Email}
					subject := fmt.Sprintf("Process Activity/%s: Responsibility Removed", existingActivity.TitleEn)
					appURL := utility.GetAppURL()
					emailData := map[string]interface{}{
						"FirstNameEn":          user.FirstNameEn,
						"LastNameEn":           user.LastNameEn,
						"ProcessActivityTitle": existingActivity.TitleEn,
						"RemovedOn":            time.Now().Format("2006-01-02"),
						"RedirectURL":          fmt.Sprintf("%sorganization/process/%d/activities", appURL, processID),
					}

					emailDTO := &commonDTOs.EmailNotificationDTO{
						EmailTo:              emailTo,
						Subject:              subject,
						TemplateRelativePath: "organization/templates/email/process_activity_unassign.html",
						EmailData:            emailData,
					}
					dbNotification := &commonDTOs.DBNotificationDTO{
						UserIDs:              []uint{uid},
						ModuleIdentifier:     "organization",
						NotificationType:     string(commonDTOs.AssignmentType),
						HeadingEn:            "Activity Responsibility Removed",
						HeadingAr:            "تمت إزالة مسؤولية النشاط",
						TitleEn:              existingActivity.TitleEn,
						TitleAr:              existingActivity.TitleAr,
						Entity:               &entity,
						Action:               &action,
						URLParam:             urlParam,
						IsSystemNotification: false,
						CreatedByID:          &currentUserID,
					}

					_ = pas.notificationSvc.SendNotification(dbNotification, emailDTO, "process-activity")
				}
			}

		}

		if len(changedData) > 0 {
			err := pas.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.ProcessActivity{}.TableName(),
				activityID,
				changedData,
				currentUserID,
				nil,
			)
			if err != nil {
				utility.LogError("Failed to log update activity", map[string]interface{}{
					"action":   models.UpdateAction,
					"table":    models.ProcessActivity{}.TableName(),
					"recordID": activityID,
					"error":    err.Error(),
				})
			}
		}
	}()

	return nil
}

func (pas *ProcessActivityService) DeleteActivityUser(activityUserID, processID, activityID, currentUserID uint) error {
	existing, err := pas.repo.GetActivityUserByID(activityUserID)
	if err != nil {
		return err
	}
	if existing == nil {
		return utility.NewValidationError(fmt.Errorf("activity user not set for user ID %d", activityUserID), "not_exist", "activity user")
	}

	// Fetch old user map for comparison before delete
	oldActivity, err := pas.GetFullDetailsByID(processID, activityID)
	if err != nil {
		return err
	}
	oldUserMap := groupUserDTOsByRole(oldActivity.Users)

	// Delete the user
	if err := pas.repo.DeleteProcessActivityUser(activityUserID); err != nil {
		return err
	}

	// Fire async activity log comparison
	go func() {
		existingActivity, err := pas.GetFullDetailsByID(processID, activityID)
		if err != nil {
			utility.LogError("Failed to fetch updated activity for logging", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}
		newUserMap := groupUserDTOsByRole(existingActivity.Users)

		changedData := make(map[string]interface{})
		roleTypes := []string{"responsible", "accountable", "consulted", "informed"}
		entity := "process-activity"
		action := "responsible-user-remove"
		for _, role := range roleTypes {
			added := difference(newUserMap[role], oldUserMap[role])
			removed := difference(oldUserMap[role], newUserMap[role])

			if len(added) > 0 || len(removed) > 0 {
				changedData[role+"_user_ids"] = map[string]interface{}{
					"added":   added,
					"removed": removed,
				}
			}
		}

		if len(changedData) > 0 {
			err := pas.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.ProcessActivity{}.TableName(),
				activityID,
				changedData,
				currentUserID,
				nil,
			)
			if err != nil {
				utility.LogError("Failed to log update activity", map[string]interface{}{
					"action":   models.UpdateAction,
					"table":    models.ProcessActivity{}.TableName(),
					"recordID": activityID,
					"error":    err.Error(),
				})
			}
		}

		removed := difference(oldUserMap["responsible"], newUserMap["responsible"])
		for _, uid := range removed {
			user, err := pas.userSvc.GetByID(uid)
			if err != nil {
				utility.LogError("Failed to fetch removed responsible user", map[string]interface{}{
					"userID": uid, "error": err.Error(),
				})
				continue
			}

			emailTo := []string{user.Email}
			subject := fmt.Sprintf("Process Activity/%s: Responsibility Removed", existingActivity.TitleEn)
			appURL := utility.GetAppURL()
			emailData := map[string]interface{}{
				"FirstNameEn":          user.FirstNameEn,
				"LastNameEn":           user.LastNameEn,
				"ProcessActivityTitle": existingActivity.TitleEn,
				"RemovedOn":            time.Now().Format("2006-01-02"),
				"RedirectURL":          fmt.Sprintf("%sorganization/process/%d/activities", appURL, processID),
			}

			emailDTO := &commonDTOs.EmailNotificationDTO{
				EmailTo:              emailTo,
				Subject:              subject,
				TemplateRelativePath: "organization/templates/email/process_activity_unassign.html",
				EmailData:            emailData,
			}

			dbNotification := &commonDTOs.DBNotificationDTO{
				UserIDs:              []uint{uid},
				ModuleIdentifier:     "organization",
				NotificationType:     string(commonDTOs.AssignmentType),
				HeadingEn:            "Activity Responsibility Removed",
				HeadingAr:            "تمت إزالة مسؤولية النشاط",
				TitleEn:              existingActivity.TitleEn,
				TitleAr:              existingActivity.TitleAr,
				Entity:               &entity,
				Action:               &action,
				URLParam:             fmt.Sprintf(`{"id":%d}`, processID),
				IsSystemNotification: false,
				CreatedByID:          &currentUserID,
			}

			_ = pas.notificationSvc.SendNotification(dbNotification, emailDTO, "process-activity")
		}
	}()

	return nil
}

func groupUserDTOsByRole(users []*processDTOs.ProcessActivityUserDTO) map[string][]uint {
	roleMap := map[string][]uint{}
	for _, u := range users {
		role := strings.ToLower(u.RoleType)
		roleMap[role] = append(roleMap[role], u.UserID)
	}
	return roleMap
}

func difference(slice1, slice2 []uint) []uint {
	diff := []uint{}
	exists := make(map[uint]bool)

	for _, val := range slice2 {
		exists[val] = true
	}

	for _, val := range slice1 {
		if !exists[val] {
			diff = append(diff, val)
		}
	}

	return diff
}

func (pas *ProcessActivityService) validateProcessActivityUserRACIDTO(form *processDTOs.ProcessActivityUserRACIDTO) error {
	userIDSet := make(map[uint]bool)
	allUserIDs := append(append(form.ResponsibleUserIDs, form.AccountableUserIDs...), form.ConsultedUserIDs...)

	for _, userID := range allUserIDs {
		if userIDSet[userID] {
			return utility.NewValidationError(fmt.Errorf("duplicate user ID: %d", userID), "duplicate", "user_ids")
		}
		if _, err := pas.userSvc.GetByID(userID); err != nil {
			return utility.NewValidationError(fmt.Errorf("user does not exist for ID: %d", userID), "not_exist", "user_ids")
		}
		userIDSet[userID] = true
	}

	return nil
}
