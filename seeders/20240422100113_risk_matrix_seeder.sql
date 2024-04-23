-- +goose Up
-- +goose StatementBegin
INSERT INTO dbo.risk_matrix (title, [calculation_type], [type]) VALUES
('ERM', 'multiplication', 'erm'),
('ISMS', 'multiplication', 'isms');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE risk_matrix;
-- +goose StatementEnd
