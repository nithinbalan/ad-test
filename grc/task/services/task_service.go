package services

import (
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/models"
	taskDTOs "jethur/grc/task/dtos"
	"jethur/grc/task/repositories"
	"jethur/utility"
)

type TaskServiceInterface interface {
	GetAll(params *taskDTOs.TaskListQueryParamsDTO) (*taskDTOs.TaskListResponseDTO, error)
	CreateTask(taskDTO commonDTOs.CreateTaskDTO, currentUserID uint, tx ...utility.DatabaseTransaction) error
	GetModuleIDByIdentifier(identifier string) (uint, error)
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

func (ts *TaskService) CreateTask(taskDTO commonDTOs.CreateTaskDTO, currentUserID uint, tx ...utility.DatabaseTransaction) error {
	taskStatus, err := ts.taskStatusSvc.GetByType("open")
	if err != nil {
		return err
	}

	task := &models.Task{
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
		CreatedBy:              currentUserID,
	}
	return ts.repo.CreateTask(task, tx...)

}

func (ts *TaskService) GetModuleIDByIdentifier(identifier string) (uint, error) {
	return ts.repo.GetModuleIDByIdentifier(identifier)
}
