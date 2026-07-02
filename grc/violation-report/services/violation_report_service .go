package services

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"strings"
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
	CreateExternal(violationReportFormDTO *violationReportDTOs.ViolationReportFormDTO, files []*multipart.FileHeader) (uint, error)
	Update(id uint, currentUserID uint, violationReport *violationReportDTOs.ViolationReportFormDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	GetDocumentPath(id uint, token string) (string, error)
	ExportViolationReports(data []*violationReportDTOs.ViolationReportListDTO, currentUserID uint) (*os.File, error)
	CreateViolationReportComment(violationReportID uint, currentUserID uint, dto *dtos.ViolationReportCommentFormDTO) (uint, error)
	UpdateViolationReportComment(violationReportID uint, commentID uint, currentUserID uint, dto *dtos.ViolationReportCommentUpdateFormDTO) error
	GetAllViolationReportComments(violationReportID uint) ([]*dtos.ViolationReportCommentListDTO, error)
	GetCommentDetailsByID(commentID uint) (*dtos.ViolationReportCommentDetailDTO, error)
	GetViolationReportCommentDocumentPath(violationReportID uint, commentID uint, token string) (string, error)
	DeleteComment(violationReportID uint, ids []uint, currentUserID uint) ([]uint, error)
	ExportViolationReportToPDF(id uint, currentUserID uint, investigationSvc ViolationReportInvestigationServiceInterface) (string, error)
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

	violationReportDTO.ReportedAtString = violationReport.ReportedAt.Format("2006-01-02")
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
		ReportedEntity:               violationReportFormDTO.ReportedEntity,
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
					urlParam := fmt.Sprintf(`{"violation_report_id":%d}`, id)
					queryParam := ""
					entity := "violation-report"
					action := "new-violation-report-create"

					// Collect all user IDs for bulk DB notification
					var userIDs []uint
					for _, user := range departmentUsers {
						userIDs = append(userIDs, user.ID)
					}

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

					if err := irs.notificationSvc.SendNotification(dbNotificationDTO, nil, "violation_report"); err != nil {
						fmt.Printf("failed to send violation report DB notifications: %v\n", err)
						utility.LogError("Failed to send DB notifications", map[string]interface{}{
							"error": err.Error(),
						})
					}

					// Get creator details for email
					createdBy, err := irs.userSvc.GetByID(currentUserID)
					if err != nil {
						utility.LogError("Failed to fetch creator user", map[string]interface{}{
							"error": err.Error(),
						})
						return
					}

					subject := fmt.Sprintf("New Violation Report: %s - %s", referenceNumber, violationReportFormDTO.TitleEn)
					appURL := utility.GetAppURL()
					sentEmails := map[string]bool{}

					// Send individual email per user so each recipient gets a personalised email
					for _, user := range departmentUsers {
						if user.Email == "" {
							continue
						}
						sentEmails[strings.ToLower(strings.TrimSpace(user.Email))] = true
						emailData := map[string]interface{}{
							"RecipientFirstNameEn":           user.FirstNameEn,
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
							EmailTo:              []string{user.Email},
							Subject:              subject,
							TemplateRelativePath: "violation-report/templates/email/violation_report_create_mail.html",
							EmailData:            emailData,
						}
						if err := irs.notificationSvc.SendNotification(nil, emailNotificationDTO, "violation_report"); err != nil {
							fmt.Printf("failed to send violation report email to user %d (%s): %v\n", user.ID, user.Email, err)
							utility.LogError("Failed to send email notification for user", map[string]interface{}{
								"error":  err.Error(),
								"userID": user.ID,
								"email":  user.Email,
							})
						}
					}

					// Send email to reporter as well, if provided and not already included
					if violationReportFormDTO.EmailOfReporter != nil {
						reporterEmail := strings.TrimSpace(*violationReportFormDTO.EmailOfReporter)
						if reporterEmail != "" && !sentEmails[strings.ToLower(reporterEmail)] {
							reporterFirstName := "User"
							if violationReportFormDTO.NameOfReporter != nil {
								nameParts := strings.Fields(strings.TrimSpace(*violationReportFormDTO.NameOfReporter))
								if len(nameParts) > 0 {
									reporterFirstName = nameParts[0]
								}
							}

							emailData := map[string]interface{}{
								"RecipientFirstNameEn":           reporterFirstName,
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
								EmailTo:              []string{reporterEmail},
								Subject:              subject,
								TemplateRelativePath: "violation-report/templates/email/violation_report_create_mail.html",
								EmailData:            emailData,
							}
							if err := irs.notificationSvc.SendNotification(nil, emailNotificationDTO, "violation_report"); err != nil {
								fmt.Printf("failed to send violation report email to reporter (%s): %v\n", reporterEmail, err)
								utility.LogError("Failed to send email notification for reporter", map[string]interface{}{
									"error": err.Error(),
									"email": reporterEmail,
								})
							}
						}
					}
				}
			}
		}

	}()

	return id, nil
}

