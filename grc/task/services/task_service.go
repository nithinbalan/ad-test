package services

import (
	"errors"

	commonDTOs "jethur/grc/common/dtos"
	masterServices "jethur/grc/masters/services"
	"jethur/grc/models"
	settingService "jethur/grc/settings/services"
	taskDTOs "jethur/grc/task/dtos"
	"jethur/grc/task/repositories"
	"jethur/utility"

	"gorm.io/gorm"
)

type TaskServiceInterface interface {
	GetAll(params *taskDTOs.TaskListQueryParamsDTO) (*taskDTOs.TaskListResponseDTO, error)
	GetByID(id uint) (*models.Task, error)
	GetByFieldID(field string, value interface{}) (*models.Task, error)
	Create(taskDTO commonDTOs.TaskFormDTO, currentUserID uint, tx ...utility.DatabaseTransaction) (uint, error)
	Update(currentUserID uint, taskDTO commonDTOs.TaskFormDTO, fieldName string, id uint, tx ...utility.DatabaseTransaction) error
	Delete(id uint, currentUserID uint) (uint, error)

	MapActionPlanStatusToTaskStatus(actionPlanStatusID uint, isReviewUserNull bool) (*models.TaskStatus, error)
	UpdateStatus(currentUserID uint, taskID uint, actionPlanStatusID uint, isReviewUserNull bool, tx ...utility.DatabaseTransaction) error
	UpdateTaskWithRelatedIDs(taskID uint, relatedTaskIDsJSON string, tx ...utility.DatabaseTransaction) error
	UpdateDocumentTaskStatus(documentID uint, userID uint, documentWorkflowID uint) (*models.Task, error)
	UpdateAuditReportTaskStatus(auditFinalReportID uint, userID uint, auditReportWorkflowID uint) (*models.Task, error)
}

type TaskService struct {
	repo                repositories.TaskRepositoryInteface
	taskStatusSvc       TaskStatusServiceInterface
	moduleSvc           settingService.ModuleServiceInterface
	actionPlanStatusSvc masterServices.ActionPlanStatusServiceInterface
}

func NewTaskService(repo repositories.TaskRepositoryInteface,
	taskStatusSvc TaskStatusServiceInterface,
	moduleSvc settingService.ModuleServiceInterface,
	actionPlanStatusSvc masterServices.ActionPlanStatusServiceInterface,
) TaskServiceInterface {
	return &TaskService{
		repo:                repo,
		taskStatusSvc:       taskStatusSvc,
		moduleSvc:           moduleSvc,
		actionPlanStatusSvc: actionPlanStatusSvc,
	}
}

func (ts *TaskService) GetAll(params *taskDTOs.TaskListQueryParamsDTO) (*taskDTOs.TaskListResponseDTO, error) {
	return ts.repo.GetAll(params)
}

func (ts TaskService) GetByID(id uint) (*models.Task, error) {
	task, err := ts.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return task, nil
}

func (ts TaskService) GetByFieldID(field string, value interface{}) (*models.Task, error) {
	task, err := ts.repo.GetByFieldID(field, value)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, value)
		}
		return nil, err
	}
	return task, nil
}

