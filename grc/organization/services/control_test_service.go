package services

import (
	"errors"
	"fmt"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	organizationDTOs "jethur/grc/organization/dtos"
	organizationRepositories "jethur/grc/organization/repositories"
	"jethur/utility"
	"path/filepath"
	"strconv"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type ControlTestServiceInterface interface {
	GetAll(controlID uint, params *organizationDTOs.ControlTestListQueryParamsDTO) (*organizationDTOs.ControlTestListResponseDTO, error)
	GetByID(id uint) (*models.ControlTest, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*organizationDTOs.ControlTestDetailsDTO, error)
	Create(controlID uint, currentUserID uint, controlTestFormDTO *organizationDTOs.ControlTestFormDTO) (uint, error)
	Update(controlID uint, id uint, currentUserID uint, controlTestFormDTO *organizationDTOs.ControlTestFormDTO) error
	Delete(controlID uint, ids []uint, currentUserID uint) ([]uint, error)
	GetTrendAnalysis(controlID uint) ([]*organizationDTOs.ControlTrendAnalysisDTO, error)
}

type ControlTestService struct {
	repo                           organizationRepositories.ControlTestRepositoryInterface
	activityLogSvc                 logServices.ActivityLogServiceInterface
	controlSvc                     ControlServiceInterface
	controlImplementationStatusSvc ControlImplementationStatusServiceInterface
	controlEfficiencyStatusSvc     ControlEfficiencyStatusServiceInterface
}

func NewControlTestService(
	repo organizationRepositories.ControlTestRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	controlSvc ControlServiceInterface,
	controlImplementationStatusSvc ControlImplementationStatusServiceInterface,
	controlEfficiencyStatusSvc ControlEfficiencyStatusServiceInterface) ControlTestServiceInterface {
	return &ControlTestService{
		repo:                           repo,
		activityLogSvc:                 activityLogSvc,
		controlSvc:                     controlSvc,
		controlImplementationStatusSvc: controlImplementationStatusSvc,
		controlEfficiencyStatusSvc:     controlEfficiencyStatusSvc,
	}
}

func (cts *ControlTestService) GetAll(controlID uint, params *organizationDTOs.ControlTestListQueryParamsDTO) (*organizationDTOs.ControlTestListResponseDTO, error) {
	_, err := cts.controlSvc.GetByID(controlID)
	if err != nil {
		return nil, err
	}

	return cts.repo.GetAll(controlID, params)
}

func (cts *ControlTestService) GetByID(id uint) (*models.ControlTest, error) {
	controlTest, err := cts.repo.GetByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return controlTest, nil
}

func (cts *ControlTestService) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*organizationDTOs.ControlTestDetailsDTO, error) {
	var controlTestDetailsDTO organizationDTOs.ControlTestDetailsDTO
	controlTest, err := cts.repo.GetFullDetailsByID(id, tx...)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	err = mapstructure.Decode(controlTest, &controlTestDetailsDTO)
	if err != nil {
		return nil, err
	}
	return &controlTestDetailsDTO, nil
}

