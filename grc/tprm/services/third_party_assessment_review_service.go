package services

import (
	"encoding/json"
	"fmt"
	commonDTOs "jethur/grc/common/dtos"
	commonServices "jethur/grc/common/services"
	logServices "jethur/grc/logs/services"
	masterServices "jethur/grc/masters/services"
	"jethur/grc/models"
	organizationSvc "jethur/grc/organization/services"
	riskServices "jethur/grc/risk/services"
	taskServices "jethur/grc/task/services"
	tprmDTOs "jethur/grc/tprm/dtos"
	"jethur/grc/tprm/repositories"
	"jethur/utility"
	"time"
)

type ThirdPartyAssessmentReviewServiceInterface interface {
	Submit(assessmentID uint, currentUserID uint) error
	Approve(assessmentID uint, currentUserID uint) error
	Reject(assessmentID uint, currentUserID uint) error
	Reopen(assessmentID uint, currentUserID uint) error
}

type ThirdPartyAssessmentReviewService struct {
	repo                      repositories.ThirdPartyAssessmentReviewRepositoryInterface
	logService                logServices.ActivityLogServiceInterface
	tprmAssessmentSvc         ThirdPartyAssessmentServiceInterface
	tprmAssessmentWorkflowSvc ThirdPartyAssessmentWorkflowServiceInterface
	workflowStatusSvc         masterServices.WorkflowStatusServiceInterface
	tprmAssessmentStatusSvc   ThirdPartyAssessmentStatusServiceInterface
	userSvc                   organizationSvc.UserServiceInterface
	notificationSvc           commonServices.NotificationServiceInterface
	taskSvc                   taskServices.TaskServiceInterface
	tprmSettingSvc            TPRMSettingServiceInterface
	tprmComplianceResponseSvc TPRMComplianceResponseServiceInterface
	riskLevelSvc              riskServices.RiskLevelServiceInterface
	thirdPartyRegisterSvc     ThirdPartyRegisterServiceInterface
}

func NewThirdPartyAssessmentReviewService(
	repo repositories.ThirdPartyAssessmentReviewRepositoryInterface,
	logService logServices.ActivityLogServiceInterface,
	tprmAssessmentSvc ThirdPartyAssessmentServiceInterface,
	tprmAssessmentWorkflowSvc ThirdPartyAssessmentWorkflowServiceInterface,
	workflowStatusSvc masterServices.WorkflowStatusServiceInterface,
	tprmAssessmentStatusSvc ThirdPartyAssessmentStatusServiceInterface,
	userSvc organizationSvc.UserServiceInterface,
	notificationSvc commonServices.NotificationServiceInterface,
	taskSvc taskServices.TaskServiceInterface,
	tprmSettingSvc TPRMSettingServiceInterface,
	tprmComplianceResponseSvc TPRMComplianceResponseServiceInterface,
	riskLevelSvc riskServices.RiskLevelServiceInterface,
	thirdPartyRegisterSvc ThirdPartyRegisterServiceInterface,
) ThirdPartyAssessmentReviewServiceInterface {
	return &ThirdPartyAssessmentReviewService{
		repo:                      repo,
		logService:                logService,
		tprmAssessmentSvc:         tprmAssessmentSvc,
		tprmAssessmentWorkflowSvc: tprmAssessmentWorkflowSvc,
		workflowStatusSvc:         workflowStatusSvc,
		tprmAssessmentStatusSvc:   tprmAssessmentStatusSvc,
		userSvc:                   userSvc,
		notificationSvc:           notificationSvc,
		taskSvc:                   taskSvc,
		tprmSettingSvc:            tprmSettingSvc,
		tprmComplianceResponseSvc: tprmComplianceResponseSvc,
		riskLevelSvc:              riskLevelSvc,
		thirdPartyRegisterSvc:     thirdPartyRegisterSvc,
	}
}

