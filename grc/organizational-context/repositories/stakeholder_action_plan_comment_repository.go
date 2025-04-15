package repositories

import (
	"jethur/grc/models"
	organizationalContextDTOs "jethur/grc/organizational-context/dtos"
	"jethur/utility"

	"gorm.io/gorm"
)

type StakeholderActionPlanCommentRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(actionPlanID uint) ([]*organizationalContextDTOs.StakeholderActionPlanCommentListDTO, error)
	GetByID(id uint) (*models.StakeholderActionPlanComment, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.StakeholderActionPlanComment, error)
	Create(comment *models.StakeholderActionPlanComment, tx ...utility.DatabaseTransaction) (uint, error)
	Update(comment *models.StakeholderActionPlanComment, tx ...utility.DatabaseTransaction) error
	Delete(actionPlanID uint, tx ...utility.DatabaseTransaction) error

	CreateDocuments(docs []*models.StakeholderActionPlanCommentDocument, tx ...utility.DatabaseTransaction) error
	DeleteDocuments(docIds []*uint, tx ...utility.DatabaseTransaction) error
	GetDocumentByID(id uint) (*models.StakeholderActionPlanCommentDocument, error)
}

type StakeholderActionPlanCommentRepository struct {
	db *gorm.DB
}

func NewStakeholderActionPlanCommentRepository(db *gorm.DB) StakeholderActionPlanCommentRepositoryInterface {
	return &StakeholderActionPlanCommentRepository{
		db: db,
	}
}

func (repo *StakeholderActionPlanCommentRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *StakeholderActionPlanCommentRepository) GetAll(actionPlanID uint) ([]*organizationalContextDTOs.StakeholderActionPlanCommentListDTO, error) {
	var comments []*organizationalContextDTOs.StakeholderActionPlanCommentListDTO
	query := repo.db.Table("stakeholder_action_plan_comments").
		Select(`
	stakeholder_action_plan_comments.id,
	stakeholder_action_plan_comments.comment,
	stakeholder_action_plan_comments.stakeholder_action_plan_id,
	stakeholder_action_plan_comments.created_at,
	stakeholder_action_plan_comments.created_by_id,
	users.first_name_en AS created_user_first_name_en,
	users.first_name_ar AS created_user_first_name_ar,
	users.last_name_en AS created_user_last_name_en,
	users.last_name_ar AS created_user_last_name_ar,
	users.initial AS created_user_initial,
	users.color AS created_user_color,
	users.email AS created_user_email,
	users.profile_image_path AS created_user_profile_image_path
`).
		Joins("LEFT JOIN users ON users.id = stakeholder_action_plan_comments.created_by_id").
		Where("stakeholder_action_plan_comments.stakeholder_action_plan_id = ?", actionPlanID)

	if err := query.Find(&comments).Error; err != nil {
		return nil, err
	}

	for _, comment := range comments {
		var documents []*organizationalContextDTOs.StakeholderActionPlanCommentDocumentDTO
		repo.db.Model(&models.StakeholderActionPlanCommentDocument{}).
			Where("stakeholder_action_plan_comment_id = ?", comment.ID).
			Find(&documents)
		comment.Documents = documents
	}

	return comments, nil
}

func (repo *StakeholderActionPlanCommentRepository) GetByID(id uint) (*models.StakeholderActionPlanComment, error) {
	var stakeholderActionPlanComment *models.StakeholderActionPlanComment
	result := repo.db.Preload("Documents").First(&stakeholderActionPlanComment, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return stakeholderActionPlanComment, nil
}

func (repo *StakeholderActionPlanCommentRepository) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.StakeholderActionPlanComment, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var comment models.StakeholderActionPlanComment
	result := db.
		Preload("StakeholderActionPlan").
		Preload("CreatedByUser").
		Preload("Documents").
		Where("id = ?", id).
		First(&comment)

	if result.Error != nil {
		return nil, result.Error
	}

	return &comment, nil
}

func (repo *StakeholderActionPlanCommentRepository) Create(comment *models.StakeholderActionPlanComment, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(comment).Error; err != nil {
		return 0, err
	}

	return comment.ID, nil
}

func (repo *StakeholderActionPlanCommentRepository) Update(comment *models.StakeholderActionPlanComment, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Model(&models.StakeholderActionPlanComment{}).Where("id = ?", comment.ID).Updates(comment).Error; err != nil {
		return err
	}

	return nil
}

func (repo *StakeholderActionPlanCommentRepository) Delete(actionPlanID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("id = ?", actionPlanID).Delete(&models.StakeholderActionPlanComment{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *StakeholderActionPlanCommentRepository) CreateDocuments(docs []*models.StakeholderActionPlanCommentDocument, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Create(&docs).Error; err != nil {
		return err
	}
	return nil
}

func (repo *StakeholderActionPlanCommentRepository) DeleteDocuments(docIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("id IN (?)", docIds).Delete(&models.StakeholderActionPlanCommentDocument{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *StakeholderActionPlanCommentRepository) GetDocumentByID(id uint) (*models.StakeholderActionPlanCommentDocument, error) {
	var document models.StakeholderActionPlanCommentDocument
	if err := repo.db.Where("id = ?", id).First(&document).Error; err != nil {
		return nil, err
	}
	return &document, nil
}
