package services

import (
	"errors"
	"fmt"
	incidentDTOs "jethur/grc/incident/dtos"
	incidentRepositories "jethur/grc/incident/repositories"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	"jethur/utility"
	"time"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type IncidentPostReviewServiceInterface interface {
	GetByID(incidentID uint) (*models.IncidentPostReview, error)
	GetFullDetailsByID(incidentID uint, tx ...utility.DatabaseTransaction) (*incidentDTOs.IncidentPostReviewDetailDTO, error)
	Create(incidentID uint, currentUserID uint, IncidentPostReviewFormDTO *incidentDTOs.IncidentPostReviewFormDTO) (uint, error)
	Update(incidentID uint, currentUserID uint, incident *incidentDTOs.IncidentPostReviewFormDTO) error
	Delete(incidentID uint, currentUserID uint) (uint, error)
}

// Service and implementations
type IncidentPostReviewService struct {
	repo                               incidentRepositories.IncidentPostReviewRepositoryInterface
	activityLogSvc                     logServices.ActivityLogServiceInterface
	userSvc                            organizationServices.UserServiceInterface
	incidentPostReviewEffectivenessSvc IncidentPostReviewEffectivenessServiceInterface
	incidentRegisterSvc                IncidentRegisterServiceInterface
	incidentTimelineSvc                IncidentTimelineServiceInterface
}

func NewIncidentPostReviewService(
	repo incidentRepositories.IncidentPostReviewRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	incidentPostReviewEffectivenessSvc IncidentPostReviewEffectivenessServiceInterface,
	incidentRegisterSvc IncidentRegisterServiceInterface,
	incidentTimelineSvc IncidentTimelineServiceInterface,
) IncidentPostReviewServiceInterface {
	return &IncidentPostReviewService{
		repo:                               repo,
		activityLogSvc:                     activityLogSvc,
		userSvc:                            userSvc,
		incidentPostReviewEffectivenessSvc: incidentPostReviewEffectivenessSvc,
		incidentRegisterSvc:                incidentRegisterSvc,
		incidentTimelineSvc:                incidentTimelineSvc,
	}
}

func (iprs IncidentPostReviewService) GetByID(incidentID uint) (*models.IncidentPostReview, error) {
	if _, err := iprs.incidentRegisterSvc.GetByID(incidentID); err != nil {
		return nil, err
	}

	incidentPostReview, err := iprs.repo.GetByID(incidentID)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, incidentID)
		}
		return nil, err
	}
	return incidentPostReview, nil
}

func (iprs IncidentPostReviewService) GetFullDetailsByID(incidentID uint, tx ...utility.DatabaseTransaction) (*incidentDTOs.IncidentPostReviewDetailDTO, error) {
	var incidentPostReviewDTO incidentDTOs.IncidentPostReviewDetailDTO

	if _, err := iprs.incidentRegisterSvc.GetByID(incidentID); err != nil {
		return nil, err
	}

	incidentPostReview, err := iprs.repo.GetFullDetailsByID(incidentID, tx...)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, incidentID)
		}
		return nil, err
	}

	err = mapstructure.Decode(incidentPostReview, &incidentPostReviewDTO)
	if err != nil {
		return nil, err
	}

	//todo :- also add extra fields which required in details screen
	if incidentPostReview.PostReviewDate != nil {
		postReviewDate := incidentPostReview.PostReviewDate.Format("2006-01-02")
		incidentPostReviewDTO.PostReviewDateString = &postReviewDate
	}

	incidentPostReviewDTO.CreatedAt = incidentPostReview.CreatedAt
	incidentPostReviewDTO.UpdatedAt = incidentPostReview.UpdatedAt
	return &incidentPostReviewDTO, nil
}

