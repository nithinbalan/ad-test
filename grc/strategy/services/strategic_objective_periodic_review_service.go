package services

import (
	"errors"
	"fmt"
	"jethur/config"
	"jethur/grc/logs/services"
	masterService "jethur/grc/masters/services"
	"jethur/grc/models"
	strategyDTOs "jethur/grc/strategy/dtos"
	strategyRepositories "jethur/grc/strategy/repositories"
	"jethur/utility"
	"os"
	"path/filepath"
	"strconv"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type StrategicObjectivePeriodicReviewServiceInterface interface {
	GetAll(params *strategyDTOs.StrategicObjectivePeriodicReviewQueryParamsDTO, documentId uint) (*strategyDTOs.StrategicObjectivePeriodicReviewListResponseDTO, error)
	GetByID(id uint) (*models.StrategicObjectiveReview, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*strategyDTOs.StrategicObjectivePeriodicReviewDetailsDTO, error)
	Create(currentUserID uint, strategicObjectiveID uint, formDTO *strategyDTOs.StrategicObjectivePeriodicReviewFormDTO) (uint, error)
	Update(id uint, strategicObjectiveID uint, currentUserID uint, formDTO *strategyDTOs.StrategicObjectivePeriodicReviewFormDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)

	GetDocumentPath(id uint, token string, tx ...utility.DatabaseTransaction) (string, error)
}

type StrategicObjectivePeriodicReviewService struct {
	repo                  strategyRepositories.StrategicObjectivePeriodicReviewRepositoryInterface
	activityLogSvc        services.ActivityLogServiceInterface
	strategicObjectiveSvc StrategicObjectiveServiceInterface
	reviewStatusSvc       masterService.ReviewStatusServiceInterface
}

func NewStrategicObjectivePeriodicReviewService(
	repo strategyRepositories.StrategicObjectivePeriodicReviewRepositoryInterface,
	activityLogSvc services.ActivityLogServiceInterface,
	strategicObjectiveSvc StrategicObjectiveServiceInterface,
	reviewStatusSvc masterService.ReviewStatusServiceInterface,
) StrategicObjectivePeriodicReviewServiceInterface {
	return &StrategicObjectivePeriodicReviewService{
		repo:                  repo,
		activityLogSvc:        activityLogSvc,
		strategicObjectiveSvc: strategicObjectiveSvc,
		reviewStatusSvc:       reviewStatusSvc,
	}
}

func (soprs *StrategicObjectivePeriodicReviewService) GetAll(params *strategyDTOs.StrategicObjectivePeriodicReviewQueryParamsDTO, documentId uint) (*strategyDTOs.StrategicObjectivePeriodicReviewListResponseDTO, error) {
	return soprs.repo.GetAll(params, documentId)
}

func (soprs *StrategicObjectivePeriodicReviewService) GetByID(id uint) (*models.StrategicObjectiveReview, error) {
	review, err := soprs.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return review, nil
}

func (soprs *StrategicObjectivePeriodicReviewService) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*strategyDTOs.StrategicObjectivePeriodicReviewDetailsDTO, error) {
	var detailsDTO strategyDTOs.StrategicObjectivePeriodicReviewDetailsDTO
	review, err := soprs.repo.GetFullDetailsByID(id, tx...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	// Check and format ReviewDate explicitly

	detailsDTO.ReviewDateString = review.ReviewDate.Format("2006-01-02")

	// Decode other fields
	err = mapstructure.Decode(review, &detailsDTO)
	if err != nil {
		return nil, err
	}

	detailsDTO.CreatedAt = review.CreatedAt
	detailsDTO.UpdatedAt = review.UpdatedAt

	return &detailsDTO, nil
}

func (soprs *StrategicObjectivePeriodicReviewService) Create(currentUserID uint, strategicObjectiveID uint, formDTO *strategyDTOs.StrategicObjectivePeriodicReviewFormDTO) (uint, error) {
	tx := soprs.repo.BeginTransaction()

	dates, err := soprs.parseControlDates(formDTO)
	if err != nil {
		return 0, err
	}

	// Update the document's review status to "reviewed"
	reviewStatus, err := soprs.reviewStatusSvc.GetByType("reviewed")
	if err != nil {
		return 0, err
	}

	strategicObjective, err := soprs.strategicObjectiveSvc.GetByID(strategicObjectiveID)
	if err != nil {
		return 0, err
	}

	review := &models.StrategicObjectiveReview{
		StrategicObjectiveID: strategicObjectiveID,
		ReviewDate:           *dates.ReviewDate,
		Description:          formDTO.Description,
		CreatedBYID:          currentUserID,
	}

	id, err := soprs.repo.Create(review, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	strategicObjective.ReviewStatusID = &reviewStatus.ID

	if err := soprs.strategicObjectiveSvc.UpdateStrategicObjective(strategicObjective, tx); err != nil {
		tx.Rollback()
		return 0, err
	}

	var documents []*models.StrategicObjectiveReviewDocument
	for _, doc := range formDTO.SupportDocuments {
		dataToEncrypt := []byte(doc.FileName)
		encryptedFileName, err := utility.EncryptData(dataToEncrypt)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
		url := utility.GenerateDocumentURL("strategy", "strategic-objective-reviews", strconv.FormatUint(uint64(id), 10), "documents", encryptedFileName)
		document := &models.StrategicObjectiveReviewDocument{
			StrategicObjectiveReviewID: id,
			URL:                        url,
			Name:                       doc.FileName,
			Size:                       doc.Size,
			Token:                      encryptedFileName,
			OriginalName:               &doc.OriginalName,
			Type:                       doc.Type,
		}
		documents = append(documents, document)
	}

	if len(documents) > 0 {
		if err := soprs.repo.CreateDocuments(documents, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	tx.Commit()

	for _, doc := range documents {
		storagePath := filepath.Join("strategy", "strategic-objective-reviews", strconv.FormatUint(uint64(id), 10), "documents")
		if err := utility.StoreFile(doc.Name, storagePath); err != nil {
			return 0, err
		}
	}

	return id, nil
}

func (soprs *StrategicObjectivePeriodicReviewService) Update(id uint, strategicObjectiveID uint, currentUserID uint, formDTO *strategyDTOs.StrategicObjectivePeriodicReviewFormDTO) error {
	tx := soprs.repo.BeginTransaction()

	existing, err := soprs.GetByID(id)
	if err != nil {
		return err
	}

	dates, err := soprs.parseControlDates(formDTO)
	if err != nil {
		return err
	}

	existing.ReviewDate = *dates.ReviewDate
	existing.Description = formDTO.Description
	// existing.UpdatedBYID = &currentUserID

	if err := soprs.repo.Update(existing, tx); err != nil {
		tx.Rollback()
		return err
	}

	var documents []*models.StrategicObjectiveReviewDocument
	for _, doc := range formDTO.SupportDocuments {
		dataToEncrypt := []byte(doc.FileName)
		encryptedFileName, err := utility.EncryptData(dataToEncrypt)
		if err != nil {
			tx.Rollback()
			return err
		}
		url := utility.GenerateDocumentURL("strategy", "strategic-objective-reviews", strconv.FormatUint(uint64(id), 10), "documents", encryptedFileName)
		document := &models.StrategicObjectiveReviewDocument{
			StrategicObjectiveReviewID: id,
			URL:                        url,
			Name:                       doc.FileName,
			Size:                       doc.Size,
			Token:                      encryptedFileName,
			OriginalName:               &doc.OriginalName,
			Type:                       doc.Type,
		}
		documents = append(documents, document)
	}

	if len(documents) > 0 {
		if err := soprs.repo.CreateDocuments(documents, tx); err != nil {
			tx.Rollback()
			return err
		}
	}
	tx.Commit()

	for _, doc := range documents {
		storagePath := filepath.Join("strategy", "strategic-objective-reviews", strconv.FormatUint(uint64(id), 10), "documents")
		if err := utility.StoreFile(doc.Name, storagePath); err != nil {
			fmt.Printf("Failed to store file %s: %v\n", doc.Name, err)
		}
	}

	return nil
}

func (soprs *StrategicObjectivePeriodicReviewService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	tx := soprs.repo.BeginTransaction()

	var failedIds []uint
	var successfulIds []uint

	for _, id := range ids {
		if err := soprs.repo.Delete(id, tx); err != nil {
			failedIds = append(failedIds, id)
		} else {
			successfulIds = append(successfulIds, id)
		}
	}

	if len(failedIds) > 0 {
		tx.Rollback()
		return failedIds, nil
	}

	tx.Commit()

	return failedIds, nil
}

func (soprs *StrategicObjectivePeriodicReviewService) GetDocumentPath(id uint, token string, tx ...utility.DatabaseTransaction) (string, error) {
	// Retrieve the full details of the document review
	existingStrategicObjectiveReview, err := soprs.GetFullDetailsByID(id, tx...)
	if err != nil {
		return "", err
	}

	// Decrypt the provided token to get the file name
	decryptedFilenameBytes, err := utility.DecryptData(token)
	if err != nil {
		return "", fmt.Errorf("resource not found")
	}

	decryptedFilename := string(decryptedFilenameBytes)
	if decryptedFilename == "" {
		return "", fmt.Errorf("resource not found")
	}

	// Construct the storage path
	storagePath := filepath.Join("strategy", "strategic-objective-reviews", strconv.FormatUint(uint64(existingStrategicObjectiveReview.ID), 10), "documents")

	// Retrieve the base storage path from environment variables
	baseStoragePath := config.GetEnv("STORAGE_PATH", "")
	fullPath := filepath.Join(baseStoragePath, storagePath, decryptedFilename)

	// Check if the file exists at the constructed path
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return "", fmt.Errorf("file not found")
	} else if err != nil {
		return "", fmt.Errorf("error accessing file: %w", err)
	}

	// Return the full path to the file
	return fullPath, nil
}

func (soprs *StrategicObjectivePeriodicReviewService) parseControlDates(dto *strategyDTOs.StrategicObjectivePeriodicReviewFormDTO) (*strategyDTOs.StrategicObjectivePeriodicReviewParsedDates, error) {
	parsed := &strategyDTOs.StrategicObjectivePeriodicReviewParsedDates{}

	var err error

	parsed.ReviewDate, err = utility.ParseDate(dto.ReviewDate)
	if err != nil {
		return nil, utility.NewValidationError(fmt.Errorf("invalid last review date"), "invalid_date", "last review date")
	}

	return parsed, nil
}
