package services

import (
	"fmt"
	"jethur/config"
	commonDtos "jethur/grc/common/dtos"
	"jethur/grc/common/repositories"
	"jethur/utility"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type FileServiceInterface interface {
	TempFileUpload(file *multipart.FileHeader) (*commonDtos.FileDTO, error)
	GetFilePreview(filePreviewDTO *commonDtos.FilePreviewDTO) (string, error)
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

func (service *FileService) GetFilePreview(filePreviewDTO *commonDtos.FilePreviewDTO) (string, error) {
	decryptedFilenameBytes, err := utility.DecryptData(filePreviewDTO.Token)
	if err != nil {
		return "", fmt.Errorf("invalid token: %w", err)
	}

	decryptedFilename := string(decryptedFilenameBytes)

	if len(decryptedFilename) == 0 {
		return "", fmt.Errorf("invalid decrypted filename")
	}

	baseStoragePath := config.GetEnv("STORAGE_PATH", "")
	fullPath := filepath.Join(baseStoragePath, filePreviewDTO.URL, decryptedFilename)

	// Check if the file exists
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return "", fmt.Errorf("file not found")
	} else if err != nil {
		return "", fmt.Errorf("error accessing file: %w", err)
	}

	return fullPath, nil
}
