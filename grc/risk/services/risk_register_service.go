package services

import (
	"errors"
	"fmt"
	complianceServices "jethur/grc/compliance/services"
	logServices "jethur/grc/logs/services"
	masterServices "jethur/grc/masters/services"
	organizationServices "jethur/grc/organization/services"
	riskDTOs "jethur/grc/risk/dtos"
	riskRepositories "jethur/grc/risk/repositories"
	"path/filepath"
	"strconv"

	"jethur/grc/models"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type RiskRegisterServiceInterface interface {
	GetAll(params *riskDTOs.RiskRegisterListQueryParamsDTO) (*riskDTOs.RiskRegisterListResponseDTO, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*riskDTOs.RiskRegisterDetailDTO, error)
	GetByID(id uint) (*models.RiskRegister, error)
	Create(currentUserID uint, riskRegisterFormDTO *riskDTOs.RiskRegisterFormDTO) (uint, error)
	Update(id uint, currentUserID uint, risk *riskDTOs.RiskRegisterFormDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	UpdateResponseStrategy(riskID uint, currentUserID uint, riskResponseStrategyFormDTO *riskDTOs.RiskResponseStrategyFormDTO) error
	UpdateRiskRating(riskID uint, currentUserID uint, riskEvaluationDetails *models.RiskEvaluation, tx utility.DatabaseTransaction) error
}

// Service and implementations
type RiskRegisterService struct {
	repo                     riskRepositories.RiskRegisterRepositoryInterface
	activityLogSvc           logServices.ActivityLogServiceInterface
	riskCategorySvc          RiskCategoryServiceInterface
	riskSourceSvc            RiskSourceServiceInterface
	userSvc                  organizationServices.UserServiceInterface
	riskStatusSvc            RiskStatusServiceInterface
	reviewFrequencySvc       masterServices.ReviewFrequencyServiceInterface
	controlSvc               organizationServices.ControlServiceInterface
	complianceRequirementSvc complianceServices.ComplianceRequirementServiceInterface
	processSvc               organizationServices.ProcessServiceInterface
	strategicObjectiveSvc    organizationServices.StrategicObjectiveServiceInterface
	slaRegisterSvc           complianceServices.SLARegisterServiceInterface
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
) RiskRegisterServiceInterface {
	return &RiskRegisterService{
		repo:                     repo,
		activityLogSvc:           activityLogSvc,
		riskCategorySvc:          riskCategorySvc,
		riskSourceSvc:            riskSourceSvc,
		userSvc:                  userSvc,
		riskStatusSvc:            riskStatusSvc,
		reviewFrequencySvc:       reviewFrequencySvc,
		controlSvc:               controlSvc,
		complianceRequirementSvc: complianceRequirementSvc,
		processSvc:               processSvc,
		strategicObjectiveSvc:    strategicObjectiveSvc,
		slaRegisterSvc:           slaRegisterSvc,
	}
}

func (rrs *RiskRegisterService) GetAll(params *riskDTOs.RiskRegisterListQueryParamsDTO) (*riskDTOs.RiskRegisterListResponseDTO, error) {
	return rrs.repo.GetAll(params)
}

func (rrs RiskRegisterService) GetByID(id uint) (*models.RiskRegister, error) {
	risk, err := rrs.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return risk, nil
}

func (rrs RiskRegisterService) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*riskDTOs.RiskRegisterDetailDTO, error) {
	// Initialize your DTO
	var riskDTO riskDTOs.RiskRegisterDetailDTO

	risk, err := rrs.repo.GetFullDetailsByID(id, tx...)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	// Map the result to your DTO
	err = mapstructure.Decode(risk, &riskDTO)
	if err != nil {
		return nil, err
	}

	//todo :- also add extra fields which required in details screen

	return &riskDTO, nil
}

func (rrs *RiskRegisterService) Create(currentUserID uint, riskRegisterFormDTO *riskDTOs.RiskRegisterFormDTO) (uint, error) {
	tx := rrs.repo.BeginTransaction()

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

	riskRegister := &models.RiskRegister{
		ReferenceNumber:   referenceNumber,
		RiskCategoryID:    riskRegisterFormDTO.RiskCategoryID,
		RiskSourceID:      riskRegisterFormDTO.RiskSourceID,
		OwnerUnitID:       riskRegisterFormDTO.OwnerUnitID,
		OwnerID:           riskRegisterFormDTO.OwnerID,
		RiskStatusID:      riskRegisterFormDTO.RiskStatusID,
		Title:             riskRegisterFormDTO.Title,
		Description:       riskRegisterFormDTO.Description,
		ReportedAt:        dates.ReportedAt,
		ReviewFrequencyID: riskRegisterFormDTO.ReviewFrequencyID,
		NextReviewDate:    dates.NextReviewDate,
		Comments:          riskRegisterFormDTO.Comments,
		CreatedByID:       currentUserID,
	}

	id, err := rrs.repo.Create(riskRegister, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	// Save ComplianceRequirements
	var complianceRequirements []*models.RiskComplianceRequirement
	for _, requirementID := range riskRegisterFormDTO.ComplianceRequirementIDs {
		complianceRequirement, err := rrs.complianceRequirementSvc.GetByComplianceRequirementID(*requirementID)
		if err != nil {
			return 0, err
		}
		requirement := &models.RiskComplianceRequirement{
			RiskRegisterID:          id,
			ComplianceRegisterID:    complianceRequirement.ComplianceClause.ComplianceRegisterID,
			ComplianceRequirementID: *requirementID,
			CreatedByID:             currentUserID,
		}
		complianceRequirements = append(complianceRequirements, requirement)
	}

	if len(complianceRequirements) > 0 {
		if err := rrs.repo.CreateRiskComplianceRequirements(complianceRequirements, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
	}
	// Save Controls
	var controls []*models.RiskControl
	for _, controlID := range riskRegisterFormDTO.ControlIDs {
		control := &models.RiskControl{
			RiskRegisterID: id,
			ControlID:      *controlID,
			CreatedByID:    currentUserID,
		}
		controls = append(controls, control)
	}

	if len(controls) > 0 {
		if err := rrs.repo.CreateRiskControls(controls, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	// Save Processes
	var processes []*models.RiskProcess
	for _, processID := range riskRegisterFormDTO.ProcessIDs {
		process := &models.RiskProcess{
			RiskRegisterID: id,
			ProcessID:      *processID,
			CreatedByID:    currentUserID,
		}
		processes = append(processes, process)
	}

	if len(processes) > 0 {
		if err := rrs.repo.CreateRiskProcesses(processes, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
	}
	// Add Strategic Objectives
	var strategicObjectives []*models.RiskStrategicObjective
	for _, strategicObjectiveID := range riskRegisterFormDTO.StrategicObjectiveIDs {
		strategicObjective := &models.RiskStrategicObjective{
			RiskRegisterID:       id,
			StrategicObjectiveID: *strategicObjectiveID,
			CreatedByID:          currentUserID,
		}
		strategicObjectives = append(strategicObjectives, strategicObjective)
	}

	if len(strategicObjectives) > 0 {
		if err := rrs.repo.CreateRiskStrategicObjectives(strategicObjectives, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
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
		url := utility.GenerateDocumentURL("risk-register", strconv.FormatUint(uint64(riskRegister.ID), 10), encryptedFilename)

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

	tx.Commit()

	go func() {
		err := rrs.activityLogSvc.LogActivity(
			models.CreateAction,
			models.RiskRegister{}.TableName(),
			id,
			riskRegisterFormDTO,
			currentUserID,
			nil,
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
	}()

	return id, nil
}

func (rrs *RiskRegisterService) Update(id uint, currentUserID uint, riskRegisterFormDTO *riskDTOs.RiskRegisterFormDTO) error {
	tx := rrs.repo.BeginTransaction()

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

	existingRiskRegister.RiskCategoryID = riskRegisterFormDTO.RiskCategoryID
	existingRiskRegister.RiskSourceID = riskRegisterFormDTO.RiskSourceID
	existingRiskRegister.OwnerUnitID = riskRegisterFormDTO.OwnerUnitID
	existingRiskRegister.OwnerID = riskRegisterFormDTO.OwnerID
	existingRiskRegister.RiskStatusID = riskRegisterFormDTO.RiskStatusID
	existingRiskRegister.Title = riskRegisterFormDTO.Title
	existingRiskRegister.Description = riskRegisterFormDTO.Description
	existingRiskRegister.ReportedAt = dates.ReportedAt
	existingRiskRegister.ReviewFrequencyID = riskRegisterFormDTO.ReviewFrequencyID
	existingRiskRegister.NextReviewDate = dates.ReportedAt
	existingRiskRegister.Comments = riskRegisterFormDTO.Comments
	existingRiskRegister.UpdatedByID = &currentUserID

	if err := rrs.repo.Update(existingRiskRegister, tx); err != nil {
		tx.Rollback()
		return err
	}

	existingRisk, err := rrs.GetFullDetailsByID(id, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	//Delete all existing compliance requirements and controls and create again
	if len(existingRisk.ComplianceRequirements) > 0 {
		var deleteRequirementIds []*uint
		for _, req := range existingRisk.ComplianceRequirements {
			deleteRequirementIds = append(deleteRequirementIds, req.ID)
		}
		if err := rrs.repo.DeleteRiskComplianceRequirements(deleteRequirementIds, tx); err != nil {
			tx.Rollback()
			return err
		}

	}

	if len(existingRisk.Controls) > 0 {

		var deleteControlIDs []*uint
		for _, ctrl := range existingRisk.Controls {
			deleteControlIDs = append(deleteControlIDs, ctrl.ID)
		}
		if err := rrs.repo.DeleteRiskControls(deleteControlIDs, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	// Create ComplianceRequirements
	var complianceRequirements []*models.RiskComplianceRequirement
	for _, requirementID := range riskRegisterFormDTO.ComplianceRequirementIDs {
		requirement := &models.RiskComplianceRequirement{
			RiskRegisterID:          id,
			ComplianceRequirementID: *requirementID,
			CreatedByID:             currentUserID,
		}
		complianceRequirements = append(complianceRequirements, requirement)
	}

	if len(complianceRequirements) > 0 {
		if err := rrs.repo.CreateRiskComplianceRequirements(complianceRequirements, tx); err != nil {
			tx.Rollback()
			return err
		}
	}
	// Create Controls
	var controls []*models.RiskControl
	for _, controlID := range riskRegisterFormDTO.ControlIDs {
		control := &models.RiskControl{
			RiskRegisterID: id,
			ControlID:      *controlID,
			CreatedByID:    currentUserID,
		}
		controls = append(controls, control)
	}

	if len(controls) > 0 {
		if err := rrs.repo.CreateRiskControls(controls, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	existingProcessIDs := make(map[uint]bool)
	for _, proc := range existingRisk.Processes {
		existingProcessIDs[proc.ProcessID] = true
	}

	// Collect new process IDs
	newProcessIDs := make(map[uint]bool)
	for _, proc := range riskRegisterFormDTO.ProcessIDs {
		newProcessIDs[*proc] = true
	}

	// Add new processes
	var processesToAdd []*models.RiskProcess
	for _, newProcessID := range riskRegisterFormDTO.ProcessIDs {
		if _, exists := existingProcessIDs[*newProcessID]; !exists {
			// Create a new risk process record
			process := &models.RiskProcess{
				RiskRegisterID: id,
				ProcessID:      *newProcessID,
				CreatedByID:    currentUserID,
			}
			processesToAdd = append(processesToAdd, process)
		}
	}

	// Save the new processes
	if len(processesToAdd) > 0 {
		if err := rrs.repo.CreateRiskProcesses(processesToAdd, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	// Remove processes not in the new list
	var deleteProcessIDs []*uint

	for _, existingProcess := range existingRisk.Processes {
		if _, exists := newProcessIDs[existingProcess.ProcessID]; !exists {
			// Append process to delete list
			deleteProcessIDs = append(deleteProcessIDs, existingProcess.ID)
		}
	}

	// Delete processes not in the new list
	if len(deleteProcessIDs) > 0 {
		if err := rrs.repo.DeleteRiskProcesses(deleteProcessIDs, tx); err != nil {
			tx.Rollback()
			return err
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
	// Update Strategic Objectives
	existingStrategicObjectiveIDs := make(map[uint]bool)
	for _, obj := range existingRisk.StrategicObjectives {
		existingStrategicObjectiveIDs[obj.StrategicObjectiveID] = true
	}

	// Collect new strategic objective IDs
	newStrategicObjectiveIDs := make(map[uint]bool)
	for _, obj := range riskRegisterFormDTO.StrategicObjectiveIDs {
		newStrategicObjectiveIDs[*obj] = true
	}

	// Add new strategic objectives
	var strategicObjectivesToAdd []*models.RiskStrategicObjective
	for _, newStrategicObjectiveID := range riskRegisterFormDTO.StrategicObjectiveIDs {
		if _, exists := existingStrategicObjectiveIDs[*newStrategicObjectiveID]; !exists {
			// Create a new risk strategic objective record
			strategicObjective := &models.RiskStrategicObjective{
				RiskRegisterID:       id,
				StrategicObjectiveID: *newStrategicObjectiveID,
				CreatedByID:          currentUserID,
			}
			strategicObjectivesToAdd = append(strategicObjectivesToAdd, strategicObjective)
		}
	}

	// Save the new strategic objectives
	if len(strategicObjectivesToAdd) > 0 {
		if err := rrs.repo.CreateRiskStrategicObjectives(strategicObjectivesToAdd, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	// Remove strategic objectives not in the new list
	var deleteStrategicObjectiveIDs []*uint

	for _, existingStrategicObjective := range existingRisk.StrategicObjectives {
		if _, exists := newStrategicObjectiveIDs[existingStrategicObjective.StrategicObjectiveID]; !exists {
			// Append strategic objective to delete list
			deleteStrategicObjectiveIDs = append(deleteStrategicObjectiveIDs, existingStrategicObjective.ID)
		}
	}

	// Delete strategic objectives not in the new list
	if len(deleteStrategicObjectiveIDs) > 0 {
		if err := rrs.repo.DeleteRiskStrategicObjectives(deleteStrategicObjectiveIDs, tx); err != nil {
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

			url := utility.GenerateDocumentURL("risk-register", strconv.FormatUint(uint64(id), 10), encryptedFilename)
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
		err := rrs.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.RiskRegister{}.TableName(),
			id,
			riskRegisterFormDTO,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.UpdateAction,
				"tableName": models.RiskRegister{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return nil
}

func (rrs *RiskRegisterService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	tx := rrs.repo.BeginTransaction()

	var failedIds []uint
	var successfulIds []uint
	documentsToDelete := make(map[uint][]*riskDTOs.RiskDocumentDTO)

	for _, id := range ids {
		risk, err := rrs.GetFullDetailsByID(id, tx)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		// Check if controls exist
		if len(risk.Controls) > 0 {
			var deleteControlIDs []*uint
			for _, ctrl := range risk.Controls {
				deleteControlIDs = append(deleteControlIDs, ctrl.ID)
			}
			if err := rrs.repo.DeleteRiskControls(deleteControlIDs, tx); err != nil {
				tx.Rollback()
				return nil, err
			}
		}

		// Delete each compliance requirement if requirements exist
		if len(risk.ComplianceRequirements) > 0 {
			var deleteRequirementIds []*uint
			for _, req := range risk.ComplianceRequirements {
				deleteRequirementIds = append(deleteRequirementIds, req.ID)
			}
			if err := rrs.repo.DeleteRiskComplianceRequirements(deleteRequirementIds, tx); err != nil {
				tx.Rollback()
				return nil, err
			}

		}

		// Check if controls exist
		if len(risk.Processes) > 0 {
			var deleteProcessIDs []*uint
			for _, process := range risk.Processes {
				deleteProcessIDs = append(deleteProcessIDs, process.ID)
			}
			if err := rrs.repo.DeleteRiskProcesses(deleteProcessIDs, tx); err != nil {
				tx.Rollback()
				return nil, err
			}
		}
		if len(risk.SLARegisters) > 0 {
			var deleteSLARegisterIDs []*uint
			for _, sla := range risk.SLARegisters {
				deleteSLARegisterIDs = append(deleteSLARegisterIDs, sla.ID)
			}
			if err := rrs.repo.DeleteRiskSLAs(deleteSLARegisterIDs, tx); err != nil {
				tx.Rollback()
				return nil, err
			}
		}
		if len(risk.StrategicObjectives) > 0 {
			var deleteStrategicObjectiveIds []*uint
			for _, obj := range risk.StrategicObjectives {
				deleteStrategicObjectiveIds = append(deleteStrategicObjectiveIds, obj.ID)
			}
			if err := rrs.repo.DeleteRiskStrategicObjectives(deleteStrategicObjectiveIds, tx); err != nil {
				failedIds = append(failedIds, id)
				continue
			}
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
				tx.Rollback()
				return nil, err
			}

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

func (raps *RiskRegisterService) UpdateResponseStrategy(riskID uint, currentUserID uint, riskResponseStrategyFormDTO *riskDTOs.RiskResponseStrategyFormDTO) error {
	tx := raps.repo.BeginTransaction()

	existingRisk, err := raps.GetByID(riskID)
	if err != nil {
		tx.Rollback()
		return err
	}

	existingRisk.RiskResponseStrategyID = &riskResponseStrategyFormDTO.RiskResponseStrategyID

	if err := raps.repo.Update(existingRisk, tx); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil
}

func (raps *RiskRegisterService) UpdateRiskRating(riskID uint, currentUserID uint, riskEvaluationDetails *models.RiskEvaluation, tx utility.DatabaseTransaction) error {

	existingRisk, err := raps.GetByID(riskID)
	if err != nil {
		tx.Rollback()
		return err
	}

	existingRisk.CurrentImpactID = &riskEvaluationDetails.RiskMatrixImpactID
	existingRisk.CurrentLikelihoodID = &riskEvaluationDetails.RiskMatrixLikelihoodID
	existingRisk.CurrentRiskRatingID = &riskEvaluationDetails.RiskMatrixRiskRatingID
	existingRisk.CurrentRiskRatingValueID = &riskEvaluationDetails.RiskMatrixRiskRatingValueID

	if err := raps.repo.Update(existingRisk, tx); err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (rrs *RiskRegisterService) validateRiskRegisterFormDTO(riskRegisterFormDTO *riskDTOs.RiskRegisterFormDTO) error {
	if _, err := rrs.riskCategorySvc.GetByID(riskRegisterFormDTO.RiskCategoryID); err != nil {
		return utility.NewValidationError(fmt.Errorf("risk category does not exist"), "not_exist", "risk category")
	}

	if _, err := rrs.riskSourceSvc.GetByID(riskRegisterFormDTO.RiskSourceID); err != nil {
		return utility.NewValidationError(fmt.Errorf("risk source does not exist"), "not_exist", "risk source")
	}

	//@todo org unit validation

	if _, err := rrs.userSvc.GetByID(riskRegisterFormDTO.OwnerID); err != nil {
		return utility.NewValidationError(fmt.Errorf("owner does not exist"), "not_exist", "owner")
	}

	if _, err := rrs.riskStatusSvc.GetByID(riskRegisterFormDTO.RiskStatusID); err != nil {
		return utility.NewValidationError(fmt.Errorf("risk status does not exist"), "not_exist", "risk status")
	}

	if riskRegisterFormDTO.ReviewFrequencyID != nil {
		if _, err := rrs.reviewFrequencySvc.GetByID(*riskRegisterFormDTO.ReviewFrequencyID); err != nil {
			return utility.NewValidationError(fmt.Errorf("review frequency does not exist"), "not_exist", "review frequency")
		}
	}

	controlIds := riskRegisterFormDTO.ControlIDs

	controlIDs := make(map[uint]bool)

	for _, controlId := range controlIds {
		if controlIDs[*controlId] {
			return utility.NewValidationError(fmt.Errorf("duplicate control ID: %d", controlId), "duplicate", "control")
		}

		if _, err := rrs.controlSvc.GetByID(*controlId); err != nil {
			return utility.NewValidationError(fmt.Errorf("control does not exist for ID: %d", controlId), "not_exist", "control")
		}

		controlIDs[*controlId] = true
	}

	requirementIds := riskRegisterFormDTO.ComplianceRequirementIDs

	requirementIDs := make(map[uint]bool)

	for _, requirementId := range requirementIds {
		if requirementIDs[*requirementId] {
			return utility.NewValidationError(fmt.Errorf("duplicate compliance requirement ID: %d", requirementId), "duplicate", "compliance requirement")
		}

		if _, err := rrs.complianceRequirementSvc.GetByComplianceRequirementID(*requirementId); err != nil {
			return utility.NewValidationError(fmt.Errorf("compliance requirement does not exist for ID: %d", requirementId), "not_exist", "compliance requirement")
		}

		requirementIDs[*requirementId] = true
	}

	processIds := riskRegisterFormDTO.ProcessIDs

	processIDs := make(map[uint]bool)

	for _, processId := range processIds {
		if processIDs[*processId] {
			return utility.NewValidationError(fmt.Errorf("duplicate process ID: %d", processId), "duplicate", "process")
		}

		if _, err := rrs.processSvc.GetByID(*processId); err != nil {
			return utility.NewValidationError(fmt.Errorf("process does not exist for ID: %d", processId), "not_exist", "process")
		}

		processIDs[*processId] = true
	}

	slaRegisterIds := riskRegisterFormDTO.ProcessIDs

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

	strategicObjectiveIds := riskRegisterFormDTO.StrategicObjectiveIDs

	strategicObjectiveIDs := make(map[uint]bool)

	for _, strategicObjectiveId := range strategicObjectiveIds {
		if strategicObjectiveIDs[*strategicObjectiveId] {
			return utility.NewValidationError(fmt.Errorf("duplicate strategic objective ID: %d", strategicObjectiveId), "duplicate", "strategic objective")
		}

		if _, err := rrs.strategicObjectiveSvc.GetByID(*strategicObjectiveId); err != nil {
			return utility.NewValidationError(fmt.Errorf("strategic objective does not exist for ID: %d", strategicObjectiveId), "not_exist", "strategic objective")
		}

		strategicObjectiveIDs[*strategicObjectiveId] = true
	}

	return nil
}

func (rrs *RiskRegisterService) parseRiskRegisterDates(dto *riskDTOs.RiskRegisterFormDTO) (*riskDTOs.RiskRegisterParsedDates, error) {
	parsed := &riskDTOs.RiskRegisterParsedDates{}

	var err error
	if *dto.ReportedAt != "" {
		parsed.ReportedAt, err = utility.ParseDate(*dto.ReportedAt)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid reported at"), "invalid_date", "reported at")
		}
	}

	if *dto.NextReviewDate != "" {
		parsed.NextReviewDate, err = utility.ParseDate(*dto.NextReviewDate)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid next review date"), "invalid_date", "next review date")
		}
	}

	return parsed, nil
}
