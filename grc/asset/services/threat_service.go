package services

import (
	"errors"
	"fmt"
	"jethur/grc/asset/dtos"
	"jethur/grc/asset/repositories"
	logServices "jethur/grc/logs/services"
	masterServices "jethur/grc/masters/services"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	"jethur/utility"
	"path/filepath"
	"strconv"
	"time"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type ThreatServiceInterface interface {
	GetAll(params *dtos.ThreatListQueryParamsDTO) (*dtos.ThreatListResponseDTO, error)
	GetByID(id uint) (*models.Threat, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*dtos.ThreatDetailDTO, error)
	Create(threatForm *dtos.ThreatFormDTO, currentUserID uint) (uint, error)
	Update(id uint, threatForm *dtos.ThreatFormDTO, currentUserID uint) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)

	GetThreatSummary(assetId uint) (*dtos.ThreatSummaryDTO, error)
}

type ThreatService struct {
	repo             repositories.ThreatRepositoryInterface
	activityLogSvc   logServices.ActivityLogServiceInterface
	threatTypeSvc    ThreatTypeServiceInterface
	severityLevelSvc masterServices.SeverityLevelServiceInterface
	threatStatusSvc  ThreatStatusServiceInterface
	userSvc          organizationServices.UserServiceInterface
	vulnerabilitySvc VulnerabilityServiceInterface
}

func NewThreatService(
	repo repositories.ThreatRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	threatTypeSvc ThreatTypeServiceInterface,
	severityLevelSvc masterServices.SeverityLevelServiceInterface,
	threatStatusSvc ThreatStatusServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	vulnerabilitySvc VulnerabilityServiceInterface,
) ThreatServiceInterface {
	return &ThreatService{
		repo:             repo,
		activityLogSvc:   activityLogSvc,
		threatTypeSvc:    threatTypeSvc,
		severityLevelSvc: severityLevelSvc,
		threatStatusSvc:  threatStatusSvc,
		userSvc:          userSvc,
		vulnerabilitySvc: vulnerabilitySvc,
	}
}

func (ts *ThreatService) GetAll(params *dtos.ThreatListQueryParamsDTO) (*dtos.ThreatListResponseDTO, error) {
	return ts.repo.GetAll(params)
}

func (ts *ThreatService) GetByID(id uint) (*models.Threat, error) {
	threat, err := ts.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return threat, nil
}

