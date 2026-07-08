package services

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	businessContinuityDTOs "jethur/grc/business-continuity/dtos"
	TestAndExerciseRepositories "jethur/grc/business-continuity/repositories"
	commonDTOs "jethur/grc/common/dtos"
	logServices "jethur/grc/logs/services"
	masterServices "jethur/grc/masters/services"
	"jethur/grc/models"
	organizationServices "jethur/grc/organization/services"
	taskServices "jethur/grc/task/services"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type TestAndExerciseServiceInterface interface {
	GetAll(params *businessContinuityDTOs.TestAndExerciseListQueryParamsDTO) (*businessContinuityDTOs.TestAndExerciseListResponseDTO, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*businessContinuityDTOs.TestAndExerciseDetailDTO, error)
	GetByID(id uint) (*models.TestAndExercise, error)
	Create(currentUserID uint, businessContinuityDTOs *businessContinuityDTOs.TestAndExerciseFormDTO) (uint, error)
	Update(id uint, currentUserID uint, TestAndExercise *businessContinuityDTOs.TestAndExerciseFormDTO) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)

	GetTestAndExerciseProcessByID(id uint) (*models.TestAndExerciseProcess, error)
	GetTestAndExerciseSystemRegisterByID(id uint) (*models.TestAndExerciseSystemRegister, error)
	ExportTestAndExercise(data []*businessContinuityDTOs.TestAndExerciseListDTO, currentUserID uint) (*os.File, error)

	GetTestResultSummary(testAndExerciseID uint) (*businessContinuityDTOs.TestResultSummaryDTO, error)
}

// Service and implementations
type TestAndExerciseService struct {
	repo                             TestAndExerciseRepositories.TestAndExerciseRepositoryInterface
	activityLogSvc                   logServices.ActivityLogServiceInterface
	userSvc                          organizationServices.UserServiceInterface
	bcpTestAndExerciseSvc            BcpTestAndExerciseServiceInterface
	businessContinuityPlanSvc        BusinessContinuityPlanServiceInterface
	testAndExerciseStatusSvc         TestAndExerciseStatusServiceInterface
	communicationActionPlanStatusSvc CommunicationActionPlanStatusServiceInterface
	actionPlanSvc                    masterServices.ActionPlanStatusServiceInterface
	userSettingSvc                   organizationServices.UserSettingServiceInterface
	taskSvc                          taskServices.TaskServiceInterface
}

func NewTestAndExerciseService(
	repo TestAndExerciseRepositories.TestAndExerciseRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
	userSvc organizationServices.UserServiceInterface,
	bcpTestAndExerciseSvc BcpTestAndExerciseServiceInterface,
	businessContinuityPlanSvc BusinessContinuityPlanServiceInterface,
	testAndExerciseStatusSvc TestAndExerciseStatusServiceInterface,
	communicationActionPlanStatusSvc CommunicationActionPlanStatusServiceInterface,
	actionPlanSvc masterServices.ActionPlanStatusServiceInterface,
	userSettingSvc organizationServices.UserSettingServiceInterface,
	taskSvc taskServices.TaskServiceInterface,
) TestAndExerciseServiceInterface {
	return &TestAndExerciseService{
		repo:                             repo,
		activityLogSvc:                   activityLogSvc,
		userSvc:                          userSvc,
		bcpTestAndExerciseSvc:            bcpTestAndExerciseSvc,
		businessContinuityPlanSvc:        businessContinuityPlanSvc,
		testAndExerciseStatusSvc:         testAndExerciseStatusSvc,
		communicationActionPlanStatusSvc: communicationActionPlanStatusSvc,
		actionPlanSvc:                    actionPlanSvc,
		userSettingSvc:                   userSettingSvc,
		taskSvc:                          taskSvc,
	}
}

func (tes *TestAndExerciseService) GetAll(params *businessContinuityDTOs.TestAndExerciseListQueryParamsDTO) (*businessContinuityDTOs.TestAndExerciseListResponseDTO, error) {
	return tes.repo.GetAll(params)
}

func (tes TestAndExerciseService) GetByID(id uint) (*models.TestAndExercise, error) {
	TestAndExercise, err := tes.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	return TestAndExercise, nil
}

