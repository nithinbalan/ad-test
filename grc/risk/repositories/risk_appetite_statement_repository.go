package repositories

import (
	"fmt"
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/models"
	riskDTOs "jethur/grc/risk/dtos"
	"jethur/utility"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type RiskAppetiteStatementRepositoryInterface interface {
	GetAll(params *riskDTOs.RiskAppetiteStatementListQueryParamsDTO) (*riskDTOs.RiskAppetiteStatementListResponseDTO, error)
	GetByID(id uint) (*models.RiskAppetiteStatement, error)
	GetFullDetailsByID(id uint) (*models.RiskAppetiteStatement, error)
	Create(RiskAppetiteStatement *models.RiskAppetiteStatement) (uint, error)
	Update(RiskAppetiteStatement *models.RiskAppetiteStatement) error
	Delete(id uint) error
	Activate(ids []uint) ([]uint, error)
	Deactivate(ids []uint) ([]uint, error)
	ExistsByStatement(StatementEn, StatementAr string, RiskCategoryID uint) (bool, error)
	ExistsByStatementExcludingID(id uint, StatementEn, StatementAr string, RiskCategoryID uint) (bool, error)
	GenerateReferenceNumber() (string, error)
}

type RiskAppetiteStatementRepository struct {
	db *gorm.DB
}

func NewRiskAppetiteStatementRepository(db *gorm.DB) *RiskAppetiteStatementRepository {
	return &RiskAppetiteStatementRepository{
		db: db,
	}
}

func (repo *RiskAppetiteStatementRepository) GetAll(params *riskDTOs.RiskAppetiteStatementListQueryParamsDTO) (*riskDTOs.RiskAppetiteStatementListResponseDTO, error) {

	var riskAppetiteStatements []*riskDTOs.RiskAppetiteStatementListDTO
	var totalCount int64

	query := repo.db.Table("risk_appetite_statements").
		Select(`
            risk_appetite_statements.*,
			risk_categories.title_en AS risk_category_title_en,
			risk_categories.title_ar AS risk_category_title_ar,
			created_user.first_name_en AS created_user_first_name_en,
            created_user.first_name_ar AS created_user_first_name_ar,
            created_user.last_name_en AS created_user_last_name_en,
            created_user.last_name_ar AS created_user_last_name_ar,   
			updated_user.first_name_en AS updated_user_first_name_en,
            updated_user.first_name_ar AS updated_user_first_name_ar,
            updated_user.last_name_en AS updated_user_last_name_en,
            updated_user.last_name_ar AS updated_user_last_name_ar          
        `).
		Joins(`
		    LEFT JOIN risk_categories ON risk_categories.id = risk_appetite_statements.risk_category_id
			LEFT JOIN users AS created_user ON created_user.id = risk_appetite_statements.created_by
			LEFT JOIN users AS updated_user ON updated_user.id = risk_appetite_statements.updated_by
		`)

	// Apply dynamic filters to the query
	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where(columnName+" IN (?)", values)
		}
	}

	// Apply quick filters
	for _, filter := range params.QuickFilter {
		query = query.Where("risk_categories.title_en = ? OR risk_categories.title_ar = ?", filter, filter)
	}

	// Apply search keyword if exists
	if params.Keyword != "" {
		query = query.Where(`
			(risk_appetite_statements.statement_en LIKE ? OR
			risk_appetite_statements.statement_ar LIKE ? OR
			risk_categories.title_en LIKE ? OR
			risk_categories.title_ar LIKE ?)`,
			"%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%")
	}

	// Calculate total count for pagination
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	// Sort by the specified column and order
	if params.SortBy == "" {
		params.SortBy = "risk_appetite_statements.statement_en" // Default sort by risk_appetite_statements.statement_en if no sortBy is specified
	}
	if params.Sort != "asc" && params.Sort != "desc" {
		params.Sort = "asc" // Default order is ascending if not specified
	}

	query = query.Order(params.SortBy + " " + params.Sort)

	offset := (params.Page - 1) * params.PageSize

	// Execute the queryc
	err := query.Offset(offset).Limit(params.PageSize).Find(&riskAppetiteStatements).Error
	if err != nil {
		return nil, err
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	response := &riskDTOs.RiskAppetiteStatementListResponseDTO{
		Items: riskAppetiteStatements,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}

	//@todo The response will be changed based on risk implementation

	return response, nil
}

