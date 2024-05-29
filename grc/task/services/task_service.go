package services

import (
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/models"
	"jethur/grc/task/repositories"

	"gorm.io/gorm"
)

type TaskServiceInterface interface {
	CreateTask(taskDTO commonDTOs.CreateTaskDTO, currentUserID uint, tx *gorm.DB) error
	GetModuleIDByIdentifier(identifier string) (uint, error)
}

type TaskService struct {
	repo          repositories.TaskRepositoryInteface
	taskStatusSvc TaskStatusServiceInterface
}

func NewTaskService(repo repositories.TaskRepositoryInteface, taskStatusSvc TaskStatusServiceInterface) TaskServiceInterface {
	return &TaskService{repo: repo, taskStatusSvc: taskStatusSvc}
}

func (ts *TaskService) CreateTask(taskDTO commonDTOs.CreateTaskDTO, currentUserID uint, tx *gorm.DB) error {
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
	return ts.repo.CreateTask(task, tx)

}

func (ts *TaskService) GetModuleIDByIdentifier(identifier string) (uint, error) {
	return ts.repo.GetModuleIDByIdentifier(identifier)
}
