package models

type BCMBIAOutageBucket struct {
	ID             uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	TitleEn        string `json:"title_en" binding:"required" gorm:"type:nvarchar(500);not null"`
	TitleAr        string `json:"title_ar" binding:"required" gorm:"type:nvarchar(500);not null"`
	FromValue      int    `json:"from_value" binding:"required" gorm:"type:int;not null"`
	ToValue        int    `json:"to_value" binding:"required" gorm:"type:int;not null"`
	DurationUnitID uint   `json:"duration_unit_id" binding:"required" gorm:"not null"`
	Order          int    `json:"order" gorm:"type:int;not null"`
	Status         string `json:"status" binding:"required,oneof=active inactive" gorm:"type:varchar(10);not null;default:'active'"`

	// Relationships
	DurationUnit *BusinessContinuityDurationUnit `gorm:"foreignKey:DurationUnitID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

func (BCMBIAOutageBucket) TableName() string {
	return "bcm_bia_outage_buckets"
}
