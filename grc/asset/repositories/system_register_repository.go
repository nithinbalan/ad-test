package repositories

import (
	"fmt"
	"jethur/grc/asset/dtos"
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/models"
	"jethur/utility"
	"strconv"
	"time"

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

	GenerateReferenceNumber() (string, error)

	CreateJobTitles(jobTitles []*models.SystemRegisterJobTitle, tx ...utility.DatabaseTransaction) error
	DeleteJobTitles(systemRegisterID uint, tx ...utility.DatabaseTransaction) error

	CreateOperationalLocations(locations []*models.SystemRegisterOperationalLocation, tx ...utility.DatabaseTransaction) error
	DeleteOperationalLocations(systemRegisterID uint, tx ...utility.DatabaseTransaction) error

	CreateRecoveryLocations(locations []*models.SystemRegisterRecoveryLocation, tx ...utility.DatabaseTransaction) error
	DeleteRecoveryLocations(systemRegisterID uint, tx ...utility.DatabaseTransaction) error

	CreateRecoveryNeededTools(tools []*models.SystemRegisterRecoveryNeededTool, tx ...utility.DatabaseTransaction) error
	DeleteRecoveryNeededTools(systemRegisterID uint, tx ...utility.DatabaseTransaction) error
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
			system_registers.reference_number,
			system_registers.title_en,
			system_registers.title_ar,
			system_registers.number_of_people_necessary,
			system_registers.third_party_dependencies,
			system_registers.outsourcing_partners,
			system_registers.customer_user,
			system_registers.recovery_manager_contact_details,
			business_owner_unit.title_en AS business_owner_unit_title_en,
			business_owner_unit.title_ar AS business_owner_unit_title_ar,
			it_responsible_unit.title_en AS it_responsible_unit_title_en,
			it_responsible_unit.title_ar AS it_responsible_unit_title_ar,
			process_operation_team.id AS process_operation_team_id,
			process_operation_team.title_en AS process_operation_team_title_en,
			process_operation_team.title_ar AS process_operation_team_title_ar,
			process_recovery_team.id AS process_recovery_team_id,
			process_recovery_team.title_en AS process_recovery_team_title_en,
			process_recovery_team.title_ar AS process_recovery_team_title_ar,
			recovery_manager.first_name_en AS recovery_manager_first_name_en,
			recovery_manager.first_name_ar AS recovery_manager_first_name_ar,
			recovery_manager.last_name_en AS recovery_manager_last_name_en,
			recovery_manager.last_name_ar AS recovery_manager_last_name_ar,
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
		Joins("LEFT JOIN organization_units AS business_owner_unit ON business_owner_unit.id = system_registers.business_owner_unit_id").
		Joins("LEFT JOIN organization_units AS it_responsible_unit ON it_responsible_unit.id = system_registers.it_responsible_unit_id").
		Joins("LEFT JOIN job_teams AS process_operation_team ON process_operation_team.id = system_registers.process_operation_team_id").
		Joins("LEFT JOIN job_teams AS process_recovery_team ON process_recovery_team.id = system_registers.process_recovery_team_id").
		Joins("LEFT JOIN users AS recovery_manager ON recovery_manager.id = system_registers.recovery_manager_id").
		Joins("LEFT JOIN users AS created_user ON created_user.id = system_registers.created_by_id").
		Joins("LEFT JOIN users AS updated_user ON updated_user.id = system_registers.updated_by_id")

	// Apply filters
	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where("system_registers."+columnName+" IN (?)", values)
		}
	}

	// Apply search keyword if exists
	if params.Keyword != "" {
		query = query.Where(`
			(system_registers.reference_number LIKE ? OR
			system_registers.title_en LIKE ? OR
			system_registers.title_ar LIKE ? OR
			system_registers.third_party_dependencies LIKE ? OR
			system_registers.outsourcing_partners LIKE ? OR
			system_registers.customer_user LIKE ? OR
			business_owner_unit.title_en LIKE ? OR
			business_owner_unit.title_ar LIKE ? OR
			it_responsible_unit.title_en LIKE ? OR
			it_responsible_unit.title_ar LIKE ?)`,
			"%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%",
			"%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%")
	}

	// Calculate total count for pagination
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	// Validate sort order to prevent SQL injection (defense in depth)
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "desc"
	}

	// Apply sorting and pagination
	query = query.Order(params.SortBy + " " + params.Sort)
	if !params.Export {
		offset := (params.Page - 1) * params.PageSize
		query = query.Offset(offset).Limit(params.PageSize)
	}

	// Execute query
	err := query.Find(&systemRegisters).Error
	if err != nil {
		return nil, err
	}

	// Fetch related data for each system register
	for _, systemRegister := range systemRegisters {
		// Get Job Titles
		var jobTitles []*dtos.SystemRegisterJobTitleListDTO
		err := repo.db.Table("system_register_job_titles AS srjt").
			Select(`
				srjt.id,
				d.title_en AS title_en,
				d.title_ar AS title_ar
			`).
			Joins("INNER JOIN designations AS d ON d.id = srjt.designation_id").
			Where("srjt.system_register_id = ?", systemRegister.ID).
			Scan(&jobTitles).Error
		if err == nil {
			systemRegister.JobTitles = jobTitles
		}

		// Get Operational Locations
		var operationalLocations []*dtos.SystemRegisterLocationListDTO
		err = repo.db.Table("system_register_operational_locations AS srol").
			Select(`
				srol.id,
				l.title_en AS title_en,
				l.title_ar AS title_ar
			`).
			Joins("INNER JOIN locations AS l ON l.id = srol.operational_location_id").
			Where("srol.system_register_id = ?", systemRegister.ID).
			Scan(&operationalLocations).Error
		if err == nil {
			systemRegister.OperationalLocations = operationalLocations
		}

		// Get Recovery Locations
		var recoveryLocations []*dtos.SystemRegisterLocationListDTO
		err = repo.db.Table("system_register_recovery_locations AS srrl").
			Select(`
				srrl.id,
				l.title_en AS title_en,
				l.title_ar AS title_ar
			`).
			Joins("INNER JOIN locations AS l ON l.id = srrl.recovery_location_id").
			Where("srrl.system_register_id = ?", systemRegister.ID).
			Scan(&recoveryLocations).Error
		if err == nil {
			systemRegister.RecoveryLocations = recoveryLocations
		}

		// Get Recovery Needed Tools
		var recoveryNeededTools []*dtos.SystemRegisterRecoveryToolListDTO
		err = repo.db.Table("system_register_recovery_needed_tools AS srrnt").
			Select(`
				srrnt.id,
				srrt.title_en AS title_en,
				srrt.title_ar AS title_ar
			`).
			Joins("INNER JOIN system_register_recovery_tools AS srrt ON srrt.id = srrnt.recovery_needed_tool_id").
			Where("srrnt.system_register_id = ?", systemRegister.ID).
			Scan(&recoveryNeededTools).Error
		if err == nil {
			systemRegister.RecoveryNeededTools = recoveryNeededTools
		}
	}

	// Calculate total pages
	totalPages := 1
	if !params.Export {
		totalPages = int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))
	}

	// Prepare response
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
		Preload("BusinessOwnerUnit").
		Preload("ITResponsibleUnit").
		Preload("ProcessOperationTeam").
		Preload("ProcessOperationTeam.TeamLeadDesignation").
		Preload("ProcessOperationTeam.JobTeamMembers").
		Preload("ProcessOperationTeam.JobTeamMembers.Designation").
		Preload("ProcessRecoveryTeam").
		Preload("ProcessRecoveryTeam.TeamLeadDesignation").
		Preload("ProcessRecoveryTeam.JobTeamMembers").
		Preload("ProcessRecoveryTeam.JobTeamMembers.Designation").
		Preload("RecoveryManager").
		Preload("SystemRegisterJobTitles").
		Preload("SystemRegisterJobTitles.Designation").
		Preload("SystemRegisterOperationalLocations").
		Preload("SystemRegisterOperationalLocations.OperationalLocation").
		Preload("SystemRegisterRecoveryLocations").
		Preload("SystemRegisterRecoveryLocations.RecoveryLocation").
		Preload("SystemRegisterRecoveryNeededTools").
		Preload("SystemRegisterRecoveryNeededTools.RecoveryNeededTool").
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

