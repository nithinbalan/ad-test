package repositories

import (
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/compliance/dtos"
	"jethur/grc/models"
	"jethur/utility"
	"time"

	"gorm.io/gorm"
)

type ComplianceReportRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(params *dtos.ComplianceReportListQueryParamsDTO) (*dtos.ComplianceReportListResponseDTO, error)
	Create(complianceReport *models.ComplianceReport, tx ...utility.DatabaseTransaction) (uint, error)
	CreateSection(section *models.ComplianceReportSection, tx ...utility.DatabaseTransaction) (uint, error)
	GetComplianceAssessmentByComplianceRegisterID(complianceRegisterID uint, startDate time.Time, endDate time.Time, tx ...utility.DatabaseTransaction) ([]models.ComplianceAssessment, error)
	CreateComplianceReportDepartment(reportDepartment *models.ComplianceReportDepartment, tx ...utility.DatabaseTransaction) (uint, error)
	CreateComplianceReportClause(reportClause *models.ComplianceReportClause, tx ...utility.DatabaseTransaction) (uint, error)
	CreateComplianceReportRequirement(reportRequirement *models.ComplianceReportRequirement, tx ...utility.DatabaseTransaction) (uint, error)
	CreateComplianceReportActionPlan(reportActionPlan *models.ComplianceReportActionPlan, tx ...utility.DatabaseTransaction) (uint, error)

	GetComplianceReportByID(complianceReportID uint) (*models.ComplianceReport, error)
	GetComplianceReportDepartments(complianceReportID uint) ([]*models.ComplianceReportDepartment, error)
	GetComplianceReportSectionByID(sectionID uint) (*models.ComplianceReportSection, error)
	UpdateComplianceReportSection(section *models.ComplianceReportSection, tx ...utility.DatabaseTransaction) error
	DeleteComplianceReport(complianceReportID uint, tx ...utility.DatabaseTransaction) error
	DeleteComplianceReportSections(complianceReportID uint, tx ...utility.DatabaseTransaction) error
	DeleteComplianceReportDepartments(complianceReportID uint, tx ...utility.DatabaseTransaction) error
	DeleteComplianceReportClauses(complianceReportID uint, tx ...utility.DatabaseTransaction) error
	DeleteComplianceReportRequirements(complianceReportID uint, tx ...utility.DatabaseTransaction) error
	DeleteComplianceReportActionPlans(complianceReportID uint, tx ...utility.DatabaseTransaction) error
}

type ComplianceReportRepository struct {
	db *gorm.DB
}

func NewComplianceReportRepository(db *gorm.DB) *ComplianceReportRepository {
	return &ComplianceReportRepository{
		db: db,
	}
}

