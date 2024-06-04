package repositories

import (
	complianceDTOs "jethur/grc/compliance/dtos"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type ComplianceReviewRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(complianceID uint, sort string, sortBy string, keyword string) ([]*complianceDTOs.ComplianceReviewListDTO, error)
	GetByID(id uint) (*models.ComplianceReview, error)
	Create(ComplianceReview *models.ComplianceReview, tx ...utility.DatabaseTransaction) (uint, error)
	Update(ComplianceReview *models.ComplianceReview, tx ...utility.DatabaseTransaction) error

	//Document
	CreateDocuments(docs []*models.ComplianceReviewDocument, tx ...utility.DatabaseTransaction) error
	GetDocumentByID(id uint) (*models.ComplianceReviewDocument, error)
	DeleteDocuments(docIds []*uint, tx ...utility.DatabaseTransaction) error
}

type ComplianceReviewRepository struct {
	db *gorm.DB
}

func NewComplianceReviewRepository(db *gorm.DB) *ComplianceReviewRepository {
	return &ComplianceReviewRepository{
		db: db,
	}
}

// BeginTransaction starts a new transaction.
func (repo *ComplianceReviewRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
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
			(compliance_register.title_en LIKE ? OR
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

	if result.Error != nil {
		return nil, result.Error
	}

	// Fetch documents separately and aggregate them into the complianceReviews slice
	// for _, review := range complianceReviews {
	// 	var documents []*complianceDTOs.DocumentDTO
	// 	docQuery := repo.db.Table("compliance_review_documents").
	// 		Select(`
	// 			url,
	// 			name,
	// 			size,
	// 			token,
	// 			type,
	// 			original_name
	// 		`).
	// 		Where("compliance_review_id = ?", review.ID).
	// 		Find(&documents)

	// 	if docQuery.Error != nil {
	// 		return nil, docQuery.Error
	// 	}

	// 	review.Documents = documents
	// }

	return complianceReviews, nil
}


func (repo *ComplianceReviewRepository) GetByID(id uint) (*models.ComplianceReview, error) {
	var complianceReview *models.ComplianceReview
	result := repo.db.First(&complianceReview, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return complianceReview, nil
}

func (repo *ComplianceReviewRepository) Create(complianceReview *models.ComplianceReview, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(complianceReview).Error; err != nil {
		return 0, err
	}

	return complianceReview.ID, nil
}

func (repo *ComplianceReviewRepository) Update(ComplianceReview *models.ComplianceReview, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(ComplianceReview).Error; err != nil {
		return err
	}
	return nil
}

func (repo *ComplianceReviewRepository) CreateDocuments(docs []*models.ComplianceReviewDocument, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	// Save all documents in a single transaction
	if err := db.Create(&docs).Error; err != nil {
		return err
	}
	return nil
}

func (repo *ComplianceReviewRepository) GetDocumentByID(id uint) (*models.ComplianceReviewDocument, error) {
	var document models.ComplianceReviewDocument
	if err := repo.db.Where("id = ?", id).First(&document).Error; err != nil {
		return nil, err
	}
	return &document, nil
}

func (repo *ComplianceReviewRepository) DeleteDocuments(docIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	// Delete documents by IDs
	if err := db.Where("id IN (?)", docIds).Delete(&models.ComplianceReviewDocument{}).Error; err != nil {
		return err
	}

	return nil
}
