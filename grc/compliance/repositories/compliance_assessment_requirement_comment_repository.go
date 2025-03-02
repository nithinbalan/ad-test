package repositories

import (
	complianceDTOs "jethur/grc/compliance/dtos"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type ComplianceAssessmentRequirementCommentRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(assessmentRequirementID uint) ([]*complianceDTOs.ComplianceAssessmentRequirementCommentListDTO, error)
	GetByID(id uint) (*models.ComplianceAssessmentRequirementComment, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.ComplianceAssessmentRequirementComment, error)
	Create(comment *models.ComplianceAssessmentRequirementComment, tx ...utility.DatabaseTransaction) (uint, error)
	Update(comment *models.ComplianceAssessmentRequirementComment, tx ...utility.DatabaseTransaction) error
	Delete(assessmentRequirementID uint, tx ...utility.DatabaseTransaction) error
}

type ComplianceAssessmentRequirementCommentRepository struct {
	db *gorm.DB
}

func NewComplianceAssessmentRequirementCommentRepository(db *gorm.DB) ComplianceAssessmentRequirementCommentRepositoryInterface {
	return &ComplianceAssessmentRequirementCommentRepository{
		db: db,
	}
}

func (repo *ComplianceAssessmentRequirementCommentRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *ComplianceAssessmentRequirementCommentRepository) GetAll(assessmentRequirementID uint) ([]*complianceDTOs.ComplianceAssessmentRequirementCommentListDTO, error) {
	var comments []*complianceDTOs.ComplianceAssessmentRequirementCommentListDTO
	query := repo.db.Table("compliance_assessment_requirement_comments").
		Select(`
	compliance_assessment_requirement_comments.id,
	compliance_assessment_requirement_comments.comment,
	compliance_assessment_requirement_comments.compliance_assessment_requirement_id,
	compliance_assessment_requirement_comments.created_at,
	compliance_assessment_requirement_comments.created_by_id,
	users.first_name_en AS created_user_first_name_en,
	users.first_name_ar AS created_user_first_name_ar,
	users.last_name_en AS created_user_last_name_en,
	users.last_name_ar AS created_user_last_name_ar,
	users.initial AS created_user_initial,
	users.color AS created_user_color,
	users.email AS created_user_email,
	users.profile_image_path AS created_user_profile_image_path
`).
		Joins("LEFT JOIN users ON users.id = compliance_assessment_requirement_comments.created_by_id").
		Where("compliance_assessment_requirement_comments.compliance_assessment_requirement_id = ?", assessmentRequirementID)

	if err := query.Find(&comments).Error; err != nil {
		return nil, err
	}

	return comments, nil
}

func (repo *ComplianceAssessmentRequirementCommentRepository) GetByID(id uint) (*models.ComplianceAssessmentRequirementComment, error) {
	var complianceAssessmentRequirementComment *models.ComplianceAssessmentRequirementComment
	result := repo.db.First(&complianceAssessmentRequirementComment, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return complianceAssessmentRequirementComment, nil
}

func (repo *ComplianceAssessmentRequirementCommentRepository) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.ComplianceAssessmentRequirementComment, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var comment models.ComplianceAssessmentRequirementComment
	result := db.
		Preload("ComplianceAssessmentRequirement").
		Preload("CreatedByUser").
		Where("id = ?", id).
		First(&comment)

	if result.Error != nil {
		return nil, result.Error
	}

	return &comment, nil
}

func (repo *ComplianceAssessmentRequirementCommentRepository) Create(comment *models.ComplianceAssessmentRequirementComment, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(comment).Error; err != nil {
		return 0, err
	}

	return comment.ID, nil
}

func (repo *ComplianceAssessmentRequirementCommentRepository) Update(comment *models.ComplianceAssessmentRequirementComment, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Model(&models.ComplianceAssessmentRequirementComment{}).Where("id = ?", comment.ID).Updates(comment).Error; err != nil {
		return err
	}

	return nil
}

func (repo *ComplianceAssessmentRequirementCommentRepository) Delete(assessmentRequirementID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("id = ?", assessmentRequirementID).Delete(&models.ComplianceAssessmentRequirementComment{}).Error; err != nil {
		return err
	}

	return nil
}
