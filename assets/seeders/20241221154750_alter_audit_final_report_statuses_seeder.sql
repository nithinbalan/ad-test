-- +goose Up
-- +goose StatementBegin
UPDATE audit_final_report_statuses
SET 
    title_en = 'Reverted',
    title_ar = 'تم إرجاعه'
WHERE 
    type = 'revert';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
UPDATE audit_final_report_statuses
SET 
    title_en = 'Revert',
    title_ar = 'إرجاع'
WHERE 
    type = 'revert';
-- +goose StatementEnd
