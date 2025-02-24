package services

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"strings"
	"time"

	commonDTOs "jethur/grc/common/dtos"
	commonServices "jethur/grc/common/services"
	logServices "jethur/grc/logs/services"
	masterServices "jethur/grc/masters/services"
	"jethur/grc/models"
	organizationDTOs "jethur/grc/organization/dtos"
	organizationRepositories "jethur/grc/organization/repositories"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type ControlServiceInterface interface {
	GetAll(params *organizationDTOs.ControlListQueryParamsDTO) (*organizationDTOs.ControlListResponseDTO, error)
	GetAllControls() ([]*models.Control, error)
	GetByID(id uint, tx ...utility.DatabaseTransaction) (*models.Control, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*organizationDTOs.ControlDetailDTO, error)
	Create(currentUserID uint, controlFormDTO *organizationDTOs.ControlFormDTO) (uint, error)
	Update(id uint, currentUserID uint, controlFormDTO *organizationDTOs.ControlFormDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)

	Archive(ids []uint, currentUserID uint) ([]uint, error)
	Restore(ids []uint, currentUserID uint) ([]uint, error)
	TestStatusUpdate(control *models.Control, tx utility.DatabaseTransaction) error
	ReviewStatusUpdate(control *models.Control, tx utility.DatabaseTransaction) error

	ExportControls(data []*organizationDTOs.ControlListDTO, currentUserID uint) (*os.File, error)
	ImportProcesses(file multipart.File, currentUserID uint) ([]int, error)

	UpdateControlRiskRating(control *models.Control, tx utility.DatabaseTransaction) error
	MarkControlAsApplicable(controlID uint, currentUserID uint, controlApplicableFormDTO *organizationDTOs.ControlApplicableFormDTO) error
	MarkControlAsNotApplicable(controlID uint, currentUserID uint, controlApplicableFormDTO *organizationDTOs.ControlApplicableFormDTO) error

	GenerateControlTemplate() (*os.File, error)
	ValidateControlTemplate(file io.Reader) ([]map[string]interface{}, error)
	ValidateImportControls(input []map[string]interface{}) ([]map[string]interface{}, error)
	ImportControls(input []map[string]interface{}, currentUserID uint) ([]uint, error)
}

type ControlService struct {
	repo                           organizationRepositories.ControlRepositoryInterface
	activityLogSvc                 logServices.ActivityLogServiceInterface
	controlTypeSvc                 ControlTypeServiceInterface
	controlCategorySvc             ControlCategoryServiceInterface
	controlImplementationStatusSvc ControlImplementationStatusServiceInterface
	controlEfficiencyStatusSvc     ControlEfficiencyStatusServiceInterface
	userSvc                        UserServiceInterface
	reviewFrequencySvc             masterServices.ReviewFrequencyServiceInterface
	lifecycleStatusSvc             masterServices.LifecycleStatusServiceInterface
	notificationSvc                commonServices.NotificationServiceInterface
	userSettingSvc                 UserSettingServiceInterface
}

func NewControlService(
	repo organizationRepositories.ControlRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	controlTypeSvc ControlTypeServiceInterface,
	controlCategorySvc ControlCategoryServiceInterface,
	controlImplementationStatusSvc ControlImplementationStatusServiceInterface,
	controlEfficiencyStatusSvc ControlEfficiencyStatusServiceInterface,
	userSvc UserServiceInterface,
	reviewFrequencySvc masterServices.ReviewFrequencyServiceInterface,
	lifecycleStatusSvc masterServices.LifecycleStatusServiceInterface,
	notificationSvc commonServices.NotificationServiceInterface,
	userSettingSvc UserSettingServiceInterface,
) ControlServiceInterface {
	return &ControlService{
		repo:                           repo,
		activityLogSvc:                 activityLogSvc,
		controlTypeSvc:                 controlTypeSvc,
		controlCategorySvc:             controlCategorySvc,
		controlImplementationStatusSvc: controlImplementationStatusSvc,
		controlEfficiencyStatusSvc:     controlEfficiencyStatusSvc,
		userSvc:                        userSvc,
		reviewFrequencySvc:             reviewFrequencySvc,
		lifecycleStatusSvc:             lifecycleStatusSvc,
		notificationSvc:                notificationSvc,
		userSettingSvc:                 userSettingSvc,
	}
}

func (cs *ControlService) GetAll(params *organizationDTOs.ControlListQueryParamsDTO) (*organizationDTOs.ControlListResponseDTO, error) {
	return cs.repo.GetAll(params)
}

func (cs *ControlService) GetAllControls() ([]*models.Control, error) {
	return cs.repo.GetAllControls()
}

