package services

import (
	"fmt"
	"jethur/grc/common/repositories"
	"jethur/grc/models"
	"jethur/utility"
)

// Interface definition
type AssociatedItemValidationServiceInterface interface {
	ValidateRiskMapping(masterID uint, riskIDs []uint, masterName string, columnName string) error
	ValidateDocumentMapping(masterID uint, documentIDs []uint, tableName string, columnName string) error
	ValidateIncidentMapping(masterID uint, incidentIDs []uint, tableName string, columnName string) error
	ValidateStrtegicObjectiveMapping(masterID uint, strategicObjectiveIds []uint, tableName string, columnName string) error
	ValidateAssetMapping(masterID uint, assetIDs []uint, tableName string, columnName string) error
	ValidateProcessMapping(masterID uint, processIds []uint, tableName string, columnName string) error
	ValidateKPIMapping(masterID uint, kpiIds []uint, tableName string, columnName string) error
	ValidateControlMapping(masterID uint, controlIds []uint, tableName string, columnName string) error
	ValidateComplianceMapping(masterID uint, complianceRegisterIDs []uint, tableName string, columnName string) error
	ValidateProjectMapping(masterID uint, projectIds []uint, tableName string, columnName string) error
	ValidateBCPMapping(masterID uint, bcppIds []uint, tableName string, columnName string) error
	ValidateClauseIDs(clauseIDs []uint) ([]*models.ComplianceClause, error)
}

// Service struct
type AssociatedItemValidationService struct {
	repo repositories.AssociatedItemValidationRepositoryInterface
}

// Constructor
func NewAssociatedItemValidationService(repo repositories.AssociatedItemValidationRepositoryInterface) AssociatedItemValidationServiceInterface {
	return &AssociatedItemValidationService{
		repo: repo,
	}
}

// Main validation logic
func (aivs *AssociatedItemValidationService) ValidateRiskMapping(masterID uint, riskIDs []uint, masterName string, columnName string) error {
	if len(riskIDs) == 0 {
		return utility.NewValidationError(fmt.Errorf("no risks provided"), "required", "risk_ids")
	}

	seen := make(map[uint]bool)
	for _, id := range riskIDs {
		if seen[id] {
			return utility.NewValidationError(fmt.Errorf("duplicate risk ID: %d", id), "duplicate", "risk_ids")
		}

		_, err := aivs.repo.RiskExists(id)
		if err != nil {
			return utility.NewValidationError(fmt.Errorf("risk does not exist: %d", id), "not_exist", "risk_ids")
		}

		alreadyMapped, err := aivs.repo.IsRiskMapped(masterName, columnName, id, masterID)
		if err != nil {
			return fmt.Errorf("failed to check mapping: %v", err)
		}
		if alreadyMapped {
			return utility.NewValidationError(
				fmt.Errorf("risk ID %d already mapped to this %s", id, masterName),
				"duplicate", "risk_ids")
		}

		seen[id] = true
	}

	return nil
}

func (aivs *AssociatedItemValidationService) ValidateDocumentMapping(masterID uint, documentIDs []uint, tableName string, columnName string) error {
	if len(documentIDs) == 0 {
		return utility.NewValidationError(fmt.Errorf("no document provided"), "required", "document_ids")
	}

	seen := make(map[uint]bool)
	for _, id := range documentIDs {
		if seen[id] {
			return utility.NewValidationError(fmt.Errorf("duplicate document ID: %d", id), "duplicate", "document_ids")
		}

		_, status, err := aivs.repo.DocumentExists(id)
		if err != nil {
			return utility.NewValidationError(fmt.Errorf("document does not exist: %d", id), "not_exist", "document_ids")
		}

		if status != "published" {
			return utility.NewValidationError(fmt.Errorf("document not published : %d", id), "document_not_published", "document_id")

		}

		alreadyMapped, err := aivs.repo.IsDocumentMapped(tableName, columnName, id, masterID)
		if err != nil {
			return fmt.Errorf("failed to check mapping: %v", err)
		}
		if alreadyMapped {
			return utility.NewValidationError(
				fmt.Errorf("document ID %d already mapped to this %s", id, tableName),
				"duplicate", "document")
		}

		seen[id] = true
	}

	return nil
}

