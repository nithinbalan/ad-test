package services

import (
	"errors"
	"fmt"
	"os"
	"time"

	commonDTOs "jethur/grc/common/dtos"
	commonServices "jethur/grc/common/services"
	logServices "jethur/grc/logs/services"
	masterServices "jethur/grc/masters/services"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	riskDTOs "jethur/grc/risk/dtos"
	riskRepositories "jethur/grc/risk/repositories"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type KRIRegisterServiceInterface interface {
	GetAll(params *riskDTOs.KRIRegisterListQueryParamsDTO) (*riskDTOs.KRIRegisterListResponseDTO, error)
	GetFullDetailsByID(id uint) (*riskDTOs.KRIRegisterDetailsDTO, error)
	GetByID(id uint) (*models.KRIRegister, error)
	Create(currentUserID uint, kriRegisterFormDTO *riskDTOs.KRIRegisterFormDTO) (uint, error)
	Update(id uint, currentUserID uint, kriRegisterFormDTO *riskDTOs.KRIRegisterFormDTO) error
	UpdateValue(kriRegister *models.KRIRegister, tx utility.DatabaseTransaction) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)

	ExportKriRegister(data []*riskDTOs.KRIRegisterListDTO, currentUserID uint) (*os.File, error)
}

// Service and implementations
type KRIRegisterService struct {
	repo               riskRepositories.KRIRegisterRepositoryInterface
	activityLogSvc     logServices.ActivityLogServiceInterface
	kriCategorySvc     KRICategoryServiceInterface
	userSvc            organizationServices.UserServiceInterface
	riskLevelSvc       RiskLevelServiceInterface
	reviewFrequencySvc masterServices.ReviewFrequencyServiceInterface
	unitSvc            masterServices.UnitServiceInterface
	userSettingSvc     organizationServices.UserSettingServiceInterface
	notificationSvc    commonServices.NotificationServiceInterface
}

func NewKRIRegisterService(
	repo riskRepositories.KRIRegisterRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	kriCategorySvc KRICategoryServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	riskLevelSvc RiskLevelServiceInterface,
	reviewFrequencySvc masterServices.ReviewFrequencyServiceInterface,
	unitSvc masterServices.UnitServiceInterface,
	userSettingSvc organizationServices.UserSettingServiceInterface,
	notificationSvc commonServices.NotificationServiceInterface,
) KRIRegisterServiceInterface {
	return &KRIRegisterService{
		repo:               repo,
		activityLogSvc:     activityLogSvc,
		kriCategorySvc:     kriCategorySvc,
		userSvc:            userSvc,
		riskLevelSvc:       riskLevelSvc,
		reviewFrequencySvc: reviewFrequencySvc,
		unitSvc:            unitSvc,
		userSettingSvc:     userSettingSvc,
		notificationSvc:    notificationSvc,
	}
}

func (krs *KRIRegisterService) GetAll(params *riskDTOs.KRIRegisterListQueryParamsDTO) (*riskDTOs.KRIRegisterListResponseDTO, error) {
	return krs.repo.GetAll(params)
}

func (krs KRIRegisterService) GetByID(id uint) (*models.KRIRegister, error) {
	kriRegister, err := krs.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return kriRegister, nil
}

func (krs KRIRegisterService) GetFullDetailsByID(id uint) (*riskDTOs.KRIRegisterDetailsDTO, error) {
	var kriRegisterDTO riskDTOs.KRIRegisterDetailsDTO

	kriRegister, err := krs.repo.GetFullDetailsByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	err = mapstructure.Decode(kriRegister, &kriRegisterDTO)
	if err != nil {
		return nil, err
	}

	kriRegisterNextReviewDate := kriRegister.NextReviewDate.Format("2006-01-02")
	kriRegisterDTO.NextReviewDateString = kriRegisterNextReviewDate

	if kriRegister.LastReviewedDate != nil {
		kriRegisterLastReviewedDate := kriRegister.LastReviewedDate.Format("2006-01-02")
		kriRegisterDTO.LastReviewedDateString = kriRegisterLastReviewedDate
	}

	kriRegisterDTO.CreatedAt = kriRegister.CreatedAt
	kriRegisterDTO.UpdatedAt = kriRegister.UpdatedAt

	//todo :- also add extra fields which required in details screen

	return &kriRegisterDTO, nil
}

