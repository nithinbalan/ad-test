package repositories

import (
	"errors"
	"jethur/grc/compliance/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type ComplianceRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ComplianceRegister, error)
	GetByID(id uint) (*models.ComplianceRegister, error)
	Create(compliance *models.ComplianceRegister) (*models.ComplianceRegister, error)
	Update(id uint, compliance *models.ComplianceRegister) (*models.ComplianceRegister, error)
	Delete(id uint) error
	Activate(ids []uint) ([]uint, error)
	Deactivate(ids []uint) ([]uint, error)
	ExistsByTitle(titleEn string, titleAr string) (bool, error)
	ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error)
}

type ComplianceRepository struct {
	db *gorm.DB
}

func NewComplianceRepository(db *gorm.DB) *ComplianceRepository {
	return &ComplianceRepository{
		db: db,
	}
}

func (repo *ComplianceRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ComplianceRegister, error) {
	var compliances []*models.ComplianceRegister

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
	result := query.Find(&compliances)

	return compliances, result.Error
}

func (repo *ComplianceRepository) GetByID(id uint) (*models.ComplianceRegister, error) {
	var compliance *models.ComplianceRegister
	result := repo.db.First(&compliance, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, id)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return compliance, result.Error
}

func (repo *ComplianceRepository) Create(compliance *models.ComplianceRegister) (*models.ComplianceRegister, error) {
	result := repo.db.Create(&compliance)

	return compliance, result.Error
}

func (repo *ComplianceRepository) Update(id uint, updatedCompliance *models.ComplianceRegister) (*models.ComplianceRegister, error) {
	var compliance *models.ComplianceRegister

	if err := repo.db.First(&compliance, id).Error; err != nil {
		return nil, err
	}

	if err := repo.db.Model(&compliance).Updates(updatedCompliance).Error; err != nil {
		return nil, err
	}

	return compliance, nil
}

func (repo *ComplianceRepository) Delete(id uint) error {

	var compliance *models.ComplianceRegister
	if err := repo.db.First(&compliance, id).Error; err != nil {
		return err
	}

	result := repo.db.Delete(&models.ComplianceRegister{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *ComplianceRepository) Activate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint

	if err := repo.db.Model(&models.ComplianceRegister{}).
		Where("id IN ? AND status = ?", ids, "inactive").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err
	}

	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.ComplianceRegister{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "active").Error; err != nil {
			return nil, err
		}
	}

	return idsToUpdate, nil
}

func (repo *ComplianceRepository) Deactivate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint

	if err := repo.db.Model(&models.ComplianceRegister{}).
		Where("id IN ? AND status = ?", ids, "active").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err
	}

	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.ComplianceRegister{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "inactive").Error; err != nil {
			return nil, err
		}
	}

	return idsToUpdate, nil
}

func (repo *ComplianceRepository) ExistsByTitle(titleEn, titleAr string) (bool, error) {
	var count int64

	err := repo.db.Model(&models.ComplianceRegister{}).
		Where("title_en = ? OR title_ar = ?", titleEn, titleAr).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)
	}

	return count > 0, nil
}

func (repo *ComplianceRepository) ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error) {
	var count int64

	err := repo.db.Model(&models.ComplianceRegister{}).
		Where("id <> ? AND (title_en = ? OR title_ar = ?)", id, titleEn, titleAr).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)
	}

	return count > 0, nil
}
