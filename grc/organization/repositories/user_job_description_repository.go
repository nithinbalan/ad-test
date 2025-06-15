package repositories

import (
	"jethur/grc/models"
	"jethur/grc/organization/dtos"
	"jethur/utility"

	"gorm.io/gorm"
)

type UserJobDescriptionRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetAll(sort string, sortBy string, keyword string, userID uint) ([]*dtos.UserJobDescriptionListDTO, error)
	GetFullDetailsByID(id uint) (*models.UserJobDescription, error)
	CreateUserJobDescriptions(tx utility.DatabaseTransaction, userJDs []*models.UserJobDescription) error
	CreateUserJobDescriptionCompetency(tx utility.DatabaseTransaction, ujdc *models.UserJobDescriptionCompetency) error
	CreateUserJobDescriptionCertificates(tx utility.DatabaseTransaction, userJDCerts []*models.UserJobDescriptionCertificate) error

	GetByUserAndJD(userID, jdID uint) (*models.UserJobDescription, error)
	DeleteUserJobDescriptionCertificatesByUserJDID(tx utility.DatabaseTransaction, userJDID uint) error
	DeleteUserJobDescriptionCompetenciesByUserJDID(tx utility.DatabaseTransaction, userJDID uint) error
	DeleteUserJobDescription(tx utility.DatabaseTransaction, userJDID uint) error

	GetCompetencyByID(id uint) (*models.UserJobDescriptionCompetency, error)
	UpdateUserJDCompetency(comp *models.UserJobDescriptionCompetency, tx utility.DatabaseTransaction) error
	GetCertificateByID(id uint) (*models.UserJobDescriptionCertificate, error)
	DeleteUserJobDescriptionCertificates(ids []*uint, tx utility.DatabaseTransaction) error
	UpdateUserJobDescriptionCertificates(tx utility.DatabaseTransaction, certs []*models.UserJobDescriptionCertificate) error
	GetCertificatesByUserJDID(userJDID uint) ([]*models.UserJobDescriptionCertificate, error)
}

type UserJobDescriptionRepository struct {
	db *gorm.DB
}

func NewUserJobDescriptionRepository(db *gorm.DB) *UserJobDescriptionRepository {
	return &UserJobDescriptionRepository{
		db: db,
	}
}

func (repo *UserJobDescriptionRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *UserJobDescriptionRepository) GetAll(sort string, sortBy string, keyword string, userID uint) ([]*dtos.UserJobDescriptionListDTO, error) {
	var jobDescs []*dtos.UserJobDescriptionListDTO

	query := repo.db.Table("user_job_descriptions").
		Select(`
			user_job_descriptions.*,
			job_descriptions.title_en AS title_en,
			job_descriptions.title_ar AS title_ar,
			job_descriptions.description_en AS description_en,
			job_descriptions.description_ar AS description_ar,
			job_descriptions.roles_and_responsibilities AS roles_and_responsibilities,
			job_descriptions.status AS status,
			user_job_descriptions.created_at AS created_at
		`).
		Joins(`
			LEFT JOIN job_descriptions ON job_descriptions.id = user_job_descriptions.job_description_id
		`).
		Where("user_job_descriptions.user_id = ?", userID)

	if keyword != "" {
		keywordFilter := "%" + keyword + "%"
		query = query.Where(`
			(job_descriptions.title_en LIKE ? OR
			job_descriptions.title_ar LIKE ?)`,
			keywordFilter, keywordFilter)
	}

	if sortBy == "" {
		sortBy = "user_job_descriptions.id"
	}
	if sort != "asc" && sort != "desc" {
		sort = "asc"
	}
	query = query.Order(sortBy + " " + sort)

	err := query.Find(&jobDescs).Error
	if err != nil {
		return nil, err
	}

	return jobDescs, nil
}

func (repo *UserJobDescriptionRepository) GetFullDetailsByID(id uint) (*models.UserJobDescription, error) {
	var model models.UserJobDescription
	err := repo.db.
		Preload("JobDescription").
		Preload("User").
		Preload("Competencies.Competency").
		Preload("Competencies.ActualLevel").
		Preload("Competencies.ExpectedLevel").
		Preload("Competencies.Certificates").
		Preload("CreatedByUser").
		First(&model, id).Error
	if err != nil {
		return nil, err
	}
	return &model, nil
}

