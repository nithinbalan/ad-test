package repositories

import (
	"jethur/grc/models" // Ensure this import path is correct for your project structure

	"gorm.io/gorm"
)

// Interfaces
type WorkflowLevelTypeRepositoryInterface interface {
	GetAll(sort string, sortBy string, keyword string, status string, moduleIdentifier string) ([]*models.WorkflowLevelType, error)
	GetByID(id uint) (*models.WorkflowLevelType, error)
	GetByType(levelType string) (*models.WorkflowLevelType, error)
}

// Implementation
type WorkflowLevelTypeRepository struct {
	db *gorm.DB
}

func NewWorkflowLevelTypeRepository(db *gorm.DB) *WorkflowLevelTypeRepository {
	return &WorkflowLevelTypeRepository{
		db: db,
	}
}

// moduleIdentifierTokenCondition builds a SQL condition (and its args) that matches
// the given identifier as a standalone token within a comma-separated module_identifier
// column value (e.g. "organization,compliance" matches identifier "organization").
func moduleIdentifierTokenCondition(identifier string) (string, []interface{}) {
	cond := "(module_identifier = ? OR module_identifier LIKE ? OR module_identifier LIKE ? OR module_identifier LIKE ?)"
	args := []interface{}{
		identifier,
		identifier + ",%",
		"%," + identifier,
		"%," + identifier + ",%",
	}
	return cond, args
}

func (repo *WorkflowLevelTypeRepository) GetAll(sort string, sortBy string, keyword string, status string, moduleIdentifier string) ([]*models.WorkflowLevelType, error) {
	var workflowLevelTypes []*models.WorkflowLevelType

	query := repo.db

	if keyword != "" {
		query = query.Where("(title_en LIKE ? OR title_ar LIKE ?)", "%"+keyword+"%", "%"+keyword+"%")
	}
	if sortBy == "" {
		sortBy = "title_en"
	}
	if sort != "asc" && sort != "desc" {
		sort = "asc"
	}

	if status == "active" || status == "inactive" {
		query = query.Where("status = ?", status)
	}

	if moduleIdentifier == "" {
		moduleIdentifier = "common"
	}

	switch moduleIdentifier {
	case "common":
		cond, args := moduleIdentifierTokenCondition("common")
		query = query.Where(cond, args...)
	case "organization":
		commonCond, commonArgs := moduleIdentifierTokenCondition("common")
		orgCond, orgArgs := moduleIdentifierTokenCondition("organization")
		query = query.Where("("+commonCond+" OR "+orgCond+")", append(commonArgs, orgArgs...)...)

	case "compliance":
		commonCond, commonArgs := moduleIdentifierTokenCondition("common")
		complianceCond, complianceArgs := moduleIdentifierTokenCondition("compliance")
		query = query.Where("("+commonCond+" OR "+complianceCond+")", append(commonArgs, complianceArgs...)...)
	}

	query = query.Order(sortBy + " " + sort)
	result := query.Find(&workflowLevelTypes)

	return workflowLevelTypes, result.Error
}

func (repo *WorkflowLevelTypeRepository) GetByID(id uint) (*models.WorkflowLevelType, error) {
	var levelType *models.WorkflowLevelType
	result := repo.db.First(&levelType, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return levelType, nil
}

func (repo *WorkflowLevelTypeRepository) GetByType(levelType string) (*models.WorkflowLevelType, error) {
	var workflowLevelType *models.WorkflowLevelType

	result := repo.db.Where("type = ?", levelType).First(&workflowLevelType)

	if result.Error != nil {
		return nil, result.Error
	}

	return workflowLevelType, nil
}
