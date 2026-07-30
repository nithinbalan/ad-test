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
	assetServices "jethur/grc/asset/services"
	commonDTOs "jethur/grc/common/dtos"
	commonServices "jethur/grc/common/services"
	complianceServices "jethur/grc/compliance/services"
	logServices "jethur/grc/logs/services"
	masterServices "jethur/grc/masters/services"
	organizationServices "jethur/grc/organization/services"
	riskDTOs "jethur/grc/risk/dtos"
	riskRepositories "jethur/grc/risk/repositories"
	settingsDTOs "jethur/grc/settings/dtos"
	settingsServices "jethur/grc/settings/services"

	"jethur/grc/models"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type RiskRegisterServiceInterface interface {
	SetRiskWorkflowService(svc RiskWorkflowServiceInterface)
	GetAll(params *riskDTOs.RiskRegisterListQueryParamsDTO) (*riskDTOs.RiskRegisterListResponseDTO, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*riskDTOs.RiskRegisterDetailDTO, error)
	GetByID(id uint, tx ...utility.DatabaseTransaction) (*models.RiskRegister, error)
	Create(currentUserID uint, riskRegisterFormDTO *riskDTOs.RiskRegisterFormDTO) (uint, error)
	Update(id uint, currentUserID uint, risk *riskDTOs.RiskRegisterFormDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	ForceDelete(ids []uint, currentUserID uint) ([]uint, error)
	UpdateResponseStrategy(riskID uint, currentUserID uint, riskResponseStrategyFormDTO *riskDTOs.RiskResponseStrategyFormDTO, tx ...utility.DatabaseTransaction) error
	DeleteResponseStrategy(riskID uint, currentUserID uint) error
	GetControlPerformance(riskID uint) (map[string]float64, error)
	UpdateRiskRating(riskID uint, currentUserID uint, riskEvaluationDetails *models.RiskEvaluation, tx utility.DatabaseTransaction, evaluationType string) error
	ClearRiskRating(riskID uint, updateDTO *riskDTOs.RiskRegisterRiskRatingUpdateDTO, tx utility.DatabaseTransaction) error

	GetDocumentPath(id uint, token string) (string, error)

	ExportRiskRegisterPDF(id uint, currentUserID uint) (string, error)

	ExportRiskRegister(data []*riskDTOs.RiskRegisterListDTO, currentUserID uint) (*os.File, error)

	SetRiskInherentEvaluationService(evaluationSvc RiskInherentEvaluationServiceInterface)
	SetRiskResidualEvaluationService(evaluationSvc RiskResidualEvaluationServiceInterface)

	GenerateRiskTemplate() (*os.File, error)
	ValidateRiskTemplate(file io.Reader) ([]map[string]interface{}, error)
	ValidateImportRisks(input []map[string]interface{}) ([]map[string]interface{}, error)
	ImportRisks(input []map[string]interface{}, currentUserID uint) ([]uint, error)

	GenerateReferenceNumber() (string, error)

	GetAllIdentifiedRisks(params *riskDTOs.RiskRegisterListQueryParamsDTO) (*riskDTOs.RiskRegisterListResponseDTO, error)

	UpdateRisk(risk *models.RiskRegister, tx ...utility.DatabaseTransaction) error
}

// Service and implementations
type RiskRegisterService struct {
	repo                      riskRepositories.RiskRegisterRepositoryInterface
	activityLogSvc            logServices.ActivityLogServiceInterface
	riskCategorySvc           RiskCategoryServiceInterface
	riskSourceSvc             RiskSourceServiceInterface
	userSvc                   organizationServices.UserServiceInterface
	riskStatusSvc             RiskStatusServiceInterface
	reviewFrequencySvc        masterServices.ReviewFrequencyServiceInterface
	controlSvc                organizationServices.ControlServiceInterface
	complianceRequirementSvc  complianceServices.ComplianceRequirementServiceInterface
	processSvc                organizationServices.ProcessServiceInterface
	strategicObjectiveSvc     organizationServices.StrategicObjectiveServiceInterface
	slaRegisterSvc            complianceServices.SLARegisterServiceInterface
	riskMatrixSvc             RiskMatrixServiceInterface
	riskTimelineSvc           RiskTimelineServiceInterface
	userSettingSvc            organizationServices.UserSettingServiceInterface
	notificationSvc           commonServices.NotificationServiceInterface
	assetRegisterSvc          assetServices.AssetRegisterServiceInterface
	riskInherentEvaluationSvc RiskInherentEvaluationServiceInterface
	riskSettingsSvc           RiskSettingsServiceInterface
	organizationUnitSvc       settingsServices.OrganizationUnitServiceInterface
	organizationUnitsSvc      organizationServices.OrganizationUnitsServiceInterface
	riskMatrixLikelihoodSvc   RiskMatrixLikelihoodServiceInterface
	riskMatrixImpactSvc       RiskMatrixImpactServiceInterface
	riskResponseStrategySvc   RiskResponseStrategyServiceInterface
	riskResidualEvaluationSvc RiskResidualEvaluationServiceInterface
	riskWorkflowSvc           RiskWorkflowServiceInterface
	severityLevelSvc          masterServices.SeverityLevelServiceInterface
	riskTemplateSettingSvc    settingsServices.RiskTemplateSettingServiceInterface
	riskReportRepo            riskRepositories.RiskReportRepositoryInterface
	riskResponseSvc           RiskResponseServiceInterface
}

func NewRiskRegisterService(
	repo riskRepositories.RiskRegisterRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	riskCategorySvc RiskCategoryServiceInterface,
	riskSourceSvc RiskSourceServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	riskStatusSvc RiskStatusServiceInterface,
	reviewFrequencySvc masterServices.ReviewFrequencyServiceInterface,
	controlSvc organizationServices.ControlServiceInterface,
	complianceRequirementSvc complianceServices.ComplianceRequirementServiceInterface,
	processSvc organizationServices.ProcessServiceInterface,
	strategicObjectiveSvc organizationServices.StrategicObjectiveServiceInterface,
	slaRegisterSvc complianceServices.SLARegisterServiceInterface,
	riskMatrixSvc RiskMatrixServiceInterface,
	riskTimelineSvc RiskTimelineServiceInterface,
	userSettingSvc organizationServices.UserSettingServiceInterface,
	notificationSvc commonServices.NotificationServiceInterface,
	assetRegisterSvc assetServices.AssetRegisterServiceInterface,
	riskInherentEvaluationSvc RiskInherentEvaluationServiceInterface,
	riskSettingsSvc RiskSettingsServiceInterface,
	organizationUnitSvc settingsServices.OrganizationUnitServiceInterface,
	organizationUnitsSvc organizationServices.OrganizationUnitsServiceInterface,
	riskMatrixLikelihoodSvc RiskMatrixLikelihoodServiceInterface,
	riskMatrixImpactSvc RiskMatrixImpactServiceInterface,
	riskResponseStrategySvc RiskResponseStrategyServiceInterface,
	riskResidualEvaluationSvc RiskResidualEvaluationServiceInterface,
	severityLevelSvc masterServices.SeverityLevelServiceInterface,
	riskReportRepo riskRepositories.RiskReportRepositoryInterface,
	riskTemplateSettingSvc settingsServices.RiskTemplateSettingServiceInterface,
	riskResponseSvc RiskResponseServiceInterface,
) RiskRegisterServiceInterface {
	return &RiskRegisterService{
		repo:                      repo,
		activityLogSvc:            activityLogSvc,
		riskCategorySvc:           riskCategorySvc,
		riskSourceSvc:             riskSourceSvc,
		userSvc:                   userSvc,
		riskStatusSvc:             riskStatusSvc,
		reviewFrequencySvc:        reviewFrequencySvc,
		controlSvc:                controlSvc,
		complianceRequirementSvc:  complianceRequirementSvc,
		processSvc:                processSvc,
		strategicObjectiveSvc:     strategicObjectiveSvc,
		slaRegisterSvc:            slaRegisterSvc,
		riskMatrixSvc:             riskMatrixSvc,
		riskTimelineSvc:           riskTimelineSvc,
		userSettingSvc:            userSettingSvc,
		notificationSvc:           notificationSvc,
		assetRegisterSvc:          assetRegisterSvc,
		riskInherentEvaluationSvc: riskInherentEvaluationSvc,
		riskSettingsSvc:           riskSettingsSvc,
		organizationUnitSvc:       organizationUnitSvc,
		organizationUnitsSvc:      organizationUnitsSvc,
		riskMatrixLikelihoodSvc:   riskMatrixLikelihoodSvc,
		riskMatrixImpactSvc:       riskMatrixImpactSvc,
		riskResponseStrategySvc:   riskResponseStrategySvc,
		riskResidualEvaluationSvc: riskResidualEvaluationSvc,
		severityLevelSvc:          severityLevelSvc,
		riskReportRepo:            riskReportRepo,
		riskTemplateSettingSvc:    riskTemplateSettingSvc,
		riskResponseSvc:           riskResponseSvc,
	}
}

func (rrs *RiskRegisterService) SetRiskWorkflowService(svc RiskWorkflowServiceInterface) {
	rrs.riskWorkflowSvc = svc
}

func (rrs *RiskRegisterService) GetAll(params *riskDTOs.RiskRegisterListQueryParamsDTO) (*riskDTOs.RiskRegisterListResponseDTO, error) {
	return rrs.repo.GetAll(params)
}

func (rrs RiskRegisterService) GetByID(id uint, tx ...utility.DatabaseTransaction) (*models.RiskRegister, error) {
	risk, err := rrs.repo.GetByID(id, tx...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return risk, nil
}

func (rrs RiskRegisterService) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*riskDTOs.RiskRegisterDetailDTO, error) {
	var riskDTO riskDTOs.RiskRegisterDetailDTO

	risk, err := rrs.repo.GetFullDetailsByID(id, tx...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	err = mapstructure.Decode(risk, &riskDTO)
	if err != nil {
		return nil, err
	}

	// Map the latest auditable item
	if len(risk.AuditableItems) > 0 {
		latestAuditableItem := risk.AuditableItems[0]
		var auditableItemDTO riskDTOs.AuditableItemMinimalRiskDetailDTO
		err := mapstructure.Decode(latestAuditableItem, &auditableItemDTO)
		if err != nil {
			return nil, err
		}
		riskDTO.AuditableItem = &auditableItemDTO
		riskDTO.AuditableItem.CreatedAt = latestAuditableItem.CreatedAt
		riskDTO.AuditableItem.UpdatedAt = latestAuditableItem.UpdatedAt
	}

	// Format date fields
	if risk.ReportedAt != nil {
		reportedAt := risk.ReportedAt.Format("2006-01-02")
		riskDTO.ReportedAtString = &reportedAt
	}

	if risk.RegistrationDate != nil {
		registrationDate := risk.RegistrationDate.Format("2006-01-02")
		riskDTO.RegistrationDateString = &registrationDate
	}

	if risk.NextReviewDate != nil {
		nextReviewDate := risk.NextReviewDate.Format("2006-01-02")
		riskDTO.NextReviewDateString = &nextReviewDate
	}

	if risk.SubmittedAt != nil {
		submittedAt := risk.SubmittedAt.Format("2006-01-02")
		riskDTO.SubmittedAtString = &submittedAt
	}

	riskDTO.CreatedAt = risk.CreatedAt
	riskDTO.UpdatedAt = risk.UpdatedAt

	return &riskDTO, nil
}

func (rrs *RiskRegisterService) Create(currentUserID uint, riskRegisterFormDTO *riskDTOs.RiskRegisterFormDTO) (uint, error) {
	tx := rrs.repo.BeginTransaction()

	err := rrs.riskWorkflowSvc.ValidateWorkflowBeforeCreate("risk_register", tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	if err := rrs.validateRiskRegisterFormDTO(riskRegisterFormDTO); err != nil {
		return 0, err
	}

	dates, err := rrs.parseRiskRegisterDates(riskRegisterFormDTO)
	if err != nil {
		return 0, err
	}

	referenceNumber, err := rrs.repo.GenerateReferenceNumber()
	if err != nil {
		return 0, err
	}

	// responsibleUnitID, err := rrs.userSvc.GetUserPrimaryUnitID(uint(riskRegisterFormDTO.ResponsibleUserID))
	// if err != nil {
	// 	return 0, utility.NewValidationError(fmt.Errorf("organization unit not set for responsible user"), "not_exist", "risk responsible unit")
	// }

	riskMatrix, err := rrs.riskMatrixSvc.GetByType("erm")
	if err != nil {
		return 0, err
	}

	riskStatus, err := rrs.riskStatusSvc.GetByType("identified")
	if err != nil {
		return 0, err
	}

	riskRegister := &models.RiskRegister{
		ReferenceNumber:   &referenceNumber,
		RiskCategoryID:    riskRegisterFormDTO.RiskCategoryID,
		RiskSourceID:      riskRegisterFormDTO.RiskSourceID,
		RiskMatrixID:      riskMatrix.ID,
		ResponsibleUnitID: riskRegisterFormDTO.ResponsibleUnitID,
		ResponsibleUserID: riskRegisterFormDTO.ResponsibleUserID,
		RiskStatusID:      riskStatus.ID,
		Title:             riskRegisterFormDTO.Title,
		Description:       riskRegisterFormDTO.Description,
		ReportedAt:        dates.ReportedAt,
		ReviewFrequencyID: riskRegisterFormDTO.ReviewFrequencyID,
		NextReviewDate:    dates.NextReviewDate,
		LastReviewNote:    riskRegisterFormDTO.LastReviewNote,
		SeverityLevelID:   riskRegisterFormDTO.SeverityLevelID,
		RiskCode:          riskRegisterFormDTO.RiskCode,
		DivisionID:        riskRegisterFormDTO.DivisionID,
		SectionID:         riskRegisterFormDTO.SectionID,
		SubSectionID:      riskRegisterFormDTO.SubSectionID,
		IsStrategic:       riskRegisterFormDTO.IsStrategic,
		IsOperational:     riskRegisterFormDTO.IsOperational,
		ExpectedCause:     riskRegisterFormDTO.ExpectedCause,
		Consequence:       riskRegisterFormDTO.Consequence,
		IsInternal:        riskRegisterFormDTO.IsInternal,
		IsExternal:        riskRegisterFormDTO.IsExternal,
		RegistrationDate:  dates.RegistrationDate,
		CreatedBYID:       currentUserID,
	}

	id, err := rrs.repo.Create(riskRegister, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	// Add workflow
	err = rrs.riskWorkflowSvc.CreateRiskWorkflow(id, "risk_register", currentUserID, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	// Add SLA Register
	var slaRegisters []*models.RiskSLA
	for _, slaRegisterID := range riskRegisterFormDTO.SLARegisterIDs {
		slaRegister := &models.RiskSLA{
			RiskRegisterID: id,
			SLARegisterID:  *slaRegisterID,
			CreatedByID:    currentUserID,
		}
		slaRegisters = append(slaRegisters, slaRegister)
	}

	if len(slaRegisters) > 0 {
		if err := rrs.repo.CreateRiskSLAs(slaRegisters, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	// Add timeline
	// todo
	timelineTypes, err := rrs.riskTimelineSvc.GetAllTypes()
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	// Get timeline statuses
	completedStatus, err := rrs.riskTimelineSvc.GetByStatusType("completed")
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	pendingStatus, err := rrs.riskTimelineSvc.GetByStatusType("pending")
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	var timelineFormDTOs []*riskDTOs.RiskTimelineFormDTO

	for _, timelineType := range timelineTypes {
		var (
			statusID      uint
			performedByID *uint
			currentTime   *time.Time
		)

		if timelineType.Type == "identified" {
			statusID = completedStatus.ID
			performedByID = &currentUserID
			now := time.Now()
			currentTime = &now
		} else {
			statusID = pendingStatus.ID
			performedByID = nil
			currentTime = nil
		}

		timelineFormDTO := &riskDTOs.RiskTimelineFormDTO{
			RiskRegisterID:     id,
			RiskTimelineTypeID: timelineType.ID,
			TimelineStatusID:   statusID,
			PerformedByID:      performedByID,
			PerformedAt:        currentTime,
		}
		timelineFormDTOs = append(timelineFormDTOs, timelineFormDTO)
	}
	if err := rrs.riskTimelineSvc.Create(timelineFormDTOs, tx); err != nil {
		tx.Rollback()
		return 0, err
	}

	var documents []*models.RiskDocument

	// Iterate over each document in the array
	for _, doc := range riskRegisterFormDTO.Documents {

		// Encrypt the filename
		dataToEncrypt := []byte(doc.FileName)
		encryptedFilename, err := utility.EncryptData(dataToEncrypt)
		if err != nil {
			tx.Rollback()
			return 0, err
		}

		// Construct the URL for the document
		url := utility.GenerateDocumentURL("risk", "risk-register", strconv.FormatUint(uint64(riskRegister.ID), 10), "documents")

		// Create a new RiskDocument instance
		document := &models.RiskDocument{
			RiskRegisterID: id,
			URL:            url,
			Name:           doc.FileName,
			Size:           doc.Size,
			Token:          encryptedFilename,
			OriginalName:   doc.OriginalName,
			Type:           doc.Type,
			CreatedByID:    currentUserID,
		}

		// Append the document to the slice
		documents = append(documents, document)
	}
	if len(documents) > 0 {
		// Save the RiskDocument
		if err := rrs.repo.CreateDocuments(documents, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
		for _, doc := range documents {
			storagePath := filepath.Join("risk", "risk-register", strconv.FormatUint(uint64(riskRegister.ID), 10), "documents")
			if err := utility.StoreFile(doc.Name, storagePath); err != nil {
				tx.Rollback()
				return 0, err
			}
		}
	}

	// update evaluation
	if riskRegisterFormDTO.InherentRiskMatrixLikelihoodID != nil && riskRegisterFormDTO.InherentRiskMatrixImpactID != nil {
		inherentEvaluationFormDTO := &riskDTOs.RiskInherentEvaluationFormDTO{
			RiskMatrixLikelihoodID: *riskRegisterFormDTO.InherentRiskMatrixLikelihoodID,
			RiskMatrixImpactID:     *riskRegisterFormDTO.InherentRiskMatrixImpactID,
		}

		_, err := rrs.riskInherentEvaluationSvc.Create(id, currentUserID, inherentEvaluationFormDTO, tx)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	if riskRegisterFormDTO.ResidualRiskMatrixLikelihoodID != nil && riskRegisterFormDTO.ResidualRiskMatrixImpactID != nil {
		residualEvaluationFormDTO := &riskDTOs.RiskResidualEvaluationFormDTO{
			RiskMatrixLikelihoodID: *riskRegisterFormDTO.ResidualRiskMatrixLikelihoodID,
			RiskMatrixImpactID:     *riskRegisterFormDTO.ResidualRiskMatrixImpactID,
		}

		_, err := rrs.riskResidualEvaluationSvc.Create(id, currentUserID, residualEvaluationFormDTO, tx)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	if riskRegisterFormDTO.RiskResponseStrategyID != nil {
		responseStrategyFormDTO := &riskDTOs.RiskResponseStrategyFormDTO{
			RiskResponseStrategyID: *riskRegisterFormDTO.RiskResponseStrategyID,
		}

		err := rrs.UpdateResponseStrategy(id, currentUserID, responseStrategyFormDTO, tx)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	tx.Commit()

	go func() {
		existingRiskRegister, err := rrs.GetFullDetailsByID(riskRegister.ID)
		if err != nil {
			utility.LogError("Failed to fetch responsible user", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		// log
		err = rrs.activityLogSvc.LogActivity(
			models.CreateAction,
			models.RiskRegister{}.TableName(),
			id,
			existingRiskRegister,
			currentUserID,
			"risk", "risk_register", "", nil,
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

		createdBy, err := rrs.userSvc.GetByID(currentUserID)
		if err != nil {
			utility.LogError("Failed to fetch created by user", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		// Send notification to responsible user and/or responsible unit head
		action := "responsible-user-assign"
		var notificationUserIDs []uint
		var emailTo []string

		// If responsible user is assigned, add them to notification
		if riskRegisterFormDTO.ResponsibleUserID != nil {
			notificationUserIDs = append(notificationUserIDs, *riskRegisterFormDTO.ResponsibleUserID)
			responsibleUser, err := rrs.userSvc.GetByID(*riskRegisterFormDTO.ResponsibleUserID)
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
			headUsers, err := rrs.organizationUnitsSvc.GetParentOrganizationUnitAndCurrentUnitHead(*riskRegisterFormDTO.ResponsibleUnitID)
			if err != nil {
				utility.LogError("Failed to get responsible unit head users for notification", map[string]interface{}{
					"unit_id": *riskRegisterFormDTO.ResponsibleUnitID,
					"error":   err.Error(),
				})
			} else {
				if headUsers.HeadUser != nil {
					// Avoid duplicate notifications if the head user is also the responsible user
					alreadyNotified := false
					for _, notifiedUserID := range notificationUserIDs {
						if notifiedUserID == headUsers.HeadUser.ID {
							alreadyNotified = true
							break
						}
					}
					if !alreadyNotified {
						notificationUserIDs = append(notificationUserIDs, headUsers.HeadUser.ID)
						emailTo = append(emailTo, headUsers.HeadUser.Email)
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
				ReferenceNumber:      referenceNumber,
				Entity:               &entity,
				Action:               &action,
				URLParam:             urlParam,
				IsSystemNotification: false,
				CreatedByID:          &currentUserID,
			}

			// Use responsible user details for email if available, otherwise use unit head
			var firstNameEn, lastNameEn string
			if riskRegisterFormDTO.ResponsibleUserID != nil {
				responsibleUser, _ := rrs.userSvc.GetByID(*riskRegisterFormDTO.ResponsibleUserID)
				if responsibleUser != nil {
					firstNameEn = responsibleUser.FirstNameEn
					lastNameEn = responsibleUser.LastNameEn
				}
			}
			if firstNameEn == "" && len(emailTo) > 0 {
				// Fallback to first recipient's name
				firstUser, _ := rrs.userSvc.GetByID(notificationUserIDs[0])
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
			err = rrs.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "risk")
			if err != nil {
				utility.LogError("Failed to send notifications", map[string]interface{}{
					"error": err.Error(),
				})
			}
		}

		// Fetch risk settings
		riskSetting, err := rrs.riskSettingsSvc.Get()
		if err != nil {
			utility.LogError("Failed to fetch risk setting", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		// Get Review User IDs
		reviewUserIDs, err := rrs.userSvc.GetUserIdsByRole(riskSetting.ReviewRoleID)
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
			user, err := rrs.userSvc.GetByID(userID)
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

		err = rrs.notificationSvc.SendNotification(dbNotificationReview, nil, "risk")
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

		err = rrs.notificationSvc.SendNotification(nil, emailNotificationReview, "risk")
		if err != nil {
			utility.LogError("Failed to send review user email", map[string]interface{}{
				"error":       err.Error(),
				"user_emails": reviewEmailTo,
			})
		}
	}()

	return id, nil
}

func (rrs *RiskRegisterService) Update(id uint, currentUserID uint, riskRegisterFormDTO *riskDTOs.RiskRegisterFormDTO) error {
	tx := rrs.repo.BeginTransaction()

	oldRiskRegister, err := rrs.GetFullDetailsByID(id, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := rrs.validateRiskRegisterFormDTO(riskRegisterFormDTO); err != nil {
		return err
	}

	existingRiskRegister, err := rrs.GetByID(id)
	if err != nil {
		return err
	}

	dates, err := rrs.parseRiskRegisterDates(riskRegisterFormDTO)
	if err != nil {
		return err
	}

	oldResponsibleUserID := existingRiskRegister.ResponsibleUserID

	// if existingRiskRegister.ResponsibleUserID != riskRegisterFormDTO.ResponsibleUserID {
	// 	responsibleUnitID, err := rrs.userSvc.GetUserPrimaryUnitID(uint(riskRegisterFormDTO.ResponsibleUserID))
	// 	if err != nil {
	// 		return utility.NewValidationError(fmt.Errorf("organization unit not set for responsible user"), "not_exist", "risk responsible unit")
	// 	}
	// 	existingRiskRegister.ResponsibleUnitID = responsibleUnitID
	// }

	existingRiskRegister.RiskCategoryID = riskRegisterFormDTO.RiskCategoryID
	existingRiskRegister.RiskSourceID = riskRegisterFormDTO.RiskSourceID
	existingRiskRegister.ResponsibleUserID = riskRegisterFormDTO.ResponsibleUserID
	existingRiskRegister.ResponsibleUnitID = riskRegisterFormDTO.ResponsibleUnitID
	existingRiskRegister.Title = riskRegisterFormDTO.Title
	existingRiskRegister.Description = riskRegisterFormDTO.Description
	existingRiskRegister.ReportedAt = dates.ReportedAt
	existingRiskRegister.ReviewFrequencyID = riskRegisterFormDTO.ReviewFrequencyID
	existingRiskRegister.NextReviewDate = dates.NextReviewDate
	existingRiskRegister.LastReviewNote = riskRegisterFormDTO.LastReviewNote
	existingRiskRegister.SeverityLevelID = riskRegisterFormDTO.SeverityLevelID
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

	if err := rrs.repo.Update(existingRiskRegister, tx); err != nil {
		tx.Rollback()
		return err
	}

	existingRisk, err := rrs.GetFullDetailsByID(id, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	if riskRegisterFormDTO.InherentRiskMatrixLikelihoodID != nil && riskRegisterFormDTO.InherentRiskMatrixImpactID != nil {
		inherentEvaluationFormDTO := &riskDTOs.RiskInherentEvaluationFormDTO{
			RiskMatrixLikelihoodID: *riskRegisterFormDTO.InherentRiskMatrixLikelihoodID,
			RiskMatrixImpactID:     *riskRegisterFormDTO.InherentRiskMatrixImpactID,
		}

		if existingRiskRegister.InherentLikelihoodID != nil && existingRiskRegister.InherentImpactID != nil {
			existingRiskEvaluation, err := rrs.riskInherentEvaluationSvc.GetFullDetailsByID(id, tx)
			if err != nil {
				tx.Rollback()
				return err
			}
			_, err = rrs.riskInherentEvaluationSvc.Delete(id, existingRiskEvaluation.ID, currentUserID, tx)
			if err != nil {
				tx.Rollback()
				return err
			}
		}

		_, err = rrs.riskInherentEvaluationSvc.Create(id, currentUserID, inherentEvaluationFormDTO, tx)
		if err != nil {
			tx.Rollback()
			return err
		}
	} else {

		if existingRiskRegister.InherentLikelihoodID != nil && existingRiskRegister.InherentImpactID != nil {
			existingRiskEvaluation, err := rrs.riskInherentEvaluationSvc.GetFullDetailsByID(id, tx)
			if err != nil {
				tx.Rollback()
				return err
			}
			_, err = rrs.riskInherentEvaluationSvc.Delete(id, existingRiskEvaluation.ID, currentUserID, tx)
			if err != nil {
				tx.Rollback()
				return err
			}
		}

	}

	existingSLARegisterIDs := make(map[uint]bool)
	for _, sla := range existingRisk.SLARegisters {
		existingSLARegisterIDs[sla.SLARegisterID] = true
	}

	// Collect new slaess IDs
	newSLARegisterIDs := make(map[uint]bool)
	for _, sla := range riskRegisterFormDTO.SLARegisterIDs {
		newSLARegisterIDs[*sla] = true
	}

	// Add new processes
	var slaRegistersToAdd []*models.RiskSLA
	for _, newSLARegisterID := range riskRegisterFormDTO.SLARegisterIDs {
		if _, exists := existingSLARegisterIDs[*newSLARegisterID]; !exists {
			// Create a new risk slaRegister record
			slaRegister := &models.RiskSLA{
				RiskRegisterID: id,
				SLARegisterID:  *newSLARegisterID,
				CreatedByID:    currentUserID,
			}
			slaRegistersToAdd = append(slaRegistersToAdd, slaRegister)
		}
	}

	// Save the new slaRegisters
	if len(slaRegistersToAdd) > 0 {
		if err := rrs.repo.CreateRiskSLAs(slaRegistersToAdd, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	// Remove slaRegisters not in the new list
	var deleteSLARegisterIDs []*uint

	for _, existingSLARegister := range existingRisk.SLARegisters {
		if _, exists := newSLARegisterIDs[existingSLARegister.SLARegisterID]; !exists {
			// Append slaRegister to delete list
			deleteSLARegisterIDs = append(deleteSLARegisterIDs, existingSLARegister.ID)
		}
	}

	// Delete slaRegisters not in the new list
	if len(deleteSLARegisterIDs) > 0 {
		if err := rrs.repo.DeleteRiskSLAs(deleteSLARegisterIDs, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	// Save Documents
	// documents update
	var documents []*models.RiskDocument
	var deleteDocumentIds []*uint
	var documentsToDelete []*models.RiskDocument

	// Iterate over each document in the array
	for _, doc := range riskRegisterFormDTO.Documents {
		if doc.IsNewFile != nil && *doc.IsNewFile {

			dataToEncrypt := []byte(doc.FileName)
			encryptedFilename, err := utility.EncryptData(dataToEncrypt)
			if err != nil {
				tx.Rollback()
				return err
			}

			url := utility.GenerateDocumentURL("risk", "risk-register", strconv.FormatUint(uint64(id), 10), "documents")
			// Create a new RiskDocument instance
			document := &models.RiskDocument{
				RiskRegisterID: id,
				URL:            url,
				Name:           doc.FileName,
				Size:           doc.Size,
				Token:          encryptedFilename,
				OriginalName:   doc.OriginalName,
				Type:           doc.Type,
				CreatedByID:    currentUserID,
			}
			documents = append(documents, document)

		} else if doc.IsDeletedFile != nil && *doc.IsDeletedFile {
			// Handle deleted file case
			if doc.ID != nil {
				existingDocument, err := rrs.repo.GetDocumentByID(*doc.ID)
				if err != nil {
					if errors.Is(err, gorm.ErrRecordNotFound) {
						return utility.NewResourceNotFoundError(err, id)
					}
					return err
				}

				documentsToDelete = append(documentsToDelete, existingDocument)
				deleteDocumentIds = append(deleteDocumentIds, &existingDocument.ID)
			}
		}
	}

	if len(documents) > 0 {
		if err := rrs.repo.CreateDocuments(documents, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	if len(deleteDocumentIds) > 0 {
		if err := rrs.repo.DeleteDocuments(deleteDocumentIds, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	if len(documents) > 0 {
		for _, doc := range documents {
			storagePath := filepath.Join("risk", "risk-register", strconv.FormatUint(uint64(existingRiskRegister.ID), 10), "documents")
			if err := utility.StoreFile(doc.Name, storagePath); err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	if len(documentsToDelete) > 0 {
		for _, document := range documentsToDelete {
			storagePath := filepath.Join("risk", "risk-register", strconv.FormatUint(uint64(existingRiskRegister.ID), 10), "documents")
			if err := utility.RemoveFile(document.Name, storagePath); err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	tx.Commit()

	go func() {
		SLARegisterAdded := []uint{}
		SLARegisterRemoved := []uint{}
		for id := range newSLARegisterIDs {
			if !existingSLARegisterIDs[id] {
				SLARegisterAdded = append(SLARegisterAdded, id)
			}
		}
		for id := range existingSLARegisterIDs {
			if !newSLARegisterIDs[id] {
				SLARegisterRemoved = append(SLARegisterRemoved, id)
			}
		}

		// Prepare log data
		changedData := make(map[string]interface{})
		associatedItems := make(map[string]interface{})

		if len(SLARegisterAdded) > 0 || len(SLARegisterRemoved) > 0 {
			associatedItems["sla_register_ids"] = map[string]interface{}{
				"added":   SLARegisterAdded,
				"removed": SLARegisterRemoved,
			}
		}

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
			err := rrs.activityLogSvc.LogActivity(
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
			updatedByUser, err := rrs.userSvc.GetByID(currentUserID)
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
					ReferenceNumber:      *existingRiskRegister.ReferenceNumber,
					Entity:               &entity,
					Action:               &oldResponsibleUserAction,
					URLParam:             urlParam,
					IsSystemNotification: false,
					CreatedByID:          &currentUserID,
				}

				oldResponsibleUser, err := rrs.userSvc.GetByID(*oldResponsibleUserID)
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

					err = rrs.notificationSvc.SendNotification(oldResponsibleUserNotificationDTO, oldResponsibleUserEmailNotificationDTO, "risk")
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
				newResponsibleUser, err := rrs.userSvc.GetByID(*riskRegisterFormDTO.ResponsibleUserID)
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
				headUsers, err := rrs.organizationUnitsSvc.GetParentOrganizationUnitAndCurrentUnitHead(*riskRegisterFormDTO.ResponsibleUnitID)
				if err != nil {
					utility.LogError("Failed to get responsible unit head users for notification", map[string]interface{}{
						"unit_id": *riskRegisterFormDTO.ResponsibleUnitID,
						"error":   err.Error(),
					})
				} else {
					if headUsers.HeadUser != nil {
						// Avoid duplicate notifications if the head user is also the responsible user
						alreadyNotified := false
						for _, notifiedUserID := range notificationUserIDs {
							if notifiedUserID == headUsers.HeadUser.ID {
								alreadyNotified = true
								break
							}
						}
						if !alreadyNotified {
							notificationUserIDs = append(notificationUserIDs, headUsers.HeadUser.ID)
							emailTo = append(emailTo, headUsers.HeadUser.Email)
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
					ReferenceNumber:      *existingRiskRegister.ReferenceNumber,
					Entity:               &entity,
					Action:               &newResponsibleUserAction,
					URLParam:             urlParam,
					IsSystemNotification: false,
					CreatedByID:          &currentUserID,
				}

				// Use responsible user details for email if available, otherwise use unit head
				var firstNameEn, lastNameEn string
				if riskRegisterFormDTO.ResponsibleUserID != nil {
					newResponsibleUser, _ := rrs.userSvc.GetByID(*riskRegisterFormDTO.ResponsibleUserID)
					if newResponsibleUser != nil {
						firstNameEn = newResponsibleUser.FirstNameEn
						lastNameEn = newResponsibleUser.LastNameEn
					}
				}
				if firstNameEn == "" && len(emailTo) > 0 {
					// Fallback to first recipient's name
					firstUser, _ := rrs.userSvc.GetByID(notificationUserIDs[0])
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

				err = rrs.notificationSvc.SendNotification(newResponsibleUserNotificationDTO, newResponsibleUserEmailNotificationDTO, "risk")
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

func (rrs *RiskRegisterService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	tx := rrs.repo.BeginTransaction()

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

	documentsToDelete := make(map[uint][]*riskDTOs.RiskDocumentDTO)

	for _, id := range ids {
		risk, err := rrs.GetFullDetailsByID(id, tx)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		riskRegister, err := rrs.repo.GetFullDetailsByID(id, tx)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		// Prepare notification data before deletion
		notifications = append(notifications, struct {
			OwnerID           *uint
			ResponsibleUnitID *uint
			TitleEn           string
			TitleAr           string
			Description       *string
			ReferenceNumber   string
		}{
			OwnerID:           &risk.ResponsibleUserID,
			ResponsibleUnitID: riskRegister.ResponsibleUnitID,
			TitleEn:           risk.Title,
			Description:       risk.Description,
			ReferenceNumber:   risk.ReferenceNumber,
		})

		deletedRiskRegister = append(deletedRiskRegister, map[string]interface{}{
			"id":    id,
			"title": risk.Title,
		})

		if len(risk.SLARegisters) > 0 {
			var deleteSLARegisterIDs []*uint
			for _, sla := range risk.SLARegisters {
				deleteSLARegisterIDs = append(deleteSLARegisterIDs, sla.ID)
			}
			if err := rrs.repo.DeleteRiskSLAs(deleteSLARegisterIDs, tx); err != nil {
				failedIds = append(failedIds, id)
				continue
			}
		}

		// Delete timeline
		// todo
		if err := rrs.riskTimelineSvc.Delete(id, tx); err != nil {
			failedIds = append(failedIds, id)
			continue
		}
		// Delete each risk document if documents exist
		if len(risk.Documents) > 0 {
			var deleteDocumentIds []*uint
			for _, doc := range risk.Documents {
				if doc != nil {
					documentsToDelete[id] = append(documentsToDelete[id], doc)
					deleteDocumentIds = append(deleteDocumentIds, doc.ID)
				}
			}
			if err := rrs.repo.DeleteDocuments(deleteDocumentIds, tx); err != nil {
				failedIds = append(failedIds, id)
				continue
			}

		}

		// delete workflow
		var deleteRiskWorkflowIds []*uint
		var deleteRiskWorkflowParticipantIds []*uint
		var deleteRiskWorkflowReviewIds []*uint
		if len(riskRegister.RiskWorkflows) > 0 {
			for _, riskWorkflow := range riskRegister.RiskWorkflows {
				deleteRiskWorkflowIds = append(deleteRiskWorkflowIds, &riskWorkflow.ID)

				if len(riskWorkflow.RiskWorkflowParticipants) > 0 {
					for _, participant := range riskWorkflow.RiskWorkflowParticipants {
						deleteRiskWorkflowParticipantIds = append(deleteRiskWorkflowParticipantIds, &participant.ID)
					}
				}

				if len(riskWorkflow.RiskWorkflowReviews) > 0 {
					for _, reviews := range riskWorkflow.RiskWorkflowReviews {
						deleteRiskWorkflowReviewIds = append(deleteRiskWorkflowReviewIds, &reviews.ID)
					}
				}
			}
		}
		if err := rrs.riskWorkflowSvc.DeleteRiskWorkflowParticipants(deleteRiskWorkflowParticipantIds, tx); err != nil {
			failedIds = append(failedIds, id)
			continue
		}

		if err := rrs.riskWorkflowSvc.DeleteRiskReviews(deleteRiskWorkflowReviewIds, tx); err != nil {
			failedIds = append(failedIds, id)
			continue
		}

		if err := rrs.riskWorkflowSvc.DeleteRiskWorkflows(deleteRiskWorkflowIds, tx); err != nil {
			failedIds = append(failedIds, id)
			continue
		}

		if err := rrs.repo.DeleteRiskUpdatePlanRisksByRiskID(id, tx); err != nil {
			failedIds = append(failedIds, id)
			continue
		}

		if err := rrs.repo.DeleteRiskEvaluationsByRiskID(id, tx); err != nil {
			failedIds = append(failedIds, id)
			continue
		}

		// Delete the risk
		if err := rrs.repo.Delete(id, tx); err != nil {
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

	if len(successfulIds) > 0 {
		for riskID, documents := range documentsToDelete {
			for _, doc := range documents {
				storagePath := filepath.Join("risk", "risk-register", strconv.FormatUint(uint64(riskID), 10), "documents")
				if err := utility.RemoveFile(doc.Name, storagePath); err != nil {
					tx.Rollback()
				}
			}
			destinationPath := filepath.Join("risk", "risk-register", strconv.FormatUint(uint64(riskID), 10))
			if err := utility.RemoveFolder(destinationPath); err != nil {
				tx.Rollback()
			}
		}
	}

	tx.Commit()

	// Log bulk delete activities for successful deletions
	if len(successfulIds) > 0 {
		go func() {
			err := rrs.activityLogSvc.LogBulkActivities(
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
					ownerUser, err := rrs.userSvc.GetByID(*notification.OwnerID)
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
					headUsers, err := rrs.organizationUnitsSvc.GetParentOrganizationUnitAndCurrentUnitHead(*notification.ResponsibleUnitID)
					if err != nil {
						utility.LogError("Failed to get responsible unit head users for notification", map[string]interface{}{
							"unit_id": *notification.ResponsibleUnitID,
							"error":   err.Error(),
						})
					} else {
						if headUsers.HeadUser != nil {
							// Avoid duplicate notifications if the head user is also the responsible user
							alreadyNotified := false
							for _, notifiedUserID := range notificationUserIDs {
								if notifiedUserID == headUsers.HeadUser.ID {
									alreadyNotified = true
									break
								}
							}
							if !alreadyNotified {
								notificationUserIDs = append(notificationUserIDs, headUsers.HeadUser.ID)
								emailTo = append(emailTo, headUsers.HeadUser.Email)
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
						ownerUser, _ := rrs.userSvc.GetByID(*notification.OwnerID)
						if ownerUser != nil {
							firstNameEn = ownerUser.FirstNameEn
							lastNameEn = ownerUser.LastNameEn
						}
					}
					if firstNameEn == "" && len(emailTo) > 0 {
						// Fallback to first recipient's name
						firstUser, _ := rrs.userSvc.GetByID(notificationUserIDs[0])
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
					err = rrs.notificationSvc.SendNotification(
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

func (rrs *RiskRegisterService) UpdateResponseStrategy(riskID uint, currentUserID uint, riskResponseStrategyFormDTO *riskDTOs.RiskResponseStrategyFormDTO, tx ...utility.DatabaseTransaction) error {
	var transaction utility.DatabaseTransaction

	if len(tx) > 0 && tx[0] != nil {
		transaction = tx[0]
	} else {
		transaction = rrs.repo.BeginTransaction()
		defer func() {
			if r := recover(); r != nil {
				transaction.Rollback()
			}
		}()
	}

	existingRisk, err := rrs.GetByID(riskID, transaction)
	if err != nil {
		transaction.Rollback()
		return err
	}

	existingRisk.RiskResponseStrategyID = &riskResponseStrategyFormDTO.RiskResponseStrategyID

	if err := rrs.repo.Update(existingRisk, transaction); err != nil {
		transaction.Rollback()
		return err
	}

	// Create RiskResponse if not exists (only with risk_id, other fields null)
	if err := rrs.riskResponseSvc.CreateIfNotExists(riskID, currentUserID, transaction); err != nil {
		transaction.Rollback()
		return err
	}

	// update timeline
	// todo
	timelineType, err := rrs.riskTimelineSvc.GetByType("risk_response")
	if err != nil {
		transaction.Rollback()
		return err
	}
	timelineStatus, err := rrs.riskTimelineSvc.GetByStatusType("completed")
	if err != nil {
		transaction.Rollback()
		return err
	}
	currentTime := time.Now()
	timelineFormDTO := &riskDTOs.RiskTimelineFormDTO{
		PerformedByID:    &currentUserID,
		TimelineStatusID: timelineStatus.ID,
		PerformedAt:      &currentTime,
	}
	if err := rrs.riskTimelineSvc.Update(riskID, timelineFormDTO, timelineType.ID, transaction); err != nil {
		transaction.Rollback()
		return err
	}

	transaction.Commit()

	return nil
}

func (rrs *RiskRegisterService) DeleteResponseStrategy(riskID uint, currentUserID uint) error {
	tx := rrs.repo.BeginTransaction()

	existingRisk, err := rrs.GetByID(riskID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Get risk response by risk register ID
	riskResponse, err := rrs.riskResponseSvc.GetByRiskRegisterID(riskID, tx)
	if err != nil {
		// If no risk response exists, just clear the strategy
		existingRisk.RiskResponseStrategyID = nil
		if err := rrs.repo.Update(existingRisk, tx); err != nil {
			tx.Rollback()
			return err
		}
	} else {
		// Delete workflow-related data
		if err := rrs.repo.DeleteByRiskResponseID(riskResponse.ID, tx); err != nil {
			tx.Rollback()
			return err
		}

		// Delete review comments
		if err := rrs.repo.DeleteByRiskResponseID(riskResponse.ID, tx); err != nil {
			tx.Rollback()
			return err
		}

		// Delete risk response
		if err := rrs.riskResponseSvc.DeleteByRiskRegisterID(riskID, tx); err != nil {
			tx.Rollback()
			return err
		}

		existingRisk.RiskResponseStrategyID = nil

		if err := rrs.repo.Update(existingRisk, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	// update timeline
	// todo
	timelineType, err := rrs.riskTimelineSvc.GetByType("risk_response")
	if err != nil {
		tx.Rollback()
		return err
	}
	timelineStatus, err := rrs.riskTimelineSvc.GetByStatusType("pending")
	if err != nil {
		tx.Rollback()
		return err
	}
	timelineFormDTO := &riskDTOs.RiskTimelineFormDTO{
		PerformedByID:    nil,
		TimelineStatusID: timelineStatus.ID,
		PerformedAt:      nil,
	}
	if err := rrs.riskTimelineSvc.Update(riskID, timelineFormDTO, timelineType.ID, tx); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil
}

func (rrs *RiskRegisterService) GetControlPerformance(riskID uint) (map[string]float64, error) {
	_, err := rrs.GetByID(riskID)
	if err != nil {
		return nil, err
	}

	inherentScore, err := rrs.getControlScoreByType(riskID, "inherent")
	if err != nil {
		return nil, err
	}

	residualScore, err := rrs.getControlScoreByType(riskID, "residual")
	if err != nil {
		return nil, err
	}

	scores := map[string]float64{
		"inherent_control_score": inherentScore,
		"residual_control_score": residualScore,
	}

	return scores, nil
}

func (rrs *RiskRegisterService) getControlScoreByType(riskID uint, controlType string) (float64, error) {
	riskControlAnalyses, err := rrs.repo.GetControlSummary(riskID, controlType)
	if err != nil {
		return 0, err
	}

	if len(riskControlAnalyses) == 0 {
		return 0, nil
	}

	implementationStatusCount := make(map[string]int)
	totalControls := len(riskControlAnalyses)

	for _, rca := range riskControlAnalyses {
		implStatus := rca.ControlImplementationStatus.Type
		implementationStatusCount[implStatus]++
	}

	controlScore := (float64(implementationStatusCount["implemented"]) / float64(totalControls)) * 100

	return controlScore, nil
}

func (rrs *RiskRegisterService) UpdateRiskRating(riskID uint, currentUserID uint, riskEvaluationDetails *models.RiskEvaluation, tx utility.DatabaseTransaction, evaluationType string) error {
	existingRisk, err := rrs.GetByID(riskID, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	if evaluationType == "inherent" {
		existingRisk.InherentImpactID = &riskEvaluationDetails.RiskMatrixImpactID
		existingRisk.InherentLikelihoodID = &riskEvaluationDetails.RiskMatrixLikelihoodID
		existingRisk.InherentRiskRatingID = &riskEvaluationDetails.RiskMatrixRiskRatingID
		existingRisk.InherentRiskRatingValueID = &riskEvaluationDetails.RiskMatrixRiskRatingValueID
	} else if evaluationType == "residual" {
		existingRisk.ResidualImpactID = &riskEvaluationDetails.RiskMatrixImpactID
		existingRisk.ResidualLikelihoodID = &riskEvaluationDetails.RiskMatrixLikelihoodID
		existingRisk.ResidualRiskRatingID = &riskEvaluationDetails.RiskMatrixRiskRatingID
		existingRisk.ResidualRiskRatingValueID = &riskEvaluationDetails.RiskMatrixRiskRatingValueID
	}

	existingRisk.CurrentImpactID = &riskEvaluationDetails.RiskMatrixImpactID
	existingRisk.CurrentLikelihoodID = &riskEvaluationDetails.RiskMatrixLikelihoodID
	existingRisk.CurrentRiskRatingID = &riskEvaluationDetails.RiskMatrixRiskRatingID
	existingRisk.CurrentRiskRatingValueID = &riskEvaluationDetails.RiskMatrixRiskRatingValueID

	if err := rrs.repo.Update(existingRisk, tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := rrs.UpdateProcessRiskRating(existingRisk.ID, tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := rrs.UpdateControlRiskRating(existingRisk.ID, tx); err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (rrs *RiskRegisterService) UpdateProcessRiskRating(riskID uint, tx utility.DatabaseTransaction) error {
	// Get the risk register details for the given riskID
	risk, err := rrs.repo.GetFullDetailsByID(riskID, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Iterate through each associated process and update their risk rating
	for _, riskProcess := range risk.Processes {
		// Get the full process details (including associated risks)
		process, err := rrs.processSvc.GetFullDetailsByID(riskProcess.ProcessID, tx)
		if err != nil {
			tx.Rollback()
			return err
		}

		// Find the maximum risk rating value from all associated risks
		var maxRiskRatingValue *models.RiskMatrixRiskRatingValue
		for _, processRisk := range process.Risks {

			riskDetail, err := rrs.repo.GetFullDetailsByID(processRisk.RiskRegisterID, tx)
			if err != nil {
				tx.Rollback()
				return err
			}

			if maxRiskRatingValue == nil || (riskDetail.CurrentRiskRatingValue != nil && riskDetail.CurrentRiskRatingValue.Value > maxRiskRatingValue.Value) {
				maxRiskRatingValue = riskDetail.CurrentRiskRatingValue
			}
		}

		existingProcess, err := rrs.processSvc.GetByID(riskProcess.ProcessID, tx)
		if err != nil {
			tx.Rollback()
			return err
		}

		// If we found a max risk rating value, update the process with the corresponding rating
		if maxRiskRatingValue != nil {
			// Update the process risk rating values
			existingProcess.RiskMatrixRiskRatingID = &maxRiskRatingValue.RiskMatrixRiskRatingID
			existingProcess.RiskMatrixRiskRatingValueID = &maxRiskRatingValue.ID

			// Update the process risk rating in the database
			err = rrs.processSvc.UpdateProcessRiskRating(existingProcess, tx)
			if err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	return nil
}

func (rrs *RiskRegisterService) UpdateControlRiskRating(riskID uint, tx utility.DatabaseTransaction) error {
	// Get the risk register details for the given riskID
	risk, err := rrs.repo.GetFullDetailsByID(riskID, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Iterate through each associated control and update their risk rating
	for _, riskControl := range risk.Controls {
		// Get the full control details (including associated risks)
		control, err := rrs.controlSvc.GetFullDetailsByID(riskControl.ControlID, tx)
		if err != nil {
			tx.Rollback()
			return err
		}

		// Find the maximum risk rating value from all associated risks
		var maxRiskRatingValue *models.RiskMatrixRiskRatingValue
		for _, controlRisk := range control.Risks {

			riskDetail, err := rrs.repo.GetFullDetailsByID(controlRisk.RiskRegisterID, tx)
			if err != nil {
				tx.Rollback()
				return err
			}

			if maxRiskRatingValue == nil || (riskDetail.CurrentRiskRatingValue != nil && riskDetail.CurrentRiskRatingValue.Value > maxRiskRatingValue.Value) {
				maxRiskRatingValue = riskDetail.CurrentRiskRatingValue
			}
		}

		existingControl, err := rrs.controlSvc.GetByID(riskControl.ControlID, tx)
		if err != nil {
			tx.Rollback()
			return err
		}

		// If we found a max risk rating value, update the control with the corresponding rating
		if maxRiskRatingValue != nil {
			// Update the control risk rating values
			existingControl.RiskMatrixRiskRatingID = &maxRiskRatingValue.RiskMatrixRiskRatingID
			existingControl.RiskMatrixRiskRatingValueID = &maxRiskRatingValue.ID

			// Update the control risk rating in the database
			err = rrs.controlSvc.UpdateControlRiskRating(existingControl, tx)
			if err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	return nil
}

func (rrs *RiskRegisterService) ClearRiskRating(riskID uint, updateDTO *riskDTOs.RiskRegisterRiskRatingUpdateDTO, tx utility.DatabaseTransaction) error {
	riskRegister, err := rrs.GetByID(riskID, tx)
	if err != nil {
		return err
	}

	riskRegister.InherentImpactID = updateDTO.InherentImpactID
	riskRegister.InherentLikelihoodID = updateDTO.InherentLikelihoodID
	riskRegister.InherentRiskRatingID = updateDTO.InherentRiskRatingID
	riskRegister.InherentRiskRatingValueID = updateDTO.InherentRiskRatingValueID

	riskRegister.CurrentImpactID = updateDTO.CurrentImpactID
	riskRegister.CurrentLikelihoodID = updateDTO.CurrentLikelihoodID
	riskRegister.CurrentRiskRatingID = updateDTO.CurrentRiskRatingID
	riskRegister.CurrentRiskRatingValueID = updateDTO.CurrentRiskRatingValueID

	riskRegister.ResidualImpactID = updateDTO.ResidualImpactID
	riskRegister.ResidualLikelihoodID = updateDTO.ResidualLikelihoodID
	riskRegister.ResidualRiskRatingID = updateDTO.ResidualRiskRatingID
	riskRegister.ResidualRiskRatingValueID = updateDTO.ResidualRiskRatingValueID

	if err := rrs.repo.Update(riskRegister, tx); err != nil {
		return err
	}

	if err := rrs.UpdateProcessRiskRating(riskRegister.ID, tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := rrs.UpdateControlRiskRating(riskRegister.ID, tx); err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (rrs *RiskRegisterService) validateRiskRegisterFormDTO(riskRegisterFormDTO *riskDTOs.RiskRegisterFormDTO) error {

	if riskRegisterFormDTO.RiskCategoryID == nil {
		if _, err := rrs.riskCategorySvc.GetByID(*riskRegisterFormDTO.RiskCategoryID); err != nil {
			return utility.NewValidationError(fmt.Errorf("risk category does not exist"), "not_exist", "risk category")
		}
	}

	if riskRegisterFormDTO.RiskSourceID != nil {
		if _, err := rrs.riskSourceSvc.GetByID(*riskRegisterFormDTO.RiskSourceID); err != nil {
			return utility.NewValidationError(fmt.Errorf("risk source does not exist"), "not_exist", "risk source")
		}
	}

	if riskRegisterFormDTO.ResponsibleUserID != nil {
		if _, err := rrs.userSvc.GetByID(*riskRegisterFormDTO.ResponsibleUserID); err != nil {
			return utility.NewValidationError(fmt.Errorf("responsible user does not exist"), "not_exist", "responsible user")
		}
	}

	if riskRegisterFormDTO.ReviewFrequencyID != nil {
		if _, err := rrs.reviewFrequencySvc.GetByID(*riskRegisterFormDTO.ReviewFrequencyID); err != nil {
			return utility.NewValidationError(fmt.Errorf("review frequency does not exist"), "not_exist", "review frequency")
		}
	}

	slaRegisterIds := riskRegisterFormDTO.SLARegisterIDs

	slaRegisterIDs := make(map[uint]bool)

	for _, slaRegisterId := range slaRegisterIds {
		if slaRegisterIDs[*slaRegisterId] {
			return utility.NewValidationError(fmt.Errorf("duplicate slaRegister ID: %d", slaRegisterId), "duplicate", "slaRegister")
		}

		if _, err := rrs.slaRegisterSvc.GetByID(*slaRegisterId); err != nil {
			return utility.NewValidationError(fmt.Errorf("slaRegister does not exist for ID: %d", slaRegisterId), "not_exist", "slaRegister")
		}

		slaRegisterIDs[*slaRegisterId] = true
	}

	assets := riskRegisterFormDTO.AssetIds
	assetIDs := make(map[uint]bool)

	for _, asset := range assets {
		if assetIDs[*asset] {
			return utility.NewValidationError(fmt.Errorf("duplicate asset ID: %d", asset), "duplicate", "asset")
		}

		if _, err := rrs.assetRegisterSvc.GetByID(*asset); err != nil {
			return utility.NewValidationError(fmt.Errorf("asset does not exist for ID: %d", asset), "not_exist", "asset")
		}

		// Add the asset ID to the map
		assetIDs[*asset] = true
	}

	// Validate that division contains the responsible unit

	if !rrs.organizationUnitSvc.IsUnitUnderParent(*riskRegisterFormDTO.ResponsibleUnitID, riskRegisterFormDTO.DivisionID) {
		return utility.NewValidationError(fmt.Errorf("responsible unit must be under the selected division"), "invalid_hierarchy", "responsible_unit")
	}

	// Validate that section is under responsible unit
	if riskRegisterFormDTO.SectionID != nil && riskRegisterFormDTO.ResponsibleUnitID != nil {
		if !rrs.organizationUnitSvc.IsUnitUnderParent(*riskRegisterFormDTO.SectionID, *riskRegisterFormDTO.ResponsibleUnitID) {
			return utility.NewValidationError(fmt.Errorf("section must be under the responsible unit"), "invalid_hierarchy", "section")
		}
	}

	// Validate that sub section is under section
	if riskRegisterFormDTO.SubSectionID != nil && riskRegisterFormDTO.SectionID != nil {
		if !rrs.organizationUnitSvc.IsUnitUnderParent(*riskRegisterFormDTO.SubSectionID, *riskRegisterFormDTO.SectionID) {
			return utility.NewValidationError(fmt.Errorf("sub section must be under the section"), "invalid_hierarchy", "sub_section")
		}
	}

	return nil
}

func (rrs *RiskRegisterService) parseRiskRegisterDates(dto *riskDTOs.RiskRegisterFormDTO) (*riskDTOs.RiskRegisterParsedDates, error) {
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

func (rrs *RiskRegisterService) GetDocumentPath(id uint, token string) (string, error) {
	// Get full details of the risk register by ID
	existingRiskRegister, err := rrs.GetFullDetailsByID(id)
	if err != nil {
		return "", err
	}

	// Check if there are any associated documents
	if len(existingRiskRegister.Documents) == 0 {
		return "", fmt.Errorf("no documents found for the risk register")
	}

	// Decrypt the token to get the file name
	decryptedFilenameBytes, err := utility.DecryptData(token)
	if err != nil {
		return "", fmt.Errorf("resource not found")
	}

	decryptedFilename := string(decryptedFilenameBytes)
	if decryptedFilename == "" {
		return "", fmt.Errorf("resource not found")
	}

	// Define the storage path
	storagePath := filepath.Join("risk", "risk-register", strconv.FormatUint(uint64(id), 10), "documents")

	// Get the base storage path from environment variables
	baseStoragePath := config.GetEnv("STORAGE_PATH", "")
	fullPath := filepath.Join(baseStoragePath, storagePath, decryptedFilename)

	// Check if the file exists at the specified path
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return "", fmt.Errorf("file not found")
	} else if err != nil {
		return "", fmt.Errorf("error accessing file: %w", err)
	}

	return fullPath, nil
}

func (cs *RiskRegisterService) ExportRiskRegister(data []*riskDTOs.RiskRegisterListDTO, currentUserID uint) (*os.File, error) {
	// Fetch user setting for language preference
	userSetting, err := cs.userSettingSvc.GetUserSettingsByUserID(currentUserID)
	if err != nil {
		return nil, err
	}

	// Define headers based on user language setting
	var headers []string
	if userSetting.Language == "ar" {
		headers = []string{
			"SL No.",
			"Ref Number",
			"Risk Code",
			"Registration Date",
			"Division",
			"Section",
			"Sub Section",
			"Is Strategic",
			"Is Operational",
			"Expected Cause",
			"Consequence",
			"Is Internal",
			"Is External",
			"Risk Title",
			"Description",
			"Risk Category",
			"Responsible Unit (TBD)",
			"Responsible User",
			"Risk Source",
			"Risk Status",
			"Likelihood",
			"Impact",
			"Risk Rating",
			"Risk Response",
			"Residual Likelihood",
			"Residual Impact",
			"Residual Risk Rating",
			"Reported At",
			"Review Frequency",
			"Next Review Date",
			"Severity Level",
		}
	} else {
		headers = []string{
			"SL No.",
			"Ref Number",
			"Risk Code",
			"Registration Date",
			"Division",
			"Section",
			"Sub Section",
			"Is Strategic",
			"Is Operational",
			"Expected Cause",
			"Consequence",
			"Is Internal",
			"Is External",
			"Risk Title",
			"Description",
			"Risk Category",
			"Responsible Unit (TBD)",
			"Responsible User",
			"Risk Source",
			"Risk Status",
			"Likelihood",
			"Impact",
			"Risk Rating",
			"Risk Response",
			"Residual Likelihood",
			"Residual Impact",
			"Residual Risk Rating",
			"Reported At",
			"Review Frequency",
			"Next Review Date",
			"Severity Level",
		}
	}
	// Prepare the data rows
	rows := make([][]interface{}, len(data))
	for i, risk := range data {
		// Initialize data variables with default values
		var riskTitle, riskCategory, responsibleUnit, responsibleUser, riskSource, riskStatus string
		var riskResponse, residualLikelihood, residualImpact, residualRiskRating, severityLevel string
		var description, likelihood, impact, riskRating string
		var reviewFrequency, nextReviewDate, reportedAt string

		// Assign values based on language preference and check for nil pointers
		if userSetting.Language == "ar" {
			if risk.RiskCategoryTitleAr != nil {
				riskCategory = *risk.RiskCategoryTitleAr
			}
			if risk.ResponsibleUnitTitleAr != nil {
				responsibleUnit = *risk.ResponsibleUnitTitleAr
			}
			if risk.ResponsibleUserFirstNameAr != nil && risk.ResponsibleUserLastNameAr != nil {
				responsibleUser = *risk.ResponsibleUserFirstNameAr + " " + *risk.ResponsibleUserLastNameAr
			}
			riskSource = utility.GetStringValue(risk.RiskSourceTitleAr)
			riskStatus = risk.RiskStatusTitleAr
			riskResponse = risk.RiskResponseStrategyTitleAr
			likelihood = utility.GetStringValue(risk.InherentLikelihoodTitleAr)
			impact = utility.GetStringValue(risk.InherentImpactTitleAr)
			riskRating = utility.GetStringValue(risk.InherentRiskRatingTitleAr)
			residualLikelihood = utility.GetStringValue(risk.ResidualLikelihoodTitleAr)
			residualImpact = utility.GetStringValue(risk.ResidualImpactTitleAr)
			residualRiskRating = utility.GetStringValue(risk.ResidualRiskRatingTitleAr)
			reviewFrequency = utility.GetStringValue(risk.ReviewFrequencyTitleAr)
			severityLevel = utility.GetStringValue(risk.SeverityLevelTitleAr)
		} else {
			if risk.RiskCategoryTitleEn != nil {
				riskCategory = *risk.RiskCategoryTitleEn
			}
			if risk.ResponsibleUnitTitleEn != nil {
				responsibleUnit = *risk.ResponsibleUnitTitleEn
			}
			if risk.ResponsibleUserFirstNameEn != nil && risk.ResponsibleUserLastNameEn != nil {
				responsibleUser = *risk.ResponsibleUserFirstNameEn + " " + *risk.ResponsibleUserLastNameEn
			}
			riskSource = utility.GetStringValue(risk.RiskSourceTitleEn)
			riskStatus = risk.RiskStatusTitleEn
			riskResponse = risk.RiskResponseStrategyTitleEn
			likelihood = utility.GetStringValue(risk.InherentLikelihoodTitleEn)
			impact = utility.GetStringValue(risk.InherentImpactTitleEn)
			riskRating = utility.GetStringValue(risk.InherentRiskRatingTitleEn)
			residualLikelihood = utility.GetStringValue(risk.ResidualLikelihoodTitleEn)
			residualImpact = utility.GetStringValue(risk.ResidualImpactTitleEn)
			residualRiskRating = utility.GetStringValue(risk.ResidualRiskRatingTitleEn)
			reviewFrequency = utility.GetStringValue(risk.ReviewFrequencyTitleEn)
			severityLevel = utility.GetStringValue(risk.SeverityLevelTitleEn)
		}

		riskTitle = risk.Title
		description = utility.GetStringValue(risk.Description)
		nextReviewDate = utility.GetStringValue(risk.NextReviewDate)
		reportedAt = utility.GetStringValue(risk.ReportedAt)

		// New fields
		var riskCode, registrationDate, division, section, subSection string
		var isStrategic, isOperational, expectedCause, consequence string
		var isInternal, isExternal string

		if risk.RiskCode != nil {
			riskCode = *risk.RiskCode
		}
		if risk.RegistrationDate != nil {
			registrationDate = *risk.RegistrationDate
		}

		if userSetting.Language == "ar" {
			if risk.DivisionTitleAr != nil {
				division = *risk.DivisionTitleAr
			}
			if risk.SectionTitleAr != nil {
				section = *risk.SectionTitleAr
			}
			if risk.SubSectionTitleAr != nil {
				subSection = *risk.SubSectionTitleAr
			}
		} else {
			if risk.DivisionTitleEn != nil {
				division = *risk.DivisionTitleEn
			}
			if risk.SectionTitleEn != nil {
				section = *risk.SectionTitleEn
			}
			if risk.SubSectionTitleEn != nil {
				subSection = *risk.SubSectionTitleEn
			}
		}

		if risk.IsStrategic {
			isStrategic = "Yes"
		} else {
			isStrategic = "No"
		}

		if risk.IsOperational {
			isOperational = "Yes"
		} else {
			isOperational = "No"
		}

		if risk.ExpectedCause != nil {
			expectedCause = *risk.ExpectedCause
		}
		if risk.Consequence != nil {
			consequence = *risk.Consequence
		}

		if risk.IsInternal {
			isInternal = "Yes"
		} else {
			isInternal = "No"
		}

		if risk.IsExternal {
			isExternal = "Yes"
		} else {
			isExternal = "No"
		}

		rows[i] = []interface{}{
			i + 1,
			risk.ReferenceNumber,
			riskCode,
			registrationDate,
			division,
			section,
			subSection,
			isStrategic,
			isOperational,
			expectedCause,
			consequence,
			isInternal,
			isExternal,
			riskTitle,
			description,
			riskCategory,
			responsibleUnit,
			responsibleUser,
			riskSource,
			riskStatus,
			likelihood,
			impact,
			riskRating,
			riskResponse,
			residualLikelihood,
			residualImpact,
			residualRiskRating,
			reportedAt,
			reviewFrequency,
			nextReviewDate,
			severityLevel,
		}
	}

	// Generate the Excel file
	return utility.GenerateExcelFile(headers, rows)
}

func (rrs *RiskRegisterService) SetRiskInherentEvaluationService(evaluationSvc RiskInherentEvaluationServiceInterface) {
	rrs.riskInherentEvaluationSvc = evaluationSvc
}

func (rrs *RiskRegisterService) SetRiskResidualEvaluationService(evaluationSvc RiskResidualEvaluationServiceInterface) {
	rrs.riskResidualEvaluationSvc = evaluationSvc
}

func (rrs *RiskRegisterService) GenerateRiskTemplate() (*os.File, error) {
	templateHeaders := []string{
		"RiskClassification", "ResponsibleUnit", "ResponsibleUser", "Title", "Description",
		"RiskCode", "RegistrationDate", "Division", "Section", "SubSection",
		"IsStrategic", "IsOperational", "ExpectedCause", "Consequence", "IsInternal", "IsExternal",
	}

	templateData := [][]interface{}{}

	// Fetch Risk Classifications (Categories)
	categories, err := rrs.riskCategorySvc.GetAll("", "", "", "active")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch risk classifications: %v", err)
	}
	categoryData := [][]interface{}{{"RiskClassification"}}
	for _, cat := range categories {
		categoryData = append(categoryData, []interface{}{cat.TitleEn})
	}

	// Fetch Responsible Units (Department - Level 2)
	responsibleUnits, err := rrs.organizationUnitSvc.GetAll(&settingsDTOs.OrganizationUnitListQueryParamsDTO{
		Status: "active",
		Level:  2,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch responsible units: %v", err)
	}
	responsibleUnitData := [][]interface{}{{"ResponsibleUnit"}}
	for _, unit := range responsibleUnits {
		responsibleUnitData = append(responsibleUnitData, []interface{}{unit.TitleEn})
	}

	// Fetch Users
	users, err := rrs.userSvc.GetAllUsers()
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

	// Fetch Organization Units (Divisions - Level 1)
	divisions, err := rrs.organizationUnitSvc.GetAll(&settingsDTOs.OrganizationUnitListQueryParamsDTO{
		Status: "active",
		Level:  1,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch divisions: %v", err)
	}
	divisionData := [][]interface{}{{"Division"}}
	for _, div := range divisions {
		divisionData = append(divisionData, []interface{}{div.TitleEn})
	}

	// Fetch Sections (Level 3)
	sections, err := rrs.organizationUnitSvc.GetAll(&settingsDTOs.OrganizationUnitListQueryParamsDTO{
		Status: "active",
		Level:  3,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch sections: %v", err)
	}
	sectionData := [][]interface{}{{"Section"}}
	for _, sec := range sections {
		sectionData = append(sectionData, []interface{}{sec.TitleEn})
	}

	// Fetch Sub Sections (Level 4)
	subSections, err := rrs.organizationUnitSvc.GetAll(&settingsDTOs.OrganizationUnitListQueryParamsDTO{
		Status: "active",
		Level:  4,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch sub sections: %v", err)
	}
	subSectionData := [][]interface{}{{"SubSection"}}
	for _, subSec := range subSections {
		subSectionData = append(subSectionData, []interface{}{subSec.TitleEn})
	}

	// Configure which columns should have dropdowns
	dropDownConfig := map[string]bool{
		"RiskClassification": true,
		"ResponsibleUnit":    true,
		"ResponsibleUser":    true,
		"Division":           true,
		"Section":            true,
		"SubSection":         true,
		"IsStrategic":        true,
		"IsOperational":      true,
		"IsInternal":         true,
		"IsExternal":         true,
	}

	// Add boolean dropdown data
	boolData := [][]interface{}{{"True", "False"}, {"true", "false"}}

	file, err := utility.GenerateExcelFileWithSheetsAndDropDown(map[string][][]interface{}{
		"Template":           append([][]interface{}{utility.ConvertToInterfaceSlice(templateHeaders)}, templateData...),
		"RiskClassification": categoryData,
		"ResponsibleUnit":    responsibleUnitData,
		"ResponsibleUser":    userData,
		"Division":           divisionData,
		"Section":            sectionData,
		"SubSection":         subSectionData,
		"IsStrategic":        boolData,
		"IsOperational":      boolData,
		"IsInternal":         boolData,
		"IsExternal":         boolData,
	}, dropDownConfig)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func (rrs *RiskRegisterService) ValidateRiskTemplate(file io.Reader) ([]map[string]interface{}, error) {
	rows, err := utility.ReadExcelFile(file)
	if err != nil {
		return nil, err
	}

	headers := []string{
		"RiskClassification", "ResponsibleUnit", "ResponsibleUser", "Title", "Description",
		"RiskCode", "RegistrationDate", "Division", "Section", "SubSection",
		"IsStrategic", "IsOperational", "ExpectedCause", "Consequence", "IsInternal", "IsExternal",
	}

	snakeHeaders := []string{
		"risk_classification", "responsible_unit", "responsible_user", "title", "description",
		"risk_code", "registration_date", "division", "section", "sub_section",
		"is_strategic", "is_operational", "expected_cause", "consequence", "is_internal", "is_external",
	}

	if len(rows) == 0 {
		return nil, fmt.Errorf("empty Excel file")
	}
	if !utility.ValidateHeaders(headers, rows[0]) {
		return nil, utility.NewValidationError(fmt.Errorf("invalid headers"), "invalid_header", "Risk Template")
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

	validatedData, err := rrs.validateRisks(input)
	if err != nil {
		return nil, err
	}

	return validatedData, nil
}

func (cs *RiskRegisterService) ValidateImportRisks(input []map[string]interface{}) ([]map[string]interface{}, error) {
	validatedData, err := cs.validateRisks(input)
	if err != nil {
		return nil, err
	}

	return validatedData, nil
}

func (rs *RiskRegisterService) ImportRisks(input []map[string]interface{}, currentUserID uint) ([]uint, error) {
	validatedData, err := rs.validateRisks(input)
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

		registrationDateStr, _ := record["registration_date"].(string)
		registrationDateStr = convertDateFormat(registrationDateStr)

		// DivisionID is required, default to 0 if not present
		divisionID := uint(0)
		if divID, ok := record["division_id"].(uint); ok {
			divisionID = divID
		}

		// Parse boolean fields from string to bool
		parseBool := func(val interface{}) bool {
			if str, ok := val.(string); ok {
				return strings.ToLower(str) == "true"
			}
			if b, ok := val.(bool); ok {
				return b
			}
			return false
		}

		dto := &riskDTOs.RiskRegisterFormDTO{
			RiskCategoryID:    utility.GetOptionalUintValue(record["risk_category_id"]),
			ResponsibleUnitID: utility.GetOptionalUintValue(record["responsible_unit_id"]),
			ResponsibleUserID: utility.GetOptionalUintValue(record["responsible_user_id"]),
			Title:             record["title"].(string),
			Description:       utility.GetOptionalStringValue(record["description"]),
			RiskCode:          utility.GetOptionalStringValue(record["risk_code"]),
			RegistrationDate:  registrationDateStr,
			DivisionID:        divisionID,
			SectionID:         utility.GetOptionalUintValue(record["section_id"]),
			SubSectionID:      utility.GetOptionalUintValue(record["sub_section_id"]),
			IsStrategic:       parseBool(record["is_strategic"]),
			IsOperational:     parseBool(record["is_operational"]),
			ExpectedCause:     utility.GetOptionalStringValue(record["expected_cause"]),
			Consequence:       utility.GetOptionalStringValue(record["consequence"]),
			IsInternal:        parseBool(record["is_internal"]),
			IsExternal:        parseBool(record["is_external"]),
		}

		id, err := rs.Create(currentUserID, dto)
		if err != nil {
			return nil, fmt.Errorf("failed to create risk for record: %+v, error: %v", record, err)
		}
		createdIDs = append(createdIDs, id)
	}

	return createdIDs, nil
}

func parseBoolValue(val interface{}) bool {
	if val == nil {
		return false
	}
	if b, ok := val.(bool); ok {
		return b
	}
	if str, ok := val.(string); ok {
		lowerValue := strings.ToLower(strings.TrimSpace(str))
		return lowerValue == "true" || lowerValue == "yes" || lowerValue == "1"
	}
	if num, ok := val.(float64); ok {
		return num == 1
	}
	if num, ok := val.(int); ok {
		return num == 1
	}
	return false
}

func (rs *RiskRegisterService) validateRisks(input []map[string]interface{}) ([]map[string]interface{}, error) {
	headers := []string{
		"risk_classification", "responsible_unit", "responsible_user", "title", "description",
		"risk_code", "registration_date", "division", "section", "sub_section",
		"is_strategic", "is_operational", "expected_cause", "consequence", "is_internal", "is_external",
	}

	// Fetch risk categories
	riskCategories, err := rs.riskCategorySvc.GetAll("", "", "", "active")
	if err != nil {
		return nil, fmt.Errorf("error fetching risk categories: %v", err)
	}
	riskCategoryTitles := make(map[string]uint)
	for _, cat := range riskCategories {
		riskCategoryTitles[strings.TrimSpace(cat.TitleEn)] = cat.ID
	}

	// Fetch Responsible Units (Department - Level 2)
	responsibleUnits, err := rs.organizationUnitSvc.GetAll(&settingsDTOs.OrganizationUnitListQueryParamsDTO{
		Status: "active",
		Level:  2,
	})
	if err != nil {
		return nil, fmt.Errorf("error fetching responsible units: %v", err)
	}
	responsibleUnitTitles := make(map[string]uint)
	for _, unit := range responsibleUnits {
		responsibleUnitTitles[strings.TrimSpace(unit.TitleEn)] = unit.ID
	}

	// Fetch users
	users, err := rs.userSvc.GetAllUsers()
	if err != nil {
		return nil, fmt.Errorf("error fetching users: %v", err)
	}
	userFullNames := make(map[string]uint)
	for _, user := range users {
		fullName := strings.TrimSpace(user.FirstNameEN + " " + user.LastNameEN)
		userFullNames[fullName] = uint(user.ID)
	}

	// Fetch Organization Units (Divisions - Level 1)
	divisions, err := rs.organizationUnitSvc.GetAll(&settingsDTOs.OrganizationUnitListQueryParamsDTO{
		Status: "active",
		Level:  1,
	})
	if err != nil {
		return nil, fmt.Errorf("error fetching divisions: %v", err)
	}
	divisionTitles := make(map[string]uint)
	for _, div := range divisions {
		divisionTitles[strings.TrimSpace(div.TitleEn)] = div.ID
	}

	// Fetch Sections (Level 3)
	sections, err := rs.organizationUnitSvc.GetAll(&settingsDTOs.OrganizationUnitListQueryParamsDTO{
		Status: "active",
		Level:  3,
	})
	if err != nil {
		return nil, fmt.Errorf("error fetching sections: %v", err)
	}
	sectionTitles := make(map[string]uint)
	for _, sec := range sections {
		sectionTitles[strings.TrimSpace(sec.TitleEn)] = sec.ID
	}

	// Fetch Sub Sections (Level 4)
	subSections, err := rs.organizationUnitSvc.GetAll(&settingsDTOs.OrganizationUnitListQueryParamsDTO{
		Status: "active",
		Level:  4,
	})
	if err != nil {
		return nil, fmt.Errorf("error fetching sub sections: %v", err)
	}
	subSectionTitles := make(map[string]uint)
	for _, subSec := range subSections {
		subSectionTitles[strings.TrimSpace(subSec.TitleEn)] = subSec.ID
	}

	// for checking unique risks
	// 1. Fetch existing titles from DB
	existingRisks, err := rs.repo.GetAllRiskTitles() // Implement this function to return []string
	if err != nil {
		return nil, fmt.Errorf("error fetching existing risk titles: %v", err)
	}
	existingTitleSet := make(map[string]bool)
	for _, title := range existingRisks {
		existingTitleSet[strings.TrimSpace(strings.ToLower(title))] = true
	}

	// 2. Track titles seen in the current Excel
	sheetTitleSet := make(map[string]bool)

	// Perform validation
	var validatedData []map[string]interface{}
	for _, record := range input {
		rowErrors := make(map[string]interface{})
		riskRecord := make(map[string]interface{})

		// Copy original data
		for key, value := range record {
			riskRecord[key] = value
		}

		for _, header := range headers {
			value, ok := record[header].(string)
			switch header {
			case "risk_classification":
				if !ok || value == "" {
					rowErrors[header] = "Risk Classification is required"
				} else if id, exists := riskCategoryTitles[strings.TrimSpace(value)]; !exists {
					rowErrors[header] = fmt.Sprintf("Invalid risk classification: %s", value)
				} else {
					riskRecord["risk_category_id"] = id
				}

			case "responsible_unit":
				if !ok || value == "" {
					rowErrors[header] = "Responsible Unit is required"
				} else if id, exists := responsibleUnitTitles[strings.TrimSpace(value)]; !exists {
					rowErrors[header] = fmt.Sprintf("Invalid responsible unit: %s", value)
				} else {
					riskRecord["responsible_unit_id"] = id
				}

			case "responsible_user":
				if !ok || value == "" {
					rowErrors[header] = "Responsible User is required"
				} else if id, exists := userFullNames[strings.TrimSpace(value)]; !exists {
					rowErrors[header] = fmt.Sprintf("Invalid responsible user: %s", value)
				} else {
					riskRecord["responsible_user_id"] = id
				}

			case "title":
				if !ok || strings.TrimSpace(value) == "" {
					rowErrors[header] = "Title is required"
				} else {
					trimmedTitle := strings.ToLower(strings.TrimSpace(value))
					if _, exists := sheetTitleSet[trimmedTitle]; exists {
						rowErrors[header] = "Duplicate title in the sheet"
					} else if _, exists := existingTitleSet[trimmedTitle]; exists {
						rowErrors[header] = "Title already exists in the system"
					} else {
						sheetTitleSet[trimmedTitle] = true
					}
				}

			case "description":

			case "risk_code":
				// Optional field, just store the value
				if value != "" {
					riskRecord["risk_code"] = value
				}

			case "registration_date":
				// Required field
				if !ok || value == "" {
					rowErrors[header] = "Registration Date is required"
				} else {
					riskRecord["registration_date"] = value
				}

			case "division":
				// Required field
				if !ok || value == "" {
					rowErrors[header] = "Division is required"
				} else if id, exists := divisionTitles[strings.TrimSpace(value)]; !exists {
					rowErrors[header] = fmt.Sprintf("Invalid division: %s", value)
				} else {
					riskRecord["division_id"] = id
				}

			case "section":
				if value != "" {
					if id, exists := sectionTitles[strings.TrimSpace(value)]; !exists {
						rowErrors[header] = fmt.Sprintf("Invalid section: %s", value)
					} else {
						riskRecord["section_id"] = id
					}
				}

			case "sub_section":
				if value != "" {
					if id, exists := subSectionTitles[strings.TrimSpace(value)]; !exists {
						rowErrors[header] = fmt.Sprintf("Invalid sub section: %s", value)
					} else {
						riskRecord["sub_section_id"] = id
					}
				}

			case "is_strategic", "is_operational", "is_internal", "is_external":
				riskRecord[header] = parseBoolValue(record[header])

			case "expected_cause":
				// Optional field
				if value != "" {
					riskRecord["expected_cause"] = value
				}

			case "consequence":
				// Optional field
				if value != "" {
					riskRecord["consequence"] = value
				}

			}

		}

		// Validate organization unit hierarchy after all IDs are collected
		// Note: IsUnitUnderParent returns bool without error. If DB lookup fails internally,
		// it returns false which will be treated as validation error. This is acceptable as
		// the units have already been validated to exist in the maps above.
		if len(rowErrors) == 0 {
			// Validate that responsible unit is under division
			if divisionID, hasDivision := riskRecord["division_id"].(uint); hasDivision && divisionID > 0 {
				if responsibleUnitID, hasResponsibleUnit := riskRecord["responsible_unit_id"].(uint); hasResponsibleUnit && responsibleUnitID > 0 {
					if !rs.organizationUnitSvc.IsUnitUnderParent(responsibleUnitID, divisionID) {
						rowErrors["responsible_unit"] = "Responsible unit must be under the selected division"
					}
				}
			}

			// Validate that section is under responsible unit
			if sectionID, hasSection := riskRecord["section_id"].(uint); hasSection && sectionID > 0 {
				if responsibleUnitID, hasResponsibleUnit := riskRecord["responsible_unit_id"].(uint); hasResponsibleUnit && responsibleUnitID > 0 {
					if !rs.organizationUnitSvc.IsUnitUnderParent(sectionID, responsibleUnitID) {
						rowErrors["section"] = "Section must be under the responsible unit"
					}
				}
			}

			// Validate that sub section is under section
			if subSectionID, hasSubSection := riskRecord["sub_section_id"].(uint); hasSubSection && subSectionID > 0 {
				if sectionID, hasSection := riskRecord["section_id"].(uint); hasSection && sectionID > 0 {
					if !rs.organizationUnitSvc.IsUnitUnderParent(subSectionID, sectionID) {
						rowErrors["sub_section"] = "Sub section must be under the section"
					}
				}
			}
		}

		if len(rowErrors) == 0 {
			riskRecord["errors"] = nil
		} else {
			for _, header := range headers {
				if _, exists := rowErrors[header]; !exists {
					rowErrors[header] = ""
				}
			}
			riskRecord["errors"] = rowErrors
		}

		validatedData = append(validatedData, riskRecord)
	}

	return validatedData, nil
}

func (rrs *RiskRegisterService) GetAllIdentifiedRisks(params *riskDTOs.RiskRegisterListQueryParamsDTO) (*riskDTOs.RiskRegisterListResponseDTO, error) {
	return rrs.repo.GetAllIdentifiedRisks(params)
}

func (rrs *RiskRegisterService) ForceDelete(ids []uint, currentUserID uint) ([]uint, error) {
	tx := rrs.repo.BeginTransaction()
	var failedIds []uint
	var deletedData []map[string]interface{}

	for _, id := range ids {
		if err := rrs.repo.ForceDelete(id, tx); err != nil {
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
		err := rrs.activityLogSvc.LogBulkActivities(
			models.DeleteAction, // Make sure this action constant is defined appropriately
			models.RiskRegister{}.TableName(),
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

func (rrs RiskRegisterService) ExportRiskRegisterPDF(id uint, currentUserID uint) (string, error) {
	// Validate that template images are configured
	if _, err := rrs.riskTemplateSettingSvc.GetImagePath("logo"); err != nil {
		return "", utility.NewForbiddenError("logo_image_missing")
	}
	if _, err := rrs.riskTemplateSettingSvc.GetImagePath("cover-page-bg"); err != nil {
		return "", utility.NewForbiddenError("cover_image_missing")
	}

	riskDTO, err := rrs.GetFullDetailsByID(id)
	if err != nil {
		return "", err
	}

	// Inline build of export view (align with incident pattern)
	createdAt := riskDTO.CreatedAt.Format("02 Jan 2006")
	generatedAt := time.Now().Format("Jan 02, 2006")
	preparedFirst := riskDTO.CreatedByUser.FirstNameEN
	preparedLast := riskDTO.CreatedByUser.LastNameEN

	var orgUnit, orgUnitAr string
	if riskDTO.ResponsibleUnit != nil {
		orgUnit = riskDTO.ResponsibleUnit.TitleEn
		orgUnitAr = riskDTO.ResponsibleUnit.TitleAr
	}

	var divisionTitleEn, divisionTitleAr string
	if riskDTO.Division != nil {
		divisionTitleEn = riskDTO.Division.TitleEn
		divisionTitleAr = riskDTO.Division.TitleAr
	}

	var sectionTitleEn, sectionTitleAr string
	if riskDTO.Section != nil {
		sectionTitleEn = riskDTO.Section.TitleEn
		sectionTitleAr = riskDTO.Section.TitleAr
	}

	reportedAt := ""
	if riskDTO.ReportedAtString != nil {
		reportedAt = *riskDTO.ReportedAtString
	}

	// var riskSourceTitleEn *string
	// var riskSourceTitleAr *string
	// if len(riskDTO.RiskRegisterSources) > 0 {
	// 	var enList []string
	// 	var arList []string

	// 	for _, rs := range riskDTO.RiskRegisterSources {
	// 		enList = append(enList, rs.RiskSource.TitleEn)
	// 		arList = append(arList, rs.RiskSource.TitleAr)
	// 	}

	// 	enStr := strings.Join(enList, ", ")
	// 	arStr := strings.Join(arList, ", ")

	// 	riskSourceTitleEn = &enStr
	// 	riskSourceTitleAr = &arStr
	// }

	var riskResponseStrategyTitleEn *string
	var riskResponseStrategyTitleAr *string
	var riskResponseStrategyDescription *string
	var riskResponseStrategyTreatmentGuidance *string
	if riskDTO.RiskResponseStrategy != nil {
		riskResponseStrategyTitleEn = &riskDTO.RiskResponseStrategy.TitleEn
		riskResponseStrategyTitleAr = &riskDTO.RiskResponseStrategy.TitleAr
	}

	var responsibleUserFirstNameEn, responsibleUserLastNameEn, responsibleUserInitial string
	if riskDTO.ResponsibleUser != nil {
		responsibleUserFirstNameEn = riskDTO.ResponsibleUser.FirstNameEN
		responsibleUserLastNameEn = riskDTO.ResponsibleUser.LastNameEN
		responsibleUserInitial = riskDTO.ResponsibleUser.Initial
	}

	var currentRiskRatingTitleEn *string
	var currentRiskRatingTitleAr *string
	var currentRiskRatingColorCode *string
	var currentRiskRatingValue *uint
	var currentLikelihoodTitleEn *string
	var currentLikelihoodTitleAr *string
	var currentLikelihoodScore *uint
	var currentLikelihoodDescriptionEn *string
	var currentImpactTitleEn *string
	var currentImpactTitleAr *string
	var currentImpactScore *uint
	var currentImpactDescriptionEn *string
	if riskDTO.CurrentRiskRating != nil {
		currentRiskRatingTitleEn = &riskDTO.CurrentRiskRating.TitleEn
		currentRiskRatingTitleAr = &riskDTO.CurrentRiskRating.TitleAr
		currentRiskRatingColorCode = &riskDTO.CurrentRiskRating.ColorCode
	}
	if riskDTO.CurrentRiskRatingValue != nil {
		currentRiskRatingValue = &riskDTO.CurrentRiskRatingValue.Value
	}
	if riskDTO.CurrentLikelihood != nil {
		currentLikelihoodTitleEn = &riskDTO.CurrentLikelihood.TitleEn
		currentLikelihoodTitleAr = &riskDTO.CurrentLikelihood.TitleAr
		currentLikelihoodScore = &riskDTO.CurrentLikelihood.Score
		currentLikelihoodDescriptionEn = &riskDTO.CurrentLikelihood.DescriptionEn
	}
	if riskDTO.CurrentImpact != nil {
		currentImpactTitleEn = &riskDTO.CurrentImpact.TitleEn
		currentImpactTitleAr = &riskDTO.CurrentImpact.TitleAr
		currentImpactScore = &riskDTO.CurrentImpact.Score
		currentImpactDescriptionEn = &riskDTO.CurrentImpact.DescriptionEn
	}

	var inherentRiskRatingTitleEn *string
	var inherentRiskRatingTitleAr *string
	var inherentRiskRatingColorCode *string
	if riskDTO.InherentRiskRating != nil {
		inherentRiskRatingTitleEn = &riskDTO.InherentRiskRating.TitleEn
		inherentRiskRatingTitleAr = &riskDTO.InherentRiskRating.TitleAr
		inherentRiskRatingColorCode = &riskDTO.InherentRiskRating.ColorCode
	}

	var inherentRiskRatingValue *uint
	if riskDTO.InherentRiskRatingValue != nil {
		inherentRiskRatingValue = &riskDTO.InherentRiskRatingValue.Value
	}

	var inherentLikelihoodTitleEn *string
	var inherentLikelihoodTitleAr *string
	var inherentLikelihoodScore *uint
	var inherentLikelihoodDescriptionEn *string
	if riskDTO.InherentLikelihood != nil {
		inherentLikelihoodTitleEn = &riskDTO.InherentLikelihood.TitleEn
		inherentLikelihoodTitleAr = &riskDTO.InherentLikelihood.TitleAr
		inherentLikelihoodScore = &riskDTO.InherentLikelihood.Score
		inherentLikelihoodDescriptionEn = &riskDTO.InherentLikelihood.DescriptionEn
	}

	var inherentImpactTitleEn *string
	var inherentImpactTitleAr *string
	var inherentImpactScore *uint
	var inherentImpactDescriptionEn *string
	if riskDTO.InherentImpact != nil {
		inherentImpactTitleEn = &riskDTO.InherentImpact.TitleEn
		inherentImpactTitleAr = &riskDTO.InherentImpact.TitleAr
		inherentImpactScore = &riskDTO.InherentImpact.Score
		inherentImpactDescriptionEn = &riskDTO.InherentImpact.DescriptionEn
	}

	// Fetch inherent evaluation to get impact category and guideline
	var inherentImpactCategoryTitleEn *string
	var inherentImpactCategoryTitleAr *string
	var inherentImpactGuideline *string
	if riskDTO.InherentImpact != nil {
		inherentEvaluation, err := rrs.riskInherentEvaluationSvc.GetFullDetailsByID(id)
		if err == nil && inherentEvaluation != nil && len(inherentEvaluation.RiskImpactAnalyses) > 0 {
			// Get the first impact analysis (typically there's one per evaluation)
			impactAnalysis := inherentEvaluation.RiskImpactAnalyses[0]
			inherentImpactCategoryTitleEn = &impactAnalysis.RiskImpactCategory.TitleEn
			inherentImpactCategoryTitleAr = &impactAnalysis.RiskImpactCategory.TitleAr

			// Get the guideline from the first guideline entry
			if len(impactAnalysis.RiskImpactAnalysisGuidelines) > 0 {
				guidelineEntry := impactAnalysis.RiskImpactAnalysisGuidelines[0]
				inherentImpactGuideline = &guidelineEntry.RiskMatrixImpactGuideline.Guideline
			}
		}
	}

	var residualRiskRatingTitleEn *string
	var residualRiskRatingTitleAr *string
	var residualRiskRatingColorCode *string
	if riskDTO.ResidualRiskRating != nil {
		residualRiskRatingTitleEn = &riskDTO.ResidualRiskRating.TitleEn
		residualRiskRatingTitleAr = &riskDTO.ResidualRiskRating.TitleAr
		residualRiskRatingColorCode = &riskDTO.ResidualRiskRating.ColorCode
	}

	var residualRiskRatingValue *uint
	if riskDTO.ResidualRiskRatingValue != nil {
		residualRiskRatingValue = &riskDTO.ResidualRiskRatingValue.Value
	}

	var residualLikelihoodTitleEn *string
	var residualLikelihoodTitleAr *string
	var residualLikelihoodScore *uint
	var residualLikelihoodDescriptionEn *string
	if riskDTO.ResidualLikelihood != nil {
		residualLikelihoodTitleEn = &riskDTO.ResidualLikelihood.TitleEn
		residualLikelihoodTitleAr = &riskDTO.ResidualLikelihood.TitleAr
		residualLikelihoodScore = &riskDTO.ResidualLikelihood.Score
		residualLikelihoodDescriptionEn = &riskDTO.ResidualLikelihood.DescriptionEn
	}

	var residualImpactTitleEn *string
	var residualImpactTitleAr *string
	var residualImpactScore *uint
	var residualImpactDescriptionEn *string
	if riskDTO.ResidualImpact != nil {
		residualImpactTitleEn = &riskDTO.ResidualImpact.TitleEn
		residualImpactTitleAr = &riskDTO.ResidualImpact.TitleAr
		residualImpactScore = &riskDTO.ResidualImpact.Score
		residualImpactDescriptionEn = &riskDTO.ResidualImpact.DescriptionEn
	}

	var severityLevelTitleEn *string
	var severityLevelTitleAr *string
	if riskDTO.SeverityLevel != nil {
		severityLevelTitleEn = &riskDTO.SeverityLevel.TitleEn
		severityLevelTitleAr = &riskDTO.SeverityLevel.TitleAr
	}

	var reviewFrequencyTitleEn *string
	var reviewFrequencyTitleAr *string
	if riskDTO.ReviewFrequency != nil {
		reviewFrequencyTitleEn = &riskDTO.ReviewFrequency.TitleEn
		reviewFrequencyTitleAr = &riskDTO.ReviewFrequency.TitleAr
	}

	var submittedByUserFirstNameEn *string
	var submittedByUserLastNameEn *string
	if riskDTO.SubmittedByUser != nil {
		submittedByUserFirstNameEn = &riskDTO.SubmittedByUser.FirstNameEN
		submittedByUserLastNameEn = &riskDTO.SubmittedByUser.LastNameEN
	}
	var updatedByUserFirstNameEn *string
	var updatedByUserLastNameEn *string
	if riskDTO.UpdatedByUser.FirstNameEN != "" {
		updatedByUserFirstNameEn = &riskDTO.UpdatedByUser.FirstNameEN
		updatedByUserLastNameEn = &riskDTO.UpdatedByUser.LastNameEN
	}
	var updatedAtString *string
	if riskDTO.UpdatedAt != nil {
		updatedAtStr := riskDTO.UpdatedAt.Format("2006-01-02 15:04:05")
		updatedAtString = &updatedAtStr
	}

	var existingControls *string
	if len(riskDTO.Controls) > 0 {
		controlTitles := make([]string, 0, len(riskDTO.Controls))
		for _, control := range riskDTO.Controls {
			if control.Control.TitleEn != "" {
				controlTitles = append(controlTitles, control.Control.TitleEn)
			}
		}
		if len(controlTitles) > 0 {
			controlsStr := strings.Join(controlTitles, ", ")
			existingControls = &controlsStr
		}
	}

	riskControls, err := rrs.riskReportRepo.GetAllRiskControls(id)
	if err != nil {
		return "", fmt.Errorf("failed to fetch controls: %w", err)
	}

	controls := make([]riskDTOs.RiskRegisterExportControlDTO, 0, len(riskControls))
	for _, control := range riskControls {
		var testScore *uint
		var implementationStatusTitleEn *string
		var controlEffectivenessRatingTitleEn *string
		var testedAtString *string
		implementationStatusTitleEnValue := ""
		controlEffectivenessRatingTitleEnValue := ""

		if control.ControlEffectivenessRatingTitleEn != "" {
			controlEffectivenessRatingTitleEn = &control.ControlEffectivenessRatingTitleEn
			controlEffectivenessRatingTitleEnValue = control.ControlEffectivenessRatingTitleEn
		}
		if control.ControlImplementationStatusTitleEn != "" {
			implementationStatusTitleEn = &control.ControlImplementationStatusTitleEn
			implementationStatusTitleEnValue = control.ControlImplementationStatusTitleEn
		}
		if control.RiskRatingValue != nil {
			testScore = control.RiskRatingValue
		}
		if control.LastTestDate != nil {
			testDateStr := control.LastTestDate.Format("2006-01-02")
			testedAtString = &testDateStr
		} else if control.LastReviewDate != nil {
			reviewDateStr := control.LastReviewDate.Format("2006-01-02")
			testedAtString = &reviewDateStr
		}

		code := ""
		if control.Code != nil {
			code = *control.Code
		}
		var controlTypeTitleEn *string
		if control.ControlTypeTitleEn != "" {
			controlTypeTitleEn = &control.ControlTypeTitleEn
		}
		var controlStatusTitleEn *string
		var controlStatusTitleAr *string
		var controlStatusColorCode *string
		var controlStatusType *string
		var controlTestDescription *string
		controlStatusTitleEnValue := ""
		controlStatusTypeValue := ""

		if control.ControlStatusTitleEn != "" {
			controlStatusTitleEn = &control.ControlStatusTitleEn
			controlStatusTitleEnValue = control.ControlStatusTitleEn
		}
		if control.ControlStatusTitleAr != "" {
			controlStatusTitleAr = &control.ControlStatusTitleAr
		}
		if control.ControlStatusColorCode != "" {
			controlStatusColorCode = &control.ControlStatusColorCode
		}
		if control.ControlStatusType != "" {
			controlStatusType = &control.ControlStatusType
			controlStatusTypeValue = control.ControlStatusType
		}
		if control.ControlTestDescription != "" {
			controlTestDescription = &control.ControlTestDescription
		}

		controlDTO := riskDTOs.RiskRegisterExportControlDTO{
			ReferenceNumber:                        control.ReferenceNumber,
			TitleEn:                                control.TitleEn,
			Code:                                   code,
			ControlTypeTitleEn:                     controlTypeTitleEn,
			TestScore:                              testScore,
			ImplementationStatusTitleEn:            implementationStatusTitleEn,
			ImplementationStatusTitleEnValue:       implementationStatusTitleEnValue,
			ControlEffectivenessRatingTitleEn:      controlEffectivenessRatingTitleEn,
			ControlEffectivenessRatingTitleEnValue: controlEffectivenessRatingTitleEnValue,
			ControlEffectivenessRatingColorCode:    control.ControlEffectivenessRatingColorCode,
			ControlEffectivenessRatingType:         control.ControlEffectivenessRatingType,
			ControlEffectivenessScore:              control.ControlEffectivenessScore,
			RiskReduction:                          control.RiskReduction,
			TestedAtString:                         testedAtString,
			TestedByUserFirstNameEn:                control.TestResponsibleUserFirstNameEn,
			TestedByUserLastNameEn:                 control.TestResponsibleUserLastNameEn,
			TestedByUserInitial:                    control.TestResponsibleUserInitial,
			ControlStatusTitleEn:                   controlStatusTitleEn,
			ControlStatusTitleEnValue:              controlStatusTitleEnValue,
			ControlStatusTitleAr:                   controlStatusTitleAr,
			ControlStatusColorCode:                 controlStatusColorCode,
			ControlStatusType:                      controlStatusType,
			ControlStatusTypeValue:                 controlStatusTypeValue,
			ControlTestDescription:                 controlTestDescription,
		}
		controls = append(controls, controlDTO)
	}

	// Compute Control Statistics
	var controlCoveragePercent *float64
	var ineffectiveControlsCount *uint
	var implementedPercent *float64
	var effectiveControlsPercent *float64
	var effectiveControlsCount *uint
	var controlScoreAchieved *uint
	var controlScoreRemaining *uint

	if len(controls) > 0 {
		totalControls := float64(len(controls))
		coverage := 100.0
		controlCoveragePercent = &coverage

		implementedCount := 0
		ineffectiveCount := 0
		effectiveCount := 0
		totalScore := uint(0)

		for _, ctrl := range controls {
			if ctrl.ImplementationStatusTitleEn != nil && *ctrl.ImplementationStatusTitleEn != "" {
				implementedCount++
			}
			if ctrl.TestScore != nil {
				totalScore += *ctrl.TestScore
			}
		}

		implPercent := (float64(implementedCount) / totalControls) * 100
		effectivePercent := float64((effectiveCount * 100) / int(totalControls))
		ineffectiveCountUint := uint(ineffectiveCount)
		effectiveCountUint := uint(effectiveCount)
		maxScore := totalControls * 100
		remainingScore := uint(maxScore) - totalScore

		implementedPercent = &implPercent
		effectiveControlsPercent = &effectivePercent
		ineffectiveControlsCount = &ineffectiveCountUint
		effectiveControlsCount = &effectiveCountUint
		controlScoreAchieved = &totalScore
		controlScoreRemaining = &remainingScore
	}

	// Fetch Control Effectiveness Ratings
	controlEffectivenessRatings := []riskDTOs.ControlEffectivenessRatingExportDTO{
		{TitleEn: "Control missing", TitleAr: "Control missing", MinScore: 0, MaxScore: 0.1, ReductionPercent: 0, ColorCode: "#FF0000"},
		{TitleEn: "Improvement Needed", TitleAr: "Improvement Needed", MinScore: 2, MaxScore: 3.5, ReductionPercent: 55, ColorCode: "#FFC107"},
		{TitleEn: "Ineffective", TitleAr: "Ineffective", MinScore: 0.1, MaxScore: 2, ReductionPercent: 25, ColorCode: "#DC3545"},
		{TitleEn: "Partially Effective", TitleAr: "Partially Effective", MinScore: 3.5, MaxScore: 4.99, ReductionPercent: 65, ColorCode: "#90DD75"},
		{TitleEn: "Very Effective", TitleAr: "Very Effective", MinScore: 5, MaxScore: 5, ReductionPercent: 75, ColorCode: "#28A745"},
	}

	// Get AverageControlEffectivenessScore from residual evaluation
	var averageControlEffectivenessScore *float64
	residualEval, err := rrs.riskResidualEvaluationSvc.GetFullDetailsByID(id)
	if err == nil && residualEval != nil {
		averageControlEffectivenessScore = &residualEval.AverageControlEffectivenessScore
	}

	riskAssets, err := rrs.riskReportRepo.GetAllRiskAssets(id)
	if err != nil {
		return "", fmt.Errorf("failed to fetch assets: %w", err)
	}
	assets := make([]riskDTOs.RiskRegisterExportAssetDTO, 0, len(riskAssets))
	for _, asset := range riskAssets {
		var assetTypeTitleEn *string
		if asset.AssetCategoryTitleEn != nil && *asset.AssetCategoryTitleEn != "" {
			assetTypeTitleEn = asset.AssetCategoryTitleEn
		}
		var categoryTitleEn *string
		if asset.AssetCategoryTitleEn != nil && *asset.AssetCategoryTitleEn != "" {
			categoryTitleEn = asset.AssetCategoryTitleEn
		}
		assetDTO := riskDTOs.RiskRegisterExportAssetDTO{
			ReferenceNumber:            asset.ReferenceNumber,
			TitleEn:                    asset.TitleEn,
			TitleAr:                    asset.TitleAr,
			AssetTypeTitleEn:           assetTypeTitleEn,
			CategoryTitleEn:            categoryTitleEn,
			Description:                asset.Description,
			ResponsibleUserFirstNameEn: asset.ResponsibleUserFirstNameEn,
			ResponsibleUserLastNameEn:  asset.ResponsibleUserLastNameEn,
		}
		assets = append(assets, assetDTO)
	}

	documents, err := rrs.riskReportRepo.GetAllRiskDocuments(id)
	if err != nil {
		return "", fmt.Errorf("failed to fetch documents: %w", err)
	}

	processes, err := rrs.riskReportRepo.GetAllRiskProcesses(id)
	if err != nil {
		return "", fmt.Errorf("failed to fetch processes: %w", err)
	}

	kpis, err := rrs.riskReportRepo.GetAllRiskKPIs(id)
	if err != nil {
		return "", fmt.Errorf("failed to fetch KPIs: %w", err)
	}

	strategicObjectives, err := rrs.riskReportRepo.GetAllRiskStrategicObjectives(id)
	if err != nil {
		return "", fmt.Errorf("failed to fetch strategic objectives: %w", err)
	}

	for i := range strategicObjectives {

		if strategicObjectives[i].StartDate != nil && *strategicObjectives[i].StartDate != "" {
			formattedDate, err := time.Parse(time.RFC3339, *strategicObjectives[i].StartDate)
			if err != nil {
				return "", fmt.Errorf("failed to parse start date: %w", err)
			}
			formattedDateStr := formattedDate.Format("02 Jan 2006")
			strategicObjectives[i].StartDate = &formattedDateStr
		}

		if strategicObjectives[i].EndDate != nil && *strategicObjectives[i].EndDate != "" {
			formattedDate, err := time.Parse(time.RFC3339, *strategicObjectives[i].EndDate)
			if err != nil {
				return "", fmt.Errorf("failed to parse end date: %w", err)
			}
			formattedDateStr := formattedDate.Format("02 Jan 2006")
			strategicObjectives[i].EndDate = &formattedDateStr
		}
	}

	complianceRequirements, err := rrs.riskReportRepo.GetAllRiskComplianceRequirements(id)
	if err != nil {
		return "", fmt.Errorf("failed to fetch compliance requirements: %w", err)
	}

	incidents, err := rrs.riskReportRepo.GetAllRiskIncidents(id)
	if err != nil {
		return "", fmt.Errorf("failed to fetch incidents: %w", err)
	}

	rootCauseAnalyses, err := fetchRootCauseAnalysesForExport(id, rrs.riskReportRepo)
	if err != nil {
		return "", fmt.Errorf("failed to fetch root cause analyses: %w", err)
	}

	actionPlans, err := fetchActionPlansForExport(id, riskDTO.ReferenceNumber, riskDTO.Title, rrs.riskReportRepo)
	if err != nil {
		return "", fmt.Errorf("failed to fetch action plans: %w", err)
	}

	kriDTOs, err := rrs.riskReportRepo.GetKRIsByRiskID(id)
	if err != nil {
		return "", fmt.Errorf("failed to fetch KRIs: %w", err)
	}

	kris := make([]riskDTOs.RiskReportKriResponseDTO, 0, len(kriDTOs))
	for _, kri := range kriDTOs {
		var unitTitleEnPtr *string
		var unitTitleArPtr *string
		if kri.UnitTitleEn != "" {
			unitTitleEnPtr = &kri.UnitTitleEn
		}
		if kri.UnitTitleAr != "" {
			unitTitleArPtr = &kri.UnitTitleAr
		}

		lowThresholdPtr := &kri.LowThreshold
		mediumThresholdPtr := &kri.MediumThreshold
		currentValuePtr := &kri.CurrentValue

		kriResponse := riskDTOs.RiskReportKriResponseDTO{
			ID:                         0,
			RiskReportIdentificationID: 0,
			RiskReferenceNumber:        riskDTO.ReferenceNumber,
			RiskTitle:                  riskDTO.Title,
			Title:                      kri.Title,
			Description:                kri.Description,
			LowThreshold:               lowThresholdPtr,
			MediumThreshold:            mediumThresholdPtr,
			UnitTitleEn:                unitTitleEnPtr,
			UnitTitleAr:                unitTitleArPtr,
			CurrentValue:               currentValuePtr,
		}
		kris = append(kris, kriResponse)
	}

	// Fetch Risk Appetite Statements (Inherent - for Current Risk Analysis)
	riskAppetiteStatementsRaw, err := rrs.riskReportRepo.GetAllRiskAppetiteStatements(id, "inherent")
	if err != nil {
		return "", fmt.Errorf("failed to fetch inherent risk appetite statements: %w", err)
	}
	riskAppetiteStatements := make([]riskDTOs.RiskRegisterExportRiskAppetiteStatementDTO, 0, len(riskAppetiteStatementsRaw))
	for _, stmt := range riskAppetiteStatementsRaw {
		riskAppetiteStatements = append(riskAppetiteStatements, *stmt)
	}

	// Fetch Risk Appetite Statements (Residual - for Residual Risk Analysis)
	residualRiskAppetiteStatementsRaw, err := rrs.riskReportRepo.GetAllRiskAppetiteStatements(id, "residual")
	if err != nil {
		return "", fmt.Errorf("failed to fetch residual risk appetite statements: %w", err)
	}
	residualRiskAppetiteStatements := make([]riskDTOs.RiskRegisterExportRiskAppetiteStatementDTO, 0, len(residualRiskAppetiteStatementsRaw))
	for _, stmt := range residualRiskAppetiteStatementsRaw {
		residualRiskAppetiteStatements = append(residualRiskAppetiteStatements, *stmt)
	}

	// Fetch Risk Closure (get latest closure)
	var closureJustification *string
	var lessonsLearned *string
	var whatWorkedWell *string
	var whatNeedsImprovement *string
	var keyTakeaways *string

	riskClosure, err := rrs.riskReportRepo.GetLatestRiskClosure(id)
	if err != nil {
		return "", fmt.Errorf("failed to fetch risk closure: %w", err)
	}
	if riskClosure != nil {
		closureJustification = riskClosure.Justification
		lessonsLearned = riskClosure.LessonsLearned
		whatWorkedWell = riskClosure.WhatWorkedWell
		whatNeedsImprovement = riskClosure.WhatNeedsImprovement
		keyTakeaways = riskClosure.KeyTakeaways
	}

	// Reviews - placeholder for now (need to check what Reviews are)
	var reviews []riskDTOs.RiskRegisterExportReviewDTO

	// Calculate risk pointer positions
	currentRiskPointerPosition := rrs.getRiskPointerPosition(currentRiskRatingTitleEn)
	inherentRiskPointerPosition := rrs.getRiskPointerPosition(inherentRiskRatingTitleEn)
	residualRiskPointerPosition := rrs.getRiskPointerPosition(residualRiskRatingTitleEn)

	risk := riskDTOs.RiskRegisterExportRiskDTO{
		ReferenceNumber:     riskDTO.ReferenceNumber,
		Title:               riskDTO.Title,
		Description:         riskDTO.Description,
		ReportedAt:          reportedAt,
		RiskCategoryTitleEn: riskDTO.RiskCategory.TitleEn,
		RiskCategoryTitleAr: riskDTO.RiskCategory.TitleAr,
		// RiskSourceTitleEn:               riskSourceTitleEn,
		// RiskSourceTitleAr:               riskSourceTitleAr,
		ResponsibleUnitTitleEn:          orgUnit,
		ResponsibleUnitTitleAr:          orgUnitAr,
		ResponsibleUserFirstNameEn:      responsibleUserFirstNameEn,
		ResponsibleUserLastNameEn:       responsibleUserLastNameEn,
		ResponsibleUserInitial:          responsibleUserInitial,
		CurrentRiskRatingTitleEn:        currentRiskRatingTitleEn,
		CurrentRiskRatingTitleAr:        currentRiskRatingTitleAr,
		CurrentRiskRatingColorCode:      currentRiskRatingColorCode,
		CurrentRiskRatingValue:          currentRiskRatingValue,
		CurrentLikelihoodTitleEn:        currentLikelihoodTitleEn,
		CurrentLikelihoodTitleAr:        currentLikelihoodTitleAr,
		CurrentLikelihoodScore:          currentLikelihoodScore,
		CurrentLikelihoodDescriptionEn:  currentLikelihoodDescriptionEn,
		CurrentImpactTitleEn:            currentImpactTitleEn,
		CurrentImpactTitleAr:            currentImpactTitleAr,
		CurrentImpactScore:              currentImpactScore,
		CurrentImpactDescriptionEn:      currentImpactDescriptionEn,
		InherentRiskRatingTitleEn:       inherentRiskRatingTitleEn,
		InherentRiskRatingTitleAr:       inherentRiskRatingTitleAr,
		InherentRiskRatingColorCode:     inherentRiskRatingColorCode,
		InherentRiskRatingValue:         inherentRiskRatingValue,
		InherentRiskPointerPosition:     inherentRiskPointerPosition,
		InherentLikelihoodTitleEn:       inherentLikelihoodTitleEn,
		InherentLikelihoodTitleAr:       inherentLikelihoodTitleAr,
		InherentLikelihoodScore:         inherentLikelihoodScore,
		InherentLikelihoodDescriptionEn: inherentLikelihoodDescriptionEn,
		InherentImpactTitleEn:           inherentImpactTitleEn,
		InherentImpactTitleAr:           inherentImpactTitleAr,
		InherentImpactScore:             inherentImpactScore,
		InherentImpactDescriptionEn:     inherentImpactDescriptionEn,
		InherentImpactCategoryTitleEn:   inherentImpactCategoryTitleEn,
		InherentImpactCategoryTitleAr:   inherentImpactCategoryTitleAr,
		InherentImpactGuideline:         inherentImpactGuideline,
		ResidualRiskRatingTitleEn:       residualRiskRatingTitleEn,
		ResidualRiskRatingTitleAr:       residualRiskRatingTitleAr,
		ResidualRiskRatingColorCode:     residualRiskRatingColorCode,
		ResidualRiskRatingValue:         residualRiskRatingValue,
		ResidualLikelihoodTitleEn:       residualLikelihoodTitleEn,
		ResidualLikelihoodTitleAr:       residualLikelihoodTitleAr,
		ResidualLikelihoodScore:         residualLikelihoodScore,
		ResidualLikelihoodDescriptionEn: residualLikelihoodDescriptionEn,
		ResidualImpactTitleEn:           residualImpactTitleEn,
		ResidualImpactTitleAr:           residualImpactTitleAr,
		ResidualImpactScore:             residualImpactScore,
		ResidualImpactDescriptionEn:     residualImpactDescriptionEn,
		SeverityLevelTitleEn:            severityLevelTitleEn,
		SeverityLevelTitleAr:            severityLevelTitleAr,
		ReviewFrequencyTitleEn:          reviewFrequencyTitleEn,
		ReviewFrequencyTitleAr:          reviewFrequencyTitleAr,
		NextReviewDateString:            riskDTO.NextReviewDateString,
		LastReviewNote:                  riskDTO.LastReviewNote,
		NextReviewLevel:                 riskDTO.NextReviewLevel,
		SubmittedByUserFirstNameEn:      submittedByUserFirstNameEn,
		SubmittedByUserLastNameEn:       submittedByUserLastNameEn,
		SubmittedAtString:               riskDTO.SubmittedAtString,
		UpdatedByUserFirstNameEn:        utility.GetStringValue(updatedByUserFirstNameEn),
		UpdatedByUserLastNameEn:         utility.GetStringValue(updatedByUserLastNameEn),
		UpdatedAtString:                 updatedAtString,
		ExistingControls:                existingControls,
		ProposedControls:                nil,
		RiskStatusTitleEn:               riskDTO.RiskStatus.TitleEn,
		RiskStatusTitleAr:               riskDTO.RiskStatus.TitleAr,
		RiskStatusColorCode:             riskDTO.RiskStatus.ColorCode,
		RiskResponseStrategyTitleEn:     riskResponseStrategyTitleEn,
		RiskResponseStrategyTitleAr:     riskResponseStrategyTitleAr,
		ActionPlans:                     actionPlans,
		KeyRiskIdentifiers:              kris,
		Controls:                        controls,
		Assets:                          assets,
		RootCauseAnalyses:               rootCauseAnalyses,
	}

	section := riskDTOs.RiskRegisterExportSectionDTO{
		Type:        "identification",
		TitleEn:     "Risk Register",
		Description: "",
		Risks:       []riskDTOs.RiskRegisterExportRiskDTO{risk},
	}

	viewModel := riskDTOs.RiskRegisterExportPDFDTO{
		Title:                    riskDTO.Title,
		CreatedAt:                createdAt,
		GeneratedAt:              generatedAt,
		ReferenceNumber:          riskDTO.ReferenceNumber,
		Description:              riskDTO.Description,
		CreatedUserFirstNameEn:   preparedFirst,
		CreatedUserLastNameEn:    preparedLast,
		OrganizationUnitTitlesEn: orgUnit,
		StartDate:                reportedAt,
		EndDate:                  reportedAt,
		RiskStatusTitleEn:        riskDTO.RiskStatus.TitleEn,
		RiskStatusTitleAr:        riskDTO.RiskStatus.TitleAr,
		RiskStatusColorCode:      riskDTO.RiskStatus.ColorCode,
		RiskCategoryTitleEn:      riskDTO.RiskCategory.TitleEn,
		RiskCategoryTitleAr:      riskDTO.RiskCategory.TitleAr,
		// RiskSourceTitleEn:                     riskSourceTitleEn,
		// RiskSourceTitleAr:                     riskSourceTitleAr,
		ResponsibleUnitTitleEn:                orgUnit,
		ResponsibleUnitTitleAr:                orgUnitAr,
		ResponsibleUserFirstNameEn:            responsibleUserFirstNameEn,
		ResponsibleUserLastNameEn:             responsibleUserLastNameEn,
		ResponsibleUserInitial:                responsibleUserInitial,
		RiskCode:                              riskDTO.RiskCode,
		DivisionTitleEn:                       divisionTitleEn,
		DivisionTitleAr:                       divisionTitleAr,
		SectionTitleEn:                        sectionTitleEn,
		SectionTitleAr:                        sectionTitleAr,
		RegistrationDate:                      riskDTO.RegistrationDateString,
		IsInternal:                            riskDTO.IsInternal,
		IsExternal:                            riskDTO.IsExternal,
		IsStrategic:                           riskDTO.IsStrategic,
		IsOperational:                         riskDTO.IsOperational,
		ExpectedCause:                         riskDTO.ExpectedCause,
		Consequence:                           riskDTO.Consequence,
		ReviewFrequencyTitleEn:                reviewFrequencyTitleEn,
		ReviewFrequencyTitleAr:                reviewFrequencyTitleAr,
		NextReviewDateString:                  riskDTO.NextReviewDateString,
		LastReviewNote:                        riskDTO.LastReviewNote,
		NextReviewLevel:                       riskDTO.NextReviewLevel,
		SeverityLevelTitleEn:                  severityLevelTitleEn,
		SeverityLevelTitleAr:                  severityLevelTitleAr,
		SubmittedByUserFirstNameEn:            submittedByUserFirstNameEn,
		SubmittedByUserLastNameEn:             submittedByUserLastNameEn,
		SubmittedAtString:                     riskDTO.SubmittedAtString,
		UpdatedByUserFirstNameEn:              updatedByUserFirstNameEn,
		UpdatedByUserLastNameEn:               updatedByUserLastNameEn,
		UpdatedAtString:                       updatedAtString,
		RiskResponseStrategyTitleEn:           riskResponseStrategyTitleEn,
		RiskResponseStrategyTitleAr:           riskResponseStrategyTitleAr,
		RiskResponseStrategyDescription:       riskResponseStrategyDescription,
		RiskResponseStrategyTreatmentGuidance: riskResponseStrategyTreatmentGuidance,
		CurrentRiskRatingTitleEn:              currentRiskRatingTitleEn,
		CurrentRiskRatingTitleAr:              currentRiskRatingTitleAr,
		CurrentRiskRatingColorCode:            currentRiskRatingColorCode,
		CurrentRiskRatingValue:                currentRiskRatingValue,
		CurrentRiskPointerPosition:            currentRiskPointerPosition,
		CurrentLikelihoodTitleEn:              currentLikelihoodTitleEn,
		CurrentLikelihoodTitleAr:              currentLikelihoodTitleAr,
		CurrentLikelihoodScore:                currentLikelihoodScore,
		CurrentLikelihoodDescriptionEn:        currentLikelihoodDescriptionEn,
		CurrentImpactTitleEn:                  currentImpactTitleEn,
		CurrentImpactTitleAr:                  currentImpactTitleAr,
		CurrentImpactScore:                    currentImpactScore,
		CurrentImpactDescriptionEn:            currentImpactDescriptionEn,
		InherentRiskRatingTitleEn:             inherentRiskRatingTitleEn,
		InherentRiskRatingTitleAr:             inherentRiskRatingTitleAr,
		InherentRiskRatingColorCode:           inherentRiskRatingColorCode,
		InherentRiskRatingValue:               inherentRiskRatingValue,
		InherentRiskPointerPosition:           inherentRiskPointerPosition,
		InherentLikelihoodTitleEn:             inherentLikelihoodTitleEn,
		InherentLikelihoodTitleAr:             inherentLikelihoodTitleAr,
		InherentLikelihoodScore:               inherentLikelihoodScore,
		InherentLikelihoodDescriptionEn:       inherentLikelihoodDescriptionEn,
		InherentImpactTitleEn:                 inherentImpactTitleEn,
		InherentImpactTitleAr:                 inherentImpactTitleAr,
		InherentImpactScore:                   inherentImpactScore,
		InherentImpactDescriptionEn:           inherentImpactDescriptionEn,
		InherentImpactCategoryTitleEn:         inherentImpactCategoryTitleEn,
		InherentImpactCategoryTitleAr:         inherentImpactCategoryTitleAr,
		InherentImpactGuideline:               inherentImpactGuideline,
		ResidualRiskRatingTitleEn:             residualRiskRatingTitleEn,
		ResidualRiskRatingTitleAr:             residualRiskRatingTitleAr,
		ResidualRiskRatingColorCode:           residualRiskRatingColorCode,
		ResidualRiskRatingValue:               residualRiskRatingValue,
		ResidualRiskPointerPosition:           residualRiskPointerPosition,
		ResidualLikelihoodTitleEn:             residualLikelihoodTitleEn,
		ResidualLikelihoodTitleAr:             residualLikelihoodTitleAr,
		ResidualLikelihoodScore:               residualLikelihoodScore,
		ResidualLikelihoodDescriptionEn:       residualLikelihoodDescriptionEn,
		ResidualImpactTitleEn:                 residualImpactTitleEn,
		ResidualImpactTitleAr:                 residualImpactTitleAr,
		ResidualImpactScore:                   residualImpactScore,
		ResidualImpactDescriptionEn:           residualImpactDescriptionEn,
		AverageControlEffectivenessScore:      averageControlEffectivenessScore,
		ControlEffectivenessRatings:           controlEffectivenessRatings,
		ActionPlans:                           actionPlans,
		Controls:                              controls,
		KeyRiskIdentifiers:                    kris,
		Assets:                                assets,
		Documents:                             documents,
		Processes:                             processes,
		KPIs:                                  kpis,
		StrategicObjectives:                   strategicObjectives,
		ComplianceRequirements:                complianceRequirements,
		Incidents:                             incidents,
		RootCauseAnalyses:                     rootCauseAnalyses,
		RiskAppetiteStatements:                riskAppetiteStatements,
		ResidualRiskAppetiteStatements:        residualRiskAppetiteStatements,
		ClosureJustification:                  closureJustification,
		LessonsLearned:                        lessonsLearned,
		WhatWorkedWell:                        whatWorkedWell,
		WhatNeedsImprovement:                  whatNeedsImprovement,
		KeyTakeaways:                          keyTakeaways,
		ControlCoveragePercent:                controlCoveragePercent,
		IneffectiveControlsCount:              ineffectiveControlsCount,
		ImplementedPercent:                    implementedPercent,
		EffectiveControlsPercent:              effectiveControlsPercent,
		EffectiveControlsCount:                effectiveControlsCount,
		ControlScoreAchieved:                  controlScoreAchieved,
		ControlScoreRemaining:                 controlScoreRemaining,
		Reviews:                               reviews,
		Sections:                              []riskDTOs.RiskRegisterExportSectionDTO{section},
	}

	// Body
	bodyTemplates := []string{"risk/templates/pdf/risk_register_body.html"}
	htmlContent, err := utility.GeneratePDFContent(viewModel, bodyTemplates, "", "")
	if err != nil {
		return "", fmt.Errorf("failed to generate HTML content: %w", err)
	}

	// Cover
	coverTemplates := []string{"risk/templates/pdf/risk_register_cover.html"}
	coverContent, err := utility.GeneratePDFContent(viewModel, coverTemplates, "", "")
	if err != nil {
		return "", fmt.Errorf("failed to generate cover content: %w", err)
	}

	// Header
	headerTemplates := []string{"risk/templates/pdf/risk_register_header.html"}
	headerContent, err := utility.GeneratePDFContent(viewModel, headerTemplates, "", "")
	if err != nil {
		return "", fmt.Errorf("failed to generate header content: %w", err)
	}

	// Footer
	footerTemplates := []string{"risk/templates/pdf/risk_register_footer.html"}
	footerContent, err := utility.GeneratePDFContent(viewModel, footerTemplates, "", "")
	if err != nil {
		return "", fmt.Errorf("failed to generate footer content: %w", err)
	}

	pdfFileName := fmt.Sprintf("RiskRegister_%d.pdf", id)
	if err := utility.GeneratePDFWithCoverHeaderFooter(coverContent, htmlContent, pdfFileName, headerContent, footerContent); err != nil {
		return "", fmt.Errorf("failed to generate PDF: %w", err)
	}

	return pdfFileName, nil
}

func (rrs RiskRegisterService) getRiskPointerPosition(riskRatingTitle *string) *string {
	if riskRatingTitle == nil {
		return nil
	}

	title := strings.ToLower(strings.TrimSpace(*riskRatingTitle))

	var position float64

	switch title {
	case "very low":
		position = 10
	case "low":
		position = 25
	case "medium":
		position = 50
	case "high":
		position = 75
	case "very high", "extreme":
		position = 90
	default:
		return nil
	}

	value := fmt.Sprintf("%.1f", position)

	return &value
}

// fetchActionPlansForExport fetches action plans for a risk register and maps them to the export DTO format
func fetchActionPlansForExport(riskID uint, riskReferenceNumber string, riskTitle string, repo riskRepositories.RiskReportRepositoryInterface) ([]riskDTOs.RiskReportActionPlanResponseDTO, error) {
	actionPlanDTOs, err := repo.GetActionPlansByRiskID(riskID)
	if err != nil {
		return nil, err
	}

	actionPlans := make([]riskDTOs.RiskReportActionPlanResponseDTO, 0, len(actionPlanDTOs))
	for _, ap := range actionPlanDTOs {
		// Fetch proposed controls for this action plan
		proposedControls, err := repo.GetActionPlanControlsByActionPlanID(ap.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch controls for action plan %d: %w", ap.ID, err)
		}

		var proposedControlsPtr *string
		if proposedControls != "" {
			proposedControlsPtr = &proposedControls
		}

		var severityLevelTitleEnPtr *string
		var severityLevelTitleArPtr *string
		if ap.SeverityLevelTitleEn != "" {
			severityLevelTitleEnPtr = &ap.SeverityLevelTitleEn
		}
		if ap.SeverityLevelTitleAr != "" {
			severityLevelTitleArPtr = &ap.SeverityLevelTitleAr
		}

		targetDatePtr := &ap.TargetDate

		actionPlan := riskDTOs.RiskReportActionPlanResponseDTO{
			ID:                         ap.ID,
			RiskReportIdentificationID: 0,
			RiskReferenceNumber:        riskReferenceNumber,
			RiskTitle:                  riskTitle,
			ActionPlanTitle:            ap.ActionPlanTitle,
			ActionPlanDescription:      ap.ActionPlanDescription,
			SeverityLevelTitleEn:       severityLevelTitleEnPtr,
			SeverityLevelTitleAr:       severityLevelTitleArPtr,
			ResponsibleUserNameEn:      ap.ResponsibleUserNameEn,
			ResponsibleUserNameAr:      ap.ResponsibleUserNameAr,
			TargetDate:                 targetDatePtr,
			ActionPlanStatusTitleEn:    ap.ActionPlanStatusTitleEn,
			ActionPlanStatusTitleAr:    ap.ActionPlanStatusTitleAr,
			ActionPlanStatusColorCode:  ap.ActionPlanStatusColorCode,
			ProposedControls:           proposedControlsPtr,
		}
		actionPlans = append(actionPlans, actionPlan)
	}

	return actionPlans, nil
}

// fetchRootCauseAnalysesForExport fetches root cause analyses for a risk register
func fetchRootCauseAnalysesForExport(riskID uint, riskReportRepo riskRepositories.RiskReportRepositoryInterface) ([]riskDTOs.RiskRegisterExportRootCauseAnalysisDTO, error) {
	rootCauses, err := riskReportRepo.GetAllRiskRootCauseAnalyses(riskID)
	if err != nil {
		return nil, err
	}

	result := make([]riskDTOs.RiskRegisterExportRootCauseAnalysisDTO, 0, len(rootCauses))
	for _, rc := range rootCauses {
		result = append(result, riskDTOs.RiskRegisterExportRootCauseAnalysisDTO{
			Why:      rc.Why,
			Question: rc.Question,
			Answer:   rc.Answer,
		})
	}

	return result, nil
}

func (s *RiskRegisterService) GenerateReferenceNumber() (string, error) {
	return s.repo.GenerateReferenceNumber()
}

func (s *RiskRegisterService) UpdateRisk(risk *models.RiskRegister, tx ...utility.DatabaseTransaction) error {
	return s.repo.Update(risk, tx...)
}
