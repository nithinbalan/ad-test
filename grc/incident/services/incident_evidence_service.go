package services

import (
	"errors"
	"fmt"
	"jethur/config"
	"jethur/grc/incident/dtos"
	"jethur/grc/incident/repositories"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	"jethur/utility"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type IncidentEvidenceServiceInterface interface {
	GetAll(incidentRegisterID uint, keyword string) ([]*dtos.IncidentEvidenceListDTO, error)
	GetByID(incidentRegisterID, id uint) (*models.IncidentEvidence, error)
	GetFullDetailsByID(incidentRegisterID uint, id uint, tx ...utility.DatabaseTransaction) (*dtos.IncidentEvidenceDetailDTO, error)
	Create(incidentRegisterID uint, form *dtos.IncidentEvidenceFormDTO, currentUserID uint) ([]uint, error)
	Update(incidentRegisterID uint, id uint, form *dtos.IncidentEvidenceFormDTO, currentUserID uint) error
	Delete(incidentRegisterID uint, ids []uint, currentUserID uint) ([]uint, error)
	GetIncidentEvidencePath(incidentRegisterID, evidenceID uint, token string, currentUserID uint) (string, error)
	GetDownloadHistories(incidentRegisterID, evidenceID uint) ([]*dtos.IncidentEvidenceDownloadHistoryListDTO, error)
}

type IncidentEvidenceService struct {
	repo                repositories.IncidentEvidenceRepositoryInterface
	activityLogSvc      logServices.ActivityLogServiceInterface
	incidentRegisterSvc IncidentRegisterServiceInterface
	evidenceSourceSvc   IncidentEvidenceSourceServiceInterface
	evidenceTypeSvc     IncidentEvidenceTypeServiceInterface
}

func NewIncidentEvidenceService(
	repo repositories.IncidentEvidenceRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	incidentRegisterSvc IncidentRegisterServiceInterface,
	evidenceSourceSvc IncidentEvidenceSourceServiceInterface,
	evidenceTypeSvc IncidentEvidenceTypeServiceInterface,
) IncidentEvidenceServiceInterface {
	return &IncidentEvidenceService{
		repo:                repo,
		activityLogSvc:      activityLogSvc,
		incidentRegisterSvc: incidentRegisterSvc,
		evidenceSourceSvc:   evidenceSourceSvc,
		evidenceTypeSvc:     evidenceTypeSvc,
	}
}

func (ies *IncidentEvidenceService) GetAll(incidentRegisterID uint, keyword string) ([]*dtos.IncidentEvidenceListDTO, error) {
	if _, err := ies.incidentRegisterSvc.GetByID(incidentRegisterID); err != nil {
		return []*dtos.IncidentEvidenceListDTO{}, err
	}

	evidences, err := ies.repo.GetAll(incidentRegisterID, keyword)
	if err != nil {
		return []*dtos.IncidentEvidenceListDTO{}, err
	}
	if evidences == nil {
		evidences = []*dtos.IncidentEvidenceListDTO{}
	}

	groupMap := make(map[uint][]*dtos.IncidentEvidenceListDTO)
	for _, evidence := range evidences {
		groupMap[evidence.IncidentEvidenceSourceID] = append(groupMap[evidence.IncidentEvidenceSourceID], evidence)
	}

	var result []*dtos.IncidentEvidenceListDTO
	for _, group := range groupMap {
		parent := group[0]

		var cleanedEvidences []*dtos.IncidentEvidenceListDTO
		for _, evidence := range group {
			cleanedEvidences = append(cleanedEvidences, &dtos.IncidentEvidenceListDTO{
				ID:                        evidence.ID,
				IncidentRegisterID:        evidence.IncidentRegisterID,
				IncidentEvidenceSourceID:  evidence.IncidentEvidenceSourceID,
				IncidentEvidenceTypeID:    evidence.IncidentEvidenceTypeID,
				IncidentEvidenceTypeEn:    evidence.IncidentEvidenceTypeEn,
				IncidentEvidenceTypeAr:    evidence.IncidentEvidenceTypeAr,
				URL:                       evidence.URL,
				Name:                      evidence.Name,
				Size:                      evidence.Size,
				Token:                     evidence.Token,
				OriginalName:              evidence.OriginalName,
				Type:                      evidence.Type,
				IsReadOnly:                evidence.IsReadOnly,
				IsDownload:                evidence.IsDownload,
				CreatedAt:                 evidence.CreatedAt,
				CreatedByID:               evidence.CreatedByID,
				CreatedByFirstNameEn:      evidence.CreatedByFirstNameEn,
				CreatedByFirstNameAr:      evidence.CreatedByFirstNameAr,
				CreatedByLastNameEn:       evidence.CreatedByLastNameEn,
				CreatedByLastNameAr:       evidence.CreatedByLastNameAr,
				CreatedByInitial:          evidence.CreatedByInitial,
				CreatedByColor:            evidence.CreatedByColor,
				CreatedByEmail:            evidence.CreatedByEmail,
				CreatedByProfileImagePath: evidence.CreatedByProfileImagePath,
				UpdatedAt:                 evidence.UpdatedAt,
				UpdatedByID:               evidence.UpdatedByID,
				UpdatedByFirstNameEn:      evidence.UpdatedByFirstNameEn,
				UpdatedByFirstNameAr:      evidence.UpdatedByFirstNameAr,
				UpdatedByLastNameEn:       evidence.UpdatedByLastNameEn,
				UpdatedByLastNameAr:       evidence.UpdatedByLastNameAr,
				UpdatedByInitial:          evidence.UpdatedByInitial,
				UpdatedByColor:            evidence.UpdatedByColor,
				UpdatedByEmail:            evidence.UpdatedByEmail,
				UpdatedByProfileImagePath: evidence.UpdatedByProfileImagePath,
			})
		}

		parentCopy := *parent
		parentCopy.Documents = cleanedEvidences
		result = append(result, &parentCopy)
	}

	if result == nil {
		result = []*dtos.IncidentEvidenceListDTO{}
	}

	return result, nil
}

func (ies *IncidentEvidenceService) GetByID(incidentRegisterID uint, id uint) (*models.IncidentEvidence, error) {
	if _, err := ies.incidentRegisterSvc.GetByID(incidentRegisterID); err != nil {
		return nil, err
	}
	evidence, err := ies.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return evidence, nil
}

func (ies *IncidentEvidenceService) GetFullDetailsByID(incidentRegisterID uint, id uint, tx ...utility.DatabaseTransaction) (*dtos.IncidentEvidenceDetailDTO, error) {
	if _, err := ies.incidentRegisterSvc.GetByID(incidentRegisterID); err != nil {
		return nil, err
	}
	evidence, err := ies.repo.GetFullDetailsByID(id, tx...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	var dto dtos.IncidentEvidenceDetailDTO
	err = mapstructure.Decode(evidence, &dto)
	if err != nil {
		return nil, err
	}
	return &dto, nil
}

func (ies *IncidentEvidenceService) Create(incidentRegisterID uint, form *dtos.IncidentEvidenceFormDTO, currentUserID uint) ([]uint, error) {
	tx := ies.repo.BeginTransaction()

	if _, err := ies.incidentRegisterSvc.GetByID(incidentRegisterID); err != nil {
		tx.Rollback()
		return nil, err
	}

	if _, err := ies.evidenceSourceSvc.GetByID(form.IncidentEvidenceSourceID); err != nil {
		tx.Rollback()
		return nil, utility.NewValidationError(fmt.Errorf("invalid incident evidence source id"), "not_exist", "incident_evidence_source_id")
	}

	if _, err := ies.evidenceTypeSvc.GetByID(form.IncidentEvidenceTypeID); err != nil {
		tx.Rollback()
		return nil, utility.NewValidationError(fmt.Errorf("invalid incident evidence type id"), "not_exist", "incident_evidence_type_id")
	}

	var createdIDs []uint
	var evidences []*models.IncidentEvidence

	for _, doc := range form.Documents {
		encryptedToken, err := utility.EncryptData([]byte(doc.FileName))
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		url := utility.GenerateDocumentURL("incident", "incident-register", strconv.FormatUint(uint64(incidentRegisterID), 10), "evidences")

		evidence := &models.IncidentEvidence{
			IncidentRegisterID:       incidentRegisterID,
			IncidentEvidenceSourceID: form.IncidentEvidenceSourceID,
			IncidentEvidenceTypeID:   form.IncidentEvidenceTypeID,
			URL:                      &url,
			Name:                     &doc.FileName,
			Size:                     &doc.Size,
			Token:                    &encryptedToken,
			OriginalName:             &doc.OriginalName,
			Type:                     &doc.Type,
			IsReadOnly:               form.IsReadOnly,
			IsDownload:               form.IsDownload,
			CreatedBYID:              currentUserID,
		}

		evidences = append(evidences, evidence)
	}
	if len(evidences) > 0 {
		if err := ies.repo.CreateIncidentEvidences(evidences, tx); err != nil {
			tx.Rollback()
			return nil, err
		}
		for _, evidence := range evidences {
			storagePath := filepath.Join("incident", "incident-register", strconv.FormatUint(uint64(incidentRegisterID), 10), "evidences")
			if err := utility.StoreFile(*evidence.Name, storagePath); err != nil {
				tx.Rollback()
				return nil, err
			}
			createdIDs = append(createdIDs, evidence.ID)
		}
	}

	tx.Commit()

	return createdIDs, nil
}

func (ies *IncidentEvidenceService) Update(incidentRegisterID uint, id uint, form *dtos.IncidentEvidenceFormDTO, currentUserID uint) error {
	if _, err := ies.incidentRegisterSvc.GetByID(incidentRegisterID); err != nil {
		return err
	}

	if _, err := ies.evidenceSourceSvc.GetByID(form.IncidentEvidenceSourceID); err != nil {
		return utility.NewValidationError(fmt.Errorf("invalid incident evidence source id"), "not_exist", "incident_evidence_source_id")
	}

	if _, err := ies.evidenceTypeSvc.GetByID(form.IncidentEvidenceTypeID); err != nil {
		return utility.NewValidationError(fmt.Errorf("invalid incident evidence type id"), "not_exist", "incident_evidence_type_id")
	}

	tx := ies.repo.BeginTransaction()

	var newEvidences []*models.IncidentEvidence
	var deleteEvidenceIds []*uint
	var evidencesToDelete []*models.IncidentEvidence

	for _, doc := range form.Documents {
		if doc.IsNewFile != nil && *doc.IsNewFile {

			encryptedFileName, err := utility.EncryptData([]byte(doc.FileName))
			if err != nil {
				tx.Rollback()
				return err
			}

			url := utility.GenerateDocumentURL("incident", "incident-register", strconv.FormatUint(uint64(incidentRegisterID), 10), "evidences")

			evidence := &models.IncidentEvidence{
				IncidentRegisterID:       incidentRegisterID,
				IncidentEvidenceSourceID: form.IncidentEvidenceSourceID,
				IncidentEvidenceTypeID:   form.IncidentEvidenceTypeID,
				URL:                      &url,
				Name:                     &doc.FileName,
				Size:                     &doc.Size,
				Token:                    &encryptedFileName,
				OriginalName:             &doc.OriginalName,
				Type:                     &doc.Type,
				IsReadOnly:               form.IsReadOnly,
				IsDownload:               form.IsDownload,
				CreatedBYID:              currentUserID,
			}
			newEvidences = append(newEvidences, evidence)
		} else if doc.IsDeletedFile != nil && *doc.IsDeletedFile {

			if doc.ID != nil {
				existingEvidence, err := ies.repo.GetByID(*doc.ID)
				if err != nil {
					tx.Rollback()
					return err
				}
				evidencesToDelete = append(evidencesToDelete, existingEvidence)
				deleteEvidenceIds = append(deleteEvidenceIds, &existingEvidence.ID)
			}
		} else {
			if doc.ID != nil {
				existingEvidence, err := ies.repo.GetByID(*doc.ID)
				if err != nil {
					tx.Rollback()
					return err
				}

				if existingEvidence.IncidentEvidenceSourceID != form.IncidentEvidenceSourceID ||
					existingEvidence.IncidentEvidenceTypeID != form.IncidentEvidenceTypeID ||
					existingEvidence.IsReadOnly != form.IsReadOnly ||
					existingEvidence.IsDownload != form.IsDownload {
					existingEvidence.IncidentEvidenceSourceID = form.IncidentEvidenceSourceID
					existingEvidence.IncidentEvidenceTypeID = form.IncidentEvidenceTypeID
					existingEvidence.IsReadOnly = form.IsReadOnly
					existingEvidence.IsDownload = form.IsDownload
					existingEvidence.UpdatedBYID = &currentUserID

					if err := ies.repo.Update(existingEvidence); err != nil {
						tx.Rollback()
						return err
					}
				}
			}
		}
	}

	if len(newEvidences) > 0 {
		if err := ies.repo.CreateIncidentEvidences(newEvidences, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	if len(deleteEvidenceIds) > 0 {
		if err := ies.repo.DeleteIncidentEvidences(deleteEvidenceIds, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, evidence := range newEvidences {
		storagePath := filepath.Join("incident", "incident-register", strconv.FormatUint(uint64(incidentRegisterID), 10), "evidences")
		if err := utility.StoreFile(*evidence.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, evidence := range evidencesToDelete {
		storagePath := filepath.Join("incident", "incident-register", strconv.FormatUint(uint64(incidentRegisterID), 10), "evidences")
		if err := utility.RemoveFile(*evidence.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()
	return nil
}

func (ies *IncidentEvidenceService) Delete(incidentRegisterID uint, ids []uint, currentUserID uint) ([]uint, error) {
	var failedIds, successIds []uint
	var deletedItems []map[string]interface{}
	evidencesToDelete := make([]*models.IncidentEvidence, 0)

	tx := ies.repo.BeginTransaction()

	for _, id := range ids {
		evidence, err := ies.repo.GetByID(id)
		if err != nil {
			failedIds = append(failedIds, id)
			continue
		}

		evidencesToDelete = append(evidencesToDelete, evidence)

		if err := ies.repo.Delete(id, tx); err != nil {
			tx.Rollback()
			failedIds = append(failedIds, id)
			continue
		}

		successIds = append(successIds, id)
		deletedItems = append(deletedItems, map[string]interface{}{
			"id":   evidence.ID,
			"name": evidence.Name,
			"type": evidence.Type,
		})
	}

	if len(successIds) == 0 {
		tx.Rollback()
		return failedIds, nil
	}

	tx.Commit()

	for _, evidence := range evidencesToDelete {
		storagePath := filepath.Join("incident", "incident-register", strconv.FormatUint(uint64(incidentRegisterID), 10), "evidences")
		if err := utility.RemoveFile(*evidence.Name, storagePath); err != nil {
			fmt.Printf("⚠️ Failed to remove file %s: %v\n", *evidence.Name, err)
		}
	}

	if len(evidencesToDelete) > 0 {
		folderPath := filepath.Join("incident", "incident-register", strconv.FormatUint(uint64(incidentRegisterID), 10), "evidences")
		if err := utility.RemoveFolder(folderPath); err != nil {
			fmt.Printf("⚠️ Failed to remove folder %s: %v\n", folderPath, err)
		}
	}

	go func() {
		_ = ies.activityLogSvc.LogBulkActivities(
			models.DeleteAction,
			models.IncidentEvidence{}.TableName(),
			successIds,
			deletedItems,
			currentUserID,
			nil,
		)
	}()

	return failedIds, nil
}

func (ies *IncidentEvidenceService) GetIncidentEvidencePath(incidentRegisterID, evidenceID uint, token string, currentUserID uint) (string, error) {
	evidence, err := ies.repo.GetByID(evidenceID)
	if err != nil {
		return "", err
	}

	if evidence.Name == nil {
		return "", fmt.Errorf("evidence filename not found")
	}

	decryptedFileNameBytes, err := utility.DecryptData(token)
	if err != nil {
		return "", fmt.Errorf("resource not found")
	}

	decryptedFileName := string(decryptedFileNameBytes)
	if decryptedFileName == "" {
		return "", fmt.Errorf("resource not found")
	}

	storagePath := filepath.Join("incident", "incident-register", strconv.FormatUint(uint64(incidentRegisterID), 10), "evidences")
	baseStoragePath := config.GetEnv("STORAGE_PATH", "")
	fullPath := filepath.Join(baseStoragePath, storagePath, decryptedFileName)

	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return "", fmt.Errorf("file not found")
	} else if err != nil {
		return "", fmt.Errorf("error accessing file: %w", err)
	}

	if evidence.IsDownload && currentUserID != 0 {
		history := &models.IncidentEvidenceDownloadHistory{
			IncidentEvidenceID: evidence.ID,
			DownloadedAt:       time.Now(),
			DownloadedBy:       currentUserID,
		}
		if err := ies.repo.CreateDownloadHistory(history); err != nil {
			fmt.Printf("⚠️ Failed to record download history for evidence %d: %v\n", evidenceID, err)
		}
	}

	return fullPath, nil
}

func (ies *IncidentEvidenceService) GetDownloadHistories(incidentRegisterID, evidenceID uint) ([]*dtos.IncidentEvidenceDownloadHistoryListDTO, error) {
	if _, err := ies.GetByID(incidentRegisterID, evidenceID); err != nil {
		return []*dtos.IncidentEvidenceDownloadHistoryListDTO{}, err
	}

	histories, err := ies.repo.GetDownloadHistories(evidenceID)
	if err != nil {
		return []*dtos.IncidentEvidenceDownloadHistoryListDTO{}, err
	}
	if histories == nil {
		histories = []*dtos.IncidentEvidenceDownloadHistoryListDTO{}
	}
	return histories, nil
}
