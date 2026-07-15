package repositories

import (
	"fmt"
	businessContinuityDTOs "jethur/grc/business-continuity/dtos"
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/models"
	"jethur/utility"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type BusinessContinuityPlanRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(params *businessContinuityDTOs.BusinessContinuityPlanListQueryParamsDTO) (*businessContinuityDTOs.BusinessContinuityPlanListResponseDTO, error)
	GetByID(id uint) (*models.BusinessContinuityPlan, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.BusinessContinuityPlan, error)
	Create(businessContinuityPlan *models.BusinessContinuityPlan, tx ...utility.DatabaseTransaction) (uint, error)
	Update(businessContinuityPlan *models.BusinessContinuityPlan, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error
	GenerateReferenceNumber() (string, error)

	GetBusinessContinuityPlanByBusinessFunctionID(businessFunctionId uint) (*models.BusinessContinuityPlan, error)
	CreateBcpProcesses(bcpProcess *models.BcpProcess, tx ...utility.DatabaseTransaction) (uint, error)
	DeleteBcpProcesses(bcpProcessIds []*uint, tx ...utility.DatabaseTransaction) error
	CreateDocuments(docs []*models.BcpDocument, tx ...utility.DatabaseTransaction) error
	GetDocumentByID(id uint) (*models.BcpDocument, error)
	DeleteDocuments(docIds []*uint, tx ...utility.DatabaseTransaction) error

	GetBcpRisksByBcpId(bcpId uint) ([]models.BcpRisk, error)
	GetTestAndExercisesByBcpID(bcpID uint) ([]*models.TestAndExercise, error)

	// BCP Process child table snapshot methods
	GetProcessWithChildrenByID(processID uint) (*models.Process, error)
	CreateBcpProcessExecutionRecord(record *models.BcpProcessExecutionRecord, tx ...utility.DatabaseTransaction) error
	CreateBcpProcessManualProcedureForm(form *models.BcpProcessManualProcedureForm, tx ...utility.DatabaseTransaction) error
	CreateBcpProcessSystemRegister(sr *models.BcpProcessSystemRegister, tx ...utility.DatabaseTransaction) error
	CreateBcpProcessInternalDependency(dep *models.BcpProcessInternalDependency, tx ...utility.DatabaseTransaction) error
	CreateBcpProcessExternalDependency(dep *models.BcpProcessExternalDependency, tx ...utility.DatabaseTransaction) error
	CreateBcpProcessEmployeeAnalysis(analysis *models.BcpProcessEmployeeAnalysis, tx ...utility.DatabaseTransaction) (uint, error)
	CreateBcpProcessEmployeeAnalysisEmployee(employee *models.BcpProcessEmployeeAnalysisEmployee, tx ...utility.DatabaseTransaction) (uint, error)
	CreateBcpProcessEmployeeAnalysisEmployeeTool(tool *models.BcpProcessEmployeeAnalysisEmployeeTool, tx ...utility.DatabaseTransaction) error
	CreateBcpProcessStaffAllocation(alloc *models.BcpProcessStaffAllocation, tx ...utility.DatabaseTransaction) (uint, error)
	CreateBcpProcessAllocationOutageBucket(bucket *models.BcpProcessAllocationOutageBucket, tx ...utility.DatabaseTransaction) (uint, error)
	CreateBcpProcessAllocationOutageBucketUser(user *models.BcpProcessAllocationOutageBucketUser, tx ...utility.DatabaseTransaction) error
	DeleteBcpProcessChildrenByBcpProcessIDs(bcpProcessIDs []uint, tx ...utility.DatabaseTransaction) error
	GetAllEquipmentTools(tx ...utility.DatabaseTransaction) ([]models.ProcessEquipmentTool, error)
	GetAllOutageBuckets(tx ...utility.DatabaseTransaction) ([]models.BCMBIAOutageBucket, error)
}

type BusinessContinuityPlanRepository struct {
	db *gorm.DB
}

func NewBusinessContinuityPlanRepository(db *gorm.DB) *BusinessContinuityPlanRepository {
	return &BusinessContinuityPlanRepository{
		db: db,
	}
}

// BeginTransaction starts a new transaction.
func (repo *BusinessContinuityPlanRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *BusinessContinuityPlanRepository) GetAll(params *businessContinuityDTOs.BusinessContinuityPlanListQueryParamsDTO) (*businessContinuityDTOs.BusinessContinuityPlanListResponseDTO, error) {
	var businessContinuityPlans []*businessContinuityDTOs.BusinessContinuityPlanListDTO
	var totalCount int64

	query := repo.db.Table("business_continuity_plans").
		Select(`
            business_continuity_plans.*,
			business_functions.title AS business_function_title,
			bcp_statuses.title_en AS bcp_status_title_en,
			bcp_statuses.title_ar AS bcp_status_title_ar,
			bcp_statuses.color_code AS bcp_status_color_code,
			bcp_statuses.type AS bcp_status_type,
			rto_time_unit.title_en AS rto_unit_title_en,
			rto_time_unit.title_ar AS rto_unit_title_ar,
			mtpd_time_unit.title_en AS mtpd_unit_title_en,
			mtpd_time_unit.title_ar AS mtpd_unit_title_ar,
			rpo_time_unit.title_en AS rpo_unit_title_en,
			rpo_time_unit.title_ar AS rpo_unit_title_ar,
			responsible_user.id AS responsible_user_id,
			responsible_user.first_name_en AS responsible_user_first_name_en,
            responsible_user.first_name_ar AS responsible_user_first_name_ar,
            responsible_user.last_name_en AS responsible_user_last_name_en,
            responsible_user.last_name_ar AS responsible_user_last_name_ar,
			responsible_user.initial AS responsible_user_initial,
            responsible_user.color AS responsible_user_color,
            responsible_user.email AS responsible_user_email,
            responsible_user.profile_image_path AS responsible_user_profile_image_path,
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
		Joins(`
		    LEFT JOIN business_functions ON business_functions.id = business_continuity_plans.business_function_id
			LEFT JOIN bcp_statuses ON bcp_statuses.id = business_continuity_plans.bcp_status_id
			LEFT JOIN users AS responsible_user ON responsible_user.id = business_continuity_plans.responsible_user_id
			LEFT JOIN time_units AS rto_time_unit ON rto_time_unit.id = business_continuity_plans.rto_unit_id
			LEFT JOIN time_units AS mtpd_time_unit ON mtpd_time_unit.id = business_continuity_plans.mtpd_unit_id
			LEFT JOIN time_units AS rpo_time_unit ON rpo_time_unit.id = business_continuity_plans.rpo_unit_id
			LEFT JOIN users AS created_user ON created_user.id = business_continuity_plans.created_by_id
			LEFT JOIN users AS updated_user ON updated_user.id = business_continuity_plans.updated_by_id
        `)
	// Apply dynamic filters to the query
	joinedProcessTables := false
	for columnName, values := range params.Filters {
		if len(values) > 0 {
			if columnName == "bcp_status_id" {
				query = query.Where("bcp_statuses.id IN (?)", values)
			} else {
				query = query.Where("business_continuity_plans."+columnName+" IN (?)", values)
			}
		}
	}
	// Apply quick filters
	for _, filter := range params.QuickFilter {
		if !joinedProcessTables {
			query = query.
				Joins(`LEFT JOIN bcp_processes ON bcp_processes.business_continuity_plan_id = business_continuity_plans.id`).
				Joins(`LEFT JOIN processes AS bcp_proc ON bcp_proc.id = bcp_processes.process_id`).
				Joins(`LEFT JOIN bcm_impact_analysis_levels ON bcm_impact_analysis_levels.id = bcp_proc.impact_analysis_level_id`)
			joinedProcessTables = true
		}
		switch filter {
		case "is_critical":
			query = query.Where("bcp_proc.impact_analysis_level_id IS NOT NULL AND bcm_impact_analysis_levels.type != ?", "non_applicable")
		case "is_not_critical":
			query = query.Where("(bcp_proc.impact_analysis_level_id IS NULL OR bcm_impact_analysis_levels.type = ?)", "non_applicable")
		}
	}
	// Calculate total count for pagination
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}
	// Apply search keyword if exists
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
			(business_continuity_plans.reference_number LIKE ? OR
			business_functions.title LIKE ? OR
			created_user.first_name_en LIKE ? OR
			created_user.first_name_ar LIKE ? OR
			created_user.last_name_en LIKE ? OR
			created_user.last_name_ar LIKE ? OR
    		updated_user.first_name_en LIKE ? OR
			updated_user.first_name_ar LIKE ? OR
			updated_user.last_name_en LIKE ? OR
			updated_user.last_name_ar LIKE ?)`,
			keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword,
			keyword)
	}

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	// Sort by the specified column and order
	if params.SortBy == "" {
		params.SortBy = "business_continuity_plans.created_at"
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "desc"
	}

	query = query.Order(params.SortBy + " " + params.Sort)

	if !params.Export {
		offset := (params.Page - 1) * params.PageSize
		query = query.Offset(offset).Limit(params.PageSize)
	}

	err := query.Find(&businessContinuityPlans).Error
	if err != nil {
		return nil, err
	}

	totalPages := 1
	if !params.Export {
		totalPages = int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))
	}

	response := &businessContinuityDTOs.BusinessContinuityPlanListResponseDTO{
		Items: businessContinuityPlans,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}
	return response, nil
}

