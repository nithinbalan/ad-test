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
	GetAll(params *policyDTOs.PolicyClauseListQueryParamsDTO, documentVersionId uint) ([]*policyDTOs.PolicyClauseListDTO, error)
	GetByID(id uint, tx ...utility.DatabaseTransaction) (*models.DocumentVersionClause, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.DocumentVersionClause, error)
	GetByClauseNumber(clauseNumber string, tx ...utility.DatabaseTransaction) (*models.DocumentVersionClause, error)
	Create(documentVersionClause *models.DocumentVersionClause, tx ...utility.DatabaseTransaction) (uint, error)
	Update(documentVersionClause *models.DocumentVersionClause, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error

	GetNextClauseSequence(parentClauseID uint, documentID uint, tx ...utility.DatabaseTransaction) (int, error)
	GetNextTopLevelClauseSequence(documentID uint, tx ...utility.DatabaseTransaction) (int, error)
	ShiftClausesAfter(documentVersionID uint, afterClauseNo string, tx ...utility.DatabaseTransaction) error
	shiftClauseNumber(clauseNo string, parentPrefix []int, insertPos int) (string, bool)

	HasSubClauses(clauseID uint, tx ...utility.DatabaseTransaction) (bool, error)
	ShiftClauseNumberAfterDeleting(documentVersionID uint, parentClauseID *uint, deletedClauseNumber string, tx ...utility.DatabaseTransaction) error
	shiftSubClauseNumberAfterDeleting(db *gorm.DB, oldPrefix, newPrefix string, parentClauseID uint, documentVersionID uint) error

	CreateClauseHistory(clauseHistory []*models.DocumentVersionClauseHistory, tx ...utility.DatabaseTransaction) error
	GetClauseHistories(clauseID uint) ([]*policyDTOs.ClauseHistoryDTO, error)
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

func (repo *PolicyClauseRepository) GetAll(params *policyDTOs.PolicyClauseListQueryParamsDTO, documentVersionId uint) ([]*policyDTOs.PolicyClauseListDTO, error) {
	var clauses []*policyDTOs.PolicyClauseListDTO

	selectClause := `
		document_version_clauses.*,

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
		(SELECT COUNT(1) 
     FROM document_version_clause_comments c 
     WHERE c.document_version_clause_id = document_version_clauses.id) AS comments_count
	`

	query := repo.db.Table("document_version_clauses").
		Select(selectClause).
		Joins("LEFT JOIN users AS created_user ON created_user.id = document_version_clauses.created_by_id").
		Joins("LEFT JOIN users AS updated_user ON updated_user.id = document_version_clauses.updated_by_id").
		Where("document_version_clauses.document_version_id = ?", documentVersionId)

	// Apply keyword filter
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
		params.SortBy = "clause_number"
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "asc"
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
		Preload("ParentClause").
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

func (repo *PolicyClauseRepository) HasSubClauses(clauseID uint, tx ...utility.DatabaseTransaction) (bool, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var count int64
	err := db.Model(&models.DocumentVersionClause{}).
		Where("document_version_clause_id = ?", clauseID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (repo *PolicyClauseRepository) ShiftClauseNumberAfterDeleting(documentVersionID uint, parentClauseID *uint, deletedClauseNumber string, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	// Build query based on parentClauseID
	query := db.Where("document_version_id = ? AND clause_number > ?", documentVersionID, deletedClauseNumber)

	if parentClauseID == nil {
		query = query.Where("document_version_clause_id IS NULL")
	} else {
		query = query.Where("document_version_clause_id = ?", *parentClauseID)
	}

	var siblings []*models.DocumentVersionClause
	if err := query.Order("clause_number ASC").Find(&siblings).Error; err != nil {
		return err
	}

	for _, sibling := range siblings {
		if sibling.ClauseNumber == nil {
			continue
		}

		// Example: 3.3 → 3.2, 3.4 → 3.3
		parts := strings.Split(*sibling.ClauseNumber, ".")
		last, _ := strconv.Atoi(parts[len(parts)-1])
		parts[len(parts)-1] = strconv.Itoa(last - 1)
		newNumber := strings.Join(parts, ".")
		oldNumber := *sibling.ClauseNumber
		sibling.ClauseNumber = &newNumber

		if err := db.Save(sibling).Error; err != nil {
			return err
		}

		// Also update all subclauses recursively
		if err := repo.shiftSubClauseNumberAfterDeleting(db, oldNumber, newNumber, sibling.ID, documentVersionID); err != nil {
			return err
		}
	}

	return nil
}

func (repo *PolicyClauseRepository) shiftSubClauseNumberAfterDeleting(db *gorm.DB, oldPrefix, newPrefix string, parentClauseID uint, documentVersionID uint) error {
	var subClauses []*models.DocumentVersionClause
	if err := db.Where("document_version_id = ? AND document_version_clause_id = ?", documentVersionID, parentClauseID).
		Find(&subClauses).Error; err != nil {
		return err
	}

	for _, sub := range subClauses {
		if sub.ClauseNumber != nil && strings.HasPrefix(*sub.ClauseNumber, oldPrefix) {
			newNumber := strings.Replace(*sub.ClauseNumber, oldPrefix, newPrefix, 1)
			sub.ClauseNumber = &newNumber
			if err := db.Save(sub).Error; err != nil {
				return err
			}
			// Recurse deeper
			if err := repo.shiftSubClauseNumberAfterDeleting(db, *sub.ClauseNumber, newNumber, sub.ID, documentVersionID); err != nil {
				return err
			}
		}
	}
	return nil
}

func (repo *PolicyClauseRepository) CreateClauseHistory(clauseHistory []*models.DocumentVersionClauseHistory, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(clauseHistory).Error; err != nil {
		return err
	}

	return nil
}

func (repo *PolicyClauseRepository) GetClauseHistories(clauseID uint) ([]*policyDTOs.ClauseHistoryDTO, error) {
	var histories []*policyDTOs.ClauseHistoryDTO
	selectClause := `
		document_version_clause_histories.*,

		document_version_clauses.clause_number AS document_version_clause_number,
		created_user.id AS created_user_id,
		created_user.first_name_en AS created_user_first_name_en,
		created_user.first_name_ar AS created_user_first_name_ar,
		created_user.last_name_en AS created_user_last_name_en,
		created_user.last_name_ar AS created_user_last_name_ar,
		created_user.initial AS created_user_initial,
		created_user.color AS created_user_color,
		created_user.email AS created_user_email,
		created_user.profile_image_path AS created_user_profile_image_path
	`

	query := repo.db.Table("document_version_clause_histories").
		Select(selectClause).
		Joins("LEFT JOIN users AS created_user ON created_user.id = document_version_clause_histories.created_by_id").
		Joins("LEFT JOIN document_version_clauses ON document_version_clauses.id = document_version_clause_histories.document_version_clause_id").
		Where("document_version_clause_histories.document_version_clause_id = ?", clauseID)

	if err := query.Find(&histories).Error; err != nil {
		return nil, err
	}

	return histories, nil
}
