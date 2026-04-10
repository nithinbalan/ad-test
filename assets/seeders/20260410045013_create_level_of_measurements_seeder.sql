-- +goose Up
-- +goose StatementBegin
INSERT INTO level_of_measurements (title_en, title_ar, type, status) VALUES
('Board Level', N'مستوى مجلس الإدارة', 'board_of_directors', 'active'),
('Committee Level', N'مستوى اللجنة', 'grc_committee', 'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM level_of_measurements;
-- +goose StatementEnd