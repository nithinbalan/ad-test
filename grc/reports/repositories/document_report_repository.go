package repositories

import (
	"fmt"
	reportDTOs "jethur/grc/reports/dtos"

	"gorm.io/gorm"
)

type DocumentReportRepositoryInterface interface {
	GetDocumentsByStatus(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.DocumentByStatusDTO, error)
	GetDocumentsByType(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.DocumentByTypeDTO, error)
	GetDocumentsByResponsibleUnit(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.DocumentByResponsibleUnitDTO, error)
	GetDocumentsByCategory(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.DocumentByCategoryDTO, error)
	GetDocumentsByResponsibleUser(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.DocumentByResponsibleUserDTO, error)
	GetDocumentsByAccess(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.DocumentByAccessDTO, error)
	GetDocumentsByPolicyAndProcedure(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.DocumentByTypeDTO, error)
	GetDocumentsByPublishedAndArchived(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.DocumentByStatusDTO, error)
	GetDocumentsPublishedByYear(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.DocumentByYearDTO, error)
	GetDocumentsPublishedByMonth(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.DocumentByMonthDTO, error)
	GetDocumentChangeRequestsByType(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.DocumentChangeRequestByTypeDTO, error)
	GetDocumentChangeRequestsByStatus(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.DocumentChangeRequestByStatusDTO, error)
	GetDocumentsUpdatedByChangeType(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.DocumentChangeRequestByTypeDTO, error)
	GetDocumentsArchivedByQuarter(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.DocumentsArchivedByQuarterDTO, error)
}

type DocumentReportRepository struct {
	db *gorm.DB
}

func NewDocumentReportRepository(db *gorm.DB) *DocumentReportRepository {
	return &DocumentReportRepository{
		db: db,
	}
}

