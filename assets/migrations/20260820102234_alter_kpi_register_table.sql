-- +goose Up
-- +goose StatementBegin
ALTER TABLE kpi_register
ADD identifier VARCHAR(100) NULL,
    type VARCHAR(20) NOT NULL CONSTRAINT DF_kpi_register_type DEFAULT 'manual',
    kpi_review_status_id INT NULL,
    previous_value DECIMAL(18, 2) NULL,
    status_id INT NULL,
    target_date DATE NULL,
    kpi_achievement_logic_id INT NULL,
    target_description NVARCHAR(500) NULL,
    review_user_id INT NULL,
    is_kri BIT NOT NULL
        CONSTRAINT DF_kpi_register_is_kri DEFAULT 0;
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE kpi_register
ADD CONSTRAINT CHK_kpi_register_type CHECK ([type] IN ('automated', 'manual')),
    CONSTRAINT FK_kpi_register_kpi_review_statuses
        FOREIGN KEY (kpi_review_status_id) REFERENCES kpi_review_statuses(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_kpi_register_lifecycle_statuses
        FOREIGN KEY (status_id) REFERENCES lifecycle_statuses(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_kpi_register_kpi_achievement_logics
        FOREIGN KEY (kpi_achievement_logic_id) REFERENCES kpi_achievement_logics(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_kpi_register_users_4
        FOREIGN KEY (review_user_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION;

ALTER TABLE kpi_register ALTER COLUMN actual_value DECIMAL(18, 2) NULL;
ALTER TABLE kpi_register ALTER COLUMN achievement DECIMAL(18, 2) NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE kpi_register DROP CONSTRAINT FK_kpi_register_users_4;
ALTER TABLE kpi_register DROP CONSTRAINT FK_kpi_register_kpi_achievement_logics;
ALTER TABLE kpi_register DROP CONSTRAINT FK_kpi_register_lifecycle_statuses;
ALTER TABLE kpi_register DROP CONSTRAINT FK_kpi_register_kpi_review_statuses;
ALTER TABLE kpi_register DROP CONSTRAINT CHK_kpi_register_type;
ALTER TABLE kpi_register DROP CONSTRAINT DF_kpi_register_is_kri;
ALTER TABLE kpi_register DROP CONSTRAINT DF_kpi_register_type;

ALTER TABLE kpi_register ALTER COLUMN actual_value FLOAT NULL;
ALTER TABLE kpi_register ALTER COLUMN achievement FLOAT NULL;

ALTER TABLE kpi_register DROP COLUMN
    identifier, type, kpi_review_status_id, previous_value, status_id,
    target_date, kpi_achievement_logic_id, target_description, review_user_id,
    is_kri;
-- +goose StatementEnd
