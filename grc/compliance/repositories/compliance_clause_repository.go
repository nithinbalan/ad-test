package repositories

import (
	"errors"
	commonDTOs "jethur/grc/common/dtos"
	complianceDTOs "jethur/grc/compliance/dtos"

	"jethur/grc/models"
	"jethur/utility"
	"sort"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

// ComplianceClauseRepositoryInterface defines the methods for compliance clause repository
type ComplianceClauseRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(complianceID uint) ([]*models.ComplianceClause, error)
	GetByID(id uint, tx ...utility.DatabaseTransaction) (*models.ComplianceClause, error)
	GetFullDetailsByID(clauseID uint) (*models.ComplianceClause, error)
	Create(clause *models.ComplianceClause, tx ...utility.DatabaseTransaction) (uint, error)
	Update(clause *models.ComplianceClause, tx ...utility.DatabaseTransaction) error
	Delete(clauseID uint, tx ...utility.DatabaseTransaction) error

	GetNextClauseSequence(parentClauseID uint, complianceRegisterID uint, tx ...utility.DatabaseTransaction) (int, error)
	GetNextTopLevelClauseSequence(complianceRegisterID uint, tx ...utility.DatabaseTransaction) (int, error)
	GetByClauseNo(clauseNo string, complianceRegisterID uint, tx ...utility.DatabaseTransaction) (*models.ComplianceClause, error)
	GetClauseWithChildren(clauseID uint, tx ...utility.DatabaseTransaction) (*models.ComplianceClause, error)

	GetAssociatedComplianceRequirements(complianceID uint, clauseID uint) ([]complianceDTOs.ComplianceRequirementListDTO, error)
	GetAssociatedControls(complianceID uint, clauseID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedControlListResponseDTO, error)
	CreateComplianceControls(controls []*models.ComplianceControl, tx ...utility.DatabaseTransaction) ([]uint, error)
	DeleteComplianceControlByClauseAndControl(clauseID, controlID uint, tx ...utility.DatabaseTransaction) error

	IsControlMappedToClause(clauseID, controlID uint) (bool, error)
	GetClausesByIDs(ids []uint, tx ...utility.DatabaseTransaction) ([]*models.ComplianceClause, error)

	GetAssociatedRisks(complianceID uint, clauseID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedRiskListResponseDTO, error)
	CreateComplianceClauseRisk(controls []*models.RiskComplianceRequirement, tx ...utility.DatabaseTransaction) error
	DeleteComplianceClauseRisk(clauseID, controlID uint, tx ...utility.DatabaseTransaction) error

	GetAssociatedDocuments(complianceID uint, clauseID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedDocumentListResponseDTO, error)
	CreateComplianceClauseDocument(controls []*models.DocumentComplianceRequirement, tx ...utility.DatabaseTransaction) error
	DeleteComplianceClauseDocument(clauseID, controlID uint, tx ...utility.DatabaseTransaction) error

	GetAssociatedIncident(complianceID uint, clauseID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedIncidentListResponseDTO, error)
	CreateComplianceClauseIncident(incidents []*models.IncidentComplianceRequirement, tx ...utility.DatabaseTransaction) error
	DeleteComplianceClauseIncident(clauseID, incidentID uint, tx ...utility.DatabaseTransaction) error

	GetAssociatedStrategicObjective(complianceID uint, clauseID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedItemStrategicObjectiveListResponseDTO, error)
	CreateComplianceClauseStrategicObjective(strategicObjectives []*models.ComplianceRequirementStrategicObjective, tx ...utility.DatabaseTransaction) error
	DeleteComplianceClausStrategicObjective(clauseID, strategicObjectiveId uint, tx ...utility.DatabaseTransaction) error

	GetAssociatedAssets(complianceID uint, clauseID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedItemAssetListResponseDTO, error)
	CreateComplianceClauseAssets(assets []*models.AssetComplianceRequirement, tx ...utility.DatabaseTransaction) error
	DeleteComplianceAssets(clauseID, assetID uint, tx ...utility.DatabaseTransaction) error

	GetAssociatedProcesses(complianceID uint, clauseID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedItemProcessListResponseDTO, error)
	CreateComplianceClausProcesses(process []*models.ProcessComplianceRequirement, tx ...utility.DatabaseTransaction) error
	DeleteComplianceClausProcesses(clauseID, processID uint, tx ...utility.DatabaseTransaction) error

	GetAssociatedKPI(complianceID uint, clauseID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedKPIRegisterListResponseDTO, error)
	CreateComplianceClausKPI(kpis []*models.KPICompliance, tx ...utility.DatabaseTransaction) error
	DeleteComplianceClausKPI(clauseID, kpiID uint, tx ...utility.DatabaseTransaction) error
}

type ComplianceClauseRepository struct {
	db *gorm.DB
}

func NewComplianceClauseRepository(db *gorm.DB) *ComplianceClauseRepository {
	return &ComplianceClauseRepository{db: db}
}

func (repo *ComplianceClauseRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *ComplianceClauseRepository) GetAll(complianceID uint) ([]*models.ComplianceClause, error) {
	var clauses []*models.ComplianceClause

	// Fetch all clauses with preloaded child clauses (recursive loading)
	err := repo.db.
		Preload("ChildClauses", func(db *gorm.DB) *gorm.DB {
			return db.Preload("ChildClauses") // Recursive preload
		}).
		Where("compliance_register_id = ?", complianceID).
		Find(&clauses).Error

	if err != nil {
		return nil, err
	}

	// Sort clauses correctly using numeric sorting
	sort.SliceStable(clauses, func(i, j int) bool {
		return compareClauseNumbers(clauses[i].ClauseNo, clauses[j].ClauseNo)
	})

	// Recursively sort child clauses
	for _, clause := range clauses {
		sortChildClauses(clause)
	}

	return clauses, nil
}

func (repo *ComplianceClauseRepository) GetByID(id uint, tx ...utility.DatabaseTransaction) (*models.ComplianceClause, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var complianceClause models.ComplianceClause
	result := db.First(&complianceClause, id)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// Return an error instead of nil to indicate record not found
			return nil, errors.New("compliance clause not found")
		}
		return nil, result.Error
	}

	return &complianceClause, nil
}

func (repo *ComplianceClauseRepository) GetFullDetailsByID(clauseID uint) (*models.ComplianceClause, error) {
	var clause models.ComplianceClause

	err := repo.db.
		Preload("ComplianceControls").
		Preload("ComplianceControls.Control").
		Preload("ComplianceRequirements").
		Preload("ComplianceRequirements.RequiredDocumentTitles").
		Preload("ComplianceRequirements.ComplianceRequirementCategory").
		Preload("ComplianceRequirements.ResponsibleUnits").
		Preload("ComplianceRequirements.ResponsibleUnits.ResponsibleUnit").
		Preload("ComplianceRequirements.ResponsibleUnits.ResponsibleUser").
		Preload("ComplianceRequirements.ResponsibleUnits.ComplianceAssessmentFrameworkValue").
		Where("id = ?", clauseID).
		First(&clause).Error

	if err != nil {
		return nil, err
	}

	return &clause, nil
}

