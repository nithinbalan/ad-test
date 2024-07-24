package repositories

import (
	"fmt"
	commonDTOs "jethur/grc/common/dtos"
	documentDTOs "jethur/grc/documents/dtos"
	"jethur/grc/models"
	"jethur/utility"
	"strconv"
	"time"

	"gorm.io/gorm"
)

// Interfaces
type DocumentRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(params *documentDTOs.DocumentListQueryParamsDTO) (*documentDTOs.DocumentListResponseDTO, error)
	GetByID(id uint, tx ...utility.DatabaseTransaction) (*models.Document, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.Document, error)
	Create(document *models.Document, tx ...utility.DatabaseTransaction) (uint, error)
	Update(document *models.Document, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error
	GenerateReferenceNumber() (string, error)

	CreateSharedUsers(sharedUsers []*models.DocumentSharedUser, tx ...utility.DatabaseTransaction) error
	DeleteSharedUsers(sharedUserIds []*uint, tx ...utility.DatabaseTransaction) error
	CreateSharedUnits(sharedUnits []*models.DocumentSharedUnit, tx ...utility.DatabaseTransaction) error
	DeleteSharedUnits(sharedUnitIds []*uint, tx ...utility.DatabaseTransaction) error

	GetDocumentVersionByID(id uint, tx ...utility.DatabaseTransaction) (*models.DocumentVersion, error)
	CreateDocumentVersion(documentVersion *models.DocumentVersion, tx ...utility.DatabaseTransaction) (uint, error)
	UpdateDocumentVersion(documentVersion *models.DocumentVersion, tx ...utility.DatabaseTransaction) error
	DeleteDocumentVersions(docIds []*uint, tx ...utility.DatabaseTransaction) error

	CreateDocumentEditHistory(documentHistory *models.DocumentEditHistory, tx ...utility.DatabaseTransaction) (uint, error)
	DeleteDocumentEditHistory(docHistoryEditIds []*uint, tx ...utility.DatabaseTransaction) error
}

// Implementation
type DocumentRepository struct {
	db *gorm.DB
}

func NewDocumentRepository(db *gorm.DB) *DocumentRepository {
	return &DocumentRepository{
		db: db,
	}
}

func (repo *DocumentRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *DocumentRepository) GetAll(params *documentDTOs.DocumentListQueryParamsDTO) (*documentDTOs.DocumentListResponseDTO, error) {
	var documents []*documentDTOs.DocumentListDTO
	var totalCount int64

	query := repo.db.Table("documents").
		Select(`
			documents.*,
			document_categories.title_en AS document_category_title_en,
			document_categories.title_ar AS document_category_title_ar,
			document_types.title_en AS document_type_title_en,
			document_types.title_ar AS document_type_title_ar,
			document_access_types.title_en AS document_access_type_title_en,
			document_access_types.title_ar AS document_type_access_title_ar,
			organization_units.title_en AS responsible_unit_title_en,
			organization_units.title_ar AS responsible_unit_title_ar,
			responsible_user.first_name_en AS responsible_user_first_name_en,
            responsible_user.first_name_ar AS responsible_user_first_name_ar,
            responsible_user.last_name_en AS responsible_user_last_name_en,
            responsible_user.last_name_ar AS responsible_user_last_name_ar,
			responsible_user.initial AS responsible_user_initial,
            responsible_user.color AS responsible_user_color,
            responsible_user.email AS responsible_user_email,
			document_versions.url AS document_url,
			document_versions.name AS document_name,
			document_versions.size AS document_size,
			document_versions.token AS document_token,
			document_versions.original_name AS document_original_name,
			document_versions.type AS document_type,
			document_versions.version AS document_version,
			document_statuses.title_en AS document_status_title_en,
			document_statuses.title_ar AS document_status_title_ar,
			document_statuses.color_code AS document_status_color_code,
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
			LEFT JOIN document_categories ON document_categories.id = documents.document_category_id
			LEFT JOIN document_types ON document_types.id = documents.document_type_id
			LEFT JOIN document_access_types ON document_access_types.id = documents.document_access_type_id
			LEFT JOIN organization_units ON organization_units.id = documents.responsible_unit_id
			LEFT JOIN users AS responsible_user ON responsible_user.id = documents.responsible_user_id
			LEFT JOIN document_versions ON document_versions.document_id = documents.id AND document_versions.is_latest_version = 1
			LEFT JOIN document_statuses ON document_statuses.id = documents.document_status_id
			LEFT JOIN users AS created_user ON created_user.id = documents.created_by_id
			LEFT JOIN users AS updated_user ON updated_user.id = documents.updated_by_id
		`)

		// Apply dynamic filters to the query
	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where(columnName+" IN (?)", values)
		}
	}

	// Apply quick filters
	for _, filter := range params.QuickFilter {
		switch filter {
		case "policy":
			query = query.Where("document_types.type = ?", "policy")
		case "procedure":
			query = query.Where("document_types.type = ?", "procedure")
		case "in_review":
			query = query.Where("document_statuses.type = ?", "in-review")
		case "draft":
			query = query.Where("document_statuses.type = ?", "draft")

		}
	}

	// Calculate total count for pagination
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	// Apply search keyword if exists
	if params.Keyword != "" {
		query = query.Where(`
				(documents.title_en LIKE ? OR
				documents.title_ar LIKE ? OR
				documents.reference_number LIKE ? OR
				documents.description_en LIKE ? OR
				documents.description_ar LIKE ? OR
				document_categories.title_en LIKE ? OR
				document_categories.title_ar LIKE ? OR
				document_statuses.title_en LIKE ? OR
				document_statuses.title_ar LIKE ? OR
				document_types.title_en LIKE ? OR
				document_types.title_ar LIKE ? OR
				organization_units.title_en LIKE ? OR
				organization_units.title_ar LIKE ? OR
				responsible_user.first_name_en LIKE ? OR
				responsible_user.first_name_ar LIKE ? OR
				responsible_user.last_name_en LIKE ? OR
				responsible_user.last_name_ar LIKE ? OR
				responsible_user.email LIKE ? OR
				responsible_user.initial LIKE ?)`,
			"%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%",
			"%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%",
			"%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%",
			"%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%",
			"%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%")
	}

	// Sort by the specified column and order
	if params.SortBy == "" {
		params.SortBy = "title_en" // Default sort by title_en if no sortBy is specified
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "asc" // Default order is ascending if not specified
	}

	query = query.Order(params.SortBy + " " + params.Sort)

	offset := (params.Page - 1) * params.PageSize

	// Execute the queryc
	err := query.Offset(offset).Limit(params.PageSize).Find(&documents).Error
	if err != nil {
		return nil, err
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	response := &documentDTOs.DocumentListResponseDTO{
		Items: documents,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}
	return response, nil

}

