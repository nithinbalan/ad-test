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
	settingServices "jethur/grc/settings/services"
	taskServices "jethur/grc/task/services"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type TestAndExerciseServiceInterface interface {
	GetAll(params *businessContinuityDTOs.TestAndExerciseListQueryParamsDTO, currentUserID uint) (*businessContinuityDTOs.TestAndExerciseListResponseDTO, error)
	GetFullDetailsByID(id uint, tx ...utility.DatabaseTransaction) (*businessContinuityDTOs.TestAndExerciseDetailDTO, error)
	GetByID(id uint) (*models.TestAndExercise, error)
	Create(currentUserID uint, businessContinuityDTOs *businessContinuityDTOs.TestAndExerciseFormDTO) (uint, error)
	Update(id uint, currentUserID uint, TestAndExercise *businessContinuityDTOs.TestAndExerciseFormDTO) error
	MarkAsCompleted(id uint, currentUserID uint) error
	Delete(ids []uint, currentUserID uint) ([]uint, error)

	GetTestAndExerciseProcessByID(id uint) (*models.TestAndExerciseProcess, error)
	GetTestAndExerciseSystemRegisterByID(id uint) (*models.TestAndExerciseSystemRegister, error)
	ExportTestAndExercise(data []*businessContinuityDTOs.TestAndExerciseListDTO, currentUserID uint) (*os.File, error)
	ExportTestAndExercisePDF(id uint, currentUserID uint) (string, error)

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
	reportTemplateSvc                settingServices.ReportTemplateSettingServiceInterface
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
	reportTemplateSvc settingServices.ReportTemplateSettingServiceInterface,
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
		reportTemplateSvc:                reportTemplateSvc,
	}
}

