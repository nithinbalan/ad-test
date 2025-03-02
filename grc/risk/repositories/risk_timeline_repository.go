package repositories

import (
	"errors"
	"jethur/grc/models"
	riskDTOs "jethur/grc/risk/dtos"
	"jethur/utility"

	"gorm.io/gorm"
)

type RiskTimelineRepositoryInterface interface {
	GetAll(riskID uint, params *riskDTOs.RiskTimelineListQueryParamsDTO) (*riskDTOs.RiskTimelineListResponseDTO, error)
	GetAllTypes() ([]models.RiskTimelineType, error)
	GetByID(riskID uint, timelineTypeID uint, tx ...utility.DatabaseTransaction) (*models.RiskTimeline, error)
	GetByType(timelineType string) (*models.RiskTimelineType, error)
	GetByStatusType(status string) (*models.TimelineStatus, error)
	Create(timelines []models.RiskTimeline, tx ...utility.DatabaseTransaction) error
	Update(riskTimeline *models.RiskTimeline, tx ...utility.DatabaseTransaction) error
	Delete(riskID uint, tx ...utility.DatabaseTransaction) error
}

type RiskTimelineRepository struct {
	db *gorm.DB
}

func NewRiskTimelineRepository(db *gorm.DB) *RiskTimelineRepository {
	return &RiskTimelineRepository{
		db: db,
	}
}

func (repo *RiskTimelineRepository) GetAll(riskID uint, params *riskDTOs.RiskTimelineListQueryParamsDTO) (*riskDTOs.RiskTimelineListResponseDTO, error) {
	var timelineList []*riskDTOs.RiskTimelineListDTO
	var totalCount int64

	query := repo.db.Table("risk_timelines").
		Select(`
			risk_timelines.id AS id,
			risk_timelines.risk_register_id AS risk_register_id,
			risk_timeline_types.type AS risk_timeline_type,
			risk_timeline_types.title_en AS risk_timeline_type_title_en,
			risk_timeline_types.title_ar AS risk_timeline_type_title_ar,
			risk_timeline_types.color_code AS risk_timeline_type_color_code,
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
			risk_timelines.performed_at AS performed_at
		`).
		Joins(`
			LEFT JOIN risk_timeline_types ON risk_timeline_types.id = risk_timelines.risk_timeline_type_id
			LEFT JOIN timeline_statuses ON timeline_statuses.id = risk_timelines.timeline_status_id
			LEFT JOIN users AS performed_by ON performed_by.id = risk_timelines.performed_by_id
		`).
		Where("risk_timelines.risk_register_id = ?", riskID)

	// Calculate total count for pagination
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	// Apply search keyword if exists
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
			(risk_timeline_types.title_en LIKE ? OR
			risk_timeline_types.title_ar LIKE ? OR
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
		params.SortBy = "risk_timelines.performed_at"
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

	return &riskDTOs.RiskTimelineListResponseDTO{
		Items: timelineList,
	}, nil
}

func (repo *RiskTimelineRepository) GetAllTypes() ([]models.RiskTimelineType, error) {
	var timelineTypes []models.RiskTimelineType
	if err := repo.db.Find(&timelineTypes).Error; err != nil {
		return nil, err
	}
	return timelineTypes, nil
}

func (repo *RiskTimelineRepository) GetByID(riskID uint, timelineTypeID uint, tx ...utility.DatabaseTransaction) (*models.RiskTimeline, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var riskTimeline models.RiskTimeline

	result := db.Where("risk_register_id = ? AND risk_timeline_type_id = ?", riskID, timelineTypeID).First(&riskTimeline)
	if result.Error != nil {
		return nil, result.Error
	}

	return &riskTimeline, nil
}

func (repo *RiskTimelineRepository) GetByType(timelineType string) (*models.RiskTimelineType, error) {
	var riskTimelineType *models.RiskTimelineType

	result := repo.db.Where("type = ?", timelineType).First(&riskTimelineType)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, "Risk type with type "+timelineType)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return riskTimelineType, nil
}

func (repo *RiskTimelineRepository) GetByStatusType(status string) (*models.TimelineStatus, error) {
	var timelineStatus *models.TimelineStatus

	result := repo.db.Where("type = ?", status).First(&timelineStatus)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, "Timeline Status with status "+status)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return timelineStatus, nil
}

func (repo *RiskTimelineRepository) Create(timelines []models.RiskTimeline, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Create(&timelines).Error
}

func (repo *RiskTimelineRepository) Update(riskTimeline *models.RiskTimeline, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(riskTimeline).Error; err != nil {
		return err
	}

	return nil
}

func (repo *RiskTimelineRepository) Delete(riskID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	// Delete the Risk timelines where risk_register_id = riskID
	if err := db.Where("risk_register_id = ?", riskID).Delete(&models.RiskTimeline{}).Error; err != nil {
		return err
	}

	return nil
}
