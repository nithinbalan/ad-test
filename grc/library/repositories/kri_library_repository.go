package repositories

import (
	"errors"

	commonDTOs "jethur/grc/common/dtos"
	libraryDTOs "jethur/grc/library/dtos"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type KRILibraryRepositoryInterface interface {
	GetAll(params *libraryDTOs.KRILibraryListQueryParamsDTO) (*libraryDTOs.KRILibraryListResponseDTO, error)
	GetByID(id uint) (*models.KRILibrary, error)
	Create(kriLibrary *models.KRILibrary) (uint, error)
	Update(kriLibrary *models.KRILibrary) error
	Delete(id uint) error
	Activate(ids []uint) ([]uint, error)
	Deactivate(ids []uint) ([]uint, error)
	ExistsByTitle(Title string, KRICategoryID uint) (bool, error)
	ExistsByTitleExcludingID(id uint, Title string, KRICategoryID uint) (bool, error)
	GetKRICategoryCounts() ([]libraryDTOs.KRICategoryCountDTO, error)
}

type KRILibraryRepository struct {
	db *gorm.DB
}

func NewKRILibraryRepository(db *gorm.DB) *KRILibraryRepository {
	return &KRILibraryRepository{
		db: db,
	}
}

func (repo *KRILibraryRepository) GetAll(params *libraryDTOs.KRILibraryListQueryParamsDTO) (*libraryDTOs.KRILibraryListResponseDTO, error) {
	var kriLibrary []*libraryDTOs.KRILibraryDTO
	var totalCount int64

	query := repo.db.Table("kri_library").
		Select(`
			id,
			kri_category,
			title,
			description,
			kri_unit
		`)

	// Filter by keyword (title, description, unit, category)
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
			title LIKE ? OR
			description LIKE ? OR
			kri_unit LIKE ? OR
			kri_category LIKE ?
		`, keyword, keyword, keyword, keyword)
	}

	// Filter by category
	if len(params.Category) > 0 {
		query = query.Where("kri_category IN ?", params.Category)
	}

	// Count total for pagination
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	// Sorting
	sortBy := params.SortBy
	if sortBy == "" {
		sortBy = "title"
	}
	sortOrder := params.Sort
	if sortOrder != "desc" {
		sortOrder = "asc"
	}
	query = query.Order(sortBy + " " + sortOrder)

	// Pagination
	if !params.Export {
		offset := (params.Page - 1) * params.PageSize
		query = query.Offset(offset).Limit(params.PageSize)
	}

	err := query.Find(&kriLibrary).Error
	if err != nil {
		return nil, err
	}

	// Response
	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))
	response := &libraryDTOs.KRILibraryListResponseDTO{
		Items: kriLibrary,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}

	return response, nil
}

func (repo *KRILibraryRepository) GetKRICategoryCounts() ([]libraryDTOs.KRICategoryCountDTO, error) {
	var results []libraryDTOs.KRICategoryCountDTO

	err := repo.db.
		Table("kri_library").
		Select("kri_category, COUNT(*) as count").
		Group("kri_category").
		Order("kri_category ASC").
		Scan(&results).Error
	if err != nil {
		return nil, err
	}

	return results, nil
}

func (repo *KRILibraryRepository) GetByID(id uint) (*models.KRILibrary, error) {
	var kriLibrary *models.KRILibrary

	result := repo.db.
		Preload("KRICategory").
		Preload("ReviewFrequency").
		First(&kriLibrary, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, id)
	} else if result.Error != nil {
		return nil, result.Error
	}
	return kriLibrary, result.Error
}

func (repo *KRILibraryRepository) Create(kriLibrary *models.KRILibrary) (uint, error) {
	if err := repo.db.Create(kriLibrary).Error; err != nil {
		return 0, err
	}

	return kriLibrary.ID, nil
}

func (repo *KRILibraryRepository) Update(kriLibrary *models.KRILibrary) error {
	if err := repo.db.Save(kriLibrary).Error; err != nil {
		return err
	}

	return nil
}

func (repo *KRILibraryRepository) Delete(id uint) error {
	var kriLibrary *models.KRILibrary
	if err := repo.db.First(&kriLibrary, id).Error; err != nil {
		return err
	}

	result := repo.db.Delete(&models.KRILibrary{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *KRILibraryRepository) Activate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint

	// Step 1: Fetch the IDs of records that are currently 'inactive' and within the provided 'ids' slice
	if err := repo.db.Model(&models.KRILibrary{}).
		Where("id IN ? AND status = ?", ids, "inactive").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err // Return error if the query fails
	}

	// Step 2: If there are any IDs to update, perform the update operation
	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.KRILibrary{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "active").Error; err != nil {
			return nil, err // Return error if the update operation fails
		}
	}

	// Return the IDs of the records that were activated
	return idsToUpdate, nil
}

func (repo *KRILibraryRepository) Deactivate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint

	// Step 1: Fetch the IDs of records that are currently 'inactive' and within the provided 'ids' slice
	if err := repo.db.Model(&models.KRILibrary{}).
		Where("id IN ? AND status = ?", ids, "active").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err // Return error if the query fails
	}

	// Step 2: If there are any IDs to update, perform the update operation
	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.KRILibrary{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "inactive").Error; err != nil {
			return nil, err // Return error if the update operation fails
		}
	}

	// Return the IDs of the records that were activated
	return idsToUpdate, nil
}

func (repo *KRILibraryRepository) ExistsByTitle(Title string, KRICategoryID uint) (bool, error) {
	var count int64

	// Check for the existence of Title or QuestionAr in a single query
	err := repo.db.Model(&models.KRILibrary{}).
		Where("title = ? AND kri_category_id = ?", Title, KRICategoryID).
		Count(&count).Error
	if err != nil {
		return false, utility.NewInternalServerError(err)
	}
	// Return true if either title exists, otherwise false
	return count > 0, nil
}

func (repo *KRILibraryRepository) ExistsByTitleExcludingID(id uint, Title string, KRICategoryID uint) (bool, error) {
	var count int64
	err := repo.db.Model(&models.KRILibrary{}).
		Where("id <> ? AND title = ? AND kri_category_id = ?", id, Title, KRICategoryID).
		Count(&count).Error
	if err != nil {
		return false, utility.NewInternalServerError(err)
	}

	return count > 0, nil
}
