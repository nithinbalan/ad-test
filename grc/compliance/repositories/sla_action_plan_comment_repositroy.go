package repositories

import (
	complianceDTOs "jethur/grc/compliance/dtos"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type SLAActionPlanCommentRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(actionPlanID uint) ([]*complianceDTOs.SLAActionPlanCommentListDTO, error)
	GetByID(id uint) (*models.SLAActionPlanComment, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.SLAActionPlanComment, error)
	Create(comment *models.SLAActionPlanComment, tx ...utility.DatabaseTransaction) (uint, error)
	Update(comment *models.SLAActionPlanComment, tx ...utility.DatabaseTransaction) error
	Delete(actionPlanID uint, tx ...utility.DatabaseTransaction) error

	CreateDocuments(docs []*models.SLAActionPlanCommentDocument, tx ...utility.DatabaseTransaction) error
	DeleteDocuments(docIds []*uint, tx ...utility.DatabaseTransaction) error
	GetDocumentByID(id uint) (*models.SLAActionPlanCommentDocument, error)
}

type SLAActionPlanCommentRepository struct {
	db *gorm.DB
}

func NewSLAActionPlanCommentRepository(db *gorm.DB) SLAActionPlanCommentRepositoryInterface {
	return &SLAActionPlanCommentRepository{
		db: db,
	}
}

func (repo *SLAActionPlanCommentRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *SLAActionPlanCommentRepository) GetAll(actionPlanID uint) ([]*complianceDTOs.SLAActionPlanCommentListDTO, error) {
	var comments []*complianceDTOs.SLAActionPlanCommentListDTO
	query := repo.db.Table("sla_action_plan_comments").
		Select(`
	sla_action_plan_comments.id,
	sla_action_plan_comments.comment,
	sla_action_plan_comments.sla_action_plan_id,
	sla_action_plan_comments.created_at,
	sla_action_plan_comments.created_by_id,
	users.first_name_en AS created_user_first_name_en,
	users.first_name_ar AS created_user_first_name_ar,
	users.last_name_en AS created_user_last_name_en,
	users.last_name_ar AS created_user_last_name_ar,
	users.initial AS created_user_initial,
	users.color AS created_user_color,
	users.email AS created_user_email,
	users.profile_image_path AS created_user_profile_image_path
`).
		Joins("LEFT JOIN users ON users.id = sla_action_plan_comments.created_by_id").
		Where("sla_action_plan_comments.sla_action_plan_id = ?", actionPlanID)

	if err := query.Find(&comments).Error; err != nil {
		return nil, err
	}

	for _, comment := range comments {
		var documents []*complianceDTOs.SLAActionPlanCommentDocumentDTO
		repo.db.Model(&models.SLAActionPlanCommentDocument{}).
			Where("sla_action_plan_comment_id = ?", comment.ID).
			Find(&documents)
		comment.Documents = documents
	}

	return comments, nil
}

func (repo *SLAActionPlanCommentRepository) GetByID(id uint) (*models.SLAActionPlanComment, error) {
	var slaActionPlanComment *models.SLAActionPlanComment
	result := repo.db.Preload("Documents").First(&slaActionPlanComment, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return slaActionPlanComment, nil
}

func (repo *SLAActionPlanCommentRepository) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.SLAActionPlanComment, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var comment models.SLAActionPlanComment
	result := db.
		Preload("SLAActionPlan").
		Preload("CreatedByUser").
		Preload("Documents").
		Where("id = ?", id).
		First(&comment)

	if result.Error != nil {
		return nil, result.Error
	}

	return &comment, nil
}

func (repo *SLAActionPlanCommentRepository) Create(comment *models.SLAActionPlanComment, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(comment).Error; err != nil {
		return 0, err
	}

	return comment.ID, nil
}

func (repo *SLAActionPlanCommentRepository) Update(comment *models.SLAActionPlanComment, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Model(&models.SLAActionPlanComment{}).Where("id = ?", comment.ID).Updates(comment).Error; err != nil {
		return err
	}

	return nil
}

func (repo *SLAActionPlanCommentRepository) Delete(actionPlanID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("id = ?", actionPlanID).Delete(&models.SLAActionPlanComment{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *SLAActionPlanCommentRepository) CreateDocuments(docs []*models.SLAActionPlanCommentDocument, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Create(&docs).Error; err != nil {
		return err
	}
	return nil
}

func (repo *SLAActionPlanCommentRepository) DeleteDocuments(docIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("id IN (?)", docIds).Delete(&models.SLAActionPlanCommentDocument{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *SLAActionPlanCommentRepository) GetDocumentByID(id uint) (*models.SLAActionPlanCommentDocument, error) {
	var document models.SLAActionPlanCommentDocument
	if err := repo.db.Where("id = ?", id).First(&document).Error; err != nil {
		return nil, err
	}
	return &document, nil
}
