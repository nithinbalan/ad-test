package initializers

import "gorm.io/gorm"

type LibraryModule struct {
}

func NewLibraryModule(db *gorm.DB, sharedModule *SharedModule, MasterModule *MasterModule) *LibraryModule {
	return &LibraryModule{}
}
