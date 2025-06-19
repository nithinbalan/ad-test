package repositories

import (
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interfaces
type CompetencyGroupRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.CompetencyGroup, error)
	GetByID(id uint) (*models.CompetencyGroup, error)
	Create(competencyGroup *models.CompetencyGroup) (*models.CompetencyGroup, error)
	Update(id uint, updatedCompetencyGroup *models.CompetencyGroup) (*models.CompetencyGroup, error)
	Delete(id uint) error
	Activate(ids []uint) ([]uint, error)
	Deactivate(ids []uint) ([]uint, error)
	ExistsByTitle(titleEn string, titleAr string) (bool, error)
	ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error)
}

// implementation
type CompetencyGroupRepository struct {
	db *gorm.DB
}

func NewCompetencyGroupRepository(db *gorm.DB) *CompetencyGroupRepository {
	return &CompetencyGroupRepository{
		db: db,
	}
}

func (repo *CompetencyGroupRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.CompetencyGroup, error) {
	var competencyGroups []*models.CompetencyGroup

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
	result := query.Find(&competencyGroups)

	return competencyGroups, result.Error
}

func (repo *CompetencyGroupRepository) GetByID(id uint) (*models.CompetencyGroup, error) {
	var competencyGroup *models.CompetencyGroup
	result := repo.db.First(&competencyGroup, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return competencyGroup, nil
}

func (repo *CompetencyGroupRepository) Create(competencyGroup *models.CompetencyGroup) (*models.CompetencyGroup, error) {
	result := repo.db.Create(&competencyGroup)

	return competencyGroup, result.Error
}

func (repo *CompetencyGroupRepository) Update(id uint, updatedCompetencyGroup *models.CompetencyGroup) (*models.CompetencyGroup, error) {
	var competencyGroup *models.CompetencyGroup

	if err := repo.db.First(&competencyGroup, id).Error; err != nil {
		return nil, err
	}

	if err := repo.db.Model(&competencyGroup).Updates(updatedCompetencyGroup).Error; err != nil {
		return nil, err
	}

	return competencyGroup, nil
}

func (repo *CompetencyGroupRepository) Delete(id uint) error {
	var competencyGroup *models.CompetencyGroup
	if err := repo.db.First(&competencyGroup, id).Error; err != nil {
		return err
	}

	result := repo.db.Delete(&models.CompetencyGroup{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *CompetencyGroupRepository) Activate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint

	if err := repo.db.Model(&models.CompetencyGroup{}).
		Where("id IN ? AND status = ?", ids, "inactive").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err
	}

	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.CompetencyGroup{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "active").Error; err != nil {
			return nil, err
		}
	}

	return idsToUpdate, nil
}

func (repo *CompetencyGroupRepository) Deactivate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint

	if err := repo.db.Model(&models.CompetencyGroup{}).
		Where("id IN ? AND status = ?", ids, "active").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err
	}

	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.CompetencyGroup{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "inactive").Error; err != nil {
			return nil, err
		}
	}

	return idsToUpdate, nil

}

func (repo *CompetencyGroupRepository) ExistsByTitle(titleEn, titleAr string) (bool, error) {
	var count int64

	err := repo.db.Model(&models.CompetencyGroup{}).
		Where("title_en = ? OR title_ar = ?", titleEn, titleAr).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)

	}

	return count > 0, nil
}

func (repo *CompetencyGroupRepository) ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error) {
	var count int64

	err := repo.db.Model(&models.CompetencyGroup{}).
		Where("id <> ? AND (title_en = ? OR title_ar = ?)", id, titleEn, titleAr).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)
	}

	return count > 0, nil
}
