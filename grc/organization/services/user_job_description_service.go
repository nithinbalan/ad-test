package services

import (
	"errors"
	"fmt"
	"jethur/grc/logs/services"
	"jethur/grc/models"
	"jethur/grc/organization/dtos"
	"jethur/grc/organization/repositories"
	"jethur/utility"
	"path/filepath"
	"strconv"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type UserJobDescriptionServiceInterface interface {
	GetAll(sort string, sortBy string, keyword string, userID uint) ([]*dtos.UserJobDescriptionListDTO, error)
	GetFullDetailsByID(id uint) (*dtos.UserJobDescriptionDetailDTO, error)
	Create(userID uint, input *dtos.AssignJobDescriptionsDTO, currentUserID uint) (uint, error)
	Delete(userID, jobDescriptionID, currentUserID uint) error
	Update(jdcID uint, dto dtos.EvaluateUserJDCompetencyDTO, currentUserID uint) error

	CreateUserJDFromDesignation(userID uint, designationID uint, currentUserID uint, tx utility.DatabaseTransaction) error
	DeleteUserJDFromDesignation(userID uint, designationID uint, currentUserID uint) error
}

type UserJobDescriptionService struct {
	repo           repositories.UserJobDescriptionRepositoryInterface
	activityLogSvc services.ActivityLogServiceInterface
	jdService      JobDescriptionServiceInterface
	designationSvc DesignationServiceInterface
}

func NewUserJobDescriptionService(repo repositories.UserJobDescriptionRepositoryInterface,
	activityLogSvc services.ActivityLogServiceInterface,
	jdService JobDescriptionServiceInterface,
	designationSvc DesignationServiceInterface) UserJobDescriptionServiceInterface {
	return &UserJobDescriptionService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
		jdService:      jdService,
		designationSvc: designationSvc,
	}
}

func (ujds *UserJobDescriptionService) GetAll(sort string, sortBy string, keyword string, userID uint) ([]*dtos.UserJobDescriptionListDTO, error) {
	return ujds.repo.GetAll(sort, sortBy, keyword, userID)
}

func (ujds *UserJobDescriptionService) GetFullDetailsByID(id uint) (*dtos.UserJobDescriptionDetailDTO, error) {
	model, err := ujds.repo.GetFullDetailsByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	var dto dtos.UserJobDescriptionDetailDTO
	err = mapstructure.Decode(model, &dto)
	if err != nil {
		return nil, err
	}
	return &dto, nil
}

