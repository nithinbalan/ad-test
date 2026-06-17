package repositories

import (
	"jethur/grc/models"

	"gorm.io/gorm"
)

type AwarenessCampaignStatusRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.AwarenessCampaignStatus, error)
	GetByID(id uint) (*models.AwarenessCampaignStatus, error)
	GetByType(status string) (*models.AwarenessCampaignStatus, error)
}

type AwarenessCampaignStatusRepository struct {
	db *gorm.DB
}

func NewAwarenessCampaignStatusRepository(db *gorm.DB) *AwarenessCampaignStatusRepository {
	return &AwarenessCampaignStatusRepository{
		db: db,
	}
}

func (repo *AwarenessCampaignStatusRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.AwarenessCampaignStatus, error) {
	var awarenessCampaignStatuses []*models.AwarenessCampaignStatus

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

	query = query.Order(sortBy + " " + sort)

	result := query.Find(&awarenessCampaignStatuses)

	return awarenessCampaignStatuses, result.Error
}

func (repo *AwarenessCampaignStatusRepository) GetByID(id uint) (*models.AwarenessCampaignStatus, error) {
	var awarenessCampaignStatus *models.AwarenessCampaignStatus

	result := repo.db.First(&awarenessCampaignStatus, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return awarenessCampaignStatus, nil
}

func (repo *AwarenessCampaignStatusRepository) GetByType(status string) (*models.AwarenessCampaignStatus, error) {
	var awarenessCampaignStatus *models.AwarenessCampaignStatus

	result := repo.db.Where("type = ?", status).First(&awarenessCampaignStatus)

	if result.Error != nil {
		return nil, result.Error
	}

	return awarenessCampaignStatus, nil
}