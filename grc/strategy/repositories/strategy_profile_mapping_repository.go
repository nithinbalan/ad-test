package repositories

import (
	"jethur/grc/models"
	"jethur/utility"

	"gorm.io/gorm"
)

type StrategyProfileMappingRepositoryInterface interface {
	BeginTransaction() utility.DatabaseTransaction

	GetStrategyProfileWithRelations(profileId uint, enabled map[string]bool) (*models.StrategyProfile, error)
	// FetchFocusAreasByProfile(profileID uint) ([]strategyDTOs.StrategyHierarchyNode, error)
	// FetchChildrenByProfile(level string, parentId int, profileID uint) ([]strategyDTOs.StrategyHierarchyNode, error)
}

type StrategyProfileMappingRepository struct {
	db *gorm.DB
}

func NewStrategyProfileMappingRepository(db *gorm.DB) *StrategyProfileMappingRepository {
	return &StrategyProfileMappingRepository{
		db: db,
	}
}

// BeginTransaction starts a new transaction.
func (repo *StrategyProfileMappingRepository) BeginTransaction() utility.DatabaseTransaction {
	return utility.NewGormDatabaseTransaction(repo.db)
}
func (r *StrategyProfileMappingRepository) GetStrategyProfileWithRelations(profileID uint, enabled map[string]bool) (*models.StrategyProfile, error) {
	db := r.db

	// Define preload paths in hierarchy order
	hierarchy := []string{
		"focus-area",
		"goal",
		"objective",
		"initiative",
		"project",
		"milestone",
		"action-plan",
	}

	// Mapping from flat keys to model struct fields
	keyToField := map[string]string{
		"focus-area":  "StrategyProfileFocusAreas",
		"goal":        "StrategyProfileGoals",
		"objective":   "StrategyProfileObjectives",
		"initiative":  "StrategyProfileInitiatives",
		"project":     "StrategyProfileProjects",
		"milestone":   "StrategyProfileMilestones",
		"action-plan": "StrategyProfileActionPlans",
	}

	// Build preload path step-by-step
	var path string
	for _, key := range hierarchy {
		if !enabled[key] {
			continue
		}
		if path == "" {
			path = keyToField[key]
		} else {
			path += "." + keyToField[key]
		}
		db = db.Preload(path)
	}

	var profile *models.StrategyProfile
	err := db.First(&profile, profileID).Error
	if err != nil {
		return nil, err
	}

	return profile, nil
}

// func (r *StrategyProfileMappingRepository) GetStrategyProfileWithRelations(profileID uint, settings map[string]bool) (*models.StrategyProfile, error) {
// 	db := r.db

