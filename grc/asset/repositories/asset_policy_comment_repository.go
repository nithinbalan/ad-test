package repositories

import (
	"jethur/grc/asset/dtos"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type AssetPolicyCommentRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(assetPolicyID uint) ([]*dtos.AssetPolicyCommentListDTO, error)
	GetByID(id uint) (*models.AssetPolicyComment, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.AssetPolicyComment, error)
	Create(comment *models.AssetPolicyComment, tx ...utility.DatabaseTransaction) (uint, error)
	Update(comment *models.AssetPolicyComment, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error
}

type AssetPolicyCommentRepository struct {
	db *gorm.DB
}

func NewAssetPolicyCommentRepository(db *gorm.DB) AssetPolicyCommentRepositoryInterface {
	return &AssetPolicyCommentRepository{db: db}
}

func (repo *AssetPolicyCommentRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *AssetPolicyCommentRepository) GetAll(assetPolicyID uint) ([]*dtos.AssetPolicyCommentListDTO, error) {
	var comments []*dtos.AssetPolicyCommentListDTO

	query := repo.db.Table("asset_policy_comments").
		Select(`
			asset_policy_comments.id,
			asset_policy_comments.comments,
			asset_policy_comments.asset_policy_id,
			asset_policy_comments.created_at,
			asset_policy_comments.updated_at,
			asset_policy_comments.created_by_id,
			asset_policy_comments.updated_by_id,
			created_users.first_name_en AS created_user_first_name_en,
			created_users.first_name_ar AS created_user_first_name_ar,
			created_users.last_name_en AS created_user_last_name_en,
			created_users.last_name_ar AS created_user_last_name_ar,
			created_users.initial AS created_user_initial,
			created_users.color AS created_user_color,
			created_users.email AS created_user_email,
			created_users.profile_image_path AS created_user_profile_image,
			updated_users.first_name_en AS updated_user_first_name_en,
			updated_users.first_name_ar AS updated_user_first_name_ar,
			updated_users.last_name_en AS updated_user_last_name_en,
			updated_users.last_name_ar AS updated_user_last_name_ar,
			updated_users.initial AS updated_user_initial,
			updated_users.color AS updated_user_color,
			updated_users.email AS updated_user_email,
			updated_users.profile_image_path AS updated_user_profile_image
		`).
		Joins("LEFT JOIN users AS created_users ON created_users.id = asset_policy_comments.created_by_id").
		Joins("LEFT JOIN users AS updated_users ON updated_users.id = asset_policy_comments.updated_by_id").
		Where("asset_policy_comments.asset_policy_id = ?", assetPolicyID)

	if err := query.Find(&comments).Error; err != nil {
		return nil, err
	}

	return comments, nil
}

func (repo *AssetPolicyCommentRepository) GetByID(id uint) (*models.AssetPolicyComment, error) {
	var comment models.AssetPolicyComment
	if err := repo.db.First(&comment, id).Error; err != nil {
		return nil, err
	}
	return &comment, nil
}

func (repo *AssetPolicyCommentRepository) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.AssetPolicyComment, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var comment models.AssetPolicyComment
	result := db.
		Preload("AssetPolicy").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		Where("id = ?", id).
		First(&comment)

	if result.Error != nil {
		return nil, result.Error
	}

	return &comment, nil
}

func (repo *AssetPolicyCommentRepository) Create(comment *models.AssetPolicyComment, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if err := db.Create(comment).Error; err != nil {
		return 0, err
	}
	return comment.ID, nil
}

func (repo *AssetPolicyCommentRepository) Update(comment *models.AssetPolicyComment, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	return db.Save(comment).Error
}

func (repo *AssetPolicyCommentRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	return db.Delete(&models.AssetPolicyComment{}, id).Error
}
