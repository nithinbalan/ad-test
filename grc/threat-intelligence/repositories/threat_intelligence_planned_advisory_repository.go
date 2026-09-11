package repositories

import (
	"fmt"
	"strconv"
	"time"

	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/models"
	tiDTOs "jethur/grc/threat-intelligence/dtos"
	"jethur/utility"

	"gorm.io/gorm"
)

type ThreatIntelligencePlannedAdvisoryRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(params *tiDTOs.PlannedAdvisoryListQueryParamsDTO) (*tiDTOs.PlannedAdvisoryListResponseDTO, error)
	GetByID(id uint, tx ...utility.DatabaseTransaction) (*models.ThreatIntelligencePlannedAdvisory, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.ThreatIntelligencePlannedAdvisory, error)
	Create(plannedAdvisory *models.ThreatIntelligencePlannedAdvisory, tx ...utility.DatabaseTransaction) (uint, error)
	Update(plannedAdvisory *models.ThreatIntelligencePlannedAdvisory, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error
	GenerateReferenceNumber() (string, error)
	GetDocumentsByPlannedAdvisoryID(plannedAdvisoryID uint) ([]*tiDTOs.PlannedAdvisoryDocumentDTO, error)
	CreatePlannedAdvisoryDocuments(documents []*models.ThreatIntelligencePlannedAdvisoryDocument, tx ...utility.DatabaseTransaction) error
	DeletePlannedAdvisoryDocuments(docIDs []*uint, tx ...utility.DatabaseTransaction) error
	GetPlannedAdvisoryDocumentByID(plannedAdvisoryID uint, documentID uint, tx ...utility.DatabaseTransaction) (*models.ThreatIntelligencePlannedAdvisoryDocument, error)
}

type ThreatIntelligencePlannedAdvisoryRepository struct {
	db *gorm.DB
}

func NewThreatIntelligencePlannedAdvisoryRepository(db *gorm.DB) *ThreatIntelligencePlannedAdvisoryRepository {
	return &ThreatIntelligencePlannedAdvisoryRepository{
		db: db,
	}
}

