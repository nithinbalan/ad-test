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

type BusinessFunctionRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(params *businessContinuityDTOs.BusinessFunctionListQueryParamsDTO) (*businessContinuityDTOs.BusinessFunctionListResponseDTO, error)
	GetByID(id uint) (*models.BusinessFunction, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.BusinessFunction, error)
	Create(businessFunction *models.BusinessFunction, tx ...utility.DatabaseTransaction) (uint, error)
	Update(businessFunction *models.BusinessFunction, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error
	GenerateReferenceNumber() (string, error)

	CreateBusinessFunctionProcesses(processes []*models.BusinessFunctionProcess, tx ...utility.DatabaseTransaction) error
	DeleteBusinessFunctionProcesses(processIds []*uint, tx ...utility.DatabaseTransaction) error
	GetTestAndExercisesByBcpID(bcpID uint) ([]*models.TestAndExercise, error)
	GetBcpTestAndExercisesByBcpID(bcpID uint) ([]*models.BcpTestAndExercise, error)
	GetProcessWithCriticality(id uint, tx ...utility.DatabaseTransaction) (*models.Process, error)
}

type BusinessFunctionRepository struct {
	db *gorm.DB
}

func NewBusinessFunctionRepository(db *gorm.DB) *BusinessFunctionRepository {
	return &BusinessFunctionRepository{
		db: db,
	}
}

// BeginTransaction starts a new transaction.
func (repo *BusinessFunctionRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *BusinessFunctionRepository) GetAll(params *businessContinuityDTOs.BusinessFunctionListQueryParamsDTO) (*businessContinuityDTOs.BusinessFunctionListResponseDTO, error) {
	var businessFunctions []*businessContinuityDTOs.BusinessFunctionListDTO
	var totalCount int64

	query := repo.db.Table("business_functions").
		Select(`
            business_functions.*,
			severity_levels.title_en AS severity_level_title_en,
			severity_levels.title_ar AS severity_level_title_ar,
			severity_levels.color_code AS severity_level_color_code,
			review_frequencies.title_en AS review_frequency_title_en,
			review_frequencies.title_ar AS review_frequency_title_ar,
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
			LEFT JOIN severity_levels ON severity_levels.id = business_functions.severity_level_id
			LEFT JOIN review_frequencies ON review_frequencies.id = business_functions.review_frequency_id
		    LEFT JOIN users AS responsible_user ON responsible_user.id = business_functions.responsible_user_id
			LEFT JOIN users AS created_user ON created_user.id = business_functions.created_by_id
			LEFT JOIN users AS updated_user ON updated_user.id = business_functions.updated_by_id
        `)
	// Apply dynamic filters to the query
	joinedProcessTables := false

	for columnName, values := range params.Filters {
		if len(values) > 0 {
			// 🔥 Join only once, not on every filter
			if (columnName == "process_id" || columnName == "responsible_unit_id" || columnName == "process_criticality_level_id") && !joinedProcessTables {
				query = query.Joins(`LEFT JOIN business_function_processes ON business_function_processes.business_function_id = business_functions.id`)
				query = query.Joins(`LEFT JOIN processes ON processes.id = business_function_processes.process_id`)
				query = query.Joins(`LEFT JOIN organization_units ON organization_units.id = processes.responsible_unit_id`)
				query = query.Joins(`LEFT JOIN process_criticality_levels ON process_criticality_levels.id = processes.process_criticality_level_id`)
				joinedProcessTables = true // ✅ mark as joined
			}

			query = query.Where(columnName+" IN (?)", values)
		}
	}

	// Calculate total count for pagination
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}
	// Apply quick filters
	for _, filter := range params.QuickFilter {
		switch filter {
		case "is_not_bcp":
			query = query.Where("business_functions.bcp_status = ?", false)
		case "is_bcp":
			query = query.Where("business_functions.bcp_status = ?", true)
		case "high", "medium", "low", "is_not_bcp_critical", "is_bcp_non_critical":
			// ⚡ These quick filters depend on process criticality, JOIN needed
			query = query.Joins(`LEFT JOIN business_function_processes ON business_function_processes.business_function_id = business_functions.id`)
			query = query.Joins(`LEFT JOIN processes ON processes.id = business_function_processes.process_id`)
			query = query.Joins(`LEFT JOIN organization_units ON organization_units.id = processes.responsible_unit_id`)
			query = query.Joins(`LEFT JOIN process_criticality_levels ON process_criticality_levels.id = processes.process_criticality_level_id`)

			switch filter {
			case "high":
				query = query.Where("process_criticality_levels.type = ?", "high")
			case "medium":
				query = query.Where("process_criticality_levels.type = ?", "medium")
			case "low":
				query = query.Where("process_criticality_levels.type = ?", "low")
			case "is_not_bcp_critical":
				query = query.Where("process_criticality_levels.type = ? AND business_functions.bcp_status = ?", "high", false)
			case "is_bcp_non_critical":
				query = query.Where("process_criticality_levels.type = ? AND business_functions.bcp_status = ?", "low", true)
			}
		case "all":
			continue
		}
	}
	// Apply search keyword if exists
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
			(business_functions.title LIKE ? OR
			business_functions.description LIKE ? OR
			business_functions.reference_number LIKE ? OR
			severity_levels.title_en LIKE ? OR
			severity_levels.title_ar LIKE ? OR
			review_frequencies.title_en LIKE ? OR
			review_frequencies.title_ar LIKE ? OR
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
			keyword, keyword, keyword, keyword, keyword, keyword, keyword)
	}

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	query = query.Order(params.SortBy + " " + params.Sort)

	err := query.Find(&businessFunctions).Error
	if err != nil {
		return nil, err
	}

	totalPages := 1

	response := &businessContinuityDTOs.BusinessFunctionListResponseDTO{
		Items: businessFunctions,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}
	return response, nil
}

