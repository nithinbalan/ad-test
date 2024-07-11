package services

import (
	"errors"
	"fmt"
	incidentDTOs "jethur/grc/incident/dtos"
	incidentRepositories "jethur/grc/incident/repositories"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type IncidentClosureServiceInterface interface {
	GetByID(incidentID uint) (*models.IncidentClosure, error)
	GetFullDetailsByID(incidentID uint, tx ...utility.DatabaseTransaction) (*incidentDTOs.IncidentClosureDetailDTO, error)
	Create(incidentID uint, currentUserID uint, IncidentClosureFormDTO *incidentDTOs.IncidentClosureFormDTO) (uint, error)
	Update(incidentID uint, currentUserID uint, incident *incidentDTOs.IncidentClosureFormDTO) error
	Delete(incidentID uint, currentUserID uint) (uint, error)
}

// Service and implementations
type IncidentClosureService struct {
	repo                incidentRepositories.IncidentClosureRepositoryInterface
	activityLogSvc      logServices.ActivityLogServiceInterface
	incidentRegisterSvc IncidentRegisterServiceInterface
}

func NewIncidentClosureService(
	repo incidentRepositories.IncidentClosureRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	incidentRegisterSvc IncidentRegisterServiceInterface) IncidentClosureServiceInterface {
	return &IncidentClosureService{
		repo:                repo,
		activityLogSvc:      activityLogSvc,
		incidentRegisterSvc: incidentRegisterSvc,
	}
}

func (ics IncidentClosureService) GetByID(incidentID uint) (*models.IncidentClosure, error) {
	if _, err := ics.incidentRegisterSvc.GetByID(incidentID); err != nil {
		return nil, err
	}

	incidentClosure, err := ics.repo.GetByID(incidentID)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, incidentID)
		}
		return nil, err
	}
	return incidentClosure, nil
}

func (ics IncidentClosureService) GetFullDetailsByID(incidentID uint, tx ...utility.DatabaseTransaction) (*incidentDTOs.IncidentClosureDetailDTO, error) {
	var incidentClosureDTO incidentDTOs.IncidentClosureDetailDTO

	if _, err := ics.incidentRegisterSvc.GetByID(incidentID); err != nil {
		return nil, err
	}

	incidentClosure, err := ics.repo.GetFullDetailsByID(incidentID, tx...)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, incidentID)
		}
		return nil, err
	}

	err = mapstructure.Decode(incidentClosure, &incidentClosureDTO)
	if err != nil {
		return nil, err
	}
	if incidentClosure.CompletionDate != nil {
		completionDate := incidentClosure.CompletionDate.Format("2006-01-02")
		incidentClosureDTO.CompletionDateString = &completionDate
	}

	incidentClosureDTO.CreatedAt = incidentClosure.CreatedAt
	incidentClosureDTO.UpdatedAt = incidentClosure.UpdatedAt

	//todo :- also add extra fields which required in details screen
	return &incidentClosureDTO, nil
}

