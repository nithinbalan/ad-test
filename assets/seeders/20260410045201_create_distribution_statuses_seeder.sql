-- +goose Up
-- +goose StatementBegin
INSERT INTO distribution_statuses (title_en, title_ar, type, status) VALUES
('Cumulative Distribution', N'توزيع تراكمي', 'cumulative', 'active'),
('Individual Distribution', N'توزيع فردي', 'individual', 'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM distribution_statuses;
-- +goose StatementEnd