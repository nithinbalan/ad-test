package models

import "time"

type RiskClosureDocument struct {
	ID            uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	RiskClosureID uint       `json:"risk_closure_id" binding:"required"`
	URL           string     `json:"url" binding:"required" gorm:"type:nvarchar(1000)"`
	Name          string     `json:"name" binding:"required" gorm:"type:nvarchar(255)"`
	Size          int        `json:"size" binding:"required"`
	Token         string     `json:"token" binding:"required" gorm:"type:varchar(255)"`
	OriginalName  string     `json:"original_name" gorm:"type:nvarchar(255)"`
	Type          string     `json:"type" binding:"required" gorm:"type:varchar(255)"`
	CreatedAt     time.Time  `json:"created_at" binding:"required" gorm:"type:datetime;not null"`
	CreatedBYID   uint       `json:"created_by_id" binding:"required"`
	UpdatedAt     *time.Time `json:"updated_at,omitempty" gorm:"type:datetime;default:null"`
	UpdatedBYID   *uint      `json:"updated_by_id"`

	// Relationships
	RiskClosure   *RiskClosure `gorm:"foreignKey:RiskClosureID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	CreatedByUser *User        `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	UpdatedByUser *User        `gorm:"foreignKey:UpdatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

// TableName overrides the table name used by RiskClosureDocument to `risk_closure_documents`
func (RiskClosureDocument) TableName() string {
	return "risk_closure_documents"
}