func (repo *BusinessContinuityPlanRepository) GetByID(id uint) (*models.BusinessContinuityPlan, error) {
	var businessContinuityPlan *models.BusinessContinuityPlan

	result := repo.db.First(&businessContinuityPlan, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return businessContinuityPlan, nil

}

func (repo *BusinessContinuityPlanRepository) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.BusinessContinuityPlan, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var businessContinuityPlan *models.BusinessContinuityPlan

	result := db.
		Preload("BusinessFunction").
		Preload("BusinessFunction.SeverityLevel").
		Preload("BusinessFunction.ReviewFrequency").
		Preload("BusinessFunction.ResponsibleUser").
		Preload("ResponsibleUser").
		Preload("ResponsibleUnit").
		Preload("Documents").
		Preload("BcpStatus").
		Preload("RtoUnit").
		Preload("MtpdUnit").
		Preload("RpoUnit").
		Preload("Processes").
		Preload("Processes.ImpactAnalysisLevel").
		Preload("Processes.RtoUnit").
		Preload("Processes.MtpdUnit").
		Preload("Processes.RpoUnit").
		Preload("Processes.RspUnit").
		Preload("Processes.Process").
		Preload("Processes.Process.ResponsibleUser").
		Preload("Processes.ExecutionRecords").
		Preload("Processes.ExecutionRecords.BCMExecutionRecord").
		Preload("Processes.ExecutionRecords.BCMExecutionRecord.OwnerEntity").
		Preload("Processes.ExecutionRecords.BCMExecutionRecord.Location").
		Preload("Processes.ExecutionRecords.BCMExecutionRecord.RecordType").
		Preload("Processes.ManualProcedureForms").
		Preload("Processes.ManualProcedureForms.BCMManualProcedureForm").
		Preload("Processes.ManualProcedureForms.BCMManualProcedureForm.OwnerEntity").
		Preload("Processes.ManualProcedureForms.BCMManualProcedureForm.Location").
		Preload("Processes.SystemRegisters").
		Preload("Processes.SystemRegisters.SystemRegister").
		Preload("Processes.SystemRegisters.SystemRegister.RTOUnit").
		Preload("Processes.SystemRegisters.SystemRegister.MTPDUnit").
		Preload("Processes.SystemRegisters.SystemRegister.RPOUnit").
		Preload("Processes.SystemRegisters.SystemRegister.RecoveryManager").
		Preload("Processes.InternalDependencies").
		Preload("Processes.InternalDependencies.SupportingUnit").
		Preload("Processes.InternalDependencies.DependencyTier").
		Preload("Processes.ExternalDependencies").
		Preload("Processes.ExternalDependencies.EntityType").
		Preload("Processes.ExternalDependencies.DependencyTier").
		Preload("Processes.ExternalDependencies.SLARegister").
		Preload("Processes.ExternalDependencies.VendorRTOUnit").
		Preload("Processes.EmployeeAnalyses").
		Preload("Processes.EmployeeAnalyses.Designation").
		Preload("Processes.EmployeeAnalyses.Employees").
		Preload("Processes.EmployeeAnalyses.Employees.User").
		Preload("Processes.EmployeeAnalyses.Employees.Tools").
		Preload("Processes.EmployeeAnalyses.Employees.Tools.EquipmentTool").
		Preload("Processes.StaffAllocations").
		Preload("Processes.StaffAllocations.BcpProcessEmployeeAnalysis").
		Preload("Processes.StaffAllocations.OutageBuckets.OutageBucket").
		Preload("Processes.StaffAllocations.OutageBuckets.Users.User").

		// Preload("BcpRecoveryPlans").
		// Preload("BcpRecoveryPlans.SeverityLevel").
		// Preload("BcpRecoveryPlans.BcpProcess").
		// Preload("BcpRecoveryPlans.BcpProcess.Process").
		// Preload("BcpRecoveryPlans.TargetCompletionUnit").
		// Preload("BcpCommunicationPlans").
		// Preload("BcpCommunicationPlans.BcpProcess").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		First(&businessContinuityPlan, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return businessContinuityPlan, nil
}

