package repositories

import (
	"errors"
	"fmt"
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/models"
	penetrationTestDTOs "jethur/grc/penetration-testing/dtos"
	"jethur/utility"

	"gorm.io/gorm"
)

type PenetrationTestRetestRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(params *penetrationTestDTOs.PenetrationTestRetestListQueryParamsDTO) (*penetrationTestDTOs.PenetrationTestRetestListResponseDTO, error)
	GetActiveByPenetrationTestID(penetrationTestID uint) (*models.PenetrationTestRetest, error)
	GetByID(id uint, tx ...utility.DatabaseTransaction) (*models.PenetrationTestRetest, error)
	Create(penetrationTestRetest *models.PenetrationTestRetest, tx ...utility.DatabaseTransaction) (uint, error)
	Update(penetrationTestRetest *models.PenetrationTestRetest, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error
	CreateRetestFinding(penetrationTestRetestFinding *models.PenetrationTestRetestFinding, tx ...utility.DatabaseTransaction) (uint, error)
	UpdateRetestFinding(penetrationTestRetestFinding *models.PenetrationTestRetestFinding, tx ...utility.DatabaseTransaction) error
	DeleteRetestFindings(penetrationTestRetestID uint, tx ...utility.DatabaseTransaction) error
	GenerateReferenceNumber(penetrationTestID uint, tx ...utility.DatabaseTransaction) (string, error)
	ValidateReferences(penetrationTestRetest *models.PenetrationTestRetest, findingIDs []uint) error
}

type PenetrationTestRetestRepository struct {
	db *gorm.DB
}

func NewPenetrationTestRetestRepository(db *gorm.DB) *PenetrationTestRetestRepository {
	return &PenetrationTestRetestRepository{
		db: db,
	}
}

