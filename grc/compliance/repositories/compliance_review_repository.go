package repositories

import (
	complianceDTOs "jethur/grc/compliance/dtos"
	"jethur/grc/models"

	"gorm.io/gorm"
)

type ComplianceReviewRepositoryInterface interface {
	GetAll(complianceID uint, sort string, sortBy string, keyword string) ([]*complianceDTOs.ComplianceReviewListDTO, error)
	GetByID(id uint) (*models.ComplianceReview, error)
	Create(ComplianceReview *models.ComplianceReview) (uint, error)
	Update(ComplianceReview *models.ComplianceReview) error
}

type ComplianceReviewRepository struct {
	db *gorm.DB
}

func NewComplianceReviewRepository(db *gorm.DB) *ComplianceReviewRepository {
	return &ComplianceReviewRepository{
		db: db,
	}
}

func (repo *ComplianceReviewRepository) GetAll(complianceID uint, sort string, sortBy string, keyword string) ([]*complianceDTOs.ComplianceReviewListDTO, error) {

	var complianceReviews []*complianceDTOs.ComplianceReviewListDTO

	query := repo.db.Table("compliance_reviews").
		Select(`
            compliance_reviews.*,
			compliance_register.title_en AS compliance_title_en,
			compliance_register.title_ar AS compliance_title_ar,
			compliance_register.description_en AS compliance_description_en,
			compliance_register.description_ar AS compliance_description_ar,
			reviewed_user.first_name_en AS reviewed_user_first_name_en,
            reviewed_user.first_name_ar AS reviewed_user_first_name_ar,
            reviewed_user.last_name_en AS reviewed_user_last_name_en,
            reviewed_user.last_name_ar AS reviewed_user_last_name_ar 
        `).
		Joins(`
		    LEFT JOIN compliance_register ON compliance_register.id = compliance_reviews.compliance_register_id
			LEFT JOIN users AS reviewed_user ON reviewed_user.id = compliance_reviews.reviewed_by
		`).
		Where("compliance_reviews.compliance_register_id = ?", complianceID)

	if keyword != "" {
		query = query.Where(`
			(compliance_register.title LIKE ? OR
			compliance_reviews.comment LIKE ?)`,
			"%"+keyword+"%", "%"+keyword+"%")
	}
	if sortBy == "" {
		sortBy = "compliance_reviews.id"
	}
	if sort != "asc" && sort != "desc" {
		sort = "asc"
	}
	query = query.Order(sortBy + " " + sort)
	result := query.Find(&complianceReviews)

	return complianceReviews, result.Error
}

func (repo *ComplianceReviewRepository) GetByID(id uint) (*models.ComplianceReview, error) {
	var complianceReview *models.ComplianceReview
	result := repo.db.First(&complianceReview, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return complianceReview, nil
}

func (repo *ComplianceReviewRepository) Create(complianceReview *models.ComplianceReview) (uint, error) {

	if err := repo.db.Create(complianceReview).Error; err != nil {
		return 0, err
	}

	return complianceReview.ID, nil
}

func (repo *ComplianceReviewRepository) Update(ComplianceReview *models.ComplianceReview) error {
	if err := repo.db.Save(ComplianceReview).Error; err != nil {
		return err
	}
	return nil
}
