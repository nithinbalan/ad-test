package repositories

import (
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type InstalledRunningApplicationRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.InstalledRunningApplication, error)
	GetByID(id uint) (*models.InstalledRunningApplication, error)
	Create(installedRunningApplication *models.InstalledRunningApplication) (*models.InstalledRunningApplication, error)
	Update(id uint, installedRunningApplication *models.InstalledRunningApplication) (*models.InstalledRunningApplication, error)
	Delete(id uint) error
	Activate(ids []uint) ([]uint, error)
	Deactivate(ids []uint) ([]uint, error)
	ExistsByTitle(titleEn string, titleAr string) (bool, error)
	ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error)
}

type InstalledRunningApplicationRepository struct {
	db *gorm.DB
}

func NewInstalledRunningApplicationRepository(db *gorm.DB) *InstalledRunningApplicationRepository {
	return &InstalledRunningApplicationRepository{
		db: db,
	}
}

func (repo *InstalledRunningApplicationRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.InstalledRunningApplication, error) {
	var installedRunningApplications []*models.InstalledRunningApplication
	query := repo.db

	if keyword != "" {
		query = query.Where(
			"(title_en LIKE ? OR title_ar LIKE ?)",
			"%"+keyword+"%",
			"%"+keyword+"%",
		)
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

	result := query.Order(sortBy + " " + sort).Find(&installedRunningApplications)

	return installedRunningApplications, result.Error
}

func (repo *InstalledRunningApplicationRepository) GetByID(id uint) (*models.InstalledRunningApplication, error) {
	var installedRunningApplication models.InstalledRunningApplication

	if err := repo.db.First(&installedRunningApplication, id).Error; err != nil {
		return nil, err
	}

	return &installedRunningApplication, nil
}

func (repo *InstalledRunningApplicationRepository) Create(installedRunningApplication *models.InstalledRunningApplication) (*models.InstalledRunningApplication, error) {
	if err := repo.db.Create(installedRunningApplication).Error; err != nil {
		return nil, err
	}

	return installedRunningApplication, nil
}

func (repo *InstalledRunningApplicationRepository) Update(id uint, updatedApplication *models.InstalledRunningApplication) (*models.InstalledRunningApplication, error) {
	var installedRunningApplication models.InstalledRunningApplication

	if err := repo.db.First(&installedRunningApplication, id).Error; err != nil {
		return nil, err
	}

	if err := repo.db.Model(&installedRunningApplication).Updates(updatedApplication).Error; err != nil {
		return nil, err
	}

	return &installedRunningApplication, nil
}

func (repo *InstalledRunningApplicationRepository) Delete(id uint) error {
	var installedRunningApplication models.InstalledRunningApplication

	if err := repo.db.First(&installedRunningApplication, id).Error; err != nil {
		return err
	}

	return repo.db.Delete(&installedRunningApplication).Error
}

func (repo *InstalledRunningApplicationRepository) Activate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint

	if err := repo.db.Model(&models.InstalledRunningApplication{}).
		Where("id IN ? AND status = ?", ids, "inactive").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err
	}

	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.InstalledRunningApplication{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "active").Error; err != nil {
			return nil, err
		}
	}

	return idsToUpdate, nil
}

func (repo *InstalledRunningApplicationRepository) Deactivate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint

	if err := repo.db.Model(&models.InstalledRunningApplication{}).
		Where("id IN ? AND status = ?", ids, "active").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err
	}

	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.InstalledRunningApplication{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "inactive").Error; err != nil {
			return nil, err
		}
	}

	return idsToUpdate, nil
}

func (repo *InstalledRunningApplicationRepository) ExistsByTitle(titleEn, titleAr string) (bool, error) {
	var count int64

	err := repo.db.Model(&models.InstalledRunningApplication{}).
		Where("title_en = ? OR title_ar = ?", titleEn, titleAr).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)
	}
	return count > 0, nil
}

func (repo *InstalledRunningApplicationRepository) ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error) {
	var count int64
	err := repo.db.Model(&models.InstalledRunningApplication{}).
		Where("id <> ? AND (title_en = ? OR title_ar = ?)", id, titleEn, titleAr).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)
	}

	return count > 0, nil
}