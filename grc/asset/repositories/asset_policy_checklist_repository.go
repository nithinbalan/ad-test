package repositories

import (
	"jethur/grc/asset/dtos"
	complianceDTOs "jethur/grc/compliance/dtos"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type AssetPolicyChecklistRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetByID(id uint) (*models.AssetPolicyChecklist, error)
	Create(assetPolicyChecklist *models.AssetPolicyChecklist, tx ...utility.DatabaseTransaction) (uint, error)
	Update(checklist *models.AssetPolicyChecklist, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error

	DeleteAssetPolicyChecklists(ids []uint, tx ...utility.DatabaseTransaction) error
	BulkUpdateAssetPolicyChecklists(checklists []*models.AssetPolicyChecklist, tx ...utility.DatabaseTransaction) error
	GetAssetPolicyChecklistSummary(assetPolicyID uint) (*dtos.AssetPolicyChecklistSummaryDTO, error)
}

type AssetPolicyChecklistRepository struct {
	db *gorm.DB
}

func NewAssetPolicyChecklistRepository(db *gorm.DB) *AssetPolicyChecklistRepository {
	return &AssetPolicyChecklistRepository{db: db}
}

func (repo *AssetPolicyChecklistRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *AssetPolicyChecklistRepository) GetByID(id uint) (*models.AssetPolicyChecklist, error) {
	var assetPolicyChecklist *models.AssetPolicyChecklist
	result := repo.db.First(&assetPolicyChecklist, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return assetPolicyChecklist, nil
}

func (repo *AssetPolicyChecklistRepository) Create(assetPolicyChecklist *models.AssetPolicyChecklist, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(assetPolicyChecklist).Error; err != nil {
		return 0, err
	}
	return assetPolicyChecklist.ID, nil
}

func (repo *AssetPolicyChecklistRepository) Update(checklist *models.AssetPolicyChecklist, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Save(checklist).Error
}

func (repo *AssetPolicyChecklistRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Where("id = ?", id).Delete(&models.AssetPolicyChecklist{}).Error
}

func (repo *AssetPolicyChecklistRepository) DeleteAssetPolicyChecklists(ids []uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	// Delete checklists in bulk
	return db.Where("id IN ?", ids).Delete(&models.AssetPolicyChecklist{}).Error
}

func (repo *AssetPolicyChecklistRepository) BulkUpdateAssetPolicyChecklists(checklists []*models.AssetPolicyChecklist, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	// Save each checklist individually in a loop
	for _, checklist := range checklists {
		if err := db.Save(checklist).Error; err != nil {
			return err
		}
	}

	return nil
}

func (repo *AssetPolicyChecklistRepository) GetAssetPolicyChecklistSummary(assetPolicyID uint) (*dtos.AssetPolicyChecklistSummaryDTO, error) {
	summary := &dtos.AssetPolicyChecklistSummaryDTO{}
	checklistCounts := dtos.ChecklistSummaryDTO{}

	// Query to get checklist status counts
	err := repo.db.
		Table("asset_policy_checklists").
		Select(`
			COUNT(*) AS total_checklist,
			COUNT(CASE WHEN compliance_statuses.type = 'compliant' THEN 1 END) AS compliant,
			COUNT(CASE WHEN compliance_statuses.type = 'partially_compliant' THEN 1 END) AS partially_compliant,
			COUNT(CASE WHEN compliance_statuses.type = 'not_applicable' THEN 1 END) AS not_applicable,
			COUNT(CASE WHEN compliance_statuses.type NOT IN ('compliant', 'partially_compliant', 'not_applicable') OR compliance_statuses.type IS NULL THEN 1 END) AS non_compliant
		`).
		Joins("LEFT JOIN compliance_statuses ON compliance_statuses.id = asset_policy_checklists.compliance_status_id").
		Where("asset_policy_checklists.asset_policy_id = ?", assetPolicyID).
		Scan(&checklistCounts).Error
	if err != nil {
		return nil, err
	}
	summary.Summary = checklistCounts

	// Query to get asset policy's overall compliance status
	var complianceStatus complianceDTOs.ComplianceStatusDTO
	err = repo.db.
		Table("asset_policies").
		Select(`
			compliance_statuses.id,
			compliance_statuses.title_en,
			compliance_statuses.title_ar,
			compliance_statuses.color_code,
			compliance_statuses.type
		`).
		Joins("LEFT JOIN compliance_statuses ON compliance_statuses.id = asset_policies.compliance_status_id").
		Where("asset_policies.id = ?", assetPolicyID).
		Scan(&complianceStatus).Error
	if err == nil && complianceStatus.ID != 0 {
		summary.ComplianceStatus = &complianceStatus
	}

	return summary, nil
}
