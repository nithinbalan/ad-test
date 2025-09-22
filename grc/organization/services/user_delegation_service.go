package services

import (
	"errors"
	commonServices "jethur/grc/common/services"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	organizationDTOs "jethur/grc/organization/dtos"
	organizationRepositories "jethur/grc/organization/repositories"
	taskServices "jethur/grc/task/services"
	"jethur/utility"
	"time"

	"gorm.io/gorm"
)

type UserDelegationServiceInterface interface {
	GetAll(userID uint) ([]*organizationDTOs.UserDelegationHistoryListDTO, error)
	GetByID(id uint) (*models.UserDelegationHistory, error)
	Create(currentUserID uint, userDelegationForm *organizationDTOs.UserDelegationFormDTO) (uint, error)

	GetLatestUserDelegation(userID uint) (*organizationDTOs.UserDelegationHistoryListDTO, error)
	EndUserDelegation(userID uint, currentUserID uint) error
}

type UserDelegationService struct {
	repo            organizationRepositories.UserDelegationRepositoryInterface
	activityLogSvc  logServices.ActivityLogServiceInterface
	userSvc         UserServiceInterface
	notificationSvc commonServices.NotificationServiceInterface
	userSettingSvc  UserSettingServiceInterface
	taskSvc         taskServices.TaskServiceInterface
}

func NewUserDelegationService(
	repo organizationRepositories.UserDelegationRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	userSvc UserServiceInterface,
	notificationSvc commonServices.NotificationServiceInterface,
	userSettingSvc UserSettingServiceInterface,
	taskSvc taskServices.TaskServiceInterface,
) UserDelegationServiceInterface {
	return &UserDelegationService{
		repo:            repo,
		activityLogSvc:  activityLogSvc,
		userSvc:         userSvc,
		notificationSvc: notificationSvc,
		userSettingSvc:  userSettingSvc,
		taskSvc:         taskSvc,
	}
}

func (ushs *UserDelegationService) GetAll(userID uint) ([]*organizationDTOs.UserDelegationHistoryListDTO, error) {
	return ushs.repo.GetAll(userID)
}

func (ushs *UserDelegationService) GetByID(id uint) (*models.UserDelegationHistory, error) {
	process, err := ushs.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return process, nil
}

func (ushs *UserDelegationService) Create(currentUserID uint, userDelegationForm *organizationDTOs.UserDelegationFormDTO) (uint, error) {

	taskCount := len(userDelegationForm.TaskIDs)
	moduleCount, err := ushs.taskSvc.GetModuleCountByTaskIDs(userDelegationForm.TaskIDs)
	if err != nil {
		return 0, err
	}

	currentDateStr := time.Now().Format("2006-01-02")
	assignedDate, err := utility.ParseDate(currentDateStr)

	user, err := ushs.userSvc.GetByID(userDelegationForm.UserID)

	if err != nil {
		return 0, err
	}

	if user.DelegationUserID != nil {
		return 0, utility.NewForbiddenError("User is already in user delegation stage")
	}

	userDelegation := &models.UserDelegationHistory{
		UserID:          userDelegationForm.UserID,
		DelegateeUserID: userDelegationForm.DelegateeUserID,
		Comments:        userDelegationForm.Comments,
		TaskCount:       taskCount,
		ModuleCount:     moduleCount,
		AssignedDate:    *assignedDate,
		CreatedBYID:     currentUserID,
	}

	oldUser, err := ushs.userSvc.GetByID(userDelegationForm.UserID)
	if err != nil {
		return 0, err
	}

	tx := ushs.repo.BeginTransaction()

	id, err := ushs.repo.Create(userDelegation, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	oldUser.DelegationUserID = &userDelegationForm.DelegateeUserID

	if err := ushs.userSvc.UpdateDelegateeUser(oldUser, tx); err != nil {
		tx.Rollback()
		return 0, err
	}

	for _, task := range userDelegationForm.TaskIDs {
		_, err = ushs.taskSvc.UpdateDelegateeUser(task, &userDelegationForm.DelegateeUserID, currentUserID, tx)
		if err != nil {
			return 0, err
		}
	}

	tables := []string{
		"threat_action_plans",
		"vulnerability_action_plans",
		"audit_finding_action_plans",
		"test_and_exercise_recovery_action_plans",
		"compliance_action_plans",
		"compliance_assessment_requirement_action_plans",
		"sla_action_plans",
		"incident_action_plans",
		"meeting_action_plans",
		"control_action_plans",
		"risk_action_plans",
		"strategy_profile_action_plans",
		"stakeholder_action_plans",
		"project_action_plans",
		"document_change_request_workflow_participants",
		"document_workflow_participants",
	}

	if err := ushs.repo.UpdateActionPlanDelegationUser(userDelegationForm.UserID, tables, &userDelegationForm.DelegateeUserID, tx); err != nil {
		tx.Rollback()
		return 0, err
	}

	tx.Commit()
	go func() {

		err = ushs.activityLogSvc.LogActivity(
			models.CreateAction,
			models.UserDelegationHistory{}.TableName(),
			id,
			userDelegation,
			currentUserID,
			"", "", "", nil,
		)
		if err != nil {
			utility.LogError("Failed to log activity", map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.UserDelegationHistory{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			})
		}
	}()

	return id, nil
}

func (ushs *UserDelegationService) GetLatestUserDelegation(userID uint) (*organizationDTOs.UserDelegationHistoryListDTO, error) {
	_, err := ushs.userSvc.GetByID(userID)
	if err != nil {
		return nil, err
	}

	return ushs.repo.GetLatestUserDelegation(userID)
}

func (ushs *UserDelegationService) EndUserDelegation(userID uint, currentUserID uint) error {

	responsibleUser, err := ushs.userSvc.GetByID(userID)
	if err != nil {
		return err
	}
	tx := ushs.repo.BeginTransaction()

	tasks, _ := ushs.taskSvc.GetAllTasksByResponsibleUserID(userID, *responsibleUser.DelegationUserID)

	for _, task := range tasks {
		_, err = ushs.taskSvc.UpdateDelegateeUser(task.ID, nil, currentUserID, tx)
		if err != nil {
			return err
		}
	}

	responsibleUser.DelegationUserID = nil

	if err := ushs.userSvc.UpdateDelegateeUser(responsibleUser, tx); err != nil {
		tx.Rollback()
		return err
	}

	tables := []string{
		"threat_action_plans",
		"vulnerability_action_plans",
		"audit_finding_action_plans",
		"test_and_exercise_recovery_action_plans",
		"compliance_action_plans",
		"compliance_assessment_requirement_action_plans",
		"sla_action_plans",
		"incident_action_plans",
		"meeting_action_plans",
		"control_action_plans",
		"risk_action_plans",
		"strategy_profile_action_plans",
		"stakeholder_action_plans",
		"project_action_plans",
		"document_change_request_workflow_participants",
		"document_workflow_participants",
	}

	if err := ushs.repo.UpdateActionPlanDelegationUser(userID, tables, nil, tx); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}
