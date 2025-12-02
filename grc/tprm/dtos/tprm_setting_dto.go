package dtos

type TprmSettingDTO struct {
	ID               int  `json:"id"`
	IsExport         bool `json:"is_export"`
	IsExternalVendor bool `json:"is_external_vendor"`
}

type TprmSettingFormDTO struct {
	IsExport         bool `json:"is_export"`
	IsExternalVendor bool `json:"is_external_vendor"`
}
