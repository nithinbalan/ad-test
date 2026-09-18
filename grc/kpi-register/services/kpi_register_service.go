package services

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	commonDTOs "jethur/grc/common/dtos"
	commonServices "jethur/grc/common/services"
	"jethur/grc/kpi-register/dtos"
	kpiDTOs "jethur/grc/kpi-register/dtos"
	kpiRegisterRepositories "jethur/grc/kpi-register/repositories"
	logServices "jethur/grc/logs/services"
	masters "jethur/grc/masters/services"
	"jethur/grc/models"
	userDTOs "jethur/grc/organization/dtos"
	organizationServices "jethur/grc/organization/services"
	settingServices "jethur/grc/settings/services"
	"jethur/utility"

	"github.com/Knetic/govaluate"
	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type KPIRegisterServiceInterface interface {
	GetAll(params *kpiDTOs.KPIRegisterListQueryParamsDTO) (*kpiDTOs.KPIRegisterListResponseDTO, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*kpiDTOs.KPIRegisterDetailDTO, error)
	GetByID(id uint) (*models.KPIRegister, error)
	Create(currentUserID uint, kpiDTOs *kpiDTOs.KPIRegisterFormDTO) (uint, error)
	Update(id uint, currentUserID uint, kpiRegister *kpiDTOs.KPIRegisterFormDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	ForceDelete(ids []uint, currentUserID uint) ([]uint, error)

	UpdateValue(kpiRegister *models.KPIRegister, tx utility.DatabaseTransaction) error
	GetTrendAnalysis(kpiID uint, limit int) ([]*kpiDTOs.KPITrendAnalysisDTO, error)
	GetLatestValueUpdateByKPIRegisterID(kpiRegisterID uint) (*kpiDTOs.KPIValueUpdateDetailsDTO, error)

	ImportKPIRegisters(input []map[string]interface{}, currentUserID uint) ([]uint, error)
	ValidateImportKPIRegisters(input []map[string]interface{}) ([]map[string]interface{}, error)
	ValidateKPIRegisterTemplate(file io.Reader) ([]map[string]interface{}, error)
	GenerateKPIRegisterTemplate() (*os.File, error)
	ExportKPIRegister(data []*kpiDTOs.KPIRegisterListDTO, currentUserID uint) (*os.File, error)

	GetKPICategorySummary(params *kpiDTOs.KPIDashboardQueryParamsDTO, currentUserID uint) ([]*dtos.KPICategorySummaryDTO, error)
	GetKPIStatusSummary(params *kpiDTOs.KPIDashboardQueryParamsDTO, currentUserID uint) ([]*dtos.KPIStatusSummaryDTO, error)
	GetKPIOverallSummary(params *kpiDTOs.KPIDashboardQueryParamsDTO, currentUserID uint) (*dtos.KPIOverallSummaryDTO, error)
	GetKPIByReviewStatus(params *kpiDTOs.KPIDashboardQueryParamsDTO, currentUserID uint) (*kpiDTOs.KPIReviewStatusDistributionResponseDTO, error)
	ActivateKPIRegisters(ids []uint, currentUserID uint) error
	ArchiveKPIRegisters(ids []uint, currentUserID uint) error
	GetAllResponsibleUsers(params *userDTOs.UserListQueryParamsDTO) (*userDTOs.UserListResponseDTO, error)
}

// Service and implementations
type KPIRegisterService struct {
	repo                   kpiRegisterRepositories.KPIRegisterRepositoryInterface
	activityLogSvc         logServices.ActivityLogServiceInterface
	userSvc                organizationServices.UserServiceInterface
	userSettingSvc         organizationServices.UserSettingServiceInterface
	organizationUnitSvc    settingServices.OrganizationUnitServiceInterface
	kpiStatusSvc           KPIStatusServiceInterface
	kpiReviewStatusSvc     KPIReviewStatusServiceInterface
	entryTypeSvc           EntryTypeServiceInterface
	dataSourceSvc          DataSourceServiceInterface
	kpiCategorySvc         KPICategoryServiceInterface
	calculationTypeSvc     masters.CalculationTypeServiceInterface
	unitsvc                masters.UnitServiceInterface
	reviewFrequencySvc     masters.ReviewFrequencyServiceInterface
	lifeCycleStatusSvc     masters.LifecycleStatusServiceInterface
	kpiAchievementLogicSvc KPIAchievementLogicServiceInterface
	notificationSvc        commonServices.NotificationServiceInterface
}

func NewKPIRegisterService(
	repo kpiRegisterRepositories.KPIRegisterRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	userSettingSvc organizationServices.UserSettingServiceInterface,
	organizationUnitSvc settingServices.OrganizationUnitServiceInterface,
	kpiStatusSvc KPIStatusServiceInterface,
	kpiReviewStatusSvc KPIReviewStatusServiceInterface,
	entryTypeSvc EntryTypeServiceInterface,
	dataSourceSvc DataSourceServiceInterface,
	kpiCategorySvc KPICategoryServiceInterface,
	calculationTypeSvc masters.CalculationTypeServiceInterface,
	unitsvc masters.UnitServiceInterface,
	reviewFrequencySvc masters.ReviewFrequencyServiceInterface,
	lifeCycleStatusSvc masters.LifecycleStatusServiceInterface,
	kpiAchievementLogicSvc KPIAchievementLogicServiceInterface,
	notificationSvc commonServices.NotificationServiceInterface,
) KPIRegisterServiceInterface {
	return &KPIRegisterService{
		repo:                   repo,
		activityLogSvc:         activityLogSvc,
		userSvc:                userSvc,
		userSettingSvc:         userSettingSvc,
		organizationUnitSvc:    organizationUnitSvc,
		kpiStatusSvc:           kpiStatusSvc,
		kpiReviewStatusSvc:     kpiReviewStatusSvc,
		entryTypeSvc:           entryTypeSvc,
		dataSourceSvc:          dataSourceSvc,
		kpiCategorySvc:         kpiCategorySvc,
		calculationTypeSvc:     calculationTypeSvc,
		unitsvc:                unitsvc,
		reviewFrequencySvc:     reviewFrequencySvc,
		lifeCycleStatusSvc:     lifeCycleStatusSvc,
		kpiAchievementLogicSvc: kpiAchievementLogicSvc,
		notificationSvc:        notificationSvc,
	}
}

func (krs *KPIRegisterService) GetAll(params *kpiDTOs.KPIRegisterListQueryParamsDTO) (*kpiDTOs.KPIRegisterListResponseDTO, error) {
	return krs.repo.GetAll(params)
}

func (krs KPIRegisterService) GetByID(id uint) (*models.KPIRegister, error) {
	kpiRegister, err := krs.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return kpiRegister, nil
}

func (krs KPIRegisterService) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*kpiDTOs.KPIRegisterDetailDTO, error) {
	// Initialize your DTO
	var kpiRegisterDTO kpiDTOs.KPIRegisterDetailDTO

	kpiRegister, err := krs.repo.GetFullDetailsByID(id, tx...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	// Map the result to your DTO
	err = mapstructure.Decode(kpiRegister, &kpiRegisterDTO)
	if err != nil {
		return nil, err
	}

	if kpiRegister.StartDate != nil {
		startDate := kpiRegister.StartDate.Format("2006-01-02")
		kpiRegisterDTO.StartDateString = &startDate
	}
	if kpiRegister.EndDate != nil {
		endDate := kpiRegister.EndDate.Format("2006-01-02")
		kpiRegisterDTO.EndDateString = &endDate
	}
	if kpiRegister.NextReviewDate != nil {
		nextReviewDate := kpiRegister.NextReviewDate.Format("2006-01-02")
		kpiRegisterDTO.NextReviewDateString = &nextReviewDate
	}
	if kpiRegister.LastReviewedDate != nil {
		lastReviewedDate := kpiRegister.LastReviewedDate.Format("2006-01-02")
		kpiRegisterDTO.LastReviewDateString = &lastReviewedDate
	}

	if kpiRegister.TargetDate != nil {
		targetDate := kpiRegister.TargetDate.Format("2006-01-02")
		kpiRegisterDTO.TargetDateString = &targetDate
	}

	kpiRegisterDTO.CreatedAt = kpiRegister.CreatedAt
	kpiRegisterDTO.UpdatedAt = kpiRegister.UpdatedAt

	return &kpiRegisterDTO, nil
}

func (krs *KPIRegisterService) Create(currentUserID uint, kpiRegisterFormDTO *kpiDTOs.KPIRegisterFormDTO) (uint, error) {
	tx := krs.repo.BeginTransaction()

	dates, err := krs.parseKPIRegisterDates(kpiRegisterFormDTO)
	if err != nil {
		return 0, err
	}
	entryType, err := krs.entryTypeSvc.GetByType("manual")
	if err != nil {
		return 0, err
	}

	dataSource, err := krs.dataSourceSvc.GetByType("user_input")
	if err != nil {
		return 0, err
	}
	kpiStatus, err := krs.kpiStatusSvc.GetByType("new")
	if err != nil {
		return 0, err
	}
	kpiReviewStatus, err := krs.kpiReviewStatusSvc.GetByType("awaiting-update")
	if err != nil {
		return 0, err
	}

	lifeCycleStatus, err := krs.lifeCycleStatusSvc.GetByType("active")
	if err != nil {
		return 0, err
	}

	referenceNumber, err := krs.repo.GenerateReferenceNumber()
	if err != nil {
		return 0, err
	}

	var formula *string

	if kpiRegisterFormDTO.Formula != nil {
		lowered := strings.ToUpper(strings.TrimSpace(*kpiRegisterFormDTO.Formula))
		formula = &lowered
	}

	kpiRegister := &models.KPIRegister{
		ReferenceNumber:       &referenceNumber,
		Title:                 kpiRegisterFormDTO.Title,
		IsKRI:                 kpiRegisterFormDTO.IsKRI,
		Description:           kpiRegisterFormDTO.Description,
		KPICategoryID:         kpiRegisterFormDTO.KPICategoryID,
		ResponsibleUserID:     kpiRegisterFormDTO.ResponsibleUserID,
		ReviewFrequencyID:     kpiRegisterFormDTO.ReviewFrequencyID,
		CalculationTypeID:     kpiRegisterFormDTO.CalculationTypeID,
		Formula:               formula,
		UnitID:                kpiRegisterFormDTO.UnitID,
		EntryTypeID:           &entryType.ID,
		DataSourceID:          &dataSource.ID,
		TargetValue:           kpiRegisterFormDTO.TargetValue,
		StartDate:             dates.StartDate,
		EndDate:               dates.EndDate,
		Remarks:               kpiRegisterFormDTO.Remarks,
		KPIStatusID:           kpiStatus.ID,
		KPIReviewStatusID:     &kpiReviewStatus.ID,
		KPIAchievementLogicID: kpiRegisterFormDTO.KPIAchievementLogicID,
		StatusID:              &lifeCycleStatus.ID,
		ReviewUserID:          kpiRegisterFormDTO.ReviewUserID,
		CreatedBYID:           currentUserID,
	}
	var periodStartDate *time.Time
	var periodEndDate *time.Time
	var periodLabel *string
	today := time.Now()

	if kpiRegisterFormDTO.ReviewFrequencyID != nil {
		frequency, err := krs.reviewFrequencySvc.GetByID(*kpiRegisterFormDTO.ReviewFrequencyID)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
		// endDate := krs.calculateNextReviewDate(today, frequency.Type)
		// endDate := utility.CalculateNextReviewDate(frequency.Type)

		periodStartDate, periodEndDate = krs.calculatePeriodDates(today, frequency.Type)
		periodLabel = krs.formatPeriodLabel(periodStartDate, periodEndDate, frequency.Type)

		if dates.StartDate != nil {
			periodStartDate = dates.StartDate
		}
		if dates.EndDate != nil {
			periodEndDate = dates.EndDate
		}
		// Assign EndDate (current period end)
		if kpiRegister.StartDate == nil {
			kpiRegister.StartDate = periodStartDate
		}

		if kpiRegister.EndDate == nil {
			kpiRegister.EndDate = periodEndDate
		}
		kpiRegister.LastReviewedDate = nil

		// Next review = next day after period end
		nextReviewDate := periodEndDate.AddDate(0, 0, 1)
		kpiRegister.NextReviewDate = &nextReviewDate

	}
	id, err := krs.repo.Create(kpiRegister, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	// Prepare input data records
	var inputs []*models.KPIRegisterDataInput
	for _, input := range kpiRegisterFormDTO.DataInputs {
		if input != nil {
			inputData := &models.KPIRegisterDataInput{
				KPIRegisterID: id,
				Input:         input.Input,
				Variable:      input.Variable,
				CreatedBYID:   currentUserID,
			}
			inputs = append(inputs, inputData)
		}
	}

	if len(inputs) > 0 {
		if err := krs.repo.CreateKPIRegisterInpuData(inputs, tx); err != nil {
			return 0, err
		}
	}
	// kpiValueUpdateFormDTO.Value = calculatedValue
	kpiValueUpdate := &models.KPIValueUpdate{
		KPIRegisterID:     id,
		ResponsibleUserID: kpiRegisterFormDTO.ResponsibleUserID,
		ReviewUserID:      kpiRegisterFormDTO.ReviewUserID,
		UnitID:            *kpiRegisterFormDTO.UnitID,
		ReviewedDate:      nil,
		// Value:         &calculatedValue,
		// Remarks:       kpiValueUpdateFormDTO.Remarks,
		IsLatest:          true, // Always set to true (1) when creating
		StartDate:         periodStartDate,
		EndDate:           periodEndDate,
		PeriodLabel:       periodLabel,
		KPIReviewStatusID: &kpiReviewStatus.ID,
		CreatedBYID:       currentUserID,
	}

	// Set all previous value updates for this KPIRegister to IsLatest = false
	if err := krs.repo.SetAllIsLatestFalse(id, tx); err != nil {
		tx.Rollback()
		return 0, err
	}

	updateid, err := krs.repo.CreateKPIValueUpdate(kpiValueUpdate, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	var valueUpdateDataInputs []*models.KPIValueUpdateDataInput
	for _, input := range kpiRegisterFormDTO.DataInputs {
		if input != nil {
			inputData := &models.KPIValueUpdateDataInput{
				KPIValueUpdateID: updateid,
				Input:            input.Input,
				Variable:         input.Variable,
				Value:            nil,
				CreatedBYID:      currentUserID,
			}
			valueUpdateDataInputs = append(valueUpdateDataInputs, inputData)

		}
	}

	if len(valueUpdateDataInputs) > 0 {
		if err := krs.repo.CreateKPIValueUpdateInputData(valueUpdateDataInputs, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
	}
	tx.Commit()
	go func() {
		existingKPIRegister, err := krs.GetFullDetailsByID(id)
		if err != nil {
			utility.LogError("Failed to fetch kpi register", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}
		err = krs.activityLogSvc.LogActivity(
			models.CreateAction,
			models.KPIRegister{}.TableName(),
			id,
			existingKPIRegister,
			currentUserID,
			"strategy", "kpi", "kpi_register", nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.KPIRegister{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}

		if existingKPIRegister.ResponsibleUser == nil {
			return
		}

		responsibleUser := existingKPIRegister.ResponsibleUser

		reviewPeriod := "N/A"
		if existingKPIRegister.StartDateString != nil && existingKPIRegister.EndDateString != nil {
			reviewPeriod = fmt.Sprintf("%s to %s", *existingKPIRegister.StartDateString, *existingKPIRegister.EndDateString)
		}

		urlParam := fmt.Sprintf(`{"kpi_register_id":%d}`, existingKPIRegister.ID)
		queryParam := ""
		entity := "kpi-register"
		action := "kpi-register-created"

		dbNotificationDTO := &commonDTOs.DBNotificationDTO{
			UserIDs:              []uint{uint(responsibleUser.ID)},
			ModuleIdentifier:     "kpi_register",
			NotificationType:     string(commonDTOs.AssignmentType),
			HeadingEn:            "KPI Review Assigned",
			HeadingAr:            "تم تعيين مراجعة KPI",
			TitleEn:              existingKPIRegister.Title,
			TitleAr:              existingKPIRegister.Title,
			Entity:               &entity,
			Action:               &action,
			URLParam:             urlParam,
			QueryParam:           queryParam,
			IsSystemNotification: false,
			CreatedByID:          &currentUserID,
		}

		emailTo := []string{responsibleUser.Email}

		subject := fmt.Sprintf(
			"KPI Review Assigned – %s %s",
			utility.GetStringValue(existingKPIRegister.ReferenceNumber),
			existingKPIRegister.Title,
		)

		appURL := utility.GetAppURL()

		targetDate := "N/A"
		if existingKPIRegister.EndDateString != nil {
			targetDate = *existingKPIRegister.EndDateString
		}

		emailData := map[string]interface{}{
			"FirstNameEn":     responsibleUser.FirstNameEN,
			"LastNameEn":      responsibleUser.LastNameEN,
			"ReferenceNumber": utility.GetStringValue(existingKPIRegister.ReferenceNumber),
			"Title":           existingKPIRegister.Title,
			"ReviewPeriod":    reviewPeriod,
			"StartDate":       existingKPIRegister.StartDateString,
			"EndDate":         existingKPIRegister.EndDateString,
			"TargetDate":      targetDate,
			"RedirectURL":     fmt.Sprintf("%s/key-performance-indicator/kpi-register/%d", appURL, existingKPIRegister.ID),
		}

		emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
			EmailTo:              emailTo,
			Subject:              subject,
			TemplateRelativePath: "key-performance-indicator/templates/email/kpi_responsible_user_mail.html",
			EmailData:            emailData,
		}

		err = krs.notificationSvc.SendNotification(
			dbNotificationDTO,
			emailNotificationDTO,
			"kpi_register",
		)
		if err != nil {
			utility.LogError("Failed to send KPI notification", map[string]interface{}{
				"error": err.Error(),
			})
		}

	}()
	return id, nil
}

func (krs *KPIRegisterService) Update(id uint, currentUserID uint, kpiRegisterFormDTO *kpiDTOs.KPIRegisterFormDTO) error {
	existingKPIRegister, err := krs.repo.GetByID(id)
	if err != nil {
		return utility.NewResourceNotFoundError(err, id)
	}

	oldKPIRegister, err := krs.GetFullDetailsByID(id)
	if err != nil {
		return err
	}

	latestValueUpdate, err := krs.repo.GetLatestValueUpdate(id)
	if err != nil {
		return err
	}

	dates, err := krs.parseKPIRegisterDates(kpiRegisterFormDTO)
	if err != nil {
		return err
	}

	// Determine what changed before mutating the model
	frequencyChanged := !uintPtrEqual(existingKPIRegister.ReviewFrequencyID, kpiRegisterFormDTO.ReviewFrequencyID)
	inputsChanged := dataInputsChanged(oldKPIRegister.DataInputs, kpiRegisterFormDTO.DataInputs)

	var formula *string
	if kpiRegisterFormDTO.Formula != nil {
		lowered := strings.ToUpper(strings.TrimSpace(*kpiRegisterFormDTO.Formula))
		formula = &lowered
	}

	existingKPIRegister.Title = kpiRegisterFormDTO.Title
	existingKPIRegister.IsKRI = kpiRegisterFormDTO.IsKRI
	existingKPIRegister.Description = kpiRegisterFormDTO.Description
	existingKPIRegister.KPICategoryID = kpiRegisterFormDTO.KPICategoryID
	existingKPIRegister.ResponsibleUserID = kpiRegisterFormDTO.ResponsibleUserID
	existingKPIRegister.ReviewFrequencyID = kpiRegisterFormDTO.ReviewFrequencyID
	existingKPIRegister.CalculationTypeID = kpiRegisterFormDTO.CalculationTypeID
	existingKPIRegister.Formula = formula
	existingKPIRegister.UnitID = kpiRegisterFormDTO.UnitID
	existingKPIRegister.TargetValue = kpiRegisterFormDTO.TargetValue
	existingKPIRegister.StartDate = dates.StartDate
	existingKPIRegister.EndDate = dates.EndDate
	existingKPIRegister.Remarks = kpiRegisterFormDTO.Remarks
	existingKPIRegister.KPIAchievementLogicID = kpiRegisterFormDTO.KPIAchievementLogicID
	existingKPIRegister.ReviewUserID = kpiRegisterFormDTO.ReviewUserID
	existingKPIRegister.UpdatedBYID = &currentUserID

	var periodStartDate *time.Time
	var periodEndDate *time.Time
	var periodLabel *string

	if frequencyChanged && kpiRegisterFormDTO.ReviewFrequencyID != nil {
		frequency, err := krs.reviewFrequencySvc.GetByID(*kpiRegisterFormDTO.ReviewFrequencyID)
		if err != nil {
			return err
		}
		today := time.Now()
		periodStartDate, periodEndDate = krs.calculatePeriodDates(today, frequency.Type)
		periodLabel = krs.formatPeriodLabel(periodStartDate, periodEndDate, frequency.Type)

		existingKPIRegister.StartDate = periodStartDate
		existingKPIRegister.EndDate = periodEndDate

		nextReviewDate := periodEndDate.AddDate(0, 0, 1)
		existingKPIRegister.NextReviewDate = &nextReviewDate
	} else if dates.EndDate != nil {
		nextReviewDate := dates.EndDate.AddDate(0, 0, 1)
		existingKPIRegister.NextReviewDate = &nextReviewDate
	}

	tx := krs.repo.BeginTransaction()

	if err = krs.repo.Update(existingKPIRegister, tx); err != nil {
		tx.Rollback()
		return err
	}

	if inputsChanged {
		if err := krs.repo.DeleteKPIRegisterInputData(id, tx); err != nil {
			tx.Rollback()
			return err
		}

		var inputs []*models.KPIRegisterDataInput
		for _, input := range kpiRegisterFormDTO.DataInputs {
			if input != nil {
				inputs = append(inputs, &models.KPIRegisterDataInput{
					KPIRegisterID: id,
					Input:         input.Input,
					Variable:      input.Variable,
					CreatedBYID:   currentUserID,
				})
			}
		}

		if len(inputs) > 0 {
			if err := krs.repo.CreateKPIRegisterInpuData(inputs, tx); err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	if latestValueUpdate != nil {
		valueUpdateChanged := false

		if kpiRegisterFormDTO.ResponsibleUserID != nil && (latestValueUpdate.ResponsibleUserID == nil || *latestValueUpdate.ResponsibleUserID != *kpiRegisterFormDTO.ResponsibleUserID) {
			latestValueUpdate.ResponsibleUserID = kpiRegisterFormDTO.ResponsibleUserID
			valueUpdateChanged = true
		}
		if kpiRegisterFormDTO.ReviewUserID != nil && (latestValueUpdate.ReviewUserID == nil || *latestValueUpdate.ReviewUserID != *kpiRegisterFormDTO.ReviewUserID) {
			latestValueUpdate.ReviewUserID = kpiRegisterFormDTO.ReviewUserID
			valueUpdateChanged = true
		}
		if kpiRegisterFormDTO.UnitID != nil && latestValueUpdate.UnitID != *kpiRegisterFormDTO.UnitID {
			latestValueUpdate.UnitID = *kpiRegisterFormDTO.UnitID
			valueUpdateChanged = true
		}
		if frequencyChanged {
			latestValueUpdate.StartDate = periodStartDate
			latestValueUpdate.EndDate = periodEndDate
			latestValueUpdate.PeriodLabel = periodLabel
			valueUpdateChanged = true
		}

		if valueUpdateChanged {
			latestValueUpdate.UpdatedBYID = &currentUserID
			if err := krs.repo.UpdateKPIValueUpdate(latestValueUpdate, tx); err != nil {
				tx.Rollback()
				return err
			}
		}

		if inputsChanged {
			if err := krs.repo.DeleteKPIValueUpdateInputData(latestValueUpdate.ID, tx); err != nil {
				tx.Rollback()
				return err
			}

			var valueUpdateDataInputs []*models.KPIValueUpdateDataInput
			for _, input := range kpiRegisterFormDTO.DataInputs {
				if input != nil {
					valueUpdateDataInputs = append(valueUpdateDataInputs, &models.KPIValueUpdateDataInput{
						KPIValueUpdateID: latestValueUpdate.ID,
						Input:            input.Input,
						Variable:         input.Variable,
						Value:            nil,
						CreatedBYID:      currentUserID,
					})
				}
			}

			if len(valueUpdateDataInputs) > 0 {
				if err := krs.repo.CreateKPIValueUpdateInputData(valueUpdateDataInputs, tx); err != nil {
					tx.Rollback()
					return err
				}
			}
		}
	}

	tx.Commit()

	go func() {
		updatedKPIRegister, err := krs.GetFullDetailsByID(id)
		if err != nil {
			utility.LogError("Failed to fetch kpi register after update", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		changedData := make(map[string]interface{})

		oldData := map[string]interface{}{
			"title":            oldKPIRegister.Title,
			"is_kri":           oldKPIRegister.IsKRI,
			"description":      oldKPIRegister.Description,
			"kpi_category":     oldKPIRegister.KPICategory.TitleEn,
			"responsible_user": utility.GetFieldValue(oldKPIRegister.ResponsibleUser),
			"review_user":      utility.GetFieldValue(oldKPIRegister.ReviewUser),
			"review_frequency": utility.GetFieldValue(oldKPIRegister.ReviewFrequency),
			"calculation_type": utility.GetFieldValue(oldKPIRegister.CalculationType),
			"formula":          oldKPIRegister.Formula,
			"unit":             utility.GetFieldValue(oldKPIRegister.Unit),
			"target_value":     oldKPIRegister.TargetValue,
			"start_date":       oldKPIRegister.StartDateString,
			"end_date":         oldKPIRegister.EndDateString,
			"remarks":          oldKPIRegister.Remarks,
		}

		newData := map[string]interface{}{
			"title":            updatedKPIRegister.Title,
			"is_kri":           updatedKPIRegister.IsKRI,
			"description":      updatedKPIRegister.Description,
			"kpi_category":     updatedKPIRegister.KPICategory.TitleEn,
			"responsible_user": utility.GetFieldValue(updatedKPIRegister.ResponsibleUser),
			"review_user":      utility.GetFieldValue(updatedKPIRegister.ReviewUser),
			"review_frequency": utility.GetFieldValue(updatedKPIRegister.ReviewFrequency),
			"calculation_type": utility.GetFieldValue(updatedKPIRegister.CalculationType),
			"formula":          updatedKPIRegister.Formula,
			"unit":             utility.GetFieldValue(updatedKPIRegister.Unit),
			"target_value":     updatedKPIRegister.TargetValue,
			"start_date":       updatedKPIRegister.StartDateString,
			"end_date":         updatedKPIRegister.EndDateString,
			"remarks":          updatedKPIRegister.Remarks,
		}

		fieldChanges := utility.CompareFieldChanges(oldData, newData)
		for k, v := range fieldChanges {
			changedData[k] = v
		}

		if len(changedData) > 0 {
			err := krs.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.KPIRegister{}.TableName(),
				id,
				changedData,
				currentUserID,
				"strategy", "kpi", "kpi_register", nil,
			)
			if err != nil {
				utility.LogError("Failed to log activity", map[string]interface{}{
					"action":    models.UpdateAction,
					"tableName": models.KPIRegister{}.TableName(),
					"recordID":  id,
					"userID":    currentUserID,
					"error":     err.Error(),
				})
			}
		}

		// Notify new responsible user if assignment changed
		responsibleUserChanged := oldKPIRegister.ResponsibleUserID != nil &&
			updatedKPIRegister.ResponsibleUserID != nil &&
			*oldKPIRegister.ResponsibleUserID != *updatedKPIRegister.ResponsibleUserID

		if responsibleUserChanged && updatedKPIRegister.ResponsibleUser != nil {
			responsibleUser := updatedKPIRegister.ResponsibleUser

			reviewPeriod := "N/A"
			if updatedKPIRegister.StartDateString != nil && updatedKPIRegister.EndDateString != nil {
				reviewPeriod = fmt.Sprintf("%s to %s", *updatedKPIRegister.StartDateString, *updatedKPIRegister.EndDateString)
			}

			urlParam := fmt.Sprintf(`{"kpi_register_id":%d}`, updatedKPIRegister.ID)
			queryParam := ""
			entity := "kpi-register"
			action := "kpi-register-updated"

			dbNotificationDTO := &commonDTOs.DBNotificationDTO{
				UserIDs:              []uint{uint(responsibleUser.ID)},
				ModuleIdentifier:     "kpi_register",
				NotificationType:     string(commonDTOs.AssignmentType),
				HeadingEn:            "KPI Review Assigned",
				HeadingAr:            "تم تعيين مراجعة KPI",
				TitleEn:              updatedKPIRegister.Title,
				TitleAr:              updatedKPIRegister.Title,
				Entity:               &entity,
				Action:               &action,
				URLParam:             urlParam,
				QueryParam:           queryParam,
				IsSystemNotification: false,
				CreatedByID:          &currentUserID,
			}

			targetDate := "N/A"
			if updatedKPIRegister.EndDateString != nil {
				targetDate = *updatedKPIRegister.EndDateString
			}

			appURL := utility.GetAppURL()

			emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
				EmailTo: []string{responsibleUser.Email},
				Subject: fmt.Sprintf(
					"KPI Review Assigned – %s %s",
					utility.GetStringValue(updatedKPIRegister.ReferenceNumber),
					updatedKPIRegister.Title,
				),
				TemplateRelativePath: "key-performance-indicator/templates/email/kpi_responsible_user_mail.html",
				EmailData: map[string]interface{}{
					"FirstNameEn":     responsibleUser.FirstNameEN,
					"LastNameEn":      responsibleUser.LastNameEN,
					"ReferenceNumber": utility.GetStringValue(updatedKPIRegister.ReferenceNumber),
					"Title":           updatedKPIRegister.Title,
					"StartDate":       updatedKPIRegister.StartDateString,
					"EndDate":         updatedKPIRegister.EndDateString,
					"ReviewPeriod":    reviewPeriod,
					"TargetDate":      targetDate,
					"RedirectURL":     fmt.Sprintf("%s/key-performance-indicator/kpi-register/%d", appURL, updatedKPIRegister.ID),
				},
			}

			if err := krs.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "kpi_register"); err != nil {
				utility.LogError("Failed to send KPI notification", map[string]interface{}{
					"error": err.Error(),
				})
			}
		}
	}()

	return nil
}

func (krs *KPIRegisterService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	tx := krs.repo.BeginTransaction()
	var failedIds []uint
	var successfulIds []uint
	var deletedKPIRegister []map[string]interface{}

	for _, id := range ids {
		kpiRegister, err := krs.GetFullDetailsByID(id, tx)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		deletedKPIRegister = append(deletedKPIRegister, map[string]interface{}{
			"id":    id,
			"title": kpiRegister.Title,
		})

		// Delete the business
		if err := krs.repo.Delete(id, tx); err != nil {
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

	tx.Commit()

	// Log bulk delete activities for successful deletions
	if len(successfulIds) > 0 {
		go func() {
			err := krs.activityLogSvc.LogBulkActivities(
				models.DeleteAction, // Make sure this action constant is defined appropriately
				models.KPIRegister{}.TableName(),
				successfulIds,
				deletedKPIRegister,
				currentUserID,
				nil,
			)
			if err != nil {
				// @todo - save to log file
				fmt.Printf("Failed to log bulk delete activity: %v\n", err)
			}
		}()
	}

	return failedIds, nil
}

func (krs *KPIRegisterService) parseKPIRegisterDates(dto *kpiDTOs.KPIRegisterFormDTO) (*kpiDTOs.KPIRegisterParsedDates, error) {
	parsed := &kpiDTOs.KPIRegisterParsedDates{}

	var err error

	if dto.StartDate != nil {
		parsed.StartDate, err = utility.ParseDate(*dto.StartDate)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid start date"), "invalid_date", "start date")
		}
	}

	if dto.EndDate != nil {
		parsed.EndDate, err = utility.ParseDate(*dto.EndDate)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid end date"), "invalid_date", "end date")
		}
	}

	return parsed, nil
}

func (krs *KPIRegisterService) UpdateValue(kpiRegister *models.KPIRegister, tx utility.DatabaseTransaction) error {
	err := krs.repo.Update(kpiRegister, tx)
	if err != nil {
		return err
	}

	return nil
}

func (krs *KPIRegisterService) GetTrendAnalysis(kpiID uint, limit int) ([]*kpiDTOs.KPITrendAnalysisDTO, error) {
	kpiValueUpdates, err := krs.repo.GetTrendAnalysis(kpiID, limit)
	if err != nil {
		return nil, err
	}

	if len(kpiValueUpdates) == 0 {
		return []*kpiDTOs.KPITrendAnalysisDTO{}, nil
	}

	// Already sorted and limited in repo (DESC), so just map to DTOs
	var trend []*kpiDTOs.KPITrendAnalysisDTO
	for _, valueUpdate := range kpiValueUpdates {
		var kpiUpdate kpiDTOs.KPITrendAnalysisDTO
		err := mapstructure.Decode(valueUpdate, &kpiUpdate)
		if err != nil {
			return nil, err
		}
		kpiUpdate.Value = float64(*valueUpdate.Value)
		reviewedAt := valueUpdate.ReviewedDate.Format("2006-01-02")
		kpiUpdate.ReviewedAtString = reviewedAt
		trend = append(trend, &kpiUpdate)
	}

	return trend, nil
}

func (krs *KPIRegisterService) GetLatestValueUpdateByKPIRegisterID(kpiRegisterID uint) (*kpiDTOs.KPIValueUpdateDetailsDTO, error) {
	valueUpdate, err := krs.repo.GetLatestValueUpdate(kpiRegisterID)
	if err != nil {
		return nil, err
	}

	// If no value update is found
	if valueUpdate == nil {
		return nil, nil
	}

	var dto kpiDTOs.KPIValueUpdateDetailsDTO

	// ✅ Fix: check ReviewedDate
	if valueUpdate.ReviewedDate != nil {
		reviewedDate := valueUpdate.ReviewedDate.Format("2006-01-02")
		dto.ReviewedDateString = &reviewedDate
	}

	if valueUpdate.StartDate != nil {
		startDate := valueUpdate.StartDate.Format("2006-01-02")
		dto.StartDateString = &startDate
	}

	if valueUpdate.EndDate != nil {
		endDate := valueUpdate.EndDate.Format("2006-01-02")
		dto.EndDateString = &endDate
	}

	err = mapstructure.Decode(valueUpdate, &dto)
	if err != nil {
		return nil, err
	}

	return &dto, nil
}

func (krs *KPIRegisterService) GenerateKPIRegisterTemplate() (*os.File, error) {

	// Fetch dropdown master data

	kpiCategories, err := krs.kpiCategorySvc.GetAll("", "", "", "active", "")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch KPI categories: %v", err)
	}

	users, err := krs.userSvc.GetAllUsers()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch users: %v", err)
	}

	reviewFrequencies, err := krs.reviewFrequencySvc.GetAll("", "", "", "active", "")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch review frequencies: %v", err)
	}

	calculationTypes, err := krs.calculationTypeSvc.GetAll("", "", "", "active")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch calculation types: %v", err)
	}

	units, err := krs.unitsvc.GetAll("", "", "", "active")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch units: %v", err)
	}

	kpiAchievementLogics, err := krs.kpiAchievementLogicSvc.GetAll("", "", "", "active")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch kpi achievement logics: %v", err)
	}

	// Template Sheet (SINGLE SHEET)

	templateHeaders := []string{
		"Title",
		"Type",
		"Objective",
		"Domain",
		"Responsible User",
		"Review User",
		"Review Frequency",
		"Calculation Type",
		"Formula",
		"Unit",
		"Data Input A",
		"Data Input B",
		"Target Value",
		"Achievement Type",
		"Remarks",
	}

	templateData := [][]interface{}{}

	// Dropdown Sheets

	typeSheet := [][]interface{}{
		{"Type"},
		{"KPI"},
		{"KRI"},
	}

	categorySheet := [][]interface{}{{"Domain"}}
	for _, c := range kpiCategories {
		categorySheet = append(categorySheet, []interface{}{c.TitleEn})
	}

	userSheet := [][]interface{}{{"Responsible User"}}
	for _, u := range users {
		if strings.TrimSpace(u.PrimaryUnit) == "" {
			continue
		}
		name := strings.TrimSpace(u.FirstNameEN + " " + u.LastNameEN)
		if name != "" {
			userSheet = append(userSheet, []interface{}{name})
		}
	}

	reviewUserSheet := [][]interface{}{{"Review User"}}
	for _, u := range users {
		if strings.TrimSpace(u.PrimaryUnit) == "" {
			continue
		}
		name := strings.TrimSpace(u.FirstNameEN + " " + u.LastNameEN)
		if name != "" {
			reviewUserSheet = append(reviewUserSheet, []interface{}{name})
		}
	}

	reviewFrequencySheet := [][]interface{}{{"Review Frequency"}}
	for _, rf := range reviewFrequencies {
		reviewFrequencySheet = append(reviewFrequencySheet, []interface{}{rf.TitleEn})
	}

	calculationTypeSheet := [][]interface{}{{"Calculation Type"}}
	for _, ct := range calculationTypes {
		calculationTypeSheet = append(calculationTypeSheet, []interface{}{ct.TitleEn})
	}

	unitSheet := [][]interface{}{{"Unit"}}
	for _, u := range units {
		unitSheet = append(unitSheet, []interface{}{u.TitleEn})
	}

	achievementTypeSheet := [][]interface{}{{"Achievement Type"}}
	for _, a := range kpiAchievementLogics {
		if a.Type == "higher" || a.Type == "lower" {
			achievementTypeSheet = append(achievementTypeSheet, []interface{}{a.TitleEn})
		}
	}
	if len(achievementTypeSheet) == 1 {
		achievementTypeSheet = append(achievementTypeSheet,
			[]interface{}{"Higher is Better"},
			[]interface{}{"Lower is Better"},
		)
	}

	// ------------------------------------------------
	// Dropdown Configuration
	// ------------------------------------------------

	dropDownConfig := map[string]bool{
		"Type":             true,
		"Domain":           true,
		"Responsible User": true,
		"Review User":      true,
		"Review Frequency": true,
		"Calculation Type": true,
		"Unit":             true,
		"Achievement Type": true,
	}

	// ------------------------------------------------
	// Generate Excel (ONE FILE)
	// ------------------------------------------------

	file, err := utility.GenerateExcelFileWithSheetsAndDropDown(
		map[string][][]interface{}{
			"Template": append(
				[][]interface{}{utility.ConvertToInterfaceSlice(templateHeaders)},
				templateData...,
			),
			"Type":             typeSheet,
			"Domain":           categorySheet,
			"Responsible User": userSheet,
			"Review User":      reviewUserSheet,
			"Review Frequency": reviewFrequencySheet,
			"Calculation Type": calculationTypeSheet,
			"Unit":             unitSheet,
			"Achievement Type": achievementTypeSheet,
		},
		dropDownConfig,
	)

	if err != nil {
		return nil, err
	}

	return file, nil
}

