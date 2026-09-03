package repositories

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"time"

	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/models"
	"jethur/grc/project/dtos"
	"jethur/utility"

	"gorm.io/gorm"
)

type ProjectDecisionLogRepositoryInterface interface {
	GetAll(projectID uint, params *dtos.ProjectDecisionLogListQueryParamsDTO) (*dtos.ProjectDecisionLogListResponseDTO, error)
	GetByID(projectID uint, id uint) (*models.ProjectDecisionLog, error)
	GetFullDetailsByID(projectID uint, id uint) (*models.ProjectDecisionLog, error)
	GenerateReferenceNumber() (string, error)
	Create(decisionLog *models.ProjectDecisionLog, tx ...utility.DatabaseTransaction) (uint, error)
	Update(decisionLog *models.ProjectDecisionLog, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error
	CreateDocuments(documents []*models.ProjectDecisionLogDocument, tx ...utility.DatabaseTransaction) error
	DeleteDocuments(ids []*uint, tx ...utility.DatabaseTransaction) error
	DeleteDocumentsByDecisionLogID(decisionLogID uint, tx ...utility.DatabaseTransaction) error
	GetDocumentByID(id uint) (*models.ProjectDecisionLogDocument, error)
	BeginTransaction() utility.DatabaseTransaction
}

type ProjectDecisionLogRepository struct {
	db *gorm.DB
}

func NewProjectDecisionLogRepository(db *gorm.DB) *ProjectDecisionLogRepository {
	return &ProjectDecisionLogRepository{db: db}
}

func (pdlr *ProjectDecisionLogRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(pdlr.db)
}

func (pdlr *ProjectDecisionLogRepository) GetAll(projectID uint, params *dtos.ProjectDecisionLogListQueryParamsDTO) (*dtos.ProjectDecisionLogListResponseDTO, error) {
	var decisionLogs []*dtos.ProjectDecisionLogListDTO
	var totalCount int64

	query := pdlr.db.Table("project_decision_logs").
		Select(`
			project_decision_logs.id,
			project_decision_logs.reference_number,
			project_decision_logs.project_id,
			project_decision_logs.decision,
			project_decision_logs.decision_notes,
			project_decision_logs.reason,
			CONVERT(VARCHAR, project_decision_logs.decision_date, 23) AS decision_date,
			project_decision_logs.requested_by_user_id,
			requested_user.first_name_en AS requested_by_user_first_name_en,
			requested_user.first_name_ar AS requested_by_user_first_name_ar,
			requested_user.last_name_en AS requested_by_user_last_name_en,
			requested_user.last_name_ar AS requested_by_user_last_name_ar,
			requested_user.initial AS requested_by_user_initial,
			requested_user.color AS requested_by_user_color,
			requested_user.email AS requested_by_user_email,
			requested_user.profile_image_path AS requested_by_user_profile_image_path,
			created_user.id AS created_user_id,
			created_user.first_name_en AS created_user_first_name_en,
			created_user.first_name_ar AS created_user_first_name_ar,
			created_user.last_name_en AS created_user_last_name_en,
			created_user.last_name_ar AS created_user_last_name_ar,
			created_user.initial AS created_user_initial,
			created_user.color AS created_user_color,
			created_user.email AS created_user_email,
			created_user.profile_image_path AS created_user_profile_image_path,
			updated_user.id AS updated_user_id,
			updated_user.first_name_en AS updated_user_first_name_en,
			updated_user.first_name_ar AS updated_user_first_name_ar,
			updated_user.last_name_en AS updated_user_last_name_en,
			updated_user.last_name_ar AS updated_user_last_name_ar,
			updated_user.initial AS updated_user_initial,
			updated_user.color AS updated_user_color,
			updated_user.email AS updated_user_email,
			updated_user.profile_image_path AS updated_user_profile_image_path,
			CONVERT(VARCHAR, project_decision_logs.created_at, 120) AS created_at,
			CONVERT(VARCHAR, project_decision_logs.updated_at, 120) AS updated_at
		`).
		Joins(`LEFT JOIN users AS requested_user ON project_decision_logs.requested_by_user_id = requested_user.id`).
		Joins(`LEFT JOIN users AS created_user ON project_decision_logs.created_by_id = created_user.id`).
		Joins(`LEFT JOIN users AS updated_user ON project_decision_logs.updated_by_id = updated_user.id`)

	if projectID != 0 {
		query = query.Where("project_decision_logs.project_id = ?", projectID)
	}

	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where("project_decision_logs."+columnName+" IN (?)", values)
		}
	}

	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("project_decision_logs.decision_date BETWEEN ? AND ?", params.DateFrom, params.DateTo)
	}

	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
			(project_decision_logs.reference_number LIKE ? OR
			project_decision_logs.decision LIKE ? OR
			project_decision_logs.decision_notes LIKE ? OR
			project_decision_logs.reason LIKE ?)
		`, keyword, keyword, keyword, keyword)
	}

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	page := params.Page
	if page <= 0 {
		page = 1
	}
	pageSize := params.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}

	sort := params.Sort
	if sort == "" {
		sort = "DESC"
	}
	sortBy := params.SortBy
	if sortBy == "" {
		sortBy = "project_decision_logs.id"
	} else {
		sortBy = "project_decision_logs." + sortBy
	}

	offset := (page - 1) * pageSize
	if err := query.Order(fmt.Sprintf("%s %s", sortBy, sort)).
		Offset(offset).
		Limit(pageSize).
		Scan(&decisionLogs).Error; err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(totalCount) / float64(pageSize)))
	return &dtos.ProjectDecisionLogListResponseDTO{
		Data: decisionLogs,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: page,
			PageSize:    pageSize,
			TotalCount:  totalCount,
			TotalPages:  totalPages,
		},
	}, nil
}

func (pdlr *ProjectDecisionLogRepository) GetByID(projectID uint, id uint) (*models.ProjectDecisionLog, error) {
	var decisionLog models.ProjectDecisionLog
	err := pdlr.db.Where("project_id = ? AND id = ?", projectID, id).First(&decisionLog).Error
	if err != nil {
		return nil, err
	}
	return &decisionLog, nil
}

func (pdlr *ProjectDecisionLogRepository) GetFullDetailsByID(projectID uint, id uint) (*models.ProjectDecisionLog, error) {
	var decisionLog models.ProjectDecisionLog
	err := pdlr.db.Where("project_id = ? AND id = ?", projectID, id).
		Preload("Project").
		Preload("RequestedByUser").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		Preload("Documents").
		First(&decisionLog).Error
	if err != nil {
		return nil, err
	}
	return &decisionLog, nil
}

func (pdlr *ProjectDecisionLogRepository) GenerateReferenceNumber() (string, error) {
	var lastReferenceNumber string
	err := pdlr.db.Model(&models.ProjectDecisionLog{}).
		Select("reference_number").
		Where("reference_number IS NOT NULL").
		Order("reference_number DESC").
		Limit(1).
		Pluck("reference_number", &lastReferenceNumber).
		Error

	if err != nil {
		return "", err
	}

	lastDigits := "0000"
	if lastReferenceNumber != "" && len(lastReferenceNumber) >= 4 {
		lastDigits = lastReferenceNumber[len(lastReferenceNumber)-4:]
	}

	currentYear := time.Now().Year()
	lastReferenceID, err := strconv.ParseUint(lastDigits, 10, 64)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("PDL-%d-%04d", currentYear, lastReferenceID+1), nil
}

func (pdlr *ProjectDecisionLogRepository) Create(decisionLog *models.ProjectDecisionLog, tx ...utility.DatabaseTransaction) (uint, error) {
	db := pdlr.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(decisionLog).Error; err != nil {
		return 0, err
	}
	return decisionLog.ID, nil
}

func (pdlr *ProjectDecisionLogRepository) Update(decisionLog *models.ProjectDecisionLog, tx ...utility.DatabaseTransaction) error {
	db := pdlr.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Save(decisionLog).Error
}

func (pdlr *ProjectDecisionLogRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := pdlr.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Delete(&models.ProjectDecisionLog{}, id).Error
}

func (pdlr *ProjectDecisionLogRepository) CreateDocuments(documents []*models.ProjectDecisionLogDocument, tx ...utility.DatabaseTransaction) error {
	db := pdlr.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Create(&documents).Error
}

func (pdlr *ProjectDecisionLogRepository) DeleteDocuments(ids []*uint, tx ...utility.DatabaseTransaction) error {
	db := pdlr.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Where("id IN (?)", ids).Delete(&models.ProjectDecisionLogDocument{}).Error
}

func (pdlr *ProjectDecisionLogRepository) DeleteDocumentsByDecisionLogID(decisionLogID uint, tx ...utility.DatabaseTransaction) error {
	db := pdlr.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Where("project_decision_log_id = ?", decisionLogID).Delete(&models.ProjectDecisionLogDocument{}).Error
}

func (pdlr *ProjectDecisionLogRepository) GetDocumentByID(id uint) (*models.ProjectDecisionLogDocument, error) {
	var document models.ProjectDecisionLogDocument
	if err := pdlr.db.Where("id = ?", id).First(&document).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &utility.ResourceNotFoundError{}
		}
		return nil, err
	}
	return &document, nil
}
