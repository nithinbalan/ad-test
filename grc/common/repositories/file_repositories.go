package repositories

import (
	"io"
	"jethur/config"
	commonDtos "jethur/grc/common/dtos"
	"mime/multipart"
	"os"
	"path/filepath"
)

type FileRepositoryInterface interface {
	SaveFile(file *multipart.FileHeader, dst string) (*commonDtos.FileDTO, error)
}

type FileRepository struct{}

func NewFileRepository() *FileRepository {
	return &FileRepository{}
}

func (r *FileRepository) SaveFile(file *multipart.FileHeader, dst string) (*commonDtos.FileDTO, error) {
	storagePath := config.GetEnv("STORAGE_PATH", "")
	fullPath := filepath.Join(storagePath, dst)

	dir := filepath.Dir(fullPath)

	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return nil, err
	}

	// Open the source file
	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	// Create the destination file
	out, err := os.Create(fullPath)
	if err != nil {
		return nil, err
	}
	defer out.Close()

	// Copy the file content to the new file
	_, err = io.Copy(out, src)
	if err != nil {
		return nil, err
	}

	fileDTO := &commonDtos.FileDTO{
		FileName:     filepath.Base(fullPath),
		OriginalName: file.Filename,
		FileType:     filepath.Ext(fullPath),
		SizeInBytes:  file.Size,
	}

	return fileDTO, nil
}
