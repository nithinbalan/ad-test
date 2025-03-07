package repositories

import (
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interfaces
type SwotPestelQuestionRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string, isSwot string) ([]*models.SwotPestelQuestion, error)
	GetAllOrdered(swotPestelQuestionGroupId uint) ([]*models.SwotPestelQuestion, error)
	GetByID(id uint) (*models.SwotPestelQuestion, error)
	Create(swotPestelQuestion *models.SwotPestelQuestion) (*models.SwotPestelQuestion, error)
	Update(id uint, updatedSwotPestelQuestion *models.SwotPestelQuestion) (*models.SwotPestelQuestion, error)
	Delete(id uint) error
	Activate(ids []uint) ([]uint, error)
	Deactivate(ids []uint) ([]uint, error)
	ExistsByTitle(swotPestelQuestionGroupId uint, titleEn string, titleAr string) (bool, error)
	ExistsByTitleExcludingID(swotPestelQuestionGroupId uint, id uint, titleEn, titleAr string) (bool, error)
}

// implementation
type SwotPestelQuestionRepository struct {
	db *gorm.DB
}

func NewSwotPestelQuestionRepository(db *gorm.DB) *SwotPestelQuestionRepository {
	return &SwotPestelQuestionRepository{
		db: db,
	}
}

func (repo *SwotPestelQuestionRepository) GetAll(sort string, sortBy string, keyword string, status string, isSwot string) ([]*models.SwotPestelQuestion, error) {
	var swotPestelQuestions []*models.SwotPestelQuestion

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
	} else if isSwot == "false" {
		query = query.Where("is_swot = ?", false)
	}

	query = query.Order(sortBy + " " + sort)
	result := query.Find(&swotPestelQuestions)

	return swotPestelQuestions, result.Error
}

func (repo *SwotPestelQuestionRepository) GetAllOrdered(swotPestelQuestionGroupID uint) ([]*models.SwotPestelQuestion, error) {
	var questions []*models.SwotPestelQuestion

	err := repo.db.
		Where("swot_question_group_id = ?", swotPestelQuestionGroupID).
		Order("[order] ASC").
		Find(&questions).Error

	if err != nil {
		return nil, err
	}

	return questions, nil
}

func (repo *SwotPestelQuestionRepository) GetByID(id uint) (*models.SwotPestelQuestion, error) {
	var swotPestelQuestion *models.SwotPestelQuestion
	result := repo.db.First(&swotPestelQuestion, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return swotPestelQuestion, nil
}

func (repo *SwotPestelQuestionRepository) Create(swotPestelQuestion *models.SwotPestelQuestion) (*models.SwotPestelQuestion, error) {
	result := repo.db.Create(&swotPestelQuestion)

	return swotPestelQuestion, result.Error
}

func (repo *SwotPestelQuestionRepository) Update(id uint, updatedSwotPestelQuestion *models.SwotPestelQuestion) (*models.SwotPestelQuestion, error) {
	var swotPestelQuestion *models.SwotPestelQuestion

	if err := repo.db.First(&swotPestelQuestion, id).Error; err != nil {
		return nil, err
	}

	if err := repo.db.Model(&swotPestelQuestion).Updates(updatedSwotPestelQuestion).Error; err != nil {
		return nil, err
	}

	return swotPestelQuestion, nil
}

func (repo *SwotPestelQuestionRepository) Delete(id uint) error {

	var swotPestelQuestion *models.SwotPestelQuestion
	if err := repo.db.First(&swotPestelQuestion, id).Error; err != nil {
		return err
	}

	result := repo.db.Delete(&models.SwotPestelQuestion{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *SwotPestelQuestionRepository) Activate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint

	if err := repo.db.Model(&models.SwotPestelQuestion{}).
		Where("id IN ? AND status = ?", ids, "inactive").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err
	}

	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.SwotPestelQuestion{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "active").Error; err != nil {
			return nil, err
		}
	}

	return idsToUpdate, nil
}

func (repo *SwotPestelQuestionRepository) Deactivate(ids []uint) ([]uint, error) {

	var idsToUpdate []uint

	if err := repo.db.Model(&models.SwotPestelQuestion{}).
		Where("id IN ? AND status = ?", ids, "active").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err
	}

	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.SwotPestelQuestion{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "inactive").Error; err != nil {
			return nil, err
		}
	}

	return idsToUpdate, nil

}

func (repo *SwotPestelQuestionRepository) ExistsByTitle(swotPestelQuestionGroupId uint, titleEn, titleAr string) (bool, error) {
	var count int64

	err := repo.db.Model(&models.SwotPestelQuestion{}).
		Where("(title_en = ? OR title_ar = ?) AND swot_question_group_id = ?", titleEn, titleAr, swotPestelQuestionGroupId).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)

	}
	return count > 0, nil
}

func (repo *SwotPestelQuestionRepository) ExistsByTitleExcludingID(swotPestelQuestionGroupId uint, id uint, titleEn, titleAr string) (bool, error) {
	var count int64

	err := repo.db.Model(&models.SwotPestelQuestion{}).
		Where("id <> ? AND swot_question_group_id = ? AND (title_en = ? OR title_ar = ?)", id, swotPestelQuestionGroupId, titleEn, titleAr).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)
	}

	return count > 0, nil
}
