package repositories

import (
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type AwarenessChannelRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.AwarenessChannel, error)
	GetByID(id uint) (*models.AwarenessChannel, error)
	Create(awarenessChannel *models.AwarenessChannel) (*models.AwarenessChannel, error)
	Update(id uint, awarenessChannel *models.AwarenessChannel) (*models.AwarenessChannel, error)
	Delete(id uint) error
	Activate(ids []uint) ([]uint, error)
	Deactivate(ids []uint) ([]uint, error)
	ExistsByTitle(titleEn string, titleAr string) (bool, error)
	ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error)
}

type AwarenessChannelRepository struct {
	db *gorm.DB
}

func NewAwarenessChannelRepository(db *gorm.DB) *AwarenessChannelRepository {
	return &AwarenessChannelRepository{
		db: db,
	}
}

func (repo *AwarenessChannelRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.AwarenessChannel, error) {
	var awarenessChannels []*models.AwarenessChannel
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

	result := query.Order(sortBy + " " + sort).Find(&awarenessChannels)

	return awarenessChannels, result.Error
}

func (repo *AwarenessChannelRepository) GetByID(id uint) (*models.AwarenessChannel, error) {
	var awarenessChannel models.AwarenessChannel

	if err := repo.db.First(&awarenessChannel, id).Error; err != nil {
		return nil, err
	}

	return &awarenessChannel, nil
}

func (repo *AwarenessChannelRepository) Create(awarenessChannel *models.AwarenessChannel) (*models.AwarenessChannel, error) {
	if err := repo.db.Create(awarenessChannel).Error; err != nil {
		return nil, err
	}

	return awarenessChannel, nil
}

func (repo *AwarenessChannelRepository) Update(id uint, updatedChannel *models.AwarenessChannel) (*models.AwarenessChannel, error) {
	var awarenessChannel models.AwarenessChannel

	if err := repo.db.First(&awarenessChannel, id).Error; err != nil {
		return nil, err
	}

	if err := repo.db.Model(&awarenessChannel).Updates(updatedChannel).Error; err != nil {
		return nil, err
	}

	return &awarenessChannel, nil
}

func (repo *AwarenessChannelRepository) Delete(id uint) error {
	var awarenessChannel models.AwarenessChannel

	if err := repo.db.First(&awarenessChannel, id).Error; err != nil {
		return err
	}

	return repo.db.Delete(&awarenessChannel).Error
}

func (repo *AwarenessChannelRepository) Activate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint

	if err := repo.db.Model(&models.AwarenessChannel{}).
		Where("id IN ? AND status = ?", ids, "inactive").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err
	}

	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.AwarenessChannel{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "active").Error; err != nil {
			return nil, err
		}
	}

	return idsToUpdate, nil
}

func (repo *AwarenessChannelRepository) Deactivate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint

	if err := repo.db.Model(&models.AwarenessChannel{}).
		Where("id IN ? AND status = ?", ids, "active").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err
	}

	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.AwarenessChannel{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "inactive").Error; err != nil {
			return nil, err
		}
	}

	return idsToUpdate, nil
}

func (repo *AwarenessChannelRepository) ExistsByTitle(titleEn, titleAr string) (bool, error) {
	var count int64

	// Check for the existence of titleEn or titleAr in a single query
	err := repo.db.Model(&models.AwarenessChannel{}).
		Where("title_en = ? OR title_ar = ?", titleEn, titleAr).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)

	}
	return count > 0, nil
}

func (repo *AwarenessChannelRepository) ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error) {
	var count int64
	err := repo.db.Model(&models.AwarenessChannel{}).
		Where("id <> ? AND (title_en = ? OR title_ar = ?)", id, titleEn, titleAr).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)
	}

	return count > 0, nil
}
