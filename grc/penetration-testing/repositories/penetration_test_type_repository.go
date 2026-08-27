package repositories

import (
	"errors"
	"gorm.io/gorm"
	"jethur/grc/models"
	"jethur/utility"
)

type PenetrationTestTypeRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.PenetrationTestType, error)
	GetByID(id uint) (*models.PenetrationTestType, error)
	Create(penetrationTestType *models.PenetrationTestType) (*models.PenetrationTestType, error)
	Update(id uint, penetrationTestType *models.PenetrationTestType) (*models.PenetrationTestType, error)
	Delete(id uint) error
	Activate(ids []uint) ([]uint, error)
	Deactivate(ids []uint) ([]uint, error)
	ExistsByTitle(titleEn string, titleAr string) (bool, error)
	ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error)
	ExistsByCode(code string) (bool, error)
	ExistsByCodeExcludingID(id uint, code string) (bool, error)
}

type PenetrationTestTypeRepository struct{ db *gorm.DB }

func NewPenetrationTestTypeRepository(db *gorm.DB) *PenetrationTestTypeRepository {
	return &PenetrationTestTypeRepository{db: db}
}

func (repo *PenetrationTestTypeRepository) GetAll(sort, sortBy, keyword, status string) ([]*models.PenetrationTestType, error) {
	var values []*models.PenetrationTestType
	query := repo.db
	if keyword != "" {
		query = query.Where("(title_en LIKE ? OR title_ar LIKE ? OR code LIKE ?)", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
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
	return values, query.Order(sortBy + " " + sort).Find(&values).Error
}

func (repo *PenetrationTestTypeRepository) GetByID(id uint) (*models.PenetrationTestType, error) {
	var value *models.PenetrationTestType
	result := repo.db.First(&value, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, id)
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return value, nil
}

func (repo *PenetrationTestTypeRepository) Create(value *models.PenetrationTestType) (*models.PenetrationTestType, error) {
	result := repo.db.Create(&value)
	return value, result.Error
}

func (repo *PenetrationTestTypeRepository) Update(id uint, updated *models.PenetrationTestType) (*models.PenetrationTestType, error) {
	var value *models.PenetrationTestType
	if err := repo.db.First(&value, id).Error; err != nil {
		return nil, err
	}
	if err := repo.db.Model(&value).Updates(updated).Error; err != nil {
		return nil, err
	}
	return value, nil
}

func (repo *PenetrationTestTypeRepository) Delete(id uint) error {
	var value *models.PenetrationTestType
	if err := repo.db.First(&value, id).Error; err != nil {
		return err
	}
	return repo.db.Delete(&models.PenetrationTestType{}, id).Error
}

func (repo *PenetrationTestTypeRepository) Activate(ids []uint) ([]uint, error) {
	return repo.updateStatus(ids, "inactive", "active")
}

func (repo *PenetrationTestTypeRepository) Deactivate(ids []uint) ([]uint, error) {
	return repo.updateStatus(ids, "active", "inactive")
}

func (repo *PenetrationTestTypeRepository) updateStatus(ids []uint, from, to string) ([]uint, error) {
	var updated []uint
	if err := repo.db.Model(&models.PenetrationTestType{}).Where("id IN ? AND status = ?", ids, from).Pluck("id", &updated).Error; err != nil {
		return nil, err
	}
	if len(updated) > 0 {
		if err := repo.db.Model(&models.PenetrationTestType{}).Where("id IN ?", updated).Update("status", to).Error; err != nil {
			return nil, err
		}
	}
	return updated, nil
}

func (repo *PenetrationTestTypeRepository) ExistsByTitle(titleEn, titleAr string) (bool, error) {
	var count int64
	err := repo.db.Model(&models.PenetrationTestType{}).Where("title_en = ? OR title_ar = ?", titleEn, titleAr).Count(&count).Error
	if err != nil {
		return false, utility.NewInternalServerError(err)
	}
	return count > 0, nil
}

func (repo *PenetrationTestTypeRepository) ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error) {
	var count int64
	err := repo.db.Model(&models.PenetrationTestType{}).Where("id <> ? AND (title_en = ? OR title_ar = ?)", id, titleEn, titleAr).Count(&count).Error
	if err != nil {
		return false, utility.NewInternalServerError(err)
	}
	return count > 0, nil
}

func (repo *PenetrationTestTypeRepository) ExistsByCode(code string) (bool, error) {
	var count int64
	err := repo.db.Model(&models.PenetrationTestType{}).Where("code = ?", code).Count(&count).Error
	if err != nil {
		return false, utility.NewInternalServerError(err)
	}
	return count > 0, nil
}

func (repo *PenetrationTestTypeRepository) ExistsByCodeExcludingID(id uint, code string) (bool, error) {
	var count int64
	err := repo.db.Model(&models.PenetrationTestType{}).Where("id <> ? AND code = ?", id, code).Count(&count).Error
	if err != nil {
		return false, utility.NewInternalServerError(err)
	}
	return count > 0, nil
}
