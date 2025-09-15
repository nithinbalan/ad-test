package repositories

import (
	"jethur/grc/models"
	policyDTOs "jethur/grc/policy/dtos"
	"jethur/utility"

	"gorm.io/gorm"
)

type PolicyVersionClauseCommentRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(actionPlanID uint, params *policyDTOs.PolicyVersionClauseCommentListQueryParamsDTO) ([]*policyDTOs.PolicyVersionClauseCommentListDTO, error)
	GetByID(id uint) (*models.DocumentVersionClauseComment, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.DocumentVersionClauseComment, error)
	Create(comment *models.DocumentVersionClauseComment, tx ...utility.DatabaseTransaction) (uint, error)
	Update(comment *models.DocumentVersionClauseComment, tx ...utility.DatabaseTransaction) error
	Delete(actionPlanID uint, tx ...utility.DatabaseTransaction) error
}

type PolicyVersionClauseCommentRepository struct {
	db *gorm.DB
}

func NewPolicyVersionClauseCommentRepository(db *gorm.DB) PolicyVersionClauseCommentRepositoryInterface {
	return &PolicyVersionClauseCommentRepository{
		db: db,
	}
}

func (repo *PolicyVersionClauseCommentRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *PolicyVersionClauseCommentRepository) GetAll(actionPlanID uint, params *policyDTOs.PolicyVersionClauseCommentListQueryParamsDTO) ([]*policyDTOs.PolicyVersionClauseCommentListDTO, error) {
	var comments []*policyDTOs.PolicyVersionClauseCommentListDTO
	query := repo.db.Table("document_version_clause_comments").
		Select(`
	document_version_clause_comments.id,
	document_version_clause_comments.comment,
	document_version_clause_comments.document_version_clause_id,
	document_version_clause_comments.level,
	document_version_clause_comments.created_at,
	document_version_clause_comments.created_by_id,
	comment_statuses.title_en AS comment_status_title_en,
	comment_statuses.title_ar AS comment_status_title_ar,
	comment_statuses.type AS comment_status_type,
	comment_statuses.color_code AS comment_status_color_code,
	users.first_name_en AS created_user_first_name_en,
	users.first_name_ar AS created_user_first_name_ar,
	users.last_name_en AS created_user_last_name_en,
	users.last_name_ar AS created_user_last_name_ar,
	users.initial AS created_user_initial,
	users.color AS created_user_color,
	users.email AS created_user_email,
	users.profile_image_path AS created_user_profile_image_path
`).
		Joins("LEFT JOIN users ON users.id = document_version_clause_comments.created_by_id").
		Joins("LEFT JOIN comment_statuses ON comment_statuses.id = document_version_clause_comments.comment_status_id").
		Where("document_version_clause_comments.document_version_clause_id = ?", actionPlanID)

		// Apply dynamic filters to the query
	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where(columnName+" IN (?)", values)
		}
	}

	// Apply search keyword if exists
	if params.Keyword != "" {
		query = query.Where(`
				(document_version_clause_comments.comment LIKE ? OR
				documents.title_ar LIKE ? OR
				documents.reference_number LIKE ? OR
				documents.description_en LIKE ? OR
				documents.description_ar LIKE ? OR
				document_categories.title_en LIKE ? OR
				document_categories.title_ar LIKE ? OR
				document_statuses.title_en LIKE ? OR
				document_statuses.title_ar LIKE ? OR
				document_types.title_en LIKE ? OR
				document_types.title_ar LIKE ? OR
				organization_units.title_en LIKE ? OR
				organization_units.title_ar LIKE ? OR
				responsible_user.first_name_en LIKE ? OR
				responsible_user.first_name_ar LIKE ? OR
				responsible_user.last_name_en LIKE ? OR
				responsible_user.last_name_ar LIKE ? OR
				responsible_user.email LIKE ? OR
				responsible_user.initial LIKE ?)`,
			"%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%",
			"%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%",
			"%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%",
			"%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%",
			"%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%")
	}
	if err := query.Find(&comments).Error; err != nil {
		return nil, err
	}

	return comments, nil
}

func (repo *PolicyVersionClauseCommentRepository) GetByID(id uint) (*models.DocumentVersionClauseComment, error) {
	var PolicyVersionClauseComment *models.DocumentVersionClauseComment
	result := repo.db.First(&PolicyVersionClauseComment, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return PolicyVersionClauseComment, nil
}

func (repo *PolicyVersionClauseCommentRepository) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.DocumentVersionClauseComment, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var comment models.DocumentVersionClauseComment
	result := db.
		Preload("DocumentVersionClause").
		Preload("CreatedByUser").
		Where("id = ?", id).
		First(&comment)

	if result.Error != nil {
		return nil, result.Error
	}

	return &comment, nil
}

func (repo *PolicyVersionClauseCommentRepository) Create(comment *models.DocumentVersionClauseComment, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(comment).Error; err != nil {
		return 0, err
	}

	return comment.ID, nil
}

func (repo *PolicyVersionClauseCommentRepository) Update(comment *models.DocumentVersionClauseComment, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Model(&models.DocumentVersionClauseComment{}).Where("id = ?", comment.ID).Updates(comment).Error; err != nil {
		return err
	}

	return nil
}

func (repo *PolicyVersionClauseCommentRepository) Delete(actionPlanID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("id = ?", actionPlanID).Delete(&models.DocumentVersionClauseComment{}).Error; err != nil {
		return err
	}

	return nil
}
