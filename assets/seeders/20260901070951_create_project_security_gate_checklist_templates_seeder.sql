-- +goose Up
-- +goose StatementBegin
INSERT INTO project_security_gate_checklist_templates (
    stage_id,
    title_en,
    title_ar
)
VALUES
(
    (
        SELECT id
        FROM project_security_gate_stages
        WHERE type = 'pre_implementation'
    ),
    N'Pre-Implementation Security Checklist',
    N'قائمة التحقق الأمنية قبل التنفيذ'
),
(
    (
        SELECT id
        FROM project_security_gate_stages
        WHERE type = 'during_implementation'
    ),
    N'During Implementation Security Checklist',
    N'قائمة التحقق الأمنية أثناء التنفيذ'
),
(
    (
        SELECT id
        FROM project_security_gate_stages
        WHERE type = 'post_implementation'
    ),
    N'Post-Implementation Security Checklist',
    N'قائمة التحقق الأمنية بعد التنفيذ'
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM project_security_gate_checklist_templates;
-- +goose StatementEnd
