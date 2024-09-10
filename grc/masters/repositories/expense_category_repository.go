package repositories

import (
	"errors"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interfaces
type ExpenseCategoryRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ExpenseCategory, error)
	GetByID(id uint) (*models.ExpenseCategory, error)
	Create(article *models.ExpenseCategory) (*models.ExpenseCategory, error)
	Update(id uint, expenseCategory *models.ExpenseCategory) (*models.ExpenseCategory, error)
	Delete(id uint) error
	Activate(ids []uint) ([]uint, error)
	Deactivate(ids []uint) ([]uint, error)
	ExistsByTitle(titleEn string, titleAr string) (bool, error)
	ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error)
}

// implementation
type ExpenseCategoryRepository struct {
	db *gorm.DB
}

func NewExpenseCategoryRepository(db *gorm.DB) *ExpenseCategoryRepository {
	return &ExpenseCategoryRepository{
		db: db,
	}
}

func (repo *ExpenseCategoryRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ExpenseCategory, error) {

	var expenseCategories []*models.ExpenseCategory

	query := repo.db

	if keyword != "" {
		query = query.Where("(title_en LIKE ? OR title_ar LIKE  ? OR description_en LIKE  ? OR description_ar LIKE  ?)", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}
	if sortBy == "" {
		sortBy = "title_en"
	}
	if sort != "asc" && sort != "desc" {
		sort = "asc"
	}

	if status == "active" || status == "inactive" {
		query = query.Where("status = ?", status)
	}

	query = query.Order(sortBy + " " + sort)
	result := query.Find(&expenseCategories)

	return expenseCategories, result.Error
}

func (repo *ExpenseCategoryRepository) GetByID(id uint) (*models.ExpenseCategory, error) {
	var expenseCategory *models.ExpenseCategory
	result := repo.db.First(&expenseCategory, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, id)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return expenseCategory, result.Error
}

func (repo *ExpenseCategoryRepository) Create(expenseCategory *models.ExpenseCategory) (*models.ExpenseCategory, error) {
	result := repo.db.Create(&expenseCategory)

	return expenseCategory, result.Error
}

func (repo *ExpenseCategoryRepository) Update(id uint, updatedExpenseCategory *models.ExpenseCategory) (*models.ExpenseCategory, error) {
	var expenseCategory *models.ExpenseCategory

	if err := repo.db.First(&expenseCategory, id).Error; err != nil {
		return nil, err
	}

	if err := repo.db.Model(&expenseCategory).Updates(updatedExpenseCategory).Error; err != nil {
		return nil, err
	}

	return expenseCategory, nil
}

func (repo *ExpenseCategoryRepository) Delete(id uint) error {

	var expenseCategory *models.ExpenseCategory
	if err := repo.db.First(&expenseCategory, id).Error; err != nil {
		return err
	}

	result := repo.db.Delete(&models.ExpenseCategory{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *ExpenseCategoryRepository) Activate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint

	// Step 1: Fetch the IDs of records that are currently 'inactive' and within the provided 'ids' slice
	if err := repo.db.Model(&models.ExpenseCategory{}).
		Where("id IN ? AND status = ?", ids, "inactive").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err // Return error if the query fails
	}

	// Step 2: If there are any IDs to update, perform the update operation
	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.ExpenseCategory{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "active").Error; err != nil {
			return nil, err // Return error if the update operation fails
		}
	}

	// Return the IDs of the records that were activated
	return idsToUpdate, nil
}

func (repo *ExpenseCategoryRepository) Deactivate(ids []uint) ([]uint, error) {

	var idsToUpdate []uint

	// Step 1: Fetch the IDs of records that are currently 'inactive' and within the provided 'ids' slice
	if err := repo.db.Model(&models.ExpenseCategory{}).
		Where("id IN ? AND status = ?", ids, "active").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err // Return error if the query fails
	}

	// Step 2: If there are any IDs to update, perform the update operation
	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.ExpenseCategory{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "inactive").Error; err != nil {
			return nil, err // Return error if the update operation fails
		}
	}

	// Return the IDs of the records that were activated
	return idsToUpdate, nil

}

func (repo *ExpenseCategoryRepository) ExistsByTitle(titleEn, titleAr string) (bool, error) {
	var count int64

	// Check for the existence of titleEn or titleAr in a single query
	err := repo.db.Model(&models.ExpenseCategory{}).
		Where("title_en = ? OR title_ar = ?", titleEn, titleAr).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)

	}
	// Return true if either title exists, otherwise false
	return count > 0, nil
}

func (repo *ExpenseCategoryRepository) ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error) {
	var count int64

	err := repo.db.Model(&models.ExpenseCategory{}).
		Where("id <> ? AND (title_en = ? OR title_ar = ?)", id, titleEn, titleAr).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)
	}

	return count > 0, nil
}
