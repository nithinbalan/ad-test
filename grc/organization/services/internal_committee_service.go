package services

import (
	"errors"
	"fmt"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	organizationDTOs "jethur/grc/organization/dtos"
	organizationRepositories "jethur/grc/organization/repositories"
	settingsDTOs "jethur/grc/settings/dtos"
	settingsServices "jethur/grc/settings/services"
	"jethur/utility"

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

	// Participant methods
	CreateParticipants(currentUserID uint, nomineeDepartmentID uint, dto *organizationDTOs.InternalCommitteeParticipantFormDTO) ([]uint, error)
	DeleteParticipant(nomineeDepartmentID uint, participantID uint, currentUserID uint) error
}

type InternalCommitteeService struct {
	repo                       organizationRepositories.InternalCommitteeRepositoryInterface
	activityLogSvc             logServices.ActivityLogServiceInterface
	userSvc                    UserServiceInterface
	organizationUnitSvc        settingsServices.OrganizationUnitServiceInterface
	committeeTypeSvc           CommitteeTypeServiceInterface
	committeeTermSvc           CommitteeTermServiceInterface
	committeeClassificationSvc CommitteeClassificationServiceInterface
	internalCommitteeStatusSvc InternalCommitteeStatusServiceInterface
}

func NewInternalCommitteeService(
	repo organizationRepositories.InternalCommitteeRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	userSvc UserServiceInterface,
	organizationUnitSvc settingsServices.OrganizationUnitServiceInterface,
	committeeTypeSvc CommitteeTypeServiceInterface,
	committeeTermSvc CommitteeTermServiceInterface,
	committeeClassificationSvc CommitteeClassificationServiceInterface,
	internalCommitteeStatusSvc InternalCommitteeStatusServiceInterface,
) InternalCommitteeServiceInterface {
	return &InternalCommitteeService{
		repo:                       repo,
		activityLogSvc:             activityLogSvc,
		userSvc:                    userSvc,
		organizationUnitSvc:        organizationUnitSvc,
		committeeTypeSvc:           committeeTypeSvc,
		committeeTermSvc:           committeeTermSvc,
		committeeClassificationSvc: committeeClassificationSvc,
		internalCommitteeStatusSvc: internalCommitteeStatusSvc,
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

		// Manually map nominee departments with participants
		if len(latestVersion.InternalCommitteeVersionNomineeDepartments) > 0 {
			var nomineeDepartmentsDTO []*organizationDTOs.InternalCommitteeVersionNomineeDepartmentDetailDTO
			for _, nomineeDept := range latestVersion.InternalCommitteeVersionNomineeDepartments {
				nomineeDeptDTO := &organizationDTOs.InternalCommitteeVersionNomineeDepartmentDetailDTO{
					ID:                         nomineeDept.ID,
					InternalCommitteeVersionID: nomineeDept.InternalCommitteeVersionID,
					NomineeDepartmentID:        nomineeDept.NomineeDepartmentID,
					MembershipType:             nomineeDept.MembershipType,
				}

				// Map nominee department
				if nomineeDept.NomineeDepartment != nil {
					nomineeDeptDTO.NomineeDepartment = &settingsDTOs.OrganizationUnitDTO{
						ID:      int(nomineeDept.NomineeDepartment.ID),
						TitleEn: nomineeDept.NomineeDepartment.TitleEn,
						TitleAr: nomineeDept.NomineeDepartment.TitleAr,
						Status:  nomineeDept.NomineeDepartment.Status,
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
		nomineeDepartments = append(nomineeDepartments, &models.InternalCommitteeVersionNomineeDepartment{
			InternalCommitteeVersionID: versionID,
			NomineeDepartmentID:        nominee.NomineeDepartmentID,
			MembershipType:             nominee.MembershipType,
			CreatedBYID:                &currentUserID,
		})
	}

	if err := ics.repo.CreateNomineeDepartments(nomineeDepartments, tx); err != nil {
		tx.Rollback()
		return 0, err
	}

	tx.Commit()

	// Log activity asynchronously
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
		nomineeDepartments = append(nomineeDepartments, &models.InternalCommitteeVersionNomineeDepartment{
			InternalCommitteeVersionID: latestVersion.ID,
			NomineeDepartmentID:        nominee.NomineeDepartmentID,
			MembershipType:             nominee.MembershipType,
			CreatedBYID:                &currentUserID,
		})
	}

	if err := ics.repo.CreateNomineeDepartments(nomineeDepartments, tx); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	// Log activity asynchronously with changed fields only
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
		newDept := &models.InternalCommitteeVersionNomineeDepartment{
			InternalCommitteeVersionID: newVersionID,
			NomineeDepartmentID:        newDeptDTO.NomineeDepartmentID,
			MembershipType:             newDeptDTO.MembershipType,
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
					UserID:      oldParticipant.UserID,
					CreatedBYID: &currentUserID,
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
	}

	return nil
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

	tx := ics.repo.BeginTransaction()

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
			UserID:      userID,
			CreatedBYID: &currentUserID,
		}

		id, err := ics.repo.CreateParticipant(participant, tx)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		createdIDs = append(createdIDs, id)
	}

	tx.Commit()

	// Log activity asynchronously
	go func() {
		for _, id := range createdIDs {
			for _, userID := range dto.UserIDs {
				err := ics.activityLogSvc.LogActivity(
					models.CreateAction,
					models.InternalCommitteeVersionNomineeDepartmentParticipant{}.TableName(),
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
	tx := ics.repo.BeginTransaction()

	if err := ics.repo.DeleteParticipant(nomineeDepartmentID, participantID, tx); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	// Log activity asynchronously
	go func() {
		err := ics.activityLogSvc.LogActivity(
			models.DeleteAction,
			models.InternalCommitteeVersionNomineeDepartmentParticipant{}.TableName(),
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
