package services

import (
	"crypto/rand"
	"fmt"
	logModels "jethur/grc/logs/models"
	logServices "jethur/grc/logs/services"
	userDTOs "jethur/grc/organization/dtos"
	organizationModels "jethur/grc/organization/models"
	"jethur/grc/organization/repositories"
	"jethur/utility"
	"log"

	"golang.org/x/crypto/bcrypt"
)

type UserServiceInterface interface {
	GetAll(params *userDTOs.UserListQueryParamsDTO) (*userDTOs.UserListResponseDTO, error)
	// GetByID(id uint) (*models.ComplianceRegister, error)
	Create(currentUserID uint, userFormDTO *userDTOs.UserFormDTO) (uint, error)
	Update(id uint, currentUserID uint, user *userDTOs.UserFormDTO) error
	// Delete(ids []uint, currentUserID uint) ([]uint, error)
}

type UserService struct {
	repo           repositories.UserRepositoryInterface
	designationSvc DesignationServiceInterface
	userStatusSvc  UserStatusServiceInterface
	activityLogSvc logServices.ActivityLogServiceInterface
}

func NewUserService(
	repo repositories.UserRepositoryInterface,
	designationSvc DesignationServiceInterface,
	userStatusSvc UserStatusServiceInterface,
	activityLogSvc logServices.ActivityLogServiceInterface) UserServiceInterface {
	return &UserService{
		repo:           repo,
		designationSvc: designationSvc,
		userStatusSvc:  userStatusSvc,
		activityLogSvc: activityLogSvc,
	}
}

func (us *UserService) GetAll(params *userDTOs.UserListQueryParamsDTO) (*userDTOs.UserListResponseDTO, error) {
	return us.repo.GetAll(params)
}

func (us *UserService) Create(currentUserID uint, userFormDTO *userDTOs.UserFormDTO) (uint, error) {

	if err := us.validateUserFormDTO(userFormDTO); err != nil {
		return 0, err
	}

	// Generate a random password
	password, err := us.generateRandomPassword(12) // Adjust length as needed
	if err != nil {
		return 0, err
	}

	// Hash the generated password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, fmt.Errorf("failed to hash password: %w", err)
	}

	userStatus, err := us.userStatusSvc.GetByStatus("active")
	if err != nil {
		log.Println("Error parsing next review date:", err)
		return 0, err
	}

	// Direct mapping from UserFormDTO to User model
	user := &organizationModels.User{
		OrganizationUnit1ID: userFormDTO.OrganizationUnit1ID,
		OrganizationUnit2ID: userFormDTO.OrganizationUnit2ID,
		OrganizationUnit3ID: userFormDTO.OrganizationUnit3ID,
		OrganizationUnit4ID: userFormDTO.OrganizationUnit4ID,
		OrganizationUnit5ID: userFormDTO.OrganizationUnit5ID,
		DesignationID:       userFormDTO.DesignationID,
		UserStatusID:        userStatus.ID,
		FirstNameEn:         userFormDTO.FirstNameEn,
		LastNameEn:          userFormDTO.LastNameEn,
		FirstNameAr:         userFormDTO.FirstNameAr,
		LastNameAr:          userFormDTO.LastNameAr,
		Email:               userFormDTO.Email,
		MobileNumber:        userFormDTO.MobileNumber,
		ProfileImagePath:    userFormDTO.ProfileImagePath,
		PasswordHash:        string(passwordHash), // Set the hashed password here
		CreatedBy:           currentUserID,
		// Initialize other fields as necessary
	}
	id, err := us.repo.Create(user)
	if err != nil {
		return 0, err
	}

	go func() {
		err := us.activityLogSvc.LogActivity(
			logModels.CreateAction,
			organizationModels.User{}.TableName(),
			id,
			userFormDTO,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    logModels.CreateAction,
				"tableName": organizationModels.User{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return id, nil
}

func (us *UserService) Update(id uint, currentUserID uint, userFormDTO *userDTOs.UserFormDTO) error {
	fmt.Println("hiiiiii")
	existingUser, err := us.repo.GetByID(id)

	if err != nil {
		return err
	}

	if err := us.validateUserFormDTO(userFormDTO); err != nil {
		return err
	}

	// Update the existingUser fields
	existingUser.OrganizationUnit1ID = userFormDTO.OrganizationUnit1ID
	existingUser.OrganizationUnit2ID = userFormDTO.OrganizationUnit2ID
	existingUser.OrganizationUnit3ID = userFormDTO.OrganizationUnit3ID
	existingUser.OrganizationUnit4ID = userFormDTO.OrganizationUnit4ID
	existingUser.OrganizationUnit5ID = userFormDTO.OrganizationUnit5ID
	existingUser.DesignationID = userFormDTO.DesignationID
	existingUser.FirstNameEn = userFormDTO.FirstNameEn
	existingUser.LastNameEn = userFormDTO.LastNameEn
	existingUser.FirstNameAr = userFormDTO.FirstNameAr
	existingUser.LastNameAr = userFormDTO.LastNameAr
	existingUser.Email = userFormDTO.Email
	existingUser.MobileNumber = userFormDTO.MobileNumber
	existingUser.ProfileImagePath = userFormDTO.ProfileImagePath
	existingUser.UpdatedBy = currentUserID

	// Make sure to call the update method on the correct repository
	err = us.repo.Update(existingUser)
	if err != nil {
		return err
	}

	return nil
}

func (us *UserService) validateUserFormDTO(userFormDTO *userDTOs.UserFormDTO) error {
	if userFormDTO.DesignationID != 0 {
		_, err := us.designationSvc.GetByID(userFormDTO.DesignationID)
		if err != nil {
			return utility.NewValidationError(fmt.Errorf("designation does not exist"), "not_exist", "designation")
		}
	}

	return nil
}

// generateRandomPassword creates a secure, random password of a specified length.
func (us *UserService) generateRandomPassword(length int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*"
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return "", fmt.Errorf("failed to generate random password: %w", err)
	}

	for i, v := range b {
		b[i] = charset[v%byte(len(charset))]
	}

	return string(b), nil
}
