package repositories

import (
	"fmt"
	businessContinuityDTOs "jethur/grc/business-continuity/dtos"
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/models"
	"jethur/utility"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type TestAndExerciseRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(params *businessContinuityDTOs.TestAndExerciseListQueryParamsDTO) (*businessContinuityDTOs.TestAndExerciseListResponseDTO, error)
	GetByID(id uint) (*models.TestAndExercise, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.TestAndExercise, error)
	Create(TestAndExercise *models.TestAndExercise, tx ...utility.DatabaseTransaction) (uint, error)
	Update(TestAndExercise *models.TestAndExercise, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error
	GenerateReferenceNumber() (string, error)

	GetTestAndExerciseProcessByID(id uint) (*models.TestAndExerciseProcess, error)
	CreateTestAndExerciseProcesses(bcpProcesses []*models.TestAndExerciseProcess, tx ...utility.DatabaseTransaction) error
	DeleteTestAndExerciseProcesses(bcpProcessIds []*uint, tx ...utility.DatabaseTransaction) error

	CreateTestAndExerciseRecoveryActionPlans(bcpRecoveryActionPlans []*models.TestAndExerciseRecoveryActionPlan, tx ...utility.DatabaseTransaction) error
	DeleteTestAndExerciseRecoveryActionPlans(bcpRecoveryActionPlanIds []*uint, tx ...utility.DatabaseTransaction) error

	CreateTestAndExerciseCommunicationActionPlans(bcpCommunicationActionPlans []*models.TestAndExerciseCommunicationActionPlan, tx ...utility.DatabaseTransaction) error
	DeleteTestAndExerciseCommunicationActionPlans(bcpCommunicationActionPlanIds []*uint, tx ...utility.DatabaseTransaction) error
}

type TestAndExerciseRepository struct {
	db *gorm.DB
}

func NewTestAndExerciseRepository(db *gorm.DB) *TestAndExerciseRepository {
	return &TestAndExerciseRepository{
		db: db,
	}
}

// BeginTransaction starts a new transaction.
func (repo *TestAndExerciseRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *TestAndExerciseRepository) GetAll(params *businessContinuityDTOs.TestAndExerciseListQueryParamsDTO) (*businessContinuityDTOs.TestAndExerciseListResponseDTO, error) {
	var TestAndExercises []*businessContinuityDTOs.TestAndExerciseListDTO
	var totalCount int64

	query := repo.db.Table("test_and_exercises").
		Select(`
            test_and_exercises.*,
			CONVERT(VARCHAR, test_and_exercises.testing_date, 23) AS testing_date,
			business_functions.id AS business_function_id,
			business_functions.title AS business_function_title,
			responsible_user.id AS responsible_user_id,
			responsible_user.first_name_en AS responsible_user_first_name_en,
            responsible_user.first_name_ar AS responsible_user_first_name_ar,
            responsible_user.last_name_en AS responsible_user_last_name_en,
            responsible_user.last_name_ar AS responsible_user_last_name_ar,
			responsible_user.initial AS responsible_user_initial,
            responsible_user.color AS responsible_user_color,
            responsible_user.email AS responsible_user_email,
            responsible_user.profile_image_path AS responsible_user_profile_image_path,
			bcp_test_and_exercises.title AS bcp_test_and_exercise_title,
			test_and_exercise_statuses.title_en AS test_and_exercise_status_title_en,
			test_and_exercise_statuses.title_ar AS test_and_exercise_status_title_ar,
			test_and_exercise_statuses.color_code AS test_and_exercise_status_color_code,
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
			LEFT JOIN business_continuity_plans ON business_continuity_plans.id = test_and_exercises.business_continuity_plan_id
			LEFT JOIN business_functions ON business_functions.id = business_continuity_plans.business_function_id
			LEFT JOIN bcp_test_and_exercises ON bcp_test_and_exercises.id = test_and_exercises.bcp_test_and_exercise_id
		    LEFT JOIN users AS responsible_user ON responsible_user.id = test_and_exercises.responsible_user_id
			LEFT JOIN test_and_exercise_statuses ON test_and_exercise_statuses.id = test_and_exercises.test_and_exercise_status_id
			LEFT JOIN users AS created_user ON created_user.id = test_and_exercises.created_by_id
			LEFT JOIN users AS updated_user ON updated_user.id = test_and_exercises.updated_by_id
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
			(business_functions.title LIKE ? OR
			bcp_test_and_exercises.title LIKE ? OR
			responsible_user.first_name_en LIKE ? OR
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
			keyword, keyword, keyword, keyword)
	}

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	query = query.Order(params.SortBy + " " + params.Sort)

	offset := (params.Page - 1) * params.PageSize
	query = query.Offset(offset).Limit(params.PageSize)

	err := query.Find(&TestAndExercises).Error
	if err != nil {
		return nil, err
	}

	totalPages := 1
	totalPages = int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	response := &businessContinuityDTOs.TestAndExerciseListResponseDTO{
		Items: TestAndExercises,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}
	return response, nil
}

func (repo *TestAndExerciseRepository) GetByID(id uint) (*models.TestAndExercise, error) {
	var TestAndExercise *models.TestAndExercise

	result := repo.db.First(&TestAndExercise, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return TestAndExercise, nil

}

func (repo *TestAndExerciseRepository) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.TestAndExercise, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var TestAndExercise *models.TestAndExercise

	result := db.
		Preload("BusinessContinuityPlan").
		Preload("BusinessContinuityPlan.BusinessFunction").
		Preload("ResponsibleUser").
		Preload("BcpTestAndExercise").
		Preload("TestAndExerciseProcesses").
		Preload("TestAndExerciseProcesses.Process").
		Preload("TestAndExerciseProcesses.PlannedRtoUnit").
		Preload("TestAndExerciseProcesses.PlannedMtpdUnit").
		Preload("TestAndExerciseProcesses.PlannedRpoUnit").
		Preload("TestAndExerciseProcesses.ActualRtoUnit").
		Preload("TestAndExerciseProcesses.ActualMtpdUnit").
		Preload("TestAndExerciseProcesses.ActualRpoUnit").
		Preload("TestAndExerciseProcesses.Timeliness").
		Preload("TestAndExerciseRecoveryActionPlans").
		Preload("TestAndExerciseCommunicationActionPlans").
		Preload("TestAndExerciseStatus").
		Preload("TestAndExerciseTestResult").
		Preload("TestAndExerciseTestResult.CommunicationFeedbackOption").
		Preload("TestAndExerciseTestResult.RecoveryFeedbackOption").
		Preload("TestAndExerciseTestResult.TestResultOption").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		First(&TestAndExercise, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return TestAndExercise, nil
}

func (repo *TestAndExerciseRepository) Create(TestAndExercise *models.TestAndExercise, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(TestAndExercise).Error; err != nil {
		return 0, err
	}

	return TestAndExercise.ID, nil
}

func (repo *TestAndExerciseRepository) Update(TestAndExercise *models.TestAndExercise, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(TestAndExercise).Error; err != nil {
		return err
	}

	return nil
}

func (repo *TestAndExerciseRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	// Delete the Business Communication Plan
	if err := db.Delete(&models.TestAndExercise{}, id).Error; err != nil {
		return err
	}

	return nil
}

func (repo *TestAndExerciseRepository) CreateTestAndExerciseProcesses(processes []*models.TestAndExerciseProcess, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(&processes).Error; err != nil {
		return err
	}
	return nil
}

func (repo *TestAndExerciseRepository) DeleteTestAndExerciseProcesses(processIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	// Delete processes by IDs
	if err := db.Where("id IN (?)", processIds).Delete(&models.TestAndExerciseProcess{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *TestAndExerciseRepository) CreateTestAndExerciseRecoveryActionPlans(recovertActionPlans []*models.TestAndExerciseRecoveryActionPlan, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(&recovertActionPlans).Error; err != nil {
		return err
	}
	return nil
}

func (repo *TestAndExerciseRepository) DeleteTestAndExerciseRecoveryActionPlans(recovertActionPlanIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	// Delete recovery action plan by IDs
	if err := db.Where("id IN (?)", recovertActionPlanIds).Delete(&models.TestAndExerciseRecoveryActionPlan{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *TestAndExerciseRepository) CreateTestAndExerciseCommunicationActionPlans(communicationActionPlans []*models.TestAndExerciseCommunicationActionPlan, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(&communicationActionPlans).Error; err != nil {
		return err
	}
	return nil
}

func (repo *TestAndExerciseRepository) DeleteTestAndExerciseCommunicationActionPlans(communicationActionPlanIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	// Delete communication action plan by IDs
	if err := db.Where("id IN (?)", communicationActionPlanIds).Delete(&models.TestAndExerciseCommunicationActionPlan{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *TestAndExerciseRepository) GenerateReferenceNumber() (string, error) {
	var lastReferenceNumber string
	err := repo.db.Model(&models.TestAndExercise{}).
		Select("reference_number").
		Order("reference_number desc").
		Limit(1).
		Pluck("reference_number", &lastReferenceNumber).
		Error

	if err != nil {
		return "", err
	}

	lastDigits := "0000"
	if lastReferenceNumber != "" && len(lastReferenceNumber) >= 4 {
		lastDigits = lastReferenceNumber[len(lastReferenceNumber)-4:]
	}

	currentYear := time.Now().Year()
	lastReferenceID, err := strconv.ParseUint(lastDigits, 10, 64)
	if err != nil {
		return "", err
	}
	newReferenceNumber := fmt.Sprintf("%s-%d-%04d", "BCM-TEST", currentYear, lastReferenceID+1)

	return newReferenceNumber, nil

}

func (repo *TestAndExerciseRepository) GetTestAndExerciseProcessByID(id uint) (*models.TestAndExerciseProcess, error) {
	var TestAndExerciseProcess *models.TestAndExerciseProcess

	result := repo.db.First(&TestAndExerciseProcess, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return TestAndExerciseProcess, nil

}
