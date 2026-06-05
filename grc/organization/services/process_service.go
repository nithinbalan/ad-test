package services

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	commonDTOs "jethur/grc/common/dtos"
	commonServices "jethur/grc/common/services"
	logServices "jethur/grc/logs/services"
	masterServices "jethur/grc/masters/services"
	"jethur/grc/models"
	organizationDTOs "jethur/grc/organization/dtos"
	organizationRepositories "jethur/grc/organization/repositories"
	settingsServices "jethur/grc/settings/services"
	"jethur/utility"

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
	ValidateProcessTemplate(file io.Reader) ([]map[string]interface{}, error)
	ValidateImportProcess(input []map[string]interface{}) ([]map[string]interface{}, error)
	ImportProcesses(input []map[string]interface{}, currentUserID uint) ([]uint, error)

	UpdateProcessRiskRating(process *models.Process, tx utility.DatabaseTransaction) error
	UpdateProcessRecoveryObjective(process *models.Process, tx utility.DatabaseTransaction) error
}

type ProcessService struct {
	repo                 organizationRepositories.ProcessRepositoryInterface
	activityLogSvc       logServices.ActivityLogServiceInterface
	processGroupSvc      ProcessGroupServiceInterface
	userSvc              UserServiceInterface
	lifecycleStatusSvc   masterServices.LifecycleStatusServiceInterface
	controlSvc           ControlServiceInterface
	notificationSvc      commonServices.NotificationServiceInterface
	organizationUnitSvc  settingsServices.OrganizationUnitServiceInterface
	userSettingSvc       UserSettingServiceInterface
	organizationUnitsSvc OrganizationUnitsServiceInterface
}