func (krs *KPIRegisterService) ValidateKPIRegisterTemplate(file io.Reader) ([]map[string]interface{}, error) {

	rows, err := utility.ReadExcelFile(file)
	if err != nil {
		return nil, err
	}

	if len(rows) == 0 {
		return nil, fmt.Errorf("Template sheet is empty")
	}

	headersWithType := []string{
		"Title",
		"Type",
		"Objective",
		"Domain",
		"Responsible User",
		"Review User",
		"Review Frequency",
		"Calculation Type",
		"Formula",
		"Unit",
		"Data Input A",
		"Data Input B",
		"Target Value",
		"Achievement Type",
		"Remarks",
	}

	snakeHeadersWithType := []string{
		"title",
		"type",
		"objective",
		"kpi_category",
		"responsible_user",
		"review_user",
		"review_frequency",
		"calculation_type",
		"formula",
		"unit",
		"a",
		"b",
		"target_value",
		"kpi_achievement_type",
		"remarks",
	}

	headersWithoutType := []string{
		"Title",
		"Objective",
		"Domain",
		"Responsible User",
		"Review User",
		"Review Frequency",
		"Calculation Type",
		"Formula",
		"Unit",
		"Data Input A",
		"Data Input B",
		"Target Value",
		"Achievement Type",
		"Remarks",
	}

	snakeHeadersWithoutType := []string{
		"title",
		"objective",
		"kpi_category",
		"responsible_user",
		"review_user",
		"review_frequency",
		"calculation_type",
		"formula",
		"unit",
		"a",
		"b",
		"target_value",
		"kpi_achievement_type",
		"remarks",
	}

	excelHeaders := make([]string, len(rows[0]))
	for i, h := range rows[0] {
		excelHeaders[i] = strings.TrimSpace(h)
	}

	var activeSnakeHeaders []string

	if len(excelHeaders) == len(headersWithType) {
		headersValid := utility.ValidateHeaders(headersWithType, excelHeaders)
		if !headersValid {
			temp := make([]string, len(headersWithType))
			copy(temp, headersWithType)
			if strings.EqualFold(excelHeaders[1], "Indicator Type") || strings.EqualFold(excelHeaders[1], "Type") {
				temp[1] = excelHeaders[1]
			}
			if strings.EqualFold(excelHeaders[13], "Type") || strings.EqualFold(excelHeaders[13], "Achievement Type") {
				temp[13] = excelHeaders[13]
			}
			headersValid = utility.ValidateHeaders(temp, excelHeaders)
		}
		if headersValid {
			activeSnakeHeaders = snakeHeadersWithType
		}
	} else if len(excelHeaders) == len(headersWithoutType) {
		headersValid := utility.ValidateHeaders(headersWithoutType, excelHeaders)
		if !headersValid {
			temp := make([]string, len(headersWithoutType))
			copy(temp, headersWithoutType)
			if strings.EqualFold(excelHeaders[12], "Type") || strings.EqualFold(excelHeaders[12], "Achievement Type") {
				temp[12] = excelHeaders[12]
			}
			headersValid = utility.ValidateHeaders(temp, excelHeaders)
		}
		if headersValid {
			activeSnakeHeaders = snakeHeadersWithoutType
		}
	}

	if activeSnakeHeaders == nil {
		return nil, utility.NewValidationError(
			fmt.Errorf("invalid headers"),
			"invalid_header",
			"KPI Register Template",
		)
	}

	// ------------------------------------------------
	// Parse rows + build DataInputs from A–B
	// ------------------------------------------------

	var input []map[string]interface{}

	for _, row := range rows[1:] {

		record := make(map[string]interface{})

		// map base fields
		for i, key := range activeSnakeHeaders {
			if i < len(row) {
				record[key] = row[i]
			} else {
				record[key] = ""
			}
		}

		if obj, ok := record["objective"]; ok && obj != nil {
			record["description"] = obj
		}

		// --------------------------------------------
		// Build DataInputs from A, B
		// --------------------------------------------

		var dataInputs []*kpiDTOs.InputDataFormDTO

		for _, v := range []string{"a", "b"} {
			val := strings.TrimSpace(fmt.Sprint(record[v]))
			if val != "" && val != "<nil>" {
				dataInputs = append(dataInputs, &kpiDTOs.InputDataFormDTO{
					Variable: strings.ToUpper(v),
					Input:    val,
				})
			}
		}

		record["data_inputs"] = dataInputs

		input = append(input, record)
	}

	// ------------------------------------------------
	// Single source of truth
	// ------------------------------------------------

	return krs.validateKPIRegisters(input)
}

