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
}

type ThreatService struct {
	repo             repositories.ThreatRepositoryInterface
	activityLogSvc   logServices.ActivityLogServiceInterface
	threatTypeSvc    ThreatTypeServiceInterface
	severityLevelSvc masterServices.SeverityLevelServiceInterface
	threatStatusSvc  ThreatStatusServiceInterface
	userSvc          organizationServices.UserServiceInterface
}

func NewThreatService(
	repo repositories.ThreatRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	threatTypeSvc ThreatTypeServiceInterface,
	severityLevelSvc masterServices.SeverityLevelServiceInterface,
	threatStatusSvc ThreatStatusServiceInterface,
	userSvc organizationServices.UserServiceInterface,
) ThreatServiceInterface {
	return &ThreatService{
		repo:             repo,
		activityLogSvc:   activityLogSvc,
		threatTypeSvc:    threatTypeSvc,
		severityLevelSvc: severityLevelSvc,
		threatStatusSvc:  threatStatusSvc,
		userSvc:          userSvc,
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
		err := ts.activityLogSvc.LogActivity(
			models.CreateAction,
			threat.TableName(),
			id,
			threatFormDTO,
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
				deleteDocumentIds = append(deleteDocumentIds, doc.ID)
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
	for _, docID := range deleteDocumentIds {
		doc, err := ts.repo.GetDocumentByID(*docID)
		if err == nil {
			storagePath := filepath.Join("asset", "threat", strconv.FormatUint(uint64(id), 10), "documents")
			if err := utility.RemoveFile(*doc.Name, storagePath); err != nil {
				fmt.Printf("Failed to delete file %s: %v\n", *doc.Name, err)
			}
		}
	}
	go func() {
		err := ts.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.Threat{}.TableName(),
			id,
			threatFormDTO,
			currentUserID,
			nil,
		)
		if err != nil {
			fmt.Printf("Failed to log activity: %v\n", err)
		}
	}()

	return nil

}

func (ts *ThreatService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	tx := ts.repo.BeginTransaction()

	var failedIds []uint
	var successfulIds []uint
	documentsToDelete := make(map[uint][]*dtos.ThreatDocumentDTO)

	for _, id := range ids {
		var deleteDocumentIds []*uint
		existingThreatDetails, err := ts.GetFullDetailsByID(id, tx)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		if existingThreatDetails.Documents != nil && len(existingThreatDetails.Documents) > 0 {
			for _, document := range existingThreatDetails.Documents {
				documentsToDelete[id] = append(documentsToDelete[id], document)
				deleteDocumentIds = append(deleteDocumentIds, document.ID)
			}

			if len(deleteDocumentIds) > 0 {
				if err := ts.repo.DeleteDocuments(deleteDocumentIds, tx); err != nil {
					failedIds = append(failedIds, id)
					continue
				}
			}
		}
		err = ts.repo.Delete(id)
		if err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}
	}

	if len(successfulIds) > 0 {
		for threatId, documents := range documentsToDelete {
			for _, doc := range documents {
				storagePath := filepath.Join("asset", "threat", strconv.FormatUint(uint64(threatId), 10), "documents")
				if err := utility.RemoveFile(doc.Name, storagePath); err != nil {
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
				nil,
				currentUserID,
				nil,
			)
			if err != nil {
				fmt.Printf("Failed to log bulk delete activity: %v\n", err)
			}
		}()
	}

	tx.Commit()

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
