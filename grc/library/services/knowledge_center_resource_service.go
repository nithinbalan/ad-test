package services

import (
	libraryDTOs "jethur/grc/library/dtos"
	"jethur/grc/library/repositories"
)

type KnowledgeCenterResourceServiceInterface interface {
	GetAll() ([]*libraryDTOs.KnowledgeCenterResourceListDTO, error)
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

func (kcrs *KnowledgeCenterResourceService) GetAll() ([]*libraryDTOs.KnowledgeCenterResourceListDTO, error) {
	return kcrs.repo.GetAll()
}