func (cs *ControlService) GetByID(id uint, tx ...utility.DatabaseTransaction) (*models.Control, error) {
	control, err := cs.repo.GetByID(id, tx...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return control, nil
}

func (cs *ControlService) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*organizationDTOs.ControlDetailDTO, error) {
	var controlDetailDTO organizationDTOs.ControlDetailDTO

	control, err := cs.repo.GetFullDetailsByID(id, tx...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	err = mapstructure.Decode(control, &controlDetailDTO)
	if err != nil {
		return nil, err
	}

	if control.LastReviewDate != nil {
		lastReviewDate := control.LastReviewDate.Format("2006-01-02")
		controlDetailDTO.LastReviewDateString = &lastReviewDate
	}
	if control.NextReviewDate != nil {
		nextReviewDate := control.NextReviewDate.Format("2006-01-02")
		controlDetailDTO.NextReviewDateString = &nextReviewDate
	}
	if control.LastTestDate != nil {
		lastTestDate := control.LastTestDate.Format("2006-01-02")
		controlDetailDTO.LastTestDateString = &lastTestDate
	}
	if control.NextTestDate != nil {
		nextTestDateString := control.NextTestDate.Format("2006-01-02")
		controlDetailDTO.NextTestDateString = &nextTestDateString
	}
	if control.ImplementedAt != nil {
		implementedAtString := control.ImplementedAt.Format("2006-01-02")
		controlDetailDTO.ImplementedAtString = &implementedAtString
	}
	if control.ApplicabilityChangeDate != nil {
		applicabilityChangeDateString := control.ApplicabilityChangeDate.Format("2006-01-02")
		controlDetailDTO.ApplicabilityChangeDateString = &applicabilityChangeDateString
	}

	// todo :- also add extra fields which required in details screen

	return &controlDetailDTO, nil
}

func (cs *ControlService) Create(currentUserID uint, controlFormDTO *organizationDTOs.ControlFormDTO) (uint, error) {
	if err := cs.validateControlFormDTO(controlFormDTO); err != nil {
		return 0, err
	}

	dates, err := cs.parseControlDates(controlFormDTO)
	if err != nil {
		return 0, err
	}

	referenceNumber, err := cs.repo.GenerateReferenceNumber()
	if err != nil {
		return 0, err
	}

	lifecycleStatus, err := cs.lifecycleStatusSvc.GetByType("active")
	if err != nil {
		return 0, err
	}
	organizationUnitID, err := cs.userSvc.GetUserPrimaryUnitID(controlFormDTO.OwnerID)
	if err != nil {
		return 0, utility.NewValidationError(fmt.Errorf("organization unit not set for owner"), "not_exist", "control owner unit")
	}

	control := &models.Control{
		ControlTypeID:                 controlFormDTO.ControlTypeID,
		ControlCategoryID:             controlFormDTO.ControlCategoryID,
		OwnerID:                       controlFormDTO.OwnerID,
		OwnerUnitID:                   organizationUnitID,
		ControlImplementationStatusID: controlFormDTO.ControlImplementationStatusID,
		ControlEfficiencyStatusID:     controlFormDTO.ControlEfficiencyStatusID,
		ReferenceNumber:               referenceNumber,
		Code:                          controlFormDTO.Code,
		TitleEn:                       controlFormDTO.TitleEn,
		TitleAr:                       controlFormDTO.TitleAr,
		Description:                   controlFormDTO.Description,
		Objective:                     controlFormDTO.Objective,
		EfficiencyScore:               controlFormDTO.EfficiencyScore,
		ReviewFrequencyID:             controlFormDTO.ReviewFrequencyID,
		LastReviewDate:                dates.LastReviewDate,
		TestResponsibleUserID:         controlFormDTO.TestResponsibleUserID,
		NextReviewDate:                dates.NextReviewDate,
		LastReviewNote:                controlFormDTO.LastReviewNote,
		TestFrequencyID:               controlFormDTO.TestFrequencyID,
		ImplementedAt:                 dates.ImplementedAt,
		LastTestDate:                  dates.LastTestDate,
		NextTestDate:                  dates.NextTestDate,
		LifecycleStatusID:             lifecycleStatus.ID,
		CreatedByID:                   currentUserID,
		IsApplicable:                  controlFormDTO.IsApplicable,
		Justification:                 controlFormDTO.Justification,
		ApplicabilityChangeDate:       dates.ApplicabilityChangeDate,
		ImplementationGuideline:       controlFormDTO.ImplementationGuideline,
	}

	id, err := cs.repo.Create(control)
	if err != nil {
		return 0, err
	}

	go func() {
		urlParam := fmt.Sprintf(`{"id":%d}`, id)
		entity := "control"
		action := "responsible-user-assign"

		// Send notification to new owner
		dbNotificationDTO := &commonDTOs.DBNotificationDTO{
			UserIDs:              []uint{controlFormDTO.OwnerID},
			ModuleIdentifier:     "organization",
			NotificationType:     string(commonDTOs.AssignmentType),
			ReferenceNumber:      referenceNumber,
			HeadingEn:            "Control Responsibility Assigned",
			HeadingAr:            "مسؤولية التحكم المعينة",
			TitleEn:              controlFormDTO.TitleEn,
			TitleAr:              controlFormDTO.TitleAr,
			Entity:               &entity,
			Action:               &action,
			URLParam:             urlParam,
			IsSystemNotification: false,
			CreatedByID:          &currentUserID,
		}

		owner, err := cs.userSvc.GetByID(controlFormDTO.OwnerID)
		if err != nil {
			utility.LogError("Failed to fetch responsible user", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		createdBy, err := cs.userSvc.GetByID(currentUserID)
		if err != nil {
			utility.LogError("Failed to fetch responsible user", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		emailTo := []string{owner.Email}
		subject := fmt.Sprintf("Control/%s: Responsibility Assigned", control.TitleEn)
		appURL := utility.GetAppURL()

		emailData := map[string]interface{}{
			"FirstNameEn":          owner.FirstNameEn,
			"LastNameEn":           owner.LastNameEn,
			"ReferenceNumber":      control.ReferenceNumber,
			"TitleEn":              control.TitleEn,
			"TitleAr":              control.TitleAr,
			"Description":          control.Description,
			"CreatedByFirstNameEn": createdBy.FirstNameEn,
			"CreatedByLastNameEn":  createdBy.LastNameEn,
			"AssignedOn":           time.Now().Format("2006-01-02"),
			"RedirectURL":          fmt.Sprintf("%scompliance/controls/%d/overview", appURL, id),
		}

		emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
			EmailTo:              emailTo,
			Subject:              subject,
			TemplateRelativePath: "organization/templates/email/control_create_owner_mail.html",
			EmailData:            emailData,
		}

		err = cs.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "control")
		if err != nil {
			utility.LogError("Failed to send notifications", map[string]interface{}{
				"error": err.Error(),
			})
		}
	}()

	go func() {
		err := cs.activityLogSvc.LogActivity(
			models.CreateAction,
			models.Control{}.TableName(),
			id,
			controlFormDTO,
			currentUserID,
			nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.Control{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log control activity", fields)
		}
	}()

	return id, nil
}

func (cs *ControlService) Update(id uint, currentUserID uint, controlFormDTO *organizationDTOs.ControlFormDTO) error {
	existingControl, err := cs.GetByID(id)
	if err != nil {
		return err
	}

	oldResponsibleUserID := existingControl.OwnerID

	if err := cs.validateControlFormDTO(controlFormDTO); err != nil {
		return err
	}

	dates, err := cs.parseControlDates(controlFormDTO)
	if err != nil {
		return err
	}

	if existingControl.OwnerID != controlFormDTO.OwnerID {
		organizationUnitID, err := cs.userSvc.GetUserPrimaryUnitID(controlFormDTO.OwnerID)
		if err != nil {
			return utility.NewValidationError(fmt.Errorf("organization unit not set for owner"), "not_exist", "control owner unit")
		}
		existingControl.OwnerUnitID = organizationUnitID
	}
	existingControl.ControlTypeID = controlFormDTO.ControlTypeID
	existingControl.ControlCategoryID = controlFormDTO.ControlCategoryID
	existingControl.OwnerID = controlFormDTO.OwnerID
	existingControl.ControlImplementationStatusID = controlFormDTO.ControlImplementationStatusID
	existingControl.ControlEfficiencyStatusID = controlFormDTO.ControlEfficiencyStatusID
	existingControl.Code = controlFormDTO.Code
	existingControl.TitleEn = controlFormDTO.TitleEn
	existingControl.TitleAr = controlFormDTO.TitleAr
	existingControl.Description = controlFormDTO.Description
	existingControl.Objective = controlFormDTO.Objective
	existingControl.ReviewFrequencyID = controlFormDTO.ReviewFrequencyID
	existingControl.LastReviewDate = dates.LastReviewDate
	existingControl.TestResponsibleUserID = controlFormDTO.TestResponsibleUserID
	existingControl.NextReviewDate = dates.NextReviewDate
	existingControl.LastReviewNote = controlFormDTO.LastReviewNote
	existingControl.TestFrequencyID = controlFormDTO.TestFrequencyID
	existingControl.ImplementedAt = dates.ImplementedAt
	existingControl.LastTestDate = dates.LastTestDate
	existingControl.NextTestDate = dates.NextTestDate
	existingControl.IsApplicable = controlFormDTO.IsApplicable
	existingControl.Justification = controlFormDTO.Justification
	existingControl.ImplementationGuideline = controlFormDTO.ImplementationGuideline
	existingControl.ApplicabilityChangeDate = dates.ApplicabilityChangeDate
	existingControl.UpdatedByID = &currentUserID

	err = cs.repo.Update(existingControl)
	if err != nil {
		return err
	}

	// If the owner has changed, send notifications to the old and new owner
	go func() {
		if oldResponsibleUserID != controlFormDTO.OwnerID {
			urlParam := fmt.Sprintf(`{"id":%d}`, id)
			entity := "control"
			removeAction := "responsible-user-remove"
			assignAction := "responsible-user-assign"

			// Notify the old owner
			oldResponsibleUserNotificationDTO := &commonDTOs.DBNotificationDTO{
				UserIDs:              []uint{oldResponsibleUserID},
				ModuleIdentifier:     "organization",
				NotificationType:     string(commonDTOs.InformationType),
				ReferenceNumber:      existingControl.ReferenceNumber,
				HeadingEn:            "Control Responsibility Removed",
				HeadingAr:            "تمت إزالة مسؤولية التحكم",
				TitleEn:              controlFormDTO.TitleEn,
				TitleAr:              controlFormDTO.TitleAr,
				Entity:               &entity,
				Action:               &removeAction,
				URLParam:             urlParam,
				IsSystemNotification: false,
				CreatedByID:          &currentUserID,
			}

			// Notify the new owner
			newResponsibleUserNotificationDTO := &commonDTOs.DBNotificationDTO{
				UserIDs:              []uint{controlFormDTO.OwnerID},
				ModuleIdentifier:     "organization",
				NotificationType:     string(commonDTOs.AssignmentType),
				ReferenceNumber:      existingControl.ReferenceNumber,
				HeadingEn:            "Control Responsibility Assigned",
				HeadingAr:            "مسؤولية التحكم المعينة",
				TitleEn:              controlFormDTO.TitleEn,
				TitleAr:              controlFormDTO.TitleAr,
				Entity:               &entity,
				Action:               &assignAction,
				URLParam:             urlParam,
				IsSystemNotification: false,
				CreatedByID:          &currentUserID,
			}

			// Fetch the responsible users' data
			oldResponsibleUser, err := cs.userSvc.GetByID(oldResponsibleUserID)
			if err != nil {
				utility.LogError("Failed to fetch old responsible user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}

			newResponsibleUser, err := cs.userSvc.GetByID(controlFormDTO.OwnerID)
			if err != nil {
				utility.LogError("Failed to fetch new responsible user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}

			// Fetch details of the user making the assignment
			updatedByUser, err := cs.userSvc.GetByID(currentUserID)
			if err != nil {
				utility.LogError("Failed to fetch assigning user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}

			appURL := utility.GetAppURL()

			// Email notification for the old responsible user
			oldResponsibleUserEmailTo := []string{oldResponsibleUser.Email}
			oldResponsibleUserSubject := fmt.Sprintf("Control/%s: Responsibility Removed", existingControl.TitleEn)
			oldResponsibleUserEmailData := map[string]interface{}{
				"FirstNameEn":     oldResponsibleUser.FirstNameEn,
				"LastNameEn":      oldResponsibleUser.LastNameEn,
				"ReferenceNumber": existingControl.ReferenceNumber,
				"TitleEn":         existingControl.TitleEn,
				"TitleAr":         existingControl.TitleAr,
				"Description":     existingControl.Description,
				"RedirectURL":     fmt.Sprintf("%scompliance/controls/%d/overview", appURL, id),
			}

			oldResponsibleUserEmailNotificationDTO := &commonDTOs.EmailNotificationDTO{
				EmailTo:              oldResponsibleUserEmailTo,
				Subject:              oldResponsibleUserSubject,
				TemplateRelativePath: "organization/templates/email/control_remove_owner_mail.html",
				EmailData:            oldResponsibleUserEmailData,
			}

			// Email notification for the new responsible user
			newResponsibleUserEmailTo := []string{newResponsibleUser.Email}
			newResponsibleUserSubject := fmt.Sprintf("Control/%s: Responsibility Assigned", existingControl.TitleEn)
			newResponsibleUserEmailData := map[string]interface{}{
				"FirstNameEn":           newResponsibleUser.FirstNameEn,
				"LastNameEn":            newResponsibleUser.LastNameEn,
				"ReferenceNumber":       existingControl.ReferenceNumber,
				"TitleEn":               existingControl.TitleEn,
				"TitleAr":               existingControl.TitleAr,
				"Description":           existingControl.Description,
				"AssignedByFirstNameEn": updatedByUser.FirstNameEn,
				"AssignedByLastNameEn":  updatedByUser.LastNameEn,
				"AssignedOn":            time.Now().Format("2006-01-02"),
				"RedirectURL":           fmt.Sprintf("%scompliance/controls/%d/overview", appURL, id),
			}

			newResponsibleUserEmailNotificationDTO := &commonDTOs.EmailNotificationDTO{
				EmailTo:              newResponsibleUserEmailTo,
				Subject:              newResponsibleUserSubject,
				TemplateRelativePath: "organization/templates/email/control_change_owner_mail.html",
				EmailData:            newResponsibleUserEmailData,
			}

			// Send notifications and emails
			err = cs.notificationSvc.SendNotification(newResponsibleUserNotificationDTO, newResponsibleUserEmailNotificationDTO, "control")
			if err != nil {
				utility.LogError("Failed to send notifications to new responsible user", map[string]interface{}{
					"error": err.Error(),
				})
			}

			err = cs.notificationSvc.SendNotification(oldResponsibleUserNotificationDTO, oldResponsibleUserEmailNotificationDTO, "control")
			if err != nil {
				utility.LogError("Failed to send notifications to old responsible user", map[string]interface{}{
					"error": err.Error(),
				})
			}

		}
	}()

	return nil
}

func (cs *ControlService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIds []uint
	var successfulIds []uint
	var notifications []struct {
		OwnerID         *uint
		TitleEn         string
		TitleAr         string
		Description     string
		ReferenceNumber string
	}

	for _, id := range ids {
		control, err := cs.GetFullDetailsByID(id)
		if err != nil {
			failedIds = append(failedIds, id)
			continue
		}

		// Prepare notification data before deletion
		notifications = append(notifications, struct {
			OwnerID         *uint
			TitleEn         string
			TitleAr         string
			Description     string
			ReferenceNumber string
		}{
			OwnerID:         control.OwnerID,
			TitleEn:         control.TitleEn,
			TitleAr:         control.TitleAr,
			Description:     *control.Description,
			ReferenceNumber: control.ReferenceNumber,
		})

		// Perform deletion
		err = cs.repo.Delete(id)
		if err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}
	}

	if len(successfulIds) > 0 {
		go func() {
			// Log bulk delete activities
			err := cs.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.Control{}.TableName(),
				successfulIds,
				nil,
				currentUserID,
				nil,
			)
			if err != nil {
				fmt.Printf("Failed to log bulk delete activity: %v\n", err)
			}

			// Send notifications to the owners of the deleted controls
			for _, notification := range notifications {
				urlParam := ""
				queryParam := ""
				entity := "control"
				action := "delete"

				dbNotificationDTO := &commonDTOs.DBNotificationDTO{
					UserIDs:              []uint{*notification.OwnerID},
					ModuleIdentifier:     "organization",
					NotificationType:     string(commonDTOs.DeletionType),
					ReferenceNumber:      notification.ReferenceNumber,
					HeadingEn:            "Control Deleted",
					HeadingAr:            "تم حذف التحكم",
					TitleEn:              notification.TitleEn,
					TitleAr:              notification.TitleAr,
					Entity:               &entity,
					Action:               &action,
					URLParam:             urlParam,
					QueryParam:           queryParam,
					IsSystemNotification: false,
					CreatedByID:          &currentUserID,
				}

				ownerUser, err := cs.userSvc.GetByID(*notification.OwnerID)
				if err != nil {
					utility.LogError("Failed to fetch owner details for notification", map[string]interface{}{
						"error": err.Error(),
					})
					continue
				}

				// Prepare email notification
				ownerEmailTo := []string{ownerUser.Email}
				subject := fmt.Sprintf("Control/%s: Deleted", notification.TitleEn)
				templateRelativePath := "organization/templates/email/control_delete_mail.html"

				emailData := map[string]interface{}{
					"FirstNameEn":     ownerUser.FirstNameEn,
					"LastNameEn":      ownerUser.LastNameEn,
					"ReferenceNumber": notification.ReferenceNumber,
					"TitleEn":         notification.TitleEn,
					"TitleAr":         notification.TitleAr,
					"Description":     notification.Description,
				}

				emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
					EmailTo:              ownerEmailTo,
					Subject:              subject,
					TemplateRelativePath: templateRelativePath,
					EmailData:            emailData,
				}

				// Send DB and email notifications
				err = cs.notificationSvc.SendNotification(
					dbNotificationDTO, emailNotificationDTO, "control")
				if err != nil {
					utility.LogError("Failed to send deletion notifications", map[string]interface{}{
						"error": err.Error(),
					})
				}
			}
		}()
	}

	return failedIds, nil
}

func (cs *ControlService) Archive(ids []uint, currentUserID uint) ([]uint, error) {
	lifecycleStatus, err := cs.lifecycleStatusSvc.GetByType("archived")
	if err != nil {
		return nil, err
	}

	idsToUpdate, err := cs.repo.GetIDsByLifecycleStatus(ids, "active")
	if err != nil {
		return nil, err
	}

	successfulIds, err := cs.repo.UpdateLifecycleStatusForIDs(idsToUpdate, lifecycleStatus.ID)
	if err != nil {
		return nil, err
	}

	if len(successfulIds) > 0 {
		go func() {
			err := cs.activityLogSvc.LogBulkActivities(
				models.ArchiveAction,
				models.Control{}.TableName(),
				successfulIds,
				nil,
				currentUserID,
				nil,
			)
			if err != nil {
				// @todo - save to log file
				fmt.Printf("Failed to log bulk archive activity: %v\n", err)
			}
		}()
	}

	return successfulIds, nil
}

func (cs *ControlService) Restore(ids []uint, currentUserID uint) ([]uint, error) {
	lifecycleStatus, err := cs.lifecycleStatusSvc.GetByType("active")
	if err != nil {
		return nil, err
	}

	idsToUpdate, err := cs.repo.GetIDsByLifecycleStatus(ids, "archived")
	if err != nil {
		return nil, err
	}

	successfulIds, err := cs.repo.UpdateLifecycleStatusForIDs(idsToUpdate, lifecycleStatus.ID)
	if err != nil {
		return nil, err
	}

	if len(successfulIds) > 0 {
		go func() {
			err := cs.activityLogSvc.LogBulkActivities(
				models.RestoreAction,
				models.Control{}.TableName(),
				successfulIds,
				nil,
				currentUserID,
				nil,
			)
			if err != nil {
				// @todo - save to log file
				fmt.Printf("Failed to log bulk archive activity: %v\n", err)
			}
		}()
	}

	return successfulIds, nil
}

func (cs *ControlService) TestStatusUpdate(control *models.Control, tx utility.DatabaseTransaction) error {
	err := cs.repo.Update(control, tx)
	if err != nil {
		return err
	}

	return nil
}

func (cs *ControlService) ReviewStatusUpdate(control *models.Control, tx utility.DatabaseTransaction) error {
	err := cs.repo.Update(control, tx)
	if err != nil {
		return err
	}

	return nil
}

func (cs *ControlService) validateControlFormDTO(controlFormDTO *organizationDTOs.ControlFormDTO) error {
	if _, err := cs.controlTypeSvc.GetByID(uint(controlFormDTO.ControlTypeID)); err != nil {
		return utility.NewValidationError(fmt.Errorf("control type does not exist"), "not_exist", "control type")
	}

	if controlFormDTO.ControlCategoryID != nil {
		if _, err := cs.controlCategorySvc.GetByID(uint(*controlFormDTO.ControlCategoryID)); err != nil {
			return utility.NewValidationError(fmt.Errorf("control category does not exist"), "not_exist", "control category")
		}
	}

	if _, err := cs.userSvc.GetByID(uint(controlFormDTO.OwnerID)); err != nil {
		return utility.NewValidationError(fmt.Errorf("control owner does not exist"), "not_exist", "control owner")
	}

	if controlFormDTO.ControlImplementationStatusID != nil {
		if _, err := cs.controlImplementationStatusSvc.GetByID(uint(*controlFormDTO.ControlImplementationStatusID)); err != nil {
			return utility.NewValidationError(fmt.Errorf("control implementation status does not exist"), "not_exist", "control implementation status")
		}
	}

	if controlFormDTO.ControlEfficiencyStatusID != nil {
		if _, err := cs.controlImplementationStatusSvc.GetByID(uint(*controlFormDTO.ControlEfficiencyStatusID)); err != nil {
			return utility.NewValidationError(fmt.Errorf("control efficiency status does not exist"), "not_exist", "control efficiency status")
		}
	}

	if controlFormDTO.ReviewFrequencyID != nil {
		if _, err := cs.reviewFrequencySvc.GetByID(uint(*controlFormDTO.ReviewFrequencyID)); err != nil {
			return utility.NewValidationError(fmt.Errorf("review frequiency does not exist"), "not_exist", "review frequiency")
		}
	}

	if _, err := cs.userSvc.GetByID(uint(controlFormDTO.TestResponsibleUserID)); err != nil {
		return utility.NewValidationError(fmt.Errorf("control test responsible user does not exist"), "not_exist", "test responsible user")
	}

	if controlFormDTO.TestFrequencyID != nil {
		if _, err := cs.reviewFrequencySvc.GetByID(uint(*controlFormDTO.TestFrequencyID)); err != nil {
			return utility.NewValidationError(fmt.Errorf("review test frequiency does not exist"), "not_exist", "test frequiency")
		}
	}

	if controlFormDTO.IsApplicable {
		if controlFormDTO.ControlImplementationStatusID == nil {
			return utility.NewValidationError(fmt.Errorf("control implementation status is required when control is applicable"), "required", "control implementation status id")
		}
		if controlFormDTO.ControlEfficiencyStatusID == nil {
			return utility.NewValidationError(fmt.Errorf("control efficiency status is required when control is applicable"), "required", "control efficiency status id")
		}
	}

	return nil
}

func (cs *ControlService) parseControlDates(dto *organizationDTOs.ControlFormDTO) (*organizationDTOs.ControlParsedDates, error) {
	parsed := &organizationDTOs.ControlParsedDates{}

	var err error
	if dto.LastReviewDate != nil {
		parsed.LastReviewDate, err = utility.ParseDate(*dto.LastReviewDate)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid last review date"), "invalid_date", "last review date")
		}
	}

	if dto.NextReviewDate != nil {
		parsed.NextReviewDate, err = utility.ParseDate(*dto.NextReviewDate)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid next review date"), "invalid_date", "next review date")
		}
	}

	if dto.LastTestDate != nil {
		parsed.LastTestDate, err = utility.ParseDate(*dto.LastTestDate)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid last test date"), "invalid_date", "last test date")
		}
	}

	if dto.NextTestDate != nil {
		parsed.NextTestDate, err = utility.ParseDate(*dto.NextTestDate)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid next test date"), "invalid_date", "next test date")
		}
	}

	if dto.ImplementedAt != nil {
		parsed.ImplementedAt, err = utility.ParseDate(*dto.ImplementedAt)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid implemented at date"), "invalid_date", "implemented at date")
		}
	}

	if dto.ApplicabilityChangeDate != nil {
		parsed.ApplicabilityChangeDate, err = utility.ParseDate(*dto.ApplicabilityChangeDate)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid applicability change date"), "invalid_date", "applicability change date")
		}
	}

	return parsed, nil
}

func (cs *ControlService) ExportControls(data []*organizationDTOs.ControlListDTO, currentUserID uint) (*os.File, error) {
	// Fetch user setting for language preference
	userSetting, err := cs.userSettingSvc.GetUserSettingsByUserID(currentUserID)
	if err != nil {
		return nil, err
	}

	fmt.Println(userSetting.Language)
	// Define headers based on user language setting
	var headers []string
	if userSetting.Language == "ar" {
		headers = []string{
			"SL No.",
			"Reference Number",
			"Control Code",
			"Title",
			"Description",
			"Objective",
			"Control Category",
			"Control Type",
			"Owner",
			"Review Frequency",
			"Implementation Status",
			"Implementation Date",
			"Efficiency Status",
			"Last Test Date",
			"Test Frequency",
			"Next Test Date",
			"Last Review Date",
			"Next Review Date",
			"Status",
		}
	} else {
		// Default to English
		headers = []string{
			"SL No.",
			"Reference Number",
			"Control Code",
			"Title",
			"Description",
			"Objective",
			"Control Category",
			"Control Type",
			"Owner",
			"Review Frequency",
			"Implementation Status",
			"Implementation Date",
			"Efficiency Status",
			"Last Test Date",
			"Test Frequency",
			"Next Test Date",
			"Last Review Date",
			"Next Review Date",
			"Status",
		}
	}

	// Convert the control data to [][]interface{} based on language setting
	rows := make([][]interface{}, len(data))
	for i, control := range data {
		if userSetting.Language == "ar" {
			rows[i] = []interface{}{
				i + 1,
				control.ReferenceNumber,
				utility.GetStringValue(control.Code),
				control.TitleAr,
				utility.GetStringValue(control.Description),
				utility.GetStringValue(control.Objective),
				control.ControlCategoryTitleAr,
				control.ControlTypeTitleAr,
				utility.GetStringValue(control.OwnerFirstNameAr) + " " + utility.GetStringValue(control.OwnerLastNameAr),
				control.ReviewFrequencyTitleAr,
				control.ControlImplementationStatusTitleAr,
				utility.FormatDate(control.ImplementedAt),
				control.ControlEfficiencyStatusTitleAr,
				utility.FormatDate(control.LastTestDate),
				utility.GetStringValue(control.TestFrequencyTitleAr),
				utility.FormatDate(control.NextTestDate),
				utility.FormatDate(control.LastReviewDate),
				utility.FormatDate(control.NextReviewDate),
				control.LifecycleStatusTitleAr,
			}
		} else {
			rows[i] = []interface{}{
				i + 1,
				control.ReferenceNumber,
				utility.GetStringValue(control.Code),
				control.TitleEn,
				utility.GetStringValue(control.Description),
				utility.GetStringValue(control.Objective),
				control.ControlCategoryTitleEn,
				control.ControlTypeTitleEn,
				utility.GetStringValue(control.OwnerFirstNameEn) + " " + utility.GetStringValue(control.OwnerLastNameEn),
				control.ReviewFrequencyTitleEn,
				control.ControlImplementationStatusTitleEn,
				utility.FormatDate(control.ImplementedAt),
				control.ControlEfficiencyStatusTitleEn,
				utility.FormatDate(control.LastTestDate),
				utility.GetStringValue(control.TestFrequencyTitleAr),
				utility.FormatDate(control.NextTestDate),
				utility.FormatDate(control.LastReviewDate),
				utility.FormatDate(control.NextReviewDate),
				control.LifecycleStatusTitleEn,
			}
		}
	}

	// Generate the Excel file
	return utility.GenerateExcelFile(headers, rows)
}

func (cs *ControlService) UpdateControlRiskRating(control *models.Control, tx utility.DatabaseTransaction) error {
	err := cs.repo.Update(control, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (cs *ControlService) MarkControlAsApplicable(controlID uint, currentUserID uint, controlApplicableFormDTO *organizationDTOs.ControlApplicableFormDTO) error {
	existingControl, err := cs.GetByID(controlID)
	if err != nil {
		return err
	}

	parsedApplicabilityChangeDate, err := utility.ParseDate(*controlApplicableFormDTO.ApplicabilityChangeDate)
	if err != nil {
		return utility.NewValidationError(fmt.Errorf("invalid applicability change date"), "invalid_date", "applicability change date")
	}

	defaultImplementationStatus, err := cs.controlImplementationStatusSvc.GetByType("not-implemented")
	if err != nil {
		return err
	}

	defaultEfficiencyStatus, err := cs.controlEfficiencyStatusSvc.GetByType("inefficient")
	if err != nil {
		return err
	}

	existingControl.ControlImplementationStatusID = &defaultImplementationStatus.ID
	existingControl.ControlEfficiencyStatusID = &defaultEfficiencyStatus.ID
	existingControl.Justification = controlApplicableFormDTO.Justification
	existingControl.ApplicabilityChangeDate = parsedApplicabilityChangeDate
	existingControl.IsApplicable = true
	existingControl.UpdatedByID = &currentUserID

	err = cs.repo.Update(existingControl)
	if err != nil {
		return err
	}

	return nil
}

func (cs *ControlService) MarkControlAsNotApplicable(controlID uint, currentUserID uint, controlApplicableFormDTO *organizationDTOs.ControlApplicableFormDTO) error {
	existingControl, err := cs.GetByID(controlID)
	if err != nil {
		return err
	}

	parsedApplicabilityChangeDate, err := utility.ParseDate(*controlApplicableFormDTO.ApplicabilityChangeDate)
	if err != nil {
		return utility.NewValidationError(fmt.Errorf("invalid applicability change date"), "invalid_date", "applicability change date")
	}

	existingControl.Justification = controlApplicableFormDTO.Justification
	existingControl.ApplicabilityChangeDate = parsedApplicabilityChangeDate
	existingControl.IsApplicable = false
	existingControl.UpdatedByID = &currentUserID

	err = cs.repo.Update(existingControl)
	if err != nil {
		return err
	}

	return nil
}

func (ps *ControlService) ImportProcesses(file multipart.File, currentUserID uint) ([]int, error) {
	mandatoryFields := []string{"Code", "ControlTypeID", "ControlCategoryID", "TitleEn", "TitleAr", "OwnerID", "TestResponsibleUserID"}
	// // Parse the Excel file
	fmt.Println(mandatoryFields)
	data, err := utility.ParseExcelFile(file, mandatoryFields)
	if err != nil {
		return nil, fmt.Errorf("invalid file format: %v", err)
	}

	// Track failed rows
	var failedRows []int

	// Process each row
	for i, row := range data {
		// Map the row data to the ProcessFormDTO
		processForm := organizationDTOs.ControlFormDTO{
			Code:                  utility.ParseOptionalString(row["Code"]),
			ControlTypeID:         utility.ParseUint(row["ControlTypeID"]),
			ControlCategoryID:     utility.ParseOptionalUint(row["ControlCategoryID"]),
			OwnerID:               utility.ParseUint(row["OwnerID"]),
			TitleEn:               row["TitleEn"],
			TitleAr:               row["TitleAr"],
			TestResponsibleUserID: utility.ParseUint(row["TestResponsibleUserID"]),
		}

		// Validate and create the process
		_, err := ps.Create(currentUserID, &processForm)
		if err != nil {
			// Log the row number that failed (+2 to account for the header row)
			failedRows = append(failedRows, i+2)
		}
	}

	return failedRows, nil
}

func (cs *ControlService) GenerateControlTemplate() (*os.File, error) {
	// Headers for the first sheet (Template)
	templateHeaders := []string{
		"Code", "ControlType", "ControlCategory", "Owner", "TitleEn", "TitleAr",
		"Description", "ImplementationStatus", "EfficiencyStatus", "TestResponsibleUser",
	}
	templateData := [][]interface{}{} // Empty template data

	// Fetch Control Types
	controlTypes, err := cs.controlTypeSvc.GetAll("", "", "", "active")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch control types: %v", err)
	}
	controlTypeHeaders := []string{"ControlType"}
	controlTypeData := [][]interface{}{}
	for _, ct := range controlTypes {
		controlTypeData = append(controlTypeData, []interface{}{ct.TitleEn})
	}

	// Fetch Control Categories
	controlCategories, err := cs.controlCategorySvc.GetAll("", "", "", "active")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch control categories: %v", err)
	}
	controlCategoryHeaders := []string{"ControlCategory"}
	controlCategoryData := [][]interface{}{}
	for _, cc := range controlCategories {
		controlCategoryData = append(controlCategoryData, []interface{}{cc.TitleEn})
	}

	// Fetch Owners
	owners, err := cs.userSvc.GetAllUsers()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch owners: %v", err)
	}
	ownerHeaders := []string{"OwnerName", "Email", "Designation", "PrimaryDepartment"}
	ownerData := [][]interface{}{}
	for _, owner := range owners {
		ownerData = append(ownerData, []interface{}{
			fmt.Sprintf("%s %s", owner.FirstNameEN, owner.LastNameEN),
			owner.Email,
			owner.Designation,
			owner.PrimaryUnit,
		})
	}

	// Fetch Control Implementation Status
	implementationStatuses, err := cs.controlImplementationStatusSvc.GetAll("", "", "", "active")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch control implementation statuses: %v", err)
	}
	implementationStatusHeaders := []string{"ImplementationStatus"}
	implementationStatusData := [][]interface{}{}
	for _, status := range implementationStatuses {
		implementationStatusData = append(implementationStatusData, []interface{}{status.TitleEn})
	}

	// Fetch Control Efficiency Status
	efficiencyStatuses, err := cs.controlEfficiencyStatusSvc.GetAll("", "", "", "active")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch control efficiency statuses: %v", err)
	}
	efficiencyStatusHeaders := []string{"EfficiencyStatus"}
	efficiencyStatusData := [][]interface{}{}
	for _, status := range efficiencyStatuses {
		efficiencyStatusData = append(efficiencyStatusData, []interface{}{status.TitleEn})
	}

	// Fetch Test Responsible Users
	testResponsibleUsers, err := cs.userSvc.GetAllUsers()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch test responsible users: %v", err)
	}
	testResponsibleUserHeaders := []string{"TestResponsibleUserName", "Email", "Designation", "PrimaryDepartment"}
	testResponsibleUserData := [][]interface{}{}
	for _, user := range testResponsibleUsers {
		testResponsibleUserData = append(testResponsibleUserData, []interface{}{
			fmt.Sprintf("%s %s", user.FirstNameEN, user.LastNameEN),
			user.Email,
			user.Designation,
			user.PrimaryUnit,
		})
	}

	// Generate the Excel file with multiple sheets
	file, err := utility.GenerateExcelFileWithSheets(map[string][][]interface{}{
		"Template":             append([][]interface{}{utility.ConvertToInterfaceSlice(templateHeaders)}, templateData...),
		"ControlType":          append([][]interface{}{utility.ConvertToInterfaceSlice(controlTypeHeaders)}, controlTypeData...),
		"ControlCategory":      append([][]interface{}{utility.ConvertToInterfaceSlice(controlCategoryHeaders)}, controlCategoryData...),
		"Owner":                append([][]interface{}{utility.ConvertToInterfaceSlice(ownerHeaders)}, ownerData...),
		"ImplementationStatus": append([][]interface{}{utility.ConvertToInterfaceSlice(implementationStatusHeaders)}, implementationStatusData...),
		"EfficiencyStatus":     append([][]interface{}{utility.ConvertToInterfaceSlice(efficiencyStatusHeaders)}, efficiencyStatusData...),
		"TestResponsibleUser":  append([][]interface{}{utility.ConvertToInterfaceSlice(testResponsibleUserHeaders)}, testResponsibleUserData...),
	})
	if err != nil {
		return nil, err
	}

	return file, nil
}

