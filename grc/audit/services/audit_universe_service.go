package services

import (
	"fmt"
	"jethur/grc/audit/dtos"
	"jethur/grc/audit/repositories"
	complianceServices "jethur/grc/compliance/services"
	documentServices "jethur/grc/documents/services"
	incidentServices "jethur/grc/incident/services"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	riskServices "jethur/grc/risk/services"
	"jethur/utility"
)

type AuditUniverseServiceInterface interface {
	Create(currentUserID uint, auditUniverseFormDTO *dtos.AuditUniverseFormDTO) error
	MarkedAsAuditableItem(ids []*uint, entityType interface{}, tx utility.DatabaseTransaction) error
	Delete(currentUserID uint, auditUniverseFormDTO *dtos.AuditUniverseFormDTO) error
	UnmarkedAsAuditableItem(ids []*uint, entityType interface{}, tx utility.DatabaseTransaction) error

	getOrganizationUnitID(entityType string, id uint) (uint, error)
}

// Service and implementations
type AuditUniverseService struct {
	repo                     repositories.AuditUniverseRepositoryInterface
	activityLogSvc           logServices.ActivityLogServiceInterface
	auditStatusSvc           AuditStatusServiceInterface
	auditorRiskRatingSvc     AuditorRiskRatingServiceInterface
	processSvc               organizationServices.ProcessServiceInterface
	controlSvc               organizationServices.ControlServiceInterface
	userSvc                  organizationServices.UserServiceInterface
	complianceRequirementSvc complianceServices.ComplianceRequirementServiceInterface
	riskRegisterSvc          riskServices.RiskRegisterServiceInterface
	incidentRegisterSvc      incidentServices.IncidentRegisterServiceInterface
	documentSvc              documentServices.DocumentServiceInterface
}

func NewAuditUniverseService(
	repo repositories.AuditUniverseRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	auditStatusSvc AuditStatusServiceInterface,
	auditorRiskRatingSvc AuditorRiskRatingServiceInterface,
	processSvc organizationServices.ProcessServiceInterface,
	controlSvc organizationServices.ControlServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	complianceRequirementSvc complianceServices.ComplianceRequirementServiceInterface,
	riskRegisterSvc riskServices.RiskRegisterServiceInterface,
	incidentRegisterSvc incidentServices.IncidentRegisterServiceInterface,
	documentSvc documentServices.DocumentServiceInterface,

) AuditUniverseServiceInterface {
	return &AuditUniverseService{
		repo:                     repo,
		activityLogSvc:           activityLogSvc,
		auditStatusSvc:           auditStatusSvc,
		auditorRiskRatingSvc:     auditorRiskRatingSvc,
		processSvc:               processSvc,
		controlSvc:               controlSvc,
		userSvc:                  userSvc,
		complianceRequirementSvc: complianceRequirementSvc,
		riskRegisterSvc:          riskRegisterSvc,
		incidentRegisterSvc:      incidentRegisterSvc,
		documentSvc:              documentSvc,
	}
}

