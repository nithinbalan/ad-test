package services

import (
	"errors"
	"fmt"
	"jethur/config"
	businessContinuityDTOs "jethur/grc/business-continuity/dtos"
	TestAndExerciseTestResultRepositories "jethur/grc/business-continuity/repositories"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	"jethur/utility"
	"os"
	"path/filepath"
	"strconv"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type TestAndExerciseTestResultServiceInterface interface {
	GetFullDetailsByID(testAndExerciseId uint, tx ...utility.DatabaseTransaction) (*businessContinuityDTOs.TestAndExerciseTestResultDetailDTO, error)
	GetByID(testAndExerciseId uint) (*models.TestAndExerciseTestResult, error)
	Create(testAndExerciseId uint, currentUserID uint, businessContinuityDTOs *businessContinuityDTOs.TestAndExerciseTestResultFormDTO) (uint, error)
	Update(testAndExerciseId uint, currentUserID uint, TestAndExerciseTestResult *businessContinuityDTOs.TestAndExerciseTestResultFormDTO) error
	Delete(testAndExerciseId uint, currentUserID uint) (uint, error)

	GetDocumentPath(testAndExerciseId uint, token string) (string, error)
}

// Service and implementations
type TestAndExerciseTestResultService struct {
	repo               TestAndExerciseTestResultRepositories.TestAndExerciseTestResultRepositoryInterface
	activityLogSvc     logServices.ActivityLogServiceInterface
	userSvc            organizationServices.UserServiceInterface
	testAndExerciseSvc TestAndExerciseServiceInterface
}

func NewTestAndExerciseTestResultService(
	repo TestAndExerciseTestResultRepositories.TestAndExerciseTestResultRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	testAndExerciseSvc TestAndExerciseServiceInterface,
) TestAndExerciseTestResultServiceInterface {
	return &TestAndExerciseTestResultService{
		repo:               repo,
		activityLogSvc:     activityLogSvc,
		userSvc:            userSvc,
		testAndExerciseSvc: testAndExerciseSvc,
	}
}

func (tetrs TestAndExerciseTestResultService) GetByID(testAndExerciseId uint) (*models.TestAndExerciseTestResult, error) {
	TestAndExerciseTestResult, err := tetrs.repo.GetByID(testAndExerciseId)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, testAndExerciseId)
		}
		return nil, err
	}

	return TestAndExerciseTestResult, nil
}

func (tetrs TestAndExerciseTestResultService) GetFullDetailsByID(testAndExerciseId uint, tx ...utility.DatabaseTransaction) (*businessContinuityDTOs.TestAndExerciseTestResultDetailDTO, error) {
	// Initialize your DTO
	var testAndExerciseTestResultDTO businessContinuityDTOs.TestAndExerciseTestResultDetailDTO

	testAndExerciseTestResult, err := tetrs.repo.GetFullDetailsByID(testAndExerciseId, tx...)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, testAndExerciseId)
		}
		return nil, err
	}

	// Map the result to your DTO
	err = mapstructure.Decode(testAndExerciseTestResult, &testAndExerciseTestResultDTO)
	if err != nil {
		return nil, err
	}

	if testAndExerciseTestResult.TestCompletionDate != nil {
		testCompletionDate := testAndExerciseTestResult.TestCompletionDate.Format("2006-01-02")
		testAndExerciseTestResultDTO.TestCompletionDateString = testCompletionDate
	}

	if testAndExerciseTestResult.TestAndExercise.TestingDate != nil {
		testingDate := testAndExerciseTestResult.TestAndExercise.TestingDate.Format("2006-01-02")
		testAndExerciseTestResultDTO.TestAndExercise.TestingDateString = testingDate
	}

	testAndExerciseTestResultDTO.CreatedAt = testAndExerciseTestResult.CreatedAt
	testAndExerciseTestResultDTO.UpdatedAt = testAndExerciseTestResult.UpdatedAt

	return &testAndExerciseTestResultDTO, nil
}

