-- +goose Up
-- +goose StatementBegin

CREATE TABLE risk_register_update_plan_statuses (
    id INT PRIMARY KEY IDENTITY(1, 1),
    title_en NVARCHAR(500) NOT NULL,
    title_ar NVARCHAR(500) NOT NULL,
    color_code VARCHAR(25) NOT NULL,
    type VARCHAR(100) NOT NULL,
        CONSTRAINT CHK_risk_register_update_plan_statuses_type CHECK (
            type IN (
                'not_started',
                'in_progress',
                'overdue',
                'completed'
            )
        ),
    status VARCHAR(10) NOT NULL DEFAULT 'active',
        CONSTRAINT CHK_risk_register_update_plan_statuses_status CHECK (status IN ('active', 'inactive'))
);

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

DROP TABLE risk_register_update_plan_statuses;

-- +goose StatementEnd


