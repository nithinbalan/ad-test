package repositories

import (
	"fmt"
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/models"
	"jethur/grc/organization/dtos"

	"jethur/utility"
	"time"

	"gorm.io/gorm"
)

type GovernanceReportRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction

	GetAll(params *dtos.GovernanceReportListQueryParamsDTO) (*dtos.GovernanceReportListResponseDTO, error)
	GetSectionByID(sectionID uint, tx ...utility.DatabaseTransaction) (*models.GovernanceReportSection, error)

	GetGovernanceReportByID(governanceReportID uint, tx ...utility.DatabaseTransaction) (*dtos.GovernanceReportListDTO, error)
	GetGovernanceReportSections(governanceReportID uint) ([]dtos.GovernanceReportSectionResponseDTO, error)
	GetGovernanceReportDocuments(reportID uint) ([]dtos.GovernanceDocumentResponseDTO, error)
	GetGovernanceReportMeetings(reportID uint) ([]dtos.GovernanceReportMeetingResponseDTO, error)

	Create(governanceReport *models.GovernanceReport, tx ...utility.DatabaseTransaction) (uint, error)
	CreateSection(section *models.GovernanceReportSection, tx ...utility.DatabaseTransaction) error
	GetPolicyDocumentByDateRange(startDate, endDate time.Time, tx ...utility.DatabaseTransaction) ([]*models.Document, error)
	GetFullDetailsByID(id uint, currentUserID *uint, tx ...utility.DatabaseTransaction) (*models.Document, error)
	GetDocumentAcknowledgmentSummary(documentID uint, primaryOrgUnitField string) (*dtos.DocumentAcknowledgmentSummaryDTO, error)
	GetPublicDocumentAcknowledgmentSummary(documentID uint) (*dtos.DocumentAcknowledgmentSummaryDTO, error)
	CreateGovernanceIDocument(document *models.GovernanceReportDocument, tx ...utility.DatabaseTransaction) error

	UpdateSection(section *models.GovernanceReportSection, tx ...utility.DatabaseTransaction) error
	DeleteGovernanceReportSections(governanceReportID uint, tx ...utility.DatabaseTransaction) error
	DeleteGovernanceReportDocuments(governanceReportID uint, tx ...utility.DatabaseTransaction) error
	DeleteGovernanceReportMeetings(governanceReportID uint, tx ...utility.DatabaseTransaction) error

	DeleteGovernanceReport(governanceReportID uint, tx ...utility.DatabaseTransaction) error
	DeleteGovernanceReportSection(sectionID uint, tx ...utility.DatabaseTransaction) error
	ReorderGovernanceReportSections(governanceReportID uint, deletedOrderNum int, tx ...utility.DatabaseTransaction) error

	GetMeetingsByDateRange(startDate, endDate time.Time, tx ...utility.DatabaseTransaction) ([]*models.Meeting, error)
	CreateGovernanceReportMeeting(meeting *models.GovernanceReportMeeting, tx ...utility.DatabaseTransaction) error
}

type GovernanceReportRepository struct {
	db *gorm.DB
}

func NewGovernanceReportRepository(db *gorm.DB) GovernanceReportRepositoryInterface {
	return &GovernanceReportRepository{db: db}
}

