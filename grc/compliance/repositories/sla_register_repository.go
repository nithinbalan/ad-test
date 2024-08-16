package repositories

import (
	"fmt"
	commonDTOs "jethur/grc/common/dtos"
	slaDTOs "jethur/grc/compliance/dtos"
	"jethur/grc/models"
	"jethur/utility"
	"strconv"
	"time"

	"gorm.io/gorm"
)

// Interfaces
type SLARegisterRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(params *slaDTOs.SLARegisterListQueryParamsDTO) (*slaDTOs.SLARegisterListResponseDTO, error)
	GetByID(id uint) (*models.SLARegister, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.SLARegister, error)
	Create(slaRegister *models.SLARegister, tx ...utility.DatabaseTransaction) (uint, error)
	Update(slaRegister *models.SLARegister, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error
	GetIDsByLifecycleStatus(ids []uint, status string) ([]uint, error)
	UpdateLifecycleStatusForIDs(ids []uint, statusId uint) ([]uint, error)
	GenerateReferenceNumber() (string, error)
	CreateSLACheck(slaCheck *models.SLACheck, tx ...utility.DatabaseTransaction) error
	DeleteSLAChecksBySLAID(slaID uint, tx ...utility.DatabaseTransaction) error

	//Document
	CreateDocument(doc *models.SLADocument, tx ...utility.DatabaseTransaction) error
	GetDocumentByID(id uint) (*models.SLADocument, error)
	DeleteDocuments(docIds []*uint, tx ...utility.DatabaseTransaction) error
	GetLatestDocumentBySLAID(id uint) (*models.SLADocument, error)

	CreateSLAProcesses(processes []*models.SLAProcess, tx ...utility.DatabaseTransaction) error
	DeleteSLAProcesses(processIds []*uint, tx ...utility.DatabaseTransaction) error
}

// implementation
type SLARegisterRepository struct {
	db *gorm.DB
}

func NewSLARegisterRepository(db *gorm.DB) *SLARegisterRepository {
	return &SLARegisterRepository{
		db: db,
	}
}

