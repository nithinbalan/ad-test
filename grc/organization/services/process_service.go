package services

import (
	"errors"
	"fmt"
	commonDTOs "jethur/grc/common/dtos"
	commonServices "jethur/grc/common/services"
	logServices "jethur/grc/logs/services"
	masterServices "jethur/grc/masters/services"
	"jethur/grc/models"
	organizationDTOs "jethur/grc/organization/dtos"
	organizationRepositories "jethur/grc/organization/repositories"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type ProcessServiceInterface interface {
	GetAll(params *organizationDTOs.ProcessListQueryParamsDTO) (*organizationDTOs.ProcessListResponseDTO, error)
	GetByID(id uint) (*models.Process, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*organizationDTOs.ProcessDetailDTO, error)
	Create(currentUserID uint, processFormDTO *organizationDTOs.ProcessFormDTO) (uint, error)
	Update(id uint, currentUserID uint, processFormDTO *organizationDTOs.ProcessFormDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)

	Archive(ids []uint, currentUserID uint) ([]uint, error)
	Restore(ids []uint, currentUserID uint) ([]uint, error)

	GetControlSummary(id uint) (map[string]interface{}, error)
}

type ProcessService struct {
	repo               organizationRepositories.ProcessRepositoryInterface
	activityLogSvc     logServices.ActivityLogServiceInterface
	processGroupSvc    ProcessGroupServiceInterface
	userSvc            UserServiceInterface
	lifecycleStatusSvc masterServices.LifecycleStatusServiceInterface
	controlSvc         ControlServiceInterface
	notificationSvc    commonServices.NotificationServiceInterface
}

func NewProcessService(
	repo organizationRepositories.ProcessRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	processGroupSvc ProcessGroupServiceInterface,
	userSvc UserServiceInterface,
	lifecycleStatusSvc masterServices.LifecycleStatusServiceInterface,
	controlSvc ControlServiceInterface,
	notificationSvc commonServices.NotificationServiceInterface) ProcessServiceInterface {
	return &ProcessService{
		repo:               repo,
		activityLogSvc:     activityLogSvc,
		processGroupSvc:    processGroupSvc,
		userSvc:            userSvc,
		lifecycleStatusSvc: lifecycleStatusSvc,
		controlSvc:         controlSvc,
		notificationSvc:    notificationSvc,
	}
}

func (ps *ProcessService) GetAll(params *organizationDTOs.ProcessListQueryParamsDTO) (*organizationDTOs.ProcessListResponseDTO, error) {
	return ps.repo.GetAll(params)
}

func (ps *ProcessService) GetByID(id uint) (*models.Process, error) {
	process, err := ps.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return process, nil
}

func (ps *ProcessService) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*organizationDTOs.ProcessDetailDTO, error) {
	var processDetailDTO organizationDTOs.ProcessDetailDTO

	process, err := ps.repo.GetFullDetailsByID(id, tx...)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	err = mapstructure.Decode(process, &processDetailDTO)
	if err != nil {
		return nil, err
	}

	//todo :- also add extra fields which required in details screen

	return &processDetailDTO, nil
}

