package repositories

import (
	"errors"
	"fmt"
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/models"
	penetrationTestDTOs "jethur/grc/penetration-testing/dtos"
	"jethur/utility"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type PenetrationTestRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(params *penetrationTestDTOs.PenetrationTestListQueryParamsDTO) (*penetrationTestDTOs.PenetrationTestListResponseDTO, error)
	GetByID(id uint, tx ...utility.DatabaseTransaction) (*models.PenetrationTest, error)
	Create(penetrationTest *models.PenetrationTest, tx ...utility.DatabaseTransaction) (uint, error)
	Update(penetrationTest *models.PenetrationTest, tx ...utility.DatabaseTransaction) error
	UpdateStatus(id uint, statusID uint, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error
	DeletePenetrationTestDetails(id uint, tx ...utility.DatabaseTransaction) error
	GenerateReferenceNumber(tx ...utility.DatabaseTransaction) (string, error)
	GetSeverityCounts(penetrationTestID uint) (*PenetrationTestSeverityCounts, error)
	GetSeverityBreakdown(penetrationTestID uint) (*penetrationTestDTOs.PenetrationTestSeverityBreakdownDTO, error)
}

type PenetrationTestSeverityCounts struct {
	CriticalCount int
	HighCount     int
	MediumCount   int
	LowCount      int
	TotalCount    int
}

type PenetrationTestRepository struct {
	db *gorm.DB
}

func NewPenetrationTestRepository(db *gorm.DB) *PenetrationTestRepository {
	return &PenetrationTestRepository{
		db: db,
	}
}

func (repo *PenetrationTestRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *PenetrationTestRepository) GetAll(params *penetrationTestDTOs.PenetrationTestListQueryParamsDTO) (*penetrationTestDTOs.PenetrationTestListResponseDTO, error) {
	var penetrationTests []*penetrationTestDTOs.PenetrationTestListDTO
	var totalCount int64

	query := repo.db.Table("penetration_tests").
		Select(`
			penetration_tests.id,
			penetration_tests.reference_number,
			penetration_tests.system_application_name,
			penetration_tests.penetration_test_type_id,
			penetration_test_types.title_en AS penetration_test_type_title_en,
			penetration_test_types.title_ar AS penetration_test_type_title_ar,
			penetration_tests.penetration_test_plan_id,
			penetration_test_plans.title_en AS penetration_test_plan_title_en,
			penetration_test_plans.title_ar AS penetration_test_plan_title_ar,
			penetration_tests.year,
			penetration_tests.period,
			penetration_tests.penetration_test_status_id,
			penetration_test_statuses.title_en AS status_title_en,
			penetration_test_statuses.title_ar AS status_title_ar,
			penetration_test_statuses.color_code AS status_color_code,
			penetration_test_statuses.type AS status_type,
			penetration_tests.penetration_test_classification_id,
			penetration_test_classifications.title_en AS penetration_test_classification_title_en,
			penetration_test_classifications.title_ar AS penetration_test_classification_title_ar,
			penetration_tests.penetration_test_executor_id,
			penetration_test_executors.title_en AS penetration_test_executor_title_en,
			penetration_test_executors.title_ar AS penetration_test_executor_title_ar,
			penetration_tests.cs_responsible_employee_id,
			cs_responsible_employee.first_name_en AS cs_responsible_employee_first_name_en,
			cs_responsible_employee.first_name_ar AS cs_responsible_employee_first_name_ar,
			cs_responsible_employee.last_name_en AS cs_responsible_employee_last_name_en,
			cs_responsible_employee.last_name_ar AS cs_responsible_employee_last_name_ar,
			cs_responsible_employee.initial AS cs_responsible_employee_initial,
			cs_responsible_employee.color AS cs_responsible_employee_color,
			cs_responsible_employee.email AS cs_responsible_employee_email,
			cs_responsible_employee.profile_image_path AS cs_responsible_employee_profile_image_path,
			environment_preparations.environment_readiness_id,
			executions.test_start_date,
			executions.test_end_date,
			executions.execution_status_id,
			execution_statuses.title_en AS execution_status_title_en,
			execution_statuses.title_ar AS execution_status_title_ar,
			executions.security_result_id,
			security_results.title_en AS security_result_title_en,
			security_results.title_ar AS security_result_title_ar,
			management_sharings.results_shared_with_management,
			(
				SELECT COUNT(*)
				FROM penetration_test_findings
				INNER JOIN severity_levels ON severity_levels.id = penetration_test_findings.severity_level_id
				WHERE penetration_test_findings.penetration_test_id = penetration_tests.id
				AND severity_levels.type = 'critical'
			) AS critical_count,
			(
				SELECT COUNT(*)
				FROM penetration_test_findings
				INNER JOIN severity_levels ON severity_levels.id = penetration_test_findings.severity_level_id
				WHERE penetration_test_findings.penetration_test_id = penetration_tests.id
				AND severity_levels.type = 'high'
			) AS high_count,
			(
				SELECT COUNT(*)
				FROM penetration_test_findings
				INNER JOIN severity_levels ON severity_levels.id = penetration_test_findings.severity_level_id
				WHERE penetration_test_findings.penetration_test_id = penetration_tests.id
				AND severity_levels.type = 'medium'
			) AS medium_count,
			(
				SELECT COUNT(*)
				FROM penetration_test_findings
				INNER JOIN severity_levels ON severity_levels.id = penetration_test_findings.severity_level_id
				WHERE penetration_test_findings.penetration_test_id = penetration_tests.id
				AND severity_levels.type = 'low'
			) AS low_count,
			(
				SELECT COUNT(*)
				FROM penetration_test_findings
				WHERE penetration_test_findings.penetration_test_id = penetration_tests.id
			) AS total_count
		`).
		Joins("LEFT JOIN penetration_test_environment_preparations AS environment_preparations ON environment_preparations.penetration_test_id = penetration_tests.id").
		Joins("LEFT JOIN penetration_test_executions AS executions ON executions.penetration_test_id = penetration_tests.id").
		Joins("LEFT JOIN penetration_test_execution_statuses AS execution_statuses ON execution_statuses.id = executions.execution_status_id").
		Joins("LEFT JOIN penetration_test_security_results AS security_results ON security_results.id = executions.security_result_id").
		Joins("LEFT JOIN penetration_test_management_sharings AS management_sharings ON management_sharings.penetration_test_id = penetration_tests.id").
		Joins("LEFT JOIN penetration_test_types ON penetration_test_types.id = penetration_tests.penetration_test_type_id").
		Joins("LEFT JOIN penetration_test_classifications ON penetration_test_classifications.id = penetration_tests.penetration_test_classification_id").
		Joins("LEFT JOIN penetration_test_plans ON penetration_test_plans.id = penetration_tests.penetration_test_plan_id").
		Joins("LEFT JOIN penetration_test_executors ON penetration_test_executors.id = penetration_tests.penetration_test_executor_id").
		Joins("LEFT JOIN users AS cs_responsible_employee ON cs_responsible_employee.id = penetration_tests.cs_responsible_employee_id").
		Joins("LEFT JOIN penetration_test_statuses ON penetration_test_statuses.id = penetration_tests.penetration_test_status_id")
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where("reference_number LIKE ? OR system_application_name LIKE ? OR scope LIKE ?", keyword, keyword, keyword)
	}

	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where("penetration_tests."+columnName+" IN (?)", values)
		}
	}

	for _, filter := range params.QuickFilter {
		switch filter {
		case "Q1", "Q2", "Q3", "Q4":
			query = query.Where("penetration_tests.period = ?", filter)
		case "all":
			continue
		}
	}

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	if params.SortBy == "" {
		params.SortBy = "penetration_tests.created_at"
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "desc"
	}
	if params.Page < 1 {
		params.Page = 1
	}
	if params.PageSize < 1 {
		params.PageSize = 10
	}

	query = query.Order(params.SortBy + " " + params.Sort)
	if err := query.Offset((params.Page - 1) * params.PageSize).Limit(params.PageSize).Find(&penetrationTests).Error; err != nil {
		return nil, err
	}

	return &penetrationTestDTOs.PenetrationTestListResponseDTO{
		Items: penetrationTests,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize)),
			TotalCount:  totalCount,
		},
	}, nil
}