func (aus *AuditUniverseService) Create(currentUserID uint, auditUniverseFormDTO *dtos.AuditUniverseFormDTO) error {
	if err := aus.validateAuditUniverseFormDTO(auditUniverseFormDTO); err != nil {
		return err
	}

	// Retrieve audit status
	auditStatus, err := aus.auditStatusSvc.GetByType("pending")
	if err != nil {
		return err
	}

	var auditableItems []*models.AuditableItem

	// Helper function to add items
	addItems := func(ids []*uint, fieldName string) error {
		for _, id := range ids {
			auditableItem := &models.AuditableItem{
				Comment:             auditUniverseFormDTO.Comment,
				AuditorRiskRatingID: auditUniverseFormDTO.AuditorRiskRatingID,
				AuditStatusID:       uint(auditStatus.ID),
				CreatedByID:         currentUserID,
			}
			var orgUnitID uint
			switch fieldName {
			case "process":
				auditableItem.ProcessID = id
				auditableItem.Type = "process"
			case "control":
				auditableItem.ControlID = id
				auditableItem.Type = "control"
			case "document":
				auditableItem.DocumentID = id
				auditableItem.Type = "document"
			case "incident":
				auditableItem.IncidentRegisterID = id
				auditableItem.Type = "incident"
			case "compliance":
				auditableItem.ComplianceRequirementID = id
				auditableItem.Type = "compliance_requirement"
			case "risk":
				auditableItem.RiskRegisterID = id
				auditableItem.Type = "risk"
			case "user":
				auditableItem.UserID = id
				auditableItem.Type = "user"
			}
			orgUnitID, err = aus.getOrganizationUnitID(fieldName, *id)
			if err != nil {
				return err
			}
			if orgUnitID != 0 {
				auditableItem.ResponsibleUnitID = orgUnitID
			}
			auditableItems = append(auditableItems, auditableItem)
		}
		return nil
	}

	if err := addItems(auditUniverseFormDTO.ProcessIDs, "process"); err != nil {
		return err
	}
	if err := addItems(auditUniverseFormDTO.ControlIDs, "control"); err != nil {
		return err
	}
	if err := addItems(auditUniverseFormDTO.DocumentIDs, "document"); err != nil {
		return err
	}
	if err := addItems(auditUniverseFormDTO.IncidentRegisterIDs, "incident"); err != nil {
		return err
	}
	if err := addItems(auditUniverseFormDTO.ComplianceRequirementIDs, "compliance"); err != nil {
		return err
	}
	if err := addItems(auditUniverseFormDTO.RiskRegisterIDs, "risk"); err != nil {
		return err
	}
	if err := addItems(auditUniverseFormDTO.UserIDs, "user"); err != nil {
		return err
	}
	// Create transaction for critical operations
	tx := aus.repo.BeginTransaction()

	// Create all auditable items in one go
	if len(auditableItems) > 0 {
		if err := aus.repo.Create(auditableItems, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := aus.MarkedAsAuditableItem(auditUniverseFormDTO.ProcessIDs, &models.Process{}, tx); err != nil {
		tx.Rollback()
		return err
	}
	if err := aus.MarkedAsAuditableItem(auditUniverseFormDTO.ControlIDs, &models.Control{}, tx); err != nil {
		tx.Rollback()
		return err
	}
	if err := aus.MarkedAsAuditableItem(auditUniverseFormDTO.DocumentIDs, &models.Document{}, tx); err != nil {
		tx.Rollback()
		return err
	}
	if err := aus.MarkedAsAuditableItem(auditUniverseFormDTO.IncidentRegisterIDs, &models.IncidentRegister{}, tx); err != nil {
		tx.Rollback()
		return err
	}
	if err := aus.MarkedAsAuditableItem(auditUniverseFormDTO.ComplianceRequirementIDs, &models.ComplianceRequirement{}, tx); err != nil {
		tx.Rollback()
		return err
	}
	if err := aus.MarkedAsAuditableItem(auditUniverseFormDTO.RiskRegisterIDs, &models.RiskRegister{}, tx); err != nil {
		tx.Rollback()
		return err
	}
	if err := aus.MarkedAsAuditableItem(auditUniverseFormDTO.UserIDs, &models.User{}, tx); err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	go func() {
		err := aus.activityLogSvc.LogActivity(
			models.CreateAction,
			models.AuditableItem{}.TableName(),
			0, //todo
			auditUniverseFormDTO,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.AuditableItem{}.TableName(),
				"recordID":  0, //todo
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return nil
}

func (aus *AuditUniverseService) getOrganizationUnitID(entityType string, id uint) (uint, error) {
	var orgUnitID uint

	switch entityType {
	case "process":
		process, err := aus.processSvc.GetByID(id)
		if err != nil {
			return 0, fmt.Errorf("failed to get process: %w", err)
		}
		if process == nil {
			return 0, nil
		}
		orgUnitID = process.OwnerUnitID
	case "control":
		control, err := aus.controlSvc.GetByID(id)
		if err != nil {
			return 0, fmt.Errorf("failed to get control: %w", err)
		}
		if control == nil {
			return 0, nil
		}
		orgUnitID = control.OwnerUnitID
	case "document":
		document, err := aus.documentSvc.GetByID(id)
		if err != nil {
			return 0, fmt.Errorf("failed to get document: %w", err)
		}
		if document == nil {
			return 0, nil
		}
		orgUnitID = document.ResponsibleUnitID
	case "incident":
		incident, err := aus.incidentRegisterSvc.GetByID(id)
		if err != nil {
			return 0, fmt.Errorf("failed to get incident: %w", err)
		}
		if incident == nil {
			return 0, nil
		}
		orgUnitID = incident.ResponsibleUnitID
	case "compliance":
		compliance, err := aus.complianceRequirementSvc.GetByComplianceRequirementID(id)
		if err != nil {
			return 0, fmt.Errorf("failed to get compliance requirement: %w", err)
		}
		if compliance == nil {
			return 0, nil
		}
		userUnit, err := aus.userSvc.GetUserPrimaryUnitID(compliance.ResponsibleUserID)
		if err != nil {
			return 0, fmt.Errorf("failed to get user primary unit ID: %w", err)
		}
		orgUnitID = userUnit
	case "risk":
		risk, err := aus.riskRegisterSvc.GetByID(id)
		if err != nil {
			return 0, fmt.Errorf("failed to get risk: %w", err)
		}
		if risk == nil {
			return 0, nil
		}
		orgUnitID = risk.ResponsibleUnitID
	case "user":
		userUnit, err := aus.userSvc.GetUserPrimaryUnitID(id)
		if err != nil {
			return 0, fmt.Errorf("failed to get user primary unit ID: %w", err)
		}
		orgUnitID = userUnit
	default:
		return 0, fmt.Errorf("unknown entity type: %s", entityType)
	}

	return orgUnitID, nil
}

func (aus *AuditUniverseService) MarkedAsAuditableItem(ids []*uint, entityType interface{}, tx utility.DatabaseTransaction) error {
	// Convert []*uint to []uint
	var uintIDs []uint
	for _, id := range ids {
		if id != nil {
			uintIDs = append(uintIDs, *id)
		}
	}

	if len(uintIDs) > 0 {
		updates := map[string]interface{}{"marked_as_auditable_item": true}
		if err := aus.repo.UpdateMarkedAsAuditableItem(uintIDs, entityType, updates, tx); err != nil {
			return err
		}
	}

	return nil
}

func (aus *AuditUniverseService) UnmarkedAsAuditableItem(ids []*uint, entityType interface{}, tx utility.DatabaseTransaction) error {
	// Convert []*uint to []uint
	var uintIDs []uint
	for _, id := range ids {
		if id != nil {
			uintIDs = append(uintIDs, *id)
		}
	}

	if len(uintIDs) > 0 {
		updates := map[string]interface{}{"marked_as_auditable_item": false}
		if err := aus.repo.UpdateMarkedAsAuditableItem(uintIDs, entityType, updates, tx); err != nil {
			return err
		}
	}
	return nil
}

func (aus *AuditUniverseService) Delete(currentUserID uint, auditUniverseFormDTO *dtos.AuditUniverseFormDTO) error {
	tx := aus.repo.BeginTransaction()
	if err := aus.validateAuditUniverseFormDTO(auditUniverseFormDTO); err != nil {
		return err
	}
	// Initialize slices to hold all IDs
	var processIDs, controlIDs, documentIDs, incidentIDs, complianceIDs, riskIDs, userIDs []uint

	// Helper function to collect IDs
	collectIDs := func(ids []*uint, target *[]uint) {
		for _, id := range ids {
			if id != nil {
				*target = append(*target, *id)
			}
		}
	}

	// Collect all relevant IDs
	collectIDs(auditUniverseFormDTO.ProcessIDs, &processIDs)
	collectIDs(auditUniverseFormDTO.ControlIDs, &controlIDs)
	collectIDs(auditUniverseFormDTO.DocumentIDs, &documentIDs)
	collectIDs(auditUniverseFormDTO.IncidentRegisterIDs, &incidentIDs)
	collectIDs(auditUniverseFormDTO.ComplianceRequirementIDs, &complianceIDs)
	collectIDs(auditUniverseFormDTO.RiskRegisterIDs, &riskIDs)
	collectIDs(auditUniverseFormDTO.UserIDs, &userIDs)

	// Common function to get auditable item IDs and delete them
	getAuditableItemIDs := func(ids []uint, fieldName string) ([]uint, error) {
		if len(ids) == 0 {
			return nil, nil
		}

		auditableItems, err := aus.repo.GetByField(fieldName, ids, tx)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		var itemIDs []uint
		for _, item := range auditableItems {
			itemIDs = append(itemIDs, item.ID)
		}

		return itemIDs, nil
	}

	deleteItems := func(ids []uint, fieldName string) error {
		itemIDs, err := getAuditableItemIDs(ids, fieldName)
		if err != nil {
			return err
		}

		if len(itemIDs) > 0 {
			// Now delete the auditable items themselves
			if err := aus.repo.DeleteByIDs(itemIDs, tx); err != nil {
				tx.Rollback()
				return err
			}
		}
		return nil
	}

	// Delete all auditable items by IDs
	if err := deleteItems(processIDs, "process_id"); err != nil {
		return err
	}
	if err := deleteItems(controlIDs, "control_id"); err != nil {
		return err
	}
	if err := deleteItems(documentIDs, "document_id"); err != nil {
		return err
	}
	if err := deleteItems(incidentIDs, "incident_register_id"); err != nil {
		return err
	}
	if err := deleteItems(complianceIDs, "compliance_requirement_id"); err != nil {
		return err
	}
	if err := deleteItems(riskIDs, "risk_register_id"); err != nil {
		return err
	}
	if err := deleteItems(userIDs, "user_id"); err != nil {
		return err
	}

	// Unmark entities as auditable
	if err := aus.UnmarkedAsAuditableItem(auditUniverseFormDTO.ProcessIDs, &models.Process{}, tx); err != nil {
		tx.Rollback()
		return err
	}
	if err := aus.UnmarkedAsAuditableItem(auditUniverseFormDTO.ControlIDs, &models.Control{}, tx); err != nil {
		tx.Rollback()
		return err
	}
	if err := aus.UnmarkedAsAuditableItem(auditUniverseFormDTO.DocumentIDs, &models.Document{}, tx); err != nil {
		tx.Rollback()
		return err
	}
	if err := aus.UnmarkedAsAuditableItem(auditUniverseFormDTO.IncidentRegisterIDs, &models.IncidentRegister{}, tx); err != nil {
		tx.Rollback()
		return err
	}
	if err := aus.UnmarkedAsAuditableItem(auditUniverseFormDTO.ComplianceRequirementIDs, &models.ComplianceRequirement{}, tx); err != nil {
		tx.Rollback()
		return err
	}
	if err := aus.UnmarkedAsAuditableItem(auditUniverseFormDTO.RiskRegisterIDs, &models.RiskRegister{}, tx); err != nil {
		tx.Rollback()
		return err
	}
	if err := aus.UnmarkedAsAuditableItem(auditUniverseFormDTO.UserIDs, &models.User{}, tx); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	// Asynchronous activity logging
	go func() {
		err := aus.activityLogSvc.LogActivity(
			models.DeleteAction,
			models.AuditableItem{}.TableName(),
			0, // No specific item to log for delete
			auditUniverseFormDTO,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.DeleteAction,
				"tableName": models.AuditableItem{}.TableName(),
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return nil
}

func (aus *AuditUniverseService) validateAuditUniverseFormDTO(auditUniverseFormDTO *dtos.AuditUniverseFormDTO) error {
	// Validate Auditor Risk Rating
	if _, err := aus.auditorRiskRatingSvc.GetByID(auditUniverseFormDTO.AuditorRiskRatingID); err != nil {
		return utility.NewValidationError(fmt.Errorf("auditor risk rating does not exist"), "not_exist", "auditor risk rating")
	}

	// Validate Process IDs
	processIDs := auditUniverseFormDTO.ProcessIDs
	processIDMap := make(map[uint]bool)

	for _, id := range processIDs {
		if processIDMap[*id] {
			return utility.NewValidationError(fmt.Errorf("duplicate process ID: %d", *id), "duplicate", "process")
		}

		if _, err := aus.processSvc.GetByID(*id); err != nil {
			return utility.NewValidationError(fmt.Errorf("process does not exist for ID: %d", *id), "not_exist", "process")
		}

		processIDMap[*id] = true
	}

	// Validate Control IDs
	controlIDs := auditUniverseFormDTO.ControlIDs
	controlIDMap := make(map[uint]bool)

	for _, id := range controlIDs {
		if controlIDMap[*id] {
			return utility.NewValidationError(fmt.Errorf("duplicate control ID: %d", *id), "duplicate", "control")
		}

		if _, err := aus.controlSvc.GetByID(*id); err != nil {
			return utility.NewValidationError(fmt.Errorf("control does not exist for ID: %d", *id), "not_exist", "control")
		}

		controlIDMap[*id] = true
	}

	// Validate Document IDs
	documentIDs := auditUniverseFormDTO.DocumentIDs
	documentIDMap := make(map[uint]bool)

	for _, id := range documentIDs {
		if documentIDMap[*id] {
			return utility.NewValidationError(fmt.Errorf("duplicate document ID: %d", *id), "duplicate", "document")
		}

		if _, err := aus.documentSvc.GetByID(*id); err != nil {
			return utility.NewValidationError(fmt.Errorf("document does not exist for ID: %d", *id), "not_exist", "document")
		}

		documentIDMap[*id] = true
	}

	// Validate Incident Register IDs
	incidentRegisterIDs := auditUniverseFormDTO.IncidentRegisterIDs
	incidentRegisterIDMap := make(map[uint]bool)

	for _, id := range incidentRegisterIDs {
		if incidentRegisterIDMap[*id] {
			return utility.NewValidationError(fmt.Errorf("duplicate incident register ID: %d", *id), "duplicate", "incident")
		}

		if _, err := aus.incidentRegisterSvc.GetByID(*id); err != nil {
			return utility.NewValidationError(fmt.Errorf("incident register does not exist for ID: %d", *id), "not_exist", "incident")
		}

		incidentRegisterIDMap[*id] = true
	}

	// Validate Compliance Requirement IDs
	complianceRequirementIDs := auditUniverseFormDTO.ComplianceRequirementIDs
	complianceRequirementIDMap := make(map[uint]bool)

	for _, id := range complianceRequirementIDs {
		if complianceRequirementIDMap[*id] {
			return utility.NewValidationError(fmt.Errorf("duplicate compliance requirement ID: %d", *id), "duplicate", "compliance requirement")
		}

		if _, err := aus.complianceRequirementSvc.GetByComplianceRequirementID(*id); err != nil {
			return utility.NewValidationError(fmt.Errorf("compliance requirement does not exist for ID: %d", *id), "not_exist", "compliance requirement")
		}

		complianceRequirementIDMap[*id] = true
	}

	// Validate Risk Register IDs
	riskRegisterIDs := auditUniverseFormDTO.RiskRegisterIDs
	riskRegisterIDMap := make(map[uint]bool)

	for _, id := range riskRegisterIDs {
		if riskRegisterIDMap[*id] {
			return utility.NewValidationError(fmt.Errorf("duplicate risk register ID: %d", *id), "duplicate", "risk register")
		}

		if _, err := aus.riskRegisterSvc.GetByID(*id); err != nil {
			return utility.NewValidationError(fmt.Errorf("risk register does not exist for ID: %d", *id), "not_exist", "risk register")
		}

		riskRegisterIDMap[*id] = true
	}

	// Validate User IDs
	userIDs := auditUniverseFormDTO.UserIDs
	userIDMap := make(map[uint]bool)

	for _, id := range userIDs {
		if userIDMap[*id] {
			return utility.NewValidationError(fmt.Errorf("duplicate user ID: %d", *id), "duplicate", "user")
		}

		if _, err := aus.userSvc.GetByID(*id); err != nil {
			return utility.NewValidationError(fmt.Errorf("user does not exist for ID: %d", *id), "not_exist", "user")
		}

		userIDMap[*id] = true
	}

	return nil
}
