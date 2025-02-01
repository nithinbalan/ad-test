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
	Create(comment *models.ComplianceActionPlanComment) (uint, error)
	Update(comment *models.ComplianceActionPlanComment, tx ...utility.DatabaseTransaction) error
	Delete(actionPlanID uint, tx ...utility.DatabaseTransaction) error
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
	return comments, nil
}

func (repo *ComplianceActionPlanCommentRepository) GetByID(id uint) (*models.ComplianceActionPlanComment, error) {
	var complianceActionPlan *models.ComplianceActionPlanComment
	result := repo.db.First(&complianceActionPlan, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return complianceActionPlan, nil
}

func (repo *ComplianceActionPlanCommentRepository) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.ComplianceActionPlanComment, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var comment *models.ComplianceActionPlanComment
	result := db.
		Preload("ComplianceActionPlan").
		Preload("CreatedByUser").
		Where("id = ?", id).
		First(&comment)

	if result.Error != nil {
		return nil, result.Error
	}

	return comment, nil
}

func (repo *ComplianceActionPlanCommentRepository) Create(comment *models.ComplianceActionPlanComment) (uint, error) {
	if err := repo.db.Create(comment).Error; err != nil {
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
