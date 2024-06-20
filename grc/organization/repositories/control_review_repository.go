package repositories

import (
	"jethur/grc/models"
	controlDTOs "jethur/grc/organization/dtos"
	"jethur/utility"

	"gorm.io/gorm"
)

type ControlReviewRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(controlID uint, sort string, sortBy string, keyword string) ([]*controlDTOs.ControlReviewListDTO, error)
	GetByID(id uint) (*models.ControlReview, error)
	Create(ControlReview *models.ControlReview, tx ...utility.DatabaseTransaction) (uint, error)
	Update(ControlReview *models.ControlReview, tx ...utility.DatabaseTransaction) error

	//Document
	CreateDocuments(docs []*models.ControlReviewDocument, tx ...utility.DatabaseTransaction) error
	GetDocumentByID(id uint) (*models.ControlReviewDocument, error)
	DeleteDocuments(docIds []*uint, tx ...utility.DatabaseTransaction) error
}

type ControlReviewRepository struct {
	db *gorm.DB
}

func NewControlReviewRepository(db *gorm.DB) *ControlReviewRepository {
	return &ControlReviewRepository{
		db: db,
	}
}

// BeginTransaction starts a new transaction.
func (repo *ControlReviewRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *ControlReviewRepository) GetAll(controlID uint, sort string, sortBy string, keyword string) ([]*controlDTOs.ControlReviewListDTO, error) {
	var controlReviews []*controlDTOs.ControlReviewListDTO

	query := repo.db.Table("control_reviews").
		Select(`
            control_reviews.*,
            controls.title_en AS control_title_en,
            controls.title_ar AS control_title_ar,
            controls.description AS control_description,
            controls.code AS control_code,
            reviewed_user.first_name_en AS reviewed_user_first_name_en,
            reviewed_user.first_name_ar AS reviewed_user_first_name_ar,
            reviewed_user.last_name_en AS reviewed_user_last_name_en,
            reviewed_user.last_name_ar AS reviewed_user_last_name_ar,
			reviewed_user.initial AS reviewed_user_initial,
        	reviewed_user.color AS reviewed_user_color,
        	reviewed_user.email AS reviewed_user_email,
        	reviewed_user.profile_image_path AS reviewed_user_profile_image_path 
        `).
		Joins(`
            LEFT JOIN controls ON controls.id = control_reviews.control_id
            LEFT JOIN users AS reviewed_user ON reviewed_user.id = control_reviews.reviewed_by
        `).
		Where("control_reviews.control_id = ?", controlID)

	if keyword != "" {
		query = query.Where(`
            (controls.title_en LIKE ? OR
            control_reviews.comment LIKE ?)`,
			"%"+keyword+"%", "%"+keyword+"%")
	}
	if sortBy == "" {
		sortBy = "control_reviews.id"
	}
	if sort != "asc" && sort != "desc" {
		sort = "asc"
	}
	query = query.Order(sortBy + " " + sort)
	result := query.Find(&controlReviews)

	if result.Error != nil {
		return nil, result.Error
	}

	// Fetching documents for each control review
	for _, review := range controlReviews {
		var documents []*models.ControlReviewDocument
		repo.db.Model(&models.ControlReviewDocument{}).Where("control_review_id = ?", review.ID).Find(&documents)
		review.Documents = documents
	}

	return controlReviews, nil
}

func (repo *ControlReviewRepository) GetByID(id uint) (*models.ControlReview, error) {
	var controlReview *models.ControlReview
	result := repo.db.First(&controlReview, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return controlReview, nil
}

func (repo *ControlReviewRepository) Create(controlReview *models.ControlReview, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(controlReview).Error; err != nil {
		return 0, err
	}

	return controlReview.ID, nil
}

func (repo *ControlReviewRepository) Update(ControlReview *models.ControlReview, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(ControlReview).Error; err != nil {
		return err
	}
	return nil
}

func (repo *ControlReviewRepository) CreateDocuments(docs []*models.ControlReviewDocument, tx ...utility.DatabaseTransaction) error {
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

func (repo *ControlReviewRepository) GetDocumentByID(id uint) (*models.ControlReviewDocument, error) {
	var document models.ControlReviewDocument
	if err := repo.db.Where("id = ?", id).First(&document).Error; err != nil {
		return nil, err
	}
	return &document, nil
}

func (repo *ControlReviewRepository) DeleteDocuments(docIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	// Delete documents by IDs
	if err := db.Where("id IN (?)", docIds).Delete(&models.ControlReviewDocument{}).Error; err != nil {
		return err
	}

	return nil
}
