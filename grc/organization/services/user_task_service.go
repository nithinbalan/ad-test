package services

import (
	"jethur/grc/organization/dtos"
	taskDTOs "jethur/grc/task/dtos"
	taskService "jethur/grc/task/services"

	"jethur/grc/organization/repositories"
)

type UserTaskServiceInterface interface {
	GetAll(userID uint, params *taskDTOs.TaskListQueryParamsDTO) (*taskDTOs.TaskListResponseDTO, error)
	GetUserTaskSummary(userID uint) (*dtos.UserTaskStatusCountResponseDTO, error)
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

func (uts *UserTaskService) GetUserTaskSummary(userID uint) (*dtos.UserTaskStatusCountResponseDTO, error) {
	rawCounts, err := uts.repo.GetUserTaskCountByStatus(userID)
	if err != nil {
		return nil, err
	}

	response := &dtos.UserTaskStatusCountResponseDTO{}
	for _, count := range rawCounts {
		switch count.Type {
		case "open":
			response.Open = count.Count
		case "in_progress":
			response.InProgress = count.Count
		case "closed":
			response.Closed = count.Count
		}
		response.Total += count.Count
	}

	return response, nil
}
