package repositories

import (
	"jethur/grc/asset/dtos"
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type AssetPolicyRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(params *dtos.AssetPolicyListQueryParamsDTO, assetRegisterID uint) (*dtos.AssetPolicyListResponseDTO, error)
	GetByID(id uint) (*models.AssetPolicy, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.AssetPolicy, error)
	Create(assetPolicy *models.AssetPolicy, tx ...utility.DatabaseTransaction) (uint, error)
	Update(assetPolicy *models.AssetPolicy, tx ...utility.DatabaseTransaction) error
	Delete(id uint) error
	GetByAssetID(assetID uint, tx ...utility.DatabaseTransaction) ([]*models.AssetPolicy, error)

	CreateAssetPolicyChecklists(checklists []*models.AssetPolicyChecklist, tx ...utility.DatabaseTransaction) error
	DeleteAssetPolicyChecklists(policyIDs []uint, tx ...utility.DatabaseTransaction) error
	DeleteAssetPolicyComments(policyIDs []uint, tx ...utility.DatabaseTransaction) error
	DeleteAssetPolicies(policyIDs []uint, tx ...utility.DatabaseTransaction) error
	GetAssetPolicyComplianceSummary(assetRegisterID uint) (*dtos.AssetPolicyByComplianceDTO, error)
}

type AssetPolicyRepository struct {
	db *gorm.DB
}

func NewAssetPolicyRepository(db *gorm.DB) *AssetPolicyRepository {
	return &AssetPolicyRepository{db: db}
}

func (repo *AssetPolicyRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *AssetPolicyRepository) GetAll(params *dtos.AssetPolicyListQueryParamsDTO, assetRegisterID uint) (*dtos.AssetPolicyListResponseDTO, error) {
	var policies []*dtos.AssetPolicyListDTO
	var totalCount int64

	query := repo.db.Table("asset_policies").
		Select(`
			asset_policies.id,
			asset_policies.document_id,
			asset_policies.asset_register_id,
			asset_register.title_en AS asset_register_title_en,
			asset_register.title_ar AS asset_register_title_ar,
			asset_categories.title_en AS asset_category_title_en,
			asset_categories.title_ar AS asset_category_title_ar,
			asset_sub_categories.title_en AS asset_sub_category_title_en,
			asset_sub_categories.title_ar AS asset_sub_category_title_ar,
			compliance_statuses.title_en AS compliance_status_title_en,
			compliance_statuses.title_ar AS compliance_status_title_ar,
			compliance_statuses.color_code AS compliance_status_color_code,
			documents.title_en AS document_title_en,
			documents.title_ar AS document_title_ar,
			document_versions.url AS document_url,
			document_versions.name AS document_name,
			document_versions.size AS document_size,
			document_versions.token AS document_token,
			document_versions.original_name AS document_original_name,
			asset_policies.created_at
		`).
		Joins("LEFT JOIN asset_register ON asset_register.id = asset_policies.asset_register_id").
		Joins("LEFT JOIN asset_categories ON asset_categories.id = asset_register.asset_category_id").
		Joins("LEFT JOIN asset_sub_categories ON asset_sub_categories.id = asset_register.asset_sub_category_id").
		Joins("LEFT JOIN documents ON documents.id = asset_policies.document_id").
		Joins("LEFT JOIN document_versions ON document_versions.document_id = documents.id AND document_versions.is_latest_version = 1").
		Joins("LEFT JOIN compliance_statuses ON asset_policies.compliance_status_id = compliance_statuses.id").
		Where("asset_policies.asset_register_id = ?", assetRegisterID)

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}
	// Apply search filter
	if params.Keyword != "" {
		query = query.Where(`
			documents.title_en LIKE ? OR
			documents.title_ar LIKE ? OR
			compliance_statuses.title_en LIKE ? OR
			compliance_statuses.title_ar LIKE ?
		`, "%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%")
	}

	query = query.Order(params.SortBy + " " + params.Sort).
		Offset((params.Page - 1) * params.PageSize).
		Limit(params.PageSize)

	if err := query.Find(&policies).Error; err != nil {
		return nil, err
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	response := &dtos.AssetPolicyListResponseDTO{
		Items: policies,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}
	return response, nil
}

