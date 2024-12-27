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
	"mime/multipart"
	"os"
	"time"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type ProcessServiceInterface interface {
	GetAll(params *organizationDTOs.ProcessListQueryParamsDTO) (*organizationDTOs.ProcessListResponseDTO, error)
	GetByID(id uint, tx ...utility.DatabaseTransaction) (*models.Process, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*organizationDTOs.ProcessDetailDTO, error)
	Create(currentUserID uint, processFormDTO *organizationDTOs.ProcessFormDTO) (uint, error)
	Update(id uint, currentUserID uint, processFormDTO *organizationDTOs.ProcessFormDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)

	Archive(ids []uint, currentUserID uint) ([]uint, error)
	Restore(ids []uint, currentUserID uint) ([]uint, error)

	GetControlSummary(id uint) (map[string]interface{}, error)

	ExportProcesses(data []*organizationDTOs.ProcessListDTO, currentUserID uint) (*os.File, error)
	GenerateProcessTemplate() (*os.File, error)
	ImportProcesses(file multipart.File, currentUserID uint) ([]int, error)

	UpdateProcessRiskRating(process *models.Process, tx utility.DatabaseTransaction) error
	UpdateProcessRecoveryObjective(process *models.Process, tx utility.DatabaseTransaction) error
}

type ProcessService struct {
	repo               organizationRepositories.ProcessRepositoryInterface
	activityLogSvc     logServices.ActivityLogServiceInterface
	processGroupSvc    ProcessGroupServiceInterface
	userSvc            UserServiceInterface
	lifecycleStatusSvc masterServices.LifecycleStatusServiceInterface
	controlSvc         ControlServiceInterface
	notificationSvc    commonServices.NotificationServiceInterface
	userSettingSvc     UserSettingServiceInterface
}

func NewProcessService(
	repo organizationRepositories.ProcessRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	processGroupSvc ProcessGroupServiceInterface,
	userSvc UserServiceInterface,
	lifecycleStatusSvc masterServices.LifecycleStatusServiceInterface,
	controlSvc ControlServiceInterface,
	notificationSvc commonServices.NotificationServiceInterface,
	userSettingSvc UserSettingServiceInterface) ProcessServiceInterface {
	return &ProcessService{
		repo:               repo,
		activityLogSvc:     activityLogSvc,
		processGroupSvc:    processGroupSvc,
		userSvc:            userSvc,
		lifecycleStatusSvc: lifecycleStatusSvc,
		controlSvc:         controlSvc,
		notificationSvc:    notificationSvc,
		userSettingSvc:     userSettingSvc,
	}
}

func (ps *ProcessService) GetAll(params *organizationDTOs.ProcessListQueryParamsDTO) (*organizationDTOs.ProcessListResponseDTO, error) {
	return ps.repo.GetAll(params)
}

