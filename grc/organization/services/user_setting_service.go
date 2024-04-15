package services

import (
	organizationRepositories "jethur/grc/organization/repositories"
)

type UserSettingServiceInterface interface {
}

type UserSettingService struct {
	repo organizationRepositories.UserSettingRepositoryInterface
}

func NewUserSettingService(repo organizationRepositories.UserSettingRepositoryInterface) UserSettingServiceInterface {
	return &UserSettingService{
		repo: repo,
	}
}
