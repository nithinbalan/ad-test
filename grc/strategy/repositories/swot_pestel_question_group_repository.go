package repositories

import (
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interfaces
type SwotPestelQuestionGroupRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string, isSwot string, isPestel string) ([]*models.SwotPestelQuestionGroup, error)
	GetByID(id uint) (*models.SwotPestelQuestionGroup, error)
	Create(swotPestelQuestionGroup *models.SwotPestelQuestionGroup) (*models.SwotPestelQuestionGroup, error)
	Update(id uint, updatedSwotPestelQuestionGroup *models.SwotPestelQuestionGroup) (*models.SwotPestelQuestionGroup, error)
	Delete(id uint) error
	Activate(ids []uint) ([]uint, error)
	Deactivate(ids []uint) ([]uint, error)
	ExistsByTitle(titleEn string, titleAr string) (bool, error)
	ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error)
}

// implementation
type SwotPestelQuestionGroupRepository struct {
	db *gorm.DB
}

func NewSwotPestelQuestionGroupRepository(db *gorm.DB) *SwotPestelQuestionGroupRepository {
	return &SwotPestelQuestionGroupRepository{
		db: db,
	}
}

func (repo *SwotPestelQuestionGroupRepository) GetAll(sort string, sortBy string, keyword string, status string, isSwot string, isPestel string) ([]*models.SwotPestelQuestionGroup, error) {
	var swotPestelQuestionGroups []*models.SwotPestelQuestionGroup

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

	if isSwot == "true" {
		query = query.Where("is_swot = ?", true)
	}

	if isPestel == "true" {
		query = query.Where("is_pestel = ?", true)
	}

	query = query.Order(sortBy + " " + sort)
	result := query.Find(&swotPestelQuestionGroups)

	return swotPestelQuestionGroups, result.Error
}

func (repo *SwotPestelQuestionGroupRepository) GetByID(id uint) (*models.SwotPestelQuestionGroup, error) {
	var swotPestelQuestionGroup *models.SwotPestelQuestionGroup
	result := repo.db.First(&swotPestelQuestionGroup, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return swotPestelQuestionGroup, nil
}

func (repo *SwotPestelQuestionGroupRepository) Create(swotPestelQuestionGroup *models.SwotPestelQuestionGroup) (*models.SwotPestelQuestionGroup, error) {
	result := repo.db.Create(&swotPestelQuestionGroup)

	return swotPestelQuestionGroup, result.Error
}

func (repo *SwotPestelQuestionGroupRepository) Update(id uint, updatedSwotPestelQuestionGroup *models.SwotPestelQuestionGroup) (*models.SwotPestelQuestionGroup, error) {
	var swotPestelQuestionGroup *models.SwotPestelQuestionGroup

	if err := repo.db.First(&swotPestelQuestionGroup, id).Error; err != nil {
		return nil, err
	}

	if err := repo.db.Model(&swotPestelQuestionGroup).Updates(updatedSwotPestelQuestionGroup).Error; err != nil {
		return nil, err
	}

	return swotPestelQuestionGroup, nil
}

func (repo *SwotPestelQuestionGroupRepository) Delete(id uint) error {

	var swotPestelQuestionGroup *models.SwotPestelQuestionGroup
	if err := repo.db.First(&swotPestelQuestionGroup, id).Error; err != nil {
		return err
	}

	result := repo.db.Delete(&models.SwotPestelQuestionGroup{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *SwotPestelQuestionGroupRepository) Activate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint

	if err := repo.db.Model(&models.SwotPestelQuestionGroup{}).
		Where("id IN ? AND status = ?", ids, "inactive").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err
	}

	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.SwotPestelQuestionGroup{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "active").Error; err != nil {
			return nil, err
		}
	}

	return idsToUpdate, nil
}

func (repo *SwotPestelQuestionGroupRepository) Deactivate(ids []uint) ([]uint, error) {

	var idsToUpdate []uint

	if err := repo.db.Model(&models.SwotPestelQuestionGroup{}).
		Where("id IN ? AND status = ?", ids, "active").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err
	}

	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.SwotPestelQuestionGroup{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "inactive").Error; err != nil {
			return nil, err // Return error if the update operation fails
		}
	}

	return idsToUpdate, nil

}

func (repo *SwotPestelQuestionGroupRepository) ExistsByTitle(titleEn, titleAr string) (bool, error) {
	var count int64

	err := repo.db.Model(&models.SwotPestelQuestionGroup{}).
		Where("title_en = ? OR title_ar = ?", titleEn, titleAr).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)

	}
	return count > 0, nil
}

func (repo *SwotPestelQuestionGroupRepository) ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error) {
	var count int64

	err := repo.db.Model(&models.SwotPestelQuestionGroup{}).
		Where("id <> ? AND (title_en = ? OR title_ar = ?)", id, titleEn, titleAr).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)
	}

	return count > 0, nil
}
