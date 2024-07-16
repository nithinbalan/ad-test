package services

import (
	"errors"
	"fmt"
	complianceServices "jethur/grc/compliance/services"
	incidentDTOs "jethur/grc/incident/dtos"
	incidentRepositories "jethur/grc/incident/repositories"
	logServices "jethur/grc/logs/services"
	masterRepositories "jethur/grc/masters/repositories"
	masterServices "jethur/grc/masters/services"
	organizationServices "jethur/grc/organization/services"
	riskServices "jethur/grc/risk/services"
	"time"

	"jethur/grc/models"
	"jethur/utility"
	"path/filepath"
	"strconv"

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
}

// Service and implementations
type IncidentRegisterService struct {
	repo                     incidentRepositories.IncidentRegisterRepositoryInterface
	activityLogSvc           logServices.ActivityLogServiceInterface
	incidentCategorySvc      IncidentCategoryServiceInterface
	incidentSubCategorySvc   IncidentSubCategoryServiceInterface
	incidentTypeSvc          IncidentTypeServiceInterface
	incidentStatusSvc        IncidentStatusServiceInterface
	severityLevelSvc         masterServices.SeverityLevelServiceInterface
	controlSvc               organizationServices.ControlServiceInterface
	complianceRequirementSvc complianceServices.ComplianceRequirementServiceInterface
	riskRegisterSvc          riskServices.RiskRegisterServiceInterface
	processSvc               organizationServices.ProcessServiceInterface
	slaRegisterSvc           complianceServices.SLARegisterServiceInterface
	incidentTimelineSvc      IncidentTimelineServiceInterface
}

