-- +goose Up
-- +goose StatementBegin

CREATE TABLE bcp_process_internal_dependencies (
    id INT PRIMARY KEY IDENTITY(1, 1),
    bcp_process_id INT NOT NULL,
    supporting_unit_id INT NOT NULL,
    dependency_tier_id INT NOT NULL,
    dependency_percentage_from INT NULL,
    dependency_percentage_to INT NULL,
    required_support NVARCHAR(MAX) NULL,
    notes NVARCHAR(MAX) NULL,
    created_by_id INT NULL,
    created_at DATETIME NOT NULL,

    CONSTRAINT FK_bcp_process_internal_dependencies_process FOREIGN KEY (bcp_process_id) REFERENCES bcp_processes(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_bcp_process_internal_dependencies_org_unit FOREIGN KEY (supporting_unit_id) REFERENCES organization_units(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_bcp_process_internal_dependencies_tier FOREIGN KEY (dependency_tier_id) REFERENCES process_dependency_tiers(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_bcp_process_internal_dependencies_users_1 FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

DROP TABLE bcp_process_internal_dependencies;

-- +goose StatementEnd
