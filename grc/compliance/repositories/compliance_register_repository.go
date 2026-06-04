package repositories

import (
	"fmt"
	commonDTOs "jethur/grc/common/dtos"
	complianceDTOs "jethur/grc/compliance/dtos"
	"jethur/grc/models"
	"jethur/utility"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type ComplianceRegisterRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(params *complianceDTOs.ComplianceRegisterListQueryParamsDTO) (*complianceDTOs.ComplianceRegisterListResponseDTO, error)
	GetByID(id uint) (*models.ComplianceRegister, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.ComplianceRegister, error)
	Create(complianceRegister *models.ComplianceRegister, tx ...utility.DatabaseTransaction) (uint, error)
	Update(complianceRegister *models.ComplianceRegister, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error
	ForceDelete(id uint, tx ...utility.DatabaseTransaction) error
	GenerateReferenceNumber() (string, error)

	//lifecycle
	GetIDsByLifecycleStatus(ids []uint, status string) ([]uint, error)
	UpdateLifecycleStatusForIDs(ids []uint, statusId uint) ([]uint, error)

	//Document
	CreateDocument(doc *models.ComplianceDocument, tx ...utility.DatabaseTransaction) error
	GetDocumentByID(id uint) (*models.ComplianceDocument, error)
	DeleteDocuments(docIds []*uint, tx ...utility.DatabaseTransaction) error
	GetLatestDocumentByComplianceID(id uint, language string) (*models.ComplianceDocument, error)

	CreateInitialComplianceStatusUpdate(complianceStatusUpdate *models.ComplianceStatusUpdate, tx ...utility.DatabaseTransaction) (uint, error)
	DeleteComplianceStatusUpdates(ids []*uint, tx ...utility.DatabaseTransaction) error

	UpdateComplianceTotalScore(complianceRegisterID uint, totalScore float64, tx utility.DatabaseTransaction) error
	HasComplianceAssessments(registerID uint, tx ...utility.DatabaseTransaction) (bool, error)
	GetAllVersions(parentID uint) ([]*models.ComplianceRegister, error)

	// Issuing Authorities
	CreateIssuingAuthority(authority *models.ComplianceRegisterIssuingAuthority, tx ...utility.DatabaseTransaction) (uint, error)
	UpdateIssuingAuthority(authority *models.ComplianceRegisterIssuingAuthority, tx ...utility.DatabaseTransaction) error
	DeleteIssuingAuthority(id uint, tx ...utility.DatabaseTransaction) error
	GetIssuingAuthorityByID(id uint, tx ...utility.DatabaseTransaction) (*models.ComplianceRegisterIssuingAuthority, error)
	GetIssuingAuthoritiesByComplianceRegisterID(complianceRegisterID uint, tx ...utility.DatabaseTransaction) ([]*models.ComplianceRegisterIssuingAuthority, error)
	CreateIssuingAuthorityEntity(entity *models.ComplianceRegisterIssuingAuthorityEntity, tx ...utility.DatabaseTransaction) (uint, error)
	UpdateIssuingAuthorityEntity(entity *models.ComplianceRegisterIssuingAuthorityEntity, tx ...utility.DatabaseTransaction) error
	DeleteIssuingAuthorityEntity(id uint, tx ...utility.DatabaseTransaction) error
	GetIssuingAuthorityEntityByID(id uint, tx ...utility.DatabaseTransaction) (*models.ComplianceRegisterIssuingAuthorityEntity, error)
	CreateIssuingAuthorityDocument(document *models.ComplianceRegisterIssuingAuthorityDocument, tx ...utility.DatabaseTransaction) (uint, error)
	UpdateIssuingAuthorityDocument(document *models.ComplianceRegisterIssuingAuthorityDocument, tx ...utility.DatabaseTransaction) error
	DeleteIssuingAuthorityDocument(id uint, tx ...utility.DatabaseTransaction) error
	GetIssuingAuthorityDocumentByID(id uint, tx ...utility.DatabaseTransaction) (*models.ComplianceRegisterIssuingAuthorityDocument, error)
	GetDocumentsByAuthorityID(authorityID uint, tx ...utility.DatabaseTransaction) ([]*models.ComplianceRegisterIssuingAuthorityDocument, error)
	DeleteIssuingAuthorityDocumentsByAuthorityID(authorityID uint, tx ...utility.DatabaseTransaction) error
	DeleteIssuingAuthorityEntitiesByAuthorityID(authorityID uint, tx ...utility.DatabaseTransaction) error

	// Shared Units
	CreateSharedUnits(sharedUnits []*models.ComplianceRegisterSharedUnit, tx ...utility.DatabaseTransaction) error
	DeleteSharedUnits(complianceRegisterID uint, departmentIDs []uint, tx ...utility.DatabaseTransaction) error
	GetSharedUnitsByComplianceRegisterID(complianceRegisterID uint, tx ...utility.DatabaseTransaction) ([]*models.ComplianceRegisterSharedUnit, error)

	// Bulk Import
	GetAllComplianceTitles() ([]string, error)

	// Note and Observations
	UpdateNote(complianceRegisterID uint, note *string, tx ...utility.DatabaseTransaction) error
}

type ComplianceRegisterRepository struct {
	db *gorm.DB
}

func NewComplianceRegisterRepository(db *gorm.DB) *ComplianceRegisterRepository {
	return &ComplianceRegisterRepository{
		db: db,
	}
}

// BeginTransaction starts a new transaction.
func (repo *ComplianceRegisterRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *ComplianceRegisterRepository) GetAll(params *complianceDTOs.ComplianceRegisterListQueryParamsDTO) (*complianceDTOs.ComplianceRegisterListResponseDTO, error) {

	var complianceRegisters []*complianceDTOs.ComplianceRegisterListDTO
	var totalCount int64

	query := repo.db.Table("compliance_register").
		Select(`
            compliance_register.*,
			CONVERT(VARCHAR, compliance_register.deadline, 23) AS deadline,
			compliance_register.subject,
    		compliance_register.is_applicable,
            compliance_register.is_not_applicable,
            compliance_register.is_regulation_notification_assessment,
            compliance_register.is_risk_based_assessment,
            compliance_register.is_public,
            compliance_register.is_shared,
        	-- CONVERT(VARCHAR, compliance_register.issue_date, 23) AS issue_date,
        	-- CONVERT(VARCHAR, compliance_register.expiry_date, 23) AS expiry_date,
       		-- CONVERT(VARCHAR, compliance_register.next_review_date, 23) AS next_review_date,
			-- CONVERT(VARCHAR, compliance_register.last_reviewed_date, 23) AS last_reviewed_date,
			-- compliance_assessment_frameworks.id AS compliance_assessment_framework_id,
			-- compliance_assessment_frameworks.title_en AS compliance_assessment_framework_title_en,
			-- compliance_assessment_frameworks.title_ar AS compliance_assessment_framework_title_ar,
			-- compliance_categories.title_en AS compliance_category_title_en,
			-- compliance_categories.title_ar AS compliance_category_title_ar,
			compliance_statuses.title_en AS compliance_status_title_en,
			compliance_statuses.title_ar AS compliance_status_title_ar,
			compliance_statuses.color_code AS compliance_status_color_code,
			 organization_units.title_en AS responsible_unit_title_en,
			organization_units.title_ar AS responsible_unit_title_ar,
			responsible_user.id AS responsible_user_id,
			responsible_user.first_name_en AS responsible_user_first_name_en,
            responsible_user.first_name_ar AS responsible_user_first_name_ar,
            responsible_user.last_name_en AS responsible_user_last_name_en,
            responsible_user.last_name_ar AS responsible_user_last_name_ar,
			responsible_user.initial AS responsible_user_initial,
            responsible_user.color AS responsible_user_color,
            responsible_user.email AS responsible_user_email,
            responsible_user.profile_image_path AS responsible_user_profile_image_path,
			-- review_frequencies.title_en AS review_frequency_title_en,
			-- review_frequencies.title_ar AS review_frequency_title_ar,
			-- review_responsible_user.id AS review_responsible_user_id,
			-- review_responsible_user.first_name_en AS review_responsible_user_first_name_en,
            -- review_responsible_user.first_name_ar AS review_responsible_user_first_name_ar,
            -- review_responsible_user.last_name_en AS review_responsible_user_last_name_en,
            -- review_responsible_user.last_name_ar AS review_responsible_user_last_name_ar,
			-- review_responsible_user.initial AS review_responsible_user_initial,
        	-- review_responsible_user.color AS review_responsible_user_color,
        	-- review_responsible_user.email AS review_responsible_user_email,
        	-- review_responsible_user.profile_image_path AS review_responsible_user_profile_image_path,  
			-- severity_levels.title_en AS severity_level_title_en,
			-- severity_levels.title_ar AS severity_level_title_ar,
			-- severity_levels.color_code AS severity_level_color_code,
			-- lifecycle_statuses.title_en AS lifecycle_status_title_en,
			-- lifecycle_statuses.title_ar AS lifecycle_status_title_ar,
			-- lifecycle_statuses.color_code AS lifecycle_status_color_code,
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
        	updated_user.profile_image_path AS updated_user_profile_image_path,
        	
        	-- New compliance fields
        	legislative_bodies.id AS legislative_body_id,
        	legislative_bodies.title_en AS legislative_body_title_en,
        	legislative_bodies.title_ar AS legislative_body_title_ar,
        	legislation_types.id AS legislation_type_id,
        	legislation_types.title_en AS legislation_type_title_en,
        	legislation_types.title_ar AS legislation_type_title_ar,
        	importance_classifications.id AS importance_classification_id,
        	importance_classifications.title_en AS importance_classification_title_en,
        	importance_classifications.title_ar AS importance_classification_title_ar,
        	secrecy_classifications.id AS secrecy_classification_id,
        	secrecy_classifications.title_en AS secrecy_classification_title_en,
        	secrecy_classifications.title_ar AS secrecy_classification_title_ar,
        	average_residual_rating.id AS average_residual_risk_rating_id,
        	average_residual_rating.title_en AS average_residual_risk_rating_title_en,
        	average_residual_rating.title_ar AS average_residual_risk_rating_title_ar,
        	average_residual_rating.color_code AS average_residual_risk_rating_color_code,
        	average_residual_rating.type AS average_residual_risk_rating_type        
        `).
		Joins(`
		    -- LEFT JOIN compliance_categories ON compliance_categories.id = compliance_register.compliance_category_id
		    LEFT JOIN compliance_statuses ON compliance_statuses.id = compliance_register.compliance_status_id
		 	LEFT JOIN organization_units ON organization_units.id = compliance_register.responsible_unit_id
		    LEFT JOIN users AS responsible_user ON responsible_user.id = compliance_register.responsible_user_id
		    -- LEFT JOIN review_frequencies ON review_frequencies.id = compliance_register.review_frequency_id
			-- LEFT JOIN users AS review_responsible_user ON review_responsible_user.id = compliance_register.review_responsible_user_id
		    -- LEFT JOIN severity_levels ON severity_levels.id = compliance_register.severity_level_id
			-- LEFT JOIN compliance_assessment_frameworks ON compliance_assessment_frameworks.id = compliance_register.compliance_assessment_framework_id
			-- LEFT JOIN lifecycle_statuses ON lifecycle_statuses.id = compliance_register.lifecycle_status_id
			LEFT JOIN users AS created_user ON created_user.id = compliance_register.created_by_id
			LEFT JOIN users AS updated_user ON updated_user.id = compliance_register.updated_by_id
			
			-- New compliance field joins
			LEFT JOIN compliance_legislative_bodies AS legislative_bodies ON legislative_bodies.id = compliance_register.legislative_body_id
			LEFT JOIN compliance_legislation_types AS legislation_types ON legislation_types.id = compliance_register.legislation_type_id
			LEFT JOIN compliance_importance_classifications AS importance_classifications ON importance_classifications.id = compliance_register.importance_classification_id
			LEFT JOIN compliance_secrecy_classifications AS secrecy_classifications ON secrecy_classifications.id = compliance_register.secrecy_classification_id
			LEFT JOIN compliance_assessment_ratings AS average_residual_rating ON average_residual_rating.id = compliance_register.average_residual_risk_rating_id
		`)

	// Apply dynamic filters to the query
	for columnName, values := range params.Filters {
		if len(values) > 0 {
			if columnName == "average_residual_risk_rating_id" {
				query = query.Where("compliance_register.average_residual_risk_rating_id IN (?)", values)
			} else {
				query = query.Where(columnName+" IN (?)", values)
			}
		}
	}

	// Apply deadline filters
	if params.DeadlineFrom != "" {
		query = query.Where("deadline >= ?", params.DeadlineFrom)
	}
	if params.DeadlineTo != "" {
		query = query.Where("deadline <= ?", params.DeadlineTo)
	}

	// Apply issue_date filters
	if params.IssueDateFrom != "" {
		query = query.Where("issue_date >= ?", params.IssueDateFrom)
	}
	if params.IssueDateTo != "" {
		query = query.Where("issue_date <= ?", params.IssueDateTo)
	}

	// Apply expiry_date filters
	if params.IssueDateFrom != "" {
		query = query.Where("expiry_date >= ?", params.IssueDateFrom)
	}
	if params.ExpiryDateTo != "" {
		query = query.Where("expiry_date <= ?", params.ExpiryDateTo)
	}

	// Apply next_review_date filters
	if params.NextReviewDateFrom != "" {
		query = query.Where("next_review_date >= ?", params.NextReviewDateFrom)
	}
	if params.NextReviewDateTo != "" {
		query = query.Where("next_review_date <= ?", params.NextReviewDateTo)
	}

	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("compliance_register.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}
	if params.IsControlAssessment != nil {
		query = query.Where("compliance_register.is_control_assessment = ?", *params.IsControlAssessment)
	}

	if params.IsRequirementAssessment != nil {
		query = query.Where("compliance_register.is_requirement_assessment = ?", *params.IsRequirementAssessment)
	}

	if params.IsLatest != nil {
		query = query.Where("compliance_register.is_latest = ?", true)
	}

	if params.IsParent != nil {
		query = query.Where("compliance_register.compliance_register_id IS NULL")
	}

	// Apply quick filters
	for _, filter := range params.QuickFilter {
		switch filter {
		case "compliant":
			query = query.Where("compliance_statuses.type = ?", "compliant")
		case "partially_compliant":
			query = query.Where("compliance_statuses.type = ?", "partially_compliant")
		case "non_compliant":
			query = query.Where("compliance_statuses.type = ?", "not_implemented")
		case "critical":
			query = query.Where("average_residual_rating.type = ?", "critical")
		case "high":
			query = query.Where("average_residual_rating.type = ?", "high")
			// case "archived":
			//	query = query.Where("lifecycle_statuses.type = ?", "archived")
		}
	}

	// Calculate total count for pagination
	if err := query.Count(&totalCount).Error; err != nil {
		// ...
		return nil, err
	}

	// Apply search keyword if exists
	if params.Keyword != "" {
		query = query.Where(`
			(compliance_register.title_en LIKE ? OR
			compliance_register.title_ar LIKE ? OR
			compliance_register.reference_number LIKE ? OR
			compliance_register.description_en LIKE ? OR
			compliance_register.description_ar LIKE ? OR
			compliance_register.standard_code LIKE ? OR
			compliance_register.legislation_no LIKE ? OR
			compliance_register.subject LIKE ? OR
			compliance_statuses.title_en LIKE ? OR
			compliance_statuses.title_ar LIKE ? OR
			legislative_bodies.title_en LIKE ? OR
			legislative_bodies.title_ar LIKE ? OR
			legislation_types.title_en LIKE ? OR
			legislation_types.title_ar LIKE ? OR
			importance_classifications.title_en LIKE ? OR
			importance_classifications.title_ar LIKE ? OR
			secrecy_classifications.title_en LIKE ? OR
			secrecy_classifications.title_ar LIKE ? OR
			created_user.first_name_en LIKE ? OR
			created_user.first_name_ar LIKE ? OR
			created_user.last_name_en LIKE ? OR
			created_user.last_name_ar LIKE ? OR
			created_user.email LIKE ? OR
			updated_user.first_name_en LIKE ? OR
			updated_user.first_name_ar LIKE ? OR
			updated_user.last_name_en LIKE ? OR
			updated_user.last_name_ar LIKE ? OR
			updated_user.email LIKE ?)`,
			"%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%",
			"%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%",
			"%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%",
			"%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%",
			"%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%",
			"%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%",
			"%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%")
	}

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	query = query.Order(params.SortBy + " " + params.Sort)

	if !params.Export {
		offset := (params.Page - 1) * params.PageSize
		query = query.Offset(offset).Limit(params.PageSize)
	}

	err := query.Find(&complianceRegisters).Error
	if err != nil {
		return nil, err
	}

	totalPages := 1
	if !params.Export {
		totalPages = int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))
	}

	response := &complianceDTOs.ComplianceRegisterListResponseDTO{
		Items: complianceRegisters,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}
	return response, nil
}

