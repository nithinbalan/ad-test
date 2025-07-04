package repositories

import (
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/models"
	"jethur/grc/strategy/dtos"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interfaces
type FocusAreaRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction

	GetAll(params *dtos.FocusAreaListQueryParamsDTO) (*dtos.FocusAreaListResponseDTO, error)
	GetByID(id uint) (*models.FocusArea, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.FocusArea, error)
	Create(focusArea *models.FocusArea, tx ...utility.DatabaseTransaction) (uint, error)
	Update(focusArea *models.FocusArea, tx ...utility.DatabaseTransaction) error
	Delete(id uint) error
	Activate(ids []uint) ([]uint, error)
	Deactivate(ids []uint) ([]uint, error)
	ExistsByTitle(titleEn string, titleAr string) (bool, error)
	ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error)
}

// implementation
type FocusAreaRepository struct {
	db *gorm.DB
}

func NewFocusAreaRepository(db *gorm.DB) *FocusAreaRepository {
	return &FocusAreaRepository{
		db: db,
	}
}
func (repo *FocusAreaRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *FocusAreaRepository) GetAll(params *dtos.FocusAreaListQueryParamsDTO) (*dtos.FocusAreaListResponseDTO, error) {
	var focusAreas []*dtos.FocusAreaListDTO
	var totalCount int64

	query := repo.db.Table("focus_areas").
		Select(`
			focus_areas.*,
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
		Joins(`LEFT JOIN users AS created_user ON created_user.id = focus_areas.created_by_id`).
		Joins(`LEFT JOIN users AS updated_user ON updated_user.id = focus_areas.updated_by_id`)

	if params.Keyword != "" {
		kw := "%" + params.Keyword + "%"
		query = query.Where(`
			focus_areas.title_en LIKE ? OR
			focus_areas.title_ar LIKE ?
		`, kw, kw)
	}

	for col, val := range params.Filters {
		if len(val) > 0 {
			query = query.Where(col+" IN (?)", val)
		}
	}

	// Sorting
	query = query.Order(params.SortBy + " " + params.Sort)

	// Pagination
	query.Count(&totalCount)
	offset := (params.Page - 1) * params.PageSize
	query = query.Offset(offset).Limit(params.PageSize)

	if err := query.Find(&focusAreas).Error; err != nil {
		return nil, err
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))
	response := &dtos.FocusAreaListResponseDTO{
		Items: focusAreas,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}
	return response, nil
}

func (repo *FocusAreaRepository) GetByID(id uint) (*models.FocusArea, error) {
	var focusArea *models.FocusArea
	result := repo.db.First(&focusArea, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return focusArea, nil
}

func (repo *FocusAreaRepository) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.FocusArea, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var focusArea *models.FocusArea

	result := db.
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		First(&focusArea, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return focusArea, nil
}

func (repo *FocusAreaRepository) Create(focusArea *models.FocusArea, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Create(focusArea).Error; err != nil {
		return 0, err
	}

	return focusArea.ID, nil
}

func (repo *FocusAreaRepository) Update(focusArea *models.FocusArea, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Save(focusArea).Error; err != nil {
		return err
	}

	return nil
}

func (repo *FocusAreaRepository) Delete(id uint) error {
	var focusArea *models.FocusArea
	if err := repo.db.First(&focusArea, id).Error; err != nil {
		return err
	}

	result := repo.db.Delete(&models.FocusArea{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *FocusAreaRepository) Activate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint

	if err := repo.db.Model(&models.FocusArea{}).
		Where("id IN ? AND status = ?", ids, "inactive").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err
	}

	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.FocusArea{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "active").Error; err != nil {
			return nil, err
		}
	}

	return idsToUpdate, nil
}

func (repo *FocusAreaRepository) Deactivate(ids []uint) ([]uint, error) {

	var idsToUpdate []uint

	if err := repo.db.Model(&models.FocusArea{}).
		Where("id IN ? AND status = ?", ids, "active").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err
	}

	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.FocusArea{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "inactive").Error; err != nil {
			return nil, err
		}
	}

	return idsToUpdate, nil

}

func (repo *FocusAreaRepository) ExistsByTitle(titleEn, titleAr string) (bool, error) {
	var count int64

	err := repo.db.Model(&models.FocusArea{}).
		Where("title_en = ? OR title_ar = ?", titleEn, titleAr).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)

	}

	return count > 0, nil
}

func (repo *FocusAreaRepository) ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error) {
	var count int64

	err := repo.db.Model(&models.FocusArea{}).
		Where("id <> ? AND (title_en = ? OR title_ar = ?)", id, titleEn, titleAr).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)
	}

	return count > 0, nil
}
