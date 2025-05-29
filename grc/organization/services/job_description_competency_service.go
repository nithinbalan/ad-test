package services

import (
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
	repo      repositories.JobDescriptionCompetencyRepositoryInterface
	jdService JobDescriptionServiceInterface
}

func NewJobDescriptionCompetencyService(
	repo repositories.JobDescriptionCompetencyRepositoryInterface,
	jdService JobDescriptionServiceInterface,
) JobDescriptionCompetencyServiceInterface {
	return &JobDescriptionCompetencyService{
		repo:      repo,
		jdService: jdService,
	}
}

func (jdcs *JobDescriptionCompetencyService) MapCompetencies(jobDescriptionID uint, currentUserID uint, input *dtos.JobDescriptionMapCompetenciesDTO) error {
	// Validate JD existence via service
	if _, err := jdcs.jdService.GetByID(jobDescriptionID); err != nil {
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

	return nil
}

func (jdcs *JobDescriptionCompetencyService) UpdateCompetencyDetails(jdcID uint, currentUserID uint, input *dtos.UpdateJDCompetencyDTO) error {
	existing, err := jdcs.repo.GetByID(jdcID)
	if err != nil {
		return utility.NewResourceNotFoundError(err, jdcID)
	}

	existing.CompetencyLevelID = input.CompetencyLevelID
	existing.UpdatedByID = &currentUserID

	if err := jdcs.repo.Update(existing); err != nil {
		return err
	}

	if err := jdcs.repo.DeleteCertifications(jdcID); err != nil {
		return err
	}

	for _, cert := range input.Certifications {
		certModel := models.JobDescriptionCompetencyCertification{
			JobDescriptionCompetencyID: jdcID,
			Title:                      cert.Title,
			IsMandatory:                cert.IsMandatory,
			CreatedByID:                currentUserID,
		}
		if err := jdcs.repo.CreateCertification(&certModel); err != nil {
			return err
		}
	}

	return nil
}

func (jdcs *JobDescriptionCompetencyService) DeleteCompetency(jdcID uint, currentUserID uint) error {
	if _, err := jdcs.repo.GetByID(jdcID); err != nil {
		return utility.NewResourceNotFoundError(err, jdcID)
	}
	return jdcs.repo.Delete(jdcID)
}
