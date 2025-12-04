-- +goose Up
-- +goose StatementBegin
INSERT INTO third_party_assessment_statuses (title_en, title_ar, color_code, type, status)
VALUES 
    ('Approved', 'موافقة', '#28A745', 'approved', 'active'),
    ('Rejected', 'مرفوض', '#DC3545', 'rejected', 'active');

UPDATE third_party_assessment_statuses
SET color_code = '#007BFF' WHERE [type] = 'submitted';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM third_party_assessment_statuses
WHERE [type] IN ('approved', 'rejected');

UPDATE third_party_assessment_statuses
SET color_code = '#28A745' WHERE [type] = 'submitted';
-- +goose StatementEnd
