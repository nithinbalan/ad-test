package repositories

import (
	meetingDTOs "jethur/grc/meeting/dtos"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type MeetingActionPlanStatusUpdateRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(meetingActionPlanID uint, sort string, sortBy string, keyword string) ([]*meetingDTOs.MeetingActionPlanStatusUpdateListDTO, error)
	GetByID(id uint) (*models.MeetingActionPlanStatusUpdate, error)
	Create(meetingActionPlanStatusUpdate *models.MeetingActionPlanStatusUpdate, tx ...utility.DatabaseTransaction) (uint, error)
	Update(meetingActionPlanStatusUpdate *models.MeetingActionPlanStatusUpdate, tx ...utility.DatabaseTransaction) error

	//Document
	CreateDocuments(docs []*models.MeetingActionPlanStatusUpdateDocument, tx ...utility.DatabaseTransaction) error
	GetDocumentByID(id uint) (*models.MeetingActionPlanStatusUpdateDocument, error)
	DeleteDocuments(docIds []*uint, tx ...utility.DatabaseTransaction) error

	GetLatestMeetingActionPlanStatusUpdate(meetingActionPlanID uint) (*meetingDTOs.MeetingActionPlanStatusUpdateListDTO, error)
	CheckAllActionPlansCompletedStatusCount(meetingMinuteID uint, excludeActionPlanID uint, completedStatuses []string, tx ...utility.DatabaseTransaction) (int64, error)
}

type MeetingActionPlanStatusUpdateRepository struct {
	db *gorm.DB
}

func NewMeetingActionPlanStatusUpdateRepository(db *gorm.DB) *MeetingActionPlanStatusUpdateRepository {
	return &MeetingActionPlanStatusUpdateRepository{
		db: db,
	}
}

