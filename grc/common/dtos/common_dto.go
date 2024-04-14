package dtos

type PaginationInfoDTO struct {
	CurrentPage int   `json:"current_page"`
	PageSize    int   `json:"page_size"`
	TotalPages  int   `json:"total_pages"`
	TotalCount  int64 `json:"total_count"`
}

type IDResponse struct {
	ID uint `json:"id"`
}

type SuccessfulIDsResponse struct {
	SuccessfulIds []uint `json:"successfull_ids"`
}

type FailedIDsResponse struct {
	FailedIds []uint `json:"failed_ids"`
}

type FileDTO struct {
	FileName     string `json:"file_name"`
	OriginalName string `json:"original_name"`
	FileType     string `json:"type"`
	SizeInBytes  int64  `json:"size"`
}
