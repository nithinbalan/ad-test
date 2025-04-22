package repositories

import (
	"jethur/grc/models"
	strategyDTOs "jethur/grc/strategy/dtos"
	"jethur/utility"

	"gorm.io/gorm"
)

type StrategyProfileActionPlanCommentRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(actionPlanID uint) ([]*strategyDTOs.StrategyProfileActionPlanCommentListDTO, error)
	GetByID(id uint) (*models.StrategyProfileActionPlanComment, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.StrategyProfileActionPlanComment, error)
	Create(comment *models.StrategyProfileActionPlanComment, tx ...utility.DatabaseTransaction) (uint, error)
	Update(comment *models.StrategyProfileActionPlanComment, tx ...utility.DatabaseTransaction) error
	Delete(actionPlanID uint, tx ...utility.DatabaseTransaction) error

	CreateDocuments(docs []*models.StrategyProfileActionPlanCommentDocument, tx ...utility.DatabaseTransaction) error
	DeleteDocuments(docIds []*uint, tx ...utility.DatabaseTransaction) error
	GetDocumentByID(id uint) (*models.StrategyProfileActionPlanCommentDocument, error)
}

type StrategyProfileActionPlanCommentRepository struct {
	db *gorm.DB
}

func NewStrategyProfileActionPlanCommentRepository(db *gorm.DB) StrategyProfileActionPlanCommentRepositoryInterface {
	return &StrategyProfileActionPlanCommentRepository{
		db: db,
	}
}

func (repo *StrategyProfileActionPlanCommentRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *StrategyProfileActionPlanCommentRepository) GetAll(actionPlanID uint) ([]*strategyDTOs.StrategyProfileActionPlanCommentListDTO, error) {
	var comments []*strategyDTOs.StrategyProfileActionPlanCommentListDTO
	query := repo.db.Table("strategy_profile_action_plan_comments").
		Select(`
	strategy_profile_action_plan_comments.id,
	strategy_profile_action_plan_comments.comment,
	strategy_profile_action_plan_comments.strategy_profile_action_plan_id,
	strategy_profile_action_plan_comments.created_at,
	strategy_profile_action_plan_comments.created_by_id,
	users.first_name_en AS created_user_first_name_en,
	users.first_name_ar AS created_user_first_name_ar,
	users.last_name_en AS created_user_last_name_en,
	users.last_name_ar AS created_user_last_name_ar,
	users.initial AS created_user_initial,
	users.color AS created_user_color,
	users.email AS created_user_email,
	users.profile_image_path AS created_user_profile_image_path
`).
		Joins("LEFT JOIN users ON users.id = strategy_profile_action_plan_comments.created_by_id").
		Where("strategy_profile_action_plan_comments.strategy_profile_action_plan_id = ?", actionPlanID)

	if err := query.Find(&comments).Error; err != nil {
		return nil, err
	}

	for _, comment := range comments {
		var documents []*strategyDTOs.StrategyProfileActionPlanCommentDocumentDTO
		repo.db.Model(&models.StrategyProfileActionPlanCommentDocument{}).
			Where("strategy_profile_action_plan_comment_id = ?", comment.ID).
			Find(&documents)
		comment.Documents = documents
	}

	return comments, nil
}

func (repo *StrategyProfileActionPlanCommentRepository) GetByID(id uint) (*models.StrategyProfileActionPlanComment, error) {
	var strategyProfileActionPlanComment *models.StrategyProfileActionPlanComment
	result := repo.db.Preload("Documents").First(&strategyProfileActionPlanComment, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return strategyProfileActionPlanComment, nil
}

func (repo *StrategyProfileActionPlanCommentRepository) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.StrategyProfileActionPlanComment, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var comment models.StrategyProfileActionPlanComment
	result := db.
		Preload("StrategyProfileActionPlan").
		Preload("CreatedByUser").
		Preload("Documents").
		Where("id = ?", id).
		First(&comment)

	if result.Error != nil {
		return nil, result.Error
	}

	return &comment, nil
}

func (repo *StrategyProfileActionPlanCommentRepository) Create(comment *models.StrategyProfileActionPlanComment, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(comment).Error; err != nil {
		return 0, err
	}

	return comment.ID, nil
}

func (repo *StrategyProfileActionPlanCommentRepository) Update(comment *models.StrategyProfileActionPlanComment, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Model(&models.StrategyProfileActionPlanComment{}).Where("id = ?", comment.ID).Updates(comment).Error; err != nil {
		return err
	}

	return nil
}

func (repo *StrategyProfileActionPlanCommentRepository) Delete(actionPlanID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("id = ?", actionPlanID).Delete(&models.StrategyProfileActionPlanComment{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *StrategyProfileActionPlanCommentRepository) CreateDocuments(docs []*models.StrategyProfileActionPlanCommentDocument, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Create(&docs).Error; err != nil {
		return err
	}
	return nil
}

func (repo *StrategyProfileActionPlanCommentRepository) DeleteDocuments(docIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("id IN (?)", docIds).Delete(&models.StrategyProfileActionPlanCommentDocument{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *StrategyProfileActionPlanCommentRepository) GetDocumentByID(id uint) (*models.StrategyProfileActionPlanCommentDocument, error) {
	var document models.StrategyProfileActionPlanCommentDocument
	if err := repo.db.Where("id = ?", id).First(&document).Error; err != nil {
		return nil, err
	}
	return &document, nil
}
