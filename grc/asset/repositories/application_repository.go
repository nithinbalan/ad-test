package repositories

import (
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type ApplicationRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.Application, error)
	GetByID(id uint) (*models.Application, error)
	Create(application *models.Application) (*models.Application, error)
	Update(id uint, application *models.Application) (*models.Application, error)
	Delete(id uint) error
	Activate(ids []uint) ([]uint, error)
	Deactivate(ids []uint) ([]uint, error)
	ExistsByTitle(titleEn string, titleAr string) (bool, error)
	ExistsByTitleExcludingID(id uint, titleEn string, titleAr string) (bool, error)
}

type ApplicationRepository struct {
	db *gorm.DB
}

func NewApplicationRepository(db *gorm.DB) *ApplicationRepository {
	return &ApplicationRepository{
		db: db,
	}
}

func (repo *ApplicationRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.Application, error) {
	var applications []*models.Application
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

	result := query.Order(sortBy + " " + sort).Find(&applications)

	return applications, result.Error
}

func (repo *ApplicationRepository) GetByID(id uint) (*models.Application, error) {
	var application models.Application

	if err := repo.db.First(&application, id).Error; err != nil {
		return nil, err
	}

	return &application, nil
}

func (repo *ApplicationRepository) Create(application *models.Application) (*models.Application, error) {
	if err := repo.db.Create(application).Error; err != nil {
		return nil, err
	}

	return application, nil
}

func (repo *ApplicationRepository) Update(id uint, updatedApplication *models.Application) (*models.Application, error) {
	var application models.Application

	if err := repo.db.First(&application, id).Error; err != nil {
		return nil, err
	}

	if err := repo.db.Model(&application).Updates(updatedApplication).Error; err != nil {
		return nil, err
	}

	return &application, nil
}

func (repo *ApplicationRepository) Delete(id uint) error {
	var application models.Application

	if err := repo.db.First(&application, id).Error; err != nil {
		return err
	}

	return repo.db.Delete(&application).Error
}

func (repo *ApplicationRepository) Activate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint

	if err := repo.db.Model(&models.Application{}).
		Where("id IN ? AND status = ?", ids, "inactive").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err
	}

	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.Application{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "active").Error; err != nil {
			return nil, err
		}
	}

	return idsToUpdate, nil
}

func (repo *ApplicationRepository) Deactivate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint

	if err := repo.db.Model(&models.Application{}).
		Where("id IN ? AND status = ?", ids, "active").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err
	}

	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.Application{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "inactive").Error; err != nil {
			return nil, err
		}
	}

	return idsToUpdate, nil
}

func (repo *ApplicationRepository) ExistsByTitle(titleEn string, titleAr string) (bool, error) {
	var count int64

	err := repo.db.Model(&models.Application{}).
		Where("title_en = ? OR title_ar = ?", titleEn, titleAr).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)
	}

	return count > 0, nil
}

func (repo *ApplicationRepository) ExistsByTitleExcludingID(id uint, titleEn string, titleAr string) (bool, error) {
	var count int64

	err := repo.db.Model(&models.Application{}).
		Where("id <> ? AND (title_en = ? OR title_ar = ?)", id, titleEn, titleAr).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)
	}

	return count > 0, nil
}