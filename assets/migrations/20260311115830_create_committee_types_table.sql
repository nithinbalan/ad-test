-- +goose Up
-- +goose StatementBegin

CREATE TABLE committee_types (
    id INT PRIMARY KEY IDENTITY(1, 1),
    title_en NVARCHAR(500) NOT NULL,
    title_ar NVARCHAR(500) NOT NULL,
    type VARCHAR(100) NOT NULL,
        CONSTRAINT CHK_committee_types_type CHECK (
            type IN (
                'statutory_reference',
                'no_statutory_reference'
            )
        ),
    status VARCHAR(10) NOT NULL DEFAULT 'active',
        CONSTRAINT CHK_committee_types_status CHECK (status IN ('active', 'inactive'))
);

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

DROP TABLE dbo.committee_types;

-- +goose StatementEnd
