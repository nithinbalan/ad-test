package services

import (
	"errors"
	"fmt"
	"log"
	"time"

	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	socDtos "jethur/grc/security-operations/dtos"
	"jethur/grc/security-operations/repositories"
	"jethur/utility"

	"gorm.io/gorm"
)

type CSLogSourceDirectoryServiceInterface interface {
	GetAll(params *socDtos.CSLogSourceDirectoryListQueryParamsDTO) (*socDtos.CSLogSourceDirectoryListResponseDTO, error)
	GetByID(id uint) (*models.CSLogSourceDirectory, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*socDtos.CSLogSourceDirectoryDetailDTO, error)
	Create(currentUserID uint, form *socDtos.CSLogSourceDirectoryFormDTO) (uint, error)
	Update(id uint, currentUserID uint, form *socDtos.CSLogSourceDirectoryFormDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)
}

type CSLogSourceDirectoryService struct {
	repo           repositories.CSLogSourceDirectoryRepositoryInterface
	activityLogSvc logServices.ActivityLogServiceInterface
}

func NewCSLogSourceDirectoryService(repo repositories.CSLogSourceDirectoryRepositoryInterface, activityLogSvc logServices.ActivityLogServiceInterface) CSLogSourceDirectoryServiceInterface {
	return &CSLogSourceDirectoryService{repo: repo, activityLogSvc: activityLogSvc}
}

func (csds *CSLogSourceDirectoryService) GetAll(params *socDtos.CSLogSourceDirectoryListQueryParamsDTO) (*socDtos.CSLogSourceDirectoryListResponseDTO, error) {
	return csds.repo.GetAll(params)
}

func (csds *CSLogSourceDirectoryService) GetByID(id uint) (*models.CSLogSourceDirectory, error) {
	entry, err := csds.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return entry, nil
}

