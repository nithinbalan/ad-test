-- +goose Up
-- +goose StatementBegin
CREATE TABLE audit_plan_commencement_letter_cc_users (
    id INT PRIMARY KEY IDENTITY(1,1),
    audit_plan_commencement_letter_id INT NOT NULL,
    user_id INT NOT NULL,
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,
    updated_by_id INT NULL,
    updated_at DATETIME NULL,

    CONSTRAINT FK_audit_plan_commencement_letter_cc_users_audit_plan_commencement_letters FOREIGN KEY (audit_plan_commencement_letter_id) REFERENCES audit_plan_commencement_letters(id) ON DELETE CASCADE ON UPDATE NO ACTION,
    CONSTRAINT FK_audit_plan_commencement_letter_cc_users_users FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_audit_plan_commencement_letter_cc_users_created_by FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_audit_plan_commencement_letter_cc_users_updated_by FOREIGN KEY (updated_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.audit_plan_commencement_letter_cc_users;
-- +goose StatementEnd
