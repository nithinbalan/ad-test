package repositories

import (
	"errors"
	complianceDTOs "jethur/grc/compliance/dtos"
	"jethur/grc/models"
	"jethur/utility"

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

	// Fetch all clauses with preloaded child clauses
	err := repo.db.
		Preload("ChildClauses", func(db *gorm.DB) *gorm.DB {
			return db.Preload("ChildClauses") // Recursive preload
		}).
		Where("compliance_register_id = ?", complianceID).
		Order("clause_no").
		Find(&clauses).Error

	if err != nil {
		return nil, err
	}

	return clauses, nil
}

func (repo *ComplianceClauseRepository) GetByID(id uint, tx ...utility.DatabaseTransaction) (*models.ComplianceClause, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var complianceClause *models.ComplianceClause
	result := db.First(&complianceClause, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return complianceClause, nil
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
			cc.compliance_clause_no,
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
			u.initials AS responsible_user_initial,
			u.color AS responsible_user_color,
			u.email AS responsible_user_email,
			u.profile_image_path AS responsible_user_profile_image_path,
			DATE_FORMAT(cr.deadline, '%Y-%m-%d') AS deadline,
			cs.title_en AS compliance_status_title_en,
			cs.title_ar AS compliance_status_title_ar,
			cs.color_code AS compliance_status_color_code,
			DATE_FORMAT(cr.created_at, '%Y-%m-%d %H:%i:%s') AS created_at,
			DATE_FORMAT(cr.updated_at, '%Y-%m-%d %H:%i:%s') AS updated_at,
			uc.id AS created_user_id,
			uc.first_name_en AS created_user_first_name_en,
			uc.first_name_ar AS created_user_first_name_ar,
			uc.last_name_en AS created_user_last_name_en,
			uc.last_name_ar AS created_user_last_name_ar,
			uc.initials AS created_user_initial,
			uc.color AS created_user_color,
			uc.email AS created_user_email,
			uc.profile_image_path AS created_user_profile_image_path,
			uu.id AS updated_user_id,
			uu.first_name_en AS updated_user_first_name_en,
			uu.first_name_ar AS updated_user_first_name_ar,
			uu.last_name_en AS updated_user_last_name_en,
			uu.last_name_ar AS updated_user_last_name_ar,
			uu.initials AS updated_user_initial,
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
		WHERE cr.compliance_clause_id = ? AND cc.compliance_id = ?
		ORDER BY cr.created_at DESC
	`

	err := repo.db.Raw(query, clauseID, complianceID).Scan(&requirements).Error
	return requirements, err
}
