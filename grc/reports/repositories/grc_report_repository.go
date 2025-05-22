package repositories

import (
	"fmt"
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/models"
	"jethur/grc/reports/dtos"

	"jethur/utility"
	"time"

	"gorm.io/gorm"
)

type GRCReportRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction

	GetAll(params *dtos.GRCReportListQueryParamsDTO) (*dtos.GRCReportListResponseDTO, error)

	GetGRCReportByID(grcReportID uint, tx ...utility.DatabaseTransaction) (*dtos.GRCReportListDTO, error)
	GetGRCReportSections(grcReportID uint) ([]dtos.GRCReportSectionResponseDTO, error)
	GetSectionByID(sectionID uint, tx ...utility.DatabaseTransaction) (*models.GRCReportSection, error)
	GetGRCReportGovernanceDocuments(reportID uint) ([]dtos.GRCGovernanceDocumentResponseDTO, error)
	GetGRCReportGovernanceTasks(grcReportID uint) ([]dtos.GRCGovernanceTaskResponseDTO, error)
	GetGRCReportGovernanceSummary(grcReportID uint) (*dtos.GRCGovernanceSummaryResponseDTO, error)
	GetGRCReportGovernanceSummaryByID(grcReportID uint, tx ...utility.DatabaseTransaction) (*models.GRCReportGovernanceSummary, error)

	AllUserCount(tx ...utility.DatabaseTransaction) (int64, error)
	GetPolicyDocumentByDateRange(startDate, endDate time.Time, tx ...utility.DatabaseTransaction) ([]*models.Document, error)
	GetDocumentAcknowledgmentSummary(documentID uint, primaryOrgUnitField string) (*dtos.DocumentAcknowledgmentSummaryDTO, error)
	GetPublicDocumentAcknowledgmentSummary(documentID uint) (*dtos.DocumentAcknowledgmentSummaryDTO, error)
	GetTasksByDateRange(startDate, endDate time.Time, tx ...utility.DatabaseTransaction) ([]*models.Task, error)

	Create(grcReport *models.GRCReport, tx ...utility.DatabaseTransaction) (uint, error)
	CreateSection(section *models.GRCReportSection, tx ...utility.DatabaseTransaction) error
	CreateGRCReportGovernanceSummary(summary *models.GRCReportGovernanceSummary, tx ...utility.DatabaseTransaction) error
	CreateGRCReportGovernanceDocument(document *models.GRCReportGovernanceDocument, tx ...utility.DatabaseTransaction) error
	CreateGRCReportGovernanceTask(task *models.GRCReportGovernanceTask, tx ...utility.DatabaseTransaction) error

	UpdateSection(section *models.GRCReportSection, tx ...utility.DatabaseTransaction) error
	UpdateGRCReportGovernanceSummary(governanceSummary *models.GRCReportGovernanceSummary, tx ...utility.DatabaseTransaction) error

	DeleteGRCReport(grcReportID uint, tx ...utility.DatabaseTransaction) error
	DeleteGRCReportSection(sectionID uint, tx ...utility.DatabaseTransaction) error
	DeleteGRCReportSections(grcReportID uint, tx ...utility.DatabaseTransaction) error
	DeleteGRCReportGovernanceDocuments(grcReportID uint, tx ...utility.DatabaseTransaction) error
	DeleteGRCReportGovernanceTasks(grcReportID uint, tx ...utility.DatabaseTransaction) error

	ReorderGRCReportSections(grcReportID uint, deletedOrderNum int, tx ...utility.DatabaseTransaction) error

	// risk section
	GetGRCReportRiskTopFive(grcReportID uint) ([]dtos.GRCReportRiskTopFiveRiskResponseDTO, error)
	GetGRCReportRiskKRIs(grcReportID uint) ([]dtos.GRCReportRiskKriResponseDTO, error)
	GetGRCReportRiskControls(grcReportID uint) ([]dtos.GRCReportRiskControlResponseDTO, error)
	GetGRCReportRiskSummary(grcReportID uint) (*dtos.GRCRiskSummaryResponseDTO, error)
	GetRiskTrend(startDate time.Time, endDate time.Time, tx ...utility.DatabaseTransaction) ([]dtos.RiskTrendDTO, error)

	CreateGRCReportRiskSummary(summary *models.GRCReportRiskSummary, tx ...utility.DatabaseTransaction) error
	GetRisksByDateRange(startDate, endDate time.Time, tx ...utility.DatabaseTransaction) ([]*models.RiskRegister, error)
	GetHighestActiveRiskRatingType(tx ...utility.DatabaseTransaction) (string, error)
	GetTopFiveRisksByRatingValue(startDate time.Time, endDate time.Time, tx ...utility.DatabaseTransaction) ([]dtos.TopFiveRiskDTO, error)
	CreateGRCReportRiskTopFive(topFive *models.GRCReportRiskTopFive, tx ...utility.DatabaseTransaction) error
	GetControlTitlesByRiskID(riskID uint, tx ...utility.DatabaseTransaction) ([]string, error)
	GetKRIsByDateRange(startDate, endDate time.Time, tx ...utility.DatabaseTransaction) ([]dtos.GRCReportRiskKRIDTO, error)
	CreateGRCReportRiskKRI(riskKri *models.GRCReportRiskKRI, tx ...utility.DatabaseTransaction) error
	GetRiskControlsByDateRange(startDate, endDate time.Time, tx ...utility.DatabaseTransaction) ([]dtos.RiskControlReportDTO, error)
	CreateGRCReportRiskControl(control *models.GRCReportRiskControl, tx ...utility.DatabaseTransaction) error

	// compliance and sla section
	GetAllSLAs(tx ...utility.DatabaseTransaction) ([]*models.SLARegister, error)
	CreateGRCReportSLA(sla *models.GRCReportSLA, tx ...utility.DatabaseTransaction) (uint, error)
	GetAllComplianceAssessments(startDate time.Time, endDate time.Time, tx ...utility.DatabaseTransaction) ([]*models.ComplianceAssessment, error)
	CreateGRCReportComplianceSummary(summary *models.GRCReportComplianceSummary, tx ...utility.DatabaseTransaction) (uint, error)
	CreateGRCReportComplianceDepartment(department *models.GRCReportComplianceDepartment, tx ...utility.DatabaseTransaction) (uint, error)
	GetOrganizationUnitByID(id uint) (*models.OrganizationUnit, error)

	GetGRCReportComplianceSummary(grcReportID uint) ([]*models.GRCReportComplianceSummary, error)
	GetGRCReportSLA(grcReportID uint) ([]*dtos.GRCReportSLAResponseDTO, error)
}

