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
	settingsServices "jethur/grc/settings/services"
	"jethur/grc/tprm/dtos"
	tprmRepo "jethur/grc/tprm/repositories"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type ThirdPartyRegisterServiceInterface interface {
	GetAll(params *dtos.ThirdPartyRegisterListQueryParamsDTO) (*dtos.ThirdPartyRegisterListResponseDTO, error)
	GetByID(id uint) (*models.ThirdPartyRegister, error)
	GetFullDetailsByID(id uint) (*dtos.ThirdPartyRegisterDetailDTO, error)
	Create(currentUserID uint, dto *dtos.ThirdPartyRegisterFormDTO) (uint, error)
	Update(id uint, currentUserID uint, dto *dtos.ThirdPartyRegisterFormDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	ForceDelete(ids []uint, currentUserID uint) ([]uint, error)
	ExportThirdPartyRegister(data []*dtos.ThirdPartyRegisterListDTO, currentUserID uint) (*os.File, error)
}

type ThirdPartyRegisterService struct {
	repo              tprmRepo.ThirdPartyRegisterRepositoryInterface
	activityLogSvc    logServices.ActivityLogServiceInterface
	severityLevelSvc  masterServices.SeverityLevelServiceInterface
	userSvc           organizationServices.UserServiceInterface
	orgUnitSvc        settingsServices.OrganizationUnitServiceInterface
	thirdPartyTypeSvc ThirdPartyTypeServiceInterface
	countrySvc        masterServices.CountryServiceInterface
	userSettingSvc    organizationServices.UserSettingServiceInterface
	notificationSvc   commonServices.NotificationServiceInterface
}

func NewThirdPartyRegisterService(
	repo tprmRepo.ThirdPartyRegisterRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	severityLevelSvc masterServices.SeverityLevelServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	orgUnitSvc settingsServices.OrganizationUnitServiceInterface,
	thirdPartyTypeSvc ThirdPartyTypeServiceInterface,
	countrySvc masterServices.CountryServiceInterface,
	userSettingSvc organizationServices.UserSettingServiceInterface,
	notificationSvc commonServices.NotificationServiceInterface,
) *ThirdPartyRegisterService {
	return &ThirdPartyRegisterService{
		repo:              repo,
		activityLogSvc:    activityLogSvc,
		severityLevelSvc:  severityLevelSvc,
		userSvc:           userSvc,
		orgUnitSvc:        orgUnitSvc,
		thirdPartyTypeSvc: thirdPartyTypeSvc,
		countrySvc:        countrySvc,
		userSettingSvc:    userSettingSvc,
		notificationSvc:   notificationSvc,
	}
}

func (tprs *ThirdPartyRegisterService) GetAll(params *dtos.ThirdPartyRegisterListQueryParamsDTO) (*dtos.ThirdPartyRegisterListResponseDTO, error) {
	return tprs.repo.GetAll(params)
}

func (tprs *ThirdPartyRegisterService) GetByID(id uint) (*models.ThirdPartyRegister, error) {
	tpr, err := tprs.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return tpr, nil
}

func (tprs *ThirdPartyRegisterService) GetFullDetailsByID(id uint) (*dtos.ThirdPartyRegisterDetailDTO, error) {
	var dto dtos.ThirdPartyRegisterDetailDTO

	tpr, err := tprs.repo.GetFullDetailsByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	err = mapstructure.Decode(tpr, &dto)
	if err != nil {
		return nil, err
	}

	if tpr.OnboardingDate != nil {
		formatted := tpr.OnboardingDate.Format("2006-01-02")
		dto.OnboardingDateString = &formatted
	}
	if tpr.ServiceStartDate != nil {
		formatted := tpr.ServiceStartDate.Format("2006-01-02")
		dto.ServiceStartDateString = &formatted
	}
	if tpr.ServiceEndDate != nil {
		formatted := tpr.ServiceEndDate.Format("2006-01-02")
		dto.ServiceEndDateString = &formatted
	}

	dto.CreatedAt = tpr.CreatedAt

	return &dto, nil
}

func (tprs *ThirdPartyRegisterService) Create(currentUserID uint, form *dtos.ThirdPartyRegisterFormDTO) (uint, error) {
	tx := tprs.repo.BeginTransaction()

	responsibleUnitID, err := tprs.userSvc.GetUserPrimaryUnitID(*form.ResponsibleUserID)
	if err != nil {
		return 0, utility.NewValidationError(fmt.Errorf("organization unit not set for responsible user"), "not_exist", "bcp communication plan responsible unit")
	}

	if err := tprs.validateForm(form); err != nil {
		return 0, err
	}

	dates, err := tprs.parseDates(form)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	referenceNumber, err := tprs.repo.GenerateReferenceNumber()
	if err != nil {
		return 0, err
	}

	tpr := &models.ThirdPartyRegister{
		ReferenceNumber:    referenceNumber,
		Title:              form.Title,
		ThirdPartyTypeID:   form.ThirdPartyTypeID,
		CountryID:          form.CountryID,
		Address:            form.Address,
		Website:            form.Website,
		OrgUnitID:          form.OrganizationUnitID,
		RepresentativeName: form.RepresentativeName,
		RepresentativeRole: form.RepresentativeRole,
		Email:              form.Email,
		PhoneNumber:        form.PhoneNumber,
		ResponsibleUserID:  form.ResponsibleUserID,
		ResponsibleUnitID:  &responsibleUnitID,
		OnboardingDate:     dates.OnboardingDate,
		ServiceStartDate:   dates.ServiceStartDate,
		ServiceEndDate:     dates.ServiceEndDate,
		SeverityLevelID:    form.SeverityLevelID,
		CreatedBYID:        currentUserID,
	}

	id, err := tprs.repo.Create(tpr, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	tx.Commit()
	go func() {
		urlParam := fmt.Sprintf(`{"id":%d}`, id)
		entity := "third-party-register"
		action := "responsible-user-assign"
		if form.ResponsibleUserID != nil {

			thirdParty, err := tprs.GetByID(id)
			if err != nil {
				return
			}

			dbNotificationDTO := &commonDTOs.DBNotificationDTO{
				UserIDs:              []uint{*form.ResponsibleUserID},
				ModuleIdentifier:     "tprm",
				NotificationType:     string(commonDTOs.AssignmentType),
				ReferenceNumber:      thirdParty.ReferenceNumber,
				HeadingEn:            "Third Party Register Assigned",
				HeadingAr:            "تم تعيين سجل طرف ثالث",
				TitleEn:              thirdParty.Title,
				TitleAr:              thirdParty.Title,
				Entity:               &entity,
				Action:               &action,
				URLParam:             urlParam,
				IsSystemNotification: false,
				CreatedByID:          &currentUserID,
			}

			responsibleUser, err := tprs.userSvc.GetByID(*thirdParty.ResponsibleUserID)
			if err != nil {
				utility.LogError("Failed to fetch responsible user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}

			assignedBy, err := tprs.userSvc.GetByID(currentUserID)
			if err != nil {
				utility.LogError("Failed to fetch assigning user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}

			emailTo := []string{responsibleUser.Email}
			subject := fmt.Sprintf("Third Party Register Responsiblity Assigned: %s", thirdParty.Title)
			appURL := utility.GetAppURL()
			emailData := map[string]interface{}{
				"FirstNameEn":           responsibleUser.FirstNameEn,
				"LastNameEn":            responsibleUser.LastNameEn,
				"Title":                 thirdParty.Title,
				"ReferenceNumber":       thirdParty.ReferenceNumber,
				"AssignedByFirstNameEn": assignedBy.FirstNameEn,
				"AssignedByLastNameEn":  assignedBy.LastNameEn,
				"AssignedOn":            time.Now().Format("2006-01-02"),
				"RedirectURL":           fmt.Sprintf("%s/tprm/third-party-register/%d/overview", appURL, thirdParty.ID),
			}

			emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
				EmailTo:              emailTo,
				Subject:              subject,
				TemplateRelativePath: "tprm/templates/email/third_party_register_responsible_user_mail.html",
				EmailData:            emailData,
			}

			err = tprs.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "third_party_register")
			if err != nil {
				utility.LogError("Failed to send notifications", map[string]interface{}{
					"error": err.Error(),
				})
			}
		}
	}()

	return id, nil
}

func (tprs *ThirdPartyRegisterService) Update(id uint, currentUserID uint, form *dtos.ThirdPartyRegisterFormDTO) error {
	tx := tprs.repo.BeginTransaction()
	existing, err := tprs.GetByID(id)
	if err != nil {
		return err
	}

	if err := tprs.validateForm(form); err != nil {
		return err
	}

	dates, err := tprs.parseDates(form)
	if err != nil {
		tx.Rollback()
		return err
	}

	if existing.ResponsibleUserID != form.ResponsibleUserID {
		responsibleUnitID, err := tprs.userSvc.GetUserPrimaryUnitID(*form.ResponsibleUserID)
		if err != nil {
			return utility.NewValidationError(fmt.Errorf("organization unit not set for responsible user"), "not_exist", "risk responsible unit")
		}
		existing.ResponsibleUnitID = &responsibleUnitID
	}

	existing.Title = form.Title
	existing.ThirdPartyTypeID = form.ThirdPartyTypeID
	existing.CountryID = form.CountryID
	existing.Address = form.Address
	existing.Website = form.Website
	existing.OrgUnitID = form.OrganizationUnitID
	existing.RepresentativeName = form.RepresentativeName
	existing.RepresentativeRole = form.RepresentativeRole
	existing.Email = form.Email
	existing.PhoneNumber = form.PhoneNumber
	existing.ResponsibleUserID = form.ResponsibleUserID
	existing.OnboardingDate = dates.OnboardingDate
	existing.ServiceStartDate = dates.ServiceStartDate
	existing.ServiceEndDate = dates.ServiceEndDate
	existing.SeverityLevelID = form.SeverityLevelID
	existing.UpdatedBYID = &currentUserID

	err = tprs.repo.Update(existing, tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	return nil
}

func (tprs *ThirdPartyRegisterService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIDs []uint

	for _, id := range ids {
		_, err := tprs.GetByID(id)
		if err != nil {
			failedIDs = append(failedIDs, id)
			continue
		}

		err = tprs.repo.Delete(id)
		if err != nil {
			failedIDs = append(failedIDs, id)
			continue
		}

	}

	return failedIDs, nil
}

func (tprs *ThirdPartyRegisterService) validateForm(form *dtos.ThirdPartyRegisterFormDTO) error {
	if _, err := tprs.thirdPartyTypeSvc.GetByID(form.ThirdPartyTypeID); err != nil {
		return utility.NewValidationError(errors.New("third party type not found"), "not_exist", "third_party_type")
	}

	if _, err := tprs.countrySvc.GetByID(form.CountryID); err != nil {
		return utility.NewValidationError(errors.New("country not found"), "not_exist", "country")
	}

	if form.ResponsibleUserID != nil {
		if _, err := tprs.userSvc.GetByID(*form.ResponsibleUserID); err != nil {
			return utility.NewValidationError(errors.New("responsible user not found"), "not_exist", "responsible_user")
		}
	}

	if form.OrganizationUnitID != nil {
		if _, err := tprs.orgUnitSvc.GetByID(*form.OrganizationUnitID); err != nil {
			return utility.NewValidationError(errors.New("organization unit not found"), "not_exist", "organization_unit")
		}
	}

	if form.SeverityLevelID != nil {
		if _, err := tprs.severityLevelSvc.GetByID(*form.SeverityLevelID); err != nil {
			return utility.NewValidationError(errors.New("severity level not found"), "not_exist", "severity_level")
		}
	}

	return nil
}

func (tprs *ThirdPartyRegisterService) parseDates(dto *dtos.ThirdPartyRegisterFormDTO) (*dtos.ThirdPartyRegisterParsedDates, error) {
	parsed := &dtos.ThirdPartyRegisterParsedDates{}
	var err error

	if dto.OnboardingDate != nil {
		parsed.OnboardingDate, err = utility.ParseDate(*dto.OnboardingDate)
		if err != nil {
			return nil, utility.NewValidationError(errors.New("invalid onboarding date"), "invalid_date", "onboarding_date")
		}
	}
	if dto.ServiceStartDate != nil {
		parsed.ServiceStartDate, err = utility.ParseDate(*dto.ServiceStartDate)
		if err != nil {
			return nil, utility.NewValidationError(errors.New("invalid service start date"), "invalid_date", "service_start_date")
		}
	}
	if dto.ServiceEndDate != nil {
		parsed.ServiceEndDate, err = utility.ParseDate(*dto.ServiceEndDate)
		if err != nil {
			return nil, utility.NewValidationError(errors.New("invalid service end date"), "invalid_date", "service_end_date")
		}
	}
	return parsed, nil
}

func (tps *ThirdPartyRegisterService) ExportThirdPartyRegister(data []*dtos.ThirdPartyRegisterListDTO, currentUserID uint) (*os.File, error) {
	// Fetch user setting for language preference
	userSetting, err := tps.userSettingSvc.GetUserSettingsByUserID(currentUserID)
	if err != nil {
		return nil, err
	}

	// Define headers based on language preference
	var headers []string
	if userSetting.Language == "ar" {
		headers = []string{
			"SL No.",
			"رقم المرجع",
			"العنوان",
			"يكتب",
			"الدولة",
			"قسم",
			"مالك داخلي",
			"المستخدم المسؤول",
			"تاريخ بدء الخدمة",
			"تاريخ انتهاء الخدمة",
			"تاريخ الصعود",
			"الحرجية",
		}
	} else {
		headers = []string{
			"SL No.",
			"Reference Number",
			"Name",
			"Type",
			"Country",
			"Department",
			"Representative Name",
			"Internal Owner",
			"Service Start Date",
			"Service End Date",
			"Onboarding Date",
			"Criticality",
		}
	}

	// Prepare rows
	rows := make([][]interface{}, len(data))
	for i, item := range data {
		var (
			thirdPartyType, country, orgUnit, severityLevel, responsibleUser string
		)

		if userSetting.Language == "ar" {
			thirdPartyType = item.ThirdPartyTypeTitleAr
			country = item.CountryTitleAr
			orgUnit = item.OrganizationUnitTitleAr
			severityLevel = item.SeverityLevelTitleAr
			responsibleUser = item.ResponsibleUserFirstNameAr + " " + item.ResponsibleUserLastNameAr
		} else {
			thirdPartyType = item.ThirdPartyTypeTitleEn
			country = item.CountryTitleEn
			orgUnit = item.OrganizationUnitTitleEn
			severityLevel = item.SeverityLevelTitleEn
			responsibleUser = item.ResponsibleUserFirstNameEn + " " + item.ResponsibleUserLastNameEn
		}

		rows[i] = []interface{}{
			i + 1,
			item.ReferenceNumber,
			item.Title,
			thirdPartyType,
			country,
			orgUnit,
			utility.GetStringValue(item.RepresentativeName),
			responsibleUser,
			utility.GetStringValue(item.ServiceStartDate),
			utility.GetStringValue(item.ServiceEndDate),
			utility.GetStringValue(item.OnboardingDate),
			severityLevel,
		}
	}

	return utility.GenerateExcelFile(headers, rows)
}

func (tps *ThirdPartyRegisterService) ForceDelete(ids []uint, currentUserID uint) ([]uint, error) {
	tx := tps.repo.BeginTransaction()
	var failedIds []uint
	var deletedData []map[string]interface{}

	for _, id := range ids {
		if err := tps.repo.ForceDelete(id, tx); err != nil {
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
		if err := tps.activityLogSvc.LogBulkActivities(
			models.DeleteAction,
			models.ThirdPartyRegister{}.TableName(),
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
