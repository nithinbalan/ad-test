package services

import (
	"fmt"
	documentDTOs "jethur/grc/documents/dtos"
	"jethur/grc/documents/repositories"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	organizationDTOs "jethur/grc/organization/dtos"
	settingDTOs "jethur/grc/settings/dtos"
	"jethur/utility"
)

type DocumentAccessServiceInterface interface {
	GetSharedUsers(documentId uint, page uint, pageSize uint, tx ...utility.DatabaseTransaction) (*documentDTOs.DocumentSharedUserListResponseDTO, error)
	GetReviewUsers(documentId uint, tx ...utility.DatabaseTransaction) ([]*organizationDTOs.UserMinimalDTO, error)
	GetSharedUnits(documentId uint, tx ...utility.DatabaseTransaction) ([]*settingDTOs.OrganizationUnitDTO, error)

	CreateSharedUsers(documentID uint, userIDs []uint, createdBy uint) error
	CreateSharedUnits(documentID uint, unitIDs []uint, createdBy uint) error
	DeleteSharedUsers(documentID uint, userIDs []*uint, tx ...utility.DatabaseTransaction) error
	DeleteAllSharedUsers(documentID uint, userIDs []*uint, tx ...utility.DatabaseTransaction) error
	DeleteSharedUnits(documentID uint, unitIDs []*uint, tx ...utility.DatabaseTransaction) error

	DeleteReviewUsers(documentID uint, userIDs []*uint, tx ...utility.DatabaseTransaction) error
	DeleteResponsibleUser(documentID uint, userID *uint, tx ...utility.DatabaseTransaction) error

	CreateUsers(sharedUsers []*models.DocumentSharedUser, tx ...utility.DatabaseTransaction) error
	CreateUnits(sharedUnits []*models.DocumentSharedUnit, tx ...utility.DatabaseTransaction) error
}

type DocumentAccessService struct {
	repo        repositories.DocumentAccessRepositoryInterface
	logService  logServices.ActivityLogServiceInterface
	documentSvc DocumentServiceInterface
}

func NewDocumentAccessService(
	repo repositories.DocumentAccessRepositoryInterface,
	logService logServices.ActivityLogServiceInterface,
	documentSvc DocumentServiceInterface) DocumentAccessServiceInterface {
	return &DocumentAccessService{
		repo:        repo,
		logService:  logService,
		documentSvc: documentSvc,
	}
}

func (das *DocumentAccessService) GetSharedUsers(documentId uint, page uint, pageSize uint, tx ...utility.DatabaseTransaction) (*documentDTOs.DocumentSharedUserListResponseDTO, error) {
	users, pagination, err := das.repo.GetSharedUsers(documentId, page, pageSize, tx...)
	if err != nil {
		return nil, err
	}

	dtos := make([]*organizationDTOs.UserMinimalDTO, 0) // Initialize as empty slice
	for _, user := range users {
		dto := &organizationDTOs.UserMinimalDTO{
			ID:               int(user.ID),
			FirstNameEN:      user.FirstNameEn,
			LastNameEN:       user.LastNameEn,
			FirstNameAR:      user.FirstNameAr,
			LastNameAR:       user.LastNameAr,
			ProfileImagePath: user.ProfileImagePath,
			Email:            user.Email,
			Initial:          user.Initial,
			Color:            user.Color,
		}
		dtos = append(dtos, dto)
	}

	responseDTO := &documentDTOs.DocumentSharedUserListResponseDTO{
		Items:      dtos,
		Pagination: pagination,
	}

	return responseDTO, nil
}

func (das *DocumentAccessService) GetReviewUsers(documentId uint, tx ...utility.DatabaseTransaction) ([]*organizationDTOs.UserMinimalDTO, error) {
	users, err := das.repo.GetReviewUsers(documentId, tx...)
	if err != nil {
		return nil, err
	}

	dtos := make([]*organizationDTOs.UserMinimalDTO, 0) // Initialize as empty slice, not nil
	for _, user := range users {
		dto := &organizationDTOs.UserMinimalDTO{
			ID:               int(user.ID),
			FirstNameEN:      user.FirstNameEn,
			LastNameEN:       user.LastNameEn,
			FirstNameAR:      user.FirstNameAr,
			LastNameAR:       user.LastNameAr,
			ProfileImagePath: user.ProfileImagePath,
			Email:            user.Email,
			Initial:          user.Initial,
			Color:            user.Color,
		}
		dtos = append(dtos, dto)
	}
	return dtos, nil
}

