package repositories

import (
	"errors"
	complianceDTOs "jethur/grc/compliance/dtos"
	"jethur/grc/models"
	organizationDTOs "jethur/grc/organization/dtos"
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
	Create(clause *models.ComplianceClause, tx ...utility.DatabaseTransaction) (uint, error)
	Update(clause *models.ComplianceClause, tx ...utility.DatabaseTransaction) error
	Delete(clauseID uint, tx ...utility.DatabaseTransaction) error

	GetNextClauseSequence(parentClauseID uint, complianceRegisterID uint, tx ...utility.DatabaseTransaction) (int, error)
	GetNextTopLevelClauseSequence(complianceRegisterID uint, tx ...utility.DatabaseTransaction) (int, error)
	GetByClauseNo(clauseNo string, complianceRegisterID uint, tx ...utility.DatabaseTransaction) (*models.ComplianceClause, error)
	GetClauseWithChildren(clauseID uint, tx ...utility.DatabaseTransaction) (*models.ComplianceClause, error)

	GetAssociatedComplianceRequirements(complianceID uint, clauseID uint) ([]complianceDTOs.ComplianceRequirementListDTO, error)
	GetAssociatedControls(complianceID uint, clauseID uint) ([]organizationDTOs.ControlListDTO, error)
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
		u.id AS responsible_user_id,
		u.first_name_en AS responsible_user_first_name_en,
		u.first_name_ar AS responsible_user_first_name_ar,
		u.last_name_en AS responsible_user_last_name_en,
		u.last_name_ar AS responsible_user_last_name_ar,
		u.initial AS responsible_user_initial,
		u.color AS responsible_user_color,
		u.email AS responsible_user_email,
		u.profile_image_path AS responsible_user_profile_image_path,
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
	LEFT JOIN users u ON cr.responsible_user_id = u.id
	LEFT JOIN users uc ON cr.created_by_id = uc.id
	LEFT JOIN users uu ON cr.updated_by_id = uu.id
	WHERE cr.compliance_clause_id = ? AND cc.compliance_register_id = ?
	ORDER BY cr.created_at DESC
`

	err := repo.db.Raw(query, clauseID, complianceID).Scan(&requirements).Error
	return requirements, err
}

func (r *ComplianceClauseRepository) GetAssociatedControls(complianceID uint, clauseID uint) ([]organizationDTOs.ControlListDTO, error) {
	controls := make([]organizationDTOs.ControlListDTO, 0)

	err := r.db.Table("compliance_controls cc").
		Select(`
			c.id,
			ct.title_en AS control_type_title_en,
			ct.title_ar AS control_type_title_ar,
			cc2.title_en AS control_category_title_en,
			cc2.title_ar AS control_category_title_ar,
			ou.title_en AS owner_unit_title_en,
			ou.title_ar AS owner_unit_title_ar,
			owner.id AS owner_id,
			owner.first_name_en AS owner_first_name_en,
			owner.first_name_ar AS owner_first_name_ar,
			owner.last_name_en AS owner_last_name_en,
			owner.last_name_ar AS owner_last_name_ar,
			owner.initial AS owner_initial,
			owner.color AS owner_color,
			owner.email AS owner_email,
			owner.profile_image_path AS owner_profile_image_path,
			cis.title_en AS control_implementation_status_title_en,
			cis.title_ar AS control_implementation_status_title_ar,
			cis.color_code AS control_implementation_status_color_code,
			ces.title_en AS control_efficiency_status_title_en,
			ces.title_ar AS control_efficiency_status_title_ar,
			ces.color_code AS control_efficiency_status_color_code,
			c.reference_number,
			c.code,
			c.title_en,
			c.title_ar,
			c.description,
			c.objective,
			rf.title_en AS review_frequency_title_en,
			rf.title_ar AS review_frequency_title_ar,
			c.last_review_date,
			c.next_review_date,
			c.last_review_note,
			test_resp.id AS test_responsible_user_id,
			test_resp.first_name_en AS test_responsible_user_first_name_en,
			test_resp.first_name_ar AS test_responsible_user_first_name_ar,
			test_resp.last_name_en AS test_responsible_user_last_name_en,
			test_resp.last_name_ar AS test_responsible_user_last_name_ar,
			test_resp.initial AS test_responsible_user_initial,
			test_resp.color AS test_responsible_user_color,
			test_resp.email AS test_responsible_user_email,
			test_resp.profile_image_path AS test_responsible_user_profile_image_path,
			tf.title_en AS test_frequency_title_en,
			tf.title_ar AS test_frequency_title_ar,
			c.last_test_date,
			c.next_test_date,
			c.implemented_at,
			c.is_applicable,
			c.justification,
			c.implementation_guideline,
			c.applicability_change_date,
			cls.title_en AS lifecycle_status_title_en,
			cls.title_ar AS lifecycle_status_title_ar,
			cls.color_code AS lifecycle_status_color_code,
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
			c.created_at,
			c.updated_at,
			rr.title_en AS risk_rating_title_en,
			rr.title_ar AS risk_rating_title_ar,
			rr.color_code AS risk_rating_color_code,
			rr.value AS risk_rating_value
		`).
		Joins("JOIN controls c ON c.id = cc.control_id").
		Joins("JOIN compliance_requirements cr ON cr.id = cc.compliance_requirement_id").
		Joins("JOIN compliance_clauses cl ON cl.id = cr.compliance_clause_id").
		Joins("LEFT JOIN control_types ct ON ct.id = c.control_type_id").
		Joins("LEFT JOIN control_categories cc2 ON cc2.id = c.control_category_id").
		Joins("LEFT JOIN organization_units ou ON ou.id = c.owner_unit_id").
		Joins("LEFT JOIN users owner ON owner.id = c.owner_id").
		Joins("LEFT JOIN control_implementation_statuses cis ON cis.id = c.control_implementation_status_id").
		Joins("LEFT JOIN control_efficiency_statuses ces ON ces.id = c.control_efficiency_status_id").
		Joins("LEFT JOIN review_frequencies rf ON rf.id = c.review_frequency_id").
		Joins("LEFT JOIN users test_resp ON test_resp.id = c.test_responsible_user_id").
		Joins("LEFT JOIN test_frequencies tf ON tf.id = c.test_frequency_id").
		Joins("LEFT JOIN control_lifecycle_statuses cls ON cls.id = c.lifecycle_status_id").
		Joins("LEFT JOIN users created_user ON created_user.id = c.created_by_id").
		Joins("LEFT JOIN users updated_user ON updated_user.id = c.updated_by_id").
		Joins("LEFT JOIN risk_ratings rr ON rr.id = c.risk_rating_id").
		Where("cc.compliance_register_id = ? AND cl.id = ?", complianceID, clauseID).
		Scan(&controls).Error

	if err != nil {
		return nil, err
	}

	return controls, nil
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
