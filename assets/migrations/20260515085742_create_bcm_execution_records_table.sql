-- +goose Up
-- +goose StatementBegin

CREATE TABLE bcm_execution_records (
    id INT PRIMARY KEY IDENTITY(1, 1),
    title_en VARCHAR(1000) NOT NULL,  
    title_ar VARCHAR(1000) NOT NULL,  
    owner_entity_id INT NOT NULL,
    record_type_id INT NULL,
    is_back_up_available BIT DEFAULT 0,
    location_id INT NULL,
    recovery_plan NVARCHAR(MAX) NULL,
    created_by_id INT,
    updated_by_id INT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,

    CONSTRAINT FK_bcm_execution_records_organization_units_1 FOREIGN KEY (owner_entity_id) REFERENCES organization_units (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_bcm_execution_records_bcm_execution_record_types_1 FOREIGN KEY (record_type_id) REFERENCES bcm_execution_record_types (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_bcm_execution_records_locations_1 FOREIGN KEY (location_id) REFERENCES locations (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_bcm_execution_records_users_1 FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_bcm_execution_records_users_2 FOREIGN KEY (updated_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION  
    
);

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DROP TABLE bcm_execution_records;
-- +goose StatementEnd