package repositories

import (
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/models"
	organizationDTOs "jethur/grc/organization/dtos"
	violationReportDTOs "jethur/grc/violation-report/dtos"
	"jethur/utility"

	"gorm.io/gorm"
)

type ViolationReportInvestigationRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(violationReportID uint, params *violationReportDTOs.ViolationReportInvestigationListQueryParamsDTO) (*violationReportDTOs.ViolationReportInvestigationListResponseDTO, error)
	GetByID(violationReportID uint) (*models.ViolationReportInvestigation, error)
	GetFullDetailsByID(violationReportID uint, tx ...utility.DatabaseTransaction) (*models.ViolationReportInvestigation, error)
	Create(violationReportInvestigation *models.ViolationReportInvestigation, tx ...utility.DatabaseTransaction) (uint, error)
	Update(violationReportInvestigation *models.ViolationReportInvestigation, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error

	//Department
	CreateDepartments(departments []*models.ViolationReportInvestigationDepartment, tx ...utility.DatabaseTransaction) error
	DeleteDepartments(departmentIds []*uint, tx ...utility.DatabaseTransaction) error

	//Document
	CreateDocuments(docs []*models.ViolationReportInvestigationDocument, tx ...utility.DatabaseTransaction) error
	GetDocumentByID(id uint) (*models.ViolationReportInvestigationDocument, error)
	DeleteDocuments(docIds []*uint, tx ...utility.DatabaseTransaction) error
}

type ViolationReportInvestigationRepository struct {
	db *gorm.DB
}

func NewViolationReportInvestigationRepository(db *gorm.DB) *ViolationReportInvestigationRepository {
	return &ViolationReportInvestigationRepository{
		db: db,
	}
}

