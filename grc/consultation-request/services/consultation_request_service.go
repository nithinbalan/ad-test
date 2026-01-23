package services

import (
	"errors"
	"fmt"
	"log"
	"time"

	"jethur/grc/consultation-request/dtos"
	"jethur/grc/consultation-request/repositories"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	settingsServices "jethur/grc/settings/services"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

// Interface for dependency injection
type ConsultationRequestServiceInterface interface {
	GetAll(params *dtos.ConsultationRequestListQueryParamsDTO) (*dtos.ConsultationRequestListResponseDTO, error)
	GetByID(id uint) (*models.ConsultationRequest, error)
	GetFullDetailsByID(id uint) (*dtos.ConsultationRequestDetailDTO, error)
	Create(currentUserID uint, dto *dtos.ConsultationRequestFormDTO) (uint, error)
	Update(id uint, dto *dtos.ConsultationRequestFormDTO, currentUserID uint) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	UpdateResponsibleUser(id uint, newResponsibleUserID uint, currentUserID uint) error
	UpdateNotes(consultationRequestID uint, notes string, updatedByID uint) error
	CreateConsultationRequestComment(consultationRequestID uint, currentUserID uint, dto *dtos.ConsultationRequestCommentFormDTO) (uint, error)
	GetAllConsultationRequestComments(consultationRequestID uint) ([]*dtos.ConsultationRequestCommentListDTO, error)
}

// Implementation
type ConsultationRequestService struct {
	repo           repositories.ConsultationRequestRepositoryInterface
	userSvc        organizationServices.UserServiceInterface
	activityLogSvc logServices.ActivityLogServiceInterface
	orgUnitSvc     settingsServices.OrganizationUnitServiceInterface
	statusSvc      ConsultationRequestStatusServiceInterface
}

// Constructor
func NewConsultationRequestService(
	repo repositories.ConsultationRequestRepositoryInterface,
	userSvc organizationServices.UserServiceInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	orgUnitSvc settingsServices.OrganizationUnitServiceInterface,
	statusSvc ConsultationRequestStatusServiceInterface,
) *ConsultationRequestService {
	return &ConsultationRequestService{
		repo:           repo,
		userSvc:        userSvc,
		activityLogSvc: activityLogSvc,
		orgUnitSvc:     orgUnitSvc,
		statusSvc:      statusSvc,
	}
}

func (s *ConsultationRequestService) GetAll(params *dtos.ConsultationRequestListQueryParamsDTO) (*dtos.ConsultationRequestListResponseDTO, error) {
	return s.repo.GetAll(params)
}
func (s *ConsultationRequestService) GetAllConsultationRequestComments(consultationRequestID uint) ([]*dtos.ConsultationRequestCommentListDTO, error) {
	return s.repo.GetAllConsultationRequestComments(consultationRequestID)
}

func (s *ConsultationRequestService) GetByID(id uint) (*models.ConsultationRequest, error) {
	request, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return request, nil
}

func (s *ConsultationRequestService) GetFullDetailsByID(id uint) (*dtos.ConsultationRequestDetailDTO, error) {
	var detailDTO dtos.ConsultationRequestDetailDTO

	request, err := s.repo.GetFullDetailsByID(id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	err = mapstructure.Decode(request, &detailDTO)
	if err != nil {
		fmt.Println("Error occurred:", err)
		return nil, err
	}

	if !request.DateRaised.IsZero() {
		detailDTO.RaisedDate = request.DateRaised.Format("2006-01-02")
	}

	if request.ReplyDate != nil && !request.ReplyDate.IsZero() {
		detailDTO.ReplyedDate = request.ReplyDate.Format("2006-01-02")
	}

	if !request.CreatedAt.IsZero() {
		detailDTO.CreatedDate = request.CreatedAt.Format("2006-01-02")
	}
	return &detailDTO, nil
}

// Create
func (s *ConsultationRequestService) Create(currentUserID uint, dto *dtos.ConsultationRequestFormDTO) (uint, error) {

	if err := s.validateConsultationRequestDTO(dto); err != nil {

		return 0, err
	}

	referenceNumber, err := s.repo.GenerateReferenceNumber()
	if err != nil {
		return 0, err
	}

	responsibleUnitID, err := s.userSvc.GetUserPrimaryUnitID(currentUserID)
	if err != nil {
		return 0, utility.NewValidationError(fmt.Errorf("organization unit not set for responsible user"), "not_exist", "bcp communication plan responsible unit")
	}
	tx := s.repo.BeginTransaction()

	request := &models.ConsultationRequest{
		Title:                      dto.Title,
		ReferenceNumber:            referenceNumber,
		DateRaised:                 time.Now(),
		RequestingDepartmentUnitID: responsibleUnitID,
		ConsultedDepartmentUnitID:  dto.ConsultedDepartmentID,
		RequestDetails:             dto.RequestDetails,
		CreatedBYID:                currentUserID,
		CreatedAt:                  time.Now(),
	}

	id, err := s.repo.Create(request, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	tx.Commit()

	// Async activity logging
	go func() {
		fullDetails, err := s.GetFullDetailsByID(id)
		if err != nil {
			return
		}
		err = s.activityLogSvc.LogActivity(
			models.CreateAction,
			models.ConsultationRequest{}.TableName(),
			id,
			fullDetails,
			currentUserID,
			"consultation", "consultation_request", "", nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.ConsultationRequest{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return id, nil
}

func (s *ConsultationRequestService) CreateConsultationRequestComment(consultationRequestID uint, currentUserID uint, dto *dtos.ConsultationRequestCommentFormDTO) (uint, error) {

	// 1️⃣ Validate consultation request exists
	existingRequest, err := s.repo.GetByID(consultationRequestID)
	if err != nil {
		return 0, utility.NewResourceNotFoundError(err, consultationRequestID)
	}
	// 2️⃣ Authorization: user must belong to consulted department
	reqUsers, err := s.userSvc.GetUsersByOrganizationUnit(existingRequest.RequestingDepartmentUnitID)
	if err != nil {
		return 0, err
	}
	// 2️⃣ Authorization: user must belong to consulted department
	consUsers, err := s.userSvc.GetUsersByOrganizationUnit(existingRequest.ConsultedDepartmentUnitID)
	if err != nil {
		return 0, err
	}
	authorized := false

	status, err := s.statusSvc.GetByID(dto.ConsultationRequestStatusID)
	if status.Type == "closed" {
		for _, u := range reqUsers {
			if u.ID == currentUserID {
				authorized = true
				break
			}
		}
		if !authorized {
			return 0, utility.NewForbiddenError("you do not have authority to update responsible user")
		}
	} else {
		for _, u := range consUsers {
			if u.ID == currentUserID {
				authorized = true
				break
			}
		}
		if !authorized {
			return 0, utility.NewForbiddenError("you do not have authority to update responsible user")
		}
	}

	// 3️⃣ Start transaction
	tx := s.repo.BeginTransaction()
	now := time.Now()

	// 4️⃣ Create comment model
	comment := &models.ConsultationRequestComment{
		ConsultationRequestID:       consultationRequestID,
		ConsultationRequestStatusID: dto.ConsultationRequestStatusID,
		Comment:                     dto.Comment,
		CreatedByID:                 currentUserID,
		CreatedAt:                   now,
	}

	commentID, err := s.repo.CreateComment(comment, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	// 5️⃣ Optional: update consultation request status
	existingRequest.ConsultationRequestStatusID = &dto.ConsultationRequestStatusID
	existingRequest.UpdatedByID = &currentUserID
	existingRequest.UpdatedAt = &now

	if err := s.repo.Update(existingRequest, tx); err != nil {
		tx.Rollback()
		return 0, err
	}

	tx.Commit()

	// 6️⃣ Async activity logging
	go func() {
		fullDetails, err := s.GetFullDetailsByID(consultationRequestID)
		if err != nil {
			return
		}

		err = s.activityLogSvc.LogActivity(
			models.CreateAction,
			models.ConsultationRequestComment{}.TableName(),
			commentID,
			fullDetails,
			currentUserID,
			"consultation_request",
			"consultation",
			"",
			nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":   models.CreateAction,
				"recordID": commentID,
				"userID":   currentUserID,
				"error":    err.Error(),
			}
			utility.LogError("Failed to log consultation request comment activity", fields)
		}
	}()

	return commentID, nil
}

// Update
func (s *ConsultationRequestService) Update(id uint, dto *dtos.ConsultationRequestFormDTO, currentUserID uint) error {
	existing, err := s.repo.GetByID(id)

	if err != nil {
		return utility.NewResourceNotFoundError(err, id)
	}

	if err := s.validateConsultationRequestDTO(dto); err != nil {
		return err
	}

	tx := s.repo.BeginTransaction()
	now := time.Now()

	existing.RequestDetails = dto.RequestDetails
	existing.ConsultedDepartmentUnitID = dto.ConsultedDepartmentID
	existing.Title = dto.Title
	existing.UpdatedByID = &currentUserID
	existing.UpdatedAt = &now
	if err := s.repo.Update(existing, tx); err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	updated, err := s.repo.GetByID(id)

	fieldsToCompare := []string{"Title", "ConsultedDepartmentUnitID", "RequestDetails"}
	changes, err := utility.CompareFields(existing, updated, fieldsToCompare)
	if err != nil {
		log.Println("error :", err.Error())
	}

	changedDataJSON, err := utility.ChangesToJSON(changes)
	if err != nil {
		log.Println("error :", err.Error())
	}

	go func() {

		err := s.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.ConsultationRequest{}.TableName(),
			id,
			changedDataJSON,
			currentUserID,
			"consultation_request", "consultation", "", nil,
		)
		if err != nil {
			fmt.Printf("Failed to log bulk delete activity: %v\n", err)
		}
	}()

	return nil
}

func (s *ConsultationRequestService) UpdateResponsibleUser(id uint, responsibleUserID uint, currentUserID uint) error {
	existing, err := s.repo.GetByID(id)
	if err != nil {
		return utility.NewResourceNotFoundError(err, id)
	}

	users, err := s.userSvc.GetUsersByOrganizationUnit(existing.ConsultedDepartmentUnitID)
	if err != nil {
		return err
	}

	authorized := false
	for _, u := range users {
		if u.ID == currentUserID {
			authorized = true
			break
		}
	}

	if !authorized {
		return utility.NewForbiddenError("you do not have authority to update responsible user")
	}

	// Check that the new responsible user also belongs to the same department
	userAuthorized := false
	for _, u := range users {
		if u.ID == responsibleUserID {
			userAuthorized = true
			break
		}
	}

	if !userAuthorized {
		return utility.NewValidationError(
			fmt.Errorf("the selected responsible user does not belong to the consulted department"),
			"invalid",
			"responsible_user_id",
		)
	}
	_, err = s.userSvc.GetByID(responsibleUserID)
	if err != nil {
		return utility.NewValidationError(
			fmt.Errorf("user with ID %d does not exist", responsibleUserID),
			"not_exist",
			"responsible_user_id",
		)
	}

	tx := s.repo.BeginTransaction()
	now := time.Now()

	existing.ResponsibleUserID = &responsibleUserID
	existing.UpdatedByID = &currentUserID
	existing.UpdatedAt = &now

	if err := s.repo.Update(existing, tx); err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	go func() {
		fullDetails, _ := s.GetFullDetailsByID(id)

		err := s.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.ConsultationRequest{}.TableName(),
			id,
			fullDetails,
			currentUserID,
			"consultation_request",
			"consultation",
			"",
			nil,
		)
		if err != nil {
			fmt.Printf("Failed to log update responsible user activity: %v\n", err)
		}
	}()

	return nil
}

// Delete
func (s *ConsultationRequestService) Delete(ids []uint, currentUserID uint) ([]uint, error) {

	var failedIds []uint
	var successfulIds []uint
	var deletedConsultationRequests []map[string]interface{}

	for _, id := range ids {
		consultationRequest, err := s.repo.GetByID(id)
		if err != nil {
			failedIds = append(failedIds, id)
			continue
		}

		err = s.repo.Delete(id)
		if err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
			deletedConsultationRequests = append(deletedConsultationRequests, map[string]interface{}{
				"id":               id,
				"reference_number": consultationRequest.ReferenceNumber,
				"request_details":  consultationRequest.RequestDetails,
			})
		}
	}
	if len(successfulIds) > 0 {
		go func() {
			err := s.activityLogSvc.LogBulkActivities(
				models.DeleteAction, // Action type
				models.ConsultationRequest{}.TableName(),
				successfulIds,
				deletedConsultationRequests, // The slice of deleted request details
				currentUserID,
				nil,
			)
			if err != nil {
				// Optional: save to log file or just print
				fmt.Printf("Failed to log ConsultationRequest delete activity: %v\n", err)
			}
		}()
	}

	return failedIds, nil
}
func (s *ConsultationRequestService) UpdateNotes(id uint, notes string, currentUserID uint) error {
	// Fetch existing consultation request
	existing, err := s.repo.GetByID(id)
	if err != nil {
		return utility.NewResourceNotFoundError(err, id)
	}

	// Optional: add validation if needed (e.g., notes length)
	if len(notes) == 0 {
		return utility.NewValidationError(fmt.Errorf("notes cannot be empty"), "required", "notes")
	}
	consultedDepartmentUnit := existing.ConsultedDepartmentUnitID
	users, err := s.userSvc.GetUsersByOrganizationUnit(consultedDepartmentUnit)
	if err != nil {
		return err
	}

	authorized := false
	for _, u := range users {
		if u.ID == currentUserID {
			authorized = true
			break
		}
	}

	if !authorized {
		return utility.NewForbiddenError("you do not have authority to update responsible user")
	}

	tx := s.repo.BeginTransaction()
	now := time.Now()

	// Update fields
	existing.Notes = &notes
	existing.UpdatedByID = &currentUserID
	existing.UpdatedAt = &now

	// Save update
	if err := s.repo.Update(existing, tx); err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	// Log activity asynchronously
	go func() {
		fullDetails, _ := s.GetFullDetailsByID(id)

		err := s.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.ConsultationRequest{}.TableName(),
			id,
			fullDetails,
			currentUserID,
			"consultation_request",
			"consultation",
			"",
			nil,
		)
		if err != nil {
			fmt.Printf("Failed to log update notes activity: %v\n", err)
		}
	}()

	return nil
}

// Validation
func (s *ConsultationRequestService) validateConsultationRequestDTO(dto *dtos.ConsultationRequestFormDTO) error {

	if _, err := s.orgUnitSvc.GetByID(dto.ConsultedDepartmentID); err != nil {
		return utility.NewValidationError(fmt.Errorf("consulted organization unit does not exist"), "not_exist", "consulted_organization_unit")
	}

	// if _, err := s.statusSvc.GetByID(dto.ConsultationRequestStatusID); err != nil {
	// 	return utility.NewValidationError(errors.New("This status not exist"), "not_exist", "consultation_request_status")
	// }

	return nil
}
