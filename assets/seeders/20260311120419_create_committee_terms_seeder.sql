-- +goose Up
-- +goose StatementBegin
INSERT INTO committee_terms (title_en, title_ar, type, status) VALUES
    ('Temporary', 'مؤقت', 'temporary', 'active'),
    ('Permanent', 'دائم', 'permanent', 'active');
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DELETE FROM committee_terms
WHERE type IN (
    'temporary',
    'permanent'
);
-- +goose StatementEnd
