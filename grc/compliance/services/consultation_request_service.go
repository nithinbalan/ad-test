package services

import (
	"errors"
	"fmt"
	"log"
	"time"

	"jethur/grc/compliance/dtos"
	"jethur/grc/compliance/repositories"
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
		ReferenceNumber:             referenceNumber,
		DateRaised:                  time.Now(),
		RequestingDepartmentUnitID:  responsibleUnitID,
		ConsultedDepartmentUnitID:   dto.ConsultedDepartmentID,
		RequestDetails:              dto.RequestDetails,
		ConsultationRequestStatusID: &dto.ConsultationRequestStatusID,
		Reply:                       dto.Reply,
		Notes:                       dto.Notes,
		CreatedBYID:                 currentUserID,
		CreatedAt:                   time.Now(),
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
			"compliance", "consultation_request", "", nil,
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

	if dto.Reply != nil {
		existing.ReplyDate = &now
	}

	existing.RequestDetails = dto.RequestDetails
	existing.ConsultationRequestStatusID = &dto.ConsultationRequestStatusID
	existing.Reply = dto.Reply
	existing.Notes = dto.Notes
	existing.UpdatedByID = &currentUserID
	existing.UpdatedAt = &now
	if err := s.repo.Update(existing, tx); err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	updated, err := s.repo.GetByID(id)

	fieldsToCompare := []string{"RequestDetails", "ConsultationRequestStatusID", "Reply", "ReplyDate", "Notes"}
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
			"compliance", "consultation_request", "", nil,
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
			"compliance",
			"consultation_request",
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

// Validation
func (s *ConsultationRequestService) validateConsultationRequestDTO(dto *dtos.ConsultationRequestFormDTO) error {

	if _, err := s.orgUnitSvc.GetByID(dto.ConsultedDepartmentID); err != nil {
		return utility.NewValidationError(fmt.Errorf("consulted organization unit does not exist"), "not_exist", "consulted_organization_unit")
	}

	if _, err := s.statusSvc.GetByID(dto.ConsultationRequestStatusID); err != nil {
		return utility.NewValidationError(errors.New("This status not exist"), "not_exist", "consultation_request_status")
	}

	return nil
}
