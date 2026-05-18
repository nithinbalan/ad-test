package services

import (
	"errors"
	"fmt"
	commonDTOs "jethur/grc/common/dtos"
	commonServices "jethur/grc/common/services"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	riskDTOs "jethur/grc/risk/dtos"
	riskRepositories "jethur/grc/risk/repositories"
	settingsServices "jethur/grc/settings/services"
	"jethur/utility"
	"time"

	"gorm.io/gorm"
)

type RiskUpdatePlanRiskIdentificationServiceInterface interface {
	GetAllRiskIdentifications(planID uint, params *riskDTOs.RiskRegisterUpdatePlanListQueryParamsDTO) (*riskDTOs.RiskIdentificationListResponseDTO, error)
	GetRiskIdentificationDetails(planID uint, id uint) (*models.RiskRegister, error)
	CreateRiskIdentification(planID uint, currentUserID uint, riskRegisterFormDTO *riskDTOs.UpdatePlanRiskRegisterFormDTO) (uint, error)
	UpdateRiskIdentification(planID uint, id uint, currentUserID uint, risk *riskDTOs.UpdatePlanRiskRegisterFormDTO) error
	DeleteRiskIdentification(planID uint, ids []uint, currentUserID uint) ([]uint, error)
	MarkRiskIdentificationAsNotApplicable(planID uint, currentUserID uint, formDTO *riskDTOs.MarkTimelineNotApplicableDTO) error
}

type RiskUpdatePlanRiskIdentificationService struct {
	repo                 riskRepositories.RiskUpdatePlanRiskIdentificationRepositoryInterface
	riskUpdatePlanSvc    RiskRegisterUpdatePlanServiceInterface
	timelineTypeSvc      RiskUpdatePlanTimelineTypeServiceInterface
	riskRegisterSvc      RiskRegisterServiceInterface
	activityLogSvc       logServices.ActivityLogServiceInterface
	userSvc              organizationServices.UserServiceInterface
	organizationUnitSvc  settingsServices.OrganizationUnitServiceInterface
	riskStatusSvc        RiskStatusServiceInterface
	riskMatrixSvc        RiskMatrixServiceInterface
	organizationUnitsSvc organizationServices.OrganizationUnitsServiceInterface
	notificationSvc      commonServices.NotificationServiceInterface
	riskSettingsSvc      RiskSettingsServiceInterface
	riskCategorySvc      RiskCategoryServiceInterface
}

func NewRiskUpdatePlanRiskIdentificationService(
	repo riskRepositories.RiskUpdatePlanRiskIdentificationRepositoryInterface,
	riskUpdatePlanSvc RiskRegisterUpdatePlanServiceInterface,
	timelineTypeSvc RiskUpdatePlanTimelineTypeServiceInterface,
	riskRegisterSvc RiskRegisterServiceInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	organizationUnitSvc settingsServices.OrganizationUnitServiceInterface,
	riskStatusSvc RiskStatusServiceInterface,
	riskMatrixSvc RiskMatrixServiceInterface,
	organizationUnitsSvc organizationServices.OrganizationUnitsServiceInterface,
	notificationSvc commonServices.NotificationServiceInterface,
	riskSettingsSvc RiskSettingsServiceInterface,
	riskCategorySvc RiskCategoryServiceInterface,
) RiskUpdatePlanRiskIdentificationServiceInterface {
	return &RiskUpdatePlanRiskIdentificationService{
		repo:                 repo,
		riskUpdatePlanSvc:    riskUpdatePlanSvc,
		timelineTypeSvc:      timelineTypeSvc,
		riskRegisterSvc:      riskRegisterSvc,
		activityLogSvc:       activityLogSvc,
		userSvc:              userSvc,
		organizationUnitSvc:  organizationUnitSvc,
		riskStatusSvc:        riskStatusSvc,
		riskMatrixSvc:        riskMatrixSvc,
		organizationUnitsSvc: organizationUnitsSvc,
		notificationSvc:      notificationSvc,
		riskSettingsSvc:      riskSettingsSvc,
		riskCategorySvc:      riskCategorySvc,
	}
}

func (rrupids *RiskUpdatePlanRiskIdentificationService) GetAllRiskIdentifications(planID uint, params *riskDTOs.RiskRegisterUpdatePlanListQueryParamsDTO) (*riskDTOs.RiskIdentificationListResponseDTO, error) {
	return rrupids.repo.GetAllRiskIdentifications(planID, params)
}

func (rrupids *RiskUpdatePlanRiskIdentificationService) GetRiskIdentificationDetails(planID uint, id uint) (*models.RiskRegister, error) {
	// Verify plan exists
	_, err := rrupids.riskUpdatePlanSvc.GetByID(planID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, planID)
		}
		return nil, err
	}

	// Get risk details
	risk, err := rrupids.repo.GetRiskIdentificationDetailsByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	// Verify risk belongs to the plan
	if risk.RiskRegisterUpdatePlanID != nil && *risk.RiskRegisterUpdatePlanID != planID {
		return nil, utility.NewResourceNotFoundError(errors.New("risk does not belong to this plan"), id)
	}

	return risk, nil
}

