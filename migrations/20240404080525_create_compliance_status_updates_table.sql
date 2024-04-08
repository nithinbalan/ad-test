-- +goose Up
-- +goose StatementBegin
CREATE TABLE compliance_status_updates (
    id INT PRIMARY KEY IDENTITY(1, 1),
    compliance_register_id INT NOT NULL,
    compliance_status_id INT NOT NULL,
    comment NVARCHAR(MAX) NOT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,
    created_by INT NOT NULL,
    updated_by INT NULL,
    CONSTRAINT FK_compliance_status_updates_compliance_register FOREIGN KEY (compliance_register_id) REFERENCES compliance_register (id),
    CONSTRAINT FK_compliance_status_updates_compliance_statuses FOREIGN KEY (compliance_status_id) REFERENCES compliance_statuses (id),
    CONSTRAINT FK_compliance_status_updates_users1 FOREIGN KEY (created_by) REFERENCES users (id),
    CONSTRAINT FK_compliance_status_updates_users2 FOREIGN KEY (updated_by) REFERENCES users (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE compliance_status_updates;
-- +goose StatementEnd
