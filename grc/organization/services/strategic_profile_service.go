package services

import (
	"errors"
	"fmt"
	organizationDTOs "jethur/grc/organization/dtos"
	"jethur/utility"
	"net/url"
	"strings"

	organizationModels "jethur/grc/organization/models"
	organizationRepositories "jethur/grc/organization/repositories"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type StrategicProfileServiceInterface interface {
	GetStrategicProfile() (*organizationDTOs.StrategicProfileDetailsDTO, error)
	Update(id uint, currentUserID uint, strategicProfile *organizationDTOs.StrategicProfileFormDTO) error
	UpdateVision(id uint, currentUserID uint, strategicProfileVision *organizationDTOs.StrategicProfileVisionFormDTO) error
	UpdateMission(id uint, currentUserID uint, strategicProfileMission *organizationDTOs.StrategicProfileMisionFormDTO) error
	UpdateValue(id uint, strategicProfileValues []*organizationDTOs.StrategicProfileValueFormDTO) error
}

type StrategicProfileService struct {
	repo organizationRepositories.StrategicProfileRepositoryInterface
}

func NewStrategicProfileService(repo organizationRepositories.StrategicProfileRepositoryInterface) StrategicProfileServiceInterface {

	return &StrategicProfileService{
		repo: repo,
	}
}

func (sps *StrategicProfileService) GetStrategicProfile() (*organizationDTOs.StrategicProfileDetailsDTO, error) {
	// Initialize your DTO
	var strategicProfileDTO organizationDTOs.StrategicProfileDetailsDTO

	strategicProfile, err := sps.repo.GetStrategicProfile()

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, nil)
		}
		return nil, err
	}

	// Map the result to your DTO
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

