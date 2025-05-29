package repositories

import (
	"jethur/grc/models"
	organizationalContextDTOs "jethur/grc/organizational-context/dtos"
	"jethur/utility"

	"gorm.io/gorm"
)

type StakeholderNeedsAndExpectationRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(stakeholderId uint, params *organizationalContextDTOs.StakeholderNeedsAndExpectationListQueryParamsDTO) ([]*organizationalContextDTOs.StakeholderNeedsAndExpectationListDTO, error)
	GetByID(id uint) (*models.StakeholderNeedsAndExpectation, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.StakeholderNeedsAndExpectation, error)
	Create(stakeholderNeedsAndExpectation *models.StakeholderNeedsAndExpectation, tx ...utility.DatabaseTransaction) (uint, error)
	Update(stakeholderNeedsAndExpectation *models.StakeholderNeedsAndExpectation, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error
}

type StakeholderNeedsAndExpectationRepository struct {
	db *gorm.DB
}

func NewStakeholderNeedsAndExpectationRepository(db *gorm.DB) *StakeholderNeedsAndExpectationRepository {
	return &StakeholderNeedsAndExpectationRepository{
		db: db,
	}
}

// BeginTransaction starts a new transaction.
func (repo *StakeholderNeedsAndExpectationRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *StakeholderNeedsAndExpectationRepository) GetAll(stakeholderId uint, params *organizationalContextDTOs.StakeholderNeedsAndExpectationListQueryParamsDTO) ([]*organizationalContextDTOs.StakeholderNeedsAndExpectationListDTO, error) {
	var stakeholderNeedsAndExpectations []*organizationalContextDTOs.StakeholderNeedsAndExpectationListDTO
	var totalCount int64

	query := repo.db.Table("stakeholder_needs_and_expectations").
		Select(`
            stakeholder_needs_and_expectations.*,
			stakeholders.title_en as stakeholder_title_en,
			stakeholders.title_ar as stakeholder_title_ar,
			severity_levels.title_en AS severity_level_title_en,
			severity_levels.title_ar AS severity_level_title_ar,
			severity_levels.color_code AS severity_level_color_code,
			severity_levels.type AS severity_level_type,
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
			LEFT JOIN stakeholders ON stakeholders.id = stakeholder_needs_and_expectations.stakeholder_id
			LEFT JOIN severity_levels ON severity_levels.id = stakeholder_needs_and_expectations.severity_level_id
			LEFT JOIN users AS created_user ON created_user.id = stakeholder_needs_and_expectations.created_by_id
			LEFT JOIN users AS updated_user ON updated_user.id = stakeholder_needs_and_expectations.updated_by_id
        `).Where(`stakeholder_needs_and_expectations.stakeholder_id`, stakeholderId)
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

	err := query.Find(&stakeholderNeedsAndExpectations).Error
	if err != nil {
		return nil, err
	}

	return stakeholderNeedsAndExpectations, nil
}

func (repo *StakeholderNeedsAndExpectationRepository) GetByID(id uint) (*models.StakeholderNeedsAndExpectation, error) {
	var stakeholderNeedsAndExpectation *models.StakeholderNeedsAndExpectation

	result := repo.db.First(&stakeholderNeedsAndExpectation, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return stakeholderNeedsAndExpectation, nil

}

func (repo *StakeholderNeedsAndExpectationRepository) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.StakeholderNeedsAndExpectation, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var stakeholderNeedsAndExpectation *models.StakeholderNeedsAndExpectation

	result := db.
		Preload("SeverityLevel").
		Preload("Stakeholder").
		Preload("Stakeholder.Type").
		Preload("Stakeholder.StakeholderCategory").
		Preload("Stakeholder.InfluenceLevel").
		Preload("Stakeholder.ImpactLevel").
		Preload("Stakeholder.OrganizationUnit").
		Preload("Stakeholder.ThirdPartyRegister").
		Preload("Stakeholder.Client").
		Preload("Stakeholder.User").
		Preload("Evaluations").
		Preload("Evaluations.EvaluationStatus").
		Preload("Evaluations.EvaluatedByUser").
		Preload("Evaluations.CreatedByUser").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		First(&stakeholderNeedsAndExpectation, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return stakeholderNeedsAndExpectation, nil
}

func (repo *StakeholderNeedsAndExpectationRepository) Create(stakeholderNeedsAndExpectation *models.StakeholderNeedsAndExpectation, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(stakeholderNeedsAndExpectation).Error; err != nil {
		return 0, err
	}

	return stakeholderNeedsAndExpectation.ID, nil
}

func (repo *StakeholderNeedsAndExpectationRepository) Update(stakeholderNeedsAndExpectation *models.StakeholderNeedsAndExpectation, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(stakeholderNeedsAndExpectation).Error; err != nil {
		return err
	}

	return nil
}

func (repo *StakeholderNeedsAndExpectationRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	// Delete the Swot Pestel Question Factor
	if err := db.Delete(&models.StakeholderNeedsAndExpectation{}, id).Error; err != nil {
		return err
	}

	return nil
}
