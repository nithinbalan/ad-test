package dtos

import "time"

type BCMBIAImpactLevelDTO struct {
	ID                              uint    `json:"id"`
	ImpactAnalysisLevelID           uint    `json:"impact_analysis_level_id"`
	ImpactAnalysisLevelTitleEn      *string `json:"impact_analysis_level_title_en"`
	ImpactAnalysisLevelTitleAr      *string `json:"impact_analysis_level_title_ar"`
	ImpactAnalysisLevelType         *string `json:"impact_analysis_level_type"`
	ImpactAnalysisLevelColorCode    *string `json:"impact_analysis_level_color_code"`
	OutageBucketID                  *uint   `json:"outage_bucket_id"`
	OutageBucketTitleEn             *string `json:"outage_bucket_title_en"`
	OutageBucketTitleAr             *string `json:"outage_bucket_title_ar"`
	OutageBucketFromValue           *int    `json:"outage_bucket_from_value"`
	OutageBucketToValue             *int    `json:"outage_bucket_to_value"`
	OutageBucketDurationUnitID      *uint   `json:"outage_bucket_duration_unit_id"`
	OutageBucketDurationUnitTitleEn *string `json:"outage_bucket_duration_unit_title_en"`
	OutageBucketDurationUnitTitleAr *string `json:"outage_bucket_duration_unit_title_ar"`
	OutageBucketDurationUnitMinutes *int    `json:"outage_bucket_duration_unit_minutes"`
	OutageBucketDurationUnitType    *string `json:"outage_bucket_duration_unit_type"`
	ClassificationEn                string  `json:"classification_en"`
	ClassificationAr                string  `json:"classification_ar"`
	IsCritical                      bool    `json:"is_critical"`
	RtoFromValue                    *int    `json:"rto_from_value"`
	RtoToValue                      *int    `json:"rto_to_value"`
	RtoCondition                    *string `json:"rto_condition"`
	RtoUnitID                       *uint   `json:"rto_unit_id"`
	RtoUnitTitleEn                  *string `json:"rto_unit_title_en"`
	RtoUnitTitleAr                  *string `json:"rto_unit_title_ar"`
	RtoUnitMinutes                  *int    `json:"rto_unit_minutes"`
	RtoUnitType                     *string `json:"rto_unit_type"`
	MtpdFromValue                   *int    `json:"mtpd_from_value"`
	MtpdToValue                     *int    `json:"mtpd_to_value"`
	MtpdCondition                   *string `json:"mtpd_condition"`
	MtpdUnitID                      *uint   `json:"mtpd_unit_id"`
	MtpdUnitTitleEn                 *string `json:"mtpd_unit_title_en"`
	MtpdUnitTitleAr                 *string `json:"mtpd_unit_title_ar"`
	MtpdUnitMinutes                 *int    `json:"mtpd_unit_minutes"`
	MtpdUnitType                    *string `json:"mtpd_unit_type"`
	RpoFromValue                    *int    `json:"rpo_from_value"`
	RpoToValue                      *int    `json:"rpo_to_value"`
	RpoCondition                    *string `json:"rpo_condition"`
	RpoUnitID                       *uint   `json:"rpo_unit_id"`
	RpoUnitTitleEn                  *string `json:"rpo_unit_title_en"`
	RpoUnitTitleAr                  *string `json:"rpo_unit_title_ar"`
	RpoUnitMinutes                  *int    `json:"rpo_unit_minutes"`
	RpoUnitType                     *string `json:"rpo_unit_type"`
	RspFromValue                    *int    `json:"rsp_from_value"`
	RspToValue                      *int    `json:"rsp_to_value"`
	RspCondition                    *string `json:"rsp_condition"`
	RspUnitID                       *uint   `json:"rsp_unit_id"`
	RspUnitTitleEn                  *string `json:"rsp_unit_title_en"`
	RspUnitTitleAr                  *string `json:"rsp_unit_title_ar"`
	RspUnitMinutes                  *int    `json:"rsp_unit_minutes"`
	RspUnitType                     *string `json:"rsp_unit_type"`
	Order                           int     `json:"order"`
	Status                          string  `json:"status"`
}

