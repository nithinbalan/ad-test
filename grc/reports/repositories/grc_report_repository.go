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

	DeleteGRCReport(grcReportID uint, tx ...utility.DatabaseTransaction) error
	DeleteGRCReportSection(sectionID uint, tx ...utility.DatabaseTransaction) error
	DeleteGRCReportSections(grcReportID uint, tx ...utility.DatabaseTransaction) error
	DeleteGRCReportGovernanceDocuments(grcReportID uint, tx ...utility.DatabaseTransaction) error
	DeleteGRCReportGovernanceTasks(grcReportID uint, tx ...utility.DatabaseTransaction) error

	ReorderGRCReportSections(grcReportID uint, deletedOrderNum int, tx ...utility.DatabaseTransaction) error
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
			document_status_title_en, document_status_title_ar, document_status_color_code,
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
		Where("dt.type = ?", "policy").
		Where("documents.created_at BETWEEN ? AND ?", startDate, endDate).
		Preload("ResponsibleUser").
		Preload("ResponsibleUnit").
		Preload("DocumentAccessType").
		Preload("DocumentStatus").
		Preload("DocumentType").
		Preload("DocumentCategory").
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
