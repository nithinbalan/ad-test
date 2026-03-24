package services

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"jethur/config"
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
	GetAll(params *dtos.ConsultationRequestListQueryParamsDTO, currentUserID uint) (*dtos.ConsultationRequestListResponseDTO, error)
	GetByID(id uint) (*models.ConsultationRequest, error)
	GetFullDetailsByID(id uint) (*dtos.ConsultationRequestDetailDTO, error)
	Create(currentUserID uint, dto *dtos.ConsultationRequestFormDTO) (uint, error)
	Update(id uint, dto *dtos.ConsultationRequestFormDTO, currentUserID uint) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	UpdateResponsibleUser(id uint, newResponsibleUserID uint, currentUserID uint) error
	UpdateNotes(consultationRequestID uint, notes string, updatedByID uint) error
	CreateConsultationRequestComment(consultationRequestID uint, currentUserID uint, dto *dtos.ConsultationRequestCommentFormDTO) (uint, error)
	UpdateConsultationRequestComment(consultationRequestID uint, commentID uint, currentUserID uint, dto *dtos.ConsultationRequestCommentUpdateFormDTO) error
	DeleteConsultationRequestComments(consultationRequestID uint, ids []uint, currentUserID uint) ([]uint, error)
	GetAllConsultationRequestComments(consultationRequestID uint) ([]*dtos.ConsultationRequestCommentListDTO, error)
	ExportConsultationRequests(data []*dtos.ConsultationRequestListDTO, currentUserID uint) (*os.File, error)
	GetDocumentPath(consultationRequestID uint, commentID uint, token string) (string, error)
}

// Implementation
type ConsultationRequestService struct {
	repo           repositories.ConsultationRequestRepositoryInterface
	userSvc        organizationServices.UserServiceInterface
	activityLogSvc logServices.ActivityLogServiceInterface
	orgUnitSvc     settingsServices.OrganizationUnitServiceInterface
	statusSvc      ConsultationRequestStatusServiceInterface
	userSettingSvc organizationServices.UserSettingServiceInterface
}

// Constructor
func NewConsultationRequestService(
	repo repositories.ConsultationRequestRepositoryInterface,
	userSvc organizationServices.UserServiceInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	orgUnitSvc settingsServices.OrganizationUnitServiceInterface,
	statusSvc ConsultationRequestStatusServiceInterface,
	userSettingSvc organizationServices.UserSettingServiceInterface,
) *ConsultationRequestService {
	return &ConsultationRequestService{
		repo:           repo,
		userSvc:        userSvc,
		activityLogSvc: activityLogSvc,
		orgUnitSvc:     orgUnitSvc,
		statusSvc:      statusSvc,
		userSettingSvc: userSettingSvc,
	}
}

