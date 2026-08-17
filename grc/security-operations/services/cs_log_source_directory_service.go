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
	GetAll(sort, sortBy, keyword string, dataEntityTypeID, responsibleEntityID, beneficiaryEntityID *uint) ([]*models.CSLogSourceDirectory, error)
	GetByID(id uint) (*models.CSLogSourceDirectory, error)
	Create(currentUserID uint, form *socDtos.CSLogSourceDirectoryFormDTO) (*models.CSLogSourceDirectory, error)
	Update(id uint, currentUserID uint, form *socDtos.CSLogSourceDirectoryFormDTO) (*models.CSLogSourceDirectory, error)
	Delete(ids []uint, currentUserID uint) ([]uint, error)
}

type CSLogSourceDirectoryService struct {
	repo           repositories.CSLogSourceDirectoryRepositoryInterface
	activityLogSvc logServices.ActivityLogServiceInterface
}

func NewCSLogSourceDirectoryService(repo repositories.CSLogSourceDirectoryRepositoryInterface, activityLogSvc logServices.ActivityLogServiceInterface) CSLogSourceDirectoryServiceInterface {
	return &CSLogSourceDirectoryService{repo: repo, activityLogSvc: activityLogSvc}
}

func (csds *CSLogSourceDirectoryService) GetAll(sort, sortBy, keyword string, dataEntityTypeID, responsibleEntityID, beneficiaryEntityID *uint) ([]*models.CSLogSourceDirectory, error) {
	return csds.repo.GetAll(sort, sortBy, keyword, dataEntityTypeID, responsibleEntityID, beneficiaryEntityID)
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

func (csds *CSLogSourceDirectoryService) Create(currentUserID uint, form *socDtos.CSLogSourceDirectoryFormDTO) (*models.CSLogSourceDirectory, error) {
	if exists, err := csds.repo.ExistsByEntityName(form.EntityName); err != nil {
		return nil, err
	} else if exists {
		return nil, utility.NewValidationError(fmt.Errorf("the entity name must be unique"), "unique", "entity_name")
	}

	// reference_number is set from the row id post insert so it never reuses on delete.
	entry := &models.CSLogSourceDirectory{
		ReferenceNumber:     fmt.Sprintf("PENDING-%d", time.Now().UnixNano()),
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

	created, err := csds.repo.Create(entry, tx)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	referenceNumber := fmt.Sprintf("SRC-%d-%04d", time.Now().Year(), created.ID)
	if err := csds.repo.SetReferenceNumber(created.ID, referenceNumber, tx); err != nil {
		tx.Rollback()
		return nil, err
	}
	created.ReferenceNumber = referenceNumber

	if err := csds.syncReceptionMechanisms(created.ID, form.ReceptionMechanismIDs, tx); err != nil {
		tx.Rollback()
		return nil, err
	}
	if err := csds.syncCommunicationChannels(created.ID, form.CommunicationChannelIDs, tx); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, utility.NewInternalServerError(err)
	}

	// Reload with all preloads so the response includes junction data
	fullEntry, err := csds.repo.GetByID(created.ID)
	if err != nil {
		return nil, err
	}

	go func() {
		if err := csds.activityLogSvc.LogActivity(
			models.CreateAction,
			models.CSLogSourceDirectory{}.TableName(),
			fullEntry.ID,
			fullEntry,
			currentUserID,
			"", "", "", nil,
		); err != nil {
			log.Printf("Failed to log activity: %v", err)
		}
	}()

	return fullEntry, nil
}

func (csds *CSLogSourceDirectoryService) Update(id uint, currentUserID uint, form *socDtos.CSLogSourceDirectoryFormDTO) (*models.CSLogSourceDirectory, error) {
	existing, err := csds.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	if form.EntityName != existing.EntityName {
		if exists, err := csds.repo.ExistsByEntityNameExcludingID(id, form.EntityName); err != nil {
			return nil, err
		} else if exists {
			return nil, utility.NewValidationError(fmt.Errorf("the entity name must be unique"), "unique", "entity_name")
		}
	}

	now := time.Now()
	updated := &models.CSLogSourceDirectory{
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

	tx := csds.repo.BeginTransaction()

	if _, err := csds.repo.Update(id, updated, tx); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := csds.repo.DeleteReceptionMechanisms(id, tx); err != nil {
		tx.Rollback()
		return nil, err
	}
	if err := csds.syncReceptionMechanisms(id, form.ReceptionMechanismIDs, tx); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := csds.repo.DeleteCommunicationChannels(id, tx); err != nil {
		tx.Rollback()
		return nil, err
	}
	if err := csds.syncCommunicationChannels(id, form.CommunicationChannelIDs, tx); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, utility.NewInternalServerError(err)
	}

	fullEntry, err := csds.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	fieldsToCompare := []string{"EntityName", "DataEntityTypeID", "Location", "ResponsibleEntityID", "BeneficiaryEntityID", "ContactName", "ContactMobile", "ContactEmail", "Website", "Notes"}
	changes, err := utility.CompareFields(existing, updated, fieldsToCompare)
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
			"", "", "", nil,
		); err != nil {
			log.Printf("Failed to log update activity: %v", err)
		}
	}()

	return fullEntry, nil
}

func (csds *CSLogSourceDirectoryService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	var failedIds []uint
	var successfulIds []uint
	var deleted []map[string]any

	for _, id := range ids {
		entry, err := csds.repo.GetByID(id)
		if err != nil {
			continue
		}
		if err := csds.repo.Delete(id); err != nil {
			failedIds = append(failedIds, id)
		} else {
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
