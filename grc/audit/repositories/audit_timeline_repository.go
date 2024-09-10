package repositories

import (
	"errors"
	auditDTOs "jethur/grc/audit/dtos"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type AuditTimelineRepositoryInterface interface {
	GetAll(auditID uint, params *auditDTOs.AuditTimelineListQueryParamsDTO) (*auditDTOs.AuditTimelineListResponseDTO, error)
	GetAllTypes() ([]models.AuditTimelineType, error)
	GetByID(auditID uint, timelineTypeID uint) (*models.AuditTimeline, error)
	GetByType(timelineType string) (*models.AuditTimelineType, error)
	GetByStatusType(status string) (*models.TimelineStatus, error)
	Create(timelines []models.AuditTimeline, tx ...utility.DatabaseTransaction) error
	Update(auditTimeline *models.AuditTimeline, tx ...utility.DatabaseTransaction) error
	Delete(auditID uint, tx ...utility.DatabaseTransaction) error
}

type AuditTimelineRepository struct {
	db *gorm.DB
}

func NewAuditTimelineRepository(db *gorm.DB) *AuditTimelineRepository {
	return &AuditTimelineRepository{
		db: db,
	}
}

func (repo *AuditTimelineRepository) GetAll(auditID uint, params *auditDTOs.AuditTimelineListQueryParamsDTO) (*auditDTOs.AuditTimelineListResponseDTO, error) {
	var timelineList []*auditDTOs.AuditTimelineListDTO
	var totalCount int64

	query := repo.db.Table("audit_timelines").
		Select(`
			audit_timelines.id AS id,
			audit_timelines.audit_id AS audit_id,
			audit_timeline_types.title_en AS audit_timeline_type_title_en,
			audit_timeline_types.title_ar AS audit_timeline_type_title_ar,
			audit_timeline_types.color_code AS audit_timeline_type_color_code,
			timeline_statuses.title_en AS timeline_status_title_en,
			timeline_statuses.title_ar AS timeline_status_title_ar,
			timeline_statuses.type AS timeline_status_type,
			timeline_statuses.color_code AS timeline_status_color_code,
			performed_by.id AS performed_by_id,
			performed_by.first_name_en AS performed_by_first_name_en,
			performed_by.first_name_ar AS performed_by_first_name_ar,
			performed_by.last_name_en AS performed_by_last_name_en,
			performed_by.last_name_ar AS performed_by_last_name_ar,
			performed_by.email AS performed_by_email,
			performed_by.initial AS performed_by_initial,
			performed_by.color AS performed_by_color_code,
			performed_by.profile_image_path AS performed_by_profile_image_path,
			audit_timelines.performed_at AS performed_at
		`).
		Joins(`
			LEFT JOIN audit_timeline_types ON audit_timeline_types.id = audit_timelines.audit_timeline_type_id
			LEFT JOIN timeline_statuses ON timeline_statuses.id = audit_timelines.timeline_status_id
			LEFT JOIN users AS performed_by ON performed_by.id = audit_timelines.performed_by_id
		`).
		Where("audit_timelines.audit_id = ?", auditID)

	// Calculate total count for pagination
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	// Apply search keyword if exists
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
			(audit_timeline_types.title_en LIKE ? OR
			audit_timeline_types.title_ar LIKE ? OR
			timeline_statuses.title_en LIKE ? OR
			timeline_statuses.title_ar LIKE ? OR
			performed_by.first_name_en LIKE ? OR
			performed_by.first_name_ar LIKE ? OR
			performed_by.last_name_en LIKE ? OR
			performed_by.last_name_ar LIKE ? OR
			performed_by.email LIKE ?)`,
			keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword)
	}

	if params.SortBy == "" {
		params.SortBy = "audit_timelines.performed_at"
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "asc"
	}

	query = query.Order(params.SortBy + " " + params.Sort)

	// Execute the query
	err := query.Find(&timelineList).Error
	if err != nil {
		return nil, err
	}

	return &auditDTOs.AuditTimelineListResponseDTO{
		Items: timelineList,
	}, nil
}

func (repo *AuditTimelineRepository) GetAllTypes() ([]models.AuditTimelineType, error) {
	var timelineTypes []models.AuditTimelineType
	if err := repo.db.Find(&timelineTypes).Error; err != nil {
		return nil, err
	}
	return timelineTypes, nil
}

func (repo *AuditTimelineRepository) GetByID(auditID uint, timelineTypeID uint) (*models.AuditTimeline, error) {
	var auditTimeline models.AuditTimeline

	result := repo.db.Where("audit_id = ? AND audit_timeline_type_id = ?", auditID, timelineTypeID).First(&auditTimeline)
	if result.Error != nil {
		return nil, result.Error
	}

	return &auditTimeline, nil
}

func (repo *AuditTimelineRepository) GetByType(timelineType string) (*models.AuditTimelineType, error) {
	var auditTimelineType *models.AuditTimelineType

	result := repo.db.Where("type = ?", timelineType).First(&auditTimelineType)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, "Audit type with type "+timelineType)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return auditTimelineType, nil
}

func (repo *AuditTimelineRepository) GetByStatusType(status string) (*models.TimelineStatus, error) {
	var timelineStatus *models.TimelineStatus

	result := repo.db.Where("type = ?", status).First(&timelineStatus)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, "Timeline Status with status "+status)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return timelineStatus, nil
}

func (repo *AuditTimelineRepository) Create(timelines []models.AuditTimeline, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Create(&timelines).Error
}

func (repo *AuditTimelineRepository) Update(auditTimeline *models.AuditTimeline, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(auditTimeline).Error; err != nil {
		return err
	}

	return nil
}

func (repo *AuditTimelineRepository) Delete(auditID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	// Delete the Audit timelines where audit_id = auditID
	if err := db.Where("audit_id = ?", auditID).Delete(&models.AuditTimeline{}).Error; err != nil {
		return err
	}

	return nil
}
