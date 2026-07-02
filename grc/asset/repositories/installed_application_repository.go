package repositories

import (
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type InstalledApplicationRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.InstalledApplication, error)
	GetByID(id uint) (*models.InstalledApplication, error)
	Create(installedApplication *models.InstalledApplication) (*models.InstalledApplication, error)
	Update(id uint, installedApplication *models.InstalledApplication) (*models.InstalledApplication, error)
	Delete(id uint) error
	Activate(ids []uint) ([]uint, error)
	Deactivate(ids []uint) ([]uint, error)
	ExistsByTitle(titleEn string, titleAr string) (bool, error)
	ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error)
}

type InstalledApplicationRepository struct {
	db *gorm.DB
}

func NewInstalledApplicationRepository(db *gorm.DB) *InstalledApplicationRepository {
	return &InstalledApplicationRepository{
		db: db,
	}
}

func (repo *InstalledApplicationRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.InstalledApplication, error) {
	var installedApplications []*models.InstalledApplication
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

	result := query.Order(sortBy + " " + sort).Find(&installedApplications)

	return installedApplications, result.Error
}

func (repo *InstalledApplicationRepository) GetByID(id uint) (*models.InstalledApplication, error) {
	var installedApplication models.InstalledApplication

	if err := repo.db.First(&installedApplication, id).Error; err != nil {
		return nil, err
	}

	return &installedApplication, nil
}

func (repo *InstalledApplicationRepository) Create(installedApplication *models.InstalledApplication) (*models.InstalledApplication, error) {
	if err := repo.db.Create(installedApplication).Error; err != nil {
		return nil, err
	}

	return installedApplication, nil
}

func (repo *InstalledApplicationRepository) Update(id uint, updatedApplication *models.InstalledApplication) (*models.InstalledApplication, error) {
	var installedApplication models.InstalledApplication

	if err := repo.db.First(&installedApplication, id).Error; err != nil {
		return nil, err
	}

	if err := repo.db.Model(&installedApplication).Updates(updatedApplication).Error; err != nil {
		return nil, err
	}

	return &installedApplication, nil
}

func (repo *InstalledApplicationRepository) Delete(id uint) error {
	var installedApplication models.InstalledApplication

	if err := repo.db.First(&installedApplication, id).Error; err != nil {
		return err
	}

	return repo.db.Delete(&installedApplication).Error
}

func (repo *InstalledApplicationRepository) Activate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint

	if err := repo.db.Model(&models.InstalledApplication{}).
		Where("id IN ? AND status = ?", ids, "inactive").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err
	}

	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.InstalledApplication{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "active").Error; err != nil {
			return nil, err
		}
	}

	return idsToUpdate, nil
}

func (repo *InstalledApplicationRepository) Deactivate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint

	if err := repo.db.Model(&models.InstalledApplication{}).
		Where("id IN ? AND status = ?", ids, "active").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err
	}

	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.InstalledApplication{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "inactive").Error; err != nil {
			return nil, err
		}
	}

	return idsToUpdate, nil
}

func (repo *InstalledApplicationRepository) ExistsByTitle(titleEn, titleAr string) (bool, error) {
	var count int64

	err := repo.db.Model(&models.InstalledApplication{}).
		Where("title_en = ? OR title_ar = ?", titleEn, titleAr).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)
	}
	return count > 0, nil
}

func (repo *InstalledApplicationRepository) ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error) {
	var count int64
	err := repo.db.Model(&models.InstalledApplication{}).
		Where("id <> ? AND (title_en = ? OR title_ar = ?)", id, titleEn, titleAr).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)
	}

	return count > 0, nil
}