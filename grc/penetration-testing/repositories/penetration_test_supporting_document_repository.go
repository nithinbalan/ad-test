package repositories

import (
	"fmt"
	"math"
	"time"

	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/models"
	"jethur/grc/penetration-testing/dtos"
	"jethur/utility"

	"gorm.io/gorm"
)

type PenetrationTestSupportingDocumentRepositoryInterface interface {
	GetAll(penetrationTestID uint, params *dtos.PenetrationTestSupportingDocumentListQueryParamsDTO) (*dtos.PenetrationTestSupportingDocumentListResponseDTO, error)
	GetByID(penetrationTestID uint, id uint) (*models.PenetrationTestSupportingDocument, error)
	GetFullDetailsByID(penetrationTestID uint, id uint) (*models.PenetrationTestSupportingDocument, error)
	Create(supportingDoc *models.PenetrationTestSupportingDocument, tx ...utility.DatabaseTransaction) (uint, error)
	Update(supportingDoc *models.PenetrationTestSupportingDocument, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error
	DeleteBulk(ids []uint, tx ...utility.DatabaseTransaction) error
	CreateDocuments(documents []*models.PenetrationTestSupportingDocumentDocument, tx ...utility.DatabaseTransaction) error
	DeleteDocuments(ids []*uint, tx ...utility.DatabaseTransaction) error
	DeleteDocumentsBySupportingDocumentID(supportingDocID uint, tx ...utility.DatabaseTransaction) error
	GetDocumentByID(id uint) (*models.PenetrationTestSupportingDocumentDocument, error)
	BeginTransaction() utility.DatabaseTransaction
}

type PenetrationTestSupportingDocumentRepository struct {
	db *gorm.DB
}

func NewPenetrationTestSupportingDocumentRepository(db *gorm.DB) *PenetrationTestSupportingDocumentRepository {
	return &PenetrationTestSupportingDocumentRepository{db: db}
}

func (ptsdr *PenetrationTestSupportingDocumentRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(ptsdr.db)
}

func (ptsdr *PenetrationTestSupportingDocumentRepository) GetAll(penetrationTestID uint, params *dtos.PenetrationTestSupportingDocumentListQueryParamsDTO) (*dtos.PenetrationTestSupportingDocumentListResponseDTO, error) {
	var supportingDocs []*dtos.PenetrationTestSupportingDocumentListDTO
	var totalCount int64

	query := ptsdr.db.Table("penetration_test_supporting_documents").
		Select(`
			penetration_test_supporting_documents.id,
			penetration_test_supporting_documents.penetration_test_id,
			penetration_test_supporting_documents.title_en,
			penetration_test_supporting_documents.title_ar,
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
			CONVERT(VARCHAR, penetration_test_supporting_documents.created_at, 23) AS created_at,
			CONVERT(VARCHAR, penetration_test_supporting_documents.updated_at, 23) AS updated_at
		`).
		Joins("LEFT JOIN users AS created_user ON created_user.id = penetration_test_supporting_documents.created_by_id").
		Joins("LEFT JOIN users AS updated_user ON updated_user.id = penetration_test_supporting_documents.updated_by_id").
		Where("penetration_test_supporting_documents.penetration_test_id = ?", penetrationTestID)

	if params.Keyword != "" {
		keywordPattern := "%" + params.Keyword + "%"
		query = query.Where(
			"(penetration_test_supporting_documents.title_en LIKE ? OR penetration_test_supporting_documents.title_ar LIKE ?)",
			keywordPattern, keywordPattern,
		)
	}

	if params.DateFrom != "" && params.DateTo != "" {
		dateFrom, errFrom := time.Parse("2006-01-02", params.DateFrom)
		dateTo, errTo := time.Parse("2006-01-02", params.DateTo)
		if errFrom == nil && errTo == nil {
			dateTo = dateTo.Add(24 * time.Hour)
			query = query.Where("penetration_test_supporting_documents.created_at >= ? AND penetration_test_supporting_documents.created_at < ?", dateFrom, dateTo)
		}
	} else if params.DateFrom != "" {
		dateFrom, err := time.Parse("2006-01-02", params.DateFrom)
		if err == nil {
			query = query.Where("penetration_test_supporting_documents.created_at >= ?", dateFrom)
		}
	} else if params.DateTo != "" {
		dateTo, err := time.Parse("2006-01-02", params.DateTo)
		if err == nil {
			dateTo = dateTo.Add(24 * time.Hour)
			query = query.Where("penetration_test_supporting_documents.created_at < ?", dateTo)
		}
	}

	for field, values := range params.Filters {
		if len(values) > 0 {
			switch field {
			case "created_by_ids":
				query = query.Where("penetration_test_supporting_documents.created_by_id IN (?)", values)
			case "updated_by_ids":
				query = query.Where("penetration_test_supporting_documents.updated_by_id IN (?)", values)
			}
		}
	}

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "desc"
	}

	orderClause := "penetration_test_supporting_documents.id DESC"
	if params.SortBy != "" && params.Sort != "" {
		orderClause = fmt.Sprintf("%s %s", params.SortBy, params.Sort)
	}

	offset := (params.Page - 1) * params.PageSize
	if err := query.Order(orderClause).Offset(offset).Limit(params.PageSize).Scan(&supportingDocs).Error; err != nil {
		return nil, err
	}

	for _, docItem := range supportingDocs {
		var docs []*models.PenetrationTestSupportingDocumentDocument
		if err := ptsdr.db.Where("penetration_test_supporting_document_id = ?", docItem.ID).Find(&docs).Error; err == nil {
			for _, d := range docs {
				docItem.Documents = append(docItem.Documents, &dtos.PenetrationTestSupportingDocumentDocumentDTO{
					ID:           d.ID,
					Name:         d.Name,
					OriginalName: d.OriginalName,
					Size:         d.Size,
					Type:         d.Type,
					URL:          d.URL,
					Token:        d.Token,
				})
			}
		}
		if docItem.Documents == nil {
			docItem.Documents = []*dtos.PenetrationTestSupportingDocumentDocumentDTO{}
		}
	}

	if supportingDocs == nil {
		supportingDocs = []*dtos.PenetrationTestSupportingDocumentListDTO{}
	}

	totalPages := int(math.Ceil(float64(totalCount) / float64(params.PageSize)))
	return &dtos.PenetrationTestSupportingDocumentListResponseDTO{
		Data: supportingDocs,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalCount:  totalCount,
			TotalPages:  totalPages,
		},
	}, nil
}

