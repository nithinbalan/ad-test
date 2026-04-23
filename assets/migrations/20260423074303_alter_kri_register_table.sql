-- +goose Up
-- +goose StatementBegin
ALTER TABLE kri_register
ADD
    submitted_by INT NULL,
    submitted_at DATETIME NULL,
    next_review_level INT NULL,
    workflow_status_id INT NULL,
    approved_by INT NULL,
    approved_at DATETIME NULL,

    CONSTRAINT FK_kri_register_users_submitted_by
        FOREIGN KEY (submitted_by)
        REFERENCES users(id)
        ON DELETE NO ACTION
        ON UPDATE NO ACTION,

    CONSTRAINT FK_kri_register_workflow_statuses
        FOREIGN KEY (workflow_status_id)
        REFERENCES workflow_statuses(id)
        ON DELETE NO ACTION
        ON UPDATE NO ACTION,

    CONSTRAINT FK_kri_register_users_approved_by
        FOREIGN KEY (approved_by)
        REFERENCES users(id)
        ON DELETE NO ACTION
        ON UPDATE NO ACTION;
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
ALTER TABLE kri_register
DROP CONSTRAINT FK_kri_register_users_submitted_by;

ALTER TABLE kri_register
DROP CONSTRAINT FK_kri_register_workflow_statuses;

ALTER TABLE kri_register
DROP CONSTRAINT FK_kri_register_users_approved_by;

ALTER TABLE kri_register
DROP COLUMN submitted_by,
            submitted_at,
            next_review_level,
            workflow_status_id,
            approved_by,
            approved_at;
-- +goose StatementEnd