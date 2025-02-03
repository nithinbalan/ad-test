package repositories

import (
	complianceDTOs "jethur/grc/compliance/dtos"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type ComplianceActionPlanCommentRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(actionPlanID uint) ([]*complianceDTOs.ComplianceActionPlanCommentListDTO, error)
	GetByID(id uint) (*models.ComplianceActionPlanComment, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.ComplianceActionPlanComment, error)
	Create(comment *models.ComplianceActionPlanComment, tx ...utility.DatabaseTransaction) (uint, error)
	Update(comment *models.ComplianceActionPlanComment, tx ...utility.DatabaseTransaction) error
	Delete(actionPlanID uint, tx ...utility.DatabaseTransaction) error

	CreateDocuments(docs []*models.ComplianceActionPlanCommentDocument, tx ...utility.DatabaseTransaction) error
	DeleteDocuments(docIds []*uint, tx ...utility.DatabaseTransaction) error
	GetDocumentByID(id uint) (*models.ComplianceActionPlanCommentDocument, error)
}

type ComplianceActionPlanCommentRepository struct {
	db *gorm.DB
}

func NewComplianceActionPlanCommentRepository(db *gorm.DB) ComplianceActionPlanCommentRepositoryInterface {
	return &ComplianceActionPlanCommentRepository{
		db: db,
	}
}

func (repo *ComplianceActionPlanCommentRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *ComplianceActionPlanCommentRepository) GetAll(actionPlanID uint) ([]*complianceDTOs.ComplianceActionPlanCommentListDTO, error) {
	var comments []*complianceDTOs.ComplianceActionPlanCommentListDTO
	query := repo.db.Table("compliance_action_plan_comments").
		Select(`
	compliance_action_plan_comments.id,
	compliance_action_plan_comments.comment,
	compliance_action_plan_comments.compliance_action_plan_id,
	compliance_action_plan_comments.created_at,
	compliance_action_plan_comments.created_by_id,
	users.first_name_en AS created_user_first_name_en,
	users.first_name_ar AS created_user_first_name_ar,
	users.last_name_en AS created_user_last_name_en,
	users.last_name_ar AS created_user_last_name_ar,
	users.initial AS created_user_initial,
	users.color AS created_user_color,
	users.email AS created_user_email,
	users.profile_image_path AS created_user_profile_image_path
`).
		Joins("LEFT JOIN users ON users.id = compliance_action_plan_comments.created_by_id").
		Where("compliance_action_plan_comments.compliance_action_plan_id = ?", actionPlanID)

	if err := query.Find(&comments).Error; err != nil {
		return nil, err
	}

	for _, comment := range comments {
		var documents []*complianceDTOs.ComplianceActionPlanCommentDocumentDTO
		repo.db.Model(&models.ComplianceActionPlanCommentDocument{}).
			Where("compliance_action_plan_comment_id = ?", comment.ID).
			Find(&documents)
		comment.Documents = documents
	}

	return comments, nil
}

func (repo *ComplianceActionPlanCommentRepository) GetByID(id uint) (*models.ComplianceActionPlanComment, error) {
	var complianceActionPlanComment *models.ComplianceActionPlanComment
	result := repo.db.Preload("Documents").First(&complianceActionPlanComment, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return complianceActionPlanComment, nil
}

func (repo *ComplianceActionPlanCommentRepository) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.ComplianceActionPlanComment, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var comment models.ComplianceActionPlanComment
	result := db.
		Preload("ComplianceActionPlan").
		Preload("CreatedByUser").
		Preload("Documents").
		Where("id = ?", id).
		First(&comment)

	if result.Error != nil {
		return nil, result.Error
	}

	return &comment, nil
}

func (repo *ComplianceActionPlanCommentRepository) Create(comment *models.ComplianceActionPlanComment, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(comment).Error; err != nil {
		return 0, err
	}

	return comment.ID, nil
}

func (repo *ComplianceActionPlanCommentRepository) Update(comment *models.ComplianceActionPlanComment, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Model(&models.ComplianceActionPlanComment{}).Where("id = ?", comment.ID).Updates(comment).Error; err != nil {
		return err
	}

	return nil
}

func (repo *ComplianceActionPlanCommentRepository) Delete(actionPlanID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("id = ?", actionPlanID).Delete(&models.ComplianceActionPlanComment{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *ComplianceActionPlanCommentRepository) CreateDocuments(docs []*models.ComplianceActionPlanCommentDocument, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Create(&docs).Error; err != nil {
		return err
	}
	return nil
}

func (repo *ComplianceActionPlanCommentRepository) DeleteDocuments(docIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("id IN (?)", docIds).Delete(&models.ComplianceActionPlanCommentDocument{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *ComplianceActionPlanCommentRepository) GetDocumentByID(id uint) (*models.ComplianceActionPlanCommentDocument, error) {
	var document models.ComplianceActionPlanCommentDocument
	if err := repo.db.Where("id = ?", id).First(&document).Error; err != nil {
		return nil, err
	}
	return &document, nil
}
