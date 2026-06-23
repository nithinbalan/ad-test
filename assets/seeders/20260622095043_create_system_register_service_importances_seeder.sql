-- +goose Up
-- +goose StatementBegin
DELETE FROM  system_register_service_importances;
-- +goose StatementEnd

-- +goose StatementBegin
DBCC CHECKIDENT ('system_register_service_importances', RESEED, 1);
-- +goose StatementEnd

-- +goose StatementBegin
INSERT INTO system_register_service_importances (title_en, title_ar, created_by_id, created_at, status) VALUES
    ('Critical', 'حرج', 1, GETDATE(), 'active'),
    ('Important', 'مهم', 1, GETDATE(), 'active'),
    ('Medium', 'متوسط', 1, GETDATE(), 'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM system_register_service_importances;
-- +goose StatementEnd
