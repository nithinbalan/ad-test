-- +goose Up
-- +goose StatementBegin

CREATE TABLE bia_processes (
    id INT PRIMARY KEY IDENTITY(1,1),
    process_id INT NOT NULL,
    workflow_status_id INT NULL,
    next_review_level INT NULL,
    submitted_by_id INT NULL,
    submitted_at DATETIME NULL,
    approved_by_id INT NULL,
    approved_at DATETIME NULL,
    created_by_id INT NULL,
    updated_by_id INT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL

    CONSTRAINT FK_bia_processes_process
        FOREIGN KEY (process_id)
        REFERENCES processes(id)
        ON DELETE NO ACTION
        ON UPDATE NO ACTION, 
     
     CONSTRAINT FK_bia_processes_users_submitted_by
        FOREIGN KEY (submitted_by_id)
        REFERENCES users(id)
        ON DELETE NO ACTION
        ON UPDATE NO ACTION,

    CONSTRAINT FK_bia_processes_workflow_statuses
        FOREIGN KEY (workflow_status_id)
        REFERENCES workflow_statuses(id)
        ON DELETE NO ACTION
        ON UPDATE NO ACTION,

    CONSTRAINT FK_bia_processes_users_approved_by
        FOREIGN KEY (approved_by_id)
        REFERENCES users(id)
        ON DELETE NO ACTION
        ON UPDATE NO ACTION,
        
        CONSTRAINT FK_bia_processes_users_1
        FOREIGN KEY (created_by_id)
        REFERENCES users(id)
        ON DELETE NO ACTION
        ON UPDATE NO ACTION,

    CONSTRAINT FK_bia_processes_users_2
        FOREIGN KEY (updated_by_id)
        REFERENCES users(id)
        ON DELETE NO ACTION
        ON UPDATE NO ACTION
);

-- +goose StatementEnd

    
-- +goose Down
-- +goose StatementBegin

DROP TABLE bia_processes;

-- +goose StatementEnd