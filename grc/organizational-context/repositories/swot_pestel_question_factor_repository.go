package repositories

import (
	"jethur/grc/models"
	organizationalContextDTOs "jethur/grc/organizational-context/dtos"
	"jethur/utility"

	"gorm.io/gorm"
)

type SwotPestelQuestionFactorRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(questionId uint, params *organizationalContextDTOs.SwotPestelQuestionFactorListQueryParamsDTO) ([]*organizationalContextDTOs.SwotPestelQuestionFactorListDTO, error)
	GetAllFactors(params *organizationalContextDTOs.SwotPestelQuestionFactorListQueryParamsDTO) ([]*organizationalContextDTOs.SwotPestelQuestionFactorListDTO, error)
	GetByID(id uint) (*models.SwotPestelQuestionFactor, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.SwotPestelQuestionFactor, error)
	Create(swotPestelQuestionFactor *models.SwotPestelQuestionFactor, tx ...utility.DatabaseTransaction) (uint, error)
	Update(swotPestelQuestionFactor *models.SwotPestelQuestionFactor, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error
}

type SwotPestelQuestionFactorRepository struct {
	db *gorm.DB
}

func NewSwotPestelQuestionFactorRepository(db *gorm.DB) *SwotPestelQuestionFactorRepository {
	return &SwotPestelQuestionFactorRepository{
		db: db,
	}
}

// BeginTransaction starts a new transaction.
func (repo *SwotPestelQuestionFactorRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *SwotPestelQuestionFactorRepository) GetAll(questionId uint, params *organizationalContextDTOs.SwotPestelQuestionFactorListQueryParamsDTO) ([]*organizationalContextDTOs.SwotPestelQuestionFactorListDTO, error) {
	var swotPestelQuestionFactors []*organizationalContextDTOs.SwotPestelQuestionFactorListDTO
	var totalCount int64

	query := repo.db.Table("swot_pestel_question_factors").
		Select(`
            swot_pestel_question_factors.*,
			swot_pestel_questions.title_en AS swot_pestel_question_title_en,
			swot_pestel_questions.title_ar AS swot_pestel_question_title_ar,
			swot_pestel_categories.title_en AS swot_pestel_category_title_en,
			swot_pestel_categories.title_ar AS swot_pestel_category_title_ar,
			swot_pestel_categories.color_code AS swot_pestel_category_color_code,
			swot_pestel_categories.type AS swot_pestel_category_type,
			created_user.id AS created_user_id,
        	created_user.first_name_en AS created_user_first_name_en,
            created_user.first_name_ar AS created_user_first_name_ar,
            created_user.last_name_en AS created_user_last_name_en,
            created_user.last_name_ar AS created_user_last_name_ar,
			created_user.initial AS created_user_initial,
            created_user.color AS created_user_color,
            created_user.email AS created_user_email,
            created_user.profile_image_path AS created_user_profile_image_path,   
			updated_user.id AS updated_user_id,
    		updated_user.first_name_en AS updated_user_first_name_en,
            updated_user.first_name_ar AS updated_user_first_name_ar,
            updated_user.last_name_en AS updated_user_last_name_en,
            updated_user.last_name_ar AS updated_user_last_name_ar,
			updated_user.initial AS updated_user_initial,
            updated_user.color AS updated_user_color,
            updated_user.email AS updated_user_email,
            updated_user.profile_image_path AS updated_user_profile_image_path  
        `).
		Joins(`
		    LEFT JOIN swot_pestel_questions ON swot_pestel_questions.id = swot_pestel_question_factors.swot_pestel_question_id
			LEFT JOIN swot_pestel_categories ON swot_pestel_categories.id = swot_pestel_question_factors.swot_pestel_category_id
			LEFT JOIN users AS created_user ON created_user.id = swot_pestel_question_factors.created_by_id
			LEFT JOIN users AS updated_user ON updated_user.id = swot_pestel_question_factors.updated_by_id
        `).Where("swot_pestel_question_factors.swot_pestel_question_id = ?", questionId)
	// Apply dynamic filters to the query
	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where(columnName+" IN (?)", values)
		}
	}
	// Calculate total count for pagination
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}
	// Apply search keyword if exists
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
			(swot_pestel_question_factors.reference_number LIKE ? OR
			swot_pestel_questions.title LIKE ? OR
			created_user.first_name_en LIKE ? OR
			created_user.first_name_ar LIKE ? OR
			created_user.last_name_en LIKE ? OR
			created_user.last_name_ar LIKE ? OR
    		updated_user.first_name_en LIKE ? OR
			updated_user.first_name_ar LIKE ? OR
			updated_user.last_name_en LIKE ? OR
			updated_user.last_name_ar LIKE ?)`,
			keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword,
			keyword)
	}

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	query = query.Order(params.SortBy + " " + params.Sort)

	if !params.Export {
		offset := (params.Page - 1) * params.PageSize
		query = query.Offset(offset).Limit(params.PageSize)
	}

	err := query.Find(&swotPestelQuestionFactors).Error
	if err != nil {
		return nil, err
	}
	return swotPestelQuestionFactors, nil
}

func (repo *SwotPestelQuestionFactorRepository) GetAllFactors(params *organizationalContextDTOs.SwotPestelQuestionFactorListQueryParamsDTO) ([]*organizationalContextDTOs.SwotPestelQuestionFactorListDTO, error) {
	var swotPestelQuestionFactors []*organizationalContextDTOs.SwotPestelQuestionFactorListDTO
	var totalCount int64

	query := repo.db.Table("swot_pestel_question_factors").
		Select(`
            swot_pestel_question_factors.*,
			swot_pestel_questions.title_en AS swot_pestel_question_title_en,
			swot_pestel_questions.title_ar AS swot_pestel_question_title_ar,
			swot_pestel_questions.answer AS swot_pestel_question_answer,
			swot_pestel_categories.title_en AS swot_pestel_category_title_en,
			swot_pestel_categories.title_ar AS swot_pestel_category_title_ar,
			swot_pestel_categories.color_code AS swot_pestel_category_color_code,
			swot_pestel_categories.type AS swot_pestel_category_type,
			created_user.id AS created_user_id,
        	created_user.first_name_en AS created_user_first_name_en,
            created_user.first_name_ar AS created_user_first_name_ar,
            created_user.last_name_en AS created_user_last_name_en,
            created_user.last_name_ar AS created_user_last_name_ar,
			created_user.initial AS created_user_initial,
            created_user.color AS created_user_color,
            created_user.email AS created_user_email,
            created_user.profile_image_path AS created_user_profile_image_path,   
			updated_user.id AS updated_user_id,
    		updated_user.first_name_en AS updated_user_first_name_en,
            updated_user.first_name_ar AS updated_user_first_name_ar,
            updated_user.last_name_en AS updated_user_last_name_en,
            updated_user.last_name_ar AS updated_user_last_name_ar,
			updated_user.initial AS updated_user_initial,
            updated_user.color AS updated_user_color,
            updated_user.email AS updated_user_email,
            updated_user.profile_image_path AS updated_user_profile_image_path  
        `).
		Joins(`
		    LEFT JOIN swot_pestel_questions ON swot_pestel_questions.id = swot_pestel_question_factors.swot_pestel_question_id
			LEFT JOIN swot_pestel_categories ON swot_pestel_categories.id = swot_pestel_question_factors.swot_pestel_category_id
			LEFT JOIN users AS created_user ON created_user.id = swot_pestel_question_factors.created_by_id
			LEFT JOIN users AS updated_user ON updated_user.id = swot_pestel_question_factors.updated_by_id
        `)
	// Apply dynamic filters to the query
	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where(columnName+" IN (?)", values)
		}
	}
	// Calculate total count for pagination
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	if params.CategoryType != "" {
		query = query.Where("swot_pestel_categories.type = ?", params.CategoryType)
	}

	// Apply search keyword if exists
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
			(swot_pestel_question_factors.reference_number LIKE ? OR
			swot_pestel_questions.title LIKE ? OR
			created_user.first_name_en LIKE ? OR
			created_user.first_name_ar LIKE ? OR
			created_user.last_name_en LIKE ? OR
			created_user.last_name_ar LIKE ? OR
    		updated_user.first_name_en LIKE ? OR
			updated_user.first_name_ar LIKE ? OR
			updated_user.last_name_en LIKE ? OR
			updated_user.last_name_ar LIKE ?)`,
			keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword,
			keyword)
	}

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	query = query.Order(params.SortBy + " " + params.Sort)

	if !params.Export {
		offset := (params.Page - 1) * params.PageSize
		query = query.Offset(offset).Limit(params.PageSize)
	}

	err := query.Find(&swotPestelQuestionFactors).Error
	if err != nil {
		return nil, err
	}
	return swotPestelQuestionFactors, nil
}

