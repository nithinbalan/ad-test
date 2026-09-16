package repositories

import (
	"fmt"
	reportDTOs "jethur/grc/reports/dtos"

	"gorm.io/gorm"
)

type PenetrationTestReportRepositoryInterface interface {
	GetPenetrationTestByType(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.PenetrationTestByTypeDTO, error)
	GetPenetrationTestByStatus(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.PenetrationTestByStatusDTO, error)
	GetPenetrationTestByClassification(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.PenetrationTestByClassificationDTO, error)
	GetPenetrationTestByPlan(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.PenetrationTestByPlanDTO, error)
	GetPenetrationTestByExecutor(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.PenetrationTestByExecutorDTO, error)
	GetPenetrationTestByResponsibleUser(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.PenetrationTestByResponsibleUserDTO, error)
	GetPenetrationTestByCreatedUser(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.PenetrationTestByCreatedUserDTO, error)
	GetPenetrationTestByResponsibleUnit(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.PenetrationTestByResponsibleUnitDTO, error)
	GetPenetrationTestByYear(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.PenetrationTestByYearDTO, error)
	GetPenetrationTestByPeriod(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.PenetrationTestByPeriodDTO, error)
}

type PenetrationTestReportRepository struct {
	db *gorm.DB
}

func NewPenetrationTestReportRepository(db *gorm.DB) *PenetrationTestReportRepository {
	return &PenetrationTestReportRepository{db: db}
}

