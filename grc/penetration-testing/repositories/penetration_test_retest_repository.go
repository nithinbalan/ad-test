package repositories

import (
	"errors"
	"fmt"
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/models"
	penetrationTestDTOs "jethur/grc/penetration-testing/dtos"
	"jethur/utility"
	"strconv"

	"gorm.io/gorm"
)

type PenetrationTestRetestRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(params *penetrationTestDTOs.PenetrationTestRetestListQueryParamsDTO) (*penetrationTestDTOs.PenetrationTestRetestListResponseDTO, error)
	GetActiveByPenetrationTestID(penetrationTestID uint) (*models.PenetrationTestRetest, error)
	GetByID(id uint, tx ...utility.DatabaseTransaction) (*models.PenetrationTestRetest, error)
	Create(penetrationTestRetest *models.PenetrationTestRetest, tx ...utility.DatabaseTransaction) (uint, error)
	Update(penetrationTestRetest *models.PenetrationTestRetest, tx ...utility.DatabaseTransaction) error
	UpdateStatus(id uint, statusID uint, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error
	CreateRetestFinding(penetrationTestRetestFinding *models.PenetrationTestRetestFinding, tx ...utility.DatabaseTransaction) (uint, error)
	UpdateRetestFinding(penetrationTestRetestFinding *models.PenetrationTestRetestFinding, tx ...utility.DatabaseTransaction) error
	DeleteRetestFindings(penetrationTestRetestID uint, tx ...utility.DatabaseTransaction) error
	GetRetestFindingByID(penetrationTestRetestID uint, penetrationTestRetestFindingID uint, tx ...utility.DatabaseTransaction) (*models.PenetrationTestRetestFinding, error)
	DeleteRetestFindingByID(penetrationTestRetestID uint, penetrationTestRetestFindingID uint, tx ...utility.DatabaseTransaction) error
	CreateDocuments(docs []*models.PenetrationTestRetestFindingDocument, tx ...utility.DatabaseTransaction) error
	GetDocumentByID(id uint) (*models.PenetrationTestRetestFindingDocument, error)
	DeleteDocuments(docIDs []*uint, tx ...utility.DatabaseTransaction) error
	DeleteRetestFindingDocumentsByRetestID(penetrationTestRetestID uint, tx ...utility.DatabaseTransaction) error
	GetEligibleOpenFindings(penetrationTestRetestID uint, penetrationTestID uint, closedStatusID uint) ([]*penetrationTestDTOs.PenetrationTestFindingListDTO, error)
	GenerateReferenceNumber(penetrationTestID uint, tx ...utility.DatabaseTransaction) (string, error)
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
			penetration_test_retests.penetration_test_retest_status_id,
			penetration_test_retest_statuses.title_en AS status_title_en,
			penetration_test_retest_statuses.title_ar AS status_title_ar,
			penetration_test_retest_statuses.color_code AS status_color_code,
			penetration_test_retest_statuses.type AS status_type,
			penetration_test_retests.notes,
			(
				SELECT COUNT(*)
				FROM penetration_test_retest_findings
				INNER JOIN penetration_test_findings ON penetration_test_findings.id = penetration_test_retest_findings.penetration_test_finding_id
				INNER JOIN severity_levels ON severity_levels.id = penetration_test_findings.severity_level_id
				LEFT JOIN penetration_test_retest_results ON penetration_test_retest_results.id = penetration_test_retest_findings.retest_result_id
				WHERE penetration_test_retest_findings.penetration_test_retest_id = penetration_test_retests.id
				AND (penetration_test_retest_results.type IS NULL OR penetration_test_retest_results.type <> 'pass')
				AND severity_levels.type = 'critical'
			) AS remaining_critical,
			(
				SELECT COUNT(*)
				FROM penetration_test_retest_findings
				INNER JOIN penetration_test_findings ON penetration_test_findings.id = penetration_test_retest_findings.penetration_test_finding_id
				INNER JOIN severity_levels ON severity_levels.id = penetration_test_findings.severity_level_id
				LEFT JOIN penetration_test_retest_results ON penetration_test_retest_results.id = penetration_test_retest_findings.retest_result_id
				WHERE penetration_test_retest_findings.penetration_test_retest_id = penetration_test_retests.id
				AND (penetration_test_retest_results.type IS NULL OR penetration_test_retest_results.type <> 'pass')
				AND severity_levels.type = 'high'
			) AS remaining_high,
			(
				SELECT COUNT(*)
				FROM penetration_test_retest_findings
				INNER JOIN penetration_test_findings ON penetration_test_findings.id = penetration_test_retest_findings.penetration_test_finding_id
				INNER JOIN severity_levels ON severity_levels.id = penetration_test_findings.severity_level_id
				LEFT JOIN penetration_test_retest_results ON penetration_test_retest_results.id = penetration_test_retest_findings.retest_result_id
				WHERE penetration_test_retest_findings.penetration_test_retest_id = penetration_test_retests.id
				AND (penetration_test_retest_results.type IS NULL OR penetration_test_retest_results.type <> 'pass')
				AND severity_levels.type = 'medium'
			) AS remaining_medium,
			(
				SELECT COUNT(*)
				FROM penetration_test_retest_findings
				INNER JOIN penetration_test_findings ON penetration_test_findings.id = penetration_test_retest_findings.penetration_test_finding_id
				INNER JOIN severity_levels ON severity_levels.id = penetration_test_findings.severity_level_id
				LEFT JOIN penetration_test_retest_results ON penetration_test_retest_results.id = penetration_test_retest_findings.retest_result_id
				WHERE penetration_test_retest_findings.penetration_test_retest_id = penetration_test_retests.id
				AND (penetration_test_retest_results.type IS NULL OR penetration_test_retest_results.type <> 'pass')
				AND severity_levels.type = 'low'
			) AS remaining_low,
			(
				SELECT COUNT(*)
				FROM penetration_test_retest_findings
				LEFT JOIN penetration_test_retest_results ON penetration_test_retest_results.id = penetration_test_retest_findings.retest_result_id
				WHERE penetration_test_retest_findings.penetration_test_retest_id = penetration_test_retests.id
				AND (penetration_test_retest_results.type IS NULL OR penetration_test_retest_results.type <> 'pass')
			) AS remaining_total,
			CONVERT(varchar(30), penetration_test_retests.created_at, 126) AS created_at,
			CONVERT(varchar(30), penetration_test_retests.updated_at, 126) AS updated_at,
			penetration_test_retests.created_by_id,
			created_user.first_name_en AS created_user_first_name_en,
			created_user.last_name_en AS created_user_last_name_en,
			created_user.first_name_ar AS created_user_first_name_ar,
			created_user.last_name_ar AS created_user_last_name_ar,
			created_user.initial AS created_user_initial,
			created_user.color AS created_user_color,
			created_user.email AS created_user_email,
			created_user.profile_image_path AS created_user_profile_image_path
		`).
		Joins("LEFT JOIN penetration_test_executors ON penetration_test_executors.id = penetration_test_retests.penetration_test_executor_id").
		Joins("LEFT JOIN penetration_test_retest_statuses ON penetration_test_retest_statuses.id = penetration_test_retests.penetration_test_retest_status_id").
		Joins("LEFT JOIN users AS created_user ON created_user.id = penetration_test_retests.created_by_id")

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
		query = query.Where("penetration_test_retest_statuses.type = ?", params.Status)
	}

	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where(columnName+" IN (?)", values)
		}
	}
	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("penetration_test_retests.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	for _, filter := range params.QuickFilter {
		switch filter {
		case "in_progress":
			query = query.Where("penetration_test_retest_statuses.type = ?", "in-progress")
		case "completed":
			query = query.Where("penetration_test_retest_statuses.type = ?", "completed")
		}
	}

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	if params.SortBy == "" {
		params.SortBy = "penetration_test_retests.start_date"
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
		Preload("PenetrationTestRetestStatus").
		Preload("RetestFindings").
		Preload("RetestFindings.Documents").
		Preload("RetestFindings.PenetrationTestFinding").
		Preload("RetestFindings.PenetrationTestFinding.SeverityLevel").
		Preload("RetestFindings.CommentByUser").
		Preload("RetestFindings.ResultByUser").
		Preload("RetestFindings.RetestResult").
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
	result := repo.db.
		Joins("LEFT JOIN penetration_test_retest_statuses ON penetration_test_retest_statuses.id = penetration_test_retests.penetration_test_retest_status_id").
		Where("penetration_test_retests.penetration_test_id = ? AND penetration_test_retest_statuses.type IN ?", penetrationTestID, []string{"draft", "in-progress"}).
		First(&penetrationTestRetest)
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

func (repo *PenetrationTestRetestRepository) UpdateStatus(id uint, statusID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Model(&models.PenetrationTestRetest{}).Where("id = ?", id).Update("penetration_test_retest_status_id", statusID).Error
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

func (repo *PenetrationTestRetestRepository) GetRetestFindingByID(penetrationTestRetestID uint, penetrationTestRetestFindingID uint, tx ...utility.DatabaseTransaction) (*models.PenetrationTestRetestFinding, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var penetrationTestRetestFinding *models.PenetrationTestRetestFinding
	if err := db.Preload("PenetrationTestFinding").Preload("Documents").Where("penetration_test_retest_id = ? AND id = ?", penetrationTestRetestID, penetrationTestRetestFindingID).First(&penetrationTestRetestFinding).Error; err != nil {
		return nil, err
	}
	return penetrationTestRetestFinding, nil
}

func (repo *PenetrationTestRetestRepository) DeleteRetestFindingByID(penetrationTestRetestID uint, penetrationTestRetestFindingID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Where("penetration_test_retest_id = ? AND id = ?", penetrationTestRetestID, penetrationTestRetestFindingID).Delete(&models.PenetrationTestRetestFinding{}).Error
}

func (repo *PenetrationTestRetestRepository) CreateDocuments(docs []*models.PenetrationTestRetestFindingDocument, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Create(&docs).Error
}

func (repo *PenetrationTestRetestRepository) GetDocumentByID(id uint) (*models.PenetrationTestRetestFindingDocument, error) {
	var document models.PenetrationTestRetestFindingDocument
	if err := repo.db.Where("id = ?", id).First(&document).Error; err != nil {
		return nil, err
	}
	return &document, nil
}

func (repo *PenetrationTestRetestRepository) DeleteDocuments(docIDs []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Where("id IN (?)", docIDs).Delete(&models.PenetrationTestRetestFindingDocument{}).Error
}

func (repo *PenetrationTestRetestRepository) DeleteRetestFindingDocumentsByRetestID(penetrationTestRetestID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Where("penetration_test_retest_finding_id IN (?)", db.Model(&models.PenetrationTestRetestFinding{}).Select("id").Where("penetration_test_retest_id = ?", penetrationTestRetestID)).Delete(&models.PenetrationTestRetestFindingDocument{}).Error
}

func (repo *PenetrationTestRetestRepository) GetEligibleOpenFindings(penetrationTestRetestID uint, penetrationTestID uint, closedStatusID uint) ([]*penetrationTestDTOs.PenetrationTestFindingListDTO, error) {
	var penetrationTestFindings []*penetrationTestDTOs.PenetrationTestFindingListDTO

	query := repo.db.Table("penetration_test_findings").
		Select(`
			penetration_test_findings.id,
			penetration_test_findings.reference_number,
			penetration_test_findings.penetration_test_id,
			penetration_test_findings.title,
			penetration_test_findings.action_plan_status_id,
			penetration_test_findings.percentage,
			action_plan_statuses.title_en AS action_plan_status_title_en,
			action_plan_statuses.title_ar AS action_plan_status_title_ar,
			action_plan_statuses.color_code AS action_plan_status_color_code,
			action_plan_statuses.type AS action_plan_status_type,
			penetration_test_findings.description,
			penetration_test_findings.recommendation,
			penetration_test_findings.severity_level_id,
			severity_levels.title_en AS severity_level_title_en,
			severity_levels.title_ar AS severity_level_title_ar,
			severity_levels.color_code AS severity_level_color_code,
			penetration_test_findings.owner_id,
			owner_user.first_name_en AS owner_first_name_en,
			owner_user.last_name_en AS owner_last_name_en,
			owner_user.first_name_ar AS owner_first_name_ar,
			owner_user.last_name_ar AS owner_last_name_ar,
			owner_user.initial AS owner_initial,
			owner_user.color AS owner_color,
			owner_user.email AS owner_email,
			owner_user.profile_image_path AS owner_profile_image_path,
			CONVERT(varchar(10), penetration_test_findings.target_remediation_date, 23) AS target_remediation_date,
			CONVERT(varchar(30), penetration_test_findings.created_at, 126) AS created_at,
			CONVERT(varchar(30), penetration_test_findings.updated_at, 126) AS updated_at
		`).
		Joins("LEFT JOIN severity_levels ON severity_levels.id = penetration_test_findings.severity_level_id").
		Joins("LEFT JOIN action_plan_statuses ON action_plan_statuses.id = penetration_test_findings.action_plan_status_id").
		Joins("LEFT JOIN users AS owner_user ON owner_user.id = penetration_test_findings.owner_id").
		Where("penetration_test_findings.penetration_test_id = ?", penetrationTestID).
		Where("penetration_test_findings.action_plan_status_id <> ?", closedStatusID).
		Where("penetration_test_findings.id NOT IN (?)", repo.db.Model(&models.PenetrationTestRetestFinding{}).Select("penetration_test_finding_id").Where("penetration_test_retest_id = ?", penetrationTestRetestID))

	if err := query.Order("penetration_test_findings.created_at desc").Find(&penetrationTestFindings).Error; err != nil {
		return nil, err
	}

	return penetrationTestFindings, nil
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

	var lastReferenceNumber string
	if err := db.Model(&models.PenetrationTestRetest{}).
		Select("reference_number").
		Where("penetration_test_id = ?", penetrationTestID).
		Order("reference_number desc").
		Limit(1).
		Pluck("reference_number", &lastReferenceNumber).Error; err != nil {
		return "", err
	}

	lastDigits := "00"
	if lastReferenceNumber != "" && len(lastReferenceNumber) >= 2 {
		lastDigits = lastReferenceNumber[len(lastReferenceNumber)-2:]
	}

	lastReferenceID, err := strconv.ParseUint(lastDigits, 10, 64)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("RT-%s-%02d", penetrationTest.ReferenceNumber[3:], lastReferenceID+1), nil
}
