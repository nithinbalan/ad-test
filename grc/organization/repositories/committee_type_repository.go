package repositories

import (
	"jethur/grc/models"
	"jethur/grc/organization/dtos"

	"gorm.io/gorm"
)

// Interfaces
type CommitteeTypeRepositoryInterface interface {
	GetAll(params *dtos.CommitteeTypeListQueryParamsDTO) ([]*models.CommitteeType, error)
	GetByID(id uint) (*models.CommitteeType, error)
	GetByType(typeFilter string) (*models.CommitteeType, error)
	Create(ct *models.CommitteeType) (uint, error)
	Update(ct *models.CommitteeType) error
	Delete(id uint) error
	ExistsByTitle(titleEn, titleAr string) (bool, error)
	ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error)
}

// implementation
type CommitteeTypeRepository struct {
	db *gorm.DB
}

func NewCommitteeTypeRepository(db *gorm.DB) *CommitteeTypeRepository {
	return &CommitteeTypeRepository{
		db: db,
	}
}

func (repo *CommitteeTypeRepository) GetAll(params *dtos.CommitteeTypeListQueryParamsDTO) ([]*models.CommitteeType, error) {
	var committeeTypes []*models.CommitteeType

	query := repo.db

	if params.Keyword != "" {
		query = query.Where("(title_en LIKE ? OR title_ar LIKE ?)", "%"+params.Keyword+"%", "%"+params.Keyword+"%")
	}
	if params.SortBy == "" {
		params.SortBy = "title_en"
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "asc"
	}

	if params.Status == "active" || params.Status == "inactive" {
		query = query.Where("status = ?", params.Status)
	}

	if params.Type == "statutory_reference" || params.Type == "no_statutory_reference" || params.Type == "other" {
		query = query.Where("type = ?", params.Type)
	}

	query = query.Order(params.SortBy + " " + params.Sort)
	result := query.Find(&committeeTypes)

	return committeeTypes, result.Error
}

func (repo *CommitteeTypeRepository) GetByID(id uint) (*models.CommitteeType, error) {
	var committeeType *models.CommitteeType

	result := repo.db.First(&committeeType, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return committeeType, result.Error
}

func (repo *CommitteeTypeRepository) GetByType(typeFilter string) (*models.CommitteeType, error) {
	var committeeType *models.CommitteeType

	result := repo.db.Where("type = ?", typeFilter).First(&committeeType)

	if result.Error != nil {
		return nil, result.Error
	}

	return committeeType, nil
}

func (repo *CommitteeTypeRepository) Create(ct *models.CommitteeType) (uint, error) {
	if err := repo.db.Create(ct).Error; err != nil {
		return 0, err
	}

	return ct.ID, nil
}

func (repo *CommitteeTypeRepository) Update(ct *models.CommitteeType) error {
	return repo.db.Save(ct).Error
}

func (repo *CommitteeTypeRepository) Delete(id uint) error {
	return repo.db.Delete(&models.CommitteeType{}, id).Error
}

func (repo *CommitteeTypeRepository) ExistsByTitle(titleEn, titleAr string) (bool, error) {
	var count int64
	err := repo.db.Model(&models.CommitteeType{}).
		Where("title_en = ? OR title_ar = ?", titleEn, titleAr).
		Count(&count).Error

	return count > 0, err
}

func (repo *CommitteeTypeRepository) ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error) {
	var count int64
	err := repo.db.Model(&models.CommitteeType{}).
		Where("id != ? AND (title_en = ? OR title_ar = ?)", id, titleEn, titleAr).
		Count(&count).Error

	return count > 0, err
}
