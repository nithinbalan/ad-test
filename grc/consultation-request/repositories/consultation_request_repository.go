package repositories

import (
	"fmt"
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/consultation-request/dtos"

	"jethur/grc/models"
	"jethur/utility"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type ConsultationRequestRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(params *dtos.ConsultationRequestListQueryParamsDTO, departmentID *uint, currentUserID uint) (*dtos.ConsultationRequestListResponseDTO, error)
	GetByID(id uint) (*models.ConsultationRequest, error)
	GetFullDetailsByID(id uint) (*models.ConsultationRequest, error)
	Create(cr *models.ConsultationRequest, tx ...utility.DatabaseTransaction) (uint, error)
	Update(cr *models.ConsultationRequest, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error
	GenerateReferenceNumber() (string, error)

	///
	CreateComment(comment *models.ConsultationRequestComment, tx ...utility.DatabaseTransaction) (uint, error)
	GetAllConsultationRequestComments(consultationRequestID uint) ([]*dtos.ConsultationRequestCommentListDTO, error)
	GetCommentByID(id uint) (*models.ConsultationRequestComment, error)
	UpdateComment(comment *models.ConsultationRequestComment, tx ...utility.DatabaseTransaction) error
	DeleteComment(id uint, tx ...utility.DatabaseTransaction) error
	CreateDocuments(documents []*models.ConsultationRequestCommentDocument, tx ...utility.DatabaseTransaction) error
	DeleteDocuments(documentIDs []*uint, tx ...utility.DatabaseTransaction) error
	GetCommentDocumentByID(id uint) (*models.ConsultationRequestCommentDocument, error)
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
func (repo *ConsultationRequestRepository) GetAll(params *dtos.ConsultationRequestListQueryParamsDTO, departmentID *uint, currentUserID uint) (*dtos.ConsultationRequestListResponseDTO, error) {
	var requests []*dtos.ConsultationRequestListDTO
	var totalCount int64
	query := repo.db.Table("consultation_requests").
		Select(`
		consultation_requests.*,
		  
		CONVERT(VARCHAR, consultation_requests.date_raised, 23) AS date_raised,
		CONVERT(VARCHAR, consultation_requests.reply_date, 23) AS reply_date,
		consultation_request_statuses.title_en AS status_title_en,
		consultation_request_statuses.title_ar AS status_title_ar,
	
		req_unit.title_en AS requesting_department_unit_title_en,
		req_unit.title_ar AS requesting_department_unit_title_ar,
		con_unit.title_en AS consulted_department_unit_title_en,
		con_unit.title_ar AS consulted_department_unit_title_ar,

		resp_user.id           AS responsible_user_id,
		resp_user.first_name_en AS responsible_first_name_en,
		resp_user.first_name_en AS responsible_first_name_en,
		resp_user.last_name_en  AS responsible_last_name_en,
		resp_user.first_name_ar AS responsible_first_name_ar,   -- Arabic first name
		resp_user.last_name_ar  AS responsible_last_name_ar,    -- Arabic last name
		resp_user.initial       AS responsible_initial,        -- initials
		resp_user.email         AS responsible_email,
		resp_user.color         AS responsible_color,
		resp_user.profile_image_path AS responsible_profile_image_path,
	
		created_user.id           AS created_by_id,
		created_user.first_name_en AS created_by_first_name_en,
		created_user.last_name_en  AS created_by_last_name_en,
		created_user.first_name_ar AS created_by_first_name_ar,   
		created_user.last_name_ar  AS created_by_last_name_ar,    
		created_user.initial       AS created_by_initial,         
		created_user.email         AS created_by_email,
		created_user.color         AS created_by_color,
		created_user.profile_image_path AS created_by_profile_image_path,


		updated_user.id           AS updated_by_id,
		updated_user.first_name_en AS updated_by_first_name_en,
		updated_user.last_name_en  AS updated_by_last_name_en,
		updated_user.first_name_ar AS updated_by_first_name_ar,   
		updated_user.last_name_ar  AS updated_by_last_name_ar,    
		updated_user.initial       AS updated_by_initial,         
		updated_user.email         AS updated_by_email,
		updated_user.color         AS updated_by_color,
		updated_user.profile_image_path AS updated_by_profile_image_path
	`).
		Joins(`LEFT JOIN consultation_request_statuses ON consultation_request_statuses.id = consultation_requests.consultation_request_status_id`).
		Joins(`LEFT JOIN organization_units AS req_unit ON req_unit.id = consultation_requests.requesting_department_unit_id`).
		Joins(`LEFT JOIN organization_units AS con_unit ON con_unit.id = consultation_requests.consulted_department_unit_id`).
		Joins(`LEFT JOIN users AS resp_user ON resp_user.id = consultation_requests.responsible_user_id`).
		Joins(`LEFT JOIN users AS created_user ON created_user.id = consultation_requests.created_by_id`).
		Joins(`LEFT JOIN users AS updated_user ON updated_user.id = consultation_requests.updated_by_id`)

	// Add department filter if departmentID is provided (non-admin users)
	if departmentID != nil {
		query = query.Where(
			repo.db.Where("consultation_requests.created_by_id = ?", currentUserID).
				Or("consultation_requests.requesting_department_unit_id = ?", *departmentID).
				Or("consultation_requests.consulted_department_unit_id = ?", *departmentID),
		)
	}

	if params.Keyword != "" {
		kw := "%" + params.Keyword + "%"

		query = query.Where(`
		consultation_requests.reference_number LIKE ? OR
		consultation_requests.request_details LIKE ? OR

		req_unit.title_en LIKE ? OR
		req_unit.title_ar LIKE ? OR

		con_unit.title_en LIKE ? OR
		con_unit.title_ar LIKE ?
	`,
			kw, kw,
			kw, kw,
			kw, kw,
		)
	}

	// Apply dynamic filters
	for column, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where("consultation_requests."+column+" IN (?)", values)
		}
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

	// Skip pagination for export
	if !params.Export {
		offset := (params.Page - 1) * params.PageSize
		query = query.Offset(offset).Limit(params.PageSize)
	}

	if err := query.Find(&requests).Error; err != nil {
		return nil, err
	}

	var totalPages int
	if !params.Export {
		totalPages = int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))
	} else {
		totalPages = 1
	}

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