func (repo *ComplianceClauseRepository) Create(clause *models.ComplianceClause, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(clause).Error; err != nil {
		return 0, err
	}

	return clause.ID, nil
}

func (repo *ComplianceClauseRepository) Update(clause *models.ComplianceClause, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	// Perform update
	if err := db.Save(clause).Error; err != nil {
		return err
	}

	return nil
}

func (repo *ComplianceClauseRepository) Delete(clauseID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	// Use a CTE to find all descendant clauses, including the parent clause
	query := `
		WITH RecursiveClauses AS (
			SELECT id
			FROM compliance_clauses
			WHERE id = ?

			UNION ALL

			SELECT c.id
			FROM compliance_clauses c
			INNER JOIN RecursiveClauses rc ON c.compliance_clause_id = rc.id
		)
		DELETE FROM compliance_clauses
		WHERE id IN (SELECT id FROM RecursiveClauses)
	`

	// Execute the query
	err := db.Exec(query, clauseID).Error
	if err != nil {
		return err
	}

	return nil
}

func (repo *ComplianceClauseRepository) GetNextClauseSequence(parentClauseID uint, complianceRegisterID uint, tx ...utility.DatabaseTransaction) (int, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var maxSequence int

	// Use TRY_CAST to safely extract and convert the last part of the clause number
	query := `
		SELECT COALESCE(MAX(TRY_CAST(RIGHT(clause_no, CHARINDEX('.', REVERSE(clause_no)) - 1) AS INT)), 0) 
		FROM compliance_clauses 
		WHERE compliance_clause_id = ? AND compliance_register_id = ?`

	err := db.Raw(query, parentClauseID, complianceRegisterID).Row().Scan(&maxSequence)
	if err != nil {
		return 0, err
	}

	return maxSequence + 1, nil
}

func (repo *ComplianceClauseRepository) GetNextTopLevelClauseSequence(complianceRegisterID uint, tx ...utility.DatabaseTransaction) (int, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var maxSequence int

	query := `
		SELECT COALESCE(MAX(TRY_CAST(clause_no AS INT)), 0) 
		FROM compliance_clauses 
		WHERE compliance_clause_id IS NULL AND compliance_register_id = ?`

	err := db.Raw(query, complianceRegisterID).Row().Scan(&maxSequence)
	if err != nil {
		return 0, err
	}

	return maxSequence + 1, nil
}

func (repo *ComplianceClauseRepository) GetByClauseNo(clauseNo string, complianceRegisterID uint, tx ...utility.DatabaseTransaction) (*models.ComplianceClause, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var clause models.ComplianceClause
	err := db.Where("clause_no = ? AND compliance_register_id = ?", clauseNo, complianceRegisterID).First(&clause).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &clause, nil
}

func (repo *ComplianceClauseRepository) GetClauseWithChildren(clauseID uint, tx ...utility.DatabaseTransaction) (*models.ComplianceClause, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var clause models.ComplianceClause
	err := db.
		Preload("ChildClauses", func(db *gorm.DB) *gorm.DB {
			return db.Preload("ChildClauses") // Recursive preload
		}).
		First(&clause, clauseID).Error

	if err != nil {
		return nil, err
	}

	return &clause, nil
}

func (repo *ComplianceClauseRepository) GetAssociatedComplianceRequirements(complianceID uint, clauseID uint) ([]complianceDTOs.ComplianceRequirementListDTO, error) {
	var requirements []complianceDTOs.ComplianceRequirementListDTO

	query := `
	SELECT 
		cr.id, 
		cr.compliance_clause_id AS clause_id,
		cc.clause_no AS compliance_clause_no,
		cc.title_en AS compliance_clause_title_en,
		cc.title_ar AS compliance_clause_title_ar,
		cr.title_en,
		cr.title_ar,
		crc.title_en AS compliance_requirement_category_title_en,
		crc.title_ar AS compliance_requirement_category_title_ar,
		sl.title_en AS severity_level_title_en,
		sl.title_ar AS severity_level_title_ar,
		sl.color_code AS severity_level_color_code,
		FORMAT(cr.deadline, 'yyyy-MM-dd') AS deadline,
		cs.title_en AS compliance_status_title_en,
		cs.title_ar AS compliance_status_title_ar,
		cs.color_code AS compliance_status_color_code,
		FORMAT(cr.created_at, 'yyyy-MM-dd HH:mm:ss') AS created_at,
		FORMAT(cr.updated_at, 'yyyy-MM-dd HH:mm:ss') AS updated_at,
		uc.id AS created_user_id,
		uc.first_name_en AS created_user_first_name_en,
		uc.first_name_ar AS created_user_first_name_ar,
		uc.last_name_en AS created_user_last_name_en,
		uc.last_name_ar AS created_user_last_name_ar,
		uc.initial AS created_user_initial,
		uc.color AS created_user_color,
		uc.email AS created_user_email,
		uc.profile_image_path AS created_user_profile_image_path,
		uu.id AS updated_user_id,
		uu.first_name_en AS updated_user_first_name_en,
		uu.first_name_ar AS updated_user_first_name_ar,
		uu.last_name_en AS updated_user_last_name_en,
		uu.last_name_ar AS updated_user_last_name_ar,
		uu.initial AS updated_user_initial,
		uu.color AS updated_user_color,
		uu.email AS updated_user_email,
		uu.profile_image_path AS updated_user_profile_image_path
	FROM compliance_requirements cr
	JOIN compliance_clauses cc ON cr.compliance_clause_id = cc.id
	JOIN compliance_requirement_categories crc ON cr.compliance_requirement_category_id = crc.id
	JOIN severity_levels sl ON cr.severity_level_id = sl.id
	JOIN compliance_statuses cs ON cr.compliance_status_id = cs.id
	LEFT JOIN users uc ON cr.created_by_id = uc.id
	LEFT JOIN users uu ON cr.updated_by_id = uu.id
	WHERE cr.compliance_clause_id = ? AND cc.compliance_register_id = ?
	ORDER BY cr.created_at DESC
`

	err := repo.db.Raw(query, clauseID, complianceID).Scan(&requirements).Error
	return requirements, err
}

