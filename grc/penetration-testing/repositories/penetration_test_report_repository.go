package repositories

import (
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/models"
	penetrationTestDTOs "jethur/grc/penetration-testing/dtos"
	"jethur/utility"
	"time"

	"gorm.io/gorm"
)

type PenetrationTestReportRepositoryInterface interface {
	GetFullReportDetailsByID(id uint) (*models.PenetrationTest, error)
	BeginTransaction() utility.DatabaseTransaction
	GetAll(params *penetrationTestDTOs.PenetrationTestReportListQueryParamsDTO) (*penetrationTestDTOs.PenetrationTestReportListResponseDTO, error)
	Create(report *models.PenetrationTestReport, tx utility.DatabaseTransaction) (uint, error)
	CreateSection(section *models.PenetrationTestReportSection, tx utility.DatabaseTransaction) (uint, error)
	GetPenetrationTestReportByID(penetrationTestReportID uint, tx ...utility.DatabaseTransaction) (*models.PenetrationTestReport, error)
	GetSectionByID(sectionID uint, tx ...utility.DatabaseTransaction) (*models.PenetrationTestReportSection, error)
	UpdateSection(section *models.PenetrationTestReportSection, tx ...utility.DatabaseTransaction) error
	DeletePenetrationTestReportSection(sectionID uint, tx ...utility.DatabaseTransaction) error
	ReorderPenetrationTestReportSections(penetrationTestReportID uint, deletedOrderNum int, tx ...utility.DatabaseTransaction) error
	GetPenetrationTestsByIDs(ids []uint, tx utility.DatabaseTransaction) ([]penetrationTestDTOs.PenetrationTestReportPenetrationTestSourceDTO, error)
	GetReportPenetrationTestsByDateRange(start time.Time, end time.Time, tx utility.DatabaseTransaction) ([]penetrationTestDTOs.PenetrationTestReportPenetrationTestSourceDTO, error)
	GetFindingsByPenetrationTestID(penetrationTestID uint, tx utility.DatabaseTransaction) ([]penetrationTestDTOs.PenetrationTestReportFindingSourceDTO, error)
	GetRetestsByPenetrationTestID(penetrationTestID uint, tx utility.DatabaseTransaction) ([]penetrationTestDTOs.PenetrationTestReportRetestSourceDTO, error)
	GetRetestFindingsByRetestID(retestID uint, tx utility.DatabaseTransaction) ([]penetrationTestDTOs.PenetrationTestReportRetestFindingSourceDTO, error)
	GetDocumentsByRetestFindingID(retestFindingID uint, tx utility.DatabaseTransaction) ([]penetrationTestDTOs.PenetrationTestReportDocumentSourceDTO, error)
	GetDocumentsByPenetrationTestID(penetrationTestID uint, tx utility.DatabaseTransaction) ([]penetrationTestDTOs.PenetrationTestReportDocumentSourceDTO, error)
	GetPenetrationTestReportDetailsByID(id uint) (*models.PenetrationTestReport, error)
	CreateReportPenetrationTest(penetrationTest *models.PenetrationTestReportPenetrationTest, tx utility.DatabaseTransaction) error
	CreateReportFinding(finding *models.PenetrationTestReportFinding, tx utility.DatabaseTransaction) error
	CreateReportRetest(retest *models.PenetrationTestReportRetest, tx utility.DatabaseTransaction) error
	CreateReportRetestFinding(retestFinding *models.PenetrationTestReportRetestFinding, tx utility.DatabaseTransaction) error
	CreatePenetrationTestDocuments(docs []*models.PenetrationTestReportDocument, tx utility.DatabaseTransaction) error
	DeletePenetrationTestReportDocuments(penetrationTestReportID uint, tx ...utility.DatabaseTransaction) error
	DeletePenetrationTestReportRetestFindings(penetrationTestReportID uint, tx ...utility.DatabaseTransaction) error
	DeletePenetrationTestReportRetests(penetrationTestReportID uint, tx ...utility.DatabaseTransaction) error
	DeletePenetrationTestReportFindings(penetrationTestReportID uint, tx ...utility.DatabaseTransaction) error
	DeletePenetrationTestReportPenetrationTests(penetrationTestReportID uint, tx ...utility.DatabaseTransaction) error
	DeletePenetrationTestReportSections(penetrationTestReportID uint, tx ...utility.DatabaseTransaction) error
	DeletePenetrationTestReport(penetrationTestReportID uint, tx ...utility.DatabaseTransaction) error
}

