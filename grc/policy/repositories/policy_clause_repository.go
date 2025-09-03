package repositories

import (
	"jethur/grc/models"
	policyDTOs "jethur/grc/policy/dtos"
	"jethur/utility"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

// Interface
type PolicyClauseRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(params *policyDTOs.PolicyClauseListQueryParamsDTO, documentVersionId uint) ([]*models.DocumentVersionClause, error)
	GetByID(id uint, tx ...utility.DatabaseTransaction) (*models.DocumentVersionClause, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.DocumentVersionClause, error)
	GetByClauseNumber(clauseNumber string, tx ...utility.DatabaseTransaction) (*models.DocumentVersionClause, error)
	Create(documentVersionClause *models.DocumentVersionClause, tx ...utility.DatabaseTransaction) (uint, error)
	Update(documentVersionClause *models.DocumentVersionClause, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error

	GetNextClauseSequence(parentClauseID uint, complianceRegisterID uint, tx ...utility.DatabaseTransaction) (int, error)
	GetNextTopLevelClauseSequence(complianceRegisterID uint, tx ...utility.DatabaseTransaction) (int, error)
	ShiftClausesAfter(documentVersionID uint, afterClauseNo string, tx ...utility.DatabaseTransaction) error
	shiftClauseNumber(clauseNo string, parentPrefix []int, insertPos int) (string, bool)
}

// Repository implementation
type PolicyClauseRepository struct {
	db *gorm.DB
}

// Constructor for PolicyClauseRepository
func NewPolicyClauseRepository(db *gorm.DB) *PolicyClauseRepository {
	return &PolicyClauseRepository{
		db: db,
	}
}

func (repo *PolicyClauseRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *PolicyClauseRepository) GetAll(params *policyDTOs.PolicyClauseListQueryParamsDTO, documentVersionId uint) ([]*models.DocumentVersionClause, error) {
	var clauses []*models.DocumentVersionClause
	query := repo.db.Model(&models.DocumentVersionClause{}).
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		Where("document_version_id = ?", documentVersionId)

	// Dynamic filtering based on the keyword
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
			(document_version_clauses.clause_number LIKE ? OR
			document_version_clauses.clause_title_en LIKE ? OR
			document_version_clauses.clause_title_ar LIKE ?)`,
			keyword, keyword, keyword)
	}

	// Sorting
	if params.SortBy == "" {
		params.SortBy = "clause_number" // Default sort by clause_number if no sortBy is specified
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "asc" // Default order is ascending if not specified
	}
	query = query.Order(params.SortBy + " " + params.Sort)

	if err := query.Find(&clauses).Error; err != nil {
		return nil, err
	}

	return clauses, nil
}

func (repo *PolicyClauseRepository) GetByID(id uint, tx ...utility.DatabaseTransaction) (*models.DocumentVersionClause, error) {
	var clause models.DocumentVersionClause
	db := repo.db

	// Check if transaction is provided and use it
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	err := db.Where("id = ?", id).First(&clause).Error
	if err != nil {
		return nil, err
	}
	return &clause, nil
}

func (repo *PolicyClauseRepository) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.DocumentVersionClause, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var documentVersionClause *models.DocumentVersionClause

	result := db.
		Preload("ParantClause").
		Preload("ChildClauses").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		Preload("DocumentVersion").
		First(&documentVersionClause, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return documentVersionClause, nil
}

func (repo *PolicyClauseRepository) GetByClauseNumber(clauseNumber string, tx ...utility.DatabaseTransaction) (*models.DocumentVersionClause, error) {
	var clause models.DocumentVersionClause
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	err := db.Where("clause_number = ?", clauseNumber).First(&clause).Error
	if err != nil {
		return nil, err
	}
	return &clause, nil
}

func (repo *PolicyClauseRepository) Create(documentVersionClause *models.DocumentVersionClause, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(documentVersionClause).Error; err != nil {
		return 0, err
	}

	return documentVersionClause.ID, nil
}

func (repo *PolicyClauseRepository) Update(documentVersionClause *models.DocumentVersionClause, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(documentVersionClause).Error; err != nil {
		return err
	}

	return nil
}

func (repo *PolicyClauseRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Delete(&models.DocumentVersionClause{}, id).Error; err != nil {
		return err
	}

	return nil
}

func (repo *PolicyClauseRepository) GetNextClauseSequence(parentClauseID uint, documentID uint, tx ...utility.DatabaseTransaction) (int, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var childClauses []string
	query := `
		SELECT clause_number
		FROM document_version_clauses dvc
		INNER JOIN document_versions dv ON dv.id = dvc.document_version_id
		WHERE dv.document_id = ? 
		  AND dv.is_latest_version = 1
		  AND dvc.document_version_clause_id = ?
	`
	if err := db.Raw(query, documentID, parentClauseID).Scan(&childClauses).Error; err != nil {
		return 0, err
	}

	maxSeq := 0
	for _, cn := range childClauses {
		parts := strings.Split(cn, ".")
		if n, err := strconv.Atoi(parts[len(parts)-1]); err == nil && n > maxSeq {
			maxSeq = n
		}
	}

	return maxSeq + 1, nil
}

func (repo *PolicyClauseRepository) GetNextTopLevelClauseSequence(documentID uint, tx ...utility.DatabaseTransaction) (int, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var maxSequence int

	query := `
		SELECT COALESCE(MAX(TRY_CAST(clause_number AS INT)), 0)
		FROM document_version_clauses dvc
		INNER JOIN document_versions dv ON dv.id = dvc.document_version_id
		WHERE dv.document_id = ?
		  AND dv.is_latest_version = 1
		  AND dvc.document_version_clause_id IS NULL
	`
	if err := db.Raw(query, documentID).Row().Scan(&maxSequence); err != nil {
		return 0, err
	}

	return maxSequence + 1, nil
}

func (repo *PolicyClauseRepository) ShiftClausesAfter(documentVersionID uint, afterClauseNo string, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var clauses []struct {
		ID           uint
		ClauseNumber string
	}
	if err := db.Table("document_version_clauses").
		Select("id, clause_number").
		Where("document_version_id = ?", documentVersionID).
		Order("LEN(clause_number) DESC, clause_number DESC").
		Scan(&clauses).Error; err != nil {
		return err
	}

	parentParts := strings.Split(afterClauseNo, ".")
	parentInts := make([]int, len(parentParts))
	for i, p := range parentParts {
		n, _ := strconv.Atoi(p)
		parentInts[i] = n
	}

	// Insert position = last part + 1
	insertPos := parentInts[len(parentInts)-1] + 1

	updates := make(map[uint]string)

	for _, clause := range clauses {
		if newNo, ok := repo.shiftClauseNumber(clause.ClauseNumber, parentInts[:len(parentInts)-1], insertPos); ok {
			updates[clause.ID] = newNo
		}
	}

	// Apply updates
	for id, newNo := range updates {
		if err := db.Model(&models.DocumentVersionClause{}).
			Where("id = ?", id).
			Update("clause_number", newNo).Error; err != nil {
			return err
		}
	}

	return nil
}

func (repo *PolicyClauseRepository) shiftClauseNumber(clauseNo string, parentPrefix []int, insertPos int) (string, bool) {
	parts := strings.Split(clauseNo, ".")
	ints := make([]int, len(parts))
	for i, p := range parts {
		n, _ := strconv.Atoi(p)
		ints[i] = n
	}

	// Must share the same prefix (same parent)
	if len(ints) < len(parentPrefix)+1 {
		return clauseNo, false
	}
	for i := range parentPrefix {
		if ints[i] != parentPrefix[i] {
			return clauseNo, false
		}
	}

	// If sibling index >= insertPos, shift
	if ints[len(parentPrefix)] >= insertPos {
		ints[len(parentPrefix)]++
		newParts := make([]string, len(ints))
		for i, n := range ints {
			newParts[i] = strconv.Itoa(n)
		}
		return strings.Join(newParts, "."), true
	}

	return clauseNo, false
}
