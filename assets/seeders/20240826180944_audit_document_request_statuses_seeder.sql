-- +goose Up
-- +goose StatementBegin
INSERT INTO audit_document_request_statuses
    (title_en, title_ar, type, color_code, status)
VALUES
    ('Open', 'مفتوح', 'open', '#28A745', 'active'),
    ('Draft', 'مسودة', 'draft', '#FFC107', 'active'),
    ('Closed', 'مغلق', 'closed', '#DC3545', 'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE audit_document_request_statuses;
-- +goose StatementEnd
