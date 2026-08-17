-- +goose Up
-- +goose StatementBegin
CREATE TABLE cs_log_source_directory_communication_channels (
    id INT PRIMARY KEY IDENTITY(1, 1),
    source_directory_id INT NOT NULL,
    communication_channel_id INT NOT NULL,

    CONSTRAINT UQ_cs_log_source_directory_communication_channels_pair
        UNIQUE (source_directory_id, communication_channel_id),

    CONSTRAINT FK_cs_log_source_directory_communication_channels_cs_log_source_directory
        FOREIGN KEY (source_directory_id) REFERENCES cs_log_source_directory(id)
        ON DELETE CASCADE ON UPDATE NO ACTION,

    CONSTRAINT FK_cs_log_source_directory_communication_channels_cs_log_communication_channels
        FOREIGN KEY (communication_channel_id) REFERENCES cs_log_communication_channels(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE cs_log_source_directory_communication_channels;
-- +goose StatementEnd
