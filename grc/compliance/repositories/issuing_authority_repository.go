package repositories

import (
	"errors"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interfaces
type IssuingAuthorityRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.IssuingAuthority, error)
	GetByID(id uint) (*models.IssuingAuthority, error)
	Create(authority *models.IssuingAuthority) (*models.IssuingAuthority, error)
	Update(id uint, authority *models.IssuingAuthority) (*models.IssuingAuthority, error)
	Delete(id uint) error
	Activate(ids []uint) ([]uint, error)
	Deactivate(ids []uint) ([]uint, error)
	ExistsByTitle(titleEn string, titleAr string) (bool, error)
	ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error)
}

// implementation
type IssuingAuthorityRepository struct {
	db *gorm.DB
}

func NewIssuingAuthorityRepository(db *gorm.DB) *IssuingAuthorityRepository {
	return &IssuingAuthorityRepository{
		db: db,
	}
}

func (repo *IssuingAuthorityRepository) GetAll(
	sort string,
	sortBy string,
	keyword string,
	status string,
) ([]*models.IssuingAuthority, error) {

	var authorities []*models.IssuingAuthority

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

	query = query.Order(sortBy + " " + sort)
	result := query.Find(&authorities)

	return authorities, result.Error
}

func (repo *IssuingAuthorityRepository) GetByID(id uint) (*models.IssuingAuthority, error) {
	var authority *models.IssuingAuthority
	result := repo.db.First(&authority, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, id)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return authority, result.Error
}

func (repo *IssuingAuthorityRepository) Create(authority *models.IssuingAuthority) (*models.IssuingAuthority, error) {
	result := repo.db.Create(&authority)
	return authority, result.Error
}

func (repo *IssuingAuthorityRepository) Update(
	id uint,
	updatedAuthority *models.IssuingAuthority,
) (*models.IssuingAuthority, error) {

	var authority *models.IssuingAuthority

	if err := repo.db.First(&authority, id).Error; err != nil {
		return nil, err
	}

	if err := repo.db.Model(&authority).Updates(updatedAuthority).Error; err != nil {
		return nil, err
	}

	return authority, nil
}

func (repo *IssuingAuthorityRepository) Delete(id uint) error {
	var authority *models.IssuingAuthority

	if err := repo.db.First(&authority, id).Error; err != nil {
		return err
	}

	result := repo.db.Delete(&models.IssuingAuthority{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *IssuingAuthorityRepository) Activate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint

	// Fetch inactive records
	if err := repo.db.Model(&models.IssuingAuthority{}).
		Where("id IN ? AND status = ?", ids, "inactive").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err
	}

	// Activate
	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.IssuingAuthority{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "active").Error; err != nil {
			return nil, err
		}
	}

	return idsToUpdate, nil
}

func (repo *IssuingAuthorityRepository) Deactivate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint

	// Fetch active records
	if err := repo.db.Model(&models.IssuingAuthority{}).
		Where("id IN ? AND status = ?", ids, "active").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err
	}

	// Deactivate
	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.IssuingAuthority{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "inactive").Error; err != nil {
			return nil, err
		}
	}

	return idsToUpdate, nil
}

func (repo *IssuingAuthorityRepository) ExistsByTitle(titleEn, titleAr string) (bool, error) {
	var count int64

	err := repo.db.Model(&models.IssuingAuthority{}).
		Where("title_en = ? OR title_ar = ?", titleEn, titleAr).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)
	}

	return count > 0, nil
}

func (repo *IssuingAuthorityRepository) ExistsByTitleExcludingID(
	id uint,
	titleEn,
	titleAr string,
) (bool, error) {

	var count int64

	err := repo.db.Model(&models.IssuingAuthority{}).
		Where("id <> ? AND (title_en = ? OR title_ar = ?)", id, titleEn, titleAr).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)
	}

	return count > 0, nil
}