func (tetrs *TestAndExerciseTestResultService) Create(testAndExerciseId uint, currentUserID uint, testAndExerciseTestResultFormDTO *businessContinuityDTOs.TestAndExerciseTestResultFormDTO) (uint, error) {
	tx := tetrs.repo.BeginTransaction()

	dates, err := tetrs.parseTestAndExerciseTestResultDates(testAndExerciseTestResultFormDTO)
	if err != nil {
		return 0, err
	}

	TestAndExerciseTestResult := &models.TestAndExerciseTestResult{
		TestAndExerciseID:             testAndExerciseId,
		TestCompletionDate:            dates.TestCompletionDate,
		CommunicationFeedbackOptionID: *testAndExerciseTestResultFormDTO.CommunicationFeedbackOptionID,
		RecoveryFeedbackOptionID:      *testAndExerciseTestResultFormDTO.RecoveryFeedbackOptionID,
		TestResultOptionID:            *testAndExerciseTestResultFormDTO.TestResultOptionID,
		TestObservations:              testAndExerciseTestResultFormDTO.TestObservations,
		Recommendations:               testAndExerciseTestResultFormDTO.Recommendations,
		Findings:                      testAndExerciseTestResultFormDTO.Findings,
		CreatedBYID:                   currentUserID,
	}

	id, err := tetrs.repo.Create(TestAndExerciseTestResult, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	for _, testAndExerciseProcess := range testAndExerciseTestResultFormDTO.TestAndExerciseProcesses {

		existingProcess, err := tetrs.testAndExerciseSvc.GetTestAndExerciseProcessByID(testAndExerciseProcess.ID)
		if err != nil {
			tx.Rollback()
			return 0, err
		}

		existingProcess.ActualRto = testAndExerciseProcess.ActualRto
		existingProcess.ActualRtoUnitId = testAndExerciseProcess.ActualRtoUnitId
		existingProcess.ActualMtpd = testAndExerciseProcess.ActualMtpd
		existingProcess.ActualMtpdUnitId = testAndExerciseProcess.ActualMtpdUnitId
		existingProcess.ActualRpo = testAndExerciseProcess.ActualRpo
		existingProcess.ActualRpoUnitId = testAndExerciseProcess.ActualRpoUnitId
		existingProcess.TimelinessID = testAndExerciseProcess.TimelinessID
		existingProcess.IsRecovered = testAndExerciseProcess.IsRecovered
		existingProcess.UpdatedBYID = &currentUserID

		err = tetrs.repo.UpdateTestAndExerciseProcess(existingProcess, tx)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	var documents []*models.TestAndExerciseTestResultDocument

	// Iterate over each document in the array
	for _, doc := range testAndExerciseTestResultFormDTO.Documents {
		// Encrypt the filename
		dataToEncrypt := []byte(doc.FileName)
		encryptedFilename, err := utility.EncryptData(dataToEncrypt)
		if err != nil {
			tx.Rollback()
			return 0, err
		}

		// Construct the URL for the document
		url := utility.GenerateDocumentURL("business-continuity", "test-and-exercise", strconv.FormatUint(uint64(testAndExerciseId), 10), "test-result", strconv.FormatUint(uint64(id), 10), "documents")

		// Create a new ComplianceRequirementStatusUpdateDocument instance
		document := &models.TestAndExerciseTestResultDocument{
			TestAndExerciseTestResultID: id,
			URL:                         url,
			Name:                        doc.FileName,
			Size:                        doc.Size,
			Token:                       encryptedFilename,
			OriginalName:                doc.OriginalName,
			Type:                        doc.Type,
			CreatedByID:                 currentUserID,
		}
		// Append the document to the slice
		documents = append(documents, document)

	}
	if len(documents) > 0 {
		// Save the ComplianceRequirementStatusUpdateDocument
		if err := tetrs.repo.CreateDocuments(documents, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
		//move file to storage
		for _, doc := range documents {
			storagePath := filepath.Join("business-continuity", "test-and-exercise", strconv.FormatUint(uint64(testAndExerciseId), 10), "test-result", strconv.FormatUint(uint64(id), 10), "documents")
			if err := utility.StoreFile(doc.Name, storagePath); err != nil {
				tx.Rollback()
				return 0, err
			}
		}

	}

	tx.Commit()
	go func() {
		existingTestAndExerciseTestResult, err := tetrs.GetFullDetailsByID(testAndExerciseId)
		if err != nil {
			utility.LogError("Failed to fetch test and exercise", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		err = tetrs.activityLogSvc.LogActivity(
			models.CreateAction,
			models.TestAndExerciseTestResult{}.TableName(),
			id,
			existingTestAndExerciseTestResult,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.TestAndExerciseTestResult{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()
	return id, nil
}
func (tetrs *TestAndExerciseTestResultService) Update(testAndExerciseId uint, currentUserID uint, testAndExerciseTestResultFormDTO *businessContinuityDTOs.TestAndExerciseTestResultFormDTO) error {
	tx := tetrs.repo.BeginTransaction()

	existingTestAndExerciseTestResult, err := tetrs.repo.GetByID(testAndExerciseId)
	if err != nil {
		return utility.NewResourceNotFoundError(err, testAndExerciseId)
	}

	oldTestAndExerciseTestResult, err := tetrs.GetFullDetailsByID(testAndExerciseId, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	dates, err := tetrs.parseTestAndExerciseTestResultDates(testAndExerciseTestResultFormDTO)
	if err != nil {
		return err
	}

	existingTestAndExerciseTestResult.TestAndExerciseID = testAndExerciseId
	existingTestAndExerciseTestResult.TestCompletionDate = dates.TestCompletionDate
	existingTestAndExerciseTestResult.CommunicationFeedbackOptionID = *testAndExerciseTestResultFormDTO.CommunicationFeedbackOptionID
	existingTestAndExerciseTestResult.RecoveryFeedbackOptionID = *testAndExerciseTestResultFormDTO.RecoveryFeedbackOptionID
	existingTestAndExerciseTestResult.TestResultOptionID = *testAndExerciseTestResultFormDTO.TestResultOptionID
	existingTestAndExerciseTestResult.TestObservations = testAndExerciseTestResultFormDTO.TestObservations
	existingTestAndExerciseTestResult.Recommendations = testAndExerciseTestResultFormDTO.Recommendations
	existingTestAndExerciseTestResult.Findings = testAndExerciseTestResultFormDTO.Findings
	existingTestAndExerciseTestResult.UpdatedBYID = &currentUserID

	err = tetrs.repo.Update(existingTestAndExerciseTestResult, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	for _, testAndExerciseProcess := range testAndExerciseTestResultFormDTO.TestAndExerciseProcesses {

		existingProcess, err := tetrs.testAndExerciseSvc.GetTestAndExerciseProcessByID(testAndExerciseProcess.ID)
		if err != nil {
			tx.Rollback()
			return err
		}

		existingProcess.ActualRto = testAndExerciseProcess.ActualRto
		existingProcess.ActualRtoUnitId = testAndExerciseProcess.ActualRtoUnitId
		existingProcess.ActualMtpd = testAndExerciseProcess.ActualMtpd
		existingProcess.ActualMtpdUnitId = testAndExerciseProcess.ActualMtpdUnitId
		existingProcess.ActualRpo = testAndExerciseProcess.ActualRpo
		existingProcess.ActualRpoUnitId = testAndExerciseProcess.ActualRpoUnitId
		existingProcess.TimelinessID = testAndExerciseProcess.TimelinessID
		existingProcess.IsRecovered = testAndExerciseProcess.IsRecovered
		existingProcess.UpdatedBYID = &currentUserID

		err = tetrs.repo.UpdateTestAndExerciseProcess(existingProcess, tx)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	var documents []*models.TestAndExerciseTestResultDocument
	var deleteDocumentIds []*uint
	var documentsToDelete []*models.TestAndExerciseTestResultDocument

	for _, doc := range testAndExerciseTestResultFormDTO.Documents {

		if doc.IsNewFile != nil && *doc.IsNewFile {
			// Encrypt the filename
			dataToEncrypt := []byte(doc.FileName)
			encryptedFileName, err := utility.EncryptData(dataToEncrypt)

			if err != nil {
				tx.Rollback()
				return err
			}

			//genereate URL
			url := utility.GenerateDocumentURL("business-continuity", "test-and-exercise", strconv.FormatUint(uint64(testAndExerciseId), 10), "test-result", strconv.FormatUint(uint64(existingTestAndExerciseTestResult.ID), 10), "documents")
			// Create a new Control test Document instance
			document := &models.TestAndExerciseTestResultDocument{
				TestAndExerciseTestResultID: existingTestAndExerciseTestResult.ID,
				URL:                         url,
				Name:                        doc.FileName,
				Size:                        doc.Size,
				Token:                       encryptedFileName,
				OriginalName:                doc.OriginalName,
				Type:                        doc.Type,
				CreatedByID:                 currentUserID,
			}

			// Append the document to the slice
			documents = append(documents, document)

		} else if doc.IsDeletedFile != nil && *doc.IsDeletedFile {
			if doc.ID != nil {
				existingDocument, err := tetrs.repo.GetDocumentByID(*doc.ID)
				if err != nil {
					if errors.Is(err, gorm.ErrRecordNotFound) {
						return utility.NewResourceNotFoundError(err, existingTestAndExerciseTestResult.ID)
					}
					return err
				}

				documentsToDelete = append(documentsToDelete, existingDocument)
				deleteDocumentIds = append(deleteDocumentIds, &existingDocument.ID)
			}
		}
	}

	if len(documents) > 0 {
		if err := tetrs.repo.CreateDocuments(documents, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	if len(deleteDocumentIds) > 0 {
		if err := tetrs.repo.DeleteDocuments(deleteDocumentIds, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, doc := range documents {
		storagePath := filepath.Join("business-continuity", "test-and-exercise", strconv.FormatUint(uint64(testAndExerciseId), 10), "test-result", strconv.FormatUint(uint64(existingTestAndExerciseTestResult.ID), 10), "documents")

		if err := utility.StoreFile(doc.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, document := range documentsToDelete {
		storagePath := filepath.Join("business-continuity", "test-and-exercise", strconv.FormatUint(uint64(testAndExerciseId), 10), "test-result", strconv.FormatUint(uint64(existingTestAndExerciseTestResult.ID), 10), "documents")
		if err := utility.RemoveFile(document.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()
	go func() {
		existingTestAndExerciseTestResult, err := tetrs.GetFullDetailsByID(testAndExerciseId)
		if err != nil {
			return
		}

		changedData := make(map[string]interface{})

		//log changes
		oldData := map[string]interface{}{
			"communication_feedback_option": oldTestAndExerciseTestResult.CommunicationFeedbackOption.TitleEn,
			"recovery_feedback_option":      oldTestAndExerciseTestResult.RecoveryFeedbackOption.TitleEn,
			"test_result_option":            oldTestAndExerciseTestResult.TestResultOption.TitleEn,
			"test_observationd":             oldTestAndExerciseTestResult.TestObservations,
			"recommendations":               oldTestAndExerciseTestResult.Recommendations,
			"findings":                      oldTestAndExerciseTestResult.Findings,
			"test_completion_date":          oldTestAndExerciseTestResult.TestCompletionDateString,
		}

		newData := map[string]interface{}{
			"communication_feedback_option": existingTestAndExerciseTestResult.CommunicationFeedbackOption.TitleEn,
			"recovery_feedback_option":      existingTestAndExerciseTestResult.RecoveryFeedbackOption.TitleEn,
			"test_result_option":            existingTestAndExerciseTestResult.TestResultOption.TitleEn,
			"test_observationd":             existingTestAndExerciseTestResult.TestObservations,
			"recommendations":               existingTestAndExerciseTestResult.Recommendations,
			"findings":                      existingTestAndExerciseTestResult.Findings,
			"test_completion_date":          existingTestAndExerciseTestResult.TestCompletionDateString,
		}
		// Compare other fields and add them to changedData
		fieldChanges := utility.CompareFieldChanges(oldData, newData)
		for k, v := range fieldChanges {
			changedData[k] = v
		}
		// Log activity if there are any changes
		if len(changedData) > 0 {
			err := tetrs.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.TestAndExerciseTestResult{}.TableName(),
				existingTestAndExerciseTestResult.ID,
				changedData,
				currentUserID,
				nil,
			)

			if err != nil {
				fields := map[string]interface{}{
					"action":    models.UpdateAction,
					"tableName": models.TestAndExerciseTestResult{}.TableName(),
					"recordID":  existingTestAndExerciseTestResult.ID,
					"userID":    currentUserID,
					"error":     err.Error(),
				}

				utility.LogError("Failed to log activity", fields)
			}
		}
	}()

	return nil
}

// func (tetrs *TestAndExerciseTestResultService) Delete(testAndExerciseId uint, ids []uint, currentUserID uint) ([]uint, error) {
// 	tx := tetrs.repo.BeginTransaction()
// 	var failedIds []uint
// 	var successfulIds []uint
// 	var deletedTestAndExerciseTestResult []map[string]interface{}

// 	for _, id := range ids {
// 		testAndExercise, err := tetrs.GetFullDetailsByID(testAndExerciseId, tx)
// 		if err != nil {
// 			tx.Rollback()
// 			return nil, err
// 		}

// 		deletedTestAndExerciseTestResult = append(deletedTestAndExerciseTestResult, map[string]interface{}{
// 			"id":                id,
// 			"test_and_exercise": testAndExercise.ID,
// 		})

// 		// Delete the test and exercise
// 		if err := tetrs.repo.Delete(id, tx); err != nil {
// 			failedIds = append(failedIds, id)
// 			continue
// 		} else {
// 			successfulIds = append(successfulIds, id)
// 		}
// 	}
// 	if len(failedIds) > 0 {
// 		tx.Rollback()
// 		return failedIds, nil
// 	}

// 	tx.Commit()

// 	// Log bulk delete activities for successful deletions
// 	if len(successfulIds) > 0 {
// 		go func() {
// 			err := tetrs.activityLogSvc.LogBulkActivities(
// 				models.DeleteAction, // Make sure this action constant is defined appropriately
// 				models.TestAndExerciseTestResult{}.TableName(),
// 				successfulIds,
// 				deletedTestAndExerciseTestResult,
// 				currentUserID,
// 				nil,
// 			)
// 			if err != nil {
// 				// @todo - save to log file
// 				fmt.Printf("Failed to log bulk delete activity: %v\n", err)
// 			}

// 		}()
// 	}

// 	return failedIds, nil
// }

func (tetrs *TestAndExerciseTestResultService) Delete(testAndExerciseId uint, currentUserID uint) (uint, error) {
	tx := tetrs.repo.BeginTransaction()
	var failedId uint
	var deletedTestResult []map[string]interface{}
	documentsToDelete := make(map[uint][]*businessContinuityDTOs.TestAndExerciseTestResultDocumentDTO)

	if _, err := tetrs.testAndExerciseSvc.GetByID(testAndExerciseId); err != nil {
		return 0, err
	}
	testAndExercise, err := tetrs.testAndExerciseSvc.GetFullDetailsByID(testAndExerciseId)

	testResult, err := tetrs.GetFullDetailsByID(testAndExerciseId, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	//Delete-> update testAndExerciseProcess

	for _, testAndExerciseProcess := range testAndExercise.TestAndExerciseProcesses {

		existingProcess, err := tetrs.testAndExerciseSvc.GetTestAndExerciseProcessByID(*testAndExerciseProcess.ID)
		if err != nil {
			tx.Rollback()
			return 0, err
		}

		existingProcess.ActualRto = nil
		existingProcess.ActualRtoUnitId = nil
		existingProcess.ActualMtpd = nil
		existingProcess.ActualMtpdUnitId = nil
		existingProcess.ActualRpo = nil
		existingProcess.ActualRpoUnitId = nil
		existingProcess.TimelinessID = nil
		existingProcess.IsRecovered = false
		existingProcess.UpdatedBYID = &currentUserID

		fmt.Println(existingProcess)

		err = tetrs.repo.UpdateTestAndExerciseProcess(existingProcess, tx)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	// Delete each incidentInvestigation document if documents exist
	if len(testResult.Documents) > 0 {
		var deleteDocumentIds []*uint
		for _, doc := range testResult.Documents {
			if doc != nil {
				documentsToDelete[testResult.ID] = append(documentsToDelete[testResult.ID], doc)
				deleteDocumentIds = append(deleteDocumentIds, &doc.ID)
			}
		}
		if err := tetrs.repo.DeleteDocuments(deleteDocumentIds, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
	}
	if err = tetrs.repo.Delete(uint(testResult.ID), tx); err != nil {
		failedId = uint(testResult.ID)
		return failedId, err
	}
	if len(documentsToDelete) > 0 {
		for testResultID, documents := range documentsToDelete {
			for _, doc := range documents {
				storagePath := filepath.Join("business-continuity", "test-and-exercise", strconv.FormatUint(uint64(testAndExerciseId), 10), "test-result", strconv.FormatUint(uint64(testResultID), 10), "documents")
				if err := utility.RemoveFile(doc.Name, storagePath); err != nil {
					tx.Rollback()
					fmt.Printf("Failed to delete file at URL %s: %v\n", doc.URL, err)
				}
			}
		}
	}
	deletedTestResult = append(deletedTestResult, map[string]interface{}{
		"id":                   testResult.ID,
		"test_and_exercise_id": testAndExerciseId,
	})
	tx.Commit()
	// Log successful delete activity
	go func() {
		err := tetrs.activityLogSvc.LogActivity(
			models.DeleteAction,
			models.IncidentInvestigation{}.TableName(),
			uint(testResult.ID),
			deletedTestResult,
			currentUserID,
			nil,
		)
		if err != nil {
			// Handle error if logging fails
			fmt.Printf("Failed to log delete activity: %v\n", err)
		}
	}()

	return failedId, nil
}

func (tetrs *TestAndExerciseTestResultService) parseTestAndExerciseTestResultDates(dto *businessContinuityDTOs.TestAndExerciseTestResultFormDTO) (*businessContinuityDTOs.TestAndExerciseTestResultParsedDates, error) {
	parsed := &businessContinuityDTOs.TestAndExerciseTestResultParsedDates{}

	var err error

	if &dto.TestCompletionDate != nil {
		parsed.TestCompletionDate, err = utility.ParseDate(dto.TestCompletionDate)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid test completion date"), "invalid_date", "test completion date")
		}
	}

	return parsed, nil
}

func (tetrs *TestAndExerciseTestResultService) GetDocumentPath(testAndExerciseID uint, token string) (string, error) {
	existingTestResult, err := tetrs.GetFullDetailsByID(testAndExerciseID)
	if err != nil {
		return "", err
	}

	if existingTestResult == nil {
		return "", errors.New("existing test result not found")
	}

	if len(existingTestResult.Documents) == 0 {
		return "", fmt.Errorf("no documents found for the test result")
	}

	decryptedFilenameBytes, err := utility.DecryptData(token)
	if err != nil {
		return "", fmt.Errorf("resource not found")
	}

	decryptedFilename := string(decryptedFilenameBytes)
	if decryptedFilename == "" {
		return "", fmt.Errorf("resource not found")
	}

	storagePath := filepath.Join("business-continuity", "test-and-exercise", strconv.FormatUint(uint64(testAndExerciseID), 10), "test-result", strconv.FormatUint(uint64(existingTestResult.ID), 10), "documents")

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
