-- +goose Up
-- +goose StatementBegin
INSERT INTO committee_types (title_en, title_ar, type, status) VALUES
    ('Statutory Reference', 'مرجع قانوني', 'statutory_reference', 'active'),
    ('No Statutory Reference', 'لا يوجد مرجع قانوني', 'no_statutory_reference', 'active');
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DELETE FROM committee_types
WHERE type IN (
    'statutory_reference',
    'no_statutory_reference'
);
-- +goose StatementEnd