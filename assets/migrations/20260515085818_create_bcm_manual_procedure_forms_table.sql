-- +goose Up
-- +goose StatementBegin

CREATE TABLE bcm_manual_procedure_forms (
    id INT PRIMARY KEY IDENTITY(1, 1),
    title_en VARCHAR(1000) NOT NULL,
    title_ar VARCHAR(1000) NOT NULL,
    code VARCHAR(1000) NULL,
    owner_id INT NOT NULL,
    location_id INT NULL,
    stakeholder_id INT NULL,
    notes NVARCHAR(MAX) NULL,
    created_by_id INT,
    updated_by_id INT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,

    CONSTRAINT FK_bcm_manual_procedure_forms_users_1 FOREIGN KEY (owner_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_bcm_manual_procedure_forms_stakeholders_1 FOREIGN KEY (stakeholder_id) REFERENCES stakeholders (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_bcm_manual_procedure_forms_locations_1 FOREIGN KEY (location_id) REFERENCES locations (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_bcm_manual_procedure_forms_users_2 FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_bcm_manual_procedure_forms_users_3 FOREIGN KEY (updated_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION  
    
);

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DROP TABLE bcm_manual_procedure_forms;
-- +goose StatementEnd