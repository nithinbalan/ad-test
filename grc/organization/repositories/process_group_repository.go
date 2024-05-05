package repositories

import (
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interfaces
type ProcessGroupRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ProcessGroup, error)
	GetByID(id uint) (*models.ProcessGroup, error)
	Create(processGroup *models.ProcessGroup) (*models.ProcessGroup, error)
	Update(id uint, updatedProcessGroup *models.ProcessGroup) (*models.ProcessGroup, error)
	Delete(id uint) error
	Activate(ids []uint) ([]uint, error)
	Deactivate(ids []uint) ([]uint, error)
	ExistsByTitle(titleEn string, titleAr string) (bool, error)
	ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error)
}

// implementation
type ProcessGroupRepository struct {
	db *gorm.DB
}

func NewProcessGroupRepository(db *gorm.DB) *ProcessGroupRepository {
	return &ProcessGroupRepository{
		db: db,
	}
}

func (repo *ProcessGroupRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ProcessGroup, error) {
	var processGroups []*models.ProcessGroup

	query := repo.db

	if keyword != "" {
		query = query.Where("(title_en LIKE ? OR title_ar LIKE ?)", "%"+keyword+"%", "%"+keyword+"%")
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
	result := query.Find(&processGroups)

	return processGroups, result.Error
}

func (repo *ProcessGroupRepository) GetByID(id uint) (*models.ProcessGroup, error) {
	var processGroup *models.ProcessGroup
	result := repo.db.First(&processGroup, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return processGroup, nil
}

func (repo *ProcessGroupRepository) Create(processGroup *models.ProcessGroup) (*models.ProcessGroup, error) {
	result := repo.db.Create(&processGroup)

	return processGroup, result.Error
}

func (repo *ProcessGroupRepository) Update(id uint, updatedProcessGroup *models.ProcessGroup) (*models.ProcessGroup, error) {
	var processGroup *models.ProcessGroup

	if err := repo.db.First(&processGroup, id).Error; err != nil {
		return nil, err
	}

	if err := repo.db.Model(&processGroup).Updates(updatedProcessGroup).Error; err != nil {
		return nil, err
	}

	return processGroup, nil
}

func (repo *ProcessGroupRepository) Delete(id uint) error {
	var processGroup *models.ProcessGroup
	if err := repo.db.First(&processGroup, id).Error; err != nil {
		return err
	}

	result := repo.db.Delete(&models.ProcessGroup{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *ProcessGroupRepository) Activate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint

	if err := repo.db.Model(&models.ProcessGroup{}).
		Where("id IN ? AND status = ?", ids, "inactive").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err
	}

	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.ProcessGroup{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "active").Error; err != nil {
			return nil, err
		}
	}

	return idsToUpdate, nil
}

func (repo *ProcessGroupRepository) Deactivate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint

	if err := repo.db.Model(&models.ProcessGroup{}).
		Where("id IN ? AND status = ?", ids, "active").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err
	}

	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.ProcessGroup{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "inactive").Error; err != nil {
			return nil, err
		}
	}

	return idsToUpdate, nil

}

func (repo *ProcessGroupRepository) ExistsByTitle(titleEn, titleAr string) (bool, error) {
	var count int64

	err := repo.db.Model(&models.ProcessGroup{}).
		Where("title_en = ? OR title_ar = ?", titleEn, titleAr).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)

	}

	return count > 0, nil
}

func (repo *ProcessGroupRepository) ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error) {
	var count int64

	err := repo.db.Model(&models.ProcessGroup{}).
		Where("id <> ? AND (title_en = ? OR title_ar = ?)", id, titleEn, titleAr).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)
	}

	return count > 0, nil
}
