package services

import (
	"errors"
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/models"
	taskDTOs "jethur/grc/task/dtos"
	"jethur/grc/task/repositories"
	"jethur/utility"

	"gorm.io/gorm"
)

type TaskServiceInterface interface {
	GetAll(params *taskDTOs.TaskListQueryParamsDTO) (*taskDTOs.TaskListResponseDTO, error)
	GetByID(id uint) (*models.Task, error)
	GetByFieldID(field string, value interface{}) (*models.Task, error)
	CreateTask(taskDTO commonDTOs.CreateTaskDTO, currentUserID uint, tx ...utility.DatabaseTransaction) error
	UpdateTask(currentUserID uint, task *models.Task, tx ...utility.DatabaseTransaction) error
	DeleteTask(id uint, currentUserID uint) (uint, error)
}

type TaskService struct {
	repo          repositories.TaskRepositoryInteface
	taskStatusSvc TaskStatusServiceInterface
}

func NewTaskService(repo repositories.TaskRepositoryInteface, taskStatusSvc TaskStatusServiceInterface) TaskServiceInterface {
	return &TaskService{repo: repo, taskStatusSvc: taskStatusSvc}
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
	task, err := ts.repo.GetByFieldID(field,value)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, value)
		}
		return nil, err
	}
	return task, nil
}

func (ts *TaskService) CreateTask(taskDTO commonDTOs.CreateTaskDTO, currentUserID uint, tx ...utility.DatabaseTransaction) error {
	taskStatus, err := ts.taskStatusSvc.GetByType("open")
	if err != nil {
		return err
	}
	referenceNumber, err := ts.repo.GenerateReferenceNumber()
	if err != nil {
		return err
	}
	task := &models.Task{
		ReferenceNumber:        referenceNumber,
		ModuleID:               taskDTO.ModuleID,
		ComplianceActionPlanID: taskDTO.ComplianceActionPlanID,
		SLAActionPlanID:        taskDTO.SLAActionPlanID,
		IncidentActionPlanID:   taskDTO.IncidentActionPlanID,
		ControlActionPlanID:    taskDTO.ControlActionPlanID,
		Title:                  taskDTO.Title,
		Description:            taskDTO.Description,
		TargetDate:             taskDTO.TargetDate,
		URL:                    taskDTO.URL,
		SeverityLevelID:        taskDTO.SeverityLevelID,
		ResponsibleUserID:      taskDTO.ResponsibleUserID,
		TaskStatusID:           taskStatus.ID,
		CreatedByID:            currentUserID,
	}
	return ts.repo.CreateTask(task, tx...)

}

func (ts *TaskService) UpdateTask(currentUserID uint, task *models.Task, tx ...utility.DatabaseTransaction) error {
	err := ts.repo.Update(task, tx...)

	if err != nil {
		return err
	}

	return nil
}

func (ts *TaskService) DeleteTask(id uint, currentUserID uint) (uint, error) {
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
