package repositories

import (
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/models"
	riskDTOs "jethur/grc/risk/dtos"
	"jethur/utility"
	"time"

	"gorm.io/gorm"
)

type RiskRegisterUpdatePlanRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(params *riskDTOs.RiskRegisterUpdatePlanListQueryParamsDTO) (*riskDTOs.RiskRegisterUpdatePlanListResponseDTO, error)
	GetByID(id uint) (*models.RiskRegisterUpdatePlan, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.RiskRegisterUpdatePlan, error)
	Create(plan *models.RiskRegisterUpdatePlan, tx ...utility.DatabaseTransaction) (uint, error)
	Update(plan *models.RiskRegisterUpdatePlan, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error

	// Team Members
	CreateTeamMembers(members []*models.RiskRegisterUpdatePlanTeamMember, tx ...utility.DatabaseTransaction) error
	DeleteTeamMembersByPlanID(planID uint, tx ...utility.DatabaseTransaction) error
	GetTeamMembersByPlanID(planID uint, tx ...utility.DatabaseTransaction) ([]*models.RiskRegisterUpdatePlanTeamMember, error)

	// Timeline Operations
	GetTimelineByPlanIDAndTypeID(planID uint, timelineTypeID uint, tx ...utility.DatabaseTransaction) (*models.RiskUpdatePlanTimeline, error)
	UpdateTimelineCompletion(timelineID uint, currentUserID uint, tx ...utility.DatabaseTransaction) error
	GetAllTimelineTypes() ([]*models.RiskUpdatePlanTimelineType, error)
	CreateTimeline(timeline *models.RiskUpdatePlanTimeline, tx ...utility.DatabaseTransaction) error
	CalculateAndUpdatePercentage(planID uint, tx ...utility.DatabaseTransaction) error
	DeleteTimelinesByPlanID(planID uint, tx ...utility.DatabaseTransaction) error
	MarkTimelineAsIncomplete(timelineID uint, tx ...utility.DatabaseTransaction) error
	GetByPlanID(id uint, tx ...utility.DatabaseTransaction) (*models.RiskRegisterUpdatePlan, error)

	// Plan Status Operations
	UpdatePlanStatus(planID uint, statusID uint, tx ...utility.DatabaseTransaction) error

	// Risk Operations
	GetRisksByPlanID(planID uint, tx ...utility.DatabaseTransaction) ([]*models.RiskRegister, error)
}

type RiskRegisterUpdatePlanRepository struct {
	db *gorm.DB
}

func NewRiskRegisterUpdatePlanRepository(db *gorm.DB) *RiskRegisterUpdatePlanRepository {
	return &RiskRegisterUpdatePlanRepository{
		db: db,
	}
}

func (repo *RiskRegisterUpdatePlanRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *RiskRegisterUpdatePlanRepository) GetAll(params *riskDTOs.RiskRegisterUpdatePlanListQueryParamsDTO) (*riskDTOs.RiskRegisterUpdatePlanListResponseDTO, error) {
	var plans []*riskDTOs.RiskRegisterUpdatePlanListDTO
	var totalCount int64

	query := repo.db.Table("risk_register_update_plans").
		Select(`
			risk_register_update_plans.id,
			risk_register_update_plans.division_id,
			div.title_en AS division_title_en,
			div.title_ar AS division_title_ar,
			risk_register_update_plans.business_unit_id,
			bu.title_en AS business_unit_title_en,
			bu.title_ar AS business_unit_title_ar,
			risk_register_update_plans.risk_register_update_plan_status_id,
			status.title_en AS status_title_en,
			status.title_ar AS status_title_ar,
			status.type AS status_type,
			status.color_code AS status_color_code,
			risk_register_update_plans.number_of_operating_entities,
			risk_register_update_plans.workdays,
			risk_register_update_plans.risk_champion_id,
			champion.first_name_en AS risk_champion_first_name_en,
			champion.first_name_ar AS risk_champion_first_name_ar,
			champion.last_name_en AS risk_champion_last_name_en,
			champion.last_name_ar AS risk_champion_last_name_ar,
			champion.email AS risk_champion_email,
			CONVERT(VARCHAR, risk_register_update_plans.start_date, 23) AS start_date,
			CONVERT(VARCHAR, risk_register_update_plans.end_date, 23) AS end_date,
			risk_register_update_plans.notes,
			risk_register_update_plans.total_percentage_completed,
			risk_register_update_plans.created_at,
			risk_register_update_plans.updated_at,
			created_by.id AS created_user_id,
			created_by.first_name_en AS created_user_first_name_en,
			created_by.first_name_ar AS created_user_first_name_ar,
			created_by.last_name_en AS created_user_last_name_en,
			created_by.last_name_ar AS created_user_last_name_ar,
			created_by.initial AS created_user_initial,
			created_by.color AS created_user_color,
			created_by.email AS created_user_email,
			created_by.profile_image_path AS created_user_profile_image_path,
			updated_by.id AS updated_user_id,
			updated_by.first_name_en AS updated_user_first_name_en,
			updated_by.first_name_ar AS updated_user_first_name_ar,
			updated_by.last_name_en AS updated_user_last_name_en,
			updated_by.last_name_ar AS updated_user_last_name_ar,
			updated_by.initial AS updated_user_initial,
			updated_by.color AS updated_user_color,
			updated_by.email AS updated_user_email,
			updated_by.profile_image_path AS updated_user_profile_image_path
		`).
		Joins("LEFT JOIN organization_units div ON risk_register_update_plans.division_id = div.id").
		Joins("LEFT JOIN organization_units bu ON risk_register_update_plans.business_unit_id = bu.id").
		Joins("LEFT JOIN risk_register_update_plan_statuses status ON risk_register_update_plans.risk_register_update_plan_status_id = status.id").
		Joins("LEFT JOIN users champion ON risk_register_update_plans.risk_champion_id = champion.id").
		Joins("LEFT JOIN users created_by ON risk_register_update_plans.created_by_id = created_by.id").
		Joins("LEFT JOIN users updated_by ON risk_register_update_plans.updated_by_id = updated_by.id")

	// Apply filters
	for column, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where("risk_register_update_plans."+column+" IN (?)", values)
		}
	}

	// Apply keyword search
	if params.Keyword != "" {
		searchPattern := "%" + params.Keyword + "%"
		query = query.Where(`
			div.title_en LIKE ? OR 
			div.title_ar LIKE ? OR 
			bu.title_en LIKE ? OR 
			bu.title_ar LIKE ? OR
			champion.first_name_en LIKE ? OR
			champion.last_name_en LIKE ? OR
			risk_register_update_plans.notes LIKE ?
		`, searchPattern, searchPattern, searchPattern, searchPattern, searchPattern, searchPattern, searchPattern)
	}

	// Count total records
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	// Apply sorting
	query = query.Order(params.SortBy + " " + params.Sort)

	// Execute query
	if err := query.Find(&plans).Error; err != nil {
		return nil, err
	}

	totalPages := 1

	totalPages = int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))

	return &riskDTOs.RiskRegisterUpdatePlanListResponseDTO{
		Items: plans,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalCount:  totalCount,
			TotalPages:  totalPages,
		},
	}, nil
}

func (repo *RiskRegisterUpdatePlanRepository) GetByID(id uint) (*models.RiskRegisterUpdatePlan, error) {
	db := repo.db
	var plan models.RiskRegisterUpdatePlan
	if err := db.First(&plan, id).Error; err != nil {
		return nil, err
	}

	return &plan, nil
}

