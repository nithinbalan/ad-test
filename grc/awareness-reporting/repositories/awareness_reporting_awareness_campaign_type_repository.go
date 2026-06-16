package repositories

import (
	"jethur/grc/models"

	"gorm.io/gorm"
)


type AwarenessCampaignTypeRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.AwarenessCampaignType, error)
	GetByID(id uint) (*models.AwarenessCampaignType, error)
	Create(awarenessCampaignType *models.AwarenessCampaignType) (*models.AwarenessCampaignType, error)
	Update(id uint, awarenessCampaignType *models.AwarenessCampaignType) (*models.AwarenessCampaignType, error)
	Delete(id uint) error
	Activate(ids []uint) ([]uint, error)
	Deactivate(ids []uint) ([]uint, error)
}

type AwarenessCampaignTypeRepository struct {
	db *gorm.DB
}

func NewAwarenessCampaignTypeRepository(db *gorm.DB) *AwarenessCampaignTypeRepository {
	return &AwarenessCampaignTypeRepository{
		db: db,
	}
}

func (repo *AwarenessCampaignTypeRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.AwarenessCampaignType, error) {
	var awarenessCampaignTypes []*models.AwarenessCampaignType
	query := repo.db

	if keyword != "" {
		query = query.Where(
			"(title_en LIKE ? OR title_ar LIKE ?)",
			"%"+keyword+"%",
			"%"+keyword+"%",
		)
	}

	if sortBy == "" {
		sortBy = "title_en"
	}

	if sort != "asc" && sort != "desc" {
		sort = "asc"
	}

	if status == "active" || status == "inactive" {
		query = query.Where("status = ?", status)
	}

	result := query.Order(sortBy + " " + sort).Find(&awarenessCampaignTypes)

	return awarenessCampaignTypes, result.Error
}

func (repo *AwarenessCampaignTypeRepository) GetByID(id uint) (*models.AwarenessCampaignType, error) {
	var awarenessCampaignType models.AwarenessCampaignType

	if err := repo.db.First(&awarenessCampaignType, id).Error; err != nil {
		return nil, err
	}

	return &awarenessCampaignType, nil
}

func (repo *AwarenessCampaignTypeRepository) Create(awarenessCampaignType *models.AwarenessCampaignType) (*models.AwarenessCampaignType, error) {
	if err := repo.db.Create(awarenessCampaignType).Error; err != nil {
		return nil, err
	}

	return awarenessCampaignType, nil
}

func (repo *AwarenessCampaignTypeRepository) Update(id uint, updatedType *models.AwarenessCampaignType) (*models.AwarenessCampaignType, error) {
	var awarenessCampaignType models.AwarenessCampaignType

	if err := repo.db.First(&awarenessCampaignType, id).Error; err != nil {
		return nil, err
	}

	if err := repo.db.Model(&awarenessCampaignType).Updates(updatedType).Error; err != nil {
		return nil, err
	}

	return &awarenessCampaignType, nil
}

func (repo *AwarenessCampaignTypeRepository) Delete(id uint) error {
	var awarenessCampaignType models.AwarenessCampaignType

	if err := repo.db.First(&awarenessCampaignType, id).Error; err != nil {
		return err
	}

	return repo.db.Delete(&awarenessCampaignType).Error
}

func (repo *AwarenessCampaignTypeRepository) Activate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint

	if err := repo.db.Model(&models.AwarenessCampaignType{}).
		Where("id IN ? AND status = ?", ids, "inactive").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err
	}

	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.AwarenessCampaignType{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "active").Error; err != nil {
			return nil, err
		}
	}

	return idsToUpdate, nil
}

func (repo *AwarenessCampaignTypeRepository) Deactivate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint

	if err := repo.db.Model(&models.AwarenessCampaignType{}).
		Where("id IN ? AND status = ?", ids, "active").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err
	}

	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.AwarenessCampaignType{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "inactive").Error; err != nil {
			return nil, err
		}
	}

	return idsToUpdate, nil
}