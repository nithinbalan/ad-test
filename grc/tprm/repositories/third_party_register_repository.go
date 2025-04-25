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

type ThirdPartyRegisterRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(params *thirdPartyDTOs.ThirdPartyRegisterListQueryParamsDTO) (*thirdPartyDTOs.ThirdPartyRegisterListResponseDTO, error)
	GetByID(id uint) (*models.ThirdPartyRegister, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.ThirdPartyRegister, error)
	Create(register *models.ThirdPartyRegister, tx ...utility.DatabaseTransaction) (uint, error)
	Update(register *models.ThirdPartyRegister, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error
	GenerateReferenceNumber() (string, error)
}

type ThirdPartyRegisterRepository struct {
	db *gorm.DB
}

func NewThirdPartyRegisterRepository(db *gorm.DB) *ThirdPartyRegisterRepository {
	return &ThirdPartyRegisterRepository{db: db}
}

func (repo *ThirdPartyRegisterRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *ThirdPartyRegisterRepository) GetAll(params *thirdPartyDTOs.ThirdPartyRegisterListQueryParamsDTO) (*thirdPartyDTOs.ThirdPartyRegisterListResponseDTO, error) {
	var results []*thirdPartyDTOs.ThirdPartyRegisterListDTO
	var totalCount int64

	query := repo.db.Table("third_party_register").
		Select(`
			third_party_register.*,
			CONVERT(VARCHAR, onboarding_date, 23) AS onboarding_date,
			CONVERT(VARCHAR, service_start_date, 23) AS service_start_date,
			CONVERT(VARCHAR, service_end_date, 23) AS service_end_date,
			countries.id AS country_id,
			countries.title_en AS country_title_en,
			countries.title_ar AS country_title_ar,
			third_party_types.title_en AS third_party_type_title_en,
			third_party_types.title_ar AS third_party_type_title_ar,
			organization_units.title_en AS organization_unit_title_en,
			organization_units.title_ar AS organization_unit_title_ar,
			severity_levels.title_en AS severity_level_title_en,
			severity_levels.title_ar AS severity_level_title_ar,
			created_user.first_name_en AS created_user_first_name_en,
			created_user.last_name_en AS created_user_last_name_en,
			updated_user.first_name_en AS updated_user_first_name_en,
			updated_user.last_name_en AS updated_user_last_name_en
		`).
		Joins("LEFT JOIN third_party_types ON third_party_types.id = third_party_register.third_party_type_id").
		Joins("LEFT JOIN countries ON countries.id = third_party_register.country_id").
		Joins("LEFT JOIN organization_units ON organization_units.id = third_party_register.organization_unit_id").
		Joins("LEFT JOIN severity_levels ON severity_levels.id = third_party_register.severity_level_id").
		Joins("LEFT JOIN users AS created_user ON created_user.id = third_party_register.created_by_id").
		Joins("LEFT JOIN users AS updated_user ON updated_user.id = third_party_register.updated_by_id")

	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where(columnName+" IN (?)", values)
		}
	}

	if params.Keyword != "" {
		query = query.Where(`
				third_party_register.title LIKE ? OR
				third_party_register.reference_number LIKE ? OR
				third_party_register.website LIKE ? OR
				countries.title_en LIKE ? OR
				countries.title_ar LIKE ?`,
			"%"+params.Keyword+"%",
			"%"+params.Keyword+"%",
			"%"+params.Keyword+"%",
			"%"+params.Keyword+"%",
			"%"+params.Keyword+"%",
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

	response := &thirdPartyDTOs.ThirdPartyRegisterListResponseDTO{
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

func (repo *ThirdPartyRegisterRepository) GetByID(id uint) (*models.ThirdPartyRegister, error) {
	var entity *models.ThirdPartyRegister
	if err := repo.db.First(&entity, id).Error; err != nil {
		return nil, err
	}
	return entity, nil
}

func (repo *ThirdPartyRegisterRepository) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.ThirdPartyRegister, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var entity *models.ThirdPartyRegister
	if err := db.Preload("ThirdPartyType").
		Preload("Country").
		Preload("OrganizationUnit").
		Preload("SeverityLevel").
		Preload("ResponsibleUser").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		First(&entity, id).Error; err != nil {
		return nil, err
	}

	return entity, nil
}

func (repo *ThirdPartyRegisterRepository) Create(entity *models.ThirdPartyRegister, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Create(entity).Error; err != nil {
		return 0, err
	}
	return entity.ID, nil
}

func (repo *ThirdPartyRegisterRepository) Update(entity *models.ThirdPartyRegister, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Save(entity).Error
}

func (repo *ThirdPartyRegisterRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Delete(&models.ThirdPartyRegister{}, id).Error
}

func (repo *ThirdPartyRegisterRepository) GenerateReferenceNumber() (string, error) {
	var lastReference string
	if err := repo.db.Model(&models.ThirdPartyRegister{}).Select("reference_number").Order("id desc").Limit(1).Pluck("reference_number", &lastReference).Error; err != nil {
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

	ref := fmt.Sprintf("TPR-%d-%04d", time.Now().Year(), lastNum+1)
	return ref, nil
}