func (repo *SwotPestelQuestionFactorRepository) GetByID(id uint) (*models.SwotPestelQuestionFactor, error) {
	var swotPestelQuestionFactor *models.SwotPestelQuestionFactor

	result := repo.db.First(&swotPestelQuestionFactor, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return swotPestelQuestionFactor, nil

}

func (repo *SwotPestelQuestionFactorRepository) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.SwotPestelQuestionFactor, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var swotPestelQuestionFactor *models.SwotPestelQuestionFactor

	result := db.
		Preload("SwotPestelQuestion").
		Preload("SwotPestelCategory").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		First(&swotPestelQuestionFactor, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return swotPestelQuestionFactor, nil
}

func (repo *SwotPestelQuestionFactorRepository) Create(swotPestelQuestionFactor *models.SwotPestelQuestionFactor, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(swotPestelQuestionFactor).Error; err != nil {
		return 0, err
	}

	return swotPestelQuestionFactor.ID, nil
}

func (repo *SwotPestelQuestionFactorRepository) Update(swotPestelQuestionFactor *models.SwotPestelQuestionFactor, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(swotPestelQuestionFactor).Error; err != nil {
		return err
	}

	return nil
}

func (repo *SwotPestelQuestionFactorRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	// Delete the Swot Pestel Question Factor
	if err := db.Delete(&models.SwotPestelQuestionFactor{}, id).Error; err != nil {
		return err
	}

	return nil
}
