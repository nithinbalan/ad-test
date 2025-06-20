-- +goose Up
-- +goose StatementBegin
CREATE TABLE compliance_requirement_status_updates (
    id INT PRIMARY KEY IDENTITY(1, 1),
    compliance_requirement_id INT NOT NULL,
    compliance_status_id INT NOT NULL,
    comment NVARCHAR(MAX) NOT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,
    created_by_id INT NOT NULL,
    updated_by_id INT NULL,
    CONSTRAINT FK_compliance_requirement_status_updates_compliance_requirements FOREIGN KEY (compliance_requirement_id) REFERENCES compliance_requirements (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_compliance_requirement_status_updates_compliance_statuses FOREIGN KEY (compliance_status_id) REFERENCES compliance_statuses (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_compliance_requirement_status_updates_users_1 FOREIGN KEY (created_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_compliance_requirement_status_updates_users_2 FOREIGN KEY (updated_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE compliance_requirement_status_updates;
-- +goose StatementEnd
