package services

import (
	"crypto/rand"
	"fmt"
	logModels "jethur/grc/logs/models"
	logServices "jethur/grc/logs/services"
	userDTOs "jethur/grc/organization/dtos"
	organizationModels "jethur/grc/organization/models"
	organizationRepositories "jethur/grc/organization/repositories"
	"jethur/utility"
	"reflect"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type UserServiceInterface interface {
	GetAll(params *userDTOs.UserListQueryParamsDTO) (*userDTOs.UserListResponseDTO, error)
	GetByID(id uint) (*organizationModels.User, error)
	Create(currentUserID uint, userFormDTO *userDTOs.UserFormDTO) (uint, error)
	Update(id uint, currentUserID uint, user *userDTOs.UserFormDTO) error
	UpdatePartial(id uint, dto *userDTOs.UserUpdateDTO) error
	GetByEmail(email string) (*organizationModels.User, error)
	GetUserSecurityByUserID(userID uint) (*organizationModels.UserSecurity, error)
	UpdateUserSecurity(userSecurity *organizationModels.UserSecurity) error
	// Delete(ids []uint, currentUserID uint) ([]uint, error)
}

type UserService struct {
	repo           organizationRepositories.UserRepositoryInterface
	designationSvc DesignationServiceInterface
	userStatusSvc  UserStatusServiceInterface
	activityLogSvc logServices.ActivityLogServiceInterface
}

func NewUserService(
	repo organizationRepositories.UserRepositoryInterface,
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

	password, err := us.generateRandomPassword(12)
	if err != nil {
		return 0, err
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}

	userStatus, err := us.userStatusSvc.GetByType("active")
	if err != nil {
		return 0, err
	}

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
		CreatedBy:           currentUserID,
	}

	userID, err := us.repo.Create(user)
	if err != nil {
		return 0, err // Only return the error
	}

	userSecurity := &organizationModels.UserSecurity{
		UserID:       userID,
		PasswordHash: string(passwordHash),
	}

	if err = us.repo.CreateUserSecurity(userSecurity); err != nil {
		return userID, err
	}

	userSetting := &organizationModels.UserSetting{
		UserID:     userID,
		TimezoneID: 5,
		Language:   "en",
		Timeout:    5,
	}
	err = us.repo.CreateUserSetting(userSetting)
	if err != nil {
		return userID, err
	}

	go func() {
		err := us.activityLogSvc.LogActivity(
			logModels.CreateAction,
			organizationModels.User{}.TableName(),
			userID, // use userID instead of id
			userFormDTO,
			currentUserID,
			nil,
		)

		if err != nil {
			fields := map[string]interface{}{
				"action":    logModels.CreateAction,
				"tableName": organizationModels.User{}.TableName(),
				"recordID":  userID, // use userID instead of id
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()

	return userID, nil // Use userID, which was modified within the transaction block
}

func (us *UserService) Update(id uint, currentUserID uint, userFormDTO *userDTOs.UserFormDTO) error {

	existingUser, err := us.repo.GetByID(id)

	if err != nil {
		return err
	}

	if err := us.validateUserFormDTO(userFormDTO); err != nil {
		return err
	}

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
	existingUser.UpdatedBy = &currentUserID

	err = us.repo.Update(existingUser)
	if err != nil {
		return err
	}

	return nil
}

func (us *UserService) UpdatePartial(id uint, userUpdateDTO *userDTOs.UserUpdateDTO) error {

	updates := make(map[string]interface{})

	v := reflect.ValueOf(userUpdateDTO).Elem()
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if field.IsValid() && !field.IsNil() {
			jsonTag := t.Field(i).Tag.Get("json")
			tagKey := strings.Split(jsonTag, ",")[0]
			updates[tagKey] = field.Elem().Interface()
		}
	}
	err := us.repo.UpdatePartial(id, updates)
	if err != nil {
		return err
	}

	return nil
}

func (us *UserService) GetByID(id uint) (*organizationModels.User, error) {
	return us.repo.GetByID(id)
}

func (us *UserService) validateUserFormDTO(userFormDTO *userDTOs.UserFormDTO) error {
	_, err := us.designationSvc.GetByID(userFormDTO.DesignationID)
	if err != nil {
		return utility.NewValidationError(fmt.Errorf("designation does not exist"), "not_exist", "designation")
	}
	return nil
}

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

func (us *UserService) GetByEmail(email string) (*organizationModels.User, error) {
	return us.repo.GetByEmail(email)
}

func (us *UserService) GetUserSecurityByUserID(userId uint) (*organizationModels.UserSecurity, error) {
	return us.repo.GetUserSecurityByUserID(userId)
}

func (us *UserService) UpdateUserSecurity(userSecurity *organizationModels.UserSecurity) error {
	return us.repo.UpdateUserSecurity(userSecurity)
}