func (das *DocumentAccessService) GetSharedUnits(documentId uint, tx ...utility.DatabaseTransaction) ([]*settingDTOs.OrganizationUnitDTO, error) {
	units, err := das.repo.GetSharedUnits(documentId, tx...)
	if err != nil {
		return nil, err
	}

	dtos := make([]*settingDTOs.OrganizationUnitDTO, 0)
	for _, unit := range units {
		dto := &settingDTOs.OrganizationUnitDTO{
			ID:      int(unit.ID),
			TitleEn: unit.TitleEn,
			TitleAr: unit.TitleAr,
		}
		dtos = append(dtos, dto)
	}
	return dtos, nil
}

func (das *DocumentAccessService) CreateSharedUsers(documentID uint, userIDs []uint, createdBy uint) error {
	if _, err := das.documentSvc.GetByID(documentID); err != nil {
		return err
	}

	existingUserIDs, err := das.repo.GetExistingSharedUserIDs(documentID, userIDs)
	if err != nil {
		return fmt.Errorf("failed to fetch existing shared users: %w", err)
	}
	existingIDsMap := make(map[uint]bool)
	for _, id := range existingUserIDs {
		existingIDsMap[id] = true
	}

	// Filter out user IDs that already exist
	var newSharedUsers []*models.DocumentSharedUser
	for _, userID := range userIDs {
		if !existingIDsMap[userID] {
			newSharedUsers = append(newSharedUsers, &models.DocumentSharedUser{
				DocumentID:  documentID,
				UserID:      userID,
				Type:        "shared-user",
				CreatedBYID: createdBy,
			})
		}
	}

	if len(newSharedUsers) == 0 {
		return utility.NewForbiddenError("all provided user IDs already exist as shared users for this document.")
	}

	return das.CreateUsers(newSharedUsers)
}

func (das *DocumentAccessService) CreateSharedUnits(documentID uint, unitIDs []uint, createdBy uint) error {
	if _, err := das.documentSvc.GetByID(documentID); err != nil {
		return err
	}

	existingUnitIDs, err := das.repo.GetExistingSharedUnitIDs(documentID, unitIDs)
	if err != nil {
		return fmt.Errorf("failed to fetch existing shared units: %w", err)
	}
	existingIDsMap := make(map[uint]bool)
	for _, id := range existingUnitIDs {
		existingIDsMap[id] = true
	}

	// Filter out unit IDs that already exist
	var newSharedUnits []*models.DocumentSharedUnit
	for _, unitID := range unitIDs {
		if !existingIDsMap[unitID] {
			newSharedUnits = append(newSharedUnits, &models.DocumentSharedUnit{
				DocumentID:         documentID,
				OrganizationUnitID: unitID,
				CreatedBYID:        createdBy,
			})
		}
	}

	if len(newSharedUnits) == 0 {
		return utility.NewForbiddenError("all provided unit IDs already exist as shared units for this document.")
	}

	return das.CreateUnits(newSharedUnits)
}

func (das *DocumentAccessService) DeleteSharedUsers(documentID uint, userIDs []*uint, tx ...utility.DatabaseTransaction) error {
	if len(userIDs) == 0 {
		return nil
	}
	return das.repo.DeleteSharedUsers(documentID, userIDs, tx...)
}

func (das *DocumentAccessService) DeleteAllSharedUsers(documentID uint, userIDs []*uint, tx ...utility.DatabaseTransaction) error {
	if len(userIDs) == 0 {
		return nil
	}
	return das.repo.DeleteAllUsers(documentID, userIDs, tx...)
}

func (das *DocumentAccessService) DeleteSharedUnits(documentID uint, unitIDs []*uint, tx ...utility.DatabaseTransaction) error {
	if len(unitIDs) == 0 {
		return nil
	}
	return das.repo.DeleteSharedUnits(documentID, unitIDs, tx...)
}

func (das *DocumentAccessService) CreateUsers(sharedUsers []*models.DocumentSharedUser, tx ...utility.DatabaseTransaction) error {
	if len(sharedUsers) > 0 {
		if err := das.repo.CreateSharedUsers(sharedUsers, tx...); err != nil {
			return err
		}
	}
	return nil
}

func (das *DocumentAccessService) CreateUnits(sharedUnits []*models.DocumentSharedUnit, tx ...utility.DatabaseTransaction) error {
	if len(sharedUnits) > 0 {
		if err := das.repo.CreateSharedUnits(sharedUnits, tx...); err != nil {
			return err
		}
	}
	return nil
}

func (das *DocumentAccessService) DeleteReviewUsers(documentID uint, userIDs []*uint, tx ...utility.DatabaseTransaction) error {
	if len(userIDs) == 0 {
		return nil
	}
	return das.repo.DeleteReviewUsers(documentID, userIDs, tx...)
}

func (das *DocumentAccessService) DeleteResponsibleUser(documentID uint, userID *uint, tx ...utility.DatabaseTransaction) error {
	return das.repo.DeleteResponsibleUser(documentID, userID, tx...)
}
