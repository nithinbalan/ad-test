package repositories

import (
	"jethur/grc/models"
	organizationalContextDTOs "jethur/grc/organizational-context/dtos"
	"jethur/utility"

	"gorm.io/gorm"
)

type StakeholderEngagementEvaluationRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(stakeholderID uint, needsAndExpectationID uint, sort string, sortBy string, keyword string) ([]*organizationalContextDTOs.StakeholderEngagementEvaluationListDTO, error)
	GetByID(id uint) (*models.StakeholderEngagementEvaluation, error)
	Create(StakeholderEngagementEvaluation *models.StakeholderEngagementEvaluation, tx ...utility.DatabaseTransaction) (uint, error)
	Update(StakeholderEngagementEvaluation *models.StakeholderEngagementEvaluation, tx ...utility.DatabaseTransaction) error

	GetLatestStakeholderEngagementEvaluation(needsAndExpectationID uint) (*organizationalContextDTOs.StakeholderEngagementEvaluationListDTO, error)
	EngagementUpdateStatus(needsAndExpectation *models.StakeholderEngagement, tx ...utility.DatabaseTransaction) error
}

type StakeholderEngagementEvaluationRepository struct {
	db *gorm.DB
}

func NewStakeholderEngagementEvaluationRepository(db *gorm.DB) *StakeholderEngagementEvaluationRepository {
	return &StakeholderEngagementEvaluationRepository{
		db: db,
	}
}

// BeginTransaction starts a new transaction.
func (repo *StakeholderEngagementEvaluationRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *StakeholderEngagementEvaluationRepository) GetAll(stakeholderID uint, needsAndExpectationID uint, sort string, sortBy string, keyword string) ([]*organizationalContextDTOs.StakeholderEngagementEvaluationListDTO, error) {

	var needsAndExpectationUpdates []*organizationalContextDTOs.StakeholderEngagementEvaluationListDTO

	query := repo.db.Table("stakeholder_engagement_evaluations").
		Select(`
            stakeholder_engagement_evaluations.*,
            stakeholder_engagements.title AS stakeholder_engagement_title,
			CONVERT(VARCHAR, stakeholder_engagement_evaluations.evaluated_at, 23) AS evaluated_at,
			evaluated_user.id AS evaluated_user_id,
        	evaluated_user.first_name_en AS evaluated_user_first_name_en,
        	evaluated_user.first_name_ar AS evaluated_user_first_name_ar,
        	evaluated_user.last_name_en AS evaluated_user_last_name_en,
        	evaluated_user.last_name_ar AS evaluated_user_last_name_ar,   
        	evaluated_user.initial AS evaluated_user_initial,
        	evaluated_user.color AS evaluated_user_color,
        	evaluated_user.email AS evaluated_user_email,
        	evaluated_user.profile_image_path AS evaluated_user_profile_image_path,
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
            LEFT JOIN stakeholder_engagements ON stakeholder_engagements.id = stakeholder_engagement_evaluations.stakeholder_engagement_id
            LEFT JOIN users AS created_user ON created_user.id = stakeholder_engagement_evaluations.created_by_id
			LEFT JOIN users AS evaluated_user ON evaluated_user.id = stakeholder_engagement_evaluations.evaluated_by_id
            LEFT JOIN users AS updated_user ON updated_user.id = stakeholder_engagement_evaluations.updated_by_id
        `).
		Where("stakeholder_engagement_evaluations.stakeholder_engagement_id = ?", needsAndExpectationID)

	if keyword != "" {
		query = query.Where(`
            (stakeholder_engagements.title LIKE ? OR
            stakeholder_engagement_evaluations.comment LIKE ?)`,
			"%"+keyword+"%", "%"+keyword+"%")
	}
	if sortBy == "" {
		sortBy = "stakeholder_engagement_evaluations.id"
	}
	if sort != "asc" && sort != "desc" {
		sort = "asc"
	}
	query = query.Order(sortBy + " " + sort)
	result := query.Find(&needsAndExpectationUpdates)

	return needsAndExpectationUpdates, result.Error
}

func (repo *StakeholderEngagementEvaluationRepository) GetByID(id uint) (*models.StakeholderEngagementEvaluation, error) {
	var needsAndExpectationStatusUpdate *models.StakeholderEngagementEvaluation
	result := repo.db.First(&needsAndExpectationStatusUpdate, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return needsAndExpectationStatusUpdate, nil
}

func (repo *StakeholderEngagementEvaluationRepository) Create(StakeholderEngagementEvaluation *models.StakeholderEngagementEvaluation, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(StakeholderEngagementEvaluation).Error; err != nil {
		return 0, err
	}

	return StakeholderEngagementEvaluation.ID, nil
}

func (repo *StakeholderEngagementEvaluationRepository) Update(StakeholderEngagementEvaluation *models.StakeholderEngagementEvaluation, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(StakeholderEngagementEvaluation).Error; err != nil {
		return err
	}
	return nil
}

func (repo *StakeholderEngagementEvaluationRepository) GetLatestStakeholderEngagementEvaluation(needsAndExpectationID uint) (*organizationalContextDTOs.StakeholderEngagementEvaluationListDTO, error) {

	var needsAndExpectationUpdate *organizationalContextDTOs.StakeholderEngagementEvaluationListDTO

	query := repo.db.Table("stakeholder_engagement_evaluations").
		Select(`
            stakeholder_engagement_evaluations.*,
            stakeholder_engagements.title AS stakeholder_engagement_title,
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
            LEFT JOIN stakeholder_engagements ON stakeholder_engagements.id = stakeholder_engagement_evaluations.stakeholder_engagement_id
            LEFT JOIN users AS created_user ON created_user.id = stakeholder_engagement_evaluations.created_by_id
            LEFT JOIN users AS updated_user ON updated_user.id = stakeholder_engagement_evaluations.updated_by_id
        `).
		Where("stakeholder_engagement_evaluations.stakeholder_engagement_id = ?", needsAndExpectationID).
		Order("stakeholder_engagement_evaluations.created_at DESC").
		First(&needsAndExpectationUpdate)

	if query.Error != nil {
		return nil, query.Error
	}

	return needsAndExpectationUpdate, nil

}

func (repo *StakeholderEngagementEvaluationRepository) EngagementUpdateStatus(StakeholderEngagement *models.StakeholderEngagement, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(StakeholderEngagement).Error; err != nil {
		return err
	}
	return nil
}