func (krs *KPIRegisterService) ValidateImportKPIRegisters(input []map[string]interface{}) ([]map[string]interface{}, error) {

	// Single source of truth for validation
	return krs.validateKPIRegisters(input)
}

func (krs *KPIRegisterService) validateKPIRegisters(input []map[string]interface{}) ([]map[string]interface{}, error) {

	kpiCategories, _ := krs.kpiCategorySvc.GetAll("", "", "", "active", "")
	reviewFrequencies, _ := krs.reviewFrequencySvc.GetAll("", "", "", "active", "")
	calculationTypes, _ := krs.calculationTypeSvc.GetAll("", "", "", "active")
	units, _ := krs.unitsvc.GetAll("", "", "", "active")
	users, _ := krs.userSvc.GetAllUsers()
	kpiAchievementLogics, _ := krs.kpiAchievementLogicSvc.GetAll("", "", "", "active")

	categoryMap := map[string]uint{}
	for _, c := range kpiCategories {
		categoryMap[normalize(c.TitleEn)] = c.ID
	}

	reviewFrequencyMap := map[string]uint{}
	for _, rf := range reviewFrequencies {
		reviewFrequencyMap[normalize(rf.TitleEn)] = rf.ID
	}

	calculationTypeMap := map[string]uint{}
	calculationTypeTypeMap := map[string]string{}
	for _, ct := range calculationTypes {
		calculationTypeMap[normalize(ct.TitleEn)] = ct.ID
		calculationTypeTypeMap[normalize(ct.TitleEn)] = ct.Type
	}

	unitMap := map[string]uint{}
	for _, u := range units {
		unitMap[normalize(u.TitleEn)] = u.ID
	}

	userMap := map[string]uint{}
	for _, u := range users {
		if strings.TrimSpace(u.PrimaryUnit) == "" {
			continue
		}
		name := strings.TrimSpace(u.FirstNameEN + " " + u.LastNameEN)
		if name != "" {
			userMap[normalize(name)] = uint(u.ID)
		}
	}

	validAchievementLogicIDs := map[uint]*models.KPIAchievementLogic{}
	achievementLogicTypeMap := map[string]uint{}
	for _, a := range kpiAchievementLogics {
		if a.Type == "higher" || a.Type == "lower" {
			validAchievementLogicIDs[a.ID] = a
			achievementLogicTypeMap[normalize(a.TitleEn)] = uint(a.ID)
			achievementLogicTypeMap[normalize(a.TitleAr)] = uint(a.ID)
			achievementLogicTypeMap[normalize(a.Type)] = uint(a.ID)
		}
	}

	allowedVars := map[string]bool{
		"A": true,
		"B": true,
		"C": true,
		"D": true,
	}

	var validatedData []map[string]interface{}

	// Validate rows

	for index, record := range input {

		rowErrors := make(map[string]interface{})
		out := make(map[string]interface{})

		for k, v := range record {
			out[k] = v
		}

		// ---------- Title ----------
		title := strings.TrimSpace(fmt.Sprint(record["title"]))
		if title == "" {
			rowErrors["title"] = "Title is required"
		} else {
			out["title"] = title
		}

		// ---------- Type (KPI / KRI) ----------
		isKRI := false
		if rawIsKRI, ok := record["is_kri"]; ok && rawIsKRI != nil {
			switch v := rawIsKRI.(type) {
			case bool:
				isKRI = v
			case string:
				s := strings.ToLower(strings.TrimSpace(v))
				if s == "true" || s == "1" || s == "kri" {
					isKRI = true
				}
			case int:
				if v == 1 {
					isKRI = true
				}
			case float64:
				if v == 1 {
					isKRI = true
				}
			}
		}

		if rawType, ok := record["type"]; ok && rawType != nil {
			s := strings.ToUpper(strings.TrimSpace(fmt.Sprint(rawType)))
			if s == "KRI" {
				isKRI = true
			} else if s == "KPI" {
				isKRI = false
			} else if s != "" && s != "<NIL>" && s != "MANUAL" && s != "AUTOMATED" && s != "HIGHER" && s != "LOWER" && s != "HIGHER IS BETTER" && s != "LOWER IS BETTER" {
				rowErrors["type"] = "Type must be KPI or KRI"
			}
		} else if rawIndType, ok := record["indicator_type"]; ok && rawIndType != nil {
			s := strings.ToUpper(strings.TrimSpace(fmt.Sprint(rawIndType)))
			if s == "KRI" {
				isKRI = true
			} else if s == "KPI" {
				isKRI = false
			}
		}

		out["is_kri"] = isKRI
		if isKRI {
			out["type"] = "KRI"
			out["indicator_type"] = "KRI"
		} else {
			out["type"] = "KPI"
			out["indicator_type"] = "KPI"
		}

		// ---------- Objective / Description (optional) ----------
		var descStr string
		if rawDesc, ok := record["description"]; ok && rawDesc != nil {
			s := strings.TrimSpace(fmt.Sprint(rawDesc))
			if s != "" && s != "<nil>" {
				descStr = s
			}
		}
		if descStr == "" {
			if rawObj, ok := record["objective"]; ok && rawObj != nil {
				s := strings.TrimSpace(fmt.Sprint(rawObj))
				if s != "" && s != "<nil>" {
					descStr = s
				}
			}
		}
		out["description"] = descStr
		out["objective"] = descStr

		// ---------- KPI Category ----------
		if key := normalize(fmt.Sprint(record["kpi_category"])); key == "" {
			rowErrors["kpi_category"] = "Domain is required"
		} else if id, ok := categoryMap[key]; !ok {
			rowErrors["kpi_category"] = "Invalid KPI Category"
		} else {
			out["kpi_category_id"] = id
		}

		// ---------- Responsible User ----------
		if key := normalize(fmt.Sprint(record["responsible_user"])); key == "" {
			rowErrors["responsible_user"] = "Responsible User is required"
		} else if id, ok := userMap[key]; !ok {
			rowErrors["responsible_user"] = "Invalid Responsible User"
		} else {
			out["responsible_user_id"] = id
		}

		// ---------- Review User (optional) ----------
		if rawRevUser, ok := record["review_user"]; ok && rawRevUser != nil {
			revUserStr := strings.TrimSpace(fmt.Sprint(rawRevUser))
			if revUserStr != "" && revUserStr != "<nil>" {
				key := normalize(revUserStr)
				if id, ok := userMap[key]; ok {
					out["review_user_id"] = id
					out["review_user"] = revUserStr
				} else {
					rowErrors["review_user"] = "Invalid Review User"
				}
			} else {
				out["review_user"] = ""
			}
		} else if rawRevID, ok := record["review_user_id"]; ok && rawRevID != nil {
			var rID uint
			switch v := rawRevID.(type) {
			case uint:
				rID = v
			case int:
				rID = uint(v)
			case float64:
				rID = uint(v)
			case string:
				if parsed, err := strconv.ParseUint(strings.TrimSpace(v), 10, 32); err == nil {
					rID = uint(parsed)
				}
			}
			if rID > 0 {
				found := false
				for _, uid := range userMap {
					if uid == rID {
						found = true
						break
					}
				}
				if found {
					out["review_user_id"] = rID
				} else {
					rowErrors["review_user"] = "Invalid Review User"
				}
			}
		}

		// ---------- Review Frequency ----------
		if key := normalize(fmt.Sprint(record["review_frequency"])); key == "" {
			rowErrors["review_frequency"] = "Review Frequency is required"
		} else if id, ok := reviewFrequencyMap[key]; !ok {
			rowErrors["review_frequency"] = "Invalid Review Frequency"
		} else {
			out["review_frequency_id"] = id
		}

		// ---------- Calculation Type ----------
		isDirect := false
		if key := normalize(fmt.Sprint(record["calculation_type"])); key == "" {
			rowErrors["calculation_type"] = "Calculation Type is required"
		} else if id, ok := calculationTypeMap[key]; !ok {
			rowErrors["calculation_type"] = "Invalid Calculation Type"
		} else {
			out["calculation_type_id"] = id
			isDirect = calculationTypeTypeMap[key] == "direct"
		}

		// ---------- Unit ----------
		if key := normalize(fmt.Sprint(record["unit"])); key == "" {
			rowErrors["unit"] = "Unit is required"
		} else if id, ok := unitMap[key]; !ok {
			rowErrors["unit"] = "Invalid Unit"
		} else {
			out["unit_id"] = id
		}

		// ---------- Achievement Type ----------
		var achKey string
		for _, k := range []string{"kpi_achievement_type", "achievement_type"} {
			if rawAch, ok := record[k]; ok && rawAch != nil {
				s := strings.TrimSpace(fmt.Sprint(rawAch))
				if s != "" && s != "<nil>" {
					achKey = normalize(s)
					break
				}
			}
		}
		if achKey == "" {
			if rawType, ok := record["type"]; ok && rawType != nil {
				s := normalize(fmt.Sprint(rawType))
				if s == "higher" || s == "lower" || s == "higher is better" || s == "lower is better" {
					achKey = s
				}
			}
		}

		if achKey != "" {
			if id, ok := achievementLogicTypeMap[achKey]; !ok {
				rowErrors["kpi_achievement_type"] = "Invalid Achievement Type"
			} else {
				out["kpi_achievement_logic_id"] = id
				if a, ok := validAchievementLogicIDs[id]; ok {
					out["kpi_achievement_type"] = a.TitleEn
					out["achievement_type"] = a.TitleEn
					out["kpi_achievement_logic_type"] = a.Type
				}
			}
		} else if rawID, ok := record["kpi_achievement_logic_id"]; ok && rawID != nil {
			var logicID uint
			switch v := rawID.(type) {
			case float64:
				logicID = uint(v)
			case int:
				logicID = uint(v)
			case uint:
				logicID = v
			case string:
				if parsed, err := strconv.ParseUint(strings.TrimSpace(v), 10, 32); err == nil {
					logicID = uint(parsed)
				}
			}

			if logicID > 0 {
				if a, ok := validAchievementLogicIDs[logicID]; !ok {
					rowErrors["kpi_achievement_type"] = "Invalid Achievement Type"
				} else {
					out["kpi_achievement_logic_id"] = a.ID
					out["kpi_achievement_type"] = a.TitleEn
					out["achievement_type"] = a.TitleEn
					out["kpi_achievement_logic_type"] = a.Type
				}
			} else {
				rowErrors["kpi_achievement_type"] = "Achievement Type is required"
			}
		} else {
			rowErrors["kpi_achievement_type"] = "Achievement Type is required"
		}

		// ---------- Target Value ----------
		var targetValStr string
		if rawTarget, ok := record["target_value"]; ok && rawTarget != nil {
			s := strings.TrimSpace(fmt.Sprint(rawTarget))
			if s != "" && s != "<nil>" {
				targetValStr = s
			}
		}

		if targetValStr == "" {
			rowErrors["target_value"] = "Target Value is required"
		} else {
			if f, err := strconv.ParseFloat(targetValStr, 64); err != nil {
				rowErrors["target_value"] = "Target Value must be a number"
			} else {
				out["target_value"] = f
			}
		}

		// ---------- Formula & Data Inputs (formula type only) ----------
		var dataInputs []*kpiDTOs.InputDataFormDTO

		if !isDirect {
			rawFormula := strings.TrimSpace(fmt.Sprint(record["formula"]))
			if rawFormula == "" || rawFormula == "<nil>" {
				rowErrors["formula"] = "Formula is required"
			} else {
				formula := strings.ToUpper(rawFormula)
				out["formula"] = formula

				expr, err := govaluate.NewEvaluableExpression(formula)
				if err != nil {
					rowErrors["formula"] = "Invalid formula expression"
				} else {
					vars := expr.Vars()
					if len(vars) == 0 {
						rowErrors["formula"] = "Formula must contain at least one variable (A or B)"
					} else {
						valueMap := map[string]string{}

						for _, v := range []string{"a", "b"} {
							val := strings.TrimSpace(fmt.Sprint(record[v]))
							if val != "" && val != "<nil>" {
								varName := strings.ToUpper(v)
								valueMap[varName] = val

								dataInputs = append(dataInputs, &kpiDTOs.InputDataFormDTO{
									Variable: varName,
									Input:    val,
								})
							}
						}

						for _, v := range vars {
							vUpper := strings.ToUpper(v)
							if !allowedVars[vUpper] {
								rowErrors["formula"] = "Formula can use only A, B"
							} else if val, exists := valueMap[vUpper]; !exists || val == "" {
								if vUpper == "A" {
									rowErrors["a"] = "Data Input A is required"
								} else if vUpper == "B" {
									rowErrors["b"] = "Data Input B is required"
								} else {
									rowErrors["formula"] = fmt.Sprintf("Data Input %s is required", vUpper)
								}
							}
						}
					}
				}
			}
		} else {
			delete(out, "formula")
			dataInputs = nil
		}

		out["data_inputs"] = dataInputs

		if len(rowErrors) > 0 {
			out["row_number"] = index + 2
			out["errors"] = rowErrors
		}

		validatedData = append(validatedData, out)
	}

	return validatedData, nil
}