func (tars *ThirdPartyAssessmentReviewService) Submit(assessmentID uint, currentUserID uint) error {
	tprmAssesment, err := tars.tprmAssessmentSvc.GetByID(assessmentID)
	if err != nil {
		return err
	}

	tx := tars.repo.BeginTransaction()

	currentDateStr := time.Now()

	assessmentStatus, err := tars.tprmAssessmentStatusSvc.GetByType("submitted")
	if err != nil {
		return err
	}

	level := 1

	tprmAssesment.ThirdPartyAssessmentStatusID = assessmentStatus.ID
	tprmAssesment.SubmittedByID = &currentUserID
	tprmAssesment.SubmittedAt = &currentDateStr
	tprmAssesment.NextReviewLevel = &level
	tprmAssesment.UpdatedBYID = &currentUserID

	err = tars.repo.UpdateAssessment(tprmAssesment, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	assessmentWorkflow, err := tars.tprmAssessmentWorkflowSvc.GetTprmAssessmentWorkflowByLevel(tprmAssesment.ID, level, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	workflowStatus, err := tars.workflowStatusSvc.GetByType("submit")
	if err != nil {
		return err
	}

	tprmReview := &models.ThirdPartyAssessmentWorkflowReview{
		TprmAssessmentWorkflowID: assessmentWorkflow.ID,
		WorkflowStatusID:         workflowStatus.ID,
		CreatedBYID:              currentUserID,
	}

	err = tars.repo.CreateTprmReview(tprmReview, tx)
	if err != nil {
		return err
	}

	tx.Commit()

	urlParam := fmt.Sprintf(`{"id":%d}`, tprmAssesment.ID)
	var taskIDs []uint

	for _, participant := range assessmentWorkflow.ThirdPartyAssessmentWorkflowParticipants {
		taskDTO := commonDTOs.TaskFormDTO{
			ModuleIdentifier:   "tprm",
			DocumentID:         &assessmentID,
			DocumentWorkflowID: &assessmentWorkflow.ID,
			Title:              tprmAssesment.Title,
			URLParam:           &urlParam,
			QueryParam:         nil,
			ResponsibleUserID:  participant.User.ID,
		}

		taskID, err := tars.taskSvc.Create(taskDTO, currentUserID)
		if err != nil {
			return err
		}
		taskIDs = append(taskIDs, taskID)

	}

	for _, taskID := range taskIDs {
		relatedTaskIDs := []uint{}
		for _, otherTaskID := range taskIDs {
			if otherTaskID != taskID {
				relatedTaskIDs = append(relatedTaskIDs, otherTaskID)
			}
		}

		relatedTaskIDsJSON, err := json.Marshal(relatedTaskIDs)
		if err != nil {
			return err
		}

		err = tars.taskSvc.UpdateTaskWithRelatedIDs(taskID, string(relatedTaskIDsJSON))
		if err != nil {
			return err
		}
	}

	go func() {
		urlParam = fmt.Sprintf(`{"id":%d}`, assessmentID)
		entity := "tprm-assessment-review"
		action := "review-user-assign"

		for _, participant := range assessmentWorkflow.ThirdPartyAssessmentWorkflowParticipants {
			dbNotificationDTO := &commonDTOs.DBNotificationDTO{
				UserIDs:              []uint{participant.User.ID},
				ModuleIdentifier:     "tprm",
				NotificationType:     string(commonDTOs.ReviewType),
				ReferenceNumber:      tprmAssesment.ReferenceNumber,
				HeadingEn:            "Third Party Assessment Review",
				HeadingAr:            "مراجعة تقييم الطرف الثالث",
				TitleEn:              tprmAssesment.Title,
				TitleAr:              tprmAssesment.Title,
				Entity:               &entity,
				Action:               &action,
				URLParam:             urlParam,
				IsSystemNotification: false,
				CreatedByID:          &currentUserID,
			}

			responsibleUser, err := tars.userSvc.GetByID(participant.User.ID)
			if err != nil {
				utility.LogError("Failed to fetch review user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}

			assignedBy, err := tars.userSvc.GetByID(currentUserID)
			if err != nil {
				utility.LogError("Failed to fetch assigning user", map[string]interface{}{
					"error": err.Error(),
				})
				return
			}

			emailTo := []string{responsibleUser.Email}
			subject := fmt.Sprintf("Third Party Assessment %s: Review Assigned", tprmAssesment.Title)
			appURL := utility.GetAppURL()
			emailData := map[string]interface{}{
				"FirstNameEn":           responsibleUser.FirstNameEn,
				"LastNameEn":            responsibleUser.LastNameEn,
				"Title":                 tprmAssesment.Title,
				"ReferenceNumber":       tprmAssesment.ReferenceNumber,
				"AssignedByFirstNameEn": assignedBy.FirstNameEn,
				"AssignedByLastNameEn":  assignedBy.LastNameEn,
				"AssignedOn":            time.Now().Format("2006-01-02"),
				"RedirectURL":           fmt.Sprintf("%s/tprm/assessment/%d/overview", appURL, assessmentID),
			}

			emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
				EmailTo:              emailTo,
				Subject:              subject,
				TemplateRelativePath: "tprm/templates/email/third_party_assessment_workflow_review_user_mail.html",
				EmailData:            emailData,
			}

			err = tars.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "tprm_assessment")
			if err != nil {
				utility.LogError("Failed to send notifications", map[string]interface{}{
					"error": err.Error(),
				})
			}
		}
	}()
	return nil
}

func (tars *ThirdPartyAssessmentReviewService) Approve(assessmentID uint, currentUserID uint) error {
	tprmAssessment, err := tars.tprmAssessmentSvc.GetFullDetailsByID(assessmentID)
	if err != nil {
		return err
	}

	thirdPartyRegister, err := tars.thirdPartyRegisterSvc.GetByID(tprmAssessment.ThirdPartyRegister.ID)
	if err != nil {
		return err
	}

	settings, err := tars.tprmSettingSvc.GetTPRMSetting()
	if err != nil {
		return err
	}

	complianceResponses, err := tars.tprmComplianceResponseSvc.GetAll("", "", "", "active")
	if err != nil {
		return err
	}

	riskLevels, err := tars.riskLevelSvc.GetAll("", "", "", "active")
	if err != nil {
		return err
	}

	tx := tars.repo.BeginTransaction()
	currentApprovalLevel := tprmAssessment.NextReviewLevel
	currentWorkflow, err := tars.tprmAssessmentWorkflowSvc.GetTprmAssessmentWorkflowByLevel(assessmentID, *currentApprovalLevel, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	workflowStatus, err := tars.workflowStatusSvc.GetByType("approve")
	if err != nil {
		tx.Rollback()
		return err
	}

	switch currentWorkflow.WorkflowLevelType.Type {
	case "individual":
		if currentUserID != currentWorkflow.ThirdPartyAssessmentWorkflowParticipants[0].UserID {
			tx.Rollback()
			return utility.NewForbiddenError("you have no permission to approve")
		}

		assessmentReview := &models.ThirdPartyAssessmentWorkflowReview{
			TprmAssessmentWorkflowID: currentWorkflow.ID,
			WorkflowStatusID:         workflowStatus.ID,
			CreatedBYID:              currentUserID,
		}

		if err := tars.repo.CreateTprmReview(assessmentReview, tx); err != nil {
			tx.Rollback()
			return err
		}

		assessmentWorkflow, err := tars.tprmAssessmentWorkflowSvc.GetTrpmAsessmentWorkflowByID(currentWorkflow.ID, tx)
		if err != nil {
			tx.Rollback()
			return err
		}

		assessmentWorkflow.WorkflowStatusID = workflowStatus.ID
		if err := tars.tprmAssessmentWorkflowSvc.UpdateTprmAssessmentWorkflow(assessmentWorkflow, tx); err != nil {
			tx.Rollback()
			return err
		}

		_, err = tars.taskSvc.UpdateDocumentTaskStatus(assessmentID, currentUserID, currentWorkflow.ID)
		if err != nil {
			return err
		}

		nextWorkflowLevel, err := tars.tprmAssessmentWorkflowSvc.GetNextLevel(tprmAssessment.ID, *currentApprovalLevel, tx)
		if err != nil {
			tx.Rollback()
			return err
		}

		if nextWorkflowLevel == nil {
			assessmentStatus, err := tars.tprmAssessmentStatusSvc.GetByType("approved")
			if err != nil {
				tx.Rollback()
				return err
			}

			existingAssessment, err := tars.tprmAssessmentSvc.GetByID(assessmentID)
			if err != nil {
				tx.Rollback()
				return err
			}

			if settings.IsExport {
				var maxValue int
				var percentage float64
				maxPossibleScore := 0
				totalScore := 0

				for _, res := range complianceResponses {
					if res.Score > maxValue {
						maxValue = res.Score
					}
				}

				for _, sec := range tprmAssessment.Sections {
					GetAssessmentSectionDetails, err := tars.tprmAssessmentSvc.GetAssessmentSectionDetails(assessmentID, sec.ID)
					if err != nil {
						utility.LogError("Failed to fetch section", map[string]interface{}{
							"error": err.Error(),
						})
					}
					for _, qstn := range GetAssessmentSectionDetails.Questions {
						GetAssessmentQuestionDetails, err := tars.tprmAssessmentSvc.GetAssessmentSectionQuestionDetails(assessmentID, sec.ID, qstn.ID)
						if err != nil {
							utility.LogError("Failed to fetch question", map[string]interface{}{
								"error": err.Error(),
							})
						}
						notApplicable := false
						for _, opt := range GetAssessmentQuestionDetails.Answer.AnswerOptions {
							totalScore += *opt.AnswerOption.Score
							if *opt.AnswerOption.Score == 0 {
								notApplicable = true
							}

						}
						if !notApplicable {
							maxPossibleScore += maxValue
						}
					}
				}

				percentage = float64(totalScore) / float64(maxPossibleScore) * 100

				existingAssessment.Score = &percentage
				thirdPartyRegister.Score = &percentage

				for _, riskLevel := range riskLevels {
					min := float64(*riskLevel.ScoreFrom)
					max := float64(*riskLevel.ScoreTo)
					if percentage >= min && (percentage < max || (percentage == max && max == 100)) {
						existingAssessment.RiskLevelID = &riskLevel.ID
						thirdPartyRegister.RiskLevelID = &riskLevel.ID
						break
					}
				}

			}

			existingAssessment.ThirdPartyAssessmentStatusID = assessmentStatus.ID
			if err := tars.repo.UpdateAssessment(existingAssessment, tx); err != nil {
				tx.Rollback()
				return err
			}

			existingAssessment.ThirdPartyAssessmentStatusID = assessmentStatus.ID
			if err := tars.repo.UpdateThirdPartyRegister(thirdPartyRegister, tx); err != nil {
				tx.Rollback()
				return err
			}
		} else {
			existingAssessment, err := tars.tprmAssessmentSvc.GetByID(assessmentID)
			if err != nil {
				tx.Rollback()
				return err
			}

			existingAssessment.NextReviewLevel = &nextWorkflowLevel.Level
			if err := tars.repo.UpdateAssessment(existingAssessment, tx); err != nil {
				tx.Rollback()
				return err
			}

			nextWorkflow, err := tars.tprmAssessmentWorkflowSvc.GetTprmAssessmentWorkflowByLevel(assessmentID, nextWorkflowLevel.Level, tx)
			if err != nil {
				tx.Rollback()
				return err
			}

			err = tars.CreateTasksForNextLevelUsers(assessmentID, nextWorkflow, tprmAssessment)
			if err != nil {
				return err
			}

			tars.notifyNextLevelUsers(assessmentID, nextWorkflow, tprmAssessment)
		}
	case "multiple-users", "role":
		allowed := false

		for _, particpant := range currentWorkflow.ThirdPartyAssessmentWorkflowParticipants {
			if particpant.UserID == currentUserID {
				allowed = true
				break
			}
		}

		if !allowed {
			tx.Rollback()
			return utility.NewForbiddenError("you have no permission to approve")
		}

		assessmentReview := &models.ThirdPartyAssessmentWorkflowReview{
			TprmAssessmentWorkflowID: currentWorkflow.ID,
			WorkflowStatusID:         workflowStatus.ID,
			CreatedBYID:              currentUserID,
		}

		if err := tars.repo.CreateTprmReview(assessmentReview, tx); err != nil {
			tx.Rollback()
			return err
		}

		if currentWorkflow.ApprovalTypes == "all" {
			_, err := tars.taskSvc.UpdateDocumentTaskStatus(assessmentID, currentUserID, currentWorkflow.ID)
			if err != nil {
				return err
			}
			allApproved, err := tars.tprmAssessmentWorkflowSvc.CheckAllParticipantsApproved(currentWorkflow.ID, currentWorkflow.ThirdPartyAssessmentWorkflowParticipants, tx)
			if err != nil {
				tx.Rollback()
				return err
			}

			if !allApproved {
				tx.Commit()
				return nil
			}

			assessmentWorkflow, err := tars.tprmAssessmentWorkflowSvc.GetTrpmAsessmentWorkflowByID(currentWorkflow.ID, tx)
			if err != nil {
				tx.Rollback()
				return err
			}

			assessmentWorkflow.WorkflowStatusID = workflowStatus.ID
			if err := tars.tprmAssessmentWorkflowSvc.UpdateTprmAssessmentWorkflow(assessmentWorkflow, tx); err != nil {
				tx.Rollback()
				return err
			}

			nextWorkflowLevel, err := tars.tprmAssessmentWorkflowSvc.GetNextLevel(assessmentID, *currentApprovalLevel, tx)
			if err != nil {
				tx.Rollback()
				return err
			}

			if nextWorkflowLevel == nil {
				assessmentStatus, err := tars.tprmAssessmentStatusSvc.GetByType("approved")
				if err != nil {
					tx.Rollback()
					return err
				}

				existingAssessment, err := tars.tprmAssessmentSvc.GetByID(assessmentID)
				if err != nil {
					tx.Rollback()
					return err
				}

				if settings.IsExport {
					var maxValue int
					var percentage float64
					maxPossibleScore := 0
					totalScore := 0

					for _, res := range complianceResponses {
						if res.Score > maxValue {
							maxValue = res.Score
						}
					}

					for _, sec := range tprmAssessment.Sections {
						GetAssessmentSectionDetails, err := tars.tprmAssessmentSvc.GetAssessmentSectionDetails(assessmentID, sec.ID)
						if err != nil {
							utility.LogError("Failed to fetch section", map[string]interface{}{
								"error": err.Error(),
							})
						}
						for _, qstn := range GetAssessmentSectionDetails.Questions {
							GetAssessmentQuestionDetails, err := tars.tprmAssessmentSvc.GetAssessmentSectionQuestionDetails(assessmentID, sec.ID, qstn.ID)
							if err != nil {
								utility.LogError("Failed to fetch question", map[string]interface{}{
									"error": err.Error(),
								})
							}
							notApplicable := false
							for _, opt := range GetAssessmentQuestionDetails.Answer.AnswerOptions {
								totalScore += *opt.AnswerOption.Score
								if *opt.AnswerOption.Score == 0 {
									notApplicable = true
								}

							}
							if !notApplicable {
								maxPossibleScore += maxValue
							}
						}
					}

					percentage = float64(totalScore) / float64(maxPossibleScore) * 100

					existingAssessment.Score = &percentage
					thirdPartyRegister.Score = &percentage

					for _, riskLevel := range riskLevels {
						min := float64(*riskLevel.ScoreFrom)
						max := float64(*riskLevel.ScoreTo)
						if percentage >= min && (percentage < max || (percentage == max && max == 100)) {
							existingAssessment.RiskLevelID = &riskLevel.ID
							thirdPartyRegister.RiskLevelID = &riskLevel.ID
							break
						}
					}

				}

				existingAssessment.ThirdPartyAssessmentStatusID = assessmentStatus.ID
				if err := tars.repo.UpdateAssessment(existingAssessment, tx); err != nil {
					tx.Rollback()
					return err
				}

				existingAssessment.ThirdPartyAssessmentStatusID = assessmentStatus.ID
				if err := tars.repo.UpdateThirdPartyRegister(thirdPartyRegister, tx); err != nil {
					tx.Rollback()
					return err
				}
			} else {
				assessmentStatus, err := tars.tprmAssessmentStatusSvc.GetByType("in-progress")
				if err != nil {
					tx.Rollback()
					return err
				}

				existingAssessment, err := tars.tprmAssessmentSvc.GetByID(assessmentID)
				if err != nil {
					tx.Rollback()
					return err
				}

				existingAssessment.NextReviewLevel = &nextWorkflowLevel.Level
				existingAssessment.ThirdPartyAssessmentStatusID = assessmentStatus.ID
				if err := tars.repo.UpdateAssessment(existingAssessment, tx); err != nil {
					tx.Rollback()
					return err
				}

				nextWorkflow, err := tars.tprmAssessmentWorkflowSvc.GetTprmAssessmentWorkflowByLevel(tprmAssessment.ID, nextWorkflowLevel.Level, tx)
				if err != nil {
					tx.Rollback()
					return err
				}

				err = tars.CreateTasksForNextLevelUsers(assessmentID, nextWorkflow, tprmAssessment)
				if err != nil {
					return err
				}
				tars.notifyNextLevelUsers(tprmAssessment.ID, nextWorkflow, tprmAssessment)
			}
		} else {
			assessmentWorkflow, err := tars.tprmAssessmentWorkflowSvc.GetTrpmAsessmentWorkflowByID(currentWorkflow.ID, tx)
			if err != nil {
				tx.Rollback()
				return err
			}

			assessmentWorkflow.WorkflowStatusID = workflowStatus.ID
			if err := tars.tprmAssessmentWorkflowSvc.UpdateTprmAssessmentWorkflow(assessmentWorkflow, tx); err != nil {
				tx.Rollback()
				return err
			}

			task, err := tars.taskSvc.UpdateDocumentTaskStatus(assessmentID, currentUserID, currentWorkflow.ID)
			if err != nil {
				return err
			}

			var relatedTaskIDs []uint
			if task.RelatedTaskIDs != nil {
				if err := json.Unmarshal([]byte(*task.RelatedTaskIDs), &relatedTaskIDs); err != nil {
					return err
				}

				// Close each related task
				for _, relatedTaskID := range relatedTaskIDs {
					relatedTask, err := tars.taskSvc.GetByID(relatedTaskID)
					if err != nil {
						return err
					}
					_, err = tars.taskSvc.UpdateDocumentTaskStatus(assessmentID, relatedTask.ResponsibleUserID, currentWorkflow.ID)
					if err != nil {
						return err
					}
				}
			}

			nextWorkflowLevel, err := tars.tprmAssessmentWorkflowSvc.GetNextLevel(assessmentID, *currentApprovalLevel, tx)
			if err != nil {
				tx.Rollback()
				return err
			}

			if nextWorkflowLevel == nil {
				assessmentStatus, err := tars.tprmAssessmentStatusSvc.GetByType("approved")
				if err != nil {
					tx.Rollback()
					return err
				}

				existingAssessment, err := tars.tprmAssessmentSvc.GetByID(assessmentID)
				if err != nil {
					tx.Rollback()
					return err
				}

				if settings.IsExport {
					var maxValue int
					var percentage float64
					maxPossibleScore := 0
					totalScore := 0

					for _, res := range complianceResponses {
						if res.Score > maxValue {
							maxValue = res.Score
						}
					}

					for _, sec := range tprmAssessment.Sections {
						GetAssessmentSectionDetails, err := tars.tprmAssessmentSvc.GetAssessmentSectionDetails(assessmentID, sec.ID)
						if err != nil {
							utility.LogError("Failed to fetch section", map[string]interface{}{
								"error": err.Error(),
							})
						}
						for _, qstn := range GetAssessmentSectionDetails.Questions {
							GetAssessmentQuestionDetails, err := tars.tprmAssessmentSvc.GetAssessmentSectionQuestionDetails(assessmentID, sec.ID, qstn.ID)
							if err != nil {
								utility.LogError("Failed to fetch question", map[string]interface{}{
									"error": err.Error(),
								})
							}
							notApplicable := false
							for _, opt := range GetAssessmentQuestionDetails.Answer.AnswerOptions {
								totalScore += *opt.AnswerOption.Score
								if *opt.AnswerOption.Score == 0 {
									notApplicable = true
								}

							}
							if !notApplicable {
								maxPossibleScore += maxValue
							}
						}
					}

					percentage = float64(totalScore) / float64(maxPossibleScore) * 100

					existingAssessment.Score = &percentage
					thirdPartyRegister.Score = &percentage

					for _, riskLevel := range riskLevels {
						min := float64(*riskLevel.ScoreFrom)
						max := float64(*riskLevel.ScoreTo)
						if percentage >= min && (percentage < max || (percentage == max && max == 100)) {
							existingAssessment.RiskLevelID = &riskLevel.ID
							thirdPartyRegister.RiskLevelID = &riskLevel.ID
							break
						}
					}

				}

				existingAssessment.ThirdPartyAssessmentStatusID = assessmentStatus.ID
				if err := tars.repo.UpdateAssessment(existingAssessment, tx); err != nil {
					tx.Rollback()
					return err
				}

				existingAssessment.ThirdPartyAssessmentStatusID = assessmentStatus.ID
				if err := tars.repo.UpdateThirdPartyRegister(thirdPartyRegister, tx); err != nil {
					tx.Rollback()
					return err
				}
			} else {
				assessmentStatus, err := tars.tprmAssessmentStatusSvc.GetByType("in-progress")
				if err != nil {
					tx.Rollback()
					return err
				}

				existingAssessment, err := tars.tprmAssessmentSvc.GetByID(assessmentID)
				if err != nil {
					tx.Rollback()
					return err
				}

				existingAssessment.NextReviewLevel = &nextWorkflowLevel.Level
				existingAssessment.ThirdPartyAssessmentStatusID = assessmentStatus.ID
				if err := tars.repo.UpdateAssessment(existingAssessment, tx); err != nil {
					tx.Rollback()
					return err
				}

				nextWorkflow, err := tars.tprmAssessmentWorkflowSvc.GetTprmAssessmentWorkflowByLevel(tprmAssessment.ID, nextWorkflowLevel.Level, tx)
				if err != nil {
					tx.Rollback()
					return err
				}

				err = tars.CreateTasksForNextLevelUsers(assessmentID, nextWorkflow, tprmAssessment)
				if err != nil {
					return err
				}
				tars.notifyNextLevelUsers(tprmAssessment.ID, nextWorkflow, tprmAssessment)
			}

		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	assessmentDetails, err := tars.tprmAssessmentSvc.GetFullDetailsByID(assessmentID)
	if err != nil {
		tx.Rollback()
		return err
	}
	tars.sendApprovalNotification(assessmentID, currentUserID, assessmentDetails)

	return nil
}

func (tars *ThirdPartyAssessmentReviewService) Reject(assessmentID uint, currentUserID uint) error {
	workflowStatus, err := tars.workflowStatusSvc.GetByType("reject")
	if err != nil {
		return err
	}

	assessmentStatus, err := tars.tprmAssessmentStatusSvc.GetByType("rejected")
	if err != nil {
		return err
	}

	tx := tars.repo.BeginTransaction()
	assessment, err := tars.tprmAssessmentSvc.GetFullDetailsByID(assessmentID)
	if err != nil {
		tx.Rollback()
		return err
	}

	currentApprovalLevel := assessment.NextReviewLevel
	currentWorkflow, err := tars.tprmAssessmentWorkflowSvc.GetTprmAssessmentWorkflowByLevel(assessmentID, *currentApprovalLevel, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	assessmentReview := &models.ThirdPartyAssessmentWorkflowReview{
		TprmAssessmentWorkflowID: currentWorkflow.ID,
		WorkflowStatusID:         workflowStatus.ID,
		CreatedBYID:              currentUserID,
	}

	if err := tars.repo.CreateTprmReview(assessmentReview, tx); err != nil {
		tx.Rollback()
		return err
	}

	assessmentWorkflow, err := tars.tprmAssessmentWorkflowSvc.GetTrpmAsessmentWorkflowByID(currentWorkflow.ID, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	assessmentWorkflow.WorkflowStatusID = workflowStatus.ID
	if err := tars.tprmAssessmentWorkflowSvc.UpdateTprmAssessmentWorkflow(assessmentWorkflow, tx); err != nil {
		tx.Rollback()
		return err
	}

	existingAssessment, err := tars.tprmAssessmentSvc.GetByID(assessmentID)
	if err != nil {
		tx.Rollback()
		return err
	}

	existingAssessment.ThirdPartyAssessmentStatusID = assessmentStatus.ID
	if err := tars.repo.UpdateAssessment(existingAssessment, tx); err != nil {
		tx.Rollback()
		return err
	}

	task, err := tars.taskSvc.UpdateDocumentTaskStatus(assessmentID, currentUserID, currentWorkflow.ID)
	if err != nil {
		return err
	}

	// Parse the related_task_ids from the task
	var relatedTaskIDs []uint
	if task.RelatedTaskIDs != nil {
		if err := json.Unmarshal([]byte(*task.RelatedTaskIDs), &relatedTaskIDs); err != nil {
			return err
		}

		// Close each related task
		for _, relatedTaskID := range relatedTaskIDs {
			relatedTask, err := tars.taskSvc.GetByID(relatedTaskID)
			if err != nil {
				return err
			}
			_, err = tars.taskSvc.UpdateDocumentTaskStatus(assessmentID, relatedTask.ResponsibleUserID, currentWorkflow.ID)
			if err != nil {
				return err
			}
		}
	}
	tx.Commit()
	tars.sendRejectedNotification(assessmentID, currentUserID, assessment)
	return nil
}

func (tars *ThirdPartyAssessmentReviewService) notifyNextLevelUsers(id uint, nextWorkflowLevel *models.ThirdPartyAssessmentWorkflow, assessmentDetail *tprmDTOs.ThirdPartyAssessmentDetailDTO) {
	go func() {
		urlParam := fmt.Sprintf(`{"id":%d}`, assessmentDetail.ID)
		entity := "tprm-assessment-review"
		action := "review-user-assign"

		for _, participant := range nextWorkflowLevel.ThirdPartyAssessmentWorkflowParticipants {
			responsibleUser, err := tars.userSvc.GetByID(participant.UserID)
			if err != nil {
				utility.LogError("Failed to fetch review user", map[string]interface{}{
					"error": err.Error(),
				})
				continue
			}

			createdByID := uint(assessmentDetail.CreatedByUser.ID)
			// Prepare database notification DTO
			dbNotificationDTO := &commonDTOs.DBNotificationDTO{
				UserIDs:              []uint{participant.UserID},
				ModuleIdentifier:     "tprm",
				NotificationType:     string(commonDTOs.ReviewType),
				HeadingEn:            "Third Party Assessment Review",
				HeadingAr:            "مراجعة تقييم الطرف الثالث",
				ReferenceNumber:      assessmentDetail.ReferenceNumber,
				TitleEn:              assessmentDetail.Title,
				TitleAr:              assessmentDetail.Title,
				Entity:               &entity,
				Action:               &action,
				URLParam:             urlParam,
				IsSystemNotification: false,
				CreatedByID:          &createdByID,
			}

			// Prepare email notification DTO
			emailTo := []string{responsibleUser.Email}
			subject := fmt.Sprintf("Third Party Assessment %s: Review Assigned", assessmentDetail.Title)
			appURL := utility.GetAppURL()
			urlParam = fmt.Sprintf("%s/tprm/assessment/%d/overview", appURL, assessmentDetail.ID)

			emailData := map[string]interface{}{
				"FirstNameEn":           responsibleUser.FirstNameEn,
				"LastNameEn":            responsibleUser.LastNameEn,
				"Title":                 assessmentDetail.Title,
				"ReferenceNumber":       assessmentDetail.ReferenceNumber,
				"AssignedByFirstNameEn": assessmentDetail.CreatedByUser.FirstNameEN,
				"AssignedByLastNameEn":  assessmentDetail.CreatedByUser.LastNameEN,
				"AssignedOn":            time.Now().Format("2006-01-02"),
				"RedirectURL":           urlParam,
			}

			emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
				EmailTo:              emailTo,
				Subject:              subject,
				TemplateRelativePath: "tprm/templates/email/third_party_assessment_workflow_review_user_mail.html",
				EmailData:            emailData,
			}

			// Send both database and email notifications
			err = tars.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "tprm_assessment")
			if err != nil {
				utility.LogError("Failed to send next level review notification", map[string]interface{}{
					"error": err.Error(),
				})
			}
		}
	}()
}

func (rrs *ThirdPartyAssessmentReviewService) sendApprovalNotification(id uint, currentUserID uint, assessmentDetail *tprmDTOs.ThirdPartyAssessmentDetailDTO) {
	go func() {
		// Prepare URL and parameters for the notification
		urlParam := fmt.Sprintf(`{"id":%d}`, assessmentDetail.ID)
		entity := "tprm-assessment-review"
		action := "approval"

		// Create the database notification DTO
		dbNotificationDTO := &commonDTOs.DBNotificationDTO{
			UserIDs:              []uint{uint(assessmentDetail.SubmittedByUser.ID)},
			ModuleIdentifier:     "tprm",
			NotificationType:     string(commonDTOs.ApprovalType),
			HeadingEn:            "Third Party Assessment Approval",
			HeadingAr:            "اعتماد تقييم الطرف الثالث",
			ReferenceNumber:      assessmentDetail.ReferenceNumber,
			TitleEn:              assessmentDetail.Title,
			TitleAr:              assessmentDetail.Title,
			Entity:               &entity,
			Action:               &action,
			URLParam:             urlParam,
			IsSystemNotification: false,
			CreatedByID:          &currentUserID,
		}

		// Get the approving user's details
		approvedByUser, err := rrs.userSvc.GetByID(currentUserID)
		if err != nil {
			utility.LogError("Failed to fetch approving user", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		// Prepare email details
		appURL := utility.GetAppURL()
		emailTo := []string{assessmentDetail.SubmittedByUser.Email}
		subject := fmt.Sprintf("Third Party Assessment %s: Approval Notification", assessmentDetail.Title)
		urlParam = fmt.Sprintf("%s/tprm/assessment/%d/overview", appURL, assessmentDetail.ID)

		emailData := map[string]interface{}{
			"Title":                 assessmentDetail.Title,
			"ReferenceNumber":       assessmentDetail.ReferenceNumber,
			"FirstNameEn":           assessmentDetail.SubmittedByUser.FirstNameEN,
			"LastNameEn":            assessmentDetail.SubmittedByUser.LastNameEN,
			"ApprovedByFirstNameEn": approvedByUser.FirstNameEn,
			"ApprovedByLastNameEn":  approvedByUser.LastNameEn,
			"ApprovedAt":            time.Now().Format("2006-01-02"),
			"RedirectURL":           urlParam,
		}

		// Create the email notification DTO
		emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
			EmailTo:              emailTo,
			Subject:              subject,
			TemplateRelativePath: "tprm/templates/email/third_party_assessment_workflow_approve_user_mail.html",
			EmailData:            emailData,
		}

		// Send the notification
		err = rrs.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "tprm_assessment")
		if err != nil {
			utility.LogError("Failed to send approval notification email", map[string]interface{}{
				"error": err.Error(),
			})
		}
	}()
}

func (tars *ThirdPartyAssessmentReviewService) sendRejectedNotification(id uint, currentUserID uint, assessmentDetail *tprmDTOs.ThirdPartyAssessmentDetailDTO) {
	go func() {
		// Prepare URL and parameters for the notification
		urlParam := fmt.Sprintf(`{"id":%d}`, id)
		entity := "tprm-assessment-review"
		action := "rejected"

		// Create the database notification DTO
		dbNotificationDTO := &commonDTOs.DBNotificationDTO{
			UserIDs:              []uint{uint(assessmentDetail.SubmittedByUser.ID)},
			ModuleIdentifier:     "tprm",
			NotificationType:     string(commonDTOs.RejectType),
			HeadingEn:            "Third Party Assessment Rejected",
			HeadingAr:            "تم رفض تقييم الطرف الثالث",
			ReferenceNumber:      assessmentDetail.ReferenceNumber,
			TitleEn:              assessmentDetail.Title,
			TitleAr:              assessmentDetail.Title,
			Entity:               &entity,
			Action:               &action,
			URLParam:             urlParam,
			IsSystemNotification: false,
			CreatedByID:          &currentUserID,
		}

		rejectedByUser, err := tars.userSvc.GetByID(currentUserID)
		if err != nil {
			utility.LogError("Failed to fetch review performed user", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		// Prepare email details
		appURL := utility.GetAppURL()
		emailTo := []string{assessmentDetail.SubmittedByUser.Email}
		subject := fmt.Sprintf("Third Party Assessment %s: Rejected Notification", assessmentDetail.Title)
		urlParam = fmt.Sprintf("%s/tprm/assessment/%d/overview", appURL, assessmentDetail.ID)

		emailData := map[string]interface{}{
			"FirstNameEn":           assessmentDetail.SubmittedByUser.FirstNameEN,
			"LastNameEn":            assessmentDetail.SubmittedByUser.LastNameEN,
			"Title":                 assessmentDetail.Title,
			"RejectedByFirstNameEn": rejectedByUser.FirstNameEn,
			"RejectedByLastNameEn":  rejectedByUser.LastNameEn,
			"ApprovedAt":            time.Now().Format("2006-01-02"),
			"RedirectURL":           urlParam,
		}

		// Create the email notification DTO
		emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
			EmailTo:              emailTo,
			Subject:              subject,
			TemplateRelativePath: "tprm/templates/email/third_party_assessment_workflow_reject_user_mail.html",
			EmailData:            emailData,
		}

		// Send the notification
		err = tars.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "tprm_assessment")
		if err != nil {
			utility.LogError("Failed to send rejected notification email", map[string]interface{}{
				"error": err.Error(),
			})
		}
	}()
}

func (tars *ThirdPartyAssessmentReviewService) CreateTasksForNextLevelUsers(assessmentID uint, workflowLevel *models.ThirdPartyAssessmentWorkflow, assessment *tprmDTOs.ThirdPartyAssessmentDetailDTO) error {
	urlParam := fmt.Sprintf(`{"id":%d}`, assessmentID)
	var taskIDs []uint

	createdByID := uint(assessment.SubmittedByUser.ID)
	// Loop through each participant and create a task
	for _, participant := range workflowLevel.ThirdPartyAssessmentWorkflowParticipants {
		taskDTO := commonDTOs.TaskFormDTO{
			ModuleIdentifier:   "tprm",
			DocumentID:         &assessmentID,
			DocumentWorkflowID: &workflowLevel.ID,
			Title:              assessment.Title,
			URLParam:           &urlParam,
			QueryParam:         nil,
			ResponsibleUserID:  participant.User.ID,
		}

		taskID, err := tars.taskSvc.Create(taskDTO, createdByID)
		if err != nil {
			return err
		}
		taskIDs = append(taskIDs, taskID)
	}

	// Set up related task IDs for each task
	for _, taskID := range taskIDs {
		relatedTaskIDs := []uint{}
		for _, otherTaskID := range taskIDs {
			if otherTaskID != taskID {
				relatedTaskIDs = append(relatedTaskIDs, otherTaskID)
			}
		}

		relatedTaskIDsJSON, err := json.Marshal(relatedTaskIDs)
		if err != nil {
			return err
		}

		err = tars.taskSvc.UpdateTaskWithRelatedIDs(taskID, string(relatedTaskIDsJSON))
		if err != nil {
			return err
		}
	}

	return nil
}

func (tars *ThirdPartyAssessmentReviewService) Reopen(assessmentID uint, currentUserID uint) error {

	rejectedStatus, err := tars.tprmAssessmentStatusSvc.GetByType("rejected")
	if err != nil {
		return err
	}

	assessment, err := tars.tprmAssessmentSvc.GetFullDetailsByID(assessmentID)
	if err != nil {
		return err
	}

	if assessment.ThirdPartyAssessmentStatus.ID != rejectedStatus.ID {
		return utility.NewForbiddenError("not in rejected status")
	}

	assessmentStatus, err := tars.tprmAssessmentStatusSvc.GetByType("draft")
	if err != nil {
		return err
	}

	tx := tars.repo.BeginTransaction()

	assessmentDetails, err := tars.repo.GetAssessmentDetails(assessmentID, tx)
	if err != nil {
		tx.Rollback()
	}

	var deleteWorkflowReviewIds []*uint

	if len(assessmentDetails.ThirdPartyAssessmentWorkflows) > 0 {
		for _, workflow := range assessmentDetails.ThirdPartyAssessmentWorkflows {

			if len(workflow.ThirdPartyAssessmentWorkflowReviews) > 0 {
				for _, reviews := range workflow.ThirdPartyAssessmentWorkflowReviews {
					deleteWorkflowReviewIds = append(deleteWorkflowReviewIds, &reviews.ID)
				}
			}
		}
	}

	fmt.Println(deleteWorkflowReviewIds)
	if err := tars.repo.DeleteTprmAssessmentWorkflowReviews(deleteWorkflowReviewIds); err != nil {
		tx.Rollback()
	}

	existingAssessment, err := tars.tprmAssessmentSvc.GetByID(assessmentID)
	if err != nil {
		tx.Rollback()
		return err
	}

	existingAssessment.SubmittedByID = nil
	existingAssessment.SubmittedAt = nil
	existingAssessment.ThirdPartyAssessmentStatusID = assessmentStatus.ID
	if err := tars.repo.UpdateAssessment(existingAssessment, tx); err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	tars.sendReopenNotification(assessmentID, currentUserID, assessment)
	return nil
}

func (tars *ThirdPartyAssessmentReviewService) sendReopenNotification(id uint, currentUserID uint, assessmentDetail *tprmDTOs.ThirdPartyAssessmentDetailDTO) {
	go func() {
		// Prepare URL and parameters for the notification
		urlParam := fmt.Sprintf(`{"id":%d}`, id)
		entity := "tprm-assessment-review"
		action := "reopen"

		// Create the database notification DTO
		dbNotificationDTO := &commonDTOs.DBNotificationDTO{
			UserIDs:              []uint{uint(assessmentDetail.SubmittedByUser.ID)},
			ModuleIdentifier:     "tprm",
			NotificationType:     string(commonDTOs.ReviewType),
			HeadingEn:            "Third Party Assessment Reopened",
			HeadingAr:            "إعادة فتح تقييم الطرف الثالث",
			ReferenceNumber:      assessmentDetail.ReferenceNumber,
			TitleEn:              assessmentDetail.Title,
			TitleAr:              assessmentDetail.Title,
			Entity:               &entity,
			Action:               &action,
			URLParam:             urlParam,
			IsSystemNotification: false,
			CreatedByID:          &currentUserID,
		}

		reopenedByUser, err := tars.userSvc.GetByID(currentUserID)
		if err != nil {
			utility.LogError("Failed to fetch review performed user", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		// Prepare email details
		appURL := utility.GetAppURL()
		emailTo := []string{assessmentDetail.CreatedByUser.Email}
		subject := fmt.Sprintf("Third Party Assessment Reopened: %s", assessmentDetail.Title)
		urlParam = fmt.Sprintf("%s/tprm/assessment/%d/overview", appURL, assessmentDetail.ID)

		emailData := map[string]interface{}{
			"FirstNameEn":           assessmentDetail.CreatedByUser.FirstNameEN,
			"LastNameEn":            assessmentDetail.CreatedByUser.LastNameEN,
			"Title":                 assessmentDetail.Title,
			"ReopenedByFirstNameEn": reopenedByUser.FirstNameEn,
			"ReopenedByLastNameEn":  reopenedByUser.LastNameEn,
			"ReopenedAt":            time.Now().Format("2006-01-02"),
			"RedirectURL":           urlParam,
		}

		// Create the email notification DTO
		emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
			EmailTo:              emailTo,
			Subject:              subject,
			TemplateRelativePath: "tprm/templates/email/third_party_assessment_workflow_reopen_user_mail.html",
			EmailData:            emailData,
		}

		// Send the notification
		err = tars.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "tprm_assessment")
		if err != nil {
			utility.LogError("Failed to send reopen notification email", map[string]interface{}{
				"error": err.Error(),
			})
		}
	}()
}
