package repositories

import (
	commonDTOs "jethur/grc/common/dtos"
	slaDTOs "jethur/grc/compliance/dtos"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interfaces
type SLARenewRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(id uint, params *slaDTOs.SLARenewListQueryParamsDTO) (*slaDTOs.SLARenewListResponseDTO, error)

	//Document
	Renew(doc *models.SLADocument, tx ...utility.DatabaseTransaction) error
	SetOlderSLADocsToNotLatest(currentDocID uint, slaRegisterID uint, tx ...utility.DatabaseTransaction) error
	FindPreviousDocument(docId uint, slaRegisterID uint, tx ...utility.DatabaseTransaction) (*models.SLADocument, error)
	GetLatestSLADocumentBySLARegisterID(id uint) (*models.SLADocument, error)
	DeleteDocumentByID(docId uint, tx ...utility.DatabaseTransaction) error
}

// implementation
type SLARenewRepository struct {
	db *gorm.DB
}

func NewSLARenewRepository(db *gorm.DB) *SLARenewRepository {
	return &SLARenewRepository{
		db: db,
	}
}

// BeginTransaction starts a new transaction.
func (repo *SLARenewRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *SLARenewRepository) GetAll(id uint, params *slaDTOs.SLARenewListQueryParamsDTO) (*slaDTOs.SLARenewListResponseDTO, error) {
	var slaRenews []*slaDTOs.SLARenewListDTO
	var totalCount int64
	query := repo.db.Table("sla_documents").
		Select(`
        sla_documents.*,
        created_by_users.first_name_en AS created_by_first_name_en,
        created_by_users.first_name_ar AS created_by_first_name_ar,
        created_by_users.last_name_en AS created_by_last_name_en,
        created_by_users.last_name_ar AS created_by_last_name_ar,
        updated_by_users.first_name_en AS updated_by_first_name_en,
        updated_by_users.first_name_ar AS updated_by_first_name_ar,
        updated_by_users.last_name_en AS updated_by_last_name_en,
        updated_by_users.last_name_ar AS updated_by_last_name_ar
    `).
		Joins(`
        JOIN sla_register ON sla_register.id = sla_documents.sla_register_id
        JOIN users AS created_by_users ON sla_register.created_by = created_by_users.id
        LEFT JOIN users AS updated_by_users ON sla_register.updated_by = updated_by_users.id
    `).
		Where("sla_documents.sla_register_id = ? AND sla_documents.is_latest = ?", id, false)

	if params.IssueDateFrom != "" {
		query = query.Where("issue_date >= ?", params.IssueDateFrom)
	}
	if params.IssueDateTo != "" {
		query = query.Where("issue_date <= ?", params.IssueDateTo)
	}

	if params.IssueDateFrom != "" {
		query = query.Where("expiry_date >= ?", params.IssueDateFrom)
	}
	if params.ExpiryDateTo != "" {
		query = query.Where("expiry_date <= ?", params.ExpiryDateTo)
	}

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	if params.Keyword != "" {
		query = query.Where(`
			(sla_documents.issue_date LIKE ? OR
			sla_documents.expiry_date LIKE ? OR
			sla_documents.sa1_date LIKE ? OR
			sla_documents.sa2_date LIKE ? OR
			sla_documents.name LIKE ? OR
			sla_documents.type LIKE ?)`,
			"%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%",
			"%"+params.Keyword+"%", "%"+params.Keyword+"%")
	}

	if params.SortBy == "" {
		params.SortBy = "name"
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "asc"
	}

	query = query.Order(params.SortBy + " " + params.Sort)

	offset := (params.Page - 1) * params.PageSize

	err := query.Offset(offset).Limit(params.PageSize).Find(&slaRenews).Error

	if err != nil {
		return nil, err
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	response := &slaDTOs.SLARenewListResponseDTO{
		Items: slaRenews,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}
	return response, nil
}

func (repo *SLARenewRepository) Renew(doc *models.SLADocument, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(doc).Error; err != nil {
		return err
	}
	return nil
}

func (repo *SLARenewRepository) SetOlderSLADocsToNotLatest(currentDocID uint, slaRegisterID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Model(&models.SLADocument{}).
		Where("id != ? AND sla_register_id = ?", currentDocID, slaRegisterID).
		Update("is_latest", false).
		Error; err != nil {
		return err
	}
	return nil
}

func (repo *SLARenewRepository) FindPreviousDocument(currentDocID uint, slaRegisterID uint, tx ...utility.DatabaseTransaction) (*models.SLADocument, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var previousDocument models.SLADocument
	if err := db.Model(&models.SLADocument{}).
		Where("sla_register_id = ? AND id != ?", slaRegisterID, currentDocID).
		First(&previousDocument).
		Error; err != nil {
		return nil, err
	}
	return &previousDocument, nil
}

func (repo *SLARenewRepository) DeleteDocumentByID(docID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Delete(&models.SLADocument{}, docID).Error; err != nil {
		return err
	}

	return nil
}

func (repo *SLARenewRepository) GetLatestSLADocumentBySLARegisterID(id uint) (*models.SLADocument, error) {
	var document models.SLADocument
	if err := repo.db.Where("sla_register_id = ?", id).
		Where("is_latest = ?", true).
		First(&document).Error; err != nil {
		return nil, err
	}
	return &document, nil
}
