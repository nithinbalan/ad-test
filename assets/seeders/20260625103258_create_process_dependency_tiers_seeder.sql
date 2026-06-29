-- +goose Up
-- +goose StatementBegin
INSERT INTO process_dependency_tiers (title_en, title_ar, created_by_id, created_at, status) VALUES
    ('Low', 'منخفض', 1, GETDATE(), 'active'),
    ('Medium', 'متوسط', 1, GETDATE(), 'active'),
    ('High', 'عالي', 1, GETDATE(), 'active'),
    ('Critical', 'حرج', 1, GETDATE(), 'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM process_dependency_tiers;
-- +goose StatementEnd
