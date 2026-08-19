-- +goose Up
-- +goose StatementBegin
CREATE TABLE cs_log_soc_room_access_logs (
    id INT PRIMARY KEY IDENTITY(1, 1),
    reference_number NVARCHAR(45) NOT NULL,
    visitor_full_name NVARCHAR(255) NOT NULL,
    entry_time DATETIME NOT NULL,
    exit_time DATETIME NOT NULL,
    visitor_id_number NVARCHAR(100) NOT NULL,
    department_id INT NOT NULL,
    external_party_name NVARCHAR(255) NOT NULL,
    purpose_of_entry NVARCHAR(MAX) NOT NULL,
    authorizing_officer_id INT NOT NULL,
    notes NVARCHAR(MAX) NULL,
    created_by_id INT NOT NULL,
    updated_by_id INT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,

    CONSTRAINT UQ_cs_log_soc_room_access_logs_reference_number UNIQUE (reference_number),

    CONSTRAINT FK_cs_log_soc_room_access_logs_organization_units
        FOREIGN KEY (department_id) REFERENCES organization_units(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_cs_log_soc_room_access_logs_users_1
        FOREIGN KEY (authorizing_officer_id) REFERENCES users(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_cs_log_soc_room_access_logs_users_2
        FOREIGN KEY (created_by_id) REFERENCES users(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_cs_log_soc_room_access_logs_users_3
        FOREIGN KEY (updated_by_id) REFERENCES users(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE cs_log_soc_room_access_logs;
-- +goose StatementEnd