func (repo *UserJobDescriptionRepository) CreateUserJobDescriptions(tx utility.DatabaseTransaction, userJDs []*models.UserJobDescription) error {
	db := repo.db
	if tx != nil {
		db = tx.DB()
	}
	return db.Create(userJDs).Error
}
func (repo *UserJobDescriptionRepository) CreateUserJobDescriptionCompetency(tx utility.DatabaseTransaction, ujdc *models.UserJobDescriptionCompetency) error {
	db := repo.db
	if tx != nil {
		db = tx.DB()
	}
	return db.Create(ujdc).Error
}

func (repo *UserJobDescriptionRepository) CreateUserJobDescriptionCertificates(tx utility.DatabaseTransaction, userJDCerts []*models.UserJobDescriptionCertificate) error {
	db := repo.db
	if tx != nil {
		db = tx.DB()
	}
	return db.Create(userJDCerts).Error
}

func (repo *UserJobDescriptionRepository) GetByUserAndJD(userID, jdID uint) (*models.UserJobDescription, error) {
	var userJD models.UserJobDescription
	err := repo.db.Where("user_id = ? AND job_description_id = ?", userID, jdID).First(&userJD).Error
	if err != nil {
		return nil, err
	}
	return &userJD, nil
}

func (repo *UserJobDescriptionRepository) DeleteUserJobDescriptionCertificatesByUserJDID(tx utility.DatabaseTransaction, userJDID uint) error {
	db := repo.db
	if tx != nil {
		db = tx.DB()
	}
	return db.Where("user_job_description_competency_id IN (?)",
		db.Model(&models.UserJobDescriptionCompetency{}).Select("id").Where("user_job_description_id = ?", userJDID),
	).Delete(&models.UserJobDescriptionCertificate{}).Error
}

func (repo *UserJobDescriptionRepository) DeleteUserJobDescriptionCompetenciesByUserJDID(tx utility.DatabaseTransaction, userJDID uint) error {
	db := repo.db
	if tx != nil {
		db = tx.DB()
	}
	return db.Where("user_job_description_id = ?", userJDID).
		Delete(&models.UserJobDescriptionCompetency{}).Error
}

func (repo *UserJobDescriptionRepository) DeleteUserJobDescription(tx utility.DatabaseTransaction, userJDID uint) error {
	db := repo.db
	if tx != nil {
		db = tx.DB()
	}
	return db.Delete(&models.UserJobDescription{}, userJDID).Error
}

func (r *UserJobDescriptionRepository) GetCompetencyByID(id uint) (*models.UserJobDescriptionCompetency, error) {
	var comp models.UserJobDescriptionCompetency
	if err := r.db.First(&comp, id).Error; err != nil {
		return nil, err
	}
	return &comp, nil
}

func (r *UserJobDescriptionRepository) UpdateUserJDCompetency(comp *models.UserJobDescriptionCompetency, tx utility.DatabaseTransaction) error {
	db := r.db
	if tx != nil {
		db = tx.DB()
	}
	return db.Save(comp).Error
}

func (r *UserJobDescriptionRepository) GetCertificateByID(id uint) (*models.UserJobDescriptionCertificate, error) {
	var cert models.UserJobDescriptionCertificate
	if err := r.db.First(&cert, id).Error; err != nil {
		return nil, err
	}
	return &cert, nil
}

func (r *UserJobDescriptionRepository) DeleteUserJobDescriptionCertificates(ids []*uint, tx utility.DatabaseTransaction) error {
	db := r.db
	if tx != nil {
		db = tx.DB()
	}
	return db.Where("id IN ?", ids).Delete(&models.UserJobDescriptionCertificate{}).Error
}

func (repo *UserJobDescriptionRepository) UpdateUserJobDescriptionCertificates(tx utility.DatabaseTransaction, certs []*models.UserJobDescriptionCertificate) error {
	db := repo.db
	if tx != nil {
		db = tx.DB()
	}
	for _, cert := range certs {
		if err := db.Save(cert).Error; err != nil {
			return err
		}
	}
	return nil
}

func (repo *UserJobDescriptionRepository) GetCertificatesByUserJDID(userJDID uint) ([]*models.UserJobDescriptionCertificate, error) {
	var certs []*models.UserJobDescriptionCertificate
	err := repo.db.
		Where("user_job_description_competency_id IN (?)",
			repo.db.Table("user_job_description_competencies").Select("id").Where("user_job_description_id = ?", userJDID),
		).Find(&certs).Error
	return certs, err
}
