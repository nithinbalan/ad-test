package dtos

type PaginationInfoDTO struct {
	CurrentPage int   `json:"currentPage"`
	PageSize    int   `json:"pageSize"`
	TotalPages  int   `json:"totalPages"`
	TotalCount  int64 `json:"totalCount"`
}

type IDResponse struct {
	ID uint `json:"id"`
}

type FileDTO struct {
	FileName     string `json:"file_name"`
	OriginalName string `json:"original_name"`
	FileType     string `json:"type"`
	SizeInBytes  int64  `json:"size"`
}
