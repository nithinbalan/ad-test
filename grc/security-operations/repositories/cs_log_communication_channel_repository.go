package repositories

import (
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type CSLogCommunicationChannelRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.CSLogCommunicationChannel, error)
	GetByID(id uint) (*models.CSLogCommunicationChannel, error)
	Create(channel *models.CSLogCommunicationChannel) (*models.CSLogCommunicationChannel, error)
	Update(id uint, channel *models.CSLogCommunicationChannel) (*models.CSLogCommunicationChannel, error)
	Delete(id uint) error
	Activate(ids []uint) ([]uint, error)
	Deactivate(ids []uint) ([]uint, error)
	ExistsByTitle(titleEn string, titleAr string) (bool, error)
	ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error)
}

type CSLogCommunicationChannelRepository struct {
	db *gorm.DB
}

func NewCSLogCommunicationChannelRepository(db *gorm.DB) *CSLogCommunicationChannelRepository {
	return &CSLogCommunicationChannelRepository{
		db: db,
	}
}

func (repo *CSLogCommunicationChannelRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.CSLogCommunicationChannel, error) {
	var channels []*models.CSLogCommunicationChannel
	query := repo.db

	if keyword != "" {
		query = query.Where("(title_en LIKE ? OR title_ar LIKE ?)", "%"+keyword+"%", "%"+keyword+"%")
	}
	if sortBy == "" {
		sortBy = "id"
	}
	if sort != "asc" && sort != "desc" {
		sort = "desc"
	}

	if status == "active" || status == "inactive" {
		query = query.Where("status = ?", status)
	}

	query = query.Order(sortBy + " " + sort)
	result := query.Find(&channels)

	return channels, result.Error
}

func (repo *CSLogCommunicationChannelRepository) GetByID(id uint) (*models.CSLogCommunicationChannel, error) {
	var channel *models.CSLogCommunicationChannel
	result := repo.db.First(&channel, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return channel, nil
}

func (repo *CSLogCommunicationChannelRepository) Create(channel *models.CSLogCommunicationChannel) (*models.CSLogCommunicationChannel, error) {
	result := repo.db.Create(&channel)

	return channel, result.Error
}

func (repo *CSLogCommunicationChannelRepository) Update(id uint, updated *models.CSLogCommunicationChannel) (*models.CSLogCommunicationChannel, error) {
	var channel *models.CSLogCommunicationChannel
	if err := repo.db.First(&channel, id).Error; err != nil {
		return nil, err
	}

	if err := repo.db.Model(&channel).Updates(updated).Error; err != nil {
		return nil, err
	}

	return channel, nil
}

func (repo *CSLogCommunicationChannelRepository) Delete(id uint) error {
	var channel *models.CSLogCommunicationChannel
	if err := repo.db.First(&channel, id).Error; err != nil {
		return err
	}

	result := repo.db.Delete(&models.CSLogCommunicationChannel{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *CSLogCommunicationChannelRepository) Activate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint
	if err := repo.db.Model(&models.CSLogCommunicationChannel{}).
		Where("id IN ? AND status = ?", ids, "inactive").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err
	}

	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.CSLogCommunicationChannel{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "active").Error; err != nil {
			return nil, err
		}
	}
	return idsToUpdate, nil
}

func (repo *CSLogCommunicationChannelRepository) Deactivate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint
	if err := repo.db.Model(&models.CSLogCommunicationChannel{}).
		Where("id IN ? AND status = ?", ids, "active").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err
	}

	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.CSLogCommunicationChannel{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "inactive").Error; err != nil {
			return nil, err
		}
	}
	return idsToUpdate, nil
}

func (repo *CSLogCommunicationChannelRepository) ExistsByTitle(titleEn, titleAr string) (bool, error) {
	var count int64
	err := repo.db.Model(&models.CSLogCommunicationChannel{}).
		Where("title_en = ? OR title_ar = ?", titleEn, titleAr).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)
	}
	return count > 0, nil
}

func (repo *CSLogCommunicationChannelRepository) ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error) {
	var count int64
	err := repo.db.Model(&models.CSLogCommunicationChannel{}).
		Where("id <> ? AND (title_en = ? OR title_ar = ?)", id, titleEn, titleAr).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)
	}
	return count > 0, nil
}
