package repositories

import (
	"errors"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interface
type ViolationReportReasonRepositoryInterface interface {
	GetAll(sort, sortBy, keyword, status string) ([]*models.ViolationReportReason, error)
	GetByID(id uint) (*models.ViolationReportReason, error)
	Create(vrr *models.ViolationReportReason) (*models.ViolationReportReason, error)
	Update(id uint, vrr *models.ViolationReportReason) (*models.ViolationReportReason, error)
	Delete(id uint) error
	ExistsByTitle(titleEn, titleAr string) (bool, error)
	ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error)
}

// Implementation
type ViolationReportReasonRepository struct {
	db *gorm.DB
}

func NewViolationReportReasonRepository(db *gorm.DB) *ViolationReportReasonRepository {
	return &ViolationReportReasonRepository{db: db}
}

func (repo *ViolationReportReasonRepository) GetAll(sort, sortBy, keyword, status string) ([]*models.ViolationReportReason, error) {
	var vrrs []*models.ViolationReportReason
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
	result := query.Find(&vrrs)

	return vrrs, result.Error
}

func (repo *ViolationReportReasonRepository) GetByID(id uint) (*models.ViolationReportReason, error) {
	var vrr *models.ViolationReportReason
	result := repo.db.First(&vrr, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, id)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return vrr, result.Error
}

func (repo *ViolationReportReasonRepository) Create(vrr *models.ViolationReportReason) (*models.ViolationReportReason, error) {
	result := repo.db.Create(&vrr)
	return vrr, result.Error
}

func (repo *ViolationReportReasonRepository) Update(id uint, updatedVRR *models.ViolationReportReason) (*models.ViolationReportReason, error) {
	var vrr *models.ViolationReportReason
	if err := repo.db.First(&vrr, id).Error; err != nil {
		return nil, err
	}

	if err := repo.db.Model(&vrr).Updates(updatedVRR).Error; err != nil {
		return nil, err
	}

	return vrr, nil
}

func (repo *ViolationReportReasonRepository) Delete(id uint) error {
	var vrr *models.ViolationReportReason
	if err := repo.db.First(&vrr, id).Error; err != nil {
		return err
	}
	return repo.db.Delete(&models.ViolationReportReason{}, id).Error
}

func (repo *ViolationReportReasonRepository) ExistsByTitle(titleEn, titleAr string) (bool, error) {
	var count int64
	err := repo.db.Model(&models.ViolationReportReason{}).
		Where("title_en = ? OR title_ar = ?", titleEn, titleAr).
		Count(&count).Error
	if err != nil {
		return false, utility.NewInternalServerError(err)
	}
	return count > 0, nil
}

func (repo *ViolationReportReasonRepository) ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error) {
	var count int64
	err := repo.db.Model(&models.ViolationReportReason{}).
		Where("id <> ? AND (title_en = ? OR title_ar = ?)", id, titleEn, titleAr).
		Count(&count).Error
	if err != nil {
		return false, utility.NewInternalServerError(err)
	}
	return count > 0, nil
}
