package services

import (
	"jethur/grc/organization/dtos"
	taskDTOs "jethur/grc/task/dtos"
	taskService "jethur/grc/task/services"

	"jethur/grc/organization/repositories"
)

type UserTaskServiceInterface interface {
	GetAll(userID uint, params *taskDTOs.TaskListQueryParamsDTO) (*taskDTOs.TaskListResponseDTO, error)
	GetUserTaskSummary(userID uint) (*dtos.UserTaskSummaryResponseDTO, error)
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


func (uts *UserTaskService) GetUserTaskSummary(userID uint) (*dtos.UserTaskSummaryResponseDTO, error) {
	rawCounts, err := uts.repo.GetUserTaskCountByStatus(userID)
	if err != nil {
		return nil, err
	}

	// Calculate the total count
	totalCount := 0
	for _, count := range rawCounts {
		totalCount += count.Count
	}

	// Ensure that status_counts is initialized as an empty slice if no data is found
	if rawCounts == nil {
		rawCounts = []dtos.UserTaskStatusCountResponseDTO{}
	}

	// Create the response DTO
	response := &dtos.UserTaskSummaryResponseDTO{
		StatusCounts: rawCounts,
		TotalCount:   totalCount,
	}

	return response, nil
}