// 	// Base preload condition
// 	if settings["focus-area"] {
// 		db = db.Preload("StrategyProfileFocusAreas")
// 		if settings["goal"] {
// 			db = db.Preload("StrategyProfileFocusAreas.StrategyProfileGoals")
// 			if settings["objective"] {
// 				db = db.Preload("StrategyProfileFocusAreas.StrategyProfileGoals.StrategyProfileObjectives")
// 				if settings["initiative"] {
// 					db = db.Preload("StrategyProfileFocusAreas.StrategyProfileGoals.StrategyProfileObjectives.StrategyProfileInitiatives")
// 					if settings["project"] {
// 						db = db.Preload("StrategyProfileFocusAreas.StrategyProfileGoals.StrategyProfileObjectives.StrategyProfileInitiatives.StrategyProfileProjects")
// 						if settings["milestone"] {
// 							db = db.Preload("StrategyProfileFocusAreas.StrategyProfileGoals.StrategyProfileObjectives.StrategyProfileInitiatives.StrategyProfileProjects.StrategyProfileMilestones")
// 							if settings["action-plan"] {
// 								db = db.Preload("StrategyProfileFocusAreas.StrategyProfileGoals.StrategyProfileObjectives.StrategyProfileInitiatives.StrategyProfileProjects.StrategyProfileMilestones.StrategyProfileActionPlans")
// 							}
// 						} else if settings["action-plan"] { //milestone disable
// 							db = db.Preload("StrategyProfileFocusAreas.StrategyProfileGoals.StrategyProfileObjectives.StrategyProfileInitiatives.StrategyProfileProjects.StrategyProfileActionPlans")
// 						}
// 					} else if settings["milestone"] { //project disable
// 						db = db.Preload("StrategyProfileFocusAreas.StrategyProfileGoals.StrategyProfileObjectives.StrategyProfileInitiatives.StrategyProfileMilestones")
// 						if settings["actipn-plan"] {
// 							db = db.Preload("StrategyProfileFocusAreas.StrategyProfileGoals.StrategyProfileObjectives.StrategyProfileInitiatives.StrategyProfileMilestones.StrategyProfileActionPlans")
// 						}
// 					} else if settings["action-plan"] { //project milestone disable
// 						db = db.Preload("StrategyProfileFocusAreas.StrategyProfileGoals.StrategyProfileObjectives.StrategyProfileInitiatives.StrategyProfileActionPlans")
// 					}
// 				} else if settings["project"] {
// 					db = db.Preload("StrategyProfileFocusAreas.StrategyProfileGoals.StrategyProfileObjectives.StrategyProfileProjects")
// 					if settings["milestone"] {
// 						db = db.Preload("StrategyProfileFocusAreas.StrategyProfileGoals.StrategyProfileObjectives.StrategyProfileProjects.StrategyProfileMilestones")
// 						if settings["action-plan"] {
// 							db = db.Preload("StrategyProfileFocusAreas.StrategyProfileGoals.StrategyProfileObjectives.StrategyProfileProjects.StrategyProfileMilestones.StrategyProfileActionPlans")
// 						}
// 					} else if settings["action-plan"] { //milestone disable
// 						db = db.Preload("StrategyProfileFocusAreas.StrategyProfileGoals.StrategyProfileObjectives.StrategyProfileProjects.StrategyProfileActionPlans")
// 					}
// 				} else if settings["milestone"] {
// 					db = db.Preload("StrategyProfileFocusAreas.StrategyProfileGoals.StrategyProfileObjectives.StrategyProfileMilestones")
// 					if settings["action-plan"] {
// 						db = db.Preload("StrategyProfileFocusAreas.StrategyProfileGoals.StrategyProfileObjectives.StrategyProfileMilestones.StrategyProfileActionPlans")
// 					}
// 				} else if settings["action-plan"] { //milestone disable
// 					db = db.Preload("StrategyProfileFocusAreas.StrategyProfileGoals.StrategyProfileObjectives.StrategyProfileActionPlans")
// 				}
// 			} else if settings["initiative"] {
// 				db = db.Preload("StrategyProfileFocusAreas.StrategyProfileGoals.StrategyProfileInitiatives")
// 				if settings["project"] {
// 					db = db.Preload("StrategyProfileFocusAreas.StrategyProfileGoals.StrategyProfileInitiatives.StrategyProfileProjects")
// 					if settings["milestone"] {
// 						db = db.Preload("StrategyProfileFocusAreas.StrategyProfileGoals.StrategyProfileInitiatives.StrategyProfileProjects.StrategyProfileMilestones")
// 						if settings["action-plan"] {
// 							db = db.Preload("StrategyProfileFocusAreas.StrategyProfileGoals.StrategyProfileInitiatives.StrategyProfileProjects.StrategyProfileMilestones.StrategyProfileActionPlans")
// 						}
// 					} else if settings["action-plan"] { //milestone disable
// 						db = db.Preload("StrategyProfileFocusAreas.StrategyProfileGoals.StrategyProfileInitiatives.StrategyProfileProjects.StrategyProfileActionPlans")
// 					}
// 				} else if settings["milestone"] { //project disable
// 					db = db.Preload("StrategyProfileFocusAreas.StrategyProfileGoals.StrategyProfileInitiatives.StrategyProfileMilestones")
// 					if settings["action-plan"] {
// 						db = db.Preload("StrategyProfileFocusAreas.StrategyProfileGoals.StrategyProfileInitiatives.StrategyProfileMilestones.StrategyProfileActionPlans")
// 					}
// 				} else if settings["action-plan"] { //project milestone disable
// 					db = db.Preload("StrategyProfileFocusAreas.StrategyProfileGoals.StrategyProfileInitiatives.StrategyProfileActionPlans")
// 				}
// 				// Repeat nesting logic from here downward...
// 			} else if settings["project"] {
// 				db = db.Preload("StrategyProfileFocusAreas.StrategyProfileGoals.StrategyProfileProjects")
// 				if settings["milestone"] {
// 					db = db.Preload("StrategyProfileFocusAreas.StrategyProfileGoals.StrategyProfileProjects.StrategyProfileMilestones")
// 					if settings["action-plan"] {
// 						db = db.Preload("StrategyProfileFocusAreas.StrategyProfileGoals.StrategyProfileProjects.StrategyProfileMilestones.StrategyProfileActionPlans")
// 					}
// 				} else if settings["action-plan"] { //milestone disable
// 					db = db.Preload("StrategyProfileFocusAreas.StrategyProfileGoals.StrategyProfileProjects.StrategyProfileActionPlans")
// 				}
// 			} else if settings["milestone"] {
// 				db = db.Preload("StrategyProfileFocusAreas.StrategyProfileGoals.StrategyProfileMilestones")
// 				if settings["action-plan"] {
// 					db = db.Preload("StrategyProfileFocusAreas.StrategyProfileGoals.StrategyProfileMilestones.StrategyProfileActionPlans")
// 				}
// 			} else if settings["action-plan"] { //milestone disable
// 				db = db.Preload("StrategyProfileFocusAreas.StrategyProfileGoals.StrategyProfileActionPlans")
// 			}
// 		} else if settings["objective"] {
// 			db = db.Preload("StrategyProfileFocusAreas.StrategyProfileObjectives")
// 			if settings["initiative"] {
// 				db = db.Preload("StrategyProfileFocusAreas.StrategyProfileObjectives.StrategyProfileInitiatives")
// 				if settings["project"] {
// 					db = db.Preload("StrategyProfileFocusAreas.StrategyProfileObjectives.StrategyProfileInitiatives.StrategyProfileProjects")
// 					if settings["milestone"] {
// 						db = db.Preload("StrategyProfileFocusAreas.StrategyProfileObjectives.StrategyProfileInitiatives.StrategyProfileProjects.StrategyProfileMilestones")
// 						if settings["action-plan"] {
// 							db = db.Preload("StrategyProfileFocusAreas.StrategyProfileObjectives.StrategyProfileInitiatives.StrategyProfileProjects.StrategyProfileMilestones.StrategyProfileActionPlans")
// 						}
// 					} else if settings["action-plan"] { //milestone disable
// 						db = db.Preload("StrategyProfileFocusAreas.StrategyProfileObjectives.StrategyProfileInitiatives.StrategyProfileProjects.StrategyProfileActionPlans")
// 					}
// 				} else if settings["milestone"] { //project disable
// 					db = db.Preload("StrategyProfileFocusAreas.StrategyProfileObjectives.StrategyProfileInitiatives.StrategyProfileMilestones")
// 					if settings["action-plan"] {
// 						db = db.Preload("StrategyProfileFocusAreas.StrategyProfileObjectives.StrategyProfileInitiatives.StrategyProfileMilestones.StrategyProfileActionPlans")
// 					}
// 				} else if settings["action-plan"] { //project milestone disable
// 					db = db.Preload("StrategyProfileFocusAreas.StrategyProfileObjectives.StrategyProfileInitiatives.StrategyProfileActionPlans")
// 				}
// 			} else if settings["project"] {
// 				db = db.Preload("StrategyProfileFocusAreas.StrategyProfileObjectives.StrategyProfileProjects")
// 				if settings["milestone"] {
// 					db = db.Preload("StrategyProfileFocusAreas.StrategyProfileObjectives.StrategyProfileProjects.StrategyProfileMilestones")
// 					if settings["action-plan"] {
// 						db = db.Preload("StrategyProfileFocusAreas.StrategyProfileObjectives.StrategyProfileProjects.StrategyProfileMilestones.StrategyProfileActionPlans")
// 					}
// 				} else if settings["action-plan"] { //milestone disable
// 					db = db.Preload("StrategyProfileFocusAreas.StrategyProfileObjectives.StrategyProfileProjects.StrategyProfileActionPlans")
// 				}
// 			} else if settings["milestone"] {
// 				db = db.Preload("StrategyProfileFocusAreas.StrategyProfileObjectives.StrategyProfileMilestones")
// 				if settings["action-plan"] {
// 					db = db.Preload("StrategyProfileFocusAreas.StrategyProfileObjectives.StrategyProfileMilestones.StrategyProfileActionPlans")
// 				}
// 			} else if settings["action-plan"] { //milestone disable
// 				db = db.Preload("StrategyProfileFocusAreas.StrategyProfileObjectives.StrategyProfileActionPlans")
// 			}
// 		} else if settings["initiative"] {
// 			db = db.Preload("StrategyProfileFocusAreas.StrategyProfileInitiatives")
// 			if settings["project"] {
// 				db = db.Preload("StrategyProfileFocusAreas.StrategyProfileInitiatives.StrategyProfileProjects")
// 				if settings["milestone"] {
// 					db = db.Preload("StrategyProfileFocusAreas.StrategyProfileInitiatives.StrategyProfileProjects.StrategyProfileMilestones")
// 					if settings["action-plan"] {
// 						db = db.Preload("StrategyProfileFocusAreas.StrategyProfileInitiatives.StrategyProfileProjects.StrategyProfileMilestones.StrategyProfileActionPlans")
// 					}
// 				} else if settings["action-plan"] { //milestone disable
// 					db = db.Preload("StrategyProfileFocusAreas.StrategyProfileInitiatives.StrategyProfileProjects.StrategyProfileActionPlans")
// 				}
// 			} else if settings["milestone"] { //project disable
// 				db = db.Preload("StrategyProfileFocusAreas.StrategyProfileInitiatives.StrategyProfileMilestones")
// 				if settings["action-plan"] {
// 					db = db.Preload("StrategyProfileFocusAreas.StrategyProfileInitiatives.StrategyProfileMilestones.StrategyProfileActionPlans")
// 				}
// 			} else if settings["action-plan"] { //project milestone disable
// 				db = db.Preload("StrategyProfileFocusAreas.StrategyProfileInitiatives.StrategyProfileActionPlans")
// 			}
// 			// Repeat nesting logic from here downward...
// 		} else if settings["project"] {
// 			db = db.Preload("StrategyProfileFocusAreas.StrategyProfileProjects")
// 			if settings["milestone"] {
// 				db = db.Preload("StrategyProfileFocusAreas.StrategyProfileProjects.StrategyProfileMilestones")
// 				if settings["action-plan"] {
// 					db = db.Preload("StrategyProfileFocusAreas.StrategyProfileProjects.StrategyProfileMilestones.StrategyProfileActionPlans")
// 				}
// 			} else if settings["action-plan"] { //milestone disable
// 				db = db.Preload("StrategyProfileFocusAreas.StrategyProfileProjects.StrategyProfileActionPlans")
// 			}
// 		} else if settings["milestone"] {
// 			db = db.Preload("StrategyProfileFocusAreas.StrategyProfileMilestones")
// 			if settings["action-plan"] {
// 				db = db.Preload("StrategyProfileFocusAreas.StrategyProfileMilestones.StrategyProfileActionPlans")
// 			}
// 		} else if settings["action-plan"] { //milestone disable
// 			db = db.Preload("StrategyProfileFocusAreas.StrategyProfileActionPlans")
// 		}
// 	} else {
// 		// Focus Area is disabled
// 		if settings["goal"] {
// 			db = db.Preload("StrategyProfileGoals")
// 			if settings["objective"] {
// 				db = db.Preload("StrategyProfileGoals.StrategyProfileObjectives")
// 				if settings["initiative"] {
// 					db = db.Preload("StrategyProfileGoals.StrategyProfileObjectives.StrategyProfileInitiatives")
// 					if settings["project"] {
// 						db = db.Preload("StrategyProfileGoals.StrategyProfileObjectives.StrategyProfileInitiatives.StrategyProfileProjects")
// 						if settings["milestone"] {
// 							db = db.Preload("StrategyProfileGoals.StrategyProfileObjectives.StrategyProfileInitiatives.StrategyProfileProjects.StrategyProfileStrategyProfileMilestones")
// 							if settings["action-plan"] {
// 								db = db.Preload("StrategyProfileGoals.StrategyProfileObjectives.StrategyProfileInitiatives.StrategyProfileProjects.StrategyProfileMilestones.StrategyProfileActionPlans")
// 							}
// 						} else if settings["action-plan"] { //milestone disable
// 							db = db.Preload("StrategyProfileGoals.StrategyProfileObjectives.StrategyProfileInitiatives.StrategyProfileProjects.StrategyProfileActionPlans")
// 						}
// 					} else if settings["milestone"] { //project disable
// 						db = db.Preload("StrategyProfileGoals.StrategyProfileObjectives.StrategyProfileInitiatives.StrategyProfileMilestones")
// 						if settings["action-plan"] {
// 							db = db.Preload("StrategyProfileGoals.StrategyProfileObjectives.StrategyProfileInitiatives.StrategyProfileMilestones.StrategyProfileActionPlans")
// 						}
// 					} else if settings["action-plan"] { //project milestone disable
// 						db = db.Preload("StrategyProfileGoals.StrategyProfileObjectives.StrategyProfileInitiatives.StrategyProfileActionPlans")
// 					}
// 				} else if settings["project"] {
// 					db = db.Preload("StrategyProfileGoals.StrategyProfileObjectives.StrategyProfileProjects")
// 					if settings["milestone"] {
// 						db = db.Preload("StrategyProfileGoals.StrategyProfileObjectives.StrategyProfileProjects.StrategyProfileMilestones")
// 						if settings["action-plan"] {
// 							db = db.Preload("StrategyProfileGoals.StrategyProfileObjectives.StrategyProfileProjects.StrategyProfileMilestones.StrategyProfileActionPlans")
// 						}
// 					} else if settings["action-plan"] { //milestone disable
// 						db = db.Preload("StrategyProfileGoals.StrategyProfileObjectives.StrategyProfileProjects.StrategyProfileActionPlans")
// 					}
// 				} else if settings["milestone"] {
// 					db = db.Preload("StrategyProfileGoals.StrategyProfileObjectives.StrategyProfileMilestones")
// 					if settings["action-plan"] {
// 						db = db.Preload("StrategyProfileGoals.StrategyProfileObjectives.StrategyProfileMilestones.StrategyProfileActionPlans")
// 					}
// 				} else if settings["action-plan"] { //milestone disable
// 					db = db.Preload("StrategyProfileGoals.StrategyProfileObjectives.StrategyProfileActionPlans")
// 				}
// 			} else if settings["initiative"] {
// 				db = db.Preload("StrategyProfileGoals.StrategyProfileInitiatives")
// 				if settings["project"] {
// 					db = db.Preload("StrategyProfileGoals.StrategyProfileInitiatives.StrategyProfileProjects")
// 					if settings["milestone"] {
// 						db = db.Preload("StrategyProfileGoals.StrategyProfileInitiatives.StrategyProfileProjects.StrategyProfileMilestones")
// 						if settings["action-plan"] {
// 							db = db.Preload("StrategyProfileGoals.StrategyProfileInitiatives.StrategyProfileProjects.StrategyProfileMilestones.StrategyProfileActionPlans")
// 						}
// 					} else if settings["action-plan"] { //milestone disable
// 						db = db.Preload("StrategyProfileGoals.StrategyProfileInitiatives.StrategyProfileProjects.StrategyProfileActionPlans")
// 					}
// 				} else if settings["milestone"] { //project disable
// 					db = db.Preload("StrategyProfileGoals.StrategyProfileInitiatives.StrategyProfileMilestones")
// 					if settings["action-plan"] {
// 						db = db.Preload("StrategyProfileGoals.StrategyProfileInitiatives.StrategyProfileMilestones.StrategyProfileActionPlans")
// 					}
// 				} else if settings["action-plan"] { //project milestone disable
// 					db = db.Preload("StrategyProfileGoals.StrategyProfileInitiatives.StrategyProfileActionPlans")
// 				}
// 				// Repeat nesting logic from here downward...
// 			} else if settings["project"] {
// 				db = db.Preload("StrategyProfileGoals.StrategyProfileProjects")
// 				if settings["milestone"] {
// 					db = db.Preload("StrategyProfileGoals.StrategyProfileProjects.StrategyProfileMilestones")
// 					if settings["action-plan"] {
// 						db = db.Preload("StrategyProfileGoals.StrategyProfileProjects.StrategyProfileMilestones.StrategyProfileActionPlans")
// 					}
// 				} else if settings["action-plan"] { //milestone disable
// 					db = db.Preload("StrategyProfileGoals.StrategyProfileProjects.StrategyProfileActionPlans")
// 				}
// 			} else if settings["milestone"] {
// 				db = db.Preload("StrategyProfileGoals.StrategyProfileMilestones")
// 				if settings["action-plan"] {
// 					db = db.Preload("StrategyProfileGoals.StrategyProfileMilestones.StrategyProfileActionPlans")
// 				}
// 			} else if settings["action-plan"] { //milestone disable
// 				db = db.Preload("StrategyProfileGoals.StrategyProfileActionPlans")
// 			}
// 		} else if settings["objective"] {
// 			db = db.Preload("StrategyProfileObjectives")
// 			if settings["initiative"] {
// 				db = db.Preload("StrategyProfileObjectives.StrategyProfileInitiatives")
// 				if settings["project"] {
// 					db = db.Preload("StrategyProfileObjectives.StrategyProfileInitiatives.StrategyProfileProjects")
// 					if settings["milestone"] {
// 						db = db.Preload("StrategyProfileObjectives.StrategyProfileInitiatives.StrategyProfileProjects.StrategyProfileMilestones")
// 						if settings["action-plan"] {
// 							db = db.Preload("StrategyProfileObjectives.StrategyProfileInitiatives.StrategyProfileProjects.StrategyProfileMilestones.StrategyProfileActionPlans")
// 						}
// 					} else if settings["action-plan"] { //milestone disable
// 						db = db.Preload("StrategyProfileObjectives.StrategyProfileInitiatives.StrategyProfileProjects.StrategyProfileActionPlans")
// 					}
// 				} else if settings["milestone"] { //project disable
// 					db = db.Preload("StrategyProfileObjectives.StrategyProfileInitiatives.StrategyProfileMilestones")
// 					if settings["action-plan"] {
// 						db = db.Preload("StrategyProfileObjectives.StrategyProfileInitiatives.StrategyProfileMilestones.StrategyProfileActionPlans")
// 					}
// 				} else if settings["action-plan"] { //project milestone disable
// 					db = db.Preload("StrategyProfileObjectives.StrategyProfileInitiatives.StrategyProfileActionPlans")
// 				}
// 			} else if settings["project"] {
// 				db = db.Preload("StrategyProfileObjectives.StrategyProfileProjects")
// 				if settings["milestone"] {
// 					db = db.Preload("StrategyProfileObjectives.StrategyProfileProjects.StrategyProfileMilestones")
// 					if settings["action-plan"] {
// 						db = db.Preload("StrategyProfileObjectives.StrategyProfileProjects.StrategyProfileMilestones.StrategyProfileActionPlans")
// 					}
// 				} else if settings["action-plan"] { //milestone disable
// 					db = db.Preload("StrategyProfileObjectives.StrategyProfileProjects.StrategyProfileActionPlans")
// 				}
// 			} else if settings["milestone"] {
// 				db = db.Preload("StrategyProfileObjectives.StrategyProfileMilestones")
// 				if settings["action-plan"] {
// 					db = db.Preload("StrategyProfileObjectives.StrategyProfileMilestones.StrategyProfileActionPlans")
// 				}
// 			} else if settings["action-plan"] { //milestone disable
// 				db = db.Preload("StrategyProfileObjectives.StrategyProfileActionPlans")
// 			}
// 		} else if settings["initiative"] {
// 			db = db.Preload("StrategyProfileInitiatives")
// 			if settings["project"] {
// 				db = db.Preload("StrategyProfileInitiatives.StrategyProfileProjects")
// 				if settings["milestone"] {
// 					db = db.Preload("StrategyProfileInitiatives.StrategyProfileProjects.StrategyProfileMilestones")
// 					if settings["action-plan"] {
// 						db = db.Preload("StrategyProfileInitiatives.StrategyProfileProjects.StrategyProfileMilestones.StrategyProfileActionPlans")
// 					}
// 				} else if settings["action-plan"] { //milestone disable
// 					db = db.Preload("StrategyProfileInitiatives.StrategyProfileProjects.StrategyProfileActionPlans")
// 				}
// 			} else if settings["milestone"] { //project disable
// 				db = db.Preload("StrategyProfileInitiatives.StrategyProfileMilestones")
// 				if settings["action-plan"] {
// 					db = db.Preload("StrategyProfileInitiatives.StrategyProfileMilestones.StrategyProfileActionPlans")
// 				}
// 			} else if settings["action-plan"] { //project milestone disable
// 				db = db.Preload("StrategyProfileInitiatives.StrategyProfileActionPlans")
// 			}
// 			// Repeat nesting logic from here downward...
// 		} else if settings["project"] {
// 			db = db.Preload("StrategyProfileProjects")
// 			if settings["milestone"] {
// 				db = db.Preload("StrategyProfileProjects.StrategyProfileMilestones")
// 				if settings["action-plan"] {
// 					db = db.Preload("StrategyProfileProjects.StrategyProfileMilestones.StrategyProfileActionPlans")
// 				}
// 			} else if settings["action-plan"] { //milestone disable
// 				db = db.Preload("StrategyProfileProjects.StrategyProfileActionPlans")
// 			}
// 		} else if settings["milestone"] {
// 			db = db.Preload("StrategyProfileMilestones")
// 			if settings["action-plan"] {
// 				db = db.Preload("StrategyProfileMilestones.StrategyProfileActionPlans")
// 			}
// 		} else if settings["action-plan"] { //milestone disable
// 			db = db.Preload("StrategyProfileActionPlans")
// 		}
// 	}

// 	var profile *models.StrategyProfile
// 	err := db.Where("id = ?", profileID).First(&profile).Error
// 	if err != nil {
// 		return nil, err
// 	}

// 	return profile, nil
// }
