-- +goose Up
-- +goose StatementBegin
CREATE TABLE internal_committee_versions
(
    id INT PRIMARY KEY IDENTITY(1, 1),
    internal_committee_id INT NOT NULL,
    committee_name_en NVARCHAR(MAX) NOT NULL,
    committee_name_ar NVARCHAR(MAX) NOT NULL,
    committee_type_id INT NOT NULL,
    committee_term_id INT NOT NULL,
    term_end_date DATE NULL,
    committee_classification_id INT NULL,
    justification NVARCHAR(MAX) NULL,
    roles_and_responsibilities NVARCHAR(MAX) NULL,
    committee_status_id INT NOT NULL,
    version DECIMAL(4, 2) NULL,
    is_latest_version BIT NULL,
    created_by_id INT NOT NULL,
    updated_by_id INT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,
    submitted_by_id INT NULL,
    approved_by_id INT NULL,
    submitted_at DATETIME NULL,
    approved_at DATETIME NULL,
    next_review_level INT NULL,
    workflow_status_id INT NULL,
    CONSTRAINT FK_internal_committee_versions_internal_committees FOREIGN KEY (internal_committee_id) REFERENCES internal_committees (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_internal_committee_versions_committee_types FOREIGN KEY (committee_type_id) REFERENCES committee_types (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_internal_committee_versions_committee_terms FOREIGN KEY (committee_term_id) REFERENCES committee_terms (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_internal_committee_versions_committee_classifications FOREIGN KEY (committee_classification_id) REFERENCES committee_classifications (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_internal_committee_versions_internal_committee_statuses FOREIGN KEY (committee_status_id) REFERENCES internal_committee_statuses (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_internal_committee_versions_users_1 FOREIGN KEY (created_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_internal_committee_versions_users_2 FOREIGN KEY (updated_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_internal_committee_versions_users_3 FOREIGN KEY (submitted_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_internal_committee_versions_users_4 FOREIGN KEY (approved_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_internal_committee_versions_workflow_statuses FOREIGN KEY (workflow_status_id) REFERENCES workflow_statuses (id) ON DELETE NO ACTION ON UPDATE NO ACTION
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE internal_committee_versions;
-- +goose StatementEnd