func NewProcessService(
	repo organizationRepositories.ProcessRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	processGroupSvc ProcessGroupServiceInterface,
	userSvc UserServiceInterface,
	lifecycleStatusSvc masterServices.LifecycleStatusServiceInterface,
	controlSvc ControlServiceInterface,
	notificationSvc commonServices.NotificationServiceInterface,
	organizationUnitSvc settingsServices.OrganizationUnitServiceInterface,
	userSettingSvc UserSettingServiceInterface,
	organizationUnitsSvc OrganizationUnitsServiceInterface,
) ProcessServiceInterface {
	return &ProcessService{
		repo:                 repo,
		activityLogSvc:       activityLogSvc,
		processGroupSvc:      processGroupSvc,
		userSvc:              userSvc,
		lifecycleStatusSvc:   lifecycleStatusSvc,
		controlSvc:           controlSvc,
		notificationSvc:      notificationSvc,
		organizationUnitSvc:  organizationUnitSvc,
		userSettingSvc:       userSettingSvc,
		organizationUnitsSvc: organizationUnitsSvc,
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
	// organizationUnitID, err := ps.userSvc.GetUserPrimaryUnitID(uint(processFormDTO.OwnerID))
	// if err != nil {
	// 	return 0, utility.NewValidationError(fmt.Errorf("organization unit not set for owner"), "not_exist", "process owner unit")
	// }

	process := &models.Process{
		ProcessGroupID:    processFormDTO.ProcessGroupID,
		ResponsibleUserID: processFormDTO.OwnerID,
		ResponsibleUnitID: processFormDTO.OwnerUnitID,
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

	tx.Commit()

	go func() {
		existingProcess, err := ps.GetFullDetailsByID(process.ID)
		if err != nil {
			utility.LogError("Failed to fetch responsible user", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		// log
		err = ps.activityLogSvc.LogActivity(
			models.CreateAction,
			models.Process{}.TableName(),
			id,
			existingProcess,
			currentUserID,
			"organization", "process", "create", nil,
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

		// notification
		urlParam := fmt.Sprintf(`{"id":%d}`, id)
		entity := "process"
		action := "responsible-user-assign"

		var emailTo []string
		var userIDs []uint
		var emailData map[string]interface{}
		var subject string

		if processFormDTO.OwnerID != nil {
			// Send notification to responsible user
			userIDs = []uint{*processFormDTO.OwnerID}

			dbNotificationDTO := &commonDTOs.DBNotificationDTO{
				UserIDs:              userIDs,
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

			responsibleUser, err := ps.userSvc.GetByID(*processFormDTO.OwnerID)
			if err != nil {
				utility.LogError("Failed to fetch responsible user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}

			emailTo = []string{responsibleUser.Email}
			subject = fmt.Sprintf("Process/%s: Ownership Assigned", processFormDTO.TitleEn)
			appURL := utility.GetAppURL()
			emailData = map[string]interface{}{
				"FirstNameEn":            responsibleUser.FirstNameEn,
				"LastNameEn":             responsibleUser.LastNameEn,
				"ProcessReferenceNumber": existingProcess.ReferenceNumber,
				"ProcessName":            existingProcess.TitleEn,
				"ProcessGroup":           existingProcess.ProcessGroup.TitleEn,
				"AssignedByFirstNameEn":  existingProcess.CreatedByUser.FirstNameEN,
				"AssignedByLastNameEn":   existingProcess.CreatedByUser.LastNameEN,
				"AssignedOn":             existingProcess.CreatedAt.Format("2006-01-02"),
				"RedirectURL":            fmt.Sprintf("%s/organization/process/%d/overview", appURL, id),
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
				utility.LogError("Failed to send notification", map[string]interface{}{
					"error": err.Error(),
				})
			}
		} else {
			// Send notification to department head since no responsible user is assigned
			headUsers, err := ps.organizationUnitsSvc.GetParentOrganizationUnitAndCurrentUnitHead(processFormDTO.OwnerUnitID)
			if err != nil {
				utility.LogError("Failed to fetch department head users for notification", map[string]interface{}{
					"error":       err.Error(),
					"ownerUnitID": processFormDTO.OwnerUnitID,
				})
				return
			}

			if headUsers.HeadUser == nil {
				utility.LogInfo("No head users found for department", map[string]interface{}{
					"ownerUnitID": processFormDTO.OwnerUnitID,
				})
				return
			}

			// Get department details
			department, err := ps.organizationUnitSvc.GetByID(processFormDTO.OwnerUnitID)
			if err != nil {
				utility.LogError("Failed to fetch department details for notification", map[string]interface{}{
					"error":       err.Error(),
					"ownerUnitID": processFormDTO.OwnerUnitID,
				})
				return
			}

			if headUsers.HeadUser != nil {
				userIDs = append(userIDs, headUsers.HeadUser.ID)
				emailTo = append(emailTo, headUsers.HeadUser.Email)
			}

			dbNotificationDTO := &commonDTOs.DBNotificationDTO{
				UserIDs:              userIDs,
				ModuleIdentifier:     "organization",
				NotificationType:     string(commonDTOs.InformationType),
				ReferenceNumber:      process.ReferenceNumber,
				HeadingEn:            "Process Created - Responsible Unit Assigned",
				HeadingAr:            "تم إنشاء العملية - تم تعيين وحدة مسؤولية",
				TitleEn:              fmt.Sprintf("Process '%s' has been created with a responsible unit", processFormDTO.TitleEn),
				TitleAr:              fmt.Sprintf("تم إنشاء العملية '%s' مع وحدة مسؤولية", processFormDTO.TitleAr),
				Entity:               &entity,
				Action:               &action,
				URLParam:             urlParam,
				IsSystemNotification: false,
				CreatedByID:          &currentUserID,
			}

			appURL := utility.GetAppURL()
			emailData = map[string]interface{}{
				"DepartmentNameEn":       department.TitleEn,
				"DepartmentNameAr":       department.TitleAr,
				"ProcessReferenceNumber": existingProcess.ReferenceNumber,
				"ProcessName":            existingProcess.TitleEn,
				"ProcessGroup":           existingProcess.ProcessGroup.TitleEn,
				"AssignedByFirstNameEn":  existingProcess.CreatedByUser.FirstNameEN,
				"AssignedByLastNameEn":   existingProcess.CreatedByUser.LastNameEN,
				"AssignedOn":             existingProcess.CreatedAt.Format("2006-01-02"),
				"RedirectURL":            fmt.Sprintf("%s/organization/process/%d/overview", appURL, id),
			}

			subject = fmt.Sprintf("Process/%s: Created - Action Required", existingProcess.TitleEn)
			emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
				EmailTo:              emailTo,
				Subject:              subject,
				TemplateRelativePath: "organization/templates/email/process_no_owner_mail.html",
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
				utility.LogError("Failed to send notification", map[string]interface{}{
					"error": err.Error(),
				})
			}
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

	oldOwnerID := process.ResponsibleUserID

	if err := ps.validateProcessFormDTO(processFormDTO); err != nil {
		tx.Rollback()
		return err
	}

	// If owner has changed, update the OwnerUnitID
	// if process.ResponsibleUserID != processFormDTO.OwnerID {
	// 	organizationUnitID, err := ps.userSvc.GetUserPrimaryUnitID(uint(processFormDTO.OwnerID))
	// 	if err != nil {
	// 		tx.Rollback()
	// 		return utility.NewValidationError(fmt.Errorf("organization unit not set for owner"), "not_exist", "process owner unit")
	// 	}
	// 	process.ResponsibleUnitID = organizationUnitID
	// }

	process.ProcessGroupID = processFormDTO.ProcessGroupID
	process.ResponsibleUserID = processFormDTO.OwnerID
	process.ResponsibleUnitID = processFormDTO.OwnerUnitID
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

	tx.Commit()

	go func() {

		changedData := make(map[string]interface{})
		// log changes
		oldOwnerName := "No Owner"
		if oldProcess.ResponsibleUser != nil {
			oldOwnerName = fmt.Sprintf("%s %s", oldProcess.ResponsibleUser.FirstNameEN, oldProcess.ResponsibleUser.LastNameEN)
		}

		newOwnerName := "No Owner"
		if existingProcess.ResponsibleUser != nil {
			newOwnerName = fmt.Sprintf("%s %s", existingProcess.ResponsibleUser.FirstNameEN, existingProcess.ResponsibleUser.LastNameEN)
		}

		oldData := map[string]interface{}{
			"process_group":    oldProcess.ProcessGroup.TitleEn,
			"owner_unit":       oldProcess.ResponsibleUnit.TitleEn,
			"owner":            oldOwnerName,
			"lifecycle_status": oldProcess.LifecycleStatus.TitleEn,
			"title_en":         oldProcess.TitleEn,
			"title_ar":         oldProcess.TitleAr,
			"description":      oldProcess.Description,
		}

		newData := map[string]interface{}{
			"process_group":    existingProcess.ProcessGroup.TitleEn,
			"owner_unit":       existingProcess.ResponsibleUnit.TitleEn,
			"owner":            newOwnerName,
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
				"", "", "", nil,
			)
			if err != nil {
				utility.LogError("Failed to log activity", map[string]interface{}{
					"action":   models.UpdateAction,
					"error":    err.Error(),
					"recordID": id,
				})
			}
		}

		// notification - handle owner changes
		ownerChanged := false
		if (oldOwnerID == nil && processFormDTO.OwnerID != nil) ||
			(oldOwnerID != nil && processFormDTO.OwnerID == nil) ||
			(oldOwnerID != nil && processFormDTO.OwnerID != nil && *oldOwnerID != *processFormDTO.OwnerID) {
			ownerChanged = true
		}

		if ownerChanged {
			urlParam := fmt.Sprintf(`{"id":%d}`, id)
			entity := "process"
			appURL := utility.GetAppURL()

			// Send notification to old owner if they existed
			if oldOwnerID != nil {
				action := "responsible-user-remove"

				oldOwnerNotificationDTO := &commonDTOs.DBNotificationDTO{
					UserIDs:              []uint{*oldOwnerID},
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
					IsSystemNotification: false,
					CreatedByID:          &currentUserID,
				}

				oldOwnerUser, err := ps.userSvc.GetByID(*oldOwnerID)
				if err != nil {
					utility.LogError("Failed to fetch old responsible user", map[string]interface{}{
						"error": err.Error(),
					})
				} else {
					oldOwnerEmailTo := []string{oldOwnerUser.Email}
					oldOwnerSubject := fmt.Sprintf("Process/%s: Ownership Removed", process.TitleEn)
					oldOwnerTemplateRelativePath := "organization/templates/email/process_remove_owner_mail.html"

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

					err = ps.notificationSvc.SendNotification(oldOwnerNotificationDTO, oldOwnerEmailNotificationDTO, "process")
					if err != nil {
						utility.LogError("Failed to send old owner notifications", map[string]interface{}{
							"error": err.Error(),
						})
					}
				}
			}

			// Send notification to new owner or department head
			if processFormDTO.OwnerID != nil {
				// Send to new responsible user
				newOwnerAction := "responsible-user-assign"

				dbNotificationDTO := &commonDTOs.DBNotificationDTO{
					UserIDs:              []uint{*processFormDTO.OwnerID},
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
					IsSystemNotification: false,
					CreatedByID:          &currentUserID,
				}

				responsibleUser, err := ps.userSvc.GetByID(*processFormDTO.OwnerID)
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
					"RedirectURL":            fmt.Sprintf("%s/organization/process/%d/overview", appURL, id),
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
			} else {
				// OwnerID is null, send notification to department head
				headUsers, err := ps.organizationUnitsSvc.GetParentOrganizationUnitAndCurrentUnitHead(processFormDTO.OwnerUnitID)
				if err != nil {
					utility.LogError("Failed to fetch department head users for notification", map[string]interface{}{
						"error":       err.Error(),
						"ownerUnitID": processFormDTO.OwnerUnitID,
					})
					return
				}

				if headUsers.HeadUser == nil {
					utility.LogInfo("No head users found for department", map[string]interface{}{
						"ownerUnitID": processFormDTO.OwnerUnitID,
					})
					return
				}

				// Get department details
				department, err := ps.organizationUnitSvc.GetByID(processFormDTO.OwnerUnitID)
				if err != nil {
					utility.LogError("Failed to fetch department details for notification", map[string]interface{}{
						"error":       err.Error(),
						"ownerUnitID": processFormDTO.OwnerUnitID,
					})
					return
				}

				var userIDs []uint
				var emailTo []string

				if headUsers.HeadUser != nil {
					userIDs = append(userIDs, headUsers.HeadUser.ID)
					emailTo = append(emailTo, headUsers.HeadUser.Email)
				}

				newOwnerAction := "responsible-user-assign"

				dbNotificationDTO := &commonDTOs.DBNotificationDTO{
					UserIDs:              userIDs,
					ModuleIdentifier:     "organization",
					NotificationType:     string(commonDTOs.InformationType),
					ReferenceNumber:      process.ReferenceNumber,
					HeadingEn:            "Process Updated - Responsible Unit Assigned",
					HeadingAr:            "تم تحديث العملية - تم تعيين وحدة مسؤولية",
					TitleEn:              fmt.Sprintf("Process '%s' has been updated with a responsible unit", processFormDTO.TitleEn),
					TitleAr:              fmt.Sprintf("تم تحديث العملية '%s' مع وحدة مسؤولية", processFormDTO.TitleAr),
					Entity:               &entity,
					Action:               &newOwnerAction,
					URLParam:             urlParam,
					IsSystemNotification: false,
					CreatedByID:          &currentUserID,
				}

				emailData := map[string]interface{}{
					"DepartmentNameEn":       department.TitleEn,
					"DepartmentNameAr":       department.TitleAr,
					"ProcessReferenceNumber": existingProcess.ReferenceNumber,
					"ProcessName":            existingProcess.TitleEn,
					"ProcessGroup":           existingProcess.ProcessGroup.TitleEn,
					"AssignedByFirstNameEn":  existingProcess.UpdatedByUser.FirstNameEN,
					"AssignedByLastNameEn":   existingProcess.UpdatedByUser.LastNameEN,
					"AssignedOn":             existingProcess.UpdatedAt.Format("2006-01-02"),
					"RedirectURL":            fmt.Sprintf("%s/organization/process/%d/overview", appURL, id),
				}

				subject := fmt.Sprintf("Process/%s: Updated - Action Required", existingProcess.TitleEn)
				emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
					EmailTo:              emailTo,
					Subject:              subject,
					TemplateRelativePath: "organization/templates/email/process_no_owner_mail.html",
					EmailData:            emailData,
				}

				err = ps.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "process")
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
			OwnerID:         process.ResponsibleUserID,
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
					UserIDs:              []uint{*notification.OwnerID},
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

	// Validate OwnerUnitID (required)
	if _, err := ps.organizationUnitSvc.GetByID(processFormDTO.OwnerUnitID); err != nil {
		return utility.NewValidationError(fmt.Errorf("owner unit does not exist"), "not_exist", "owner unit")
	}

	// Validate OwnerID only if it's provided (optional)
	if processFormDTO.OwnerID != nil {
		if _, err := ps.userSvc.GetByID(*processFormDTO.OwnerID); err != nil {
			return utility.NewValidationError(fmt.Errorf("process owner does not exist"), "not_exist", "process owner")
		}
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
	// Headers for the first sheet (Template)
	templateHeaders := []string{
		"ProcessGroup", "Owner", "TitleEn", "TitleAr", "Description", "Controls",
	}
	templateData := [][]interface{}{} // Empty template data

	// Fetch process groups for the second sheet
	processGroups, err := ps.processGroupSvc.GetAll("", "", "", "active")
	if err != nil {
		return nil, err
	}

	// Headers and data for the second sheet (ProcessGroup)
	processGroupHeaders := []string{"ProcessGroupTitleEn", "ProcessGroupTitleAr"}
	processGroupData := [][]interface{}{}
	for _, group := range processGroups {
		processGroupData = append(processGroupData, []interface{}{group.TitleEn, group.TitleAr})
	}

	// Headers and data for the third sheet (Owner)
	// Fetch user details
	users, err := ps.userSvc.GetAllUsers()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch owners: %v", err)
	}

	// Prepare data for Owner sheet
	ownerHeaders := []string{"OwnerName", "Email", "Designation", "PrimaryDepartment"}
	ownerData := [][]interface{}{}
	for _, user := range users {
		ownerData = append(ownerData, []interface{}{
			fmt.Sprintf("%s %s", user.FirstNameEN, user.LastNameEN),
			user.Email,
			user.Designation,
			user.PrimaryUnit,
		})
	}

	// Fetch controls for the fourth sheet
	controls, err := ps.controlSvc.GetAllControls()
	if err != nil {
		return nil, err
	}

	// Headers and data for the fourth sheet (Controls)
	controlHeaders := []string{"Controls"}
	controlData := [][]interface{}{}
	for _, control := range controls {
		controlData = append(controlData, []interface{}{control.TitleEn})
	}

	// Generate the Excel file with "Template," "ProcessGroup," "Owner," and "Controls" sheets
	file, err := utility.GenerateExcelFileWithSheets(map[string][][]interface{}{
		"Template":     append([][]interface{}{utility.ConvertToInterfaceSlice(templateHeaders)}, templateData...),
		"ProcessGroup": append([][]interface{}{utility.ConvertToInterfaceSlice(processGroupHeaders)}, processGroupData...),
		"Owner":        append([][]interface{}{utility.ConvertToInterfaceSlice(ownerHeaders)}, ownerData...),
		"Controls":     append([][]interface{}{utility.ConvertToInterfaceSlice(controlHeaders)}, controlData...),
	})
	if err != nil {
		return nil, err
	}

	return file, nil
}

func (ps *ProcessService) ValidateProcessTemplate(file io.Reader) ([]map[string]interface{}, error) {
	// Read Excel file
	rows, err := utility.ReadExcelFile(file)
	if err != nil {
		return nil, err
	}

	// Define headers in the input format and map them to snake case
	inputHeaders := []string{"ProcessGroup", "Owner", "TitleEn", "TitleAr", "Description", "Controls"}
	snakeCaseHeaders := []string{"process_group", "owner", "title_en", "title_ar", "description", "controls"}

	if len(rows) == 0 {
		return nil, fmt.Errorf("empty Excel file")
	}

	// Validate headers
	if !utility.ValidateHeaders(inputHeaders, rows[0]) {
		return nil, utility.NewValidationError(fmt.Errorf("invalid headers in the Excel file"), "invalid_header", "Process Template")
	}

	// Convert rows to input format
	var input []map[string]interface{}
	for _, row := range rows[1:] { // Skip header row
		record := make(map[string]interface{})
		for i, header := range snakeCaseHeaders {
			if i < len(row) {
				record[header] = row[i]
			} else {
				record[header] = "" // Pad missing columns
			}
		}
		input = append(input, record)
	}

	validatedData, err := ps.validateProcesses(input)
	if err != nil {
		return nil, err
	}

	return validatedData, nil
}

func (ps *ProcessService) ValidateImportProcess(input []map[string]interface{}) ([]map[string]interface{}, error) {
	validatedData, err := ps.validateProcesses(input)
	if err != nil {
		return nil, err
	}

	return validatedData, nil
}

func (ps *ProcessService) ImportProcesses(input []map[string]interface{}, currentUserID uint) ([]uint, error) {
	validatedData, err := ps.validateProcesses(input)
	if err != nil {
		return nil, err
	}

	for _, record := range validatedData {
		if errors, ok := record["errors"].(map[string]interface{}); ok && len(errors) > 0 {
			return nil, utility.NewValidationError(fmt.Errorf("invalid data in the import file"), "invalid_data", "Invalid Data")
		}
	}

	var createdProcessIDs []uint
	for _, record := range validatedData {
		// Handle optional OwnerID
		var ownerID *uint
		if ownerIDVal, ok := record["owner_id"]; ok && ownerIDVal != nil {
			if oid, ok := ownerIDVal.(uint); ok {
				ownerID = &oid
			}
		}

		processFormDTO := &organizationDTOs.ProcessFormDTO{
			ProcessGroupID: record["process_group_id"].(uint),
			OwnerUnitID:    record["owner_unit_id"].(uint),
			OwnerID:        ownerID,
			TitleEn:        record["title_en"].(string),
			TitleAr:        record["title_ar"].(string),
			Description:    utility.GetOptionalStringValue(record["description"]),
		}

		if controls, ok := record["control_ids"].([]uint); ok {
			var controlIDs []*uint
			for _, id := range controls {
				controlID := id
				controlIDs = append(controlIDs, &controlID)
			}
			processFormDTO.ControlIds = controlIDs
		} else {
			processFormDTO.ControlIds = nil
		}

		processID, err := ps.Create(currentUserID, processFormDTO)
		if err != nil {
			return nil, fmt.Errorf("failed to create process for record: %+v, error: %v", record, err)
		}

		createdProcessIDs = append(createdProcessIDs, processID)
	}

	return createdProcessIDs, nil
}

func (ps *ProcessService) validateProcesses(input []map[string]interface{}) ([]map[string]interface{}, error) {
	headers := []string{"process_group", "owner", "title_en", "title_ar", "description", "controls"}

	// Fetch process group titles and map to IDs
	processGroups, err := ps.processGroupSvc.GetAll("", "", "", "active")
	if err != nil {
		return nil, fmt.Errorf("error fetching process groups: %v", err)
	}
	processGroupTitles := make(map[string]uint)
	for _, group := range processGroups {
		processGroupTitles[strings.TrimSpace(group.TitleEn)] = group.ID
	}

	// Fetch user full names and map to IDs
	users, err := ps.userSvc.GetAllUsers()
	if err != nil {
		return nil, fmt.Errorf("error fetching users: %v", err)
	}
	userFullNames := make(map[string]uint)
	for _, user := range users {
		fullName := strings.TrimSpace(user.FirstNameEN + " " + user.LastNameEN)
		userFullNames[fullName] = uint(user.ID)
	}

	// Fetch control titles and map to IDs
	controls, err := ps.controlSvc.GetAllControls()
	if err != nil {
		return nil, fmt.Errorf("error fetching controls: %v", err)
	}
	controlTitles := make(map[string]uint)
	for _, control := range controls {
		controlTitles[strings.TrimSpace(control.TitleEn)] = control.ID
	}

	// Perform validation
	var validatedData []map[string]interface{}

	for _, record := range input {
		rowErrors := make(map[string]interface{})
		processedRecord := make(map[string]interface{})

		// Copy the original record for the response
		for key, value := range record {
			processedRecord[key] = value
		}

		// Validate fields and map IDs
		for _, header := range headers {
			value, ok := record[header].(string)
			switch header {
			case "process_group":
				if !ok || value == "" {
					rowErrors["process_group"] = "Field is required"
				} else if id, exists := processGroupTitles[strings.TrimSpace(value)]; !exists {
					rowErrors["process_group"] = fmt.Sprintf("Invalid process group title: %s", value)
				} else {
					processedRecord["process_group_id"] = id // Add the ID
				}

			case "owner":
				if !ok || value == "" {
					rowErrors["owner"] = "Field is required"
				} else if id, exists := userFullNames[strings.TrimSpace(value)]; !exists {
					rowErrors["owner"] = fmt.Sprintf("Invalid owner name: %s", value)
				} else {
					processedRecord["owner_id"] = id // Add the ID
				}

			case "controls":
				if ok && value != "" { // Validate only if controls field is provided
					controlList := utility.ParseCommaSeparatedValues(value)
					var controlIDs []uint
					invalidControls := []string{}
					for _, control := range controlList {
						trimmedControl := strings.TrimSpace(control)
						if id, exists := controlTitles[trimmedControl]; exists {
							controlIDs = append(controlIDs, id)
						} else {
							invalidControls = append(invalidControls, fmt.Sprintf("Invalid control: %s", trimmedControl))
						}
					}
					if len(invalidControls) > 0 {
						rowErrors["controls"] = strings.Join(invalidControls, "; ") // Combine errors into a single string
					} else {
						processedRecord["control_ids"] = controlIDs // Add the IDs as []uint
					}
				}

			case "title_en", "title_ar":
				if !ok || value == "" {
					rowErrors[header] = "Field is required"
				}

			case "description":
				// Optional field, no validation needed
			}
		}

		// Add errors or success record
		if len(rowErrors) == 0 {
			processedRecord["errors"] = nil
		} else {
			// Ensure all keys are present in the `errors` map
			for _, header := range headers {
				if _, exists := rowErrors[header]; !exists {
					rowErrors[header] = "" // Empty string for other fields
				}
			}
			processedRecord["errors"] = rowErrors
		}
		validatedData = append(validatedData, processedRecord)
	}

	return validatedData, nil
}

func (ps *ProcessService) UpdateProcessRecoveryObjective(process *models.Process, tx utility.DatabaseTransaction) error {
	err := ps.repo.Update(process, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}
