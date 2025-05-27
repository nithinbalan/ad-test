package repositories

import (
	libraryDTOs "jethur/grc/library/dtos"

	"gorm.io/gorm"
)

type KnowledgeCenterResourceRepositoryInterface interface {
	GetAll(params *libraryDTOs.KnowledgeCenterResourceListQueryParamsDTO) ([]*libraryDTOs.KnowledgeCenterResourceListDTO, error)
}

type KnowledgeCenterResourceRepository struct {
	db *gorm.DB
}

func NewKnowledgeCenterResourceRepository(db *gorm.DB) *KnowledgeCenterResourceRepository {
	return &KnowledgeCenterResourceRepository{
		db: db,
	}
}

func (repo *KnowledgeCenterResourceRepository) GetAll(params *libraryDTOs.KnowledgeCenterResourceListQueryParamsDTO) ([]*libraryDTOs.KnowledgeCenterResourceListDTO, error) {
	var resources []*libraryDTOs.KnowledgeCenterResourceListDTO

	query := repo.db.Table("knowledge_center_resources").
		Select(`
			knowledge_center_resources.id,
			knowledge_center_resources.title_en,
			knowledge_center_resources.title_ar,
			knowledge_center_resources.icon,
			knowledge_center_resources.item_count,
			knowledge_center_resources.type,
			knowledge_center_resources.created_at,
			knowledge_center_resources.updated_at,
			modules.identifier AS module_identifier,
			modules.title_en AS module_title_en,
			modules.title_ar AS module_title_ar,
			knowledge_center_resource_statuses.title_en AS knowledge_center_resource_status_title_en,
			knowledge_center_resource_statuses.title_ar AS knowledge_center_resource_status_title_ar,
			knowledge_center_resource_statuses.color_code AS color_code
		`).
		Joins("JOIN modules ON modules.id = knowledge_center_resources.module_id").
		Joins("JOIN knowledge_center_resource_statuses ON knowledge_center_resource_statuses.id = knowledge_center_resources.knowledge_center_resource_status_id")

	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where(`
			knowledge_center_resources.title_en LIKE ? OR
			knowledge_center_resources.title_ar LIKE ? OR
			modules.title_en LIKE ? OR
			modules.title_ar LIKE ?
		`, keyword, keyword, keyword, keyword)
	}

	err := query.Find(&resources).Error
	if err != nil {
		return nil, err
	}

	return resources, nil
}
