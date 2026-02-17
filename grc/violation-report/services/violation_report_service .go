package services

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"jethur/config"
	commonDTOs "jethur/grc/common/dtos"
	commonServices "jethur/grc/common/services"
	logServices "jethur/grc/logs/services"
	masterRepositories "jethur/grc/masters/repositories"
	masterServices "jethur/grc/masters/services"
	organizationServices "jethur/grc/organization/services"
	settingServices "jethur/grc/settings/services"
	"jethur/grc/violation-report/dtos"
	violationReportDTOs "jethur/grc/violation-report/dtos"
	violationReportRepositories "jethur/grc/violation-report/repositories"

	"jethur/grc/models"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type ViolationReportServiceInterface interface {
	GetAll(params *violationReportDTOs.ViolationReportListQueryParamsDTO, currentUserID uint) (*violationReportDTOs.ViolationReportListResponseDTO, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*violationReportDTOs.ViolationReportDetailDTO, error)
	GetByID(id uint) (*models.ViolationReport, error)
	Create(currentUserID uint, violationReportFormDTO *violationReportDTOs.ViolationReportFormDTO) (uint, error)
	Update(id uint, currentUserID uint, violationReport *violationReportDTOs.ViolationReportFormDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	GetDocumentPath(id uint, token string) (string, error)
	ExportViolationReports(data []*violationReportDTOs.ViolationReportListDTO, currentUserID uint) (*os.File, error)
	CreateViolationReportComment(violationReportID uint, currentUserID uint, dto *dtos.ViolationReportCommentFormDTO) (uint, error)
	GetAllViolationReportComments(violationReportID uint) ([]*dtos.ViolationReportCommentListDTO, error)
}

// Service and implementations
type ViolationReportService struct {
	repo                          violationReportRepositories.ViolationReportRepositoryInterface
	activityLogSvc                logServices.ActivityLogServiceInterface
	violationReportChannelSvc     ViolationReportChannelServiceInterface
	violationReportReasonSvc      ViolationReportReasonServiceInterface
	violationReportCategorySvc    ViolationReportCategoryServiceInterface
	violationReportSubCategorySvc ViolationReportSubCategoryServiceInterface
	violationReportStatusSvc      ViolationReportStatusServiceInterface
	severityLevelSvc              masterServices.SeverityLevelServiceInterface
	userSvc                       organizationServices.UserServiceInterface
	userSettingSvc                organizationServices.UserSettingServiceInterface
	notificationSvc               commonServices.NotificationServiceInterface
	systemPreferenceSvc           settingServices.SystemPreferenceServiceInterface
}

func NewViolationReportService(
	repo violationReportRepositories.ViolationReportRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	violationReportChannelSvc ViolationReportChannelServiceInterface,
	violationReportReasonSvc ViolationReportReasonServiceInterface,
	violationReportCategorySvc ViolationReportCategoryServiceInterface,
	violationReportSubCategorySvc ViolationReportSubCategoryServiceInterface,
	violationReportStatusSvc ViolationReportStatusServiceInterface,
	severityLevelSvc masterRepositories.SeverityLevelRepositoryInterface,
	userSvc organizationServices.UserServiceInterface,
	userSettingSvc organizationServices.UserSettingServiceInterface,
	notificationSvc commonServices.NotificationServiceInterface,
	systemPreferenceSvc settingServices.SystemPreferenceServiceInterface,
) ViolationReportServiceInterface {
	return &ViolationReportService{
		repo:                          repo,
		activityLogSvc:                activityLogSvc,
		violationReportChannelSvc:     violationReportChannelSvc,
		violationReportReasonSvc:      violationReportReasonSvc,
		violationReportCategorySvc:    violationReportCategorySvc,
		violationReportSubCategorySvc: violationReportSubCategorySvc,
		violationReportStatusSvc:      violationReportStatusSvc,
		severityLevelSvc:              severityLevelSvc,
		userSvc:                       userSvc,
		userSettingSvc:                userSettingSvc,
		notificationSvc:               notificationSvc,
		systemPreferenceSvc:           systemPreferenceSvc,
	}
}

func (irs *ViolationReportService) GetAll(params *violationReportDTOs.ViolationReportListQueryParamsDTO, currentUserID uint) (*violationReportDTOs.ViolationReportListResponseDTO, error) {

	user, err := irs.userSvc.GetFullDetailsByID(currentUserID)
	if err != nil {
		return nil, err
	}

	var departmentID *uint
	if user.OrganizationUnit2 != nil {
		id := uint(user.OrganizationUnit2.ID)
		departmentID = &id
	}

	if user.Roles != nil {
		for _, role := range *user.Roles {
			if *role.Type == "super-admin" {
				departmentID = nil
				break
			}
		}
	}
	return irs.repo.GetAll(params, departmentID, currentUserID)
}

func (irs *ViolationReportService) GetByID(id uint) (*models.ViolationReport, error) {
	violationReport, err := irs.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return violationReport, nil
}

func (irs *ViolationReportService) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*violationReportDTOs.ViolationReportDetailDTO, error) {
	// Initialize your DTO
	var violationReportDTO violationReportDTOs.ViolationReportDetailDTO

	violationReport, err := irs.repo.GetFullDetailsByID(id, tx...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	// Map the result to your DTO
	err = mapstructure.Decode(violationReport, &violationReportDTO)
	if err != nil {
		return nil, err
	}

	violationReportDTO.ReportedAt = violationReport.ReportedAt
	violationReportDTO.CreatedAt = violationReport.CreatedAt
	violationReportDTO.UpdatedAt = violationReport.UpdatedAt

	return &violationReportDTO, nil
}

func (irs *ViolationReportService) Create(currentUserID uint, violationReportFormDTO *violationReportDTOs.ViolationReportFormDTO) (uint, error) {
	tx := irs.repo.BeginTransaction()
	if err := irs.validateViolationReportFormDTO(violationReportFormDTO); err != nil {
		return 0, err
	}

	dates, err := irs.parseViolationReportDates(violationReportFormDTO, currentUserID)
	if err != nil {
		return 0, err
	}

	referenceNumber, err := irs.repo.GenerateReferenceNumber()
	if err != nil {
		return 0, err
	}

	// update violation report status as new
	violationReportStatus, err := irs.violationReportStatusSvc.GetByType("new")
	if err != nil {
		return 0, err
	}
	violationReport := &models.ViolationReport{
		ReferenceNumber:              referenceNumber,
		ViolationReportChannelID:     violationReportFormDTO.ViolationReportChannelID,
		ViolationReportReasonID:      violationReportFormDTO.ViolationReportReasonID,
		ViolationReportCategoryID:    violationReportFormDTO.ViolationReportCategoryID,
		ViolationReportSubCategoryID: violationReportFormDTO.ViolationReportSubCategoryID,
		SeverityLevelID:              violationReportFormDTO.SeverityLevelID,
		ViolationReportStatusID:      violationReportStatus.ID,
		ReportedEntityID:             violationReportFormDTO.ReportedEntityID,
		IDOfReporter:                 violationReportFormDTO.IDOfReporter,
		NameOfReporter:               violationReportFormDTO.NameOfReporter,
		PhoneOfReporter:              violationReportFormDTO.PhoneOfReporter,
		EmailOfReporter:              violationReportFormDTO.EmailOfReporter,
		TitleEn:                      violationReportFormDTO.TitleEn,
		TitleAr:                      violationReportFormDTO.TitleAr,
		DescriptionEn:                violationReportFormDTO.DescriptionEn,
		DescriptionAr:                violationReportFormDTO.DescriptionAr,
		EvidenceDescription:          violationReportFormDTO.EvidenceDescription,
		ReportedAt:                   dates.ReportedAt,
		CreatedBYID:                  currentUserID,
	}

	id, err := irs.repo.Create(violationReport, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	var documents []*models.ViolationReportDocument

	// Iterate over each document in the array
	for _, doc := range violationReportFormDTO.Documents {

		// Encrypt the filename
		dataToEncrypt := []byte(doc.FileName)
		encryptedFilename, err := utility.EncryptData(dataToEncrypt)
		if err != nil {
			tx.Rollback()
			return 0, err
		}

		// Construct the URL for the document
		url := utility.GenerateDocumentURL("violation-report", "violation-report-register", strconv.FormatUint(uint64(id), 10), "documents")

		// Create a new ViolationReportDocument instance
		document := &models.ViolationReportDocument{
			ViolationReportID: id,
			URL:               url,
			Name:              doc.FileName,
			Size:              doc.Size,
			Token:             encryptedFilename,
			OriginalName:      doc.OriginalName,
			Type:              doc.Type,
			CreatedByID:       currentUserID,
		}

		// Append the document to the slice
		documents = append(documents, document)
	}
	if len(documents) > 0 {
		// Save the documents
		if err := irs.repo.CreateDocuments(documents, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
		for _, doc := range documents {
			storagePath := filepath.Join("violation-report", "violation-report-register", strconv.FormatUint(uint64(id), 10), "documents")
			if err := utility.StoreFile(doc.Name, storagePath); err != nil {
				tx.Rollback()
				return 0, err
			}
		}
	}

	tx.Commit()

	go func() {
		existingViolationReport, err := irs.GetFullDetailsByID(id)
		if err != nil {
			utility.LogError("Failed to fetch violation report register", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		err = irs.activityLogSvc.LogActivity(
			models.CreateAction,
			models.ViolationReport{}.TableName(),
			id,
			existingViolationReport,
			currentUserID,
			"violation_report", "violation_report_register", "violation_report", nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.ViolationReport{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}

		// Send notifications to users in the category department
		if violationReportFormDTO.ViolationReportCategoryID != nil && *violationReportFormDTO.ViolationReportCategoryID != 0 {
			// Get the violation report with category and department details
			violationReportWithCategory, err := irs.GetFullDetailsByID(id)
			if err != nil {
				utility.LogError("Failed to fetch violation report with category", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}

			// Check if category and department exist
			if violationReportWithCategory.ViolationReportCategory != nil &&
				violationReportWithCategory.ViolationReportCategory.Department != nil {

				departmentID := violationReportWithCategory.ViolationReportCategory.Department.ID
				departmentUsers, err := irs.userSvc.GetUsersByOrganizationUnit(departmentID)
				if err != nil {
					utility.LogError("Failed to fetch department users", map[string]interface{}{
						"error":        err.Error(),
						"departmentID": departmentID,
					})
					return
				}

				if len(departmentUsers) > 0 {
					// Prepare user IDs for notification
					var userIDs []uint
					for _, user := range departmentUsers {
						userIDs = append(userIDs, user.ID)
					}

					// Prepare notification
					urlParam := fmt.Sprintf(`{"violation_report_id":%d}`, id)
					queryParam := ""
					entity := "violation-report"
					action := "new-violation-report-create"

					dbNotificationDTO := &commonDTOs.DBNotificationDTO{
						UserIDs:              userIDs,
						ModuleIdentifier:     "violation_report",
						NotificationType:     string(commonDTOs.AssignmentType),
						HeadingEn:            "New Violation Report Assigned to Your Department",
						HeadingAr:            "تم تعيين تقرير مخالفة جديد لقسمك",
						TitleEn:              violationReportFormDTO.TitleEn,
						TitleAr:              violationReportFormDTO.TitleAr,
						Entity:               &entity,
						Action:               &action,
						URLParam:             urlParam,
						QueryParam:           queryParam,
						IsSystemNotification: false,
						CreatedByID:          &currentUserID,
					}

					// Get creator details for email
					createdBy, err := irs.userSvc.GetByID(currentUserID)
					if err != nil {
						utility.LogError("Failed to fetch creator user", map[string]interface{}{
							"error": err.Error(),
						})
						return
					}

					// Prepare email list for all department users
					var emailTo []string
					for _, user := range departmentUsers {
						emailTo = append(emailTo, user.Email)
					}

					subject := fmt.Sprintf("New Violation Report: %s", violationReportFormDTO.TitleEn)
					appURL := utility.GetAppURL()
					emailData := map[string]interface{}{
						"ViolationReportReferenceNumber": referenceNumber,
						"ViolationReportTitle":           violationReportFormDTO.TitleEn,
						"ViolationReportDescription":     violationReportFormDTO.DescriptionEn,
						"ReportedAt":                     dates.ReportedAt.Format("2006-01-02"),
						"CreatedByFirstNameEn":           createdBy.FirstNameEn,
						"CreatedByLastNameEn":            createdBy.LastNameEn,
						"CreatedOn":                      time.Now().Format("2006-01-02"),
						"RedirectURL":                    fmt.Sprintf("%s/violation-report/register/%d/overview", appURL, id),
					}

					emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
						EmailTo:              emailTo,
						Subject:              subject,
						TemplateRelativePath: "violation-report/templates/email/violation_report_create_mail.html",
						EmailData:            emailData,
					}

					// Send notification and mail
					err = irs.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "violation_report")
					if err != nil {
						utility.LogError("Failed to send notifications", map[string]interface{}{
							"error": err.Error(),
						})
					}
				}
			}
		}

	}()

	return id, nil
}

func (irs *ViolationReportService) Update(id uint, currentUserID uint, violationReportFormDTO *violationReportDTOs.ViolationReportFormDTO) error {

	tx := irs.repo.BeginTransaction()

	existingViolationReport, err := irs.repo.GetByID(id)
	if err != nil {
		return utility.NewResourceNotFoundError(err, id)
	}

	oldViolationReport, err := irs.GetFullDetailsByID(id, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := irs.validateViolationReportFormDTO(violationReportFormDTO); err != nil {
		tx.Rollback()
		return err
	}

	// Note: ReportedAt is not updated during update operations - it's set only during creation
	existingViolationReport.ViolationReportChannelID = violationReportFormDTO.ViolationReportChannelID
	existingViolationReport.ViolationReportReasonID = violationReportFormDTO.ViolationReportReasonID
	existingViolationReport.ViolationReportCategoryID = violationReportFormDTO.ViolationReportCategoryID
	existingViolationReport.ViolationReportSubCategoryID = violationReportFormDTO.ViolationReportSubCategoryID
	existingViolationReport.SeverityLevelID = violationReportFormDTO.SeverityLevelID
	existingViolationReport.TitleEn = violationReportFormDTO.TitleEn
	existingViolationReport.TitleAr = violationReportFormDTO.TitleAr
	existingViolationReport.DescriptionEn = violationReportFormDTO.DescriptionEn
	existingViolationReport.DescriptionAr = violationReportFormDTO.DescriptionAr
	existingViolationReport.EvidenceDescription = violationReportFormDTO.EvidenceDescription
	existingViolationReport.ReportedEntityID = violationReportFormDTO.ReportedEntityID
	existingViolationReport.IDOfReporter = violationReportFormDTO.IDOfReporter
	existingViolationReport.NameOfReporter = violationReportFormDTO.NameOfReporter
	existingViolationReport.PhoneOfReporter = violationReportFormDTO.PhoneOfReporter
	existingViolationReport.EmailOfReporter = violationReportFormDTO.EmailOfReporter
	existingViolationReport.UpdatedBYID = &currentUserID

	if err := irs.repo.Update(existingViolationReport, tx); err != nil {
		tx.Rollback()
		return err
	}

	// ---------- DOCUMENTS ----------
	var documents []*models.ViolationReportDocument
	var deleteDocumentIDs []*uint
	var documentsToDelete []*models.ViolationReportDocument

	for _, doc := range violationReportFormDTO.Documents {

		// NEW FILE
		if doc.IsNewFile != nil && *doc.IsNewFile {

			encryptedFilename, err := utility.EncryptData([]byte(doc.FileName))
			if err != nil {
				tx.Rollback()
				return err
			}

			url := utility.GenerateDocumentURL(
				"violation-report",
				"violation-report-register",
				strconv.FormatUint(uint64(id), 10),
				"documents",
			)

			document := &models.ViolationReportDocument{
				ViolationReportID: id,
				URL:               url,
				Name:              doc.FileName,
				Size:              doc.Size,
				Token:             encryptedFilename,
				OriginalName:      doc.OriginalName,
				Type:              doc.Type,
				CreatedByID:       currentUserID,
			}

			documents = append(documents, document)

			// DELETE FILE
		} else if doc.IsDeletedFile != nil && *doc.IsDeletedFile && doc.ID != nil {

			existingDocument, err := irs.repo.GetDocumentByID(*doc.ID)
			if err != nil {
				tx.Rollback()
				return err
			}

			deleteDocumentIDs = append(deleteDocumentIDs, &existingDocument.ID)
			documentsToDelete = append(documentsToDelete, existingDocument)
		}
	}

	if len(documents) > 0 {
		if err := irs.repo.CreateDocuments(documents, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	if len(deleteDocumentIDs) > 0 {
		if err := irs.repo.DeleteDocuments(deleteDocumentIDs, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	// ---------- FILE SYSTEM ----------
	for _, doc := range documents {
		storagePath := filepath.Join("violation-report", "violation-report-register", strconv.FormatUint(uint64(id), 10), "documents")
		if err := utility.StoreFile(doc.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, doc := range documentsToDelete {
		storagePath := filepath.Join("violation-report", "violation-report-register", strconv.FormatUint(uint64(id), 10), "documents")
		if err := utility.RemoveFile(doc.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()

	go func() {

		newViolationReport, err := irs.GetFullDetailsByID(id)
		if err != nil {
			return
		}

		oldData := map[string]interface{}{
			"title_en":                      oldViolationReport.TitleEn,
			"title_ar":                      oldViolationReport.TitleAr,
			"violation_report_channel":      oldViolationReport.ViolationReportChannel.TitleEn,
			"violation_report_reason":       oldViolationReport.ViolationReportReason.TitleEn,
			"violation_report_category":     oldViolationReport.ViolationReportCategory.TitleEn,
			"violation_report_sub_category": oldViolationReport.ViolationReportSubCategory.TitleEn,
			"description_en":                oldViolationReport.DescriptionEn,
			"description_ar":                oldViolationReport.DescriptionAr,
			"severity_level":                oldViolationReport.SeverityLevel.TitleEn,
			"reported_entity":               oldViolationReport.ReportedEntity.TitleEn,
			"reported_at":                   oldViolationReport.ReportedAt,
		}

		newData := map[string]interface{}{
			"title_en":                      newViolationReport.TitleEn,
			"title_ar":                      newViolationReport.TitleAr,
			"violation_report_channel":      newViolationReport.ViolationReportChannel.TitleEn,
			"violation_report_reason":       newViolationReport.ViolationReportReason.TitleEn,
			"violation_report_category":     newViolationReport.ViolationReportCategory.TitleEn,
			"violation_report_sub_category": newViolationReport.ViolationReportSubCategory.TitleEn,
			"description_en":                newViolationReport.DescriptionEn,
			"description_ar":                newViolationReport.DescriptionAr,
			"severity_level":                newViolationReport.SeverityLevel.TitleEn,
			"reported_entity":               newViolationReport.ReportedEntity.TitleEn,
			"reported_at":                   newViolationReport.ReportedAt,
		}

		changedData := utility.CompareFieldChanges(oldData, newData)

		if len(changedData) > 0 {
			_ = irs.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.ViolationReport{}.TableName(),
				id,
				changedData,
				currentUserID,
				"violation_report", "violation_report_register", "violation_report", nil,
			)
		}
	}()

	return nil
}

func (irs *ViolationReportService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	tx := irs.repo.BeginTransaction()
	var failedIds []uint
	var successfulIds []uint
	var deletedViolationReports []map[string]interface{}

	documentsToDelete := make(map[uint][]*violationReportDTOs.ViolationReportDocumentDTO)

	for _, id := range ids {
		violationReport, err := irs.GetFullDetailsByID(id, tx)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		deletedViolationReports = append(deletedViolationReports, map[string]interface{}{
			"id":    id,
			"title": violationReport.TitleEn,
		})

		// Delete each document if documents exist
		if len(violationReport.Documents) > 0 {
			var deleteDocumentIds []*uint
			for _, doc := range violationReport.Documents {
				if doc != nil {
					documentsToDelete[id] = append(documentsToDelete[id], doc)
					deleteDocumentIds = append(deleteDocumentIds, doc.ID)
				}
			}
			if err := irs.repo.DeleteDocuments(deleteDocumentIds, tx); err != nil {
				tx.Rollback()
				return nil, err
			}

		}

		// Delete the violation report
		if err := irs.repo.Delete(id, tx); err != nil {
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
		for violationReportID, documents := range documentsToDelete {
			for _, doc := range documents {
				fmt.Println("filename=", doc.Name)
				storagePath := filepath.Join("violation-report", "violation-report-register", strconv.FormatUint(uint64(violationReportID), 10), "documents")
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
			err := irs.activityLogSvc.LogBulkActivities(
				models.DeleteAction, // Make sure this action constant is defined appropriately
				models.ViolationReport{}.TableName(),
				successfulIds,
				deletedViolationReports,
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

func (irs *ViolationReportService) validateViolationReportFormDTO(violationReportFormDTO *violationReportDTOs.ViolationReportFormDTO) error {
	// Validate violation report channel
	if _, err := irs.violationReportChannelSvc.GetByID(violationReportFormDTO.ViolationReportChannelID); err != nil {
		return utility.NewValidationError(fmt.Errorf("violation report channel does not exist"), "not_exist", "violation_report_channel")
	}

	// Validate violation report category
	if _, err := irs.violationReportCategorySvc.GetByID(*violationReportFormDTO.ViolationReportCategoryID); err != nil {
		return utility.NewValidationError(fmt.Errorf("violation report category does not exist"), "not_exist", "violation_report_category")
	}

	// Validate violation report sub category if provided
	if violationReportFormDTO.ViolationReportSubCategoryID != nil {
		if _, err := irs.violationReportSubCategorySvc.GetByID(*violationReportFormDTO.ViolationReportSubCategoryID); err != nil {
			return utility.NewValidationError(fmt.Errorf("violation report sub category does not exist"), "not_exist", "violation_report_sub_category")
		}
	}

	// Validate violation report reason
	if _, err := irs.violationReportReasonSvc.GetByID(violationReportFormDTO.ViolationReportReasonID); err != nil {
		return utility.NewValidationError(fmt.Errorf("violation report reason does not exist"), "not_exist", "violation_report_reason")
	}

	// Validate severity level if provided
	if violationReportFormDTO.SeverityLevelID != nil {
		if _, err := irs.severityLevelSvc.GetByID(*violationReportFormDTO.SeverityLevelID); err != nil {
			return utility.NewValidationError(fmt.Errorf("severity level does not exist"), "not_exist", "severity_level")
		}
	}

	// Validate documents if provided
	if len(violationReportFormDTO.Documents) > 0 {
		if err := irs.validateViolationReportDocumentFormDTO(violationReportFormDTO.Documents, violationReportFormDTO.EvidenceDescription); err != nil {
			return err
		}
	}

	return nil
}

func (irs *ViolationReportService) parseViolationReportDates(dto *violationReportDTOs.ViolationReportFormDTO, currentUserID uint) (*violationReportDTOs.ViolationReportParsedDates, error) {
	parsed := &violationReportDTOs.ViolationReportParsedDates{}

	userPreference, err := irs.userSvc.GetCurrentUserDetails(currentUserID)
	if err != nil {
		return nil, err
	}

	timezoneName := userPreference.UserSetting.Timezone.Name

	// If ReportedAt is provided, parse it; otherwise use current time
	if dto.ReportedAt != "" {
		reportedAtPtr, err := utility.ConvertToUTC(dto.ReportedAt, timezoneName)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid reported date: %w", err), "invalid_date", "reported_at")
		}
		parsed.ReportedAt = *reportedAtPtr
	} else {
		// Automatically set ReportedAt to current time if not provided
		parsed.ReportedAt = time.Now()
	}

	return parsed, nil
}

func (irs *ViolationReportService) validateViolationReportDocumentFormDTO(documents []*violationReportDTOs.ViolationReportDocumentFormDTO, evidenceDescription *string) error {
	// Documents are optional, but if provided, evidence_description is required
	if len(documents) > 0 {
		// If documents are present, evidence_description is required
		if evidenceDescription == nil || *evidenceDescription == "" {
			return utility.NewValidationError(fmt.Errorf("evidence description is required when documents are provided"), "required", "evidence_description")
		}
	}

	return nil
}

func (irs *ViolationReportService) GetDocumentPath(id uint, token string) (string, error) {
	existingViolationReport, err := irs.GetFullDetailsByID(id)
	if err != nil {
		return "", err
	}

	if len(existingViolationReport.Documents) == 0 {
		return "", fmt.Errorf("no documents found for the violation report")
	}

	decryptedFilenameBytes, err := utility.DecryptData(token)
	if err != nil {
		return "", fmt.Errorf("resource not found")
	}

	decryptedFilename := string(decryptedFilenameBytes)
	if decryptedFilename == "" {
		return "", fmt.Errorf("resource not found")
	}

	storagePath := filepath.Join("violation-report", "violation-report-register", strconv.FormatUint(uint64(id), 10), "documents")

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

func (irs *ViolationReportService) ExportViolationReports(data []*violationReportDTOs.ViolationReportListDTO, currentUserID uint) (*os.File, error) {
	// Fetch user setting for language preference
	userSetting, err := irs.userSettingSvc.GetUserSettingsByUserID(currentUserID)
	if err != nil {
		return nil, err
	}
	var headers []string
	// Define headers based on the user's language setting
	if userSetting.Language == "ar" {
		headers = []string{
			"SL No.",
			"Reference Number",
			"Title",
			"Reported At",
			"Channel",
			"Category",
			"Sub Category",
			"Reason",
			"Severity Level",
			"Reported Entity",
			"Reporter Name",
			"Reporter Phone",
			"Description",
			"Current Status",
			"Created By",
			"Created At",
		}
	} else {
		headers = []string{
			"SL No.",
			"Reference Number",
			"Title",
			"Reported At",
			"Channel",
			"Category",
			"Sub Category",
			"Reason",
			"Severity Level",
			"Reported Entity",
			"Reporter Name",
			"Reporter Phone",
			"Description",
			"Current Status",
			"Created By",
			"Created At",
		}
	}
	// Convert the violation report data to [][]interface{} based on language setting
	rows := make([][]interface{}, len(data))
	for i, report := range data {
		if userSetting.Language == "ar" {
			// Populate row data in Arabic
			rows[i] = []interface{}{
				i + 1,
				report.ReferenceNumber,
				report.TitleAr,
				report.ReportedAt.Format("2006/01/02 03:04 PM"),
				report.ViolationReportChannelTitleAr,
				report.ViolationReportCategoryTitleAr,
				utility.GetStringValue(report.ViolationReportSubCategoryTitleAr),
				report.ViolationReportReasonTitleAr,
				utility.GetStringValue(report.SeverityLevelTitleAr),
				report.ReportedEntityTitleAr,
				// utility.GetStringValue(report.NameOfReporter),
				// utility.GetStringValue(report.PhoneOfReporter),
				utility.GetStringValue(report.DescriptionAr),
				report.ViolationReportStatusTitleAr,
				report.CreatedUserFirstNameAr + " " + report.CreatedUserLastNameAr,
				report.CreatedAt.Format("2006/01/02 03:04 PM"),
			}
		} else {
			// Populate row data in English
			rows[i] = []interface{}{
				i + 1,
				report.ReferenceNumber,
				report.TitleEn,
				report.ReportedAt.Format("2006/01/02 03:04 PM"),
				report.ViolationReportChannelTitleEn,
				report.ViolationReportCategoryTitleEn,
				utility.GetStringValue(report.ViolationReportSubCategoryTitleEn),
				report.ViolationReportReasonTitleEn,
				utility.GetStringValue(report.SeverityLevelTitleEn),
				report.ReportedEntityTitleEn,
				// utility.GetStringValue(report.NameOfReporter),
				// utility.GetStringValue(report.PhoneOfReporter),
				utility.GetStringValue(report.DescriptionEn),
				report.ViolationReportStatusTitleEn,
				report.CreatedUserFirstNameEn + " " + report.CreatedUserLastNameEn,
				report.CreatedAt.Format("2006/01/02 03:04 PM"),
			}
		}
	}

	// Generate the Excel file
	return utility.GenerateExcelFile(headers, rows)
}

func (s *ViolationReportService) CreateViolationReportComment(violationReportID uint, currentUserID uint, dto *dtos.ViolationReportCommentFormDTO) (uint, error) {

	existingReport, err := s.repo.GetByID(violationReportID)
	if err != nil {
		return 0, utility.NewResourceNotFoundError(err, violationReportID)
	}

	//  Start transaction
	tx := s.repo.BeginTransaction()
	now := time.Now()

	//  Create comment model
	comment := &models.ViolationReportComment{
		ViolationReportID:       violationReportID,
		ViolationReportStatusID: dto.ViolationReportStatusID,
		Comment:                 dto.Comment,
		CreatedByID:             currentUserID,
		CreatedAt:               now,
	}

	commentID, err := s.repo.CreateComment(comment, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	// Update main report status if provided
	if dto.ViolationReportStatusID != nil {
		existingReport.ViolationReportStatusID = *dto.ViolationReportStatusID
	}

	existingReport.UpdatedBYID = &currentUserID
	existingReport.UpdatedAt = &now

	if err := s.repo.Update(existingReport, tx); err != nil {
		tx.Rollback()
		return 0, err
	}

	tx.Commit()

	//  Async activity log
	go func() {
		fullDetails, err := s.GetFullDetailsByID(violationReportID)
		if err != nil {
			return
		}

		err = s.activityLogSvc.LogActivity(
			models.CreateAction,
			models.ViolationReportComment{}.TableName(),
			commentID,
			fullDetails,
			currentUserID,
			"violation_report",
			"violation",
			"",
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":   models.CreateAction,
				"recordID": commentID,
				"userID":   currentUserID,
				"error":    err.Error(),
			}
			utility.LogError("Failed to log violation report comment activity", fields)
		}
	}()

	return commentID, nil
}

func (s *ViolationReportService) GetAllViolationReportComments(violationReportID uint) ([]*dtos.ViolationReportCommentListDTO, error) {
	return s.repo.GetAllViolationReportComments(violationReportID)
}
