package repositories

import (
	"errors"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interface
type ViolationReportChannelRepositoryInterface interface {
	GetAll(sort, sortBy, keyword, status string) ([]*models.ViolationReportChannel, error)
	GetByID(id uint) (*models.ViolationReportChannel, error)
	Create(vrc *models.ViolationReportChannel) (*models.ViolationReportChannel, error)
	Update(id uint, vrc *models.ViolationReportChannel) (*models.ViolationReportChannel, error)
	Delete(id uint) error
	ExistsByTitle(titleEn, titleAr string) (bool, error)
	ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error)
}

// Implementation
type ViolationReportChannelRepository struct {
	db *gorm.DB
}

func NewViolationReportChannelRepository(db *gorm.DB) *ViolationReportChannelRepository {
	return &ViolationReportChannelRepository{db: db}
}

func (repo *ViolationReportChannelRepository) GetAll(sort, sortBy, keyword, status string) ([]*models.ViolationReportChannel, error) {
	var vrcs []*models.ViolationReportChannel
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
	result := query.Find(&vrcs)

	return vrcs, result.Error
}

func (repo *ViolationReportChannelRepository) GetByID(id uint) (*models.ViolationReportChannel, error) {
	var vrc *models.ViolationReportChannel
	result := repo.db.First(&vrc, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, id)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return vrc, result.Error
}

func (repo *ViolationReportChannelRepository) Create(vrc *models.ViolationReportChannel) (*models.ViolationReportChannel, error) {
	result := repo.db.Create(&vrc)
	return vrc, result.Error
}

func (repo *ViolationReportChannelRepository) Update(id uint, updatedVRC *models.ViolationReportChannel) (*models.ViolationReportChannel, error) {
	var vrc *models.ViolationReportChannel
	if err := repo.db.First(&vrc, id).Error; err != nil {
		return nil, err
	}

	if err := repo.db.Model(&vrc).Updates(updatedVRC).Error; err != nil {
		return nil, err
	}

	return vrc, nil
}

func (repo *ViolationReportChannelRepository) Delete(id uint) error {
	var vrc *models.ViolationReportChannel
	if err := repo.db.First(&vrc, id).Error; err != nil {
		return err
	}
	return repo.db.Delete(&models.ViolationReportChannel{}, id).Error
}

func (repo *ViolationReportChannelRepository) ExistsByTitle(titleEn, titleAr string) (bool, error) {
	var count int64
	err := repo.db.Model(&models.ViolationReportChannel{}).
		Where("title_en = ? OR title_ar = ?", titleEn, titleAr).
		Count(&count).Error
	if err != nil {
		return false, utility.NewInternalServerError(err)
	}
	return count > 0, nil
}

func (repo *ViolationReportChannelRepository) ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error) {
	var count int64
	err := repo.db.Model(&models.ViolationReportChannel{}).
		Where("id <> ? AND (title_en = ? OR title_ar = ?)", id, titleEn, titleAr).
		Count(&count).Error
	if err != nil {
		return false, utility.NewInternalServerError(err)
	}
	return count > 0, nil
}
