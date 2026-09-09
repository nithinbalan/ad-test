package services

import (
	"errors"
	"time"

	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	"jethur/grc/project/dtos"
	"jethur/grc/project/repositories"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type ProjectClosureServiceInterface interface {
	GetAll(projectID uint) ([]*dtos.ProjectClosureListDTO, error)
	GetFullDetailsByID(projectID uint, id uint) (*dtos.ProjectClosureDetailDTO, error)
	Create(projectID uint, currentUserID uint, form *dtos.ProjectClosureFormDTO) (uint, error)
	Update(projectID uint, id uint, currentUserID uint, form *dtos.ProjectClosureFormDTO) error
	Delete(projectID uint, id uint, currentUserID uint) error
	Approve(projectID uint, id uint, currentUserID uint) error
}

type ProjectClosureService struct {
	repo            repositories.ProjectClosureRepositoryInterface
	projectSvc      ProjectServiceInterface
	projectPhaseSvc ProjectPhaseServiceInterface
	userSvc         organizationServices.UserServiceInterface
	activityLog     logServices.ActivityLogServiceInterface
}

func NewProjectClosureService(
	repo repositories.ProjectClosureRepositoryInterface,
	projectSvc ProjectServiceInterface,
	projectPhaseSvc ProjectPhaseServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	activityLog logServices.ActivityLogServiceInterface,
) ProjectClosureServiceInterface {
	return &ProjectClosureService{
		repo:            repo,
		projectSvc:      projectSvc,
		projectPhaseSvc: projectPhaseSvc,
		userSvc:         userSvc,
		activityLog:     activityLog,
	}
}

func (pcs *ProjectClosureService) GetAll(projectID uint) ([]*dtos.ProjectClosureListDTO, error) {
	if _, err := pcs.projectSvc.GetByID(projectID); err != nil {
		return nil, err
	}
	return pcs.repo.GetAll(projectID)
}

func (pcs *ProjectClosureService) GetFullDetailsByID(projectID uint, id uint) (*dtos.ProjectClosureDetailDTO, error) {
	if _, err := pcs.projectSvc.GetByID(projectID); err != nil {
		return nil, err
	}

	closure, err := pcs.repo.GetFullDetailsByID(projectID, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &utility.ResourceNotFoundError{}
		}
		return nil, err
	}

	var detailDTO dtos.ProjectClosureDetailDTO
	err = mapstructure.Decode(closure, &detailDTO)
	if err != nil {
		return nil, err
	}

	detailDTO.CreatedAt = closure.CreatedAt
	detailDTO.UpdatedAt = closure.UpdatedAt

	return &detailDTO, nil
}

func (pcs *ProjectClosureService) Create(projectID uint, currentUserID uint, form *dtos.ProjectClosureFormDTO) (uint, error) {
	if _, err := pcs.projectSvc.GetByID(projectID); err != nil {
		return 0, err
	}

	if _, err := pcs.userSvc.GetByID(form.ApproverID); err != nil {
		return 0, utility.NewValidationError(errors.New("approver user not found"), "user_not_found", "approver_id")
	}

	closure := &models.ProjectClosure{
		ProjectID:             projectID,
		OverallProjectOutcome: form.OverallProjectOutcome,
		ClosureReason:         form.ClosureReason,
		ClosureComments:       form.ClosureComments,
		LessonLearned:         form.LessonLearned,
		ApproverID:            form.ApproverID,
		CreatedByID:           currentUserID,
		CreatedAt:             time.Now(),
	}

	tx := pcs.repo.BeginTransaction()
	id, err := pcs.repo.Create(closure, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	go func() {
		_ = pcs.activityLog.LogActivity(
			models.CreateAction,
			models.ProjectClosure{}.TableName(),
			id,
			closure,
			currentUserID,
			"project", "project_register", "", nil,
		)
	}()

	return id, nil
}

func (pcs *ProjectClosureService) Update(projectID uint, id uint, currentUserID uint, form *dtos.ProjectClosureFormDTO) error {
	if _, err := pcs.projectSvc.GetByID(projectID); err != nil {
		return err
	}

	closure, err := pcs.repo.GetByID(projectID, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &utility.ResourceNotFoundError{}
		}
		return err
	}

	if _, err := pcs.userSvc.GetByID(form.ApproverID); err != nil {
		return utility.NewValidationError(errors.New("approver user not found"), "user_not_found", "approver_id")
	}

	oldData := map[string]interface{}{
		"overall_project_outcome": closure.OverallProjectOutcome,
		"closure_reason":          closure.ClosureReason,
		"closure_comments":        closure.ClosureComments,
		"lesson_learned":          closure.LessonLearned,
		"approver_id":             closure.ApproverID,
	}

	closure.OverallProjectOutcome = form.OverallProjectOutcome
	closure.ClosureReason = form.ClosureReason
	closure.ClosureComments = form.ClosureComments
	closure.LessonLearned = form.LessonLearned
	closure.ApproverID = form.ApproverID
	closure.UpdatedByID = &currentUserID
	updatedAt := time.Now()
	closure.UpdatedAt = &updatedAt

	tx := pcs.repo.BeginTransaction()
	if err := pcs.repo.Update(closure, tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	go func() {
		newData := map[string]interface{}{
			"overall_project_outcome": form.OverallProjectOutcome,
			"closure_reason":          form.ClosureReason,
			"closure_comments":        form.ClosureComments,
			"lesson_learned":          form.LessonLearned,
			"approver_id":             form.ApproverID,
		}
		changes := utility.CompareFieldChanges(oldData, newData)
		if len(changes) > 0 {
			_ = pcs.activityLog.LogActivity(
				models.UpdateAction,
				models.ProjectClosure{}.TableName(),
				id,
				changes,
				currentUserID,
				"project", "project_register", "", nil,
			)
		}
	}()

	return nil
}

func (pcs *ProjectClosureService) Delete(projectID uint, id uint, currentUserID uint) error {
	if _, err := pcs.projectSvc.GetByID(projectID); err != nil {
		return err
	}

	closure, err := pcs.repo.GetByID(projectID, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &utility.ResourceNotFoundError{}
		}
		return err
	}

	tx := pcs.repo.BeginTransaction()
	if err := pcs.repo.Delete(id, tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	go func() {
		_ = pcs.activityLog.LogActivity(
			models.DeleteAction,
			models.ProjectClosure{}.TableName(),
			id,
			closure,
			currentUserID,
			"project", "project_register", "", nil,
		)
	}()

	return nil
}

func (pcs *ProjectClosureService) Approve(projectID uint, id uint, currentUserID uint) error {
	if _, err := pcs.projectSvc.GetByID(projectID); err != nil {
		return err
	}

	closure, err := pcs.repo.GetByID(projectID, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &utility.ResourceNotFoundError{}
		}
		return err
	}

	if closure.ApproverID != currentUserID {
		return utility.NewForbiddenError("only the assigned approver can approve this project closure")
	}

	goLivePhase, err := pcs.projectPhaseSvc.GetByType("go_live")
	if err != nil {
		return err
	}

	tx := pcs.repo.BeginTransaction()

	if err := pcs.repo.UpdateProjectPhase(projectID, goLivePhase.ID, tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	go func() {
		_ = pcs.activityLog.LogActivity(
			models.UpdateAction,
			models.Project{}.TableName(),
			projectID,
			map[string]interface{}{"project_phase_id": goLivePhase.ID},
			currentUserID,
			"project", "project_register", "", nil,
		)
	}()

	return nil
}
