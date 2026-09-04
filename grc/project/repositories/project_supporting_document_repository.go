package repositories

import (
	"fmt"
	"math"
	"time"

	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/models"
	"jethur/grc/project/dtos"
	"jethur/utility"

	"gorm.io/gorm"
)

type ProjectSupportingDocumentRepositoryInterface interface {
	GetAll(projectID uint, params *dtos.ProjectSupportingDocumentListQueryParamsDTO) (*dtos.ProjectSupportingDocumentListResponseDTO, error)
	GetByID(projectID uint, id uint) (*models.ProjectSupportingDocument, error)
	GetFullDetailsByID(projectID uint, id uint) (*models.ProjectSupportingDocument, error)
	Create(supportingDoc *models.ProjectSupportingDocument, tx ...utility.DatabaseTransaction) (uint, error)
	Update(supportingDoc *models.ProjectSupportingDocument, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error
	DeleteBulk(ids []uint, tx ...utility.DatabaseTransaction) error
	CreateDocuments(documents []*models.ProjectSupportingDocumentDocument, tx ...utility.DatabaseTransaction) error
	DeleteDocuments(ids []*uint, tx ...utility.DatabaseTransaction) error
	DeleteDocumentsBySupportingDocumentID(supportingDocID uint, tx ...utility.DatabaseTransaction) error
	GetDocumentByID(id uint) (*models.ProjectSupportingDocumentDocument, error)
	BeginTransaction() utility.DatabaseTransaction
}

type ProjectSupportingDocumentRepository struct {
	db *gorm.DB
}

func NewProjectSupportingDocumentRepository(db *gorm.DB) *ProjectSupportingDocumentRepository {
	return &ProjectSupportingDocumentRepository{db: db}
}

func (psdr *ProjectSupportingDocumentRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(psdr.db)
}

func (psdr *ProjectSupportingDocumentRepository) GetAll(projectID uint, params *dtos.ProjectSupportingDocumentListQueryParamsDTO) (*dtos.ProjectSupportingDocumentListResponseDTO, error) {
	var supportingDocs []*dtos.ProjectSupportingDocumentListDTO
	var totalCount int64

	query := psdr.db.Table("project_supporting_documents").
		Select(`
			project_supporting_documents.id,
			project_supporting_documents.project_id,
			project_supporting_documents.title_en,
			project_supporting_documents.title_ar,
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
			CONVERT(VARCHAR, project_supporting_documents.created_at, 23) AS created_at,
			CONVERT(VARCHAR, project_supporting_documents.updated_at, 23) AS updated_at
		`).
		Joins("LEFT JOIN users AS created_user ON created_user.id = project_supporting_documents.created_by_id").
		Joins("LEFT JOIN users AS updated_user ON updated_user.id = project_supporting_documents.updated_by_id").
		Where("project_supporting_documents.project_id = ?", projectID)

	if params.Keyword != "" {
		keywordPattern := "%" + params.Keyword + "%"
		query = query.Where(
			"project_supporting_documents.title_en LIKE ? OR project_supporting_documents.title_ar LIKE ?",
			keywordPattern, keywordPattern,
		)
	}

	if params.DateFrom != "" && params.DateTo != "" {
		dateFrom, errFrom := time.Parse("2006-01-02", params.DateFrom)
		dateTo, errTo := time.Parse("2006-01-02", params.DateTo)
		if errFrom == nil && errTo == nil {
			dateTo = dateTo.Add(24 * time.Hour)
			query = query.Where("project_supporting_documents.created_at >= ? AND project_supporting_documents.created_at < ?", dateFrom, dateTo)
		}
	} else if params.DateFrom != "" {
		dateFrom, err := time.Parse("2006-01-02", params.DateFrom)
		if err == nil {
			query = query.Where("project_supporting_documents.created_at >= ?", dateFrom)
		}
	} else if params.DateTo != "" {
		dateTo, err := time.Parse("2006-01-02", params.DateTo)
		if err == nil {
			dateTo = dateTo.Add(24 * time.Hour)
			query = query.Where("project_supporting_documents.created_at < ?", dateTo)
		}
	}

	for field, values := range params.Filters {
		if len(values) > 0 {
			switch field {
			case "created_by_ids":
				query = query.Where("project_supporting_documents.created_by_id IN (?)", values)
			case "updated_by_ids":
				query = query.Where("project_supporting_documents.updated_by_id IN (?)", values)
			}
		}
	}

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	orderClause := "project_supporting_documents.id DESC"
	if params.SortBy != "" && params.Sort != "" {
		orderClause = fmt.Sprintf("%s %s", params.SortBy, params.Sort)
	}

	offset := (params.Page - 1) * params.PageSize
	if err := query.Order(orderClause).Offset(offset).Limit(params.PageSize).Scan(&supportingDocs).Error; err != nil {
		return nil, err
	}

	// Fetch documents for each item
	for _, docItem := range supportingDocs {
		var docs []*models.ProjectSupportingDocumentDocument
		if err := psdr.db.Where("project_supporting_document_id = ?", docItem.ID).Find(&docs).Error; err == nil {
			for _, d := range docs {
				docItem.Documents = append(docItem.Documents, &dtos.ProjectSupportingDocumentDocumentDTO{
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
			docItem.Documents = []*dtos.ProjectSupportingDocumentDocumentDTO{}
		}
	}

	if supportingDocs == nil {
		supportingDocs = []*dtos.ProjectSupportingDocumentListDTO{}
	}

	totalPages := int(math.Ceil(float64(totalCount) / float64(params.PageSize)))
	return &dtos.ProjectSupportingDocumentListResponseDTO{
		Data: supportingDocs,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalCount:  totalCount,
			TotalPages:  totalPages,
		},
	}, nil
}

func (psdr *ProjectSupportingDocumentRepository) GetByID(projectID uint, id uint) (*models.ProjectSupportingDocument, error) {
	var supportingDoc models.ProjectSupportingDocument
	err := psdr.db.Where("id = ? AND project_id = ?", id, projectID).First(&supportingDoc).Error
	if err != nil {
		return nil, err
	}
	return &supportingDoc, nil
}

func (psdr *ProjectSupportingDocumentRepository) GetFullDetailsByID(projectID uint, id uint) (*models.ProjectSupportingDocument, error) {
	var supportingDoc models.ProjectSupportingDocument
	err := psdr.db.
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		Preload("Documents").
		Where("id = ? AND project_id = ?", id, projectID).
		First(&supportingDoc).Error
	if err != nil {
		return nil, err
	}
	return &supportingDoc, nil
}

func (psdr *ProjectSupportingDocumentRepository) Create(supportingDoc *models.ProjectSupportingDocument, tx ...utility.DatabaseTransaction) (uint, error) {
	db := psdr.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Create(supportingDoc).Error; err != nil {
		return 0, err
	}
	return supportingDoc.ID, nil
}

func (psdr *ProjectSupportingDocumentRepository) Update(supportingDoc *models.ProjectSupportingDocument, tx ...utility.DatabaseTransaction) error {
	db := psdr.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Save(supportingDoc).Error
}

func (psdr *ProjectSupportingDocumentRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := psdr.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Where("id = ?", id).Delete(&models.ProjectSupportingDocument{}).Error
}

func (psdr *ProjectSupportingDocumentRepository) DeleteBulk(ids []uint, tx ...utility.DatabaseTransaction) error {
	db := psdr.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Where("id IN (?)", ids).Delete(&models.ProjectSupportingDocument{}).Error
}

func (psdr *ProjectSupportingDocumentRepository) CreateDocuments(documents []*models.ProjectSupportingDocumentDocument, tx ...utility.DatabaseTransaction) error {
	db := psdr.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Create(&documents).Error
}

func (psdr *ProjectSupportingDocumentRepository) DeleteDocuments(ids []*uint, tx ...utility.DatabaseTransaction) error {
	db := psdr.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Where("id IN (?)", ids).Delete(&models.ProjectSupportingDocumentDocument{}).Error
}

func (psdr *ProjectSupportingDocumentRepository) DeleteDocumentsBySupportingDocumentID(supportingDocID uint, tx ...utility.DatabaseTransaction) error {
	db := psdr.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Where("project_supporting_document_id = ?", supportingDocID).Delete(&models.ProjectSupportingDocumentDocument{}).Error
}

func (psdr *ProjectSupportingDocumentRepository) GetDocumentByID(id uint) (*models.ProjectSupportingDocumentDocument, error) {
	var document models.ProjectSupportingDocumentDocument
	err := psdr.db.Where("id = ?", id).First(&document).Error
	if err != nil {
		return nil, err
	}
	return &document, nil
}
