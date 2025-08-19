package services

import (
	libraryDTOs "jethur/grc/library/dtos"
	"jethur/grc/library/repositories"
)

type ArticleLibraryServiceInterface interface {
	GetAll(params *libraryDTOs.ArticleLibraryListQueryParamsDTO) (*libraryDTOs.ArticleLibraryListResponseDTO, error)
}

// Service and implementations
type ArticleLibraryService struct {
	repo repositories.ArticleLibraryRepositoryInterface
}

func NewArticleLibraryService(
	repo repositories.ArticleLibraryRepositoryInterface,
) ArticleLibraryServiceInterface {
	return &ArticleLibraryService{
		repo: repo,
	}
}

func (vls *ArticleLibraryService) GetAll(params *libraryDTOs.ArticleLibraryListQueryParamsDTO) (*libraryDTOs.ArticleLibraryListResponseDTO, error) {
	return vls.repo.GetAll(params)
}
