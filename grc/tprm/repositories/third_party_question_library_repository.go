package repositories

import (
	"fmt"
	commonDTOs "jethur/grc/common/dtos"
	"strconv"
	"time"

	"jethur/grc/models"
	"jethur/grc/tprm/dtos"
	"jethur/utility"

	"gorm.io/gorm"
)

type ThirdPartyQuestionLibraryRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(params *dtos.ThirdPartyQuestionLibraryListQueryParamsDTO) (*dtos.ThirdPartyQuestionLibraryListResponseDTO, error)
	GetByID(id uint) (*models.ThirdPartyQuestionLibrary, error)
	GetFullDetailsByID(id uint) (*models.ThirdPartyQuestionLibrary, error)
	Create(entity *models.ThirdPartyQuestionLibrary, tx ...utility.DatabaseTransaction) (uint, error)
	Update(entity *models.ThirdPartyQuestionLibrary, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error
	GenerateReferenceNumber() (string, error)

	CreateAnswerOptions(options []*models.ThirdPartyQuestionLibraryAnswerOption, tx ...utility.DatabaseTransaction) error
	DeleteAnswerOptionsByLibraryID(libraryID uint, tx ...utility.DatabaseTransaction) error
}

type ThirdPartyQuestionLibraryRepository struct {
	db *gorm.DB
}

func NewThirdPartyQuestionLibraryRepository(db *gorm.DB) *ThirdPartyQuestionLibraryRepository {
	return &ThirdPartyQuestionLibraryRepository{db: db}
}

func (repo *ThirdPartyQuestionLibraryRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *ThirdPartyQuestionLibraryRepository) GetAll(params *dtos.ThirdPartyQuestionLibraryListQueryParamsDTO) (*dtos.ThirdPartyQuestionLibraryListResponseDTO, error) {
	var results []*dtos.ThirdPartyQuestionLibraryListDTO
	var totalCount int64

	query := repo.db.Table("third_party_question_libraries").
		Select(`
			third_party_question_libraries.*,
			third_party_question_categories.title_en AS third_party_question_category_title_en,
			third_party_question_categories.title_ar AS third_party_question_category_title_ar,
			third_party_question_types.title_en AS third_party_question_type_title_en,
			third_party_question_types.title_ar AS third_party_question_type_title_ar
		`).
		Joins("LEFT JOIN third_party_question_categories ON third_party_question_categories.id = third_party_question_libraries.third_party_question_category_id").
		Joins("LEFT JOIN third_party_question_types ON third_party_question_types.id = third_party_question_libraries.third_party_question_type_id")

	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
			third_party_question_libraries.question_title_en LIKE ? OR
			third_party_question_libraries.question_title_ar LIKE ?
		`, keyword, keyword)
	}

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	query = query.Order(params.SortBy + " " + params.Sort)
	offset := (params.Page - 1) * params.PageSize
	query = query.Offset(offset).Limit(params.PageSize)

	if err := query.Find(&results).Error; err != nil {
		return nil, err
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	response := &dtos.ThirdPartyQuestionLibraryListResponseDTO{
		Items: results,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}
	return response, nil
}

func (repo *ThirdPartyQuestionLibraryRepository) GetByID(id uint) (*models.ThirdPartyQuestionLibrary, error) {
	var entity *models.ThirdPartyQuestionLibrary
	if err := repo.db.First(&entity, id).Error; err != nil {
		return nil, err
	}
	return entity, nil
}

func (repo *ThirdPartyQuestionLibraryRepository) GetFullDetailsByID(id uint) (*models.ThirdPartyQuestionLibrary, error) {
	var entity *models.ThirdPartyQuestionLibrary
	if err := repo.db.Preload("QuestionCategory").Preload("QuestionType").First(&entity, id).Error; err != nil {
		return nil, err
	}
	return entity, nil
}

func (repo *ThirdPartyQuestionLibraryRepository) Create(entity *models.ThirdPartyQuestionLibrary, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(entity).Error; err != nil {
		return 0, err
	}
	return entity.ID, nil
}

func (repo *ThirdPartyQuestionLibraryRepository) Update(entity *models.ThirdPartyQuestionLibrary, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Save(entity).Error
}

func (repo *ThirdPartyQuestionLibraryRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Delete(&models.ThirdPartyQuestionLibrary{}, id).Error
}

func (repo *ThirdPartyQuestionLibraryRepository) GenerateReferenceNumber() (string, error) {
	var lastReference string
	if err := repo.db.Model(&models.ThirdPartyQuestionLibrary{}).Select("reference_number").Order("id desc").Limit(1).Pluck("reference_number", &lastReference).Error; err != nil {
		return "", err
	}

	lastDigits := "0000"
	if lastReference != "" && len(lastReference) >= 4 {
		lastDigits = lastReference[len(lastReference)-4:]
	}
	lastNum, err := strconv.Atoi(lastDigits)
	if err != nil {
		return "", err
	}

	ref := fmt.Sprintf("TPQL-%d-%04d", time.Now().Year(), lastNum+1)
	return ref, nil
}

func (repo *ThirdPartyQuestionLibraryRepository) CreateAnswerOptions(options []*models.ThirdPartyQuestionLibraryAnswerOption, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Create(&options).Error
}

func (repo *ThirdPartyQuestionLibraryRepository) DeleteAnswerOptionsByLibraryID(libraryID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Where("third_party_question_library_id = ?", libraryID).Delete(&models.ThirdPartyQuestionLibraryAnswerOption{}).Error
}