func (rrupids *RiskUpdatePlanRiskIdentificationService) CreateRiskIdentification(planID uint, currentUserID uint, riskRegisterFormDTO *riskDTOs.UpdatePlanRiskRegisterFormDTO) (uint, error) {
	// Verify plan exists
	_, err := rrupids.riskUpdatePlanSvc.GetByID(planID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, utility.NewResourceNotFoundError(err, planID)
		}
		return 0, err
	}

	// Get timeline type for risk_identification
	timelineType, err := rrupids.timelineTypeSvc.GetByType("risk_identification")
	if err != nil {
		return 0, err
	}

	// Get timeline by plan ID and timeline type ID
	timeline, err := rrupids.riskUpdatePlanSvc.GetTimelineByPlanIDAndTypeID(planID, timelineType.ID)
	if err != nil {
		return 0, err
	}

	tx := rrupids.repo.BeginTransaction()

	if err := rrupids.validateRiskRegisterFormDTO(riskRegisterFormDTO); err != nil {
		return 0, err
	}

	dates, err := rrupids.parseUpdatePlanRiskRegisterDates(riskRegisterFormDTO)
	if err != nil {
		return 0, err
	}

	riskMatrix, err := rrupids.riskMatrixSvc.GetByType("erm")
	if err != nil {
		return 0, err
	}

	riskStatus, err := rrupids.riskStatusSvc.GetByType("open")
	if err != nil {
		return 0, err
	}

	riskRegister := &models.RiskRegister{
		RiskCategoryID:           riskRegisterFormDTO.RiskCategoryID,
		RiskSourceID:             riskRegisterFormDTO.RiskSourceID,
		RiskMatrixID:             riskMatrix.ID,
		ResponsibleUnitID:        riskRegisterFormDTO.ResponsibleUnitID,
		ResponsibleUserID:        riskRegisterFormDTO.ResponsibleUserID,
		RiskStatusID:             riskStatus.ID,
		Title:                    riskRegisterFormDTO.Title,
		Description:              riskRegisterFormDTO.Description,
		ReviewFrequencyID:        riskRegisterFormDTO.ReviewFrequencyID,
		LastReviewNote:           riskRegisterFormDTO.LastReviewNote,
		RiskCode:                 riskRegisterFormDTO.RiskCode,
		RegistrationDate:         dates.RegistrationDate,
		DivisionID:               riskRegisterFormDTO.DivisionID,
		SectionID:                riskRegisterFormDTO.SectionID,
		SubSectionID:             riskRegisterFormDTO.SubSectionID,
		IsStrategic:              riskRegisterFormDTO.IsStrategic,
		IsOperational:            riskRegisterFormDTO.IsOperational,
		ExpectedCause:            riskRegisterFormDTO.ExpectedCause,
		Consequence:              riskRegisterFormDTO.Consequence,
		IsInternal:               riskRegisterFormDTO.IsInternal,
		IsExternal:               riskRegisterFormDTO.IsExternal,
		RiskRegisterUpdatePlanID: &planID,
		CreatedBYID:              currentUserID,
	}

	id, err := rrupids.repo.CreateRiskIdentification(riskRegister, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	// Check if this is the first risk for this plan
	riskCount, err := rrupids.repo.CountRisksByPlanID(planID, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	// Update timeline completion and percentage only if this is the first risk
	if riskCount == 1 {
		if err := rrupids.riskUpdatePlanSvc.UpdateTimelineCompletion(timeline.ID, currentUserID, tx); err != nil {
			tx.Rollback()
			return 0, err
		}

		if err := rrupids.riskUpdatePlanSvc.CalculateAndUpdatePercentage(planID, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	tx.Commit()

	go func() {
		existingRiskRegister, err := rrupids.riskRegisterSvc.GetFullDetailsByID(riskRegister.ID)
		if err != nil {
			utility.LogError("Failed to fetch responsible user", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		// log
		err = rrupids.activityLogSvc.LogActivity(
			models.CreateAction,
			models.RiskRegister{}.TableName(),
			id,
			existingRiskRegister,
			currentUserID,
			"risk", "risk_register", "risk_register_update_plan", nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.RiskRegister{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}

		// Prepare common variables for notifications
		urlParam := fmt.Sprintf(`{"id":%d}`, id)
		entity := "risk-register"

		createdBy, err := rrupids.userSvc.GetByID(currentUserID)
		if err != nil {
			utility.LogError("Failed to fetch created by user", map[string]interface{}{
				"error": err.Error(),
			})
		}

		// Send notification to responsible user and/or responsible unit head
		action := "responsible-user-assign"
		var notificationUserIDs []uint
		var emailTo []string

		// If responsible user is assigned, add them to notification
		if riskRegisterFormDTO.ResponsibleUserID != nil {
			notificationUserIDs = append(notificationUserIDs, *riskRegisterFormDTO.ResponsibleUserID)
			responsibleUser, err := rrupids.userSvc.GetByID(*riskRegisterFormDTO.ResponsibleUserID)
			if err != nil {
				utility.LogError("Failed to fetch responsible user", map[string]interface{}{
					"error": err.Error(),
				})
			} else {
				emailTo = append(emailTo, responsibleUser.Email)
			}
		}

		// Also send notification to responsible unit head
		if riskRegisterFormDTO.ResponsibleUnitID != nil {
			headUsers, err := rrupids.organizationUnitsSvc.GetHeadUsersByOrganizationUnitID(*riskRegisterFormDTO.ResponsibleUnitID)
			if err != nil {
				utility.LogError("Failed to get responsible unit head users for notification", map[string]interface{}{
					"unit_id": *riskRegisterFormDTO.ResponsibleUnitID,
					"error":   err.Error(),
				})
			} else {
				for _, headUser := range headUsers {
					// Avoid duplicate notifications if the head user is also the responsible user
					alreadyNotified := false
					for _, notifiedUserID := range notificationUserIDs {
						if notifiedUserID == headUser.ID {
							alreadyNotified = true
							break
						}
					}
					if !alreadyNotified {
						notificationUserIDs = append(notificationUserIDs, headUser.ID)
						emailTo = append(emailTo, headUser.Email)
					}
				}
			}
		}

		// Send notification if there are recipients
		if len(notificationUserIDs) > 0 {
			dbNotificationDTO := &commonDTOs.DBNotificationDTO{
				UserIDs:              notificationUserIDs,
				ModuleIdentifier:     "risk",
				NotificationType:     string(commonDTOs.AssignmentType),
				HeadingEn:            "Risk Register Responsibility Assigned",
				HeadingAr:            "مسؤولية تسجيل المخاطر المعينة",
				TitleEn:              riskRegisterFormDTO.Title,
				Entity:               &entity,
				Action:               &action,
				URLParam:             urlParam,
				IsSystemNotification: false,
				CreatedByID:          &currentUserID,
			}

			// Use responsible user details for email if available, otherwise use unit head
			var firstNameEn, lastNameEn string
			if riskRegisterFormDTO.ResponsibleUserID != nil {
				responsibleUser, _ := rrupids.userSvc.GetByID(*riskRegisterFormDTO.ResponsibleUserID)
				if responsibleUser != nil {
					firstNameEn = responsibleUser.FirstNameEn
					lastNameEn = responsibleUser.LastNameEn
				}
			}
			if firstNameEn == "" && len(emailTo) > 0 {
				// Fallback to first recipient's name
				firstUser, _ := rrupids.userSvc.GetByID(notificationUserIDs[0])
				if firstUser != nil {
					firstNameEn = firstUser.FirstNameEn
					lastNameEn = firstUser.LastNameEn
				}
			}

			subject := fmt.Sprintf("Risk Register/%s: Responsibility Assigned", riskRegister.Title)
			appURL := utility.GetAppURL()
			emailData := map[string]interface{}{
				"FirstNameEn":          firstNameEn,
				"LastNameEn":           lastNameEn,
				"RiskReferenceNumber":  existingRiskRegister.ReferenceNumber,
				"RiskTitle":            existingRiskRegister.Title,
				"RiskDescription":      utility.GetStringValue(existingRiskRegister.Description),
				"CreatedByFirstNameEn": createdBy.FirstNameEn,
				"CreatedByLastNameEn":  createdBy.LastNameEn,
				"AssignedOn":           time.Now().Format("2006-01-02"),
				"RedirectURL":          fmt.Sprintf("%s/risk/identified-risk/%d/workflow", appURL, id),
			}
			emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
				EmailTo:              emailTo,
				Subject:              subject,
				TemplateRelativePath: "risk/templates/email/risk_register_create_responsible_user_mail.html",
				EmailData:            emailData,
			}
			err = rrupids.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "risk")
			if err != nil {
				utility.LogError("Failed to send notifications", map[string]interface{}{
					"error": err.Error(),
				})
			}
		}

		// Fetch risk settings
		riskSetting, err := rrupids.riskSettingsSvc.Get()
		if err != nil {
			utility.LogError("Failed to fetch risk setting", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		// Get Review User IDs
		reviewUserIDs, err := rrupids.userSvc.GetUserIdsByRole(riskSetting.ReviewRoleID)
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

		reviewEmailTo := make([]string, len(reviewUserIDs))
		reviewAction := "review-user-assign"

		for i, userID := range reviewUserIDs {
			user, err := rrupids.userSvc.GetByID(userID)
			if err != nil {
				utility.LogError("Failed to fetch reviewer details", map[string]interface{}{
					"error":  err.Error(),
					"userID": userID,
				})
				continue
			}
			reviewEmailTo[i] = user.Email
		}

		dbNotificationReview := &commonDTOs.DBNotificationDTO{
			UserIDs:              reviewUserIDs,
			ModuleIdentifier:     "risk",
			NotificationType:     string(commonDTOs.InformationType),
			HeadingEn:            "Risk Register Review Assigned",
			HeadingAr:            "تم تعيين مراجعة سجل المخاطر",
			TitleEn:              riskRegister.Title,
			TitleAr:              riskRegister.Title,
			ReferenceNumber:      existingRiskRegister.ReferenceNumber,
			Entity:               &entity,
			Action:               &reviewAction,
			URLParam:             urlParam,
			IsSystemNotification: false,
			CreatedByID:          &currentUserID,
		}

		err = rrupids.notificationSvc.SendNotification(dbNotificationReview, nil, "risk")
		if err != nil {
			utility.LogError("Failed to send DB notification", map[string]interface{}{
				"error": err.Error(),
			})
		}

		emailNotificationReview := &commonDTOs.EmailNotificationDTO{
			EmailTo:              reviewEmailTo,
			Subject:              fmt.Sprintf("Risk Register/%s: Review Assigned", riskRegister.Title),
			TemplateRelativePath: "risk/templates/email/risk_register_review_user_mail.html",
			EmailData: map[string]interface{}{
				"FirstNameEn":          riskSetting.ReviewRole.TitleEn,
				"RiskReferenceNumber":  existingRiskRegister.ReferenceNumber,
				"RiskTitle":            existingRiskRegister.Title,
				"RiskDescription":      utility.GetStringValue(existingRiskRegister.Description),
				"CreatedByFirstNameEn": createdBy.FirstNameEn,
				"CreatedByLastNameEn":  createdBy.LastNameEn,
				"AssignedOn":           time.Now().Format("2006-01-02"),
				"RedirectURL":          fmt.Sprintf("%s/risk/identified-risk/%d/workflow", utility.GetAppURL(), riskRegister.ID),
			},
		}

		err = rrupids.notificationSvc.SendNotification(nil, emailNotificationReview, "risk")
		if err != nil {
			utility.LogError("Failed to send review user email", map[string]interface{}{
				"error":       err.Error(),
				"user_emails": reviewEmailTo,
			})
		}
	}()

	return id, nil
}

func (rrupids *RiskUpdatePlanRiskIdentificationService) UpdateRiskIdentification(planID uint, id uint, currentUserID uint, riskRegisterFormDTO *riskDTOs.UpdatePlanRiskRegisterFormDTO) error {
	// Verify plan exists
	_, err := rrupids.riskUpdatePlanSvc.GetByID(planID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utility.NewResourceNotFoundError(err, planID)
		}
		return err
	}

	tx := rrupids.repo.BeginTransaction()

	oldRiskRegister, err := rrupids.riskRegisterSvc.GetFullDetailsByID(id, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := rrupids.validateRiskRegisterFormDTO(riskRegisterFormDTO); err != nil {
		return err
	}

	dates, err := rrupids.parseUpdatePlanRiskRegisterDates(riskRegisterFormDTO)
	if err != nil {
		return err
	}

	existingRiskRegister, err := rrupids.riskRegisterSvc.GetByID(id)
	if err != nil {
		return err
	}

	// Verify risk belongs to this plan
	if existingRiskRegister.RiskRegisterUpdatePlanID == nil || *existingRiskRegister.RiskRegisterUpdatePlanID != planID {
		tx.Rollback()
		return utility.NewResourceNotFoundError(fmt.Errorf("risk not found in this plan"), id)
	}

	oldResponsibleUserID := existingRiskRegister.ResponsibleUserID

	existingRiskRegister.RiskCategoryID = riskRegisterFormDTO.RiskCategoryID
	existingRiskRegister.RiskSourceID = riskRegisterFormDTO.RiskSourceID
	existingRiskRegister.ResponsibleUserID = riskRegisterFormDTO.ResponsibleUserID
	existingRiskRegister.ResponsibleUnitID = riskRegisterFormDTO.ResponsibleUnitID
	existingRiskRegister.Title = riskRegisterFormDTO.Title
	existingRiskRegister.Description = riskRegisterFormDTO.Description
	existingRiskRegister.ReviewFrequencyID = riskRegisterFormDTO.ReviewFrequencyID
	existingRiskRegister.LastReviewNote = riskRegisterFormDTO.LastReviewNote
	existingRiskRegister.RiskCode = riskRegisterFormDTO.RiskCode
	existingRiskRegister.DivisionID = riskRegisterFormDTO.DivisionID
	existingRiskRegister.SectionID = riskRegisterFormDTO.SectionID
	existingRiskRegister.SubSectionID = riskRegisterFormDTO.SubSectionID
	existingRiskRegister.IsStrategic = riskRegisterFormDTO.IsStrategic
	existingRiskRegister.IsOperational = riskRegisterFormDTO.IsOperational
	existingRiskRegister.ExpectedCause = riskRegisterFormDTO.ExpectedCause
	existingRiskRegister.Consequence = riskRegisterFormDTO.Consequence
	existingRiskRegister.IsInternal = riskRegisterFormDTO.IsInternal
	existingRiskRegister.IsExternal = riskRegisterFormDTO.IsExternal
	existingRiskRegister.RegistrationDate = dates.RegistrationDate
	existingRiskRegister.UpdatedBYID = &currentUserID

	if err := rrupids.repo.UpdateRiskIdentification(existingRiskRegister, tx); err != nil {
		tx.Rollback()
		return err
	}

	existingRisk, err := rrupids.riskRegisterSvc.GetFullDetailsByID(id, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	go func() {

		// Prepare log data
		changedData := make(map[string]interface{})

		// log changes
		oldData := map[string]interface{}{
			"title":            oldRiskRegister.Title,
			"description":      oldRiskRegister.Description,
			"risk_categrory":   oldRiskRegister.RiskCategory.TitleEn,
			"risk_source":      utility.GetFieldValue(oldRiskRegister.RiskSource).(string),
			"reported_at":      oldRiskRegister.ReportedAtString,
			"next_review_date": oldRiskRegister.NextReviewDateString,
		}

		if oldRiskRegister.ResponsibleUser != nil {
			oldData["responsible_user"] = fmt.Sprintf("%s %s", oldRiskRegister.ResponsibleUser.FirstNameEN, oldRiskRegister.ResponsibleUser.LastNameEN)
		}

		newData := map[string]interface{}{
			"title":            existingRisk.Title,
			"description":      existingRisk.Description,
			"risk_categrory":   existingRisk.RiskCategory.TitleEn,
			"risk_source":      utility.GetFieldValue(existingRisk.RiskSource).(string),
			"reported_at":      existingRisk.ReportedAtString,
			"next_review_date": existingRisk.NextReviewDateString,
		}

		if existingRisk.ResponsibleUser != nil {
			newData["responsible_user"] = fmt.Sprintf("%s %s", existingRisk.ResponsibleUser.FirstNameEN, existingRisk.ResponsibleUser.LastNameEN)
		}

		// Compare other fields and add them to changedData
		fieldChanges := utility.CompareFieldChanges(oldData, newData)
		for k, v := range fieldChanges {
			changedData[k] = v
		}

		// Log activity if there are any changes
		if len(changedData) > 0 {
			err := rrupids.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.RiskRegister{}.TableName(),
				id,
				changedData,
				currentUserID,
				"risk", "risk_register", "", nil,
			)
			if err != nil {
				utility.LogError("Failed to log activity", map[string]interface{}{
					"action":   models.UpdateAction,
					"error":    err.Error(),
					"recordID": id,
				})
			}
		}

		// Send notifications to old and new responsible users if the responsible user has changed
		if oldResponsibleUserID != existingRiskRegister.ResponsibleUserID {
			// URL param for notifications
			urlParam := fmt.Sprintf(`{"id":%d}`, id)
			entity := "risk-register"

			// Fetch details of the user making the assignment
			updatedByUser, err := rrupids.userSvc.GetByID(currentUserID)
			if err != nil {
				utility.LogError("Failed to fetch assigning user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}

			appURL := utility.GetAppURL()

			// Send notification to old responsible user if exists
			if oldResponsibleUserID != nil {
				oldResponsibleUserAction := "responsible-user-remove"
				oldResponsibleUserNotificationDTO := &commonDTOs.DBNotificationDTO{
					UserIDs:              []uint{*oldResponsibleUserID},
					ModuleIdentifier:     "risk",
					NotificationType:     string(commonDTOs.InformationType),
					HeadingEn:            "Risk Register Responsibility Removed",
					HeadingAr:            "تمت إزالة مسؤولية تسجيل المخاطر",
					TitleEn:              existingRiskRegister.Title,
					Entity:               &entity,
					Action:               &oldResponsibleUserAction,
					URLParam:             urlParam,
					IsSystemNotification: false,
					CreatedByID:          &currentUserID,
				}

				oldResponsibleUser, err := rrupids.userSvc.GetByID(*oldResponsibleUserID)
				if err != nil {
					utility.LogError("Failed to fetch old responsible user", map[string]interface{}{
						"error": err.Error(),
					})
				} else {
					oldResponsibleUserEmailTo := []string{oldResponsibleUser.Email}
					oldResponsibleUserSubject := fmt.Sprintf("Risk Register/%s: Responsibility Removed", existingRiskRegister.Title)
					oldResponsibleUserEmailData := map[string]interface{}{
						"FirstNameEn":          oldResponsibleUser.FirstNameEn,
						"LastNameEn":           oldResponsibleUser.LastNameEn,
						"RiskReferenceNumber":  existingRiskRegister.ReferenceNumber,
						"RiskTitle":            existingRiskRegister.Title,
						"RiskDescription":      utility.GetStringValue(existingRiskRegister.Description),
						"UpdatedByFirstNameEn": updatedByUser.FirstNameEn,
						"UpdatedByLastNameEn":  updatedByUser.LastNameEn,
						"RedirectURL":          fmt.Sprintf("%s/risk/identified-risk/%d/workflow", appURL, id),
					}

					oldResponsibleUserEmailNotificationDTO := &commonDTOs.EmailNotificationDTO{
						EmailTo:              oldResponsibleUserEmailTo,
						Subject:              oldResponsibleUserSubject,
						TemplateRelativePath: "risk/templates/email/risk_register_remove_responsible_user_mail.html",
						EmailData:            oldResponsibleUserEmailData,
					}

					err = rrupids.notificationSvc.SendNotification(oldResponsibleUserNotificationDTO, oldResponsibleUserEmailNotificationDTO, "risk")
					if err != nil {
						utility.LogError("Failed to send notifications to old responsible user", map[string]interface{}{
							"error": err.Error(),
						})
					}
				}
			}

			// Send notification to new responsible user and/or responsible unit head
			var notificationUserIDs []uint
			var emailTo []string
			newResponsibleUserAction := "responsible-user-assign"

			// If responsible user is assigned, add them to notification
			if riskRegisterFormDTO.ResponsibleUserID != nil {
				notificationUserIDs = append(notificationUserIDs, *riskRegisterFormDTO.ResponsibleUserID)
				newResponsibleUser, err := rrupids.userSvc.GetByID(*riskRegisterFormDTO.ResponsibleUserID)
				if err != nil {
					utility.LogError("Failed to fetch new responsible user", map[string]interface{}{
						"error": err.Error(),
					})
				} else {
					emailTo = append(emailTo, newResponsibleUser.Email)
				}
			}

			// Also send notification to responsible unit head
			if riskRegisterFormDTO.ResponsibleUnitID != nil {
				headUsers, err := rrupids.organizationUnitsSvc.GetHeadUsersByOrganizationUnitID(*riskRegisterFormDTO.ResponsibleUnitID)
				if err != nil {
					utility.LogError("Failed to get responsible unit head users for notification", map[string]interface{}{
						"unit_id": *riskRegisterFormDTO.ResponsibleUnitID,
						"error":   err.Error(),
					})
				} else {
					for _, headUser := range headUsers {
						// Avoid duplicate notifications if the head user is also the responsible user
						alreadyNotified := false
						for _, notifiedUserID := range notificationUserIDs {
							if notifiedUserID == headUser.ID {
								alreadyNotified = true
								break
							}
						}
						if !alreadyNotified {
							notificationUserIDs = append(notificationUserIDs, headUser.ID)
							emailTo = append(emailTo, headUser.Email)
						}
					}
				}
			}

			// Send notification if there are recipients
			if len(notificationUserIDs) > 0 {
				newResponsibleUserNotificationDTO := &commonDTOs.DBNotificationDTO{
					UserIDs:              notificationUserIDs,
					ModuleIdentifier:     "risk",
					NotificationType:     string(commonDTOs.AssignmentType),
					HeadingEn:            "Risk Register Responsibility Assigned",
					HeadingAr:            "مسؤولية تسجيل المخاطر المعينة",
					TitleEn:              existingRiskRegister.Title,
					Entity:               &entity,
					Action:               &newResponsibleUserAction,
					URLParam:             urlParam,
					IsSystemNotification: false,
					CreatedByID:          &currentUserID,
				}

				// Use responsible user details for email if available, otherwise use unit head
				var firstNameEn, lastNameEn string
				if riskRegisterFormDTO.ResponsibleUserID != nil {
					newResponsibleUser, _ := rrupids.userSvc.GetByID(*riskRegisterFormDTO.ResponsibleUserID)
					if newResponsibleUser != nil {
						firstNameEn = newResponsibleUser.FirstNameEn
						lastNameEn = newResponsibleUser.LastNameEn
					}
				}
				if firstNameEn == "" && len(emailTo) > 0 {
					// Fallback to first recipient's name
					firstUser, _ := rrupids.userSvc.GetByID(notificationUserIDs[0])
					if firstUser != nil {
						firstNameEn = firstUser.FirstNameEn
						lastNameEn = firstUser.LastNameEn
					}
				}

				newResponsibleUserSubject := fmt.Sprintf("Risk Register/%s: Responsibility Assigned", existingRiskRegister.Title)
				newResponsibleUserEmailData := map[string]interface{}{
					"FirstNameEn":           firstNameEn,
					"LastNameEn":            lastNameEn,
					"RiskReferenceNumber":   existingRiskRegister.ReferenceNumber,
					"RiskTitle":             existingRiskRegister.Title,
					"RiskDescription":       utility.GetStringValue(existingRiskRegister.Description),
					"AssignedByFirstNameEn": updatedByUser.FirstNameEn,
					"AssignedByLastNameEn":  updatedByUser.LastNameEn,
					"AssignedOn":            time.Now().Format("2006-01-02"),
					"RedirectURL":           fmt.Sprintf("%s/risk/identified-risk/%d/workflow", appURL, id),
				}

				newResponsibleUserEmailNotificationDTO := &commonDTOs.EmailNotificationDTO{
					EmailTo:              emailTo,
					Subject:              newResponsibleUserSubject,
					TemplateRelativePath: "risk/templates/email/risk_register_change_responsible_user_mail.html",
					EmailData:            newResponsibleUserEmailData,
				}

				err = rrupids.notificationSvc.SendNotification(newResponsibleUserNotificationDTO, newResponsibleUserEmailNotificationDTO, "risk")
				if err != nil {
					utility.LogError("Failed to send notifications to new responsible user", map[string]interface{}{
						"error": err.Error(),
					})
				}
			}
		}
	}()

	return nil
}

func (rrupids *RiskUpdatePlanRiskIdentificationService) DeleteRiskIdentification(planID uint, ids []uint, currentUserID uint) ([]uint, error) {
	// Verify plan exists
	_, err := rrupids.riskUpdatePlanSvc.GetByID(planID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, planID)
		}
		return nil, err
	}

	// Get timeline type for risk_identification
	timelineType, err := rrupids.timelineTypeSvc.GetByType("risk_identification")
	if err != nil {
		return nil, err
	}

	// Get timeline by plan ID and timeline type ID
	timeline, err := rrupids.riskUpdatePlanSvc.GetTimelineByPlanIDAndTypeID(planID, timelineType.ID)
	if err != nil {
		return nil, err
	}

	tx := rrupids.repo.BeginTransaction()

	var failedIds []uint
	var successfulIds []uint
	var notifications []struct {
		OwnerID           *uint
		ResponsibleUnitID *uint
		TitleEn           string
		TitleAr           string
		Description       *string
		ReferenceNumber   string
	}

	var deletedRiskRegister []map[string]interface{}

	for _, id := range ids {
		// Get risk model to verify plan ownership
		riskModel, err := rrupids.riskRegisterSvc.GetByID(id)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		// Verify risk belongs to this plan
		if riskModel.RiskRegisterUpdatePlanID == nil || *riskModel.RiskRegisterUpdatePlanID != planID {
			failedIds = append(failedIds, id)
			continue
		}

		risk, err := rrupids.riskRegisterSvc.GetFullDetailsByID(id, tx)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		riskRegister, err := rrupids.riskRegisterSvc.GetFullDetailsByID(id, tx)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		// Prepare notification data before deletion
		var responsibleUnitID *uint
		if riskRegister.ResponsibleUnit != nil {
			unitID := uint(riskRegister.ResponsibleUnit.ID)
			responsibleUnitID = &unitID
		}

		notifications = append(notifications, struct {
			OwnerID           *uint
			ResponsibleUnitID *uint
			TitleEn           string
			TitleAr           string
			Description       *string
			ReferenceNumber   string
		}{
			OwnerID:           &risk.ResponsibleUserID,
			ResponsibleUnitID: responsibleUnitID,
			TitleEn:           risk.Title,
			Description:       risk.Description,
			ReferenceNumber:   risk.ReferenceNumber,
		})

		deletedRiskRegister = append(deletedRiskRegister, map[string]interface{}{
			"id":    id,
			"title": risk.Title,
		})

		// Delete the risk
		if err := rrupids.repo.DeleteRiskIdentification(id, tx); err != nil {
			failedIds = append(failedIds, id)
			continue
		} else {
			successfulIds = append(successfulIds, id)
		}
	}
	if len(failedIds) > 0 {
		tx.Rollback()
		return failedIds, nil
	}

	// Check remaining risks after deletion
	remainingRiskCount, err := rrupids.repo.CountRisksByPlanID(planID, tx)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// If no risks remain, mark timeline as incomplete
	if remainingRiskCount == 0 {
		if err := rrupids.riskUpdatePlanSvc.MarkTimelineAsIncomplete(timeline.ID, tx); err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	// Recalculate percentage
	if err := rrupids.riskUpdatePlanSvc.CalculateAndUpdatePercentage(planID, tx); err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()

	// Log bulk delete activities for successful deletions
	if len(successfulIds) > 0 {
		go func() {
			err := rrupids.activityLogSvc.LogBulkActivities(
				models.DeleteAction, // Make sure this action constant is defined appropriately
				models.RiskRegister{}.TableName(),
				successfulIds,
				deletedRiskRegister,
				currentUserID,
				nil,
			)
			if err != nil {
				// @todo - save to log file
				fmt.Printf("Failed to log bulk delete activity: %v\n", err)
			}

			// Send DB and email notifications for each successfully deleted risk register
			for _, notification := range notifications {
				urlParam := ""
				entity := "risk-register"
				action := "delete"

				var notificationUserIDs []uint
				var emailTo []string

				// If responsible user is assigned, add them to notification
				if notification.OwnerID != nil {
					notificationUserIDs = append(notificationUserIDs, *notification.OwnerID)
					ownerUser, err := rrupids.userSvc.GetByID(*notification.OwnerID)
					if err != nil {
						utility.LogError("Failed to fetch owner details for notification", map[string]interface{}{
							"error": err.Error(),
						})
					} else {
						emailTo = append(emailTo, ownerUser.Email)
					}
				}

				// Also send notification to responsible unit head
				if notification.ResponsibleUnitID != nil {
					headUsers, err := rrupids.organizationUnitsSvc.GetHeadUsersByOrganizationUnitID(*notification.ResponsibleUnitID)
					if err != nil {
						utility.LogError("Failed to get responsible unit head users for notification", map[string]interface{}{
							"unit_id": *notification.ResponsibleUnitID,
							"error":   err.Error(),
						})
					} else {
						for _, headUser := range headUsers {
							// Avoid duplicate notifications if the head user is also the responsible user
							alreadyNotified := false
							for _, notifiedUserID := range notificationUserIDs {
								if notifiedUserID == headUser.ID {
									alreadyNotified = true
									break
								}
							}
							if !alreadyNotified {
								notificationUserIDs = append(notificationUserIDs, headUser.ID)
								emailTo = append(emailTo, headUser.Email)
							}
						}
					}
				}

				// Send notification if there are recipients
				if len(notificationUserIDs) > 0 {
					dbNotificationDTO := &commonDTOs.DBNotificationDTO{
						UserIDs:              notificationUserIDs,
						ModuleIdentifier:     "risk",
						NotificationType:     string(commonDTOs.DeletionType),
						HeadingEn:            "Risk Register Deleted",
						HeadingAr:            "تم حذف سجل المخاطر",
						TitleEn:              notification.TitleEn,
						ReferenceNumber:      notification.ReferenceNumber,
						Entity:               &entity,
						Action:               &action,
						URLParam:             urlParam,
						IsSystemNotification: false,
						CreatedByID:          &currentUserID,
					}

					// Use responsible user details for email if available, otherwise use unit head
					var firstNameEn, lastNameEn string
					if notification.OwnerID != nil {
						ownerUser, _ := rrupids.userSvc.GetByID(*notification.OwnerID)
						if ownerUser != nil {
							firstNameEn = ownerUser.FirstNameEn
							lastNameEn = ownerUser.LastNameEn
						}
					}
					if firstNameEn == "" && len(emailTo) > 0 {
						// Fallback to first recipient's name
						firstUser, _ := rrupids.userSvc.GetByID(notificationUserIDs[0])
						if firstUser != nil {
							firstNameEn = firstUser.FirstNameEn
							lastNameEn = firstUser.LastNameEn
						}
					}

					subject := fmt.Sprintf("Risk Register/%s: Deleted", notification.TitleEn)
					templateRelativePath := "risk/templates/email/risk_register_delete_mail.html"

					emailData := map[string]interface{}{
						"FirstNameEn":         firstNameEn,
						"LastNameEn":          lastNameEn,
						"RiskReferenceNumber": notification.ReferenceNumber,
						"RiskTitle":           notification.TitleEn,
						"RiskDescription":     notification.Description,
					}

					emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
						EmailTo:              emailTo,
						Subject:              subject,
						TemplateRelativePath: templateRelativePath,
						EmailData:            emailData,
					}

					// Send DB and email notifications
					err = rrupids.notificationSvc.SendNotification(
						dbNotificationDTO, emailNotificationDTO, "risk")
					if err != nil {
						utility.LogError("Failed to send deletion notifications", map[string]interface{}{
							"error": err.Error(),
						})
					}
				}
			}
		}()
	}

	return failedIds, nil
}

func (rrupids *RiskUpdatePlanRiskIdentificationService) MarkRiskIdentificationAsNotApplicable(planID uint, currentUserID uint, formDTO *riskDTOs.MarkTimelineNotApplicableDTO) error {
	_, err := rrupids.riskUpdatePlanSvc.GetByID(planID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utility.NewResourceNotFoundError(err, planID)
		}
		return err
	}

	timelineType, err := rrupids.timelineTypeSvc.GetByType("risk_identification")
	if err != nil {
		return err
	}

	timeline, err := rrupids.riskUpdatePlanSvc.GetTimelineByPlanIDAndTypeID(planID, timelineType.ID)
	if err != nil {
		return err
	}

	tx := rrupids.repo.BeginTransaction()

	now := time.Now()
	timeline.IsNotApplicable = true
	timeline.IsCompleted = true
	timeline.Justification = &formDTO.Justification
	timeline.UpdatedBYID = &currentUserID
	timeline.UpdatedAt = &now

	if err := tx.DB().Save(&timeline).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := rrupids.riskUpdatePlanSvc.CalculateAndUpdatePercentage(planID, tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (rrupids *RiskUpdatePlanRiskIdentificationService) validateRiskRegisterFormDTO(riskRegisterFormDTO *riskDTOs.UpdatePlanRiskRegisterFormDTO) error {

	if riskRegisterFormDTO.RiskCategoryID == nil {
		if _, err := rrupids.riskCategorySvc.GetByID(*riskRegisterFormDTO.RiskCategoryID); err != nil {
			return utility.NewValidationError(fmt.Errorf("risk category does not exist"), "not_exist", "risk category")
		}
	}

	if riskRegisterFormDTO.ResponsibleUserID != nil {
		if _, err := rrupids.userSvc.GetByID(*riskRegisterFormDTO.ResponsibleUserID); err != nil {
			return utility.NewValidationError(fmt.Errorf("responsible user does not exist"), "not_exist", "responsible user")
		}
	}

	// Validate that division contains the responsible unit

	if !rrupids.organizationUnitSvc.IsUnitUnderParent(*riskRegisterFormDTO.ResponsibleUnitID, riskRegisterFormDTO.DivisionID) {
		return utility.NewValidationError(fmt.Errorf("responsible unit must be under the selected division"), "invalid_hierarchy", "responsible_unit")
	}

	// Validate that section is under responsible unit
	if riskRegisterFormDTO.SectionID != nil && riskRegisterFormDTO.ResponsibleUnitID != nil {
		if !rrupids.organizationUnitSvc.IsUnitUnderParent(*riskRegisterFormDTO.SectionID, *riskRegisterFormDTO.ResponsibleUnitID) {
			return utility.NewValidationError(fmt.Errorf("section must be under the responsible unit"), "invalid_hierarchy", "section")
		}
	}

	return nil
}

func (rrupids *RiskUpdatePlanRiskIdentificationService) parseRiskRegisterDates(dto *riskDTOs.RiskRegisterFormDTO) (*riskDTOs.RiskRegisterParsedDates, error) {
	parsed := &riskDTOs.RiskRegisterParsedDates{}

	var err error

	if dto.ReportedAt != nil && *dto.ReportedAt != "" {
		parsed.ReportedAt, err = utility.ParseDate(*dto.ReportedAt)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid reported at"), "invalid_date", "reported at")
		}
	}

	if dto.NextReviewDate != nil && *dto.NextReviewDate != "" {
		parsed.NextReviewDate, err = utility.ParseDate(*dto.NextReviewDate)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid next review date"), "invalid_date", "next review date")
		}
	}

	if dto.RegistrationDate != "" {
		parsed.RegistrationDate, err = utility.ParseDate(dto.RegistrationDate)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid registration date"), "invalid_date", "registration date")
		}
	}

	return parsed, nil
}

func (rrupids *RiskUpdatePlanRiskIdentificationService) parseUpdatePlanRiskRegisterDates(dto *riskDTOs.UpdatePlanRiskRegisterFormDTO) (*riskDTOs.RiskRegisterParsedDates, error) {
	parsed := &riskDTOs.RiskRegisterParsedDates{}

	var err error

	if dto.ReportedAt != nil && *dto.ReportedAt != "" {
		parsed.ReportedAt, err = utility.ParseDate(*dto.ReportedAt)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid reported at"), "invalid_date", "reported at")
		}
	}

	if dto.NextReviewDate != nil && *dto.NextReviewDate != "" {
		parsed.NextReviewDate, err = utility.ParseDate(*dto.NextReviewDate)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid next review date"), "invalid_date", "next review date")
		}
	}

	if dto.RegistrationDate != "" {
		parsed.RegistrationDate, err = utility.ParseDate(dto.RegistrationDate)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid registration date"), "invalid_date", "registration date")
		}
	}

	return parsed, nil
}