func (r *ComplianceClauseRepository) GetAssociatedControls(complianceID uint, clauseID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedControlListResponseDTO, error) {
	var controlList []*commonDTOs.AssociatedControlListDTO
	var totalCount int64
	query := r.db.Table("compliance_controls cc").
		Select(`
			controls.*,
			cc.control_id,
			control_types.title_en AS control_type_title_en,
			control_types.title_ar AS control_type_title_ar,
			control_categories.title_en AS control_category_title_en,
			control_categories.title_ar AS control_category_title_ar,
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
			review_frequencies.title_en AS review_frequency_title_en,
			review_frequencies.title_ar AS review_frequency_title_ar,
			test_responsible_user.id AS test_responsible_user_id,
			test_responsible_user.first_name_en AS test_responsible_user_first_name_en,
			test_responsible_user.first_name_ar AS test_responsible_user_first_name_ar,
			test_responsible_user.last_name_en AS test_responsible_user_last_name_en,
			test_responsible_user.last_name_ar AS test_responsible_user_last_name_ar,
			test_responsible_user.initial AS test_responsible_user_initial,
			test_responsible_user.color AS test_responsible_user_color,
			test_responsible_user.email AS test_responsible_user_email,
			test_responsible_user.profile_image_path AS test_responsible_user_profile_image_path,
			test_frequencies.title_en AS test_frequency_title_en,
			test_frequencies.title_ar AS test_frequency_title_ar,
			lifecycle_statuses.title_en AS lifecycle_status_title_en,
			lifecycle_statuses.title_ar AS lifecycle_status_title_ar,
			lifecycle_statuses.color_code AS lifecycle_status_color_code,
			control_implementation_statuses.title_en AS control_implementation_status_title_en,
			control_implementation_statuses.title_ar AS control_implementation_status_title_ar,
			control_implementation_statuses.color_code AS control_implementation_status_color_code,
			control_efficiency_statuses.title_en AS control_efficiency_status_title_en,
			control_efficiency_statuses.title_ar AS control_efficiency_status_title_ar,
			control_efficiency_statuses.color_code AS control_efficiency_status_color_code,
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
			risk_matrix_risk_ratings.title_en AS risk_rating_title_en,
			risk_matrix_risk_ratings.title_ar AS risk_rating_title_ar,
			risk_matrix_risk_ratings.color_code AS risk_rating_color_code,
			risk_matrix_risk_rating_values.value AS risk_rating_value
		`).
		Joins("JOIN controls ON controls.id = cc.control_id").
		// Joins("JOIN compliance_requirements cr ON cr.id = cc.compliance_requirement_id").
		// Joins("JOIN compliance_clauses cl ON cl.id = cr.compliance_clause_id").
		Joins("JOIN compliance_clauses cl ON cl.id = cc.compliance_clause_id").
		Joins("LEFT JOIN control_types ON control_types.id = controls.control_type_id").
		Joins("LEFT JOIN control_categories ON control_categories.id = controls.control_category_id").
		Joins("LEFT JOIN organization_units ON organization_units.id = controls.responsible_unit_id").
		Joins("LEFT JOIN users AS owner ON owner.id = controls.responsible_user_id").
		Joins("LEFT JOIN control_implementation_statuses ON control_implementation_statuses.id = controls.control_implementation_status_id").
		Joins("LEFT JOIN control_efficiency_statuses ON control_efficiency_statuses.id = controls.control_efficiency_status_id").
		Joins("LEFT JOIN review_frequencies ON review_frequencies.id = controls.review_frequency_id").
		Joins("LEFT JOIN users AS test_responsible_user ON test_responsible_user.id = controls.test_responsible_user_id").
		Joins("LEFT JOIN review_frequencies AS test_frequencies ON test_frequencies.id = controls.test_frequency_id").
		Joins("LEFT JOIN lifecycle_statuses ON lifecycle_statuses.id = controls.lifecycle_status_id").
		Joins("LEFT JOIN users AS created_user ON created_user.id = controls.created_by_id").
		Joins("LEFT JOIN users AS updated_user ON updated_user.id = controls.updated_by_id").
		Joins("LEFT JOIN risk_matrix_risk_ratings ON risk_matrix_risk_ratings.id = controls.risk_matrix_risk_rating_id").
		Joins("LEFT JOIN risk_matrix_risk_rating_values ON risk_matrix_risk_rating_values.id = controls.risk_matrix_risk_rating_value_id").
		Where("cc.compliance_register_id = ? AND cl.id = ?", complianceID, clauseID)
		// Apply search keyword if exists
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
			(controls.title_en LIKE ? OR
			controls.title_ar LIKE ? OR
			controls.reference_number LIKE ? OR
			controls.code LIKE ? OR
			control_categories.title_en LIKE ? OR
			control_categories.title_ar LIKE ? OR
			control_types.title_en LIKE ? OR
			control_types.title_ar LIKE ? OR
			control_efficiency_statuses.title_en LIKE ? OR
			control_efficiency_statuses.title_ar LIKE ? OR
			control_implementation_statuses.title_en LIKE ? OR
			control_implementation_statuses.title_ar LIKE ? OR
			owner.first_name_en LIKE ? OR
			owner.first_name_ar LIKE ? OR
			owner.last_name_en LIKE ? OR
			owner.last_name_ar LIKE ? OR
			owner.email LIKE ?)`,
			keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword,
			keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword)
	}

	if params.SortBy == "" {
		params.SortBy = "controls.title_en"
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "asc"
	}
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}
	query = query.Order(params.SortBy + " " + params.Sort)

	offset := (params.Page - 1) * params.PageSize

	// Execute the query
	err := query.Offset(offset).Limit(params.PageSize).Find(&controlList).Error
	if err != nil {
		return nil, err
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	response := &commonDTOs.AssociatedControlListResponseDTO{
		Items: controlList,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}
	return response, nil
}

func compareClauseNumbers(a, b string) bool {
	aParts := strings.Split(a, ".")
	bParts := strings.Split(b, ".")

	for i := 0; i < len(aParts) && i < len(bParts); i++ {
		aNum, _ := strconv.Atoi(aParts[i])
		bNum, _ := strconv.Atoi(bParts[i])

		if aNum != bNum {
			return aNum < bNum
		}
	}

	return len(aParts) < len(bParts)
}

func sortChildClauses(clause *models.ComplianceClause) {
	if len(clause.ChildClauses) == 0 {
		return
	}

	sort.SliceStable(clause.ChildClauses, func(i, j int) bool {
		return compareClauseNumbers(clause.ChildClauses[i].ClauseNo, clause.ChildClauses[j].ClauseNo)
	})

	// Recursively sort nested child clauses
	for _, child := range clause.ChildClauses {
		sortChildClauses(child)
	}
}

func (repo *ComplianceClauseRepository) CreateComplianceControls(controls []*models.ComplianceControl, tx ...utility.DatabaseTransaction) ([]uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var createdIDs []uint
	for _, control := range controls {
		if err := db.Create(control).Error; err != nil {
			return nil, err
		}
		createdIDs = append(createdIDs, control.ID)
	}
	return createdIDs, nil
}

func (repo *ComplianceClauseRepository) DeleteComplianceControlByClauseAndControl(clauseID, controlID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Where("compliance_clause_id = ? AND control_id = ?", clauseID, controlID).
		Delete(&models.ComplianceControl{}).Error
}

func (repo *ComplianceClauseRepository) IsControlMappedToClause(clauseID, controlID uint) (bool, error) {
	var count int64
	err := repo.db.
		Model(&models.ComplianceControl{}).
		Where("compliance_clause_id = ? AND control_id = ?", clauseID, controlID).
		Count(&count).Error

	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (repo *ComplianceClauseRepository) GetClausesByIDs(ids []uint, tx ...utility.DatabaseTransaction) ([]*models.ComplianceClause, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var clauses []*models.ComplianceClause
	if err := db.
		Select("id, compliance_register_id").
		Where("id IN ?", ids).
		Find(&clauses).Error; err != nil {
		return nil, err
	}

	return clauses, nil
}

func (repo *ComplianceClauseRepository) GetAssociatedRisks(complianceID uint, clauseID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedRiskListResponseDTO, error) {
	var riskList []*commonDTOs.AssociatedRiskListDTO
	var totalCount int64

	query := repo.db.Table("risk_compliance_requirements").
		Select(`
			risk_compliance_requirements.id,
			risk_compliance_requirements.risk_register_id,
			risk_compliance_requirements.compliance_clause_id,
			 risk_register.reference_number,
            risk_register.title,
            risk_register.description,
			CONVERT(VARCHAR, risk_register.reported_at, 23) AS reported_at,
            risk_register.created_at,
            risk_register.updated_at,
			risk_register.current_risk_rating_value_id,
            risk_categories.title_en AS risk_category_title_en,
            risk_categories.title_ar AS risk_category_title_ar,
            risk_sources.title_en AS risk_source_title_en,
            risk_sources.title_ar AS risk_source_title_ar,
			organization_units.id AS responsible_unit_id,
            organization_units.title_en AS responsible_unit_title_en,
            organization_units.title_ar AS responsible_unit_title_ar,
			users.id AS responsible_user_id,
            users.first_name_en AS responsible_user_first_name_en,
            users.first_name_ar AS responsible_user_first_name_ar,
            users.last_name_en AS responsible_user_last_name_en,
            users.last_name_ar AS responsible_user_last_name_ar,
			users.initial AS responsible_user_initial,
        	users.color AS responsible_user_color,
        	users.email AS responsible_user_email,
        	users.profile_image_path AS responsible_user_profile_image_path,
            risk_statuses.title_en AS risk_status_title_en,
            risk_statuses.title_ar AS risk_status_title_ar,
			risk_response_strategies.title_en AS risk_response_strategy_title_en,
            risk_response_strategies.title_ar AS risk_response_strategy_title_ar,
            risk_statuses.color_code AS risk_status_color_code,
			risk_statuses.type AS risk_status_type,
            review_frequencies.title_en AS review_frequency_title_en,
            review_frequencies.title_ar AS review_frequency_title_ar,
			CONVERT(VARCHAR, risk_register.next_review_date, 23) AS next_review_date,
            risk_register.last_review_note,
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
            current_likelihood.title_en AS current_likelihood_title_en,
            current_likelihood.title_ar AS current_likelihood_title_ar,
            current_impact.title_en AS current_impact_title_en,
            current_impact.title_ar AS current_impact_title_ar,
            current_risk_rating.title_en AS current_risk_rating_title_en,
            current_risk_rating.title_ar AS current_risk_rating_title_ar,
            current_risk_rating_value.value AS current_risk_rating_value,
            current_risk_rating.color_code AS current_risk_rating_color_code,
            inherent_likelihood.title_en AS inherent_likelihood_title_en,
            inherent_likelihood.title_ar AS inherent_likelihood_title_ar,
            inherent_impact.title_en AS inherent_impact_title_en,
            inherent_impact.title_ar AS inherent_impact_title_ar,
            inherent_risk_rating.title_en AS inherent_risk_rating_title_en,
            inherent_risk_rating.title_ar AS inherent_risk_rating_title_ar,
            inherent_risk_rating_value.value AS inherent_risk_rating_value,
            inherent_risk_rating.color_code AS inherent_risk_rating_color_code,
            residual_likelihood.title_en AS residual_likelihood_title_en,
            residual_likelihood.title_ar AS residual_likelihood_title_ar,
            residual_impact.title_en AS residual_impact_title_en,
            residual_impact.title_ar AS residual_impact_title_ar,
            residual_risk_rating.title_en AS residual_risk_rating_title_en,
            residual_risk_rating.title_ar AS residual_risk_rating_title_ar,
            residual_risk_rating_value.value AS residual_risk_rating_value,
            residual_risk_rating.color_code AS residual_risk_rating_color_code
		`).
		Joins(`
			LEFT JOIN risk_register ON risk_register.id = risk_compliance_requirements.risk_register_id
			LEFT JOIN risk_categories ON risk_categories.id = risk_register.risk_category_id
            LEFT JOIN risk_sources ON risk_sources.id = risk_register.risk_source_id
            LEFT JOIN organization_units ON organization_units.id = risk_register.responsible_unit_id
            LEFT JOIN users ON users.id = risk_register.responsible_user_id
            LEFT JOIN risk_statuses ON risk_statuses.id = risk_register.risk_status_id
			LEFT JOIN risk_response_strategies ON risk_response_strategies.id = risk_register.risk_response_strategy_id
            LEFT JOIN users AS created_user ON created_user.id = risk_register.created_by_id
            LEFT JOIN users AS updated_user ON updated_user.id = risk_register.updated_by_id
            LEFT JOIN review_frequencies ON review_frequencies.id = risk_register.review_frequency_id
            LEFT JOIN risk_matrix_likelihoods AS current_likelihood ON current_likelihood.id = risk_register.current_likelihood_id
            LEFT JOIN risk_matrix_impacts AS current_impact ON current_impact.id = risk_register.current_impact_id
            LEFT JOIN risk_matrix_risk_ratings AS current_risk_rating ON current_risk_rating.id = risk_register.current_risk_rating_id
            LEFT JOIN risk_matrix_risk_rating_values AS current_risk_rating_value ON current_risk_rating_value.id = risk_register.current_risk_rating_value_id
            LEFT JOIN risk_matrix_likelihoods AS inherent_likelihood ON inherent_likelihood.id = risk_register.inherent_likelihood_id
            LEFT JOIN risk_matrix_impacts AS inherent_impact ON inherent_impact.id = risk_register.inherent_impact_id
            LEFT JOIN risk_matrix_risk_ratings AS inherent_risk_rating ON inherent_risk_rating.id = risk_register.inherent_risk_rating_id
            LEFT JOIN risk_matrix_risk_rating_values AS inherent_risk_rating_value ON inherent_risk_rating_value.id = risk_register.inherent_risk_rating_value_id
            LEFT JOIN risk_matrix_likelihoods AS residual_likelihood ON residual_likelihood.id = risk_register.residual_likelihood_id
            LEFT JOIN risk_matrix_impacts AS residual_impact ON residual_impact.id = risk_register.residual_impact_id
            LEFT JOIN risk_matrix_risk_ratings AS residual_risk_rating ON residual_risk_rating.id = risk_register.residual_risk_rating_id
            LEFT JOIN risk_matrix_risk_rating_values AS residual_risk_rating_value ON residual_risk_rating_value.id = risk_register.residual_risk_rating_value_id
		`).Where("risk_compliance_requirements.compliance_register_id = ? AND risk_compliance_requirements.compliance_clause_id = ?", complianceID, clauseID)

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
			(risk_register.reference_number LIKE ? OR
			risk_register.title LIKE ? OR
			risk_register.description LIKE ? OR
			risk_categories.title_en LIKE ? OR
			risk_categories.title_ar LIKE ? OR
			risk_sources.title_en LIKE ? OR
			risk_sources.title_ar LIKE ? OR
			organization_units.title_en LIKE ? OR
			organization_units.title_ar LIKE ? OR
			users.first_name_en LIKE ? OR
			users.first_name_ar LIKE ? OR
			users.last_name_en LIKE ? OR
			users.last_name_ar LIKE ? OR
			risk_statuses.title_en LIKE ? OR
			risk_statuses.title_ar LIKE ? OR
			risk_statuses.color_code LIKE ?)`,
			keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword)
	}

	if params.SortBy == "" {
		params.SortBy = "risk_register.title"
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "asc"
	}

	query = query.Order(params.SortBy + " " + params.Sort)
	offset := (params.Page - 1) * params.PageSize

	// Execute the query
	err := query.Offset(offset).Limit(params.PageSize).Find(&riskList).Error
	if err != nil {
		return nil, err
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	response := &commonDTOs.AssociatedRiskListResponseDTO{
		Items: riskList,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}
	return response, nil
}

