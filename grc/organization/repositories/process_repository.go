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
	GetDocumentByID(id uint) (*models.Document, error)

	//lifecycle status
	GetIDsByLifecycleStatus(ids []uint, status string) ([]uint, error)
	UpdateLifecycleStatusForIDs(ids []uint, statusId uint) ([]uint, error)

	GetControlSummary(processID uint) ([]models.ProcessControl, error)

	// Process BCM associations
	GetProcessExecutionRecordsByProcessID(processID uint, tx ...utility.DatabaseTransaction) ([]*models.ProcessExecutionRecord, error)
	CreateProcessExecutionRecord(record *models.ProcessExecutionRecord, tx ...utility.DatabaseTransaction) error
	DeleteProcessExecutionRecordsByProcessID(processID uint, tx ...utility.DatabaseTransaction) error

	GetProcessManualProcedureFormsByProcessID(processID uint, tx ...utility.DatabaseTransaction) ([]*models.ProcessManualProcedureForm, error)
	CreateProcessManualProcedureForm(form *models.ProcessManualProcedureForm, tx ...utility.DatabaseTransaction) error
	DeleteProcessManualProcedureFormsByProcessID(processID uint, tx ...utility.DatabaseTransaction) error

	GetProcessSystemRegistersByProcessID(processID uint, tx ...utility.DatabaseTransaction) ([]*models.ProcessSystemRegister, error)
	CreateProcessSystemRegister(record *models.ProcessSystemRegister, tx ...utility.DatabaseTransaction) error
	DeleteProcessSystemRegistersByProcessID(processID uint, tx ...utility.DatabaseTransaction) error

	// Associated entity validation helpers
	GetBCMProcessTypeByID(id uint, tx ...utility.DatabaseTransaction) (*models.BCMProcessType, error)
	GetBCMExecutionRecordByID(id uint, tx ...utility.DatabaseTransaction) (*models.BCMExecutionRecord, error)
	GetBCMManualProcedureFormByID(id uint, tx ...utility.DatabaseTransaction) (*models.BCMManualProcedureForm, error)
	GetSystemRegisterByID(id uint, tx ...utility.DatabaseTransaction) (*models.SystemRegister, error)

	// GetAll helpers for import template dropdowns
	GetAllBCMProcessTypes() ([]*models.BCMProcessType, error)
	GetAllBCMExecutionRecords() ([]*models.BCMExecutionRecord, error)
	GetAllBCMManualProcedureForms() ([]*models.BCMManualProcedureForm, error)
	GetAllSystemRegisters() ([]*models.SystemRegister, error)
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
		processes.rto_from_value as rto_from_value,
		processes.rto_to_value as rto_to_value,
		processes.rto_unit_id as rto_unit_id,
		rto_bc_unit.title_en AS rto_unit_title_en,
		rto_bc_unit.title_ar AS rto_unit_title_ar,
		processes.rto_condition as rto_condition,
		processes.mtpd_from_value as mtpd_from_value,
		processes.mtpd_to_value as mtpd_to_value,
		processes.mtpd_unit_id as mtpd_unit_id,
		mtpd_bc_unit.title_en AS mtpd_unit_title_en,
		mtpd_bc_unit.title_ar AS mtpd_unit_title_ar,
		processes.mtpd_condition as mtpd_condition,
		processes.rpo_from_value as rpo_from_value,
		processes.rpo_to_value as rpo_to_value,
		processes.rpo_unit_id as rpo_unit_id,
		rpo_bc_unit.title_en AS rpo_unit_title_en,
		rpo_bc_unit.title_ar AS rpo_unit_title_ar,
		processes.rpo_condition as rpo_condition,
		processes.rsp_from_value as rsp_from_value,
		processes.rsp_to_value as rsp_to_value,
		processes.rsp_unit_id as rsp_unit_id,
		rsp_bc_unit.title_en AS rsp_unit_title_en,
		rsp_bc_unit.title_ar AS rsp_unit_title_ar,
		processes.rsp_condition as rsp_condition,
		processes.impact_analysis_level_id AS impact_analysis_level_id,
		bcm_impact_analysis_levels.title_en AS impact_analysis_level_title_en,
		bcm_impact_analysis_levels.title_ar AS impact_analysis_level_title_ar,
		bcm_impact_analysis_levels.color_code AS impact_analysis_level_color_code,
		bcm_impact_analysis_levels.type AS impact_analysis_level_type,
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
		process_criticality_levels.id AS process_criticality_level_id,
		process_criticality_levels.title_en AS process_criticality_level_title_en,
		process_criticality_levels.title_ar AS process_criticality_level_title_ar,
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
		risk_matrix_risk_rating_values.value AS risk_rating_value,
		risk_levels.title_en AS risk_level_title_en,
        risk_levels.title_ar AS risk_level_title_ar,
		risk_levels.color_code AS risk_level_color_code,
		risk_levels.type AS risk_level_type,
		processes.division_id,
		divisions.title_en AS division_title_en,
		divisions.title_ar AS division_title_ar,
		processes.is_automatic,
		processes.is_manual,
		processes.process_type_id,
		bcm_process_types.title_en AS process_type_title_en,
		bcm_process_types.title_ar AS process_type_title_ar,
		processes.location_id,
		locations.title_en AS location_title_en,
		locations.title_ar AS location_title_ar,
		processes.alternative_location,
		processes.rush_hour,
		processes.review_frequency_id,
		review_frequencies.title_en AS review_frequency_title_en,
		review_frequencies.title_ar AS review_frequency_title_ar,
		processes.is_can_be_done_manually,
		processes.time_to_perform_procedure_manually,
		processes.steps_to_perform_procedure_manually,
		processes.requirement_to_perform_procedure_manually,
		bia_processes.workflow_status_id AS workflow_status_id,
		workflow_statuses.title_en AS workflow_status_title_en,
		workflow_statuses.title_ar AS workflow_status_title_ar,
		workflow_statuses.type AS workflow_status_type,
		workflow_statuses.color_code AS workflow_status_color_code,
		bia_processes.approved_by_id AS approved_by_id,
		approved_by.first_name_en AS approved_by_first_name_en,
		approved_by.first_name_ar AS approved_by_first_name_ar,
		approved_by.last_name_en AS approved_by_last_name_en,
		approved_by.last_name_ar AS approved_by_last_name_ar,
		approved_by.initial AS approved_by_initial,
		approved_by.color AS approved_by_color,
		approved_by.email AS approved_by_email,
		approved_by.profile_image_path AS approved_by_profile_image_path,
		CONVERT(VARCHAR, bia_processes.approved_at, 23) AS approved_at
	`).
		Joins("LEFT JOIN process_groups ON process_groups.id = processes.process_group_id").
		Joins("LEFT JOIN bia_processes ON bia_processes.process_id = processes.id").
		Joins("LEFT JOIN workflow_statuses ON workflow_statuses.id = bia_processes.workflow_status_id").
		Joins("LEFT JOIN users AS approved_by ON approved_by.id = bia_processes.approved_by_id").
		Joins("LEFT JOIN organization_units ON organization_units.id = processes.responsible_unit_id").
		Joins("LEFT JOIN organization_units AS divisions ON divisions.id = processes.division_id").
		Joins("LEFT JOIN users AS owner ON owner.id = processes.responsible_user_id").
		Joins("LEFT JOIN lifecycle_statuses ON lifecycle_statuses.id = processes.lifecycle_status_id").
		Joins("LEFT JOIN users AS created_user ON created_user.id = processes.created_by_id").
		Joins("LEFT JOIN users AS updated_user ON updated_user.id = processes.updated_by_id").
		Joins("LEFT JOIN users AS last_assessed_user ON last_assessed_user.id = processes.last_assessed_by_id").
		Joins("LEFT JOIN business_continuity_duration_units AS rto_bc_unit ON rto_bc_unit.id = processes.rto_unit_id").
		Joins("LEFT JOIN business_continuity_duration_units AS mtpd_bc_unit ON mtpd_bc_unit.id = processes.mtpd_unit_id").
		Joins("LEFT JOIN business_continuity_duration_units AS rpo_bc_unit ON rpo_bc_unit.id = processes.rpo_unit_id").
		Joins("LEFT JOIN business_continuity_duration_units AS rsp_bc_unit ON rsp_bc_unit.id = processes.rsp_unit_id").
		Joins("LEFT JOIN bcm_impact_analysis_levels ON bcm_impact_analysis_levels.id = processes.impact_analysis_level_id").
		Joins("LEFT JOIN risk_matrix_risk_ratings ON risk_matrix_risk_ratings.id = processes.risk_matrix_risk_rating_id").
		Joins("LEFT JOIN risk_matrix_risk_rating_values ON risk_matrix_risk_rating_values.id = processes.risk_matrix_risk_rating_value_id").
		Joins("LEFT JOIN risk_levels ON risk_levels.id = processes.risk_level_id").
		Joins("LEFT JOIN process_criticality_levels ON process_criticality_levels.id = processes.process_criticality_level_id").
		Joins("LEFT JOIN bcm_process_types ON bcm_process_types.id = processes.process_type_id").
		Joins("LEFT JOIN locations ON locations.id = processes.location_id").
		Joins("LEFT JOIN review_frequencies ON review_frequencies.id = processes.review_frequency_id")

	// Apply dynamic filters to the query
	filterColumnMap := map[string]string{
		"process_group_id":    "processes.process_group_id",
		"responsible_unit_id": "processes.responsible_unit_id",
		"responsible_user_id": "processes.responsible_user_id",
		"lifecycle_status_id": "processes.lifecycle_status_id",
		"division_id":         "processes.division_id",
		"process_type_id":     "processes.process_type_id",
		"location_id":         "processes.location_id",
		"review_frequency_id": "processes.review_frequency_id",
		"workflow_status_id":  "bia_processes.workflow_status_id",
	}
	for paramName, values := range params.Filters {
		if len(values) == 0 {
			continue
		}
		if col, ok := filterColumnMap[paramName]; ok {
			query = query.Where(col+" IN (?)", values)
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
		case "is_bia_performed":
			query = query.Where("processes.is_bia_performed = ?", true)
		case "bia_not_performed":
			query = query.Where("processes.is_bia_performed = ?", false)
		case "is_automatic":
			query = query.Where("processes.is_automatic = ?", true)
		case "is_manual":
			query = query.Where("processes.is_manual = ?", true)
		case "is_can_be_done_manually":
			query = query.Where("processes.is_can_be_done_manually = ?", true)
		case "is_critical":
			query = query.Where("bcm_impact_analysis_levels.type != ? AND bcm_impact_analysis_levels.type IS NOT NULL AND bia_processes.approved_by_id IS NOT NULL", "non_applicable")
		case "high":
			query = query.Where("process_criticality_levels.type = ?", "high")
		case "medium":
			query = query.Where("process_criticality_levels.type = ?", "medium")
		case "low":
			query = query.Where("process_criticality_levels.type = ?", "low")
		}
	}

	// Apply search keyword if exists
	if params.Keyword != "" {
		kw := "%" + params.Keyword + "%"
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
		 lifecycle_statuses.title_ar LIKE ? OR
		 divisions.title_en LIKE ? OR
		 divisions.title_ar LIKE ? OR
		 bcm_process_types.title_en LIKE ? OR
		 bcm_process_types.title_ar LIKE ? OR
		 locations.title_en LIKE ? OR
		 locations.title_ar LIKE ? OR
		 review_frequencies.title_en LIKE ? OR
		 review_frequencies.title_ar LIKE ? OR
		 workflow_statuses.title_en LIKE ? OR
		 workflow_statuses.title_ar LIKE ? OR
		 approved_by.first_name_en LIKE ? OR
		 approved_by.first_name_ar LIKE ? OR
		 approved_by.last_name_en LIKE ? OR
		 approved_by.last_name_ar LIKE ?)`,
			kw, kw, kw, kw, kw, kw, kw, kw, kw, kw, kw, kw,
			kw, kw, kw, kw, kw, kw, kw, kw, kw, kw, kw, kw, kw, kw)
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

	if err := repo.populateProcessListAssociations(processes); err != nil {
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

func (repo *ProcessRepository) populateProcessListAssociations(processes []*organizationDTOs.ProcessListDTO) error {
	if len(processes) == 0 {
		return nil
	}

	processIDs := make([]uint, len(processes))
	for i, p := range processes {
		processIDs[i] = p.ID
	}

	var executionRecords []*models.ProcessExecutionRecord
	if err := repo.db.Where("process_id IN ?", processIDs).Preload("BCMExecutionRecord").Find(&executionRecords).Error; err != nil {
		return err
	}
	execMap := make(map[uint][]*organizationDTOs.ProcessExecutionRecordDTO)
	for _, er := range executionRecords {
		if er == nil || er.BCMExecutionRecord == nil {
			continue
		}
		execMap[er.ProcessID] = append(execMap[er.ProcessID], &organizationDTOs.ProcessExecutionRecordDTO{
			ID:                   er.ID,
			ProcessID:            er.ProcessID,
			BCMExecutionRecordID: er.BCMExecutionRecordID,
			TitleEn:              er.BCMExecutionRecord.TitleEn,
			TitleAr:              er.BCMExecutionRecord.TitleAr,
		})
	}

	var manualForms []*models.ProcessManualProcedureForm
	if err := repo.db.Where("process_id IN ?", processIDs).Preload("BCMManualProcedureForm").Find(&manualForms).Error; err != nil {
		return err
	}
	formMap := make(map[uint][]*organizationDTOs.ProcessManualProcedureFormDTO)
	for _, f := range manualForms {
		if f == nil || f.BCMManualProcedureForm == nil {
			continue
		}
		formMap[f.ProcessID] = append(formMap[f.ProcessID], &organizationDTOs.ProcessManualProcedureFormDTO{
			ID:                       f.ID,
			ProcessID:                f.ProcessID,
			BCMManualProcedureFormID: f.BCMManualProcedureFormID,
			TitleEn:                  f.BCMManualProcedureForm.TitleEn,
			TitleAr:                  f.BCMManualProcedureForm.TitleAr,
			Code:                     f.BCMManualProcedureForm.Code,
		})
	}

	var systemRegisters []*models.ProcessSystemRegister
	if err := repo.db.Where("process_id IN ?", processIDs).Preload("SystemRegister").Find(&systemRegisters).Error; err != nil {
		return err
	}
	srMap := make(map[uint]*organizationDTOs.ProcessSystemRegisterDTO)
	for _, sr := range systemRegisters {
		if sr == nil || sr.SystemRegister == nil {
			continue
		}
		if _, exists := srMap[sr.ProcessID]; !exists {
			srMap[sr.ProcessID] = &organizationDTOs.ProcessSystemRegisterDTO{
				ID:               sr.ID,
				ProcessID:        sr.ProcessID,
				SystemRegisterID: sr.SystemRegisterID,
				TitleEn:          sr.SystemRegister.TitleEn,
				TitleAr:          sr.SystemRegister.TitleAr,
			}
		}
	}

	for _, p := range processes {
		if p == nil {
			continue
		}
		p.ExecutionRecords = execMap[p.ID]
		p.ManualProcedureForms = formMap[p.ID]
		p.SystemRegister = srMap[p.ID]
	}

	return nil
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
		Preload("ResponsibleUnit").
		Preload("ResponsibleUser").
		Preload("ResponsibleUser.OrganizationUnit1").
		Preload("LifecycleStatus").
		Preload("Division").
		Preload("ProcessType").
		Preload("Location").
		Preload("ReviewFrequency").
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
		Preload("RspUnit").
		Preload("ImpactAnalysisLevel").
		Preload("RiskLevel").
		Preload("ProcessCriticalityLevel").
		Preload("LastAssessedByUser").
		Preload("AuditableItems", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at DESC").Limit(5) // Fetch the latest 5 auditable items
		}).
		Preload("AuditableItems.ResponsibleUser").
		Preload("AuditableItems.ResponsibleUnit").
		Preload("AuditableItems.AuditorRiskRating").
		Preload("AuditableItems.CreatedByUser").
		Preload("ProcessExecutionRecords").
		Preload("ProcessExecutionRecords.BCMExecutionRecord").
		Preload("ProcessExecutionRecords.BCMExecutionRecord.OwnerEntity").
		Preload("ProcessExecutionRecords.BCMExecutionRecord.Location").
		Preload("ProcessExecutionRecords.BCMExecutionRecord.RecordType").
		Preload("ProcessManualProcedureForms").
		Preload("ProcessManualProcedureForms.BCMManualProcedureForm").
		Preload("ProcessManualProcedureForms.BCMManualProcedureForm.OwnerEntity").
		Preload("ProcessManualProcedureForms.BCMManualProcedureForm.Location").
		Preload("ProcessSystemRegister").
		Preload("ProcessSystemRegister.SystemRegister").
		Preload("ProcessSystemRegister.SystemRegister.RTOUnit").
		Preload("ProcessSystemRegister.SystemRegister.MTPDUnit").
		Preload("ProcessSystemRegister.SystemRegister.RPOUnit").
		Preload("ProcessSystemRegister.SystemRegister.RecoveryManager").
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

	if err := db.Where("process_id = ?", id).Delete(&models.ProcessExecutionRecord{}).Error; err != nil {
		return err
	}

	if err := db.Where("process_id = ?", id).Delete(&models.ProcessManualProcedureForm{}).Error; err != nil {
		return err
	}

	if err := db.Where("process_id = ?", id).Delete(&models.ProcessSystemRegister{}).Error; err != nil {
		return err
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

// ProcessExecutionRecord association methods
func (repo *ProcessRepository) GetProcessExecutionRecordsByProcessID(processID uint, tx ...utility.DatabaseTransaction) ([]*models.ProcessExecutionRecord, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var records []*models.ProcessExecutionRecord
	if err := db.Where("process_id = ?", processID).Find(&records).Error; err != nil {
		return nil, err
	}
	return records, nil
}

func (repo *ProcessRepository) CreateProcessExecutionRecord(record *models.ProcessExecutionRecord, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Create(record).Error
}

func (repo *ProcessRepository) DeleteProcessExecutionRecordsByProcessID(processID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Where("process_id = ?", processID).Delete(&models.ProcessExecutionRecord{}).Error
}

// ProcessManualProcedureForm association methods
func (repo *ProcessRepository) GetProcessManualProcedureFormsByProcessID(processID uint, tx ...utility.DatabaseTransaction) ([]*models.ProcessManualProcedureForm, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var forms []*models.ProcessManualProcedureForm
	if err := db.Where("process_id = ?", processID).Find(&forms).Error; err != nil {
		return nil, err
	}
	return forms, nil
}

func (repo *ProcessRepository) CreateProcessManualProcedureForm(form *models.ProcessManualProcedureForm, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Create(form).Error
}

func (repo *ProcessRepository) DeleteProcessManualProcedureFormsByProcessID(processID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Where("process_id = ?", processID).Delete(&models.ProcessManualProcedureForm{}).Error
}

// ProcessSystemRegister association methods
func (repo *ProcessRepository) GetProcessSystemRegistersByProcessID(processID uint, tx ...utility.DatabaseTransaction) ([]*models.ProcessSystemRegister, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var records []*models.ProcessSystemRegister
	if err := db.Where("process_id = ?", processID).Find(&records).Error; err != nil {
		return nil, err
	}
	return records, nil
}

func (repo *ProcessRepository) CreateProcessSystemRegister(record *models.ProcessSystemRegister, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Create(record).Error
}

func (repo *ProcessRepository) DeleteProcessSystemRegistersByProcessID(processID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Where("process_id = ?", processID).Delete(&models.ProcessSystemRegister{}).Error
}

func (repo *ProcessRepository) GetBCMProcessTypeByID(id uint, tx ...utility.DatabaseTransaction) (*models.BCMProcessType, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var item models.BCMProcessType
	if err := db.First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (repo *ProcessRepository) GetBCMExecutionRecordByID(id uint, tx ...utility.DatabaseTransaction) (*models.BCMExecutionRecord, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var item models.BCMExecutionRecord
	if err := db.First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (repo *ProcessRepository) GetBCMManualProcedureFormByID(id uint, tx ...utility.DatabaseTransaction) (*models.BCMManualProcedureForm, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var item models.BCMManualProcedureForm
	if err := db.First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (repo *ProcessRepository) GetSystemRegisterByID(id uint, tx ...utility.DatabaseTransaction) (*models.SystemRegister, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var item models.SystemRegister
	if err := db.First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (repo *ProcessRepository) GetAllBCMProcessTypes() ([]*models.BCMProcessType, error) {
	var items []*models.BCMProcessType
	if err := repo.db.Where("status = ?", "active").Order("title_en").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (repo *ProcessRepository) GetAllBCMExecutionRecords() ([]*models.BCMExecutionRecord, error) {
	var items []*models.BCMExecutionRecord
	if err := repo.db.Order("title_en").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (repo *ProcessRepository) GetAllBCMManualProcedureForms() ([]*models.BCMManualProcedureForm, error) {
	var items []*models.BCMManualProcedureForm
	if err := repo.db.Order("title_en").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (repo *ProcessRepository) GetAllSystemRegisters() ([]*models.SystemRegister, error) {
	var items []*models.SystemRegister
	if err := repo.db.Order("title_en").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}