func normalize(s string) string {
	return strings.ToLower(strings.Join(strings.Fields(strings.TrimSpace(s)), " "))
}

func uintPtrEqual(a, b *uint) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return *a == *b
}

func dataInputsChanged(old []*kpiDTOs.KPIRegisterDataInputDTO, incoming []*kpiDTOs.InputDataFormDTO) bool {
	if len(old) != len(incoming) {
		return true
	}
	oldMap := make(map[string]string, len(old))
	for _, o := range old {
		if o != nil {
			oldMap[o.Variable] = o.Input
		}
	}
	for _, n := range incoming {
		if n == nil {
			continue
		}
		if oldMap[n.Variable] != n.Input {
			return true
		}
	}
	return false
}

func (krs *KPIRegisterService) ImportKPIRegisters(input []map[string]interface{}, currentUserID uint) ([]uint, error) {

	var createdIDs []uint

	// helper functions
	strPtr := func(s string) *string {
		s = strings.TrimSpace(s)
		if s == "" {
			return nil
		}
		return &s
	}

	uintPtr := func(u uint) *uint {
		return &u
	}

	kpiAchievementLogics, _ := krs.kpiAchievementLogicSvc.GetAll("", "", "", "active")
	validAchievementLogicIDs := map[uint]*models.KPIAchievementLogic{}
	achievementLogicTypeMap := map[string]uint{}
	for _, a := range kpiAchievementLogics {
		if a.Type == "higher" || a.Type == "lower" {
			validAchievementLogicIDs[a.ID] = a
			achievementLogicTypeMap[normalize(a.TitleEn)] = uint(a.ID)
			achievementLogicTypeMap[normalize(a.TitleAr)] = uint(a.ID)
			achievementLogicTypeMap[normalize(a.Type)] = uint(a.ID)
		}
	}

	users, _ := krs.userSvc.GetAllUsers()
	userMap := map[string]uint{}
	for _, u := range users {
		if strings.TrimSpace(u.PrimaryUnit) == "" {
			continue
		}
		name := strings.TrimSpace(u.FirstNameEN + " " + u.LastNameEN)
		if name != "" {
			userMap[normalize(name)] = uint(u.ID)
		}
	}

	for _, record := range input {

		var kpiCategoryID uint
		switch v := record["kpi_category_id"].(type) {
		case uint:
			kpiCategoryID = v
		case int:
			kpiCategoryID = uint(v)
		case float64:
			kpiCategoryID = uint(v)
		default:
			return nil, fmt.Errorf("invalid kpi_category_id in record: %+v", record)
		}

		var responsibleUserID uint
		switch v := record["responsible_user_id"].(type) {
		case uint:
			responsibleUserID = v
		case int:
			responsibleUserID = uint(v)
		case float64:
			responsibleUserID = uint(v)
		default:
			return nil, fmt.Errorf("invalid responsible_user_id in record: %+v", record)
		}

		var reviewFrequencyID uint
		switch v := record["review_frequency_id"].(type) {
		case uint:
			reviewFrequencyID = v
		case int:
			reviewFrequencyID = uint(v)
		case float64:
			reviewFrequencyID = uint(v)
		default:
			return nil, fmt.Errorf("invalid review_frequency_id in record: %+v", record)
		}

		var calculationTypeID uint
		switch v := record["calculation_type_id"].(type) {
		case uint:
			calculationTypeID = v
		case int:
			calculationTypeID = uint(v)
		case float64:
			calculationTypeID = uint(v)
		default:
			return nil, fmt.Errorf("invalid calculation_type_id in record: %+v", record)
		}

		var unitID uint
		switch v := record["unit_id"].(type) {
		case uint:
			unitID = v
		case int:
			unitID = uint(v)
		case float64:
			unitID = uint(v)
		default:
			return nil, fmt.Errorf("invalid unit_id in record: %+v", record)
		}

		calcType, err := krs.calculationTypeSvc.GetByID(calculationTypeID)
		if err != nil {
			return nil, fmt.Errorf("failed to get calculation type: %v", err)
		}
		isDirect := calcType.Type == "direct"

		var formula *string
		if !isDirect {
			if v, ok := record["formula"].(string); ok {
				if upper := strings.ToUpper(strings.TrimSpace(v)); upper != "" && upper != "<NIL>" {
					formula = &upper
				}
			}
			if formula == nil {
				return nil, fmt.Errorf("formula is required for formula based calculation in record: %+v", record)
			}
		}

		var targetValue *float64
		if v, ok := record["target_value"]; ok && v != nil {
			switch tv := v.(type) {
			case float64:
				targetValue = &tv
			case int:
				f := float64(tv)
				targetValue = &f
			case string:
				s := strings.TrimSpace(tv)
				if s != "" && s != "<nil>" {
					f, err := strconv.ParseFloat(s, 64)
					if err != nil {
						return nil, fmt.Errorf("invalid target_value in record: %+v", record)
					}
					targetValue = &f
				}
			}
		}

		if targetValue == nil {
			return nil, fmt.Errorf("target_value is required in record: %+v", record)
		}

		var kpiAchievementLogicID uint
		if v, ok := record["kpi_achievement_logic_id"]; ok && v != nil {
			switch val := v.(type) {
			case uint:
				kpiAchievementLogicID = val
			case int:
				kpiAchievementLogicID = uint(val)
			case float64:
				kpiAchievementLogicID = uint(val)
			case string:
				if parsed, err := strconv.ParseUint(strings.TrimSpace(val), 10, 32); err == nil {
					kpiAchievementLogicID = uint(parsed)
				}
			}
		}

		if kpiAchievementLogicID == 0 {
			for _, k := range []string{"kpi_achievement_type", "achievement_type", "type"} {
				if achTypeStr, ok := record[k].(string); ok {
					if id, ok := achievementLogicTypeMap[normalize(achTypeStr)]; ok {
						kpiAchievementLogicID = id
						break
					}
				}
			}
		}

		if kpiAchievementLogicID == 0 {
			return nil, fmt.Errorf("kpi_achievement_logic_id is required in record: %+v", record)
		}

		if _, ok := validAchievementLogicIDs[kpiAchievementLogicID]; !ok {
			return nil, fmt.Errorf("invalid kpi_achievement_logic_id in record: %+v", record)
		}

		// ---------- Start Date (optional) ----------
		var startDate *string
		if v, ok := record["start_date"].(string); ok {
			if s := strings.TrimSpace(v); s != "" {
				startDate = &s
			}
		}

		// ---------- Data Inputs (formula type only) ----------
		var dataInputs []*kpiDTOs.InputDataFormDTO

		if !isDirect {
			if record["data_inputs"] != nil {
				if rawInputs, ok := record["data_inputs"].([]interface{}); ok {
					for _, item := range rawInputs {
						if m, ok := item.(map[string]interface{}); ok {
							variable := strings.ToUpper(strings.TrimSpace(fmt.Sprint(m["variable"])))
							inputVal := strings.TrimSpace(fmt.Sprint(m["input"]))
							if variable != "" && inputVal != "" && inputVal != "<nil>" {
								dataInputs = append(dataInputs, &kpiDTOs.InputDataFormDTO{
									Variable: variable,
									Input:    inputVal,
								})
							}
						}
					}
				}
			}

			if len(dataInputs) == 0 {
				for _, v := range []string{"a", "b"} {
					val := strings.TrimSpace(fmt.Sprint(record[v]))
					if val != "" && val != "<nil>" {
						dataInputs = append(dataInputs, &kpiDTOs.InputDataFormDTO{
							Variable: strings.ToUpper(v),
							Input:    val,
						})
					}
				}
			}

			if formula != nil {
				expr, err := govaluate.NewEvaluableExpression(*formula)
				if err == nil {
					inputVarMap := map[string]bool{}
					for _, inp := range dataInputs {
						inputVarMap[inp.Variable] = true
					}
					for _, v := range expr.Vars() {
						vUpper := strings.ToUpper(v)
						if !inputVarMap[vUpper] {
							return nil, fmt.Errorf("data input for variable %s is required in record: %+v", vUpper, record)
						}
					}
				}
			}
		}

		// ---------- Type (KPI / KRI) ----------
		isKRI := false
		if rawIsKRI, ok := record["is_kri"]; ok && rawIsKRI != nil {
			switch v := rawIsKRI.(type) {
			case bool:
				isKRI = v
			case string:
				s := strings.ToLower(strings.TrimSpace(v))
				if s == "true" || s == "1" || s == "kri" {
					isKRI = true
				}
			case int:
				if v == 1 {
					isKRI = true
				}
			case float64:
				if v == 1 {
					isKRI = true
				}
			}
		}

		if rawType, ok := record["type"]; ok && rawType != nil {
			s := strings.ToUpper(strings.TrimSpace(fmt.Sprint(rawType)))
			if s == "KRI" {
				isKRI = true
			} else if s == "KPI" {
				isKRI = false
			}
		} else if rawIndType, ok := record["indicator_type"]; ok && rawIndType != nil {
			s := strings.ToUpper(strings.TrimSpace(fmt.Sprint(rawIndType)))
			if s == "KRI" {
				isKRI = true
			} else if s == "KPI" {
				isKRI = false
			}
		}

		// ---------- Review User (optional) ----------
		var reviewUserID *uint
		if v, ok := record["review_user_id"]; ok && v != nil {
			switch r := v.(type) {
			case uint:
				if r > 0 {
					reviewUserID = &r
				}
			case int:
				if r > 0 {
					u := uint(r)
					reviewUserID = &u
				}
			case float64:
				if r > 0 {
					u := uint(r)
					reviewUserID = &u
				}
			case string:
				if parsed, err := strconv.ParseUint(strings.TrimSpace(r), 10, 32); err == nil && parsed > 0 {
					u := uint(parsed)
					reviewUserID = &u
				}
			}
		}

		if reviewUserID == nil {
			if revStr, ok := record["review_user"].(string); ok {
				revStr = strings.TrimSpace(revStr)
				if revStr != "" && revStr != "<nil>" {
					if id, ok := userMap[normalize(revStr)]; ok {
						reviewUserID = &id
					}
				}
			}
		}

		// ---------- Objective / Description (optional) ----------
		var desc *string
		if v, ok := record["description"]; ok && v != nil {
			s := strings.TrimSpace(fmt.Sprint(v))
			if s != "" && s != "<nil>" {
				desc = &s
			}
		}
		if desc == nil {
			if v, ok := record["objective"]; ok && v != nil {
				s := strings.TrimSpace(fmt.Sprint(v))
				if s != "" && s != "<nil>" {
					desc = &s
				}
			}
		}

		// ---------- Build DTO ----------
		dto := &kpiDTOs.KPIRegisterFormDTO{
			Title:                 strings.TrimSpace(fmt.Sprint(record["title"])),
			IsKRI:                 isKRI,
			Description:           desc,
			KPICategoryID:         kpiCategoryID,
			ResponsibleUserID:     uintPtr(responsibleUserID),
			ReviewUserID:          reviewUserID,
			ReviewFrequencyID:     uintPtr(reviewFrequencyID),
			CalculationTypeID:     uintPtr(calculationTypeID),
			Formula:               formula,
			UnitID:                uintPtr(unitID),
			TargetValue:           targetValue,
			StartDate:             startDate,
			Remarks:               strPtr(fmt.Sprint(record["remarks"])),
			DataInputs:            dataInputs, // ✅ NOW CORRECT
			KPIAchievementLogicID: &kpiAchievementLogicID,
		}

		id, err := krs.Create(currentUserID, dto)
		if err != nil {
			return nil, fmt.Errorf(
				"failed to create KPI Register for record %+v: %v",
				record,
				err,
			)
		}

		createdIDs = append(createdIDs, id)
	}

	return createdIDs, nil
}