func (ts *ThreatService) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*dtos.ThreatDetailDTO, error) {
	var threatDetailDTO dtos.ThreatDetailDTO

	threat, err := ts.repo.GetFullDetailsByID(id, tx...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	err = mapstructure.Decode(threat, &threatDetailDTO)
	if err != nil {
		return nil, err
	}

	if threat.DetectionDate != nil {
		detectionDateString := threat.DetectionDate.Format("2006-01-02")
		threatDetailDTO.DetectionDateString = &detectionDateString
	}
	if threat.ResolutionDate != nil {
		resolutionDateString := threat.ResolutionDate.Format("2006-01-02")
		threatDetailDTO.ResolutionDateString = &resolutionDateString
	}

	// Add additional fields as necessary for the details screen
	// For example, setting relationship data or additional fields as required by the DTO

	return &threatDetailDTO, nil
}

func (ts *ThreatService) Create(threatFormDTO *dtos.ThreatFormDTO, currentUserID uint) (uint, error) {
	var tx = ts.repo.BeginTransaction()

	// Validate the DTO
	if err := ts.validateThreatFormDTO(threatFormDTO); err != nil {
		return 0, err
	}

	// Parse dates from the DTO
	dates, err := ts.parseThreatDates(threatFormDTO)
	if err != nil {
		return 0, err
	}

	// Retrieve the primary organization unit for the responsible user if not provided
	organizationUnitID, err := ts.userSvc.GetUserPrimaryUnitID(threatFormDTO.ResponsibleUserID)
	if err != nil {
		return 0, utility.NewValidationError(fmt.Errorf("organization unit not set for responsible user"), "not_exist", "responsible user unit")
	}

	// Populate the Threat model with data from the DTO
	threat := &models.Threat{
		TitleEn:           threatFormDTO.TitleEn,
		TitleAr:           threatFormDTO.TitleAr,
		ThreatTypeID:      threatFormDTO.ThreatTypeID,
		SeverityLevelID:   threatFormDTO.SeverityLevelID,
		ThreatStatusID:    threatFormDTO.ThreatStatusID,
		ResponsibleUserID: threatFormDTO.ResponsibleUserID,
		ResponsibleUnitID: organizationUnitID,
		VulnerabilityID:   threatFormDTO.VulnerabilityID,
		Description:       threatFormDTO.Description,
		Source:            threatFormDTO.Source,
		DetectionDate:     dates.DetectionDate,
		ResolutionDate:    dates.ResolutionDate,
		CreatedBYID:       currentUserID,
		CreatedAt:         time.Now(),
	}

	// Save the threat record in the repository
	id, err := ts.repo.Create(threat)
	if err != nil {
		return 0, err
	}
	var threatDocuments []*models.ThreatDocument

	for _, doc := range threatFormDTO.Documents {
		// Encrypt the document file name
		dataToEncrypt := []byte(doc.FileName)
		encryptedFileName, err := utility.EncryptData(dataToEncrypt)
		if err != nil {
			return 0, err
		}

		// Generate the document URL
		url := utility.GenerateDocumentURL("asset", "threat", strconv.FormatUint(uint64(id), 10), "documents", encryptedFileName)

		// Create a ThreatDocument model instance
		document := &models.ThreatDocument{
			ThreatID:     id,
			URL:          &url,
			Name:         &doc.FileName,
			Size:         &doc.Size,
			Token:        &encryptedFileName,
			OriginalName: &doc.OriginalName,
			Type:         &doc.Type,
			CreatedBYID:  currentUserID,
		}
		threatDocuments = append(threatDocuments, document)
	}

	// Save documents to the database
	if len(threatDocuments) > 0 {
		if err := ts.repo.CreateDocuments(threatDocuments, tx); err != nil {
			return 0, err
		}
	}

	// Store files in the specified path
	for _, doc := range threatDocuments {
		storagePath := filepath.Join("asset", "threat", strconv.FormatUint(uint64(id), 10), "documents")
		if err := utility.StoreFile(*doc.Name, storagePath); err != nil {
			return 0, err
		}
	}
	tx.Commit()

	// Log the activity asynchronously
	go func() {
		existingThreat, err := ts.GetFullDetailsByID(id)
		if err != nil {
			utility.LogError("Failed to fetch threat", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		err = ts.activityLogSvc.LogActivity(
			models.CreateAction,
			threat.TableName(),
			id,
			existingThreat,
			currentUserID,
			nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": threat.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}
			utility.LogError("Failed to log threat activity", fields)
		}
	}()

	return id, nil
}

func (ts *ThreatService) Update(id uint, threatFormDTO *dtos.ThreatFormDTO, currentUserID uint) error {
	var tx = ts.repo.BeginTransaction()

	oldThreat, err := ts.GetFullDetailsByID(id, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Retrieve the existing threat record
	existingThreat, err := ts.repo.GetByID(id)
	if err != nil {
		return err
	}

	// Validate the incoming DTO
	if err := ts.validateThreatFormDTO(threatFormDTO); err != nil {
		return err
	}

	// Parse dates for DetectionDate and ResolutionDate from DTO
	dates, err := ts.parseThreatDates(threatFormDTO)
	if err != nil {
		return err
	}

	// Check if ResponsibleUserID has changed and update ResponsibleUnitID if necessary
	if existingThreat.ResponsibleUserID != threatFormDTO.ResponsibleUserID {
		organizationUnitID, err := ts.userSvc.GetUserPrimaryUnitID(threatFormDTO.ResponsibleUserID)
		if err != nil {
			return utility.NewValidationError(fmt.Errorf("organization unit not set for responsible user"), "not_exist", "responsible user unit")
		}
		existingThreat.ResponsibleUnitID = organizationUnitID
	}

	// Update the fields with new values from the DTO
	existingThreat.TitleEn = threatFormDTO.TitleEn
	existingThreat.TitleAr = threatFormDTO.TitleAr
	existingThreat.ThreatTypeID = threatFormDTO.ThreatTypeID
	existingThreat.SeverityLevelID = threatFormDTO.SeverityLevelID
	existingThreat.ThreatStatusID = threatFormDTO.ThreatStatusID
	existingThreat.ResponsibleUserID = threatFormDTO.ResponsibleUserID
	existingThreat.VulnerabilityID = threatFormDTO.VulnerabilityID
	existingThreat.Description = threatFormDTO.Description
	existingThreat.Source = threatFormDTO.Source
	existingThreat.DetectionDate = dates.DetectionDate
	existingThreat.ResolutionDate = dates.ResolutionDate
	existingThreat.UpdatedBYID = &currentUserID

	// Update
	if err := ts.repo.Update(existingThreat, tx); err != nil {
		tx.Rollback()
		return err
	}
	// Prepare slices for new and deleted documents
	var newDocuments []*models.ThreatDocument
	var deleteDocumentIds []*uint
	var documentsToDelete []*models.ThreatDocument

	// Process each document in the ThreatFormDTO
	for _, doc := range threatFormDTO.Documents {
		if doc.IsNewFile != nil && *doc.IsNewFile {
			// Encrypt the filename
			dataToEncrypt := []byte(doc.FileName)
			encryptedFileName, err := utility.EncryptData(dataToEncrypt)
			if err != nil {
				tx.Rollback()
				return err
			}

			// Generate URL for the document
			url := utility.GenerateDocumentURL("asset", "threat", strconv.FormatUint(uint64(id), 10), "documents", encryptedFileName)

			// Create a new ThreatDocument instance
			document := &models.ThreatDocument{
				ThreatID:     id,
				URL:          &url,
				Name:         &doc.FileName,
				Size:         &doc.Size,
				Token:        &encryptedFileName,
				OriginalName: &doc.OriginalName,
				Type:         &doc.Type,
				CreatedBYID:  currentUserID,
			}

			// Append to the list of new documents
			newDocuments = append(newDocuments, document)
		} else if doc.IsDeletedFile != nil && *doc.IsDeletedFile {
			if doc.ID != nil {
				existingDocument, err := ts.repo.GetDocumentByID(*doc.ID)
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

	// Create new documents if any
	if len(newDocuments) > 0 {
		if err := ts.repo.CreateDocuments(newDocuments, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	// Delete documents if any
	if len(deleteDocumentIds) > 0 {
		if err := ts.repo.DeleteDocuments(deleteDocumentIds, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	// Commit transaction
	tx.Commit()

	// Store new files in storage
	for _, doc := range newDocuments {
		storagePath := filepath.Join("asset", "threat", strconv.FormatUint(uint64(id), 10), "documents")
		if err := utility.StoreFile(*doc.Name, storagePath); err != nil {
			fmt.Printf("Failed to store file %s: %v\n", *doc.Name, err)
		}
	}

	// Remove deleted files from storage
	for _, doc := range documentsToDelete {
		storagePath := filepath.Join("asset", "threat", strconv.FormatUint(uint64(id), 10), "documents")
		if err := utility.RemoveFile(*doc.Name, storagePath); err != nil {
			fmt.Printf("Failed to delete file %s: %v\n", *doc.Name, err)
		}
	}
	go func() {
		existingThreat, err := ts.GetFullDetailsByID(id)
		if err != nil {
			return
		}
		// Prepare log data
		changedData := make(map[string]interface{})
		//log changes
		oldData := map[string]interface{}{
			"threat_type":      oldThreat.ThreatType.TitleEn,
			"severity_level":   oldThreat.SeverityLevel.TitleEn,
			"responsible_user": fmt.Sprintf("%s %s", oldThreat.ResponsibleUser.FirstNameEN, oldThreat.ResponsibleUser.LastNameEN),
			"title_en":         oldThreat.TitleEn,
			"title_ar":         oldThreat.TitleAr,
			"description":      oldThreat.Description,
			"threat_status":    oldThreat.ThreatStatus.TitleEn,
			"vulnerability":    utility.GetFieldValue(oldThreat.Vulnerability),
			"source":           oldThreat.Source,
			"detection_date":   oldThreat.DetectionDateString,
			"resolution_date":  oldThreat.ResolutionDateString,
		}

		newData := map[string]interface{}{
			"threat_type":      existingThreat.ThreatType.TitleEn,
			"severity_level":   existingThreat.SeverityLevel.TitleEn,
			"responsible_user": fmt.Sprintf("%s %s", existingThreat.ResponsibleUser.FirstNameEN, existingThreat.ResponsibleUser.LastNameEN),
			"title_en":         existingThreat.TitleEn,
			"title_ar":         existingThreat.TitleAr,
			"description":      existingThreat.Description,
			"threat_status":    existingThreat.ThreatStatus.TitleEn,
			"vulnerability":    utility.GetFieldValue(existingThreat.Vulnerability),
			"source":           existingThreat.Source,
			"detection_date":   existingThreat.DetectionDateString,
			"resolution_date":  existingThreat.ResolutionDateString,
		}
		// Compare other fields and add them to changedData
		fieldChanges := utility.CompareFieldChanges(oldData, newData)
		for k, v := range fieldChanges {
			changedData[k] = v
		}
		// Log activity if there are any changes
		if len(changedData) > 0 {
			err := ts.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.Threat{}.TableName(),
				id,
				changedData,
				currentUserID,
				nil,
			)
			if err != nil {
				fields := map[string]interface{}{
					"action":    models.UpdateAction,
					"tableName": models.Threat{}.TableName(),
					"recordID":  id,
					"userID":    currentUserID,
					"error":     err.Error(),
				}
				fmt.Printf("Failed to log activity: %v\n", fields)
			}
		}
	}()

	return nil

}

func (ts *ThreatService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	tx := ts.repo.BeginTransaction()

	var failedIds []uint
	var successfulIds []uint
	var deletedThreat []map[string]interface{}
	documentsToDelete := make(map[uint][]*dtos.ThreatDocumentDTO)
	for _, id := range ids {
		var deleteDocumentIds []*uint
		existingThreatDetails, err := ts.GetFullDetailsByID(id, tx)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		if existingThreatDetails.Documents != nil && len(existingThreatDetails.Documents) > 0 {
			for _, doc := range existingThreatDetails.Documents {
				documentsToDelete[id] = append(documentsToDelete[id], doc)
				deleteDocumentIds = append(deleteDocumentIds, doc.ID)

			}
			if err := ts.repo.DeleteDocuments(deleteDocumentIds, tx); err != nil {
				tx.Rollback()
				return nil, err
			}

		}

		if err := ts.repo.Delete(id, tx); err != nil {
			failedIds = append(failedIds, id)
			continue
		} else {
			successfulIds = append(successfulIds, id)
		}

		deletedThreat = append(deletedThreat, map[string]interface{}{
			"id":       id,
			"title_en": existingThreatDetails.TitleEn,
			"title_ar": existingThreatDetails.TitleAr,
		})
	}

	if len(failedIds) > 0 {
		tx.Rollback()
		return failedIds, nil
	}
	if len(successfulIds) > 0 {
		tx.Commit()
		for threatId, documents := range documentsToDelete {
			for _, doc := range documents {
				// storagePath := filepath.Join("asset", "threat", strconv.FormatUint(uint64(threatId), 10), "documents")
				// if err := utility.RemoveFile(doc.Name, storagePath); err != nil {
				// 	tx.Rollback()
				// 	fmt.Printf("Failed to delete file at URL %s: %v\n", *doc.URL, err)
				// }
				storagePath := filepath.Join("asset", "threat", strconv.FormatUint(uint64(threatId), 10))
				if err := utility.RemoveFolder(storagePath); err != nil {
					tx.Rollback()
					fmt.Printf("Failed to delete file at URL %s: %v\n", *doc.URL, err)
				}

			}
		}
		go func() {
			err := ts.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.Threat{}.TableName(),
				successfulIds,
				deletedThreat,
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

func (ts *ThreatService) validateThreatFormDTO(threatFormDTO *dtos.ThreatFormDTO) error {
	// Validate Threat Type
	if _, err := ts.threatTypeSvc.GetByID(threatFormDTO.ThreatTypeID); err != nil {
		return utility.NewValidationError(fmt.Errorf("threat type does not exist"), "not_exist", "threat type")
	}

	// Validate Severity Level
	if _, err := ts.severityLevelSvc.GetByID(threatFormDTO.SeverityLevelID); err != nil {
		return utility.NewValidationError(fmt.Errorf("severity level does not exist"), "not_exist", "severity level")
	}

	// Validate Threat Status
	if _, err := ts.threatStatusSvc.GetByID(threatFormDTO.ThreatStatusID); err != nil {
		return utility.NewValidationError(fmt.Errorf("threat status does not exist"), "not_exist", "threat status")
	}

	// Validate Responsible User
	if _, err := ts.userSvc.GetByID(threatFormDTO.ResponsibleUserID); err != nil {
		return utility.NewValidationError(fmt.Errorf("responsible user does not exist"), "not_exist", "responsible user")
	}

	if threatFormDTO.VulnerabilityID != nil {
		if _, err := ts.vulnerabilitySvc.GetByID(*threatFormDTO.VulnerabilityID); err != nil {
			return utility.NewValidationError(fmt.Errorf("vulnerability does not exist"), "not_exist", "vulnerability")
		}
	}
	return nil
}

func (ts *ThreatService) parseThreatDates(threatRegisterFormDTO *dtos.ThreatFormDTO) (*dtos.ThreatParsedDates, error) {
	parse := &dtos.ThreatParsedDates{}

	var err error
	if threatRegisterFormDTO.DetectionDate != nil {
		parse.DetectionDate, err = utility.ParseDate(*threatRegisterFormDTO.DetectionDate)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid detection date"), "invalid_date", "detection date")
		}
	}
	if threatRegisterFormDTO.ResolutionDate != nil {
		parse.ResolutionDate, err = utility.ParseDate(*threatRegisterFormDTO.ResolutionDate)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid resolution date"), "invalid_date", "resolution date")
		}
	}
	return parse, nil
}

func (ts *ThreatService) GetThreatSummary(assetId uint) (*dtos.ThreatSummaryDTO, error) {
	openCount, closedCount, err := ts.repo.GetThreatSummary(assetId)
	if err != nil {
		return nil, err
	}
	totalCount := openCount + closedCount
	summary := &dtos.ThreatSummaryDTO{
		OpenCount:   openCount,
		ClosedCount: closedCount,
		TotalCount:  totalCount,
	}

	return summary, nil
}
