package repositories

import (
	commonDTOs "jethur/grc/common/dtos"
	logDTOs "jethur/grc/logs/dtos"

	"gorm.io/gorm"
)

type TaskSchedularLogRepositoryInterface interface {
	GetAll(params *logDTOs.TaskSchedularLogListQueryParamsDTO) (*logDTOs.TaskSchedularLogListResponseDTO, error)
}

type TaskSchedularLogRepository struct {
	db *gorm.DB
}

func NewTaskSchedularLogRepository(db *gorm.DB) TaskSchedularLogRepositoryInterface {
	return &TaskSchedularLogRepository{
		db: db,
	}
}

func (repo *TaskSchedularLogRepository) GetAll(params *logDTOs.TaskSchedularLogListQueryParamsDTO) (*logDTOs.TaskSchedularLogListResponseDTO, error) {
	var totalCount int64
	var logs []*logDTOs.TaskSchedularLogListDTO

	query := repo.db.Table("task_schedular_logs").
		Select(`
			id,
			task_name,
			status,
			message,
			run_at,
			completed_at as completedAt
		`)

	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where("task_name LIKE ? OR message LIKE ?", keyword, keyword)
	}

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	if params.SortBy == "" {
		params.SortBy = "run_at"
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "desc"
	}
	query = query.Order(params.SortBy + " " + params.Sort)

	offset := (params.Page - 1) * params.PageSize
	err := query.Offset(offset).Limit(params.PageSize).Find(&logs).Error
	if err != nil {
		return nil, err
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	return &logDTOs.TaskSchedularLogListResponseDTO{
		Items: logs,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}, nil
}
