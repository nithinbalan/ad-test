package repositories

import (
	"fmt"
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/models"
	taskDTOs "jethur/grc/task/dtos"
	"jethur/utility"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type TaskRepositoryInteface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(params *taskDTOs.TaskListQueryParamsDTO) (*taskDTOs.TaskListResponseDTO, error)
	GetByID(id uint) (*models.Task, error)
	GetByFieldID(field string, value interface{}) (*models.Task, error)
	GetByResponsibleUserID(field string, value interface{}, responsibleuserID uint) (*models.Task, error)
	Create(task *models.Task, tx ...utility.DatabaseTransaction) (uint, error)
	Update(task *models.Task, tx ...utility.DatabaseTransaction) error
	Delete(id uint) error
	GenerateReferenceNumber() (string, error)

	UpdateTaskRelatedIDs(taskID uint, relatedTaskIDsJSON string, tx ...utility.DatabaseTransaction) error
	UpdateDocumentTaskStatus(documentID uint, userID uint, documentWorkflowID uint, status uint) (*models.Task, error)
	UpdateDocumentChangeRequestTaskStatus(documentCRID uint, userID uint, documentCRWorkflowID uint, status uint) (*models.Task, error)
	UpdateAuditReportTaskStatus(auditFinalReportID uint, userID uint, auditReportWorkflowID uint, status uint) (*models.Task, error)
	UpdateAuditDraftReportTaskStatus(auditDraftReportID uint, userID uint, auditReportWorkflowID uint, status uint) (*models.Task, error)
	UpdateComplianceActionPlanTaskStatus(compliaceActionPlanID uint, userID uint, status uint) (*models.Task, error)
}

type TaskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) *TaskRepository {
	return &TaskRepository{
		db: db,
	}
}

