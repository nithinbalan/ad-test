package dtos

import (
	commonDTOs "jethur/grc/common/dtos"
	organizationDTOs "jethur/grc/organization/dtos"
	"time"
)

type MeetingListDTO struct {
	ID                 uint       `json:"id"`
	MeetingID          *uint      `json:"meeting_id"`
	MeetingTypeID      uint       `json:"meeting_type_id"`
	MeetingFrequencyID *uint      `json:"meeting_frequency_id"`
	Title              string     `json:"title"`
	Location           *string    `json:"location"`
	Date               time.Time  `json:"date"`
	StartTime          time.Time  `json:"start_time"`
	EndTime            time.Time  `json:"end_time"`
	OrganizerID        uint       `json:"organizer_id"`
	CreatedBy          uint       `json:"created_by"`
	UpdatedBy          *uint      `json:"updated_by"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          *time.Time `json:"updated_at"`
	MeetingStatusID    uint       `json:"meeting_status_id"`
}

type MeetingListQueryParamsDTO struct {
	Sort        string
	SortBy      string
	Keyword     string
	Status      string
	Page        int
	PageSize    int
	Filters     map[string][]string
	QuickFilter []string
	DateFrom    string
	DateTo      string
}

type MeetingListResponseDTO struct {
	Items      []*MeetingListDTO            `json:"items"`
	Pagination commonDTOs.PaginationInfoDTO `json:"pagination"`
}

type MeetingDetailDTO struct {
	ID               uint                             `json:"id"`
	Meeting          *MeetingDetailDTO                `json:"meeting"`
	MeetingType      MeetingTypeDTO                   `json:"meeting_type"`
	MeetingFrequency *MeetingFrequencyDTO             `json:"meeting_frequency"`
	Title            string                           `json:"title"`
	Location         *string                          `json:"location"`
	Date             time.Time                        `json:"date"`
	StartTime        time.Time                        `json:"start_time"`
	EndTime          time.Time                        `json:"end_time"`
	Organizer        organizationDTOs.UserMinimalDTO  `json:"organizer"`
	CreatedByUser    organizationDTOs.UserMinimalDTO  `json:"created_by"`
	UpdatedByUser    *organizationDTOs.UserMinimalDTO `json:"updated_by"`
	CreatedAt        time.Time                        `json:"created_at"`
	UpdatedAt        *time.Time                       `json:"updated_at"`
	MeetingStatus    MeetingStatusDTO                 `json:"meeting_status"`
	Participants     []*MeetingParticipantDTO         `json:"participants"`
}

type MeetingParticipantDTO struct {
	ID        *uint                           `json:"id"`
	MeetingID uint                            `json:"meeting_id"`
	UserID    uint                            `json:"user_id"`
	User      organizationDTOs.UserMinimalDTO `json:"user"`
}

type MeetingFormDTO struct {
	Title              string  `json:"title" binding:"required"`
	MeetingTypeID      uint    `json:"meeting_type_id" binding:"required"`
	MeetingFrequencyID *uint   `json:"meeting_frequency_id"`
	MeetingID          *uint   `json:"meeting_id"`
	Location           *string `json:"location"`
	Date               string  `json:"date" binding:"required,datetime=2006-01-02"`
	StartTime          string  `json:"start_time" binding:"required"`
	EndTime            string  `json:"end_time" binding:"required"`
	OrganizerID        uint    `json:"organizer_id" binding:"required"`
	MeetingStatusID    uint    `json:"meeting_status_id" binding:"required"`
	UserIDs            []*uint `json:"user_ids"`
}

type MeetingParsedDates struct {
	Date      time.Time
	StartTime time.Time
	EndTime   time.Time
}
