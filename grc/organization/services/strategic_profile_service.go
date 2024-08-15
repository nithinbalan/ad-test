package services

import (
	"errors"
	"fmt"
	"jethur/config"
	organizationDTOs "jethur/grc/organization/dtos"
	"jethur/utility"
	"mime/multipart"
	"os"
	"path/filepath"

	commonRepositories "jethur/grc/common/repositories"
	"jethur/grc/models"
	organizationRepositories "jethur/grc/organization/repositories"

	"github.com/google/uuid"
	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type StrategicProfileServiceInterface interface {
	GetStrategicProfile() (*organizationDTOs.StrategicProfileDetailsDTO, error)
	Update(currentUserID uint, strategicProfile *organizationDTOs.StrategicProfileFormDTO) error
	UpdateVision(currentUserID uint, strategicProfileVision *organizationDTOs.StrategicProfileVisionFormDTO) error
	UpdateMission(currentUserID uint, strategicProfileMission *organizationDTOs.StrategicProfileMissionFormDTO) error
	UpdateValues(strategicProfileValues []*organizationDTOs.StrategicProfileValueFormDTO) error
	UploadProfileImage(file *multipart.FileHeader) error
	GetProfileImagePath() (string, error)
}

type StrategicProfileService struct {
	repo     organizationRepositories.StrategicProfileRepositoryInterface
	fileRepo commonRepositories.FileRepositoryInterface
}

func NewStrategicProfileService(repo organizationRepositories.StrategicProfileRepositoryInterface,
	fileRepo commonRepositories.FileRepositoryInterface) StrategicProfileServiceInterface {

	return &StrategicProfileService{
		repo:     repo,
		fileRepo: fileRepo,
	}
}

// @todo - image support to be added
func (sps *StrategicProfileService) GetStrategicProfile() (*organizationDTOs.StrategicProfileDetailsDTO, error) {
	var strategicProfileDTO organizationDTOs.StrategicProfileDetailsDTO

	strategicProfile, err := sps.repo.GetStrategicProfileDetails()

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, nil)
		}
		return nil, err
	}

	err = mapstructure.Decode(strategicProfile, &strategicProfileDTO)
	if err != nil {
		return nil, err
	}

	for _, value := range strategicProfile.StrategicProfileValues {
		var valueDTO organizationDTOs.StrategicProfileValueFormDTO
		err := mapstructure.Decode(value, &valueDTO)
		if err != nil {
			return nil, err // Handle error appropriately
		}

		strategicProfileDTO.Values = append(strategicProfileDTO.Values, valueDTO)
	}

	//todo :- also add extra fields which required in details screen

	return &strategicProfileDTO, nil
}

func (sps *StrategicProfileService) Update(currentUserID uint, strategicProfileDTO *organizationDTOs.StrategicProfileFormDTO) error {
	existingStrategicProfile, err := sps.repo.Get()

	if err != nil {
		return utility.NewResourceNotFoundError(err, nil)
	}

	if err := sps.validateStrategicProfileFormDTO(strategicProfileDTO); err != nil {
		return err
	}

	existingStrategicProfile.TitleEn = strategicProfileDTO.TitleEn
	existingStrategicProfile.TitleAr = strategicProfileDTO.TitleAr
	existingStrategicProfile.Address = strategicProfileDTO.Address
	existingStrategicProfile.Location = strategicProfileDTO.Location
	existingStrategicProfile.Email = strategicProfileDTO.Email
	existingStrategicProfile.Industry = strategicProfileDTO.Industry
	existingStrategicProfile.Phone = strategicProfileDTO.Phone
	existingStrategicProfile.Linkedin = strategicProfileDTO.Linkedin
	existingStrategicProfile.X = strategicProfileDTO.X
	existingStrategicProfile.Facebook = strategicProfileDTO.Facebook
	existingStrategicProfile.Instagram = strategicProfileDTO.Instagram
	existingStrategicProfile.Youtube = strategicProfileDTO.Youtube
	existingStrategicProfile.AboutEn = strategicProfileDTO.AboutEn
	existingStrategicProfile.AboutAr = strategicProfileDTO.AboutAr
	existingStrategicProfile.UpdatedByID = &currentUserID

	err = sps.repo.Update(existingStrategicProfile)
	if err != nil {
		return err
	}

	return nil

}

func (sps *StrategicProfileService) UpdateVision(currentUserID uint, strategicProfileVisionDTO *organizationDTOs.StrategicProfileVisionFormDTO) error {
	existingStrategicProfile, err := sps.repo.Get()

	if err != nil {
		return utility.NewResourceNotFoundError(err, nil)
	}

	existingStrategicProfile.VisionEn = strategicProfileVisionDTO.VisionEn
	existingStrategicProfile.VisionAr = strategicProfileVisionDTO.VisionAr
	existingStrategicProfile.UpdatedByID = &currentUserID

	err = sps.repo.Update(existingStrategicProfile)
	if err != nil {
		return err
	}

	return nil

}

