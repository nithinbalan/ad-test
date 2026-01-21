-- +goose Up
-- +goose StatementBegin
CREATE TABLE consultation_requests
(
    id INT IDENTITY(1,1) PRIMARY KEY,
    reference_number VARCHAR(45) NOT NULL UNIQUE, -- CONS-12323
    date_raised DATE NOT NULL,
    requesting_department_unit_id INT NOT NULL,
    consulted_department_unit_id INT NOT NULL,
    responsible_user_id INT NULL,
    request_details NVARCHAR(MAX) NOT NULL,
    consultation_request_status_id INT NOT NULL,
    reply NVARCHAR(MAX) NULL,
    reply_date DATE NULL,
    notes NVARCHAR(MAX) NULL,
    created_by_id INT NOT NULL,
    updated_by_id INT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,

    -- Foreign Keys
    CONSTRAINT FK_consultation_requests_department_unit
        FOREIGN KEY (requesting_department_unit_id)
        REFERENCES organization_units (id)
        ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_consultation_requests_consulted_department_unit
        FOREIGN KEY (consulted_department_unit_id )
        REFERENCES organization_units (id)
        ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_consultation_requests_responsible_users
        FOREIGN KEY (responsible_user_id)
        REFERENCES users (id)
        ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_consultation_requests_statuses
        FOREIGN KEY (consultation_request_status_id)
        REFERENCES consultation_request_statuses (id)
        ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_consultation_requests_created_by_users
        FOREIGN KEY (created_by_id)
        REFERENCES users (id)
        ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_consultation_requests_updated_by_users
        FOREIGN KEY (updated_by_id)
        REFERENCES users (id)
        ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE consultation_requests;
-- +goose StatementEnd