func (repo *PenetrationTestRetestRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *PenetrationTestRetestRepository) GetAll(params *penetrationTestDTOs.PenetrationTestRetestListQueryParamsDTO) (*penetrationTestDTOs.PenetrationTestRetestListResponseDTO, error) {
	var penetrationTestRetests []*penetrationTestDTOs.PenetrationTestRetestListDTO
	var totalCount int64

	query := repo.db.Table("penetration_test_retests AS penetration_test_retests").
		Select(`
			penetration_test_retests.id,
			penetration_test_retests.reference_number,
			penetration_test_retests.penetration_test_id,
			CONVERT(varchar(10), penetration_test_retests.start_date, 23) AS start_date,
			CONVERT(varchar(10), penetration_test_retests.end_date, 23) AS end_date,
			penetration_test_retests.penetration_test_executor_id,
			penetration_test_executors.title_en AS executor_title_en,
			penetration_test_executors.title_ar AS executor_title_ar,
			penetration_test_retests.status,
			penetration_test_retests.notes,
			(
				SELECT COUNT(*)
				FROM penetration_test_retest_findings
				INNER JOIN penetration_test_findings ON penetration_test_findings.id = penetration_test_retest_findings.penetration_test_finding_id
				INNER JOIN severity_levels ON severity_levels.id = penetration_test_findings.severity_level_id
				WHERE penetration_test_retest_findings.penetration_test_retest_id = penetration_test_retests.id
				AND penetration_test_retest_findings.is_verified = 0
				AND severity_levels.type = 'critical'
			) AS remaining_critical,
			(
				SELECT COUNT(*)
				FROM penetration_test_retest_findings
				INNER JOIN penetration_test_findings ON penetration_test_findings.id = penetration_test_retest_findings.penetration_test_finding_id
				INNER JOIN severity_levels ON severity_levels.id = penetration_test_findings.severity_level_id
				WHERE penetration_test_retest_findings.penetration_test_retest_id = penetration_test_retests.id
				AND penetration_test_retest_findings.is_verified = 0
				AND severity_levels.type = 'high'
			) AS remaining_high,
			(
				SELECT COUNT(*)
				FROM penetration_test_retest_findings
				INNER JOIN penetration_test_findings ON penetration_test_findings.id = penetration_test_retest_findings.penetration_test_finding_id
				INNER JOIN severity_levels ON severity_levels.id = penetration_test_findings.severity_level_id
				WHERE penetration_test_retest_findings.penetration_test_retest_id = penetration_test_retests.id
				AND penetration_test_retest_findings.is_verified = 0
				AND severity_levels.type = 'medium'
			) AS remaining_medium,
			(
				SELECT COUNT(*)
				FROM penetration_test_retest_findings
				INNER JOIN penetration_test_findings ON penetration_test_findings.id = penetration_test_retest_findings.penetration_test_finding_id
				INNER JOIN severity_levels ON severity_levels.id = penetration_test_findings.severity_level_id
				WHERE penetration_test_retest_findings.penetration_test_retest_id = penetration_test_retests.id
				AND penetration_test_retest_findings.is_verified = 0
				AND severity_levels.type = 'low'
			) AS remaining_low,
			(
				SELECT COUNT(*)
				FROM penetration_test_retest_findings
				WHERE penetration_test_retest_findings.penetration_test_retest_id = penetration_test_retests.id
				AND penetration_test_retest_findings.is_verified = 0
			) AS remaining_total
		`).
		Joins("LEFT JOIN penetration_test_executors ON penetration_test_executors.id = penetration_test_retests.penetration_test_executor_id")

	if params.PenetrationTestID > 0 {
		query = query.Where("penetration_test_retests.penetration_test_id = ?", params.PenetrationTestID)
	}
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
			(penetration_test_retests.reference_number LIKE ? OR
			penetration_test_retests.notes LIKE ? OR
			penetration_test_executors.title_en LIKE ? OR
			penetration_test_executors.title_ar LIKE ?)`,
			keyword, keyword, keyword, keyword)
	}
	if params.Status != "" {
		query = query.Where("penetration_test_retests.status = ?", params.Status)
	}

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	if params.SortBy == "" {
		params.SortBy = "start_date"
	}
	allowedSortColumns := map[string]bool{
		"reference_number": true,
		"start_date":       true,
		"end_date":         true,
		"status":           true,
	}
	if !allowedSortColumns[params.SortBy] {
		params.SortBy = "start_date"
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

	query = query.Order("penetration_test_retests." + params.SortBy + " " + params.Sort)
	if err := query.Offset((params.Page - 1) * params.PageSize).Limit(params.PageSize).Find(&penetrationTestRetests).Error; err != nil {
		return nil, err
	}

	return &penetrationTestDTOs.PenetrationTestRetestListResponseDTO{
		Items: penetrationTestRetests,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize)),
			TotalCount:  totalCount,
		},
	}, nil
}

func (repo *PenetrationTestRetestRepository) GetByID(id uint, tx ...utility.DatabaseTransaction) (*models.PenetrationTestRetest, error) {
	var penetrationTestRetest *models.PenetrationTestRetest
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	result := db.
		Preload("PenetrationTest").
		Preload("PenetrationTestExecutor").
		Preload("RetestFindings").
		Preload("RetestFindings.PenetrationTestFinding").
		Preload("RetestFindings.PenetrationTestFinding.SeverityLevel").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		First(&penetrationTestRetest, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(result.Error, id)
		}
		return nil, result.Error
	}
	return penetrationTestRetest, nil
}

func (repo *PenetrationTestRetestRepository) GetActiveByPenetrationTestID(penetrationTestID uint) (*models.PenetrationTestRetest, error) {
	var penetrationTestRetest *models.PenetrationTestRetest
	result := repo.db.Where("penetration_test_id = ? AND status = ?", penetrationTestID, "in-progress").First(&penetrationTestRetest)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return penetrationTestRetest, nil
}

func (repo *PenetrationTestRetestRepository) Create(penetrationTestRetest *models.PenetrationTestRetest, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(penetrationTestRetest).Error; err != nil {
		return 0, err
	}
	return penetrationTestRetest.ID, nil
}

func (repo *PenetrationTestRetestRepository) Update(penetrationTestRetest *models.PenetrationTestRetest, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Model(&models.PenetrationTestRetest{}).Where("id = ?", penetrationTestRetest.ID).Updates(penetrationTestRetest).Error
}

func (repo *PenetrationTestRetestRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	var penetrationTestRetest *models.PenetrationTestRetest
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.First(&penetrationTestRetest, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utility.NewResourceNotFoundError(err, id)
		}
		return err
	}
	return db.Delete(&models.PenetrationTestRetest{}, id).Error
}

func (repo *PenetrationTestRetestRepository) CreateRetestFinding(penetrationTestRetestFinding *models.PenetrationTestRetestFinding, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(penetrationTestRetestFinding).Error; err != nil {
		return 0, err
	}
	return penetrationTestRetestFinding.ID, nil
}

func (repo *PenetrationTestRetestRepository) UpdateRetestFinding(penetrationTestRetestFinding *models.PenetrationTestRetestFinding, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Model(&models.PenetrationTestRetestFinding{}).Where("id = ?", penetrationTestRetestFinding.ID).Updates(penetrationTestRetestFinding).Error
}

func (repo *PenetrationTestRetestRepository) DeleteRetestFindings(penetrationTestRetestID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Where("penetration_test_retest_id = ?", penetrationTestRetestID).Delete(&models.PenetrationTestRetestFinding{}).Error
}

func (repo *PenetrationTestRetestRepository) GenerateReferenceNumber(penetrationTestID uint, tx ...utility.DatabaseTransaction) (string, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var penetrationTest models.PenetrationTest
	if err := db.Select("reference_number").First(&penetrationTest, penetrationTestID).Error; err != nil {
		return "", err
	}

	var count int64
	if err := db.Model(&models.PenetrationTestRetest{}).Where("penetration_test_id = ?", penetrationTestID).Count(&count).Error; err != nil {
		return "", err
	}
	return fmt.Sprintf("RT-%s-%02d", penetrationTest.ReferenceNumber[3:], count+1), nil
}

func (repo *PenetrationTestRetestRepository) ValidateReferences(penetrationTestRetest *models.PenetrationTestRetest, findingIDs []uint) error {
	references := []struct {
		table string
		id    uint
		field string
	}{
		{"penetration_tests", penetrationTestRetest.PenetrationTestID, "penetration_test_id"},
		{"penetration_test_executors", penetrationTestRetest.PenetrationTestExecutorID, "penetration_test_executor_id"},
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

	for _, findingID := range findingIDs {
		var count int64
		if err := repo.db.Table("penetration_test_findings").Where("id = ? AND penetration_test_id = ?", findingID, penetrationTestRetest.PenetrationTestID).Count(&count).Error; err != nil {
			return err
		}
		if count == 0 {
			return utility.NewValidationError(fmt.Errorf("referenced record does not exist"), "invalid", "finding_ids")
		}
	}
	return nil
}