func (repo *ThreatIntelligencePlannedAdvisoryRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *ThreatIntelligencePlannedAdvisoryRepository) GetAll(params *tiDTOs.PlannedAdvisoryListQueryParamsDTO) (*tiDTOs.PlannedAdvisoryListResponseDTO, error) {
	var plannedAdvisories []*tiDTOs.PlannedAdvisoryListDTO
	var totalCount int64

	query := repo.db.Table("threat_intelligence_planned_advisories AS pa").
		Select(`
			pa.id,
			pa.reference_number,
			pa.case_number,
			pa.title_en,
			pa.title_ar,
			CONVERT(VARCHAR, pa.receipt_date, 120) AS receipt_date,
			pa.case_link,
			pa.threat_type,
			threat_intelligence_severities.score AS threat_severity_score,
			threat_intelligence_likelihoods.score AS likelihood_score,
			pa.risk_level_id,
			pa.risk_level_value_id,
			threat_intelligence_risk_levels.title_en AS risk_level_title_en,
			threat_intelligence_risk_levels.title_ar AS risk_level_title_ar,
			threat_intelligence_risk_levels.color_code AS risk_level_color_code,
			threat_intelligence_risk_levels.type AS risk_level_type,
			pa.processing_time,
			pa.description,
			pa.affected_assets,
			pa.tags,
			pa.iocs,
			pa.advisory_source_id,
			planned_advisory_sources.title_en AS advisory_source_title_en,
			planned_advisory_sources.title_ar AS advisory_source_title_ar,
			pa.case_classification_id,
			planned_advisory_classifications.title_en AS case_classification_title_en,
			planned_advisory_classifications.title_ar AS case_classification_title_ar,
			pa.advisory_type_id,
			planned_advisory_types.title_en AS advisory_type_title_en,
			planned_advisory_types.title_ar AS advisory_type_title_ar,
			pa.importance_id,
			planned_advisory_importances.title_en AS importance_title_en,
			planned_advisory_importances.title_ar AS importance_title_ar,
			pa.security_classification_id,
			planned_advisory_security_classifications.title_en AS security_classification_title_en,
			planned_advisory_security_classifications.title_ar AS security_classification_title_ar,
			pa.advisory_recipient_id,
			advisory_recipient.first_name_en AS advisory_recipient_first_name_en,
			advisory_recipient.first_name_ar AS advisory_recipient_first_name_ar,
			advisory_recipient.last_name_en AS advisory_recipient_last_name_en,
			advisory_recipient.last_name_ar AS advisory_recipient_last_name_ar,
			advisory_recipient.email AS advisory_recipient_email,
			advisory_recipient.initial AS advisory_recipient_initial,
			advisory_recipient.color AS advisory_recipient_color,
			advisory_recipient.profile_image_path AS advisory_recipient_profile_image_path,
			pa.organization_unit_id AS organization_unit_id,
			responsible_unit.title_en AS organization_unit_title_en,
			responsible_unit.title_ar AS organization_unit_title_ar,
			pa.executed_by_id,
			pa.advisory_status_id,
			planned_advisory_statuses.title_en AS advisory_status_title_en,
			planned_advisory_statuses.title_ar AS advisory_status_title_ar,
			pa.internal_status_id,
			planned_advisory_internal_statuses.title_en AS internal_status_title_en,
			planned_advisory_internal_statuses.title_ar AS internal_status_title_ar,
			CONVERT(VARCHAR, pa.notify_date, 120) AS notify_date,
			CONVERT(VARCHAR, pa.deadline_date, 120) AS deadline_date,
			CONVERT(VARCHAR, pa.resolve_date, 120) AS resolve_date,
			CONVERT(VARCHAR, pa.reply_date, 120) AS reply_date,
			pa.saed_ticket_no,
			pa.attachment,
			pa.report_link,
			pa.notes,
			pa.created_at,
			pa.updated_at,
			created_user.id AS created_user_id,
			created_user.first_name_en AS created_user_first_name_en,
			created_user.first_name_ar AS created_user_first_name_ar,
			created_user.last_name_en AS created_user_last_name_en,
			created_user.last_name_ar AS created_user_last_name_ar,
			created_user.initial AS created_user_initial,
			created_user.color AS created_user_color,
			created_user.email AS created_user_email,
			created_user.profile_image_path AS created_user_profile_image_path,
			updated_user.id AS updated_user_id,
			updated_user.first_name_en AS updated_user_first_name_en,
			updated_user.first_name_ar AS updated_user_first_name_ar,
			updated_user.last_name_en AS updated_user_last_name_en,
			updated_user.last_name_ar AS updated_user_last_name_ar,
			updated_user.initial AS updated_user_initial,
			updated_user.color AS updated_user_color,
			updated_user.email AS updated_user_email,
			updated_user.profile_image_path AS updated_user_profile_image_path
		`).
		Joins(`
			LEFT JOIN planned_advisory_sources ON planned_advisory_sources.id = pa.advisory_source_id
			LEFT JOIN planned_advisory_classifications ON planned_advisory_classifications.id = pa.case_classification_id
			LEFT JOIN planned_advisory_types ON planned_advisory_types.id = pa.advisory_type_id
			LEFT JOIN planned_advisory_importances ON planned_advisory_importances.id = pa.importance_id
			LEFT JOIN planned_advisory_security_classifications ON planned_advisory_security_classifications.id = pa.security_classification_id
			LEFT JOIN planned_advisory_statuses ON planned_advisory_statuses.id = pa.advisory_status_id
			LEFT JOIN planned_advisory_internal_statuses ON planned_advisory_internal_statuses.id = pa.internal_status_id
			LEFT JOIN users AS advisory_recipient ON advisory_recipient.id = pa.advisory_recipient_id
			LEFT JOIN organization_units AS responsible_unit ON responsible_unit.id = pa.organization_unit_id
			LEFT JOIN users AS created_user ON created_user.id = pa.created_by_id
			LEFT JOIN users AS updated_user ON updated_user.id = pa.updated_by_id
			LEFT JOIN threat_intelligence_severities ON threat_intelligence_severities.id = pa.threat_severity_id
			LEFT JOIN threat_intelligence_likelihoods ON threat_intelligence_likelihoods.id = pa.likelihood_id
			LEFT JOIN threat_intelligence_risk_levels ON threat_intelligence_risk_levels.id = pa.risk_level_id
			LEFT JOIN threat_intelligence_risk_level_values ON threat_intelligence_risk_level_values.id = pa.risk_level_value_id
		`)

	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where(columnName+" IN (?)", values)
		}
	}

	if params.ReceiptDateFrom != "" {
		query = query.Where("pa.receipt_date >= ?", params.ReceiptDateFrom)
	}

	if params.ReceiptDateTo != "" {
		query = query.Where("pa.receipt_date <= ?", params.ReceiptDateTo)
	}

	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("pa.receipt_date BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	if params.Status != "" {
		query = query.Where("planned_advisory_statuses.title_en = ?", params.Status)
	}

	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
			(pa.reference_number LIKE ? OR
			pa.case_number LIKE ? OR
			pa.title_en LIKE ? OR
			pa.title_ar LIKE ? OR
			pa.case_link LIKE ? OR
			pa.threat_type LIKE ? OR
			pa.description LIKE ? OR
			pa.affected_assets LIKE ? OR
			pa.tags LIKE ? OR
			pa.iocs LIKE ? OR
			pa.saed_ticket_no LIKE ? OR
			pa.notes LIKE ? OR
			planned_advisory_sources.title_en LIKE ? OR
			planned_advisory_sources.title_ar LIKE ? OR
			planned_advisory_classifications.title_en LIKE ? OR
			planned_advisory_classifications.title_ar LIKE ? OR
			planned_advisory_types.title_en LIKE ? OR
			planned_advisory_types.title_ar LIKE ? OR
			planned_advisory_importances.title_en LIKE ? OR
			planned_advisory_importances.title_ar LIKE ? OR
			planned_advisory_security_classifications.title_en LIKE ? OR
			planned_advisory_security_classifications.title_ar LIKE ? OR
			planned_advisory_statuses.title_en LIKE ? OR
			planned_advisory_statuses.title_ar LIKE ? OR
			advisory_recipient.first_name_en LIKE ? OR
			advisory_recipient.first_name_ar LIKE ? OR
			advisory_recipient.last_name_en LIKE ? OR
			advisory_recipient.last_name_ar LIKE ? OR
			responsible_unit.title_en LIKE ? OR
			responsible_unit.title_ar LIKE ? OR
			pa.receipt_date LIKE ? OR
			created_user.first_name_en LIKE ? OR
			created_user.first_name_ar LIKE ? OR
			created_user.last_name_en LIKE ? OR
			created_user.last_name_ar LIKE ? OR
			updated_user.first_name_en LIKE ? OR
			updated_user.first_name_ar LIKE ? OR
			updated_user.last_name_en LIKE ? OR
			updated_user.last_name_ar LIKE ?)`,
			keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword,
			keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword,
			keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword,
			keyword, keyword, keyword, keyword, keyword)
	}

	query = query.Order(params.SortBy + " " + params.Sort)

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	if !params.Export {
		offset := (params.Page - 1) * params.PageSize
		query = query.Offset(offset).Limit(params.PageSize)
	}

	if err := query.Find(&plannedAdvisories).Error; err != nil {
		return nil, err
	}

	if params.Export {
		return &tiDTOs.PlannedAdvisoryListResponseDTO{
			Items: plannedAdvisories,
		}, nil
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	return &tiDTOs.PlannedAdvisoryListResponseDTO{
		Items: plannedAdvisories,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}, nil
}

func (repo *ThreatIntelligencePlannedAdvisoryRepository) GetByID(id uint, tx ...utility.DatabaseTransaction) (*models.ThreatIntelligencePlannedAdvisory, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var plannedAdvisory *models.ThreatIntelligencePlannedAdvisory
	result := db.First(&plannedAdvisory, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return plannedAdvisory, nil
}

func (repo *ThreatIntelligencePlannedAdvisoryRepository) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.ThreatIntelligencePlannedAdvisory, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var plannedAdvisory *models.ThreatIntelligencePlannedAdvisory
	result := db.
		Preload("AdvisorySource").
		Preload("CaseClassification").
		Preload("AdvisoryType").
		Preload("Importance").
		Preload("SecurityClassification").
		Preload("AdvisoryRecipient").
		Preload("OrganizationUnit").
		Preload("ExecutedBy").
		Preload("AdvisoryStatus").
		Preload("InternalStatus").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		Preload("ThreatSeverity").
		Preload("Likelihood").
		Preload("RiskLevel").
		Preload("RiskLevelValue").
		Preload("Documents").
		First(&plannedAdvisory, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return plannedAdvisory, nil
}

func (repo *ThreatIntelligencePlannedAdvisoryRepository) Create(plannedAdvisory *models.ThreatIntelligencePlannedAdvisory, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(plannedAdvisory).Error; err != nil {
		return 0, err
	}
	return plannedAdvisory.ID, nil
}

func (repo *ThreatIntelligencePlannedAdvisoryRepository) Update(plannedAdvisory *models.ThreatIntelligencePlannedAdvisory, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(plannedAdvisory).Error; err != nil {
		return err
	}
	return nil
}

func (repo *ThreatIntelligencePlannedAdvisoryRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Delete(&models.ThreatIntelligencePlannedAdvisory{}, id).Error; err != nil {
		return err
	}
	return nil
}

func (repo *ThreatIntelligencePlannedAdvisoryRepository) GenerateReferenceNumber() (string, error) {
	var lastReferenceNumber string
	err := repo.db.Model(&models.ThreatIntelligencePlannedAdvisory{}).
		Select("reference_number").
		Order("reference_number desc").
		Where("reference_number LIKE ?", "PA-%").
		Limit(1).
		Pluck("reference_number", &lastReferenceNumber).
		Error

	if err != nil {
		return "", err
	}

	lastDigits := "0000"
	if lastReferenceNumber != "" && len(lastReferenceNumber) >= 4 {
		lastDigits = lastReferenceNumber[len(lastReferenceNumber)-4:]
	}

	currentYear := time.Now().Year()
	lastReferenceID, err := strconv.ParseUint(lastDigits, 10, 64)
	if err != nil {
		return "", err
	}
	newReferenceNumber := fmt.Sprintf("%s-%d-%04d", "PA", currentYear, lastReferenceID+1)

	return newReferenceNumber, nil
}

func (repo *ThreatIntelligencePlannedAdvisoryRepository) GetDocumentsByPlannedAdvisoryID(plannedAdvisoryID uint) ([]*tiDTOs.PlannedAdvisoryDocumentDTO, error) {
	var documents []*tiDTOs.PlannedAdvisoryDocumentDTO
	err := repo.db.Table("threat_intelligence_planned_advisory_documents").
		Select("id, threat_intelligence_planned_advisory_id, url, name, size, original_name, type, token, created_at").
		Where("threat_intelligence_planned_advisory_id = ?", plannedAdvisoryID).
		Find(&documents).Error
	if err != nil {
		return nil, err
	}
	return documents, nil
}

func (repo *ThreatIntelligencePlannedAdvisoryRepository) CreatePlannedAdvisoryDocuments(documents []*models.ThreatIntelligencePlannedAdvisoryDocument, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if len(documents) == 0 {
		return nil
	}
	return db.Create(&documents).Error
}

func (repo *ThreatIntelligencePlannedAdvisoryRepository) DeletePlannedAdvisoryDocuments(docIDs []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Where("id IN (?)", docIDs).Delete(&models.ThreatIntelligencePlannedAdvisoryDocument{}).Error
}

func (repo *ThreatIntelligencePlannedAdvisoryRepository) GetPlannedAdvisoryDocumentByID(plannedAdvisoryID uint, documentID uint, tx ...utility.DatabaseTransaction) (*models.ThreatIntelligencePlannedAdvisoryDocument, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var document models.ThreatIntelligencePlannedAdvisoryDocument
	err := db.Where("id = ? AND threat_intelligence_planned_advisory_id = ?", documentID, plannedAdvisoryID).First(&document).Error
	if err != nil {
		return nil, err
	}
	return &document, nil
}
