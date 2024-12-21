package repositories

import (
	"fmt"
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/models"
	organizationDTOs "jethur/grc/organization/dtos"
	"jethur/utility"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type ProcessRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(params *organizationDTOs.ProcessListQueryParamsDTO) (*organizationDTOs.ProcessListResponseDTO, error)
	GetByID(id uint, tx ...utility.DatabaseTransaction) (*models.Process, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.Process, error)
	Create(process *models.Process, tx ...utility.DatabaseTransaction) (uint, error)
	Update(process *models.Process, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error
	GenerateReferenceNumber() (string, error)

	CreateProcessControls(controls []*models.ProcessControl, tx ...utility.DatabaseTransaction) error
	DeleteProcessControls(controlIds []*uint, tx ...utility.DatabaseTransaction) error

	CreateProcessDocuments(documents []*models.ProcessAssociatedDocument, tx ...utility.DatabaseTransaction) error
	DeleteProcessDocuments(documentIds []*uint, tx ...utility.DatabaseTransaction) error

	CreateProcessAssetRegister(assetRegister []*models.ProcessAsset, tx ...utility.DatabaseTransaction) error
	DeleteProcessAssetRegister(assetRegisterIds []*uint, tx ...utility.DatabaseTransaction) error

	GetDocumentByID(id uint) (*models.Document, error)

	//lifecycle status
	GetIDsByLifecycleStatus(ids []uint, status string) ([]uint, error)
	UpdateLifecycleStatusForIDs(ids []uint, statusId uint) ([]uint, error)

	GetControlSummary(processID uint) ([]models.ProcessControl, error)
}

type ProcessRepository struct {
	db *gorm.DB
}

func NewProcessRepository(db *gorm.DB) *ProcessRepository {
	return &ProcessRepository{
		db: db,
	}
}

func (repo *ProcessRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *ProcessRepository) GetAll(params *organizationDTOs.ProcessListQueryParamsDTO) (*organizationDTOs.ProcessListResponseDTO, error) {
	var processes []*organizationDTOs.ProcessListDTO
	var totalCount int64

	query := repo.db.Table("processes").
		Select(`
		processes.id,
		processes.reference_number,
		processes.title_en,
		processes.title_ar,
		processes.description,
		processes.marked_as_auditable_item,
		process_groups.title_en AS process_group_title_en,
		process_groups.title_ar AS process_group_title_ar,
		processes.rto as rto,
		processes.rto_unit_id as rto_unit_id,
		rto_time_unit.title_en AS rto_unit_title_en,
		rto_time_unit.title_ar AS rto_unit_title_ar,
		processes.mtpd as mtpd,
		processes.mtpd_unit_id as mtpd_unit_id,
		mtpd_time_unit.title_en AS mtpd_unit_title_en,
		mtpd_time_unit.title_ar AS mtpd_unit_title_ar,
		processes.rpo as rpo,
		processes.rpo_unit_id as rpo_unit_id,
		rpo_time_unit.title_en AS rpo_unit_title_en,
		rpo_time_unit.title_ar AS rpo_unit_title_ar,
		processes.is_bia_performed,
		processes.last_assessed_at,
		last_assessed_user.id AS last_assessed_user_id,
		last_assessed_user.first_name_en AS last_assessed_user_first_name_en,
		last_assessed_user.first_name_ar AS last_assessed_user_first_name_ar,
		last_assessed_user.last_name_en AS last_assessed_user_last_name_en,
		last_assessed_user.last_name_ar AS last_assessed_user_last_name_ar,
		last_assessed_user.initial AS last_assessed_user_initial,
		last_assessed_user.color AS last_assessed_user_color,
		last_assessed_user.email AS last_assessed_user_email,
		last_assessed_user.profile_image_path AS last_assessed_user_profile_image_path,
		organization_units.title_en AS owner_unit_title_en,
		organization_units.title_ar AS owner_unit_title_ar,
		owner.id AS owner_id,
		owner.first_name_en AS owner_first_name_en,
		owner.first_name_ar AS owner_first_name_ar,
		owner.last_name_en AS owner_last_name_en,
		owner.last_name_ar AS owner_last_name_ar,
		owner.initial AS owner_initial,
		owner.color AS owner_color,
		owner.email AS owner_email,
		owner.profile_image_path AS owner_profile_image_path,
		lifecycle_statuses.title_en AS lifecycle_status_title_en,
		lifecycle_statuses.title_ar AS lifecycle_status_title_ar,
		lifecycle_statuses.color_code AS lifecycle_status_color_code,
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
		processes.created_at,
		processes.updated_at,
		risk_matrix_risk_ratings.title_en AS risk_rating_title_en,
		risk_matrix_risk_ratings.title_ar AS risk_rating_title_ar,
		risk_matrix_risk_ratings.color_code AS risk_rating_color_code,
		risk_matrix_risk_rating_values.value AS risk_rating_value
	`).
		Joins("LEFT JOIN process_groups ON process_groups.id = processes.process_group_id").
		Joins("LEFT JOIN organization_units ON organization_units.id = processes.owner_unit_id").
		Joins("LEFT JOIN users AS owner ON owner.id = processes.owner_id").
		Joins("LEFT JOIN lifecycle_statuses ON lifecycle_statuses.id = processes.lifecycle_status_id").
		Joins("LEFT JOIN users AS created_user ON created_user.id = processes.created_by_id").
		Joins("LEFT JOIN users AS updated_user ON updated_user.id = processes.updated_by_id").
		Joins("LEFT JOIN users AS last_assessed_user ON last_assessed_user.id = processes.last_assessed_by_id").
		Joins("LEFT JOIN time_units AS rto_time_unit ON rto_time_unit.id = processes.rto_unit_id").
		Joins("LEFT JOIN time_units AS mtpd_time_unit ON mtpd_time_unit.id = processes.mtpd_unit_id").
		Joins("LEFT JOIN time_units AS rpo_time_unit ON rpo_time_unit.id = processes.rpo_unit_id").
		Joins("LEFT JOIN risk_matrix_risk_ratings ON risk_matrix_risk_ratings.id = processes.risk_matrix_risk_rating_id").
		Joins("LEFT JOIN risk_matrix_risk_rating_values ON risk_matrix_risk_rating_values.id = processes.risk_matrix_risk_rating_value_id")

	// Apply dynamic filters to the query
	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where("processes."+columnName+" IN (?)", values)
		}
	}

	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("processes.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	// Apply quick filters
	for _, filter := range params.QuickFilter {
		switch filter {
		case "active":
			query = query.Where("lifecycle_statuses.type = ?", "active")
		case "archived":
			query = query.Where("lifecycle_statuses.type = ?", "archived")
		case "high_risk":
			query = query.Where("risk_matrix_risk_ratings.type IN (?)", []string{"high", "very_high"})
		}
	}

	// Apply search keyword if exists
	if params.Keyword != "" {
		query = query.Where(`
		(processes.reference_number LIKE ? OR
		 processes.title_en LIKE ? OR
		 processes.title_ar LIKE ? OR
		 processes.description LIKE ? OR
		 process_groups.title_en LIKE ? OR
		 process_groups.title_ar LIKE ? OR
		 organization_units.title_en LIKE ? OR
		 organization_units.title_ar LIKE ? OR
		 owner.first_name_en LIKE ? OR
		 owner.first_name_ar LIKE ? OR
		 lifecycle_statuses.title_en LIKE ? OR
		 lifecycle_statuses.title_ar LIKE ?)`,
			"%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%",
			"%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%",
			"%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%")
	}

	// Count total records for pagination
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	query = query.Order(params.SortBy + " " + params.Sort)

	if !params.Export {
		offset := (params.Page - 1) * params.PageSize
		query = query.Offset(offset).Limit(params.PageSize)
	}

	err := query.Find(&processes).Error
	if err != nil {
		return nil, err
	}

	totalPages := 1
	if !params.Export {
		totalPages = int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))
	}

	response := &organizationDTOs.ProcessListResponseDTO{
		Items: processes,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}
	return response, nil
}

