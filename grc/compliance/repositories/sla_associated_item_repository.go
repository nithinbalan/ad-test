package repositories

import (
	commonDTOs "jethur/grc/common/dtos"
	complianceDTOs "jethur/grc/compliance/dtos"
	"jethur/grc/models"

	"gorm.io/gorm"
)

type SLAAssociatedItemRepositoryInterface interface {
	GetAllProcesses(slaID uint, params *complianceDTOs.SLAAssociatedItemListQueryParamsDTO) (*complianceDTOs.SLAAssociatedItemProcessListResponseDTO, error)
	GetAllAssociatedCounts(slaID uint) (*complianceDTOs.SlaAssociatedItemsCountResponseDTO, error)
}

type SLAAssociatedItemRepository struct {
	db *gorm.DB
}

func NewSLAAssociatedItemRepository(db *gorm.DB) *SLAAssociatedItemRepository {
	return &SLAAssociatedItemRepository{
		db: db,
	}
}

func (repo *SLAAssociatedItemRepository) GetAllProcesses(slaID uint, params *complianceDTOs.SLAAssociatedItemListQueryParamsDTO) (*complianceDTOs.SLAAssociatedItemProcessListResponseDTO, error) {
	var processList []*complianceDTOs.SLAProcessListDTO
	var totalCount int64

	query := repo.db.Table("sla_processes").
		Select(`
			processes.id,
			processes.reference_number,
			processes.title_en,
			processes.title_ar,
			processes.description,
			process_groups.title_en AS process_group_title_en,
			process_groups.title_ar AS process_group_title_ar,
			organization_units.title_en AS owner_unit_title_en,
			organization_units.title_ar AS owner_unit_title_ar,
			owner.id AS owner_id,
			owner.first_name_en AS owner_first_name_en,
			owner.first_name_ar AS owner_first_name_ar,
			owner.last_name_en AS owner_last_name_en,
			owner.last_name_ar AS owner_last_name_ar,
			owner.initial AS control_owner_initial,
            owner.color AS control_owner_color,
            owner.email AS control_owner_email,
            owner.profile_image_path AS control_owner_profile_image_path,
			lifecycle_statuses.title_en AS lifecycle_status_title_en,
			lifecycle_statuses.title_ar AS lifecycle_status_title_ar,
			lifecycle_statuses.color_code AS lifecycle_status_color_code
		`).
		Joins(`
			LEFT JOIN processes ON processes.id = sla_processes.process_id
			LEFT JOIN process_groups ON process_groups.id = processes.process_group_id
			LEFT JOIN organization_units ON organization_units.id = processes.owner_unit_id
			LEFT JOIN users AS owner ON owner.id = processes.owner_id
			LEFT JOIN lifecycle_statuses ON lifecycle_statuses.id = processes.lifecycle_status_id
		`).Where("sla_processes.sla_register_id = ?", slaID)

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	// Apply search keyword if exists
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
			(processes.title_en LIKE ? OR
			processes.title_ar LIKE ? OR
			processes.reference_number LIKE ? OR
			process_groups.title_en LIKE ? OR
			process_groups.title_ar LIKE ? OR
			organization_units.title_en LIKE ? OR
			organization_units.title_ar LIKE ? OR
			owner.first_name_en LIKE ? OR
			owner.first_name_ar LIKE ? OR
			owner.last_name_en LIKE ? OR
			owner.last_name_ar LIKE ? OR
			lifecycle_statuses.title_en LIKE ? OR
			lifecycle_statuses.title_ar LIKE ? OR
			lifecycle_statuses.color_code LIKE ?)`,
			keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword,
			keyword, keyword, keyword, keyword, keyword)
	}

	if params.SortBy == "" {
		params.SortBy = "processes.title_en"
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "asc"
	}

	query = query.Order(params.SortBy + " " + params.Sort)

	offset := (params.Page - 1) * params.PageSize

	// Execute the query
	err := query.Offset(offset).Limit(params.PageSize).Find(&processList).Error
	if err != nil {
		return nil, err
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	response := &complianceDTOs.SLAAssociatedItemProcessListResponseDTO{
		Items: processList,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}
	return response, nil
}

func (repo *SLAAssociatedItemRepository) GetAllAssociatedCounts(slaID uint) (*complianceDTOs.SlaAssociatedItemsCountResponseDTO, error) {
	var risksCount, incidentsCount, questionsCount int64
	var err error

	// Query for questions count
	err = repo.db.Model(&models.SLACheckQuestion{}).
		Joins("JOIN sla_check_question_groups ON sla_check_questions.sla_check_question_group_id = sla_check_question_groups.id").
		Joins("JOIN sla_checks ON sla_check_question_groups.sla_check_id = sla_checks.id").
		Where("sla_checks.sla_register_id = ?", slaID).
		Count(&questionsCount).Error
	if err != nil {
		return nil, err
	}

	// Query for risks count
	err = repo.db.Model(&models.RiskSLA{}).
		Where("sla_register_id = ?", slaID).
		Count(&risksCount).Error
	if err != nil {
		return nil, err
	}

	// Query for incidents count
	err = repo.db.Model(&models.IncidentSLA{}).
		Where("sla_register_id = ?", slaID).
		Count(&incidentsCount).Error
	if err != nil {
		return nil, err
	}

	return &complianceDTOs.SlaAssociatedItemsCountResponseDTO{
		Questions: int(questionsCount),
		Risks:     int(risksCount),
		Incidents: int(incidentsCount),
	}, nil
}
