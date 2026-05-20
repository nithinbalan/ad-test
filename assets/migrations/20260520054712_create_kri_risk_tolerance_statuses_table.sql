-- +goose Up
-- +goose StatementBegin

CREATE TABLE kri_risk_tolerance_statuses (
    id INT PRIMARY KEY IDENTITY(1, 1),
    title_en NVARCHAR(500) NOT NULL,
    title_ar NVARCHAR(500) NOT NULL,
    color_code VARCHAR(25) NOT NULL,
    type VARCHAR(100) NOT NULL,
        CONSTRAINT CHK_kri_risk_tolerance_statuses_type CHECK (
            type IN (
                'within_minimum_susceptibility',
                'within_maximum_susceptibility',
                'at_maximum_tolerable',
                'exceeded_maximum_tolerable'
            )
        ),
    status VARCHAR(10) NOT NULL DEFAULT 'active',
        CONSTRAINT CHK_kri_risk_tolerance_statuses_status CHECK (status IN ('active', 'inactive'))
);

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

DROP TABLE kri_risk_tolerance_statuses;

-- +goose StatementEnd
