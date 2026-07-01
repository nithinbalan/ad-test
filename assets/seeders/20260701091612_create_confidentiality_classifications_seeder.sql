-- +goose Up
-- +goose StatementBegin
INSERT INTO confidentiality_classifications (
    title_en,
    title_ar,
    type
)
VALUES
    ('Public', N'عام', 'public'),
    ('Internal', N'داخلي', 'internal'),
    ('Confidential', N'سري', 'confidential'),
    ('Restricted', N'مقيد', 'restricted');
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DELETE FROM confidentiality_classifications;
-- +goose StatementEnd