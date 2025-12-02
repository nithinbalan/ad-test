package repositories

import (
	"jethur/grc/models"
	controlDTOs "jethur/grc/organization/dtos"
	"jethur/utility"

	"gorm.io/gorm"
)

type ControlActionPlanCommentRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(actionPlanID uint) ([]*controlDTOs.ControlActionPlanCommentListDTO, error)
	GetByID(id uint) (*models.ControlActionPlanComment, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.ControlActionPlanComment, error)
	Create(comment *models.ControlActionPlanComment, tx ...utility.DatabaseTransaction) (uint, error)
	Update(comment *models.ControlActionPlanComment, tx ...utility.DatabaseTransaction) error
	Delete(actionPlanID uint, tx ...utility.DatabaseTransaction) error

	CreateDocuments(docs []*models.ControlActionPlanCommentDocument, tx ...utility.DatabaseTransaction) error
	DeleteDocuments(docIds []*uint, tx ...utility.DatabaseTransaction) error
	GetDocumentByID(id uint) (*models.ControlActionPlanCommentDocument, error)
}

type ControlActionPlanCommentRepository struct {
	db *gorm.DB
}

func NewControlActionPlanCommentRepository(db *gorm.DB) ControlActionPlanCommentRepositoryInterface {
	return &ControlActionPlanCommentRepository{
		db: db,
	}
}

func (repo *ControlActionPlanCommentRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *ControlActionPlanCommentRepository) GetAll(actionPlanID uint) ([]*controlDTOs.ControlActionPlanCommentListDTO, error) {
	var comments []*controlDTOs.ControlActionPlanCommentListDTO
	query := repo.db.Table("control_action_plan_comments").
		Select(`
	control_action_plan_comments.id,
	control_action_plan_comments.comment,
	control_action_plan_comments.control_action_plan_id,
	control_action_plan_comments.created_at,
	control_action_plan_comments.created_by_id,
	users.first_name_en AS created_user_first_name_en,
	users.first_name_ar AS created_user_first_name_ar,
	users.last_name_en AS created_user_last_name_en,
	users.last_name_ar AS created_user_last_name_ar,
	users.initial AS created_user_initial,
	users.color AS created_user_color,
	users.email AS created_user_email,
	users.profile_image_path AS created_user_profile_image_path
`).
		Joins("LEFT JOIN users ON users.id = control_action_plan_comments.created_by_id").
		Where("control_action_plan_comments.control_action_plan_id = ?", actionPlanID)

	if err := query.Find(&comments).Error; err != nil {
		return nil, err
	}

	for _, comment := range comments {
		var documents []*controlDTOs.ControlActionPlanCommentDocumentDTO
		repo.db.Model(&models.ControlActionPlanCommentDocument{}).
			Where("control_action_plan_comment_id = ?", comment.ID).
			Find(&documents)
		comment.Documents = documents
	}

	return comments, nil
}

func (repo *ControlActionPlanCommentRepository) GetByID(id uint) (*models.ControlActionPlanComment, error) {
	var controlActionPlanComment *models.ControlActionPlanComment
	result := repo.db.Preload("Documents").First(&controlActionPlanComment, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return controlActionPlanComment, nil
}

func (repo *ControlActionPlanCommentRepository) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.ControlActionPlanComment, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var comment models.ControlActionPlanComment
	result := db.
		Preload("ControlActionPlan").
		Preload("CreatedByUser").
		Preload("Documents").
		Where("id = ?", id).
		First(&comment)

	if result.Error != nil {
		return nil, result.Error
	}

	return &comment, nil
}

func (repo *ControlActionPlanCommentRepository) Create(comment *models.ControlActionPlanComment, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(comment).Error; err != nil {
		return 0, err
	}

	return comment.ID, nil
}

func (repo *ControlActionPlanCommentRepository) Update(comment *models.ControlActionPlanComment, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Model(&models.ControlActionPlanComment{}).Where("id = ?", comment.ID).Updates(comment).Error; err != nil {
		return err
	}

	return nil
}

func (repo *ControlActionPlanCommentRepository) Delete(actionPlanID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("id = ?", actionPlanID).Delete(&models.ControlActionPlanComment{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *ControlActionPlanCommentRepository) CreateDocuments(docs []*models.ControlActionPlanCommentDocument, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Create(&docs).Error; err != nil {
		return err
	}
	return nil
}

func (repo *ControlActionPlanCommentRepository) DeleteDocuments(docIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("id IN (?)", docIds).Delete(&models.ControlActionPlanCommentDocument{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *ControlActionPlanCommentRepository) GetDocumentByID(id uint) (*models.ControlActionPlanCommentDocument, error) {
	var document models.ControlActionPlanCommentDocument
	if err := repo.db.Where("id = ?", id).First(&document).Error; err != nil {
		return nil, err
	}
	return &document, nil
}
