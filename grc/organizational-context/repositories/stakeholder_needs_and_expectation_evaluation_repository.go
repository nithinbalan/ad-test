package repositories

import (
	"jethur/grc/models"
	organizationalContextDTOs "jethur/grc/organizational-context/dtos"
	"jethur/utility"

	"gorm.io/gorm"
)

type StakeholderNeedsAndExpectationEvaluationRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(stakeholderID uint, needsAndExpectationID uint, sort string, sortBy string, keyword string) ([]*organizationalContextDTOs.StakeholderNeedsAndExpectationEvaluationListDTO, error)
	GetByID(id uint) (*models.StakeholderNeedsAndExpectationEvaluation, error)
	Create(StakeholderNeedsAndExpectationEvaluation *models.StakeholderNeedsAndExpectationEvaluation, tx ...utility.DatabaseTransaction) (uint, error)
	Update(StakeholderNeedsAndExpectationEvaluation *models.StakeholderNeedsAndExpectationEvaluation, tx ...utility.DatabaseTransaction) error

	GetLatestStakeholderNeedsAndExpectationEvaluation(needsAndExpectationID uint) (*organizationalContextDTOs.StakeholderNeedsAndExpectationEvaluationListDTO, error)
	NeedsAndExpectationUpdateStatus(needsAndExpectation *models.StakeholderNeedsAndExpectation, tx ...utility.DatabaseTransaction) error
}

type StakeholderNeedsAndExpectationEvaluationRepository struct {
	db *gorm.DB
}

func NewStakeholderNeedsAndExpectationEvaluationRepository(db *gorm.DB) *StakeholderNeedsAndExpectationEvaluationRepository {
	return &StakeholderNeedsAndExpectationEvaluationRepository{
		db: db,
	}
}