type PenetrationTestReportRepository struct {
	db *gorm.DB
}

func NewPenetrationTestReportRepository(db *gorm.DB) *PenetrationTestReportRepository {
	return &PenetrationTestReportRepository{
		db: db,
	}
}

func (repo *PenetrationTestReportRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *PenetrationTestReportRepository) GetAll(params *penetrationTestDTOs.PenetrationTestReportListQueryParamsDTO) (*penetrationTestDTOs.PenetrationTestReportListResponseDTO, error) {
	var reports []*penetrationTestDTOs.PenetrationTestReportListDTO
	var totalCount int64

	query := repo.db.Table("penetration_test_reports").
		Select(`
			penetration_test_reports.id,
			penetration_test_reports.title,
			CONVERT(VARCHAR, penetration_test_reports.start_date, 23) AS start_date,
			CONVERT(VARCHAR, penetration_test_reports.end_date, 23) AS end_date,
			penetration_test_reports.created_at,
			created_user.id AS created_user_id,
			created_user.first_name_en AS created_user_first_name_en,
			created_user.first_name_ar AS created_user_first_name_ar,
			created_user.last_name_en AS created_user_last_name_en,
			created_user.last_name_ar AS created_user_last_name_ar,
			created_user.email AS created_user_email,
			created_user.initial AS created_user_initial,
			created_user.color AS created_user_color,
			created_user.profile_image_path AS created_user_profile_image_path
		`).
		Joins("LEFT JOIN users AS created_user ON created_user.id = penetration_test_reports.created_by_id")

	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where(columnName+" IN (?)", values)
		}
	}

	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
			(penetration_test_reports.title LIKE ? OR
			created_user.first_name_en LIKE ? OR
			created_user.last_name_en LIKE ? OR
			created_user.first_name_ar LIKE ? OR
			created_user.last_name_ar LIKE ?)`,
			keyword, keyword, keyword, keyword, keyword)
	}

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	sortBy := params.SortBy
	if sortBy == "" {
		sortBy = "penetration_test_reports.id"
	}
	sort := params.Sort
	if sort == "" {
		sort = "DESC"
	}
	query = query.Order(sortBy + " " + sort)

	offset := (params.Page - 1) * params.PageSize
	query = query.Offset(offset).Limit(params.PageSize)

	if err := query.Find(&reports).Error; err != nil {
		return nil, err
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	response := &penetrationTestDTOs.PenetrationTestReportListResponseDTO{
		Items: reports,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}
	return response, nil
}

func (repo *PenetrationTestReportRepository) Create(report *models.PenetrationTestReport, tx utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if tx != nil {
		db = tx.DB()
	}

	if err := db.Create(report).Error; err != nil {
		return 0, err
	}

	return report.ID, nil
}

func (repo *PenetrationTestReportRepository) CreateSection(
	section *models.PenetrationTestReportSection,
	tx utility.DatabaseTransaction,
) (uint, error) {

	db := repo.db
	if tx != nil {
		db = tx.DB()
	}

	if err := db.Create(section).Error; err != nil {
		return 0, err
	}

	return section.ID, nil
}

func (repo *PenetrationTestReportRepository) GetPenetrationTestReportByID(penetrationTestReportID uint, tx ...utility.DatabaseTransaction) (*models.PenetrationTestReport, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var report models.PenetrationTestReport
	err := db.Where("id = ?", penetrationTestReportID).First(&report).Error
	if err != nil {
		return nil, err
	}

	return &report, nil
}

func (repo *PenetrationTestReportRepository) GetSectionByID(sectionID uint, tx ...utility.DatabaseTransaction) (*models.PenetrationTestReportSection, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var section models.PenetrationTestReportSection
	err := db.Where("id = ?", sectionID).First(&section).Error
	if err != nil {
		return nil, err
	}
	return &section, nil
}

func (repo *PenetrationTestReportRepository) UpdateSection(section *models.PenetrationTestReportSection, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Save(section).Error
}

func (repo *PenetrationTestReportRepository) DeletePenetrationTestReportSection(sectionID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Where("id = ?", sectionID).Delete(&models.PenetrationTestReportSection{}).Error
}

func (repo *PenetrationTestReportRepository) ReorderPenetrationTestReportSections(penetrationTestReportID uint, deletedOrderNum int, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	err := db.Exec(`
		UPDATE penetration_test_report_sections
		SET order_num = order_num - 1
		WHERE penetration_test_report_id = ? AND order_num > ?`,
		penetrationTestReportID, deletedOrderNum).Error

	return err
}

func (repo *PenetrationTestReportRepository) GetPenetrationTestsByIDs(ids []uint, tx utility.DatabaseTransaction) ([]penetrationTestDTOs.PenetrationTestReportPenetrationTestSourceDTO, error) {
	db := repo.db
	if tx != nil {
		db = tx.DB()
	}
	var items []penetrationTestDTOs.PenetrationTestReportPenetrationTestSourceDTO
	query := db.Table("penetration_tests AS pt").Select(`
			pt.id,
			pt.reference_number AS reference_number,
			pt.system_application_name AS system_application_name,
			pt.scope AS scope,
			pt.year AS year,
			pt.period AS period,
			ptt.title_en AS penetration_test_type_title_en,
			ptt.title_ar AS penetration_test_type_title_ar,
			test_plan.title_en AS penetration_test_plan_title_en,
			test_plan.title_ar AS penetration_test_plan_title_ar,
			classification.title_en AS penetration_test_classification_title_en,
			classification.title_ar AS penetration_test_classification_title_ar,
			executor.title_en AS penetration_test_executor_title_en,
			executor.title_ar AS penetration_test_executor_title_ar,
			CONCAT(employee.first_name_en, ' ', employee.last_name_en) AS cs_responsible_employee_name_en,
			CONCAT(employee.first_name_ar, ' ', employee.last_name_ar) AS cs_responsible_employee_name_ar,
			status.title_en AS penetration_test_status_title_en,
			status.title_ar AS penetration_test_status_title_ar,
			status.color_code AS penetration_test_status_color_code,
			readiness.title_en AS environment_readiness_title_en,
			readiness.title_ar AS environment_readiness_title_ar,
			environment.environment_notify_date AS environment_notify_date,
			environment.expected_ready_date AS expected_ready_date,
			execution.test_start_date AS execution_start_date,
			execution.test_end_date AS execution_end_date,
			execution_status.title_en AS execution_status_title_en,
			execution_status.title_ar AS execution_status_title_ar,
			security_result.title_en AS security_result_title_en,
			security_result.title_ar AS security_result_title_ar,
			COALESCE(management.results_shared_with_management, 0) AS results_shared_with_management,
			management.management_notify_date AS management_notify_date,
			management.management_response_date AS management_response_date,
			management.remediation_notes AS remediation_notes,
			closure.closure_date AS closure_date,
			closure.closure_remarks AS closure_remarks,
			CASE WHEN closure_owner.id IS NOT NULL THEN CONCAT(closure_owner.first_name_en, ' ', closure_owner.last_name_en) END AS closure_owner_name_en,
			CASE WHEN closure_owner.id IS NOT NULL THEN CONCAT(closure_owner.first_name_ar, ' ', closure_owner.last_name_ar) END AS closure_owner_name_ar,
			pt.created_at AS created_at,
			pt.updated_at AS updated_at
		`).
		Joins("LEFT JOIN penetration_test_types AS ptt ON ptt.id = pt.penetration_test_type_id").
		Joins("LEFT JOIN penetration_test_plans AS test_plan ON test_plan.id = pt.penetration_test_plan_id").
		Joins("LEFT JOIN penetration_test_classifications AS classification ON classification.id = pt.penetration_test_classification_id").
		Joins("LEFT JOIN penetration_test_executors AS executor ON executor.id = pt.penetration_test_executor_id").
		Joins("LEFT JOIN penetration_test_statuses AS status ON status.id = pt.penetration_test_status_id").
		Joins("LEFT JOIN penetration_test_environment_preparations AS environment ON environment.penetration_test_id = pt.id").
		Joins("LEFT JOIN penetration_test_environment_readinesses AS readiness ON readiness.id = environment.environment_readiness_id").
		Joins("LEFT JOIN penetration_test_executions AS execution ON execution.penetration_test_id = pt.id").
		Joins("LEFT JOIN penetration_test_execution_statuses AS execution_status ON execution_status.id = execution.execution_status_id").
		Joins("LEFT JOIN penetration_test_security_results AS security_result ON security_result.id = execution.security_result_id").
		Joins("LEFT JOIN users AS employee ON employee.id = pt.cs_responsible_employee_id").
		Joins("LEFT JOIN penetration_test_management_sharings AS management ON management.penetration_test_id = pt.id").
		Joins("LEFT JOIN penetration_test_closures AS closure ON closure.penetration_test_id = pt.id").
		Joins("LEFT JOIN users AS closure_owner ON closure_owner.id = closure.closure_owner_id")
	query = query.Where("pt.id IN ?", ids)
	if err := query.Order("pt.created_at DESC").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (repo *PenetrationTestReportRepository) GetReportPenetrationTestsByDateRange(start time.Time, end time.Time, tx utility.DatabaseTransaction) ([]penetrationTestDTOs.PenetrationTestReportPenetrationTestSourceDTO, error) {
	db := repo.db
	if tx != nil {
		db = tx.DB()
	}
	var items []penetrationTestDTOs.PenetrationTestReportPenetrationTestSourceDTO
	query := db.Table("penetration_tests AS pt").Select(`
			pt.id,
			pt.reference_number AS reference_number,
			pt.system_application_name AS system_application_name,
			pt.scope AS scope,
			pt.year AS year,
			pt.period AS period,
			ptt.title_en AS penetration_test_type_title_en,
			ptt.title_ar AS penetration_test_type_title_ar,
			test_plan.title_en AS penetration_test_plan_title_en,
			test_plan.title_ar AS penetration_test_plan_title_ar,
			classification.title_en AS penetration_test_classification_title_en,
			classification.title_ar AS penetration_test_classification_title_ar,
			executor.title_en AS penetration_test_executor_title_en,
			executor.title_ar AS penetration_test_executor_title_ar,
			CONCAT(employee.first_name_en, ' ', employee.last_name_en) AS cs_responsible_employee_name_en,
			CONCAT(employee.first_name_ar, ' ', employee.last_name_ar) AS cs_responsible_employee_name_ar,
			status.title_en AS penetration_test_status_title_en,
			status.title_ar AS penetration_test_status_title_ar,
			status.color_code AS penetration_test_status_color_code,
			readiness.title_en AS environment_readiness_title_en,
			readiness.title_ar AS environment_readiness_title_ar,
			environment.environment_notify_date AS environment_notify_date,
			environment.expected_ready_date AS expected_ready_date,
			execution.test_start_date AS execution_start_date,
			execution.test_end_date AS execution_end_date,
			execution_status.title_en AS execution_status_title_en,
			execution_status.title_ar AS execution_status_title_ar,
			security_result.title_en AS security_result_title_en,
			security_result.title_ar AS security_result_title_ar,
			COALESCE(management.results_shared_with_management, 0) AS results_shared_with_management,
			management.management_notify_date AS management_notify_date,
			management.management_response_date AS management_response_date,
			management.remediation_notes AS remediation_notes,
			closure.closure_date AS closure_date,
			closure.closure_remarks AS closure_remarks,
			CASE WHEN closure_owner.id IS NOT NULL THEN CONCAT(closure_owner.first_name_en, ' ', closure_owner.last_name_en) END AS closure_owner_name_en,
			CASE WHEN closure_owner.id IS NOT NULL THEN CONCAT(closure_owner.first_name_ar, ' ', closure_owner.last_name_ar) END AS closure_owner_name_ar,
			pt.created_at AS created_at,
			pt.updated_at AS updated_at
		`).
		Joins("LEFT JOIN penetration_test_types AS ptt ON ptt.id = pt.penetration_test_type_id").
		Joins("LEFT JOIN penetration_test_plans AS test_plan ON test_plan.id = pt.penetration_test_plan_id").
		Joins("LEFT JOIN penetration_test_classifications AS classification ON classification.id = pt.penetration_test_classification_id").
		Joins("LEFT JOIN penetration_test_executors AS executor ON executor.id = pt.penetration_test_executor_id").
		Joins("LEFT JOIN penetration_test_statuses AS status ON status.id = pt.penetration_test_status_id").
		Joins("LEFT JOIN penetration_test_environment_preparations AS environment ON environment.penetration_test_id = pt.id").
		Joins("LEFT JOIN penetration_test_environment_readinesses AS readiness ON readiness.id = environment.environment_readiness_id").
		Joins("LEFT JOIN penetration_test_executions AS execution ON execution.penetration_test_id = pt.id").
		Joins("LEFT JOIN penetration_test_execution_statuses AS execution_status ON execution_status.id = execution.execution_status_id").
		Joins("LEFT JOIN penetration_test_security_results AS security_result ON security_result.id = execution.security_result_id").
		Joins("LEFT JOIN users AS employee ON employee.id = pt.cs_responsible_employee_id").
		Joins("LEFT JOIN penetration_test_management_sharings AS management ON management.penetration_test_id = pt.id").
		Joins("LEFT JOIN penetration_test_closures AS closure ON closure.penetration_test_id = pt.id").
		Joins("LEFT JOIN users AS closure_owner ON closure_owner.id = closure.closure_owner_id")
	query = query.Where("pt.created_at BETWEEN ? AND ?", start, end)
	if err := query.Order("pt.created_at DESC").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (repo *PenetrationTestReportRepository) GetFindingsByPenetrationTestID(penetrationTestID uint, tx utility.DatabaseTransaction) ([]penetrationTestDTOs.PenetrationTestReportFindingSourceDTO, error) {
	db := repo.db
	if tx != nil {
		db = tx.DB()
	}
	var items []penetrationTestDTOs.PenetrationTestReportFindingSourceDTO
	query := db.Table("penetration_test_findings AS finding").Select(`
			finding.id,
			finding.reference_number AS reference_number,
			finding.title AS title,
			finding.description AS description,
			finding.recommendation AS recommendation,
			severity.title_en AS severity_level_title_en,
			severity.title_ar AS severity_level_title_ar,
			severity.color_code AS severity_level_color_code,
			CONCAT(owner.first_name_en, ' ', owner.last_name_en) AS owner_name_en,
			CONCAT(owner.first_name_ar, ' ', owner.last_name_ar) AS owner_name_ar,
			status.title_en AS action_plan_status_title_en,
			status.title_ar AS action_plan_status_title_ar,
			status.color_code AS action_plan_status_color_code,
			finding.percentage AS percentage,
			finding.target_remediation_date AS target_remediation_date,
			finding.created_at AS created_at,
			finding.updated_at AS updated_at
		`).
		Joins("LEFT JOIN severity_levels AS severity ON severity.id = finding.severity_level_id").
		Joins("LEFT JOIN action_plan_statuses AS status ON status.id = finding.action_plan_status_id").
		Joins("LEFT JOIN users AS owner ON owner.id = finding.owner_id")
	query = query.Where("finding.penetration_test_id = ?", penetrationTestID)
	if err := query.Order("finding.id ASC").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (repo *PenetrationTestReportRepository) GetRetestsByPenetrationTestID(penetrationTestID uint, tx utility.DatabaseTransaction) ([]penetrationTestDTOs.PenetrationTestReportRetestSourceDTO, error) {
	db := repo.db
	if tx != nil {
		db = tx.DB()
	}
	var items []penetrationTestDTOs.PenetrationTestReportRetestSourceDTO
	query := db.Table("penetration_test_retests AS retest").Select(`
			retest.id,
			retest.reference_number AS reference_number,
			retest.start_date AS start_date,
			retest.end_date AS end_date,
			executor.title_en AS executor_title_en,
			executor.title_ar AS executor_title_ar,
			status.title_en AS status_title_en,
			status.title_ar AS status_title_ar,
			status.color_code AS status_color_code,
			retest.notes AS notes,
			retest.created_at AS created_at,
			retest.updated_at AS updated_at
		`).
		Joins("LEFT JOIN penetration_test_executors AS executor ON executor.id = retest.penetration_test_executor_id").
		Joins("LEFT JOIN penetration_test_retest_statuses AS status ON status.id = retest.penetration_test_retest_status_id")
	query = query.Where("retest.penetration_test_id = ?", penetrationTestID)
	if err := query.Order("retest.id ASC").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (repo *PenetrationTestReportRepository) GetRetestFindingsByRetestID(retestID uint, tx utility.DatabaseTransaction) ([]penetrationTestDTOs.PenetrationTestReportRetestFindingSourceDTO, error) {
	db := repo.db
	if tx != nil {
		db = tx.DB()
	}
	var findings []penetrationTestDTOs.PenetrationTestReportRetestFindingSourceDTO
	err := db.Table("penetration_test_retest_findings AS retest_finding").
		Select(`
			retest_finding.id,
			retest_finding.penetration_test_finding_id,
			retest_finding.comment,
			retest_finding.comment_at,
			retest_finding.result_at,
			result.title_en AS result_title_en,
			result.title_ar AS result_title_ar
		`).
		Joins("LEFT JOIN penetration_test_retest_results AS result ON result.id = retest_finding.retest_result_id").
		Where("retest_finding.penetration_test_retest_id = ?", retestID).
		Order("retest_finding.id ASC").
		Find(&findings).Error
	if err != nil {
		return nil, err
	}
	return findings, nil
}

func (repo *PenetrationTestReportRepository) GetDocumentsByRetestFindingID(retestFindingID uint, tx utility.DatabaseTransaction) ([]penetrationTestDTOs.PenetrationTestReportDocumentSourceDTO, error) {
	db := repo.db
	if tx != nil {
		db = tx.DB()
	}
	var documents []penetrationTestDTOs.PenetrationTestReportDocumentSourceDTO
	err := db.Table("penetration_test_retest_finding_documents").
		Select(`
			penetration_test_retest_finding_documents.id,
			penetration_test_retest_finding_documents.penetration_test_retest_finding_id,
			penetration_test_retest_finding_documents.url,
			penetration_test_retest_finding_documents.name,
			penetration_test_retest_finding_documents.size,
			penetration_test_retest_finding_documents.token,
			penetration_test_retest_finding_documents.original_name,
			penetration_test_retest_finding_documents.type,
			penetration_test_retest_finding_documents.created_by_id,
			penetration_test_retest_finding_documents.updated_by_id,
			penetration_test_retest_finding_documents.created_at,
			penetration_test_retest_finding_documents.updated_at
		`).
		Where("penetration_test_retest_finding_id = ?", retestFindingID).
		Order("id ASC").Find(&documents).Error
	if err != nil {
		return nil, err
	}
	return documents, nil
}

func (repo *PenetrationTestReportRepository) GetDocumentsByPenetrationTestID(penetrationTestID uint, tx utility.DatabaseTransaction) ([]penetrationTestDTOs.PenetrationTestReportDocumentSourceDTO, error) {
	db := repo.db
	if tx != nil {
		db = tx.DB()
	}
	var documents []penetrationTestDTOs.PenetrationTestReportDocumentSourceDTO
	err := db.Table("penetration_test_supporting_document_documents").
		Select(`
			penetration_test_supporting_document_documents.id,
			penetration_test_supporting_document_documents.penetration_test_supporting_document_id,
			penetration_test_supporting_document_documents.url,
			penetration_test_supporting_document_documents.name,
			penetration_test_supporting_document_documents.size,
			penetration_test_supporting_document_documents.token,
			penetration_test_supporting_document_documents.original_name,
			penetration_test_supporting_document_documents.type,
			penetration_test_supporting_document_documents.created_by_id,
			penetration_test_supporting_document_documents.updated_by_id,
			penetration_test_supporting_document_documents.created_at,
			penetration_test_supporting_document_documents.updated_at
		`).
		Where("penetration_test_supporting_document_id IN (SELECT id FROM penetration_test_supporting_documents WHERE penetration_test_id = ?)", penetrationTestID).
		Order("id ASC").
		Find(&documents).Error
	if err != nil {
		return nil, err
	}
	return documents, nil
}

func (repo *PenetrationTestReportRepository) GetPenetrationTestReportDetailsByID(id uint) (*models.PenetrationTestReport, error) {
	var report models.PenetrationTestReport
	err := repo.db.
		Preload("CreatedByUser").
		Preload("Sections").
		Preload("PenetrationTests").
		Preload("PenetrationTests.PenetrationTest").
		Preload("PenetrationTests.PenetrationTest.Closure").
		Preload("PenetrationTests.PenetrationTest.Closure.ClosureOwner").
		Preload("PenetrationTests.PenetrationTest.Closure.CreatedByUser").
		Preload("PenetrationTests.PenetrationTest.Closure.UpdatedByUser").
		Preload("PenetrationTests.Findings").
		Preload("PenetrationTests.Retests").
		Preload("PenetrationTests.Retests.RetestFindings").
		Preload("PenetrationTests.Retests.RetestFindings.Documents").
		Preload("PenetrationTests.Documents").
		Where("id = ?", id).
		First(&report).Error
	if err != nil {
		return nil, err
	}
	return &report, nil
}

func (repo *PenetrationTestReportRepository) CreateReportPenetrationTest(penetrationTest *models.PenetrationTestReportPenetrationTest, tx utility.DatabaseTransaction) error {
	db := repo.db
	if tx != nil {
		db = tx.DB()
	}
	if err := db.Create(penetrationTest).Error; err != nil {
		return err
	}

	return nil
}

func (repo *PenetrationTestReportRepository) CreateReportFinding(finding *models.PenetrationTestReportFinding, tx utility.DatabaseTransaction) error {
	db := repo.db
	if tx != nil {
		db = tx.DB()
	}
	if err := db.Create(finding).Error; err != nil {
		return err
	}

	return nil
}

func (repo *PenetrationTestReportRepository) CreateReportRetest(retest *models.PenetrationTestReportRetest, tx utility.DatabaseTransaction) error {
	db := repo.db
	if tx != nil {
		db = tx.DB()
	}
	if err := db.Create(retest).Error; err != nil {
		return err
	}

	return nil
}

func (repo *PenetrationTestReportRepository) CreateReportRetestFinding(retestFinding *models.PenetrationTestReportRetestFinding, tx utility.DatabaseTransaction) error {
	db := repo.db
	if tx != nil {
		db = tx.DB()
	}
	if err := db.Create(retestFinding).Error; err != nil {
		return err
	}

	return nil
}

func (repo *PenetrationTestReportRepository) CreatePenetrationTestDocuments(
	docs []*models.PenetrationTestReportDocument,
	tx utility.DatabaseTransaction,
) error {
	db := repo.db
	if tx != nil {
		db = tx.DB()
	}

	if err := db.Create(&docs).Error; err != nil {
		return err
	}
	return nil
}

func (repo *PenetrationTestReportRepository) DeletePenetrationTestReportDocuments(penetrationTestReportID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Where("penetration_test_report_penetration_test_id IN (SELECT id FROM penetration_test_report_penetration_tests WHERE penetration_test_report_id = ?)", penetrationTestReportID).Delete(&models.PenetrationTestReportDocument{}).Error
}

func (repo *PenetrationTestReportRepository) DeletePenetrationTestReportRetestFindings(penetrationTestReportID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Where("penetration_test_report_retest_id IN (SELECT id FROM penetration_test_report_retests WHERE penetration_test_report_id = ?)", penetrationTestReportID).Delete(&models.PenetrationTestReportRetestFinding{}).Error
}

func (repo *PenetrationTestReportRepository) DeletePenetrationTestReportRetests(penetrationTestReportID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Where("penetration_test_report_id = ?", penetrationTestReportID).Delete(&models.PenetrationTestReportRetest{}).Error
}

func (repo *PenetrationTestReportRepository) DeletePenetrationTestReportFindings(penetrationTestReportID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Where("penetration_test_report_id = ?", penetrationTestReportID).Delete(&models.PenetrationTestReportFinding{}).Error
}

func (repo *PenetrationTestReportRepository) DeletePenetrationTestReportPenetrationTests(penetrationTestReportID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Where("penetration_test_report_id = ?", penetrationTestReportID).Delete(&models.PenetrationTestReportPenetrationTest{}).Error
}

func (repo *PenetrationTestReportRepository) DeletePenetrationTestReportSections(penetrationTestReportID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Where("penetration_test_report_id = ?", penetrationTestReportID).Delete(&models.PenetrationTestReportSection{}).Error
}

func (repo *PenetrationTestReportRepository) DeletePenetrationTestReport(penetrationTestReportID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Where("id = ?", penetrationTestReportID).Delete(&models.PenetrationTestReport{}).Error
}

func (repo *PenetrationTestReportRepository) GetFullReportDetailsByID(id uint) (*models.PenetrationTest, error) {
	var penetrationTest models.PenetrationTest
	if err := repo.db.Preload("CreatedByUser").Where("id = ?", id).First(&penetrationTest).Error; err != nil {
		return nil, err
	}
	return &penetrationTest, nil
}
