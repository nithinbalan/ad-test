package services

import (
	libraryDTOs "jethur/grc/library/dtos"
	"jethur/grc/library/repositories"
	"jethur/grc/models"
)

type ArticleLibraryServiceInterface interface {
	GetAll(params *libraryDTOs.ArticleLibraryListQueryParamsDTO) (*libraryDTOs.ArticleLibraryListResponseDTO, error)
	GetFullDetailsByID(id uint) (*models.ArticleLibrary, error)
	GetArticleCategories() ([]libraryDTOs.ArticleCategory, error)
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

func (als *ArticleLibraryService) GetAll(params *libraryDTOs.ArticleLibraryListQueryParamsDTO) (*libraryDTOs.ArticleLibraryListResponseDTO, error) {
	return als.repo.GetAll(params)
}

func (als ArticleLibraryService) GetFullDetailsByID(id uint) (*models.ArticleLibrary, error) {
	return als.repo.GetFullDetailsByID(id)
}

func (als *ArticleLibraryService) GetArticleCategories() ([]libraryDTOs.ArticleCategory, error) {
	return als.repo.GetArticleCategories()
}
