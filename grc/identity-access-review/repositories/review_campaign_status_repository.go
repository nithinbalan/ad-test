package repositories

import (
	"jethur/grc/models"

	"gorm.io/gorm"
)

type ReviewCampaignStatusRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ReviewCampaignStatus, error)
	GetByID(id uint) (*models.ReviewCampaignStatus, error)
	GetByType(reviewCampaignStatusType string) (*models.ReviewCampaignStatus, error)
}

type ReviewCampaignStatusRepository struct {
	db *gorm.DB
}

func NewReviewCampaignStatusRepository(db *gorm.DB) *ReviewCampaignStatusRepository {
	return &ReviewCampaignStatusRepository{
		db: db,
	}
}

func (repo *ReviewCampaignStatusRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ReviewCampaignStatus, error) {
	var reviewCampaignStatuses []*models.ReviewCampaignStatus

	query := repo.db

	if keyword != "" {
		query = query.Where(
			"(title_en LIKE ? OR title_ar LIKE ?)",
			"%"+keyword+"%",
			"%"+keyword+"%",
		)
	}

	if sortBy != "id" && sortBy != "title_en" && sortBy != "title_ar" && sortBy != "color_code" && sortBy != "type" && sortBy != "status" {
		sortBy = "title_en"
	}

	if sort != "asc" && sort != "desc" {
		sort = "asc"
	}

	if status == "active" || status == "inactive" {
		query = query.Where("status = ?", status)
	}

	query = query.Order(sortBy + " " + sort)

	result := query.Find(&reviewCampaignStatuses)

	return reviewCampaignStatuses, result.Error
}

func (repo *ReviewCampaignStatusRepository) GetByID(id uint) (*models.ReviewCampaignStatus, error) {
	var reviewCampaignStatus *models.ReviewCampaignStatus

	result := repo.db.First(&reviewCampaignStatus, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return reviewCampaignStatus, nil
}

func (repo *ReviewCampaignStatusRepository) GetByType(reviewCampaignStatusType string) (*models.ReviewCampaignStatus, error) {
	var reviewCampaignStatus *models.ReviewCampaignStatus

	result := repo.db.Where("type = ?", reviewCampaignStatusType).First(&reviewCampaignStatus)

	if result.Error != nil {
		return nil, result.Error
	}

	return reviewCampaignStatus, nil
}
