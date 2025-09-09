package repositories

import (
	complianceDTOs "jethur/grc/compliance/dtos"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type ComplianceAssessmentControlCommentRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(controlID uint) ([]*complianceDTOs.ComplianceAssessmentControlCommentListDTO, error)
	GetByID(id uint) (*models.ComplianceAssessmentControlComment, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.ComplianceAssessmentControlComment, error)
	Create(comment *models.ComplianceAssessmentControlComment, tx ...utility.DatabaseTransaction) (uint, error)
	Update(comment *models.ComplianceAssessmentControlComment, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error
}

type ComplianceAssessmentControlCommentRepository struct {
	db *gorm.DB
}

func NewComplianceAssessmentControlCommentRepository(db *gorm.DB) ComplianceAssessmentControlCommentRepositoryInterface {
	return &ComplianceAssessmentControlCommentRepository{
		db: db,
	}
}

func (repo *ComplianceAssessmentControlCommentRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *ComplianceAssessmentControlCommentRepository) GetAll(controlID uint) ([]*complianceDTOs.ComplianceAssessmentControlCommentListDTO, error) {
	var comments []*complianceDTOs.ComplianceAssessmentControlCommentListDTO
	query := repo.db.Table("compliance_assessment_control_comments").
		Select(`
			compliance_assessment_control_comments.id,
			compliance_assessment_control_comments.comment,
			compliance_assessment_control_comments.compliance_assessment_control_id,
			compliance_assessment_control_comments.created_at,
			compliance_assessment_control_comments.created_by_id,
			users.first_name_en AS created_user_first_name_en,
			users.first_name_ar AS created_user_first_name_ar,
			users.last_name_en AS created_user_last_name_en,
			users.last_name_ar AS created_user_last_name_ar,
			users.initial AS created_user_initial,
			users.color AS created_user_color,
			users.email AS created_user_email,
			users.profile_image_path AS created_user_profile_image_path
		`).
		Joins("LEFT JOIN users ON users.id = compliance_assessment_control_comments.created_by_id").
		Where("compliance_assessment_control_comments.compliance_assessment_control_id = ?", controlID)

	if err := query.Find(&comments).Error; err != nil {
		return nil, err
	}

	return comments, nil
}

func (repo *ComplianceAssessmentControlCommentRepository) GetByID(id uint) (*models.ComplianceAssessmentControlComment, error) {
	var comment *models.ComplianceAssessmentControlComment
	result := repo.db.First(&comment, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return comment, nil
}

func (repo *ComplianceAssessmentControlCommentRepository) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.ComplianceAssessmentControlComment, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var comment models.ComplianceAssessmentControlComment
	result := db.
		Preload("ComplianceAssessmentControl").
		Preload("ComplianceAssessmentRequirementStatus").
		Preload("CreatedByUser").
		Where("id = ?", id).
		First(&comment)

	if result.Error != nil {
		return nil, result.Error
	}

	return &comment, nil
}

func (repo *ComplianceAssessmentControlCommentRepository) Create(comment *models.ComplianceAssessmentControlComment, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(comment).Error; err != nil {
		return 0, err
	}
	return comment.ID, nil
}

func (repo *ComplianceAssessmentControlCommentRepository) Update(comment *models.ComplianceAssessmentControlComment, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Model(&models.ComplianceAssessmentControlComment{}).Where("id = ?", comment.ID).Updates(comment).Error; err != nil {
		return err
	}

	return nil
}

func (repo *ComplianceAssessmentControlCommentRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("id = ?", id).Delete(&models.ComplianceAssessmentControlComment{}).Error; err != nil {
		return err
	}

	return nil
}