func (repo *ConsultationRequestRepository) GetAllConsultationRequestComments(consultationRequestID uint) ([]*dtos.ConsultationRequestCommentListDTO, error) {

	var comments []*dtos.ConsultationRequestCommentListDTO

	query := repo.db.Table("consultation_request_comments").
		Select(`
			consultation_request_comments.id,
			consultation_request_comments.consultation_request_id,
			consultation_request_comments.comment,

			consultation_request_statuses.id AS consultation_request_status_id,
			consultation_request_statuses.title_en AS status_title_en,
			consultation_request_statuses.title_ar AS status_title_ar,
			consultation_request_statuses.description_en AS status_description_en,
			consultation_request_statuses.description_ar AS status_description_ar,
			consultation_request_statuses.type AS status_type,
			consultation_request_statuses.status AS status,
			consultation_request_statuses.created_at AS status_created_at,

			consultation_request_comments.created_at,
			users.id AS created_by_id,
			users.first_name_en AS created_by_first_name_en,
			users.last_name_en AS created_by_last_name_en,
			users.first_name_ar AS created_by_first_name_ar,
			users.last_name_ar AS created_by_last_name_ar,
			users.initial As created_by_initial
		`).
		Joins(`LEFT JOIN consultation_request_statuses ON consultation_request_statuses.id = consultation_request_comments.consultation_request_status_id`).
		Joins(`LEFT JOIN users ON users.id = consultation_request_comments.created_by_id`).
		Where("consultation_request_comments.consultation_request_id = ?", consultationRequestID).
		Order("consultation_request_comments.created_at ASC")

	if err := query.Find(&comments).Error; err != nil {
		return nil, err
	}

	for _,comment := range comments {
		var documents []*dtos.ConsultationRequestCommentDocumentDTO
		repo.db.
		Model(&models.ConsultationRequestCommentDocument{}).
		Where("consultation_request_comment_id = ?",comment.ID).
		Find(&documents)

		comment.Documents = documents
	}

	return comments, nil
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

func (repo *ConsultationRequestRepository) CreateComment(comment *models.ConsultationRequestComment, tx ...utility.DatabaseTransaction) (uint, error) {

	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Create(comment).Error; err != nil {
		return 0, err
	}

	return comment.ID, nil
}

func (repo *ConsultationRequestRepository) UpdateComment(comment *models.ConsultationRequestComment, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Model(&models.ConsultationRequestComment{}).
		Where("id = ?", comment.ID).
		Updates(comment).Error
}

func (repo *ConsultationRequestRepository) DeleteComment(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Where("id = ?", id).Delete(&models.ConsultationRequestComment{}).Error
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

func (repo *ConsultationRequestRepository) GetCommentByID(id uint) (*models.ConsultationRequestComment, error) {
	var comment models.ConsultationRequestComment
	if err := repo.db.Preload("Documents").First(&comment, id).Error; err != nil {
		return nil, err
	}
	return &comment, nil
}

func (repo *ConsultationRequestRepository) CreateDocuments(documents []*models.ConsultationRequestCommentDocument, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Create(&documents).Error
}

func (repo *ConsultationRequestRepository) DeleteDocuments(documentIDs []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Where("id IN (?)", documentIDs).Delete(&models.ConsultationRequestCommentDocument{}).Error
}

func (repo *ConsultationRequestRepository) GetCommentDocumentByID(id uint) (*models.ConsultationRequestCommentDocument, error) {
	var document models.ConsultationRequestCommentDocument
	if err := repo.db.Where("id = ?", id).First(&document).Error; err != nil {
		return nil, err
	}

	return &document, nil
}