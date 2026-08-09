package repositories

import (
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/identity-access-review/dtos"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type IdentityReviewUserRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(params *dtos.IdentityReviewUserListQueryParamsDTO) (*dtos.IdentityReviewUserListResponseDTO, error)
	GetByID(id uint) (*models.IdentityReviewUser, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.IdentityReviewUser, error)
	Create(identityReviewUser *models.IdentityReviewUser, tx ...utility.DatabaseTransaction) (uint, error)
	Update(identityReviewUser *models.IdentityReviewUser, tx ...utility.DatabaseTransaction) error
	UpdateAccountDecision(identityReviewUser *models.IdentityReviewUser, tx ...utility.DatabaseTransaction) error
	UpdatePrivilege(identityReviewUser *models.IdentityReviewUser, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error
}

func (repo *IdentityReviewUserRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

type IdentityReviewUserRepository struct {
	db *gorm.DB
}

func NewIdentityReviewUserRepository(db *gorm.DB) *IdentityReviewUserRepository {
	return &IdentityReviewUserRepository{db: db}
}

func (repo *IdentityReviewUserRepository) GetAll(params *dtos.IdentityReviewUserListQueryParamsDTO) (*dtos.IdentityReviewUserListResponseDTO, error) {
	var identityReviewUsers []*dtos.IdentityReviewUserListDTO
	var totalCount int64

	query := repo.db.Model(&models.IdentityReviewUser{}).
		Select(`
			identity_review_users.id,
			identity_review_users.identity_review_campaign_id,
			identity_review_campaigns.campaign_name AS identity_review_campaign_name,
			identity_review_users.username,
			identity_review_users.user_identifier,
			identity_review_users.last_login,
			identity_review_users.current_role,
			identity_review_users.department_id,
			departments.title_en AS department_title_en,
			departments.title_ar AS department_title_ar,
			identity_review_users.privilege_id,
			privileged_accounts.title_en AS privilege_title_en,
			privileged_accounts.title_ar AS privilege_title_ar,
			identity_review_users.account_type_id,
			account_types.title_en AS account_type_title_en,
			account_types.title_ar AS account_type_title_ar,
			identity_review_users.account_decision_id,
			account_decisions.title_en AS account_decision_title_en,
			account_decisions.title_ar AS account_decision_title_ar,
			account_decisions.type AS account_decision_type,
			identity_review_users.created_at,
			identity_review_users.updated_at,
			identity_review_users.created_by_id,
			created_by.first_name_en AS created_by_first_name_en,
			created_by.last_name_en AS created_by_last_name_en,
			created_by.initial AS created_by_initial,
			identity_review_users.updated_by_id,
			updated_by.first_name_en AS updated_by_first_name_en,
			updated_by.last_name_en AS updated_by_last_name_en,
			updated_by.initial AS updated_by_initial
		`).
		Joins("INNER JOIN identity_review_campaigns ON identity_review_campaigns.id = identity_review_users.identity_review_campaign_id").
		Joins("INNER JOIN organization_units AS departments ON departments.id = identity_review_users.department_id").
		Joins("LEFT JOIN privileged_accounts ON privileged_accounts.id = identity_review_users.privilege_id").
		Joins("LEFT JOIN account_types ON account_types.id = identity_review_users.account_type_id").
		Joins("LEFT JOIN account_decisions ON account_decisions.id = identity_review_users.account_decision_id").
		Joins("LEFT JOIN users AS created_by ON created_by.id = identity_review_users.created_by_id").
		Joins("LEFT JOIN users AS updated_by ON updated_by.id = identity_review_users.updated_by_id")

	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
			identity_review_users.username LIKE ? OR
			identity_review_users.user_identifier LIKE ? OR
			identity_review_users.last_login LIKE ? OR
			identity_review_users.current_role LIKE ? OR
			departments.title_en LIKE ? OR
			departments.title_ar LIKE ? OR
			privileged_accounts.title_en LIKE ? OR
			privileged_accounts.title_ar LIKE ? OR
			account_types.title_en LIKE ? OR
			account_types.title_ar LIKE ? OR
			account_decisions.title_en LIKE ? OR
			account_decisions.title_ar LIKE ?`,
			keyword, keyword, keyword, keyword,
			keyword, keyword, keyword, keyword,
			keyword, keyword, keyword, keyword,
		)
	}

	for column, values := range params.Filters {
		if len(values) == 0 {
			continue
		}

		switch column {
		case "identity_review_campaign_id":
			query = query.Where("identity_review_users.identity_review_campaign_id IN (?)", values)
		case "department_id":
			query = query.Where("identity_review_users.department_id IN (?)", values)
		case "privilege_id":
			query = query.Where("identity_review_users.privilege_id IN (?)", values)
		case "account_type_id":
			query = query.Where("identity_review_users.account_type_id IN (?)", values)
		case "account_decision_id":
			query = query.Where("identity_review_users.account_decision_id IN (?)", values)
		}
	}

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	offset := (params.Page - 1) * params.PageSize
	if err := query.
		Order(params.SortBy + " " + params.Sort).
		Offset(offset).
		Limit(params.PageSize).
		Find(&identityReviewUsers).Error; err != nil {
		return nil, err
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	return &dtos.IdentityReviewUserListResponseDTO{
		Items: identityReviewUsers,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}, nil
}

func (repo *IdentityReviewUserRepository) GetByID(id uint) (*models.IdentityReviewUser, error) {
	var identityReviewUser models.IdentityReviewUser

	if err := repo.db.First(&identityReviewUser, id).Error; err != nil {
		return nil, err
	}

	return &identityReviewUser, nil
}

func (repo *IdentityReviewUserRepository) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.IdentityReviewUser, error) {
	var identityReviewUser models.IdentityReviewUser
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.
		Preload("IdentityReviewCampaign").
		Preload("Department").
		Preload("Privilege").
		Preload("AccountType").
		Preload("AccountDecision").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		First(&identityReviewUser, id).Error; err != nil {
		return nil, err
	}

	return &identityReviewUser, nil
}

func (repo *IdentityReviewUserRepository) Create(identityReviewUser *models.IdentityReviewUser, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.
		Omit("updated_at", "updated_by_id").
		Create(identityReviewUser).Error; err != nil {
		return 0, err
	}

	return identityReviewUser.ID, nil
}

func (repo *IdentityReviewUserRepository) Update(identityReviewUser *models.IdentityReviewUser, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Save(identityReviewUser).Error
}

func (repo *IdentityReviewUserRepository) UpdateAccountDecision(identityReviewUser *models.IdentityReviewUser, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.
		Model(&models.IdentityReviewUser{}).
		Where("id = ?", identityReviewUser.ID).
		Updates(map[string]interface{}{
			"account_decision_id": identityReviewUser.AccountDecisionID,
			"updated_by_id":       identityReviewUser.UpdatedByID,
			"updated_at":          identityReviewUser.UpdatedAt,
		}).Error
}

func (repo *IdentityReviewUserRepository) UpdatePrivilege(identityReviewUser *models.IdentityReviewUser, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.
		Model(&models.IdentityReviewUser{}).
		Where("id = ?", identityReviewUser.ID).
		Updates(map[string]interface{}{
			"privilege_id":  identityReviewUser.PrivilegeID,
			"updated_by_id": identityReviewUser.UpdatedByID,
			"updated_at":    identityReviewUser.UpdatedAt,
		}).Error
}

func (repo *IdentityReviewUserRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Delete(&models.IdentityReviewUser{}, id).Error
}
