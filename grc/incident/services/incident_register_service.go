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
	complianceServices "jethur/grc/compliance/services"
	incidentDTOs "jethur/grc/incident/dtos"
	incidentRepositories "jethur/grc/incident/repositories"
	logServices "jethur/grc/logs/services"
	masterRepositories "jethur/grc/masters/repositories"
	masterServices "jethur/grc/masters/services"
	organizationServices "jethur/grc/organization/services"
	riskServices "jethur/grc/risk/services"

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
	StatusUpdate(id uint, currentUserID uint, incidentStatusID uint, tx utility.DatabaseTransaction) error
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
}

func NewIncidentRegisterService(
	repo incidentRepositories.IncidentRegisterRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	incidentCategorySvc IncidentCategoryServiceInterface,
	incidentTypeSvc IncidentTypeServiceInterface,
	incidentStatusSvc IncidentStatusServiceInterface,
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
) IncidentRegisterServiceInterface {
	return &IncidentRegisterService{
		repo:                     repo,
		activityLogSvc:           activityLogSvc,
		incidentCategorySvc:      incidentCategorySvc,
		incidentTypeSvc:          incidentTypeSvc,
		incidentStatusSvc:        incidentStatusSvc,
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

	dates, err := irs.parseIncidentRegisterDates(incidentRegisterFormDTO)
	if err != nil {
		return 0, err
	}

	referenceNumber, err := irs.repo.GenerateReferenceNumber()
	if err != nil {
		return 0, err
	}

	if err := irs.validateIncidentDocumentFormDTO(incidentRegisterFormDTO.Documents); err != nil {
		return 0, err
	}
	//update incident status as open
	incidentStatus, err := irs.incidentStatusSvc.GetByType("open")
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
		ReferenceNumber:    referenceNumber,
		IncidentCategoryID: incidentRegisterFormDTO.IncidentCategoryID,
		IncidentTypeID:     incidentRegisterFormDTO.IncidentTypeID,
		SeverityLevelID:    *incidentRegisterFormDTO.SeverityLevelID,
		IncidentStatusID:   &incidentStatus.ID,
		ResponsibleUserID:  incidentRegisterFormDTO.ResponsibleUserID,
		ResponsibleUnitID:  organizationUnitID,
		Title:              incidentRegisterFormDTO.Title,
		Description:        incidentRegisterFormDTO.Description,
		IncidentAt:         dates.IncidentAt,
		Location:           incidentRegisterFormDTO.Location,
		ReportedBy:         incidentRegisterFormDTO.ReportedBy,
		ReportedAt:         dates.ReportedAt,
		QuickAction:        incidentRegisterFormDTO.QuickAction,
		CreatedByID:        currentUserID,
	}

	id, err := irs.repo.Create(incidentRegister, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	// Save ComplianceRequirements
	var complianceRequirements []*models.IncidentComplianceRequirement
	for _, reqId := range incidentRegisterFormDTO.ComplianceRequirementIds {
		complianceRequirement, err := irs.complianceRequirementSvc.GetByComplianceRequirementID(*reqId)
		if err != nil {
			return 0, err
		}
		requirement := &models.IncidentComplianceRequirement{
			IncidentRegisterID:      id,
			ComplianceRegisterID:    complianceRequirement.ComplianceClause.ComplianceRegisterID,
			ComplianceRequirementID: *reqId,
			CreatedByID:             currentUserID,
		}
		complianceRequirements = append(complianceRequirements, requirement)
	}

	if len(complianceRequirements) > 0 {
		if err := irs.repo.CreateIncidentComplianceRequirements(complianceRequirements, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
	}
	// Save Controls
	var controls []*models.IncidentControl
	for _, ctrl := range incidentRegisterFormDTO.ControlIds {
		control := &models.IncidentControl{
			IncidentRegisterID: id,
			ControlID:          *ctrl,
			CreatedByID:        currentUserID,
		}
		controls = append(controls, control)
	}

	if len(controls) > 0 {
		if err := irs.repo.CreateIncidentControls(controls, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
	}
	// Save Risks
	var risks []*models.IncidentRisk
	for _, rsk := range incidentRegisterFormDTO.RiskIds {
		risk := &models.IncidentRisk{
			IncidentRegisterID: id,
			RiskRegisterID:     *rsk,
			CreatedByID:        currentUserID,
		}
		risks = append(risks, risk)
	}

	if len(risks) > 0 {
		if err := irs.repo.CreateIncidentRisks(risks, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	// Save Processes
	var processes []*models.IncidentProcess
	for _, processID := range incidentRegisterFormDTO.ProcessIDs {
		process := &models.IncidentProcess{
			IncidentRegisterID: id,
			ProcessID:          *processID,
			CreatedByID:        currentUserID,
		}
		processes = append(processes, process)
	}

	if len(processes) > 0 {
		if err := irs.repo.CreateIncidentProcesses(processes, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
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

	var associatedDocs []*models.IncidentAssociatedDocument
	for _, documentID := range incidentRegisterFormDTO.DocumentIds {
		document := &models.IncidentAssociatedDocument{
			IncidentRegisterID: id,
			DocumentID:         *documentID,
			CreatedByID:        currentUserID,
		}
		associatedDocs = append(associatedDocs, document)
	}

	if len(associatedDocs) > 0 {
		if err := irs.repo.CreateIncidentDocuments(associatedDocs, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	//Add timeline
	//todo
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
			if filepath.Ext(doc.Name) == ".docx" {
				if err := utility.ConvertDocxToPDF(storagePath, doc.Name, storagePath); err != nil {
					tx.Rollback()
					return 0, err
				}
			}
		}
	}

	tx.Commit()

	go func() {
		err := irs.activityLogSvc.LogActivity(
			models.CreateAction,
			models.IncidentRegister{}.TableName(),
			id,
			incidentRegisterFormDTO,
			currentUserID,
			nil,
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
			UserID:               incidentRegisterFormDTO.ResponsibleUserID,
			ModuleIdentifier:     "incident",
			NotificationType:     string(commonDTOs.AssignmentType),
			HeadingEn:            "Incident Register Assignment",
			HeadingAr:            "مهمة تسجيل الحوادث",
			TitleEn:              incidentRegisterFormDTO.Title,
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
			"RedirectURL":             fmt.Sprintf("%sincident/register/%d/overview", appURL, id),
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

	if err := irs.validateIncidentRegisterFormDTO(incidentRegisterFormDTO); err != nil {
		return err
	}

	dates, err := irs.parseIncidentRegisterDates(incidentRegisterFormDTO)
	if err != nil {
		return err
	}
	if err := irs.validateIncidentDocumentFormDTO(incidentRegisterFormDTO.Documents); err != nil {
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
	existingIncidentRegister.ResponsibleUserID = incidentRegisterFormDTO.ResponsibleUserID
	existingIncidentRegister.Title = incidentRegisterFormDTO.Title
	existingIncidentRegister.Description = incidentRegisterFormDTO.Description
	existingIncidentRegister.IncidentAt = dates.IncidentAt
	existingIncidentRegister.Location = incidentRegisterFormDTO.Location
	existingIncidentRegister.ReportedBy = incidentRegisterFormDTO.ReportedBy
	existingIncidentRegister.ReportedAt = dates.ReportedAt
	existingIncidentRegister.QuickAction = incidentRegisterFormDTO.QuickAction
	existingIncidentRegister.UpdatedByID = &currentUserID

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
	existingControlIDs := make(map[uint]bool)
	for _, ctrl := range existingIncident.Controls {
		existingControlIDs[ctrl.ControlID] = true
	}

	// Collect new control IDs
	newControlIDs := make(map[uint]bool)
	for _, ctrl := range incidentRegisterFormDTO.ControlIds {
		newControlIDs[*ctrl] = true
	}

	// Add new controls
	var controlsToAdd []*models.IncidentControl
	for _, newControlID := range incidentRegisterFormDTO.ControlIds {
		if _, exists := existingControlIDs[*newControlID]; !exists {
			// Create a new incident control record
			control := &models.IncidentControl{
				IncidentRegisterID: id,
				ControlID:          *newControlID,
				CreatedByID:        currentUserID,
			}
			controlsToAdd = append(controlsToAdd, control)
		}
	}

	// Save the new controls
	if len(controlsToAdd) > 0 {
		if err := irs.repo.CreateIncidentControls(controlsToAdd, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	// Remove controls not in the new list
	var deleteControlIDs []*uint

	for _, existingControl := range existingIncident.Controls {
		if _, exists := newControlIDs[existingControl.ControlID]; !exists {
			// Append control to delete list
			deleteControlIDs = append(deleteControlIDs, existingControl.ID)
		}
	}

	// Delete controls not in the new list
	if len(deleteControlIDs) > 0 {
		if err := irs.repo.DeleteIncidentControls(deleteControlIDs, tx); err != nil {
			tx.Rollback()
			return err
		}
	}
	// Collect existing compliance requirement IDs
	existingRequirementIDs := make(map[uint]bool)
	for _, req := range existingIncident.ComplianceRequirements {
		existingRequirementIDs[req.ComplianceRequirementID] = true
	}

	// Collect new compliance requirement IDs
	newRequirementIDs := make(map[uint]bool)
	for _, req := range incidentRegisterFormDTO.ComplianceRequirementIds {
		newRequirementIDs[*req] = true
	}

	// Add new compliance requirements
	var requirementsToAdd []*models.IncidentComplianceRequirement
	for _, newRequirementID := range incidentRegisterFormDTO.ComplianceRequirementIds {
		if _, exists := existingRequirementIDs[*newRequirementID]; !exists {
			// Create a new incident compliance requirement record
			requirement := &models.IncidentComplianceRequirement{
				IncidentRegisterID:      id,
				ComplianceRequirementID: *newRequirementID,
				CreatedByID:             currentUserID,
			}
			requirementsToAdd = append(requirementsToAdd, requirement)
		}
	}

	// Save the new compliance requirements
	if len(requirementsToAdd) > 0 {
		if err := irs.repo.CreateIncidentComplianceRequirements(requirementsToAdd, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	// Remove compliance requirements not in the new list
	var deleteRequirementIDs []*uint
	for _, existingRequirement := range existingIncident.ComplianceRequirements {
		if _, exists := newRequirementIDs[existingRequirement.ComplianceRequirementID]; !exists {
			// Append requirement to delete list
			deleteRequirementIDs = append(deleteRequirementIDs, existingRequirement.ID)
		}
	}

	// Delete compliance requirements not in the new list
	if len(deleteRequirementIDs) > 0 {
		if err := irs.repo.DeleteIncidentComplianceRequirements(deleteRequirementIDs, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	// Collect existing risk IDs
	existingRiskIDs := make(map[uint]bool)
	for _, risk := range existingIncident.Risks {
		existingRiskIDs[risk.RiskRegisterID] = true
	}

	// Collect new risk IDs
	newRiskIDs := make(map[uint]bool)
	for _, riskID := range incidentRegisterFormDTO.RiskIds {
		newRiskIDs[*riskID] = true
	}

	// Add new risks
	var risksToAdd []*models.IncidentRisk
	for _, newRiskID := range incidentRegisterFormDTO.RiskIds {
		if _, exists := existingRiskIDs[*newRiskID]; !exists {
			// Create a new incident risk record
			risk := &models.IncidentRisk{
				IncidentRegisterID: id,
				RiskRegisterID:     *newRiskID,
				CreatedByID:        currentUserID,
			}
			risksToAdd = append(risksToAdd, risk)
		}
	}

	// Save the new risks
	if len(risksToAdd) > 0 {
		if err := irs.repo.CreateIncidentRisks(risksToAdd, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	// Remove risks not in the new list
	var deleteRiskIDs []*uint
	for _, existingRisk := range existingIncident.Risks {
		if _, exists := newRiskIDs[existingRisk.RiskRegisterID]; !exists {
			// Append risk ID to delete list
			deleteRiskIDs = append(deleteRiskIDs, existingRisk.ID)
		}
	}

	// Delete risks not in the new list
	if len(deleteRiskIDs) > 0 {
		if err := irs.repo.DeleteIncidentRisks(deleteRiskIDs, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	existingProcessIDs := make(map[uint]bool)
	for _, proc := range existingIncident.Processes {
		existingProcessIDs[proc.ProcessID] = true
	}

	// Collect new process IDs
	newProcessIDs := make(map[uint]bool)
	for _, proc := range incidentRegisterFormDTO.ProcessIDs {
		newProcessIDs[*proc] = true
	}

	// Add new processes
	var processesToAdd []*models.IncidentProcess
	for _, newProcessID := range incidentRegisterFormDTO.ProcessIDs {
		if _, exists := existingProcessIDs[*newProcessID]; !exists {
			// Create a new incident process record
			process := &models.IncidentProcess{
				IncidentRegisterID: id,
				ProcessID:          *newProcessID,
				CreatedByID:        currentUserID,
			}
			processesToAdd = append(processesToAdd, process)
		}
	}

	// Save the new processes
	if len(processesToAdd) > 0 {
		if err := irs.repo.CreateIncidentProcesses(processesToAdd, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	// Remove processes not in the new list
	var deleteProcessIDs []*uint

	for _, existingProcess := range existingIncident.Processes {
		if _, exists := newProcessIDs[existingProcess.ProcessID]; !exists {
			// Append process to delete list
			deleteProcessIDs = append(deleteProcessIDs, existingProcess.ID)
		}
	}

	// Delete processes not in the new list
	if len(deleteProcessIDs) > 0 {
		if err := irs.repo.DeleteIncidentProcesses(deleteProcessIDs, tx); err != nil {
			tx.Rollback()
			return err
		}
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

	existingDocumentIDs := make(map[uint]bool)
	for _, doc := range existingIncident.AssociatedDocuments {
		existingDocumentIDs[doc.DocumentID] = true
	}

	// Collect new document IDs
	newDocumentIDs := make(map[uint]bool)
	for _, doc := range incidentRegisterFormDTO.DocumentIds {
		newDocumentIDs[*doc] = true
	}

	// Add new documents
	var documentsToAdd []*models.IncidentAssociatedDocument
	for _, newDocumentID := range incidentRegisterFormDTO.DocumentIds {
		if _, exists := existingDocumentIDs[*newDocumentID]; !exists {
			// Create a new incident document record
			document := &models.IncidentAssociatedDocument{
				IncidentRegisterID: id,
				DocumentID:         *newDocumentID,
				CreatedByID:        currentUserID,
			}
			documentsToAdd = append(documentsToAdd, document)
		}
	}

	// Save the new documents
	if len(documentsToAdd) > 0 {
		if err := irs.repo.CreateIncidentDocuments(documentsToAdd, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	// Remove documents not in the new list
	var deleteDocumentIDs []*uint

	for _, existingDocument := range existingIncident.AssociatedDocuments {
		if _, exists := newDocumentIDs[existingDocument.DocumentID]; !exists {
			// Append document to delete list
			deleteDocumentIDs = append(deleteDocumentIDs, existingDocument.ID)
		}
	}

	// Delete documents not in the new list
	if len(deleteDocumentIDs) > 0 {
		if err := irs.repo.DeleteIncidentDocuments(deleteDocumentIDs, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	//Update timeline
	//todo
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
	//documents update
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
		err := irs.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.IncidentRegister{}.TableName(),
			id,
			incidentRegisterFormDTO,
			currentUserID,
			nil,
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

		// Send notifications to old and new responsible users if the responsible user has changed
		if oldResponsibleUserID != existingIncidentRegister.ResponsibleUserID {
			// URL param for notifications
			urlParam := fmt.Sprintf(`{"id":%d}`, id)
			entity := "incident-register"

			// Old responsible user notification (Removal)
			oldResponsibleUserAction := "responsible-user-remove"
			oldResponsibleUserNotificationDTO := &commonDTOs.DBNotificationDTO{
				UserID:               oldResponsibleUserID,
				ModuleIdentifier:     "incident",
				NotificationType:     string(commonDTOs.InformationType),
				HeadingEn:            "Incident Register Responsibility Removed",
				HeadingAr:            "تمت إزالة مسؤولية تسجيل الحوادث",
				TitleEn:              existingIncidentRegister.Title,
				Entity:               &entity,
				Action:               &oldResponsibleUserAction,
				URLParam:             urlParam,
				IsSystemNotification: false,
				CreatedByID:          &currentUserID,
			}

			// New responsible user notification (Assignment)
			newResponsibleUserAction := "responsible-user-assign"
			newResponsibleUserNotificationDTO := &commonDTOs.DBNotificationDTO{
				UserID:               incidentRegisterFormDTO.ResponsibleUserID,
				ModuleIdentifier:     "incident",
				NotificationType:     string(commonDTOs.AssignmentType),
				HeadingEn:            "Incident Register Responsibility Assigned",
				HeadingAr:            "تسجيل الحوادث المسؤولية المعينة",
				TitleEn:              existingIncidentRegister.Title,
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
				"RedirectURL":             fmt.Sprintf("%sincident/register/%d/overview", appURL, id),
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
				"RedirectURL":             fmt.Sprintf("%sincident/register/%d/overview", appURL, id),
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

		// Check if controls exist
		if len(incident.Controls) > 0 {
			var deleteControlIds []*uint
			for _, ctrl := range incident.Controls {
				deleteControlIds = append(deleteControlIds, ctrl.ID)
			}
			if err := irs.repo.DeleteIncidentControls(deleteControlIds, tx); err != nil {

				tx.Rollback()
				return nil, err
			}
		}

		// Delete each compliance requirement if requirements exist
		if len(incident.ComplianceRequirements) > 0 {
			var deleteRequirementIds []*uint
			for _, req := range incident.ComplianceRequirements {
				deleteRequirementIds = append(deleteRequirementIds, req.ID)
			}
			if err := irs.repo.DeleteIncidentComplianceRequirements(deleteRequirementIds, tx); err != nil {
				tx.Rollback()
				return nil, err
			}

		}
		// Check if Risk exist
		if len(incident.Risks) > 0 {
			var deleteRiskIds []*uint
			for _, rsk := range incident.Risks {
				deleteRiskIds = append(deleteRiskIds, rsk.ID)
			}
			if err := irs.repo.DeleteIncidentRisks(deleteRiskIds, tx); err != nil {

				tx.Rollback()
				return nil, err
			}
		}

		if len(incident.Processes) > 0 {
			var deleteControlIDs []*uint
			for _, process := range incident.Processes {
				deleteControlIDs = append(deleteControlIDs, process.ID)
			}
			if err := irs.repo.DeleteIncidentProcesses(deleteControlIDs, tx); err != nil {
				tx.Rollback()
				return nil, err
			}
		}
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

		if len(incident.AssociatedDocuments) > 0 {
			var deleteDocumentIds []*uint
			for _, document := range incident.AssociatedDocuments {
				deleteDocumentIds = append(deleteDocumentIds, document.ID)
			}
			if err := irs.repo.DeleteIncidentDocuments(deleteDocumentIds, tx); err != nil {
				tx.Rollback()
				return nil, err
			}
		}

		//Delete timeline
		//todo
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
				nil,
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
					UserID:               *notification.OwnerID,
					ModuleIdentifier:     "incident",
					NotificationType:     string(commonDTOs.DeletionType),
					HeadingEn:            "Incident Register Deleted",
					HeadingAr:            "تم حذف سجل الحوادث",
					TitleEn:              notification.TitleEn,
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
	existingIncident.UpdatedByID = &currentUserID

	err = irs.repo.Update(existingIncident, tx)

	if err != nil {
		return err
	}

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
	//Controls exist
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
	//Compliance Requirements Exist
	requirements := incidentRegisterFormDTO.ComplianceRequirementIds

	// Create a map to store requirement IDs
	requirementIDs := make(map[uint]bool)

	// Check if each requirement exists
	for _, requirement := range requirements {
		// Check if the requirement ID already exists in the map
		if requirementIDs[*requirement] {
			// If the requirement ID already exists, return a validation error
			return utility.NewValidationError(fmt.Errorf("duplicate compliance requirement ID: %d", requirement), "duplicate", "compliance requirement")
		}

		// Check if the requirement exists in the database
		if _, err := irs.complianceRequirementSvc.GetByComplianceRequirementID(*requirement); err != nil {
			return utility.NewValidationError(fmt.Errorf("compliance requirement does not exist for ID: %d", requirement), "not_exist", "compliance requirement")
		}

		// Add the requirement ID to the map
		requirementIDs[*requirement] = true
	}
	//Risks exist
	risks := incidentRegisterFormDTO.RiskIds
	// Create a map to store risk IDs
	riskIDs := make(map[uint]bool)

	// Check if each risk exists
	for _, risk := range risks {
		// Check if the risk ID already exists in the map
		if riskIDs[*risk] {
			// If the risk ID already exists, return a validation error
			return utility.NewValidationError(fmt.Errorf("duplicate risk ID: %d", risk), "duplicate", "risk")
		}

		// Check if the risk exists in the database
		if _, err := irs.riskRegisterSvc.GetByID(*risk); err != nil {
			return utility.NewValidationError(fmt.Errorf("risk does not exist for ID: %d", risk), "not_exist", "risk")
		}

		// Add the risk ID to the map
		riskIDs[*risk] = true
	}

	processIds := incidentRegisterFormDTO.ProcessIDs

	processIDs := make(map[uint]bool)

	for _, processId := range processIds {
		if processIDs[*processId] {
			return utility.NewValidationError(fmt.Errorf("duplicate process ID: %d", processId), "duplicate", "process")
		}

		if _, err := irs.processSvc.GetByID(*processId); err != nil {
			return utility.NewValidationError(fmt.Errorf("process does not exist for ID: %d", processId), "not_exist", "process")
		}

		processIDs[*processId] = true
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

	return nil
}

func (irs *IncidentRegisterService) parseIncidentRegisterDates(dto *incidentDTOs.IncidentRegisterFormDTO) (*incidentDTOs.IncidentRegisterParsedDates, error) {
	parsed := &incidentDTOs.IncidentRegisterParsedDates{}

	var err error
	if dto.ReportedAt != "" {
		parsed.ReportedAt, err = utility.ParseDateTime(dto.ReportedAt)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid reported date"), "invalid_date", "reported at")
		}
	}

	if dto.IncidentAt != "" {
		parsed.IncidentAt, err = utility.ParseDateTime(dto.IncidentAt)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid incident date"), "invalid_date", "incident at")
		}
	}

	return parsed, nil
}

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
			"Communicated By",
			"Communicated To",
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
			"Communicated By",
			"Communicated To",
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
				incident.Location,
				incident.ReportedUserFirstNameAr + " " + incident.ReportedUserLastNameAr,
				incident.ResponsibleUserFirstNameAr + " " + incident.ResponsibleUserLastNameAr,
				incident.ReportedAt,
				incident.IncidentStatusTitleAr,
				incident.QuickAction,
				// incident.InvestigationInitiatedAt,
				// incident.InvestigationDescription,
				// *incident.CommunicatedByFirstNameAr + " " + *incident.CommunicatedByLastNameAr,
				// incident.CommunicatedViaAr,
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
				incident.Location,
				incident.ReportedUserFirstNameEn + " " + incident.ReportedUserFirstNameEn,
				incident.ResponsibleUserFirstNameEn + " " + incident.ResponsibleUserLastNameEn,
				incident.ReportedAt,
				incident.IncidentStatusTitleEn,
				incident.QuickAction,
				// incident.InvestigationInitiatedAt,
				// incident.InvestigationDescription,
				// incident.CommunicatedByFirstNameEn + " " + *incident.CommunicatedByLastNameEn,
				// incident.CommunicatedViaEn,
			}
		}
	}

	// Generate the Excel file
	return utility.GenerateExcelFile(headers, rows)
}
