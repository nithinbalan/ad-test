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
	GenerateReferenceNumber() (string, error)

	//lifecycle
	GetIDsByLifecycleStatus(ids []uint, status string) ([]uint, error)
	UpdateLifecycleStatusForIDs(ids []uint, statusId uint) ([]uint, error)

	CreateComplianceControls(controls []*models.ComplianceControl, tx ...utility.DatabaseTransaction) error
	DeleteComplianceControls(controlIds []*uint, tx ...utility.DatabaseTransaction) error

	//Document
	CreateDocument(doc *models.ComplianceDocument, tx ...utility.DatabaseTransaction) error
	GetDocumentByID(id uint) (*models.ComplianceDocument, error)
	DeleteDocuments(docIds []*uint, tx ...utility.DatabaseTransaction) error
	GetLatestDocumentByComplianceID(id uint, language string) (*models.ComplianceDocument, error)

	CreateInitialComplianceStatusUpdate(complianceStatusUpdate *models.ComplianceStatusUpdate, tx ...utility.DatabaseTransaction) (uint, error)
	DeleteComplianceStatusUpdates(ids []*uint, tx ...utility.DatabaseTransaction) error

	UpdateComplianceTotalScore(complianceRegisterID uint, totalScore float64, tx utility.DatabaseTransaction) error
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
        	CONVERT(VARCHAR, compliance_register.issue_date, 23) AS issue_date,
        	CONVERT(VARCHAR, compliance_register.expiry_date, 23) AS expiry_date,
       		CONVERT(VARCHAR, compliance_register.next_review_date, 23) AS next_review_date,
			CONVERT(VARCHAR, compliance_register.last_reviewed_date, 23) AS last_reviewed_date,
			compliance_categories.title_en AS compliance_category_title_en,
			compliance_categories.title_ar AS compliance_category_title_ar,
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
			review_frequencies.title_en AS review_frequency_title_en,
			review_frequencies.title_ar AS review_frequency_title_ar,
			review_responsible_user.id AS review_responsible_user_id,
			review_responsible_user.first_name_en AS review_responsible_user_first_name_en,
            review_responsible_user.first_name_ar AS review_responsible_user_first_name_ar,
            review_responsible_user.last_name_en AS review_responsible_user_last_name_en,
            review_responsible_user.last_name_ar AS review_responsible_user_last_name_ar,
			review_responsible_user.initial AS review_responsible_user_initial,
        	review_responsible_user.color AS review_responsible_user_color,
        	review_responsible_user.email AS review_responsible_user_email,
        	review_responsible_user.profile_image_path AS review_responsible_user_profile_image_path,  
			severity_levels.title_en AS severity_level_title_en,
			severity_levels.title_ar AS severity_level_title_ar,
			severity_levels.color_code AS severity_level_color_code,
			lifecycle_statuses.title_en AS lifecycle_status_title_en,
			lifecycle_statuses.title_ar AS lifecycle_status_title_ar,
			lifecycle_statuses.color_code AS lifecycle_status_color_code,
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
		    LEFT JOIN compliance_categories ON compliance_categories.id = compliance_register.compliance_category_id
		    LEFT JOIN compliance_statuses ON compliance_statuses.id = compliance_register.compliance_status_id
		    LEFT JOIN organization_units ON organization_units.id = compliance_register.responsible_unit_id
		    LEFT JOIN users AS responsible_user ON responsible_user.id = compliance_register.responsible_user_id
		    LEFT JOIN review_frequencies ON review_frequencies.id = compliance_register.review_frequency_id
			LEFT JOIN users AS review_responsible_user ON review_responsible_user.id = compliance_register.review_responsible_user_id
		    LEFT JOIN severity_levels ON severity_levels.id = compliance_register.severity_level_id
			LEFT JOIN lifecycle_statuses ON lifecycle_statuses.id = compliance_register.lifecycle_status_id
			LEFT JOIN users AS created_user ON created_user.id = compliance_register.created_by_id
			LEFT JOIN users AS updated_user ON updated_user.id = compliance_register.updated_by_id
		`)

	// Apply dynamic filters to the query
	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where(columnName+" IN (?)", values)
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

	// Apply quick filters
	for _, filter := range params.QuickFilter {
		switch filter {
		case "compliant":
			query = query.Where("compliance_statuses.type = ?", "compliant")
		case "partially_compliant":
			query = query.Where("compliance_statuses.type = ?", "partially_compliant")
		case "non_compliant":
			query = query.Where("compliance_statuses.type = ?", "not-implemented")
		case "archived":
			query = query.Where("lifecycle_statuses.type = ?", "archived")

		}
	}

	// Calculate total count for pagination
	if err := query.Count(&totalCount).Error; err != nil {
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
			compliance_categories.title_en LIKE ? OR
			compliance_categories.title_ar LIKE ? OR
			compliance_statuses.title_en LIKE ? OR
			compliance_statuses.title_ar LIKE ? OR
			organization_units.title_en LIKE ? OR
			organization_units.title_ar LIKE ? OR
			responsible_user.first_name_en LIKE ? OR
			responsible_user.first_name_ar LIKE ? OR
			responsible_user.last_name_en LIKE ? OR
			responsible_user.last_name_ar LIKE ? OR
			responsible_user.email LIKE ? OR
            responsible_user.initial LIKE ? OR
            responsible_user.color LIKE ? OR
            responsible_user.profile_image_path LIKE ? OR
			severity_levels.title_en LIKE ? OR
			severity_levels.title_ar LIKE ? OR
			lifecycle_statuses.title_en LIKE ? OR
			lifecycle_statuses.title_ar LIKE ?)`,
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
		Preload("ComplianceCategory").
		Preload("ComplianceStatus").
		Preload("ResponsibleUnit.ParentUnit.ParentUnit.ParentUnit.ParentUnit.ParentUnit").
		Preload("ResponsibleUser").
		Preload("ResponsibleUser.OrganizationUnit1").
		Preload("ReviewFrequency").
		Preload("ComplianceReviewStatus").
		Preload("LastReviewedUser").
		Preload("ReviewResponsibleUser").
		Preload("SeverityLevel").
		Preload("LifecycleStatus").
		Preload("DocumentEn", func(db *gorm.DB) *gorm.DB {
			return db.Where("language = ?", "en")
		}).
		Preload("DocumentAr", func(db *gorm.DB) *gorm.DB {
			return db.Where("language = ?", "ar")
		}).
		Preload("Controls").
		Preload("Controls.Control.ControlCategory").
		Preload("Controls.Control.ControlType").
		Preload("ComplianceStatusUpdates").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
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

func (repo *ComplianceRegisterRepository) CreateComplianceControls(controls []*models.ComplianceControl, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(&controls).Error; err != nil {
		return err
	}
	return nil
}

func (repo *ComplianceRegisterRepository) DeleteComplianceControls(controlIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	// Delete controls by IDs
	if err := db.Where("id IN (?)", controlIds).Delete(&models.ComplianceControl{}).Error; err != nil {
		return err
	}

	return nil
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
