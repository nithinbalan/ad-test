-- +goose Up
-- +goose StatementBegin
INSERT INTO audit_draft_finding_statuses
    (title_en, title_ar, color_code, type, is_finalized)
VALUES
    ('Open', 'مفتوح', '#007BFF', 'open', 0),-- blue
    ('Accepted', 'مقبول', '#28A745', 'accepted', 1),-- Green
    ('Challenged', 'تم الطعن عليه', '#FFC107', 'challenged', 0), --yellow
    ('Dismissed', 'رفض', '#DC3545', 'dismissed', 1);-- red

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE audit_draft_finding_statuses;
-- +goose StatementEnd