func (ps *ProcessService) GetByID(id uint, tx ...utility.DatabaseTransaction) (*models.Process, error) {
	process, err := ps.repo.GetByID(id, tx...)

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

	processDetailDTO.CreatedAt = process.CreatedAt
	processDetailDTO.UpdatedAt = process.UpdatedAt
	// Map the latest auditable item
	if len(process.AuditableItems) > 0 {
		latestAuditableItem := process.AuditableItems[0]
		var auditableItemDTO organizationDTOs.AuditableItemMinimalDetailDTO
		err := mapstructure.Decode(latestAuditableItem, &auditableItemDTO)
		if err != nil {
			return nil, err
		}
		processDetailDTO.AuditableItem = &auditableItemDTO
		processDetailDTO.AuditableItem.CreatedAt = latestAuditableItem.CreatedAt
		processDetailDTO.AuditableItem.UpdatedAt = latestAuditableItem.UpdatedAt
	}

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
		CreatedBYID:       currentUserID,
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

	var assetRegister []*models.ProcessAsset
	for _, assetRegisterID := range processFormDTO.AssetRegisterIds {
		asset := &models.ProcessAsset{
			ProcessID:       id,
			AssetRegisterID: *assetRegisterID,
			CreatedByID:     currentUserID,
		}
		assetRegister = append(assetRegister, asset)
	}

	if len(assetRegister) > 0 {
		if err := ps.repo.CreateProcessAssetRegister(assetRegister, tx); err != nil {
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

	tx.Commit()

	go func() {
		existingProcess, err := ps.GetFullDetailsByID(process.ID)
		if err != nil {
			utility.LogError("Failed to fetch responsible user", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		//log
		err = ps.activityLogSvc.LogActivity(
			models.CreateAction,
			models.Process{}.TableName(),
			id,
			existingProcess,
			currentUserID,
			nil,
		)

		if err != nil {
			utility.LogError("Failed to log activity", map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.Process{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			})
		}

		//notification
		urlParam := fmt.Sprintf(`{"id":%d}`, id)
		entity := "process"
		action := "responsible-user-assign"

		dbNotificationDTO := &commonDTOs.DBNotificationDTO{
			UserID:               processFormDTO.OwnerID,
			ModuleIdentifier:     "organization",
			NotificationType:     string(commonDTOs.AssignmentType),
			ReferenceNumber:      process.ReferenceNumber,
			HeadingEn:            "Process Responsiblity Assigned",
			HeadingAr:            "مسؤولية العملية المعينة",
			TitleEn:              processFormDTO.TitleEn,
			TitleAr:              processFormDTO.TitleAr,
			Entity:               &entity,
			Action:               &action,
			URLParam:             urlParam,
			IsSystemNotification: false,
			CreatedByID:          &currentUserID,
		}

		responsibleUser, err := ps.userSvc.GetByID(processFormDTO.OwnerID)
		if err != nil {
			utility.LogError("Failed to fetch responsible user", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		emailTo := []string{responsibleUser.Email}
		subject := fmt.Sprintf("Process/%s: Ownership Assigned", processFormDTO.TitleEn)
		appURL := utility.GetAppURL()
		emailData := map[string]interface{}{
			"FirstNameEn":            responsibleUser.FirstNameEn,
			"LastNameEn":             responsibleUser.LastNameEn,
			"ProcessReferenceNumber": existingProcess.ReferenceNumber,
			"ProcessName":            existingProcess.TitleEn,
			"ProcessGroup":           existingProcess.ProcessGroup.TitleEn,
			"AssignedByFirstNameEn":  existingProcess.CreatedByUser.FirstNameEN,
			"AssignedByLastNameEn":   existingProcess.CreatedByUser.LastNameEN,
			"AssignedOn":             existingProcess.CreatedAt.Format("2006-01-02"),
			"RedirectURL":            fmt.Sprintf("%s/process/%d", appURL, id),
		}

		emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
			EmailTo:              emailTo,
			Subject:              subject,
			TemplateRelativePath: "organization/templates/email/process_create_owner_mail.html",
			EmailData:            emailData,
		}

		// Log email before sending (status: "Pending")
		err = ps.notificationSvc.LogEmail(
			emailNotificationDTO.EmailTo[0],
			emailNotificationDTO.Subject,
			emailNotificationDTO.TemplateRelativePath,
			emailNotificationDTO.EmailData,
			"Pending",
			"",
			nil,
			currentUserID,
		)
		if err != nil {
			utility.LogError("Failed to log email", map[string]interface{}{"error": err.Error()})
		}

		err = ps.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "process")
		if err != nil {
			// Log the email as "Failed"
			ps.notificationSvc.LogEmail(
				emailNotificationDTO.EmailTo[0],
				emailNotificationDTO.Subject,
				emailNotificationDTO.TemplateRelativePath,
				emailNotificationDTO.EmailData,
				"Failed",
				err.Error(),
				nil,
				currentUserID,
			)
			utility.LogError("Failed to send notifications", map[string]interface{}{
				"error": err.Error(),
			})
		} else {
			now := time.Now()
			// Log the email as "Sent"
			ps.notificationSvc.LogEmail(
				emailNotificationDTO.EmailTo[0],
				emailNotificationDTO.Subject,
				emailNotificationDTO.TemplateRelativePath,
				emailNotificationDTO.EmailData,
				"Sent",
				"",
				&now,
				currentUserID,
			)
		}
	}()

	return id, nil
}

func (ps *ProcessService) Update(id uint, currentUserID uint, processFormDTO *organizationDTOs.ProcessFormDTO) error {
	tx := ps.repo.BeginTransaction()

	oldProcess, err := ps.GetFullDetailsByID(id, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	process, err := ps.GetByID(id)
	if err != nil {
		tx.Rollback()
		return err
	}

	oldOwnerID := process.OwnerID

	if err := ps.validateProcessFormDTO(processFormDTO); err != nil {
		tx.Rollback()
		return err
	}

	// If owner has changed, update the OwnerUnitID
	if process.OwnerID != processFormDTO.OwnerID {
		organizationUnitID, err := ps.userSvc.GetUserPrimaryUnitID(uint(processFormDTO.OwnerID))
		if err != nil {
			tx.Rollback()
			return utility.NewValidationError(fmt.Errorf("organization unit not set for owner"), "not_exist", "process owner unit")
		}
		process.OwnerUnitID = organizationUnitID
	}

	process.ProcessGroupID = processFormDTO.ProcessGroupID
	process.OwnerID = processFormDTO.OwnerID
	process.TitleEn = processFormDTO.TitleEn
	process.TitleAr = processFormDTO.TitleAr
	process.Description = processFormDTO.Description
	process.UpdatedBYID = &currentUserID

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

	newControlIDs := make(map[uint]bool)
	for _, ctrl := range processFormDTO.ControlIds {
		newControlIDs[*ctrl] = true
	}

	var controlsToAdd []*models.ProcessControl
	for _, newControlID := range processFormDTO.ControlIds {
		if _, exists := existingControlIDs[*newControlID]; !exists {
			control := &models.ProcessControl{
				ProcessID:   id,
				ControlID:   *newControlID,
				CreatedByID: currentUserID,
			}
			controlsToAdd = append(controlsToAdd, control)
		}
	}

	if len(controlsToAdd) > 0 {
		if err := ps.repo.CreateProcessControls(controlsToAdd, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	var deleteControlIDs []*uint
	for _, existingControl := range existingProcess.Controls {
		if _, exists := newControlIDs[existingControl.ControlID]; !exists {
			deleteControlIDs = append(deleteControlIDs, existingControl.ID)
		}
	}

	if len(deleteControlIDs) > 0 {
		if err := ps.repo.DeleteProcessControls(deleteControlIDs, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	existingAssetRegisterIDs := make(map[uint]bool)
	for _, ctrl := range existingProcess.Assets {
		existingAssetRegisterIDs[ctrl.AssetRegisterID] = true
	}

	newAssetRegisterIDs := make(map[uint]bool)
	for _, ctrl := range processFormDTO.AssetRegisterIds {
		newAssetRegisterIDs[*ctrl] = true
	}

	var assetRegisterToAdd []*models.ProcessAsset
	for _, newAssetRegisterID := range processFormDTO.AssetRegisterIds {
		if _, exists := existingAssetRegisterIDs[*newAssetRegisterID]; !exists {
			assetRegister := &models.ProcessAsset{
				ProcessID:       id,
				AssetRegisterID: *newAssetRegisterID,
				CreatedByID:     currentUserID,
			}
			assetRegisterToAdd = append(assetRegisterToAdd, assetRegister)
		}
	}

	if len(assetRegisterToAdd) > 0 {
		if err := ps.repo.CreateProcessAssetRegister(assetRegisterToAdd, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	var deleteAssetRegisterIDs []*uint
	for _, existingAssetRegister := range existingProcess.Assets {
		if _, exists := newAssetRegisterIDs[existingAssetRegister.AssetRegisterID]; !exists {
			deleteAssetRegisterIDs = append(deleteAssetRegisterIDs, existingAssetRegister.ID)
		}
	}

	if len(deleteAssetRegisterIDs) > 0 {
		if err := ps.repo.DeleteProcessAssetRegister(deleteAssetRegisterIDs, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	existingDocumentIDs := make(map[uint]bool)
	for _, doc := range existingProcess.Documents {
		existingDocumentIDs[doc.DocumentID] = true
	}

	newDocumentIDs := make(map[uint]bool)
	for _, doc := range processFormDTO.DocumentIds {
		newDocumentIDs[*doc] = true
	}

	var documentsToAdd []*models.ProcessAssociatedDocument
	for _, newDocumentID := range processFormDTO.DocumentIds {
		if _, exists := existingDocumentIDs[*newDocumentID]; !exists {
			document := &models.ProcessAssociatedDocument{
				ProcessID:   id,
				DocumentID:  *newDocumentID,
				CreatedByID: currentUserID,
			}
			documentsToAdd = append(documentsToAdd, document)
		}
	}

	if len(documentsToAdd) > 0 {
		if err := ps.repo.CreateProcessDocuments(documentsToAdd, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	var deleteDocumentIDs []*uint
	for _, existingDocument := range existingProcess.Documents {
		if _, exists := newDocumentIDs[existingDocument.DocumentID]; !exists {
			deleteDocumentIDs = append(deleteDocumentIDs, existingDocument.ID)
		}
	}

	if len(deleteDocumentIDs) > 0 {
		if err := ps.repo.DeleteProcessDocuments(deleteDocumentIDs, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()

	go func() {
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

		// Determine asset register and document changes
		assetRegisterAdded := []uint{}
		assetRegisterRemoved := []uint{}
		for id := range newAssetRegisterIDs {
			if !existingAssetRegisterIDs[id] {
				assetRegisterAdded = append(assetRegisterAdded, id)
			}
		}
		for id := range existingAssetRegisterIDs {
			if !newAssetRegisterIDs[id] {
				assetRegisterRemoved = append(assetRegisterRemoved, id)
			}
		}

		documentsAdded := []uint{}
		documentsRemoved := []uint{}
		for id := range newDocumentIDs {
			if !existingDocumentIDs[id] {
				documentsAdded = append(documentsAdded, id)
			}
		}
		for id := range existingDocumentIDs {
			if !newDocumentIDs[id] {
				documentsRemoved = append(documentsRemoved, id)
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
		if len(assetRegisterAdded) > 0 || len(assetRegisterRemoved) > 0 {
			associatedItems["asset_register_ids"] = map[string]interface{}{
				"added":   assetRegisterAdded,
				"removed": assetRegisterRemoved,
			}
		}
		if len(documentsAdded) > 0 || len(documentsRemoved) > 0 {
			associatedItems["document_ids"] = map[string]interface{}{
				"added":   documentsAdded,
				"removed": documentsRemoved,
			}
		}

		if len(associatedItems) > 0 {
			changedData["associated_items"] = associatedItems
		}
		//log changes
		oldData := map[string]interface{}{
			"process_group":    oldProcess.ProcessGroup.TitleEn,
			"owner_unit":       oldProcess.OwnerUnit.TitleEn,
			"owner":            fmt.Sprintf("%s %s", oldProcess.Owner.FirstNameEN, oldProcess.Owner.LastNameEN),
			"lifecycle_status": oldProcess.LifecycleStatus.TitleEn,
			"title_en":         oldProcess.TitleEn,
			"title_ar":         oldProcess.TitleAr,
			"description":      oldProcess.Description,
		}

		newData := map[string]interface{}{
			"process_group":    existingProcess.ProcessGroup.TitleEn,
			"owner_unit":       existingProcess.OwnerUnit.TitleEn,
			"owner":            fmt.Sprintf("%s %s", existingProcess.Owner.FirstNameEN, existingProcess.Owner.LastNameEN),
			"lifecycle_status": existingProcess.LifecycleStatus.TitleEn,
			"title_en":         existingProcess.TitleEn,
			"title_ar":         existingProcess.TitleAr,
			"description":      existingProcess.Description,
		}
		// Compare other fields and add them to changedData
		fieldChanges := utility.CompareFieldChanges(oldData, newData)
		for k, v := range fieldChanges {
			changedData[k] = v
		}

		// Log activity if there are any changes
		if len(changedData) > 0 {
			err := ps.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.Process{}.TableName(),
				id,
				changedData,
				currentUserID,
				nil,
			)
			if err != nil {
				utility.LogError("Failed to log activity", map[string]interface{}{
					"action":   models.UpdateAction,
					"error":    err.Error(),
					"recordID": id,
				})
			}
		}

		//notification
		if oldOwnerID != processFormDTO.OwnerID {
			urlParam := fmt.Sprintf(`{"id":%d}`, id)
			queryParam := ""
			entity := "process"
			action := "responsible-user-remove"

			oldOwnerNotificationDTO := &commonDTOs.DBNotificationDTO{
				UserID:               oldOwnerID,
				ModuleIdentifier:     "organization",
				NotificationType:     string(commonDTOs.InformationType),
				ReferenceNumber:      process.ReferenceNumber,
				HeadingEn:            "Process Responsibility Removed",
				HeadingAr:            "تمت إزالة مسؤولية العملية",
				TitleEn:              processFormDTO.TitleEn,
				TitleAr:              processFormDTO.TitleAr,
				Entity:               &entity,
				Action:               &action,
				URLParam:             urlParam,
				QueryParam:           queryParam,
				IsSystemNotification: false,
				CreatedByID:          &currentUserID,
			}

			oldOwnerUser, err := ps.userSvc.GetByID(oldOwnerID)
			if err != nil {
				utility.LogError("Failed to fetch old responsible user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}

			oldOwnerEmailTo := []string{oldOwnerUser.Email}
			oldOwnerSubject := fmt.Sprintf("Process/%s: Ownership Removed", process.TitleEn)
			oldOwnerTemplateRelativePath := "organization/templates/email/process_remove_owner_mail.html"

			appURL := utility.GetAppURL()

			oldOwnerEmailData := map[string]interface{}{
				"FirstNameEn":            oldOwnerUser.FirstNameEn,
				"LastNameEn":             oldOwnerUser.LastNameEn,
				"ProcessReferenceNumber": existingProcess.ReferenceNumber,
				"ProcessName":            existingProcess.TitleEn,
				"ProcessGroup":           existingProcess.ProcessGroup.TitleEn,
			}

			oldOwnerEmailNotificationDTO := &commonDTOs.EmailNotificationDTO{
				EmailTo:              oldOwnerEmailTo,
				Subject:              oldOwnerSubject,
				TemplateRelativePath: oldOwnerTemplateRelativePath,
				EmailData:            oldOwnerEmailData,
			}

			//new owner notification
			newOwnerAction := "responsible-user-assign"

			dbNotificationDTO := &commonDTOs.DBNotificationDTO{
				UserID:               processFormDTO.OwnerID,
				ModuleIdentifier:     "organization",
				NotificationType:     string(commonDTOs.AssignmentType),
				ReferenceNumber:      process.ReferenceNumber,
				HeadingEn:            "Process Responsibility Assigned",
				HeadingAr:            "مسؤولية العملية المعينة",
				TitleEn:              processFormDTO.TitleEn,
				TitleAr:              processFormDTO.TitleAr,
				Entity:               &entity,
				Action:               &newOwnerAction,
				URLParam:             urlParam,
				QueryParam:           queryParam,
				IsSystemNotification: false,
				CreatedByID:          &currentUserID,
			}

			responsibleUser, err := ps.userSvc.GetByID(processFormDTO.OwnerID)
			if err != nil {
				utility.LogError("Failed to fetch responsible user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}

			emailTo := []string{responsibleUser.Email}
			subject := fmt.Sprintf("Process/%s: Ownership Changed", process.TitleEn)
			templateRelativePath := "organization/templates/email/process_change_owner_mail.html"

			emailData := map[string]interface{}{
				"FirstNameEn":            responsibleUser.FirstNameEn,
				"LastNameEn":             responsibleUser.LastNameEn,
				"ProcessReferenceNumber": existingProcess.ReferenceNumber,
				"ProcessName":            existingProcess.TitleEn,
				"ProcessGroup":           existingProcess.ProcessGroup.TitleEn,
				"NewOwnerFirstNameEn":    responsibleUser.FirstNameEn,
				"NewOwnerLastNameEn":     responsibleUser.LastNameEn,
				"AssignedByFirstNameEn":  existingProcess.CreatedByUser.FirstNameEN,
				"AssignedByLastNameEn":   existingProcess.CreatedByUser.LastNameEN,
				"AssignedOn":             existingProcess.CreatedAt.Format("2006-01-02"),
				"RedirectURL":            fmt.Sprintf("%s/process/%d", appURL, id),
			}

			emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
				EmailTo:              emailTo,
				Subject:              subject,
				TemplateRelativePath: templateRelativePath,
				EmailData:            emailData,
			}

			err = ps.notificationSvc.SendNotification(
				dbNotificationDTO, emailNotificationDTO, "process")
			if err != nil {
				utility.LogError("Failed to send notifications", map[string]interface{}{
					"error": err.Error(),
				})
			}

			err = ps.notificationSvc.SendNotification(
				oldOwnerNotificationDTO, oldOwnerEmailNotificationDTO, "process")
			if err != nil {
				utility.LogError("Failed to send old owner notifications", map[string]interface{}{
					"error": err.Error(),
				})

			}
		}

	}()

	return nil
}

func (ps *ProcessService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	tx := ps.repo.BeginTransaction()

	var failedIds []uint
	var successfulIds []uint
	var notifications []struct {
		OwnerID         *uint
		TitleEn         string
		TitleAr         string
		ReferenceNumber string
		ProcessGroup    string
	}
	var deletedProcesses []map[string]interface{}

	for _, id := range ids {
		process, err := ps.GetFullDetailsByID(id, tx)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		// Prepare notification data before deletion
		notifications = append(notifications, struct {
			OwnerID         *uint
			TitleEn         string
			TitleAr         string
			ReferenceNumber string
			ProcessGroup    string
		}{
			OwnerID:         process.OwnerID,
			TitleEn:         process.TitleEn,
			TitleAr:         process.TitleAr,
			ReferenceNumber: process.ReferenceNumber,
			ProcessGroup:    process.ProcessGroup.TitleEn,
		})

		deletedProcesses = append(deletedProcesses, map[string]interface{}{
			"id":       id,
			"title_en": process.TitleEn,
			"title_ar": process.TitleAr,
		})

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

		if len(process.Assets) > 0 {
			var deleteAssetRegisterIds []*uint
			for _, ctrl := range process.Assets {
				deleteAssetRegisterIds = append(deleteAssetRegisterIds, ctrl.ID)
			}
			if err := ps.repo.DeleteProcessAssetRegister(deleteAssetRegisterIds, tx); err != nil {
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
			// Log bulk delete activities
			err := ps.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.Process{}.TableName(),
				successfulIds,
				deletedProcesses,
				currentUserID,
				nil,
			)
			if err != nil {
				fmt.Printf("Failed to log bulk delete activity: %v\n", err)
			}

			// Send notifications to the owners of the deleted processes
			for _, notification := range notifications {
				urlParam := ""
				queryParam := ""
				entity := "process"
				action := "delete"

				dbNotificationDTO := &commonDTOs.DBNotificationDTO{
					UserID:               *notification.OwnerID,
					ModuleIdentifier:     "organization",
					NotificationType:     string(commonDTOs.DeletionType),
					ReferenceNumber:      notification.ReferenceNumber,
					HeadingEn:            "Process Deleted",
					HeadingAr:            "تم حذف العملية",
					TitleEn:              notification.TitleEn,
					TitleAr:              notification.TitleAr,
					Entity:               &entity,
					Action:               &action,
					URLParam:             urlParam,
					QueryParam:           queryParam,
					IsSystemNotification: false,
					CreatedByID:          &currentUserID,
				}

				ownerUser, err := ps.userSvc.GetByID(*notification.OwnerID)
				if err != nil {
					utility.LogError("Failed to fetch owner details for notification", map[string]interface{}{
						"error": err.Error(),
					})
					continue
				}

				ownerEmailTo := []string{ownerUser.Email}
				subject := fmt.Sprintf("Process/%s: Deleted", notification.TitleEn)
				templateRelativePath := "organization/templates/email/process_delete_mail.html"

				emailData := map[string]interface{}{
					"FirstNameEn":            ownerUser.FirstNameEn,
					"LastNameEn":             ownerUser.LastNameEn,
					"ProcessReferenceNumber": notification.ReferenceNumber,
					"ProcessName":            notification.TitleEn,
					"ProcessGroup":           notification.ProcessGroup,
				}

				emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
					EmailTo:              ownerEmailTo,
					Subject:              subject,
					TemplateRelativePath: templateRelativePath,
					EmailData:            emailData,
				}

				err = ps.notificationSvc.SendNotification(
					dbNotificationDTO, emailNotificationDTO, "process")
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

func (ps *ProcessService) ExportProcesses(data []*organizationDTOs.ProcessListDTO, currentUserID uint) (*os.File, error) {
	userSetting, err := ps.userSettingSvc.GetUserSetting(currentUserID)
	if err != nil {
		return nil, err
	}

	fmt.Println(userSetting)

	var headers []string
	if userSetting.Language == "ar" {
		headers = []string{
			"SL No.",
			"Reference Number",
			"Title",
			"Description",
			"Process Group",
			"Owner Unit",
			"Owner Name",
			"Lifecycle Status",
		}
	} else {
		headers = []string{
			"SL No.",
			"Reference Number",
			"Title",
			"Description",
			"Process Group",
			"Owner Unit",
			"Owner Name",
			"Lifecycle Status",
		}
	}

	rows := make([][]interface{}, len(data))
	for i, process := range data {
		if userSetting.Language == "ar" {
			rows[i] = []interface{}{
				i + 1,
				process.ReferenceNumber,
				process.TitleAr,
				utility.GetStringValue(process.Description),
				process.ProcessGroupTitleAr,
				utility.GetStringValue(process.OwnerUnitTitleAr),
				utility.GetStringValue(process.OwnerFirstNameAr) + " " + utility.GetStringValue(process.OwnerLastNameAr),
				process.LifecycleStatusTitleAr,
			}
		} else {
			rows[i] = []interface{}{
				i + 1,
				process.ReferenceNumber,
				process.TitleEn,
				utility.GetStringValue(process.Description),
				process.ProcessGroupTitleEn,
				utility.GetStringValue(process.OwnerUnitTitleEn),
				utility.GetStringValue(process.OwnerFirstNameEn) + " " + utility.GetStringValue(process.OwnerLastNameEn),
				process.LifecycleStatusTitleEn,
			}
		}
	}

	return utility.GenerateExcelFile(headers, rows)
}

func (ps *ProcessService) UpdateProcessRiskRating(process *models.Process, tx utility.DatabaseTransaction) error {

	err := ps.repo.Update(process, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (ps *ProcessService) GenerateProcessTemplate() (*os.File, error) {
	headers := []string{
		"ProcessGroupID", "OwnerID", "TitleEn", "TitleAr", "Description", "ControlIds", "DocumentIds",
	}

	// Pass empty data for the template
	data := [][]interface{}{}

	file, err := utility.GenerateExcelFile(headers, data)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func (ps *ProcessService) ImportProcesses(file multipart.File, currentUserID uint) ([]int, error) {
	mandatoryFields := []string{"ProcessGroupID", "OwnerID", "TitleEn", "TitleAr"}
	// Parse the Excel file
	data, err := utility.ParseExcelFile(file, mandatoryFields)
	if err != nil {
		return nil, fmt.Errorf("invalid file format: %v", err)
	}

	// Track failed rows
	var failedRows []int

	// Process each row
	for i, row := range data {
		// Map the row data to the ProcessFormDTO
		processForm := organizationDTOs.ProcessFormDTO{
			ProcessGroupID: utility.ParseUint(row["ProcessGroupID"]),
			OwnerID:        utility.ParseUint(row["OwnerID"]),
			TitleEn:        row["TitleEn"],
			TitleAr:        row["TitleAr"],
			Description:    utility.GetStringPointer(row["Description"]),
			ControlIds:     utility.ParseUintSlice(row["ControlIds"]),
			DocumentIds:    utility.ParseUintSlice(row["DocumentIds"]),
		}

		// Validate and create the process
		_, err := ps.Create(currentUserID, &processForm)
		if err != nil {
			// Log the row number that failed (+2 to account for the header row)
			failedRows = append(failedRows, i+2)
		}
	}

	return failedRows, nil
}

func (ps *ProcessService) UpdateProcessRecoveryObjective(process *models.Process, tx utility.DatabaseTransaction) error {

	err := ps.repo.Update(process, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}
