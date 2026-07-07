package services

import (
	"errors"
	"fmt"
	"time"

	bcDTOs "jethur/grc/business-continuity/dtos"
	bcRepos "jethur/grc/business-continuity/repositories"
	logServices "jethur/grc/logs/services"
	"jethur/grc/models"
	"jethur/utility"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type BCMStaffAllocationServiceInterface interface {
	GetFullDetailsByProcessID(processID uint) ([]*bcDTOs.BCMStaffAllocationDetailDTO, error)
	DeleteByProcessID(processID uint, tx ...utility.DatabaseTransaction) error
	UpdateAllocation(processID uint, currentUserID uint, dto *bcDTOs.BCMStaffAllocationUpdateFormDTO) error
	GenerateAllocations(processID uint, currentUserID uint) error
}

type BCMStaffAllocationService struct {
	repo           bcRepos.BCMStaffAllocationRepositoryInterface
	activityLogSvc logServices.ActivityLogServiceInterface
}

func NewBCMStaffAllocationService(
	repo bcRepos.BCMStaffAllocationRepositoryInterface,
	activityLogSvc logServices.ActivityLogServiceInterface,
) BCMStaffAllocationServiceInterface {
	return &BCMStaffAllocationService{
		repo:           repo,
		activityLogSvc: activityLogSvc,
	}
}

func (s *BCMStaffAllocationService) DeleteByProcessID(processID uint, tx ...utility.DatabaseTransaction) error {
	return s.repo.DeleteByProcessID(processID, tx...)
}

func (s *BCMStaffAllocationService) GetFullDetailsByProcessID(processID uint) ([]*bcDTOs.BCMStaffAllocationDetailDTO, error) {
	allocations, err := s.repo.GetFullDetailsByProcessID(processID)
	if err != nil {
		return nil, err
	}

	var result []*bcDTOs.BCMStaffAllocationDetailDTO
	for _, allocation := range allocations {
		var dto bcDTOs.BCMStaffAllocationDetailDTO
		if err := mapstructure.Decode(allocation, &dto); err != nil {
			return nil, err
		}
		result = append(result, &dto)
	}

	return result, nil
}

// UpdateAllocation replaces all user assignments for a given bcm_staff_allocation_outage_bucket_id.
func (s *BCMStaffAllocationService) UpdateAllocation(processID uint, currentUserID uint, dto *bcDTOs.BCMStaffAllocationUpdateFormDTO) error {
	tx := s.repo.BeginTransaction()

	if err := s.repo.DeleteUsersByOutageBucketID(dto.BCMStaffAllocationOutageBucketID, tx); err != nil {
		tx.Rollback()
		return err
	}

	for _, u := range dto.Users {
		user := &models.BCMStaffAllocationOutageBucketUser{
			BCMStaffAllocationOutageBucketID: dto.BCMStaffAllocationOutageBucketID,
			UserID:                           u.UserID,
			IsRemote:                         u.IsRemote,
			CreatedBYID:                      currentUserID,
			CreatedAt:                        time.Now(),
		}
		if err := s.repo.CreateUser(user, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	go func() {
		_ = s.activityLogSvc.LogActivity(
			models.UpdateAction,
			models.BCMStaffAllocation{}.TableName(),
			processID,
			map[string]interface{}{"bcm_staff_allocation_outage_bucket_id": dto.BCMStaffAllocationOutageBucketID, "users": dto.Users},
			currentUserID,
			"business-continuity", "bcm-bia", "staff-allocation", nil,
		)
	}()

	return nil
}

// GenerateAllocations creates one BCMStaffAllocation per ProcessEmployeeAnalysis,
// then one BCMStaffAllocationOutageBucket child row per outage bucket.
//
// Logic:
//  1. Fetch BIA Rating → RspFromValue + RspUnit (preloaded)
//  2. Fetch all active outage buckets ordered by [order] ASC (with DurationUnit preloaded)
//  3. Convert RSP to minutes (RspFromValue * RspUnit.Minutes).
//     Find which bucket the RSP falls into (bucketFrom <= rspMinutes <= bucketTo).
//     That bucket is the first ACTIVE interval; all buckets before it have IsSkipped = true.
//     Example: RSP = 2 Days (2880 min) falls into 1-3 Days bucket → skip 0-8h, 8-24h
//  4. For each ProcessEmployeeAnalysis (ordered by ID ASC):
//     a. Create one BCMStaffAllocation parent row
//     b. For each bucket: create one BCMStaffAllocationOutageBucket child row
//     Users are assigned separately via UpdateAllocation.
func (s *BCMStaffAllocationService) GenerateAllocations(processID uint, currentUserID uint) error {
	// Step 1: get RSP from BIA rating
	rating, err := s.repo.GetRatingByProcessID(processID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utility.NewValidationError(
				fmt.Errorf("no BIA rating found for process %d, complete BIA first", processID),
				"not_found",
				"bia_rating",
			)
		}
		return err
	}

	// Step 2: get all active buckets sorted by Order ASC
	buckets, err := s.repo.GetOutageBuckets()
	if err != nil {
		return err
	}
	if len(buckets) == 0 {
		return utility.NewValidationError(
			fmt.Errorf("no active outage buckets configured"),
			"not_found",
			"outage_buckets",
		)
	}

	// Step 3: determine first active bucket order based on RSP
	// Convert RSP to minutes, find which bucket the RSP falls into,
	// then skip all buckets before it.
	// Example: RSP = 2 Days (2880 min) falls into "1-3 Days" bucket → skip 0-8h, 8-24h
	firstActiveOrder := buckets[0].Order
	if rating.RspFromValue != nil && rating.RspUnit != nil {
		rspMinutes := *rating.RspFromValue * rating.RspUnit.Minutes
		for _, b := range buckets {
			if b.DurationUnit != nil {
				bucketFromMinutes := b.FromValue * b.DurationUnit.Minutes
				bucketToMinutes := b.ToValue * b.DurationUnit.Minutes
				if rspMinutes >= bucketFromMinutes && rspMinutes <= bucketToMinutes {
					firstActiveOrder = b.Order
					break
				}
			}
		}
	}

	// Step 4: get all employee analysis rows for this process ordered by ID ASC
	analyses, err := s.repo.GetAnalysisByProcessID(processID)
	if err != nil {
		return err
	}
	if len(analyses) == 0 {
		return utility.NewValidationError(
			fmt.Errorf("no employee analysis records found for process %d", processID),
			"not_found",
			"employee_analysis",
		)
	}

	tx := s.repo.BeginTransaction()

	if err := s.repo.DeleteByProcessID(processID, tx); err != nil {
		tx.Rollback()
		return err
	}

	for _, analysis := range analyses {
		peaID := analysis.ID

		// Create one parent row per ProcessEmployeeAnalysis
		parent := &models.BCMStaffAllocation{
			ProcessID:                 processID,
			ProcessEmployeeAnalysisID: &peaID,
			CreatedBYID:               currentUserID,
			CreatedAt:                 time.Now(),
		}
		allocationID, err := s.repo.Create(parent, tx)
		if err != nil {
			tx.Rollback()
			return err
		}

		// Create one child row per outage bucket
		for _, bucket := range buckets {
			bucketID := bucket.ID
			isSkipped := bucket.Order < firstActiveOrder

			child := &models.BCMStaffAllocationOutageBucket{
				BCMStaffAllocationID: allocationID,
				OutageBucketID:       &bucketID,
				IsSkipped:            isSkipped,
				CreatedBYID:          currentUserID,
				CreatedAt:            time.Now(),
			}
			if _, err := s.repo.CreateOutageBucket(child, tx); err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	go func() {
		_ = s.activityLogSvc.LogActivity(
			models.CreateAction,
			models.BCMStaffAllocation{}.TableName(),
			processID,
			map[string]interface{}{"process_id": processID, "auto_generated": true},
			currentUserID,
			"business-continuity", "bcm-bia", "staff-allocation", nil,
		)
	}()

	return nil
}
