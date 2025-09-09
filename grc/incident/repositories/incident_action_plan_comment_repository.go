package repositories

import (
	incidentDTOs "jethur/grc/incident/dtos"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type IncidentActionPlanCommentRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(actionPlanID uint) ([]*incidentDTOs.IncidentActionPlanCommentListDTO, error)
	GetByID(id uint) (*models.IncidentActionPlanComment, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.IncidentActionPlanComment, error)
	Create(comment *models.IncidentActionPlanComment, tx ...utility.DatabaseTransaction) (uint, error)
	Update(comment *models.IncidentActionPlanComment, tx ...utility.DatabaseTransaction) error
	Delete(actionPlanID uint, tx ...utility.DatabaseTransaction) error

	CreateDocuments(docs []*models.IncidentActionPlanCommentDocument, tx ...utility.DatabaseTransaction) error
	DeleteDocuments(docIds []*uint, tx ...utility.DatabaseTransaction) error
	GetDocumentByID(id uint) (*models.IncidentActionPlanCommentDocument, error)
}

type IncidentActionPlanCommentRepository struct {
	db *gorm.DB
}

func NewIncidentActionPlanCommentRepository(db *gorm.DB) IncidentActionPlanCommentRepositoryInterface {
	return &IncidentActionPlanCommentRepository{
		db: db,
	}
}

func (repo *IncidentActionPlanCommentRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *IncidentActionPlanCommentRepository) GetAll(actionPlanID uint) ([]*incidentDTOs.IncidentActionPlanCommentListDTO, error) {
	var comments []*incidentDTOs.IncidentActionPlanCommentListDTO
	query := repo.db.Table("incident_action_plan_comments").
		Select(`
	incident_action_plan_comments.id,
	incident_action_plan_comments.comment,
	incident_action_plan_comments.incident_action_plan_id,
	incident_action_plan_comments.created_at,
	incident_action_plan_comments.created_by_id,
	users.first_name_en AS created_user_first_name_en,
	users.first_name_ar AS created_user_first_name_ar,
	users.last_name_en AS created_user_last_name_en,
	users.last_name_ar AS created_user_last_name_ar,
	users.initial AS created_user_initial,
	users.color AS created_user_color,
	users.email AS created_user_email,
	users.profile_image_path AS created_user_profile_image_path
`).
		Joins("LEFT JOIN users ON users.id = incident_action_plan_comments.created_by_id").
		Where("incident_action_plan_comments.incident_action_plan_id = ?", actionPlanID)

	if err := query.Find(&comments).Error; err != nil {
		return nil, err
	}

	for _, comment := range comments {
		var documents []*incidentDTOs.IncidentActionPlanCommentDocumentDTO
		repo.db.Model(&models.IncidentActionPlanCommentDocument{}).
			Where("incident_action_plan_comment_id = ?", comment.ID).
			Find(&documents)
		comment.Documents = documents
	}

	return comments, nil
}

func (repo *IncidentActionPlanCommentRepository) GetByID(id uint) (*models.IncidentActionPlanComment, error) {
	var incidentActionPlanComment *models.IncidentActionPlanComment
	result := repo.db.Preload("Documents").First(&incidentActionPlanComment, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return incidentActionPlanComment, nil
}

func (repo *IncidentActionPlanCommentRepository) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.IncidentActionPlanComment, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var comment models.IncidentActionPlanComment
	result := db.
		Preload("IncidentActionPlan").
		Preload("CreatedByUser").
		Preload("Documents").
		Where("id = ?", id).
		First(&comment)

	if result.Error != nil {
		return nil, result.Error
	}

	return &comment, nil
}

func (repo *IncidentActionPlanCommentRepository) Create(comment *models.IncidentActionPlanComment, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(comment).Error; err != nil {
		return 0, err
	}

	return comment.ID, nil
}

func (repo *IncidentActionPlanCommentRepository) Update(comment *models.IncidentActionPlanComment, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Model(&models.IncidentActionPlanComment{}).Where("id = ?", comment.ID).Updates(comment).Error; err != nil {
		return err
	}

	return nil
}

func (repo *IncidentActionPlanCommentRepository) Delete(actionPlanID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("id = ?", actionPlanID).Delete(&models.IncidentActionPlanComment{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *IncidentActionPlanCommentRepository) CreateDocuments(docs []*models.IncidentActionPlanCommentDocument, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Create(&docs).Error; err != nil {
		return err
	}
	return nil
}

func (repo *IncidentActionPlanCommentRepository) DeleteDocuments(docIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("id IN (?)", docIds).Delete(&models.IncidentActionPlanCommentDocument{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *IncidentActionPlanCommentRepository) GetDocumentByID(id uint) (*models.IncidentActionPlanCommentDocument, error) {
	var document models.IncidentActionPlanCommentDocument
	if err := repo.db.Where("id = ?", id).First(&document).Error; err != nil {
		return nil, err
	}
	return &document, nil
}
