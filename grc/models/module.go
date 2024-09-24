package models

type Module struct {
	ID         uint   `json:"id" gorm:"primaryKey;autoIncrement" binding:"-"`
	Identifier string `json:"identifier" gorm:"size:45;not null" binding:"required"`
	TitleEn    string `json:"title_en" gorm:"size:100;not null" binding:"required"`
	TitleAr    string `json:"title_ar" gorm:"size:100;not null" binding:"required"`
	ModuleID   *uint  `json:"module_id,omitempty" binding:"omitempty"`
	Order      int    `json:"order" gorm:"not null" binding:"required"`
	Icon       string `json:"icon" gorm:"size:45" binding:"required"`
	Route      string `json:"route" gorm:"size:255" binding:"required"`

	// Self-referencing relationship
	Parent     *Module   `json:"parent,omitempty" gorm:"foreignKey:ModuleID"`
	SubModules *[]Module `json:"sub_modules,omitempty" gorm:"foreignKey:ModuleID"`

	// One-to-many relationship
	NotificationMasters *[]NotificationMaster `json:"notification_masters,omitempty" gorm:"foreignKey:ModuleID"`
}

func (Module) TableName() string {
	return "modules"
}
