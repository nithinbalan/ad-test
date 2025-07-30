package repositories

import (
	"errors"
	"fmt"
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/models"
	strategyDTOs "jethur/grc/strategy/dtos"
	"jethur/utility"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type KPIRegisterRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(params *strategyDTOs.KPIRegisterListQueryParamsDTO) (*strategyDTOs.KPIRegisterListResponseDTO, error)
	GetByID(id uint) (*models.KPIRegister, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.KPIRegister, error)
	Create(kpiRegister *models.KPIRegister, tx ...utility.DatabaseTransaction) (uint, error)
	Update(kpiRegister *models.KPIRegister, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error

	CreateKPIRegisterInpuData(inputs []*models.KPIRegisterDataInput, tx ...utility.DatabaseTransaction) error
	DeleteKPIRegisterInputData(kpiRegisterId uint, tx ...utility.DatabaseTransaction) error

	GetTrendAnalysis(kpiId uint, limit int) ([]*models.KPIValueUpdate, error)
	GetLatestValueUpdate(kpiId uint) (*models.KPIValueUpdate, error)
	GenerateReferenceNumber() (string, error)
}

type KPIRegisterRepository struct {
	db *gorm.DB
}

func NewKPIRegisterRepository(db *gorm.DB) *KPIRegisterRepository {
	return &KPIRegisterRepository{
		db: db,
	}
}

// BeginTransaction starts a new transaction.
func (repo *KPIRegisterRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *KPIRegisterRepository) GetAll(params *strategyDTOs.KPIRegisterListQueryParamsDTO) (*strategyDTOs.KPIRegisterListResponseDTO, error) {
	var kpiRegister []*strategyDTOs.KPIRegisterListDTO
	var totalCount int64

	query := repo.db.Table("kpi_register").
		Select(`
            kpi_register.*,
			kpi_categories.title_en AS kpi_category_title_en,
			kpi_categories.title_ar AS kpi_category_title_ar,
			CONVERT(VARCHAR, kpi_register.start_date, 23) AS start_date,
			CONVERT(VARCHAR, kpi_register.end_date, 23) AS target_date,
			responsible_user.id AS responsible_user_id,
			responsible_user.first_name_en AS responsible_user_first_name_en,
			responsible_user.first_name_ar AS responsible_user_first_name_ar,
			responsible_user.last_name_en AS responsible_user_last_name_en,
			responsible_user.last_name_ar AS responsible_user_last_name_ar,
			responsible_user.initial AS responsible_user_initial,
			responsible_user.color AS responsible_user_color,
			responsible_user.email AS responsible_user_email,
			responsible_user.profile_image_path AS responsible_user_profile_image_path,
			kpi_statuses.title_en AS kpi_status_title_en,
			kpi_statuses.title_ar AS kpi_status_title_ar,
			kpi_statuses.type AS kpi_status_type,
			kpi_statuses.color_code AS kpi_status_color_code,
			review_frequencies.title_en AS review_frequency_title_en,
			review_frequencies.title_ar AS review_frequency_title_ar,
			review_frequencies.type AS review_frequency_type,
			calculation_types.title_en AS calculation_type_title_en,
			calculation_types.title_ar AS calculation_type_title_ar,
			calculation_types.type AS calculation_type_type,
			units.title_en AS unit_title_en,
			units.title_ar AS unit_title_ar,
			units.code_en AS unit_code_en,
			units.code_ar AS unit_code_ar,
			units.color_code AS unit_color_code,
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
			LEFT JOIN kpi_categories ON kpi_categories.id = kpi_register.kpi_category_id
			LEFT JOIN kpi_statuses ON kpi_statuses.id = kpi_register.kpi_status_id
			LEFT JOIN calculation_types ON calculation_types.id = kpi_register.calculation_type_id
			LEFT JOIN review_frequencies ON review_frequencies.id = kpi_register.review_frequency_id
			LEFT JOIN units ON units.id = kpi_register.unit_id
			LEFT JOIN users AS responsible_user ON responsible_user.id = kpi_register.responsible_user_id
			LEFT JOIN users AS created_user ON created_user.id = kpi_register.created_by_id
			LEFT JOIN users AS updated_user ON updated_user.id = kpi_register.updated_by_id
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

	// Apply search keyword if exists
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
			(kpi_register.title LIKE ? OR
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

	err := query.Find(&kpiRegister).Error
	if err != nil {
		return nil, err
	}

	totalPages := 1
	if !params.Export {
		totalPages = int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))
	}

	response := &strategyDTOs.KPIRegisterListResponseDTO{
		Items: kpiRegister,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}
	return response, nil
}

func (repo *KPIRegisterRepository) GetByID(id uint) (*models.KPIRegister, error) {
	var kpiRegister *models.KPIRegister

	result := repo.db.First(&kpiRegister, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return kpiRegister, nil

}

func (repo *KPIRegisterRepository) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.KPIRegister, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var kpiRegister *models.KPIRegister

	result := db.
		Preload("KPICategory").
		Preload("KPIStatus").
		Preload("ResponsibleUser").
		Preload("CalculationType").
		Preload("ReviewFrequency").
		Preload("Unit").
		Preload("DataInputs").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		First(&kpiRegister, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return kpiRegister, nil
}

func (repo *KPIRegisterRepository) Create(kpiRegister *models.KPIRegister, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(kpiRegister).Error; err != nil {
		return 0, err
	}

	return kpiRegister.ID, nil
}

func (repo *KPIRegisterRepository) Update(kpiRegister *models.KPIRegister, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(kpiRegister).Error; err != nil {
		return err
	}

	return nil
}

func (repo *KPIRegisterRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	// Delete the Strategy profile
	if err := db.Delete(&models.KPIRegister{}, id).Error; err != nil {
		return err
	}

	return nil
}

func (repo *KPIRegisterRepository) CreateKPIRegisterInpuData(inputDatas []*models.KPIRegisterDataInput, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Create(&inputDatas).Error; err != nil {
		return err
	}

	return nil
}

func (repo *KPIRegisterRepository) DeleteKPIRegisterInputData(kpiRegisterID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("kpi_register_id = ?", kpiRegisterID).
		Delete(&models.KPIRegisterDataInput{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *KPIRegisterRepository) GetTrendAnalysis(kpiID uint, limit int) ([]*models.KPIValueUpdate, error) {
	var kpiValueUpdate []*models.KPIValueUpdate

	result := repo.db.Where("kpi_register_id = ?", kpiID).
		Order("reviewed_date DESC").
		Limit(limit).
		Find(&kpiValueUpdate)

	if result.Error != nil {
		return nil, result.Error
	}

	return kpiValueUpdate, nil
}

func (repo *KPIRegisterRepository) GetLatestValueUpdate(kpiID uint) (*models.KPIValueUpdate, error) {
	var kpiValueUpdate models.KPIValueUpdate

	result := repo.db.Preload("ReviewUser").
		Preload("Documents").
		Preload("Unit").
		Preload("DataInputs").
		Where("kpi_register_id = ? AND is_latest = 1", kpiID).
		First(&kpiValueUpdate)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil, nil when no record is found
		}
		return nil, result.Error
	}

	return &kpiValueUpdate, nil
}
func (repo *KPIRegisterRepository) GenerateReferenceNumber() (string, error) {
	var lastReferenceNumber *string
	err := repo.db.Model(&models.KPIRegister{}).
		Select("reference_number").
		Order("reference_number desc").
		Limit(1).
		Scan(&lastReferenceNumber).
		Error

	if err != nil {
		return "", err
	}

	lastDigits := "0000"
	if lastReferenceNumber != nil && len(*lastReferenceNumber) >= 4 {
		lastDigits = (*lastReferenceNumber)[len(*lastReferenceNumber)-4:]
	}

	currentYear := time.Now().Year()
	lastReferenceID, err := strconv.ParseUint(lastDigits, 10, 64)
	if err != nil {
		return "", err
	}

	newReferenceNumber := fmt.Sprintf("%s-%d-%04d", "KPI", currentYear, lastReferenceID+1)
	return newReferenceNumber, nil
}
