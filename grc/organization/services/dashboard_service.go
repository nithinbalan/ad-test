package services

import (
	"jethur/grc/models"
	organizationDTOs "jethur/grc/organization/dtos"
	organizationRepositories "jethur/grc/organization/repositories"

)

type DashboardServiceInterface interface {
	GetDashboardSummary(userID uint) (*organizationDTOs.DashboardDTO, error)
	GetUserTasks(userID uint, limit int) ([]*organizationDTOs.UserTaskDTO, error)
	GetOverdueTasks(userID uint, limit int) ([]*organizationDTOs.UserTaskDTO, error)
	GetAssignedRisks(userID uint, limit int) ([]*organizationDTOs.AssignedRiskDTO, error)
	GetRecentActivity(userID uint, limit int) ([]*models.ActivityLog, error)
}

type DashboardService struct {
	repo organizationRepositories.DashboardRepositoryInterface
}

func NewDashboardService(repo organizationRepositories.DashboardRepositoryInterface) *DashboardService {
	return &DashboardService{repo: repo}
}

func (ds *DashboardService) GetDashboardSummary(userID uint) (*organizationDTOs.DashboardDTO, error) {
	return ds.repo.GetDashboardSummary(userID)
}

func (ds *DashboardService) GetUserTasks(userID uint, limit int) ([]*organizationDTOs.UserTaskDTO, error) {
	return ds.repo.GetUserTasks(userID, limit)
}

func (ds *DashboardService) GetOverdueTasks(userID uint, limit int) ([]*organizationDTOs.UserTaskDTO, error) {
	return ds.repo.GetOverdueTasks(userID, limit)
}

func (ds *DashboardService) GetAssignedRisks(userID uint, limit int) ([]*organizationDTOs.AssignedRiskDTO, error) {
	return ds.repo.GetAssignedRisks(userID, limit)
}

func (ds *DashboardService) GetRecentActivity(userID uint, limit int) ([]*models.ActivityLog, error) {
	return ds.repo.GetRecentActivity(userID, limit)
}