func (cs *ControlService) ValidateControlTemplate(file io.Reader) ([]map[string]interface{}, error) {
	// Read Excel file
	rows, err := utility.ReadExcelFile(file)
	if err != nil {
		return nil, err
	}

	// Define headers in the input format and map them to snake case
	inputHeaders := []string{"Code", "ControlType", "ControlCategory", "Owner", "TitleEn", "TitleAr", "Description", "ImplementationStatus", "EfficiencyStatus", "TestResponsibleUser"}
	snakeCaseHeaders := []string{"code", "control_type", "control_category", "owner", "title_en", "title_ar", "description", "implementation_status", "efficiency_status", "test_responsible_user"}

	if len(rows) == 0 {
		return nil, fmt.Errorf("empty Excel file")
	}

	// Validate headers
	if !utility.ValidateHeaders(inputHeaders, rows[0]) {
		return nil, utility.NewValidationError(fmt.Errorf("invalid headers in the Excel file"), "invalid_header", "Control Template")
	}

	// Convert rows to input format
	var input []map[string]interface{}
	for _, row := range rows[1:] { // Skip header row
		record := make(map[string]interface{})
		for i, header := range snakeCaseHeaders {
			if i < len(row) {
				record[header] = row[i]
			} else {
				record[header] = "" // Pad missing columns
			}
		}
		input = append(input, record)
	}

	validatedData, err := cs.validateControls(input)
	if err != nil {
		return nil, err
	}

	return validatedData, nil
}

