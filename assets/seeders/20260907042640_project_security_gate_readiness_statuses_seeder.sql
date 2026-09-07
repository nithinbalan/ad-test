-- +goose Up
-- +goose StatementBegin
INSERT INTO project_security_gate_readiness_statuses (
    title_en,
    title_ar,
    type,
    color_code,
    status
)
VALUES
('Draft', N'مسودة', 'draft', '#6C757D', 'active'),
('Go Live', N'الإطلاق', 'go_live', '#28A745', 'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM project_security_gate_readiness_statuses;
-- +goose StatementEnd