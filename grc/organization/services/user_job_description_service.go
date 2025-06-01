package services

import (
	"fmt"
	"jethur/grc/logs/services"
	"jethur/grc/models"
	"jethur/grc/organization/dtos"
	"jethur/grc/organization/repositories"
	"jethur/utility"
)

type UserJobDescriptionServiceInterface interface {
	Create(userID uint, input *dtos.AssignJobDescriptionsDTO, currentUserID uint) error
	Delete(userID, jobDescriptionID, currentUserID uint) error
}

type UserJobDescriptionService struct {
	repo           repositories.UserJobDescriptionRepositoryInterface
	activityLogSvc services.ActivityLogServiceInterface
	jdService      JobDescriptionServiceInterface
}

func NewUserJobDescriptionService(repo repositories.UserJobDescriptionRepositoryInterface,
	activityLogSvc services.ActivityLogServiceInterface,
	jdService JobDescriptionServiceInterface) UserJobDescriptionServiceInterface {
	return &UserJobDescriptionService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
		jdService:      jdService,
	}
}
func (s *UserJobDescriptionService) Create(userID uint, input *dtos.AssignJobDescriptionsDTO, currentUserID uint) error {
	if err := s.validateJobDescriptionIDs(input.JobDescriptionIDs); err != nil {
		return err
	}
	var userJDs []models.UserJobDescription
	var userJDCs []models.UserJobDescriptionCompetency
	var userJDCerts []models.UserJobDescriptionCertificate

	for _, jdID := range input.JobDescriptionIDs {
		userJD := models.UserJobDescription{
			UserID:           userID,
			JobDescriptionID: jdID,
			CreatedByID:      currentUserID,
		}
		userJDs = append(userJDs, userJD)
	}

	tx := s.repo.BeginTransaction()
	tx.Rollback()

	if err := s.repo.CreateUserJobDescriptions(tx, &userJDs); err != nil {
		tx.Rollback()
		return err
	}

	for _, userJD := range userJDs {
		jdDetails, err := s.jdService.GetFullDetailsByID(userJD.JobDescriptionID)
		if err != nil {
			tx.Rollback()
			return err
		}

		for _, comp := range jdDetails.Competencies {
			ujdc := models.UserJobDescriptionCompetency{
				UserJobDescriptionID:       userJD.ID,
				JobDescriptionCompetencyID: comp.ID,
				ActualLevelID:              *comp.CompetencyLevelID,
			}
			userJDCs = append(userJDCs, ujdc)

			for _, cert := range comp.Certifications {
				ujdcert := models.UserJobDescriptionCertificate{
					UserJobDescriptionCompetencyID:          0,
					JobDescriptionCompetencyCertificationID: cert.ID,
					CertificateTitle:                        cert.Title,
					CreatedByID:                             &currentUserID,
				}
				userJDCerts = append(userJDCerts, ujdcert)
			}
		}
	}

	if err := s.repo.CreateUserJobDescriptionCompetencies(tx, &userJDCs); err != nil {
		tx.Rollback()
		return err
	}

	for i := range userJDCerts {
		for _, ujdc := range userJDCs {
			if userJDCerts[i].JobDescriptionCompetencyCertificationID == ujdc.JobDescriptionCompetencyID {
				userJDCerts[i].UserJobDescriptionCompetencyID = ujdc.ID
				break
			}
		}
	}

	if err := s.repo.CreateUserJobDescriptionCertificates(tx, &userJDCerts); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (s *UserJobDescriptionService) Delete(userID, jobDescriptionID, currentUserID uint) error {
	userJD, err := s.repo.GetByUserAndJD(userID, jobDescriptionID)
	if err != nil {
		return utility.NewResourceNotFoundError(err, jobDescriptionID)
	}
	tx := s.repo.BeginTransaction()

	if err := s.repo.DeleteUserJobDescriptionCertificatesByUserJDID(tx, userJD.ID); err != nil {
		tx.Rollback()
		return err
	}

	if err := s.repo.DeleteUserJobDescriptionCompetenciesByUserJDID(tx, userJD.ID); err != nil {
		tx.Rollback()
		return err
	}

	if err := s.repo.DeleteUserJobDescription(tx, userJD.ID); err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	go func() {
		_ = s.activityLogSvc.LogActivity(
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

func (s *UserJobDescriptionService) validateJobDescriptionIDs(jdIDs []uint) error {
	seen := make(map[uint]bool)
	for _, jdID := range jdIDs {
		if seen[jdID] {
			return utility.NewValidationError(fmt.Errorf("duplicate job description ID: %d", jdID), "duplicate", "job_description_ids")
		}
		if _, err := s.jdService.GetByID(jdID); err != nil {
			return utility.NewValidationError(fmt.Errorf("job description does not exist for ID: %d", jdID), "not_exist", "job_description_ids")
		}
		seen[jdID] = true
	}
	return nil
}
