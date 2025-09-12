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

	SaveGeneratedPDF(content []byte, fileName string) (*commonDtos.FileDTO, error)
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

func (service *FileService) SaveGeneratedPDF(content []byte, fileName string) (*commonDtos.FileDTO, error) {
	dst := "tmp/" + strconv.FormatInt(time.Now().UnixNano(), 10) + "_" + utility.SanitizeFileName(fileName)
	storagePath := config.GetEnv("STORAGE_PATH", "")
	fullPath := filepath.Join(storagePath, dst)

	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return nil, err
	}

	// Write the generated content directly
	if err := os.WriteFile(fullPath, content, 0644); err != nil {
		return nil, err
	}

	fileDTO := &commonDtos.FileDTO{
		FileName:     filepath.Base(fullPath),
		OriginalName: fileName,
		FileType:     filepath.Ext(fullPath),
		SizeInBytes:  int64(len(content)),
	}

	return fileDTO, nil
}
