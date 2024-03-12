package repositories

import (
	"errors"
	"jethur/grc/compliance/models"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interfaces
type SLARegisterRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string) ([]*models.SLARegister, error)
	GetByID(id uint) (*models.SLARegister, error)
	Create(article *models.SLARegister) (*models.SLARegister, error)
	Update(id uint, slaRegister *models.SLARegister) (*models.SLARegister, error)
	Delete(id uint) error
	ExistsByTitle(titleEn string, titleAr string) (bool, error)
	ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error)
}

// implementation
type SLARegisterRepository struct {
	db *gorm.DB
}

func NewSLARegisterRepository(db *gorm.DB) *SLARegisterRepository {
	return &SLARegisterRepository{
		db: db,
	}
}

func (repo *SLARegisterRepository) GetAll(sort string, sortBy string, keyword string) ([]*models.SLARegister, error) {

	var slaRegisters []*models.SLARegister

	query := repo.db

	if keyword != "" {
		query = query.Where("title_en LIKE ? OR title_ar LIKE ? OR description LIKE ? OR client LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}
	if sortBy == "" {
		sortBy = "title_en"
	}
	if sort != "asc" && sort != "desc" {
		sort = "asc"
	}

	// Preload relationships
	result := query.Order(sortBy + " " + sort).
		Preload("SLACategory").
		Preload("SeverityLevel").
		Preload("ResponsibleUnit").
		Preload("ResponsibleUser").
		Preload("LifecycleStatus").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		Find(&slaRegisters)
	return slaRegisters, result.Error

}

func (repo *SLARegisterRepository) GetByID(id uint) (*models.SLARegister, error) {
	var slaRegister *models.SLARegister
	// result := repo.db.First(&slaRegister, id)
	result := repo.db.Preload("SLACategory").
		Preload("SeverityLevel").
		Preload("ResponsibleUnit").
		Preload("ResponsibleUser").
		Preload("LifecycleStatus").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		First(&slaRegister, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, id)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return slaRegister, result.Error
}

func (repo *SLARegisterRepository) Create(slaRegister *models.SLARegister) (*models.SLARegister, error) {
	result := repo.db.Create(&slaRegister)

	return slaRegister, result.Error
}

func (repo *SLARegisterRepository) Update(id uint, updatedRegister *models.SLARegister) (*models.SLARegister, error) {
	var slaRegister *models.SLARegister

	if err := repo.db.First(&slaRegister, id).Error; err != nil {
		return nil, err
	}

	if err := repo.db.Model(&slaRegister).Updates(updatedRegister).Error; err != nil {
		return nil, err
	}

	return slaRegister, nil
}

func (repo *SLARegisterRepository) Delete(id uint) error {

	var slaRegister *models.SLARegister
	if err := repo.db.First(&slaRegister, id).Error; err != nil {
		return err
	}

	result := repo.db.Delete(&models.SLARegister{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}


func (repo *SLARegisterRepository) ExistsByTitle(titleEn, titleAr string) (bool, error) {
	var count int64

	// Check for the existence of titleEn or titleAr in a single query
	err := repo.db.Model(&models.SLARegister{}).
		Where("title_en = ? OR title_ar = ?", titleEn, titleAr).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)

	}
	// Return true if either title exists, otherwise false
	return count > 0, nil
}


func (repo *SLARegisterRepository) ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error) {
	var count int64

	err := repo.db.Model(&models.SLARegister{}).
		Where("id <> ? AND (title_en = ? OR title_ar = ?)", id, titleEn, titleAr).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)
	}

	return count > 0, nil
}


