package services

import (
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/models"
	"jethur/grc/task/repositories"

	"gorm.io/gorm"
)

type TaskServiceInterface interface {
	CreateTask(taskFormDTO commonDTOs.TaskFormDTO,currentUserID uint, tx *gorm.DB) error
	GetModuleIDByIdentifier(identifier string) (uint, error)
}

type TaskService struct {
	repo repositories.TaskRepositoryInteface
	taskStatusSvc TaskStatusServiceInterface

}

func NewTaskService(repo repositories.TaskRepositoryInteface,taskStatusSvc TaskStatusServiceInterface) TaskServiceInterface {
	return &TaskService{repo: repo,taskStatusSvc:	taskStatusSvc}
}

func (ts *TaskService) CreateTask(taskFormDTO commonDTOs.TaskFormDTO,currentUserID uint, tx *gorm.DB) error {
	taskStatus, err := ts.taskStatusSvc.GetByType("new")
	if err != nil {
		return err
	}

	task := &models.Task{
        ModuleID:                taskFormDTO.ModuleID,
        ComplianceActionPlanID:  taskFormDTO.ComplianceActionPlanID,
        SLAActionPlanID:         taskFormDTO.SLAActionPlanID,
        IncidentActionPlanID:    taskFormDTO.IncidentActionPlanID,
        ControlActionPlanID:     taskFormDTO.ControlActionPlanID,
        Title:                   taskFormDTO.Title,
        Description:             taskFormDTO.Description,
        TargetDate:              taskFormDTO.TargetDate,
        URL:                     taskFormDTO.URL,
        SeverityLevelID:         taskFormDTO.SeverityLevelID,
        ReviewResponsibleUserID: taskFormDTO.ReviewResponsibleUserID,
        TaskStatusID:            taskStatus.ID,
        CreatedBy:               currentUserID,
    }
	if tx != nil {
		return ts.repo.CreateTask(task, tx)
	} else {
		return ts.repo.CreateTask(task)
	}
}

func (ts *TaskService) GetModuleIDByIdentifier(identifier string) (uint, error) {
	return ts.repo.GetModuleIDByIdentifier(identifier)
}