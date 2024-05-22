package repositories

import (
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interfaces
type MeetingTypeRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.MeetingType, error)
	GetByID(id uint) (*models.MeetingType, error)
	Create(article *models.MeetingType) (*models.MeetingType, error)
	Update(id uint, meetingType *models.MeetingType) (*models.MeetingType, error)
	Delete(id uint) error
	Activate(ids []uint) ([]uint, error)
	Deactivate(ids []uint) ([]uint, error)
	ExistsByTitle(titleEn string, titleAr string) (bool, error)
	ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error)
}

// implementation
type MeetingTypeRepository struct {
	db *gorm.DB
}

func NewMeetingTypeRepository(db *gorm.DB) *MeetingTypeRepository {
	return &MeetingTypeRepository{
		db: db,
	}
}

func (repo *MeetingTypeRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.MeetingType, error) {
	var meetingTypes []*models.MeetingType
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
	result := query.Find(&meetingTypes)

	return meetingTypes, result.Error
}

func (repo *MeetingTypeRepository) GetByID(id uint) (*models.MeetingType, error) {
	var meetingType *models.MeetingType
	result := repo.db.First(&meetingType, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return meetingType, nil
}

func (repo *MeetingTypeRepository) Create(meetingType *models.MeetingType) (*models.MeetingType, error) {
	result := repo.db.Create(&meetingType)

	return meetingType, result.Error
}

func (repo *MeetingTypeRepository) Update(id uint, updatedType *models.MeetingType) (*models.MeetingType, error) {
	var meetingType *models.MeetingType
	if err := repo.db.First(&meetingType, id).Error; err != nil {
		return nil, err
	}

	if err := repo.db.Model(&meetingType).Updates(updatedType).Error; err != nil {
		return nil, err
	}

	return meetingType, nil
}

func (repo *MeetingTypeRepository) Delete(id uint) error {
	var meetingType *models.MeetingType
	if err := repo.db.First(&meetingType, id).Error; err != nil {
		return err
	}

	result := repo.db.Delete(&models.MeetingType{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *MeetingTypeRepository) Activate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint
	// Step 1: Fetch the IDs of records that are currently 'inactive' and within the provided 'ids' slice
	if err := repo.db.Model(&models.MeetingType{}).
		Where("id IN ? AND status = ?", ids, "inactive").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err // Return error if the query fails
	}

	// Step 2: If there are any IDs to update, perform the update operation
	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.MeetingType{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "active").Error; err != nil {
			return nil, err // Return error if the update operation fails
		}
	}
	// Return the IDs of the records that were activated
	return idsToUpdate, nil
}

func (repo *MeetingTypeRepository) Deactivate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint
	// Step 1: Fetch the IDs of records that are currently 'inactive' and within the provided 'ids' slice
	if err := repo.db.Model(&models.MeetingType{}).
		Where("id IN ? AND status = ?", ids, "active").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err
	}

	// Step 2: If there are any IDs to update, perform the update operation
	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.MeetingType{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "inactive").Error; err != nil {
			return nil, err
		}
	}
	return idsToUpdate, nil
}

func (repo *MeetingTypeRepository) ExistsByTitle(titleEn, titleAr string) (bool, error) {
	var count int64

	// Check for the existence of titleEn or titleAr in a single query
	err := repo.db.Model(&models.MeetingType{}).
		Where("title_en = ? OR title_ar = ?", titleEn, titleAr).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)

	}
	return count > 0, nil
}

func (repo *MeetingTypeRepository) ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error) {
	var count int64
	err := repo.db.Model(&models.MeetingType{}).
		Where("id <> ? AND (title_en = ? OR title_ar = ?)", id, titleEn, titleAr).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)
	}

	return count > 0, nil
}
