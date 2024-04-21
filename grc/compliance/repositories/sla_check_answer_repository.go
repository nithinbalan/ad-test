package repositories

import (
	"jethur/grc/compliance/models"

	"gorm.io/gorm"
)

type SLACheckAnswerRepositoryInterface interface {
	Create(slaCheckAnswer *models.SLACheckAnswer) (*models.SLACheckAnswer, error)
}

type SLACheckAnswerRepository struct {
	db *gorm.DB
}

func NewSLACheckAnswerRepository(db *gorm.DB) *SLACheckAnswerRepository {
	return &SLACheckAnswerRepository{
		db: db,
	}
}

func (repo *SLACheckAnswerRepository) Create(slaCheckAnswer *models.SLACheckAnswer) (*models.SLACheckAnswer, error) {

	if err := repo.db.Create(slaCheckAnswer).Error; err != nil {
		return nil, err
	}

	return slaCheckAnswer, nil
}