func (cs *ControlService) ValidateImportControls(input []map[string]interface{}) ([]map[string]interface{}, error) {
	validatedData, err := cs.validateControls(input)
	if err != nil {
		return nil, err
	}

	return validatedData, nil
}

func (cs *ControlService) ImportControls(input []map[string]interface{}, currentUserID uint) ([]uint, error) {
	validatedData, err := cs.validateControls(input)
	if err != nil {
		return nil, err
	}

	// Check for any validation errors before processing the import
	for _, record := range validatedData {
		if errors, ok := record["errors"].(map[string]interface{}); ok && len(errors) > 0 {
			return nil, utility.NewValidationError(fmt.Errorf("invalid data in the import file"), "invalid_data", "Invalid Data")
		}
	}

	var createdControlIDs []uint
	for _, record := range validatedData {
		controlFormDTO := &organizationDTOs.ControlFormDTO{
			Code:                          utility.GetOptionalStringValue(record["code"]),
			ControlTypeID:                 record["control_type_id"].(uint),
			ControlCategoryID:             utility.GetOptionalUintValue(record["control_category_id"]),
			OwnerID:                       record["owner_id"].(uint),
			TitleEn:                       record["title_en"].(string),
			TitleAr:                       record["title_ar"].(string),
			Description:                   utility.GetOptionalStringValue(record["description"]),
			ControlImplementationStatusID: utility.GetOptionalUintValue(record["implementation_status_id"]),
			ControlEfficiencyStatusID:     utility.GetOptionalUintValue(record["efficiency_status_id"]),
			TestResponsibleUserID:         record["test_responsible_user_id"].(uint),
			IsApplicable:                  true,
		}

		controlID, err := cs.Create(currentUserID, controlFormDTO)
		if err != nil {
			return nil, fmt.Errorf("failed to create control for record: %+v, error: %v", record, err)
		}

		createdControlIDs = append(createdControlIDs, controlID)
	}

	return createdControlIDs, nil
}

