package repositories

import (
	commonDTOs "jethur/grc/common/dtos"
	libraryDTOs "jethur/grc/libraries/dtos"
	"jethur/grc/libraries/models"
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
            kri_library.*,
			kri_categories.title_en AS kri_category_title_en,
			kri_categories.title_ar AS kri_category_title_ar,
			review_frequencies.title_en AS review_frequency_title_en,
			review_frequencies.title_ar AS review_frequency_title_ar
        `).
		Joins(`
		    INNER JOIN kri_categories ON kri_categories.id = kri_library.kri_category_id
		    LEFT JOIN review_frequencies ON review_frequencies.id = kri_library.review_frequency_id
		`)

	// Calculate total count for pagination
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	// Apply search keyword if exists
	if params.Keyword != "" {
		query = query.Where(`
			(kri_library.title LIKE ? OR
			kri_library.description LIKE ? OR
			kri_library.calculation_type LIKE ? OR
			kri_library.source LIKE ? OR
			kri_categories.title_en LIKE ? OR
			kri_categories.title_ar LIKE ? OR
			review_frequencies.title_en LIKE ? OR
			review_frequencies.title_ar LIKE ?)`,
			"%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%",
			"%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%")
	}

	// Sort by the specified column and order
	if params.SortBy == "" {
		params.SortBy = "kri_library.title" // Default sort by title if no sortBy is specified
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "asc" // Default order is ascending if not specified
	}

	query = query.Order(params.SortBy + " " + params.Sort)

	offset := (params.Page - 1) * params.PageSize

	// Execute the queryc
	err := query.Offset(offset).Limit(params.PageSize).Find(&kriLibrary).Error
	if err != nil {
		return nil, err
	}

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

func (repo *KRILibraryRepository) GetByID(id uint) (*models.KRILibrary, error) {
	var kriLibrary *models.KRILibrary

	result := repo.db.
		Preload("KRICategory").
		Preload("ReviewFrequency").
		First(&kriLibrary, id)

	if result.Error != nil {
		return nil, result.Error
	}
	return kriLibrary, nil

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