func (repo *ComplianceClauseRepository) CreateComplianceClauseRisk(risks []*models.RiskComplianceRequirement, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(&risks).Error; err != nil {
		return err
	}
	return nil
}

func (repo *ComplianceClauseRepository) DeleteComplianceClauseRisk(clauseID, riskID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Where("compliance_clause_id = ? AND risk_register_id = ?", clauseID, riskID).
		Delete(&models.RiskComplianceRequirement{}).Error
}

func (repo *ComplianceClauseRepository) GetAssociatedDocuments(complianceID uint, clauseID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedDocumentListResponseDTO, error) {
	var documentList []*commonDTOs.AssociatedDocumentListDTO
	var totalCount int64

	query := repo.db.Table("document_compliance_requirements").
		Select(`
			document_compliance_requirements.id,
			document_compliance_requirements.document_id,
			documents.reference_number AS reference_number,
			documents.document_type_id AS document_type_id,
			documents.title_en AS title_en,
			documents.title_ar AS title_ar,
			documents.last_checkout_at AS last_checkout_at,
			documents.last_checkin_at AS last_checkin_at,
			documents.is_locked AS is_locked,
			documents.description_en AS description_en,
			documents.description_ar AS description_ar,
			documents.responsible_user_id,
			documents.created_at,
			documents.updated_at,
			documents.published_at,
			document_categories.title_en AS document_category_title_en,
			document_categories.title_ar AS document_category_title_ar,
			document_types.title_en AS document_type_title_en,
			document_types.title_ar AS document_type_title_ar,
			document_access_types.title_en AS document_access_type_title_en,
			document_access_types.title_ar AS document_type_access_title_ar,
			organization_units.title_en AS responsible_unit_title_en,
			organization_units.title_ar AS responsible_unit_title_ar,
			responsible_user.id AS responsible_user_id,
			responsible_user.first_name_en AS responsible_user_first_name_en,
			responsible_user.first_name_ar AS responsible_user_first_name_ar,
			responsible_user.last_name_en AS responsible_user_last_name_en,
			responsible_user.last_name_ar AS responsible_user_last_name_ar,
			responsible_user.initial AS responsible_user_initial,
			responsible_user.color AS responsible_user_color,
			responsible_user.email AS responsible_user_email,
			responsible_user.profile_image_path AS responsible_user_profile_image_path,
			document_versions.url AS document_url,
			document_versions.name AS document_name,
			document_versions.size AS document_size,
			document_versions.token AS document_token,
			document_versions.original_name AS document_original_name,
			document_versions.type AS document_type,
			document_versions.version AS document_version,
			document_statuses.title_en AS document_status_title_en,
			document_statuses.title_ar AS document_status_title_ar,
			document_statuses.color_code AS document_status_color_code,
			document_statuses.type AS document_status_type,
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
			published_user.first_name_en AS published_user_first_name_en,
			published_user.first_name_ar AS published_user_first_name_ar,
			published_user.last_name_en AS published_user_last_name_en,
			published_user.last_name_ar AS published_user_last_name_ar,
			published_user.initial AS published_user_initial,
			published_user.color AS published_user_color,
			published_user.email AS published_user_email,
			published_user.profile_image_path AS published_user_profile_image_path
		`).
		Joins(`
			LEFT JOIN documents ON documents.id = document_compliance_requirements.document_id
			LEFT JOIN document_categories ON document_categories.id = documents.document_category_id
			LEFT JOIN document_types ON document_types.id = documents.document_type_id
			LEFT JOIN document_access_types ON document_access_types.id = documents.document_access_type_id
			LEFT JOIN organization_units ON organization_units.id = documents.responsible_unit_id
			LEFT JOIN users AS responsible_user ON responsible_user.id = documents.responsible_user_id
			LEFT JOIN document_versions ON document_versions.document_id = documents.id AND document_versions.is_latest_version = 1
			LEFT JOIN document_statuses ON document_statuses.id = documents.document_status_id
			LEFT JOIN users AS created_user ON created_user.id = documents.created_by_id
			LEFT JOIN users AS updated_user ON updated_user.id = documents.updated_by_id
			LEFT JOIN users AS published_user ON published_user.id = documents.published_by_id
		`).Where("document_compliance_requirements.compliance_register_id = ? AND document_compliance_requirements.compliance_clause_id = ?", complianceID, clauseID)

	// Apply search keyword if exists
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
			(documents.title_en LIKE ? OR
			documents.title_ar LIKE ? OR
			documents.reference_number LIKE ? OR
			document_categories.title_en LIKE ? OR
			document_categories.title_ar LIKE ? OR
			responsible_user.first_name_en LIKE ? OR
			responsible_user.first_name_ar LIKE ? OR
			responsible_user.last_name_en LIKE ? OR
			responsible_user.last_name_ar LIKE ? OR
			organization_units.title_en LIKE ? OR
			organization_units.title_ar LIKE ? OR
			document_statuses.title_en LIKE ? OR
			document_statuses.title_ar LIKE ?)`,
			keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword)
	}

	if params.SortBy == "" {
		params.SortBy = "documents.title_en"
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "asc"
	}
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}
	query = query.Order(params.SortBy + " " + params.Sort)

	offset := (params.Page - 1) * params.PageSize

	// Execute the query
	err := query.Offset(offset).Limit(params.PageSize).Find(&documentList).Error
	if err != nil {
		return nil, err
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	response := &commonDTOs.AssociatedDocumentListResponseDTO{
		Items: documentList,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}
	return response, nil
}

func (repo *ComplianceClauseRepository) CreateComplianceClauseDocument(documents []*models.DocumentComplianceRequirement, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(&documents).Error; err != nil {
		return err
	}
	return nil
}

func (repo *ComplianceClauseRepository) DeleteComplianceClauseDocument(clauseID, documentID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Where("compliance_clause_id = ? AND document_id = ?", clauseID, documentID).
		Delete(&models.DocumentComplianceRequirement{}).Error
}

func (repo *ComplianceClauseRepository) GetAssociatedIncident(complianceID uint, clauseID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedIncidentListResponseDTO, error) {
	var totalCount int64
	var incidentList []*commonDTOs.AssociatedIncidentListDTO

	// Second Query: Fetch compliance requirement incidents using the compliance requirement IDs
	query := repo.db.Table("incident_compliance_requirements").
		Select(`
            incident_register.*,
			incident_compliance_requirements.id as id,
			incident_compliance_requirements.incident_register_id,
            incident_categories.title_en AS incident_category_title_en,
            incident_categories.title_ar AS incident_category_title_ar,
            incident_types.title_en AS incident_type_title_en,
            incident_types.title_ar AS incident_type_title_ar,
            incident_statuses.title_en AS incident_status_title_en,
            incident_statuses.title_ar AS incident_status_title_ar,
			incident_statuses.color_code AS incident_status_color_code,
			severity_levels.title_en AS severity_level_title_en,
			severity_levels.title_ar AS severity_level_title_ar,
			severity_levels.color_code AS severity_level_color_code,
			organization_units.title_en AS responsible_unit_title_en,
			organization_units.title_ar AS responsible_unit_title_ar,
			responsible_user.id AS responsible_user_id,
			responsible_user.first_name_en AS responsible_user_first_name_en,
            responsible_user.first_name_ar AS responsible_user_first_name_ar,
            responsible_user.last_name_en AS responsible_user_last_name_en,
            responsible_user.last_name_ar AS responsible_user_last_name_ar,
			responsible_user.initial AS responsible_user_initial,
            responsible_user.color AS responsible_user_color,
            responsible_user.email AS responsible_user_email,
            responsible_user.profile_image_path AS responsible_user_profile_image_path,
			reported_user.id AS reported_by_id,
            reported_user.first_name_en AS reported_user_first_name_en,
            reported_user.first_name_ar AS reported_user_first_name_ar,
            reported_user.last_name_en AS reported_user_last_name_en,
            reported_user.last_name_ar AS reported_user_last_name_ar,
			reported_user.initial AS reported_user_initial,
            reported_user.color AS reported_user_color,
            reported_user.email AS reported_user_email,
            reported_user.profile_image_path AS reported_user_profile_image_path,
			CONVERT(VARCHAR, incident_investigations.initiated_at, 23) AS investigation_initiated_at,
			incident_investigations.description AS investigation_description,
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
            LEFT JOIN incident_register ON incident_register.id = incident_compliance_requirements.incident_register_id
			 LEFT JOIN incident_categories ON incident_categories.id = incident_register.incident_category_id
            LEFT JOIN incident_sub_categories ON incident_sub_categories.id = incident_register.incident_sub_category_id
            LEFT JOIN incident_types ON incident_types.id = incident_register.incident_type_id
            LEFT JOIN incident_statuses ON incident_statuses.id = incident_register.incident_status_id
		    LEFT JOIN severity_levels ON severity_levels.id = incident_register.severity_level_id
            LEFT JOIN organization_units ON organization_units.id = incident_register.responsible_unit_id
		    LEFT JOIN users AS responsible_user ON responsible_user.id = incident_register.responsible_user_id
			LEFT JOIN users AS reported_user ON reported_user.id = incident_register.reported_by
			LEFT JOIN incident_investigations ON incident_investigations.incident_register_id = incident_register.id
			LEFT JOIN users AS created_user ON created_user.id = incident_register.created_by_id
			LEFT JOIN users AS updated_user ON updated_user.id = incident_register.updated_by_id
        `).Where("incident_compliance_requirements.compliance_register_id = ? AND incident_compliance_requirements.compliance_clause_id = ?", complianceID, clauseID)

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}
	// Apply search keyword if exists
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where("incident_register.title LIKE ?", keyword)
	}

	if params.SortBy == "" {
		params.SortBy = "incident_register.title"
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "asc"
	}

	query = query.Order(params.SortBy + " " + params.Sort)

	offset := (params.Page - 1) * params.PageSize

	err := query.Offset(offset).Limit(params.PageSize).Find(&incidentList).Error
	if err != nil {
		return nil, err
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	response := &commonDTOs.AssociatedIncidentListResponseDTO{
		Items: incidentList,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}
	return response, nil
}

func (repo *ComplianceClauseRepository) CreateComplianceClauseIncident(incidents []*models.IncidentComplianceRequirement, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(&incidents).Error; err != nil {
		return err
	}
	return nil
}

func (repo *ComplianceClauseRepository) DeleteComplianceClauseIncident(clauseID, incidentId uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Where("compliance_clause_id = ? AND incident_register_id = ?", clauseID, incidentId).
		Delete(&models.IncidentComplianceRequirement{}).Error
}

func (repo *ComplianceClauseRepository) GetAssociatedStrategicObjective(complianceID uint, clauseID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedItemStrategicObjectiveListResponseDTO, error) {
	var strategicObjectiveList []*commonDTOs.AssociatedStrategicObjectiveListDTO
	var totalCount int64

	query := repo.db.Table("compliance_requirement_strategic_objectives").
		Select(`
			strategic_objectives.*,
			compliance_requirement_strategic_objectives.strategic_objective_id,
			strategic_objective_categories.title_en AS strategic_objective_category_title_en,
			strategic_objective_categories.title_ar AS strategic_objective_category_title_ar,
			organization_units.title_en AS responsible_unit_title_en,
			organization_units.title_ar AS responsible_unit_title_ar,
			responsible_user.first_name_en AS responsible_user_first_name_en,
			responsible_user.last_name_en AS responsible_user_last_name_en,
			responsible_user.first_name_ar AS responsible_user_first_name_ar,
			responsible_user.last_name_ar AS responsible_user_last_name_ar,
			responsible_user.initial AS responsible_user_initial,
			responsible_user.email AS responsible_user_email,
			responsible_user.color AS responsible_user_color,
			responsible_user.profile_image_path AS responsible_user_profile_image_path,
			review_user.first_name_en AS review_user_first_name_en,
			review_user.first_name_ar AS review_user_first_name_ar,
			review_user.last_name_en AS review_user_last_name_en,
			review_user.last_name_ar AS review_user_last_name_ar,
			review_user.initial AS review_user_initial,
			review_user.color AS review_user_color,
			review_user.email AS review_user_email,
			review_user.profile_image_path AS review_user_profile_image_path,
			review_frequencies.title_en AS review_frequency_title_en,
			review_frequencies.title_ar AS review_frequency_title_ar,
			created_user.first_name_en AS created_user_first_name_en,
			created_user.first_name_ar AS created_user_first_name_ar,
			created_user.last_name_en AS created_user_last_name_en,
			created_user.last_name_ar AS created_user_last_name_ar,
			created_user.initial AS created_user_initial,
			created_user.color AS created_user_color,
			created_user.email AS created_user_email,
			created_user.profile_image_path AS created_user_profile_image_path,
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
			LEFT JOIN strategic_objectives ON strategic_objectives.id = compliance_requirement_strategic_objectives.strategic_objective_id
			LEFT JOIN strategic_objective_categories ON strategic_objective_categories.id = strategic_objectives.strategic_objective_category_id
			LEFT JOIN organization_units ON organization_units.id = strategic_objectives.responsible_unit_id
			LEFT JOIN users AS responsible_user ON responsible_user.id = strategic_objectives.responsible_user_id
			LEFT JOIN users AS review_user ON review_user.id = strategic_objectives.review_user_id
			LEFT JOIN review_frequencies ON review_frequencies.id = strategic_objectives.review_frequency_id
			LEFT JOIN users AS created_user ON created_user.id = strategic_objectives.created_by_id
			LEFT JOIN users AS updated_user ON updated_user.id = strategic_objectives.updated_by_id
		`).Where("compliance_requirement_strategic_objectives.compliance_register_id = ? AND compliance_requirement_strategic_objectives.compliance_clause_id = ?", complianceID, clauseID)

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	// Apply search keyword if exists
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
			(strategic_objectives.title_en LIKE ? OR
			strategic_objectives.title_ar LIKE ?)`,
			keyword, keyword)
	}

	if params.SortBy == "" {
		params.SortBy = "strategic_objectives.title_en"
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "asc"
	}

	query = query.Order(params.SortBy + " " + params.Sort)

	offset := (params.Page - 1) * params.PageSize

	// Execute the query
	err := query.Offset(offset).Limit(params.PageSize).Find(&strategicObjectiveList).Error
	if err != nil {
		return nil, err
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	response := &commonDTOs.AssociatedItemStrategicObjectiveListResponseDTO{
		Items: strategicObjectiveList,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}
	return response, nil
}

