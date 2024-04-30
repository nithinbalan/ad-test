package repositories

import (
	"jethur/grc/models"

	"gorm.io/gorm"
)

type SLACheckAnswerRepositoryInterface interface {
	GetByID(id uint) (*models.SLACheckAnswer, error)
	Create(slaCheckAnswer *models.SLACheckAnswer) (*models.SLACheckAnswer, error)
	Update(slaCheckAnswer *models.SLACheckAnswer) error
	Delete(id uint) error
}

type SLACheckAnswerRepository struct {
	db *gorm.DB
}

func NewSLACheckAnswerRepository(db *gorm.DB) *SLACheckAnswerRepository {
	return &SLACheckAnswerRepository{
		db: db,
	}
}

func (repo *SLACheckAnswerRepository) GetByID(id uint) (*models.SLACheckAnswer, error) {
	var slaCheckAnswer *models.SLACheckAnswer
	result := repo.db.First(&slaCheckAnswer, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return slaCheckAnswer, nil
}

func (repo *SLACheckAnswerRepository) Create(slaCheckAnswer *models.SLACheckAnswer) (*models.SLACheckAnswer, error) {

	if err := repo.db.Create(slaCheckAnswer).Error; err != nil {
		return nil, err
	}

	return slaCheckAnswer, nil
}

func (repo *SLACheckAnswerRepository) Update(slaCheckAnswer *models.SLACheckAnswer) error {

	if err := repo.db.Save(slaCheckAnswer).Error; err != nil {
		return err
	}

	return nil
}

func (repo *SLACheckAnswerRepository) Delete(id uint) error {

	var slaCheckAnswer *models.SLACheckAnswer
	if err := repo.db.First(&slaCheckAnswer, id).Error; err != nil {
		return err
	}

	result := repo.db.Delete(&models.SLACheckAnswer{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