func (repo *PenetrationTestRepository) GetByID(id uint, tx ...utility.DatabaseTransaction) (*models.PenetrationTest, error) {
	var penetrationTest *models.PenetrationTest
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	result := db.
		Preload("PenetrationTestType").
		Preload("PenetrationTestPlan").
		Preload("PenetrationTestClassification").
		Preload("PenetrationTestExecutor").
		Preload("CSResponsibleEmployee").
		Preload("PenetrationTestStatus").
		Preload("EnvironmentPreparation.EnvironmentReadiness").
		Preload("Execution.ExecutionStatus").
		Preload("Execution.SecurityResult").
		Preload("ManagementSharing").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		First(&penetrationTest, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(result.Error, id)
		}
		return nil, result.Error
	}
	return penetrationTest, nil
}

func (repo *PenetrationTestRepository) Create(penetrationTest *models.PenetrationTest, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(penetrationTest).Error; err != nil {
		return 0, err
	}
	return penetrationTest.ID, nil
}

func (repo *PenetrationTestRepository) Update(penetrationTest *models.PenetrationTest, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Model(&models.PenetrationTest{}).Where("id = ?", penetrationTest.ID).Updates(penetrationTest).Error
}

func (repo *PenetrationTestRepository) UpdateStatus(id uint, statusID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Model(&models.PenetrationTest{}).Where("id = ?", id).Update("penetration_test_status_id", statusID).Error
}

