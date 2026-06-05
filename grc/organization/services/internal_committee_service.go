package services

import (
	"errors"
	"fmt"
	commonDTOs "jethur/grc/common/dtos"
	commonServices "jethur/grc/common/services"
	logServices "jethur/grc/logs/services"
	masterServices "jethur/grc/masters/services"
	"jethur/grc/models"
	organizationDTOs "jethur/grc/organization/dtos"
	organizationRepositories "jethur/grc/organization/repositories"
	settingsDTOs "jethur/grc/settings/dtos"
	settingsServices "jethur/grc/settings/services"
	"jethur/utility"
	"time"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type InternalCommitteeServiceInterface interface {
	GetAll(params *organizationDTOs.InternalCommitteeListQueryParamsDTO) (*organizationDTOs.InternalCommitteeListResponseDTO, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*organizationDTOs.InternalCommitteeDetailDTO, error)
	GetByID(id uint) (*models.InternalCommittee, error)
	Create(currentUserID uint, internalCommitteeFormDTO *organizationDTOs.InternalCommitteeFormDTO) (uint, error)
	Update(id uint, currentUserID uint, internalCommitteeFormDTO *organizationDTOs.InternalCommitteeFormDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)
	ReformCommittee(currentUserID uint, reformDTO *organizationDTOs.ReformInternalCommitteeDTO) (uint, error)
	GetLatestVersionByID(id uint) (*models.InternalCommittee, error)
	DisbandCommittee(committeeID uint, currentUserID uint) error

	// Participant methods
	CreateParticipants(currentUserID uint, nomineeDepartmentID uint, dto *organizationDTOs.InternalCommitteeParticipantFormDTO) ([]uint, error)
	DeleteParticipant(nomineeDepartmentID uint, participantID uint, currentUserID uint) error
	ValidateAllDepartmentsHaveParticipants(versionID uint, tx ...utility.DatabaseTransaction) (bool, error)
}

type InternalCommitteeService struct {
	repo                       organizationRepositories.InternalCommitteeRepositoryInterface
	activityLogSvc             logServices.ActivityLogServiceInterface
	userSvc                    UserServiceInterface
	organizationUnitSvc        settingsServices.OrganizationUnitServiceInterface
	committeeTypeSvc           CommitteeTypeServiceInterface
	committeeMemberTypeSvc     CommitteeMemberTypeServiceInterface
	committeeTermSvc           CommitteeTermServiceInterface
	committeeClassificationSvc CommitteeClassificationServiceInterface
	internalCommitteeStatusSvc InternalCommitteeStatusServiceInterface
	notificationSvc            commonServices.NotificationServiceInterface
	organizationUnitsSvc       OrganizationUnitsServiceInterface
	workflowReviewSvc          InternalCommitteeWorkflowReviewServiceInterface
	committeeWorkflowSvc       InternalCommitteeWorkflowServiceInterface
	workflowStatusSvc          masterServices.WorkflowStatusServiceInterface
}

func NewInternalCommitteeService(
	repo organizationRepositories.InternalCommitteeRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	userSvc UserServiceInterface,
	organizationUnitSvc settingsServices.OrganizationUnitServiceInterface,
	committeeTypeSvc CommitteeTypeServiceInterface,
	committeeMemberTypeSvc CommitteeMemberTypeServiceInterface,
	committeeTermSvc CommitteeTermServiceInterface,
	committeeClassificationSvc CommitteeClassificationServiceInterface,
	internalCommitteeStatusSvc InternalCommitteeStatusServiceInterface,
	notificationSvc commonServices.NotificationServiceInterface,
	organizationUnitsSvc OrganizationUnitsServiceInterface,
	workflowReviewSvc InternalCommitteeWorkflowReviewServiceInterface,
	committeeWorkflowSvc InternalCommitteeWorkflowServiceInterface,
	workflowStatusSvc masterServices.WorkflowStatusServiceInterface,
) InternalCommitteeServiceInterface {
	return &InternalCommitteeService{
		repo:                       repo,
		activityLogSvc:             activityLogSvc,
		userSvc:                    userSvc,
		organizationUnitSvc:        organizationUnitSvc,
		committeeTypeSvc:           committeeTypeSvc,
		committeeMemberTypeSvc:     committeeMemberTypeSvc,
		committeeTermSvc:           committeeTermSvc,
		committeeClassificationSvc: committeeClassificationSvc,
		internalCommitteeStatusSvc: internalCommitteeStatusSvc,
		notificationSvc:            notificationSvc,
		organizationUnitsSvc:       organizationUnitsSvc,
		workflowReviewSvc:          workflowReviewSvc,
		committeeWorkflowSvc:       committeeWorkflowSvc,
		workflowStatusSvc:          workflowStatusSvc,
	}
}

func (ics *InternalCommitteeService) GetAll(params *organizationDTOs.InternalCommitteeListQueryParamsDTO) (*organizationDTOs.InternalCommitteeListResponseDTO, error) {
	return ics.repo.GetAll(params)
}

func (ics *InternalCommitteeService) GetByID(id uint) (*models.InternalCommittee, error) {
	internalCommittee, err := ics.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return internalCommittee, nil
}

func (ics *InternalCommitteeService) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*organizationDTOs.InternalCommitteeDetailDTO, error) {
	var internalCommitteeDTO organizationDTOs.InternalCommitteeDetailDTO

	internalCommittee, err := ics.repo.GetFullDetailsByID(id, tx...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	err = mapstructure.Decode(internalCommittee, &internalCommitteeDTO)
	if err != nil {
		return nil, err
	}

	// Map latest version if exists
	if len(internalCommittee.InternalCommitteeVersions) > 0 {
		latestVersion := internalCommittee.InternalCommitteeVersions[0]
		var versionDTO organizationDTOs.InternalCommitteeVersionDetailDTO

		err = mapstructure.Decode(latestVersion, &versionDTO)
		if err != nil {
			return nil, err
		}

		// Format term end date string
		if latestVersion.TermEndDate != nil {
			termEndDateStr := latestVersion.TermEndDate.Format("2006-01-02")
			versionDTO.TermEndDateString = &termEndDateStr
		}

		versionDTO.CreatedAt = latestVersion.CreatedAt
		if latestVersion.UpdatedAt != nil {
			versionDTO.UpdatedAt = latestVersion.UpdatedAt
		}

		// Manually map nominee departments with participants
		if len(latestVersion.InternalCommitteeVersionNomineeDepartments) > 0 {
			var nomineeDepartmentsDTO []*organizationDTOs.InternalCommitteeVersionNomineeDepartmentDetailDTO
			for _, nomineeDept := range latestVersion.InternalCommitteeVersionNomineeDepartments {
				nomineeDeptDTO := &organizationDTOs.InternalCommitteeVersionNomineeDepartmentDetailDTO{
					ID:                         nomineeDept.ID,
					InternalCommitteeVersionID: nomineeDept.InternalCommitteeVersionID,
					NomineeDepartmentID:        nomineeDept.NomineeDepartmentID,
					CommitteeMemberTypeID:      nomineeDept.CommitteeMemberTypeID,
					MembershipTypeID:           nomineeDept.CommitteeMemberTypeID,
				}

				if nomineeDept.CommitteeMemberType != nil {
					nomineeDeptDTO.MembershipType = &organizationDTOs.InternalCommitteeNomineeMembershipTypeDetailDTO{
						ID:      nomineeDept.CommitteeMemberType.ID,
						TitleEn: nomineeDept.CommitteeMemberType.TitleEn,
						TitleAr: nomineeDept.CommitteeMemberType.TitleAr,
						Status:  nomineeDept.CommitteeMemberType.Status,
					}
				}

				// Map nominee department
				if nomineeDept.NomineeDepartment != nil {
					nomineeDeptDTO.NomineeDepartment = &settingsDTOs.OrganizationUnitDTO{
						ID:      int(nomineeDept.NomineeDepartment.ID),
						TitleEn: nomineeDept.NomineeDepartment.TitleEn,
						TitleAr: nomineeDept.NomineeDepartment.TitleAr,
						Status:  nomineeDept.NomineeDepartment.Status,
					}
					if nomineeDept.NomineeDepartment.OrganizationUnitHead != nil {
						nomineeDeptDTO.NomineeDepartment.OrganizationUnitHead = &settingsDTOs.UserMinimalDTO{
							ID:               int(nomineeDept.NomineeDepartment.OrganizationUnitHead.ID),
							FirstNameEN:      nomineeDept.NomineeDepartment.OrganizationUnitHead.FirstNameEn,
							LastNameEN:       nomineeDept.NomineeDepartment.OrganizationUnitHead.LastNameEn,
							FirstNameAR:      nomineeDept.NomineeDepartment.OrganizationUnitHead.FirstNameAr,
							LastNameAR:       nomineeDept.NomineeDepartment.OrganizationUnitHead.LastNameAr,
							ProfileImagePath: nomineeDept.NomineeDepartment.OrganizationUnitHead.ProfileImagePath,
							Email:            nomineeDept.NomineeDepartment.OrganizationUnitHead.Email,
							Initial:          nomineeDept.NomineeDepartment.OrganizationUnitHead.Initial,
							Color:            nomineeDept.NomineeDepartment.OrganizationUnitHead.Color,
						}
					}
				}

				// Map participants
				if len(nomineeDept.Participants) > 0 {
					var participantsDTO []organizationDTOs.InternalCommitteeVersionNomineeDepartmentParticipantDTO
					for _, participant := range nomineeDept.Participants {
						participantDTO := organizationDTOs.InternalCommitteeVersionNomineeDepartmentParticipantDTO{
							ID:     participant.ID,
							UserID: participant.UserID,
						}

						// Map participant user
						if participant.User != nil {
							participantDTO.User = &organizationDTOs.UserMinimalDTO{
								ID:               int(participant.User.ID),
								FirstNameEN:      participant.User.FirstNameEn,
								LastNameEN:       participant.User.LastNameEn,
								FirstNameAR:      participant.User.FirstNameAr,
								LastNameAR:       participant.User.LastNameAr,
								Email:            participant.User.Email,
								Initial:          participant.User.Initial,
								Color:            participant.User.Color,
								ProfileImagePath: participant.User.ProfileImagePath,
							}
						}

						participantsDTO = append(participantsDTO, participantDTO)
					}
					nomineeDeptDTO.Participants = participantsDTO
				}

				nomineeDepartmentsDTO = append(nomineeDepartmentsDTO, nomineeDeptDTO)
			}
			versionDTO.NomineeDepartments = nomineeDepartmentsDTO
		}

		internalCommitteeDTO.LatestVersion = &versionDTO
	}

	internalCommitteeDTO.CreatedAt = internalCommittee.CreatedAt
	if internalCommittee.UpdatedAt != nil {
		internalCommitteeDTO.UpdatedAt = internalCommittee.UpdatedAt
	}

	return &internalCommitteeDTO, nil
}

func (ics *InternalCommitteeService) Create(currentUserID uint, internalCommitteeFormDTO *organizationDTOs.InternalCommitteeFormDTO) (uint, error) {
	// Validate form data
	if err := ics.validateInternalCommitteeForm(internalCommitteeFormDTO); err != nil {
		return 0, err
	}

	parsedDates, err := ics.parseInternalCommitteeDates(internalCommitteeFormDTO)
	if err != nil {
		return 0, err
	}

	referenceNumber, err := ics.repo.GenerateReferenceNumber()
	if err != nil {
		return 0, err
	}

	requestingEntityID, err := ics.userSvc.GetUserPrimaryUnitID(uint(currentUserID))
	if err != nil {
		return 0, utility.NewValidationError(fmt.Errorf("organization unit not set for current user"), "not_exist", "current user unit")
	}

	tx := ics.repo.BeginTransaction()
	// Create internal committee
	internalCommittee := &models.InternalCommittee{
		RequestingEntityID: requestingEntityID,
		ReferenceNumber:    referenceNumber,
		CreatedBYID:        currentUserID,
	}

	committeeID, err := ics.repo.Create(internalCommittee, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	// Get draft status
	draftStatus, err := ics.internalCommitteeStatusSvc.GetByType("draft")
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	// Create version
	version := float64(1)
	isLatest := true
	internalCommitteeVersion := &models.InternalCommitteeVersion{
		InternalCommitteeID:       committeeID,
		CommitteeNameEn:           internalCommitteeFormDTO.CommitteeNameEn,
		CommitteeNameAr:           internalCommitteeFormDTO.CommitteeNameAr,
		CommitteeTypeID:           internalCommitteeFormDTO.CommitteeTypeID,
		CommitteeTermID:           internalCommitteeFormDTO.CommitteeTermID,
		TermEndDate:               parsedDates.TermEndDate,
		CommitteeClassificationID: internalCommitteeFormDTO.CommitteeClassificationID,
		Justification:             internalCommitteeFormDTO.Justification,
		RolesAndResponsibilities:  internalCommitteeFormDTO.RolesAndResponsibilities,
		CommitteeAuthorities:      internalCommitteeFormDTO.CommitteeAuthorities,
		CommitteeStatusID:         draftStatus.ID,
		Version:                   &version,
		IsLatestVersion:           &isLatest,
		CreatedBYID:               currentUserID,
	}

	versionID, err := ics.repo.CreateVersion(internalCommitteeVersion, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	// Create nominee departments
	var nomineeDepartments []*models.InternalCommitteeVersionNomineeDepartment
	for _, nominee := range internalCommitteeFormDTO.NomineeDepartments {
		committeeMemberTypeID := ics.resolveNomineeCommitteeMemberTypeID(nominee)
		nomineeDepartments = append(nomineeDepartments, &models.InternalCommitteeVersionNomineeDepartment{
			InternalCommitteeVersionID: versionID,
			NomineeDepartmentID:        nominee.NomineeDepartmentID,
			CommitteeMemberTypeID:      committeeMemberTypeID,
			CreatedBYID:                &currentUserID,
		})
	}

	if err := ics.repo.CreateNomineeDepartments(nomineeDepartments, tx); err != nil {
		tx.Rollback()
		return 0, err
	}

	tx.Commit()

	// Send notifications to department heads
	go func() {

		newCommittee, err := ics.GetFullDetailsByID(committeeID)
		if err != nil {
			return
		}

		createdData := map[string]interface{}{
			"reference_number":           newCommittee.ReferenceNumber,
			"committee_name_en":          newCommittee.LatestVersion.CommitteeNameEn,
			"committee_name_ar":          newCommittee.LatestVersion.CommitteeNameAr,
			"committee_type":             newCommittee.LatestVersion.CommitteeType.TitleEn,
			"committee_term":             newCommittee.LatestVersion.CommitteeTerm.TitleEn,
			"term_end_date":              newCommittee.LatestVersion.TermEndDateString,
			"committee_classification":   newCommittee.LatestVersion.CommitteeClassification.TitleEn,
			"justification":              newCommittee.LatestVersion.Justification,
			"roles_and_responsibilities": newCommittee.LatestVersion.RolesAndResponsibilities,
			"committee_authorities":      newCommittee.LatestVersion.CommitteeAuthorities,
		}

		_ = ics.activityLogSvc.LogActivity(
			models.CreateAction,
			models.InternalCommittee{}.TableName(),
			committeeID,
			createdData,
			currentUserID,
			"organization",
			"committee_management",
			"internal_committee",
			nil,
		)

		for _, nominee := range internalCommitteeFormDTO.NomineeDepartments {
			// Get department head users directly
			headUsers, err := ics.organizationUnitsSvc.GetParentOrganizationUnitWithHead(nominee.NomineeDepartmentID)
			if err != nil {
				utility.LogError("Failed to get department head users for notification", map[string]interface{}{
					"department_id": nominee.NomineeDepartmentID,
					"error":         err.Error(),
				})
				continue
			}

			if headUsers.HeadUser == nil {
				continue
			}

			// Get committee details for notification
			committee, err := ics.GetFullDetailsByID(committeeID)
			if err != nil {
				utility.LogError("Failed to get committee details for notification", map[string]interface{}{
					"committee_id": committeeID,
					"error":        err.Error(),
				})
				continue
			}

			// Get department details
			department, err := ics.organizationUnitSvc.GetByID(nominee.NomineeDepartmentID)
			if err != nil {
				utility.LogError("Failed to get department details for notification", map[string]interface{}{
					"department_id": nominee.NomineeDepartmentID,
					"error":         err.Error(),
				})
				continue
			}

			// Create notifications for each department head
			if headUsers.HeadUser != nil {
				urlParam := fmt.Sprintf(`{"id":%d}`, committeeID)
				entity := "internal_committee"
				action := "department-nomination"

				dbNotificationDTO := &commonDTOs.DBNotificationDTO{
					UserIDs:              []uint{headUsers.HeadUser.ID},
					ModuleIdentifier:     "organization",
					NotificationType:     string(commonDTOs.AssignmentType),
					ReferenceNumber:      committee.ReferenceNumber,
					HeadingEn:            fmt.Sprintf("Department Nominated for Committee:%s", newCommittee.LatestVersion.CommitteeNameEn),
					HeadingAr:            fmt.Sprintf("تم ترشيح الإدارة للجنة:%s", newCommittee.LatestVersion.CommitteeNameAr),
					TitleEn:              fmt.Sprintf("Your department (%s) has been nominated for committee: %s", department.TitleEn, newCommittee.LatestVersion.CommitteeNameEn),
					TitleAr:              fmt.Sprintf("لقد تم ترشيح إدارتك (%s) للجنة: %s", department.TitleAr, newCommittee.LatestVersion.CommitteeNameAr),
					Entity:               &entity,
					Action:               &action,
					URLParam:             urlParam,
					IsSystemNotification: false,
					CreatedByID:          &currentUserID,
				}

				responsibleUser, err := ics.userSvc.GetByID(headUsers.HeadUser.ID)
				if err != nil {
					utility.LogError("Failed to fetch head user details for notification", map[string]interface{}{
						"error": err.Error(),
					})
					return
				}
				fmt.Println(newCommittee)
				emailTo := []string{responsibleUser.Email}
				subject := fmt.Sprintf("Committee:%s: Department Nomination", newCommittee.LatestVersion.CommitteeNameEn)
				appURL := utility.GetAppURL()
				emailData := map[string]interface{}{
					"FirstNameEn":        responsibleUser.FirstNameEn,
					"LastNameEn":         responsibleUser.LastNameEn,
					"DepartmentNameEn":   department.TitleEn,
					"DepartmentNameAr":   department.TitleAr,
					"CommitteeReference": committee.ReferenceNumber,
					"CommitteeNameEn":    newCommittee.LatestVersion.CommitteeNameEn,
					"CommitteeNameAr":    newCommittee.LatestVersion.CommitteeNameAr,
					"MembershipType":     ics.getCommitteeMemberTypeTitle(ics.resolveNomineeCommitteeMemberTypeID(nominee)),
					"NominatedBy":        fmt.Sprintf("%s %s", newCommittee.CreatedByUser.FirstNameEN, newCommittee.CreatedByUser.LastNameEN),
					"NominatedOn":        newCommittee.CreatedAt.Format("2006-01-02"),
					"RedirectURL":        fmt.Sprintf("%s/organization/committee/%d/overview", appURL, committeeID),
				}

				emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
					EmailTo:              emailTo,
					Subject:              subject,
					TemplateRelativePath: "organization/templates/email/internal_committee_nomination_mail.html",
					EmailData:            emailData,
				}

				err = ics.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "committee_management")
				if err != nil {
					utility.LogError("Failed to send notifications", map[string]interface{}{
						"error": err.Error(),
					})
				}
			}
		}
	}()

	return committeeID, nil
}

func (ics *InternalCommitteeService) Update(id uint, currentUserID uint, internalCommitteeFormDTO *organizationDTOs.InternalCommitteeFormDTO) error {
	tx := ics.repo.BeginTransaction()

	// Get existing internal committee
	existingCommittee, err := ics.repo.GetByID(id)
	if err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utility.NewResourceNotFoundError(err, id)
		}
		return err
	}

	// Get old data before update for activity logging
	oldCommittee, err := ics.GetFullDetailsByID(id, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Validate form data
	if err := ics.validateInternalCommitteeForm(internalCommitteeFormDTO); err != nil {
		tx.Rollback()
		return err
	}

	parsedDates, err := ics.parseInternalCommitteeDates(internalCommitteeFormDTO)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Update internal committee
	existingCommittee.UpdatedBYID = &currentUserID

	if err := ics.repo.Update(existingCommittee, tx); err != nil {
		tx.Rollback()
		return err
	}

	// Get latest version to update
	versions, err := ics.repo.GetVersionsByCommitteeID(id, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	var latestVersion *models.InternalCommitteeVersion
	for _, v := range versions {
		if v.IsLatestVersion != nil && *v.IsLatestVersion {
			latestVersion = v
			break
		}
	}

	if latestVersion == nil {
		tx.Rollback()
		return fmt.Errorf("no latest version found for committee")
	}

	// Update existing latest version instead of creating new one
	latestVersion.CommitteeNameEn = internalCommitteeFormDTO.CommitteeNameEn
	latestVersion.CommitteeNameAr = internalCommitteeFormDTO.CommitteeNameAr
	latestVersion.CommitteeTypeID = internalCommitteeFormDTO.CommitteeTypeID
	latestVersion.CommitteeTermID = internalCommitteeFormDTO.CommitteeTermID
	latestVersion.TermEndDate = parsedDates.TermEndDate
	latestVersion.CommitteeClassificationID = internalCommitteeFormDTO.CommitteeClassificationID
	latestVersion.Justification = internalCommitteeFormDTO.Justification
	latestVersion.RolesAndResponsibilities = internalCommitteeFormDTO.RolesAndResponsibilities
	latestVersion.CommitteeAuthorities = internalCommitteeFormDTO.CommitteeAuthorities
	latestVersion.UpdatedBYID = &currentUserID

	if err := ics.repo.UpdateVersion(latestVersion, tx); err != nil {
		tx.Rollback()
		return err
	}

	// Delete existing nominee departments and create new ones
	if err := ics.repo.DeleteNomineeDepartmentsByVersionID(latestVersion.ID, tx); err != nil {
		tx.Rollback()
		return err
	}

	// Create new nominee departments
	var nomineeDepartments []*models.InternalCommitteeVersionNomineeDepartment
	for _, nominee := range internalCommitteeFormDTO.NomineeDepartments {
		committeeMemberTypeID := ics.resolveNomineeCommitteeMemberTypeID(nominee)
		nomineeDepartments = append(nomineeDepartments, &models.InternalCommitteeVersionNomineeDepartment{
			InternalCommitteeVersionID: latestVersion.ID,
			NomineeDepartmentID:        nominee.NomineeDepartmentID,
			CommitteeMemberTypeID:      committeeMemberTypeID,
			CreatedBYID:                &currentUserID,
		})
	}

	if err := ics.repo.CreateNomineeDepartments(nomineeDepartments, tx); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	// Send notifications to department heads
	go func() {

		newCommittee, err := ics.GetFullDetailsByID(id)
		if err != nil {
			return
		}

		oldData := map[string]interface{}{
			"committee_name_en":          oldCommittee.LatestVersion.CommitteeNameEn,
			"committee_name_ar":          oldCommittee.LatestVersion.CommitteeNameAr,
			"committee_type":             oldCommittee.LatestVersion.CommitteeType.TitleEn,
			"committee_term":             oldCommittee.LatestVersion.CommitteeTerm.TitleEn,
			"term_end_date":              oldCommittee.LatestVersion.TermEndDateString,
			"committee_classification":   oldCommittee.LatestVersion.CommitteeClassification.TitleEn,
			"justification":              oldCommittee.LatestVersion.Justification,
			"roles_and_responsibilities": oldCommittee.LatestVersion.RolesAndResponsibilities,
			"committee_authorities":      oldCommittee.LatestVersion.CommitteeAuthorities,
		}

		newData := map[string]interface{}{
			"committee_name_en":          newCommittee.LatestVersion.CommitteeNameEn,
			"committee_name_ar":          newCommittee.LatestVersion.CommitteeNameAr,
			"committee_type":             newCommittee.LatestVersion.CommitteeType.TitleEn,
			"committee_term":             newCommittee.LatestVersion.CommitteeTerm.TitleEn,
			"term_end_date":              newCommittee.LatestVersion.TermEndDateString,
			"committee_classification":   newCommittee.LatestVersion.CommitteeClassification.TitleEn,
			"justification":              newCommittee.LatestVersion.Justification,
			"roles_and_responsibilities": newCommittee.LatestVersion.RolesAndResponsibilities,
			"committee_authorities":      newCommittee.LatestVersion.CommitteeAuthorities,
		}

		changedData := utility.CompareFieldChanges(oldData, newData)

		if len(changedData) > 0 {
			_ = ics.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.InternalCommittee{}.TableName(),
				id,
				changedData,
				currentUserID,
				"organization",
				"committee_management",
				"internal_committee",
				nil,
			)
		}

		// Compare old and new nominee departments to identify deleted ones
		oldDepartmentIDs := make(map[uint]bool)
		for _, oldNominee := range oldCommittee.LatestVersion.NomineeDepartments {
			oldDepartmentIDs[oldNominee.NomineeDepartmentID] = true
		}

		newDepartmentIDs := make(map[uint]bool)
		for _, newNominee := range internalCommitteeFormDTO.NomineeDepartments {
			newDepartmentIDs[newNominee.NomineeDepartmentID] = true
		}

		// Send notifications for deleted departments
		for _, oldNominee := range oldCommittee.LatestVersion.NomineeDepartments {
			if !newDepartmentIDs[oldNominee.NomineeDepartmentID] {
				// This department was removed
				headUsers, err := ics.organizationUnitsSvc.GetParentOrganizationUnitWithHead(oldNominee.NomineeDepartmentID)
				if err != nil {
					utility.LogError("Failed to get department head users for deletion notification", map[string]interface{}{
						"department_id": oldNominee.NomineeDepartmentID,
						"error":         err.Error(),
					})
					continue
				}

				if headUsers.HeadUser == nil {
					continue
				}

				// Get committee details for notification
				committee, err := ics.GetByID(id)
				if err != nil {
					utility.LogError("Failed to get committee details for deletion notification", map[string]interface{}{
						"committee_id": id,
						"error":        err.Error(),
					})
					continue
				}

				// Get department details
				department, err := ics.organizationUnitSvc.GetByID(oldNominee.NomineeDepartmentID)
				if err != nil {
					utility.LogError("Failed to get department details for deletion notification", map[string]interface{}{
						"department_id": oldNominee.NomineeDepartmentID,
						"error":         err.Error(),
					})
					continue
				}

				// Create notifications for each department head about removal
				if headUsers.HeadUser != nil {
					urlParam := fmt.Sprintf(`{"id":%d}`, id)
					entity := "internal_committee"
					action := "department-nomination-removed"

					dbNotificationDTO := &commonDTOs.DBNotificationDTO{
						UserIDs:              []uint{headUsers.HeadUser.ID},
						ModuleIdentifier:     "organization",
						NotificationType:     string(commonDTOs.InformationType),
						ReferenceNumber:      committee.ReferenceNumber,
						HeadingEn:            fmt.Sprintf("Department Nominated for Removed from Committee:%s", newCommittee.LatestVersion.CommitteeNameEn),
						HeadingAr:            fmt.Sprintf("تم إزالة ترشيح الإدارة من اللجنة:%s", newCommittee.LatestVersion.CommitteeNameAr),
						TitleEn:              fmt.Sprintf("Your department (%s) has been removed from committee: %s", department.TitleEn, newCommittee.LatestVersion.CommitteeNameEn),
						TitleAr:              fmt.Sprintf("لقد تم إزالة ترشيح إدارتك (%s) من اللجنة: %s", department.TitleAr, newCommittee.LatestVersion.CommitteeNameAr),
						Entity:               &entity,
						Action:               &action,
						URLParam:             urlParam,
						IsSystemNotification: false,
						CreatedByID:          &currentUserID,
					}

					responsibleUser, err := ics.userSvc.GetByID(headUsers.HeadUser.ID)
					if err != nil {
						utility.LogError("Failed to fetch head user details for deletion notification", map[string]interface{}{
							"error": err.Error(),
						})
						continue
					}

					emailTo := []string{responsibleUser.Email}
					subject := fmt.Sprintf("Committee:%s: Department Removed", newCommittee.LatestVersion.CommitteeNameEn)
					appURL := utility.GetAppURL()
					emailData := map[string]interface{}{
						"FirstNameEn":        responsibleUser.FirstNameEn,
						"LastNameEn":         responsibleUser.LastNameEn,
						"DepartmentNameEn":   department.TitleEn,
						"DepartmentNameAr":   department.TitleAr,
						"CommitteeReference": committee.ReferenceNumber,
						"CommitteeNameEn":    newCommittee.LatestVersion.CommitteeNameEn,
						"CommitteeNameAr":    newCommittee.LatestVersion.CommitteeNameAr,
						"MembershipType":     ics.getCommitteeMemberTypeTitle(oldNominee.CommitteeMemberTypeID),
						"RemovedBy":          fmt.Sprintf("%s %s", newCommittee.UpdatedByUser.FirstNameEN, newCommittee.UpdatedByUser.LastNameEN),
						"RemovedOn":          newCommittee.UpdatedAt.Format("2006-01-02"),
						"RedirectURL":        fmt.Sprintf("%s/organization/committee/%d/overview", appURL, id),
					}

					emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
						EmailTo:              emailTo,
						Subject:              subject,
						TemplateRelativePath: "organization/templates/email/internal_committee_nomination_removed_mail.html",
						EmailData:            emailData,
					}

					err = ics.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "committee_management")
					if err != nil {
						utility.LogError("Failed to send deletion notifications", map[string]interface{}{
							"error": err.Error(),
						})
					}
				}
			}
		}

		// Send notifications for new departments (use existing create template)
		for _, nominee := range internalCommitteeFormDTO.NomineeDepartments {
			if !oldDepartmentIDs[nominee.NomineeDepartmentID] {
				// This is a new department
				// Get department head users directly
				headUsers, err := ics.organizationUnitsSvc.GetParentOrganizationUnitWithHead(nominee.NomineeDepartmentID)
				if err != nil {
					utility.LogError("Failed to get department head users for notification", map[string]interface{}{
						"department_id": nominee.NomineeDepartmentID,
						"error":         err.Error(),
					})
					continue
				}

				if headUsers.HeadUser == nil {
					continue
				}

				// Get committee details for notification
				committee, err := ics.GetByID(id)
				if err != nil {
					utility.LogError("Failed to get committee details for notification", map[string]interface{}{
						"committee_id": id,
						"error":        err.Error(),
					})
					continue
				}

				// Get department details
				department, err := ics.organizationUnitSvc.GetByID(nominee.NomineeDepartmentID)
				if err != nil {
					utility.LogError("Failed to get department details for notification", map[string]interface{}{
						"department_id": nominee.NomineeDepartmentID,
						"error":         err.Error(),
					})
					continue
				}

				// Create notifications for each department head (new department - use create template)
				if headUsers.HeadUser != nil {
					urlParam := fmt.Sprintf(`{"id":%d}`, id)
					entity := "internal_committee"
					action := "department-nomination"

					dbNotificationDTO := &commonDTOs.DBNotificationDTO{
						UserIDs:              []uint{headUsers.HeadUser.ID},
						ModuleIdentifier:     "organization",
						NotificationType:     string(commonDTOs.AssignmentType),
						ReferenceNumber:      committee.ReferenceNumber,
						HeadingEn:            fmt.Sprintf("Department Nominated for Committee:%s", newCommittee.LatestVersion.CommitteeNameEn),
						HeadingAr:            fmt.Sprintf("تم ترشيح الإدارة للجنة:%s", newCommittee.LatestVersion.CommitteeNameAr),
						TitleEn:              fmt.Sprintf("Your department (%s) has been nominated for committee: %s", department.TitleEn, newCommittee.LatestVersion.CommitteeNameEn),
						TitleAr:              fmt.Sprintf("لقد تم ترشيح إدارتك (%s) للجنة: %s", department.TitleAr, newCommittee.LatestVersion.CommitteeNameAr),
						Entity:               &entity,
						Action:               &action,
						URLParam:             urlParam,
						IsSystemNotification: false,
						CreatedByID:          &currentUserID,
					}

					responsibleUser, err := ics.userSvc.GetByID(headUsers.HeadUser.ID)
					if err != nil {
						utility.LogError("Failed to fetch head user details for notification", map[string]interface{}{
							"error": err.Error(),
						})
						continue
					}

					emailTo := []string{responsibleUser.Email}
					subject := fmt.Sprintf("Committee:%s: Department Nomination", newCommittee.LatestVersion.CommitteeNameEn)
					appURL := utility.GetAppURL()
					emailData := map[string]interface{}{
						"FirstNameEn":        responsibleUser.FirstNameEn,
						"LastNameEn":         responsibleUser.LastNameEn,
						"DepartmentNameEn":   department.TitleEn,
						"DepartmentNameAr":   department.TitleAr,
						"CommitteeReference": committee.ReferenceNumber,
						"CommitteeNameEn":    newCommittee.LatestVersion.CommitteeNameEn,
						"CommitteeNameAr":    newCommittee.LatestVersion.CommitteeNameAr,
						"MembershipType":     ics.getCommitteeMemberTypeTitle(ics.resolveNomineeCommitteeMemberTypeID(nominee)),
						"NominatedBy":        fmt.Sprintf("%s %s", newCommittee.UpdatedByUser.FirstNameEN, newCommittee.UpdatedByUser.LastNameEN),
						"NominatedOn":        newCommittee.UpdatedAt.Format("2006-01-02"),
						"RedirectURL":        fmt.Sprintf("%s/organization/committee/%d/overview", appURL, id),
					}

					emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
						EmailTo:              emailTo,
						Subject:              subject,
						TemplateRelativePath: "organization/templates/email/internal_committee_nomination_mail.html",
						EmailData:            emailData,
					}

					err = ics.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "committee_management")
					if err != nil {
						utility.LogError("Failed to send notifications", map[string]interface{}{
							"error": err.Error(),
						})
					}
				}
			}
		}
	}()

	return nil
}

