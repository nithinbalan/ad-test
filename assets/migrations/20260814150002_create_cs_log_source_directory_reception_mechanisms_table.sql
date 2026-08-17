-- +goose Up
-- +goose StatementBegin
CREATE TABLE cs_log_source_directory_reception_mechanisms (
    id INT PRIMARY KEY IDENTITY(1, 1),
    source_directory_id INT NOT NULL,
    reception_mechanism_id INT NOT NULL,

    CONSTRAINT UQ_cs_log_source_directory_reception_mechanisms_pair
        UNIQUE (source_directory_id, reception_mechanism_id),

    CONSTRAINT FK_cs_log_source_directory_reception_mechanisms_cs_log_source_directory
        FOREIGN KEY (source_directory_id) REFERENCES cs_log_source_directory(id)
        ON DELETE CASCADE ON UPDATE NO ACTION,

    CONSTRAINT FK_cs_log_source_directory_reception_mechanisms_cs_log_data_reception_mechanisms
        FOREIGN KEY (reception_mechanism_id) REFERENCES cs_log_data_reception_mechanisms(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE cs_log_source_directory_reception_mechanisms;
-- +goose StatementEnd
