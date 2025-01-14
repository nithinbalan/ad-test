package repositories

import (
	businessContinuityDTOs "jethur/grc/business-continuity/dtos"
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type BcpCommunicationPlanRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(bcpId uint, params *businessContinuityDTOs.BcpCommunicationPlanListQueryParamsDTO) (*businessContinuityDTOs.BcpCommunicationPlanListResponseDTO, error)
	GetByID(id uint) (*models.BcpCommunicationPlan, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.BcpCommunicationPlan, error)
	Create(BcpCommunicationPlan *models.BcpCommunicationPlan, tx ...utility.DatabaseTransaction) (uint, error)
	Update(BcpCommunicationPlan *models.BcpCommunicationPlan, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error
}

type BcpCommunicationPlanRepository struct {
	db *gorm.DB
}

func NewBcpCommunicationPlanRepository(db *gorm.DB) *BcpCommunicationPlanRepository {
	return &BcpCommunicationPlanRepository{
		db: db,
	}
}

// BeginTransaction starts a new transaction.
func (repo *BcpCommunicationPlanRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *BcpCommunicationPlanRepository) GetAll(bcpId uint, params *businessContinuityDTOs.BcpCommunicationPlanListQueryParamsDTO) (*businessContinuityDTOs.BcpCommunicationPlanListResponseDTO, error) {
	var BcpCommunicationPlans []*businessContinuityDTOs.BcpCommunicationPlanListDTO
	var totalCount int64

	query := repo.db.Table("bcp_communication_plans").
		Select(`
            bcp_communication_plans.*,
			bcp_phases.title_en AS bcp_phase_title_en,
			bcp_phases.title_ar AS bcp_phase_title_ar,
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
			LEFT JOIN bcp_phases ON bcp_phases.id = bcp_communication_plans.bcp_phase_id
		    LEFT JOIN users AS responsible_user ON responsible_user.id = bcp_communication_plans.responsible_user_id
			LEFT JOIN users AS created_user ON created_user.id = bcp_communication_plans.created_by_id
			LEFT JOIN users AS updated_user ON updated_user.id = bcp_communication_plans.updated_by_id
        `).Where("bcp_communication_plans.business_continuity_plan_id = ?", bcpId)
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
			(bcp_communication_plans.title LIKE ? OR
			bcp_phases.title_en LIKE ? OR
			bcp_phases.title_ar LIKE ? OR
			responsible_user.first_name_en LIKE ? OR
			responsible_user.first_name_ar LIKE ? OR
			responsible_user.last_name_en LIKE ? OR
			responsible_user.last_name_ar LIKE ? OR
			responsible_user.email LIKE ? OR
            responsible_user.initial LIKE ? OR
            responsible_user.color LIKE ? OR
            responsible_user.profile_image_path LIKE ? OR
			created_user.first_name_en LIKE ? OR
			created_user.first_name_ar LIKE ? OR
			created_user.last_name_en LIKE ? OR
			created_user.last_name_ar LIKE ? OR
    		updated_user.first_name_en LIKE ? OR
			updated_user.first_name_ar LIKE ? OR
			updated_user.last_name_en LIKE ? OR
			updated_user.last_name_ar LIKE ?)`,
			keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword,
			keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword,
			keyword, keyword, keyword, keyword, keyword)
	}

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	query = query.Order(params.SortBy + " " + params.Sort)

	offset := (params.Page - 1) * params.PageSize
	query = query.Offset(offset).Limit(params.PageSize)

	err := query.Find(&BcpCommunicationPlans).Error
	if err != nil {
		return nil, err
	}

	totalPages := 1
	totalPages = int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	response := &businessContinuityDTOs.BcpCommunicationPlanListResponseDTO{
		Items: BcpCommunicationPlans,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}
	return response, nil
}

func (repo *BcpCommunicationPlanRepository) GetByID(id uint) (*models.BcpCommunicationPlan, error) {
	var BcpCommunicationPlan *models.BcpCommunicationPlan

	result := repo.db.First(&BcpCommunicationPlan, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return BcpCommunicationPlan, nil

}

func (repo *BcpCommunicationPlanRepository) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.BcpCommunicationPlan, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var BcpCommunicationPlan *models.BcpCommunicationPlan

	result := db.
		Preload("BusinessContinuityPlan").
		Preload("BusinessContinuityPlan.BusinessFunction").
		Preload("ResponsibleUser").
		Preload("ResponsibleUser.OrganizationUnit1").
		Preload("BcpPhase").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		First(&BcpCommunicationPlan, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return BcpCommunicationPlan, nil
}

func (repo *BcpCommunicationPlanRepository) Create(BcpCommunicationPlan *models.BcpCommunicationPlan, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(BcpCommunicationPlan).Error; err != nil {
		return 0, err
	}

	return BcpCommunicationPlan.ID, nil
}

func (repo *BcpCommunicationPlanRepository) Update(BcpCommunicationPlan *models.BcpCommunicationPlan, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(BcpCommunicationPlan).Error; err != nil {
		return err
	}

	return nil
}

func (repo *BcpCommunicationPlanRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	// Delete the Business Communication Plan
	if err := db.Delete(&models.BcpCommunicationPlan{}, id).Error; err != nil {
		return err
	}

	return nil
}
