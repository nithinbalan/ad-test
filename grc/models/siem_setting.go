package models

import (
	"time"
)

type SIEMSetting struct {
	ID                   uint       `json:"id" gorm:"primaryKey;autoIncrement"`
	IntegrationID        uint       `json:"integration_id"`
	SyslogServer         string     `json:"syslog_server" gorm:"size:255;not null"`
	SyslogPort           int        `json:"syslog_port" gorm:"not null"`
	SyslogProtocol       string     `json:"syslog_protocol" gorm:"size:10;not null"`
	SyslogTimezone       string     `json:"syslog_timezone,omitempty" gorm:"size:50"`
	SyslogDatetimeFormat string     `json:"syslog_datetime_format,omitempty" gorm:"size:50"`
	UpdatedByID          *uint      `json:"updated_by_id,omitempty" gorm:"default:null"`
	UpdatedAt            *time.Time `json:"updated_at,omitempty" gorm:"type:datetime;default:null"`

	// Relationships
	Integration   *Integration `gorm:"foreignKey:IntegrationID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	CreatedByUser *User        `gorm:"foreignKey:CreatedByID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	UpdatedByUser *User        `gorm:"foreignKey:UpdatedByID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

func (SIEMSetting) TableName() string {
	return "siem_settings"
}