func (ics *IncidentClosureService) Create(incidentID uint, currentUserID uint, incidentClosureFormDTO *incidentDTOs.IncidentClosureFormDTO) (uint, error) {
	tx := ics.repo.BeginTransaction()

	if _, err := ics.incidentRegisterSvc.GetByID(incidentID); err != nil {
		return 0, err
	}

	dates, err := ics.parseIncidentClosureDates(incidentClosureFormDTO)
	if err != nil {
		return 0, err
	}

	incidentClosure := &models.IncidentClosure{
		IncidentRegisterID: incidentID,
		CompletionDate:     dates.CompletionDate,
		ResolutionSummary:  incidentClosureFormDTO.ResolutionSummary,
		Notes:              incidentClosureFormDTO.Notes,
		CreatedByID:        currentUserID,
	}

	id, err := ics.repo.Create(incidentClosure, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	// Save Lessons
	var incidentClosureLessons []*models.IncidentClosureLesson
	for _, lesson := range incidentClosureFormDTO.Lessons {
		lessonObj := &models.IncidentClosureLesson{
			IncidentClosureID: id,
			Lessons:           *lesson,
		}
		incidentClosureLessons = append(incidentClosureLessons, lessonObj)
	}

	if len(incidentClosureLessons) > 0 {
		if err := ics.repo.CreateLessons(incidentClosureLessons, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
	}
	tx.Commit()

	go func() {
		err := ics.activityLogSvc.LogActivity(
			models.CreateAction,
			models.IncidentClosure{}.TableName(),
			id,
			incidentClosureFormDTO,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.IncidentClosure{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return id, nil
}

func (ics *IncidentClosureService) Update(incidentID uint, currentUserID uint, incidentClosureFormDTO *incidentDTOs.IncidentClosureFormDTO) error {
	tx := ics.repo.BeginTransaction()
	if _, err := ics.incidentRegisterSvc.GetByID(incidentID); err != nil {
		return err
	}

	existingIncidentClosure, err := ics.GetByID(incidentID)
	if err != nil {
		return err
	}

	dates, err := ics.parseIncidentClosureDates(incidentClosureFormDTO)
	if err != nil {
		return err
	}
	// Update Incident Closure
	existingIncidentClosure.CompletionDate = dates.CompletionDate
	existingIncidentClosure.ResolutionSummary = incidentClosureFormDTO.ResolutionSummary
	existingIncidentClosure.Notes = incidentClosureFormDTO.Notes
	existingIncidentClosure.UpdatedByID = &currentUserID

	err = ics.repo.Update(existingIncidentClosure, tx)

	if err != nil {
		tx.Rollback()
		return err
	}

	existingIncidentClosureDetail, err := ics.GetFullDetailsByID(incidentID, tx)

	if err != nil {
		tx.Rollback()
		return err
	}

	// Collect Existing Lesson Lessonss
	existingClosureLessons := make(map[string]bool)
	for _, lesson := range existingIncidentClosureDetail.Lessons {
		existingClosureLessons[lesson.Lessons] = true
	}

	// Collect New Lesson Lessonss
	newClosureLessons := make(map[string]bool)
	for _, lesson := range incidentClosureFormDTO.Lessons {
		newClosureLessons[*lesson] = true
	}

	// Add New Lessons
	var lessonsToAdd []*models.IncidentClosureLesson
	for _, newLesson := range incidentClosureFormDTO.Lessons {
		if _, exists := existingClosureLessons[*newLesson]; !exists {
			// Create a new incident Closure lesson record
			lesson := &models.IncidentClosureLesson{
				IncidentClosureID: uint(existingIncidentClosureDetail.ID),
				Lessons:           *newLesson,
			}
			lessonsToAdd = append(lessonsToAdd, lesson)
		}
	}

	// Save the new lessons
	if len(lessonsToAdd) > 0 {
		if err := ics.repo.CreateLessons(lessonsToAdd, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	// Remove lessons not in the new list
	var lessonsToDelete []*uint
	for _, existingLesson := range existingIncidentClosureDetail.Lessons {
		if _, exists := newClosureLessons[existingLesson.Lessons]; !exists {
			// Delete the existing lesson from the repository
			lessonsToDelete = append(lessonsToDelete, &existingLesson.ID)
		}
	}

	// Delete the lessons
	if len(lessonsToDelete) > 0 {
		if err := ics.repo.DeleteLessons(lessonsToDelete, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()

	go func() {
		err := ics.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.IncidentClosure{}.TableName(),
			existingIncidentClosure.ID,
			incidentClosureFormDTO,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.UpdateAction,
				"tableName": models.IncidentClosure{}.TableName(),
				"recordID":  existingIncidentClosure.ID,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return nil
}

func (ics *IncidentClosureService) Delete(incidentID uint, currentUserID uint) (uint, error) {
	tx := ics.repo.BeginTransaction()
	var failedId uint

	if _, err := ics.incidentRegisterSvc.GetByID(incidentID); err != nil {
		return 0, err
	}
	incidentClosure, err := ics.GetFullDetailsByID(incidentID, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	if len(incidentClosure.Lessons) > 0 {
		var deleteLessonIds []*uint
		for _, lesson := range incidentClosure.Lessons {
			deleteLessonIds = append(deleteLessonIds, &lesson.ID)
		}
		if err := ics.repo.DeleteLessons(deleteLessonIds, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	if err = ics.repo.Delete(uint(incidentClosure.ID), tx); err != nil {
		failedId = uint(incidentClosure.ID)
		return failedId, err
	}
	tx.Commit()
	// Log successful delete activity
	go func() {
		err := ics.activityLogSvc.LogActivity(
			models.DeleteAction,
			models.IncidentClosure{}.TableName(),
			uint(incidentClosure.ID),
			nil,
			currentUserID,
			nil,
		)
		if err != nil {
			// Handle error if logging fails
			fmt.Printf("Failed to log delete activity: %v\n", err)
		}
	}()

	return failedId, nil
}

func (ics *IncidentClosureService) parseIncidentClosureDates(dto *incidentDTOs.IncidentClosureFormDTO) (*incidentDTOs.IncidentClosureParsedDates, error) {
	parsed := &incidentDTOs.IncidentClosureParsedDates{}

	var err error
	if dto.CompletionDate != "" {
		parsed.CompletionDate, err = utility.ParseDate(dto.CompletionDate)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid completion date"), "invalid_date", "completion_date")
		}
	}

	return parsed, nil
}
