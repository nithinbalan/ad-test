-- +goose Up
-- +goose StatementBegin
CREATE TABLE audit_plan_commencement_letters (
    id INT PRIMARY KEY IDENTITY(1,1),
    audit_plan_id INT NOT NULL,
    responsible_unit_id INT NOT NULL,
    to_email_ids NVARCHAR(MAX) NULL,
    cc_email_ids NVARCHAR(MAX) NULL,
    subject NVARCHAR(255) NOT NULL,
    body NVARCHAR(MAX) NOT NULL,
    from_user_id INT NOT NULL,
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,
    updated_by_id INT NULL,
    updated_at DATETIME NULL,

    CONSTRAINT FK_audit_plan_commencement_letters_audit_plans FOREIGN KEY (audit_plan_id) REFERENCES audit_plans(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_audit_plan_commencement_letters_responsible_units FOREIGN KEY (responsible_unit_id) REFERENCES units(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_audit_plan_commencement_letters_users FOREIGN KEY (from_user_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_audit_plan_commencement_letters_created_by FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_audit_plan_commencement_letters_updated_by FOREIGN KEY (updated_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.audit_plan_commencement_letters;
-- +goose StatementEnd
