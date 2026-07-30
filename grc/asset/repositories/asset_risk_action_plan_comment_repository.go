package repositories

import (
	assetDTOs "jethur/grc/asset/dtos"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type AssetRiskActionPlanCommentRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(actionPlanID uint) ([]*assetDTOs.AssetRiskActionPlanCommentListDTO, error)
	GetByID(id uint) (*models.RiskActionPlanComment, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.RiskActionPlanComment, error)
	Create(comment *models.RiskActionPlanComment, tx ...utility.DatabaseTransaction) (uint, error)
	Update(comment *models.RiskActionPlanComment, tx ...utility.DatabaseTransaction) error
	Delete(actionPlanID uint, tx ...utility.DatabaseTransaction) error

	CreateDocuments(docs []*models.RiskActionPlanCommentDocument, tx ...utility.DatabaseTransaction) error
	DeleteDocuments(docIds []*uint, tx ...utility.DatabaseTransaction) error
	GetDocumentByID(id uint) (*models.RiskActionPlanCommentDocument, error)
}

type AssetRiskActionPlanCommentRepository struct {
	db *gorm.DB
}

func NewAssetRiskActionPlanCommentRepository(db *gorm.DB) AssetRiskActionPlanCommentRepositoryInterface {
	return &AssetRiskActionPlanCommentRepository{
		db: db,
	}
}

func (repo *AssetRiskActionPlanCommentRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *AssetRiskActionPlanCommentRepository) GetAll(actionPlanID uint) ([]*assetDTOs.AssetRiskActionPlanCommentListDTO, error) {
	var comments []*assetDTOs.AssetRiskActionPlanCommentListDTO
	query := repo.db.Table("risk_action_plan_comments").
		Select(`
	risk_action_plan_comments.id,
	risk_action_plan_comments.comment,
	risk_action_plan_comments.risk_action_plan_id,
	risk_action_plan_comments.created_at,
	risk_action_plan_comments.created_by_id,
	users.first_name_en AS created_user_first_name_en,
	users.first_name_ar AS created_user_first_name_ar,
	users.last_name_en AS created_user_last_name_en,
	users.last_name_ar AS created_user_last_name_ar,
	users.initial AS created_user_initial,
	users.color AS created_user_color,
	users.email AS created_user_email,
	users.profile_image_path AS created_user_profile_image_path
`).
		Joins("LEFT JOIN users ON users.id = risk_action_plan_comments.created_by_id").
		Where("risk_action_plan_comments.risk_action_plan_id = ?", actionPlanID)

	if err := query.Find(&comments).Error; err != nil {
		return nil, err
	}

	for _, comment := range comments {
		var documents []*assetDTOs.AssetRiskActionPlanCommentDocumentDTO
		repo.db.Model(&models.RiskActionPlanCommentDocument{}).
			Where("risk_action_plan_comment_id = ?", comment.ID).
			Find(&documents)
		comment.Documents = documents
	}

	return comments, nil
}

func (repo *AssetRiskActionPlanCommentRepository) GetByID(id uint) (*models.RiskActionPlanComment, error) {
	var riskActionPlanComment *models.RiskActionPlanComment
	result := repo.db.Preload("Documents").First(&riskActionPlanComment, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return riskActionPlanComment, nil
}

func (repo *AssetRiskActionPlanCommentRepository) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.RiskActionPlanComment, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var comment models.RiskActionPlanComment
	result := db.
		Preload("RiskActionPlan").
		Preload("CreatedByUser").
		Preload("Documents").
		Where("id = ?", id).
		First(&comment)

	if result.Error != nil {
		return nil, result.Error
	}

	return &comment, nil
}

func (repo *AssetRiskActionPlanCommentRepository) Create(comment *models.RiskActionPlanComment, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(comment).Error; err != nil {
		return 0, err
	}

	return comment.ID, nil
}

func (repo *AssetRiskActionPlanCommentRepository) Update(comment *models.RiskActionPlanComment, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Model(&models.RiskActionPlanComment{}).Where("id = ?", comment.ID).Updates(comment).Error; err != nil {
		return err
	}

	return nil
}

func (repo *AssetRiskActionPlanCommentRepository) Delete(actionPlanID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("id = ?", actionPlanID).Delete(&models.RiskActionPlanComment{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *AssetRiskActionPlanCommentRepository) CreateDocuments(docs []*models.RiskActionPlanCommentDocument, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Create(&docs).Error; err != nil {
		return err
	}
	return nil
}

func (repo *AssetRiskActionPlanCommentRepository) DeleteDocuments(docIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("id IN (?)", docIds).Delete(&models.RiskActionPlanCommentDocument{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *AssetRiskActionPlanCommentRepository) GetDocumentByID(id uint) (*models.RiskActionPlanCommentDocument, error) {
	var document models.RiskActionPlanCommentDocument
	if err := repo.db.Where("id = ?", id).First(&document).Error; err != nil {
		return nil, err
	}
	return &document, nil
}