func (repo *DocumentReportRepository) GetDocumentsByStatus(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.DocumentByStatusDTO, error) {
	results := []*reportDTOs.DocumentByStatusDTO{}

	query := repo.db.Table("document_statuses").
		Select(`
			document_statuses.id AS document_status_id,
			document_statuses.title_en AS document_status_title_en,
			document_statuses.title_ar AS document_status_title_ar,
			document_statuses.color_code AS document_status_color_code,
			document_statuses.type AS document_status_type,
			COALESCE(COUNT(documents.id), 0) AS count
		`).
		Joins("LEFT JOIN documents ON documents.document_status_id = document_statuses.id").
		Group("document_statuses.id, document_statuses.title_en, document_statuses.title_ar, document_statuses.color_code, document_statuses.type")

	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("documents.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.Scan(&results).Error
	return results, err
}

func (repo *DocumentReportRepository) GetDocumentsByType(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.DocumentByTypeDTO, error) {
	results := []*reportDTOs.DocumentByTypeDTO{}

	dateCondition := ""
	if params.DateFrom != "" && params.DateTo != "" {
		dateCondition = fmt.Sprintf("AND documents.created_at BETWEEN '%s 00:00:00' AND '%s 23:59:59'", params.DateFrom, params.DateTo)
	}

	query := repo.db.Raw(`
		SELECT 
			dt.id as document_type_id,
			dt.title_en as document_type_title_en,
			dt.title_ar as document_type_title_ar,
			COUNT(d.id) as count
		FROM document_types dt
		LEFT JOIN (
			SELECT documents.id, documents.document_type_id 
			FROM documents 
			JOIN document_statuses ON document_statuses.id = documents.document_status_id
			WHERE document_statuses.type = 'published' ` + dateCondition + `
		) d ON d.document_type_id = dt.id
		GROUP BY dt.id, dt.title_en, dt.title_ar
		ORDER BY dt.title_en
	`)

	err := query.Scan(&results).Error
	return results, err
}

func (repo *DocumentReportRepository) GetDocumentsByResponsibleUnit(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.DocumentByResponsibleUnitDTO, error) {
	results := []*reportDTOs.DocumentByResponsibleUnitDTO{}

	query := repo.db.Table("documents").
		Select(`
			organization_units.id AS responsible_unit_id,
			organization_units.title_en AS responsible_unit_title_en,
			organization_units.title_ar AS responsible_unit_title_ar,
			organization_units.color_code AS responsible_unit_color_code,
			COUNT(documents.id) AS count
		`).
		Joins("JOIN organization_units ON organization_units.id = documents.responsible_unit_id").
		Joins("JOIN document_statuses ON document_statuses.id = documents.document_status_id").
		Where("document_statuses.type = ?", "published").
		Group("organization_units.id, organization_units.title_en, organization_units.title_ar, organization_units.color_code")

	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("compliance_register.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.Scan(&results).Error
	return results, err
}

func (repo *DocumentReportRepository) GetDocumentsByCategory(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.DocumentByCategoryDTO, error) {
	results := []*reportDTOs.DocumentByCategoryDTO{}

	dateCondition := ""
	if params.DateFrom != "" && params.DateTo != "" {
		dateCondition = fmt.Sprintf("AND documents.created_at BETWEEN '%s 00:00:00' AND '%s 23:59:59'", params.DateFrom, params.DateTo)
	}

	query := repo.db.Raw(`
		SELECT 
			dc.id as category_id,
			dc.title_en as category_title_en,
			dc.title_ar as category_title_ar,
			COUNT(d.id) as count
		FROM document_categories dc
		LEFT JOIN (
			SELECT documents.id, documents.document_category_id 
			FROM documents 
			JOIN document_statuses ON document_statuses.id = documents.document_status_id
			WHERE document_statuses.type = 'published' ` + dateCondition + `
		) d ON d.document_category_id = dc.id
		GROUP BY dc.id, dc.title_en, dc.title_ar
		ORDER BY dc.title_en
	`)

	err := query.Scan(&results).Error
	return results, err
}

func (repo *DocumentReportRepository) GetDocumentsByResponsibleUser(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.DocumentByResponsibleUserDTO, error) {
	results := []*reportDTOs.DocumentByResponsibleUserDTO{}

	query := repo.db.Table("documents").
		Select(`
			users.id AS responsible_user_id,
			users.first_name_en AS responsible_user_first_name_en,
			users.first_name_ar AS responsible_user_first_name_ar,
			users.last_name_en AS responsible_user_last_name_en,
			users.last_name_ar AS responsible_user_last_name_ar,
			CONCAT(users.first_name_en, ' ', users.last_name_en) AS responsible_user_title_en,
			CONCAT(users.first_name_ar, ' ', users.last_name_ar) AS responsible_user_title_ar,
			users.initial AS responsible_user_initial,
			users.color AS responsible_user_color,
			users.email AS responsible_user_email,
			users.profile_image_path AS responsible_user_profile_image_path,
			COUNT(documents.id) AS count
		`).
		Joins("JOIN users ON users.id = documents.responsible_user_id").
		Joins("JOIN document_statuses ON document_statuses.id = documents.document_status_id").
		Where("document_statuses.type = ?", "published").
		Group("users.id, users.first_name_en, users.first_name_ar, users.last_name_en, users.last_name_ar, users.initial, users.color, users.email, users.profile_image_path")

	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("documents.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.Scan(&results).Error
	return results, err
}

func (repo *DocumentReportRepository) GetDocumentsByAccess(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.DocumentByAccessDTO, error) {
	results := []*reportDTOs.DocumentByAccessDTO{}

	dateCondition := ""
	if params.DateFrom != "" && params.DateTo != "" {
		dateCondition = fmt.Sprintf("AND documents.created_at BETWEEN '%s 00:00:00' AND '%s 23:59:59'", params.DateFrom, params.DateTo)
	}

	query := repo.db.Raw(`
		SELECT 
			dat.id as access_id,
			dat.title_en as access_type_title_en,
			dat.title_ar as access_type_title_ar,
			dat.type as access_type,
			dat.color_code as access_type_color_code,
			COUNT(d.id) as count
		FROM document_access_types dat
		LEFT JOIN (
			SELECT documents.id, documents.document_access_type_id 
			FROM documents 
			JOIN document_statuses ON document_statuses.id = documents.document_status_id
			WHERE document_statuses.type = 'published' ` + dateCondition + `
		) d ON d.document_access_type_id = dat.id
		GROUP BY dat.id, dat.title_en, dat.title_ar, dat.type, dat.color_code
		ORDER BY dat.title_en
	`)

	err := query.Scan(&results).Error
	return results, err
}

func (repo *DocumentReportRepository) GetDocumentsByPolicyAndProcedure(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.DocumentByTypeDTO, error) {
	results := []*reportDTOs.DocumentByTypeDTO{}

	dateCondition := ""
	if params.DateFrom != "" && params.DateTo != "" {
		dateCondition = fmt.Sprintf("AND documents.created_at BETWEEN '%s 00:00:00' AND '%s 23:59:59'", params.DateFrom, params.DateTo)
	}

	query := repo.db.Raw(`
		SELECT 
			dt.id as document_type_id,
			dt.title_en as document_type_title_en,
			dt.title_ar as document_type_title_ar,
			COUNT(d.id) as count
		FROM document_types dt
		LEFT JOIN (
			SELECT documents.id, documents.document_type_id 
			FROM documents 
			JOIN document_statuses ON document_statuses.id = documents.document_status_id
			WHERE document_statuses.type = 'published' ` + dateCondition + `
		) d ON d.document_type_id = dt.id
		WHERE dt.type IN ('policy', 'procedure')
		GROUP BY dt.id, dt.title_en, dt.title_ar
		ORDER BY dt.title_en
	`)

	err := query.Scan(&results).Error
	return results, err
}

func (repo *DocumentReportRepository) GetDocumentsByPublishedAndArchived(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.DocumentByStatusDTO, error) {
	results := []*reportDTOs.DocumentByStatusDTO{}

	dateCondition := ""
	if params.DateFrom != "" && params.DateTo != "" {
		dateCondition = fmt.Sprintf("AND documents.created_at BETWEEN '%s 00:00:00' AND '%s 23:59:59'", params.DateFrom, params.DateTo)
	}

	query := repo.db.Raw(`
		SELECT 
			ds.id as document_status_id,
			ds.title_en as document_status_title_en,
			ds.title_ar as document_status_title_ar,
			ds.type as document_status_type,
			ds.color_code as document_status_color_code,
			COUNT(d.id) as count
		FROM document_statuses ds
		LEFT JOIN (
			SELECT documents.id, documents.document_status_id 
			FROM documents 
			WHERE 1=1 ` + dateCondition + `
		) d ON d.document_status_id = ds.id
		WHERE ds.type IN ('published', 'archived')
		GROUP BY ds.id, ds.title_en, ds.title_ar, ds.type, ds.color_code
		ORDER BY ds.title_en
	`)

	err := query.Scan(&results).Error
	return results, err
}

func (repo *DocumentReportRepository) GetDocumentsPublishedByYear(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.DocumentByYearDTO, error) {
	results := []*reportDTOs.DocumentByYearDTO{}

	// Fetch yearly document count
	var rawResults []struct {
		Year  int  `gorm:"column:year"`
		Count uint `gorm:"column:count"`
	}

	query := repo.db.Table("documents").
		Select(`YEAR(documents.published_at) AS year, COUNT(documents.id) AS count`).
		Joins("JOIN document_statuses ON document_statuses.id = documents.document_status_id").
		Where("document_statuses.type = ? AND documents.published_at IS NOT NULL", "published").
		Group("YEAR(documents.published_at)").
		Order("YEAR(documents.published_at) ASC")

	// Apply date filter if provided
	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("documents.published_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.Scan(&rawResults).Error
	if err != nil {
		return nil, err
	}

	// Convert year to English and Arabic title
	for _, record := range rawResults {
		results = append(results, &reportDTOs.DocumentByYearDTO{
			TitleEn: fmt.Sprintf("%d", record.Year), // English Title (Year)
			TitleAr: fmt.Sprintf("%d", record.Year), // Arabic Title (Same as Year)
			Count:   int64(record.Count),
		})
	}

	return results, nil
}

func (repo *DocumentReportRepository) GetDocumentsPublishedByMonth(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.DocumentByMonthDTO, error) {
	results := []*reportDTOs.DocumentByMonthDTO{}

	// Fetch monthly document count for last 12 months
	var rawResults []struct {
		Year  int  `gorm:"column:year"`
		Month int  `gorm:"column:month"`
		Count uint `gorm:"column:count"`
	}

	query := repo.db.Raw(`
		WITH months AS (
			SELECT 
				YEAR(GETDATE()) as year,
				MONTH(GETDATE()) as month
			UNION ALL
			SELECT 
				CASE 
					WHEN month = 1 THEN year - 1 
					ELSE year 
				END,
				CASE 
					WHEN month = 1 THEN 12 
					ELSE month - 1 
				END
			FROM months
			WHERE (year * 12 + month) > (YEAR(GETDATE()) * 12 + MONTH(GETDATE()) - 12)
		)
		SELECT 
			m.year,
			m.month,
			COALESCE(COUNT(d.id), 0) as count
		FROM months m
		LEFT JOIN (
			SELECT 
				documents.id,
				YEAR(documents.published_at) as doc_year,
				MONTH(documents.published_at) as doc_month
			FROM documents
			JOIN document_statuses ON document_statuses.id = documents.document_status_id
			WHERE document_statuses.type = 'published'
			AND documents.published_at IS NOT NULL
		) d ON d.doc_year = m.year AND d.doc_month = m.month
		GROUP BY m.year, m.month
		ORDER BY m.year DESC, m.month DESC
		OPTION (MAXRECURSION 12)
	`)

	err := query.Scan(&rawResults).Error
	if err != nil {
		return nil, err
	}

	// Month names in English and Arabic
	monthNamesEN := []string{"January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"}
	monthNamesAR := []string{"يناير", "فبراير", "مارس", "أبريل", "مايو", "يونيو", "يوليو", "أغسطس", "سبتمبر", "أكتوبر", "نوفمبر", "ديسمبر"}

	// Convert results to DTO format
	for _, record := range rawResults {
		monthIndex := record.Month - 1 // Convert to 0-based index
		results = append(results, &reportDTOs.DocumentByMonthDTO{
			TitleEn: fmt.Sprintf("%s %d", monthNamesEN[monthIndex], record.Year),
			TitleAr: fmt.Sprintf("%s %d", monthNamesAR[monthIndex], record.Year),
			Count:   int64(record.Count),
		})
	}

	return results, nil
}

func (repo *DocumentReportRepository) GetDocumentChangeRequestsByType(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.DocumentChangeRequestByTypeDTO, error) {
	results := []*reportDTOs.DocumentChangeRequestByTypeDTO{}

	dateCondition := ""
	if params.DateFrom != "" && params.DateTo != "" {
		dateCondition = fmt.Sprintf("AND document_change_requests.created_at BETWEEN '%s 00:00:00' AND '%s 23:59:59'", params.DateFrom, params.DateTo)
	}

	query := repo.db.Raw(`
		SELECT 
			dcrt.id as document_change_request_type_id,
			dcrt.title_en as document_change_request_type_title_en,
			dcrt.title_ar as document_change_request_type_title_ar,
			dcrt.type as document_change_request_type_type,
			dcrt.status as document_change_request_type_status,
			COUNT(dcr.id) as count
		FROM document_change_request_types dcrt
		LEFT JOIN (
			SELECT id, document_change_request_type_id 
			FROM document_change_requests 
			WHERE 1=1 ` + dateCondition + `
		) dcr ON dcr.document_change_request_type_id = dcrt.id
		WHERE dcrt.status = 'active'
		GROUP BY dcrt.id, dcrt.title_en, dcrt.title_ar, dcrt.type, dcrt.status
		ORDER BY dcrt.title_en
	`)

	err := query.Scan(&results).Error
	return results, err
}

func (repo *DocumentReportRepository) GetDocumentChangeRequestsByStatus(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.DocumentChangeRequestByStatusDTO, error) {
	results := []*reportDTOs.DocumentChangeRequestByStatusDTO{}

	dateCondition := ""
	if params.DateFrom != "" && params.DateTo != "" {
		dateCondition = fmt.Sprintf("AND document_change_requests.created_at BETWEEN '%s 00:00:00' AND '%s 23:59:59'", params.DateFrom, params.DateTo)
	}

	query := repo.db.Raw(`
		SELECT 
			dcrs.id as document_change_request_status_id,
			dcrs.title_en as document_change_request_status_title_en,
			dcrs.title_ar as document_change_request_status_title_ar,
			dcrs.type as document_change_request_status_type,
			dcrs.color_code as document_change_request_status_color_code,
			COUNT(dcr.id) as count
		FROM document_change_request_statuses dcrs
		LEFT JOIN (
			SELECT id, document_change_request_status_id 
			FROM document_change_requests 
			WHERE 1=1 ` + dateCondition + `
		) dcr ON dcr.document_change_request_status_id = dcrs.id
		GROUP BY dcrs.id, dcrs.title_en, dcrs.title_ar, dcrs.type, dcrs.color_code
		ORDER BY dcrs.title_en
	`)

	err := query.Scan(&results).Error
	return results, err
}

func (repo *DocumentReportRepository) GetDocumentsUpdatedByChangeType(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.DocumentChangeRequestByTypeDTO, error) {
	results := []*reportDTOs.DocumentChangeRequestByTypeDTO{}

	dateCondition := ""
	if params.DateFrom != "" && params.DateTo != "" {
		dateCondition = fmt.Sprintf("AND dcr.created_at BETWEEN '%s 00:00:00' AND '%s 23:59:59'", params.DateFrom, params.DateTo)
	}

	query := repo.db.Raw(`
		SELECT 
			dcrt.id as document_change_request_type_id,
			dcrt.title_en as document_change_request_type_title_en,
			dcrt.title_ar as document_change_request_type_title_ar,
			dcrt.type as document_change_request_type_type,
			COUNT(dcr.id) as count
		FROM document_change_request_types dcrt
		LEFT JOIN (
			SELECT 
				document_change_requests.id,
				document_change_requests.document_change_request_type_id
			FROM document_change_requests 
			JOIN document_change_request_statuses dcrs 
				ON dcrs.id = document_change_requests.document_change_request_status_id
			WHERE dcrs.type = 'published' ` + dateCondition + `
		) dcr ON dcr.document_change_request_type_id = dcrt.id
		WHERE dcrt.type IN ('major', 'minor', 'archive')
		AND dcrt.status = 'active'
		GROUP BY dcrt.id, dcrt.title_en, dcrt.title_ar, dcrt.type
		ORDER BY dcrt.title_en
	`)

	err := query.Scan(&results).Error
	return results, err
}

func (repo *DocumentReportRepository) GetDocumentsArchivedByQuarter(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.DocumentsArchivedByQuarterDTO, error) {
	results := []*reportDTOs.DocumentsArchivedByQuarterDTO{}

	query := repo.db.Raw(`
		WITH quarters AS (
			SELECT 
				YEAR(DATEADD(QUARTER, -3, GETDATE())) as year,
				DATEPART(QUARTER, DATEADD(QUARTER, -3, GETDATE())) as quarter
			UNION ALL
			SELECT 
				YEAR(DATEADD(QUARTER, -2, GETDATE())) as year,
				DATEPART(QUARTER, DATEADD(QUARTER, -2, GETDATE())) as quarter
			UNION ALL
			SELECT 
				YEAR(DATEADD(QUARTER, -1, GETDATE())) as year,
				DATEPART(QUARTER, DATEADD(QUARTER, -1, GETDATE())) as quarter
			UNION ALL
			SELECT 
				YEAR(GETDATE()) as year,
				DATEPART(QUARTER, GETDATE()) as quarter
		)
		SELECT 
			q.year,
			q.quarter,
			CONCAT('Q', q.quarter, ' ', q.year) as title_en,
			CONCAT('الربع ', q.quarter, ' ', q.year) as title_ar,
			COALESCE(COUNT(dcr.id), 0) as count
		FROM quarters q
		LEFT JOIN (
			SELECT 
				dcr.id,
				YEAR(dcr.published_at) as pub_year,
				DATEPART(QUARTER, dcr.published_at) as pub_quarter
			FROM document_change_requests dcr
			JOIN document_change_request_types dcrt ON dcrt.id = dcr.document_change_request_type_id
			JOIN document_change_request_statuses dcrs ON dcrs.id = dcr.document_change_request_status_id
			WHERE dcrt.type = 'archive'
			AND dcrs.type = 'published'
			AND dcr.published_at IS NOT NULL
		) dcr ON dcr.pub_year = q.year AND dcr.pub_quarter = q.quarter
		GROUP BY q.year, q.quarter
		ORDER BY q.year DESC, q.quarter DESC
	`)

	err := query.Scan(&results).Error
	return results, err
}