// BeginTransaction starts a new transaction.
func (repo *ViolationReportInvestigationRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *ViolationReportInvestigationRepository) GetAll(violationReportID uint, params *violationReportDTOs.ViolationReportInvestigationListQueryParamsDTO) (*violationReportDTOs.ViolationReportInvestigationListResponseDTO, error) {
	var totalCount int64

	// Main table query with all joins and selections
	query := repo.db.Table("violation_report_investigations").
		Select(`
        violation_report_investigations.id,
        violation_report_investigations.violation_report_id,
        violation_report_investigations.steps_taken,
        violation_report_investigations.suggested_recommendation,
        violation_report_investigations.finding_summary,
        violation_report_investigations.disciplinary_procedure,
        violation_report_investigations.is_disciplinary_regulatory_compliance,
        violation_report_investigations.actual_recommendation,
        violation_report_investigations.note_and_observations,
		violation_report_investigations.id_of_reporter,
		violation_report_investigations.name_of_reporter,
		violation_report_investigations.phone_of_reporter,
		violation_report_investigations.email_of_reporter,
        violation_report_investigations.created_by_id,
        violation_report_investigations.updated_by_id,
        violation_report_investigations.created_at,
        violation_report_investigations.updated_at,
		CONVERT(VARCHAR, violation_report_investigations.investigation_date, 23) AS investigation_date,
        violation_report.title_en as violation_title,
		violation_report.title_ar as violation_report_title_ar,
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
        LEFT JOIN violation_reports ON violation_reports.id = violation_report_investigations.violation_report_id
		LEFT JOIN users AS created_user ON created_user.id = violation_report_investigations.created_by_id
		LEFT JOIN users AS updated_user ON updated_user.id = violation_report_investigations.updated_by_id
    `).Where("violation_report_investigations.violation_report_id = ?", violationReportID)

	// Base query for count (after main table query)
	baseQuery := repo.db.Table("violation_report_investigations").
		Where("violation_report_investigations.violation_report_id = ?", violationReportID)

	// Apply dynamic filters to the count query
	for columnName, values := range params.Filters {
		if columnName == "department_id" && len(values) > 0 {
			// For department filter, we need to join with departments table
			baseQuery = baseQuery.Joins("INNER JOIN violation_report_investigation_departments vrid ON vrid.violation_report_investigation_id = violation_report_investigations.id").
				Where("vrid.organization_unit_id IN (?)", values)
		} else if len(values) > 0 {
			baseQuery = baseQuery.Where(columnName+" IN (?)", values)
		}
	}

	// Apply other filters
	if params.InvestigationDate != "" {
		baseQuery = baseQuery.Where("violation_report_investigations.investigation_date >= ?", params.InvestigationDate)
	}

	// Calculate total count for pagination
	if err := baseQuery.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	// Apply dynamic filters to the query (except department_id which is handled separately)
	for columnName, values := range params.Filters {
		if columnName != "department_id" && len(values) > 0 {
			query = query.Where(columnName+" IN (?)", values)
		}
	}

	// Apply department filter if specified
	if deptIDs, ok := params.Filters["department_id"]; ok && len(deptIDs) > 0 {
		query = query.Joins("INNER JOIN violation_report_investigation_departments vrid ON vrid.violation_report_investigation_id = violation_report_investigations.id").
			Where("vrid.organization_unit_id IN (?)", deptIDs)
	}

	// Apply other filters
	if params.InvestigationDate != "" {
		query = query.Where("violation_report_investigations.investigation_date >= ?", params.InvestigationDate)
	}

	// Apply search keyword if exists
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
			(violation_report_investigations.steps_taken LIKE ? OR
			violation_report.title_en LIKE ? OR
			violation_report.title_ar LIKE ? OR
			violation_report_investigations.suggested_recommendation LIKE ? OR
			violation_report_investigations.finding_summary LIKE ? OR
			violation_report_investigations.disciplinary_procedure LIKE ? OR
			violation_report_investigations.actual_recommendation LIKE ? OR
			violation_report_investigations.note_and_observations LIKE ? OR
			created_user.first_name_en LIKE ? OR
			created_user.first_name_ar LIKE ? OR
			created_user.last_name_en LIKE ? OR
			created_user.last_name_ar LIKE ? OR
			updated_user.first_name_en LIKE ? OR
			updated_user.first_name_ar LIKE ? OR
			updated_user.last_name_en LIKE ? OR
			updated_user.last_name_ar LIKE ?)`,
			keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword,
			keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword)
	}

	// Sort by the specified column and order
	if params.SortBy == "" {
		params.SortBy = "created_at" // Default sort by id if no sortBy is specified
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "desc" // Default order is ascending if not specified
	}

	query = query.Order("violation_report_investigations." + params.SortBy + " " + params.Sort)

	offset := (params.Page - 1) * params.PageSize

	// Execute the query to get basic investigation data
	var violationReportInvestigations []*violationReportDTOs.ViolationReportInvestigationListDTO
	err := query.Offset(offset).Limit(params.PageSize).Find(&violationReportInvestigations).Error
	if err != nil {
		return nil, err
	}

	// Get department data for each investigation
	for _, investigation := range violationReportInvestigations {
		var departments []struct {
			ID      uint   `json:"id"`
			TitleEn string `json:"title_en"`
			TitleAr string `json:"title_ar"`
		}

		err := repo.db.Table("violation_report_investigation_departments vrid").
			Select("ou.id, ou.title_en, ou.title_ar").
			Joins("INNER JOIN organization_units ou ON ou.id = vrid.organization_unit_id").
			Where("vrid.violation_report_investigation_id = ?", investigation.ID).
			Find(&departments).Error

		if err != nil {
			return nil, err
		}

		// Convert to DTO format
		for _, dept := range departments {
			investigation.Departments = append(investigation.Departments, organizationDTOs.OrganizationUnitDTO{
				ID:      dept.ID,
				TitleEn: dept.TitleEn,
				TitleAr: dept.TitleAr,
			})
		}
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	response := &violationReportDTOs.ViolationReportInvestigationListResponseDTO{
		Items: violationReportInvestigations,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}
	return response, nil
}

func (repo *ViolationReportInvestigationRepository) GetByID(violationReportID uint) (*models.ViolationReportInvestigation, error) {
	var violationReportInvestigation *models.ViolationReportInvestigation

	result := repo.db.Where("violation_report_id = ?", violationReportID).First(&violationReportInvestigation)
	if result.Error != nil {
		return nil, result.Error
	}

	return violationReportInvestigation, nil

}

func (repo *ViolationReportInvestigationRepository) GetFullDetailsByID(violationReportID uint, tx ...utility.DatabaseTransaction) (*models.ViolationReportInvestigation, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var violationReportInvestigation *models.ViolationReportInvestigation

	result := db.Where("violation_report_id = ?", violationReportID).
		Preload("ViolationReport").
		Preload("Documents").
		Preload("Departments").
		Preload("Departments.Department").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		First(&violationReportInvestigation)

	if result.Error != nil {
		return nil, result.Error
	}

	return violationReportInvestigation, nil
}

func (repo *ViolationReportInvestigationRepository) Create(violationReportInvestigation *models.ViolationReportInvestigation, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(violationReportInvestigation).Error; err != nil {
		return 0, err
	}

	return violationReportInvestigation.ID, nil
}

func (repo *ViolationReportInvestigationRepository) Update(violationReportInvestigation *models.ViolationReportInvestigation, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(violationReportInvestigation).Error; err != nil {
		return err
	}

	return nil
}

func (repo *ViolationReportInvestigationRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	// Delete the violation report investigation
	if err := db.Delete(&models.ViolationReportInvestigation{}, id).Error; err != nil {
		return err
	}

	return nil
}

// Department func
func (repo *ViolationReportInvestigationRepository) CreateDepartments(departments []*models.ViolationReportInvestigationDepartment, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	// Save all departments in a single transaction
	if err := db.Create(&departments).Error; err != nil {
		return err
	}
	return nil
}

func (repo *ViolationReportInvestigationRepository) DeleteDepartments(departmentIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	// Delete departments by IDs
	if err := db.Where("id IN (?)", departmentIds).Delete(&models.ViolationReportInvestigationDepartment{}).Error; err != nil {
		return err
	}

	return nil
}

// Document func
func (repo *ViolationReportInvestigationRepository) GetDocumentByID(id uint) (*models.ViolationReportInvestigationDocument, error) {
	var document models.ViolationReportInvestigationDocument
	if err := repo.db.Where("id = ?", id).First(&document).Error; err != nil {
		return nil, err
	}
	return &document, nil
}

func (repo *ViolationReportInvestigationRepository) CreateDocuments(docs []*models.ViolationReportInvestigationDocument, tx ...utility.DatabaseTransaction) error {
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

func (repo *ViolationReportInvestigationRepository) DeleteDocuments(docIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	// Delete documents by IDs
	if err := db.Where("id IN (?)", docIds).Delete(&models.ViolationReportInvestigationDocument{}).Error; err != nil {
		return err
	}

	return nil
}