func (aivs *AssociatedItemValidationService) ValidateIncidentMapping(masterID uint, incidentIDs []uint, tableName string, columnName string) error {
	if len(incidentIDs) == 0 {
		return utility.NewValidationError(fmt.Errorf("no incidents provided"), "required", "incident_ids")
	}

	seen := make(map[uint]bool)
	for _, id := range incidentIDs {
		if seen[id] {
			return utility.NewValidationError(fmt.Errorf("duplicate incident ID: %d", id), "duplicate", "incident_ids")
		}

		_, err := aivs.repo.IncidentExists(id)
		if err != nil {
			return utility.NewValidationError(fmt.Errorf("incident does not exist: %d", id), "not_exist", "incident_ids")
		}

		alreadyMapped, err := aivs.repo.IsIncidentMapped(tableName, columnName, id, masterID)
		if err != nil {
			return fmt.Errorf("failed to check mapping: %v", err)
		}
		if alreadyMapped {
			return utility.NewValidationError(
				fmt.Errorf("document ID %d already mapped to this %s", id, tableName),
				"duplicate", "incident")
		}

		seen[id] = true
	}

	return nil
}

func (aivs *AssociatedItemValidationService) ValidateStrtegicObjectiveMapping(masterID uint, strategicObjectiveIds []uint, tableName string, columnName string) error {
	if len(strategicObjectiveIds) == 0 {
		return utility.NewValidationError(fmt.Errorf("no strategic objective provided"), "required", "strategic_objective_id")
	}

	seen := make(map[uint]bool)
	for _, id := range strategicObjectiveIds {
		if seen[id] {
			return utility.NewValidationError(fmt.Errorf("duplicate strategic objective ID: %d", id), "duplicate", "strategic_objective_id")
		}

		_, err := aivs.repo.StrategicObjectiveExists(id)
		if err != nil {
			return utility.NewValidationError(fmt.Errorf("strategic objective does not exist: %d", id), "not_exist", "strategic_objective_id")
		}

		alreadyMapped, err := aivs.repo.IsstrategicObjectiveMapped(tableName, columnName, id, masterID)
		if err != nil {
			return fmt.Errorf("failed to check mapping: %v", err)
		}
		if alreadyMapped {
			return utility.NewValidationError(
				fmt.Errorf("strategic objective ID %d already mapped to this %s", id, tableName),
				"duplicate", "strategic objective")
		}

		seen[id] = true
	}

	return nil
}

func (aivs *AssociatedItemValidationService) ValidateAssetMapping(masterID uint, assetIDs []uint, tableName string, columnName string) error {
	if len(assetIDs) == 0 {
		return utility.NewValidationError(fmt.Errorf("no asset provided"), "required", "asset_id")
	}

	seen := make(map[uint]bool)
	for _, id := range assetIDs {
		if seen[id] {
			return utility.NewValidationError(fmt.Errorf("duplicate asset ID: %d", id), "duplicate", "asset_id")
		}

		_, err := aivs.repo.AssetExists(id)
		if err != nil {
			return utility.NewValidationError(fmt.Errorf("asset does not exist: %d", id), "not_exist", "asset_id")
		}

		alreadyMapped, err := aivs.repo.IsAssetMapped(tableName, columnName, id, masterID)
		if err != nil {
			return fmt.Errorf("failed to check mapping: %v", err)
		}
		if alreadyMapped {
			return utility.NewValidationError(
				fmt.Errorf("asset ID %d already mapped to this %s", id, tableName),
				"duplicate", "asset")
		}

		seen[id] = true
	}

	return nil
}

