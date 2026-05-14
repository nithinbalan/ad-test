package services

import (
	"errors"
	"fmt"
	"os"
	"time"

	"jethur/grc/asset/dtos"
	"jethur/grc/asset/repositories"
	logServices "jethur/grc/logs/services"
	masterServices "jethur/grc/masters/services"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type SystemRegisterServiceInterface interface {
	GetAll(params *dtos.SystemRegisterListQueryParamsDTO) (*dtos.SystemRegisterListResponseDTO, error)
	GetByID(id uint) (*models.SystemRegister, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*dtos.SystemRegisterDetailDTO, error)
	Create(systemRegisterForm *dtos.SystemRegisterFormDTO, currentUserID uint) (uint, error)
	Update(id uint, systemRegisterForm *dtos.SystemRegisterFormDTO, currentUserID uint) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	ExportSystemRegister(data []*dtos.SystemRegisterListDTO, currentUserID uint) (*os.File, error)
}

type SystemRegisterService struct {
	repo            repositories.SystemRegisterRepositoryInterface
	activityLogSvc  logServices.ActivityLogServiceInterface
	userSvc         organizationServices.UserServiceInterface
	designationSvc  organizationServices.DesignationServiceInterface
	locationSvc     masterServices.LocationServiceInterface
	recoveryToolSvc SystemRegisterRecoveryToolServiceInterface
	userSettingSvc  organizationServices.UserSettingServiceInterface
}

func NewSystemRegisterService(
	repo repositories.SystemRegisterRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	designationSvc organizationServices.DesignationServiceInterface,
	locationSvc masterServices.LocationServiceInterface,
	recoveryToolSvc SystemRegisterRecoveryToolServiceInterface,
	userSettingSvc organizationServices.UserSettingServiceInterface,
) SystemRegisterServiceInterface {
	return &SystemRegisterService{
		repo:            repo,
		activityLogSvc:  activityLogSvc,
		userSvc:         userSvc,
		designationSvc:  designationSvc,
		locationSvc:     locationSvc,
		recoveryToolSvc: recoveryToolSvc,
		userSettingSvc:  userSettingSvc,
	}
}

func (srs *SystemRegisterService) GetAll(params *dtos.SystemRegisterListQueryParamsDTO) (*dtos.SystemRegisterListResponseDTO, error) {
	return srs.repo.GetAll(params)
}

func (srs *SystemRegisterService) GetByID(id uint) (*models.SystemRegister, error) {
	systemRegister, err := srs.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return systemRegister, nil
}

func (srs *SystemRegisterService) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*dtos.SystemRegisterDetailDTO, error) {
	var systemRegisterDetailDTO dtos.SystemRegisterDetailDTO

	systemRegister, err := srs.repo.GetFullDetailsByID(id, tx...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	err = mapstructure.Decode(systemRegister, &systemRegisterDetailDTO)
	if err != nil {
		return nil, err
	}

	systemRegisterDetailDTO.CreatedAt = systemRegister.CreatedAt
	return &systemRegisterDetailDTO, nil
}

func (srs *SystemRegisterService) Create(systemRegisterFormDTO *dtos.SystemRegisterFormDTO, currentUserID uint) (uint, error) {
	tx := srs.repo.BeginTransaction()

	if err := srs.validateSystemRegisterFormDTO(systemRegisterFormDTO); err != nil {
		return 0, err
	}

	referenceNumber, err := srs.repo.GenerateReferenceNumber()
	if err != nil {
		return 0, err
	}

	systemRegister := &models.SystemRegister{
		ReferenceNumber:               referenceNumber,
		NumberOfPeopleNecessary:       systemRegisterFormDTO.NumberOfPeopleNecessary,
		BusinessOwnerUnitID:           systemRegisterFormDTO.BusinessOwnerUnitID,
		ITResponsibleUnitID:           systemRegisterFormDTO.ITResponsibleUnitID,
		ThirdPartyDependencies:        systemRegisterFormDTO.ThirdPartyDependencies,
		OutsourcingPartners:           systemRegisterFormDTO.OutsourcingPartners,
		CustomerUser:                  systemRegisterFormDTO.CustomerUser,
		ProcessOperationTeamID:        systemRegisterFormDTO.ProcessOperationTeamID,
		ProcessRecoveryTeamID:         systemRegisterFormDTO.ProcessRecoveryTeamID,
		RecoveryManagerID:             systemRegisterFormDTO.RecoveryManagerID,
		RecoveryManagerContactDetails: systemRegisterFormDTO.RecoveryManagerContactDetails,
		CreatedBYID:                   currentUserID,
		CreatedAt:                     time.Now(),
	}

	systemRegisterID, err := srs.repo.Create(systemRegister, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	// Handle Job Titles
	if len(systemRegisterFormDTO.JobTitleIDs) > 0 {
		var jobTitles []*models.SystemRegisterJobTitle
		for _, jobTitleID := range systemRegisterFormDTO.JobTitleIDs {
			jobTitles = append(jobTitles, &models.SystemRegisterJobTitle{
				SystemRegisterID: systemRegisterID,
				DesignationID:    jobTitleID,
				CreatedBYID:      currentUserID,
				CreatedAt:        time.Now(),
			})
		}

		if err := srs.repo.CreateJobTitles(jobTitles, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	// Handle Operational Locations
	if len(systemRegisterFormDTO.OperationalLocationIDs) > 0 {
		var operationalLocations []*models.SystemRegisterOperationalLocation
		for _, locationID := range systemRegisterFormDTO.OperationalLocationIDs {
			operationalLocations = append(operationalLocations, &models.SystemRegisterOperationalLocation{
				SystemRegisterID:      systemRegisterID,
				OperationalLocationID: locationID,
				CreatedBYID:           currentUserID,
				CreatedAt:             time.Now(),
			})
		}

		if err := srs.repo.CreateOperationalLocations(operationalLocations, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	// Handle Recovery Locations
	if len(systemRegisterFormDTO.RecoveryLocationIDs) > 0 {
		var recoveryLocations []*models.SystemRegisterRecoveryLocation
		for _, locationID := range systemRegisterFormDTO.RecoveryLocationIDs {
			recoveryLocations = append(recoveryLocations, &models.SystemRegisterRecoveryLocation{
				SystemRegisterID:   systemRegisterID,
				RecoveryLocationID: locationID,
				CreatedBYID:        currentUserID,
				CreatedAt:          time.Now(),
			})
		}

		if err := srs.repo.CreateRecoveryLocations(recoveryLocations, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	// Handle Recovery Needed Tools
	if len(systemRegisterFormDTO.RecoveryNeededToolIDs) > 0 {
		var recoveryNeededTools []*models.SystemRegisterRecoveryNeededTool
		for _, toolID := range systemRegisterFormDTO.RecoveryNeededToolIDs {
			recoveryNeededTools = append(recoveryNeededTools, &models.SystemRegisterRecoveryNeededTool{
				SystemRegisterID:     systemRegisterID,
				RecoveryNeededToolID: toolID,
				CreatedBYID:          currentUserID,
				CreatedAt:            time.Now(),
			})
		}

		if err := srs.repo.CreateRecoveryNeededTools(recoveryNeededTools, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	tx.Commit()

	go func() {
		existingSystemRegister, err := srs.GetFullDetailsByID(systemRegisterID)
		if err != nil {
			utility.LogError("Failed to fetch system register", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}
		err = srs.activityLogSvc.LogActivity(
			models.CreateAction,
			models.SystemRegister{}.TableName(),
			systemRegisterID,
			existingSystemRegister,
			currentUserID,
			"asset",
			"system_register",
			"",
			nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.SystemRegister{}.TableName(),
				"recordID":  systemRegisterID,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log system register activity", fields)
		}
	}()

	return systemRegisterID, nil
}

func (srs *SystemRegisterService) Update(id uint, systemRegisterFormDTO *dtos.SystemRegisterFormDTO, currentUserID uint) error {
	tx := srs.repo.BeginTransaction()
	existingSystemRegister, err := srs.repo.GetByID(id)
	if err != nil {
		return err
	}

	oldSystemRegister, err := srs.GetFullDetailsByID(id, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := srs.validateSystemRegisterFormDTO(systemRegisterFormDTO); err != nil {
		return err
	}

	// Update fields
	existingSystemRegister.NumberOfPeopleNecessary = systemRegisterFormDTO.NumberOfPeopleNecessary
	existingSystemRegister.BusinessOwnerUnitID = systemRegisterFormDTO.BusinessOwnerUnitID
	existingSystemRegister.ITResponsibleUnitID = systemRegisterFormDTO.ITResponsibleUnitID
	existingSystemRegister.ThirdPartyDependencies = systemRegisterFormDTO.ThirdPartyDependencies
	existingSystemRegister.OutsourcingPartners = systemRegisterFormDTO.OutsourcingPartners
	existingSystemRegister.CustomerUser = systemRegisterFormDTO.CustomerUser
	existingSystemRegister.ProcessOperationTeamID = systemRegisterFormDTO.ProcessOperationTeamID
	existingSystemRegister.ProcessRecoveryTeamID = systemRegisterFormDTO.ProcessRecoveryTeamID
	existingSystemRegister.RecoveryManagerID = systemRegisterFormDTO.RecoveryManagerID
	existingSystemRegister.RecoveryManagerContactDetails = systemRegisterFormDTO.RecoveryManagerContactDetails
	existingSystemRegister.UpdatedBYID = &currentUserID
	now := time.Now()
	existingSystemRegister.UpdatedAt = &now

	err = srs.repo.Update(existingSystemRegister, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Update Job Titles
	if err := srs.repo.DeleteJobTitles(id, tx); err != nil {
		tx.Rollback()
		return err
	}

	if len(systemRegisterFormDTO.JobTitleIDs) > 0 {
		var jobTitles []*models.SystemRegisterJobTitle
		for _, jobTitleID := range systemRegisterFormDTO.JobTitleIDs {
			jobTitles = append(jobTitles, &models.SystemRegisterJobTitle{
				SystemRegisterID: id,
				DesignationID:    jobTitleID,
				CreatedBYID:      currentUserID,
				CreatedAt:        time.Now(),
			})
		}

		if err := srs.repo.CreateJobTitles(jobTitles, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	// Update Operational Locations
	if err := srs.repo.DeleteOperationalLocations(id, tx); err != nil {
		tx.Rollback()
		return err
	}

	if len(systemRegisterFormDTO.OperationalLocationIDs) > 0 {
		var operationalLocations []*models.SystemRegisterOperationalLocation
		for _, locationID := range systemRegisterFormDTO.OperationalLocationIDs {
			operationalLocations = append(operationalLocations, &models.SystemRegisterOperationalLocation{
				SystemRegisterID:      id,
				OperationalLocationID: locationID,
				CreatedBYID:           currentUserID,
				CreatedAt:             time.Now(),
			})
		}

		if err := srs.repo.CreateOperationalLocations(operationalLocations, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	// Update Recovery Locations
	if err := srs.repo.DeleteRecoveryLocations(id, tx); err != nil {
		tx.Rollback()
		return err
	}

	if len(systemRegisterFormDTO.RecoveryLocationIDs) > 0 {
		var recoveryLocations []*models.SystemRegisterRecoveryLocation
		for _, locationID := range systemRegisterFormDTO.RecoveryLocationIDs {
			recoveryLocations = append(recoveryLocations, &models.SystemRegisterRecoveryLocation{
				SystemRegisterID:   id,
				RecoveryLocationID: locationID,
				CreatedBYID:        currentUserID,
				CreatedAt:          time.Now(),
			})
		}

		if err := srs.repo.CreateRecoveryLocations(recoveryLocations, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	// Update Recovery Needed Tools
	if err := srs.repo.DeleteRecoveryNeededTools(id, tx); err != nil {
		tx.Rollback()
		return err
	}

	if len(systemRegisterFormDTO.RecoveryNeededToolIDs) > 0 {
		var recoveryNeededTools []*models.SystemRegisterRecoveryNeededTool
		for _, toolID := range systemRegisterFormDTO.RecoveryNeededToolIDs {
			recoveryNeededTools = append(recoveryNeededTools, &models.SystemRegisterRecoveryNeededTool{
				SystemRegisterID:     id,
				RecoveryNeededToolID: toolID,
				CreatedBYID:          currentUserID,
				CreatedAt:            time.Now(),
			})
		}

		if err := srs.repo.CreateRecoveryNeededTools(recoveryNeededTools, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()

	go func() {
		existingSystemRegister, err := srs.GetFullDetailsByID(id)
		if err != nil {
			return
		}

		// Prepare log data
		changedData := make(map[string]interface{})

		// log changes
		oldData := map[string]interface{}{
			"number_of_people_necessary": oldSystemRegister.NumberOfPeopleNecessary,
			"third_party_dependencies":   oldSystemRegister.ThirdPartyDependencies,
			"outsourcing_partners":       oldSystemRegister.OutsourcingPartners,
			"customer_user":              oldSystemRegister.CustomerUser,
		}

		newData := map[string]interface{}{
			"number_of_people_necessary": existingSystemRegister.NumberOfPeopleNecessary,
			"third_party_dependencies":   existingSystemRegister.ThirdPartyDependencies,
			"outsourcing_partners":       existingSystemRegister.OutsourcingPartners,
			"customer_user":              existingSystemRegister.CustomerUser,
		}

		// Compare other fields and add them to changedData
		fieldChanges := utility.CompareFieldChanges(oldData, newData)
		for k, v := range fieldChanges {
			changedData[k] = v
		}

		// Log activity if there are any changes
		if len(changedData) > 0 {
			err := srs.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.SystemRegister{}.TableName(),
				id,
				changedData,
				currentUserID,
				"asset",
				"system_register",
				"",
				nil,
			)
			if err != nil {
				fields := map[string]interface{}{
					"action":    models.UpdateAction,
					"tableName": models.SystemRegister{}.TableName(),
					"recordID":  id,
					"userID":    currentUserID,
					"error":     err.Error(),
				}

				utility.LogError("Failed to log activity", fields)
			}
		}
	}()

	return nil
}

func (srs *SystemRegisterService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIds []uint
	var successfulIds []uint
	var deletedSystemRegisters []map[string]interface{}

	tx := srs.repo.BeginTransaction()

	for _, id := range ids {
		systemRegister, err := srs.GetFullDetailsByID(id)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		// Delete related entities
		if err := srs.repo.DeleteJobTitles(id, tx); err != nil {
			tx.Rollback()
			failedIds = append(failedIds, id)
			continue
		}

		if err := srs.repo.DeleteOperationalLocations(id, tx); err != nil {
			tx.Rollback()
			failedIds = append(failedIds, id)
			continue
		}

		if err := srs.repo.DeleteRecoveryLocations(id, tx); err != nil {
			tx.Rollback()
			failedIds = append(failedIds, id)
			continue
		}

		if err := srs.repo.DeleteRecoveryNeededTools(id, tx); err != nil {
			tx.Rollback()
			failedIds = append(failedIds, id)
			continue
		}

		// Delete system register record
		err = srs.repo.Delete(id, tx)
		if err != nil {
			tx.Rollback()
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}

		// Store deleted system register details for logging
		deletedSystemRegisters = append(deletedSystemRegisters, map[string]interface{}{
			"id":               id,
			"reference_number": systemRegister.ReferenceNumber,
		})
	}

	tx.Commit()

	// Log activity in the background
	if len(successfulIds) > 0 {
		go func() {
			err := srs.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.SystemRegister{}.TableName(),
				successfulIds,
				deletedSystemRegisters,
				currentUserID,
				nil,
			)
			if err != nil {
				fmt.Printf("Failed to log bulk delete activity: %v\n", err)
			}
		}()
	}

	return failedIds, nil
}

func (srs *SystemRegisterService) ExportSystemRegister(data []*dtos.SystemRegisterListDTO, currentUserID uint) (*os.File, error) {
	userSetting, err := srs.userSettingSvc.GetUserSetting(currentUserID)
	if err != nil {
		return nil, err
	}

	var headers []string
	if userSetting.Language == "ar" {
		headers = []string{
			"SL No.",
			"Reference Number",
			"Number of People Necessary",
			"Business Owner Unit",
			"IT Responsible Unit",
			"Process Operation Team",
			"Process Recovery Team",
			"Recovery Manager",
			"Third Party Dependencies",
			"Outsourcing Partners",
			"Customer User",
		}
	} else {
		headers = []string{
			"SL No.",
			"Reference Number",
			"Number of People Necessary",
			"Business Owner Unit",
			"IT Responsible Unit",
			"Process Operation Team",
			"Process Recovery Team",
			"Recovery Manager",
			"Third Party Dependencies",
			"Outsourcing Partners",
			"Customer User",
		}
	}

	rows := make([][]interface{}, len(data))
	for i, sr := range data {
		numPeople := ""
		if sr.NumberOfPeopleNecessary != nil {
			numPeople = fmt.Sprintf("%d", *sr.NumberOfPeopleNecessary)
		}

		if userSetting.Language == "ar" {
			rows[i] = []interface{}{
				i + 1,
				sr.ReferenceNumber,
				numPeople,
				utility.GetStringValue(sr.BusinessOwnerUnitTitleAr),
				utility.GetStringValue(sr.ITResponsibleUnitTitleAr),
				utility.GetStringValue(sr.ProcessOperationTeamTitleAr),
				utility.GetStringValue(sr.ProcessRecoveryTeamTitleAr),
				fmt.Sprintf("%s %s", utility.GetStringValue(sr.RecoveryManagerFirstNameAr), utility.GetStringValue(sr.RecoveryManagerLastNameAr)),
				utility.GetStringValue(sr.ThirdPartyDependencies),
				utility.GetStringValue(sr.OutsourcingPartners),
				utility.GetStringValue(sr.CustomerUser),
			}
		} else {
			rows[i] = []interface{}{
				i + 1,
				sr.ReferenceNumber,
				numPeople,
				utility.GetStringValue(sr.BusinessOwnerUnitTitleEn),
				utility.GetStringValue(sr.ITResponsibleUnitTitleEn),
				utility.GetStringValue(sr.ProcessOperationTeamTitleEn),
				utility.GetStringValue(sr.ProcessRecoveryTeamTitleEn),
				fmt.Sprintf("%s %s", utility.GetStringValue(sr.RecoveryManagerFirstNameEn), utility.GetStringValue(sr.RecoveryManagerLastNameEn)),
				utility.GetStringValue(sr.ThirdPartyDependencies),
				utility.GetStringValue(sr.OutsourcingPartners),
				utility.GetStringValue(sr.CustomerUser),
			}
		}
	}

	return utility.GenerateExcelFile(headers, rows)
}

func (srs *SystemRegisterService) validateSystemRegisterFormDTO(systemRegisterFormDTO *dtos.SystemRegisterFormDTO) error {
	// Validate recovery manager
	if systemRegisterFormDTO.RecoveryManagerID != nil {
		if _, err := srs.userSvc.GetByID(*systemRegisterFormDTO.RecoveryManagerID); err != nil {
			return utility.NewValidationError(fmt.Errorf("recovery manager does not exist"), "not_exist", "recovery manager")
		}
	}

	// Validate job title IDs (designations)
	for _, jobTitleID := range systemRegisterFormDTO.JobTitleIDs {
		if _, err := srs.designationSvc.GetByID(jobTitleID); err != nil {
			return utility.NewValidationError(fmt.Errorf("designation with ID %d does not exist", jobTitleID), "not_exist", "designation")
		}
	}

	// Validate operational location IDs
	for _, locationID := range systemRegisterFormDTO.OperationalLocationIDs {
		if _, err := srs.locationSvc.GetByID(locationID); err != nil {
			return utility.NewValidationError(fmt.Errorf("operational location with ID %d does not exist", locationID), "not_exist", "operational location")
		}
	}

	// Validate recovery location IDs
	for _, locationID := range systemRegisterFormDTO.RecoveryLocationIDs {
		if _, err := srs.locationSvc.GetByID(locationID); err != nil {
			return utility.NewValidationError(fmt.Errorf("recovery location with ID %d does not exist", locationID), "not_exist", "recovery location")
		}
	}

	// Validate recovery needed tool IDs
	for _, toolID := range systemRegisterFormDTO.RecoveryNeededToolIDs {
		if _, err := srs.recoveryToolSvc.GetByID(toolID); err != nil {
			return utility.NewValidationError(fmt.Errorf("recovery tool with ID %d does not exist", toolID), "not_exist", "recovery tool")
		}
	}

	return nil
}