// BeginTransaction starts a new transaction.
func (repo *SLARegisterRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *SLARegisterRepository) GetAll(params *slaDTOs.SLARegisterListQueryParamsDTO) (*slaDTOs.SLARegisterListResponseDTO, error) {
	var slaRegisters []*slaDTOs.SLARegisterListDTO
	var totalCount int64
	query := repo.db.Table("sla_register").
		Select(`
            sla_register.*,
			sla_documents.version,
			sla_documents.issue_date,
			sla_documents.expiry_date,
			sla_documents.sa1_date,
			sla_documents.sa2_date,
			sla_documents.comment,
			sla_documents.url,
			sla_documents.name,
			sla_documents.size,
			sla_documents.token,
			sla_documents.type,
			sla_documents.original_name,
			sla_documents.is_latest,
            sla_categories.title_en AS sla_category_title_en,
            sla_categories.title_ar AS sla_category_title_ar,
            severity_levels.title_en AS severity_level_title_en,
            severity_levels.title_ar AS severity_level_title_ar,
            organization_units.title_en AS responsible_unit_title_en,
            organization_units.title_ar AS responsible_unit_title_ar,
            responsible_users.first_name_en AS responsible_user_first_name_en,
            responsible_users.first_name_ar AS responsible_user_first_name_ar,
            responsible_users.last_name_en AS responsible_user_last_name_en,
            responsible_users.last_name_ar AS responsible_user_last_name_ar,
			responsible_users.initial AS responsible_user_initial,
            responsible_users.color AS responsible_user_color,
            responsible_users.email AS responsible_user_email,
            responsible_users.profile_image_path AS responsible_user_profile_image_path,
            lifecycle_statuses.title_en AS lifecycle_status_title_en,
            lifecycle_statuses.title_ar AS lifecycle_status_title_ar,
			lifecycle_statuses.color_code AS lifecycle_status_color_code,
			compliance_statuses.title_en AS compliance_status_title_en,
			compliance_statuses.title_ar AS compliance_status_title_ar,
			compliance_statuses.color_code AS compliance_status_color_code,
            created_user.id AS created_user_id,
        	created_user.first_name_en AS created_user_first_name_en,
            created_user.first_name_ar AS created_user_first_name_ar,
            created_user.last_name_en AS created_user_last_name_en,
            created_user.last_name_ar AS created_user_last_name_ar,
			created_user.initial AS created_user_initial,
            created_user.color AS created_user_color,
            created_user.email AS created_user_email,
            created_user.profile_image_path AS created_user_profile_image_path,
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
            LEFT JOIN sla_categories ON sla_register.sla_category_id = sla_categories.id
            LEFT JOIN severity_levels ON sla_register.severity_level_id = severity_levels.id
            LEFT JOIN organization_units ON sla_register.responsible_unit_id = organization_units.id
            LEFT JOIN users AS responsible_users ON sla_register.responsible_user_id = responsible_users.id
            JOIN lifecycle_statuses ON sla_register.lifecycle_status_id = lifecycle_statuses.id
		    LEFT JOIN compliance_statuses ON sla_register.compliance_status_id = compliance_statuses.id 
            JOIN users AS created_user ON sla_register.created_by_id = created_user.id
            LEFT JOIN users AS updated_user ON sla_register.updated_by_id = updated_user.id
			LEFT JOIN sla_documents ON sla_register.id = sla_documents.sla_register_id AND sla_documents.is_latest = 1

    `)

	// Apply dynamic filters to the query
	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where(columnName+" IN (?)", values)
		}
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

	// Apply quick filters
	for _, filter := range params.QuickFilter {
		switch filter {
		case "compliant":
			query = query.Where("compliance_statuses.type = ?", "compliant")
		case "non_compliant":
			query = query.Where("compliance_statuses.type = ?", "not_implemented")
		case "archived":
			query = query.Where("lifecycle_statuses.type = ?", "archived")

		}
	}

	// Calculate total count for pagination
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}
	// Apply search Keyword if exists
	if params.Keyword != "" {
		query = query.Where(`
			(sla_register.title_en LIKE ? OR
			sla_register.title_ar LIKE ? OR
			sla_register.reference_number LIKE ? OR
			sla_register.description LIKE ? OR
			sla_register.client LIKE ? OR
			sla_documents.issue_date LIKE ? OR
			sla_documents.expiry_date LIKE ? OR
			sla_documents.sa1_date LIKE ? OR
			sla_documents.sa2_date LIKE ? OR
			sla_documents.name LIKE ? OR
			sla_documents.type LIKE ? OR
			sla_categories.title_en LIKE ? OR
			sla_categories.title_ar LIKE ? OR
			organization_units.title_en LIKE ? OR
			organization_units.title_ar LIKE ? OR
			responsible_users.first_name_en LIKE ? OR
			responsible_users.first_name_ar LIKE ? OR
			responsible_users.last_name_en LIKE ? OR
			responsible_users.last_name_ar LIKE ? OR
			responsible_users.email LIKE ? OR
            responsible_users.initial LIKE ? OR
            responsible_users.color LIKE ? OR
            responsible_users.profile_image_path LIKE ? OR
			severity_levels.title_en LIKE ? OR
			severity_levels.title_ar LIKE ? OR
			compliance_statuses.title_en LIKE ? OR
			compliance_statuses.title_ar LIKE ? OR
			lifecycle_statuses.title_en LIKE ? OR
			lifecycle_statuses.title_ar LIKE ?)`,
			"%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%",
			"%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%",
			"%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%",
			"%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%",
			"%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%",
			"%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%",
			"%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%")
	}

	// Sort by the specified column and order
	if params.SortBy == "" {
		params.SortBy = "title_en" // Default Sort by title_en if no SortBy is specified
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "asc" // Default order is ascending if not specified
	}

	query = query.Order(params.SortBy + " " + params.Sort)

	offset := (params.Page - 1) * params.PageSize

	// Execute the queryc
	err := query.Offset(offset).Limit(params.PageSize).Find(&slaRegisters).Error

	if err != nil {
		return nil, err
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	response := &slaDTOs.SLARegisterListResponseDTO{
		Items: slaRegisters,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}
	return response, nil
}

