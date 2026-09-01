-- +goose Up
-- +goose StatementBegin
INSERT INTO project_security_gate_stages (
    title_en,
    title_ar,
    type,
    color_code,
    status
)
VALUES
('Pre-implementation', N'قبل التنفيذ', 'pre_implementation', '#007BFF', 'active'),
('During implementation', N'أثناء التنفيذ', 'during_implementation', '#FFC107', 'active'),
('Post-implementation', N'بعد التنفيذ', 'post_implementation', '#28A745', 'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM project_security_gate_stages;
-- +goose StatementEnd