type BCMBIAImpactLevelFormDTO struct {
	OutageBucketID        uint    `json:"outage_bucket_id" binding:"required"`
	ImpactAnalysisLevelID uint    `json:"impact_analysis_level_id" binding:"required"`
	ClassificationEn      string  `json:"classification_en" binding:"required"`
	ClassificationAr      string  `json:"classification_ar" binding:"required"`
	IsCritical            bool    `json:"is_critical"`
	RtoFromValue          *int    `json:"rto_from_value"`
	RtoToValue            *int    `json:"rto_to_value"`
	RtoCondition          *string `json:"rto_condition"`
	RtoUnitID             *uint   `json:"rto_unit_id"`
	MtpdFromValue         *int    `json:"mtpd_from_value"`
	MtpdToValue           *int    `json:"mtpd_to_value"`
	MtpdCondition         *string `json:"mtpd_condition"`
	MtpdUnitID            *uint   `json:"mtpd_unit_id"`
	RpoFromValue          *int    `json:"rpo_from_value"`
	RpoToValue            *int    `json:"rpo_to_value"`
	RpoCondition          *string `json:"rpo_condition"`
	RpoUnitID             *uint   `json:"rpo_unit_id"`
	RspFromValue          *int    `json:"rsp_from_value"`
	RspToValue            *int    `json:"rsp_to_value"`
	RspCondition          *string `json:"rsp_condition"`
	RspUnitID             *uint   `json:"rsp_unit_id"`
}

type BCMBIAImpactLevelBulkUpdateDTO struct {
	Settings []*BCMBIAImpactLevelFormDTO `json:"settings" binding:"required,dive"`
}

type BCMBIAImpactLevelDetailDTO struct {
	ID                  uint                               `json:"id"`
	OutageBucketID      *uint                              `json:"outage_bucket_id"`
	OutageBucket        *BCMBIAOutageBucketDTO             `json:"outage_bucket"`
	ClassificationEn    string                             `json:"classification_en"`
	ClassificationAr    string                             `json:"classification_ar"`
	IsCritical          bool                               `json:"is_critical"`
	RtoFromValue        *int                               `json:"rto_from_value"`
	RtoToValue          *int                               `json:"rto_to_value"`
	RtoCondition        *string                            `json:"rto_condition"`
	MtpdFromValue       *int                               `json:"mtpd_from_value"`
	MtpdToValue         *int                               `json:"mtpd_to_value"`
	MtpdCondition       *string                            `json:"mtpd_condition"`
	RpoFromValue        *int                               `json:"rpo_from_value"`
	RpoToValue          *int                               `json:"rpo_to_value"`
	RpoCondition        *string                            `json:"rpo_condition"`
	RspFromValue        *int                               `json:"rsp_from_value"`
	RspToValue          *int                               `json:"rsp_to_value"`
	RspCondition        *string                            `json:"rsp_condition"`
	Order               int                                `json:"order"`
	Status              string                             `json:"status"`
	ImpactAnalysisLevel *BCMImpactAnalysisLevelDTO         `json:"impact_analysis_level"`
	RtoUnit             *BusinessContinuityDurationUnitDTO `json:"rto_unit"`
	MtpdUnit            *BusinessContinuityDurationUnitDTO `json:"mtpd_unit"`
	RpoUnit             *BusinessContinuityDurationUnitDTO `json:"rpo_unit"`
	RspUnit             *BusinessContinuityDurationUnitDTO `json:"rsp_unit"`
	CreatedAt           time.Time                          `json:"created_at"`
	UpdatedAt           *time.Time                         `json:"updated_at"`
}