func (ics *InternalCommitteeService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	tx := ics.repo.BeginTransaction()

	var failedIds []uint
	var successfulIds []uint
	var deletedCommittees []map[string]interface{}

	for _, id := range ids {

		// Get committee details
		committee, err := ics.GetFullDetailsByID(id, tx)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		deletedCommittees = append(deletedCommittees, map[string]interface{}{
			"id":               id,
			"reference_number": committee.ReferenceNumber,
			"committee_name":   committee.LatestVersion.CommitteeNameEn,
		})

		// Get versions
		versions, err := ics.repo.GetVersionsByCommitteeID(id, tx)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		// Delete nominee departments and their participants
		for _, version := range versions {
			// Get nominee departments for this version
			nomineeDepartments, err := ics.repo.GetNomineeDepartmentsByVersionID(version.ID, tx)
			if err != nil {
				tx.Rollback()
				return nil, err
			}

			// Delete participants for each nominee department
			for _, nomineeDept := range nomineeDepartments {
				if err := ics.repo.DeleteParticipantsByNomineeDepartmentID(nomineeDept.ID, tx); err != nil {
					tx.Rollback()
					return nil, err
				}
			}

			// Now delete the nominee departments
			if err := ics.repo.DeleteNomineeDepartmentsByVersionID(version.ID, tx); err != nil {
				tx.Rollback()
				return nil, err
			}
		}

		// Delete versions
		for _, version := range versions {
			if err := ics.repo.DeleteVersion(version.ID, tx); err != nil {
				tx.Rollback()
				return nil, err
			}
		}

		// Delete main committee
		if err := ics.repo.Delete(id, tx); err != nil {
			failedIds = append(failedIds, id)
			continue
		} else {
			successfulIds = append(successfulIds, id)
		}
	}

	// If any failed → rollback everything
	if len(failedIds) > 0 {
		tx.Rollback()
		return failedIds, nil
	}

	tx.Commit()

	// Log bulk delete activities
	if len(successfulIds) > 0 {
		go func() {
			err := ics.activityLogSvc.LogBulkActivities(
				models.DeleteAction,
				models.InternalCommittee{}.TableName(),
				successfulIds,
				deletedCommittees,
				currentUserID,
				nil,
			)
			if err != nil {
				utility.LogError("Failed to log bulk delete activity", map[string]interface{}{
					"error": err.Error(),
				})
			}
		}()
	}

	return failedIds, nil
}