func (tes TestAndExerciseService) GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*businessContinuityDTOs.TestAndExerciseDetailDTO, error) {
	// Initialize your DTO
	var testAndExerciseDTO businessContinuityDTOs.TestAndExerciseDetailDTO

	testAndExercise, err := tes.repo.GetFullDetailsByID(id, tx...)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	// Map the result to your DTO
	err = mapstructure.Decode(testAndExercise, &testAndExerciseDTO)
	if err != nil {
		return nil, err
	}

	if testAndExercise.TestingDate != nil {
		testingDate := testAndExercise.TestingDate.Format("2006-01-02")
		testAndExerciseDTO.TestingDateString = testingDate
	}

	testAndExerciseDTO.CreatedAt = testAndExercise.CreatedAt
	testAndExerciseDTO.UpdatedAt = testAndExercise.UpdatedAt

	return &testAndExerciseDTO, nil
}

func (tes *TestAndExerciseService) Create(currentUserID uint, testAndExerciseFormDTO *businessContinuityDTOs.TestAndExerciseFormDTO) (uint, error) {

	responsibleUnitID, err := tes.userSvc.GetUserPrimaryUnitID(testAndExerciseFormDTO.ResponsibleUserID)
	if err != nil {
		return 0, utility.NewValidationError(fmt.Errorf("organization unit not set for responsible user"), "not_exist", "test and exercise responsible unit")
	}

	testAndExerciseStatus, err := tes.testAndExerciseStatusSvc.GetByType("scheduled")
	if err != nil {
		return 0, err
	}

	actionPlanStatus, err := tes.actionPlanSvc.GetByType("new")
	if err != nil {
		return 0, err
	}

	communicationActionPlanStatus, err := tes.communicationActionPlanStatusSvc.GetByType("pending")
	if err != nil {
		return 0, err
	}

	tx := tes.repo.BeginTransaction()

	referenceNumber, err := tes.repo.GenerateReferenceNumber()
	if err != nil {
		return 0, err
	}

	dates, err := tes.parseTestAndExerciseDates(testAndExerciseFormDTO)
	if err != nil {
		return 0, err
	}

	TestAndExercise := &models.TestAndExercise{
		ReferenceNumber:          referenceNumber,
		BusinessContinuityPlanID: testAndExerciseFormDTO.BusinessContinuityPlanID,
		BcpTestAndExerciseID:     testAndExerciseFormDTO.BcpTestAndExerciseID,
		ResponsibleUserID:        &testAndExerciseFormDTO.ResponsibleUserID,
		ResponsibleUnitID:        &responsibleUnitID,
		TestingDate:              dates.TestingDate,
		Description:              testAndExerciseFormDTO.Description,
		TestAndExerciseStatusID:  testAndExerciseStatus.ID,
		CreatedBYID:              currentUserID,
	}

	id, err := tes.repo.Create(TestAndExercise, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	// Save Processes

	bcpTestAndExercise, err := tes.bcpTestAndExerciseSvc.GetFullDetailsByID(testAndExerciseFormDTO.BusinessContinuityPlanID, testAndExerciseFormDTO.BcpTestAndExerciseID, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	var testAndExerciseProcesses []*models.TestAndExerciseProcess
	processIds := make(map[uint]bool)
	for _, bcpProcess := range bcpTestAndExercise.BcpTestAndExerciseBcpProcesses {
		process := &models.TestAndExerciseProcess{
			TestAndExerciseID: id,
			ProcessID:         bcpProcess.BcpProcess.ProcessID,
			PlannedRto:        bcpProcess.BcpProcess.RtoFromValue,
			PlannedRtoUnitId:  bcpProcess.BcpProcess.RtoUnitId,
			PlannedMtpd:       bcpProcess.BcpProcess.MtpdFromValue,
			PlannedMtpdUnitId: bcpProcess.BcpProcess.MtpdUnitId,
			PlannedRpo:        bcpProcess.BcpProcess.RpoFromValue,
			PlannedRpoUnitId:  bcpProcess.BcpProcess.RpoUnitId,
			CreatedBYID:       currentUserID,
		}
		testAndExerciseProcesses = append(testAndExerciseProcesses, process)
		processIds[bcpProcess.BcpProcess.ProcessID] = true
	}
	if len(testAndExerciseProcesses) > 0 {
		if err := tes.repo.CreateTestAndExerciseProcesses(testAndExerciseProcesses, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	// Save System Registers
	var testAndExerciseSystemRegisters []*models.TestAndExerciseSystemRegister
	systemRegisterIds := make(map[uint]bool)
	for _, bcpSystemRegister := range bcpTestAndExercise.BcpTestAndExerciseSystemRegisters {
		systemRegister := &models.TestAndExerciseSystemRegister{
			TestAndExerciseID: id,
			SystemRegisterID:  bcpSystemRegister.SystemRegisterID,
			PlannedRto:        bcpSystemRegister.SystemRegister.RTO,
			PlannedRtoUnitId:  bcpSystemRegister.SystemRegister.RTOUnitID,
			PlannedMtpd:       bcpSystemRegister.SystemRegister.MTPD,
			PlannedMtpdUnitId: bcpSystemRegister.SystemRegister.MTPDUnitID,
			PlannedRpo:        bcpSystemRegister.SystemRegister.RPO,
			PlannedRpoUnitId:  bcpSystemRegister.SystemRegister.RPOUnitID,
			CreatedBYID:       currentUserID,
		}
		testAndExerciseSystemRegisters = append(testAndExerciseSystemRegisters, systemRegister)
		systemRegisterIds[bcpSystemRegister.SystemRegisterID] = true
	}
	if len(testAndExerciseSystemRegisters) > 0 {
		if err := tes.repo.CreateTestAndExerciseSystemRegisters(testAndExerciseSystemRegisters, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	// Auto-create TestAndExerciseTestResult with one child row per system register
	testResult := &models.TestAndExerciseTestResult{
		TestAndExerciseID: id,
		CreatedBYID:       currentUserID,
	}
	testResultID, err := tes.repo.CreateTestAndExerciseTestResult(testResult, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	if len(bcpTestAndExercise.BcpTestAndExerciseSystemRegisters) > 0 {
		var testResultSystemRegisters []*models.TestAndExerciseTestResultSystemRegister
		for _, bcpSR := range bcpTestAndExercise.BcpTestAndExerciseSystemRegisters {
			item := &models.TestAndExerciseTestResultSystemRegister{
				TestAndExerciseTestResultID: testResultID,
				SystemRegisterID:            bcpSR.SystemRegisterID,
				CreatedBYID:                 currentUserID,
			}
			testResultSystemRegisters = append(testResultSystemRegisters, item)
		}
		if err := tes.repo.CreateTestAndExerciseTestResultSystemRegisters(testResultSystemRegisters, tx); err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	//  save recovery action plan

	businessContinuityPlan, err := tes.businessContinuityPlanSvc.GetFullDetailsByID(testAndExerciseFormDTO.BusinessContinuityPlanID, tx)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	if businessContinuityPlan.BcpRecoveryPlans != nil {
		for _, bcpRecoveryPlan := range businessContinuityPlan.BcpRecoveryPlans {
			if _, exists := processIds[bcpRecoveryPlan.BcpProcess.Process.ID]; !exists {
				continue
			}
			testingDate, err := utility.ParseDate(testAndExerciseFormDTO.TestingDate) // time.Parse("2006-01-02", testAndExerciseFormDTO.TestingDate)
			if err != nil {
				// Handle parsing error
				continue
			}
			var targetDate time.Time
			if bcpRecoveryPlan.TargetCompletion != nil {
				// Convert TargetCompletion to an integer
				targetCompletion, err := strconv.Atoi(*bcpRecoveryPlan.TargetCompletion)
				if err != nil {
					// Handle parsing error for TargetCompletion
					continue
				}

				switch bcpRecoveryPlan.TargetCompletionUnit.TitleEn {
				case "Minutes":
					targetDate = testingDate.Add(time.Duration(targetCompletion) * time.Minute)
				case "Hours":
					targetDate = testingDate.Add(time.Duration(targetCompletion) * time.Hour)
				case "Days":
					targetDate = testingDate.AddDate(0, 0, targetCompletion)
				case "Weeks":
					targetDate = testingDate.AddDate(0, 0, targetCompletion*7)
				default:
					// Handle unsupported unit
					continue
				}
			} else {
				targetDate = *testingDate
			}

			process := &models.TestAndExerciseRecoveryActionPlan{
				TestAndExerciseID:  id,
				ProcessID:          &bcpRecoveryPlan.BcpProcess.Process.ID,
				Title:              bcpRecoveryPlan.Title,
				SeverityLevelID:    &bcpRecoveryPlan.SeverityLevel.ID,
				BcpPhaseID:         &bcpRecoveryPlan.BcpPhaseID,
				ResponsibleUserID:  bcpRecoveryPlan.ResponsibleUserID,
				ResponsibleUnitID:  bcpRecoveryPlan.ResponsibleUnitID,
				TargetDate:         &targetDate,
				ActionPlanStatusId: &actionPlanStatus.ID,
				IsRecovery:         true,
				CreatedBYID:        currentUserID,
			}

			actionPlanId, err := tes.repo.CreateTestAndExerciseRecoveryActionPlans(process, tx)
			if err != nil {
				tx.Rollback()
				return 0, err
			}
			urlParam := fmt.Sprintf(`{"test_and_exercise_id":%d, "recovery_action_plan_id":%d}`, id, actionPlanId)
			taskDTO := commonDTOs.TaskFormDTO{
				ModuleIdentifier:     "test_and_exercise",
				RecoveryActionPlanID: &actionPlanId,
				Title:                bcpRecoveryPlan.Title,
				TargetDate:           &targetDate,
				URLParam:             &urlParam,
				QueryParam:           nil,
				SeverityLevelID:      &bcpRecoveryPlan.SeverityLevel.ID,
				ResponsibleUserID:    bcpRecoveryPlan.ResponsibleUserID,
			}

			if _, err := tes.taskSvc.Create(taskDTO, currentUserID, tx); err != nil {
				tx.Rollback()
				return 0, err
			}
		}
	}

	//  save communication action plan
	if businessContinuityPlan.BcpCommunicationPlans != nil {
		var testAndExerciseCommunicationActionPlan []*models.TestAndExerciseCommunicationActionPlan

		for _, bcpCommunicationPlan := range businessContinuityPlan.BcpCommunicationPlans {
			if _, exists := processIds[bcpCommunicationPlan.BcpProcess.ProcessID]; !exists {
				continue
			}

			process := &models.TestAndExerciseCommunicationActionPlan{
				TestAndExerciseID:               id,
				ProcessID:                       &bcpCommunicationPlan.BcpProcess.ProcessID,
				Title:                           bcpCommunicationPlan.Title,
				BcpPhaseID:                      bcpCommunicationPlan.BcpPhaseID,
				ResponsibleUserID:               bcpCommunicationPlan.ResponsibleUserID,
				ResponsibleUnitID:               bcpCommunicationPlan.ResponsibleUnitID,
				Recipient:                       bcpCommunicationPlan.Recipient,
				Frequency:                       bcpCommunicationPlan.Frequency,
				Channel:                         bcpCommunicationPlan.Channel,
				CommunicationActionPlanStatusID: &communicationActionPlanStatus.ID,
				CreatedBYID:                     currentUserID,
			}
			testAndExerciseCommunicationActionPlan = append(testAndExerciseCommunicationActionPlan, process)
		}

		if len(testAndExerciseCommunicationActionPlan) > 0 {
			if err := tes.repo.CreateTestAndExerciseCommunicationActionPlans(testAndExerciseCommunicationActionPlan, tx); err != nil {
				tx.Rollback()
				return 0, err
			}
		}
	}

	tx.Commit()
	go func() {
		existingTestAndExercise, err := tes.GetFullDetailsByID(TestAndExercise.ID)
		if err != nil {
			utility.LogError("Failed to fetch test and exercise", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}

		err = tes.activityLogSvc.LogActivity(
			models.CreateAction,
			models.TestAndExercise{}.TableName(),
			id,
			existingTestAndExercise,
			currentUserID,
			"business_continuity",
			"test_and_exercise",
			"test_and_exercise",
			nil,
		)
		if err != nil {
			fields := map[string]interface{}{
				"action":    models.CreateAction,
				"tableName": models.TestAndExercise{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			}

			utility.LogError("Failed to log activity", fields)
		}
	}()
	return id, nil
}

func (tes *TestAndExerciseService) Update(id uint, currentUserID uint, testAndExerciseFormDTO *businessContinuityDTOs.TestAndExerciseFormDTO) error {
	tx := tes.repo.BeginTransaction()

	existingTestAndExercise, err := tes.repo.GetByID(id)
	if err != nil {
		return utility.NewResourceNotFoundError(err, id)
	}

	oldTestAndExercise, err := tes.GetFullDetailsByID(id, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	if existingTestAndExercise.ResponsibleUserID != &testAndExerciseFormDTO.ResponsibleUserID {
		responsibleUnitID, err := tes.userSvc.GetUserPrimaryUnitID(*existingTestAndExercise.ResponsibleUserID)
		if err != nil {
			return utility.NewValidationError(fmt.Errorf("organization unit not set for responsible user"), "not_exist", "risk responsible unit")
		}
		existingTestAndExercise.ResponsibleUnitID = &responsibleUnitID
	}

	dates, err := tes.parseTestAndExerciseDates(testAndExerciseFormDTO)
	if err != nil {
		return err
	}

	existingTestAndExercise.BusinessContinuityPlanID = testAndExerciseFormDTO.BusinessContinuityPlanID
	existingTestAndExercise.ResponsibleUserID = &testAndExerciseFormDTO.ResponsibleUserID
	existingTestAndExercise.Description = testAndExerciseFormDTO.Description
	existingTestAndExercise.TestingDate = dates.TestingDate
	existingTestAndExercise.BcpTestAndExerciseID = testAndExerciseFormDTO.BcpTestAndExerciseID
	existingTestAndExercise.UpdatedBYID = &currentUserID

	err = tes.repo.Update(existingTestAndExercise, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	if existingTestAndExercise.BusinessContinuityPlanID != oldTestAndExercise.BusinessContinuityPlanID {
		// Collect existing Process IDs
		existingTestAndExercise, err := tes.GetFullDetailsByID(id, tx)
		if err != nil {
			tx.Rollback()
			return err
		}

		existingBcpProcessIDs := make(map[uint]bool)
		for _, proc := range existingTestAndExercise.TestAndExerciseProcesses {
			existingBcpProcessIDs[proc.ProcessID] = true
		}

		bcpTestAndExercise, err := tes.bcpTestAndExerciseSvc.GetFullDetailsByID(testAndExerciseFormDTO.BusinessContinuityPlanID, testAndExerciseFormDTO.BcpTestAndExerciseID)
		if err != nil {
			tx.Rollback()
			return err
		}

		// Collect new Process IDs
		newBcpProcessIDs := make(map[uint]bool)
		for _, proc := range bcpTestAndExercise.BcpTestAndExerciseBcpProcesses {
			newBcpProcessIDs[proc.BcpProcess.ProcessID] = true
		}

		// Add new Processes
		var processesToAdd []*models.TestAndExerciseProcess
		for _, newProcess := range bcpTestAndExercise.BcpTestAndExerciseBcpProcesses {
			if _, exists := existingBcpProcessIDs[newProcess.BcpProcess.ProcessID]; !exists {
				// Create a new TestAndExercise Process record
				bcpProcess := &models.TestAndExerciseProcess{
					TestAndExerciseID: id,
					ProcessID:         *newProcess.BcpProcess.ID,
					PlannedRto:        newProcess.BcpProcess.RtoFromValue,
					PlannedRtoUnitId:  newProcess.BcpProcess.RtoUnitId,
					PlannedMtpd:       newProcess.BcpProcess.MtpdFromValue,
					PlannedMtpdUnitId: newProcess.BcpProcess.MtpdUnitId,
					PlannedRpo:        newProcess.BcpProcess.RpoFromValue,
					PlannedRpoUnitId:  newProcess.BcpProcess.RpoUnitId,
					CreatedBYID:       currentUserID,
				}
				processesToAdd = append(processesToAdd, bcpProcess)
			}
		}

		// Save the new Processes
		if len(processesToAdd) > 0 {
			if err := tes.repo.CreateTestAndExerciseProcesses(processesToAdd, tx); err != nil {
				tx.Rollback()
				return err
			}
		}

		// Remove Processes not in the new list
		var deleteProcessIDs []*uint

		for _, existingProcess := range existingTestAndExercise.TestAndExerciseProcesses {
			if _, exists := newBcpProcessIDs[existingProcess.ProcessID]; !exists {
				// Append bcpProcess to delete list
				deleteProcessIDs = append(deleteProcessIDs, existingProcess.ID)
			}
		}

		// Delete bcpProcesses not in the new list
		if len(deleteProcessIDs) > 0 {
			if err := tes.repo.DeleteTestAndExerciseProcesses(deleteProcessIDs, tx); err != nil {
				tx.Rollback()
				return err
			}
		}

		// Handle System Registers
		existingSystemRegisterIDs := make(map[uint]bool)
		for _, sr := range existingTestAndExercise.TestAndExerciseSystemRegisters {
			existingSystemRegisterIDs[sr.SystemRegisterID] = true
		}

		newSystemRegisterIDs := make(map[uint]bool)
		for _, sr := range bcpTestAndExercise.BcpTestAndExerciseSystemRegisters {
			newSystemRegisterIDs[sr.SystemRegisterID] = true
		}

		// Add new System Registers
		var systemRegistersToAdd []*models.TestAndExerciseSystemRegister
		for _, newSR := range bcpTestAndExercise.BcpTestAndExerciseSystemRegisters {
			if _, exists := existingSystemRegisterIDs[newSR.SystemRegisterID]; !exists {
				systemRegister := &models.TestAndExerciseSystemRegister{
					TestAndExerciseID: id,
					SystemRegisterID:  newSR.SystemRegisterID,
					PlannedRto:        newSR.SystemRegister.RTO,
					PlannedRtoUnitId:  newSR.SystemRegister.RTOUnitID,
					PlannedMtpd:       newSR.SystemRegister.MTPD,
					PlannedMtpdUnitId: newSR.SystemRegister.MTPDUnitID,
					PlannedRpo:        newSR.SystemRegister.RPO,
					PlannedRpoUnitId:  newSR.SystemRegister.RPOUnitID,
					CreatedBYID:       currentUserID,
				}
				systemRegistersToAdd = append(systemRegistersToAdd, systemRegister)
			}
		}

		if len(systemRegistersToAdd) > 0 {
			if err := tes.repo.CreateTestAndExerciseSystemRegisters(systemRegistersToAdd, tx); err != nil {
				tx.Rollback()
				return err
			}
		}

		// Remove System Registers not in the new list
		var deleteSystemRegisterIDs []*uint
		for _, existingSR := range existingTestAndExercise.TestAndExerciseSystemRegisters {
			if _, exists := newSystemRegisterIDs[existingSR.SystemRegisterID]; !exists {
				deleteSystemRegisterIDs = append(deleteSystemRegisterIDs, existingSR.ID)
			}
		}

		if len(deleteSystemRegisterIDs) > 0 {
			if err := tes.repo.DeleteTestAndExerciseSystemRegisters(deleteSystemRegisterIDs, tx); err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	tx.Commit()
	go func() {
		existingTestAndExercise, err := tes.GetFullDetailsByID(id)
		if err != nil {
			return
		}

		changedData := make(map[string]interface{})

		// log changes
		oldData := map[string]interface{}{
			"responsible_user":            fmt.Sprintf("%s %s", oldTestAndExercise.ResponsibleUser.FirstNameEN, oldTestAndExercise.ResponsibleUser.LastNameEN),
			"business_continuity_plan_id": oldTestAndExercise.BusinessContinuityPlanID,
			"bcp_test_and_exercise_id":    oldTestAndExercise.BcpTestAndExerciseID,
			"description":                 oldTestAndExercise.Description,
			"testing_date":                oldTestAndExercise.TestingDateString,
		}

		newData := map[string]interface{}{
			"responsible_user":            fmt.Sprintf("%s %s", existingTestAndExercise.ResponsibleUser.FirstNameEN, existingTestAndExercise.ResponsibleUser.LastNameEN),
			"business_continuity_plan_id": existingTestAndExercise.BusinessContinuityPlanID,
			"bcp_test_and_exercise_id":    existingTestAndExercise.BcpTestAndExerciseID,
			"description":                 existingTestAndExercise.Description,
			"testing_date":                existingTestAndExercise.TestingDateString,
		}
		// Compare other fields and add them to changedData
		fieldChanges := utility.CompareFieldChanges(oldData, newData)
		for k, v := range fieldChanges {
			changedData[k] = v
		}
		// Log activity if there are any changes
		if len(changedData) > 0 {
			err := tes.activityLogSvc.LogActivity(
				models.UpdateAction,
				models.TestAndExercise{}.TableName(),
				id,
				changedData,
				currentUserID,
				"business_continuity",
				"test_and_exercise",
				"test_and_exercise",
				nil,
			)
			if err != nil {
				fields := map[string]interface{}{
					"action":    models.UpdateAction,
					"tableName": models.TestAndExercise{}.TableName(),
					"recordID":  id,
					"userID":    currentUserID,
					"error":     err.Error(),
				}

				utility.LogError("Failed to log activity", fields)
			}
		}
	}()

	return nil
}

func (tes *TestAndExerciseService) Delete(ids []uint, currentUserID uint) ([]uint, error) {
	tx := tes.repo.BeginTransaction()
	var failedIds []uint
	var successfulIds []uint
	var deletedTestAndExercise []map[string]interface{}

	for _, id := range ids {
		testAndExercise, err := tes.GetFullDetailsByID(id, tx)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		deletedTestAndExercise = append(deletedTestAndExercise, map[string]interface{}{
			"id":                id,
			"test_and_exercise": testAndExercise.BcpTestAndExercise.Title,
		})

		if len(testAndExercise.TestAndExerciseProcesses) > 0 {
			var deleteprocessIDs []*uint
			for _, process := range testAndExercise.TestAndExerciseProcesses {
				deleteprocessIDs = append(deleteprocessIDs, process.ID)
			}
			if err := tes.repo.DeleteTestAndExerciseProcesses(deleteprocessIDs, tx); err != nil {
				tx.Rollback()
				return nil, err
			}
		}

		if len(testAndExercise.TestAndExerciseSystemRegisters) > 0 {
			var deleteSystemRegisterIDs []*uint
			for _, sr := range testAndExercise.TestAndExerciseSystemRegisters {
				deleteSystemRegisterIDs = append(deleteSystemRegisterIDs, sr.ID)
			}
			if err := tes.repo.DeleteTestAndExerciseSystemRegisters(deleteSystemRegisterIDs, tx); err != nil {
				tx.Rollback()
				return nil, err
			}
		}

		if len(testAndExercise.TestAndExerciseRecoveryActionPlans) > 0 {
			var deleteRecoveryActionPlanIDs []*uint
			for _, recoveryActionPlan := range testAndExercise.TestAndExerciseRecoveryActionPlans {
				deleteRecoveryActionPlanIDs = append(deleteRecoveryActionPlanIDs, recoveryActionPlan.ID)
			}
			if err := tes.repo.DeleteTestAndExerciseRecoveryActionPlans(deleteRecoveryActionPlanIDs, tx); err != nil {
				tx.Rollback()
				return nil, err
			}
		}
		if len(testAndExercise.TestAndExerciseCommunicationActionPlans) > 0 {
			var deleteCommunicationActionPlanIDs []*uint
			for _, communicationActionPlan := range testAndExercise.TestAndExerciseCommunicationActionPlans {
				deleteCommunicationActionPlanIDs = append(deleteCommunicationActionPlanIDs, communicationActionPlan.ID)
			}
			if err := tes.repo.DeleteTestAndExerciseCommunicationActionPlans(deleteCommunicationActionPlanIDs, tx); err != nil {
				tx.Rollback()
				return nil, err
			}
		}

		// Delete the test and exercise
		if err := tes.repo.Delete(id, tx); err != nil {
			failedIds = append(failedIds, id)
			continue
		} else {
			successfulIds = append(successfulIds, id)
		}
	}
	if len(failedIds) > 0 {
		tx.Rollback()
		return failedIds, nil
	}

	tx.Commit()

	// Log bulk delete activities for successful deletions
	if len(successfulIds) > 0 {
		go func() {
			err := tes.activityLogSvc.LogBulkActivities(
				models.DeleteAction, // Make sure this action constant is defined appropriately
				models.TestAndExercise{}.TableName(),
				successfulIds,
				deletedTestAndExercise,
				currentUserID,
				nil,
			)
			if err != nil {
				// @todo - save to log file
				fmt.Printf("Failed to log bulk delete activity: %v\n", err)
			}
		}()
	}

	return failedIds, nil
}

func (tes *TestAndExerciseService) parseTestAndExerciseDates(dto *businessContinuityDTOs.TestAndExerciseFormDTO) (*businessContinuityDTOs.TestAndExerciseParsedDates, error) {
	parsed := &businessContinuityDTOs.TestAndExerciseParsedDates{}

	var err error

	if &dto.TestingDate != nil {
		parsed.TestingDate, err = utility.ParseDate(dto.TestingDate)
		if err != nil {
			return nil, utility.NewValidationError(fmt.Errorf("invalid testing date"), "invalid_date", "testing date")
		}
	}

	return parsed, nil
}

func (tes TestAndExerciseService) GetTestAndExerciseProcessByID(id uint) (*models.TestAndExerciseProcess, error) {
	TestAndExerciseProcess, err := tes.repo.GetTestAndExerciseProcessByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	return TestAndExerciseProcess, nil
}

func (tes TestAndExerciseService) GetTestAndExerciseSystemRegisterByID(id uint) (*models.TestAndExerciseSystemRegister, error) {
	testAndExerciseSystemRegister, err := tes.repo.GetTestAndExerciseSystemRegisterByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utility.NewResourceNotFoundError(err, id)
		}
		return nil, err
	}

	return testAndExerciseSystemRegister, nil
}

func (tes *TestAndExerciseService) ExportTestAndExercise(data []*businessContinuityDTOs.TestAndExerciseListDTO, currentUserID uint) (*os.File, error) {
	userSetting, err := tes.userSettingSvc.GetUserSetting(currentUserID)
	if err != nil {
		return nil, err
	}

	var headers []string
	if userSetting.Language == "ar" {
		headers = []string{
			"SL No.",
			"Reference Number",
			"Title",
			"Business Function",
			"Testing Date",
		}
	} else {
		headers = []string{
			"SL No.",
			"Reference Number",
			"Title",
			"Business Function",
			"Testing Date",
		}
	}

	rows := make([][]interface{}, len(data))
	for i, testAndExercise := range data {
		if userSetting.Language == "ar" {
			rows[i] = []interface{}{
				i + 1,
				testAndExercise.ReferenceNumber,
				testAndExercise.BcpTestAndExerciseTitle,
				testAndExercise.BusinessFunctionTitle,
				testAndExercise.TestingDate,
			}
		} else {
			rows[i] = []interface{}{
				i + 1,
				testAndExercise.ReferenceNumber,
				testAndExercise.BcpTestAndExerciseTitle,
				testAndExercise.BusinessFunctionTitle,
				testAndExercise.TestingDate,
			}
		}
	}

	return utility.GenerateExcelFile(headers, rows)
}

func (tes *TestAndExerciseService) GetTestResultSummary(testAndExerciseID uint) (*businessContinuityDTOs.TestResultSummaryDTO, error) {
	items, err := tes.repo.GetTestResultSystemRegistersByTestAndExerciseID(testAndExerciseID)
	if err != nil {
		return nil, err
	}

	optionMap := make(map[uint]*businessContinuityDTOs.TestResultSummaryOptionDTO)
	var optionOrder []uint

	for _, sr := range items {
		if sr.TestResultOption == nil {
			continue
		}
		optID := sr.TestResultOption.ID
		if _, exists := optionMap[optID]; !exists {
			optionMap[optID] = &businessContinuityDTOs.TestResultSummaryOptionDTO{
				TitleEn:         sr.TestResultOption.TitleEn,
				TitleAr:         sr.TestResultOption.TitleAr,
				ColorCode:       sr.TestResultOption.ColorCode,
				Type:            sr.TestResultOption.Type,
				Count:           0,
				SystemRegisters: []*businessContinuityDTOs.TestResultSummarySystemRegisterDTO{},
			}
			optionOrder = append(optionOrder, optID)
		}
		opt := optionMap[optID]
		opt.Count++
		if sr.SystemRegister != nil {
			opt.SystemRegisters = append(opt.SystemRegisters, &businessContinuityDTOs.TestResultSummarySystemRegisterDTO{
				ID:      sr.SystemRegister.ID,
				TitleEn: sr.SystemRegister.TitleEn,
				TitleAr: sr.SystemRegister.TitleAr,
			})
		}
	}

	var options []*businessContinuityDTOs.TestResultSummaryOptionDTO
	for _, id := range optionOrder {
		options = append(options, optionMap[id])
	}
	if options == nil {
		options = []*businessContinuityDTOs.TestResultSummaryOptionDTO{}
	}

	return &businessContinuityDTOs.TestResultSummaryDTO{
		Options: options,
	}, nil
}
