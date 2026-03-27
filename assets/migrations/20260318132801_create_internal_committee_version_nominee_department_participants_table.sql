-- +goose Up
-- +goose StatementBegin
CREATE TABLE internal_committee_version_nominee_department_participants (
    id INT PRIMARY KEY IDENTITY(1,1),
    internal_committee_version_nominee_department_id INT NOT NULL,
    user_id INT NOT NULL,
    department_id INT NOT NULL,
    created_by_id INT,
    updated_by_id INT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME,
    CONSTRAINT FK_internal_committee_version_nominee_department_participants_internal_committee_version_nominee_departments FOREIGN KEY (internal_committee_version_nominee_department_id) REFERENCES internal_committee_version_nominee_departments(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_internal_committee_version_nominee_department_participants_users_1 FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_internal_committee_version_nominee_department_participants_organization_units FOREIGN KEY (department_id) REFERENCES organization_units(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_internal_committee_version_nominee_department_participants_users_2 FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_internal_committee_version_nominee_department_participants_users_3 FOREIGN KEY (updated_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.internal_committee_version_nominee_department_participants;
-- +goose StatementEnd