func (repo *BusinessFunctionRepository) GetByID(id uint) (*models.BusinessFunction, error) {
	var businessFunction *models.BusinessFunction

	result := repo.db.First(&businessFunction, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return businessFunction, nil

}

func (repo *BusinessFunctionRepository) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.BusinessFunction, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var businessFunction *models.BusinessFunction

	result := db.
		Preload("SeverityLevel").
		Preload("Processes").
		Preload("Processes.Process").
		Preload("Processes.Process.MtpdUnit").
		Preload("Processes.Process.RpoUnit").
		Preload("Processes.Process.RtoUnit").
		Preload("Processes.Process.ProcessGroup").
		Preload("Processes.Process.RiskLevel").
		Preload("Processes.Process.ResponsibleUnit").
		Preload("Processes.Process.ResponsibleUser").
		Preload("Processes.Process.LifecycleStatus").
		Preload("Processes.Process.RiskMatrixRating").
		Preload("Processes.Process.RiskMatrixRatingValue").
		Preload("ResponsibleUser", models.WithPrimaryOrganizationUnit()).
		Preload("ReviewFrequency").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		Preload("LatestBusinessContinuityPlan").
		Preload("LatestBusinessContinuityPlan.MtpdUnit").
		Preload("LatestBusinessContinuityPlan.RpoUnit").
		Preload("LatestBusinessContinuityPlan.RtoUnit").
		Preload("LatestBusinessContinuityPlan.Processes").
		Preload("LatestBusinessContinuityPlan.Processes.Process").
		Preload("LatestBusinessContinuityPlan.Processes.MtpdUnit").
		Preload("LatestBusinessContinuityPlan.Processes.RpoUnit").
		Preload("LatestBusinessContinuityPlan.Processes.RtoUnit").
		Preload("LatestBusinessContinuityPlan.Processes.RiskLevel").
		First(&businessFunction, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return businessFunction, nil
}

func (repo *BusinessFunctionRepository) Create(businessFunction *models.BusinessFunction, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(businessFunction).Error; err != nil {
		return 0, err
	}

	return businessFunction.ID, nil
}

func (repo *BusinessFunctionRepository) Update(businessFunction *models.BusinessFunction, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(businessFunction).Error; err != nil {
		return err
	}

	return nil
}

func (repo *BusinessFunctionRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	// Delete the Business Function
	if err := db.Delete(&models.BusinessFunction{}, id).Error; err != nil {
		return err
	}

	return nil
}

func (repo *BusinessFunctionRepository) GenerateReferenceNumber() (string, error) {
	var lastReferenceNumber string
	err := repo.db.Model(&models.BusinessFunction{}).
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
	newReferenceNumber := fmt.Sprintf("%s-%d-%04d", "BF", currentYear, lastReferenceID+1)

	return newReferenceNumber, nil

}

func (repo *BusinessFunctionRepository) CreateBusinessFunctionProcesses(processes []*models.BusinessFunctionProcess, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(&processes).Error; err != nil {
		return err
	}
	return nil
}

func (repo *BusinessFunctionRepository) DeleteBusinessFunctionProcesses(processIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	// Delete processes by IDs
	if err := db.Where("id IN (?)", processIds).Delete(&models.BusinessFunctionProcess{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *BusinessFunctionRepository) GetTestAndExercisesByBcpID(bcpID uint) ([]*models.TestAndExercise, error) {
	var testAndExercises []*models.TestAndExercise

	result := repo.db.
		Preload("TestAndExerciseTestResult").
		Preload("TestAndExerciseTestResult.TestResultOption").
		Where("business_continuity_plan_id = ?", bcpID).
		Find(&testAndExercises)

	if result.Error != nil {
		return nil, result.Error
	}

	return testAndExercises, nil
}

func (repo *BusinessFunctionRepository) GetBcpTestAndExercisesByBcpID(bcpID uint) ([]*models.BcpTestAndExercise, error) {
	var bcpTestAndExercises []*models.BcpTestAndExercise

	result := repo.db.
		Where("business_continuity_plan_id = ?", bcpID).
		Find(&bcpTestAndExercises)

	if result.Error != nil {
		return nil, result.Error
	}

	return bcpTestAndExercises, nil
}
func (repo *BusinessFunctionRepository) GetProcessWithCriticality(id uint, tx ...utility.DatabaseTransaction) (*models.Process, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var process models.Process
	err := db.Preload("ProcessCriticalityLevel").First(&process, id).Error
	if err != nil {
		return nil, err
	}
	return &process, nil
}
