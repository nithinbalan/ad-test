-- +goose Up
-- +goose StatementBegin
CREATE TABLE units (
    id INT PRIMARY KEY IDENTITY(1,1),
    title_en NVARCHAR(100) NOT NULL,
    title_ar NVARCHAR(100) NOT NULL,
    code_en NVARCHAR(100) NOT NULL,
    code_ar NVARCHAR(100) NOT NULL,
    type VARCHAR(20) NOT NULL DEFAULT 'other',
       CONSTRAINT CHK_units_type CHECK (type IN ('length', 'area', 'volume', 'mass', 'time', 'temperature', 'speed', 'pressure', 'energy', 'power', 'frequency', 'electrical', 'light', 'data','other')),
    color_code VARCHAR(10) NULL,
    status VARCHAR(10) NOT NULL DEFAULT 'active',
    CONSTRAINT CHK_units_status CHECK (status IN ('active', 'inactive'))
   
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE units;
-- +goose StatementEnd
