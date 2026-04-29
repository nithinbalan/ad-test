package services

import (
	"errors"

	commonDTOs "jethur/grc/common/dtos"
	logDTOs "jethur/grc/logs/dtos"
	"jethur/grc/logs/repositories"
	orgDTOs "jethur/grc/organization/dtos"

	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type EmailLogServiceInterface interface {
	GetAll(params *logDTOs.EmailLogListQueryParamsDTO) (*logDTOs.EmailLogListResponseDTO, error)
	GetFullDetailsByID(id uint) (*logDTOs.EmailLogDetailDTO, error)
	GetSummary(params *logDTOs.EmailLogListQueryParamsDTO) (*logDTOs.EmailLogSummaryDTO, error)
}

type EmailLogService struct {
	repo repositories.EmailLogRepositoryInteface
}

func NewEmailLogService(repo repositories.EmailLogRepositoryInteface) EmailLogServiceInterface {
	return &EmailLogService{
		repo: repo,
	}
}

func (els *EmailLogService) GetAll(params *logDTOs.EmailLogListQueryParamsDTO) (*logDTOs.EmailLogListResponseDTO, error) {
	emails, totalCount, err := els.repo.GetAll(params)
	if err != nil {
		return nil, err
	}

	totalPages := 0
	if params.PageSize > 0 {
		totalPages = int((totalCount + int64(params.PageSize) - 1) / int64(params.PageSize))
	}

	items := make([]*logDTOs.EmailLogDTO, 0, len(emails))

	for _, email := range emails {
		// var sentAt *string
		// if email.SentAt != nil {
		// 	formatted := email.SentAt.Format("2006-01-02")
		// 	sentAt = &formatted
		// }

		// var lastAttemptAt *string
		// if email.LastAttemptAt != nil {
		// 	formatted := email.LastAttemptAt.Format("2006-01-02")
		// 	lastAttemptAt = &formatted
		// }
		dto := &logDTOs.EmailLogDTO{
			ID:                email.ID,
			ModuleID:          email.ModuleID,
			SubModuleID:       email.SubModuleID,
			Subject:           email.Subject,
			IsPermanentFailed: email.IsPermanentFailed,
			AttemptCount:      email.AttemptCount,
			MaxAttempts:       email.MaxAttempts,
			LastAttemptAt:     email.LastAttemptAt,
			ErrorMessage:      email.ErrorMessage,
			SentAt:            email.SentAt,
			CreatedAt:         email.CreatedAt,
			Recipients:        make([]logDTOs.EmailLogRecipientDTO, 0, len(email.Recipients)),
		}

		if email.EmailStatus != nil {
			dto.EmailStatusID = &email.EmailStatus.ID
			dto.EmailStatusTitleEn = &email.EmailStatus.TitleEn
			dto.EmailStatusTitleAr = &email.EmailStatus.TitleAr
			dto.EmailStatusColorCode = &email.EmailStatus.ColorCode
			dto.EmailStatusType = &email.EmailStatus.Type
		}

		for _, r := range email.Recipients {
			dto.Recipients = append(dto.Recipients, logDTOs.EmailLogRecipientDTO{
				ID:            r.ID,
				RecipientType: r.RecipientType,
				EmailAddress:  r.EmailAddress,
				DisplayName:   r.DisplayName,
				IsExternal:    r.IsExternal,
			})
		}

		items = append(items, dto)
	}

	return &logDTOs.EmailLogListResponseDTO{
		Items: items,
		Pagination: commonDTOs.PaginationInfoDTO{
			CurrentPage: params.Page,
			PageSize:    params.PageSize,
			TotalPages:  totalPages,
			TotalCount:  totalCount,
		},
	}, nil
}

func (els *EmailLogService) GetFullDetailsByID(id uint) (*logDTOs.EmailLogDetailDTO, error) {
	email, err := els.repo.GetFullDetailsByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, err
	}

	var dto logDTOs.EmailLogDetailDTO
	err = mapstructure.Decode(email, &dto)
	if err != nil {
		return nil, err
	}

	dto.CreatedAt = email.CreatedAt
	dto.CreatedAtString = email.CreatedAt.Format("2006-01-02")

	if email.LastAttemptAt != nil {
		dateStr := email.LastAttemptAt.Format("2006-01-02")
		dto.LastAttemptAtString = &dateStr
	}

	if email.NextRetryAt != nil {
		dateStr := email.NextRetryAt.Format("2006-01-02")
		dto.NextRetryAtString = &dateStr
	}

	if email.LockedAt != nil {
		dateStr := email.LockedAt.Format("2006-01-02")
		dto.LockedAtString = &dateStr
	}

	if email.SentAt != nil {
		dateStr := email.SentAt.Format("2006-01-02")
		dto.SentAtString = &dateStr
	}

	dto.Recipients = make([]logDTOs.EmailLogRecipientDTO, 0)
	dto.Attachments = make([]logDTOs.EmailLogAttachmentDTO, 0)

	if email.Module != nil {
		dto.Module = &logDTOs.ModuleMinimalDTO{
			ID:         email.Module.ID,
			Identifier: email.Module.Identifier,
			TitleEn:    email.Module.TitleEn,
			TitleAr:    email.Module.TitleAr,
		}
	}

	if email.SubModule != nil {
		dto.SubModule = &logDTOs.ModuleMinimalDTO{
			ID:         email.SubModule.ID,
			Identifier: email.SubModule.Identifier,
			TitleEn:    email.SubModule.TitleEn,
			TitleAr:    email.SubModule.TitleAr,
		}
	}

	for _, recipient := range email.Recipients {
		recDto := logDTOs.EmailLogRecipientDTO{
			ID:            recipient.ID,
			RecipientType: recipient.RecipientType,
			EmailAddress:  recipient.EmailAddress,
			DisplayName:   recipient.DisplayName,
			IsExternal:    recipient.IsExternal,
		}

		if recipient.User != nil {
			recDto.User = &orgDTOs.UserMinimalDTO{
				ID:               int(recipient.User.ID),
				FirstNameEN:      recipient.User.FirstNameEn,
				FirstNameAR:      recipient.User.FirstNameAr,
				LastNameEN:       recipient.User.LastNameEn,
				LastNameAR:       recipient.User.LastNameAr,
				Email:            recipient.User.Email,
				ProfileImagePath: recipient.User.ProfileImagePath,
				Initial:          recipient.User.Initial,
				Color:            recipient.User.Color,
			}
		}

		dto.Recipients = append(dto.Recipients, recDto)
	}

	for _, attachment := range email.Attachments {
		dto.Attachments = append(dto.Attachments, logDTOs.EmailLogAttachmentDTO{
			ID:          attachment.ID,
			FileName:    attachment.FileName,
			ContentType: attachment.ContentType,
			SizeBytes:   attachment.SizeBytes,
			StorageKey:  attachment.StorageKey,
			SortOrder:   attachment.SortOrder,
			CreatedAt:   attachment.CreatedAt,
		})
	}

	return &dto, nil
}

func (els *EmailLogService) GetSummary(params *logDTOs.EmailLogListQueryParamsDTO) (*logDTOs.EmailLogSummaryDTO, error) {
	return els.repo.GetSummary(params)
}