func (repo *ProcessRepository) GetByID(id uint, tx ...utility.DatabaseTransaction) (*models.Process, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var process *models.Process
	result := db.First(&process, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return process, nil
}

func (repo *ProcessRepository) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.Process, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var process *models.Process

	result := db.
		Preload("ProcessGroup").
		Preload("OwnerUnit").
		Preload("Owner").
		Preload("Owner.OrganizationUnit1").
		Preload("LifecycleStatus").
		Preload("Controls").
		Preload("Controls.Control").
		Preload("Assets").
		Preload("Assets.Asset").
		Preload("Documents").
		Preload("Documents.Document").
		Preload("Risks").
		Preload("RiskMatrixRating").
		Preload("RiskMatrixRatingValue").
		Preload("CreatedByUser").
		Preload("CreatedByUser.OrganizationUnit1").
		Preload("UpdatedByUser").
		Preload("UpdatedByUser.OrganizationUnit1").
		Preload("RtoUnit").
		Preload("MtpdUnit").
		Preload("RpoUnit").
		Preload("LastAssessedByUser").
		Preload("AuditableItems", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at DESC").Limit(5) // Fetch the latest 5 auditable items
		}).
		Preload("AuditableItems.ResponsibleUser").
		Preload("AuditableItems.ResponsibleUnit").
		Preload("AuditableItems.AuditorRiskRating").
		Preload("AuditableItems.CreatedByUser").
		First(&process, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return process, nil
}

