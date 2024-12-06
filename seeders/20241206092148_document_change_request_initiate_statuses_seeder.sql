-- +goose Up
-- +goose StatementBegin
INSERT INTO document_change_request_initiate_statuses (title_en, title_ar, color_code, type, status) VALUES
('Initiated', 'مبادر', '#17A2B8', 'initiated', 'active'),
('Approved', 'مقبول', '#28A745', 'approved', 'active'),
('Rejected', 'مرفوض', '#E0554D', 'rejected', 'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM dbo.document_change_request_initiate_statuses WHERE type IN ('initiated', 'approved', 'rejected');
-- +goose StatementEnd
