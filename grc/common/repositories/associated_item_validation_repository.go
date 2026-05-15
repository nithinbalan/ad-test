package repositories

import (
	"jethur/grc/models"

	"gorm.io/gorm"
)

type AssociatedItemValidationRepositoryInterface interface {
	IsRiskMapped(tableName string, masterColumn string, riskID, masterID uint) (bool, error)
	RiskExists(id uint) (*models.RiskRegister, error)
	IsDocumentMapped(masterName string, columnName string, documentID, masterID uint) (bool, error)
	DocumentExists(id uint) (*models.Document, string, error)
	IsIncidentMapped(masterName string, columnName string, incidentID, masterID uint) (bool, error)
	IncidentExists(id uint) (*models.IncidentRegister, error)
	IsstrategicObjectiveMapped(masterName string, columnName string, strategicObjectiveID, masterID uint) (bool, error)
	StrategicObjectiveExists(id uint) (*models.StrategicObjective, error)
	IsAssetMapped(masterName string, columnName string, assetID, masterID uint) (bool, error)
	AssetExists(id uint) (*models.AssetRegister, error)
	IsSLAMapped(masterName string, columnName string, slaID, masterID uint) (bool, error)
	SLAExists(id uint) (*models.SLARegister, error)
	IsProcessMapped(masterName string, columnName string, processID, masterID uint) (bool, error)
	ProcessExists(id uint) (*models.Process, error)
	IsKPIMapped(masterName string, columnName string, kpiID, masterID uint) (bool, error)
	KpiExists(id uint) (*models.KPIRegister, error)
	IsControlMapped(masterName string, columnName string, controlId, masterID uint) (bool, error)
	ControlExists(id uint) (*models.Control, error)
	IsComplianceRegisterMapped(masterName string, columnName string, complianceRegisterID, masterID uint) (bool, error)
	ComplianceRegisterExists(id uint) (*models.ComplianceRegister, error)
	IsProjectMapped(masterName string, columnName string, projectId, masterID uint) (bool, error)
	ProjectExists(id uint) (*models.Project, error)
	IsBCPMapped(masterName string, columnName string, bcpId, masterID uint) (bool, error)
	BCPExists(id uint) (*models.BusinessContinuityPlan, error)
	GetClausesByIDs(ids []uint) ([]*models.ComplianceClause, error)
}

type AssociatedItemValidationRepository struct {
	db *gorm.DB
}

func NewRiskMappingRepository(db *gorm.DB) AssociatedItemValidationRepositoryInterface {
	return &AssociatedItemValidationRepository{
		db: db,
	}
}

