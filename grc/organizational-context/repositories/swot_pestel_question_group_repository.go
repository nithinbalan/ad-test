package repositories

import (
	"jethur/grc/models"
	"jethur/utility"

	organizationalContextDTOs "jethur/grc/organizational-context/dtos"

	"gorm.io/gorm"
)

// Interfaces
type SwotPestelQuestionGroupRepositoryInterface interface {
	GetAll(params *organizationalContextDTOs.SwotPestelQuestionGroupListQueryParamsDTO) ([]*organizationalContextDTOs.SwotPestelQuestionGroupListDTO, error)
	GetByID(id uint) (*models.SwotPestelQuestionGroup, error)
	Create(swotPestelQuestionGroup *models.SwotPestelQuestionGroup) (*models.SwotPestelQuestionGroup, error)
	Update(id uint, updatedSwotPestelQuestionGroup *models.SwotPestelQuestionGroup) (*models.SwotPestelQuestionGroup, error)
	Delete(id uint) error
	GetAllSwotPestelQuestionGroups(isSwot string, isPestel string) ([]*models.SwotPestelQuestionGroup, error)
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

func (repo *SwotPestelQuestionGroupRepository) GetAll(params *organizationalContextDTOs.SwotPestelQuestionGroupListQueryParamsDTO) ([]*organizationalContextDTOs.SwotPestelQuestionGroupListDTO, error) {
	var results []*organizationalContextDTOs.SwotPestelQuestionGroupListDTO

	query := repo.db.Table("swot_pestel_question_groups").Select("swot_pestel_question_groups.*")

	if keyword := params.Keyword; keyword != "" {
		like := "%" + keyword + "%"
		query = query.Where("title_en LIKE ? OR title_ar LIKE ?", like, like)
	}

	if params.Status == "active" || params.Status == "inactive" {
		query = query.Where("status = ?", params.Status)
	}

	if params.IsSwot == "true" {
		query = query.Where("is_swot = ?", true)
	}

	if params.IsPestel == "true" {
		query = query.Where("is_pestel = ?", true)
	}

	sortBy := params.SortBy
	if sortBy == "" {
		sortBy = "title_en"
	}
	sortOrder := "asc"
	if params.Sort == "desc" {
		sortOrder = "desc"
	}

	query = query.Order(sortBy + " " + sortOrder)

	err := query.Find(&results).Error

	for _, result := range results {
		var questions []*organizationalContextDTOs.QuestionDTO

		repo.db.Model(&models.SwotPestelQuestion{}).Where("swot_pestel_question_group_id = ?", result.ID).Where("status = ?", "active").Order("[order] ASC").Find(&questions)
		result.SwotPestelQuestions = questions
	}
	return results, err
}

func (repo *SwotPestelQuestionGroupRepository) GetAllSwotPestelQuestionGroups(isSwot string, isPestel string) ([]*models.SwotPestelQuestionGroup, error) {
	var swotPestelQuestionGroups []*models.SwotPestelQuestionGroup
	query := repo.db
	if isSwot == "true" {
		query = query.Where("is_swot = ?", true)
	}

	if isPestel == "true" {
		query = query.Where("is_pestel = ?", true)
	}

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