func (repo *GovernanceReportRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *GovernanceReportRepository) GetAll(params *dtos.GovernanceReportListQueryParamsDTO) (*dtos.GovernanceReportListResponseDTO, error) {
	var reports []*dtos.GovernanceReportListDTO
	var totalCount int64

	query := repo.db.Table("governance_reports").
		Select(`
			governance_reports.id, 
			governance_reports.title, 
			CONVERT(VARCHAR, governance_reports.start_date, 23) AS start_date,
			CONVERT(VARCHAR, governance_reports.end_date, 23) AS end_date,
			governance_reports.created_at,
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
		Joins("LEFT JOIN users AS created_user ON created_user.id = governance_reports.created_by_id")

	if params.Keyword != "" {
		query = query.Where(`(governance_reports.title LIKE ? )`, "%"+params.Keyword+"%")
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

	response := &dtos.GovernanceReportListResponseDTO{
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

func (repo *GovernanceReportRepository) GetSectionByID(sectionID uint, tx ...utility.DatabaseTransaction) (*models.GovernanceReportSection, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var governanceReportSection models.GovernanceReportSection
	err := db.
		Where("id = ?", sectionID).First(&governanceReportSection).Error
	if err != nil {
		return nil, err
	}
	return &governanceReportSection, nil
}

func (repo *GovernanceReportRepository) GetGovernanceReportByID(governanceReportID uint, tx ...utility.DatabaseTransaction) (*dtos.GovernanceReportListDTO, error) {
	var report dtos.GovernanceReportListDTO

	query := repo.db.Table("governance_reports").
		Select(`
			governance_reports.id, 
			governance_reports.title, 
			CONVERT(VARCHAR, governance_reports.start_date, 23) AS start_date,
			CONVERT(VARCHAR, governance_reports.end_date, 23) AS end_date,
			governance_reports.created_at,
			users.id AS created_by_id,
			users.first_name_en AS created_user_first_name_en,
			users.first_name_ar AS created_user_first_name_ar,
			users.last_name_en AS created_user_last_name_en,
			users.last_name_ar AS created_user_last_name_ar,
			users.email AS created_user_email
		`).
		Joins("LEFT JOIN users ON users.id = governance_reports.created_by_id").
		Where("governance_reports.id = ?", governanceReportID).
		Group("governance_reports.id, governance_reports.title, governance_reports.start_date, governance_reports.end_date, governance_reports.created_at, users.id, users.first_name_en, users.first_name_ar, users.last_name_en, users.last_name_ar, users.email")

	if err := query.First(&report).Error; err != nil {

		return nil, err
	}

	return &report, nil
}

// Fetch sections for the report
func (repo *GovernanceReportRepository) GetGovernanceReportSections(governanceReportID uint) ([]dtos.GovernanceReportSectionResponseDTO, error) {
	var sections []dtos.GovernanceReportSectionResponseDTO
	err := repo.db.Table("governance_report_sections").
		Select("id, type, title_en, title_ar, description_en, description_ar, order_num").
		Where("governance_report_id = ?", governanceReportID).
		Order("order_num ASC").
		Find(&sections).Error

	if err != nil {
		return nil, err
	}
	return sections, nil
}

func (repo *GovernanceReportRepository) GetGovernanceReportDocuments(governanceReportID uint) ([]dtos.GovernanceDocumentResponseDTO, error) {
	var documents []dtos.GovernanceDocumentResponseDTO

	err := repo.db.Table("governance_report_documents").
		Select(`
			id, governance_report_id, reference_number,
			document_title_en, document_title_ar,
			responsible_user_name_en, responsible_user_name_ar,
			responsible_unit_title_en, responsible_unit_title_ar,
			total_employees, acknowledged_count, pending_count,
			created_by_id, created_at
		`).
		Where("governance_report_id = ?", governanceReportID).
		Order("created_at ASC").
		Scan(&documents).Error

	return documents, err
}

func (repo *GovernanceReportRepository) GetGovernanceReportMeetings(governanceReportID uint) ([]dtos.GovernanceReportMeetingResponseDTO, error) {
	var meetings []dtos.GovernanceReportMeetingResponseDTO

	err := repo.db.Table("governance_report_meetings").
		Select(`
			id, governance_report_id, meeting_title, meeting_date,
			participants_count, key_decisions, actions_assigned,
			created_by_id, created_at
		`).
		Where("governance_report_id = ?", governanceReportID).
		Order("meeting_date ASC").
		Scan(&meetings).Error

	return meetings, err
}

func (repo *GovernanceReportRepository) Create(governanceReport *models.GovernanceReport, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(governanceReport).Error; err != nil {
		return 0, err
	}
	return governanceReport.ID, nil
}

func (repo *GovernanceReportRepository) GetPolicyDocumentByDateRange(startDate, endDate time.Time, tx ...utility.DatabaseTransaction) ([]*models.Document, error) {
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
		Find(&documents).Error

	if err != nil {
		return nil, err
	}

	return documents, nil
}

func (repo *GovernanceReportRepository) CreateSection(section *models.GovernanceReportSection, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Create(&section).Error
}

func (repo *GovernanceReportRepository) UpdateSection(section *models.GovernanceReportSection, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Save(section).Error
}

func (repo *GovernanceReportRepository) DeleteGovernanceReportSections(governanceReportID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Where("governance_report_id = ?", governanceReportID).Delete(&models.GovernanceReportSection{}).Error
}

func (repo *GovernanceReportRepository) DeleteGovernanceReportDocuments(governanceReportID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Where("governance_report_id = ?", governanceReportID).Delete(&models.GovernanceReportDocument{}).Error
}

func (repo *GovernanceReportRepository) DeleteGovernanceReportMeetings(governanceReportID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Where("governance_report_id = ?", governanceReportID).Delete(&models.GovernanceReportMeeting{}).Error
}

func (repo *GovernanceReportRepository) DeleteGovernanceReport(governanceReportID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Where("id = ?", governanceReportID).Delete(&models.GovernanceReport{}).Error
}

func (repo *GovernanceReportRepository) DeleteGovernanceReportSection(sectionID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Where("id = ?", sectionID).Delete(&models.GovernanceReportSection{}).Error
}

func (repo *GovernanceReportRepository) ReorderGovernanceReportSections(governanceReportID uint, deletedOrderNum int, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	// Shift order numbers of remaining sections
	err := db.Exec(`
		UPDATE governance_report_sections
		SET order_num = order_num - 1
		WHERE governance_report_id = ? AND order_num > ?`,
		governanceReportID, deletedOrderNum).Error

	return err
}

func (repo *GovernanceReportRepository) GetFullDetailsByID(id uint, currentUserID *uint, tx ...utility.DatabaseTransaction) (*models.Document, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var document *models.Document
	result := db.
		Preload("DocumentAccessType").
		First(&document, id)

	if result.Error != nil {
		return nil, result.Error
	}

	// If currentUserID is provided, check the acknowledgment status
	if currentUserID != nil {
		var count int64
		db.Model(&models.DocumentAcknowledgeUser{}).
			Where("document_id = ? AND user_id = ?", id, *currentUserID).
			Count(&count)
		document.Acknowledged = count > 0
	}

	return document, nil
}

func (repo *GovernanceReportRepository) GetDocumentAcknowledgmentSummary(documentID uint, primaryOrgUnitField string) (*dtos.DocumentAcknowledgmentSummaryDTO, error) {
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

func (repo *GovernanceReportRepository) GetPublicDocumentAcknowledgmentSummary(documentID uint) (*dtos.DocumentAcknowledgmentSummaryDTO, error) {
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

func (repo *GovernanceReportRepository) CreateGovernanceIDocument(document *models.GovernanceReportDocument, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Create(&document).Error
}

func (repo *GovernanceReportRepository) GetMeetingsByDateRange(startDate, endDate time.Time, tx ...utility.DatabaseTransaction) ([]*models.Meeting, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var meetings []*models.Meeting
	err := db.
		Preload("Participants").
		Preload("MeetingMinutes").
		Preload("MeetingActionPlans.ResponsibleUser").
		Where("date BETWEEN ? AND ?", startDate, endDate).
		Find(&meetings).Error

	if err != nil {
		return nil, err
	}

	return meetings, nil
}

func (repo *GovernanceReportRepository) CreateGovernanceReportMeeting(meeting *models.GovernanceReportMeeting, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Create(&meeting).Error
}
