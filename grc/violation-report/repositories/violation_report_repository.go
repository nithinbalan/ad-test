package repositories

import (
	"fmt"
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/models"
	"jethur/grc/violation-report/dtos"
	violationReportDTOs "jethur/grc/violation-report/dtos"
	"jethur/utility"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type ViolationReportRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(params *violationReportDTOs.ViolationReportListQueryParamsDTO, departmentID *uint, currentUserID uint) (*violationReportDTOs.ViolationReportListResponseDTO, error)
	GetByID(id uint) (*models.ViolationReport, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.ViolationReport, error)
	Create(violationReport *models.ViolationReport, tx ...utility.DatabaseTransaction) (uint, error)
	Update(violationReport *models.ViolationReport, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error
	GenerateReferenceNumber() (string, error)

	//Document
	CreateDocuments(docs []*models.ViolationReportDocument, tx ...utility.DatabaseTransaction) error
	GetDocumentByID(id uint) (*models.ViolationReportDocument, error)
	DeleteDocuments(docIds []*uint, tx ...utility.DatabaseTransaction) error

	//comment
	CreateComment(comment *models.ViolationReportComment, tx ...utility.DatabaseTransaction) (uint, error)
	GetAllViolationReportComments(violationReportID uint) ([]*dtos.ViolationReportCommentListDTO, error)
	GetCommentByID(id uint) (*models.ViolationReportComment, error)
	GetFullCommentDetailsByID(id uint) (*models.ViolationReportComment, error)
	UpdateComment(comment *models.ViolationReportComment, tx ...utility.DatabaseTransaction) error
	DeleteComment(id uint, tx ...utility.DatabaseTransaction) error
	CreateCommentDocuments(docs []*models.ViolationReportCommentDocument, tx ...utility.DatabaseTransaction) error
	DeleteCommentDocuments(docIDs []*uint, tx ...utility.DatabaseTransaction) error
	GetCommentDocumentByID(id uint) (*models.ViolationReportCommentDocument, error)
}

type ViolationReportRepository struct {
	db *gorm.DB
}

func NewViolationReportRepository(db *gorm.DB) *ViolationReportRepository {
	return &ViolationReportRepository{
		db: db,
	}
}

// BeginTransaction starts a new transaction.
func (repo *ViolationReportRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *ViolationReportRepository) GetAll(params *violationReportDTOs.ViolationReportListQueryParamsDTO, departmentID *uint, currentUserID uint) (*violationReportDTOs.ViolationReportListResponseDTO, error) {

	var reports []*violationReportDTOs.ViolationReportListDTO
	var totalCount int64

	query := repo.db.Table("violation_reports").
		Select(`
			violation_reports.id,
			violation_reports.reference_number,
			violation_reports.reported_at,
			violation_reports.violation_report_channel_id,
			violation_report_channels.title_en AS violation_report_channel_title_en,
			violation_report_channels.title_ar AS violation_report_channel_title_ar,
			violation_reports.title_en,
			violation_reports.title_ar,
			violation_reports.reported_entity_id,
			violation_reports.id_of_reporter,
			violation_reports.name_of_reporter,
			violation_reports.phone_of_reporter,
			violation_reports.email_of_reporter,
			organization_units.title_en AS reported_entity_title_en,
			organization_units.title_ar AS reported_entity_title_ar,
			violation_reports.description_en,
			violation_reports.description_ar,
			violation_reports.violation_report_category_id,
			violation_report_categories.title_en AS violation_report_category_title_en,
			violation_report_categories.title_ar AS violation_report_category_title_ar,
			violation_reports.violation_report_sub_category_id,
			violation_report_sub_categories.title_en AS violation_report_sub_category_title_en,
			violation_report_sub_categories.title_ar AS violation_report_sub_category_title_ar,
			violation_reports.violation_report_reason_id,
			violation_report_reasons.title_en AS violation_report_reason_title_en,
			violation_report_reasons.title_ar AS violation_report_reason_title_ar,
			violation_reports.evidence_description AS evidence_description,
			violation_reports.severity_level_id,
			severity_levels.title_en AS severity_level_title_en,
			severity_levels.title_ar AS severity_level_title_ar,
			severity_levels.color_code AS severity_level_color_code,
			violation_reports.violation_report_status_id,
			violation_report_statuses.title_en AS violation_report_status_title_en,
			violation_report_statuses.title_ar AS violation_report_status_title_ar,
			violation_report_statuses.color_code AS violation_report_status_color_code,
			violation_reports.created_at,
			violation_reports.updated_at,
			violation_reports.created_by_id,
			created_user.first_name_en AS created_user_first_name_en,
			created_user.first_name_ar AS created_user_first_name_ar,
			created_user.last_name_en AS created_user_last_name_en,
			created_user.last_name_ar AS created_user_last_name_ar,
			created_user.initial AS created_user_initial,
			created_user.color AS created_user_color,
			created_user.email AS created_user_email,
			created_user.profile_image_path AS created_user_profile_image_path,

			violation_reports.updated_by_id,
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
			LEFT JOIN violation_report_channels ON violation_report_channels.id = violation_reports.violation_report_channel_id
			LEFT JOIN violation_report_categories ON violation_report_categories.id = violation_reports.violation_report_category_id
			LEFT JOIN violation_report_sub_categories ON violation_report_sub_categories.id = violation_reports.violation_report_sub_category_id
			LEFT JOIN violation_report_reasons ON violation_report_reasons.id = violation_reports.violation_report_reason_id
			LEFT JOIN violation_report_statuses ON violation_report_statuses.id = violation_reports.violation_report_status_id
			LEFT JOIN severity_levels ON severity_levels.id = violation_reports.severity_level_id
			LEFT JOIN organization_units ON organization_units.id = violation_reports.reported_entity_id
			LEFT JOIN users AS created_user ON created_user.id = violation_reports.created_by_id
			LEFT JOIN users AS updated_user ON updated_user.id = violation_reports.updated_by_id
		`)

		// Add department filter if departmentID is provided (non-admin users)
	if departmentID != nil {
		query = query.Where(
			repo.db.Where("violation_reports.created_by_id = ?", currentUserID).
				Or("violation_reports.reported_entity_id = ?", *departmentID),
		)
	}

	for column, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where("violation_reports."+column+" IN (?)", values)
		}
	}

	if params.ReportedAt != "" {
		query = query.Where("violation_reports.reported_at >= ?", params.ReportedAt)
	}

	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
			violation_reports.reference_number LIKE ? OR
			violation_reports.title_en LIKE ? OR
			violation_reports.title_ar LIKE ? OR
			violation_report_channels.title_en LIKE ? OR
			violation_report_channels.title_ar LIKE ? OR
			violation_reports.description_en LIKE ? OR
			violation_reports.description_ar LIKE ? OR
			violation_report_categories.title_en LIKE ? OR
			violation_report_categories.title_ar LIKE ? OR
			violation_report_sub_categories.title_en LIKE ? OR
			violation_report_sub_categories.title_ar LIKE ? OR
			violation_report_reasons.title_en LIKE ? OR
			violation_report_reasons.title_ar LIKE ? OR
			violation_report_statuses.title_en LIKE ? OR
			violation_report_statuses.title_ar LIKE ? OR
			severity_levels.title_en LIKE ? OR
			severity_levels.title_ar LIKE ? OR
			organization_units.title_en LIKE ? OR
			organization_units.title_ar LIKE ? OR
			created_user.first_name_en LIKE ? OR
			created_user.last_name_en LIKE ? OR
			created_user.email LIKE ?
		`,
			keyword, keyword, keyword, keyword, keyword,
			keyword, keyword, keyword, keyword,
			keyword, keyword, keyword, keyword,
			keyword, keyword, keyword, keyword,
			keyword, keyword, keyword, keyword, keyword,
		)
	}

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}
	query = query.Order(params.SortBy + " " + params.Sort)

	if !params.Export {
		offset := (params.Page - 1) * params.PageSize
		query = query.Offset(offset).Limit(params.PageSize)
	}

	if err := query.Find(&reports).Error; err != nil {
		return nil, err
	}

	totalPages := 1
	if !params.Export {
		totalPages = int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))
	}

	return &violationReportDTOs.ViolationReportListResponseDTO{
		Items: reports,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}, nil
}

func (repo *ViolationReportRepository) GetByID(id uint) (*models.ViolationReport, error) {
	var violationReport *models.ViolationReport

	result := repo.db.First(&violationReport, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return violationReport, nil

}

func (repo *ViolationReportRepository) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.ViolationReport, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var violationReport *models.ViolationReport

	result := db.
		Preload("ViolationReportChannel").
		Preload("ViolationReportCategory").
		Preload("ViolationReportCategory.Department").
		Preload("ViolationReportSubCategory").
		Preload("ViolationReportReason").
		Preload("ViolationReportStatus").
		Preload("SeverityLevel").
		Preload("Documents").
		Preload("ReportedEntity.ParentUnit.ParentUnit.ParentUnit.ParentUnit.ParentUnit").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		First(&violationReport, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return violationReport, nil
}

func (repo *ViolationReportRepository) Create(violationReport *models.ViolationReport, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(violationReport).Error; err != nil {
		return 0, err
	}

	return violationReport.ID, nil
}

func (repo *ViolationReportRepository) Update(violationReport *models.ViolationReport, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(violationReport).Error; err != nil {
		return err
	}

	return nil
}

func (repo *ViolationReportRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	// Delete the violation report
	if err := db.Delete(&models.ViolationReport{}, id).Error; err != nil {
		return err
	}

	return nil
}

func (repo *ViolationReportRepository) GenerateReferenceNumber() (string, error) {
	var lastReferenceNumber string
	err := repo.db.Model(&models.ViolationReport{}).
		Select("reference_number").
		Order("reference_number desc").
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
	newReferenceNumber := fmt.Sprintf("%s-%d-%04d", "VIOL", currentYear, lastReferenceID+1)

	return newReferenceNumber, nil

}

// Document func
func (repo *ViolationReportRepository) GetDocumentByID(id uint) (*models.ViolationReportDocument, error) {
	var document models.ViolationReportDocument
	if err := repo.db.Where("id = ?", id).First(&document).Error; err != nil {
		return nil, err
	}
	return &document, nil
}

func (repo *ViolationReportRepository) CreateDocuments(docs []*models.ViolationReportDocument, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	// Save all documents in a single transaction
	if err := db.Create(&docs).Error; err != nil {
		return err
	}
	return nil
}

func (repo *ViolationReportRepository) DeleteDocuments(docIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	// Delete documents by IDs
	if err := db.Where("id IN (?)", docIds).Delete(&models.ViolationReportDocument{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *ViolationReportRepository) CreateComment(comment *models.ViolationReportComment, tx ...utility.DatabaseTransaction) (uint, error) {

	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Create(comment).Error; err != nil {
		return 0, err
	}

	return comment.ID, nil
}

func (repo *ViolationReportRepository) GetAllViolationReportComments(violationReportID uint) ([]*dtos.ViolationReportCommentListDTO, error) {

	var comments []*dtos.ViolationReportCommentListDTO

	query := repo.db.Table("violation_report_comments").
		Select(`
			violation_report_comments.id,
			violation_report_comments.violation_report_id,
			violation_report_comments.comment,

			violation_report_statuses.id AS violation_report_status_id,
			violation_report_statuses.title_en AS status_title_en,
			violation_report_statuses.title_ar AS status_title_ar,
			violation_report_statuses.type AS status_type,
			violation_report_statuses.status AS status,

			violation_report_comments.created_at,
			users.id AS created_by_id,
			users.first_name_en AS created_by_first_name_en,
			users.last_name_en AS created_by_last_name_en,
			users.first_name_ar AS created_by_first_name_ar,
			users.last_name_ar AS created_by_last_name_ar,
			users.initial AS created_by_initial
		`).
		Joins(`
			LEFT JOIN violation_report_statuses 
			ON violation_report_statuses.id = violation_report_comments.violation_report_status_id
		`).
		Joins(`
			LEFT JOIN users 
			ON users.id = violation_report_comments.created_by_id
		`).
		Where("violation_report_comments.violation_report_id = ?", violationReportID).
		Order("violation_report_comments.created_at ASC")

	if err := query.Find(&comments).Error; err != nil {
		return nil, err
	}

	for _, comment := range comments {
		var documents []*dtos.ViolationReportCommentDocumentDTO
		repo.db.Model(&models.ViolationReportCommentDocument{}).
			Where("violation_report_comment_id = ?", comment.ID).
			Find(&documents)
		comment.Documents = documents
	}

	return comments, nil
}

func (repo *ViolationReportRepository) GetCommentByID(id uint) (*models.ViolationReportComment, error) {
	var comment models.ViolationReportComment
	if err := repo.db.Preload("Documents").First(&comment, id).Error; err != nil {
		return nil, err
	}

	return &comment, nil
}

func (repo *ViolationReportRepository) GetFullCommentDetailsByID(id uint) (*models.ViolationReportComment, error) {

	var comment models.ViolationReportComment

	if err := repo.db.
		Preload("ViolationReport").
		Preload("ViolationReportStatus").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		Preload("Documents").
		Where("id = ?", id).
		First(&comment).Error; err != nil {
		return nil, err
	}

	return &comment, nil
}

func (repo *ViolationReportRepository) DeleteComment(id uint, tx ...utility.DatabaseTransaction) error {

	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("id = ?", id).
		Delete(&models.ViolationReportComment{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *ViolationReportRepository) UpdateComment(comment *models.ViolationReportComment, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Model(&models.ViolationReportComment{}).
		Where("id = ?", comment.ID).
		Updates(comment).Error
}

func (repo *ViolationReportRepository) CreateCommentDocuments(docs []*models.ViolationReportCommentDocument, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Create(&docs).Error
}

func (repo *ViolationReportRepository) DeleteCommentDocuments(docIDs []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Where("id IN (?)", docIDs).Delete(&models.ViolationReportCommentDocument{}).Error
}

func (repo *ViolationReportRepository) GetCommentDocumentByID(id uint) (*models.ViolationReportCommentDocument, error) {
	var document models.ViolationReportCommentDocument
	if err := repo.db.Where("id = ?", id).First(&document).Error; err != nil {
		return nil, err
	}

	return &document, nil
}