func (irs *ViolationReportService) CreateExternal(violationReportFormDTO *violationReportDTOs.ViolationReportFormDTO, files []*multipart.FileHeader) (uint, error) {
	if err := irs.validateViolationReportFormDTO(violationReportFormDTO); err != nil {
		return 0, err
	}

	tx := irs.repo.BeginTransaction()

	superAdmin, err := irs.userSvc.GetSuperAdmin()
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	dates, err := irs.parseViolationReportDates(violationReportFormDTO, superAdmin.ID)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	referenceNumber, err := irs.repo.GenerateReferenceNumber()
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	violationReportStatus, err := irs.violationReportStatusSvc.GetByType("new")
	if err != nil {
		tx.Rollback()
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
		ReportedEntity:               violationReportFormDTO.ReportedEntity,
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
		CreatedBYID:                  superAdmin.ID,
		IsExternal:                   true,
	}

	id, err := irs.repo.Create(violationReport, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	var documents []*models.ViolationReportDocument

	for _, file := range files {
		timestamp := time.Now().UnixNano()
		sanitizedName := fmt.Sprintf("%d_%s", timestamp, utility.SanitizeFileName(file.Filename))

		dataToEncrypt := []byte(sanitizedName)
		encryptedFilename, err := utility.EncryptData(dataToEncrypt)
		if err != nil {
			tx.Rollback()
			return 0, err
		}

		url := utility.GenerateDocumentURL("violation-report", "violation-report-register", strconv.FormatUint(uint64(id), 10), "documents")

		document := &models.ViolationReportDocument{
			ViolationReportID: id,
			URL:               url,
			Name:              sanitizedName,
			Size:              int(file.Size),
			Token:             encryptedFilename,
			OriginalName:      file.Filename,
			Type:              file.Header.Get("Content-Type"),
			CreatedByID:       superAdmin.ID,
		}

		documents = append(documents, document)
	}

	if len(documents) > 0 {
		if err := irs.repo.CreateDocuments(documents, tx); err != nil {
			tx.Rollback()
			return 0, err
		}

		storagePath := config.GetEnv("STORAGE_PATH", "")
		for i, file := range files {
			doc := documents[i]
			dstDir := filepath.Join(storagePath, "violation-report", "violation-report-register", strconv.FormatUint(uint64(id), 10), "documents")
			if err := os.MkdirAll(dstDir, 0755); err != nil {
				tx.Rollback()
				return 0, err
			}
			dst := filepath.Join(dstDir, doc.Name)

			src, err := file.Open()
			if err != nil {
				tx.Rollback()
				return 0, err
			}

			out, err := os.Create(dst)
			if err != nil {
				src.Close()
				tx.Rollback()
				return 0, err
			}

			if _, err := io.Copy(out, src); err != nil {
				out.Close()
				src.Close()
				tx.Rollback()
				return 0, err
			}
			out.Close()
			src.Close()
		}
	}

	tx.Commit()

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
	existingViolationReport.ReportedEntity = violationReportFormDTO.ReportedEntity
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
			"reported_entity":               oldViolationReport.ReportedEntity,
			"reported_at":                   oldViolationReport.ReportedAtString,
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
			"reported_entity":               newViolationReport.ReportedEntity,
			"reported_at":                   newViolationReport.ReportedAtString,
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
	category, err := irs.violationReportCategorySvc.GetByID(*violationReportFormDTO.ViolationReportCategoryID)
	if err != nil {
		return utility.NewValidationError(fmt.Errorf("violation report category does not exist"), "not_exist", "violation_report_category")
	}

	// Validate violation report sub category if provided
	if violationReportFormDTO.ViolationReportSubCategoryID != nil {
		subCategory, err := irs.violationReportSubCategorySvc.GetByID(*violationReportFormDTO.ViolationReportSubCategoryID)
		if err != nil {
			return utility.NewValidationError(fmt.Errorf("violation report sub category does not exist"), "not_exist", "violation_report_sub_category")
		}

		if subCategory.ViolationReportCategoryID != category.ID {
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

	// Validate reporter email format if provided
	if violationReportFormDTO.EmailOfReporter != nil {
		email := strings.TrimSpace(*violationReportFormDTO.EmailOfReporter)
		if email != "" && !utility.IsValidEmail(email) {
			return utility.NewValidationError(fmt.Errorf("email of reporter is invalid"), "invalid", "email_of_reporter")
		}
	}

	return nil
}

func (irs *ViolationReportService) parseViolationReportDates(dto *violationReportDTOs.ViolationReportFormDTO, currentUserID uint) (*violationReportDTOs.ViolationReportParsedDates, error) {
	parsed := &violationReportDTOs.ViolationReportParsedDates{}

	// If ReportedAt is provided, parse date only; otherwise use current date
	if strings.TrimSpace(dto.ReportedAt) != "" {
		reportedAtPtr, err := utility.ParseDate(strings.TrimSpace(dto.ReportedAt))
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid reported date: %w", err), "invalid_date", "reported_at")
		}
		if reportedAtPtr != nil {
			parsed.ReportedAt = *reportedAtPtr
		}
	} else {
		currentDate, _ := utility.ParseDate(time.Now().Format("2006-01-02"))
		if currentDate != nil {
			parsed.ReportedAt = *currentDate
		}
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
			"Reporter ID",
			"Reporter Name",
			"Reporter Email",
			"Reporter Phone",
			"Evidence Description",
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
			"Reporter ID",
			"Reporter Name",
			"Reporter Email",
			"Reporter Phone",
			"Evidence Description",
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
				report.ReportedAt,
				report.ViolationReportChannelTitleAr,
				*report.ViolationReportCategoryTitleAr,
				utility.GetStringValue(report.ViolationReportSubCategoryTitleAr),
				report.ViolationReportReasonTitleAr,
				utility.GetStringValue(report.SeverityLevelTitleAr),
				utility.GetStringValue(report.ReportedEntity),
				*report.IDOfReporter,
				utility.GetStringValue(report.NameOfReporter),
				utility.GetStringValue(report.EmailOfReporter),
				utility.GetStringValue(report.PhoneOfReporter),
				utility.GetStringValue(report.EvidenceDescription),
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
				report.ReportedAt,
				report.ViolationReportChannelTitleEn,
				*report.ViolationReportCategoryTitleEn,
				utility.GetStringValue(report.ViolationReportSubCategoryTitleEn),
				report.ViolationReportReasonTitleEn,
				utility.GetStringValue(report.SeverityLevelTitleEn),
				utility.GetStringValue(report.ReportedEntity),
				*report.IDOfReporter,
				utility.GetStringValue(report.NameOfReporter),
				utility.GetStringValue(report.EmailOfReporter),
				utility.GetStringValue(report.PhoneOfReporter),
				utility.GetStringValue(report.EvidenceDescription),
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
		CreatedBYID:             currentUserID,
		CreatedAt:               now,
	}

	commentID, err := s.repo.CreateComment(comment, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	var documents []*models.ViolationReportCommentDocument
	for _, doc := range dto.Documents {
		if doc.FileName == "" || doc.Size <= 0 {
			continue
		}

		encryptedFilename, err := utility.EncryptData([]byte(doc.FileName))
		if err != nil {
			tx.Rollback()
			return 0, err
		}

		url := utility.GenerateDocumentURL(
			"violation-report",
			"violation-report-register",
			strconv.FormatUint(uint64(violationReportID), 10),
			"comments",
			strconv.FormatUint(uint64(commentID), 10),
			"documents",
		)

		document := &models.ViolationReportCommentDocument{
			ViolationReportCommentID: commentID,
			URL:                      url,
			Name:                     doc.FileName,
			Size:                     doc.Size,
			Token:                    encryptedFilename,
			OriginalName:             doc.OriginalName,
			Type:                     doc.Type,
			CreatedByID:              currentUserID,
			CreatedAt:                now,
		}

		documents = append(documents, document)
	}

	if len(documents) > 0 {
		if err := s.repo.CreateCommentDocuments(documents, tx); err != nil {
			tx.Rollback()
			return 0, err
		}

		for _, doc := range documents {
			storagePath := filepath.Join(
				"violation-report",
				"violation-report-register",
				strconv.FormatUint(uint64(violationReportID), 10),
				"comments",
				strconv.FormatUint(uint64(commentID), 10),
				"documents",
			)
			if err := utility.StoreFile(doc.Name, storagePath); err != nil {
				tx.Rollback()
				return 0, err
			}
		}
	}

	// Update main report status if provided
	if dto.ViolationReportStatusID != nil {
		existingReport.ViolationReportStatusID = *dto.ViolationReportStatusID
		existingReport.UpdatedBYID = &currentUserID
		existingReport.UpdatedAt = &now
		if err := s.repo.Update(existingReport, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
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

	// Async status update notifications/email to department users
	if dto.ViolationReportStatusID != nil {
		go func() {
			fullDetails, err := s.GetFullDetailsByID(violationReportID)
			if err != nil {
				return
			}

			if fullDetails.ViolationReportCategory != nil && fullDetails.ViolationReportCategory.Department != nil {
				departmentID := fullDetails.ViolationReportCategory.Department.ID
				departmentUsers, err := s.userSvc.GetUsersByOrganizationUnit(departmentID)
				if err != nil {
					utility.LogError("Failed to fetch department users for status update notification", map[string]interface{}{
						"error":        err.Error(),
						"departmentID": departmentID,
					})
					return
				}

				if len(departmentUsers) > 0 {
					urlParam := fmt.Sprintf(`{"violation_report_id":%d}`, violationReportID)
					queryParam := ""
					entity := "violation-report"
					action := "violation-report-status-update"

					var userIDs []uint
					for _, user := range departmentUsers {
						userIDs = append(userIDs, user.ID)
					}

					dbNotificationDTO := &commonDTOs.DBNotificationDTO{
						UserIDs:              userIDs,
						ModuleIdentifier:     "violation_report",
						NotificationType:     string(commonDTOs.AssignmentType),
						HeadingEn:            "Violation Report Status Updated",
						HeadingAr:            "تم تحديث حالة تقرير المخالفة",
						TitleEn:              fullDetails.TitleEn,
						TitleAr:              fullDetails.TitleAr,
						Entity:               &entity,
						Action:               &action,
						URLParam:             urlParam,
						QueryParam:           queryParam,
						IsSystemNotification: false,
						CreatedByID:          &currentUserID,
					}

					if err := s.notificationSvc.SendNotification(dbNotificationDTO, nil, "violation_report"); err != nil {
						fmt.Printf("failed to send violation report status update DB notifications: %v\n", err)
						utility.LogError("Failed to send status update DB notifications", map[string]interface{}{
							"error": err.Error(),
						})
					}

					createdBy, err := s.userSvc.GetByID(currentUserID)
					if err != nil {
						utility.LogError("Failed to fetch updater user for status update email", map[string]interface{}{
							"error": err.Error(),
						})
						return
					}

					subject := fmt.Sprintf("Violation Report Status Updated: %s - %s", fullDetails.ReferenceNumber, fullDetails.TitleEn)
					appURL := utility.GetAppURL()
					sentEmails := map[string]bool{}

					for _, user := range departmentUsers {
						if user.Email == "" {
							continue
						}
						sentEmails[strings.ToLower(strings.TrimSpace(user.Email))] = true

						emailData := map[string]interface{}{
							"RecipientFirstNameEn":           user.FirstNameEn,
							"ViolationReportReferenceNumber": fullDetails.ReferenceNumber,
							"ViolationReportTitle":           fullDetails.TitleEn,
							"ViolationReportStatusTitle":     fullDetails.ViolationReportStatus.TitleEn,
							"ViolationReportDescription":     fullDetails.DescriptionEn,
							"ReportedAt":                     fullDetails.ReportedAtString,
							"CreatedByFirstNameEn":           createdBy.FirstNameEn,
							"CreatedByLastNameEn":            createdBy.LastNameEn,
							"CreatedOn":                      time.Now().Format("2006-01-02"),
							"RedirectURL":                    fmt.Sprintf("%s/violation-report/register/%d/overview", appURL, violationReportID),
						}

						emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
							EmailTo:              []string{user.Email},
							Subject:              subject,
							TemplateRelativePath: "violation-report/templates/email/violation_report_status_update_mail.html",
							EmailData:            emailData,
						}
						if err := s.notificationSvc.SendNotification(nil, emailNotificationDTO, "violation_report"); err != nil {
							fmt.Printf("failed to send violation report status update email to user %d (%s): %v\n", user.ID, user.Email, err)
							utility.LogError("Failed to send status update email notification for user", map[string]interface{}{
								"error":  err.Error(),
								"userID": user.ID,
								"email":  user.Email,
							})
						}
					}

					// Send status update email to reporter as well, if provided and not already included
					if fullDetails.EmailOfReporter != nil {
						reporterEmail := strings.TrimSpace(*fullDetails.EmailOfReporter)
						if reporterEmail != "" && !sentEmails[strings.ToLower(reporterEmail)] {
							reporterFirstName := "User"
							if fullDetails.NameOfReporter != nil {
								nameParts := strings.Fields(strings.TrimSpace(*fullDetails.NameOfReporter))
								if len(nameParts) > 0 {
									reporterFirstName = nameParts[0]
								}
							}

							emailData := map[string]interface{}{
								"RecipientFirstNameEn":           reporterFirstName,
								"ViolationReportReferenceNumber": fullDetails.ReferenceNumber,
								"ViolationReportTitle":           fullDetails.TitleEn,
								"ViolationReportStatusTitle":     fullDetails.ViolationReportStatus.TitleEn,
								"ViolationReportDescription":     fullDetails.DescriptionEn,
								"ReportedAt":                     fullDetails.ReportedAtString,
								"CreatedByFirstNameEn":           createdBy.FirstNameEn,
								"CreatedByLastNameEn":            createdBy.LastNameEn,
								"CreatedOn":                      time.Now().Format("2006-01-02"),
								"RedirectURL":                    fmt.Sprintf("%s/violation-report/register/%d/overview", appURL, violationReportID),
							}

							emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
								EmailTo:              []string{reporterEmail},
								Subject:              subject,
								TemplateRelativePath: "violation-report/templates/email/violation_report_status_update_mail.html",
								EmailData:            emailData,
							}
							if err := s.notificationSvc.SendNotification(nil, emailNotificationDTO, "violation_report"); err != nil {
								fmt.Printf("failed to send violation report status update email to reporter (%s): %v\n", reporterEmail, err)
								utility.LogError("Failed to send status update email notification for reporter", map[string]interface{}{
									"error": err.Error(),
									"email": reporterEmail,
								})
							}
						}
					}
				}
			}
		}()
	}

	return commentID, nil
}

func (s *ViolationReportService) GetAllViolationReportComments(violationReportID uint) ([]*dtos.ViolationReportCommentListDTO, error) {
	return s.repo.GetAllViolationReportComments(violationReportID)
}

func (s *ViolationReportService) GetCommentDetailsByID(commentID uint) (*dtos.ViolationReportCommentDetailDTO, error) {
	comment, err := s.repo.GetFullCommentDetailsByID(commentID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, commentID)
		}
		return nil, err
	}

	var detailDTO dtos.ViolationReportCommentDetailDTO
	if err := mapstructure.Decode(comment, &detailDTO); err != nil {
		return nil, err
	}
	detailDTO.CreatedAt = comment.CreatedAt
	detailDTO.UpdatedAt = comment.UpdatedAt

	return &detailDTO, nil
}

func (s *ViolationReportService) UpdateViolationReportComment(violationReportID uint, commentID uint, currentUserID uint, dto *dtos.ViolationReportCommentUpdateFormDTO) error {
	if _, err := s.repo.GetByID(violationReportID); err != nil {
		return utility.NewResourceNotFoundError(err, violationReportID)
	}

	existingComment, err := s.repo.GetCommentByID(commentID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utility.NewResourceNotFoundError(err, commentID)
		}
		return err
	}

	if existingComment.ViolationReportID != violationReportID {
		return utility.NewResourceNotFoundError(fmt.Errorf("violation report comment not found"), commentID)
	}

	if existingComment.ViolationReportStatusID != nil {
		return utility.NewValidationError(
			fmt.Errorf("comments with status cannot be updated"),
			"invalid",
			"comment_id",
		)
	}

	tx := s.repo.BeginTransaction()
	now := time.Now()
	oldComment := existingComment.Comment

	existingComment.Comment = dto.Comment
	existingComment.IsEdited = true
	existingComment.UpdatedBYID = &currentUserID
	existingComment.UpdatedAt = &now

	if err := s.repo.UpdateComment(existingComment, tx); err != nil {
		tx.Rollback()
		return err
	}

	var newDocuments []*models.ViolationReportCommentDocument
	var deleteDocumentIDs []*uint
	var documentsToDelete []*models.ViolationReportCommentDocument

	for _, doc := range dto.Documents {
		if doc.IsNewFile != nil && *doc.IsNewFile {
			if doc.FileName == "" || doc.Size <= 0 {
				continue
			}

			encryptedFilename, err := utility.EncryptData([]byte(doc.FileName))
			if err != nil {
				tx.Rollback()
				return err
			}

			url := utility.GenerateDocumentURL(
				"violation-report",
				"violation-report-register",
				strconv.FormatUint(uint64(violationReportID), 10),
				"comments",
				strconv.FormatUint(uint64(commentID), 10),
				"documents",
			)

			document := &models.ViolationReportCommentDocument{
				ViolationReportCommentID: commentID,
				URL:                      url,
				Name:                     doc.FileName,
				Size:                     doc.Size,
				Token:                    encryptedFilename,
				OriginalName:             doc.OriginalName,
				Type:                     doc.Type,
				CreatedByID:              currentUserID,
				CreatedAt:                now,
			}
			newDocuments = append(newDocuments, document)
		} else if doc.IsDeletedFile != nil && *doc.IsDeletedFile && doc.ID != nil {
			existingDocument, err := s.repo.GetCommentDocumentByID(*doc.ID)
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					tx.Rollback()
					return utility.NewResourceNotFoundError(err, *doc.ID)
				}
				tx.Rollback()
				return err
			}

			if existingDocument.ViolationReportCommentID != commentID {
				tx.Rollback()
				return utility.NewResourceNotFoundError(fmt.Errorf("violation report comment document not found"), *doc.ID)
			}

			documentsToDelete = append(documentsToDelete, existingDocument)
			deleteDocumentIDs = append(deleteDocumentIDs, &existingDocument.ID)
		}
	}

	if len(newDocuments) > 0 {
		if err := s.repo.CreateCommentDocuments(newDocuments, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	if len(deleteDocumentIDs) > 0 {
		if err := s.repo.DeleteCommentDocuments(deleteDocumentIDs, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, doc := range newDocuments {
		storagePath := filepath.Join(
			"violation-report",
			"violation-report-register",
			strconv.FormatUint(uint64(violationReportID), 10),
			"comments",
			strconv.FormatUint(uint64(commentID), 10),
			"documents",
		)
		if err := utility.StoreFile(doc.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, doc := range documentsToDelete {
		storagePath := filepath.Join(
			"violation-report",
			"violation-report-register",
			strconv.FormatUint(uint64(violationReportID), 10),
			"comments",
			strconv.FormatUint(uint64(commentID), 10),
			"documents",
		)
		if err := utility.RemoveFile(doc.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()

	go func() {
		changedData := utility.CompareFieldChanges(
			map[string]interface{}{"comment": oldComment},
			map[string]interface{}{"comment": existingComment.Comment},
		)

		if len(changedData) == 0 {
			return
		}

		if err := s.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.ViolationReportComment{}.TableName(),
			commentID,
			changedData,
			currentUserID,
			"violation_report",
			"violation",
			"",
			nil,
		); err != nil {
			utility.LogError("Failed to log violation report comment update activity", map[string]interface{}{
				"recordID": commentID,
				"userID":   currentUserID,
				"error":    err.Error(),
			})
		}
	}()

	return nil
}

func (s *ViolationReportService) GetViolationReportCommentDocumentPath(violationReportID uint, commentID uint, token string) (string, error) {
	comment, err := s.repo.GetCommentByID(commentID)
	if err != nil {
		return "", err
	}

	if comment.ViolationReportID != violationReportID {
		return "", fmt.Errorf("resource not found")
	}

	if len(comment.Documents) == 0 {
		return "", fmt.Errorf("no documents found for this violation report comment")
	}

	decryptedFilenameBytes, err := utility.DecryptData(token)
	if err != nil {
		return "", fmt.Errorf("resource not found")
	}

	decryptedFilename := string(decryptedFilenameBytes)
	if decryptedFilename == "" {
		return "", fmt.Errorf("resource not found")
	}

	storagePath := filepath.Join(
		"violation-report",
		"violation-report-register",
		strconv.FormatUint(uint64(violationReportID), 10),
		"comments",
		strconv.FormatUint(uint64(commentID), 10),
		"documents",
	)

	baseStoragePath := config.GetEnv("STORAGE_PATH", "")
	fullPath := filepath.Join(baseStoragePath, storagePath, decryptedFilename)

	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return "", fmt.Errorf("file not found")
	} else if err != nil {
		return "", fmt.Errorf("error accessing file: %w", err)
	}

	return fullPath, nil
}

func (s *ViolationReportService) DeleteComment(violationReportID uint, ids []uint, currentUserID uint) ([]uint, error) {

	tx := s.repo.BeginTransaction()

	// Validate violation report exists
	if _, err := s.repo.GetByID(violationReportID); err != nil {
		return nil, err
	}

	var failedIDs []uint
	var successfulIDs []uint
	var deletedComments []map[string]interface{}

	for _, id := range ids {

		comment, err := s.repo.GetFullCommentDetailsByID(id)
		if err != nil {
			failedIDs = append(failedIDs, id)
			continue
		}

		// Ensure comment belongs to this violation report
		if comment.ViolationReportID != violationReportID {
			failedIDs = append(failedIDs, id)
			continue
		}

		// Prevent delete if the comment carries status update information.
		if comment.ViolationReportStatusID != nil {
			failedIDs = append(failedIDs, id)
			continue
		}

		for _, doc := range comment.Documents {
			storagePath := filepath.Join(
				"violation-report",
				"violation-report-register",
				strconv.FormatUint(uint64(violationReportID), 10),
				"comments",
				strconv.FormatUint(uint64(id), 10),
				"documents",
			)

			if err := utility.RemoveFile(doc.Name, storagePath); err != nil {
				tx.Rollback()
				return nil, err
			}
		}

		if err := s.repo.DeleteComment(id, tx); err != nil {
			failedIDs = append(failedIDs, id)
			continue
		}

		successfulIDs = append(successfulIDs, id)

		deletedComments = append(deletedComments, map[string]interface{}{
			"id":                  id,
			"violation_report_id": violationReportID,
		})
	}

	if len(failedIDs) > 0 {
		tx.Rollback()
		return failedIDs, nil
	}

	tx.Commit()

	if len(successfulIDs) > 0 {
		go func() {
			err := s.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.ViolationReportComment{}.TableName(),
				successfulIDs,
				deletedComments,
				currentUserID,
				nil,
			)
			if err != nil {
				fmt.Printf("Failed to log bulk delete activity: %v\n", err)
			}
		}()
	}

	return failedIDs, nil
}
func (irs *ViolationReportService) ExportViolationReportToPDF(id uint, currentUserID uint, investigationSvc ViolationReportInvestigationServiceInterface) (string, error) {
	violationReportDTO, err := irs.GetFullDetailsByID(id)
	if err != nil {
		return "", err
	}

	// Fetch investigation data if exists
	var investigationDTO *violationReportDTOs.ViolationReportInvestigationExportDTO
	investigationDetail, err := investigationSvc.GetFullDetailsByID(id)
	if err == nil && investigationDetail != nil {
		var departments []*violationReportDTOs.ViolationReportInvestigationDepartmentDTO
		if investigationDetail.Departments != nil {
			departments = investigationDetail.Departments
		}
		investigationDTO = &violationReportDTOs.ViolationReportInvestigationExportDTO{
			StepsTaken:                         investigationDetail.StepsTaken,
			InvestigationDateString:            investigationDetail.InvestigationDateString,
			SuggestedRecommendation:            investigationDetail.SuggestedRecommendation,
			Departments:                        departments,
			FindingSummary:                     investigationDetail.FindingSummary,
			DisciplinaryProcedure:              investigationDetail.DisciplinaryProcedure,
			IsDisciplinaryRegulatoryCompliance: investigationDetail.IsDisciplinaryRegulatoryCompliance,
			Documents:                          investigationDetail.Documents,
			NoteAndObservations:                investigationDetail.NoteAndObservations,
			CreatedByUserFirstNameEn:           investigationDetail.CreatedByUser.FirstNameEN,
			CreatedByUserLastNameEn:            investigationDetail.CreatedByUser.LastNameEN,
			CreatedAtString:                    investigationDetail.CreatedAt.Format("Jan 02, 2006 03:04 PM"),
		}
	}

	// Prepare overview data
	createdAt := violationReportDTO.CreatedAt.Format("2006-01-02")
	generatedAt := time.Now().Format("Jan 02, 2006")
	reportedAtString := violationReportDTO.ReportedAtString

	var descriptionEn *string
	if violationReportDTO.DescriptionEn != "" {
		descriptionEn = &violationReportDTO.DescriptionEn
	}

	var descriptionAr *string
	if violationReportDTO.DescriptionAr != "" {
		descriptionAr = &violationReportDTO.DescriptionAr
	}

	viewModel := violationReportDTOs.ViolationReportExportPDFDTO{
		Title:                      violationReportDTO.TitleEn,
		ReferenceNumber:            violationReportDTO.ReferenceNumber,
		CreatedAt:                  createdAt,
		GeneratedAt:                generatedAt,
		CreatedUserFirstNameEn:     violationReportDTO.CreatedByUser.FirstNameEN,
		CreatedUserLastNameEn:      violationReportDTO.CreatedByUser.LastNameEN,
		ReportedAtString:           reportedAtString,
		IDOfReporter:               violationReportDTO.IDOfReporter,
		NameOfReporter:             violationReportDTO.NameOfReporter,
		PhoneOfReporter:            violationReportDTO.PhoneOfReporter,
		EmailOfReporter:            violationReportDTO.EmailOfReporter,
		ViolationReportChannel:     violationReportDTO.ViolationReportChannel,
		ViolationReportReason:      violationReportDTO.ViolationReportReason,
		ViolationReportCategory:    violationReportDTO.ViolationReportCategory,
		ViolationReportSubCategory: violationReportDTO.ViolationReportSubCategory,
		ViolationReportStatus:      violationReportDTO.ViolationReportStatus,
		SeverityLevel:              violationReportDTO.SeverityLevel,
		ReportedByUserFirstNameEn:  violationReportDTO.CreatedByUser.FirstNameEN,
		ReportedByUserLastNameEn:   violationReportDTO.CreatedByUser.LastNameEN,
		DescriptionEn:              descriptionEn,
		DescriptionAr:              descriptionAr,
		Documents:                  violationReportDTO.Documents,
		Investigation:              investigationDTO,
	}

	// Body
	bodyTemplates := []string{"violation-report/templates/pdf/violation_report_register_body.html"}
	htmlContent, err := utility.GeneratePDFContent(viewModel, bodyTemplates, "", "")
	if err != nil {
		return "", fmt.Errorf("failed to generate HTML content: %w", err)
	}

	// Cover
	coverTemplates := []string{"violation-report/templates/pdf/violation_report_register_cover.html"}
	coverContent, err := utility.GeneratePDFContent(viewModel, coverTemplates, "", "")
	if err != nil {
		return "", fmt.Errorf("failed to generate cover content: %w", err)
	}

	// Header
	headerTemplates := []string{"violation-report/templates/pdf/violation_report_register_header.html"}
	headerContent, err := utility.GeneratePDFContent(viewModel, headerTemplates, "", "")
	if err != nil {
		return "", fmt.Errorf("failed to generate header content: %w", err)
	}

	// Footer
	footerTemplates := []string{"violation-report/templates/pdf/violation_report_register_footer.html"}
	footerContent, err := utility.GeneratePDFContent(viewModel, footerTemplates, "", "")
	if err != nil {
		return "", fmt.Errorf("failed to generate footer content: %w", err)
	}

	pdfFileName := fmt.Sprintf("ViolationReport_%d.pdf", id)
	if err := utility.GeneratePDFWithCoverHeaderFooter(coverContent, htmlContent, pdfFileName, headerContent, footerContent); err != nil {
		return "", fmt.Errorf("failed to generate PDF: %w", err)
	}

	return pdfFileName, nil
}
