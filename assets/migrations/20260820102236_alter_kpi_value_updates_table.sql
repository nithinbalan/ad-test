-- +goose Up
-- +goose StatementBegin
ALTER TABLE kpi_value_updates ALTER COLUMN value DECIMAL(18, 2) NULL;
ALTER TABLE kpi_value_updates ALTER COLUMN reviewed_date DATE NULL;

ALTER TABLE kpi_value_updates
ADD kpi_review_status_id INT NULL,
    start_date DATE NULL,
    end_date DATE NULL,
    period_label NVARCHAR(50) NULL,
    corrective_action NVARCHAR(MAX) NULL,
    approved_at DATETIME NULL,
    reverted_at DATETIME NULL;
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE kpi_value_updates
ADD CONSTRAINT FK_kpi_value_updates_kpi_review_statuses
    FOREIGN KEY (kpi_review_status_id) REFERENCES kpi_review_statuses(id) ON DELETE NO ACTION ON UPDATE NO ACTION;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE kpi_value_updates DROP CONSTRAINT FK_kpi_value_updates_kpi_review_statuses;
ALTER TABLE kpi_value_updates DROP COLUMN
    kpi_review_status_id, start_date, end_date, period_label,
    corrective_action, approved_at, reverted_at;
ALTER TABLE kpi_value_updates ALTER COLUMN value INT NULL;
ALTER TABLE kpi_value_updates ALTER COLUMN reviewed_date DATE NOT NULL;
-- +goose StatementEnd
