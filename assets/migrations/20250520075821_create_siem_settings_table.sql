-- +goose Up
-- +goose StatementBegin
CREATE TABLE siem_settings (
    id INT PRIMARY KEY IDENTITY(1,1),
    integration_id INT NULL,
    syslog_server NVARCHAR(255) NOT NULL,
    syslog_port INT NOT NULL,
    syslog_protocol NVARCHAR(10) NOT NULL,
    syslog_timezone NVARCHAR(50) NULL,
    syslog_datetime_format NVARCHAR(50) NULL,
    updated_by_id INT NULL,
    updated_at DATETIME NULL,

    CONSTRAINT FK_siem_settings_updated_by FOREIGN KEY (updated_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_siem_settings_integrations FOREIGN KEY (integration_id) REFERENCES integrations(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE siem_settings;
-- +goose StatementEnd
