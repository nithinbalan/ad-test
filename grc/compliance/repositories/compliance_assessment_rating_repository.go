package repositories

import (
	"errors"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type ComplianceAssessmentRatingRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ComplianceAssessmentRating, error)
	GetByID(id uint) (*models.ComplianceAssessmentRating, error)
	GetByType(ratingType string) (*models.ComplianceAssessmentRating, error)
}

type ComplianceAssessmentRatingRepository struct {
	db *gorm.DB
}

func NewComplianceAssessmentRatingRepository(db *gorm.DB) *ComplianceAssessmentRatingRepository {
	return &ComplianceAssessmentRatingRepository{
		db: db,
	}
}

func (repo *ComplianceAssessmentRatingRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ComplianceAssessmentRating, error) {
	var ratings []*models.ComplianceAssessmentRating
	query := repo.db

	if keyword != "" {
		query = query.Where("(title_en LIKE ? OR title_ar LIKE ?)", "%"+keyword+"%", "%"+keyword+"%")
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
	result := query.Find(&ratings)

	return ratings, result.Error
}

func (repo *ComplianceAssessmentRatingRepository) GetByID(id uint) (*models.ComplianceAssessmentRating, error) {
	var rating *models.ComplianceAssessmentRating
	result := repo.db.First(&rating, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return rating, nil
}

func (repo *ComplianceAssessmentRatingRepository) GetByType(ratingType string) (*models.ComplianceAssessmentRating, error) {
	var rating *models.ComplianceAssessmentRating
	result := repo.db.Where("type = ?", ratingType).First(&rating)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, "Compliance Assessment Rating with type "+ratingType)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return rating, nil
}
