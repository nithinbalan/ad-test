package repositories

import (
	"fmt"
	incidentDTOs "jethur/grc/incident/dtos"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type IncidentResponseTeamRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(params *incidentDTOs.IncidentResponseTeamListQueryParamsDTO) ([]*models.IncidentResponseTeam, int64, error)
	GetByID(id uint) (*models.IncidentResponseTeam, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.IncidentResponseTeam, error)
	Create(team *models.IncidentResponseTeam, tx ...utility.DatabaseTransaction) (uint, error)
	Update(team *models.IncidentResponseTeam, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error

	CreateMembers(members []*models.IncidentResponseTeamMember, tx ...utility.DatabaseTransaction) error
	DeleteMembers(teamMembersIDs []*uint, tx ...utility.DatabaseTransaction) error
}

type IncidentResponseTeamRepository struct {
	db *gorm.DB
}

func NewIncidentResponseTeamRepository(db *gorm.DB) *IncidentResponseTeamRepository {
	return &IncidentResponseTeamRepository{
		db: db,
	}
}

func (repo *IncidentResponseTeamRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *IncidentResponseTeamRepository) GetAll(params *incidentDTOs.IncidentResponseTeamListQueryParamsDTO) ([]*models.IncidentResponseTeam, int64, error) {
	var teams []*models.IncidentResponseTeam
	var totalCount int64

	query := repo.db.Preload("TeamLeadUser").
		Preload("IncidentResponseTeamMembers").
		Preload("IncidentResponseTeamMembers.User").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		Joins("LEFT JOIN users ON users.id = incident_response_teams.team_lead_user_id").
		Model(&models.IncidentResponseTeam{})

	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where(fmt.Sprintf("%s IN (?)", columnName), values)
		}
	}

	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
            (incident_response_teams.title_en LIKE ? OR
            incident_response_teams.title_ar LIKE ? OR
            incident_response_teams.scope LIKE ? OR
            incident_response_teams.objective LIKE ? OR
            users.first_name_en LIKE ? OR
            users.last_name_en LIKE ?)`,
			keyword, keyword, keyword, keyword, keyword, keyword)
	}

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	if params.SortBy == "" {
		params.SortBy = "created_at"
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "desc"
	}

	query = query.Order(params.SortBy + " " + params.Sort)

	if params.Page > 0 && params.PageSize > 0 {
		offset := (params.Page - 1) * params.PageSize
		query = query.Offset(offset).Limit(params.PageSize)
	}

	if err := query.Find(&teams).Error; err != nil {
		return nil, 0, err
	}

	return teams, totalCount, nil
}

func (repo *IncidentResponseTeamRepository) GetByID(id uint) (*models.IncidentResponseTeam, error) {
	var team *models.IncidentResponseTeam
	result := repo.db.First(&team, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return team, nil
}

func (repo *IncidentResponseTeamRepository) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.IncidentResponseTeam, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var team models.IncidentResponseTeam

	result := db.
		Preload("TeamLeadUser").
		Preload("IncidentResponseTeamMembers").
		Preload("IncidentResponseTeamMembers.User").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		First(&team, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return &team, nil
}

func (repo *IncidentResponseTeamRepository) Create(team *models.IncidentResponseTeam, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(team).Error; err != nil {
		return 0, err
	}
	return team.ID, nil
}

func (repo *IncidentResponseTeamRepository) Update(team *models.IncidentResponseTeam, tx ...utility.DatabaseTransaction) error {
	if err := repo.db.Save(team).Error; err != nil {
		return err
	}
	return nil
}

func (repo *IncidentResponseTeamRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var team *models.IncidentResponseTeam
	if err := db.First(&team, id).Error; err != nil {
		return err
	}

	result := db.Delete(&models.IncidentResponseTeam{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *IncidentResponseTeamRepository) CreateMembers(members []*models.IncidentResponseTeamMember, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(&members).Error; err != nil {
		return err
	}
	return nil
}

func (repo *IncidentResponseTeamRepository) DeleteMembers(teamMembersIDs []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("id IN (?)", teamMembersIDs).Delete(&models.IncidentResponseTeamMember{}).Error; err != nil {
		return err
	}

	return nil
}
