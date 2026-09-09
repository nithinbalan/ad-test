package repositories

import (
	"fmt"
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/models"
	socDtos "jethur/grc/security-operations/dtos"
	"jethur/utility"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type CSLogSocRoomAccessLogRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(params *socDtos.CSLogSocRoomAccessLogListQueryParamsDTO) (*socDtos.CSLogSocRoomAccessLogListResponseDTO, error)
	GetByID(id uint) (*models.CSLogSocRoomAccessLog, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.CSLogSocRoomAccessLog, error)
	Create(entry *models.CSLogSocRoomAccessLog, tx ...utility.DatabaseTransaction) (uint, error)
	Update(entry *models.CSLogSocRoomAccessLog, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error
	GenerateReferenceNumber() (string, error)

	//Documents
	CreateDocuments(docs []*models.CSLogSocRoomAccessLogDocument, tx ...utility.DatabaseTransaction) error
	GetDocumentByID(id uint) (*models.CSLogSocRoomAccessLogDocument, error)
	DeleteDocuments(docIds []*uint, tx ...utility.DatabaseTransaction) error
}

type CSLogSocRoomAccessLogRepository struct {
	db *gorm.DB
}

func NewCSLogSocRoomAccessLogRepository(db *gorm.DB) *CSLogSocRoomAccessLogRepository {
	return &CSLogSocRoomAccessLogRepository{
		db: db,
	}
}

func (repo *CSLogSocRoomAccessLogRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *CSLogSocRoomAccessLogRepository) GetAll(params *socDtos.CSLogSocRoomAccessLogListQueryParamsDTO) (*socDtos.CSLogSocRoomAccessLogListResponseDTO, error) {
	var entries []*socDtos.CSLogSocRoomAccessLogListDTO
	var totalCount int64

	query := repo.db.Table("cs_log_soc_room_access_logs AS soc").
		Select(`
			soc.id,
			soc.reference_number,
			soc.visitor_full_name,
			CONVERT(VARCHAR, soc.entry_time, 120) AS entry_time,
			CONVERT(VARCHAR, soc.exit_time, 120) AS exit_time,
			soc.visitor_id_number,
			soc.department_id,
			department.title_en AS department_title_en,
			department.title_ar AS department_title_ar,
			soc.external_party_name,
			soc.purpose_of_entry,
			soc.authorizing_officer_id,
			authorizing_officer.first_name_en AS authorizing_officer_first_name_en,
			authorizing_officer.first_name_ar AS authorizing_officer_first_name_ar,
			authorizing_officer.last_name_en AS authorizing_officer_last_name_en,
			authorizing_officer.last_name_ar AS authorizing_officer_last_name_ar,
			authorizing_officer.initial AS authorizing_officer_initial,
			authorizing_officer.color AS authorizing_officer_color,
			authorizing_officer.profile_image_path AS authorizing_officer_profile_image_path,
			CONVERT(VARCHAR, soc.created_at, 120) AS created_at,
			CONVERT(VARCHAR, soc.updated_at, 120) AS updated_at,
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
			updated_user.profile_image_path AS updated_user_profile_image_path
		`).
		Joins("LEFT JOIN organization_units AS department ON department.id = soc.department_id").
		Joins("LEFT JOIN users AS authorizing_officer ON authorizing_officer.id = soc.authorizing_officer_id").
		Joins("LEFT JOIN users AS created_user ON created_user.id = soc.created_by_id").
		Joins("LEFT JOIN users AS updated_user ON updated_user.id = soc.updated_by_id")

	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
			(soc.visitor_full_name LIKE ? OR
			soc.reference_number LIKE ? OR
			soc.visitor_id_number LIKE ? OR
			soc.external_party_name LIKE ? OR
			department.title_en LIKE ? OR
			department.title_ar LIKE ? OR
			authorizing_officer.first_name_en LIKE ? OR
			authorizing_officer.first_name_ar LIKE ? OR
			authorizing_officer.last_name_en LIKE ? OR
			authorizing_officer.last_name_ar LIKE ? OR
			CONCAT(authorizing_officer.first_name_en, ' ', authorizing_officer.last_name_en) LIKE ? OR
			CONCAT(authorizing_officer.first_name_ar, ' ', authorizing_officer.last_name_ar) LIKE ?)`,
			keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword)
	}

	filterColumns := map[string]string{
		"department_id":          "soc.department_id",
		"authorizing_officer_id": "soc.authorizing_officer_id",
	}
	for filterName, values := range params.Filters {
		if columnName, ok := filterColumns[filterName]; ok && len(values) > 0 {
			query = query.Where(columnName+" IN (?)", values)
		}
	}

	if params.DateFrom != "" {
		query = query.Where("soc.entry_time >= ?", params.DateFrom+" 00:00:00")
	}
	if params.DateTo != "" {
		query = query.Where("soc.entry_time <= ?", params.DateTo+" 23:59:59")
	}

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	query = query.Order(params.SortBy + " " + params.Sort)

	if !params.Export {
		offset := (params.Page - 1) * params.PageSize
		query = query.Offset(offset).Limit(params.PageSize)
	}

	err := query.Find(&entries).Error
	if err != nil {
		return nil, err
	}

	totalPages := 1
	if !params.Export {
		totalPages = int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))
	}

	response := &socDtos.CSLogSocRoomAccessLogListResponseDTO{
		Items: entries,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}
	return response, nil
}

func (repo *CSLogSocRoomAccessLogRepository) GetByID(id uint) (*models.CSLogSocRoomAccessLog, error) {
	var entry *models.CSLogSocRoomAccessLog

	result := repo.db.First(&entry, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return entry, nil
}

func (repo *CSLogSocRoomAccessLogRepository) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.CSLogSocRoomAccessLog, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var entry *models.CSLogSocRoomAccessLog
	result := db.
		Preload("Department").
		Preload("AuthorizingOfficer").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		Preload("Documents").
		First(&entry, id)

	if result.Error != nil {
		return nil, result.Error
	}
	return entry, nil
}

func (repo *CSLogSocRoomAccessLogRepository) Create(entry *models.CSLogSocRoomAccessLog, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(entry).Error; err != nil {
		return 0, err
	}
	return entry.ID, nil
}

func (repo *CSLogSocRoomAccessLogRepository) Update(updated *models.CSLogSocRoomAccessLog, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Save(updated).Error; err != nil {
		return err
	}
	return nil
}

func (repo *CSLogSocRoomAccessLogRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Delete(&models.CSLogSocRoomAccessLog{}, id).Error; err != nil {
		return err
	}

	return nil
}

func (repo *CSLogSocRoomAccessLogRepository) GenerateReferenceNumber() (string, error) {
	var lastReferenceNumber string
	err := repo.db.Model(&models.CSLogSocRoomAccessLog{}).
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
	lastReferenceID, err := strconv.ParseUint(lastDigits, 10, 64)
	if err != nil {
		return "", err
	}
	newReferenceNumber := fmt.Sprintf("SRA-%d-%04d", currentYear, lastReferenceID+1)

	return newReferenceNumber, nil
}

// Document func
func (repo *CSLogSocRoomAccessLogRepository) CreateDocuments(docs []*models.CSLogSocRoomAccessLogDocument, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(&docs).Error; err != nil {
		return err
	}
	return nil
}

func (repo *CSLogSocRoomAccessLogRepository) GetDocumentByID(id uint) (*models.CSLogSocRoomAccessLogDocument, error) {
	var document models.CSLogSocRoomAccessLogDocument
	if err := repo.db.Where("id = ?", id).First(&document).Error; err != nil {
		return nil, err
	}
	return &document, nil
}

func (repo *CSLogSocRoomAccessLogRepository) DeleteDocuments(docIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Where("id IN (?)", docIds).Delete(&models.CSLogSocRoomAccessLogDocument{}).Error; err != nil {
		return err
	}
	return nil
}
