package repositories

import (
	businessContinuityDTOs "jethur/grc/business-continuity/dtos"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type BcpBiaMatrixRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAllKeyProcess(bcpId uint, params *businessContinuityDTOs.BcpKeyProcessListQueryParamsDTO) ([]*businessContinuityDTOs.BcpKeyProcessListDTO, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.BusinessContinuityPlan, error)
	GetByID(bcpProcessId uint) (*models.BcpProcess, error)
	Update(bcpProcess *models.BcpProcess, tx ...utility.DatabaseTransaction) error
}

type BcpBiaMatrixRepository struct {
	db *gorm.DB
}

func NewBcpBiaMatrixRepository(db *gorm.DB) *BcpBiaMatrixRepository {
	return &BcpBiaMatrixRepository{
		db: db,
	}
}

// BeginTransaction starts a new transaction.
func (repo *BcpBiaMatrixRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *BcpBiaMatrixRepository) GetAllKeyProcess(bcpId uint, params *businessContinuityDTOs.BcpKeyProcessListQueryParamsDTO) ([]*businessContinuityDTOs.BcpKeyProcessListDTO, error) {
	var BcpKeyProcesss []*businessContinuityDTOs.BcpKeyProcessListDTO
	var totalCount int64

	query := repo.db.Table("bcp_processes").
		Select(`
            bcp_processes.*,
            processes.title_en AS process_title_en,
            processes.title_ar AS process_title_ar,
			processes.description AS process_description,
			organization_units.title_en AS process_owner_unit_title_en,
			organization_units.title_ar AS process_owner_unit_title_ar,
			owner.id AS process_owner_id,
			owner.first_name_en AS process_owner_first_name_en,
			owner.first_name_ar AS process_owner_first_name_ar,
			owner.last_name_en AS process_owner_last_name_en,
			owner.last_name_ar AS process_owner_last_name_ar,
			owner.initial AS process_owner_initial,
			owner.color AS process_owner_color,
			owner.email AS process_owner_email,
			owner.profile_image_path AS process_owner_profile_image_path
            
        `).
		Joins(`
			LEFT JOIN processes ON processes.id = bcp_processes.process_id
			LEFT JOIN users AS owner ON owner.id = processes.owner_id
			LEFT JOIN organization_units ON organization_units.id = processes.owner_unit_id
        `).Where("bcp_processes.business_continuity_plan_id = ?", bcpId)
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
			(processes.title_en LIKE ? OR
			processes.title_ar LIKE ?)`,
			keyword, keyword)
	}

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	query = query.Order(params.SortBy + " " + params.Sort)

	err := query.Find(&BcpKeyProcesss).Error
	if err != nil {
		return nil, err
	}

	return BcpKeyProcesss, nil
}

func (repo *BcpBiaMatrixRepository) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.BusinessContinuityPlan, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var businessContinuityPlan *models.BusinessContinuityPlan
	result := db.
		Preload("RtoUnit").
		Preload("MtpdUnit").
		Preload("RpoUnit").
		Preload("Processes").
		Preload("Processes.Process").
		Preload("Processes.RtoUnit").
		Preload("Processes.MtpdUnit").
		Preload("Processes.RpoUnit").
		Preload("Processes.RiskLevel").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		First(&businessContinuityPlan, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return businessContinuityPlan, nil
}

func (repo *BcpBiaMatrixRepository) GetByID(id uint) (*models.BcpProcess, error) {
	var bcpProcess *models.BcpProcess

	result := repo.db.First(&bcpProcess, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return bcpProcess, nil

}

func (repo *BcpBiaMatrixRepository) Update(bcpProcess *models.BcpProcess, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(bcpProcess).Error; err != nil {
		return err
	}

	return nil
}