func (repo *AssetPolicyRepository) GetByID(id uint) (*models.AssetPolicy, error) {
	var assetPolicy *models.AssetPolicy
	result := repo.db.First(&assetPolicy, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return assetPolicy, nil
}

func (repo *AssetPolicyRepository) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.AssetPolicy, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var assetPolicy *models.AssetPolicy

	result := db.
		Preload("AssetRegister").
		Preload("ComplianceStatus").
		Preload("Document").
		Preload("Document.LatestDocumentVersion").
		Preload("ComplianceStatus").
		Preload("AssetPolicyChecklists").
		Preload("AssetPolicyChecklists.DocumentRequirementGroup").
		Preload("AssetPolicyChecklists.ComplianceStatus").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		First(&assetPolicy, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return assetPolicy, nil
}

func (repo *AssetPolicyRepository) Create(assetPolicy *models.AssetPolicy, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(assetPolicy).Error; err != nil {
		return 0, err
	}
	return assetPolicy.ID, nil
}

func (repo *AssetPolicyRepository) Update(assetPolicy *models.AssetPolicy, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Save(assetPolicy).Error
}

func (repo *AssetPolicyRepository) Delete(id uint) error {
	return repo.db.Delete(&models.AssetPolicy{}, id).Error
}

func (repo *AssetPolicyRepository) GetByAssetID(assetID uint, tx ...utility.DatabaseTransaction) ([]*models.AssetPolicy, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var assetPolicies []*models.AssetPolicy
	err := db.Where("asset_register_id = ?", assetID).
		Preload("Document").
		Preload("CreatedByUser").
		Find(&assetPolicies).Error

	if err != nil {
		return nil, err
	}

	return assetPolicies, nil
}

// BulkCreateAssetPolicyChecklists inserts multiple AssetPolicyChecklist records in bulk.
func (repo *AssetPolicyRepository) CreateAssetPolicyChecklists(checklists []*models.AssetPolicyChecklist, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(&checklists).Error; err != nil {
		return err
	}
	return nil
}

// BulkDeleteAssetPolicyChecklists deletes all checklists associated with given AssetPolicy IDs.
func (repo *AssetPolicyRepository) DeleteAssetPolicyChecklists(policyIDs []uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Where("asset_policy_id IN (?)", policyIDs).Delete(&models.AssetPolicyChecklist{}).Error; err != nil {
		return err
	}
	return nil
}
func (repo *AssetPolicyRepository) DeleteAssetPolicyComments(policyIDs []uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Where("asset_policy_id IN (?)", policyIDs).Delete(&models.AssetPolicyComment{}).Error; err != nil {
		return err
	}
	return nil
}

// BulkDeleteAssetPolicies deletes multiple AssetPolicy records in bulk.
func (repo *AssetPolicyRepository) DeleteAssetPolicies(policyIDs []uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Where("id IN (?)", policyIDs).Delete(&models.AssetPolicy{}).Error; err != nil {
		return err
	}
	return nil
}

func (repo *AssetPolicyRepository) GetAssetPolicyComplianceSummary(assetRegisterID uint) (*dtos.AssetPolicyByComplianceDTO, error) {
	result := dtos.AssetPolicyByComplianceDTO{}

	query := repo.db.Table("asset_policies").
		Select(`
			COUNT(asset_policies.id) AS total_policies,
			COUNT(CASE WHEN compliance_statuses.type = 'compliant' THEN asset_policies.id END) AS compliant, 
			COUNT(CASE WHEN compliance_statuses.type != 'compliant' OR compliance_statuses.type IS NULL THEN asset_policies.id END) AS non_compliant,
			(100.0 * COUNT(CASE WHEN compliance_statuses.type = 'compliant' THEN asset_policies.id END) / NULLIF(COUNT(asset_policies.id), 0)) AS compliant_percentage,
			(100.0 * COUNT(CASE WHEN compliance_statuses.type != 'compliant' OR compliance_statuses.type IS NULL THEN asset_policies.id END) / NULLIF(COUNT(asset_policies.id), 0)) AS non_compliant_percentage
		`).
		Joins("LEFT JOIN compliance_statuses ON asset_policies.compliance_status_id = compliance_statuses.id").
		Where("asset_policies.asset_register_id = ?", assetRegisterID)

	err := query.Scan(&result).Error
	return &result, err
}
