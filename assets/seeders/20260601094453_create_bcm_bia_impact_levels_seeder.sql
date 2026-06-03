-- +goose Up
-- +goose StatementBegin
DECLARE @hours_id INT;
DECLARE @days_id INT;
DECLARE @weeks_id INT;
DECLARE @months_id INT;

DECLARE @non_applicable_id INT;
DECLARE @very_low_id INT;
DECLARE @low_id INT;
DECLARE @moderate_id INT;
DECLARE @high_id INT;
DECLARE @very_high_id INT;

SELECT @hours_id = id FROM business_continuity_duration_units WHERE type = 'hours';
SELECT @days_id = id FROM business_continuity_duration_units WHERE type = 'days';
SELECT @weeks_id = id FROM business_continuity_duration_units WHERE type = 'weeks';
SELECT @months_id = id FROM business_continuity_duration_units WHERE type = 'months';

SELECT @non_applicable_id = id FROM bcm_impact_analysis_levels WHERE type = 'non_applicable';
SELECT @very_low_id = id FROM bcm_impact_analysis_levels WHERE type = 'very_low';
SELECT @low_id = id FROM bcm_impact_analysis_levels WHERE type = 'low';
SELECT @moderate_id = id FROM bcm_impact_analysis_levels WHERE type = 'moderate';
SELECT @high_id = id FROM bcm_impact_analysis_levels WHERE type = 'high';
SELECT @very_high_id = id FROM bcm_impact_analysis_levels WHERE type = 'very_high';

INSERT INTO bcm_bia_impact_levels
(
    impact_analysis_level_id,
    classification_en,
    classification_ar,
    is_critical,
    rto_from_value,
    rto_to_value,
    rto_condition,
    rto_unit_id,
    mtpd_from_value,
    mtpd_to_value,
    mtpd_condition,
    mtpd_unit_id,
    rpo_from_value,
    rpo_to_value,
    rpo_condition,
    rpo_unit_id,
    rsp_from_value,
    rsp_to_value,
    rsp_condition,
    rsp_unit_id,
    [order],
    status
)
VALUES
(
    @non_applicable_id,
    'Level I',
    N'المستوى الأول',
    0,
    1, NULL, 'exact', @months_id,
    1, NULL, 'greater_than', @months_id,
    1, NULL, 'less_than', @months_id,
    7, NULL, 'greater_than', @days_id,
    1,
    'active'
),
(
    @very_low_id,
    'Level II',
    N'المستوى الثاني',
    1,
    7, NULL, 'exact', @days_id,
    14, NULL, 'exact', @days_id,
    7, NULL, 'less_than', @days_id,
    2, 3, 'range', @days_id,
    2,
    'active'
),
(
    @low_id,
    'Level III',
    N'المستوى الثالث',
    1,
    3, NULL, 'exact', @days_id,
    7, NULL, 'exact', @days_id,
    3, NULL, 'less_than', @days_id,
    1, NULL, 'exact', @days_id,
    3,
    'active'
),
(
    @moderate_id,
    'Level IV',
    N'المستوى الرابع',
    1,
    24, NULL, 'exact', @hours_id,
    3, NULL, 'exact', @days_id,
    24, NULL, 'less_than', @hours_id,
    4, NULL, 'exact', @hours_id,
    4,
    'active'
),
(
    @high_id,
    'Level V',
    N'المستوى الخامس',
    1,
    8, NULL, 'exact', @hours_id,
    24, NULL, 'exact', @hours_id,
    8, NULL, 'less_than', @hours_id,
    2, NULL, 'exact', @hours_id,
    5,
    'active'
),
(
    @very_high_id,
    'Level VI',
    N'المستوى السادس',
    1,
    4, NULL, 'exact', @hours_id,
    8, NULL, 'exact', @hours_id,
    4, NULL, 'less_than', @hours_id,
    1, NULL, 'exact', @hours_id,
    6,
    'active'
);
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DELETE FROM bcm_bia_impact_levels;
-- +goose StatementEnd
