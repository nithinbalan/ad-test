package repositories

import (
	vulnerabilityDTOs "jethur/grc/asset/dtos"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type VulnerabilityActionPlanCommentRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(actionPlanID uint) ([]*vulnerabilityDTOs.VulnerabilityActionPlanCommentListDTO, error)
	GetByID(id uint) (*models.VulnerabilityActionPlanComment, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.VulnerabilityActionPlanComment, error)
	Create(comment *models.VulnerabilityActionPlanComment, tx ...utility.DatabaseTransaction) (uint, error)
	Update(comment *models.VulnerabilityActionPlanComment, tx ...utility.DatabaseTransaction) error
	Delete(actionPlanID uint, tx ...utility.DatabaseTransaction) error

	CreateDocuments(docs []*models.VulnerabilityActionPlanCommentDocument, tx ...utility.DatabaseTransaction) error
	DeleteDocuments(docIds []*uint, tx ...utility.DatabaseTransaction) error
	GetDocumentByID(id uint) (*models.VulnerabilityActionPlanCommentDocument, error)
}

type VulnerabilityActionPlanCommentRepository struct {
	db *gorm.DB
}

func NewVulnerabilityActionPlanCommentRepository(db *gorm.DB) VulnerabilityActionPlanCommentRepositoryInterface {
	return &VulnerabilityActionPlanCommentRepository{
		db: db,
	}
}

func (repo *VulnerabilityActionPlanCommentRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *VulnerabilityActionPlanCommentRepository) GetAll(actionPlanID uint) ([]*vulnerabilityDTOs.VulnerabilityActionPlanCommentListDTO, error) {
	var comments []*vulnerabilityDTOs.VulnerabilityActionPlanCommentListDTO
	query := repo.db.Table("vulnerability_action_plan_comments").
		Select(`
	vulnerability_action_plan_comments.id,
	vulnerability_action_plan_comments.comment,
	vulnerability_action_plan_comments.vulnerability_action_plan_id,
	vulnerability_action_plan_comments.created_at,
	vulnerability_action_plan_comments.created_by_id,
	users.first_name_en AS created_user_first_name_en,
	users.first_name_ar AS created_user_first_name_ar,
	users.last_name_en AS created_user_last_name_en,
	users.last_name_ar AS created_user_last_name_ar,
	users.initial AS created_user_initial,
	users.color AS created_user_color,
	users.email AS created_user_email,
	users.profile_image_path AS created_user_profile_image_path
`).
		Joins("LEFT JOIN users ON users.id = vulnerability_action_plan_comments.created_by_id").
		Where("vulnerability_action_plan_comments.vulnerability_action_plan_id = ?", actionPlanID)

	if err := query.Find(&comments).Error; err != nil {
		return nil, err
	}

	for _, comment := range comments {
		var documents []*vulnerabilityDTOs.VulnerabilityActionPlanCommentDocumentDTO
		repo.db.Model(&models.VulnerabilityActionPlanCommentDocument{}).
			Where("vulnerability_action_plan_comment_id = ?", comment.ID).
			Find(&documents)
		comment.Documents = documents
	}

	return comments, nil
}

func (repo *VulnerabilityActionPlanCommentRepository) GetByID(id uint) (*models.VulnerabilityActionPlanComment, error) {
	var vulnerabilityActionPlanComment *models.VulnerabilityActionPlanComment
	result := repo.db.Preload("Documents").First(&vulnerabilityActionPlanComment, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return vulnerabilityActionPlanComment, nil
}

func (repo *VulnerabilityActionPlanCommentRepository) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.VulnerabilityActionPlanComment, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var comment models.VulnerabilityActionPlanComment
	result := db.
		Preload("VulnerabilityActionPlan").
		Preload("CreatedByUser").
		Preload("Documents").
		Where("id = ?", id).
		First(&comment)

	if result.Error != nil {
		return nil, result.Error
	}

	return &comment, nil
}

func (repo *VulnerabilityActionPlanCommentRepository) Create(comment *models.VulnerabilityActionPlanComment, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(comment).Error; err != nil {
		return 0, err
	}

	return comment.ID, nil
}

func (repo *VulnerabilityActionPlanCommentRepository) Update(comment *models.VulnerabilityActionPlanComment, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Model(&models.VulnerabilityActionPlanComment{}).Where("id = ?", comment.ID).Updates(comment).Error; err != nil {
		return err
	}

	return nil
}

func (repo *VulnerabilityActionPlanCommentRepository) Delete(actionPlanID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("id = ?", actionPlanID).Delete(&models.VulnerabilityActionPlanComment{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *VulnerabilityActionPlanCommentRepository) CreateDocuments(docs []*models.VulnerabilityActionPlanCommentDocument, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Create(&docs).Error; err != nil {
		return err
	}
	return nil
}

func (repo *VulnerabilityActionPlanCommentRepository) DeleteDocuments(docIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("id IN (?)", docIds).Delete(&models.VulnerabilityActionPlanCommentDocument{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *VulnerabilityActionPlanCommentRepository) GetDocumentByID(id uint) (*models.VulnerabilityActionPlanCommentDocument, error) {
	var document models.VulnerabilityActionPlanCommentDocument
	if err := repo.db.Where("id = ?", id).First(&document).Error; err != nil {
		return nil, err
	}
	return &document, nil
}