func (repo *ComplianceRegisterRepository) GetByID(id uint) (*models.ComplianceRegister, error) {
	var compliance *models.ComplianceRegister

	result := repo.db.First(&compliance, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return compliance, nil

}

func (repo *ComplianceRegisterRepository) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.ComplianceRegister, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var compliance *models.ComplianceRegister

	result := db.
		// Preload("ComplianceCategory").
		Preload("ComplianceStatus").
		Preload("ComplianceClauses").
		Preload("ComplianceClauses.ComplianceRequirements").
		Preload("ComplianceClauses.ComplianceRequirements.RequiredDocumentTitles").
		Preload("ResponsibleUnit").
		Preload("ResponsibleUnit.ParentUnit.ParentUnit.ParentUnit.ParentUnit.ParentUnit").
		Preload("ResponsibleUser").
		Preload("ResponsibleUser.OrganizationUnit1").
		Preload("AverageResidualRiskRating").
		// Preload("ReviewFrequency").
		// Preload("ComplianceReviewStatus").
		// Preload("LastReviewedUser").
		// Preload("ReviewResponsibleUser").
		// Preload("SeverityLevel").
		// Preload("LifecycleStatus").
		Preload("ComplianceFolder").
		Preload("ParentComplianceRegister").
		Preload("ChildComplianceRegister").
		// Preload("DocumentEn", func(db *gorm.DB) *gorm.DB {
		// 	return db.Where("language = ?", "en")
		// }).
		// Preload("DocumentAr", func(db *gorm.DB) *gorm.DB {
		// 	return db.Where("language = ?", "ar")
		// }).
		// Preload("Controls").
		// Preload("Controls.Control.ControlCategory").
		// Preload("Controls.Control.ControlType").
		Preload("ComplianceStatusUpdates").
		// Preload("ComplianceAssessmentFramework").
		// Preload("ComplianceAssessmentFramework.AssessmentFrameworkType").
		// Preload("ComplianceAssessmentFramework.Values").
		// Preload("ComplianceAssessmentFramework.Values.ComplianceStatus").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").

		// New preloads for compliance fields
		Preload("LegislativeBody").
		Preload("LegislationType").
		Preload("ImportanceClassification").
		Preload("SecrecyClassification").
		Preload("AverageResidualRiskRating").

		// New preloads for issuing authorities and shared units
		Preload("IssuingAuthorities").
		Preload("IssuingAuthorities.IssuingAuthority").
		Preload("IssuingAuthorities.Entity").
		Preload("IssuingAuthorities.Document").
		Preload("ComplianceRegisterSharedUnit").
		Preload("ComplianceRegisterSharedUnit.Department").
		Preload("SubmittedByUser").
		Preload("ApprovedByUser").
		Preload("WorkflowStatus").
		First(&compliance, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return compliance, nil
}

func (repo *ComplianceRegisterRepository) Create(complianceRegister *models.ComplianceRegister, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(complianceRegister).Error; err != nil {
		return 0, err
	}

	return complianceRegister.ID, nil
}

func (repo *ComplianceRegisterRepository) Update(complianceRegister *models.ComplianceRegister, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(complianceRegister).Error; err != nil {
		return err
	}

	return nil
}

func (repo *ComplianceRegisterRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	// Delete the SLA Register
	if err := db.Delete(&models.ComplianceRegister{}, id).Error; err != nil {
		return err
	}

	return nil
}

func (repo *ComplianceRegisterRepository) GetIDsByLifecycleStatus(ids []uint, status string) ([]uint, error) {
	var matchingIDs []uint

	if err := repo.db.Model(&models.ComplianceRegister{}).
		Joins("JOIN lifecycle_statuses ON compliance_register.lifecycle_status_id = lifecycle_statuses.id").
		Where("compliance_register.id IN ? AND lifecycle_statuses.type = ?", ids, status).
		Pluck("compliance_register.id", &matchingIDs).Error; err != nil {
		return nil, err
	}

	return matchingIDs, nil
}

func (repo *ComplianceRegisterRepository) UpdateLifecycleStatusForIDs(ids []uint, statusId uint) ([]uint, error) {

	if err := repo.db.Model(&models.ComplianceRegister{}).
		Where("id IN ?", ids).
		Update("lifecycle_status_id", statusId).Error; err != nil {
		return nil, err
	}

	return ids, nil
}

func (repo *ComplianceRegisterRepository) GenerateReferenceNumber() (string, error) {
	var lastReferenceNumber string
	err := repo.db.Model(&models.ComplianceRegister{}).
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
	newReferenceNumber := fmt.Sprintf("%s-%d-%04d", "CR", currentYear, lastReferenceID+1)

	return newReferenceNumber, nil
}

func (repo *ComplianceRegisterRepository) GetDocumentByID(id uint) (*models.ComplianceDocument, error) {
	var document models.ComplianceDocument
	if err := repo.db.Where("id = ?", id).First(&document).Error; err != nil {
		return nil, err
	}
	return &document, nil
}

func (repo *ComplianceRegisterRepository) CreateDocument(doc *models.ComplianceDocument, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(doc).Error; err != nil {
		return err
	}
	return nil
}

func (repo *ComplianceRegisterRepository) DeleteDocuments(docIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	// Delete Compliance documents by IDs
	if err := db.Where("id IN (?)", docIds).Delete(&models.ComplianceDocument{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *ComplianceRegisterRepository) GetLatestDocumentByComplianceID(id uint, language string) (*models.ComplianceDocument, error) {
	var document models.ComplianceDocument
	if err := repo.db.Where("compliance_register_id = ? AND language = ?", id, language).First(&document).Error; err != nil {
		return nil, err
	}
	return &document, nil
}

func (repo *ComplianceRegisterRepository) CreateInitialComplianceStatusUpdate(complianceStatusUpdate *models.ComplianceStatusUpdate, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(complianceStatusUpdate).Error; err != nil {
		return 0, err
	}
	return complianceStatusUpdate.ID, nil
}

func (repo *ComplianceRegisterRepository) DeleteComplianceStatusUpdates(ids []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("id IN (?)", ids).Delete(&models.ComplianceStatusUpdate{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *ComplianceRegisterRepository) UpdateComplianceTotalScore(complianceRegisterID uint, totalScore float64, tx utility.DatabaseTransaction) error {
	db := tx.DB()

	err := db.Model(&models.ComplianceRegister{}).
		Where("id = ?", complianceRegisterID).
		Update("total_score", totalScore).
		Error
	if err != nil {
		return err
	}

	return nil
}
func (repo *ComplianceRegisterRepository) HasComplianceAssessments(registerID uint, tx ...utility.DatabaseTransaction) (bool, error) {
	var count int64

	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	query := db.Table("compliance_assessments").Where("compliance_register_id = ?", registerID)
	if err := query.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (repo *ComplianceRegisterRepository) GetAllVersions(parentID uint) ([]*models.ComplianceRegister, error) {
	var children []*models.ComplianceRegister
	err := repo.db.Preload("ComplianceCategory").
		Preload("ComplianceStatus").
		Preload("ResponsibleUnit").
		Preload("ResponsibleUser").
		Preload("ReviewFrequency").
		Preload("ReviewResponsibleUser").
		Preload("SeverityLevel").
		Preload("LifecycleStatus").
		Preload("LastReviewedUser").
		Preload("ComplianceReviewStatus").
		Preload("ComplianceFolder").
		Preload("ComplianceAssessmentFramework").
		Preload("DocumentEn").
		Preload("DocumentAr").
		Preload("Controls").
		Preload("ComplianceStatusUpdates").
		Preload("ComplianceClauses").
		Preload("ParentComplianceRegister").
		Where("compliance_register_id = ?", parentID).
		Find(&children).Error

	if err != nil {
		return nil, err
	}

	// Recursively fetch sub-children
	for _, child := range children {
		subChildren, err := repo.GetAllVersions(child.ID)
		if err != nil {
			return nil, err
		}
		child.ChildComplianceRegister = subChildren
	}

	return children, nil
}

func (repo *ComplianceRegisterRepository) ForceDelete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	/* ================= ACTION PLAN TREE ================= */
	var registerIDs []uint

	err := db.Raw(`
    WITH RecursiveRegisters AS (
        SELECT id, compliance_register_id
        FROM compliance_register
        WHERE id = ?

        UNION ALL

        SELECT cr.id, cr.compliance_register_id
        FROM compliance_register cr
        JOIN RecursiveRegisters rr
          ON cr.compliance_register_id = rr.id
    )
    SELECT id
    FROM RecursiveRegisters
    ORDER BY
        CASE WHEN compliance_register_id IS NULL THEN 1 ELSE 0 END,
        id DESC
`, id).Scan(&registerIDs).Error

	if err != nil {
		return err
	}

	queries := []string{

		/* ========= STATUS UPDATES ========= */

		`DELETE compliance_status_update_documents
		 FROM compliance_status_update_documents
		 JOIN compliance_status_updates ON compliance_status_updates.id = compliance_status_update_documents.compliance_status_update_id
		 WHERE compliance_status_updates.compliance_register_id = ?`,

		`DELETE FROM compliance_status_updates
		 WHERE compliance_register_id = ?`,

		/* ========= REVIEWS ========= */

		`DELETE compliance_review_documents
		 FROM compliance_review_documents
		 JOIN compliance_reviews
		   ON compliance_reviews.id = compliance_review_documents.compliance_review_id
		 WHERE compliance_reviews.compliance_register_id = ?`,

		`DELETE FROM compliance_reviews
		 WHERE compliance_register_id = ?`,

		/* ========= REPORTS ========= */

		`DELETE compliance_report_action_plans
		 FROM compliance_report_action_plans
		 JOIN compliance_report_requirements ON compliance_report_requirements.id = compliance_report_action_plans.compliance_report_requirement_id
		 JOIN compliance_report_clauses  ON compliance_report_clauses.id = compliance_report_requirements.compliance_report_clause_id
		 JOIN compliance_report_departments ON compliance_report_departments.id = compliance_report_clauses.compliance_report_department_id
		 JOIN compliance_reports ON compliance_reports.id = compliance_report_departments.compliance_report_id
		 WHERE compliance_reports.compliance_register_id = ?`,

		`DELETE compliance_report_requirements
		 FROM compliance_report_requirements
		 JOIN compliance_report_clauses ON compliance_report_clauses.id = compliance_report_requirements.compliance_report_clause_id
		 JOIN compliance_report_departments ON compliance_report_departments.id = compliance_report_clauses.compliance_report_department_id
		 JOIN compliance_reports ON compliance_reports.id = compliance_report_departments.compliance_report_id
		 WHERE compliance_reports.compliance_register_id = ?`,

		`DELETE compliance_report_clauses
		 FROM compliance_report_clauses
		 JOIN compliance_report_departments ON compliance_report_departments.id = compliance_report_clauses.compliance_report_department_id
		 JOIN compliance_reports ON compliance_reports.id = compliance_report_departments.compliance_report_id
		 WHERE compliance_reports.compliance_register_id = ?`,

		`DELETE FROM compliance_report_departments
		 WHERE compliance_report_id IN (
		   SELECT id FROM compliance_reports WHERE compliance_register_id = ?
		 )`,

		`DELETE FROM compliance_report_sections
		 WHERE compliance_report_id IN (
		   SELECT id FROM compliance_reports WHERE compliance_register_id = ?
		 )`,

		`DELETE FROM compliance_reports
		 WHERE compliance_register_id = ?`,

		/* ========= DOCUMENTS ========= */
		`DELETE FROM compliance_documents
		 WHERE compliance_register_id = ?`,

		/* ========= CLAUSES ========= */
		`DELETE control_clauses
		 FROM control_clauses
		 JOIN compliance_clauses
		   ON compliance_clauses.id = control_clauses.compliance_clause_id
		 WHERE compliance_clauses.compliance_register_id = ?`,

		`DELETE compliance_requirement_status_update_documents
		 FROM compliance_requirement_status_update_documents
		 JOIN compliance_requirement_status_updates ON compliance_requirement_status_updates.id = compliance_requirement_status_update_documents.compliance_requirement_status_update_id
		 JOIN compliance_requirements ON compliance_requirements.id = compliance_requirement_status_updates.compliance_requirement_id
		 JOIN compliance_clauses ON compliance_clauses.id = compliance_requirements.compliance_clause_id
		 WHERE compliance_clauses.compliance_register_id = ?`,

		`DELETE compliance_requirement_status_updates
		 FROM compliance_requirement_status_updates
		 JOIN compliance_requirements ON compliance_requirements.id =compliance_requirement_status_updates.compliance_requirement_id
		 JOIN compliance_clauses ON compliance_clauses.id = compliance_requirements.compliance_clause_id
		 WHERE compliance_clauses.compliance_register_id = ?`,

		`DELETE compliance_assessment_control_comments
		FROM compliance_assessment_control_comments
		JOIN compliance_assessment_controls ON compliance_assessment_controls.id =compliance_assessment_control_comments.compliance_assessment_control_id
		JOIN compliance_controls ON compliance_controls.id = compliance_assessment_controls.compliance_control_id
		JOIN compliance_requirements ON compliance_requirements.id =compliance_controls.compliance_requirement_id
		JOIN compliance_clauses ON compliance_clauses.id = compliance_requirements.compliance_clause_id
		WHERE compliance_clauses.compliance_register_id = ?`,

		`DELETE compliance_action_plan_status_update_documents
		FROM compliance_action_plan_status_update_documents
		JOIN compliance_action_plan_status_updates ON compliance_action_plan_status_updates.id = compliance_action_plan_status_update_documents.compliance_action_plan_status_update_id
		JOIN compliance_action_plans ON compliance_action_plans.id = compliance_action_plan_status_updates.compliance_action_plan_id
		JOIN compliance_assessment_controls ON compliance_assessment_controls.id =compliance_action_plans.compliance_assessment_control_id
		JOIN compliance_controls ON compliance_controls.id = compliance_assessment_controls.compliance_control_id
		JOIN compliance_requirements ON compliance_requirements.id =compliance_controls.compliance_requirement_id
		JOIN compliance_clauses ON compliance_clauses.id = compliance_requirements.compliance_clause_id
		WHERE compliance_clauses.compliance_register_id = ?`,

		`DELETE compliance_action_plan_status_updates
		FROM compliance_action_plan_status_updates
		JOIN compliance_action_plans ON compliance_action_plans.id = compliance_action_plan_status_updates.compliance_action_plan_id
		JOIN compliance_assessment_controls ON compliance_assessment_controls.id =compliance_action_plans.compliance_assessment_control_id
		JOIN compliance_controls ON compliance_controls.id = compliance_assessment_controls.compliance_control_id
		JOIN compliance_requirements ON compliance_requirements.id =compliance_controls.compliance_requirement_id
		JOIN compliance_clauses ON compliance_clauses.id = compliance_requirements.compliance_clause_id
		WHERE compliance_clauses.compliance_register_id = ?`,

		`DELETE compliance_action_plan_comment_documents
		FROM compliance_action_plan_comment_documents
		JOIN compliance_action_plan_comments ON compliance_action_plan_comments.id = compliance_action_plan_comment_documents.compliance_action_plan_comment_id
		JOIN compliance_action_plans ON compliance_action_plans.id = compliance_action_plan_comments.compliance_action_plan_id
		JOIN compliance_assessment_controls ON compliance_assessment_controls.id =compliance_action_plans.compliance_assessment_control_id
		JOIN compliance_controls ON compliance_controls.id = compliance_assessment_controls.compliance_control_id
		JOIN compliance_requirements ON compliance_requirements.id =compliance_controls.compliance_requirement_id
		JOIN compliance_clauses ON compliance_clauses.id = compliance_requirements.compliance_clause_id
		WHERE compliance_clauses.compliance_register_id = ?`,

		`DELETE compliance_action_plan_comments
		FROM compliance_action_plan_comments
		JOIN compliance_action_plans ON compliance_action_plans.id = compliance_action_plan_comments.compliance_action_plan_id
		JOIN compliance_assessment_controls ON compliance_assessment_controls.id =compliance_action_plans.compliance_assessment_control_id
		JOIN compliance_controls ON compliance_controls.id = compliance_assessment_controls.compliance_control_id
		JOIN compliance_requirements ON compliance_requirements.id =compliance_controls.compliance_requirement_id
		JOIN compliance_clauses ON compliance_clauses.id = compliance_requirements.compliance_clause_id
		WHERE compliance_clauses.compliance_register_id = ?`,

		`DELETE compliance_action_plans
		FROM compliance_action_plans
		JOIN compliance_assessment_controls ON compliance_assessment_controls.id = compliance_action_plans.compliance_assessment_control_id
		JOIN compliance_controls ON compliance_controls.id = compliance_assessment_controls.compliance_control_id
		JOIN compliance_requirements ON compliance_requirements.id =compliance_controls.compliance_requirement_id
		JOIN compliance_clauses ON compliance_clauses.id = compliance_requirements.compliance_clause_id
		WHERE compliance_clauses.compliance_register_id = ?`,

		`DELETE compliance_assessment_controls
		FROM compliance_assessment_controls 
		JOIN compliance_controls ON compliance_controls.id = compliance_assessment_controls.compliance_control_id
		JOIN compliance_requirements ON compliance_requirements.id =compliance_controls.compliance_requirement_id
		JOIN compliance_clauses ON compliance_clauses.id = compliance_requirements.compliance_clause_id
		WHERE compliance_clauses.compliance_register_id = ?`,

		`DELETE compliance_controls
		FROM compliance_controls 
		JOIN compliance_requirements ON compliance_requirements.id =compliance_controls.compliance_requirement_id
		JOIN compliance_clauses ON compliance_clauses.id = compliance_requirements.compliance_clause_id
		WHERE compliance_clauses.compliance_register_id = ?`,

		`DELETE compliance_assessment_requirement_required_documents
		 FROM compliance_assessment_requirement_required_documents
		 JOIN compliance_assessment_requirement_document_titles ON compliance_assessment_requirement_document_titles.id = compliance_assessment_requirement_required_documents.compliance_assessment_requirement_document_title_id
		 JOIN compliance_assessment_requirement_responsible_units ON compliance_assessment_requirement_responsible_units.id = compliance_assessment_requirement_document_titles.compliance_assessment_requirement_responsible_unit_id
		 JOIN compliance_assessment_requirements ON compliance_assessment_requirements.id = compliance_assessment_requirement_responsible_units.compliance_assessment_requirement_id
		 JOIN compliance_requirements ON compliance_requirements.id = compliance_assessment_requirements.compliance_requirement_id
		 JOIN compliance_clauses ON compliance_clauses.id = compliance_requirements.compliance_clause_id
		 WHERE compliance_clauses.compliance_register_id = ?`,

		`DELETE compliance_assessment_requirement_document_titles
		 FROM compliance_assessment_requirement_document_titles
		 JOIN compliance_assessment_requirement_responsible_units ON compliance_assessment_requirement_responsible_units.id = compliance_assessment_requirement_document_titles.compliance_assessment_requirement_responsible_unit_id
		 JOIN compliance_assessment_requirements ON compliance_assessment_requirements.id = compliance_assessment_requirement_responsible_units.compliance_assessment_requirement_id
		 JOIN compliance_requirements ON compliance_requirements.id = compliance_assessment_requirements.compliance_requirement_id
		 JOIN compliance_clauses ON compliance_clauses.id = compliance_requirements.compliance_clause_id
		 WHERE compliance_clauses.compliance_register_id = ?`,

		`DELETE compliance_assessment_requirement_comments
		 FROM compliance_assessment_requirement_comments 
		 JOIN compliance_assessment_requirement_responsible_units ON compliance_assessment_requirement_responsible_units.id = compliance_assessment_requirement_comments.compliance_assessment_requirement_responsible_unit_id
		 JOIN compliance_assessment_requirements ON compliance_assessment_requirements.id = compliance_assessment_requirement_responsible_units.compliance_assessment_requirement_id
		 JOIN compliance_requirements ON compliance_requirements.id = compliance_assessment_requirements.compliance_requirement_id
		 JOIN compliance_clauses ON compliance_clauses.id = compliance_requirements.compliance_clause_id
		 WHERE compliance_clauses.compliance_register_id = ?`,

		`DELETE compliance_assessment_requirement_action_plans
		 FROM compliance_assessment_requirement_action_plans 
		 JOIN compliance_assessment_requirement_responsible_units ON compliance_assessment_requirement_responsible_units.id = compliance_assessment_requirement_action_plans.compliance_assessment_requirement_responsible_unit_id
		 JOIN compliance_assessment_requirements ON compliance_assessment_requirements.id = compliance_assessment_requirement_responsible_units.compliance_assessment_requirement_id
		 JOIN compliance_requirements ON compliance_requirements.id = compliance_assessment_requirements.compliance_requirement_id
		 JOIN compliance_clauses ON compliance_clauses.id = compliance_requirements.compliance_clause_id
		 WHERE compliance_clauses.compliance_register_id = ?`,

		`DELETE compliance_action_plan_status_update_documents
		FROM compliance_action_plan_status_update_documents
		JOIN compliance_action_plan_status_updates
		ON compliance_action_plan_status_updates.id =
			compliance_action_plan_status_update_documents.compliance_action_plan_status_update_id
		JOIN compliance_action_plans
		ON compliance_action_plans.id =
			compliance_action_plan_status_updates.compliance_action_plan_id
		JOIN compliance_assessment_requirement_responsible_units
		ON compliance_assessment_requirement_responsible_units.id =
			compliance_action_plans.compliance_assessment_requirement_responsible_unit_id
		JOIN compliance_assessment_requirements
		ON compliance_assessment_requirements.id =
			compliance_assessment_requirement_responsible_units.compliance_assessment_requirement_id
		JOIN compliance_requirements
		ON compliance_requirements.id =
			compliance_assessment_requirements.compliance_requirement_id
		JOIN compliance_clauses
		ON compliance_clauses.id =
			compliance_requirements.compliance_clause_id
		WHERE compliance_clauses.compliance_register_id = ?`,

		`DELETE compliance_action_plan_status_updates
		FROM compliance_action_plan_status_updates
		JOIN compliance_action_plans
		ON compliance_action_plans.id = compliance_action_plan_status_updates.compliance_action_plan_id
		JOIN compliance_assessment_requirement_responsible_units
		ON compliance_assessment_requirement_responsible_units.id =
			compliance_action_plans.compliance_assessment_requirement_responsible_unit_id
		JOIN compliance_assessment_requirements
		ON compliance_assessment_requirements.id =
			compliance_assessment_requirement_responsible_units.compliance_assessment_requirement_id
		JOIN compliance_requirements
		ON compliance_requirements.id =
			compliance_assessment_requirements.compliance_requirement_id
		JOIN compliance_clauses
		ON compliance_clauses.id =
			compliance_requirements.compliance_clause_id
		WHERE compliance_clauses.compliance_register_id = ?`,

		`DELETE tasks
		FROM tasks
		JOIN compliance_action_plans
		ON compliance_action_plans.id = tasks.compliance_action_plan_id
		JOIN compliance_assessment_requirement_responsible_units
		ON compliance_assessment_requirement_responsible_units.id =
			compliance_action_plans.compliance_assessment_requirement_responsible_unit_id
		JOIN compliance_assessment_requirements
		ON compliance_assessment_requirements.id =
			compliance_assessment_requirement_responsible_units.compliance_assessment_requirement_id
		JOIN compliance_requirements
		ON compliance_requirements.id =
			compliance_assessment_requirements.compliance_requirement_id
		JOIN compliance_clauses
		ON compliance_clauses.id =
			compliance_requirements.compliance_clause_id
		WHERE compliance_clauses.compliance_register_id = ?`,

		`DELETE compliance_action_plans
		FROM compliance_action_plans
		JOIN compliance_assessment_requirement_responsible_units
		ON compliance_assessment_requirement_responsible_units.id =
			compliance_action_plans.compliance_assessment_requirement_responsible_unit_id
		JOIN compliance_assessment_requirements
		ON compliance_assessment_requirements.id =
			compliance_assessment_requirement_responsible_units.compliance_assessment_requirement_id
		JOIN compliance_requirements
		ON compliance_requirements.id =
			compliance_assessment_requirements.compliance_requirement_id
		JOIN compliance_clauses
		ON compliance_clauses.id =
			compliance_requirements.compliance_clause_id
		WHERE compliance_clauses.compliance_register_id = ?`,

		`DELETE compliance_assessment_requirement_responsible_units
		 FROM compliance_assessment_requirement_responsible_units 
		 JOIN compliance_assessment_requirements ON compliance_assessment_requirements.id = compliance_assessment_requirement_responsible_units.compliance_assessment_requirement_id
		 JOIN compliance_requirements ON compliance_requirements.id = compliance_assessment_requirements.compliance_requirement_id
		 JOIN compliance_clauses ON compliance_clauses.id = compliance_requirements.compliance_clause_id
		 WHERE compliance_clauses.compliance_register_id = ?`,

		`DELETE compliance_assessment_requirements
		 FROM compliance_assessment_requirements 
		 JOIN compliance_requirements ON compliance_requirements.id = compliance_assessment_requirements.compliance_requirement_id
		 JOIN compliance_clauses ON compliance_clauses.id = compliance_requirements.compliance_clause_id
		 WHERE compliance_clauses.compliance_register_id = ?`,

		`DELETE compliance_requirement_required_document_titles
		 FROM compliance_requirement_required_document_titles
		 JOIN compliance_requirements ON compliance_requirements.id = compliance_requirement_required_document_titles.compliance_requirement_id
		 JOIN compliance_clauses ON compliance_clauses.id = compliance_requirements.compliance_clause_id		 
		 WHERE compliance_clauses.compliance_register_id = ?`,

		`DELETE compliance_requirement_responsible_units
		FROM compliance_requirement_responsible_units
		JOIN compliance_requirements
		ON compliance_requirements.id =
			compliance_requirement_responsible_units.compliance_requirement_id
		JOIN compliance_clauses
		ON compliance_clauses.id =
			compliance_requirements.compliance_clause_id
		WHERE compliance_clauses.compliance_register_id = ?`,

		`DELETE audit_schedule_auditable_items
		FROM audit_schedule_auditable_items
		JOIN auditable_items
		ON auditable_items.id = audit_schedule_auditable_items.auditable_item_id
		JOIN compliance_requirements
		ON compliance_requirements.id = auditable_items.compliance_requirement_id
		JOIN compliance_clauses
		ON compliance_clauses.id = compliance_requirements.compliance_clause_id
		WHERE compliance_clauses.compliance_register_id = ?`,

		`DELETE audit_plan_schedule_auditable_items
		FROM audit_plan_schedule_auditable_items
		JOIN auditable_items
		ON auditable_items.id = audit_plan_schedule_auditable_items.auditable_item_id
		JOIN compliance_requirements
		ON compliance_requirements.id = auditable_items.compliance_requirement_id
		JOIN compliance_clauses
		ON compliance_clauses.id = compliance_requirements.compliance_clause_id
		WHERE compliance_clauses.compliance_register_id = ?`,

		`DELETE audit_document_request_auditable_items
		FROM audit_document_request_auditable_items
		JOIN auditable_items
		ON auditable_items.id = audit_document_request_auditable_items.auditable_item_id
		JOIN compliance_requirements
		ON compliance_requirements.id = auditable_items.compliance_requirement_id
		JOIN compliance_clauses
		ON compliance_clauses.id = compliance_requirements.compliance_clause_id
		WHERE compliance_clauses.compliance_register_id = ?`,

		`DELETE audit_finding_strategic_objectives
		FROM audit_finding_strategic_objectives
		JOIN audit_findings
		ON audit_findings.id = audit_finding_strategic_objectives.audit_finding_id
		JOIN audit_checklists
		ON audit_checklists.id = audit_findings.audit_checklist_id
		JOIN audit_plan_checklists
		ON audit_plan_checklists.id = audit_checklists.audit_plan_checklist_id
		JOIN auditable_items
		ON auditable_items.id = audit_plan_checklists.auditable_item_id
		JOIN compliance_requirements
		ON compliance_requirements.id = auditable_items.compliance_requirement_id
		JOIN compliance_clauses
		ON compliance_clauses.id = compliance_requirements.compliance_clause_id
		WHERE compliance_clauses.compliance_register_id = ?`,

		`DELETE audit_finding_root_cause_analyses
		FROM audit_finding_root_cause_analyses
		JOIN audit_findings
		ON audit_findings.id = audit_finding_root_cause_analyses.audit_finding_id
		JOIN audit_checklists
		ON audit_checklists.id = audit_findings.audit_checklist_id
		JOIN audit_plan_checklists
		ON audit_plan_checklists.id = audit_checklists.audit_plan_checklist_id
		JOIN auditable_items
		ON auditable_items.id = audit_plan_checklists.auditable_item_id
		JOIN compliance_requirements
		ON compliance_requirements.id = auditable_items.compliance_requirement_id
		JOIN compliance_clauses
		ON compliance_clauses.id = compliance_requirements.compliance_clause_id
		WHERE compliance_clauses.compliance_register_id = ?`,

		`DELETE audit_finding_risks
		FROM audit_finding_risks
		JOIN audit_findings
		ON audit_findings.id = audit_finding_risks.audit_finding_id
		JOIN audit_checklists
		ON audit_checklists.id = audit_findings.audit_checklist_id
		JOIN audit_plan_checklists
		ON audit_plan_checklists.id = audit_checklists.audit_plan_checklist_id
		JOIN auditable_items
		ON auditable_items.id = audit_plan_checklists.auditable_item_id
		JOIN compliance_requirements
		ON compliance_requirements.id = auditable_items.compliance_requirement_id
		JOIN compliance_clauses
		ON compliance_clauses.id = compliance_requirements.compliance_clause_id
		WHERE compliance_clauses.compliance_register_id = ?`,

		`DELETE audit_finding_processes
		FROM audit_finding_processes
		JOIN audit_findings ON audit_findings.id = audit_finding_processes.audit_finding_id
		JOIN audit_checklists ON audit_checklists.id = audit_findings.audit_checklist_id
		JOIN audit_plan_checklists ON audit_plan_checklists.id = audit_checklists.audit_plan_checklist_id
		JOIN auditable_items ON auditable_items.id = audit_plan_checklists.auditable_item_id
		JOIN compliance_requirements ON compliance_requirements.id = auditable_items.compliance_requirement_id
		JOIN compliance_clauses ON compliance_clauses.id = compliance_requirements.compliance_clause_id
		WHERE compliance_clauses.compliance_register_id = ?`,

		`DELETE audit_finding_kpis
		FROM audit_finding_kpis
		JOIN audit_findings ON audit_findings.id = audit_finding_kpis.audit_finding_id
		JOIN audit_checklists ON audit_checklists.id = audit_findings.audit_checklist_id
		JOIN audit_plan_checklists ON audit_plan_checklists.id = audit_checklists.audit_plan_checklist_id
		JOIN auditable_items ON auditable_items.id = audit_plan_checklists.auditable_item_id
		JOIN compliance_requirements ON compliance_requirements.id = auditable_items.compliance_requirement_id
		JOIN compliance_clauses ON compliance_clauses.id = compliance_requirements.compliance_clause_id
		WHERE compliance_clauses.compliance_register_id = ?`,

		`DELETE audit_finding_incidents
		FROM audit_finding_incidents
		JOIN audit_findings ON audit_findings.id = audit_finding_incidents.audit_finding_id
		JOIN audit_checklists ON audit_checklists.id = audit_findings.audit_checklist_id
		JOIN audit_plan_checklists ON audit_plan_checklists.id = audit_checklists.audit_plan_checklist_id
		JOIN auditable_items ON auditable_items.id = audit_plan_checklists.auditable_item_id
		JOIN compliance_requirements ON compliance_requirements.id = auditable_items.compliance_requirement_id
		JOIN compliance_clauses ON compliance_clauses.id = compliance_requirements.compliance_clause_id
		WHERE compliance_clauses.compliance_register_id = ?`,

		`DELETE audit_finding_documents
		FROM audit_finding_documents
		JOIN audit_findings ON audit_findings.id = audit_finding_documents.audit_finding_id
		JOIN audit_checklists ON audit_checklists.id = audit_findings.audit_checklist_id
		JOIN audit_plan_checklists ON audit_plan_checklists.id = audit_checklists.audit_plan_checklist_id
		JOIN auditable_items ON auditable_items.id = audit_plan_checklists.auditable_item_id
		JOIN compliance_requirements ON compliance_requirements.id = auditable_items.compliance_requirement_id
		JOIN compliance_clauses ON compliance_clauses.id = compliance_requirements.compliance_clause_id
		WHERE compliance_clauses.compliance_register_id = ?`,

		`DELETE audit_finding_controls
		FROM audit_finding_controls
		JOIN audit_findings ON audit_findings.id = audit_finding_controls.audit_finding_id
		JOIN audit_checklists ON audit_checklists.id = audit_findings.audit_checklist_id
		JOIN audit_plan_checklists ON audit_plan_checklists.id = audit_checklists.audit_plan_checklist_id
		JOIN auditable_items ON auditable_items.id = audit_plan_checklists.auditable_item_id
		JOIN compliance_requirements ON compliance_requirements.id = auditable_items.compliance_requirement_id
		JOIN compliance_clauses ON compliance_clauses.id = compliance_requirements.compliance_clause_id
		WHERE compliance_clauses.compliance_register_id = ?`,

		`DELETE audit_finding_compliance_requirements
		FROM audit_finding_compliance_requirements
		JOIN audit_findings ON audit_findings.id = audit_finding_compliance_requirements.audit_finding_id
		JOIN audit_checklists ON audit_checklists.id = audit_findings.audit_checklist_id
		JOIN audit_plan_checklists ON audit_plan_checklists.id = audit_checklists.audit_plan_checklist_id
		JOIN auditable_items ON auditable_items.id = audit_plan_checklists.auditable_item_id
		JOIN compliance_requirements ON compliance_requirements.id = auditable_items.compliance_requirement_id
		JOIN compliance_clauses ON compliance_clauses.id = compliance_requirements.compliance_clause_id
		WHERE compliance_clauses.compliance_register_id = ?`,

		`DELETE audit_finding_associated_documents
		FROM audit_finding_associated_documents
		JOIN audit_findings ON audit_findings.id = audit_finding_associated_documents.audit_finding_id
		JOIN audit_checklists ON audit_checklists.id = audit_findings.audit_checklist_id
		JOIN audit_plan_checklists ON audit_plan_checklists.id = audit_checklists.audit_plan_checklist_id
		JOIN auditable_items ON auditable_items.id = audit_plan_checklists.auditable_item_id
		JOIN compliance_requirements ON compliance_requirements.id = auditable_items.compliance_requirement_id
		JOIN compliance_clauses ON compliance_clauses.id = compliance_requirements.compliance_clause_id
		WHERE compliance_clauses.compliance_register_id = ?`,

		`DELETE audit_finding_assets
		FROM audit_finding_assets
		JOIN audit_findings ON audit_findings.id = audit_finding_assets.audit_finding_id
		JOIN audit_checklists ON audit_checklists.id = audit_findings.audit_checklist_id
		JOIN audit_plan_checklists ON audit_plan_checklists.id = audit_checklists.audit_plan_checklist_id
		JOIN auditable_items ON auditable_items.id = audit_plan_checklists.auditable_item_id
		JOIN compliance_requirements ON compliance_requirements.id = auditable_items.compliance_requirement_id
		JOIN compliance_clauses ON compliance_clauses.id = compliance_requirements.compliance_clause_id
		WHERE compliance_clauses.compliance_register_id = ?`,

		`DELETE tasks
		FROM tasks
		JOIN audit_finding_action_plans ON audit_finding_action_plans.id = tasks.audit_finding_action_plan_id
		JOIN audit_findings ON audit_findings.id = audit_finding_action_plans.audit_finding_id
		JOIN audit_checklists ON audit_checklists.id = audit_findings.audit_checklist_id
		JOIN audit_plan_checklists ON audit_plan_checklists.id = audit_checklists.audit_plan_checklist_id
		JOIN auditable_items ON auditable_items.id = audit_plan_checklists.auditable_item_id
		JOIN compliance_requirements ON compliance_requirements.id = auditable_items.compliance_requirement_id
		JOIN compliance_clauses ON compliance_clauses.id = compliance_requirements.compliance_clause_id
		WHERE compliance_clauses.compliance_register_id = ?`,

		`DELETE audit_finding_action_plan_status_update_documents
		FROM audit_finding_action_plan_status_update_documents
		JOIN audit_finding_action_plan_status_updates
		ON audit_finding_action_plan_status_updates.id =
			audit_finding_action_plan_status_update_documents.audit_finding_action_plan_status_update_id
		JOIN audit_finding_action_plans
		ON audit_finding_action_plans.id =
			audit_finding_action_plan_status_updates.audit_finding_action_plan_id
		JOIN audit_findings ON audit_findings.id = audit_finding_action_plans.audit_finding_id
		JOIN audit_checklists ON audit_checklists.id = audit_findings.audit_checklist_id
		JOIN audit_plan_checklists ON audit_plan_checklists.id = audit_checklists.audit_plan_checklist_id
		JOIN auditable_items ON auditable_items.id = audit_plan_checklists.auditable_item_id
		JOIN compliance_requirements ON compliance_requirements.id = auditable_items.compliance_requirement_id
		JOIN compliance_clauses ON compliance_clauses.id = compliance_requirements.compliance_clause_id
		WHERE compliance_clauses.compliance_register_id = ?`,

		`DELETE audit_finding_action_plan_status_updates
		FROM audit_finding_action_plan_status_updates
		JOIN audit_finding_action_plans
		ON audit_finding_action_plans.id =
			audit_finding_action_plan_status_updates.audit_finding_action_plan_id
		JOIN audit_findings ON audit_findings.id = audit_finding_action_plans.audit_finding_id
		JOIN audit_checklists ON audit_checklists.id = audit_findings.audit_checklist_id
		JOIN audit_plan_checklists ON audit_plan_checklists.id = audit_checklists.audit_plan_checklist_id
		JOIN auditable_items ON auditable_items.id = audit_plan_checklists.auditable_item_id
		JOIN compliance_requirements ON compliance_requirements.id = auditable_items.compliance_requirement_id
		JOIN compliance_clauses ON compliance_clauses.id = compliance_requirements.compliance_clause_id
		WHERE compliance_clauses.compliance_register_id = ?`,

		`DELETE audit_finding_action_plan_comment_documents
		FROM audit_finding_action_plan_comment_documents
		JOIN audit_finding_action_plan_comments
		ON audit_finding_action_plan_comments.id =
			audit_finding_action_plan_comment_documents.audit_finding_action_plan_comment_id
		JOIN audit_finding_action_plans
		ON audit_finding_action_plans.id =
			audit_finding_action_plan_comments.audit_finding_action_plan_id
		JOIN audit_findings ON audit_findings.id = audit_finding_action_plans.audit_finding_id
		JOIN audit_checklists ON audit_checklists.id = audit_findings.audit_checklist_id
		JOIN audit_plan_checklists ON audit_plan_checklists.id = audit_checklists.audit_plan_checklist_id
		JOIN auditable_items ON auditable_items.id = audit_plan_checklists.auditable_item_id
		JOIN compliance_requirements ON compliance_requirements.id = auditable_items.compliance_requirement_id
		JOIN compliance_clauses ON compliance_clauses.id = compliance_requirements.compliance_clause_id
		WHERE compliance_clauses.compliance_register_id = ?`,

		`DELETE audit_finding_action_plan_comments
		FROM audit_finding_action_plan_comments
		JOIN audit_finding_action_plans ON audit_finding_action_plans.id = audit_finding_action_plan_comments.audit_finding_action_plan_id
		JOIN audit_findings ON audit_findings.id = audit_finding_action_plans.audit_finding_id
		JOIN audit_checklists ON audit_checklists.id = audit_findings.audit_checklist_id
		JOIN audit_plan_checklists ON audit_plan_checklists.id = audit_checklists.audit_plan_checklist_id
		JOIN auditable_items ON auditable_items.id = audit_plan_checklists.auditable_item_id
		JOIN compliance_requirements ON compliance_requirements.id = auditable_items.compliance_requirement_id
		JOIN compliance_clauses ON compliance_clauses.id = compliance_requirements.compliance_clause_id
		WHERE compliance_clauses.compliance_register_id = ?`,

		`DELETE audit_finding_action_plans
		FROM audit_finding_action_plans
		JOIN audit_findings ON audit_findings.id = audit_finding_action_plans.audit_finding_id
		JOIN audit_checklists ON audit_checklists.id = audit_findings.audit_checklist_id
		JOIN audit_plan_checklists ON audit_plan_checklists.id = audit_checklists.audit_plan_checklist_id
		JOIN auditable_items ON auditable_items.id = audit_plan_checklists.auditable_item_id
		JOIN compliance_requirements ON compliance_requirements.id = auditable_items.compliance_requirement_id
		JOIN compliance_clauses ON compliance_clauses.id = compliance_requirements.compliance_clause_id
		WHERE compliance_clauses.compliance_register_id = ?`,

		`DELETE audit_findings
		FROM audit_findings
		JOIN audit_checklists ON audit_checklists.id = audit_findings.audit_checklist_id
		JOIN audit_plan_checklists ON audit_plan_checklists.id = audit_checklists.audit_plan_checklist_id
		JOIN auditable_items ON auditable_items.id = audit_plan_checklists.auditable_item_id
		JOIN compliance_requirements ON compliance_requirements.id = auditable_items.compliance_requirement_id
		JOIN compliance_clauses ON compliance_clauses.id = compliance_requirements.compliance_clause_id
		WHERE compliance_clauses.compliance_register_id = ?`,

		`DELETE audit_draft_finding_documents
		 FROM audit_draft_finding_documents 
		 JOIN audit_draft_findings ON audit_draft_findings.id = audit_draft_finding_documents.audit_draft_finding_id
		 JOIN audit_checklists ON audit_checklists.id = audit_draft_findings.audit_checklist_id
		 JOIN audit_plan_checklists ON audit_plan_checklists.id = audit_checklists.audit_plan_checklist_id
		 JOIN auditable_items ON auditable_items.id = audit_plan_checklists.auditable_item_id
		 JOIN compliance_requirements ON compliance_requirements.id = auditable_items.compliance_requirement_id
		 JOIN compliance_clauses ON compliance_clauses.id = compliance_requirements.compliance_clause_id
		 WHERE compliance_clauses.compliance_register_id = ?`,

		`DELETE audit_draft_finding_conversation_documents
		 FROM audit_draft_finding_conversation_documents 
		 JOIN audit_draft_finding_conversations ON audit_draft_finding_conversations.id = audit_draft_finding_conversation_documents.audit_draft_finding_conversation_id
		 JOIN audit_draft_findings ON audit_draft_findings.id = audit_draft_finding_conversations.audit_draft_finding_id
		 JOIN audit_checklists ON audit_checklists.id = audit_draft_findings.audit_checklist_id
		 JOIN audit_plan_checklists ON audit_plan_checklists.id = audit_checklists.audit_plan_checklist_id
		 JOIN auditable_items ON auditable_items.id = audit_plan_checklists.auditable_item_id
		 JOIN compliance_requirements ON compliance_requirements.id = auditable_items.compliance_requirement_id
		 JOIN compliance_clauses ON compliance_clauses.id = compliance_requirements.compliance_clause_id
		 WHERE compliance_clauses.compliance_register_id = ?`,

		`DELETE audit_draft_finding_conversations
		 FROM audit_draft_finding_conversations 
		 JOIN audit_draft_findings ON audit_draft_findings.id = audit_draft_finding_conversations.audit_draft_finding_id
		 JOIN audit_checklists ON audit_checklists.id = audit_draft_findings.audit_checklist_id
		 JOIN audit_plan_checklists ON audit_plan_checklists.id = audit_checklists.audit_plan_checklist_id
		 JOIN auditable_items ON auditable_items.id = audit_plan_checklists.auditable_item_id
		 JOIN compliance_requirements ON compliance_requirements.id = auditable_items.compliance_requirement_id
		 JOIN compliance_clauses ON compliance_clauses.id = compliance_requirements.compliance_clause_id
		 WHERE compliance_clauses.compliance_register_id = ?`,

		`DELETE audit_draft_finding_action_plans
		 FROM audit_draft_finding_action_plans 
		 JOIN audit_draft_findings ON audit_draft_findings.id = audit_draft_finding_action_plans.audit_draft_finding_id
		 JOIN audit_checklists ON audit_checklists.id = audit_draft_findings.audit_checklist_id
		 JOIN audit_plan_checklists ON audit_plan_checklists.id = audit_checklists.audit_plan_checklist_id
		 JOIN auditable_items ON auditable_items.id = audit_plan_checklists.auditable_item_id
		 JOIN compliance_requirements ON compliance_requirements.id = auditable_items.compliance_requirement_id
		 JOIN compliance_clauses ON compliance_clauses.id = compliance_requirements.compliance_clause_id
		 WHERE compliance_clauses.compliance_register_id = ?`,

		`DELETE audit_draft_findings
		 FROM audit_draft_findings 
		 JOIN audit_checklists ON audit_checklists.id = audit_draft_findings.audit_checklist_id
		 JOIN audit_plan_checklists ON audit_plan_checklists.id = audit_checklists.audit_plan_checklist_id
		 JOIN auditable_items ON auditable_items.id = audit_plan_checklists.auditable_item_id
		 JOIN compliance_requirements ON compliance_requirements.id = auditable_items.compliance_requirement_id
		 JOIN compliance_clauses ON compliance_clauses.id = compliance_requirements.compliance_clause_id
		 WHERE compliance_clauses.compliance_register_id = ?`,

		`DELETE audit_checklist_documents
		 FROM audit_checklist_documents 
		 JOIN audit_checklists ON audit_checklists.id = audit_checklist_documents.audit_checklist_id
		 JOIN audit_plan_checklists ON audit_plan_checklists.id = audit_checklists.audit_plan_checklist_id
		 JOIN auditable_items ON auditable_items.id = audit_plan_checklists.auditable_item_id
		 JOIN compliance_requirements ON compliance_requirements.id = auditable_items.compliance_requirement_id
		 JOIN compliance_clauses ON compliance_clauses.id = compliance_requirements.compliance_clause_id
		 WHERE compliance_clauses.compliance_register_id = ?`,

		`DELETE audit_checklists
		 FROM audit_checklists 
		 JOIN audit_plan_checklists ON audit_plan_checklists.id = audit_checklists.audit_plan_checklist_id
		 JOIN auditable_items ON auditable_items.id = audit_plan_checklists.auditable_item_id
		 JOIN compliance_requirements ON compliance_requirements.id = auditable_items.compliance_requirement_id
		 JOIN compliance_clauses ON compliance_clauses.id = compliance_requirements.compliance_clause_id
		 WHERE compliance_clauses.compliance_register_id = ?`,

		`DELETE audit_plan_checklists
		 FROM audit_plan_checklists 
		 JOIN auditable_items ON auditable_items.id = audit_plan_checklists.auditable_item_id
		 JOIN compliance_requirements ON compliance_requirements.id = auditable_items.compliance_requirement_id
		 JOIN compliance_clauses ON compliance_clauses.id = compliance_requirements.compliance_clause_id
		 WHERE compliance_clauses.compliance_register_id = ?`,

		`DELETE audit_plan_auditable_item_risk_controls
		 FROM audit_plan_auditable_item_risk_controls 
		 JOIN audit_plan_auditable_item_risks ON audit_plan_auditable_item_risks.id = audit_plan_auditable_item_risk_controls.audit_plan_auditable_item_risk_id
		 JOIN audit_plan_auditable_items ON audit_plan_auditable_items.id = audit_plan_auditable_item_risks.audit_plan_auditable_item_id
		 JOIN auditable_items ON auditable_items.id = audit_plan_auditable_items.auditable_item_id
		 JOIN compliance_requirements ON compliance_requirements.id = auditable_items.compliance_requirement_id
		 JOIN compliance_clauses ON compliance_clauses.id = compliance_requirements.compliance_clause_id		 
		 WHERE compliance_clauses.compliance_register_id = ?`,

		`DELETE audit_plan_auditable_item_risks
		 FROM audit_plan_auditable_item_risks 
		 JOIN audit_plan_auditable_items ON audit_plan_auditable_items.id = audit_plan_auditable_item_risks.audit_plan_auditable_item_id
		 JOIN auditable_items ON auditable_items.id = audit_plan_auditable_items.auditable_item_id
		 JOIN compliance_requirements ON compliance_requirements.id = auditable_items.compliance_requirement_id
		 JOIN compliance_clauses ON compliance_clauses.id = compliance_requirements.compliance_clause_id		 
		 WHERE compliance_clauses.compliance_register_id = ?`,

		`DELETE audit_plan_auditable_items
		 FROM audit_plan_auditable_items 
		 JOIN auditable_items ON auditable_items.id = audit_plan_auditable_items.auditable_item_id
		 JOIN compliance_requirements ON compliance_requirements.id = auditable_items.compliance_requirement_id
		 JOIN compliance_clauses ON compliance_clauses.id = compliance_requirements.compliance_clause_id		 
		 WHERE compliance_clauses.compliance_register_id = ?`,

		`DELETE auditable_items
		 FROM auditable_items 
		 JOIN compliance_requirements ON compliance_requirements.id = auditable_items.compliance_requirement_id
		 JOIN compliance_clauses ON compliance_clauses.id = compliance_requirements.compliance_clause_id		 
		 WHERE compliance_clauses.compliance_register_id = ?`,

		`DELETE FROM asset_compliance_requirements WHERE compliance_register_id = ?`,
		`DELETE FROM third_party_compliance_requirements WHERE compliance_register_id = ?`,
		`DELETE FROM risk_compliance_requirements WHERE compliance_register_id = ?`,
		`DELETE FROM project_compliance_requirements WHERE compliance_register_id = ?`,
		`DELETE FROM process_compliance_requirements WHERE compliance_register_id = ?`,
		`DELETE FROM meeting_compliance_requirements WHERE compliance_register_id = ?`,
		`DELETE FROM kri_compliance_requirements WHERE compliance_register_id = ?`,
		`DELETE FROM kpi_compliances WHERE compliance_register_id = ?`,
		`DELETE FROM incident_compliance_requirements WHERE compliance_register_id = ?`,
		`DELETE FROM document_compliance_requirements WHERE compliance_register_id = ?`,
		`DELETE FROM compliance_requirement_strategic_objectives WHERE compliance_register_id = ?`,
		`DELETE FROM bcp_compliance_requirements WHERE compliance_register_id = ?`,
		`DELETE FROM compliance_requirement_associated_documents WHERE compliance_register_id = ?`,
		//////////////////

		`DELETE compliance_action_plan_status_update_documents
		FROM compliance_action_plan_status_update_documents
		JOIN compliance_action_plan_status_updates ON compliance_action_plan_status_updates.id = compliance_action_plan_status_update_documents.compliance_action_plan_status_update_id
		JOIN compliance_action_plans ON compliance_action_plans.id = compliance_action_plan_status_updates.compliance_action_plan_id
		JOIN compliance_assessment_controls ON compliance_assessment_controls.id =compliance_action_plans.compliance_assessment_control_id
		JOIN compliance_controls ON compliance_controls.id = compliance_assessment_controls.compliance_control_id
		WHERE compliance_controls.compliance_register_id = ?`,

		`DELETE compliance_action_plan_status_updates
		FROM compliance_action_plan_status_updates
		JOIN compliance_action_plans ON compliance_action_plans.id = compliance_action_plan_status_updates.compliance_action_plan_id
		JOIN compliance_assessment_controls ON compliance_assessment_controls.id =compliance_action_plans.compliance_assessment_control_id
		JOIN compliance_controls ON compliance_controls.id = compliance_assessment_controls.compliance_control_id
		WHERE compliance_controls.compliance_register_id = ?`,

		`DELETE tasks
		FROM tasks
		JOIN compliance_action_plans
		ON compliance_action_plans.id = tasks.compliance_action_plan_id
		JOIN compliance_assessment_controls
		ON compliance_assessment_controls.id =
			compliance_action_plans.compliance_assessment_control_id
		JOIN compliance_controls
		ON compliance_controls.id =
			compliance_assessment_controls.compliance_control_id
		WHERE compliance_controls.compliance_register_id = ?`,

		`DELETE compliance_assessment_controls
		FROM compliance_assessment_controls 
		JOIN compliance_controls ON compliance_controls.id = compliance_assessment_controls.compliance_control_id
		WHERE compliance_controls.compliance_register_id = ?`,

		`DELETE compliance_controls
		FROM compliance_controls 
		JOIN compliance_clauses ON compliance_clauses.id = compliance_controls.compliance_clause_id
		WHERE compliance_clauses.compliance_register_id = ?`,

		`DELETE compliance_action_plan_status_update_documents
		FROM compliance_action_plan_status_update_documents
		JOIN compliance_action_plan_status_updates
		ON compliance_action_plan_status_updates.id =
			compliance_action_plan_status_update_documents.compliance_action_plan_status_update_id
		JOIN compliance_action_plans
		ON compliance_action_plans.id =
			compliance_action_plan_status_updates.compliance_action_plan_id
		WHERE compliance_action_plans.compliance_register_id = ?`,

		`DELETE compliance_action_plan_status_updates
		FROM compliance_action_plan_status_updates
		JOIN compliance_action_plans
		ON compliance_action_plans.id =
			compliance_action_plan_status_updates.compliance_action_plan_id
		WHERE compliance_action_plans.compliance_register_id = ?`,

		`DELETE tasks
		FROM tasks
		JOIN compliance_action_plans
		ON compliance_action_plans.id =
			tasks.compliance_action_plan_id
		WHERE compliance_action_plans.compliance_register_id = ?`,

		`DELETE FROM compliance_action_plans
		WHERE compliance_register_id = ?`,
		//////////////////
		`DELETE compliance_controls
		FROM compliance_controls 
		JOIN compliance_requirements ON compliance_requirements.id =compliance_controls.compliance_requirement_id
		JOIN compliance_clauses ON compliance_clauses.id = compliance_requirements.compliance_clause_id
		WHERE compliance_clauses.compliance_register_id = ?`,

		`DELETE compliance_requirements
		 FROM compliance_requirements
		 JOIN compliance_clauses ON compliance_clauses.id = compliance_requirements.compliance_clause_id		 
		 WHERE compliance_clauses.compliance_register_id = ?`,

		`DELETE audit_finding_compliance_requirements
		FROM audit_finding_compliance_requirements
		JOIN audit_findings
		ON audit_findings.id =
			audit_finding_compliance_requirements.audit_finding_id
		JOIN audit_checklists
		ON audit_checklists.id =
			audit_findings.audit_checklist_id
		JOIN audit_plan_checklists
		ON audit_plan_checklists.id =
			audit_checklists.audit_plan_checklist_id
		JOIN auditable_items
		ON auditable_items.id =
			audit_plan_checklists.auditable_item_id
		JOIN compliance_clauses
		ON compliance_clauses.id =
			audit_finding_compliance_requirements.compliance_clause_id
		WHERE compliance_clauses.compliance_register_id = ?`,

		`DELETE audit_finding_compliance_requirements
		FROM audit_finding_compliance_requirements
		JOIN compliance_clauses
		ON compliance_clauses.id =
			audit_finding_compliance_requirements.compliance_clause_id
		WHERE compliance_clauses.compliance_register_id = ?`,

		`DELETE FROM compliance_clauses WHERE compliance_clauses.compliance_register_id = ?`,
		`DELETE FROM compliance_assessments WHERE compliance_assessments.compliance_register_id = ?`,
		`DELETE FROM compliance_controls WHERE compliance_register_id = ?`,

		/* ========= WORKFLOWS ========= */
		`DELETE compliance_register_workflow_reviews
		 FROM compliance_register_workflow_reviews
		 JOIN compliance_register_workflows
		 ON compliance_register_workflows.id = compliance_register_workflow_reviews.compliance_register_workflow_id
		 WHERE compliance_register_workflows.compliance_register_id = ?`,

		`DELETE compliance_register_workflow_participants
		 FROM compliance_register_workflow_participants
		 JOIN compliance_register_workflows
		 ON compliance_register_workflows.id = compliance_register_workflow_participants.compliance_register_workflow_id
		 WHERE compliance_register_workflows.compliance_register_id = ?`,

		`DELETE FROM compliance_register_workflow_histories WHERE compliance_register_id = ?`,
		`DELETE FROM compliance_register_workflows WHERE compliance_register_id = ?`,
	}

	return db.Transaction(func(tx *gorm.DB) error {
		for _, rid := range registerIDs {
			for i, q := range queries {
				if err := tx.Exec(q, rid).Error; err != nil {
					return fmt.Errorf("force delete failed for register %d at step %d: %w", rid, i, err)
				}
			}
			if err := tx.Delete(&models.ComplianceRegister{}, rid).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// Issuing Authority Methods
func (repo *ComplianceRegisterRepository) CreateIssuingAuthority(authority *models.ComplianceRegisterIssuingAuthority, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Create(authority).Error; err != nil {
		return 0, err
	}
	return authority.ID, nil
}

func (repo *ComplianceRegisterRepository) UpdateIssuingAuthority(authority *models.ComplianceRegisterIssuingAuthority, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Save(authority).Error
}

func (repo *ComplianceRegisterRepository) DeleteIssuingAuthority(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Delete(&models.ComplianceRegisterIssuingAuthority{}, id).Error
}

func (repo *ComplianceRegisterRepository) GetIssuingAuthorityByID(id uint, tx ...utility.DatabaseTransaction) (*models.ComplianceRegisterIssuingAuthority, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var authority models.ComplianceRegisterIssuingAuthority
	err := db.First(&authority, id).Error
	if err != nil {
		return nil, err
	}
	return &authority, nil
}

func (repo *ComplianceRegisterRepository) GetIssuingAuthoritiesByComplianceRegisterID(complianceRegisterID uint, tx ...utility.DatabaseTransaction) ([]*models.ComplianceRegisterIssuingAuthority, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var authorities []*models.ComplianceRegisterIssuingAuthority
	if err := db.Where("compliance_register_id = ?", complianceRegisterID).
		Preload("Entity").
		Preload("Document").
		Find(&authorities).Error; err != nil {
		return nil, err
	}

	return authorities, nil
}

func (repo *ComplianceRegisterRepository) CreateIssuingAuthorityEntity(entity *models.ComplianceRegisterIssuingAuthorityEntity, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Create(entity).Error; err != nil {
		return 0, err
	}
	return entity.ID, nil
}

func (repo *ComplianceRegisterRepository) UpdateIssuingAuthorityEntity(entity *models.ComplianceRegisterIssuingAuthorityEntity, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Save(entity).Error
}

func (repo *ComplianceRegisterRepository) DeleteIssuingAuthorityEntity(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Delete(&models.ComplianceRegisterIssuingAuthorityEntity{}, id).Error
}

func (repo *ComplianceRegisterRepository) GetIssuingAuthorityEntityByID(id uint, tx ...utility.DatabaseTransaction) (*models.ComplianceRegisterIssuingAuthorityEntity, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var entity *models.ComplianceRegisterIssuingAuthorityEntity
	if err := db.First(&entity, id).Error; err != nil {
		return nil, err
	}

	return entity, nil
}

func (repo *ComplianceRegisterRepository) CreateIssuingAuthorityDocument(document *models.ComplianceRegisterIssuingAuthorityDocument, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Create(document).Error; err != nil {
		return 0, err
	}
	return document.ID, nil
}

func (repo *ComplianceRegisterRepository) UpdateIssuingAuthorityDocument(document *models.ComplianceRegisterIssuingAuthorityDocument, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Save(document).Error
}

func (repo *ComplianceRegisterRepository) DeleteIssuingAuthorityDocument(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Delete(&models.ComplianceRegisterIssuingAuthorityDocument{}, id).Error
}

func (repo *ComplianceRegisterRepository) GetIssuingAuthorityDocumentByID(id uint, tx ...utility.DatabaseTransaction) (*models.ComplianceRegisterIssuingAuthorityDocument, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var document models.ComplianceRegisterIssuingAuthorityDocument
	err := db.First(&document, id).Error
	if err != nil {
		return nil, err
	}
	return &document, nil
}

func (repo *ComplianceRegisterRepository) GetDocumentsByAuthorityID(authorityID uint, tx ...utility.DatabaseTransaction) ([]*models.ComplianceRegisterIssuingAuthorityDocument, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var documents []*models.ComplianceRegisterIssuingAuthorityDocument
	err := db.Where("compliance_register_issuing_authority_id = ?", authorityID).Find(&documents).Error
	if err != nil {
		return nil, err
	}
	return documents, nil
}

func (repo *ComplianceRegisterRepository) DeleteIssuingAuthorityDocumentsByAuthorityID(authorityID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Where("compliance_register_issuing_authority_id = ?", authorityID).Delete(&models.ComplianceRegisterIssuingAuthorityDocument{}).Error
}

func (repo *ComplianceRegisterRepository) DeleteIssuingAuthorityEntitiesByAuthorityID(authorityID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Where("compliance_register_issuing_authority_id = ?", authorityID).Delete(&models.ComplianceRegisterIssuingAuthorityEntity{}).Error
}

// CreateSharedUnits creates multiple shared units for a compliance register
func (repo *ComplianceRegisterRepository) CreateSharedUnits(sharedUnits []*models.ComplianceRegisterSharedUnit, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if len(sharedUnits) == 0 {
		return nil
	}

	return db.Create(&sharedUnits).Error
}

// DeleteSharedUnits deletes shared units for a compliance register by department IDs
func (repo *ComplianceRegisterRepository) DeleteSharedUnits(complianceRegisterID uint, departmentIDs []uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if len(departmentIDs) == 0 {
		return nil
	}

	return db.Where("compliance_register_id = ? AND department_id IN ?", complianceRegisterID, departmentIDs).Delete(&models.ComplianceRegisterSharedUnit{}).Error
}

// GetSharedUnitsByComplianceRegisterID retrieves all shared units for a compliance register
func (repo *ComplianceRegisterRepository) GetSharedUnitsByComplianceRegisterID(complianceRegisterID uint, tx ...utility.DatabaseTransaction) ([]*models.ComplianceRegisterSharedUnit, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var sharedUnits []*models.ComplianceRegisterSharedUnit
	err := db.Where("compliance_register_id = ?", complianceRegisterID).Find(&sharedUnits).Error
	return sharedUnits, err
}

func (repo *ComplianceRegisterRepository) UpdateNote(id uint, note *string, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Model(&models.ComplianceRegister{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"note": note,
		}).Error
}

func (repo *ComplianceRegisterRepository) GetAllComplianceTitles() ([]string, error) {
	var titles []string
	err := repo.db.Model(&models.ComplianceRegister{}).Pluck("title_en", &titles).Error
	return titles, err
}