func (repo *BusinessContinuityPlanRepository) Create(businessContinuityPlan *models.BusinessContinuityPlan, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(businessContinuityPlan).Error; err != nil {
		return 0, err
	}

	return businessContinuityPlan.ID, nil
}

func (repo *BusinessContinuityPlanRepository) Update(businessContinuityPlan *models.BusinessContinuityPlan, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(businessContinuityPlan).Error; err != nil {
		return err
	}

	return nil
}

func (repo *BusinessContinuityPlanRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	// Delete the Business Function
	if err := db.Delete(&models.BusinessContinuityPlan{}, id).Error; err != nil {
		return err
	}

	return nil
}

func (repo *BusinessContinuityPlanRepository) GenerateReferenceNumber() (string, error) {
	var lastReferenceNumber string
	err := repo.db.Model(&models.BusinessContinuityPlan{}).
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
	newReferenceNumber := fmt.Sprintf("%s-%d-%04d", "BCP", currentYear, lastReferenceID+1)

	return newReferenceNumber, nil

}

func (repo *BusinessContinuityPlanRepository) GetBusinessContinuityPlanByBusinessFunctionID(businessFunctionId uint) (*models.BusinessContinuityPlan, error) {
	var businessContinuityPlan models.BusinessContinuityPlan

	// Query the database to find if an businessContinuityPlan exists with the given BusinessFunctionId
	result := repo.db.Where("business_function_id = ?", businessFunctionId).First(&businessContinuityPlan)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			// Return nil if no businessContinuityPlan is found (not an error, meaning it's safe to create a new businessContinuityPlan)
			return nil, nil
		}
		// Return an error if something went wrong during the query
		return nil, result.Error
	}

	// Return the found audit if it exists
	return &businessContinuityPlan, nil
}

