package services

import (
	"errors"
	documentDTOs "jethur/grc/documents/dtos"
	"jethur/grc/documents/repositories"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type DocumentReviewServiceInterface interface {
	GetComments(id uint) (*documentDTOs.DocumentVersionWorkflowDTO, error)
	AddComment(id uint, currentUserID uint, documentReviewCommentFormDTO *documentDTOs.DocumentReviewCommentFormDTO) error
}

type DocumentReviewService struct {
	repo                      repositories.DocumentReviewRepositoryInterface
	logService                logServices.ActivityLogServiceInterface
	documentSvc               DocumentServiceInterface
	documentWorkflowSvc       DocumentWorkflowServiceInterface
	documentWorkflowStatusSvc DocumentWorkflowStatusServiceInterface
	documentStatusSvc         DocumentStatusServiceInterface
}

func NewDocumentReviewService(
	repo repositories.DocumentReviewRepositoryInterface,
	logService logServices.ActivityLogServiceInterface,
	documentSvc DocumentServiceInterface,
	documentWorkflowSvc DocumentWorkflowServiceInterface,
	documentWorkflowStatusSvc DocumentWorkflowStatusServiceInterface,
	documentStatusSvc DocumentStatusServiceInterface) DocumentReviewServiceInterface {
	return &DocumentReviewService{
		repo:                      repo,
		logService:                logService,
		documentSvc:               documentSvc,
		documentWorkflowSvc:       documentWorkflowSvc,
		documentWorkflowStatusSvc: documentWorkflowStatusSvc,
		documentStatusSvc:         documentStatusSvc,
	}
}

func (drs *DocumentReviewService) GetComments(id uint) (*documentDTOs.DocumentVersionWorkflowDTO, error) {
	var documentVersionWorkflowDTO documentDTOs.DocumentVersionWorkflowDTO

	document, err := drs.documentSvc.GetFullDetailsByID(id, nil)
	if err != nil {
		return nil, err
	}

	if document.DocumentStatus.Type == "draft" {
		return nil, utility.NewForbiddenError("Document must be submitted for review to access comments.")
	}

	currentApprovalLevel := document.LatestDocumentVersion.NextReviewLevel

	documentVersionWorkflows, err := drs.repo.GetComments(document.LatestDocumentVersion.ID, *currentApprovalLevel)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	err = mapstructure.Decode(documentVersionWorkflows, &documentVersionWorkflowDTO)
	if err != nil {
		return nil, err
	}

	return &documentVersionWorkflowDTO, nil
}

func (drs *DocumentReviewService) AddComment(id uint, currentUserID uint, documentReviewCommentFormDTO *documentDTOs.DocumentReviewCommentFormDTO) error {
	document, err := drs.documentSvc.GetFullDetailsByID(id, nil)
	if err != nil {
		return err
	}

	currentApprovalLevel := document.LatestDocumentVersion.NextReviewLevel

	currentDocumentWorkflow, err := drs.documentWorkflowSvc.GetDocumentWorkflowByLevel(document.LatestDocumentVersion.ID, *currentApprovalLevel)
	if err != nil {
		return err
	}

	if documentReviewCommentFormDTO.IsApproved {
		documentWorkflowStatus, err := drs.documentWorkflowStatusSvc.GetByType("approve")
		if err != nil {
			return err
		}
		switch currentDocumentWorkflow.WorkflowLevelType.Type {
		case "individual":
			if currentUserID != currentDocumentWorkflow.DocumentWorkflowParticipants[0].UserID {
				return utility.NewForbiddenError("you have no permission to approve")
			}

			documentWorkflowReview := &models.DocumentWorkflowReview{
				DocumentWorkflowID:       currentDocumentWorkflow.ID,
				DocumentWorkflowStatusID: documentWorkflowStatus.ID,
				Comment:                  documentReviewCommentFormDTO.Comment,
				CreatedBYID:              currentUserID,
			}

			err = drs.repo.CreateDocumentWorkflowReview(documentWorkflowReview)
			if err != nil {
				return err
			}

			documentWorkflow, err := drs.documentWorkflowSvc.GetDocumentWorkflowByID(currentDocumentWorkflow.ID)
			if err != nil {
				return err
			}

			documentWorkflow.DocumentWorkflowStatusID = documentWorkflowStatus.ID
			err = drs.documentWorkflowSvc.UpdateDocumentWorkflow(documentWorkflow)
			if err != nil {
				return err
			}

			nextDocWorkflowLevel, err := drs.documentWorkflowSvc.GetNextLevel(document.LatestDocumentVersion.ID, *currentApprovalLevel)
			if err != nil {
				return err // Handle other errors that might occur
			}

			if nextDocWorkflowLevel == nil {
				documentStatus, err := drs.documentStatusSvc.GetByType("approved")
				if err != nil {
					return err
				}
				existingDocument, err := drs.documentSvc.GetByID(id)
				if err != nil {
					return err
				}

				existingDocument.DocumentStatusID = documentStatus.ID
				err = drs.documentSvc.UpdateDocument(existingDocument)
				if err != nil {
					return err
				}
			} else {
				existingdocVersion, err := drs.documentSvc.GetDocumentVersionByID(document.LatestDocumentVersion.ID)
				if err != nil {
					return err
				}

				existingdocVersion.NextReviewLevel = &nextDocWorkflowLevel.Level
				err = drs.documentSvc.UpdateDocumentVersion(existingdocVersion)
				if err != nil {
					return err
				}
			}
		case "multiple-users", "role":
			// Check if the current user ID is in the list of participants
			allowed := false
			for _, participant := range currentDocumentWorkflow.DocumentWorkflowParticipants {
				if participant.UserID == currentUserID {
					allowed = true
					break
				}
			}

			if !allowed {
				return utility.NewForbiddenError("you have no permission to approve")
			}

			documentWorkflowReview := &models.DocumentWorkflowReview{
				DocumentWorkflowID:       currentDocumentWorkflow.ID,
				DocumentWorkflowStatusID: documentWorkflowStatus.ID,
				Comment:                  documentReviewCommentFormDTO.Comment,
				CreatedBYID:              currentUserID,
			}

			err = drs.repo.CreateDocumentWorkflowReview(documentWorkflowReview)
			if err != nil {
				return err
			}

			if currentDocumentWorkflow.ApprovalTypes == "all" {
				allApproved, err := drs.documentWorkflowSvc.CheckAllParticipantsApproved(currentDocumentWorkflow.ID, currentDocumentWorkflow.DocumentWorkflowParticipants)
				if err != nil {
					return err
				}

				if !allApproved {
					return nil // Not all users have approved, so do not proceed to next level
				} else {
					documentWorkflow, err := drs.documentWorkflowSvc.GetDocumentWorkflowByID(currentDocumentWorkflow.ID)
					if err != nil {
						return err
					}

					documentWorkflow.DocumentWorkflowStatusID = documentWorkflowStatus.ID
					err = drs.documentWorkflowSvc.UpdateDocumentWorkflow(documentWorkflow)
					if err != nil {
						return err
					}

					nextDocWorkflowLevel, err := drs.documentWorkflowSvc.GetNextLevel(document.LatestDocumentVersion.ID, *currentApprovalLevel)
					if err != nil {
						return err // Handle other errors that might occur
					}

					if nextDocWorkflowLevel == nil {
						documentStatus, err := drs.documentStatusSvc.GetByType("approved")
						if err != nil {
							return err
						}
						existingDocument, err := drs.documentSvc.GetByID(id)
						if err != nil {
							return err
						}

						existingDocument.DocumentStatusID = documentStatus.ID
						err = drs.documentSvc.UpdateDocument(existingDocument)
						if err != nil {
							return err
						}
					} else {
						existingdocVersion, err := drs.documentSvc.GetDocumentVersionByID(document.LatestDocumentVersion.ID)
						if err != nil {
							return err
						}

						existingdocVersion.NextReviewLevel = &nextDocWorkflowLevel.Level
						err = drs.documentSvc.UpdateDocumentVersion(existingdocVersion)
						if err != nil {
							return err
						}
					}
				}
			} else {
				documentWorkflow, err := drs.documentWorkflowSvc.GetDocumentWorkflowByID(currentDocumentWorkflow.ID)
				if err != nil {
					return err
				}

				documentWorkflow.DocumentWorkflowStatusID = documentWorkflowStatus.ID
				err = drs.documentWorkflowSvc.UpdateDocumentWorkflow(documentWorkflow)
				if err != nil {
					return err
				}

				nextDocWorkflowLevel, err := drs.documentWorkflowSvc.GetNextLevel(document.LatestDocumentVersion.ID, *currentApprovalLevel)
				if err != nil {
					return err // Handle other errors that might occur
				}

				if nextDocWorkflowLevel == nil {
					documentStatus, err := drs.documentStatusSvc.GetByType("approved")
					if err != nil {
						return err
					}
					existingDocument, err := drs.documentSvc.GetByID(id)
					if err != nil {
						return err
					}

					existingDocument.DocumentStatusID = documentStatus.ID
					err = drs.documentSvc.UpdateDocument(existingDocument)
					if err != nil {
						return err
					}
				} else {
					existingdocVersion, err := drs.documentSvc.GetDocumentVersionByID(document.LatestDocumentVersion.ID)
					if err != nil {
						return err
					}

					existingdocVersion.NextReviewLevel = &nextDocWorkflowLevel.Level
					err = drs.documentSvc.UpdateDocumentVersion(existingdocVersion)
					if err != nil {
						return err
					}
				}
			}
		}

	} else if documentReviewCommentFormDTO.IsRejected {
		documentWorkflowStatus, err := drs.documentWorkflowStatusSvc.GetByType("reject")
		if err != nil {
			return err
		}

		documentWorkflowReview := &models.DocumentWorkflowReview{
			DocumentWorkflowID:       currentDocumentWorkflow.ID,
			DocumentWorkflowStatusID: documentWorkflowStatus.ID,
			Comment:                  documentReviewCommentFormDTO.Comment,
			CreatedBYID:              currentUserID,
		}

		err = drs.repo.CreateDocumentWorkflowReview(documentWorkflowReview)
		if err != nil {
			return err
		}

		documentWorkflow, err := drs.documentWorkflowSvc.GetDocumentWorkflowByID(currentDocumentWorkflow.ID)
		if err != nil {
			return err
		}

		documentWorkflow.DocumentWorkflowStatusID = documentWorkflowStatus.ID
		err = drs.documentWorkflowSvc.UpdateDocumentWorkflow(documentWorkflow)
		if err != nil {
			return err
		}

		documentStatus, err := drs.documentStatusSvc.GetByType("rejected")
		if err != nil {
			return err
		}
		existingDocument, err := drs.documentSvc.GetByID(id)
		if err != nil {
			return err
		}

		existingDocument.DocumentStatusID = documentStatus.ID
		err = drs.documentSvc.UpdateDocument(existingDocument)
		if err != nil {
			return err
		}

	} else if documentReviewCommentFormDTO.IsReverted {
		documentWorkflowStatus, err := drs.documentWorkflowStatusSvc.GetByType("revert")
		if err != nil {
			return err
		}
		documentWorkflowReview := &models.DocumentWorkflowReview{
			DocumentWorkflowID:       currentDocumentWorkflow.ID,
			DocumentWorkflowStatusID: documentWorkflowStatus.ID,
			Comment:                  documentReviewCommentFormDTO.Comment,
			CreatedBYID:              currentUserID,
		}

		err = drs.repo.CreateDocumentWorkflowReview(documentWorkflowReview)
		if err != nil {
			return err
		}

	} else {
		documentWorkflowStatus, err := drs.documentWorkflowStatusSvc.GetByType("comment")
		if err != nil {
			return err
		}
		documentWorkflowReview := &models.DocumentWorkflowReview{
			DocumentWorkflowID:       currentDocumentWorkflow.ID,
			DocumentWorkflowStatusID: documentWorkflowStatus.ID,
			Comment:                  documentReviewCommentFormDTO.Comment,
			CreatedBYID:              currentUserID,
		}

		err = drs.repo.CreateDocumentWorkflowReview(documentWorkflowReview)
		if err != nil {
			return err
		}
	}

	return nil
}
