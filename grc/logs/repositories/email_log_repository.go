package repositories

import (
	logDTOs "jethur/grc/logs/dtos"
	"jethur/grc/models"

	"gorm.io/gorm"
)

type EmailLogRepositoryInteface interface {
	GetAll(params *logDTOs.EmailLogListQueryParamsDTO) ([]*models.Email, int64, error)
	GetFullDetailsByID(id uint) (*models.Email, error)
	GetSummary(params *logDTOs.EmailLogListQueryParamsDTO) (*logDTOs.EmailLogSummaryDTO, error)
}

type EmailLogRepository struct {
	db *gorm.DB
}

func NewEmailLogRepository(db *gorm.DB) *EmailLogRepository {
	return &EmailLogRepository{
		db: db,
	}
}

func (repo *EmailLogRepository) GetAll(params *logDTOs.EmailLogListQueryParamsDTO) ([]*models.Email, int64, error) {
	var totalCount int64
	var emails []*models.Email

	query := repo.db.Model(&models.Email{}).
		Joins("LEFT JOIN email_statuses ON email_statuses.id = emails.email_status_id")

	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
			emails.subject LIKE ? OR 
			emails.error_message LIKE ? OR 
			EXISTS (
				SELECT 1 FROM email_statuses
				WHERE email_statuses.id = emails.email_status_id
				AND (email_statuses.title_en LIKE ? OR email_statuses.title_ar LIKE ?)
			) OR 
			EXISTS (
				SELECT 1 FROM email_recipients 
				WHERE email_recipients.email_id = emails.id 
				AND (email_recipients.email_address LIKE ? OR email_recipients.display_name LIKE ?)
			)
		`, keyword, keyword, keyword, keyword, keyword, keyword)
	}

	for columnName, values := range params.Filters {
		if len(values) == 0 {
			continue
		}

		if columnName == "email_status_id" {
			query = query.Where("emails.email_status_id IN (?)", values)
		}
	}

	if params.StartDate != "" {
		query = query.Where("emails.created_at >= ?", params.StartDate+" 00:00:00")
	}
	if params.EndDate != "" {
		query = query.Where("emails.created_at < ?", params.EndDate+" 23:59:59")
	}

	if len(params.QuickFilter) > 0 {
		for _, qf := range params.QuickFilter {
			switch qf {
			case "permanent_failed":
				query = query.Where("emails.is_permanent_failed = ?", true)
			case "retrying":
				query = query.Where("email_statuses.type = ? AND emails.attempt_count > 0 AND emails.is_permanent_failed = ?", "failed", false)
			}
		}
	}

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	if params.SortBy == "" {
		params.SortBy = "emails.created_at"
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "desc"
	}

	query = query.Order(params.SortBy + " " + params.Sort)

	offset := (params.Page - 1) * params.PageSize
	err := query.
		Preload("EmailStatus").
		Preload("Recipients").
		Offset(offset).
		Limit(params.PageSize).
		Find(&emails).Error

	if err != nil {
		return nil, 0, err
	}

	return emails, totalCount, nil
}

func (repo *EmailLogRepository) GetFullDetailsByID(id uint) (*models.Email, error) {
	var email *models.Email
	err := repo.db.
		Preload("Recipients.User").
		Preload("Attachments").
		Preload("Module").
		Preload("SubModule").
		Preload("EmailStatus").
		First(&email, id).Error
	if err != nil {
		return nil, err
	}
	return email, nil
}

func (repo *EmailLogRepository) GetSummary(params *logDTOs.EmailLogListQueryParamsDTO) (*logDTOs.EmailLogSummaryDTO, error) {
	var summary logDTOs.EmailLogSummaryDTO

	query := repo.db.Model(&models.Email{}).
		Joins("LEFT JOIN email_statuses ON email_statuses.id = emails.email_status_id")

	if params.StartDate != "" {
		query = query.Where("emails.created_at >= ?", params.StartDate+" 00:00:00")
	}

	if params.EndDate != "" {
		query = query.Where("emails.created_at < ?", params.EndDate+" 23:59:59")
	}

	err := query.Select(`
        COUNT(emails.id) as total_count,
        COALESCE(SUM(CASE WHEN email_statuses.type = 'pending' THEN 1 ELSE 0 END), 0) as pending_count,
        COALESCE(SUM(CASE WHEN email_statuses.type = 'processing' THEN 1 ELSE 0 END), 0) as processing_count,
        COALESCE(SUM(CASE WHEN email_statuses.type = 'sent' THEN 1 ELSE 0 END), 0) as sent_count,
        COALESCE(SUM(CASE WHEN email_statuses.type = 'failed' THEN 1 ELSE 0 END), 0) as failed_count
    `).Scan(&summary).Error

	if err != nil {
		return nil, err
	}

	return &summary, nil
}