func (repo *ComplianceClauseRepository) CreateComplianceClauseStrategicObjective(strtegicObjectives []*models.ComplianceRequirementStrategicObjective, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(&strtegicObjectives).Error; err != nil {
		return err
	}
	return nil
}

func (repo *ComplianceClauseRepository) DeleteComplianceClausStrategicObjective(clauseID, strategicObjectiveId uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Where("compliance_clause_id = ? AND strategic_objective_id = ?", clauseID, strategicObjectiveId).
		Delete(&models.ComplianceRequirementStrategicObjective{}).Error
}

func (repo *ComplianceClauseRepository) GetAssociatedAssets(complianceID uint, clauseID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedItemAssetListResponseDTO, error) {
	var riskAssets []*commonDTOs.AssociatedAssetListDTO
	var totalCount int64

	query := repo.db.Table("asset_compliance_requirements").
		Select(`
		asset_compliance_requirements.id,
		asset_compliance_requirements.asset_register_id,
		asset_register.reference_number,
		asset_register.title_en,
		asset_register.title_ar,
		asset_register.location,
		asset_register.description,
		CONVERT(VARCHAR, asset_register.acquisition_date, 23) AS acquisition_date,
		asset_categories.title_en AS asset_category_title_en,
		asset_categories.title_ar AS asset_category_title_ar,
		asset_statuses.title_en AS asset_status_title_en,
		asset_statuses.title_ar AS asset_status_title_ar,
		asset_statuses.color_code AS asset_status_color_code,
		organization_units.title_en AS responsible_unit_title_en,
		organization_units.title_ar AS responsible_unit_title_ar,
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
		Joins("LEFT JOIN asset_register ON asset_register.id = asset_compliance_requirements.asset_register_id").
		Joins("LEFT JOIN asset_categories ON asset_categories.id = asset_register.asset_category_id").
		Joins("LEFT JOIN asset_statuses ON asset_statuses.id = asset_register.asset_status_id").
		Joins("LEFT JOIN organization_units ON organization_units.id = asset_register.responsible_unit_id").
		Joins("LEFT JOIN users AS responsible_user ON responsible_user.id = asset_register.responsible_user_id").
		Joins("LEFT JOIN users AS created_user ON created_user.id = asset_register.created_by_id").
		Joins("LEFT JOIN users AS updated_user ON updated_user.id = asset_register.updated_by_id").
		Where("asset_compliance_requirements.compliance_register_id = ? AND asset_compliance_requirements.compliance_clause_id = ?", complianceID, clauseID)

	// Calculate total count for pagination
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	// Apply search keyword if exists
	if params.Keyword != "" {
		query = query.Where(`
			(asset_register.title_en LIKE ? OR
			asset_register.title_ar LIKE ? OR
			asset_register.reference_number LIKE ? OR
			asset_register.description LIKE ? OR
			asset_categories.title_en LIKE ? OR
			asset_categories.title_ar LIKE ? OR
			asset_statuses.title_en LIKE ? OR
			asset_statuses.title_ar LIKE ?)`,

			"%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%",
			"%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%")
	}

	if params.SortBy == "" {
		params.SortBy = "asset_register.title_en"
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "asc"
	}

	query = query.Order(params.SortBy + " " + params.Sort)

	offset := (params.Page - 1) * params.PageSize

	// Execute the query
	err := query.Offset(offset).Limit(params.PageSize).Find(&riskAssets).Error
	if err != nil {
		return nil, err
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	response := &commonDTOs.AssociatedItemAssetListResponseDTO{
		Items: riskAssets,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}

	return response, nil
}

func (repo *ComplianceClauseRepository) CreateComplianceClauseAssets(assets []*models.AssetComplianceRequirement, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(&assets).Error; err != nil {
		return err
	}
	return nil
}

func (repo *ComplianceClauseRepository) DeleteComplianceAssets(clauseID, assetID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Where("compliance_clause_id = ? AND asset_register_id = ?", clauseID, assetID).
		Delete(&models.AssetComplianceRequirement{}).Error
}

func (repo *ComplianceClauseRepository) GetAssociatedProcesses(complianceID uint, clauseID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedItemProcessListResponseDTO, error) {
	var processList []*commonDTOs.AssociatedProcessListDTO
	var totalCount int64

	query := repo.db.Table("process_compliance_requirements").
		Select(`
			processes.id,
			process_compliance_requirements.process_id,
			processes.reference_number,
			processes.title_en,
			processes.title_ar,
			processes.description,
			process_groups.title_en AS process_group_title_en,
			process_groups.title_ar AS process_group_title_ar,
			organization_units.title_en AS owner_unit_title_en,
			organization_units.title_ar AS owner_unit_title_ar,
			owner.first_name_en AS owner_first_name_en,
			owner.first_name_ar AS owner_first_name_ar,
			owner.last_name_en AS owner_last_name_en,
			owner.last_name_ar AS owner_last_name_ar,
			owner.initial AS control_owner_initial,
            owner.color AS control_owner_color,
            owner.email AS control_owner_email,
            owner.profile_image_path AS control_owner_profile_image_path,
			lifecycle_statuses.title_en AS lifecycle_status_title_en,
			lifecycle_statuses.title_ar AS lifecycle_status_title_ar,
			lifecycle_statuses.color_code AS lifecycle_status_color_code
		`).
		Joins(`
			LEFT JOIN processes ON processes.id = process_compliance_requirements.process_id
			LEFT JOIN process_groups ON process_groups.id = processes.process_group_id
			LEFT JOIN organization_units ON organization_units.id = processes.responsible_unit_id
			LEFT JOIN users AS owner ON owner.id = processes.responsible_user_id
			LEFT JOIN lifecycle_statuses ON lifecycle_statuses.id = processes.lifecycle_status_id
		`).Where("process_compliance_requirements.compliance_register_id = ? AND process_compliance_requirements.compliance_clause_id = ?", complianceID, clauseID)

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	// Apply search keyword if exists
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
			(processes.title_en LIKE ? OR
			processes.title_ar LIKE ? OR
			processes.reference_number LIKE ? OR
			process_groups.title_en LIKE ? OR
			process_groups.title_ar LIKE ? OR
			organization_units.title_en LIKE ? OR
			organization_units.title_ar LIKE ? OR
			owner.first_name_en LIKE ? OR
			owner.first_name_ar LIKE ? OR
			owner.last_name_en LIKE ? OR
			owner.last_name_ar LIKE ? OR
			lifecycle_statuses.title_en LIKE ? OR
			lifecycle_statuses.title_ar LIKE ? OR
			lifecycle_statuses.color_code LIKE ?)`,
			keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword,
			keyword, keyword, keyword, keyword, keyword)
	}

	if params.SortBy == "" {
		params.SortBy = "processes.title_en"
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "asc"
	}

	query = query.Order(params.SortBy + " " + params.Sort)

	offset := (params.Page - 1) * params.PageSize

	// Execute the query
	err := query.Offset(offset).Limit(params.PageSize).Find(&processList).Error
	if err != nil {
		return nil, err
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	response := &commonDTOs.AssociatedItemProcessListResponseDTO{
		Items: processList,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}
	return response, nil
}

