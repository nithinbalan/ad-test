-- +goose Up
-- +goose StatementBegin

CREATE TABLE business_continuity_duration_units (
    id INT PRIMARY KEY IDENTITY(1, 1),
    title_en NVARCHAR(500) NOT NULL,
    title_ar NVARCHAR(500) NOT NULL,
    minutes INT NOT NULL,
    type VARCHAR(100) NOT NULL,
        CONSTRAINT CHK_business_continuity_duration_units_type CHECK (
            type IN (
                'minutes',
                'hours',
                'days',
                'weeks'
            )
        ),
    status VARCHAR(10) NOT NULL DEFAULT 'active',
        CONSTRAINT CHK_business_continuity_duration_units_status CHECK (status IN ('active', 'inactive'))
);

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

DROP TABLE dbo.business_continuity_duration_units;

-- +goose StatementEnd
