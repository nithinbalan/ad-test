package repositories

import (
	"jethur/grc/models"

	"gorm.io/gorm"
)

type AssociatedItemValidationRepositoryInterface interface {
	IsRiskMapped(tableName string, masterColumn string, riskID, masterID uint) (bool, error)
	RiskExists(id uint) (*models.RiskRegister, error)
	IsDocumentMapped(masterName string, columnName string, documentID, masterID uint) (bool, error)
	DocumentExists(id uint) (*models.Document, error)
	IsIncidentMapped(masterName string, columnName string, incidentID, masterID uint) (bool, error)
	IncidentExists(id uint) (*models.IncidentRegister, error)
	IsstrategicObjectiveMapped(masterName string, columnName string, strategicObjectiveID, masterID uint) (bool, error)
	StrategicObjectiveExists(id uint) (*models.StrategicObjective, error)
	IsAssetMapped(masterName string, columnName string, assetID, masterID uint) (bool, error)
	AssetExists(id uint) (*models.AssetRegister, error)
	IsProcessMapped(masterName string, columnName string, processID, masterID uint) (bool, error)
	ProcessExists(id uint) (*models.Process, error)
	IsKPIMapped(masterName string, columnName string, kpiID, masterID uint) (bool, error)
	KpiExists(id uint) (*models.KPIRegister, error)
	IsControlMapped(masterName string, columnName string, controlId, masterID uint) (bool, error)
	ControlExists(id uint) (*models.Control, error)
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

func (repo *AssociatedItemValidationRepository) DocumentExists(id uint) (*models.Document, error) {
	db := repo.db

	var document *models.Document

	result := db.
		First(&document, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return document, nil

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
