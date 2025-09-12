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
	// Create a temp file to simulate upload
	tempFile, err := os.CreateTemp("", fileName)
	if err != nil {
		return nil, err
	}
	defer os.Remove(tempFile.Name())

	if _, err := tempFile.Write(content); err != nil {
		return nil, err
	}
	tempFile.Close()

	// Reopen as multipart.FileHeader
	f, err := os.Open(tempFile.Name())
	if err != nil {
		return nil, err
	}
	defer f.Close()

	fileStat, _ := f.Stat()
	header := &multipart.FileHeader{
		Filename: filepath.Base(fileName),
		Size:     fileStat.Size(),
	}

	// Now reuse TempFileUpload
	return service.TempFileUpload(header)
}
