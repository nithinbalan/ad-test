package repositories

import (
	"fmt"
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/compliance/dtos"
	"jethur/grc/models"
	"jethur/utility"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type ConsultationRequestRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(params *dtos.ConsultationRequestListQueryParamsDTO) (*dtos.ConsultationRequestListResponseDTO, error)
	GetByID(id uint) (*models.ConsultationRequest, error)
	GetFullDetailsByID(id uint) (*models.ConsultationRequest, error)
	Create(cr *models.ConsultationRequest, tx ...utility.DatabaseTransaction) (uint, error)
	Update(cr *models.ConsultationRequest, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error
	GenerateReferenceNumber() (string, error)
}

type ConsultationRequestRepository struct {
	db *gorm.DB
}

func NewConsultationRequestRepository(db *gorm.DB) *ConsultationRequestRepository {
	return &ConsultationRequestRepository{db: db}
}

func (repo *ConsultationRequestRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

// GetAll with filters, pagination, keyword search
func (repo *ConsultationRequestRepository) GetAll(params *dtos.ConsultationRequestListQueryParamsDTO) (*dtos.ConsultationRequestListResponseDTO, error) {
	var requests []*dtos.ConsultationRequestListDTO
	var totalCount int64
	query := repo.db.Table("consultation_requests").
		Select(`
		consultation_requests.*,
		consultation_request_statuses.title_en AS status_title_en,
		consultation_request_statuses.title_ar AS status_title_ar,
	
		req_unit.title_en AS requesting_department_unit_title_en,
		req_unit.title_ar AS requesting_department_unit_title_ar,
		con_unit.title_en AS consulted_department_unit_title_en,
		con_unit.title_ar AS consulted_department_unit_title_ar,
		resp_user.first_name_en AS responsible_first_name_en,
		resp_user.last_name_en  AS responsible_last_name_en,
		resp_user.email         AS responsible_email,
		created_user.first_name_en AS created_by_first_name_en,
		created_user.last_name_en  AS created_by_last_name_en,
		updated_user.first_name_en AS updated_by_first_name_en,
		updated_user.last_name_en  AS updated_by_last_name_en
	`).
		Joins(`LEFT JOIN consultation_request_statuses ON consultation_request_statuses.id = consultation_requests.consultation_request_status_id`).
		Joins(`LEFT JOIN organization_units AS req_unit ON req_unit.id = consultation_requests.requesting_department_unit_id`).
		Joins(`LEFT JOIN organization_units AS con_unit ON con_unit.id = consultation_requests.consulted_department_unit_id`).
		Joins(`LEFT JOIN users AS resp_user ON resp_user.id = consultation_requests.responsible_user_id`).
		Joins(`LEFT JOIN users AS created_user ON created_user.id = consultation_requests.created_by_id`).
		Joins(`LEFT JOIN users AS updated_user ON updated_user.id = consultation_requests.updated_by_id`)

	if params.Keyword != "" {
		kw := "%" + params.Keyword + "%"
		query = query.Where(`
			consultation_requests.reference_number LIKE ? OR
			consultation_requests.request_details LIKE ? `,
			kw, kw,
		)
	}

	// Sorting
	if params.SortBy == "" {
		params.SortBy = "consultation_requests.date_raised"
	}
	if params.Sort == "" {
		params.Sort = "DESC"
	}
	query = query.Order(params.SortBy + " " + params.Sort)

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	offset := (params.Page - 1) * params.PageSize
	query = query.Offset(offset).Limit(params.PageSize)

	if err := query.Find(&requests).Error; err != nil {
		return nil, err
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	response := &dtos.ConsultationRequestListResponseDTO{
		Items: requests,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}
	return response, nil
}

// GetByID
func (repo *ConsultationRequestRepository) GetByID(id uint) (*models.ConsultationRequest, error) {
	var cr models.ConsultationRequest
	err := repo.db.First(&cr, id).Error
	if err != nil {
		return nil, err
	}
	return &cr, nil
}

// GetFullDetailsByID with all relationships preloaded
func (repo *ConsultationRequestRepository) GetFullDetailsByID(id uint) (*models.ConsultationRequest, error) {

	db := repo.db
	var cr models.ConsultationRequest
	err := db.
		Preload("ResponsibleUser").
		Preload("ConsultedDepartmentUnit").
		Preload("RequestingDepartmentUnit").
		Preload("ConsultationRequestStatus").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		First(&cr, id).Error

	if err != nil {
		return nil, err
	}

	return &cr, nil
}

// Create
func (repo *ConsultationRequestRepository) Create(cr *models.ConsultationRequest, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(cr).Error; err != nil {
		return 0, err
	}
	return cr.ID, nil
}

// Update
func (repo *ConsultationRequestRepository) Update(cr *models.ConsultationRequest, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Save(cr).Error
}

// Delete
func (repo *ConsultationRequestRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var consultationRequest *models.ConsultationRequest

	// Check if the record exists first
	if err := db.First(&consultationRequest, id).Error; err != nil {
		return err
	}

	result := db.Delete(&models.ConsultationRequest{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

// GenerateReferenceNumber in format CR-YYYY-0001
func (repo *ConsultationRequestRepository) GenerateReferenceNumber() (string, error) {
	var lastReferenceNumber string
	err := repo.db.Model(&models.ConsultationRequest{}).
		Select("reference_number").
		Order("reference_number desc").
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
	lastID, err := strconv.ParseUint(lastDigits, 10, 64)
	if err != nil {
		return "", err
	}

	newReferenceNumber := fmt.Sprintf("CONS-%d-%04d", currentYear, lastID+1)
	return newReferenceNumber, nil
}