func NewIncidentRegisterService(
	repo incidentRepositories.IncidentRegisterRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	incidentCategorySvc IncidentCategoryServiceInterface,
	incidentSubCategorySvc IncidentSubCategoryServiceInterface,
	incidentTypeSvc IncidentTypeServiceInterface,
	incidentStatusSvc IncidentStatusServiceInterface,
	severityLevelSvc masterRepositories.SeverityLevelRepositoryInterface,
	controlSvc organizationServices.ControlServiceInterface,
	complianceRequirementSvc complianceServices.ComplianceRequirementServiceInterface,
	riskRegisterSvc riskServices.RiskRegisterServiceInterface,
	processSvc organizationServices.ProcessServiceInterface,
	slaRegisterSvc complianceServices.SLARegisterServiceInterface,
	incidentTimelineSvc IncidentTimelineServiceInterface,
) IncidentRegisterServiceInterface {
	return &IncidentRegisterService{
		repo:                     repo,
		activityLogSvc:           activityLogSvc,
		incidentCategorySvc:      incidentCategorySvc,
		incidentSubCategorySvc:   incidentSubCategorySvc,
		incidentTypeSvc:          incidentTypeSvc,
		incidentStatusSvc:        incidentStatusSvc,
		severityLevelSvc:         severityLevelSvc,
		controlSvc:               controlSvc,
		complianceRequirementSvc: complianceRequirementSvc,
		riskRegisterSvc:          riskRegisterSvc,
		processSvc:               processSvc,
		slaRegisterSvc:           slaRegisterSvc,
		incidentTimelineSvc:      incidentTimelineSvc,
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

	//todo :- also add extra fields which required in details screen
	// Format time fields without time component
	if incident.IncidentAt != nil {
		incidentAt := incident.IncidentAt.Format("2006-01-02")
		incidentDTO.IncidentAtString = incidentAt
	}
	if incident.ReportedAt != nil {
		reportedAt := incident.ReportedAt.Format("2006-01-02")
		incidentDTO.ReportedAtStaring = reportedAt
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

	incidentRegister := &models.IncidentRegister{
		ReferenceNumber:       referenceNumber,
		IncidentCategoryID:    incidentRegisterFormDTO.IncidentCategoryID,
		IncidentSubCategoryID: *incidentRegisterFormDTO.IncidentSubCategoryID,
		IncidentTypeID:        incidentRegisterFormDTO.IncidentTypeID,
		SeverityLevelID:       *incidentRegisterFormDTO.SeverityLevelID,
		IncidentStatusID:      incidentRegisterFormDTO.IncidentStatusID,
		Title:                 incidentRegisterFormDTO.Title,
		Description:           incidentRegisterFormDTO.Description,
		IncidentAt:            dates.IncidentAt,
		Location:              incidentRegisterFormDTO.Location,
		ReportedBy:            incidentRegisterFormDTO.ReportedBy,
		ReportedAt:            dates.ReportedAt,
		QuickAction:           incidentRegisterFormDTO.QuickAction,
		CreatedByID:           currentUserID,
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
	//Add timeline
	//todo
	timelineType, err := irs.incidentTimelineSvc.GetByType("reporting")
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	//todo
	timelineStatus, err := irs.incidentTimelineSvc.GetByStatusType("pending")
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	timelineFormDTO := &incidentDTOs.IncidentTimelineFormDTO{
		IncidentRegisterID:     id,
		IncidentTimelineTypeID: timelineType.ID,
		TimelineStatusID:       timelineStatus.ID,
		PerformedByID:          currentUserID,
		PerformedAt:            time.Now(),
	}
	if err := irs.incidentTimelineSvc.Create(timelineFormDTO, tx); err != nil {
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
	}()

	return id, nil
}

func (irs *IncidentRegisterService) Update(id uint, currentUserID uint, incidentRegisterFormDTO *incidentDTOs.IncidentRegisterFormDTO) error {
	tx := irs.repo.BeginTransaction()

	existingIncidentRegister, err := irs.repo.GetByID(id)

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
	existingIncidentRegister.IncidentCategoryID = incidentRegisterFormDTO.IncidentCategoryID
	existingIncidentRegister.IncidentSubCategoryID = *incidentRegisterFormDTO.IncidentSubCategoryID
	existingIncidentRegister.IncidentTypeID = incidentRegisterFormDTO.IncidentTypeID
	existingIncidentRegister.SeverityLevelID = *incidentRegisterFormDTO.SeverityLevelID
	existingIncidentRegister.IncidentStatusID = incidentRegisterFormDTO.IncidentStatusID
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
	//Update timeline
	//todo
	timelineType, err := irs.incidentTimelineSvc.GetByType("reporting")
	if err != nil {
		tx.Rollback()
		return err
	}
	timelineFormDTO := &incidentDTOs.IncidentTimelineFormDTO{
		PerformedByID: currentUserID,
		PerformedAt:   time.Now(),
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
	}()

	return nil
}

func (irs *IncidentRegisterService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	tx := irs.repo.BeginTransaction()
	var failedIds []uint
	var successfulIds []uint
	documentsToDelete := make(map[uint][]*incidentDTOs.IncidentDocumentDTO)

	for _, id := range ids {
		incident, err := irs.GetFullDetailsByID(id, tx)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
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
		//Delete timeline
		//todo
		timelineType, err := irs.incidentTimelineSvc.GetByType("reporting")
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		if err := irs.incidentTimelineSvc.Delete(id, timelineType.ID, tx); err != nil {
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
		}()
	}

	return failedIds, nil
}

func (irs *IncidentRegisterService) validateIncidentRegisterFormDTO(incidentRegisterFormDTO *incidentDTOs.IncidentRegisterFormDTO) error {
	if _, err := irs.incidentCategorySvc.GetByID(incidentRegisterFormDTO.IncidentCategoryID); err != nil {
		return utility.NewValidationError(fmt.Errorf("incident category does not exist"), "not_exist", "incident category")
	}

	if incidentRegisterFormDTO.IncidentSubCategoryID != nil {
		if _, err := irs.incidentSubCategorySvc.GetByID(*incidentRegisterFormDTO.IncidentSubCategoryID); err != nil {
			return utility.NewValidationError(fmt.Errorf("incident sub category does not exist"), "not_exist", "incident sub category")
		}
	}

	if _, err := irs.incidentTypeSvc.GetByID(incidentRegisterFormDTO.IncidentTypeID); err != nil {
		return utility.NewValidationError(fmt.Errorf("incident type does not exist"), "not_exist", "incident type")
	}

	if _, err := irs.incidentStatusSvc.GetByID(incidentRegisterFormDTO.IncidentStatusID); err != nil {
		return utility.NewValidationError(fmt.Errorf("incident status does not exist"), "not_exist", "incident status")
	}
	if incidentRegisterFormDTO.SeverityLevelID != nil {
		if _, err := irs.severityLevelSvc.GetByID(*incidentRegisterFormDTO.SeverityLevelID); err != nil {
			return utility.NewValidationError(fmt.Errorf("severity level does not exist"), "not_exist", "severity level")
		}
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
		parsed.ReportedAt, err = utility.ParseDate(dto.ReportedAt)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid reported date"), "invalid_date", "reported at")
		}
	}

	if dto.IncidentAt != "" {
		parsed.IncidentAt, err = utility.ParseDate(dto.IncidentAt)
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
