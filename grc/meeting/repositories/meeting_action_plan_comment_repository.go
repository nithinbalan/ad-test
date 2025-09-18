package repositories

import (
	meetingDTOs "jethur/grc/meeting/dtos"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type MeetingActionPlanCommentRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(actionPlanID uint) ([]*meetingDTOs.MeetingActionPlanCommentListDTO, error)
	GetByID(id uint) (*models.MeetingActionPlanComment, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.MeetingActionPlanComment, error)
	Create(comment *models.MeetingActionPlanComment, tx ...utility.DatabaseTransaction) (uint, error)
	Update(comment *models.MeetingActionPlanComment, tx ...utility.DatabaseTransaction) error
	Delete(actionPlanID uint, tx ...utility.DatabaseTransaction) error

	CreateDocuments(docs []*models.MeetingActionPlanCommentDocument, tx ...utility.DatabaseTransaction) error
	DeleteDocuments(docIds []*uint, tx ...utility.DatabaseTransaction) error
	GetDocumentByID(id uint) (*models.MeetingActionPlanCommentDocument, error)
}

type MeetingActionPlanCommentRepository struct {
	db *gorm.DB
}

func NewMeetingActionPlanCommentRepository(db *gorm.DB) MeetingActionPlanCommentRepositoryInterface {
	return &MeetingActionPlanCommentRepository{
		db: db,
	}
}

func (repo *MeetingActionPlanCommentRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *MeetingActionPlanCommentRepository) GetAll(actionPlanID uint) ([]*meetingDTOs.MeetingActionPlanCommentListDTO, error) {
	var comments []*meetingDTOs.MeetingActionPlanCommentListDTO
	query := repo.db.Table("meeting_action_plan_comments").
		Select(`
	meeting_action_plan_comments.id,
	meeting_action_plan_comments.comment,
	meeting_action_plan_comments.meeting_action_plan_id,
	meeting_action_plan_comments.created_at,
	meeting_action_plan_comments.created_by_id,
	users.first_name_en AS created_user_first_name_en,
	users.first_name_ar AS created_user_first_name_ar,
	users.last_name_en AS created_user_last_name_en,
	users.last_name_ar AS created_user_last_name_ar,
	users.initial AS created_user_initial,
	users.color AS created_user_color,
	users.email AS created_user_email,
	users.profile_image_path AS created_user_profile_image_path
`).
		Joins("LEFT JOIN users ON users.id = meeting_action_plan_comments.created_by_id").
		Where("meeting_action_plan_comments.meeting_action_plan_id = ?", actionPlanID)

	if err := query.Find(&comments).Error; err != nil {
		return nil, err
	}

	for _, comment := range comments {
		var documents []*meetingDTOs.MeetingActionPlanCommentDocumentDTO
		repo.db.Model(&models.MeetingActionPlanCommentDocument{}).
			Where("meeting_action_plan_comment_id = ?", comment.ID).
			Find(&documents)
		comment.Documents = documents
	}

	return comments, nil
}

func (repo *MeetingActionPlanCommentRepository) GetByID(id uint) (*models.MeetingActionPlanComment, error) {
	var meetingActionPlanComment *models.MeetingActionPlanComment
	result := repo.db.Preload("Documents").First(&meetingActionPlanComment, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return meetingActionPlanComment, nil
}

func (repo *MeetingActionPlanCommentRepository) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.MeetingActionPlanComment, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var comment models.MeetingActionPlanComment
	result := db.
		Preload("MeetingActionPlan").
		Preload("CreatedByUser").
		Preload("Documents").
		Where("id = ?", id).
		First(&comment)

	if result.Error != nil {
		return nil, result.Error
	}

	return &comment, nil
}

func (repo *MeetingActionPlanCommentRepository) Create(comment *models.MeetingActionPlanComment, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(comment).Error; err != nil {
		return 0, err
	}

	return comment.ID, nil
}

func (repo *MeetingActionPlanCommentRepository) Update(comment *models.MeetingActionPlanComment, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Model(&models.MeetingActionPlanComment{}).Where("id = ?", comment.ID).Updates(comment).Error; err != nil {
		return err
	}

	return nil
}

func (repo *MeetingActionPlanCommentRepository) Delete(actionPlanID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("id = ?", actionPlanID).Delete(&models.MeetingActionPlanComment{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *MeetingActionPlanCommentRepository) CreateDocuments(docs []*models.MeetingActionPlanCommentDocument, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Create(&docs).Error; err != nil {
		return err
	}
	return nil
}

func (repo *MeetingActionPlanCommentRepository) DeleteDocuments(docIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("id IN (?)", docIds).Delete(&models.MeetingActionPlanCommentDocument{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *MeetingActionPlanCommentRepository) GetDocumentByID(id uint) (*models.MeetingActionPlanCommentDocument, error) {
	var document models.MeetingActionPlanCommentDocument
	if err := repo.db.Where("id = ?", id).First(&document).Error; err != nil {
		return nil, err
	}
	return &document, nil
}