func (repo *PenetrationTestRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	var penetrationTest *models.PenetrationTest
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.First(&penetrationTest, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utility.NewResourceNotFoundError(err, id)
		}
		return err
	}
	return db.Delete(&models.PenetrationTest{}, id).Error
}

func (repo *PenetrationTestRepository) DeletePenetrationTestDetails(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Where("penetration_test_id = ?", id).Delete(&models.PenetrationTestEnvironmentPreparation{}).Error; err != nil {
		return err
	}
	if err := db.Where("penetration_test_id = ?", id).Delete(&models.PenetrationTestExecution{}).Error; err != nil {
		return err
	}
	return db.Where("penetration_test_id = ?", id).Delete(&models.PenetrationTestManagementSharing{}).Error
}

func (repo *PenetrationTestRepository) GenerateReferenceNumber(tx ...utility.DatabaseTransaction) (string, error) {
	year := time.Now().Year()
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var lastReferenceNumber string
	if err := db.Model(&models.PenetrationTest{}).
		Select("reference_number").
		Where("reference_number LIKE ?", fmt.Sprintf("PT-%d-%%", year)).
		Order("reference_number desc").
		Limit(1).
		Pluck("reference_number", &lastReferenceNumber).Error; err != nil {
		return "", err
	}

	lastDigits := "000"
	if lastReferenceNumber != "" && len(lastReferenceNumber) >= 3 {
		lastDigits = lastReferenceNumber[len(lastReferenceNumber)-3:]
	}

	lastReferenceID, err := strconv.ParseUint(lastDigits, 10, 64)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("PT-%d-%03d", year, lastReferenceID+1), nil
}

func (repo *PenetrationTestRepository) GetSeverityCounts(penetrationTestID uint) (*PenetrationTestSeverityCounts, error) {
	counts := &PenetrationTestSeverityCounts{}
	if err := repo.db.Table("penetration_test_findings").
		Select(`
			COUNT(*) AS total_count,
			SUM(CASE WHEN severity_levels.type = 'critical' THEN 1 ELSE 0 END) AS critical_count,
			SUM(CASE WHEN severity_levels.type = 'high' THEN 1 ELSE 0 END) AS high_count,
			SUM(CASE WHEN severity_levels.type = 'medium' THEN 1 ELSE 0 END) AS medium_count,
			SUM(CASE WHEN severity_levels.type = 'low' THEN 1 ELSE 0 END) AS low_count
		`).
		Joins("INNER JOIN severity_levels ON severity_levels.id = penetration_test_findings.severity_level_id").
		Where("penetration_test_findings.penetration_test_id = ?", penetrationTestID).
		Scan(counts).Error; err != nil {
		return nil, err
	}
	return counts, nil
}

func (repo *PenetrationTestRepository) GetSeverityBreakdown(penetrationTestID uint) (*penetrationTestDTOs.PenetrationTestSeverityBreakdownDTO, error) {
	result := &penetrationTestDTOs.PenetrationTestSeverityBreakdownDTO{
		SeverityLevels: []penetrationTestDTOs.PenetrationTestSeverityCountDTO{},
	}

	if err := repo.db.Table("penetration_test_findings").
		Where("penetration_test_id = ?", penetrationTestID).
		Count(&result.TotalFindings).Error; err != nil {
		return nil, err
	}

	if err := repo.db.Table("severity_levels sl").
		Select(`
			sl.id,
			sl.title_en,
			sl.title_ar,
			sl.color_code,
			sl.type,
			COUNT(ptf.id) AS count
		`).
		Joins("LEFT JOIN penetration_test_findings ptf ON ptf.severity_level_id = sl.id AND ptf.penetration_test_id = ?", penetrationTestID).
		Group("sl.id, sl.title_en, sl.title_ar, sl.color_code, sl.type, sl.weight").
		Order("sl.weight DESC, sl.id").
		Scan(&result.SeverityLevels).Error; err != nil {
		return nil, err
	}

	for i := range result.SeverityLevels {
		if result.TotalFindings > 0 {
			result.SeverityLevels[i].Percentage = float64(result.SeverityLevels[i].Count) * 100 / float64(result.TotalFindings)
		}
	}

	return result, nil
}
