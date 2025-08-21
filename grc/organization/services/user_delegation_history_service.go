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

type UserDelegationHistoryServiceInterface interface {
	GetAll(userID uint) ([]*organizationDTOs.UserDelegationHistoryListDTO, error)
	GetByID(id uint) (*models.UserDelegationHistory, error)
	Create(currentUserID uint, userDelegationHistoryForm *organizationDTOs.UserDelegationHistoryFormDTO) (uint, error)
}

type UserDelegationHistoryService struct {
	repo            organizationRepositories.UserDelegationHistoryRepositoryInterface
	activityLogSvc  logServices.ActivityLogServiceInterface
	userSvc         UserServiceInterface
	notificationSvc commonServices.NotificationServiceInterface
	userSettingSvc  UserSettingServiceInterface
	taskSvc         taskServices.TaskServiceInterface
}

func NewUserDelegationHistoryService(
	repo organizationRepositories.UserDelegationHistoryRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	userSvc UserServiceInterface,
	notificationSvc commonServices.NotificationServiceInterface,
	userSettingSvc UserSettingServiceInterface,
	taskSvc taskServices.TaskServiceInterface,
) UserDelegationHistoryServiceInterface {
	return &UserDelegationHistoryService{
		repo:            repo,
		activityLogSvc:  activityLogSvc,
		userSvc:         userSvc,
		notificationSvc: notificationSvc,
		userSettingSvc:  userSettingSvc,
		taskSvc:         taskSvc,
	}
}

func (ushs *UserDelegationHistoryService) GetAll(userID uint) ([]*organizationDTOs.UserDelegationHistoryListDTO, error) {
	return ushs.repo.GetAll(userID)
}

func (ushs *UserDelegationHistoryService) GetByID(id uint) (*models.UserDelegationHistory, error) {
	process, err := ushs.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return process, nil
}

func (ushs *UserDelegationHistoryService) Create(currentUserID uint, userDelegationHistoryForm *organizationDTOs.UserDelegationHistoryFormDTO) (uint, error) {

	taskCount := len(userDelegationHistoryForm.TaskIDs)
	moduleCount := len(userDelegationHistoryForm.ModuleIDs)

	currentDateStr := time.Now().Format("2006-01-02")
	assignedDate, err := utility.ParseDate(currentDateStr)

	user, err := ushs.userSvc.GetByID(userDelegationHistoryForm.UserID)

	if err != nil {
		return 0, err
	}

	if user.DelegateeUserID != nil {
		return 0, utility.NewForbiddenError("User is already in user delegation stage")
	}

	userDelegation := &models.UserDelegationHistory{
		UserID:          userDelegationHistoryForm.UserID,
		DelegateeUserID: userDelegationHistoryForm.DelegateeUserID,
		Comments:        userDelegationHistoryForm.Comments,
		TaskCount:       taskCount,
		ModuleCount:     moduleCount,
		AssignedDate:    *assignedDate,
		CreatedBYID:     currentUserID,
	}

	oldUser, err := ushs.userSvc.GetByID(userDelegationHistoryForm.UserID)
	if err != nil {
		return 0, err
	}

	tx := ushs.repo.BeginTransaction()

	id, err := ushs.repo.Create(userDelegation, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	oldUser.DelegateeUserID = &userDelegationHistoryForm.DelegateeUserID

	if err := ushs.userSvc.UpdateDelegateeUser(oldUser, tx); err != nil {
		tx.Rollback()
		return 0, err
	}

	for _, task := range userDelegationHistoryForm.TaskIDs {
		_, err = ushs.taskSvc.UpdateDelegateeUser(task, userDelegationHistoryForm.DelegateeUserID, currentUserID, tx)
		if err != nil {
			return 0, err
		}
	}

	tx.Commit()
	go func() {

		err = ushs.activityLogSvc.LogActivity(
			models.CreateAction,
			models.UserDelegationHistory{}.TableName(),
			id,
			userDelegation,
			currentUserID,
			nil,
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
