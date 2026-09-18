package repositories

import (
	"jethur/grc/models"
	"jethur/grc/project/dtos"
	"jethur/utility"

	"gorm.io/gorm"
)

type ProjectClosureRepositoryInterface interface {
	GetAll(projectID uint) ([]*dtos.ProjectClosureListDTO, error)
	GetByID(projectID uint, id uint) (*models.ProjectClosure, error)
	GetFullDetailsByID(projectID uint, id uint, tx ...utility.DatabaseTransaction) (*models.ProjectClosure, error)
	Create(closure *models.ProjectClosure, tx ...utility.DatabaseTransaction) (uint, error)
	Update(closure *models.ProjectClosure, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error
	UpdateProjectPhase(projectID uint, phaseID uint, tx ...utility.DatabaseTransaction) error
	BeginTransaction() utility.DatabaseTransaction
}

type ProjectClosureRepository struct {
	db *gorm.DB
}

func NewProjectClosureRepository(db *gorm.DB) *ProjectClosureRepository {
	return &ProjectClosureRepository{db: db}
}

func (pcr *ProjectClosureRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(pcr.db)
}

func (pcr *ProjectClosureRepository) GetAll(projectID uint) ([]*dtos.ProjectClosureListDTO, error) {
	var closures []*dtos.ProjectClosureListDTO

	query := pcr.db.Table("project_closures").
		Select(`
			project_closures.id,
			project_closures.project_id,
			project_closures.overall_project_outcome,
			project_closures.closure_reason,
			project_closures.closure_comments,
			project_closures.lesson_learned,
			project_closures.approver_id,
			project_closures.is_approved,
			approver_user.first_name_en AS approver_user_first_name_en,
			approver_user.first_name_ar AS approver_user_first_name_ar,
			approver_user.last_name_en AS approver_user_last_name_en,
			approver_user.last_name_ar AS approver_user_last_name_ar,
			approver_user.initial AS approver_user_initial,
			approver_user.color AS approver_user_color,
			approver_user.email AS approver_user_email,
			approver_user.profile_image_path AS approver_user_profile_image_path,
			created_user.id AS created_user_id,
			created_user.first_name_en AS created_user_first_name_en,
			created_user.first_name_ar AS created_user_first_name_ar,
			created_user.last_name_en AS created_user_last_name_en,
			created_user.last_name_ar AS created_user_last_name_ar,
			created_user.initial AS created_user_initial,
			created_user.color AS created_user_color,
			created_user.email AS created_user_email,
			created_user.profile_image_path AS created_user_profile_image_path,
			CONVERT(VARCHAR, project_closures.created_at, 23) AS created_at,
			CONVERT(VARCHAR, project_closures.updated_at, 23) AS updated_at
		`).
		Joins("LEFT JOIN users AS approver_user ON approver_user.id = project_closures.approver_id").
		Joins("LEFT JOIN users AS created_user ON created_user.id = project_closures.created_by_id").
		Where("project_closures.project_id = ?", projectID).
		Order("project_closures.id DESC")

	if err := query.Find(&closures).Error; err != nil {
		return nil, err
	}

	return closures, nil
}

func (pcr *ProjectClosureRepository) GetByID(projectID uint, id uint) (*models.ProjectClosure, error) {
	var closure models.ProjectClosure
	result := pcr.db.Where("project_id = ? AND id = ?", projectID, id).First(&closure)
	if result.Error != nil {
		return nil, result.Error
	}
	return &closure, nil
}

func (pcr *ProjectClosureRepository) GetFullDetailsByID(projectID uint, id uint, tx ...utility.DatabaseTransaction) (*models.ProjectClosure, error) {
	db := pcr.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var closure models.ProjectClosure
	result := db.
		Preload("ApproverUser").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		Where("project_id = ? AND id = ?", projectID, id).
		First(&closure)

	if result.Error != nil {
		return nil, result.Error
	}
	return &closure, nil
}

func (pcr *ProjectClosureRepository) Create(closure *models.ProjectClosure, tx ...utility.DatabaseTransaction) (uint, error) {
	db := pcr.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(closure).Error; err != nil {
		return 0, err
	}
	return closure.ID, nil
}

func (pcr *ProjectClosureRepository) Update(closure *models.ProjectClosure, tx ...utility.DatabaseTransaction) error {
	db := pcr.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Save(closure).Error
}

func (pcr *ProjectClosureRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := pcr.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Delete(&models.ProjectClosure{}, id).Error
}

func (pcr *ProjectClosureRepository) UpdateProjectPhase(projectID uint, phaseID uint, tx ...utility.DatabaseTransaction) error {
	db := pcr.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Model(&models.Project{}).Where("id = ?", projectID).Update("project_phase_id", phaseID).Error
}
