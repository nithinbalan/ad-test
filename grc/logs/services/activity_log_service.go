package services

import (
	"encoding/json"
	"errors"
	logDTOs "jethur/grc/logs/dtos"
	"jethur/grc/logs/repositories"
	"jethur/grc/models"
	"time"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type ActivityLogServiceInterface interface {
	GetAll(params *logDTOs.ActivityLogListQueryParamsDTO) (*logDTOs.ActivityLogListResponseDTO, error)
	GetFullDetailsByID(id uint) (*logDTOs.ActivityLogDetailDTO, error)
	LogActivity(action models.Action, tableName string, recordID uint, changedData interface{}, userID uint, tx *gorm.DB) error
	LogBulkActivities(action models.Action, tableName string, recordIDs []uint, changedData interface{}, userID uint, tx *gorm.DB) error // Simplified signature

}

type ActivityLogService struct {
	repo repositories.ActivityLogRepositoryInteface
}

func NewActivityLogService(repo repositories.ActivityLogRepositoryInteface) ActivityLogServiceInterface {
	return &ActivityLogService{repo: repo}
}

func (als *ActivityLogService) GetAll(params *logDTOs.ActivityLogListQueryParamsDTO) (*logDTOs.ActivityLogListResponseDTO, error) {
	return als.repo.GetAll(params)
}

func (als *ActivityLogService) GetFullDetailsByID(id uint) (*logDTOs.ActivityLogDetailDTO, error) {
	activityLog, err := als.repo.GetFullDetailsByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, err
	}

	var activityLogDTO logDTOs.ActivityLogDetailDTO
	err = mapstructure.Decode(activityLog, &activityLogDTO)
	if err != nil {
		return nil, err
	}

	return &activityLogDTO, nil
}

func (als *ActivityLogService) LogActivity(action models.Action, tableName string, recordID uint, changedData interface{}, userID uint, tx *gorm.DB) error {

	changedDataJSON, err := json.Marshal(changedData)
	if err != nil {
		return err
	}

	activityLog := models.ActivityLog{
		Action:      action,
		TableName:   tableName,
		RecordID:    recordID,
		ChangedData: string(changedDataJSON),
		UserID:      userID,
		CreatedAt:   time.Now(),
	}
	if tx != nil {
		return als.repo.SaveActivityLogs([]*models.ActivityLog{&activityLog}, tx)
	} else {
		return als.repo.SaveActivityLogs([]*models.ActivityLog{&activityLog})
	}
}

func (als *ActivityLogService) LogBulkActivities(action models.Action, tableName string, recordIDs []uint, changedData interface{}, userID uint, tx *gorm.DB) error {
	var logs []*models.ActivityLog

	changedDataJSON, err := json.Marshal(changedData)
	if err != nil {
		return err
	}

	for _, recordID := range recordIDs {
		log := models.ActivityLog{
			Action:      action,
			TableName:   tableName,
			RecordID:    recordID,
			ChangedData: string(changedDataJSON),
			UserID:      userID,
			CreatedAt:   time.Now(),
		}

		logs = append(logs, &log)
	}

	if tx != nil {
		return als.repo.SaveActivityLogs(logs, tx)
	} else {
		return als.repo.SaveActivityLogs(logs)
	}
}
