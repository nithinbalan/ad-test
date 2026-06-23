-- +goose Up
-- +goose StatementBegin
INSERT INTO system_register_disruption_levels (title_en, title_ar, created_by_id, created_at, status) VALUES
    ('Moderate', 'معتدل', 1, GETDATE(), 'active'),
    ('Heavy', 'شديد', 1, GETDATE(), 'active'),
    ('Severe', 'حاد', 1, GETDATE(), 'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM system_register_disruption_levels;
-- +goose StatementEnd
