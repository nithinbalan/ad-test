package repositories

import (
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type IdentityReviewUserCommentRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(identityReviewUserID uint) ([]*models.IdentityReviewUserComment, error)
	GetByID(identityReviewUserID uint, id uint) (*models.IdentityReviewUserComment, error)
	Create(comment *models.IdentityReviewUserComment, tx ...utility.DatabaseTransaction) (uint, error)
	Update(comment *models.IdentityReviewUserComment, tx ...utility.DatabaseTransaction) error
	Delete(identityReviewUserID uint, ids []uint, tx ...utility.DatabaseTransaction) error
}

type IdentityReviewUserCommentRepository struct {
	db *gorm.DB
}

func NewIdentityReviewUserCommentRepository(db *gorm.DB) *IdentityReviewUserCommentRepository {
	return &IdentityReviewUserCommentRepository{
		db: db,
	}
}

func (repo *IdentityReviewUserCommentRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *IdentityReviewUserCommentRepository) GetAll(identityReviewUserID uint) ([]*models.IdentityReviewUserComment, error) {
	var comments []*models.IdentityReviewUserComment

	err := repo.db.
		Preload("AccountDecision").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		Where("identity_review_user_comments.identity_review_user_id = ?", identityReviewUserID).
		Order("identity_review_user_comments.created_at DESC").
		Find(&comments).Error
	return comments, err
}

func (repo *IdentityReviewUserCommentRepository) GetByID(identityReviewUserID uint, id uint) (*models.IdentityReviewUserComment, error) {
	var comment models.IdentityReviewUserComment

	err := repo.db.
		Preload("AccountDecision").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		Where("identity_review_user_comments.identity_review_user_id = ? AND identity_review_user_comments.id = ?", identityReviewUserID, id).
		First(&comment).Error
	if err != nil {
		return nil, err
	}

	return &comment, nil
}

func (repo *IdentityReviewUserCommentRepository) Create(comment *models.IdentityReviewUserComment, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Create(comment).Error; err != nil {
		return 0, err
	}

	return comment.ID, nil
}

func (repo *IdentityReviewUserCommentRepository) Update(comment *models.IdentityReviewUserComment, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Model(comment).Updates(map[string]interface{}{
		"account_decision_id": comment.AccountDecisionID,
		"comment":             comment.Comment,
		"updated_by_id":       comment.UpdatedBYID,
		"updated_at":          comment.UpdatedAt,
	}).Error
}

func (repo *IdentityReviewUserCommentRepository) Delete(identityReviewUserID uint, ids []uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.
		Where("identity_review_user_id = ? AND id IN (?)", identityReviewUserID, ids).
		Delete(&models.IdentityReviewUserComment{}).Error
}
