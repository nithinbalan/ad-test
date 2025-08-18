-- +goose Up
-- +goose StatementBegin
INSERT INTO compliance_assessment_statuses 
(title_en, title_ar, color_code, type, status) VALUES
('In progress', 'في تَقَدم', '#FFC107', 'in_progress', 'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM compliance_assessment_statuses 
WHERE [type] = 'in_progress';
-- +goose StatementEnd
