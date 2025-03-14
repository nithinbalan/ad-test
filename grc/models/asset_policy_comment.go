package models

import "time"

// AssetPolicyComment represents comments on asset policies
type AssetPolicyComment struct {
	ID            uint       `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	AssetPolicyID uint       `json:"asset_policy_id" gorm:"not null"`
	Comments      string     `json:"comments" gorm:"type:text;not null"`
	CreatedBYID   uint       `json:"created_by_id" gorm:"not null"`
	CreatedAt     time.Time  `json:"created_at" binding:"required" gorm:"type:datetime;not null"`
	UpdatedAt     *time.Time `json:"updated_at,omitempty" gorm:"type:datetime"`
	UpdatedBYID   *uint      `json:"updated_by_id,omitempty"`

	// Relationships
	AssetPolicy   *AssetPolicy `gorm:"foreignKey:AssetPolicyID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	CreatedByUser *User        `gorm:"foreignKey:CreatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
	UpdatedByUser *User        `gorm:"foreignKey:UpdatedBYID;constraint:OnUpdate:NO ACTION,OnDelete:NO ACTION"`
}

// TableName overrides the default table name
func (AssetPolicyComment) TableName() string {
	return "asset_policy_comments"
}
