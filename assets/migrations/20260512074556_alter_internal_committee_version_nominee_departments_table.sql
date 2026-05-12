-- +goose Up
-- +goose StatementBegin

-- Add new column
ALTER TABLE internal_committee_version_nominee_departments
ADD committee_member_type_id INT NULL;

-- Add foreign key constraint
ALTER TABLE internal_committee_version_nominee_departments
ADD CONSTRAINT FK_internal_committee_version_nominee_departments_committee_member_types
FOREIGN KEY (committee_member_type_id)
REFERENCES committee_member_types(id)
ON DELETE NO ACTION
ON UPDATE NO ACTION;

-- If old membership_type data is no longer needed
ALTER TABLE internal_committee_version_nominee_departments
DROP COLUMN membership_type;

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

-- Recreate old column
ALTER TABLE internal_committee_version_nominee_departments
ADD membership_type NVARCHAR(MAX) NULL;

-- Drop foreign key constraint
ALTER TABLE internal_committee_version_nominee_departments
DROP CONSTRAINT FK_internal_committee_version_nominee_departments_committee_member_types;

-- Drop new column
ALTER TABLE internal_committee_version_nominee_departments
DROP COLUMN committee_member_type_id;

-- +goose StatementEnd