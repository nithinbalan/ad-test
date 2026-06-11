-- +goose Up
-- +goose StatementBegin

CREATE TABLE governance_settings (
    id INT PRIMARY KEY IDENTITY(1,1),

    approve_department_id INT NOT NULL,

    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,

    updated_by_id INT NULL,
    updated_at DATETIME NULL,

    CONSTRAINT FK_governance_settings_organization_units
        FOREIGN KEY (approve_department_id)
        REFERENCES organization_units(id)
        ON DELETE NO ACTION
        ON UPDATE NO ACTION,

    CONSTRAINT FK_governance_settings_created_by
        FOREIGN KEY (created_by_id)
        REFERENCES users(id)
        ON DELETE NO ACTION
        ON UPDATE NO ACTION,

    CONSTRAINT FK_governance_settings_updated_by
        FOREIGN KEY (updated_by_id)
        REFERENCES users(id)
        ON DELETE NO ACTION
        ON UPDATE NO ACTION
);

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

DROP TABLE dbo.governance_settings;

-- +goose StatementEnd