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
	settingsServices "jethur/grc/settings/services"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type KRIRegisterServiceInterface interface {
	SetKRIWorkflowService(svc KRIRegisterWorkflowServiceInterface)
	GetAll(params *riskDTOs.KRIRegisterListQueryParamsDTO) (*riskDTOs.KRIRegisterListResponseDTO, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*riskDTOs.KRIRegisterDetailsDTO, error)
	GetByID(id uint, tx ...utility.DatabaseTransaction) (*models.KRIRegister, error)
	Create(currentUserID uint, kriRegisterFormDTO *riskDTOs.KRIRegisterFormDTO) (uint, error)
	Update(id uint, currentUserID uint, kriRegisterFormDTO *riskDTOs.KRIRegisterFormDTO) error
	UpdateValue(kriRegister *models.KRIRegister, tx utility.DatabaseTransaction) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	ForceDelete(ids []uint, currentUserID uint) ([]uint, error)

	ExportKriRegister(data []*riskDTOs.KRIRegisterListDTO, currentUserID uint) (*os.File, error)
	GetTrendAnalysis(kriID uint, limit int) ([]*riskDTOs.KRITrendAnalysisDTO, error)
	GetFrequencyTrendAnalysis(kriID uint, limit int) ([]*riskDTOs.KRIFrequencyTrendAnalysisDTO, error)
	GetAllFrequencyValues(kriRegisterID uint) ([]*riskDTOs.KRIFrequencyValueListDTO, error)

	DetermineRiskLevel(value int, lowThreshold, mediumThreshold *int) (*models.RiskLevel, error)
}

// Service and implementations
type KRIRegisterService struct {
	repo                riskRepositories.KRIRegisterRepositoryInterface
	activityLogSvc      logServices.ActivityLogServiceInterface
	kriCategorySvc      KRICategoryServiceInterface
	userSvc             organizationServices.UserServiceInterface
	riskLevelSvc        RiskLevelServiceInterface
	kriSubCategorySvc   KRISubCategoryServiceInterface
	levelMeasureSvc     LevelOfMeasurementServiceInterface
	distributionSvc     DistributionStatusServiceInterface
	organizationUnitSvc settingsServices.OrganizationUnitServiceInterface
	reviewFrequencySvc  masterServices.ReviewFrequencyServiceInterface
	unitSvc             masterServices.UnitServiceInterface
	userSettingSvc      organizationServices.UserSettingServiceInterface
	notificationSvc     commonServices.NotificationServiceInterface
	kriWorkflowSvc      KRIRegisterWorkflowServiceInterface
}

func NewKRIRegisterService(
	repo riskRepositories.KRIRegisterRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	kriCategorySvc KRICategoryServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	riskLevelSvc RiskLevelServiceInterface,
	kriSubCategorySvc KRISubCategoryServiceInterface,
	levelMeasureSvc LevelOfMeasurementServiceInterface,
	distributionSvc DistributionStatusServiceInterface,
	organizationUnitSvc settingsServices.OrganizationUnitServiceInterface,
	reviewFrequencySvc masterServices.ReviewFrequencyServiceInterface,
	unitSvc masterServices.UnitServiceInterface,
	userSettingSvc organizationServices.UserSettingServiceInterface,
	notificationSvc commonServices.NotificationServiceInterface,
) KRIRegisterServiceInterface {
	return &KRIRegisterService{
		repo:                repo,
		activityLogSvc:      activityLogSvc,
		kriCategorySvc:      kriCategorySvc,
		userSvc:             userSvc,
		riskLevelSvc:        riskLevelSvc,
		kriSubCategorySvc:   kriSubCategorySvc,
		levelMeasureSvc:     levelMeasureSvc,
		distributionSvc:     distributionSvc,
		organizationUnitSvc: organizationUnitSvc,
		reviewFrequencySvc:  reviewFrequencySvc,
		unitSvc:             unitSvc,
		userSettingSvc:      userSettingSvc,
		notificationSvc:     notificationSvc,
	}
}

func (krs *KRIRegisterService) SetKRIWorkflowService(svc KRIRegisterWorkflowServiceInterface) {
	krs.kriWorkflowSvc = svc
}

func (krs *KRIRegisterService) GetAll(params *riskDTOs.KRIRegisterListQueryParamsDTO) (*riskDTOs.KRIRegisterListResponseDTO, error) {
	return krs.repo.GetAll(params)
}

