-- +goose Up
-- +goose StatementBegin
CREATE TABLE action_plan_statuses (
    id INT PRIMARY KEY IDENTITY(1, 1),
    title_en NVARCHAR(100) NOT NULL,
    title_ar NVARCHAR(100) NOT NULL,
    color_code VARCHAR(25) NOT NULL,
    type VARCHAR(20) NOT NULL DEFAULT 'new',
        CONSTRAINT CHK_action_plan_statuses_type CHECK (type IN ('new', 'in_progress', 'completed', 'on_hold', 'terminated')),
    status VARCHAR(10) NOT NULL DEFAULT 'active',
        CONSTRAINT CHK_action_plan_statuses_status CHECK (status IN ('active', 'inactive'))
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.action_plan_statuses;
-- +goose StatementEnd
