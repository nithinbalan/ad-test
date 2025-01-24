package repositories

import (
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type TestAndExerciseTestResultRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction
	GetByID(testAndExerciseID uint) (*models.TestAndExerciseTestResult, error)
	GetFullDetailsByID(testAndExerciseID uint, tx ...utility.DatabaseTransaction) (*models.TestAndExerciseTestResult, error)
	Create(TestAndExerciseTestResult *models.TestAndExerciseTestResult, tx ...utility.DatabaseTransaction) (uint, error)
	Update(TestAndExerciseTestResult *models.TestAndExerciseTestResult, tx ...utility.DatabaseTransaction) error
	Delete(id uint, tx ...utility.DatabaseTransaction) error

	UpdateTestAndExerciseProcess(process *models.TestAndExerciseProcess, tx ...utility.DatabaseTransaction) error
	CreateDocuments(docs []*models.TestAndExerciseTestResultDocument, tx ...utility.DatabaseTransaction) error
	GetDocumentByID(id uint) (*models.TestAndExerciseTestResultDocument, error)
	DeleteDocuments(docIds []*uint, tx ...utility.DatabaseTransaction) error
}

type TestAndExerciseTestResultRepository struct {
	db *gorm.DB
}

func NewTestAndExerciseTestResultRepository(db *gorm.DB) *TestAndExerciseTestResultRepository {
	return &TestAndExerciseTestResultRepository{
		db: db,
	}
}

// BeginTransaction starts a new transaction.
func (repo *TestAndExerciseTestResultRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}

func (repo *TestAndExerciseTestResultRepository) GetByID(testAndExerciseID uint) (*models.TestAndExerciseTestResult, error) {
	var testAndExerciseTestResult *models.TestAndExerciseTestResult

	result := repo.db.Where("test_and_exercise_id = ?", testAndExerciseID).First(&testAndExerciseTestResult)
	if result.Error != nil {
		return nil, result.Error
	}

	return testAndExerciseTestResult, nil

}

func (repo *TestAndExerciseTestResultRepository) GetFullDetailsByID(testAndExerciseID uint, tx ...utility.DatabaseTransaction) (*models.TestAndExerciseTestResult, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	var testAndExerciseTestResult *models.TestAndExerciseTestResult

	result := db.Where("test_and_exercise_id = ?", testAndExerciseID).
		Preload("TestAndExercise").
		Preload("TestAndExercise.BusinessContinuityPlan").
		Preload("TestAndExercise.BusinessContinuityPlan.BusinessFunction").
		Preload("TestAndExercise.TestAndExerciseStatus").
		Preload("TestAndExercise.TestAndExerciseProcesses").
		Preload("TestAndExercise.TestAndExerciseProcesses.Process").
		Preload("TestAndExercise.TestAndExerciseProcesses.PlannedRtoUnit").
		Preload("TestAndExercise.TestAndExerciseProcesses.PlannedMtpdUnit").
		Preload("TestAndExercise.TestAndExerciseProcesses.PlannedRpoUnit").
		Preload("TestAndExercise.TestAndExerciseProcesses.ActualRtoUnit").
		Preload("TestAndExercise.TestAndExerciseProcesses.ActualMtpdUnit").
		Preload("TestAndExercise.TestAndExerciseProcesses.ActualRpoUnit").
		Preload("TestAndExercise.TestAndExerciseProcesses.Timeliness").
		Preload("CommunicationFeedbackOption").
		Preload("RecoveryFeedbackOption").
		Preload("TestResultOption").
		Preload("Documents").
		Preload("CreatedByUser").
		Preload("UpdatedByUser").
		First(&testAndExerciseTestResult)

	if result.Error != nil {
		return nil, result.Error
	}

	return testAndExerciseTestResult, nil
}

func (repo *TestAndExerciseTestResultRepository) Create(TestAndExerciseTestResult *models.TestAndExerciseTestResult, tx ...utility.DatabaseTransaction) (uint, error) {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Create(TestAndExerciseTestResult).Error; err != nil {
		return 0, err
	}

	return TestAndExerciseTestResult.ID, nil
}

func (repo *TestAndExerciseTestResultRepository) Update(TestAndExerciseTestResult *models.TestAndExerciseTestResult, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(TestAndExerciseTestResult).Error; err != nil {
		return err
	}

	return nil
}

func (repo *TestAndExerciseTestResultRepository) Delete(id uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	// Delete the Business Communication Plan
	if err := db.Delete(&models.TestAndExerciseTestResult{}, id).Error; err != nil {
		return err
	}

	return nil
}

func (repo *TestAndExerciseTestResultRepository) UpdateTestAndExerciseProcess(process *models.TestAndExerciseProcess, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	if err := db.Save(process).Error; err != nil {
		return err
	}
	return nil
}

func (repo *TestAndExerciseTestResultRepository) CreateDocuments(docs []*models.TestAndExerciseTestResultDocument, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}
	// Save all documents in a single transaction
	if err := db.Create(&docs).Error; err != nil {
		return err
	}
	return nil
}

func (repo *TestAndExerciseTestResultRepository) GetDocumentByID(id uint) (*models.TestAndExerciseTestResultDocument, error) {
	var document models.TestAndExerciseTestResultDocument
	if err := repo.db.Where("id = ?", id).First(&document).Error; err != nil {
		return nil, err
	}
	return &document, nil
}

func (repo *TestAndExerciseTestResultRepository) DeleteDocuments(docIds []*uint, tx ...utility.DatabaseTransaction) error {
	db := repo.db
	if len(tx) > 0 && tx[0] != nil {
		db = tx[0].DB()
	}

	// Delete documents by IDs
	if err := db.Where("id IN (?)", docIds).Delete(&models.TestAndExerciseTestResultDocument{}).Error; err != nil {
		return err
	}

	return nil
}
