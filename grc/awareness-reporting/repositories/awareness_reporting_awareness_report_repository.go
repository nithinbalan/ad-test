package repositories

import (
	"fmt"
	awarenessDTOs "jethur/grc/awareness-reporting/dtos"
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/models"
	"jethur/utility"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type AwarenessReportRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(params *awarenessDTOs.AwarenessReportListQueryParamsDTO) (*awarenessDTOs.AwarenessReportListResponseDTO, error)
	GetByID(id uint) (*models.AwarenessReport, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.AwarenessReport, error)
	Create(report *models.AwarenessReport, tx ...utility.DatabaseTransaction) (uint, error)
	Update(report *models.AwarenessReport, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error
	GenerateReferenceNumber() (string, error)
}

type AwarenessReportRepository struct {
	db *gorm.DB
}

func NewAwarenessReportRepository(db *gorm.DB) *AwarenessReportRepository {
	return &AwarenessReportRepository{
		db: db,
	}
}

func (repo *AwarenessReportRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *AwarenessReportRepository) GetAll(params *awarenessDTOs.AwarenessReportListQueryParamsDTO) (*awarenessDTOs.AwarenessReportListResponseDTO, error) {
	var reports []*awarenessDTOs.AwarenessReportListDTO
	var totalCount int64

	query := repo.db.Table("awareness_reports").
		Select(`
			awareness_reports.*,
			awareness_campaign_types.title_en AS awareness_report_type_title_en,
			awareness_campaign_types.title_ar AS awareness_report_type_title_ar,
			awareness_report_sub_categories.title_en AS sub_category_title_en,
			awareness_report_sub_categories.title_ar AS sub_category_title_ar,
			severity_levels.title_en AS severity_title_en,
			severity_levels.title_ar AS severity_title_ar,
			severity_levels.color_code AS severity_color_code,
			awareness_handling_statuses.title_en AS awareness_handling_status_title_en,
			awareness_handling_statuses.title_ar AS awareness_handling_status_title_ar,
			handling_owner_user.id AS handling_owner_id,
			handling_owner_user.first_name_en AS handling_owner_first_name_en,
			handling_owner_user.first_name_ar AS handling_owner_first_name_ar,
			handling_owner_user.last_name_en AS handling_owner_last_name_en,
			handling_owner_user.last_name_ar AS handling_owner_last_name_ar,
			handling_owner_user.initial AS handling_owner_initial,
			handling_owner_user.color AS handling_owner_color,
			created_user.id AS created_user_id,
			created_user.first_name_en AS created_user_first_name_en,
			created_user.first_name_ar AS created_user_first_name_ar,
			created_user.last_name_en AS created_user_last_name_en,
			created_user.last_name_ar AS created_user_last_name_ar,
			created_user.initial AS created_user_initial,
			created_user.color AS created_user_color,
			updated_user.id AS updated_user_id,
			updated_user.first_name_en AS updated_user_first_name_en,
			updated_user.first_name_ar AS updated_user_first_name_ar,
			updated_user.last_name_en AS updated_user_last_name_en,
			updated_user.last_name_ar AS updated_user_last_name_ar,
			updated_user.initial AS updated_user_initial,
			updated_user.color AS updated_user_color
		`).
		Joins(`
			LEFT JOIN awareness_campaign_types ON awareness_campaign_types.id = awareness_reports.awareness_report_type_id
			LEFT JOIN awareness_report_sub_categories ON awareness_report_sub_categories.id = awareness_reports.sub_category_id
			LEFT JOIN severity_levels ON severity_levels.id = awareness_reports.severity_id
			LEFT JOIN awareness_handling_statuses ON awareness_handling_statuses.id = awareness_reports.awareness_handling_status_id
			LEFT JOIN users AS handling_owner_user ON handling_owner_user.id = awareness_reports.handling_owner_id
			LEFT JOIN users AS created_user ON created_user.id = awareness_reports.created_by_id
			LEFT JOIN users AS updated_user ON updated_user.id = awareness_reports.updated_by_id
		`)

	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where(columnName+" IN (?)", values)
		}
	}

	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("awareness_reports.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
			(awareness_reports.reference_number LIKE ? OR
			awareness_reports.report_source_channel LIKE ? OR
			awareness_campaign_types.title_en LIKE ? OR
			awareness_campaign_types.title_ar LIKE ? OR
			awareness_report_sub_categories.title_en LIKE ? OR
			awareness_report_sub_categories.title_ar LIKE ? OR
			severity_levels.title_en LIKE ? OR
			severity_levels.title_ar LIKE ? OR
			awareness_handling_statuses.title_en LIKE ? OR
			awareness_handling_statuses.title_ar LIKE ? OR
			handling_owner_user.first_name_en LIKE ? OR
			handling_owner_user.first_name_ar LIKE ? OR
			handling_owner_user.last_name_en LIKE ? OR
			handling_owner_user.last_name_ar LIKE ? OR
			handling_owner_user.email LIKE ?)`,
			keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword,
			keyword, keyword, keyword, keyword, keyword, keyword, keyword)
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

	response := &awarenessDTOs.AwarenessReportListResponseDTO{
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

func (repo *AwarenessReportRepository) GetByID(id uint) (*models.AwarenessReport, error) {
	var report *models.AwarenessReport

	result := repo.db.First(&report, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return report, nil
}

func (repo *AwarenessReportRepository) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.AwarenessReport, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var report *models.AwarenessReport

	result := db.
		Preload("AwarenessReportType").
		Preload("SubCategory").
		Preload("SeverityLevel").
		Preload("AwarenessHandlingStatus").
		Preload("HandlingOwner").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		First(&report, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return report, nil
}

func (repo *AwarenessReportRepository) Create(report *models.AwarenessReport, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(report).Error; err != nil {
		return 0, err
	}

	return report.ID, nil
}

func (repo *AwarenessReportRepository) Update(report *models.AwarenessReport, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(report).Error; err != nil {
		return err
	}

	return nil
}

func (repo *AwarenessReportRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Delete(&models.AwarenessReport{}, id).Error; err != nil {
		return err
	}

	return nil
}

func (repo *AwarenessReportRepository) GenerateReferenceNumber() (string, error) {
	var lastReferenceNumber string
	err := repo.db.Model(&models.AwarenessReport{}).
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
	newReferenceNumber := fmt.Sprintf("%s-%d-%04d", "AR", currentYear, lastReferenceID+1)

	return newReferenceNumber, nil
}