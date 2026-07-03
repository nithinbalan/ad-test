-- +goose Up
-- +goose StatementBegin
CREATE TABLE bcp_process_external_dependencies (
    id INT PRIMARY KEY IDENTITY(1, 1),
    bcp_process_id INT NOT NULL,
    entity_name NVARCHAR(1000) NOT NULL,
    entity_type_id INT NULL,
    needed_support_description NVARCHAR(MAX) NULL,
    dependency_tier_id INT NULL,
    sla_register_id INT NULL,
    vendor_rto INT NULL,
    vendor_rto_unit_id INT NULL,
    internal_stakeholder NVARCHAR(MAX) NULL,
    contact_number NVARCHAR(100) NULL,
    contact_email NVARCHAR(255) NULL,
    created_by_id INT NULL,
    created_at DATETIME NOT NULL,

    CONSTRAINT FK_bcp_process_external_dependencies_process FOREIGN KEY (bcp_process_id) REFERENCES bcp_processes(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_bcp_process_external_dependencies_entity_type FOREIGN KEY (entity_type_id) REFERENCES process_external_dependency_entity_types(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_bcp_process_external_dependencies_tier FOREIGN KEY (dependency_tier_id) REFERENCES process_dependency_tiers(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_bcp_process_external_dependencies_vendor_rto_unit FOREIGN KEY (vendor_rto_unit_id) REFERENCES business_continuity_duration_units(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_bcp_process_external_dependencies_sla FOREIGN KEY (sla_register_id) REFERENCES sla_register(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_bcp_process_external_dependencies_users_1 FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE bcp_process_external_dependencies;
-- +goose StatementEnd
