package models

import (
	"time"
)

type CSLogSocRoomAccessLogDocument struct {
	ID                 uint      `json:"id" gorm:"primary_key;auto_increment"`
	SocRoomAccessLogID uint      `json:"soc_room_access_log_id" gorm:"not null"`
	URL                *string   `json:"url" gorm:"type:varchar(1000);default:null"`
	Name               *string   `json:"name" gorm:"type:nvarchar(255);default:null"`
	Size               *int      `json:"size" gorm:"default:null"`
	Token              *string   `json:"token" gorm:"type:varchar(255);default:null"`
	OriginalName       *string   `json:"original_name" gorm:"type:nvarchar(255);default:null"`
	Type               *string   `json:"type" gorm:"type:varchar(255);default:null"`
	CreatedBYID        *uint     `json:"created_by_id" gorm:"column:created_by_id;default:null"`
	UpdatedBYID        *uint     `json:"updated_by_id" gorm:"column:updated_by_id;default:null"`
	CreatedAt          time.Time `json:"created_at" gorm:"type:datetime;not null"`
	UpdatedAt          time.Time `json:"updated_at" gorm:"type:datetime;not null"`

	// Relationships
	CSLogSocRoomAccessLog *CSLogSocRoomAccessLog `gorm:"foreignKey:SocRoomAccessLogID;constraint:OnUpdate:NO ACTION,OnDelete:CASCADE"`
	CreatedByUser         *User                  `gorm:"foreignKey:CreatedBYID"`
	UpdatedByUser         *User                  `gorm:"foreignKey:UpdatedBYID"`
}

func (CSLogSocRoomAccessLogDocument) TableName() string {
	return "cs_log_soc_room_access_log_documents"
}
