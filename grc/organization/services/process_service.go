package services

import (
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"strings"

	commonDTOs "jethur/grc/common/dtos"
	commonServices "jethur/grc/common/services"
	logServices "jethur/grc/logs/services"
	masterServices "jethur/grc/masters/services"
	"jethur/grc/models"
	organizationDTOs "jethur/grc/organization/dtos"
	organizationRepositories "jethur/grc/organization/repositories"
	settingsDTOs "jethur/grc/settings/dtos"
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
	locationSvc          masterServices.LocationServiceInterface
	reviewFrequencySvc   masterServices.ReviewFrequencyServiceInterface
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
	locationSvc masterServices.LocationServiceInterface,
	reviewFrequencySvc masterServices.ReviewFrequencyServiceInterface,
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
		locationSvc:          locationSvc,
		reviewFrequencySvc:   reviewFrequencySvc,
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

	isAutomatic, isManual := ps.resolveExecutionAutomation(processFormDTO.ExecutionAutomation)
	isCanBeDoneManually := false
	if processFormDTO.IsCanBeDoneManually != nil {
		isCanBeDoneManually = *processFormDTO.IsCanBeDoneManually
	}

	process := &models.Process{
		ProcessGroupID:                        processFormDTO.ProcessGroupID,
		ResponsibleUserID:                     processFormDTO.OwnerID,
		ResponsibleUnitID:                     processFormDTO.OwnerUnitID,
		LifecycleStatusID:                     &lifecycleStatus.ID,
		ReferenceNumber:                       referenceNumber,
		TitleEn:                               processFormDTO.TitleEn,
		TitleAr:                               processFormDTO.TitleAr,
		Description:                           processFormDTO.Description,
		DivisionID:                            processFormDTO.DivisionID,
		IsAutomatic:                           isAutomatic,
		IsManual:                              isManual,
		ProcessTypeID:                         processFormDTO.ProcessTypeID,
		LocationID:                            processFormDTO.LocationID,
		AlternativeLocation:                   processFormDTO.AlternativeLocation,
		RushHour:                              processFormDTO.RushHour,
		ReviewFrequencyID:                     processFormDTO.ReviewFrequencyID,
		IsCanBeDoneManually:                   isCanBeDoneManually,
		TimeToPerformProcedureManually:        processFormDTO.TimeToPerformProcedureManually,
		StepsToPerformProcedureManually:       processFormDTO.StepsToPerformProcedureManually,
		RequirementToPerformProcedureManually: processFormDTO.RequirementToPerformProcedureManually,
		CreatedBYID:                           currentUserID,
	}

	id, err := ps.repo.Create(process, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	// Create association records
	if err := ps.createProcessAssociations(id, currentUserID, processFormDTO, tx); err != nil {
		tx.Rollback()
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

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
			processGroupTitleEn := ""
			if existingProcess.ProcessGroup != nil {
				processGroupTitleEn = existingProcess.ProcessGroup.TitleEn
			}
			emailData = map[string]interface{}{
				"FirstNameEn":            responsibleUser.FirstNameEn,
				"LastNameEn":             responsibleUser.LastNameEn,
				"ProcessReferenceNumber": existingProcess.ReferenceNumber,
				"ProcessName":            existingProcess.TitleEn,
				"ProcessGroup":           processGroupTitleEn,
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

	process, err := ps.GetByID(id, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	oldOwnerID := process.ResponsibleUserID

	if err := ps.validateProcessFormDTO(processFormDTO); err != nil {
		tx.Rollback()
		return err
	}

	isAutomatic, isManual := ps.resolveExecutionAutomation(processFormDTO.ExecutionAutomation)
	isCanBeDoneManually := false
	if processFormDTO.IsCanBeDoneManually != nil {
		isCanBeDoneManually = *processFormDTO.IsCanBeDoneManually
	}

	process.ProcessGroupID = processFormDTO.ProcessGroupID
	process.ResponsibleUserID = processFormDTO.OwnerID
	process.ResponsibleUnitID = processFormDTO.OwnerUnitID
	process.TitleEn = processFormDTO.TitleEn
	process.TitleAr = processFormDTO.TitleAr
	process.Description = processFormDTO.Description
	process.DivisionID = processFormDTO.DivisionID
	process.IsAutomatic = isAutomatic
	process.IsManual = isManual
	process.ProcessTypeID = processFormDTO.ProcessTypeID
	process.LocationID = processFormDTO.LocationID
	process.AlternativeLocation = processFormDTO.AlternativeLocation
	process.RushHour = processFormDTO.RushHour
	process.ReviewFrequencyID = processFormDTO.ReviewFrequencyID
	process.IsCanBeDoneManually = isCanBeDoneManually
	process.TimeToPerformProcedureManually = processFormDTO.TimeToPerformProcedureManually
	process.StepsToPerformProcedureManually = processFormDTO.StepsToPerformProcedureManually
	process.RequirementToPerformProcedureManually = processFormDTO.RequirementToPerformProcedureManually
	process.UpdatedBYID = &currentUserID

	err = ps.repo.Update(process, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Synchronize association records
	if err := ps.syncProcessAssociations(id, currentUserID, processFormDTO, tx); err != nil {
		tx.Rollback()
		return err
	}

	// Collect existing control IDs
	existingProcess, err := ps.GetFullDetailsByID(id, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

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
			"process_group":                       oldProcess.ProcessGroup,
			"owner_unit":                          oldProcess.ResponsibleUnit,
			"owner":                               oldOwnerName,
			"lifecycle_status":                    oldProcess.LifecycleStatus,
			"title_en":                            oldProcess.TitleEn,
			"title_ar":                            oldProcess.TitleAr,
			"description":                         oldProcess.Description,
			"division":                            oldProcess.Division,
			"is_automatic":                        oldProcess.IsAutomatic,
			"is_manual":                           oldProcess.IsManual,
			"process_type":                        oldProcess.ProcessType,
			"location":                            oldProcess.Location,
			"alternative_location":                oldProcess.AlternativeLocation,
			"rush_hour":                           oldProcess.RushHour,
			"review_frequency":                    oldProcess.ReviewFrequency,
			"is_can_be_done_manually":             oldProcess.IsCanBeDoneManually,
			"time_to_perform_procedure_manually":  oldProcess.TimeToPerformProcedureManually,
			"steps_to_perform_procedure_manually": oldProcess.StepsToPerformProcedureManually,
			"requirement_to_perform_procedure_manually": oldProcess.RequirementToPerformProcedureManually,
		}

		newData := map[string]interface{}{
			"process_group":                       existingProcess.ProcessGroup,
			"owner_unit":                          existingProcess.ResponsibleUnit,
			"owner":                               newOwnerName,
			"lifecycle_status":                    existingProcess.LifecycleStatus,
			"title_en":                            existingProcess.TitleEn,
			"title_ar":                            existingProcess.TitleAr,
			"description":                         existingProcess.Description,
			"division":                            existingProcess.Division,
			"is_automatic":                        existingProcess.IsAutomatic,
			"is_manual":                           existingProcess.IsManual,
			"process_type":                        existingProcess.ProcessType,
			"location":                            existingProcess.Location,
			"alternative_location":                existingProcess.AlternativeLocation,
			"rush_hour":                           existingProcess.RushHour,
			"review_frequency":                    existingProcess.ReviewFrequency,
			"is_can_be_done_manually":             existingProcess.IsCanBeDoneManually,
			"time_to_perform_procedure_manually":  existingProcess.TimeToPerformProcedureManually,
			"steps_to_perform_procedure_manually": existingProcess.StepsToPerformProcedureManually,
			"requirement_to_perform_procedure_manually": existingProcess.RequirementToPerformProcedureManually,
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

					processGroupTitleEn := ""
					if existingProcess.ProcessGroup != nil {
						processGroupTitleEn = existingProcess.ProcessGroup.TitleEn
					}
					oldOwnerEmailData := map[string]interface{}{
						"FirstNameEn":            oldOwnerUser.FirstNameEn,
						"LastNameEn":             oldOwnerUser.LastNameEn,
						"ProcessReferenceNumber": existingProcess.ReferenceNumber,
						"ProcessName":            existingProcess.TitleEn,
						"ProcessGroup":           processGroupTitleEn,
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

				newOwnerProcessGroup := ""
				if existingProcess.ProcessGroup != nil {
					newOwnerProcessGroup = existingProcess.ProcessGroup.TitleEn
				}
				emailData := map[string]interface{}{
					"FirstNameEn":            responsibleUser.FirstNameEn,
					"LastNameEn":             responsibleUser.LastNameEn,
					"ProcessReferenceNumber": existingProcess.ReferenceNumber,
					"ProcessName":            existingProcess.TitleEn,
					"ProcessGroup":           newOwnerProcessGroup,
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

				deptOwnerProcessGroup := ""
				if existingProcess.ProcessGroup != nil {
					deptOwnerProcessGroup = existingProcess.ProcessGroup.TitleEn
				}
				emailData := map[string]interface{}{
					"DepartmentNameEn":       department.TitleEn,
					"DepartmentNameAr":       department.TitleAr,
					"ProcessReferenceNumber": existingProcess.ReferenceNumber,
					"ProcessName":            existingProcess.TitleEn,
					"ProcessGroup":           deptOwnerProcessGroup,
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
		processGroupTitle := ""
		if process.ProcessGroup != nil {
			processGroupTitle = process.ProcessGroup.TitleEn
		}
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
			ProcessGroup:    processGroupTitle,
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
	// Validate ProcessGroupID only if provided (nullable)
	if processFormDTO.ProcessGroupID != nil {
		if _, err := ps.processGroupSvc.GetByID(*processFormDTO.ProcessGroupID); err != nil {
			return utility.NewValidationError(fmt.Errorf("process group does not exist"), "not_exist", "process group")
		}
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

	// Validate Division if provided
	if processFormDTO.DivisionID != nil {
		if _, err := ps.organizationUnitSvc.GetByID(*processFormDTO.DivisionID); err != nil {
			return utility.NewValidationError(fmt.Errorf("division does not exist"), "not_exist", "division")
		}
	}

	// Validate Execution Automation
	if processFormDTO.ExecutionAutomation == nil || *processFormDTO.ExecutionAutomation == "" {
		return utility.NewValidationError(fmt.Errorf("execution automation is required"), "required", "execution_automation")
	}
	if *processFormDTO.ExecutionAutomation != "automatic" && *processFormDTO.ExecutionAutomation != "manual" {
		return utility.NewValidationError(fmt.Errorf("execution automation must be automatic or manual"), "invalid", "execution_automation")
	}
	if *processFormDTO.ExecutionAutomation == "automatic" {
		if processFormDTO.SystemRegisterID == nil {
			return utility.NewValidationError(fmt.Errorf("associated system is required when execution automation is automatic"), "required", "system_register_id")
		}
	}

	// Validate ProcessType if provided
	if processFormDTO.ProcessTypeID != nil {
		if _, err := ps.repo.GetBCMProcessTypeByID(*processFormDTO.ProcessTypeID); err != nil {
			return utility.NewValidationError(fmt.Errorf("process type does not exist"), "not_exist", "process_type")
		}
	}

	// Validate Location if provided
	if processFormDTO.LocationID != nil {
		if _, err := ps.locationSvc.GetByID(*processFormDTO.LocationID); err != nil {
			return utility.NewValidationError(fmt.Errorf("location does not exist"), "not_exist", "location")
		}
	}

	// Validate ReviewFrequency if provided
	if processFormDTO.ReviewFrequencyID != nil {
		if _, err := ps.reviewFrequencySvc.GetByID(*processFormDTO.ReviewFrequencyID); err != nil {
			return utility.NewValidationError(fmt.Errorf("review frequency does not exist"), "not_exist", "review_frequency")
		}
	}

	// Validate associated execution records
	for _, id := range processFormDTO.ExecutionRecordIDs {
		if id == nil {
			continue
		}
		if _, err := ps.repo.GetBCMExecutionRecordByID(*id); err != nil {
			return utility.NewValidationError(fmt.Errorf("execution record does not exist"), "not_exist", "execution_record_ids")
		}
	}

	// Validate associated manual procedure forms
	for _, id := range processFormDTO.ManualProcedureFormIDs {
		if id == nil {
			continue
		}
		if _, err := ps.repo.GetBCMManualProcedureFormByID(*id); err != nil {
			return utility.NewValidationError(fmt.Errorf("manual procedure form does not exist"), "not_exist", "manual_procedure_form_ids")
		}
	}

	// Validate associated system register (single for current use)
	if processFormDTO.SystemRegisterID != nil {
		if _, err := ps.repo.GetSystemRegisterByID(*processFormDTO.SystemRegisterID); err != nil {
			return utility.NewValidationError(fmt.Errorf("system register does not exist"), "not_exist", "system_register_id")
		}
	}

	// Validate manual execution fields if process can be done manually
	if processFormDTO.IsCanBeDoneManually != nil && *processFormDTO.IsCanBeDoneManually {
		if processFormDTO.TimeToPerformProcedureManually == nil || *processFormDTO.TimeToPerformProcedureManually == "" {
			return utility.NewValidationError(fmt.Errorf("time to perform procedure manually is required when process can be done manually"), "required", "time_to_perform_procedure_manually")
		}
		if processFormDTO.StepsToPerformProcedureManually == nil || *processFormDTO.StepsToPerformProcedureManually == "" {
			return utility.NewValidationError(fmt.Errorf("steps to perform procedure manually is required when process can be done manually"), "required", "steps_to_perform_procedure_manually")
		}
	}

	return nil
}

func (ps *ProcessService) resolveExecutionAutomation(executionAutomation *string) (bool, bool) {
	if executionAutomation == nil {
		return false, false
	}
	switch *executionAutomation {
	case "automatic":
		return true, false
	case "manual":
		return false, true
	default:
		return false, false
	}
}

func (ps *ProcessService) createProcessAssociations(processID uint, currentUserID uint, processFormDTO *organizationDTOs.ProcessFormDTO, tx utility.DatabaseTransaction) error {
	for _, recordID := range processFormDTO.ExecutionRecordIDs {
		if recordID == nil {
			continue
		}
		record := &models.ProcessExecutionRecord{
			ProcessID:            processID,
			BCMExecutionRecordID: *recordID,
			CreatedBYID:          currentUserID,
		}
		if err := ps.repo.CreateProcessExecutionRecord(record, tx); err != nil {
			return err
		}
	}

	for _, formID := range processFormDTO.ManualProcedureFormIDs {
		if formID == nil {
			continue
		}
		form := &models.ProcessManualProcedureForm{
			ProcessID:                processID,
			BCMManualProcedureFormID: *formID,
			CreatedBYID:              currentUserID,
		}
		if err := ps.repo.CreateProcessManualProcedureForm(form, tx); err != nil {
			return err
		}
	}

	if processFormDTO.SystemRegisterID != nil {
		sr := &models.ProcessSystemRegister{
			ProcessID:        processID,
			SystemRegisterID: *processFormDTO.SystemRegisterID,
			CreatedBYID:      currentUserID,
		}
		if err := ps.repo.CreateProcessSystemRegister(sr, tx); err != nil {
			return err
		}
	}

	return nil
}

func (ps *ProcessService) syncProcessAssociations(processID uint, currentUserID uint, processFormDTO *organizationDTOs.ProcessFormDTO, tx utility.DatabaseTransaction) error {
	if err := ps.repo.DeleteProcessExecutionRecordsByProcessID(processID, tx); err != nil {
		return err
	}
	if err := ps.repo.DeleteProcessManualProcedureFormsByProcessID(processID, tx); err != nil {
		return err
	}
	if err := ps.repo.DeleteProcessSystemRegistersByProcessID(processID, tx); err != nil {
		return err
	}

	return ps.createProcessAssociations(processID, currentUserID, processFormDTO, tx)
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
		"total_control_score":       math.Round((float64(implementationStatusCount["implemented"])/float64(totalControls))*10000) / 100,
		"implementation_status":     implementationStatus,
		"efficiency_percentage":     math.Round((float64(efficiencyStatusCount["efficient"])/float64(totalControls))*10000) / 100,
		"implementation_percentage": math.Round((float64(implementationStatusCount["implemented"])/float64(totalControls))*10000) / 100,
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

	headers := []string{
		"Ref. No.",
		"Title",
		"Description",
		"Responsible Unit",
		"Process Owner",
		"Type",
		"Location",
		"Alternative Location",
		"Rush Hour",
		"Frequency",
		"Associated System",
		"Execution Automation",
		"Can be done manually",
		"Time to perform procedure manually",
		"Requirement to perform procedure manually",
		"Steps to perform procedure manually",
		"Created By",
		"Status",
	}

	rows := make([][]interface{}, len(data))
	for i, process := range data {
		automation := ""
		if process.IsAutomatic && process.IsManual {
			automation = "Automatic & Manual"
		} else if process.IsAutomatic {
			automation = "Automatic"
		} else if process.IsManual {
			automation = "Manual"
		}

		canBeDoneManually := "No"
		if process.IsCanBeDoneManually {
			canBeDoneManually = "Yes"
		}

		associatedSystem := ""
		if process.SystemRegister != nil {
			if userSetting.Language == "ar" {
				associatedSystem = process.SystemRegister.TitleAr
			} else {
				associatedSystem = process.SystemRegister.TitleEn
			}
		}

		if userSetting.Language == "ar" {
			rows[i] = []interface{}{
				process.ReferenceNumber,
				process.TitleAr,
				utility.GetStringValue(process.Description),
				utility.GetStringValue(process.OwnerUnitTitleAr),
				fmt.Sprintf("%s %s", utility.GetStringValue(process.OwnerFirstNameAr), utility.GetStringValue(process.OwnerLastNameAr)),
				utility.GetStringValue(process.ProcessTypeTitleAr),
				utility.GetStringValue(process.LocationTitleAr),
				utility.GetStringValue(process.AlternativeLocation),
				utility.GetStringValue(process.RushHour),
				utility.GetStringValue(process.ReviewFrequencyTitleAr),
				associatedSystem,
				automation,
				canBeDoneManually,
				utility.GetStringValue(process.TimeToPerformProcedureManually),
				utility.GetStringValue(process.RequirementToPerformProcedureManually),
				utility.GetStringValue(process.StepsToPerformProcedureManually),
				fmt.Sprintf("%s %s", process.CreatedUserFirstNameAr, process.CreatedUserLastNameAr),
				process.LifecycleStatusTitleAr,
			}
		} else {
			rows[i] = []interface{}{
				process.ReferenceNumber,
				process.TitleEn,
				utility.GetStringValue(process.Description),
				utility.GetStringValue(process.OwnerUnitTitleEn),
				fmt.Sprintf("%s %s", utility.GetStringValue(process.OwnerFirstNameEn), utility.GetStringValue(process.OwnerLastNameEn)),
				utility.GetStringValue(process.ProcessTypeTitleEn),
				utility.GetStringValue(process.LocationTitleEn),
				utility.GetStringValue(process.AlternativeLocation),
				utility.GetStringValue(process.RushHour),
				utility.GetStringValue(process.ReviewFrequencyTitleEn),
				associatedSystem,
				automation,
				canBeDoneManually,
				utility.GetStringValue(process.TimeToPerformProcedureManually),
				utility.GetStringValue(process.RequirementToPerformProcedureManually),
				utility.GetStringValue(process.StepsToPerformProcedureManually),
				fmt.Sprintf("%s %s", process.CreatedUserFirstNameEn, process.CreatedUserLastNameEn),
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
	// Headers for the Template sheet — column names must match reference sheet names for dropdowns
	templateHeaders := []string{
		"TitleEn", "TitleAr", "Description",
		"Division", "Department", "ResponsibleUser",
		"ExecutionAutomation", "Type", "Location",
		"AlternativeLocation", "RushHour", "Frequency",
		"ExecutionRecords", "SystemRegister", "CanBeDoneManually",
		"TimeToPerformProcedureManually", "StepsToPerformProcedureManually",
		"RequirementToPerformProcedureManually", "ManualProcedureForms",
	}
	templateData := [][]interface{}{} // Empty template data

	// Sheet: ResponsibleUser (name matches template column for dropdown)
	users, err := ps.userSvc.GetAllUsers()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch users: %v", err)
	}
	responsibleUserData := [][]interface{}{{"ResponsibleUser"}}
	for _, user := range users {
		fullName := strings.TrimSpace(fmt.Sprintf("%s %s", user.FirstNameEN, user.LastNameEN))
		if fullName != "" {
			responsibleUserData = append(responsibleUserData, []interface{}{fullName})
		}
	}

	// Sheet: Division (level=1 org units)
	divisionParams := &settingsDTOs.OrganizationUnitListQueryParamsDTO{
		Sort: "title_en", SortBy: "title_en", Keyword: "", ParentUnitIDs: []uint{}, Level: 1,
	}
	divisionUnits, err := ps.organizationUnitSvc.GetAll(divisionParams)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch divisions: %v", err)
	}
	divisionData := [][]interface{}{{"Division"}}
	for _, unit := range divisionUnits {
		if unit.TitleEn != "" {
			divisionData = append(divisionData, []interface{}{unit.TitleEn})
		}
	}

	// Sheet: Department (level=2 org units)
	departmentParams := &settingsDTOs.OrganizationUnitListQueryParamsDTO{
		Sort: "title_en", SortBy: "title_en", Keyword: "", ParentUnitIDs: []uint{}, Level: 2,
	}
	departmentUnits, err := ps.organizationUnitSvc.GetAll(departmentParams)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch departments: %v", err)
	}
	departmentData := [][]interface{}{{"Department"}}
	for _, unit := range departmentUnits {
		if unit.TitleEn != "" {
			departmentData = append(departmentData, []interface{}{unit.TitleEn})
		}
	}

	// Sheet: Location
	locations, err := ps.locationSvc.GetAll("", "", "", "active")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch locations: %v", err)
	}
	locationData := [][]interface{}{{"Location"}}
	for _, loc := range locations {
		if loc.TitleEn != "" {
			locationData = append(locationData, []interface{}{loc.TitleEn})
		}
	}

	// Sheet: Frequency
	frequencies, err := ps.reviewFrequencySvc.GetAll("", "", "", "active", "")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch review frequencies: %v", err)
	}
	frequencyData := [][]interface{}{{"Frequency"}}
	for _, freq := range frequencies {
		if freq.TitleEn != "" {
			frequencyData = append(frequencyData, []interface{}{freq.TitleEn})
		}
	}

	// Sheet: Type (BCMProcessType)
	processTypes, err := ps.repo.GetAllBCMProcessTypes()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch process types: %v", err)
	}
	typeData := [][]interface{}{{"Type"}}
	for _, pt := range processTypes {
		if pt.TitleEn != "" {
			typeData = append(typeData, []interface{}{pt.TitleEn})
		}
	}

	// Sheet: ExecutionRecords (BCMExecutionRecord)
	execRecords, err := ps.repo.GetAllBCMExecutionRecords()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch execution records: %v", err)
	}
	execRecordsData := [][]interface{}{{"ExecutionRecords"}}
	for _, er := range execRecords {
		if er.TitleEn != "" {
			execRecordsData = append(execRecordsData, []interface{}{er.TitleEn})
		}
	}

	// Sheet: SystemRegister
	systemRegisters, err := ps.repo.GetAllSystemRegisters()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch system registers: %v", err)
	}
	systemRegisterData := [][]interface{}{{"SystemRegister"}}
	for _, sr := range systemRegisters {
		if sr.TitleEn != "" {
			systemRegisterData = append(systemRegisterData, []interface{}{sr.TitleEn})
		}
	}

	// Sheet: ManualProcedureForms (BCMManualProcedureForm)
	manualForms, err := ps.repo.GetAllBCMManualProcedureForms()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch manual procedure forms: %v", err)
	}
	manualFormsData := [][]interface{}{{"ManualProcedureForms"}}
	for _, mf := range manualForms {
		if mf.TitleEn != "" {
			manualFormsData = append(manualFormsData, []interface{}{mf.TitleEn})
		}
	}

	// Fixed-value sheets for enum columns
	executionAutomationData := [][]interface{}{{"ExecutionAutomation"}, {"automatic"}, {"manual"}}
	canBeDoneManuallyData := [][]interface{}{{"CanBeDoneManually"}, {"yes"}, {"no"}}

	dropDownConfig := map[string]bool{
		"Division":             true,
		"Department":           true,
		"ResponsibleUser":      true,
		"ExecutionAutomation":  true,
		"Type":                 true,
		"Location":             true,
		"Frequency":            true,
		"ExecutionRecords":     true,
		"SystemRegister":       true,
		"ManualProcedureForms": true,
		"CanBeDoneManually":    true,
	}

	file, err := utility.GenerateExcelFileWithSheetsAndDropDown(map[string][][]interface{}{
		"Template":             append([][]interface{}{utility.ConvertToInterfaceSlice(templateHeaders)}, templateData...),
		"Division":             divisionData,
		"Department":           departmentData,
		"ResponsibleUser":      responsibleUserData,
		"ExecutionAutomation":  executionAutomationData,
		"Type":                 typeData,
		"Location":             locationData,
		"Frequency":            frequencyData,
		"ExecutionRecords":     execRecordsData,
		"SystemRegister":       systemRegisterData,
		"ManualProcedureForms": manualFormsData,
		"CanBeDoneManually":    canBeDoneManuallyData,
	}, dropDownConfig)
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

	if len(rows) == 0 {
		return nil, fmt.Errorf("empty Excel file")
	}

	inputHeaders := []string{
		"TitleEn", "TitleAr", "Description",
		"Division", "Department", "ResponsibleUser",
		"ExecutionAutomation", "Type", "Location",
		"AlternativeLocation", "RushHour", "Frequency",
		"ExecutionRecords", "SystemRegister", "CanBeDoneManually",
		"TimeToPerformProcedureManually", "StepsToPerformProcedureManually",
		"RequirementToPerformProcedureManually", "ManualProcedureForms",
	}
	snakeCaseHeaders := []string{
		"title_en", "title_ar", "description",
		"division", "department", "responsible_user",
		"execution_automation", "type", "location",
		"alternative_location", "rush_hour", "frequency",
		"execution_records", "system_register", "can_be_done_manually",
		"time_to_perform_procedure_manually", "steps_to_perform_procedure_manually",
		"requirement_to_perform_procedure_manually", "manual_procedure_forms",
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
				record[header] = ""
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
		processFormDTO := &organizationDTOs.ProcessFormDTO{
			OwnerUnitID:                           record["owner_unit_id"].(uint),
			TitleEn:                               record["title_en"].(string),
			TitleAr:                               record["title_ar"].(string),
			Description:                           utility.GetOptionalStringValue(record["description"]),
			ExecutionAutomation:                   utility.GetOptionalStringValue(record["execution_automation"]),
			AlternativeLocation:                   utility.GetOptionalStringValue(record["alternative_location"]),
			RushHour:                              utility.GetOptionalStringValue(record["rush_hour"]),
			TimeToPerformProcedureManually:        utility.GetOptionalStringValue(record["time_to_perform_procedure_manually"]),
			StepsToPerformProcedureManually:       utility.GetOptionalStringValue(record["steps_to_perform_procedure_manually"]),
			RequirementToPerformProcedureManually: utility.GetOptionalStringValue(record["requirement_to_perform_procedure_manually"]),
		}

		if ownerID, ok := record["responsible_user_id"].(uint); ok {
			processFormDTO.OwnerID = &ownerID
		}
		if divisionID, ok := record["division_id"].(uint); ok {
			processFormDTO.DivisionID = &divisionID
		}
		if processTypeID, ok := record["process_type_id"].(uint); ok {
			processFormDTO.ProcessTypeID = &processTypeID
		}
		if locationID, ok := record["location_id"].(uint); ok {
			processFormDTO.LocationID = &locationID
		}
		if frequencyID, ok := record["frequency_id"].(uint); ok {
			processFormDTO.ReviewFrequencyID = &frequencyID
		}
		if systemRegisterID, ok := record["system_register_id"].(uint); ok {
			processFormDTO.SystemRegisterID = &systemRegisterID
		}
		if canBeDoneManually, ok := record["can_be_done_manually"].(bool); ok {
			processFormDTO.IsCanBeDoneManually = &canBeDoneManually
		}
		if executionRecordIDs, ok := record["execution_record_ids"].([]uint); ok {
			for _, id := range executionRecordIDs {
				id := id
				processFormDTO.ExecutionRecordIDs = append(processFormDTO.ExecutionRecordIDs, &id)
			}
		}
		if manualFormIDs, ok := record["manual_procedure_form_ids"].([]uint); ok {
			for _, id := range manualFormIDs {
				id := id
				processFormDTO.ManualProcedureFormIDs = append(processFormDTO.ManualProcedureFormIDs, &id)
			}
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
	headers := []string{
		"title_en", "title_ar", "description",
		"division", "department", "responsible_user",
		"execution_automation", "type", "location",
		"alternative_location", "rush_hour", "frequency",
		"execution_records", "system_register", "can_be_done_manually",
		"time_to_perform_procedure_manually", "steps_to_perform_procedure_manually",
		"requirement_to_perform_procedure_manually", "manual_procedure_forms",
	}

	// Fetch user full names → IDs
	users, err := ps.userSvc.GetAllUsers()
	if err != nil {
		return nil, fmt.Errorf("error fetching users: %v", err)
	}
	userFullNames := make(map[string]uint)
	for _, user := range users {
		fullName := strings.TrimSpace(user.FirstNameEN + " " + user.LastNameEN)
		userFullNames[fullName] = uint(user.ID)
	}

	// Fetch division titles (level=1) → IDs
	divisionParams := &settingsDTOs.OrganizationUnitListQueryParamsDTO{
		Sort: "", SortBy: "", Keyword: "", ParentUnitIDs: []uint{}, Level: 1,
	}
	divisionUnits, err := ps.organizationUnitSvc.GetAll(divisionParams)
	if err != nil {
		return nil, fmt.Errorf("error fetching divisions: %v", err)
	}
	divisionTitles := make(map[string]uint)
	for _, unit := range divisionUnits {
		divisionTitles[strings.TrimSpace(unit.TitleEn)] = unit.ID
	}

	// Fetch department titles (level=2) → IDs
	departmentParams := &settingsDTOs.OrganizationUnitListQueryParamsDTO{
		Sort: "", SortBy: "", Keyword: "", ParentUnitIDs: []uint{}, Level: 2,
	}
	departmentUnits, err := ps.organizationUnitSvc.GetAll(departmentParams)
	if err != nil {
		return nil, fmt.Errorf("error fetching departments: %v", err)
	}
	departmentTitles := make(map[string]uint)
	for _, unit := range departmentUnits {
		departmentTitles[strings.TrimSpace(unit.TitleEn)] = unit.ID
	}

	// Fetch location titles → IDs
	locations, err := ps.locationSvc.GetAll("", "", "", "active")
	if err != nil {
		return nil, fmt.Errorf("error fetching locations: %v", err)
	}
	locationTitles := make(map[string]uint)
	for _, loc := range locations {
		locationTitles[strings.TrimSpace(loc.TitleEn)] = loc.ID
	}

	// Fetch review frequency titles → IDs
	frequencies, err := ps.reviewFrequencySvc.GetAll("", "", "", "active", "")
	if err != nil {
		return nil, fmt.Errorf("error fetching review frequencies: %v", err)
	}
	frequencyTitles := make(map[string]uint)
	for _, freq := range frequencies {
		frequencyTitles[strings.TrimSpace(freq.TitleEn)] = freq.ID
	}

	// Fetch process type titles → IDs
	processTypes, err := ps.repo.GetAllBCMProcessTypes()
	if err != nil {
		return nil, fmt.Errorf("error fetching process types: %v", err)
	}
	processTypeTitles := make(map[string]uint)
	for _, pt := range processTypes {
		processTypeTitles[strings.TrimSpace(pt.TitleEn)] = pt.ID
	}

	// Fetch execution record titles → IDs
	execRecords, err := ps.repo.GetAllBCMExecutionRecords()
	if err != nil {
		return nil, fmt.Errorf("error fetching execution records: %v", err)
	}
	execRecordTitles := make(map[string]uint)
	for _, er := range execRecords {
		execRecordTitles[strings.TrimSpace(er.TitleEn)] = er.ID
	}

	// Fetch system register titles → IDs
	systemRegisters, err := ps.repo.GetAllSystemRegisters()
	if err != nil {
		return nil, fmt.Errorf("error fetching system registers: %v", err)
	}
	systemRegisterTitles := make(map[string]uint)
	for _, sr := range systemRegisters {
		systemRegisterTitles[strings.TrimSpace(sr.TitleEn)] = sr.ID
	}

	// Fetch manual procedure form titles → IDs
	manualForms, err := ps.repo.GetAllBCMManualProcedureForms()
	if err != nil {
		return nil, fmt.Errorf("error fetching manual procedure forms: %v", err)
	}
	manualFormTitles := make(map[string]uint)
	for _, mf := range manualForms {
		manualFormTitles[strings.TrimSpace(mf.TitleEn)] = mf.ID
	}

	var validatedData []map[string]interface{}

	for _, record := range input {
		rowErrors := make(map[string]interface{})
		processedRecord := make(map[string]interface{})

		for key, value := range record {
			processedRecord[key] = value
		}

		for _, header := range headers {
			value, _ := record[header].(string)
			value = strings.TrimSpace(value)

			switch header {
			case "title_en", "title_ar":
				if value == "" {
					rowErrors[header] = "Field is required"
				}

			case "department":
				if value == "" {
					rowErrors["department"] = "Field is required"
				} else if id, exists := departmentTitles[value]; !exists {
					rowErrors["department"] = fmt.Sprintf("Invalid department: %s", value)
				} else {
					processedRecord["owner_unit_id"] = id
					processedRecord["department_id"] = id
				}

			case "division":
				if value != "" {
					if id, exists := divisionTitles[value]; !exists {
						rowErrors["division"] = fmt.Sprintf("Invalid division: %s", value)
					} else {
						processedRecord["division_id"] = id
					}
				}

			case "responsible_user":
				if value != "" {
					if id, exists := userFullNames[value]; !exists {
						rowErrors["responsible_user"] = fmt.Sprintf("Invalid user: %s", value)
					} else {
						processedRecord["responsible_user_id"] = id
					}
				}

			case "execution_automation":
				if value == "" {
					rowErrors["execution_automation"] = "Field is required"
				} else {
					normalized := strings.ToLower(value)
					if normalized != "automatic" && normalized != "manual" {
						rowErrors["execution_automation"] = "Must be 'automatic' or 'manual'"
					} else {
						processedRecord["execution_automation"] = normalized
					}
				}

			case "location":
				if value != "" {
					if id, exists := locationTitles[value]; !exists {
						rowErrors["location"] = fmt.Sprintf("Invalid location: %s", value)
					} else {
						processedRecord["location_id"] = id
					}
				}

			case "frequency":
				if value != "" {
					if id, exists := frequencyTitles[value]; !exists {
						rowErrors["frequency"] = fmt.Sprintf("Invalid frequency: %s", value)
					} else {
						processedRecord["frequency_id"] = id
					}
				}

			case "can_be_done_manually":
				var isCanBeDoneManually bool
				switch v := record[header].(type) {
				case bool:
					isCanBeDoneManually = v
				case string:
					v = strings.TrimSpace(v)
					if v == "" {
						rowErrors["can_be_done_manually"] = "Field is required"
						break
					}
					normalized := strings.ToLower(v)
					if normalized != "yes" && normalized != "no" {
						rowErrors["can_be_done_manually"] = "Must be 'yes' or 'no'"
						break
					}
					isCanBeDoneManually = normalized == "yes"
				default:
					rowErrors["can_be_done_manually"] = "Must be 'yes' or 'no'"
				}
				if _, hasError := rowErrors["can_be_done_manually"]; !hasError {
					processedRecord["can_be_done_manually"] = isCanBeDoneManually
				}

			case "type":
				if value != "" {
					if id, exists := processTypeTitles[value]; !exists {
						rowErrors["type"] = fmt.Sprintf("Invalid type: %s", value)
					} else {
						processedRecord["process_type_id"] = id
					}
				}

			case "execution_records":
				if value != "" {
					titles := utility.ParseCommaSeparatedValues(value)
					var ids []uint
					var invalidVals []string
					for _, t := range titles {
						t = strings.TrimSpace(t)
						if id, exists := execRecordTitles[t]; exists {
							ids = append(ids, id)
						} else {
							invalidVals = append(invalidVals, t)
						}
					}
					if len(invalidVals) > 0 {
						rowErrors["execution_records"] = fmt.Sprintf("Invalid execution records: %s", strings.Join(invalidVals, ", "))
					} else {
						processedRecord["execution_record_ids"] = ids
					}
				}

			case "system_register":
				if value != "" {
					if id, exists := systemRegisterTitles[value]; !exists {
						rowErrors["system_register"] = fmt.Sprintf("Invalid system register: %s", value)
					} else {
						processedRecord["system_register_id"] = id
					}
				}

			case "manual_procedure_forms":
				if value != "" {
					titles := utility.ParseCommaSeparatedValues(value)
					var ids []uint
					var invalidVals []string
					for _, t := range titles {
						t = strings.TrimSpace(t)
						if id, exists := manualFormTitles[t]; exists {
							ids = append(ids, id)
						} else {
							invalidVals = append(invalidVals, t)
						}
					}
					if len(invalidVals) > 0 {
						rowErrors["manual_procedure_forms"] = fmt.Sprintf("Invalid manual procedure forms: %s", strings.Join(invalidVals, ", "))
					} else {
						processedRecord["manual_procedure_form_ids"] = ids
					}
				}
			}
		}

		// Cross-field: execution_automation=automatic requires system_register
		if autoVal, ok := processedRecord["execution_automation"].(string); ok && autoVal == "automatic" {
			if _, hasSR := processedRecord["system_register_id"]; !hasSR {
				rowErrors["system_register"] = "System register is required when execution automation is 'automatic'"
			}
		}

		// Cross-field: can_be_done_manually=yes requires time and steps fields
		if canBeDone, ok := processedRecord["can_be_done_manually"].(bool); ok && canBeDone {
			if v, _ := record["time_to_perform_procedure_manually"].(string); strings.TrimSpace(v) == "" {
				rowErrors["time_to_perform_procedure_manually"] = "Field is required when 'Can Be Done Manually' is yes"
			}
			if v, _ := record["steps_to_perform_procedure_manually"].(string); strings.TrimSpace(v) == "" {
				rowErrors["steps_to_perform_procedure_manually"] = "Field is required when 'Can Be Done Manually' is yes"
			}
		}

		if len(rowErrors) == 0 {
			processedRecord["errors"] = nil
		} else {
			for _, header := range headers {
				if _, exists := rowErrors[header]; !exists {
					rowErrors[header] = ""
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
