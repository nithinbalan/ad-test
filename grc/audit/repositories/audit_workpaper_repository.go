package repositories

import (
	auditDTOs "jethur/grc/audit/dtos"
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type AuditWorkpaperRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(auditID uint, params *auditDTOs.AuditWorkpaperListQueryParamsDTO) (*auditDTOs.AuditWorkpaperListResponseDTO, error)
	GetByID(id uint) (*models.AuditWorkpaper, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.AuditWorkpaper, error)
	Create(auditWorkPaper *models.AuditWorkpaper, tx ...utility.DatabaseTransaction) (uint, error)
	Update(auditWorkPaper *models.AuditWorkpaper, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error
}

type AuditWorkpaperRepository struct {
	db *gorm.DB
}

func NewAuditWorkpaperRepository(db *gorm.DB) *AuditWorkpaperRepository {
	return &AuditWorkpaperRepository{
		db: db,
	}
}

// BeginTransaction starts a new transaction.
func (repo *AuditWorkpaperRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *AuditWorkpaperRepository) GetAll(auditID uint, params *auditDTOs.AuditWorkpaperListQueryParamsDTO) (*auditDTOs.AuditWorkpaperListResponseDTO, error) {
	var auditWorkPapers []*auditDTOs.AuditWorkpaperListDTO
	var totalCount int64

	query := repo.db.Table("audit_workpapers").
		Select(`
			audit_workpapers.id,
			audit_workpapers.audit_id,
			audit_workpapers.source_id,
			sources.title_en AS source_title_en,
			sources.title_ar AS source_title_ar,
			audit_workpapers.title,
			audit_workpapers.description,
			audit_workpapers.url,
			audit_workpapers.name,
			audit_workpapers.size,
			audit_workpapers.token,
			audit_workpapers.original_name,
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
			LEFT JOIN sources ON sources.id = audit_workpapers.source_id
			LEFT JOIN users AS created_user ON created_user.id = audit_workpapers.created_by_id
			LEFT JOIN users AS updated_user ON updated_user.id = audit_workpapers.updated_by_id
		`).Where(`audit_workpapers.audit_id = ?`, auditID)

	// Apply dynamic filters to the query
	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where(columnName+" IN (?)", values)
		}
	}
	// Calculate total count for pagination
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	// Apply search keyword if exists
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
			audit_workpapers.title LIKE ? OR
			audit_workpapers.description LIKE ? OR
			sources.title_en LIKE ? OR
			sources.title_ar LIKE ? OR
			created_user.first_name LIKE ? OR
			created_user.last_name LIKE ? OR
			updated_user.first_name LIKE ? OR
			updated_user.last_name LIKE ?`,
			keyword, keyword, keyword, keyword, keyword, keyword, keyword)
	}

	// Sort by the specified column and order
	if params.SortBy == "" {
		params.SortBy = "title" // Default sort by title if no sortBy is specified
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "asc" // Default order is ascending if not specified
	}

	query = query.Order(params.SortBy + " " + params.Sort)

	offset := (params.Page - 1) * params.PageSize

	// Execute the query
	err := query.Offset(offset).Limit(params.PageSize).Find(&auditWorkPapers).Error
	if err != nil {
		return nil, err
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	response := &auditDTOs.AuditWorkpaperListResponseDTO{
		Items: auditWorkPapers,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}
	return response, nil
}

func (repo *AuditWorkpaperRepository) GetByID(id uint) (*models.AuditWorkpaper, error) {
	var auditWorkPaper *models.AuditWorkpaper

	result := repo.db.First(&auditWorkPaper, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return auditWorkPaper, nil
}

func (repo *AuditWorkpaperRepository) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.AuditWorkpaper, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var auditWorkPaper *models.AuditWorkpaper

	result := db.
		Preload("Audit").
		Preload("Source").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		First(&auditWorkPaper, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return auditWorkPaper, nil
}

func (repo *AuditWorkpaperRepository) Create(auditWorkPaper *models.AuditWorkpaper, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(auditWorkPaper).Error; err != nil {
		return 0, err
	}

	return auditWorkPaper.ID, nil
}

func (repo *AuditWorkpaperRepository) Update(auditWorkPaper *models.AuditWorkpaper, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(auditWorkPaper).Error; err != nil {
		return err
	}

	return nil
}

func (repo *AuditWorkpaperRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	// Delete the Audit WorkPaper
	if err := db.Delete(&models.AuditWorkpaper{}, id).Error; err != nil {
		return err
	}

	return nil
}
