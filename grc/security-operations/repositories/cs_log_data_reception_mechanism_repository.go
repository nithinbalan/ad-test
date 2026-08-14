package repositories

import (
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type CSLogDataReceptionMechanismRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.CSLogDataReceptionMechanism, error)
	GetByID(id uint) (*models.CSLogDataReceptionMechanism, error)
	Create(mechanism *models.CSLogDataReceptionMechanism) (*models.CSLogDataReceptionMechanism, error)
	Update(id uint, mechanism *models.CSLogDataReceptionMechanism) (*models.CSLogDataReceptionMechanism, error)
	Delete(id uint) error
	Activate(ids []uint) ([]uint, error)
	Deactivate(ids []uint) ([]uint, error)
	ExistsByTitle(titleEn string, titleAr string) (bool, error)
	ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error)
}

type CSLogDataReceptionMechanismRepository struct {
	db *gorm.DB
}

func NewCSLogDataReceptionMechanismRepository(db *gorm.DB) *CSLogDataReceptionMechanismRepository {
	return &CSLogDataReceptionMechanismRepository{
		db: db,
	}
}

func (repo *CSLogDataReceptionMechanismRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.CSLogDataReceptionMechanism, error) {
	var mechanisms []*models.CSLogDataReceptionMechanism
	query := repo.db

	if keyword != "" {
		query = query.Where("(title_en LIKE ? OR title_ar LIKE ?)", "%"+keyword+"%", "%"+keyword+"%")
	}
	if sortBy == "" {
		sortBy = "id"
	}
	if sort != "asc" && sort != "desc" {
		sort = "desc"
	}

	if status == "active" || status == "inactive" {
		query = query.Where("status = ?", status)
	}

	query = query.Order(sortBy + " " + sort)
	result := query.Find(&mechanisms)

	return mechanisms, result.Error
}

func (repo *CSLogDataReceptionMechanismRepository) GetByID(id uint) (*models.CSLogDataReceptionMechanism, error) {
	var mechanism *models.CSLogDataReceptionMechanism
	result := repo.db.First(&mechanism, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return mechanism, nil
}

func (repo *CSLogDataReceptionMechanismRepository) Create(mechanism *models.CSLogDataReceptionMechanism) (*models.CSLogDataReceptionMechanism, error) {
	result := repo.db.Create(&mechanism)

	return mechanism, result.Error
}

func (repo *CSLogDataReceptionMechanismRepository) Update(id uint, updated *models.CSLogDataReceptionMechanism) (*models.CSLogDataReceptionMechanism, error) {
	var mechanism *models.CSLogDataReceptionMechanism
	if err := repo.db.First(&mechanism, id).Error; err != nil {
		return nil, err
	}

	if err := repo.db.Model(&mechanism).Updates(updated).Error; err != nil {
		return nil, err
	}

	return mechanism, nil
}

func (repo *CSLogDataReceptionMechanismRepository) Delete(id uint) error {
	var mechanism *models.CSLogDataReceptionMechanism
	if err := repo.db.First(&mechanism, id).Error; err != nil {
		return err
	}

	result := repo.db.Delete(&models.CSLogDataReceptionMechanism{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *CSLogDataReceptionMechanismRepository) Activate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint
	if err := repo.db.Model(&models.CSLogDataReceptionMechanism{}).
		Where("id IN ? AND status = ?", ids, "inactive").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err
	}

	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.CSLogDataReceptionMechanism{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "active").Error; err != nil {
			return nil, err
		}
	}
	return idsToUpdate, nil
}

func (repo *CSLogDataReceptionMechanismRepository) Deactivate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint
	if err := repo.db.Model(&models.CSLogDataReceptionMechanism{}).
		Where("id IN ? AND status = ?", ids, "active").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err
	}

	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.CSLogDataReceptionMechanism{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "inactive").Error; err != nil {
			return nil, err
		}
	}
	return idsToUpdate, nil
}

func (repo *CSLogDataReceptionMechanismRepository) ExistsByTitle(titleEn, titleAr string) (bool, error) {
	var count int64
	err := repo.db.Model(&models.CSLogDataReceptionMechanism{}).
		Where("title_en = ? OR title_ar = ?", titleEn, titleAr).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)
	}
	return count > 0, nil
}

func (repo *CSLogDataReceptionMechanismRepository) ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error) {
	var count int64
	err := repo.db.Model(&models.CSLogDataReceptionMechanism{}).
		Where("id <> ? AND (title_en = ? OR title_ar = ?)", id, titleEn, titleAr).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)
	}
	return count > 0, nil
}