func (ts *TaskService) Create(taskDTO commonDTOs.TaskFormDTO, currentUserID uint, tx ...utility.DatabaseTransaction) (uint, error) {
	taskStatus, err := ts.taskStatusSvc.GetByType("open")
	if err != nil {
		return 0, err
	}
	referenceNumber, err := ts.repo.GenerateReferenceNumber()
	if err != nil {
		return 0, err
	}
	moduleID, err := ts.moduleSvc.GetModuleIDByIdentifier(taskDTO.ModuleIdentifier)
	if err != nil {
		return 0, err
	}
	task := &models.Task{
		ReferenceNumber:          referenceNumber,
		ModuleID:                 moduleID,
		ComplianceActionPlanID:   taskDTO.ComplianceActionPlanID,
		SLAActionPlanID:          taskDTO.SLAActionPlanID,
		IncidentActionPlanID:     taskDTO.IncidentActionPlanID,
		ControlActionPlanID:      taskDTO.ControlActionPlanID,
		RiskActionPlanID:         taskDTO.RiskActionPlanID,
		MeetingActionPlanID:      taskDTO.MeetingActionPlanID,
		AuditFindingActionPlanID: taskDTO.AuditFindingActionPlanID,
		DocumentID:               taskDTO.DocumentID,
		DocumentWorkflowID:       taskDTO.DocumentWorkflowID,
		AuditFinalReportID:       taskDTO.AuditFinalReportID,
		AuditReportWorkflowID:    taskDTO.AuditReportWorkflowID,
		Title:                    taskDTO.Title,
		Description:              taskDTO.Description,
		TargetDate:               taskDTO.TargetDate,
		URLParam:                 taskDTO.URLParam,
		QueryParam:               taskDTO.QueryParam,
		SeverityLevelID:          taskDTO.SeverityLevelID,
		ResponsibleUserID:        taskDTO.ResponsibleUserID,
		TaskStatusID:             taskStatus.ID,
		CreatedByID:              currentUserID,
	}
	id, err := ts.repo.Create(task, tx...)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (ts *TaskService) Update(currentUserID uint, taskDTO commonDTOs.TaskFormDTO, fieldName string, id uint, tx ...utility.DatabaseTransaction) error {
	existingTask, err := ts.GetByFieldID(fieldName, id)
	if err != nil {
		return err
	}
	existingTask.Title = taskDTO.Title
	existingTask.Description = taskDTO.Description
	existingTask.TargetDate = taskDTO.TargetDate
	existingTask.SeverityLevelID = taskDTO.SeverityLevelID
	existingTask.ResponsibleUserID = taskDTO.ResponsibleUserID
	existingTask.UpdatedByID = currentUserID
	err = ts.repo.Update(existingTask, tx...)
	if err != nil {
		return err
	}
	return nil
}

func (ts *TaskService) Delete(id uint, currentUserID uint) (uint, error) {
	var failedId uint

	entryToDelete, err := ts.GetByID(id)
	if err != nil {
		return 0, err
	}

	// Delete the entry
	if err := ts.repo.Delete(entryToDelete.ID); err != nil {
		failedId = id
		return failedId, err
	}

	return failedId, nil
}

func (ts *TaskService) UpdateStatus(currentUserID uint, taskID uint, actionPlanStatusID uint, isReviewUserNull bool, tx ...utility.DatabaseTransaction) error {
	existingTask, err := ts.GetByID(taskID)
	if err != nil {
		return err
	}

	taskStatus, err := ts.MapActionPlanStatusToTaskStatus(actionPlanStatusID, isReviewUserNull)
	if err != nil {
		return err
	}

	existingTask.TaskStatusID = taskStatus.ID
	existingTask.UpdatedByID = currentUserID

	return ts.repo.Update(existingTask, tx...)
}

func (ts *TaskService) MapActionPlanStatusToTaskStatus(actionPlanStatusID uint, isReviewUserNull bool) (*models.TaskStatus, error) {
	actionPlanStatus, err := ts.actionPlanStatusSvc.GetByID(actionPlanStatusID)
	if err != nil {
		return nil, err
	}

	var taskStatusType string

	if isReviewUserNull && actionPlanStatus.Type == "completed" {
		taskStatusType = "closed"
	} else {
		switch actionPlanStatus.Type {
		case "new":
			taskStatusType = "open"
		case "in_progress":
			taskStatusType = "in_progress"
		case "completed":
			taskStatusType = "in_progress"
		case "on_hold", "terminated", "closed":
			taskStatusType = "closed"
		case "reopen":
			taskStatusType = "in_progress"
		default:
			taskStatusType = "open"
		}
	}

	statusResult, err := ts.taskStatusSvc.GetByType(taskStatusType)
	if err != nil {
		return nil, err
	}

	return statusResult, nil
}

func (ts *TaskService) UpdateTaskWithRelatedIDs(taskID uint, relatedTaskIDsJSON string, tx ...utility.DatabaseTransaction) error {
	return ts.repo.UpdateTaskRelatedIDs(taskID, relatedTaskIDsJSON, tx...)
}

func (ts *TaskService) UpdateDocumentTaskStatus(documentID uint, userID uint, documentWorkflowID uint) (*models.Task, error) {
	taskStatus, err := ts.taskStatusSvc.GetByType("closed")
	if err != nil {
		return nil, err
	}
	task, err := ts.repo.UpdateDocumentTaskStatus(documentID, userID, documentWorkflowID, taskStatus.ID)
	if err != nil {
		return nil, err
	}

	return task, nil
}

func (ts *TaskService) UpdateAuditReportTaskStatus(auditFinalReportID uint, userID uint, auditReportWorkflowID uint) (*models.Task, error) {
	taskStatus, err := ts.taskStatusSvc.GetByType("closed")
	if err != nil {
		return nil, err
	}
	task, err := ts.repo.UpdateAuditReportTaskStatus(auditFinalReportID, userID, auditReportWorkflowID, taskStatus.ID)
	if err != nil {
		return nil, err
	}

	return task, nil
}
