-- +goose Up
-- +goose StatementBegin
CREATE TABLE consultation_request_comments
(
    id INT IDENTITY(1,1) PRIMARY KEY,


    consultation_request_id INT NOT NULL,
    consultation_request_status_id INT NOT NULL,
    comment NVARCHAR(MAX) NOT NULL,


    -- Audit Fields
    created_at DATETIME NOT NULL,
    created_by_id INT NOT NULL,
    updated_at DATETIME NULL,
    updated_by_id INT NULL,


    -- Foreign Keys
    CONSTRAINT FK_crc_consultation_requests
        FOREIGN KEY (consultation_request_id)
        REFERENCES consultation_requests (id)
        ON UPDATE NO ACTION
        ON DELETE NO ACTION,


    CONSTRAINT FK_crc_consultation_request_statuses
        FOREIGN KEY (consultation_request_status_id)
        REFERENCES consultation_request_statuses (id)
        ON UPDATE NO ACTION
        ON DELETE NO ACTION,


    CONSTRAINT FK_crc_created_by_users
        FOREIGN KEY (created_by_id)
        REFERENCES users (id)
        ON UPDATE NO ACTION
        ON DELETE NO ACTION,


    CONSTRAINT FK_crc_updated_by_users
        FOREIGN KEY (updated_by_id)
        REFERENCES users (id)
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
);
-- +goose StatementEnd




-- +goose Down
-- +goose StatementBegin
DROP TABLE consultation_request_comments;
-- +goose StatementEnd