func (aivs *AssociatedItemValidationService) ValidateProcessMapping(masterID uint, processIds []uint, tableName string, columnName string) error {
	if len(processIds) == 0 {
		return utility.NewValidationError(fmt.Errorf("no process provided"), "required", "procees_id")
	}

	seen := make(map[uint]bool)
	for _, id := range processIds {
		if seen[id] {
			return utility.NewValidationError(fmt.Errorf("duplicate process ID: %d", id), "duplicate", "procees_id")
		}

		_, err := aivs.repo.ProcessExists(id)
		if err != nil {
			return utility.NewValidationError(fmt.Errorf("process does not exist: %d", id), "not_exist", "procees_id")
		}

		alreadyMapped, err := aivs.repo.IsProcessMapped(tableName, columnName, id, masterID)
		if err != nil {
			return fmt.Errorf("failed to check mapping: %v", err)
		}
		if alreadyMapped {
			return utility.NewValidationError(
				fmt.Errorf("process ID %d already mapped to this %s", id, tableName),
				"duplicate", "process")
		}

		seen[id] = true
	}

	return nil
}

func (aivs *AssociatedItemValidationService) ValidateKPIMapping(masterID uint, kpiIds []uint, tableName string, columnName string) error {
	if len(kpiIds) == 0 {
		return utility.NewValidationError(fmt.Errorf("no kpi provided"), "required", "kpi_id")
	}

	seen := make(map[uint]bool)
	for _, id := range kpiIds {
		if seen[id] {
			return utility.NewValidationError(fmt.Errorf("duplicate kpi ID: %d", id), "duplicate", "kpi_id")
		}

		_, err := aivs.repo.KpiExists(id)
		if err != nil {
			return utility.NewValidationError(fmt.Errorf("kpi does not exist: %d", id), "not_exist", "kpi_id")
		}

		alreadyMapped, err := aivs.repo.IsKPIMapped(tableName, columnName, id, masterID)
		if err != nil {
			return fmt.Errorf("failed to check mapping: %v", err)
		}
		if alreadyMapped {
			return utility.NewValidationError(
				fmt.Errorf("kpi ID %d already mapped to this %s", id, tableName),
				"duplicate", "KPI")
		}

		seen[id] = true
	}

	return nil
}

func (aivs *AssociatedItemValidationService) ValidateControlMapping(masterID uint, controlIds []uint, tableName string, columnName string) error {
	if len(controlIds) == 0 {
		return utility.NewValidationError(fmt.Errorf("no controls provided"), "required", "control_id")
	}

	seen := make(map[uint]bool)
	for _, id := range controlIds {
		if seen[id] {
			return utility.NewValidationError(fmt.Errorf("duplicate contol ID: %d", id), "duplicate", "control_id")
		}

		_, err := aivs.repo.ControlExists(id)
		if err != nil {
			return utility.NewValidationError(fmt.Errorf("control does not exist: %d", id), "not_exist", "control_id")
		}

		alreadyMapped, err := aivs.repo.IsControlMapped(tableName, columnName, id, masterID)
		if err != nil {
			return fmt.Errorf("failed to check mapping: %v", err)
		}
		if alreadyMapped {
			return utility.NewValidationError(
				fmt.Errorf("control ID %d already mapped to this %s", id, tableName),
				"duplicate", "Control")
		}

		seen[id] = true
	}

	return nil
}

