package services

import (
	"errors"
	"fmt"
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
	Archive(ids []uint, currentUserID uint) ([]uint, error)
	Restore(ids []uint, currentUserID uint) ([]uint, error)
	StatusUpdate(id uint, currentUserID uint, complianceStatusID uint, tx utility.DatabaseTransaction) error
	ReviewStatusUpdate(complianceRegister *models.ComplianceRegister, tx utility.DatabaseTransaction) error

	UpdateComplianceTotalScore(complianceID uint, totalScore float64, tx utility.DatabaseTransaction) error

	GetDocumentPath(id uint, token string) (string, error)

	ExportComplianceRegisters(data []*complianceDTOs.ComplianceRegisterListDTO, currentUserID uint) (*os.File, error)
}

// Service and implementations
type ComplianceRegisterService struct {
	repo                  complianceRepositories.ComplianceRegisterRepositoryInterface
	activityLogSvc        logServices.ActivityLogServiceInterface
	lifecycleStatusSvc    masterServices.LifecycleStatusServiceInterface
	complianceCategorySvc ComplianceCategoryServiceInterface
	complianceStatusSvc   ComplianceStatusServiceInterface
	reviewFrequencySvc    masterServices.ReviewFrequencyServiceInterface
	severityLevelSvc      masterServices.SeverityLevelServiceInterface
	controlSvc            organizationServices.ControlServiceInterface
	userSettingSvc        organizationServices.UserSettingServiceInterface
	userSvc               organizationServices.UserServiceInterface
	notificationSvc       commonServices.NotificationServiceInterface
	complianceFolderSvc   ComplianceFolderServiceInterface
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
) ComplianceRegisterServiceInterface {
	return &ComplianceRegisterService{
		repo:                  repo,
		activityLogSvc:        activityLogSvc,
		lifecycleStatusSvc:    lifecycleStatusSvc,
		complianceCategorySvc: complianceCategorySvc,
		complianceStatusSvc:   complianceStatusSvc,
		reviewFrequencySvc:    reviewFrequencySvc,
		severityLevelSvc:      severityLevelSvc,
		controlSvc:            controlSvc,
		userSettingSvc:        userSettingSvc,
		userSvc:               userSvc,
		notificationSvc:       notificationSvc,
		complianceFolderSvc:   complianceFolderSvc,
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

	complianceDTO.CreatedAt = compliance.CreatedAt
	complianceDTO.UpdatedAt = compliance.UpdatedAt
	// todo :- also add extra fields which required in details screen

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

	lifecycleStatus, err := crs.lifecycleStatusSvc.GetByType("active")
	if err != nil {
		return 0, err
	}

	if err := crs.validateComplianceDocumentFormDTO(&complianceRegisterFormDTO.DocumentEn); err != nil {
		return 0, err
	}

	if err := crs.validateComplianceDocumentFormDTO(&complianceRegisterFormDTO.DocumentAr); err != nil {
		return 0, err
	}

	responsibleUnitID, err := crs.userSvc.GetUserPrimaryUnitID(uint(complianceRegisterFormDTO.ResponsibleUserID))
	if err != nil {
		return 0, utility.NewValidationError(fmt.Errorf("organization unit not set for responsible user"), "not_exist", "bcp communication plan responsible unit")
	}

	complianceRegister := &models.ComplianceRegister{
		ComplianceCategoryID:     complianceRegisterFormDTO.ComplianceCategoryID,
		ComplianceStatusID:       complianceRegisterFormDTO.ComplianceStatusID,
		ResponsibleUnitID:        responsibleUnitID,
		ResponsibleUserID:        complianceRegisterFormDTO.ResponsibleUserID,
		ReviewFrequencyID:        complianceRegisterFormDTO.ReviewFrequencyID,
		LastReviewNote:           complianceRegisterFormDTO.LastReviewNote,
		ReviewResponsibleUserID:  complianceRegisterFormDTO.ReviewResponsibleUserID,
		SeverityLevelID:          complianceRegisterFormDTO.SeverityLevelID,
		ReferenceNumber:          referenceNumber,
		StandardCode:             complianceRegisterFormDTO.StandardCode,
		TitleEn:                  complianceRegisterFormDTO.TitleEn,
		TitleAr:                  complianceRegisterFormDTO.TitleAr,
		DescriptionEn:            complianceRegisterFormDTO.DescriptionEn,
		DescriptionAr:            complianceRegisterFormDTO.DescriptionAr,
		Version:                  complianceRegisterFormDTO.Version,
		Deadline:                 dates.Deadline,
		IssueDate:                dates.IssueDate,
		ExpiryDate:               dates.ExpiryDate,
		NextReviewDate:           dates.NextReviewDate,
		LastReviewedDate:         dates.LastReviewedDate,
		LastReviewedUserID:       complianceRegisterFormDTO.LastReviewedUserID,
		ComplianceReviewStatusID: complianceRegisterFormDTO.ComplianceReviewStatusID,
		CreatedByID:              currentUserID,
		LifecycleStatusID:        lifecycleStatus.ID,
		ComplianceFolderID:       complianceRegisterFormDTO.ComplianceFolderID,
	}

	id, err := crs.repo.Create(complianceRegister, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	complianceStatusID := complianceRegisterFormDTO.ComplianceStatusID
	complianceStatusUpdate := &models.ComplianceStatusUpdate{
		ComplianceRegisterID: id,
		ComplianceStatusID:   complianceStatusID,
		Comment:              "Initial compliance status",
		CreatedByID:          currentUserID,
	}

	_, err = crs.repo.CreateInitialComplianceStatusUpdate(complianceStatusUpdate, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	// Save Controls
	var controls []*models.ComplianceControl
	for _, ctrl := range complianceRegisterFormDTO.ControlIDs {
		control := &models.ComplianceControl{
			ComplianceRegisterID:    id,
			ControlID:               *ctrl,
			ComplianceRequirementID: nil,
			CreatedByID:             currentUserID,
		}
		controls = append(controls, control)
	}

	if len(controls) > 0 {
		if err := crs.repo.CreateComplianceControls(controls, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	if err := crs.handleDocument(currentUserID, id, &complianceRegisterFormDTO.DocumentEn, tx); err != nil {
		tx.Rollback()
		return 0, err
	}

	if err := crs.handleDocument(currentUserID, id, &complianceRegisterFormDTO.DocumentAr, tx); err != nil {
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
			nil,
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
			UserIDs:              []uint{complianceRegisterFormDTO.ResponsibleUserID},
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

		responsibleUser, err := crs.userSvc.GetByID(complianceRegisterFormDTO.ResponsibleUserID)
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

		// Assuming GetByID returns only an error and no severity level
		severityLevel, err := crs.severityLevelSvc.GetByID(*complianceRegister.SeverityLevelID)
		if err != nil {
			utility.LogError("Failed to fetch severity level", map[string]interface{}{
				"error":   err.Error(),
				"levelID": complianceRegister.SeverityLevelID,
			})
			return
		}

		emailTo := []string{responsibleUser.Email}
		subject := fmt.Sprintf("Compliance Register/%s: Responsibility Assigned", complianceRegister.TitleEn)
		appURL := utility.GetAppURL()

		emailData := map[string]interface{}{
			"FirstNameEn":               responsibleUser.FirstNameEn,
			"LastNameEn":                responsibleUser.LastNameEn,
			"ComplianceReferenceNumber": complianceRegister.ReferenceNumber,
			"ComplianceTitleEn":         complianceRegister.TitleEn,
			"ComplianceTitleAr":         complianceRegister.TitleAr,
			"ComplianceDescription":     complianceRegister.DescriptionEn,
			"SeverityLevel":             severityLevel.TitleEn,
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

	if err := crs.validateComplianceDocumentFormDTO(&complianceRegisterFormDTO.DocumentEn); err != nil {
		return err
	}

	if err := crs.validateComplianceDocumentFormDTO(&complianceRegisterFormDTO.DocumentEn); err != nil {
		return err
	}

	oldComplianceRegister, err := crs.GetFullDetailsByID(id, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	if existingComplianceRegister.ResponsibleUserID != complianceRegisterFormDTO.ResponsibleUserID {
		responsibleUnitID, err := crs.userSvc.GetUserPrimaryUnitID(existingComplianceRegister.ResponsibleUserID)
		if err != nil {
			return utility.NewValidationError(fmt.Errorf("organization unit not set for responsible user"), "not_exist", "risk responsible unit")
		}
		existingComplianceRegister.ResponsibleUnitID = responsibleUnitID
	}

	existingComplianceRegister.ComplianceCategoryID = complianceRegisterFormDTO.ComplianceCategoryID
	existingComplianceRegister.ComplianceStatusID = complianceRegisterFormDTO.ComplianceStatusID
	existingComplianceRegister.ResponsibleUserID = complianceRegisterFormDTO.ResponsibleUserID
	existingComplianceRegister.ReviewFrequencyID = complianceRegisterFormDTO.ReviewFrequencyID
	existingComplianceRegister.ReviewResponsibleUserID = complianceRegisterFormDTO.ReviewResponsibleUserID
	existingComplianceRegister.SeverityLevelID = complianceRegisterFormDTO.SeverityLevelID
	existingComplianceRegister.StandardCode = complianceRegisterFormDTO.StandardCode
	existingComplianceRegister.TitleEn = complianceRegisterFormDTO.TitleEn
	existingComplianceRegister.TitleAr = complianceRegisterFormDTO.TitleAr
	existingComplianceRegister.DescriptionEn = complianceRegisterFormDTO.DescriptionEn
	existingComplianceRegister.DescriptionAr = complianceRegisterFormDTO.DescriptionAr
	existingComplianceRegister.Version = complianceRegisterFormDTO.Version
	existingComplianceRegister.Deadline = dates.Deadline
	existingComplianceRegister.IssueDate = dates.IssueDate
	existingComplianceRegister.ExpiryDate = dates.ExpiryDate
	existingComplianceRegister.NextReviewDate = dates.NextReviewDate
	existingComplianceRegister.LastReviewedDate = dates.LastReviewedDate
	existingComplianceRegister.LastReviewedUserID = complianceRegisterFormDTO.LastReviewedUserID
	existingComplianceRegister.ComplianceReviewStatusID = complianceRegisterFormDTO.ComplianceReviewStatusID
	existingComplianceRegister.LastReviewNote = complianceRegisterFormDTO.LastReviewNote
	existingComplianceRegister.UpdatedByID = &currentUserID
	existingComplianceRegister.ComplianceFolderID = complianceRegisterFormDTO.ComplianceFolderID

	err = crs.repo.Update(existingComplianceRegister, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Collect existing control IDs
	existingCompliance, err := crs.GetFullDetailsByID(id, tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	existingControlIDs := make(map[uint]bool)
	for _, ctrl := range existingCompliance.Controls {
		existingControlIDs[ctrl.ControlID] = true
	}

	// Collect new control IDs
	newControlIDs := make(map[uint]bool)
	for _, ctrl := range complianceRegisterFormDTO.ControlIDs {
		newControlIDs[*ctrl] = true
	}

	// Add new controls
	var controlsToAdd []*models.ComplianceControl
	for _, newControlID := range complianceRegisterFormDTO.ControlIDs {
		if _, exists := existingControlIDs[*newControlID]; !exists {
			// Create a new compliance control record
			control := &models.ComplianceControl{
				ComplianceRegisterID:    id,
				ComplianceRequirementID: nil,
				ControlID:               *newControlID,
				CreatedByID:             currentUserID,
			}
			controlsToAdd = append(controlsToAdd, control)
		}
	}

	// Save the new controls
	if len(controlsToAdd) > 0 {
		if err := crs.repo.CreateComplianceControls(controlsToAdd, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	// Remove controls not in the new list
	var deleteControlIDs []*uint

	for _, existingControl := range existingCompliance.Controls {
		if _, exists := newControlIDs[existingControl.ControlID]; !exists {
			// Append control to delete list
			deleteControlIDs = append(deleteControlIDs, existingControl.ID)
		}
	}

	// Delete controls not in the new list
	if len(deleteControlIDs) > 0 {
		if err := crs.repo.DeleteComplianceControls(deleteControlIDs, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	existingComplianceDocumentEn, err := crs.repo.GetLatestDocumentByComplianceID(id, "en")
	if err != nil {
		return err
	}

	existingComplianceDocumentAr, err := crs.repo.GetLatestDocumentByComplianceID(id, "ar")
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	if complianceRegisterFormDTO.DocumentEn.IsNewFile != nil && *complianceRegisterFormDTO.DocumentEn.IsNewFile {
		existingComplianceDocumentEnFileName := existingComplianceDocumentEn.Name

		if err := crs.handleDocument(currentUserID, id, &complianceRegisterFormDTO.DocumentEn, tx); err != nil {
			tx.Rollback()
			return err
		}

		storagePath := filepath.Join("compliance", "compliance-register", strconv.FormatUint(uint64(id), 10), "documents")
		if err := utility.RemoveFile(existingComplianceDocumentEnFileName, storagePath); err != nil {
			tx.Rollback()
			return err
		}

	}

	if complianceRegisterFormDTO.DocumentAr.IsNewFile != nil && *complianceRegisterFormDTO.DocumentAr.IsNewFile {
		var existingComplianceDocumentArFileName string
		if existingComplianceDocumentAr != nil {
			existingComplianceDocumentArFileName = existingComplianceDocumentAr.Name
		}

		if err := crs.handleDocument(currentUserID, id, &complianceRegisterFormDTO.DocumentAr, tx); err != nil {
			tx.Rollback()
			return err
		}

		if existingComplianceDocumentAr != nil {
			storagePath := filepath.Join("compliance", "compliance-register", strconv.FormatUint(uint64(id), 10), "documents")
			if err := utility.RemoveFile(existingComplianceDocumentArFileName, storagePath); err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	tx.Commit()

	go func() {
		existingCompliance, err := crs.GetFullDetailsByID(id)
		if err != nil {
			return
		}
		// Determine control and document changes
		controlsAdded := []uint{}
		controlsRemoved := []uint{}
		for id := range newControlIDs {
			if !existingControlIDs[id] {
				controlsAdded = append(controlsAdded, id)
			}
		}
		for id := range existingControlIDs {
			if !newControlIDs[id] {
				controlsRemoved = append(controlsRemoved, id)
			}
		}
		// Prepare log data
		changedData := make(map[string]interface{})
		associatedItems := make(map[string]interface{})

		if len(controlsAdded) > 0 || len(controlsRemoved) > 0 {
			associatedItems["control_ids"] = map[string]interface{}{
				"added":   controlsAdded,
				"removed": controlsRemoved,
			}
		}
		if len(associatedItems) > 0 {
			changedData["associated_items"] = associatedItems
		}
		// log changes
		oldData := map[string]interface{}{
			"compliance_category":      oldComplianceRegister.ComplianceCategory.TitleEn,
			"compliance_status":        oldComplianceRegister.ComplianceStatus.TitleEn,
			"responsible_unit":         oldComplianceRegister.ResponsibleUnit.TitleEn,
			"responsible_user":         fmt.Sprintf("%s %s", oldComplianceRegister.ResponsibleUser.FirstNameEN, oldComplianceRegister.ResponsibleUser.LastNameEN),
			"review_frequency":         oldComplianceRegister.ReviewFrequency.TitleEn,
			"review_responsible_user":  utility.GetFieldValue(oldComplianceRegister.ReviewResponsibleUser).(string),
			"severity_level":           oldComplianceRegister.SeverityLevel.TitleEn,
			"standard_code":            oldComplianceRegister.StandardCode,
			"title_en":                 oldComplianceRegister.TitleEn,
			"title_ar":                 oldComplianceRegister.TitleAr,
			"description_en":           oldComplianceRegister.DescriptionEn,
			"description_ar":           oldComplianceRegister.DescriptionAr,
			"version":                  oldComplianceRegister.Version,
			"dead_line":                oldComplianceRegister.DeadlineString,
			"issue_date":               oldComplianceRegister.IssueDateString,
			"expiry_date":              oldComplianceRegister.ExpiryDateString,
			"last_review_date":         oldComplianceRegister.LastReviewedDateString,
			"next_review_date":         oldComplianceRegister.NextReviewDateString,
			"last_review_note":         oldComplianceRegister.LastReviewNote,
			"last_review_user":         utility.GetFieldValue(oldComplianceRegister.LastReviewedUser).(string),
			"compliance_review_status": utility.GetFieldValue(oldComplianceRegister.ComplianceReviewStatus).(string),
		}

		newData := map[string]interface{}{
			"compliance_category":      existingCompliance.ComplianceCategory.TitleEn,
			"compliance_status":        existingCompliance.ComplianceStatus.TitleEn,
			"responsible_unit":         existingCompliance.ResponsibleUnit.TitleEn,
			"responsible_user":         fmt.Sprintf("%s %s", existingCompliance.ResponsibleUser.FirstNameEN, existingCompliance.ResponsibleUser.LastNameEN),
			"review_frequency":         existingCompliance.ReviewFrequency.TitleEn,
			"review_responsible_user":  utility.GetFieldValue(existingCompliance.ReviewResponsibleUser).(string),
			"description_ar":           existingCompliance.DescriptionAr,
			"version":                  existingCompliance.Version,
			"dead_line":                existingCompliance.DeadlineString,
			"issue_date":               existingCompliance.IssueDateString,
			"expiry_date":              existingCompliance.ExpiryDateString,
			"last_review_date":         existingCompliance.LastReviewedDateString,
			"next_review_date":         existingCompliance.NextReviewDateString,
			"last_review_note":         existingCompliance.LastReviewNote,
			"last_review_user":         utility.GetFieldValue(existingCompliance.LastReviewedUser).(string),
			"compliance_review_status": utility.GetFieldValue(existingCompliance.ComplianceReviewStatus).(string),
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
				nil,
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
				UserIDs:              []uint{oldResponsibleUserID},
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
				UserIDs:              []uint{complianceRegisterFormDTO.ResponsibleUserID},
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
			oldResponsibleUser, err := crs.userSvc.GetByID(oldResponsibleUserID)
			if err != nil {
				utility.LogError("Failed to fetch old responsible user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}

			newResponsibleUser, err := crs.userSvc.GetByID(complianceRegisterFormDTO.ResponsibleUserID)
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

		// Check if controls exist
		if len(complianceRegister.Controls) > 0 {
			var deleteControlIds []*uint
			for _, ctrl := range complianceRegister.Controls {
				deleteControlIds = append(deleteControlIds, ctrl.ID)
			}
			if err := crs.repo.DeleteComplianceControls(deleteControlIds, tx); err != nil {
				tx.Rollback()
				return nil, err
			}
		}

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

	// Deleting files should not affect the main operation result
	if len(successfulIds) > 0 {
		for complianceID, documents := range documentsToDelete {
			for _, doc := range documents {
				storagePath := filepath.Join("compliance", "compliance-register", strconv.FormatUint(uint64(complianceID), 10), "documents")
				if err := utility.RemoveFile(doc.Name, storagePath); err != nil {
					tx.Rollback()
					fmt.Printf("Failed to delete file at URL %s: %v\n", doc.URL, err)
				}
			}
			destinationPath := filepath.Join("compliance", "compliance-register", strconv.FormatUint(uint64(complianceID), 10))
			if err := utility.RemoveFolder(destinationPath); err != nil {
				tx.Rollback()
			}
		}
	}

	tx.Commit()

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

			// Send DB and email notifications for each successfully deleted compliance register
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

	existingCompliance.ComplianceStatusID = complianceStatusID
	existingCompliance.UpdatedByID = &currentUserID

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

	if _, err := crs.complianceCategorySvc.GetByID(complianceRegisterFormDTO.ComplianceCategoryID); err != nil {
		return utility.NewValidationError(fmt.Errorf("compliance category does not exist"), "not_exist", "compliance category")
	}

	if _, err := crs.complianceStatusSvc.GetByID(complianceRegisterFormDTO.ComplianceStatusID); err != nil {
		return utility.NewValidationError(fmt.Errorf("compliance status does not exist"), "not_exist", "compliance status")
	}

	if complianceRegisterFormDTO.ReviewFrequencyID != nil {
		if _, err := crs.reviewFrequencySvc.GetByID(*complianceRegisterFormDTO.ReviewFrequencyID); err != nil {
			return utility.NewValidationError(fmt.Errorf("review frequency does not exist"), "not_exist", "review frequency")
		}
	}
	if complianceRegisterFormDTO.SeverityLevelID != nil {
		if _, err := crs.severityLevelSvc.GetByID(*complianceRegisterFormDTO.SeverityLevelID); err != nil {
			return utility.NewValidationError(fmt.Errorf("severity level does not exist"), "not_exist", "severity level")
		}
	}

	// Controls exist
	controls := complianceRegisterFormDTO.ControlIDs

	controlIDs := make(map[uint]bool)

	for _, control := range controls {
		if controlIDs[*control] {
			return utility.NewValidationError(fmt.Errorf("duplicate control ID: %d", control), "duplicate", "control")
		}

		if _, err := crs.controlSvc.GetByID(*control); err != nil {
			return utility.NewValidationError(fmt.Errorf("control does not exist for ID: %d", control), "not_exist", "control")
		}

		controlIDs[*control] = true
	}

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

	if dto.IssueDate != nil {
		parsed.IssueDate, err = utility.ParseDate(*dto.IssueDate)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid issue date"), "invalid_date", "issue date")
		}
	}

	if dto.ExpiryDate != nil {
		parsed.ExpiryDate, err = utility.ParseDate(*dto.ExpiryDate)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid expiry date"), "invalid_date", "expiry date")
		}
	}

	if dto.NextReviewDate != nil {
		parsed.NextReviewDate, err = utility.ParseDate(*dto.NextReviewDate)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid next review date"), "invalid_date", "next review date")
		}
	}

	if dto.LastReviewedDate != nil {
		parsed.LastReviewedDate, err = utility.ParseDate(*dto.LastReviewedDate)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid last reviewed date"), "invalid_date", "last reviewed date")
		}
	}

	return parsed, nil
}

func (crs *ComplianceRegisterService) UpdateComplianceTotalScore(complianceID uint, totalScore float64, tx utility.DatabaseTransaction) error {
	return crs.repo.UpdateComplianceTotalScore(complianceID, totalScore, tx)
}

func (crs *ComplianceRegisterService) validateComplianceDocumentFormDTO(document *complianceDTOs.ComplianceDocumentFormDTO) error {
	// Validate document
	if document == nil {
		return utility.NewValidationError(fmt.Errorf("document is required"), "required", "document")
	}
	// Normalize the file type to lowercase for case-insensitive comparison
	fileType := strings.ToLower(document.Type)

	// Validate file type (PDF)
	if fileType != ".pdf" {
		return utility.NewValidationError(fmt.Errorf("only PDF files are allowed"), "invalid_file_type", "type")
	}

	// @todo
	// Validate file size (if needed)
	// For example:
	// if document.Size > MaxFileSize {
	//     return utility.NewValidationError(fmt.Errorf("file size exceeds the limit"), "invalid_file_size", "file_size")
	// }

	return nil
}

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

func (crs *ComplianceRegisterService) ExportComplianceRegisters(data []*complianceDTOs.ComplianceRegisterListDTO, currentUserID uint) (*os.File, error) {
	// Retrieve user settings for language preference
	userSetting, err := crs.userSettingSvc.GetUserSetting(currentUserID)
	if err != nil {
		return nil, err
	}

	var headers []string
	// Define headers based on the user's language setting
	if userSetting.Language == "ar" {
		headers = []string{
			"SL No.",
			"Reference Number",
			"Standard Code",
			"Title (AR)",
			"Category (AR)",
			"Deadline",
			"Issue Date",
			"Expiry Date",
			"Description (AR)",
			"Lifecycle Status (AR)",
			"Responsible Unit (AR)",
			"Compliance Status (AR)",
			"Severity Level (AR)",
			"Responsible User (AR)",
			"Review Frequency (AR)",
			"Frequency Responsible User (AR)",
		}
	} else {
		headers = []string{
			"SL No.",
			"Reference Number",
			"Standard Code",
			"Title (EN)",
			"Category (EN)",
			"Deadline",
			"Issue Date",
			"Expiry Date",
			"Description (EN)",
			"Lifecycle Status (EN)",
			"Responsible Unit (EN)",
			"Compliance Status (EN)",
			"Severity Level (EN)",
			"Responsible User (EN)",
			"Review Frequency (EN)",
			"Frequency Responsible User (EN)",
		}
	}

	// Convert data to [][]interface{} based on language preference
	rows := make([][]interface{}, len(data))
	for i, item := range data {
		if userSetting.Language == "ar" {
			rows[i] = []interface{}{
				i + 1,
				item.ReferenceNumber,
				item.StandardCode,
				item.TitleAr,
				item.ComplianceCategoryTitleAr,
				item.Deadline,
				item.IssueDate,
				item.ExpiryDate,
				item.DescriptionAr,
				item.LifecycleStatusTitleAr,
				item.ResponsibleUnitTitleAr,
				item.ComplianceStatusTitleAr,
				item.SeverityLevelTitleAr,
				item.ResponsibleUserFirstNameAr + " " + item.ResponsibleUserLastNameAr,
				item.ReviewFrequencyTitleAr,
				item.ReviewResponsibleUserFirstNameAr + " " + item.ReviewResponsibleUserLastNameAr,
			}
		} else {
			rows[i] = []interface{}{
				i + 1,
				item.ReferenceNumber,
				item.StandardCode,
				item.TitleEn,
				item.ComplianceCategoryTitleEn,
				item.Deadline,
				item.IssueDate,
				item.ExpiryDate,
				item.DescriptionEn,
				item.LifecycleStatusTitleEn,
				item.ResponsibleUnitTitleEn,
				item.ComplianceStatusTitleEn,
				item.SeverityLevelTitleEn,
				item.ResponsibleUserFirstNameEn + " " + item.ResponsibleUserLastNameEn,
				item.ReviewFrequencyTitleEn,
				item.ReviewResponsibleUserFirstNameEn + " " + item.ReviewResponsibleUserLastNameEn,
			}
		}
	}
	// Generate the Excel file
	return utility.GenerateExcelFile(headers, rows)
}

func (crs *ComplianceRegisterService) handleDocument(currentUserID uint, complianceRegisterID uint, documentDTO *complianceDTOs.ComplianceDocumentFormDTO, tx ...utility.DatabaseTransaction) error {
	if documentDTO == nil {
		return nil // Skip if no document is provided
	}

	// Generate encrypted token for the file name
	encryptedFilename, err := utility.EncryptData([]byte(documentDTO.FileName))
	if err != nil {
		return fmt.Errorf("failed to encrypt document filename: %w", err)
	}

	// Generate the document URL
	url := utility.GenerateDocumentURL("compliance", "compliance-register", strconv.FormatUint(uint64(complianceRegisterID), 10), "documents")

	// Create a new ComplianceDocument instance
	complianceDocument := &models.ComplianceDocument{
		ComplianceRegisterID: complianceRegisterID,
		Language:             documentDTO.Language,
		URL:                  url,
		Name:                 documentDTO.FileName,
		Size:                 documentDTO.Size,
		Token:                encryptedFilename,
		OriginalName:         documentDTO.OriginalName,
		Type:                 documentDTO.Type,
		CreatedByID:          currentUserID,
	}

	// Save the document record in the database
	if err := crs.repo.CreateDocument(complianceDocument, tx...); err != nil {
		return fmt.Errorf("failed to save document record: %w", err)
	}

	// Store the document file in the storage
	storagePath := filepath.Join("compliance", "compliance-register", strconv.FormatUint(uint64(complianceRegisterID), 10), "documents")
	if err := utility.StoreFile(documentDTO.FileName, storagePath); err != nil {
		return fmt.Errorf("failed to store document file: %w", err)
	}

	return nil
}
