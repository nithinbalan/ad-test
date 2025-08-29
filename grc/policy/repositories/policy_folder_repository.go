package repositories

import (
	"jethur/grc/models"
	policyDTOs "jethur/grc/policy/dtos"
	"jethur/utility"
	"strings"

	"gorm.io/gorm"
)

// Interfaces
type PolicyFolderRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*policyDTOs.PolicyFolderListDTO, error)
	GetByID(id uint) (*models.PolicyFolder, error)
	Create(folder *models.PolicyFolder) (*models.PolicyFolder, error)
	Update(id uint, folder *models.PolicyFolder) (*models.PolicyFolder, error)
	Delete(id uint) error
	Activate(ids []uint) ([]uint, error)
	Deactivate(ids []uint) ([]uint, error)
	ExistsByTitle(titleEn string, titleAr string) (bool, error)
	ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error)
}

// Implementation
type PolicyFolderRepository struct {
	db *gorm.DB
}

func NewPolicyFolderRepository(db *gorm.DB) *PolicyFolderRepository {
	return &PolicyFolderRepository{
		db: db,
	}
}

func (repo *PolicyFolderRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*policyDTOs.PolicyFolderListDTO, error) {
	var folders []*policyDTOs.PolicyFolderListDTO

	query := `
		SELECT 
			cf.id, 
			cf.title_en, 
			cf.title_ar, 
			cf.status, 
			COALESCE(cf.created_at, '1900-01-01 00:00:00') AS created_at, 
			COUNT(cr.id) AS policy_register_count
		FROM policy_folders cf
		LEFT JOIN documents cr ON cr.policy_folder_id = cf.id
	`

	// Apply filters dynamically
	whereClauses := []string{}
	args := []interface{}{}

	if keyword != "" {
		whereClauses = append(whereClauses, "(cf.title_en LIKE ? OR cf.title_ar LIKE ?)")
		args = append(args, "%"+keyword+"%", "%"+keyword+"%")
	}

	if status == "active" || status == "inactive" {
		whereClauses = append(whereClauses, "cf.status = ?")
		args = append(args, status)
	}

	// Construct WHERE clause
	if len(whereClauses) > 0 {
		query += " WHERE " + strings.Join(whereClauses, " AND ")
	}

	// Add created_at to GROUP BY to resolve the error
	query += " GROUP BY cf.id, cf.title_en, cf.title_ar, cf.status, cf.created_at"

	// Apply sorting
	if sortBy == "" {
		sortBy = "cf.created_at"
	}
	if sort != "asc" && sort != "desc" {
		sort = "asc"
	}
	query += " ORDER BY " + sortBy + " " + sort

	// Execute query
	err := repo.db.Raw(query, args...).Scan(&folders).Error

	// Ensure folders is never nil
	if folders == nil {
		folders = []*policyDTOs.PolicyFolderListDTO{}
	}

	return folders, err
}

func (repo *PolicyFolderRepository) GetByID(id uint) (*models.PolicyFolder, error) {
	var folder *models.PolicyFolder
	result := repo.db.First(&folder, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return folder, nil
}

func (repo *PolicyFolderRepository) Create(folder *models.PolicyFolder) (*models.PolicyFolder, error) {
	result := repo.db.Create(&folder)

	return folder, result.Error
}

func (repo *PolicyFolderRepository) Update(id uint, updatedFolder *models.PolicyFolder) (*models.PolicyFolder, error) {
	var folder *models.PolicyFolder

	if err := repo.db.First(&folder, id).Error; err != nil {
		return nil, err
	}

	if err := repo.db.Model(&folder).Updates(updatedFolder).Error; err != nil {
		return nil, err
	}

	return folder, nil
}

func (repo *PolicyFolderRepository) Delete(id uint) error {
	var folder *models.PolicyFolder
	if err := repo.db.First(&folder, id).Error; err != nil {
		return err
	}

	result := repo.db.Delete(&models.PolicyFolder{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *PolicyFolderRepository) Activate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint

	// Step 1: Fetch the IDs of records that are currently 'inactive' and within the provided 'ids' slice
	if err := repo.db.Model(&models.PolicyFolder{}).
		Where("id IN ? AND status = ?", ids, "inactive").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err
	}

	// Step 2: If there are any IDs to update, perform the update operation
	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.PolicyFolder{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "active").Error; err != nil {
			return nil, err
		}
	}

	// Return the IDs of the records that were activated
	return idsToUpdate, nil
}

func (repo *PolicyFolderRepository) Deactivate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint

	// Step 1: Fetch the IDs of records that are currently 'active' and within the provided 'ids' slice
	if err := repo.db.Model(&models.PolicyFolder{}).
		Where("id IN ? AND status = ?", ids, "active").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err
	}

	// Step 2: If there are any IDs to update, perform the update operation
	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.PolicyFolder{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "inactive").Error; err != nil {
			return nil, err
		}
	}

	// Return the IDs of the records that were deactivated
	return idsToUpdate, nil
}

func (repo *PolicyFolderRepository) ExistsByTitle(titleEn, titleAr string) (bool, error) {
	var count int64

	// Check for the existence of titleEn or titleAr in a single query
	err := repo.db.Model(&models.PolicyFolder{}).
		Where("title_en = ? OR title_ar = ?", titleEn, titleAr).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)
	}

	// Return true if either title exists, otherwise false
	return count > 0, nil
}

func (repo *PolicyFolderRepository) ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error) {
	var count int64

	err := repo.db.Model(&models.PolicyFolder{}).
		Where("id <> ? AND (title_en = ? OR title_ar = ?)", id, titleEn, titleAr).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)
	}

	return count > 0, nil
}
