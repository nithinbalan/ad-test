-- +goose Up
-- +goose StatementBegin
DECLARE @non_applicable_id INT;
DECLARE @very_low_id INT;
DECLARE @low_id INT;
DECLARE @moderate_id INT;
DECLARE @high_id INT;
DECLARE @very_high_id INT;

DECLARE @bucket_0_8_hours_id INT;
DECLARE @bucket_8_24_hours_id INT;
DECLARE @bucket_1_3_days_id INT;
DECLARE @bucket_3_7_days_id INT;
DECLARE @bucket_7_30_days_id INT;

SELECT @non_applicable_id = id FROM bcm_impact_analysis_levels WHERE type = 'non_applicable';
SELECT @very_low_id = id FROM bcm_impact_analysis_levels WHERE type = 'very_low';
SELECT @low_id = id FROM bcm_impact_analysis_levels WHERE type = 'low';
SELECT @moderate_id = id FROM bcm_impact_analysis_levels WHERE type = 'moderate';
SELECT @high_id = id FROM bcm_impact_analysis_levels WHERE type = 'high';
SELECT @very_high_id = id FROM bcm_impact_analysis_levels WHERE type = 'very_high';

SELECT @bucket_0_8_hours_id = id FROM bcm_bia_outage_buckets WHERE title_en = '0-8 Hours';
SELECT @bucket_8_24_hours_id = id FROM bcm_bia_outage_buckets WHERE title_en = '8-24 Hours';
SELECT @bucket_1_3_days_id = id FROM bcm_bia_outage_buckets WHERE title_en = '1-3 Days';
SELECT @bucket_3_7_days_id = id FROM bcm_bia_outage_buckets WHERE title_en = '3-7 Days';
SELECT @bucket_7_30_days_id = id FROM bcm_bia_outage_buckets WHERE title_en = '7-30 Days';

UPDATE bcm_bia_impact_levels
SET classification_en = 'Level I', classification_ar = N'المستوى الأول', [order] = 1, outage_bucket_id = @bucket_0_8_hours_id
WHERE impact_analysis_level_id = @very_high_id;

UPDATE bcm_bia_impact_levels
SET classification_en = 'Level II', classification_ar = N'المستوى الثاني', [order] = 2, outage_bucket_id = @bucket_8_24_hours_id
WHERE impact_analysis_level_id = @high_id;

UPDATE bcm_bia_impact_levels
SET classification_en = 'Level III', classification_ar = N'المستوى الثالث', [order] = 3, outage_bucket_id = @bucket_1_3_days_id
WHERE impact_analysis_level_id = @moderate_id;

UPDATE bcm_bia_impact_levels
SET classification_en = 'Level IV', classification_ar = N'المستوى الرابع', [order] = 4, outage_bucket_id = @bucket_3_7_days_id
WHERE impact_analysis_level_id = @low_id;

UPDATE bcm_bia_impact_levels
SET classification_en = 'Level V', classification_ar = N'المستوى الخامس', [order] = 5, outage_bucket_id = @bucket_7_30_days_id
WHERE impact_analysis_level_id = @very_low_id;

UPDATE bcm_bia_impact_levels
SET classification_en = 'Level VI', classification_ar = N'المستوى السادس', [order] = 6, outage_bucket_id = NULL
WHERE impact_analysis_level_id = @non_applicable_id;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DECLARE @non_applicable_id INT;
DECLARE @very_low_id INT;
DECLARE @low_id INT;
DECLARE @moderate_id INT;
DECLARE @high_id INT;
DECLARE @very_high_id INT;

SELECT @non_applicable_id = id FROM bcm_impact_analysis_levels WHERE type = 'non_applicable';
SELECT @very_low_id = id FROM bcm_impact_analysis_levels WHERE type = 'very_low';
SELECT @low_id = id FROM bcm_impact_analysis_levels WHERE type = 'low';
SELECT @moderate_id = id FROM bcm_impact_analysis_levels WHERE type = 'moderate';
SELECT @high_id = id FROM bcm_impact_analysis_levels WHERE type = 'high';
SELECT @very_high_id = id FROM bcm_impact_analysis_levels WHERE type = 'very_high';

UPDATE bcm_bia_impact_levels
SET classification_en = 'Level I', classification_ar = N'المستوى الأول', [order] = 1, outage_bucket_id = NULL
WHERE impact_analysis_level_id = @non_applicable_id;

UPDATE bcm_bia_impact_levels
SET classification_en = 'Level II', classification_ar = N'المستوى الثاني', [order] = 2, outage_bucket_id = NULL
WHERE impact_analysis_level_id = @very_low_id;

UPDATE bcm_bia_impact_levels
SET classification_en = 'Level III', classification_ar = N'المستوى الثالث', [order] = 3, outage_bucket_id = NULL
WHERE impact_analysis_level_id = @low_id;

UPDATE bcm_bia_impact_levels
SET classification_en = 'Level IV', classification_ar = N'المستوى الرابع', [order] = 4, outage_bucket_id = NULL
WHERE impact_analysis_level_id = @moderate_id;

UPDATE bcm_bia_impact_levels
SET classification_en = 'Level V', classification_ar = N'المستوى الخامس', [order] = 5, outage_bucket_id = NULL
WHERE impact_analysis_level_id = @high_id;

UPDATE bcm_bia_impact_levels
SET classification_en = 'Level VI', classification_ar = N'المستوى السادس', [order] = 6, outage_bucket_id = NULL
WHERE impact_analysis_level_id = @very_high_id;
-- +goose StatementEnd