type GRCReportRepository struct {
	db *gorm.DB
}

func NewGRCReportRepository(db *gorm.DB) GRCReportRepositoryInterface {
	return &GRCReportRepository{db: db}
}

func (repo *GRCReportRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *GRCReportRepository) GetAll(params *dtos.GRCReportListQueryParamsDTO) (*dtos.GRCReportListResponseDTO, error) {
	var reports []*dtos.GRCReportListDTO
	var totalCount int64

	query := repo.db.Table("grc_reports").
		Select(`
			grc_reports.id, 
			grc_reports.title, 
			CONVERT(VARCHAR, grc_reports.start_date, 23) AS start_date,
			CONVERT(VARCHAR, grc_reports.end_date, 23) AS end_date,
			grc_reports.created_at,
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
		Joins("LEFT JOIN users AS created_user ON created_user.id = grc_reports.created_by_id")

	if params.Keyword != "" {
		query = query.Where(`(grc_reports.title LIKE ? )`, "%"+params.Keyword+"%")
	}

	query = query.Order(params.SortBy + " " + params.Sort)

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	offset := (params.Page - 1) * params.PageSize
	query = query.Offset(offset).Limit(params.PageSize)

	if err := query.Find(&reports).Error; err != nil {
		return nil, err
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	response := &dtos.GRCReportListResponseDTO{
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

func (repo *GRCReportRepository) GetSectionByID(sectionID uint, tx ...utility.DatabaseTransaction) (*models.GRCReportSection, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var grcReportSection models.GRCReportSection
	err := db.
		Where("id = ?", sectionID).First(&grcReportSection).Error
	if err != nil {
		return nil, err
	}
	return &grcReportSection, nil
}

func (repo *GRCReportRepository) GetGRCReportByID(grcReportID uint, tx ...utility.DatabaseTransaction) (*dtos.GRCReportListDTO, error) {
	var report dtos.GRCReportListDTO

	query := repo.db.Table("grc_reports").
		Select(`
			grc_reports.id, 
			grc_reports.title, 
			CONVERT(VARCHAR, grc_reports.start_date, 23) AS start_date,
			CONVERT(VARCHAR, grc_reports.end_date, 23) AS end_date,
			grc_reports.created_at,
			users.id AS created_by_id,
			users.first_name_en AS created_user_first_name_en,
			users.first_name_ar AS created_user_first_name_ar,
			users.last_name_en AS created_user_last_name_en,
			users.last_name_ar AS created_user_last_name_ar,
			users.email AS created_user_email
		`).
		Joins("LEFT JOIN users ON users.id = grc_reports.created_by_id").
		Where("grc_reports.id = ?", grcReportID).
		Group("grc_reports.id, grc_reports.title, grc_reports.start_date, grc_reports.end_date, grc_reports.created_at, users.id, users.first_name_en, users.first_name_ar, users.last_name_en, users.last_name_ar, users.email")

	if err := query.First(&report).Error; err != nil {

		return nil, err
	}

	return &report, nil
}

// Fetch sections for the report
func (repo *GRCReportRepository) GetGRCReportSections(grcReportID uint) ([]dtos.GRCReportSectionResponseDTO, error) {
	var sections []dtos.GRCReportSectionResponseDTO
	err := repo.db.Table("grc_report_sections").
		Select("id, type, title_en, title_ar, description_en, description_ar, order_num").
		Where("grc_report_id = ?", grcReportID).
		Order("order_num ASC").
		Find(&sections).Error

	if err != nil {
		return nil, err
	}
	return sections, nil
}

func (repo *GRCReportRepository) GetGRCReportGovernanceDocuments(grcReportID uint) ([]dtos.GRCGovernanceDocumentResponseDTO, error) {
	var documents = []dtos.GRCGovernanceDocumentResponseDTO{}

	err := repo.db.Table("grc_report_governance_documents").
		Select(`
			id, grc_report_id, reference_number,
			document_title_en, document_title_ar,
			responsible_user_name_en, responsible_user_name_ar,
			responsible_unit_title_en, responsible_unit_title_ar,
			total_employees, acknowledged_count, pending_count,
			document_status_title_en, document_status_title_ar, document_status_color_code,acknowledge_due_date,
			created_by_id, created_at
		`).
		Where("grc_report_id = ?", grcReportID).
		Order("created_at ASC").
		Scan(&documents).Error

	return documents, err
}

func (repo *GRCReportRepository) GetGRCReportGovernanceTasks(grcReportID uint) ([]dtos.GRCGovernanceTaskResponseDTO, error) {
	var tasks = []dtos.GRCGovernanceTaskResponseDTO{}

	err := repo.db.Table("grc_report_governance_tasks").
		Select(`
			grc_report_governance_tasks.id,
			grc_report_governance_tasks.grc_report_id,
			grc_report_governance_tasks.task_title,
			grc_report_governance_tasks.task_description,
			grc_report_governance_tasks.responsible_user_name_en,
			grc_report_governance_tasks.responsible_user_name_ar,
			grc_report_governance_tasks.task_status,
			grc_report_governance_tasks.task_status_title_en AS task_status_title_en,
			grc_report_governance_tasks.task_status_title_ar AS task_status_title_ar,
			grc_report_governance_tasks.task_status_color_code AS task_status_color_code,
			grc_report_governance_tasks.severity_level_title_en AS severity_level_title_en,
			grc_report_governance_tasks.severity_level_title_ar AS severity_level_title_ar,
			grc_report_governance_tasks.severity_level_color_code AS severity_level_color_code,
			CONVERT(VARCHAR, grc_report_governance_tasks.target_date, 23) AS target_date,
			grc_report_governance_tasks.created_by_id,
			CONVERT(VARCHAR, grc_report_governance_tasks.created_at, 120) AS created_at,
			created_user.first_name_en AS created_user_first_name_en,
			created_user.first_name_ar AS created_user_first_name_ar,
			created_user.last_name_en AS created_user_last_name_en,
			created_user.last_name_ar AS created_user_last_name_ar,
			created_user.email AS created_user_email,
			created_user.initial AS created_user_initial,
			created_user.color AS created_user_color,
			created_user.profile_image_path AS created_user_profile_image_path
		`).
		Joins("LEFT JOIN users AS created_user ON created_user.id = grc_report_governance_tasks.created_by_id").
		Where("grc_report_governance_tasks.grc_report_id = ?", grcReportID).
		Order("grc_report_governance_tasks.target_date ASC").
		Scan(&tasks).Error

	return tasks, err
}

func (repo *GRCReportRepository) GetGRCReportGovernanceSummary(grcReportID uint) (*dtos.GRCGovernanceSummaryResponseDTO, error) {
	var summary dtos.GRCGovernanceSummaryResponseDTO

	err := repo.db.Table("grc_report_governance_summary").
		Select(`
			id,
			grc_report_id AS reference_number,
			policies_in_system,
			total_policies,
			roles_defined,
			policy_approved_rate,
			highlights,
			challenges,
			document_overview_chart,
			document_type_chart,
			top_published_documents,
			created_by_id,
			CONVERT(VARCHAR, created_at, 120) AS created_at
		`).
		Where("grc_report_id = ?", grcReportID).
		First(&summary).Error

	if err != nil {
		return nil, err
	}

	return &summary, nil
}

func (repo *GRCReportRepository) Create(grcReport *models.GRCReport, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(grcReport).Error; err != nil {
		return 0, err
	}
	return grcReport.ID, nil
}

func (repo *GRCReportRepository) GetPolicyDocumentByDateRange(startDate, endDate time.Time, tx ...utility.DatabaseTransaction) ([]*models.Document, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var documents []*models.Document
	err := db.
		Joins("JOIN document_types dt ON dt.id = documents.document_type_id").
		Joins("JOIN document_versions dv ON dv.document_id = documents.id").
		Where("dt.type = ?", "policy").
		Where("dv.acknowledge_due_date BETWEEN ? AND ?", startDate, endDate).
		Where("dv.is_latest_version = ?", true).
		Preload("ResponsibleUser").
		Preload("ResponsibleUnit").
		Preload("DocumentAccessType").
		Preload("DocumentStatus").
		Preload("DocumentType").
		Preload("DocumentCategory").
		Preload("LatestDocumentVersion", "is_latest_version = ?", true).
		Find(&documents).Error

	if err != nil {
		return nil, err
	}

	return documents, nil
}

func (repo *GRCReportRepository) CreateSection(section *models.GRCReportSection, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Create(&section).Error
}

func (repo *GRCReportRepository) UpdateSection(section *models.GRCReportSection, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Save(section).Error
}

func (repo *GRCReportRepository) DeleteGRCReportSections(grcReportID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Where("grc_report_id = ?", grcReportID).Delete(&models.GRCReportSection{}).Error
}

func (repo *GRCReportRepository) DeleteGRCReportGovernanceDocuments(grcReportID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Where("grc_report_id = ?", grcReportID).Delete(&models.GRCReportGovernanceDocument{}).Error
}

func (repo *GRCReportRepository) DeleteGRCReportGovernanceTasks(grcReportID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Where("grc_report_id = ?", grcReportID).Delete(&models.GRCReportGovernanceTask{}).Error
}

func (repo *GRCReportRepository) DeleteGRCReport(grcReportID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Where("id = ?", grcReportID).Delete(&models.GRCReport{}).Error
}

func (repo *GRCReportRepository) DeleteGRCReportSection(sectionID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Where("id = ?", sectionID).Delete(&models.GRCReportSection{}).Error
}

func (repo *GRCReportRepository) ReorderGRCReportSections(grcReportID uint, deletedOrderNum int, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	// Shift order numbers of remaining sections
	err := db.Exec(`
		UPDATE grc_report_sections
		SET order_num = order_num - 1
		WHERE grc_report_id = ? AND order_num > ?`,
		grcReportID, deletedOrderNum).Error

	return err
}

func (repo *GRCReportRepository) GetDocumentAcknowledgmentSummary(documentID uint, primaryOrgUnitField string) (*dtos.DocumentAcknowledgmentSummaryDTO, error) {
	var totalUsers int64
	var acknowledgedUsers int64

	totalUsersQuery := fmt.Sprintf(`
		SELECT COUNT(DISTINCT user_id) FROM (
			SELECT user_id 
			FROM document_shared_users 
			JOIN users ON users.id = document_shared_users.user_id
			JOIN user_statuses ON user_statuses.id = users.user_status_id
			WHERE document_shared_users.document_id = ? AND user_statuses.is_active = 1
			UNION
			SELECT users.id AS user_id 
			FROM document_shared_units
			JOIN organization_units ON organization_units.id = document_shared_units.organization_unit_id
			JOIN users ON users.%s = organization_units.id
			JOIN user_statuses ON user_statuses.id = users.user_status_id
			WHERE document_shared_units.document_id = ? AND user_statuses.is_active = 1
		) AS combined_users
	`, primaryOrgUnitField)

	acknowledgedUsersQuery := fmt.Sprintf(`
		SELECT COUNT(DISTINCT user_id) FROM (
			SELECT document_acknowledge_users.user_id 
			FROM document_acknowledge_users
			JOIN users ON users.id = document_acknowledge_users.user_id
			JOIN user_statuses ON user_statuses.id = users.user_status_id
			WHERE document_acknowledge_users.document_id = ? 
			AND user_statuses.is_active = 1
			AND EXISTS (
				SELECT 1 
				FROM document_shared_users 
				WHERE document_shared_users.user_id = document_acknowledge_users.user_id 
				AND document_shared_users.document_id = document_acknowledge_users.document_id
				UNION
				SELECT 1 
				FROM document_shared_units
				JOIN organization_units ON organization_units.id = document_shared_units.organization_unit_id
				JOIN users AS unit_users ON unit_users.%s = organization_units.id
				JOIN user_statuses AS unit_status ON unit_status.id = unit_users.user_status_id
				WHERE document_shared_units.document_id = document_acknowledge_users.document_id
				AND unit_users.id = document_acknowledge_users.user_id AND unit_status.is_active = 1
			)
		) AS acknowledged_users
	`, primaryOrgUnitField)

	err := repo.db.Raw(totalUsersQuery, documentID, documentID).Scan(&totalUsers).Error
	if err != nil {
		return nil, fmt.Errorf("error fetching total users: %v", err)
	}

	err = repo.db.Raw(acknowledgedUsersQuery, documentID).Scan(&acknowledgedUsers).Error
	if err != nil {
		return nil, fmt.Errorf("error fetching acknowledged users: %v", err)
	}

	notAcknowledgedUsers := totalUsers - acknowledgedUsers

	summary := &dtos.DocumentAcknowledgmentSummaryDTO{
		TotalUsers:             totalUsers,
		AcknowledgedUsers:      acknowledgedUsers,
		NotAcknowledgedUsers:   notAcknowledgedUsers,
		AcknowledgedPercent:    float64(acknowledgedUsers) / float64(totalUsers) * 100,
		NotAcknowledgedPercent: float64(notAcknowledgedUsers) / float64(totalUsers) * 100,
	}

	return summary, nil
}

func (repo *GRCReportRepository) GetPublicDocumentAcknowledgmentSummary(documentID uint) (*dtos.DocumentAcknowledgmentSummaryDTO, error) {
	var totalUsers int64
	var acknowledgedUsers int64

	// Query to count all active users
	totalUsersQuery := `
		SELECT COUNT(*)
		FROM users
		JOIN user_statuses ON user_statuses.id = users.user_status_id
		WHERE user_statuses.is_active = 1
	`

	// Query to count acknowledged users for the public document
	acknowledgedUsersQuery := `
		SELECT COUNT(*)
		FROM document_acknowledge_users
		JOIN users ON users.id = document_acknowledge_users.user_id
		JOIN user_statuses ON user_statuses.id = users.user_status_id
		WHERE document_acknowledge_users.document_id = ? AND user_statuses.is_active = 1
	`

	// Fetch total users
	err := repo.db.Raw(totalUsersQuery).Scan(&totalUsers).Error
	if err != nil {
		return nil, fmt.Errorf("error fetching total users for public document: %v", err)
	}

	// Fetch acknowledged users
	err = repo.db.Raw(acknowledgedUsersQuery, documentID).Scan(&acknowledgedUsers).Error
	if err != nil {
		return nil, fmt.Errorf("error fetching acknowledged users for public document: %v", err)
	}

	// Calculate not acknowledged users
	notAcknowledgedUsers := totalUsers - acknowledgedUsers

	// Prepare and return the summary
	summary := &dtos.DocumentAcknowledgmentSummaryDTO{
		TotalUsers:             totalUsers,
		AcknowledgedUsers:      acknowledgedUsers,
		NotAcknowledgedUsers:   notAcknowledgedUsers,
		AcknowledgedPercent:    float64(acknowledgedUsers) / float64(totalUsers) * 100,
		NotAcknowledgedPercent: float64(notAcknowledgedUsers) / float64(totalUsers) * 100,
	}

	return summary, nil
}

func (repo *GRCReportRepository) CreateGRCReportGovernanceDocument(document *models.GRCReportGovernanceDocument, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Create(&document).Error
}

func (repo *GRCReportRepository) GetTasksByDateRange(startDate, endDate time.Time, tx ...utility.DatabaseTransaction) ([]*models.Task, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var tasks []*models.Task

	err := db.
		Preload("ResponsibleUser").
		Preload("TaskStatus").
		Preload("SeverityLevel").
		Joins("LEFT JOIN users ON users.id = tasks.responsible_user_id").
		Where("tasks.created_at BETWEEN ? AND ?", startDate, endDate).
		Find(&tasks).Error

	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (repo *GRCReportRepository) CreateGRCReportGovernanceTask(task *models.GRCReportGovernanceTask, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Create(task).Error
}

func (repo *GRCReportRepository) CreateGRCReportGovernanceSummary(summary *models.GRCReportGovernanceSummary, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Create(summary).Error
}

func (repo *GRCReportRepository) AllUserCount(tx ...utility.DatabaseTransaction) (int64, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var count int64
	err := db.Model(&models.User{}).Count(&count).Error
	return count, err
}

//risk report section

func (repo *GRCReportRepository) GetRisksByDateRange(startDate, endDate time.Time, tx ...utility.DatabaseTransaction) ([]*models.RiskRegister, error) {
	var risks []*models.RiskRegister

	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	err := db.Preload("RiskCategory").
		Preload("RiskResponseStrategy").
		Preload("CurrentRiskRating").
		Preload("KRIs.KriRegister").
		Where("reported_at BETWEEN ? AND ?", startDate, endDate).
		Find(&risks).Error

	return risks, err
}

func (repo *GRCReportRepository) CreateGRCReportRiskSummary(summary *models.GRCReportRiskSummary, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Create(summary).Error
}

func (repo *GRCReportRepository) GetHighestActiveRiskRatingType(tx ...utility.DatabaseTransaction) (string, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var ratingType string
	err := db.Table("risk_matrix_risk_ratings").
		Where("status = ?", "active").
		Order("CASE type WHEN 'very_high' THEN 1 WHEN 'high' THEN 2 WHEN 'medium' THEN 3 WHEN 'low' THEN 4 WHEN 'very_low' THEN 5 END").
		Limit(1).
		Pluck("type", &ratingType).Error
	if err != nil || ratingType == "" {
		return "", fmt.Errorf("failed to fetch highest risk rating type: %w", err)
	}
	return ratingType, nil
}

func (repo *GRCReportRepository) GetTopFiveRisksByRatingValue(startDate time.Time, endDate time.Time, tx ...utility.DatabaseTransaction) ([]dtos.TopFiveRiskDTO, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var risks []dtos.TopFiveRiskDTO

	err := db.Table("risk_register").
		Select(`
			risk_register.id,
			risk_register.reference_number,
			risk_register.title,
			risk_register.description,
			risk_categories.title_en AS risk_category_title_en,
			risk_categories.title_ar AS risk_category_title_ar,
			(CONCAT(responsible_user.first_name_en, ' ', responsible_user.last_name_en)) AS responsible_user_name_en,
			(CONCAT(responsible_user.first_name_ar, ' ', responsible_user.last_name_ar)) AS responsible_user_name_ar,
			organization_units.title_en AS responsible_unit_title_en,
			organization_units.title_ar AS responsible_unit_title_ar,
			risk_statuses.type AS risk_status_type,
			inherent_risk_rating.title_en AS inherent_risk_rating_title_en,
			inherent_risk_rating.title_ar AS inherent_risk_rating_title_ar,
			inherent_risk_rating.color_code AS inherent_risk_rating_color_code,
			inherent_rating_value.value AS inherent_risk_rating_score,
			inherent_likelihood.title_en AS inherent_likelihood_title_en,
			inherent_likelihood.title_ar AS inherent_likelihood_title_ar,
			inherent_impact.title_en AS inherent_impact_title_en,
			inherent_impact.title_ar AS inherent_impact_title_ar,
			residual_risk_rating.title_en AS residual_risk_rating_title_en,
			residual_risk_rating.title_ar AS residual_risk_rating_title_ar,
			residual_risk_rating.color_code AS residual_risk_rating_color_code,
			residual_rating_value.value AS residual_risk_rating_score,
			residual_likelihood.title_en AS residual_likelihood_title_en,
			residual_likelihood.title_ar AS residual_likelihood_title_ar,
			residual_impact.title_en AS residual_impact_title_en,
			residual_impact.title_ar AS residual_impact_title_ar,
			risk_response_strategies.title_en AS risk_response_strategy_title_en,
			risk_response_strategies.title_ar AS risk_response_strategy_title_ar
		`).
		Joins("LEFT JOIN risk_categories ON risk_categories.id = risk_register.risk_category_id").
		Joins("LEFT JOIN users AS responsible_user ON responsible_user.id = risk_register.responsible_user_id").
		Joins("LEFT JOIN organization_units ON organization_units.id = risk_register.responsible_unit_id").
		Joins("LEFT JOIN risk_statuses ON risk_statuses.id = risk_register.risk_status_id").
		Joins("LEFT JOIN risk_matrix_risk_ratings AS inherent_risk_rating ON inherent_risk_rating.id = risk_register.inherent_risk_rating_id").
		Joins("LEFT JOIN risk_matrix_risk_rating_values AS inherent_rating_value ON inherent_rating_value.id = risk_register.inherent_risk_rating_value_id").
		Joins("LEFT JOIN risk_matrix_likelihoods AS inherent_likelihood ON inherent_likelihood.id = risk_register.inherent_likelihood_id").
		Joins("LEFT JOIN risk_matrix_impacts AS inherent_impact ON inherent_impact.id = risk_register.inherent_impact_id").
		Joins("LEFT JOIN risk_matrix_risk_ratings AS residual_risk_rating ON residual_risk_rating.id = risk_register.residual_risk_rating_id").
		Joins("LEFT JOIN risk_matrix_risk_rating_values AS residual_rating_value ON residual_rating_value.id = risk_register.residual_risk_rating_value_id").
		Joins("LEFT JOIN risk_matrix_likelihoods AS residual_likelihood ON residual_likelihood.id = risk_register.residual_likelihood_id").
		Joins("LEFT JOIN risk_matrix_impacts AS residual_impact ON residual_impact.id = risk_register.residual_impact_id").
		Joins("LEFT JOIN risk_response_strategies ON risk_response_strategies.id = risk_register.risk_response_strategy_id").
		Where("risk_register.reported_at BETWEEN ? AND ?", startDate, endDate).
		Order("risk_register.current_risk_rating_value_id DESC").
		Limit(5).
		Find(&risks).Error

	if err != nil {
		return nil, err
	}

	return risks, nil
}

func (repo *GRCReportRepository) CreateGRCReportRiskTopFive(topFive *models.GRCReportRiskTopFive, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Create(topFive).Error
}

func (repo *GRCReportRepository) GetControlTitlesByRiskID(riskID uint, tx ...utility.DatabaseTransaction) ([]string, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var titles []string
	err := db.Table("risk_controls").
		Select("controls.title_en").
		Joins("JOIN controls ON controls.id = risk_controls.control_id").
		Where("risk_controls.risk_register_id = ?", riskID).
		Pluck("controls.title_en", &titles).Error

	return titles, err
}

func (repo *GRCReportRepository) GetKRIsByDateRange(startDate, endDate time.Time, tx ...utility.DatabaseTransaction) ([]dtos.GRCReportRiskKRIDTO, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var results []dtos.GRCReportRiskKRIDTO

	err := db.Table("risk_kris").
		Select(`
			risk_register.title AS risk_register_title,
			kri_register.title AS kri_register_title,
			kri_register.description AS kri_register_description,
			kri_register.low_threshold,
			kri_register.medium_threshold,
			kri_register.current_value,
			risk_statuses.title_en AS risk_status_title_en,
			risk_statuses.title_ar AS risk_status_title_ar,
			risk_statuses.color_code AS risk_status_title_color_code
		`).
		Joins("INNER JOIN risk_register ON risk_register.id = risk_kris.risk_register_id").
		Joins("INNER JOIN kri_register ON kri_register.id = risk_kris.kri_register_id").
		Joins("LEFT JOIN risk_statuses ON risk_statuses.id = risk_register.risk_status_id").
		Where("risk_register.reported_at BETWEEN ? AND ?", startDate, endDate).
		Find(&results).Error

	if err != nil {
		return nil, err
	}

	return results, nil
}

func (repo *GRCReportRepository) CreateGRCReportRiskKRI(riskKri *models.GRCReportRiskKRI, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Create(riskKri).Error
}

func (repo *GRCReportRepository) GetRiskControlsByDateRange(startDate, endDate time.Time, tx ...utility.DatabaseTransaction) ([]dtos.RiskControlReportDTO, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var controls []dtos.RiskControlReportDTO

	err := db.Table("risk_controls").
		Select(`
			controls.code AS control_code,
			control_categories.title_en AS control_category_title_en,
			control_categories.title_ar AS control_category_title_ar,
			controls.title_en AS title,
			controls.description,
			responsible_user.first_name_en + ' ' + responsible_user.last_name_en AS responsible_user_name_en,
			responsible_user.first_name_ar + ' ' + responsible_user.last_name_ar AS responsible_user_name_ar,
			owner_units.title_en AS responsible_unit_title_en,
			owner_units.title_ar AS responsible_unit_title_ar,
			implementation_status.title_en AS control_implementation_status_title_en,
			implementation_status.title_ar AS control_implementation_status_title_ar,
			implementation_status.color_code AS control_implementation_status_color_code,
			efficiency_status.title_en AS control_efficiency_status_title_en,
			efficiency_status.title_ar AS control_efficiency_status_title_ar,
			efficiency_status.color_code AS control_efficiency_status_color_code,
			test_frequency.title_en AS test_frequency_title_en,
			test_frequency.title_ar AS test_frequency_title_ar,
			controls.next_test_date,
			test_user.first_name_en + ' ' + test_user.last_name_en AS test_responsible_user_name_en,
			test_user.first_name_ar + ' ' + test_user.last_name_ar AS test_responsible_user_name_ar
		`).
		Joins("LEFT JOIN controls ON controls.id = risk_controls.control_id").
		Joins("LEFT JOIN control_categories ON control_categories.id = controls.control_category_id").
		Joins("LEFT JOIN users AS responsible_user ON responsible_user.id = controls.responsible_user_id").
		Joins("LEFT JOIN organization_units AS owner_units ON owner_units.id = controls.responsible_unit_id").
		Joins("LEFT JOIN control_implementation_statuses AS implementation_status ON implementation_status.id = controls.control_implementation_status_id").
		Joins("LEFT JOIN control_efficiency_statuses AS efficiency_status ON efficiency_status.id = controls.control_efficiency_status_id").
		Joins("LEFT JOIN review_frequencies AS test_frequency ON test_frequency.id = controls.test_frequency_id").
		Joins("LEFT JOIN users AS test_user ON test_user.id = controls.test_responsible_user_id").
		Joins("LEFT JOIN risk_register ON risk_register.id = risk_controls.risk_register_id").
		Where("risk_register.reported_at BETWEEN ? AND ?", startDate, endDate).
		Find(&controls).Error

	if err != nil {
		return nil, err
	}

	return controls, nil
}

func (repo *GRCReportRepository) CreateGRCReportRiskControl(control *models.GRCReportRiskControl, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Create(control).Error
}

func (repo *GRCReportRepository) GetGRCReportRiskTopFive(grcReportID uint) ([]dtos.GRCReportRiskTopFiveRiskResponseDTO, error) {
	var risks []dtos.GRCReportRiskTopFiveRiskResponseDTO
	err := repo.db.Table("grc_report_risk_top_five").
		Where("grc_report_id = ?", grcReportID).
		Find(&risks).Error

	if err != nil {
		return nil, err
	}
	return risks, nil
}

func (repo *GRCReportRepository) GetGRCReportRiskKRIs(grcReportID uint) ([]dtos.GRCReportRiskKriResponseDTO, error) {
	var kris []dtos.GRCReportRiskKriResponseDTO
	err := repo.db.Table("grc_report_risk_kris").
		Where("grc_report_id = ?", grcReportID).
		Find(&kris).Error

	if err != nil {
		return nil, err
	}
	return kris, nil
}

func (repo *GRCReportRepository) GetGRCReportRiskControls(grcReportID uint) ([]dtos.GRCReportRiskControlResponseDTO, error) {
	var controls []dtos.GRCReportRiskControlResponseDTO
	err := repo.db.Table("grc_report_risk_controls").
		Where("grc_report_id = ?", grcReportID).
		Find(&controls).Error

	if err != nil {
		return nil, err
	}
	return controls, nil
}

func (repo *GRCReportRepository) GetGRCReportRiskSummary(grcReportID uint) (*dtos.GRCRiskSummaryResponseDTO, error) {
	var summary dtos.GRCRiskSummaryResponseDTO

	err := repo.db.Table("grc_report_risk_summary").
		Select(`
			risk_identified_count,
			critical_risks,
			treatment_strategies,
			risk_score_summary,
			kris,
			risk_overview_chart,
			risk_heatmap_chart,
			top_risks,
			risk_trend
		`).
		Where("grc_report_id = ?", grcReportID).
		Scan(&summary).Error

	if err != nil {
		return nil, err
	}

	return &summary, nil
}

// compliance and sla section
func (repo *GRCReportRepository) GetAllComplianceAssessments(startDate time.Time, endDate time.Time, tx ...utility.DatabaseTransaction) ([]*models.ComplianceAssessment, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var assessments []*models.ComplianceAssessment
	// First get the latest assessment IDs for each department
	var latestAssessmentIDs []uint
	subQuery := db.Table("compliance_assessments ca1").
		Select("ca1.id").
		Joins("JOIN compliance_assessment_statuses ON ca1.compliance_assessment_status_id = compliance_assessment_statuses.id").
		Where("ca1.created_at > ? AND ca1.created_at < ?", startDate, endDate).
		Where("compliance_assessment_statuses.type = ?", "approved").
		Where("ca1.created_at = (SELECT MAX(ca2.created_at) FROM compliance_assessments ca2 WHERE ca2.organization_unit_id = ca1.organization_unit_id AND ca2.compliance_register_id = ca1.compliance_register_id)").
		Pluck("ca1.id", &latestAssessmentIDs)

	if subQuery.Error != nil {
		return nil, subQuery.Error
	}

	// Then get the full assessment data for these IDs
	result := db.Preload("AssessmentRequirements").
		Preload("AssessmentRequirements.ComplianceRequirement").
		Preload("AssessmentRequirements.ActionPlans").
		Preload("AssessmentRequirements.ComplianceRequirement.ComplianceClause").
		Preload("AssessmentRequirements.ComplianceStatus").
		Preload("ComplianceAssessmentStatus").
		Where("id IN ?", latestAssessmentIDs).
		Find(&assessments)

	if result.Error != nil {
		return nil, result.Error
	}

	return assessments, nil
}

func (repo *GRCReportRepository) CreateGRCReportComplianceSummary(summary *models.GRCReportComplianceSummary, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	result := db.Create(summary)
	if result.Error != nil {
		return 0, result.Error
	}
	return summary.ID, nil
}

func (repo *GRCReportRepository) CreateGRCReportComplianceDepartment(department *models.GRCReportComplianceDepartment, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	result := db.Create(department)
	if result.Error != nil {
		return 0, result.Error
	}
	return uint(department.ID), nil
}

func (repo *GRCReportRepository) GetOrganizationUnitByID(id uint) (*models.OrganizationUnit, error) {
	var orgUnit models.OrganizationUnit
	result := repo.db.First(&orgUnit, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &orgUnit, nil
}

func (repo *GRCReportRepository) GetAllSLAs(tx ...utility.DatabaseTransaction) ([]*models.SLARegister, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var slas []*models.SLARegister
	err := db.Preload("Documents", func(db *gorm.DB) *gorm.DB {
		return db.Order("created_at DESC").Limit(1)
	}).
		Preload("ResponsibleUnit").
		Preload("ResponsibleUser").
		Preload("LifecycleStatus").
		Preload("ComplianceStatus").
		Find(&slas).Error

	if err != nil {
		return nil, err
	}
	return slas, nil
}

func (repo *GRCReportRepository) CreateGRCReportSLA(sla *models.GRCReportSLA, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	result := db.Create(sla)
	if result.Error != nil {
		return 0, result.Error
	}
	return sla.ID, nil
}

func (repo *GRCReportRepository) GetGRCReportComplianceSummary(grcReportID uint) ([]*models.GRCReportComplianceSummary, error) {
	var summary []*models.GRCReportComplianceSummary
	err := repo.db.
		Preload("OrganizationUnit").
		Preload("ApproveUser").
		Preload("ReviewUser").
		Preload("GRCReportComplianceDepartments", func(db *gorm.DB) *gorm.DB {
			return db.Preload("OrganizationUnit").Preload("ComplianceStatus")
		}).
		Where("grc_report_id = ?", grcReportID).
		Find(&summary).Error
	if err != nil {
		return nil, err
	}

	return summary, nil
}

func (repo *GRCReportRepository) GetGRCReportSLA(grcReportID uint) ([]*dtos.GRCReportSLAResponseDTO, error) {
	var slas []*dtos.GRCReportSLAResponseDTO
	err := repo.db.Table("grc_report_sla").
		Select(`
			grc_report_sla.id,
			grc_report_sla.grc_report_id,
			grc_report_sla.title_en,
			grc_report_sla.title_ar,
			parent_ou.title_en as parent_organization_unit_title_en,
			parent_ou.title_ar as parent_organization_unit_title_ar,
			ou.title_en as organization_unit_title_en,
			ou.title_ar as organization_unit_title_ar,
			compliance_status.title_en as compliance_status_title_en,
			compliance_status.title_ar as compliance_status_title_ar,
			compliance_status.color_code as compliance_status_color_code,
			lifecycle_status.title_en as life_cycle_status_title_en,
			lifecycle_status.title_ar as life_cycle_status_title_ar,
			lifecycle_status.color_code as life_cycle_status_color_code,
			grc_report_sla.expiry_date as expery_date,
			grc_report_sla.issue_date as issue_date
		`).
		Joins("LEFT JOIN organization_units parent_ou ON parent_ou.id = grc_report_sla.parent_organization_unit_id").
		Joins("LEFT JOIN organization_units ou ON ou.id = grc_report_sla.organization_unit_id").
		Joins("LEFT JOIN compliance_statuses compliance_status ON compliance_status.id = grc_report_sla.compliance_status_id").
		Joins("LEFT JOIN lifecycle_statuses lifecycle_status ON lifecycle_status.id = grc_report_sla.lifecycle_status_id").
		Where("grc_report_sla.grc_report_id = ?", grcReportID).
		Find(&slas).Error
	if err != nil {
		return nil, err
	}
	return slas, nil
}

func (repo *GRCReportRepository) GetRiskTrend(startDate time.Time, endDate time.Time, tx ...utility.DatabaseTransaction) ([]dtos.RiskTrendDTO, error) {

	var closedStatusIDs []uint
	if err := repo.db.
		Table("risk_statuses").
		Where("type = ?", "closed").
		Pluck("id", &closedStatusIDs).Error; err != nil {
		return nil, err
	}

	var results []dtos.MonthlyRiskTrendData
	query := `
        SELECT 
            MONTH(created_at) AS month,
            SUM(CASE WHEN created_at BETWEEN ? AND ? THEN 1 ELSE 0 END) AS new_risks,
            SUM(CASE 
                WHEN created_at BETWEEN ? AND ? AND risk_status_id IN (?)
                THEN 1 ELSE 0 END
            ) AS closure_rate
        FROM risk_register
        WHERE created_at BETWEEN ? AND ?
        GROUP BY MONTH(created_at)
        ORDER BY MONTH(created_at)
    `

	if err := repo.db.
		Raw(query,
			startDate, endDate, // for new_risks
			startDate, endDate, // for closure_rate
			closedStatusIDs,    // for IN clause
			startDate, endDate, // for WHERE clause
		).
		Scan(&results).Error; err != nil {
		return nil, err
	}

	// Convert numeric month to label and fill missing months
	monthLabels := [12]string{"Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"}
	dataMap := make(map[int]dtos.MonthlyRiskTrendData)
	for _, res := range results {
		dataMap[res.Month] = res
	}

	var final []dtos.RiskTrendDTO
	for i := 1; i <= 12; i++ {
		m := dataMap[i]
		final = append(final, dtos.RiskTrendDTO{
			Month:       monthLabels[i-1],
			NewRisks:    m.NewRisks,
			ClosureRate: m.ClosureRate,
		})
	}

	return final, nil

}

func (repo *GRCReportRepository) GetGRCReportGovernanceSummaryByID(grcReportID uint, tx ...utility.DatabaseTransaction) (*models.GRCReportGovernanceSummary, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var grcReportGovernanceSummary models.GRCReportGovernanceSummary
	err := db.
		Where("grc_report_id = ?", grcReportID).First(&grcReportGovernanceSummary).Error
	if err != nil {
		return nil, err
	}
	return &grcReportGovernanceSummary, nil
}

func (repo *GRCReportRepository) UpdateGRCReportGovernanceSummary(section *models.GRCReportGovernanceSummary, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Save(section).Error
}
