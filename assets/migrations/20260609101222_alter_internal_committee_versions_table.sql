-- +goose Up
-- +goose StatementBegin

ALTER TABLE internal_committee_versions
ADD canceled_by_id INT NULL,
    canceled_at DATETIME NULL,
    disband_by_id INT NULL,
    disband_at DATETIME NULL,
    completed_by_id INT NULL,
    completed_at DATETIME NULL;

ALTER TABLE internal_committee_versions
ADD CONSTRAINT FK_internal_committee_versions_canceled_by
FOREIGN KEY (canceled_by_id) REFERENCES users(id);

ALTER TABLE internal_committee_versions
ADD CONSTRAINT FK_internal_committee_versions_disband_by
FOREIGN KEY (disband_by_id) REFERENCES users(id);

ALTER TABLE internal_committee_versions
ADD CONSTRAINT FK_internal_committee_versions_completed_by
FOREIGN KEY (completed_by_id) REFERENCES users(id);

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

ALTER TABLE internal_committee_versions
DROP CONSTRAINT FK_internal_committee_versions_canceled_by;

ALTER TABLE internal_committee_versions
DROP CONSTRAINT FK_internal_committee_versions_disband_by;

ALTER TABLE internal_committee_versions
DROP CONSTRAINT FK_internal_committee_versions_completed_by;

ALTER TABLE internal_committee_versions
DROP COLUMN canceled_by_id,
            canceled_at,
            disband_by_id,
            disband_at,
            completed_by_id,
            completed_at;

-- +goose StatementEnd