func (krs *KPIRegisterService) ExportKPIRegister(data []*kpiDTOs.KPIRegisterListDTO, currentUserID uint) (*os.File, error) {

	userSetting, err := krs.userSettingSvc.GetUserSettingsByUserID(currentUserID)
	if err != nil {
		return nil, err
	}

	var headers []string

	if userSetting.Language == "ar" {
		headers = []string{
			"SL No.",
			"رقم المرجع",
			"العنوان",
			"النوع",
			"النطاق",
			"المسؤول",
			"الحالة",
			"تكرار المراجعة",
			"نوع الحساب",
			"نوع الإنجاز",
			"القيمة المستهدفة",
			"الوحدة",
			"تاريخ البداية",
			"تاريخ النهاية",
			"نسبة الإنجاز",
			"تاريخ الإنشاء",
		}
	} else {
		headers = []string{
			"SL No.",
			"Reference Number",
			"Title",
			"Type",
			"Domain",
			"Responsible User",
			"Status",
			"Review Frequency",
			"Calculation Type",
			"Achievement Type",
			"Target Value",
			"Unit",
			"Start Date",
			"End Date",
			"Achievement %",
			"Created By",
		}
	}

	rows := make([][]interface{}, len(data))

	for i, item := range data {

		var (
			category        string
			status          string
			reviewFrequency string
			calculationType string
			achievementType string
			typeStr         string
			unit            string
			responsibleUser string
			createdBy       string
		)

		if userSetting.Language == "ar" {
			category = item.KPICategoryTitleAr
			status = item.KPIStatusTitleAr
			reviewFrequency = utility.GetStringValue(item.ReviewFrequencyTitleAr)
			calculationType = utility.GetStringValue(item.CalculationTypeTitleAr)
			achievementType = item.KPIAchievementLogicTitleAr
			if achievementType == "" {
				achievementType = item.KPIAchievementLogicTitleEn
			}
			if item.IsKRI {
				typeStr = "مؤشر المخاطر الرئيسي (KRI)"
			} else {
				typeStr = "مؤشر الأداء الرئيسي (KPI)"
			}
			unit = utility.GetStringValue(item.UnitTitleAr)

			responsibleUser =
				utility.GetStringValue(item.ResponsibleUserFirstNameAr) + " " +
					utility.GetStringValue(item.ResponsibleUserLastNameAr)

			createdBy =
				item.CreatedUserFirstNameAr + " " +
					item.CreatedUserLastNameAr

		} else {
			category = item.KPICategoryTitleEn
			status = item.KPIStatusTitleEn
			reviewFrequency = utility.GetStringValue(item.ReviewFrequencyTitleEn)
			calculationType = utility.GetStringValue(item.CalculationTypeTitleEn)
			achievementType = item.KPIAchievementLogicTitleEn
			if item.IsKRI {
				typeStr = "KRI"
			} else {
				typeStr = "KPI"
			}
			unit = utility.GetStringValue(item.UnitTitleEn)

			responsibleUser =
				utility.GetStringValue(item.ResponsibleUserFirstNameEn) + " " +
					utility.GetStringValue(item.ResponsibleUserLastNameEn)

			createdBy =
				item.CreatedUserFirstNameEn + " " +
					item.CreatedUserLastNameEn
		}

		rows[i] = []interface{}{
			i + 1,
			utility.GetStringValue(item.ReferenceNumber),
			item.Title,
			typeStr,
			category,
			responsibleUser,
			status,
			reviewFrequency,
			calculationType,
			achievementType,
			floatPointerToString(item.TargetValue),
			unit,
			utility.GetStringValue(item.StartDate),
			utility.GetStringValue(item.EndDate),
			floatPointerToString(item.Achievement),
			createdBy,
		}
	}

	return utility.GenerateExcelFile(headers, rows)
}