// BeginTransaction starts a new transaction.
func (repo *MeetingActionPlanStatusUpdateRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *MeetingActionPlanStatusUpdateRepository) GetAll(meetingActionPlanID uint, sort string, sortBy string, keyword string) ([]*meetingDTOs.MeetingActionPlanStatusUpdateListDTO, error) {

	var meetingActionPlanUpdates []*meetingDTOs.MeetingActionPlanStatusUpdateListDTO

	query := repo.db.Table("meeting_action_plan_status_updates").
		Select(`
            meeting_action_plan_status_updates.*,
            meeting_action_plans.title AS meeting_action_plan_title,
            action_plan_statuses.title_en AS action_plan_status_title_en,
            action_plan_statuses.title_ar AS action_plan_status_title_ar,
            action_plan_statuses.color_code AS action_plan_status_color_code,
            created_user.id AS created_user_id,
        	created_user.first_name_en AS created_user_first_name_en,
            created_user.first_name_ar AS created_user_first_name_ar,
            created_user.last_name_en AS created_user_last_name_en,
            created_user.last_name_ar AS created_user_last_name_ar,
			created_user.initial AS created_user_initial,
        	created_user.color AS created_user_color,
        	created_user.email AS created_user_email,
        	created_user.profile_image_path AS created_user_profile_image_path,
            updated_user.id AS updated_user_id,
    		updated_user.first_name_en AS updated_user_first_name_en,
            updated_user.first_name_ar AS updated_user_first_name_ar,
            updated_user.last_name_en AS updated_user_last_name_en,
            updated_user.last_name_ar AS updated_user_last_name_ar,
			updated_user.initial AS updated_user_initial,
        	updated_user.color AS updated_user_color,
        	updated_user.email AS updated_user_email,
        	updated_user.profile_image_path AS updated_user_profile_image_path 
        `).
		Joins(`
            LEFT JOIN meeting_action_plans ON meeting_action_plans.id = meeting_action_plan_status_updates.meeting_action_plan_id
            LEFT JOIN action_plan_statuses ON action_plan_statuses.id = meeting_action_plan_status_updates.action_plan_status_id
            LEFT JOIN users AS created_user ON created_user.id = meeting_action_plan_status_updates.created_by_id
            LEFT JOIN users AS updated_user ON updated_user.id = meeting_action_plan_status_updates.updated_by_id
        `).
		Where("meeting_action_plan_status_updates.meeting_action_plan_id = ?", meetingActionPlanID)

	if keyword != "" {
		query = query.Where(`
            (meeting_action_plans.title LIKE ? OR
            meeting_action_plan_status_updates.comment LIKE ? OR
            action_plan_statuses.title_en LIKE ? OR
            action_plan_statuses.title_ar LIKE ?)`,
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}
	if sortBy == "" {
		sortBy = "meeting_action_plan_status_updates.id"
	}
	if sort != "asc" && sort != "desc" {
		sort = "asc"
	}
	query = query.Order(sortBy + " " + sort)

	result := query.Find(&meetingActionPlanUpdates)
	return meetingActionPlanUpdates, result.Error
}

func (repo *MeetingActionPlanStatusUpdateRepository) GetByID(id uint) (*models.MeetingActionPlanStatusUpdate, error) {
	var meetingActionPlanStatusUpdate *models.MeetingActionPlanStatusUpdate
	result := repo.db.Preload("Documents").First(&meetingActionPlanStatusUpdate, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return meetingActionPlanStatusUpdate, nil
}

func (repo *MeetingActionPlanStatusUpdateRepository) Create(meetingActionPlanStatusUpdate *models.MeetingActionPlanStatusUpdate, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(meetingActionPlanStatusUpdate).Error; err != nil {
		return 0, err
	}

	return meetingActionPlanStatusUpdate.ID, nil
}

func (repo *MeetingActionPlanStatusUpdateRepository) Update(meetingActionPlanStatusUpdate *models.MeetingActionPlanStatusUpdate, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(meetingActionPlanStatusUpdate).Error; err != nil {
		return err
	}
	return nil
}

func (repo *MeetingActionPlanStatusUpdateRepository) CreateDocuments(docs []*models.MeetingActionPlanStatusUpdateDocument, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	// Save all documents in a single transaction
	if err := db.Create(&docs).Error; err != nil {
		return err
	}
	return nil
}

func (repo *MeetingActionPlanStatusUpdateRepository) GetDocumentByID(id uint) (*models.MeetingActionPlanStatusUpdateDocument, error) {
	var document models.MeetingActionPlanStatusUpdateDocument
	if err := repo.db.Where("id = ?", id).First(&document).Error; err != nil {
		return nil, err
	}
	return &document, nil
}

func (repo *MeetingActionPlanStatusUpdateRepository) DeleteDocuments(docIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	// Delete documents by IDs
	if err := db.Where("id IN (?)", docIds).Delete(&models.MeetingActionPlanStatusUpdateDocument{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *MeetingActionPlanStatusUpdateRepository) GetLatestMeetingActionPlanStatusUpdate(meetingActionPlanID uint) (*meetingDTOs.MeetingActionPlanStatusUpdateListDTO, error) {

	var meetingActionPlanUpdate *meetingDTOs.MeetingActionPlanStatusUpdateListDTO

	query := repo.db.Table("meeting_action_plan_status_updates").
		Select(`
            meeting_action_plan_status_updates.*,
            meeting_action_plans.title AS meeting_action_plan_title,
            action_plan_statuses.title_en AS action_plan_status_title_en,
            action_plan_statuses.title_ar AS action_plan_status_title_ar,
            action_plan_statuses.color_code AS action_plan_status_color_code,
            created_user.id AS created_user_id,
        	created_user.first_name_en AS created_user_first_name_en,
            created_user.first_name_ar AS created_user_first_name_ar,
            created_user.last_name_en AS created_user_last_name_en,
            created_user.last_name_ar AS created_user_last_name_ar,
			created_user.initial AS created_user_initial,
        	created_user.color AS created_user_color,
        	created_user.email AS created_user_email,
        	created_user.profile_image_path AS created_user_profile_image_path,
            updated_user.id AS updated_user_id,
    		updated_user.first_name_en AS updated_user_first_name_en,
            updated_user.first_name_ar AS updated_user_first_name_ar,
            updated_user.last_name_en AS updated_user_last_name_en,
            updated_user.last_name_ar AS updated_user_last_name_ar,
			updated_user.initial AS updated_user_initial,
        	updated_user.color AS updated_user_color,
        	updated_user.email AS updated_user_email,
        	updated_user.profile_image_path AS updated_user_profile_image_path 
        `).
		Joins(`
            LEFT JOIN meeting_action_plans ON meeting_action_plans.id = meeting_action_plan_status_updates.meeting_action_plan_id
            LEFT JOIN action_plan_statuses ON action_plan_statuses.id = meeting_action_plan_status_updates.action_plan_status_id
            LEFT JOIN users AS created_user ON created_user.id = meeting_action_plan_status_updates.created_by_id
            LEFT JOIN users AS updated_user ON updated_user.id = meeting_action_plan_status_updates.updated_by_id
        `).
		Where("meeting_action_plan_status_updates.meeting_action_plan_id = ?", meetingActionPlanID).
		Order("meeting_action_plan_status_updates.created_at DESC").
		First(&meetingActionPlanUpdate)

	if query.Error != nil {
		return nil, query.Error
	}

	return meetingActionPlanUpdate, nil

}

func (repo *MeetingActionPlanStatusUpdateRepository) CheckAllActionPlansCompletedStatusCount(meetingMinuteID uint, excludeActionPlanID uint, completedStatuses []string, tx ...utility.DatabaseTransaction) (int64, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	// Count completed/terminated action plans, excluding the current action plan being updated
	var count int64
	if err := db.Model(&models.MeetingActionPlan{}).
		Where("meeting_minute_id = ?", meetingMinuteID).
		Where("id != ?", excludeActionPlanID).
		Where("action_plan_status_id NOT IN (?)", completedStatuses).
		Count(&count).Error; err != nil {
		return 0, err
	}

	// Return the count of completed/terminated action plans
	return count, nil
}
