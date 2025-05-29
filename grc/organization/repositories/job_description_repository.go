package repositories

import (
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interface
type JobDescriptionRepositoryInterface interface {
	GetAll(sort, sortBy, keyword, status string) ([]*models.JobDescription, error)
	GetByID(id uint) (*models.JobDescription, error)
	Create(jobDesc *models.JobDescription) (*models.JobDescription, error)
	Update(id uint, updatedJobDesc *models.JobDescription) (*models.JobDescription, error)
	Delete(id uint) error
	Activate(ids []uint) ([]uint, error)
	Deactivate(ids []uint) ([]uint, error)
	ExistsByTitle(titleEn, titleAr string) (bool, error)
	ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error)
}

// Implementation
type JobDescriptionRepository struct {
	db *gorm.DB
}

func NewJobDescriptionRepository(db *gorm.DB) *JobDescriptionRepository {
	return &JobDescriptionRepository{db: db}
}

func (repo *JobDescriptionRepository) GetAll(sort, sortBy, keyword, status string) ([]*models.JobDescription, error) {
	var jobDescs []*models.JobDescription

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
	err := query.Find(&jobDescs).Error
	return jobDescs, err
}

func (repo *JobDescriptionRepository) GetByID(id uint) (*models.JobDescription, error) {
	var jobDesc models.JobDescription
	err := repo.db.First(&jobDesc, id).Error
	if err != nil {
		return nil, err
	}
	return &jobDesc, nil
}

func (repo *JobDescriptionRepository) Create(jobDesc *models.JobDescription) (*models.JobDescription, error) {
	err := repo.db.Create(jobDesc).Error
	return jobDesc, err
}

func (repo *JobDescriptionRepository) Update(id uint, updatedJobDesc *models.JobDescription) (*models.JobDescription, error) {
	var existing models.JobDescription
	if err := repo.db.First(&existing, id).Error; err != nil {
		return nil, err
	}

	if err := repo.db.Model(&existing).Updates(updatedJobDesc).Error; err != nil {
		return nil, err
	}

	return &existing, nil
}

func (repo *JobDescriptionRepository) Delete(id uint) error {
	if err := repo.db.First(&models.JobDescription{}, id).Error; err != nil {
		return err
	}
	return repo.db.Delete(&models.JobDescription{}, id).Error
}

func (repo *JobDescriptionRepository) Activate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint
	err := repo.db.Model(&models.JobDescription{}).
		Where("id IN ? AND status = ?", ids, "inactive").
		Pluck("id", &idsToUpdate).Error
	if err != nil {
		return nil, err
	}

	if len(idsToUpdate) > 0 {
		err = repo.db.Model(&models.JobDescription{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "active").Error
		if err != nil {
			return nil, err
		}
	}

	return idsToUpdate, nil
}

func (repo *JobDescriptionRepository) Deactivate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint
	err := repo.db.Model(&models.JobDescription{}).
		Where("id IN ? AND status = ?", ids, "active").
		Pluck("id", &idsToUpdate).Error
	if err != nil {
		return nil, err
	}

	if len(idsToUpdate) > 0 {
		err = repo.db.Model(&models.JobDescription{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "inactive").Error
		if err != nil {
			return nil, err
		}
	}

	return idsToUpdate, nil
}

func (repo *JobDescriptionRepository) ExistsByTitle(titleEn, titleAr string) (bool, error) {
	var count int64
	err := repo.db.Model(&models.JobDescription{}).
		Where("title_en = ? OR title_ar = ?", titleEn, titleAr).
		Count(&count).Error
	if err != nil {
		return false, utility.NewInternalServerError(err)
	}
	return count > 0, nil
}

func (repo *JobDescriptionRepository) ExistsByTitleExcludingID(id uint, titleEn, titleAr string) (bool, error) {
	var count int64
	err := repo.db.Model(&models.JobDescription{}).
		Where("id <> ? AND (title_en = ? OR title_ar = ?)", id, titleEn, titleAr).
		Count(&count).Error
	if err != nil {
		return false, utility.NewInternalServerError(err)
	}
	return count > 0, nil
}