func (krs *KPIRegisterService) calculateNextReviewDate(now time.Time, frequencyType string) time.Time {
	year := now.Year()
	location := now.Location()

	var reviewDates []time.Time

	switch frequencyType {

	case "monthly":
		for month := 1; month <= 12; month++ {
			reviewDates = append(reviewDates,
				time.Date(year, time.Month(month), 1, 0, 0, 0, 0, location))
		}

	case "quarterly":
		for _, month := range []int{1, 4, 7, 10} {
			reviewDates = append(reviewDates,
				time.Date(year, time.Month(month), 1, 0, 0, 0, 0, location))
		}

	case "half-yearly":
		for _, month := range []int{1, 7} {
			reviewDates = append(reviewDates,
				time.Date(year, time.Month(month), 1, 0, 0, 0, 0, location))
		}

	case "annually":
		reviewDates = append(reviewDates,
			time.Date(year, time.January, 1, 0, 0, 0, 0, location))

	case "weekly":
		return now.AddDate(0, 0, 7)

	case "biweekly":
		return now.AddDate(0, 0, 14)

	case "daily":
		return now.AddDate(0, 0, 1)

	default:
		return now
	}

	// Find next upcoming review date
	for _, date := range reviewDates {
		if !date.Before(now) {
			return date
		}
	}

	// If no date left this year → first slot next year
	switch frequencyType {
	case "monthly":
		return time.Date(year+1, time.January, 1, 0, 0, 0, 0, location)
	case "quarterly":
		return time.Date(year+1, time.January, 1, 0, 0, 0, 0, location)
	case "half-yearly":
		return time.Date(year+1, time.January, 1, 0, 0, 0, 0, location)
	case "annually":
		return time.Date(year+1, time.January, 1, 0, 0, 0, 0, location)
	}

	return now
}

