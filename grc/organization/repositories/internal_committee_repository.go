package repositories

import (
	"fmt"
	commonDTOs "jethur/grc/common/dtos"
	"jethur/grc/models"
	organizationDTOs "jethur/grc/organization/dtos"
	"jethur/utility"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type InternalCommitteeRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(params *organizationDTOs.InternalCommitteeListQueryParamsDTO) (*organizationDTOs.InternalCommitteeListResponseDTO, error)
	GetByID(id uint) (*models.InternalCommittee, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.InternalCommittee, error)
	Create(internalCommittee *models.InternalCommittee, tx ...utility.DatabaseTransaction) (uint, error)
	Update(internalCommittee *models.InternalCommittee, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error
	GenerateReferenceNumber() (string, error)
	GetLatestVersionByID(id uint) (*models.InternalCommittee, error)

	CreateVersion(version *models.InternalCommitteeVersion, tx ...utility.DatabaseTransaction) (uint, error)
	UpdateVersion(version *models.InternalCommitteeVersion, tx ...utility.DatabaseTransaction) error
	DeleteVersion(id uint, tx ...utility.DatabaseTransaction) error
	GetVersionsByCommitteeID(committeeID uint, tx ...utility.DatabaseTransaction) ([]*models.InternalCommitteeVersion, error)
	UpdateVersionStatus(versionID uint, statusID uint, currentUserID uint, tx ...utility.DatabaseTransaction) error
	CreateNomineeDepartments(nomineeDepartments []*models.InternalCommitteeVersionNomineeDepartment, tx ...utility.DatabaseTransaction) error
	GetNomineeDepartmentsByVersionID(versionID uint, tx ...utility.DatabaseTransaction) ([]*models.InternalCommitteeVersionNomineeDepartment, error)
	DeleteNomineeDepartmentsByVersionID(versionID uint, tx ...utility.DatabaseTransaction) error
	GetNomineeDepartmentByID(id uint, tx ...utility.DatabaseTransaction) (*models.InternalCommitteeVersionNomineeDepartment, error)

	// Participant methods
	CreateParticipant(participant *models.InternalCommitteeVersionNomineeDepartmentParticipant, tx ...utility.DatabaseTransaction) (uint, error)
	GetParticipantsByNomineeDepartmentID(nomineeDepartmentID uint, tx ...utility.DatabaseTransaction) ([]*models.InternalCommitteeVersionNomineeDepartmentParticipant, error)
	DeleteParticipantsByNomineeDepartmentID(nomineeDepartmentID uint, tx ...utility.DatabaseTransaction) error
	DeleteParticipant(nomineeDepartmentID uint, participantID uint, tx ...utility.DatabaseTransaction) error
	HasParticipantsByVersionID(versionID uint, tx ...utility.DatabaseTransaction) (bool, error)
	ValidateAllDepartmentsHaveParticipants(versionID uint, tx ...utility.DatabaseTransaction) (bool, error)
}

type InternalCommitteeRepository struct {
	db *gorm.DB
}

func NewInternalCommitteeRepository(db *gorm.DB) *InternalCommitteeRepository {
	return &InternalCommitteeRepository{
		db: db,
	}
}

func (repo *InternalCommitteeRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *InternalCommitteeRepository) GetAll(params *organizationDTOs.InternalCommitteeListQueryParamsDTO) (*organizationDTOs.InternalCommitteeListResponseDTO, error) {
	var internalCommittees []*organizationDTOs.InternalCommitteeListDTO
	var totalCount int64

	query := repo.db.Table("internal_committees").
		Select(`
			internal_committees.id,
			internal_committees.reference_number,
			internal_committees.requesting_entity_id,
			requesting_entity.title_en AS requesting_entity_title_en,
			requesting_entity.title_ar AS requesting_entity_title_ar,
			latest_version.committee_name_en,
			latest_version.committee_name_ar,
			committee_types.title_en AS committee_type_title_en,
			committee_types.title_ar AS committee_type_title_ar,
			committee_terms.title_en AS committee_term_title_en,
			committee_terms.title_ar AS committee_term_title_ar,
			committee_classifications.title_en AS committee_classification_title_en,
			committee_classifications.title_ar AS committee_classification_title_ar,
			internal_committee_statuses.title_en AS committee_status_title_en,
			internal_committee_statuses.title_ar AS committee_status_title_ar,
			internal_committee_statuses.type AS committee_status_type,
			internal_committee_statuses.color_code AS committee_status_color_code,
			CAST(latest_version.version AS VARCHAR) AS version,
			CONVERT(VARCHAR, latest_version.term_end_date, 23) AS term_end_date,
			latest_version.justification,
			latest_version.roles_and_responsibilities,
			latest_version.committee_authorities,
			latest_version.is_disband_requested,
			CONVERT(VARCHAR, internal_committees.created_at, 120) AS created_at,
			CONVERT(VARCHAR, internal_committees.updated_at, 120) AS updated_at,
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
		Joins(`
			LEFT JOIN organization_units AS requesting_entity ON requesting_entity.id = internal_committees.requesting_entity_id
			LEFT JOIN internal_committee_versions AS latest_version ON latest_version.internal_committee_id = internal_committees.id AND latest_version.is_latest_version = 1
			LEFT JOIN committee_types ON committee_types.id = latest_version.committee_type_id
			LEFT JOIN committee_terms ON committee_terms.id = latest_version.committee_term_id
			LEFT JOIN committee_classifications ON committee_classifications.id = latest_version.committee_classification_id
			LEFT JOIN internal_committee_statuses ON internal_committee_statuses.id = latest_version.committee_status_id
			LEFT JOIN users AS created_user ON created_user.id = internal_committees.created_by_id
			LEFT JOIN users AS updated_user ON updated_user.id = internal_committees.updated_by_id
		`)

	// Apply dynamic filters
	for columnName, values := range params.Filters {
		if len(values) > 0 {
			switch columnName {
			case "committee_member_type_id":
				query = query.Where(`
					EXISTS (
						SELECT 1
						FROM internal_committee_version_nominee_departments AS nominee_departments
						WHERE nominee_departments.internal_committee_version_id = latest_version.id
						AND nominee_departments.committee_member_type_id IN (?)
					)
				`, values)
			default:
				query = query.Where(columnName+" IN (?)", values)
			}
		}
	}

	// Apply date filters
	if params.DateFrom != "" && params.DateTo != "" {
		query = query.Where("internal_committees.created_at BETWEEN ? AND ?", params.DateFrom+" 00:00:00", params.DateTo+" 23:59:59")
	}

	// Apply keyword search
	if params.Keyword != "" {
		query = query.Where(`(
			internal_committees.reference_number LIKE ? OR
			requesting_entity.title_en LIKE ? OR
			requesting_entity.title_ar LIKE ? OR
			latest_version.committee_name_en LIKE ? OR
			latest_version.committee_name_ar LIKE ?
		)`, "%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%", "%"+params.Keyword+"%")
	}

	// Apply status filter
	if params.Status == "active" || params.Status == "inactive" {
		query = query.Where("internal_committee_statuses.status = ?", params.Status)
	}

	// Apply quick filters
	for _, quickFilter := range params.QuickFilter {
		switch quickFilter {
		case "disbanded":
			query = query.Where("internal_committee_statuses.type = ?", "disbanded")
		case "approved":
			query = query.Where("internal_committee_statuses.type = ?", "approved")
		case "permanent":
			query = query.Where("committee_terms.type = ?", "permanent")
		}
	}

	// Get total count
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, err
	}
	// Apply sorting
	if params.SortBy != "" {
		query = query.Order(params.SortBy + " " + params.Sort)
	} else {
		query = query.Order("internal_committees.created_at desc")
	}

	// Apply pagination
	offset := (params.Page - 1) * params.PageSize
	query = query.Offset(offset).Limit(params.PageSize)

	// Execute query
	if err := query.Scan(&internalCommittees).Error; err != nil {
		return nil, err
	}

	// Calculate pagination info
	totalPages := 1
	if !params.Export {
		totalPages = int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))
	}

	return &organizationDTOs.InternalCommitteeListResponseDTO{
		Items: internalCommittees,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}, nil
}

func (repo *InternalCommitteeRepository) GetByID(id uint) (*models.InternalCommittee, error) {
	var internalCommittee *models.InternalCommittee
	if err := repo.db.First(&internalCommittee, id).Error; err != nil {
		return nil, err
	}
	return internalCommittee, nil
}

func (repo *InternalCommitteeRepository) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*models.InternalCommittee, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var internalCommittee *models.InternalCommittee
	if err := db.Preload("RequestingEntity").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		Preload("InternalCommitteeVersions", func(db *gorm.DB) *gorm.DB {
			return db.Where("is_latest_version = ?", true).
				Preload("CommitteeType").
				Preload("CommitteeTerm").
				Preload("CommitteeClassification").
				Preload("CommitteeStatus").
				Preload("CreatedByUser").
				Preload("UpdatedByUser").
				Preload("SubmittedByUser").
				Preload("ApprovedByUser").
				Preload("WorkflowStatus").
				Preload("InternalCommitteeVersionNomineeDepartments").
				Preload("InternalCommitteeVersionNomineeDepartments.CommitteeMemberType").
				Preload("InternalCommitteeVersionNomineeDepartments.NomineeDepartment").
				Preload("InternalCommitteeVersionNomineeDepartments.Participants").
				Preload("InternalCommitteeVersionNomineeDepartments.Participants.User").
				Preload("InternalCommitteeVersionNomineeDepartments.Participants.Department")
		}).
		First(&internalCommittee, id).
		Error; err != nil {
		return nil, err
	}

	return internalCommittee, nil
}

func (repo *InternalCommitteeRepository) Create(internalCommittee *models.InternalCommittee, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Create(internalCommittee).Error; err != nil {
		return 0, err
	}

	return internalCommittee.ID, nil
}

func (repo *InternalCommitteeRepository) Update(internalCommittee *models.InternalCommittee, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Save(internalCommittee).Error; err != nil {
		return err
	}

	return nil
}

func (repo *InternalCommitteeRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	// Delete internal committee
	if err := db.Delete(&models.InternalCommittee{}, id).Error; err != nil {
		return err
	}

	return nil
}

func (repo *InternalCommitteeRepository) CreateVersion(version *models.InternalCommitteeVersion, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Create(version).Error; err != nil {
		return 0, err
	}

	return version.ID, nil
}

func (repo *InternalCommitteeRepository) UpdateVersion(version *models.InternalCommitteeVersion, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Save(version).Error; err != nil {
		return err
	}

	return nil
}

func (repo *InternalCommitteeRepository) DeleteVersion(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Delete(&models.InternalCommitteeVersion{}, id).Error; err != nil {
		return err
	}

	return nil
}

func (repo *InternalCommitteeRepository) GetVersionsByCommitteeID(committeeID uint, tx ...utility.DatabaseTransaction) ([]*models.InternalCommitteeVersion, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var versions []*models.InternalCommitteeVersion
	if err := db.Where("internal_committee_id = ?", committeeID).Find(&versions).Error; err != nil {
		return nil, err
	}

	return versions, nil
}

func (repo *InternalCommitteeRepository) UpdateVersionStatus(versionID uint, statusID uint, currentUserID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	now := time.Now()
	updates := map[string]interface{}{
		"committee_status_id": statusID,
		"updated_by_id":       currentUserID,
		"updated_at":          now,
	}

	return db.Model(&models.InternalCommitteeVersion{}).
		Where("id = ?", versionID).
		Updates(updates).Error
}

func (repo *InternalCommitteeRepository) CreateNomineeDepartments(nomineeDepartments []*models.InternalCommitteeVersionNomineeDepartment, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	for _, dept := range nomineeDepartments {
		if err := db.Create(dept).Error; err != nil {
			return err
		}
	}

	return nil
}

func (repo *InternalCommitteeRepository) GetNomineeDepartmentsByVersionID(versionID uint, tx ...utility.DatabaseTransaction) ([]*models.InternalCommitteeVersionNomineeDepartment, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var nomineeDepartments []*models.InternalCommitteeVersionNomineeDepartment
	if err := db.Where("internal_committee_version_id = ?", versionID).
		Find(&nomineeDepartments).Error; err != nil {
		return nil, err
	}

	return nomineeDepartments, nil
}

func (repo *InternalCommitteeRepository) DeleteNomineeDepartmentsByVersionID(versionID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Where("internal_committee_version_id = ?", versionID).
		Delete(&models.InternalCommitteeVersionNomineeDepartment{}).Error
}

func (repo *InternalCommitteeRepository) GetNomineeDepartmentByID(id uint, tx ...utility.DatabaseTransaction) (*models.InternalCommitteeVersionNomineeDepartment, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var nomineeDepartment models.InternalCommitteeVersionNomineeDepartment
	if err := db.Where("id = ?", id).
		First(&nomineeDepartment).Error; err != nil {
		return nil, err
	}

	return &nomineeDepartment, nil
}

func (repo *InternalCommitteeRepository) GenerateReferenceNumber() (string, error) {
	var lastReferenceNumber string

	currentYear := time.Now().Year()
	prefix := fmt.Sprintf("COM-%d-", currentYear)

	// Get last reference number for current year
	err := repo.db.Model(&models.InternalCommittee{}).
		Select("reference_number").
		Where("reference_number LIKE ?", prefix+"%").
		Order("reference_number desc").
		Limit(1).
		Pluck("reference_number", &lastReferenceNumber).Error

	if err != nil {
		return "", err
	}

	lastDigits := "0000"

	if lastReferenceNumber != "" && len(lastReferenceNumber) >= 4 {
		lastDigits = lastReferenceNumber[len(lastReferenceNumber)-4:]
	}

	lastID, err := strconv.ParseUint(lastDigits, 10, 64)
	if err != nil {
		return "", err
	}

	newReferenceNumber := fmt.Sprintf("COM-%d-%04d", currentYear, lastID+1)

	return newReferenceNumber, nil
}

// Participant Methods

func (repo *InternalCommitteeRepository) CreateParticipant(participant *models.InternalCommitteeVersionNomineeDepartmentParticipant, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	if err := db.Create(participant).Error; err != nil {
		return 0, err
	}

	return participant.ID, nil
}

func (repo *InternalCommitteeRepository) GetParticipantsByNomineeDepartmentID(nomineeDepartmentID uint, tx ...utility.DatabaseTransaction) ([]*models.InternalCommitteeVersionNomineeDepartmentParticipant, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var participants []*models.InternalCommitteeVersionNomineeDepartmentParticipant
	err := db.Preload("NomineeDepartment").
		Preload("NomineeDepartment.CommitteeMemberType").
		Preload("NomineeDepartment.InternalCommitteeVersion").
		Where("internal_committee_version_nominee_department_id = ?", nomineeDepartmentID).
		Find(&participants).Error

	return participants, err
}

func (repo *InternalCommitteeRepository) DeleteParticipantsByNomineeDepartmentID(nomineeDepartmentID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Where("internal_committee_version_nominee_department_id = ?", nomineeDepartmentID).
		Delete(&models.InternalCommitteeVersionNomineeDepartmentParticipant{}).Error
}

func (repo *InternalCommitteeRepository) DeleteParticipant(nomineeDepartmentID uint, participantID uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	return db.Where("internal_committee_version_nominee_department_id = ? AND id = ?", nomineeDepartmentID, participantID).
		Delete(&models.InternalCommitteeVersionNomineeDepartmentParticipant{}).Error
}

func (repo *InternalCommitteeRepository) HasParticipantsByVersionID(versionID uint, tx ...utility.DatabaseTransaction) (bool, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	var count int64
	err := db.Table("internal_committee_version_nominee_department_participants").
		Joins("JOIN internal_committee_version_nominee_departments ON internal_committee_version_nominee_departments.id = internal_committee_version_nominee_department_participants.internal_committee_version_nominee_department_id").
		Where("internal_committee_version_nominee_departments.internal_committee_version_id = ?", versionID).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (repo *InternalCommitteeRepository) ValidateAllDepartmentsHaveParticipants(versionID uint, tx ...utility.DatabaseTransaction) (bool, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	// Get count of all nominee departments for this version
	var totalDepartments int64
	err := db.Table("internal_committee_version_nominee_departments").
		Where("internal_committee_version_id = ?", versionID).
		Count(&totalDepartments).Error

	if err != nil {
		return false, err
	}

	// Get count of nominee departments that have at least one participant
	var departmentsWithParticipants int64
	err = db.Table("internal_committee_version_nominee_departments").
		Joins("JOIN internal_committee_version_nominee_department_participants ON internal_committee_version_nominee_department_participants.internal_committee_version_nominee_department_id = internal_committee_version_nominee_departments.id").
		Where("internal_committee_version_nominee_departments.internal_committee_version_id = ?", versionID).
		Distinct("internal_committee_version_nominee_departments.id").
		Count(&departmentsWithParticipants).Error

	if err != nil {
		return false, err
	}

	// All departments have participants if counts match
	return totalDepartments == departmentsWithParticipants, nil
}

func (repo *InternalCommitteeRepository) GetLatestVersionByID(id uint) (*models.InternalCommittee, error) {
	var internalCommittee *models.InternalCommittee
	if err := repo.db.Preload("InternalCommitteeVersions", func(db *gorm.DB) *gorm.DB {
		return db.Where("is_latest_version = ?", true)
	}).First(&internalCommittee, id).Error; err != nil {
		return nil, err
	}
	return internalCommittee, nil
}
