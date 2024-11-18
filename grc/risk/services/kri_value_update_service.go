package services

import (
	"errors"
	"fmt"
	"jethur/config"
	commonDTOs "jethur/grc/common/dtos"
	commonServices "jethur/grc/common/services"
	logServices "jethur/grc/logs/services"
	masterServices "jethur/grc/masters/services"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	riskDTOs "jethur/grc/risk/dtos"
	riskRepositories "jethur/grc/risk/repositories"
	"jethur/utility"
	"log"
	"path/filepath"
	"strconv"
	"time"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type KRIValueUpdateServiceInterface interface {
	GetAll(kriRegisterID uint, params *riskDTOs.KRIValueUpdateListQueryParamsDTO) (*riskDTOs.KRIValueUpdateListResponseDTO, error)
	GetFullDetailsByID(id uint) (*riskDTOs.KRIValueUpdateDetailsDTO, error)
	Create(kriRegisterID, currentUserID uint, kriValueUpdateFormDTO *riskDTOs.KRIValueUpdateFormDTO) (uint, error)
	Update(kriRegisterID, updateId, currentUserID uint, kriValueUpdateFormDTO *riskDTOs.KRIValueUpdateFormDTO) error
	Delete(kriRegisterID uint, ids []uint, currentUserID uint) ([]uint, error)
}

// Service and implementations
type KRIValueUpdateService struct {
	repo               riskRepositories.KRIValueUpdateRepositoryInterface
	activityLogSvc     logServices.ActivityLogServiceInterface
	kriRegisterSvc     KRIRegisterServiceInterface
	reviewFrequencySvc masterServices.ReviewFrequencyServiceInterface
	userSvc            organizationServices.UserServiceInterface
	notificationSvc    commonServices.NotificationServiceInterface
}

func NewKRIValueUpdateService(
	repo riskRepositories.KRIValueUpdateRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	kriRegisterSvc KRIRegisterServiceInterface,
	reviewFrequencySvc masterServices.ReviewFrequencyServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	notificationSvc commonServices.NotificationServiceInterface,
) KRIValueUpdateServiceInterface {
	return &KRIValueUpdateService{
		repo:               repo,
		activityLogSvc:     activityLogSvc,
		kriRegisterSvc:     kriRegisterSvc,
		reviewFrequencySvc: reviewFrequencySvc,
		userSvc:            userSvc,
		notificationSvc:    notificationSvc,
	}
}

func (kvus *KRIValueUpdateService) GetAll(kriRegisterID uint, params *riskDTOs.KRIValueUpdateListQueryParamsDTO) (*riskDTOs.KRIValueUpdateListResponseDTO, error) {
	_, err := kvus.kriRegisterSvc.GetByID(kriRegisterID)
	if err != nil {
		return nil, err
	}

	return kvus.repo.GetAll(kriRegisterID, params)
}

func (kvus *KRIValueUpdateService) Create(kriRegisterID, currentUserID uint, kriValueUpdateFormDTO *riskDTOs.KRIValueUpdateFormDTO) (uint, error) {
	tx := kvus.repo.BeginTransaction()
	existingKRIRegister, err := kvus.kriRegisterSvc.GetByID(kriRegisterID)
	if err != nil {
		return 0, err
	}

	reviewFrequency, err := kvus.reviewFrequencySvc.GetByID(existingKRIRegister.ReviewFrequencyID)
	if err != nil {
		return 0, err
	}

	nextReviewDate := utility.CalculateNextReviewDate(reviewFrequency.Type)

	dates, err := kvus.parseKRIValueUpdateDates(kriValueUpdateFormDTO)
	if err != nil {
		return 0, err
	}
	kriValueUpdate := &models.KRIValueUpdate{
		KRIRegisterID:   kriRegisterID,
		UnitID:          *existingKRIRegister.UnitID,
		ReviewedDate:    *dates.ReviewedDate,
		Value:           &kriValueUpdateFormDTO.Value,
		LowThreshold:    existingKRIRegister.LowThreshold,
		MediumThreshold: existingKRIRegister.MediumThreshold,
		HighThreshold:   existingKRIRegister.HighThreshold,
		Remarks:         kriValueUpdateFormDTO.Remarks,
		CreatedByID:     currentUserID,
	}

	id, err := kvus.repo.Create(kriValueUpdate, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	existingKRIRegister.CurrentValue = &kriValueUpdateFormDTO.Value
	existingKRIRegister.LastReviewNote = kriValueUpdateFormDTO.Remarks
	existingKRIRegister.LastReviewedDate = dates.ReviewedDate
	existingKRIRegister.NextReviewDate = nextReviewDate

	err = kvus.kriRegisterSvc.UpdateValue(existingKRIRegister, tx)
	if err != nil {
		return 0, err
	}

	var documents []*models.KRIValueUpdateDocument
	// Iterate over each document in the array
	for _, doc := range kriValueUpdateFormDTO.Documents {

		// Encrypt the filename
		dataToEncrypt := []byte(doc.FileName)
		encryptedFilename, err := utility.EncryptData(dataToEncrypt)
		if err != nil {
			tx.Rollback()
			return 0, err
		}

		// Construct the URL for the document
		url := fmt.Sprintf("%s/kri-register/%d/updates/%s", config.GetEnv("API_URL", "http://localhost:8080"), kriRegisterID, encryptedFilename)

		// Create a new KRIValueUpdateDocument instance
		document := &models.KRIValueUpdateDocument{
			KRIValueUpdateID: id,
			URL:              url,
			Name:             doc.FileName,
			Size:             doc.Size,
			Token:            encryptedFilename,
			OriginalName:     doc.OriginalName,
			Type:             doc.Type,
		}
		// Append the document to the slice
		documents = append(documents, document)

	}
	if len(documents) > 0 {
		if err := kvus.repo.CreateDocuments(documents, tx); err != nil {
			tx.Rollback()
			return 0, err
		}

		for _, doc := range kriValueUpdateFormDTO.Documents {
			storagePath := filepath.Join("risk", "kri-register", strconv.FormatUint(uint64(kriRegisterID), 10), "updates", strconv.FormatUint(uint64(id), 10), "documents")
			if err := utility.StoreFile(doc.FileName, storagePath); err != nil {
				tx.Rollback()
				return 0, err
			}
		}
	}

	tx.Commit()

	go func() {
		existingKriValueUpdate, err := kvus.GetFullDetailsByID(kriValueUpdate.ID)
		if err != nil {
			utility.LogError("Failed to fetch kri value update", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		err = kvus.activityLogSvc.LogActivity(
			models.CreateAction,
			models.KRIValueUpdate{}.TableName(),
			id,
			existingKriValueUpdate,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.KRIValueUpdate{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
		kriValue := kriValueUpdateFormDTO.Value
		if existingKRIRegister.HighThreshold != nil && kriValue > *existingKRIRegister.HighThreshold {
			// Get Risks associated with the given kri
			existingRisks, err := kvus.repo.GetRisksByKriID(kriRegisterID)
			if err != nil {
				log.Println("error :", err.Error())
			}

			for _, existingRisk := range existingRisks {
				existingKRIRegister, err := kvus.kriRegisterSvc.GetByID(kriRegisterID)
				if err != nil {
					return
				}
				// Prepare notification and email
				urlParam := fmt.Sprintf(`{"id":%d}`, kriRegisterID)
				entity := "kri-value-update"
				action := "value-update-high"
				dbNotificationDTO := &commonDTOs.DBNotificationDTO{
					UserID:               existingRisk.ResponsibleUserID,
					ModuleIdentifier:     "risk",
					NotificationType:     string(commonDTOs.AssignmentType),
					HeadingEn:            "Kri Value Update",
					HeadingAr:            "تحديث قيمة كري",
					TitleEn:              existingKRIRegister.Title,
					Entity:               &entity,
					Action:               &action,
					URLParam:             urlParam,
					IsSystemNotification: false,
					CreatedByID:          &currentUserID,
				}
				responsibleUser, err := kvus.userSvc.GetByID(existingRisk.ResponsibleUserID)
				if err != nil {
					utility.LogError("Failed to fetch responsible user", map[string]interface{}{
						"error": err.Error(),
					})
					return
				}
				emailTo := []string{responsibleUser.Email}
				subject := fmt.Sprintf("KRI Register %s: The value for the KRI has gone above the acceptable limit.", existingKRIRegister.Title)
				appURL := utility.GetAppURL()
				emailData := map[string]interface{}{
					"FirstNameEn":         responsibleUser.FirstNameEn,
					"LastNameEn":          responsibleUser.LastNameEn,
					"KRIReferenceNumber":  existingKRIRegister.ReferenceNumber,
					"KRITitle":            existingKRIRegister.Title,
					"KRIDescription":      existingKRIRegister.Description,
					"RiskReferenceNumber": existingRisk.ReferenceNumber,
					"RiskTitle":           existingRisk.Title,
					"RiskDescription":     existingRisk.Description,
					"CurrentValue":        kriValue,
					"AssignedOn":          time.Now().Format("2006-01-02"),
					"RedirectURL":         fmt.Sprintf("%srisk/kri-register/%d/overview", appURL, kriRegisterID),
				}
				emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
					EmailTo:              emailTo,
					Subject:              subject,
					TemplateRelativePath: "risk/templates/email/kri_value_update_mail.html",
					EmailData:            emailData,
				}
				err = kvus.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "kri_register")
				if err != nil {
					utility.LogError("Failed to send notifications", map[string]interface{}{
						"error": err.Error(),
					})
				}
			}

		}
	}()

	return id, nil
}

func (kvus KRIValueUpdateService) GetFullDetailsByID(id uint) (*riskDTOs.KRIValueUpdateDetailsDTO, error) {
	var valueUpdate riskDTOs.KRIValueUpdateDetailsDTO

	existingKRIUpdate, err := kvus.repo.GetFullDetailsByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	err = mapstructure.Decode(existingKRIUpdate, &valueUpdate)
	if err != nil {
		return nil, err
	}

	reviewedDate := existingKRIUpdate.ReviewedDate.Format("2006-01-02")
	valueUpdate.ReviewedDateString = &reviewedDate

	valueUpdate.CreatedAt = existingKRIUpdate.CreatedAt
	valueUpdate.UpdatedAt = existingKRIUpdate.UpdatedAt

	return &valueUpdate, nil
}

func (kvus *KRIValueUpdateService) Update(kriRegisterID, updateId uint, currentUserID uint, kriValueUpdateFormDTO *riskDTOs.KRIValueUpdateFormDTO) error {
	tx := kvus.repo.BeginTransaction()

	oldKriValueUpdate, err := kvus.GetFullDetailsByID(updateId)
	if err != nil {
		return err
	}

	existingKRIRegister, err := kvus.kriRegisterSvc.GetByID(kriRegisterID)
	if err != nil {
		tx.Rollback()
		return err
	}

	dates, err := kvus.parseKRIValueUpdateDates(kriValueUpdateFormDTO)
	if err != nil {
		return err
	}

	existingKRIUpdate, err := kvus.repo.GetByID(updateId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utility.NewResourceNotFoundError(err, updateId)
		}
		return err
	}

	// Update KRI Value Update
	existingKRIUpdate.KRIRegisterID = kriRegisterID
	existingKRIUpdate.ReviewedDate = *dates.ReviewedDate
	existingKRIUpdate.UnitID = *existingKRIRegister.UnitID
	existingKRIUpdate.Value = &kriValueUpdateFormDTO.Value
	existingKRIUpdate.Remarks = kriValueUpdateFormDTO.Remarks

	err = kvus.repo.Update(existingKRIUpdate, tx)

	if err != nil {
		tx.Rollback()
		return err
	}

	existingKRIRegister.CurrentValue = &kriValueUpdateFormDTO.Value
	existingKRIRegister.LastReviewNote = kriValueUpdateFormDTO.Remarks

	err = kvus.kriRegisterSvc.UpdateValue(existingKRIRegister, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	var documents []*models.KRIValueUpdateDocument
	var deleteDocumentIds []*uint
	var documentsToDelete []*models.KRIValueUpdateDocument
	// Iterate over each document in the array
	for _, doc := range kriValueUpdateFormDTO.Documents {
		if doc.IsNewFile != nil && *doc.IsNewFile {
			dataToEncrypt := []byte(doc.FileName)
			encryptedFilename, err := utility.EncryptData(dataToEncrypt)
			if err != nil {
				tx.Rollback()
				return err
			}
			url := fmt.Sprintf("%s/kri-register/%d/updates/%s", config.GetEnv("API_URL", "http://localhost:8080"), kriRegisterID, encryptedFilename)

			// Create a new KRIValueUpdateDocument instance
			document := &models.KRIValueUpdateDocument{
				KRIValueUpdateID: updateId,
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
		storagePath := filepath.Join("risk", "kri-register", strconv.FormatUint(uint64(kriRegisterID), 10), "updates", strconv.FormatUint(uint64(updateId), 10), "documents")

		if err := utility.StoreFile(doc.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, document := range documentsToDelete {
		storagePath := filepath.Join("risk", "kri-register", strconv.FormatUint(uint64(kriRegisterID), 10), "updates", strconv.FormatUint(uint64(updateId), 10), "documents")

		if err := utility.RemoveFile(document.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()
	go func() {
		existingKriValueUpdate, err := kvus.GetFullDetailsByID(updateId)
		if err != nil {
			return
		}
		changedData := make(map[string]interface{})
		//log changes
		oldData := map[string]interface{}{
			"review_date": oldKriValueUpdate.ReviewedDateString,
			"value":       oldKriValueUpdate.Value,
			"remarks":     oldKriValueUpdate.Remarks,
		}

		newData := map[string]interface{}{
			"review_date": existingKriValueUpdate.ReviewedDateString,
			"value":       existingKriValueUpdate.Value,
			"remarks":     existingKriValueUpdate.Remarks,
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
				models.KRIValueUpdate{}.TableName(),
				updateId,
				changedData,
				currentUserID,
				nil,
			)

			if err != nil {
				fields := map[string]interface{}{
					"action":    models.UpdateAction,
					"tableName": models.KRIValueUpdate{}.TableName(),
					"recordID":  updateId,
					"userID":    currentUserID,
					"error":     err.Error(),
				}

				utility.LogError("Failed to log activity", fields)
			}
		}
		kriValue := kriValueUpdateFormDTO.Value
		if existingKRIRegister.HighThreshold != nil && kriValue > *existingKRIRegister.HighThreshold {
			// Get Risks associated with the given kri
			existingRisks, err := kvus.repo.GetRisksByKriID(kriRegisterID)
			if err != nil {
				log.Println("error :", err.Error())
			}

			for _, existingRisk := range existingRisks {
				existingKRIRegister, err := kvus.kriRegisterSvc.GetByID(kriRegisterID)
				if err != nil {
					return
				}
				// Prepare notification and email
				urlParam := fmt.Sprintf(`{"id":%d}`, kriRegisterID)
				entity := "kri-value-update"
				action := "value-update-high"
				dbNotificationDTO := &commonDTOs.DBNotificationDTO{
					UserID:               existingRisk.ResponsibleUserID,
					ModuleIdentifier:     "risk",
					NotificationType:     string(commonDTOs.AssignmentType),
					HeadingEn:            "Kri Value Update",
					HeadingAr:            "تحديث قيمة كري",
					TitleEn:              existingKRIRegister.Title,
					Entity:               &entity,
					Action:               &action,
					URLParam:             urlParam,
					IsSystemNotification: false,
					CreatedByID:          &currentUserID,
				}
				responsibleUser, err := kvus.userSvc.GetByID(existingRisk.ResponsibleUserID)
				if err != nil {
					utility.LogError("Failed to fetch responsible user", map[string]interface{}{
						"error": err.Error(),
					})
					return
				}
				emailTo := []string{responsibleUser.Email}
				subject := fmt.Sprintf("KRI Register %s: The value for the KRI has gone above the acceptable limit.", existingKRIRegister.Title)
				appURL := utility.GetAppURL()
				emailData := map[string]interface{}{
					"FirstNameEn":         responsibleUser.FirstNameEn,
					"LastNameEn":          responsibleUser.LastNameEn,
					"KRIReferenceNumber":  existingKRIRegister.ReferenceNumber,
					"KRITitle":            existingKRIRegister.Title,
					"KRIDescription":      existingKRIRegister.Description,
					"RiskReferenceNumber": existingRisk.ReferenceNumber,
					"RiskTitle":           existingRisk.Title,
					"RiskDescription":     existingRisk.Description,
					"CurrentValue":        kriValue,
					"AssignedOn":          time.Now().Format("2006-01-02"),
					"RedirectURL":         fmt.Sprintf("%srisk/kri-register/%d/overview", appURL, kriRegisterID),
				}
				emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
					EmailTo:              emailTo,
					Subject:              subject,
					TemplateRelativePath: "risk/templates/email/kri_value_update_mail.html",
					EmailData:            emailData,
				}
				err = kvus.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "kri_register")
				if err != nil {
					utility.LogError("Failed to send notifications", map[string]interface{}{
						"error": err.Error(),
					})
				}
			}

		}
	}()

	return nil
}

func (kvus *KRIValueUpdateService) Delete(kriRegisterID uint, ids []uint, currentUserID uint) ([]uint, error) {
	tx := kvus.repo.BeginTransaction()

	existingKRIRegister, err := kvus.kriRegisterSvc.GetByID(kriRegisterID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	var failedIds []uint
	var successfulIds []uint

	var notifications []struct {
		OwnerID         *uint
		TitleEn         string
		TitleAr         string
		Description     string
		ReferenceNumber string
	}
	var deletedKriValueUpdate []map[string]interface{}

	documentsToDelete := make(map[uint][]*models.KRIValueUpdateDocument)

	for _, id := range ids {
		var deleteDocumentIds []*uint
		valueUpdate, err := kvus.repo.GetFullDetailsByID(id, tx)
		if err != nil {
			failedIds = append(failedIds, id)
			continue
		}

		if valueUpdate.Documents != nil && len(valueUpdate.Documents) > 0 {
			for _, updateDocument := range valueUpdate.Documents {
				documentsToDelete[id] = append(documentsToDelete[id], &updateDocument)
				deleteDocumentIds = append(deleteDocumentIds, &updateDocument.ID)
			}
			if len(deleteDocumentIds) > 0 {
				if err := kvus.repo.DeleteDocuments(deleteDocumentIds, tx); err != nil {
					failedIds = append(failedIds, id)
					continue
				}
			}
		}

		if err := kvus.repo.Delete(id, tx); err != nil {
			failedIds = append(failedIds, id)
			continue
		} else {
			successfulIds = append(successfulIds, id)
		}

		deletedKriValueUpdate = append(deletedKriValueUpdate, map[string]interface{}{
			"id":              id,
			"kri_register_id": kriRegisterID,
		})
	}

	if len(failedIds) > 0 {
		tx.Rollback()
		return failedIds, nil
	}

	// Deleting files should not affect the main operation result
	if len(successfulIds) > 0 {
		for updateId, documents := range documentsToDelete {
			for _, doc := range documents {
				storagePath := filepath.Join("risk", "kri-register", strconv.FormatUint(uint64(kriRegisterID), 10), "updates", strconv.FormatUint(uint64(updateId), 10), "documents")
				if err := utility.RemoveFile(doc.Name, storagePath); err != nil {
					tx.Rollback()
					fmt.Printf("Failed to delete file at URL %s: %v\n", doc.URL, err)
				}
			}
		}

		// Update KRIRegister with the next available KRIValueUpdate value
		nextValueUpdate, err := kvus.repo.GetNextKRIValueUpdate(kriRegisterID, successfulIds, tx)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		reviewFrequency, err := kvus.reviewFrequencySvc.GetByID(existingKRIRegister.ReviewFrequencyID)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		if nextValueUpdate != nil && nextValueUpdate.ID != 0 {
			existingKRIRegister.CurrentValue = nextValueUpdate.Value
			existingKRIRegister.LastReviewNote = nextValueUpdate.Remarks
			existingKRIRegister.LastReviewedDate = &nextValueUpdate.ReviewedDate
			existingKRIRegister.NextReviewDate = utility.CalculateNextReviewDate(reviewFrequency.Type)
			// getting risks related to jri
			existingRisks, err := kvus.repo.GetRisksByKriID(kriRegisterID)
			if err != nil {
				log.Println("error :", err.Error())
			}

			for _, existingRisk := range existingRisks {
				existingKRIRegister, err := kvus.kriRegisterSvc.GetByID(kriRegisterID)
				if err != nil {
					return nil, err
				}

				// Prepare notification data before deletion
				notifications = append(notifications, struct {
					OwnerID         *uint
					TitleEn         string
					TitleAr         string
					Description     string
					ReferenceNumber string
				}{
					OwnerID:         &existingRisk.ResponsibleUserID,
					TitleEn:         existingKRIRegister.Title,
					Description:     *existingKRIRegister.Description,
					ReferenceNumber: *existingKRIRegister.ReferenceNumber,
				})

				// Send DB and email notifications for each successfully deleted risk register
				for _, notification := range notifications {
					urlParam := ""
					entity := "kri-value-update"
					action := "delete"

					// Prepare DB notification
					dbNotificationDTO := &commonDTOs.DBNotificationDTO{
						UserID:               *notification.OwnerID,
						ModuleIdentifier:     "risk",
						NotificationType:     string(commonDTOs.DeletionType),
						HeadingEn:            "KRi Value Update Deleted",
						HeadingAr:            "تم حذف تحديث قيمة KR",
						TitleEn:              notification.TitleEn,
						Entity:               &entity,
						Action:               &action,
						URLParam:             urlParam,
						IsSystemNotification: false,
						CreatedByID:          &currentUserID,
					}

					// Fetch owner details (responsible user)
					ownerUser, err := kvus.userSvc.GetByID(*notification.OwnerID)
					if err != nil {
						utility.LogError("Failed to fetch owner details for notification", map[string]interface{}{
							"error": err.Error(),
						})
						continue
					}

					// Prepare email notification
					ownerEmailTo := []string{ownerUser.Email}
					subject := fmt.Sprintf("KRI:%s,KRI value is updated.", notification.TitleEn)
					appURL := utility.GetAppURL()
					templateRelativePath := "risk/templates/email/kri_value_update_delete_mail.html"

					emailData := map[string]interface{}{
						"FirstNameEn":         ownerUser.FirstNameEn,
						"LastNameEn":          ownerUser.LastNameEn,
						"KRIReferenceNumber":  existingKRIRegister.ReferenceNumber,
						"KRITitle":            existingKRIRegister.Title,
						"KRIDescription":      existingKRIRegister.Description,
						"RiskReferenceNumber": existingRisk.ReferenceNumber,
						"RiskTitle":           existingRisk.Title,
						"RiskDescription":     existingRisk.Description,
						"CurrentValue":        &nextValueUpdate.Value,
						"AssignedOn":          time.Now().Format("2006-01-02"),
						"RedirectURL":         fmt.Sprintf("%srisk/kri-register/%d/overview", appURL, kriRegisterID),
					}

					emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
						EmailTo:              ownerEmailTo,
						Subject:              subject,
						TemplateRelativePath: templateRelativePath,
						EmailData:            emailData,
					}

					// Send DB and email notifications
					err = kvus.notificationSvc.SendNotification(
						dbNotificationDTO, emailNotificationDTO, "kri_register")
					if err != nil {
						utility.LogError("Failed to send deletion notifications", map[string]interface{}{
							"error": err.Error(),
						})
					}
				}
			}
		} else {
			// Set fields to nil if no next value update exists
			existingKRIRegister.CurrentValue = nil
			existingKRIRegister.LastReviewNote = nil
			// existingKRIRegister.LastReviewedDate = nil
			existingKRIRegister.NextReviewDate = utility.CalculateNextReviewDate(reviewFrequency.Type)

			// getting risks related to jri
			existingRisks, err := kvus.repo.GetRisksByKriID(kriRegisterID)
			if err != nil {
				log.Println("error :", err.Error())
			}

			for _, existingRisk := range existingRisks {
				existingKRIRegister, err := kvus.kriRegisterSvc.GetByID(kriRegisterID)
				if err != nil {
					return nil, err
				}

				// Prepare notification data before deletion
				notifications = append(notifications, struct {
					OwnerID         *uint
					TitleEn         string
					TitleAr         string
					Description     string
					ReferenceNumber string
				}{
					OwnerID:         &existingRisk.ResponsibleUserID,
					TitleEn:         existingKRIRegister.Title,
					Description:     *existingKRIRegister.Description,
					ReferenceNumber: *existingKRIRegister.ReferenceNumber,
				})

				// Send DB and email notifications for each successfully deleted risk register
				for _, notification := range notifications {
					urlParam := ""
					entity := "kri-value-update"
					action := "delete"

					// Prepare DB notification
					dbNotificationDTO := &commonDTOs.DBNotificationDTO{
						UserID:               *notification.OwnerID,
						ModuleIdentifier:     "risk",
						NotificationType:     string(commonDTOs.DeletionType),
						HeadingEn:            "KRi Value Update Deleted",
						HeadingAr:            "تم حذف تحديث قيمة KR",
						TitleEn:              notification.TitleEn,
						Entity:               &entity,
						Action:               &action,
						URLParam:             urlParam,
						IsSystemNotification: false,
						CreatedByID:          &currentUserID,
					}

					// Fetch owner details (responsible user)
					ownerUser, err := kvus.userSvc.GetByID(*notification.OwnerID)
					if err != nil {
						utility.LogError("Failed to fetch owner details for notification", map[string]interface{}{
							"error": err.Error(),
						})
						continue
					}

					// Prepare email notification
					ownerEmailTo := []string{ownerUser.Email}
					subject := fmt.Sprintf("KRI:%s,KRI value is updated.", notification.TitleEn)
					appURL := utility.GetAppURL()
					templateRelativePath := "risk/templates/email/kri_value_update_delete_mail.html"
					currentValue := "NIL"
					emailData := map[string]interface{}{
						"FirstNameEn":         ownerUser.FirstNameEn,
						"LastNameEn":          ownerUser.LastNameEn,
						"KRIReferenceNumber":  existingKRIRegister.ReferenceNumber,
						"KRITitle":            existingKRIRegister.Title,
						"KRIDescription":      existingKRIRegister.Description,
						"RiskReferenceNumber": existingRisk.ReferenceNumber,
						"RiskTitle":           existingRisk.Title,
						"RiskDescription":     existingRisk.Description,
						"CurrentValue":        currentValue,
						"AssignedOn":          time.Now().Format("2006-01-02"),
						"RedirectURL":         fmt.Sprintf("%srisk/kri-register/%d/overview", appURL, kriRegisterID),
					}

					emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
						EmailTo:              ownerEmailTo,
						Subject:              subject,
						TemplateRelativePath: templateRelativePath,
						EmailData:            emailData,
					}

					// Send DB and email notifications
					err = kvus.notificationSvc.SendNotification(
						dbNotificationDTO, emailNotificationDTO, "kri_register")
					if err != nil {
						utility.LogError("Failed to send deletion notifications", map[string]interface{}{
							"error": err.Error(),
						})
					}
				}
			}
		}

		if err := kvus.kriRegisterSvc.UpdateValue(existingKRIRegister, tx); err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	tx.Commit()

	if len(successfulIds) > 0 {
		go func() {
			err := kvus.activityLogSvc.LogBulkActivities(
				models.DeleteAction, // Make sure this action constant is defined appropriately
				models.KRIValueUpdate{}.TableName(),
				successfulIds,
				deletedKriValueUpdate,
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

func (kvus *KRIValueUpdateService) parseKRIValueUpdateDates(dto *riskDTOs.KRIValueUpdateFormDTO) (*riskDTOs.KRIValueUpdateParsedDates, error) {
	parsed := &riskDTOs.KRIValueUpdateParsedDates{}

	var err error
	if dto.ReviewedDate != "" {
		parsed.ReviewedDate, err = utility.ParseDate(dto.ReviewedDate)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid review date"), "invalid_date", "review date")
		}
	}

	return parsed, nil
}