func (ics *InternalCommitteeService) ReformCommittee(currentUserID uint, reformDTO *organizationDTOs.ReformInternalCommitteeDTO) (uint, error) {
	// Validate form data
	formDTO := &organizationDTOs.InternalCommitteeFormDTO{
		CommitteeNameEn:           reformDTO.CommitteeNameEn,
		CommitteeNameAr:           reformDTO.CommitteeNameAr,
		CommitteeTypeID:           reformDTO.CommitteeTypeID,
		CommitteeTermID:           reformDTO.CommitteeTermID,
		TermEndDate:               reformDTO.TermEndDate,
		CommitteeClassificationID: reformDTO.CommitteeClassificationID,
		Justification:             reformDTO.Justification,
		RolesAndResponsibilities:  reformDTO.RolesAndResponsibilities,
		CommitteeAuthorities:      reformDTO.CommitteeAuthorities,
		NomineeDepartments:        reformDTO.NomineeDepartments,
	}

	if err := ics.validateInternalCommitteeForm(formDTO); err != nil {
		return 0, err
	}

	parsedDates, err := ics.parseInternalCommitteeDates(formDTO)
	if err != nil {
		return 0, err
	}

	tx := ics.repo.BeginTransaction()

	// Get existing committee
	existingCommittee, err := ics.repo.GetByID(reformDTO.InternalCommitteeID)
	if err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, utility.NewResourceNotFoundError(err, reformDTO.InternalCommitteeID)
		}
		return 0, err
	}

	// Update committee
	existingCommittee.UpdatedBYID = &currentUserID
	if err := ics.repo.Update(existingCommittee, tx); err != nil {
		tx.Rollback()
		return 0, err
	}

	// Get existing versions
	versions, err := ics.repo.GetVersionsByCommitteeID(reformDTO.InternalCommitteeID, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	// Find latest version and get its nominee departments
	var latestVersion *models.InternalCommitteeVersion
	var oldNomineeDepartments []*models.InternalCommitteeVersionNomineeDepartment
	for _, v := range versions {
		if v.IsLatestVersion != nil && *v.IsLatestVersion {
			latestVersion = v
			oldNomineeDepartments, err = ics.repo.GetNomineeDepartmentsByVersionID(v.ID, tx)
			if err != nil {
				tx.Rollback()
				return 0, err
			}
			break
		}
	}

	if latestVersion == nil {
		tx.Rollback()
		return 0, fmt.Errorf("no latest version found for committee")
	}

	// Set old version's is_latest to false
	isLatestFalse := false
	latestVersion.IsLatestVersion = &isLatestFalse
	latestVersion.UpdatedBYID = &currentUserID
	if err := ics.repo.UpdateVersion(latestVersion, tx); err != nil {
		tx.Rollback()
		return 0, err
	}

	// Get reformation status
	reformationStatus, err := ics.internalCommitteeStatusSvc.GetByType("reformation")
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	// Calculate new version number
	newVersionNumber := *latestVersion.Version + 1
	isLatestTrue := true

	// Create new version
	newVersion := &models.InternalCommitteeVersion{
		InternalCommitteeID:       reformDTO.InternalCommitteeID,
		CommitteeNameEn:           reformDTO.CommitteeNameEn,
		CommitteeNameAr:           reformDTO.CommitteeNameAr,
		CommitteeTypeID:           reformDTO.CommitteeTypeID,
		CommitteeTermID:           reformDTO.CommitteeTermID,
		TermEndDate:               parsedDates.TermEndDate,
		CommitteeClassificationID: reformDTO.CommitteeClassificationID,
		Justification:             reformDTO.Justification,
		RolesAndResponsibilities:  reformDTO.RolesAndResponsibilities,
		CommitteeAuthorities:      reformDTO.CommitteeAuthorities,
		CommitteeStatusID:         reformationStatus.ID,
		Version:                   &newVersionNumber,
		IsLatestVersion:           &isLatestTrue,
		CreatedBYID:               currentUserID,
	}

	newVersionID, err := ics.repo.CreateVersion(newVersion, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	// Create map of old departments for quick lookup
	oldDepartmentsMap := make(map[uint][]*models.InternalCommitteeVersionNomineeDepartmentParticipant)
	for _, oldDept := range oldNomineeDepartments {
		participants, err := ics.repo.GetParticipantsByNomineeDepartmentID(oldDept.ID, tx)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
		oldDepartmentsMap[oldDept.NomineeDepartmentID] = participants
	}

	// Create new nominee departments and copy participants if department matches
	for _, newDeptDTO := range reformDTO.NomineeDepartments {
		committeeMemberTypeID := ics.resolveNomineeCommitteeMemberTypeID(newDeptDTO)
		newDept := &models.InternalCommitteeVersionNomineeDepartment{
			InternalCommitteeVersionID: newVersionID,
			NomineeDepartmentID:        newDeptDTO.NomineeDepartmentID,
			CommitteeMemberTypeID:      committeeMemberTypeID,
			CreatedBYID:                &currentUserID,
		}

		// Create nominee department
		var nomineeDepartments []*models.InternalCommitteeVersionNomineeDepartment
		nomineeDepartments = append(nomineeDepartments, newDept)
		if err := ics.repo.CreateNomineeDepartments(nomineeDepartments, tx); err != nil {
			tx.Rollback()
			return 0, err
		}

		// Check if this department existed in old version
		if oldParticipants, exists := oldDepartmentsMap[newDeptDTO.NomineeDepartmentID]; exists && len(oldParticipants) > 0 {
			// Copy participants from old department to new department
			for _, oldParticipant := range oldParticipants {
				newParticipant := &models.InternalCommitteeVersionNomineeDepartmentParticipant{
					InternalCommitteeVersionNomineeDepartmentID: newDept.ID,
					UserID:       oldParticipant.UserID,
					DepartmentID: oldParticipant.DepartmentID,
					CreatedBYID:  &currentUserID,
				}

				if _, err := ics.repo.CreateParticipant(newParticipant, tx); err != nil {
					tx.Rollback()
					return 0, err
				}
			}
		}
	}

	tx.Commit()

	// Log activity asynchronously
	go func() {
		newCommittee, err := ics.GetFullDetailsByID(reformDTO.InternalCommitteeID)
		if err != nil {
			return
		}

		reformData := map[string]interface{}{
			"reference_number":           newCommittee.ReferenceNumber,
			"committee_name_en":          newCommittee.LatestVersion.CommitteeNameEn,
			"committee_name_ar":          newCommittee.LatestVersion.CommitteeNameAr,
			"version":                    newCommittee.LatestVersion.Version,
			"committee_status":           "reformation",
			"committee_type":             newCommittee.LatestVersion.CommitteeType.TitleEn,
			"committee_term":             newCommittee.LatestVersion.CommitteeTerm.TitleEn,
			"term_end_date":              newCommittee.LatestVersion.TermEndDateString,
			"committee_classification":   newCommittee.LatestVersion.CommitteeClassification.TitleEn,
			"justification":              newCommittee.LatestVersion.Justification,
			"roles_and_responsibilities": newCommittee.LatestVersion.RolesAndResponsibilities,
			"committee_authorities":      newCommittee.LatestVersion.CommitteeAuthorities,
		}

		_ = ics.activityLogSvc.LogActivity(
			models.CreateAction,
			models.InternalCommittee{}.TableName(),
			reformDTO.InternalCommitteeID,
			reformData,
			currentUserID,
			"organization",
			"committee_management",
			"reform_internal_committee",
			nil,
		)

		for _, nominee := range reformDTO.NomineeDepartments {
			// Get department head users directly
			headUsers, err := ics.organizationUnitsSvc.GetParentOrganizationUnitWithHead(nominee.NomineeDepartmentID)
			if err != nil {
				utility.LogError("Failed to get department head users for notification", map[string]interface{}{
					"department_id": nominee.NomineeDepartmentID,
					"error":         err.Error(),
				})
				continue
			}

			if headUsers.HeadUser == nil {
				continue
			}

			// Get committee details for notification
			committee, err := ics.GetByID(newCommittee.ID)
			if err != nil {
				utility.LogError("Failed to get committee details for notification", map[string]interface{}{
					"committee_id": newCommittee.ID,
					"error":        err.Error(),
				})
				continue
			}

			// Get department details
			department, err := ics.organizationUnitSvc.GetByID(nominee.NomineeDepartmentID)
			if err != nil {
				utility.LogError("Failed to get department details for notification", map[string]interface{}{
					"department_id": nominee.NomineeDepartmentID,
					"error":         err.Error(),
				})
				continue
			}

			// Create notifications for each department head
			if headUsers.HeadUser != nil {
				urlParam := fmt.Sprintf(`{"id":%d}`, newCommittee.ID)
				entity := "internal_committee"
				action := "department-nomination"

				dbNotificationDTO := &commonDTOs.DBNotificationDTO{
					UserIDs:              []uint{headUsers.HeadUser.ID},
					ModuleIdentifier:     "organization",
					NotificationType:     string(commonDTOs.AssignmentType),
					ReferenceNumber:      committee.ReferenceNumber,
					HeadingEn:            fmt.Sprintf("Department Nominated for Committee:%s", newCommittee.LatestVersion.CommitteeNameEn),
					HeadingAr:            fmt.Sprintf("تم ترشيح الإدارة للجنة:%s", newCommittee.LatestVersion.CommitteeNameAr),
					TitleEn:              fmt.Sprintf("Your department (%s) has been nominated for committee: %s", department.TitleEn, newCommittee.LatestVersion.CommitteeNameEn),
					TitleAr:              fmt.Sprintf("لقد تم ترشيح إدارتك (%s) للجنة: %s", department.TitleAr, newCommittee.LatestVersion.CommitteeNameAr),
					Entity:               &entity,
					Action:               &action,
					URLParam:             urlParam,
					IsSystemNotification: false,
					CreatedByID:          &currentUserID,
				}

				responsibleUser, err := ics.userSvc.GetByID(headUsers.HeadUser.ID)
				if err != nil {
					utility.LogError("Failed to fetch head user details for notification", map[string]interface{}{
						"error": err.Error(),
					})
					return
				}

				emailTo := []string{responsibleUser.Email}
				subject := fmt.Sprintf("Committee/%s: Department Nomination", committee.ReferenceNumber)
				appURL := utility.GetAppURL()
				emailData := map[string]interface{}{
					"FirstNameEn":        responsibleUser.FirstNameEn,
					"LastNameEn":         responsibleUser.LastNameEn,
					"DepartmentNameEn":   department.TitleEn,
					"DepartmentNameAr":   department.TitleAr,
					"CommitteeReference": committee.ReferenceNumber,
					"CommitteeNameEn":    newCommittee.LatestVersion.CommitteeNameEn,
					"CommitteeNameAr":    newCommittee.LatestVersion.CommitteeNameAr,
					"MembershipType":     ics.getCommitteeMemberTypeTitle(ics.resolveNomineeCommitteeMemberTypeID(nominee)),
					"NominatedBy":        fmt.Sprintf("%s %s", newCommittee.CreatedByUser.FirstNameEN, newCommittee.CreatedByUser.LastNameEN),
					"NominatedOn":        newCommittee.CreatedAt.Format("2006-01-02"),
					"RedirectURL":        fmt.Sprintf("%s/organization/committee/%d/overview", appURL, newCommittee.ID),
				}

				emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
					EmailTo:              emailTo,
					Subject:              subject,
					TemplateRelativePath: "organization/templates/email/internal_committee_nomination_mail.html",
					EmailData:            emailData,
				}

				err = ics.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "committee_management")
				if err != nil {
					utility.LogError("Failed to send notifications", map[string]interface{}{
						"error": err.Error(),
					})
				}
			}
		}
	}()

	return newVersionID, nil
}

