package repositories

import (
	meetingDTOs "jethur/grc/meeting/dtos"

	"gorm.io/gorm"
)

type MeetingAssociatedItemRepositoryInterface interface {
	GetAllControls(meetingID uint, params *meetingDTOs.MeetingAssociatedItemListQueryParamsDTO) ([]*meetingDTOs.MeetingControlListDTO, error)
	GetAllComplianceRequirements(meetingID uint, params *meetingDTOs.MeetingAssociatedItemListQueryParamsDTO) ([]*meetingDTOs.MeetingComplianceRequirementListDTO, error)
	GetAllRisks(meetingID uint, params *meetingDTOs.MeetingAssociatedItemListQueryParamsDTO) ([]*meetingDTOs.MeetingRiskListDTO, error)
	GetAllProcesses(meetingID uint, params *meetingDTOs.MeetingAssociatedItemListQueryParamsDTO) ([]*meetingDTOs.MeetingProcessListDTO, error)
}

type MeetingAssociatedItemRepository struct {
	db *gorm.DB
}

func NewMeetingAssociatedItemRepository(db *gorm.DB) *MeetingAssociatedItemRepository {
	return &MeetingAssociatedItemRepository{
		db: db,
	}
}

func (repo *MeetingAssociatedItemRepository) GetAllControls(meetingID uint, params *meetingDTOs.MeetingAssociatedItemListQueryParamsDTO) ([]*meetingDTOs.MeetingControlListDTO, error) {
	var controlList []*meetingDTOs.MeetingControlListDTO

	query := repo.db.Table("meeting_controls").
		Select(`
			controls.id AS control_id,
            control_categories.title_en AS control_category_title_en,
            control_categories.title_ar AS control_category_title_ar,
            control_types.title_en AS control_type_title_en,
            control_types.title_ar AS control_type_title_ar,
            controls.title_en AS control_title_en,
            controls.title_ar AS control_title_ar,
            controls.description AS control_description,
            control_efficiency_statuses.title_en AS control_efficiency_status_en,
            control_efficiency_statuses.title_ar AS control_efficiency_status_ar,
            control_implementation_statuses.title_en AS control_implementation_status_en,
            control_implementation_statuses.title_ar AS control_implementation_status_ar
        `).
		Joins(`
			LEFT JOIN controls ON controls.id = meeting_controls.control_id
            LEFT JOIN control_categories ON controls.control_category_id = control_categories.id
            LEFT JOIN control_types ON controls.control_type_id = control_types.id
            LEFT JOIN control_efficiency_statuses ON controls.control_efficiency_status_id = control_efficiency_statuses.id
            LEFT JOIN control_implementation_statuses ON controls.control_implementation_status_id = control_implementation_statuses.id
        `).
		Where("meeting_controls.meeting_id = ?", meetingID)

	// Apply search keyword if exists
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
			(controls.title_en LIKE ? OR
			controls.title_ar LIKE ? OR
			controls.reference_number LIKE ? OR
			control_categories.title_en LIKE ? OR
			control_categories.title_ar LIKE ? OR
			control_types.title_en LIKE ? OR
			control_types.title_ar LIKE ? OR
			control_efficiency_statuses.title_en LIKE ? OR
			control_efficiency_statuses.title_ar LIKE ? OR
			control_implementation_statuses.title_en LIKE ? OR
			control_implementation_statuses.title_ar LIKE ?)`,
			keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword,
			keyword, keyword)
	}

	if params.SortBy == "" {
		params.SortBy = "controls.title_en"
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "asc"
	}

	query = query.Order(params.SortBy + " " + params.Sort)

	if err := query.Find(&controlList).Error; err != nil {
		return nil, err
	}

	return controlList, nil
}

func (repo *MeetingAssociatedItemRepository) GetAllComplianceRequirements(meetingID uint, params *meetingDTOs.MeetingAssociatedItemListQueryParamsDTO) ([]*meetingDTOs.MeetingComplianceRequirementListDTO, error) {
	var complianceRequirements []*meetingDTOs.MeetingComplianceRequirementListDTO

	query := repo.db.Table("meeting_compliance_requirements").
		Select(`
			meeting_compliance_requirements.compliance_requirement_id AS compliance_requirement_id,
            compliance_clauses.title AS compliance_clause_title,
            compliance_requirement_categories.title_en AS compliance_requirement_category_title_en,
            compliance_requirement_categories.title_ar AS compliance_requirement_category_title_ar,
            compliance_requirements.title AS compliance_requirement_title,
            compliance_statuses.title_en AS compliance_status_title_en,
            compliance_statuses.title_ar AS compliance_status_title_ar,
			compliance_statuses.color_code AS compliance_status_color_code

        `).
		Joins(`
			LEFT JOIN compliance_requirements ON compliance_requirements.id = meeting_compliance_requirements.compliance_requirement_id
            LEFT JOIN compliance_clauses ON compliance_requirements.compliance_clause_id = compliance_clauses.id
            LEFT JOIN compliance_statuses ON compliance_requirements.compliance_status_id = compliance_statuses.id
            LEFT JOIN compliance_requirement_categories ON compliance_requirements.compliance_requirement_category_id = compliance_requirement_categories.id
        `).
		Where("meeting_compliance_requirements.meeting_id = ?", meetingID)

	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
			(compliance_requirements.title LIKE ?)`,
			keyword)
	}

	if params.SortBy == "" {
		params.SortBy = "compliance_requirements.title"
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "asc"
	}

	query = query.Order(params.SortBy + " " + params.Sort)

	if err := query.Find(&complianceRequirements).Error; err != nil {
		return nil, err
	}

	return complianceRequirements, nil
}

