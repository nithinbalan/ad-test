package services

import (
	"jethur/grc/models"
	organizationDTOs "jethur/grc/organization/dtos"
	organizationRepositories "jethur/grc/organization/repositories"
	"math"
)

type MyDashboardServiceInterface interface {
	GetDashboardSummary(userID uint) (*organizationDTOs.MyDashboardDTO, error)
	GetKPIs(userID uint) (*organizationDTOs.MyDashboardKPIDTO, error)
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

func (ds *MyDashboardService) GetKPIs(userID uint) (*organizationDTOs.MyDashboardKPIDTO, error) {
	totalTasks, completedTasks, err := ds.repo.GetTaskCompletionData(userID)
	if err != nil {
		return nil, err
	}
	var taskCompletion float64
	if totalTasks > 0 {
		taskCompletion = float64(completedTasks) / float64(totalTasks) * 100
	} else {
		taskCompletion = 0
	}

	totalActionPlans, closedActionPlans, err := ds.repo.GetRiskMitigationData(userID)
	if err != nil {
		return nil, err
	}
	var riskMitigation float64
	if totalActionPlans > 0 {
		riskMitigation = float64(closedActionPlans) / float64(totalActionPlans) * 100
	} else {
		riskMitigation = 0
	}

	complianceScore, err := ds.repo.GetComplianceScoreData(userID)
	if err != nil {
		return nil, err
	}

	kpis := &organizationDTOs.MyDashboardKPIDTO{
		TaskCompletion:  int(math.Ceil(taskCompletion)),
		RiskMitigation:  int(math.Ceil(riskMitigation)),
		ComplianceScore: int(math.Ceil(complianceScore)),
	}

	return kpis, nil
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


