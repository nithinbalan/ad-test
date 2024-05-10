package initializers

import "gorm.io/gorm"

type LibraryModule struct {
}

func NewLibraryModule(db *gorm.DB) *LibraryModule {
	return &LibraryModule{}
}
