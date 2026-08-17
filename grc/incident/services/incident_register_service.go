package services

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"jethur/config"
	assetServices "jethur/grc/asset/services"
	commonDTOs "jethur/grc/common/dtos"
	commonServices "jethur/grc/common/services"
	complianceServices "jethur/grc/compliance/services"
	incidentDTOs "jethur/grc/incident/dtos"
	incidentRepositories "jethur/grc/incident/repositories"
	logServices "jethur/grc/logs/services"
	masterRepositories "jethur/grc/masters/repositories"
	masterServices "jethur/grc/masters/services"
	organizationServices "jethur/grc/organization/services"
	riskServices "jethur/grc/risk/services"
	settingServices "jethur/grc/settings/services"

	"jethur/grc/models"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type IncidentRegisterServiceInterface interface {
	GetAll(params *incidentDTOs.IncidentRegisterListQueryParamsDTO) (*incidentDTOs.IncidentRegisterListResponseDTO, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*incidentDTOs.IncidentRegisterDetailDTO, error)
	GetByID(id uint) (*models.IncidentRegister, error)
	Create(currentUserID uint, incidentRegisterFormDTO *incidentDTOs.IncidentRegisterFormDTO) (uint, error)
	Update(id uint, currentUserID uint, incident *incidentDTOs.IncidentRegisterFormDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	ForceDelete(ids []uint, currentUserID uint) ([]uint, error)
	StatusUpdate(id uint, currentUserID uint, incidentStatusID uint, tx utility.DatabaseTransaction) error
	UpdateIncidentStatus(id uint, currentUserID uint, incidentStatusID uint) error
	AcknowledgeIncident(id uint, currentUserID uint) error
	GetDocumentPath(id uint, token string) (string, error)
	ExportIncidentRegister(data []*incidentDTOs.IncidentRegisterListDTO, currentUserID uint) (*os.File, error)
}

// Service and implementations
type IncidentRegisterService struct {
	repo                     incidentRepositories.IncidentRegisterRepositoryInterface
	activityLogSvc           logServices.ActivityLogServiceInterface
	incidentCategorySvc      IncidentCategoryServiceInterface
	incidentTypeSvc          IncidentTypeServiceInterface
	incidentStatusSvc        IncidentStatusServiceInterface
	incidentResponseTeamSvc  IncidentResponseTeamServiceInterface
	severityLevelSvc         masterServices.SeverityLevelServiceInterface
	controlSvc               organizationServices.ControlServiceInterface
	complianceRequirementSvc complianceServices.ComplianceRequirementServiceInterface
	riskRegisterSvc          riskServices.RiskRegisterServiceInterface
	processSvc               organizationServices.ProcessServiceInterface
	slaRegisterSvc           complianceServices.SLARegisterServiceInterface
	incidentTimelineSvc      IncidentTimelineServiceInterface
	userSvc                  organizationServices.UserServiceInterface
	userSettingSvc           organizationServices.UserSettingServiceInterface
	notificationSvc          commonServices.NotificationServiceInterface
	assetRegisterSvc         assetServices.AssetRegisterServiceInterface
	systemPreferenceSvc      settingServices.SystemPreferenceServiceInterface
}

func NewIncidentRegisterService(
	repo incidentRepositories.IncidentRegisterRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	incidentCategorySvc IncidentCategoryServiceInterface,
	incidentTypeSvc IncidentTypeServiceInterface,
	incidentStatusSvc IncidentStatusServiceInterface,
	incidentResponseTeamSvc IncidentResponseTeamServiceInterface,
	severityLevelSvc masterRepositories.SeverityLevelRepositoryInterface,
	controlSvc organizationServices.ControlServiceInterface,
	complianceRequirementSvc complianceServices.ComplianceRequirementServiceInterface,
	riskRegisterSvc riskServices.RiskRegisterServiceInterface,
	processSvc organizationServices.ProcessServiceInterface,
	slaRegisterSvc complianceServices.SLARegisterServiceInterface,
	incidentTimelineSvc IncidentTimelineServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	userSettingSvc organizationServices.UserSettingServiceInterface,
	notificationSvc commonServices.NotificationServiceInterface,
	assetRegisterSvc assetServices.AssetRegisterServiceInterface,
	systemPreferenceSvc settingServices.SystemPreferenceServiceInterface,
) IncidentRegisterServiceInterface {
	return &IncidentRegisterService{
		repo:                     repo,
		activityLogSvc:           activityLogSvc,
		incidentCategorySvc:      incidentCategorySvc,
		incidentTypeSvc:          incidentTypeSvc,
		incidentStatusSvc:        incidentStatusSvc,
		incidentResponseTeamSvc:  incidentResponseTeamSvc,
		severityLevelSvc:         severityLevelSvc,
		controlSvc:               controlSvc,
		complianceRequirementSvc: complianceRequirementSvc,
		riskRegisterSvc:          riskRegisterSvc,
		processSvc:               processSvc,
		slaRegisterSvc:           slaRegisterSvc,
		incidentTimelineSvc:      incidentTimelineSvc,
		userSvc:                  userSvc,
		userSettingSvc:           userSettingSvc,
		notificationSvc:          notificationSvc,
		assetRegisterSvc:         assetRegisterSvc,
		systemPreferenceSvc:      systemPreferenceSvc,
	}
}

func (irs *IncidentRegisterService) GetAll(params *incidentDTOs.IncidentRegisterListQueryParamsDTO) (*incidentDTOs.IncidentRegisterListResponseDTO, error) {
	return irs.repo.GetAll(params)
}

func (irs IncidentRegisterService) GetByID(id uint) (*models.IncidentRegister, error) {
	incident, err := irs.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return incident, nil
}

func (irs IncidentRegisterService) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*incidentDTOs.IncidentRegisterDetailDTO, error) {
	// Initialize your DTO
	var incidentDTO incidentDTOs.IncidentRegisterDetailDTO

	incident, err := irs.repo.GetFullDetailsByID(id, tx...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	// Map the result to your DTO
	err = mapstructure.Decode(incident, &incidentDTO)
	if err != nil {
		return nil, err
	}

	incidentDTO.CreatedAt = incident.CreatedAt
	incidentDTO.UpdatedAt = incident.UpdatedAt

	return &incidentDTO, nil
}

func (irs *IncidentRegisterService) Create(currentUserID uint, incidentRegisterFormDTO *incidentDTOs.IncidentRegisterFormDTO) (uint, error) {
	tx := irs.repo.BeginTransaction()
	if err := irs.validateIncidentRegisterFormDTO(incidentRegisterFormDTO); err != nil {
		return 0, err
	}

	if incidentRegisterFormDTO.IncidentResponseTeamID == nil {
		return 0, utility.NewValidationError(fmt.Errorf("incident response team is required"), "required", "incident_response_team_id")
	}

	// Derive the response team lead and member users automatically from the selected master response team.
	selectedResponseTeam, err := irs.incidentResponseTeamSvc.GetFullDetailsByID(*incidentRegisterFormDTO.IncidentResponseTeamID)
	if err != nil {
		return 0, err
	}

	var responseTeamLeadID *uint
	if selectedResponseTeam.TeamLeadUser != nil {
		leadID := uint(selectedResponseTeam.TeamLeadUser.ID)
		responseTeamLeadID = &leadID
	}

	var responseTeamMemberUserIDs []*uint
	for _, member := range selectedResponseTeam.IncidentResponseTeamMembers {
		userID := member.UserID
		responseTeamMemberUserIDs = append(responseTeamMemberUserIDs, &userID)
	}

	dates, err := irs.parseIncidentRegisterDates(incidentRegisterFormDTO, currentUserID)
	if err != nil {
		return 0, err
	}

	referenceNumber, err := irs.repo.GenerateReferenceNumber()
	if err != nil {
		return 0, err
	}

	// update incident status as new
	incidentStatus, err := irs.incidentStatusSvc.GetByType("new")
	if err != nil {
		return 0, err
	}
	var organizationUnitID uint
	unitID, err := irs.userSvc.GetUserPrimaryUnitID(incidentRegisterFormDTO.ResponsibleUserID)
	if err != nil {
		return 0, utility.NewValidationError(fmt.Errorf("organization unit not set for owner"), "not_exist", "responsible unit")
	}
	organizationUnitID = unitID
	incidentRegister := &models.IncidentRegister{
		ReferenceNumber:                     referenceNumber,
		IncidentCategoryID:                  incidentRegisterFormDTO.IncidentCategoryID,
		IncidentTypeID:                      incidentRegisterFormDTO.IncidentTypeID,
		SeverityLevelID:                     *incidentRegisterFormDTO.SeverityLevelID,
		IncidentStatusID:                    &incidentStatus.ID,
		BusinessOwnerID:                     incidentRegisterFormDTO.BusinessOwnerID,
		IncidentResponseTeamID:              incidentRegisterFormDTO.IncidentResponseTeamID,
		ResponseTeamLeadID:                  responseTeamLeadID,
		IncidentScopeID:                     incidentRegisterFormDTO.IncidentScopeID,
		IncidentSensitivityClassificationID: incidentRegisterFormDTO.IncidentSensitivityClassificationID,
		ConfidentialitySeverityLevelID:      incidentRegisterFormDTO.ConfidentialitySeverityLevelID,
		IntegritySeverityLevelID:            incidentRegisterFormDTO.IntegritySeverityLevelID,
		AvailabilitySeverityLevelID:         incidentRegisterFormDTO.AvailabilitySeverityLevelID,
		ResponsibleUserID:                   incidentRegisterFormDTO.ResponsibleUserID,
		ResponsibleUnitID:                   organizationUnitID,
		Title:                               incidentRegisterFormDTO.Title,
		Description:                         incidentRegisterFormDTO.Description,
		IncidentAt:                          dates.IncidentAt,
		Location:                            incidentRegisterFormDTO.Location,
		ReportedBy:                          incidentRegisterFormDTO.ReportedBy,
		ReportedAt:                          dates.ReportedAt,
		QuickAction:                         incidentRegisterFormDTO.QuickAction,
		CreatedBYID:                         currentUserID,
	}

	id, err := irs.repo.Create(incidentRegister, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	// Add SLA Register
	var slaRegisters []*models.IncidentSLA
	for _, slaRegisterID := range incidentRegisterFormDTO.SLARegisterIDs {
		slaRegister := &models.IncidentSLA{
			IncidentRegisterID: id,
			SLARegisterID:      *slaRegisterID,
			CreatedByID:        currentUserID,
		}
		slaRegisters = append(slaRegisters, slaRegister)
	}

	if len(slaRegisters) > 0 {
		if err := irs.repo.CreateIncidentSLAs(slaRegisters, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	// Add response team members (derived automatically from the selected master response team)
	var responseTeamMembers []*models.IncidentRegisterResponseTeam
	for _, userID := range responseTeamMemberUserIDs {
		responseTeamMember := &models.IncidentRegisterResponseTeam{
			IncidentRegisterID: id,
			UserID:             *userID,
			CreatedByID:        currentUserID,
			CreatedAt:          time.Now(),
		}
		responseTeamMembers = append(responseTeamMembers, responseTeamMember)
	}

	if len(responseTeamMembers) > 0 {
		if err := irs.repo.CreateResponseTeamMembers(responseTeamMembers, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	// Add timeline
	// todo
	timelineTypes, err := irs.incidentTimelineSvc.GetAllTypes()
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	// Get timeline statuses
	completedStatus, err := irs.incidentTimelineSvc.GetByStatusType("completed")
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	pendingStatus, err := irs.incidentTimelineSvc.GetByStatusType("pending")
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	var timelineFormDTOs []*incidentDTOs.IncidentTimelineFormDTO

	for _, timelineType := range timelineTypes {
		var (
			statusID      uint
			performedByID *uint
			currentTime   *time.Time
		)

		if timelineType.Type == "reporting" {
			statusID = completedStatus.ID
			performedByID = &currentUserID
			now := time.Now()
			currentTime = &now
		} else {
			statusID = pendingStatus.ID
			performedByID = nil
			currentTime = nil
		}

		timelineFormDTO := &incidentDTOs.IncidentTimelineFormDTO{
			IncidentRegisterID:     id,
			IncidentTimelineTypeID: timelineType.ID,
			TimelineStatusID:       statusID,
			PerformedByID:          performedByID,
			PerformedAt:            currentTime,
		}
		timelineFormDTOs = append(timelineFormDTOs, timelineFormDTO)
	}
	if err := irs.incidentTimelineSvc.Create(timelineFormDTOs, tx); err != nil {
		tx.Rollback()
		return 0, err
	}

	var documents []*models.IncidentDocument

	// Iterate over each document in the array
	for _, doc := range incidentRegisterFormDTO.Documents {

		// Encrypt the filename
		dataToEncrypt := []byte(doc.FileName)
		encryptedFilename, err := utility.EncryptData(dataToEncrypt)
		if err != nil {
			tx.Rollback()
			return 0, err
		}

		// Construct the URL for the document
		url := utility.GenerateDocumentURL("incident", "incident-register", strconv.FormatUint(uint64(incidentRegister.ID), 10), "documents")

		// Create a new IncidentDocument instance
		document := &models.IncidentDocument{
			IncidentRegisterID: id,
			URL:                url,
			Name:               doc.FileName,
			Size:               doc.Size,
			Token:              encryptedFilename,
			OriginalName:       doc.OriginalName,
			Type:               doc.Type,
			CreatedByID:        currentUserID,
		}

		// Append the document to the slice
		documents = append(documents, document)
	}
	if len(documents) > 0 {
		// Save the IncidentDocument
		if err := irs.repo.CreateDocuments(documents, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
		for _, doc := range documents {
			storagePath := filepath.Join("incident", "incident-register", strconv.FormatUint(uint64(incidentRegister.ID), 10), "documents")
			if err := utility.StoreFile(doc.Name, storagePath); err != nil {
				tx.Rollback()
				return 0, err
			}
			// Convert DOCX to PDF if the document is a DOCX file
			// if filepath.Ext(doc.Name) == ".docx" {
			// 	if err := utility.ConvertDocxToPDF(storagePath, doc.Name, storagePath); err != nil {
			// 		tx.Rollback()
			// 		return 0, err
			// 	}
			// }
		}
	}

	tx.Commit()

	go func() {
		existingIncident, err := irs.GetFullDetailsByID(incidentRegister.ID)
		if err != nil {
			utility.LogError("Failed to fetch incident register", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		err = irs.activityLogSvc.LogActivity(
			models.CreateAction,
			models.IncidentRegister{}.TableName(),
			id,
			existingIncident,
			currentUserID,
			"incident", "incident_register", "register", nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.IncidentRegister{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}

		// Prepare notification and email
		urlParam := fmt.Sprintf(`{"id":%d}`, id)
		entity := "incident-register"
		action := "responsible-user-assign"
		dbNotificationDTO := &commonDTOs.DBNotificationDTO{
			UserIDs:              []uint{incidentRegisterFormDTO.ResponsibleUserID},
			ModuleIdentifier:     "incident",
			NotificationType:     string(commonDTOs.AssignmentType),
			HeadingEn:            "Incident Register Responsibility Assigned",
			HeadingAr:            "تسجيل الحوادث المسؤولية المعينة",
			TitleEn:              incidentRegisterFormDTO.Title,
			TitleAr:              incidentRegisterFormDTO.Title,
			ReferenceNumber:      incidentRegister.ReferenceNumber,
			Entity:               &entity,
			Action:               &action,
			URLParam:             urlParam,
			IsSystemNotification: false,
			CreatedByID:          &currentUserID,
		}
		responsibleUser, err := irs.userSvc.GetByID(incidentRegisterFormDTO.ResponsibleUserID)
		if err != nil {
			utility.LogError("Failed to fetch responsible user", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}
		createdBy, err := irs.userSvc.GetByID(currentUserID)
		if err != nil {
			utility.LogError("Failed to fetch responsible user", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}
		emailTo := []string{responsibleUser.Email}
		subject := fmt.Sprintf("Incident Register/%s: Responsibility Assigned", incidentRegister.Title)
		appURL := utility.GetAppURL()
		emailData := map[string]interface{}{
			"FirstNameEn":             responsibleUser.FirstNameEn,
			"LastNameEn":              responsibleUser.LastNameEn,
			"IncidentReferenceNumber": incidentRegister.ReferenceNumber,
			"IncidentTitle":           incidentRegister.Title,
			"IncidentDescription":     incidentRegister.Description,
			"CreatedByFirstNameEn":    createdBy.FirstNameEn,
			"CreatedByLastNameEn":     createdBy.LastNameEn,
			"AssignedOn":              time.Now().Format("2006-01-02"),
			"RedirectURL":             fmt.Sprintf("%s/incident/register/%d/overview", appURL, id),
		}
		emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
			EmailTo:              emailTo,
			Subject:              subject,
			TemplateRelativePath: "incident/templates/email/incident_register_create_responsible_user_mail.html",
			EmailData:            emailData,
		}
		err = irs.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "incident_register")
		if err != nil {
			utility.LogError("Failed to send notifications", map[string]interface{}{
				"error": err.Error(),
			})
		}
	}()

	return id, nil
}

func (irs *IncidentRegisterService) Update(id uint, currentUserID uint, incidentRegisterFormDTO *incidentDTOs.IncidentRegisterFormDTO) error {
	tx := irs.repo.BeginTransaction()

	existingIncidentRegister, err := irs.repo.GetByID(id)

	oldResponsibleUserID := existingIncidentRegister.ResponsibleUserID

	if err != nil {
		return utility.NewResourceNotFoundError(err, id)
	}

	oldIncidentRegister, err := irs.GetFullDetailsByID(id, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := irs.validateIncidentRegisterFormDTO(incidentRegisterFormDTO); err != nil {
		return err
	}

	dates, err := irs.parseIncidentRegisterDates(incidentRegisterFormDTO, currentUserID)
	if err != nil {
		return err
	}

	if existingIncidentRegister.ResponsibleUserID != incidentRegisterFormDTO.ResponsibleUserID {
		organizationUnitID, err := irs.userSvc.GetUserPrimaryUnitID(incidentRegisterFormDTO.ResponsibleUserID)
		if err != nil {
			return utility.NewValidationError(fmt.Errorf("organization unit not set for owner"), "not_exist", "control owner unit")
		}
		existingIncidentRegister.ResponsibleUnitID = organizationUnitID
	}
	existingIncidentRegister.IncidentCategoryID = incidentRegisterFormDTO.IncidentCategoryID
	existingIncidentRegister.IncidentTypeID = incidentRegisterFormDTO.IncidentTypeID
	existingIncidentRegister.SeverityLevelID = *incidentRegisterFormDTO.SeverityLevelID
	existingIncidentRegister.BusinessOwnerID = incidentRegisterFormDTO.BusinessOwnerID
	// IncidentResponseTeamID is intentionally not updated here; it is locked after creation.
	existingIncidentRegister.ResponseTeamLeadID = incidentRegisterFormDTO.ResponseTeamLeadID
	existingIncidentRegister.IncidentScopeID = incidentRegisterFormDTO.IncidentScopeID
	existingIncidentRegister.IncidentSensitivityClassificationID = incidentRegisterFormDTO.IncidentSensitivityClassificationID
	existingIncidentRegister.ConfidentialitySeverityLevelID = incidentRegisterFormDTO.ConfidentialitySeverityLevelID
	existingIncidentRegister.IntegritySeverityLevelID = incidentRegisterFormDTO.IntegritySeverityLevelID
	existingIncidentRegister.AvailabilitySeverityLevelID = incidentRegisterFormDTO.AvailabilitySeverityLevelID
	existingIncidentRegister.ResponsibleUserID = incidentRegisterFormDTO.ResponsibleUserID
	existingIncidentRegister.Title = incidentRegisterFormDTO.Title
	existingIncidentRegister.Description = incidentRegisterFormDTO.Description
	existingIncidentRegister.IncidentAt = dates.IncidentAt
	existingIncidentRegister.Location = incidentRegisterFormDTO.Location
	existingIncidentRegister.ReportedBy = incidentRegisterFormDTO.ReportedBy
	existingIncidentRegister.ReportedAt = dates.ReportedAt
	existingIncidentRegister.QuickAction = incidentRegisterFormDTO.QuickAction
	existingIncidentRegister.UpdatedBYID = &currentUserID

	err = irs.repo.Update(existingIncidentRegister, tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	// Collect existing control IDs
	existingIncident, err := irs.GetFullDetailsByID(id, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	existingSLARegisterIDs := make(map[uint]bool)
	for _, sla := range existingIncident.SLARegisters {
		existingSLARegisterIDs[sla.SLARegisterID] = true
	}

	// Collect new slaess IDs
	newSLARegisterIDs := make(map[uint]bool)
	for _, sla := range incidentRegisterFormDTO.SLARegisterIDs {
		newSLARegisterIDs[*sla] = true
	}

	// Add new SLA
	var slaRegistersToAdd []*models.IncidentSLA
	for _, newSLARegisterID := range incidentRegisterFormDTO.SLARegisterIDs {
		if _, exists := existingSLARegisterIDs[*newSLARegisterID]; !exists {
			// Create a new risk slaRegister record
			slaRegister := &models.IncidentSLA{
				IncidentRegisterID: id,
				SLARegisterID:      *newSLARegisterID,
				CreatedByID:        currentUserID,
			}
			slaRegistersToAdd = append(slaRegistersToAdd, slaRegister)
		}
	}

	// Save the new slaRegisters
	if len(slaRegistersToAdd) > 0 {
		if err := irs.repo.CreateIncidentSLAs(slaRegistersToAdd, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	// Remove slaRegisters not in the new list
	var deleteSLARegisterIDs []*uint

	for _, existingSLARegister := range existingIncident.SLARegisters {
		if _, exists := newSLARegisterIDs[existingSLARegister.SLARegisterID]; !exists {
			// Append slaRegister to delete list
			deleteSLARegisterIDs = append(deleteSLARegisterIDs, existingSLARegister.ID)
		}
	}

	// Delete slaRegisters not in the new list
	if len(deleteSLARegisterIDs) > 0 {
		if err := irs.repo.DeleteIncidentSLAs(deleteSLARegisterIDs, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	// Sync response team members
	existingResponseTeamMemberUserIDs := make(map[uint]bool)
	for _, responseTeamMember := range existingIncident.ResponseTeamMembers {
		existingResponseTeamMemberUserIDs[responseTeamMember.UserID] = true
	}

	newResponseTeamMemberUserIDs := make(map[uint]bool)
	for _, userID := range incidentRegisterFormDTO.ResponseTeamMemberUserIDs {
		newResponseTeamMemberUserIDs[*userID] = true
	}

	var responseTeamMembersToAdd []*models.IncidentRegisterResponseTeam
	for _, newUserID := range incidentRegisterFormDTO.ResponseTeamMemberUserIDs {
		if _, exists := existingResponseTeamMemberUserIDs[*newUserID]; !exists {
			responseTeamMember := &models.IncidentRegisterResponseTeam{
				IncidentRegisterID: id,
				UserID:             *newUserID,
				CreatedByID:        currentUserID,
				CreatedAt:          time.Now(),
			}
			responseTeamMembersToAdd = append(responseTeamMembersToAdd, responseTeamMember)
		}
	}

	if len(responseTeamMembersToAdd) > 0 {
		if err := irs.repo.CreateResponseTeamMembers(responseTeamMembersToAdd, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	var deleteResponseTeamMemberIDs []*uint
	for _, existingResponseTeamMember := range existingIncident.ResponseTeamMembers {
		if _, exists := newResponseTeamMemberUserIDs[existingResponseTeamMember.UserID]; !exists {
			deleteResponseTeamMemberIDs = append(deleteResponseTeamMemberIDs, &existingResponseTeamMember.ID)
		}
	}

	if len(deleteResponseTeamMemberIDs) > 0 {
		if err := irs.repo.DeleteResponseTeamMembers(deleteResponseTeamMemberIDs, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	// Update timeline
	// todo
	timelineType, err := irs.incidentTimelineSvc.GetByType("reporting")
	if err != nil {
		tx.Rollback()
		return err
	}
	timelineStatus, err := irs.incidentTimelineSvc.GetByStatusType("completed")
	if err != nil {
		tx.Rollback()
		return err
	}
	currentTime := time.Now()
	timelineFormDTO := &incidentDTOs.IncidentTimelineFormDTO{
		PerformedByID:    &currentUserID,
		TimelineStatusID: timelineStatus.ID,
		PerformedAt:      &currentTime,
	}
	if err := irs.incidentTimelineSvc.Update(id, timelineFormDTO, timelineType.ID, tx); err != nil {
		tx.Rollback()
		return err
	}
	// documents update
	var documents []*models.IncidentDocument
	var deleteDocumentIds []*uint
	var documentsToDelete []*models.IncidentDocument

	// Iterate over each document in the array
	for _, doc := range incidentRegisterFormDTO.Documents {
		if doc.IsNewFile != nil && *doc.IsNewFile {
			dataToEncrypt := []byte(doc.FileName)
			encryptedFilename, err := utility.EncryptData(dataToEncrypt)
			if err != nil {
				tx.Rollback()
				return err
			}
			url := utility.GenerateDocumentURL("incident", "incident-register", strconv.FormatUint(uint64(id), 10), "documents")

			// Create a new IncidentDocument instance
			document := &models.IncidentDocument{
				IncidentRegisterID: id,
				URL:                url,
				Name:               doc.FileName,
				Size:               doc.Size,
				Token:              encryptedFilename,
				OriginalName:       doc.OriginalName,
				Type:               doc.Type,
				CreatedByID:        currentUserID,
			}
			documents = append(documents, document)

		} else if doc.IsDeletedFile != nil && *doc.IsDeletedFile {
			// Handle deleted file case
			if doc.ID != nil {
				existingDocument, err := irs.repo.GetDocumentByID(*doc.ID)
				if err != nil {
					return err
				}
				if existingDocument == nil {
					tx.Rollback()
					return fmt.Errorf("existing document not found for ID: %d", *doc.ID)
				}
				deleteDocumentIds = append(deleteDocumentIds, &existingDocument.ID)
				documentsToDelete = append(documentsToDelete, existingDocument)
			}
		}
	}
	if len(documents) > 0 {
		if err := irs.repo.CreateDocuments(documents, tx); err != nil {
			tx.Rollback()
			return err
		}
	}
	if len(deleteDocumentIds) > 0 {
		if err := irs.repo.DeleteDocuments(deleteDocumentIds, tx); err != nil {
			tx.Rollback()
			return err
		}
	}
	if len(documents) > 0 {
		for _, doc := range documents {
			storagePath := filepath.Join("incident", "incident-register", strconv.FormatUint(uint64(existingIncidentRegister.ID), 10), "documents")
			if err := utility.StoreFile(doc.Name, storagePath); err != nil {
				tx.Rollback()
				return err
			}
		}
	}
	if len(documentsToDelete) > 0 {
		for _, document := range documentsToDelete {
			storagePath := filepath.Join("incident", "incident-register", strconv.FormatUint(uint64(existingIncidentRegister.ID), 10), "documents")
			if err := utility.RemoveFile(document.Name, storagePath); err != nil {
				tx.Rollback()
				return err
			}
		}
	}
	tx.Commit()

	go func() {
		existingIncident, err := irs.GetFullDetailsByID(id)
		if err != nil {
			return
		}
		slaRegistersAdded := []uint{}
		slaRegistersRemoved := []uint{}
		for id := range newSLARegisterIDs {
			if !existingSLARegisterIDs[id] {
				slaRegistersAdded = append(slaRegistersAdded, id)
			}
		}
		for id := range existingSLARegisterIDs {
			if !newSLARegisterIDs[id] {
				slaRegistersRemoved = append(slaRegistersRemoved, id)
			}
		}

		// Prepare log data
		changedData := make(map[string]interface{})
		associatedItems := make(map[string]interface{})

		if len(associatedItems) > 0 {
			changedData["associated_items"] = associatedItems
		}
		// log changes
		oldData := map[string]interface{}{
			"incident_category": oldIncidentRegister.IncidentCategory.TitleEn,
			"incident_type":     oldIncidentRegister.IncidentType.TitleEn,
			"severity_level":    oldIncidentRegister.SeverityLevel.TitleEn,
			"responsible_user":  fmt.Sprintf("%s %s", oldIncidentRegister.ResponsibleUser.FirstNameEN, oldIncidentRegister.ResponsibleUser.LastNameEN),
			"title":             oldIncidentRegister.Title,
			"description":       oldIncidentRegister.Description,
			"incident_at":       oldIncidentRegister.IncidentAt,
			"location":          oldIncidentRegister.Location,
			"reported_by":       fmt.Sprintf("%s %s", oldIncidentRegister.ReportedByUser.FirstNameEN, oldIncidentRegister.ReportedByUser.LastNameEN),
			"quick_action":      oldIncidentRegister.QuickAction,
		}

		newData := map[string]interface{}{
			"incident_category": existingIncident.IncidentCategory.TitleEn,
			"incident_type":     existingIncident.IncidentType.TitleEn,
			"severity_level":    existingIncident.SeverityLevel.TitleEn,
			"responsible_user":  fmt.Sprintf("%s %s", existingIncident.ResponsibleUser.FirstNameEN, existingIncident.ResponsibleUser.LastNameEN),
			"title":             existingIncident.Title,
			"description":       existingIncident.Description,
			"incident_at":       existingIncident.IncidentAt,
			"location":          existingIncident.Location,
			"reported_by":       fmt.Sprintf("%s %s", existingIncident.ReportedByUser.FirstNameEN, existingIncident.ReportedByUser.LastNameEN),
			"quick_action":      existingIncident.QuickAction,
		}
		// Compare other fields and add them to changedData
		fieldChanges := utility.CompareFieldChanges(oldData, newData)
		for k, v := range fieldChanges {
			changedData[k] = v
		}
		// Log activity if there are any changes
		if len(changedData) > 0 {
			err := irs.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.IncidentRegister{}.TableName(),
				id,
				changedData,
				currentUserID,
				"incident", "incident_register", "register", nil,
			)
			if err != nil {
				fields := map[string]interface{}{
					"action":    models.UpdateAction,
					"tableName": models.IncidentRegister{}.TableName(),
					"recordID":  id,
					"userID":    currentUserID,
					"error":     err.Error(),
				}

				utility.LogError("Failed to log activity", fields)
			}
		}

		// Send notifications to old and new responsible users if the responsible user has changed
		if oldResponsibleUserID != existingIncidentRegister.ResponsibleUserID {
			// URL param for notifications
			urlParam := fmt.Sprintf(`{"id":%d}`, id)
			entity := "incident-register"

			// Old responsible user notification (Removal)
			oldResponsibleUserAction := "responsible-user-remove"
			oldResponsibleUserNotificationDTO := &commonDTOs.DBNotificationDTO{
				UserIDs:              []uint{oldResponsibleUserID},
				ModuleIdentifier:     "incident",
				NotificationType:     string(commonDTOs.InformationType),
				HeadingEn:            "Incident Register Responsibility Removed",
				HeadingAr:            "تمت إزالة مسؤولية تسجيل الحوادث",
				TitleEn:              existingIncidentRegister.Title,
				TitleAr:              existingIncidentRegister.Title,
				ReferenceNumber:      existingIncidentRegister.ReferenceNumber,
				Entity:               &entity,
				Action:               &oldResponsibleUserAction,
				URLParam:             urlParam,
				IsSystemNotification: false,
				CreatedByID:          &currentUserID,
			}

			// New responsible user notification (Assignment)
			newResponsibleUserAction := "responsible-user-assign"
			newResponsibleUserNotificationDTO := &commonDTOs.DBNotificationDTO{
				UserIDs:              []uint{incidentRegisterFormDTO.ResponsibleUserID},
				ModuleIdentifier:     "incident",
				NotificationType:     string(commonDTOs.AssignmentType),
				HeadingEn:            "Incident Register Responsibility Assigned",
				HeadingAr:            "تسجيل الحوادث المسؤولية المعينة",
				TitleEn:              existingIncidentRegister.Title,
				TitleAr:              existingIncidentRegister.Title,
				ReferenceNumber:      existingIncidentRegister.ReferenceNumber,
				Entity:               &entity,
				Action:               &newResponsibleUserAction,
				URLParam:             urlParam,
				IsSystemNotification: false,
				CreatedByID:          &currentUserID,
			}

			// Fetch the responsible users' data
			oldResponsibleUser, err := irs.userSvc.GetByID(oldResponsibleUserID)
			if err != nil {
				utility.LogError("Failed to fetch old responsible user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}

			newResponsibleUser, err := irs.userSvc.GetByID(incidentRegisterFormDTO.ResponsibleUserID)
			if err != nil {
				utility.LogError("Failed to fetch new responsible user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}

			// Fetch details of the user making the assignment
			updatedByUser, err := irs.userSvc.GetByID(currentUserID)
			if err != nil {
				utility.LogError("Failed to fetch assigning user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}

			appURL := utility.GetAppURL()

			// Email notification for the old responsible user
			oldResponsibleUserEmailTo := []string{oldResponsibleUser.Email}
			oldResponsibleUserSubject := fmt.Sprintf("Incident Register/%s: Responsibility Removed", existingIncidentRegister.Title)
			oldResponsibleUserEmailData := map[string]interface{}{
				"FirstNameEn":             oldResponsibleUser.FirstNameEn,
				"LastNameEn":              oldResponsibleUser.LastNameEn,
				"IncidentReferenceNumber": existingIncidentRegister.ReferenceNumber,
				"IncidentTitle":           existingIncidentRegister.Title,
				"IncidentDescription":     existingIncidentRegister.Description,
				"UpdatedByFirstNameEn":    updatedByUser.FirstNameEn,
				"UpdatedByLastNameEn":     updatedByUser.LastNameEn,
				"RedirectURL":             fmt.Sprintf("%s/incident/register/%d/overview", appURL, id),
			}

			oldResponsibleUserEmailNotificationDTO := &commonDTOs.EmailNotificationDTO{
				EmailTo:              oldResponsibleUserEmailTo,
				Subject:              oldResponsibleUserSubject,
				TemplateRelativePath: "incident/templates/email/incident_register_remove_responsible_user_mail.html",
				EmailData:            oldResponsibleUserEmailData,
			}

			// Email notification for the new responsible user
			newResponsibleUserEmailTo := []string{newResponsibleUser.Email}
			newResponsibleUserSubject := fmt.Sprintf("Incident Register/%s: Responsibility Assigned", existingIncidentRegister.Title)
			newResponsibleUserEmailData := map[string]interface{}{
				"FirstNameEn":             newResponsibleUser.FirstNameEn,
				"LastNameEn":              newResponsibleUser.LastNameEn,
				"IncidentReferenceNumber": existingIncidentRegister.ReferenceNumber,
				"IncidentTitle":           existingIncidentRegister.Title,
				"IncidentDescription":     existingIncidentRegister.Description,
				"AssignedByFirstNameEn":   updatedByUser.FirstNameEn,
				"AssignedByLastNameEn":    updatedByUser.LastNameEn,
				"AssignedOn":              time.Now().Format("2006-01-02"),
				"RedirectURL":             fmt.Sprintf("%s/incident/register/%d/overview", appURL, id),
			}

			newResponsibleUserEmailNotificationDTO := &commonDTOs.EmailNotificationDTO{
				EmailTo:              newResponsibleUserEmailTo,
				Subject:              newResponsibleUserSubject,
				TemplateRelativePath: "incident/templates/email/incident_register_change_responsible_user_mail.html",
				EmailData:            newResponsibleUserEmailData,
			}

			// Send notifications and emails
			err = irs.notificationSvc.SendNotification(newResponsibleUserNotificationDTO, newResponsibleUserEmailNotificationDTO, "incident_register")
			if err != nil {
				utility.LogError("Failed to send notifications to new responsible user", map[string]interface{}{
					"error": err.Error(),
				})
			}

			err = irs.notificationSvc.SendNotification(oldResponsibleUserNotificationDTO, oldResponsibleUserEmailNotificationDTO, "incident_register")
			if err != nil {
				utility.LogError("Failed to send notifications to old responsible user", map[string]interface{}{
					"error": err.Error(),
				})
			}
		}
	}()

	return nil
}

func (irs *IncidentRegisterService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	tx := irs.repo.BeginTransaction()
	var failedIds []uint
	var successfulIds []uint
	var deletedIncidentRegister []map[string]interface{}

	var notifications []struct {
		OwnerID         *uint
		TitleEn         string
		TitleAr         string
		Description     string
		ReferenceNumber string
	}

	documentsToDelete := make(map[uint][]*incidentDTOs.IncidentDocumentDTO)

	for _, id := range ids {
		incident, err := irs.GetFullDetailsByID(id, tx)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		deletedIncidentRegister = append(deletedIncidentRegister, map[string]interface{}{
			"id":    id,
			"title": incident.Title,
		})

		// Prepare notification data before deletion
		notifications = append(notifications, struct {
			OwnerID         *uint
			TitleEn         string
			TitleAr         string
			Description     string
			ReferenceNumber string
		}{
			OwnerID:         &incident.ResponsibleUserID,
			TitleEn:         incident.Title,
			Description:     incident.Description,
			ReferenceNumber: incident.ReferenceNumber,
		})

		if len(incident.SLARegisters) > 0 {
			var deleteSLARegisterIDs []*uint
			for _, sla := range incident.SLARegisters {
				deleteSLARegisterIDs = append(deleteSLARegisterIDs, sla.ID)
			}
			if err := irs.repo.DeleteIncidentSLAs(deleteSLARegisterIDs, tx); err != nil {
				tx.Rollback()
				return nil, err
			}
		}

		// Delete response team members
		if len(incident.ResponseTeamMembers) > 0 {
			var deleteResponseTeamMemberIDs []*uint
			for _, responseTeamMember := range incident.ResponseTeamMembers {
				deleteResponseTeamMemberIDs = append(deleteResponseTeamMemberIDs, &responseTeamMember.ID)
			}
			if err := irs.repo.DeleteResponseTeamMembers(deleteResponseTeamMemberIDs, tx); err != nil {
				tx.Rollback()
				return nil, err
			}
		}

		// Delete timeline
		// todo
		if err := irs.incidentTimelineSvc.Delete(id, tx); err != nil {
			tx.Rollback()
			return nil, err
		}

		// Delete each incident document if documents exist
		if len(incident.Documents) > 0 {
			var deleteDocumentIds []*uint
			for _, doc := range incident.Documents {
				if doc != nil {
					documentsToDelete[id] = append(documentsToDelete[id], doc)
					deleteDocumentIds = append(deleteDocumentIds, doc.ID)
				}
			}
			if err := irs.repo.DeleteDocuments(deleteDocumentIds, tx); err != nil {
				tx.Rollback()
				return nil, err
			}

		}

		// Delete the incident
		if err := irs.repo.Delete(id, tx); err != nil {
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

	// Deleting files should not affect the main operation result
	if len(successfulIds) > 0 {
		for incidentID, documents := range documentsToDelete {
			for _, doc := range documents {
				fmt.Println("filename=", doc.Name)
				storagePath := filepath.Join("incident", "incident-register", strconv.FormatUint(uint64(incidentID), 10), "documents")
				if err := utility.RemoveFile(doc.Name, storagePath); err != nil {
					tx.Rollback()
					fmt.Printf("Failed to delete file at URL %s: %v\n", *doc.URL, err)
				}
			}
		}
	}

	tx.Commit()

	// Log bulk delete activities for successful deletions
	if len(successfulIds) > 0 {
		go func() {
			err := irs.activityLogSvc.LogBulkActivities(
				models.DeleteAction, // Make sure this action constant is defined appropriately
				models.IncidentRegister{}.TableName(),
				successfulIds,
				deletedIncidentRegister,
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
				entity := "incident-register"
				action := "delete"

				// Prepare DB notification
				dbNotificationDTO := &commonDTOs.DBNotificationDTO{
					UserIDs:              []uint{*notification.OwnerID},
					ModuleIdentifier:     "incident",
					NotificationType:     string(commonDTOs.DeletionType),
					HeadingEn:            "Incident Register Deleted",
					HeadingAr:            "تم حذف سجل الحوادث",
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
				ownerUser, err := irs.userSvc.GetByID(*notification.OwnerID)
				if err != nil {
					utility.LogError("Failed to fetch owner details for notification", map[string]interface{}{
						"error": err.Error(),
					})
					continue
				}

				// Prepare email notification
				ownerEmailTo := []string{ownerUser.Email}
				subject := fmt.Sprintf("Incident Register/%s: Deleted", notification.TitleEn)
				templateRelativePath := "incident/templates/email/incident_register_delete_mail.html"

				emailData := map[string]interface{}{
					"FirstNameEn":             ownerUser.FirstNameEn,
					"LastNameEn":              ownerUser.LastNameEn,
					"IncidentReferenceNumber": notification.ReferenceNumber,
					"IncidentTitle":           notification.TitleEn,
					"IncidentDescription":     notification.Description,
				}

				emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
					EmailTo:              ownerEmailTo,
					Subject:              subject,
					TemplateRelativePath: templateRelativePath,
					EmailData:            emailData,
				}

				// Send DB and email notifications
				err = irs.notificationSvc.SendNotification(
					dbNotificationDTO, emailNotificationDTO, "incident_register")
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

func (irs *IncidentRegisterService) StatusUpdate(id uint, currentUserID uint, incidentStatusID uint, tx utility.DatabaseTransaction) error {
	existingIncident, err := irs.GetByID(id)
	if err != nil {
		return err
	}

	existingIncident.IncidentStatusID = &incidentStatusID
	existingIncident.UpdatedBYID = &currentUserID

	err = irs.repo.Update(existingIncident, tx)
	if err != nil {
		return err
	}

	return nil
}

func (irs *IncidentRegisterService) UpdateIncidentStatus(id uint, currentUserID uint, incidentStatusID uint) error {
	newIncidentStatus, err := irs.incidentStatusSvc.GetByID(incidentStatusID)
	if err != nil {
		return utility.NewValidationError(fmt.Errorf("incident status does not exist"), "not_exist", "incident_status_id")
	}

	existingIncident, err := irs.GetByID(id)
	if err != nil {
		return err
	}

	if existingIncident.IncidentStatusID != nil && *existingIncident.IncidentStatusID == newIncidentStatus.ID {
		return utility.NewValidationError(fmt.Errorf("incident is already in %s status", newIncidentStatus.Type), "already_in_status", "incident_status_id")
	}

	tx := irs.repo.BeginTransaction()

	if err := irs.StatusUpdate(id, currentUserID, newIncidentStatus.ID, tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	go func() {
		changedData := map[string]interface{}{
			"incident_status_id":    newIncidentStatus.ID,
			"incident_status_type":  newIncidentStatus.Type,
			"incident_status_title": newIncidentStatus.TitleEn,
		}

		err := irs.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.IncidentRegister{}.TableName(),
			id,
			changedData,
			currentUserID,
			"incident", "incident_register", "incident_status",
			nil,
		)
		if err != nil {
			utility.LogError("Failed to log incident status activity", map[string]interface{}{
				"incident_id": id,
				"user_id":     currentUserID,
				"error":       err.Error(),
			})
		}
	}()

	return nil
}

func (irs *IncidentRegisterService) AcknowledgeIncident(id uint, currentUserID uint) error {
	existingIncident, err := irs.GetByID(id)
	if err != nil {
		return err
	}

	if existingIncident.IsAcknowledge {
		return utility.NewValidationError(fmt.Errorf("incident is already acknowledged"), "already_acknowledged", "is_acknowledge")
	}

	tx := irs.repo.BeginTransaction()

	now := time.Now()
	existingIncident.IsAcknowledge = true
	existingIncident.AcknowledgedAt = &now
	existingIncident.AcknowledgedByID = &currentUserID
	existingIncident.UpdatedBYID = &currentUserID

	if err := irs.repo.Update(existingIncident, tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	go func() {
		changedData := map[string]interface{}{
			"is_acknowledge":     true,
			"acknowledged_at":    now,
			"acknowledged_by_id": currentUserID,
		}

		err := irs.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.IncidentRegister{}.TableName(),
			id,
			changedData,
			currentUserID,
			"incident", "incident_register", "acknowledge",
			nil,
		)
		if err != nil {
			utility.LogError("Failed to log incident acknowledge activity", map[string]interface{}{
				"incident_id": id,
				"user_id":     currentUserID,
				"error":       err.Error(),
			})
		}
	}()

	return nil
}

func (irs *IncidentRegisterService) validateIncidentRegisterFormDTO(incidentRegisterFormDTO *incidentDTOs.IncidentRegisterFormDTO) error {
	if _, err := irs.incidentCategorySvc.GetByID(incidentRegisterFormDTO.IncidentCategoryID); err != nil {
		return utility.NewValidationError(fmt.Errorf("incident category does not exist"), "not_exist", "incident category")
	}

	if incidentRegisterFormDTO.IncidentTypeID != nil {
		if _, err := irs.incidentTypeSvc.GetByID(*incidentRegisterFormDTO.IncidentTypeID); err != nil {
			return utility.NewValidationError(fmt.Errorf("incident type does not exist"), "not_exist", "incident type")
		}
	}
	if incidentRegisterFormDTO.SeverityLevelID != nil {
		if _, err := irs.severityLevelSvc.GetByID(*incidentRegisterFormDTO.SeverityLevelID); err != nil {
			return utility.NewValidationError(fmt.Errorf("severity level does not exist"), "not_exist", "severity level")
		}
	}
	if _, err := irs.userSvc.GetByID(uint(incidentRegisterFormDTO.ResponsibleUserID)); err != nil {
		return utility.NewValidationError(fmt.Errorf("responsible user does not exist"), "not_exist", "responsible user")
	}
	if incidentRegisterFormDTO.BusinessOwnerID != nil {
		if _, err := irs.userSvc.GetByID(*incidentRegisterFormDTO.BusinessOwnerID); err != nil {
			return utility.NewValidationError(fmt.Errorf("business owner does not exist"), "not_exist", "business owner")
		}
	}
	if incidentRegisterFormDTO.ConfidentialitySeverityLevelID != nil {
		if _, err := irs.severityLevelSvc.GetByID(*incidentRegisterFormDTO.ConfidentialitySeverityLevelID); err != nil {
			return utility.NewValidationError(fmt.Errorf("confidentiality severity level does not exist"), "not_exist", "confidentiality severity level")
		}
	}
	if incidentRegisterFormDTO.IntegritySeverityLevelID != nil {
		if _, err := irs.severityLevelSvc.GetByID(*incidentRegisterFormDTO.IntegritySeverityLevelID); err != nil {
			return utility.NewValidationError(fmt.Errorf("integrity severity level does not exist"), "not_exist", "integrity severity level")
		}
	}
	if incidentRegisterFormDTO.AvailabilitySeverityLevelID != nil {
		if _, err := irs.severityLevelSvc.GetByID(*incidentRegisterFormDTO.AvailabilitySeverityLevelID); err != nil {
			return utility.NewValidationError(fmt.Errorf("availability severity level does not exist"), "not_exist", "availability severity level")
		}
	}

	// Controls exist
	controls := incidentRegisterFormDTO.ControlIds

	// Create a map to store control IDs
	controlIDs := make(map[uint]bool)

	// Check if each control exists
	for _, control := range controls {
		// Check if the control ID already exists in the map
		if controlIDs[*control] {
			// If the control ID already exists, return a validation error
			return utility.NewValidationError(fmt.Errorf("duplicate control ID: %d", control), "duplicate", "control")
		}

		// Check if the control exists in the database
		if _, err := irs.controlSvc.GetByID(*control); err != nil {
			return utility.NewValidationError(fmt.Errorf("control does not exist for ID: %d", control), "not_exist", "control")
		}

		// Add the control ID to the map
		controlIDs[*control] = true
	}
	slaRegisterIds := incidentRegisterFormDTO.SLARegisterIDs

	slaRegisterIDs := make(map[uint]bool)

	for _, slaRegisterId := range slaRegisterIds {
		if slaRegisterIDs[*slaRegisterId] {
			return utility.NewValidationError(fmt.Errorf("duplicate slaRegister ID: %d", slaRegisterId), "duplicate", "sla register")
		}

		if _, err := irs.slaRegisterSvc.GetByID(*slaRegisterId); err != nil {
			return utility.NewValidationError(fmt.Errorf("slaRegister does not exist for ID: %d", slaRegisterId), "not_exist", "sla register")
		}

		slaRegisterIDs[*slaRegisterId] = true
	}

	// Response team
	if incidentRegisterFormDTO.IncidentResponseTeamID != nil {
		if _, err := irs.incidentResponseTeamSvc.GetByID(*incidentRegisterFormDTO.IncidentResponseTeamID); err != nil {
			return utility.NewValidationError(fmt.Errorf("incident response team does not exist for ID: %d", *incidentRegisterFormDTO.IncidentResponseTeamID), "not_exist", "incident response team")
		}
	}

	// Response team lead
	if incidentRegisterFormDTO.ResponseTeamLeadID != nil {
		if _, err := irs.userSvc.GetByID(*incidentRegisterFormDTO.ResponseTeamLeadID); err != nil {
			return utility.NewValidationError(fmt.Errorf("response team lead does not exist for ID: %d", *incidentRegisterFormDTO.ResponseTeamLeadID), "not_exist", "response team lead")
		}
	}

	// Response team members
	responseTeamMemberUserIDs := make(map[uint]bool)
	for _, userID := range incidentRegisterFormDTO.ResponseTeamMemberUserIDs {
		if responseTeamMemberUserIDs[*userID] {
			return utility.NewValidationError(fmt.Errorf("duplicate response team member user ID: %d", *userID), "duplicate", "response team member")
		}

		if _, err := irs.userSvc.GetByID(*userID); err != nil {
			return utility.NewValidationError(fmt.Errorf("response team member does not exist for ID: %d", *userID), "not_exist", "response team member")
		}

		responseTeamMemberUserIDs[*userID] = true
	}

	return nil
}

func (irs *IncidentRegisterService) parseIncidentRegisterDates(dto *incidentDTOs.IncidentRegisterFormDTO, currentUserID uint) (*incidentDTOs.IncidentRegisterParsedDates, error) {
	parsed := &incidentDTOs.IncidentRegisterParsedDates{}

	userPreference, err := irs.userSvc.GetCurrentUserDetails(currentUserID)
	if err != nil {
		return nil, err
	}

	timezoneName := userPreference.UserSetting.Timezone.Name

	fmt.Println(timezoneName)

	if dto.ReportedAt != "" {
		parsed.ReportedAt, err = utility.ConvertToUTC(dto.ReportedAt, timezoneName)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid reported date: %w", err), "invalid_date", "reported_at")
		}
	}

	if dto.IncidentAt != "" {
		parsed.IncidentAt, err = utility.ConvertToUTC(dto.IncidentAt, timezoneName)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid incident date: %w", err), "invalid_date", "incident_at")
		}
	}

	return parsed, nil
}

// func (irs *IncidentRegisterService) parseIncidentRegisterDates(dto *incidentDTOs.IncidentRegisterFormDTO) (*incidentDTOs.IncidentRegisterParsedDates, error) {
// 	parsed := &incidentDTOs.IncidentRegisterParsedDates{}

// 	var err error
// 	if dto.ReportedAt != "" {
// 		parsed.ReportedAt, err = utility.ParseDateTime(dto.ReportedAt)
// 		if err != nil {
// 			return nil, utility.NewValidationError(fmt.Errorf("invalid reported date"), "invalid_date", "reported at")
// 		}
// 	}

// 	if dto.IncidentAt != "" {
// 		parsed.IncidentAt, err = utility.ParseDateTime(dto.IncidentAt)
// 		if err != nil {
// 			return nil, utility.NewValidationError(fmt.Errorf("invalid incident date"), "invalid_date", "incident at")
// 		}
// 	}

// 	return parsed, nil
// }

func (irs *IncidentRegisterService) validateIncidentDocumentFormDTO(documents []*incidentDTOs.IncidentDocumentFormDTO) error {
	// Check if documents array is empty
	if len(documents) == 0 {
		return utility.NewValidationError(fmt.Errorf("documents array is required"), "required", "documents")
	}

	return nil
}

func (irs *IncidentRegisterService) GetDocumentPath(id uint, token string) (string, error) {
	existingIncidentRegister, err := irs.GetFullDetailsByID(id)
	if err != nil {
		return "", err
	}

	if len(existingIncidentRegister.Documents) == 0 {
		return "", fmt.Errorf("no documents found for the incident register")
	}

	decryptedFilenameBytes, err := utility.DecryptData(token)
	if err != nil {
		return "", fmt.Errorf("resource not found")
	}

	decryptedFilename := string(decryptedFilenameBytes)
	if decryptedFilename == "" {
		return "", fmt.Errorf("resource not found")
	}

	storagePath := filepath.Join("incident", "incident-register", strconv.FormatUint(uint64(id), 10), "documents")

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

func (is *IncidentRegisterService) ExportIncidentRegister(data []*incidentDTOs.IncidentRegisterListDTO, currentUserID uint) (*os.File, error) {
	// Fetch user setting for language preference
	userSetting, err := is.userSettingSvc.GetUserSettingsByUserID(currentUserID)
	if err != nil {
		return nil, err
	}
	var headers []string
	// Define headers based on the user's language setting
	if userSetting.Language == "ar" {
		headers = []string{
			"SL No.",
			"Reference Number",
			"Title",
			"Incident AT",
			"Category",
			"Severity Level",
			"Location",
			"Reported By",
			"Responsible User",
			"Reported Date & Time",
			"Current Status",
			"Action Taken",
			"Investigation Date & Time",
			"Investigation Details",
		}
	} else {
		headers = []string{
			"SL No.",
			"Reference Number",
			"Title",
			"Incident AT",
			"Category",
			"Severity Level",
			"Location",
			"Reported By",
			"Responsible User",
			"Reported Date & Time",
			"Current Status",
			"Action Taken",
			"Investigation Date & Time",
			"Investigation Details",
		}
	}
	// Convert the incident register data to [][]interface{} based on language setting
	rows := make([][]interface{}, len(data))
	for i, incident := range data {
		if userSetting.Language == "ar" {
			// Populate row data in Arabic
			rows[i] = []interface{}{
				i + 1,
				incident.ReferenceNumber,
				incident.Title,
				incident.IncidentAt,
				incident.IncidentCategoryTitleAr,
				incident.SeverityLevelTitleAr,
				utility.GetStringValue(incident.Location),
				incident.ReportedUserFirstNameAr + " " + incident.ReportedUserLastNameAr,
				incident.ResponsibleUserFirstNameAr + " " + incident.ResponsibleUserLastNameAr,
				incident.ReportedAt,
				incident.IncidentStatusTitleAr,
				utility.GetStringValue(incident.QuickAction),
				utility.GetStringValue(incident.InvestigationInitiatedAt),
				utility.GetStringValue(incident.InvestigationDescription),
			}
		} else {
			// Populate row data in English
			rows[i] = []interface{}{
				i + 1,
				incident.ReferenceNumber,
				incident.Title,
				incident.IncidentAt,
				incident.IncidentCategoryTitleEn,
				incident.SeverityLevelTitleEn,
				utility.GetStringValue(incident.Location),
				incident.ReportedUserFirstNameEn + " " + incident.ReportedUserLastNameEn,
				incident.ResponsibleUserFirstNameEn + " " + incident.ResponsibleUserLastNameEn,
				incident.ReportedAt,
				incident.IncidentStatusTitleEn,
				utility.GetStringValue(incident.QuickAction),
				utility.GetStringValue(incident.InvestigationInitiatedAt),
				utility.GetStringValue(incident.InvestigationDescription),
			}
		}
	}

	// Generate the Excel file
	return utility.GenerateExcelFile(headers, rows)
}

func (is *IncidentRegisterService) ForceDelete(ids []uint, currentUserID uint) ([]uint, error) {
	tx := is.repo.BeginTransaction()
	var failedIds []uint
	var deletedData []map[string]interface{}

	for _, id := range ids {
		if err := is.repo.ForceDelete(id, tx); err != nil {
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
		err := is.activityLogSvc.LogBulkActivities(
			models.DeleteAction, // Make sure this action constant is defined appropriately
			models.IncidentRegister{}.TableName(),
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