func (repo *SystemRegisterRepository) GenerateReferenceNumber() (string, error) {
	var lastReferenceNumber string
	err := repo.db.Model(&models.SystemRegister{}).
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
	newReferenceNumber := fmt.Sprintf("%s-%d-%04d", "SYS", currentYear, lastReferenceID+1)

	return newReferenceNumber, nil
}

func (repo *SystemRegisterRepository) CreateJobTitles(jobTitles []*models.SystemRegisterJobTitle, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(&jobTitles).Error; err != nil {
		return err
	}
	return nil
}

func (repo *SystemRegisterRepository) DeleteJobTitles(systemRegisterID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("system_register_id = ?", systemRegisterID).Delete(&models.SystemRegisterJobTitle{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *SystemRegisterRepository) CreateOperationalLocations(locations []*models.SystemRegisterOperationalLocation, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(&locations).Error; err != nil {
		return err
	}
	return nil
}

func (repo *SystemRegisterRepository) DeleteOperationalLocations(systemRegisterID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("system_register_id = ?", systemRegisterID).Delete(&models.SystemRegisterOperationalLocation{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *SystemRegisterRepository) CreateRecoveryLocations(locations []*models.SystemRegisterRecoveryLocation, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(&locations).Error; err != nil {
		return err
	}
	return nil
}

func (repo *SystemRegisterRepository) DeleteRecoveryLocations(systemRegisterID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("system_register_id = ?", systemRegisterID).Delete(&models.SystemRegisterRecoveryLocation{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *SystemRegisterRepository) CreateRecoveryNeededTools(tools []*models.SystemRegisterRecoveryNeededTool, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(&tools).Error; err != nil {
		return err
	}
	return nil
}

func (repo *SystemRegisterRepository) DeleteRecoveryNeededTools(systemRegisterID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("system_register_id = ?", systemRegisterID).Delete(&models.SystemRegisterRecoveryNeededTool{}).Error; err != nil {
		return err
	}

	return nil
}