func (sps *StrategicProfileService) UpdateMission(currentUserID uint, strategicProfileMissionDTO *organizationDTOs.StrategicProfileMissionFormDTO) error {
	existingStrategicProfile, err := sps.repo.Get()

	if err != nil {
		return utility.NewResourceNotFoundError(err, nil)
	}

	existingStrategicProfile.MissionEn = strategicProfileMissionDTO.MissionEn
	existingStrategicProfile.MissionAr = strategicProfileMissionDTO.MissionAr
	existingStrategicProfile.UpdatedByID = &currentUserID

	err = sps.repo.Update(existingStrategicProfile)
	if err != nil {
		return err
	}

	return nil

}

func (sps *StrategicProfileService) UpdateValues(strategicProfileValuesDTO []*organizationDTOs.StrategicProfileValueFormDTO) error {
	existingStrategicProfile, err := sps.repo.Get()

	if err != nil {
		return utility.NewResourceNotFoundError(err, nil)
	}

	err = utility.CheckUniqueFields(strategicProfileValuesDTO, "TitleEn", "TitleAr")
	if err != nil {
		return utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title")
	}

	var values []*models.StrategicProfileValue
	for _, dto := range strategicProfileValuesDTO {

		value := &models.StrategicProfileValue{}

		if dto.ID != nil {
			value.ID = *dto.ID // Existing record
		}
		value.StrategicProfileID = existingStrategicProfile.ID
		value.TitleEn = dto.TitleEn
		value.TitleAr = dto.TitleAr
		value.DescriptionEn = dto.DescriptionEn
		value.DescriptionAr = dto.DescriptionAr

		values = append(values, value)

	}
	err = sps.repo.UpdateValues(values)
	if err != nil {
		return err
	}
	//@todo- activity log need to capture
	return nil
}

func (sps *StrategicProfileService) UploadProfileImage(file *multipart.FileHeader) error {
	allowedTypes := map[string]bool{"image/jpeg": true, "image/png": true}
	if !allowedTypes[file.Header.Get("Content-Type")] {
		return utility.NewValidationError(fmt.Errorf("only jpeg/png files are allowed"), "invalid_image_type", "type")
	}

	maxSize := int64(2 * 1024 * 1024)
	if file.Size > maxSize {
		return utility.NewValidationError(fmt.Errorf("only 2MB files are allowed"), "invalid_profile_file_size", "size")
	}

	fileExtension := filepath.Ext(file.Filename)
	newFileName := uuid.New().String() + fileExtension
	storagePath := "organization/profile_image/" + newFileName

	_, err := sps.fileRepo.SaveFile(file, storagePath)
	if err != nil {
		return err
	}

	profile, err := sps.repo.Get()
	if err != nil {
		return utility.NewResourceNotFoundError(err, nil)
	}

	if profile.ProfileImage != nil {
		oldFilePath := "organization/profile_image/" + *profile.ProfileImage
		err := os.Remove(filepath.Join(config.GetEnv("STORAGE_PATH", ""), oldFilePath))
		if err != nil && !os.IsNotExist(err) {
			return err
		}
	}

	// Generate token for the new file
	dataToEncrypt := []byte(newFileName)
	encryptedFilename, err := utility.EncryptData(dataToEncrypt)
	if err != nil {
		return err
	}

	// Update the profile with the new image
	profile.ProfileImage = &newFileName
	profile.ProfileImageToken = &encryptedFilename
	if err := sps.repo.Update(profile); err != nil {
		return err
	}

	return nil
}

func (sps *StrategicProfileService) GetProfileImagePath() (string, error) {
	profile, err := sps.repo.Get()
	if err != nil {
		return "", utility.NewResourceNotFoundError(err, nil)
	}

	if profile.ProfileImage == nil {
		return "", utility.NewResourceNotFoundError(fmt.Errorf("profile image not found"), nil)
	}

	storagePath := config.GetEnv("STORAGE_PATH", "")
	filePath := filepath.Join(storagePath, "organization/profile_image", *profile.ProfileImage)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return "", utility.NewResourceNotFoundError(fmt.Errorf("file not found"), nil)
	}

	return filePath, nil
}

func (us *StrategicProfileService) validateStrategicProfileFormDTO(strategicProfileDTO *organizationDTOs.StrategicProfileFormDTO) error {
	validateURL := func(url *string, validateFunc func(string) bool, fieldName string) error {
		if url != nil && *url != "" && !validateFunc(*url) {
			return utility.NewValidationError(fmt.Errorf("invalid %s URL", fieldName), "not_exist", fieldName)
		}
		return nil
	}

	if err := validateURL(strategicProfileDTO.Linkedin, utility.IsValidLinkedInURL, "linkedin"); err != nil {
		return err
	}
	if err := validateURL(strategicProfileDTO.Instagram, utility.IsValidInstagramURL, "instagram"); err != nil {
		return err
	}
	if err := validateURL(strategicProfileDTO.Youtube, utility.IsValidYouTubeURL, "youtube"); err != nil {
		return err
	}
	if err := validateURL(strategicProfileDTO.Facebook, utility.IsValidFacebookURL, "facebook"); err != nil {
		return err
	}
	if err := validateURL(strategicProfileDTO.X, utility.IsValidTwitterURL, "twitter"); err != nil { // Assuming X represents Twitter here
		return err
	}

	return nil
}