// Generic function to check if a risk is already mapped to a master entity
func (repo *AssociatedItemValidationRepository) IsRiskMapped(masterName string, columnName string, riskID, masterID uint) (bool, error) {
	var count int64
	err := repo.db.
		Table(masterName).
		Where("risk_register_id = ? AND "+columnName+" = ?", riskID, masterID).
		Count(&count).Error

	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (repo *AssociatedItemValidationRepository) RiskExists(id uint) (*models.RiskRegister, error) {
	db := repo.db

	var risk *models.RiskRegister

	result := db.
		First(&risk, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return risk, nil

}

func (repo *AssociatedItemValidationRepository) IsDocumentMapped(masterName string, columnName string, documentID, masterID uint) (bool, error) {
	var count int64
	err := repo.db.
		Table(masterName).
		Where("document_id = ? AND "+columnName+" = ?", documentID, masterID).
		Count(&count).Error

	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// func (repo *AssociatedItemValidationRepository) DocumentExists(id uint) (*models.Document, error) {
// 	db := repo.db

// 	var document *models.Document

// 	result := db.
// 		First(&document, id)
// 	if result.Error != nil {
// 		return nil, result.Error
// 	}

// 	return document, nil

// }

func (repo *AssociatedItemValidationRepository) DocumentExists(id uint) (*models.Document, string, error) {
	db := repo.db

	var document models.Document
	result := db.First(&document, id)
	if result.Error != nil {
		return nil, "", result.Error
	}

	var statusType string
	err := db.
		Table("document_statuses").
		Select("type").
		Where("id = ?", document.DocumentStatusID).
		Scan(&statusType).Error
	if err != nil {
		return nil, "", err
	}

	return &document, statusType, nil
}

func (repo *AssociatedItemValidationRepository) IsIncidentMapped(masterName string, columnName string, incidentID, masterID uint) (bool, error) {
	var count int64
	err := repo.db.
		Table(masterName).
		Where("incident_register_id = ? AND "+columnName+" = ?", incidentID, masterID).
		Count(&count).Error

	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (repo *AssociatedItemValidationRepository) IncidentExists(id uint) (*models.IncidentRegister, error) {
	db := repo.db

	var incident *models.IncidentRegister

	result := db.
		First(&incident, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return incident, nil

}

func (repo *AssociatedItemValidationRepository) IsstrategicObjectiveMapped(masterName string, columnName string, strategicObjectiveID, masterID uint) (bool, error) {
	var count int64
	err := repo.db.
		Table(masterName).
		Where("strategic_objective_id = ? AND "+columnName+" = ?", strategicObjectiveID, masterID).
		Count(&count).Error

	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (repo *AssociatedItemValidationRepository) StrategicObjectiveExists(id uint) (*models.StrategicObjective, error) {
	db := repo.db

	var strategicObj *models.StrategicObjective

	result := db.
		First(&strategicObj, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return strategicObj, nil

}

func (repo *AssociatedItemValidationRepository) IsAssetMapped(masterName string, columnName string, assetID, masterID uint) (bool, error) {
	var count int64
	err := repo.db.
		Table(masterName).
		Where("asset_register_id = ? AND "+columnName+" = ?", assetID, masterID).
		Count(&count).Error

	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (repo *AssociatedItemValidationRepository) AssetExists(id uint) (*models.AssetRegister, error) {
	db := repo.db

	var asset *models.AssetRegister

	result := db.
		First(&asset, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return asset, nil

}

func (repo *AssociatedItemValidationRepository) SLAExists(id uint) (*models.SLARegister, error) {
	db := repo.db

	var sla *models.SLARegister

	result := db.
		First(&sla, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return sla, nil
}

func (repo *AssociatedItemValidationRepository) IsSLAMapped(masterName string, columnName string, slaID, masterID uint) (bool, error) {
	var count int64
	err := repo.db.
		Table(masterName).
		Where("sla_register_id = ? AND "+columnName+" = ?", slaID, masterID).
		Count(&count).Error

	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (repo *AssociatedItemValidationRepository) IsProcessMapped(masterName string, columnName string, processID, masterID uint) (bool, error) {
	var count int64
	err := repo.db.
		Table(masterName).
		Where("process_id = ? AND "+columnName+" = ?", processID, masterID).
		Count(&count).Error

	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (repo *AssociatedItemValidationRepository) ProcessExists(id uint) (*models.Process, error) {
	db := repo.db

	var process *models.Process

	result := db.
		First(&process, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return process, nil

}

func (repo *AssociatedItemValidationRepository) IsKPIMapped(masterName string, columnName string, kpiID, masterID uint) (bool, error) {
	var count int64
	err := repo.db.
		Table(masterName).
		Where("kpi_register_id = ? AND "+columnName+" = ?", kpiID, masterID).
		Count(&count).Error

	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (repo *AssociatedItemValidationRepository) KpiExists(id uint) (*models.KPIRegister, error) {
	db := repo.db

	var kpi *models.KPIRegister

	result := db.
		First(&kpi, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return kpi, nil

}

func (repo *AssociatedItemValidationRepository) IsControlMapped(masterName string, columnName string, controlId, masterID uint) (bool, error) {
	var count int64
	err := repo.db.
		Table(masterName).
		Where("control_id = ? AND "+columnName+" = ?", controlId, masterID).
		Count(&count).Error

	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (repo *AssociatedItemValidationRepository) ControlExists(id uint) (*models.Control, error) {
	db := repo.db

	var contol *models.Control

	result := db.
		First(&contol, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return contol, nil

}

func (repo *AssociatedItemValidationRepository) IsComplianceRegisterMapped(masterName string, columnName string, complianceRegisterID, masterID uint) (bool, error) {
	var count int64
	err := repo.db.
		Table(masterName).
		Where("compliance_register_id = ? AND "+columnName+" = ?", complianceRegisterID, masterID).
		Count(&count).Error

	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (repo *AssociatedItemValidationRepository) ComplianceRegisterExists(id uint) (*models.ComplianceRegister, error) {
	db := repo.db

	var complianceRegister *models.ComplianceRegister

	result := db.
		First(&complianceRegister, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return complianceRegister, nil

}

func (repo *AssociatedItemValidationRepository) IsProjectMapped(masterName string, columnName string, projectId, masterID uint) (bool, error) {
	var count int64
	err := repo.db.
		Table(masterName).
		Where("project_id = ? AND "+columnName+" = ?", projectId, masterID).
		Count(&count).Error

	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (repo *AssociatedItemValidationRepository) ProjectExists(id uint) (*models.Project, error) {
	db := repo.db

	var project *models.Project

	result := db.
		First(&project, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return project, nil

}

func (repo *AssociatedItemValidationRepository) IsBCPMapped(masterName string, columnName string, bcpId, masterID uint) (bool, error) {
	var count int64
	err := repo.db.
		Table(masterName).
		Where("business_continuity_plan_id = ? AND "+columnName+" = ?", bcpId, masterID).
		Count(&count).Error

	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (repo *AssociatedItemValidationRepository) BCPExists(id uint) (*models.BusinessContinuityPlan, error) {
	db := repo.db

	var bcp *models.BusinessContinuityPlan

	result := db.
		First(&bcp, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return bcp, nil

}
func (repo *AssociatedItemValidationRepository) GetClausesByIDs(ids []uint) ([]*models.ComplianceClause, error) {
	db := repo.db
	var clauses []*models.ComplianceClause
	if err := db.
		Select("id, compliance_register_id").
		Where("id IN ?", ids).
		Find(&clauses).Error; err != nil {
		return nil, err
	}

	return clauses, nil
}
