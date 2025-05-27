package repositories

import (
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/models"
	"jethur/grc/organization/dtos"
	"jethur/utility"

	"gorm.io/gorm"
)

type ProcessActivityRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(processID uint, params *dtos.ProcessActivityListQueryParamsDTO) (*dtos.ProcessActivityListResponseDTO, error)
	GetByID(id uint) (*models.ProcessActivity, error)
	GetFullDetailsByID(id uint) (*models.ProcessActivity, error)
	Create(activity *models.ProcessActivity, tx ...utility.DatabaseTransaction) (uint, error)
	Update(activity *models.ProcessActivity, tx ...utility.DatabaseTransaction) error
	Delete(id uint) error

	CreateProcessActivityUsers(users []*models.ProcessActivityUser, tx ...utility.DatabaseTransaction) error
	DeleteProcessActivityUsers(userIDs []uint, tx ...utility.DatabaseTransaction) error
	DeleteProcessActivityUsersByActivityID(activityID uint) error
	GetActivityUsers(activityID uint) ([]*models.ProcessActivityUser, error)
	GetActivityUserByID(id uint) (*models.ProcessActivityUser, error)
	DeleteProcessActivityUser(id uint) error
}

type ProcessActivityRepository struct {
	db *gorm.DB
}

func NewProcessActivityRepository(db *gorm.DB) *ProcessActivityRepository {
	return &ProcessActivityRepository{db: db}
}

func (repo *ProcessActivityRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *ProcessActivityRepository) GetAll(processID uint, params *dtos.ProcessActivityListQueryParamsDTO) (*dtos.ProcessActivityListResponseDTO, error) {
	var activities []*dtos.ProcessActivityListDTO
	var totalCount int64

	query := repo.db.Table("process_activities").
		Select(`
			process_activities.*,
			created_user.id AS created_user_id,
			created_user.first_name_en AS created_user_first_name_en,
			created_user.first_name_ar AS created_user_first_name_ar,
			created_user.last_name_en AS created_user_last_name_en,
			created_user.last_name_ar AS created_user_last_name_ar,
			created_user.initial AS created_user_initial,
			created_user.color AS created_user_color,
			created_user.email AS created_user_email,
			created_user.profile_image_path AS created_user_profile_image_path,
			updated_user.id AS updated_user_id,
			updated_user.first_name_en AS updated_user_first_name_en,
			updated_user.first_name_ar AS updated_user_first_name_ar,
			updated_user.last_name_en AS updated_user_last_name_en,
			updated_user.last_name_ar AS updated_user_last_name_ar,
			updated_user.initial AS updated_user_initial,
			updated_user.color AS updated_user_color,
			updated_user.email AS updated_user_email,
			updated_user.profile_image_path AS updated_user_profile_image_path
		`).
		Joins("LEFT JOIN users AS created_user ON created_user.id = process_activities.created_by_id").
		Joins("LEFT JOIN users AS updated_user ON updated_user.id = process_activities.updated_by_id").
		Where("process_activities.process_id = ?", processID)

	// Apply dynamic filters to the query
	for columnName, values := range params.Filters {
		if len(values) > 0 {
			query = query.Where("process_activities."+columnName+" IN (?)", values)
		}
	}

	// Apply keyword search
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
			(process_activities.title_en LIKE ? OR
			process_activities.title_ar LIKE ? OR
			process_activities.expected_output LIKE ? OR
			process_activities.estimated_duration LIKE ? OR
			process_activities.description_en LIKE ? OR
			process_activities.description_ar LIKE ?)`,
			keyword, keyword, keyword, keyword, keyword, keyword,
		)
	}

	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}

	query = query.Order(params.SortBy + " " + params.Sort)

	offset := (params.Page - 1) * params.PageSize
	query = query.Offset(offset).Limit(params.PageSize)

	if err := query.Find(&activities).Error; err != nil {
		return nil, err
	}

	totalPages := int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))
	response := &dtos.ProcessActivityListResponseDTO{
		Items: activities,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}
	return response, nil
}

func (repo *ProcessActivityRepository) GetByID(id uint) (*models.ProcessActivity, error) {
	var activity models.ProcessActivity
	if err := repo.db.First(&activity, id).Error; err != nil {
		return nil, err
	}
	return &activity, nil
}

func (repo *ProcessActivityRepository) GetFullDetailsByID(id uint) (*models.ProcessActivity, error) {
	var activity models.ProcessActivity
	if err := repo.db.Preload("Users").
		Preload("Users.User").
		Preload("Users.ResponsibleUnit").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		First(&activity, id).Error; err != nil {
		return nil, err
	}
	return &activity, nil
}

func (repo *ProcessActivityRepository) Create(activity *models.ProcessActivity, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(activity).Error; err != nil {
		return 0, err
	}
	return activity.ID, nil
}

func (repo *ProcessActivityRepository) Update(activity *models.ProcessActivity, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Save(activity).Error
}

func (repo *ProcessActivityRepository) Delete(id uint) error {
	return repo.db.Delete(&models.ProcessActivity{}, id).Error
}

func (repo *ProcessActivityRepository) CreateProcessActivityUsers(users []*models.ProcessActivityUser, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Create(&users).Error
}

func (repo *ProcessActivityRepository) DeleteProcessActivityUsers(userIDs []uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	return db.Where("id IN ?", userIDs).Delete(&models.ProcessActivityUser{}).Error
}

func (repo *ProcessActivityRepository) DeleteProcessActivityUsersByActivityID(activityID uint) error {
	return repo.db.Where("process_activity_id = ?", activityID).Delete(&models.ProcessActivityUser{}).Error
}

func (repo *ProcessActivityRepository) GetActivityUsers(activityID uint) ([]*models.ProcessActivityUser, error) {
	var mappings []*models.ProcessActivityUser
	if err := repo.db.Where("process_activity_id = ?", activityID).Find(&mappings).Error; err != nil {
		return nil, err
	}
	return mappings, nil
}

func (repo *ProcessActivityRepository) GetActivityUserByID(id uint) (*models.ProcessActivityUser, error) {
	var user models.ProcessActivityUser
	if err := repo.db.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (repo *ProcessActivityRepository) DeleteProcessActivityUser(id uint) error {
	return repo.db.Delete(&models.ProcessActivityUser{}, id).Error
}
