package repositories

import (
	"errors"
	"fmt"
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/models"
	penetrationTestDTOs "jethur/grc/penetration-testing/dtos"
	"jethur/utility"
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
	GenerateReferenceNumber(tx ...utility.DatabaseTransaction) (string, error)
	ValidateReferences(penetrationTest *models.PenetrationTest) error
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

	query := repo.db.Table("penetration_tests")
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
		Preload("EnvironmentReadiness").
		Preload("ExecutionStatus").
		Preload("SecurityResult").
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
	return db.Delete(&models.PenetrationTest{}, id).Error
}

func (repo *PenetrationTestRepository) GenerateReferenceNumber(tx ...utility.DatabaseTransaction) (string, error) {
	year := time.Now().Year()
	var count int64
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Model(&models.PenetrationTest{}).
		Where("reference_number LIKE ?", fmt.Sprintf("PT-%d-%%", year)).
		Count(&count).Error; err != nil {
		return "", err
	}
	return fmt.Sprintf("PT-%d-%03d", year, count+1), nil
}

func (repo *PenetrationTestRepository) ValidateReferences(penetrationTest *models.PenetrationTest) error {
	references := []struct {
		table string
		id    uint
		field string
	}{
		{"penetration_test_types", penetrationTest.PenetrationTestTypeID, "penetration_test_type_id"},
		{"penetration_test_plans", penetrationTest.PenetrationTestPlanID, "penetration_test_plan_id"},
		{"penetration_test_classifications", penetrationTest.PenetrationTestClassificationID, "penetration_test_classification_id"},
		{"penetration_test_executors", penetrationTest.PenetrationTestExecutorID, "penetration_test_executor_id"},
		{"users", penetrationTest.CSResponsibleEmployeeID, "cs_responsible_employee_id"},
		{"penetration_test_environment_readinesses", penetrationTest.EnvironmentReadinessID, "environment_readiness_id"},
		{"penetration_test_execution_statuses", penetrationTest.ExecutionStatusID, "execution_status_id"},
		{"penetration_test_security_results", penetrationTest.SecurityResultID, "security_result_id"},
	}
	for _, reference := range references {
		var count int64
		if err := repo.db.Table(reference.table).Where("id = ?", reference.id).Count(&count).Error; err != nil {
			return err
		}
		if count == 0 {
			return utility.NewValidationError(fmt.Errorf("referenced record does not exist"), "invalid", reference.field)
		}
	}
	return nil
}