func (s *ConsultationRequestService) GetAll(params *dtos.ConsultationRequestListQueryParamsDTO, currentUserID uint) (*dtos.ConsultationRequestListResponseDTO, error) {
	user, err := s.userSvc.GetFullDetailsByID(currentUserID)
	if err != nil {
		return nil, err
	}

	var departmentID *uint
	if user.OrganizationUnit2 != nil {
		id := uint(user.OrganizationUnit2.ID)
		departmentID = &id
	}

	if user.Roles != nil {
		for _, role := range *user.Roles {
			if *role.Type == "super-admin" {
				departmentID = nil
				break
			}
		}
	}
	return s.repo.GetAll(params, departmentID, currentUserID)
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
	status, err := s.statusSvc.GetByType("new")
	if err != nil {
		return 0, err
	}
	tx := s.repo.BeginTransaction()

	request := &models.ConsultationRequest{
		Title:                       dto.Title,
		ReferenceNumber:             referenceNumber,
		DateRaised:                  time.Now(),
		RequestingDepartmentUnitID:  responsibleUnitID,
		ConsultedDepartmentUnitID:   dto.ConsultedDepartmentID,
		ConsultationRequestStatusID: &status.ID,
		RequestDetails:              dto.RequestDetails,
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

	existingRequest, err := s.repo.GetByID(consultationRequestID)
	if err != nil {
		return 0, utility.NewResourceNotFoundError(err, consultationRequestID)
	}

	reqUsers, err := s.userSvc.GetUsersByOrganizationUnit(existingRequest.RequestingDepartmentUnitID)
	if err != nil {
		return 0, err
	}

	consUsers, err := s.userSvc.GetUsersByOrganizationUnit(existingRequest.ConsultedDepartmentUnitID)
	if err != nil {
		return 0, err
	}

	if dto.ConsultationRequestStatusID != nil {

		status, err := s.statusSvc.GetByID(*dto.ConsultationRequestStatusID)
		if err != nil {
			return 0, err
		}

		authorized := false

		if status.Type == "closed" || status.Type == "reopened" {
			for _, u := range reqUsers {
				if u.ID == currentUserID {
					authorized = true
					break
				}
			}
			if !authorized {
				return 0, utility.NewForbiddenError(
					"you do not have authority to update status",
				)
			}
		} else {
			for _, u := range consUsers {
				if u.ID == currentUserID {
					authorized = true
					break
				}
			}
			if !authorized {
				return 0, utility.NewForbiddenError(
					"you do not have authority to update status",
				)
			}
		}
		if status.Type == "reply_sent" {

			existingRequest.Reply = &dto.Comment
			now := time.Now()
			existingRequest.ReplyDate = &now

		}
	}

	tx := s.repo.BeginTransaction()
	now := time.Now()

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

	var documents []*models.ConsultationRequestCommentDocument

	for _,doc := range dto.Documents {
		if doc.FileName == "" || doc.Size <= 0 {
			continue
		}

		encryptedFileName, err := utility.EncryptData([]byte(doc.FileName))
		if err != nil {
			tx.Rollback()
			return 0, err
		}

		url := utility.GenerateDocumentURL(
			"consultation-request",
			"consultation-request",
			strconv.FormatUint(uint64(consultationRequestID), 10),
			"comments",
			strconv.FormatUint(uint64(commentID), 10),
			"documents",
		)

		document := &models.ConsultationRequestCommentDocument{
			ConsultationRequestCommentID: commentID,
			URL:                          url,
			Name:                         doc.FileName,
			Size:                         doc.Size,
			Token:                        encryptedFileName,
			OriginalName:                 doc.OriginalName,
			Type:                         doc.Type,
			CreatedByID:                  currentUserID,
			CreatedAt:                    now,
		}

		documents = append(documents,document)
	}

	if len(documents) > 0 {
		if err := s.repo.CreateDocuments(documents,tx); err != nil {
			tx.Rollback()
			return 0,err
		}
	}

	if dto.ConsultationRequestStatusID != nil {
		existingRequest.ConsultationRequestStatusID = dto.ConsultationRequestStatusID
	}
	existingRequest.UpdatedByID = &currentUserID
	existingRequest.UpdatedAt = &now

	if err := s.repo.Update(existingRequest, tx); err != nil {
		tx.Rollback()
		return 0, err
	}

	for _,doc := range documents {
		storagePath := filepath.Join(
			"consultation-request",
			"consultation-request",
			strconv.FormatUint(uint64(consultationRequestID), 10),
			"comments",
			strconv.FormatUint(uint64(commentID), 10),
			"documents",
		)

		if err := utility.StoreFile(doc.Name, storagePath); err != nil {
			tx.Rollback()
			return 0, err
		}
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

func (s *ConsultationRequestService) UpdateConsultationRequestComment(consultationRequestID uint, commentID uint, currentUserID uint, dto *dtos.ConsultationRequestCommentUpdateFormDTO) error {
	if _, err := s.repo.GetByID(consultationRequestID); err != nil {
		return utility.NewResourceNotFoundError(err, consultationRequestID)
	}

	existingComment, err := s.repo.GetCommentByID(commentID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utility.NewResourceNotFoundError(err, commentID)
		}
		return err
	}

	if existingComment.ConsultationRequestID != consultationRequestID {
		return utility.NewResourceNotFoundError(fmt.Errorf("consultation request comment not found"), commentID)
	}

	if existingComment.ConsultationRequestStatusID != nil {
		return utility.NewValidationError(
			fmt.Errorf("comments with status cannot be updated"),
			"invalid",
			"comment_id",
		)
	}

	tx := s.repo.BeginTransaction()
	now := time.Now()
	oldComment := existingComment.Comment

	existingComment.Comment = dto.Comment
	existingComment.UpdatedByID = &currentUserID
	existingComment.UpdatedAt = &now

	if err := s.repo.UpdateComment(existingComment, tx); err != nil {
		tx.Rollback()
		return err
	}

	var newDocuments []*models.ConsultationRequestCommentDocument
	var deleteDocumentIDs []*uint
	var documentsToDelete []*models.ConsultationRequestCommentDocument

	for _, doc := range dto.Documents {
		if doc.IsNewFile != nil && *doc.IsNewFile {
			if doc.FileName == "" || doc.Size <= 0 {
				continue
			}

			encryptedFileName, err := utility.EncryptData([]byte(doc.FileName))
			if err != nil {
				tx.Rollback()
				return err
			}

			url := utility.GenerateDocumentURL(
				"consultation-request",
				"consultation-request",
				strconv.FormatUint(uint64(consultationRequestID), 10),
				"comments",
				strconv.FormatUint(uint64(commentID), 10),
				"documents",
			)

			document := &models.ConsultationRequestCommentDocument{
				ConsultationRequestCommentID: commentID,
				URL:                          url,
				Name:                         doc.FileName,
				Size:                         doc.Size,
				Token:                        encryptedFileName,
				OriginalName:                 doc.OriginalName,
				Type:                         doc.Type,
				CreatedByID:                  currentUserID,
				CreatedAt:                    now,
			}

			newDocuments = append(newDocuments, document)
		} else if doc.IsDeletedFile != nil && *doc.IsDeletedFile && doc.ID != nil {
			existingDocument, err := s.repo.GetCommentDocumentByID(*doc.ID)
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					tx.Rollback()
					return utility.NewResourceNotFoundError(err, *doc.ID)
				}
				tx.Rollback()
				return err
			}

			if existingDocument.ConsultationRequestCommentID != commentID {
				tx.Rollback()
				return utility.NewResourceNotFoundError(fmt.Errorf("consultation request comment document not found"), *doc.ID)
			}

			documentsToDelete = append(documentsToDelete, existingDocument)
			deleteDocumentIDs = append(deleteDocumentIDs, &existingDocument.ID)
		}
	}

	if len(newDocuments) > 0 {
		if err := s.repo.CreateDocuments(newDocuments, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	if len(deleteDocumentIDs) > 0 {
		if err := s.repo.DeleteDocuments(deleteDocumentIDs, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, doc := range newDocuments {
		storagePath := filepath.Join(
			"consultation-request",
			"consultation-request",
			strconv.FormatUint(uint64(consultationRequestID), 10),
			"comments",
			strconv.FormatUint(uint64(commentID), 10),
			"documents",
		)

		if err := utility.StoreFile(doc.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, document := range documentsToDelete {
		storagePath := filepath.Join(
			"consultation-request",
			"consultation-request",
			strconv.FormatUint(uint64(consultationRequestID), 10),
			"comments",
			strconv.FormatUint(uint64(commentID), 10),
			"documents",
		)

		if err := utility.RemoveFile(document.Name, storagePath); err != nil {
			tx.Rollback()
			return err
		}
	}

	tx.Commit()

	go func() {
		changedData := utility.CompareFieldChanges(
			map[string]interface{}{"comment": oldComment},
			map[string]interface{}{"comment": existingComment.Comment},
		)

		if len(changedData) == 0 {
			return
		}

		err := s.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.ConsultationRequestComment{}.TableName(),
			commentID,
			changedData,
			currentUserID,
			"consultation_request",
			"consultation",
			"",
			nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.UpdateAction,
				"tableName": models.ConsultationRequestComment{}.TableName(),
				"recordID":  commentID,
				"userID":    currentUserID,
				"error":     err.Error(),
			}
			utility.LogError("Failed to log consultation request comment update activity", fields)
		}
	}()

	return nil
}

func (s *ConsultationRequestService) DeleteConsultationRequestComments(consultationRequestID uint, ids []uint, currentUserID uint) ([]uint, error) {
	tx := s.repo.BeginTransaction()

	if _, err := s.repo.GetByID(consultationRequestID); err != nil {
		return nil, utility.NewResourceNotFoundError(err, consultationRequestID)
	}

	var failedIDs []uint
	var successfulIDs []uint
	var deletedComments []map[string]interface{}

	for _, id := range ids {
		comment, err := s.repo.GetCommentByID(id)
		if err != nil {
			failedIDs = append(failedIDs, id)
			continue
		}

		if comment.ConsultationRequestID != consultationRequestID {
			failedIDs = append(failedIDs, id)
			continue
		}

		if comment.ConsultationRequestStatusID != nil {
			failedIDs = append(failedIDs, id)
			continue
		}

		if err := s.repo.DeleteComment(id, tx); err != nil {
			failedIDs = append(failedIDs, id)
			continue
		}

		successfulIDs = append(successfulIDs, id)
		deletedComments = append(deletedComments, map[string]interface{}{
			"id":                      id,
			"consultation_request_id": consultationRequestID,
		})
	}

	if len(failedIDs) > 0 {
		tx.Rollback()
		return failedIDs, nil
	}

	tx.Commit()

	if len(successfulIDs) > 0 {
		go func() {
			err := s.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.ConsultationRequestComment{}.TableName(),
				successfulIDs,
				deletedComments,
				currentUserID,
				nil,
			)
			if err != nil {
				fmt.Printf("Failed to log consultation request comment bulk delete activity: %v\n", err)
			}
		}()
	}

	return failedIDs, nil
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

func (s *ConsultationRequestService) ExportConsultationRequests(data []*dtos.ConsultationRequestListDTO, currentUserID uint) (*os.File, error) {
	// Fetch user setting for language preference
	userSetting, err := s.userSettingSvc.GetUserSettingsByUserID(currentUserID)
	if err != nil {
		return nil, err
	}
	var headers []string
	// Define headers based on the user's language setting
	if userSetting.Language == "ar" {
		headers = []string{
			"SL No.",
			"Reference Number",
			"Date Raised",
			"Reply Date",
			"Status",
			"Requesting Department",
			"Consulted Department",
			"Responsible User",
			"Request Details",
			"Notes",
			"Created By",
			"Created At",
		}
	} else {
		headers = []string{
			"SL No.",
			"Reference Number",
			"Date Raised",
			"Reply Date",
			"Status",
			"Requesting Department",
			"Consulted Department",
			"Responsible User",
			"Request Details",
			"Notes",
			"Created By",
			"Created At",
		}
	}
	// Convert the consultation request data to [][]interface{} based on language setting
	rows := make([][]interface{}, len(data))
	for i, request := range data {
		statusTitle := request.StatusTitleEn
		requestingDept := request.RequestingDepartmentUnitTitleEn
		consultedDept := request.ConsultedDepartmentUnitTitleEn
		responsibleUser := request.ResponsibleFirstNameEn + " " + request.ResponsibleLastNameEn
		createdBy := request.CreatedByFirstNameEn + " " + request.CreatedByLastNameEn

		if userSetting.Language == "ar" {
			statusTitle = request.StatusTitleAr
			requestingDept = request.RequestingDepartmentUnitTitleAr
			consultedDept = request.ConsultedDepartmentUnitTitleAr
			responsibleUser = request.ResponsibleFirstNameAr + " " + request.ResponsibleLastNameAr
			createdBy = request.CreatedByFirstNameAr + " " + request.CreatedByLastNameAr
		}

		rows[i] = []interface{}{
			i + 1,
			request.ReferenceNumber,
			request.DateRaised,
			request.ReplyDate,
			statusTitle,
			requestingDept,
			consultedDept,
			responsibleUser,
			request.RequestDetails,
			request.Notes,
			createdBy,
			request.CreatedAt,
		}
	}

	// Generate the Excel file
	return utility.GenerateExcelFile(headers, rows)
}

func (s *ConsultationRequestService) GetDocumentPath(consultationRequestID uint, commentID uint, token string) (string, error) {
	existingComment, err := s.repo.GetCommentByID(commentID)
	if err != nil {
		return "", err
	}

	fmt.Println("consultationRequestID", consultationRequestID)
	if existingComment.ConsultationRequestID != consultationRequestID {
		return "", fmt.Errorf("resource not found")
	}

	fmt.Println("existingComment.Documents", existingComment.Documents)
	if len(existingComment.Documents) == 0 {
		return "", fmt.Errorf("no documents found for this consultation request comment")
	}

	fmt.Println("token", token)
	decryptedFilenameBytes, err := utility.DecryptData(token)
	if err != nil {
		return "", fmt.Errorf("resource not found")
	}

	fmt.Println("decryptedFilenameBytes", decryptedFilenameBytes)
	decryptedFilename := string(decryptedFilenameBytes)
	if decryptedFilename == "" {
		return "", fmt.Errorf("resource not found")
	}

	fmt.Println("decryptedFilename", decryptedFilename)
	storagePath := filepath.Join(
		"consultation-request",
		"consultation-request",
		strconv.FormatUint(uint64(consultationRequestID), 10),
		"comments",
		strconv.FormatUint(uint64(commentID), 10),
		"documents",
	)

	fmt.Println("storagePath", storagePath)
	baseStoragePath := config.GetEnv("STORAGE_PATH", "")
	fullPath := filepath.Join(baseStoragePath, storagePath, decryptedFilename)

	fmt.Println("fullPath", fullPath)
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return "", fmt.Errorf("file not found")
	} else if err != nil {
		return "", fmt.Errorf("error accessing file: %w", err)
	}

	return fullPath, nil
}