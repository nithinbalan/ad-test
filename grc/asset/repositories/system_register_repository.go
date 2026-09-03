package repositories

import (
	"jethur/grc/asset/dtos"
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

// Interfaces
type SystemRegisterRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(params *dtos.SystemRegisterListQueryParamsDTO) (*dtos.SystemRegisterListResponseDTO, error)
	GetByID(id uint) (*models.SystemRegister, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.SystemRegister, error)
	Create(systemRegister *models.SystemRegister, tx ...utility.DatabaseTransaction) (uint, error)
	Update(systemRegister *models.SystemRegister, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error
	CreateSLA(sla *models.SystemRegisterSLA, tx ...utility.DatabaseTransaction) error
	DeleteSLAs(systemRegisterID uint, tx ...utility.DatabaseTransaction) error
	GetSLARegisterByID(slaRegisterID uint, tx ...utility.DatabaseTransaction) (*models.SLARegister, error)
	CreateSystemRegisterImpactType(impactType *models.SystemRegisterImpactType, tx ...utility.DatabaseTransaction) error
	DeleteSystemRegisterImpactTypesBySystemRegisterID(systemRegisterID uint, tx ...utility.DatabaseTransaction) error
	GetImpactCategoryByID(id uint, tx ...utility.DatabaseTransaction) (*models.ImpactCategory, error)
}

// Implementation
type SystemRegisterRepository struct {
	db *gorm.DB
}

func NewSystemRegisterRepository(db *gorm.DB) *SystemRegisterRepository {
	return &SystemRegisterRepository{
		db: db,
	}
}

// BeginTransaction starts a new transaction.
func (repo *SystemRegisterRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *SystemRegisterRepository) GetAll(params *dtos.SystemRegisterListQueryParamsDTO) (*dtos.SystemRegisterListResponseDTO, error) {
	var systemRegisters []*dtos.SystemRegisterListDTO
	var totalCount int64

	query := repo.db.Table("system_registers").
		Select(`
			system_registers.id,
			system_registers.title_en,
			system_registers.title_ar,
			system_registers.service_description,
			system_registers.modules_components,
			system_registers.third_party_dependencies,
			system_registers.service_importance_id,
			service_importance.title_en AS service_importance_title_en,
			service_importance.title_ar AS service_importance_title_ar,
			system_registers.responsible_it_team_id,
			responsible_it_team.title_en AS responsible_it_team_title_en,
			responsible_it_team.title_ar AS responsible_it_team_title_ar,
			system_registers.outsourcing_partner,
			system_registers.operation_location_id,
			operation_location.title_en AS operation_location_title_en,
			operation_location.title_ar AS operation_location_title_ar,
			system_registers.business_owner_id,
			business_owner.title_en AS business_owner_title_en,
			business_owner.title_ar AS business_owner_title_ar,
			system_registers.customers,
			system_registers.disruption_impact_level_id,
			disruption_impact_level.title_en AS disruption_impact_level_title_en,
			disruption_impact_level.title_ar AS disruption_impact_level_title_ar,
			system_registers.rto,
			system_registers.rto_unit_id,
			rto_unit.title_en AS rto_unit_title_en,
			rto_unit.title_ar AS rto_unit_title_ar,
			system_registers.mtpd,
			system_registers.mtpd_unit_id,
			mtpd_unit.title_en AS mtpd_unit_title_en,
			mtpd_unit.title_ar AS mtpd_unit_title_ar,
			system_registers.rpo,
			system_registers.rpo_unit_id,
			rpo_unit.title_en AS rpo_unit_title_en,
			rpo_unit.title_ar AS rpo_unit_title_ar,
			system_registers.it_processes,
			system_registers.it_service_owner_id,
			it_service_owner.title_en AS it_service_owner_title_en,
			it_service_owner.title_ar AS it_service_owner_title_ar,
			system_registers.process_operation_team_id,
			process_operation_team.title_en AS process_operation_team_title_en,
			process_operation_team.title_ar AS process_operation_team_title_ar,
			system_registers.process_recovery_team_id,
			process_recovery_team.title_en AS process_recovery_team_title_en,
			process_recovery_team.title_ar AS process_recovery_team_title_ar,
			system_registers.recovery_manager_id,
			recovery_manager.title_en AS recovery_manager_title_en,
			recovery_manager.title_ar AS recovery_manager_title_ar,
			system_registers.recovery_manager_phone,
			system_registers.recovery_manager_email,
			system_registers.recovery_location,
			system_registers.vendor_rto,
			system_registers.vendor_rto_unit_id,
			vendor_rto_unit.title_en AS vendor_rto_unit_title_en,
			vendor_rto_unit.title_ar AS vendor_rto_unit_title_ar,
			created_user.id AS created_user_id,
			created_user.first_name_en AS created_user_first_name_en,
			created_user.first_name_ar AS created_user_first_name_ar,
			created_user.last_name_en AS created_user_last_name_en,
			created_user.last_name_ar AS created_user_last_name_ar,
			created_user.initial AS created_user_initial,
			created_user.color AS created_user_color,
			created_user.email AS created_user_email,
			created_user.profile_image_path AS created_user_profile_image_path,
			system_registers.created_at,
			updated_user.id AS updated_user_id,
			updated_user.first_name_en AS updated_user_first_name_en,
			updated_user.first_name_ar AS updated_user_first_name_ar,
			updated_user.last_name_en AS updated_user_last_name_en,
			updated_user.last_name_ar AS updated_user_last_name_ar,
			updated_user.initial AS updated_user_initial,
			updated_user.color AS updated_user_color,
			updated_user.email AS updated_user_email,
			updated_user.profile_image_path AS updated_user_profile_image_path,
			system_registers.updated_at
		`).
		Joins("LEFT JOIN system_register_service_importances AS service_importance ON service_importance.id = system_registers.service_importance_id").
		Joins("LEFT JOIN job_teams AS responsible_it_team ON responsible_it_team.id = system_registers.responsible_it_team_id").
		Joins("LEFT JOIN locations AS operation_location ON operation_location.id = system_registers.operation_location_id").
		Joins("LEFT JOIN organization_units AS business_owner ON business_owner.id = system_registers.business_owner_id").
		Joins("LEFT JOIN system_register_disruption_levels AS disruption_impact_level ON disruption_impact_level.id = system_registers.disruption_impact_level_id").
		Joins("LEFT JOIN designations AS it_service_owner ON it_service_owner.id = system_registers.it_service_owner_id").
		Joins("LEFT JOIN job_teams AS process_operation_team ON process_operation_team.id = system_registers.process_operation_team_id").
		Joins("LEFT JOIN job_teams AS process_recovery_team ON process_recovery_team.id = system_registers.process_recovery_team_id").
		Joins("LEFT JOIN designations AS recovery_manager ON recovery_manager.id = system_registers.recovery_manager_id").
		Joins("LEFT JOIN business_continuity_duration_units AS rto_unit ON rto_unit.id = system_registers.rto_unit_id").
		Joins("LEFT JOIN business_continuity_duration_units AS mtpd_unit ON mtpd_unit.id = system_registers.mtpd_unit_id").
		Joins("LEFT JOIN business_continuity_duration_units AS rpo_unit ON rpo_unit.id = system_registers.rpo_unit_id").
		Joins("LEFT JOIN business_continuity_duration_units AS vendor_rto_unit ON vendor_rto_unit.id = system_registers.vendor_rto_unit_id").
		Joins("LEFT JOIN users AS created_user ON created_user.id = system_registers.created_by_id").
		Joins("LEFT JOIN users AS updated_user ON updated_user.id = system_registers.updated_by_id")

	// Apply filters
	for columnName, values := range params.Filters {
		if len(values) > 0 {
			// Special handling for impact_type_ids filter
			if columnName == "impact_type_ids" {
				query = query.Where("system_registers.id IN (SELECT system_register_id FROM system_register_impact_types WHERE impact_type_id IN (?))", values)
			} else {
				query = query.Where("system_registers."+columnName+" IN (?)", values)
			}
		}
	}

	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
			(system_registers.title_en LIKE ? OR
			system_registers.title_ar LIKE ? OR
			system_registers.service_description LIKE ? OR
			system_registers.modules_components LIKE ? OR
			system_registers.third_party_dependencies LIKE ? OR
			system_registers.outsourcing_partner LIKE ? OR
			system_registers.customers LIKE ? OR
			system_registers.it_processes LIKE ? OR
			system_registers.recovery_manager_phone LIKE ? OR
			system_registers.recovery_manager_email LIKE ? OR
			system_registers.recovery_location LIKE ? OR
			service_importance.title_en LIKE ? OR
			service_importance.title_ar LIKE ? OR
			responsible_it_team.title_en LIKE ? OR
			responsible_it_team.title_ar LIKE ? OR
			operation_location.title_en LIKE ? OR
			operation_location.title_ar LIKE ? OR
			business_owner.title_en LIKE ? OR
			business_owner.title_ar LIKE ? OR
			disruption_impact_level.title_en LIKE ? OR
			disruption_impact_level.title_ar LIKE ? OR
			it_service_owner.title_en LIKE ? OR
			it_service_owner.title_ar LIKE ? OR
			process_operation_team.title_en LIKE ? OR
			process_operation_team.title_ar LIKE ? OR
			process_recovery_team.title_en LIKE ? OR
			process_recovery_team.title_ar LIKE ? OR
			recovery_manager.title_en LIKE ? OR
			recovery_manager.title_ar LIKE ? OR
			rto_unit.title_en LIKE ? OR
			rto_unit.title_ar LIKE ? OR
			mtpd_unit.title_en LIKE ? OR
			mtpd_unit.title_ar LIKE ? OR
			rpo_unit.title_en LIKE ? OR
			rpo_unit.title_ar LIKE ? OR
			vendor_rto_unit.title_en LIKE ? OR
			vendor_rto_unit.title_ar LIKE ? OR
			created_user.first_name_en LIKE ? OR
			created_user.first_name_ar LIKE ? OR
			created_user.last_name_en LIKE ? OR
			created_user.last_name_ar LIKE ? OR
			created_user.email LIKE ? OR
			updated_user.first_name_en LIKE ? OR
			updated_user.first_name_ar LIKE ? OR
			updated_user.last_name_en LIKE ? OR
			updated_user.last_name_ar LIKE ? OR
			updated_user.email LIKE ?)`,
			keyword, keyword, keyword, keyword, keyword,
			keyword, keyword, keyword, keyword, keyword,
			keyword, keyword, keyword, keyword, keyword,
			keyword, keyword, keyword, keyword, keyword,
			keyword, keyword, keyword, keyword, keyword,
			keyword, keyword, keyword, keyword, keyword,
			keyword, keyword, keyword, keyword, keyword,
			keyword, keyword, keyword, keyword, keyword,
			keyword, keyword, keyword, keyword, keyword,
			keyword, keyword, keyword, keyword, keyword)
	}

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "desc"
	}

	query = query.Order(params.SortBy + " " + params.Sort)
	if !params.Export {
		offset := (params.Page - 1) * params.PageSize
		query = query.Offset(offset).Limit(params.PageSize)
	}

	err := query.Find(&systemRegisters).Error
	if err != nil {
		return nil, err
	}

	// Fetch impact types for each system register
	for _, systemRegister := range systemRegisters {
		var impactTypes []dtos.SystemRegisterImpactTypeListDTO

		err := repo.db.
			Table("system_register_impact_types").
			Select(`
				impact_categories.id,
				impact_categories.title_en,
				impact_categories.title_ar
			`).
			Joins("LEFT JOIN impact_categories ON impact_categories.id = system_register_impact_types.impact_type_id").
			Where("system_register_impact_types.system_register_id = ?", systemRegister.ID).
			Scan(&impactTypes).Error

		if err != nil {
			return nil, err
		}

		systemRegister.ImpactTypes = impactTypes
	}

	totalPages := 1
	if !params.Export {
		totalPages = int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))
	}

	response := &dtos.SystemRegisterListResponseDTO{
		Items: systemRegisters,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}

	return response, nil
}

func (repo *SystemRegisterRepository) GetByID(id uint) (*models.SystemRegister, error) {
	var systemRegister *models.SystemRegister
	result := repo.db.First(&systemRegister, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return systemRegister, nil
}

func (repo *SystemRegisterRepository) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.SystemRegister, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var systemRegister *models.SystemRegister

	result := db.
		Preload("ServiceImportance").
		Preload("ResponsibleITTeam").
		Preload("ResponsibleITTeam.TeamLeadDesignation").
		Preload("ResponsibleITTeam.JobTeamMembers").
		Preload("ResponsibleITTeam.JobTeamMembers.Designation").
		Preload("ResponsibleITTeam.TeamLeadUser").
		Preload("OperationLocation").
		Preload("BusinessOwner").
		Preload("SystemRegisterImpactTypes").
		Preload("SystemRegisterImpactTypes.ImpactType").
		Preload("DisruptionImpactLevel").
		Preload("RTOUnit").
		Preload("MTPDUnit").
		Preload("RPOUnit").
		Preload("ITServiceOwner").
		Preload("ProcessOperationTeam").
		Preload("ProcessOperationTeam.TeamLeadDesignation").
		Preload("ProcessOperationTeam.JobTeamMembers").
		Preload("ProcessOperationTeam.JobTeamMembers.Designation").
		Preload("ProcessOperationTeam.TeamLeadUser").
		Preload("ProcessRecoveryTeam").
		Preload("ProcessRecoveryTeam.TeamLeadDesignation").
		Preload("ProcessRecoveryTeam.JobTeamMembers").
		Preload("ProcessRecoveryTeam.JobTeamMembers.Designation").
		Preload("ProcessRecoveryTeam.TeamLeadUser").
		Preload("RecoveryManager").
		Preload("VendorRTOUnit").
		Preload("SystemRegisterSLA").
		Preload("SystemRegisterSLA.SLARegister").
		Preload("ProcessSystemRegisters.Process.ImpactAnalysisLevel").
		Preload("ProcessSystemRegisters.Process.Division").
		Preload("ProcessSystemRegisters.Process.ResponsibleUnit").
		Preload("ProcessSystemRegisters.Process.Section").
		Preload("ProcessSystemRegisters.Process.SubSection").
		Preload("ProcessSystemRegisters.Process.ResponsibleUser").
		Preload("ProcessSystemRegisters.Process.LifecycleStatus").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		First(&systemRegister, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return systemRegister, nil
}

func (repo *SystemRegisterRepository) Create(systemRegister *models.SystemRegister, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(systemRegister).Error; err != nil {
		return 0, err
	}
	return systemRegister.ID, nil
}

func (repo *SystemRegisterRepository) Update(systemRegister *models.SystemRegister, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(systemRegister).Error; err != nil {
		return err
	}
	return nil
}

func (repo *SystemRegisterRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var systemRegister *models.SystemRegister
	if err := db.First(&systemRegister, id).Error; err != nil {
		return err
	}
	return db.Delete(&models.SystemRegister{}, id).Error
}

func (repo *SystemRegisterRepository) CreateSLA(sla *models.SystemRegisterSLA, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(sla).Error; err != nil {
		return err
	}
	return nil
}

func (repo *SystemRegisterRepository) DeleteSLAs(systemRegisterID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Where("system_register_id = ?", systemRegisterID).Delete(&models.SystemRegisterSLA{}).Error; err != nil {
		return err
	}
	return nil
}

func (repo *SystemRegisterRepository) GetSLARegisterByID(slaRegisterID uint, tx ...utility.DatabaseTransaction) (*models.SLARegister, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var slaRegister *models.SLARegister
	result := db.
		Preload("RTOUnit").
		Preload("MTPDUnit").
		Preload("RPOUnit").
		First(&slaRegister, slaRegisterID)

	if result.Error != nil {
		return nil, result.Error
	}

	return slaRegister, nil
}

func (repo *SystemRegisterRepository) CreateSystemRegisterImpactType(impactType *models.SystemRegisterImpactType, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(impactType).Error; err != nil {
		return err
	}
	return nil
}

func (repo *SystemRegisterRepository) DeleteSystemRegisterImpactTypesBySystemRegisterID(systemRegisterID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Where("system_register_id = ?", systemRegisterID).Delete(&models.SystemRegisterImpactType{}).Error; err != nil {
		return err
	}
	return nil
}

func (repo *SystemRegisterRepository) GetImpactCategoryByID(id uint, tx ...utility.DatabaseTransaction) (*models.ImpactCategory, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var impactCategory *models.ImpactCategory
	result := db.First(&impactCategory, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return impactCategory, nil
}
