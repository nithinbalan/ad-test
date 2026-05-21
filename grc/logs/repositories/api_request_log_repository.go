package repositories

import (
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/logs/dtos"
	"strings"

	"gorm.io/gorm"
)

type APIRequestLogRepositoryInterface interface {
	GetAll(params *dtos.APIRequestLogListQueryParamsDTO) (*dtos.APIRequestLogListResponseDTO, error)
	GetSummary(params *dtos.APIRequestLogListQueryParamsDTO) (*dtos.APIRequestLogSummaryDTO, error)
}

type APIRequestLogRepository struct {
	db *gorm.DB
}

func NewAPIRequestLogRepository(db *gorm.DB) *APIRequestLogRepository {
	return &APIRequestLogRepository{db: db}
}

func (repo *APIRequestLogRepository) GetAll(params *dtos.APIRequestLogListQueryParamsDTO) (*dtos.APIRequestLogListResponseDTO, error) {

	var totalCount int64
	var logs []*dtos.APIRequestLogListDTO

	query := applyAPIRequestLogFilters(repo.db.Table("api_request_logs"), params)

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "desc"
	}
	if params.SortBy == "" {
		params.SortBy = "api_request_logs.requested_at"
	}
	if params.PageSize <= 0 {
		params.PageSize = 10
	}
	if params.Page < 1 {
		params.Page = 1
	}

	query = query.Order(params.SortBy + " " + params.Sort)
	offset := (params.Page - 1) * params.PageSize

	if err := query.Offset(offset).Limit(params.PageSize).Find(&logs).Error; err != nil {
		return nil, err
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))
	return &dtos.APIRequestLogListResponseDTO{
		Items: logs,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}, nil
}

func (repo *APIRequestLogRepository) GetSummary(params *dtos.APIRequestLogListQueryParamsDTO) (*dtos.APIRequestLogSummaryDTO, error) {

	var summary dtos.APIRequestLogSummaryDTO
	query := applyAPIRequestLogFilters(repo.db.Table("api_request_logs"), params)

	err := query.Select(`
		COUNT(*) as total_requests,
		ISNULL(SUM(CASE WHEN log_status = 'completed' THEN 1 ELSE 0 END), 0) as completed,
		ISNULL(SUM(CASE WHEN log_status = 'pending' THEN 1 ELSE 0 END), 0) as pending,
		ISNULL(SUM(CASE WHEN log_status = 'failed' THEN 1 ELSE 0 END), 0) as failed,
		ISNULL(AVG(duration_ms), 0) as avg_duration_ms
	`).Scan(&summary).Error

	if err != nil {
		return nil, err
	}

	return &summary, nil
}

func applyAPIRequestLogFilters(query *gorm.DB, params *dtos.APIRequestLogListQueryParamsDTO) *gorm.DB {
	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where(
			"requested_at BETWEEN ? AND ?",
			params.DateFrom+" 00:00:00",
			params.DateTo+" 23:59:59",
		)
	}

	var (
		statuses    []string
		dateFilter  string
		codeRanges  [][2]int
		slowRequest bool
	)

	for _, filter := range params.QuickFilter {
		switch filter {
		case "completed":
			statuses = append(statuses, "completed")
		case "pending":
			statuses = append(statuses, "pending")
		case "aborted":
			statuses = append(statuses, "aborted")
		case "failed":
			statuses = append(statuses, "failed")
		case "today", "yesterday":
			dateFilter = filter
		case "success_2xx":
			codeRanges = append(codeRanges, [2]int{200, 299})
		case "client_error_4xx":
			codeRanges = append(codeRanges, [2]int{400, 499})
		case "server_error_5xx":
			codeRanges = append(codeRanges, [2]int{500, 599})
		case "slow":
			slowRequest = true
		}
	}

	if len(statuses) > 0 {
		query = query.Where("log_status IN (?)", statuses)
	}

	switch dateFilter {
	case "today":
		query = query.Where(`
			requested_at >= CAST(GETDATE() AS DATE)
			AND requested_at < DATEADD(DAY, 1, CAST(GETDATE() AS DATE))
		`)
	case "yesterday":
		query = query.Where(`
			requested_at >= DATEADD(DAY, -1, CAST(GETDATE() AS DATE))
			AND requested_at < CAST(GETDATE() AS DATE)
		`)
	}

	if len(codeRanges) > 0 {
		clauses := make([]string, 0, len(codeRanges))
		args := make([]any, 0, len(codeRanges)*2)
		for _, r := range codeRanges {
			clauses = append(clauses, "status_code BETWEEN ? AND ?")
			args = append(args, r[0], r[1])
		}
		query = query.Where(strings.Join(clauses, " OR "), args...)
	}

	if slowRequest {
		query = query.Where("duration_ms > ?", 1000)
	}

	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(
			"url LIKE ? OR method LIKE ? OR request_id LIKE ? OR error_message LIKE ?",
			keyword, keyword, keyword, keyword,
		)
	}

	return query
}
