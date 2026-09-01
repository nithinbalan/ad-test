-- +goose Up
-- +goose StatementBegin
INSERT INTO project_security_gate_classifications (
    title_en,
    title_ar,
    type,
    color_code,
    status
)
VALUES
('Cloud', N'السحابة', 'cloud', '#007BFF', 'active'),
('Remote', N'الوصول عن بعد', 'remote', '#D6A2E8', 'active'),
('Data', N'البيانات', 'data', '#28A745', 'active'),
('Social', N'وسائل التواصل الاجتماعي', 'social', '#FFC107', 'active'),
('Third Party', N'طرف ثالث', 'third_party', '#FD7E14', 'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM project_security_gate_classifications;
-- +goose StatementEnd
