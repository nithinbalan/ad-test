package repositories

import (
	"errors"

	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type ProjectSecurityGateReportRepositoryInterface interface {
	GetFullDetailsByID(id uint) (*models.ProjectSecurityGate, error)
	GetAllChecklists(securityGateID uint) ([]*models.ProjectSecurityGateChecklist, error)
	GetAllExceptionWorkflowHistories(securityGateID uint) ([]*models.ProjectSecurityGateChecklistExceptionWorkflowHistory, error)

	GetAllExceptionChecklists() ([]*models.ProjectSecurityGateChecklist, error)
	GetAllExceptionWorkflowHistoriesForAllGates() ([]*models.ProjectSecurityGateChecklistExceptionWorkflowHistory, error)

	GetAllResponseDocuments(securityGateID uint) ([]*models.ProjectSecurityGateChecklistResponseDocument, error)
	GetAllWorkflowCommentDocuments(securityGateID uint) ([]*models.ProjectSecurityGateChecklistCommentDocument, error)
}

type ProjectSecurityGateReportRepository struct {
	db *gorm.DB
}

func NewProjectSecurityGateReportRepository(db *gorm.DB) ProjectSecurityGateReportRepositoryInterface {
	return &ProjectSecurityGateReportRepository{
		db: db,
	}
}

func (repo *ProjectSecurityGateReportRepository) GetFullDetailsByID(id uint) (*models.ProjectSecurityGate, error) {
	var projectSecurityGate *models.ProjectSecurityGate
	result := repo.db.
		Preload("Stage").
		Preload("TemplateVersion").
		Preload("BusinessOwner").
		Preload("ReadinessStatus").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		Preload("Classifications").
		Preload("Classifications.Classification").
		First(&projectSecurityGate, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, utility.NewResourceNotFoundError(result.Error, id)
	} else if result.Error != nil {
		return nil, result.Error
	}

	return projectSecurityGate, nil
}

func (repo *ProjectSecurityGateReportRepository) GetAllChecklists(securityGateID uint) ([]*models.ProjectSecurityGateChecklist, error) {
	var checklists []*models.ProjectSecurityGateChecklist
	result := repo.db.
		Select("project_security_gate_checklists.*").
		Preload("Requirement").
		Preload("ResponseStatus").
		Preload("WorkflowStatus").
		Preload("SubmittedByUser").
		Preload("ApprovedByUser").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		Preload("Documents").
		Preload("Documents.CreatedByUser").
		Preload("Controls").
		Preload("Controls.Control").
		Preload("Risks").
		Preload("Risks.RiskRegister").
		Preload("Risks.RiskRegister.RiskStatus").
		Preload("Risks.RiskRegister.RiskCategory").
		Preload("Risks.RiskRegister.ResponsibleUnit").
		Preload("Risks.RiskRegister.ResponsibleUser").
		Joins("JOIN project_security_gate_checklist_template_requirements AS requirements ON requirements.id = project_security_gate_checklists.requirement_id").
		Where("project_security_gate_checklists.security_gate_id = ?", securityGateID).
		Order("requirements.is_critical DESC, project_security_gate_checklists.id ASC").
		Find(&checklists)

	if result.Error != nil {
		return nil, result.Error
	}

	return checklists, nil
}

func (repo *ProjectSecurityGateReportRepository) GetAllExceptionWorkflowHistories(securityGateID uint) ([]*models.ProjectSecurityGateChecklistExceptionWorkflowHistory, error) {
	var histories []*models.ProjectSecurityGateChecklistExceptionWorkflowHistory
	result := repo.db.
		Select("project_security_gate_checklist_exception_workflow_histories.*").
		Preload("WorkflowStatus").
		Preload("CreatedByUser").
		Joins("JOIN project_security_gate_checklists AS checklists ON checklists.id = project_security_gate_checklist_exception_workflow_histories.checklist_id").
		Where("checklists.security_gate_id = ?", securityGateID).
		Order("project_security_gate_checklist_exception_workflow_histories.checklist_id ASC, project_security_gate_checklist_exception_workflow_histories.created_at ASC").
		Find(&histories)

	if result.Error != nil {
		return nil, result.Error
	}

	return histories, nil
}

func (repo *ProjectSecurityGateReportRepository) GetAllExceptionChecklists() ([]*models.ProjectSecurityGateChecklist, error) {
	var checklists []*models.ProjectSecurityGateChecklist
	result := repo.db.
		Select("project_security_gate_checklists.*").
		Preload("Requirement").
		Preload("ResponseStatus").
		Preload("WorkflowStatus").
		Preload("SubmittedByUser").
		Preload("ApprovedByUser").
		Preload("SecurityGate").
		Preload("SecurityGate.Stage").
		Preload("SecurityGate.ReadinessStatus").
		Preload("SecurityGate.BusinessOwner").
		Joins("JOIN project_security_gates AS security_gates ON security_gates.id = project_security_gate_checklists.security_gate_id").
		Joins("JOIN project_security_gate_checklist_template_requirements AS requirements ON requirements.id = project_security_gate_checklists.requirement_id").
		Where("project_security_gate_checklists.is_exception = ?", true).
		Order("security_gates.project_title_en ASC, requirements.is_critical DESC, project_security_gate_checklists.id ASC").
		Find(&checklists)

	if result.Error != nil {
		return nil, result.Error
	}

	return checklists, nil
}

func (repo *ProjectSecurityGateReportRepository) GetAllExceptionWorkflowHistoriesForAllGates() ([]*models.ProjectSecurityGateChecklistExceptionWorkflowHistory, error) {
	var histories []*models.ProjectSecurityGateChecklistExceptionWorkflowHistory
	result := repo.db.
		Select("project_security_gate_checklist_exception_workflow_histories.*").
		Preload("WorkflowStatus").
		Preload("CreatedByUser").
		Joins("JOIN project_security_gate_checklists AS checklists ON checklists.id = project_security_gate_checklist_exception_workflow_histories.checklist_id").
		Where("checklists.is_exception = ?", true).
		Order("project_security_gate_checklist_exception_workflow_histories.checklist_id ASC, project_security_gate_checklist_exception_workflow_histories.created_at ASC").
		Find(&histories)

	if result.Error != nil {
		return nil, result.Error
	}

	return histories, nil
}

func (repo *ProjectSecurityGateReportRepository) GetAllResponseDocuments(securityGateID uint) ([]*models.ProjectSecurityGateChecklistResponseDocument, error) {
	var documents []*models.ProjectSecurityGateChecklistResponseDocument
	result := repo.db.
		Select("project_security_gate_checklist_response_documents.*").
		Preload("Checklist").
		Preload("Checklist.Requirement").
		Joins("JOIN project_security_gate_checklists AS checklists ON checklists.id = project_security_gate_checklist_response_documents.project_security_gate_checklist_id").
		Where("checklists.security_gate_id = ?", securityGateID).
		Order("project_security_gate_checklist_response_documents.project_security_gate_checklist_id ASC, project_security_gate_checklist_response_documents.id ASC").
		Find(&documents)

	if result.Error != nil {
		return nil, result.Error
	}

	return documents, nil
}

// GetAllWorkflowCommentDocuments skips plain discussion comments
func (repo *ProjectSecurityGateReportRepository) GetAllWorkflowCommentDocuments(securityGateID uint) ([]*models.ProjectSecurityGateChecklistCommentDocument, error) {
	var documents []*models.ProjectSecurityGateChecklistCommentDocument
	result := repo.db.
		Select("project_security_gate_checklist_comment_documents.*").
		Preload("Comment").
		Preload("Comment.Checklist").
		Preload("Comment.Checklist.Requirement").
		Joins("JOIN project_security_gate_checklist_comments AS comments ON comments.id = project_security_gate_checklist_comment_documents.project_security_gate_checklist_comment_id").
		Joins("JOIN project_security_gate_checklists AS checklists ON checklists.id = comments.checklist_id").
		Where("checklists.security_gate_id = ? AND comments.workflow_history_id IS NOT NULL", securityGateID).
		Order("comments.checklist_id ASC, comments.id ASC, project_security_gate_checklist_comment_documents.id ASC").
		Find(&documents)

	if result.Error != nil {
		return nil, result.Error
	}

	return documents, nil
}
