package repositories

import (
	"fmt"
	"jethur/grc/models"
	strategyDTOs "jethur/grc/strategy/dtos"
	"jethur/utility"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type StrategyProfileRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(params *strategyDTOs.StrategyProfileListQueryParamsDTO) ([]*strategyDTOs.StrategyProfileListDTO, error)
	GetByID(id uint) (*models.StrategyProfile, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.StrategyProfile, error)
	Create(strategyProfile *models.StrategyProfile, tx ...utility.DatabaseTransaction) (uint, error)
	Update(strategyProfile *models.StrategyProfile, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error

	GenerateReferenceNumber() (string, error)
}

type StrategyProfileRepository struct {
	db *gorm.DB
}

func NewStrategyProfileRepository(db *gorm.DB) *StrategyProfileRepository {
	return &StrategyProfileRepository{
		db: db,
	}
}

// BeginTransaction starts a new transaction.
func (repo *StrategyProfileRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *StrategyProfileRepository) GetAll(params *strategyDTOs.StrategyProfileListQueryParamsDTO) ([]*strategyDTOs.StrategyProfileListDTO, error) {
	var strategyProfiles []*strategyDTOs.StrategyProfileListDTO
	var totalCount int64

	query := repo.db.Table("strategy_profiles").
		Select(`
            strategy_profiles.*,
			CONVERT(VARCHAR, strategy_profiles.start_date, 23) AS start_date,
			CONVERT(VARCHAR, strategy_profiles.end_date, 23) AS end_date,
			responsible_user.id AS responsible_user_id,
			responsible_user.first_name_en AS responsible_user_first_name_en,
			responsible_user.first_name_ar AS responsible_user_first_name_ar,
			responsible_user.last_name_en AS responsible_user_last_name_en,
			responsible_user.last_name_ar AS responsible_user_last_name_ar,
			responsible_user.initial AS responsible_user_initial,
			responsible_user.color AS responsible_user_color,
			responsible_user.email AS responsible_user_email,
			responsible_user.profile_image_path AS responsible_user_profile_image_path,
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
			LEFT JOIN users AS responsible_user ON responsible_user.id = strategy_profiles.responsible_user_id
			LEFT JOIN users AS created_user ON created_user.id = strategy_profiles.created_by_id
			LEFT JOIN users AS updated_user ON updated_user.id = strategy_profiles.updated_by_id
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
	// Apply date filters if provided
	if params.StartDate != "" {
		query = query.Where("date >= ?", params.StartDate)
	}
	if params.EndDate != "" {
		query = query.Where("date <= ?", params.EndDate)
	}
	// Apply search keyword if exists
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
			(strategy_profiles.title LIKE ? OR
			created_user.first_name_en LIKE ? OR
			created_user.first_name_ar LIKE ? OR
			created_user.last_name_en LIKE ? OR
			created_user.last_name_ar LIKE ? OR
    		updated_user.first_name_en LIKE ? OR
			updated_user.first_name_ar LIKE ? OR
			updated_user.last_name_en LIKE ? OR
			updated_user.last_name_ar LIKE ?)`,
			keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword)
	}

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	query = query.Order(params.SortBy + " " + params.Sort)

	if !params.Export {
		offset := (params.Page - 1) * params.PageSize
		query = query.Offset(offset).Limit(params.PageSize)
	}

	err := query.Find(&strategyProfiles).Error
	if err != nil {
		return nil, err
	}
	return strategyProfiles, nil
}

func (repo *StrategyProfileRepository) GetByID(id uint) (*models.StrategyProfile, error) {
	var strategyProfile *models.StrategyProfile

	result := repo.db.First(&strategyProfile, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return strategyProfile, nil

}

func (repo *StrategyProfileRepository) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.StrategyProfile, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var strategyProfile *models.StrategyProfile

	result := db.
		Preload("ResponsibleUser").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		First(&strategyProfile, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return strategyProfile, nil
}

func (repo *StrategyProfileRepository) Create(strategyProfile *models.StrategyProfile, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(strategyProfile).Error; err != nil {
		return 0, err
	}

	return strategyProfile.ID, nil
}

func (repo *StrategyProfileRepository) Update(strategyProfile *models.StrategyProfile, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(strategyProfile).Error; err != nil {
		return err
	}

	return nil
}

func (repo *StrategyProfileRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	// Delete the Strategy profile
	if err := db.Delete(&models.StrategyProfile{}, id).Error; err != nil {
		return err
	}

	return nil
}

func (repo *StrategyProfileRepository) GenerateReferenceNumber() (string, error) {
	var lastReferenceNumber string
	err := repo.db.Model(&models.StrategyProfile{}).
		Select("reference_number").
		Order("reference_number desc").
		Limit(1).
		Pluck("reference_number", &lastReferenceNumber).
		Error

	if err != nil {
		return "", err
	}

	lastDigits := "0000"
	if lastReferenceNumber != "" && len(lastReferenceNumber) >= 4 {
		lastDigits = lastReferenceNumber[len(lastReferenceNumber)-4:]
	}

	currentYear := time.Now().Year()
	lastReferenceID, err := strconv.ParseUint(lastDigits, 10, 64)
	if err != nil {
		return "", err
	}
	newReferenceNumber := fmt.Sprintf("%s-%d-%04d", "SP", currentYear, lastReferenceID+1)

	return newReferenceNumber, nil
}