func (repo *ComplianceClauseRepository) CreateComplianceClausProcesses(processes []*models.ProcessComplianceRequirement, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(&processes).Error; err != nil {
		return err
	}
	return nil
}

func (repo *ComplianceClauseRepository) DeleteComplianceClausProcesses(clauseID, processID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Where("compliance_clause_id = ? AND process_id = ?", clauseID, processID).
		Delete(&models.ProcessComplianceRequirement{}).Error
}

func (repo *ComplianceClauseRepository) GetAssociatedKPI(complianceID uint, clauseID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedKPIRegisterListResponseDTO, error) {
	var kpiRegister []*commonDTOs.AssociatedKPIRegisterListDTO
	var totalCount int64

	query := repo.db.Table("kpi_compliances").
		Select(`
            kpi_register.*,
			kpi_compliances.kpi_register_id,
			kpi_categories.title_en AS kpi_category_title_en,
			kpi_categories.title_ar AS kpi_category_title_ar,
			CONVERT(VARCHAR, kpi_register.start_date, 23) AS start_date,
			CONVERT(VARCHAR, kpi_register.end_date, 23) AS target_date,
			responsible_user.id AS responsible_user_id,
			responsible_user.first_name_en AS responsible_user_first_name_en,
			responsible_user.first_name_ar AS responsible_user_first_name_ar,
			responsible_user.last_name_en AS responsible_user_last_name_en,
			responsible_user.last_name_ar AS responsible_user_last_name_ar,
			responsible_user.initial AS responsible_user_initial,
			responsible_user.color AS responsible_user_color,
			responsible_user.email AS responsible_user_email,
			responsible_user.profile_image_path AS responsible_user_profile_image_path,
			kpi_statuses.title_en AS kpi_status_title_en,
			kpi_statuses.title_ar AS kpi_status_title_ar,
			kpi_statuses.type AS kpi_status_type,
			kpi_statuses.color_code AS kpi_status_color_code,
			review_frequencies.title_en AS review_frequency_title_en,
			review_frequencies.title_ar AS review_frequency_title_ar,
			review_frequencies.type AS review_frequency_type,
			calculation_types.title_en AS calculation_type_title_en,
			calculation_types.title_ar AS calculation_type_title_ar,
			calculation_types.type AS calculation_type_type,
			units.title_en AS unit_title_en,
			units.title_ar AS unit_title_ar,
			units.code_en AS unit_code_en,
			units.code_ar AS unit_code_ar,
			units.color_code AS unit_color_code,
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
		    LEFT JOIN kpi_register ON kpi_register.id = kpi_compliances.kpi_register_id
			LEFT JOIN kpi_categories ON kpi_categories.id = kpi_register.kpi_category_id
			LEFT JOIN kpi_statuses ON kpi_statuses.id = kpi_register.kpi_status_id
			LEFT JOIN calculation_types ON calculation_types.id = kpi_register.calculation_type_id
			LEFT JOIN review_frequencies ON review_frequencies.id = kpi_register.review_frequency_id
			LEFT JOIN units ON units.id = kpi_register.unit_id
			LEFT JOIN users AS responsible_user ON responsible_user.id = kpi_register.responsible_user_id
			LEFT JOIN users AS created_user ON created_user.id = kpi_register.created_by_id
			LEFT JOIN users AS updated_user ON updated_user.id = kpi_register.updated_by_id
        `).Where("kpi_compliances.compliance_register_id = ? AND kpi_compliances.compliance_clause_id = ?", complianceID, clauseID)

	// Apply search keyword if exists
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
			(kpi_register.title LIKE ? OR
			created_user.first_name_en LIKE ? OR
			created_user.first_name_ar LIKE ? OR
			created_user.last_name_en LIKE ? OR
			created_user.last_name_ar LIKE ? OR
    		updated_user.first_name_en LIKE ? OR
			updated_user.first_name_ar LIKE ? OR
			updated_user.last_name_en LIKE ? OR
			updated_user.last_name_ar LIKE ?)`,
			keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword, keyword)
	}

	if params.SortBy == "" {
		params.SortBy = "kpi_register.title"
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "asc"
	}
	// Calculate total count for pagination
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}
	query = query.Order(params.SortBy + " " + params.Sort)

	// Apply pagination
	offset := (params.Page - 1) * params.PageSize
	query = query.Offset(offset).Limit(params.PageSize)

	// Fetch the data
	if err := query.Find(&kpiRegister).Error; err != nil {
		return nil, err
	}

	// Calculate total pages
	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	response := &commonDTOs.AssociatedKPIRegisterListResponseDTO{
		Items: kpiRegister,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}
	return response, nil
}

func (repo *ComplianceClauseRepository) CreateComplianceClausKPI(kpis []*models.KPICompliance, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(&kpis).Error; err != nil {
		return err
	}
	return nil
}

func (repo *ComplianceClauseRepository) DeleteComplianceClausKPI(clauseID, kpiID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Where("compliance_clause_id = ? AND kpi_register_id = ?", clauseID, kpiID).
		Delete(&models.KPICompliance{}).Error
}