func (cs *ControlService) validateControls(input []map[string]interface{}) ([]map[string]interface{}, error) {
	headers := []string{"code", "control_type", "control_category", "owner", "title_en", "title_ar", "description", "implementation_status", "efficiency_status", "test_responsible_user"}

	// Fetch control types
	controlTypes, err := cs.controlTypeSvc.GetAll("", "", "", "active")
	if err != nil {
		return nil, fmt.Errorf("error fetching control types: %v", err)
	}
	controlTypeTitles := make(map[string]uint)
	for _, controlType := range controlTypes {
		controlTypeTitles[strings.TrimSpace(controlType.TitleEn)] = controlType.ID
	}

	// Fetch control categories
	controlCategories, err := cs.controlCategorySvc.GetAll("", "", "", "active")
	if err != nil {
		return nil, fmt.Errorf("error fetching control categories: %v", err)
	}
	controlCategoryTitles := make(map[string]uint)
	for _, controlCategory := range controlCategories {
		controlCategoryTitles[strings.TrimSpace(controlCategory.TitleEn)] = controlCategory.ID
	}

	// Fetch control implementation statuses
	implementationStatuses, err := cs.controlImplementationStatusSvc.GetAll("", "", "", "active")
	if err != nil {
		return nil, fmt.Errorf("error fetching implementation statuses: %v", err)
	}
	implementationStatusTitles := make(map[string]uint)
	for _, status := range implementationStatuses {
		implementationStatusTitles[strings.TrimSpace(status.TitleEn)] = status.ID
	}

	// Fetch control efficiency statuses
	efficiencyStatuses, err := cs.controlEfficiencyStatusSvc.GetAll("", "", "", "active")
	if err != nil {
		return nil, fmt.Errorf("error fetching efficiency statuses: %v", err)
	}
	efficiencyStatusTitles := make(map[string]uint)
	for _, status := range efficiencyStatuses {
		efficiencyStatusTitles[strings.TrimSpace(status.TitleEn)] = status.ID
	}

	// Fetch users
	users, err := cs.userSvc.GetAllUsers()
	if err != nil {
		return nil, fmt.Errorf("error fetching users: %v", err)
	}
	userFullNames := make(map[string]uint)
	for _, user := range users {
		fullName := strings.TrimSpace(user.FirstNameEN + " " + user.LastNameEN)
		userFullNames[fullName] = uint(user.ID)
	}

	// Perform validation
	var validatedData []map[string]interface{}

	for _, record := range input {
		rowErrors := make(map[string]interface{})
		controlRecord := make(map[string]interface{})

		// Copy original data
		for key, value := range record {
			controlRecord[key] = value
		}

		// Validation
		for _, header := range headers {
			value, ok := record[header].(string)
			switch header {

			case "control_type":
				if !ok || value == "" {
					rowErrors[header] = "Control Type is required"
				} else if id, exists := controlTypeTitles[strings.TrimSpace(value)]; !exists {
					rowErrors[header] = fmt.Sprintf("Invalid control type: %s", value)
				} else {
					controlRecord["control_type_id"] = id
				}

			case "control_category":
				if value == "" {
					// Allow empty control category, set to nil in the record
					controlRecord["control_category_id"] = nil
				} else if id, exists := controlCategoryTitles[strings.TrimSpace(value)]; !exists {
					rowErrors[header] = fmt.Sprintf("Invalid control category: %s", value)
				} else {
					controlRecord["control_category_id"] = id
				}

			case "owner":
				if !ok || value == "" {
					rowErrors[header] = "Owner is required"
				} else if id, exists := userFullNames[strings.TrimSpace(value)]; !exists {
					rowErrors[header] = fmt.Sprintf("Invalid owner: %s", value)
				} else {
					controlRecord["owner_id"] = id
				}

			case "title_en", "title_ar":
				if !ok || value == "" {
					rowErrors[header] = "Field is required"
				}

			case "description":
				// Optional, no validation required
			case "implementation_status":
				if !ok || value == "" {
					rowErrors[header] = "Implementation Status is required"
				} else if id, exists := implementationStatusTitles[strings.TrimSpace(value)]; !exists {
					rowErrors[header] = fmt.Sprintf("Invalid implementation status: %s", value)
				} else {
					controlRecord["implementation_status_id"] = id
				}

			case "efficiency_status":
				if !ok || value == "" {
					rowErrors[header] = "Efficiency Status is required"
				} else if id, exists := efficiencyStatusTitles[strings.TrimSpace(value)]; !exists {
					rowErrors[header] = fmt.Sprintf("Invalid efficiency status: %s", value)
				} else {
					controlRecord["efficiency_status_id"] = id
				}

			case "test_responsible_user":
				if id, exists := userFullNames[strings.TrimSpace(value)]; !exists {
					rowErrors[header] = fmt.Sprintf("Invalid test responsible user: %s", value)
				} else {
					controlRecord["test_responsible_user_id"] = id
				}
			}
		}

		// Attach errors
		if len(rowErrors) == 0 {
			controlRecord["errors"] = nil
		} else {
			for _, header := range headers {
				if _, exists := rowErrors[header]; !exists {
					rowErrors[header] = "" // Fill empty error messages
				}
			}
			controlRecord["errors"] = rowErrors
		}

		validatedData = append(validatedData, controlRecord)
	}

	return validatedData, nil
}
