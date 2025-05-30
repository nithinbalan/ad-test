package services

import (
	"fmt"
	"jethur/grc/logs/services"
	"jethur/grc/models"
	"jethur/grc/organization/dtos"
	"jethur/grc/organization/repositories"
	"jethur/utility"
)

type JobDescriptionCompetencyServiceInterface interface {
	MapCompetencies(jobDescriptionID uint, currentUserID uint, input *dtos.JobDescriptionMapCompetenciesDTO) error
	UpdateCompetencyDetails(jdcID uint, currentUserID uint, input *dtos.UpdateJDCompetencyDTO) error
	DeleteCompetency(jdcID uint, currentUserID uint) error
}

type JobDescriptionCompetencyService struct {
	repo               repositories.JobDescriptionCompetencyRepositoryInterface
	jdService          JobDescriptionServiceInterface
	competencyLevelSvc CompetencyLevelServiceInterface
	competencySvc      CompetencyServiceInterface
	activityLogSvc     services.ActivityLogServiceInterface
}

func NewJobDescriptionCompetencyService(
	repo repositories.JobDescriptionCompetencyRepositoryInterface,
	jdService JobDescriptionServiceInterface,
	competencyLevelSvc CompetencyLevelServiceInterface,
	competencySvc CompetencyServiceInterface,
	activityLogSvc services.ActivityLogServiceInterface,
) JobDescriptionCompetencyServiceInterface {
	return &JobDescriptionCompetencyService{
		repo:               repo,
		jdService:          jdService,
		competencyLevelSvc: competencyLevelSvc,
		competencySvc:      competencySvc,
		activityLogSvc:     activityLogSvc,
	}
}

func (jdcs *JobDescriptionCompetencyService) MapCompetencies(jobDescriptionID uint, currentUserID uint, input *dtos.JobDescriptionMapCompetenciesDTO) error {
	// Validate JD existence via service
	if _, err := jdcs.jdService.GetByID(jobDescriptionID); err != nil {
		return err
	}

	// Validate competencies
	if err := jdcs.validateCompetencyIDs(input.CompetencyIDs); err != nil {
		return err
	}

	var jdCompetencies []models.JobDescriptionCompetency
	for _, compID := range input.CompetencyIDs {
		jdCompetencies = append(jdCompetencies, models.JobDescriptionCompetency{
			JobDescriptionID: jobDescriptionID,
			CompetencyID:     compID,
			CreatedByID:      currentUserID,
		})
	}

	if len(jdCompetencies) > 0 {
		if err := jdcs.repo.Create(jdCompetencies); err != nil {
			return err
		}
	}
	go func() {
		ids := make([]uint, len(jdCompetencies))
		for i, jc := range jdCompetencies {
			ids[i] = jc.ID
		}
		err := jdcs.activityLogSvc.LogBulkActivities(models.CreateAction, models.JobDescriptionCompetency{}.TableName(), ids, jdCompetencies, currentUserID, nil)
		if err != nil {
			utility.LogError("Failed to log JD Competency map activity", map[string]interface{}{"error": err})
		}
	}()
	return nil
}

func (jdcs *JobDescriptionCompetencyService) UpdateCompetencyDetails(jdcID uint, currentUserID uint, input *dtos.UpdateJDCompetencyDTO) error {
	// Fetch existing competency
	existing, err := jdcs.repo.GetByID(jdcID)
	if err != nil {
		return utility.NewResourceNotFoundError(err, jdcID)
	}
	oldData := map[string]interface{}{
		"competency_level_id": existing.CompetencyLevelID,
	}
	existing.CompetencyLevelID = &input.CompetencyLevelID
	existing.UpdatedByID = &currentUserID

	if err := jdcs.repo.Update(existing); err != nil {
		return err
	}

	// Fetch existing certifications from DB
	existingCerts, err := jdcs.repo.GetCertificationsByCompetencyID(jdcID)
	if err != nil {
		return err
	}

	// Create maps for comparison
	inputMap := make(map[string]dtos.JobDescriptionCompetencyCertFormDTO)
	for _, cert := range input.Certifications {
		inputMap[cert.Title] = cert
	}

	existingMap := make(map[string]models.JobDescriptionCompetencyCertification)
	for _, cert := range existingCerts {
		existingMap[cert.Title] = cert
	}

	var toCreate []models.JobDescriptionCompetencyCertification
	var toDeleteIDs []uint
	var toUpdate []models.JobDescriptionCompetencyCertification

	// Identify certs to create or update
	for title, inputCert := range inputMap {
		if existingCert, exists := existingMap[title]; exists {
			// Update if IsMandatory changed
			if existingCert.IsMandatory != inputCert.IsMandatory {
				existingCert.IsMandatory = inputCert.IsMandatory
				toUpdate = append(toUpdate, existingCert)
			}
			// else: do nothing (no change)
		} else {
			// New cert
			toCreate = append(toCreate, models.JobDescriptionCompetencyCertification{
				JobDescriptionCompetencyID: jdcID,
				Title:                      inputCert.Title,
				IsMandatory:                inputCert.IsMandatory,
				CreatedByID:                currentUserID,
			})
		}
	}

	// Identify certs to delete
	for title, existingCert := range existingMap {
		if _, stillExists := inputMap[title]; !stillExists {
			toDeleteIDs = append(toDeleteIDs, existingCert.ID)
		}
	}

	// Perform DB operations
	if len(toDeleteIDs) > 0 {
		if err := jdcs.repo.DeleteCertificationsByIDs(toDeleteIDs); err != nil {
			return err
		}
	}

	if len(toCreate) > 0 {
		if err := jdcs.repo.BulkCreateCertifications(toCreate); err != nil {
			return err
		}
	}

	for _, cert := range toUpdate {
		if err := jdcs.repo.UpdateCertification(&cert); err != nil {
			return err
		}
	}
	changed := utility.CompareFieldChanges(oldData, map[string]interface{}{
		"competency_level_id": input.CompetencyLevelID,
	})
	go func() {
		if len(changed) > 0 {
			_ = jdcs.activityLogSvc.LogActivity(models.UpdateAction, models.JobDescriptionCompetency{}.TableName(), jdcID, changed, currentUserID, nil)
		}
	}()
	return nil
}

func (jdcs *JobDescriptionCompetencyService) DeleteCompetency(jdcID uint, currentUserID uint) error {
	existing, err := jdcs.repo.GetByID(jdcID)
	if err != nil {
		return utility.NewResourceNotFoundError(err, jdcID)
	}

	// Delete certifications first to avoid FK conflict
	if err := jdcs.repo.DeleteCertifications(jdcID); err != nil {
		return err
	}
	go func() {
		_ = jdcs.activityLogSvc.LogActivity(
			models.DeleteAction,
			models.JobDescriptionCompetency{}.TableName(),
			jdcID,
			existing,
			currentUserID,
			nil,
		)
	}()
	// Then delete the competency
	return jdcs.repo.Delete(jdcID)
}

func (jdcs *JobDescriptionCompetencyService) validateCompetencyIDs(competencyIDs []uint) error {
	seen := make(map[uint]bool)
	for _, compID := range competencyIDs {
		if seen[compID] {
			return utility.NewValidationError(fmt.Errorf("duplicate competency ID: %d", compID), "duplicate", "competency_ids")
		}
		if _, err := jdcs.competencySvc.GetByID(compID); err != nil {
			return utility.NewValidationError(fmt.Errorf("competency does not exist for ID: %d", compID), "not_exist", "competency_ids")
		}
		seen[compID] = true
	}
	return nil
}