func (csds *CSLogSourceDirectoryService) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*socDtos.CSLogSourceDirectoryDetailDTO, error) {
	entry, err := csds.repo.GetFullDetailsByID(id, tx...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return toDetailDTO(entry), nil
}

func (csds *CSLogSourceDirectoryService) Create(currentUserID uint, form *socDtos.CSLogSourceDirectoryFormDTO) (uint, error) {
	if err := csds.validateForm(form); err != nil {
		return 0, err
	}
	referenceNumber, err := csds.repo.GenerateReferenceNumber()
	if err != nil {
		return 0, err
	}

	entry := &models.CSLogSourceDirectory{
		ReferenceNumber:     referenceNumber,
		EntityName:          form.EntityName,
		DataEntityTypeID:    form.DataEntityTypeID,
		Location:            form.Location,
		ResponsibleEntityID: form.ResponsibleEntityID,
		BeneficiaryEntityID: form.BeneficiaryEntityID,
		ContactName:         form.ContactName,
		ContactMobile:       form.ContactMobile,
		ContactEmail:        form.ContactEmail,
		Website:             form.Website,
		Notes:               form.Notes,
		CreatedBYID:         currentUserID,
		CreatedAt:           time.Now(),
	}

	tx := csds.repo.BeginTransaction()

	createdID, err := csds.repo.Create(entry, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	if err := csds.syncReceptionMechanisms(createdID, form.ReceptionMechanismIDs, tx); err != nil {
		tx.Rollback()
		return 0, err
	}
	if err := csds.syncCommunicationChannels(createdID, form.CommunicationChannelIDs, tx); err != nil {
		tx.Rollback()
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		return 0, utility.NewInternalServerError(err)
	}

	go func() {
		fullEntry, err := csds.repo.GetFullDetailsByID(createdID)
		if err != nil {
			log.Printf("Failed to fetch created source directory for activity: %v", err)
			return
		}
		if err := csds.activityLogSvc.LogActivity(
			models.CreateAction,
			models.CSLogSourceDirectory{}.TableName(),
			createdID,
			fullEntry,
			currentUserID,
			"security_operations", "cs_log_source_directory", "create source directory", nil,
		); err != nil {
			log.Printf("Failed to log activity: %v", err)
		}
	}()

	return createdID, nil
}

func (csds *CSLogSourceDirectoryService) Update(id uint, currentUserID uint, form *socDtos.CSLogSourceDirectoryFormDTO) error {
	_, err := csds.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utility.NewResourceNotFoundError(err, id)
		}
		return err
	}
	tx := csds.repo.BeginTransaction()
	oldEntry, err := csds.repo.GetFullDetailsByID(id, tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	if err := csds.validateForm(form); err != nil {
		tx.Rollback()
		return err
	}

	now := time.Now()
	updated := &models.CSLogSourceDirectory{
		ID:                  id,
		EntityName:          form.EntityName,
		DataEntityTypeID:    form.DataEntityTypeID,
		Location:            form.Location,
		ResponsibleEntityID: form.ResponsibleEntityID,
		BeneficiaryEntityID: form.BeneficiaryEntityID,
		ContactName:         form.ContactName,
		ContactMobile:       form.ContactMobile,
		ContactEmail:        form.ContactEmail,
		Website:             form.Website,
		Notes:               form.Notes,
		UpdatedBYID:         &currentUserID,
		UpdatedAt:           &now,
	}

	if err := csds.repo.Update(updated, tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := csds.repo.DeleteReceptionMechanisms(id, tx); err != nil {
		tx.Rollback()
		return err
	}
	if err := csds.syncReceptionMechanisms(id, form.ReceptionMechanismIDs, tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := csds.repo.DeleteCommunicationChannels(id, tx); err != nil {
		tx.Rollback()
		return err
	}
	if err := csds.syncCommunicationChannels(id, form.CommunicationChannelIDs, tx); err != nil {
		tx.Rollback()
		return err
	}
	updatedEntry, err := csds.repo.GetFullDetailsByID(id, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return utility.NewInternalServerError(err)
	}

	fieldsToCompare := []string{"EntityName", "DataEntityTypeID", "Location", "ResponsibleEntityID", "BeneficiaryEntityID", "ContactName", "ContactMobile", "ContactEmail", "Website", "Notes"}
	changes, err := utility.CompareFields(oldEntry, updatedEntry, fieldsToCompare)
	if err != nil {
		log.Println("error :", err.Error())
	}
	changedDataJSON, err := utility.ChangesToJSON(changes)
	if err != nil {
		log.Println("error :", err.Error())
	}

	go func() {
		if err := csds.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.CSLogSourceDirectory{}.TableName(),
			id,
			changedDataJSON,
			currentUserID,
			"security_operations", "cs_log_source_directory", "update source directory", nil,
		); err != nil {
			log.Printf("Failed to log update activity: %v", err)
		}
	}()

	return nil
}

func (csds *CSLogSourceDirectoryService) validateForm(form *socDtos.CSLogSourceDirectoryFormDTO) error {
	if form.DataEntityTypeID != nil {
		exists, err := csds.repo.ExistsDataEntityType(*form.DataEntityTypeID)
		if err != nil {
			return err
		}
		if !exists {
			return utility.NewValidationError(fmt.Errorf("data entity type does not exist"), "not_exist", "data_entity_type_id")
		}
	}
	for _, item := range []struct {
		id    *uint
		field string
	}{{form.ResponsibleEntityID, "responsible_entity_id"}, {form.BeneficiaryEntityID, "beneficiary_entity_id"}} {
		if item.id == nil {
			continue
		}
		exists, err := csds.repo.ExistsOrganizationUnit(*item.id)
		if err != nil {
			return err
		}
		if !exists {
			return utility.NewValidationError(fmt.Errorf("organization unit does not exist"), "not_exist", item.field)
		}
	}
	for _, id := range form.ReceptionMechanismIDs {
		exists, err := csds.repo.ExistsReceptionMechanism(id)
		if err != nil {
			return err
		}
		if !exists {
			return utility.NewValidationError(fmt.Errorf("reception mechanism does not exist"), "not_exist", "reception_mechanism_ids")
		}
	}
	for _, id := range form.CommunicationChannelIDs {
		exists, err := csds.repo.ExistsCommunicationChannel(id)
		if err != nil {
			return err
		}
		if !exists {
			return utility.NewValidationError(fmt.Errorf("communication channel does not exist"), "not_exist", "communication_channel_ids")
		}
	}
	return nil
}

func (csds *CSLogSourceDirectoryService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIds []uint
	var successfulIds []uint
	var deleted []map[string]any

	for _, id := range ids {
		tx := csds.repo.BeginTransaction()
		entry, err := csds.repo.GetFullDetailsByID(id, tx)
		if err != nil {
			tx.Rollback()
			failedIds = append(failedIds, id)
			continue
		}
		if err := csds.repo.DeleteReceptionMechanisms(id, tx); err != nil {
			tx.Rollback()
			failedIds = append(failedIds, id)
			continue
		}
		if err := csds.repo.DeleteCommunicationChannels(id, tx); err != nil {
			tx.Rollback()
			failedIds = append(failedIds, id)
			continue
		}
		if err := csds.repo.Delete(id, tx); err != nil {
			tx.Rollback()
			failedIds = append(failedIds, id)
		} else {
			if err := tx.Commit(); err != nil {
				failedIds = append(failedIds, id)
				continue
			}
			successfulIds = append(successfulIds, id)
		}
		deleted = append(deleted, map[string]any{
			"id":               id,
			"reference_number": entry.ReferenceNumber,
			"entity_name":      entry.EntityName,
		})
	}

	if len(successfulIds) > 0 {
		go func() {
			if err := csds.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.CSLogSourceDirectory{}.TableName(),
				successfulIds,
				deleted,
				currentUserID,
				nil,
			); err != nil {
				log.Printf("Failed to log bulk delete activity: %v", err)
			}
		}()
	}

	return failedIds, nil
}

func (csds *CSLogSourceDirectoryService) syncReceptionMechanisms(sourceDirectoryID uint, ids []uint, tx utility.DatabaseTransaction) error {
	links := make([]*models.CSLogSourceDirectoryReceptionMechanism, 0, len(ids))
	seen := map[uint]bool{}
	for _, mid := range ids {
		if seen[mid] {
			continue
		}
		seen[mid] = true
		links = append(links, &models.CSLogSourceDirectoryReceptionMechanism{
			SourceDirectoryID:    sourceDirectoryID,
			ReceptionMechanismID: mid,
		})
	}
	return csds.repo.CreateReceptionMechanisms(links, tx)
}

func (csds *CSLogSourceDirectoryService) syncCommunicationChannels(sourceDirectoryID uint, ids []uint, tx utility.DatabaseTransaction) error {
	links := make([]*models.CSLogSourceDirectoryCommunicationChannel, 0, len(ids))
	seen := map[uint]bool{}
	for _, cid := range ids {
		if seen[cid] {
			continue
		}
		seen[cid] = true
		links = append(links, &models.CSLogSourceDirectoryCommunicationChannel{
			SourceDirectoryID:      sourceDirectoryID,
			CommunicationChannelID: cid,
		})
	}
	return csds.repo.CreateCommunicationChannels(links, tx)
}

func toListDTO(entry *models.CSLogSourceDirectory) *socDtos.CSLogSourceDirectoryListDTO {
	dto := &socDtos.CSLogSourceDirectoryListDTO{
		ID:                  entry.ID,
		ReferenceNumber:     entry.ReferenceNumber,
		EntityName:          entry.EntityName,
		DataEntityTypeID:    entry.DataEntityTypeID,
		ResponsibleEntityID: entry.ResponsibleEntityID,
		BeneficiaryEntityID: entry.BeneficiaryEntityID,
		CreatedAt:           entry.CreatedAt,
		UpdatedAt:           entry.UpdatedAt,
	}
	if entry.DataEntityType != nil {
		dto.DataEntityTypeTitleEn = &entry.DataEntityType.TitleEn
		dto.DataEntityTypeTitleAr = &entry.DataEntityType.TitleAr
	}
	if entry.ResponsibleEntity != nil {
		dto.ResponsibleEntityTitleEn = &entry.ResponsibleEntity.TitleEn
		dto.ResponsibleEntityTitleAr = &entry.ResponsibleEntity.TitleAr
	}
	if entry.BeneficiaryEntity != nil {
		dto.BeneficiaryEntityTitleEn = &entry.BeneficiaryEntity.TitleEn
		dto.BeneficiaryEntityTitleAr = &entry.BeneficiaryEntity.TitleAr
	}
	return dto
}

func toDetailDTO(entry *models.CSLogSourceDirectory) *socDtos.CSLogSourceDirectoryDetailDTO {
	dto := &socDtos.CSLogSourceDirectoryDetailDTO{
		ID:                    entry.ID,
		ReferenceNumber:       entry.ReferenceNumber,
		EntityName:            entry.EntityName,
		DataEntityTypeID:      entry.DataEntityTypeID,
		Location:              entry.Location,
		ResponsibleEntityID:   entry.ResponsibleEntityID,
		BeneficiaryEntityID:   entry.BeneficiaryEntityID,
		ContactName:           entry.ContactName,
		ContactMobile:         entry.ContactMobile,
		ContactEmail:          entry.ContactEmail,
		Website:               entry.Website,
		Notes:                 entry.Notes,
		CreatedByID:           entry.CreatedBYID,
		UpdatedByID:           entry.UpdatedBYID,
		CreatedAt:             entry.CreatedAt,
		UpdatedAt:             entry.UpdatedAt,
		ReceptionMechanisms:   make([]*socDtos.CSLogSourceDirectoryReceptionMechanismDTO, 0, len(entry.ReceptionMechanisms)),
		CommunicationChannels: make([]*socDtos.CSLogSourceDirectoryCommunicationChannelDTO, 0, len(entry.CommunicationChannels)),
	}
	if entry.DataEntityType != nil {
		dto.DataEntityType = &socDtos.CSLogSourceDirectoryLookupDTO{ID: entry.DataEntityType.ID, TitleEn: entry.DataEntityType.TitleEn, TitleAr: entry.DataEntityType.TitleAr}
	}
	if entry.ResponsibleEntity != nil {
		dto.ResponsibleEntity = &socDtos.CSLogSourceDirectoryOrganizationUnitDTO{ID: entry.ResponsibleEntity.ID, TitleEn: entry.ResponsibleEntity.TitleEn, TitleAr: entry.ResponsibleEntity.TitleAr}
	}
	if entry.BeneficiaryEntity != nil {
		dto.BeneficiaryEntity = &socDtos.CSLogSourceDirectoryOrganizationUnitDTO{ID: entry.BeneficiaryEntity.ID, TitleEn: entry.BeneficiaryEntity.TitleEn, TitleAr: entry.BeneficiaryEntity.TitleAr}
	}
	if entry.CreatedByUser != nil {
		dto.CreatedBy = toUserDTO(entry.CreatedByUser)
	}
	if entry.UpdatedByUser != nil {
		dto.UpdatedBy = toUserDTO(entry.UpdatedByUser)
	}
	for _, relation := range entry.ReceptionMechanisms {
		item := &socDtos.CSLogSourceDirectoryReceptionMechanismDTO{ID: relation.ID, SourceDirectoryID: relation.SourceDirectoryID, ReceptionMechanismID: relation.ReceptionMechanismID}
		if relation.ReceptionMechanism != nil {
			item.ReceptionMechanism = &socDtos.CSLogSourceDirectoryLookupDTO{ID: relation.ReceptionMechanism.ID, TitleEn: relation.ReceptionMechanism.TitleEn, TitleAr: relation.ReceptionMechanism.TitleAr}
		}
		dto.ReceptionMechanisms = append(dto.ReceptionMechanisms, item)
	}
	for _, relation := range entry.CommunicationChannels {
		item := &socDtos.CSLogSourceDirectoryCommunicationChannelDTO{ID: relation.ID, SourceDirectoryID: relation.SourceDirectoryID, CommunicationChannelID: relation.CommunicationChannelID}
		if relation.CommunicationChannel != nil {
			item.CommunicationChannel = &socDtos.CSLogSourceDirectoryLookupDTO{ID: relation.CommunicationChannel.ID, TitleEn: relation.CommunicationChannel.TitleEn, TitleAr: relation.CommunicationChannel.TitleAr}
		}
		dto.CommunicationChannels = append(dto.CommunicationChannels, item)
	}
	return dto
}

func toUserDTO(user *models.User) *socDtos.CSLogSourceDirectoryUserDTO {
	return &socDtos.CSLogSourceDirectoryUserDTO{ID: user.ID, FirstNameEn: user.FirstNameEn, FirstNameAr: user.FirstNameAr, LastNameEn: user.LastNameEn, LastNameAr: user.LastNameAr, Email: user.Email}
}