func (repo *DocumentRepository) GetByID(id uint, tx ...utility.DatabaseTransaction) (*models.Document, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var document *models.Document

	result := db.First(&document, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return document, nil
}

func (repo *DocumentRepository) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.Document, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var document *models.Document
	result := db.
		Preload("DocumentType").
		Preload("DocumentCategory").
		Preload("ResponsibleUser").
		Preload("ResponsibleUnit").
		Preload("DocumentAccessType").
		Preload("DocumentStatus").
		Preload("LatestDocumentVersion", "is_latest_version = ?", true).
		Preload("LatestDocumentVersion.CreatedByUser").
		Preload("LatestDocumentVersion.UpdatedByUser").
		Preload("DocumentVersions").
		Preload("DocumentVersions.DocumentEditHistory").
		Preload("DocumentSharedUsers").
		Preload("DocumentSharedUsers.User").
		Preload("DocumentSharedUsers.CreatedByUser").
		Preload("DocumentSharedUnits").
		Preload("DocumentSharedUnits.OrganizationUnit").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		First(&document, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return document, nil
}

func (repo *DocumentRepository) Create(document *models.Document, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(document).Error; err != nil {
		return 0, err
	}

	return document.ID, nil
}

func (repo *DocumentRepository) Update(document *models.Document, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(document).Error; err != nil {
		return err
	}

	return nil
}

func (repo *DocumentRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	// Delete document
	if err := db.Delete(&models.Document{}, id).Error; err != nil {
		return err
	}

	return nil
}

func (repo *DocumentRepository) GetDocumentVersionByID(id uint, tx ...utility.DatabaseTransaction) (*models.DocumentVersion, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var documentVersion *models.DocumentVersion

	result := db.First(&documentVersion, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return documentVersion, nil
}

func (repo *DocumentRepository) CreateDocumentVersion(documentVersion *models.DocumentVersion, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(documentVersion).Error; err != nil {
		return 0, err
	}
	return documentVersion.ID, nil
}

func (repo *DocumentRepository) UpdateDocumentVersion(documentVersion *models.DocumentVersion, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(documentVersion).Error; err != nil {
		return err
	}

	return nil
}

func (repo *DocumentRepository) DeleteDocumentVersions(docIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("id IN (?)", docIds).Delete(&models.DocumentVersion{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *DocumentRepository) CreateSharedUsers(sharedUsers []*models.DocumentSharedUser, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Create(&sharedUsers).Error; err != nil {
		return err
	}
	return nil
}

func (repo *DocumentRepository) DeleteSharedUsers(sharedUserIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("id IN (?)", sharedUserIds).Delete(&models.DocumentSharedUser{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *DocumentRepository) CreateSharedUnits(sharedUnits []*models.DocumentSharedUnit, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(&sharedUnits).Error; err != nil {
		return err
	}
	return nil
}

func (repo *DocumentRepository) DeleteSharedUnits(sharedUnitIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("id IN (?)", sharedUnitIds).Delete(&models.DocumentSharedUnit{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *DocumentRepository) GenerateReferenceNumber() (string, error) {
	var lastReferenceNumber string
	err := repo.db.Model(&models.Document{}).
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
	newReferenceNumber := fmt.Sprintf("%s-%d-%04d", "DC", currentYear, lastReferenceID+1)

	return newReferenceNumber, nil
}

func (repo *DocumentRepository) CreateDocumentEditHistory(documentHistory *models.DocumentEditHistory, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(documentHistory).Error; err != nil {
		return 0, err
	}
	return documentHistory.ID, nil
}

func (repo *DocumentRepository) DeleteDocumentEditHistory(docEditHistoryIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("id IN (?)", docEditHistoryIds).Delete(&models.DocumentEditHistory{}).Error; err != nil {
		return err
	}

	return nil
}
