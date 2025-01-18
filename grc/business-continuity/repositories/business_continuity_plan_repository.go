package repositories

import (
	"fmt"
	businessContinuityDTOs "jethur/grc/business-continuity/dtos"
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/models"
	"jethur/utility"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type BusinessContinuityPlanRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(params *businessContinuityDTOs.BusinessContinuityPlanListQueryParamsDTO) (*businessContinuityDTOs.BusinessContinuityPlanListResponseDTO, error)
	GetByID(id uint) (*models.BusinessContinuityPlan, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.BusinessContinuityPlan, error)
	Create(businessContinuityPlan *models.BusinessContinuityPlan, tx ...utility.DatabaseTransaction) (uint, error)
	Update(businessContinuityPlan *models.BusinessContinuityPlan, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error
	GenerateReferenceNumber() (string, error)

	GetBusinessContinuityPlanByBusinessFunctionID(businessFunctionId uint) (*models.BusinessContinuityPlan, error)
	CreateBcpProcesses(bcpProcess *models.BcpProcess, tx ...utility.DatabaseTransaction) error
	DeleteBcpProcesses(bcpProcessIds []*uint, tx ...utility.DatabaseTransaction) error
}

type BusinessContinuityPlanRepository struct {
	db *gorm.DB
}

func NewBusinessContinuityPlanRepository(db *gorm.DB) *BusinessContinuityPlanRepository {
	return &BusinessContinuityPlanRepository{
		db: db,
	}
}

// BeginTransaction starts a new transaction.
func (repo *BusinessContinuityPlanRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *BusinessContinuityPlanRepository) GetAll(params *businessContinuityDTOs.BusinessContinuityPlanListQueryParamsDTO) (*businessContinuityDTOs.BusinessContinuityPlanListResponseDTO, error) {
	var businessContinuityPlans []*businessContinuityDTOs.BusinessContinuityPlanListDTO
	var totalCount int64

	query := repo.db.Table("business_continuity_plans").
		Select(`
            business_continuity_plans.*,
			business_functions.title AS business_function_title,
			bcp_statuses.title_en AS bcp_status_title_en,
			bcp_statuses.title_ar AS bcp_status_title_ar,
			bcp_statuses.color_code AS bcp_status_color_code,
			bcp_statuses.type AS bcp_status_type,
			rto_time_unit.title_en AS rto_unit_title_en,
			rto_time_unit.title_ar AS rto_unit_title_ar,
			mtpd_time_unit.title_en AS mtpd_unit_title_en,
			mtpd_time_unit.title_ar AS mtpd_unit_title_ar,
			rpo_time_unit.title_en AS rpo_unit_title_en,
			rpo_time_unit.title_ar AS rpo_unit_title_ar,
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
		    LEFT JOIN business_functions ON business_functions.id = business_continuity_plans.business_function_id
			LEFT JOIN bcp_statuses ON bcp_statuses.id = business_continuity_plans.bcp_status_id
			LEFT JOIN time_units AS rto_time_unit ON rto_time_unit.id = business_continuity_plans.rto_unit_id
			LEFT JOIN time_units AS mtpd_time_unit ON mtpd_time_unit.id = business_continuity_plans.mtpd_unit_id
			LEFT JOIN time_units AS rpo_time_unit ON rpo_time_unit.id = business_continuity_plans.rpo_unit_id
			LEFT JOIN users AS created_user ON created_user.id = business_continuity_plans.created_by_id
			LEFT JOIN users AS updated_user ON updated_user.id = business_continuity_plans.updated_by_id
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
			(business_continuity_plans.reference_number LIKE ? OR
			business_functions.title LIKE ? OR
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

	err := query.Find(&businessContinuityPlans).Error
	if err != nil {
		return nil, err
	}

	totalPages := 1
	if !params.Export {
		totalPages = int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))
	}

	response := &businessContinuityDTOs.BusinessContinuityPlanListResponseDTO{
		Items: businessContinuityPlans,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}
	return response, nil
}

func (repo *BusinessContinuityPlanRepository) GetByID(id uint) (*models.BusinessContinuityPlan, error) {
	var businessContinuityPlan *models.BusinessContinuityPlan

	result := repo.db.First(&businessContinuityPlan, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return businessContinuityPlan, nil

}

func (repo *BusinessContinuityPlanRepository) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.BusinessContinuityPlan, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var businessContinuityPlan *models.BusinessContinuityPlan

	result := db.
		Preload("BusinessFunction").
		Preload("BusinessFunction.SeverityLevel").
		Preload("BusinessFunction.ReviewFrequency").
		Preload("BcpStatus").
		Preload("RtoUnit").
		Preload("MtpdUnit").
		Preload("RpoUnit").
		Preload("Processes").
		Preload("Processes.Process").
		Preload("Processes.Process.Owner").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		First(&businessContinuityPlan, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return businessContinuityPlan, nil
}

func (repo *BusinessContinuityPlanRepository) Create(businessContinuityPlan *models.BusinessContinuityPlan, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(businessContinuityPlan).Error; err != nil {
		return 0, err
	}

	return businessContinuityPlan.ID, nil
}

func (repo *BusinessContinuityPlanRepository) Update(businessContinuityPlan *models.BusinessContinuityPlan, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(businessContinuityPlan).Error; err != nil {
		return err
	}

	return nil
}

func (repo *BusinessContinuityPlanRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	// Delete the Business Function
	if err := db.Delete(&models.BusinessContinuityPlan{}, id).Error; err != nil {
		return err
	}

	return nil
}

func (repo *BusinessContinuityPlanRepository) GenerateReferenceNumber() (string, error) {
	var lastReferenceNumber string
	err := repo.db.Model(&models.BusinessContinuityPlan{}).
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
	newReferenceNumber := fmt.Sprintf("%s-%d-%04d", "BCP", currentYear, lastReferenceID+1)

	return newReferenceNumber, nil

}

func (repo *BusinessContinuityPlanRepository) GetBusinessContinuityPlanByBusinessFunctionID(businessFunctionId uint) (*models.BusinessContinuityPlan, error) {
	var businessContinuityPlan models.BusinessContinuityPlan

	// Query the database to find if an businessContinuityPlan exists with the given BusinessFunctionId
	result := repo.db.Where("business_function_id = ?", businessFunctionId).First(&businessContinuityPlan)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			// Return nil if no businessContinuityPlan is found (not an error, meaning it's safe to create a new businessContinuityPlan)
			return nil, nil
		}
		// Return an error if something went wrong during the query
		return nil, result.Error
	}

	// Return the found audit if it exists
	return &businessContinuityPlan, nil
}

func (repo *BusinessContinuityPlanRepository) CreateBcpProcesses(bcpProcess *models.BcpProcess, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(bcpProcess).Error; err != nil {
		return err
	}
	return nil
}

func (repo *BusinessContinuityPlanRepository) DeleteBcpProcesses(bcpProcessIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	// Delete processes by IDs
	if err := db.Where("id IN (?)", bcpProcessIds).Delete(&models.BcpProcess{}).Error; err != nil {
		return err
	}

	return nil
}