// BeginTransaction starts a new transaction.
func (repo *TaskRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *TaskRepository) GetAll(params *taskDTOs.TaskListQueryParamsDTO) (*taskDTOs.TaskListResponseDTO, error) {
	var tasks []*taskDTOs.TaskListDTO
	var totalCount int64

	query := repo.db.Table("tasks").
		Select(`
            tasks.*,
            modules.title_en AS module_title_en,
            modules.title_ar AS module_title_ar,
			CONVERT(VARCHAR, tasks.target_date, 23) AS target_date,
            task_statuses.title_en AS task_status_title_en,
            task_statuses.title_ar AS task_status_title_ar,
			task_statuses.color_code AS task_status_color_code,
			severity_levels.title_en AS severity_level_title_en,
			severity_levels.title_ar AS severity_level_title_ar,
			severity_levels.color_code AS severity_level_color_code,
			responsible_user.id AS responsible_user_id,
            responsible_user.first_name_en AS responsible_user_first_name_en,
            responsible_user.first_name_ar AS responsible_user_first_name_ar,
            responsible_user.last_name_en AS responsible_user_last_name_en,
            responsible_user.last_name_ar AS responsible_user_last_name_ar,
            responsible_user.initial AS responsible_user_initial,
            responsible_user.color AS responsible_user_color,
            responsible_user.email AS responsible_user_email,
            responsible_user.profile_image_path AS responsible_user_profile_image_path,
			created_user.id AS created_user_id,
       		created_user.first_name_en AS created_user_first_name_en,
            created_user.first_name_ar AS created_user_first_name_ar,
            created_user.last_name_en AS created_user_last_name_en,
            created_user.last_name_ar AS created_user_last_name_ar,
			created_user.initial AS created_user_initial,
            created_user.color AS created_user_color,
            created_user.email AS created_user_email,
            created_user.profile_image_path AS created_user_profile_image_path,   
			updated_user.id AS updated_user_id,
    		updated_user.first_name_en AS updated_user_first_name_en,
            updated_user.first_name_ar AS updated_user_first_name_ar,
            updated_user.last_name_en AS updated_user_last_name_en,
            updated_user.last_name_ar AS updated_user_last_name_ar,
			updated_user.initial AS updated_user_initial,
            updated_user.color AS updated_user_color,
            updated_user.email AS updated_user_email,
            updated_user.profile_image_path AS updated_user_profile_image_path     
        `).
		Joins(`
            LEFT JOIN modules ON modules.id = tasks.module_id
            LEFT JOIN task_statuses ON task_statuses.id = tasks.task_status_id
		    LEFT JOIN severity_levels ON severity_levels.id = tasks.severity_level_id
            LEFT JOIN users AS responsible_user ON responsible_user.id = tasks.responsible_user_id
			LEFT JOIN users AS created_user ON created_user.id = tasks.created_by_id
			LEFT JOIN users AS updated_user ON updated_user.id = tasks.updated_by_id
        `)
	// Apply dynamic filters to the query
	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where(columnName+" IN (?)", values)
		}
	}

	// Apply target_date filters
	if params.TargetDateFrom != "" {
		query = query.Where("target_date >= ?", params.TargetDateFrom)
	}
	if params.TargetDateTo != "" {
		query = query.Where("target_date <= ?", params.TargetDateTo)
	}

	// Calculate total count for pagination
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}
	// Apply search keyword if exists
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
			(tasks.title LIKE ? OR
			tasks.description LIKE ? OR
			modules.title_en LIKE ? OR
			modules.title_ar LIKE ? OR
			task_statuses.title_en LIKE ? OR
			task_statuses.title_ar LIKE ? OR
			severity_levels.title_en LIKE ? OR
			severity_levels.title_ar LIKE ? OR
			responsible_user.first_name_en LIKE ? OR
			responsible_user.first_name_ar LIKE ? OR
			responsible_user.last_name_en LIKE ? OR
			responsible_user.last_name_ar LIKE ? OR
			responsible_user.email LIKE ? OR
			created_user.first_name_en LIKE ? OR
			created_user.first_name_ar LIKE ? OR
			created_user.last_name_en LIKE ? OR
			created_user.last_name_ar LIKE ? OR
			created_user.email LIKE ? OR
			updated_user.id AS updated_user_id,
    		updated_user.first_name_en LIKE ? OR
			updated_user.first_name_ar LIKE ? OR
			updated_user.last_name_en LIKE ? OR
			updated_user.last_name_en LIKE ? OR
			updated_user.email LIKE ?)`,
			keyword, keyword, keyword, keyword, keyword, keyword, keyword,
			keyword, keyword, keyword, keyword, keyword, keyword, keyword,
			keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword)
	}

	// Sort by the specified column and order
	if params.SortBy == "" {
		params.SortBy = "title" // Default sort by title if no sortBy is specified
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "asc" // Default order is ascending if not specified
	}

	query = query.Order(params.SortBy + " " + params.Sort)

	offset := (params.Page - 1) * params.PageSize

	// Execute the query
	err := query.Offset(offset).Limit(params.PageSize).Find(&tasks).Error
	if err != nil {
		return nil, err
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	response := &taskDTOs.TaskListResponseDTO{
		Items: tasks,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}
	return response, nil
}

func (repo *TaskRepository) GetByID(id uint) (*models.Task, error) {
	var task *models.Task

	result := repo.db.First(&task, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return task, nil

}

func (repo *TaskRepository) GetByFieldID(field string, value interface{}) (*models.Task, error) {
	var task models.Task

	result := repo.db.Where(field+" = ?", value).First(&task)
	if result.Error != nil {
		return nil, result.Error
	}

	return &task, nil
}

func (repo *TaskRepository) GetByResponsibleUserID(field string, value interface{}, responsibleUserID uint) (*models.Task, error) {
	var task models.Task

	result := repo.db.Where(field+" = ? AND responsible_user_id = ?", value, responsibleUserID).Last(&task)
	if result.Error != nil {
		return nil, result.Error
	}

	return &task, nil
}

func (repo *TaskRepository) Create(task *models.Task, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(task).Error; err != nil {
		return 0, err
	}
	return task.ID, nil
}

func (repo *TaskRepository) GenerateReferenceNumber() (string, error) {
	var lastReferenceNumber string
	err := repo.db.Model(&models.Task{}).
		Select("reference_number").
		Order("reference_number desc").
		Limit(1).
		Pluck("reference_number", &lastReferenceNumber).
		Error

	if err != nil {
		return "", err
	}

	lastDigits := "0000"
	if lastReferenceNumber != "" && len(lastReferenceNumber) >= 4 {
		lastDigits = lastReferenceNumber[len(lastReferenceNumber)-4:]
	}

	currentYear := time.Now().Year()
	lastReferenceID, err := strconv.ParseUint(lastDigits, 10, 64)
	if err != nil {
		return "", err
	}
	newReferenceNumber := fmt.Sprintf("%s-%d-%04d", "T", currentYear, lastReferenceID+1)

	return newReferenceNumber, nil
}

func (repo *TaskRepository) Update(task *models.Task, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(task).Error; err != nil {
		return err
	}
	return nil
}

func (repo *TaskRepository) Delete(id uint) error {

	var task *models.Task
	if err := repo.db.First(&task, id).Error; err != nil {
		return err
	}

	result := repo.db.Delete(&models.Task{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *TaskRepository) UpdateTaskRelatedIDs(taskID uint, relatedTaskIDsJSON string, tx ...utility.DatabaseTransaction) error {
	// Use the transaction's database if provided
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	// Update the related_task_ids column for the given task ID
	return db.Model(&models.Task{}).
		Where("id = ?", taskID).
		Update("related_task_ids", relatedTaskIDsJSON).Error
}

func (repo *TaskRepository) UpdateDocumentTaskStatus(documentID uint, userID uint, documentWorkflowID uint, status uint) (*models.Task, error) {
	var task models.Task

	// Find the task by document and user ID
	if err := repo.db.Where("document_id = ? AND responsible_user_id = ? AND document_workflow_id = ?", documentID, userID, documentWorkflowID).First(&task).Error; err != nil {
		return nil, err
	}

	// Update the task status
	task.TaskStatusID = status
	task.UpdatedByID = userID
	if err := repo.db.Save(&task).Error; err != nil {
		return nil, err
	}

	return &task, nil
}

func (repo *TaskRepository) UpdateDocumentChangeRequestTaskStatus(documentCRID uint, userID uint, documentCRWorkflowID uint, status uint) (*models.Task, error) {
	var task models.Task

	// Find the task by document and user ID
	if err := repo.db.Where("document_change_request_id = ? AND responsible_user_id = ? AND document_change_request_workflow_id = ?", documentCRID, userID, documentCRWorkflowID).First(&task).Error; err != nil {
		return nil, err
	}

	// Update the task status
	task.TaskStatusID = status
	task.UpdatedByID = userID
	if err := repo.db.Save(&task).Error; err != nil {
		return nil, err
	}

	return &task, nil
}

func (repo *TaskRepository) UpdateAuditReportTaskStatus(auditFinalReportID uint, userID uint, auditReportWorkflowID uint, status uint) (*models.Task, error) {
	var task models.Task

	// Find the task by audit_final_report_id and user ID
	if err := repo.db.Where("audit_final_report_id = ? AND responsible_user_id = ? AND audit_report_workflow_id = ?", auditFinalReportID, userID, auditReportWorkflowID).First(&task).Error; err != nil {
		return nil, err
	}

	// Update the task status
	task.TaskStatusID = status
	task.UpdatedByID = userID
	if err := repo.db.Save(&task).Error; err != nil {
		return nil, err
	}

	return &task, nil
}

func (repo *TaskRepository) UpdateAuditDraftReportTaskStatus(auditDraftReportID uint, userID uint, auditReportWorkflowID uint, status uint) (*models.Task, error) {
	var task models.Task

	// Find the task by audit_final_report_id and user ID
	if err := repo.db.Where("audit_draft_report_id = ? AND responsible_user_id = ? AND audit_report_workflow_id = ?", auditDraftReportID, userID, auditReportWorkflowID).First(&task).Error; err != nil {
		return nil, err
	}

	// Update the task status
	task.TaskStatusID = status
	task.UpdatedByID = userID
	if err := repo.db.Save(&task).Error; err != nil {
		return nil, err
	}

	return &task, nil
}

func (repo *TaskRepository) UpdateComplianceActionPlanTaskStatus(compliaceActionPlanID uint, userID uint, status uint) (*models.Task, error) {
	var task models.Task

	// Find the task by compliance_action_plan_id and user ID
	if err := repo.db.Where("compliance_action_plan_id = ? AND responsible_user_id = ?", compliaceActionPlanID, userID).First(&task).Error; err != nil {
		return nil, err
	}

	// Update the task status
	task.TaskStatusID = status
	task.UpdatedByID = userID
	if err := repo.db.Save(&task).Error; err != nil {
		return nil, err
	}

	return &task, nil
}
