package services

import (
	"errors"
	"fmt"
	commonDTOs "jethur/grc/common/dtos"
	commonServices "jethur/grc/common/services"
	logServices "jethur/grc/logs/services"
	masterServices "jethur/grc/masters/services"
	"jethur/grc/models"
	organizationDTOs "jethur/grc/organization/dtos"
	organizationRepositories "jethur/grc/organization/repositories"
	"jethur/utility"
	"os"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type ControlServiceInterface interface {
	GetAll(params *organizationDTOs.ControlListQueryParamsDTO) (*organizationDTOs.ControlListResponseDTO, error)
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

	UpdateControlRiskRating(control *models.Control, tx utility.DatabaseTransaction) error
	MarkControlAsApplicable(controlID uint, currentUserID uint, controlApplicableFormDTO *organizationDTOs.ControlApplicableFormDTO) error
	MarkControlAsNotApplicable(controlID uint, currentUserID uint, controlApplicableFormDTO *organizationDTOs.ControlApplicableFormDTO) error
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
	userSettingSvc UserSettingServiceInterface) ControlServiceInterface {
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

	//todo :- also add extra fields which required in details screen

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
			UserID:               controlFormDTO.OwnerID,
			ModuleIdentifier:     "organization",
			NotificationType:     string(commonDTOs.AssignmentType),
			ReferenceNumber:      referenceNumber,
			HeadingEn:            "Control Assignment",
			HeadingAr:            "مهمة التحكم",
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

		emailTo := []string{owner.Email}
		subject := fmt.Sprintf("Control/%s: Added", control.TitleEn)
		templateRelativePath := "organization/templates/email/control_create_owner_mail.html"
		emailData := organizationDTOs.ControlMailDTO{
			FirstNameEn: owner.FirstNameEn,
			LastNameEn:  owner.LastNameEn,
			MessageEn:   control.TitleEn,
		}

		emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
			EmailTo:              emailTo,
			Subject:              subject,
			TemplateRelativePath: templateRelativePath,
			EmailData:            emailData,
		}
		err = cs.notificationSvc.SendNotification(
			dbNotificationDTO, emailNotificationDTO, "control")
		if err != nil {
			utility.LogError("Failed to send control notifications", map[string]interface{}{
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

	oldOwnerID := existingControl.OwnerID

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
	existingControl.ApplicabilityChangeDate = dates.ApplicabilityChangeDate
	existingControl.UpdatedByID = &currentUserID

	err = cs.repo.Update(existingControl)
	if err != nil {
		return err
	}

	// If the owner has changed, send notifications to the old and new owner
	go func() {
		if oldOwnerID != controlFormDTO.OwnerID {
			urlParam := fmt.Sprintf(`{"id":%d}`, id)
			entity := "control"
			removeAction := "responsible-user-remove"
			assignAction := "responsible-user-assign"

			// Notify the old owner
			oldOwnerNotification := &commonDTOs.DBNotificationDTO{
				UserID:               oldOwnerID,
				ModuleIdentifier:     "organization",
				NotificationType:     string(commonDTOs.InformationType),
				ReferenceNumber:      existingControl.ReferenceNumber,
				HeadingEn:            "Control Ownership Removed",
				HeadingAr:            "إزالة مسؤولية التحكم",
				TitleEn:              controlFormDTO.TitleEn,
				TitleAr:              controlFormDTO.TitleAr,
				Entity:               &entity,
				Action:               &removeAction,
				URLParam:             urlParam,
				IsSystemNotification: false,
				CreatedByID:          &currentUserID,
			}

			oldOwner, err := cs.userSvc.GetByID(oldOwnerID)
			if err != nil {
				utility.LogError("Failed to fetch old owner details", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}

			oldOwnerEmail := &commonDTOs.EmailNotificationDTO{
				EmailTo: []string{oldOwner.Email},
				Subject: fmt.Sprintf("Control/%s: Ownership Removed", existingControl.TitleEn),
				EmailData: organizationDTOs.ControlMailDTO{
					FirstNameEn: oldOwner.FirstNameEn,
					LastNameEn:  oldOwner.LastNameEn,
					MessageEn:   existingControl.TitleEn,
				},
				TemplateRelativePath: "organization/templates/email/control_remove_owner_mail.html",
			}

			// Notify the new owner
			newOwnerNotification := &commonDTOs.DBNotificationDTO{
				UserID:               controlFormDTO.OwnerID,
				ModuleIdentifier:     "organization",
				NotificationType:     string(commonDTOs.AssignmentType),
				ReferenceNumber:      existingControl.ReferenceNumber,
				HeadingEn:            "Control Ownership Assigned",
				HeadingAr:            "تم تعيين مسؤولية التحكم",
				TitleEn:              controlFormDTO.TitleEn,
				TitleAr:              controlFormDTO.TitleAr,
				Entity:               &entity,
				Action:               &assignAction,
				URLParam:             urlParam,
				IsSystemNotification: false,
				CreatedByID:          &currentUserID,
			}

			newOwner, err := cs.userSvc.GetByID(controlFormDTO.OwnerID)
			if err != nil {
				utility.LogError("Failed to fetch new owner details", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}

			newOwnerEmail := &commonDTOs.EmailNotificationDTO{
				EmailTo: []string{newOwner.Email},
				Subject: fmt.Sprintf("Control/%s: Ownership Assigned", existingControl.TitleEn),
				EmailData: organizationDTOs.ControlMailDTO{
					FirstNameEn: newOwner.FirstNameEn,
					LastNameEn:  newOwner.LastNameEn,
					MessageEn:   existingControl.TitleEn,
				},
				TemplateRelativePath: "organization/templates/email/control_change_owner_mail.html",
			}

			// Send notifications for old and new owners
			err = cs.notificationSvc.SendNotification(newOwnerNotification, newOwnerEmail, "control")
			if err != nil {
				utility.LogError("Failed to send new owner notifications", map[string]interface{}{
					"error": err.Error(),
				})
			}

			err = cs.notificationSvc.SendNotification(oldOwnerNotification, oldOwnerEmail, "control")
			if err != nil {
				utility.LogError("Failed to send old owner notifications", map[string]interface{}{
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
			ReferenceNumber string
		}{
			OwnerID:         control.OwnerID,
			TitleEn:         control.TitleEn,
			TitleAr:         control.TitleAr,
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
					UserID:               *notification.OwnerID,
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

				ownerEmailTo := []string{ownerUser.Email}
				subject := fmt.Sprintf("Control/%s: Deleted", notification.TitleEn)
				templateRelativePath := "organization/templates/email/control_delete_mail.html"
				redirectURL := fmt.Sprintf("control/deleted/%d", notification.OwnerID)

				emailData := organizationDTOs.ControlMailDTO{
					FirstNameEn: ownerUser.FirstNameEn,
					LastNameEn:  ownerUser.LastNameEn,
					MessageEn:   notification.TitleEn,
					RedirectURL: redirectURL,
				}

				emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
					EmailTo:              ownerEmailTo,
					Subject:              subject,
					TemplateRelativePath: templateRelativePath,
					EmailData:            emailData,
				}

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

	if _, err := cs.controlCategorySvc.GetByID(uint(controlFormDTO.ControlCategoryID)); err != nil {
		return utility.NewValidationError(fmt.Errorf("control category does not exist"), "not_exist", "control category")
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
