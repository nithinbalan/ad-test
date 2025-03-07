package repositories

import (
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interfaces
type SwotQuestionRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.SwotQuestion, error)
	GetAllOrdered(swotQuestionGroupId uint) ([]*models.SwotQuestion, error)
	GetByID(id uint) (*models.SwotQuestion, error)
	Create(swotQuestion *models.SwotQuestion) (*models.SwotQuestion, error)
	Update(id uint, updatedSwotQuestion *models.SwotQuestion) (*models.SwotQuestion, error)
	Delete(id uint) error
	Activate(ids []uint) ([]uint, error)
	Deactivate(ids []uint) ([]uint, error)
	ExistsByTitle(swotQuestionGroupId uint, titleEn string, titleAr string) (bool, error)
	ExistsByTitleExcludingID(swotQuestionGroupId uint, id uint, titleEn, titleAr string) (bool, error)
}

// implementation
type SwotQuestionRepository struct {
	db *gorm.DB
}

func NewSwotQuestionRepository(db *gorm.DB) *SwotQuestionRepository {
	return &SwotQuestionRepository{
		db: db,
	}
}

func (repo *SwotQuestionRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.SwotQuestion, error) {
	var swotQuestions []*models.SwotQuestion

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
	result := query.Find(&swotQuestions)

	return swotQuestions, result.Error
}

func (repo *SwotQuestionRepository) GetAllOrdered(swotQuestionGroupID uint) ([]*models.SwotQuestion, error) {
	var questions []*models.SwotQuestion

	err := repo.db.
		Where("swot_question_group_id = ?", swotQuestionGroupID).
		Order("[order] ASC").
		Find(&questions).Error

	if err != nil {
		return nil, err
	}

	return questions, nil
}

func (repo *SwotQuestionRepository) GetByID(id uint) (*models.SwotQuestion, error) {
	var swotQuestion *models.SwotQuestion
	result := repo.db.First(&swotQuestion, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return swotQuestion, nil
}

func (repo *SwotQuestionRepository) Create(swotQuestion *models.SwotQuestion) (*models.SwotQuestion, error) {
	result := repo.db.Create(&swotQuestion)

	return swotQuestion, result.Error
}

func (repo *SwotQuestionRepository) Update(id uint, updatedSwotQuestion *models.SwotQuestion) (*models.SwotQuestion, error) {
	var swotQuestion *models.SwotQuestion

	if err := repo.db.First(&swotQuestion, id).Error; err != nil {
		return nil, err
	}

	if err := repo.db.Model(&swotQuestion).Updates(updatedSwotQuestion).Error; err != nil {
		return nil, err
	}

	return swotQuestion, nil
}

func (repo *SwotQuestionRepository) Delete(id uint) error {

	var swotQuestion *models.SwotQuestion
	if err := repo.db.First(&swotQuestion, id).Error; err != nil {
		return err
	}

	result := repo.db.Delete(&models.SwotQuestion{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *SwotQuestionRepository) Activate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint

	if err := repo.db.Model(&models.SwotQuestion{}).
		Where("id IN ? AND status = ?", ids, "inactive").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err
	}

	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.SwotQuestion{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "active").Error; err != nil {
			return nil, err
		}
	}

	return idsToUpdate, nil
}

func (repo *SwotQuestionRepository) Deactivate(ids []uint) ([]uint, error) {

	var idsToUpdate []uint

	if err := repo.db.Model(&models.SwotQuestion{}).
		Where("id IN ? AND status = ?", ids, "active").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err
	}

	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.SwotQuestion{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "inactive").Error; err != nil {
			return nil, err
		}
	}

	return idsToUpdate, nil

}

func (repo *SwotQuestionRepository) ExistsByTitle(swotQuestionGroupId uint, titleEn, titleAr string) (bool, error) {
	var count int64

	err := repo.db.Model(&models.SwotQuestion{}).
		Where("(title_en = ? OR title_ar = ?) AND swot_question_group_id = ?", titleEn, titleAr, swotQuestionGroupId).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)

	}
	return count > 0, nil
}

func (repo *SwotQuestionRepository) ExistsByTitleExcludingID(swotQuestionGroupId uint, id uint, titleEn, titleAr string) (bool, error) {
	var count int64

	err := repo.db.Model(&models.SwotQuestion{}).
		Where("id <> ? AND swot_question_group_id = ? AND (title_en = ? OR title_ar = ?)", id, swotQuestionGroupId, titleEn, titleAr).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)
	}

	return count > 0, nil
}