func (sps *StrategicProfileService) Update(id uint, currentUserID uint, strategicProfileDTO *organizationDTOs.StrategicProfileFormDTO) error {

	existingStrategicProfile, err := sps.repo.GetByID(id)

	if err != nil {
		return err
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
	existingStrategicProfile.UpdatedBy = &currentUserID

	err = sps.repo.Update(existingStrategicProfile)
	if err != nil {
		return err
	}

	return nil

}

func (sps *StrategicProfileService) UpdateVision(id uint, currentUserID uint, strategicProfileVisionDTO *organizationDTOs.StrategicProfileVisionFormDTO) error {

	existingStrategicProfile, err := sps.repo.GetByID(id)

	if err != nil {
		return err
	}

	existingStrategicProfile.VisionEn = strategicProfileVisionDTO.VisionEn
	existingStrategicProfile.VisionAr = strategicProfileVisionDTO.VisionAr
	existingStrategicProfile.UpdatedBy = &currentUserID

	err = sps.repo.Update(existingStrategicProfile)
	if err != nil {
		return err
	}

	return nil

}

func (sps *StrategicProfileService) UpdateMission(id uint, currentUserID uint, strategicProfileMissionDTO *organizationDTOs.StrategicProfileMisionFormDTO) error {

	existingStrategicProfile, err := sps.repo.GetByID(id)

	if err != nil {
		return err
	}

	existingStrategicProfile.MissionEn = strategicProfileMissionDTO.MissionEn
	existingStrategicProfile.MissionAr = strategicProfileMissionDTO.MissionAr
	existingStrategicProfile.UpdatedBy = &currentUserID

	err = sps.repo.Update(existingStrategicProfile)
	if err != nil {
		return err
	}

	return nil

}

func (sps *StrategicProfileService) UpdateValue(strategicProfileID uint, strategicProfileValuesDTO []*organizationDTOs.StrategicProfileValueFormDTO) error {
	seenTitlesEn := make(map[string]bool)
	seenTitlesAr := make(map[string]bool)

	var values []*organizationModels.StrategicProfileValue
	for _, dto := range strategicProfileValuesDTO {

		if _, found := seenTitlesEn[dto.TitleEn]; found {
			return utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title_en")
		}
		seenTitlesEn[dto.TitleEn] = true

		if _, found := seenTitlesAr[dto.TitleAr]; found {
			return utility.NewValidationError(fmt.Errorf("the title must be unique"), "unique", "title_ar")
		}
		seenTitlesAr[dto.TitleAr] = true

		value := &organizationModels.StrategicProfileValue{}

		if dto.ID != nil {
			value.ID = *dto.ID // Existing record
		}
		value.StrategicProfileID = strategicProfileID
		value.TitleEn = dto.TitleEn
		value.TitleAr = dto.TitleAr
		value.DescriptionEn = dto.DescriptionEn
		value.DescriptionAr = dto.DescriptionAr

		values = append(values, value)

	}
	err := sps.repo.UpdateValue(strategicProfileID, values)
	if err != nil {
		return err
	}
	//@todo- activity log need to capture
	return nil
}

func (us *StrategicProfileService) validateStrategicProfileFormDTO(strategicProfileDTO *organizationDTOs.StrategicProfileFormDTO) error {
	if strategicProfileDTO.Linkedin != nil && !isValidLinkedInURL(*strategicProfileDTO.Linkedin) {
		return utility.NewValidationError(fmt.Errorf("invalid LinkedIn URL"), "not_exist", "linkedin")
	}
	if strategicProfileDTO.Instagram != nil && !isValidInstagramURL(*strategicProfileDTO.Instagram) {
		return utility.NewValidationError(fmt.Errorf("invalid Instagram URL"), "not_exist", "instagram")
	}
	if strategicProfileDTO.Youtube != nil && !isValidYouTubeURL(*strategicProfileDTO.Youtube) {
		return utility.NewValidationError(fmt.Errorf("invalid YouTube URL"), "not_exist", "youtube")
	}
	if strategicProfileDTO.Facebook != nil && !isValidFacebookURL(*strategicProfileDTO.Facebook) {
		return utility.NewValidationError(fmt.Errorf("invalid Facebook URL"), "not_exist", "facebook")
	}
	if strategicProfileDTO.X != nil && !isValidTwitterURL(*strategicProfileDTO.X) {
		return utility.NewValidationError(fmt.Errorf("invalid Twitter URL"), "not_exist", "twitter")
	}

	return nil

}

// Helper function to parse and validate a URL's scheme and host.
func isValidURL(toTest string) (bool, *url.URL) {
	u, err := url.Parse(toTest)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return false, nil
	}
	return true, u
}

// Specific validation functions for each platform
func isValidLinkedInURL(toTest string) bool {
	ok, u := isValidURL(toTest)
	if !ok {
		return false
	}
	host := strings.ToLower(u.Host)
	return host == "www.linkedin.com" || host == "linkedin.com" || strings.HasSuffix(host, ".linkedin.com")
}

func isValidInstagramURL(toTest string) bool {
	ok, u := isValidURL(toTest)
	if !ok {
		return false
	}
	host := strings.ToLower(u.Host)
	return host == "www.instagram.com" || host == "instagram.com"
}

func isValidYouTubeURL(toTest string) bool {
	ok, u := isValidURL(toTest)
	if !ok {
		return false
	}
	host := strings.ToLower(u.Host)
	return host == "www.youtube.com" || host == "youtube.com" || strings.HasPrefix(host, "youtu.be")
}

func isValidFacebookURL(toTest string) bool {
	ok, u := isValidURL(toTest)
	if !ok {
		return false
	}
	host := strings.ToLower(u.Host)
	return host == "www.facebook.com" || host == "facebook.com"
}

func isValidTwitterURL(toTest string) bool {
	ok, u := isValidURL(toTest)
	if !ok {
		return false
	}
	host := strings.ToLower(u.Host)
	return host == "www.twitter.com" || host == "twitter.com" || strings.HasPrefix(host, "m.twitter.com")
}
