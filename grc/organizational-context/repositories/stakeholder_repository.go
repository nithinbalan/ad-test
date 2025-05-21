package repositories

import (
	"jethur/grc/models"
	organizationalContextDTOs "jethur/grc/organizational-context/dtos"
	"jethur/utility"

	"gorm.io/gorm"
)

type StakeholderRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(params *organizationalContextDTOs.StakeholderListQueryParamsDTO) ([]*organizationalContextDTOs.StakeholderListDTO, error)
	GetByID(id uint) (*models.Stakeholder, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.Stakeholder, error)
	Create(stkeholder *models.Stakeholder, tx ...utility.DatabaseTransaction) (uint, error)
	Update(stkeholder *models.Stakeholder, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error
}

type StakeholderRepository struct {
	db *gorm.DB
}

func NewStakeholderRepository(db *gorm.DB) *StakeholderRepository {
	return &StakeholderRepository{
		db: db,
	}
}

// BeginTransaction starts a new transaction.
func (repo *StakeholderRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *StakeholderRepository) GetAll(params *organizationalContextDTOs.StakeholderListQueryParamsDTO) ([]*organizationalContextDTOs.StakeholderListDTO, error) {
	var stkeholders []*organizationalContextDTOs.StakeholderListDTO
	var totalCount int64

	query := repo.db.Table("stakeholders").
		Select(`
            stakeholders.*,
			process_dependency_types.title_en AS type_title_en,
			process_dependency_types.title_ar AS type_title_ar,
			influence_level.title_en AS influence_level_title_en,
			influence_level.title_ar AS influence_level_title_ar,
			influence_level.type AS influence_level_type,
			impact_level.title_en AS impact_level_title_en,
			impact_level.title_ar AS impact_level_title_ar,
			impact_level.type AS impact_level_type,
			stakeholder_categories.title_en AS stakeholder_category_title_en,
			stakeholder_categories.title_ar AS stakeholder_category_title_ar,
			stakeholder_categories.type AS stakeholder_category_type,
			organization_units.title_en AS organization_unit_title_en,
			organization_units.title_ar AS organization_unit_title_ar,
			third_party_register.title AS third_party_register_title,
			client.title AS client_title,
			individual_user.first_name_en AS user_first_name_en,
            individual_user.first_name_ar AS user_first_name_ar,
            individual_user.last_name_en AS user_last_name_en,
            individual_user.last_name_ar AS user_last_name_ar,
			individual_user.initial AS user_initial,
            individual_user.color AS user_color,
            individual_user.email AS user_email,
            individual_user.profile_image_path AS user_profile_image_path,   
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
			LEFT JOIN process_dependency_types ON process_dependency_types.id = stakeholders.type_id
			LEFT JOIN stakeholder_categories ON stakeholder_categories.id = stakeholders.stakeholder_category_id
			LEFT JOIN severity_levels AS influence_level ON influence_level.id = stakeholders.influence_level_id
			LEFT JOIN severity_levels AS impact_level ON impact_level.id = stakeholders.impact_level_id
			LEFT JOIN organization_units ON organization_units.id = stakeholders.organization_unit_id
			LEFT JOIN third_party_register AS third_party_register ON third_party_register.id = stakeholders.third_party_register_id
			LEFT JOIN third_party_register AS client ON client.id = stakeholders.client_id
			LEFT JOIN users AS individual_user ON individual_user.id = stakeholders.user_id
			LEFT JOIN users AS created_user ON created_user.id = stakeholders.created_by_id
			LEFT JOIN users AS updated_user ON updated_user.id = stakeholders.updated_by_id
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
			(created_user.first_name_en LIKE ? OR
			created_user.first_name_ar LIKE ? OR
			created_user.last_name_en LIKE ? OR
			created_user.last_name_ar LIKE ? OR
    		updated_user.first_name_en LIKE ? OR
			updated_user.first_name_ar LIKE ? OR
			updated_user.last_name_en LIKE ? OR
			updated_user.last_name_ar LIKE ?)`,
			keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword)
	}

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	query = query.Order(params.SortBy + " " + params.Sort)

	if !params.Export {
		offset := (params.Page - 1) * params.PageSize
		query = query.Offset(offset).Limit(params.PageSize)
	}

	err := query.Find(&stkeholders).Error
	if err != nil {
		return nil, err
	}
	return stkeholders, nil
}

func (repo *StakeholderRepository) GetByID(id uint) (*models.Stakeholder, error) {
	var stkeholder *models.Stakeholder

	result := repo.db.First(&stkeholder, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return stkeholder, nil

}

func (repo *StakeholderRepository) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.Stakeholder, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var stkeholder *models.Stakeholder

	result := db.
		Preload("Type").
		Preload("StakeholderCategory").
		Preload("InfluenceLevel").
		Preload("ImpactLevel").
		Preload("OrganizationUnit").
		Preload("ThirdPartyRegister").
		Preload("Client").
		Preload("User").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		First(&stkeholder, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return stkeholder, nil
}

func (repo *StakeholderRepository) Create(stkeholder *models.Stakeholder, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(stkeholder).Error; err != nil {
		return 0, err
	}

	return stkeholder.ID, nil
}

func (repo *StakeholderRepository) Update(stkeholder *models.Stakeholder, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(stkeholder).Error; err != nil {
		return err
	}

	return nil
}

func (repo *StakeholderRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	// Delete the Swot Pestel Question Factor
	if err := db.Delete(&models.Stakeholder{}, id).Error; err != nil {
		return err
	}

	return nil
}
