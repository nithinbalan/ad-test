package services

import (
	"jethur/grc/logs/dtos"
	"jethur/grc/logs/repositories"
)

type APIRequestLogServiceInterface interface {
	GetAll(params *dtos.APIRequestLogListQueryParamsDTO) (*dtos.APIRequestLogListResponseDTO, error)
	GetSummary(params *dtos.APIRequestLogListQueryParamsDTO) (*dtos.APIRequestLogSummaryDTO, error)
}

type APIRequestLogService struct {
	repo repositories.APIRequestLogRepositoryInterface
}

func NewAPIRequestLogService(repo repositories.APIRequestLogRepositoryInterface) APIRequestLogServiceInterface {
	return &APIRequestLogService{repo: repo}
}

func (arls *APIRequestLogService) GetAll(params *dtos.APIRequestLogListQueryParamsDTO) (*dtos.APIRequestLogListResponseDTO, error) {
	return arls.repo.GetAll(params)
}

func (arls *APIRequestLogService) GetSummary(params *dtos.APIRequestLogListQueryParamsDTO) (*dtos.APIRequestLogSummaryDTO, error) {
	return arls.repo.GetSummary(params)
}