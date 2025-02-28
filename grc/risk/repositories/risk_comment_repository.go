package repositories

import (
	"jethur/grc/models"
	"jethur/grc/risk/dtos"
	"jethur/utility"

	"gorm.io/gorm"
)

type RiskCommentRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction

	GetAll(riskID uint) ([]*dtos.RiskCommentListDTO, error)
	GetByID(id uint) (*models.RiskComment, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.RiskComment, error)
	Create(comment *models.RiskComment, tx ...utility.DatabaseTransaction) (uint, error)
	Update(comment *models.RiskComment, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error

	UpdateRiskStatus(risk *models.RiskRegister, tx ...utility.DatabaseTransaction) error
}

type RiskCommentRepository struct {
	db *gorm.DB
}

func NewRiskCommentRepository(db *gorm.DB) RiskCommentRepositoryInterface {
	return &RiskCommentRepository{
		db: db,
	}
}

func (repo *RiskCommentRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

// Get all comments for a specific risk
func (repo *RiskCommentRepository) GetAll(riskID uint) ([]*dtos.RiskCommentListDTO, error) {
	var comments []*dtos.RiskCommentListDTO
	query := repo.db.Table("risk_comments").
		Select(`
			risk_comments.id,
			risk_comments.comment,
			risk_comments.risk_register_id,
			risk_comments.created_at,
			risk_comments.created_by_id,
			users.first_name_en AS created_user_first_name_en,
			users.first_name_ar AS created_user_first_name_ar,
			users.last_name_en AS created_user_last_name_en,
			users.last_name_ar AS created_user_last_name_ar,
			users.initial AS created_user_initial,
			users.color AS created_user_color,
			users.email AS created_user_email,
			users.profile_image_path AS created_user_profile_image_path
		`).
		Joins("LEFT JOIN users ON users.id = risk_comments.created_by_id").
		Where("risk_comments.risk_register_id = ?", riskID)

	if err := query.Find(&comments).Error; err != nil {
		return nil, err
	}

	return comments, nil
}

// Get a specific comment by ID
func (repo *RiskCommentRepository) GetByID(id uint) (*models.RiskComment, error) {
	var riskComment *models.RiskComment
	result := repo.db.First(&riskComment, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return riskComment, nil
}

// Get detailed information about a specific comment
func (repo *RiskCommentRepository) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.RiskComment, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var comment models.RiskComment
	result := db.
		Preload("RiskRegister").
		Preload("CreatedByUser").
		Where("id = ?", id).
		First(&comment)

	if result.Error != nil {
		return nil, result.Error
	}

	return &comment, nil
}

// Create a new risk comment
func (repo *RiskCommentRepository) Create(comment *models.RiskComment, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(comment).Error; err != nil {
		return 0, err
	}

	return comment.ID, nil
}

// Update an existing risk comment
func (repo *RiskCommentRepository) Update(comment *models.RiskComment, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Model(&models.RiskComment{}).Where("id = ?", comment.ID).Updates(comment).Error; err != nil {
		return err
	}

	return nil
}

// Delete a risk comment by ID
func (repo *RiskCommentRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("id = ?", id).Delete(&models.RiskComment{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *RiskCommentRepository) UpdateRiskStatus(risk *models.RiskRegister, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Model(risk).Update("risk_status_id", risk.RiskStatusID).Error; err != nil {
		return err
	}

	return nil
}
