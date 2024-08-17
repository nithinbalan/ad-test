package repositories

import (
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/models"
	organizationDTOs "jethur/grc/organization/dtos"
	"jethur/utility"

	"gorm.io/gorm"
)

type ControlTestRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(controlID uint, params *organizationDTOs.ControlTestListQueryParamsDTO) (*organizationDTOs.ControlTestListResponseDTO, error)
	GetByID(id uint) (*models.ControlTest, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.ControlTest, error)
	Create(controlTest *models.ControlTest, tx ...utility.DatabaseTransaction) (uint, error)
	Update(controlTest *models.ControlTest, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error

	GetDocumentByID(id uint) (*models.ControlTestDocument, error)
	CreateDocuments(docs []*models.ControlTestDocument, tx ...utility.DatabaseTransaction) error
	DeleteDocuments(docIds []*uint, tx ...utility.DatabaseTransaction) error

	GetTrendAnalysis(controlID uint) ([]*models.ControlTest, error)
}

type ControlTestRepository struct {
	db *gorm.DB
}

func NewControlTestRepository(db *gorm.DB) *ControlTestRepository {
	return &ControlTestRepository{
		db: db,
	}
}

// BeginTransaction starts a new transaction.
func (repo *ControlTestRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *ControlTestRepository) GetAll(controlID uint, params *organizationDTOs.ControlTestListQueryParamsDTO) (*organizationDTOs.ControlTestListResponseDTO, error) {
	var controlsTests []*organizationDTOs.ControlTestListDTO
	var totalCount int64

	query := repo.db.Table("control_tests").
		Select(`
			control_tests.id,
			control_tests.control_id,
			control_tests.title,
			control_implementation_statuses.title_en AS control_implementation_status_title_en,
			control_implementation_statuses.title_ar AS control_implementation_status_title_ar,
			control_implementation_statuses.color_code AS control_implementation_status_color_code,
			control_efficiency_statuses.title_en AS control_efficiency_status_title_en,
			control_efficiency_statuses.title_ar AS control_efficiency_status_title_ar,
			control_efficiency_statuses.color_code AS control_efficiency_status_color_code,
			control_tests.description,
			control_tests.score,
			control_tests.comments,
			control_tests.tested_source,
			tested_by_user.id AS tested_by_user_id,
			tested_by_user.first_name_en AS tested_by_user_first_name_en,
			tested_by_user.first_name_ar AS tested_by_user_first_name_ar,
			tested_by_user.last_name_en AS tested_by_user_last_name_en,
			tested_by_user.last_name_ar AS tested_by_user_last_name_ar,
			tested_by_user.initial AS tested_by_user_initial,
        	tested_by_user.color AS tested_by_user_color,
        	tested_by_user.email AS tested_by_user_email,
        	tested_by_user.profile_image_path AS tested_by_user_profile_image_path,
			updated_by_user.id AS updated_user_id,
			updated_by_user.first_name_en AS updated_user_first_name_en,
			updated_by_user.first_name_ar AS updated_user_first_name_ar,
			updated_by_user.last_name_en AS updated_user_last_name_en,
			updated_by_user.last_name_ar AS updated_user_last_name_ar,
			updated_by_user.initial AS updated_user_initial,
        	updated_by_user.color AS updated_user_color,
        	updated_by_user.email AS updated_user_email,
        	updated_by_user.profile_image_path AS updated_user_profile_image_path, 
        	CONVERT(VARCHAR, control_tests.tested_at, 23) AS tested_at,
			control_tests.updated_at
        `).
		Joins("LEFT JOIN control_implementation_statuses ON control_implementation_statuses.id = control_tests.control_implementation_status_id").
		Joins("LEFT JOIN control_efficiency_statuses ON control_efficiency_statuses.id = control_tests.control_efficiency_status_id").
		Joins("LEFT JOIN users AS tested_by_user ON tested_by_user.id = control_tests.tested_by").
		Joins("LEFT JOIN users AS updated_by_user ON updated_by_user.id = control_tests.updated_by_id").
		Where("control_tests.control_id = ?", controlID)

	// Apply dynamic filters to the query
	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where(columnName+" IN (?)", values)
		}
	}

	// Apply tested_at filters
	if params.TestedAtFrom != "" {
		query = query.Where("tested_at >= ?", params.TestedAtFrom)
	}
	if params.TestedAtTo != "" {
		query = query.Where("tested_at <= ?", params.TestedAtTo)
	}

	//Apply quick filters
	for _, filter := range params.QuickFilter {
		switch filter {
		case "implemented":
			query = query.Where("control_implementation_statuses.type = ?", "implemented")
		case "not_implemented":
			query = query.Where("control_implementation_statuses.type = ?", "not_implemented")
		case "partially_implemented":
			query = query.Where("control_implementation_statuses.type = ?", "partially_implemented")
		case "efficient":
			query = query.Where("control_efficiency_statuses.type = ?", "efficient")
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
			  (control_tests.title LIKE ? OR
				control_tests.description LIKE ? OR
				control_tests.tested_source LIKE ? OR
				control_implementation_statuses.title_en LIKE ? OR
				control_implementation_statuses.title_ar LIKE ? OR
				control_efficiency_statuses.title_en LIKE ? OR
				control_efficiency_statuses.title_ar LIKE ? OR
				tested_by_user.first_name_en LIKE ? OR
				tested_by_user.first_name_ar LIKE ? OR
				tested_by_user.last_name_en LIKE ? OR
				tested_by_user.last_name_ar LIKE ?)`,
			"%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%",
			"%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%",
			"%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%")
	}

	// Sort by the specified column and order
	if params.SortBy == "" {
		params.SortBy = "title" // Default sort by title_en if no sortBy is specified
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "asc" // Default order is ascending if not specified
	}

	query = query.Order(params.SortBy + " " + params.Sort)

	offset := (params.Page - 1) * params.PageSize

	// Execute the queryc
	err := query.Offset(offset).Limit(params.PageSize).Find(&controlsTests).Error
	if err != nil {
		return nil, err
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	response := &organizationDTOs.ControlTestListResponseDTO{
		Items: controlsTests,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}
	return response, nil
}

func (repo *ControlTestRepository) GetByID(id uint) (*models.ControlTest, error) {
	var controlTest *models.ControlTest

	result := repo.db.First(&controlTest, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return controlTest, nil

}

func (repo *ControlTestRepository) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.ControlTest, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var controlTest *models.ControlTest

	result := db.
		Preload("Documents").
		Preload("ControlImplementationStatus").
		Preload("ControlEfficiencyStatus").
		Preload("TestedByUser").
		Preload("UpdatedByUser").
		First(&controlTest, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return controlTest, nil

}

func (repo *ControlTestRepository) Create(controlTest *models.ControlTest, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(controlTest).Error; err != nil {
		return 0, err
	}

	return controlTest.ID, nil
}

func (repo *ControlTestRepository) Update(controlTest *models.ControlTest, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(controlTest).Error; err != nil {
		return err
	}

	return nil
}

func (repo *ControlTestRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var controlTest *models.ControlTest
	if err := db.First(&controlTest, id).Error; err != nil {
		return err
	}

	result := db.Delete(&models.ControlTest{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *ControlTestRepository) GetDocumentByID(id uint) (*models.ControlTestDocument, error) {
	var document models.ControlTestDocument
	if err := repo.db.Where("id = ?", id).First(&document).Error; err != nil {
		return nil, err
	}
	return &document, nil
}

func (repo *ControlTestRepository) CreateDocuments(docs []*models.ControlTestDocument, tx ...utility.DatabaseTransaction) error {
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

func (repo *ControlTestRepository) DeleteDocuments(docIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	// Delete documents by IDs
	if err := db.Where("id IN (?)", docIds).Delete(&models.ControlTestDocument{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *ControlTestRepository) GetTrendAnalysis(controlID uint) ([]*models.ControlTest, error) {
	var controlTests []*models.ControlTest

	// Use GORM to fetch the last 6 control tests
	result := repo.db.Where("control_id = ?", controlID).
		Order("tested_at DESC").
		Limit(6).
		Find(&controlTests)

	if result.Error != nil {
		return nil, result.Error
	}

	return controlTests, nil
}