func (ics *InternalCommitteeService) validateInternalCommitteeForm(internalCommitteeFormDTO *organizationDTOs.InternalCommitteeFormDTO) error {
	// Validate committee type
	if _, err := ics.committeeTypeSvc.GetByID(internalCommitteeFormDTO.CommitteeTypeID); err != nil {
		return utility.NewValidationError(err, "not_exist", "committee_type")
	}

	// Validate committee term
	committeeTerm, err := ics.committeeTermSvc.GetByID(internalCommitteeFormDTO.CommitteeTermID)
	if err != nil {
		return utility.NewValidationError(err, "not_exist", "committee_term")
	}

	// If committee term is temporary, TermEndDate is required
	if committeeTerm.Type == "temporary" {
		if internalCommitteeFormDTO.TermEndDate == nil || *internalCommitteeFormDTO.TermEndDate == "" {
			return utility.NewValidationError(fmt.Errorf("term end date is required for temporary committee term"), "required", "term_end_date")
		}
	}

	// Validate committee classification
	if _, err := ics.committeeClassificationSvc.GetByID(internalCommitteeFormDTO.CommitteeClassificationID); err != nil {
		return utility.NewValidationError(err, "not_exist", "committee_classification")
	}

	// Validate nominee departments
	for _, nominee := range internalCommitteeFormDTO.NomineeDepartments {
		if _, err := ics.organizationUnitSvc.GetByID(nominee.NomineeDepartmentID); err != nil {
			return utility.NewValidationError(err, "not_exist", "nominee_department")
		}

		committeeMemberTypeID := ics.resolveNomineeCommitteeMemberTypeID(nominee)
		if committeeMemberTypeID != nil {
			if _, err := ics.committeeMemberTypeSvc.GetByID(*committeeMemberTypeID); err != nil {
				return utility.NewValidationError(err, "not_exist", "committee_member_type")
			}
		}
	}

	return nil
}