func parseExcelDate(value string) (time.Time, error) {
	value = strings.TrimSpace(value)

	layouts := []string{
		"2006-01-02", // yyyy-mm-dd
		"01-02-06",   // mm-dd-yy
		"01-02-2006", // mm-dd-yyyy
		"02-01-2006", // dd-mm-yyyy
		"02/01/2006", // dd/mm/yyyy
		"01/02/2006", // mm/dd/yyyy
		"2/1/2006",
		"1/2/2006",
	}

	for _, layout := range layouts {
		if t, err := time.Parse(layout, value); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("invalid date format")
}

func floatPointerToString(val *float64) string {
	if val == nil {
		return ""
	}
	return fmt.Sprintf("%.2f", *val)
}

func (krs KPIRegisterService) GetKPICategorySummary(params *kpiDTOs.KPIDashboardQueryParamsDTO, currentUserID uint) ([]*dtos.KPICategorySummaryDTO, error) {

	return krs.repo.GetKPICategorySummary(params, currentUserID)
}

func (krs KPIRegisterService) GetKPIStatusSummary(params *kpiDTOs.KPIDashboardQueryParamsDTO, currentUserID uint) ([]*dtos.KPIStatusSummaryDTO, error) {

	return krs.repo.GetKPIStatusSummary(params, currentUserID)
}

func (krs KPIRegisterService) GetKPIOverallSummary(params *kpiDTOs.KPIDashboardQueryParamsDTO, currentUserID uint) (*dtos.KPIOverallSummaryDTO, error) {
	return krs.repo.GetKPIOverallSummary(params, currentUserID)
}

func (krs *KPIRegisterService) GetKPIByReviewStatus(params *kpiDTOs.KPIDashboardQueryParamsDTO, currentUserID uint) (*kpiDTOs.KPIReviewStatusDistributionResponseDTO, error) {

	results, err := krs.repo.GetKPIByReviewStatus(params, currentUserID)
	if err != nil {
		return nil, err
	}

	var total int64
	for _, r := range results {
		total += r.Count
	}

	return &kpiDTOs.KPIReviewStatusDistributionResponseDTO{
		Total: total,
		Items: results,
	}, nil
}

func (krs *KPIRegisterService) calculatePeriodDates(baseDate time.Time, frequencyType string) (*time.Time, *time.Time) {
	location := baseDate.Location()
	year, month, day := baseDate.Date()

	var start time.Time
	var end time.Time

	switch frequencyType {

	case "daily":
		start = time.Date(year, month, day, 0, 0, 0, 0, location)
		end = start

	case "weekly":
		// Week starts Monday
		offset := (int(baseDate.Weekday()) + 6) % 7
		start = baseDate.AddDate(0, 0, -offset)
		start = time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, location)
		end = start.AddDate(0, 0, 6)

	case "biweekly":
		// Start from 1st of month, 2-week blocks
		firstDayOfMonth := time.Date(year, month, 1, 0, 0, 0, 0, location)
		daysFromStart := int(baseDate.Sub(firstDayOfMonth).Hours() / 24)
		blockStart := (daysFromStart / 14) * 14

		start = firstDayOfMonth.AddDate(0, 0, blockStart)
		end = start.AddDate(0, 0, 13)

		// cap to month end
		monthEnd := time.Date(year, month+1, 0, 0, 0, 0, 0, location)
		if end.After(monthEnd) {
			end = monthEnd
		}

	case "monthly":
		start = time.Date(year, month, 1, 0, 0, 0, 0, location)
		end = time.Date(year, month+1, 0, 0, 0, 0, 0, location)

	case "quarterly":
		var startMonth time.Month
		switch {
		case month <= 3:
			startMonth = 1
		case month <= 6:
			startMonth = 4
		case month <= 9:
			startMonth = 7
		default:
			startMonth = 10
		}
		start = time.Date(year, startMonth, 1, 0, 0, 0, 0, location)
		end = time.Date(year, startMonth+3, 0, 0, 0, 0, 0, location)

	case "half-yearly":
		if month <= 6 {
			start = time.Date(year, 1, 1, 0, 0, 0, 0, location)
			end = time.Date(year, 7, 0, 0, 0, 0, 0, location) // June end
		} else {
			start = time.Date(year, 7, 1, 0, 0, 0, 0, location)
			end = time.Date(year, 13, 0, 0, 0, 0, 0, location) // Dec end
		}

	case "annually":
		start = time.Date(year, 1, 1, 0, 0, 0, 0, location)
		end = time.Date(year, 12, 31, 0, 0, 0, 0, location)

	default:
		start = baseDate
		end = baseDate
	}

	return &start, &end
}

