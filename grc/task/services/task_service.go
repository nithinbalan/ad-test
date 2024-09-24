package services

import (
	"errors"
	commonDTOs "jethur/grc/common/dtos"
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
	Create(taskDTO commonDTOs.TaskFormDTO, currentUserID uint, tx ...utility.DatabaseTransaction) error
	Update(currentUserID uint, taskDTO commonDTOs.TaskFormDTO, fieldName string, id uint, tx ...utility.DatabaseTransaction) error
	Delete(id uint, currentUserID uint) (uint, error)
}

type TaskService struct {
	repo          repositories.TaskRepositoryInteface
	taskStatusSvc TaskStatusServiceInterface
	moduleSvc     settingService.ModuleServiceInterface
}

func NewTaskService(repo repositories.TaskRepositoryInteface,
	taskStatusSvc TaskStatusServiceInterface,
	moduleSvc settingService.ModuleServiceInterface) TaskServiceInterface {
	return &TaskService{repo: repo,
		taskStatusSvc: taskStatusSvc,
		moduleSvc:     moduleSvc}
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

func (ts *TaskService) Create(taskDTO commonDTOs.TaskFormDTO, currentUserID uint, tx ...utility.DatabaseTransaction) error {
	taskStatus, err := ts.taskStatusSvc.GetByType("open")
	if err != nil {
		return err
	}
	referenceNumber, err := ts.repo.GenerateReferenceNumber()
	if err != nil {
		return err
	}
	moduleID, err := ts.moduleSvc.GetModuleIDByIdentifier(taskDTO.ModuleIdentifier)
	if err != nil {
		return err
	}
	task := &models.Task{
		ReferenceNumber:        referenceNumber,
		ModuleID:               moduleID,
		ComplianceActionPlanID: taskDTO.ComplianceActionPlanID,
		SLAActionPlanID:        taskDTO.SLAActionPlanID,
		IncidentActionPlanID:   taskDTO.IncidentActionPlanID,
		ControlActionPlanID:    taskDTO.ControlActionPlanID,
		RiskActionPlanID:       taskDTO.RiskActionPlanID,
		Title:                  taskDTO.Title,
		Description:            taskDTO.Description,
		TargetDate:             taskDTO.TargetDate,
		URLParam:               taskDTO.URLParam,
		QueryParam:             taskDTO.QueryParam,
		SeverityLevelID:        taskDTO.SeverityLevelID,
		ResponsibleUserID:      taskDTO.ResponsibleUserID,
		TaskStatusID:           taskStatus.ID,
		CreatedByID:            currentUserID,
	}
	return ts.repo.Create(task, tx...)

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
