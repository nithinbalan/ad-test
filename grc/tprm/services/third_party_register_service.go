package services

import (
	"errors"
	"os"

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

	return &dto, nil
}

func (tprs *ThirdPartyRegisterService) Create(currentUserID uint, form *dtos.ThirdPartyRegisterFormDTO) (uint, error) {
	tx := tprs.repo.BeginTransaction()

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
		OrganizationUnitID: form.OrganizationUnitID,
		RepresentativeName: form.RepresentativeName,
		RepresentativeRole: form.RepresentativeRole,
		Email:              form.Email,
		PhoneNumber:        form.PhoneNumber,
		ResponsibleUserID:  form.ResponsibleUserID,
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

	existing.Title = form.Title
	existing.ThirdPartyTypeID = form.ThirdPartyTypeID
	existing.CountryID = form.CountryID
	existing.Address = form.Address
	existing.Website = form.Website
	existing.OrganizationUnitID = form.OrganizationUnitID
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
			"نوع الطرف الثالث",
			"الدولة",
			"الوحدة التنظيمية",
			"اسم الممثل",
			"دور الممثل",
			"البريد الإلكتروني",
			"رقم الهاتف",
			"المستخدم المسؤول",
			"تاريخ بدء الخدمة",
			"تاريخ انتهاء الخدمة",
			"تاريخ الانضمام",
			"مستوى الشدة",
		}
	} else {
		headers = []string{
			"SL No.",
			"Reference Number",
			"Title",
			"Third Party Type",
			"Country",
			"Organization Unit",
			"Representative Name",
			"Representative Role",
			"Email",
			"Phone Number",
			"Responsible User",
			"Service Start Date",
			"Service End Date",
			"Onboarding Date",
			"Severity Level",
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
			utility.GetStringValue(item.RepresentativeRole),
			utility.GetStringValue(item.Email),
			utility.GetStringValue(item.PhoneNumber),
			responsibleUser,
			utility.GetStringValue(item.ServiceStartDate),
			utility.GetStringValue(item.ServiceEndDate),
			utility.GetStringValue(item.OnboardingDate),
			severityLevel,
		}
	}

	return utility.GenerateExcelFile(headers, rows)
}