func (repo *SLARegisterRepository) GetByID(id uint) (*models.SLARegister, error) {
	var slaRegister *models.SLARegister

	result := repo.db.First(&slaRegister, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return slaRegister, nil
}

func (repo *SLARegisterRepository) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.SLARegister, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var slaRegister *models.SLARegister

	result := db.
		Preload("SLACategory").
		Preload("SeverityLevel").
		Preload("ResponsibleUnit").
		Preload("ResponsibleUser").
		Preload("LifecycleStatus").
		Preload("ComplianceStatus").
		Preload("Documents").
		Preload("Checks").
		Preload("Processes").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		First(&slaRegister, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return slaRegister, nil
}

func (repo *SLARegisterRepository) Create(slaRegister *models.SLARegister, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(slaRegister).Error; err != nil {
		return 0, err
	}

	return slaRegister.ID, nil
}

func (repo *SLARegisterRepository) Update(slaRegister *models.SLARegister, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(slaRegister).Error; err != nil {
		return err
	}

	return nil
}

func (repo *SLARegisterRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	// Delete the SLA Register
	if err := db.Delete(&models.SLARegister{}, id).Error; err != nil {
		return err
	}

	return nil
}

func (repo *SLARegisterRepository) GetIDsByLifecycleStatus(ids []uint, status string) ([]uint, error) {
	var matchingIDs []uint

	if err := repo.db.Model(&models.SLARegister{}).
		Joins("JOIN lifecycle_statuses ON sla_register.lifecycle_status_id = lifecycle_statuses.id").
		Where("sla_register.id IN ? AND lifecycle_statuses.type = ?", ids, status).
		Pluck("sla_register.id", &matchingIDs).Error; err != nil {
		return nil, err
	}

	return matchingIDs, nil
}

func (repo *SLARegisterRepository) UpdateLifecycleStatusForIDs(ids []uint, statusId uint) ([]uint, error) {

	if err := repo.db.Model(&models.SLARegister{}).
		Where("id IN ?", ids).
		Update("lifecycle_status_id", statusId).Error; err != nil {
		return nil, err
	}

	return ids, nil
}
func (repo *SLARegisterRepository) GenerateReferenceNumber() (string, error) {
	var lastReferenceNumber string
	err := repo.db.Model(&models.SLARegister{}).
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
	lastReferenceID, _ := strconv.ParseUint(lastDigits, 10, 64)

	newReferenceNumber := fmt.Sprintf("%s-%d-%04d", "SLA", currentYear, lastReferenceID+1)

	return newReferenceNumber, nil
}

func (repo *SLARegisterRepository) GetDocumentByID(id uint) (*models.SLADocument, error) {
	var document models.SLADocument
	if err := repo.db.Where("id = ?", id).First(&document).Error; err != nil {
		return nil, err
	}
	return &document, nil
}

func (repo *SLARegisterRepository) CreateDocument(doc *models.SLADocument, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(doc).Error; err != nil {
		return err
	}
	return nil
}

func (repo *SLARegisterRepository) CreateSLACheck(slaCheck *models.SLACheck, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(slaCheck).Error; err != nil {
		return err
	}

	return nil
}

func (repo *SLARegisterRepository) DeleteDocuments(docIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	// Delete SLA documents by IDs
	if err := db.Where("id IN (?)", docIds).Delete(&models.SLADocument{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *SLARegisterRepository) DeleteSLAChecksBySLAID(slaID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Where("sla_register_id = ?", slaID).Delete(&models.SLACheck{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *SLARegisterRepository) GetLatestDocumentBySLAID(id uint) (*models.SLADocument, error) {
	var document models.SLADocument
	if err := repo.db.Where("sla_register_id = ? AND is_latest = ?", id, true).First(&document).Error; err != nil {
		return nil, err
	}
	return &document, nil
}

func (repo *SLARegisterRepository) CreateSLAProcesses(processes []*models.SLAProcess, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(&processes).Error; err != nil {
		return err
	}
	return nil
}

func (repo *SLARegisterRepository) DeleteSLAProcesses(processIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	// Delete processes by IDs
	if err := db.Where("id IN (?)", processIds).Delete(&models.SLAProcess{}).Error; err != nil {
		return err
	}

	return nil
}
