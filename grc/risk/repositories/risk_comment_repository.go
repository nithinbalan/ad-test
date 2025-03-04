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

func (repo *RiskCommentRepository) GetAll(riskID uint) ([]*dtos.RiskCommentListDTO, error) {
	var comments []*dtos.RiskCommentListDTO
	query := repo.db.Table("risk_comments").
		Select(`
			risk_comments.id,
			risk_comments.comments,
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
			users.profile_image_path AS created_user_profile_image_path,
			risk_statuses.title_en AS risk_status_title_en,
			risk_statuses.title_ar AS risk_status_title_ar,
			risk_statuses.color_code AS risk_status_color_code,
			risk_timeline_types.title_en AS risk_timeline_title_en,
			risk_timeline_types.title_ar AS risk_timeline_title_ar,
			risk_timeline_types.color_code AS risk_timeline_color_code,
			risk_timeline_types.type AS risk_timeline_type
		`).
		Joins("LEFT JOIN users ON users.id = risk_comments.created_by_id").
		Joins("LEFT JOIN risk_statuses ON risk_statuses.id = risk_comments.risk_status_id").
		Joins("LEFT JOIN risk_timeline_types ON risk_timeline_types.id = risk_comments.risk_timeline_type_id").
		Where("risk_comments.risk_register_id = ?", riskID)

	if err := query.Find(&comments).Error; err != nil {
		return nil, err
	}

	return comments, nil
}

func (repo *RiskCommentRepository) GetByID(id uint) (*models.RiskComment, error) {
	var riskComment *models.RiskComment
	result := repo.db.First(&riskComment, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return riskComment, nil
}

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

	if err := db.Model(&models.RiskRegister{}).Where("id = ?", risk.ID).Updates(risk).Error; err != nil {
		return err
	}

	return nil
}