func (krs *KPIRegisterService) formatPeriodLabel(startDate, endDate *time.Time, frequencyType string) *string {
	if startDate == nil {
		return nil
	}

	switch frequencyType {

	case "daily":
		val := startDate.Format("02 Jan 2006")
		return &val

	case "weekly", "biweekly":
		if endDate != nil {
			val := fmt.Sprintf("%s - %s",
				startDate.Format("02 Jan"),
				endDate.Format("02 Jan 2006"))
			return &val
		}
		val := startDate.Format("02 Jan 2006")
		return &val

	case "monthly":
		val := startDate.Format("Jan 2006")
		return &val

	case "quarterly":
		quarter := ((int(startDate.Month()) - 1) / 3) + 1
		val := fmt.Sprintf("Q%d %d", quarter, startDate.Year())
		return &val

	case "half-yearly":
		half := "H1"
		if startDate.Month() > 6 {
			half = "H2"
		}
		val := fmt.Sprintf("%s %d", half, startDate.Year())
		return &val

	case "annually":
		val := fmt.Sprintf("%d", startDate.Year())
		return &val

	default:
		val := startDate.Format("02 Jan 2006")
		return &val
	}
}

func (krs *KPIRegisterService) ActivateKPIRegisters(ids []uint, currentUserID uint) error {
	activeStatus, err := krs.lifeCycleStatusSvc.GetByType("active")
	if err != nil {
		return err
	}
	archiveStatus, err := krs.lifeCycleStatusSvc.GetByType("archived")
	if err != nil {
		return err
	}

	successfulIds, err := krs.repo.ActivateKPIRegisters(ids, archiveStatus.ID, activeStatus.ID)

	if len(successfulIds) > 0 {
		go func() {
			err := krs.activityLogSvc.LogBulkActivities(
				models.ActivateAction,
				models.KPIRegister{}.TableName(),
				successfulIds,
				nil,
				currentUserID,
				nil,
			)
			if err != nil {
				fmt.Printf("Failed to log bulk delete activity: %v\n", err)
			}
		}()
	}

	return err
}

func (krs *KPIRegisterService) ArchiveKPIRegisters(ids []uint, currentUserID uint) error {
	activeStatus, err := krs.lifeCycleStatusSvc.GetByType("active")
	if err != nil {
		return err
	}
	archiveStatus, err := krs.lifeCycleStatusSvc.GetByType("archived")
	if err != nil {
		return err
	}

	successfulIds, err := krs.repo.ArchiveKPIRegisters(ids, archiveStatus.ID, activeStatus.ID)

	if len(successfulIds) > 0 {
		go func() {
			err := krs.activityLogSvc.LogBulkActivities(
				models.ArchiveAction,
				models.KPIRegister{}.TableName(),
				successfulIds,
				nil,
				currentUserID,
				nil,
			)
			if err != nil {
				fmt.Printf("Failed to log bulk delete activity: %v\n", err)
			}
		}()
	}

	return err
}

func (krs *KPIRegisterService) ForceDelete(ids []uint, currentUserID uint) ([]uint, error) {

	tx := krs.repo.BeginTransaction()

	var failedIds []uint
	var successIds []uint
	var deletedData []map[string]interface{}

	for _, id := range ids {

		kpiRegister, err := krs.repo.GetFullDetailsByID(id)
		if err != nil {
			failedIds = append(failedIds, id)
			continue
		}

		successIds = append(successIds, id)

		deletedData = append(deletedData, map[string]interface{}{
			"id":       id,
			"title_en": kpiRegister.Title,
			"title_ar": kpiRegister.Title,
		})
	}

	if len(failedIds) > 0 {
		tx.Rollback()
		return failedIds, nil
	}

	if err := krs.repo.ForceDelete(successIds, tx); err != nil {
		tx.Rollback()
		return successIds, err
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return successIds, err
	}

	go func() {
		err := krs.activityLogSvc.LogBulkActivities(
			models.DeleteAction,
			models.KPIRegister{}.TableName(),
			successIds,
			deletedData,
			currentUserID,
			nil,
		)
		if err != nil {
			fmt.Printf("Failed to log bulk delete activity: %v\n", err)
		}
	}()

	return nil, nil
}

func (krs *KPIRegisterService) GetAllResponsibleUsers(params *userDTOs.UserListQueryParamsDTO) (*userDTOs.UserListResponseDTO, error) {
	return krs.repo.GetAllResponsibleUsers(params)
}