func (cts *ControlTestService) Create(controlID uint, currentUserID uint, controlTestFormDTO *organizationDTOs.ControlTestFormDTO) (uint, error) {
	tx := cts.repo.BeginTransaction()
	if err := cts.validateControlTestFormDTO(controlTestFormDTO); err != nil {
		return 0, err
	}

	dates, err := cts.parseControlTestDates(controlTestFormDTO)
	if err != nil {
		return 0, err
	}

	control := &models.ControlTest{
		ControlID:                     controlID,
		ControlImplementationStatusID: controlTestFormDTO.ControlImplementationStatusID,
		ControlEfficiencyStatusID:     controlTestFormDTO.ControlEfficiencyStatusID,
		Title:                         controlTestFormDTO.Title,
		Description:                   controlTestFormDTO.Description,
		Score:                         controlTestFormDTO.Score,
		Comments:                      controlTestFormDTO.Comments,
		TestedSource:                  "direct",
		TestedBy:                      currentUserID,
		TestedAt:                      dates.TestedAt,
	}

	id, err := cts.repo.Create(control, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	err = cts.controlSvc.StatusUpdate(controlID, currentUserID, controlTestFormDTO, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	var documents []*models.ControlTestDocument

	for _, doc := range controlTestFormDTO.Documents {
		// Encrypt the filename
		dataToEncrypt := []byte(doc.FileName)
		encryptedFileName, err := utility.EncryptData(dataToEncrypt)
		if err != nil {
			tx.Rollback()
			return 0, err
		}

		//genereate URL
		url := utility.GenerateDocumentURL("organization", "control", strconv.FormatUint(uint64(controlID), 10), "control-test", strconv.FormatUint(uint64(id), 10), "documents", encryptedFileName)

		// Create a new Control test Document instance
		document := &models.ControlTestDocument{
			ControlTestID: id,
			URL:           url,
			Name:          doc.FileName,
			Size:          doc.Size,
			Token:         encryptedFileName,
			OriginalName:  doc.OriginalName,
			Type:          doc.Type,
		}

		// Append the document to the slice
		documents = append(documents, document)
	}

	if len(documents) > 0 {
		if err := cts.repo.CreateDocuments(documents, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	for _, doc := range documents {
		storagePath := filepath.Join("organization", "control", strconv.FormatUint(uint64(controlID), 10), "control-test", strconv.FormatUint(uint64(id), 10), "documents")
		if err := utility.StoreFile(doc.Name, storagePath); err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	tx.Commit()

	go func() {
		err := cts.activityLogSvc.LogActivity(
			models.CreateAction,
			models.ControlTest{}.TableName(),
			id,
			controlTestFormDTO,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.ControlTest{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return id, nil
}

func (cts *ControlTestService) Update(controlID uint, id uint, currentUserID uint, controlTestFormDTO *organizationDTOs.ControlTestFormDTO) error {
	tx := cts.repo.BeginTransaction()
	existingControlTest, err := cts.GetByID(id)

	if err != nil {
		return err
	}

	if err := cts.validateControlTestFormDTO(controlTestFormDTO); err != nil {
		return err
	}

	dates, err := cts.parseControlTestDates(controlTestFormDTO)
	if err != nil {
		return err
	}

	existingControlTest.ControlImplementationStatusID = controlTestFormDTO.ControlImplementationStatusID
	existingControlTest.ControlEfficiencyStatusID = controlTestFormDTO.ControlEfficiencyStatusID
	existingControlTest.Title = controlTestFormDTO.Title
	existingControlTest.Description = controlTestFormDTO.Description
	existingControlTest.Comments = controlTestFormDTO.Comments
	existingControlTest.Score = controlTestFormDTO.Score
	existingControlTest.TestedAt = dates.TestedAt
	existingControlTest.UpdatedByID = &currentUserID

	err = cts.repo.Update(existingControlTest, tx)

	if err != nil {
		tx.Rollback()
		return err
	}

	err = cts.controlSvc.StatusUpdate(controlID, currentUserID, controlTestFormDTO, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	var documents []*models.ControlTestDocument
	var deleteDocumentIds []*uint
	var documentsToDelete []*models.ControlTestDocument

	for _, doc := range controlTestFormDTO.Documents {

		if doc.IsNewFile != nil && *doc.IsNewFile {
			// Encrypt the filename
			dataToEncrypt := []byte(doc.FileName)
			encryptedFileName, err := utility.EncryptData(dataToEncrypt)
			if err != nil {
				tx.Rollback()
				return err
			}

			//genereate URL
			url := utility.GenerateDocumentURL("organization", "control", strconv.FormatUint(uint64(existingControlTest.ControlID), 10), "control-test", strconv.FormatUint(uint64(id), 10), "documents", encryptedFileName)

			// Create a new Control test Document instance
			document := &models.ControlTestDocument{
				ControlTestID: id,
				URL:           url,
				Name:          doc.FileName,
				Size:          doc.Size,
				Token:         encryptedFileName,
				OriginalName:  doc.OriginalName,
				Type:          doc.Type,
			}

			// Append the document to the slice
			documents = append(documents, document)

		} else if doc.IsDeletedFile != nil && *doc.IsDeletedFile {
			if doc.ID != nil {
				existingDocument, err := cts.repo.GetDocumentByID(*doc.ID)
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
		if err := cts.repo.CreateDocuments(documents, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	if len(deleteDocumentIds) > 0 {
		if err := cts.repo.DeleteDocuments(deleteDocumentIds, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, doc := range documents {
		storagePath := filepath.Join("organization", "control", strconv.FormatUint(uint64(existingControlTest.ControlID), 10), "control-test", strconv.FormatUint(uint64(id), 10), "documents")
		if err := utility.StoreFile(doc.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, document := range documentsToDelete {
		storagePath := filepath.Join("organization", "control", strconv.FormatUint(uint64(existingControlTest.ControlID), 10), "control-test", strconv.FormatUint(uint64(id), 10), "documents")
		fmt.Println(storagePath)
		if err := utility.RemoveFile(document.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()
	return nil
}

func (cts *ControlTestService) Delete(controlID uint, ids []uint, currentUserID uint) ([]uint, error) {
	tx := cts.repo.BeginTransaction()

	_, err := cts.controlSvc.GetByID(controlID)
	if err != nil {
		return nil, err
	}

	var failedIds []uint
	var successfulIds []uint
	documentsToDelete := make(map[uint][]*organizationDTOs.ControlTestDocumentDTO)

	for _, id := range ids {
		var deleteDocumentIds []*uint
		controlTest, err := cts.GetFullDetailsByID(id, tx)
		if err != nil {
			failedIds = append(failedIds, id)
			continue
		}

		if controlTest.Documents != nil && len(controlTest.Documents) > 0 {
			for _, document := range controlTest.Documents {
				documentsToDelete[id] = append(documentsToDelete[id], document)
				deleteDocumentIds = append(deleteDocumentIds, document.ID)
			}

			if len(deleteDocumentIds) > 0 {
				if err := cts.repo.DeleteDocuments(deleteDocumentIds, tx); err != nil {
					failedIds = append(failedIds, id)
					continue
				}
			}
		}

		if err := cts.repo.Delete(id, tx); err != nil {
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
		for controlTestId, documents := range documentsToDelete {
			for _, doc := range documents {
				storagePath := filepath.Join("organization", "control", strconv.FormatUint(uint64(controlID), 10), "control-test", strconv.FormatUint(uint64(controlTestId), 10), "documents")
				if err := utility.RemoveFile(doc.Name, storagePath); err != nil {
					tx.Rollback()
					fmt.Printf("Failed to delete file at URL %s: %v\n", *doc.URL, err)
				}
			}
		}
	}

	tx.Commit()

	go func() {
		err := cts.activityLogSvc.LogBulkActivities(
			models.DeleteAction,
			models.ControlTest{}.TableName(),
			successfulIds,
			nil,
			currentUserID,
			nil,
		)
		if err != nil {
			fmt.Printf("Failed to log bulk delete activity: %v\n", err)
		}
	}()

	return failedIds, nil
}

func (cts *ControlTestService) validateControlTestFormDTO(controlTestFormDTO *organizationDTOs.ControlTestFormDTO) error {
	if _, err := cts.controlImplementationStatusSvc.GetByID(uint(controlTestFormDTO.ControlImplementationStatusID)); err != nil {
		return utility.NewValidationError(fmt.Errorf("control implementation status does not exist"), "not_exist", "control implementation status")
	}

	if _, err := cts.controlEfficiencyStatusSvc.GetByID(uint(controlTestFormDTO.ControlEfficiencyStatusID)); err != nil {
		return utility.NewValidationError(fmt.Errorf("control efficiency status does not exist"), "not_exist", "control efficiency status")
	}

	return nil
}

func (cts *ControlTestService) parseControlTestDates(dto *organizationDTOs.ControlTestFormDTO) (*organizationDTOs.ControlTestParsedDates, error) {
	parsed := &organizationDTOs.ControlTestParsedDates{}

	var err error
	if dto.TestedAt != nil {
		parsed.TestedAt, err = utility.ParseDate(*dto.TestedAt)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid tested at date"), "invalid_date", "tested at")
		}
	}
	return parsed, nil
}

func (cts *ControlTestService) GetTrendAnalysis(controlID uint) ([]*organizationDTOs.ControlTrendAnalysisDTO, error) {
	controlTests, err := cts.repo.GetTrendAnalysis(controlID)
	if err != nil {
		return nil, err
	}

	if len(controlTests) == 0 {
		return []*organizationDTOs.ControlTrendAnalysisDTO{}, nil
	}

	statuses, err := cts.controlEfficiencyStatusSvc.GetAll("", "", "", "")
	if err != nil {
		return nil, err
	}

	var scores []*organizationDTOs.ControlTrendAnalysisDTO

	// Map the ControlTest models to ControlTrendAnalysisDTO using mapstructure
	for _, ct := range controlTests {
		if ct.Score != nil && ct.TestedAt != nil {
			var dto organizationDTOs.ControlTrendAnalysisDTO
			err := mapstructure.Decode(ct, &dto)
			if err != nil {
				return nil, err
			}
			dto.Score = float64(*ct.Score)

			testedAt := ct.TestedAt.Format("2006-01-02")
			dto.TestedAtString = testedAt
			// Map the score to its respective color code
			for _, status := range statuses {
				if dto.Score >= float64(status.ScoreFrom) && dto.Score <= float64(status.ScoreTo) {
					dto.ColorCode = status.ColorCode
					break
				}
			}
			scores = append(scores, &dto)
		}
	}

	// Reverse the slice to get the tests in chronological order
	for i, j := 0, len(scores)-1; i < j; i, j = i+1, j-1 {
		scores[i], scores[j] = scores[j], scores[i]
	}

	return scores, nil
}
