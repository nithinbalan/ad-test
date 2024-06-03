package services

import (
	"errors"
	"fmt"
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/incident/dtos"
	incidentDTOs "jethur/grc/incident/dtos"
	incidentRepositories "jethur/grc/incident/repositories"
	logServices "jethur/grc/logs/services"
	organizationServices "jethur/grc/organization/services"

	"jethur/grc/models"
	"jethur/utility"
	"path/filepath"
	"strconv"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type IncidentInvestigationServiceInterface interface {
	GetAll(incidentID uint, params *incidentDTOs.IncidentInvestigationListQueryParamsDTO) (*incidentDTOs.IncidentInvestigationListResponseDTO, error)
	GetFullDetailsByID(incidentID uint, id uint, tx ...utility.DatabaseTransaction) (*incidentDTOs.IncidentInvestigationDetailDTO, error)
	GetByID(id uint) (*models.IncidentInvestigation, error)
	Create(incidentID uint, currentUserID uint, incidentInvestigationFormDTO *incidentDTOs.IncidentInvestigationFormDTO) (uint, error)
	Update(incidentID uint, id uint, currentUserID uint, incident *incidentDTOs.IncidentInvestigationFormDTO) error
	Delete(incidentID uint, ids []uint, currentUserID uint) ([]uint, error)
}

// Service and implementations
type IncidentInvestigationService struct {
	repo                           incidentRepositories.IncidentInvestigationRepositoryInterface
	activityLogSvc                 logServices.ActivityLogServiceInterface
	incidentRegisterSvc            IncidentRegisterServiceInterface
	userSvc                        organizationServices.UserServiceInterface
	incidentInvestigationStatusSvc IncidentInvestigationStatusServiceInterface
}

func NewIncidentInvestigationService(
	repo incidentRepositories.IncidentInvestigationRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	incidentRegisterSvc IncidentRegisterServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	incidentInvestigationStatusSvc IncidentInvestigationStatusServiceInterface,
) IncidentInvestigationServiceInterface {
	return &IncidentInvestigationService{
		repo:                           repo,
		activityLogSvc:                 activityLogSvc,
		incidentRegisterSvc:            incidentRegisterSvc,
		userSvc:                        userSvc,
		incidentInvestigationStatusSvc: incidentInvestigationStatusSvc,
	}
}

func (iis *IncidentInvestigationService) GetAll(incidentID uint, params *incidentDTOs.IncidentInvestigationListQueryParamsDTO) (*incidentDTOs.IncidentInvestigationListResponseDTO, error) {
	if _, err := iis.incidentRegisterSvc.GetByID(incidentID); err != nil {
		return nil, err
	}
	return iis.repo.GetAll(incidentID, params)
}

func (iis IncidentInvestigationService) GetByID(id uint) (*models.IncidentInvestigation, error) {
	incident, err := iis.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return incident, nil
}

func (iis IncidentInvestigationService) GetFullDetailsByID(incidentID uint, id uint, tx ...utility.DatabaseTransaction) (*incidentDTOs.IncidentInvestigationDetailDTO, error) {
	if _, err := iis.incidentRegisterSvc.GetByID(incidentID); err != nil {
		return nil, err
	}
	// Initialize your DTO
	var incidentInvestigationDTO incidentDTOs.IncidentInvestigationDetailDTO

	incident, err := iis.repo.GetFullDetailsByID(id, tx...)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	// Map the result to your DTO
	err = mapstructure.Decode(incident, &incidentInvestigationDTO)
	if err != nil {
		return nil, err
	}

	//todo :- also add extra fields which required in details screen

	return &incidentInvestigationDTO, nil
}

func (iis *IncidentInvestigationService) Create(incidentID uint, currentUserID uint, incidentInvestigationFormDTO *dtos.IncidentInvestigationFormDTO) (uint, error) {
	tx := iis.repo.BeginTransaction()
	if _, err := iis.incidentRegisterSvc.GetByID(incidentID); err != nil {
		return 0, err
	}
	if err := iis.validateIncidentInvestigationFormDTO(incidentInvestigationFormDTO); err != nil {
		return 0, err
	}
	dates, err := iis.parseIncidentInvestigationDates(incidentInvestigationFormDTO)
	if err != nil {
		return 0, err
	}

	if err := iis.validateIncidentInvestigationDocumentFormDTO(incidentInvestigationFormDTO.Documents); err != nil {
		return 0, err
	}

	//Create Incident Investigation
	incidentInvestigation := &models.IncidentInvestigation{
		IncidentRegisterID:            incidentID,
		Description:                   incidentInvestigationFormDTO.Description,
		InitiatedAt:                   *dates.InitiatedAt,
		TeamLeadID:                    incidentInvestigationFormDTO.TeamLeadID,
		ActionRecommendation:          incidentInvestigationFormDTO.ActionRecommendation,
		IncidentInvestigationStatusID: incidentInvestigationFormDTO.IncidentInvestigationStatusID,
		CreatedBy:                     currentUserID,
	}

	id, err := iis.repo.Create(incidentInvestigation, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	// Save Team Members
	var users []*models.IncidentInvestigationUser
	for _, userId := range incidentInvestigationFormDTO.UserIds {
		user := &models.IncidentInvestigationUser{
			IncidentInvestigationID: id,
			UserID:                  *userId,
			CreatedBy:               currentUserID,
		}
		users = append(users, user)
	}

	if len(users) > 0 {
		if err := iis.repo.CreateMembers(users, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	// Save Findings
	var findings []*models.IncidentInvestigationFinding
	for _, finding := range incidentInvestigationFormDTO.Findings {
		findingObj := &models.IncidentInvestigationFinding{
			IncidentInvestigationID: id,
			Title:                   *finding,
			CreatedBy:               currentUserID,
		}
		findings = append(findings, findingObj)
	}

	if len(findings) > 0 {
		if err := iis.repo.CreateFindings(findings, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
	}
	// Save Documents
	var documents []*models.IncidentInvestigationDocument
	// Iterate over each document in the array
	for _, doc := range incidentInvestigationFormDTO.Documents {
		// Encrypt the filename
		dataToEncrypt := []byte(doc.FileName)
		encryptedFilename, err := utility.EncryptData(dataToEncrypt)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
		// Construct the URL for the document
		url := utility.GenerateDocumentURL("incident-investigations", strconv.FormatUint(uint64(incidentInvestigation.ID), 10), encryptedFilename)

		// Create a new Incident Investigation Document instance
		document := &models.IncidentInvestigationDocument{
			IncidentInvestigationID: id,
			URL:                     url,
			Name:                    doc.FileName,
			Size:                    doc.Size,
			Token:                   encryptedFilename,
			OriginalName:            doc.OriginalName,
			Type:                    doc.Type,
			CreatedBy:               currentUserID,
		}

		// Append the document to the slice
		documents = append(documents, document)
	}
	if len(documents) > 0 {
		if err := iis.repo.CreateDocuments(documents, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
		for _, doc := range documents {
			storagePath := filepath.Join("incident", "investigations", strconv.FormatUint(uint64(incidentInvestigation.ID), 10), "documents")
			if err := utility.StoreFile(doc.Name, storagePath); err != nil {
				tx.Rollback()
				return 0, err
			}
		}
	}

	tx.Commit()

	go func() {
		err := iis.activityLogSvc.LogActivity(
			models.CreateAction,
			models.IncidentInvestigation{}.TableName(),
			id,
			incidentInvestigationFormDTO,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.IncidentInvestigation{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return id, nil
}

func (iis *IncidentInvestigationService) Update(incidentID uint, id uint, currentUserID uint, incidentInvestigationFormDTO *dtos.IncidentInvestigationFormDTO) error {
	tx := iis.repo.BeginTransaction()
	if _, err := iis.incidentRegisterSvc.GetByID(incidentID); err != nil {
		return err
	}
	if err := iis.validateIncidentInvestigationFormDTO(incidentInvestigationFormDTO); err != nil {
		return err
	}
	dates, err := iis.parseIncidentInvestigationDates(incidentInvestigationFormDTO)
	if err != nil {
		return err
	}
	if err := iis.validateIncidentInvestigationDocumentFormDTO(incidentInvestigationFormDTO.Documents); err != nil {
		return err
	}
	// Update existing Incident Investigation
	existingIncidentInvestigation, err := iis.repo.GetByID(id)
	if err != nil {
		return err
	}
	if existingIncidentInvestigation == nil {
		return errors.New("existing incident investigation not found")
	}

	existingIncidentInvestigation.IncidentRegisterID = incidentID
	existingIncidentInvestigation.Description = incidentInvestigationFormDTO.Description
	existingIncidentInvestigation.InitiatedAt = *dates.InitiatedAt
	existingIncidentInvestigation.TeamLeadID = incidentInvestigationFormDTO.TeamLeadID
	existingIncidentInvestigation.ActionRecommendation = incidentInvestigationFormDTO.ActionRecommendation
	existingIncidentInvestigation.IncidentInvestigationStatusID = incidentInvestigationFormDTO.IncidentInvestigationStatusID
	existingIncidentInvestigation.UpdatedBy = &currentUserID

	err = iis.repo.Update(existingIncidentInvestigation, tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	// Collect existing user IDs
	existingIncidentInvestigationDetails, err := iis.repo.GetFullDetailsByID(id, tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	existingUserIDs := make(map[uint]bool)
	for _, user := range existingIncidentInvestigationDetails.Users {
		existingUserIDs[user.ID] = true
	}
	// Collect new user IDs
	newUserIDs := make(map[uint]bool)
	for _, user := range incidentInvestigationFormDTO.UserIds {
		newUserIDs[*user] = true
	}
	// Add new users
	var users []*models.IncidentInvestigationUser
	for _, newUser := range incidentInvestigationFormDTO.UserIds {
		if _, exists := existingUserIDs[*newUser]; !exists {
			// Create a new incident investigation user record
			user := &models.IncidentInvestigationUser{
				IncidentInvestigationID: existingIncidentInvestigationDetails.ID,
				UserID:                  *newUser,
				CreatedBy:               currentUserID,
			}
			users = append(users, user)
		}
	}
	// Save the new user
	if len(users) > 0 {
		if err := iis.repo.CreateMembers(users, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	// Remove users not in the new list
	var deleteUserIDs []*uint
	for _, existingUser := range existingIncidentInvestigationDetails.Users {
		if _, exists := newUserIDs[existingUser.ID]; !exists {
			// Delete the existing user from the repository
			deleteUserIDs = append(deleteUserIDs, &existingUser.ID)
		}
	}
	// Delete the user
	if len(deleteUserIDs) > 0 {
		if err := iis.repo.DeleteMembers(deleteUserIDs, tx); err != nil {
			tx.Rollback()
			return err
		}
	}
	// Collect Existing Finding Titles
	existingFindingTitles := make(map[string]bool)
	for _, finding := range existingIncidentInvestigationDetails.Findings {
		existingFindingTitles[finding.Title] = true
	}

	// Collect New Finding Titles
	newFindingTitles := make(map[string]bool)
	for _, finding := range incidentInvestigationFormDTO.Findings {
		newFindingTitles[*finding] = true
	}

	// Add New Findings
	var findingsToAdd []*models.IncidentInvestigationFinding
	for _, newFinding := range incidentInvestigationFormDTO.Findings {
		if _, exists := existingFindingTitles[*newFinding]; !exists {
			// Create a new incident investigation finding record
			finding := &models.IncidentInvestigationFinding{
				IncidentInvestigationID: existingIncidentInvestigationDetails.ID,
				Title:                   *newFinding,
				CreatedBy:               currentUserID,
			}
			findingsToAdd = append(findingsToAdd, finding)
		}
	}

	// Save the new findings
	if len(findingsToAdd) > 0 {
		if err := iis.repo.CreateFindings(findingsToAdd, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	// Remove findings not in the new list
	var findingsToDelete []*uint
	for _, existingFinding := range existingIncidentInvestigationDetails.Findings {
		if _, exists := newFindingTitles[existingFinding.Title]; !exists {
			// Delete the existing finding from the repository
			findingsToDelete = append(findingsToDelete, &existingFinding.ID)
		}
	}

	// Delete the findings
	if len(findingsToDelete) > 0 {
		if err := iis.repo.DeleteFindings(findingsToDelete, tx); err != nil {
			tx.Rollback()
			return err
		}
	}
	// Update or delete documents
	var documents []*models.IncidentInvestigationDocument
	var deleteDocumentIDs []*uint
	var documentsToDelete []*models.IncidentInvestigationDocument

	for _, doc := range incidentInvestigationFormDTO.Documents {
		if doc.IsNewFile != nil && *doc.IsNewFile {
			dataToEncrypt := []byte(doc.FileName)
			encryptedFilename, err := utility.EncryptData(dataToEncrypt)
			if err != nil {
				tx.Rollback()
				return err
			}
			url := utility.GenerateDocumentURL("incident-investigations", strconv.FormatUint(uint64(existingIncidentInvestigation.ID), 10), encryptedFilename)

			document := &models.IncidentInvestigationDocument{
				IncidentInvestigationID: existingIncidentInvestigation.ID,
				URL:                     url,
				Name:                    doc.FileName,
				Size:                    doc.Size,
				Token:                   encryptedFilename,
				Type:                    doc.Type,
				OriginalName:            doc.OriginalName,
				CreatedBy:               currentUserID,
			}
			documents = append(documents, document)
		} else if doc.IsDeletedFile != nil && *doc.IsDeletedFile {
			if doc.ID != nil {
				existingDocument, err := iis.repo.GetDocumentByID(*doc.ID)
				if err != nil {
					return err
				}
				if existingDocument == nil {
					tx.Rollback()
					return fmt.Errorf("existing document not found for ID: %d", *doc.ID)
				}
				deleteDocumentIDs = append(deleteDocumentIDs, &existingDocument.ID)
				documentsToDelete = append(documentsToDelete, existingDocument)
			}
		}
	}

	//Save Document
	if len(documents) > 0 {
		if err := iis.repo.CreateDocuments(documents, tx); err != nil {
			tx.Rollback()
			return err
		}
	}
	if len(deleteDocumentIDs) > 0 {
		if err := iis.repo.DeleteDocuments(deleteDocumentIDs, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	//Store File
	if len(documents) > 0 {
		for _, doc := range documents {
			storagePath := filepath.Join("incident", "investigations", strconv.FormatUint(uint64(existingIncidentInvestigation.ID), 10), "documents")
			if err := utility.StoreFile(doc.Name, storagePath); err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	//Remove File
	if len(documentsToDelete) > 0 {
		for _, document := range documentsToDelete {
			storagePath := filepath.Join("incident", "investigations", strconv.FormatUint(uint64(existingIncidentInvestigation.ID), 10), "documents")
			if err := utility.RemoveFile(document.Name, storagePath); err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	tx.Commit()

	go func() {
		err := iis.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.IncidentInvestigation{}.TableName(),
			id,
			incidentInvestigationFormDTO,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.UpdateAction,
				"tableName": models.IncidentInvestigation{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return nil
}

func (iis *IncidentInvestigationService) Delete(incidentID uint, ids []uint, currentUserID uint) ([]uint, error) {
	tx := iis.repo.BeginTransaction()
	if _, err := iis.incidentRegisterSvc.GetByID(incidentID); err != nil {
		return nil, err
	}
	var failedIds []uint
	var successfulIds []uint
	documentsToDelete := make(map[uint][]*incidentDTOs.IncidentInvestigationDocumentDTO)

	for _, id := range ids {
		incidentInvestigation, err := iis.GetFullDetailsByID(incidentID, id, tx)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		// Check if user exist
		if len(incidentInvestigation.Users) > 0 {
			var deleteUserIds []*uint
			for _, user := range incidentInvestigation.Users {
				deleteUserIds = append(deleteUserIds, &user.ID)
			}
			if err := iis.repo.DeleteMembers(deleteUserIds, tx); err != nil {
				tx.Rollback()
				return nil, err
			}
		}

		if len(incidentInvestigation.Findings) > 0 {
			var deleteFindingIds []*uint
			for _, finding := range incidentInvestigation.Findings {
				deleteFindingIds = append(deleteFindingIds, &finding.ID)
			}
			if err := iis.repo.DeleteFindings(deleteFindingIds, tx); err != nil {
				tx.Rollback()
				return nil, err
			}
		}

		// Delete each incidentInvestigation document if documents exist
		if len(incidentInvestigation.Documents) > 0 {
			var deleteDocumentIds []*uint
			for _, doc := range incidentInvestigation.Documents {
				if doc != nil {
					documentsToDelete[id] = append(documentsToDelete[id], doc)
					deleteDocumentIds = append(deleteDocumentIds, doc.ID)
				}
			}
			if err := iis.repo.DeleteDocuments(deleteDocumentIds, tx); err != nil {
				tx.Rollback()
				return nil, err
			}
		}

		// Delete the incidentInvestigation
		if err := iis.repo.Delete(id, tx); err != nil {
			failedIds = append(failedIds, id)
			continue // Move this line inside the loop
		}
		successfulIds = append(successfulIds, id)
	}

	// Rollback and return failedIds if any failed
	if len(failedIds) > 0 {
		tx.Rollback()
		return failedIds, nil
	}

	// Deleting files should not affect the main operation result
	if len(successfulIds) > 0 {
		for incidentInvestigationID, documents := range documentsToDelete {
			for _, doc := range documents {
				storagePath := filepath.Join("incident", "investigations", strconv.FormatUint(uint64(incidentInvestigationID), 10), "documents")
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
			err := iis.activityLogSvc.LogBulkActivities(
				models.DeleteAction, // Make sure this action constant is defined appropriately
				models.IncidentInvestigation{}.TableName(),
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

func (iis *IncidentInvestigationService) validateIncidentInvestigationFormDTO(incidentInvestigationFormDTO *incidentDTOs.IncidentInvestigationFormDTO) error {
	if _, err := iis.incidentInvestigationStatusSvc.GetByID(uint(incidentInvestigationFormDTO.IncidentInvestigationStatusID)); err != nil {
		return utility.NewValidationError(fmt.Errorf("incident investigation status does not exist"), "not_exist", "incident investigation status")
	}
	//User exist
	users := incidentInvestigationFormDTO.UserIds
	// If no user are provided, simply return nil to exit the function
	if len(users) == 0 {
		return nil
	}

	// Create a map to store user IDs
	userIDs := make(map[uint]bool)

	// Check if each user exists
	for _, user := range users {
		// Check if the user ID already exists in the map
		if userIDs[*user] {
			// If the user ID already exists, return a validation error
			return utility.NewValidationError(fmt.Errorf("duplicate user ID: %d", user), "duplicate", "user")
		}

		// Check if the user exists in the database
		if _, err := iis.userSvc.GetByID(*user); err != nil {
			return utility.NewValidationError(fmt.Errorf("user does not exist for ID: %d", user), "not_exist", "user")
		}

		// Add the user ID to the map
		userIDs[*user] = true
	}

	return nil
}

func (iis *IncidentInvestigationService) parseIncidentInvestigationDates(dto *incidentDTOs.IncidentInvestigationFormDTO) (*incidentDTOs.IncidentInvestigationParsedDates, error) {
	parsed := &incidentDTOs.IncidentInvestigationParsedDates{}

	var err error
	if dto.InitiatedAt != "" {
		parsed.InitiatedAt, err = utility.ParseDate(dto.InitiatedAt)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid initiate date"), "invalid_date", "initiate at")
		}
	}

	return parsed, nil
}

func (iis *IncidentInvestigationService) validateIncidentInvestigationDocumentFormDTO(documents []*commonDTOs.DocumentFormDTO) error {
	// Check if documents array is empty
	if len(documents) == 0 {
		return utility.NewValidationError(fmt.Errorf("documents array is required"), "required", "documents")
	}

	return nil
}
