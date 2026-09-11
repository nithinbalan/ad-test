package services

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"jethur/config"
	commonDTOs "jethur/grc/common/dtos"
	commonServices "jethur/grc/common/services"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	settingsServices "jethur/grc/settings/services"
	tiDTOs "jethur/grc/threat-intelligence/dtos"
	tiRepositories "jethur/grc/threat-intelligence/repositories"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type ThreatIntelligencePlannedAdvisoryServiceInterface interface {
	GetAll(params *tiDTOs.PlannedAdvisoryListQueryParamsDTO) (*tiDTOs.PlannedAdvisoryListResponseDTO, error)
	GetByID(id uint, tx ...utility.DatabaseTransaction) (*models.ThreatIntelligencePlannedAdvisory, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*tiDTOs.PlannedAdvisoryDetailDTO, error)
	Create(currentUserID uint, formDTO *tiDTOs.PlannedAdvisoryFormDTO) (uint, error)
	Update(id uint, currentUserID uint, formDTO *tiDTOs.PlannedAdvisoryFormDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	UpdateStatus(id uint, currentUserID uint, formDTO *tiDTOs.PlannedAdvisoryStatusUpdateFormDTO) error
	GetDocumentPath(plannedAdvisoryID uint, token string) (string, error)
	ExportPlannedAdvisory(data []*tiDTOs.PlannedAdvisoryListDTO, currentUserID uint) (*os.File, error)
}

type ThreatIntelligencePlannedAdvisoryService struct {
	repo                              tiRepositories.ThreatIntelligencePlannedAdvisoryRepositoryInterface
	activityLogSvc                    logServices.ActivityLogServiceInterface
	notificationSvc                   commonServices.NotificationServiceInterface
	advisorySourceSvc                 PlannedAdvisorySourceServiceInterface
	advisoryClassificationSvc         PlannedAdvisoryClassificationServiceInterface
	advisoryTypeSvc                   PlannedAdvisoryTypeServiceInterface
	advisoryImportanceSvc             PlannedAdvisoryImportanceServiceInterface
	advisorySecurityClassificationSvc PlannedAdvisorySecurityClassificationServiceInterface
	advisoryStatusSvc                 PlannedAdvisoryStatusServiceInterface
	advisoryInternalStatusSvc         PlannedAdvisoryInternalStatusServiceInterface
	likelihoodSvc                     ThreatIntelligenceLikelihoodServiceInterface
	severitySvc                       ThreatIntelligenceSeverityServiceInterface
	riskLevelSvc                      ThreatIntelligenceRiskLevelServiceInterface
	riskLevelValueSvc                 ThreatIntelligenceRiskLevelValueServiceInterface
	userSvc                           organizationServiceProvider
	organizationUnitSvc               settingsServices.OrganizationUnitServiceInterface
	userSettingSvc                    organizationServices.UserSettingServiceInterface
}

// organizationServiceProvider exposes only the methods we need from organization services
type organizationServiceProvider interface {
	GetByID(id uint) (*models.User, error)
}

func NewThreatIntelligencePlannedAdvisoryService(
	repo tiRepositories.ThreatIntelligencePlannedAdvisoryRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	notificationSvc commonServices.NotificationServiceInterface,
	advisorySourceSvc PlannedAdvisorySourceServiceInterface,
	advisoryClassificationSvc PlannedAdvisoryClassificationServiceInterface,
	advisoryTypeSvc PlannedAdvisoryTypeServiceInterface,
	advisoryImportanceSvc PlannedAdvisoryImportanceServiceInterface,
	advisorySecurityClassificationSvc PlannedAdvisorySecurityClassificationServiceInterface,
	advisoryStatusSvc PlannedAdvisoryStatusServiceInterface,
	advisoryInternalStatusSvc PlannedAdvisoryInternalStatusServiceInterface,
	likelihoodSvc ThreatIntelligenceLikelihoodServiceInterface,
	severitySvc ThreatIntelligenceSeverityServiceInterface,
	riskLevelSvc ThreatIntelligenceRiskLevelServiceInterface,
	riskLevelValueSvc ThreatIntelligenceRiskLevelValueServiceInterface,
	userSvc organizationServiceProvider,
	organizationUnitSvc settingsServices.OrganizationUnitServiceInterface,
	userSettingSvc organizationServices.UserSettingServiceInterface,
) ThreatIntelligencePlannedAdvisoryServiceInterface {
	return &ThreatIntelligencePlannedAdvisoryService{
		repo:                              repo,
		activityLogSvc:                    activityLogSvc,
		notificationSvc:                   notificationSvc,
		advisorySourceSvc:                 advisorySourceSvc,
		advisoryClassificationSvc:         advisoryClassificationSvc,
		advisoryTypeSvc:                   advisoryTypeSvc,
		advisoryImportanceSvc:             advisoryImportanceSvc,
		advisorySecurityClassificationSvc: advisorySecurityClassificationSvc,
		advisoryStatusSvc:                 advisoryStatusSvc,
		advisoryInternalStatusSvc:         advisoryInternalStatusSvc,
		likelihoodSvc:                     likelihoodSvc,
		severitySvc:                       severitySvc,
		riskLevelSvc:                      riskLevelSvc,
		riskLevelValueSvc:                 riskLevelValueSvc,
		userSvc:                           userSvc,
		organizationUnitSvc:               organizationUnitSvc,
		userSettingSvc:                    userSettingSvc,
	}
}

func (pas *ThreatIntelligencePlannedAdvisoryService) GetAll(params *tiDTOs.PlannedAdvisoryListQueryParamsDTO) (*tiDTOs.PlannedAdvisoryListResponseDTO, error) {
	return pas.repo.GetAll(params)
}

func (pas *ThreatIntelligencePlannedAdvisoryService) GetByID(id uint, tx ...utility.DatabaseTransaction) (*models.ThreatIntelligencePlannedAdvisory, error) {
	plannedAdvisory, err := pas.repo.GetByID(id, tx...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return plannedAdvisory, nil
}

func (pas *ThreatIntelligencePlannedAdvisoryService) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*tiDTOs.PlannedAdvisoryDetailDTO, error) {
	var detailDTO tiDTOs.PlannedAdvisoryDetailDTO

	plannedAdvisory, err := pas.repo.GetFullDetailsByID(id, tx...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	if err := mapstructure.Decode(plannedAdvisory, &detailDTO); err != nil {
		return nil, err
	}

	receiptDateStr := plannedAdvisory.ReceiptDate.Format("2006-01-02 15:04:05")
	detailDTO.ReceiptDateString = &receiptDateStr

	if plannedAdvisory.NotifyDate != nil {
		s := plannedAdvisory.NotifyDate.Format("2006-01-02 15:04:05")
		detailDTO.NotifyDateString = &s
	}
	if plannedAdvisory.DeadlineDate != nil {
		s := plannedAdvisory.DeadlineDate.Format("2006-01-02 15:04:05")
		detailDTO.DeadlineDateString = &s
	}
	if plannedAdvisory.ResolveDate != nil {
		s := plannedAdvisory.ResolveDate.Format("2006-01-02 15:04:05")
		detailDTO.ResolveDateString = &s
	}
	if plannedAdvisory.ReplyDate != nil {
		s := plannedAdvisory.ReplyDate.Format("2006-01-02 15:04:05")
		detailDTO.ReplyDateString = &s
	}

	detailDTO.CreatedAt = plannedAdvisory.CreatedAt
	detailDTO.UpdatedAt = plannedAdvisory.UpdatedAt

	return &detailDTO, nil
}

func (pas *ThreatIntelligencePlannedAdvisoryService) Create(currentUserID uint, formDTO *tiDTOs.PlannedAdvisoryFormDTO) (uint, error) {
	if err := pas.validatePlannedAdvisoryFormDTO(formDTO); err != nil {
		return 0, err
	}

	dates, err := pas.parsePlannedAdvisoryDates(formDTO)
	if err != nil {
		return 0, err
	}

	referenceNumber, err := pas.repo.GenerateReferenceNumber()
	if err != nil {
		return 0, err
	}

	openStatus, err := pas.advisoryStatusSvc.GetByType("open")
	if err != nil {
		return 0, utility.NewValidationError(fmt.Errorf("open advisory status does not exist"), "not_exist", "advisory_status")
	}

	plannedAdvisory := &models.ThreatIntelligencePlannedAdvisory{
		ReferenceNumber:          referenceNumber,
		CaseNumber:               formDTO.CaseNumber,
		TitleEn:                  formDTO.TitleEn,
		TitleAr:                  formDTO.TitleAr,
		ReceiptDate:              dates.ReceiptDate,
		CaseLink:                 formDTO.CaseLink,
		AdvisorySourceID:         formDTO.AdvisorySourceID,
		CaseClassificationID:     formDTO.CaseClassificationID,
		AdvisoryTypeID:           formDTO.AdvisoryTypeID,
		ImportanceID:             formDTO.ImportanceID,
		SecurityClassificationID: formDTO.SecurityClassificationID,
		ThreatType:               formDTO.ThreatType,
		ThreatSeverityID:         formDTO.ThreatSeverityID,
		LikelihoodID:             formDTO.LikelihoodID,
		Description:              formDTO.Description,
		AffectedAssets:           formDTO.AffectedAssets,
		Tags:                     formDTO.Tags,
		IOCs:                     formDTO.IOCs,
		AdvisoryRecipientID:      formDTO.AdvisoryRecipientID,
		OrganizationUnitID:       formDTO.OrganizationUnitID,
		ExecutedByID:             formDTO.ExecutedByID,
		AdvisoryStatusID:         openStatus.ID,
		InternalStatusID:         formDTO.InternalStatusID,
		NotifyDate:               dates.NotifyDate,
		DeadlineDate:             dates.DeadlineDate,
		ResolveDate:              dates.ResolveDate,
		ReplyDate:                dates.ReplyDate,
		SAEDTicketNo:             formDTO.SAEDTicketNo,
		Attachment:               formDTO.Attachment,
		ReportLink:               formDTO.ReportLink,
		Notes:                    formDTO.Notes,
		CreatedBYID:              currentUserID,
	}

	if err := pas.calculateAndApplyRiskLevel(plannedAdvisory); err != nil {
		return 0, err
	}
	pas.calculateProcessingTime(plannedAdvisory)

	tx := pas.repo.BeginTransaction()

	id, err := pas.repo.Create(plannedAdvisory, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	if err := pas.createPlannedAdvisoryDocuments(id, currentUserID, formDTO.Documents, tx); err != nil {
		tx.Rollback()
		return 0, err
	}

	tx.Commit()

	go func() {
		detail, err := pas.GetFullDetailsByID(id)
		if err != nil {
			utility.LogError("Failed to fetch planned advisory for activity log", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}
		err = pas.activityLogSvc.LogActivity(
			models.CreateAction,
			models.ThreatIntelligencePlannedAdvisory{}.TableName(),
			id,
			detail,
			currentUserID,
			"threat_intelligence", "planned_advisory", "", nil,
		)
		if err != nil {
			utility.LogError("Failed to log planned advisory creation activity", map[string]interface{}{
				"error": err.Error(),
			})
		}
	}()

	return id, nil
}

func (pas *ThreatIntelligencePlannedAdvisoryService) Update(id uint, currentUserID uint, formDTO *tiDTOs.PlannedAdvisoryFormDTO) error {
	tx := pas.repo.BeginTransaction()

	oldDetail, err := pas.GetFullDetailsByID(id, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := pas.validatePlannedAdvisoryFormDTO(formDTO); err != nil {
		tx.Rollback()
		return err
	}

	dates, err := pas.parsePlannedAdvisoryDates(formDTO)
	if err != nil {
		tx.Rollback()
		return err
	}

	existing, err := pas.repo.GetByID(id, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	existing.CaseNumber = formDTO.CaseNumber
	existing.TitleEn = formDTO.TitleEn
	existing.TitleAr = formDTO.TitleAr
	existing.ReceiptDate = dates.ReceiptDate
	existing.CaseLink = formDTO.CaseLink
	existing.AdvisorySourceID = formDTO.AdvisorySourceID
	existing.CaseClassificationID = formDTO.CaseClassificationID
	existing.AdvisoryTypeID = formDTO.AdvisoryTypeID
	existing.ImportanceID = formDTO.ImportanceID
	existing.SecurityClassificationID = formDTO.SecurityClassificationID
	existing.ThreatType = formDTO.ThreatType
	existing.ThreatSeverityID = formDTO.ThreatSeverityID
	existing.LikelihoodID = formDTO.LikelihoodID
	existing.Description = formDTO.Description
	existing.AffectedAssets = formDTO.AffectedAssets
	existing.Tags = formDTO.Tags
	existing.IOCs = formDTO.IOCs
	existing.AdvisoryRecipientID = formDTO.AdvisoryRecipientID
	existing.OrganizationUnitID = formDTO.OrganizationUnitID
	existing.ExecutedByID = formDTO.ExecutedByID
	existing.InternalStatusID = formDTO.InternalStatusID
	existing.NotifyDate = dates.NotifyDate
	existing.DeadlineDate = dates.DeadlineDate
	existing.ResolveDate = dates.ResolveDate
	existing.ReplyDate = dates.ReplyDate
	existing.SAEDTicketNo = formDTO.SAEDTicketNo
	existing.Attachment = formDTO.Attachment
	existing.ReportLink = formDTO.ReportLink
	existing.Notes = formDTO.Notes
	existing.UpdatedBYID = &currentUserID

	deletedFileNames, err := pas.updatePlannedAdvisoryDocuments(id, currentUserID, formDTO.Documents, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := pas.calculateAndApplyRiskLevel(existing); err != nil {
		tx.Rollback()
		return err
	}

	pas.calculateProcessingTime(existing)

	if err := pas.repo.Update(existing, tx); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	storagePath := filepath.Join("threat-intelligence", "planned-advisories", strconv.FormatUint(uint64(id), 10), "documents")
	for _, fileName := range deletedFileNames {
		if err := utility.RemoveFile(fileName, storagePath); err != nil {
			utility.LogError("Failed to remove deleted planned advisory document", map[string]interface{}{
				"error": err.Error(),
				"file":  fileName,
			})
		}
	}

	go func() {
		newDetail, err := pas.GetFullDetailsByID(id)
		if err != nil {
			return
		}

		changedData := pas.buildChangedData(oldDetail, newDetail)
		if len(changedData) > 0 {
			err = pas.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.ThreatIntelligencePlannedAdvisory{}.TableName(),
				id,
				changedData,
				currentUserID,
				"threat_intelligence", "planned_advisory", "", nil,
			)
			if err != nil {
				utility.LogError("Failed to log planned advisory update activity", map[string]interface{}{
					"error": err.Error(),
				})
			}
		}
	}()

	return nil
}

func (pas *ThreatIntelligencePlannedAdvisoryService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	tx := pas.repo.BeginTransaction()

	var failedIds []uint
	var deletedRecords []struct {
		id    uint
		title string
	}

	for _, id := range ids {
		detail, err := pas.GetFullDetailsByID(id, tx)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		if len(detail.Documents) > 0 {
			var deleteDocumentIds []*uint
			for _, doc := range detail.Documents {
				if doc != nil {
					deleteDocumentIds = append(deleteDocumentIds, &doc.ID)
				}
			}
			if err := pas.repo.DeletePlannedAdvisoryDocuments(deleteDocumentIds, tx); err != nil {
				failedIds = append(failedIds, id)
				continue
			}
		}

		if err := pas.repo.Delete(id, tx); err != nil {
			failedIds = append(failedIds, id)
			continue
		}

		deletedRecords = append(deletedRecords, struct {
			id    uint
			title string
		}{id: id, title: detail.TitleEn})
	}

	tx.Commit()

	for _, record := range deletedRecords {
		storagePath := filepath.Join("threat-intelligence", "planned-advisories", strconv.FormatUint(uint64(record.id), 10), "documents")
		if err := utility.RemoveFolder(storagePath); err != nil {
			utility.LogError("Failed to remove planned advisory documents folder", map[string]interface{}{
				"error": err.Error(),
				"path":  storagePath,
			})
		}
	}

	go func() {
		for _, record := range deletedRecords {
			err := pas.activityLogSvc.LogActivity(
				models.DeleteAction,
				models.ThreatIntelligencePlannedAdvisory{}.TableName(),
				record.id,
				map[string]interface{}{"id": record.id, "title_en": record.title},
				currentUserID,
				"threat_intelligence", "planned_advisory", "", nil,
			)
			if err != nil {
				utility.LogError("Failed to log planned advisory deletion activity", map[string]interface{}{
					"error": err.Error(),
				})
			}
		}
	}()

	return failedIds, nil
}

func (pas *ThreatIntelligencePlannedAdvisoryService) UpdateStatus(id uint, currentUserID uint, formDTO *tiDTOs.PlannedAdvisoryStatusUpdateFormDTO) error {
	existing, err := pas.GetFullDetailsByID(id)
	if err != nil {
		return err
	}

	newStatus, err := pas.advisoryStatusSvc.GetByID(formDTO.AdvisoryStatusID)
	if err != nil {
		return utility.NewValidationError(fmt.Errorf("advisory status does not exist"), "not_exist", "advisory_status_id")
	}

	if existing.AdvisoryStatus != nil && existing.AdvisoryStatus.ID == newStatus.ID {
		return utility.NewValidationError(fmt.Errorf("advisory is already in the requested status"), "duplicate", "advisory_status_id")
	}

	oldStatusID := uint(0)
	if existing.AdvisoryStatus != nil {
		oldStatusID = existing.AdvisoryStatus.ID
	}

	tx := pas.repo.BeginTransaction()

	model, err := pas.repo.GetByID(id, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	model.AdvisoryStatusID = newStatus.ID
	model.UpdatedBYID = &currentUserID

	if err := pas.repo.Update(model, tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	go func() {
		changedData := map[string]interface{}{
			"advisory_status_id": map[string]interface{}{
				"old": oldStatusID,
				"new": newStatus.ID,
			},
		}

		logErr := pas.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.ThreatIntelligencePlannedAdvisory{}.TableName(),
			id,
			changedData,
			currentUserID,
			"threat_intelligence", "planned_advisory", "status_update", nil,
		)
		if logErr != nil {
			utility.LogError("Failed to log planned advisory status update activity", map[string]interface{}{
				"error": logErr.Error(),
			})
		}

		creator, err := pas.userSvc.GetByID(uint(existing.CreatedByUser.ID))
		if err != nil {
			utility.LogError("Failed to fetch planned advisory creator for notification", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		updater, err := pas.userSvc.GetByID(currentUserID)
		if err != nil {
			utility.LogError("Failed to fetch planned advisory status updater for notification", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		oldStatusTitle := ""
		if existing.AdvisoryStatus != nil {
			oldStatusTitle = existing.AdvisoryStatus.TitleEn
		}

		entity := "planned-advisory"
		action := "status-update"
		urlParam := fmt.Sprintf(`{"planned_advisory_id":%d}`, id)

		dbNotificationDTO := &commonDTOs.DBNotificationDTO{
			UserIDs:              []uint{uint(creator.ID)},
			ModuleIdentifier:     "threat_intelligence",
			NotificationType:     string(commonDTOs.ActionType),
			ReferenceNumber:      existing.ReferenceNumber,
			HeadingEn:            "Planned Advisory Status Update",
			HeadingAr:            "تحديث حالة التنبيه المخطط",
			TitleEn:              existing.TitleEn,
			TitleAr:              existing.TitleAr,
			Entity:               &entity,
			Action:               &action,
			URLParam:             urlParam,
			QueryParam:           "",
			IsSystemNotification: false,
			CreatedByID:          &currentUserID,
		}

		emailTo := []string{creator.Email}
		subject := fmt.Sprintf("Planned Advisory %s: Status Updated", existing.ReferenceNumber)
		appURL := utility.GetAppURL()
		emailData := map[string]interface{}{
			"FirstNameEn":     creator.FirstNameEn,
			"LastNameEn":      creator.LastNameEn,
			"ReferenceNumber": existing.ReferenceNumber,
			"Title":           existing.TitleEn,
			"OldStatus":       oldStatusTitle,
			"NewStatus":       newStatus.TitleEn,
			"UpdatedBy":       fmt.Sprintf("%s %s", updater.FirstNameEn, updater.LastNameEn),
			"UpdatedOn":       time.Now().Format("2006-01-02"),
			"RedirectURL":     fmt.Sprintf("%s/threat-intelligence/planned-advisories/%d", appURL, id),
		}

		emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
			EmailTo:              emailTo,
			Subject:              subject,
			TemplateRelativePath: "threat-intelligence/templates/email/planned_advisory_status_update_mail.html",
			EmailData:            emailData,
		}

		if err := pas.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "planned_advisory"); err != nil {
			utility.LogError("Failed to send planned advisory status update notification", map[string]interface{}{
				"error": err.Error(),
			})
		}
	}()

	return nil
}

func (pas *ThreatIntelligencePlannedAdvisoryService) validatePlannedAdvisoryFormDTO(formDTO *tiDTOs.PlannedAdvisoryFormDTO) error {
	if _, err := pas.advisorySourceSvc.GetByID(formDTO.AdvisorySourceID); err != nil {
		return utility.NewValidationError(fmt.Errorf("advisory source does not exist"), "not_exist", "advisory_source")
	}

	if _, err := pas.advisoryClassificationSvc.GetByID(formDTO.CaseClassificationID); err != nil {
		return utility.NewValidationError(fmt.Errorf("case classification does not exist"), "not_exist", "case_classification")
	}

	if _, err := pas.advisoryTypeSvc.GetByID(formDTO.AdvisoryTypeID); err != nil {
		return utility.NewValidationError(fmt.Errorf("advisory type does not exist"), "not_exist", "advisory_type")
	}

	if _, err := pas.advisoryImportanceSvc.GetByID(formDTO.ImportanceID); err != nil {
		return utility.NewValidationError(fmt.Errorf("importance does not exist"), "not_exist", "importance")
	}

	if _, err := pas.advisorySecurityClassificationSvc.GetByID(formDTO.SecurityClassificationID); err != nil {
		return utility.NewValidationError(fmt.Errorf("security classification does not exist"), "not_exist", "security_classification")
	}

	if formDTO.InternalStatusID != nil {
		if _, err := pas.advisoryInternalStatusSvc.GetByID(*formDTO.InternalStatusID); err != nil {
			return utility.NewValidationError(fmt.Errorf("internal status does not exist"), "not_exist", "internal_status")
		}
	}

	if formDTO.AdvisoryRecipientID > 0 {
		if _, err := pas.userSvc.GetByID(formDTO.AdvisoryRecipientID); err != nil {
			return utility.NewValidationError(fmt.Errorf("advisory recipient does not exist"), "not_exist", "advisory_recipient")
		}
	}

	if formDTO.ExecutedByID != nil {
		if _, err := pas.userSvc.GetByID(*formDTO.ExecutedByID); err != nil {
			return utility.NewValidationError(fmt.Errorf("executed by user does not exist"), "not_exist", "executed_by")
		}
	}

	if formDTO.OrganizationUnitID > 0 {
		if _, err := pas.organizationUnitSvc.GetByID(formDTO.OrganizationUnitID); err != nil {
			return utility.NewValidationError(fmt.Errorf("organization unit does not exist"), "not_exist", "organization_unit")
		}
	}

	if formDTO.ThreatSeverityID != nil {
		if _, err := pas.severitySvc.GetByID(*formDTO.ThreatSeverityID); err != nil {
			return utility.NewValidationError(fmt.Errorf("threat severity does not exist"), "not_exist", "threat_severity")
		}
	}

	if formDTO.LikelihoodID != nil {
		if _, err := pas.likelihoodSvc.GetByID(*formDTO.LikelihoodID); err != nil {
			return utility.NewValidationError(fmt.Errorf("likelihood does not exist"), "not_exist", "likelihood")
		}
	}

	if (formDTO.ThreatSeverityID != nil && formDTO.LikelihoodID == nil) || (formDTO.ThreatSeverityID == nil && formDTO.LikelihoodID != nil) {
		return utility.NewValidationError(fmt.Errorf("both threat severity and likelihood are required for risk calculation"), "required", "risk_calculation")
	}

	return nil
}

func (pas *ThreatIntelligencePlannedAdvisoryService) parsePlannedAdvisoryDates(formDTO *tiDTOs.PlannedAdvisoryFormDTO) (*tiDTOs.PlannedAdvisoryParsedDates, error) {
	parsed := &tiDTOs.PlannedAdvisoryParsedDates{}
	var err error

	receiptDate, err := utility.ParseDate(formDTO.ReceiptDate)
	if err != nil {
		return nil, utility.NewValidationError(fmt.Errorf("invalid receipt date"), "invalid_date", "receipt_date")
	}
	if receiptDate == nil {
		return nil, utility.NewValidationError(fmt.Errorf("receipt date is required"), "required", "receipt_date")
	}
	parsed.ReceiptDate = *receiptDate

	if formDTO.NotifyDate != nil && *formDTO.NotifyDate != "" {
		parsed.NotifyDate, err = utility.ParseDate(*formDTO.NotifyDate)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid notify date"), "invalid_date", "notify_date")
		}
	}

	if formDTO.DeadlineDate != nil && *formDTO.DeadlineDate != "" {
		parsed.DeadlineDate, err = utility.ParseDate(*formDTO.DeadlineDate)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid deadline date"), "invalid_date", "deadline_date")
		}
	}

	if formDTO.ResolveDate != nil && *formDTO.ResolveDate != "" {
		parsed.ResolveDate, err = utility.ParseDate(*formDTO.ResolveDate)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid resolve date"), "invalid_date", "resolve_date")
		}
	}

	if formDTO.ReplyDate != nil && *formDTO.ReplyDate != "" {
		parsed.ReplyDate, err = utility.ParseDate(*formDTO.ReplyDate)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid reply date"), "invalid_date", "reply_date")
		}
	}

	return parsed, nil
}

func (pas *ThreatIntelligencePlannedAdvisoryService) calculateAndApplyRiskLevel(plannedAdvisory *models.ThreatIntelligencePlannedAdvisory) error {
	if plannedAdvisory.ThreatSeverityID == nil || plannedAdvisory.LikelihoodID == nil {
		plannedAdvisory.RiskLevelID = nil
		plannedAdvisory.RiskLevelValueID = nil
		return nil
	}

	severity, err := pas.severitySvc.GetByID(*plannedAdvisory.ThreatSeverityID)
	if err != nil {
		return err
	}

	likelihood, err := pas.likelihoodSvc.GetByID(*plannedAdvisory.LikelihoodID)
	if err != nil {
		return err
	}

	score := pas.riskLevelValueSvc.EvaluateScore(
		&models.ThreatIntelligenceLikelihood{Score: likelihood.Score},
		&models.ThreatIntelligenceSeverity{Score: severity.Score},
		"multiplication",
	)

	riskLevelValueArray, err := pas.buildRiskLevelValueArray()
	if err != nil {
		return err
	}

	riskLevelValue, found := pas.riskLevelValueSvc.FindRiskLevelValue(riskLevelValueArray, score)
	if !found || riskLevelValue == nil {
		riskLevelValue, found = pas.riskLevelValueSvc.FindClosestRiskLevelValue(riskLevelValueArray, score)
		if !found || riskLevelValue == nil {
			return utility.NewValidationError(fmt.Errorf("no matching risk level value found for score %d", score), "not_exist", "risk_level_value")
		}
	}

	riskLevelID := riskLevelValue.ThreatIntelligenceRiskLevelID
	riskLevelValueID := riskLevelValue.ThreatIntelligenceRiskLevelValueID
	plannedAdvisory.RiskLevelID = &riskLevelID
	plannedAdvisory.RiskLevelValueID = &riskLevelValueID

	return nil
}

func (pas *ThreatIntelligencePlannedAdvisoryService) buildRiskLevelValueArray() ([]*tiDTOs.ThreatIntelligenceRiskLevelDataDTO, error) {
	riskLevels, err := pas.riskLevelSvc.GetAll(&tiDTOs.ThreatIntelligenceRiskLevelListQueryParamsDTO{Status: "active"})
	if err != nil {
		return nil, err
	}

	var riskLevelValueArray []*tiDTOs.ThreatIntelligenceRiskLevelDataDTO
	for _, riskLevel := range riskLevels {
		riskLevelDetails, err := pas.riskLevelSvc.GetByID(riskLevel.ID)
		if err != nil {
			continue
		}

		for _, value := range riskLevelDetails.ThreatIntelligenceRiskLevelValue {
			riskLevelValueArray = append(riskLevelValueArray, &tiDTOs.ThreatIntelligenceRiskLevelDataDTO{
				ThreatIntelligenceRiskLevelID:      riskLevel.ID,
				ThreatIntelligenceRiskLevelValueID: value.ID,
				RiskLevelTitleEN:                   riskLevel.TitleEn,
				RiskLevelTitleAR:                   riskLevel.TitleAr,
				RiskLevelType:                      riskLevel.Type,
				ColorCode:                          riskLevel.ColorCode,
				Value:                              value.Value,
			})
		}
	}

	return riskLevelValueArray, nil
}

func (pas *ThreatIntelligencePlannedAdvisoryService) calculateProcessingTime(plannedAdvisory *models.ThreatIntelligencePlannedAdvisory) {
	if plannedAdvisory.ResolveDate == nil {
		plannedAdvisory.ProcessingTime = nil
		return
	}

	elapsed := (*plannedAdvisory.ResolveDate).Sub(plannedAdvisory.ReceiptDate)
	plannedAdvisory.ProcessingTime = formatDuration(elapsed)
}

func formatDuration(d time.Duration) *string {
	days := int(d.Hours()) / 24
	hours := int(d.Hours()) % 24
	minutes := int(d.Minutes()) % 60

	result := ""
	if days > 0 {
		result += fmt.Sprintf("%d day", days)
		if days > 1 {
			result += "s"
		}
		result += " "
	}
	if hours > 0 {
		result += fmt.Sprintf("%d hour", hours)
		if hours > 1 {
			result += "s"
		}
		if minutes > 0 {
			result += " "
		}
	}
	if minutes > 0 || (days == 0 && hours == 0) {
		result += fmt.Sprintf("%d minute", minutes)
		if minutes != 1 {
			result += "s"
		}
	}

	if len(result) > 0 && result[len(result)-1] == ' ' {
		result = result[:len(result)-1]
	}

	return &result
}

func (pas *ThreatIntelligencePlannedAdvisoryService) createPlannedAdvisoryDocuments(plannedAdvisoryID uint, currentUserID uint, docs []*commonDTOs.DocumentFormDTO, tx utility.DatabaseTransaction) error {
	if len(docs) == 0 {
		return nil
	}

	var documents []*models.ThreatIntelligencePlannedAdvisoryDocument
	for _, doc := range docs {
		encryptedFilename, err := utility.EncryptData([]byte(doc.FileName))
		if err != nil {
			return err
		}

		url := utility.GenerateDocumentURL("threat-intelligence", "planned-advisories", strconv.FormatUint(uint64(plannedAdvisoryID), 10), "documents")

		documents = append(documents, &models.ThreatIntelligencePlannedAdvisoryDocument{
			ThreatIntelligencePlannedAdvisoryID: plannedAdvisoryID,
			URL:                                 url,
			Name:                                doc.FileName,
			Size:                                doc.Size,
			Token:                               encryptedFilename,
			OriginalName:                        doc.OriginalName,
			Type:                                doc.Type,
			CreatedByID:                         currentUserID,
		})
	}

	if err := pas.repo.CreatePlannedAdvisoryDocuments(documents, tx); err != nil {
		return err
	}

	storagePath := filepath.Join("threat-intelligence", "planned-advisories", strconv.FormatUint(uint64(plannedAdvisoryID), 10), "documents")
	for _, doc := range documents {
		if err := utility.StoreFile(doc.Name, storagePath); err != nil {
			return err
		}
	}

	return nil
}

func (pas *ThreatIntelligencePlannedAdvisoryService) updatePlannedAdvisoryDocuments(plannedAdvisoryID uint, currentUserID uint, docs []*commonDTOs.DocumentFormDTO, tx utility.DatabaseTransaction) ([]string, error) {
	var deletedFileNames []string
	if len(docs) == 0 {
		return deletedFileNames, nil
	}

	var documents []*models.ThreatIntelligencePlannedAdvisoryDocument
	var deleteDocumentIds []*uint

	for _, doc := range docs {
		if (doc.IsNewFile != nil && *doc.IsNewFile) || (doc.ID == nil && doc.FileName != "") {
			encryptedFilename, err := utility.EncryptData([]byte(doc.FileName))
			if err != nil {
				return deletedFileNames, err
			}

			url := utility.GenerateDocumentURL("threat-intelligence", "planned-advisories", strconv.FormatUint(uint64(plannedAdvisoryID), 10), "documents")

			documents = append(documents, &models.ThreatIntelligencePlannedAdvisoryDocument{
				ThreatIntelligencePlannedAdvisoryID: plannedAdvisoryID,
				URL:                                 url,
				Name:                                doc.FileName,
				Size:                                doc.Size,
				Token:                               encryptedFilename,
				OriginalName:                        doc.OriginalName,
				Type:                                doc.Type,
				CreatedByID:                         currentUserID,
			})
		} else if doc.IsDeletedFile != nil && *doc.IsDeletedFile && doc.ID != nil {
			existingDocument, err := pas.repo.GetPlannedAdvisoryDocumentByID(plannedAdvisoryID, *doc.ID, tx)
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return deletedFileNames, utility.NewResourceNotFoundError(err, *doc.ID)
				}
				return deletedFileNames, err
			}
			deleteDocumentIds = append(deleteDocumentIds, doc.ID)
			deletedFileNames = append(deletedFileNames, existingDocument.Name)
		}
	}

	if len(documents) > 0 {
		if err := pas.repo.CreatePlannedAdvisoryDocuments(documents, tx); err != nil {
			return deletedFileNames, err
		}
		storagePath := filepath.Join("threat-intelligence", "planned-advisories", strconv.FormatUint(uint64(plannedAdvisoryID), 10), "documents")
		for _, doc := range documents {
			if err := utility.StoreFile(doc.Name, storagePath); err != nil {
				return deletedFileNames, err
			}
		}
	}

	if len(deleteDocumentIds) > 0 {
		if err := pas.repo.DeletePlannedAdvisoryDocuments(deleteDocumentIds, tx); err != nil {
			return deletedFileNames, err
		}
	}

	return deletedFileNames, nil
}

func (pas *ThreatIntelligencePlannedAdvisoryService) GetDocumentPath(plannedAdvisoryID uint, token string) (string, error) {
	documents, err := pas.repo.GetDocumentsByPlannedAdvisoryID(plannedAdvisoryID)
	if err != nil {
		return "", err
	}
	if len(documents) == 0 {
		return "", fmt.Errorf("no documents found for the planned advisory")
	}

	decryptedFilenameBytes, err := utility.DecryptData(token)
	if err != nil {
		return "", fmt.Errorf("resource not found")
	}

	decryptedFilename := string(decryptedFilenameBytes)
	if decryptedFilename == "" {
		return "", fmt.Errorf("resource not found")
	}

	storagePath := filepath.Join("threat-intelligence", "planned-advisories", strconv.FormatUint(uint64(plannedAdvisoryID), 10), "documents")
	baseStoragePath := config.GetEnv("STORAGE_PATH", "")
	fullPath := filepath.Join(baseStoragePath, storagePath, decryptedFilename)

	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return "", fmt.Errorf("file not found")
	} else if err != nil {
		return "", fmt.Errorf("error accessing file: %w", err)
	}

	return fullPath, nil
}

func (pas *ThreatIntelligencePlannedAdvisoryService) ExportPlannedAdvisory(data []*tiDTOs.PlannedAdvisoryListDTO, currentUserID uint) (*os.File, error) {
	userSetting, err := pas.userSettingSvc.GetUserSettingsByUserID(currentUserID)
	if err != nil {
		return nil, err
	}

	var headers []string
	if userSetting.Language == "ar" {
		headers = []string{
			"SL No.",
			"Reference Number",
			"Case Number",
			"Title",
			"Receipt Date",
			"Case Link",
			"Threat Type",
			"Advisory Source",
			"Case Classification",
			"Advisory Type",
			"Importance",
			"Security Classification",
			"Advisory Recipient",
			"Organization Unit",
			"Executed By ID",
			"Advisory Status",
			"Internal Status",
			"Threat Severity Score",
			"Likelihood Score",
			"Risk Level",
			"Risk Level Value ID",
			"Description",
			"Affected Assets",
			"Tags",
			"IOCs",
			"Notify Date",
			"Deadline Date",
			"Resolve Date",
			"Reply Date",
			"SAED Ticket No",
			"Report Link",
			"Notes",
			"Created At",
			"Created By",
			"Updated At",
			"Updated By",
		}
	} else {
		headers = []string{
			"SL No.",
			"Reference Number",
			"Case Number",
			"Title",
			"Receipt Date",
			"Case Link",
			"Threat Type",
			"Advisory Source",
			"Case Classification",
			"Advisory Type",
			"Importance",
			"Security Classification",
			"Advisory Recipient",
			"Organization Unit",
			"Executed By ID",
			"Advisory Status",
			"Internal Status",
			"Threat Severity Score",
			"Likelihood Score",
			"Risk Level",
			"Risk Level Value ID",
			"Description",
			"Affected Assets",
			"Tags",
			"IOCs",
			"Notify Date",
			"Deadline Date",
			"Resolve Date",
			"Reply Date",
			"SAED Ticket No",
			"Report Link",
			"Notes",
			"Created At",
			"Created By",
			"Updated At",
			"Updated By",
		}
	}

	rows := make([][]interface{}, len(data))
	for i, advisory := range data {
		var title, advisorySource, caseClassification, advisoryType, importance, securityClassification string
		var advisoryRecipient, organizationUnit, advisoryStatus, internalStatus, riskLevel, createdBy, updatedBy string
		var threatSeverityScore, likelihoodScore, executedByID, riskLevelValueID string

		if userSetting.Language == "ar" {
			title = advisory.TitleAr
			advisorySource = advisory.AdvisorySourceTitleAr
			caseClassification = advisory.CaseClassificationTitleAr
			advisoryType = advisory.AdvisoryTypeTitleAr
			importance = advisory.ImportanceTitleAr
			securityClassification = advisory.SecurityClassificationTitleAr
			if advisory.AdvisoryRecipientFirstNameAr != nil && advisory.AdvisoryRecipientLastNameAr != nil {
				advisoryRecipient = *advisory.AdvisoryRecipientFirstNameAr + " " + *advisory.AdvisoryRecipientLastNameAr
			}
			organizationUnit = advisory.OrganizationUnitTitleAr
			advisoryStatus = advisory.AdvisoryStatusTitleAr
			if advisory.InternalStatusTitleAr != nil {
				internalStatus = *advisory.InternalStatusTitleAr
			}
			if advisory.RiskLevelTitleAr != nil {
				riskLevel = *advisory.RiskLevelTitleAr
			}
			createdBy = advisory.CreatedUserFirstNameAr + " " + advisory.CreatedUserLastNameAr
			if advisory.UpdatedUserFirstNameAr != nil && advisory.UpdatedUserLastNameAr != nil {
				updatedBy = *advisory.UpdatedUserFirstNameAr + " " + *advisory.UpdatedUserLastNameAr
			}
		} else {
			title = advisory.TitleEn
			advisorySource = advisory.AdvisorySourceTitleEn
			caseClassification = advisory.CaseClassificationTitleEn
			advisoryType = advisory.AdvisoryTypeTitleEn
			importance = advisory.ImportanceTitleEn
			securityClassification = advisory.SecurityClassificationTitleEn
			if advisory.AdvisoryRecipientFirstNameEn != nil && advisory.AdvisoryRecipientLastNameEn != nil {
				advisoryRecipient = *advisory.AdvisoryRecipientFirstNameEn + " " + *advisory.AdvisoryRecipientLastNameEn
			}
			organizationUnit = advisory.OrganizationUnitTitleEn
			advisoryStatus = advisory.AdvisoryStatusTitleEn
			if advisory.InternalStatusTitleEn != nil {
				internalStatus = *advisory.InternalStatusTitleEn
			}
			if advisory.RiskLevelTitleEn != nil {
				riskLevel = *advisory.RiskLevelTitleEn
			}
			createdBy = advisory.CreatedUserFirstNameEn + " " + advisory.CreatedUserLastNameEn
			if advisory.UpdatedUserFirstNameEn != nil && advisory.UpdatedUserLastNameEn != nil {
				updatedBy = *advisory.UpdatedUserFirstNameEn + " " + *advisory.UpdatedUserLastNameEn
			}
		}

		if advisory.ThreatSeverityScore != nil {
			threatSeverityScore = strconv.Itoa(*advisory.ThreatSeverityScore)
		}
		if advisory.LikelihoodScore != nil {
			likelihoodScore = strconv.Itoa(*advisory.LikelihoodScore)
		}
		if advisory.ExecutedByID != nil {
			executedByID = strconv.FormatUint(uint64(*advisory.ExecutedByID), 10)
		}
		if advisory.RiskLevelValueID != nil {
			riskLevelValueID = strconv.FormatUint(uint64(*advisory.RiskLevelValueID), 10)
		}

		rows[i] = []interface{}{
			i + 1,
			advisory.ReferenceNumber,
			advisory.CaseNumber,
			title,
			utility.GetStringValue(advisory.ReceiptDate),
			advisory.CaseLink,
			utility.GetStringValue(advisory.ThreatType),
			advisorySource,
			caseClassification,
			advisoryType,
			importance,
			securityClassification,
			advisoryRecipient,
			organizationUnit,
			executedByID,
			advisoryStatus,
			internalStatus,
			threatSeverityScore,
			likelihoodScore,
			riskLevel,
			riskLevelValueID,
			utility.GetStringValue(advisory.Description),
			utility.GetStringValue(advisory.AffectedAssets),
			utility.GetStringValue(advisory.Tags),
			utility.GetStringValue(advisory.IOCs),
			utility.GetStringValue(advisory.NotifyDate),
			utility.GetStringValue(advisory.DeadlineDate),
			utility.GetStringValue(advisory.ResolveDate),
			utility.GetStringValue(advisory.ReplyDate),
			utility.GetStringValue(advisory.SAEDTicketNo),
			utility.GetStringValue(advisory.ReportLink),
			utility.GetStringValue(advisory.Notes),
			advisory.CreatedAt,
			createdBy,
			utility.GetStringValue(advisory.UpdatedAt),
			updatedBy,
		}
	}

	return utility.GenerateExcelFile(headers, rows)
}

func (pas *ThreatIntelligencePlannedAdvisoryService) buildChangedData(oldDetail, newDetail *tiDTOs.PlannedAdvisoryDetailDTO) map[string]interface{} {
	changedData := make(map[string]interface{})

	if oldDetail.CaseNumber != newDetail.CaseNumber {
		changedData["case_number"] = map[string]interface{}{"from": oldDetail.CaseNumber, "to": newDetail.CaseNumber}
	}
	if oldDetail.TitleEn != newDetail.TitleEn {
		changedData["title_en"] = map[string]interface{}{"from": oldDetail.TitleEn, "to": newDetail.TitleEn}
	}
	if oldDetail.TitleAr != newDetail.TitleAr {
		changedData["title_ar"] = map[string]interface{}{"from": oldDetail.TitleAr, "to": newDetail.TitleAr}
	}
	if oldDetail.ReceiptDateString != newDetail.ReceiptDateString {
		changedData["receipt_date"] = map[string]interface{}{"from": oldDetail.ReceiptDateString, "to": newDetail.ReceiptDateString}
	}
	if oldDetail.CaseLink != newDetail.CaseLink {
		changedData["case_link"] = map[string]interface{}{"from": oldDetail.CaseLink, "to": newDetail.CaseLink}
	}
	if oldDetail.RiskLevelValueID != newDetail.RiskLevelValueID {
		changedData["risk_level_value_id"] = map[string]interface{}{"from": oldDetail.RiskLevelValueID, "to": newDetail.RiskLevelValueID}
	}
	if oldDetail.ProcessingTime != newDetail.ProcessingTime {
		changedData["processing_time"] = map[string]interface{}{"from": oldDetail.ProcessingTime, "to": newDetail.ProcessingTime}
	}
	if oldDetail.ResolveDateString != newDetail.ResolveDateString {
		changedData["resolve_date"] = map[string]interface{}{"from": oldDetail.ResolveDateString, "to": newDetail.ResolveDateString}
	}

	return changedData
}
