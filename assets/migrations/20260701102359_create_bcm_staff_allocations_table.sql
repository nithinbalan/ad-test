-- +goose Up
-- +goose StatementBegin
CREATE TABLE bcm_staff_allocations (
    id INT IDENTITY(1,1) PRIMARY KEY,
    process_id INT NOT NULL,
    process_employee_analysis_id INT NULL,
    created_by_id INT,
    updated_by_id INT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,   

    CONSTRAINT FK_bcm_staff_allocations_process FOREIGN KEY (process_id)
        REFERENCES processes(id) ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_bcm_staff_allocations_process_employee_analysis FOREIGN KEY (process_employee_analysis_id)
        REFERENCES process_employee_analysis(id) ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_bcm_staff_allocations_users_1 FOREIGN KEY (created_by_id)
        REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    
    CONSTRAINT FK_bcm_staff_allocations_users_2 FOREIGN KEY (updated_by_id)
        REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE bcm_staff_allocations;
-- +goose StatementEnd