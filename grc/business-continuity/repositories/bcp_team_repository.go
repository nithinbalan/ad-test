package repositories

import (
	businessContinuityDTOs "jethur/grc/business-continuity/dtos"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type BcpTeamsRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(params *businessContinuityDTOs.BcpTeamsListQueryParamsDTO, bcpId uint) ([]*businessContinuityDTOs.BcpTeamsListDTO, error)
	GetByID(id uint) (*models.BcpTeam, error)
	Create(bcpTeam *models.BcpTeam, tx ...utility.DatabaseTransaction) (uint, error)
	Update(bcpTeam *models.BcpTeam, tx ...utility.DatabaseTransaction) error
	Delete(TeamID uint, tx ...utility.DatabaseTransaction) error
}

type BcpTeamsRepository struct {
	db *gorm.DB
}

func NewBcpTeamsRepository(db *gorm.DB) *BcpTeamsRepository {
	return &BcpTeamsRepository{
		db: db,
	}
}

func (repo *BcpTeamsRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *BcpTeamsRepository) GetAll(params *businessContinuityDTOs.BcpTeamsListQueryParamsDTO, bcpId uint) ([]*businessContinuityDTOs.BcpTeamsListDTO, error) {
	var teams []*businessContinuityDTOs.BcpTeamsListDTO

	query := repo.db.Table("bcp_teams").
		Select(`
            bcp_teams.*,
            responsible_user.first_name_en AS responsible_user_first_name_en,
            responsible_user.first_name_ar AS responsible_user_first_name_ar,
            responsible_user.last_name_en AS responsible_user_last_name_en,
            responsible_user.last_name_ar AS responsible_user_last_name_ar,
            responsible_user.email AS responsible_user_email,
            responsible_user.profile_image_path AS responsible_user_profile_image_path,
            backup_user.first_name_en AS backup_user_first_name_en,
            backup_user.first_name_ar AS backup_user_first_name_ar,
            backup_user.last_name_en AS backup_user_last_name_en,
            backup_user.last_name_ar AS backup_user_last_name_ar,
            backup_user.email AS backup_user_email,
            backup_user.profile_image_path AS backup_user_profile_image_path,
            created_user.first_name_en AS created_user_first_name_en,
            created_user.first_name_ar AS created_user_first_name_ar,
            created_user.last_name_en AS created_user_last_name_en,
            created_user.last_name_ar AS created_user_last_name_ar,
            created_user.email AS created_user_email,
            organization_units.title_en AS responsible_unit_title_en,
            organization_units.title_ar AS responsible_unit_title_ar
        `).
		Joins(`LEFT JOIN users AS responsible_user ON responsible_user.id = bcp_teams.responsible_user_id`).
		Joins(`LEFT JOIN users AS backup_user ON backup_user.id = bcp_teams.backup_user_id`).
		Joins(`LEFT JOIN users AS created_user ON created_user.id = bcp_teams.created_by_id`).
		Joins(`LEFT JOIN organization_units ON organization_units.id = bcp_teams.responsible_unit_id`).
		Where("bcp_teams.business_continuity_plan_id = ?", bcpId)

	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where(columnName+" IN (?)", values)
		}
	}

	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
			(bcp_teams.responsibility LIKE ? OR
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

	if params.SortBy != "" && params.Sort != "" {
		query = query.Order(params.SortBy + " " + params.Sort)
	}

	if err := query.Find(&teams).Error; err != nil {
		return nil, err
	}

	return teams, nil
}

func (repo *BcpTeamsRepository) GetByID(id uint) (*models.BcpTeam, error) {
	var Team models.BcpTeam

	if err := repo.db.First(&Team, id).Error; err != nil {
		return nil, err
	}

	return &Team, nil
}

func (repo *BcpTeamsRepository) Create(bcpTeam *models.BcpTeam, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Create(bcpTeam).Error; err != nil {
		return 0, err
	}

	return bcpTeam.ID, nil
}

func (repo *BcpTeamsRepository) Update(bcpTeam *models.BcpTeam, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Save(bcpTeam).Error; err != nil {
		return err
	}

	return nil
}

func (repo *BcpTeamsRepository) Delete(TeamID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("id = ?", TeamID).Delete(&models.BcpTeam{}).Error; err != nil {
		return err
	}

	return nil
}
