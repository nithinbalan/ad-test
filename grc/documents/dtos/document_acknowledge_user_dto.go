package dtos

type DocumentAcknowledgeListQueryParamsDTO struct {
	Sort     string
	SortBy   string
	Keyword  string
	Page     int
	PageSize int
	Filters  map[string][]string
}

type DocumentAcknowledgementListDTO struct {
	ID           uint   `json:"user_id"`
	FirstNameEn  string `json:"first_name_en"`
	FirstNameAr  string `json:"first_name_ar"`
	Email        string `json:"email"`
	Acknowledged bool   `json:"acknowledged"`
}

type DocumentAcknowledgementListResponseDTO struct {
	Acknowledgements []DocumentAcknowledgementListDTO `json:"acknowledgements"`
}
