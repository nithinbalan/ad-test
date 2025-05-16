package repositories

import (
	"fmt"
	"strconv"
	"time"

	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/models"
	thirdPartyDTOs "jethur/grc/tprm/dtos"
	"jethur/utility"

	"gorm.io/gorm"
)

type ThirdPartyAssessmentRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(params *thirdPartyDTOs.ThirdPartyAssessmentListQueryParamsDTO) (*thirdPartyDTOs.ThirdPartyAssessmentListResponseDTO, error)
	GetByID(id uint) (*models.ThirdPartyAssessment, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.ThirdPartyAssessment, error)
	Create(assessments *models.ThirdPartyAssessment, tx ...utility.DatabaseTransaction) (uint, error)
	Update(assessments *models.ThirdPartyAssessment, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error
	GenerateReferenceNumber() (string, error)
}

type ThirdPartyAssessmentRepository struct {
	db *gorm.DB
}

func NewThirdPartyAssessmentRepository(db *gorm.DB) *ThirdPartyAssessmentRepository {
	return &ThirdPartyAssessmentRepository{db: db}
}

func (repo *ThirdPartyAssessmentRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *ThirdPartyAssessmentRepository) GetAll(params *thirdPartyDTOs.ThirdPartyAssessmentListQueryParamsDTO) (*thirdPartyDTOs.ThirdPartyAssessmentListResponseDTO, error) {
	var results []*thirdPartyDTOs.ThirdPartyAssessmentListDTO
	var totalCount int64

	query := repo.db.Table("third_party_assessments").
		Select(`
			third_party_assessments.*,
			CONVERT(VARCHAR, due_date, 23) AS due_date,
			question_sets.id AS question_set_id,
			question_sets.title_en AS question_set_title_en,
			question_sets.title_ar AS question_set_title_ar,
			third_party_register.title AS third_party_register_title,
			created_user.first_name_en AS created_user_first_name_en,
			created_user.last_name_en AS created_user_last_name_en,
			updated_user.first_name_en AS updated_user_first_name_en,
			updated_user.last_name_en AS updated_user_last_name_en
		`).
		Joins("LEFT JOIN third_party_register ON third_party_register.id = third_party_assessments.third_party_register_id").
		Joins("LEFT JOIN question_sets ON question_sets.id = third_party_assessments.question_set_id").
		Joins("LEFT JOIN users AS created_user ON created_user.id = third_party_assessments.created_by_id").
		Joins("LEFT JOIN users AS updated_user ON updated_user.id = third_party_assessments.updated_by_id")

	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where(columnName+" IN (?)", values)
		}
	}

	if params.Keyword != "" {
		likePattern := "%" + params.Keyword + "%"
		query = query.Where(`
			third_party_assessments.title LIKE ? OR
			third_party_assessments.reference_number LIKE ? OR
			third_party_assessments.due_date LIKE ? OR
	
			third_party_register.title LIKE ? OR
	
			question_sets.title_en LIKE ? OR
			question_sets.title_ar LIKE ? OR
		`,
			likePattern, likePattern,
			likePattern, likePattern,
			likePattern, likePattern,
		)
	}

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	query = query.Order(params.SortBy + " " + params.Sort)
	offset := (params.Page - 1) * params.PageSize
	query = query.Offset(offset).Limit(params.PageSize)

	err := query.Find(&results).Error
	if err != nil {
		return nil, err
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	response := &thirdPartyDTOs.ThirdPartyAssessmentListResponseDTO{
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

func (repo *ThirdPartyAssessmentRepository) GetByID(id uint) (*models.ThirdPartyAssessment, error) {
	var entity *models.ThirdPartyAssessment
	if err := repo.db.First(&entity, id).Error; err != nil {
		return nil, err
	}
	return entity, nil
}

func (repo *ThirdPartyAssessmentRepository) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.ThirdPartyAssessment, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var entity *models.ThirdPartyAssessment
	if err := db.Preload("ThirdPartyRegister").
		Preload("ThirdPartyRegister.ThirdPartyType").
		Preload("ThirdPartyRegister.Country").
		Preload("ThirdPartyRegister.OrganizationUnit").
		Preload("ThirdPartyRegister.ResponsibleUser").
		Preload("ThirdPartyRegister.ResponsibleUnit").
		Preload("ThirdPartyRegister.SeverityLevel").
		Preload("QuestionSet").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		First(&entity, id).Error; err != nil {
		return nil, err
	}

	return entity, nil
}

func (repo *ThirdPartyAssessmentRepository) Create(entity *models.ThirdPartyAssessment, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Create(entity).Error; err != nil {
		return 0, err
	}
	return entity.ID, nil
}

func (repo *ThirdPartyAssessmentRepository) Update(entity *models.ThirdPartyAssessment, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Save(entity).Error
}

func (repo *ThirdPartyAssessmentRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Delete(&models.ThirdPartyAssessment{}, id).Error
}

func (repo *ThirdPartyAssessmentRepository) GenerateReferenceNumber() (string, error) {
	var lastReference string
	if err := repo.db.Model(&models.ThirdPartyAssessment{}).Select("reference_number").Order("id desc").Limit(1).Pluck("reference_number", &lastReference).Error; err != nil {
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

	ref := fmt.Sprintf("TPA-%d-%04d", time.Now().Year(), lastNum+1)
	return ref, nil
}
