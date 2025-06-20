package repositories

import (
	"jethur/grc/models"

	"gorm.io/gorm"
)

type JobDescriptionCompetencyRepositoryInterface interface {
	Create(competencies []models.JobDescriptionCompetency) error
	GetByID(id uint) (*models.JobDescriptionCompetency, error)
	Update(jdc *models.JobDescriptionCompetency) error
	Delete(id uint) error
	DeleteCertifications(jdcID uint) error
	CreateCertification(cert *models.JobDescriptionCompetencyCertification) error

	GetCertificationsByCompetencyID(jdcID uint) ([]models.JobDescriptionCompetencyCertification, error)
	DeleteCertificationsByIDs(ids []uint) error
	BulkCreateCertifications(certs []models.JobDescriptionCompetencyCertification) error
	UpdateCertification(cert *models.JobDescriptionCompetencyCertification) error
}
type JobDescriptionCompetencyRepository struct {
	db *gorm.DB
}

func NewJobDescriptionCompetencyRepository(db *gorm.DB) *JobDescriptionCompetencyRepository {
	return &JobDescriptionCompetencyRepository{db: db}
}

func (repo *JobDescriptionCompetencyRepository) Create(competencies []models.JobDescriptionCompetency) error {
	return repo.db.Create(&competencies).Error
}

func (repo *JobDescriptionCompetencyRepository) GetByID(id uint) (*models.JobDescriptionCompetency, error) {
	var item models.JobDescriptionCompetency
	err := repo.db.Preload("Certifications").First(&item, id).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (repo *JobDescriptionCompetencyRepository) Update(jdc *models.JobDescriptionCompetency) error {
	return repo.db.Model(&models.JobDescriptionCompetency{}).
		Where("id = ?", jdc.ID).
		Updates(jdc).Error
}

func (repo *JobDescriptionCompetencyRepository) Delete(id uint) error {
	return repo.db.Delete(&models.JobDescriptionCompetency{}, id).Error
}

func (repo *JobDescriptionCompetencyRepository) DeleteCertifications(jdcID uint) error {
	return repo.db.Where("job_description_competency_id = ?", jdcID).
		Delete(&models.JobDescriptionCompetencyCertification{}).Error
}

func (repo *JobDescriptionCompetencyRepository) CreateCertification(cert *models.JobDescriptionCompetencyCertification) error {
	return repo.db.Create(cert).Error
}

func (repo *JobDescriptionCompetencyRepository) GetCertificationsByCompetencyID(jdcID uint) ([]models.JobDescriptionCompetencyCertification, error) {
	var certs []models.JobDescriptionCompetencyCertification
	err := repo.db.Where("job_description_competency_id = ?", jdcID).Find(&certs).Error
	return certs, err
}

func (repo *JobDescriptionCompetencyRepository) DeleteCertificationsByIDs(ids []uint) error {
	return repo.db.Where("id IN ?", ids).Delete(&models.JobDescriptionCompetencyCertification{}).Error
}

func (repo *JobDescriptionCompetencyRepository) BulkCreateCertifications(certs []models.JobDescriptionCompetencyCertification) error {
	return repo.db.Create(&certs).Error
}

func (repo *JobDescriptionCompetencyRepository) UpdateCertification(cert *models.JobDescriptionCompetencyCertification) error {
	return repo.db.Save(cert).Error
}