func (repo *ProcessRepository) Create(process *models.Process, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(process).Error; err != nil {
		return 0, err
	}
	return process.ID, nil
}

func (repo *ProcessRepository) Update(process *models.Process, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(process).Error; err != nil {
		return err
	}
	return nil
}

func (repo *ProcessRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Delete(&models.Process{}, id).Error; err != nil {
		return err
	}

	return nil
}

func (repo *ProcessRepository) GetIDsByLifecycleStatus(ids []uint, status string) ([]uint, error) {
	var matchingIDs []uint

	if err := repo.db.Model(&models.Process{}).
		Joins("JOIN lifecycle_statuses ON process.lifecycle_status_id = lifecycle_statuses.id").
		Where("process.id IN ? AND lifecycle_statuses.type = ?", ids, status).
		Pluck("process.id", &matchingIDs).Error; err != nil {
		return nil, err
	}

	return matchingIDs, nil
}

func (repo *ProcessRepository) UpdateLifecycleStatusForIDs(ids []uint, statusId uint) ([]uint, error) {
	if err := repo.db.Model(&models.Process{}).
		Where("id IN ?", ids).
		Update("lifecycle_status_id", statusId).Error; err != nil {
		return nil, err
	}

	return ids, nil
}

func (repo *ProcessRepository) GenerateReferenceNumber() (string, error) {
	var lastReferenceNumber string
	err := repo.db.Model(&models.Process{}).
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
	newReferenceNumber := fmt.Sprintf("%s-%d-%04d", "PR", currentYear, lastReferenceID+1)

	return newReferenceNumber, nil
}

func (repo *ProcessRepository) CreateProcessControls(controls []*models.ProcessControl, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(&controls).Error; err != nil {
		return err
	}
	return nil
}

func (repo *ProcessRepository) DeleteProcessControls(controlIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	// Delete controls by IDs
	if err := db.Where("id IN (?)", controlIds).Delete(&models.ProcessControl{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *ProcessRepository) CreateProcessAssetRegister(assetRegister []*models.ProcessAsset, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(&assetRegister).Error; err != nil {
		return err
	}
	return nil
}

func (repo *ProcessRepository) DeleteProcessAssetRegister(AssetRegisterIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	// Delete asset register by IDs
	if err := db.Where("id IN (?)", AssetRegisterIds).Delete(&models.ProcessAsset{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *ProcessRepository) CreateProcessDocuments(documents []*models.ProcessAssociatedDocument, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(&documents).Error; err != nil {
		return err
	}
	return nil
}

func (repo *ProcessRepository) DeleteProcessDocuments(documentIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("id IN (?)", documentIds).Delete(&models.ProcessAssociatedDocument{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *ProcessRepository) GetDocumentByID(id uint) (*models.Document, error) {
	var document *models.Document
	result := repo.db.First(&document, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return document, nil
}

func (repo *ProcessRepository) GetControlSummary(processID uint) ([]models.ProcessControl, error) {
	var processControls []models.ProcessControl

	result := repo.db.
		Preload("Control").
		Preload("Control.ControlImplementationStatus").
		Preload("Control.ControlEfficiencyStatus").
		Where("process_id = ?", processID).
		Find(&processControls)
	if result.Error != nil {
		return nil, result.Error
	}

	return processControls, nil
}
