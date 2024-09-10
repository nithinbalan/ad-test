package repositories

import (
	"errors"
	incidentDTOs "jethur/grc/incident/dtos"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type IncidentTimelineRepositoryInterface interface {
	GetAll(incidentID uint, params *incidentDTOs.IncidentTimelineListQueryParamsDTO) (*incidentDTOs.IncidentTimelineListResponseDTO, error)
	GetAllTypes() ([]models.IncidentTimelineType, error)
	GetByID(incidentID uint, timelineTypeID uint) (*models.IncidentTimeline, error)
	GetByType(timelineType string) (*models.IncidentTimelineType, error)
	GetByStatusType(status string) (*models.TimelineStatus, error)
	Create(timelines []models.IncidentTimeline, tx ...utility.DatabaseTransaction) error
	Update(incidentTimeline *models.IncidentTimeline, tx ...utility.DatabaseTransaction) error
	Delete(incidentID uint, tx ...utility.DatabaseTransaction) error
}

type IncidentTimelineRepository struct {
	db *gorm.DB
}

func NewIncidentTimelineRepository(db *gorm.DB) *IncidentTimelineRepository {
	return &IncidentTimelineRepository{
		db: db,
	}
}

func (repo *IncidentTimelineRepository) GetAll(incidentID uint, params *incidentDTOs.IncidentTimelineListQueryParamsDTO) (*incidentDTOs.IncidentTimelineListResponseDTO, error) {
	var timelineList []*incidentDTOs.IncidentTimelineListDTO
	var totalCount int64

	query := repo.db.Table("incident_timelines").
		Select(`
			incident_timelines.id AS id,
			incident_timelines.incident_register_id AS incident_register_id,
			incident_timeline_types.title_en AS incident_timeline_type_title_en,
			incident_timeline_types.title_ar AS incident_timeline_type_title_ar,
			incident_timeline_types.color_code AS incident_timeline_type_color_code,
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
			performed_by.color AS performed_by_color,
			performed_by.profile_image_path AS performed_by_profile_image_path,
			incident_timelines.performed_at AS performed_at
		`).
		Joins(`
			LEFT JOIN incident_timeline_types ON incident_timeline_types.id = incident_timelines.incident_timeline_type_id
			LEFT JOIN timeline_statuses ON timeline_statuses.id = incident_timelines.timeline_status_id
			LEFT JOIN users AS performed_by ON performed_by.id = incident_timelines.performed_by_id
		`).
		Where("incident_timelines.incident_register_id = ?", incidentID)

	// Calculate total count for pagination
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	// Apply search keyword if exists
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
			(incident_timeline_types.title_en LIKE ? OR
			incident_timeline_types.title_ar LIKE ? OR
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
		params.SortBy = "incident_timelines.performed_at"
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

	return &incidentDTOs.IncidentTimelineListResponseDTO{
		Items: timelineList,
	}, nil
}

func (repo *IncidentTimelineRepository) GetAllTypes() ([]models.IncidentTimelineType, error) {
	var timelineTypes []models.IncidentTimelineType
	if err := repo.db.Find(&timelineTypes).Error; err != nil {
		return nil, err
	}
	return timelineTypes, nil
}

func (repo *IncidentTimelineRepository) GetByID(incidentID uint, timelineTypeID uint) (*models.IncidentTimeline, error) {
	var incidentTimeline models.IncidentTimeline

	result := repo.db.Where("incident_register_id = ? AND incident_timeline_type_id = ?", incidentID, timelineTypeID).First(&incidentTimeline)
	if result.Error != nil {
		return nil, result.Error
	}

	return &incidentTimeline, nil
}

func (repo *IncidentTimelineRepository) GetByType(timelineType string) (*models.IncidentTimelineType, error) {
	var incidentTimelineType *models.IncidentTimelineType

	result := repo.db.Where("type = ?", timelineType).First(&incidentTimelineType)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, "Incident type with type "+timelineType)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return incidentTimelineType, nil
}

func (repo *IncidentTimelineRepository) GetByStatusType(status string) (*models.TimelineStatus, error) {
	var timelineStatus *models.TimelineStatus

	result := repo.db.Where("type = ?", status).First(&timelineStatus)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, "Timeline Status with status "+status)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return timelineStatus, nil
}

func (repo *IncidentTimelineRepository) Create(timelines []models.IncidentTimeline, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Create(&timelines).Error
}

func (repo *IncidentTimelineRepository) Update(incidentTimeline *models.IncidentTimeline, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(incidentTimeline).Error; err != nil {
		return err
	}

	return nil
}

func (repo *IncidentTimelineRepository) Delete(incidentID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	// Delete the Incident timelines where incident_register_id = incidentID
	if err := db.Where("incident_register_id = ?", incidentID).Delete(&models.IncidentTimeline{}).Error; err != nil {
		return err
	}

	return nil
}