func (ps *ProcessService) Create(currentUserID uint, processFormDTO *organizationDTOs.ProcessFormDTO) (uint, error) {
	tx := ps.repo.BeginTransaction()

	if err := ps.validateProcessFormDTO(processFormDTO); err != nil {
		return 0, err
	}

	referenceNumber, err := ps.repo.GenerateReferenceNumber()
	if err != nil {
		return 0, err
	}

	lifecycleStatus, err := ps.lifecycleStatusSvc.GetByType("active")
	if err != nil {
		return 0, err
	}
	organizationUnitID, err := ps.userSvc.GetUserPrimaryUnitID(uint(processFormDTO.OwnerID))
	if err != nil {
		return 0, utility.NewValidationError(fmt.Errorf("organization unit not set for owner"), "not_exist", "process owner unit")
	}

	process := &models.Process{
		ProcessGroupID:    processFormDTO.ProcessGroupID,
		OwnerID:           processFormDTO.OwnerID,
		OwnerUnitID:       organizationUnitID,
		LifecycleStatusID: lifecycleStatus.ID,
		ReferenceNumber:   referenceNumber,
		TitleEn:           processFormDTO.TitleEn,
		TitleAr:           processFormDTO.TitleAr,
		Description:       processFormDTO.Description,
		CreatedByID:       currentUserID,
	}

	id, err := ps.repo.Create(process)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	var controls []*models.ProcessControl
	for _, controlID := range processFormDTO.ControlIds {
		control := &models.ProcessControl{
			ProcessID:   id,
			ControlID:   *controlID,
			CreatedByID: currentUserID,
		}
		controls = append(controls, control)
	}

	if len(controls) > 0 {
		if err := ps.repo.CreateProcessControls(controls, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	var documents []*models.ProcessAssociatedDocument
	for _, documentID := range processFormDTO.DocumentIds {
		document := &models.ProcessAssociatedDocument{
			ProcessID:   id,
			DocumentID:  *documentID,
			CreatedByID: currentUserID,
		}
		documents = append(documents, document)
	}

	if len(documents) > 0 {
		if err := ps.repo.CreateProcessDocuments(documents, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	urlParam := fmt.Sprintf(`{"id":%d}`, id)
	queryParam := ""

	dbNotificationDTO := &commonDTOs.DBNotificationDTO{
		UserID:               processFormDTO.OwnerID,
		ModuleIdentifier:     "organization",
		NotificationType:     string(commonDTOs.AssignmentType),
		ReferenceNumber:      referenceNumber,
		HeadingEn:            "Process Assignment",
		HeadingAr:            "مهمة العملية",
		TitleEn:              processFormDTO.TitleEn,
		TitleAr:              processFormDTO.TitleAr,
		URLParam:             urlParam,
		QueryParam:           queryParam,
		IsSystemNotification: false,
		CreatedByID:          &currentUserID,
	}

	//email notification
	resposibleUser, err := ps.userSvc.GetByID(processFormDTO.OwnerID)
	if err != nil {
		return 0, err
	}

	redirectURL := fmt.Sprintf("process/%d", id)
	emailTo := []string{resposibleUser.Email}
	subject := fmt.Sprintf("Process/%s: Added", process.TitleEn)
	templateRelativePath := "organization/templates/email/create_process_mail.html"

	emailData := organizationDTOs.SendCreateProcessMailDTO{
		FirstNameEn: resposibleUser.FirstNameEn,
		LastNameEn:  resposibleUser.LastNameEn,
		MessageEn:   process.TitleEn,
		RedirectURL: redirectURL,
	}

	emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
		EmailTo:              emailTo,
		Subject:              subject,
		TemplateRelativePath: templateRelativePath,
		EmailData:            emailData,
	}

	tx.Commit()

	go func() {
		err := ps.notificationSvc.SendNotification(
			dbNotificationDTO, emailNotificationDTO, "process")
		if err != nil {
			utility.LogError("Failed to send notifications", map[string]interface{}{
				"error": err.Error(),
			})
		}
	}()

	go func() {
		err := ps.activityLogSvc.LogActivity(
			models.CreateAction,
			models.Process{}.TableName(),
			id,
			processFormDTO,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.Control{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return id, nil
}

func (ps *ProcessService) Update(id uint, currentUserID uint, processFormDTO *organizationDTOs.ProcessFormDTO) error {
	tx := ps.repo.BeginTransaction()

	process, err := ps.GetByID(id)

	if err != nil {
		return err
	}

	if err := ps.validateProcessFormDTO(processFormDTO); err != nil {
		return err
	}
	if process.OwnerID != processFormDTO.OwnerID {
		organizationUnitID, err := ps.userSvc.GetUserPrimaryUnitID(uint(processFormDTO.OwnerID))
		if err != nil {
			return utility.NewValidationError(fmt.Errorf("organization unit not set for owner"), "not_exist", "process owner unit")
		}
		process.OwnerUnitID = organizationUnitID
	}
	process.ProcessGroupID = processFormDTO.ProcessGroupID
	process.OwnerID = processFormDTO.OwnerID
	process.TitleEn = processFormDTO.TitleEn
	process.TitleAr = processFormDTO.TitleAr
	process.Description = processFormDTO.Description
	process.UpdatedByID = &currentUserID

	err = ps.repo.Update(process)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Collect existing control IDs
	existingProcess, err := ps.GetFullDetailsByID(id, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	existingControlIDs := make(map[uint]bool)
	for _, ctrl := range existingProcess.Controls {
		existingControlIDs[ctrl.ControlID] = true
	}

	// Collect new control IDs
	newControlIDs := make(map[uint]bool)
	for _, ctrl := range processFormDTO.ControlIds {
		newControlIDs[*ctrl] = true
	}

	// Add new controls
	var controlsToAdd []*models.ProcessControl
	for _, newControlID := range processFormDTO.ControlIds {
		if _, exists := existingControlIDs[*newControlID]; !exists {
			// Create a new process control record
			control := &models.ProcessControl{
				ProcessID:   id,
				ControlID:   *newControlID,
				CreatedByID: currentUserID,
			}
			controlsToAdd = append(controlsToAdd, control)
		}
	}

	// Save the new controls
	if len(controlsToAdd) > 0 {
		if err := ps.repo.CreateProcessControls(controlsToAdd, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	// Remove controls not in the new list
	var deleteControlIDs []*uint

	for _, existingControl := range existingProcess.Controls {
		if _, exists := newControlIDs[existingControl.ControlID]; !exists {
			// Append control to delete list
			deleteControlIDs = append(deleteControlIDs, existingControl.ID)
		}
	}

	// Delete controls not in the new list
	if len(deleteControlIDs) > 0 {
		if err := ps.repo.DeleteProcessControls(deleteControlIDs, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	existingDocumentIDs := make(map[uint]bool)
	for _, doc := range existingProcess.Documents {
		existingDocumentIDs[doc.DocumentID] = true
	}

	// Collect new document IDs
	newDocumentIDs := make(map[uint]bool)
	for _, doc := range processFormDTO.DocumentIds {
		newDocumentIDs[*doc] = true
	}

	// Add new documents
	var documentsToAdd []*models.ProcessAssociatedDocument
	for _, newDocumentID := range processFormDTO.DocumentIds {
		if _, exists := existingDocumentIDs[*newDocumentID]; !exists {
			// Create a new process document record
			document := &models.ProcessAssociatedDocument{
				ProcessID:   id,
				DocumentID:  *newDocumentID,
				CreatedByID: currentUserID,
			}
			documentsToAdd = append(documentsToAdd, document)
		}
	}

	// Save the new documents
	if len(documentsToAdd) > 0 {
		if err := ps.repo.CreateProcessDocuments(documentsToAdd, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	// Remove documents not in the new list
	var deleteDocumentIDs []*uint

	for _, existingDocument := range existingProcess.Documents {
		if _, exists := newDocumentIDs[existingDocument.DocumentID]; !exists {
			// Append document to delete list
			deleteDocumentIDs = append(deleteDocumentIDs, existingDocument.ID)
		}
	}

	// Delete documents not in the new list
	if len(deleteDocumentIDs) > 0 {
		if err := ps.repo.DeleteProcessDocuments(deleteDocumentIDs, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()
	return nil
}

func (ps *ProcessService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	tx := ps.repo.BeginTransaction()

	var failedIds []uint
	var successfulIds []uint

	for _, id := range ids {
		process, err := ps.GetFullDetailsByID(id, tx)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		if len(process.Controls) > 0 {
			var deleteControlIds []*uint
			for _, ctrl := range process.Controls {
				deleteControlIds = append(deleteControlIds, ctrl.ID)
			}
			if err := ps.repo.DeleteProcessControls(deleteControlIds, tx); err != nil {
				tx.Rollback()
				return nil, err
			}
		}

		if len(process.Documents) > 0 {
			var deleteDocumentIds []*uint
			for _, document := range process.Documents {
				deleteDocumentIds = append(deleteDocumentIds, document.ID)
			}
			if err := ps.repo.DeleteProcessDocuments(deleteDocumentIds, tx); err != nil {
				tx.Rollback()
				return nil, err
			}
		}

		if err := ps.repo.Delete(id, tx); err != nil {
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

	tx.Commit()

	if len(successfulIds) > 0 {
		go func() {
			err := ps.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.Process{}.TableName(),
				successfulIds,
				nil,
				currentUserID,
				nil,
			)
			if err != nil {
				fmt.Printf("Failed to log bulk delete activity: %v\n", err)
			}
		}()
	}

	return failedIds, nil
}

func (ps *ProcessService) Archive(ids []uint, currentUserID uint) ([]uint, error) {
	lifecycleStatus, err := ps.lifecycleStatusSvc.GetByType("archived")
	if err != nil {
		return nil, err
	}

	idsToUpdate, err := ps.repo.GetIDsByLifecycleStatus(ids, "active")
	if err != nil {
		return nil, err
	}

	successfulIds, err := ps.repo.UpdateLifecycleStatusForIDs(idsToUpdate, lifecycleStatus.ID)

	if err != nil {
		return nil, err
	}

	if len(successfulIds) > 0 {
		go func() {
			err := ps.activityLogSvc.LogBulkActivities(
				models.ArchiveAction,
				models.Process{}.TableName(),
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

func (ps *ProcessService) Restore(ids []uint, currentUserID uint) ([]uint, error) {
	lifecycleStatus, err := ps.lifecycleStatusSvc.GetByType("active")
	if err != nil {
		return nil, err
	}

	idsToUpdate, err := ps.repo.GetIDsByLifecycleStatus(ids, "archived")
	if err != nil {
		return nil, err
	}

	successfulIds, err := ps.repo.UpdateLifecycleStatusForIDs(idsToUpdate, lifecycleStatus.ID)

	if err != nil {
		return nil, err
	}

	if len(successfulIds) > 0 {
		go func() {
			err := ps.activityLogSvc.LogBulkActivities(
				models.RestoreAction,
				models.Process{}.TableName(),
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

func (ps *ProcessService) validateProcessFormDTO(processFormDTO *organizationDTOs.ProcessFormDTO) error {
	if _, err := ps.processGroupSvc.GetByID(uint(processFormDTO.ProcessGroupID)); err != nil {
		return utility.NewValidationError(fmt.Errorf("process group does not exist"), "not_exist", "process group")
	}

	if _, err := ps.userSvc.GetByID(uint(processFormDTO.OwnerID)); err != nil {
		return utility.NewValidationError(fmt.Errorf("process owner does not exist"), "not_exist", "process owner")
	}

	controls := processFormDTO.ControlIds

	controlIDs := make(map[uint]bool)

	for _, control := range controls {
		if controlIDs[*control] {
			return utility.NewValidationError(fmt.Errorf("duplicate control ID: %d", control), "duplicate", "control")
		}

		if _, err := ps.controlSvc.GetByID(*control); err != nil {
			return utility.NewValidationError(fmt.Errorf("control does not exist for ID: %d", control), "not_exist", "control")
		}

		controlIDs[*control] = true
	}

	documents := processFormDTO.DocumentIds

	documentIDs := make(map[uint]bool)

	for _, document := range documents {
		if documentIDs[*document] {
			return utility.NewValidationError(fmt.Errorf("duplicate document ID: %d", document), "duplicate", "document")
		}

		if _, err := ps.repo.GetDocumentByID(*document); err != nil {
			return utility.NewValidationError(fmt.Errorf("document does not exist for ID: %d", document), "not_exist", "document")
		}

		documentIDs[*document] = true
	}

	return nil
}

func (ps *ProcessService) GetControlSummary(id uint) (map[string]interface{}, error) {
	processControls, err := ps.repo.GetControlSummary(id)
	if err != nil {
		return nil, err
	}

	if len(processControls) == 0 {
		return map[string]interface{}{
			"total_control_score":       0,
			"implementation_status":     []map[string]interface{}{},
			"efficiency_percentage":     0,
			"implementation_percentage": 0,
			"total_efficiency":          []map[string]interface{}{},
			"total_control_count":       0,
		}, nil
	}

	implementationStatusCount := make(map[string]int)
	efficiencyStatusCount := make(map[string]int)
	implementationDetails := make(map[string]models.ControlImplementationStatus)
	efficiencyDetails := make(map[string]models.ControlEfficiencyStatus)

	for _, pc := range processControls {
		implStatus := pc.Control.ControlImplementationStatus.Type
		implementationStatusCount[implStatus]++
		implementationDetails[implStatus] = *pc.Control.ControlImplementationStatus

		if pc.Control.ControlEfficiencyStatus != nil {
			effStatus := pc.Control.ControlEfficiencyStatus.Type
			efficiencyStatusCount[effStatus]++
			efficiencyDetails[effStatus] = *pc.Control.ControlEfficiencyStatus
		}
	}

	totalControls := len(processControls)

	implementationStatus := make([]map[string]interface{}, 0)
	for status, count := range implementationStatusCount {
		implDetail := implementationDetails[status]
		implementationStatus = append(implementationStatus, map[string]interface{}{
			"id":         implDetail.ID,
			"type":       implDetail.Type,
			"title_en":   implDetail.TitleEn,
			"title_ar":   implDetail.TitleAr,
			"count":      count,
			"color_code": implDetail.ColorCode,
		})
	}

	efficiencyStatus := make([]map[string]interface{}, 0)
	for status, count := range efficiencyStatusCount {
		effDetail := efficiencyDetails[status]
		efficiencyStatus = append(efficiencyStatus, map[string]interface{}{
			"id":         effDetail.ID,
			"type":       effDetail.Type,
			"title_en":   effDetail.TitleEn,
			"title_ar":   effDetail.TitleAr,
			"count":      count,
			"color_code": effDetail.ColorCode,
		})
	}

	summary := map[string]interface{}{
		"total_control_score":       (float64(implementationStatusCount["implemented"]) / float64(totalControls)) * 100,
		"implementation_status":     implementationStatus,
		"efficiency_percentage":     (float64(efficiencyStatusCount["efficient"]) / float64(totalControls)) * 100,
		"implementation_percentage": (float64(implementationStatusCount["implemented"]) / float64(totalControls)) * 100,
		"total_efficiency":          efficiencyStatus,
		"total_control_count":       totalControls,
	}

	return summary, nil
}
