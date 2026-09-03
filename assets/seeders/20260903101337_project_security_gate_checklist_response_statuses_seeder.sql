-- +goose Up
-- +goose StatementBegin
INSERT INTO project_security_gate_checklist_response_statuses (
    title_en,
    title_ar,
    type,
    color_code,
    status
)
VALUES
('Implemented',N'تم التنفيذ','implemented','#28A745','active'),
('Partially Implemented',N'تم التنفيذ جزئياً','partially_implemented','#FFC107','active'),
('Not Implemented',N'لم يتم التنفيذ','not_implemented','#DC3545','active'),
('Not Applicable',N'لا ينطبق','not_applicable','#6C757D','active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM project_security_gate_checklist_response_statuses;
-- +goose StatementEnd