func (repo *BusinessContinuityPlanRepository) CreateBcpProcesses(bcpProcess *models.BcpProcess, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(bcpProcess).Error; err != nil {
		return 0, err
	}
	return bcpProcess.ID, nil
}

func (repo *BusinessContinuityPlanRepository) DeleteBcpProcesses(bcpProcessIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	// Delete processes by IDs
	if err := db.Where("id IN (?)", bcpProcessIds).Delete(&models.BcpProcess{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *BusinessContinuityPlanRepository) CreateDocuments(docs []*models.BcpDocument, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	// Save all documents in a single transaction
	if err := db.Create(&docs).Error; err != nil {
		return err
	}
	return nil
}

func (repo *BusinessContinuityPlanRepository) GetDocumentByID(id uint) (*models.BcpDocument, error) {
	var document models.BcpDocument
	if err := repo.db.Where("id = ?", id).First(&document).Error; err != nil {
		return nil, err
	}
	return &document, nil
}

func (repo *BusinessContinuityPlanRepository) DeleteDocuments(docIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	// Delete documents by IDs
	if err := db.Where("id IN (?)", docIds).Delete(&models.BcpDocument{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *BusinessContinuityPlanRepository) GetProcessWithChildrenByID(processID uint) (*models.Process, error) {
	var process models.Process
	err := repo.db.
		Preload("ProcessExecutionRecords").
		Preload("ProcessManualProcedureForms").
		Preload("ProcessSystemRegister").
		Preload("InternalDependencies").
		Preload("ExternalDependencies").
		Preload("ProcessEmployeeAnalysis").
		Preload("ProcessEmployeeAnalysis.Employees").
		Preload("ProcessEmployeeAnalysis.Employees.Tools").
		Preload("BCMStaffAllocations").
		Preload("BCMStaffAllocations.OutageBuckets").
		Preload("BCMStaffAllocations.OutageBuckets.Users").
		First(&process, processID).Error
	return &process, err
}

func (repo *BusinessContinuityPlanRepository) CreateBcpProcessExecutionRecord(record *models.BcpProcessExecutionRecord, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Create(record).Error
}

func (repo *BusinessContinuityPlanRepository) CreateBcpProcessManualProcedureForm(form *models.BcpProcessManualProcedureForm, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Create(form).Error
}

func (repo *BusinessContinuityPlanRepository) CreateBcpProcessSystemRegister(sr *models.BcpProcessSystemRegister, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Create(sr).Error
}

func (repo *BusinessContinuityPlanRepository) CreateBcpProcessInternalDependency(dep *models.BcpProcessInternalDependency, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Create(dep).Error
}

func (repo *BusinessContinuityPlanRepository) CreateBcpProcessExternalDependency(dep *models.BcpProcessExternalDependency, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Create(dep).Error
}

func (repo *BusinessContinuityPlanRepository) CreateBcpProcessEmployeeAnalysis(analysis *models.BcpProcessEmployeeAnalysis, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	err := db.Create(analysis).Error
	return analysis.ID, err
}

func (repo *BusinessContinuityPlanRepository) CreateBcpProcessEmployeeAnalysisEmployee(employee *models.BcpProcessEmployeeAnalysisEmployee, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	err := db.Create(employee).Error
	return employee.ID, err
}

func (repo *BusinessContinuityPlanRepository) CreateBcpProcessEmployeeAnalysisEmployeeTool(tool *models.BcpProcessEmployeeAnalysisEmployeeTool, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Create(tool).Error
}

func (repo *BusinessContinuityPlanRepository) CreateBcpProcessStaffAllocation(alloc *models.BcpProcessStaffAllocation, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	err := db.Create(alloc).Error
	return alloc.ID, err
}

func (repo *BusinessContinuityPlanRepository) CreateBcpProcessAllocationOutageBucket(bucket *models.BcpProcessAllocationOutageBucket, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	err := db.Create(bucket).Error
	return bucket.ID, err
}

func (repo *BusinessContinuityPlanRepository) CreateBcpProcessAllocationOutageBucketUser(user *models.BcpProcessAllocationOutageBucketUser, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Create(user).Error
}

// DeleteBcpProcessChildrenByBcpProcessIDs deletes all child table records for the given bcp_process IDs
// in the correct order to respect FK constraints.
func (repo *BusinessContinuityPlanRepository) DeleteBcpProcessChildrenByBcpProcessIDs(bcpProcessIDs []uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if len(bcpProcessIDs) == 0 {
		return nil
	}

	// Delete outage bucket users
	if err := db.Where(
		"bcp_process_allocation_outage_bucket_id IN (SELECT id FROM bcp_process_allocation_outage_buckets WHERE bcp_process_allocation_id IN (SELECT id FROM bcp_process_staff_allocations WHERE bcp_process_id IN (?)))",
		bcpProcessIDs).Delete(&models.BcpProcessAllocationOutageBucketUser{}).Error; err != nil {
		return err
	}
	// Delete outage buckets
	if err := db.Where(
		"bcp_process_allocation_id IN (SELECT id FROM bcp_process_staff_allocations WHERE bcp_process_id IN (?))",
		bcpProcessIDs).Delete(&models.BcpProcessAllocationOutageBucket{}).Error; err != nil {
		return err
	}
	// Delete staff allocations
	if err := db.Where("bcp_process_id IN (?)", bcpProcessIDs).Delete(&models.BcpProcessStaffAllocation{}).Error; err != nil {
		return err
	}
	// Delete employee analysis employee tools
	if err := db.Where(
		"bcp_process_employee_analysis_employee_id IN (SELECT id FROM bcp_process_employee_analysis_employees WHERE bcp_process_employee_analysis_id IN (SELECT id FROM bcp_process_employee_analysis WHERE bcp_process_id IN (?)))",
		bcpProcessIDs).Delete(&models.BcpProcessEmployeeAnalysisEmployeeTool{}).Error; err != nil {
		return err
	}
	// Delete employee analysis employees
	if err := db.Where(
		"bcp_process_employee_analysis_id IN (SELECT id FROM bcp_process_employee_analysis WHERE bcp_process_id IN (?))",
		bcpProcessIDs).Delete(&models.BcpProcessEmployeeAnalysisEmployee{}).Error; err != nil {
		return err
	}
	// Delete employee analysis
	if err := db.Where("bcp_process_id IN (?)", bcpProcessIDs).Delete(&models.BcpProcessEmployeeAnalysis{}).Error; err != nil {
		return err
	}
	// Delete execution records
	if err := db.Where("bcp_process_id IN (?)", bcpProcessIDs).Delete(&models.BcpProcessExecutionRecord{}).Error; err != nil {
		return err
	}
	// Delete manual procedure forms
	if err := db.Where("bcp_process_id IN (?)", bcpProcessIDs).Delete(&models.BcpProcessManualProcedureForm{}).Error; err != nil {
		return err
	}
	// Delete system registers
	if err := db.Where("bcp_process_id IN (?)", bcpProcessIDs).Delete(&models.BcpProcessSystemRegister{}).Error; err != nil {
		return err
	}
	// Delete internal dependencies
	if err := db.Where("bcp_process_id IN (?)", bcpProcessIDs).Delete(&models.BcpProcessInternalDependency{}).Error; err != nil {
		return err
	}
	// Delete external dependencies
	if err := db.Where("bcp_process_id IN (?)", bcpProcessIDs).Delete(&models.BcpProcessExternalDependency{}).Error; err != nil {
		return err
	}
	return nil
}

func (repo *BusinessContinuityPlanRepository) GetBcpRisksByBcpId(bcpId uint) ([]models.BcpRisk, error) {
	var bcpRisks []models.BcpRisk
	if err := repo.db.Where("business_continuity_plan_id = ?", bcpId).Find(&bcpRisks).Error; err != nil {
		return nil, err
	}

	return bcpRisks, nil
}

func (repo *BusinessContinuityPlanRepository) GetTestAndExercisesByBcpID(bcpID uint) ([]*models.TestAndExercise, error) {
	var testAndExercises []*models.TestAndExercise

	result := repo.db.
		Preload("TestAndExerciseTestResult").
		Preload("TestAndExerciseTestResult.TestResultSystemRegisters").
		Preload("TestAndExerciseTestResult.TestResultSystemRegisters.TestResultOption").
		Where("business_continuity_plan_id = ?", bcpID).
		Find(&testAndExercises)

	if result.Error != nil {
		return nil, result.Error
	}

	return testAndExercises, nil
}

func (repo *BusinessContinuityPlanRepository) GetAllEquipmentTools(tx ...utility.DatabaseTransaction) ([]models.ProcessEquipmentTool, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var equipmentTools []models.ProcessEquipmentTool
	result := db.Where("status = ?", "active").Find(&equipmentTools)

	if result.Error != nil {
		return nil, result.Error
	}

	return equipmentTools, nil
}

func (repo *BusinessContinuityPlanRepository) GetAllOutageBuckets(tx ...utility.DatabaseTransaction) ([]models.BCMBIAOutageBucket, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var outageBuckets []models.BCMBIAOutageBucket
	result := db.Where("status = ?", "active").
		Preload("DurationUnit").
		Order("[order] ASC").
		Find(&outageBuckets)

	if result.Error != nil {
		return nil, result.Error
	}

	return outageBuckets, nil
}
