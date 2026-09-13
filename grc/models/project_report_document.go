package models

import "time"

type ProjectReportDocument struct {
	ID                     uint       `json:"id" gorm:"primary_key;auto_increment"`
	ProjectReportProjectID uint       `json:"project_report_project_id" gorm:"not null"`
	Name                   *string    `json:"name" gorm:"type:nvarchar(255)"`
	Size                   *int       `json:"size" gorm:"type:int"`
	Token                  *string    `json:"token" gorm:"type:varchar(max)"`
	OriginalName           *string    `json:"original_name" gorm:"type:nvarchar(255)"`
	Type                   *string    `json:"type" gorm:"type:varchar(255)"`
	CreatedAt              time.Time  `json:"created_at" gorm:"type:datetime;not null;default:getdate()"`
	UpdatedAt              *time.Time `json:"updated_at" gorm:"type:datetime"`

	ProjectReportProject *ProjectReportProject `gorm:"foreignKey:ProjectReportProjectID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

func (ProjectReportDocument) TableName() string {
	return "project_report_documents"
}
