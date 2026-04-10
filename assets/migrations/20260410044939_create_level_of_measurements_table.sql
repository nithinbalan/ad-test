-- +goose Up
-- +goose StatementBegin
CREATE TABLE level_of_measurements (
    id INT PRIMARY KEY IDENTITY(1,1),
    title_en NVARCHAR(100) NOT NULL,
    title_ar NVARCHAR(100) NOT NULL,
    type VARCHAR(50) NOT NULL,
        CONSTRAINT CHK_level_of_measurements_type 
        CHECK (type IN ('board_of_directors', 'grc_committee')),
    status VARCHAR(10) NOT NULL DEFAULT 'active',
        CONSTRAINT CHK_level_of_measurements_status 
        CHECK (status IN ('active', 'inactive'))
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE level_of_measurements;
-- +goose StatementEnd