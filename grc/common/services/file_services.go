package services

import (
	commonDtos "jethur/grc/common/dtos"
	"jethur/grc/common/repositories"
	"jethur/utility"
	"mime/multipart"
	"strconv"
	"time"
)

type FileServiceInterface interface {
	TempFileUpload(file *multipart.FileHeader) (*commonDtos.FileDTO, error)
}

type FileService struct {
	repo repositories.FileRepositoryInterface
}

func NewFileService(repo repositories.FileRepositoryInterface) FileServiceInterface {
	return &FileService{
		repo: repo,
	}
}

func (service *FileService) TempFileUpload(file *multipart.FileHeader) (*commonDtos.FileDTO, error) {
	dst := "tmp/" + strconv.FormatInt(time.Now().UnixNano(), 10) + "_" + utility.SanitizeFileName(file.Filename)
	uploadedFile, err := service.repo.SaveFile(file, dst)
	if err != nil {
		return nil, err
	}
	return uploadedFile, nil
}