func (iprs *IncidentPostReviewService) Create(incidentID uint, currentUserID uint, incidentPostReviewFormDTO *incidentDTOs.IncidentPostReviewFormDTO) (uint, error) {
	tx := iprs.repo.BeginTransaction()

	if _, err := iprs.incidentRegisterSvc.GetByID(incidentID); err != nil {
		return 0, err
	}
	if err := iprs.validateIncidentPostReviewFormDTO(incidentPostReviewFormDTO); err != nil {
		return 0, err
	}

	dates, err := iprs.parseIncidentPostReviewDates(incidentPostReviewFormDTO)
	if err != nil {
		return 0, err
	}

	incidentPostReview := &models.IncidentPostReview{
		IncidentRegisterID:                incidentID,
		IncidentPostReviewEffectivenessID: int(incidentPostReviewFormDTO.IncidentPostReviewEffectivenessID),
		PostReviewDate:                    dates.PostReviewDate,
		Summary:                           incidentPostReviewFormDTO.Summary,
		CreatedByID:                       currentUserID,
	}

	id, err := iprs.repo.Create(incidentPostReview, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	// Save Participants
	var postReviewUsers []*models.IncidentPostReviewUser
	for _, userId := range incidentPostReviewFormDTO.UserIds {
		user := &models.IncidentPostReviewUser{
			IncidentPostReviewID: id,
			UserID:               *userId,
			CreatedByID:          currentUserID,
		}
		postReviewUsers = append(postReviewUsers, user)
	}

	if len(postReviewUsers) > 0 {
		if err := iprs.repo.CreateParticipants(postReviewUsers, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	// Save Lessons
	var incidentPostReviewLessons []*models.IncidentPostReviewLesson
	for _, lesson := range incidentPostReviewFormDTO.Lessons {
		lessonObj := &models.IncidentPostReviewLesson{
			IncidentPostReviewID: id,
			Lesson:               *lesson,
			CreatedByID:          currentUserID,
		}
		incidentPostReviewLessons = append(incidentPostReviewLessons, lessonObj)
	}

	if len(incidentPostReviewLessons) > 0 {
		if err := iprs.repo.CreateLessons(incidentPostReviewLessons, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
	}
	// Save Recommendations
	var incidentPostReviewRecommendations []*models.IncidentPostReviewRecommendation
	for _, recommendation := range incidentPostReviewFormDTO.Recommendations {
		recommendationObj := &models.IncidentPostReviewRecommendation{
			IncidentPostReviewID: id,
			Recommendation:       *recommendation,
			CreatedByID:          currentUserID,
		}
		incidentPostReviewRecommendations = append(incidentPostReviewRecommendations, recommendationObj)
	}

	if len(incidentPostReviewRecommendations) > 0 {
		if err := iprs.repo.CreateRecommendations(incidentPostReviewRecommendations, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
	}
	//Add timeline
	//todo
	timelineType, err := iprs.incidentTimelineSvc.GetByType("pir")
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	//todo
	timelineStatus, err := iprs.incidentTimelineSvc.GetByStatusType("completed")
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	currentTime := time.Now()
	timelineFormDTO := &incidentDTOs.IncidentTimelineFormDTO{
		IncidentRegisterID:     incidentID,
		IncidentTimelineTypeID: timelineType.ID,
		TimelineStatusID:       timelineStatus.ID,
		PerformedByID:          currentUserID,
		PerformedAt:            &currentTime,
	}
	if err := iprs.incidentTimelineSvc.Create(timelineFormDTO, tx); err != nil {
		tx.Rollback()
		return 0, err
	}
	tx.Commit()

	go func() {
		err := iprs.activityLogSvc.LogActivity(
			models.CreateAction,
			models.IncidentPostReview{}.TableName(),
			id,
			incidentPostReviewFormDTO,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.IncidentPostReview{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return id, nil
}

func (iprs *IncidentPostReviewService) Update(incidentID uint, currentUserID uint, incidentPostReviewFormDTO *incidentDTOs.IncidentPostReviewFormDTO) error {
	tx := iprs.repo.BeginTransaction()
	if _, err := iprs.incidentRegisterSvc.GetByID(incidentID); err != nil {
		return err
	}

	existingIncidentPostReview, err := iprs.GetByID(incidentID)
	if err != nil {
		return err
	}
	if err := iprs.validateIncidentPostReviewFormDTO(incidentPostReviewFormDTO); err != nil {
		return err
	}

	dates, err := iprs.parseIncidentPostReviewDates(incidentPostReviewFormDTO)
	if err != nil {
		return err
	}
	// Update
	existingIncidentPostReview.PostReviewDate = dates.PostReviewDate
	existingIncidentPostReview.IncidentPostReviewEffectivenessID = int(incidentPostReviewFormDTO.IncidentPostReviewEffectivenessID)
	existingIncidentPostReview.Summary = incidentPostReviewFormDTO.Summary
	existingIncidentPostReview.UpdatedByID = &currentUserID

	err = iprs.repo.Update(existingIncidentPostReview, tx)

	if err != nil {
		tx.Rollback()
		return err
	}

	existingIncidentPostReviewDetail, err := iprs.GetFullDetailsByID(incidentID, tx)

	if err != nil {
		tx.Rollback()
		return err
	}
	//Update participants
	existingUserIDs := make(map[uint]bool)
	for _, user := range existingIncidentPostReviewDetail.Users {
		existingUserIDs[user.ID] = true
	}
	// Collect new user IDs
	newUserIDs := make(map[uint]bool)
	for _, user := range incidentPostReviewFormDTO.UserIds {
		newUserIDs[*user] = true
	}
	// Add new users
	var users []*models.IncidentPostReviewUser
	for _, newUser := range incidentPostReviewFormDTO.UserIds {
		if _, exists := existingUserIDs[*newUser]; !exists {
			// Create a new IncidentPostReviewUser  record
			user := &models.IncidentPostReviewUser{
				IncidentPostReviewID: existingIncidentPostReview.ID,
				UserID:               *newUser,
				CreatedByID:          currentUserID,
			}
			users = append(users, user)
		}
	}
	// Save the new user
	if len(users) > 0 {
		if err := iprs.repo.CreateParticipants(users, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	// Remove users not in the new list
	var deleteUserIDs []*uint
	for _, existingUser := range existingIncidentPostReviewDetail.Users {
		if _, exists := newUserIDs[existingUser.ID]; !exists {
			// Delete the existing user from the repository
			deleteUserIDs = append(deleteUserIDs, &existingUser.ID)
		}
	}
	// Delete the user
	if len(deleteUserIDs) > 0 {
		if err := iprs.repo.DeleteParticipants(deleteUserIDs, tx); err != nil {
			tx.Rollback()
			return err
		}
	}
	//Update Lessons
	// Collect Existing Lesson Lessons
	existingClosureLessons := make(map[string]bool)
	for _, lesson := range existingIncidentPostReviewDetail.Lessons {
		existingClosureLessons[lesson.Lesson] = true
	}

	// Collect New Lesson Lessonss
	newClosureLessons := make(map[string]bool)
	for _, lesson := range incidentPostReviewFormDTO.Lessons {
		newClosureLessons[*lesson] = true
	}

	// Add New Lessons
	var lessonsToAdd []*models.IncidentPostReviewLesson
	for _, newLesson := range incidentPostReviewFormDTO.Lessons {
		if _, exists := existingClosureLessons[*newLesson]; !exists {
			// Create a new incident Closure lesson record
			lesson := &models.IncidentPostReviewLesson{
				IncidentPostReviewID: existingIncidentPostReview.ID,
				Lesson:               *newLesson,
				CreatedByID:          currentUserID,
			}
			lessonsToAdd = append(lessonsToAdd, lesson)
		}
	}

	// Save the new lessons
	if len(lessonsToAdd) > 0 {
		if err := iprs.repo.CreateLessons(lessonsToAdd, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	// Remove lessons not in the new list
	var lessonsToDelete []*uint
	for _, existingLesson := range existingIncidentPostReviewDetail.Lessons {
		if _, exists := newClosureLessons[existingLesson.Lesson]; !exists {
			// Delete the existing lesson from the repository
			lessonsToDelete = append(lessonsToDelete, &existingLesson.ID)
		}
	}

	// Delete the lessons
	if len(lessonsToDelete) > 0 {
		if err := iprs.repo.DeleteLessons(lessonsToDelete, tx); err != nil {
			tx.Rollback()
			return err
		}
	}
	//Update Reccomendations
	// Collect Existing Recommendations
	existingRecommendations := make(map[string]bool)
	for _, recommendation := range existingIncidentPostReviewDetail.Recommendations {
		existingRecommendations[recommendation.Recommendation] = true
	}

	// Collect New Recommendations
	newRecommendations := make(map[string]bool)
	for _, recommendation := range incidentPostReviewFormDTO.Recommendations {
		newRecommendations[*recommendation] = true
	}

	// Update Existing Recommendations
	var recommendationsToAdd []*models.IncidentPostReviewRecommendation
	for _, newRec := range incidentPostReviewFormDTO.Recommendations {
		if _, exists := existingRecommendations[*newRec]; !exists {
			// Create a new recommendation record
			recommendation := &models.IncidentPostReviewRecommendation{
				IncidentPostReviewID: existingIncidentPostReview.ID,
				Recommendation:       *newRec,
				CreatedByID:          currentUserID,
			}
			recommendationsToAdd = append(recommendationsToAdd, recommendation)
		}
	}

	// Save the updated recommendations
	if len(recommendationsToAdd) > 0 {
		if err := iprs.repo.CreateRecommendations(recommendationsToAdd, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	// Delete Recommendations not in the new list
	var recommendationsToDelete []*uint
	for _, existingRecommendation := range existingIncidentPostReviewDetail.Recommendations {
		if _, exists := newRecommendations[existingRecommendation.Recommendation]; !exists {
			// Delete the existing recommendation from the repository
			recommendationsToDelete = append(recommendationsToDelete, &existingRecommendation.ID)
		}
	}

	// Delete the recommendations
	if len(recommendationsToDelete) > 0 {
		if err := iprs.repo.DeleteRecommendations(recommendationsToDelete, tx); err != nil {
			tx.Rollback()
			return err
		}
	}
	//Update timeline
	//todo
	timelineType, err := iprs.incidentTimelineSvc.GetByType("pir")
	if err != nil {
		tx.Rollback()
		return err
	}
	timelineStatus, err := iprs.incidentTimelineSvc.GetByStatusType("completed")
	if err != nil {
		tx.Rollback()
		return err
	}
	currentTime := time.Now()
	timelineFormDTO := &incidentDTOs.IncidentTimelineFormDTO{
		PerformedByID:    currentUserID,
		TimelineStatusID: timelineStatus.ID,
		PerformedAt:      &currentTime,
	}
	if err := iprs.incidentTimelineSvc.Update(incidentID, timelineFormDTO, timelineType.ID, tx); err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	go func() {
		err := iprs.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.IncidentPostReview{}.TableName(),
			existingIncidentPostReview.ID,
			incidentPostReviewFormDTO,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    models.UpdateAction,
				"tableName": models.IncidentPostReview{}.TableName(),
				"recordID":  existingIncidentPostReview.ID,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return nil
}

func (iprs *IncidentPostReviewService) Delete(incidentID uint, currentUserID uint) (uint, error) {
	tx := iprs.repo.BeginTransaction()
	var failedId uint

	if _, err := iprs.incidentRegisterSvc.GetByID(incidentID); err != nil {
		return 0, err
	}
	incidentPostReview, err := iprs.GetFullDetailsByID(incidentID, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	if len(incidentPostReview.Users) > 0 {
		var deleteUserIds []*uint
		for _, user := range incidentPostReview.Users {
			deleteUserIds = append(deleteUserIds, &user.ID)
		}
		if err := iprs.repo.DeleteParticipants(deleteUserIds, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
	}
	if len(incidentPostReview.Lessons) > 0 {
		var deleteLessonIds []*uint
		for _, lesson := range incidentPostReview.Lessons {
			deleteLessonIds = append(deleteLessonIds, &lesson.ID)
		}
		if err := iprs.repo.DeleteLessons(deleteLessonIds, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
	}
	if len(incidentPostReview.Recommendations) > 0 {
		var deleteRecommendationIds []*uint
		for _, recommendation := range incidentPostReview.Recommendations {
			deleteRecommendationIds = append(deleteRecommendationIds, &recommendation.ID)
		}
		if err := iprs.repo.DeleteRecommendations(deleteRecommendationIds, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
	}
	if err = iprs.repo.Delete(uint(incidentPostReview.ID), tx); err != nil {
		failedId = uint(incidentPostReview.ID)
		return failedId, err
	}
	//Delete-> update timeline performedAt as null
	//todo
	timelineType, err := iprs.incidentTimelineSvc.GetByType("pir")
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	//todo
	timelineStatus, err := iprs.incidentTimelineSvc.GetByStatusType("pending")
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	timelineFormDTO := &incidentDTOs.IncidentTimelineFormDTO{
		PerformedByID:    currentUserID,
		TimelineStatusID: timelineStatus.ID,
		PerformedAt:      nil,
	}
	if err := iprs.incidentTimelineSvc.Update(incidentID, timelineFormDTO, timelineType.ID, tx); err != nil {
		tx.Rollback()
		return 0, err
	}
	tx.Commit()
	// Log successful delete activity
	go func() {
		err := iprs.activityLogSvc.LogActivity(
			models.DeleteAction,
			models.IncidentPostReview{}.TableName(),
			uint(incidentPostReview.ID),
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
func (iprs *IncidentPostReviewService) parseIncidentPostReviewDates(dto *incidentDTOs.IncidentPostReviewFormDTO) (*incidentDTOs.IncidentPostReviewParsedDates, error) {
	parsed := &incidentDTOs.IncidentPostReviewParsedDates{}

	var err error
	if dto.PostReviewDate != "" {
		parsed.PostReviewDate, err = utility.ParseDate(dto.PostReviewDate)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid post review date"), "invalid_date", "post_review_date")
		}
	}

	return parsed, nil
}

func (iprs *IncidentPostReviewService) validateIncidentPostReviewFormDTO(incidentPostReviewFormDTO *incidentDTOs.IncidentPostReviewFormDTO) error {
	if _, err := iprs.incidentPostReviewEffectivenessSvc.GetByID(uint(incidentPostReviewFormDTO.IncidentPostReviewEffectivenessID)); err != nil {
		return utility.NewValidationError(fmt.Errorf("incident incident post review effectiveness id does not exist"), "not_exist", "incident post review effectiveness")
	}
	//User exist
	users := incidentPostReviewFormDTO.UserIds
	// If no user are provided, simply return nil to exit the function
	if len(users) == 0 {
		return nil
	}

	// Create a map to store user IDs
	userIDs := make(map[uint]bool)

	// Check if each user exists
	for _, user := range users {
		// Check if the user ID already exists in the map
		if userIDs[*user] {
			// If the user ID already exists, return a validation error
			return utility.NewValidationError(fmt.Errorf("duplicate user ID: %d", user), "duplicate", "user")
		}

		// Check if the user exists in the database
		if _, err := iprs.userSvc.GetByID(*user); err != nil {
			return utility.NewValidationError(fmt.Errorf("user does not exist for ID: %d", user), "not_exist", "user")
		}

		// Add the user ID to the map
		userIDs[*user] = true
	}

	return nil
}