func (repo *RiskRegisterUpdatePlanRepository) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.RiskRegisterUpdatePlan, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var plan models.RiskRegisterUpdatePlan
	if err := db.Preload("Division").
		Preload("BusinessUnit").
		Preload("RiskRegisterUpdatePlanStatus").
		Preload("RiskChampion").
		Preload("TeamMembers.User").
		Preload("Timelines").
		Preload("Timelines.RiskUpdatePlanTimelineType").
		Preload("Timelines.ErmDeptDocuments").
		Preload("Timelines.IntroductoryMeetingDocuments").
		Preload("Timelines.BusinessContexts").
		Preload("Timelines.Risks").
		Preload("Timelines.Risks.RiskRegister").
		Preload("Timelines.Controls").
		Preload("Timelines.Controls.Control").
		Preload("Timelines.Processes").
		Preload("Timelines.Processes.Process").
		Preload("Timelines.LeadershipAssessments").
		Preload("Timelines.LeadershipAssessments.User").
		Preload("Timelines.LeadershipDocuments").
		Preload("Timelines.RiskAssessmentWorkshopDocuments").
		Preload("Timelines.RiskTeamAssessments").
		Preload("Timelines.RiskTeamAssessments.RiskChampion").
		Preload("Timelines.RiskTeamDocuments").
		Preload("Timelines.RiskWorkshopAssessmentDocuments").
		Preload("Timelines.SelfAwarenessTrainings").
		Preload("Timelines.CreatedByUser").
		Preload("Timelines.UpdatedByUser").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		Preload("SubmittedByUser").
		Preload("ApprovedByUser").
		Preload("WorkflowStatus").
		First(&plan, id).Error; err != nil {
		return nil, err
	}

	return &plan, nil
}

func (repo *RiskRegisterUpdatePlanRepository) Create(plan *models.RiskRegisterUpdatePlan, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Create(plan).Error; err != nil {
		return 0, err
	}

	return plan.ID, nil
}

func (repo *RiskRegisterUpdatePlanRepository) Update(plan *models.RiskRegisterUpdatePlan, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Save(plan).Error; err != nil {
		return err
	}

	return nil
}

func (repo *RiskRegisterUpdatePlanRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Delete(&models.RiskRegisterUpdatePlan{}, id).Error; err != nil {
		return err
	}

	return nil
}

// Team Members
func (repo *RiskRegisterUpdatePlanRepository) CreateTeamMembers(members []*models.RiskRegisterUpdatePlanTeamMember, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if len(members) == 0 {
		return nil
	}

	if err := db.Create(&members).Error; err != nil {
		return err
	}

	return nil
}

