package services

import (
	"errors"
	"fmt"
	organizationDTOs "jethur/grc/organization/dtos"
	"jethur/utility"

	"jethur/grc/models"
	organizationRepositories "jethur/grc/organization/repositories"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type StrategicProfileServiceInterface interface {
	GetStrategicProfile() (*organizationDTOs.StrategicProfileDetailsDTO, error)
	Update(currentUserID uint, strategicProfile *organizationDTOs.StrategicProfileFormDTO) error
	UpdateVision(currentUserID uint, strategicProfileVision *organizationDTOs.StrategicProfileVisionFormDTO) error
	UpdateMission(currentUserID uint, strategicProfileMission *organizationDTOs.StrategicProfileMissionFormDTO) error
	UpdateValues(strategicProfileValues []*organizationDTOs.StrategicProfileValueFormDTO) error
}

type StrategicProfileService struct {
	repo organizationRepositories.StrategicProfileRepositoryInterface
}

func NewStrategicProfileService(repo organizationRepositories.StrategicProfileRepositoryInterface) StrategicProfileServiceInterface {

	return &StrategicProfileService{
		repo: repo,
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
	existingStrategicProfile.UpdatedBy = &currentUserID

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
	existingStrategicProfile.UpdatedBy = &currentUserID

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
	existingStrategicProfile.UpdatedBy = &currentUserID

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

func (us *StrategicProfileService) validateStrategicProfileFormDTO(strategicProfileDTO *organizationDTOs.StrategicProfileFormDTO) error {
	if strategicProfileDTO.Linkedin != nil && !utility.IsValidLinkedInURL(*strategicProfileDTO.Linkedin) {
		return utility.NewValidationError(fmt.Errorf("invalid LinkedIn URL"), "not_exist", "linkedin")
	}

	if strategicProfileDTO.Instagram != nil && !utility.IsValidInstagramURL(*strategicProfileDTO.Instagram) {
		return utility.NewValidationError(fmt.Errorf("invalid Instagram URL"), "not_exist", "instagram")
	}

	if strategicProfileDTO.Youtube != nil && !utility.IsValidYouTubeURL(*strategicProfileDTO.Youtube) {
		return utility.NewValidationError(fmt.Errorf("invalid YouTube URL"), "not_exist", "youtube")
	}

	if strategicProfileDTO.Facebook != nil && !utility.IsValidFacebookURL(*strategicProfileDTO.Facebook) {
		return utility.NewValidationError(fmt.Errorf("invalid Facebook URL"), "not_exist", "facebook")
	}

	if strategicProfileDTO.X != nil && !utility.IsValidTwitterURL(*strategicProfileDTO.X) {
		return utility.NewValidationError(fmt.Errorf("invalid X URL"), "not_exist", "twitter")
	}

	return nil

}
