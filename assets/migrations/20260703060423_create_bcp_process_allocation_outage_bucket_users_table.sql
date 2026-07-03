-- +goose Up
-- +goose StatementBegin
CREATE TABLE bcp_process_allocation_outage_bucket_users (
    id INT IDENTITY(1,1) PRIMARY KEY,
    bcp_process_allocation_outage_bucket_id INT NOT NULL,
    user_id INT NULL,
    is_remote BIT NOT NULL DEFAULT 0,
    created_by_id INT,
    created_at DATETIME NOT NULL, 

    CONSTRAINT FK_bcp_process_allocation_outage_bucket_users_outage_bucket FOREIGN KEY (bcp_process_allocation_outage_bucket_id)
        REFERENCES bcp_process_allocation_outage_buckets(id) ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_bcp_process_allocation_outage_bucket_users_user FOREIGN KEY (user_id)
        REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_bcp_process_allocation_outage_bucket_users_created_by FOREIGN KEY (created_by_id)
        REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE bcp_process_allocation_outage_bucket_users;
-- +goose StatementEnd