func (ics *InternalCommitteeService) getCommitteeMemberTypeTitle(committeeMemberTypeID *uint) string {
	if committeeMemberTypeID == nil {
		return ""
	}

	committeeMemberType, err := ics.committeeMemberTypeSvc.GetByID(*committeeMemberTypeID)
	if err != nil || committeeMemberType == nil {
		return ""
	}

	return committeeMemberType.TitleEn
}

func (ics *InternalCommitteeService) resolveNomineeCommitteeMemberTypeID(nominee organizationDTOs.InternalCommitteeVersionNomineeDepartmentFormDTO) *uint {
	return nominee.CommitteeMemberTypeID
}

func (ics *InternalCommitteeService) getCommitteeMemberTypeTitleFromNomineeDepartment(nomineeDepartment *models.InternalCommitteeVersionNomineeDepartment) string {
	if nomineeDepartment == nil {
		return ""
	}

	if nomineeDepartment.CommitteeMemberType != nil {
		return nomineeDepartment.CommitteeMemberType.TitleEn
	}

	return ics.getCommitteeMemberTypeTitle(nomineeDepartment.CommitteeMemberTypeID)
}

func (ics *InternalCommitteeService) parseInternalCommitteeDates(internalCommitteeFormDTO *organizationDTOs.InternalCommitteeFormDTO) (*organizationDTOs.InternalCommitteeParsedDates, error) {
	parsed := &organizationDTOs.InternalCommitteeParsedDates{}

	var err error

	if internalCommitteeFormDTO.TermEndDate != nil && *internalCommitteeFormDTO.TermEndDate != "" {
		parsed.TermEndDate, err = utility.ParseDate(*internalCommitteeFormDTO.TermEndDate)
		if err != nil {
			return nil, utility.NewValidationError(
				fmt.Errorf("invalid term end date"),
				"invalid_date",
				"term_end_date",
			)
		}
	}

	return parsed, nil
}

