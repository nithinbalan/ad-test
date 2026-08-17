package services

import (
	"errors"
	"fmt"

	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	socDtos "jethur/grc/security-operations/dtos"
	"jethur/grc/security-operations/repositories"
	settingServices "jethur/grc/settings/services"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type CSLogSourceDirectoryServiceInterface interface {
	GetAll(params *socDtos.CSLogSourceDirectoryListQueryParamsDTO) (*socDtos.CSLogSourceDirectoryListResponseDTO, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*socDtos.CSLogSourceDirectoryDetailDTO, error)
	GetByID(id uint) (*models.CSLogSourceDirectory, error)
	Create(currentUserID uint, form *socDtos.CSLogSourceDirectoryFormDTO) (uint, error)
	Update(id uint, currentUserID uint, form *socDtos.CSLogSourceDirectoryFormDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)
}

type CSLogSourceDirectoryService struct {
	repo                    repositories.CSLogSourceDirectoryRepositoryInterface
	activityLogSvc          logServices.ActivityLogServiceInterface
	dataEntityTypeSvc       CSLogDataEntityTypeServiceInterface
	receptionMechanismSvc   CSLogDataReceptionMechanismServiceInterface
	communicationChannelSvc CSLogCommunicationChannelServiceInterface
	organizationUnitSvc     settingServices.OrganizationUnitServiceInterface
}

func NewCSLogSourceDirectoryService(
	repo repositories.CSLogSourceDirectoryRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	dataEntityTypeSvc CSLogDataEntityTypeServiceInterface,
	receptionMechanismSvc CSLogDataReceptionMechanismServiceInterface,
	communicationChannelSvc CSLogCommunicationChannelServiceInterface,
	organizationUnitSvc settingServices.OrganizationUnitServiceInterface,
) CSLogSourceDirectoryServiceInterface {
	return &CSLogSourceDirectoryService{
		repo:                    repo,
		activityLogSvc:          activityLogSvc,
		dataEntityTypeSvc:       dataEntityTypeSvc,
		receptionMechanismSvc:   receptionMechanismSvc,
		communicationChannelSvc: communicationChannelSvc,
		organizationUnitSvc:     organizationUnitSvc,
	}
}

func (csds *CSLogSourceDirectoryService) GetAll(params *socDtos.CSLogSourceDirectoryListQueryParamsDTO) (*socDtos.CSLogSourceDirectoryListResponseDTO, error) {
	return csds.repo.GetAll(params)
}

func (csds *CSLogSourceDirectoryService) GetByID(id uint) (*models.CSLogSourceDirectory, error) {
	sourceDirectory, err := csds.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return sourceDirectory, nil
}

func (csds *CSLogSourceDirectoryService) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*socDtos.CSLogSourceDirectoryDetailDTO, error) {
	var sourceDirectoryDTO socDtos.CSLogSourceDirectoryDetailDTO

	sourceDirectory, err := csds.repo.GetFullDetailsByID(id, tx...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	err = mapstructure.Decode(sourceDirectory, &sourceDirectoryDTO)
	if err != nil {
		return nil, err
	}

	sourceDirectoryDTO.CreatedAt = sourceDirectory.CreatedAt
	sourceDirectoryDTO.UpdatedAt = sourceDirectory.UpdatedAt

	return &sourceDirectoryDTO, nil
}

func (csds *CSLogSourceDirectoryService) Create(currentUserID uint, sourceDirectoryFormDTO *socDtos.CSLogSourceDirectoryFormDTO) (uint, error) {
	tx := csds.repo.BeginTransaction()
	if err := csds.validateCSLogSourceDirectoryFormDTO(sourceDirectoryFormDTO); err != nil {
		tx.Rollback()
		return 0, err
	}

	referenceNumber, err := csds.repo.GenerateReferenceNumber()
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	sourceDirectory := &models.CSLogSourceDirectory{
		ReferenceNumber:     referenceNumber,
		EntityName:          sourceDirectoryFormDTO.EntityName,
		DataEntityTypeID:    sourceDirectoryFormDTO.DataEntityTypeID,
		Location:            sourceDirectoryFormDTO.Location,
		ResponsibleEntityID: sourceDirectoryFormDTO.ResponsibleEntityID,
		BeneficiaryEntityID: sourceDirectoryFormDTO.BeneficiaryEntityID,
		ContactName:         sourceDirectoryFormDTO.ContactName,
		ContactMobile:       sourceDirectoryFormDTO.ContactMobile,
		ContactEmail:        sourceDirectoryFormDTO.ContactEmail,
		Website:             sourceDirectoryFormDTO.Website,
		Notes:               sourceDirectoryFormDTO.Notes,
		CreatedBYID:         currentUserID,
	}

	id, err := csds.repo.Create(sourceDirectory, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	var receptionMechanisms []*models.CSLogSourceDirectoryReceptionMechanism
	for _, receptionMechanismID := range sourceDirectoryFormDTO.ReceptionMechanismIDs {
		receptionMechanism := &models.CSLogSourceDirectoryReceptionMechanism{
			SourceDirectoryID:    id,
			ReceptionMechanismID: *receptionMechanismID,
		}
		receptionMechanisms = append(receptionMechanisms, receptionMechanism)
	}
	if len(receptionMechanisms) > 0 {
		if err := csds.repo.CreateReceptionMechanisms(receptionMechanisms, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	var communicationChannels []*models.CSLogSourceDirectoryCommunicationChannel
	for _, communicationChannelID := range sourceDirectoryFormDTO.CommunicationChannelIDs {
		communicationChannel := &models.CSLogSourceDirectoryCommunicationChannel{
			SourceDirectoryID:      id,
			CommunicationChannelID: *communicationChannelID,
		}
		communicationChannels = append(communicationChannels, communicationChannel)
	}
	if len(communicationChannels) > 0 {
		if err := csds.repo.CreateCommunicationChannels(communicationChannels, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	if err := tx.Commit(); err != nil {
		return 0, utility.NewInternalServerError(err)
	}

	go func() {
		existingSourceDirectory, err := csds.GetFullDetailsByID(id)
		if err != nil {
			utility.LogError("Failed to fetch source directory", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		err = csds.activityLogSvc.LogActivity(
			models.CreateAction,
			models.CSLogSourceDirectory{}.TableName(),
			id,
			existingSourceDirectory,
			currentUserID,
			"security_operations", "security_operations_source_directory", "create source directory", nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.CSLogSourceDirectory{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return id, nil
}

func (csds *CSLogSourceDirectoryService) Update(id uint, currentUserID uint, sourceDirectoryFormDTO *socDtos.CSLogSourceDirectoryFormDTO) error {
	tx := csds.repo.BeginTransaction()
	existingSourceDirectory, err := csds.repo.GetByID(id)
	if err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utility.NewResourceNotFoundError(err, id)
		}
		return err
	}
	oldSourceDirectory, err := csds.GetFullDetailsByID(id, tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	if err := csds.validateCSLogSourceDirectoryFormDTO(sourceDirectoryFormDTO); err != nil {
		tx.Rollback()
		return err
	}

	existingSourceDirectory.EntityName = sourceDirectoryFormDTO.EntityName
	existingSourceDirectory.DataEntityTypeID = sourceDirectoryFormDTO.DataEntityTypeID
	existingSourceDirectory.Location = sourceDirectoryFormDTO.Location
	existingSourceDirectory.ResponsibleEntityID = sourceDirectoryFormDTO.ResponsibleEntityID
	existingSourceDirectory.BeneficiaryEntityID = sourceDirectoryFormDTO.BeneficiaryEntityID
	existingSourceDirectory.ContactName = sourceDirectoryFormDTO.ContactName
	existingSourceDirectory.ContactMobile = sourceDirectoryFormDTO.ContactMobile
	existingSourceDirectory.ContactEmail = sourceDirectoryFormDTO.ContactEmail
	existingSourceDirectory.Website = sourceDirectoryFormDTO.Website
	existingSourceDirectory.Notes = sourceDirectoryFormDTO.Notes
	existingSourceDirectory.UpdatedBYID = &currentUserID

	if err := csds.repo.Update(existingSourceDirectory, tx); err != nil {
		tx.Rollback()
		return err
	}

	existingReceptionMechanismIDs := make(map[uint]bool)
	for _, receptionMechanism := range oldSourceDirectory.ReceptionMechanisms {
		existingReceptionMechanismIDs[receptionMechanism.ReceptionMechanismID] = true
	}
	newReceptionMechanismIDs := make(map[uint]bool)
	for _, receptionMechanismID := range sourceDirectoryFormDTO.ReceptionMechanismIDs {
		newReceptionMechanismIDs[*receptionMechanismID] = true
	}

	var receptionMechanismsToAdd []*models.CSLogSourceDirectoryReceptionMechanism
	for _, receptionMechanismID := range sourceDirectoryFormDTO.ReceptionMechanismIDs {
		if _, exists := existingReceptionMechanismIDs[*receptionMechanismID]; !exists {
			receptionMechanism := &models.CSLogSourceDirectoryReceptionMechanism{
				SourceDirectoryID:    id,
				ReceptionMechanismID: *receptionMechanismID,
			}
			receptionMechanismsToAdd = append(receptionMechanismsToAdd, receptionMechanism)
		}
	}
	if len(receptionMechanismsToAdd) > 0 {
		if err := csds.repo.CreateReceptionMechanisms(receptionMechanismsToAdd, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	var receptionMechanismIDsToDelete []*uint
	for _, receptionMechanism := range oldSourceDirectory.ReceptionMechanisms {
		if _, exists := newReceptionMechanismIDs[receptionMechanism.ReceptionMechanismID]; !exists {
			receptionMechanismIDsToDelete = append(receptionMechanismIDsToDelete, receptionMechanism.ID)
		}
	}
	if len(receptionMechanismIDsToDelete) > 0 {
		if err := csds.repo.DeleteReceptionMechanisms(receptionMechanismIDsToDelete, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	existingCommunicationChannelIDs := make(map[uint]bool)
	for _, communicationChannel := range oldSourceDirectory.CommunicationChannels {
		existingCommunicationChannelIDs[communicationChannel.CommunicationChannelID] = true
	}
	newCommunicationChannelIDs := make(map[uint]bool)
	for _, communicationChannelID := range sourceDirectoryFormDTO.CommunicationChannelIDs {
		newCommunicationChannelIDs[*communicationChannelID] = true
	}

	var communicationChannelsToAdd []*models.CSLogSourceDirectoryCommunicationChannel
	for _, communicationChannelID := range sourceDirectoryFormDTO.CommunicationChannelIDs {
		if _, exists := existingCommunicationChannelIDs[*communicationChannelID]; !exists {
			communicationChannel := &models.CSLogSourceDirectoryCommunicationChannel{
				SourceDirectoryID:      id,
				CommunicationChannelID: *communicationChannelID,
			}
			communicationChannelsToAdd = append(communicationChannelsToAdd, communicationChannel)
		}
	}
	if len(communicationChannelsToAdd) > 0 {
		if err := csds.repo.CreateCommunicationChannels(communicationChannelsToAdd, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	var communicationChannelIDsToDelete []*uint
	for _, communicationChannel := range oldSourceDirectory.CommunicationChannels {
		if _, exists := newCommunicationChannelIDs[communicationChannel.CommunicationChannelID]; !exists {
			communicationChannelIDsToDelete = append(communicationChannelIDsToDelete, communicationChannel.ID)
		}
	}
	if len(communicationChannelIDsToDelete) > 0 {
		if err := csds.repo.DeleteCommunicationChannels(communicationChannelIDsToDelete, tx); err != nil {
			tx.Rollback()
			return err
		}
	}
	if err := tx.Commit(); err != nil {
		return utility.NewInternalServerError(err)
	}

	go func() {
		existingSourceDirectory, err := csds.GetFullDetailsByID(id)
		if err != nil {
			utility.LogError("Failed to fetch source directory", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		receptionMechanismsAdded := []uint{}
		receptionMechanismsRemoved := []uint{}
		for receptionMechanismID := range newReceptionMechanismIDs {
			if _, exists := existingReceptionMechanismIDs[receptionMechanismID]; !exists {
				receptionMechanismsAdded = append(receptionMechanismsAdded, receptionMechanismID)
			}
		}
		for receptionMechanismID := range existingReceptionMechanismIDs {
			if _, exists := newReceptionMechanismIDs[receptionMechanismID]; !exists {
				receptionMechanismsRemoved = append(receptionMechanismsRemoved, receptionMechanismID)
			}
		}

		communicationChannelsAdded := []uint{}
		communicationChannelsRemoved := []uint{}
		for communicationChannelID := range newCommunicationChannelIDs {
			if _, exists := existingCommunicationChannelIDs[communicationChannelID]; !exists {
				communicationChannelsAdded = append(communicationChannelsAdded, communicationChannelID)
			}
		}
		for communicationChannelID := range existingCommunicationChannelIDs {
			if _, exists := newCommunicationChannelIDs[communicationChannelID]; !exists {
				communicationChannelsRemoved = append(communicationChannelsRemoved, communicationChannelID)
			}
		}

		changedData := make(map[string]interface{})
		if len(receptionMechanismsAdded) > 0 || len(receptionMechanismsRemoved) > 0 {
			changedData["reception_mechanisms"] = map[string]interface{}{
				"added":   receptionMechanismsAdded,
				"removed": receptionMechanismsRemoved,
			}
		}
		if len(communicationChannelsAdded) > 0 || len(communicationChannelsRemoved) > 0 {
			changedData["communication_channels"] = map[string]interface{}{
				"added":   communicationChannelsAdded,
				"removed": communicationChannelsRemoved,
			}
		}

		oldData := map[string]interface{}{
			"entity_name":        oldSourceDirectory.EntityName,
			"data_entity_type":   utility.GetFieldValue(oldSourceDirectory.DataEntityType),
			"location":           oldSourceDirectory.Location,
			"responsible_entity": utility.GetFieldValue(oldSourceDirectory.ResponsibleEntity),
			"beneficiary_entity": utility.GetFieldValue(oldSourceDirectory.BeneficiaryEntity),
			"contact_name":       oldSourceDirectory.ContactName,
			"contact_mobile":     oldSourceDirectory.ContactMobile,
			"contact_email":      oldSourceDirectory.ContactEmail,
			"website":            oldSourceDirectory.Website,
			"notes":              oldSourceDirectory.Notes,
		}
		newData := map[string]interface{}{
			"entity_name":        existingSourceDirectory.EntityName,
			"data_entity_type":   utility.GetFieldValue(existingSourceDirectory.DataEntityType),
			"location":           existingSourceDirectory.Location,
			"responsible_entity": utility.GetFieldValue(existingSourceDirectory.ResponsibleEntity),
			"beneficiary_entity": utility.GetFieldValue(existingSourceDirectory.BeneficiaryEntity),
			"contact_name":       existingSourceDirectory.ContactName,
			"contact_mobile":     existingSourceDirectory.ContactMobile,
			"contact_email":      existingSourceDirectory.ContactEmail,
			"website":            existingSourceDirectory.Website,
			"notes":              existingSourceDirectory.Notes,
		}
		fieldChanges := utility.CompareFieldChanges(oldData, newData)
		for k, v := range fieldChanges {
			changedData[k] = v
		}

		if len(changedData) > 0 {
			err := csds.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.CSLogSourceDirectory{}.TableName(),
				id,
				changedData,
				currentUserID,
				"security_operations", "security_operations_source_directory", "update source directory", nil,
			)
			if err != nil {
				fields := map[string]interface{}{
					"action":    models.UpdateAction,
					"tableName": models.CSLogSourceDirectory{}.TableName(),
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

func (csds *CSLogSourceDirectoryService) validateCSLogSourceDirectoryFormDTO(sourceDirectoryFormDTO *socDtos.CSLogSourceDirectoryFormDTO) error {
	if sourceDirectoryFormDTO.DataEntityTypeID != nil {
		if _, err := csds.dataEntityTypeSvc.GetByID(*sourceDirectoryFormDTO.DataEntityTypeID); err != nil {
			return utility.NewValidationError(fmt.Errorf("data entity type does not exist"), "not_exist", "data entity type")
		}
	}
	if sourceDirectoryFormDTO.ResponsibleEntityID != nil {
		if _, err := csds.organizationUnitSvc.GetByID(*sourceDirectoryFormDTO.ResponsibleEntityID); err != nil {
			return utility.NewValidationError(fmt.Errorf("responsible entity does not exist"), "not_exist", "responsible entity")
		}
	}
	if sourceDirectoryFormDTO.BeneficiaryEntityID != nil {
		if _, err := csds.organizationUnitSvc.GetByID(*sourceDirectoryFormDTO.BeneficiaryEntityID); err != nil {
			return utility.NewValidationError(fmt.Errorf("beneficiary entity does not exist"), "not_exist", "beneficiary entity")
		}
	}

	receptionMechanismIDs := make(map[uint]bool)
	for _, receptionMechanismID := range sourceDirectoryFormDTO.ReceptionMechanismIDs {
		if receptionMechanismIDs[*receptionMechanismID] {
			return utility.NewValidationError(fmt.Errorf("duplicate reception mechanism ID: %d", receptionMechanismID), "duplicate", "reception mechanism")
		}
		if _, err := csds.receptionMechanismSvc.GetByID(*receptionMechanismID); err != nil {
			return utility.NewValidationError(fmt.Errorf("reception mechanism does not exist for ID: %d", receptionMechanismID), "not_exist", "reception mechanism")
		}
		receptionMechanismIDs[*receptionMechanismID] = true
	}

	communicationChannelIDs := make(map[uint]bool)
	for _, communicationChannelID := range sourceDirectoryFormDTO.CommunicationChannelIDs {
		if communicationChannelIDs[*communicationChannelID] {
			return utility.NewValidationError(fmt.Errorf("duplicate communication channel ID: %d", communicationChannelID), "duplicate", "communication channel")
		}
		if _, err := csds.communicationChannelSvc.GetByID(*communicationChannelID); err != nil {
			return utility.NewValidationError(fmt.Errorf("communication channel does not exist for ID: %d", communicationChannelID), "not_exist", "communication channel")
		}
		communicationChannelIDs[*communicationChannelID] = true
	}

	return nil
}

func (csds *CSLogSourceDirectoryService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIds []uint
	var successfulIds []uint
	var deletedSourceDirectories []map[string]interface{}

	for _, id := range ids {
		tx := csds.repo.BeginTransaction()
		sourceDirectory, err := csds.GetFullDetailsByID(id, tx)
		if err != nil {
			tx.Rollback()
			failedIds = append(failedIds, id)
			continue
		}

		if err := csds.repo.Delete(id, tx); err != nil {
			tx.Rollback()
			failedIds = append(failedIds, id)
			continue
		}
		if err := tx.Commit(); err != nil {
			failedIds = append(failedIds, id)
			continue
		}

		successfulIds = append(successfulIds, id)
		deletedSourceDirectories = append(deletedSourceDirectories, map[string]interface{}{
			"id":               id,
			"reference_number": sourceDirectory.ReferenceNumber,
			"entity_name":      sourceDirectory.EntityName,
		})
	}

	if len(successfulIds) > 0 {
		go func() {
			err := csds.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.CSLogSourceDirectory{}.TableName(),
				successfulIds,
				deletedSourceDirectories,
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
