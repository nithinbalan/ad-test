package models

type Currency struct {
	ID    uint   `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	Title string `json:"title" binding:"required,max=45" gorm:"type:varchar(45);not null"`
	Code  string `json:"code" binding:"required,max=10" gorm:"type:varchar(10);not null"`
	Icon  string `json:"icon" binding:"required,max=45" gorm:"type:varchar(45);not null"`
}

// TableName overrides the table name used by Currency to `currencies`
func (Currency) TableName() string {
	return "currencies"
}
