package repositories

import (
	"fmt"
	reportDTOs "jethur/grc/reports/dtos"

	"gorm.io/gorm"
)

type MeetingReportRepositoryInterface interface {
	GetMeetingByType(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.MeetingByTypeDTO, error)
	GetMeetingByStatus(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.MeetingByStatusDTO, error)
	GetMeetingByOrganizer(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.MeetingByOrganizerDTO, error)
	GetMeetingByResponsibleUnit(params *reportDTOs.ReportListQueryParamsDTO, primaryOrgUnitField string) ([]*reportDTOs.MeetingByResponsibleUnitDTO, error)
	GetMeetingMinutesByStatus(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.MeetingMinutesByStatusDTO, error)

	GetMeetingActionPlanByStatus(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.MeetingActionPlanByStatusDTO, error)
	GetMeetingActionPlanByResponsibleUser(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.MeetingActionPlanByResponsibleUserDTO, error)
	GetMeetingActionPlanByResponsibleUnit(params *reportDTOs.ReportListQueryParamsDTO, primaryOrgUnitField string) ([]*reportDTOs.MeetingActionPlanByResponsibleUnitDTO, error)
	GetActionPlanByOpenVsClosed(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.ActionPlanByOpenVsClosedDTO, error)

	GetMeetingMOMCount(params *reportDTOs.ReportListQueryParamsDTO) (*reportDTOs.MeetingMOMCountDTO, error)
	GetMeetingsWithMOM(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.MeetingWithMOMDTO, error)
	GetMeetingsWithOpenActionPlans(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.MeetingWithOpenActionPlanDTO, error)
	GetMeetingsWithOverdueActionPlans(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.MeetingWithOverdueActionPlanDTO, error)
	GetMeetingCountByYear(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.MeetingCountByYearDTO, error)
	GetMeetingCountByUser(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.MeetingCountByUserDTO, error)
}

type MeetingReportRepository struct {
	db *gorm.DB
}

func NewMeetingReportRepository(db *gorm.DB) *MeetingReportRepository {
	return &MeetingReportRepository{db: db}
}

