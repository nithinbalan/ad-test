-- +goose Up
-- +goose StatementBegin

CREATE TABLE bcm_impact_analysis_levels (
    id INT PRIMARY KEY IDENTITY(1, 1),
    title_en VARCHAR(45) NOT NULL,
    title_ar NVARCHAR(45) NOT NULL,
    score INT NOT NULL,
    color_code VARCHAR(25) NOT NULL,
    type VARCHAR(20) NOT NULL,
        CONSTRAINT CHK_bcm_impact_analysis_levels_type CHECK (
            type IN (
                'non_applicable',
                'very_low',
                'low',
                'moderate',
                'high',
                'very_high'
            )
        ),
    status VARCHAR(10) NOT NULL DEFAULT 'active',
        CONSTRAINT CHK_bcm_impact_analysis_levels_status CHECK (status IN ('active', 'inactive'))
);

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

DROP TABLE dbo.bcm_impact_analysis_levels;

-- +goose StatementEnd