func (ptsdr *PenetrationTestSupportingDocumentRepository) GetByID(penetrationTestID uint, id uint) (*models.PenetrationTestSupportingDocument, error) {
	var supportingDoc models.PenetrationTestSupportingDocument
	err := ptsdr.db.Where("id = ? AND penetration_test_id = ?", id, penetrationTestID).First(&supportingDoc).Error
	if err != nil {
		return nil, err
	}
	return &supportingDoc, nil
}

func (ptsdr *PenetrationTestSupportingDocumentRepository) GetFullDetailsByID(penetrationTestID uint, id uint) (*models.PenetrationTestSupportingDocument, error) {
	var supportingDoc models.PenetrationTestSupportingDocument
	err := ptsdr.db.
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		Preload("Documents").
		Where("id = ? AND penetration_test_id = ?", id, penetrationTestID).
		First(&supportingDoc).Error
	if err != nil {
		return nil, err
	}
	return &supportingDoc, nil
}

func (ptsdr *PenetrationTestSupportingDocumentRepository) Create(supportingDoc *models.PenetrationTestSupportingDocument, tx ...utility.DatabaseTransaction) (uint, error) {
	db := ptsdr.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Create(supportingDoc).Error; err != nil {
		return 0, err
	}
	return supportingDoc.ID, nil
}

func (ptsdr *PenetrationTestSupportingDocumentRepository) Update(supportingDoc *models.PenetrationTestSupportingDocument, tx ...utility.DatabaseTransaction) error {
	db := ptsdr.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Save(supportingDoc).Error
}

func (ptsdr *PenetrationTestSupportingDocumentRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := ptsdr.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Where("id = ?", id).Delete(&models.PenetrationTestSupportingDocument{}).Error
}

func (ptsdr *PenetrationTestSupportingDocumentRepository) DeleteBulk(ids []uint, tx ...utility.DatabaseTransaction) error {
	db := ptsdr.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Where("id IN (?)", ids).Delete(&models.PenetrationTestSupportingDocument{}).Error
}

func (ptsdr *PenetrationTestSupportingDocumentRepository) CreateDocuments(documents []*models.PenetrationTestSupportingDocumentDocument, tx ...utility.DatabaseTransaction) error {
	db := ptsdr.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Create(&documents).Error
}

func (ptsdr *PenetrationTestSupportingDocumentRepository) DeleteDocuments(ids []*uint, tx ...utility.DatabaseTransaction) error {
	db := ptsdr.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Where("id IN (?)", ids).Delete(&models.PenetrationTestSupportingDocumentDocument{}).Error
}

func (ptsdr *PenetrationTestSupportingDocumentRepository) DeleteDocumentsBySupportingDocumentID(supportingDocID uint, tx ...utility.DatabaseTransaction) error {
	db := ptsdr.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Where("penetration_test_supporting_document_id = ?", supportingDocID).Delete(&models.PenetrationTestSupportingDocumentDocument{}).Error
}

func (ptsdr *PenetrationTestSupportingDocumentRepository) GetDocumentByID(id uint) (*models.PenetrationTestSupportingDocumentDocument, error) {
	var document models.PenetrationTestSupportingDocumentDocument
	err := ptsdr.db.Where("id = ?", id).First(&document).Error
	if err != nil {
		return nil, err
	}
	return &document, nil
}
