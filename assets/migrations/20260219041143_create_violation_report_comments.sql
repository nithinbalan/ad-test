-- +goose Up
-- +goose StatementBegin

CREATE TABLE violation_report_comments
(
    id INT IDENTITY(1,1) PRIMARY KEY,

    violation_report_id INT NOT NULL,
    violation_report_status_id INT NULL,
    comment NVARCHAR(MAX) NOT NULL,

    created_at DATETIME NOT NULL,
    created_by_id INT NOT NULL,
    updated_at DATETIME NULL,
    updated_by_id INT NULL,

    -- Foreign Keys
    CONSTRAINT FK_vrc_violation_reports
        FOREIGN KEY (violation_report_id)
        REFERENCES violation_reports (id)
        ON UPDATE NO ACTION
        ON DELETE NO ACTION,

    CONSTRAINT FK_vrc_violation_report_statuses
        FOREIGN KEY (violation_report_status_id)
        REFERENCES violation_report_statuses (id)
        ON UPDATE NO ACTION
        ON DELETE NO ACTION,

    CONSTRAINT FK_vrc_created_by_users
        FOREIGN KEY (created_by_id)
        REFERENCES users (id)
        ON UPDATE NO ACTION
        ON DELETE NO ACTION,

    CONSTRAINT FK_vrc_updated_by_users
        FOREIGN KEY (updated_by_id)
        REFERENCES users (id)
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE violation_report_comments;

-- +goose StatementEnd
