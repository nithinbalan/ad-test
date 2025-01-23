package repositories

import (
	businessContinuityDTOs "jethur/grc/business-continuity/dtos"
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type TestAndExerciseRecoveryActionPlanRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(testAndExerciseID uint, params *businessContinuityDTOs.TestAndExerciseRecoveryActionPlanQueryParamsDTO) (*businessContinuityDTOs.TestAndExerciseRecoveryActionPlanListResponseDTO, error)
	GetByID(id uint) (*models.TestAndExerciseRecoveryActionPlan, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.TestAndExerciseRecoveryActionPlan, error)
	Create(actionPlan *models.TestAndExerciseRecoveryActionPlan, tx ...utility.DatabaseTransaction) (uint, error)
	Update(actionPlan *models.TestAndExerciseRecoveryActionPlan, tx ...utility.DatabaseTransaction) error
	Delete(actionPlanID uint, tx ...utility.DatabaseTransaction) error
}

type TestAndExerciseRecoveryActionPlanRepository struct {
	db *gorm.DB
}

func NewTestAndExerciseRecoveryActionPlanRepository(db *gorm.DB) *TestAndExerciseRecoveryActionPlanRepository {
	return &TestAndExerciseRecoveryActionPlanRepository{
		db: db,
	}
}

func (repo *TestAndExerciseRecoveryActionPlanRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *TestAndExerciseRecoveryActionPlanRepository) GetAll(testAndExerciseID uint, params *businessContinuityDTOs.TestAndExerciseRecoveryActionPlanQueryParamsDTO) (*businessContinuityDTOs.TestAndExerciseRecoveryActionPlanListResponseDTO, error) {
	var testAndExerciseRecoveryActionPlans []*businessContinuityDTOs.TestAndExerciseRecoveryActionPlanListDTO
	var totalCount int64

	query := repo.db.Table("test_and_exercise_recovery_action_plans").
		Select(`
            test_and_exercise_recovery_action_plans.id,
            test_and_exercise_recovery_action_plans.test_and_exercise_id,
            test_and_exercise_recovery_action_plans.bcp_phase_id,
            test_and_exercise_recovery_action_plans.process_id,
            processes.title_en AS process_title_en, 
            processes.title_ar AS process_title_ar, 
            test_and_exercise_recovery_action_plans.severity_level_id,
            severity_levels.title_en AS severity_level_title_en, 
            severity_levels.title_ar AS severity_level_title_ar, 
            test_and_exercise_recovery_action_plans.title,
            test_and_exercise_recovery_action_plans.responsible_user_id,
            test_and_exercise_recovery_action_plans.responsible_unit_id,
            organization_units.title_en AS responsible_unit_title_en, 
            organization_units.title_ar AS responsible_unit_title_ar,
            responsible_user.first_name_en AS responsible_user_first_name_en,
            responsible_user.last_name_en AS responsible_user_last_name_en, 
            responsible_user.first_name_ar AS responsible_user_first_name_ar,
            responsible_user.last_name_ar AS responsible_user_last_name_ar,
            responsible_user.initial AS responsible_user_initial, 
            responsible_user.email AS responsible_user_email, 
            responsible_user.color AS responsible_user_color, 
            responsible_user.profile_image_path AS responsible_user_profile_image_path, 
            test_and_exercise_recovery_action_plans.target_completion,
            test_and_exercise_recovery_action_plans.target_completion_unit_id,
            time_units.title_en AS target_completion_unit_en,
            time_units.title_ar AS target_completion_unit_ar,
            test_and_exercise_recovery_action_plans.action_plan_status_id,
            action_plan_statuses.title_en AS action_plan_status_title_en, 
            action_plan_statuses.title_ar AS action_plan_status_title_ar, 
            test_and_exercise_recovery_action_plans.created_by_id,
            created_user.first_name_en AS created_user_first_name_en,
            created_user.last_name_en AS created_user_last_name_en, 
            created_user.first_name_ar AS created_user_first_name_ar,
            created_user.last_name_ar AS created_user_last_name_ar,
            created_user.profile_image_path AS created_user_profile_image_path, 
            created_user.email AS created_user_email, 
            created_user.initial AS created_user_initial, 
            created_user.color AS created_user_color, 
            test_and_exercise_recovery_action_plans.created_at,
            test_and_exercise_recovery_action_plans.updated_by_id,
            updated_user.first_name_en AS updated_user_first_name_en, 
            updated_user.last_name_en AS updated_user_last_name_en,
            updated_user.first_name_ar AS updated_user_first_name_ar,
            updated_user.last_name_ar AS updated_user_last_name_ar,
            updated_user.profile_image_path AS updated_user_profile_image_path, 
            updated_user.email AS updated_user_email,
            updated_user.initial AS updated_user_initial,
            updated_user.color AS updated_user_color, 
            test_and_exercise_recovery_action_plans.updated_at
        `).
		Joins(`LEFT JOIN processes ON processes.id = test_and_exercise_recovery_action_plans.process_id`).
		Joins(`LEFT JOIN severity_levels ON severity_levels.id = test_and_exercise_recovery_action_plans.severity_level_id`).
		Joins(`LEFT JOIN users AS responsible_user ON responsible_user.id = test_and_exercise_recovery_action_plans.responsible_user_id`).
		Joins(`LEFT JOIN organization_units ON organization_units.id = test_and_exercise_recovery_action_plans.responsible_unit_id`).
		Joins(`LEFT JOIN time_units ON time_units.id = test_and_exercise_recovery_action_plans.target_completion_unit_id`).
		Joins(`LEFT JOIN action_plan_statuses ON action_plan_statuses.id = test_and_exercise_recovery_action_plans.action_plan_status_id`).
		Joins(`LEFT JOIN users AS created_user ON created_user.id = test_and_exercise_recovery_action_plans.created_by_id`).
		Joins(`LEFT JOIN users AS updated_user ON updated_user.id = test_and_exercise_recovery_action_plans.updated_by_id`).
		Where("test_and_exercise_recovery_action_plans.test_and_exercise_id = ?", testAndExerciseID)

	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where(columnName+" IN (?)", values)
		}
	}

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
			(processes.title_en LIKE ? OR
			processes.title_ar LIKE ? OR
			organization_units.title_en LIKE ? OR
			organization_units.title_ar LIKE ? OR
			responsible_user.first_name_en LIKE ? OR
			responsible_user.last_name_en LIKE ? OR
			responsible_user.first_name_ar LIKE ? OR
			responsible_user.last_name_ar LIKE ? OR
			responsible_user.email LIKE ? OR
			action_plan_statuses.title_en LIKE ? OR
			action_plan_statuses.title_ar LIKE ? OR
			test_and_exercise_recovery_action_plans.title LIKE ? OR
			time_units.title_en LIKE ? OR
			time_units.title_ar LIKE ? OR
			created_user.first_name_en LIKE ? OR
			created_user.last_name_en LIKE ? OR
			created_user.first_name_ar LIKE ? OR
			created_user.last_name_ar LIKE ? OR
			created_user.email LIKE ? OR
			updated_user.first_name_en LIKE ? OR
			updated_user.last_name_en LIKE ? OR
			updated_user.first_name_ar LIKE ? OR
			updated_user.last_name_ar LIKE ? OR
			updated_user.email LIKE ?)`,
			keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword,
			keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword,
			keyword, keyword, keyword, keyword, keyword, keyword, keyword)
	}

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	query = query.Order(params.SortBy + " " + params.Sort)

	offset := (params.Page - 1) * params.PageSize
	query = query.Offset(offset).Limit(params.PageSize)

	err := query.Find(&testAndExerciseRecoveryActionPlans).Error
	if err != nil {
		return nil, err
	}

	totalPages := 1
	totalPages = int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	response := &businessContinuityDTOs.TestAndExerciseRecoveryActionPlanListResponseDTO{
		Items: testAndExerciseRecoveryActionPlans,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}
	return response, nil
}

func (repo *TestAndExerciseRecoveryActionPlanRepository) GetByID(id uint) (*models.TestAndExerciseRecoveryActionPlan, error) {
	var actionPlan models.TestAndExerciseRecoveryActionPlan
	if err := repo.db.First(&actionPlan, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &actionPlan, nil
}

func (repo *TestAndExerciseRecoveryActionPlanRepository) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.TestAndExerciseRecoveryActionPlan, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var actionPlan *models.TestAndExerciseRecoveryActionPlan
	result := db.
		Preload("TestAndExercise").
		Preload("BcpPhase").
		Preload("Process").
		Preload("SeverityLevel").
		Preload("TargetCompletionUnit").
		Preload("ActionPlanStatus").
		Preload("ResponsibleUser").
		Preload("ResponsibleUnit").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		Where("id = ?", id).
		First(&actionPlan)

	if result.Error != nil {
		return nil, result.Error
	}

	return actionPlan, nil
}

func (repo *TestAndExerciseRecoveryActionPlanRepository) Create(actionPlan *models.TestAndExerciseRecoveryActionPlan, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Create(actionPlan).Error; err != nil {
		return 0, err
	}

	return actionPlan.ID, nil
}

func (repo *TestAndExerciseRecoveryActionPlanRepository) Update(actionPlan *models.TestAndExerciseRecoveryActionPlan, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Model(&models.TestAndExerciseRecoveryActionPlan{}).Where("id = ?", actionPlan.ID).Updates(actionPlan).Error; err != nil {
		return err
	}

	return nil
}

func (repo *TestAndExerciseRecoveryActionPlanRepository) Delete(actionPlanID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("id = ?", actionPlanID).Delete(&models.TestAndExerciseRecoveryActionPlan{}).Error; err != nil {
		return err
	}

	return nil
}
