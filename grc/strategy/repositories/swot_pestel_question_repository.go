package repositories

import (
	"jethur/grc/models"
	strategyDTOs "jethur/grc/strategy/dtos"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interfaces
type SwotPestelQuestionRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string, isSwot string, isPestel string) ([]*strategyDTOs.SwotPestelQuestionListDTO, error)
	GetAllOrdered(swotPestelQuestionGroupId uint) ([]*models.SwotPestelQuestion, error)
	GetByID(id uint) (*models.SwotPestelQuestion, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.SwotPestelQuestion, error)
	Create(swotPestelQuestion *models.SwotPestelQuestion) (*models.SwotPestelQuestion, error)
	Update(id uint, updatedSwotPestelQuestion *models.SwotPestelQuestion) (*models.SwotPestelQuestion, error)
	Delete(id uint) error
	Activate(ids []uint) ([]uint, error)
	Deactivate(ids []uint) ([]uint, error)
	ExistsByTitle(swotPestelQuestionGroupId uint, titleEn string, titleAr string) (bool, error)
	ExistsByTitleExcludingID(swotPestelQuestionGroupId uint, id uint, titleEn, titleAr string) (bool, error)

	CreateSwotPestelQuestionInvitees(invitees []*models.SwotPestelQuestionInvitee, tx ...utility.DatabaseTransaction) error
	DeleteSwotPestelQuestionInvitees(swotPestelQuestionID uint, inviteeIDs []*uint, tx ...utility.DatabaseTransaction) error
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

func (repo *SwotPestelQuestionRepository) GetAll(sort string, sortBy string, keyword string, status string, isSwot string, isPestel string) ([]*strategyDTOs.SwotPestelQuestionListDTO, error) {
	var swotPestelQuestions []*strategyDTOs.SwotPestelQuestionListDTO

	query := repo.db.Table("swot_pestel_questions").
		Select(`
			swot_pestel_questions.*,
			swot_pestel_question_groups.title_en AS swot_pestel_question_group_title_en,
			swot_pestel_question_groups.title_ar AS swot_pestel_question_group_title_ar,
			answered_user.id AS answered_by_id,
			answered_user.first_name_en AS answered_user_first_name_en,
			answered_user.first_name_ar AS answered_user_first_name_ar,
			answered_user.last_name_en AS answered_user_last_name_en,
			answered_user.last_name_ar AS answered_user_last_name_ar,
			answered_user.initial AS answered_user_initial,
			answered_user.color AS answered_user_color,
			answered_user.email AS answered_user_email,
			answered_user.profile_image_path AS answered_user_profile_image_path
		`).
		Joins("LEFT JOIN swot_pestel_question_groups ON swot_pestel_question_groups.id = swot_pestel_questions.swot_pestel_question_group_id").
		Joins("LEFT JOIN users AS answered_user ON answered_user.id = swot_pestel_questions.answered_by_id")

	// Apply search keyword if exists
	if keyword != "" {
		keyword := "%" + keyword + "%"
		query = query.Where(`
			(swot_pestel_questions.title_en LIKE ? OR
			swot_pestel_questions.title_ar LIKE ? OR
			swot_pestel_question_groups.title_en LIKE ? OR
			swot_pestel_question_groups.title_ar LIKE ?)`,
			keyword, keyword, keyword, keyword)
	}

	// Apply status filter
	if status == "active" || status == "inactive" {
		query = query.Where("swot_pestel_questions.status = ?", status)
	}

	// Apply is_swot filter
	if isSwot == "true" {
		query = query.Where("swot_pestel_questions.is_swot = ?", true)
	}

	// Apply is_pestel filter
	if isPestel == "true" {
		query = query.Where("swot_pestel_questions.is_pestel = ?", true)
	}

	// Sort by the specified column and order
	if sortBy == "" {
		sortBy = "title_en" // Default sorting
	}
	if sort != "asc" && sort != "desc" {
		sort = "asc"
	}

	query = query.Order(sortBy + " " + sort)

	// Execute query
	err := query.Find(&swotPestelQuestions).Error
	if err != nil {
		return nil, err
	}

	for _, question := range swotPestelQuestions {
		var invitees []*strategyDTOs.SwotPestelQuestionInviteeDTO

		repo.db.Model(&models.SwotPestelQuestionInvitee{}).Where("swot_pestel_question_id = ?", question.ID).Find(&invitees)
		question.SwotPestelQuestionInvitees = invitees
	}

	return swotPestelQuestions, nil
}

func (repo *SwotPestelQuestionRepository) GetAllOrdered(swotPestelQuestionGroupID uint) ([]*models.SwotPestelQuestion, error) {
	var questions []*models.SwotPestelQuestion

	err := repo.db.
		Where("swot_pestel_question_group_id = ?", swotPestelQuestionGroupID).
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

func (repo *SwotPestelQuestionRepository) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.SwotPestelQuestion, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var swotPestelQuestion *models.SwotPestelQuestion

	result := db.
		Preload("SwotPestelQuestionGroup").
		Preload("AnsweredByUser").
		Preload("SwotPestelQuestionInvitees").
		Preload("SwotPestelQuestionInvitees.User").
		First(&swotPestelQuestion, id)

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
		Where("(title_en = ? OR title_ar = ?) AND swot_pestel_question_group_id = ?", titleEn, titleAr, swotPestelQuestionGroupId).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)

	}
	return count > 0, nil
}

func (repo *SwotPestelQuestionRepository) ExistsByTitleExcludingID(swotPestelQuestionGroupId uint, id uint, titleEn, titleAr string) (bool, error) {
	var count int64

	err := repo.db.Model(&models.SwotPestelQuestion{}).
		Where("id <> ? AND swot_pestel_question_group_id = ? AND (title_en = ? OR title_ar = ?)", id, swotPestelQuestionGroupId, titleEn, titleAr).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)
	}

	return count > 0, nil
}

func (repo *SwotPestelQuestionRepository) CreateSwotPestelQuestionInvitees(invitees []*models.SwotPestelQuestionInvitee, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Create(&invitees).Error; err != nil {
		return err
	}

	return nil
}

func (repo *SwotPestelQuestionRepository) DeleteSwotPestelQuestionInvitees(swotPestelQuestionID uint, inviteeIDs []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("swot_pestel_question_id = ? AND user_id IN (?)", swotPestelQuestionID, inviteeIDs).
		Delete(&models.SwotPestelQuestionInvitee{}).Error; err != nil {
		return err
	}

	return nil
}
