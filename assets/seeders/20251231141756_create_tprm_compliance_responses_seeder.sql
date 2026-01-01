-- +goose Up
-- +goose StatementBegin

INSERT INTO tprm_compliance_responses (title_en, title_ar, score, status,color_code,created_by_id,created_at) VALUES
('Not Applicable', 'غير قابل للتطبيق', 0, 'active','#6C757D',1,GETDATE()),
('Comply', 'امتثال', 3, 'active','#28A745',1,GETDATE()),
('Partially Comply', 'الامتثال الجزئي', 2, 'active','#FFC107',1,GETDATE()),
('Not Comply', 'لا يمتثل', 1, 'active','#DC3545',1,GETDATE());

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE tprm_compliance_responses;
-- +goose StatementEnd
