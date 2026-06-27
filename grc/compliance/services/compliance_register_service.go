package services

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"jethur/config"
	commonDTOs "jethur/grc/common/dtos"
	commonServices "jethur/grc/common/services"
	complianceDTOs "jethur/grc/compliance/dtos"
	complianceRepositories "jethur/grc/compliance/repositories"
	logServices "jethur/grc/logs/services"
	masterRepositories "jethur/grc/masters/repositories"
	masterServices "jethur/grc/masters/services"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	settingsDTOs "jethur/grc/settings/dtos"
	settingsServices "jethur/grc/settings/services"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type ComplianceRegisterServiceInterface interface {
	GetAll(params *complianceDTOs.ComplianceRegisterListQueryParamsDTO) (*complianceDTOs.ComplianceRegisterListResponseDTO, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*complianceDTOs.ComplianceRegisterDetailDTO, error)
	GetByID(id uint) (*models.ComplianceRegister, error)
	Create(currentUserID uint, complianceRegisterFormDTO *complianceDTOs.ComplianceRegisterFormDTO) (uint, error)
	Update(id uint, currentUserID uint, compliance *complianceDTOs.ComplianceRegisterFormDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	ForceDelete(ids []uint, currentUserID uint) ([]uint, error)
	Archive(ids []uint, currentUserID uint) ([]uint, error)
	Restore(ids []uint, currentUserID uint) ([]uint, error)
	StatusUpdate(id uint, currentUserID uint, complianceStatusID uint, tx utility.DatabaseTransaction) error
	ReviewStatusUpdate(complianceRegister *models.ComplianceRegister, tx utility.DatabaseTransaction) error

	UpdateComplianceTotalScore(complianceID uint, totalScore float64, tx utility.DatabaseTransaction) error

	GetDocumentPath(id uint, token string) (string, error)
	GetIssuingAuthorityDocumentPath(authorityID uint, token string) (string, error)

	ExportComplianceRegisters(data []*complianceDTOs.ComplianceRegisterListDTO, currentUserID uint) (*os.File, error)

	CreateComplianceVersion(currentUserID uint, complianceRegisterFormDTO *complianceDTOs.ComplianceVersionFormDTO) (uint, error)
	GetAllVersions(id uint) ([]*complianceDTOs.ComplianceRegisterDetailDTO, error)

	UpdateNote(id uint, currentUserID uint, dto *complianceDTOs.ComplianceRegisterNoteDTO) error

	GenerateComplianceTemplate() (*os.File, error)
	ValidateComplianceTemplate(file io.Reader) ([]map[string]interface{}, error)
	ValidateImportCompliances(input []map[string]interface{}) ([]map[string]interface{}, error)
	ImportCompliances(input []map[string]interface{}, currentUserID uint) ([]uint, error)
}

// Service and implementations
type ComplianceRegisterService struct {
	repo                             complianceRepositories.ComplianceRegisterRepositoryInterface
	activityLogSvc                   logServices.ActivityLogServiceInterface
	lifecycleStatusSvc               masterServices.LifecycleStatusServiceInterface
	complianceCategorySvc            ComplianceCategoryServiceInterface
	complianceStatusSvc              ComplianceStatusServiceInterface
	reviewFrequencySvc               masterServices.ReviewFrequencyServiceInterface
	severityLevelSvc                 masterServices.SeverityLevelServiceInterface
	controlSvc                       organizationServices.ControlServiceInterface
	userSettingSvc                   organizationServices.UserSettingServiceInterface
	userSvc                          organizationServices.UserServiceInterface
	notificationSvc                  commonServices.NotificationServiceInterface
	complianceFolderSvc              ComplianceFolderServiceInterface
	complianceAssessmentFrameworkSvc ComplianceAssessmentFrameworkServiceInterface
	legislativeBodySvc               LegislativeBodyServiceInterface
	legislationTypeSvc               LegislationTypeServiceInterface
	importanceClassificationSvc      ImportanceClassificationServiceInterface
	secrecyClassificationSvc         SecrecyClassificationServiceInterface
	issuingAuthoritySvc              IssuingAuthorityServiceInterface
	organizationUnitSvc              settingsServices.OrganizationUnitServiceInterface
}

func NewComplianceRegisterService(
	repo complianceRepositories.ComplianceRegisterRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	lifecycleStatusSvc masterServices.LifecycleStatusServiceInterface,
	complianceCategorySvc ComplianceCategoryServiceInterface,
	complianceStatusSvc ComplianceStatusServiceInterface,
	reviewFrequencySvc masterServices.ReviewFrequencyServiceInterface,
	severityLevelSvc masterRepositories.SeverityLevelRepositoryInterface,
	controlSvc organizationServices.ControlServiceInterface,
	userSettingSvc organizationServices.UserSettingServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	notificationSvc commonServices.NotificationServiceInterface,
	complianceFolderSvc ComplianceFolderServiceInterface,
	complianceAssessmentFrameworkSvc ComplianceAssessmentFrameworkServiceInterface,
	legislativeBodySvc LegislativeBodyServiceInterface,
	legislationTypeSvc LegislationTypeServiceInterface,
	importanceClassificationSvc ImportanceClassificationServiceInterface,
	secrecyClassificationSvc SecrecyClassificationServiceInterface,
	issuingAuthoritySvc IssuingAuthorityServiceInterface,
	organizationUnitSvc settingsServices.OrganizationUnitServiceInterface,
) ComplianceRegisterServiceInterface {
	return &ComplianceRegisterService{
		repo:                             repo,
		activityLogSvc:                   activityLogSvc,
		lifecycleStatusSvc:               lifecycleStatusSvc,
		complianceCategorySvc:            complianceCategorySvc,
		complianceStatusSvc:              complianceStatusSvc,
		reviewFrequencySvc:               reviewFrequencySvc,
		severityLevelSvc:                 severityLevelSvc,
		controlSvc:                       controlSvc,
		userSettingSvc:                   userSettingSvc,
		userSvc:                          userSvc,
		notificationSvc:                  notificationSvc,
		complianceFolderSvc:              complianceFolderSvc,
		complianceAssessmentFrameworkSvc: complianceAssessmentFrameworkSvc,
		legislativeBodySvc:               legislativeBodySvc,
		legislationTypeSvc:               legislationTypeSvc,
		importanceClassificationSvc:      importanceClassificationSvc,
		secrecyClassificationSvc:         secrecyClassificationSvc,
		issuingAuthoritySvc:              issuingAuthoritySvc,
		organizationUnitSvc:              organizationUnitSvc,
	}
}

func (crs *ComplianceRegisterService) GetAll(params *complianceDTOs.ComplianceRegisterListQueryParamsDTO) (*complianceDTOs.ComplianceRegisterListResponseDTO, error) {
	return crs.repo.GetAll(params)
}

func (crs ComplianceRegisterService) GetByID(id uint) (*models.ComplianceRegister, error) {
	compliance, err := crs.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return compliance, nil
}

func (crs ComplianceRegisterService) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*complianceDTOs.ComplianceRegisterDetailDTO, error) {
	var complianceDTO complianceDTOs.ComplianceRegisterDetailDTO

	compliance, err := crs.repo.GetFullDetailsByID(id, tx...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	err = mapstructure.Decode(compliance, &complianceDTO)
	if err != nil {
		return nil, err
	}

	// Format time fields without time component
	if compliance.Deadline != nil {
		deadline := compliance.Deadline.Format("2006-01-02")
		complianceDTO.DeadlineString = &deadline
	}
	if compliance.IssueDate != nil {
		issueDate := compliance.IssueDate.Format("2006-01-02")
		complianceDTO.IssueDateString = &issueDate
	}
	if compliance.ExpiryDate != nil {
		expiryDate := compliance.ExpiryDate.Format("2006-01-02")
		complianceDTO.ExpiryDateString = &expiryDate
	}
	if compliance.NextReviewDate != nil {
		nextReviewDate := compliance.NextReviewDate.Format("2006-01-02")
		complianceDTO.NextReviewDateString = &nextReviewDate
	}
	if compliance.LastReviewedDate != nil {
		lastReviewDate := compliance.LastReviewedDate.Format("2006-01-02")
		complianceDTO.LastReviewedDateString = &lastReviewDate
	}

	complianceDTO.CreatedAtDateString = compliance.CreatedAt.Format("2006-01-02")
	complianceDTO.UpdatedAt = compliance.UpdatedAt

	// Handle issuing authorities date formatting
	if len(compliance.IssuingAuthorities) > 0 {
		for i, authority := range compliance.IssuingAuthorities {
			if authority.Entity != nil {
				if !authority.Entity.DateOfReceiving.IsZero() {
					dateStr := authority.Entity.DateOfReceiving.Format("2006-01-02")
					complianceDTO.IssuingAuthorities[i].Entity.DateOfReceivingString = dateStr
				}
				if authority.Entity.HijriDateOfReceiving != "" {
					complianceDTO.IssuingAuthorities[i].Entity.HijriDateOfReceivingString = authority.Entity.HijriDateOfReceiving
				}
			}
		}
	}

	return &complianceDTO, nil
}

func (crs *ComplianceRegisterService) Create(currentUserID uint, complianceRegisterFormDTO *complianceDTOs.ComplianceRegisterFormDTO) (uint, error) {
	tx := crs.repo.BeginTransaction()

	if err := crs.validateComplianceRegisterFormDTO(complianceRegisterFormDTO); err != nil {
		return 0, err
	}

	dates, err := crs.parseComplianceRegisterDates(complianceRegisterFormDTO)
	if err != nil {
		return 0, err
	}

	referenceNumber, err := crs.repo.GenerateReferenceNumber()
	if err != nil {
		return 0, err
	}

	// lifecycleStatus, err := crs.lifecycleStatusSvc.GetByType("active")
	// if err != nil {
	// 	return 0, err
	// }

	// if err := crs.validateComplianceDocumentFormDTO(complianceRegisterFormDTO.DocumentEn); err != nil {
	// 	return 0, err
	// }

	// if err := crs.validateComplianceDocumentFormDTO(complianceRegisterFormDTO.DocumentAr); err != nil {
	// 	return 0, err
	// }

	responsibleUnitID, err := crs.userSvc.GetUserPrimaryUnitID(uint(*complianceRegisterFormDTO.ResponsibleUserID))
	if err != nil {
		return 0, utility.NewValidationError(fmt.Errorf("organization unit not set for responsible user"), "not_exist", "bcp communication plan responsible unit")
	}

	complianceRegister := &models.ComplianceRegister{
		//ComplianceCategoryID:            complianceRegisterFormDTO.ComplianceCategoryID,
		ComplianceStatusID: complianceRegisterFormDTO.ComplianceStatusID,
		ResponsibleUnitID:  &responsibleUnitID,
		ResponsibleUserID:  complianceRegisterFormDTO.ResponsibleUserID,
		//ReviewFrequencyID:               complianceRegisterFormDTO.ReviewFrequencyID,
		//LastReviewNote:                  complianceRegisterFormDTO.LastReviewNote,
		//ReviewResponsibleUserID:         complianceRegisterFormDTO.ReviewResponsibleUserID,
		//SeverityLevelID:                 complianceRegisterFormDTO.SeverityLevelID,
		ReferenceNumber: referenceNumber,
		//StandardCode:                    complianceRegisterFormDTO.StandardCode,
		TitleEn:       complianceRegisterFormDTO.TitleEn,
		TitleAr:       complianceRegisterFormDTO.TitleAr,
		DescriptionEn: complianceRegisterFormDTO.DescriptionEn,
		DescriptionAr: complianceRegisterFormDTO.DescriptionAr,
		Version:       complianceRegisterFormDTO.Version,
		Deadline:      dates.Deadline,
		//IssueDate:                       dates.IssueDate,
		//ExpiryDate:                      dates.ExpiryDate,
		//NextReviewDate:                  dates.NextReviewDate,
		//LastReviewedDate:                dates.LastReviewedDate,
		//LastReviewedUserID:              complianceRegisterFormDTO.LastReviewedUserID,
		//ComplianceReviewStatusID:        complianceRegisterFormDTO.ComplianceReviewStatusID,
		CreatedBYID: currentUserID,
		//LifecycleStatusID:               &lifecycleStatus.ID,
		ComplianceFolderID: complianceRegisterFormDTO.ComplianceFolderID,
		//ComplianceAssessmentFrameworkID: complianceRegisterFormDTO.ComplianceAssessmentFrameworkID,
		//IsControlAssessment:             complianceRegisterFormDTO.IsControlAssessment,
		IsRequirementAssessment: true,
		IsLatest:                true,

		//New Fields
		LegislativeBodyID:          complianceRegisterFormDTO.LegislativeBodyID,
		LegislationTypeID:          complianceRegisterFormDTO.LegislationTypeID,
		LegislationNo:              complianceRegisterFormDTO.LegislationNo,
		ImportanceClassificationID: complianceRegisterFormDTO.ImportanceClassificationID,
		SecrecyClassificationID:    complianceRegisterFormDTO.SecrecyClassificationID,
		Subject:                    complianceRegisterFormDTO.Subject,
		IsApplicable:               complianceRegisterFormDTO.IsApplicable,
		IsNotApplicable:            complianceRegisterFormDTO.IsNotApplicable,
		// IsRegulationNotificationAssessment: complianceRegisterFormDTO.IsRegulationNotificationAssessment,
		// IsRiskBasedAssessment:              complianceRegisterFormDTO.IsRiskBasedAssessment,
		IsPublic: complianceRegisterFormDTO.IsPublic,
		IsShared: complianceRegisterFormDTO.IsShared,
	}

	id, err := crs.repo.Create(complianceRegister, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	if complianceRegisterFormDTO.ComplianceStatusID != nil {

		complianceStatusID := complianceRegisterFormDTO.ComplianceStatusID
		complianceStatusUpdate := &models.ComplianceStatusUpdate{
			ComplianceRegisterID: id,
			ComplianceStatusID:   *complianceStatusID,
			Comment:              "Initial compliance status",
			CreatedByID:          currentUserID,
		}

		_, err = crs.repo.CreateInitialComplianceStatusUpdate(complianceStatusUpdate, tx)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	// if err := crs.handleDocument(currentUserID, id, complianceRegisterFormDTO.DocumentEn, tx); err != nil {
	// 	tx.Rollback()
	// 	return 0, err
	// }

	// if err := crs.handleDocument(currentUserID, id, complianceRegisterFormDTO.DocumentAr, tx); err != nil {
	// 	tx.Rollback()
	// 	return 0, err
	// }

	// Handle issuing authorities
	if err := crs.handleIssuingAuthorities(currentUserID, id, complianceRegisterFormDTO.IssuingAuthorities, complianceRegisterFormDTO.SecrecyClassificationID, tx); err != nil {
		tx.Rollback()
		return 0, err
	}

	// Validate shared units requirement
	if complianceRegisterFormDTO.IsShared && len(complianceRegisterFormDTO.SharedUnitIDs) == 0 {
		tx.Rollback()
		return 0, utility.NewValidationError(fmt.Errorf("shared units are required when IsShared is true"), "required", "shared_unit_ids")
	}

	// Handle shared units
	if err := crs.handleSharedUnits(currentUserID, id, complianceRegisterFormDTO.SharedUnitIDs, complianceRegisterFormDTO.IsShared, tx); err != nil {
		tx.Rollback()
		return 0, err
	}

	tx.Commit()

	go func() {
		existingCompliance, err := crs.GetFullDetailsByID(complianceRegister.ID)
		if err != nil {
			utility.LogError("Failed to fetch compliance register", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		err = crs.activityLogSvc.LogActivity(
			models.CreateAction,
			models.ComplianceRegister{}.TableName(),
			id,
			existingCompliance,
			currentUserID,
			"compliance", "compliance_register", "create compliance", nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.ComplianceRegister{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}

		// Prepare notification and email
		urlParam := fmt.Sprintf(`{"id":%d}`, id)
		entity := "compliance-register"
		action := "responsible-user-assign"

		dbNotificationDTO := &commonDTOs.DBNotificationDTO{
			UserIDs:              []uint{*complianceRegisterFormDTO.ResponsibleUserID},
			ModuleIdentifier:     "compliance",
			NotificationType:     string(commonDTOs.AssignmentType),
			HeadingEn:            "Compliance Register Responsibility Assigned",
			HeadingAr:            "تم تعيين مسؤولية سجل الامتثال",
			TitleEn:              complianceRegisterFormDTO.TitleEn,
			TitleAr:              complianceRegisterFormDTO.TitleAr,
			ReferenceNumber:      complianceRegister.ReferenceNumber,
			Entity:               &entity,
			Action:               &action,
			URLParam:             urlParam,
			IsSystemNotification: false,
			CreatedByID:          &currentUserID,
		}

		responsibleUser, err := crs.userSvc.GetByID(*complianceRegisterFormDTO.ResponsibleUserID)
		if err != nil {
			utility.LogError("Failed to fetch responsible user", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		createdBy, err := crs.userSvc.GetByID(currentUserID)
		if err != nil {
			utility.LogError("Failed to fetch responsible user", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		// Check if severity level ID is present before fetching
		var severityLevel *models.SeverityLevel
		if complianceRegister.SeverityLevelID != nil {
			severityLevel, err = crs.severityLevelSvc.GetByID(*complianceRegister.SeverityLevelID)
			if err != nil {
				utility.LogError("Failed to fetch severity level", map[string]interface{}{
					"error":   err.Error(),
					"levelID": complianceRegister.SeverityLevelID,
				})
				return
			}
		}

		emailTo := []string{responsibleUser.Email}
		subject := fmt.Sprintf("Compliance Register/%s: Responsibility Assigned", complianceRegister.TitleEn)
		appURL := utility.GetAppURL()

		severityLevelTitle := "N/A"
		if severityLevel != nil {
			severityLevelTitle = severityLevel.TitleEn
		}

		emailData := map[string]interface{}{
			"FirstNameEn":               responsibleUser.FirstNameEn,
			"LastNameEn":                responsibleUser.LastNameEn,
			"ComplianceReferenceNumber": complianceRegister.ReferenceNumber,
			"ComplianceTitleEn":         complianceRegister.TitleEn,
			"ComplianceTitleAr":         complianceRegister.TitleAr,
			"ComplianceDescription":     complianceRegister.DescriptionEn,
			"SeverityLevel":             severityLevelTitle,
			"CreatedByFirstNameEn":      createdBy.FirstNameEn,
			"CreatedByLastNameEn":       createdBy.LastNameEn,
			"AssignedOn":                time.Now().Format("2006-01-02"),
			"RedirectURL":               fmt.Sprintf("%s/compliance/register/%d", appURL, id),
		}

		emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
			EmailTo:              emailTo,
			Subject:              subject,
			TemplateRelativePath: "compliance/templates/email/compliance_register_create_responsible_user_mail.html",
			EmailData:            emailData,
		}

		err = crs.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "compliance_register")
		if err != nil {
			utility.LogError("Failed to send notifications", map[string]interface{}{
				"error": err.Error(),
			})
		}
	}()

	return id, nil
}

func (crs *ComplianceRegisterService) Update(id uint, currentUserID uint, complianceRegisterFormDTO *complianceDTOs.ComplianceRegisterFormDTO) error {
	tx := crs.repo.BeginTransaction()

	existingComplianceRegister, err := crs.GetByID(id)
	if err != nil {
		return err
	}

	oldResponsibleUserID := existingComplianceRegister.ResponsibleUserID

	if err := crs.validateComplianceRegisterFormDTO(complianceRegisterFormDTO); err != nil {
		return err
	}

	dates, err := crs.parseComplianceRegisterDates(complianceRegisterFormDTO)
	if err != nil {
		return err
	}

	// if err := crs.validateComplianceDocumentFormDTO(complianceRegisterFormDTO.DocumentEn); err != nil {
	// 	return err
	// }

	// if err := crs.validateComplianceDocumentFormDTO(complianceRegisterFormDTO.DocumentAr); err != nil {
	// 	return err
	// }

	oldComplianceRegister, err := crs.GetFullDetailsByID(id, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	if existingComplianceRegister.ResponsibleUserID != complianceRegisterFormDTO.ResponsibleUserID {
		responsibleUnitID, err := crs.userSvc.GetUserPrimaryUnitID(*complianceRegisterFormDTO.ResponsibleUserID)
		if err != nil {
			return utility.NewValidationError(fmt.Errorf("organization unit not set for responsible user"), "not_exist", "risk responsible unit")
		}
		existingComplianceRegister.ResponsibleUnitID = &responsibleUnitID
	}
	// hasAssessments, err := crs.repo.HasComplianceAssessments(id, tx)
	// if err != nil {
	// 	return err
	// }

	// if hasAssessments && existingComplianceRegister.ComplianceAssessmentFrameworkID != nil &&
	// 	complianceRegisterFormDTO.ComplianceAssessmentFrameworkID != nil &&
	// 	existingComplianceRegister.ComplianceAssessmentFrameworkID != complianceRegisterFormDTO.ComplianceAssessmentFrameworkID {
	// 	return utility.NewForbiddenError("compliance_assessment_framework_cannot_be_changed")
	// }
	// // Prevent changing IsControlAssessment if assessments exist
	// if hasAssessments && existingComplianceRegister.IsControlAssessment != complianceRegisterFormDTO.IsControlAssessment {
	// 	return utility.NewForbiddenError("is_control_assessment_cannot_be_changed")
	// }

	// // Prevent changing IsRequirementAssessment if assessments exist
	// if hasAssessments && existingComplianceRegister.IsRequirementAssessment != complianceRegisterFormDTO.IsRequirementAssessment {
	// 	return utility.NewForbiddenError("is_requirement_assessment_cannot_be_changed")
	// }

	//existingComplianceRegister.ComplianceCategoryID = complianceRegisterFormDTO.ComplianceCategoryID
	existingComplianceRegister.ComplianceStatusID = complianceRegisterFormDTO.ComplianceStatusID
	existingComplianceRegister.ResponsibleUserID = complianceRegisterFormDTO.ResponsibleUserID
	// existingComplianceRegister.ReviewFrequencyID = complianceRegisterFormDTO.ReviewFrequencyID
	// existingComplianceRegister.ReviewResponsibleUserID = complianceRegisterFormDTO.ReviewResponsibleUserID
	// existingComplianceRegister.SeverityLevelID = complianceRegisterFormDTO.SeverityLevelID
	// existingComplianceRegister.StandardCode = complianceRegisterFormDTO.StandardCode
	existingComplianceRegister.TitleEn = complianceRegisterFormDTO.TitleEn
	existingComplianceRegister.TitleAr = complianceRegisterFormDTO.TitleAr
	existingComplianceRegister.DescriptionEn = complianceRegisterFormDTO.DescriptionEn
	existingComplianceRegister.DescriptionAr = complianceRegisterFormDTO.DescriptionAr
	existingComplianceRegister.Version = complianceRegisterFormDTO.Version
	existingComplianceRegister.Deadline = dates.Deadline
	// existingComplianceRegister.IssueDate = dates.IssueDate
	// existingComplianceRegister.ExpiryDate = dates.ExpiryDate
	// existingComplianceRegister.NextReviewDate = dates.NextReviewDate
	// existingComplianceRegister.LastReviewedDate = dates.LastReviewedDate
	// existingComplianceRegister.LastReviewedUserID = complianceRegisterFormDTO.LastReviewedUserID
	// existingComplianceRegister.ComplianceReviewStatusID = complianceRegisterFormDTO.ComplianceReviewStatusID
	// existingComplianceRegister.LastReviewNote = complianceRegisterFormDTO.LastReviewNote
	existingComplianceRegister.UpdatedBYID = &currentUserID
	existingComplianceRegister.ComplianceFolderID = complianceRegisterFormDTO.ComplianceFolderID
	// existingComplianceRegister.ComplianceAssessmentFrameworkID = complianceRegisterFormDTO.ComplianceAssessmentFrameworkID
	// existingComplianceRegister.IsControlAssessment = complianceRegisterFormDTO.IsControlAssessment
	//existingComplianceRegister.IsRequirementAssessment = complianceRegisterFormDTO.IsRequirementAssessment

	//New Fields
	existingComplianceRegister.LegislativeBodyID = complianceRegisterFormDTO.LegislativeBodyID
	existingComplianceRegister.LegislationTypeID = complianceRegisterFormDTO.LegislationTypeID
	existingComplianceRegister.LegislationNo = complianceRegisterFormDTO.LegislationNo
	existingComplianceRegister.ImportanceClassificationID = complianceRegisterFormDTO.ImportanceClassificationID
	existingComplianceRegister.SecrecyClassificationID = complianceRegisterFormDTO.SecrecyClassificationID
	existingComplianceRegister.Subject = complianceRegisterFormDTO.Subject
	existingComplianceRegister.IsApplicable = complianceRegisterFormDTO.IsApplicable
	existingComplianceRegister.IsNotApplicable = complianceRegisterFormDTO.IsNotApplicable
	// existingComplianceRegister.IsRegulationNotificationAssessment = complianceRegisterFormDTO.IsRegulationNotificationAssessment
	// existingComplianceRegister.IsRiskBasedAssessment = complianceRegisterFormDTO.IsRiskBasedAssessment
	existingComplianceRegister.IsPublic = complianceRegisterFormDTO.IsPublic
	existingComplianceRegister.IsShared = complianceRegisterFormDTO.IsShared

	err = crs.repo.Update(existingComplianceRegister, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	// existingComplianceDocumentEn, err := crs.repo.GetLatestDocumentByComplianceID(id, "en")
	// if err != nil {
	// 	return err
	// }

	// existingComplianceDocumentAr, err := crs.repo.GetLatestDocumentByComplianceID(id, "ar")
	// if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
	// 	return err
	// }

	// if complianceRegisterFormDTO.DocumentEn.IsNewFile != nil && *complianceRegisterFormDTO.DocumentEn.IsNewFile {
	// 	existingComplianceDocumentEnFileName := existingComplianceDocumentEn.Name

	// 	if err := crs.handleDocument(currentUserID, id, complianceRegisterFormDTO.DocumentEn, tx); err != nil {
	// 		tx.Rollback()
	// 		return err
	// 	}

	// 	storagePath := filepath.Join("compliance", "compliance-register", strconv.FormatUint(uint64(id), 10), "documents")
	// 	if err := utility.RemoveFile(existingComplianceDocumentEnFileName, storagePath); err != nil {
	// 		tx.Rollback()
	// 		return err
	// 	}

	// }

	// if complianceRegisterFormDTO.DocumentAr.IsNewFile != nil && *complianceRegisterFormDTO.DocumentAr.IsNewFile {
	// 	var existingComplianceDocumentArFileName string
	// 	if existingComplianceDocumentAr != nil {
	// 		existingComplianceDocumentArFileName = existingComplianceDocumentAr.Name
	// 	}

	// 	if err := crs.handleDocument(currentUserID, id, complianceRegisterFormDTO.DocumentAr, tx); err != nil {
	// 		tx.Rollback()
	// 		return err
	// 	}

	// 	if existingComplianceDocumentAr != nil {
	// 		storagePath := filepath.Join("compliance", "compliance-register", strconv.FormatUint(uint64(id), 10), "documents")
	// 		if err := utility.RemoveFile(existingComplianceDocumentArFileName, storagePath); err != nil {
	// 			tx.Rollback()
	// 			return err
	// 		}
	// 	}
	// }

	// Handle issuing authorities
	if err := crs.handleIssuingAuthorities(currentUserID, id, complianceRegisterFormDTO.IssuingAuthorities, complianceRegisterFormDTO.SecrecyClassificationID, tx); err != nil {
		tx.Rollback()
		return err
	}

	// Validate shared units requirement
	if complianceRegisterFormDTO.IsShared && len(complianceRegisterFormDTO.SharedUnitIDs) == 0 {
		tx.Rollback()
		return utility.NewValidationError(fmt.Errorf("shared units are required when IsShared is true"), "required", "shared_unit_ids")
	}

	// Handle shared units
	if err := crs.handleSharedUnits(currentUserID, id, complianceRegisterFormDTO.SharedUnitIDs, complianceRegisterFormDTO.IsShared, tx); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	go func() {
		// existingCompliance, err := crs.GetFullDetailsByID(id)
		// if err != nil {
		// 	return
		// }

		// Prepare log data
		changedData := make(map[string]interface{})

		// log changes
		oldComplianceStatus := ""
		if oldComplianceRegister.ComplianceStatus != nil {
			oldComplianceStatus = oldComplianceRegister.ComplianceStatus.TitleEn
		}

		oldResponsibleUnit := ""
		if oldComplianceRegister.ResponsibleUnit.ID != 0 {
			oldResponsibleUnit = oldComplianceRegister.ResponsibleUnit.TitleEn
		}

		oldResponsibleUser := ""
		if oldComplianceRegister.ResponsibleUser.ID != 0 {
			oldResponsibleUser = fmt.Sprintf("%s %s", oldComplianceRegister.ResponsibleUser.FirstNameEN, oldComplianceRegister.ResponsibleUser.LastNameEN)
		}

		oldData := map[string]interface{}{
			// "compliance_category":      oldComplianceRegister.ComplianceCategory.TitleEn,
			"compliance_status": oldComplianceStatus,
			"responsible_unit":  oldResponsibleUnit,
			"responsible_user":  oldResponsibleUser,
			// "review_frequency":         utility.GetFieldValue(oldComplianceRegister.ReviewFrequency).(string),
			// "review_responsible_user":  utility.GetFieldValue(oldComplianceRegister.ReviewResponsibleUser).(string),
			// "severity_level":           utility.GetFieldValue(oldComplianceRegister.SeverityLevel).(string),
			// "standard_code":            oldComplianceRegister.StandardCode,
			"title_en":       oldComplianceRegister.TitleEn,
			"title_ar":       oldComplianceRegister.TitleAr,
			"description_en": oldComplianceRegister.DescriptionEn,
			"description_ar": oldComplianceRegister.DescriptionAr,
			"version":        oldComplianceRegister.Version,
			"dead_line":      oldComplianceRegister.DeadlineString,
			// "issue_date":               oldComplianceRegister.IssueDateString,
			// "expiry_date":              oldComplianceRegister.ExpiryDateString,
			// "last_review_date":         oldComplianceRegister.LastReviewedDateString,
			// "next_review_date":         oldComplianceRegister.NextReviewDateString,
			// "last_review_note":         oldComplianceRegister.LastReviewNote,
			// "last_review_user":         utility.GetFieldValue(oldComplianceRegister.LastReviewedUser).(string),
			"compliance_review_status":     utility.GetFieldValue(oldComplianceRegister.ComplianceReviewStatus).(string),
			"legislative_body_id":          oldComplianceRegister.LegislativeBodyID,
			"legislation_type_id":          oldComplianceRegister.LegislationTypeID,
			"legislation_no":               oldComplianceRegister.LegislationNo,
			"importance_classification_id": oldComplianceRegister.ImportanceClassificationID,
			"secrecy_classification_id":    oldComplianceRegister.SecrecyClassificationID,
			"subject":                      oldComplianceRegister.Subject,
			"is_applicable":                oldComplianceRegister.IsApplicable,
			"is_not_applicable":            oldComplianceRegister.IsNotApplicable,
			// "is_regulation_notification_assessment": oldComplianceRegister.IsRegulationNotificationAssessment,
			// "is_risk_based_assessment":              oldComplianceRegister.IsRiskBasedAssessment,
		}

		newData := map[string]interface{}{
			// "compliance_category":      oldComplianceRegister.ComplianceCategory.TitleEn,
			"compliance_status": oldComplianceStatus,
			"responsible_unit":  oldResponsibleUnit,
			// "responsible_user":         fmt.Sprintf("%s %s", oldComplianceRegister.ResponsibleUser.FirstNameEN, oldComplianceRegister.ResponsibleUser.LastNameEN),
			// "review_frequency":         utility.GetFieldValue(oldComplianceRegister.ReviewFrequency).(string),
			// "review_responsible_user":  utility.GetFieldValue(oldComplianceRegister.ReviewResponsibleUser).(string),
			// "severity_level":           utility.GetFieldValue(oldComplianceRegister.SeverityLevel).(string),
			// "standard_code":            oldComplianceRegister.StandardCode,
			"title_en":       existingComplianceRegister.TitleEn,
			"title_ar":       existingComplianceRegister.TitleAr,
			"description_en": existingComplianceRegister.DescriptionEn,
			"description_ar": existingComplianceRegister.DescriptionAr,
			"version":        existingComplianceRegister.Version,
			"dead_line":      existingComplianceRegister.Deadline,
			// "issue_date":               oldComplianceRegister.IssueDateString,
			// "expiry_date":              oldComplianceRegister.ExpiryDateString,
			// "last_review_date":         oldComplianceRegister.LastReviewedDateString,
			// "next_review_date":         oldComplianceRegister.NextReviewDateString,
			// "last_review_note":         oldComplianceRegister.LastReviewNote,
			// "last_review_user":         utility.GetFieldValue(oldComplianceRegister.LastReviewedUser).(string),
			"compliance_review_status":     utility.GetFieldValue(existingComplianceRegister.ComplianceReviewStatus).(string),
			"legislative_body_id":          existingComplianceRegister.LegislativeBodyID,
			"legislation_type_id":          existingComplianceRegister.LegislationTypeID,
			"legislation_no":               existingComplianceRegister.LegislationNo,
			"importance_classification_id": existingComplianceRegister.ImportanceClassificationID,
			"secrecy_classification_id":    existingComplianceRegister.SecrecyClassificationID,
			"subject":                      existingComplianceRegister.Subject,
			"is_applicable":                existingComplianceRegister.IsApplicable,
			"is_not_applicable":            existingComplianceRegister.IsNotApplicable,
			// "is_regulation_notification_assessment": existingComplianceRegister.IsRegulationNotificationAssessment,
			// "is_risk_based_assessment":              existingComplianceRegister.IsRiskBasedAssessment,
		}
		// Compare other fields and add them to changedData
		fieldChanges := utility.CompareFieldChanges(oldData, newData)
		for k, v := range fieldChanges {
			changedData[k] = v
		}
		// Log activity if there are any changes
		if len(changedData) > 0 {
			// Log the update activity
			err := crs.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.ComplianceRegister{}.TableName(),
				id,
				changedData,
				currentUserID,
				"compliance", "compliance_register", "", nil,
			)
			if err != nil {
				utility.LogError("Failed to log activity", map[string]interface{}{
					"action":    models.UpdateAction,
					"tableName": models.ComplianceRegister{}.TableName(),
					"recordID":  id,
					"userID":    currentUserID,
					"error":     err.Error(),
				})
			}
		}

		// Send notifications to old and new responsible users if the responsible user has changed
		if oldResponsibleUserID != existingComplianceRegister.ResponsibleUserID {
			// URL param for notifications
			urlParam := fmt.Sprintf(`{"id":%d}`, id)
			entity := "compliance-register"

			// Old responsible user notification (Removal)
			oldResponsibleUserAction := "responsible-user-remove"
			oldResponsibleUserNotificationDTO := &commonDTOs.DBNotificationDTO{
				UserIDs:              []uint{*oldResponsibleUserID},
				ModuleIdentifier:     "compliance",
				NotificationType:     string(commonDTOs.InformationType),
				HeadingEn:            "Compliance Register Responsibility Removed",
				HeadingAr:            "تمت إزالة مسؤولية تسجيل الامتثال",
				TitleEn:              existingComplianceRegister.TitleEn,
				TitleAr:              existingComplianceRegister.TitleAr,
				ReferenceNumber:      existingComplianceRegister.ReferenceNumber,
				Entity:               &entity,
				Action:               &oldResponsibleUserAction,
				URLParam:             urlParam,
				IsSystemNotification: false,
				CreatedByID:          &currentUserID,
			}

			// New responsible user notification (Assignment)
			newResponsibleUserAction := "responsible-user-assign"
			newResponsibleUserNotificationDTO := &commonDTOs.DBNotificationDTO{
				UserIDs:              []uint{*complianceRegisterFormDTO.ResponsibleUserID},
				ModuleIdentifier:     "compliance",
				NotificationType:     string(commonDTOs.AssignmentType),
				HeadingEn:            "Compliance Register Responsibility Assigned",
				HeadingAr:            "تم تعيين مسؤولية سجل الامتثال",
				TitleEn:              existingComplianceRegister.TitleEn,
				TitleAr:              existingComplianceRegister.TitleAr,
				ReferenceNumber:      existingComplianceRegister.ReferenceNumber,
				Entity:               &entity,
				Action:               &newResponsibleUserAction,
				URLParam:             urlParam,
				IsSystemNotification: false,
				CreatedByID:          &currentUserID,
			}

			// Fetch the responsible users' data
			oldResponsibleUser, err := crs.userSvc.GetByID(*oldResponsibleUserID)
			if err != nil {
				utility.LogError("Failed to fetch old responsible user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}

			newResponsibleUser, err := crs.userSvc.GetByID(*complianceRegisterFormDTO.ResponsibleUserID)
			if err != nil {
				utility.LogError("Failed to fetch new responsible user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}

			// Fetch details of the user making the assignment
			updatedByUser, err := crs.userSvc.GetByID(currentUserID)
			if err != nil {
				utility.LogError("Failed to fetch assigning user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}

			appURL := utility.GetAppURL()

			// Email notification for the old responsible user
			oldResponsibleUserEmailTo := []string{oldResponsibleUser.Email}
			oldResponsibleUserSubject := fmt.Sprintf("Compliance Register/%s: Responsibility Removed", existingComplianceRegister.TitleEn)
			oldResponsibleUserEmailData := map[string]interface{}{
				"FirstNameEn":               oldResponsibleUser.FirstNameEn,
				"LastNameEn":                oldResponsibleUser.LastNameEn,
				"ComplianceReferenceNumber": existingComplianceRegister.ReferenceNumber,
				"ComplianceTitle":           existingComplianceRegister.TitleEn,
				"ComplianceTitleAr":         existingComplianceRegister.TitleAr,
				"ComplianceDescription":     existingComplianceRegister.DescriptionEn,
				"UpdatedByFirstNameEn":      updatedByUser.FirstNameEn,
				"UpdatedByLastNameEn":       updatedByUser.LastNameEn,
				"RedirectURL":               fmt.Sprintf("%s/compliance/register/%d", appURL, id),
			}

			oldResponsibleUserEmailNotificationDTO := &commonDTOs.EmailNotificationDTO{
				EmailTo:              oldResponsibleUserEmailTo,
				Subject:              oldResponsibleUserSubject,
				TemplateRelativePath: "compliance/templates/email/compliance_register_remove_responsible_user_mail.html",
				EmailData:            oldResponsibleUserEmailData,
			}

			// Email notification for the new responsible user
			newResponsibleUserEmailTo := []string{newResponsibleUser.Email}
			newResponsibleUserSubject := fmt.Sprintf("Compliance Register/%s: Responsibility Assigned", existingComplianceRegister.TitleEn)
			newResponsibleUserEmailData := map[string]interface{}{
				"FirstNameEn":               newResponsibleUser.FirstNameEn,
				"LastNameEn":                newResponsibleUser.LastNameEn,
				"ComplianceReferenceNumber": existingComplianceRegister.ReferenceNumber,
				"ComplianceTitle":           existingComplianceRegister.TitleEn,
				"ComplianceTitleAr":         existingComplianceRegister.TitleAr,
				"ComplianceDescription":     existingComplianceRegister.DescriptionEn,
				"AssignedByFirstNameEn":     updatedByUser.FirstNameEn,
				"AssignedByLastNameEn":      updatedByUser.LastNameEn,
				"AssignedOn":                time.Now().Format("2006-01-02"),
				"RedirectURL":               fmt.Sprintf("%s/compliance/register/%d", appURL, id),
			}

			newResponsibleUserEmailNotificationDTO := &commonDTOs.EmailNotificationDTO{
				EmailTo:              newResponsibleUserEmailTo,
				Subject:              newResponsibleUserSubject,
				TemplateRelativePath: "compliance/templates/email/compliance_register_change_responsible_user_mail.html",
				EmailData:            newResponsibleUserEmailData,
			}

			// Send notifications and emails
			err = crs.notificationSvc.SendNotification(newResponsibleUserNotificationDTO, newResponsibleUserEmailNotificationDTO, "compliance_register")
			if err != nil {
				utility.LogError("Failed to send notifications to new responsible user", map[string]interface{}{
					"error": err.Error(),
				})
			}

			err = crs.notificationSvc.SendNotification(oldResponsibleUserNotificationDTO, oldResponsibleUserEmailNotificationDTO, "compliance_register")
			if err != nil {
				utility.LogError("Failed to send notifications to old responsible user", map[string]interface{}{
					"error": err.Error(),
				})
			}
		}
	}()

	return nil
}

func (crs *ComplianceRegisterService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	tx := crs.repo.BeginTransaction()

	var failedIds []uint
	var successfulIds []uint
	var deletedComplianceRegister []map[string]interface{}

	var notifications []struct {
		OwnerID         *uint
		TitleEn         string
		TitleAr         string
		Description     *string
		ReferenceNumber string
	}

	documentsToDelete := make(map[uint][]*complianceDTOs.ComplianceDocumentDTO)

	for _, id := range ids {
		complianceRegister, err := crs.GetFullDetailsByID(id, tx)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		// Prepare notification data before deletion
		notifications = append(notifications, struct {
			OwnerID         *uint
			TitleEn         string
			TitleAr         string
			Description     *string
			ReferenceNumber string
		}{
			OwnerID:         &complianceRegister.ResponsibleUserID,
			TitleEn:         complianceRegister.TitleEn,
			TitleAr:         complianceRegister.TitleAr,
			Description:     complianceRegister.DescriptionEn,
			ReferenceNumber: complianceRegister.ReferenceNumber,
		})

		if len(complianceRegister.ComplianceStatusUpdates) > 0 {
			var deleteStatusUpdateIds []*uint
			for _, update := range complianceRegister.ComplianceStatusUpdates {
				deleteStatusUpdateIds = append(deleteStatusUpdateIds, &update.ID)
			}
			if err := crs.repo.DeleteComplianceStatusUpdates(deleteStatusUpdateIds, tx); err != nil {
				tx.Rollback()
				return nil, err
			}
		}

		// Delete shared units
		existingSharedUnits, err := crs.repo.GetSharedUnitsByComplianceRegisterID(id, tx)
		if err == nil && len(existingSharedUnits) > 0 {
			var departmentIDs []uint
			for _, sharedUnit := range existingSharedUnits {
				departmentIDs = append(departmentIDs, sharedUnit.DepartmentID)
			}
			if err := crs.repo.DeleteSharedUnits(id, departmentIDs, tx); err != nil {
				tx.Rollback()
				return nil, err
			}
		}

		// Delete issuing authorities and their related records
		if len(complianceRegister.IssuingAuthorities) > 0 {
			for _, authority := range complianceRegister.IssuingAuthorities {
				// Delete issuing authority document if exists
				if authority.Document != nil {
					if err := crs.repo.DeleteIssuingAuthorityDocument(authority.Document.ID, tx); err != nil {
						tx.Rollback()
						return nil, err
					}
				}

				// Delete issuing authority entity if exists
				if authority.Entity != nil {
					if err := crs.repo.DeleteIssuingAuthorityEntity(authority.Entity.ID, tx); err != nil {
						tx.Rollback()
						return nil, err
					}
				}

				// Delete the issuing authority itself
				if err := crs.repo.DeleteIssuingAuthority(authority.ID, tx); err != nil {
					tx.Rollback()
					return nil, err
				}
			}
		}

		// Delete each Compliance document if documents exist
		var docIdsToDelete []*uint

		if complianceRegister.DocumentEn != nil {
			documentsToDelete[id] = append(documentsToDelete[id], complianceRegister.DocumentEn)
			docIdsToDelete = append(docIdsToDelete, complianceRegister.DocumentEn.ID)
		}

		if complianceRegister.DocumentAr != nil {
			documentsToDelete[id] = append(documentsToDelete[id], complianceRegister.DocumentAr)
			docIdsToDelete = append(docIdsToDelete, complianceRegister.DocumentAr.ID)
		}

		// Delete all collected documents in bulk
		if len(docIdsToDelete) > 0 {
			if err := crs.repo.DeleteDocuments(docIdsToDelete, tx); err != nil {
				tx.Rollback()
				return nil, err
			}
		}

		// Delete the Compliance Register
		if err := crs.repo.Delete(id, tx); err != nil {
			failedIds = append(failedIds, id)
			continue
		} else {
			successfulIds = append(successfulIds, id)
		}

		deletedComplianceRegister = append(deletedComplianceRegister, map[string]interface{}{
			"id":       id,
			"title_en": complianceRegister.TitleEn,
			"title_ar": complianceRegister.TitleAr,
		})
	}

	if len(failedIds) > 0 {
		tx.Rollback()
		return failedIds, nil
	}

	// Commit DB changes before file cleanup
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	// Deleting files should not affect the main operation result (cleanup after commit)
	if len(successfulIds) > 0 {
		for complianceID, documents := range documentsToDelete {
			for _, doc := range documents {
				storagePath := filepath.Join("compliance", "compliance-register", strconv.FormatUint(uint64(complianceID), 10), "documents")
				if err := utility.RemoveFile(doc.Name, storagePath); err != nil {
					// Log error but don't fail the operation since DB is already committed
					utility.LogError("Failed to delete file from storage", map[string]interface{}{
						"url":   doc.URL,
						"error": err.Error(),
					})
				}
			}
			destinationPath := filepath.Join("compliance", "compliance-register", strconv.FormatUint(uint64(complianceID), 10))
			if err := utility.RemoveFolder(destinationPath); err != nil {
				// Log error but don't fail the operation since DB is already committed
				utility.LogError("Failed to remove folder", map[string]interface{}{
					"path":  destinationPath,
					"error": err.Error(),
				})
			}
		}
	}

	// Log bulk delete activities for successful deletions
	if len(successfulIds) > 0 {
		go func() {
			err := crs.activityLogSvc.LogBulkActivities(
				models.DeleteAction, // Make sure this action constant is defined appropriately
				models.ComplianceRegister{}.TableName(),
				successfulIds,
				deletedComplianceRegister,
				currentUserID,
				nil,
			)
			if err != nil {
				// @todo - save to log file
				fmt.Printf("Failed to log bulk delete activity: %v\n", err)
			}

			//Send DB and email notifications for each successfully deleted compliance register
			for _, notification := range notifications {
				urlParam := ""
				entity := "compliance-register"
				action := "delete"

				// Prepare DB notification
				dbNotificationDTO := &commonDTOs.DBNotificationDTO{
					UserIDs:              []uint{*notification.OwnerID},
					ModuleIdentifier:     "compliance",
					NotificationType:     string(commonDTOs.DeletionType),
					HeadingEn:            "Compliance Register Deleted",
					HeadingAr:            "تم حذف سجل الامتثال",
					TitleEn:              notification.TitleEn,
					TitleAr:              notification.TitleAr,
					ReferenceNumber:      notification.ReferenceNumber,
					Entity:               &entity,
					Action:               &action,
					URLParam:             urlParam,
					IsSystemNotification: false,
					CreatedByID:          &currentUserID,
				}

				// Fetch owner details (responsible user)
				ownerUser, err := crs.userSvc.GetByID(*notification.OwnerID)
				if err != nil {
					utility.LogError("Failed to fetch owner details for notification", map[string]interface{}{
						"error": err.Error(),
					})
					continue
				}

				// Prepare email notification
				ownerEmailTo := []string{ownerUser.Email}
				subject := fmt.Sprintf("Compliance Register/%s: Deleted", notification.TitleEn)
				templateRelativePath := "compliance/templates/email/compliance_register_delete_mail.html"

				emailData := map[string]interface{}{
					"FirstNameEn":               ownerUser.FirstNameEn,
					"LastNameEn":                ownerUser.LastNameEn,
					"ComplianceReferenceNumber": notification.ReferenceNumber,
					"ComplianceTitle":           notification.TitleEn,
					"ComplianceTitleAr":         notification.TitleAr,
					"ComplianceDescription":     notification.Description,
				}

				emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
					EmailTo:              ownerEmailTo,
					Subject:              subject,
					TemplateRelativePath: templateRelativePath,
					EmailData:            emailData,
				}

				// Send DB and email notifications
				err = crs.notificationSvc.SendNotification(
					dbNotificationDTO, emailNotificationDTO, "compliance_register")
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

func (crs *ComplianceRegisterService) Archive(ids []uint, currentUserID uint) ([]uint, error) {
	lifecycleStatus, err := crs.lifecycleStatusSvc.GetByType("archived")
	if err != nil {
		return nil, err
	}
	idsToUpdate, err := crs.repo.GetIDsByLifecycleStatus(ids, "active")
	if err != nil {
		return nil, err
	}

	successfulIds, err := crs.repo.UpdateLifecycleStatusForIDs(idsToUpdate, lifecycleStatus.ID)
	if err != nil {
		return nil, err
	}

	if len(successfulIds) > 0 {
		go func() {
			err := crs.activityLogSvc.LogBulkActivities(
				models.ArchiveAction,
				models.ComplianceRegister{}.TableName(),
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

func (crs *ComplianceRegisterService) Restore(ids []uint, currentUserID uint) ([]uint, error) {
	lifecycleStatus, err := crs.lifecycleStatusSvc.GetByType("active")
	if err != nil {
		return nil, err
	}
	idsToUpdate, err := crs.repo.GetIDsByLifecycleStatus(ids, "archived")
	if err != nil {
		return nil, err
	}

	successfulIds, err := crs.repo.UpdateLifecycleStatusForIDs(idsToUpdate, lifecycleStatus.ID)
	if err != nil {
		return nil, err
	}

	if len(successfulIds) > 0 {
		go func() {
			err := crs.activityLogSvc.LogBulkActivities(
				models.RestoreAction,
				models.ComplianceRegister{}.TableName(),
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

func (crs *ComplianceRegisterService) StatusUpdate(id uint, currentUserID uint, complianceStatusID uint, tx utility.DatabaseTransaction) error {
	existingCompliance, err := crs.GetByID(id)
	if err != nil {
		return err
	}

	existingCompliance.ComplianceStatusID = &complianceStatusID
	existingCompliance.UpdatedBYID = &currentUserID

	err = crs.repo.Update(existingCompliance, tx)
	if err != nil {
		return err
	}

	return nil
}

func (crs *ComplianceRegisterService) ReviewStatusUpdate(complianceRegister *models.ComplianceRegister, tx utility.DatabaseTransaction) error {
	err := crs.repo.Update(complianceRegister, tx)
	if err != nil {
		return err
	}

	return nil
}

func (crs *ComplianceRegisterService) validateComplianceRegisterFormDTO(complianceRegisterFormDTO *complianceDTOs.ComplianceRegisterFormDTO) error {
	if complianceRegisterFormDTO.ComplianceFolderID != nil {
		if _, err := crs.complianceFolderSvc.GetByID(*complianceRegisterFormDTO.ComplianceFolderID); err != nil {
			return utility.NewValidationError(fmt.Errorf("compliance folder does not exist"), "not_exist", "compliance folder")
		}
	}

	// if _, err := crs.complianceCategorySvc.GetByID(*complianceRegisterFormDTO.ComplianceCategoryID); err != nil {
	// 	return utility.NewValidationError(fmt.Errorf("compliance category does not exist"), "not_exist", "compliance category")
	// }

	if complianceRegisterFormDTO.ComplianceStatusID != nil {
		if _, err := crs.complianceStatusSvc.GetByID(*complianceRegisterFormDTO.ComplianceStatusID); err != nil {
			return utility.NewValidationError(fmt.Errorf("compliance status does not exist"), "not_exist", "compliance status")
		}
	}

	// if complianceRegisterFormDTO.ReviewFrequencyID != nil {
	// 	if _, err := crs.reviewFrequencySvc.GetByID(*complianceRegisterFormDTO.ReviewFrequencyID); err != nil {
	// 		return utility.NewValidationError(fmt.Errorf("review frequency does not exist"), "not_exist", "review frequency")
	// 	}
	// }
	// if complianceRegisterFormDTO.SeverityLevelID != nil {
	// 	if _, err := crs.severityLevelSvc.GetByID(*complianceRegisterFormDTO.SeverityLevelID); err != nil {
	// 		return utility.NewValidationError(fmt.Errorf("severity level does not exist"), "not_exist", "severity level")
	// 	}
	// }
	// if complianceRegisterFormDTO.ComplianceAssessmentFrameworkID != nil {
	// 	if _, err := crs.complianceAssessmentFrameworkSvc.GetByID(*complianceRegisterFormDTO.ComplianceAssessmentFrameworkID); err != nil {
	// 		return utility.NewValidationError(fmt.Errorf("compliance assessment framework does not exist"), "not_exist", "compliance assessment framework")
	// 	}
	// }

	// Validate new compliance fields
	if complianceRegisterFormDTO.LegislativeBodyID != nil {
		if _, err := crs.legislativeBodySvc.GetByID(*complianceRegisterFormDTO.LegislativeBodyID); err != nil {
			return utility.NewValidationError(fmt.Errorf("legislative body does not exist"), "not_exist", "legislative body")
		}
	}

	if complianceRegisterFormDTO.LegislationTypeID != nil {
		if _, err := crs.legislationTypeSvc.GetByID(*complianceRegisterFormDTO.LegislationTypeID); err != nil {
			return utility.NewValidationError(fmt.Errorf("legislation type does not exist"), "not_exist", "legislation type")
		}
	}

	if complianceRegisterFormDTO.ImportanceClassificationID != nil {
		if _, err := crs.importanceClassificationSvc.GetByID(*complianceRegisterFormDTO.ImportanceClassificationID); err != nil {
			return utility.NewValidationError(fmt.Errorf("importance classification does not exist"), "not_exist", "importance classification")
		}
	}

	if complianceRegisterFormDTO.SecrecyClassificationID != nil {
		if _, err := crs.secrecyClassificationSvc.GetByID(*complianceRegisterFormDTO.SecrecyClassificationID); err != nil {
			return utility.NewValidationError(fmt.Errorf("secrecy classification does not exist"), "not_exist", "secrecy classification")
		}
	}

	// Validate issuing authorities
	if err := crs.validateIssuingAuthorities(complianceRegisterFormDTO.IssuingAuthorities); err != nil {
		return err
	}

	return nil
}

// validateIssuingAuthorities validates issuing authorities and their entities/documents
func (crs *ComplianceRegisterService) validateIssuingAuthorities(issuingAuthorities []complianceDTOs.ComplianceRegisterIssuingAuthorityFormDTO) error {
	for i, authority := range issuingAuthorities {
		// Validate issuing authority ID is required
		if authority.IssuingAuthorityID == nil {
			return utility.NewValidationError(fmt.Errorf("issuing authority ID is required"), "required", fmt.Sprintf("issuing_authorities[%d].issuing_authority_id", i))
		}

		// TODO: Add issuing authority service validation when available
		// if _, err := crs.issuingAuthoritySvc.GetByID(*authority.IssuingAuthorityID); err != nil {
		// 	return utility.NewValidationError(fmt.Errorf("issuing authority does not exist"), "not_exist", fmt.Sprintf("issuing_authorities[%d].issuing_authority_id", i))
		// }

		// Validate entity if provided
		if authority.Entity != nil {
			if err := crs.validateIssuingAuthorityEntity(authority.Entity, i); err != nil {
				return err
			}
		}

		// Validate document if provided
		if authority.Document != nil {
			if err := crs.validateIssuingAuthorityDocument(authority.Document, i); err != nil {
				return err
			}
		}
	}

	return nil
}

// validateIssuingAuthorityEntity validates issuing authority entity
func (crs *ComplianceRegisterService) validateIssuingAuthorityEntity(entity *complianceDTOs.ComplianceRegisterIssuingAuthorityEntityFormDTO, authorityIndex int) error {
	// Circular number is required if entity is provided
	if entity.CircularNo == nil || *entity.CircularNo == "" {
		return utility.NewValidationError(fmt.Errorf("circular number is required"), "required", fmt.Sprintf("issuing_authorities[%d].entity.circular_no", authorityIndex))
	}

	// Validate date formats if provided
	if entity.DateOfReceiving != nil && *entity.DateOfReceiving != "" {
		if _, err := time.Parse("2006-01-02", *entity.DateOfReceiving); err != nil {
			return utility.NewValidationError(fmt.Errorf("invalid date of receiving format"), "invalid_date", fmt.Sprintf("issuing_authorities[%d].entity.date_of_receiving", authorityIndex))
		}
	}

	// Hijri date is stored as string, no parsing validation needed
	// Frontend should ensure proper format (YYYY-MM-DD)

	return nil
}

// validateIssuingAuthorityDocument validates issuing authority document following the same pattern as validateComplianceDocumentFormDTO
func (crs *ComplianceRegisterService) validateIssuingAuthorityDocument(document *complianceDTOs.ComplianceRegisterIssuingAuthorityDocumentFormDTO, authorityIndex int) error {
	// Document name is required if document is provided
	// Validate document
	if document == nil {
		return utility.NewValidationError(fmt.Errorf("document is required"), "required", "document")
	}

	// Validate file type if provided (following same pattern as main document validation)
	// if document.Type != nil {
	// 	fileType := strings.ToLower(*document.Type)
	// 	if fileType != ".pdf" {
	// 		return utility.NewValidationError(fmt.Errorf("only PDF files are allowed"), "invalid_file_type", fmt.Sprintf("issuing_authorities[%d].document.type", authorityIndex))
	// 	}
	// }

	return nil
}

func (crs *ComplianceRegisterService) parseComplianceRegisterDates(dto *complianceDTOs.ComplianceRegisterFormDTO) (*complianceDTOs.ComplianceRegisterParsedDates, error) {
	parsed := &complianceDTOs.ComplianceRegisterParsedDates{}

	var err error

	if dto.Deadline != nil {
		parsed.Deadline, err = utility.ParseDate(*dto.Deadline)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid deadline"), "invalid_date", "deadline")
		}
	}

	// if dto.IssueDate != nil {
	// 	parsed.IssueDate, err = utility.ParseDate(*dto.IssueDate)
	// 	if err != nil {
	// 		return nil, utility.NewValidationError(fmt.Errorf("invalid issue date"), "invalid_date", "issue date")
	// 	}
	// }

	// if dto.ExpiryDate != nil {
	// 	parsed.ExpiryDate, err = utility.ParseDate(*dto.ExpiryDate)
	// 	if err != nil {
	// 		return nil, utility.NewValidationError(fmt.Errorf("invalid expiry date"), "invalid_date", "expiry date")
	// 	}
	// }

	// if dto.NextReviewDate != nil {
	// 	parsed.NextReviewDate, err = utility.ParseDate(*dto.NextReviewDate)
	// 	if err != nil {
	// 		return nil, utility.NewValidationError(fmt.Errorf("invalid next review date"), "invalid_date", "next review date")
	// 	}
	// }

	// if dto.LastReviewedDate != nil {
	// 	parsed.LastReviewedDate, err = utility.ParseDate(*dto.LastReviewedDate)
	// 	if err != nil {
	// 		return nil, utility.NewValidationError(fmt.Errorf("invalid last reviewed date"), "invalid_date", "last reviewed date")
	// 	}
	// }

	return parsed, nil
}

func (crs *ComplianceRegisterService) UpdateComplianceTotalScore(complianceID uint, totalScore float64, tx utility.DatabaseTransaction) error {
	return crs.repo.UpdateComplianceTotalScore(complianceID, totalScore, tx)
}

// func (crs *ComplianceRegisterService) validateComplianceDocumentFormDTO(document *complianceDTOs.ComplianceDocumentFormDTO) error {
// 	// Validate document
// 	if document == nil {
// 		return utility.NewValidationError(fmt.Errorf("document is required"), "required", "document")
// 	}
// 	// Normalize the file type to lowercase for case-insensitive comparison
// 	fileType := strings.ToLower(document.Type)

// 	// Validate file type (PDF)
// 	if fileType != ".pdf" {
// 		return utility.NewValidationError(fmt.Errorf("only PDF files are allowed"), "invalid_file_type", "type")
// 	}

// 	// @todo
// 	// Validate file size (if needed)
// 	// For example:
// 	// if document.Size > MaxFileSize {
// 	//     return utility.NewValidationError(fmt.Errorf("file size exceeds the limit"), "invalid_file_size", "file_size")
// 	// }

// 	return nil
// }

func (crs *ComplianceRegisterService) GetDocumentPath(id uint, token string) (string, error) {
	// existingComplianceRegister, err := crs.GetFullDetailsByID(id)
	// if err != nil {
	// 	return "", err
	// }

	// if len(existingComplianceRegister.Documents) == 0 {
	// 	return "", fmt.Errorf("no documents found for the compliance register")
	// }

	decryptedFilenameBytes, err := utility.DecryptData(token)
	if err != nil {
		return "", fmt.Errorf("resource not found")
	}

	decryptedFilename := string(decryptedFilenameBytes)
	if decryptedFilename == "" {
		return "", fmt.Errorf("resource not found")
	}

	storagePath := filepath.Join("compliance", "compliance-register", strconv.FormatUint(uint64(id), 10), "documents")

	baseStoragePath := config.GetEnv("STORAGE_PATH", "")
	fullPath := filepath.Join(baseStoragePath, storagePath, decryptedFilename)

	// Check if the file exists
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return "", fmt.Errorf("file not found")
	} else if err != nil {
		return "", fmt.Errorf("error accessing file: %w", err)
	}

	return fullPath, nil
}

func (crs *ComplianceRegisterService) GetIssuingAuthorityDocumentPath(authorityID uint, token string) (string, error) {
	// Get the authority to find the compliance register ID and document
	authority, err := crs.repo.GetIssuingAuthorityByID(authorityID)
	if err != nil {
		return "", fmt.Errorf("authority not found")
	}

	// Get the document for this authority
	documents, err := crs.repo.GetDocumentsByAuthorityID(authorityID)
	if err != nil {
		return "", fmt.Errorf("failed to get documents")
	}

	if len(documents) == 0 {
		return "", fmt.Errorf("no document found for this authority")
	}

	// Use the first document (since there should be only one per authority)
	document := documents[0]

	// Decrypt and validate the token
	decryptedFilenameBytes, err := utility.DecryptData(token)
	if err != nil {
		return "", fmt.Errorf("invalid token")
	}

	decryptedFilename := string(decryptedFilenameBytes)
	if decryptedFilename == "" || decryptedFilename != document.Name {
		return "", fmt.Errorf("invalid token")
	}

	// Build the storage path for issuing authority documents
	storagePath := filepath.Join("compliance", "compliance-register", strconv.FormatUint(uint64(authority.ComplianceRegisterID), 10), "issuing-authority-documents")

	baseStoragePath := config.GetEnv("STORAGE_PATH", "")
	fullPath := filepath.Join(baseStoragePath, storagePath, decryptedFilename)

	// Check if the file exists
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return "", fmt.Errorf("file not found")
	} else if err != nil {
		return "", fmt.Errorf("error accessing file: %w", err)
	}

	return fullPath, nil
}

// OLD EXPORT FUNCTION - COMMENTED OUT
// func (crs *ComplianceRegisterService) ExportComplianceRegisters(data []*complianceDTOs.ComplianceRegisterListDTO, currentUserID uint) (*os.File, error) {
// 	// Retrieve user settings for language preference
// 	// userSetting, err := crs.userSettingSvc.GetUserSetting(currentUserID)
// 	// if err != nil {
// 	// 	return nil, err
// 	// }

// 	// var headers []string
// 	// // Define headers based on the user's language setting
// 	// if userSetting.Language == "ar" {
// 	// 	headers = []string{
// 	// 		"SL No.",
// 	// 		"Reference Number",
// 	// 		"Standard Code",
// 	// 		"Title (AR)",
// 	// 		"Category (AR)",
// 	// 		"Deadline",
// 	// 		"Issue Date",
// 	// 		"Expiry Date",
// 	// 		"Description (AR)",
// 	// 		"Lifecycle Status (AR)",
// 	// 		"Responsible Unit (AR)",
// 	// 		"Compliance Status (AR)",
// 	// 		"Severity Level (AR)",
// 	// 		"Responsible User (AR)",
// 	// 		"Review Frequency (AR)",
// 	// 		"Review Responsible User (AR)",
// 	// 	}
// 	// } else {
// 	// 	headers = []string{
// 	// 		"SL No.",
// 	// 		"Reference Number",
// 	// 		"Standard Code",
// 	// 		"Title (EN)",
// 	// 		"Category (EN)",
// 	// 		"Deadline",
// 	// 		"Issue Date",
// 	// 		"Expiry Date",
// 	// 		"Description (EN)",
// 	// 		"Lifecycle Status (EN)",
// 	// 		"Responsible Unit (EN)",
// 	// 		"Compliance Status (EN)",
// 	// 		"Severity Level (EN)",
// 	// 		"Responsible User (EN)",
// 	// 		"Review Frequency (EN)",
// 	// 		"Review Responsible User (EN)",
// 	// 	}
// 	// }

// 	// // Convert data to [][]interface{} based on language preference
// 	// rows := make([][]interface{}, len(data))
// 	// for i, item := range data {
// 	// 	if userSetting.Language == "ar" {
// 	// 		rows[i] = []interface{}{
// 	// 			i + 1,
// 	// 			item.ReferenceNumber,
// 	// 			item.StandardCode,
// 	// 			item.TitleAr,
// 	// 			item.ComplianceCategoryTitleAr,
// 	// 			item.Deadline,
// 	// 			item.IssueDate,
// 	// 			item.ExpiryDate,
// 	// 			item.DescriptionAr,
// 	// 			item.LifecycleStatusTitleAr,
// 	// 			item.ResponsibleUnitTitleAr,
// 	// 			item.ComplianceStatusTitleAr,
// 	// 			item.SeverityLevelTitleAr,
// 	// 			item.ResponsibleUserFirstNameAr + " " + item.ResponsibleUserLastNameAr,
// 	// 			item.ReviewFrequencyTitleAr,
// 	// 			item.ReviewResponsibleUserFirstNameAr + " " + item.ReviewResponsibleUserLastNameAr,
// 	// 		}
// 	// 	} else {
// 	// 		rows[i] = []interface{}{
// 	// 			i + 1,
// 	// 			item.ReferenceNumber,
// 	// 			item.StandardCode,
// 	// 			item.TitleEn,
// 	// 			item.ComplianceCategoryTitleEn,
// 	// 			item.Deadline,
// 	// 			item.IssueDate,
// 	// 			item.ExpiryDate,
// 	// 			item.DescriptionEn,
// 	// 			item.LifecycleStatusTitleEn,
// 	// 			item.ResponsibleUnitTitleEn,
// 	// 			item.ComplianceStatusTitleEn,
// 	// 			item.SeverityLevelTitleEn,
// 	// 			item.ResponsibleUserFirstNameEn + " " + item.ResponsibleUserLastNameEn,
// 	// 			item.ReviewFrequencyTitleEn,
// 	// 			item.ReviewResponsibleUserFirstNameEn + " " + item.ReviewResponsibleUserLastNameEn,
// 	// 		}
// 	// 	}
// 	// }
// 	// Generate the Excel file
// 	// return utility.GenerateExcelFile(headers, rows)
// 	return nil, nil
// }

func (crs *ComplianceRegisterService) ExportComplianceRegisters(data []*complianceDTOs.ComplianceRegisterListDTO, currentUserID uint) (*os.File, error) {
	userSetting, err := crs.userSettingSvc.GetUserSetting(currentUserID)
	if err != nil {
		return nil, err
	}

	headers := []string{
		"#",
		"Reference Number",
		"Title",
		"Description",
		"Deadline",
		"Legislative Body",
		"Legislation Type",
		"Secrecy Classification",
		"Residual Risk Score",
		"Residual Risk Rating",
		"Importance Classification",
	}

	rows := make([][]interface{}, len(data))
	for i, item := range data {
		var title, description, legislativeBody, legislationType, secrecyClassification, residualRiskRating, importanceClassification string

		if userSetting.Language == "ar" {
			title = item.TitleAr
			description = item.DescriptionAr
			legislativeBody = item.LegislativeBodyTitleAr
			legislationType = item.LegislationTypeTitleAr
			secrecyClassification = item.SecrecyClassificationTitleAr
			importanceClassification = item.ImportanceClassificationTitleAr
			residualRiskRating = utility.GetStringValue(item.AverageResidualRiskRatingTitleAr)
		} else {
			title = item.TitleEn
			description = item.DescriptionEn
			legislativeBody = item.LegislativeBodyTitleEn
			legislationType = item.LegislationTypeTitleEn
			secrecyClassification = item.SecrecyClassificationTitleEn
			importanceClassification = item.ImportanceClassificationTitleEn
			residualRiskRating = utility.GetStringValue(item.AverageResidualRiskRatingTitleEn)
		}

		var residualRiskScore interface{}
		if item.AverageResidualRiskScore != nil {
			residualRiskScore = *item.AverageResidualRiskScore
		} else {
			residualRiskScore = ""
		}

		rows[i] = []interface{}{
			i + 1,
			item.ReferenceNumber,
			title,
			description,
			item.Deadline,
			legislativeBody,
			legislationType,
			secrecyClassification,
			residualRiskScore,
			residualRiskRating,
			importanceClassification,
		}
	}

	return utility.GenerateExcelFile(headers, rows)
}

// func (crs *ComplianceRegisterService) handleDocument(currentUserID uint, complianceRegisterID uint, documentDTO *complianceDTOs.ComplianceDocumentFormDTO, tx ...utility.DatabaseTransaction) error {
// 	if documentDTO == nil {
// 		return nil // Skip if no document is provided
// 	}

// 	// Generate encrypted token for the file name
// 	encryptedFilename, err := utility.EncryptData([]byte(documentDTO.FileName))
// 	if err != nil {
// 		return fmt.Errorf("failed to encrypt document filename: %w", err)
// 	}

// 	// Generate the document URL
// 	url := utility.GenerateDocumentURL("compliance", "compliance-register", strconv.FormatUint(uint64(complianceRegisterID), 10), "documents")

// 	// Create a new ComplianceDocument instance
// 	complianceDocument := &models.ComplianceDocument{
// 		ComplianceRegisterID: complianceRegisterID,
// 		Language:             documentDTO.Language,
// 		URL:                  url,
// 		Name:                 documentDTO.FileName,
// 		Size:                 documentDTO.Size,
// 		Token:                encryptedFilename,
// 		OriginalName:         documentDTO.OriginalName,
// 		Type:                 documentDTO.Type,
// 		CreatedByID:          currentUserID,
// 	}

// 	// Save the document record in the database
// 	if err := crs.repo.CreateDocument(complianceDocument, tx...); err != nil {
// 		return fmt.Errorf("failed to save document record: %w", err)
// 	}

// 	// Store the document file in the storage
// 	storagePath := filepath.Join("compliance", "compliance-register", strconv.FormatUint(uint64(complianceRegisterID), 10), "documents")
// 	if err := utility.StoreFile(documentDTO.FileName, storagePath); err != nil {
// 		return fmt.Errorf("failed to store document file: %w", err)
// 	}

// 	return nil
// }

func (crs *ComplianceRegisterService) CreateComplianceVersion(currentUserID uint, complianceRegisterFormDTO *complianceDTOs.ComplianceVersionFormDTO) (uint, error) {
	parentCompliance, err := crs.GetByID(*complianceRegisterFormDTO.ParentComplianceRegisterID)
	if err != nil {
		return 0, err
	}

	tx := crs.repo.BeginTransaction()

	referenceNumber, err := crs.repo.GenerateReferenceNumber()
	if err != nil {
		return 0, err
	}

	// lifecycleStatus, err := crs.lifecycleStatusSvc.GetByType("active")
	// if err != nil {
	// 	return 0, err
	// }

	// if err := crs.validateComplianceDocumentFormDTO(&complianceRegisterFormDTO.DocumentEn); err != nil {
	// 	return 0, err
	// }

	// if err := crs.validateComplianceDocumentFormDTO(&complianceRegisterFormDTO.DocumentAr); err != nil {
	// 	return 0, err
	// }

	// responsibleUnitID, err := crs.userSvc.GetUserPrimaryUnitID(uint(*parentCompliance.ResponsibleUserID))
	// if err != nil {
	// 	return 0, utility.NewValidationError(fmt.Errorf("organization unit not set for responsible user"), "not_exist", "bcp communication plan responsible unit")
	// }

	complianceRegister := &models.ComplianceRegister{
		// ComplianceCategoryID:            parentCompliance.ComplianceCategoryID,
		ComplianceStatusID: parentCompliance.ComplianceStatusID,
		// ResponsibleUnitID:               &responsibleUnitID,
		// ResponsibleUserID:               parentCompliance.ResponsibleUserID,
		// ReviewFrequencyID:               parentCompliance.ReviewFrequencyID,
		// LastReviewNote:                  parentCompliance.LastReviewNote,
		// ReviewResponsibleUserID:         parentCompliance.ReviewResponsibleUserID,
		// SeverityLevelID:                 parentCompliance.SeverityLevelID,
		ReferenceNumber: referenceNumber,
		// StandardCode:                    parentCompliance.StandardCode,
		ParentComplianceRegisterID: complianceRegisterFormDTO.ParentComplianceRegisterID,
		TitleEn:                    parentCompliance.TitleEn,
		TitleAr:                    parentCompliance.TitleAr,
		DescriptionEn:              complianceRegisterFormDTO.DescriptionEn,
		DescriptionAr:              complianceRegisterFormDTO.DescriptionAr,
		Version:                    complianceRegisterFormDTO.Version,
		Deadline:                   parentCompliance.Deadline,
		// IssueDate:                       parentCompliance.IssueDate,
		// ExpiryDate:                      parentCompliance.ExpiryDate,
		// NextReviewDate:                  parentCompliance.NextReviewDate,
		// LastReviewedDate:                parentCompliance.LastReviewedDate,
		// LastReviewedUserID:              parentCompliance.LastReviewedUserID,
		// ComplianceReviewStatusID:        parentCompliance.ComplianceReviewStatusID,
		CreatedBYID: currentUserID,
		// LifecycleStatusID:               &lifecycleStatus.ID,
		ComplianceFolderID: parentCompliance.ComplianceFolderID,
		// ComplianceAssessmentFrameworkID: parentCompliance.ComplianceAssessmentFrameworkID,
		// IsControlAssessment:             parentCompliance.IsControlAssessment,
		// IsRequirementAssessment:         parentCompliance.IsRequirementAssessment,
		IsLatest: true,

		// New Compliance Fields from parent
		LegislativeBodyID: parentCompliance.LegislativeBodyID,
		LegislationTypeID: parentCompliance.LegislationTypeID,
		LegislationNo:     parentCompliance.LegislationNo,
		//ComplianceDeadline:                 parentCompliance.ComplianceDeadline,
		ImportanceClassificationID: parentCompliance.ImportanceClassificationID,
		SecrecyClassificationID:    parentCompliance.SecrecyClassificationID,
		Subject:                    parentCompliance.Subject,
		IsApplicable:               parentCompliance.IsApplicable,
		IsNotApplicable:            parentCompliance.IsNotApplicable,
		// IsRegulationNotificationAssessment: parentCompliance.IsRegulationNotificationAssessment,
		// IsRiskBasedAssessment:              parentCompliance.IsRiskBasedAssessment,
		IsPublic: parentCompliance.IsPublic,
		IsShared: parentCompliance.IsShared,
	}

	parentComplianceRegister, err := crs.GetFullDetailsByID(*complianceRegisterFormDTO.ParentComplianceRegisterID)
	if err != nil {
		return 0, err
	}

	for _, version := range parentComplianceRegister.ChildComplianceRegister {
		versionDetails, err := crs.GetByID(uint(version.ID))
		if err != nil {
			return 0, err
		}

		versionDetails.IsLatest = false

		err = crs.repo.Update(versionDetails, tx)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	parentCompliance.IsLatest = false

	err = crs.repo.Update(parentCompliance, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	// complianceRegister.StandardCode = &parentComplianceRegister.StandardCode

	id, err := crs.repo.Create(complianceRegister, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	if parentCompliance.ComplianceStatusID != nil {
		complianceStatusID := parentCompliance.ComplianceStatusID
		complianceStatusUpdate := &models.ComplianceStatusUpdate{
			ComplianceRegisterID: id,
			ComplianceStatusID:   *complianceStatusID,
			Comment:              "Initial compliance status",
			CreatedByID:          currentUserID,
		}

		_, err = crs.repo.CreateInitialComplianceStatusUpdate(complianceStatusUpdate, tx)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	// if err := crs.handleDocument(currentUserID, id, &complianceRegisterFormDTO.DocumentEn, tx); err != nil {
	// 	tx.Rollback()
	// 	return 0, err
	// }

	// if err := crs.handleDocument(currentUserID, id, &complianceRegisterFormDTO.DocumentAr, tx); err != nil {
	// 	tx.Rollback()
	// 	return 0, err
	// }

	// Copy issuing authorities from parent compliance register
	if len(parentComplianceRegister.IssuingAuthorities) > 0 {
		parentIssuingAuthorities := make([]complianceDTOs.ComplianceRegisterIssuingAuthorityFormDTO, 0, len(parentComplianceRegister.IssuingAuthorities))
		for _, parentAuthority := range parentComplianceRegister.IssuingAuthorities {
			var issuingAuthorityID *uint
			if parentAuthority.IssuingAuthority != nil {
				issuingAuthorityID = &parentAuthority.IssuingAuthority.ID
			}

			authorityDTO := complianceDTOs.ComplianceRegisterIssuingAuthorityFormDTO{
				IssuingAuthorityID: issuingAuthorityID,
			}

			// Copy entity if exists
			if parentAuthority.Entity != nil {
				authorityDTO.Entity = &complianceDTOs.ComplianceRegisterIssuingAuthorityEntityFormDTO{
					CircularNo:           &parentAuthority.Entity.CircularNo,
					DateOfReceiving:      &parentAuthority.Entity.DateOfReceivingString,
					HijriDateOfReceiving: &parentAuthority.Entity.HijriDateOfReceivingString,
				}
			}

			// Copy document if exists
			if parentAuthority.Document != nil {
				authorityDTO.Document = &complianceDTOs.ComplianceRegisterIssuingAuthorityDocumentFormDTO{
					FileName:     &parentAuthority.Document.Name,
					URL:          &parentAuthority.Document.URL,
					OriginalName: &parentAuthority.Document.OriginalName,
					Type:         &parentAuthority.Document.Type,
					Size:         &parentAuthority.Document.Size,
				}
			}

			parentIssuingAuthorities = append(parentIssuingAuthorities, authorityDTO)
		}

		// Handle issuing authorities for the new version
		if err := crs.handleIssuingAuthorities(currentUserID, id, parentIssuingAuthorities, parentCompliance.SecrecyClassificationID, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	// Copy shared units from parent compliance register if it's shared
	if parentCompliance.IsShared {
		existingSharedUnits, err := crs.repo.GetSharedUnitsByComplianceRegisterID(*complianceRegisterFormDTO.ParentComplianceRegisterID, tx)
		if err != nil {
			tx.Rollback()
			return 0, err
		}

		if len(existingSharedUnits) > 0 {
			sharedUnitIDs := make([]*uint, 0, len(existingSharedUnits))
			for _, unit := range existingSharedUnits {
				sharedUnitIDs = append(sharedUnitIDs, &unit.DepartmentID)
			}

			// Handle shared units for the new version
			if err := crs.handleSharedUnits(currentUserID, id, sharedUnitIDs, parentCompliance.IsShared, tx); err != nil {
				tx.Rollback()
				return 0, err
			}
		}
	}

	tx.Commit()

	go func() {
		existingCompliance, err := crs.GetFullDetailsByID(complianceRegister.ID)
		if err != nil {
			utility.LogError("Failed to fetch compliance register", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		err = crs.activityLogSvc.LogActivity(
			models.CreateAction,
			models.ComplianceRegister{}.TableName(),
			id,
			existingCompliance,
			currentUserID,
			"compliance", "compliance_register", "create compliance", nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.ComplianceRegister{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}

		// Prepare notification and email
		// urlParam := fmt.Sprintf(`{"id":%d}`, id)
		// entity := "compliance-register"
		// action := "responsible-user-assign"

		// dbNotificationDTO := &commonDTOs.DBNotificationDTO{
		// 	UserIDs:              []uint{*parentCompliance.ResponsibleUserID},
		// 	ModuleIdentifier:     "compliance",
		// 	NotificationType:     string(commonDTOs.AssignmentType),
		// 	HeadingEn:            "Compliance Register Responsibility Assigned",
		// 	HeadingAr:            "تم تعيين مسؤولية سجل الامتثال",
		// 	TitleEn:              parentCompliance.TitleEn,
		// 	TitleAr:              parentCompliance.TitleAr,
		// 	ReferenceNumber:      complianceRegister.ReferenceNumber,
		// 	Entity:               &entity,
		// 	Action:               &action,
		// 	URLParam:             urlParam,
		// 	IsSystemNotification: false,
		// 	CreatedByID:          &currentUserID,
		// }

		// responsibleUser, err := crs.userSvc.GetByID(*parentCompliance.ResponsibleUserID)
		// if err != nil {
		// 	utility.LogError("Failed to fetch responsible user", map[string]interface{}{
		// 		"error": err.Error(),
		// 	})
		// 	return
		// }

		// createdBy, err := crs.userSvc.GetByID(currentUserID)
		// if err != nil {
		// 	utility.LogError("Failed to fetch responsible user", map[string]interface{}{
		// 		"error": err.Error(),
		// 	})
		// 	return
		// }

		// // Assuming GetByID returns only an error and no severity level
		// severityLevel, err := crs.severityLevelSvc.GetByID(*complianceRegister.SeverityLevelID)
		// if err != nil {
		// 	utility.LogError("Failed to fetch severity level", map[string]interface{}{
		// 		"error":   err.Error(),
		// 		"levelID": complianceRegister.SeverityLevelID,
		// 	})
		// 	return
		// }

		// emailTo := []string{responsibleUser.Email}
		// subject := fmt.Sprintf("Compliance Register/%s: Responsibility Assigned", complianceRegister.TitleEn)
		// appURL := utility.GetAppURL()

		// emailData := map[string]interface{}{
		// 	"FirstNameEn":               responsibleUser.FirstNameEn,
		// 	"LastNameEn":                responsibleUser.LastNameEn,
		// 	"ComplianceReferenceNumber": complianceRegister.ReferenceNumber,
		// 	"ComplianceTitleEn":         complianceRegister.TitleEn,
		// 	"ComplianceTitleAr":         complianceRegister.TitleAr,
		// 	"ComplianceDescription":     complianceRegister.DescriptionEn,
		// 	"SeverityLevel":             severityLevel.TitleEn,
		// 	"CreatedByFirstNameEn":      createdBy.FirstNameEn,
		// 	"CreatedByLastNameEn":       createdBy.LastNameEn,
		// 	"AssignedOn":                time.Now().Format("2006-01-02"),
		// 	"RedirectURL":               fmt.Sprintf("%s/compliance/register/%d", appURL, id),
		// }

		// emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
		// 	EmailTo:              emailTo,
		// 	Subject:              subject,
		// 	TemplateRelativePath: "compliance/templates/email/compliance_register_create_responsible_user_mail.html",
		// 	EmailData:            emailData,
		// }

		// err = crs.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "compliance_register")
		// if err != nil {
		// 	utility.LogError("Failed to send notifications", map[string]interface{}{
		// 		"error": err.Error(),
		// 	})
		// }
	}()

	return id, nil
}

func (crs ComplianceRegisterService) GetAllVersions(id uint) ([]*complianceDTOs.ComplianceRegisterDetailDTO, error) {
	var complianceDTO []*complianceDTOs.ComplianceRegisterDetailDTO
	filteredCompliance := make([]*models.ComplianceRegister, 0)

	existingComplianceRegister, err := crs.GetByID(id)
	if err != nil {
		return nil, err
	}

	complianceRegisterID := id
	if existingComplianceRegister.ParentComplianceRegisterID != nil {
		complianceRegisterID = *existingComplianceRegister.ParentComplianceRegisterID
		parentRegister, err := crs.repo.GetFullDetailsByID(complianceRegisterID)
		if err != nil {
			return nil, err
		}
		filteredCompliance = append(filteredCompliance, parentRegister)
	}

	compliance, err := crs.repo.GetAllVersions(complianceRegisterID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	for _, c := range compliance {
		if c.ID != id { // exclude the current one
			filteredCompliance = append(filteredCompliance, c)
		}
	}

	err = mapstructure.Decode(filteredCompliance, &complianceDTO)
	if err != nil {
		return nil, err
	}

	return complianceDTO, nil
}

func (crs *ComplianceRegisterService) ForceDelete(ids []uint, currentUserID uint) ([]uint, error) {

	tx := crs.repo.BeginTransaction()
	var failedIds []uint
	var deletedData []map[string]interface{}

	for _, id := range ids {
		if err := crs.repo.ForceDelete(id, tx); err != nil {
			failedIds = append(failedIds, id)
			continue
		}

		deletedData = append(deletedData, map[string]interface{}{
			"id": id,
		})
	}

	if len(failedIds) > 0 {
		tx.Rollback()
		return failedIds, nil
	}

	tx.Commit()

	go func() {
		err := crs.activityLogSvc.LogBulkActivities(
			models.DeleteAction,
			models.ComplianceRegister{}.TableName(),
			ids,
			deletedData,
			currentUserID,
			nil,
		)
		if err != nil {
			// @todo - save to log file
			fmt.Printf("Failed to log bulk delete activity: %v\n", err)
		}
	}()

	return nil, nil
}

// handleIssuingAuthorities handles the creation and update of issuing authorities for a compliance register
func (crs *ComplianceRegisterService) handleIssuingAuthorities(currentUserID uint, complianceRegisterID uint, issuingAuthorities []complianceDTOs.ComplianceRegisterIssuingAuthorityFormDTO, secrecyClassificationID *uint, tx utility.DatabaseTransaction) error {
	// Get existing issuing authorities for this compliance register
	existingAuthorities, err := crs.repo.GetIssuingAuthoritiesByComplianceRegisterID(complianceRegisterID, tx)
	if err != nil {
		return err
	}

	// Create maps to track existing records
	existingAuthorityIDs := make(map[uint]uint) // map[authorityID]recordID
	for _, authority := range existingAuthorities {
		if authority.ID != 0 {
			existingAuthorityIDs[authority.ID] = authority.ID
		}
	}

	// Get new authority IDs from form - also track by issuing_authority_id for authorities without ID
	newAuthorityIDs := make(map[uint]bool)
	newIssuingAuthorityIDs := make(map[uint]bool) // Track by issuing_authority_id
	for _, authorityDTO := range issuingAuthorities {
		if authorityDTO.ID != nil {
			newAuthorityIDs[*authorityDTO.ID] = true
		}
		if authorityDTO.IssuingAuthorityID != nil {
			newIssuingAuthorityIDs[*authorityDTO.IssuingAuthorityID] = true
		}
	}

	// Process each authority from the request
	for _, authorityDTO := range issuingAuthorities {
		if authorityDTO.IssuingAuthorityID == nil {
			continue
		}

		authority := &models.ComplianceRegisterIssuingAuthority{
			ComplianceRegisterID: complianceRegisterID,
			IssuingAuthorityID:   *authorityDTO.IssuingAuthorityID,
			CreatedBYID:          currentUserID,
		}

		// Create or update the issuing authority
		var authorityID uint
		if authorityDTO.ID != nil {
			// Update existing by ID
			authority.ID = *authorityDTO.ID
			// Find existing authority to preserve created_at and created_by_id
			for _, existing := range existingAuthorities {
				if existing.ID == *authorityDTO.ID {
					authority.CreatedAt = existing.CreatedAt     // Preserve original created_at
					authority.CreatedBYID = existing.CreatedBYID // Preserve original created_by_id
					break
				}
			}
			err := crs.repo.UpdateIssuingAuthority(authority, tx)
			if err != nil {
				return err
			}
			authorityID = *authorityDTO.ID
		} else {
			// Check if authority already exists by issuing_authority_id for this compliance register
			var existingAuthority *models.ComplianceRegisterIssuingAuthority
			for _, existing := range existingAuthorities {
				if existing.IssuingAuthorityID == *authorityDTO.IssuingAuthorityID {
					existingAuthority = existing
					break
				}
			}

			if existingAuthority != nil {
				// Update existing authority found by issuing_authority_id
				authority.ID = existingAuthority.ID
				authority.CreatedAt = existingAuthority.CreatedAt     // Preserve original created_at
				authority.CreatedBYID = existingAuthority.CreatedBYID // Preserve original created_by_id
				err := crs.repo.UpdateIssuingAuthority(authority, tx)
				if err != nil {
					return err
				}
				authorityID = existingAuthority.ID
			} else {
				// Create new authority
				id, err := crs.repo.CreateIssuingAuthority(authority, tx)
				if err != nil {
					return err
				}
				authorityID = id
			}
		}

		// Handle single entity per authority
		if authorityDTO.Entity != nil {
			if authorityDTO.Entity.CircularNo != nil && *authorityDTO.Entity.CircularNo != "" {
				var dateOfReceiving time.Time
				var hijriDateOfReceiving string

				if authorityDTO.Entity.DateOfReceiving != nil && *authorityDTO.Entity.DateOfReceiving != "" {
					parsedDate, err := time.Parse("2006-01-02", *authorityDTO.Entity.DateOfReceiving)
					if err != nil {
						return utility.NewValidationError(fmt.Errorf("invalid date_of_receiving format"), "invalid_date", "date_of_receiving")
					}
					// Convert to UTC to avoid timezone offset issues with SQL Server
					dateOfReceiving = time.Date(parsedDate.Year(), parsedDate.Month(), parsedDate.Day(), 0, 0, 0, 0, time.UTC)
				}

				if authorityDTO.Entity.HijriDateOfReceiving != nil && *authorityDTO.Entity.HijriDateOfReceiving != "" {
					// Store hijri date as string without parsing
					hijriDateOfReceiving = *authorityDTO.Entity.HijriDateOfReceiving
				}

				if authorityDTO.Entity.ID != nil {
					// Update existing entity - fetch first to preserve CreatedAt and other fields
					existingEntity, err := crs.repo.GetIssuingAuthorityEntityByID(*authorityDTO.Entity.ID, tx)
					if err != nil {
						return fmt.Errorf("failed to get existing entity: %w", err)
					}

					// Update only the fields that can be changed
					existingEntity.CircularNo = *authorityDTO.Entity.CircularNo
					existingEntity.DateOfReceiving = dateOfReceiving
					existingEntity.HijriDateOfReceiving = hijriDateOfReceiving
					// Set updated_at and updated_by fields
					now := time.Now()
					existingEntity.UpdatedAt = &now
					existingEntity.UpdatedBYID = &currentUserID

					err = crs.repo.UpdateIssuingAuthorityEntity(existingEntity, tx)
					if err != nil {
						return err
					}
				} else {
					// Create new entity
					entity := &models.ComplianceRegisterIssuingAuthorityEntity{
						ComplianceRegisterIssuingAuthorityID: authorityID,
						CircularNo:                           *authorityDTO.Entity.CircularNo,
						DateOfReceiving:                      dateOfReceiving,
						HijriDateOfReceiving:                 hijriDateOfReceiving,
						CreatedBYID:                          currentUserID,
						CreatedAt:                            time.Now(),
					}

					_, err := crs.repo.CreateIssuingAuthorityEntity(entity, tx)
					if err != nil {
						return err
					}
				}
			}
		}

		// Handle single document per authority - ONLY process when explicit flags are provided
		if authorityDTO.Document != nil &&
			((authorityDTO.Document.IsDeletedFile != nil && *authorityDTO.Document.IsDeletedFile) ||
				(authorityDTO.Document.IsNewFile != nil && *authorityDTO.Document.IsNewFile)) {

			// Handle document deletion - ID is required
			if authorityDTO.Document.IsDeletedFile != nil && *authorityDTO.Document.IsDeletedFile {
				if authorityDTO.Document.ID == nil {
					return fmt.Errorf("document ID is required when is_deleted_file is true")
				}

				// Get existing document to retrieve file name for deletion
				existingDoc, err := crs.repo.GetIssuingAuthorityDocumentByID(*authorityDTO.Document.ID, tx)
				if err != nil {
					return fmt.Errorf("failed to get existing document for deletion: %w", err)
				}

				// Delete the document record from database
				err = crs.repo.DeleteIssuingAuthorityDocument(*authorityDTO.Document.ID, tx)
				if err != nil {
					return fmt.Errorf("failed to delete document record: %w", err)
				}

				// Remove the physical file
				storagePath := filepath.Join("compliance", "compliance-register", strconv.FormatUint(uint64(complianceRegisterID), 10), "issuing-authority-documents")
				if err := utility.RemoveFile(existingDoc.Name, storagePath); err != nil {
					return fmt.Errorf("failed to remove document file: %w", err)
				}
			} else if authorityDTO.Document.IsNewFile != nil && *authorityDTO.Document.IsNewFile {
				// Handle new document creation - since only single document is allowed per authority,
				// delete existing document first if it exists, then create new one
				if authorityDTO.Document.FileName != nil && *authorityDTO.Document.FileName != "" {
					// Check secrecy classification to prevent document upload for Secret/Top Secret
					if secrecyClassificationID != nil {
						secrecyClassification, err := crs.secrecyClassificationSvc.GetByID(*secrecyClassificationID)
						if err != nil {
							return fmt.Errorf("failed to get secrecy classification: %w", err)
						}

						// CR-4: Disable document upload for Secret or Top Secret classifications
						if secrecyClassification.Type == "secret" || secrecyClassification.Type == "top_secret" {
							return utility.NewValidationError(
								fmt.Errorf("document upload is disabled for Secret and Top Secret classifications"),
								"document_upload_disabled",
								"secrecy_classification",
							)
						}
					}

					// Check if there's an existing document for this authority and delete it first
					// Since only single document is allowed, we need to clean up the old one
					existingAuthorities, err := crs.repo.GetIssuingAuthoritiesByComplianceRegisterID(complianceRegisterID, tx)
					if err != nil {
						return fmt.Errorf("failed to get existing authorities: %w", err)
					}

					// Find the current authority and check if it has an existing document
					for _, existingAuth := range existingAuthorities {
						if existingAuth.ID == authorityID && existingAuth.Document != nil {
							// Delete existing document record
							err = crs.repo.DeleteIssuingAuthorityDocument(existingAuth.Document.ID, tx)
							if err != nil {
								return fmt.Errorf("failed to delete existing document: %w", err)
							}

							// Remove existing physical file
							storagePath := filepath.Join("compliance", "compliance-register", strconv.FormatUint(uint64(complianceRegisterID), 10), "issuing-authority-documents")
							if err := utility.RemoveFile(existingAuth.Document.Name, storagePath); err != nil {
								return fmt.Errorf("failed to remove existing document file: %w", err)
							}
							break
						}
					}

					// Generate encrypted token for the new file
					encryptedFilename, err := utility.EncryptData([]byte(*authorityDTO.Document.FileName))
					if err != nil {
						return fmt.Errorf("failed to encrypt document filename: %w", err)
					}

					// Generate the document URL
					url := utility.GenerateDocumentURL("compliance", "compliance-register", strconv.FormatUint(uint64(complianceRegisterID), 10), "issuing-authority-documents")

					// Create new document
					document := &models.ComplianceRegisterIssuingAuthorityDocument{
						ComplianceRegisterIssuingAuthorityID: authorityID,
						URL:                                  url,
						Name:                                 *authorityDTO.Document.FileName,
						Token:                                encryptedFilename,
						OriginalName:                         *authorityDTO.Document.OriginalName,
						Type:                                 *authorityDTO.Document.Type,
						CreatedByID:                          currentUserID,
						UpdatedByID:                          currentUserID,
					}

					if authorityDTO.Document.Size != nil {
						document.Size = *authorityDTO.Document.Size
					}

					// Create the new document
					_, err = crs.repo.CreateIssuingAuthorityDocument(document, tx)
					if err != nil {
						return err
					}

					// Store the new document file
					storagePath := filepath.Join("compliance", "compliance-register", strconv.FormatUint(uint64(complianceRegisterID), 10), "issuing-authority-documents")
					if err := utility.StoreFile(*authorityDTO.Document.FileName, storagePath); err != nil {
						return fmt.Errorf("failed to store issuing authority document file: %w", err)
					}
				}
			}
			// Note: When no document flags are passed, we intentionally do nothing
			// This preserves existing documents without any modifications
		}
	}

	// Delete removed issuing authorities and their related data
	var authorityIDsToDelete []uint
	for authorityID := range existingAuthorityIDs {
		// Check if authority exists in new request by ID or by issuing_authority_id
		existsByID := newAuthorityIDs[authorityID]

		// Also check if this authority's issuing_authority_id is still present in the request
		var existsByIssuingAuthorityID bool
		for _, existingAuth := range existingAuthorities {
			if existingAuth.ID == authorityID && newIssuingAuthorityIDs[existingAuth.IssuingAuthorityID] {
				existsByIssuingAuthorityID = true
				break
			}
		}

		// Only delete if authority doesn't exist by either ID or issuing_authority_id
		if !existsByID && !existsByIssuingAuthorityID {
			authorityIDsToDelete = append(authorityIDsToDelete, authorityID)
		}
	}

	// Delete removed authorities and their related entities/documents
	for _, authorityID := range authorityIDsToDelete {
		// First, get all documents for this authority to remove physical files
		documents, err := crs.repo.GetDocumentsByAuthorityID(authorityID, tx)
		if err != nil {
			return fmt.Errorf("failed to get documents for authority deletion: %w", err)
		}

		// Remove all physical files for this authority
		storagePath := filepath.Join("compliance", "compliance-register", strconv.FormatUint(uint64(complianceRegisterID), 10), "issuing-authority-documents")
		for _, document := range documents {
			if err := utility.RemoveFile(document.Name, storagePath); err != nil {
				// Log error but don't fail the transaction - file might already be deleted
				fmt.Printf("Warning: failed to remove document file %s: %v\n", document.Name, err)
			}
		}

		// Delete all documents for this authority from database
		if err := crs.repo.DeleteIssuingAuthorityDocumentsByAuthorityID(authorityID, tx); err != nil {
			return fmt.Errorf("failed to delete issuing authority documents: %w", err)
		}

		// Delete all entities for this authority
		if err := crs.repo.DeleteIssuingAuthorityEntitiesByAuthorityID(authorityID, tx); err != nil {
			return fmt.Errorf("failed to delete issuing authority entities: %w", err)
		}

		// Delete the authority itself
		if err := crs.repo.DeleteIssuingAuthority(authorityID, tx); err != nil {
			return fmt.Errorf("failed to delete issuing authority: %w", err)
		}
	}

	return nil
}

// handleSharedUnits handles the creation and update of shared units for a compliance register
func (crs *ComplianceRegisterService) handleSharedUnits(currentUserID uint, complianceRegisterID uint, sharedUnitIDs []*uint, isShared bool, tx utility.DatabaseTransaction) error {
	// If not shared or no shared unit IDs, skip processing
	if !isShared || len(sharedUnitIDs) == 0 {
		return nil
	}

	// Get existing shared units for this compliance register
	existingSharedUnits, err := crs.repo.GetSharedUnitsByComplianceRegisterID(complianceRegisterID, tx)
	if err != nil {
		return fmt.Errorf("failed to get existing shared units: %w", err)
	}

	// Create a map of existing department IDs for quick lookup
	existingDeptIDs := make(map[uint]bool)
	for _, unit := range existingSharedUnits {
		existingDeptIDs[unit.DepartmentID] = true
	}

	// Prepare new shared units to create
	var newSharedUnits []*models.ComplianceRegisterSharedUnit
	var deptIDsToKeep []uint

	for _, unitID := range sharedUnitIDs {
		if unitID == nil {
			continue
		}

		deptIDsToKeep = append(deptIDsToKeep, *unitID)

		// Only create if it doesn't already exist
		if !existingDeptIDs[*unitID] {
			sharedUnit := &models.ComplianceRegisterSharedUnit{
				ComplianceRegisterID: complianceRegisterID,
				DepartmentID:         *unitID,
				CreatedBYID:          currentUserID,
				CreatedAt:            time.Now(),
			}
			newSharedUnits = append(newSharedUnits, sharedUnit)
		}
	}

	// Create new shared units
	if len(newSharedUnits) > 0 {
		if err := crs.repo.CreateSharedUnits(newSharedUnits, tx); err != nil {
			return fmt.Errorf("failed to create shared units: %w", err)
		}
	}

	// Delete shared units that are no longer needed
	var deptIDsToDelete []uint
	keepMap := make(map[uint]bool)
	for _, id := range deptIDsToKeep {
		keepMap[id] = true
	}

	for _, unit := range existingSharedUnits {
		if !keepMap[unit.DepartmentID] {
			deptIDsToDelete = append(deptIDsToDelete, unit.DepartmentID)
		}
	}

	if len(deptIDsToDelete) > 0 {
		if err := crs.repo.DeleteSharedUnits(complianceRegisterID, deptIDsToDelete, tx); err != nil {
			return fmt.Errorf("failed to delete shared units: %w", err)
		}
	}

	return nil
}

func (crs *ComplianceRegisterService) UpdateNote(id uint, currentUserID uint, dto *complianceDTOs.ComplianceRegisterNoteDTO) error {
	// Get existing investigation for audit log
	existing, err := crs.repo.GetByID(id)
	if err != nil {
		return err
	}
	tx := crs.repo.BeginTransaction()
	// Update the note
	if err := crs.repo.UpdateNote(id, dto.Note, tx); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	// Log the activity
	go func() {
		changedData := map[string]interface{}{
			"note_and_observations": map[string]interface{}{
				"old": existing.Note,
				"new": dto.Note,
			},
		}

		_ = crs.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.ComplianceRegister{}.TableName(),
			id,
			changedData,
			currentUserID,
			"compliance", "compliance_register", "update compliance register note", nil,
		)
	}()

	return nil
}

func (crs *ComplianceRegisterService) GenerateComplianceTemplate() (*os.File, error) {
	templateHeaders := []string{
		"LegislativeBody", "LegislationType", "LegislationNo", "ResponsibleUser",
		"Title", "Description", "Version", "Deadline",
		"ImportanceClassification", "SecrecyClassification", "ComplianceFolder", "IssuingAuthority",
		"IssuingAuthorityCircularNo", "IssuingAuthorityDateOfReceiving", "IssuingAuthorityHijriDateOfReceiving",
		"ApplicabilityStatement", "AssessSettings", "SharedUnit",
	}

	templateData := [][]interface{}{}

	// Fetch Legislative Bodies
	legislativeBodies, err := crs.legislativeBodySvc.GetAll("", "", "", "active")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch legislative bodies: %v", err)
	}
	legislativeBodyData := [][]interface{}{{"LegislativeBody"}}
	for _, body := range legislativeBodies {
		legislativeBodyData = append(legislativeBodyData, []interface{}{body.TitleEn})
	}

	// Fetch Legislation Types
	legislationTypes, err := crs.legislationTypeSvc.GetAll("", "", "", "active")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch legislation types: %v", err)
	}
	legislationTypeData := [][]interface{}{{"LegislationType"}}
	for _, lType := range legislationTypes {
		legislationTypeData = append(legislationTypeData, []interface{}{lType.TitleEn})
	}

	// Fetch Users
	users, err := crs.userSvc.GetAllUsers()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch users: %v", err)
	}
	userData := [][]interface{}{{"ResponsibleUser", "Email"}}
	for _, u := range users {
		userData = append(userData, []interface{}{
			fmt.Sprintf("%s %s", u.FirstNameEN, u.LastNameEN),
			u.Email,
		})
	}

	// Fetch Importance Classifications
	importanceClassifications, err := crs.importanceClassificationSvc.GetAll("", "", "", "active")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch importance classifications: %v", err)
	}
	importanceClassificationData := [][]interface{}{{"ImportanceClassification"}}
	for _, ic := range importanceClassifications {
		importanceClassificationData = append(importanceClassificationData, []interface{}{ic.TitleEn})
	}

	// Fetch Secrecy Classifications
	secrecyClassifications, err := crs.secrecyClassificationSvc.GetAll("", "", "", "active")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch secrecy classifications: %v", err)
	}
	secrecyClassificationData := [][]interface{}{{"SecrecyClassification"}}
	for _, sc := range secrecyClassifications {
		secrecyClassificationData = append(secrecyClassificationData, []interface{}{sc.TitleEn})
	}

	// Fetch Issuing Authorities
	issuingAuthorities, err := crs.issuingAuthoritySvc.GetAll("", "", "", "active")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch issuing authorities: %v", err)
	}
	issuingAuthorityData := [][]interface{}{{"IssuingAuthority"}}
	for _, ia := range issuingAuthorities {
		issuingAuthorityData = append(issuingAuthorityData, []interface{}{ia.TitleEn})
	}

	// Fetch Compliance Folders
	complianceFolders, err := crs.complianceFolderSvc.GetAll("", "", "", "active")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch compliance folders: %v", err)
	}
	complianceFolderData := [][]interface{}{{"ComplianceFolder"}}
	for _, folder := range complianceFolders {
		complianceFolderData = append(complianceFolderData, []interface{}{folder.TitleEn})
	}

	// Fetch Shared Units (Departments - Level 2)
	sharedUnits, err := crs.organizationUnitSvc.GetAll(&settingsDTOs.OrganizationUnitListQueryParamsDTO{
		Status: "active",
		Level:  2,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch shared units: %v", err)
	}
	sharedUnitData := [][]interface{}{{"SharedUnit"}}
	for _, unit := range sharedUnits {
		sharedUnitData = append(sharedUnitData, []interface{}{unit.TitleEn})
	}

	// Configure which columns should have dropdowns
	dropDownConfig := map[string]bool{
		"LegislativeBody":          true,
		"LegislationType":          true,
		"ResponsibleUser":          true,
		"ImportanceClassification": true,
		"SecrecyClassification":    true,
		"ComplianceFolder":         true,
		"IssuingAuthority":         true,
		"ApplicabilityStatement":   true,
		"AssessSettings":           true,
		"SharedUnit":               true,
	}

	// Add dropdown data for grouped fields
	applicabilityData := [][]interface{}{{"ApplicabilityStatement"}, {"Applicable"}, {"Not Applicable"}}
	assessSettingsData := [][]interface{}{{"AssessSettings"}, {"Public"}, {"Shared"}}

	file, err := utility.GenerateExcelFileWithSheetsAndDropDown(map[string][][]interface{}{
		"Template":                 append([][]interface{}{utility.ConvertToInterfaceSlice(templateHeaders)}, templateData...),
		"LegislativeBody":          legislativeBodyData,
		"LegislationType":          legislationTypeData,
		"ResponsibleUser":          userData,
		"ImportanceClassification": importanceClassificationData,
		"SecrecyClassification":    secrecyClassificationData,
		"ComplianceFolder":         complianceFolderData,
		"IssuingAuthority":         issuingAuthorityData,
		"ApplicabilityStatement":   applicabilityData,
		"AssessSettings":           assessSettingsData,
		"SharedUnit":               sharedUnitData,
	}, dropDownConfig)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func (crs *ComplianceRegisterService) ValidateComplianceTemplate(file io.Reader) ([]map[string]interface{}, error) {
	rows, err := utility.ReadExcelFile(file)
	if err != nil {
		return nil, err
	}

	headers := []string{
		"LegislativeBody", "LegislationType", "LegislationNo", "ResponsibleUser",
		"Title", "Description", "Version", "Deadline",
		"ImportanceClassification", "SecrecyClassification", "ComplianceFolder", "IssuingAuthority",
		"IssuingAuthorityCircularNo", "IssuingAuthorityDateOfReceiving", "IssuingAuthorityHijriDateOfReceiving",
		"ApplicabilityStatement", "AssessSettings", "SharedUnit",
	}

	snakeHeaders := []string{
		"legislative_body", "legislation_type", "legislation_no", "responsible_user",
		"title", "description", "version", "deadline",
		"importance_classification", "secrecy_classification", "compliance_folder", "issuing_authority",
		"issuing_authority_circular_no", "issuing_authority_date_of_receiving", "issuing_authority_hijri_date_of_receiving",
		"applicability_statement", "assess_settings", "shared_unit",
	}

	if len(rows) == 0 {
		return nil, fmt.Errorf("empty Excel file")
	}
	if !utility.ValidateHeaders(headers, rows[0]) {
		return nil, utility.NewValidationError(fmt.Errorf("invalid headers"), "invalid_header", "Compliance Template")
	}

	var input []map[string]interface{}
	for _, row := range rows[1:] {
		record := make(map[string]interface{})
		for i, header := range snakeHeaders {
			if i < len(row) {
				record[header] = row[i]
			} else {
				record[header] = ""
			}
		}
		input = append(input, record)
	}

	validatedData, err := crs.validateCompliances(input)
	if err != nil {
		return nil, err
	}

	return validatedData, nil
}

func (crs *ComplianceRegisterService) ValidateImportCompliances(input []map[string]interface{}) ([]map[string]interface{}, error) {
	validatedData, err := crs.validateCompliances(input)
	if err != nil {
		return nil, err
	}

	return validatedData, nil
}

func (crs *ComplianceRegisterService) ImportCompliances(input []map[string]interface{}, currentUserID uint) ([]uint, error) {
	validatedData, err := crs.validateCompliances(input)
	if err != nil {
		return nil, err
	}

	for _, record := range validatedData {
		if errors, ok := record["errors"].(map[string]interface{}); ok && len(errors) > 0 {
			return nil, utility.NewValidationError(fmt.Errorf("invalid data in the import file"), "invalid_data", "Invalid Data")
		}
	}

	var createdIDs []uint
	for _, record := range validatedData {

		// Convert date format from MM-DD-YY to YYYY-MM-DD
		convertDateFormat := func(dateStr string) string {
			if dateStr == "" {
				return ""
			}
			// Try to parse MM-DD-YY format
			if parsed, err := time.Parse("01-02-06", dateStr); err == nil {
				return parsed.Format("2006-01-02")
			}
			// Try to parse YYYY-MM-DD format (already correct)
			if _, err := time.Parse("2006-01-02", dateStr); err == nil {
				return dateStr
			}
			// Return as-is if parsing fails
			return dateStr
		}

		// Build IssuingAuthorities array if issuing_authority_id exists
		var issuingAuthorities []complianceDTOs.ComplianceRegisterIssuingAuthorityFormDTO
		if issuingAuthorityID := utility.GetOptionalUintValue(record["issuing_authority_id"]); issuingAuthorityID != nil {
			dateOfReceiving := utility.GetOptionalStringValue(record["issuing_authority_date_of_receiving"])
			if dateOfReceiving != nil {
				convertedDate := convertDateFormat(*dateOfReceiving)
				dateOfReceiving = &convertedDate
			}
			entity := &complianceDTOs.ComplianceRegisterIssuingAuthorityEntityFormDTO{
				CircularNo:           utility.GetOptionalStringValue(record["issuing_authority_circular_no"]),
				DateOfReceiving:      dateOfReceiving,
				HijriDateOfReceiving: utility.GetOptionalStringValue(record["issuing_authority_hijri_date_of_receiving"]),
			}
			issuingAuthorities = append(issuingAuthorities, complianceDTOs.ComplianceRegisterIssuingAuthorityFormDTO{
				IssuingAuthorityID: issuingAuthorityID,
				Entity:             entity,
			})
		}

		// Build SharedUnitIDs array if shared_unit_id exists
		var sharedUnitIDs []*uint
		if sharedUnitID := utility.GetOptionalUintValue(record["shared_unit_id"]); sharedUnitID != nil {
			sharedUnitIDs = append(sharedUnitIDs, sharedUnitID)
		}

		// Get title and description values - will be stored in both En and Ar fields
		title := record["title"].(string)
		description := utility.GetOptionalStringValue(record["description"])

		// Convert deadline date format
		deadline := utility.GetOptionalStringValue(record["deadline"])
		if deadline != nil {
			convertedDeadline := convertDateFormat(*deadline)
			deadline = &convertedDeadline
		}

		dto := &complianceDTOs.ComplianceRegisterFormDTO{
			LegislativeBodyID: utility.GetOptionalUintValue(record["legislative_body_id"]),
			LegislationTypeID: utility.GetOptionalUintValue(record["legislation_type_id"]),
			LegislationNo:     utility.GetOptionalStringValue(record["legislation_no"]),
			// ResponsibleUnitID is auto-derived from ResponsibleUserID in Create method
			ResponsibleUserID:          utility.GetOptionalUintValue(record["responsible_user_id"]),
			TitleEn:                    title,
			TitleAr:                    title,
			DescriptionEn:              description,
			DescriptionAr:              description,
			Version:                    utility.GetOptionalStringValue(record["version"]),
			Deadline:                   deadline,
			ImportanceClassificationID: utility.GetOptionalUintValue(record["importance_classification_id"]),
			SecrecyClassificationID:    utility.GetOptionalUintValue(record["secrecy_classification_id"]),
			ComplianceFolderID:         utility.GetOptionalUintValue(record["compliance_folder_id"]),
			IsApplicable: func() bool {
				if v, ok := record["is_applicable"].(bool); ok {
					return v
				}
				return false
			}(),
			IsNotApplicable: func() bool {
				if v, ok := record["is_not_applicable"].(bool); ok {
					return v
				}
				return false
			}(),
			IsPublic: func() bool {
				if v, ok := record["is_public"].(bool); ok {
					return v
				}
				return false
			}(),
			IsShared: func() bool {
				if v, ok := record["is_shared"].(bool); ok {
					return v
				}
				return false
			}(),
			IsRequirementAssessment: true,
			IssuingAuthorities:      issuingAuthorities,
			SharedUnitIDs:           sharedUnitIDs,
		}

		id, err := crs.Create(currentUserID, dto)
		if err != nil {
			return nil, fmt.Errorf("failed to create compliance for record: %+v, error: %v", record, err)
		}
		createdIDs = append(createdIDs, id)
	}

	return createdIDs, nil
}

func (crs *ComplianceRegisterService) validateCompliances(input []map[string]interface{}) ([]map[string]interface{}, error) {
	headers := []string{
		"legislative_body", "legislation_type", "legislation_no", "responsible_user",
		"title", "description", "version", "deadline",
		"importance_classification", "secrecy_classification", "compliance_folder", "issuing_authority",
		"issuing_authority_circular_no", "issuing_authority_date_of_receiving", "issuing_authority_hijri_date_of_receiving",
		"applicability_statement", "assess_settings", "shared_unit",
	}

	// Fetch Legislative Bodies
	legislativeBodies, err := crs.legislativeBodySvc.GetAll("", "", "", "active")
	if err != nil {
		return nil, fmt.Errorf("error fetching legislative bodies: %v", err)
	}
	legislativeBodyTitles := make(map[string]uint)
	for _, body := range legislativeBodies {
		legislativeBodyTitles[strings.TrimSpace(body.TitleEn)] = body.ID
	}

	// Fetch Legislation Types
	legislationTypes, err := crs.legislationTypeSvc.GetAll("", "", "", "active")
	if err != nil {
		return nil, fmt.Errorf("error fetching legislation types: %v", err)
	}
	legislationTypeTitles := make(map[string]uint)
	for _, lType := range legislationTypes {
		legislationTypeTitles[strings.TrimSpace(lType.TitleEn)] = lType.ID
	}

	// Fetch users
	users, err := crs.userSvc.GetAllUsers()
	if err != nil {
		return nil, fmt.Errorf("error fetching users: %v", err)
	}
	userFullNames := make(map[string]uint)
	for _, user := range users {
		fullName := strings.TrimSpace(user.FirstNameEN + " " + user.LastNameEN)
		userFullNames[fullName] = uint(user.ID)
	}

	// Fetch Importance Classifications
	importanceClassifications, err := crs.importanceClassificationSvc.GetAll("", "", "", "active")
	if err != nil {
		return nil, fmt.Errorf("error fetching importance classifications: %v", err)
	}
	importanceClassificationTitles := make(map[string]uint)
	for _, ic := range importanceClassifications {
		importanceClassificationTitles[strings.TrimSpace(ic.TitleEn)] = ic.ID
	}

	// Fetch Secrecy Classifications
	secrecyClassifications, err := crs.secrecyClassificationSvc.GetAll("", "", "", "active")
	if err != nil {
		return nil, fmt.Errorf("error fetching secrecy classifications: %v", err)
	}
	secrecyClassificationTitles := make(map[string]uint)
	for _, sc := range secrecyClassifications {
		secrecyClassificationTitles[strings.TrimSpace(sc.TitleEn)] = sc.ID
	}

	// Fetch Compliance Folders
	complianceFolders, err := crs.complianceFolderSvc.GetAll("", "", "", "active")
	if err != nil {
		return nil, fmt.Errorf("error fetching compliance folders: %v", err)
	}
	complianceFolderTitles := make(map[string]uint)
	for _, folder := range complianceFolders {
		complianceFolderTitles[strings.TrimSpace(folder.TitleEn)] = folder.ID
	}

	// Fetch Issuing Authorities
	issuingAuthorities, err := crs.issuingAuthoritySvc.GetAll("", "", "", "active")
	if err != nil {
		return nil, fmt.Errorf("error fetching issuing authorities: %v", err)
	}
	issuingAuthorityTitles := make(map[string]uint)
	for _, ia := range issuingAuthorities {
		issuingAuthorityTitles[strings.TrimSpace(ia.TitleEn)] = ia.ID
	}

	// Fetch Shared Units (Departments - Level 2)
	sharedUnits, err := crs.organizationUnitSvc.GetAll(&settingsDTOs.OrganizationUnitListQueryParamsDTO{
		Status: "active",
		Level:  2,
	})
	if err != nil {
		return nil, fmt.Errorf("error fetching shared units: %v", err)
	}
	sharedUnitTitles := make(map[string]uint)
	for _, unit := range sharedUnits {
		sharedUnitTitles[strings.TrimSpace(unit.TitleEn)] = unit.ID
	}

	// For checking unique compliance titles
	existingCompliances, err := crs.repo.GetAllComplianceTitles()
	if err != nil {
		return nil, fmt.Errorf("error fetching existing compliance titles: %v", err)
	}
	existingTitleSet := make(map[string]bool)
	for _, title := range existingCompliances {
		existingTitleSet[strings.TrimSpace(strings.ToLower(title))] = true
	}

	// Track titles seen in the current Excel
	sheetTitleSet := make(map[string]bool)

	// Perform validation
	var validatedData []map[string]interface{}
	for _, record := range input {
		rowErrors := make(map[string]interface{})
		complianceRecord := make(map[string]interface{})

		// Copy original data
		for key, value := range record {
			complianceRecord[key] = value
		}

		for _, header := range headers {
			value, ok := record[header].(string)
			switch header {
			case "legislative_body":
				if !ok || value == "" {
					rowErrors[header] = "Legislative Body is required"
				} else if id, exists := legislativeBodyTitles[strings.TrimSpace(value)]; !exists {
					rowErrors[header] = fmt.Sprintf("Invalid legislative body: %s", value)
				} else {
					complianceRecord["legislative_body_id"] = id
				}

			case "legislation_type":
				if !ok || value == "" {
					rowErrors[header] = "Legislation Type is required"
				} else if id, exists := legislationTypeTitles[strings.TrimSpace(value)]; !exists {
					rowErrors[header] = fmt.Sprintf("Invalid legislation type: %s", value)
				} else {
					complianceRecord["legislation_type_id"] = id
				}

			case "legislation_no":
				if !ok || value == "" {
					rowErrors[header] = "Legislation No is required"
				}

			case "responsible_user":
				if !ok || value == "" {
					rowErrors[header] = "Responsible User is required"
				} else if id, exists := userFullNames[strings.TrimSpace(value)]; !exists {
					rowErrors[header] = fmt.Sprintf("Invalid responsible user: %s", value)
				} else {
					complianceRecord["responsible_user_id"] = id
				}

			case "importance_classification":
				if !ok || value == "" {
					rowErrors[header] = "Importance Classification is required"
				} else if id, exists := importanceClassificationTitles[strings.TrimSpace(value)]; !exists {
					rowErrors[header] = fmt.Sprintf("Invalid importance classification: %s", value)
				} else {
					complianceRecord["importance_classification_id"] = id
				}

			case "secrecy_classification":
				if !ok || value == "" {
					rowErrors[header] = "Secrecy Classification is required"
				} else if id, exists := secrecyClassificationTitles[strings.TrimSpace(value)]; !exists {
					rowErrors[header] = fmt.Sprintf("Invalid secrecy classification: %s", value)
				} else {
					complianceRecord["secrecy_classification_id"] = id
				}

			case "compliance_folder":
				if ok && value != "" {
					if id, exists := complianceFolderTitles[strings.TrimSpace(value)]; !exists {
						rowErrors[header] = fmt.Sprintf("Invalid compliance folder: %s", value)
					} else {
						complianceRecord["compliance_folder_id"] = id
					}
				}

			case "issuing_authority":
				if !ok || value == "" {
					rowErrors[header] = "Issuing Authority is required"
				} else if id, exists := issuingAuthorityTitles[strings.TrimSpace(value)]; !exists {
					rowErrors[header] = fmt.Sprintf("Invalid issuing authority: %s", value)
				} else {
					complianceRecord["issuing_authority_id"] = id
				}

			case "shared_unit":
				// Check if Shared is selected in AssessSettings
				isShared := false
				if sharedVal, ok := complianceRecord["is_shared"].(bool); ok && sharedVal {
					isShared = true
				}

				if isShared && (!ok || value == "") {
					rowErrors[header] = "Shared Unit is required when Assess Settings is 'Shared'"
				} else if ok && value != "" {
					if id, exists := sharedUnitTitles[strings.TrimSpace(value)]; !exists {
						rowErrors[header] = fmt.Sprintf("Invalid shared unit: %s", value)
					} else {
						complianceRecord["shared_unit_id"] = id
					}
				}

			case "title":
				if !ok || value == "" {
					rowErrors[header] = "Title is required"
				} else {
					trimmedTitle := strings.TrimSpace(strings.ToLower(value))
					if existingTitleSet[trimmedTitle] {
						rowErrors[header] = fmt.Sprintf("Compliance with title '%s' already exists in the system", value)
					} else if sheetTitleSet[trimmedTitle] {
						rowErrors[header] = fmt.Sprintf("Duplicate title '%s' found in the sheet", value)
					} else {
						sheetTitleSet[trimmedTitle] = true
					}
				}

			case "applicability_statement":
				if ok && value != "" {
					trimmedValue := strings.TrimSpace(value)
					if trimmedValue == "Applicable" {
						complianceRecord["is_applicable"] = true
						complianceRecord["is_not_applicable"] = false
					} else if trimmedValue == "Not Applicable" {
						complianceRecord["is_applicable"] = false
						complianceRecord["is_not_applicable"] = true
					} else {
						rowErrors[header] = fmt.Sprintf("Invalid applicability statement: %s (must be 'Applicable' or 'Not Applicable')", value)
					}
				}

			case "assess_settings":
				if ok && value != "" {
					trimmedValue := strings.TrimSpace(value)
					if trimmedValue == "Public" {
						complianceRecord["is_public"] = true
						complianceRecord["is_shared"] = false
					} else if trimmedValue == "Shared" {
						complianceRecord["is_public"] = false
						complianceRecord["is_shared"] = true
					} else {
						rowErrors[header] = fmt.Sprintf("Invalid assess settings: %s (must be 'Public' or 'Shared')", value)
					}
				}
			}
		}

		complianceRecord["errors"] = rowErrors
		validatedData = append(validatedData, complianceRecord)
	}

	return validatedData, nil
}