func (aivs *AssociatedItemValidationService) ValidateComplianceMapping(masterID uint, complianceRegisterIDs []uint, tableName string, columnName string) error {
	if len(complianceRegisterIDs) == 0 {
		return utility.NewValidationError(fmt.Errorf("no compliance registers provided"), "required", "clause_ids")
	}

	seen := make(map[uint]bool)
	for _, id := range complianceRegisterIDs {
		if seen[id] {
			return utility.NewValidationError(fmt.Errorf("duplicate clause ID: %d", id), "duplicate", "clause_ids")
		}

		_, err := aivs.repo.ComplianceRegisterExists(id)
		if err != nil {
			return utility.NewValidationError(fmt.Errorf("clause does not exist: %d", id), "not_exist", "clause_ids")
		}

		alreadyMapped, err := aivs.repo.IsComplianceRegisterMapped(tableName, columnName, id, masterID)
		if err != nil {
			return fmt.Errorf("failed to check mapping: %v", err)
		}
		if alreadyMapped {
			return utility.NewValidationError(
				fmt.Errorf("clause ID %d already mapped to this %s", id, tableName),
				"duplicate", "clause_ids")
		}

		seen[id] = true
	}

	return nil
}

func (aivs *AssociatedItemValidationService) ValidateProjectMapping(masterID uint, projectIds []uint, tableName string, columnName string) error {
	if len(projectIds) == 0 {
		return utility.NewValidationError(fmt.Errorf("no projects provided"), "required", "project_id")
	}

	seen := make(map[uint]bool)
	for _, id := range projectIds {
		if seen[id] {
			return utility.NewValidationError(fmt.Errorf("duplicate project ID: %d", id), "duplicate", "project_id")
		}

		_, err := aivs.repo.ProjectExists(id)
		if err != nil {
			return utility.NewValidationError(fmt.Errorf("project does not exist: %d", id), "not_exist", "project_id")
		}

		alreadyMapped, err := aivs.repo.IsProjectMapped(tableName, columnName, id, masterID)
		if err != nil {
			return fmt.Errorf("failed to check mapping: %v", err)
		}
		if alreadyMapped {
			return utility.NewValidationError(
				fmt.Errorf("project ID %d already mapped to this %s", id, tableName),
				"duplicate", "Project")
		}

		seen[id] = true
	}

	return nil
}

func (aivs *AssociatedItemValidationService) ValidateBCPMapping(masterID uint, bcppIds []uint, tableName string, columnName string) error {
	if len(bcppIds) == 0 {
		return utility.NewValidationError(fmt.Errorf("no bcp provided"), "required", "bcp_id")
	}

	seen := make(map[uint]bool)
	for _, id := range bcppIds {
		if seen[id] {
			return utility.NewValidationError(fmt.Errorf("duplicate bcp ID: %d", id), "duplicate", "bcp_id")
		}

		_, err := aivs.repo.BCPExists(id)
		if err != nil {
			return utility.NewValidationError(fmt.Errorf("bcp does not exist: %d", id), "not_exist", "bcp_id")
		}

		alreadyMapped, err := aivs.repo.IsBCPMapped(tableName, columnName, id, masterID)
		if err != nil {
			return fmt.Errorf("failed to check mapping: %v", err)
		}
		if alreadyMapped {
			return utility.NewValidationError(
				fmt.Errorf("bcp ID %d already mapped to this %s", id, tableName),
				"duplicate", "BCP")
		}

		seen[id] = true
	}

	return nil
}

func (aivs *AssociatedItemValidationService) ValidateClauseIDs(clauseIDs []uint) ([]*models.ComplianceClause, error) {
	seen := make(map[uint]bool)
	for _, id := range clauseIDs {
		if seen[id] {
			return nil, utility.NewValidationError(fmt.Errorf("duplicate clause ID: %d", id), "duplicate", "clause_id")
		}
		seen[id] = true
	}

	clauses, err := aivs.repo.GetClausesByIDs(clauseIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch clauses: %w", err)
	}

	foundClauseIDs := make(map[uint]bool)
	for _, clause := range clauses {
		foundClauseIDs[clause.ID] = true
	}

	for _, id := range clauseIDs {
		if !foundClauseIDs[id] {
			return nil, utility.NewValidationError(
				fmt.Errorf("clause does not exist for ID: %d", id), "not_exist", "clause_id")
		}
	}

	return clauses, nil
}
