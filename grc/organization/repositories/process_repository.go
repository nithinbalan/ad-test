package repositories

import (
	"jethur/grc/models"
	organizationDTOs "jethur/grc/organization/dtos"
	"jethur/utility"

	"gorm.io/gorm"
)

type ProcessRepositoryInterface interface {
	GetAll(unitID uint, params *organizationDTOs.ProcessListQueryParamsDTO) (*organizationDTOs.ProcessListResponseDTO, error)
	GetByID(id uint) (*models.Process, error)
	Create(process *models.Process) (uint, error)
	Update(process *models.Process, tx ...utility.DatabaseTransaction) error
	Delete(id uint) error
}

type ProcessRepository struct {
	db *gorm.DB
}

func NewProcessRepository(db *gorm.DB) *ProcessRepository {
	return &ProcessRepository{
		db: db,
	}
}

// func (repo *ProcessRepository) GetAll(unitID uint, params *organizationDTOs.ProcessListQueryParamsDTO) (*organizationDTOs.ProcessListResponseDTO, error) {
// 	var processes []*models.Process
// 	var totalCount int64

// 	query := repo.db.Table("processes").Where("owner_unit_id = ?", unitID)

// 	// Dynamic filters, target date filters, and quick filters implementation would be similar to ControlActionPlanRepository
// 	// Implement additional filtering logic as required by your application here

// 	// Calculate total count for pagination
// 	if err := query.Count(&totalCount).Error; err != nil {
// 		return nil, err
// 	}

// 	// Apply search, sort, and pagination
// 	if params.Keyword != "" {
// 		keyword := "%" + params.Keyword + "%"
// 		query = query.Where("title_en LIKE ? OR title_ar LIKE ?", keyword, keyword)
// 	}

// 	if params.SortBy == "" {
// 		params.SortBy = "title_en" // Default sort by title_en if no sortBy is specified
// 	}

// 	query = query.Order(params.SortBy + " " + params.Sort).Offset((params.Page - 1) * params.PageSize).Limit(params.PageSize)

// 	// Execute the query
// 	err := query.Find(&processes).Error
// 	if err != nil {
// 		return nil, err
// 	}

// 	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

// 	response := &organizationDTOs.ProcessListResponseDTO{
// 		Items: processes,
// 		Pagination: commonDTOs.PaginationInfoDTO{
// 			CurrentPage: params.Page,
// 			PageSize:    params.PageSize,
// 			TotalPages:  totalPages,
// 			TotalCount:  totalCount,
// 		},
// 	}
// 	return response, nil
// }

// func (repo *ProcessRepository) GetByID(id uint) (*models.Process, error) {
// 	var process *models.Process
// 	result := repo.db.Preload("Owner").Preload("LifecycleStatus").First(&process, id)
// 	if result.Error != nil {
// 		return nil, result.Error
// 	}
// 	return process, nil
// }

func (repo *ProcessRepository) Create(process *models.Process) (uint, error) {
	if err := repo.db.Create(process).Error; err != nil {
		return 0, err
	}
	return process.ID, nil
}

func (repo *ProcessRepository) Update(process *models.Process, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(process).Error; err != nil {
		return err
	}
	return nil
}

func (repo *ProcessRepository) Delete(id uint) error {
	result := repo.db.Delete(&models.Process{}, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
