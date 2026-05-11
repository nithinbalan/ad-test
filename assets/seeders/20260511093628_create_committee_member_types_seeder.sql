-- +goose Up
-- +goose StatementBegin
INSERT INTO committee_member_types (title_en, title_ar, status) VALUES
    ('Committee Chair', 'رئيس اللجنة', 'active'),
    ('Committee Vice-Chair', 'نائب رئيس اللجنة', 'active'),
    ('Member', 'عضو', 'active'),
    ('Alternate Member', 'عضو احتياطي', 'active'),
    ('Committee Secretary', 'أمين اللجنة', 'active');
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DELETE FROM committee_member_types
WHERE title_en IN (
    'Committee Chair',
    'Committee Vice-Chair',
    'Member',
    'Alternate Member',
    'Committee Secretary'
);
-- +goose StatementEnd