func (krs KRIRegisterService) GetByID(id uint, tx ...utility.DatabaseTransaction) (*models.KRIRegister, error) {
	kriRegister, err := krs.repo.GetByID(id, tx...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return kriRegister, nil
}

func (krs KRIRegisterService) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*riskDTOs.KRIRegisterDetailsDTO, error) {
	var kriRegisterDTO riskDTOs.KRIRegisterDetailsDTO

	kriRegister, err := krs.repo.GetFullDetailsByID(id, tx...)
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
	if kriRegister.SubmittedAt != nil {
		submittedAt := kriRegister.SubmittedAt.Format("2006-01-02 15:04:05")
		kriRegisterDTO.SubmittedAtString = &submittedAt
	}
	if kriRegister.ApprovedAt != nil {
		approvedAt := kriRegister.ApprovedAt.Format("2006-01-02 15:04:05")
		kriRegisterDTO.ApprovedAtString = &approvedAt
	}

	kriRegisterDTO.CreatedAt = kriRegister.CreatedAt
	kriRegisterDTO.UpdatedAt = kriRegister.UpdatedAt

	// Compute actual threshold values from stored percentages and annual_target
	if kriRegister.AnnualTarget != nil {
		annualTarget := *kriRegister.AnnualTarget
		if kriRegister.MinimumSusceptibilityPercent != nil {
			v := annualTarget * (*kriRegister.MinimumSusceptibilityPercent / 100)
			kriRegisterDTO.MinimumSusceptibility = &v
		}
		if kriRegister.MaximumSusceptibilityPercent != nil {
			v := annualTarget * (*kriRegister.MaximumSusceptibilityPercent / 100)
			kriRegisterDTO.MaximumSusceptibility = &v
		}
		if kriRegister.MaximumTolerablePercent != nil {
			v := annualTarget * (*kriRegister.MaximumTolerablePercent / 100)
			kriRegisterDTO.MaximumTolerable = &v
		}
	}

	return &kriRegisterDTO, nil
}