func (ujds *UserJobDescriptionService) Create(userID uint, input *dtos.AssignJobDescriptionsDTO, currentUserID uint) (uint, error) {
	if err := ujds.validateJobDescriptionIDs(input.JobDescriptionIDs); err != nil {
		return 0, err
	}

	var userJDs []*models.UserJobDescription
	var userJDCerts []*models.UserJobDescriptionCertificate

	// Step 1: Prepare UserJobDescriptions
	for _, jdID := range input.JobDescriptionIDs {
		userJD := models.UserJobDescription{
			UserID:           userID,
			JobDescriptionID: jdID,
			CreatedByID:      currentUserID,
		}
		userJDs = append(userJDs, &userJD)
	}

	tx := ujds.repo.BeginTransaction()

	// Step 2: Insert UserJobDescriptions
	if err := ujds.repo.CreateUserJobDescriptions(tx, userJDs); err != nil {
		tx.Rollback()
		return 0, err
	}

	// Step 3: Insert competencies and prepare certificates
	for _, userJD := range userJDs {
		jdDetails, err := ujds.jdService.GetFullDetailsByID(userJD.JobDescriptionID)
		if err != nil {
			tx.Rollback()
			return 0, err
		}

		for _, comp := range jdDetails.Competencies {
			ujdc := models.UserJobDescriptionCompetency{
				UserJobDescriptionID:   userJD.ID,
				CompetencyID:           comp.CompetencyID,
				CompetencyTitleEn:      &comp.Competency.TitleEn,
				CompetencyTitleAr:      &comp.Competency.TitleAr,
				CompetencyGroupTitleEn: &comp.Competency.CompetencyGroup.TitleEn,
				CompetencyGroupTitleAr: &comp.Competency.CompetencyGroup.TitleAr,
			}
			if err := ujds.repo.CreateUserJobDescriptionCompetency(tx, &ujdc); err != nil {
				tx.Rollback()
				return 0, err
			}

			for _, cert := range comp.Certifications {
				ujdcert := models.UserJobDescriptionCertificate{
					UserJobDescriptionCompetencyID:          ujdc.ID,
					JobDescriptionCompetencyCertificationID: cert.ID,
					CertificateTitle:                        cert.Title,
					CreatedByID:                             &currentUserID,
				}
				userJDCerts = append(userJDCerts, &ujdcert)
			}
		}
	}

	// Step 4: Bulk insert certificates
	if len(userJDCerts) > 0 {
		if err := ujds.repo.CreateUserJobDescriptionCertificates(tx, userJDCerts); err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	tx.Commit()
	return userID, nil
}

func (ujds *UserJobDescriptionService) Delete(userID, jobDescriptionID, currentUserID uint) error {
	userJD, err := ujds.repo.GetByUserAndJD(userID, jobDescriptionID)
	if err != nil {
		return utility.NewResourceNotFoundError(err, jobDescriptionID)
	}

	tx := ujds.repo.BeginTransaction()
	certs, err := ujds.repo.GetCertificatesByUserJDID(userJD.ID)
	if err != nil {
		tx.Rollback()
		return err
	}
	if err := ujds.repo.DeleteUserJobDescriptionCertificatesByUserJDID(tx, userJD.ID); err != nil {
		tx.Rollback()
		return err
	}

	if err := ujds.repo.DeleteUserJobDescriptionCompetenciesByUserJDID(tx, userJD.ID); err != nil {
		tx.Rollback()
		return err
	}

	if err := ujds.repo.DeleteUserJobDescription(tx, userJD.ID); err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	go func() {
		for _, cert := range certs {
			if cert.Name != nil {
				storagePath := filepath.Join("job-description", "competency", strconv.Itoa(int(userJD.ID)), "certificates")
				_ = utility.RemoveFile(*cert.Name, storagePath)
			}
		}

		_ = ujds.activityLogSvc.LogActivity(
			models.DeleteAction,
			models.UserJobDescription{}.TableName(),
			userJD.ID,
			userJD,
			currentUserID,
			nil,
		)
	}()

	return nil
}

func (ujds *UserJobDescriptionService) Update(jdcID uint, dto dtos.EvaluateUserJDCompetencyDTO, currentUserID uint) error {
	comp, err := ujds.repo.GetCompetencyByID(jdcID)
	if err != nil {
		return utility.NewResourceNotFoundError(err, jdcID)
	}

	comp.ActualLevelID = &dto.ActualLevelID
	comp.Comments = dto.Comments

	tx := ujds.repo.BeginTransaction()

	if err := ujds.repo.UpdateUserJDCompetency(comp, tx); err != nil {
		tx.Rollback()
		return err
	}

	// --- Handle Certificate Docs ---
	var certsToUpdate []*models.UserJobDescriptionCertificate
	var certFilesToStore []*models.UserJobDescriptionCertificate
	var certFilesToDelete []*models.UserJobDescriptionCertificate

	for _, cert := range dto.Documents {
		// Load the existing certificate model by ID
		existingDoc, err := ujds.repo.GetCertificateByID(cert.ID)
		if err != nil {
			tx.Rollback()
			return err
		}

		if cert.IsNewFile != nil && *cert.IsNewFile {
			// Update existing record with new file metadata
			encrypted, err := utility.EncryptData([]byte(cert.FileName))
			if err != nil {
				tx.Rollback()
				return err
			}
			url := utility.GenerateDocumentURL("job-description", "competency", strconv.Itoa(int(jdcID)), "certificates")

			existingDoc.CertificateTitle = cert.FileName
			existingDoc.Name = &cert.FileName
			existingDoc.Token = &encrypted
			existingDoc.OriginalName = &cert.OriginalName
			existingDoc.URL = &url
			existingDoc.Type = &cert.Type
			existingDoc.Size = &cert.Size

			certsToUpdate = append(certsToUpdate, existingDoc)
			certFilesToStore = append(certFilesToStore, existingDoc)

		} else if cert.IsDeletedFile != nil && *cert.IsDeletedFile {
			// Only clear document-related fields
			existingDoc.Token = nil
			existingDoc.Name = nil
			existingDoc.OriginalName = nil
			existingDoc.URL = nil
			existingDoc.Type = nil
			existingDoc.Size = nil

			certsToUpdate = append(certsToUpdate, existingDoc)
			certFilesToDelete = append(certFilesToDelete, existingDoc)
		}
	}

	// Insert new files
	if len(certsToUpdate) > 0 {
		if err := ujds.repo.UpdateUserJobDescriptionCertificates(tx, certsToUpdate); err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, doc := range certFilesToStore {
		storagePath := filepath.Join("job-description", "competency", strconv.Itoa(int(jdcID)), "certificates")
		if err := utility.StoreFile(*doc.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}
	for _, doc := range certFilesToDelete {
		storagePath := filepath.Join("job-description", "competency", strconv.Itoa(int(jdcID)), "certificates")
		if err := utility.RemoveFile(*doc.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (ujds *UserJobDescriptionService) CreateUserJDFromDesignation(userID uint, designationID uint, currentUserID uint, tx utility.DatabaseTransaction) error {
	designation, err := ujds.designationSvc.GetByID(designationID)
	if err != nil {
		return err
	}

	if designation.JobDescriptionID == nil {
		return nil // No JD attached to designation
	}

	// Step 1: Create UserJobDescription as slice (to reuse bulk insert logic)
	userJDs := []*models.UserJobDescription{
		{
			UserID:           userID,
			JobDescriptionID: *designation.JobDescriptionID,
			CreatedByID:      currentUserID,
		},
	}

	if err := ujds.repo.CreateUserJobDescriptions(tx, userJDs); err != nil {
		return err
	}

	// Step 2: Fetch competencies and certifications for that JD
	jdDetails, err := ujds.jdService.GetFullDetailsByID(*designation.JobDescriptionID)
	if err != nil {
		return err
	}

	var userJDCerts []*models.UserJobDescriptionCertificate

	// Step 3: Prepare competencies and certificates
	for _, comp := range jdDetails.Competencies {
		ujdc := &models.UserJobDescriptionCompetency{
			UserJobDescriptionID:   userJDs[0].ID,
			CompetencyID:           comp.CompetencyID,
			CompetencyTitleEn:      &comp.Competency.TitleEn,
			CompetencyTitleAr:      &comp.Competency.TitleAr,
			CompetencyGroupTitleEn: &comp.Competency.CompetencyGroup.TitleEn,
			CompetencyGroupTitleAr: &comp.Competency.CompetencyGroup.TitleAr,
		}

		if err := ujds.repo.CreateUserJobDescriptionCompetency(tx, ujdc); err != nil {
			return err
		}

		for _, cert := range comp.Certifications {
			ujdcert := &models.UserJobDescriptionCertificate{
				UserJobDescriptionCompetencyID:          ujdc.ID,
				JobDescriptionCompetencyCertificationID: cert.ID,
				CertificateTitle:                        cert.Title,
				CreatedByID:                             &currentUserID,
			}
			userJDCerts = append(userJDCerts, ujdcert)
		}
	}

	// Step 4: Bulk insert certificates
	if len(userJDCerts) > 0 {
		if err := ujds.repo.CreateUserJobDescriptionCertificates(tx, userJDCerts); err != nil {
			return err
		}
	}

	return nil
}

func (ujds *UserJobDescriptionService) DeleteUserJDFromDesignation(userID uint, designationID uint, currentUserID uint) error {
	designation, err := ujds.designationSvc.GetByID(designationID)
	if err != nil {
		return err
	}

	if designation.JobDescriptionID == nil {
		return nil // No JD linked
	}

	userJD, err := ujds.repo.GetByUserAndJD(userID, *designation.JobDescriptionID)
	if err != nil {
		return utility.NewResourceNotFoundError(err, *designation.JobDescriptionID)
	}

	tx := ujds.repo.BeginTransaction()

	// Load certs before deletion
	certs, err := ujds.repo.GetCertificatesByUserJDID(userJD.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := ujds.repo.DeleteUserJobDescriptionCertificatesByUserJDID(tx, userJD.ID); err != nil {
		tx.Rollback()
		return err
	}

	if err := ujds.repo.DeleteUserJobDescriptionCompetenciesByUserJDID(tx, userJD.ID); err != nil {
		tx.Rollback()
		return err
	}

	if err := ujds.repo.DeleteUserJobDescription(tx, userJD.ID); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	go func() {
		for _, cert := range certs {
			if cert.Name != nil {
				storagePath := filepath.Join("job-description", "competency", strconv.Itoa(int(userJD.ID)), "certificates")
				_ = utility.RemoveFile(*cert.Name, storagePath)
			}
		}

		_ = ujds.activityLogSvc.LogActivity(
			models.DeleteAction,
			models.UserJobDescription{}.TableName(),
			userJD.ID,
			userJD,
			currentUserID,
			nil,
		)
	}()

	return nil
}

func (ujds *UserJobDescriptionService) validateJobDescriptionIDs(jdIDs []uint) error {
	seen := make(map[uint]bool)

	for _, jdID := range jdIDs {
		if seen[jdID] {
			return utility.NewValidationError(fmt.Errorf("duplicate job description ID: %d", jdID), "duplicate", "job_description_ids")
		}
		seen[jdID] = true

		// Load full JD with Competencies and their Certifications
		jdDetails, err := ujds.jdService.GetFullDetailsByID(jdID)
		if err != nil {
			return utility.NewValidationError(fmt.Errorf("unable to load full details for job description ID: %d", jdID), "load_error", "job_description_ids")
		}

		if len(jdDetails.Competencies) == 0 {
			return utility.NewValidationError(fmt.Errorf("no competencies found for job description ID: %d", jdID), "no_competencies", "job_description_ids")
		}

	}

	return nil
}