// Participant Methods

func (ics *InternalCommitteeService) CreateParticipants(currentUserID uint, nomineeDepartmentID uint, dto *organizationDTOs.InternalCommitteeParticipantFormDTO) ([]uint, error) {
	// Validate users exist
	for _, userID := range dto.UserIDs {
		if _, err := ics.userSvc.GetByID(userID); err != nil {
			return nil, utility.NewValidationError(err, "not_exist", "user")
		}
	}

	// Fetch nominee department to get NomineeDepartmentID
	nomineeDepartment, err := ics.repo.GetNomineeDepartmentByID(nomineeDepartmentID)
	if err != nil {
		return nil, utility.NewResourceNotFoundError(err, nomineeDepartmentID)
	}

	tx := ics.repo.BeginTransaction()

	// Check if any participants exist for this version before deletion
	hasParticipants, err := ics.repo.HasParticipantsByVersionID(nomineeDepartment.InternalCommitteeVersionID, tx)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// Delete existing participants for this nominee department
	if err := ics.repo.DeleteParticipantsByNomineeDepartmentID(nomineeDepartmentID, tx); err != nil {
		tx.Rollback()
		return nil, err
	}

	// Create new participants
	var createdIDs []uint
	for _, userID := range dto.UserIDs {
		participant := &models.InternalCommitteeVersionNomineeDepartmentParticipant{
			InternalCommitteeVersionNomineeDepartmentID: nomineeDepartmentID,
			UserID:       userID,
			DepartmentID: nomineeDepartment.NomineeDepartmentID, // Set NomineeDepartmentID to DepartmentID
			CreatedBYID:  &currentUserID,
		}

		id, err := ics.repo.CreateParticipant(participant, tx)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		createdIDs = append(createdIDs, id)
	}

	// Update status to under_nomination only if this is the first time participants are being added
	if !hasParticipants {
		underNominationStatus, err := ics.internalCommitteeStatusSvc.GetByType("under_nomination")
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		err = ics.repo.UpdateVersionStatus(nomineeDepartment.InternalCommitteeVersionID, underNominationStatus.ID, currentUserID, tx)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	tx.Commit()

	// Send notifications to assigned participants
	go func() {
		// Get participants to access nominee department relationship
		participants, err := ics.repo.GetParticipantsByNomineeDepartmentID(nomineeDepartmentID)
		if err != nil || len(participants) == 0 {
			utility.LogError("Failed to get participants for committee details", map[string]interface{}{
				"nominee_department_id": nomineeDepartmentID,
				"error":                 err.Error(),
			})
			return
		}

		// Use first participant to get committee details (all participants belong to same committee)
		firstParticipant := participants[0]
		if firstParticipant.NomineeDepartment == nil || firstParticipant.NomineeDepartment.InternalCommitteeVersion == nil {
			utility.LogError("Missing nominee department or version relationship", map[string]interface{}{
				"nominee_department_id": nomineeDepartmentID,
			})
			return
		}

		committeeID := firstParticipant.NomineeDepartment.InternalCommitteeVersion.InternalCommitteeID

		// Get full committee details
		committee, err := ics.GetFullDetailsByID(committeeID)
		if err != nil {
			utility.LogError("Failed to get committee details for notification", map[string]interface{}{
				"committee_id": committeeID,
				"error":        err.Error(),
			})
			return
		}

		// Get department details
		department, err := ics.organizationUnitSvc.GetByID(firstParticipant.NomineeDepartment.NomineeDepartmentID)
		if err != nil {
			utility.LogError("Failed to get department details for notification", map[string]interface{}{
				"department_id": firstParticipant.NomineeDepartment.NomineeDepartmentID,
				"error":         err.Error(),
			})
			return
		}

		// Get department head (current user) who is assigning the participants
		departmentHead, err := ics.userSvc.GetByID(currentUserID)
		if err != nil {
			utility.LogError("Failed to get department head details for notification", map[string]interface{}{
				"user_id": currentUserID,
				"error":   err.Error(),
			})
			return
		}

		// Send notifications to each assigned participant
		for _, userID := range dto.UserIDs {
			// Get user details
			user, err := ics.userSvc.GetByID(userID)
			if err != nil {
				utility.LogError("Failed to get user details for notification", map[string]interface{}{
					"user_id": userID,
					"error":   err.Error(),
				})
				continue
			}

			urlParam := fmt.Sprintf(`{"id":%d}`, committeeID)
			entity := "internal_committee"
			action := "participant-assigned"

			dbNotificationDTO := &commonDTOs.DBNotificationDTO{
				UserIDs:              []uint{userID},
				ModuleIdentifier:     "organization",
				NotificationType:     string(commonDTOs.AssignmentType),
				ReferenceNumber:      committee.ReferenceNumber,
				HeadingEn:            fmt.Sprintf("Department Nominated for Committee:%s", committee.LatestVersion.CommitteeNameEn),
				HeadingAr:            fmt.Sprintf("تم ترشيح الإدارة للجنة:%s", committee.LatestVersion.CommitteeNameAr),
				TitleEn:              fmt.Sprintf("Your department (%s) has been nominated for committee: %s", department.TitleEn, committee.LatestVersion.CommitteeNameEn),
				TitleAr:              fmt.Sprintf("لقد تم ترشيح إدارتك (%s) للجنة: %s", department.TitleAr, committee.LatestVersion.CommitteeNameAr),
				Entity:               &entity,
				Action:               &action,
				URLParam:             urlParam,
				IsSystemNotification: false,
				CreatedByID:          &currentUserID,
			}

			emailTo := []string{user.Email}
			subject := fmt.Sprintf("Committee/%s: Participant Assignment", committee.LatestVersion.CommitteeNameEn)
			appURL := utility.GetAppURL()
			emailData := map[string]interface{}{
				"FirstNameEn":        user.FirstNameEn,
				"LastNameEn":         user.LastNameEn,
				"CommitteeReference": committee.ReferenceNumber,
				"CommitteeNameEn":    committee.LatestVersion.CommitteeNameEn,
				"CommitteeNameAr":    committee.LatestVersion.CommitteeNameAr,
				"DepartmentNameEn":   department.TitleEn,
				"DepartmentNameAr":   department.TitleAr,
				"MembershipType":     ics.getCommitteeMemberTypeTitleFromNomineeDepartment(firstParticipant.NomineeDepartment),
				"AssignedBy":         fmt.Sprintf("%s %s", departmentHead.FirstNameEn, departmentHead.LastNameEn),
				"AssignedOn":         time.Now().Format("2006-01-02"),
				"RedirectURL":        fmt.Sprintf("%s/organization/committee/%d/overview", appURL, committeeID),
			}

			emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
				EmailTo:              emailTo,
				Subject:              subject,
				TemplateRelativePath: "organization/templates/email/internal_committee_participant_assigned_mail.html",
				EmailData:            emailData,
			}

			err = ics.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "committee_management")
			if err != nil {
				utility.LogError("Failed to send participant assignment notification", map[string]interface{}{
					"error":   err.Error(),
					"user_id": userID,
				})
			}
		}

		// Log activity asynchronously
		for _, id := range createdIDs {
			for _, userID := range dto.UserIDs {
				err := ics.activityLogSvc.LogActivity(
					models.CreateAction,
					"committee_participants",
					id,
					map[string]interface{}{
						"nominee_department_id": nomineeDepartmentID,
						"user_id":               userID,
					},
					currentUserID,
					"organization",
					"committee_management",
					"create committee participant",
					nil,
				)
				if err != nil {
					utility.LogError("Failed to log activity", map[string]interface{}{
						"error": err.Error(),
					})
				}
			}
		}
	}()

	return createdIDs, nil
}

