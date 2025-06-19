package repositories

import (
	"errors"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interface
type ThirdPartyQuestionTypeRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ThirdPartyQuestionType, error)
	GetByID(id uint) (*models.ThirdPartyQuestionType, error)
	GetByType(questionType string) (*models.ThirdPartyQuestionType, error)
}

// Implementation
type ThirdPartyQuestionTypeRepository struct {
	db *gorm.DB
}

func NewThirdPartyQuestionTypeRepository(db *gorm.DB) *ThirdPartyQuestionTypeRepository {
	return &ThirdPartyQuestionTypeRepository{
		db: db,
	}
}

func (repo *ThirdPartyQuestionTypeRepository) GetAll(sort string, sortBy string, keyword string, status string) ([]*models.ThirdPartyQuestionType, error) {
	var questionTypes []*models.ThirdPartyQuestionType

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
	result := query.Find(&questionTypes)

	return questionTypes, result.Error
}

func (repo *ThirdPartyQuestionTypeRepository) GetByID(id uint) (*models.ThirdPartyQuestionType, error) {
	var questionType *models.ThirdPartyQuestionType
	result := repo.db.First(&questionType, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, id)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return questionType, result.Error
}

func (repo *ThirdPartyQuestionTypeRepository) GetByType(questionType string) (*models.ThirdPartyQuestionType, error) {
	var tpqt *models.ThirdPartyQuestionType

	result := repo.db.Where("type = ?", questionType).First(&tpqt)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, "Third Party Question Type with type "+questionType)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return tpqt, nil
}
