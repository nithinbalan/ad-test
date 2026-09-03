-- +goose Up
-- +goose StatementBegin
CREATE TABLE planned_advisory_security_classifications (
    id INT PRIMARY KEY IDENTITY(1, 1),
    title_en NVARCHAR(MAX) NOT NULL,
    title_ar NVARCHAR(MAX) NOT NULL,
    status VARCHAR(10) NOT NULL DEFAULT 'active',
        CONSTRAINT CHK_planned_advisory_security_classifications_status CHECK (status IN ('active', 'inactive'))
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE planned_advisory_security_classifications;
-- +goose StatementEnd