package repositories

import (
	threatDTOs "jethur/grc/asset/dtos"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type ThreatActionPlanCommentRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(actionPlanID uint) ([]*threatDTOs.ThreatActionPlanCommentListDTO, error)
	GetByID(id uint) (*models.ThreatActionPlanComment, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.ThreatActionPlanComment, error)
	Create(comment *models.ThreatActionPlanComment, tx ...utility.DatabaseTransaction) (uint, error)
	Update(comment *models.ThreatActionPlanComment, tx ...utility.DatabaseTransaction) error
	Delete(actionPlanID uint, tx ...utility.DatabaseTransaction) error

	CreateDocuments(docs []*models.ThreatActionPlanCommentDocument, tx ...utility.DatabaseTransaction) error
	DeleteDocuments(docIds []*uint, tx ...utility.DatabaseTransaction) error
	GetDocumentByID(id uint) (*models.ThreatActionPlanCommentDocument, error)
}

type ThreatActionPlanCommentRepository struct {
	db *gorm.DB
}

func NewThreatActionPlanCommentRepository(db *gorm.DB) ThreatActionPlanCommentRepositoryInterface {
	return &ThreatActionPlanCommentRepository{
		db: db,
	}
}

func (repo *ThreatActionPlanCommentRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *ThreatActionPlanCommentRepository) GetAll(actionPlanID uint) ([]*threatDTOs.ThreatActionPlanCommentListDTO, error) {
	var comments []*threatDTOs.ThreatActionPlanCommentListDTO
	query := repo.db.Table("threat_action_plan_comments").
		Select(`
	threat_action_plan_comments.id,
	threat_action_plan_comments.comment,
	threat_action_plan_comments.threat_action_plan_id,
	threat_action_plan_comments.created_at,
	threat_action_plan_comments.created_by_id,
	users.first_name_en AS created_user_first_name_en,
	users.first_name_ar AS created_user_first_name_ar,
	users.last_name_en AS created_user_last_name_en,
	users.last_name_ar AS created_user_last_name_ar,
	users.initial AS created_user_initial,
	users.color AS created_user_color,
	users.email AS created_user_email,
	users.profile_image_path AS created_user_profile_image_path
`).
		Joins("LEFT JOIN users ON users.id = threat_action_plan_comments.created_by_id").
		Where("threat_action_plan_comments.threat_action_plan_id = ?", actionPlanID)

	if err := query.Find(&comments).Error; err != nil {
		return nil, err
	}

	for _, comment := range comments {
		var documents []*threatDTOs.ThreatActionPlanCommentDocumentDTO
		repo.db.Model(&models.ThreatActionPlanCommentDocument{}).
			Where("threat_action_plan_comment_id = ?", comment.ID).
			Find(&documents)
		comment.Documents = documents
	}

	return comments, nil
}

func (repo *ThreatActionPlanCommentRepository) GetByID(id uint) (*models.ThreatActionPlanComment, error) {
	var threatActionPlanComment *models.ThreatActionPlanComment
	result := repo.db.Preload("Documents").First(&threatActionPlanComment, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return threatActionPlanComment, nil
}

func (repo *ThreatActionPlanCommentRepository) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.ThreatActionPlanComment, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var comment models.ThreatActionPlanComment
	result := db.
		Preload("ThreatActionPlan").
		Preload("CreatedByUser").
		Preload("Documents").
		Where("id = ?", id).
		First(&comment)

	if result.Error != nil {
		return nil, result.Error
	}

	return &comment, nil
}

func (repo *ThreatActionPlanCommentRepository) Create(comment *models.ThreatActionPlanComment, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(comment).Error; err != nil {
		return 0, err
	}

	return comment.ID, nil
}

func (repo *ThreatActionPlanCommentRepository) Update(comment *models.ThreatActionPlanComment, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Model(&models.ThreatActionPlanComment{}).Where("id = ?", comment.ID).Updates(comment).Error; err != nil {
		return err
	}

	return nil
}

func (repo *ThreatActionPlanCommentRepository) Delete(actionPlanID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("id = ?", actionPlanID).Delete(&models.ThreatActionPlanComment{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *ThreatActionPlanCommentRepository) CreateDocuments(docs []*models.ThreatActionPlanCommentDocument, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Create(&docs).Error; err != nil {
		return err
	}
	return nil
}

func (repo *ThreatActionPlanCommentRepository) DeleteDocuments(docIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("id IN (?)", docIds).Delete(&models.ThreatActionPlanCommentDocument{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *ThreatActionPlanCommentRepository) GetDocumentByID(id uint) (*models.ThreatActionPlanCommentDocument, error) {
	var document models.ThreatActionPlanCommentDocument
	if err := repo.db.Where("id = ?", id).First(&document).Error; err != nil {
		return nil, err
	}
	return &document, nil
}
