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

type PenetrationTestFindingRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(params *penetrationTestDTOs.PenetrationTestFindingListQueryParamsDTO) (*penetrationTestDTOs.PenetrationTestFindingListResponseDTO, error)
	GetByID(id uint, tx ...utility.DatabaseTransaction) (*models.PenetrationTestFinding, error)
	Create(penetrationTestFinding *models.PenetrationTestFinding, tx ...utility.DatabaseTransaction) (uint, error)
	Update(penetrationTestFinding *models.PenetrationTestFinding, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error
	GenerateReferenceNumber(tx ...utility.DatabaseTransaction) (string, error)
}

type PenetrationTestFindingRepository struct {
	db *gorm.DB
}

func NewPenetrationTestFindingRepository(db *gorm.DB) *PenetrationTestFindingRepository {
	return &PenetrationTestFindingRepository{
		db: db,
	}
}

func (repo *PenetrationTestFindingRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *PenetrationTestFindingRepository) GetAll(params *penetrationTestDTOs.PenetrationTestFindingListQueryParamsDTO) (*penetrationTestDTOs.PenetrationTestFindingListResponseDTO, error) {
	var penetrationTestFindings []*penetrationTestDTOs.PenetrationTestFindingListDTO
	var totalCount int64

	query := repo.db.Table("penetration_test_findings").
		Select(`
			penetration_test_findings.id,
			penetration_test_findings.reference_number,
			penetration_test_findings.penetration_test_id,
			penetration_test_findings.title,
			penetration_test_findings.action_plan_status_id,
			action_plan_statuses.title_en AS action_plan_status_title_en,
			action_plan_statuses.title_ar AS action_plan_status_title_ar,
			action_plan_statuses.color_code AS action_plan_status_color_code,
			action_plan_statuses.type AS action_plan_status_type,
			penetration_test_findings.severity_level_id,
			severity_levels.title_en AS severity_level_title_en,
			severity_levels.title_ar AS severity_level_title_ar,
			severity_levels.color_code AS severity_level_color_code,
			penetration_test_findings.owner_id,
			owner_user.first_name_en AS owner_first_name_en,
			owner_user.last_name_en AS owner_last_name_en,
			owner_user.first_name_ar AS owner_first_name_ar,
			owner_user.last_name_ar AS owner_last_name_ar,
			CONVERT(varchar(10), penetration_test_findings.target_remediation_date, 23) AS target_remediation_date,
			CONVERT(varchar(30), penetration_test_findings.created_at, 126) AS created_at,
			CONVERT(varchar(30), penetration_test_findings.updated_at, 126) AS updated_at
		`).
		Joins("LEFT JOIN severity_levels ON severity_levels.id = penetration_test_findings.severity_level_id").
		Joins("LEFT JOIN action_plan_statuses ON action_plan_statuses.id = penetration_test_findings.action_plan_status_id").
		Joins("LEFT JOIN users AS owner_user ON owner_user.id = penetration_test_findings.owner_id")

	if params.PenetrationTestID > 0 {
		query = query.Where("penetration_test_findings.penetration_test_id = ?", params.PenetrationTestID)
	}
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
			(penetration_test_findings.reference_number LIKE ? OR
			penetration_test_findings.title LIKE ? OR
			penetration_test_findings.description LIKE ? OR
			penetration_test_findings.recommendation LIKE ? OR
			severity_levels.title_en LIKE ? OR
			owner_user.first_name_en LIKE ? OR
			owner_user.last_name_en LIKE ?)`,
			keyword, keyword, keyword, keyword, keyword, keyword, keyword)
	}

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	if params.SortBy == "" {
		params.SortBy = "created_at"
	}
	allowedSortColumns := map[string]bool{
		"reference_number":        true,
		"title":                   true,
		"severity_level_id":       true,
		"action_plan_status_id":   true,
		"target_remediation_date": true,
		"created_at":              true,
	}
	if !allowedSortColumns[params.SortBy] {
		params.SortBy = "created_at"
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

	query = query.Order("penetration_test_findings." + params.SortBy + " " + params.Sort)
	if err := query.Offset((params.Page - 1) * params.PageSize).Limit(params.PageSize).Find(&penetrationTestFindings).Error; err != nil {
		return nil, err
	}

	return &penetrationTestDTOs.PenetrationTestFindingListResponseDTO{
		Items: penetrationTestFindings,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize)),
			TotalCount:  totalCount,
		},
	}, nil
}

func (repo *PenetrationTestFindingRepository) GetByID(id uint, tx ...utility.DatabaseTransaction) (*models.PenetrationTestFinding, error) {
	var penetrationTestFinding *models.PenetrationTestFinding
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	result := db.
		Preload("PenetrationTest").
		Preload("SeverityLevel").
		Preload("ActionPlanStatus").
		Preload("Owner").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		First(&penetrationTestFinding, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(result.Error, id)
		}
		return nil, result.Error
	}
	return penetrationTestFinding, nil
}

func (repo *PenetrationTestFindingRepository) Create(penetrationTestFinding *models.PenetrationTestFinding, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(penetrationTestFinding).Error; err != nil {
		return 0, err
	}
	return penetrationTestFinding.ID, nil
}

func (repo *PenetrationTestFindingRepository) Update(penetrationTestFinding *models.PenetrationTestFinding, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Model(&models.PenetrationTestFinding{}).Where("id = ?", penetrationTestFinding.ID).Updates(penetrationTestFinding).Error
}

func (repo *PenetrationTestFindingRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	var penetrationTestFinding *models.PenetrationTestFinding
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.First(&penetrationTestFinding, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utility.NewResourceNotFoundError(err, id)
		}
		return err
	}
	return db.Delete(&models.PenetrationTestFinding{}, id).Error
}

func (repo *PenetrationTestFindingRepository) GenerateReferenceNumber(tx ...utility.DatabaseTransaction) (string, error) {
	var count int64
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Model(&models.PenetrationTestFinding{}).Count(&count).Error; err != nil {
		return "", err
	}
	return fmt.Sprintf("FIND-%03d", count+1), nil
}
