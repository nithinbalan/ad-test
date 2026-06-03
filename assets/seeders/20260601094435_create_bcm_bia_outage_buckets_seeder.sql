-- +goose Up
-- +goose StatementBegin
DECLARE @hours_id INT;
DECLARE @days_id INT;

SELECT @hours_id = id FROM business_continuity_duration_units WHERE type = 'hours';
SELECT @days_id = id FROM business_continuity_duration_units WHERE type = 'days';

INSERT INTO bcm_bia_outage_buckets (title_en, title_ar, from_value, to_value, duration_unit_id, [order], status) VALUES
    ('0-8 Hours', '0-8 ساعات', 0, 8, @hours_id, 1, 'active'),
    ('8-24 Hours', '8-24 ساعة', 8, 24, @hours_id, 2, 'active'),
    ('1-3 Days', '1-3 أيام', 1, 3, @days_id, 3, 'active'),
    ('3-7 Days', '3-7 أيام', 3, 7, @days_id, 4, 'active'),
    ('7-30 Days', '7-30 يوم', 7, 30, @days_id, 5, 'active');
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DELETE FROM bcm_bia_outage_buckets;
-- +goose StatementEnd
