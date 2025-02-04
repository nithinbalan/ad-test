package repositories

import (
	"errors"
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
