package services

import (
	meetingDTOs "jethur/grc/meeting/dtos"
	meetingRepositories "jethur/grc/meeting/repositories"
)

type MeetingAssociatedItemServiceInterface interface {
	GetAllControls(meetingID uint, params *meetingDTOs.MeetingAssociatedItemListQueryParamsDTO) ([]*meetingDTOs.MeetingControlListDTO, error)
	GetAllComplianceRequirements(meetingID uint, params *meetingDTOs.MeetingAssociatedItemListQueryParamsDTO) ([]*meetingDTOs.MeetingComplianceRequirementListDTO, error)
	GetAllRisks(meetingID uint, params *meetingDTOs.MeetingAssociatedItemListQueryParamsDTO) ([]*meetingDTOs.MeetingRiskListDTO, error)
	GetAllProcesses(meetingID uint, params *meetingDTOs.MeetingAssociatedItemListQueryParamsDTO) ([]*meetingDTOs.MeetingProcessListDTO, error)
}

// Service and implementations
type MeetingAssociatedItemService struct {
	repo       meetingRepositories.MeetingAssociatedItemRepositoryInterface
	meetingSvc MeetingServiceInterface
}

func NewMeetingAssociatedItemService(
	repo meetingRepositories.MeetingAssociatedItemRepositoryInterface,
	meetingSvc MeetingServiceInterface,

) MeetingAssociatedItemServiceInterface {
	return &MeetingAssociatedItemService{
		repo:       repo,
		meetingSvc: meetingSvc,
	}
}

func (mais *MeetingAssociatedItemService) GetAllControls(meetingID uint, params *meetingDTOs.MeetingAssociatedItemListQueryParamsDTO) ([]*meetingDTOs.MeetingControlListDTO, error) {
	if _, err := mais.meetingSvc.GetByID(meetingID); err != nil {
		return nil, err
	}
	return mais.repo.GetAllControls(meetingID, params)
}

func (mais *MeetingAssociatedItemService) GetAllComplianceRequirements(meetingID uint, params *meetingDTOs.MeetingAssociatedItemListQueryParamsDTO) ([]*meetingDTOs.MeetingComplianceRequirementListDTO, error) {
	if _, err := mais.meetingSvc.GetByID(meetingID); err != nil {
		return nil, err
	}
	return mais.repo.GetAllComplianceRequirements(meetingID, params)
}

func (mais *MeetingAssociatedItemService) GetAllRisks(meetingID uint, params *meetingDTOs.MeetingAssociatedItemListQueryParamsDTO) ([]*meetingDTOs.MeetingRiskListDTO, error) {
	if _, err := mais.meetingSvc.GetByID(meetingID); err != nil {
		return nil, err
	}
	return mais.repo.GetAllRisks(meetingID, params)
}

func (mais *MeetingAssociatedItemService) GetAllProcesses(meetingID uint, params *meetingDTOs.MeetingAssociatedItemListQueryParamsDTO) ([]*meetingDTOs.MeetingProcessListDTO, error) {
	if _, err := mais.meetingSvc.GetByID(meetingID); err != nil {
		return nil, err
	}
	return mais.repo.GetAllProcesses(meetingID, params)
}