// BeginTransaction starts a new transaction.
func (repo *StakeholderNeedsAndExpectationEvaluationRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *StakeholderNeedsAndExpectationEvaluationRepository) GetAll(stakeholderID uint, needsAndExpectationID uint, sort string, sortBy string, keyword string) ([]*organizationalContextDTOs.StakeholderNeedsAndExpectationEvaluationListDTO, error) {

	var needsAndExpectationUpdates []*organizationalContextDTOs.StakeholderNeedsAndExpectationEvaluationListDTO

	query := repo.db.Table("stakeholder_needs_and_expectation_evaluations").
		Select(`
            stakeholder_needs_and_expectation_evaluations.*,
            stakeholder_needs_and_expectations.title AS stakeholder_needs_and_expectation_title,
			CONVERT(VARCHAR, stakeholder_needs_and_expectation_evaluations.evaluated_at, 23) AS evaluated_at,
            evaluation_statuses.title_en AS evaluation_status_title_en,
            evaluation_statuses.title_ar AS evaluation_status_title_ar,
            evaluation_statuses.color_code AS evaluation_status_color_code,
			evaluation_statuses.type AS evaluation_status_type,
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
            LEFT JOIN stakeholder_needs_and_expectations ON stakeholder_needs_and_expectations.id = stakeholder_needs_and_expectation_evaluations.stakeholder_needs_and_expectation_id
            LEFT JOIN evaluation_statuses ON evaluation_statuses.id = stakeholder_needs_and_expectation_evaluations.evaluation_status_id
			LEFT JOIN users AS evaluated_user ON evaluated_user.id = stakeholder_needs_and_expectation_evaluations.evaluated_by_id
            LEFT JOIN users AS created_user ON created_user.id = stakeholder_needs_and_expectation_evaluations.created_by_id
            LEFT JOIN users AS updated_user ON updated_user.id = stakeholder_needs_and_expectation_evaluations.updated_by_id
        `).
		Where("stakeholder_needs_and_expectation_evaluations.stakeholder_needs_and_expectation_id = ?", needsAndExpectationID)

	if keyword != "" {
		query = query.Where(`
            (stakeholder_needs_and_expectations.title LIKE ? OR
            stakeholder_needs_and_expectation_evaluations.comment LIKE ? OR
            evaluation_statuses.title_en LIKE ? OR
            evaluation_statuses.title_ar LIKE ?)`,
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}
	if sortBy == "" {
		sortBy = "stakeholder_needs_and_expectation_evaluations.id"
	}
	if sort != "asc" && sort != "desc" {
		sort = "asc"
	}
	query = query.Order(sortBy + " " + sort)
	result := query.Find(&needsAndExpectationUpdates)

	return needsAndExpectationUpdates, result.Error
}

func (repo *StakeholderNeedsAndExpectationEvaluationRepository) GetByID(id uint) (*models.StakeholderNeedsAndExpectationEvaluation, error) {
	var needsAndExpectationStatusUpdate *models.StakeholderNeedsAndExpectationEvaluation
	result := repo.db.First(&needsAndExpectationStatusUpdate, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return needsAndExpectationStatusUpdate, nil
}

func (repo *StakeholderNeedsAndExpectationEvaluationRepository) Create(StakeholderNeedsAndExpectationEvaluation *models.StakeholderNeedsAndExpectationEvaluation, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(StakeholderNeedsAndExpectationEvaluation).Error; err != nil {
		return 0, err
	}

	return StakeholderNeedsAndExpectationEvaluation.ID, nil
}

func (repo *StakeholderNeedsAndExpectationEvaluationRepository) Update(StakeholderNeedsAndExpectationEvaluation *models.StakeholderNeedsAndExpectationEvaluation, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(StakeholderNeedsAndExpectationEvaluation).Error; err != nil {
		return err
	}
	return nil
}

func (repo *StakeholderNeedsAndExpectationEvaluationRepository) GetLatestStakeholderNeedsAndExpectationEvaluation(needsAndExpectationID uint) (*organizationalContextDTOs.StakeholderNeedsAndExpectationEvaluationListDTO, error) {

	var needsAndExpectationUpdate *organizationalContextDTOs.StakeholderNeedsAndExpectationEvaluationListDTO

	query := repo.db.Table("stakeholder_needs_and_expectation_evaluations").
		Select(`
            stakeholder_needs_and_expectation_evaluations.*,
            stakeholder_needs_and_expectations.title AS stakeholder_needs_and_expectation_title,
            evaluation_statuses.title_en AS evaluation_status_title_en,
            evaluation_statuses.title_ar AS evaluation_status_title_ar,
            evaluation_statuses.color_code AS evaluation_status_color_code,
			evaluation_statuses.type AS evaluation_status_type,
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
            LEFT JOIN stakeholder_needs_and_expectations ON stakeholder_needs_and_expectations.id = stakeholder_needs_and_expectation_evaluations.stakeholder_needs_and_expectation_id
            LEFT JOIN evaluation_statuses ON evaluation_statuses.id = stakeholder_needs_and_expectation_evaluations.evaluation_status_id
            LEFT JOIN users AS created_user ON created_user.id = stakeholder_needs_and_expectation_evaluations.created_by_id
            LEFT JOIN users AS updated_user ON updated_user.id = stakeholder_needs_and_expectation_evaluations.updated_by_id
        `).
		Where("stakeholder_needs_and_expectation_evaluations.stakeholder_needs_and_expectation_id = ?", needsAndExpectationID).
		Order("stakeholder_needs_and_expectation_evaluations.created_at DESC").
		First(&needsAndExpectationUpdate)

	if query.Error != nil {
		return nil, query.Error
	}

	return needsAndExpectationUpdate, nil

}

func (repo *StakeholderNeedsAndExpectationEvaluationRepository) NeedsAndExpectationUpdateStatus(StakeholderNeedsAndExpectation *models.StakeholderNeedsAndExpectation, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(StakeholderNeedsAndExpectation).Error; err != nil {
		return err
	}
	return nil
}
