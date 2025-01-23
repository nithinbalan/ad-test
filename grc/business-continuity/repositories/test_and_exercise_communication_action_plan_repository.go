package repositories

import (
	businessContinuityDTOs "jethur/grc/business-continuity/dtos"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type TestAndExerciseCommunicationActionPlanRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(params *businessContinuityDTOs.TestAndExerciseCommunicationActionPlanListQueryParamsDTO) ([]*businessContinuityDTOs.TestAndExerciseCommunicationActionPlanListDTO, error)
	GetByID(id uint) (*models.TestAndExerciseCommunicationActionPlan, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.TestAndExerciseCommunicationActionPlan, error)
	Create(TestAndExerciseCommunicationActionPlan *models.TestAndExerciseCommunicationActionPlan, tx ...utility.DatabaseTransaction) (uint, error)
	Update(TestAndExerciseCommunicationActionPlan *models.TestAndExerciseCommunicationActionPlan, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error
}

type TestAndExerciseCommunicationActionPlanRepository struct {
	db *gorm.DB
}

func NewTestAndExerciseCommunicationActionPlanRepository(db *gorm.DB) *TestAndExerciseCommunicationActionPlanRepository {
	return &TestAndExerciseCommunicationActionPlanRepository{
		db: db,
	}
}

// BeginTransaction starts a new transaction.
func (repo *TestAndExerciseCommunicationActionPlanRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *TestAndExerciseCommunicationActionPlanRepository) GetAll(params *businessContinuityDTOs.TestAndExerciseCommunicationActionPlanListQueryParamsDTO) ([]*businessContinuityDTOs.TestAndExerciseCommunicationActionPlanListDTO, error) {
	var testAndExerciseCommunicationActionPlans []*businessContinuityDTOs.TestAndExerciseCommunicationActionPlanListDTO
	var totalCount int64

	query := repo.db.Table("test_and_exercise_communication_action_plans").
		Select(`
            test_and_exercise_communication_action_plans.*,
			bcp_phases.title_en AS bcp_phase_title_en,
			bcp_phases.title_ar AS bcp_phase_title_ar,
			processes.title_en AS process_title_en,
			processes.title_ar AS process_title_ar,
			communication_action_plan_statuses.title_en AS communication_action_plan_status_title_en,
			communication_action_plan_statuses.title_ar AS communication_action_plan_status_title_ar,
			communication_action_plan_statuses.color_code AS communication_action_plan_status_color_code,
			responsible_user.id AS responsible_user_id,
			responsible_user.first_name_en AS responsible_user_first_name_en,
            responsible_user.first_name_ar AS responsible_user_first_name_ar,
            responsible_user.last_name_en AS responsible_user_last_name_en,
            responsible_user.last_name_ar AS responsible_user_last_name_ar,
			responsible_user.initial AS responsible_user_initial,
            responsible_user.color AS responsible_user_color,
            responsible_user.email AS responsible_user_email,
            responsible_user.profile_image_path AS responsible_user_profile_image_path,
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
			LEFT JOIN bcp_phases ON bcp_phases.id = test_and_exercise_communication_action_plans.bcp_phase_id
			LEFT JOIN processes ON processes.id = test_and_exercise_communication_action_plans.process_id
		    LEFT JOIN users AS responsible_user ON responsible_user.id = test_and_exercise_communication_action_plans.responsible_user_id
			LEFT JOIN communication_action_plan_statuses ON communication_action_plan_statuses.id = test_and_exercise_communication_action_plans.communication_action_plan_status_id
			LEFT JOIN users AS created_user ON created_user.id = test_and_exercise_communication_action_plans.created_by_id
			LEFT JOIN users AS updated_user ON updated_user.id = test_and_exercise_communication_action_plans.updated_by_id
        `)
	// Apply dynamic filters to the query
	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where(columnName+" IN (?)", values)
		}
	}
	// Calculate total count for pagination
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}
	// Apply search keyword if exists
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
			(responsible_user.first_name_en LIKE ? OR
			responsible_user.first_name_ar LIKE ? OR
			responsible_user.last_name_en LIKE ? OR
			responsible_user.last_name_ar LIKE ? OR
			responsible_user.email LIKE ? OR
            responsible_user.initial LIKE ? OR
            responsible_user.color LIKE ? OR
            responsible_user.profile_image_path LIKE ? OR
			created_user.first_name_en LIKE ? OR
			created_user.first_name_ar LIKE ? OR
			created_user.last_name_en LIKE ? OR
			created_user.last_name_ar LIKE ? OR
    		updated_user.first_name_en LIKE ? OR
			updated_user.first_name_ar LIKE ? OR
			updated_user.last_name_en LIKE ? OR
			updated_user.last_name_ar LIKE ?)`,
			keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword,
			keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword,
			keyword, keyword)
	}

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	query = query.Order(params.SortBy + " " + params.Sort)

	err := query.Find(&testAndExerciseCommunicationActionPlans).Error
	if err != nil {
		return nil, err
	}

	return testAndExerciseCommunicationActionPlans, nil
}

func (repo *TestAndExerciseCommunicationActionPlanRepository) GetByID(id uint) (*models.TestAndExerciseCommunicationActionPlan, error) {
	var testAndExerciseCommunicationActionPlan *models.TestAndExerciseCommunicationActionPlan

	result := repo.db.First(&testAndExerciseCommunicationActionPlan, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return testAndExerciseCommunicationActionPlan, nil

}

func (repo *TestAndExerciseCommunicationActionPlanRepository) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.TestAndExerciseCommunicationActionPlan, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var testAndExerciseCommunicationActionPlan *models.TestAndExerciseCommunicationActionPlan

	result := db.
		Preload("TestAndExercise").
		Preload("BcpPhase").
		Preload("Process").
		Preload("ResponsibleUser").
		Preload("CommunicationActionPlanStatus").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		First(&testAndExerciseCommunicationActionPlan, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return testAndExerciseCommunicationActionPlan, nil
}

func (repo *TestAndExerciseCommunicationActionPlanRepository) Create(TestAndExerciseCommunicationActionPlan *models.TestAndExerciseCommunicationActionPlan, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(TestAndExerciseCommunicationActionPlan).Error; err != nil {
		return 0, err
	}

	return TestAndExerciseCommunicationActionPlan.ID, nil
}

func (repo *TestAndExerciseCommunicationActionPlanRepository) Update(TestAndExerciseCommunicationActionPlan *models.TestAndExerciseCommunicationActionPlan, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(TestAndExerciseCommunicationActionPlan).Error; err != nil {
		return err
	}

	return nil
}

func (repo *TestAndExerciseCommunicationActionPlanRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	// Delete the Business Communication Plan
	if err := db.Delete(&models.TestAndExerciseCommunicationActionPlan{}, id).Error; err != nil {
		return err
	}

	return nil
}
