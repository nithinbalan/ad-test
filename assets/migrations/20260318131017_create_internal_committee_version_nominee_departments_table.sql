-- +goose Up
-- +goose StatementBegin
CREATE TABLE internal_committee_version_nominee_departments (
    id INT PRIMARY KEY IDENTITY(1,1),
    internal_committee_version_id INT NOT NULL,
    nominee_department_id INT NOT NULL,
    membership_type NVARCHAR(MAX) NULL,
    created_by_id INT,
    updated_by_id INT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME,
    CONSTRAINT FK_internal_committee_version_nominee_departments_internal_committee_versions FOREIGN KEY (internal_committee_version_id) REFERENCES internal_committee_versions(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_internal_committee_version_nominee_departments_users_1 FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_internal_committee_version_nominee_departments_users_2 FOREIGN KEY (updated_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.internal_committee_version_nominee_departments;
-- +goose StatementEnd
