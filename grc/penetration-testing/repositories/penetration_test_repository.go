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
	Delete(id uint, tx ...utility.DatabaseTransaction) error
	DeletePenetrationTestDetails(id uint, tx ...utility.DatabaseTransaction) error
	GenerateReferenceNumber(tx ...utility.DatabaseTransaction) (string, error)
	GetSeverityCounts(penetrationTestID uint) (*PenetrationTestSeverityCounts, error)
}

type PenetrationTestSeverityCounts struct {
	CriticalCount int
	HighCount     int
	MediumCount   int
	LowCount      int
	TotalFindings int
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
			penetration_tests.penetration_test_plan_id,
			penetration_tests.year,
			penetration_tests.period,
			penetration_tests.penetration_test_classification_id,
			penetration_tests.penetration_test_executor_id,
			penetration_tests.cs_responsible_employee_id,
			environment_preparations.environment_readiness_id,
			executions.test_start_date,
			executions.test_end_date,
			executions.execution_status_id,
			executions.security_result_id,
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
			) AS total_findings
		`).
		Joins("LEFT JOIN penetration_test_environment_preparations AS environment_preparations ON environment_preparations.penetration_test_id = penetration_tests.id").
		Joins("LEFT JOIN penetration_test_executions AS executions ON executions.penetration_test_id = penetration_tests.id").
		Joins("LEFT JOIN penetration_test_management_sharings AS management_sharings ON management_sharings.penetration_test_id = penetration_tests.id")
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where("reference_number LIKE ? OR system_application_name LIKE ? OR scope LIKE ?", keyword, keyword, keyword)
	}

	if params.SortBy == "" {
		params.SortBy = "created_at"
	}
	allowedSortColumns := map[string]bool{
		"reference_number":        true,
		"system_application_name": true,
		"year":                    true,
		"period":                  true,
		"created_at":              true,
	}
	if !allowedSortColumns[params.SortBy] {
		params.SortBy = "created_at"
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "desc"
	}

	query = query.Order("penetration_tests." + params.SortBy + " " + params.Sort)
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}
	if params.Page < 1 {
		params.Page = 1
	}
	if params.PageSize < 1 {
		params.PageSize = 10
	}
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
	if err := repo.DeletePenetrationTestDetails(id, tx...); err != nil {
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
			COUNT(*) AS total_findings,
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
