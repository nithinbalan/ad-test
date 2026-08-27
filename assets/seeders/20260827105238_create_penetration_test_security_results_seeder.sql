-- +goose Up
-- +goose StatementBegin
INSERT INTO penetration_test_security_results (title_en, title_ar, type, status)
VALUES ('No Findings', N'لم يتم العثور على نتائج', 'no-findings', 'active'), ('Findings Identified', N'تم تحديد نتائج', 'findings-identified', 'active'), ('Remediation In Progress', N'المعالجة قيد التنفيذ', 'remediation-in-progress', 'active'), ('Remediation Completed', N'اكتملت المعالجة', 'remediation-completed', 'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE penetration_test_security_results;
-- +goose StatementEnd
