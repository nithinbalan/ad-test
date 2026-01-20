-- +goose Up
-- +goose StatementBegin

CREATE TABLE action_plan_master (
    id INT IDENTITY(1,1) PRIMARY KEY,
    module_id INT NOT NULL,
    incident_category_id INT NULL,
    title NVARCHAR(100) NOT NULL,
    description NVARCHAR(MAX) NULL,
    severity_level_id INT NULL,
    status VARCHAR(10) NOT NULL DEFAULT 'active',
      CONSTRAINT CHK_action_plan_master_status CHECK (status IN ('active', 'inactive')),

    CONSTRAINT FK_action_plan_master_modules
        FOREIGN KEY (module_id)
        REFERENCES modules(id) ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_action_plan_master_incident_categories
        FOREIGN KEY (incident_category_id)
        REFERENCES incident_categories(id) ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_action_plan_master_severity_levels
        FOREIGN KEY (severity_level_id)
        REFERENCES severity_levels(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE action_plan_master;

-- +goose StatementEnd