func (krs *KRIRegisterService) Create(currentUserID uint, kriRegisterFormDTO *riskDTOs.KRIRegisterFormDTO) (uint, error) {
	tx := krs.repo.BeginTransaction()

	err := krs.kriWorkflowSvc.ValidateWorkflowBeforeCreate("kri_register", tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	if err := krs.validateKRIRegisterFormDTO(kriRegisterFormDTO, nil); err != nil {
		tx.Rollback()
		return 0, err
	}
	dates, err := krs.parseKRIRegisterDates(kriRegisterFormDTO)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	referenceNumber, err := krs.repo.GenerateReferenceNumber()
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	var organizationUnitID uint
	unitID, err := krs.userSvc.GetUserPrimaryUnitID(kriRegisterFormDTO.ResponsibleUserID)
	if err != nil {
		tx.Rollback()
		return 0, utility.NewValidationError(fmt.Errorf("responsible unit not set forthe responsible user"), "not_exist", "responsible unit")
	}
	organizationUnitID = unitID

	// var riskLevelID *uint
	// if kriRegisterFormDTO.CurrentValue != nil {
	// 	riskLevel, err := krs.DetermineRiskLevel(*kriRegisterFormDTO.CurrentValue, kriRegisterFormDTO.LowThreshold, kriRegisterFormDTO.MediumThreshold)
	// 	if err != nil {
	// 		tx.Rollback()
	// 		return 0, err
	// 	}
	// 	riskLevelID = &riskLevel.ID
	// }

	kriRegister := &models.KRIRegister{
		ReferenceNumber:   &referenceNumber,
		KRICategoryID:     kriRegisterFormDTO.KRICategoryID,
		ResponsibleUnitID: organizationUnitID,
		ResponsibleUserID: kriRegisterFormDTO.ResponsibleUserID,
		ReviewFrequencyID: kriRegisterFormDTO.ReviewFrequencyID,
		// RiskLevelID:                  riskLevelID,
		UnitID:                       &kriRegisterFormDTO.UnitID,
		LevelOfMeasurementID:         kriRegisterFormDTO.LevelOfMeasurementID,
		RiskResponsibleUnitID:        kriRegisterFormDTO.RiskResponsibleUnitID,
		KRISubCategoryID:             kriRegisterFormDTO.KRISubCategoryID,
		AnnualTarget:                 kriRegisterFormDTO.AnnualTarget,
		ReviewTarget:                 kriRegisterFormDTO.ReviewTarget,
		DistributionStatusID:         kriRegisterFormDTO.DistributionStatusID,
		Title:                        kriRegisterFormDTO.Title,
		Description:                  &kriRegisterFormDTO.Description,
		CalculationType:              kriRegisterFormDTO.CalculationType,
		LowThreshold:                 kriRegisterFormDTO.LowThreshold,
		MediumThreshold:              kriRegisterFormDTO.MediumThreshold,
		CurrentValue:                 kriRegisterFormDTO.CurrentValue,
		NextReviewDate:               *dates.NextReviewDate,
		LastReviewNote:               kriRegisterFormDTO.LastReviewNote,
		MinimumSusceptibilityPercent: kriRegisterFormDTO.MinimumSusceptibilityPercent,
		MaximumSusceptibilityPercent: kriRegisterFormDTO.MaximumSusceptibilityPercent,
		MaximumTolerablePercent:      kriRegisterFormDTO.MaximumTolerablePercent,
		IsLowerBetter:                kriRegisterFormDTO.IsLowerBetter,
		IsHigherBetter:               kriRegisterFormDTO.IsHigherBetter,
		CreatedBYID:                  currentUserID,
		UpdatedAt:                    nil,
	}
	if kriRegisterFormDTO.LastReviewedDate != "" {
		kriRegister.LastReviewedDate = dates.LastReviewedDate
	}

	id, err := krs.repo.Create(kriRegister, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	err = krs.kriWorkflowSvc.CreateKRIRegisterWorkflow(id, "kri_register", currentUserID, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	initialValue := "Initial kri value"
	kriValueUpdate := &models.KRIValueUpdate{
		KRIRegisterID:   id,
		UnitID:          kriRegisterFormDTO.UnitID,
		ReviewedDate:    *dates.NextReviewDate,
		Value:           kriRegisterFormDTO.CurrentValue,
		LowThreshold:    kriRegisterFormDTO.LowThreshold,
		MediumThreshold: kriRegisterFormDTO.MediumThreshold,
		Remarks:         &initialValue,
		CreatedByID:     currentUserID,
	}
	_, err = krs.repo.CreateInitialKriValueUpdate(kriValueUpdate, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	// Process frequency values if provided
	if len(kriRegisterFormDTO.FrequencyValues) > 0 {
		err = krs.processFrequencyValues(id, kriRegisterFormDTO.ReviewFrequencyID, kriRegisterFormDTO.FrequencyValues, currentUserID, tx)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	go func() {
		existingKriRegister, err := krs.GetFullDetailsByID(kriRegister.ID)
		if err != nil {
			utility.LogError("Failed to fetch kri register", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		err = krs.activityLogSvc.LogActivity(
			models.CreateAction,
			models.KRIRegister{}.TableName(),
			id,
			existingKriRegister,
			currentUserID,
			"risk", "kri", "", nil,
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

		urlParam := fmt.Sprintf(`{"id":%d}`, id)
		entity := "kri-register"
		action := "responsible-user-assign"

		dbNotificationDTO := &commonDTOs.DBNotificationDTO{
			UserIDs:              []uint{kriRegisterFormDTO.ResponsibleUserID},
			ModuleIdentifier:     "risk",
			NotificationType:     string(commonDTOs.AssignmentType),
			HeadingEn:            "KRI Register Responsibililty Assigned",
			HeadingAr:            "KRI Register Responsibililty Assigned",
			TitleEn:              kriRegisterFormDTO.Title,
			ReferenceNumber:      referenceNumber,
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
			"RedirectURL":          fmt.Sprintf("%s/risk/kri-register/%d/overview", appURL, id),
		}

		emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
			EmailTo:              emailTo,
			Subject:              subject,
			TemplateRelativePath: "risk/templates/email/kri_register_create_responsible_user_mail.html",
			EmailData:            emailData,
		}

		err = krs.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "kri_register")
		if err != nil {
			utility.LogError("Failed to send notifications", map[string]interface{}{
				"error": err.Error(),
			})
		}
	}()

	return id, nil
}

func (krs *KRIRegisterService) Update(id uint, currentUserID uint, kriRegisterFormDTO *riskDTOs.KRIRegisterFormDTO) error {
	oldKriRegister, err := krs.GetFullDetailsByID(id)
	if err != nil {
		return err
	}

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
	existingKRIRegister.UnitID = &kriRegisterFormDTO.UnitID
	existingKRIRegister.LevelOfMeasurementID = kriRegisterFormDTO.LevelOfMeasurementID
	existingKRIRegister.RiskResponsibleUnitID = kriRegisterFormDTO.RiskResponsibleUnitID
	existingKRIRegister.KRISubCategoryID = kriRegisterFormDTO.KRISubCategoryID
	existingKRIRegister.AnnualTarget = kriRegisterFormDTO.AnnualTarget
	existingKRIRegister.ReviewTarget = kriRegisterFormDTO.ReviewTarget
	existingKRIRegister.DistributionStatusID = kriRegisterFormDTO.DistributionStatusID
	existingKRIRegister.Title = kriRegisterFormDTO.Title
	existingKRIRegister.Description = &kriRegisterFormDTO.Description
	existingKRIRegister.CalculationType = kriRegisterFormDTO.CalculationType
	existingKRIRegister.LowThreshold = kriRegisterFormDTO.LowThreshold
	existingKRIRegister.MediumThreshold = kriRegisterFormDTO.MediumThreshold
	existingKRIRegister.MinimumSusceptibilityPercent = kriRegisterFormDTO.MinimumSusceptibilityPercent
	existingKRIRegister.MaximumSusceptibilityPercent = kriRegisterFormDTO.MaximumSusceptibilityPercent
	existingKRIRegister.MaximumTolerablePercent = kriRegisterFormDTO.MaximumTolerablePercent
	existingKRIRegister.IsLowerBetter = kriRegisterFormDTO.IsLowerBetter
	existingKRIRegister.IsHigherBetter = kriRegisterFormDTO.IsHigherBetter

	if kriRegisterFormDTO.LastReviewedDate != "" {
		existingKRIRegister.LastReviewedDate = dates.LastReviewedDate
	}
	existingKRIRegister.NextReviewDate = *dates.NextReviewDate
	existingKRIRegister.LastReviewNote = kriRegisterFormDTO.LastReviewNote
	existingKRIRegister.UpdatedBYID = &currentUserID

	tx := krs.repo.BeginTransaction()

	err = krs.repo.Update(existingKRIRegister, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Handle frequency values update
	if len(kriRegisterFormDTO.FrequencyValues) > 0 {
		// Delete existing frequency values
		err = krs.repo.DeleteFrequencyValuesByKRIID(id, tx)
		if err != nil {
			tx.Rollback()
			return err
		}

		// Create new frequency values
		err = krs.processFrequencyValues(id, kriRegisterFormDTO.ReviewFrequencyID, kriRegisterFormDTO.FrequencyValues, currentUserID, tx)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	existingKri, err := krs.GetFullDetailsByID(id)
	if err != nil {
		return err
	}

	go func() {
		changedData := make(map[string]interface{})
		// log changes
		oldData := map[string]interface{}{
			"title":                          oldKriRegister.Title,
			"description":                    oldKriRegister.Description,
			"kri_category":                   oldKriRegister.KRICategory.TitleEn,
			"responsible_user":               fmt.Sprintf("%s %s", oldKriRegister.ResponsibleUser.FirstNameEN, oldKriRegister.ResponsibleUser.LastNameEN),
			"review_frequency":               oldKriRegister.ReviewFrequency.TitleEn,
			"risk_level":                     utility.GetFieldValue(oldKriRegister.RiskLevel).(string),
			"unit":                           oldKriRegister.Unit,
			"level_of_measurement":           utility.GetFieldValue(oldKriRegister.LevelOfMeasurement),
			"risk_responsible_unit":          utility.GetFieldValue(oldKriRegister.RiskResponsibleUnit),
			"kri_sub_category":               utility.GetFieldValue(oldKriRegister.KRISubCategory),
			"annual_target":                  oldKriRegister.AnnualTarget,
			"review_target":                  oldKriRegister.ReviewTarget,
			"distribution_status":            utility.GetFieldValue(oldKriRegister.DistributionStatus),
			"calculation_type":               oldKriRegister.CalculationType,
			"low_threshold":                  oldKriRegister.LowThreshold,
			"medium_threshold":               oldKriRegister.MediumThreshold,
			"high_threshold":                 oldKriRegister.HighThreshold,
			"current_value":                  oldKriRegister.CurrentValue,
			"last_review_note":               oldKriRegister.LastReviewNote,
			"minimum_susceptibility_percent": oldKriRegister.MinimumSusceptibilityPercent,
			"maximum_susceptibility_percent": oldKriRegister.MaximumSusceptibilityPercent,
			"maximum_tolerable_percent":      oldKriRegister.MaximumTolerablePercent,
		}

		newData := map[string]interface{}{
			"title":                          existingKri.Title,
			"description":                    existingKri.Description,
			"kri_category":                   existingKri.KRICategory.TitleEn,
			"responsible_user":               fmt.Sprintf("%s %s", existingKri.ResponsibleUser.FirstNameEN, existingKri.ResponsibleUser.LastNameEN),
			"review_frequency":               existingKri.ReviewFrequency.TitleEn,
			"risk_level":                     utility.GetFieldValue(existingKri.RiskLevel).(string),
			"unit":                           existingKri.Unit,
			"level_of_measurement":           utility.GetFieldValue(existingKri.LevelOfMeasurement),
			"risk_responsible_unit":          utility.GetFieldValue(existingKri.RiskResponsibleUnit),
			"kri_sub_category":               utility.GetFieldValue(existingKri.KRISubCategory),
			"annual_target":                  existingKri.AnnualTarget,
			"review_target":                  existingKri.ReviewTarget,
			"distribution_status":            utility.GetFieldValue(existingKri.DistributionStatus),
			"calculation_type":               existingKri.CalculationType,
			"low_threshold":                  existingKri.LowThreshold,
			"medium_threshold":               existingKri.MediumThreshold,
			"high_threshold":                 existingKri.HighThreshold,
			"current_value":                  existingKri.CurrentValue,
			"last_review_note":               existingKri.LastReviewNote,
			"minimum_susceptibility_percent": existingKri.MinimumSusceptibilityPercent,
			"maximum_susceptibility_percent": existingKri.MaximumSusceptibilityPercent,
			"maximum_tolerable_percent":      existingKri.MaximumTolerablePercent,
		}
		// Compare other fields and add them to changedData
		fieldChanges := utility.CompareFieldChanges(oldData, newData)
		for k, v := range fieldChanges {
			changedData[k] = v
		}
		// Log activity if there are any changes
		if len(changedData) > 0 {
			// Log the update activity
			err := krs.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.KRIRegister{}.TableName(),
				id,
				changedData,
				currentUserID,
				"risk", "kri", "", nil,
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
		}

		// Send notifications to old and new responsible users if the responsible user has changed
		if oldResponsibleUserID != existingKRIRegister.ResponsibleUserID {
			// URL param for notifications
			urlParam := fmt.Sprintf(`{"id":%d}`, id)
			entity := "kri-register"

			// Old responsible user notification (Removal)
			oldResponsibleUserAction := "responsible-user-remove"
			oldResponsibleUserNotificationDTO := &commonDTOs.DBNotificationDTO{
				UserIDs:              []uint{oldResponsibleUserID},
				ModuleIdentifier:     "risk",
				NotificationType:     string(commonDTOs.InformationType),
				HeadingEn:            "Kri Register Responsibility Removed",
				HeadingAr:            "تمت إزالة مسؤولية تسجيل Kri",
				TitleEn:              existingKRIRegister.Title,
				ReferenceNumber:      *existingKRIRegister.ReferenceNumber,
				Entity:               &entity,
				Action:               &oldResponsibleUserAction,
				URLParam:             urlParam,
				IsSystemNotification: false,
				CreatedByID:          &currentUserID,
			}

			// New responsible user notification (Assignment)
			newResponsibleUserAction := "responsible-user-assign"
			newResponsibleUserNotificationDTO := &commonDTOs.DBNotificationDTO{
				UserIDs:              []uint{kriRegisterFormDTO.ResponsibleUserID},
				ModuleIdentifier:     "risk",
				NotificationType:     string(commonDTOs.AssignmentType),
				HeadingEn:            "Kri Register Responsibility Assigned",
				HeadingAr:            "Kri تسجيل المسؤولية المعينة",
				TitleEn:              existingKRIRegister.Title,
				ReferenceNumber:      *existingKRIRegister.ReferenceNumber,
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
				"RedirectURL":          fmt.Sprintf("%s/risk/kri-register/%d/overview", appURL, id),
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
				"RedirectURL":           fmt.Sprintf("%s/risk/kri-register/%d/overview", appURL, id),
			}

			newResponsibleUserEmailNotificationDTO := &commonDTOs.EmailNotificationDTO{
				EmailTo:              newResponsibleUserEmailTo,
				Subject:              newResponsibleUserSubject,
				TemplateRelativePath: "risk/templates/email/kri_register_change_responsible_user_mail.html",
				EmailData:            newResponsibleUserEmailData,
			}

			// Send notifications and emails
			err = krs.notificationSvc.SendNotification(newResponsibleUserNotificationDTO, newResponsibleUserEmailNotificationDTO, "kri_register")
			if err != nil {
				utility.LogError("Failed to send notifications to new responsible user", map[string]interface{}{
					"error": err.Error(),
				})
			}

			err = krs.notificationSvc.SendNotification(oldResponsibleUserNotificationDTO, oldResponsibleUserEmailNotificationDTO, "kri_register")
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
	tx := krs.repo.BeginTransaction()

	var failedIds []uint
	var successfulIds []uint

	var notifications []struct {
		OwnerID         *uint
		TitleEn         string
		TitleAr         string
		Description     string
		ReferenceNumber string
	}

	var deletedKriRegister []map[string]interface{}

	for _, id := range ids {
		_, err := krs.repo.GetByID(id, tx)
		if err != nil {
			continue
		}

		kriRegister, err := krs.repo.GetFullDetailsByID(id, tx)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		if len(kriRegister.ValueUpdates) > 0 {
			var deleteStatusUpdateIds []*uint
			for _, update := range kriRegister.ValueUpdates {
				deleteStatusUpdateIds = append(deleteStatusUpdateIds, &update.ID)
			}
			if err := krs.repo.DeleteKriValueUpdates(deleteStatusUpdateIds, tx); err != nil {
				failedIds = append(failedIds, id)
				continue
			}
		}

		var deleteKRIWorkflowIds []*uint
		var deleteKRIWorkflowParticipantIds []*uint
		var deleteKRIWorkflowReviewIds []*uint
		if len(kriRegister.KRIRegisterWorkflows) > 0 {
			for _, kriWorkflow := range kriRegister.KRIRegisterWorkflows {
				deleteKRIWorkflowIds = append(deleteKRIWorkflowIds, &kriWorkflow.ID)

				if len(kriWorkflow.KRIRegisterWorkflowParticipants) > 0 {
					for _, participant := range kriWorkflow.KRIRegisterWorkflowParticipants {
						deleteKRIWorkflowParticipantIds = append(deleteKRIWorkflowParticipantIds, &participant.ID)
					}
				}

				if len(kriWorkflow.KRIRegisterWorkflowReviews) > 0 {
					for _, review := range kriWorkflow.KRIRegisterWorkflowReviews {
						deleteKRIWorkflowReviewIds = append(deleteKRIWorkflowReviewIds, &review.ID)
					}
				}
			}
		}

		if err := krs.kriWorkflowSvc.DeleteKRIRegisterWorkflowParticipants(deleteKRIWorkflowParticipantIds, tx); err != nil {
			failedIds = append(failedIds, id)
			continue
		}

		if err := krs.kriWorkflowSvc.DeleteKRIRegisterReviews(deleteKRIWorkflowReviewIds, tx); err != nil {
			failedIds = append(failedIds, id)
			continue
		}

		if err := krs.kriWorkflowSvc.DeleteKRIRegisterWorkflows(deleteKRIWorkflowIds, tx); err != nil {
			failedIds = append(failedIds, id)
			continue
		}

		// Delete frequency values
		if err := krs.repo.DeleteFrequencyValuesByKRIID(id, tx); err != nil {
			failedIds = append(failedIds, id)
			continue
		}

		err = krs.repo.Delete(id, tx)
		if err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}

		notifications = append(notifications, struct {
			OwnerID         *uint
			TitleEn         string
			TitleAr         string
			Description     string
			ReferenceNumber string
		}{
			OwnerID:         &kriRegister.ResponsibleUserID,
			TitleEn:         kriRegister.Title,
			Description:     utility.GetStringValue(kriRegister.Description),
			ReferenceNumber: utility.GetStringValue(kriRegister.ReferenceNumber),
		})

		deletedKriRegister = append(deletedKriRegister, map[string]interface{}{
			"id":    id,
			"title": kriRegister.Title,
		})
	}

	if len(failedIds) > 0 {
		tx.Rollback()
		return failedIds, nil
	}

	if err := tx.Commit(); err != nil {
		return ids, err
	}

	if len(successfulIds) > 0 {
		go func() {
			err := krs.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.KRIRegister{}.TableName(),
				successfulIds,
				deletedKriRegister,
				currentUserID,
				nil,
			)
			if err != nil {
				fmt.Printf("Failed to log bulk delete activity: %v\n", err)
			}

			for _, notification := range notifications {
				urlParam := ""
				entity := "kri-register"
				action := "delete"

				dbNotificationDTO := &commonDTOs.DBNotificationDTO{
					UserIDs:              []uint{*notification.OwnerID},
					ModuleIdentifier:     "risk",
					NotificationType:     string(commonDTOs.DeletionType),
					HeadingEn:            "KRI Register Deleted",
					HeadingAr:            "KRI Register Deleted",
					TitleEn:              notification.TitleEn,
					TitleAr:              notification.TitleAr,
					ReferenceNumber:      notification.ReferenceNumber,
					Entity:               &entity,
					Action:               &action,
					URLParam:             urlParam,
					IsSystemNotification: false,
					CreatedByID:          &currentUserID,
				}

				ownerUser, err := krs.userSvc.GetByID(*notification.OwnerID)
				if err != nil {
					utility.LogError("Failed to fetch owner details for notification", map[string]interface{}{
						"error": err.Error(),
					})
					continue
				}

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

				err = krs.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "kri_register")
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
	if _, err := krs.unitSvc.GetByID(kriRegisterFormDTO.UnitID); err != nil {
		return utility.NewValidationError(fmt.Errorf("unit does not exist"), "not_exist", "unit")
	}

	if _, err := krs.reviewFrequencySvc.GetByID(kriRegisterFormDTO.ReviewFrequencyID); err != nil {
		return utility.NewValidationError(fmt.Errorf("review frequency does not exist"), "not_exist", "review frequency")
	}

	if kriRegisterFormDTO.LevelOfMeasurementID != nil {
		if _, err := krs.levelMeasureSvc.GetByID(*kriRegisterFormDTO.LevelOfMeasurementID); err != nil {
			return utility.NewValidationError(fmt.Errorf("level of measurement does not exist"), "not_exist", "level of measurement")
		}
	}

	if kriRegisterFormDTO.RiskResponsibleUnitID != nil {
		if _, err := krs.organizationUnitSvc.GetByID(*kriRegisterFormDTO.RiskResponsibleUnitID); err != nil {
			return utility.NewValidationError(fmt.Errorf("risk responsible unit does not exist"), "not_exist", "risk responsible unit")
		}
	}

	if kriRegisterFormDTO.KRISubCategoryID != nil {
		if _, err := krs.kriSubCategorySvc.GetByID(*kriRegisterFormDTO.KRISubCategoryID); err != nil {
			return utility.NewValidationError(fmt.Errorf("kri sub category does not exist"), "not_exist", "kri sub category")
		}
	}

	if kriRegisterFormDTO.DistributionStatusID != nil {
		if _, err := krs.distributionSvc.GetByID(*kriRegisterFormDTO.DistributionStatusID); err != nil {
			return utility.NewValidationError(fmt.Errorf("distribution status does not exist"), "not_exist", "distribution status")
		}
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
			unit = utility.GetStringValue(kri.UnitTitleAr)
			riskLevel = utility.GetStringValue(kri.RiskLevelTitleAr)
			reviewFrequency = kri.ReviewFrequencyTitleAr
		} else {
			category = kri.KRICategoryTitleEn
			responsibleUser = kri.ResponsibleUserFirstNameEn + " " + kri.ResponsibleUserLastNameEn
			unit = utility.GetStringValue(kri.UnitTitleEn)
			riskLevel = utility.GetStringValue(kri.RiskLevelTitleEn)
			reviewFrequency = kri.ReviewFrequencyTitleEn
		}

		rows[i] = []interface{}{
			i + 1,
			utility.GetStringValue(kri.ReferenceNumber),
			kri.Title,
			utility.GetStringValue(kri.Description),
			category,
			kri.CalculationType,
			responsibleUser,
			fmt.Sprintf("<%s", utility.GetIntValueAsString(kri.LowThreshold)),
			fmt.Sprintf("%s-%s", utility.GetIntValueAsString(kri.LowThreshold), utility.GetIntValueAsString(kri.MediumThreshold)),
			fmt.Sprintf(">%s", utility.GetIntValueAsString(kri.MediumThreshold)),
			utility.GetIntValueAsString(kri.CurrentValue),
			unit,
			riskLevel,
			reviewFrequency,
			utility.GetStringValue(kri.LastReviewedDate),
			kri.NextReviewDate,
		}
	}

	// Generate the Excel file
	return utility.GenerateExcelFile(headers, rows)
}

func (krs *KRIRegisterService) DetermineRiskLevel(value int, lowThreshold, mediumThreshold *int) (*models.RiskLevel, error) {
	var riskLevelType string

	// Determine risk level based on the thresholds
	if lowThreshold != nil && value < *lowThreshold {
		riskLevelType = "low"
	} else if mediumThreshold != nil && value < *mediumThreshold {
		riskLevelType = "medium"
	} else if mediumThreshold != nil && value >= *mediumThreshold {
		riskLevelType = "high"
	} else {
		// Log an error if thresholds are not properly defined
		utility.LogError("Undefined thresholds", map[string]interface{}{
			"value":           value,
			"LowThreshold":    lowThreshold,
			"MediumThreshold": mediumThreshold,
		})
		return nil, fmt.Errorf("undefined thresholds for value: %d", value)
	}

	// Fetch RiskLevel model from the database
	riskLevel, err := krs.riskLevelSvc.GetByType(riskLevelType)
	if err != nil {
		// Log the error and return it
		utility.LogError("Failed to fetch RiskLevel", map[string]interface{}{
			"error":         err.Error(),
			"riskLevelType": riskLevelType,
		})
		return nil, err
	}

	return riskLevel, nil
}

// Old implementation - kept for reference
func (krs *KRIRegisterService) GetTrendAnalysis(kriID uint, limit int) ([]*riskDTOs.KRITrendAnalysisDTO, error) {
	kriValueUpdates, err := krs.repo.GetTrendAnalysis(kriID, limit)
	if err != nil {
		return nil, err
	}

	if len(kriValueUpdates) == 0 {
		return []*riskDTOs.KRITrendAnalysisDTO{}, nil
	}

	// Already sorted and limited in repo (DESC), so just map to DTOs
	var trend []*riskDTOs.KRITrendAnalysisDTO
	for _, valueUpdate := range kriValueUpdates {
		var kriUpdate riskDTOs.KRITrendAnalysisDTO
		err := mapstructure.Decode(valueUpdate, &kriUpdate)
		if err != nil {
			return nil, err
		}
		if valueUpdate.Value != nil {
			kriUpdate.Value = float64(*valueUpdate.Value)
		} else {
			kriUpdate.Value = 0
		}
		reviewedAt := valueUpdate.ReviewedDate.Format("2006-01-02")
		kriUpdate.ReviewedAtString = reviewedAt
		trend = append(trend, &kriUpdate)
	}

	return trend, nil
}

// New implementation - fetches from kri_frequency_value with month-wise data and risk tolerance status
func (krs *KRIRegisterService) GetFrequencyTrendAnalysis(kriID uint, limit int) ([]*riskDTOs.KRIFrequencyTrendAnalysisDTO, error) {
	frequencyValues, err := krs.repo.GetFrequencyTrendAnalysis(kriID, limit)
	if err != nil {
		return nil, err
	}

	if len(frequencyValues) == 0 {
		return []*riskDTOs.KRIFrequencyTrendAnalysisDTO{}, nil
	}

	// Map frequency values to DTOs with risk tolerance status details
	var trend []*riskDTOs.KRIFrequencyTrendAnalysisDTO
	for _, fv := range frequencyValues {
		trendDTO := &riskDTOs.KRIFrequencyTrendAnalysisDTO{
			ID:               fv.ID,
			KRIRegisterID:    fv.KRIRegisterID,
			Year:             fv.Year,
			MonthEn:          fv.MonthEn,
			MonthAr:          fv.MonthAr,
			Type:             fv.Type,
			TargetValue:      fv.TargetValue,
			ActualValue:      fv.ActualValue,
			AchievementValue: fv.AchievementValue,
			IsAchieved:       fv.IsAchieved,
			IsMissed:         fv.IsMissed,
			CreatedAt:        fv.CreatedAt.Format("2006-01-02 15:04:05"),
		}

		// Add risk tolerance status details if available
		if fv.RiskToleranceStatus != nil {
			trendDTO.RiskToleranceStatusID = &fv.RiskToleranceStatus.ID
			trendDTO.RiskToleranceStatusEn = &fv.RiskToleranceStatus.TitleEn
			trendDTO.RiskToleranceStatusAr = &fv.RiskToleranceStatus.TitleAr
			trendDTO.RiskToleranceColorCode = &fv.RiskToleranceStatus.ColorCode
			trendDTO.RiskToleranceType = &fv.RiskToleranceStatus.Type
		}

		trend = append(trend, trendDTO)
	}

	return trend, nil
}

func (krs *KRIRegisterService) ForceDelete(ids []uint, currentUserID uint) ([]uint, error) {
	tx := krs.repo.BeginTransaction()
	var failedIds []uint
	var deletedData []map[string]interface{}

	for _, id := range ids {
		if err := krs.repo.ForceDelete(id, tx); err != nil {
			failedIds = append(failedIds, id)
			continue
		}

		deletedData = append(deletedData, map[string]interface{}{
			"id": id,
		})
	}

	if len(failedIds) > 0 {
		return failedIds, nil
	}

	if err := tx.Commit(); err != nil {
		return ids, err
	}

	go func() {
		if err := krs.activityLogSvc.LogBulkActivities(
			models.DeleteAction,
			models.KRIRegister{}.TableName(),
			ids,
			deletedData,
			currentUserID,
			nil,
		); err != nil {
			fmt.Printf("Failed to log bulk delete activity: %v\n", err)
		}
	}()

	return nil, nil
}

// Helper function to get month name in English and Arabic
func getMonthNames(monthType string) (string, string) {
	monthMap := map[string][2]string{
		"jan": {"Jan", "يناير"},
		"feb": {"Feb", "فبراير"},
		"mar": {"Mar", "مارس"},
		"apr": {"Apr", "أبريل"},
		"may": {"May", "مايو"},
		"jun": {"Jun", "يونيو"},
		"jul": {"Jul", "يوليو"},
		"aug": {"Aug", "أغسطس"},
		"sep": {"Sep", "سبتمبر"},
		"oct": {"Oct", "أكتوبر"},
		"nov": {"Nov", "نوفمبر"},
		"dec": {"Dec", "ديسمبر"},
	}

	if names, ok := monthMap[monthType]; ok {
		return names[0], names[1]
	}
	return "", ""
}

// Helper function to process frequency values based on review frequency type
func (krs *KRIRegisterService) processFrequencyValues(kriRegisterID uint, reviewFrequencyID uint, frequencyValues []*riskDTOs.FrequencyValueFormDTO, currentUserID uint, tx utility.DatabaseTransaction) error {
	// Get review frequency to check type
	reviewFrequency, err := krs.reviewFrequencySvc.GetByID(reviewFrequencyID)
	if err != nil {
		return err
	}

	// Only process if review frequency is monthly or quarterly
	if reviewFrequency.Type != "monthly" && reviewFrequency.Type != "quarterly" {
		return nil
	}

	// Validate frequency values based on type
	if reviewFrequency.Type == "monthly" && len(frequencyValues) != 12 {
		return utility.NewValidationError(fmt.Errorf("monthly frequency requires 12 values"), "invalid_count", "frequency_values")
	}
	if reviewFrequency.Type == "quarterly" && len(frequencyValues) != 4 {
		return utility.NewValidationError(fmt.Errorf("quarterly frequency requires 4 values (mar, jun, sep, dec)"), "invalid_count", "frequency_values")
	}

	// Create frequency value records
	var kriFrequencyValues []*models.KRIFrequencyValue
	currentYear := time.Now().Year()
	for _, fv := range frequencyValues {
		monthEn, monthAr := getMonthNames(fv.Type)
		if monthEn == "" {
			return utility.NewValidationError(fmt.Errorf("invalid month type: %s", fv.Type), "invalid_type", "frequency_values")
		}

		kriFrequencyValue := &models.KRIFrequencyValue{
			KRIRegisterID: kriRegisterID,
			Year:          currentYear,
			MonthEn:       monthEn,
			MonthAr:       monthAr,
			Type:          fv.Type,
			TargetValue:   fv.TargetValue,
			CreatedBYID:   &currentUserID,
			CreatedAt:     time.Now(),
		}
		kriFrequencyValues = append(kriFrequencyValues, kriFrequencyValue)
	}

	return krs.repo.CreateFrequencyValues(kriFrequencyValues, tx)
}

func (krs *KRIRegisterService) GetAllFrequencyValues(kriRegisterID uint) ([]*riskDTOs.KRIFrequencyValueListDTO, error) {
	return krs.repo.GetAllFrequencyValues(kriRegisterID)
}