func (repo *RiskAppetiteStatementRepository) GetByID(id uint) (*models.RiskAppetiteStatement, error) {
	var riskAppetiteStatement *models.RiskAppetiteStatement

	result := repo.db.First(&riskAppetiteStatement, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return riskAppetiteStatement, nil

}

func (repo *RiskAppetiteStatementRepository) GetFullDetailsByID(id uint) (*models.RiskAppetiteStatement, error) {
	var riskAppetiteStatement *models.RiskAppetiteStatement // Declare the variable inside the function

	result := repo.db.
		Preload("RiskCategory").
		First(&riskAppetiteStatement, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return riskAppetiteStatement, nil
}

func (repo *RiskAppetiteStatementRepository) Create(riskAppetiteStatement *models.RiskAppetiteStatement) (uint, error) {

	if err := repo.db.Create(riskAppetiteStatement).Error; err != nil {
		return 0, err
	}

	return riskAppetiteStatement.ID, nil
}

func (repo *RiskAppetiteStatementRepository) Update(riskAppetiteStatement *models.RiskAppetiteStatement) error {

	if err := repo.db.Save(riskAppetiteStatement).Error; err != nil {
		return err
	}

	return nil
}

func (repo *RiskAppetiteStatementRepository) Delete(id uint) error {

	var riskAppetiteStatement *models.RiskAppetiteStatement
	if err := repo.db.First(&riskAppetiteStatement, id).Error; err != nil {
		return err
	}

	result := repo.db.Delete(&models.RiskAppetiteStatement{}, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *RiskAppetiteStatementRepository) Activate(ids []uint) ([]uint, error) {
	var idsToUpdate []uint

	// Step 1: Fetch the IDs of records that are currently 'inactive' and within the provided 'ids' slice
	if err := repo.db.Model(&models.RiskAppetiteStatement{}).
		Where("id IN ? AND status = ?", ids, "inactive").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err // Return error if the query fails
	}

	// Step 2: If there are any IDs to update, perform the update operation
	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.RiskAppetiteStatement{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "active").Error; err != nil {
			return nil, err // Return error if the update operation fails
		}
	}

	// Return the IDs of the records that were activated
	return idsToUpdate, nil
}

func (repo *RiskAppetiteStatementRepository) Deactivate(ids []uint) ([]uint, error) {

	var idsToUpdate []uint

	// Step 1: Fetch the IDs of records that are currently 'inactive' and within the provided 'ids' slice
	if err := repo.db.Model(&models.RiskAppetiteStatement{}).
		Where("id IN ? AND status = ?", ids, "active").
		Pluck("id", &idsToUpdate).Error; err != nil {
		return nil, err // Return error if the query fails
	}

	// Step 2: If there are any IDs to update, perform the update operation
	if len(idsToUpdate) > 0 {
		if err := repo.db.Model(&models.RiskAppetiteStatement{}).
			Where("id IN ?", idsToUpdate).
			Update("status", "inactive").Error; err != nil {
			return nil, err // Return error if the update operation fails
		}
	}

	// Return the IDs of the records that were activated
	return idsToUpdate, nil

}

func (repo *RiskAppetiteStatementRepository) ExistsByStatement(StatementEn, StatementAr string, RiskCategoryID uint) (bool, error) {
	var count int64

	// Check for the existence of Title or QuestionAr in a single query
	err := repo.db.Model(&models.RiskAppetiteStatement{}).
		Where("(statement_en = ? OR statement_ar = ?) AND risk_category_id = ?", StatementEn, StatementAr, RiskCategoryID).
		Count(&count).Error

	if err != nil {
		return false, utility.NewInternalServerError(err)

	}
	// Return true if either title exists, otherwise false
	return count > 0, nil
}

func (repo *RiskAppetiteStatementRepository) ExistsByStatementExcludingID(id uint, StatementEn, StatementAr string, RiskCategoryID uint) (bool, error) {
	var count int64
	err := repo.db.Model(&models.RiskAppetiteStatement{}).
		Where("id <> ? AND (statement_en = ? OR statement_ar = ?) AND risk_category_id = ?", id, StatementEn, StatementAr, RiskCategoryID).
		Count(&count).Error
	if err != nil {
		return false, utility.NewInternalServerError(err)
	}

	return count > 0, nil
}

func (repo *RiskAppetiteStatementRepository) GenerateReferenceNumber() (string, error) {
	var lastReferenceNumber string
	err := repo.db.Model(&models.RiskAppetiteStatement{}).
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
	newReferenceNumber := fmt.Sprintf("%s-%d-%04d", "RAS", currentYear, lastReferenceID+1)

	return newReferenceNumber, nil
}