func (repo *MeetingReportRepository) GetMeetingByType(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.MeetingByTypeDTO, error) {
	results := []*reportDTOs.MeetingByTypeDTO{}

	query := repo.db.Table("meeting_types").
		Select(`meeting_types.id,
                meeting_types.title_en,
                meeting_types.title_ar,
                COUNT(meetings.id) AS count`).
		Joins("LEFT JOIN meetings ON meetings.meeting_type_id = meeting_types.id").
		Group("meeting_types.id, meeting_types.title_en, meeting_types.title_ar")

	// Apply date filter if provided
	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("meetings.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.Scan(&results).Error
	return results, err
}

func (repo *MeetingReportRepository) GetMeetingByStatus(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.MeetingByStatusDTO, error) {
	results := []*reportDTOs.MeetingByStatusDTO{}

	query := repo.db.Table("meeting_statuses").
		Select(`meeting_statuses.id,
                meeting_statuses.title_en,
                meeting_statuses.title_ar,
                meeting_statuses.color_code,
                COUNT(meetings.id) AS count`).
		Joins("LEFT JOIN meetings ON meetings.meeting_status_id = meeting_statuses.id").
		Group("meeting_statuses.id, meeting_statuses.title_en, meeting_statuses.title_ar, meeting_statuses.color_code")

	// Apply date filter if provided
	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("meetings.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.Scan(&results).Error
	return results, err
}

func (repo *MeetingReportRepository) GetMeetingByOrganizer(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.MeetingByOrganizerDTO, error) {
	results := []*reportDTOs.MeetingByOrganizerDTO{}

	query := repo.db.Table("meetings").
		Select(`
            users.id AS user_id,
            users.first_name_en AS user_first_name_en,
            users.first_name_ar AS user_first_name_ar,
            users.last_name_en AS user_last_name_en,
            users.last_name_ar AS user_last_name_ar,
            CONCAT(users.first_name_en, ' ', users.last_name_en) AS user_title_en,
            CONCAT(users.first_name_ar, ' ', users.last_name_ar) AS user_title_ar,
            users.initial AS user_initial,
            users.color AS user_color,
            users.email AS user_email,
            users.profile_image_path AS user_profile_image_path,
            COUNT(meetings.id) AS count
        `).
		Joins("JOIN users ON users.id = meetings.organizer_id").
		Group("users.id, users.first_name_en, users.first_name_ar, users.last_name_en, users.last_name_ar, users.initial, users.color, users.email, users.profile_image_path")

	// Apply date filter if provided
	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("meetings.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.Scan(&results).Error
	return results, err
}

func (repo *MeetingReportRepository) GetMeetingByResponsibleUnit(params *reportDTOs.ReportListQueryParamsDTO, primaryOrgUnitField string) ([]*reportDTOs.MeetingByResponsibleUnitDTO, error) {
	results := []*reportDTOs.MeetingByResponsibleUnitDTO{}

	// Validate the primaryOrgUnitField to prevent SQL injection
	validFields := map[string]bool{
		"organization_unit_1_id": true,
		"organization_unit_2_id": true,
		"organization_unit_3_id": true,
		"organization_unit_4_id": true,
		"organization_unit_5_id": true,
	}

	if !validFields[primaryOrgUnitField] {
		return nil, fmt.Errorf("invalid primary organization unit field: %s", primaryOrgUnitField)
	}

	query := repo.db.Table("meetings").
		Select(`
            organization_units.id AS responsible_unit_id,
            organization_units.title_en AS responsible_unit_title_en,
            organization_units.title_ar AS responsible_unit_title_ar,
            COUNT(meetings.id) AS count
        `).
		Joins(`JOIN users ON users.id = meetings.organizer_id`).
		Joins(fmt.Sprintf(`
            JOIN organization_units ON organization_units.id = users.%s
        `, primaryOrgUnitField)).
		Group("organization_units.id, organization_units.title_en, organization_units.title_ar")

	// Apply date filter if provided
	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("meetings.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.Scan(&results).Error
	return results, err
}

func (repo *MeetingReportRepository) GetMeetingMinutesByStatus(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.MeetingMinutesByStatusDTO, error) {
	results := []*reportDTOs.MeetingMinutesByStatusDTO{}

	query := repo.db.Table("meeting_minute_statuses").
		Select(`
            meeting_minute_statuses.id,
            meeting_minute_statuses.title_en,
            meeting_minute_statuses.title_ar,
            meeting_minute_statuses.color_code,
            COUNT(meeting_minutes.id) AS count
        `).
		Joins("LEFT JOIN meeting_minutes ON meeting_minutes.meeting_minute_status_id = meeting_minute_statuses.id").
		Group("meeting_minute_statuses.id, meeting_minute_statuses.title_en, meeting_minute_statuses.title_ar, meeting_minute_statuses.color_code")

	// Apply date filter if provided
	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("meeting_minutes.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.Scan(&results).Error
	return results, err
}

func (repo *MeetingReportRepository) GetMeetingActionPlanByStatus(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.MeetingActionPlanByStatusDTO, error) {
	results := []*reportDTOs.MeetingActionPlanByStatusDTO{}

	query := repo.db.Table("action_plan_statuses").
		Select(`
            action_plan_statuses.id,
            action_plan_statuses.title_en,
            action_plan_statuses.title_ar,
            action_plan_statuses.color_code,
            COUNT(meeting_action_plans.id) AS count
        `).
		Joins("LEFT JOIN meeting_action_plans ON meeting_action_plans.action_plan_status_id = action_plan_statuses.id").
		Group("action_plan_statuses.id, action_plan_statuses.title_en, action_plan_statuses.title_ar, action_plan_statuses.color_code")

	// Apply date filter if provided
	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("meeting_action_plans.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.Scan(&results).Error
	return results, err
}

func (repo *MeetingReportRepository) GetMeetingActionPlanByResponsibleUser(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.MeetingActionPlanByResponsibleUserDTO, error) {
	results := []*reportDTOs.MeetingActionPlanByResponsibleUserDTO{}

	query := repo.db.Table("meeting_action_plans").
		Select(`
            users.id AS id,
            users.first_name_en AS user_first_name_en,
            users.first_name_ar AS user_first_name_ar,
            users.last_name_en AS user_last_name_en,
            users.last_name_ar AS user_last_name_ar,
            CONCAT(users.first_name_en, ' ', users.last_name_en) AS user_title_en,
            CONCAT(users.first_name_ar, ' ', users.last_name_ar) AS user_title_ar,
            users.initial AS user_initial,
            users.color AS user_color,
            users.email AS user_email,
            users.profile_image_path AS user_profile_image_path,
            COUNT(meeting_action_plans.id) AS count
        `).
		Joins("LEFT JOIN users ON users.id = meeting_action_plans.responsible_user_id").
		Group("users.id, users.first_name_en, users.last_name_en, users.first_name_ar, users.last_name_ar, users.initial, users.color, users.email, users.profile_image_path")

	// Apply date filter if provided
	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("meeting_action_plans.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.Scan(&results).Error
	return results, err
}

func (repo *MeetingReportRepository) GetMeetingActionPlanByResponsibleUnit(params *reportDTOs.ReportListQueryParamsDTO, primaryOrgUnitField string) ([]*reportDTOs.MeetingActionPlanByResponsibleUnitDTO, error) {
	results := []*reportDTOs.MeetingActionPlanByResponsibleUnitDTO{}

	// Validate the primaryOrgUnitField to prevent SQL injection
	validFields := map[string]bool{
		"organization_unit_1_id": true,
		"organization_unit_2_id": true,
		"organization_unit_3_id": true,
		"organization_unit_4_id": true,
		"organization_unit_5_id": true,
	}

	if !validFields[primaryOrgUnitField] {
		return nil, fmt.Errorf("invalid primary organization unit field: %s", primaryOrgUnitField)
	}

	query := repo.db.Table("meeting_action_plans").
		Select(`
            organization_units.id AS responsible_unit_id,
            organization_units.title_en AS responsible_unit_title_en,
            organization_units.title_ar AS responsible_unit_title_ar,
            COUNT(meeting_action_plans.id) AS count
        `).
		Joins(`JOIN users ON users.id = meeting_action_plans.responsible_user_id`).
		Joins(fmt.Sprintf(`
            JOIN organization_units ON organization_units.id = users.%s
        `, primaryOrgUnitField)).
		Group("organization_units.id, organization_units.title_en, organization_units.title_ar")

	// Apply date filter if provided
	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("meeting_action_plans.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.Scan(&results).Error
	return results, err
}

func (repo *MeetingReportRepository) GetActionPlanByOpenVsClosed(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.ActionPlanByOpenVsClosedDTO, error) {
	var rawResults []struct {
		ID        uint   `gorm:"column:id"`
		Type      string `gorm:"column:type"`
		TitleEn   string `gorm:"column:title_en"`
		TitleAr   string `gorm:"column:title_ar"`
		ColorCode string `gorm:"column:color_code"`
	}

	// Step 1: Fetch action plan statuses with their details
	query := repo.db.Table("meeting_action_plans").
		Select(`
            action_plan_statuses.id, 
            action_plan_statuses.type, 
            action_plan_statuses.title_en, 
            action_plan_statuses.title_ar, 
            action_plan_statuses.color_code
        `).
		Joins("JOIN action_plan_statuses ON action_plan_statuses.id = meeting_action_plans.action_plan_status_id")

	// Apply date filter if provided
	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("meeting_action_plans.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.Scan(&rawResults).Error
	if err != nil {
		return nil, err
	}

	// Step 2: Initialize default statuses with 0 counts
	openStatus := &reportDTOs.ActionPlanByOpenVsClosedDTO{
		TitleEn:   "Open",
		TitleAr:   "مفتوح",
		Type:      "open",
		ColorCode: "", // Default, will be overwritten if records exist
		Count:     0,
	}
	closedStatus := &reportDTOs.ActionPlanByOpenVsClosedDTO{
		TitleEn:   "Closed",
		TitleAr:   "مغلق",
		Type:      "closed",
		ColorCode: "", // Default, will be overwritten if records exist
		Count:     0,
	}

	// Step 3: Process data in Go
	for _, record := range rawResults {
		// Completed, Terminated, and Closed are treated as "closed"
		if record.Type == "completed" || record.Type == "terminated" || record.Type == "close" {
			closedStatus.Count++
			closedStatus.ID = record.ID
			closedStatus.ColorCode = record.ColorCode
		} else {
			openStatus.Count++
			openStatus.ID = record.ID
			openStatus.ColorCode = record.ColorCode
		}
	}

	// Step 4: Prepare the response as a slice of pointers
	results := []*reportDTOs.ActionPlanByOpenVsClosedDTO{
		openStatus,
		closedStatus,
	}

	return results, nil
}

func (repo *MeetingReportRepository) GetMeetingMOMCount(params *reportDTOs.ReportListQueryParamsDTO) (*reportDTOs.MeetingMOMCountDTO, error) {
	var result struct {
		WithMOMCount    int `gorm:"column:with_mom_count"`
		WithoutMOMCount int `gorm:"column:without_mom_count"`
	}

	query := repo.db.Raw(`
        SELECT 
            (SELECT COUNT(DISTINCT meetings.id)
             FROM meetings
             JOIN meeting_minutes ON meeting_minutes.meeting_id = meetings.id
             WHERE meeting_minutes.minutes IS NOT NULL) AS with_mom_count,
            (SELECT COUNT(DISTINCT meetings.id)
             FROM meetings
             LEFT JOIN meeting_minutes ON meeting_minutes.meeting_id = meetings.id
             WHERE meeting_minutes.id IS NULL) AS without_mom_count
    `)

	// Apply date filter if provided
	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where(`
            meetings.created_at BETWEEN ? AND ?`,
			params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.Scan(&result).Error
	if err != nil {
		return nil, err
	}

	// Map results to DTO
	response := &reportDTOs.MeetingMOMCountDTO{
		WithMOMCount:    result.WithMOMCount,
		WithoutMOMCount: result.WithoutMOMCount,
	}

	return response, nil
}

func (repo *MeetingReportRepository) GetMeetingsWithMOM(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.MeetingWithMOMDTO, error) {
	results := []*reportDTOs.MeetingWithMOMDTO{}

	query := repo.db.Raw(`
        WITH MeetingMinutesCTE AS (
            SELECT
                meeting_minutes.meeting_id,
                STRING_AGG(meeting_minutes.minutes, ', ') AS minutes
            FROM meeting_minutes
            WHERE meeting_minutes.minutes IS NOT NULL
            GROUP BY meeting_minutes.meeting_id
        )
        SELECT 
            meetings.id AS meeting_id,
            meetings.title AS meeting_title,
            COALESCE(MeetingMinutesCTE.minutes, '') AS minutes,
            meeting_statuses.title_en AS status
        FROM 
            meetings
        LEFT JOIN 
            MeetingMinutesCTE ON MeetingMinutesCTE.meeting_id = meetings.id
        JOIN 
            meeting_statuses ON meeting_statuses.id = meetings.meeting_status_id
        WHERE 
            MeetingMinutesCTE.minutes IS NOT NULL
        GROUP BY 
            meetings.id, meetings.title, MeetingMinutesCTE.minutes, meeting_statuses.title_en
    `)

	// Apply date filter if provided
	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("meetings.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.Scan(&results).Error
	return results, err
}

func (repo *MeetingReportRepository) GetMeetingsWithOpenActionPlans(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.MeetingWithOpenActionPlanDTO, error) {
	results := []*reportDTOs.MeetingWithOpenActionPlanDTO{}

	query := repo.db.Raw(`
        SELECT 
            meetings.id AS meeting_id,
            meetings.title AS meeting_title,
            COUNT(meeting_action_plans.id) AS open_action_plan_count
        FROM 
            meetings
        JOIN 
            meeting_action_plans ON meeting_action_plans.meeting_id = meetings.id
        JOIN 
            action_plan_statuses ON action_plan_statuses.id = meeting_action_plans.action_plan_status_id
        WHERE 
            action_plan_statuses.type IN ('new', 'in_progress', 'reopen')
        GROUP BY 
            meetings.id, meetings.title
    `)

	// Apply date filter if provided
	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("meetings.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.Scan(&results).Error
	return results, err
}

func (repo *MeetingReportRepository) GetMeetingsWithOverdueActionPlans(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.MeetingWithOverdueActionPlanDTO, error) {
	results := []*reportDTOs.MeetingWithOverdueActionPlanDTO{}

	query := repo.db.Raw(`
        SELECT 
            meetings.id AS meeting_id,
            meetings.title AS meeting_title,
            COUNT(meeting_action_plans.id) AS overdue_action_plan_count
        FROM 
            meetings
        JOIN 
            meeting_action_plans ON meeting_action_plans.meeting_id = meetings.id
        JOIN 
            action_plan_statuses ON action_plan_statuses.id = meeting_action_plans.action_plan_status_id
        WHERE 
            meeting_action_plans.target_date < CAST(GETDATE() AS DATE) -- Fix: Use GETDATE() instead of CURRENT_DATE
            AND action_plan_statuses.type IN ('new', 'in_progress', 'reopen')
        GROUP BY 
            meetings.id, meetings.title
    `)

	// Apply date filter if provided
	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("meetings.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.Scan(&results).Error
	return results, err
}

func (repo *MeetingReportRepository) GetMeetingCountByYear(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.MeetingCountByYearDTO, error) {
	results := []*reportDTOs.MeetingCountByYearDTO{}

	// Fetch yearly meeting count
	var rawResults []struct {
		Year  int  `gorm:"column:year"`
		Count uint `gorm:"column:count"`
	}

	query := repo.db.Table("meetings").
		Select(`YEAR(meetings.created_at) AS year, COUNT(meetings.id) AS count`).
		Group("YEAR(meetings.created_at)").
		Order("YEAR(meetings.created_at) ASC")

	// Apply date filter if provided
	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("meetings.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.Scan(&rawResults).Error
	if err != nil {
		return nil, err
	}

	// Convert year to English and Arabic title
	for _, record := range rawResults {
		results = append(results, &reportDTOs.MeetingCountByYearDTO{
			TitleEn: fmt.Sprintf("%d", record.Year), // English Title (Year)
			TitleAr: fmt.Sprintf("%d", record.Year), // Arabic Title (Same as Year)
			Count:   record.Count,
		})
	}

	return results, nil
}

func (repo *MeetingReportRepository) GetMeetingCountByUser(params *reportDTOs.ReportListQueryParamsDTO) ([]*reportDTOs.MeetingCountByUserDTO, error) {
	results := []*reportDTOs.MeetingCountByUserDTO{}

	query := repo.db.Table("meeting_participants").
		Select(`
            users.id AS user_id,
            users.first_name_en AS user_first_name_en,
            users.first_name_ar AS user_first_name_ar,
            users.last_name_en AS user_last_name_en,
            users.last_name_ar AS user_last_name_ar,
            CONCAT(users.first_name_en, ' ', users.last_name_en) AS user_title_en,
            CONCAT(users.first_name_ar, ' ', users.last_name_ar) AS user_title_ar,
            users.initial AS user_initial,
            users.color AS user_color,
            users.email AS user_email,
            users.profile_image_path AS user_profile_image_path,
            COUNT(meeting_participants.meeting_id) AS meeting_count
        `).
		Joins("LEFT JOIN users ON users.id = meeting_participants.user_id").
		Group("users.id, users.first_name_en, users.last_name_en, users.first_name_ar, users.last_name_ar, users.initial, users.color, users.email, users.profile_image_path").
		Order("users.first_name_en ASC")

	// Apply date filter if provided
	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("meeting_participants.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	err := query.Scan(&results).Error
	if err != nil {
		return nil, err
	}

	return results, nil
}
