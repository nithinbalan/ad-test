package services

import (
	"jethur/grc/models"
	organizationDTOs "jethur/grc/organization/dtos"
	organizationRepositories "jethur/grc/organization/repositories"
)

type MyDashboardServiceInterface interface {
	GetDashboardSummary(userID uint) (*organizationDTOs.MyDashboardDTO, error)
	GetUserTasks(userID uint, limit int) ([]organizationDTOs.UserTaskDTO, error)
	GetOverdueTasks(userID uint, limit int) ([]organizationDTOs.UserTaskDTO, error)
	GetAssignedRisks(userID uint, limit int) ([]organizationDTOs.AssignedRiskDTO, error)
	GetRecentActivity(userID uint, limit int) ([]*models.ActivityLog, error)
}

type MyDashboardService struct {
	repo organizationRepositories.MyDashboardRepositoryInterface
}

func NewMyDashboardService(repo organizationRepositories.MyDashboardRepositoryInterface) *MyDashboardService {
	return &MyDashboardService{repo: repo}
}

func (ds *MyDashboardService) GetDashboardSummary(userID uint) (*organizationDTOs.MyDashboardDTO, error) {
	tasksAssigned, tasksCompleted, overdueItems, newNotifications, err := ds.repo.GetDashboardSummary(userID)
	if err != nil {
		return nil, err
	}

	summary := &organizationDTOs.MyDashboardDTO{
		TasksAssigned:    tasksAssigned,
		TasksCompleted:   tasksCompleted,
		OverdueItems:     overdueItems,
		NewNotifications: newNotifications,
	}

	return summary, nil
}

func (ds *MyDashboardService) GetUserTasks(userID uint, limit int) ([]organizationDTOs.UserTaskDTO, error) {
	return ds.repo.GetUserTasks(userID, limit)
}

func (ds *MyDashboardService) GetOverdueTasks(userID uint, limit int) ([]organizationDTOs.UserTaskDTO, error) {
	return ds.repo.GetOverdueTasks(userID, limit)
}

func (ds *MyDashboardService) GetAssignedRisks(userID uint, limit int) ([]organizationDTOs.AssignedRiskDTO, error) {
	return ds.repo.GetAssignedRisks(userID, limit)
}


func (ds *MyDashboardService) GetRecentActivity(userID uint, limit int) ([]*models.ActivityLog, error) {
	return ds.repo.GetRecentActivity(userID, limit)
}


