package services

import (
	"fmt"
	commonDTOs "jethur/grc/common/dtos"
	commonServices "jethur/grc/common/services"
	"jethur/grc/models"
	"jethur/grc/project/repositories"
	"jethur/utility"
	"time"
)

type ProjectAssociatedItemServiceInterface interface {
	GetAllProjectDocuments(projectID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedDocumentListResponseDTO, error)
	CreateProjectDocument(projectID uint, userID uint, dto *commonDTOs.AssociatedDocumentMapDTO) error
	DeleteProjectDocument(projectID uint, documentID uint) error
	GetAllProjectControls(projectID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedControlListResponseDTO, error)
	CreateProjectControl(projectID uint, userID uint, dto *commonDTOs.AssociatedControlMapDTO) error
	DeleteProjectContol(projectID uint, controlID uint) error
	GetAllProjectProcesses(projectID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedItemProcessListResponseDTO, error)
	CreateProjectProcess(projectID uint, userID uint, dto *commonDTOs.AssociatedProcessMapDTO) error
	DeleteProjectProcess(projectID uint, processID uint) error
	GetAllProjectStrategicObjectives(projectID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedItemStrategicObjectiveListResponseDTO, error)
	CreateProjectStrategicObjective(projectID uint, userID uint, dto *commonDTOs.AssociatedStrategicObjectiveMapDTO) error
	DeleteProjectStrategicObjective(projectID uint, strategicObjectiveId uint) error
	GetAllProjectKPI(projectID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedKPIRegisterListResponseDTO, error)
	CreateProjectKPI(projectID uint, userID uint, dto *commonDTOs.AssociatedKPIMapDTO) error
	DeleteProjectKPI(projectID uint, kpiID uint) error
	GetAllProjectAssets(projectID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedItemAssetListResponseDTO, error)
	CreateProjectAsset(projectID uint, userID uint, dto *commonDTOs.AssociatedAssetMapDTO) error
	DeleteProjectAsset(projectID uint, kpiID uint) error
	GetAllProjectIncidents(projectID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedIncidentListResponseDTO, error)
	CreateProjectIncident(projectID uint, userID uint, dto *commonDTOs.AssociatedIncidentMapDTO) error
	DeleteProjectIncident(projectID uint, incidentID uint) error
}

// Service and implementations
type ProjectAssociatedItemService struct {
	repo          repositories.ProjectAssociatedItemRepositoryInterface
	projectSvc    ProjectServiceInterface
	validationSvc commonServices.AssociatedItemValidationServiceInterface
}

func NewProjectAssociatedItemService(repo repositories.ProjectAssociatedItemRepositoryInterface,
	projectSvc ProjectServiceInterface,
	validationSvc commonServices.AssociatedItemValidationServiceInterface,
) ProjectAssociatedItemServiceInterface {
	return &ProjectAssociatedItemService{
		repo:          repo,
		projectSvc:    projectSvc,
		validationSvc: validationSvc,
	}
}

func (pais *ProjectAssociatedItemService) GetAllProjectDocuments(assetID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedDocumentListResponseDTO, error) {
	if _, err := pais.projectSvc.GetByID(assetID); err != nil {
		return nil, err
	}
	return pais.repo.GetAllProjectDocuments(assetID, params)
}