func (ics *InternalCommitteeService) DeleteParticipant(nomineeDepartmentID uint, participantID uint, currentUserID uint) error {
	// Get participant details before deletion for notification
	participants, err := ics.repo.GetParticipantsByNomineeDepartmentID(nomineeDepartmentID)
	if err != nil {
		return err
	}

	var deletedParticipant *models.InternalCommitteeVersionNomineeDepartmentParticipant
	for _, participant := range participants {
		if participant.ID == participantID {
			deletedParticipant = participant
			break
		}
	}

	if deletedParticipant == nil {
		return utility.NewResourceNotFoundError(fmt.Errorf("participant not found"), participantID)
	}

	tx := ics.repo.BeginTransaction()

	if err := ics.repo.DeleteParticipant(nomineeDepartmentID, participantID, tx); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	// Send removal notification to the deleted participant
	go func() {
		// Get committee details through the relationship chain
		if deletedParticipant.NomineeDepartment == nil || deletedParticipant.NomineeDepartment.InternalCommitteeVersion == nil {
			utility.LogError("Missing nominee department or version relationship for removal notification", map[string]interface{}{
				"nominee_department_id": nomineeDepartmentID,
			})
			return
		}

		committeeID := deletedParticipant.NomineeDepartment.InternalCommitteeVersion.InternalCommitteeID

		// Get full committee details
		committee, err := ics.GetFullDetailsByID(committeeID)
		if err != nil {
			utility.LogError("Failed to get committee details for removal notification", map[string]interface{}{
				"committee_id": committeeID,
				"error":        err.Error(),
			})
			return
		}

		// Get department details
		department, err := ics.organizationUnitSvc.GetByID(deletedParticipant.NomineeDepartment.NomineeDepartmentID)
		if err != nil {
			utility.LogError("Failed to get department details for removal notification", map[string]interface{}{
				"department_id": deletedParticipant.NomineeDepartment.NomineeDepartmentID,
				"error":         err.Error(),
			})
			return
		}

		// Get user details
		user, err := ics.userSvc.GetByID(deletedParticipant.UserID)
		if err != nil {
			utility.LogError("Failed to get user details for removal notification", map[string]interface{}{
				"user_id": deletedParticipant.UserID,
				"error":   err.Error(),
			})
			return
		}

		// Get department head (current user) who is removing the participant
		departmentHead, err := ics.userSvc.GetByID(currentUserID)
		if err != nil {
			utility.LogError("Failed to get department head details for removal notification", map[string]interface{}{
				"user_id": currentUserID,
				"error":   err.Error(),
			})
			return
		}

		urlParam := fmt.Sprintf(`{"id":%d}`, committeeID)
		entity := "internal_committee"
		action := "participant-removed"

		dbNotificationDTO := &commonDTOs.DBNotificationDTO{
			UserIDs:              []uint{deletedParticipant.UserID},
			ModuleIdentifier:     "organization",
			NotificationType:     string(commonDTOs.InformationType),
			ReferenceNumber:      committee.ReferenceNumber,
			HeadingEn:            fmt.Sprintf("Committee Participant Removed: %s", committee.LatestVersion.CommitteeNameEn),
			HeadingAr:            fmt.Sprintf("تم إزالة عضو اللجنة: %s", committee.LatestVersion.CommitteeNameAr),
			TitleEn:              fmt.Sprintf("You have been removed as a participant from committee: %s", committee.LatestVersion.CommitteeNameEn),
			TitleAr:              fmt.Sprintf("لقد تم إزالتك كمشارك من اللجنة: %s", committee.LatestVersion.CommitteeNameAr),
			Entity:               &entity,
			Action:               &action,
			URLParam:             urlParam,
			IsSystemNotification: false,
			CreatedByID:          &currentUserID,
		}

		emailTo := []string{user.Email}
		subject := fmt.Sprintf("Committee/%s: Participant Removal", committee.LatestVersion.CommitteeNameEn)
		appURL := utility.GetAppURL()
		emailData := map[string]interface{}{
			"FirstNameEn":        user.FirstNameEn,
			"LastNameEn":         user.LastNameEn,
			"CommitteeReference": committee.ReferenceNumber,
			"CommitteeNameEn":    committee.LatestVersion.CommitteeNameEn,
			"CommitteeNameAr":    committee.LatestVersion.CommitteeNameAr,
			"DepartmentNameEn":   department.TitleEn,
			"DepartmentNameAr":   department.TitleAr,
			"MembershipType":     ics.getCommitteeMemberTypeTitleFromNomineeDepartment(deletedParticipant.NomineeDepartment),
			"RemovedBy":          fmt.Sprintf("%s %s", departmentHead.FirstNameEn, departmentHead.LastNameEn),
			"RemovedOn":          time.Now().Format("2006-01-02"),
			"RedirectURL":        fmt.Sprintf("%s/organization/committee/%d/overview", appURL, committeeID),
		}

		emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
			EmailTo:              emailTo,
			Subject:              subject,
			TemplateRelativePath: "organization/templates/email/internal_committee_participant_removed_mail.html",
			EmailData:            emailData,
		}

		err = ics.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "committee_management")
		if err != nil {
			utility.LogError("Failed to send participant removal notification", map[string]interface{}{
				"error":   err.Error(),
				"user_id": deletedParticipant.UserID,
			})
		}
	}()

	// Log activity asynchronously
	go func() {
		err := ics.activityLogSvc.LogActivity(
			models.DeleteAction,
			"committee_participants",
			nomineeDepartmentID,
			map[string]interface{}{
				"nominee_department_id": nomineeDepartmentID,
			},
			currentUserID,
			"organization",
			"committee_management",
			"delete committee participants by department",
			nil,
		)
		if err != nil {
			utility.LogError("Failed to log activity", map[string]interface{}{
				"error": err.Error(),
			})
		}
	}()

	return nil
}

func (ics *InternalCommitteeService) GetLatestVersionByID(id uint) (*models.InternalCommittee, error) {
	internalCommittee, err := ics.repo.GetLatestVersionByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}
	return internalCommittee, nil
}

