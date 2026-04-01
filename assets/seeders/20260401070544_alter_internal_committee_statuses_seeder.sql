-- +goose Up
-- +goose StatementBegin

INSERT INTO internal_committee_statuses (title_en, title_ar, color_code, type, status)
VALUES ('Under Nomination', 'قيد الترشيح', '#D2B48C', 'under_nomination', 'active');

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

DELETE FROM internal_committee_statuses
WHERE type = 'under_nomination';

-- +goose StatementEnd