func (repo *PenetrationTestReportRepository) GetPenetrationTestByType(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.PenetrationTestByTypeDTO, error) {
	results := []*reportDTOs.PenetrationTestByTypeDTO{}

	query := repo.db.Table("penetration_test_types").
		Select(`
			penetration_test_types.id AS penetration_test_type_id,
			penetration_test_types.title_en AS penetration_test_type_title_en,
			penetration_test_types.title_ar AS penetration_test_type_title_ar,
			COUNT(penetration_tests.id) AS count
		`).
		Joins("LEFT JOIN penetration_tests ON penetration_tests.penetration_test_type_id = penetration_test_types.id").
		Group("penetration_test_types.id, penetration_test_types.title_en, penetration_test_types.title_ar").
		Order("COUNT(penetration_tests.id) DESC")

	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("penetration_tests.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.Scan(&results).Error
	if err != nil {
		return nil, err
	}

	return results, nil
}

func (repo *PenetrationTestReportRepository) GetPenetrationTestByStatus(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.PenetrationTestByStatusDTO, error) {
	results := []*reportDTOs.PenetrationTestByStatusDTO{}

	query := repo.db.Table("penetration_test_statuses").
		Select(`
			penetration_test_statuses.id,
			penetration_test_statuses.title_en AS title_en,
			penetration_test_statuses.title_ar AS title_ar,
			penetration_test_statuses.color_code,
			COUNT(penetration_tests.id) AS count
		`).
		Joins("LEFT JOIN penetration_tests ON penetration_tests.penetration_test_status_id = penetration_test_statuses.id").
		Group("penetration_test_statuses.id, penetration_test_statuses.title_en, penetration_test_statuses.title_ar, penetration_test_statuses.color_code")

	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("penetration_tests.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.Scan(&results).Error
	return results, err
}

func (repo *PenetrationTestReportRepository) GetPenetrationTestByClassification(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.PenetrationTestByClassificationDTO, error) {
	results := []*reportDTOs.PenetrationTestByClassificationDTO{}

	query := repo.db.Table("penetration_test_classifications").
		Select(`
			penetration_test_classifications.id,
			penetration_test_classifications.title_en AS title_en,
			penetration_test_classifications.title_ar AS title_ar,
			COUNT(penetration_tests.id) AS count
		`).
		Joins("LEFT JOIN penetration_tests ON penetration_tests.penetration_test_classification_id = penetration_test_classifications.id").
		Group("penetration_test_classifications.id, penetration_test_classifications.title_en, penetration_test_classifications.title_ar")

	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("penetration_tests.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.Scan(&results).Error
	return results, err
}

func (repo *PenetrationTestReportRepository) GetPenetrationTestByPlan(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.PenetrationTestByPlanDTO, error) {
	results := []*reportDTOs.PenetrationTestByPlanDTO{}

	query := repo.db.Table("penetration_test_plans").
		Select(`
			penetration_test_plans.id,
			penetration_test_plans.title_en AS title_en,
			penetration_test_plans.title_ar AS title_ar,
			COUNT(penetration_tests.id) AS count
		`).
		Joins("LEFT JOIN penetration_tests ON penetration_tests.penetration_test_plan_id = penetration_test_plans.id").
		Group("penetration_test_plans.id, penetration_test_plans.title_en, penetration_test_plans.title_ar")

	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("penetration_tests.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.Scan(&results).Error
	return results, err
}

func (repo *PenetrationTestReportRepository) GetPenetrationTestByExecutor(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.PenetrationTestByExecutorDTO, error) {
	results := []*reportDTOs.PenetrationTestByExecutorDTO{}

	query := repo.db.Table("penetration_test_executors").
		Select(`
			penetration_test_executors.id,
			penetration_test_executors.title_en AS title_en,
			penetration_test_executors.title_ar AS title_ar,
			COUNT(penetration_tests.id) AS count
		`).
		Joins("LEFT JOIN penetration_tests ON penetration_tests.penetration_test_executor_id = penetration_test_executors.id").
		Group("penetration_test_executors.id, penetration_test_executors.title_en, penetration_test_executors.title_ar")

	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("penetration_tests.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.Scan(&results).Error
	return results, err
}

func (repo *PenetrationTestReportRepository) GetPenetrationTestByResponsibleUser(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.PenetrationTestByResponsibleUserDTO, error) {
	results := []*reportDTOs.PenetrationTestByResponsibleUserDTO{}

	query := repo.db.Table("penetration_tests").
		Select(`
			users.id AS responsible_user_id,
			users.first_name_en AS responsible_user_first_name_en,
			users.last_name_en AS responsible_user_last_name_en,
			users.first_name_ar AS responsible_user_first_name_ar,
			users.last_name_ar AS responsible_user_last_name_ar,
			users.initial AS responsible_user_initial,
			users.color AS responsible_user_color,
			users.profile_image_path AS responsible_user_profile_image_path,
			CONCAT(users.first_name_en, ' ', users.last_name_en) AS responsible_user_title_en,
			CONCAT(users.first_name_ar, ' ', users.last_name_ar) AS responsible_user_title_ar,
			users.email AS responsible_user_email,
			COUNT(penetration_tests.id) AS count
		`).
		Joins("JOIN users ON users.id = penetration_tests.cs_responsible_employee_id").
		Group("users.id, users.first_name_en, users.last_name_en, users.first_name_ar, users.last_name_ar, users.initial, users.color, users.profile_image_path, users.email")

	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("penetration_tests.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.Scan(&results).Error
	return results, err
}

func (repo *PenetrationTestReportRepository) GetPenetrationTestByCreatedUser(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.PenetrationTestByCreatedUserDTO, error) {
	results := []*reportDTOs.PenetrationTestByCreatedUserDTO{}

	query := repo.db.Table("penetration_tests").
		Select(`
			users.id AS created_user_id,
			users.first_name_en AS created_user_first_name_en,
			users.last_name_en AS created_user_last_name_en,
			users.first_name_ar AS created_user_first_name_ar,
			users.last_name_ar AS created_user_last_name_ar,
			users.initial AS created_user_initial,
			users.color AS created_user_color,
			users.profile_image_path AS created_user_profile_image_path,
			CONCAT(users.first_name_en, ' ', users.last_name_en) AS created_user_title_en,
			CONCAT(users.first_name_ar, ' ', users.last_name_ar) AS created_user_title_ar,
			users.email AS created_user_email,
			COUNT(penetration_tests.id) AS count
		`).
		Joins("JOIN users ON users.id = penetration_tests.created_by_id").
		Group("users.id, users.first_name_en, users.last_name_en, users.first_name_ar, users.last_name_ar, users.initial, users.color, users.profile_image_path, users.email")

	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("penetration_tests.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.Scan(&results).Error
	return results, err
}

func (repo *PenetrationTestReportRepository) GetPenetrationTestByResponsibleUnit(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.PenetrationTestByResponsibleUnitDTO, error) {
	results := []*reportDTOs.PenetrationTestByResponsibleUnitDTO{}

	query := repo.db.Table("penetration_tests").
		Select(`
			organization_units.id AS responsible_unit_id,
			organization_units.title_en AS responsible_unit_title_en,
			organization_units.title_ar AS responsible_unit_title_ar,
			organization_units.color_code AS responsible_unit_color_code,
			COUNT(penetration_tests.id) AS count
		`).
		Joins("JOIN users ON users.id = penetration_tests.cs_responsible_employee_id").
		Joins("JOIN organization_units ON organization_units.id = users.organization_unit_1_id").
		Group("organization_units.id, organization_units.title_en, organization_units.title_ar, organization_units.color_code")

	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("penetration_tests.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.Scan(&results).Error
	return results, err
}

func (repo *PenetrationTestReportRepository) GetPenetrationTestByYear(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.PenetrationTestByYearDTO, error) {
	results := []*reportDTOs.PenetrationTestByYearDTO{}
	var rawResults []struct {
		Year  int  `gorm:"column:year"`
		Count uint `gorm:"column:count"`
	}

	query := repo.db.Table("penetration_tests").
		Select(`
			year,
			COUNT(id) AS count
		`).
		Group("year").
		Order("year ASC")

	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("penetration_tests.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.Scan(&rawResults).Error
	if err != nil {
		return nil, err
	}

	for _, record := range rawResults {
		results = append(results, &reportDTOs.PenetrationTestByYearDTO{
			TitleEn: fmt.Sprintf("%d", record.Year),
			TitleAr: fmt.Sprintf("%d", record.Year),
			Count:   record.Count,
		})
	}

	return results, nil
}

func (repo *PenetrationTestReportRepository) GetPenetrationTestByPeriod(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.PenetrationTestByPeriodDTO, error) {
	results := []*reportDTOs.PenetrationTestByPeriodDTO{}
	var rawResults []struct {
		Period string `gorm:"column:period"`
		Count  uint   `gorm:"column:count"`
	}

	query := repo.db.Table("penetration_tests").
		Select(`
			period,
			COUNT(id) AS count
		`).
		Group("period").
		Order("period ASC")

	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("penetration_tests.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.Scan(&rawResults).Error
	if err != nil {
		return nil, err
	}

	for _, record := range rawResults {
		results = append(results, &reportDTOs.PenetrationTestByPeriodDTO{
			TitleEn: fmt.Sprintf("%s", record.Period),
			TitleAr: fmt.Sprintf("%s", record.Period),
			Count:   record.Count,
		})
	}

	return results, nil
}
