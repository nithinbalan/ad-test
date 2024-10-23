package repositories

import (
	"jethur/grc/reports/dtos"

	"gorm.io/gorm"
)

type OrganizationReportRepositoryInterface interface {
	GetUsersByStatus(params *dtos.ReportListQueryParamsDTO) ([]*dtos.UsersByStatusDTO, error)
}

type OrganizationReportRepository struct {
	db *gorm.DB
}

func NewOrganizationReportRepository(db *gorm.DB) *OrganizationReportRepository {
	return &OrganizationReportRepository{db: db}
}

func (repo *OrganizationReportRepository) GetUsersByStatus(params *dtos.ReportListQueryParamsDTO) ([]*dtos.UsersByStatusDTO, error) {
	var results []*dtos.UsersByStatusDTO

	query := repo.db.Table("users").
		Select(`
			user_statuses.id AS id,
			user_statuses.title_en AS user_status_title_en,
			user_statuses.title_ar AS user_status_title_ar,
			user_statuses.color_code AS user_status_color_code,
			user_statuses.type AS user_status_type,
			COUNT(users.id) AS count`).
		Joins("JOIN user_statuses ON user_statuses.id = users.user_status_id").
		Group("user_statuses.id, user_statuses.title_en, user_statuses.title_ar, user_statuses.color_code, user_statuses.type")

	// Apply date filters if provided
	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("users.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.Scan(&results).Error
	if err != nil {
		return nil, err
	}

	return results, nil
}