func (repo *MeetingAssociatedItemRepository) GetAllRisks(meetingID uint, params *meetingDTOs.MeetingAssociatedItemListQueryParamsDTO) ([]*meetingDTOs.MeetingRiskListDTO, error) {
	var meetingRisk []*meetingDTOs.MeetingRiskListDTO

	query := repo.db.Table("meeting_risks").
		Select(`
			meeting_risks.risk_register_id AS risk_register_id,
			risk_register.title AS risk_register_title,
			risk_statuses.title_en AS risk_status_title_en,
			risk_statuses.title_ar AS risk_status_title_ar,
			risk_statuses.color_code AS risk_status_color_code
		`).
		Joins(`
			LEFT JOIN risk_register ON risk_register.id = meeting_risks.risk_register_id
			LEFT JOIN risk_statuses ON risk_statuses.id = risk_register.risk_status_id
		`).
		Where("meeting_risks.meeting_id = ?", meetingID)

	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
			(risk_register.title LIKE ?)`,
			keyword)
	}

	if params.SortBy == "" {
		params.SortBy = "risk_register.title"
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "asc"
	}

	query = query.Order(params.SortBy + " " + params.Sort)

	if err := query.Find(&meetingRisk).Error; err != nil {
		return nil, err
	}

	return meetingRisk, nil
}

func (repo *MeetingAssociatedItemRepository) GetAllProcesses(meetingID uint, params *meetingDTOs.MeetingAssociatedItemListQueryParamsDTO) ([]*meetingDTOs.MeetingProcessListDTO, error) {
	var meetingProcess []*meetingDTOs.MeetingProcessListDTO

	query := repo.db.Table("meeting_risks").
		Select(`
			meeting_risks.risk_register_id AS risk_register_id,
			risk_register.title AS risk_register_title,
			risk_statuses.title_en AS risk_status_title_en,
			risk_statuses.title_ar AS risk_status_title_ar,
			risk_statuses.color_code AS risk_status_color_code
		`).
		Joins(`
			LEFT JOIN risk_register ON risk_register.id = meeting_risks.risk_register_id
			LEFT JOIN risk_statuses ON risk_statuses.id = risk_register.risk_status_id
		`).
		Where("meeting_risks.meeting_id = ?", meetingID)

	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
			(risk_register.title LIKE ?)`,
			keyword)
	}

	if params.SortBy == "" {
		params.SortBy = "risk_register.title"
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "asc"
	}

	query = query.Order(params.SortBy + " " + params.Sort)

	if err := query.Find(&meetingProcess).Error; err != nil {
		return nil, err
	}

	return meetingProcess, nil
}