func (repo *RiskRegisterUpdatePlanRepository) DeleteTeamMembersByPlanID(planID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("risk_register_update_plan_id = ?", planID).Delete(&models.RiskRegisterUpdatePlanTeamMember{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *RiskRegisterUpdatePlanRepository) GetTeamMembersByPlanID(planID uint, tx ...utility.DatabaseTransaction) ([]*models.RiskRegisterUpdatePlanTeamMember, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var members []*models.RiskRegisterUpdatePlanTeamMember
	if err := db.Preload("User").Where("risk_register_update_plan_id = ?", planID).Find(&members).Error; err != nil {
		return nil, err
	}

	return members, nil
}

// Timeline Operations

func (repo *RiskRegisterUpdatePlanRepository) GetAllTimelineTypes() ([]*models.RiskUpdatePlanTimelineType, error) {
	var timelineTypes []*models.RiskUpdatePlanTimelineType
	if err := repo.db.Find(&timelineTypes).Error; err != nil {
		return nil, err
	}
	return timelineTypes, nil
}

func (repo *RiskRegisterUpdatePlanRepository) CreateTimeline(timeline *models.RiskUpdatePlanTimeline, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Create(timeline).Error; err != nil {
		return err
	}

	return nil
}

func (repo *RiskRegisterUpdatePlanRepository) GetTimelineByPlanIDAndTypeID(planID uint, timelineTypeID uint, tx ...utility.DatabaseTransaction) (*models.RiskUpdatePlanTimeline, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var timeline models.RiskUpdatePlanTimeline
	if err := db.Where("risk_update_plan_id = ? AND risk_update_plan_timeline_type_id = ?", planID, timelineTypeID).
		First(&timeline).Error; err != nil {
		return nil, err
	}

	return &timeline, nil
}

func (repo *RiskRegisterUpdatePlanRepository) UpdateTimelineCompletion(timelineID uint, currentUserID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	now := time.Now()
	updates := map[string]interface{}{
		"is_completed":  true,
		"updated_by_id": currentUserID,
		"updated_at":    now,
	}

	return db.Model(&models.RiskUpdatePlanTimeline{}).
		Where("id = ?", timelineID).
		Updates(updates).Error
}

func (repo *RiskRegisterUpdatePlanRepository) MarkTimelineAsIncomplete(timelineID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	now := time.Now()
	if err := db.Model(&models.RiskUpdatePlanTimeline{}).
		Where("id = ?", timelineID).
		Updates(map[string]interface{}{
			"is_completed": false,
			"updated_at":   &now,
		}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *RiskRegisterUpdatePlanRepository) CalculateAndUpdatePercentage(planID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	// Count total timelines for this plan
	var totalCount int64
	if err := db.Model(&models.RiskUpdatePlanTimeline{}).
		Where("risk_update_plan_id = ?", planID).
		Count(&totalCount).Error; err != nil {
		return err
	}

	if totalCount == 0 {
		return nil
	}

	// Count completed timelines (is_completed = true)
	var completedCount int64
	if err := db.Model(&models.RiskUpdatePlanTimeline{}).
		Where("risk_update_plan_id = ? AND is_completed = ?", planID, true).
		Count(&completedCount).Error; err != nil {
		return err
	}

	// Get the plan to check if ApprovedBYID is set
	var plan models.RiskRegisterUpdatePlan
	if err := db.Select("approved_by_id").Where("id = ?", planID).First(&plan).Error; err != nil {
		return err
	}

	// Total steps = timeline count + 1 (final approval step)
	totalSteps := float64(totalCount + 1)

	// Completed steps = completed timelines + (1 if approved, else 0)
	completedSteps := float64(completedCount)
	if plan.ApprovedBYID != nil {
		completedSteps += 1
	}

	// Calculate percentage: (completedSteps / totalSteps) * 100
	percentage := int((completedSteps / totalSteps) * 100)

	// Update the plan's total_percentage_completed
	if err := db.Model(&models.RiskRegisterUpdatePlan{}).
		Where("id = ?", planID).
		Update("total_percentage_completed", percentage).Error; err != nil {
		return err
	}

	return nil
}

func (repo *RiskRegisterUpdatePlanRepository) DeleteTimelinesByPlanID(planID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Where("risk_update_plan_id = ?", planID).Delete(&models.RiskUpdatePlanTimeline{}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *RiskRegisterUpdatePlanRepository) GetByPlanID(id uint, tx ...utility.DatabaseTransaction) (*models.RiskRegisterUpdatePlan, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var plan models.RiskRegisterUpdatePlan
	if err := db.First(&plan, id).Error; err != nil {
		return nil, err
	}

	return &plan, nil
}

func (repo *RiskRegisterUpdatePlanRepository) UpdatePlanStatus(planID uint, statusID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Model(&models.RiskRegisterUpdatePlan{}).
		Where("id = ?", planID).
		Update("risk_register_update_plan_status_id", statusID).Error; err != nil {
		return err
	}

	return nil
}

func (repo *RiskRegisterUpdatePlanRepository) GetRisksByPlanID(planID uint, tx ...utility.DatabaseTransaction) ([]*models.RiskRegister, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var risks []*models.RiskRegister
	err := db.Where("risk_register_update_plan_id = ?", planID).Find(&risks).Error
	if err != nil {
		return nil, err
	}

	return risks, nil
}