func (tes *TestAndExerciseService) GetAll(params *businessContinuityDTOs.TestAndExerciseListQueryParamsDTO, currentUserID uint) (*businessContinuityDTOs.TestAndExerciseListResponseDTO, error) {
	if currentUserID != 0 {
		isCyberUser, err := tes.userSvc.IsCyberUser(currentUserID)
		if err != nil {
			return nil, err
		}
		params.IsCyberUser = isCyberUser
	}
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
		tx.Rollback()
		return 0, err
	}

	dates, err := tes.parseTestAndExerciseDates(testAndExerciseFormDTO)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	isCyberUser := false
	if currentUserID != 0 {
		isCyberUser, err = tes.userSvc.IsCyberUser(currentUserID)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	TestAndExercise := &models.TestAndExercise{
		ReferenceNumber:          referenceNumber,
		BusinessContinuityPlanID: testAndExerciseFormDTO.BusinessContinuityPlanID,
		BcpTestAndExerciseID:     testAndExerciseFormDTO.BcpTestAndExerciseID,
		Objective:                testAndExerciseFormDTO.Objective,
		Scenario:                 testAndExerciseFormDTO.Scenario,
		SuccessCriteria:          testAndExerciseFormDTO.SuccessCriteria,
		ResponsibleUserID:        &testAndExerciseFormDTO.ResponsibleUserID,
		ResponsibleUnitID:        &responsibleUnitID,
		TestingDate:              dates.TestingDate,
		Description:              testAndExerciseFormDTO.Description,
		TestAndExerciseStatusID:  testAndExerciseStatus.ID,
		CreatedBYID:              currentUserID,
		IsCyber:                  isCyberUser,
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
		tx.Rollback()
		return utility.NewResourceNotFoundError(err, id)
	}

	oldTestAndExercise, err := tes.GetFullDetailsByID(id, tx)
	if err != nil {
		tx.Rollback()
		return err
	}

	if existingTestAndExercise.ResponsibleUserID == nil || *existingTestAndExercise.ResponsibleUserID != testAndExerciseFormDTO.ResponsibleUserID {
		responsibleUnitID, err := tes.userSvc.GetUserPrimaryUnitID(testAndExerciseFormDTO.ResponsibleUserID)
		if err != nil {
			tx.Rollback()
			return utility.NewValidationError(fmt.Errorf("organization unit not set for responsible user"), "not_exist", "risk responsible unit")
		}
		existingTestAndExercise.ResponsibleUnitID = &responsibleUnitID
	}

	dates, err := tes.parseTestAndExerciseDates(testAndExerciseFormDTO)
	if err != nil {
		tx.Rollback()
		return err
	}

	existingTestAndExercise.BusinessContinuityPlanID = testAndExerciseFormDTO.BusinessContinuityPlanID
	existingTestAndExercise.ResponsibleUserID = &testAndExerciseFormDTO.ResponsibleUserID
	existingTestAndExercise.Objective = testAndExerciseFormDTO.Objective
	existingTestAndExercise.Scenario = testAndExerciseFormDTO.Scenario
	existingTestAndExercise.SuccessCriteria = testAndExerciseFormDTO.SuccessCriteria
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
			"Test Name",
			"Comment",
			"Business Function",
			"Responsible User",
			"Testing Date",
			"Status",
			"Created Date",
			"Created By",
		}
	} else {
		headers = []string{
			"SL No.",
			"Reference Number",
			"Test Name",
			"Comment",
			"Business Function",
			"Responsible User",
			"Testing Date",
			"Status",
			"Created Date",
			"Created By",
		}
	}

	rows := make([][]interface{}, len(data))
	for i, testAndExercise := range data {
		createdDate := testAndExercise.CreatedAt
		if t, err := time.Parse(time.RFC3339, testAndExercise.CreatedAt); err == nil {
			createdDate = t.Format("2006-01-02")
		} else if t, err := time.Parse("2006-01-02T15:04:05.999Z", testAndExercise.CreatedAt); err == nil {
			createdDate = t.Format("2006-01-02")
		}
		if userSetting.Language == "ar" {
			rows[i] = []interface{}{
				i + 1,
				testAndExercise.ReferenceNumber,
				testAndExercise.BcpTestAndExerciseTitle,
				utility.GetStringValue(testAndExercise.Description),
				testAndExercise.BusinessFunctionTitle,
				fmt.Sprintf("%s %s", utility.GetStringValue(testAndExercise.ResponsibleUserFirstNameAr), utility.GetStringValue(testAndExercise.ResponsibleUserLastNameAr)),
				testAndExercise.TestingDate,
				testAndExercise.TestAndExerciseStatusTitleAr,
				createdDate,
				fmt.Sprintf("%s %s", testAndExercise.CreatedUserFirstNameAr, testAndExercise.CreatedUserLastNameAr),
			}
		} else {
			rows[i] = []interface{}{
				i + 1,
				testAndExercise.ReferenceNumber,
				testAndExercise.BcpTestAndExerciseTitle,
				utility.GetStringValue(testAndExercise.Description),
				testAndExercise.BusinessFunctionTitle,
				fmt.Sprintf("%s %s", utility.GetStringValue(testAndExercise.ResponsibleUserFirstNameEn), utility.GetStringValue(testAndExercise.ResponsibleUserLastNameEn)),
				testAndExercise.TestingDate,
				testAndExercise.TestAndExerciseStatusTitleEn,
				createdDate,
				fmt.Sprintf("%s %s", testAndExercise.CreatedUserFirstNameEn, testAndExercise.CreatedUserLastNameEn),
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

func (tes *TestAndExerciseService) MarkAsCompleted(id uint, currentUserID uint) error {
	tx := tes.repo.BeginTransaction()

	testAndExercise, err := tes.repo.GetByID(id)
	if err != nil {
		tx.Rollback()
		return utility.NewResourceNotFoundError(err, id)
	}

	completedStatus, err := tes.testAndExerciseStatusSvc.GetByType("completed")
	if err != nil {
		tx.Rollback()
		return err
	}

	businessContinuityPlan, err := tes.businessContinuityPlanSvc.GetByID(testAndExercise.BusinessContinuityPlanID)
	if err != nil {
		tx.Rollback()
		return err
	}

	oldTestAndExerciseStatusID := testAndExercise.TestAndExerciseStatusID
	oldIsTestAndExerciseCompleted := businessContinuityPlan.IsTestAndExerciseCompleted

	testAndExercise.TestAndExerciseStatusID = completedStatus.ID
	testAndExercise.UpdatedBYID = &currentUserID
	if err := tes.repo.Update(testAndExercise, tx); err != nil {
		tx.Rollback()
		return err
	}

	businessContinuityPlan.IsTestAndExerciseCompleted = true
	businessContinuityPlan.UpdatedBYID = &currentUserID
	if err := tes.businessContinuityPlanSvc.UpdateRecoveryObjectives(businessContinuityPlan, tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	go func() {
		changedData := utility.CompareFieldChanges(
			map[string]interface{}{
				"test_and_exercise_status_id":    oldTestAndExerciseStatusID,
				"is_test_and_exercise_completed": oldIsTestAndExerciseCompleted,
			},
			map[string]interface{}{
				"test_and_exercise_status_id":    completedStatus.ID,
				"is_test_and_exercise_completed": true,
			},
		)

		if len(changedData) == 0 {
			return
		}

		if err := tes.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.TestAndExercise{}.TableName(),
			id,
			changedData,
			currentUserID,
			"business_continuity",
			"test_and_exercise",
			"test_and_exercise",
			nil,
		); err != nil {
			utility.LogError("Failed to log activity", map[string]interface{}{
				"action":    models.UpdateAction,
				"tableName": models.TestAndExercise{}.TableName(),
				"recordID":  id,
				"userID":    currentUserID,
				"error":     err.Error(),
			})
		}
	}()

	return nil
}

func (tes *TestAndExerciseService) ExportTestAndExercisePDF(id uint, currentUserID uint) (string, error) {
	// Validate that template images are configured
	if _, err := tes.reportTemplateSvc.GetImagePath("logo"); err != nil {
		return "", utility.NewForbiddenError("logo_image_missing")
	}
	if _, err := tes.reportTemplateSvc.GetImagePath("cover-page-bg"); err != nil {
		return "", utility.NewForbiddenError("cover_image_missing")
	}

	// Get complete Test & Exercise details
	testExerciseDTO, err := tes.GetFullDetailsByID(id)
	if err != nil {
		return "", err
	}

	// Format dates
	createdAt := testExerciseDTO.CreatedAt.Format("02 Jan 2006")
	generatedAt := time.Now().Format("Jan 02, 2006")
	preparedFirst := testExerciseDTO.CreatedByUser.FirstNameEN
	preparedLast := testExerciseDTO.CreatedByUser.LastNameEN

	var responsibleFirst, responsibleLast string
	if testExerciseDTO.ResponsibleUser != nil {
		responsibleFirst = testExerciseDTO.ResponsibleUser.FirstNameEN
		responsibleLast = testExerciseDTO.ResponsibleUser.LastNameEN
	}

	var updatedFirst, updatedLast, updatedAtStr *string
	if testExerciseDTO.UpdatedByUser != nil {
		f := testExerciseDTO.UpdatedByUser.FirstNameEN
		l := testExerciseDTO.UpdatedByUser.LastNameEN
		updatedFirst = &f
		updatedLast = &l
	}
	if testExerciseDTO.UpdatedAt != nil {
		s := testExerciseDTO.UpdatedAt.Format("02 Jan 2006")
		updatedAtStr = &s
	}

	// Extract Test Result data
	testObservations := ""
	recommendations := ""
	findings := ""
	var performanceMeasures []businessContinuityDTOs.TestAndExercisePerformanceMeasureExportDTO
	var systemRegisters []businessContinuityDTOs.TestAndExerciseSystemRegisterExportDTO

	if testExerciseDTO.TestAndExerciseTestResult != nil {
		if testExerciseDTO.TestAndExerciseTestResult.TestObservations != nil {
			testObservations = utility.StripHTMLTags(*testExerciseDTO.TestAndExerciseTestResult.TestObservations)
		}
		if testExerciseDTO.TestAndExerciseTestResult.Recommendations != nil {
			recommendations = utility.StripHTMLTags(*testExerciseDTO.TestAndExerciseTestResult.Recommendations)
		}
		if testExerciseDTO.TestAndExerciseTestResult.Findings != nil {
			findings = utility.StripHTMLTags(*testExerciseDTO.TestAndExerciseTestResult.Findings)
		}

		// Build Performance Measures table
		for _, pm := range testExerciseDTO.TestAndExerciseTestResult.PerformanceMeasures {
			criterion := ""
			if pm.Criterion != nil {
				criterion = utility.StripHTMLTags(*pm.Criterion)
			}
			target := ""
			if pm.Target != nil {
				target = utility.StripHTMLTags(*pm.Target)
			}
			achievedResult := ""
			if pm.AchievedResult != nil {
				achievedResult = utility.StripHTMLTags(*pm.AchievedResult)
			}
			remarks := ""
			if pm.Remarks != nil {
				remarks = utility.StripHTMLTags(*pm.Remarks)
			}

			performanceMeasures = append(performanceMeasures, businessContinuityDTOs.TestAndExercisePerformanceMeasureExportDTO{
				Criterion:      criterion,
				Target:         target,
				AchievedResult: achievedResult,
				Remarks:        remarks,
			})
		}

		// Build System Registers table
		for _, sr := range testExerciseDTO.TestAndExerciseTestResult.TestResultSystemRegisters {
			systemTitle := sr.SystemRegister.TitleEn
			testCompletionDate := sr.TestCompletionDateString

			commFeedbackTitle := sr.CommunicationFeedbackOption.TitleEn
			commFeedbackColor := sr.CommunicationFeedbackOption.ColorCode

			recoveryFeedbackTitle := sr.RecoveryFeedbackOption.TitleEn
			recoveryFeedbackColor := sr.RecoveryFeedbackOption.ColorCode

			testResultTitle := sr.TestResultOption.TitleEn
			testResultColor := sr.TestResultOption.ColorCode

			systemRegisters = append(systemRegisters, businessContinuityDTOs.TestAndExerciseSystemRegisterExportDTO{
				SystemTitle:                      systemTitle,
				TestCompletionDate:               testCompletionDate,
				CommunicationFeedbackOptionTitle: commFeedbackTitle,
				CommunicationFeedbackOptionColor: commFeedbackColor,
				RecoveryFeedbackOptionTitle:      recoveryFeedbackTitle,
				RecoveryFeedbackOptionColor:      recoveryFeedbackColor,
				TestResultOptionTitle:            testResultTitle,
				TestResultOptionColor:            testResultColor,
			})
		}
	}

	// Build Recovery Action Plans table
	var recoveryActionPlans []businessContinuityDTOs.TestAndExerciseRecoveryActionPlanExportDTO
	for _, rap := range testExerciseDTO.TestAndExerciseRecoveryActionPlans {
		processTitle := ""
		if rap.Process != nil {
			processTitle = rap.Process.TitleEn
		}
		systemTitle := ""
		if rap.SystemRegister != nil {
			systemTitle = rap.SystemRegister.TitleEn
		}
		bcpPhaseTitle := rap.BcpPhase.TitleEn
		severityLevel := ""
		if rap.SeverityLevel != nil {
			severityLevel = rap.SeverityLevel.TitleEn
		}
		responsibleUser := fmt.Sprintf("%s %s", rap.ResponsibleUser.FirstNameEN, rap.ResponsibleUser.LastNameEN)
		responsibleUnit := rap.ResponsibleUnit.TitleEn
		targetDate := ""
		if rap.TargetDateString != nil {
			targetDate = *rap.TargetDateString
		}

		recoveryActionPlans = append(recoveryActionPlans, businessContinuityDTOs.TestAndExerciseRecoveryActionPlanExportDTO{
			ProcessTitle:    processTitle,
			SystemTitle:     systemTitle,
			BcpPhaseTitle:   bcpPhaseTitle,
			SeverityLevel:   severityLevel,
			ResponsibleUser: responsibleUser,
			ResponsibleUnit: responsibleUnit,
			TargetDate:      targetDate,
		})
	}

	// Build Communication Action Plans table
	var communicationActionPlans []businessContinuityDTOs.TestAndExerciseCommunicationActionPlanExportDTO
	for _, cap := range testExerciseDTO.TestAndExerciseCommunicationActionPlans {
		processTitle := ""
		if cap.Process != nil {
			processTitle = cap.Process.TitleEn
		}
		systemTitle := ""
		if cap.SystemRegister != nil {
			systemTitle = cap.SystemRegister.TitleEn
		}
		bcpPhaseTitle := cap.BcpPhase.TitleEn
		responsibleUser := fmt.Sprintf("%s %s", cap.ResponsibleUser.FirstNameEN, cap.ResponsibleUser.LastNameEN)
		responsibleUnit := cap.ResponsibleUnit.TitleEn
		frequency := ""
		if cap.Frequency != nil {
			frequency = *cap.Frequency
		}
		recipient := ""
		if cap.Recipient != nil {
			recipient = *cap.Recipient
		}
		channel := ""
		if cap.Channel != nil {
			channel = *cap.Channel
		}

		communicationActionPlans = append(communicationActionPlans, businessContinuityDTOs.TestAndExerciseCommunicationActionPlanExportDTO{
			ProcessTitle:    processTitle,
			SystemTitle:     systemTitle,
			BcpPhaseTitle:   bcpPhaseTitle,
			ResponsibleUser: responsibleUser,
			ResponsibleUnit: responsibleUnit,
			Frequency:       frequency,
			Recipient:       recipient,
			Channel:         channel,
		})
	}

	// Prepare view model
	objective := ""
	if testExerciseDTO.Objective != nil {
		objective = utility.StripHTMLTags(*testExerciseDTO.Objective)
	}
	scenario := ""
	if testExerciseDTO.Scenario != nil {
		scenario = utility.StripHTMLTags(*testExerciseDTO.Scenario)
	}
	successCriteria := ""
	if testExerciseDTO.SuccessCriteria != nil {
		successCriteria = utility.StripHTMLTags(*testExerciseDTO.SuccessCriteria)
	}
	description := ""
	if testExerciseDTO.Description != nil {
		description = utility.StripHTMLTags(*testExerciseDTO.Description)
	}

	responsibleUnitTitle := ""
	if testExerciseDTO.BusinessContinuityPlan.ResponsibleUnit.ID != 0 {
		responsibleUnitTitle = testExerciseDTO.BusinessContinuityPlan.ResponsibleUnit.TitleEn
	}

	viewModel := businessContinuityDTOs.TestAndExerciseExportPDFDTO{
		Title:                      testExerciseDTO.BcpTestAndExercise.Title,
		ReferenceNumber:            testExerciseDTO.ReferenceNumber,
		CreatedAt:                  createdAt,
		GeneratedAt:                generatedAt,
		CreatedUserFirstNameEn:     preparedFirst,
		CreatedUserLastNameEn:      preparedLast,
		UpdatedUserFirstNameEn:     updatedFirst,
		UpdatedUserLastNameEn:      updatedLast,
		UpdatedAtString:            updatedAtStr,
		BcpTestAndExerciseTitle:    testExerciseDTO.BcpTestAndExercise.Title,
		Objective:                  objective,
		Scenario:                   scenario,
		SuccessCriteria:            successCriteria,
		TestingDate:                testExerciseDTO.TestingDateString,
		Description:                description,
		ResponsibleUserFirstNameEn: responsibleFirst,
		ResponsibleUserLastNameEn:  responsibleLast,
		ResponsibleUnitTitleEn:     responsibleUnitTitle,
		StatusTitleEn:              testExerciseDTO.TestAndExerciseStatus.TitleEn,
		StatusColorCode:            testExerciseDTO.TestAndExerciseStatus.ColorCode,
		StatusType:                 testExerciseDTO.TestAndExerciseStatus.Type,
		TestObservations:           testObservations,
		Recommendations:            recommendations,
		Findings:                   findings,
		PerformanceMeasures:        performanceMeasures,
		SystemRegisters:            systemRegisters,
		RecoveryActionPlans:        recoveryActionPlans,
		CommunicationActionPlans:   communicationActionPlans,
	}

	// Body
	bodyTemplates := []string{"business-continuity/templates/pdf/test_and_exercise_body.html"}
	htmlContent, err := utility.GeneratePDFContent(viewModel, bodyTemplates, "", "")
	if err != nil {
		return "", fmt.Errorf("failed to generate HTML content: %w", err)
	}

	// Cover
	coverTemplates := []string{"business-continuity/templates/pdf/test_and_exercise_cover.html"}
	coverContent, err := utility.GeneratePDFContent(viewModel, coverTemplates, "", "")
	if err != nil {
		return "", fmt.Errorf("failed to generate cover content: %w", err)
	}

	// Header
	headerTemplates := []string{"business-continuity/templates/pdf/test_and_exercise_header.html"}
	headerContent, err := utility.GeneratePDFContent(viewModel, headerTemplates, "", "")
	if err != nil {
		return "", fmt.Errorf("failed to generate header content: %w", err)
	}

	// Footer
	footerTemplates := []string{"business-continuity/templates/pdf/test_and_exercise_footer.html"}
	footerContent, err := utility.GeneratePDFContent(viewModel, footerTemplates, "", "")
	if err != nil {
		return "", fmt.Errorf("failed to generate footer content: %w", err)
	}

	pdfFileName := fmt.Sprintf("%s.pdf", viewModel.ReferenceNumber)
	if err := utility.GeneratePDFWithCoverHeaderFooter(coverContent, htmlContent, pdfFileName, headerContent, footerContent); err != nil {
		return "", fmt.Errorf("failed to generate PDF: %w", err)
	}

	return pdfFileName, nil
}
