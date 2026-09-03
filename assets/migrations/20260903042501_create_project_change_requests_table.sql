-- +goose Up
-- +goose StatementBegin
CREATE TABLE project_change_requests (
    id INT PRIMARY KEY IDENTITY(1,1),
    reference_number VARCHAR(50) NOT NULL,
    project_id INT NOT NULL,
    title NVARCHAR(255) NOT NULL,
    description NVARCHAR(MAX) NULL,
    business_justification NVARCHAR(MAX) NULL,
    anticipated_benefits NVARCHAR(MAX) NULL,
    risks_of_deferral NVARCHAR(MAX) NULL,
    requested_by_user_id INT NOT NULL,
    requested_on DATE NOT NULL,
    workflow_status_id INT NULL,
    submitted_by_id INT NULL,
    submitted_at DATETIME NULL,
    next_review_level INT NULL,
    approved_by_id INT NULL,
    approved_at DATETIME NULL,
    rejected_by_id INT NULL,
    rejected_at DATETIME NULL,
    created_by_id INT NOT NULL,
    updated_by_id INT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,

    CONSTRAINT FK_project_change_requests_projects FOREIGN KEY (project_id) 
        REFERENCES projects(id) ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_project_change_requests_workflow_status FOREIGN KEY (workflow_status_id) 
        REFERENCES workflow_statuses(id) ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_project_change_requests_users_1 FOREIGN KEY (requested_by_user_id) 
        REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_project_change_requests_users_2 FOREIGN KEY (submitted_by_id) 
        REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_project_change_requests_users_3 FOREIGN KEY (approved_by_id) 
        REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_project_change_requests_users_4 FOREIGN KEY (rejected_by_id) 
        REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_project_change_requests_users_5 FOREIGN KEY (created_by_id) 
        REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_project_change_requests_users_6 FOREIGN KEY (updated_by_id) 
        REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE project_change_requests;
-- +goose StatementEnd
