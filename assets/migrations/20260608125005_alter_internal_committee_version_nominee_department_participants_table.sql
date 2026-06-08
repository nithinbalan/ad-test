-- +goose Up
-- +goose StatementBegin

ALTER TABLE internal_committee_version_nominee_department_participants
ALTER COLUMN user_id INT NULL;

ALTER TABLE internal_committee_version_nominee_department_participants
ADD external_name VARCHAR(255) NULL;

ALTER TABLE internal_committee_version_nominee_department_participants
ADD external_email VARCHAR(255) NULL;

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

ALTER TABLE internal_committee_version_nominee_department_participants
DROP COLUMN external_email;

ALTER TABLE internal_committee_version_nominee_department_participants
DROP COLUMN external_name;

ALTER TABLE internal_committee_version_nominee_department_participants
ALTER COLUMN user_id INT NULL;

-- +goose StatementEnd