func (repo *ComplianceReportRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *ComplianceReportRepository) GetAll(params *dtos.ComplianceReportListQueryParamsDTO) (*dtos.ComplianceReportListResponseDTO, error) {
	db := repo.db
	var reports []*dtos.ComplianceReportListDTO
	var totalCount int64
	query := db.Table("compliance_reports").
		Select(`
			compliance_reports.id,
			compliance_reports.title,
			CONVERT(VARCHAR, compliance_reports.start_date, 23) AS start_date,
			CONVERT(VARCHAR, compliance_reports.end_date, 23) AS end_date,
			compliance_reports.created_at,
			compliance_register.title_en AS compliance_register_title_en,
			compliance_register.title_ar AS compliance_register_title_ar,
			created_user.id AS created_user_id,
			created_user.first_name_en AS created_user_first_name_en,
			created_user.first_name_ar AS created_user_first_name_ar,
			created_user.last_name_en AS created_user_last_name_en,	
			created_user.last_name_ar AS created_user_last_name_ar,
			created_user.initial AS created_user_initial,
			created_user.color AS created_user_color,
			created_user.email AS created_user_email,
			created_user.profile_image_path AS created_user_profile_image_path
		`).
		Joins("LEFT JOIN compliance_register ON compliance_register.id = compliance_reports.compliance_register_id").
		Joins("LEFT JOIN users AS created_user ON created_user.id = compliance_reports.created_by_id")

	if params.Keyword != "" {
		query = query.Where(`
			(compliance_reports.title LIKE ? OR 
			compliance_register.title_en LIKE ? OR 
			compliance_register.title_ar LIKE ? OR 
			created_user.first_name_en LIKE ? OR 
			created_user.first_name_ar LIKE ? OR 
			created_user.last_name_en LIKE ? OR 
			created_user.last_name_ar LIKE ? OR 
			created_user.initial LIKE ? OR 
			created_user.color LIKE ? OR 
			created_user.email LIKE ? OR 	
			created_user.profile_image_path LIKE ?)`,
			"%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%",
			"%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%",
			"%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%",
			"%"+params.Keyword+"%", "%"+params.Keyword+"%")
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

	response := &dtos.ComplianceReportListResponseDTO{
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

func (repo *ComplianceReportRepository) GetComplianceAssessmentByComplianceRegisterID(complianceRegisterID uint, startDate time.Time, endDate time.Time, tx ...utility.DatabaseTransaction) ([]models.ComplianceAssessment, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var complianceAssessments []models.ComplianceAssessment
	result := db.Preload("AssessmentRequirements").
		Preload("AssessmentRequirements.ComplianceRequirement").
		Preload("AssessmentRequirements.ActionPlans").
		Preload("AssessmentRequirements.ComplianceRequirement.ComplianceClause").
		Preload("AssessmentRequirements.ComplianceStatus").
		Where("compliance_register_id = ?", complianceRegisterID).
		Where("created_at > ? AND created_at < ?", startDate, endDate).
		Find(&complianceAssessments)
	if result.Error != nil {
		return nil, result.Error
	}

	return complianceAssessments, nil
}

func (repo *ComplianceReportRepository) Create(complianceReport *models.ComplianceReport, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	result := db.Create(complianceReport)
	if result.Error != nil {
		return 0, result.Error
	}

	return complianceReport.ID, nil

}

func (repo *ComplianceReportRepository) CreateSection(section *models.ComplianceReportSection, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	result := db.Create(section)
	if result.Error != nil {
		return 0, result.Error
	}

	return section.ID, nil
}

func (repo *ComplianceReportRepository) CreateComplianceReportDepartment(reportDepartment *models.ComplianceReportDepartment, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	result := db.Create(reportDepartment)
	if result.Error != nil {
		return 0, result.Error
	}

	return reportDepartment.ID, nil
}

func (repo *ComplianceReportRepository) CreateComplianceReportClause(reportClause *models.ComplianceReportClause, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	result := db.Create(reportClause)
	if result.Error != nil {
		return 0, result.Error
	}

	return reportClause.ID, nil
}

func (repo *ComplianceReportRepository) CreateComplianceReportRequirement(reportRequirement *models.ComplianceReportRequirement, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	result := db.Create(reportRequirement)
	if result.Error != nil {
		return 0, result.Error
	}

	return reportRequirement.ID, nil
}

func (repo *ComplianceReportRepository) CreateComplianceReportActionPlan(reportActionPlan *models.ComplianceReportActionPlan, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	result := db.Create(reportActionPlan)
	if result.Error != nil {
		return 0, result.Error
	}

	return reportActionPlan.ID, nil
}

func (repo *ComplianceReportRepository) GetComplianceReportByID(complianceReportID uint) (*models.ComplianceReport, error) {

	db := repo.db

	var complianceReport models.ComplianceReport
	result := db.Preload("Sections").
		Preload("ComplianceRegister").
		Preload("CreatedByUser").
		Where("id = ?", complianceReportID).First(&complianceReport)

	if result.Error != nil {
		return nil, result.Error
	}

	return &complianceReport, nil
}

func (repo *ComplianceReportRepository) GetComplianceReportDepartments(complianceReportID uint) ([]*models.ComplianceReportDepartment, error) {
	var departments []*models.ComplianceReportDepartment

	result := repo.db.
		Preload("OrganizationUnit").
		Preload("ComplianceReportClauses").
		Preload("ComplianceReportClauses.ComplianceReportNonCompliantRequirements", func(db *gorm.DB) *gorm.DB {
			return db.Where("score < ? AND compliance_status_id IS NOT NULL AND compliance_status_id NOT IN (SELECT id FROM compliance_statuses WHERE type = 'not_applicable')", 60)
		}).
		Preload("ComplianceReportClauses.ComplianceReportNonCompliantRequirements.ComplianceReportActionPlans").
		Preload("ComplianceReportClauses.ComplianceReportNonCompliantRequirements.ComplianceReportActionPlans.ActionPlanStatus").
		Preload("ComplianceReportClauses.ComplianceReportNonCompliantRequirements.ComplianceReportActionPlans.SeverityLevel").
		Preload("ComplianceReportClauses.ComplianceReportRequirements").
		Preload("ComplianceReportClauses.ComplianceReportRequirements.ComplianceStatus").
		Preload("ComplianceReportClauses.ComplianceReportRequirements.ComplianceReportActionPlans").
		Preload("ComplianceReportClauses.ComplianceReportRequirements.ComplianceReportActionPlans.ActionPlanStatus").
		Preload("ComplianceReportClauses.ComplianceReportRequirements.ComplianceReportActionPlans.SeverityLevel").
		Where("compliance_report_id = ?", complianceReportID).
		Find(&departments)

	if result.Error != nil {
		return nil, result.Error
	}

	return departments, nil
}

func (repo *ComplianceReportRepository) GetComplianceReportSectionByID(sectionID uint) (*models.ComplianceReportSection, error) {
	var section models.ComplianceReportSection
	result := repo.db.Where("id = ?", sectionID).First(&section)
	if result.Error != nil {
		return nil, result.Error
	}

	return &section, nil
}

func (repo *ComplianceReportRepository) UpdateComplianceReportSection(section *models.ComplianceReportSection, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	result := db.Save(section)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *ComplianceReportRepository) DeleteComplianceReport(complianceReportID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	result := db.Delete(&models.ComplianceReport{}, complianceReportID)
	return result.Error
}

func (repo *ComplianceReportRepository) DeleteComplianceReportSections(complianceReportID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	result := db.Where("compliance_report_id = ?", complianceReportID).Delete(&models.ComplianceReportSection{})
	return result.Error
}

func (repo *ComplianceReportRepository) DeleteComplianceReportDepartments(complianceReportID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	result := db.Where("compliance_report_id = ?", complianceReportID).Delete(&models.ComplianceReportDepartment{})
	return result.Error
}

func (repo *ComplianceReportRepository) DeleteComplianceReportClauses(complianceReportID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	result := db.Where("compliance_report_id = ?", complianceReportID).Delete(&models.ComplianceReportClause{})
	return result.Error
}

func (repo *ComplianceReportRepository) DeleteComplianceReportRequirements(complianceReportID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	result := db.Where("compliance_report_id = ?", complianceReportID).Delete(&models.ComplianceReportRequirement{})
	return result.Error
}

func (repo *ComplianceReportRepository) DeleteComplianceReportActionPlans(complianceReportID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	result := db.Where("compliance_report_id = ?", complianceReportID).Delete(&models.ComplianceReportActionPlan{})
	return result.Error
}