func (pais *ProjectAssociatedItemService) CreateProjectDocument(projectID, userID uint, dto *commonDTOs.AssociatedDocumentMapDTO) error {

	if _, err := pais.projectSvc.GetByID(projectID); err != nil {
		return utility.NewValidationError(fmt.Errorf("project does not exist for ID: %d", projectID), "not_exist", "asset_id")
	}

	if err := pais.validationSvc.ValidateDocumentMapping(projectID, dto.DocumentIDs, "project_associated_documents", "project_id"); err != nil {
		return err
	}

	tx := pais.repo.BeginTransaction()
	var mappings []*models.ProjectAssociatedDocument
	currentTime := time.Now()

	for _, documentID := range dto.DocumentIDs {
		mappings = append(mappings, &models.ProjectAssociatedDocument{
			DocumentID:  documentID,
			ProjectID:   projectID,
			CreatedByID: userID,
			CreatedAt:   currentTime,
		})
	}

	if err := pais.repo.CreateProjectDocument(mappings, tx); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (pais *ProjectAssociatedItemService) DeleteProjectDocument(projectID uint, documentID uint) error {
	// Validate control existence
	if _, err := pais.projectSvc.GetByID(projectID); err != nil {
		return utility.NewValidationError(fmt.Errorf("project does not exist for ID: %d", projectID), "not_exist", "asset_register_id")
	}

	tx := pais.repo.BeginTransaction()

	if err := pais.repo.DeleteProjectDocument(projectID, documentID, tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
func (pais *ProjectAssociatedItemService) GetAllProjectControls(assetID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedControlListResponseDTO, error) {
	if _, err := pais.projectSvc.GetByID(assetID); err != nil {
		return nil, err
	}
	return pais.repo.GetAllProjectControls(assetID, params)
}

func (pais *ProjectAssociatedItemService) CreateProjectControl(projectID, userID uint, dto *commonDTOs.AssociatedControlMapDTO) error {

	if _, err := pais.projectSvc.GetByID(projectID); err != nil {
		return utility.NewValidationError(fmt.Errorf("project does not exist for ID: %d", projectID), "not_exist", "asset_id")
	}

	if err := pais.validationSvc.ValidateControlMapping(projectID, dto.ControlIDs, "project_controls", "project_id"); err != nil {
		return err
	}

	tx := pais.repo.BeginTransaction()
	var mappings []*models.ProjectControl
	currentTime := time.Now()

	for _, controlID := range dto.ControlIDs {
		mappings = append(mappings, &models.ProjectControl{
			ControlID:   controlID,
			ProjectID:   projectID,
			CreatedBYID: userID,
			CreatedAt:   currentTime,
		})
	}

	if err := pais.repo.CreateProjectControl(mappings, tx); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (pais *ProjectAssociatedItemService) DeleteProjectContol(projectID uint, controlID uint) error {
	// Validate control existence
	if _, err := pais.projectSvc.GetByID(projectID); err != nil {
		return utility.NewValidationError(fmt.Errorf("project does not exist for ID: %d", projectID), "not_exist", "asset_register_id")
	}

	tx := pais.repo.BeginTransaction()

	if err := pais.repo.DeleteProjectContol(projectID, controlID, tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (pais *ProjectAssociatedItemService) GetAllProjectProcesses(assetID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedItemProcessListResponseDTO, error) {
	if _, err := pais.projectSvc.GetByID(assetID); err != nil {
		return nil, err
	}
	return pais.repo.GetAllProjectProcesses(assetID, params)
}

func (pais *ProjectAssociatedItemService) CreateProjectProcess(projectID, userID uint, dto *commonDTOs.AssociatedProcessMapDTO) error {

	if _, err := pais.projectSvc.GetByID(projectID); err != nil {
		return utility.NewValidationError(fmt.Errorf("project does not exist for ID: %d", projectID), "not_exist", "asset_id")
	}

	if err := pais.validationSvc.ValidateProcessMapping(projectID, dto.ProcessIDs, "project_processes", "project_id"); err != nil {
		return err
	}

	tx := pais.repo.BeginTransaction()
	var mappings []*models.ProjectProcess
	currentTime := time.Now()

	for _, processID := range dto.ProcessIDs {
		mappings = append(mappings, &models.ProjectProcess{
			ProcessID:   processID,
			ProjectID:   projectID,
			CreatedByID: userID,
			CreatedAt:   currentTime,
		})
	}

	if err := pais.repo.CreateProjectProcess(mappings, tx); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (pais *ProjectAssociatedItemService) DeleteProjectProcess(projectID uint, processID uint) error {
	// Validate control existence
	if _, err := pais.projectSvc.GetByID(projectID); err != nil {
		return utility.NewValidationError(fmt.Errorf("project does not exist for ID: %d", projectID), "not_exist", "asset_register_id")
	}

	tx := pais.repo.BeginTransaction()

	if err := pais.repo.DeleteProjectProcess(projectID, processID, tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (pais *ProjectAssociatedItemService) GetAllProjectStrategicObjectives(assetID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedItemStrategicObjectiveListResponseDTO, error) {
	if _, err := pais.projectSvc.GetByID(assetID); err != nil {
		return nil, err
	}
	return pais.repo.GetAllProjectStrategicObjectives(assetID, params)
}

func (pais *ProjectAssociatedItemService) CreateProjectStrategicObjective(projectID, userID uint, dto *commonDTOs.AssociatedStrategicObjectiveMapDTO) error {

	if _, err := pais.projectSvc.GetByID(projectID); err != nil {
		return utility.NewValidationError(fmt.Errorf("project does not exist for ID: %d", projectID), "not_exist", "asset_id")
	}

	if err := pais.validationSvc.ValidateStrtegicObjectiveMapping(projectID, dto.StrategicObjectiveIDs, "project_strategic_objectives", "project_id"); err != nil {
		return err
	}

	tx := pais.repo.BeginTransaction()
	var mappings []*models.ProjectStrategicObjective
	currentTime := time.Now()

	for _, StrategicObjectiveID := range dto.StrategicObjectiveIDs {
		mappings = append(mappings, &models.ProjectStrategicObjective{
			StrategicObjectiveID: StrategicObjectiveID,
			ProjectID:            projectID,
			CreatedBYID:          userID,
			CreatedAt:            currentTime,
		})
	}

	if err := pais.repo.CreateProjectStrategicObjective(mappings, tx); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (pais *ProjectAssociatedItemService) DeleteProjectStrategicObjective(projectID uint, strategicObjectiveId uint) error {
	// Validate control existence
	if _, err := pais.projectSvc.GetByID(projectID); err != nil {
		return utility.NewValidationError(fmt.Errorf("project does not exist for ID: %d", projectID), "not_exist", "asset_register_id")
	}

	tx := pais.repo.BeginTransaction()

	if err := pais.repo.DeleteProjectStrategicObjective(projectID, strategicObjectiveId, tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (pais *ProjectAssociatedItemService) GetAllProjectKPI(projectID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedKPIRegisterListResponseDTO, error) {
	if _, err := pais.projectSvc.GetByID(projectID); err != nil {
		return nil, err
	}
	return pais.repo.GetAllProjectKPI(projectID, params)
}

func (pais *ProjectAssociatedItemService) CreateProjectKPI(projectID, userID uint, dto *commonDTOs.AssociatedKPIMapDTO) error {

	if _, err := pais.projectSvc.GetByID(projectID); err != nil {
		return err
	}

	if err := pais.validationSvc.ValidateKPIMapping(projectID, dto.KpiIDs, "kpi_projects", "project_id"); err != nil {
		return err
	}

	tx := pais.repo.BeginTransaction()
	var mappings []*models.KPIProject
	currentTime := time.Now()

	for _, KpiID := range dto.KpiIDs {
		mappings = append(mappings, &models.KPIProject{
			KPIRegisterID: KpiID,
			ProjectID:     projectID,
			CreatedByID:   userID,
			CreatedAt:     currentTime,
		})
	}

	if err := pais.repo.CreateProjectKPI(mappings, tx); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (pais *ProjectAssociatedItemService) DeleteProjectKPI(projectID uint, kpiID uint) error {
	// Validate control existence
	if _, err := pais.projectSvc.GetByID(projectID); err != nil {
		return err
	}

	tx := pais.repo.BeginTransaction()

	if err := pais.repo.DeleteProjectKPI(projectID, kpiID, tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (pais *ProjectAssociatedItemService) GetAllProjectAssets(projectID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedItemAssetListResponseDTO, error) {
	if _, err := pais.projectSvc.GetByID(projectID); err != nil {
		return nil, err
	}
	return pais.repo.GetAllProjectAssets(projectID, params)
}

func (pais *ProjectAssociatedItemService) CreateProjectAsset(projectID, userID uint, dto *commonDTOs.AssociatedAssetMapDTO) error {

	if _, err := pais.projectSvc.GetByID(projectID); err != nil {
		return err
	}

	if err := pais.validationSvc.ValidateAssetMapping(projectID, dto.AssetIDs, "project_assets", "project_id"); err != nil {
		return err
	}

	tx := pais.repo.BeginTransaction()
	var mappings []*models.ProjectAsset
	currentTime := time.Now()

	for _, AssetID := range dto.AssetIDs {
		mappings = append(mappings, &models.ProjectAsset{
			AssetRegisterID: AssetID,
			ProjectID:       projectID,
			CreatedBYID:     userID,
			CreatedAt:       currentTime,
		})
	}

	if err := pais.repo.CreateProjectAsset(mappings, tx); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (pais *ProjectAssociatedItemService) DeleteProjectAsset(projectID uint, assetID uint) error {
	// Validate control existence
	if _, err := pais.projectSvc.GetByID(projectID); err != nil {
		return err
	}

	tx := pais.repo.BeginTransaction()

	if err := pais.repo.DeleteProjectAsset(projectID, assetID, tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (pais *ProjectAssociatedItemService) GetAllProjectIncidents(projectID uint, params *commonDTOs.AssociatedItemListQueryParamsDTO) (*commonDTOs.AssociatedIncidentListResponseDTO, error) {
	if _, err := pais.projectSvc.GetByID(projectID); err != nil {
		return nil, err
	}
	return pais.repo.GetAllProjectIncidents(projectID, params)
}

func (pais *ProjectAssociatedItemService) CreateProjectIncident(projectID, userID uint, dto *commonDTOs.AssociatedIncidentMapDTO) error {

	if _, err := pais.projectSvc.GetByID(projectID); err != nil {
		return err
	}

	if err := pais.validationSvc.ValidateIncidentMapping(projectID, dto.IncidentIDs, "project_incidents", "project_id"); err != nil {
		return err
	}

	tx := pais.repo.BeginTransaction()
	var mappings []*models.ProjectIncident
	currentTime := time.Now()

	for _, incidentID := range dto.IncidentIDs {
		mappings = append(mappings, &models.ProjectIncident{
			IncidentRegisterID: incidentID,
			ProjectID:          projectID,
			CreatedBYID:        userID,
			CreatedAt:          currentTime,
		})
	}

	if err := pais.repo.CreateProjectIncident(mappings, tx); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (pais *ProjectAssociatedItemService) DeleteProjectIncident(projectID uint, incidentID uint) error {
	// Validate control existence
	if _, err := pais.projectSvc.GetByID(projectID); err != nil {
		return err
	}

	tx := pais.repo.BeginTransaction()

	if err := pais.repo.DeleteProjectIncident(projectID, incidentID, tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
