package services

import (
	taskDTOs "jethur/grc/task/dtos"
	taskService "jethur/grc/task/services"

	"jethur/grc/organization/repositories"
)

type UserTaskServiceInterface interface {
	GetAll(userID uint, params *taskDTOs.TaskListQueryParamsDTO) (*taskDTOs.TaskListResponseDTO, error)
	// GetUserTaskSummary(userID uint) ([]taskDTOs.UserTaskStatusCountDTO, error)
}

type UserTaskService struct {
	repo    repositories.UserTaskRepositoryInteface
	taskSvc taskService.TaskServiceInterface
}

func NewUserTaskService(repo repositories.UserTaskRepositoryInteface, taskSvc taskService.TaskServiceInterface) UserTaskServiceInterface {
	return &UserTaskService{repo: repo, taskSvc: taskSvc}
}

func (uts *UserTaskService) GetAll(userID uint, params *taskDTOs.TaskListQueryParamsDTO) (*taskDTOs.TaskListResponseDTO, error) {
	return uts.taskSvc.GetAll(params)
}


// func (uts *UserTaskService) GetUserTaskSummary(userID uint) ([]taskDTOs.UserTaskStatusCountDTO, error) {
// 	taskStatusCounts, err := uts.repo.GetUserTaskCountByStatus(userID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return taskStatusCounts, nil
// }