func (ics *InternalCommitteeService) DisbandCommittee(committeeID uint, currentUserID uint) error {
	tx := ics.repo.BeginTransaction()

	// Get existing committee
	existingCommittee, err := ics.repo.GetByID(committeeID)
	if err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utility.NewResourceNotFoundError(err, committeeID)
		}
		return err
	}

	// Update committee
	existingCommittee.UpdatedBYID = &currentUserID
	if err := ics.repo.Update(existingCommittee, tx); err != nil {
		tx.Rollback()
		return err
	}

	// Get existing versions
	versions, err := ics.repo.GetVersionsByCommitteeID(committeeID, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Find latest version and get its details
	var latestVersion *models.InternalCommitteeVersion
	var oldNomineeDepartments []*models.InternalCommitteeVersionNomineeDepartment
	for _, v := range versions {
		if v.IsLatestVersion != nil && *v.IsLatestVersion {
			latestVersion = v
			oldNomineeDepartments, err = ics.repo.GetNomineeDepartmentsByVersionID(v.ID, tx)
			if err != nil {
				tx.Rollback()
				return err
			}
			break
		}
	}

	if latestVersion == nil {
		tx.Rollback()
		return fmt.Errorf("no latest version found for committee")
	}

	// Set old version's is_latest to false
	isLatestFalse := false
	latestVersion.IsLatestVersion = &isLatestFalse
	latestVersion.UpdatedBYID = &currentUserID
	if err := ics.repo.UpdateVersion(latestVersion, tx); err != nil {
		tx.Rollback()
		return err
	}

	// Get draft status
	draftStatus, err := ics.internalCommitteeStatusSvc.GetByType("draft")
	if err != nil {
		tx.Rollback()
		return err
	}

	// Calculate new version number
	newVersionNumber := *latestVersion.Version + 1
	isLatestTrue := true

	// Create new version with disband_required status
	newVersion := &models.InternalCommitteeVersion{
		InternalCommitteeID:       committeeID,
		CommitteeNameEn:           latestVersion.CommitteeNameEn,
		CommitteeNameAr:           latestVersion.CommitteeNameAr,
		CommitteeTypeID:           latestVersion.CommitteeTypeID,
		CommitteeTermID:           latestVersion.CommitteeTermID,
		TermEndDate:               latestVersion.TermEndDate,
		CommitteeClassificationID: latestVersion.CommitteeClassificationID,
		Justification:             latestVersion.Justification,
		RolesAndResponsibilities:  latestVersion.RolesAndResponsibilities,
		CommitteeAuthorities:      latestVersion.CommitteeAuthorities,
		CommitteeStatusID:         draftStatus.ID,
		IsDisbandRequested:        true,
		Version:                   &newVersionNumber,
		IsLatestVersion:           &isLatestTrue,
		CreatedBYID:               currentUserID,
	}

	newVersionID, err := ics.repo.CreateVersion(newVersion, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Copy nominee departments and participants from old version
	for _, oldDept := range oldNomineeDepartments {
		newDept := &models.InternalCommitteeVersionNomineeDepartment{
			InternalCommitteeVersionID: newVersionID,
			NomineeDepartmentID:        oldDept.NomineeDepartmentID,
			CommitteeMemberTypeID:      oldDept.CommitteeMemberTypeID,
			CreatedBYID:                &currentUserID,
		}

		// Create nominee department
		var nomineeDepartments []*models.InternalCommitteeVersionNomineeDepartment
		nomineeDepartments = append(nomineeDepartments, newDept)
		if err := ics.repo.CreateNomineeDepartments(nomineeDepartments, tx); err != nil {
			tx.Rollback()
			return err
		}

		// Copy participants from old department
		oldParticipants, err := ics.repo.GetParticipantsByNomineeDepartmentID(oldDept.ID, tx)
		if err != nil {
			tx.Rollback()
			return err
		}

		for _, oldParticipant := range oldParticipants {
			newParticipant := &models.InternalCommitteeVersionNomineeDepartmentParticipant{
				InternalCommitteeVersionNomineeDepartmentID: newDept.ID,
				UserID:       oldParticipant.UserID,
				DepartmentID: oldParticipant.DepartmentID,
				CreatedBYID:  &currentUserID,
			}

			if _, err := ics.repo.CreateParticipant(newParticipant, tx); err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	tx.Commit()

	// Automatically submit for workflow approval
	submitErr := ics.submitDisbandForWorkflow(committeeID, newVersionID, currentUserID)
	if submitErr != nil {
		return submitErr
	}

	// Send notifications to nominee departments asynchronously
	go func() {

		committeeDetails, err := ics.GetFullDetailsByID(committeeID)
		if err != nil {
			utility.LogError("Failed to get committee details for activity logging", map[string]interface{}{
				"error":       err.Error(),
				"committeeID": committeeID,
			})
			return
		}

		err = ics.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.InternalCommitteeVersion{}.TableName(),
			latestVersion.ID,
			committeeDetails,
			currentUserID,
			"organization",
			"committee_management",
			"disband committee",
			nil,
		)
		if err != nil {
			utility.LogError("Failed to log disband activity", map[string]interface{}{
				"error":       err.Error(),
				"committeeID": committeeID,
			})
		}
		// Get nominee departments for the latest version
		nomineeDepartments, err := ics.repo.GetNomineeDepartmentsByVersionID(latestVersion.ID)
		if err != nil {
			utility.LogError("Failed to fetch nominee departments for disband notifications", map[string]interface{}{
				"error":       err.Error(),
				"committeeID": committeeID,
			})
			return
		}

		// Get user making the disband request
		updatedByUser, err := ics.userSvc.GetByID(currentUserID)
		if err != nil {
			utility.LogError("Failed to fetch user making disband request", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		appURL := utility.GetAppURL()
		urlParam := fmt.Sprintf(`{"id":%d}`, committeeID)
		entity := "internal_committee"

		for _, nomineeDept := range nomineeDepartments {
			// Get department head users (following existing pattern)
			headUsers, err := ics.organizationUnitsSvc.GetParentOrganizationUnitWithHead(nomineeDept.NomineeDepartmentID)
			if err != nil {
				utility.LogError("Failed to fetch department head users for disband notification", map[string]interface{}{
					"error":         err.Error(),
					"nomineeDeptID": nomineeDept.NomineeDepartmentID,
					"committeeID":   committeeID,
				})
				continue
			}

			if headUsers.HeadUser == nil {
				continue
			}

			// Get department details
			department, err := ics.organizationUnitSvc.GetByID(nomineeDept.NomineeDepartmentID)
			if err != nil {
				utility.LogError("Failed to fetch department details for disband notification", map[string]interface{}{
					"error":         err.Error(),
					"nomineeDeptID": nomineeDept.NomineeDepartmentID,
					"committeeID":   committeeID,
				})
				continue
			}

			// Collect user IDs and emails for notifications
			var userIDs []uint
			var userEmails []string

			if headUsers.HeadUser != nil {
				userIDs = append(userIDs, headUsers.HeadUser.ID)
				userEmails = append(userEmails, headUsers.HeadUser.Email)
			}

			// DB Notification
			action := "committee-disband-required"
			dbNotificationDTO := &commonDTOs.DBNotificationDTO{
				UserIDs:              userIDs,
				ModuleIdentifier:     "organization",
				NotificationType:     string(commonDTOs.InformationType),
				HeadingEn:            "Committee Disband Required",
				HeadingAr:            "مطلوب حل اللجنة",
				TitleEn:              committeeDetails.LatestVersion.CommitteeNameEn,
				TitleAr:              committeeDetails.LatestVersion.CommitteeNameAr,
				ReferenceNumber:      committeeDetails.ReferenceNumber,
				Entity:               &entity,
				Action:               &action,
				URLParam:             urlParam,
				IsSystemNotification: false,
				CreatedByID:          &currentUserID,
			}

			// Email Notification
			subject := fmt.Sprintf("Committee/%s: Disband Required", committeeDetails.LatestVersion.CommitteeNameEn)
			emailData := map[string]interface{}{
				"DepartmentNameEn": department.TitleEn,
				"DepartmentNameAr": department.TitleAr,
				"CommitteeNameEn":  committeeDetails.LatestVersion.CommitteeNameEn,
				"CommitteeNameAr":  committeeDetails.LatestVersion.CommitteeNameAr,
				"ReferenceNumber":  committeeDetails.ReferenceNumber,
				"FirstNameEn":      updatedByUser.FirstNameEn,
				"LastNameEn":       updatedByUser.LastNameEn,
				"DisbandDate":      time.Now().Format("2006-01-02"),
				"RedirectURL":      fmt.Sprintf("%s/organization/internal-committee/%d", appURL, committeeID),
			}

			emailNotificationDTO := &commonDTOs.EmailNotificationDTO{
				EmailTo:              userEmails,
				Subject:              subject,
				TemplateRelativePath: "organization/templates/email/internal_committee_disbanded_mail.html",
				EmailData:            emailData,
			}

			// Send notifications
			err = ics.notificationSvc.SendNotification(dbNotificationDTO, emailNotificationDTO, "committee_management")
			if err != nil {
				utility.LogError("Failed to send disband notifications", map[string]interface{}{
					"error":         err.Error(),
					"nomineeDeptID": nomineeDept.NomineeDepartmentID,
					"committeeID":   committeeID,
				})
			}
		}
	}()

	return nil
}

func (ics *InternalCommitteeService) submitDisbandForWorkflow(committeeID uint, versionID uint, currentUserID uint) error {
	// Simply call the Submit method from the workflow review service
	return ics.workflowReviewSvc.Submit(committeeID, currentUserID)
}

func (ics *InternalCommitteeService) ValidateAllDepartmentsHaveParticipants(versionID uint, tx ...utility.DatabaseTransaction) (bool, error) {
	return ics.repo.ValidateAllDepartmentsHaveParticipants(versionID, tx...)
}
