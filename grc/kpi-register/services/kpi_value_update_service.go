package services

import (
	"errors"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"jethur/config"
	commonDTOs "jethur/grc/common/dtos"
	commonServices "jethur/grc/common/services"
	kpiDTOs "jethur/grc/kpi-register/dtos"
	kpiRepositories "jethur/grc/kpi-register/repositories"
	logServices "jethur/grc/logs/services"
	masterServices "jethur/grc/masters/services"
	"jethur/grc/models"
	organizationDTOs "jethur/grc/organization/dtos"
	organizationServices "jethur/grc/organization/services"
	"jethur/utility"

	"github.com/Knetic/govaluate"
	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type KPIValueUpdateServiceInterface interface {
	GetAll(params *kpiDTOs.KPIValueUpdateListQueryParamsDTO, currentUserID uint) (*kpiDTOs.KPIValueUpdateListResponseDTO, error)
	GetFullDetailsByID(id uint) (*kpiDTOs.KPIValueUpdateDetailsDTO, error)
	Update(updateId, currentUserID uint, kpiValueUpdateFormDTO *kpiDTOs.KPIValueUpdateFormDTO) error

	GetDocumentPath(kpiValueUpdateId uint, token string) (string, error)
	GetLatestValueByKPIRegisterID(kpiRegisterID uint, tx ...utility.DatabaseTransaction) (*float64, error)
	EvaluateKPIValueUpdateInputs(kpiValueUpdateId uint, dataInputs *kpiDTOs.KPIValueUpdateEvaluationFormDTO) (*kpiDTOs.KPIValueUpdateEvaluationDetailDTO, error)
	ExportKPIValueUpdate(data []*kpiDTOs.KPIValueUpdateListDTO, currentUserID uint) (*os.File, error)

	Submit(updateId, currentUserID uint) error
	Approve(updateId, currentUserID uint, form *kpiDTOs.KPIValueUpdateCommentFormDTO) error
	Revert(updateId, currentUserID uint, form *kpiDTOs.KPIValueUpdateCommentFormDTO) error
}

// Service and implementations
type KPIValueUpdateService struct {
	repo               kpiRepositories.KPIValueUpdateRepositoryInterface
	activityLogSvc     logServices.ActivityLogServiceInterface
	kpiRegisterSvc     KPIRegisterServiceInterface
	reviewFrequencySvc masterServices.ReviewFrequencyServiceInterface
	userSvc            organizationServices.UserServiceInterface
	notificationSvc    commonServices.NotificationServiceInterface
	kpiStatusSvc       KPIStatusServiceInterface
	kpiReviewStatusSvc KPIReviewStatusServiceInterface
	userSettingSvc     organizationServices.UserSettingServiceInterface
}

func NewKPIValueUpdateService(
	repo kpiRepositories.KPIValueUpdateRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	kpiRegisterSvc KPIRegisterServiceInterface,
	reviewFrequencySvc masterServices.ReviewFrequencyServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	notificationSvc commonServices.NotificationServiceInterface,
	kpiStatusSvc KPIStatusServiceInterface,
	kpiReviewStatusSvc KPIReviewStatusServiceInterface,
	userSettingSvc organizationServices.UserSettingServiceInterface,
) KPIValueUpdateServiceInterface {
	return &KPIValueUpdateService{
		repo:               repo,
		activityLogSvc:     activityLogSvc,
		kpiRegisterSvc:     kpiRegisterSvc,
		reviewFrequencySvc: reviewFrequencySvc,
		userSvc:            userSvc,
		notificationSvc:    notificationSvc,
		kpiStatusSvc:       kpiStatusSvc,
		kpiReviewStatusSvc: kpiReviewStatusSvc,
		userSettingSvc:     userSettingSvc,
	}
}
func (kvus *KPIValueUpdateService) GetAll(params *kpiDTOs.KPIValueUpdateListQueryParamsDTO, currentUserID uint) (*kpiDTOs.KPIValueUpdateListResponseDTO, error) {

	isPassiveUser, err := kvus.userSvc.IsPassiveUser(currentUserID)
	if err != nil {
		return nil, err
	}
	params.IsPassiveUser = isPassiveUser
	return kvus.repo.GetAll(params, currentUserID)
}

func (kvus KPIValueUpdateService) GetFullDetailsByID(id uint) (*kpiDTOs.KPIValueUpdateDetailsDTO, error) {
	var valueUpdate kpiDTOs.KPIValueUpdateDetailsDTO

	existingKPIUpdate, err := kvus.repo.GetFullDetailsByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	err = mapstructure.Decode(existingKPIUpdate, &valueUpdate)
	if err != nil {
		return nil, err
	}

	if valueUpdate.ResponsibleUser == nil && existingKPIUpdate.KPIRegister != nil && existingKPIUpdate.KPIRegister.ResponsibleUser != nil {
		var respUser organizationDTOs.UserMinimalDTO
		if err := mapstructure.Decode(existingKPIUpdate.KPIRegister.ResponsibleUser, &respUser); err == nil {
			valueUpdate.ResponsibleUser = &respUser
		}
	}
	if valueUpdate.ReviewUser == nil && existingKPIUpdate.KPIRegister != nil && existingKPIUpdate.KPIRegister.ReviewUser != nil {
		var revUser organizationDTOs.UserMinimalDTO
		if err := mapstructure.Decode(existingKPIUpdate.KPIRegister.ReviewUser, &revUser); err == nil {
			valueUpdate.ReviewUser = &revUser
		}
	}

	if existingKPIUpdate.ReviewedDate != nil {
		reviewedDate := existingKPIUpdate.ReviewedDate.Format("2006-01-02")
		valueUpdate.ReviewedDateString = &reviewedDate
	}

	if existingKPIUpdate.StartDate != nil {
		startDate := existingKPIUpdate.StartDate.Format("2006-01-02")
		valueUpdate.StartDateString = &startDate
	}
	if existingKPIUpdate.EndDate != nil {
		endDate := existingKPIUpdate.EndDate.Format("2006-01-02")
		valueUpdate.EndDateString = &endDate
	}

	for i, history := range valueUpdate.ValueUpdateHistories {
		history.CreatedAt = existingKPIUpdate.ValueUpdateHistories[i].CreatedAt
	}

	valueUpdate.CreatedAt = existingKPIUpdate.CreatedAt
	valueUpdate.UpdatedAt = existingKPIUpdate.UpdatedAt

	return &valueUpdate, nil
}

func (kvus *KPIValueUpdateService) Update(updateId uint, currentUserID uint, kpiValueUpdateFormDTO *kpiDTOs.KPIValueUpdateFormDTO) error {
	tx := kvus.repo.BeginTransaction()

	oldKpiValueUpdate, err := kvus.GetFullDetailsByID(updateId)
	if err != nil {
		return err
	}

	existingKPIRegister, err := kvus.kpiRegisterSvc.GetByID(oldKpiValueUpdate.KPIRegisterID)
	if err != nil {
		tx.Rollback()
		return err
	}
	existingKPIRegisterGetFullDetails, err := kvus.kpiRegisterSvc.GetFullDetailsByID(oldKpiValueUpdate.KPIRegisterID, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	existingKPIUpdate, err := kvus.repo.GetByID(updateId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utility.NewResourceNotFoundError(err, updateId)
		}
		return err
	}
	kpiReviewStatus, err := kvus.kpiReviewStatusSvc.GetByType("updated")
	if err != nil {
		return err
	}
	calculatedValue := kpiValueUpdateFormDTO.Value
	inputMap := make(map[string]float64)

	if existingKPIRegister.CalculationType.Type == "formula" {

		for _, input := range kpiValueUpdateFormDTO.DataInputs {
			if input != nil && input.Variable != "" {
				inputMap[strings.ToUpper(input.Variable)] = input.Value
			}
		}

		// 3. Calculate value in API
		calculatedValue, err = CalculateKPIValue(
			*existingKPIRegister.Formula,
			inputMap,
		)
		if err != nil {
			return err
		}

	}
	// Update KPI Value Update
	// existingKPIUpdate.KPIRegisterID = kpiRegisterID
	reviewedDate := time.Now()
	existingKPIUpdate.ReviewedDate = &reviewedDate
	existingKPIUpdate.UnitID = *existingKPIRegister.UnitID
	existingKPIUpdate.Value = &calculatedValue
	existingKPIUpdate.Remarks = kpiValueUpdateFormDTO.Remarks
	existingKPIUpdate.KPIReviewStatusID = &kpiReviewStatus.ID
	existingKPIUpdate.CorrectiveAction = kpiValueUpdateFormDTO.CorrectiveAction
	if existingKPIUpdate.ResponsibleUserID == nil {
		if existingKPIRegister.ResponsibleUserID != nil {
			existingKPIUpdate.ResponsibleUserID = existingKPIRegister.ResponsibleUserID
		} else {
			existingKPIUpdate.ResponsibleUserID = &currentUserID
		}
	}
	if existingKPIUpdate.ReviewUserID == nil && existingKPIRegister.ReviewUserID != nil {
		existingKPIUpdate.ReviewUserID = existingKPIRegister.ReviewUserID
	}
	existingKPIUpdate.UpdatedBYID = &currentUserID

	err = kvus.repo.Update(existingKPIUpdate, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	valueUpdateHistory := &models.KPIValueUpdateHistory{
		KPIValueUpdateID: existingKPIUpdate.ID,
		Value:            &calculatedValue,
		UnitID:           existingKPIRegister.UnitID,
		Remarks:          kpiValueUpdateFormDTO.Remarks,
		CorrectiveAction: kpiValueUpdateFormDTO.CorrectiveAction,
		CreatedBYID:      &currentUserID,
	}

	if err := kvus.repo.CreateKPIValueUpdateHistory(valueUpdateHistory, tx); err != nil {
		tx.Rollback()
		return err
	}
	// keep old actual value as previous value before updating
	if existingKPIRegister.ActualValue != nil {
		oldActualValue := *existingKPIRegister.ActualValue
		existingKPIRegister.PreviousValue = &oldActualValue
	} else {
		existingKPIRegister.PreviousValue = nil
	}

	actual := calculatedValue
	existingKPIRegister.ActualValue = &actual
	existingKPIRegister.LastReviewNote = kpiValueUpdateFormDTO.Remarks

	// =========================
	// Achievement Calculation
	// =========================

	if existingKPIRegister.TargetValue != nil {

		target := *existingKPIRegister.TargetValue

		// ✅ Safe logic (avoid nil panic)
		logic := "higher"

		if existingKPIRegisterGetFullDetails.KPIAchievementLogic != nil {
			logic = existingKPIRegisterGetFullDetails.KPIAchievementLogic.Type
		}
		achievement := kvus.CalculateAchievementPercent(
			actual,
			target,
			logic,
		)

		existingKPIRegister.Achievement = &achievement

		// =========================
		// Status Mapping
		// =========================

		statusType := kvus.GetKPIStatusType(achievement)

		kpiStatus, err := kvus.kpiStatusSvc.GetByType(statusType)
		if err != nil {
			tx.Rollback()
			return err
		}

		existingKPIRegister.KPIStatusID = kpiStatus.ID

	} else {
		existingKPIRegister.Achievement = nil
	}
	existingKPIRegister.KPIReviewStatusID = &kpiReviewStatus.ID

	err = kvus.kpiRegisterSvc.UpdateValue(existingKPIRegister, tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	// initiativeIDs, err := kvus.repo.GetStrategicInitiativeIDsByKPIRegisterID(kpiRegisterID, tx)
	// if err != nil {
	// 	tx.Rollback()
	// 	return err
	// }

	// for _, initiativeID := range initiativeIDs {
	// 	if err := kvus.initiativeSvc.recalculateAverageActualValue(initiativeID, tx); err != nil {
	// 		tx.Rollback()
	// 		return err
	// 	}
	// }
	// if err := kvus.updateObjectiveScores(kpiRegisterID, tx); err != nil {
	// 	tx.Rollback()
	// 	return err
	// }

	var documents []*models.KPIValueUpdateDocument
	var deleteDocumentIds []*uint
	var documentsToDelete []*models.KPIValueUpdateDocument
	// Iterate over each document in the array
	for _, doc := range kpiValueUpdateFormDTO.Documents {
		if doc.IsNewFile != nil && *doc.IsNewFile {
			dataToEncrypt := []byte(doc.FileName)
			encryptedFilename, err := utility.EncryptData(dataToEncrypt)
			if err != nil {
				tx.Rollback()
				return err
			}

			url := utility.GenerateDocumentURL(
				"kpis",
				"kpi-register",
				strconv.FormatUint(uint64(existingKPIRegister.KPICategoryID), 10),
				"updates",
				strconv.FormatUint(uint64(updateId), 10),
				"documents",
			)
			// Create a new KPIValueUpdateDocument instance
			document := &models.KPIValueUpdateDocument{
				KPIValueUpdateID: updateId,
				URL:              url,
				Name:             doc.FileName,
				Size:             doc.Size,
				Token:            encryptedFilename,
				OriginalName:     doc.OriginalName,
				Type:             doc.Type,
			}

			// Append the document to the slice
			documents = append(documents, document)

		} else if doc.IsDeletedFile != nil && *doc.IsDeletedFile {
			// Handle deleted file case
			if doc.ID != nil {
				existingDocument, err := kvus.repo.GetDocumentByID(*doc.ID)
				if err != nil {
					if errors.Is(err, gorm.ErrRecordNotFound) {
						return utility.NewResourceNotFoundError(err, doc.ID)
					}
					return err
				}

				documentsToDelete = append(documentsToDelete, existingDocument)
				deleteDocumentIds = append(deleteDocumentIds, &existingDocument.ID)
			}
		}
	}
	if len(documents) > 0 {
		if err := kvus.repo.CreateDocuments(documents, tx); err != nil {
			tx.Rollback()
			return err
		}
	}
	if len(deleteDocumentIds) > 0 {
		if err := kvus.repo.DeleteDocuments(deleteDocumentIds, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, doc := range documents {
		storagePath := filepath.Join("kpis", "kpi-register", strconv.FormatUint(uint64(existingKPIRegister.ID), 10), "updates", strconv.FormatUint(uint64(updateId), 10), "documents")

		if err := utility.StoreFile(doc.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, document := range documentsToDelete {
		storagePath := filepath.Join("kpis", "kpi-register", strconv.FormatUint(uint64(existingKPIRegister.ID), 10), "updates", strconv.FormatUint(uint64(updateId), 10), "documents")

		if err := utility.RemoveFile(document.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := kvus.repo.DeleteKPIValueUpdateInputData(updateId, tx); err != nil {
		tx.Rollback()
		return err
	}

	// Prepare input data records
	var inputs []*models.KPIValueUpdateDataInput
	var valueUpdateHistoryInputs []*models.KPIValueUpdateHistoryInput
	for _, input := range kpiValueUpdateFormDTO.DataInputs {
		if input != nil {
			inputData := &models.KPIValueUpdateDataInput{
				KPIValueUpdateID: updateId,
				Input:            input.Input,
				Variable:         input.Variable,
				Value:            &input.Value,
				CreatedBYID:      currentUserID,
			}
			inputs = append(inputs, inputData)

			valueHistoryInput := &models.KPIValueUpdateHistoryInput{
				KPIValueUpdateHistoryID: valueUpdateHistory.ID,
				Variable:                input.Variable,
				Value:                   input.Value,
				Input:                   input.Input,
				CreatedBYID:             &currentUserID,
			}

			valueUpdateHistoryInputs = append(valueUpdateHistoryInputs, valueHistoryInput)
		}
	}

	if len(inputs) > 0 {
		if err := kvus.repo.CreateKPIValueUpdateInputData(inputs, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	if len(valueUpdateHistoryInputs) > 0 {
		if err := kvus.repo.CreateKPIValueUpdateHistoryInputs(valueUpdateHistoryInputs, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()
	go func() {
		existingKpiValueUpdate, err := kvus.GetFullDetailsByID(updateId)
		if err != nil {
			return
		}
		changedData := make(map[string]interface{})
		// log changes
		oldData := map[string]interface{}{
			"review_date": oldKpiValueUpdate.ReviewedDateString,
			"value":       oldKpiValueUpdate.Value,
			"remarks":     oldKpiValueUpdate.Remarks,
		}

		newData := map[string]interface{}{
			"review_date": existingKpiValueUpdate.ReviewedDateString,
			"value":       existingKpiValueUpdate.Value,
			"remarks":     existingKpiValueUpdate.Remarks,
		}
		// Compare other fields and add them to changedData
		fieldChanges := utility.CompareFieldChanges(oldData, newData)
		for k, v := range fieldChanges {
			changedData[k] = v
		}
		// Log activity if there are any changes
		if len(changedData) > 0 {
			err = kvus.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.KPIValueUpdate{}.TableName(),
				updateId,
				changedData,
				currentUserID,
				"strategy", "kpi", "reviews", nil,
			)
			if err != nil {
				fields := map[string]interface{}{
					"action":    models.UpdateAction,
					"tableName": models.KPIValueUpdate{}.TableName(),
					"recordID":  updateId,
					"userID":    currentUserID,
					"error":     err.Error(),
				}

				utility.LogError("Failed to log activity", fields)
			}
		}
	}()

	return nil
}

func (kvus *KPIValueUpdateService) GetDocumentPath(kpiValueUpdateId uint, token string) (string, error) {
	// Get full details of the kpi register by ID
	existingKPIUpdate, err := kvus.GetFullDetailsByID(kpiValueUpdateId)
	if err != nil {
		return "", err
	}

	// Check if there are any associated documents
	if len(existingKPIUpdate.Documents) == 0 {
		return "", fmt.Errorf("no documents found for the kpi value update")
	}

	// Decrypt the token to get the file name
	decryptedFilenameBytes, err := utility.DecryptData(token)
	if err != nil {
		return "", fmt.Errorf("resource not found")
	}

	decryptedFilename := string(decryptedFilenameBytes)
	if decryptedFilename == "" {
		return "", fmt.Errorf("resource not found")
	}

	// Define the storage path
	storagePath := filepath.Join("kpis", "kpi-register", strconv.FormatUint(uint64(existingKPIUpdate.KPIRegisterID), 10), "updates", strconv.FormatUint(uint64(kpiValueUpdateId), 10), "documents")

	// Get the base storage path from environment variables
	baseStoragePath := config.GetEnv("STORAGE_PATH", "")
	fullPath := filepath.Join(baseStoragePath, storagePath, decryptedFilename)

	// Check if the file exists at the specified path
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return "", fmt.Errorf("file not found")
	} else if err != nil {
		return "", fmt.Errorf("error accessing file: %w", err)
	}

	return fullPath, nil
}

func CalculateKPIValue(formula string, inputs map[string]float64) (float64, error) {

	expression, err := govaluate.NewEvaluableExpression(formula)
	if err != nil {
		return 0, err
	}

	parameters := make(map[string]interface{})
	for k, v := range inputs {
		parameters[k] = v
	}

	// Validate all required variables exist
	for _, v := range expression.Vars() {
		if _, ok := parameters[v]; !ok {
			return 0, fmt.Errorf("missing input value for variable: %s", v)
		}
	}

	result, err := expression.Evaluate(parameters)
	if err != nil {
		return 0, err
	}

	value, ok := result.(float64)
	if !ok {
		return 0, errors.New("formula did not return numeric result")
	}

	if math.IsInf(value, 0) || math.IsNaN(value) {
		return 0, nil
	}

	return value, nil
}

func (kvus *KPIValueUpdateService) GetLatestValueByKPIRegisterID(kpiRegisterID uint, tx ...utility.DatabaseTransaction) (*float64, error) {
	return kvus.repo.GetLatestValueByKPIRegisterID(kpiRegisterID, tx...)
}

func (kvus *KPIValueUpdateService) Submit(updateId, currentUserID uint) error {
	tx := kvus.repo.BeginTransaction()

	existingKPIUpdate, err := kvus.repo.GetByID(updateId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utility.NewResourceNotFoundError(err, updateId)
		}
		return err
	}

	kpiReviewStatus, err := kvus.kpiReviewStatusSvc.GetByType("submitted")
	if err != nil {
		return err
	}

	existingKPIUpdate.KPIReviewStatusID = &kpiReviewStatus.ID
	existingKPIUpdate.UpdatedBYID = &currentUserID

	if err = kvus.repo.Update(existingKPIUpdate, tx); err != nil {
		tx.Rollback()
		return err
	}

	kpiRegister, err := kvus.kpiRegisterSvc.GetByID(existingKPIUpdate.KPIRegisterID)
	if err != nil {
		tx.Rollback()
		return err
	}
	kpiRegister.KPIReviewStatusID = &kpiReviewStatus.ID
	if err = kvus.kpiRegisterSvc.UpdateValue(kpiRegister, tx); err != nil {
		tx.Rollback()
		return err
	}

	review := &models.KPIValueUpdateReview{
		KPIValueUpdateID:  updateId,
		KPIReviewStatusID: kpiReviewStatus.ID,
		CreatedBYID:       currentUserID,
	}
	if err = kvus.repo.CreateKPIValueUpdateReview(review, tx); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	kvus.logAndNotifyKPIValueReview(updateId, currentUserID, "submitted", existingKPIUpdate.KPIRegisterID)
	return nil
}

func (kvus *KPIValueUpdateService) Approve(updateId, currentUserID uint, form *kpiDTOs.KPIValueUpdateCommentFormDTO) error {
	tx := kvus.repo.BeginTransaction()

	existingKPIUpdate, err := kvus.repo.GetByID(updateId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utility.NewResourceNotFoundError(err, updateId)
		}
		return err
	}

	kpiReviewStatus, err := kvus.kpiReviewStatusSvc.GetByType("approved")
	if err != nil {
		return err
	}

	now := time.Now()
	existingKPIUpdate.KPIReviewStatusID = &kpiReviewStatus.ID
	existingKPIUpdate.ApprovedAt = &now
	existingKPIUpdate.RevertedAt = nil
	existingKPIUpdate.UpdatedBYID = &currentUserID

	if err = kvus.repo.Update(existingKPIUpdate, tx); err != nil {
		tx.Rollback()
		return err
	}

	kpiRegister, err := kvus.kpiRegisterSvc.GetByID(existingKPIUpdate.KPIRegisterID)
	if err != nil {
		tx.Rollback()
		return err
	}
	kpiRegister.KPIReviewStatusID = &kpiReviewStatus.ID
	if err = kvus.kpiRegisterSvc.UpdateValue(kpiRegister, tx); err != nil {
		tx.Rollback()
		return err
	}

	comment := form.Comment
	review := &models.KPIValueUpdateReview{
		KPIValueUpdateID:  updateId,
		KPIReviewStatusID: kpiReviewStatus.ID,
		Comment:           &comment,
		CreatedBYID:       currentUserID,
	}
	if err = kvus.repo.CreateKPIValueUpdateReview(review, tx); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	kvus.logAndNotifyKPIValueReview(updateId, currentUserID, "approved", existingKPIUpdate.KPIRegisterID)
	return nil
}

func (kvus *KPIValueUpdateService) Revert(updateId, currentUserID uint, form *kpiDTOs.KPIValueUpdateCommentFormDTO) error {
	tx := kvus.repo.BeginTransaction()

	existingKPIUpdate, err := kvus.repo.GetByID(updateId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utility.NewResourceNotFoundError(err, updateId)
		}
		return err
	}

	kpiReviewStatus, err := kvus.kpiReviewStatusSvc.GetByType("reverted")
	if err != nil {
		return err
	}

	now := time.Now()
	existingKPIUpdate.KPIReviewStatusID = &kpiReviewStatus.ID
	existingKPIUpdate.RevertedAt = &now
	existingKPIUpdate.ApprovedAt = nil
	existingKPIUpdate.UpdatedBYID = &currentUserID

	if err = kvus.repo.Update(existingKPIUpdate, tx); err != nil {
		tx.Rollback()
		return err
	}

	kpiRegister, err := kvus.kpiRegisterSvc.GetByID(existingKPIUpdate.KPIRegisterID)
	if err != nil {
		tx.Rollback()
		return err
	}
	kpiRegister.KPIReviewStatusID = &kpiReviewStatus.ID
	if err = kvus.kpiRegisterSvc.UpdateValue(kpiRegister, tx); err != nil {
		tx.Rollback()
		return err
	}

	comment := form.Comment
	review := &models.KPIValueUpdateReview{
		KPIValueUpdateID:  updateId,
		KPIReviewStatusID: kpiReviewStatus.ID,
		Comment:           &comment,
		CreatedBYID:       currentUserID,
	}
	if err = kvus.repo.CreateKPIValueUpdateReview(review, tx); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	kvus.logAndNotifyKPIValueReview(updateId, currentUserID, "reverted", existingKPIUpdate.KPIRegisterID)
	return nil
}

func (kvus *KPIValueUpdateService) logAndNotifyKPIValueReview(updateId uint, currentUserID uint, action string, kpiRegisterID uint) {
	go func() {
		// Log Activity
		err := kvus.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.KPIValueUpdate{}.TableName(),
			updateId,
			nil,
			currentUserID,
			"key_performance_indicator", "kpi_review", "reviews", nil,
		)
		if err != nil {
			utility.LogError("Failed to log activity", map[string]interface{}{
				"action":    models.UpdateAction,
				"tableName": models.KPIValueUpdate{}.TableName(),
				"recordID":  updateId,
				"userID":    currentUserID,
				"error":     err.Error(),
			})
		}

		// Notify
		kpiRegister, err := kvus.kpiRegisterSvc.GetFullDetailsByID(kpiRegisterID)
		if err != nil {
			utility.LogError("Failed to fetch kpi register for notification", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		var targetUser *organizationDTOs.UserMinimalDTO
		var headingEn, headingAr, emailSubject string
		var actionString string
		var templatePath string

		switch action {
		case "submitted":
			targetUser = kpiRegister.ReviewUser
			headingEn = "KPI Review Submitted"
			headingAr = "تم تقديم مراجعة مؤشرات الأداء الرئيسية"
			emailSubject = fmt.Sprintf("KPI Review Submitted – %s %s", utility.GetStringValue(kpiRegister.ReferenceNumber), kpiRegister.Title)
			actionString = "submitted"
			templatePath = "key-performance-indicator/templates/email/kpi_value_update_submit_mail.html"
		case "approved":
			targetUser = kpiRegister.ResponsibleUser
			headingEn = "KPI Review Approved"
			headingAr = "تمت الموافقة على مراجعة مؤشرات الأداء الرئيسية"
			emailSubject = fmt.Sprintf("KPI Review Approved – %s %s", utility.GetStringValue(kpiRegister.ReferenceNumber), kpiRegister.Title)
			actionString = "approved"
			templatePath = "key-performance-indicator/templates/email/kpi_value_update_approve_mail.html"
		case "reverted":
			targetUser = kpiRegister.ResponsibleUser
			headingEn = "KPI Review Reverted"
			headingAr = "تم التراجع عن مراجعة مؤشرات الأداء الرئيسية"
			emailSubject = fmt.Sprintf("KPI Review Reverted – %s %s", utility.GetStringValue(kpiRegister.ReferenceNumber), kpiRegister.Title)
			actionString = "reverted"
			templatePath = "key-performance-indicator/templates/email/kpi_value_update_revert_mail.html"
		}

		if targetUser == nil {
			return
		}

		urlParam := fmt.Sprintf(`{"kpi_value_update_id":%d, "kpi_register_id":%d}`, updateId, kpiRegisterID)
		entity := "kpi-value-update"

		dbNotificationDTO := &commonDTOs.DBNotificationDTO{
			UserIDs:              []uint{uint(targetUser.ID)},
			ModuleIdentifier:     "kpi_review",
			NotificationType:     string(commonDTOs.ReviewType),
			HeadingEn:            headingEn,
			HeadingAr:            headingAr,
			TitleEn:              kpiRegister.Title,
			TitleAr:              kpiRegister.Title,
			Entity:               &entity,
			Action:               &actionString,
			URLParam:             urlParam,
			QueryParam:           "",
			IsSystemNotification: false,
			CreatedByID:          &currentUserID,
		}

		appURL := utility.GetAppURL()
		redirectURL := fmt.Sprintf("%s/key-performance-indicator/kpi-review/%d/overview", appURL, updateId)

		reviewPeriod := "N/A"
		if kpiRegister.StartDateString != nil && kpiRegister.EndDateString != nil {
			reviewPeriod = fmt.Sprintf("%s to %s", *kpiRegister.StartDateString, *kpiRegister.EndDateString)
		}

		targetDate := "N/A"
		if kpiRegister.EndDateString != nil {
			targetDate = *kpiRegister.EndDateString
		}

		emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
			EmailTo:              []string{targetUser.Email},
			Subject:              emailSubject,
			TemplateRelativePath: templatePath,
			EmailData: map[string]interface{}{
				"FirstNameEn":     targetUser.FirstNameEN,
				"LastNameEn":      targetUser.LastNameEN,
				"ReferenceNumber": utility.GetStringValue(kpiRegister.ReferenceNumber),
				"Title":           kpiRegister.Title,
				"ReviewPeriod":    reviewPeriod,
				"StartDate":       kpiRegister.StartDateString,
				"EndDate":         kpiRegister.EndDateString,
				"TargetDate":      targetDate,
				"RedirectURL":     redirectURL,
			},
		}

		err = kvus.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "kpi_register")
		if err != nil {
			utility.LogError("Failed to send KPI review notification", map[string]interface{}{
				"error": err.Error(),
			})
		}
	}()
}

func (kvus *KPIValueUpdateService) EvaluateKPIValueUpdateInputs(kpiValueUpdateId uint, dataInputs *kpiDTOs.KPIValueUpdateEvaluationFormDTO) (*kpiDTOs.KPIValueUpdateEvaluationDetailDTO, error) {
	kpiValueUpdate, err := kvus.GetFullDetailsByID(kpiValueUpdateId)
	if err != nil {
		return nil, err
	}

	kpiRegister, err := kvus.kpiRegisterSvc.GetFullDetailsByID(kpiValueUpdate.KPIRegisterID)
	if err != nil {
		return nil, err
	}
	inputMap := make(map[string]float64)

	for _, input := range dataInputs.DataInputs {
		if input.Variable != "" {
			inputMap[strings.ToUpper(input.Variable)] = input.Value
		}
	}

	var calculatedValue float64

	if kpiRegister.CalculationType.Type == "formula" {
		calculatedValue, err = CalculateKPIValue(
			*kpiRegister.Formula,
			inputMap,
		)
		if err != nil {
			return nil, err
		}
	}

	return &kpiDTOs.KPIValueUpdateEvaluationDetailDTO{
		DataInputs: dataInputs.DataInputs,
		Value:      calculatedValue,
	}, nil
}

func (kvus *KPIValueUpdateService) CalculateAchievementPercent(actual, target float64, logic string) float64 {

	if target == 0 && actual == 0 {
		return 100
	}

	switch logic {

	case "higher":
		if target == 0 {
			return 150
		}
		return math.Min((actual/target)*100, 150)

	case "lower":
		if actual == 0 {
			return 150
		}
		return math.Min((target/actual)*100, 150)

	default:
		return 0
	}
}

func (kvus *KPIValueUpdateService) GetKPIStatusType(achievement float64) string {

	switch {
	case achievement > 110:
		return "exceeded"
	case achievement >= 95:
		return "achieved"
	case achievement >= 85:
		return "near-target"
	case achievement >= 70:
		return "needs-improvement"
	default:
		return "low"
	}
}

func (kvus *KPIValueUpdateService) ExportKPIValueUpdate(data []*kpiDTOs.KPIValueUpdateListDTO, currentUserID uint) (*os.File, error) {
	userSetting, err := kvus.userSettingSvc.GetUserSettingsByUserID(currentUserID)
	if err != nil {
		return nil, err
	}

	var headers []string

	if userSetting.Language == "ar" {
		headers = []string{
			"SL No.",
			"مؤشر الأداء",
			"النطاق",
			"تاريخ المراجعة",
			"القيمة",
			"ملاحظات",
			"الفترة",
			"تاريخ البداية",
			"تاريخ النهاية",
			"تمت مراجعته بواسطة",
			"حالة المراجعة",
			"تاريخ الإنشاء",
		}
	} else {
		headers = []string{
			"SL No.",
			"KPI Register",
			"Domain",
			"Reviewed Date",
			"Value",
			"Remarks",
			"Period",
			"Start Date",
			"End Date",
			"Reviewed By",
			"Review Status",
			"Created By",
		}
	}

	rows := make([][]interface{}, len(data))

	for i, item := range data {
		var (
			category     string
			reviewStatus string
			reviewedBy   string
			createdBy    string
		)

		if userSetting.Language == "ar" {
			category = item.KPICategoryTitleAr
			reviewStatus = item.KPIReviewStatusTitleAr
			reviewedBy = utility.GetStringValue(item.ReviewUserFirstNameAr) + " " + utility.GetStringValue(item.ReviewUserLastNameAr)
			createdBy = item.CreatedUserFirstNameAr + " " + utility.GetStringValue(item.CreatedUserLastNameAr)
		} else {
			category = item.KPICategoryTitleEn
			reviewStatus = item.KPIReviewStatusTitleEn
			reviewedBy = utility.GetStringValue(item.ReviewUserFirstNameEn) + " " + utility.GetStringValue(item.ReviewUserLastNameEn)
			createdBy = item.CreatedUserFirstNameEn + " " + item.CreatedUserLastNameEn
		}

		rows[i] = []interface{}{
			i + 1,
			item.KPIRegisterTitle,
			category,
			item.ReviewedDate,
			floatPointerToString(item.Value),
			item.Remarks,
			utility.GetStringValue(item.PeriodLabel),
			utility.GetStringValue(item.StartDate),
			utility.GetStringValue(item.EndDate),
			reviewedBy,
			reviewStatus,
			createdBy,
		}
	}

	return utility.GenerateExcelFile(headers, rows)
}
