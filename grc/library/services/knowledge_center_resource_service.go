package services

import (
	libraryDTOs "jethur/grc/library/dtos"
	"jethur/grc/library/repositories"
)

type KnowledgeCenterResourceServiceInterface interface {
	GetAll(params *libraryDTOs.KnowledgeCenterResourceListQueryParamsDTO) ([]*libraryDTOs.KnowledgeCenterResourceListDTO, error)
}

// Service and implementations
type KnowledgeCenterResourceService struct {
	repo repositories.KnowledgeCenterResourceRepositoryInterface
}

func NewKnowledgeCenterResourceService(
	repo repositories.KnowledgeCenterResourceRepositoryInterface,
) KnowledgeCenterResourceServiceInterface {
	return &KnowledgeCenterResourceService{
		repo: repo,
	}
}

func (kcrs *KnowledgeCenterResourceService) GetAll(params *libraryDTOs.KnowledgeCenterResourceListQueryParamsDTO) ([]*libraryDTOs.KnowledgeCenterResourceListDTO, error) {
	return kcrs.repo.GetAll(params)
}
