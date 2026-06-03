package dtos

type BCMBIAOutageBucketDTO struct {
	ID             uint                               `json:"id"`
	TitleEn        string                             `json:"title_en"`
	TitleAr        string                             `json:"title_ar"`
	FromValue      int                                `json:"from_value"`
	ToValue        int                                `json:"to_value"`
	DurationUnitID uint                               `json:"duration_unit_id"`
	DurationUnit   *BusinessContinuityDurationUnitDTO `json:"duration_unit"`
	Order          int                                `json:"order"`
	Status         string                             `json:"status"`
}
