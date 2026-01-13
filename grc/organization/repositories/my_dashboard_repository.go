package repositories

import (
	"jethur/grc/models"
	organizationDTOs "jethur/grc/organization/dtos"
	"time"

	"gorm.io/gorm"
)

type MyDashboardRepositoryInterface interface {
	GetDashboardSummary(userID uint) (tasksAssigned, tasksCompleted, overdueItems, newNotifications int64, err error)
	GetUserTasks(userID uint, limit int) ([]organizationDTOs.UserTaskDTO, error)
	GetOverdueTasks(userID uint, limit int) ([]organizationDTOs.UserTaskDTO, error)
	GetAssignedRisks(userID uint, limit int) ([]organizationDTOs.AssignedRiskFlatDTO, error)
	GetRecentActivity(userID uint, limit int) ([]*models.ActivityLog, error)
}

type MyDashboardRepository struct {
	db *gorm.DB
}

func NewMyDashboardRepository(db *gorm.DB) *MyDashboardRepository {
	return &MyDashboardRepository{db: db}
}

func (dr *MyDashboardRepository) GetDashboardSummary(userID uint) (tasksAssigned, tasksCompleted, overdueItems, newNotifications int64, err error) {
	today := time.Now().Format("2006-01-02")

	taskCountsQuery := `
		SELECT
			COUNT(*) as total_assigned,
			COUNT(CASE WHEN ts.type = 'closed' THEN 1 END) as total_completed,
			COUNT(CASE WHEN t.target_date < ? AND t.target_date IS NOT NULL AND ts.type != 'closed' THEN 1 END) as total_overdue
		FROM tasks t
		LEFT JOIN task_statuses ts ON t.task_status_id = ts.id
		WHERE t.responsible_user_id = ?
	`

	row := dr.db.Raw(taskCountsQuery, today, userID).Row()
	if scanErr := row.Scan(&tasksAssigned, &tasksCompleted, &overdueItems); scanErr != nil {
		err = scanErr
		return
	}

	if err = dr.db.Model(&models.UserNotification{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Count(&newNotifications).Error; err != nil {
		return
	}

	return
}

func (dr *MyDashboardRepository) GetUserTasks(userID uint, limit int) ([]organizationDTOs.UserTaskDTO, error) {
	var userTasks []organizationDTOs.UserTaskDTO

	query := `
		SELECT
			t.id, t.reference_number, t.module_id, t.title, t.description,
			sl.title_en as severity_level_title_en, sl.title_ar as severity_level_title_ar, sl.color_code as severity_level_color_code,
			t.responsible_user_id,
			responsible_user.first_name_en as responsible_user_first_name_en,
			responsible_user.first_name_ar as responsible_user_first_name_ar,
			responsible_user.last_name_en as responsible_user_last_name_en,
			responsible_user.last_name_ar as responsible_user_last_name_ar,
			responsible_user.initial as responsible_user_initial,
			responsible_user.color as responsible_user_color,
			responsible_user.email as responsible_user_email,
			responsible_user.profile_image_path as responsible_user_profile_image_path,
			t.delegatee_user_id,
			delegatee_user.first_name_en as delegatee_user_first_name_en,
			delegatee_user.first_name_ar as delegatee_user_first_name_ar,
			delegatee_user.last_name_en as delegatee_user_last_name_en,
			delegatee_user.last_name_ar as delegatee_user_last_name_ar,
			delegatee_user.initial as delegatee_user_initial,
			delegatee_user.color as delegatee_user_color,
			delegatee_user.email as delegatee_user_email,
			delegatee_user.profile_image_path as delegatee_user_profile_image_path,
			t.target_date as target_date,
			t.url_param,
			ts.title_en as task_status_title_en,
			ts.title_ar as task_status_title_ar,
			ts.color_code as task_status_color_code,
			t.created_at as created_at,
			t.updated_at as updated_at,
			t.created_by_id as created_user_id,
			created_user.first_name_en as created_user_first_name_en,
			created_user.first_name_ar as created_user_first_name_ar,
			created_user.last_name_en as created_user_last_name_en,
			created_user.last_name_ar as created_user_last_name_ar,
			created_user.initial as created_user_initial,
			created_user.color as created_user_color,
			created_user.email as created_user_email,
			created_user.profile_image_path as created_user_profile_image_path,
			t.updated_by_id as updated_user_id,
			updated_user.first_name_en as updated_user_first_name_en,
			updated_user.first_name_ar as updated_user_first_name_ar,
			updated_user.last_name_en as updated_user_last_name_en,
			updated_user.last_name_ar as updated_user_last_name_ar,
			updated_user.initial as updated_user_initial,
			updated_user.color as updated_user_color,
			updated_user.email as updated_user_email,
			updated_user.profile_image_path as updated_user_profile_image_path
		FROM tasks t
		LEFT JOIN users responsible_user ON t.responsible_user_id = responsible_user.id
		LEFT JOIN users delegatee_user ON t.delegatee_user_id = delegatee_user.id
		LEFT JOIN users created_user ON t.created_by_id = created_user.id
		LEFT JOIN users updated_user ON t.updated_by_id = updated_user.id
		LEFT JOIN task_statuses ts ON t.task_status_id = ts.id
		LEFT JOIN severity_levels sl ON t.severity_level_id = sl.id
		WHERE t.responsible_user_id = ?
		ORDER BY t.created_at DESC
		OFFSET 0 ROWS FETCH NEXT ? ROWS ONLY
	`

	if err := dr.db.Raw(query, userID, limit).Scan(&userTasks).Error; err != nil {
		return nil, err
	}

	return userTasks, nil
}

func (dr *MyDashboardRepository) GetOverdueTasks(userID uint, limit int) ([]organizationDTOs.UserTaskDTO, error) {
	var overdueTasks []organizationDTOs.UserTaskDTO
	today := time.Now().Format("2006-01-02")

	query := `
		SELECT
			t.id, t.reference_number, t.module_id, t.title, t.description,
			sl.title_en as severity_level_title_en, sl.title_ar as severity_level_title_ar, sl.color_code as severity_level_color_code,
			t.responsible_user_id,
			responsible_user.first_name_en as responsible_user_first_name_en,
			responsible_user.first_name_ar as responsible_user_first_name_ar,
			responsible_user.last_name_en as responsible_user_last_name_en,
			responsible_user.last_name_ar as responsible_user_last_name_ar,
			responsible_user.initial as responsible_user_initial,
			responsible_user.color as responsible_user_color,
			responsible_user.email as responsible_user_email,
			responsible_user.profile_image_path as responsible_user_profile_image_path,
			t.delegatee_user_id,
			delegatee_user.first_name_en as delegatee_user_first_name_en,
			delegatee_user.first_name_ar as delegatee_user_first_name_ar,
			delegatee_user.last_name_en as delegatee_user_last_name_en,
			delegatee_user.last_name_ar as delegatee_user_last_name_ar,
			delegatee_user.initial as delegatee_user_initial,
			delegatee_user.color as delegatee_user_color,
			delegatee_user.email as delegatee_user_email,
			delegatee_user.profile_image_path as delegatee_user_profile_image_path,
			t.target_date as target_date,
			t.url_param,
			ts.title_en as task_status_title_en,
			ts.title_ar as task_status_title_ar,
			ts.color_code as task_status_color_code,
			t.created_at as created_at,
			t.updated_at as updated_at,
			t.created_by_id as created_user_id,
			created_user.first_name_en as created_user_first_name_en,
			created_user.first_name_ar as created_user_first_name_ar,
			created_user.last_name_en as created_user_last_name_en,
			created_user.last_name_ar as created_user_last_name_ar,
			created_user.initial as created_user_initial,
			created_user.color as created_user_color,
			created_user.email as created_user_email,
			created_user.profile_image_path as created_user_profile_image_path,
			t.updated_by_id as updated_user_id,
			updated_user.first_name_en as updated_user_first_name_en,
			updated_user.first_name_ar as updated_user_first_name_ar,
			updated_user.last_name_en as updated_user_last_name_en,
			updated_user.last_name_ar as updated_user_last_name_ar,
			updated_user.initial as updated_user_initial,
			updated_user.color as updated_user_color,
			updated_user.email as updated_user_email,
			updated_user.profile_image_path as updated_user_profile_image_path
		FROM tasks t
		LEFT JOIN users responsible_user ON t.responsible_user_id = responsible_user.id
		LEFT JOIN users delegatee_user ON t.delegatee_user_id = delegatee_user.id
		LEFT JOIN users created_user ON t.created_by_id = created_user.id
		LEFT JOIN users updated_user ON t.updated_by_id = updated_user.id
		LEFT JOIN task_statuses ts ON t.task_status_id = ts.id
		LEFT JOIN severity_levels sl ON t.severity_level_id = sl.id
		WHERE t.responsible_user_id = ?
			AND t.target_date < ?
			AND t.target_date IS NOT NULL
			AND ts.type != 'closed'
		ORDER BY t.target_date ASC, t.created_at DESC
		OFFSET 0 ROWS FETCH NEXT ? ROWS ONLY
	`

	if err := dr.db.Raw(query, userID, today, limit).Scan(&overdueTasks).Error; err != nil {
		return nil, err
	}

	return overdueTasks, nil
}

func (dr *MyDashboardRepository) GetAssignedRisks(userID uint, limit int) ([]organizationDTOs.AssignedRiskFlatDTO, error) {
	var flatRisks []organizationDTOs.AssignedRiskFlatDTO

	query := `
		SELECT
			rr.id,
			rr.reference_number,
			rr.title,
			ISNULL(sl.title_en, 'Not Set') as severity_title,
			ISNULL(sl.color_code, '#808080') as severity_color,
			rr.updated_at as last_updated
		FROM risk_register rr
		LEFT JOIN severity_levels sl ON rr.severity_level_id = sl.id
		WHERE rr.responsible_user_id = ?
		ORDER BY rr.updated_at DESC
		OFFSET 0 ROWS FETCH NEXT ? ROWS ONLY
	`

	if err := dr.db.Raw(query, userID, limit).Scan(&flatRisks).Error; err != nil {
		return nil, err
	}

	return flatRisks, nil
}

func (dr *MyDashboardRepository) GetRecentActivity(userID uint, limit int) ([]*models.ActivityLog, error) {
	var activityLogs []*models.ActivityLog

	query := dr.db.Table("activity_logs").
		Select(`
			activity_logs.*,
			users.id AS user_id,
			users.first_name_en AS user_first_name_en,
			users.first_name_ar AS user_first_name_ar,
			users.last_name_en AS user_last_name_en,
			users.last_name_ar AS user_last_name_ar,
			users.email AS user_email,
			users.profile_image_path AS user_profile_image_path,
			users.initial AS user_initial,
			users.color AS user_color,
			modules.id AS module_id,
			modules.identifier AS module_identifier,
			modules.title_en AS module_title_en,
			modules.title_ar AS module_title_ar,
			sub_modules.id AS sub_module_id,
			sub_modules.identifier AS sub_module_identifier,
			sub_modules.title_en AS sub_module_title_en,
			sub_modules.title_ar AS sub_module_title_ar
		`).
		Joins(`
			LEFT JOIN users ON users.id = activity_logs.user_id
		`).
		Joins(`
			LEFT JOIN modules ON modules.id = activity_logs.module_id
		`).
		Joins(`
			LEFT JOIN modules AS sub_modules ON sub_modules.id = activity_logs.sub_module_id
		`).
		Where("activity_logs.user_id = ?", userID).
		Order("activity_logs.created_at DESC").
		Limit(limit)

	err := query.Find(&activityLogs).Error
	if err != nil {
		return nil, err
	}

	return activityLogs, nil
}

