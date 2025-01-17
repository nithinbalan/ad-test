package repositories

import (
	businessContinuityDTOs "jethur/grc/business-continuity/dtos"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type BcpRolesRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(params *businessContinuityDTOs.BcpRolesListQueryParamsDTO) (*businessContinuityDTOs.BcpRolesListResponseDTO, error)
	GetByID(id uint) (*models.BcpRole, error)
	Create(bcpRole *models.BcpRole, tx ...utility.DatabaseTransaction) (uint, error)
	Update(bcpRole *models.BcpRole, tx ...utility.DatabaseTransaction) error
	Delete(roleID uint, tx ...utility.DatabaseTransaction) error
}

type BcpRolesRepository struct {
	db *gorm.DB
}

func NewBcpRolesRepository(db *gorm.DB) *BcpRolesRepository {
	return &BcpRolesRepository{
		db: db,
	}
}

func (repo *BcpRolesRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *BcpRolesRepository) GetAll(params *businessContinuityDTOs.BcpRolesListQueryParamsDTO) (*businessContinuityDTOs.BcpRolesListResponseDTO, error) {
	var roles []*businessContinuityDTOs.BcpRolesListDTO

	query := repo.db.Table("bcp_roles").
		Select(`
            bcp_roles.*,
            responsible_user.first_name_en AS responsible_user_first_name_en,
            responsible_user.first_name_ar AS responsible_user_first_name_ar,
            responsible_user.last_name_en AS responsible_user_last_name_en,
            responsible_user.last_name_ar AS responsible_user_last_name_ar,
            responsible_user.email AS responsible_user_email,
            backup_user.first_name_en AS backup_user_first_name_en,
            backup_user.first_name_ar AS backup_user_first_name_ar,
            backup_user.last_name_en AS backup_user_last_name_en,
            backup_user.last_name_ar AS backup_user_last_name_ar,
            backup_user.email AS backup_user_email,
            created_user.first_name_en AS created_user_first_name_en,
            created_user.first_name_ar AS created_user_first_name_ar,
            created_user.last_name_en AS created_user_last_name_en,
            created_user.last_name_ar AS created_user_last_name_ar,
            created_user.email AS created_user_email,
            organization_units.title_en AS responsible_unit_title_en
        `).
		Joins(`LEFT JOIN users AS responsible_user ON responsible_user.id = bcp_roles.responsible_user_id`).
		Joins(`LEFT JOIN users AS backup_user ON backup_user.id = bcp_roles.backup_user_id`).
		Joins(`LEFT JOIN users AS created_user ON created_user.id = bcp_roles.created_by_id`).
		Joins(`LEFT JOIN organization_units ON organization_units.id = bcp_roles.responsible_unit_id`)

	// Apply filters dynamically, including `bcpId`
	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where(columnName+" IN (?)", values)
		}
	}

	// Apply keyword search
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
			(bcp_roles.responsibility LIKE ? OR
			responsible_user.first_name_en LIKE ? OR
			responsible_user.first_name_ar LIKE ? OR
			responsible_user.last_name_en LIKE ? OR
			responsible_user.last_name_ar LIKE ? OR
			responsible_user.email LIKE ? OR
			backup_user.first_name_en LIKE ? OR
			backup_user.first_name_ar LIKE ? OR
			backup_user.last_name_en LIKE ? OR
			backup_user.last_name_ar LIKE ? OR
			backup_user.email LIKE ? OR
			organization_units.title_en LIKE ? OR
			created_user.first_name_en LIKE ? OR
			created_user.first_name_ar LIKE ? OR
			created_user.last_name_en LIKE ? OR
			created_user.last_name_ar LIKE ? OR
			created_user.email LIKE ?)`,
			keyword, keyword, keyword, keyword, keyword, keyword,
			keyword, keyword, keyword, keyword, keyword,
			keyword, keyword, keyword, keyword, keyword, keyword)
	}

	// Fetch all results
	if err := query.Find(&roles).Error; err != nil {
		return nil, err
	}

	// Construct response
	response := &businessContinuityDTOs.BcpRolesListResponseDTO{
		Items: roles,
	}

	return response, nil
}

// GetByID fetches a BCP role by its ID.
func (repo *BcpRolesRepository) GetByID(id uint) (*models.BcpRole, error) {
	var role models.BcpRole

	if err := repo.db.First(&role, id).Error; err != nil {
		return nil, err
	}

	return &role, nil
}

// Create inserts a new BCP role.
func (repo *BcpRolesRepository) Create(bcpRole *models.BcpRole, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Create(bcpRole).Error; err != nil {
		return 0, err
	}

	return bcpRole.ID, nil
}

// Update modifies an existing BCP role.
func (repo *BcpRolesRepository) Update(bcpRole *models.BcpRole, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Save(bcpRole).Error; err != nil {
		return err
	}

	return nil
}

func (repo *BcpRolesRepository) Delete(roleID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("id = ?", roleID).Delete(&models.BcpRole{}).Error; err != nil {
		return err
	}

	return nil
}
