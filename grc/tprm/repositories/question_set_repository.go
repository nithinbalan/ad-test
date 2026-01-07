package repositories

import (
	commonDTOs "jethur/grc/common/dtos"

	"jethur/grc/models"
	"jethur/grc/tprm/dtos"
	"jethur/utility"

	"gorm.io/gorm"
)

type QuestionSetRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(params *dtos.QuestionSetListQueryParamsDTO) (*dtos.QuestionSetListResponseDTO, error)
	GetByID(id uint) (*models.QuestionSet, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.QuestionSet, error)
	Create(entity *models.QuestionSet, tx ...utility.DatabaseTransaction) (uint, error)
	Update(entity *models.QuestionSet, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error

	GetSectionByID(sectionID uint, tx ...utility.DatabaseTransaction) (*models.QuestionSetSection, error)
	CreateSection(section *models.QuestionSetSection, tx ...utility.DatabaseTransaction) (uint, error)
	UpdateSection(section *models.QuestionSetSection, tx ...utility.DatabaseTransaction) error
	CreateBulk(questions []*models.QuestionSetSectionQuestion, tx ...utility.DatabaseTransaction) error
	DeleteByQuestionSetID(questionSetID uint, tx ...utility.DatabaseTransaction) error
	DeleteByID(id uint, tx ...utility.DatabaseTransaction) error
	GetByQuestionSetID(questionSetID uint) ([]*models.QuestionSetSectionQuestion, error)
	GetSectionIDByQuestionID(questionID uint) (uint, error)
	GetQuestionsBySectionID(sectionID uint, tx ...utility.DatabaseTransaction) ([]*models.QuestionSetSectionQuestion, error)
	DeleteSection(sectionID uint, tx ...utility.DatabaseTransaction) error
	GetByQuestionID(questionSetID uint, tx ...utility.DatabaseTransaction) (*models.QuestionSetSectionQuestion, error)
	UpdateQuestion(question *models.QuestionSetSectionQuestion, tx ...utility.DatabaseTransaction) error
}

type QuestionSetRepository struct {
	db *gorm.DB
}

func NewQuestionSetRepository(db *gorm.DB) *QuestionSetRepository {
	return &QuestionSetRepository{db: db}
}

func (repo *QuestionSetRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *QuestionSetRepository) GetAll(params *dtos.QuestionSetListQueryParamsDTO) (*dtos.QuestionSetListResponseDTO, error) {
	var results []*dtos.QuestionSetListDTO
	var totalCount int64

	subQuery := repo.db.
		Table("question_set_section_questions").
		Select("question_set_id, COUNT(*) as total_questions").
		Group("question_set_id")

	query := repo.db.Table("question_sets").
		Select(`
			question_sets.id,
			question_sets.title_en,
			question_sets.title_ar,
			question_sets.description,
			FORMAT(question_sets.created_at, 'yyyy-MM-dd HH:mm:ss') AS created_at,
			FORMAT(question_sets.updated_at, 'yyyy-MM-dd HH:mm:ss') AS updated_at,
			question_sets.created_by_id,
			question_sets.updated_by_id,
			ISNULL(qsq.total_questions, 0) AS total_questions,

			created_user.first_name_en AS created_user_first_name_en,
			created_user.first_name_ar AS created_user_first_name_ar,
			created_user.last_name_en AS created_user_last_name_en,
			created_user.last_name_ar AS created_user_last_name_ar,
			created_user.initial AS created_user_initial,
			created_user.color AS created_user_color,
			created_user.email AS created_user_email,
			created_user.profile_image_path AS created_user_profile_image_path,

			updated_user.first_name_en AS updated_user_first_name_en,
			updated_user.first_name_ar AS updated_user_first_name_ar,
			updated_user.last_name_en AS updated_user_last_name_en,
			updated_user.last_name_ar AS updated_user_last_name_ar,
			updated_user.initial AS updated_user_initial,
			updated_user.color AS updated_user_color,
			updated_user.email AS updated_user_email,
			updated_user.profile_image_path AS updated_user_profile_image_path
		`).
		Joins("LEFT JOIN (?) AS qsq ON qsq.question_set_id = question_sets.id", subQuery).
		Joins("LEFT JOIN users AS created_user ON created_user.id = question_sets.created_by_id").
		Joins("LEFT JOIN users AS updated_user ON updated_user.id = question_sets.updated_by_id")

	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where("question_sets."+columnName+" IN (?)", values)
		}
	}

	if params.Keyword != "" {
		like := "%" + params.Keyword + "%"
		query = query.Where(`
			question_sets.title_en LIKE ? OR
			question_sets.title_ar LIKE ? OR
			question_sets.description LIKE ?
		`, like, like, like)
	}

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	query = query.Order(params.SortBy + " " + params.Sort)

	if !params.Export {
		offset := (params.Page - 1) * params.PageSize
		query = query.Offset(offset).Limit(params.PageSize)
	}

	if err := query.Find(&results).Error; err != nil {
		return nil, err
	}

	if params.Export {
		return &dtos.QuestionSetListResponseDTO{
			Items: results,
		}, nil
	}
	
	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	return &dtos.QuestionSetListResponseDTO{
		Items: results,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}, nil
}