func (krs *KRIRegisterService) Create(currentUserID uint, kriRegisterFormDTO *riskDTOs.KRIRegisterFormDTO) (uint, error) {
	if err := krs.validateKRIRegisterFormDTO(kriRegisterFormDTO, nil); err != nil {
		return 0, err
	}
	dates, err := krs.parseKRIRegisterDates(kriRegisterFormDTO)
	if err != nil {
		return 0, err
	}
	referenceNumber, err := krs.repo.GenerateReferenceNumber()
	if err != nil {
		return 0, err
	}

	var organizationUnitID uint
	unitID, err := krs.userSvc.GetUserPrimaryUnitID(kriRegisterFormDTO.ResponsibleUserID)
	if err != nil {
		return 0, utility.NewValidationError(fmt.Errorf("responsible unit not set forthe responsible user"), "not_exist", "responsible unit")
	}
	organizationUnitID = unitID
	kriRegister := &models.KRIRegister{
		ReferenceNumber:   referenceNumber,
		KRICategoryID:     kriRegisterFormDTO.KRICategoryID,
		ResponsibleUnitID: organizationUnitID,
		ResponsibleUserID: kriRegisterFormDTO.ResponsibleUserID,
		ReviewFrequencyID: kriRegisterFormDTO.ReviewFrequencyID,
		RiskLevelID:       kriRegisterFormDTO.RiskLevelID,
		UnitID:            kriRegisterFormDTO.UnitID,
		Title:             kriRegisterFormDTO.Title,
		Description:       kriRegisterFormDTO.Description,
		CalculationType:   kriRegisterFormDTO.CalculationType,
		LowThreshold:      kriRegisterFormDTO.LowThreshold,
		MediumThreshold:   kriRegisterFormDTO.MediumThreshold,
		HighThreshold:     kriRegisterFormDTO.HighThreshold,
		CurrentValue:      kriRegisterFormDTO.CurrentValue,
		NextReviewDate:    *dates.NextReviewDate,
		LastReviewNote:    kriRegisterFormDTO.LastReviewNote,
		CreatedBYID:       currentUserID,
		UpdatedAt:         nil,
	}
	if kriRegisterFormDTO.LastReviewedDate != "" {
		kriRegister.LastReviewedDate = dates.LastReviewedDate
	}

	id, err := krs.repo.Create(kriRegister)
	if err != nil {
		return 0, err
	}
	go func() {
		err := krs.activityLogSvc.LogActivity(
			models.CreateAction,
			models.KRIRegister{}.TableName(),
			id,
			kriRegisterFormDTO,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.KRIRegister{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}

		// Prepare notification and email
		urlParam := fmt.Sprintf(`{"id":%d}`, id)
		entity := "kri-register"
		action := "responsible-user-assign"

		dbNotificationDTO := &commonDTOs.DBNotificationDTO{
			UserID:               kriRegisterFormDTO.ResponsibleUserID,
			ModuleIdentifier:     "risk",
			NotificationType:     string(commonDTOs.AssignmentType),
			HeadingEn:            "KRI Register Assignment",
			HeadingAr:            "مهمة تسجيل KRI",
			TitleEn:              kriRegisterFormDTO.Title,
			Entity:               &entity,
			Action:               &action,
			URLParam:             urlParam,
			IsSystemNotification: false,
			CreatedByID:          &currentUserID,
		}

		responsibleUser, err := krs.userSvc.GetByID(kriRegisterFormDTO.ResponsibleUserID)
		if err != nil {
			utility.LogError("Failed to fetch responsible user", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		createdBy, err := krs.userSvc.GetByID(currentUserID)
		if err != nil {
			utility.LogError("Failed to fetch responsible user", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		emailTo := []string{responsibleUser.Email}
		subject := fmt.Sprintf("KRI Register/%s: Responsibility Assigned", kriRegister.Title)
		appURL := utility.GetAppURL()

		emailData := map[string]interface{}{
			"FirstNameEn":          responsibleUser.FirstNameEn,
			"LastNameEn":           responsibleUser.LastNameEn,
			"KriReferenceNumber":   kriRegister.ReferenceNumber,
			"KriTitleEn":           kriRegister.Title,
			"KriDescription":       kriRegister.Description,
			"CreatedByFirstNameEn": createdBy.FirstNameEn,
			"CreatedByLastNameEn":  createdBy.LastNameEn,
			"AssignedOn":           time.Now().Format("2006-01-02"),
			"RedirectURL":          fmt.Sprintf("%srisk/kri-register/%d/overview", appURL, id),
		}

		emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
			EmailTo:              emailTo,
			Subject:              subject,
			TemplateRelativePath: "risk/templates/email/kri_register_create_responsible_user_mail.html",
			EmailData:            emailData,
		}

		err = krs.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "risk")
		if err != nil {
			utility.LogError("Failed to send notifications", map[string]interface{}{
				"error": err.Error(),
			})
		}
	}()

	return id, nil
}

func (krs *KRIRegisterService) Update(id uint, currentUserID uint, kriRegisterFormDTO *riskDTOs.KRIRegisterFormDTO) error {
	existingKRIRegister, err := krs.GetByID(id)

	if err != nil {
		return err
	}

	if err := krs.validateKRIRegisterFormDTO(kriRegisterFormDTO, &id); err != nil {
		return err
	}

	dates, err := krs.parseKRIRegisterDates(kriRegisterFormDTO)
	if err != nil {
		return err
	}
	if existingKRIRegister.ResponsibleUserID != kriRegisterFormDTO.ResponsibleUserID {
		organizationUnitID, err := krs.userSvc.GetUserPrimaryUnitID(uint(kriRegisterFormDTO.ResponsibleUserID))
		if err != nil {
			return utility.NewValidationError(fmt.Errorf("responsible unit not set for the responsible user"), "not_exist", "responsible unit")
		}
		existingKRIRegister.ResponsibleUnitID = organizationUnitID
	}

	oldResponsibleUserID := existingKRIRegister.ResponsibleUserID

	existingKRIRegister.KRICategoryID = kriRegisterFormDTO.KRICategoryID
	existingKRIRegister.ResponsibleUserID = kriRegisterFormDTO.ResponsibleUserID
	existingKRIRegister.ReviewFrequencyID = kriRegisterFormDTO.ReviewFrequencyID
	existingKRIRegister.RiskLevelID = kriRegisterFormDTO.RiskLevelID
	existingKRIRegister.UnitID = kriRegisterFormDTO.UnitID
	existingKRIRegister.Title = kriRegisterFormDTO.Title
	existingKRIRegister.Description = kriRegisterFormDTO.Description
	existingKRIRegister.CalculationType = kriRegisterFormDTO.CalculationType
	existingKRIRegister.LowThreshold = kriRegisterFormDTO.LowThreshold
	existingKRIRegister.MediumThreshold = kriRegisterFormDTO.MediumThreshold
	existingKRIRegister.HighThreshold = kriRegisterFormDTO.HighThreshold
	existingKRIRegister.CurrentValue = kriRegisterFormDTO.CurrentValue

	if kriRegisterFormDTO.LastReviewedDate != "" {
		existingKRIRegister.LastReviewedDate = dates.LastReviewedDate
	}
	existingKRIRegister.NextReviewDate = *dates.NextReviewDate
	existingKRIRegister.LastReviewNote = kriRegisterFormDTO.LastReviewNote
	existingKRIRegister.UpdatedBYID = &currentUserID

	err = krs.repo.Update(existingKRIRegister)
	if err != nil {
		return err
	}

	go func() {
		// Log the update activity
		err := krs.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.KRIRegister{}.TableName(),
			id,
			kriRegisterFormDTO,
			currentUserID,
			nil,
		)
		if err != nil {
			utility.LogError("Failed to log activity", map[string]interface{}{
				"action":    models.UpdateAction,
				"tableName": models.KRIRegister{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			})
		}

		// Send notifications to old and new responsible users if the responsible user has changed
		if oldResponsibleUserID != existingKRIRegister.ResponsibleUserID {
			// URL param for notifications
			urlParam := fmt.Sprintf(`{"id":%d}`, id)
			entity := "kri-register"

			// Old responsible user notification (Removal)
			oldResponsibleUserAction := "responsible-user-remove"
			oldResponsibleUserNotificationDTO := &commonDTOs.DBNotificationDTO{
				UserID:               oldResponsibleUserID,
				ModuleIdentifier:     "risk",
				NotificationType:     string(commonDTOs.InformationType),
				HeadingEn:            "Kri Register Responsibility Removed",
				HeadingAr:            "تمت إزالة مسؤولية تسجيل Kri",
				TitleEn:              existingKRIRegister.Title,
				Entity:               &entity,
				Action:               &oldResponsibleUserAction,
				URLParam:             urlParam,
				IsSystemNotification: false,
				CreatedByID:          &currentUserID,
			}

			// New responsible user notification (Assignment)
			newResponsibleUserAction := "responsible-user-assign"
			newResponsibleUserNotificationDTO := &commonDTOs.DBNotificationDTO{
				UserID:               kriRegisterFormDTO.ResponsibleUserID,
				ModuleIdentifier:     "risk",
				NotificationType:     string(commonDTOs.AssignmentType),
				HeadingEn:            "Kri Register Responsibility Assigned",
				HeadingAr:            "Kri تسجيل المسؤولية المعينة",
				TitleEn:              existingKRIRegister.Title,
				Entity:               &entity,
				Action:               &newResponsibleUserAction,
				URLParam:             urlParam,
				IsSystemNotification: false,
				CreatedByID:          &currentUserID,
			}

			// Fetch the responsible users' data
			oldResponsibleUser, err := krs.userSvc.GetByID(oldResponsibleUserID)
			if err != nil {
				utility.LogError("Failed to fetch old responsible user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}

			newResponsibleUser, err := krs.userSvc.GetByID(kriRegisterFormDTO.ResponsibleUserID)
			if err != nil {
				utility.LogError("Failed to fetch new responsible user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}

			// Fetch details of the user making the assignment
			updatedByUser, err := krs.userSvc.GetByID(currentUserID)
			if err != nil {
				utility.LogError("Failed to fetch assigning user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}

			appURL := utility.GetAppURL()

			// Email notification for the old responsible user
			oldResponsibleUserEmailTo := []string{oldResponsibleUser.Email}
			oldResponsibleUserSubject := fmt.Sprintf("Kri Register/%s: Responsibility Removed", existingKRIRegister.Title)
			oldResponsibleUserEmailData := map[string]interface{}{
				"FirstNameEn":          oldResponsibleUser.FirstNameEn,
				"LastNameEn":           oldResponsibleUser.LastNameEn,
				"KriReferenceNumber":   existingKRIRegister.ReferenceNumber,
				"KriTitle":             existingKRIRegister.Title,
				"KriDescription":       existingKRIRegister.Description,
				"UpdatedByFirstNameEn": updatedByUser.FirstNameEn,
				"UpdatedByLastNameEn":  updatedByUser.LastNameEn,
				"RedirectURL":          fmt.Sprintf("%srisk/kri-register/%d/overview", appURL, id),
			}

			oldResponsibleUserEmailNotificationDTO := &commonDTOs.EmailNotificationDTO{
				EmailTo:              oldResponsibleUserEmailTo,
				Subject:              oldResponsibleUserSubject,
				TemplateRelativePath: "risk/templates/email/kri_register_remove_responsible_user_mail.html",
				EmailData:            oldResponsibleUserEmailData,
			}

			// Email notification for the new responsible user
			newResponsibleUserEmailTo := []string{newResponsibleUser.Email}
			newResponsibleUserSubject := fmt.Sprintf("KRI Register/%s: Responsibility Assigned", existingKRIRegister.Title)
			newResponsibleUserEmailData := map[string]interface{}{
				"FirstNameEn":           newResponsibleUser.FirstNameEn,
				"LastNameEn":            newResponsibleUser.LastNameEn,
				"KriReferenceNumber":    existingKRIRegister.ReferenceNumber,
				"KriTitle":              existingKRIRegister.Title,
				"KriDescription":        existingKRIRegister.Description,
				"AssignedByFirstNameEn": updatedByUser.FirstNameEn,
				"AssignedByLastNameEn":  updatedByUser.LastNameEn,
				"AssignedOn":            time.Now().Format("2006-01-02"),
				"RedirectURL":           fmt.Sprintf("%srisk/kri-register/%d/overview", appURL, id),
			}

			newResponsibleUserEmailNotificationDTO := &commonDTOs.EmailNotificationDTO{
				EmailTo:              newResponsibleUserEmailTo,
				Subject:              newResponsibleUserSubject,
				TemplateRelativePath: "risk/templates/email/kri_register_change_responsible_user_mail.html",
				EmailData:            newResponsibleUserEmailData,
			}

			// Send notifications and emails
			err = krs.notificationSvc.SendNotification(newResponsibleUserNotificationDTO, newResponsibleUserEmailNotificationDTO, "risk")
			if err != nil {
				utility.LogError("Failed to send notifications to new responsible user", map[string]interface{}{
					"error": err.Error(),
				})
			}

			err = krs.notificationSvc.SendNotification(oldResponsibleUserNotificationDTO, oldResponsibleUserEmailNotificationDTO, "risk")
			if err != nil {
				utility.LogError("Failed to send notifications to old responsible user", map[string]interface{}{
					"error": err.Error(),
				})
			}
		}
	}()

	return nil
}

func (krs *KRIRegisterService) UpdateValue(kriRegister *models.KRIRegister, tx utility.DatabaseTransaction) error {
	err := krs.repo.Update(kriRegister, tx)
	if err != nil {
		return err
	}

	return nil
}

func (krs *KRIRegisterService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
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

		_, err := krs.repo.GetByID(id)

		if err != nil {
			continue
		}

		kriRegister, err := krs.GetFullDetailsByID(id)
		if err != nil {
			return nil, err
		}

		err = krs.repo.Delete(id)
		if err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}

		// Prepare notification data before deletion
		notifications = append(notifications, struct {
			OwnerID         *uint
			TitleEn         string
			TitleAr         string
			Description     string
			ReferenceNumber string
		}{
			OwnerID:         &kriRegister.ResponsibleUserID,
			TitleEn:         kriRegister.Title,
			Description:     *kriRegister.Description,
			ReferenceNumber: kriRegister.ReferenceNumber,
		})
	}

	if len(successfulIds) > 0 {

		go func() {
			err := krs.activityLogSvc.LogBulkActivities(
				models.DeleteAction, // Make sure this action constant is defined appropriately
				models.KRIRegister{}.TableName(),
				successfulIds,
				nil,
				currentUserID,
				nil,
			)
			if err != nil {
				// @todo - save to log file
				fmt.Printf("Failed to log bulk delete activity: %v\n", err)
			}

			// Send DB and email notifications for each successfully deleted kri register
			for _, notification := range notifications {
				urlParam := ""
				entity := "kri-register"
				action := "delete"

				// Prepare DB notification
				dbNotificationDTO := &commonDTOs.DBNotificationDTO{
					UserID:               *notification.OwnerID,
					ModuleIdentifier:     "risk",
					NotificationType:     string(commonDTOs.DeletionType),
					HeadingEn:            "KRI Register Deleted",
					HeadingAr:            "تم حذف تسجيل KRI",
					TitleEn:              notification.TitleEn,
					TitleAr:              notification.TitleAr,
					Entity:               &entity,
					Action:               &action,
					URLParam:             urlParam,
					IsSystemNotification: false,
					CreatedByID:          &currentUserID,
				}

				// Fetch owner details (responsible user)
				ownerUser, err := krs.userSvc.GetByID(*notification.OwnerID)
				if err != nil {
					utility.LogError("Failed to fetch owner details for notification", map[string]interface{}{
						"error": err.Error(),
					})
					continue
				}

				// Prepare email notification
				ownerEmailTo := []string{ownerUser.Email}
				subject := fmt.Sprintf("KRI Register/%s: Deleted", notification.TitleEn)
				templateRelativePath := "risk/templates/email/kri_register_delete_mail.html"

				emailData := map[string]interface{}{
					"FirstNameEn":        ownerUser.FirstNameEn,
					"LastNameEn":         ownerUser.LastNameEn,
					"KriReferenceNumber": notification.ReferenceNumber,
					"KriTitle":           notification.TitleEn,
					"KriDescription":     notification.Description,
				}

				emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
					EmailTo:              ownerEmailTo,
					Subject:              subject,
					TemplateRelativePath: templateRelativePath,
					EmailData:            emailData,
				}

				// Send DB and email notifications
				err = krs.notificationSvc.SendNotification(
					dbNotificationDTO, emailNotificationDTO, "risk")
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

func (krs *KRIRegisterService) validateKRIRegisterFormDTO(kriRegisterFormDTO *riskDTOs.KRIRegisterFormDTO, id *uint) error {
	if _, err := krs.kriCategorySvc.GetByID(kriRegisterFormDTO.KRICategoryID); err != nil {
		return utility.NewValidationError(fmt.Errorf("kri category does not exist"), "not_exist", "kri category")
	}

	if _, err := krs.userSvc.GetByID(kriRegisterFormDTO.ResponsibleUserID); err != nil {
		return utility.NewValidationError(fmt.Errorf("responsible user does not exist"), "not_exist", "responsible user")
	}
	if kriRegisterFormDTO.RiskLevelID != nil {
		if _, err := krs.riskLevelSvc.GetByID(*kriRegisterFormDTO.RiskLevelID); err != nil {
			return utility.NewValidationError(fmt.Errorf("risk level does not exist"), "not_exist", "risk level")
		}
	}
	if _, err := krs.unitSvc.GetByID(kriRegisterFormDTO.UnitID); err != nil {
		return utility.NewValidationError(fmt.Errorf("unit does not exist"), "not_exist", "unit")
	}

	if _, err := krs.reviewFrequencySvc.GetByID(kriRegisterFormDTO.ReviewFrequencyID); err != nil {
		return utility.NewValidationError(fmt.Errorf("review frequency does not exist"), "not_exist", "review frequency")
	}

	if id == nil {
		if exists, err := krs.repo.ExistsByTitle(kriRegisterFormDTO.Title, kriRegisterFormDTO.KRICategoryID); err != nil {
			return err
		} else if exists {
			return utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
		}

	} else {
		if exists, err := krs.repo.ExistsByTitleExcludingID(*id, kriRegisterFormDTO.Title, kriRegisterFormDTO.KRICategoryID); err != nil {
			return err
		} else if exists {
			return utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
		}
	}

	return nil
}

func (krs *KRIRegisterService) parseKRIRegisterDates(dto *riskDTOs.KRIRegisterFormDTO) (*riskDTOs.KRIRegisterParsedDates, error) {
	parsed := &riskDTOs.KRIRegisterParsedDates{}

	var err error
	if dto.LastReviewedDate != "" {
		parsed.LastReviewedDate, err = utility.ParseDate(dto.LastReviewedDate)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid last review date"), "invalid_date", "last review date")
		}
	}

	if dto.NextReviewDate != "" {
		parsed.NextReviewDate, err = utility.ParseDate(dto.NextReviewDate)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid next review date"), "invalid_date", "next review date")
		}
	}

	return parsed, nil
}

func (krs *KRIRegisterService) ExportKriRegister(data []*riskDTOs.KRIRegisterListDTO, currentUserID uint) (*os.File, error) {
	// Fetch user setting for language preference
	userSetting, err := krs.userSettingSvc.GetUserSettingsByUserID(currentUserID)
	if err != nil {
		return nil, err
	}

	// Define headers based on user language setting
	var headers []string
	if userSetting.Language == "ar" {
		headers = []string{
			"SL No.",
			"Ref No",
			"Title",
			"Description",
			"Category",
			"Calculation Type",
			"Responsible User",
			"Low Threshold",
			"Medium Threshold",
			"High Threshold",
			"Current Value",
			"Unit",
			"Risk Level",
			"Review Frequency",
			"Reviewed At",
			"Next Review Date",
		}
	} else {
		headers = []string{
			"SL No.",
			"Ref No",
			"Title",
			"Description",
			"Category",
			"Calculation Type",
			"Responsible User",
			"Low Threshold",
			"Medium Threshold",
			"High Threshold",
			"Current Value",
			"Unit",
			"Risk Level",
			"Review Frequency",
			"Reviewed At",
			"Next Review Date",
		}
	}
	// Prepare the data rows
	rows := make([][]interface{}, len(data))
	for i, kri := range data {
		var category, responsibleUser, unit, riskLevel, reviewFrequency string

		// Assign values based on language preference and check for nil pointers
		if userSetting.Language == "ar" {
			category = kri.KRICategoryTitleAr
			responsibleUser = kri.ResponsibleUserFirstNameAr + " " + kri.ResponsibleUserLastNameAr
			unit = kri.UnitTitleAr
			riskLevel = kri.RiskLevelTitleAr
			reviewFrequency = kri.ReviewFrequencyTitleAr
		} else {
			category = kri.KRICategoryTitleEn
			responsibleUser = kri.ResponsibleUserFirstNameEn + " " + kri.ResponsibleUserLastNameEn
			unit = kri.UnitTitleEn
			riskLevel = kri.RiskLevelTitleEn
			reviewFrequency = kri.ReviewFrequencyTitleEn
		}

		rows[i] = []interface{}{
			i + 1,
			kri.ReferenceNumber,
			kri.Title,
			utility.GetStringValue(kri.Description),
			category,
			kri.CalculationType,
			responsibleUser,
			utility.GetIntValueAsString(kri.LowThreshold),
			utility.GetIntValueAsString(kri.MediumThreshold),
			utility.GetIntValueAsString(kri.HighThreshold),
			utility.GetIntValueAsString(kri.CurrentValue),
			unit,
			riskLevel,
			reviewFrequency,
			kri.LastReviewedDate,
			kri.NextReviewDate,
		}
	}

	// Generate the Excel file
	return utility.GenerateExcelFile(headers, rows)
}
