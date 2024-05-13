-- +goose Up
-- +goose StatementBegin
CREATE TABLE incident_action_plan_status_update_documents (
    id INT PRIMARY KEY IDENTITY(1,1),
    incident_action_plan_status_update_id INT NOT NULL,
    url VARCHAR(1000) DEFAULT NULL,
    name NVARCHAR(255) NULL,
    size INT DEFAULT NULL,
    token VARCHAR(255) DEFAULT NULL,
    type VARCHAR(255) DEFAULT NULL,
    icon VARCHAR(10) NOT NULL, -- Assuming icon names like 'pdf', 'excel', etc.
        CONSTRAINT CHK_incident_action_plan_status_update_documents_icon CHECK (icon IN ('pdf', 'excel', 'ppt', 'image', 'text', 'other')),
    created_by INT,
    updated_by INT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,
    CONSTRAINT FK_incident_action_plan_status_update_documents FOREIGN KEY (incident_action_plan_status_update_id) REFERENCES incident_action_plan_status_updates(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_incident_action_plan_status_update_documents_users_1 FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_incident_action_plan_status_update_documents_users_2 FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.incident_action_plan_status_update_documents;
-- +goose StatementEnd