func (repo *QuestionSetRepository) GetByID(id uint) (*models.QuestionSet, error) {
	var entity *models.QuestionSet
	if err := repo.db.First(&entity, id).Error; err != nil {
		return nil, err
	}
	return entity, nil
}

func (repo *QuestionSetRepository) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.QuestionSet, error) {
	var entity models.QuestionSet

	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	err := db.
		Preload("Sections", func(db *gorm.DB) *gorm.DB {
			return db.Order("[question_set_sections].[order] ASC")
		}).
		Preload("Sections.Questions", func(db *gorm.DB) *gorm.DB {
			return db.Order("[question_set_section_questions].[order] ASC")
		}).
		// Preload("Sections").Order("order ASC").
		// Preload("Sections.Questions").
		Preload("Sections.Questions.QuestionLibrary").
		Preload("Sections.Questions.QuestionLibrary.QuestionCategory").
		Preload("Sections.Questions.QuestionLibrary.QuestionType").
		Preload("Sections.Questions.QuestionLibrary.AnswerOptions").
		Preload("Sections.Questions.QuestionLibrary.AnswerOptions.RiskLevel").
		Preload("Sections.Questions.QuestionLibrary.AnswerOptions.TPRMComplianceResponse").
		// Preload("Questions").
		// Preload("Questions.QuestionLibrary").
		// Preload("Questions.QuestionLibrary.QuestionCategory").
		// Preload("Questions.QuestionLibrary.QuestionType").
		// Preload("Questions.QuestionLibrary.AnswerOptions").
		// Preload("Questions.QuestionLibrary.AnswerOptions.RiskLevel").
		Where("id = ?", id).
		First(&entity).Error

	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (repo *QuestionSetRepository) Create(entity *models.QuestionSet, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(entity).Error; err != nil {
		return 0, err
	}
	return entity.ID, nil
}

func (repo *QuestionSetRepository) Update(entity *models.QuestionSet, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Save(entity).Error
}

func (repo *QuestionSetRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Delete(&models.QuestionSet{}, id).Error
}

func (repo *QuestionSetRepository) GetSectionByID(id uint, tx ...utility.DatabaseTransaction) (*models.QuestionSetSection, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var section *models.QuestionSetSection
	if err := db.First(&section, id).Error; err != nil {
		return nil, err
	}
	return section, nil
}

func (repo *QuestionSetRepository) CreateSection(section *models.QuestionSetSection, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(section).Error; err != nil {
		return 0, err
	}
	return section.ID, nil
}

func (repo *QuestionSetRepository) UpdateSection(section *models.QuestionSetSection, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Save(section).Error
}

func (repo *QuestionSetRepository) CreateBulk(questions []*models.QuestionSetSectionQuestion, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Create(&questions).Error
}

func (repo *QuestionSetRepository) DeleteByQuestionSetID(questionSetID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Where("question_set_id = ?", questionSetID).Delete(&models.QuestionSetSectionQuestion{}).Error
}

func (repo *QuestionSetRepository) DeleteByID(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Delete(&models.QuestionSetSectionQuestion{}, id).Error
}

func (repo *QuestionSetRepository) GetByQuestionSetID(questionSetID uint) ([]*models.QuestionSetSectionQuestion, error) {
	var questions []*models.QuestionSetSectionQuestion
	if err := repo.db.Where("question_set_id = ?", questionSetID).Find(&questions).Error; err != nil {
		return nil, err
	}
	return questions, nil
}

func (repo *QuestionSetRepository) GetSectionIDByQuestionID(questionID uint) (uint, error) {
	var sectionID uint
	query := `SELECT question_set_section_id FROM question_set_section_questions WHERE id = ?`

	err := repo.db.Raw(query, questionID).Scan(&sectionID).Error
	return sectionID, err
}

func (repo *QuestionSetRepository) DeleteSection(sectionID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Where("id = ?", sectionID).Delete(&models.QuestionSetSection{}).Error
}

func (repo *QuestionSetRepository) GetQuestionsBySectionID(sectionID uint, tx ...utility.DatabaseTransaction) ([]*models.QuestionSetSectionQuestion, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var questions []*models.QuestionSetSectionQuestion
	if err := db.Where("question_set_section_id = ?", sectionID).Find(&questions).Error; err != nil {
		return nil, err
	}
	return questions, nil
}

func (repo *QuestionSetRepository) GetByQuestionID(questionID uint, tx ...utility.DatabaseTransaction) (*models.QuestionSetSectionQuestion, error) {
	var question *models.QuestionSetSectionQuestion
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Where("id = ?", questionID).Find(&question).Error; err != nil {
		return nil, err
	}
	return question, nil
}

func (repo *QuestionSetRepository) UpdateQuestion(question *models.QuestionSetSectionQuestion, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Save(question).Error
}
